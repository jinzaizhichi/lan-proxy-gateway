package gateway

import (
	"fmt"
	"strings"
)

// DeviceGuide 返回可贴到设备上的接入说明。
// mixedPort 是 mihomo 的 HTTP+SOCKS5 混合端口，提供给方式 2（局域网代理）用。
func DeviceGuide(status Status, mixedPort int) string {
	ip := status.LocalIP
	if ip == "" {
		ip = "<本机局域网 IP>"
	}
	router := firstNonEmpty(status.Router, "<路由器 IP>")

	var b strings.Builder
	b.WriteString("=== 设备接入指引 ===\n")
	b.WriteString(fmt.Sprintf("  本机局域网 IP : %s\n", ip))
	b.WriteString(fmt.Sprintf("  路由器 IP     : %s\n", router))
	b.WriteString("\n下面两种方式任选，看设备本身支持啥：\n")

	// 方式 1：TUN 网关
	b.WriteString("\n─────────────── 方式 1：TUN 网关（覆盖面最广） ───────────────\n")
	b.WriteString("  适用：Switch / PS5 / Apple TV / 智能电视 等「只能改网关」的设备。\n")
	b.WriteString("  在设备网络设置里填：\n")
	b.WriteString(fmt.Sprintf("    网关 (Gateway)  →  %s\n", ip))
	b.WriteString(fmt.Sprintf("    DNS 服务器       →  %s\n", ip))
	b.WriteString("    子网掩码          →  255.255.255.0      （99% 家庭网就这个）\n")
	b.WriteString("    IP 地址           →  保留 DHCP 或设同网段静态 IP\n")
	b.WriteString("  保存并重连 Wi-Fi 后：所有流量（含 YouTube / 游戏 / 各类 App）自动走代理。\n")
	b.WriteString("  前提：本机 TUN 要开、DNS 代理要开。\n")

	// 方式 2：局域网 HTTP/SOCKS5 代理
	b.WriteString("\n─────────────── 方式 2：局域网代理（App 级粒度） ───────────────\n")
	b.WriteString("  适用：iPhone / 电脑 App / 浏览器插件 —— 这些能手动填代理服务器的场景。\n")
	b.WriteString("  在设备 / App 的代理设置里填：\n")
	b.WriteString(fmt.Sprintf("    代理服务器地址  →  %s\n", ip))
	b.WriteString(fmt.Sprintf("    端口            →  %d        （HTTP + SOCKS5 混合，任选）\n", mixedPort))
	b.WriteString("    不用填用户名 / 密码\n")
	b.WriteString("  只走**那个 App 自己的代理流量**，不支持代理的流量不会被劫持。\n")
	b.WriteString("  例：Switch 走方式 2 的话打不开 YouTube（Switch 不会主动把流量塞给代理）。\n")
	b.WriteString("  不需要 TUN，只要 gateway 在跑就能用。\n")
	return b.String()
}

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
