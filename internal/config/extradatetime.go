package config

import "time"

type ExtraDateTime struct {
	Time
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

func (t ExtraDateTime) ToTime() time.Time {
	return time.Date(t.Year, time.Month(t.Month), t.Day, t.H, t.M, 0, 0, time.Local)
}
