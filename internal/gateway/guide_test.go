package gateway

import (
	"runtime"
	"strings"
	"testing"
)

func TestDeviceGuideIsCompactTable(t *testing.T) {
	out := DeviceGuide(Status{
		LocalIP: "192.168.12.100",
		Router:  "192.168.12.1",
	}, 17890)

	// 公共断言：参数表头 + 本机 IP + 代理端口 + 「填代理」字眼必须出现。
	for _, want := range []string{
		"参数",
		"接入方式",
		"填代理",
		"主机=192.168.12.100",
		"端口=17890",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("guide missing %q:\n%s", want, out)
		}
	}

	// 平台特定断言：Mac 上指引按 mode 分场景，强调「网关 + DNS 都要改」；
	// 其它 unix 平台沿用紧凑的"改网关 / 填代理 / 本机使用"三行表格。
	switch runtime.GOOS {
	case "darwin":
		for _, want := range []string{
			"TUN 旁路由",
			"端口模式",
			"DNS 也必须",
			"改网关 + 改 DNS",
		} {
			if !strings.Contains(out, want) {
				t.Fatalf("darwin guide missing %q:\n%s", want, out)
			}
		}
	case "windows":
		if !strings.Contains(out, "Windows 不支持改网关") {
			t.Fatalf("windows guide missing Windows-specific copy:\n%s", out)
		}
	default:
		for _, want := range []string{"改网关", "停止 gateway 会自动恢复本机 DNS"} {
			if !strings.Contains(out, want) {
				t.Fatalf("unix guide missing %q:\n%s", want, out)
			}
		}
	}

	if strings.Contains(out, "验证方式 3") {
		t.Fatalf("guide should stay compact, got old verbose hint:\n%s", out)
	}
}
