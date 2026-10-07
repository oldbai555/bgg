// Package fitnessconv 身材管理 gateway 侧 iam-rpc 响应 ↔ HTTP types 的转换，fitness/* 各 logic 包共用。
package fitnessconv

import (
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/services/iam/iamclient"
)

func DayContent(c *iamclient.FitnessDayContent) *types.FitnessDayContent {
	if c == nil {
		return nil
	}
	out := &types.FitnessDayContent{
		TrainingType: c.TrainingType, TrainingTime: c.TrainingTime, StepGoal: c.StepGoal,
		Warmup: c.Warmup, Stretch: c.Stretch, TrainingNote: c.TrainingNote, Tip: c.Tip,
		Exercises: make([]types.FitnessExercise, 0, len(c.Exercises)),
		Meals:     make([]types.FitnessMeal, 0, len(c.Meals)),
	}
	for _, e := range c.Exercises {
		out.Exercises = append(out.Exercises, types.FitnessExercise{Name: e.Name, Sets: e.Sets, Reps: e.Reps, Duration: e.Duration, Note: e.Note})
	}
	for _, m := range c.Meals {
		alts := m.Alternatives
		if alts == nil {
			alts = []string{}
		}
		out.Meals = append(out.Meals, types.FitnessMeal{
			Slot: m.Slot, Time: m.Time, Place: m.Place, Food: m.Food, Portion: m.Portion,
			Kcal: m.Kcal, Protein: m.Protein, HowToOrder: m.HowToOrder, Alternatives: alts,
		})
	}
	return out
}

func DayContentToRPC(c *types.FitnessDayContent) *iamclient.FitnessDayContent {
	if c == nil {
		return nil
	}
	out := &iamclient.FitnessDayContent{
		TrainingType: c.TrainingType, TrainingTime: c.TrainingTime, StepGoal: c.StepGoal,
		Warmup: c.Warmup, Stretch: c.Stretch, TrainingNote: c.TrainingNote, Tip: c.Tip,
	}
	for _, e := range c.Exercises {
		out.Exercises = append(out.Exercises, &iamclient.FitnessExercise{Name: e.Name, Sets: e.Sets, Reps: e.Reps, Duration: e.Duration, Note: e.Note})
	}
	for _, m := range c.Meals {
		out.Meals = append(out.Meals, &iamclient.FitnessMeal{
			Slot: m.Slot, Time: m.Time, Place: m.Place, Food: m.Food, Portion: m.Portion,
			Kcal: m.Kcal, Protein: m.Protein, HowToOrder: m.HowToOrder, Alternatives: m.Alternatives,
		})
	}
	return out
}

func TipItem(t *iamclient.FitnessTipItem) types.FitnessTipItem {
	return types.FitnessTipItem{
		Id: t.Id, Code: t.Code, Title: t.Title, Content: t.Content, Sort: t.Sort, Status: t.Status,
		CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt,
	}
}

func TipItems(list []*iamclient.FitnessTipItem) []types.FitnessTipItem {
	out := make([]types.FitnessTipItem, 0, len(list))
	for _, t := range list {
		out = append(out, TipItem(t))
	}
	return out
}

func CheckinItem(c *iamclient.FitnessCheckinItem) *types.FitnessCheckinItem {
	if c == nil {
		return nil
	}
	out := &types.FitnessCheckinItem{
		Date: c.Date, TemplateId: c.TemplateId, TrainingType: c.TrainingType, TrainingStatus: c.TrainingStatus,
		ExerciseDone: c.ExerciseDone, ExerciseTotal: c.ExerciseTotal, Steps: c.Steps, StepGoal: c.StepGoal,
		WaterCups: c.WaterCups, UpdatedAt: c.UpdatedAt, Meals: make([]types.FitnessMealCheck, 0, len(c.Meals)),
	}
	if out.ExerciseDone == nil {
		out.ExerciseDone = []int64{}
	}
	for _, m := range c.Meals {
		out.Meals = append(out.Meals, types.FitnessMealCheck{Slot: m.Slot, Status: m.Status, Note: m.Note})
	}
	return out
}

func BodyItems(list []*iamclient.FitnessBodyRecordItem) []types.FitnessBodyRecordItem {
	out := make([]types.FitnessBodyRecordItem, 0, len(list))
	for _, b := range list {
		out = append(out, types.FitnessBodyRecordItem{Date: b.Date, WeightKg: b.WeightKg, WaistCm: b.WaistCm, WeightAvg7: b.WeightAvg7})
	}
	return out
}

func SuggestionItem(s *iamclient.FitnessSuggestionItem) *types.FitnessSuggestionItem {
	if s == nil {
		return nil
	}
	return &types.FitnessSuggestionItem{
		Id: s.Id, WeekStart: s.WeekStart, RuleCode: s.RuleCode,
		FromTemplateId: s.FromTemplateId, FromTemplateName: s.FromTemplateName,
		ToTemplateId: s.ToTemplateId, ToTemplateName: s.ToTemplateName,
		Reason: s.Reason, Status: s.Status,
	}
}

func TemplateBrief(t *iamclient.FitnessTemplateBrief) *types.FitnessTemplateBrief {
	if t == nil {
		return nil
	}
	return &types.FitnessTemplateBrief{Id: t.Id, Code: t.Code, Name: t.Name, Stage: t.Stage, Summary: t.Summary}
}

func TemplateItem(t *iamclient.FitnessTemplateItem) types.FitnessTemplateItem {
	return types.FitnessTemplateItem{
		Id: t.Id, Code: t.Code, Name: t.Name, Level: t.Level, Stage: t.Stage, Summary: t.Summary,
		DailyKcal: t.DailyKcal, DailyProtein: t.DailyProtein, Sort: t.Sort, Status: t.Status,
		CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt,
	}
}

func TemplateDayItems(list []*iamclient.FitnessTemplateDayItem) []types.FitnessTemplateDayItem {
	out := make([]types.FitnessTemplateDayItem, 0, len(list))
	for _, d := range list {
		out = append(out, types.FitnessTemplateDayItem{
			Id: d.Id, TemplateId: d.TemplateId, Weekday: d.Weekday, PlanDate: d.PlanDate,
			Content: DayContent(d.Content), UpdatedAt: d.UpdatedAt,
		})
	}
	return out
}

func MemberItem(m *iamclient.FitnessMemberItem) types.FitnessMemberItem {
	if m == nil {
		return types.FitnessMemberItem{}
	}
	return types.FitnessMemberItem{
		UserId: m.UserId, Username: m.Username, Nickname: m.Nickname, Avatar: m.Avatar,
		TemplateId: m.TemplateId, TemplateName: m.TemplateName, HeightCm: m.HeightCm,
		TargetWeightKg: m.TargetWeightKg, LatestWeightKg: m.LatestWeightKg, StreakDays: m.StreakDays,
		WeekRate: m.WeekRate, LastCheckinDate: m.LastCheckinDate, CreatedAt: m.CreatedAt,
	}
}

func SwitchItems(list []*iamclient.FitnessSwitchItem) []types.FitnessSwitchItem {
	out := make([]types.FitnessSwitchItem, 0, len(list))
	for _, s := range list {
		out = append(out, types.FitnessSwitchItem{
			Id: s.Id, UserId: s.UserId, Nickname: s.Nickname,
			FromTemplateId: s.FromTemplateId, FromTemplateName: s.FromTemplateName,
			ToTemplateId: s.ToTemplateId, ToTemplateName: s.ToTemplateName,
			EffectiveDate: s.EffectiveDate, Source: s.Source, Reason: s.Reason,
			OperatorId: s.OperatorId, OperatorName: s.OperatorName, CreatedAt: s.CreatedAt,
		})
	}
	return out
}
