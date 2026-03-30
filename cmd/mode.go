package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tght/lan-proxy-gateway/internal/config"
	"github.com/tght/lan-proxy-gateway/internal/platform"
	"github.com/tght/lan-proxy-gateway/internal/proxy"
	tmpl "github.com/tght/lan-proxy-gateway/internal/template"
	"github.com/tght/lan-proxy-gateway/internal/ui"
)

var modeDescriptions = map[string]string{
	config.ProxyModeRule:      "规则分流（国内直连，国外代理）",
	config.ProxyModeGlobal:    "全局机场代理（所有流量走机场节点）",
	config.ProxyModeGlobalISP: "全局住宅ISP（所有流量经机场→住宅IP出口）",
	config.ProxyModeAIProxy:   "AI工作流代理（AI服务走住宅IP，阿里系直连，其余走住宅IP）",
}

var allModes = []string{
	config.ProxyModeRule,
	config.ProxyModeGlobal,
	config.ProxyModeGlobalISP,
	config.ProxyModeAIProxy,
}

var modeCmd = &cobra.Command{
	Use:   "mode [rule|global|global_isp|ai_proxy]",
	Short: "切换代理模式",
	Long: `切换代理模式。

可用模式:
  rule        规则分流模式（默认），国内直连、国外走代理
  global      全局机场代理，所有流量通过机场节点
  global_isp  全局住宅ISP，所有流量经由 机场→住宅IP 链式代理
  ai_proxy    AI工作流代理，AI服务强制走住宅IP，阿里系办公直连，
              其余流量也走住宅IP（需要配置 chain_proxy）

用法:
  gateway mode              # 查看当前模式
  gateway mode rule         # 切换到规则分流
  gateway mode global       # 切换到全局机场
  gateway mode global_isp   # 切换到全局住宅ISP
  gateway mode ai_proxy     # 切换到AI工作流代理`,
	Args: cobra.MaximumNArgs(1),
	Run:  runMode,
}

func init() {
	rootCmd.AddCommand(modeCmd)
}

func runMode(cmd *cobra.Command, args []string) {
	cfgPath := resolveConfigPath()
	cfg := loadConfigOrDefault()

	if len(args) == 0 {
		showCurrentMode(cfg)
		return
	}

	target := args[0]
	oldMode := cfg.EffectiveProxyMode()
	cfg.ProxyMode = target

	if err := cfg.ValidateProxyMode(); err != nil {
		ui.Error("%s", err)
		os.Exit(1)
	}

	if cfgPath == ".secret" {
		cfgPath = "gateway.yaml"
	}
	if err := config.Save(cfg, cfgPath); err != nil {
		ui.Error("保存配置失败: %s", err)
		os.Exit(1)
	}

	if oldMode == target {
		ui.Info("当前已是 %s 模式", target)
	} else {
		ui.Success("代理模式已切换: %s → %s", oldMode, target)
	}
	if desc, ok := modeDescriptions[target]; ok {
		fmt.Printf("  %s\n", color.New(color.Faint).Sprint(desc))
	}

	fmt.Println()
	fmt.Print("是否立即重新生成配置？[Y/n] ")
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(answer)

	if answer == "" || strings.ToLower(answer) == "y" {
		regenerateConfig(cfg)
	}
}

func showCurrentMode(cfg *config.Config) {
	mode := cfg.EffectiveProxyMode()
	fmt.Println()
	fmt.Printf("  %s %s\n", color.New(color.Bold).Sprint("当前模式:"), color.CyanString(mode))
	if desc, ok := modeDescriptions[mode]; ok {
		fmt.Printf("  %s %s\n", color.New(color.Bold).Sprint("说明:    "), desc)
	}
	if cfg.ChainProxy != nil && cfg.ChainProxy.Enabled {
		fmt.Printf("  %s %s (%s:%d)\n",
			color.New(color.Bold).Sprint("链式代理:"),
			cfg.ChainProxy.Name, cfg.ChainProxy.Server, cfg.ChainProxy.Port)
	}
	fmt.Println()
	fmt.Println("  可用模式:")
	for _, m := range allModes {
		marker := "  "
		if m == mode {
			marker = color.GreenString("▸ ")
		}
		fmt.Printf("    %s%-12s %s\n", marker, m, color.New(color.Faint).Sprint(modeDescriptions[m]))
	}
	fmt.Println()
	fmt.Printf("  %s\n", color.New(color.Faint).Sprint("切换: gateway mode [rule|global|global_isp|ai_proxy]"))
	fmt.Println()
}

func regenerateConfig(cfg *config.Config) {
	p := platform.New()
	dDir := ensureDataDir()
	iface, _ := p.DetectDefaultInterface()
	ip, _ := p.DetectInterfaceIP(iface)

	if cfg.ProxySource == "file" {
		providerFile := filepath.Join(dDir, "proxy_provider", cfg.SubscriptionName+".yaml")
		count, err := proxy.ExtractProxies(cfg.ProxyConfigFile, providerFile)
		if err != nil {
			ui.Error("提取代理节点失败: %s", err)
			os.Exit(1)
		}
		ui.Success("已提取 %d 个代理节点", count)
	}

	configPath := filepath.Join(dDir, "config.yaml")
	if err := tmpl.RenderTemplate(cfg, iface, ip, configPath); err != nil {
		ui.Error("配置生成失败: %s", err)
		os.Exit(1)
	}
	ui.Success("配置文件已生成")
	fmt.Println()
	ui.Info("如需生效，请重启网关: sudo gateway restart")
}
