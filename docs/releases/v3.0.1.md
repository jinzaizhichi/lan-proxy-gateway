# LAN Proxy Gateway v3.0.1

这一版主要修掉 v3.0.0 控制台在退出 / 重进流程上的体验坑，让「退出操作」和「退出网关」两件事分清楚。

## 重点变更

### 1. `Q` 只退控制台，不再停 mihomo

- 以前在主菜单按 `Q` 会把后台 mihomo 一起关掉，主人抱怨「我只是退出操作，又不是退出程序」
- 现在 `Q` 的语义统一为「返回 / 退出控制台」，mihomo 继续在后台跑
- 想要停 mihomo 请走 `4 启动 / 重启 / 停止` 里的「停止」

### 2. 主菜单新增 `6 关闭 gateway 并退出`

- 一站式入口：停 mihomo + 退出控制台
- 有 mihomo 在跑时会二次确认，避免误操作导致 LAN 设备突然断网

### 3. 新进程自动接管后台 mihomo

- 之前 `sudo gateway` 再进来时，主菜单会错显「未启动」，按启动又被 preflight 当成端口冲突、把自己的 mihomo 杀了重启
- 现在 `engine.New()` 会预先读 `mihomo.pid`，如果后台有活着的 mihomo 就直接接管
- `Engine.Start()` 检测到已在运行会短路返回，不会再误杀
- 按「启动」时如果已在跑，会直接提示「已经在运行，无需再启动」

### 4. `gateway install` 装完直接进主菜单

- 以前装完会 `waitForSignal()` 卡在前台等 `Ctrl+C`
- 现在和默认 `gateway` 一致，装完直接进主菜单，用户可以就地查状态、换源、装 service

## 下载

| 系统 | 直接下载 |
|---|---|
| macOS Apple Silicon | [gateway-darwin-arm64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.1/gateway-darwin-arm64) / [gateway-darwin-arm64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.1/gateway-darwin-arm64.tar.gz) |
| macOS Intel | [gateway-darwin-amd64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.1/gateway-darwin-amd64) / [gateway-darwin-amd64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.1/gateway-darwin-amd64.tar.gz) |
| Linux x86_64 | [gateway-linux-amd64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.1/gateway-linux-amd64) / [gateway-linux-amd64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.1/gateway-linux-amd64.tar.gz) |
| Linux ARM64 | [gateway-linux-arm64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.1/gateway-linux-arm64) / [gateway-linux-arm64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.1/gateway-linux-arm64.tar.gz) |
| Windows x86_64 | [gateway-windows-amd64.exe](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.1/gateway-windows-amd64.exe) / [gateway-windows-amd64.zip](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.1/gateway-windows-amd64.zip) |

校验文件: [SHA256SUMS](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.1/SHA256SUMS)

## 建议验证

1. `sudo gateway` 进主菜单，按 `4 → 1` 启动网关，确认状态变「运行中」
2. 按 `Q` 退控制台；重新 `sudo gateway` 进来，主菜单应直接显示「运行中」而不是「未启动」
3. 再按一次 `4 → 1`，应提示「已经在运行，无需再启动」（不会杀掉自己的 mihomo）
4. 按 `6` 关闭 gateway，确认有二次确认，选 y 后 mihomo 正常停掉、程序退出
