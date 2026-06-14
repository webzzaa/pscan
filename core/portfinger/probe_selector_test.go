package portfinger

import (
	"testing"
)

func TestPortInRange(t *testing.T) {
	tests := []struct {
		name     string
		port     int
		portsStr string
		expected bool
	}{
		{
			name:     "单个端口匹配",
			port:     80,
			portsStr: "80",
			expected: true,
		},
		{
			name:     "单个端口不匹配",
			port:     81,
			portsStr: "80",
			expected: false,
		},
		{
			name:     "端口列表匹配",
			port:     443,
			portsStr: "80,443,8080",
			expected: true,
		},
		{
			name:     "端口列表不匹配",
			port:     8443,
			portsStr: "80,443,8080",
			expected: false,
		},
		{
			name:     "端口范围匹配-起点",
			port:     1000,
			portsStr: "1000-2000",
			expected: true,
		},
		{
			name:     "端口范围匹配-终点",
			port:     2000,
			portsStr: "1000-2000",
			expected: true,
		},
		{
			name:     "端口范围匹配-中间",
			port:     1500,
			portsStr: "1000-2000",
			expected: true,
		},
		{
			name:     "端口范围不匹配-小于起点",
			port:     999,
			portsStr: "1000-2000",
			expected: false,
		},
		{
			name:     "端口范围不匹配-大于终点",
			port:     2001,
			portsStr: "1000-2000",
			expected: false,
		},
		{
			name:     "混合格式匹配-单端口",
			port:     22,
			portsStr: "22,80,443,1000-2000,8080",
			expected: true,
		},
		{
			name:     "混合格式匹配-范围内",
			port:     1234,
			portsStr: "22,80,443,1000-2000,8080",
			expected: true,
		},
		{
			name:     "混合格式不匹配",
			port:     3000,
			portsStr: "22,80,443,1000-2000,8080",
			expected: false,
		},
		{
			name:     "空字符串",
			port:     80,
			portsStr: "",
			expected: false,
		},
		{
			name:     "带空格的端口列表",
			port:     443,
			portsStr: "80, 443, 8080",
			expected: true,
		},
		{
			name:     "Nmap格式-GetRequest探测器端口",
			port:     8080,
			portsStr: "80,81,82,83,84,85,86,87,88,89,90,280,443,591,593,623,664,777,808,832,888,901,981,1010,1080,1100,1241,1311,1352,1434,1944,2301,2381,2574,3000,3128,3268,4000,4001,4002,4100,4444,5000,5050,5432,5555,5800,5801,5802,5803,6080,7000,7001,7002,7103,7201,7777,7778,8000,8001,8002,8003,8006,8008,8009,8014,8042,8080,8081,8082,8083,8084,8085,8087,8088,8089,8090,8091,8100,8118,8123,8172,8180,8181,8200,8222,8243,8280,8281,8333,8383,8400,8443,8500,8509,8787,8800,8888,8899,8983,9000,9001,9002,9080,9090,9091,9100,9200,9443,9990,9999,10000,10443,12443,16080,18091,18092,20720,28017",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PortInRange(tt.port, tt.portsStr)
			if result != tt.expected {
				t.Errorf("PortInRange(%d, %q) = %v, want %v", tt.port, tt.portsStr, result, tt.expected)
			}
		})
	}
}

func TestGetProbesForPort(t *testing.T) {
	// 确保全局 VScan 已初始化
	InitializeGlobalVScan()
	v := GetGlobalVScan()

	// 测试常见端口
	// 注意：SSH(22) 和 MySQL(3306) 等服务在 nmap 规则中不使用 ports 字段
	// 它们依赖 NULL 探测器（等待服务主动发送 banner）
	tests := []struct {
		port        int
		expectFound bool
		description string
	}{
		{port: 80, expectFound: true, description: "HTTP端口应该有探测器"},
		{port: 22, expectFound: false, description: "SSH端口使用NULL探测（无ports字段）"},
		{port: 443, expectFound: true, description: "HTTPS端口应该有探测器"},
		{port: 3306, expectFound: false, description: "MySQL端口使用NULL探测（无ports字段）"},
		{port: 1, expectFound: true, description: "端口1有GetRequest和Help探测器"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			probes := v.GetProbesForPort(tt.port)
			if tt.expectFound && len(probes) == 0 {
				t.Errorf("端口 %d: 期望找到探测器，但找到 %d 个", tt.port, len(probes))
			}
			if len(probes) > 0 {
				t.Logf("端口 %d: 找到 %d 个探测器", tt.port, len(probes))
				for i, p := range probes {
					t.Logf("  [%d] %s (rarity=%d)", i+1, p.Name, p.Rarity)
				}
			} else {
				t.Logf("端口 %d: 无特定探测器（使用NULL探测）", tt.port)
			}
		})
	}
}

func TestGetProbesForPort_RaritySorting(t *testing.T) {
	InitializeGlobalVScan()
	v := GetGlobalVScan()

	// 获取端口80的探测器（应该有多个）
	probes := v.GetProbesForPort(80)
	if len(probes) < 2 {
		t.Skip("端口80的探测器数量不足，跳过排序测试")
	}

	// 验证按 rarity 排序（从低到高）
	for i := 1; i < len(probes); i++ {
		prev := probes[i-1].Rarity
		curr := probes[i].Rarity
		// 将0视为10（最低优先级）
		if prev == 0 {
			prev = 10
		}
		if curr == 0 {
			curr = 10
		}
		if prev > curr {
			t.Errorf("探测器未按rarity排序: probes[%d].Rarity=%d > probes[%d].Rarity=%d",
				i-1, probes[i-1].Rarity, i, probes[i].Rarity)
		}
	}
}

func TestFilterProbesByIntensity(t *testing.T) {
	// 创建模拟探测器
	probes := []*Probe{
		{Name: "p1", Rarity: 1},
		{Name: "p2", Rarity: 3},
		{Name: "p3", Rarity: 5},
		{Name: "p4", Rarity: 7},
		{Name: "p5", Rarity: 9},
		{Name: "p6", Rarity: 0}, // 0 视为 1
	}

	tests := []struct {
		intensity     int
		expectedCount int
	}{
		{intensity: 1, expectedCount: 2},  // p1, p6
		{intensity: 3, expectedCount: 3},  // p1, p2, p6
		{intensity: 5, expectedCount: 4},  // p1, p2, p3, p6
		{intensity: 7, expectedCount: 5},  // p1, p2, p3, p4, p6
		{intensity: 9, expectedCount: 6},  // all
		{intensity: 0, expectedCount: 5},  // 默认7，所以 p1, p2, p3, p4, p6
		{intensity: -1, expectedCount: 5}, // 默认7
		{intensity: 10, expectedCount: 6}, // 截断到9
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := FilterProbesByIntensity(probes, tt.intensity)
			if len(result) != tt.expectedCount {
				t.Errorf("FilterProbesByIntensity(intensity=%d): got %d probes, want %d",
					tt.intensity, len(result), tt.expectedCount)
			}
		})
	}
}

func TestGetSSLProbesForPort(t *testing.T) {
	InitializeGlobalVScan()
	v := GetGlobalVScan()

	// 测试SSL端口
	sslPorts := []int{443, 465, 636, 993, 995}
	for _, port := range sslPorts {
		probes := v.GetSSLProbesForPort(port)
		t.Logf("SSL端口 %d: 找到 %d 个SSL探测器", port, len(probes))
	}
}

func TestGetAllProbesSortedByRarity(t *testing.T) {
	InitializeGlobalVScan()
	v := GetGlobalVScan()

	probes := v.GetAllProbesSortedByRarity()
	if len(probes) == 0 {
		t.Fatal("GetAllProbesSortedByRarity 返回空列表")
	}

	t.Logf("总共 %d 个TCP探测器", len(probes))

	// 验证排序
	for i := 1; i < len(probes); i++ {
		prev := probes[i-1].Rarity
		curr := probes[i].Rarity
		if prev == 0 {
			prev = 10
		}
		if curr == 0 {
			curr = 10
		}
		if prev > curr {
			t.Errorf("探测器未按rarity排序: probes[%d].Rarity=%d > probes[%d].Rarity=%d",
				i-1, probes[i-1].Rarity, i, probes[i].Rarity)
		}
	}

	// 打印前10个探测器
	t.Log("前10个探测器（按rarity排序）:")
	for i := 0; i < 10 && i < len(probes); i++ {
		t.Logf("  [%d] %s (rarity=%d)", i+1, probes[i].Name, probes[i].Rarity)
	}
}

// TestFallbacksCompilation 验证 fallback 数组编译
func TestFallbacksCompilation(t *testing.T) {
	InitializeGlobalVScan()
	v := GetGlobalVScan()

	// 获取 NULL 探测器
	nullProbe, hasNull := v.ProbesMapKName["NULL"]
	if !hasNull {
		t.Fatal("NULL 探测器不存在")
	}

	// 验证 NULL 探测器的 fallback 只包含自身
	if nullProbe.Fallbacks[0] == nil {
		t.Error("NULL 探测器的 Fallbacks[0] 为 nil")
	} else if nullProbe.Fallbacks[0].Name != "NULL" {
		t.Errorf("NULL 探测器的 Fallbacks[0] 应该是自身，实际是 %s", nullProbe.Fallbacks[0].Name)
	}
	t.Log("✓ NULL 探测器的 fallback 只包含自身")

	// 验证 GetRequest 探测器（TCP，无 fallback 指令）
	getReq, hasGetReq := v.ProbesMapKName["GetRequest"]
	if hasGetReq {
		// fallbacks[0] 应该是自身
		if getReq.Fallbacks[0] == nil || getReq.Fallbacks[0].Name != "GetRequest" {
			t.Error("GetRequest 的 Fallbacks[0] 应该是自身")
		}
		// fallbacks[1] 应该是 NULL（TCP 探测器）
		if getReq.Protocol == "tcp" && getReq.Fallbacks[1] != nil {
			t.Logf("✓ GetRequest (TCP) 的 Fallbacks[1] = %s", getReq.Fallbacks[1].Name)
		}
	}

	// 统计有 fallback 数组的探测器数量
	countWithFallbacks := 0
	countWithNullFallback := 0
	for _, probe := range v.Probes {
		if probe.Fallbacks[0] != nil {
			countWithFallbacks++
		}
		// 检查 TCP 探测器是否有 NULL fallback
		if probe.Protocol == "tcp" {
			for i := 0; i < MaxFallbacks+1; i++ {
				if probe.Fallbacks[i] == nil {
					break
				}
				if probe.Fallbacks[i].Name == "NULL" {
					countWithNullFallback++
					break
				}
			}
		}
	}

	t.Logf("✓ %d 个探测器有 fallback 数组", countWithFallbacks)
	t.Logf("✓ %d 个 TCP 探测器有 NULL fallback", countWithNullFallback)
}

// TestFallbacksWithDirective 验证有 fallback 指令的探测器
func TestFallbacksWithDirective(t *testing.T) {
	InitializeGlobalVScan()
	v := GetGlobalVScan()

	// 查找有 fallback 指令的探测器
	for _, probe := range v.Probes {
		if probe.Fallback != "" {
			t.Logf("探测器 %s 有 fallback 指令: %s", probe.Name, probe.Fallback)

			// 验证 fallback 数组
			t.Logf("  Fallbacks 数组:")
			for i := 0; i < MaxFallbacks+1; i++ {
				if probe.Fallbacks[i] == nil {
					break
				}
				t.Logf("    [%d] %s", i, probe.Fallbacks[i].Name)
			}
		}
	}
}
