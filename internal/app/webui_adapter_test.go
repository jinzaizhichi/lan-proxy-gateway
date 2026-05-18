package app

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tght/lan-proxy-gateway/internal/config"
	"github.com/tght/lan-proxy-gateway/internal/gateway"
	"github.com/tght/lan-proxy-gateway/internal/webui"
)

// newCtrlForTest 构造一个最小可用的 WebUIController，cfg 落到 t.TempDir。
// 不起 Engine，避免 saveAndReload 因 Engine.Running()=false 直接返回，便于单测纯
// 校验逻辑。
func newCtrlForTest(t *testing.T) *WebUIController {
	t.Helper()
	dir := t.TempDir()
	cfg := config.Default()
	config.Normalize(cfg)
	a := &App{
		Cfg:     cfg,
		Paths:   config.Paths{Root: dir, ConfigFile: filepath.Join(dir, "gateway.yaml"), MihomoDir: dir, CacheDir: dir},
		Gateway: gateway.New(),
	}
	return NewWebUIController(a)
}

// ─────────── P0-4: SetPorts 校验 ───────────

func TestSetPorts_RejectsOutOfRange(t *testing.T) {
	c := newCtrlForTest(t)
	if err := c.SetPorts(context.Background(), webui.Ports{Mixed: 70000}); err == nil {
		t.Fatal("Mixed=70000 必须被拒")
	}
	if err := c.SetPorts(context.Background(), webui.Ports{API: 80}); err == nil {
		t.Fatal("API=80 必须被拒（< 1024）")
	}
	if err := c.SetPorts(context.Background(), webui.Ports{Redir: -1}); err == nil {
		t.Fatal("Redir=-1 必须被拒")
	}
}

func TestSetPorts_AllowsDNSPrivilegedPort(t *testing.T) {
	c := newCtrlForTest(t)
	// DNS 默认 53 是特权端口，但 gateway 一直允许（已经是默认值）
	if err := c.SetPorts(context.Background(), webui.Ports{DNS: 53}); err != nil {
		t.Fatalf("DNS=53 应允许；got %v", err)
	}
}

func TestSetPorts_RejectsConflict(t *testing.T) {
	c := newCtrlForTest(t)
	// 把 Mixed 改成跟 API 一样
	target := c.app.Cfg.Runtime.Ports.API
	err := c.SetPorts(context.Background(), webui.Ports{Mixed: target})
	if err == nil || !strings.Contains(err.Error(), "冲突") {
		t.Fatalf("端口冲突必须报错；got %v", err)
	}
}

func TestSetPorts_RejectsChangingWebUI(t *testing.T) {
	c := newCtrlForTest(t)
	cur := c.app.Cfg.Runtime.Ports.WebUI
	err := c.SetPorts(context.Background(), webui.Ports{WebUI: cur + 1})
	if err == nil {
		t.Fatal("从 Web 端改 WebUI 端口必须被拒（会让 HTTP server 自杀）")
	}
}

func TestSetPorts_ZeroMeansUnchanged(t *testing.T) {
	c := newCtrlForTest(t)
	orig := c.app.Cfg.Runtime.Ports.Mixed
	if err := c.SetPorts(context.Background(), webui.Ports{Redir: 17893}); err != nil {
		t.Fatalf("SetPorts: %v", err)
	}
	if c.app.Cfg.Runtime.Ports.Mixed != orig {
		t.Fatalf("Mixed=0 应保持原值；before=%d after=%d", orig, c.app.Cfg.Runtime.Ports.Mixed)
	}
	if c.app.Cfg.Runtime.Ports.Redir != 17893 {
		t.Fatalf("Redir 应被改成 17893；got %d", c.app.Cfg.Runtime.Ports.Redir)
	}
}

// ─────────── P0-4: SetScript custom 路径校验 ───────────

func TestSetScript_CustomRejectsRelativePath(t *testing.T) {
	c := newCtrlForTest(t)
	err := c.SetScript(context.Background(), webui.ScriptPayload{
		Mode: "custom", CustomPath: "./foo.js",
	})
	if err == nil || !strings.Contains(err.Error(), "绝对路径") {
		t.Fatalf("相对路径必须被拒；got %v", err)
	}
}

func TestSetScript_CustomRejectsPathTraversal(t *testing.T) {
	c := newCtrlForTest(t)
	err := c.SetScript(context.Background(), webui.ScriptPayload{
		Mode: "custom", CustomPath: "/etc/../etc/passwd",
	})
	if err == nil || !strings.Contains(err.Error(), "..") {
		t.Fatalf("含 .. 必须被拒；got %v", err)
	}
}

func TestSetScript_CustomRejectsNonJSExtension(t *testing.T) {
	c := newCtrlForTest(t)
	tmp := filepath.Join(t.TempDir(), "evil.sh")
	if err := os.WriteFile(tmp, []byte("#!/bin/sh\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := c.SetScript(context.Background(), webui.ScriptPayload{
		Mode: "custom", CustomPath: tmp,
	})
	if err == nil || !strings.Contains(err.Error(), ".js") {
		t.Fatalf("非 .js 文件必须被拒；got %v", err)
	}
}

func TestSetScript_CustomRejectsMissingFile(t *testing.T) {
	c := newCtrlForTest(t)
	err := c.SetScript(context.Background(), webui.ScriptPayload{
		Mode: "custom", CustomPath: "/nonexistent/path/to/script.js",
	})
	if err == nil {
		t.Fatal("不存在的文件必须被拒")
	}
}

func TestSetScript_CustomAcceptsValidPath(t *testing.T) {
	c := newCtrlForTest(t)
	tmp := filepath.Join(t.TempDir(), "ok.js")
	if err := os.WriteFile(tmp, []byte("// ok"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := c.SetScript(context.Background(), webui.ScriptPayload{
		Mode: "custom", CustomPath: tmp,
	})
	if err != nil {
		t.Fatalf("合法 .js 路径应通过；got %v", err)
	}
	if c.app.Cfg.Source.ScriptPath != tmp {
		t.Fatalf("ScriptPath 应被写成 %q；got %q", tmp, c.app.Cfg.Source.ScriptPath)
	}
}

// ─────────── P0-2: Snapshot RLock + Set* Lock 并发不 panic ───────────
// 没法精确断言 race，只能跑 -race 看；这里至少证明在大量并发下不 panic、不死锁。
func TestController_ConcurrentReadsAndWritesDoNotPanic(t *testing.T) {
	c := newCtrlForTest(t)
	done := make(chan struct{})
	go func() {
		for i := 0; i < 50; i++ {
			_ = c.Snapshot()
		}
		done <- struct{}{}
	}()
	go func() {
		for i := 0; i < 50; i++ {
			_ = c.SetAdblock(context.Background(), i%2 == 0)
		}
		done <- struct{}{}
	}()
	<-done
	<-done
}
