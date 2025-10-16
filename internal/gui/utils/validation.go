package utils

import (
	"errors"
	"regexp"

	"fyne.io/fyne/v2/widget"
)

// ValidationEntry wrapper *widget.Entry with validations.
type ValidationEntry struct {
	Entry     *widget.Entry
	Validator func(s string) error
	invalid   bool
}

// ValidationForm stored entries for form validation.
type ValidationForm struct {
	validateErrFn  func()
	validatePassFn func()
	entries        []*ValidationEntry
}

// NewValidationForm create a new form with validations.
func NewValidationForm(validateErrFn func(), validatePassFn func()) *ValidationForm {
	return &ValidationForm{
		validateErrFn:  validateErrFn,
		validatePassFn: validatePassFn,
	}
}

// AddValidationEntries add ValidationEntry to ValidationForm.
func (v *ValidationForm) AddValidationEntries(entries ...*ValidationEntry) {
	v.entries = append(v.entries, entries...)
}

// SetOrRefreshValidate set or refresh validation for entry.
func (v *ValidationForm) SetOrRefreshValidate() {
	for _, e := range v.entries {
		e.Entry.SetOnValidationChanged(func(err error) {
			if err != nil {
				e.Entry.SetValidationError(err)
				e.Entry.Refresh()
				e.invalid = true
				v.validateErrFn()
			} else {
				e.invalid = false
				e.Entry.SetValidationError(nil)
				e.Entry.Refresh()
				if v.isNotExistsInvalidEntry() {
					v.validatePassFn()
				}
			}
		})
		e.Entry.Validator = e.Validator
	}
}

// isNotExistsInvalidEntry checks not exists invalid entries.
func (v *ValidationForm) isNotExistsInvalidEntry() bool {
	for _, e := range v.entries {
		if e.invalid {
			return false
		}
	}

	return true
}

// NumberValidation validation for number entry.
func NumberValidation() func(s string) error {
	return func(s string) error {
		if s == "" {
			return nil
		}
		numberRegex := regexp.MustCompile(`^\d+$`)
		if !numberRegex.MatchString(s) {
			return errors.New("field is must be numeric")
		}

		return nil
	}
}
