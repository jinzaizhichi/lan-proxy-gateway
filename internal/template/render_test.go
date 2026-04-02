package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tght/lan-proxy-gateway/internal/config"
)

func TestRenderTemplateIncludesInterfaceScopedTunWhenBypassLocalEnabled(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Runtime.Tun.Enabled = true
	cfg.Runtime.Tun.BypassLocal = true
	cfg.Proxy.SubscriptionURL = "https://example.com/sub"

	outPath := filepath.Join(t.TempDir(), "config.yaml")
	if err := RenderTemplate(cfg, "en0", "192.168.12.100", outPath); err != nil {
		t.Fatalf("RenderTemplate() error = %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	rendered := string(data)

	if !strings.Contains(rendered, "include-interface:\n    - en0") {
		t.Fatalf("rendered config missing include-interface block:\n%s", rendered)
	}
}

func TestRenderTemplateOmitsInterfaceScopedTunWhenBypassLocalDisabled(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Runtime.Tun.Enabled = true
	cfg.Runtime.Tun.BypassLocal = false
	cfg.Proxy.SubscriptionURL = "https://example.com/sub"

	outPath := filepath.Join(t.TempDir(), "config.yaml")
	if err := RenderTemplate(cfg, "en0", "192.168.12.100", outPath); err != nil {
		t.Fatalf("RenderTemplate() error = %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	rendered := string(data)

	if strings.Contains(rendered, "include-interface:") {
		t.Fatalf("rendered config should not contain include-interface block:\n%s", rendered)
	}
}
