package config

import (
	"strconv"
	"strings"
	"testing"
)

/*
constants_test.go - 配置常量测试

测试目标：端口组、探测器配置、字典数据
价值：配置错误会导致：
  - 端口组错误 → 扫描范围错误（用户遗漏目标）
  - 字典错误 → 暴力破解失败（无法登录系统）
  - 探测器配置错误 → 服务识别失败

"配置是数据，但数据也会有bug。端口范围错误、字典重复、
空值遗漏——这些都是真实问题。测试数据和测试代码一样重要。"
*/

// =============================================================================
// 端口组测试
// =============================================================================

// TestPortGroups_Format 测试端口组格式
//
// 验证：所有端口组字符串格式正确（可解析为端口列表）
func TestPortGroups_Format(t *testing.T) {
	tests := []struct {
		name      string
		portGroup string
	}{
		{"WebPorts", WebPorts},
		{"MainPorts", MainPorts},
		{"DbPorts", DbPorts},
		{"ServicePorts", ServicePorts},
		{"CommonPorts", CommonPorts},
		{"AllPorts", AllPorts},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 验证格式：逗号分隔的数字或范围
			if tt.portGroup == "" {
				t.Error("端口组不应为空")
				return
			}

			// AllPorts是特殊格式"1-65535"
			if tt.portGroup == "1-65535" {
				t.Logf("✓ %s 格式正确（范围格式）", tt.name)
				return
			}

			// 其他端口组应该是逗号分隔的数字
			ports := strings.Split(tt.portGroup, ",")
			if len(ports) == 0 {
				t.Error("端口组应该包含至少一个端口")
				return
			}

			// 验证每个端口都是有效数字
			for i, portStr := range ports {
				port, err := strconv.Atoi(strings.TrimSpace(portStr))
				if err != nil {
					t.Errorf("第%d个端口 '%s' 不是有效数字: %v", i+1, portStr, err)
					continue
				}

				// 验证端口范围
				if port < 1 || port > 65535 {
					t.Errorf("第%d个端口 %d 超出有效范围 [1-65535]", i+1, port)
				}
			}

			t.Logf("✓ %s 格式正确（%d个端口）", tt.name, len(ports))
		})
	}
}

// TestPortGroups_NoEmpty 测试端口组非空
func TestPortGroups_NoEmpty(t *testing.T) {
	groups := map[string]string{
		"WebPorts":     WebPorts,
		"MainPorts":    MainPorts,
		"DbPorts":      DbPorts,
		"ServicePorts": ServicePorts,
		"CommonPorts":  CommonPorts,
		"AllPorts":     AllPorts,
	}

	for name, ports := range groups {
		if ports == "" {
			t.Errorf("%s 不应为空字符串", name)
		}
	}

	t.Logf("✓ 所有端口组非空")
}

// TestPortGroups_NoDuplicates 测试端口组无重复
func TestPortGroups_NoDuplicates(t *testing.T) {
	tests := []struct {
		name      string
		portGroup string
	}{
		{"WebPorts", WebPorts},
		{"MainPorts", MainPorts},
		{"DbPorts", DbPorts},
		{"ServicePorts", ServicePorts},
		{"CommonPorts", CommonPorts},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.portGroup == "1-65535" {
				t.Skip("范围格式无需检查重复")
				return
			}

			ports := strings.Split(tt.portGroup, ",")
			seen := make(map[string]bool)
			duplicates := []string{}

			for _, port := range ports {
				port = strings.TrimSpace(port)
				if seen[port] {
					duplicates = append(duplicates, port)
				}
				seen[port] = true
			}

			if len(duplicates) > 0 {
				t.Errorf("%s 包含重复端口: %v", tt.name, duplicates)
			} else {
				t.Logf("✓ %s 无重复端口", tt.name)
			}
		})
	}
}

// TestGetPortGroups_Completeness 测试GetPortGroups完整性
//
// 验证：返回的map包含所有预定义的端口组
func TestGetPortGroups_Completeness(t *testing.T) {
	groups := GetPortGroups()

	expectedKeys := []string{"web", "main", "db", "service", "common", "all"}
	for _, key := range expectedKeys {
		if _, ok := groups[key]; !ok {
			t.Errorf("GetPortGroups缺少键: %s", key)
		}
	}

	if len(groups) != len(expectedKeys) {
		t.Errorf("GetPortGroups返回%d个组，期望%d个", len(groups), len(expectedKeys))
	}

	t.Logf("✓ GetPortGroups包含所有%d个端口组", len(expectedKeys))
}

// TestGetPortGroups_Values 测试GetPortGroups返回正确的值
func TestGetPortGroups_Values(t *testing.T) {
	groups := GetPortGroups()

	tests := []struct {
		key      string
		expected string
	}{
		{"web", WebPorts},
		{"main", MainPorts},
		{"db", DbPorts},
		{"service", ServicePorts},
		{"common", CommonPorts},
		{"all", AllPorts},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			actual, ok := groups[tt.key]
			if !ok {
				t.Fatalf("GetPortGroups缺少键: %s", tt.key)
			}

			if actual != tt.expected {
				t.Errorf("GetPortGroups[%s] 值不匹配\n期望前20字符: %s...\n实际前20字符: %s...",
					tt.key, tt.expected[:20], actual[:20])
			}

			t.Logf("✓ %s 映射正确", tt.key)
		})
	}
}

// =============================================================================
// 探测器配置测试
// =============================================================================

// TestDefaultProbeMap_NoEmpty 测试默认探测器列表非空
func TestDefaultProbeMap_NoEmpty(t *testing.T) {
	if len(DefaultProbeMap) == 0 {
		t.Error("DefaultProbeMap不应为空")
		return
	}

	// 验证每个探测器名称非空
	for i, probe := range DefaultProbeMap {
		if probe == "" {
			t.Errorf("第%d个探测器名称为空", i+1)
		}
	}

	t.Logf("✓ DefaultProbeMap包含%d个探测器", len(DefaultProbeMap))
}

// TestDefaultPortMap_ValidKeys 测试DefaultPortMap的键有效
func TestDefaultPortMap_ValidKeys(t *testing.T) {
	invalidPorts := []int{}

	for port := range DefaultPortMap {
		if port < 1 || port > 65535 {
			invalidPorts = append(invalidPorts, port)
		}
	}

	if len(invalidPorts) > 0 {
		t.Errorf("DefaultPortMap包含无效端口号: %v", invalidPorts)
	} else {
		t.Logf("✓ DefaultPortMap的%d个端口号都有效", len(DefaultPortMap))
	}
}

// TestDefaultPortMap_NoEmptyValues 测试DefaultPortMap值非空
func TestDefaultPortMap_NoEmptyValues(t *testing.T) {
	emptyPorts := []int{}

	for port, probes := range DefaultPortMap {
		if len(probes) == 0 {
			emptyPorts = append(emptyPorts, port)
		}
	}

	if len(emptyPorts) > 0 {
		t.Errorf("以下端口的探测器列表为空: %v", emptyPorts)
	} else {
		t.Logf("✓ DefaultPortMap所有端口都有探测器")
	}
}

// =============================================================================
// 字典数据测试
// =============================================================================

// TestDefaultUserDict_NoEmptyKeys 测试DefaultUserDict键非空
func TestDefaultUserDict_NoEmptyKeys(t *testing.T) {
	for service, users := range DefaultUserDict {
		if service == "" {
			t.Error("DefaultUserDict包含空服务名")
		}

		if len(users) == 0 {
			t.Errorf("服务 '%s' 的用户列表为空", service)
		}
	}

	t.Logf("✓ DefaultUserDict包含%d个服务", len(DefaultUserDict))
}

// TestDefaultUserDict_CommonServices 测试DefaultUserDict包含常见服务
func TestDefaultUserDict_CommonServices(t *testing.T) {
	commonServices := []string{"ftp", "mysql", "mssql", "ssh", "redis", "mongodb"}

	for _, service := range commonServices {
		if _, ok := DefaultUserDict[service]; !ok {
			t.Errorf("DefaultUserDict缺少常见服务: %s", service)
		}
	}

	t.Logf("✓ DefaultUserDict包含所有常见服务")
}

// TestDefaultUserDict_AllowsEmptyUser 测试DefaultUserDict允许空用户名
//
// 验证：某些服务（如redis）允许空用户名
func TestDefaultUserDict_AllowsEmptyUser(t *testing.T) {
	// redis服务应该包含空用户名
	redisUsers, ok := DefaultUserDict["redis"]
	if !ok {
		t.Skip("DefaultUserDict不包含redis，跳过测试")
		return
	}

	hasEmptyUser := false
	for _, user := range redisUsers {
		if user == "" {
			hasEmptyUser = true
			break
		}
	}

	if !hasEmptyUser {
		t.Error("redis用户列表应该包含空用户名（默认无认证）")
	} else {
		t.Logf("✓ redis用户列表正确包含空用户名")
	}
}

// TestDefaultPasswords_NoEmpty 测试DefaultPasswords非空
func TestDefaultPasswords_NoEmpty(t *testing.T) {
	if len(DefaultPasswords) == 0 {
		t.Error("DefaultPasswords不应为空")
		return
	}

	t.Logf("✓ DefaultPasswords包含%d个密码", len(DefaultPasswords))
}

// TestDefaultPasswords_AllowsEmptyPassword 测试DefaultPasswords允许空密码
func TestDefaultPasswords_AllowsEmptyPassword(t *testing.T) {
	// 应该包含空密码（某些服务默认无密码）
	hasEmptyPassword := false
	for _, pass := range DefaultPasswords {
		if pass == "" {
			hasEmptyPassword = true
			break
		}
	}

	if !hasEmptyPassword {
		t.Error("DefaultPasswords应该包含空密码（某些服务默认无密码）")
	} else {
		t.Logf("✓ DefaultPasswords正确包含空密码")
	}
}

// TestDefaultPasswords_HasPlaceholder 测试DefaultPasswords包含占位符
func TestDefaultPasswords_HasPlaceholder(t *testing.T) {
	// 应该包含{user}占位符（密码=用户名的场景）
	hasPlaceholder := false
	for _, pass := range DefaultPasswords {
		if strings.Contains(pass, "{user}") {
			hasPlaceholder = true
			break
		}
	}

	if !hasPlaceholder {
		t.Error("DefaultPasswords应该包含{user}占位符（密码=用户名变体）")
	} else {
		t.Logf("✓ DefaultPasswords正确包含{user}占位符")
	}
}

// =============================================================================
// 结构体测试
// =============================================================================

// TestPocInfo_Fields 测试PocInfo结构体字段
func TestPocInfo_Fields(t *testing.T) {
	poc := PocInfo{
		Target:  "http://example.com",
		PocName: "test-poc",
	}

	if poc.Target != "http://example.com" {
		t.Error("PocInfo.Target赋值失败")
	}

	if poc.PocName != "test-poc" {
		t.Error("PocInfo.PocName赋值失败")
	}

	t.Logf("✓ PocInfo结构体正常工作")
}

// TestCredentialPair_Fields 测试CredentialPair结构体字段
func TestCredentialPair_Fields(t *testing.T) {
	cred := CredentialPair{
		Username: "admin",
		Password: "password123",
	}

	if cred.Username != "admin" {
		t.Error("CredentialPair.Username赋值失败")
	}

	if cred.Password != "password123" {
		t.Error("CredentialPair.Password赋值失败")
	}

	t.Logf("✓ CredentialPair结构体正常工作")
}
