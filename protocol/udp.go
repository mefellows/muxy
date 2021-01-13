// See https://ops.tips/blog/udp-client-and-server-in-go/#tcp-dialing-vs-udp-dialing and https://gist.github.com/mike-zhang/3853251

package protocol

import (
	"fmt"
	"net"
	"sync"

	"github.com/mefellows/muxy/log"
	"github.com/mefellows/muxy/muxy"
	"github.com/mefellows/plugo/plugo"
)

// UDPProxy implements a UDP proxy as a top-level Muxy protocol
type UDPProxy struct {
	Port       int    `required:"true"`
	Host       string `required:"true" default:"localhost"`
	ProxyHost  string `required:"true" mapstructure:"proxy_host"`
	ProxyPort  int    `required:"true" mapstructure:"proxy_port"`
	HexOutput  bool   `mapstructure:"hex_output"`
	PacketSize int    `mapstructure:"packet_size" default:"64" required:"true"`
	connID     uint64
	middleware []muxy.Middleware
	proxy      *udpProxy
}

func init() {
	plugo.PluginFactories.Register(func() (interface{}, error) {
		return &UDPProxy{}, nil
	}, "udp_proxy")
}

// Setup the UDP proxy
func (p *UDPProxy) Setup(middleware []muxy.Middleware) {
	p.middleware = middleware
}

// Teardown the UDP proxy
func (p *UDPProxy) Teardown() {
	log.Info("UDP Proxy closed (%d bytes sent, %d bytes received)", p.proxy.sentBytes, p.proxy.receivedBytes)
}

// Proxy runs the UDP proxy
func (p *UDPProxy) Proxy() {
	log.Trace("Checking connection: %s:%d", p.Host, p.Port)
	laddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", p.Host, p.Port))
	check("udp", err)
	raddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", p.ProxyHost, p.ProxyPort))
	check("udp", err)
	conn, err := net.ListenUDP("udp", laddr)
	check("udp", err)
	p.connID++

	log.Info("UDP Proxy proxy listening on %s", log.Colorize(log.BLUE, fmt.Sprintf("udp://%s:%d", p.Host, p.Port)))

	proxy := &udpProxy{
		lconn:      conn,
		laddr:      laddr,
		raddr:      raddr,
		packetSize: p.PacketSize,
		erred:      false,
		prefix:     fmt.Sprintf("Connection #%03d ", p.connID),
		hex:        p.HexOutput,
		middleware: p.middleware,
		clients:    make(connectionMap),
		clientLock: new(sync.Mutex),
	}

	p.proxy = proxy

	proxy.start()
}

//A udpProxy represents connections and their state
type udpProxy struct {
	middleware    []muxy.Middleware
	sentBytes     uint64
	receivedBytes uint64
	laddr, raddr  *net.UDPAddr
	lconn         *net.UDPConn
	protocol      string
	erred         bool
	prefix        string
	hex           bool
	packetSize    int
	clients       connectionMap
	clientLock    *sync.Mutex
}

func (p *udpProxy) start() {
	log.Trace("UDP Proxy Starting UDP Proxy")

	p.listen()
	p.lconn.Close()

	log.Info("UDP Proxy closed (%d bytes sent, %d bytes received)", p.sentBytes, p.receivedBytes)
}

// Information maintained for each client/server connection
type connection struct {
	clientAddr *net.UDPAddr // Address of the client initiating the request
	serverConn *net.UDPConn // UDP connection to the proxied server
}

type connectionMap map[string]*connection

// Mutex used to serialize access to the dictionary
var connectionLocker *sync.Mutex = new(sync.Mutex)

// Handle errors
func handleConnectionErr(err error) bool {
	if err == nil {
		return false
	}
	log.Info("Error: %s", err.Error())
	return true
}

// Go routine which manages connection from server to single client
func (p *udpProxy) waitForServer(conn *connection) {
	var buffer = make([]byte, p.packetSize)

	for {
		log.Info("waiting for remote server's (%s) response", conn.serverConn.RemoteAddr().String())
		// Read from proxied server
		n, err := conn.serverConn.Read(buffer[0:])
		if handleConnectionErr(err) {
			log.Error("error reading from server: %s", err)
			continue
		}

		// Extract message
		b := buffer[0:n]
		log.Info("response from server: %b", b)

		// Apply middleware to modify the remote server response
		ctx := &muxy.Context{Bytes: b}
		for _, middleware := range p.middleware {
			log.Trace("UDP Proxy applying middleware %v", middleware)
			middleware.HandleEvent(muxy.EventPostDispatch, ctx)
			log.Trace("UDP Proxy overwriting bytes sent back to originating client: %s", ctx.Bytes)
			b = ctx.Bytes
		}

		// Relay it to client
		n, err = p.lconn.WriteToUDP(b, conn.clientAddr)
		p.sentBytes += uint64(n)

		log.Info("sending back to client %s", conn.clientAddr.String())
		if handleConnectionErr(err) {
			log.Error("error relaying to client: %s", err)
			continue
		}
		log.Info("Relayed '%s' from server to %s.\n", string(b), conn.clientAddr.String())
	}
}

// Generate a new connection by opening a UDP connection to the server
func (p *udpProxy) newConnection(serverAddress, clientAddress *net.UDPAddr) *connection {
	serverConn, err := net.DialUDP("udp", nil, serverAddress)

	if handleConnectionErr(err) {
		return nil
	}

	conn := &connection{
		clientAddr: clientAddress,
		serverConn: serverConn,
	}

	return conn
}

func (p *udpProxy) listen() {
	buff := make([]byte, p.packetSize)

	for {

		// 1. Wait for a new connection from a client
		n, clientaddr, err := p.lconn.ReadFromUDP(buff)
		if err != nil {
			fmt.Printf("Some error reading from client UDP  %v", err)
			continue
		}

		fmt.Printf("Read a message from %v %v \n", clientaddr, p)
		b := buff[:n]
		p.receivedBytes += uint64(n)

		saddr := clientaddr.String()
		p.clientLock.Lock()
		conn, found := p.clients[saddr]

		// Create a new connection if not found
		if !found {
			conn = p.newConnection(p.raddr, clientaddr)
			if conn == nil {
				p.clientLock.Unlock()
				continue
			}
			p.clients[saddr] = conn
			log.Info("Created new connection for client %s\n", saddr)

			// Create a new connection to track the comms
			go p.waitForServer(conn)
		} else {
			log.Info("Found connection for client %s\n", saddr)
		}
		p.clientLock.Unlock()

		// Modify the incoming payload before sending to proxied server
		ctx := &muxy.Context{Bytes: b}
		for _, middleware := range p.middleware {
			log.Trace("UDP Proxy applying middleware %v", middleware)
			middleware.HandleEvent(muxy.EventPreDispatch, ctx)
			log.Trace("UDP Proxy overwriting bytes sent to target: %s", ctx.Bytes)
			b = ctx.Bytes
		}

		// Relay originating message, with any symptoms applied, to remote server
		_, err = conn.serverConn.Write(b)
		if handleConnectionErr(err) {
			log.Debug("Error relaying message to client %v", err)
			continue
		}
	}
}
