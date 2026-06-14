package output

import (
	"os"
)

// =============================================================================
// 输出格式常量
// =============================================================================

// Format 输出格式类型
type Format string

const (
	// FormatTXT 文本格式输出
	FormatTXT Format = "txt"
	// FormatJSON JSON格式输出
	FormatJSON Format = "json"
	// FormatCSV CSV格式输出
	FormatCSV Format = "csv"
)

// =============================================================================
// 结果类型常量
// =============================================================================

// ResultType 定义结果类型
type ResultType string

const (
	// TypeHost 主机存活
	TypeHost ResultType = "HOST"
	// TypePort 端口开放
	TypePort ResultType = "PORT"
	// TypeService 服务识别
	TypeService ResultType = "SERVICE"
	// TypeVuln 漏洞发现
	TypeVuln ResultType = "VULN"
)

// =============================================================================
// 文件操作常量
// =============================================================================

const (
	// DefaultFilePermissions 文件操作权限
	DefaultFilePermissions = 0644
	// DefaultDirPermissions 目录操作权限
	DefaultDirPermissions = 0755

	// DefaultFileFlags 文件打开标志
	DefaultFileFlags = os.O_CREATE | os.O_WRONLY | os.O_APPEND

	// JSONIndentPrefix JSON格式化前缀
	JSONIndentPrefix = ""
	// JSONIndentString JSON格式化缩进字符串
	JSONIndentString = "  "
)
