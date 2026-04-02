# Apple TV Setup Guide

[Chinese version](../appletv-setup.md)

After LAN Proxy Gateway is running on your computer, follow these steps to configure your Apple TV.

## Requirements

- `sudo gateway start` is already running on the computer
- note the computer's LAN IP from the startup screen, for example `192.168.1.2`
- the Apple TV and the computer are connected to the same Wi-Fi or router

## Steps

### 1. Open network settings

```text
Settings -> Network -> Wi-Fi
```

### 2. Select the Wi-Fi network

Open the currently connected Wi-Fi network.

### 3. Configure IP

Change **Configure IP** to **Manual**:

| Field | Value | Notes |
|---|---|---|
| IP Address | `192.168.x.104` | An unused IP in the same subnet as the computer |
| Subnet Mask | `255.255.255.0` | Keep the default |
| Router | `192.168.x.2` | Your computer IP |

> Example: if the computer IP is `192.168.1.2`, set Apple TV to `192.168.1.104` and Router to `192.168.1.2`.

### 4. Configure DNS

Change **Configure DNS** to **Manual**:

- set DNS to your computer IP, such as `192.168.1.2`
- edit the existing entry or delete and add it again

### 5. Finish

Go back and Apple TV will apply the settings automatically.

## Verify

- open YouTube and confirm videos load
- open Netflix and confirm browsing and playback work
- test Disney+ or other streaming apps if needed

## About Netflix

Netflix is strict about proxy detection. Not every node will work. If Netflix shows a region error or refuses playback:

1. open the panel at `http://<your-computer-ip>:9090/ui`
2. look for nodes tagged with terms like "unlock", "streaming", or "Netflix"
3. switch nodes and try again

Unlock support depends on your provider, so you may need to test several nodes.

## About 4K Playback

4K streaming needs stable bandwidth, usually around 25Mbps or higher. If playback buffers or stalls, switch to a faster node in the panel.

## Restore Automatic Networking

After stopping the gateway, switch Apple TV back:

```text
Settings -> Network -> Wi-Fi -> choose your network
-> Configure IP: Automatic
-> Configure DNS: Automatic
```

## Common Issues

**Netflix says a proxy was detected or shows the wrong region**  
Switch to a node that supports streaming unlock in `http://<your-computer-ip>:9090/ui`.

**4K playback buffers**  
The current node likely does not have enough bandwidth. Switch to a faster node.

**No internet connection**  
1. Run `gateway status` to confirm the gateway is running.  
2. Confirm both Router and DNS on Apple TV point to your computer IP.  
3. Confirm the first three IP segments match your computer, such as `192.168.1.x`.  
4. Restart Wi-Fi on Apple TV.
