package cards

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/AndreyNiki/grpc-highloader/internal/gui/components/highloader/config"
	"github.com/AndreyNiki/grpc-highloader/internal/gui/utils"
)

const (
	buttonRemoveRequestsCard = "X"
	buttonAddRequestCard     = "Add Request"
	titleRequestsCard        = "Requests"
)

// ProtoCardsHolder struct for management proto cards.
type ProtoCardsHolder struct {
	Cards *Cards[*ProtoCard]
}

// NewProtoCardsHolder create a new ProtoCardsHolder.
func NewProtoCardsHolder() *ProtoCardsHolder {
	return &ProtoCardsHolder{
		Cards: NewCards[*ProtoCard](),
	}
}

// Add new card to the parent element.
func (p *ProtoCardsHolder) Add(containerCards *ContainerCards, proto *config.Proto) *ProtoCard {
	c := newProtoCard(containerCards, proto)
	c.parent.Add(c.card)

	id := p.Cards.Add(c)
	c.buttonRemove.OnTapped = func() {
		p.Cards.Remove(id)
		c.parent.Remove(c.card)
	}
	return c
}

// ProtoCard struct with info for proto card.
type ProtoCard struct {
	FilePath          string
	RequestCardHolder *RequestsCardsHolder
	card              *widget.Card
	parent            *fyne.Container
	buttonRemove      *widget.Button
}

// newProtoCard create a new ProtoCard.
func newProtoCard(containerCards *ContainerCards, proto *config.Proto) *ProtoCard {
	buttonRemove := widget.NewButton(buttonRemoveRequestsCard, nil)
	buttonAddReq := widget.NewButton(buttonAddRequestCard, nil)
	buttonRemove.Importance = widget.DangerImportance

	newButtonRemove := container.NewGridWithColumns(5, utils.NewObjectWithSpacers(4, buttonRemove)...)
	requestCardBox := container.NewVBox()
	cardBox := container.NewVBox(newButtonRemove, requestCardBox, buttonAddReq)
	card := widget.NewCard(titleRequestsCard, containerCards.Proto.FilePath, cardBox)

	requestCardHolder := NewRequestsCardsHolder()
	newContainer := &ContainerCards{
		Parent:        requestCardBox,
		Proto:         containerCards.Proto,
		Host:          containerCards.Host,
		LoaderFactory: containerCards.LoaderFactory,
	}
	buttonAddReq.OnTapped = func() {
		requestCardHolder.Add(newContainer, nil)
	}

	if proto != nil {
		for _, r := range proto.Requests {
			requestCardHolder.Add(newContainer, &r)
		}
	}
	return &ProtoCard{
		FilePath:          containerCards.Proto.FilePath,
		RequestCardHolder: requestCardHolder,
		card:              card,
		parent:            containerCards.Parent,
		buttonRemove:      buttonRemove,
	}
}
