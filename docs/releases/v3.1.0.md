# LAN Proxy Gateway v3.1.0

一次攒了很多改动的版本，围绕两条主线：**小白友好**（菜单、文案、引导）+ **功能完整**（Web 控制台内嵌、订阅 inline、链式代理预设、节点测速、代理源自愈）。

## 1. 主菜单重新命名 + 白话描述

```
  1  设备接入指引      Switch / PS5 / 手机怎么连到这里
  2  分流 & 规则        国内直连 / 国外走代理 / 广告拦截 / TUN 开关
  3  代理 & 订阅        换代理 · 切节点 · 连通测试 · 全局扩展脚本
  4  启动 / 重启 / 停止
  5  看日志
  6  关闭 gateway 并退出（停 mihomo）
  Q  退出控制台（mihomo 留在后台继续跑）
```

- `流量控制` → **`分流 & 规则`**
- `换代理源` → **`代理 & 订阅`**
- `增强脚本` → **`全局扩展脚本`**（对齐 Clash 用户习惯）
- 右侧白话说明一眼看清能干啥

## 2. 内嵌 Web 控制台 · `metacubexd`

```
Web 控制台（浏览器打开，切节点 / 查流量 / 改规则）
  本机    http://127.0.0.1:19090/ui/          ✓
  局域网  http://192.168.12.100:19090/ui/   ✓   （手机 / 平板也能用）
```

- `metacubexd` 的 gh-pages dist（2.2 MB）直接 `go:embed` 进 binary
- gateway 启动时 `deployWebUI` 自动释放到 `workdir/ui/`，`template.yaml` 补 `external-ui: ui`，`/ui` 开箱即用
- 浏览器里能切节点、改模式、看连接、查规则、读日志、改 DNS —— 主菜单里能做的事都能在 UI 里做
- binary 从 ~6.5MB 长到 ~14MB（加上 goja / supervisor / metacubexd 等）

## 3. 代理源层重构：订阅 / 文件改成 inline 模式

以前订阅和本地文件走 `proxy-provider` 模式，mihomo 只读 `proxies:` 字段，**订阅自带的 `proxy-groups` 和 `rules` 全部被扔掉**，用户看到订阅配的"飞天猪"之类分组根本不存在、规则集也没生效，主人真机测了是"怎么用都没效果"。

现在改成 **inline**：
- 读用户 yaml → 抽 `proxies + proxy-groups + rules` → 直接嵌进最终 mihomo config.yaml
- 用户订阅的 `DOMAIN-SUFFIX,xxx,🚀 节点选择` 这种规则直接生效
- 顶级字段（`mode / dns / tun / port / external-controller`）**不从订阅读**，由 gateway 的 base template 渲染
- 如果用户 yaml 没 `Proxy` 组，自动补一个 select 兜底

## 4. 单点代理合并 + 每项自动连通性测试

原本 5 种源（external / subscription / file / remote / none）压成 4 种：
- **单点代理**（external + remote 合并，填用户名才启用认证）
- **机场订阅**
- **本地配置文件**（原「Clash 配置文件」改名）
- **暂不配置**

「代理 & 订阅」菜单进入时并发探测每项可达性：

```
    1  单点代理      ✓  127.0.0.1:6578 socks5
    2  机场订阅      ·  (未配置)
    3  本地配置文件  ✓  ~/Documents/clash/long.yaml
    4  暂不配置      ·  全部走直连
```

- 单点代理 → TCP 拨号
- 订阅 → HTTP GET 状态码
- 本地文件 → 存在性 + YAML 粗检
- 按 `T` 重跑

## 5. 切换节点菜单 · 自动测速 + 延迟排序

```
分组：🛬 AI落地节点  (当前：🏠 住宅IP)

   1  ✓ 🏠 住宅IP                    186 ms
   2    🛫 AI起飞节点                 —

  ── 操作 ── R 刷新测速   <编号> 切到该节点   0 返回
```

- 进入分组自动调 mihomo `/group/<name>/delay` 测速（8s 超时）
- 按延迟升序排（超时丢后）；< 300 ms 绿 / < 1000 ms 黄 / 其余红 / 灰 `—`
- 按 R 重测，选数字立即切节点

## 6. 链式代理预设（住宅 IP 落地）

`S 全局扩展脚本 → 1 预设 · 链式代理`，引导式填住宅 IP 节点字段：

```
请填写住宅 IP 落地节点：
  节点名称（回车=🏠 住宅IP）: 
  类型 (http / socks5)（回车=socks5）: 
  服务器地址: xxx.xxx.xxx.xxx
  端口（回车=443）: 
  用户名（无需认证回车）: 
  密码: 
```

运行时把用户填的值实例化到内嵌的 `residential-chain.tmpl.js` 模板 → 生成 `workdir/residential-chain.js` → goja 在最终 mihomo config 上跑。

脚本功能：聚合订阅节点为「🛫 AI起飞节点」组，聚合住宅 IP 为「🛬 AI落地节点」，把 Claude / OpenAI / Cursor / Termius 等 AI 相关域名路由到 AI落地节点（链式代理），阿里系直连。

`2 自定义 .js 文件路径` 入口仍保留给高级用户自己写脚本（`function main(config) { return config }` 格式）。

## 7. 自定义规则编辑器

「分流 & 规则 → 4 自定义规则」进入：

```
自定义规则
  [直连] (0)  (无)
  [走代理] (1)
    1  DOMAIN-SUFFIX,github.com
  [拒绝] (0)  (无)

  ── 操作 ── A 添加一条   D <编号> 删除某条   0 返回
```

- `A` 添加 → 选类型（DOMAIN-SUFFIX / DOMAIN / DOMAIN-KEYWORD / IP-CIDR / PROCESS-NAME / 手写）→ 填目标 → 选去向（直连 / 代理 / 拒绝）
- `D 3` 删除第 3 条
- 自动 save + reload；优先级最高（比内置 `china_direct` / `adblock` 更靠前）

## 8. 代理源自愈 · supervisor 自动 fallback 到 direct

后台 goroutine 每 30 秒跑一次 `source.Test`：

- 代理上游挂了（7890 connection refused / 订阅 502 / 文件消失等）→ 立即通过 mihomo `/configs` PATCH 把 `mode` 切到 `direct`，LAN 设备继续走直连，**不断网**
- 恢复后自动切回原 mode
- 主菜单状态栏红字高亮告警；ticker 2 秒检测 health 状态变化立即重绘（用户不用戳菜单就能看到）

```
● 运行中    本机 IP: 192.168.12.100
代理源: 单点代理（本机，external）
⚠ 代理源异常 · 已临时切到直连（LAN 设备不会断网，但不再走代理）
  原因: 连不上 127.0.0.1:7890 → ...
  修复后会自动切回
```

## 9. 设备接入指引 · 两种方式并排

「1 设备接入指引」展示两种接入方式，明写端口：

- **方式 1 · TUN 网关**（Switch / PS5 / Apple TV）：改网关 + DNS 到本机 IP
- **方式 2 · 局域网代理**（iPhone / 电脑 App）：填 `本机 IP:17890`（HTTP + SOCKS5 混合）

明说差异：Switch 改网关能走 YouTube，App 填代理只走自己的代理流量。

## 10. 日志菜单 · 易读视图

默认中文简要视图（之前版本加的，顺带提一下），按 `r` 切原始：

```
🟡 01:27:55  TCP 直连 115.190.130.195:21114 （GeoIP=cn）  → 目标无响应（超时）
🟡 01:25:53  UDP 走代理[Proxy] 17.248.216.69:443 （兜底）  → 本机代理拒绝连接（检查代理软件是否启动）
```

识别 DIRECT / Proxy / REJECT + GeoIP / Match / DomainSuffix / Keyword / IPCIDR / ProcessName + 常见错误。

## 11. 默认端口避开 Clash

运行时端口从 `mixed=7890 / redir=7892 / api=9090` 改成 `17890 / 17892 / 19090`，避开 Clash / mihomo / Clash Verge 常占端口。`Source.External.Port` 默认仍是 `7890`（要和本机 Clash 的上游端口保持一致）。

## 12. 其他修复

- mihomo 顶级项（`external-ui` / `tun` / `port` 等）API reload 不更新，改成 `Reload` 统一走 Stop+Start
- 本地配置文件先复制到 workdir 内再让 mihomo 读，绕开 mihomo 的 SAFE_PATHS 检查
- template 补 `proxy-server-nameserver`（国内 DNS），避免代理服务器域名自己查自己的 DNS 死锁
- 主菜单从 `sudo gateway` 进来自动 reattach 到后台 mihomo（v3.0.1 做过一次，这版保持）

## 下载

| 系统 | 直接下载 |
|---|---|
| macOS Apple Silicon | [gateway-darwin-arm64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.1.0/gateway-darwin-arm64) / [.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.1.0/gateway-darwin-arm64.tar.gz) |
| macOS Intel | [gateway-darwin-amd64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.1.0/gateway-darwin-amd64) / [.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.1.0/gateway-darwin-amd64.tar.gz) |
| Linux x86_64 | [gateway-linux-amd64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.1.0/gateway-linux-amd64) / [.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.1.0/gateway-linux-amd64.tar.gz) |
| Linux ARM64 | [gateway-linux-arm64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.1.0/gateway-linux-arm64) / [.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.1.0/gateway-linux-arm64.tar.gz) |
| Windows x86_64 | [gateway-windows-amd64.exe](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.1.0/gateway-windows-amd64.exe) / [.zip](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.1.0/gateway-windows-amd64.zip) |

校验文件: [SHA256SUMS](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v3.1.0/SHA256SUMS)

## 兼容性 / 升级提示

- **老 `gateway.yaml` 继续认**，external / remote / file / subscription 所有 type 保留；UI 展示统一成新名字
- **老默认端口**（7890/7892/9090）在已有 config 里会被 Load 的兜底逻辑保留，不会强制改；要换新端口去「分流 & 规则 → 9 高级」
- **开机自启** 如果从 3.0.2 之前跨版本升，建议 `sudo gateway service uninstall && sudo gateway service install` 重装一下（新版 service unit 用 `start --foreground`）
