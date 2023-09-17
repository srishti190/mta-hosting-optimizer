package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	expectedPort := 8082

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Error loading configuration: %v", err)
	}

	if cfg.Port != expectedPort {
		t.Errorf("Expected Port to be %d, got %d", expectedPort, cfg.Port)
	}
}
