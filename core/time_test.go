package core

import "time"

func Date(year, months, day int) time.Time {
	return time.Date(year, time.Month(months), day, 0, 0, 0, 0, time.UTC)
}
