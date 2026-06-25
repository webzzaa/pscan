package logging

/*
constants.go - 日志系统常量定义

统一管理common/logging包中的所有常量，便于查看和编辑。
*/

import (
	"time"

	"github.com/fatih/color"
)

// =============================================================================
// 日志级别常量 - 层级设计
// =============================================================================

// LogLevel 日志级别类型（数值越小越详细）
type LogLevel int

// 定义系统支持的日志级别常量（层级：Debug < Base < Info < Success < Vuln < Error）
const (
	LevelDebug   LogLevel = 0 // 调试信息（最详细）
	LevelBase    LogLevel = 1 // 基础信息（扫描进度等）
	LevelInfo    LogLevel = 2 // 一般信息（端口开放、服务识别等）
	LevelSuccess LogLevel = 3 // 成功结果（Web指纹等）
	LevelVuln    LogLevel = 4 // 重要发现（弱密码、漏洞等）
	LevelError   LogLevel = 5 // 错误信息（始终显示）
)

// 向后兼容的别名
const (
	LevelAll             LogLevel = LevelDebug // ALL 等同于 Debug（显示所有）
	LevelInfoSuccess     LogLevel = LevelInfo  // 废弃，映射到 Info
	LevelBaseInfoSuccess LogLevel = LevelBase  // 废弃，映射到 Base
)

// =============================================================================
// 时间显示常量 (从Formatter.go迁移)
// =============================================================================

const (
	// MaxMillisecondDisplay 毫秒显示的最大时长
	MaxMillisecondDisplay = time.Second
	// MaxSecondDisplay 秒显示的最大时长
	MaxSecondDisplay = time.Minute
	// MaxMinuteDisplay 分钟显示的最大时长
	MaxMinuteDisplay = time.Hour

	// SlowOutputDelay 慢速输出延迟
	SlowOutputDelay = 50 * time.Millisecond

	// ProgressClearDelay 进度条清除延迟
	ProgressClearDelay = 10 * time.Millisecond
)

// =============================================================================
// 日志前缀常量 (从Formatter.go迁移)
// =============================================================================

const (
	// PrefixDebug 调试日志前缀
	PrefixDebug = "[.]"
	// PrefixInfo 信息日志前缀
	PrefixInfo = "[*]"
	// PrefixSuccess 成功日志前缀
	PrefixSuccess = "[+]"
	// PrefixVuln 漏洞/重要发现前缀
	PrefixVuln = "[!]"
	// PrefixError 错误日志前缀
	PrefixError = "[-]"
)

// =============================================================================
// 默认配置常量
// =============================================================================

const (
	// DefaultLevel 默认日志级别
	DefaultLevel = LevelAll
	// DefaultEnableColor 默认启用彩色输出
	DefaultEnableColor = true
	// DefaultSlowOutput 默认不启用慢速输出
	DefaultSlowOutput = false
	// DefaultShowProgress 默认显示进度条
	DefaultShowProgress = true
)

// =============================================================================
// 默认颜色映射
// =============================================================================

// GetDefaultLevelColors 获取默认的日志级别颜色映射
func GetDefaultLevelColors() map[LogLevel]interface{} {
	return map[LogLevel]interface{}{
		LevelError:   color.FgYellow, // 错误日志显示黄色
		LevelVuln:    color.FgRed,    // 漏洞/重要发现显示红色（密码成功、漏洞等）
		LevelBase:    color.FgWhite,  // 基础日志显示白色（普通信息）
		LevelInfo:    color.FgWhite,  // 信息日志显示白色（普通信息）
		LevelSuccess: color.FgGreen,  // 成功日志显示绿色（Web指纹等）
		LevelDebug:   color.FgWhite,  // 调试日志显示白色
	}
}
