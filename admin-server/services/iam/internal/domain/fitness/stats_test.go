package fitness

import (
	"math"
	"testing"

	"postapocgame/admin-server/services/iam/internal/consts"
	fitnessmodel "postapocgame/admin-server/services/iam/internal/model/fitness"
)

func checkin(date string, trainingType, status int64, done []int64, total int64) *fitnessmodel.FitnessCheckin {
	return &fitnessmodel.FitnessCheckin{
		CheckinDate:    date,
		TemplateId:     10,
		TrainingType:   trainingType,
		TrainingStatus: status,
		ExerciseDone:   EncodeJSON(done),
		ExerciseTotal:  total,
		Meals:          "[]",
	}
}

func almost(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestDayScore(t *testing.T) {
	strength := consts.FitnessTrainingStrength
	cases := []struct {
		name        string
		fact        DayFact
		wantScore   float64
		wantCounted bool
	}{
		{"完成", DayFact{Checkin: checkin("d", strength, consts.FitnessTrainingStatusDone, nil, 6)}, 1, true},
		{"部分完成按动作比例", DayFact{Checkin: checkin("d", strength, consts.FitnessTrainingStatusPartial, []int64{0, 1, 2}, 6)}, 0.5, true},
		{"部分完成没有动作清单按 0.5", DayFact{Checkin: checkin("d", consts.FitnessTrainingCardio, consts.FitnessTrainingStatusPartial, nil, 0)}, 0.5, true},
		{"跳过", DayFact{Checkin: checkin("d", strength, consts.FitnessTrainingStatusSkipped, nil, 6)}, 0, true},
		{"训练日没打卡", DayFact{PlannedType: strength}, 0, true},
		{"休息日不计入", DayFact{PlannedType: consts.FitnessTrainingRest}, 0, false},
		{"打卡快照是休息日，即使计划后来改成力量也不计入", DayFact{PlannedType: strength, Checkin: checkin("d", consts.FitnessTrainingRest, consts.FitnessTrainingStatusDone, nil, 0)}, 0, false},
		{"没有计划不计入", DayFact{}, 0, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			score, counted := DayScore(c.fact)
			if !almost(score, c.wantScore) || counted != c.wantCounted {
				t.Fatalf("DayScore = (%v, %v), want (%v, %v)", score, counted, c.wantScore, c.wantCounted)
			}
		})
	}
}

func TestCompletionRate(t *testing.T) {
	today := "2026-10-07" // 周三
	strength, rest := consts.FitnessTrainingStrength, consts.FitnessTrainingRest
	facts := []DayFact{
		{Date: "2026-10-05", PlannedType: strength, Checkin: checkin("2026-10-05", strength, consts.FitnessTrainingStatusDone, nil, 6)},
		{Date: "2026-10-06", PlannedType: rest},
		{Date: "2026-10-07", PlannedType: strength}, // 今天还没打训练卡：不计入
		{Date: "2026-10-08", PlannedType: strength}, // 未来：不计入
	}
	rate, counted := CompletionRate(facts, today)
	if counted != 1 || !almost(rate, 1) {
		t.Fatalf("今天未打卡不应拉低完成率, got rate=%v counted=%d", rate, counted)
	}

	facts[2].Checkin = checkin(today, strength, consts.FitnessTrainingStatusPartial, []int64{0, 1}, 4)
	rate, counted = CompletionRate(facts, today)
	if counted != 2 || !almost(rate, 0.75) {
		t.Fatalf("今天部分完成 0.5 → (1+0.5)/2, got rate=%v counted=%d", rate, counted)
	}

	missed := []DayFact{{Date: "2026-10-05", PlannedType: strength}, {Date: "2026-10-06", PlannedType: strength, Checkin: checkin("2026-10-06", strength, consts.FitnessTrainingStatusDone, nil, 6)}}
	if rate, counted := CompletionRate(missed, today); counted != 2 || !almost(rate, 0.5) {
		t.Fatalf("过去没打卡按 0 计, got rate=%v counted=%d", rate, counted)
	}

	if rate, counted := CompletionRate([]DayFact{{Date: "2026-10-06", PlannedType: rest}}, today); counted != 0 || rate != 0 {
		t.Fatalf("全是休息日应无数据, got rate=%v counted=%d", rate, counted)
	}
}

func TestIsActive(t *testing.T) {
	if IsActive(nil) {
		t.Fatal("没有打卡记录不算")
	}
	empty := &fitnessmodel.FitnessCheckin{Meals: "[]", ExerciseDone: "[]"}
	if IsActive(empty) {
		t.Fatal("全空记录不算打卡")
	}
	mealOnly := &fitnessmodel.FitnessCheckin{Meals: EncodeJSON([]MealCheck{{Slot: consts.FitnessMealBreakfast, Status: consts.FitnessMealStatusOnPlan}})}
	if !IsActive(mealOnly) {
		t.Fatal("只打了早餐也算打卡")
	}
	if !IsActive(&fitnessmodel.FitnessCheckin{WaterCups: 2, Meals: "[]"}) {
		t.Fatal("只打了喝水也算打卡")
	}
}

func TestStreak(t *testing.T) {
	today := "2026-10-07"
	active := map[string]bool{"2026-10-07": true, "2026-10-06": true, "2026-10-05": true, "2026-10-03": true}
	if got := Streak(active, today); got != 3 {
		t.Fatalf("今天打了卡从今天往回数, got %d", got)
	}
	delete(active, "2026-10-07")
	if got := Streak(active, today); got != 2 {
		t.Fatalf("今天还没打从昨天往回数, got %d", got)
	}
	if got := Streak(map[string]bool{"2026-10-05": true}, today); got != 0 {
		t.Fatalf("昨天断了就是 0, got %d", got)
	}
	cross := map[string]bool{"2027-01-01": true, "2026-12-31": true, "2026-12-30": true}
	if got := Streak(cross, "2027-01-01"); got != 3 {
		t.Fatalf("跨年连续, got %d", got)
	}
}

func TestMovingAverage(t *testing.T) {
	points := []WeightPoint{
		{Date: "2026-10-01", Weight: 80},
		{Date: "2026-10-02", Weight: 79},
		{Date: "2026-10-03", Weight: 0}, // 只填了腰围
		{Date: "2026-10-07", Weight: 78},
		{Date: "2026-10-08", Weight: 77}, // 窗口 [10-02, 10-08]，10-01 被挤出
	}
	got := MovingAverage(points)
	want := []float64{80, 79.5, 0, 79, 78}
	for i := range want {
		if !almost(got[i], want[i]) {
			t.Fatalf("MovingAverage[%d] = %v, want %v (all=%v)", i, got[i], want[i], got)
		}
	}
}

func TestAvgWeightBetween(t *testing.T) {
	points := []WeightPoint{{"2026-09-28", 80}, {"2026-09-30", 79}, {"2026-10-05", 78}, {"2026-10-06", 0}}
	if got := AvgWeightBetween(points, "2026-09-28", "2026-10-04"); !almost(got, 79.5) {
		t.Fatalf("got %v", got)
	}
	if got := AvgWeightBetween(points, "2026-10-12", "2026-10-18"); got != 0 {
		t.Fatalf("没有数据应为 0, got %v", got)
	}
}
