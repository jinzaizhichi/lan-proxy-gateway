# LAN Proxy Gateway v3.3.2

**IPv6 回归修复**。v2.x 时内嵌的 mihomo 模板是默认允许 IPv6 的；v3.0 重写模板时从 v1 legacy 继承了 `ipv6: false` + `dns.ipv6: false`，导致 3.x 上双栈目标（Cloudflare / ChatGPT / Google 等）在 IPv6 优先的客户端上解析不到 AAAA、走 IPv6 偏好的 App 直接超时。本版恢复 v2 行为。

## 🐛 修复

- **恢复全局 IPv6 支持** ([embed/template.yaml](embed/template.yaml))：两处强制禁用改回开启
  - `ipv6: false` → `ipv6: true`（根层，允许 mihomo 处理 IPv6 目标地址）
  - `dns.ipv6: false` → `dns.ipv6: true`（DNS 返回 AAAA 记录）
  - 有用户反馈「2.x 版本之前支持 IPv6，更新到 3.x 后不支持了」就是这个回归。

## 🛡 未改动 / 设计抉择

- **TUN 块保持不劫持 IPv6 路由**（`renderTUNBlock` 不写 `inet6-address`）。打开 IPv6 后如果同时让 TUN 接管 IPv6 默认路由，会重演 [issue #2](https://github.com/Tght1211/lan-proxy-gateway/issues/2) 的 Linux 场景：本机 IPv6 对外可达性被切断，LAN 里用 IPv6 ping 不到这台网关。当前选择：IPv6 流量继续走系统默认路由，本机 IPv6 服务保持外部可达，mihomo 只接管 IPv4 TUN 劫持路径（绝大多数用户的 LAN / 家宽在 IPv4 上），v2.x 行为 ≈ 对齐。

  **如果你确实想让 TUN 也接管 IPv6**（某些纯 IPv6 环境可能需要）：等下一个小版本加 `gateway.tun.ipv6` 开关，或直接自己改 `~/.config/lan-proxy-gateway/mihomo/config.yaml` 里的 `tun:` 块加 `inet6-address: fdfe:dcba:9876::1/126`。

## ✅ 测试

- 新增 `TestRenderIPv6Enabled`：render 输出必须同时包含 `ipv6: true`（全局）和 DNS 块的 `ipv6: true`，且不再残留任何 `ipv6: false`
- 全仓 `go test ./...` + `-race` 全过

## 下载

| 系统 | 直接下载 |
|---|---|
| macOS Apple Silicon | [gateway-darwin-arm64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.3.2/gateway-darwin-arm64) / [.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.3.2/gateway-darwin-arm64.tar.gz) |
| macOS Intel | [gateway-darwin-amd64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.3.2/gateway-darwin-amd64) / [.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.3.2/gateway-darwin-amd64.tar.gz) |
| Linux x86_64 | [gateway-linux-amd64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.3.2/gateway-linux-amd64) / [.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.3.2/gateway-linux-amd64.tar.gz) |
| Linux ARM64 | [gateway-linux-arm64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.3.2/gateway-linux-arm64) / [.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.3.2/gateway-linux-arm64.tar.gz) |
| Windows x86_64 | [gateway-windows-amd64.exe](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.3.2/gateway-windows-amd64.exe) / [.zip](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.3.2/gateway-windows-amd64.zip) |

校验文件: [SHA256SUMS](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.3.2/SHA256SUMS)

## 升级

**macOS / Linux**：

```bash
sudo gateway stop
curl -fsSL https://raw.githubusercontent.com/Tght1211/lan-proxy-gateway/main/install.sh | sudo bash
sudo gateway start
```

**Windows**：

```powershell
gateway stop
irm https://raw.githubusercontent.com/Tght1211/lan-proxy-gateway/main/install.ps1 | iex
gateway start
```

旧用户**必须重启 mihomo**（`sudo gateway stop` 然后 `sudo gateway start`，或菜单 `[M] → 4 → 2`）才能让新模板生效，升级完二进制但不重启的话旧 config.yaml 还在 disk 上。
