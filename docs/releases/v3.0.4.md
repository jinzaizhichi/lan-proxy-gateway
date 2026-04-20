# LAN Proxy Gateway v3.0.4

单点修复版：v3.0.3 里 `sudo gateway install` 最后一步「装开机自启」在 macOS 上概率性报 `Bootstrap failed: 5: Input/output error`。

## 修复

### macOS launchd bootstrap 加鲁棒性

mac `launchctl bootstrap` 如果发现目标 label（`com.lan-proxy-gateway`）已经存在「半注册」状态（上一次 bootstrap 没完全走完、或者之前用旧版装过没清干净），会直接返回 exit 5 / Input/output error。常见触发条件：

- 老版本装过 service（plist 里还是 `start` 前台），升级到 v3.0.3 后 plist 覆盖了但 label 注册状态没更新
- 之前手动跑过 `launchctl load` / `bootstrap` 中断过

v3.0.4 在 `InstallService` 里做了三件事：

1. bootstrap 前先 `bootout` 清可能残留的 label 注册（没注册时忽略返回值）
2. bootstrap 成功后追加 `launchctl enable system/com.lan-proxy-gateway`，防止被加进 disabled 列表的情况
3. 把「already bootstrapped」/「service already bootstrapped」识别为成功，而不是报错

Linux 和 Windows 路径不受影响。

## 手动兜底（升级前就用 v3.0.3 + 这次装 service 失败了的用户）

不一定要升级到 v3.0.4 才能装上 service，直接手动跑这三条也可以恢复：

```bash
sudo launchctl bootout system/com.lan-proxy-gateway 2>/dev/null || true
sudo launchctl bootstrap system /Library/LaunchDaemons/com.lan-proxy-gateway.plist
sudo launchctl enable system/com.lan-proxy-gateway
```

## 下载

| 系统 | 直接下载 |
|---|---|
| macOS Apple Silicon | [gateway-darwin-arm64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.4/gateway-darwin-arm64) / [gateway-darwin-arm64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.4/gateway-darwin-arm64.tar.gz) |
| macOS Intel | [gateway-darwin-amd64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.4/gateway-darwin-amd64) / [gateway-darwin-amd64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.4/gateway-darwin-amd64.tar.gz) |
| Linux x86_64 | [gateway-linux-amd64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.4/gateway-linux-amd64) / [gateway-linux-amd64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.4/gateway-linux-amd64.tar.gz) |
| Linux ARM64 | [gateway-linux-arm64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.4/gateway-linux-arm64) / [gateway-linux-arm64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.4/gateway-linux-arm64.tar.gz) |
| Windows x86_64 | [gateway-windows-amd64.exe](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.4/gateway-windows-amd64.exe) / [gateway-windows-amd64.zip](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.4/gateway-windows-amd64.zip) |

校验文件: [SHA256SUMS](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.0.4/SHA256SUMS)
