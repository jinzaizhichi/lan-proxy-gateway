# LAN Proxy Gateway v3.2.0

Windows 支持**重大修复版**。本版基于 Windows 11 实机从头走完「脚本安装 → 下载 mihomo → 配置订阅 → 启动 → 主菜单操作 → 手机/手持设备接入 → Web 控制台」全链路，修掉 10 个阻塞级 Windows bug，补齐多处 macOS/Linux 早已能用但 Windows 半残的功能。Mac/Linux 用户零回归。

## 🪟 Windows 重大修复（10 项）

| # | 位置 | 症状 | 修复 |
|---|---|---|---|
| 1 | `internal/mihomo/download.go` | 下载 mihomo 只打一次 GitHub，超时就死 | 加 4 个镜像 fallback（`ghfast.top` / `hub.gitmirror.com` / `github.moeyy.xyz` / `ghp.ci`），支持 `GITHUB_MIRROR` env 覆盖 + `HTTP_PROXY` 走已有代理；15s 拨号超时快速切下一个 |
| 2 | `internal/platform/platform_windows.go` | 装完找不到 mihomo（`%LOCALAPPDATA%\...\bin\` 不在 PATH）| `ResolveMihomoPath` 显式查真实安装目录 |
| 3 | `internal/engine/process_windows.go` | `terminateProcess` 只 TerminateProcess 主进程，TUN 子进程活着继续占 17890 | 改 `taskkill /PID x /T /F` 杀进程树 |
| 4 | `internal/engine/process_windows.go` | `pidAlive` 用 `Signal(0)` —— Windows 一律不支持，导致菜单状态永远"○ 未启动" | 改用 `tasklist /FI "PID eq x"` |
| 5 | `internal/engine/preflight.go` | `FindStaleMihomoPIDs` Windows 分支 `return nil // TODO`，"清理残留 mihomo" 按钮永远失效 | 实装 `tasklist /FI "IMAGENAME eq mihomo.exe"` |
| 6 | `internal/engine/preflight.go` | `LookupPortOwner` Windows 也是 return nil，端口冲突消息看不到占用者 | 实装 `netstat -ano` + `tasklist` 反查 PID → 进程名 |
| 7 | `internal/engine/preflight.go` | 端口冲突错误文案硬编码"主菜单→流量控制"（菜单不存在）| 改真实路径"2 分流 & 规则 → 9 高级设置" / "4 启动/重启/停止 → 4 清理残留 mihomo" |
| 8 | `internal/platform/platform_windows.go` | mihomo TUN 启动后 DetectNetwork 选中 TUN 网卡 → 主菜单显示"本机 IP: 198.18.0.1"（CGNAT）| 过滤 TUN/TAP 适配器名字 + 跳过 198.18.0.0/15 网段 |
| 9 | `cmd/install.go` / `cmd/start.go` | 尾部指引硬编码 `sudo gateway` / `sudo gateway stop`（Windows 没 sudo） | 抽 `elevatedCmd()` 按平台分支 |
| 10 | `embed/embed.go` | `//go:embed webui` 默认跳过 `_nuxt/` 和 `_fonts/`（下划线开头），Web UI 缺 2MB 资源，浏览器打开 `/ui/` **一片空白** | 改 `//go:embed all:webui` 强制打包所有文件；binary 22.5MB → 24.5MB |

## ✨ 新增 / 增强

- **下载镜像**：mihomo 二进制走 4 个镜像 + `HTTP_PROXY` env + `GITHUB_MIRROR` 覆盖，国内用户开箱即用
- **下载进度**：`Installer.Logf` 字段，用户看到"↓ 下载 mihomo v1.18.6 (直连 github) / 失败 → 镜像 ghfast.top / ✓"的实时反馈
- **启动失败提示更聪明**：`a.Start()` 返回的错误本身就带完整解决方案（如端口冲突给出具体菜单路径），install.go 不再乱甩重复/误导的"端口冲突关 DNS"建议
- **设备接入指引按平台分化**：Windows 直接说"不支持方式 1 改网关；请用方式 2 设 HTTP 代理"，并给出 Android/iOS 具体步骤，不再给出 `dscacheutil` 等 macOS 独有命令
- **主菜单操作行统一**：4 处"0 返回"补齐"（或按 Q）"提示
- **macOS `SetLocalDNSToLoopback` 说明**：菜单里提示会自动改**所有活跃网卡**（Wi-Fi + Ethernet + Thunderbolt 等），不用担心多网卡场景

## 📝 README 改版

- **突出招牌能力「AI 服务走纯净住宅 IP 链式代理」**：顶部 tagline 加 AI 场景二号位；feature list 开辟"🌟 招牌能力"独立板块，讲清 "机房 IP 被 AI 风控 → 机场起飞 + 住宅 IP 落地" 的价值；"场景三"从 3 行扩到完整教程（流量路径图、操作步骤、自动生成内容、验证方法）
- 顶部加 **Windows 限制** callout：家用版无 RRAS/NAT，方式 1 走不通；Switch/PS5 透明代理需 Linux/macOS/软路由
- 方式对比表增加 macOS / Linux / Windows 三列支持情况
- 方式 3 重写：TUN 开着 vs 改本机 DNS 两条路的取舍讲清楚；Linux 给出遍历 NetworkManager 连接的真命令；Windows 明说没实装
- 镜像列表跟代码对齐（旧文写 `ghproxy.com / ddlc.top`，代码实际是 `ghfast.top / ghp.ci`）
- 加 Clash Verge 共存 env 示例（`HTTP_PROXY=http://127.0.0.1:7897` 让 gateway 走现有代理下载）

## 🐛 非 Windows 回归修复

- `internal/engine/engine.go` `deployWebUI` 的 `filepath.Rel` 在混合分隔符场景下可能产出错误路径；改用 `path.Clean` + `strings.TrimPrefix` + `filepath.FromSlash`，加回归测试

## ✅ 测试覆盖

新增 3 组单元测试：

- `internal/mihomo/download_test.go` — 镜像候选列表、`GITHUB_MIRROR` env 覆盖、`Install` 回退路径、全失败路径（httptest 起假服务器验证）
- `internal/platform/platform_windows_test.go` — 虚拟适配器识别（Mihomo / Wintun / TAP / Tailscale）、TUN IP 网段判定（198.18.0.0/15）
- `internal/engine/deploy_webui_test.go` — Web UI 释放后 `index.html` / `_nuxt/` / `favicon.svg` 必在预期路径，且没泄漏到 workdir 外（专门抓 10. 号 bug 的回归测试）

全仓 `go test ./...` 在 Windows 11 + Go 1.25 通过。

## 🧪 真机验证

**主要测试环境：Windows 11 Pro 22H2 · Go 1.25.0**

走完全链路 ✅：
- `install.ps1` 安装 gateway.exe
- `gateway install`：下载 mihomo（直连 github 超时 → 自动切 ghfast.top 成功）、下载 geoip/geosite 数据
- 配置本地 .yaml 订阅源
- 首次启动：port 53 冲突（Clash Verge 在跑）→ 按提示 2→9→1 关 DNS 代理 → 启动成功
- Web UI `http://127.0.0.1:19090/ui/`：metacubexd 完整加载 ✓
- 手机 Android Wi-Fi 高级选项设 HTTP 代理 `192.168.12.109:17890` → 连通 ✓
- 主菜单状态行正确显示运行/未启动（修 #4 前永远显示"未启动"）
- "清理残留 mihomo" 按钮能真的杀掉僵尸进程（修 #5 前永远"没发现"）

**macOS**：本版未做针对性修改（原本在 macOS 上各项能跑通）。本次交叉编译产物通过 `go test ./...`。

**Linux**：⚠️ **本版未做真机测试**，仅交叉编译 + 单元测试通过。如果你在 Linux 上跑本版遇到问题，请开 issue —— 很欢迎 Linux 用户贡献实机验证结果。

## 下载

| 系统 | 直接下载 |
|---|---|
| macOS Apple Silicon | [gateway-darwin-arm64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.2.0/gateway-darwin-arm64) / [.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.2.0/gateway-darwin-arm64.tar.gz) |
| macOS Intel | [gateway-darwin-amd64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.2.0/gateway-darwin-amd64) / [.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.2.0/gateway-darwin-amd64.tar.gz) |
| Linux x86_64 | [gateway-linux-amd64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.2.0/gateway-linux-amd64) / [.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.2.0/gateway-linux-amd64.tar.gz) |
| Linux ARM64 | [gateway-linux-arm64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.2.0/gateway-linux-arm64) / [.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.2.0/gateway-linux-arm64.tar.gz) |
| Windows x86_64 | [gateway-windows-amd64.exe](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.2.0/gateway-windows-amd64.exe) / [.zip](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.2.0/gateway-windows-amd64.zip) |

校验文件: [SHA256SUMS](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.2.0/SHA256SUMS)

## 升级

**Windows 重点**（旧版清理残留是 no-op，升级前建议手动清干净再覆盖）：

```powershell
gateway stop
taskkill /IM mihomo.exe /T /F 2>$null
# 重新跑安装脚本（会自动拉 v3.2.0）
irm https://raw.githubusercontent.com/Tght1211/lan-proxy-gateway/main/install.ps1 | iex
gateway install
```

**macOS / Linux**：

```bash
sudo gateway stop
curl -fsSL https://raw.githubusercontent.com/Tght1211/lan-proxy-gateway/main/install.sh | sudo bash
sudo gateway install
```
