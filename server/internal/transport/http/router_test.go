package httpserver

import "testing"

func TestNewHumaConfigUsesRelativeServerURLWhenBackendBaseURLUnset(t *testing.T) {
	config := newHumaConfig(Dependencies{APIVersion: "v1"})

	if len(config.Servers) != 1 {
		t.Fatalf("expected one server, got %d", len(config.Servers))
	}
	if config.Servers[0].URL != "/api/v1" {
		t.Fatalf("expected relative server URL, got %q", config.Servers[0].URL)
	}
}

func TestNewHumaConfigUsesBackendBaseURLServerURL(t *testing.T) {
	config := newHumaConfig(Dependencies{
		APIVersion:     "v1",
		BackendBaseURL: "https://api.netstamp.dev/",
	})

	if len(config.Servers) != 1 {
		t.Fatalf("expected one server, got %d", len(config.Servers))
	}
	if config.Servers[0].URL != "https://api.netstamp.dev/api/v1" {
		t.Fatalf("expected absolute server URL, got %q", config.Servers[0].URL)
	}
}
