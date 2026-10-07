package fitness

import (
	"postapocgame/admin-server/services/iam/internal/consts"
	fitnessmodel "postapocgame/admin-server/services/iam/internal/model/fitness"
)

// TemplateAt 返回 date 当天生效的切换记录：生效日期 <= date 的最后一条，同一生效日期取 ID 最大。
// 切换记录只追加不修改，所以「某天用哪个模板」永远可以从历史推出来，切换不会改写已经过去的日子。
func TemplateAt(switches []fitnessmodel.FitnessTemplateSwitch, date string) (fitnessmodel.FitnessTemplateSwitch, bool) {
	var best fitnessmodel.FitnessTemplateSwitch
	found := false
	for _, s := range switches {
		if s.EffectiveDate > date {
			continue
		}
		if !found || s.EffectiveDate > best.EffectiveDate || (s.EffectiveDate == best.EffectiveDate && s.Id > best.Id) {
			best = s
			found = true
		}
	}
	return best, found
}

// PendingAfter 返回 today 之后才生效的最后一条切换（「明天起换成 XX」）。
func PendingAfter(switches []fitnessmodel.FitnessTemplateSwitch, today string) (fitnessmodel.FitnessTemplateSwitch, bool) {
	var best fitnessmodel.FitnessTemplateSwitch
	found := false
	for _, s := range switches {
		if s.EffectiveDate <= today {
			continue
		}
		if !found || s.EffectiveDate > best.EffectiveDate || (s.EffectiveDate == best.EffectiveDate && s.Id > best.Id) {
			best = s
			found = true
		}
	}
	return best, found
}

// EffectiveDateFor 切换生效日期：首次登录默认模板当天生效，其余（手动/采纳建议/管理员指定）一律次日生效，
// 保证今天已经开始执行、已经打过的卡对应的计划不会被中途换掉。
func EffectiveDateFor(source int64, today string) string {
	if source == consts.FitnessSwitchSourceDefault {
		return today
	}
	return AddDays(today, 1)
}

// ResolveDay 解析某模板某天的计划：该日期的覆盖优先，其次周模板对应星期几。
// rows 可以混有多个模板的数据，按 templateID 过滤。
func ResolveDay(rows []fitnessmodel.FitnessTemplateDay, templateID uint64, date string) (*fitnessmodel.FitnessTemplateDay, bool) {
	weekday := Weekday(date)
	var weekly *fitnessmodel.FitnessTemplateDay
	for i := range rows {
		row := &rows[i]
		if row.TemplateId != templateID {
			continue
		}
		if row.Weekday == 0 && row.PlanDate == date {
			return row, true
		}
		if row.Weekday == weekday && row.PlanDate == "" && weekly == nil {
			weekly = row
		}
	}
	return weekly, false
}

// PlannedTrainingType 返回某人某天计划的训练类型，没有模板/没有该天内容返回 0。
func PlannedTrainingType(switches []fitnessmodel.FitnessTemplateSwitch, rows []fitnessmodel.FitnessTemplateDay, date string) int64 {
	sw, ok := TemplateAt(switches, date)
	if !ok {
		return 0
	}
	day, _ := ResolveDay(rows, sw.ToTemplateId, date)
	if day == nil {
		return 0
	}
	return day.TrainingType
}
