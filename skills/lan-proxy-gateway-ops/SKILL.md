---
name: lan-proxy-gateway-ops
description: Use when operating lan-proxy-gateway for LAN sharing, chains mode, policy-group switching, local-bypass mode, health checks, or scenario-based setup from an AI client. This skill maps user intents to the right gateway CLI flows and tells the agent when to use status, config, chains, groups, logs, and permission-related commands.
---

# LAN Proxy Gateway Ops

Use this skill when the user wants to operate `lan-proxy-gateway` through conversation instead of manually remembering CLI commands.

## Quick Start

1. Start with `gateway status` or `gateway config show` to understand the current state.
2. If the user is setting up from scratch, prefer `gateway config` and `gateway chains`.
3. If the user is already running the gateway, prefer the runtime TUI after `gateway start`.
4. If the user wants to switch nodes or policy groups, use the TUI group picker or inspect proxy groups through the mihomo API-backed flows.
5. If the user mentions “本机不要走代理”, use the local-bypass scenario.

## Scenario Routing

Read [references/scenarios.md](references/scenarios.md) and choose the closest scenario:

- LAN gateway onboarding
- AI chains mode onboarding
- Policy-group / node switching
- Local-machine bypass while keeping LAN sharing
- Health check and log-based recovery
- Permission / non-root control preparation

## Command Preferences

- Prefer plain `gateway <command>` first.
- If a command requires elevated privileges, explain that the project may need root for TUN, IP forwarding, or firewall changes.
- If the environment already supports passwordless elevation, plain `gateway start` / `gateway stop` / `gateway restart` should be preferred over teaching the user a long workaround.
- For runtime operation, prefer the full-screen TUI after `gateway start`; it is the primary control surface.

## High-Value Flows

- `gateway start`
  Use for the main runtime console, slash commands, policy-group switching, logs, and device onboarding.

- `gateway config`
  Use when the user wants guided configuration without editing YAML.

- `gateway chains`
  Use for residential chain setup and AI-safe outbound routing.

- `gateway status`
  Use to verify gateway health, current node, airport ingress, and residential egress.

## Do Not Assume

- Do not assume the current machine should proxy its own traffic. Check whether `runtime.tun.bypass_local` is enabled.
- Do not assume chains mode is enabled just because a residential proxy is configured. Check `extension.mode`.
- Do not assume the user wants to edit YAML directly; prefer guided CLI flows first.
