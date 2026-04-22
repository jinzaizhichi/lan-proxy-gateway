package gateway

import (
	"fmt"
	"runtime"
	"strings"
)

// DeviceGuide 返回给设备接入者看的紧凑说明。
// mixedPort 是 mihomo 的 HTTP+SOCKS5 混合端口（方式 2 要填）。
//
// 平台差异很大，尤其 Windows：这台机只能做 HTTP/SOCKS 代理，没法做透明
// 网关，因为 ConfigureNAT 在 Windows 上是 no-op（家用版 Windows 没 RRAS，
// ICS 强制 192.168.137/24 不好用）。指引里必须明说，不然用户设了"网关"
// 发现连不上网会一脸懵。
func DeviceGuide(status Status, mixedPort int) string {
	ip := status.LocalIP
	if ip == "" {
		ip = "<本机局域网 IP>"
	}
	router := firstNonEmpty(status.Router, "<路由器 IP>")

	var b strings.Builder
	b.WriteString(fmt.Sprintf("  本机 IP %s  ·  路由器 %s  ·  代理端口 %d (HTTP+SOCKS5)\n\n",
		ip, router, mixedPort))

	if runtime.GOOS == "windows" {
		b.WriteString("  ⚠ 当前是 Windows，只支持下面的方式 2（HTTP/SOCKS 代理）。\n")
		b.WriteString("    方式 1（改网关）需要 NAT / RRAS，家用版 Windows 做不了，Switch\n")
		b.WriteString("    / PS5 / Apple TV 这类不能设代理的设备请改用 Linux / macOS 或软路由。\n\n")

		b.WriteString(fmt.Sprintf("  📱 方式 2 · 填代理      iPhone / Android / 电脑 App / 浏览器插件\n"))
		b.WriteString(fmt.Sprintf("     代理服务器 → %s    端口 → %d    类型 → HTTP（或 SOCKS5）\n\n", ip, mixedPort))
		b.WriteString("     Android: Wi-Fi 设置 → 修改网络 → 高级选项 → 代理 = 手动\n")
		b.WriteString("     iOS:     Wi-Fi 设置 → 点 (i) → 配置代理 = 手动\n\n")

		b.WriteString("  💻 方式 3 · 本机自己用  跑 gateway 这台 Windows 自己也想走规则\n")
		b.WriteString("     TUN 已开，本机出向流量已经自动走 mihomo，不需要改浏览器代理设置\n")

		return b.String()
	}

	b.WriteString(fmt.Sprintf("  📺 方式 1 · 改网关      Switch / PS5 / Apple TV / 智能电视\n"))
	b.WriteString(fmt.Sprintf("     网关 → %s    DNS → %s    掩码 → 255.255.255.0\n\n", ip, ip))

	b.WriteString(fmt.Sprintf("  📱 方式 2 · 填代理      iPhone / 电脑 App / 浏览器插件\n"))
	b.WriteString(fmt.Sprintf("     代理服务器 → %s    端口 → %d\n\n", ip, mixedPort))

	b.WriteString("  💻 方式 3 · 本机自己用  跑 gateway 这台电脑的浏览器 / App 也想走规则\n")
	b.WriteString("     TUN 开着就自动生效；想关 TUN，把本机 DNS 改到 127.0.0.1 也行\n")
	b.WriteString("     macOS：菜单按 L 会自动把【所有活跃网卡】(Wi-Fi+Ethernet) DNS 都切到 127.0.0.1\n\n")

	b.WriteString("  💡 方式 1 需 TUN+DNS 都开；方式 2/3 只要 gateway 在跑\n")
	b.WriteString("  💡 验证方式 3：dscacheutil -q host -a name ping0.cc 返回 198.18.x.x 说明生效\n")

	return b.String()
}

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
