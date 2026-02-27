package datetimeinterval

import (
	"cmp"
	"time"
)

type DateTimeInterval struct {
	Start, End time.Time
}

func dticmp(dti1, dti2 DateTimeInterval) int {
	return cmp.Or(
		dti1.Start.Compare(dti2.Start),
		dti1.End.Compare(dti2.End),
	)
}
