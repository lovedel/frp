package legacy

import (
	"testing"
)

func TestUnmarshalClientConfFromIni_ServerPortString(t *testing.T) {
	content := []byte(`
[common]
server_addr = example.com
server_port = https://example.com/port
`)
	cfg, err := UnmarshalClientConfFromIni(content)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if cfg.ServerPort != "https://example.com/port" {
		t.Fatalf("expected ServerPort https://example.com/port, got %q", cfg.ServerPort)
	}
}

func TestUnmarshalClientConfFromIni_ServerPortNumber(t *testing.T) {
	content := []byte(`
[common]
server_addr = example.com
server_port = 7000
`)
	cfg, err := UnmarshalClientConfFromIni(content)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if cfg.ServerPort != "7000" {
		t.Fatalf("expected ServerPort 7000, got %q", cfg.ServerPort)
	}
}
