package datetimeinterval

import "slices"

type DateTimeIntervals []DateTimeInterval

func (d *DateTimeIntervals) Add(dti DateTimeInterval) {
	*d = append(*d, dti)
}

func (d *DateTimeIntervals) AddAll(dtis DateTimeIntervals) {
	*d = append(*d, dtis...)
}

func (d *DateTimeIntervals) Copy() (d2 DateTimeIntervals) {
	d2 = make(DateTimeIntervals, len(*d))
	copy(d2, *d)
	return
}

func (d *DateTimeIntervals) GetNextExtended() (dti DateTimeInterval, ok bool) {
	if len(*d) == 0 {
		return
	}
	ok = true
	d2 := d.Copy()
	slices.SortFunc(d2, dticmp)
	dti = d2[0]
	//iterate over the next intervals
	for _, dti2 := range d2[1:] {
		//dti2 starts before dti ends, so we can extend dti to include dti2
		if dti.End.Compare(dti2.Start) >= 0 {
			dti.End = dti2.End
		} else {
			// no overlap. Ignore all other ones
			break
		}
	}
	return
}
