package output

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Manager 简化的输出管理器
type Manager struct {
	mu     sync.RWMutex
	config *ManagerConfig
	writer Writer
	closed bool
}

// NewManager 创建新的输出管理器
func NewManager(config *ManagerConfig) (*Manager, error) {
	if config == nil {
		return nil, fmt.Errorf("output config cannot be nil")
	}

	// 创建输出目录
	if err := createOutputDir(config.OutputPath); err != nil {
		return nil, err
	}

	manager := &Manager{
		config: config,
	}

	// 初始化写入器（内部会验证格式）
	if err := manager.initializeWriter(); err != nil {
		return nil, err
	}

	return manager, nil
}

// createOutputDir 创建输出目录
func createOutputDir(outputPath string) error {
	dir := filepath.Dir(outputPath)
	return os.MkdirAll(dir, DefaultDirPermissions)
}

// initializeWriter 初始化写入器
func (m *Manager) initializeWriter() error {
	var writer Writer
	var err error

	switch m.config.Format {
	case FormatTXT:
		writer, err = NewTXTWriter(m.config.OutputPath)
	case FormatJSON:
		writer, err = NewJSONWriter(m.config.OutputPath)
	case FormatCSV:
		writer, err = NewCSVWriter(m.config.OutputPath)
	default:
		return fmt.Errorf("unsupported format: %s", m.config.Format)
	}

	if err != nil {
		return err
	}

	m.writer = writer
	return m.writer.WriteHeader()
}

// SaveResult 保存扫描结果
func (m *Manager) SaveResult(result *ScanResult) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return fmt.Errorf("output manager is closed")
	}

	if result == nil {
		return fmt.Errorf("result cannot be nil")
	}

	return m.writer.Write(result)
}

// Flush 刷新输出
func (m *Manager) Flush() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return fmt.Errorf("output manager is closed")
	}

	return m.writer.Flush()
}

// Close 关闭输出管理器
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return nil
	}

	m.closed = true
	if m.writer != nil {
		return m.writer.Close()
	}
	return nil
}
