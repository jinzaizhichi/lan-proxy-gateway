# LAN Proxy Gateway v3.1.2

文档版：README 全面精致化 + 新增 AI 运维 Skill。代码零功能改动，跨平台（darwin / linux / windows）全部验证通过。

## README 全面重写

- 首屏 `✨ 能干什么` 罗列 10+ 项核心能力（TUN 网关、Web 控制台、自动自愈、链式代理预设、节点测速、日志易读视图…）
- 三种设备接入方式独立成节，表格化展示「填什么 / 适合谁」
- 新增 `🔧 典型场景` 四个示例（已有 Clash Verge / 纯订阅 / AI 住宅 IP / 本机 Mac 浏览器）
- Web 控制台 / 主菜单 / 命令行 / 跨平台支持各自独立成章
- 三层架构图加第四块「健壮性 · supervisor」
- 目录结构同步到 v3.1（新增 `script/presets` / `embed/webui`）

## AI 运维 Skill

新增 `~/.claude/skills/lan-proxy-gateway-ops/SKILL.md` 让 Claude 等 AI 代理能直接运维本机 gateway：

- **无需 sudo 路径（推荐 AI 用）**：走 mihomo REST API（`http://127.0.0.1:19090`）做日常操作 — 切节点 / 换模式 / 测延迟 / 查连接 / 读规则
- **需要 sudo 的命令**：配了 sudoers NOPASSWD 白名单后 AI 可免密跑；否则走系统 service 一次装好让 mihomo 常驻
- 典型场景推荐做法（换节点 / 排查 / 加规则 / 本机 DNS 修复）全部列出
- 注意事项（URL 编码 CJK / Emoji 分组名、reload 范围、supervisor 自动 fallback 的状态理解）

## 跨平台验证

本版二进制三平台全部编译通过：

- `darwin-arm64` / `darwin-amd64`（macOS Apple Silicon + Intel）
- `linux-amd64` / `linux-arm64`
- `windows-amd64`

**Linux/Windows 兼容性现状**：
- ✅ 启动 / 停止 / TUN / 代理 / Web 控制台 / script 预设 全部可用
- ✅ 服务管理：systemd unit / Windows schtasks
- ⚠️ 本机 DNS 一键切（方式 3）只有 macOS 实装（`networksetup`）；Linux（systemd-resolved / NetworkManager）和 Windows（netsh）菜单里显示手动命令，不阻塞使用

## 下载

| 系统 | 直接下载 |
|---|---|
| macOS Apple Silicon | [gateway-darwin-arm64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.1.2/gateway-darwin-arm64) / [.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.1.2/gateway-darwin-arm64.tar.gz) |
| macOS Intel | [gateway-darwin-amd64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.1.2/gateway-darwin-amd64) / [.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.1.2/gateway-darwin-amd64.tar.gz) |
| Linux x86_64 | [gateway-linux-amd64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.1.2/gateway-linux-amd64) / [.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.1.2/gateway-linux-amd64.tar.gz) |
| Linux ARM64 | [gateway-linux-arm64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.1.2/gateway-linux-arm64) / [.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.1.2/gateway-linux-arm64.tar.gz) |
| Windows x86_64 | [gateway-windows-amd64.exe](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.1.2/gateway-windows-amd64.exe) / [.zip](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.1.2/gateway-windows-amd64.zip) |

校验文件: [SHA256SUMS](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.1.2/SHA256SUMS)

## 升级

二进制：`curl | bash` 一键装脚本会自动装上 v3.1.2。

配置 / 运行时无影响，无须重启 mihomo。
