# LAN Proxy Gateway v2.2.7

这一版把 Windows 兼容性做了一次完整补验，不只是修文案，而是把安装、更新、状态检测、开机自启和提示链路一起收了一轮。

## 重点修复

### 1. 终端与权限判断更稳

- 启动时会主动切到 UTF-8 代码页 `65001`
- 同时开启 ANSI 虚拟终端支持
- 管理员权限判断不再依赖脆弱的 `net session`，改成 Windows Token 成员检测

### 2. 安装和更新链路补齐

- `install.ps1` 安装完成后会立刻刷新当前 PowerShell 的 `$env:Path`
- Windows 下的 `gateway update` 改成后台替换锁定中的 `.exe`
- 当前进程退出后会自动完成替换并重新启动网关

### 3. 中文 Windows 下的状态检测更可靠

- `IsRunning()` 不再依赖 `"没有"` / `"No tasks"` 这类本地化文本
- `IsIPForwardingEnabled()` 不再解析 `netsh` 的中英文输出
- 默认网卡与默认网关改成按默认路由表识别，不再容易误选 VPN / 虚拟网卡

### 4. Windows 的开机自启和命令提示更合理

- `gateway service install` 在 Windows 下改成安装任务计划程序启动任务
- 不再把普通 CLI 伪装成 `sc.exe` 服务
- 日志和运行提示也按平台输出：
  - Windows：`gateway <subcommand>`、`Get-Content -Wait`
  - macOS / Linux：`sudo gateway <subcommand>`、`tail -f`

### 5. 路径细节补齐

- 日志路径统一走 `os.TempDir()`
- 本地配置路径现在同时支持 `~/...` 和 `~\...`

## 下载

| 系统 | 直接下载 |
|---|---|
| macOS Apple Silicon | [gateway-darwin-arm64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.7/gateway-darwin-arm64) / [gateway-darwin-arm64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.7/gateway-darwin-arm64.tar.gz) |
| macOS Intel | [gateway-darwin-amd64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.7/gateway-darwin-amd64) / [gateway-darwin-amd64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.7/gateway-darwin-amd64.tar.gz) |
| Linux x86_64 | [gateway-linux-amd64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.7/gateway-linux-amd64) / [gateway-linux-amd64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.7/gateway-linux-amd64.tar.gz) |
| Linux ARM64 | [gateway-linux-arm64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.7/gateway-linux-arm64) / [gateway-linux-arm64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.7/gateway-linux-arm64.tar.gz) |
| Windows x86_64 | [gateway-windows-amd64.exe](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.7/gateway-windows-amd64.exe) / [gateway-windows-amd64.zip](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.7/gateway-windows-amd64.zip) |

校验文件: [SHA256SUMS](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.7/SHA256SUMS)

## 建议验证

升级到 `v2.2.7` 后，建议顺手确认这几处：

1. 在 Windows PowerShell 里安装完成后，直接输入 `gateway`，确认当前窗口无需重开就能识别命令
2. 执行 `gateway update`，确认更新流程会在后台替换并重新启动，而不是卡在当前 `.exe`
3. 执行 `gateway status`，确认默认网卡、IP 转发和进程状态都能正常显示
4. 执行 `gateway service install`，确认开机自启任务能被正常创建
