package config

import (
	"time"

	"github.com/ski7777/nuki-cron/internal/recurringinterval"
)

type Schedule struct {
	Start RecurringDateTime `json:"start"`
	End   RecurringDateTime `json:"end"`
}

func (s *Schedule) ToRecurringInterval() (ri recurringinterval.RecurringInterval, err error) {
	start_rds, err := s.Start.ToRecurringDateSchedule()
	if err != nil {
		return
	}
	end_rds, err := s.End.ToRecurringDateSchedule()
	if err != nil {
		return
	}
	ri, err = recurringinterval.NewRecurringInterval(
		start_rds,
		end_rds,
		time.Duration(s.Start.H)*time.Hour+time.Duration(s.Start.M)*time.Minute,
		time.Duration(s.End.H)*time.Hour+time.Duration(s.End.M)*time.Minute,
	)
	return
}
