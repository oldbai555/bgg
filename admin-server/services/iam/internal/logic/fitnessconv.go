package logic

import (
	fitnessdomain "postapocgame/admin-server/services/iam/internal/domain/fitness"
	fitnessmodel "postapocgame/admin-server/services/iam/internal/model/fitness"
	fitnessrepo "postapocgame/admin-server/services/iam/internal/repository/fitness"

	"postapocgame/admin-server/services/iam/iam"
)

// fitness 域 model/领域对象 → proto 的转换，Fitness*Logic 共用。

func fitnessTemplateItem(t fitnessmodel.FitnessTemplate) *iam.FitnessTemplateItem {
	return &iam.FitnessTemplateItem{
		Id: t.Id, Code: t.Code, Name: t.Name, Level: t.Level, Stage: t.Stage, Summary: t.Summary,
		DailyKcal: t.DailyKcal, DailyProtein: t.DailyProtein, Sort: t.Sort, Status: t.Status,
		CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt,
	}
}

func fitnessTemplateModel(in *iam.FitnessTemplateSaveRequest) *fitnessmodel.FitnessTemplate {
	return &fitnessmodel.FitnessTemplate{
		Id: in.Id, Code: in.Code, Name: in.Name, Level: in.Level, Stage: in.Stage, Summary: in.Summary,
		DailyKcal: in.DailyKcal, DailyProtein: in.DailyProtein, Sort: in.Sort, Status: in.Status,
	}
}

func fitnessTemplateBrief(t *fitnessmodel.FitnessTemplate) *iam.FitnessTemplateBrief {
	if t == nil {
		return nil
	}
	return &iam.FitnessTemplateBrief{Id: t.Id, Code: t.Code, Name: t.Name, Stage: t.Stage, Summary: t.Summary}
}

func fitnessDayContent(d *fitnessmodel.FitnessTemplateDay) *iam.FitnessDayContent {
	if d == nil {
		return nil
	}
	out := &iam.FitnessDayContent{
		TrainingType: d.TrainingType, TrainingTime: d.TrainingTime, StepGoal: d.StepGoal,
		Warmup: d.Warmup, Stretch: d.Stretch, TrainingNote: d.TrainingNote, Tip: d.Tip,
	}
	for _, e := range fitnessdomain.DecodeExercises(d.Exercises) {
		out.Exercises = append(out.Exercises, &iam.FitnessExercise{Name: e.Name, Sets: e.Sets, Reps: e.Reps, Duration: e.Duration, Note: e.Note})
	}
	for _, m := range fitnessdomain.DecodeMeals(d.Meals) {
		out.Meals = append(out.Meals, &iam.FitnessMeal{
			Slot: m.Slot, Time: m.Time, Place: m.Place, Food: m.Food, Portion: m.Portion,
			Kcal: m.Kcal, Protein: m.Protein, HowToOrder: m.HowToOrder, Alternatives: m.Alternatives,
		})
	}
	return out
}

func fitnessDayContentInput(c *iam.FitnessDayContent) fitnessdomain.DayContent {
	if c == nil {
		return fitnessdomain.DayContent{}
	}
	out := fitnessdomain.DayContent{
		TrainingType: c.TrainingType, TrainingTime: c.TrainingTime, StepGoal: c.StepGoal,
		Warmup: c.Warmup, Stretch: c.Stretch, TrainingNote: c.TrainingNote, Tip: c.Tip,
	}
	for _, e := range c.Exercises {
		out.Exercises = append(out.Exercises, fitnessdomain.Exercise{Name: e.Name, Sets: e.Sets, Reps: e.Reps, Duration: e.Duration, Note: e.Note})
	}
	for _, m := range c.Meals {
		out.Meals = append(out.Meals, fitnessdomain.Meal{
			Slot: m.Slot, Time: m.Time, Place: m.Place, Food: m.Food, Portion: m.Portion,
			Kcal: m.Kcal, Protein: m.Protein, HowToOrder: m.HowToOrder, Alternatives: m.Alternatives,
		})
	}
	return out
}

func fitnessTemplateDayItem(d fitnessmodel.FitnessTemplateDay) *iam.FitnessTemplateDayItem {
	return &iam.FitnessTemplateDayItem{
		Id: d.Id, TemplateId: d.TemplateId, Weekday: d.Weekday, PlanDate: d.PlanDate,
		Content: fitnessDayContent(&d), UpdatedAt: d.UpdatedAt,
	}
}

func fitnessTipItem(t fitnessmodel.FitnessTip) *iam.FitnessTipItem {
	return &iam.FitnessTipItem{
		Id: t.Id, Code: t.Code, Title: t.Title, Content: t.Content, Sort: t.Sort, Status: t.Status,
		CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt,
	}
}

func fitnessTipModel(in *iam.FitnessTipSaveRequest) *fitnessmodel.FitnessTip {
	return &fitnessmodel.FitnessTip{Id: in.Id, Code: in.Code, Title: in.Title, Content: in.Content, Sort: in.Sort, Status: in.Status}
}

func fitnessCheckinItem(c *fitnessmodel.FitnessCheckin) *iam.FitnessCheckinItem {
	if c == nil {
		return nil
	}
	out := &iam.FitnessCheckinItem{
		Date: c.CheckinDate, TemplateId: c.TemplateId, TrainingType: c.TrainingType, TrainingStatus: c.TrainingStatus,
		ExerciseDone: fitnessdomain.DecodeInts(c.ExerciseDone), ExerciseTotal: c.ExerciseTotal,
		Steps: c.Steps, StepGoal: c.StepGoal, WaterCups: c.WaterCups, UpdatedAt: c.UpdatedAt,
	}
	for _, m := range fitnessdomain.DecodeMealChecks(c.Meals) {
		out.Meals = append(out.Meals, &iam.FitnessMealCheck{Slot: m.Slot, Status: m.Status, Note: m.Note})
	}
	return out
}

func fitnessBodyItems(points []fitnessdomain.BodyPoint) []*iam.FitnessBodyRecordItem {
	out := make([]*iam.FitnessBodyRecordItem, 0, len(points))
	for _, p := range points {
		out = append(out, &iam.FitnessBodyRecordItem{Date: p.Date, WeightKg: p.WeightKg, WaistCm: p.WaistCm, WeightAvg7: p.WeightAvg7})
	}
	return out
}

func fitnessSuggestionItem(v *fitnessdomain.SuggestionView) *iam.FitnessSuggestionItem {
	if v == nil {
		return nil
	}
	s := v.Suggestion
	return &iam.FitnessSuggestionItem{
		Id: s.Id, WeekStart: s.WeekStart, RuleCode: s.RuleCode,
		FromTemplateId: s.FromTemplateId, FromTemplateName: v.FromName,
		ToTemplateId: s.ToTemplateId, ToTemplateName: v.ToName,
		Reason: s.Reason, Status: s.Status,
	}
}

func fitnessMemberItem(m fitnessrepo.MemberRow, sum *fitnessdomain.MemberSummary, tplNames map[uint64]string) *iam.FitnessMemberItem {
	out := &iam.FitnessMemberItem{
		UserId: m.UserId, Username: m.Username, Nickname: m.Nickname, Avatar: m.Avatar,
		HeightCm: m.HeightCm, TargetWeightKg: m.TargetWeightKg, CreatedAt: m.CreatedAt,
	}
	if sum != nil {
		out.TemplateId = sum.TemplateID
		out.TemplateName = tplNames[sum.TemplateID]
		out.LatestWeightKg = sum.LatestWeight
		out.StreakDays = int64(sum.StreakDays)
		out.WeekRate = sum.WeekRate
		out.LastCheckinDate = sum.LastCheckinDate
	}
	return out
}

func fitnessSwitchItem(v fitnessdomain.SwitchView) *iam.FitnessSwitchItem {
	s := v.Switch
	return &iam.FitnessSwitchItem{
		Id: s.Id, UserId: s.UserId, Nickname: v.Nickname,
		FromTemplateId: s.FromTemplateId, FromTemplateName: v.FromName,
		ToTemplateId: s.ToTemplateId, ToTemplateName: v.ToName,
		EffectiveDate: s.EffectiveDate, Source: s.Source, Reason: s.Reason,
		OperatorId: s.OperatorId, OperatorName: v.OperatorName, CreatedAt: s.CreatedAt,
	}
}
