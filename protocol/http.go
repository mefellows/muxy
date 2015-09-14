package protocol

import (
	"fmt"
	"github.com/mefellows/muxy/log"
	"github.com/mefellows/muxy/muxy"
	"github.com/mefellows/pkigo/pki"
	"github.com/mefellows/plugo/plugo"
	"net/http"
	"time"
)

type HttpProxy struct {
	Port          int    `required:"true"`
	Host          string `required:"true" default:"localhost"`
	Protocol      string `default:"http" required:"true"`
	ProxyHost     string `required:"true" mapstructure:"proxy_host"`
	ProxyPort     int    `required:"true" mapstructure:"proxy_port"`
	ProxyProtocol string `required:"true" default:"http" mapstructure:"proxy_protocol"`
	middleware    []muxy.Middleware
}

func init() {
	plugo.PluginFactories.Register(func() (interface{}, error) {
		return &HttpProxy{}, nil
	}, "http_proxy")
}

func (p *HttpProxy) Setup(middleware []muxy.Middleware) {
	p.middleware = middleware
}

func (p *HttpProxy) Teardown() {
}

func (p *HttpProxy) Proxy() {
	log.Info("HTTP proxy listening on %s", log.Colorize(log.BLUE, fmt.Sprintf("%s://%s:%d", p.Protocol, p.Host, p.Port)))
	pkiMgr, err := pki.New()
	checkHttpServerError(err)
	config, err := pkiMgr.GetClientTLSConfig()
	checkHttpServerError(err)
	config.InsecureSkipVerify = false

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
		checkHttpServerError(err)
		checkHttpServerError(http.ListenAndServeTLS(fmt.Sprintf("%s:%d", p.Host, p.Port), pkiMgr.Config.ClientCertPath, pkiMgr.Config.ClientKeyPath, mux))
	} else {
		checkHttpServerError(http.ListenAndServe(fmt.Sprintf("%s:%d", p.Host, p.Port), mux))
	}
}

func checkHttpServerError(err error) {
	if err != nil {
		log.Error("ListenAndServe error: ", err.Error())
	}
}
