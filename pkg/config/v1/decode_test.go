package v1

import (
	"testing"
)

func TestDecodeClientConfigJSON_ServerPortString(t *testing.T) {
	input := []byte(`{"serverPort": "https://example.com/port", "serverAddr": "example.com"}`)
	cfg, err := DecodeClientConfigJSON(input, DecodeOptions{DisallowUnknownFields: true})
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if cfg.ServerPort != 0 {
		t.Fatalf("expected ServerPort 0, got %d", cfg.ServerPort)
	}
	if cfg.ServerPortSource != "https://example.com/port" {
		t.Fatalf("expected ServerPortSource https://example.com/port, got %q", cfg.ServerPortSource)
	}
}

func TestDecodeClientConfigJSON_ServerPortNumber(t *testing.T) {
	input := []byte(`{"serverPort": 7000, "serverAddr": "example.com"}`)
	cfg, err := DecodeClientConfigJSON(input, DecodeOptions{DisallowUnknownFields: true})
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if cfg.ServerPort != 7000 {
		t.Fatalf("expected ServerPort 7000, got %d", cfg.ServerPort)
	}
	if cfg.ServerPortSource != "" {
		t.Fatalf("expected empty ServerPortSource, got %q", cfg.ServerPortSource)
	}
}
