package fitness

import (
	"testing"
	"time"

	"postapocgame/admin-server/services/iam/internal/consts"
	fitnessmodel "postapocgame/admin-server/services/iam/internal/model/fitness"
)

func TestToday_UsesShanghaiTimezone(t *testing.T) {
	cases := []struct {
		name string
		utc  string
		want string
	}{
		{"UTC 前一天 16:00 已经是上海次日零点", "2026-10-06T16:00:00Z", "2026-10-07"},
		{"UTC 15:59 仍是上海当天 23:59", "2026-10-07T15:59:59Z", "2026-10-07"},
		{"跨年", "2026-12-31T16:30:00Z", "2027-01-01"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			now, err := time.Parse(time.RFC3339, c.utc)
			if err != nil {
				t.Fatal(err)
			}
			if got := Today(now); got != c.want {
				t.Fatalf("Today(%s) = %s, want %s", c.utc, got, c.want)
			}
			// 机器本地时区不影响结果
			if got := Today(now.In(time.FixedZone("PST", -8*3600))); got != c.want {
				t.Fatalf("Today in PST = %s, want %s", got, c.want)
			}
		})
	}
}

func TestWeekdayAndWeekStart(t *testing.T) {
	// 2026-10-05 是周一
	if got := Weekday("2026-10-05"); got != 1 {
		t.Fatalf("Weekday(周一) = %d", got)
	}
	if got := Weekday("2026-10-11"); got != 7 {
		t.Fatalf("Weekday(周日) = %d", got)
	}
	if got := WeekStart("2026-10-11"); got != "2026-10-05" {
		t.Fatalf("WeekStart(周日) = %s", got)
	}
	if got := WeekStart("2026-10-05"); got != "2026-10-05" {
		t.Fatalf("WeekStart(周一) = %s", got)
	}
}

func dayRow(id, tplID uint64, weekday int64, planDate string, trainingType int64) fitnessmodel.FitnessTemplateDay {
	return fitnessmodel.FitnessTemplateDay{Id: id, TemplateId: tplID, Weekday: weekday, PlanDate: planDate, TrainingType: trainingType}
}

func TestResolveDay_OverrideTakesPrecedence(t *testing.T) {
	rows := []fitnessmodel.FitnessTemplateDay{
		dayRow(1, 10, 3, "", consts.FitnessTrainingStrength), // 周三
		dayRow(2, 10, 0, "2026-10-07", consts.FitnessTrainingRest),
		dayRow(3, 20, 0, "2026-10-14", consts.FitnessTrainingCardio), // 别的模板的覆盖
	}

	day, isOverride := ResolveDay(rows, 10, "2026-10-07")
	if day == nil || day.Id != 2 || !isOverride {
		t.Fatalf("2026-10-07 应命中日期覆盖, got %+v override=%v", day, isOverride)
	}

	day, isOverride = ResolveDay(rows, 10, "2026-10-14")
	if day == nil || day.Id != 1 || isOverride {
		t.Fatalf("2026-10-14 没有本模板的覆盖，应回落周三周模板, got %+v override=%v", day, isOverride)
	}

	day, _ = ResolveDay(rows, 10, "2026-10-08")
	if day != nil {
		t.Fatalf("周四没有内容应返回 nil, got %+v", day)
	}
}

func sw(id, from, to uint64, effective string, source int64) fitnessmodel.FitnessTemplateSwitch {
	return fitnessmodel.FitnessTemplateSwitch{Id: id, UserId: 1, FromTemplateId: from, ToTemplateId: to, EffectiveDate: effective, Source: source}
}

func TestTemplateAt_NextDaySwitch(t *testing.T) {
	today := "2026-10-07"
	switches := []fitnessmodel.FitnessTemplateSwitch{
		sw(1, 0, 10, "2026-09-01", consts.FitnessSwitchSourceDefault),
	}
	// 今天手动切到 20：次日生效
	switches = append(switches, sw(2, 10, 20, EffectiveDateFor(consts.FitnessSwitchSourceManual, today), consts.FitnessSwitchSourceManual))

	if got, _ := TemplateAt(switches, today); got.ToTemplateId != 10 {
		t.Fatalf("今天仍应是原模板 10, got %d", got.ToTemplateId)
	}
	if got, _ := TemplateAt(switches, "2026-10-08"); got.ToTemplateId != 20 {
		t.Fatalf("明天应是新模板 20, got %d", got.ToTemplateId)
	}
	if got, _ := TemplateAt(switches, "2026-08-31"); got.Id != 0 {
		t.Fatalf("首次分配之前没有模板, got %+v", got)
	}
	pending, ok := PendingAfter(switches, today)
	if !ok || pending.ToTemplateId != 20 || pending.EffectiveDate != "2026-10-08" {
		t.Fatalf("应有一条明天生效的待切换, got %+v ok=%v", pending, ok)
	}
}

func TestTemplateAt_SameDayLatestWins(t *testing.T) {
	switches := []fitnessmodel.FitnessTemplateSwitch{
		sw(5, 10, 30, "2026-10-08", consts.FitnessSwitchSourceManual),
		sw(1, 0, 10, "2026-09-01", consts.FitnessSwitchSourceDefault),
		sw(7, 10, 20, "2026-10-08", consts.FitnessSwitchSourceAdmin),
	}
	if got, _ := TemplateAt(switches, "2026-10-09"); got.ToTemplateId != 20 {
		t.Fatalf("同一生效日多条取 ID 最大, got %d", got.ToTemplateId)
	}
}

func TestEffectiveDateFor(t *testing.T) {
	if got := EffectiveDateFor(consts.FitnessSwitchSourceDefault, "2026-10-07"); got != "2026-10-07" {
		t.Fatalf("首次默认模板当天生效, got %s", got)
	}
	for _, src := range []int64{consts.FitnessSwitchSourceManual, consts.FitnessSwitchSourceSuggestion, consts.FitnessSwitchSourceAdmin} {
		if got := EffectiveDateFor(src, "2026-12-31"); got != "2027-01-01" {
			t.Fatalf("source=%d 应次日生效, got %s", src, got)
		}
	}
}

func TestPlannedTrainingType_FollowsSwitchAndOverride(t *testing.T) {
	switches := []fitnessmodel.FitnessTemplateSwitch{
		sw(1, 0, 10, "2026-10-01", consts.FitnessSwitchSourceDefault),
		sw(2, 10, 20, "2026-10-08", consts.FitnessSwitchSourceManual),
	}
	rows := []fitnessmodel.FitnessTemplateDay{
		dayRow(1, 10, 3, "", consts.FitnessTrainingStrength),
		dayRow(2, 20, 4, "", consts.FitnessTrainingCardio),
		dayRow(3, 20, 0, "2026-10-15", consts.FitnessTrainingRest),
	}
	if got := PlannedTrainingType(switches, rows, "2026-10-07"); got != consts.FitnessTrainingStrength {
		t.Fatalf("10-07 周三用模板 10, got %d", got)
	}
	if got := PlannedTrainingType(switches, rows, "2026-10-08"); got != consts.FitnessTrainingCardio {
		t.Fatalf("10-08 周四用模板 20, got %d", got)
	}
	if got := PlannedTrainingType(switches, rows, "2026-10-15"); got != consts.FitnessTrainingRest {
		t.Fatalf("10-15 被日期覆盖成休息, got %d", got)
	}
	if got := PlannedTrainingType(switches, rows, "2026-09-30"); got != 0 {
		t.Fatalf("开始使用之前没有计划, got %d", got)
	}
}

func TestCanCheckinOn(t *testing.T) {
	today := "2026-10-07"
	cases := map[string]bool{
		"2026-10-07": true,
		"2026-09-30": true, // 往前 7 天
		"2026-09-29": false,
		"2026-10-08": false,
	}
	for date, want := range cases {
		if got := CanCheckinOn(date, today); got != want {
			t.Fatalf("CanCheckinOn(%s) = %v, want %v", date, got, want)
		}
	}
}
