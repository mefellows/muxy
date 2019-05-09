package symptom

import (
	"log"
	"testing"

	"strings"

	"github.com/mefellows/muxy/muxy"
	"github.com/mefellows/muxy/symptom/throttler"
)

func TestNetworkShapeSymptom_Setup(t *testing.T) {
	s := NetworkShaperSymptom{}
	oldThrottler := executeThrottler
	executeThrottler = func(c *throttler.Config) {
		t.Log("Calling fake throttler")
	}
	defer func() {
		executeThrottler = oldThrottler
	}()

	s.Setup()
}

func TestNetworkShapeSymptom_Teardown(t *testing.T) {
	s := NetworkShaperSymptom{}
	oldThrottler := executeThrottler
	executeThrottler = func(c *throttler.Config) {
		t.Log("Calling fake throttler")
	}
	defer func() {
		executeThrottler = oldThrottler
	}()

	s.Teardown()
}

func TestNetworkShapeSymptom_Muck(t *testing.T) {
	s := NetworkShaperSymptom{}
	s.Muck(nil)
}

func TestNetworkShapeSymptom_HandleEvent(t *testing.T) {
	s := NetworkShaperSymptom{}
	s.HandleEvent(muxy.EventPreDispatch, nil)
}

func TestNetworkShapeSymptom_parseProtos(t *testing.T) {
	var called bool
	oldFail := fail
	fail = func(string, ...interface{}) {
		called = true
	}
	defer func() {
		fail = oldFail
	}()

	// Valid
	valid := "udp,tcp,icmp"
	parsed := parseProtos(valid)
	if len(parsed) != 3 {
		t.Fatal("Want 3, got", len(parsed))
	}

	// Invalid
	invalid := "udp,tcp,foo,icmp"
	parseProtos(invalid)
	if !called {
		t.Fatal("Expected parseProtos to fail but did not")
	}

	// Invalid - bad formatting
	called = false
	invalid = "tcp udp icmp"
	parseProtos(invalid)
	if !called {
		t.Fatal("Expected parseProtos to fail but did not")
	}
}

func TestNetworkShapeSymptom_suppressOutput(t *testing.T) {
	supressOutput(func() {
		log.Println("something")
	})
}

func TestNetworkShapeSymptom_parseLoss(t *testing.T) {
	// For the failure cases
	var called bool
	oldFail := fail
	fail = func(string, ...interface{}) {
		called = true
	}
	defer func() {
		fail = oldFail
	}()

	// Valid
	valid := map[string]float64{
		"100%":   100,
		"100":    100,
		"0":      0,
		"0%":     0,
		" 100% ": 100,
	}
	invalid := []string{
		"",
		"%",
		"abc",
	}

	for test, want := range valid {
		answer := parseLoss(test)
		if answer != want {
			t.Fatal("Want", want, "got", answer)
		}
	}

	// Failures
	for _, test := range invalid {
		parseLoss(test)
		if !called {
			t.Fatal("Expected parseProtos to fail but did not")
		}
		called = false
	}
}

func TestNetworkShapeSymptom_parseAddresses(t *testing.T) {
	// For the failure cases
	var called bool
	oldFail := fail
	fail = func(string, ...interface{}) {
		called = true
	}
	defer func() {
		fail = oldFail
	}()

	// Valid
	valid := []string{
		// IPv4s
		"127.0.0.1",
		"0.0.0.0",
		"255.255.255.255",

		// IPv6s
		"::1",
		"fe80::82e6:50ff:fe28:6c",
		"fe80::1",

		// CIDRs
		"192.168.100.14/24",
		"::1/24",
	}
	invalid := []string{
		"127",
		"%",
		"abc",
	}

	validTest := strings.Join(valid, ",")

	ipv4s, ipv6s := parseAddrs(validTest)

	if len(ipv4s) != 4 {
		t.Fatal("Want 4 got", len(ipv4s))
	}
	if len(ipv6s) != 4 {
		t.Fatal("Want 4 got", len(ipv6s))
	}
	if called {
		t.Fatal("Expected parseAddrs to pass but failed!")
	}

	// Failures: individual
	for _, test := range invalid {
		parseAddrs(test)
		if !called {
			t.Fatal("Expected parseProtos to fail but did not")
		}
		called = false
	}

	// Failures: joined
	invalidTest := strings.Join(invalid, ",")
	parseAddrs(invalidTest)
	if !called {
		t.Fatal("Expected parseProtos to fail but did not")
	}
}

func TestNetworkShapeSymptom_parsePort(t *testing.T) {
	validPorts := map[string]int{
		"80":     80,
		"  80  ": 80,
	}
	for test, want := range validPorts {
		got := parsePort(test)
		if want != got {
			log.Fatal("Want", want, "got", got)
		}
	}
	invalidPorts := []string{
		"",
		"a",
		"eighty",
	}
	for _, test := range invalidPorts {
		got := parsePort(test)
		if got != 0 {
			log.Fatal("Wan 0, got", got)
		}
	}
}

func TestNetworkShapeSymptom_validPort(t *testing.T) {
	validPorts := []string{
		"1",
		"65535",
		"80",
		" 80 ",
		" 80",
		"80 ",
	}
	for _, test := range validPorts {
		if !validPort(test) {
			t.Fatal("Expected", test, "to be valid port but was not")
		}
	}
	invalidPorts := []string{
		" ",
		"-1",
		"0",
		"65536",
		"100000000000000000",
	}
	for _, test := range invalidPorts {
		if validPort(test) {
			t.Fatal("Expected", test, "to be an invalid port but was")
		}
	}
}

func TestNetworkShapeSymptom_validRange(t *testing.T) {
	validPortRange := []string{
		"1:65535",
		"1 : 65535",
		" 1:65535 ",
		" 1 : 65535 ",
	}
	for _, test := range validPortRange {
		if !validRange(test) {
			t.Fatal("Expected", test, "to be valid port but was not")
		}
	}
	invalidPortRange := []string{
		" ",
		"-1",
		"0",
		"65536",
		"100000000000000000",
		"0:65536",
		"0:0",
		"-1:0",
		"100:1",
		"-1:1",
		"-1:65535",
	}
	for _, test := range invalidPortRange {
		if validRange(test) {
			t.Fatal("Expected", test, "to be an invalid port but was")
		}
	}
}

func TestNetworkShapeSymptom_parsePorts(t *testing.T) {
	// For the failure cases
	var called bool
	oldFail := fail
	fail = func(string, ...interface{}) {
		called = true
	}
	defer func() {
		fail = oldFail
	}()

	// Valid
	validPortRange := []string{
		"1:65535",
		"1 : 65535",
		" 1:65535 ",
		" 1 : 65535 ",
		"80",
	}
	for _, test := range validPortRange {
		res := parsePorts(test)
		if len(res) != 1 {
			t.Fatal("Want 1, got", len(res))
		}
		if called {
			t.Fatal("Expected parsePorts to pass but did not")
		}
	}
	testValidString := strings.Join(validPortRange, ",")

	res := parsePorts(testValidString)
	if len(res) != 5 {
		t.Fatal("Want 5, got", len(res))
	}
	if called {
		t.Fatal("Expected parsePorts to pass but did not")
	}

	// Invalid
	invalidPortRange := []string{
		" ",
		"-1",
		"0",
		"65536",
		"100000000000000000",
		"0:65536",
		"0:0",
		"-1:0",
		"-1:1",
		"-1:65535",
	}
	for _, test := range invalidPortRange {
		parsePorts(test)

		if !called {
			t.Fatal("Expected parsePorts to fail but did not")
		}
	}
}

func TestNetworkShapeSymptom_portHigher(t *testing.T) {
	validPorts := map[string]string{
		"0": "-1",
		"1": "0",
	}
	for p1, p2 := range validPorts {
		if !portHigher(p1, p2) {
			t.Fatal("Want", p1, "to be higher than", p2)
		}
	}
	invalidPorts := map[string]string{
		"0":  "1",
		"-1": "0",
		"1":  "65535",
	}
	for p1, p2 := range invalidPorts {
		if portHigher(p1, p2) {
			t.Fatal("Want", p2, "to be higher than", p1)
		}
	}
}
