# Nintendo Switch Setup Guide

[Chinese version](../switch-setup.md)

After LAN Proxy Gateway is running on your computer, follow these steps to configure your Switch.

## Requirements

- `sudo gateway start` is already running on the computer
- note the computer's LAN IP from the startup screen, for example `192.168.1.2`
- the Switch and the computer are connected to the same Wi-Fi or router

## Steps

### 1. Open network settings

```text
Home -> System Settings -> Internet -> Internet Settings
```

### 2. Choose your Wi-Fi

Select your Wi-Fi network, then choose **Change Settings**.

### 3. Configure IP address

Change **IP Address Settings** to **Manual**:

| Field | Value | Notes |
|---|---|---|
| IP Address | `192.168.x.100` | An unused IP in the same subnet as the computer |
| Subnet Mask | `255.255.255.0` | Keep the default |
| Gateway | `192.168.x.2` | Your computer IP |

> Example: if the computer IP is `192.168.1.2`, set the Switch to `192.168.1.100` and Gateway to `192.168.1.2`.

### 4. Leave DNS on Automatic

Keep **DNS Settings** on **Automatic**.

> On Switch, manually changing DNS can sometimes break connectivity for YouTube or eShop. In practice, changing only the gateway works better.

### 5. Save and test

Choose **Save**, then **Connect to This Network**.

The connection test should show that the console is connected to the internet.

## Verify

- eShop opens normally
- online games connect
- the YouTube app can play videos

## About NAT Type

Behind the gateway, Switch often shows NAT type B or C. That is normal, and most online games still work fine. If one specific game has trouble, switch nodes in the panel at `http://<your-computer-ip>:9090/ui`.

## Restore Automatic Networking

After stopping the gateway, change the Switch settings back:

```text
System Settings -> Internet -> Internet Settings -> your Wi-Fi -> Change Settings
-> IP Address Settings: Automatic
-> DNS Settings: Automatic
```

Or simply delete the Wi-Fi network and reconnect. That restores DHCP automatically.

## Common Issues

**Connection test fails**
1. Confirm `gateway status` shows `mihomo` running.
2. Confirm the Switch gateway points to your computer IP.
3. Confirm the first three IP segments match your computer, such as `192.168.1.x`.

**General internet works, but eShop does not open**
Your rules or node may need an update. Try `sudo gateway update`, or switch nodes in the panel.

**NAT type shows C or D**
That can happen in TUN mode and is usually expected. Most games still work. If a specific game is affected, try a different node in the panel.
