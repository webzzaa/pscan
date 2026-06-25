//go:build !web

package common

// WebMode 非Web版本永远为false
var WebMode = false

// WebPort 非Web版本不使用
var WebPort = 0
