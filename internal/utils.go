package internal

import "time"

// Checks if two time.Time instances occur on the same calendar date (year, month, day).
// Returns true if both time.Time values have the same year, month, and day; false otherwise.
func IsSameDate(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
