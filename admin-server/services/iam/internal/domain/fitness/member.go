package fitness

import (
	"context"
	"strings"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/internal/consts"
	fitnessmodel "postapocgame/admin-server/services/iam/internal/model/fitness"
)

// DayPlan 某人某天的计划 + 打卡。Day 为 nil 表示这天没有计划（没有模板或模板缺这一天）。
type DayPlan struct {
	Date       string
	Today      string
	CanCheckin bool
	Template   *fitnessmodel.FitnessTemplate
	Day        *fitnessmodel.FitnessTemplateDay
	IsOverride bool
	Tips       []fitnessmodel.FitnessTip
	Checkin    *fitnessmodel.FitnessCheckin
}

// CanCheckinOn 只能打今天和往前 FitnessBackfillDays 天的卡。
func CanCheckinOn(date, today string) bool {
	return date <= today && date >= AddDays(today, -consts.FitnessBackfillDays)
}

// ResolvePlan 解析某人某天的计划：切换历史 → 当天模板 → 日期覆盖优先，否则周模板。
func (s *Service) ResolvePlan(ctx context.Context, userID uint64, date string) (*fitnessmodel.FitnessTemplate, *fitnessmodel.FitnessTemplateDay, bool, error) {
	switches, err := s.Switch.ListByUsers(ctx, []uint64{userID})
	if err != nil {
		return nil, nil, false, dbErr(err)
	}
	sw, ok := TemplateAt(switches, date)
	if !ok {
		return nil, nil, false, nil
	}
	tpl, err := s.Template.FindByID(ctx, sw.ToTemplateId)
	if err != nil {
		return nil, nil, false, dbErr(err)
	}
	if tpl == nil {
		return nil, nil, false, nil
	}
	rows, err := s.TemplateDay.ListByTemplates(ctx, []uint64{tpl.Id}, date, date)
	if err != nil {
		return nil, nil, false, dbErr(err)
	}
	day, isOverride := ResolveDay(rows, tpl.Id, date)
	return tpl, day, isOverride, nil
}

func (s *Service) DayPlan(ctx context.Context, userID uint64, date string) (*DayPlan, error) {
	if _, err := s.EnsureMember(ctx, userID); err != nil {
		return nil, err
	}
	today := s.Today()
	if date == "" {
		date = today
	}
	if !ValidDate(date) {
		return nil, errs.New(errs.CodeBadRequest, "日期格式应为 YYYY-MM-DD")
	}
	tpl, day, isOverride, err := s.ResolvePlan(ctx, userID, date)
	if err != nil {
		return nil, err
	}
	tips, err := s.Tip.ListEnabled(ctx)
	if err != nil {
		return nil, dbErr(err)
	}
	checkin, err := s.Checkin.FindByUserDate(ctx, userID, date)
	if err != nil {
		return nil, dbErr(err)
	}
	return &DayPlan{
		Date:       date,
		Today:      today,
		CanCheckin: CanCheckinOn(date, today),
		Template:   tpl,
		Day:        day,
		IsOverride: isOverride,
		Tips:       tips,
		Checkin:    checkin,
	}, nil
}

// CheckinInput 手机端提交的一天打卡（整天覆盖式保存）。
type CheckinInput struct {
	Date           string
	TrainingStatus int64
	ExerciseDone   []int64
	Meals          []MealCheck
	Steps          int64
	WaterCups      int64
}

// SaveCheckin 保存打卡。模板/训练类型/动作数/步数目标在第一次打卡时快照，之后再编辑这天也不变——
// 次日生效的模板切换、后台改模板都不会改写已经打过的卡。
func (s *Service) SaveCheckin(ctx context.Context, userID uint64, in CheckinInput) (*fitnessmodel.FitnessCheckin, error) {
	if _, err := s.EnsureMember(ctx, userID); err != nil {
		return nil, err
	}
	today := s.Today()
	if in.Date == "" {
		in.Date = today
	}
	if !ValidDate(in.Date) {
		return nil, errs.New(errs.CodeBadRequest, "日期格式应为 YYYY-MM-DD")
	}
	if !CanCheckinOn(in.Date, today) {
		return nil, errs.New(errs.CodeBadRequest, "只能打今天和最近 7 天的卡")
	}
	if in.TrainingStatus < consts.FitnessTrainingStatusNone || in.TrainingStatus > consts.FitnessTrainingStatusSkipped {
		return nil, errs.New(errs.CodeBadRequest, "训练状态不合法")
	}
	if in.Steps < 0 || in.Steps > 200000 {
		return nil, errs.New(errs.CodeBadRequest, "步数不合法")
	}
	if in.WaterCups < 0 || in.WaterCups > 30 {
		return nil, errs.New(errs.CodeBadRequest, "喝水杯数不合法")
	}
	meals, err := normalizeMealChecks(in.Meals)
	if err != nil {
		return nil, err
	}

	existing, err := s.Checkin.FindByUserDate(ctx, userID, in.Date)
	if err != nil {
		return nil, dbErr(err)
	}
	record := existing
	if record == nil {
		record = &fitnessmodel.FitnessCheckin{UserId: userID, CheckinDate: in.Date, TrainingType: consts.FitnessTrainingRest}
		tpl, day, _, err := s.ResolvePlan(ctx, userID, in.Date)
		if err != nil {
			return nil, err
		}
		if tpl != nil {
			record.TemplateId = tpl.Id
		}
		if day != nil {
			record.TrainingType = day.TrainingType
			record.ExerciseTotal = int64(len(DecodeExercises(day.Exercises)))
			record.StepGoal = day.StepGoal
		}
	}

	done := normalizeExerciseDone(in.ExerciseDone, record.ExerciseTotal)
	status := in.TrainingStatus
	if status == consts.FitnessTrainingStatusNone && len(done) > 0 {
		status = consts.FitnessTrainingStatusPartial
		if int64(len(done)) == record.ExerciseTotal {
			status = consts.FitnessTrainingStatusDone
		}
	}
	record.TrainingStatus = status
	record.ExerciseDone = EncodeJSON(done)
	record.Meals = EncodeJSON(meals)
	record.Steps = in.Steps
	record.WaterCups = in.WaterCups

	if existing != nil {
		if err := s.Checkin.Update(ctx, record); err != nil {
			return nil, dbErr(err)
		}
	} else if err := s.Checkin.Create(ctx, record); err != nil {
		// 并发双击撞唯一键：按已存在的那条更新
		again, findErr := s.Checkin.FindByUserDate(ctx, userID, in.Date)
		if findErr != nil || again == nil {
			return nil, dbErr(err)
		}
		record.Id = again.Id
		record.TemplateId, record.TrainingType, record.ExerciseTotal, record.StepGoal = again.TemplateId, again.TrainingType, again.ExerciseTotal, again.StepGoal
		if err := s.Checkin.Update(ctx, record); err != nil {
			return nil, dbErr(err)
		}
	}
	return s.Checkin.FindByUserDate(ctx, userID, in.Date)
}

func normalizeMealChecks(in []MealCheck) ([]MealCheck, error) {
	valid := map[string]bool{}
	for _, slot := range consts.FitnessMealSlots {
		valid[slot] = true
	}
	seen := map[string]bool{}
	out := make([]MealCheck, 0, len(in))
	for _, m := range in {
		if !valid[m.Slot] {
			return nil, errs.New(errs.CodeBadRequest, "餐次不合法: "+m.Slot)
		}
		if m.Status < consts.FitnessMealStatusNone || m.Status > consts.FitnessMealStatusSkipped {
			return nil, errs.New(errs.CodeBadRequest, "餐食打卡状态不合法")
		}
		if seen[m.Slot] {
			continue
		}
		seen[m.Slot] = true
		note := strings.TrimSpace(m.Note)
		if m.Status != consts.FitnessMealStatusOther {
			note = ""
		}
		out = append(out, MealCheck{Slot: m.Slot, Status: m.Status, Note: truncate(note, 100)})
	}
	return out, nil
}

// normalizeExerciseDone 去重、去越界、排序（动作下标 0..total-1）。
func normalizeExerciseDone(in []int64, total int64) []int64 {
	seen := map[int64]bool{}
	out := make([]int64, 0, len(in))
	for i := int64(0); i < total; i++ {
		for _, v := range in {
			if v == i && !seen[v] {
				seen[v] = true
				out = append(out, v)
			}
		}
	}
	return out
}

// BodyPoint 身体数据的一天，WeightAvg7 是 7 日移动平均。
type BodyPoint struct {
	Date       string
	WeightKg   float64
	WaistCm    float64
	WeightAvg7 float64
}

// SaveBodyRecord 记录体重/腰围（同一天覆盖），日期窗口同打卡。
func (s *Service) SaveBodyRecord(ctx context.Context, userID uint64, date string, weight, waist float64) error {
	if _, err := s.EnsureMember(ctx, userID); err != nil {
		return err
	}
	today := s.Today()
	if date == "" {
		date = today
	}
	if !ValidDate(date) || !CanCheckinOn(date, today) {
		return errs.New(errs.CodeBadRequest, "只能记录今天和最近 7 天的数据")
	}
	if weight != 0 && (weight < 20 || weight > 300) {
		return errs.New(errs.CodeBadRequest, "体重需在 20~300kg 之间")
	}
	if waist != 0 && (waist < 30 || waist > 200) {
		return errs.New(errs.CodeBadRequest, "腰围需在 30~200cm 之间")
	}
	if weight == 0 && waist == 0 {
		return errs.New(errs.CodeBadRequest, "体重和腰围至少填一项")
	}
	existing, err := s.BodyRecord.FindByUserDate(ctx, userID, date)
	if err != nil {
		return dbErr(err)
	}
	if existing != nil {
		existing.WeightKg, existing.WaistCm = round1(weight), round1(waist)
		return dbErr(s.BodyRecord.Update(ctx, existing))
	}
	return dbErr(s.BodyRecord.Create(ctx, &fitnessmodel.FitnessBodyRecord{
		UserId: userID, RecordDate: date, WeightKg: round1(weight), WaistCm: round1(waist),
	}))
}

// BodyPoints [start, end] 内的身体数据；移动平均会往前多取 6 天，保证区间第一天的均值也是完整窗口。
func (s *Service) BodyPoints(ctx context.Context, userID uint64, start, end string) ([]BodyPoint, error) {
	records, err := s.BodyRecord.ListByUsersRange(ctx, []uint64{userID}, AddDays(start, -(consts.FitnessMovingAvgDays-1)), end)
	if err != nil {
		return nil, dbErr(err)
	}
	return bodyPointsFrom(records, start), nil
}

func bodyPointsFrom(records []fitnessmodel.FitnessBodyRecord, start string) []BodyPoint {
	points := make([]WeightPoint, len(records))
	for i, r := range records {
		points[i] = WeightPoint{Date: r.RecordDate, Weight: r.WeightKg}
	}
	avgs := MovingAverage(points)
	out := make([]BodyPoint, 0, len(records))
	for i, r := range records {
		if r.RecordDate < start {
			continue
		}
		out = append(out, BodyPoint{Date: r.RecordDate, WeightKg: r.WeightKg, WaistCm: r.WaistCm, WeightAvg7: avgs[i]})
	}
	return out
}

// MyBodyPoints 手机端：最近 days 天（默认 90，最多 365）。
func (s *Service) MyBodyPoints(ctx context.Context, userID uint64, days int64) ([]BodyPoint, error) {
	if _, err := s.EnsureMember(ctx, userID); err != nil {
		return nil, err
	}
	if days <= 0 {
		days = 90
	}
	if days > consts.FitnessBodyRecordMaxDays {
		days = consts.FitnessBodyRecordMaxDays
	}
	today := s.Today()
	return s.BodyPoints(ctx, userID, AddDays(today, -int(days-1)), today)
}

// MyProfile 「我的」页聚合数据。
type MyProfile struct {
	Today      string
	Member     *MemberSummary
	Profile    *fitnessmodel.FitnessProfile
	Nickname   string
	Avatar     string
	Current    *fitnessmodel.FitnessTemplate
	Pending    *fitnessmodel.FitnessTemplateSwitch
	PendingTpl *fitnessmodel.FitnessTemplate
	Suggestion *SuggestionView
	Templates  []fitnessmodel.FitnessTemplate
}

func (s *Service) MyProfile(ctx context.Context, userID uint64) (*MyProfile, error) {
	profile, err := s.EnsureMember(ctx, userID)
	if err != nil {
		return nil, err
	}
	today := s.Today()
	summaries, err := s.MemberSummaries(ctx, []uint64{userID}, today)
	if err != nil {
		return nil, err
	}
	tplMap, err := s.templateMap(ctx)
	if err != nil {
		return nil, err
	}
	switches, err := s.Switch.ListByUsers(ctx, []uint64{userID})
	if err != nil {
		return nil, dbErr(err)
	}
	out := &MyProfile{Today: today, Profile: profile, Member: summaries[userID]}
	if member, err := s.Profile.FindMember(ctx, userID); err == nil && member != nil {
		out.Nickname, out.Avatar = member.Nickname, member.Avatar
	}
	if cur, ok := TemplateAt(switches, today); ok {
		if t, ok := tplMap[cur.ToTemplateId]; ok {
			out.Current = &t
		}
	}
	if pending, ok := PendingAfter(switches, today); ok {
		out.Pending = &pending
		if t, ok := tplMap[pending.ToTemplateId]; ok {
			out.PendingTpl = &t
		}
	}
	for _, t := range tplMap {
		if t.Status == consts.FitnessStatusEnabled {
			out.Templates = append(out.Templates, t)
		}
	}
	sortTemplates(out.Templates)
	suggestion, err := s.CurrentSuggestion(ctx, userID)
	if err != nil {
		return nil, err
	}
	out.Suggestion = suggestion
	return out, nil
}

func (s *Service) UpdateProfile(ctx context.Context, userID uint64, height, target float64, preference int64) error {
	profile, err := s.EnsureMember(ctx, userID)
	if err != nil {
		return err
	}
	if height != 0 && (height < 100 || height > 250) {
		return errs.New(errs.CodeBadRequest, "身高需在 100~250cm 之间")
	}
	if target != 0 && (target < 30 || target > 300) {
		return errs.New(errs.CodeBadRequest, "目标体重需在 30~300kg 之间")
	}
	if preference != consts.FitnessPreferenceMorning && preference != consts.FitnessPreferenceEvening {
		return errs.New(errs.CodeBadRequest, "训练时间偏好不合法")
	}
	profile.HeightCm, profile.TargetWeightKg, profile.TrainingPreference = round1(height), round1(target), preference
	return dbErr(s.Profile.Update(ctx, profile))
}
