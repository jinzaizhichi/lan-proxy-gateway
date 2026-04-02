# Changelog

All notable changes to this project will be documented here.

## v2.2.0 - 2026-04-03

### Added

- Added a full-screen runtime TUI console after `gateway start`
- Added slash-command interaction and in-terminal strategy group / node switching
- Added `gateway skill` for AI-client skill discovery
- Added `gateway permission` to help configure passwordless control
- Added `runtime.tun.bypass_local` so the host machine can stay closer to direct access while LAN devices still use the gateway
- Added cached update notices on the CLI home page, start summary, and TUI console

### Changed

- Unified config structure into `proxy`, `runtime`, `rules`, and `extension`
- Kept backward compatibility for legacy flat config fields
- Moved rules into a structured rules system instead of hardcoded template blocks
- Improved chains / egress reporting to show entry node, normal exit, and residential exit more clearly
- Improved config-center summaries and config-path visibility
- Improved install and usage messaging for mainland China users

### Documentation

- Rewrote the main Chinese README
- Added `README_EN.md`
- Expanded command and advanced guides
- Added release notes and release-process documentation
