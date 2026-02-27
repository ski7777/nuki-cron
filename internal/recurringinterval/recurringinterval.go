package recurringinterval

import (
	"errors"
	"time"

	"github.com/ski7777/nuki-cron/internal/datetimeinterval"
	"github.com/ski7777/nuki-cron/internal/recurringdate"
)

type RecurringInterval struct {
	dateStart, dateEnd     recurringdate.RecurringDateSchedule
	offsetStart, offsetEnd time.Duration
}

func (ri RecurringInterval) getTimes(t time.Time) (start, end time.Time, err error) {
	start, err = ri.dateStart.GetByTime(t)
	if err != nil {
		return
	}
	end, err = ri.dateEnd.GetByTime(t)
	if err != nil {
		return
	}
	start = start.Add(ri.offsetStart)
	end = end.Add(ri.offsetEnd)
	if end.Before(start) || start.Equal(end) {
		err = errors.New("end time must be after start time")
	}
	return
}

func (ri RecurringInterval) GetNextTimes() (dti datetimeinterval.DateTimeInterval, err error) {
	// if we are before or during the event in current month, return the current month
	// if we are after the event in current month, return the next month
	rel := time.Now()
	status, err := ri.GetStatus()
	if err != nil {
		return
	}
	if status == 1 {
		rel = time.Date(rel.Year(), rel.Month()+1, 1, 0, 0, 0, 0, rel.Location())
	}
	dti.Start, dti.End, err = ri.getTimes(rel)
	return
}

func (ri RecurringInterval) GetStatus() (status int, err error) {
	// status: -1 before, 0 during, 1 after
	now := time.Now()
	start, end, err := ri.getTimes(now)
	if err != nil {
		return
	}
	status = 0
	if now.Before(start) {
		status = -1
		return
	}
	if now.After(end) {
		status = 1
		return
	}
	return
}

func NewRecurringInterval(dateStart, dateEnd recurringdate.RecurringDateSchedule, offsetStart, offsetEnd time.Duration) (ri RecurringInterval, err error) {
	ri = RecurringInterval{dateStart: dateStart, dateEnd: dateEnd, offsetStart: offsetStart, offsetEnd: offsetEnd}
	_, _, err = ri.getTimes(time.Now())
	return
}
