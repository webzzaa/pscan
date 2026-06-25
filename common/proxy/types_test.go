package proxy

import (
	"errors"
	"testing"
)

/*
types_test.go - 代理类型测试

测试目标：ProxyType枚举、ProxyConfig配置、ProxyError错误
价值：类型定义错误会导致：
  - 代理类型识别错误（连接失败）
  - 配置默认值错误（超时、重试次数）
  - 错误信息丢失（无法调试）

"类型是接口契约。枚举值错了会导致用户无法连接，
默认配置错了会导致超时异常。这些都是真实问题。"
*/

// =============================================================================
// ProxyType - 枚举测试
// =============================================================================

// TestProxyType_String 测试ProxyType.String()方法
//
// 验证：每个枚举值都有正确的字符串表示
func TestProxyType_String(t *testing.T) {
	tests := []struct {
		name      string
		proxyType ProxyType
		expected  string
	}{
		{"None", ProxyTypeNone, "none"},
		{"HTTP", ProxyTypeHTTP, "http"},
		{"HTTPS", ProxyTypeHTTPS, "https"},
		{"SOCKS5", ProxyTypeSOCKS5, "socks5"},
		{"Unknown", ProxyType(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.proxyType.String()
			if result != tt.expected {
				t.Errorf("ProxyType(%d).String() = %q, want %q",
					tt.proxyType, result, tt.expected)
			}

			t.Logf("✓ ProxyType(%d) → %q", tt.proxyType, result)
		})
	}
}

// TestProxyType_AllEnums 测试所有枚举值定义
func TestProxyType_AllEnums(t *testing.T) {
	// 验证枚举值从0开始递增
	tests := []struct {
		name     string
		value    ProxyType
		expected int
	}{
		{"ProxyTypeNone", ProxyTypeNone, 0},
		{"ProxyTypeHTTP", ProxyTypeHTTP, 1},
		{"ProxyTypeHTTPS", ProxyTypeHTTPS, 2},
		{"ProxyTypeSOCKS5", ProxyTypeSOCKS5, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.value) != tt.expected {
				t.Errorf("%s = %d, want %d", tt.name, tt.value, tt.expected)
			}
		})
	}

	t.Logf("✓ 所有ProxyType枚举值定义正确")
}

// =============================================================================
// ProxyConfig - 配置测试
// =============================================================================

// TestDefaultProxyConfig_Values 测试DefaultProxyConfig返回正确的默认值
//
// 验证：默认配置包含所有必要字段的合理值
func TestDefaultProxyConfig_Values(t *testing.T) {
	config := DefaultProxyConfig()

	if config == nil {
		t.Fatal("DefaultProxyConfig() 返回nil")
	}

	// 验证类型
	if config.Type != ProxyTypeNone {
		t.Errorf("默认Type = %v, want %v", config.Type, ProxyTypeNone)
	}

	// 验证超时
	if config.Timeout != DefaultProxyTimeout {
		t.Errorf("默认Timeout = %v, want %v", config.Timeout, DefaultProxyTimeout)
	}

	// 验证重试次数
	if config.MaxRetries != DefaultProxyMaxRetries {
		t.Errorf("默认MaxRetries = %d, want %d", config.MaxRetries, DefaultProxyMaxRetries)
	}

	// 验证KeepAlive
	if config.KeepAlive != DefaultProxyKeepAlive {
		t.Errorf("默认KeepAlive = %v, want %v", config.KeepAlive, DefaultProxyKeepAlive)
	}

	// 验证IdleTimeout
	if config.IdleTimeout != DefaultProxyIdleTimeout {
		t.Errorf("默认IdleTimeout = %v, want %v", config.IdleTimeout, DefaultProxyIdleTimeout)
	}

	// 验证MaxIdleConns
	if config.MaxIdleConns != DefaultProxyMaxIdleConns {
		t.Errorf("默认MaxIdleConns = %d, want %d", config.MaxIdleConns, DefaultProxyMaxIdleConns)
	}

	t.Logf("✓ 默认配置所有字段正确")
}

// TestDefaultProxyConfig_Reasonable 测试默认配置的合理性
func TestDefaultProxyConfig_Reasonable(t *testing.T) {
	config := DefaultProxyConfig()

	// 超时应该 > 0
	if config.Timeout <= 0 {
		t.Error("Timeout应该大于0")
	}

	// 重试次数应该 >= 0
	if config.MaxRetries < 0 {
		t.Error("MaxRetries应该 >= 0")
	}

	// KeepAlive应该 > 0
	if config.KeepAlive <= 0 {
		t.Error("KeepAlive应该大于0")
	}

	// IdleTimeout应该 > 0
	if config.IdleTimeout <= 0 {
		t.Error("IdleTimeout应该大于0")
	}

	// MaxIdleConns应该 > 0
	if config.MaxIdleConns <= 0 {
		t.Error("MaxIdleConns应该大于0")
	}

	// 超时关系：IdleTimeout > Timeout（空闲超时应该更长）
	if config.IdleTimeout < config.Timeout {
		t.Error("IdleTimeout应该大于Timeout")
	}

	t.Logf("✓ 默认配置合理性检查通过")
}

// TestProxyConfig_CustomValues 测试ProxyConfig自定义值
func TestProxyConfig_CustomValues(t *testing.T) {
	config := &ProxyConfig{
		Type:     ProxyTypeHTTP,
		Address:  "127.0.0.1:8080",
		Username: "user",
		Password: "pass",
	}

	// 测试字段值是否正确赋值
	_ = config.Type
	_ = config.Address
	_ = config.Username
	_ = config.Password

	if config.Type != ProxyTypeHTTP {
		t.Error("自定义Type赋值失败")
	}

	if config.Address != "127.0.0.1:8080" {
		t.Error("自定义Address赋值失败")
	}

	if config.Username != "user" {
		t.Error("自定义Username赋值失败")
	}

	if config.Password != "pass" {
		t.Error("自定义Password赋值失败")
	}

	t.Logf("✓ ProxyConfig自定义值测试通过")
}

// =============================================================================
// ProxyError - 错误类型测试
// =============================================================================

// TestProxyError_Error 测试ProxyError.Error()方法
//
// 验证：错误信息格式正确
func TestProxyError_Error(t *testing.T) {
	t.Run("无Cause", func(t *testing.T) {
		err := &ProxyError{
			Type:    "test_error",
			Message: "test message",
			Code:    100,
		}

		expected := "test message"
		if err.Error() != expected {
			t.Errorf("Error() = %q, want %q", err.Error(), expected)
		}

		t.Logf("✓ ProxyError无Cause时返回纯Message")
	})

	t.Run("有Cause", func(t *testing.T) {
		cause := errors.New("root cause")
		err := &ProxyError{
			Type:    "test_error",
			Message: "test message",
			Code:    100,
			Cause:   cause,
		}

		expected := "test message: root cause"
		if err.Error() != expected {
			t.Errorf("Error() = %q, want %q", err.Error(), expected)
		}

		t.Logf("✓ ProxyError有Cause时正确拼接")
	})
}

// TestNewProxyError 测试NewProxyError构造函数
func TestNewProxyError(t *testing.T) {
	t.Run("无Cause", func(t *testing.T) {
		err := NewProxyError("config_error", "invalid config", 1001, nil)

		if err.Type != "config_error" {
			t.Errorf("Type = %q, want %q", err.Type, "config_error")
		}

		if err.Message != "invalid config" {
			t.Errorf("Message = %q, want %q", err.Message, "invalid config")
		}

		if err.Code != 1001 {
			t.Errorf("Code = %d, want %d", err.Code, 1001)
		}

		if err.Cause != nil {
			t.Error("Cause应该为nil")
		}

		t.Logf("✓ NewProxyError无Cause测试通过")
	})

	t.Run("有Cause", func(t *testing.T) {
		cause := errors.New("connection refused")
		err := NewProxyError("connection_error", "failed to connect", 2001, cause)

		if err.Type != "connection_error" {
			t.Errorf("Type = %q, want %q", err.Type, "connection_error")
		}

		if err.Message != "failed to connect" {
			t.Errorf("Message = %q, want %q", err.Message, "failed to connect")
		}

		if err.Code != 2001 {
			t.Errorf("Code = %d, want %d", err.Code, 2001)
		}

		if !errors.Is(err.Cause, cause) {
			t.Error("Cause应该是传入的cause")
		}

		t.Logf("✓ NewProxyError有Cause测试通过")
	})
}

// TestProxyError_AllErrorTypes 测试所有预定义错误类型常量
func TestProxyError_AllErrorTypes(t *testing.T) {
	errorTypes := []struct {
		name     string
		constant string
	}{
		{"Config", ErrTypeConfig},
		{"Connection", ErrTypeConnection},
		{"Auth", ErrTypeAuth},
		{"Timeout", ErrTypeTimeout},
		{"Protocol", ErrTypeProtocol},
	}

	for _, et := range errorTypes {
		t.Run(et.name, func(t *testing.T) {
			if et.constant == "" {
				t.Errorf("%s错误类型常量为空", et.name)
			}

			// 使用错误类型创建错误
			err := NewProxyError(et.constant, et.name+" error", 0, nil)
			if err.Type != et.constant {
				t.Errorf("Type = %q, want %q", err.Type, et.constant)
			}

			t.Logf("✓ %s错误类型: %q", et.name, et.constant)
		})
	}
}

// =============================================================================
// 常量测试
// =============================================================================

// TestProxyConstants_Reasonable 测试常量合理性
func TestProxyConstants_Reasonable(t *testing.T) {
	// 超时常量应该大于0
	if DefaultProxyTimeout <= 0 {
		t.Error("DefaultProxyTimeout应该大于0")
	}

	// 重试次数应该 >= 0
	if DefaultProxyMaxRetries < 0 {
		t.Error("DefaultProxyMaxRetries应该 >= 0")
	}

	// KeepAlive应该大于0
	if DefaultProxyKeepAlive <= 0 {
		t.Error("DefaultProxyKeepAlive应该大于0")
	}

	// IdleTimeout应该大于0
	if DefaultProxyIdleTimeout <= 0 {
		t.Error("DefaultProxyIdleTimeout应该大于0")
	}

	// MaxIdleConns应该大于0
	if DefaultProxyMaxIdleConns <= 0 {
		t.Error("DefaultProxyMaxIdleConns应该大于0")
	}

	t.Logf("✓ 所有代理常量合理")
}

// TestProxyTypeStrings_NoEmpty 测试代理类型字符串非空
func TestProxyTypeStrings_NoEmpty(t *testing.T) {
	typeStrings := []struct {
		name  string
		value string
	}{
		{"ProxyTypeStringNone", ProxyTypeStringNone},
		{"ProxyTypeStringHTTP", ProxyTypeStringHTTP},
		{"ProxyTypeStringHTTPS", ProxyTypeStringHTTPS},
		{"ProxyTypeStringSOCKS5", ProxyTypeStringSOCKS5},
		{"ProxyTypeStringUnknown", ProxyTypeStringUnknown},
	}

	for _, ts := range typeStrings {
		t.Run(ts.name, func(t *testing.T) {
			if ts.value == "" {
				t.Errorf("%s不应为空", ts.name)
			}

			t.Logf("✓ %s = %q", ts.name, ts.value)
		})
	}
}
