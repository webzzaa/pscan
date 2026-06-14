package core

import (
	"bytes"
	"io"
	"net"
	"testing"
	"time"
)

/*
service_probe_test.go - ServiceProbe核心逻辑测试

注意：service_probe.go 包含大量网络IO和全局状态依赖。
本测试文件专注于可测试的纯逻辑和算法正确性：
1. buildServiceInfo - 数据转换逻辑
2. handleNoMatch - HTTP服务识别逻辑
3. handleHardMatch - 匹配结果处理
4. readFromConn - 缓冲区读取逻辑

不测试的部分（需要集成测试）：
- SmartIdentify, PortInfo - 网络IO + 全局探测器依赖
- Write, Read, Connect - 网络IO操作
- 探测器策略函数 - 依赖全局 PortMap 和 VScan

"这代码把数据结构和网络IO混在一起了，应该分离。
但既然现在无法重构，我们至少测试纯逻辑部分。"
*/

// =============================================================================
// 核心逻辑测试：数据转换
// =============================================================================

// TestBuildServiceInfo 测试服务信息构建逻辑
func TestBuildServiceInfo(t *testing.T) {
	tests := []struct {
		name           string
		setupInfo      func() *SmartPortInfoScanner
		expectedName   string
		expectedBanner string
		hasExtras      bool
	}{
		{
			name: "完整服务信息",
			setupInfo: func() *SmartPortInfoScanner {
				scanner := &SmartPortInfoScanner{
					Address: "192.168.1.1",
					Port:    80,
					info: &Info{
						Result: Result{
							Service: Service{
								Name: "http",
								Extras: map[string]string{
									"version": "Apache/2.4.41",
									"os":      "Linux",
								},
							},
							Banner: "Apache/2.4.41 (Ubuntu)",
						},
					},
				}
				return scanner
			},
			expectedName:   "http",
			expectedBanner: "Apache/2.4.41 (Ubuntu)",
			hasExtras:      true,
		},
		{
			name: "只有服务名称",
			setupInfo: func() *SmartPortInfoScanner {
				scanner := &SmartPortInfoScanner{
					Address: "192.168.1.1",
					Port:    22,
					info: &Info{
						Result: Result{
							Service: Service{
								Name:   "ssh",
								Extras: map[string]string{},
							},
							Banner: "",
						},
					},
				}
				return scanner
			},
			expectedName:   "ssh",
			expectedBanner: "",
			hasExtras:      false,
		},
		{
			name: "未知服务",
			setupInfo: func() *SmartPortInfoScanner {
				scanner := &SmartPortInfoScanner{
					Address: "192.168.1.1",
					Port:    9999,
					info: &Info{
						Result: Result{
							Service: Service{
								Name:   "unknown",
								Extras: map[string]string{},
							},
							Banner: "Binary data",
						},
					},
				}
				return scanner
			},
			expectedName:   "unknown",
			expectedBanner: "Binary data",
			hasExtras:      false,
		},
		{
			name: "包含版本号的服务",
			setupInfo: func() *SmartPortInfoScanner {
				scanner := &SmartPortInfoScanner{
					Address: "192.168.1.1",
					Port:    3306,
					info: &Info{
						Result: Result{
							Service: Service{
								Name: "mysql",
								Extras: map[string]string{
									"version": "5.7.33",
									"product": "MySQL",
								},
							},
							Banner: "MySQL 5.7.33",
						},
					},
				}
				return scanner
			},
			expectedName:   "mysql",
			expectedBanner: "MySQL 5.7.33",
			hasExtras:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := tt.setupInfo()
			serviceInfo := scanner.buildServiceInfo()

			// 验证服务名称
			if serviceInfo.Name != tt.expectedName {
				t.Errorf("Name = %q, 期望 %q", serviceInfo.Name, tt.expectedName)
			}

			// 验证Banner
			if serviceInfo.Banner != tt.expectedBanner {
				t.Errorf("Banner = %q, 期望 %q", serviceInfo.Banner, tt.expectedBanner)
			}

			// 验证Extras
			if tt.hasExtras && len(serviceInfo.Extras) == 0 {
				t.Error("期望有Extras数据，但为空")
			}

			// 验证Version提取
			if version, ok := serviceInfo.Extras["version"]; ok {
				if serviceInfo.Version != version {
					t.Errorf("Version = %q, 期望从Extras提取 %q", serviceInfo.Version, version)
				}
			}

			// 验证Extras不为nil
			if serviceInfo.Extras == nil {
				t.Error("Extras不应为nil")
			}
		})
	}
}

// TestBuildServiceInfo_EmptyExtras 测试空Extras的处理
func TestBuildServiceInfo_EmptyExtras(t *testing.T) {
	scanner := &SmartPortInfoScanner{
		Address: "192.168.1.1",
		Port:    80,
		info: &Info{
			Result: Result{
				Service: Service{
					Name:   "http",
					Extras: nil, // nil Extras
				},
				Banner: "Test",
			},
		},
	}

	serviceInfo := scanner.buildServiceInfo()

	// 验证不会panic
	if serviceInfo.Extras == nil {
		t.Error("Extras应被初始化，不应为nil")
	}

	// 验证Version为空
	if serviceInfo.Version != "" {
		t.Errorf("Version应为空, 实际 %q", serviceInfo.Version)
	}
}

// =============================================================================
// HTTP识别逻辑测试
// =============================================================================

// TestHandleNoMatch_HTTPDetection 测试HTTP服务识别逻辑
func TestHandleNoMatch_HTTPDetection(t *testing.T) {
	tests := []struct {
		name            string
		banner          string
		softFound       bool
		expectedService string
	}{
		{
			name:            "HTTP协议头识别-大写",
			banner:          "HTTP/1.1 200 OK Server: nginx", // TrimBanner会把\r\n替换为空格
			softFound:       false,
			expectedService: "http",
		},
		{
			name:            "HTTP协议头识别-小写http/",
			banner:          "http/1.0 404 Not Found", // 修复后支持小写
			softFound:       false,
			expectedService: "http", // 修复后大小写不敏感
		},
		{
			name:            "HTML内容识别-小写html",
			banner:          "<html><body>Test</body></html>",
			softFound:       false,
			expectedService: "http",
		},
		{
			name:            "HTML内容识别-大写HTML",
			banner:          "<!DOCTYPE HTML>", // 修复后支持大写
			softFound:       false,
			expectedService: "http", // 修复后大小写不敏感
		},
		{
			name:            "HTTP协议头-混合大小写Http/",
			banner:          "Http/2.0 200 OK",
			softFound:       false,
			expectedService: "http",
		},
		{
			name:            "HTML内容-混合大小写HtMl",
			banner:          "<HtMl><body>Test</body></HtMl>",
			softFound:       false,
			expectedService: "http",
		},
		{
			name:            "非HTTP服务",
			banner:          "SSH-2.0-OpenSSH_7.4",
			softFound:       false,
			expectedService: "unknown",
		},
		{
			name:            "空Banner",
			banner:          "",
			softFound:       false,
			expectedService: "unknown",
		},
		{
			name:            "二进制数据",
			banner:          "Binary Data", // TrimBanner把\x00\x01\x02\x03替换为空格，然后TrimSpace
			softFound:       false,
			expectedService: "unknown",
		},
		{
			name:            "软匹配覆盖-不检查HTTP",
			banner:          "HTTP/1.1 200 OK",
			softFound:       true, // 有软匹配时不应识别为HTTP
			expectedService: "test-service",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := &Info{
				Result: Result{},
			}

			// 模拟软匹配
			var softMatch Match
			if tt.softFound {
				softMatch = Match{
					Service: "test-service",
				}
			}

			// 调用handleNoMatch
			info.handleNoMatch([]byte(tt.banner), &info.Result, tt.softFound, softMatch)

			// 验证服务识别结果
			if info.Result.Service.Name != tt.expectedService {
				t.Errorf("Service.Name = %q, 期望 %q", info.Result.Service.Name, tt.expectedService)
			}

			// 验证Banner被正确设置
			if info.Result.Banner != tt.banner {
				t.Errorf("Banner = %q, 期望 %q", info.Result.Banner, tt.banner)
			}

			// 验证Found标志
			if tt.softFound && !info.Found {
				t.Error("软匹配时Found应为true")
			}
		})
	}
}

// TestHandleNoMatch_HTTPVariants 测试HTTP识别的各种变体
func TestHandleNoMatch_HTTPVariants(t *testing.T) {
	// 根据实际实现，只有包含"HTTP/"（大写）或"html"（小写）的才识别为http
	httpVariants := []string{
		"HTTP/1.0 200 OK",
		"HTTP/1.1 404 Not Found",
		"HTTP/2 500 Internal Server Error",
		"<html>",
		"<!DOCTYPE html>",
		"Content-Type: text/html",
	}

	for _, banner := range httpVariants {
		t.Run(banner, func(t *testing.T) {
			info := &Info{
				Result: Result{},
			}

			info.handleNoMatch([]byte(banner), &info.Result, false, Match{})

			if info.Result.Service.Name != "http" {
				t.Errorf("Banner %q 应识别为http, 实际 %q", banner, info.Result.Service.Name)
			}
		})
	}
}

// =============================================================================
// 匹配结果处理测试
// =============================================================================

// TestHandleHardMatch 测试硬匹配处理逻辑
func TestHandleHardMatch(t *testing.T) {
	tests := []struct {
		name             string
		response         []byte
		matchService     string
		expectedService  string
		expectedFound    bool
		checkMicrosoftDS bool
	}{
		{
			name:            "标准HTTP匹配",
			response:        []byte("HTTP/1.1 200 OK\r\nServer: nginx/1.18.0"),
			matchService:    "http",
			expectedService: "http",
			expectedFound:   true,
		},
		{
			name:            "SSH匹配",
			response:        []byte("SSH-2.0-OpenSSH_8.0"),
			matchService:    "ssh",
			expectedService: "ssh",
			expectedFound:   true,
		},
		{
			name:             "Microsoft-DS特殊处理",
			response:         []byte("SMB Domain Info"),
			matchService:     "microsoft-ds",
			expectedService:  "microsoft-ds",
			expectedFound:    true,
			checkMicrosoftDS: true,
		},
		{
			name:            "空响应",
			response:        []byte(""),
			matchService:    "unknown",
			expectedService: "unknown",
			expectedFound:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := &Info{
				Result: Result{
					Service: Service{
						Extras: make(map[string]string),
					},
				},
			}

			// 创建模拟Match
			match := &Match{
				Service: tt.matchService,
			}

			// 调用handleHardMatch
			info.handleHardMatch(tt.response, match)

			// 验证服务名称
			if info.Result.Service.Name != tt.expectedService {
				t.Errorf("Service.Name = %q, 期望 %q", info.Result.Service.Name, tt.expectedService)
			}

			// 验证Found标志
			if info.Found != tt.expectedFound {
				t.Errorf("Found = %v, 期望 %v", info.Found, tt.expectedFound)
			}

			// 验证Banner被设置
			if info.Result.Banner == "" && len(tt.response) > 0 {
				t.Error("Banner应被设置")
			}

			// 验证microsoft-ds特殊处理
			if tt.checkMicrosoftDS {
				if hostname, ok := info.Result.Service.Extras["hostname"]; !ok {
					t.Error("microsoft-ds应设置hostname字段")
				} else if hostname != info.Result.Banner {
					t.Errorf("hostname = %q, 应等于Banner %q", hostname, info.Result.Banner)
				}
			}
		})
	}
}

// =============================================================================
// 缓冲区读取逻辑测试
// =============================================================================

// mockConn 模拟网络连接
type mockConn struct {
	data        []byte
	readPos     int
	chunkSize   int // 每次Read返回的字节数
	closed      bool
	shouldError bool
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	if m.closed {
		return 0, io.EOF
	}

	if m.shouldError {
		return 0, net.ErrClosed
	}

	if m.readPos >= len(m.data) {
		return 0, io.EOF
	}

	// 模拟分块读取
	// chunkSize控制每次Read返回的字节数（不是缓冲区大小）
	readSize := m.chunkSize
	if readSize == 0 {
		// chunkSize=0表示一次性读取整个缓冲区
		readSize = len(b)
	}

	remaining := len(m.data) - m.readPos
	if readSize > remaining {
		readSize = remaining
	}
	if readSize > len(b) {
		readSize = len(b)
	}

	copy(b, m.data[m.readPos:m.readPos+readSize])
	m.readPos += readSize

	// 关键：readFromConn在 count < size 时会停止读取
	// 所以如果 chunkSize > 0，我们要么返回满缓冲区，要么返回EOF
	// 为了测试分块读取，需要让readFromConn认为还有更多数据

	return readSize, nil
}

func (m *mockConn) Write(b []byte) (n int, err error)  { return len(b), nil }
func (m *mockConn) Close() error                       { m.closed = true; return nil }
func (m *mockConn) LocalAddr() net.Addr                { return nil }
func (m *mockConn) RemoteAddr() net.Addr               { return nil }
func (m *mockConn) SetDeadline(t time.Time) error      { return nil }
func (m *mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }

// TestReadFromConn 测试连接读取逻辑
func TestReadFromConn(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		chunkSize   int
		expectedLen int
	}{
		{
			name:        "一次性读取完整数据",
			data:        []byte("Hello, World!"),
			chunkSize:   0, // 0表示一次性读取
			expectedLen: 13,
		},
		{
			name:        "分块读取-填满缓冲区才继续",
			data:        bytes.Repeat([]byte("A"), 5000), // 超过2KB，会分多次读取
			chunkSize:   2048,                            // 每次填满缓冲区
			expectedLen: 5000,
		},
		{
			name:        "分块读取-大数据",
			data:        bytes.Repeat([]byte("Test"), 2048), // 8KB数据
			chunkSize:   2048,                               // 每次2KB
			expectedLen: 8192,
		},
		{
			name:        "空数据",
			data:        []byte{},
			chunkSize:   0,
			expectedLen: 0,
		},
		{
			name:        "小于缓冲区的数据",
			data:        []byte("X"),
			chunkSize:   0,
			expectedLen: 1,
		},
		{
			name:        "恰好填满缓冲区",
			data:        bytes.Repeat([]byte("B"), 2048),
			chunkSize:   2048,
			expectedLen: 2048,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := &mockConn{
				data:      tt.data,
				chunkSize: tt.chunkSize,
			}

			result, err := readFromConn(conn)

			// 验证没有错误
			if err != nil {
				t.Errorf("readFromConn() 错误 = %v", err)
			}

			// 验证读取长度
			if len(result) != tt.expectedLen {
				t.Errorf("读取长度 = %d, 期望 %d", len(result), tt.expectedLen)
			}

			// 验证数据内容
			if !bytes.Equal(result, tt.data) {
				t.Error("读取数据与原始数据不匹配")
			}
		})
	}
}

// TestReadFromConn_EOF 测试EOF处理
func TestReadFromConn_EOF(t *testing.T) {
	conn := &mockConn{
		data:      []byte("Data before EOF"),
		chunkSize: 100,
	}

	result, err := readFromConn(conn)

	// EOF时应返回已读取的数据，不返回错误
	if err != nil {
		t.Errorf("EOF时不应返回错误, 实际 %v", err)
	}

	if len(result) != len(conn.data) {
		t.Errorf("应返回EOF前的数据, 长度 = %d, 期望 %d", len(result), len(conn.data))
	}
}

// TestReadFromConn_Error 测试错误处理
func TestReadFromConn_Error(t *testing.T) {
	conn := &mockConn{
		shouldError: true,
	}

	result, err := readFromConn(conn)

	// 应该返回错误
	if err == nil {
		t.Error("连接错误时应返回错误")
	}

	// 结果应该为空或nil
	if len(result) != 0 {
		t.Errorf("错误时应返回空数据, 实际长度 %d", len(result))
	}
}

// =============================================================================
// 边界情况测试
// =============================================================================

// TestReadFromConn_LargeData 测试大数据读取
func TestReadFromConn_LargeData(t *testing.T) {
	// 模拟10MB数据
	largeData := bytes.Repeat([]byte("X"), 10*1024*1024)

	conn := &mockConn{
		data:      largeData,
		chunkSize: 2048, // 每次读2KB
	}

	result, err := readFromConn(conn)

	if err != nil {
		t.Errorf("大数据读取错误 = %v", err)
	}

	if len(result) != len(largeData) {
		t.Errorf("大数据读取长度 = %d, 期望 %d", len(result), len(largeData))
	}
}

// TestReadFromConn_BinaryData 测试二进制数据
func TestReadFromConn_BinaryData(t *testing.T) {
	binaryData := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}

	conn := &mockConn{
		data:      binaryData,
		chunkSize: 0, // 一次性读取，避免提前终止
	}

	result, err := readFromConn(conn)

	if err != nil {
		t.Errorf("二进制数据读取错误 = %v", err)
	}

	if !bytes.Equal(result, binaryData) {
		t.Errorf("二进制数据 = %v, 期望 %v", result, binaryData)
	}
}
