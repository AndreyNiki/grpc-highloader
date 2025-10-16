package cards

import (
	"encoding/json"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	"github.com/AndreyNiki/grpc-highloader/internal/entity"
	"github.com/AndreyNiki/grpc-highloader/internal/gui/interfaces"
	"github.com/AndreyNiki/grpc-highloader/internal/gui/utils"
	"github.com/AndreyNiki/grpc-highloader/internal/logger"
	"github.com/AndreyNiki/grpc-highloader/internal/metrics"
)

// Cards type for add/remove cards.
type Cards[T comparable] struct {
	NextID int
	Holder map[int]T
}

// NewCards create a new Cards.
func NewCards[T comparable]() *Cards[T] {
	return &Cards[T]{
		NextID: 0,
		Holder: map[int]T{},
	}
}

// Add new card.
func (c *Cards[T]) Add(t T) int {
	c.Holder[c.NextID] = t
	id := c.NextID
	c.NextID++
	return id
}

// Remove card by id.
func (c *Cards[T]) Remove(id int) {
	delete(c.Holder, id)
}

// ContainerCards struct for transport between cards.
type ContainerCards struct {
	Parent        *fyne.Container
	Host          *widget.Entry
	Proto         *entity.ParsedProto
	LoaderFactory interfaces.LoaderFactory
}

// FormRequest form with info from GUI.
type FormRequest struct {
	Logger          *logger.Logger
	LogPath         *utils.Entry
	MetricsPath     *utils.Entry
	RPS             *utils.Entry
	StopAfter       *utils.EntryTime
	DeadlineReq     *utils.EntryTime
	ServicesMethods *ServicesMethods
	Metadata        *Metadata
	TimeTrackerCh   chan struct{}
	CancelCh        chan struct{}
	ButtonRemove    *widget.Button
	Host            *widget.Entry
	Metrics         *metrics.Metrics
	ParsedProto     *entity.ParsedProto
}

// ServicesMethods struct with services and method. Also stored message for method.
type ServicesMethods struct {
	Services     *widget.Select
	Methods      *widget.Select
	MessageEntry *widget.Entry
}

// preset values in form GUI.
func (s *ServicesMethods) preset(service, method, message string) {
	if service != "" {
		if slices.Contains(s.Services.Options, service) {
			s.Services.Selected = service
		} else {
			s.Services.Selected = "Not Found"
		}
	}
	if method != "" {
		if slices.Contains(s.Methods.Options, method) {
			s.Methods.Selected = method
		} else {
			s.Methods.Selected = "Not Found"
		}
	}

	if message != "" {
		var m map[string]any
		err := json.Unmarshal([]byte(message), &m)
		if err != nil {
			s.MessageEntry.SetText("Error unmarshalling message")
			return
		}

		j, err := json.MarshalIndent(m, "", "  ")
		if err != nil {
			s.MessageEntry.SetText("Error marshalling message")
			return
		}
		s.MessageEntry.SetText(string(j))
	}
}
