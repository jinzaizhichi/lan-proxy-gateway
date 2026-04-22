//go:build windows

package platform

import (
	"net"
	"testing"
)

func TestIsVirtualAdapterName(t *testing.T) {
	virtual := []string{
		"Mihomo",
		"mihomo",
		"Wintun Userspace Tunnel",
		"WINTUN",
		"Tailscale",
		"TAP-Windows Adapter V9",
		"tap 0901",
	}
	for _, n := range virtual {
		if !isVirtualAdapterName(n) {
			t.Errorf("expected %q to be classified as virtual", n)
		}
	}

	physical := []string{
		"WLAN",
		"Wi-Fi",
		"Ethernet",
		"以太网 2",
		"Local Area Connection",
	}
	for _, n := range physical {
		if isVirtualAdapterName(n) {
			t.Errorf("expected %q to NOT be classified as virtual", n)
		}
	}
}

func TestIsTunIPv4(t *testing.T) {
	cgnat := []string{
		"198.18.0.1",
		"198.18.255.255",
		"198.19.0.1",
		"198.19.200.50",
	}
	for _, s := range cgnat {
		if !isTunIPv4(net.ParseIP(s)) {
			t.Errorf("expected %s to be classified as TUN IP (198.18/15)", s)
		}
	}

	real := []string{
		"192.168.12.109",
		"10.0.0.5",
		"172.16.1.1",
		"198.17.255.255", // just outside the block
		"198.20.0.1",     // just outside the block
	}
	for _, s := range real {
		if isTunIPv4(net.ParseIP(s)) {
			t.Errorf("expected %s to NOT be classified as TUN IP", s)
		}
	}
}
