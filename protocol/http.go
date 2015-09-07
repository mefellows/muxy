package protocol

import (
	"fmt"
	"github.com/mefellows/muxy/log"
	"github.com/mefellows/muxy/muxy"
	"github.com/mefellows/plugo/plugo"
	"net/http"
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
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		director := func(req *http.Request) {
			req = r
			req.URL.Scheme = p.ProxyProtocol
			req.URL.Host = fmt.Sprintf("%s:%d", p.ProxyHost, p.ProxyPort)
		}
		proxy := &ReverseProxy{Director: director, Middleware: p.middleware}
		proxy.ServeHTTP(w, r)
	})
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", p.Host, p.Port), mux)
	if err != nil {
		log.Info("ListenAndServe error: ", err.Error())
	}
}
