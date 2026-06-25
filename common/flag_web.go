//go:build web

package common

import "flag"

// WebMode 表示是否启动Web管理界面
var WebMode bool

// WebPort Web服务器端口
var WebPort int

func init() {
	flag.BoolVar(&WebMode, "web", false, "启动Web管理界面 (Start Web UI)")
	flag.IntVar(&WebPort, "webport", 10240, "Web服务器端口 (Web server port)")
}
