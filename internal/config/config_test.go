package config

import (
	"os"
	"testing"
)

func TestDefault_ReturnsExpectedDefaults(t *testing.T) {
	os.Unsetenv("ASD_ADDR")
	os.Unsetenv("ASD_DB_PATH")

	cfg := Default()

	if cfg.Addr != "127.0.0.1:8080" {
		t.Errorf("expected default addr 127.0.0.1:8080, got %q", cfg.Addr)
	}
	if cfg.DBPath != "asd.db" {
		t.Errorf("expected default DBPath asd.db, got %q", cfg.DBPath)
	}
}

func TestDefault_RespectsEnvVars(t *testing.T) {
	tests := []struct {
		name     string
		envKey   string
		envVal   string
		checkFn  func(Config) string
		expected string
	}{
		{
			name:     "ASD_ADDR override",
			envKey:   "ASD_ADDR",
			envVal:   "0.0.0.0:9999",
			checkFn:  func(c Config) string { return c.Addr },
			expected: "0.0.0.0:9999",
		},
		{
			name:     "ASD_DB_PATH override",
			envKey:   "ASD_DB_PATH",
			envVal:   "/tmp/test.db",
			checkFn:  func(c Config) string { return c.DBPath },
			expected: "/tmp/test.db",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.envKey, tt.envVal)
			defer os.Unsetenv(tt.envKey)

			cfg := Default()
			got := tt.checkFn(cfg)
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestDefault_EmptyEnvFallsBackToDefault(t *testing.T) {
	os.Setenv("ASD_ADDR", "")
	defer os.Unsetenv("ASD_ADDR")

	cfg := Default()
	if cfg.Addr != "127.0.0.1:8080" {
		t.Errorf("empty env should fallback to default, got %q", cfg.Addr)
	}
}
