package protocol

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mefellows/muxy/log"
	"github.com/mefellows/muxy/muxy"
	"github.com/mefellows/pkigo/pki"
	"github.com/mefellows/plugo/plugo"
)

// HTTPProxy implements the proxy interface for the HTTP protocol
// nolint
type HTTPProxy struct {
	Port          int    `required:"true"`
	Host          string `required:"true" default:"localhost"`
	Protocol      string `default:"http" required:"true"`
	ProxyHost     string `required:"true" mapstructure:"proxy_host"`
	ProxyPort     int    `required:"true" mapstructure:"proxy_port"`
	ProxyProtocol string `required:"true" default:"http" mapstructure:"proxy_protocol"`
	Insecure      bool   `required:"true" default:"false" mapstructure:"insecure"`
	middleware    []muxy.Middleware
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
	pkiMgr, err := pki.New()
	checkHTTPServerError(err)
	config, err := pkiMgr.GetClientTLSConfig()
	checkHTTPServerError(err)
	config.InsecureSkipVerify = p.Insecure

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
		checkHTTPServerError(http.ListenAndServeTLS(fmt.Sprintf("%s:%d", p.Host, p.Port), pkiMgr.Config.ClientCertPath, pkiMgr.Config.ClientKeyPath, mux))
	} else {
		checkHTTPServerError(http.ListenAndServe(fmt.Sprintf("%s:%d", p.Host, p.Port), mux))
	}
}

func checkHTTPServerError(err error) {
	if err != nil {
		log.Error("ListenAndServe error: ", err.Error())
	}
}
