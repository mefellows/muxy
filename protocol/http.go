package protocol

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"regexp"

	"github.com/mefellows/muxy/log"
	"github.com/mefellows/muxy/muxy"
	"github.com/mefellows/pkigo/pki"
	"github.com/mefellows/plugo/plugo"
)

// ProxyRequest contains details of the HTTP request to match
// to determine if the proxy rule should fire
type ProxyRequest struct {
	Method string
	// Headers map[string]string
	// Cookies []http.Cookie
	Path string
	Host string
}

// ProxyPass contains details of the HTTP request to
// send to the downstream proxy target
type ProxyPass struct {
	Method string
	// Headers map[string]string
	// Cookies []http.Cookie
	Path string
	Host string

	// Scheme is one of http or https
	Scheme string

	// Body    string // TODO: Could use this to return back a stub
	// TODO: Make this a templated body i.e. Accept ProxyRequest object and
	// Muxy Context or something so that the body can be smart/intelligent
	// Stub boolean // TODO: Use this to turn proxy rule into a stub instead.
}

// ProxyRule contains the rules for proxying a target HTTP system
type ProxyRule struct {
	Request ProxyRequest
	Pass    ProxyPass
}

// HTTPProxy implements the proxy interface for the HTTP protocol
type HTTPProxy struct {
	Port                int         `required:"true"`
	Host                string      `required:"true" default:"localhost"`
	Protocol            string      `default:"http" required:"true"`
	ProxyHost           string      `required:"true" mapstructure:"proxy_host"`
	ProxyPort           int         `required:"true" mapstructure:"proxy_port"`
	ProxyProtocol       string      `required:"true" default:"http" mapstructure:"proxy_protocol"`
	Insecure            bool        `required:"true" default:"false" mapstructure:"insecure"`
	ProxySslCertificate string      `required:"false" mapstructure:"proxy_ssl_cert"`
	ProxySslKey         string      `required:"false" mapstructure:"proxy_ssl_key"`
	ProxyClientSslCert  string      `required:"false" mapstructure:"proxy_client_ssl_cert"`
	ProxyClientSslKey   string      `required:"false" mapstructure:"proxy_client_ssl_key"`
	ProxyClientSslCa    string      `required:"false" mapstructure:"proxy_client_ssl_ca"`
	ProxyRules          []ProxyRule `required:"false" mapstructure:"proxy_rules"`
	middleware          []muxy.Middleware
}

func init() {
	plugo.PluginFactories.Register(func() (interface{}, error) {
		return &HTTPProxy{}, nil
	}, "http_proxy")
}

func (p *HTTPProxy) defaultProxyRule() ProxyRule {
	return ProxyRule{
		Request: ProxyRequest{
			Path:   "/",
			Host:   ".*",
			Method: ".*",
		},
		Pass: ProxyPass{
			Host: fmt.Sprintf("%s:%d", p.ProxyHost, p.ProxyPort),
		},
	}
}

// Setup sets up the middleware
func (p *HTTPProxy) Setup(middleware []muxy.Middleware) {
	p.middleware = middleware

	// Add default (catch all) proxy rule
	if len(p.ProxyRules) == 0 {
		p.ProxyRules = []ProxyRule{
			p.defaultProxyRule(),
		}
	} else {
		p.ProxyRules = append(p.ProxyRules, p.defaultProxyRule())
	}
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
		var proxy *ReverseProxy

		for _, rule := range p.ProxyRules {
			log.Trace("Matching request %v against ProxyRule %v", r, rule)

			if MatchRule(rule, *r) {
				log.Trace("Matched ProxyRule %v", rule)
				director := func(req *http.Request) {
					req = r

					p.ApplyProxyPassRule(rule, req)
				}

				proxy = &ReverseProxy{Director: director, Middleware: p.middleware}
				proxy.Transport = &http.Transport{
					Proxy:               http.ProxyFromEnvironment,
					TLSClientConfig:     config,
					TLSHandshakeTimeout: 10 * time.Second,
				}
				proxy.ServeHTTP(w, r)
				return
			}
		}

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

// ApplyProxyPassRule applies ProxyPass rules to the outbond request
// from the reverse proxy
func (p *HTTPProxy) ApplyProxyPassRule(rule ProxyRule, req *http.Request) {
	if rule.Pass.Method != "" {
		req.Method = rule.Pass.Method
	}

	if rule.Pass.Path != "" {
		req.URL.Path = fmt.Sprintf("%s%s", rule.Pass.Path, req.URL.Path)
	}

	if rule.Pass.Scheme != "" {
		req.URL.Scheme = rule.Pass.Scheme
	} else {
		req.URL.Scheme = p.ProxyProtocol
	}

	if rule.Pass.Host != "" {
		req.URL.Host = rule.Pass.Host
	} else {
		req.URL.Host = fmt.Sprintf("%s:%d", p.ProxyHost, p.ProxyPort)
	}
}

// MatchRule compares a ProxyRule to an http request to determine a match
func MatchRule(rule ProxyRule, req http.Request) bool {
	if rule.Request.Path != "" {
		log.Debug("ProxyRule matching path '%s' with '%s'", rule.Request.Path, req.URL.Path)
		if match, _ := regexp.MatchString(rule.Request.Path, req.URL.Path); !match {
			return false
		}
	}

	if rule.Request.Host != "" {
		log.Debug("ProxyRule matching host '%s' with '%s'", rule.Request.Host, req.URL.Host)
		if match, _ := regexp.MatchString(rule.Request.Host, req.URL.Host); !match {
			return false
		}
	}

	if rule.Request.Method != "" {
		log.Debug("ProxyRule matching method '%s' with '%s'", rule.Request.Method, req.Method)
		if match, _ := regexp.MatchString(rule.Request.Method, req.Method); !match {
			return false
		}
	}
	return true
}
