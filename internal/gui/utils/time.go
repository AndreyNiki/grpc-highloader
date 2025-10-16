package utils

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// Time wrapper string for time.
type Time string

const (
	TimeNameMilliseconds Time = "Milliseconds"
	TimeNameSeconds      Time = "Seconds"
	TimeNameMinutes      Time = "Minutes"
	TimeNameHours        Time = "Hours"
)

// String convert Time to string.
func (t Time) String() string {
	return string(t)
}

// Duration return time.Duration by value.
func (t Time) Duration(value string) time.Duration {
	v, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	switch t {
	case TimeNameMilliseconds:
		return time.Millisecond * time.Duration(v)
	case TimeNameSeconds:
		return time.Second * time.Duration(v)
	case TimeNameMinutes:
		return time.Minute * time.Duration(v)
	case TimeNameHours:
		return time.Hour * time.Duration(v)
	default:
		return 0
	}
}

// DurationLabel showing time duration in real-time.
func DurationLabel(ch chan struct{}, startTime time.Time, label *widget.Label) {
	go func() {
		for {
			select {
			case <-ch:
				return
			default:
				timeDuration := time.Since(startTime)
				fyne.Do(func() {
					label.SetText(fmt.Sprintf("%.2fs", timeDuration.Seconds()))
				})
				time.Sleep(1 * time.Millisecond)
			}
		}
	}()
}
