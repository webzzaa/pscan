//go:build !debug
// +build !debug

package debug

// 生产版本：pprof 完全不编译进来

func Start() {}
func Stop()  {}
