# LAN Proxy Gateway v2.2.11

这一版把 simple 模式继续往“日常真的能用”的方向收。重点不是再堆命令，而是把高频操作收进工作台，并让节点选择本身带上延时、重测和排序能力。

## 重点变更

### 1. simple 模式帮助收口为高频操作

- 默认 `help` 现在只展示常用入口，例如 `nodes`、`subscription`、`extension`、`chain`
- 低频命令没有删掉，而是统一收进 `help all`
- 这样进入 simple 模式后，不会再先看到一整屏命令清单

### 2. simple 工作台快捷键真正可用

- 打开 `subscription / extension / chain` 后，面板里的 `1 / 2 / A / S / R ...` 现在可以直接输入
- 缺参数时会继续提示输入，而不是把快捷键当成未知命令
- 纯命令模式终于更像“工作台”，不需要反复手敲整条命令

### 3. `nodes` 改成可测速的节点工作台

- 进入分组后，会展示每个节点的延时
- 支持输入 `T` 主动重测整组节点
- 测速完成后按低延时优先排序，失败项放到最后
- 同时会过滤 `剩余流量 / 套餐到期 / 最新网址 / 重置剩余` 这类订阅说明项，避免它们混进可切换节点列表

### 4. README / 命令手册 / release 说明同步更新

- README 和 README_EN 现在写明了 simple 模式下的高频帮助、工作台快捷键和 `nodes` 延时排序行为
- `docs/commands.md` / `docs/en/commands.md` 也补上了 pure console 的节点工作台说明
- release 说明和版本规划同步更新，方便后面继续按 patch 节奏发版

## 下载

| 系统 | 直接下载 |
|---|---|
| macOS Apple Silicon | [gateway-darwin-arm64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.11/gateway-darwin-arm64) / [gateway-darwin-arm64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.11/gateway-darwin-arm64.tar.gz) |
| macOS Intel | [gateway-darwin-amd64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.11/gateway-darwin-amd64) / [gateway-darwin-amd64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.11/gateway-darwin-amd64.tar.gz) |
| Linux x86_64 | [gateway-linux-amd64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.11/gateway-linux-amd64) / [gateway-linux-amd64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.11/gateway-linux-amd64.tar.gz) |
| Linux ARM64 | [gateway-linux-arm64](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.11/gateway-linux-arm64) / [gateway-linux-arm64.tar.gz](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.11/gateway-linux-arm64.tar.gz) |
| Windows x86_64 | [gateway-windows-amd64.exe](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.11/gateway-windows-amd64.exe) / [gateway-windows-amd64.zip](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.11/gateway-windows-amd64.zip) |

校验文件: [SHA256SUMS](https://github.com/Tght1211/lan-proxy-gateway/releases/download/v2.2.11/SHA256SUMS)

## 建议验证

升级到 `v2.2.11` 后，建议顺手确认这几处：

1. 执行 `sudo gateway console`，输入 `help`，确认默认只显示高频操作；再输入 `help all`，确认完整命令仍然都在
2. 输入 `chain`、`extension` 或 `subscription` 进入工作台后，直接按面板里的 `1 / 2 / A / S / R ...`，确认快捷键会被识别并继续提示输入
3. 输入 `nodes`，选择一个节点分组，确认会展示每个节点延时，并且输入 `T` 后会重新测速并按延时排序
4. 确认 `剩余流量 / 套餐到期 / 最新网址` 这类订阅说明项不会再混进可切换节点列表
