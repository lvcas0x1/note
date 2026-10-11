// Package config loads server settings from the environment.
package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// Defaults.
const (
	DefaultListen     = "127.0.0.1:8080"
	DefaultSearchRate = 60
	DefaultData       = "data"
)

// Server holds the settings for `notes serve`.
type Server struct {
	Listen         string   // NOTES_LISTEN
	TrustedProxies []string // NOTES_TRUSTED_PROXIES, comma-separated IPs or CIDRs
	SearchRate     int      // NOTES_SEARCH_RATE, search requests per IP per minute
	DataDir        string   // NOTES_DATA
}

// Load reads the server settings. getenv is os.Getenv in production.
func Load(getenv func(string) string) (Server, error) {
	if getenv == nil {
		getenv = os.Getenv
	}
	s := Server{Listen: DefaultListen, SearchRate: DefaultSearchRate, DataDir: DefaultData}
	if v := strings.TrimSpace(getenv("NOTES_LISTEN")); v != "" {
		if _, _, err := net.SplitHostPort(v); err != nil {
			return s, fmt.Errorf("NOTES_LISTEN: %w", err)
		}
		s.Listen = v
	}
	if v := strings.TrimSpace(getenv("NOTES_SEARCH_RATE")); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			return s, fmt.Errorf("NOTES_SEARCH_RATE: want a positive integer, got %q", v)
		}
		s.SearchRate = n
	}
	for _, p := range strings.Split(getenv("NOTES_TRUSTED_PROXIES"), ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if net.ParseIP(p) == nil {
			if _, _, err := net.ParseCIDR(p); err != nil {
				return s, fmt.Errorf("NOTES_TRUSTED_PROXIES: %q is not an IP or CIDR", p)
			}
		}
		s.TrustedProxies = append(s.TrustedProxies, p)
	}
	if v := strings.TrimSpace(getenv("NOTES_DATA")); v != "" {
		s.DataDir = v
	}
	return s, nil
}
