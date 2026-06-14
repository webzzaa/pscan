//go:build !web

package web

import "errors"

// ErrWebNotSupported 非Web版本不支持Web功能
var ErrWebNotSupported = errors.New("web mode not supported in this build, rebuild with: go build -tags web")

// StartServer 非Web版本的空实现
func StartServer(port int) error {
	return ErrWebNotSupported
}
