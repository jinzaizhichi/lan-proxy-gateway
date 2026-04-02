# LAN Proxy Gateway

[中文说明](README.md)

An open-source educational project that turns your computer into a LAN-wide transparent proxy gateway powered by `mihomo`.

It is designed first for Chinese beginners, but the underlying networking ideas are universal:

- Share proxy capability with devices that cannot install proxy apps
- Use one computer as a transparent gateway for your home LAN
- Add `chains` mode for AI tools that benefit from a residential exit
- Operate everything through an interactive CLI / TUI

## Core Features

- LAN gateway sharing for `Switch / PS5 / Apple TV / smart TV / phones`
- Subscription URL or local Clash / mihomo config as proxy source
- `chains` mode for `Claude / ChatGPT / Codex / Cursor`
- Interactive config center: `gateway config`
- Runtime TUI console after `gateway start`
- Strategy group / node switching inside the terminal
- Optional local-host bypass while LAN clients still use the gateway
- AI client skill entry: `gateway skill`
- Update prompt plus `gateway update`

## Quick Start

### Install

For users in mainland China, the jsDelivr entry is usually friendlier than GitHub raw:

#### macOS / Linux

```bash
curl -fsSL https://cdn.jsdelivr.net/gh/Tght1211/lan-proxy-gateway@main/install.sh | bash
```

Fallback:

```bash
curl -fsSL https://raw.githubusercontent.com/Tght1211/lan-proxy-gateway/main/install.sh | bash
```

#### Windows PowerShell

```powershell
irm https://cdn.jsdelivr.net/gh/Tght1211/lan-proxy-gateway@main/install.ps1 | iex
```

Fallback:

```powershell
irm https://raw.githubusercontent.com/Tght1211/lan-proxy-gateway/main/install.ps1 | iex
```

### Initialize

```bash
gateway install
```

### Start

```bash
sudo gateway start
```

### Connect LAN devices

Set the other device's:

- Gateway to your computer's LAN IP
- DNS to the same IP

Device guides:

- [Phone / Tablet](docs/phone-setup.md)
- [Nintendo Switch](docs/switch-setup.md)
- [PS5](docs/ps5-setup.md)
- [Apple TV](docs/appletv-setup.md)
- [Smart TV](docs/tv-setup.md)

## Most Useful Commands

| Command | Purpose |
|---|---|
| `gateway install` | Initial setup wizard |
| `gateway config` | Interactive config center |
| `sudo gateway start` | Start gateway and open runtime TUI |
| `gateway status` | Show runtime and egress status |
| `gateway chains` | Chains / residential proxy wizard |
| `gateway switch` | Switch proxy source or extension mode |
| `gateway skill` | Show AI skill path and scenarios |
| `gateway permission install` | Install passwordless control rule |
| `sudo gateway update` | Upgrade to the latest version |

## Configuration Layout

The config is grouped into four stable sections:

```yaml
proxy:
runtime:
rules:
extension:
```

This keeps the project's two main strengths intact:

1. LAN-wide gateway sharing
2. Chains / residential exit for AI traffic

Legacy flat config fields remain backward compatible.

## Documentation

- [Command Reference](docs/commands.md)
- [Advanced Guide](docs/advanced.md)
- [FAQ](docs/faq.md)

## Releases

Releases include:

- raw binaries
- compressed archives
- `SHA256SUMS`
- upgrade notes

Download: [GitHub Releases](https://github.com/Tght1211/lan-proxy-gateway/releases)

## License

[MIT](LICENSE)
