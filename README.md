# LAN Proxy Gateway

[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/Tght1211/lan-proxy-gateway)](https://github.com/Tght1211/lan-proxy-gateway/releases)
[![License](https://img.shields.io/github/license/Tght1211/lan-proxy-gateway)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey)]()

> **把一台电脑变成整屋的代理网关** — Switch / PS5 / Apple TV / iPhone / 智能电视，改个网关 + DNS 就能科学上网，不用每台设备单独装代理 App。
>
> 🧠 **AI 场景加分**：一键配「机场起飞 + 住宅 IP 落地」链式代理 —— Claude / OpenAI / Cursor / Termius / ping0.cc 自动走**纯净住宅 IP**出口，不触发风控不封号。

面向**非编程玩家**：一键安装、中文菜单、内嵌 Web 控制台、配置全程引导式。

---

## ✨ 能干什么

### 🌟 招牌能力 · 纯净住宅 IP 链式出口

AI 网站（Claude / OpenAI / Cursor）对机房 IP 风控越来越狠，用普通机场节点经常登录就被拉黑、被要求验证、甚至封号。这个项目**内置链式代理预设**：

- **流量路径**：浏览器 / Cursor → 本机 → 机场节点（起飞，解决跨境带宽）→ 住宅 IP（落地，**真实家庭宽带出口**） → AI 网站
- **AI 网站看到的是住宅 IP**（/Comcast / AT&T / 家用 ASN），风控系统识别为"普通用户家里访问"
- **一键填就好**：菜单 → `3 代理 & 订阅 → S 全局扩展脚本 → 1 预设 · 链式代理`，填住宅 IP 服务器 / 端口 / 账号密码就行；规则（Claude 域名走住宅、其它走机场、国内直连）自动生成
- **不误伤其它流量**：YouTube / Google Drive 这类走机场，国内网站走直连，住宅 IP 只给 AI 专用 —— 流量账按最省的方式分

> 配合本机 DNS 切到 127.0.0.1（方式 3），Claude Code CLI / Cursor Agent 跑在本机上时也会自动走住宅 IP 出口。

### 其它能力

- 🏠 **LAN 透明网关** — 基于 mihomo TUN 劫持，家里的 Switch / PS5 / 智能电视只要把网关和 DNS 指向这台电脑就自动走代理（**macOS / Linux；Windows 家用版只支持方式 2 设代理**）
- 📱 **混合代理端口** — 同时开 HTTP + SOCKS5（默认 `17890`，避开 Clash 7890），iPhone / Android / 浏览器 App 直接填当代理（**全平台可用，Windows 唯一推荐方式**）
- 🌐 **内嵌 metacubexd Web 控制台** — 浏览器打开 `http://ip:19090/ui/`，切节点 / 改规则 / 看流量一应俱全，手机平板也能进
- ⚡ **自动自愈** — 代理源挂了 30 秒内 supervisor 自动切到直连保命（LAN 不断网），恢复后切回
- 🎯 **自定义规则 UI** — 菜单里增删 DOMAIN-SUFFIX / IP-CIDR / PROCESS-NAME 规则，优先级盖过内置
- 📊 **节点测速 + 排序** — 进切节点页面自动并发测延迟，按速度升序
- 🗒️ **日志易读视图** — mihomo 英文日志自动翻译成中文简要（`🟡 01:27:55 TCP 直连 xxx → 目标无响应（超时）`）
- 💻 **方式 3 · 本机也走规则** — macOS 一键切本机 DNS 到 `127.0.0.1`（自动遍历 Wi-Fi / Ethernet / Thunderbolt 全部活跃网卡），本机浏览器也能按域名分流
- 🔧 **小白友好** — 流量 / 代理 / 规则菜单重分类，右侧白话描述，告警自动重绘

---

## 🚀 一键安装

### macOS / Linux

```bash
curl -fsSL https://raw.githubusercontent.com/Tght1211/lan-proxy-gateway/main/install.sh | bash
```

脚本装完自动进入配置向导（问代理源 → 启动 → 问开机自启），整条流在一个终端里走完。

### Windows（管理员 PowerShell）

```powershell
irm https://raw.githubusercontent.com/Tght1211/lan-proxy-gateway/main/install.ps1 | iex
```

### 国内访问 GitHub 慢？

脚本 + `gateway install` 下载 mihomo 时会按顺序尝试镜像（`ghfast.top` / `hub.gitmirror.com` / `github.moeyy.xyz` / `ghp.ci`）。也可以手动指定：

```bash
GITHUB_MIRROR=https://你的镜像/ bash install.sh           # Linux / macOS
$env:GITHUB_MIRROR = "https://你的镜像/"; gateway install   # Windows PowerShell
```

如果本机已经有 Clash Verge / Clash for Windows 在跑，也可以直接让下载走它：

```powershell
# PowerShell —— 走 Clash Verge 默认混合端口
$env:HTTP_PROXY  = "http://127.0.0.1:7897"
$env:HTTPS_PROXY = "http://127.0.0.1:7897"
gateway install
```

---

## 🎬 整体架构

```mermaid
flowchart LR
    subgraph Devices["📱 不便装代理的设备"]
        S[Switch]
        P[PS5]
        T[Apple TV]
        TV[智能电视]
        M[手机 / 平板]
    end

    subgraph Host["🖥️ 跑 gateway 的电脑"]
        direction TB
        TUN[TUN 劫持]
        MH[mihomo 内核<br/>规则分流 / 广告拦截]
        UI[metacubexd<br/>Web 控制台]
        TUN --> MH
        MH --> UI
    end

    subgraph Sources["🌐 代理源（任选其一）"]
        E1[单点代理<br/>本机 / 远程]
        SUB[机场订阅<br/>URL]
        F[本地配置文件<br/>.yaml]
        SCR[全局扩展脚本<br/>链式代理预设]
    end

    Devices -- "网关+DNS → 电脑 IP" --> Host
    Host --> Sources

    style Devices fill:#fff5e6,stroke:#ff9900
    style Host fill:#e6f3ff,stroke:#0066cc
    style Sources fill:#e6ffe6,stroke:#00aa00
```

**核心思路**：电脑只要能科学上网，把它变网关就能让整屋设备一起享受，**不用每台设备重配一份订阅**。

---

## 📱 三种设备接入方式

> ⚠️ **Windows 限制**
>
> **方式 1（改网关）在 Windows 家用版上不可用**。Mac/Linux 靠 iptables/pf 做 NAT，家用版 Windows 没有 RRAS，ICS 又强制 `192.168.137.0/24` 不能自定义 —— 项目里 Windows 分支的 `ConfigureNAT` 是 no-op。手机 / iPad 请直接用**方式 2**（设 HTTP 代理）。
>
> Switch / PS5 / Apple TV / 智能电视这类**不能设代理只能改网关**的设备，在 Windows 上走不通，请改用 Linux / macOS 主机，或买个软路由。

### 📺 方式 1 · 改网关（适合 Switch / PS5 / Apple TV / 智能电视）

> ⛔ **仅 macOS / Linux**。Windows 用户请跳到方式 2。

这类设备**只能填网关 + DNS，不能填代理**，用这个方式：

| 字段 | 填什么 |
|---|---|
| 网关 (Gateway) | 电脑的局域网 IP |
| DNS 服务器 | 同一个 IP |
| 子网掩码 | `255.255.255.0` |
| IP 地址 | 保留 DHCP 或设静态 IP |

保存 → 重连 Wi-Fi，所有流量（YouTube / 游戏 / 各类 App）自动走代理。

### 📱 方式 2 · 填代理（适合 iPhone / Android / 电脑 App / 浏览器插件）

能**手动填代理服务器**的场景 —— **Windows 下唯一能用的方式**：

| 字段 | 填什么 |
|---|---|
| 代理服务器 | 电脑的局域网 IP |
| 端口 | `17890`（HTTP + SOCKS5 混合） |
| 类型 | HTTP 或 SOCKS5（端口同一个） |
| 用户名 / 密码 | 留空 |

具体到设备：
- **Android**：Wi-Fi 设置 → 长按当前网络 → 修改网络 → 高级选项 → 代理 = 手动
- **iOS / iPad**：Wi-Fi 设置 → 点当前网络右侧 (i) → 配置代理 = 手动
- **Windows / macOS 浏览器**：装 Proxy SwitchyOmega 之类扩展，填 `电脑IP:17890`

只走 App 自己发到代理的流量（Switch / PS5 填这个没用 —— 主机不会主动把流量交给代理）。

### 💻 方式 3 · 本机自己也走规则

gateway 所在电脑的浏览器 / App 也想按规则走代理 — 两条路：

**① 开着 TUN（推荐，Windows 唯一可用的路）**

TUN 接管本机所有出向流量，浏览器不用改任何代理设置。**装完默认就是这状态**。

> ⚠️ TUN 只劫持 TCP/UDP **数据包**，DNS 查询如果直连路由器，mihomo 看到的是 IP —— `GEOIP,CN,DIRECT` 能命中，但 `DOMAIN-SUFFIX,github.com,Proxy` 这类域名规则**命不中**。纯国内 / 国外分流够用；想让域名规则精确生效，继续看 ②。

**② 把本机系统 DNS 改到 `127.0.0.1`（走 mihomo fake-ip，让域名规则全部生效）**

- **macOS 一键**：主菜单 → `1 设备接入指引` → 按 `L`。代码会 `networksetup -listallnetworkservices` 枚举**所有活跃网络服务**（Wi-Fi / Ethernet / Thunderbolt 等全部一起改），插着网线又开 Wi-Fi 也能覆盖到 —— 不需要自己选。恢复按 `R`。

- **Linux**：一键切换没实现（各发行版 systemd-resolved / NetworkManager / resolvconf 风格差异大）。手动命令示例：
    ```bash
    # systemd-resolved（Ubuntu 18.04+ / Fedora）
    for link in $(nmcli -t -f DEVICE con show --active | cut -d: -f1); do
        sudo resolvectl dns "$link" 127.0.0.1
    done
    # NetworkManager（不管有没有 systemd-resolved）
    sudo nmcli con mod "Wired connection 1" ipv4.dns 127.0.0.1 && sudo nmcli con up "Wired connection 1"
    ```
  多连接场景记得每个 active 连接都改，否则系统按优先级选到没改的那个就没效。

- **Windows**：一键切换没实现；也不推荐改系统 DNS（netsh 每个接口独立、改回麻烦、还要处理 IPv6）。**推荐保持 TUN 开着**；要强制域名规则生效，浏览器装 SwitchyOmega 填 `127.0.0.1:17890` 即可。

- **验证**（macOS）：`dscacheutil -q host -a name ping0.cc` 返回 `198.18.x.x`（fake-ip）即 DNS 已走 mihomo。

### 💡 方式对比

| 方式 | 覆盖面 | 需要 TUN | 适合设备 | macOS | Linux | Windows |
|---|---|---|---|---|---|---|
| 1 改网关 | 所有流量 | ✅ | Switch / PS5 / TV | ✅ | ✅ | ❌ |
| 2 填代理 | 仅支持代理的 App | ❌ | iPhone / Android / 浏览器 | ✅ | ✅ | ✅ |
| 3 本机自己 | 本机全部 | TUN 开=✅ / 关=改 DNS | 跑 gateway 这台电脑 | ✅ | ✅ | ✅（TUN 开着自动） |

---

## 🌐 Web 控制台（`metacubexd` 内嵌）

订阅 / 本地文件源时，浏览器打开：

```
http://<电脑局域网 IP>:19090/ui/
```

内置功能：

- 📊 切代理组 / 节点（**自动测延迟 + 颜色标注**）
- 🔄 切换模式（rule / global / direct）
- 📈 实时连接 / 流量统计
- 📜 查规则命中 / 读日志
- ⚙️ 改 mihomo 运行时配置

`metacubexd` 的 dist 已经 `go:embed` 进 binary（2.2 MB），gateway 启动自动释放到 `workdir/ui/`，**开箱即用**，不需要自己 clone UI 仓库。

---

## 🖥️ 主菜单一览

```
────────────────────────────────────────────────
  LAN 代理网关
────────────────────────────────────────────────
  ● 运行中    本机 IP: 192.168.12.100
  代理源: 本地配置文件（file）

  1  设备接入指引      Switch / PS5 / 手机怎么连到这里
  2  分流 & 规则        国内直连 / 国外走代理 / 广告拦截 / TUN 开关
  3  代理 & 订阅        换代理 · 切节点 · 连通测试 · 全局扩展脚本
  4  启动 / 重启 / 停止
  5  看日志
  6  关闭 gateway 并退出（停 mihomo）
  Q  退出控制台（mihomo 留在后台继续跑）
```

每个子菜单都用「── 操作 ── 横排按键」结尾，统一体验。进入「3 代理 & 订阅」自动并发测每种源可达性，「N 切换节点」自动跑 mihomo `/group/xxx/delay` 并按延迟排序。

---

## 🔧 典型场景

### 场景一：本机已经装着 Clash Verge，想让手机也用这份代理

安装 → 选 `1) 单点代理` → 填 `127.0.0.1` + Clash 的端口（通常 `7890`）→ 启动。手机按「方式 2 填代理」填 `电脑IP:17890`，或按「方式 1 改网关」。

### 场景二：纯机场订阅用户

安装 → 选 `2) 机场订阅` → 粘订阅 URL → 启动。所有订阅规则、分组、节点**全部 inline** 到 mihomo（v3.1+ 改了渲染方式，不再走 proxy-provider，机场自定义分组完整生效）。

### 🌟 场景三：AI 服务走纯净住宅 IP 链式代理（招牌玩法）

**用机场节点访问 Claude / OpenAI 总是被风控？换这招。**

**问题**：
- 机场的出口 IP 都是**机房 ASN**（Digital Ocean / Vultr / AWS 之类），OpenAI / Anthropic 的风控系统识别度极高 —— 轻则要求验证、重则直接封号
- 但机场节点**带宽好、跨境路径稳**，直接扔掉换住宅 IP 又会卡到天荒地老
- 单独跑一个住宅 IP 落地节点？流量**全**走住宅，YouTube / Google Drive 跟着一起 —— 家庭宽带跑流量扛不住

**这个项目的解法**（`presets.RenderResidentialChain`，代码在 `internal/script/presets/`）：

```
你的浏览器 / Cursor / Claude Code
          │
          ↓
    本机 gateway (TUN 抓取)
          │
          ↓
  🛫 AI起飞节点 (机场, 带宽好)  ← 跨境这一段走这里，快
          │
          ↓
  🛬 AI落地节点 (住宅 IP)       ← 最后一跳是家庭宽带，AI 网站看到的是"普通家用 IP"
          │
          ↓
    Claude / OpenAI / Cursor
```

**操作步骤**：
1. 先配好订阅源（场景二）
2. 主菜单 → `3 代理 & 订阅 → S 全局扩展脚本 → 1 预设 · 链式代理`
3. 填住宅 IP 信息（服务器 / 端口 / 用户名 / 密码，SOCKS5 或 HTTP 都支持）

**自动生成**：
- `🛫 AI起飞节点` 组 = 订阅的所有机场节点（可按需切换）
- `🛬 AI落地节点` 组 = 住宅 IP 节点，`dialer-proxy` 指向 🛫 AI起飞节点（先机场后住宅的链式连接）
- **规则**：Claude / OpenAI / Anthropic / Cursor / Termius / ping0.cc / openai.com / anthropic.com / generativelanguage.googleapis.com 等 AI 相关域名 → `🛬 AI落地节点`
- **不误伤**：YouTube / Google Drive / GitHub 等走 `Proxy`（机场）；国内流量走 `DIRECT`；住宅 IP 带宽只给 AI 独享

**验证出口 IP**：
- 本机跑 AI CLI / Claude Code / Cursor —— 切到方式 3（macOS 菜单按 L，Windows TUN 开着就行），访问 `https://ping0.cc` 看到的就是住宅 IP
- 服务端验证：进 `https://whoer.net` 看 ASN 类型是 "Residential" 就对了

### 场景四：本机浏览器直接验证出口 IP

**macOS**：按「方式 3」切本机 DNS → 浏览器开 `https://ping0.cc`，应该显示住宅 IP（不是机场 IP、不是你家宽带 IP）。

**Windows**：TUN 开着时（默认）浏览器自动走 mihomo；开 `https://ping0.cc` 看到的也应该是住宅 IP。如果想让 `DOMAIN-SUFFIX` 规则精确命中（只让 `ping0.cc` 走住宅，其它走机场）而不是靠 GeoIP 兜底，浏览器装个 SwitchyOmega 指 `127.0.0.1:17890`。

---

## ⚙️ 配置文件

`~/.config/lan-proxy-gateway/gateway.yaml`（Windows：`%APPDATA%\lan-proxy-gateway\`）：

```yaml
version: 2

gateway:
  enabled: true
  tun: { enabled: true, bypass_local: false }
  dns: { enabled: true, port: 53 }

traffic:
  mode: rule                # rule | global | direct
  adblock: true
  rulesets:
    china_direct: true
    apple: true
    nintendo: true
    global: true
    lan_direct: true
  extras:                   # 自定义规则（优先级最高）
    direct: []
    proxy:  []
    reject: []

source:
  type: file                # external | subscription | file | remote | none
  file:
    path: /Users/xxx/clash/long.yaml
  # script_path: ...            # 手写 .js 扩展脚本路径（与 chain_residential 互斥）
  # chain_residential:          # 链式代理预设（用菜单填更方便）
  #   name: "🏠 住宅IP"
  #   kind: socks5
  #   server: 1.2.3.4
  #   port: 443

runtime:
  ports: { mixed: 17890, redir: 17892, api: 19090 }  # 避开 Clash 7890/7892/9090
  log_level: warning
```

完整示例见 [`gateway.example.yaml`](gateway.example.yaml)。v1 / v2 配置在首次加载时自动升级到 v3 并打印迁移报告。

---

## 🔌 命令行

所有操作都能在**主菜单**里完成，不用记命令。下面是供系统服务 / 脚本化用的：

| 命令 | 作用 |
|---|---|
| `gateway` | 进主菜单（首次会先跑配置向导） |
| `gateway install` | 下载 mihomo + 向导 + 启动 + 问开机自启 |
| `gateway start` | 默认后台启动；`--foreground` 给 launchd/systemd 用 |
| `gateway stop` | 停止 mihomo |
| `gateway status` | 一次性输出当前状态 + 设备接入指引 |
| `gateway service install` | 安装为开机自启（launchd / systemd / schtasks） |
| `gateway service uninstall` | 卸载系统服务 |
| `gateway service status` | 查看服务状态 |

---

## 🌍 跨平台支持

| 系统 | IP 转发 | NAT | 服务管理 | 本机 DNS 一键切 |
|---|---|---|---|---|
| **macOS** | `sysctl net.inet.ip.forwarding` | `pfctl` | `launchd` plist | ✅ `networksetup` |
| **Linux** | `/proc/sys/net/ipv4/ip_forward` | `iptables MASQUERADE` | `systemd` unit | ⚠️ 手动（见指引） |
| **Windows** | 注册表 `IPEnableRouter=1` | mihomo TUN 虚拟网卡 | `schtasks` 计划任务 | ⚠️ 手动（见指引） |

编译产物 5 平台均通过：`darwin-arm64 / darwin-amd64 / linux-amd64 / linux-arm64 / windows-amd64`。

> Linux / Windows 欢迎 PR 贡献一键 DNS / 更完善的服务管理。

---

## 🧩 三层架构

```mermaid
flowchart TB
    subgraph L1["【主】gateway.*"]
        G1[LAN 网关<br/>IP 转发]
        G2[TUN<br/>劫持全部流量]
        G3[DNS 代理<br/>端口 53]
    end

    subgraph L2["【副】traffic.*"]
        T1[规则模式<br/>国内直连国外代理]
        T2[全局 / 直连]
        T3[广告拦截]
        T4[自定义规则]
    end

    subgraph L3["【拓展】source.*"]
        P1[单点代理]
        P2[机场订阅]
        P3[本地文件]
        P4[全局扩展脚本<br/>链式代理预设]
    end

    subgraph L4["【健壮性】supervisor"]
        H1[30s 源健康检查]
        H2[挂了自动切 direct]
        H3[恢复切回原 mode]
    end

    L1 --> L2 --> L3
    L2 -.被 H1-3 保护.-> L4

    style L1 fill:#ffe0e0,stroke:#cc0000
    style L2 fill:#fff0d9,stroke:#cc6600
    style L3 fill:#e0f0ff,stroke:#0066cc
    style L4 fill:#f0e6ff,stroke:#8800cc
```

---

## 🛠️ 手动编译

```bash
git clone https://github.com/Tght1211/lan-proxy-gateway
cd lan-proxy-gateway

make build            # 当前平台 → ./gateway
make install          # 装到 /usr/local/bin/gateway（需 sudo）
make test             # 全部单元测试
make build-all        # 交叉编译 darwin / linux / windows
```

### 目录结构（v3.1）

```
cmd/                cobra 入口（5 个命令）
internal/
  app/              统一门面（console + cobra 共用）+ supervisor（代理源自愈）
  gateway/          【主】LAN 网关 + 设备接入指引
  traffic/          【副】规则 + 内置 ruleset + 自定义合并
  source/           【拓展】代理源 inline + 连通性测试
  engine/           mihomo 进程 + 渲染 + REST API（SelectNode / GroupDelay / SetMode）
  script/           goja 脚本执行器
    presets/        内嵌预设（链式代理 · 住宅 IP 落地）
  config/           v3 schema（v1/v2 自动迁移）
  platform/         跨平台（darwin/linux/windows）
  console/          菜单式交互 + 日志易读视图 + 显示宽度对齐
  mihomo/           下载 mihomo 内核
embed/
  template.yaml     mihomo config 模板
  webui/            metacubexd dist（2.2 MB，go:embed）
legacy/v1/          v1 源码留档
```

---

## 🤖 给 AI 用的 Skill

`~/.claude/skills/lan-proxy-gateway-ops/SKILL.md`：AI 代理通过 mihomo REST API 做日常运维（切节点 / 换模式 / 查日志）**无需 sudo**；需要 sudo 的动作走 sudoers NOPASSWD 白名单或建议走系统 service。

---

## 🤝 贡献

欢迎 issue / PR！

- Linux / Windows 一键 DNS 切换的实现
- 新的 ruleset 内置规则
- 英文 README / 文档
- mihomo 新 API 的 UI 接入

---

## 📜 License

[MIT](LICENSE) © 2025-2026 [Tght1211](https://github.com/Tght1211)

基于 [mihomo](https://github.com/MetaCubeX/mihomo)（Clash.Meta）内核 + [metacubexd](https://github.com/MetaCubeX/metacubexd) 控制台。

---

## ⭐ Star History

[![Star History Chart](https://api.star-history.com/svg?repos=Tght1211/lan-proxy-gateway&type=Date)](https://star-history.com/#Tght1211/lan-proxy-gateway&Date)

如果觉得有用，点个 Star ⭐ 支持一下吧~
