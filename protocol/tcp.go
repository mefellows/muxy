package protocol

import (
	"fmt"
	"io"
	"net"

	"github.com/mefellows/muxy/log"
	"github.com/mefellows/muxy/muxy"
	"github.com/mefellows/plugo/plugo"
)

// TCPProxy implements a TCP proxy
type TCPProxy struct {
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
		return &TCPProxy{}, nil
	}, "tcp_proxy")
}

var check = func(err error) {
	if err != nil {
		log.Fatalf("Error setting up TCP Proxy: %s", err.Error())
	}
}

// Setup the TCP proxy
func (p *TCPProxy) Setup(middleware []muxy.Middleware) {
	p.middleware = middleware
}

// Teardown the TCP proxy
func (p *TCPProxy) Teardown() {
}

// Proxy runs the TCP proxy
func (p *TCPProxy) Proxy() {
	log.Trace("Checking connection: %s:%d", p.Host, p.Port)
	laddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", p.Host, p.Port))
	check(err)
	raddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", p.ProxyHost, p.ProxyPort))
	check(err)
	listener, err := net.ListenTCP("tcp", laddr)
	check(err)

	for {
		log.Info("TCP Proxy proxy listening on %s", log.Colorize(log.BLUE, fmt.Sprintf("tcp://%s:%d", p.Host, p.Port)))
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Error("Failed to accept connection", err)
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
		log.Warn(p.prefix+s+"%v: ", err)
	}
	p.errsig <- true
	p.erred = true
}

func (p *proxy) start() {

	log.Trace("TCP Proxy Starting TCP Proxy")

	defer p.lconn.Close()

	// connect to remote
	log.Info("Connecting to %v", p.raddr)
	rconn, err := net.DialTCP("tcp", nil, p.raddr)
	if err != nil {
		p.err("TCP Proxy remote connection failed: %s", err)
		return
	}
	p.rconn = rconn
	defer p.rconn.Close()

	// nagles?
	if p.nagles {
		p.lconn.SetNoDelay(true)
		p.rconn.SetNoDelay(true)
	}

	// display both ends
	log.Info("TCP Proxy opened %s >>> %s", p.lconn.RemoteAddr().String(), p.rconn.RemoteAddr().String())

	// bidirectional copy
	go p.pipe(p.lconn, p.rconn)
	go p.pipe(p.rconn, p.lconn)

	//wait for close...
	<-p.errsig

	log.Info("TCP Proxy closed (%d bytes sent, %d bytes received)", p.sentBytes, p.receivedBytes)
}

func (p *proxy) pipe(src io.Reader, dst io.Writer) {
	// Direction of traffic
	islocal := src == p.lconn

	buff := make([]byte, p.packetsize)
	done := false
	for !done {
		n, readErr := src.Read(buff)
		if readErr != nil || n == 0 {
			if !islocal {
				p.err("TCP Proxy read failed: ", readErr)
			}
			done = true
		}

		b := buff[:n]

		ctx := &muxy.Context{Bytes: b}
		for _, middleware := range p.middleware {
			log.Trace("TCP Proxy applying middleware %v", middleware)
			if islocal {
				middleware.HandleEvent(muxy.EventPreDispatch, ctx)
				log.Trace("TCP Proxy overwriting bytes sent to target: %s", ctx.Bytes)
			} else {
				middleware.HandleEvent(muxy.EventPostDispatch, ctx)
				log.Trace("TCP Proxy overwriting bytes sent back to originating client: %s", ctx.Bytes)
			}
			b = ctx.Bytes
		}

		n, err := dst.Write(b)
		if err != nil {
			log.Error("TCP Proxy write failed: %s", err.Error())
			p.err("TCP Proxy write failed '%s'\n", err)

			return
		}
		if islocal {
			p.sentBytes += uint64(n)
		} else {
			p.receivedBytes += uint64(n)
		}
	}
}
