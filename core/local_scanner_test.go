package core

import (
	"testing"

	"test/common"
)

// TestNewLocalScanStrategy 测试本地扫描策略构造函数
func TestNewLocalScanStrategy(t *testing.T) {
	strategy := NewLocalScanStrategy()

	if strategy == nil {
		t.Fatal("NewLocalScanStrategy 返回 nil")
	}

	if strategy.BaseScanStrategy == nil {
		t.Error("BaseScanStrategy 未初始化")
	}

	// 验证过滤器类型
	if strategy.filterType != FilterLocal {
		t.Errorf("filterType: 期望 FilterLocal(%d), 实际 %d", FilterLocal, strategy.filterType)
	}

	// 验证策略名称
	if strategy.strategyName != "本地扫描" {
		t.Errorf("strategyName: 期望 '本地扫描', 实际 %q", strategy.strategyName)
	}
}

// TestLocalScanStrategy_PrepareTargets 测试PrepareTargets
func TestLocalScanStrategy_PrepareTargets(t *testing.T) {
	strategy := NewLocalScanStrategy()

	tests := []struct {
		name     string
		input    common.HostInfo
		expected int
	}{
		{
			name:     "空HostInfo",
			input:    common.HostInfo{},
			expected: 1,
		},
		{
			name: "带Host的HostInfo",
			input: common.HostInfo{
				Host: "localhost",
			},
			expected: 1,
		},
		{
			name: "完整HostInfo",
			input: common.HostInfo{
				Host: "127.0.0.1",
				Port: 80,
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targets := strategy.PrepareTargets(tt.input)

			// 验证返回列表长度
			if len(targets) != tt.expected {
				t.Errorf("PrepareTargets() 返回长度 = %d, 期望 %d", len(targets), tt.expected)
			}

			// 验证返回的第一个元素与输入相同
			if len(targets) > 0 {
				if targets[0].Host != tt.input.Host {
					t.Errorf("targets[0].Host = %q, 期望 %q", targets[0].Host, tt.input.Host)
				}
				if targets[0].Port != tt.input.Port {
					t.Errorf("targets[0].Port = %q, 期望 %q", targets[0].Port, tt.input.Port)
				}
			}
		})
	}
}

// TestLocalScanStrategy_PrepareTargets_ImmutabilityCheck 测试PrepareTargets不修改输入
func TestLocalScanStrategy_PrepareTargets_ImmutabilityCheck(t *testing.T) {
	strategy := NewLocalScanStrategy()

	original := common.HostInfo{
		Host: "192.168.1.1",
		Port: 22,
	}

	// 保存原始值副本
	originalHost := original.Host
	originalPort := original.Port

	// 调用PrepareTargets
	targets := strategy.PrepareTargets(original)

	// 验证原始输入未被修改
	if original.Host != originalHost {
		t.Errorf("输入被修改: original.Host = %q, 期望 %q", original.Host, originalHost)
	}
	if original.Port != originalPort {
		t.Errorf("输入被修改: original.Port = %d, 期望 %d", original.Port, originalPort)
	}

	// 验证返回值与输入相等
	if len(targets) != 1 {
		t.Fatalf("targets长度 = %d, 期望 1", len(targets))
	}

	if targets[0].Host != originalHost {
		t.Errorf("targets[0].Host = %q, 期望 %q", targets[0].Host, originalHost)
	}
}

// TestLocalScanStrategy_TypeAssertion 测试类型继承关系
func TestLocalScanStrategy_TypeAssertion(t *testing.T) {
	strategy := NewLocalScanStrategy()

	// 验证类型继承
	if strategy.BaseScanStrategy == nil {
		t.Error("LocalScanStrategy 未嵌入 BaseScanStrategy")
	}

	// 验证可以访问BaseScanStrategy的方法
	err := strategy.ValidateConfiguration()
	if err != nil {
		t.Errorf("ValidateConfiguration() 应返回 nil, 实际: %v", err)
	}
}

// TestLocalScanStrategy_FieldAccess 测试字段访问
func TestLocalScanStrategy_FieldAccess(t *testing.T) {
	strategy := NewLocalScanStrategy()

	// 通过BaseScanStrategy访问私有字段
	if strategy.strategyName == "" {
		t.Error("strategyName 不应为空")
	}

	if strategy.filterType != FilterLocal {
		t.Errorf("filterType 应为 FilterLocal, 实际 %d", strategy.filterType)
	}
}
