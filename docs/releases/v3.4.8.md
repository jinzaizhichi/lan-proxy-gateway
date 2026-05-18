# LAN Proxy Gateway v3.4.8

**Gateway Web 控制台体验增强版**。这一版把首页改成真正可操作的控制台：代理源、策略组、节点测速、接入设备、TUN 和局域网代理服务都可以在 WebUI 里直接管理。

## 新增 / 改进

- **新增 Gateway Web 控制台入口**

  `gateway start` 现在会同时启动 Gateway 自带 WebUI（默认 `http://<gatewayIP>:19091/`），启动横幅会打印带 `#token=...` 的完整访问地址。也可以用 `gateway webui` 随时重新查看当前地址。

- **首页代理源可视化与快捷切换**

  首页显示当前代理源类型、路径/订阅详情、健康状态和错误摘要；有多个有效 source profile 时可直接在首页切换，没有其它源时提供到配置库的入口。

- **Clash 风格策略组控制**

  首页不再展示单个“生效分组”。WebUI 会从 Mihomo `/proxies` 读取所有 Selector 策略组，提供可滚动分组列表、当前分组节点、快捷节点切换和一键测速；完整分组/节点列表仍可跳到“路由矩阵”。

- **透明代理与代理服务不再互斥**

  TUN 透明代理和 HTTP/SOCKS5 `mixed-port` 现在是两个独立能力。首页“接入能力”可分别进入设置，TUN 可独立开关，代理服务可独立开关。

- **局域网代理服务支持认证**

  新增 `runtime.proxy_service` 配置：

  ```yaml
  runtime:
    proxy_service:
      enabled: true
      username: ""
      password: ""
  ```

  用户名/密码留空时保持无认证；填写后会渲染 Mihomo `authentication`。关闭代理服务时会渲染 `mixed-port: 0`。

- **实时看板改为接入设备数**

  首页不再把每条网络请求算成“活跃连接”，而是按唯一来源 IP 统计接入设备数；点击看板仍可查看设备和连接明细。

- **内置规则集可完整查看和复制**

  WebUI 规则集行支持展开完整规则列表、滚动查看和一键复制，便于确认内置规则到底匹配什么。

- **左上角版本与手动升级入口**

  版本号现在展示在左上角，点击可检查 GitHub 最新 release；有新版本时可从 WebUI 触发升级流程。

## 修复

- 修复 WebUI 关闭透明代理后状态仍显示开启的问题。本机上游代理（例如 `127.0.0.1:7890`）现在只会在 TUN 本来开启时设置 `bypass_local` 防自循环，不再覆盖用户显式关闭 TUN 的选择。
- 修正首页能力文案：不再用旧 `gateway.mode` 显示“端口模式 · 零干扰”，而是按实际能力显示“透明代理已启用 / 代理服务已启用”。
- 修复快速策略组区域在分组很多时不适合横向固定卡片的问题，改为可扩展的分组列表 + 当前分组节点视图。

## 测试

- `node --check internal/webui/static/app.js`
- `go test ./internal/config ./internal/engine ./internal/app ./internal/webui`
- `go test ./...`

新增 / 更新测试覆盖：

- 本机 loopback external 源不会强制把关闭的 TUN/DNS 改回开启
- TUN 开启时本机 loopback external 源仍会强制 `bypass_local`，避免自代理回环
- 代理服务认证与关闭状态的 Mihomo 渲染
- WebUI API、token 鉴权、静态资源和控制路由

## 下载

| 系统 | 直接下载 |
|---|---|
| macOS Apple Silicon | [gateway-darwin-arm64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.4.8/gateway-darwin-arm64) / [.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.4.8/gateway-darwin-arm64.tar.gz) |
| macOS Intel | [gateway-darwin-amd64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.4.8/gateway-darwin-amd64) / [.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.4.8/gateway-darwin-amd64.tar.gz) |
| Linux x86_64 | [gateway-linux-amd64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.4.8/gateway-linux-amd64) / [.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.4.8/gateway-linux-amd64.tar.gz) |
| Linux ARM64 | [gateway-linux-arm64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.4.8/gateway-linux-arm64) / [.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.4.8/gateway-linux-arm64.tar.gz) |
| Windows x86_64 | [gateway-windows-amd64.exe](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.4.8/gateway-windows-amd64.exe) / [.zip](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.4.8/gateway-windows-amd64.zip) |

校验文件: [SHA256SUMS](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.4.8/SHA256SUMS)

## 升级

```bash
sudo gateway update
```

或者指定版本：

```bash
sudo gateway update v3.4.8
```
