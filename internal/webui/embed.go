// Package webui 是 lan-proxy-gateway 自己的 HTTP 控制台。
//
// 跟 mihomo 自己挂在 19090/ui 的那个面板不是一回事 —— 那个面板只能调 mihomo
// API（选节点、看流量），改不了 gateway 的配置（代理源、网关模式、规则集合
// 这些）。本包提供完整功能的 Web 控制台，监听 webui 端口（默认 19091），
// 直接调用 internal/app.App 来读写 gateway.yaml 并 reload mihomo。
//
// 设计原则：
//   - 单页 HTML + 原生 JS，无构建链，避免引入 npm / webpack
//   - 静态资源通过 go:embed 烧进二进制，部署零外部依赖
//   - REST API 严格走 /api/* 前缀，单页应用 fetch 它们
//   - 端口可在 gateway.yaml runtime.ports.web_ui 改，0 = 不监听
package webui

import "embed"

//go:embed static/*
var staticFS embed.FS
