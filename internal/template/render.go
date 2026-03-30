package template

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	embed "github.com/tght/lan-proxy-gateway/embed"
	"github.com/tght/lan-proxy-gateway/internal/config"
)

// RenderTemplate replaces {{VARIABLE}} placeholders with actual values
// and writes the result to outputPath.
func RenderTemplate(cfg *config.Config, iface, ip, outputPath string) error {
	result, err := renderToString(cfg, iface, ip)
	if err != nil {
		return err
	}
	return os.WriteFile(outputPath, []byte(result), 0644)
}

// renderToString 执行模板渲染并返回字符串结果（便于测试）
func renderToString(cfg *config.Config, iface, ip string) (string, error) {
	result := embed.TemplateContent

	replacements := map[string]string{
		"{{MIXED_PORT}}":        strconv.Itoa(cfg.Ports.Mixed),
		"{{REDIR_PORT}}":        strconv.Itoa(cfg.Ports.Redir),
		"{{API_PORT}}":          strconv.Itoa(cfg.Ports.API),
		"{{API_SECRET}}":        cfg.APISecret,
		"{{DNS_LISTEN_PORT}}":   strconv.Itoa(cfg.Ports.DNS),
		"{{SUBSCRIPTION_URL}}":  cfg.SubscriptionURL,
		"{{SUBSCRIPTION_NAME}}": cfg.SubscriptionName,
		"{{LAN_INTERFACE}}":     iface,
		"{{LAN_IP}}":            ip,
	}

	proxySection, groupSection, rulesSection := buildChainProxyConfig(cfg)
	replacements["{{CHAIN_PROXY_SECTION}}"] = proxySection
	replacements["{{CHAIN_PROXY_GROUP}}"] = groupSection
	replacements["{{CHAIN_PROXY_RULES}}"] = rulesSection

	for placeholder, value := range replacements {
		result = strings.ReplaceAll(result, placeholder, value)
	}

	if cfg.ProxySource == "file" {
		result = patchForFileMode(result)
	}

	mode := cfg.EffectiveProxyMode()
	switch mode {
	case config.ProxyModeGlobal, config.ProxyModeGlobalISP:
		result = patchForGlobalMode(result, cfg)
	case config.ProxyModeAIProxy:
		result = patchForAIProxyMode(result, cfg)
	}

	return result, nil
}

// buildChainProxyConfig 构建链式代理的配置片段
func buildChainProxyConfig(cfg *config.Config) (proxySection, groupSection, rulesSection string) {
	if cfg.ChainProxy == nil || !cfg.ChainProxy.Enabled {
		proxySection = "proxies: []"
		groupSection = ""
		rulesSection = ""
		return
	}

	cp := cfg.ChainProxy
	proxySection = fmt.Sprintf(`proxies:
  - name: %s
    type: %s
    server: %s
    port: %d`, cp.Name, cp.Type, cp.Server, cp.Port)
	if cp.Username != "" {
		proxySection += fmt.Sprintf("\n    username: %s", cp.Username)
	}
	if cp.Password != "" {
		proxySection += fmt.Sprintf("\n    password: %s", cp.Password)
	}
	proxySection += fmt.Sprintf("\n    udp: %t", cp.UDP)
	proxySection += "\n    dialer-proxy: Proxy"

	mode := cfg.EffectiveProxyMode()

	switch mode {
	case config.ProxyModeGlobalISP, config.ProxyModeAIProxy:
		groupSection = fmt.Sprintf(`  - name: Global ISP
    type: select
    proxies:
      - %s
      - Proxy
      - DIRECT

`, cp.Name)
		rulesSection = ""
	default:
		groupSection = fmt.Sprintf(`  - name: AI + Foreign
    type: select
    proxies:
      - %s
      - Proxy
      - DIRECT

`, cp.Name)
		rulesSection = `  # --- AI + Foreign（链式代理 → 住宅 IP）---
  - DOMAIN-SUFFIX,anthropic.com,AI + Foreign
  - DOMAIN-SUFFIX,claude.ai,AI + Foreign
  - DOMAIN-SUFFIX,claudeusercontent.com,AI + Foreign
  - DOMAIN-SUFFIX,openai.com,AI + Foreign
  - DOMAIN-SUFFIX,chatgpt.com,AI + Foreign
  - DOMAIN-SUFFIX,oaiusercontent.com,AI + Foreign
  - DOMAIN-SUFFIX,oaistatic.com,AI + Foreign
  - DOMAIN,gemini.google.com,AI + Foreign
  - DOMAIN,generativelanguage.googleapis.com,AI + Foreign
  - DOMAIN,ping0.cc,AI + Foreign
  - DOMAIN,ipinfo.io,AI + Foreign
`
	}
	return
}

// patchForGlobalMode 将规则分流模式替换为全局代理模式
func patchForGlobalMode(content string, cfg *config.Config) string {
	var target string
	switch cfg.EffectiveProxyMode() {
	case config.ProxyModeGlobal:
		target = "Proxy"
	case config.ProxyModeGlobalISP:
		target = "Global ISP"
	default:
		return content
	}

	idx := strings.Index(content, "\nrules:\n")
	if idx == -1 {
		return content
	}

	globalRules := fmt.Sprintf(`
rules:
  # --- 本地/局域网直连 ---
  - DOMAIN-SUFFIX,local,DIRECT
  - IP-CIDR,127.0.0.0/8,DIRECT
  - IP-CIDR,172.16.0.0/12,DIRECT
  - IP-CIDR,192.168.0.0/16,DIRECT
  - IP-CIDR,10.0.0.0/8,DIRECT
  - IP-CIDR,100.64.0.0/10,DIRECT
  - IP-CIDR,224.0.0.0/4,DIRECT
  - IP-CIDR6,fe80::/10,DIRECT

  # --- 全局代理（所有其他流量）---
  - MATCH,%s
`, target)

	return content[:idx] + globalRules
}

// patchForAIProxyMode 生成 AI 编码工作流专用规则：
// AI 服务 → 住宅ISP代理，阿里系办公 → 直连，中国大陆 → 直连，其余 → 住宅ISP
func patchForAIProxyMode(content string, cfg *config.Config) string {
	idx := strings.Index(content, "\nrules:\n")
	if idx == -1 {
		return content
	}

	aiRules := `
rules:
  # === 进程名直连（阿里系办公应用）===
  - PROCESS-NAME,AliLang,DIRECT
  - PROCESS-NAME,IDingTalk,DIRECT
  - PROCESS-NAME,DingTalk,DIRECT
  - PROCESS-NAME,CloudShell,DIRECT
  - PROCESS-NAME,com.alibaba.endpoint.aliedr.ne,DIRECT

  # === 阿里系域名直连 ===
  - DOMAIN-SUFFIX,alibaba-inc.com,DIRECT
  - DOMAIN-SUFFIX,alibaba.net,DIRECT
  - DOMAIN-SUFFIX,dingtalk.com,DIRECT
  - DOMAIN-SUFFIX,alipay.com,DIRECT
  - DOMAIN-SUFFIX,aliyun.com,DIRECT
  - DOMAIN-SUFFIX,alicdn.com,DIRECT
  - DOMAIN-SUFFIX,alibaba.com,DIRECT
  - DOMAIN-SUFFIX,alibabacloud.com,DIRECT
  - DOMAIN-SUFFIX,taobao.com,DIRECT
  - DOMAIN-SUFFIX,tmall.com,DIRECT
  - DOMAIN-SUFFIX,aliyuncs.com,DIRECT
  - DOMAIN-SUFFIX,alidns.com,DIRECT
  - DOMAIN-KEYWORD,alibaba,DIRECT
  - DOMAIN-KEYWORD,alicdn,DIRECT
  - DOMAIN-KEYWORD,alipay,DIRECT
  - DOMAIN-KEYWORD,dingtalk,DIRECT
  - DOMAIN-KEYWORD,taobao,DIRECT

  # === Claude / Anthropic → 住宅IP代理 ===
  - DOMAIN-SUFFIX,anthropic.com,Global ISP
  - DOMAIN-SUFFIX,claude.ai,Global ISP
  - DOMAIN-SUFFIX,claudeusercontent.com,Global ISP

  # === OpenAI / ChatGPT → 住宅IP代理 ===
  - DOMAIN-SUFFIX,openai.com,Global ISP
  - DOMAIN-SUFFIX,chatgpt.com,Global ISP
  - DOMAIN-SUFFIX,oaiusercontent.com,Global ISP
  - DOMAIN-SUFFIX,oaistatic.com,Global ISP
  - DOMAIN-SUFFIX,openaiapi.com,Global ISP

  # === Google 系列 → 住宅IP代理 ===
  - DOMAIN-SUFFIX,google.com,Global ISP
  - DOMAIN-SUFFIX,googleapis.com,Global ISP
  - DOMAIN-SUFFIX,googlevideo.com,Global ISP
  - DOMAIN-SUFFIX,gstatic.com,Global ISP
  - DOMAIN-SUFFIX,googleusercontent.com,Global ISP
  - DOMAIN-SUFFIX,google.co,Global ISP
  - DOMAIN-SUFFIX,ggpht.com,Global ISP
  - DOMAIN-SUFFIX,youtube.com,Global ISP
  - DOMAIN-SUFFIX,ytimg.com,Global ISP
  - DOMAIN-SUFFIX,youtu.be,Global ISP
  - DOMAIN-SUFFIX,gmail.com,Global ISP
  - DOMAIN-SUFFIX,deepmind.com,Global ISP
  - DOMAIN-SUFFIX,deepmind.google,Global ISP
  - DOMAIN-KEYWORD,google,Global ISP

  # === Cursor IDE → 住宅IP代理 ===
  - DOMAIN-SUFFIX,cursor.sh,Global ISP
  - DOMAIN-SUFFIX,cursor.com,Global ISP
  - DOMAIN-SUFFIX,cursorapi.com,Global ISP
  - DOMAIN-SUFFIX,cursor.so,Global ISP
  - DOMAIN-KEYWORD,cursor,Global ISP

  # === GitHub / Copilot → 住宅IP代理 ===
  - DOMAIN-SUFFIX,github.com,Global ISP
  - DOMAIN-SUFFIX,githubusercontent.com,Global ISP
  - DOMAIN-SUFFIX,githubcopilot.com,Global ISP
  - DOMAIN-SUFFIX,github.io,Global ISP
  - DOMAIN-SUFFIX,githubassets.com,Global ISP
  - DOMAIN-KEYWORD,github,Global ISP
  - DOMAIN-KEYWORD,copilot,Global ISP

  # === 其他 AI 服务 → 住宅IP代理 ===
  - DOMAIN-SUFFIX,perplexity.ai,Global ISP
  - DOMAIN-SUFFIX,huggingface.co,Global ISP
  - DOMAIN-SUFFIX,midjourney.com,Global ISP
  - DOMAIN-SUFFIX,stability.ai,Global ISP
  - DOMAIN-SUFFIX,cohere.ai,Global ISP
  - DOMAIN-SUFFIX,mistral.ai,Global ISP
  - DOMAIN-SUFFIX,together.ai,Global ISP
  - DOMAIN-SUFFIX,groq.com,Global ISP
  - DOMAIN-SUFFIX,fireworks.ai,Global ISP
  - DOMAIN-SUFFIX,replicate.com,Global ISP
  - DOMAIN-SUFFIX,vercel.com,Global ISP
  - DOMAIN-SUFFIX,v0.dev,Global ISP
  - DOMAIN-SUFFIX,sentry.io,Global ISP

  # === IP 检测站 → 住宅IP代理 ===
  - DOMAIN,ping0.cc,Global ISP
  - DOMAIN,ipinfo.io,Global ISP
  - DOMAIN,ip.sb,Global ISP
  - DOMAIN,ifconfig.me,Global ISP
  - DOMAIN,myip.com,Global ISP

  # === 本地/局域网直连 ===
  - DOMAIN-SUFFIX,local,DIRECT
  - IP-CIDR,127.0.0.0/8,DIRECT
  - IP-CIDR,172.16.0.0/12,DIRECT
  - IP-CIDR,192.168.0.0/16,DIRECT
  - IP-CIDR,10.0.0.0/8,DIRECT
  - IP-CIDR,100.64.0.0/10,DIRECT
  - IP-CIDR,224.0.0.0/4,DIRECT
  - IP-CIDR6,fe80::/10,DIRECT

  # === 中国大陆直连 ===
  - DOMAIN-SUFFIX,cn,DIRECT
  - DOMAIN-KEYWORD,-cn,DIRECT
  - GEOIP,CN,DIRECT

  # === 兜底 → 住宅IP代理 ===
  - MATCH,Global ISP
`
	return content[:idx] + aiRules
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
