package config

import "testing"

func env(m map[string]string) func(string) string { return func(k string) string { return m[k] } }

func TestDefaults(t *testing.T) {
	s, err := Load(env(nil))
	if err != nil || s.Listen != "127.0.0.1:8080" || s.SearchRate != 60 || len(s.TrustedProxies) != 0 || s.DataDir != "data" {
		t.Fatalf("defaults = %+v err = %v", s, err)
	}
}

func TestOverrides(t *testing.T) {
	s, err := Load(env(map[string]string{
		"NOTES_LISTEN": "0.0.0.0:9000", "NOTES_SEARCH_RATE": "30",
		"NOTES_TRUSTED_PROXIES": "10.0.0.0/8, 192.168.1.10", "NOTES_DATA": "/srv/notes",
	}))
	if err != nil || s.Listen != "0.0.0.0:9000" || s.SearchRate != 30 || len(s.TrustedProxies) != 2 || s.DataDir != "/srv/notes" {
		t.Fatalf("settings = %+v err = %v", s, err)
	}
}

func TestInvalid(t *testing.T) {
	for _, m := range []map[string]string{
		{"NOTES_LISTEN": "nope"},
		{"NOTES_SEARCH_RATE": "0"},
		{"NOTES_SEARCH_RATE": "fast"},
		{"NOTES_TRUSTED_PROXIES": "proxy.example.com"},
	} {
		if _, err := Load(env(m)); err == nil {
			t.Errorf("%v: expected an error", m)
		}
	}
}
