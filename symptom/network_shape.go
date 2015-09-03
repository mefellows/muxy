package symptom

import (
	"github.com/tylertreat/comcast/throttler"
	"log"
)

// Shape bandwidth to mobile, slower speeds

//
type ShittyNetworkSymptom struct {
	config throttler.Config
}

func (s *ShittyNetworkSymptom) Setup() {
	log.Printf("Setting up ShittyNetworkSymptom: Enabling firewall")
	//// TODO: Add support for other options like packet reordering, duplication, etc.
	//var (
	//    device      = flag.String("device", "", "Interface (device) to use (defaults to eth0 where applicable)")
	//    mode        = flag.String("mode", throttler.Start, "Start or stop packet controls")
	//    latency     = flag.Int("latency", -1, "Latency to add in ms")
	//    targetbw    = flag.Int("target-bw", -1, "Target bandwidth limit in kbit/s (slow-lane)")
	//    defaultbw   = flag.Int("default-bw", -1, "Default bandwidth limit in kbit/s (fast-lane)")
	//    packetLoss  = flag.String("packet-loss", "0", "Packet loss percentage (eg: 0.1%%)")
	//    targetaddr  = flag.String("target-addr", "", "Target addresses, (eg: 10.0.0.1 or 10.0.0.0/24 or 10.0.0.1,192.168.0.0/24 or 2001:db8:a::123)")
	//    targetport  = flag.String("target-port", "", "Target port(s) (eg: 80 or 1:65535 or 22,80,443,1000:1010)")
	//    targetproto = flag.String("target-proto", "tcp,udp,icmp", "Target protocol TCP/UDP (eg: tcp or tcp,udp or icmp)")
	//    dryrun      = flag.Bool("dry-run", false, "Specifies whether or not to actually commit the rule changes")
	//    //icmptype    = flag.String("icmp-type", "", "icmp message type (eg: reply or reply,request)") //TODO: Maybe later :3
	//)
	//flag.Parse()
	//
	//targetIPv4, targetIPv6 := parseAddrs(*targetaddr)

	s.config = throttler.Config{
		Device:           "eth0",
		Latency:          250,
		TargetBandwidth:  750,
		DefaultBandwidth: 740,
		PacketLoss:       1.5,
		TargetIps:        []string{},
		TargetIps6:       []string{},
		TargetPorts:      []string{"8282"},
		TargetProtos:     []string{"tcp"},
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
