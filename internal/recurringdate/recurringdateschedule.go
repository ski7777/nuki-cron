package recurringdate

import (
	"errors"
	"time"

	"github.com/ski7777/nuki-cron/internal/util"
)

type RecurringDateSchedule struct {
	wd time.Weekday
	n  int
}

func (rds RecurringDateSchedule) GetByYearMonth(year int, month time.Month) (t time.Time, err error) {
	return util.GetNthWdoM(year, month, time.Local, rds.wd, rds.n)
}

func (rds RecurringDateSchedule) GetByTime(t time.Time) (time.Time, error) {
	return util.GetNthWdoM(t.Year(), t.Month(), t.Location(), rds.wd, rds.n)
}

func (rds RecurringDateSchedule) GetNow() (time.Time, error) {
	return rds.GetByTime(time.Now())
}

func NewReccuringDateSchedule(wd time.Weekday, n int) (rds RecurringDateSchedule, err error) {
	if n == 0 || n > 4 || n < (-4) {
		err = errors.New("n must be between -4 and 4, excluding 0")
		return
	}
	rds = RecurringDateSchedule{wd: wd, n: n}
	return
}
