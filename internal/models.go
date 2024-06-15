package internal

import (
	"fmt"
	"math"
	"time"
)

type WorkTimeDuration time.Duration

type HistoryRecord struct {
	Date     time.Time        `json:"date"`
	Worked   WorkTimeDuration `json:"worked"`
	NeedWork WorkTimeDuration `json:"need_work"`
	Overwork WorkTimeDuration `json:"overwork"`
}

type Data struct {
	NeedWork WorkTimeDuration `json:"need_work"`
	Overwork WorkTimeDuration `json:"overwork"`
	History  []HistoryRecord  `json:"history"`
}

func NewData() *Data {
	return &Data{NeedWork: 0, Overwork: 0, History: []HistoryRecord{}}
}

func NewHistoryRecord(needWork, workedDuration WorkTimeDuration) *HistoryRecord {
	return &HistoryRecord{
		Date:     time.Now(),
		Worked:   workedDuration,
		NeedWork: needWork,
		Overwork: workedDuration - needWork,
	}
}

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
