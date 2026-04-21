// Package presets 内嵌一套「开箱即用」的增强脚本模板，给 gateway 主菜单
// 「增强脚本」作引导式填表的预设用。
//
// 目前只有一个：residential-chain（链式代理 · 住宅 IP 落地），作用：
//   - 把订阅节点包成 🛫 AI起飞节点 组
//   - 把用户填的住宅 IP 代理包成 🛬 AI落地节点，用 mihomo dialer-proxy 串起来
//   - AI 相关域名（Claude/OpenAI/Cursor/Termius）路由到 AI落地节点
package presets

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tght/lan-proxy-gateway/internal/config"
)

//go:embed residential-chain.tmpl.js
var residentialChainTmpl string

// ResidentialChainFile 是预设脚本渲染到 workdir 后的目标文件名。
// engine.Render 拿到它往 cfg.Source.ScriptPath 填。
const ResidentialChainFile = "residential-chain.js"

// RenderResidentialChain 用用户填的 ChainResidentialConfig 实例化预设脚本，
// 写到 workDir/residential-chain.js，返回落盘路径。
// 占位符 __RESIDENTIAL_PROXY_JSON__ 替换成一个完整的 JS/JSON 对象字面量，
// 敏感字段（username/password）在 JSON 里原样存，同时靠 0o600 文件权限护住。
func RenderResidentialChain(cfg *config.ChainResidentialConfig, workDir string) (string, error) {
	if cfg == nil {
		return "", fmt.Errorf("ChainResidentialConfig 为空")
	}
	if cfg.Server == "" || cfg.Port <= 0 {
		return "", fmt.Errorf("住宅 IP 节点信息不完整：server=%q port=%d", cfg.Server, cfg.Port)
	}
	kind := strings.ToLower(cfg.Kind)
	if kind != "http" && kind != "socks5" {
		kind = "socks5"
	}
	dialer := cfg.DialerProxy
	if dialer == "" {
		dialer = "🛫 AI起飞节点"
	}
	name := cfg.Name
	if name == "" {
		name = "🏠 住宅IP"
	}

	obj := map[string]interface{}{
		"name":         name,
		"type":         kind,
		"server":       cfg.Server,
		"port":         cfg.Port,
		"dialer-proxy": dialer,
		"udp":          false,
	}
	if cfg.Username != "" {
		obj["username"] = cfg.Username
		obj["password"] = cfg.Password
	}

	// JSON 字面量天然是合法的 JS 对象字面量，直接塞进去。
	// 缩进让脚本好读一些。
	objJSON, err := json.MarshalIndent(obj, "    ", "  ")
	if err != nil {
		return "", fmt.Errorf("住宅 IP 节点 JSON 序列化失败: %w", err)
	}
	rendered := strings.ReplaceAll(residentialChainTmpl, "__RESIDENTIAL_PROXY_JSON__", string(objJSON))

	if err := os.MkdirAll(workDir, 0o755); err != nil {
		return "", fmt.Errorf("创建 workdir: %w", err)
	}
	dst := filepath.Join(workDir, ResidentialChainFile)
	// 脚本里含明文凭证，0o600 防止 LAN 里别的用户能读。
	if err := os.WriteFile(dst, []byte(rendered), 0o600); err != nil {
		return "", fmt.Errorf("写入预设脚本: %w", err)
	}
	return dst, nil
}
