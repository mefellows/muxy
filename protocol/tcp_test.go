package protocol

import (
	"errors"
	"io"
	"log"
	"testing"
	"time"

	"fmt"

	"net"

	"github.com/mefellows/muxy/muxy"
	"github.com/mefellows/muxy/symptom"
)

func setupLocalTCP(port int) {
	go func() {
		l, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
			log.Fatal(err)
		}
		defer l.Close()

		for {
			// Wait for a connection.
			conn, err := l.Accept()
			if err != nil {
				log.Fatal(err)
			}
			// Handle the connection in a new goroutine.
			// The loop then returns to accepting, so that
			// multiple connections may be served concurrently.
			go func(c net.Conn) {
				// Echo all incoming data.
				io.Copy(c, c)
				// Shut down the connection.
				c.Close()
			}(conn)
		}
	}()
}

func TestTCPProxy_Setup(t *testing.T) {
	p := TCPProxy{}
	p.Setup([]muxy.Middleware{})
}

func TestTCPProxy_Teardown(t *testing.T) {
	p := TCPProxy{}
	p.Teardown()
}

func TestTCPProxy_Proxy(t *testing.T) {
	proxyPort := 7777
	proxyHost := "localhost"
	port := 7778
	setupLocalTCP(proxyPort)
	p := TCPProxy{
		Port:            port,
		Host:            "localhost",
		ProxyHost:       proxyHost,
		ProxyPort:       proxyPort,
		PacketSize:      64,
		NaglesAlgorithm: true,
	}
	// defer l.Close()

	// Wait for TCP server to be up
	waitForPort(proxyPort, t)

	// Start proxy and wait for it to be up
	go p.Proxy()
	waitForPort(port, t)

	// Send a message over TCP
	remoteAddr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%d", port))
	conn, _ := net.DialTCP("tcp", nil, remoteAddr)
	message := "some message"
	i, err := conn.Write([]byte(message))
	if err != nil {
		t.Fatal("Got error, want nil")
	}
	if i != len(message) {
		t.Fatal("Got", i, "want", len(message))
	}

	var b = make([]byte, 1024)
	i, err = conn.Read(b)
	if err != nil {
		t.Fatal("Got error, want nil")
	}
	if i != len(message) {
		t.Fatal("Got", i, "want", len(message))
	}

	conn.CloseRead()
}
func TestTCPProxy_ProxyWithMiddleware(t *testing.T) {
	proxyPort := 7779
	proxyHost := "localhost"
	setupLocalTCP(proxyPort)

	// Middleware
	var newRequestBody = "new request body"
	var newResponseBody = "new response body"
	tamperer := &symptom.TCPTampererSymptom{
		Request: symptom.TCPRequestConfig{
			Body:      newRequestBody,
			Randomize: false,
			Truncate:  false,
		},
		Response: symptom.TCPResponseConfig{
			Body:      newResponseBody,
			Randomize: false,
			Truncate:  false,
		},
	}
	tamperer.Setup()

	port := 7776
	p := TCPProxy{
		Port:            port,
		Host:            "localhost",
		ProxyHost:       proxyHost,
		ProxyPort:       proxyPort,
		PacketSize:      64,
		NaglesAlgorithm: true,
		middleware:      []muxy.Middleware{tamperer},
	}

	// Wait for TCP server to be up
	waitForPort(proxyPort, t)

	// Start proxy and wait for it to be up
	go p.Proxy()
	waitForPort(port, t)

	// Send a message over TCP
	remoteAddr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%d", port))
	conn, _ := net.DialTCP("tcp", nil, remoteAddr)
	message := "some message"

	i, err := conn.Write([]byte(message))
	if err != nil {
		t.Fatal("Got error, want nil")
	}
	if i != len(message) {
		t.Fatal("Got", i, "want", len(message))
	}

	var b = make([]byte, 1024)
	i, err = conn.Read(b)
	if err != nil {
		t.Fatal("Got error, want nil")
	}
	if i != len(newResponseBody) {
		t.Fatal("Got", i, "want", len(message))
	}

	conn.CloseRead()
}

func TestTCPProxy_ProxyFail(t *testing.T) {
	oldCheck := check
	doneChan := make(chan bool, 1)
	check = func(error) {
		doneChan <- true
	}
	defer func() {
		check = oldCheck
	}()

	proxyPort := 7775
	setupLocalTCP(proxyPort)
	p := TCPProxy{
		Port:            proxyPort,
		Host:            "localhost",
		ProxyHost:       "localhost",
		ProxyPort:       proxyPort,
		PacketSize:      64,
		NaglesAlgorithm: true,
	}
	// defer l.Close()

	// Start proxy and wait for it to be up
	go p.Proxy()

	for {
		select {
		case <-time.After(1 * time.Second):
			t.Fatal("Timeout waiting for Proxy to fail")
		case <-doneChan:
			return
		}
	}
}

func waitForPort(port int, t *testing.T) {
	timeout := time.After(1 * time.Second)
	for {
		select {
		case <-timeout:
			t.Fatalf("Expected server to start < 1s.")
		case <-time.After(50 * time.Millisecond):
			_, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
			if err == nil {
				return
			}
		}
	}
}

// GetFreePort Gets an available port by asking the kernal for a random port
// ready and available for use.
func GetFreePort() int {
	addr, err := net.ResolveTCPAddr("tcp", ":0")
	if err != nil {
		return 0
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func TestTCPProxy_pipe(t *testing.T) {
}

func TestTCPProxy_err(t *testing.T) {
	p := proxy{
		erred:  false,
		errsig: make(chan bool, 1),
	}
	p.err("problem", nil)
	p.err("real error", errors.New("some error"))

}
