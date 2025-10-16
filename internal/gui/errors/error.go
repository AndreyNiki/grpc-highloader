package errors

import (
	"image/color"

	"fyne.io/fyne/v2/canvas"
)

// GUIError error for show in GUI.
type GUIError struct {
	Err  error
	Text *canvas.Text
}

// NewGUIError create a new GUIError.
func NewGUIError(err error) *GUIError {
	errText := canvas.NewText(err.Error(), color.NRGBA{
		R: 188,
		G: 0,
		B: 0,
		A: 255,
	})

	return &GUIError{
		Err:  err,
		Text: errText,
	}
}
