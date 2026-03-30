package config

import (
	"os"
	"path/filepath"
	"testing"
)

// @author buchi
// @since 2026-03-30

func TestEffectiveProxyMode(t *testing.T) {
	tests := []struct {
		name     string
		mode     string
		expected string
	}{
		{"empty defaults to rule", "", ProxyModeRule},
		{"rule", ProxyModeRule, ProxyModeRule},
		{"global", ProxyModeGlobal, ProxyModeGlobal},
		{"global_isp", ProxyModeGlobalISP, ProxyModeGlobalISP},
		{"ai_proxy", ProxyModeAIProxy, ProxyModeAIProxy},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{ProxyMode: tt.mode}
			if got := cfg.EffectiveProxyMode(); got != tt.expected {
				t.Errorf("EffectiveProxyMode() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestValidateProxyMode(t *testing.T) {
	chainProxy := &ChainProxyConfig{Enabled: true, Name: "test", Type: "socks5", Server: "1.2.3.4", Port: 1080}

	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name:    "empty mode (defaults to rule) is valid",
			cfg:     &Config{},
			wantErr: false,
		},
		{
			name:    "rule mode is valid",
			cfg:     &Config{ProxyMode: ProxyModeRule},
			wantErr: false,
		},
		{
			name:    "global mode is valid",
			cfg:     &Config{ProxyMode: ProxyModeGlobal},
			wantErr: false,
		},
		{
			name:    "global_isp with chain_proxy is valid",
			cfg:     &Config{ProxyMode: ProxyModeGlobalISP, ChainProxy: chainProxy},
			wantErr: false,
		},
		{
			name:    "global_isp without chain_proxy fails",
			cfg:     &Config{ProxyMode: ProxyModeGlobalISP},
			wantErr: true,
		},
		{
			name:    "global_isp with disabled chain_proxy fails",
			cfg:     &Config{ProxyMode: ProxyModeGlobalISP, ChainProxy: &ChainProxyConfig{Enabled: false}},
			wantErr: true,
		},
		{
			name:    "ai_proxy with chain_proxy is valid",
			cfg:     &Config{ProxyMode: ProxyModeAIProxy, ChainProxy: chainProxy},
			wantErr: false,
		},
		{
			name:    "ai_proxy without chain_proxy fails",
			cfg:     &Config{ProxyMode: ProxyModeAIProxy},
			wantErr: true,
		},
		{
			name:    "unknown mode fails",
			cfg:     &Config{ProxyMode: "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.ValidateProxyMode()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateProxyMode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadAndSaveWithProxyMode(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "gateway.yaml")

	cfg := DefaultConfig()
	cfg.ProxyMode = ProxyModeGlobal
	cfg.SubscriptionURL = "https://example.com/sub"

	if err := Save(cfg, path); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.ProxyMode != ProxyModeGlobal {
		t.Errorf("loaded ProxyMode = %q, want %q", loaded.ProxyMode, ProxyModeGlobal)
	}

	data, _ := os.ReadFile(path)
	content := string(data)
	if !contains(content, "proxy_mode: global") {
		t.Errorf("saved file should contain 'proxy_mode: global', got:\n%s", content)
	}
}

func TestLoadPreservesEmptyProxyMode(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "gateway.yaml")

	yamlContent := `proxy_source: url
subscription_url: "https://example.com"
subscription_name: test
ports:
  mixed: 7890
  redir: 7892
  api: 9090
  dns: 53
`
	if err := os.WriteFile(path, []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.ProxyMode != "" {
		t.Errorf("ProxyMode should be empty for legacy config, got %q", cfg.ProxyMode)
	}
	if cfg.EffectiveProxyMode() != ProxyModeRule {
		t.Errorf("EffectiveProxyMode() = %q, want %q", cfg.EffectiveProxyMode(), ProxyModeRule)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
