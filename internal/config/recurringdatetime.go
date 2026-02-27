package config

import (
	"time"

	"github.com/ski7777/nuki-cron/internal/recurringdate"
)

type RecurringDateTime struct {
	Time
	Dow int `json:"dow"`
	N   int `json:"n"`
}

func (rdt *RecurringDateTime) ToRecurringDateSchedule() (recurringdate.RecurringDateSchedule, error) {
	return recurringdate.NewReccuringDateSchedule(time.Weekday(rdt.Dow), rdt.N)
}
