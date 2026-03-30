package template

import (
	"strings"
	"testing"

	"github.com/tght/lan-proxy-gateway/internal/config"
)

// @author buchi
// @since 2026-03-30

func newTestConfig() *config.Config {
	return &config.Config{
		ProxySource:      "url",
		SubscriptionURL:  "https://example.com/sub",
		SubscriptionName: "test-sub",
		Ports: config.PortsConfig{
			Mixed: 7890,
			Redir: 7892,
			API:   9090,
			DNS:   53,
		},
	}
}

func newChainProxy() *config.ChainProxyConfig {
	return &config.ChainProxyConfig{
		Enabled:  true,
		Name:     "isp-proxy",
		Type:     "socks5",
		Server:   "1.2.3.4",
		Port:     1080,
		Username: "user",
		Password: "pass",
		UDP:      false,
	}
}

func TestRenderRuleMode(t *testing.T) {
	cfg := newTestConfig()
	cfg.ProxyMode = config.ProxyModeRule

	result, err := renderToString(cfg, "en0", "192.168.1.100")
	if err != nil {
		t.Fatalf("renderToString() error = %v", err)
	}

	checks := []struct {
		desc     string
		contains bool
		substr   string
	}{
		{"should have mode: rule", true, "mode: rule"},
		{"should have full rule list", true, "DOMAIN-KEYWORD,google,Proxy"},
		{"should have MATCH,Proxy", true, "MATCH,Proxy"},
		{"should have Nintendo rules", true, "nintendo.net"},
		{"should have CN direct rules", true, "GEOIP,CN,DIRECT"},
		{"should not have Global ISP group", false, "Global ISP"},
	}

	for _, c := range checks {
		got := strings.Contains(result, c.substr)
		if got != c.contains {
			t.Errorf("%s: Contains(%q) = %v, want %v", c.desc, c.substr, got, c.contains)
		}
	}
}

func TestRenderRuleModeWithChainProxy(t *testing.T) {
	cfg := newTestConfig()
	cfg.ProxyMode = config.ProxyModeRule
	cfg.ChainProxy = newChainProxy()

	result, err := renderToString(cfg, "en0", "192.168.1.100")
	if err != nil {
		t.Fatalf("renderToString() error = %v", err)
	}

	checks := []struct {
		desc     string
		contains bool
		substr   string
	}{
		{"should have chain proxy node", true, "name: isp-proxy"},
		{"should have dialer-proxy", true, "dialer-proxy: Proxy"},
		{"should have AI + Foreign group", true, "AI + Foreign"},
		{"should have AI domain rules", true, "anthropic.com,AI + Foreign"},
		{"should still have full rules", true, "GEOIP,CN,DIRECT"},
		{"should have MATCH,Proxy at end", true, "MATCH,Proxy"},
	}

	for _, c := range checks {
		got := strings.Contains(result, c.substr)
		if got != c.contains {
			t.Errorf("%s: Contains(%q) = %v, want %v", c.desc, c.substr, got, c.contains)
		}
	}
}

func TestRenderGlobalMode(t *testing.T) {
	cfg := newTestConfig()
	cfg.ProxyMode = config.ProxyModeGlobal

	result, err := renderToString(cfg, "en0", "192.168.1.100")
	if err != nil {
		t.Fatalf("renderToString() error = %v", err)
	}

	checks := []struct {
		desc     string
		contains bool
		substr   string
	}{
		{"should have MATCH,Proxy", true, "MATCH,Proxy"},
		{"should have LAN direct rules", true, "IP-CIDR,192.168.0.0/16,DIRECT"},
		{"should have global comment", true, "全局代理"},
		{"should NOT have Nintendo rules", false, "nintendo.net"},
		{"should NOT have CN direct rules", false, "GEOIP,CN,DIRECT"},
		{"should NOT have domain keyword rules", false, "DOMAIN-KEYWORD,google,Proxy"},
		{"should NOT have ad block rules", false, "admarvel,REJECT"},
	}

	for _, c := range checks {
		got := strings.Contains(result, c.substr)
		if got != c.contains {
			t.Errorf("%s: Contains(%q) = %v, want %v", c.desc, c.substr, got, c.contains)
		}
	}
}

func TestRenderGlobalISPMode(t *testing.T) {
	cfg := newTestConfig()
	cfg.ProxyMode = config.ProxyModeGlobalISP
	cfg.ChainProxy = newChainProxy()

	result, err := renderToString(cfg, "en0", "192.168.1.100")
	if err != nil {
		t.Fatalf("renderToString() error = %v", err)
	}

	checks := []struct {
		desc     string
		contains bool
		substr   string
	}{
		{"should have chain proxy node", true, "name: isp-proxy"},
		{"should have dialer-proxy", true, "dialer-proxy: Proxy"},
		{"should have Global ISP group", true, "Global ISP"},
		{"should have MATCH,Global ISP", true, "MATCH,Global ISP"},
		{"should have LAN direct rules", true, "IP-CIDR,192.168.0.0/16,DIRECT"},
		{"should NOT have Nintendo rules", false, "nintendo.net"},
		{"should NOT have CN direct rules", false, "GEOIP,CN,DIRECT"},
		{"should NOT have AI + Foreign group", false, "AI + Foreign"},
	}

	for _, c := range checks {
		got := strings.Contains(result, c.substr)
		if got != c.contains {
			t.Errorf("%s: Contains(%q) = %v, want %v", c.desc, c.substr, got, c.contains)
		}
	}
}

func TestRenderGlobalModeWithChainProxy(t *testing.T) {
	cfg := newTestConfig()
	cfg.ProxyMode = config.ProxyModeGlobal
	cfg.ChainProxy = newChainProxy()

	result, err := renderToString(cfg, "en0", "192.168.1.100")
	if err != nil {
		t.Fatalf("renderToString() error = %v", err)
	}

	if !strings.Contains(result, "MATCH,Proxy") {
		t.Error("global mode should route all to Proxy even with chain proxy configured")
	}
	if !strings.Contains(result, "name: isp-proxy") {
		t.Error("chain proxy node should still be defined")
	}
	if strings.Contains(result, "nintendo.net") {
		t.Error("should not have specific domain rules in global mode")
	}
}

func TestRenderEmptyModeDefaultsToRule(t *testing.T) {
	cfg := newTestConfig()

	result, err := renderToString(cfg, "en0", "192.168.1.100")
	if err != nil {
		t.Fatalf("renderToString() error = %v", err)
	}

	if !strings.Contains(result, "GEOIP,CN,DIRECT") {
		t.Error("empty mode should default to rule mode with full rules")
	}
	if !strings.Contains(result, "nintendo.net") {
		t.Error("empty mode should have Nintendo rules like rule mode")
	}
}

func TestPatchForGlobalModePreservesNonRulesContent(t *testing.T) {
	cfg := newTestConfig()
	cfg.ProxyMode = config.ProxyModeGlobal

	result, err := renderToString(cfg, "en0", "192.168.1.100")
	if err != nil {
		t.Fatalf("renderToString() error = %v", err)
	}

	essentials := []string{
		"mixed-port: 7890",
		"tun:",
		"enable: true",
		"dns:",
		"fake-ip",
		"proxy-providers:",
		"proxy-groups:",
		"name: Proxy",
		"name: Auto",
		"name: Fallback",
	}
	for _, s := range essentials {
		if !strings.Contains(result, s) {
			t.Errorf("global mode should preserve %q in output", s)
		}
	}
}

func TestRenderAIProxyMode(t *testing.T) {
	cfg := newTestConfig()
	cfg.ProxyMode = config.ProxyModeAIProxy
	cfg.ChainProxy = newChainProxy()

	result, err := renderToString(cfg, "en0", "192.168.1.100")
	if err != nil {
		t.Fatalf("renderToString() error = %v", err)
	}

	checks := []struct {
		desc     string
		contains bool
		substr   string
	}{
		// 链式代理基础
		{"should have chain proxy node", true, "name: isp-proxy"},
		{"should have dialer-proxy", true, "dialer-proxy: Proxy"},
		{"should have Global ISP group", true, "Global ISP"},
		{"should have MATCH,Global ISP", true, "MATCH,Global ISP"},

		// 阿里系进程直连
		{"should have AliLang process rule", true, "PROCESS-NAME,AliLang,DIRECT"},
		{"should have DingTalk process rule", true, "PROCESS-NAME,DingTalk,DIRECT"},
		{"should have IDingTalk process rule", true, "PROCESS-NAME,IDingTalk,DIRECT"},
		{"should have CloudShell process rule", true, "PROCESS-NAME,CloudShell,DIRECT"},
		{"should have aliedr process rule", true, "PROCESS-NAME,com.alibaba.endpoint.aliedr.ne,DIRECT"},

		// 阿里系域名直连
		{"should have alibaba-inc direct", true, "alibaba-inc.com,DIRECT"},
		{"should have dingtalk direct", true, "dingtalk.com,DIRECT"},
		{"should have alipay direct", true, "alipay.com,DIRECT"},
		{"should have aliyun direct", true, "aliyun.com,DIRECT"},
		{"should have alibaba.net direct", true, "alibaba.net,DIRECT"},

		// Claude → 住宅IP
		{"should have claude.ai via ISP", true, "claude.ai,Global ISP"},
		{"should have anthropic via ISP", true, "anthropic.com,Global ISP"},
		{"should have claudeusercontent via ISP", true, "claudeusercontent.com,Global ISP"},

		// OpenAI → 住宅IP
		{"should have openai via ISP", true, "openai.com,Global ISP"},
		{"should have chatgpt via ISP", true, "chatgpt.com,Global ISP"},

		// Google → 住宅IP
		{"should have google via ISP", true, "google.com,Global ISP"},
		{"should have googleapis via ISP", true, "googleapis.com,Global ISP"},
		{"should have youtube via ISP", true, "youtube.com,Global ISP"},

		// Cursor → 住宅IP
		{"should have cursor.sh via ISP", true, "cursor.sh,Global ISP"},
		{"should have cursor.com via ISP", true, "cursor.com,Global ISP"},

		// GitHub → 住宅IP
		{"should have github via ISP", true, "github.com,Global ISP"},
		{"should have githubcopilot via ISP", true, "githubcopilot.com,Global ISP"},

		// 其他 AI
		{"should have perplexity via ISP", true, "perplexity.ai,Global ISP"},
		{"should have mistral via ISP", true, "mistral.ai,Global ISP"},

		// 中国大陆直连
		{"should have GEOIP CN direct", true, "GEOIP,CN,DIRECT"},
		{"should have LAN direct", true, "IP-CIDR,192.168.0.0/16,DIRECT"},

		// 不应包含旧规则模式的内容
		{"should NOT have Nintendo rules", false, "nintendo.net"},
		{"should NOT have ad block rules", false, "admarvel,REJECT"},
		{"should NOT have AI + Foreign group", false, "AI + Foreign"},
	}

	for _, c := range checks {
		got := strings.Contains(result, c.substr)
		if got != c.contains {
			t.Errorf("%s: Contains(%q) = %v, want %v", c.desc, c.substr, got, c.contains)
		}
	}
}

func TestAIProxyModePreservesProxyGroups(t *testing.T) {
	cfg := newTestConfig()
	cfg.ProxyMode = config.ProxyModeAIProxy
	cfg.ChainProxy = newChainProxy()

	result, err := renderToString(cfg, "en0", "192.168.1.100")
	if err != nil {
		t.Fatalf("renderToString() error = %v", err)
	}

	essentials := []string{
		"mixed-port: 7890",
		"tun:",
		"dns:",
		"proxy-providers:",
		"proxy-groups:",
		"name: Proxy",
		"name: Auto",
		"name: Fallback",
		"name: Global ISP",
	}
	for _, s := range essentials {
		if !strings.Contains(result, s) {
			t.Errorf("ai_proxy mode should preserve %q in output", s)
		}
	}
}

func TestGlobalISPGroupContainsChainProxyName(t *testing.T) {
	cfg := newTestConfig()
	cfg.ProxyMode = config.ProxyModeGlobalISP
	cfg.ChainProxy = newChainProxy()

	result, err := renderToString(cfg, "en0", "192.168.1.100")
	if err != nil {
		t.Fatalf("renderToString() error = %v", err)
	}

	if !strings.Contains(result, "- isp-proxy") {
		t.Error("Global ISP group should list the chain proxy name as an option")
	}
}
