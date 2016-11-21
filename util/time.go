package util

import "time"

// 精确time.Time到天
func TimeAccurateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// 精确time.Time到秒
func TimeAccurateToSecond(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
}

func Yesterday() time.Time {
	return TimeAccurateToDay(time.Now()).AddDate(0, 0, -1)
}

func Tomorrow() time.Time {
	return TimeAccurateToDay(time.Now()).AddDate(0, 0, 1)
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

// return 0-6
func MergeToPythonWeekDay(weekDay time.Weekday) int {
	day := 0
	if weekDay == time.Sunday {
		day = 7
	} else {
		day = int(weekDay)
	}
	return day - 1
}
