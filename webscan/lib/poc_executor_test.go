package lib

import (
	"testing"
)

// =============================================================================
// 重构后函数的单元测试
// =============================================================================

// TestGetRuleHash 测试规则哈希计算
func TestGetRuleHash(t *testing.T) {
	tests := []struct {
		name      string
		rule1     *Rules
		rule2     *Rules
		wantSame  bool
		wantEmpty bool
	}{
		{
			name: "相同规则产生相同哈希",
			rule1: &Rules{
				Method:  "GET",
				Path:    "/api/test",
				Headers: map[string]string{"User-Agent": "test"},
				Body:    "",
			},
			rule2: &Rules{
				Method:  "GET",
				Path:    "/api/test",
				Headers: map[string]string{"User-Agent": "test"},
				Body:    "",
			},
			wantSame: true,
		},
		{
			name: "不同Method产生不同哈希",
			rule1: &Rules{
				Method:  "GET",
				Path:    "/api/test",
				Headers: map[string]string{},
				Body:    "",
			},
			rule2: &Rules{
				Method:  "POST",
				Path:    "/api/test",
				Headers: map[string]string{},
				Body:    "",
			},
			wantSame: false,
		},
		{
			name: "不同Path产生不同哈希",
			rule1: &Rules{
				Method:  "GET",
				Path:    "/api/test1",
				Headers: map[string]string{},
				Body:    "",
			},
			rule2: &Rules{
				Method:  "GET",
				Path:    "/api/test2",
				Headers: map[string]string{},
				Body:    "",
			},
			wantSame: false,
		},
		{
			name: "空规则产生非空哈希",
			rule1: &Rules{
				Method:  "",
				Path:    "",
				Headers: map[string]string{},
				Body:    "",
			},
			rule2:     nil,
			wantSame:  false,
			wantEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash1 := getRuleHash(tt.rule1)

			// 验证哈希非空（MD5应该是32个十六进制字符）
			if !tt.wantEmpty && len(hash1) != 32 {
				t.Errorf("getRuleHash() 返回的哈希长度不正确，期望32，实际%d", len(hash1))
			}

			// 验证哈希是否相同
			if tt.rule2 != nil {
				hash2 := getRuleHash(tt.rule2)
				areSame := hash1 == hash2
				if areSame != tt.wantSame {
					t.Errorf("getRuleHash() 哈希相同性不符合预期\n规则1哈希: %s\n规则2哈希: %s\n期望相同: %v\n实际相同: %v",
						hash1, hash2, tt.wantSame, areSame)
				}
			}
		})
	}
}

// TestApplyParametersToRule 测试参数替换逻辑
func TestApplyParametersToRule(t *testing.T) {
	tests := []struct {
		name            string
		rule            Rules
		sets            ListMap
		payloads        map[string]interface{}
		variableMap     map[string]interface{}
		payloadExpr     string
		wantReplacement bool
		wantPath        string
		wantHeader      string
		wantBody        string
	}{
		{
			name: "替换Path中的参数",
			rule: Rules{
				Method:  "GET",
				Path:    "/api/{{key}}/test",
				Headers: map[string]string{},
				Body:    "",
			},
			sets: ListMap{
				{Key: "key", Value: []string{"value1"}},
			},
			payloads: map[string]interface{}{
				"key": "myvalue",
			},
			variableMap:     map[string]interface{}{},
			payloadExpr:     "",
			wantReplacement: true,
			wantPath:        "/api/myvalue/test",
		},
		{
			name: "替换Header中的参数",
			rule: Rules{
				Method:  "GET",
				Path:    "/test",
				Headers: map[string]string{"X-Custom": "{{token}}"},
				Body:    "",
			},
			sets: ListMap{
				{Key: "token", Value: []string{"abc123"}},
			},
			payloads: map[string]interface{}{
				"token": "secret123",
			},
			variableMap:     map[string]interface{}{},
			payloadExpr:     "",
			wantReplacement: true,
			wantHeader:      "secret123",
		},
		{
			name: "替换Body中的参数",
			rule: Rules{
				Method:  "POST",
				Path:    "/api/login",
				Headers: map[string]string{},
				Body:    `{"username":"{{user}}","password":"{{pass}}"}`,
			},
			sets: ListMap{
				{Key: "user", Value: []string{"admin"}},
				{Key: "pass", Value: []string{"123456"}},
			},
			payloads: map[string]interface{}{
				"user": "testuser",
				"pass": "testpass",
			},
			variableMap:     map[string]interface{}{},
			payloadExpr:     "",
			wantReplacement: true,
			wantBody:        `{"username":"testuser","password":"testpass"}`,
		},
		{
			name: "无匹配参数时不替换",
			rule: Rules{
				Method:  "GET",
				Path:    "/static/page",
				Headers: map[string]string{},
				Body:    "",
			},
			sets: ListMap{
				{Key: "key", Value: []string{"value"}},
			},
			payloads: map[string]interface{}{
				"key": "test",
			},
			variableMap:     map[string]interface{}{},
			payloadExpr:     "",
			wantReplacement: false,
			wantPath:        "/static/page",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 克隆规则避免修改原始数据
			currentRule := Rules{
				Method:  tt.rule.Method,
				Path:    tt.rule.Path,
				Headers: make(map[string]string),
				Body:    tt.rule.Body,
			}
			for k, v := range tt.rule.Headers {
				currentRule.Headers[k] = v
			}

			hasReplacement, replacedParams := applyParametersToRule(
				&currentRule,
				tt.sets,
				tt.payloads,
				tt.variableMap,
				tt.payloadExpr,
			)

			// 验证是否发生替换
			if hasReplacement != tt.wantReplacement {
				t.Errorf("applyParametersToRule() hasReplacement = %v, want %v", hasReplacement, tt.wantReplacement)
			}

			// 验证替换后的值
			if tt.wantPath != "" && currentRule.Path != tt.wantPath {
				t.Errorf("applyParametersToRule() Path = %q, want %q", currentRule.Path, tt.wantPath)
			}

			if tt.wantHeader != "" {
				found := false
				for _, v := range currentRule.Headers {
					if v == tt.wantHeader {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("applyParametersToRule() Header 未找到期望值 %q", tt.wantHeader)
				}
			}

			if tt.wantBody != "" && currentRule.Body != tt.wantBody {
				t.Errorf("applyParametersToRule() Body = %q, want %q", currentRule.Body, tt.wantBody)
			}

			// 验证替换参数列表
			if hasReplacement && len(replacedParams) == 0 {
				t.Error("applyParametersToRule() 有替换但replacedParams为空")
			}
		})
	}
}
