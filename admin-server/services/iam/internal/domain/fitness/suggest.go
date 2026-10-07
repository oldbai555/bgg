package fitness

import (
	"fmt"

	"postapocgame/admin-server/services/iam/internal/consts"
)

// WeekRate 一周的训练完成率，Counted==0 表示这周没有可统计的训练日（还没开始用/整周休息），视为无数据。
type WeekRate struct {
	Rate    float64
	Counted int
}

// SuggestInput 生成「本周建议」需要的全部输入，按周一 WeekStart 往前看已经结束的整周。
type SuggestInput struct {
	CurrentCode string
	// CurrentSince 当前模板的生效日期，用来算「在当前模板上已经待了几周」
	CurrentSince string
	// PlateauReturnCode 进入平台期之前用的模板编码（平台期结束后回到它）
	PlateauReturnCode string
	WeekStart         string
	// WeeklyRates/WeeklyAvgWeights 时间正序，最后一个元素是上周；体重 0 表示那周没填
	WeeklyRates      []WeekRate
	WeeklyAvgWeights []float64
	LatestWeight     float64
	TargetWeight     float64
}

// Suggestion 命中的建议；ToCode 是目标模板编码，Reason 给用户看。
type Suggestion struct {
	RuleCode string
	ToCode   string
	Reason   string
}

var upgradeTo = map[string]string{
	consts.FitnessTemplateStarter:  consts.FitnessTemplateStandard,
	consts.FitnessTemplateStandard: consts.FitnessTemplateAdvanced,
	consts.FitnessTemplateEvening:  consts.FitnessTemplateAdvanced,
	consts.FitnessTemplateBusy:     consts.FitnessTemplateStandard,
}

var downgradeTo = map[string]string{
	consts.FitnessTemplateAdvanced: consts.FitnessTemplateStandard,
	consts.FitnessTemplateStandard: consts.FitnessTemplateStarter,
	consts.FitnessTemplateEvening:  consts.FitnessTemplateStarter,
	consts.FitnessTemplateStarter:  consts.FitnessTemplateBusy,
}

// Suggest 按优先级依次判断，命中第一条即返回，都不命中返回 nil。系统只建议不强切。
//  1. 体重达到目标 → 维持期
//  2. 平台期用满 N 周 → 回原模板
//  3. 连续 N 周完成率 < 阈值 → 降一档（入门再降就是忙碌周）
//  4. 连续 N 周完成率 >= 阈值 → 升一档
//  5. 周均体重连续 N 周不降 → 平台期
//
// 2~5 都要求在当前模板上已经待满对应周数，刚切过去的模板不会被立刻再建议切走。
func Suggest(in SuggestInput) *Suggestion {
	weeksOnCurrent := 0
	if in.CurrentSince != "" && in.CurrentSince <= in.WeekStart {
		weeksOnCurrent = DaysBetween(in.CurrentSince, in.WeekStart) / 7
	}

	if in.TargetWeight > 0 && in.LatestWeight > 0 && in.LatestWeight <= in.TargetWeight &&
		in.CurrentCode != consts.FitnessTemplateMaintain {
		return &Suggestion{
			RuleCode: consts.FitnessRuleTargetReached,
			ToCode:   consts.FitnessTemplateMaintain,
			Reason:   fmt.Sprintf("最近体重 %.1fkg 已达到目标 %.1fkg", in.LatestWeight, in.TargetWeight),
		}
	}

	if in.CurrentCode == consts.FitnessTemplatePlateau {
		if weeksOnCurrent >= consts.FitnessPlateauStayWeeks {
			back := in.PlateauReturnCode
			if back == "" || back == consts.FitnessTemplatePlateau {
				back = consts.FitnessTemplateStandard
			}
			return &Suggestion{
				RuleCode: consts.FitnessRulePlateauExit,
				ToCode:   back,
				Reason:   fmt.Sprintf("平台期调整已满 %d 周", consts.FitnessPlateauStayWeeks),
			}
		}
		return nil
	}

	if to, ok := downgradeTo[in.CurrentCode]; ok && weeksOnCurrent >= consts.FitnessDowngradeWeeks &&
		lastWeeksAll(in.WeeklyRates, consts.FitnessDowngradeWeeks, func(r float64) bool { return r < consts.FitnessDowngradeRate }) {
		return &Suggestion{
			RuleCode: consts.FitnessRuleDowngrade,
			ToCode:   to,
			Reason:   fmt.Sprintf("连续 %d 周训练完成率低于 %d%%", consts.FitnessDowngradeWeeks, int(consts.FitnessDowngradeRate*100)),
		}
	}

	if to, ok := upgradeTo[in.CurrentCode]; ok && weeksOnCurrent >= consts.FitnessUpgradeWeeks &&
		lastWeeksAll(in.WeeklyRates, consts.FitnessUpgradeWeeks, func(r float64) bool { return r >= consts.FitnessUpgradeRate }) {
		return &Suggestion{
			RuleCode: consts.FitnessRuleUpgrade,
			ToCode:   to,
			Reason:   fmt.Sprintf("连续 %d 周训练完成率达到 %d%% 以上", consts.FitnessUpgradeWeeks, int(consts.FitnessUpgradeRate*100)),
		}
	}

	if in.CurrentCode != consts.FitnessTemplateMaintain && weeksOnCurrent >= consts.FitnessPlateauWeeks &&
		weightNotDecreasing(in.WeeklyAvgWeights, consts.FitnessPlateauWeeks) {
		return &Suggestion{
			RuleCode: consts.FitnessRulePlateau,
			ToCode:   consts.FitnessTemplatePlateau,
			Reason:   fmt.Sprintf("周平均体重连续 %d 周没有下降", consts.FitnessPlateauWeeks),
		}
	}
	return nil
}

// lastWeeksAll 最近 n 周都有数据且都满足 cond。
func lastWeeksAll(rates []WeekRate, n int, cond func(float64) bool) bool {
	if len(rates) < n {
		return false
	}
	for _, w := range rates[len(rates)-n:] {
		if w.Counted == 0 || !cond(w.Rate) {
			return false
		}
	}
	return true
}

// weightNotDecreasing 最近 n 次周环比都没下降（需要 n+1 周的周均体重，任一周没数据不判定）。
func weightNotDecreasing(avgs []float64, n int) bool {
	if len(avgs) < n+1 {
		return false
	}
	recent := avgs[len(avgs)-n-1:]
	for _, w := range recent {
		if w <= 0 {
			return false
		}
	}
	for i := 1; i < len(recent); i++ {
		if round1(recent[i]) < round1(recent[i-1]) {
			return false
		}
	}
	return true
}
