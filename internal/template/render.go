package template

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	embed "github.com/tght/lan-proxy-gateway/embed"
	"github.com/tght/lan-proxy-gateway/internal/config"
	"github.com/tght/lan-proxy-gateway/internal/rules"
	"github.com/tght/lan-proxy-gateway/internal/script"
)

// RenderTemplate replaces {{VARIABLE}} placeholders with actual values
// and writes the result to outputPath.
func RenderTemplate(cfg *config.Config, iface, ip, outputPath string) error {
	result := embed.TemplateContent

	tunConfig := "tun:\n  enable: false"
	if cfg.Runtime.Tun.Enabled {
		tunConfig = "tun:\n  enable: true\n  stack: mixed\n  auto-route: true\n  auto-detect-interface: true\n  mtu: 1500"
		if cfg.Runtime.Tun.BypassLocal && iface != "" {
			tunConfig += "\n  include-interface:\n    - " + iface
		}
	}

	replacements := map[string]string{
		"{{MIXED_PORT}}":        strconv.Itoa(cfg.Runtime.Ports.Mixed),
		"{{REDIR_PORT}}":        strconv.Itoa(cfg.Runtime.Ports.Redir),
		"{{API_PORT}}":          strconv.Itoa(cfg.Runtime.Ports.API),
		"{{API_SECRET}}":        cfg.Runtime.APISecret,
		"{{DNS_LISTEN_PORT}}":   strconv.Itoa(cfg.Runtime.Ports.DNS),
		"{{SUBSCRIPTION_URL}}":  cfg.Proxy.SubscriptionURL,
		"{{SUBSCRIPTION_NAME}}": cfg.Proxy.SubscriptionName,
		"{{LAN_INTERFACE}}":     iface,
		"{{LAN_IP}}":            ip,
		"{{TUN_CONFIG}}":        tunConfig,
		"{{RULES_BLOCK}}":       rules.Render(cfg),
	}

	for placeholder, value := range replacements {
		result = strings.ReplaceAll(result, placeholder, value)
	}

	// For file mode: patch proxy-providers from http to file type
	if cfg.Proxy.Source == "file" {
		result = patchForFileMode(result)
	}

	output := []byte(result)

	switch cfg.Extension.Mode {
	case "chains":
		if cfg.Extension.ResidentialChain == nil {
			return fmt.Errorf("extension.mode 为 chains 但未配置 extension.residential_chain")
		}
		modified, err := script.ApplyChains(cfg.Extension.ResidentialChain, output)
		if err != nil {
			return fmt.Errorf("链式代理注入失败: %w", err)
		}
		output = modified
	case "script":
		if cfg.Extension.ScriptPath == "" {
			return fmt.Errorf("extension.mode 为 script 但未配置 extension.script_path")
		}
		modified, err := script.Apply(cfg.Extension.ScriptPath, output)
		if err != nil {
			return fmt.Errorf("扩展脚本执行失败: %w", err)
		}
		output = modified
	}

	return os.WriteFile(outputPath, output, 0644)
}

// patchForFileMode modifies the generated config to use local file
// instead of HTTP subscription.
func patchForFileMode(content string) string {
	lines := strings.Split(content, "\n")
	var result []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Change type: http to type: file
		if trimmed == "type: http" {
			line = strings.Replace(line, "type: http", "type: file", 1)
		}
		// Remove url: line within proxy-providers
		if strings.HasPrefix(trimmed, "url: \"") {
			continue
		}
		// Remove interval: 3600 line
		if trimmed == "interval: 3600" {
			continue
		}
		result = append(result, line)
	}

	return strings.Join(result, "\n")
}
