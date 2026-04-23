# LAN Proxy Gateway v3.3.1

紧跟 v3.3.0 的 **首页修正小版本**。两件事：单跳代理时首页不再重复显示「起飞 == 落地」；新增 LAN 代理端口健康探测，mihomo 中途挂掉或端口被抢时首页会立刻红字警告，不再装作一切正常。

## 🐛 修复

- **首页单跳代理展示**：没配链式代理（住宅 IP 落地）时，🛫 起飞 和 🛬 落地显示的是同一个节点，用户会以为是 bug。v3.3.1 改为：
  - 单跳（`takeoff == landing`）→ 一行 **`🌐 出口节点  xxx`**，清爽明确
  - 链式代理（有 ChainResidential）→ 保留原来的 **`🛫 起飞 + 🛬 落地`** 两行
  - 「真实出口」（ipinfo.io 查到的 IP / 位置 / ISP）两种场景都在下面一行显示

- **LAN 代理端口健康探测**：新增 `mixed` 端口活性探测，`API 通 + mixed 不通`时首页红字告警：

  ```
  ⚠ 代理端口不通：LAN 设备（手机 / Switch 等）现在连不上这台网关
    常见诱因：reload 后 mihomo 没起干净 · TUN strict-route 抢端口 · 防火墙拦 LAN 入站
    建议：[M] → 4 → 2 重启；还不行就 [M] → 4 → 4 清理残留 mihomo 后再 1 启动
  ```

  之前若 mihomo 启动阶段通过 API 就绪检测（API 已 listen），随后因其他原因挂掉或 mixed 端口没成功绑上，dashboard 仍会继续显示"● 运行中"的过期状态，直到下次手动操作才发现手机连不上。v3.3.1 把 LAN 入口端口也拉进活性检查，问题当场可见。

## ✅ 测试

- 新增 `TestDrawDashboardHopCollapsesWhenNoChain` / `TestDrawDashboardHopShowsBothWhenChain`：两种模式的渲染分支
- 新增 `TestDrawDashboardWarnsOnMixedPortDown`：端口告警文案
- macOS 本机端到端验证 LAN 代理链路：本机 + 局域网 IP 走 17890 HTTP 代理访问 gstatic / baidu / claude.ai 全通
- 全仓 `go test ./...` + `-race` 全过

## 下载

| 系统 | 直接下载 |
|---|---|
| macOS Apple Silicon | [gateway-darwin-arm64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.3.1/gateway-darwin-arm64) / [.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.3.1/gateway-darwin-arm64.tar.gz) |
| macOS Intel | [gateway-darwin-amd64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.3.1/gateway-darwin-amd64) / [.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.3.1/gateway-darwin-amd64.tar.gz) |
| Linux x86_64 | [gateway-linux-amd64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.3.1/gateway-linux-amd64) / [.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.3.1/gateway-linux-amd64.tar.gz) |
| Linux ARM64 | [gateway-linux-arm64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.3.1/gateway-linux-arm64) / [.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.3.1/gateway-linux-arm64.tar.gz) |
| Windows x86_64 | [gateway-windows-amd64.exe](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.3.1/gateway-windows-amd64.exe) / [.zip](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.3.1/gateway-windows-amd64.zip) |

校验文件: [SHA256SUMS](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.3.1/SHA256SUMS)

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

完整 v3.3.0 功能列表：[v3.3.0 release notes](./v3.3.0.md)
