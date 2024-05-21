package util

import (
	"sort"
	"time"
)

// 精确time.Time到天
func TimeAccurateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// 精确time.Time到秒
func TimeAccurateToSecond(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
}

func Today() time.Time {
	return TimeAccurateToDay(time.Now())
}

func Yesterday() time.Time {
	return Today().AddDate(0, 0, -1)
}

func Tomorrow() time.Time {
	return Today().AddDate(0, 0, 1)
}

func GetWeekStartDateFromTime(t time.Time) time.Time {
	nowZeroDate := TimeAccurateToDay(t)
	weekDay := t.Weekday()
	weekDayNum := int(weekDay)
	if weekDay == time.Sunday {
		weekDayNum = 7
	}
	return nowZeroDate.AddDate(0, 0, -(weekDayNum - 1))
}

func GetWeekEndDateFromTime(t time.Time) time.Time {
	nowZeroDate := TimeAccurateToDay(t)
	weekDay := t.Weekday()
	weekDayNum := int(weekDay)
	if weekDay == time.Sunday {
		weekDayNum = 7
	}
	return nowZeroDate.AddDate(0, 0, (7 - weekDayNum))
}

func GetMonthStartDateFromTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func GetMonthEndDateFromTime(t time.Time) time.Time {
	var nextMonthStartDate time.Time
	if t.Month() == time.December {
		nextMonthStartDate = time.Date(t.Year()+1, time.January, 1, 0, 0, 0, 0, t.Location())
	} else {
		nextMonthStartDate = time.Date(t.Year(), time.Month(int(t.Month())+1), 1, 0, 0, 0, 0, t.Location())
	}
	return nextMonthStartDate.AddDate(0, 0, -1)
}

// Week number of the year (Monday as the first day of the week) as a decimal number [00,53].
// All days in a new year preceding the first Monday are considered to be in week 0.
func GetWeek(t time.Time) int {
	yday := t.YearDay()
	wday := int(t.Weekday()+6) % 7
	week := (yday - wday + 6) / 7

	return week
}

// 和主版本的周计算方法匹配
func GetMainVersionWeek(t time.Time) int {
	return GetWeek(t) + 1
}

// return 0-6 (Monday is 0 and Sunday is 6)
func MergeToPythonWeekDay(weekDay time.Weekday) int {
	day := 0
	if weekDay == time.Sunday {
		day = 7
	} else {
		day = int(weekDay)
	}
	return day - 1
}

// return 1-7 (Monday is 1 and Sunday is 7)
func MergeToPythonISOWeekDay(weekDay time.Weekday) int {
	day := 0
	if weekDay == time.Sunday {
		day = 7
	} else {
		day = int(weekDay)
	}
	return day
}

// TimeSlice attaches the methods of Interface to []time.Time, sorting in increasing order.
type TimeSlice []time.Time

func (t TimeSlice) Len() int {
	return len(t)
}
func (t TimeSlice) Less(i, j int) bool {
	return t[i].Before(t[j])
}
func (t TimeSlice) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t TimeSlice) Sort() {
	sort.Sort(t)
}

func MaxTime(times ...time.Time) time.Time {
	t := TimeSlice(times)
	re := sort.Reverse(t)
	sort.Sort(re)
	return t[0]
}

func MinTime(times ...time.Time) time.Time {
	t := TimeSlice(times)
	t.Sort()
	return t[0]
}

func PtrToTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

func TimeToPtr(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

func TimeNowPtr() *time.Time {
	now := time.Now()
	return &now
}
