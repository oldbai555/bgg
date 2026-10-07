package fitness

import (
	"testing"

	"postapocgame/admin-server/services/iam/internal/consts"
	fitnessmodel "postapocgame/admin-server/services/iam/internal/model/fitness"
)

func TestStatusOrDefault(t *testing.T) {
	cases := []struct {
		name     string
		status   int64
		fallback int64
		want     int64
	}{
		{"新建未传状态按启用", 0, consts.FitnessStatusEnabled, consts.FitnessStatusEnabled},
		{"编辑未传状态保持原来的禁用", 0, consts.FitnessStatusDisabled, consts.FitnessStatusDisabled},
		{"显式禁用", consts.FitnessStatusDisabled, consts.FitnessStatusEnabled, consts.FitnessStatusDisabled},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := statusOrDefault(c.status, c.fallback); got != c.want {
				t.Fatalf("statusOrDefault(%d, %d) = %d, want %d", c.status, c.fallback, got, c.want)
			}
		})
	}
}

func TestValidateTemplate_RejectsUnknownStatus(t *testing.T) {
	for _, status := range []int64{0, 3, -1} {
		tpl := &fitnessmodel.FitnessTemplate{Name: "测试", Status: status}
		if err := validateTemplate(tpl); err == nil {
			t.Fatalf("status=%d 应该校验失败", status)
		}
	}
	tpl := &fitnessmodel.FitnessTemplate{Name: "测试", Status: consts.FitnessStatusDisabled}
	if err := validateTemplate(tpl); err != nil {
		t.Fatalf("禁用状态应该合法: %v", err)
	}
}
