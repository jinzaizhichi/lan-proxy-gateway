package script

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestApply_InjectsProxyGroup(t *testing.T) {
	dir := t.TempDir()
	scriptPath := filepath.Join(dir, "enhance.js")
	if err := os.WriteFile(scriptPath, []byte(`
function main(config) {
  config["proxy-groups"] = config["proxy-groups"] || [];
  config["proxy-groups"].unshift({
    name: "AI",
    type: "select",
    proxies: ["DIRECT"]
  });
  config.rules = (config.rules || []);
  config.rules.unshift("DOMAIN-SUFFIX,openai.com,AI");
  return config;
}
`), 0o644); err != nil {
		t.Fatalf("write script: %v", err)
	}

	input := []byte(`
mode: rule
proxies: []
proxy-groups:
  - name: Proxy
    type: select
    proxies: [DIRECT]
rules:
  - MATCH,Proxy
`)
	out, err := Apply(scriptPath, input)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	s := string(out)
	if !strings.Contains(s, "name: AI") {
		t.Errorf("期望注入 AI 组，实际输出:\n%s", s)
	}
	if !strings.Contains(s, "DOMAIN-SUFFIX,openai.com,AI") {
		t.Errorf("期望注入优先规则，实际输出:\n%s", s)
	}
	// 原有结构保持
	if !strings.Contains(s, "name: Proxy") {
		t.Errorf("原 Proxy 组丢了:\n%s", s)
	}
}

func TestApply_MissingMainFn(t *testing.T) {
	dir := t.TempDir()
	sp := filepath.Join(dir, "bad.js")
	_ = os.WriteFile(sp, []byte(`function notMain(c){return c}`), 0o644)
	_, err := Apply(sp, []byte("proxies: []"))
	if err == nil || !strings.Contains(err.Error(), "main") {
		t.Errorf("期望找不到 main 的错误, got %v", err)
	}
}
