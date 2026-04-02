# PS5 Setup Guide

[Chinese version](../ps5-setup.md)

After LAN Proxy Gateway is running on your computer, follow these steps to configure your PS5.

## Requirements

- `sudo gateway start` is already running on the computer
- note the computer's LAN IP from the startup screen, for example `192.168.1.2`
- the PS5 and the computer are connected to the same Wi-Fi or router

## Steps

### 1. Open network settings

```text
Settings -> Network -> Set Up Internet Connection
```

### 2. Choose Wi-Fi and set it up manually

Select your Wi-Fi network. Instead of connecting right away, choose **Set Up Manually**.  
If the PS5 is already connected, press the **Options** button on the controller and open the settings for that network.

### 3. Configure IP address

Change **IP Address Settings** to **Manual**:

| Field | Value | Notes |
|---|---|---|
| IP Address | `192.168.x.103` | An unused IP in the same subnet as the computer |
| Subnet Mask | `255.255.255.0` | Keep the default |
| Default Gateway | `192.168.x.2` | Your computer IP |

> Example: if the computer IP is `192.168.1.2`, set the PS5 to `192.168.1.103` and Default Gateway to `192.168.1.2`.

### 4. Configure DNS

Change **DNS Settings** to **Manual**:

| Field | Value |
|---|---|
| Primary DNS | Your computer IP, for example `192.168.1.2` |
| Secondary DNS | `8.8.8.8` |

### 5. Other options

| Field | Value |
|---|---|
| MTU Settings | Automatic |
| Proxy Server | Do Not Use |

> Important: **Proxy Server must stay on "Do Not Use"**. Traffic is already being forwarded transparently through the gateway. Setting a proxy server on the PS5 will conflict with that path.

### 6. Test the connection

After saving, choose **Test Internet Connection**.

You should see:

- obtain IP address: successful
- internet connection: successful
- PlayStation Network sign-in: successful

## About NAT Type

When using the gateway, PS5 usually shows **Type 2**, and sometimes **Type 3**.

| NAT Type | Meaning | Impact |
|---|---|---|
| Type 1 | Direct connection, no gateway | Not the target setup here |
| Type 2 | Moderate, behind the gateway | Most online games work fine |
| Type 3 | Strict | Voice chat or some P2P games may be limited |

If you end up with Type 3 and it affects gameplay, try switching nodes in the panel at `http://<your-computer-ip>:9090/ui`.

## Restore Automatic Networking

After stopping the gateway, switch the PS5 back to automatic settings:

```text
Settings -> Network -> Set Up Internet Connection -> choose your Wi-Fi
-> IP Address Settings: Automatic
-> DNS Settings: Automatic
-> Proxy Server: Do Not Use
```

Or delete the Wi-Fi network and reconnect.

## Common Issues

**NAT Type 3 and online play is affected**  
This can happen when routing through a gateway. Try switching to a different node in the panel. Some nodes produce better NAT behavior than others.

**PSN Store is slow or does not open**  
1. Run `gateway status` to confirm the gateway is running.  
2. Switch to a lower-latency node in the panel.  
3. Some nodes are simply not ideal for PSN, so try a few.

**Game downloads are slow**  
Download speed depends on node bandwidth. Test nodes in the panel and switch to a faster one. If you only need game downloads and do not need overseas access at that moment, you can temporarily switch the PS5 back to automatic networking so it connects directly to the router.

**Cannot sign in to PlayStation Network**  
1. Confirm the gateway is running with `gateway status`.  
2. Confirm the PS5 Default Gateway and Primary DNS both point to your computer IP.  
3. Confirm the first three IP segments match your computer, such as `192.168.1.x`.  
4. Confirm Proxy Server is set to **Do Not Use**.
