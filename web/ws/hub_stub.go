//go:build !web

package ws

// Hub 非Web版本的空结构
type Hub struct{}

// NewHub 非Web版本返回nil
func NewHub() *Hub {
	return nil
}
