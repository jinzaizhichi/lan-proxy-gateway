# LAN Proxy Gateway v2.2.0

这是一次围绕“让小白也能玩转局域网共享和链式代理”的系统化升级。

## 这一版最重要的变化

### 1. 启动后不再只是打印结果，而是进入运行中控制台

- `gateway start` 成功后进入全屏 TUI 控制台
- 支持 slash 命令
- 支持确认流
- 支持 `Ctrl+P` 切换策略组和节点

### 2. 配置结构完成统一

配置现在围绕:

- `proxy`
- `runtime`
- `rules`
- `extension`

组织，同时保留对旧版顶层字段的兼容读取。

### 3. 链式代理和出口展示更清楚

`gateway status` 现在会明确展示:

- 入口节点
- 普通出口
- 住宅出口

更适合验证 `chains` 是否真的生效。

### 4. 本机绕过代理

新增:

```yaml
runtime:
  tun:
    bypass_local: true
```

适合只把当前机器当“局域网网关”，不让它自己也被代理接管。

### 5. AI Skill 和普通权限控制

新增:

- `gateway skill`
- `gateway permission`

更适合和 AI 客户端协同，也更适合减少每次手动 `sudo` 的心智负担。

### 6. 更完整的更新体验

- `gateway update`
- 首页 / 启动页 / TUI 新版提醒
- release 下载继续带镜像回退

## 下载哪个文件？

| 系统 | 推荐文件 |
|---|---|
| macOS Apple Silicon | `gateway-darwin-arm64` 或 `gateway-darwin-arm64.tar.gz` |
| macOS Intel | `gateway-darwin-amd64` 或 `gateway-darwin-amd64.tar.gz` |
| Linux x86_64 | `gateway-linux-amd64` 或 `gateway-linux-amd64.tar.gz` |
| Linux ARM64 | `gateway-linux-arm64` 或 `gateway-linux-arm64.tar.gz` |
| Windows x86_64 | `gateway-windows-amd64.exe` 或 `gateway-windows-amd64.zip` |

另外附带:

- `SHA256SUMS`

## 升级建议

从 `v2.1.1` 升级建议:

1. 先备份当前 `gateway.yaml`
2. 使用 `sudo gateway update`
3. 执行 `gateway config show`
4. 执行 `gateway status`

旧配置仍兼容读取，但建议后续逐步迁移到新的分组结构。
