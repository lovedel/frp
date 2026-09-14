package net

import (
	"context"

	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParsePortString(t *testing.T) {
	cases := []struct {
		in      string
		want    int
		wantErr bool
	}{
		{"7000", 7000, false},
		{" 7001\n", 7001, false},
		{"0", 0, true},
		{"65536", 0, true},
		{"abc", 0, true},
		{"", 0, true},
	}

	for _, c := range cases {
		got, err := parsePortString(c.in)
		if (err != nil) != c.wantErr {
			t.Fatalf("parsePortString(%q) err=%v, wantErr=%v", c.in, err, c.wantErr)
		}
		if !c.wantErr && got != c.want {
			t.Fatalf("parsePortString(%q)=%d, want %d", c.in, got, c.want)
		}
	}
}

func TestResolveServerPortSource_URL(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(" 8080\n"))
	}))
	defer s.Close()

	port, err := ResolveServerPortSource(context.Background(), s.URL)
	if err != nil {
		t.Fatalf("ResolveServerPortSource error: %v", err)
	}
	if port != 8080 {
		t.Fatalf("got port %d, want 8080", port)
	}
}

func TestResolveServerPortSource_URLBadBody(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("not-a-port"))
	}))
	defer s.Close()

	if _, err := ResolveServerPortSource(context.Background(), s.URL); err == nil {
		t.Fatalf("expected error for invalid port body")
	}
}

func TestResolveServerPortSource_URLNon2xx(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer s.Close()

	if _, err := ResolveServerPortSource(context.Background(), s.URL); err == nil {
		t.Fatalf("expected error for non-2xx status")
	}
}

func TestResolveServerPortSource_DNSLookupFail(t *testing.T) {
	_, err := ResolveServerPortSource(context.Background(), "nonexist.invalid.not_exist")
	if err == nil {
		t.Fatalf("expected lookup error for unknown domain")
	}
}

func TestResolveServerPortSource_PlaintextInt(t *testing.T) {
	port, err := ResolveServerPortSource(context.Background(), "9999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if port != 9999 {
		t.Fatalf("got port %d, want 9999", port)
	}
}
