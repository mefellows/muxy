package command

import "testing"

func TestCommands(t *testing.T) {

}

func TestCommands_Init(t *testing.T) {
	setup()

	proxy := Commands["proxy"]
	if proxy == nil {
		t.Fatal("Want proxy command, got nil")
	}
	proxy()

	pki := Commands["pki"]
	if Commands["pki"] == nil {
		t.Fatal("Want pki command, got nil")
	}
	pki()
}
