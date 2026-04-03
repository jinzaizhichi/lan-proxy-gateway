# LAN Proxy Gateway v2.2.2

这一版继续收口运行中控制台，重点修正 TUI 顶部焦点和功能导航页面。

## 这一版最重要的变化

### 1. `功能导航` 改成了 TUI 原生页面

- 不再抓取 CLI 输出再塞回 TUI
- 不会再出现分隔线乱入、标题重复或布局被打散的问题
- 导航内容现在会按当前模式给出更清楚的下一步入口

### 2. 顶部 tab 现在是真正可聚焦的区域

- 在左侧功能列表里按 `Esc`，焦点会回到顶部 tab
- 顶部区域会有明确高亮
- `←/→` 切换分区
- `↓` 或 `Enter` 回到左侧功能列表

### 3. 控制台操作提示更清楚

- README 中英版补上了运行中控制台的常用按键
- 命令文档补上了顶部 tab 与左侧列表的焦点流转说明

## 下载

| 系统 | 直接下载 |
|---|---|
| macOS Apple Silicon | [gateway-darwin-arm64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.2/gateway-darwin-arm64) / [gateway-darwin-arm64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.2/gateway-darwin-arm64.tar.gz) |
| macOS Intel | [gateway-darwin-amd64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.2/gateway-darwin-amd64) / [gateway-darwin-amd64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.2/gateway-darwin-amd64.tar.gz) |
| Linux x86_64 | [gateway-linux-amd64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.2/gateway-linux-amd64) / [gateway-linux-amd64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.2/gateway-linux-amd64.tar.gz) |
| Linux ARM64 | [gateway-linux-arm64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.2/gateway-linux-arm64) / [gateway-linux-arm64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.2/gateway-linux-arm64.tar.gz) |
| Windows x86_64 | [gateway-windows-amd64.exe](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.2/gateway-windows-amd64.exe) / [gateway-windows-amd64.zip](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.2/gateway-windows-amd64.zip) |

校验文件: [SHA256SUMS](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.2/SHA256SUMS)

## 升级建议

从 `v2.2.1` 升级:

1. 使用 `sudo gateway update`
2. 执行 `sudo gateway console`
3. 在 TUI 里验证 `功能导航` 和顶部 tab 焦点切换
4. 如需纯命令模式，执行 `sudo gateway console --simple`
