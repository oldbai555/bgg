package fitness

import (
	"postapocgame/admin-server/services/iam/internal/consts"
	fitnessmodel "postapocgame/admin-server/services/iam/internal/model/fitness"
)

// DayFact 统计用的一天：计划训练类型 + 当天打卡（没打卡 Checkin 为 nil）。
type DayFact struct {
	Date        string
	PlannedType int64
	Checkin     *fitnessmodel.FitnessCheckin
}

// trainingType 有打卡用打卡时的快照（模板后来被改也不影响历史），没打卡用计划。
func (f DayFact) trainingType() int64 {
	if f.Checkin != nil && f.Checkin.TemplateId > 0 {
		return f.Checkin.TrainingType
	}
	return f.PlannedType
}

// DayScore 训练完成度：完成=1，部分完成=已勾动作/动作总数（没有动作清单按 0.5），跳过/没打卡=0。
// 休息日和没有计划的日子不计入分母（counted=false）。
func DayScore(f DayFact) (score float64, counted bool) {
	tt := f.trainingType()
	if tt == 0 || tt == consts.FitnessTrainingRest {
		return 0, false
	}
	if f.Checkin == nil {
		return 0, true
	}
	switch f.Checkin.TrainingStatus {
	case consts.FitnessTrainingStatusDone:
		return 1, true
	case consts.FitnessTrainingStatusPartial:
		total := f.Checkin.ExerciseTotal
		if total <= 0 {
			return 0.5, true
		}
		done := int64(len(DecodeInts(f.Checkin.ExerciseDone)))
		if done > total {
			done = total
		}
		return float64(done) / float64(total), true
	default:
		return 0, true
	}
}

// CompletionRate 一段日子的训练完成率 = 训练日得分之和 / 训练日天数。
// today 当天还没打训练卡时不计入（一大早打开不应该看到完成率被今天拉低），过去的日子没打卡按 0 计。
func CompletionRate(facts []DayFact, today string) (rate float64, counted int) {
	var sum float64
	for _, f := range facts {
		if f.Date > today {
			continue
		}
		if f.Date == today && (f.Checkin == nil || f.Checkin.TrainingStatus == consts.FitnessTrainingStatusNone) {
			continue
		}
		score, ok := DayScore(f)
		if !ok {
			continue
		}
		sum += score
		counted++
	}
	if counted == 0 {
		return 0, 0
	}
	return sum / float64(counted), counted
}

// IsActive 当天算「打过卡」：训练、任意一餐、步数、喝水任一有记录。
func IsActive(c *fitnessmodel.FitnessCheckin) bool {
	if c == nil {
		return false
	}
	if c.TrainingStatus > 0 || c.Steps > 0 || c.WaterCups > 0 {
		return true
	}
	for _, m := range DecodeMealChecks(c.Meals) {
		if m.Status > 0 {
			return true
		}
	}
	return false
}

// Streak 连续打卡天数：今天打过卡从今天往回数；今天还没打从昨天往回数（今天还有机会，不算断）。
func Streak(active map[string]bool, today string) int {
	day := today
	if !active[day] {
		day = AddDays(today, -1)
	}
	n := 0
	for active[day] && n < consts.FitnessStreakLookbackDays {
		n++
		day = AddDays(day, -1)
	}
	return n
}

// WeightPoint 一天的体重（0 表示没填）。
type WeightPoint struct {
	Date   string
	Weight float64
}

// MovingAverage 每个点的 7 日移动平均：取 [d-6, d] 内所有填过体重的日子求平均，
// 没填体重的点返回 0。points 需要按日期正序。
func MovingAverage(points []WeightPoint) []float64 {
	out := make([]float64, len(points))
	for i, p := range points {
		if p.Weight <= 0 {
			continue
		}
		windowStart := AddDays(p.Date, -(consts.FitnessMovingAvgDays - 1))
		var sum float64
		n := 0
		for j := i; j >= 0 && points[j].Date >= windowStart; j-- {
			if points[j].Weight > 0 {
				sum += points[j].Weight
				n++
			}
		}
		out[i] = round1(sum / float64(n))
	}
	return out
}

// AvgWeightBetween [start, end] 内填过的体重平均值，没有数据返回 0。
func AvgWeightBetween(points []WeightPoint, start, end string) float64 {
	var sum float64
	n := 0
	for _, p := range points {
		if p.Weight > 0 && p.Date >= start && p.Date <= end {
			sum += p.Weight
			n++
		}
	}
	if n == 0 {
		return 0
	}
	return sum / float64(n)
}

func round1(v float64) float64 {
	if v < 0 {
		return -round1(-v)
	}
	return float64(int64(v*10+0.5)) / 10
}

// Round2 保留两位小数（完成率 0~1）。
func Round2(v float64) float64 {
	return float64(int64(v*100+0.5)) / 100
}
