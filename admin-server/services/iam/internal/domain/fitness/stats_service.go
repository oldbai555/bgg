package fitness

import (
	"context"
	"sort"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/internal/consts"
	fitnessmodel "postapocgame/admin-server/services/iam/internal/model/fitness"
	fitnessrepo "postapocgame/admin-server/services/iam/internal/repository/fitness"
)

// 统计口径集中在这一个文件：后台总览/明细/个人页/手机端「我的」都调这里，保证同一个人同一天的数字处处一致。
// 规模假设：使用者是公司内部人数级别（几十到几百），按人全量加载近一年打卡在内存里算，不做预聚合表。

// MemberSummary 某人截至某天的统计。
type MemberSummary struct {
	UserID          uint64
	TemplateID      uint64
	StreakDays      int
	WeekRate        float64
	WeekCounted     int
	LastCheckinDate string
	LatestWeight    float64
	ActiveOnDate    bool
}

// memberFacts 一批使用者在 [start, end] 的原始数据。
type memberFacts struct {
	switches map[uint64][]fitnessmodel.FitnessTemplateSwitch
	rows     []fitnessmodel.FitnessTemplateDay
	checkins map[uint64]map[string]*fitnessmodel.FitnessCheckin
	bodies   map[uint64][]fitnessmodel.FitnessBodyRecord
}

func (s *Service) loadFacts(ctx context.Context, userIDs []uint64, start, end string, withBody bool) (*memberFacts, error) {
	switches, err := s.Switch.ListByUsers(ctx, userIDs)
	if err != nil {
		return nil, dbErr(err)
	}
	rows, err := s.TemplateDay.ListByTemplates(ctx, templateIDsOf(switches), start, end)
	if err != nil {
		return nil, dbErr(err)
	}
	checkins, err := s.Checkin.ListByUsersRange(ctx, userIDs, start, end)
	if err != nil {
		return nil, dbErr(err)
	}
	f := &memberFacts{
		switches: groupSwitches(switches),
		rows:     rows,
		checkins: map[uint64]map[string]*fitnessmodel.FitnessCheckin{},
		bodies:   map[uint64][]fitnessmodel.FitnessBodyRecord{},
	}
	for i := range checkins {
		c := &checkins[i]
		if f.checkins[c.UserId] == nil {
			f.checkins[c.UserId] = map[string]*fitnessmodel.FitnessCheckin{}
		}
		f.checkins[c.UserId][c.CheckinDate] = c
	}
	if withBody {
		bodies, err := s.BodyRecord.ListByUsersRange(ctx, userIDs, AddDays(start, -(consts.FitnessMovingAvgDays-1)), end)
		if err != nil {
			return nil, dbErr(err)
		}
		for _, b := range bodies {
			f.bodies[b.UserId] = append(f.bodies[b.UserId], b)
		}
	}
	return f, nil
}

// dayFacts 某人在 dates 上的 DayFact。
func (f *memberFacts) dayFacts(userID uint64, dates []string) []DayFact {
	out := make([]DayFact, 0, len(dates))
	for _, d := range dates {
		out = append(out, DayFact{
			Date:        d,
			PlannedType: PlannedTrainingType(f.switches[userID], f.rows, d),
			Checkin:     f.checkins[userID][d],
		})
	}
	return out
}

func (f *memberFacts) summary(userID uint64, asOf string) *MemberSummary {
	sum := &MemberSummary{UserID: userID}
	if sw, ok := TemplateAt(f.switches[userID], asOf); ok {
		sum.TemplateID = sw.ToTemplateId
	}
	active := map[string]bool{}
	for date, c := range f.checkins[userID] {
		if date > asOf || !IsActive(c) {
			continue
		}
		active[date] = true
		if date > sum.LastCheckinDate {
			sum.LastCheckinDate = date
		}
	}
	sum.ActiveOnDate = active[asOf]
	sum.StreakDays = Streak(active, asOf)
	rate, counted := CompletionRate(f.dayFacts(userID, DateRange(WeekStart(asOf), asOf)), asOf)
	sum.WeekRate, sum.WeekCounted = Round2(rate), counted
	bodies := f.bodies[userID]
	for i := len(bodies) - 1; i >= 0; i-- {
		if bodies[i].RecordDate <= asOf && bodies[i].WeightKg > 0 {
			sum.LatestWeight = bodies[i].WeightKg
			break
		}
	}
	return sum
}

// MemberSummaries 一批使用者截至 asOf 的连续打卡/本周完成率/最新体重。
func (s *Service) MemberSummaries(ctx context.Context, userIDs []uint64, asOf string) (map[uint64]*MemberSummary, error) {
	out := map[uint64]*MemberSummary{}
	if len(userIDs) == 0 {
		return out, nil
	}
	facts, err := s.loadFacts(ctx, userIDs, AddDays(asOf, -consts.FitnessStreakLookbackDays), asOf, true)
	if err != nil {
		return nil, err
	}
	for _, id := range userIDs {
		out[id] = facts.summary(id, asOf)
	}
	return out, nil
}

// Overview 后台总览。
type Overview struct {
	Date              string
	MemberCount       int
	TodayCheckinCount int
	WeekAvgRate       float64
	Ranking           []RankItem
	DailyCounts       []DailyCount
}

type RankItem struct {
	Member  fitnessrepo.MemberRow
	Summary *MemberSummary
}

type DailyCount struct {
	Date  string
	Count int
}

func (s *Service) Overview(ctx context.Context, date string) (*Overview, error) {
	today := s.Today()
	if date == "" {
		date = today
	}
	if !ValidDate(date) {
		return nil, errs.New(errs.CodeBadRequest, "日期格式应为 YYYY-MM-DD")
	}
	members, err := s.Profile.ListMembers(ctx, nil)
	if err != nil {
		return nil, dbErr(err)
	}
	ids := make([]uint64, len(members))
	for i, m := range members {
		ids[i] = m.UserId
	}
	out := &Overview{Date: date, MemberCount: len(members)}
	facts := &memberFacts{}
	if len(ids) > 0 {
		facts, err = s.loadFacts(ctx, ids, AddDays(date, -consts.FitnessStreakLookbackDays), date, false)
		if err != nil {
			return nil, err
		}
	}

	var rateSum float64
	rateN := 0
	for _, m := range members {
		sum := facts.summary(m.UserId, date)
		if sum.ActiveOnDate {
			out.TodayCheckinCount++
		}
		if sum.WeekCounted > 0 {
			rateSum += sum.WeekRate
			rateN++
		}
		if sum.StreakDays > 0 {
			out.Ranking = append(out.Ranking, RankItem{Member: m, Summary: sum})
		}
	}
	if rateN > 0 {
		out.WeekAvgRate = Round2(rateSum / float64(rateN))
	}
	sort.SliceStable(out.Ranking, func(i, j int) bool {
		a, b := out.Ranking[i].Summary, out.Ranking[j].Summary
		if a.StreakDays != b.StreakDays {
			return a.StreakDays > b.StreakDays
		}
		return a.WeekRate > b.WeekRate
	})
	if len(out.Ranking) > consts.FitnessRankingSize {
		out.Ranking = out.Ranking[:consts.FitnessRankingSize]
	}
	for _, d := range DateRange(AddDays(date, -6), date) {
		n := 0
		for _, byDate := range facts.checkins {
			if IsActive(byDate[d]) {
				n++
			}
		}
		out.DailyCounts = append(out.DailyCounts, DailyCount{Date: d, Count: n})
	}
	return out, nil
}

// CheckinRow 后台打卡明细一行。Score < 0 表示这天不计入完成率（休息日）。
type CheckinRow struct {
	Checkin      fitnessmodel.FitnessCheckin
	Nickname     string
	TemplateName string
	ExerciseDone int64
	MealOnPlan   int64
	MealOther    int64
	MealSkipped  int64
	Score        float64
}

// ValidateRange 校验后台筛选的日期区间；都为空表示不限。
func ValidateRange(start, end string) error {
	if start != "" && !ValidDate(start) || end != "" && !ValidDate(end) {
		return errs.New(errs.CodeBadRequest, "日期格式应为 YYYY-MM-DD")
	}
	if start != "" && end != "" && start > end {
		return errs.New(errs.CodeBadRequest, "开始日期不能晚于结束日期")
	}
	return nil
}

func (s *Service) CheckinRows(ctx context.Context, filter fitnessrepo.CheckinFilter, page, pageSize int64) ([]CheckinRow, int64, error) {
	if err := ValidateRange(filter.StartDate, filter.EndDate); err != nil {
		return nil, 0, err
	}
	list, total, err := s.Checkin.FindPage(ctx, filter, page, pageSize)
	if err != nil {
		return nil, 0, dbErr(err)
	}
	tplMap, err := s.templateMap(ctx)
	if err != nil {
		return nil, 0, err
	}
	userIDs := make([]uint64, 0, len(list))
	for _, c := range list {
		userIDs = append(userIDs, c.UserId)
	}
	names, err := s.Profile.UserNicknames(ctx, userIDs)
	if err != nil {
		return nil, 0, dbErr(err)
	}
	rows := make([]CheckinRow, 0, len(list))
	for i := range list {
		c := list[i]
		row := CheckinRow{Checkin: c, Nickname: names[c.UserId], TemplateName: tplMap[c.TemplateId].Name, ExerciseDone: int64(len(DecodeInts(c.ExerciseDone)))}
		for _, m := range DecodeMealChecks(c.Meals) {
			switch m.Status {
			case consts.FitnessMealStatusOnPlan:
				row.MealOnPlan++
			case consts.FitnessMealStatusOther:
				row.MealOther++
			case consts.FitnessMealStatusSkipped:
				row.MealSkipped++
			}
		}
		score, counted := DayScore(DayFact{Date: c.CheckinDate, PlannedType: c.TrainingType, Checkin: &c})
		row.Score = -1
		if counted {
			row.Score = Round2(score)
		}
		rows = append(rows, row)
	}
	return rows, total, nil
}

// CalendarDay 个人明细热力图的一天。
type CalendarDay struct {
	Date         string
	Counted      bool
	Checked      bool
	Score        float64
	TrainingType int64
}

// MemberDetail 后台个人明细：热力图 + 体重趋势 + 切换记录。
type MemberDetail struct {
	Member    fitnessrepo.MemberRow
	Summary   *MemberSummary
	RangeRate float64
	Calendar  []CalendarDay
	Body      []BodyPoint
	Switches  []SwitchView
}

func (s *Service) MemberDetail(ctx context.Context, userID uint64, start, end string) (*MemberDetail, error) {
	today := s.Today()
	if end == "" {
		end = today
	}
	if start == "" {
		start = AddDays(end, -89)
	}
	if err := ValidateRange(start, end); err != nil {
		return nil, err
	}
	if DaysBetween(start, end) >= consts.FitnessStatsMaxRangeDays {
		return nil, errs.New(errs.CodeBadRequest, "日期跨度最多一年")
	}
	member, err := s.Profile.FindMember(ctx, userID)
	if err != nil {
		return nil, dbErr(err)
	}
	if member == nil {
		return nil, errs.New(errs.CodeNotFound, "使用者不存在")
	}
	summaries, err := s.MemberSummaries(ctx, []uint64{userID}, today)
	if err != nil {
		return nil, err
	}
	facts, err := s.loadFacts(ctx, []uint64{userID}, start, end, false)
	if err != nil {
		return nil, err
	}
	dayFacts := facts.dayFacts(userID, DateRange(start, end))
	out := &MemberDetail{Member: *member, Summary: summaries[userID]}
	out.RangeRate, _ = CompletionRate(dayFacts, today)
	out.RangeRate = Round2(out.RangeRate)
	for _, f := range dayFacts {
		score, counted := DayScore(f)
		out.Calendar = append(out.Calendar, CalendarDay{
			Date:         f.Date,
			Counted:      counted,
			Checked:      IsActive(f.Checkin),
			Score:        Round2(score),
			TrainingType: f.trainingType(),
		})
	}
	if out.Body, err = s.BodyPoints(ctx, userID, start, end); err != nil {
		return nil, err
	}
	switches, _, err := s.SwitchViews(ctx, userID, 1, 100)
	if err != nil {
		return nil, err
	}
	out.Switches = switches
	return out, nil
}

// SwitchView 切换记录 + 展示名。
type SwitchView struct {
	Switch       fitnessmodel.FitnessTemplateSwitch
	Nickname     string
	FromName     string
	ToName       string
	OperatorName string
}

func (s *Service) SwitchViews(ctx context.Context, userID uint64, page, pageSize int64) ([]SwitchView, int64, error) {
	list, total, err := s.Switch.FindPage(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, dbErr(err)
	}
	tplMap, err := s.templateMap(ctx)
	if err != nil {
		return nil, 0, err
	}
	var ids []uint64
	for _, sw := range list {
		ids = append(ids, sw.UserId)
		if sw.OperatorId > 0 {
			ids = append(ids, sw.OperatorId)
		}
	}
	names, err := s.Profile.UserNicknames(ctx, ids)
	if err != nil {
		return nil, 0, dbErr(err)
	}
	out := make([]SwitchView, 0, len(list))
	for _, sw := range list {
		out = append(out, SwitchView{
			Switch:       sw,
			Nickname:     names[sw.UserId],
			FromName:     tplMap[sw.FromTemplateId].Name,
			ToName:       tplMap[sw.ToTemplateId].Name,
			OperatorName: names[sw.OperatorId],
		})
	}
	return out, total, nil
}

// MemberPage 后台使用者列表（附带统计）。
func (s *Service) MemberPage(ctx context.Context, page, pageSize int64, keyword string) ([]fitnessrepo.MemberRow, map[uint64]*MemberSummary, int64, error) {
	members, total, err := s.Profile.FindMemberPage(ctx, page, pageSize, keyword)
	if err != nil {
		return nil, nil, 0, dbErr(err)
	}
	ids := make([]uint64, len(members))
	for i, m := range members {
		ids[i] = m.UserId
	}
	summaries, err := s.MemberSummaries(ctx, ids, s.Today())
	if err != nil {
		return nil, nil, 0, err
	}
	return members, summaries, total, nil
}

// TemplateNames 模板 ID → 名称。
func (s *Service) TemplateNames(ctx context.Context) (map[uint64]string, error) {
	tplMap, err := s.templateMap(ctx)
	if err != nil {
		return nil, err
	}
	out := make(map[uint64]string, len(tplMap))
	for id, t := range tplMap {
		out[id] = t.Name
	}
	return out, nil
}
