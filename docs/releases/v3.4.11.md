# LAN Proxy Gateway v3.4.11

**安装脚本下载校验修复版**。这一版修复 `curl | bash` 安装时镜像源返回错误内容却被当成 `.tar.gz` 解压的问题。

## 修复

- **`install.sh` 不再把镜像错误页当 release 包解压**

  部分 GitHub 加速源（这次是 `ghps.cc`）会对 release asset 返回 HTML / 错误页，但 HTTP 状态码仍让 `curl -fSL` 看起来成功。旧脚本下载后直接 `tar -xzf`，因此报：

  ```text
  tar: Error opening archive: Unrecognized archive format
  ```

  现在下载完成后会按用途校验内容：

  - GitHub API 响应必须包含 `tag_name`
  - `.tar.gz` release asset 必须同时通过 `gzip -t` 和 `tar -tzf`

  某个候选源下载成功但内容无效时，脚本会打印“下载失败或内容无效”，然后自动尝试下一个源，而不是进入解压。

## 测试

- `bash -n install.sh`
- 手动验证非 gzip 内容会被拒绝
- 手动验证本地 release tarball 可通过 `gzip -t` 和 `tar -tzf`
- `go test ./...`
- `make build-all VERSION=v3.4.11`

## 下载

| 系统 | 直接下载 |
|---|---|
| macOS Apple Silicon | [gateway-darwin-arm64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.4.11/gateway-darwin-arm64) / [.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.4.11/gateway-darwin-arm64.tar.gz) |
| macOS Intel | [gateway-darwin-amd64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.4.11/gateway-darwin-amd64) / [.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.4.11/gateway-darwin-amd64.tar.gz) |
| Linux x86_64 | [gateway-linux-amd64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.4.11/gateway-linux-amd64) / [.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.4.11/gateway-linux-amd64.tar.gz) |
| Linux ARM64 | [gateway-linux-arm64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.4.11/gateway-linux-arm64) / [.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.4.11/gateway-linux-arm64.tar.gz) |
| Windows x86_64 | [gateway-windows-amd64.exe](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.4.11/gateway-windows-amd64.exe) / [.zip](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.4.11/gateway-windows-amd64.zip) |

校验文件: [SHA256SUMS](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.4.11/SHA256SUMS)

## 升级 / 安装

安装脚本本身从 `main` 拉取，更新 main 后立即生效：

```bash
curl -fsSL https://raw.githubusercontent.com/Tght1211/lan-proxy-gateway/main/install.sh | bash
```

已安装用户也可以升级到本版：

```bash
gateway update v3.4.11
```
