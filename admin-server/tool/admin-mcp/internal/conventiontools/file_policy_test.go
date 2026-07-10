package conventiontools

import "testing"

func TestFilePolicyMatching(t *testing.T) {
	rules, err := loadFilePolicyRules("../../data")
	if err != nil {
		t.Fatalf("loadFilePolicyRules() error = %v", err)
	}

	cases := []struct {
		path       string
		wantPolicy string
	}{
		{"internal/handler/iam/user/user_create_handler.go", "forbidden_generated"},
		{"internal/handler/iam/user/custom_routes.go", "normal_business_code"}, // 命中 exceptions，跳过该规则，落到兜底规则
		{"internal/repository/iam/user_repository.go", "editable_handwritten"},
		{"internal/model/iam/adminusermodel_gen.go", "forbidden_generated"},
		{"internal/model/iam/adminusermodel.go", "editable_sibling"},
		{"internal/domain/iam/user_service.go", "editable_handwritten"},
		{"/etc/work/mysql.json", "stop_and_ask"},
		{"admin-server/internal/logic/iam/user/user_create_logic.go", "editable_generated_skeleton"},
		{"some/random/business/file.go", "normal_business_code"},
	}

	for _, c := range cases {
		normalized := normalizeRepoPath(c.path)
		got := matchPolicy(rules, normalized)
		if got != c.wantPolicy {
			t.Errorf("path=%q normalized=%q policy = %q, want %q", c.path, normalized, got, c.wantPolicy)
		}
	}
}

// matchPolicy 复用 query_file_policy 的匹配顺序逻辑，供测试直接调用，不经过 MCP 调用层。
func matchPolicy(rules *filePolicyRules, normalized string) string {
	for _, rule := range rules.Rules {
		matched, err := doublestarMatch(rule.Glob, normalized)
		if err != nil || !matched {
			continue
		}
		if isException(rule.Exceptions, normalized) {
			continue
		}
		return rule.Policy
	}
	return ""
}
