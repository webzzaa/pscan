//go:build windows

package common

import (
	"os"

	"golang.org/x/sys/windows"
)

// init 在包加载时自动启用 Windows 终端的 ANSI 支持
func init() {
	enableVirtualTerminalProcessing()
}

// enableVirtualTerminalProcessing 启用 Windows 控制台的虚拟终端处理
// 使 Windows 终端支持 ANSI 转义码（如 \r, \033[2K 等）
func enableVirtualTerminalProcessing() {
	handle := windows.Handle(os.Stdout.Fd())

	var mode uint32
	_ = windows.GetConsoleMode(handle, &mode)
	_ = windows.SetConsoleMode(handle, mode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
}
