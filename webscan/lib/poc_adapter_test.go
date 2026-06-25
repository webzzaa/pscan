package lib

import (
	"testing"
)

// TestDetectPocFormat 测试POC格式检测
func TestDetectPocFormat(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		expected PocFormat
	}{
		{
			name: "fscan格式 - 有name和rules",
			yaml: `
name: test-poc
rules:
  - method: GET
    path: /test
`,
			expected: FormatFscan,
		},
		{
			name: "fscan格式 - 有name和groups",
			yaml: `
name: test-poc
groups:
  group1:
    - method: GET
      path: /test
`,
			expected: FormatFscan,
		},
		{
			name: "Nuclei格式 - 有id和info",
			yaml: `
id: test-nuclei
info:
  name: Test Template
  author: test
  severity: info
http:
  - method: GET
    path:
      - "{{BaseURL}}"
`,
			expected: FormatNuclei,
		},
		{
			name: "未知格式",
			yaml: `
unknown: field
data: test
`,
			expected: FormatUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format := DetectPocFormat([]byte(tt.yaml))
			if format != tt.expected {
				t.Errorf("DetectPocFormat() = %v, want %v", format, tt.expected)
			}
		})
	}
}

// TestTskPocAdapter 测试fscan格式适配器
func TestTskPocAdapter(t *testing.T) {
	yaml := `
name: poc-yaml-test-fscan
set:
  rand: randomInt(10000, 99999)
rules:
  - method: GET
    path: /api/test
    expression: |
      response.status == 200
detail:
  author: test
  links:
    - https://example.com
`

	adapter, err := loadTskPoc([]byte(yaml))
	if err != nil {
		t.Fatalf("loadTskPoc() error = %v", err)
	}

	if adapter.GetFormat() != FormatFscan {
		t.Errorf("GetFormat() = %v, want %v", adapter.GetFormat(), FormatFscan)
	}

	if adapter.GetName() != "poc-yaml-test-fscan" {
		t.Errorf("GetName() = %v, want %v", adapter.GetName(), "poc-yaml-test-fscan")
	}

	poc, err := adapter.ToTskPoc()
	if err != nil {
		t.Fatalf("ToTskPoc() error = %v", err)
	}

	if poc.Name != "poc-yaml-test-fscan" {
		t.Errorf("Poc.Name = %v, want %v", poc.Name, "poc-yaml-test-fscan")
	}

	if len(poc.Rules) != 1 {
		t.Errorf("len(Poc.Rules) = %v, want %v", len(poc.Rules), 1)
	}

	if poc.Detail.Author != "test" {
		t.Errorf("Poc.Detail.Author = %v, want %v", poc.Detail.Author, "test")
	}
}

// TestNucleiPocAdapter 测试Nuclei格式适配器
func TestNucleiPocAdapter(t *testing.T) {
	yaml := `
id: test-nuclei-template
info:
  name: Test Nuclei Template
  author: pdteam
  severity: high
  description: Test template for nuclei adapter
  reference:
    - https://example.com/vuln
http:
  - method: GET
    path:
      - "{{BaseURL}}/admin"
      - "{{BaseURL}}/api"
    matchers:
      - type: word
        words:
          - "admin panel"
          - "dashboard"
      - type: status
        status:
          - 200
`

	adapter, err := loadNucleiPoc([]byte(yaml))
	if err != nil {
		t.Fatalf("loadNucleiPoc() error = %v", err)
	}

	if adapter.GetFormat() != FormatNuclei {
		t.Errorf("GetFormat() = %v, want %v", adapter.GetFormat(), FormatNuclei)
	}

	if adapter.GetName() != "Test Nuclei Template" {
		t.Errorf("GetName() = %v, want %v", adapter.GetName(), "Test Nuclei Template")
	}

	poc, err := adapter.ToTskPoc()
	if err != nil {
		t.Fatalf("ToTskPoc() error = %v", err)
	}

	if poc.Name != "Test Nuclei Template" {
		t.Errorf("Poc.Name = %v, want %v", poc.Name, "Test Nuclei Template")
	}

	// Nuclei的两个path应该转换为2个rule
	if len(poc.Rules) != 2 {
		t.Errorf("len(Poc.Rules) = %v, want %v", len(poc.Rules), 2)
	}

	if poc.Detail.Author != "pdteam" {
		t.Errorf("Poc.Detail.Author = %v, want %v", poc.Detail.Author, "pdteam")
	}

	// 验证expression包含word匹配
	if poc.Rules[0].Expression == "" {
		t.Error("Rule.Expression should not be empty")
	}
}

// TestConvertNucleiMatchers 测试Nuclei matcher转换
func TestConvertNucleiMatchers(t *testing.T) {
	tests := []struct {
		name              string
		matchers          []struct {
			Type      string   `yaml:"type"`
			Words     []string `yaml:"words"`
			Status    []int    `yaml:"status"`
			Regex     []string `yaml:"regex"`
			Condition string   `yaml:"condition"`
			Part      string   `yaml:"part"`
		}
		matchersCondition string
		wantContains      string
	}{
		{
			name: "单个word matcher",
			matchers: []struct {
				Type      string   `yaml:"type"`
				Words     []string `yaml:"words"`
				Status    []int    `yaml:"status"`
				Regex     []string `yaml:"regex"`
				Condition string   `yaml:"condition"`
				Part      string   `yaml:"part"`
			}{
				{
					Type:  "word",
					Words: []string{"admin"},
				},
			},
			matchersCondition: "",
			wantContains:      "response.body.bcontains",
		},
		{
			name: "单个status matcher",
			matchers: []struct {
				Type      string   `yaml:"type"`
				Words     []string `yaml:"words"`
				Status    []int    `yaml:"status"`
				Regex     []string `yaml:"regex"`
				Condition string   `yaml:"condition"`
				Part      string   `yaml:"part"`
			}{
				{
					Type:   "status",
					Status: []int{200},
				},
			},
			matchersCondition: "",
			wantContains:      "response.status == 200",
		},
		{
			name: "多个matcher - AND条件",
			matchers: []struct {
				Type      string   `yaml:"type"`
				Words     []string `yaml:"words"`
				Status    []int    `yaml:"status"`
				Regex     []string `yaml:"regex"`
				Condition string   `yaml:"condition"`
				Part      string   `yaml:"part"`
			}{
				{
					Type:  "word",
					Words: []string{"admin"},
				},
				{
					Type:   "status",
					Status: []int{200},
				},
			},
			matchersCondition: "",
			wantContains:      "&&",
		},
		{
			name: "多个matcher - OR条件",
			matchers: []struct {
				Type      string   `yaml:"type"`
				Words     []string `yaml:"words"`
				Status    []int    `yaml:"status"`
				Regex     []string `yaml:"regex"`
				Condition string   `yaml:"condition"`
				Part      string   `yaml:"part"`
			}{
				{
					Type:  "word",
					Words: []string{"admin"},
				},
				{
					Type:   "status",
					Status: []int{200},
				},
			},
			matchersCondition: "or",
			wantContains:      "||",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertNucleiMatchers(tt.matchers, tt.matchersCondition)
			if result == "" {
				t.Error("convertNucleiMatchers() returned empty string")
			}
			// 简单验证是否包含预期内容
			if tt.wantContains != "" {
				found := false
				for _, word := range []string{tt.wantContains} {
					if contains(result, word) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("convertNucleiMatchers() = %v, want to contain %v", result, tt.wantContains)
				}
			}
		})
	}
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && hasSubstring(s, substr))
}

func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestLoadUniversalPoc 测试通用加载器
func TestLoadUniversalPoc(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		yaml     string
		wantType PocFormat
		wantErr  bool
	}{
		{
			name:     "加载fscan格式",
			filename: "test-fscan.yml",
			yaml: `
name: test-fscan
rules:
  - method: GET
    path: /test
`,
			wantType: FormatFscan,
			wantErr:  false,
		},
		{
			name:     "加载nuclei格式",
			filename: "test-nuclei.yaml",
			yaml: `
id: test-nuclei
info:
  name: Test
http:
  - method: GET
    path:
      - "{{BaseURL}}"
`,
			wantType: FormatNuclei,
			wantErr:  false,
		},
		{
			name:     "未知格式报错",
			filename: "test-unknown.yml",
			yaml: `
unknown: format
`,
			wantType: FormatUnknown,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			poc, err := LoadUniversalPoc(tt.filename, []byte(tt.yaml))
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadUniversalPoc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if poc.GetFormat() != tt.wantType {
					t.Errorf("LoadUniversalPoc() format = %v, want %v", poc.GetFormat(), tt.wantType)
				}

				// 验证转换为fscan格式
				fscanPoc, err := poc.ToTskPoc()
				if err != nil {
					t.Errorf("ToTskPoc() error = %v", err)
				}
				if fscanPoc == nil {
					t.Error("ToTskPoc() returned nil")
				}
			}
		})
	}
}
