package symptom

import (
	"github.com/mefellows/muxy/muxy"
	"github.com/tylertreat/comcast/throttler"
	"log"
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
}

func init() {
	muxy.PluginFactories.Register(func() (interface{}, error) {
		return &ShittyNetworkSymptom{}, nil
	}, "network_shape")

}

func (s *ShittyNetworkSymptom) Setup() {
	log.Printf("Setting up ShittyNetworkSymptom: Enabling firewall")

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
	throttler.Run(&s.config)
}

func (m ShittyNetworkSymptom) HandleEvent(e muxy.ProxyEvent, ctx *muxy.Context) {
	switch e {
	case muxy.EVENT_PRE_DISPATCH:
		m.Muck(ctx)
	}
}

func (s *ShittyNetworkSymptom) Muck(ctx *muxy.Context) {
	log.Printf("Mucking...")
}

func (s *ShittyNetworkSymptom) Teardown() {
	log.Printf("Tearing down ShittyNetworkSymptom")
	s.config.Stop = true
	throttler.Run(&s.config)
}
