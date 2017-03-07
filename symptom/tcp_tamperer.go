package symptom

import (
	"math/rand"
	"time"

	"github.com/mefellows/muxy/log"
	"github.com/mefellows/muxy/muxy"
	"github.com/mefellows/plugo/plugo"
)

// TCPRequestConfig contains details of the HTTP request to tamper with prior to
// sending on to the target system
type TCPRequestConfig struct {
	// Body fixes the request message
	Body string

	// Randomize request message
	Randomize bool

	// Truncate the request message. Removes trailing char
	Truncate bool
}

// TCPResponseConfig contains details of the TCP response to tamper with prior to
// sending on to the initiating system
type TCPResponseConfig struct {
	// Body fixes the response message
	Body string

	// Randomize response message
	Randomize bool

	// Truncate the response message. Removes trailing char
	Truncate bool
}

// TCPTampererSymptom is a plugin to mess with request/responses between
// a consumer and provider system
type TCPTampererSymptom struct {
	Request  TCPRequestConfig
	Response TCPResponseConfig
}

func init() {
	plugo.PluginFactories.Register(func() (interface{}, error) {
		return &TCPTampererSymptom{}, nil
	}, "tcp_tamperer")
}

// Setup sets up the plugin
func (m TCPTampererSymptom) Setup() {
	log.Debug("TCP Tamperer Setup()")
}

// Teardown shuts down the plugin
func (m TCPTampererSymptom) Teardown() {
	log.Debug("TCP Tamperer Teardown()")
}

// HandleEvent is a hook to allow the plugin to intervene with a request/response
// event
func (m *TCPTampererSymptom) HandleEvent(e muxy.ProxyEvent, ctx *muxy.Context) {
	log.Trace("TCP Tamperer - Handle Event")
	switch e {
	case muxy.EventPreDispatch:
		log.Debug("TCP Tamperer - Handle Event - pre dispatch")
		m.MuckRequest(ctx)
	case muxy.EventPostDispatch:
		log.Debug("TCP Tamperer - Handle Event - post dispatch")
		m.MuckResponse(ctx)
	}
}

// MuckRequest adds chaos to the request
func (m *TCPTampererSymptom) MuckRequest(ctx *muxy.Context) {
	if m.Request.Body != "" {
		log.Debug("TCP Tamperer - replacing body [%s] with [%s]", ctx.Bytes, m.Request.Body)
		ctx.Bytes = []byte(m.Request.Body)
	}
	if m.Request.Randomize {
		random := randStringBytesMaskImprSrc(len(ctx.Bytes))
		log.Debug("TCP Tamperer - randomizing body [%s] with [%s]", ctx.Bytes, random)
		ctx.Bytes = random
	}
	if m.Request.Truncate {
		if len(ctx.Bytes) >= 2 {
			log.Debug("TCP Tamperer - randomizing body [%s] with [%s]", ctx.Bytes, ctx.Bytes[:len(ctx.Bytes)-2])
			ctx.Bytes = ctx.Bytes[:len(ctx.Bytes)-2]
		}
	}
}

// MuckResponse adds chaos to the response
func (m *TCPTampererSymptom) MuckResponse(ctx *muxy.Context) {
	if m.Response.Body != "" {
		log.Debug("TCP Tamperer - replacing body [%s] with [%s]", ctx.Bytes, m.Request.Body)
		ctx.Bytes = []byte(m.Response.Body)
	}
	if m.Response.Randomize {
		random := randStringBytesMaskImprSrc(len(ctx.Bytes))
		log.Debug("TCP Tamperer - randomizing body [%s] with [%s]", ctx.Bytes, random)
		ctx.Bytes = random
	}
	if m.Response.Truncate {
		// TODO: why 2 and 3?
		if len(ctx.Bytes) >= 3 {
			log.Debug("TCP Tamperer - randomizing body [%s] with [%s]", ctx.Bytes, ctx.Bytes[:len(ctx.Bytes)-3])
			ctx.Bytes = ctx.Bytes[:len(ctx.Bytes)-3]
		}
	}
}

// Randomly mess with bytes in an array
var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// Generate string of a given size
// Courtesy of http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func randStringBytesMaskImprSrc(n int) []byte {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return b
}
