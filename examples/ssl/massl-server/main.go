package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func helloServer(w http.ResponseWriter, req *http.Request) {
	log.Println("MASSL Server - /hello called")
	io.WriteString(w, "hello, world!\n")
}
func fileNotFoundServer(w http.ResponseWriter, req *http.Request) {
	log.Println("404: ", req.URL.Path, "not found")
	io.WriteString(w, fmt.Sprint("404", req.URL.Path, " not found\n"))
}

func main() {
	http.HandleFunc("/hello", helloServer)
	http.HandleFunc("/", fileNotFoundServer)

	caCert, err := ioutil.ReadFile("ca.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		ClientCAs: caCertPool,
		// NoClientCert
		// RequestClientCert
		// RequireAnyClientCert
		// VerifyClientCertIfGiven
		// RequireAndVerifyClientCert
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()

	server := &http.Server{
		Addr:      ":8080",
		TLSConfig: tlsConfig,
	}

	log.Println("MASSL Server Listening on port 8080")
	log.Println("")
	log.Println("curl --cacert ca.pem -E ./client.p12:password -v https://localhost:8080/hello")
	server.ListenAndServeTLS("server-cert.pem", "server-key.pem") //private cert
}
