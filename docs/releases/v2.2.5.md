# LAN Proxy Gateway v2.2.5

这一版主要收口 Windows 兼容性，重点是把 `gateway install` 和 `gateway update` 在 Windows x86_64 上的内核与升级链路补齐。

## 这一版最重要的变化

### 1. 修复 Windows 的 mihomo 内核解压

- 官方 `mihomo` 的 Windows zip 里，实际文件名是 `mihomo-windows-amd64-compatible.exe`
- 之前的逻辑只认 `mihomo.exe`
- 现在会正确识别官方包里的可执行文件，并统一安装为本地 `mihomo.exe`

### 2. 修复 Windows 的 gateway update

- 之前 Windows 升级时请求的是不带后缀的 `gateway-windows-amd64`
- GitHub Release 实际资产是 `gateway-windows-amd64.exe`
- 现在 `gateway update` 会下载正确的 Windows 资产

### 3. 补了一轮更实的 Windows 验证

- 核对了 MetaCubeX 官方 `mihomo v1.19.8` 的 Windows release 资产名
- 验证了官方 Windows 下载 URL 可访问
- 验证了官方 zip 的真实内容结构
- 补上了对应的自动化测试

## 下载

| 系统 | 直接下载 |
|---|---|
| macOS Apple Silicon | [gateway-darwin-arm64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.5/gateway-darwin-arm64) / [gateway-darwin-arm64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.5/gateway-darwin-arm64.tar.gz) |
| macOS Intel | [gateway-darwin-amd64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.5/gateway-darwin-amd64) / [gateway-darwin-amd64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.5/gateway-darwin-amd64.tar.gz) |
| Linux x86_64 | [gateway-linux-amd64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.5/gateway-linux-amd64) / [gateway-linux-amd64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.5/gateway-linux-amd64.tar.gz) |
| Linux ARM64 | [gateway-linux-arm64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.5/gateway-linux-arm64) / [gateway-linux-arm64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.5/gateway-linux-arm64.tar.gz) |
| Windows x86_64 | [gateway-windows-amd64.exe](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.5/gateway-windows-amd64.exe) / [gateway-windows-amd64.zip](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.5/gateway-windows-amd64.zip) |

校验文件: [SHA256SUMS](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.5/SHA256SUMS)

## 升级建议

Windows 用户升级到 `v2.2.5` 后，建议顺手验证这两条链路:

1. 运行 `gateway install`
2. 确认能自动下载并安装 `mihomo`
3. 再运行 `gateway update`
4. 确认升级链路能拿到新的 `.exe` 版本
