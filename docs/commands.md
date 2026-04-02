# 完整命令说明

## 基础操作

| 命令 | 说明 | 需要管理员权限 |
|------|------|:---:|
| `gateway install` | 初始化向导：下载 mihomo、配置代理来源、生成配置文件 | 否 |
| `sudo gateway start` | 启动网关（开启 IP 转发、启动 mihomo） | ✓ |
| `sudo gateway stop` | 停止网关（关闭 mihomo、清理路由规则） | ✓ |
| `sudo gateway restart` | 重启网关（等同于 stop + start） | ✓ |
| `gateway status` | 查看运行状态：节点、连接数、上下行流量、本机 IP | 否 |

## 维护与升级

| 命令 | 说明 | 需要管理员权限 |
|------|------|:---:|
| `sudo gateway health` | 健康检查：检测进程 / TUN 接口 / API，异常时自动重启修复 | ✓ |
| `sudo gateway update` | 一键升级到最新版本（自动下载、替换、重启），GitHub 超时自动切换镜像 | ✓ |

## 切换代理来源

| 命令 | 说明 |
|------|------|
| `gateway switch` | 查看当前代理来源（url 或 file 模式） |
| `gateway switch url` | 切换到订阅链接模式（需在 `gateway.yaml` 中已配置 `subscription_url`） |
| `gateway switch file /path/to/config.yaml` | 切换到本地 Clash/mihomo 配置文件模式，自动提取 proxies 节点 |

## TUN 模式控制

| 命令 | 说明 |
|------|------|
| `gateway tun on` | 开启 TUN 透明代理模式（**Switch / PS5 / 电视 要走代理必须开启**），需重启生效 |
| `gateway tun off` | 关闭 TUN 模式（默认），改为规则模式 |

> TUN 模式是什么：开启后 mihomo 会创建一块虚拟网卡，接管所有流经网关的流量，实现真正的透明代理。**让其他设备科学上网必须开启 TUN 模式。** 如果本机同时运行 Clash Verge 且也开了 TUN，两者会冲突——先关掉 Clash Verge 的 TUN。

## 开机自启动服务

| 命令 | 说明 | 需要管理员权限 |
|------|------|:---:|
| `sudo gateway service install` | 安装系统服务：开机自动启动 + 崩溃自动重启（macOS 用 launchd，Linux 用 systemd，Windows 用 sc.exe） | ✓ |
| `sudo gateway service uninstall` | 卸载系统服务 | ✓ |

## 全局参数

所有命令都支持以下参数（放在命令后面）：

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--config <路径>` | 指定配置文件路径 | 当前目录的 `gateway.yaml` |
| `--data-dir <路径>` | 指定数据目录路径（存放 mihomo 运行配置） | 当前目录的 `data/` |

**示例：**

```bash
# 使用自定义配置文件路径
sudo gateway start --config /etc/gateway/gateway.yaml --data-dir /var/lib/gateway
```
