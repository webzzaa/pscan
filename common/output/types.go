package output

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// ScanResult 扫描结果结构
type ScanResult struct {
	Time    time.Time              `json:"time"`    // 发现时间
	Type    ResultType             `json:"type"`    // 结果类型
	Target  string                 `json:"target"`  // 目标(IP/域名/URL)
	Status  string                 `json:"status"`  // 状态描述
	Details map[string]interface{} `json:"details"` // 详细信息
}

// FormatDetails 格式化Details为键值对字符串（排序key以保证输出稳定）
func (r *ScanResult) FormatDetails(separator, kvFormat string) string {
	if len(r.Details) == 0 {
		return ""
	}

	keys := make([]string, 0, len(r.Details))
	for key := range r.Details {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	pairs := make([]string, 0, len(keys))
	for _, key := range keys {
		pairs = append(pairs, fmt.Sprintf(kvFormat, key, r.Details[key]))
	}
	return strings.Join(pairs, separator)
}

// Writer 输出写入器接口
type Writer interface {
	Write(result *ScanResult) error
	WriteHeader() error
	Flush() error
	Close() error
	GetFormat() Format
}

// ManagerConfig 输出管理器配置
type ManagerConfig struct {
	OutputPath string `json:"output_path"` // 输出路径
	Format     Format `json:"format"`      // 输出格式
}

// DefaultManagerConfig 默认管理器配置
func DefaultManagerConfig(outputPath string, format Format) *ManagerConfig {
	return &ManagerConfig{
		OutputPath: outputPath,
		Format:     format,
	}
}
