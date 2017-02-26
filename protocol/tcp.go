package protocol

import (
	"fmt"
	"io"
	"net"

	"github.com/mefellows/muxy/log"
	"github.com/mefellows/muxy/muxy"
	"github.com/mefellows/plugo/plugo"
)

// TcpProxy implements a TCP proxy
// nolint
type TcpProxy struct {
	Port            int    `required:"true"`
	Host            string `required:"true" default:"localhost"`
	ProxyHost       string `required:"true" mapstructure:"proxy_host"`
	ProxyPort       int    `required:"true" mapstructure:"proxy_port"`
	NaglesAlgorithm bool   `mapstructure:"nagles_algorithm"`
	HexOutput       bool   `mapstructure:"hex_output"`
	PacketSize      int    `mapstructure:"packet_size" default:"64" required:"true"`
	connID          uint64
	middleware      []muxy.Middleware
}

func init() {
	plugo.PluginFactories.Register(func() (interface{}, error) {
		return &TcpProxy{}, nil
	}, "tcp_proxy")
}

func check(err error) {
	if err != nil {
		log.Fatalf("Error setting up TCP Proxy: %s", err.Error())
	}
}

// Setup sets up the TCP proxy
func (p *TcpProxy) Setup(middleware []muxy.Middleware) {
	p.middleware = middleware
}

// Teardown shuts down the TCP proxy
func (p *TcpProxy) Teardown() {
}

// Proxy runs the TCP proxy
func (p *TcpProxy) Proxy() {
	log.Trace("Checking connection: %s:%d", p.Host, p.Port)
	laddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", p.Host, p.Port))
	check(err)
	log.Trace("Checking connection: %s:%d", p.ProxyHost, p.ProxyPort)
	raddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", p.ProxyHost, p.ProxyPort))
	check(err)
	listener, err := net.ListenTCP("tcp", laddr)
	check(err)

	for {
		log.Info("TCP proxy listening on %s", log.Colorize(log.BLUE, fmt.Sprintf("tcp://%s:%d", p.Host, p.Port)))
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Printf("Failed to accept connection '%s'\n", err)
			continue
		}
		p.connID++

		p := &proxy{
			lconn:      conn,
			laddr:      laddr,
			raddr:      raddr,
			packetsize: p.PacketSize,
			erred:      false,
			errsig:     make(chan bool),
			prefix:     fmt.Sprintf("Connection #%03d ", p.connID),
			hex:        p.HexOutput,
			nagles:     p.NaglesAlgorithm,
			middleware: p.middleware,
		}
		go p.start()
	}
}

//A proxy represents a pair of connections and their state
type proxy struct {
	middleware    []muxy.Middleware
	sentBytes     uint64
	receivedBytes uint64
	laddr, raddr  *net.TCPAddr
	lconn, rconn  *net.TCPConn
	protocol      string
	erred         bool
	errsig        chan bool
	prefix        string
	matcher       func([]byte)
	replacer      func([]byte) []byte
	nagles        bool
	hex           bool
	packetsize    int
}

func (p *proxy) err(s string, err error) {
	if p.erred {
		return
	}
	if err != io.EOF {
		log.Warn(p.prefix+s, err)
	}
	p.errsig <- true
	p.erred = true
}

func (p *proxy) start() {

	log.Trace("Starting TCP Proxy")

	defer p.lconn.Close()
	//connect to remote
	rconn, err := net.DialTCP("tcp", nil, p.raddr)
	if err != nil {
		p.err("Remote connection failed: %s", err)
		return
	}
	p.rconn = rconn
	defer p.rconn.Close()
	//nagles?
	if p.nagles {
		p.lconn.SetNoDelay(true)
		p.rconn.SetNoDelay(true)
	}
	//display both ends
	log.Info("Opened %s >>> %s", p.lconn.RemoteAddr().String(), p.rconn.RemoteAddr().String())

	//bidirectional copy
	go p.pipe(p.lconn, p.rconn)
	go p.pipe(p.rconn, p.lconn)
	//wait for close...
	<-p.errsig
	log.Info("Closed (%d bytes sent, %d bytes received)", p.sentBytes, p.receivedBytes)
}

func (p *proxy) pipe(src io.Reader, dst io.Writer) {
	// Direction
	islocal := src == p.lconn

	buff := make([]byte, p.packetsize)
	done := false
	for !done {
		n, readErr := src.Read(buff)
		if readErr != nil || n == 0 {
			if !islocal {
				p.err("Read failed '%s'\n", readErr)
			}
			done = true
		}

		b := buff[:n]

		ctx := &muxy.Context{Bytes: b}
		for _, middleware := range p.middleware {
			if islocal {
				middleware.HandleEvent(muxy.EventPreDispatch, ctx)
			} else {
				middleware.HandleEvent(muxy.EventPostDispatch, ctx)
			}
		}

		n, err := dst.Write(b)
		if err != nil {
			log.Error("Write failed: %s", err.Error())
			p.err("Write failed '%s'\n", err)

			return
		}
		if islocal {
			p.sentBytes += uint64(n)
		} else {
			p.receivedBytes += uint64(n)
		}
	}
}
