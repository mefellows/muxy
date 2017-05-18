package run

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"testing"
	"time"

	_ "github.com/mefellows/muxy/middleware"
	_ "github.com/mefellows/muxy/protocol"
	_ "github.com/mefellows/muxy/symptom"
)

var proxiedServerBody = "proxied server up!"

func TestTCPProxy_New(t *testing.T) {
	config := &Config{}

	m := New(config)

	if m == nil {
		t.Fatalf("must be a Muxy")
	}
}

func TestTCPProxy_NewWithDefaultConfig(t *testing.T) {
	m := NewWithDefaultConfig()

	if m == nil {
		t.Fatalf("must be a Muxy")
	}
}

func TestTCPProxy_Run(t *testing.T) {
	m := NewWithDefaultConfig()
	m.config.ConfigFile = "muxy_test.yml"

	proxiedPort := 8889
	port := 8888
	runTestServer(proxiedPort)
	go m.Run()
	waitForPort(port, t)
	waitForPort(proxiedPort, t)

	// Do some tests
	res, err := http.Get(fmt.Sprintf("http://localhost:%d", port))
	if err != nil {
		log.Fatal(err)
	}

	// Check body
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	if string(body) != proxiedServerBody {
		t.Fatal("Want", proxiedServerBody, "got", string(body))
	}
}

func TestTCPProxy_LoadPlugins(t *testing.T) {

}

func runTestServer(port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(proxiedServerBody))
	})

	go http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
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
