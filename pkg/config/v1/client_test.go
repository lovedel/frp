package v1

import (
	"testing"
)

func TestClientCommonConfigComplete_ServerPortSource(t *testing.T) {
	cases := []struct {
		name     string
		cfg      *ClientCommonConfig
		wantPort int
		wantSrc  string
	}{
		{
			name:     "default",
			cfg:      &ClientCommonConfig{ServerAddr: "127.0.0.1"},
			wantPort: 7000,
			wantSrc:  "",
		},
		{
			name:     "explicit port",
			cfg:      &ClientCommonConfig{ServerAddr: "127.0.0.1", ServerPort: 7500},
			wantPort: 7500,
			wantSrc:  "",
		},
		{
			name:     "source set",
			cfg:      &ClientCommonConfig{ServerAddr: "127.0.0.1", ServerPortSource: "https://example.com"},
			wantPort: 0,
			wantSrc:  "https://example.com",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := c.cfg.Complete(); err != nil {
				t.Fatalf("Complete error: %v", err)
			}
			if c.cfg.ServerPort != c.wantPort {
				t.Fatalf("ServerPort=%d, want %d", c.cfg.ServerPort, c.wantPort)
			}
			if c.cfg.ServerPortSource != c.wantSrc {
				t.Fatalf("ServerPortSource=%q, want %q", c.cfg.ServerPortSource, c.wantSrc)
			}
		})
	}
}
