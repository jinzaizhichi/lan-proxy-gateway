# LAN Proxy Gateway v3.0.2

单行修复版：v3.0.1 里 `gateway -v` 错误输出 `dev` 而不是真实版本号。功能层面和 v3.0.1 完全等价，升级只是为了让版本号能正确显示。

## 重点变更

### 修复 `gateway -v` 总是显示 `dev`

- 根因：`ldflags -X main.version=...` 注入的是 main 包变量，然后 main.go 在运行时才赋值给 `cmd.Version`；但 `cmd.rootCmd.Version` 已经在 cmd 包初始化时按值拷贝了默认的 `"dev"`，main.go 里的赋值来得太晚
- 修复：ldflags 直接注入 `cmd.Version`，main.go 去中转，版本号在 cobra 读 rootCmd.Version 之前就已经是实际 tag 了
- 现在 `gateway -v` 应正确输出 `gateway version v3.0.2`

v3.0.1 相对 v3.0.0 的功能变更（Q 只退控制台 / 新增关闭 gateway 入口 / reattach 后台 mihomo / install 进主菜单）全部保留，不赘述。

## 下载

| 系统 | 直接下载 |
|---|---|
| macOS Apple Silicon | [gateway-darwin-arm64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.2/gateway-darwin-arm64) / [gateway-darwin-arm64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.2/gateway-darwin-arm64.tar.gz) |
| macOS Intel | [gateway-darwin-amd64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.2/gateway-darwin-amd64) / [gateway-darwin-amd64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.2/gateway-darwin-amd64.tar.gz) |
| Linux x86_64 | [gateway-linux-amd64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.2/gateway-linux-amd64) / [gateway-linux-amd64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.2/gateway-linux-amd64.tar.gz) |
| Linux ARM64 | [gateway-linux-arm64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.2/gateway-linux-arm64) / [gateway-linux-arm64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.2/gateway-linux-arm64.tar.gz) |
| Windows x86_64 | [gateway-windows-amd64.exe](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.2/gateway-windows-amd64.exe) / [gateway-windows-amd64.zip](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.2/gateway-windows-amd64.zip) |

校验文件: [SHA256SUMS](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.2/SHA256SUMS)

## 建议验证

1. `gateway -v` 应显示 `gateway version v3.0.2`
2. `sudo gateway` 主菜单有第 6 项「关闭 gateway 并退出（停 mihomo）」
3. `Q` 退控制台后，重新进来应直接显示「运行中」而不是「未启动」
