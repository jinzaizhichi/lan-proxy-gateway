# 进阶配置

## 开机自启与稳定性

安装系统服务后，电脑重启后网关自动启动，崩溃后自动恢复：

```bash
sudo gateway service install
```

内置稳定性保障：

| 机制 | 说明 |
|------|------|
| 节点健康检查 | 每 120 秒检测节点可用性，失效自动切换 |
| 进程崩溃自愈 | launchd / systemd 检测到崩溃后自动重启 |
| 定时健康检查 | 每天 4:00 和 12:00 自动执行健康检查 |
| 日志轮转 | 自动保留最近 3 份日志，防止磁盘占满 |
| 一键升级 | `sudo gateway update` 自动更新到最新版本 |

## 扩展脚本（高级用户）

如果你有特殊需求，比如：
- 让 ChatGPT / Claude 单独走住宅 IP 节点
- 把公司内网域名强制直连
- 在订阅规则前面插入自定义规则

可以通过扩展脚本实现。脚本格式与 Clash Verge Rev 完全兼容。详情见 [gateway.example.yaml](../gateway.example.yaml) 中的 `script_path` 配置说明。
