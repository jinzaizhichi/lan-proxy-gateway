// Package embed exposes the baked-in mihomo skeleton template and rule data.
package embed

import (
	"embed"
	_ "embed"
)

//go:embed template.yaml
var Template string

// WebUI 内嵌一份轻量 mihomo 控制台（单 index.html，样式/脚本都在里面）。
// engine 启动时把它释放到 workdir/ui/，mihomo 的 external-ui: ui 直接就能 serve，
// 浏览器访问 http://<host>:<api_port>/ui 就能切节点、切模式、看流量。
//
//go:embed webui
var WebUI embed.FS
