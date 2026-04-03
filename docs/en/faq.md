# FAQ

[Chinese version](../faq.md)

**Q: Everything works, but overseas sites are still slow. Why?**  
> It is usually a node issue. Open the panel at `http://<your-computer-ip>:9090/ui`, test latency, and switch to a lower-latency node.

**Q: Will mainland apps become slower?**  
> No. The built-in rules keep mainland traffic direct, proxy overseas traffic, and block obvious ads. Domestic sites should feel the same as before.

**Q: NAT type on PS5 or Switch got worse. Will online play break?**  
> Most games still work fine. If a specific game has trouble, switch nodes in the panel, or temporarily set the console back to automatic networking so it connects directly to the router.

**Q: Netflix on Apple TV shows the wrong region. What should I do?**  
> Netflix only works on nodes that support streaming unlock. In `http://<your-computer-ip>:9090/ui`, switch to a node tagged with something like "unlock", "streaming", or "Netflix".

**Q: The device has no network at all, not even for domestic sites.**  
> Check these in order: 1. run `gateway status` to confirm the gateway is running; 2. confirm both gateway and DNS on the device point to your computer IP; 3. confirm the first three IP segments match your computer, such as `192.168.1.x`; 4. restart once with `sudo gateway stop` and `sudo gateway start`.

**Q: Does it work on Windows?**  
> Yes. Windows is fully supported. IP forwarding is enabled through `netsh`, status checks and default-interface detection are handled in a Windows-aware way, and auto-start is installed through Task Scheduler. Run PowerShell or Command Prompt as Administrator, which is the Windows equivalent of `sudo`.

**Q: Can I run it on a router?**  
> It depends on the router system and hardware.
>
> **Possible:**
> - high-performance routers running full Linux, such as x86 or ARM OpenWrt / iKuai boxes with at least 256MB RAM
> - Synology or QNAP NAS devices
> - GL.iNet routers that run OpenWrt and allow SSH access
>
> **How to do it:** SSH into the device, download the matching `gateway-linux-arm64` or `gateway-linux-amd64`, place it under `/usr/local/bin/`, and use it like a normal Linux install.
>
> **Usually not practical:**
> - common consumer routers from Huawei, Xiaomi, TP-Link, and similar vendors, because the system is heavily customized and often lacks the pieces required to run this
> - routers with less than 128MB RAM, because `mihomo` needs a basic amount of memory
>
> **Recommended approach:** if your router supports it, that is the cleanest whole-home setup. If not, a Mac mini, Raspberry Pi, NAS, or low-power mini PC works very well.

**Q: Do I need a Mac?**  
> No. macOS, Linux, and Windows are all supported. A Mac mini is simply a convenient always-on option because of its low power usage.

**Q: Why does it need `sudo`? Is that safe?**  
> `sudo` means administrator privileges on macOS and Linux. TUN mode needs to create a virtual network interface and modify routing tables, which are system-level operations. The project is open source, so you can inspect the code on GitHub.

**Q: Is my subscription URL uploaded anywhere?**  
> No. The subscription URL is stored only in your local `gateway.yaml`. It is not uploaded anywhere, and `.gitignore` already excludes that file.

**Q: Compared with a soft router, what are the tradeoffs?**  
> | | LAN Proxy Gateway | Soft router |
> |---|---|---|
> | Cost | Reuse an existing computer | Need extra hardware |
> | Learning curve | A few commands with guides | Usually much higher |
> | Stability | Depends on the computer staying on | Better for 24/7 use |
> | Best fit | People with an available computer | People who want a dedicated always-on box |

**Q: How do other devices recover normal internet after I stop the gateway?**  
> Set networking back to automatic:
>
> | Device | Action |
> |---|---|
> | iPhone | Settings -> Wi-Fi -> info button -> Configure IP: Automatic; Configure DNS: Automatic |
> | Android | Settings -> WLAN -> long-press network -> Modify -> change IP settings back to DHCP |
> | Switch | Settings -> Internet -> Internet Settings -> Wi-Fi -> Change Settings -> set everything back to automatic |
> | PS5 | Settings -> Network -> Set Up Internet Connection -> reconnect to Wi-Fi |
> | Apple TV | Settings -> Network -> Wi-Fi -> Configure IP: Automatic; Configure DNS: Automatic |
> | Smart TV | Settings -> Network -> Wi-Fi -> set IP back to DHCP / Automatic |
>
> The easiest method is often to forget the Wi-Fi network and reconnect. That usually restores DHCP automatically.
