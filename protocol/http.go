package protocol

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/mefellows/muxy/log"
	"github.com/mefellows/muxy/muxy"
	"github.com/mefellows/pkigo/pki"
	"github.com/mefellows/plugo/plugo"
)

// HTTPProxy implements the proxy interface for the HTTP protocol
type HTTPProxy struct {
	Port                int    `required:"true"`
	Host                string `required:"true" default:"localhost"`
	Protocol            string `default:"http" required:"true"`
	ProxyHost           string `required:"true" mapstructure:"proxy_host"`
	ProxyPort           int    `required:"true" mapstructure:"proxy_port"`
	ProxyProtocol       string `required:"true" default:"http" mapstructure:"proxy_protocol"`
	Insecure            bool   `required:"true" default:"false" mapstructure:"insecure"`
	ProxySslCertificate string `required:"false" mapstructure:"proxy_ssl_cert"`
	ProxySslKey         string `required:"false" mapstructure:"proxy_ssl_key"`
	ProxyClientSslCert  string `required:"false" mapstructure:"proxy_client_ssl_cert"`
	ProxyClientSslKey   string `required:"false" mapstructure:"proxy_client_ssl_key"`
	ProxyClientSslCa    string `required:"false" mapstructure:"proxy_client_ssl_ca"`
	middleware          []muxy.Middleware
}

func init() {
	plugo.PluginFactories.Register(func() (interface{}, error) {
		return &HTTPProxy{}, nil
	}, "http_proxy")
}

// Setup sets up the middleware
func (p *HTTPProxy) Setup(middleware []muxy.Middleware) {
	p.middleware = middleware
}

// Teardown shuts down the middleware
func (p *HTTPProxy) Teardown() {
}

// Proxy performs the proxy event
func (p *HTTPProxy) Proxy() {
	log.Info("HTTP proxy listening on %s", log.Colorize(log.BLUE, fmt.Sprintf("%s://%s:%d", p.Protocol, p.Host, p.Port)))

	// TODO: Suggest we only try to configure SSL infra iff:
	//        a) https is requested; and
	//        b) custom certs/keys not provided
	pkiMgr, err := pki.New()
	checkHTTPServerError(err)

	config, err := pkiMgr.GetClientTLSConfig()
	checkHTTPServerError(err)

	// Override SSL / TLS settings
	config.InsecureSkipVerify = p.Insecure

	if p.ProxySslCertificate == "" {
		p.ProxySslCertificate = pkiMgr.Config.ServerCertPath
	}

	if p.ProxySslKey == "" {
		p.ProxySslKey = pkiMgr.Config.ServerKeyPath
	}

	// MASSL (client certiicate) setup
	if p.ProxyClientSslCert != "" && p.ProxyClientSslKey != "" && p.ProxyClientSslCa != "" {
		// Load client cert
		cert, err := tls.LoadX509KeyPair(p.ProxyClientSslCert, p.ProxyClientSslKey)
		if err != nil {
			log.Fatal(err)
		}

		// Load CA cert
		caCert, err := ioutil.ReadFile(p.ProxyClientSslCa)
		if err != nil {
			log.Fatal(err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		// Setup HTTPS client
		config.Certificates = []tls.Certificate{cert}
		config.RootCAs = caCertPool
		config.BuildNameToCertificate()
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		director := func(req *http.Request) {
			req = r
			req.URL.Scheme = p.ProxyProtocol
			req.URL.Host = fmt.Sprintf("%s:%d", p.ProxyHost, p.ProxyPort)
		}

		proxy := &ReverseProxy{Director: director, Middleware: p.middleware}
		proxy.Transport = &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			TLSClientConfig:     config,
			TLSHandshakeTimeout: 10 * time.Second,
		}

		proxy.ServeHTTP(w, r)
	})

	if p.Protocol == "https" {
		checkHTTPServerError(err)
		checkHTTPServerError(http.ListenAndServeTLS(fmt.Sprintf("%s:%d", p.Host, p.Port), p.ProxySslCertificate, p.ProxySslKey, mux))
	} else {
		checkHTTPServerError(http.ListenAndServe(fmt.Sprintf("%s:%d", p.Host, p.Port), mux))
	}
}

func checkHTTPServerError(err error) {
	if err != nil {
		log.Error("ListenAndServe error: ", err.Error())
	}
}
