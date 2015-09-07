package symptom

import (
	"bytes"
	"github.com/mefellows/muxy/log"
	"github.com/mefellows/muxy/muxy"
	"github.com/mefellows/plugo/plugo"
	"github.com/tylertreat/comcast/throttler"
	"io"
	l "log"
	"os"
)

// Shape bandwidth to mobile, slower speeds
type ShittyNetworkSymptom struct {
	config           throttler.Config
	Device           string
	Latency          int
	TargetBandwidth  int      `mapstructure:"target_bw"`
	DefaultBandwidth int      `mapstructure:"default_bw"`
	PacketLoss       float64  `mapstructure:"packet_loss"`
	TargetIps        []string `mapstructure:"target_ips"`
	TargetIps6       []string `mapstructure:"target_ips6"`
	TargetPorts      []string `mapstructure:"target_ports"`
	TargetProtos     []string `mapstructure:"target_protos" required:"true" default:"tcp,icmp"`
	out              io.Writer
	err              io.Writer
}

func init() {
	plugo.PluginFactories.Register(func() (interface{}, error) {
		return &ShittyNetworkSymptom{}, nil
	}, "network_shape")

}

func (s *ShittyNetworkSymptom) Setup() {
	log.Debug("Setting up ShittyNetworkSymptom: Enabling firewall")

	s.config = throttler.Config{
		Device:           s.Device,
		Latency:          s.Latency,
		TargetBandwidth:  s.TargetBandwidth,
		DefaultBandwidth: s.DefaultBandwidth,
		PacketLoss:       s.PacketLoss,
		TargetIps:        s.TargetIps,
		TargetIps6:       s.TargetIps6,
		TargetPorts:      s.TargetPorts,
		TargetProtos:     s.TargetProtos,
		DryRun:           false,
	}

	supressOutput(func() {
		throttler.Run(&s.config)
	})

}

func (m ShittyNetworkSymptom) HandleEvent(e muxy.ProxyEvent, ctx *muxy.Context) {
	switch e {
	case muxy.EVENT_PRE_DISPATCH:
		m.Muck(ctx)
	}
}

func (s *ShittyNetworkSymptom) Muck(ctx *muxy.Context) {
	log.Debug("ShittyNetworkSymptom Mucking...")
}

func (s *ShittyNetworkSymptom) Teardown() {
	log.Debug("Tearing down ShittyNetworkSymptom")
	s.config.Stop = true
	supressOutput(func() {
		throttler.Run(&s.config)
	})
}

// Supress output of function to keep logs clean
func supressOutput(f func()) {
	old := os.Stdout
	oldErr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w
	l.SetOutput(w)

	f()

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old
	os.Stderr = oldErr
	l.SetOutput(old)
}
