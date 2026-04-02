# Phone and Tablet Setup Guide

[Chinese version](../phone-setup.md)

After LAN Proxy Gateway is running on your computer, follow these steps to configure an iPhone, iPad, or Android phone or tablet.

Once configured, the device does not need any proxy app installed. This is especially convenient for family devices that should "just work".

## Requirements

- `sudo gateway start` is already running on the computer
- note the computer's LAN IP from the startup screen, for example `192.168.1.2`
- the phone or tablet is connected to the same Wi-Fi or router

---

## iOS (iPhone / iPad)

### 1. Open Wi-Fi settings

```text
Settings -> Wi-Fi
```

### 2. Open network details

Tap the info button next to the connected Wi-Fi network.

### 3. Configure IP

Change **Configure IP** to **Manual**:

| Field | Value | Notes |
|---|---|---|
| IP Address | `192.168.x.100` | An unused IP in the same subnet as the computer |
| Subnet Mask | `255.255.255.0` | Keep the default |
| Router | `192.168.x.2` | Your computer IP |

> Example: if the computer IP is `192.168.1.2`, you can set the iPhone to `192.168.1.100` and set Router to `192.168.1.2`.

### 4. Configure DNS

On the same page, change **Configure DNS** to **Manual**:

- remove existing DNS servers
- add your computer IP, such as `192.168.1.2`

### 5. Save

Tap **Save** to finish.

### Verify

Open Safari and visit [google.com](https://google.com). If it loads, the setup worked.

---

## Android Phone / Tablet

> Menu names vary by vendor, but the flow is usually similar across Huawei, Xiaomi, OPPO, Samsung, and others.

### 1. Open WLAN settings

```text
Settings -> WLAN (or Wi-Fi)
```

### 2. Modify the network

Long-press the connected Wi-Fi network and choose **Modify network**.  
On some phones, tap the settings icon next to the network name.

### 3. Expand advanced options

Turn on **Show advanced options** or open **Advanced**.

### 4. Switch to static IP

Change **IP settings** from DHCP to **Static**:

| Field | Value | Notes |
|---|---|---|
| IP address | `192.168.x.101` | An unused IP in the same subnet as the computer |
| Gateway | `192.168.x.2` | Your computer IP |
| Network prefix length | `24` | Same as `255.255.255.0` |
| DNS 1 | `192.168.x.2` | Your computer IP |
| DNS 2 | `8.8.8.8` | Optional |

> Example: if the computer IP is `192.168.1.2`, set the phone to `192.168.1.101`, and set both Gateway and DNS 1 to `192.168.1.2`.

### 5. Save

Tap **Save** to finish.

### Verify

Open a browser and visit [google.com](https://google.com). If it loads, the setup worked.

---

## Restore Automatic Networking

After you stop the gateway, switch the device back to automatic networking, or it may lose internet access.

**iOS**

```text
Settings -> Wi-Fi -> info button -> Configure IP: Automatic -> Configure DNS: Automatic -> Save
```

**Android**

```text
Settings -> WLAN -> long-press network -> Modify network -> change IP settings back to DHCP -> Save
```

The simplest method is often to forget the Wi-Fi network and reconnect. That restores automatic IP assignment.

---

## Common Issues

**Some apps say they detected a proxy**  
TUN mode is transparent proxying, so apps usually do not see it as a classic proxy. If a specific app still warns, it often does not affect normal use.

**Will mainland apps become slower?**  
No. The built-in smart routing keeps mainland traffic direct, so only overseas traffic goes through the proxy.

**Nothing works after setup, even domestic sites**  
1. Run `gateway status` to confirm the gateway is still running.  
2. Confirm both Gateway and DNS point to your computer IP.  
3. Confirm the first three IP segments match your computer, such as `192.168.1.x`.  
4. Turn Wi-Fi off and on again.

**What is the easiest way to set this up for family members?**  
Once configured, it keeps working as long as the gateway computer stays online. The main thing to remember is that if the computer is turned off, the phone must be switched back to automatic IP mode. Installing the system service with `sudo gateway service install` is recommended for always-on use.
