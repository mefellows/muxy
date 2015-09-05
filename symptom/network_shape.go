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
	TargetFoo        []int    `mapstructure:"target_protos" required:"true" default:"10,20,30,40"`
}

func init() {
	muxy.SymptomFactories.Register(func() (muxy.Symptom, error) {
		return &ShittyNetworkSymptom{}, nil
	}, "network_shape")

}

func (s *ShittyNetworkSymptom) Setup() {
	log.Printf("Setting up ShittyNetworkSymptom: Enabling firewall")
	log.Printf("What's my array value? %v\n", s.TargetFoo)

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

func (s *ShittyNetworkSymptom) Muck() {
	log.Printf("Mucking...")
}

func (s *ShittyNetworkSymptom) Teardown() {
	log.Printf("Tearing down ShittyNetworkSymptom")
	s.config.Stop = true
	throttler.Run(&s.config)
}
