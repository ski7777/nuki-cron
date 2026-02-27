package config

import "github.com/ski7777/nuki-cron/internal/datetimeinterval"

type ExtraEvent struct {
	Start ExtraDateTime `json:"start"`
	End   ExtraDateTime `json:"end"`
}

func (e *ExtraEvent) ToDateTimeInterval() datetimeinterval.DateTimeInterval {
	return datetimeinterval.DateTimeInterval{
		Start: e.Start.ToTime(),
		End:   e.End.ToTime(),
	}
}
