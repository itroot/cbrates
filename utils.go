package cbrates

import "time"

func formatDate(year int, month time.Month, day int) string {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Format(DateFormat)
}
