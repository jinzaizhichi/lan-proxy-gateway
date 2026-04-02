# LAN Proxy Gateway v2.2.0

这次重点是把局域网共享、chains 和终端交互整理成一套更完整的系统。

## 这一版最重要的变化

### 1. 启动后进入运行中控制台

- `gateway start` 成功后进入全屏 TUI 控制台
- 支持 slash 命令
- 支持确认流
- 支持 `Ctrl+P` 切换策略组和节点

### 2. 配置结构统一

配置现在围绕:

- `proxy`
- `runtime`
- `rules`
- `extension`

同时保留旧版顶层字段兼容读取。

### 3. chains 和出口展示更清楚

`gateway status` 现在会明确展示:

- 入口节点
- 普通出口
- 住宅出口

更方便确认 `chains` 是否真的生效。

### 4. 支持本机绕过代理

新增:

```yaml
runtime:
  tun:
    bypass_local: true
```

适合只让局域网设备走网关。

### 5. AI Skill 和权限控制

新增:

- `gateway skill`
- `gateway permission`

更方便 AI 客户端协同，也减少频繁手动 `sudo`。

### 6. 更新体验补齐

- `gateway update`
- 首页 / 启动页 / TUI 新版提醒
- release 下载继续带镜像回退

## 下载

| 系统 | 直接下载 |
|---|---|
| macOS Apple Silicon | [gateway-darwin-arm64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.0/gateway-darwin-arm64) / [gateway-darwin-arm64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.0/gateway-darwin-arm64.tar.gz) |
| macOS Intel | [gateway-darwin-amd64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.0/gateway-darwin-amd64) / [gateway-darwin-amd64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.0/gateway-darwin-amd64.tar.gz) |
| Linux x86_64 | [gateway-linux-amd64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.0/gateway-linux-amd64) / [gateway-linux-amd64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.0/gateway-linux-amd64.tar.gz) |
| Linux ARM64 | [gateway-linux-arm64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.0/gateway-linux-arm64) / [gateway-linux-arm64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.0/gateway-linux-arm64.tar.gz) |
| Windows x86_64 | [gateway-windows-amd64.exe](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.0/gateway-windows-amd64.exe) / [gateway-windows-amd64.zip](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.0/gateway-windows-amd64.zip) |

校验文件: [SHA256SUMS](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.0/SHA256SUMS)

## 升级建议

从 `v2.1.1` 升级:

1. 先备份当前 `gateway.yaml`
2. 使用 `sudo gateway update`
3. 执行 `gateway config show`
4. 执行 `gateway status`

旧配置仍兼容读取。
