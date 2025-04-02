package util

import (
	"math"
	"time"
)

type DateInterval struct {
	interval *time.Duration
}

// GetEndDate implements DateIntervalInterface.
func (d *DateInterval) GetEndDate() time.Time {
	return time.Now().Add(*d.interval)
}

// GetEndTimeMs implements DateIntervalInterface.
func (d *DateInterval) GetEndTimeMs() int64 {
	return time.Now().Add(*d.interval).UnixMilli()
}

// GetEndTimeSeconds implements DateIntervalInterface.
func (d *DateInterval) GetEndTimeSeconds() int64 {
	return int64(math.Ceil(float64(d.GetEndTimeMs()) / 1000))
}

// GetSeconds implements DateIntervalInterface.
func (d *DateInterval) GetSeconds() int64 {
	return int64(math.Ceil(float64(*d.interval) / float64(time.Second)))
}

func NewDateInterface(interval *time.Duration) *DateInterval {
	return &DateInterval{interval: interval}
}
