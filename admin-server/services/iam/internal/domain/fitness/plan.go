package fitness

import "encoding/json"

// Exercise 动作清单的一项，存在 fitness_template_day.exercises（JSON 数组）。
// 组数×次数（Sets+Reps）与时长（Duration）二选一或组合，如 3×8-12、3×30-45秒、15分钟。
type Exercise struct {
	Name     string `json:"name"`
	Sets     int64  `json:"sets"`
	Reps     string `json:"reps"`
	Duration string `json:"duration"`
	Note     string `json:"note"`
}

// Meal 一餐的计划，存在 fitness_template_day.meals（JSON 数组）。
type Meal struct {
	Slot         string   `json:"slot"`
	Time         string   `json:"time"`
	Place        string   `json:"place"`
	Food         string   `json:"food"`
	Portion      string   `json:"portion"`
	Kcal         int64    `json:"kcal"`
	Protein      int64    `json:"protein"`
	HowToOrder   string   `json:"howToOrder"`
	Alternatives []string `json:"alternatives"`
}

// MealCheck 一餐的打卡，存在 fitness_checkin.meals（JSON 数组）。
type MealCheck struct {
	Slot   string `json:"slot"`
	Status int64  `json:"status"`
	Note   string `json:"note"`
}

// DecodeExercises/DecodeMeals/DecodeMealChecks/DecodeInts 解析 JSON 列；空串或脏数据按空列表处理，
// 不让一条坏数据把整天计划/整页统计打挂。
func DecodeExercises(raw string) []Exercise {
	var out []Exercise
	decodeJSON(raw, &out)
	return out
}

func DecodeMeals(raw string) []Meal {
	var out []Meal
	decodeJSON(raw, &out)
	return out
}

func DecodeMealChecks(raw string) []MealCheck {
	var out []MealCheck
	decodeJSON(raw, &out)
	return out
}

func DecodeInts(raw string) []int64 {
	var out []int64
	decodeJSON(raw, &out)
	return out
}

// EncodeJSON 序列化 JSON 列，nil 切片写成 []。
func EncodeJSON(v any) string {
	b, err := json.Marshal(v)
	if err != nil || string(b) == "null" {
		return "[]"
	}
	return string(b)
}

func decodeJSON(raw string, out any) {
	if raw == "" {
		return
	}
	_ = json.Unmarshal([]byte(raw), out)
}
