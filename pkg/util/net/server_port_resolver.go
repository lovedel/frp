// Copyright 2026 The frp Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package net

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const serverPortResolveTimeout = 5 * time.Second

// ResolveServerPortSource resolves a dynamic server port source.
// The source can be:
//   - A plain integer string (returned as-is).
//   - An HTTP(S) URL whose response body contains the port number.
//   - A domain name whose TXT record contains the port number.
func ResolveServerPortSource(ctx context.Context, source string) (int, error) {
	source = strings.TrimSpace(source)
	if source == "" {
		return 0, fmt.Errorf("empty server port source")
	}

	if port, err := strconv.Atoi(source); err == nil {
		if port <= 0 || port > 65535 {
			return 0, fmt.Errorf("server port %d out of range", port)
		}
		return port, nil
	}

	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		return resolveServerPortFromURL(ctx, source)
	}

	return resolveServerPortFromTXT(ctx, source)
}

func resolveServerPortFromURL(ctx context.Context, url string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, serverPortResolveTimeout)
	defer cancel()

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		},
		Timeout: serverPortResolveTimeout,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("create request for server port url: %w", err)
	}
	req.Header.Set("User-Agent", "frpc")

	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("fetch server port from url: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, fmt.Errorf("server port url returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 128))
	if err != nil {
		return 0, fmt.Errorf("read server port url body: %w", err)
	}

	port, err := parsePortString(string(body))
	if err != nil {
		return 0, fmt.Errorf("parse server port from url: %w", err)
	}
	return port, nil
}

func resolveServerPortFromTXT(ctx context.Context, domain string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, serverPortResolveTimeout)
	defer cancel()

	resolver := net.DefaultResolver
	records, err := resolver.LookupTXT(ctx, domain)
	if err != nil {
		return 0, fmt.Errorf("lookup txt records for %s: %w", domain, err)
	}

	for _, record := range records {
		port, err := parsePortString(record)
		if err == nil {
			return port, nil
		}
	}

	return 0, fmt.Errorf("no valid port found in txt records for %s", domain)
}

func parsePortString(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty port string")
	}
	port, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid port string %q: %w", s, err)
	}
	if port <= 0 || port > 65535 {
		return 0, fmt.Errorf("port %d out of range", port)
	}
	return port, nil
}
