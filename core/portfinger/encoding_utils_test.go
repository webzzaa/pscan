package portfinger

import (
	"bytes"
	"testing"
)

// TestDecodePattern 测试nmap探测数据解码
func TestDecodePattern(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []byte
	}{
		{
			name:     "十六进制编码-单字节",
			input:    `\x48`,
			expected: []byte{0x48}, // 'H'
		},
		{
			name:     "十六进制编码-多字节",
			input:    `\x48\x65\x6c\x6c\x6f`,
			expected: []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f}, // "Hello"
		},
		{
			name:     "转义字符-换行",
			input:    `\n`,
			expected: []byte{'\n'},
		},
		{
			name:     "转义字符-回车",
			input:    `\r`,
			expected: []byte{'\r'},
		},
		{
			name:     "转义字符-制表符",
			input:    `\t`,
			expected: []byte{'\t'},
		},
		{
			name:     "转义字符-响铃",
			input:    `\a`,
			expected: []byte{'\a'},
		},
		{
			name:     "转义字符-换页",
			input:    `\f`,
			expected: []byte{'\f'},
		},
		{
			name:     "转义字符-垂直制表符",
			input:    `\v`,
			expected: []byte{'\v'},
		},
		{
			name:     "转义字符-反斜杠",
			input:    `\\`,
			expected: []byte{'\\'},
		},
		{
			name:     "八进制编码-单字节",
			input:    `\101`,
			expected: []byte{0101}, // 'A' (65)
		},
		{
			name:     "八进制编码-两位",
			input:    `\72`,
			expected: []byte{072}, // ':' (58)
		},
		{
			name:     "八进制编码-一位",
			input:    `\7`,
			expected: []byte{7},
		},
		{
			name:     "混合编码-nmap GET请求",
			input:    `GET / HTTP/1.0\r\n\r\n`,
			expected: []byte("GET / HTTP/1.0\r\n\r\n"),
		},
		{
			name:     "混合编码-十六进制+文本",
			input:    `\x48ello`,
			expected: []byte("Hello"),
		},
		{
			name:     "普通文本",
			input:    `Hello World`,
			expected: []byte("Hello World"),
		},
		{
			name:     "空字符串",
			input:    ``,
			expected: []byte{},
		},
		{
			name:     "复杂nmap探测数据",
			input:    `\x00\x00\x00\x01\x02\x03`,
			expected: []byte{0x00, 0x00, 0x00, 0x01, 0x02, 0x03},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := DecodePattern(tt.input)
			if err != nil {
				t.Fatalf("DecodePattern() 错误 = %v", err)
			}
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("DecodePattern() = %v (%q), 期望 %v (%q)",
					result, string(result), tt.expected, string(tt.expected))
			}
		})
	}
}

// TestDecodePattern_InvalidHex 测试非法十六进制编码
func TestDecodePattern_InvalidHex(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "不完整的十六进制-只有\\x",
			input: `\x`,
		},
		{
			name:  "不完整的十六进制-只有一位",
			input: `\xA`,
		},
		{
			name:  "非法十六进制字符",
			input: `\xGH`,
		},
		{
			name:  "十六进制后截断",
			input: `Hello\x`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := DecodePattern(tt.input)
			// 非法的十六进制应该被忽略，返回原字符
			if err != nil {
				t.Errorf("DecodePattern() 不应返回错误: %v", err)
			}
			// 验证至少有输出（即使不正确也不应panic）
			if result == nil {
				t.Error("DecodePattern() 不应返回 nil")
			}
		})
	}
}

// TestDecodePattern_OctalEdgeCases 测试八进制边界情况
func TestDecodePattern_OctalEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []byte
	}{
		{
			name:     "八进制最大值int8-127",
			input:    `\177`,
			expected: []byte{0177}, // 127, int8最大值
		},
		{
			name:     "八进制零",
			input:    `\0`,
			expected: []byte{0},
		},
		{
			name:     "八进制混合",
			input:    `\101\102\103`,
			expected: []byte{'A', 'B', 'C'},
		},
		{
			name:     "八进制后跟普通数字",
			input:    `\1018`,
			expected: []byte{0101, '8'}, // 'A' + '8'
		},
		{
			name:     "八进制最大值-255",
			input:    `\377`,
			expected: []byte{0xFF}, // 255, 八进制最大值
		},
		{
			name:     "八进制超出255-按原字符",
			input:    `\777`,
			expected: []byte{'\\', '7', '7', '7'}, // 超出范围，按原字符处理
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := DecodePattern(tt.input)
			if err != nil {
				t.Fatalf("DecodePattern() 错误 = %v", err)
			}
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("DecodePattern() = %v, 期望 %v", result, tt.expected)
			}
		})
	}
}

// TestDecodeData 测试DecodeData包装器
func TestDecodeData(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []byte
	}{
		{
			name:     "双引号包裹",
			input:    `"Hello"`,
			expected: []byte("Hello"),
		},
		{
			name:     "单引号包裹",
			input:    `'World'`,
			expected: []byte("World"),
		},
		{
			name:     "双引号包裹+转义",
			input:    `"\x48\x65\x6c\x6c\x6f"`,
			expected: []byte("Hello"),
		},
		{
			name:     "无引号",
			input:    `Hello`,
			expected: []byte("Hello"),
		},
		{
			name:     "只有开头引号",
			input:    `"Hello`,
			expected: []byte("Hello"),
		},
		{
			name:     "只有结尾引号",
			input:    `Hello"`,
			expected: []byte("Hello"),
		},
		{
			name:     "空字符串-双引号",
			input:    `""`,
			expected: []byte{},
		},
		{
			name:     "nmap探测数据格式",
			input:    `"GET / HTTP/1.0\r\n\r\n"`,
			expected: []byte("GET / HTTP/1.0\r\n\r\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := DecodeData(tt.input)
			if err != nil {
				t.Fatalf("DecodeData() 错误 = %v", err)
			}
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("DecodeData() = %v (%q), 期望 %v (%q)",
					result, string(result), tt.expected, string(tt.expected))
			}
		})
	}
}

// TestIsValidHex 测试十六进制验证
func TestIsValidHex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"合法-数字", "12", true},
		{"合法-小写字母", "ab", true},
		{"合法-大写字母", "AB", true},
		{"合法-混合", "3F", true},
		{"合法-全0", "00", true},
		{"合法-全F", "FF", true},
		{"非法-单字符", "A", false},
		{"非法-三字符", "ABC", false},
		{"非法-空字符串", "", false},
		{"非法-包含G", "AG", false},
		{"非法-包含特殊字符", "A@", false},
		{"非法-包含空格", "A ", false},
		{"非法-汉字", "中文", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidHex(tt.input)
			if result != tt.expected {
				t.Errorf("isValidHex(%q) = %v, 期望 %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestDecodePattern_RealWorldNmapData 测试真实nmap探测数据
func TestDecodePattern_RealWorldNmapData(t *testing.T) {
	tests := []struct {
		name  string
		input string
		desc  string
	}{
		{
			name:  "HTTP GET请求",
			input: `GET / HTTP/1.0\r\n\r\n`,
			desc:  "nmap HTTP探测",
		},
		{
			name:  "SSH握手",
			input: `SSH-2.0-OpenSSH_8.0\r\n`,
			desc:  "SSH版本探测",
		},
		{
			name:  "MySQL握手",
			input: `\x00\x00\x00\x0a5.7.0`,
			desc:  "MySQL协议",
		},
		{
			name:  "二进制协议",
			input: `\x00\x01\x02\x03\x04\x05`,
			desc:  "纯二进制数据",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := DecodePattern(tt.input)
			if err != nil {
				t.Errorf("%s 解码失败: %v", tt.desc, err)
			}
			if len(result) == 0 {
				t.Errorf("%s 解码结果为空", tt.desc)
			}
			t.Logf("%s 解码成功: %d 字节", tt.desc, len(result))
		})
	}
}

// BenchmarkDecodePattern 基准测试DecodePattern
func BenchmarkDecodePattern(b *testing.B) {
	input := `GET / HTTP/1.0\r\n\r\n`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DecodePattern(input)
	}
}

// BenchmarkDecodePattern_Complex 基准测试复杂编码
func BenchmarkDecodePattern_Complex(b *testing.B) {
	input := `\x48\x65\x6c\x6c\x6f\x20\x57\x6f\x72\x6c\x64\r\n`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DecodePattern(input)
	}
}
