package cards

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	buttonRemoveKeyValueMetadata = "-"
)

// KeyValue key and value for Metadata.
type KeyValue struct {
	Key   *widget.Entry
	Value *widget.Entry
}

// Metadata for request.
type Metadata struct {
	KeyValues []KeyValue
	parent    *fyne.Container
}

// NewMetaData create a new Metadata.
func NewMetaData(parent *fyne.Container) *Metadata {
	return &Metadata{
		KeyValues: []KeyValue{},
		parent:    parent,
	}
}

// AddKeyValue add new key-value elements to the parent element.
func (m *Metadata) AddKeyValue(key, value string) {
	keyEntry := widget.NewEntry()
	valueEntry := widget.NewEntry()
	if key != "" {
		keyEntry.SetText(key)
	}
	if value != "" {
		valueEntry.SetText(value)
	}
	m.KeyValues = append(m.KeyValues, KeyValue{
		Key:   keyEntry,
		Value: valueEntry,
	})

	buttonRemove := widget.NewButton(buttonRemoveKeyValueMetadata, nil)
	buttonRemove.Importance = widget.DangerImportance
	gridKeyValues := container.NewGridWithColumns(3, keyEntry, valueEntry, buttonRemove)

	m.parent.Add(gridKeyValues)
	buttonRemove.OnTapped = func() {
		m.parent.Remove(gridKeyValues)
	}
}

// MapString return metadata in map[string]string.
func (m *Metadata) MapString() map[string]string {
	result := make(map[string]string)
	for _, kv := range m.KeyValues {
		result[kv.Key.Text] = kv.Value.Text
	}

	return result
}
