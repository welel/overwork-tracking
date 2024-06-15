package internal

import (
	"fmt"
	"math"
	"time"
)

// WorkTimeDuration represents a duration of time for work-related activities.
type WorkTimeDuration time.Duration

// HistoryRecord represents a record of work history for a specific date.
type HistoryRecord struct {
	Date     time.Time        `json:"date"`      // Date of the work record.
	Worked   WorkTimeDuration `json:"worked"`    // Time duration worked on that date.
	NeedWork WorkTimeDuration `json:"need_work"` // Required work duration for that date.
	Overwork WorkTimeDuration `json:"overwork"`  // Amount of overwork (worked - required).
}

// Data represents the overall data structure containing work tracking information.
type Data struct {
	NeedWork WorkTimeDuration `json:"need_work"` // Current required work duration for today.
	Overwork WorkTimeDuration `json:"overwork"`  // Total overwork accumulated (positive is overwork, negative is opposite).
	History  []HistoryRecord  `json:"history"`   // List of historical work records.
}

// Creates a new Data structure with default values.
func NewData() *Data {
	return &Data{
		NeedWork: 0,
		Overwork: 0,
		History:  []HistoryRecord{},
	}
}

// Creates a new HistoryRecord with the given required work duration (needWork) and worked duration.
func NewHistoryRecord(needWork, workedDuration WorkTimeDuration) *HistoryRecord {
	return &HistoryRecord{
		Date:     time.Now(),
		Worked:   workedDuration,
		NeedWork: needWork,
		Overwork: workedDuration - needWork,
	}
}

// Returns a string representation of WorkTimeDuration in "-HH:MM" format.
func (d WorkTimeDuration) String() string {
	totalMinutes := int64(time.Duration(d).Minutes())
	sign := ""
	if totalMinutes < 0 {
		sign = "-"
	}
	h := int(math.Abs(float64(int(totalMinutes / 60))))
	m := int(math.Abs(float64(totalMinutes % 60)))
	return fmt.Sprintf("%s%02d:%02d", sign, h, m)
}
