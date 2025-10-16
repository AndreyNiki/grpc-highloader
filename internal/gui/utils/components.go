package utils

import (
	"image/color"
	"slices"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// Entry wrapper for *widget.Entry.
type Entry struct {
	Value        *widget.Entry
	Label        *widget.Label
	DefaultValue *string
}

// NewEntry create a new Entry.
func NewEntry(labelName string, defaultValue, placeholder *string) *Entry {
	entry := widget.NewEntry()
	if placeholder != nil {
		entry.SetPlaceHolder(*placeholder)
	}

	return &Entry{
		Value:        entry,
		Label:        widget.NewLabel(labelName),
		DefaultValue: defaultValue,
	}
}

// GetValue return value for Entry.
func (e *Entry) GetValue() string {
	if e.Value.Text != "" {
		return e.Value.Text
	}

	if e.DefaultValue != nil {
		return *e.DefaultValue
	}

	return ""
}

// EntryTime wrapper for *widget.Entry.
type EntryTime struct {
	Entry       *Entry
	Select      *widget.Select
	DefaultTime *DefaultTime
}

// DefaultTime struct for default time.
type DefaultTime struct {
	Value string
	Time  Time
}

// NewEntryTime create a new EntryTime.
func NewEntryTime(labelName string, defaultTime *DefaultTime, placeholder *string, fn func(value string)) *EntryTime {
	var defaultValue *string
	if defaultTime != nil {
		defaultValue = &defaultTime.Value
	}

	e := NewEntry(labelName, defaultValue, placeholder)
	s := NewSelectTimes(fn)

	s.SetSelectedIndex(0)
	return &EntryTime{
		Entry:       e,
		Select:      s,
		DefaultTime: defaultTime,
	}
}

// GetValue return value for EntryTime.
func (e *EntryTime) GetValue() time.Duration {
	if e.Entry.Value.Text != "" {
		return Time(e.Select.Selected).Duration(e.Entry.Value.Text)
	}
	if e.DefaultTime != nil {
		return e.DefaultTime.Time.Duration(e.DefaultTime.Value)
	}

	return 0
}

// FindAndSetOption find and set option for EntryTime.
func (e *EntryTime) FindAndSetOption(value, typ string) {
	if slices.Contains(e.Select.Options, typ) {
		e.Select.Selected = typ
	} else {
		e.Select.Selected = "not found"
	}

	e.Entry.Value.SetText(value)
}

// NewLine make new *canvas.Line.
func NewLine() *canvas.Line {
	return canvas.NewLine(color.NRGBA{
		R: 18,
		G: 18,
		B: 18,
		A: 255,
	})
}

// NewObjectWithSpacers make objects with spacers.
func NewObjectWithSpacers(len int, obj fyne.CanvasObject) []fyne.CanvasObject {
	objects := make([]fyne.CanvasObject, 0, len+1)
	for i := 0; i < len; i++ {
		objects = append(objects, layout.NewSpacer())
	}
	objects = append(objects, obj)
	return objects
}

// NewSelectTimes new times with preset select times.
func NewSelectTimes(fn func(value string)) *widget.Select {
	times := []string{
		TimeNameMilliseconds.String(),
		TimeNameSeconds.String(),
		TimeNameMinutes.String(),
		TimeNameHours.String()}
	selects := widget.NewSelect(times, fn)

	return selects
}
