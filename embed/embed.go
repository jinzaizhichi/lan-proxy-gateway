// Package embed exposes the baked-in mihomo skeleton template and rule data.
package embed

import (
	"embed"
	_ "embed"
)

//go:embed template.yaml
var Template string

// WebUI 内嵌 metacubexd 控制台。engine 启动时把它释放到 workdir/ui/，
// mihomo 的 external-ui: ui 直接就能 serve，浏览器访问
// http://<host>:<api_port>/ui 就能切节点、切模式、看流量。
//
// 必须用 all: 前缀 —— 默认的 //go:embed webui 会跳过任何以 `.` 或 `_`
// 开头的文件/目录，而 metacubexd 的全部静态资源都在 _nuxt/ 和 _fonts/
// 下，不加 all: 就只会打包几个根目录文件，浏览器打开 /ui/ 拿到壳 HTML
// 但所有 _nuxt/*.js 返回 404，页面全白。
//
//go:embed all:webui
var WebUI embed.FS
