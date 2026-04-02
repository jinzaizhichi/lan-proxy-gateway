# Scenarios

## LAN Gateway Onboarding

Goal: let other LAN devices use this machine as gateway/DNS.

Suggested flow:

1. `gateway config`
2. Enable TUN in runtime settings
3. `gateway start`
4. In the runtime TUI, open device instructions or run `gateway status`

## AI Chains Mode Onboarding

Goal: Claude / ChatGPT / Codex / Cursor traffic uses residential egress.

Suggested flow:

1. `gateway chains`
2. Confirm `extension.mode = chains`
3. `gateway start`
4. In runtime TUI, use `/chains` and `/status` to verify ingress and egress

## Policy Group / Node Switching

Goal: change airport or strategy-group node from the CLI runtime.

Suggested flow:

1. `gateway start`
2. Open the runtime TUI
3. Use `/groups` or `Ctrl+P`
4. Choose the strategy group, then choose a node

## Local Machine Bypass

Goal: keep LAN sharing on, but avoid proxying the gateway machine’s own traffic.

Suggested flow:

1. `gateway config`
2. Open `局域网共享 / TUN / 端口`
3. Enable `本机绕过代理`
4. `gateway restart`
5. Validate with `gateway status`

## Health Check And Recovery

Goal: verify mihomo, TUN, API, and network egress.

Suggested flow:

1. `gateway status`
2. `gateway health`
3. `gateway start` if not running
4. `/logs` in the runtime TUI if the user needs recent runtime output

## Permission Preparation

Goal: make AI-driven terminal operation less blocked by privilege prompts.

Suggested flow:

1. Check whether root is required for the requested operation
2. Prefer project-supported permission setup or passwordless elevation path when available
3. Keep read-only commands non-invasive and run them without elevation
