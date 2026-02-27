package util

import (
	"errors"
	"time"
)

func GetNthWdoM(year int, month time.Month, loc *time.Location, wd time.Weekday, n int) (t time.Time, err error) {
	// n=0 is disallowed
	// n>4 is disallowed

	// n>0 returns the nth occurrence of the weekday in the month
	// n<0 returns the nth last occurrence of the weekday in the month

	if n == 0 || n > 4 || n < (-4) {
		err = errors.New("n must be between -4 and 4, excluding 0")
		return
	}

	var (
		orig time.Time
		do   int
	)

	if n > 0 {
		orig = time.Date(year, month, 1, 0, 0, 0, 0, loc)
		do = (int(orig.Weekday())+int(wd))%7 + 7*(n-1)
	} else {
		//last day of current month. Don't ask...
		orig = time.Date(year, month+1, 0, 0, 0, 0, 0, loc)
		do = -((int(orig.Weekday()) - int(wd)) % 7) + 7*(n+1)
	}
	t = orig.AddDate(0, 0, do)
	return
}
