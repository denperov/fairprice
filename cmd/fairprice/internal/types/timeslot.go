package types

import "time"

type Timeslot int64

func (t Timeslot) ToTime() time.Time {
	return time.Unix(int64(t), 0).UTC()
}
