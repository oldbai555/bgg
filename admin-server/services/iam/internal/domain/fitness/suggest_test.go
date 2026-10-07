package fitness

import (
	"testing"

	"postapocgame/admin-server/services/iam/internal/consts"
)

func rates(values ...float64) []WeekRate {
	out := make([]WeekRate, len(values))
	for i, v := range values {
		out[i] = WeekRate{Rate: v, Counted: 3}
	}
	return out
}

// weekStart 2026-10-05（周一）；since 2026-09-07 表示已经在当前模板上待满 4 周
const (
	testWeekStart = "2026-10-05"
	since4Weeks   = "2026-09-07"
	since2Weeks   = "2026-09-21"
	sinceLastWeek = "2026-09-28"
)

func TestSuggest(t *testing.T) {
	cases := []struct {
		name     string
		in       SuggestInput
		wantRule string
		wantTo   string
	}{
		{
			name:     "体重达标优先于一切 → 维持期",
			in:       SuggestInput{CurrentCode: consts.FitnessTemplateAdvanced, CurrentSince: since4Weeks, WeeklyRates: rates(0.9, 0.9, 0.9, 0.9), LatestWeight: 64.8, TargetWeight: 65},
			wantRule: consts.FitnessRuleTargetReached, wantTo: consts.FitnessTemplateMaintain,
		},
		{
			name: "已经是维持期不再建议维持期",
			in:   SuggestInput{CurrentCode: consts.FitnessTemplateMaintain, CurrentSince: since4Weeks, LatestWeight: 60, TargetWeight: 65, WeeklyAvgWeights: []float64{60, 60, 60}},
		},
		{
			name:     "连续 4 周 ≥80% → 入门升标准",
			in:       SuggestInput{CurrentCode: consts.FitnessTemplateStarter, CurrentSince: since4Weeks, WeeklyRates: rates(0.8, 0.85, 1, 0.9)},
			wantRule: consts.FitnessRuleUpgrade, wantTo: consts.FitnessTemplateStandard,
		},
		{
			name:     "标准升进阶",
			in:       SuggestInput{CurrentCode: consts.FitnessTemplateStandard, CurrentSince: since4Weeks, WeeklyRates: rates(0.8, 0.8, 0.8, 0.8)},
			wantRule: consts.FitnessRuleUpgrade, wantTo: consts.FitnessTemplateAdvanced,
		},
		{
			name: "4 周里有一周 79% 不升",
			in:   SuggestInput{CurrentCode: consts.FitnessTemplateStarter, CurrentSince: since4Weeks, WeeklyRates: rates(0.9, 0.79, 0.9, 0.9)},
		},
		{
			name: "某周没有数据不升",
			in:   SuggestInput{CurrentCode: consts.FitnessTemplateStarter, CurrentSince: since4Weeks, WeeklyRates: []WeekRate{{Rate: 0, Counted: 0}, {1, 3}, {1, 3}, {1, 3}}},
		},
		{
			name: "刚切过来不满 4 周不升",
			in:   SuggestInput{CurrentCode: consts.FitnessTemplateStarter, CurrentSince: since2Weeks, WeeklyRates: rates(0.9, 0.9, 0.9, 0.9)},
		},
		{
			name: "进阶已经是最高档不升",
			in:   SuggestInput{CurrentCode: consts.FitnessTemplateAdvanced, CurrentSince: since4Weeks, WeeklyRates: rates(1, 1, 1, 1)},
		},
		{
			name:     "连续 2 周 <50% → 标准降入门",
			in:       SuggestInput{CurrentCode: consts.FitnessTemplateStandard, CurrentSince: since2Weeks, WeeklyRates: rates(0.9, 0.9, 0.4, 0.3)},
			wantRule: consts.FitnessRuleDowngrade, wantTo: consts.FitnessTemplateStarter,
		},
		{
			name:     "入门再降 → 忙碌周",
			in:       SuggestInput{CurrentCode: consts.FitnessTemplateStarter, CurrentSince: since4Weeks, WeeklyRates: rates(0.2, 0.2, 0.2, 0.2)},
			wantRule: consts.FitnessRuleDowngrade, wantTo: consts.FitnessTemplateBusy,
		},
		{
			name: "只有 1 周 <50% 不降",
			in:   SuggestInput{CurrentCode: consts.FitnessTemplateStandard, CurrentSince: since4Weeks, WeeklyRates: rates(0.9, 0.9, 0.6, 0.3)},
		},
		{
			name: "50% 整不算低于",
			in:   SuggestInput{CurrentCode: consts.FitnessTemplateStandard, CurrentSince: since4Weeks, WeeklyRates: rates(0.9, 0.9, 0.5, 0.5)},
		},
		{
			name:     "周均体重连续 2 周不降 → 平台期",
			in:       SuggestInput{CurrentCode: consts.FitnessTemplateStandard, CurrentSince: since4Weeks, WeeklyRates: rates(0.6, 0.6, 0.6, 0.6), WeeklyAvgWeights: []float64{0, 75.2, 75.2, 75.4}},
			wantRule: consts.FitnessRulePlateau, wantTo: consts.FitnessTemplatePlateau,
		},
		{
			name: "体重还在降不进平台期",
			in:   SuggestInput{CurrentCode: consts.FitnessTemplateStandard, CurrentSince: since4Weeks, WeeklyAvgWeights: []float64{0, 76, 75.5, 75.6}},
		},
		{
			name: "缺一周体重不判定平台期",
			in:   SuggestInput{CurrentCode: consts.FitnessTemplateStandard, CurrentSince: since4Weeks, WeeklyAvgWeights: []float64{0, 75, 0, 75}},
		},
		{
			name:     "平台期满 2 周 → 回原模板",
			in:       SuggestInput{CurrentCode: consts.FitnessTemplatePlateau, CurrentSince: since2Weeks, PlateauReturnCode: consts.FitnessTemplateEvening, WeeklyAvgWeights: []float64{75, 75, 75, 75}},
			wantRule: consts.FitnessRulePlateauExit, wantTo: consts.FitnessTemplateEvening,
		},
		{
			name: "平台期才 1 周不动，也不会重复建议平台期",
			in:   SuggestInput{CurrentCode: consts.FitnessTemplatePlateau, CurrentSince: sinceLastWeek, WeeklyAvgWeights: []float64{75, 75, 75, 75}},
		},
		{
			name:     "平台期原模板未知时回标准",
			in:       SuggestInput{CurrentCode: consts.FitnessTemplatePlateau, CurrentSince: since4Weeks},
			wantRule: consts.FitnessRulePlateauExit, wantTo: consts.FitnessTemplateStandard,
		},
		{
			name:     "降档优先于平台期",
			in:       SuggestInput{CurrentCode: consts.FitnessTemplateAdvanced, CurrentSince: since4Weeks, WeeklyRates: rates(0.9, 0.9, 0.1, 0.1), WeeklyAvgWeights: []float64{0, 75, 75, 75}},
			wantRule: consts.FitnessRuleDowngrade, wantTo: consts.FitnessTemplateStandard,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.in.WeekStart = testWeekStart
			got := Suggest(c.in)
			if c.wantRule == "" {
				if got != nil {
					t.Fatalf("不应命中规则, got %+v", got)
				}
				return
			}
			if got == nil || got.RuleCode != c.wantRule || got.ToCode != c.wantTo {
				t.Fatalf("Suggest = %+v, want rule=%s to=%s", got, c.wantRule, c.wantTo)
			}
			if got.Reason == "" {
				t.Fatal("建议理由不能为空")
			}
		})
	}
}
