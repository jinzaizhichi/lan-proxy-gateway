package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	ProxyModeRule      = "rule"
	ProxyModeGlobal    = "global"
	ProxyModeGlobalISP = "global_isp"
	ProxyModeAIProxy   = "ai_proxy"
)

type Config struct {
	ProxySource      string            `yaml:"proxy_source"`
	ProxyMode        string            `yaml:"proxy_mode,omitempty"`
	SubscriptionURL  string            `yaml:"subscription_url,omitempty"`
	ProxyConfigFile  string            `yaml:"proxy_config_file,omitempty"`
	SubscriptionName string            `yaml:"subscription_name"`
	Ports            PortsConfig       `yaml:"ports"`
	APISecret        string            `yaml:"api_secret,omitempty"`
	ChainProxy       *ChainProxyConfig `yaml:"chain_proxy,omitempty"`
}

func (c *Config) EffectiveProxyMode() string {
	if c.ProxyMode == "" {
		return ProxyModeRule
	}
	return c.ProxyMode
}

func (c *Config) ValidateProxyMode() error {
	mode := c.EffectiveProxyMode()
	switch mode {
	case ProxyModeRule, ProxyModeGlobal:
		return nil
	case ProxyModeGlobalISP, ProxyModeAIProxy:
		if c.ChainProxy == nil || !c.ChainProxy.Enabled {
			return fmt.Errorf("%s 模式需要启用 chain_proxy 配置", mode)
		}
		return nil
	default:
		return fmt.Errorf("不支持的代理模式: %s (可选: rule, global, global_isp, ai_proxy)", mode)
	}
}

type PortsConfig struct {
	Mixed int `yaml:"mixed"`
	Redir int `yaml:"redir"`
	API   int `yaml:"api"`
	DNS   int `yaml:"dns"`
}

type ChainProxyConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Server   string `yaml:"server"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	UDP      bool   `yaml:"udp"`
}

func DefaultConfig() *Config {
	return &Config{
		ProxySource:      "url",
		SubscriptionName: "subscription",
		Ports: PortsConfig{
			Mixed: 7890,
			Redir: 7892,
			API:   9090,
			DNS:   53,
		},
	}
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func Save(cfg *Config, path string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	header := []byte("# lan-proxy-gateway 配置文件\n# 此文件包含敏感信息，请勿提交到 Git\n\n")
	return os.WriteFile(path, append(header, data...), 0600)
}
