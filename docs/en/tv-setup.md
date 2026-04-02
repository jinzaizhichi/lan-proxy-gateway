# Smart TV Setup Guide

[Chinese version](../tv-setup.md)

Most smart TVs, including Xiaomi, Redmi, Hisense, TCL, Sony Android TV, and Philips, use an Android-like network flow. The setup is similar to an Android phone.

## Requirements

- `sudo gateway start` is already running on the computer
- note the computer's LAN IP from the startup screen, for example `192.168.1.2`
- the TV and the computer are connected to the same Wi-Fi or router

## Steps

### 1. Open network settings

Menu paths vary by brand. Common paths include:

```text
Settings -> Network -> Wireless Network (or Wi-Fi)
```

Sometimes it is under **System Settings** or **Advanced Settings**.

### 2. Open Wi-Fi details

Select the connected Wi-Fi network and look for **Advanced Settings** or **Modify Network**.

On some TVs, you need to long-press the OK button on the remote or press the menu button to reveal more options.

### 3. Switch to static IP

Change IP assignment from **Automatic (DHCP)** to **Manual** or **Static**, then enter:

| Field | Value | Notes |
|---|---|---|
| IP Address | `192.168.x.105` | An unused IP in the same subnet as the computer |
| Subnet Mask | `255.255.255.0` | Keep the default |
| Gateway | `192.168.x.2` | Your computer IP |
| DNS 1 | `192.168.x.2` | Your computer IP |
| DNS 2 | `8.8.8.8` | Optional |

> Example: if the computer IP is `192.168.1.2`, set the TV to `192.168.1.105`, and set both Gateway and DNS 1 to `192.168.1.2`.

### 4. Save and reconnect

Save the settings, then reconnect Wi-Fi so they take effect.

## Verify

Open YouTube or another streaming app. If playback works, the setup is done.

---

## Brand-Specific Notes

### Xiaomi / Redmi TV

```text
Settings -> Network -> Wireless Network -> choose connected Wi-Fi -> Advanced Settings -> set IP to Static
```

If you cannot find **Advanced Settings**, try:

```text
Settings -> More Settings -> Network -> Wi-Fi -> long-press the current network
```

### Hisense / TCL TV

```text
Menu -> Settings -> Network Connection -> Wireless Network -> select Wi-Fi -> press the settings button on the remote -> Advanced Settings
```

### Sony Android TV

```text
Settings -> Network & Internet -> Wi-Fi -> choose connected network -> edit icon
-> expand Advanced Options -> change IP settings to Static
```

### Philips Android TV

```text
Settings -> Wireless and Wired Networks -> Connect to Network -> choose your Wi-Fi
-> Advanced Settings -> Static IP
```

### Samsung Smart TV (Tizen)

Samsung uses a different UI:

```text
Settings -> General -> Network -> Network Status -> IP Settings
-> IP Settings: Manual
-> DNS Settings: Manual
```

Use the same values shown above.

---

## What If You Cannot Find Static IP Settings?

Some older TVs or heavily customized systems hide this menu. You can try:

1. search the web for `your-tv-model manual IP`
2. bind a fixed IP on your router with DHCP reservation, then use that fixed address in your TV setup

Some routers can also push custom gateway and DNS values through DHCP reservations.

---

## Restore Automatic Networking

After stopping the gateway, switch the TV back to DHCP or Automatic:

```text
Settings -> Network -> Wi-Fi -> choose current network -> Advanced Settings -> IP settings back to DHCP / Automatic
```

Or simply forget the Wi-Fi network and reconnect.

---

## Common Issues

**I cannot find advanced network settings**  
TV menus vary a lot. A common trick is to long-press the connected Wi-Fi network or press the menu button on the remote. If that still fails, search the web for your exact TV model plus `manual IP`.

**I entered everything, but videos still do not play**  
1. Confirm the gateway is running with `gateway status`.  
2. Confirm Gateway and DNS point to your computer IP, not the router IP.  
3. Open the panel at `http://<your-computer-ip>:9090/ui` and switch to a faster node.

**YouTube works, but Netflix does not**  
Netflix needs nodes that support streaming unlock. Switch to a node tagged with something like "Netflix", "streaming", or "unlock".
