package fitness

import "time"

// Shanghai 是身材管理所有「今天/某天」的唯一口径。用 FixedZone 而不是 LoadLocation：
// 运行镜像不一定带 tzdata，而中国没有夏令时，+8 固定偏移与 Asia/Shanghai 完全等价。
var Shanghai = time.FixedZone("Asia/Shanghai", 8*3600)

const dateLayout = "2006-01-02"

// Today 返回 now 在上海时区的日期 YYYY-MM-DD。
func Today(now time.Time) string {
	return now.In(Shanghai).Format(dateLayout)
}

// ParseDate 解析 YYYY-MM-DD（上海时区零点）。
func ParseDate(date string) (time.Time, error) {
	return time.ParseInLocation(dateLayout, date, Shanghai)
}

// ValidDate 判断是否是合法的 YYYY-MM-DD。
func ValidDate(date string) bool {
	t, err := ParseDate(date)
	return err == nil && t.Format(dateLayout) == date
}

// AddDays 日期加减天数；非法日期原样返回空串。
func AddDays(date string, days int) string {
	t, err := ParseDate(date)
	if err != nil {
		return ""
	}
	return t.AddDate(0, 0, days).Format(dateLayout)
}

// Weekday 返回 1(周一)~7(周日)。
func Weekday(date string) int64 {
	t, err := ParseDate(date)
	if err != nil {
		return 0
	}
	wd := int64(t.Weekday())
	if wd == 0 {
		return 7
	}
	return wd
}

// WeekStart 返回 date 所在周的周一。
func WeekStart(date string) string {
	wd := Weekday(date)
	if wd == 0 {
		return ""
	}
	return AddDays(date, -int(wd-1))
}

// DaysBetween 返回 to - from 的天数。
func DaysBetween(from, to string) int {
	a, err1 := ParseDate(from)
	b, err2 := ParseDate(to)
	if err1 != nil || err2 != nil {
		return 0
	}
	return int(b.Sub(a).Hours() / 24)
}

// DateRange 返回 [start, end] 闭区间内的每一天，start > end 时返回空。
func DateRange(start, end string) []string {
	n := DaysBetween(start, end)
	if n < 0 || !ValidDate(start) {
		return nil
	}
	out := make([]string, 0, n+1)
	for i := 0; i <= n; i++ {
		out = append(out, AddDays(start, i))
	}
	return out
}
