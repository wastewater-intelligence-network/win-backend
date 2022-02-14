package utils

import (
	"strconv"
	"strings"
	"time"
)

func splitDate(date string) (year int, month time.Month, day int, err error) {
	s := strings.Split(date, "-")
	year_, err := strconv.ParseInt(s[0], 10, 32)
	if err != nil {
		return
	}
	year = int(year_)

	month_, err := strconv.ParseInt(s[1], 10, 32)
	if err != nil {
		return
	}
	month = time.Month(month_)

	day_, err := strconv.ParseInt(s[2], 10, 32)
	if err != nil {
		return
	}
	day = int(day_)
	return
}

func GetDayTime(hour, min, sec, nsec int, date string) time.Time {
	var year, day int
	var month time.Month

	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		panic(err)
	}

	if date != "" {
		year, month, day, _ = splitDate(date)
	} else {
		t := time.Now().In(loc)
		year, month, day = t.Date()
	}

	return time.Date(year, month, day, hour, min, sec, nsec, loc)
}
