package cards

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/AndreyNiki/grpc-highloader/internal/entity"
	"github.com/AndreyNiki/grpc-highloader/internal/gui/components/highloader/config"
	"github.com/AndreyNiki/grpc-highloader/internal/gui/components/highloader/mapper"
	"github.com/AndreyNiki/grpc-highloader/internal/gui/interfaces"
	"github.com/AndreyNiki/grpc-highloader/internal/gui/utils"
	"github.com/AndreyNiki/grpc-highloader/internal/logger"
	"github.com/AndreyNiki/grpc-highloader/internal/metrics"
	"github.com/AndreyNiki/grpc-highloader/internal/utils/ptr"
)

const (
	buttonStartRequestName     = "Start"
	buttonStopRequestName      = "Stop"
	buttonRemoveRequestName    = "Remove"
	buttonAddKeyValueName      = "+"
	labelServicesName          = "Services"
	labelMethodsName           = "Methods"
	labelRPSName               = "~Req/s"
	labelDurationSecondsName   = "Stop After"
	labelRequestSettingName    = "Request Settings"
	labelMessageName           = "Message"
	labelMetadataName          = "Metadata"
	labelMetadataKeyName       = "Key"
	labelMetadataValueName     = "Value"
	labelAdditionalOptionsName = "Additional Options"
	labelRequestCardName       = "Request"
	labelRequestDeadlineName   = "Request Deadline"
)

const (
	rpsDefault = "100"
)

// RequestsCardsHolder struct for management request cards.
type RequestsCardsHolder struct {
	Cards *Cards[*RequestCard]
}

// NewRequestsCardsHolder create a new RequestsCardsHolder.
func NewRequestsCardsHolder() *RequestsCardsHolder {
	return &RequestsCardsHolder{
		Cards: NewCards[*RequestCard](),
	}
}

// Add new card to the parent element.
func (p *RequestsCardsHolder) Add(containerCards *ContainerCards, preloadRequest *config.Request) *RequestCard {
	c := newRequestCard(containerCards, preloadRequest)
	c.parent.Add(c.card)

	id := p.Cards.Add(c)
	c.buttonRemove.OnTapped = func() {
		p.Cards.Remove(id)
		c.parent.Remove(c.card)
	}
	return c
}

// RequestCard struct with info for request card.
type RequestCard struct {
	card          *widget.Card
	parent        *fyne.Container
	loaderFactory interfaces.LoaderFactory
	buttonStart   *widget.Button
	buttonStop    *widget.Button
	buttonRemove  *widget.Button
	Form          *FormRequest
}

// newRequestCard create a new RequestCard.
func newRequestCard(containerCards *ContainerCards, preloadRequest *config.Request) *RequestCard {
	r := &RequestCard{
		parent:        containerCards.Parent,
		loaderFactory: containerCards.LoaderFactory,
	}

	le := utils.NewEntry("Debug Log Path", nil, ptr.ToPtr("If no set then no saved logs"))
	me := utils.NewEntry("Metrics Path", nil, ptr.ToPtr("If no set then no saved metrics"))
	lb := container.NewVBox(le.Label, le.Value)
	mb := container.NewVBox(me.Label, me.Value)

	sm := r.makeServicesMethods(containerCards.Proto)
	bs := container.NewVBox(widget.NewLabel(labelServicesName), sm.Services)
	bm := container.NewVBox(widget.NewLabel(labelMethodsName), sm.Methods)
	bsm := container.NewVBox(bs, bm)

	rps := utils.NewEntry(labelRPSName, ptr.ToPtr(rpsDefault),
		ptr.ToPtr(fmt.Sprintf("default %q", rpsDefault)))
	geWorkers := container.NewGridWithColumns(2, rps.Value, rps.Label)
	sa := utils.NewEntryTime(labelDurationSecondsName, nil, nil, nil)
	gsa := container.NewGridWithColumns(3, sa.Entry.Label, sa.Entry.Value, sa.Select)
	dr := utils.NewEntryTime(labelRequestDeadlineName, nil, nil, nil)
	gdr := container.NewGridWithColumns(3, dr.Entry.Label, dr.Entry.Value, dr.Select)

	vBoxRS := container.NewVBox(widget.NewLabel(labelRequestSettingName), geWorkers, gsa, gdr)
	hBoxRS := container.NewHBox(utils.NewLine(), vBoxRS)

	rb := container.NewVBox(widget.NewLabel(labelMessageName), sm.MessageEntry)
	lmd := container.NewGridWithColumns(2, widget.NewLabel(labelMetadataKeyName),
		widget.NewLabel(labelMetadataValueName))

	smd := container.NewVBox()
	metadata := NewMetaData(smd)
	buttonAddMetadata := widget.NewButton(buttonAddKeyValueName, func() {
		metadata.AddKeyValue("", "")
	})
	bmd := container.NewVBox(widget.NewLabel(labelMetadataName), lmd, smd, buttonAddMetadata)

	ao := widget.NewAccordionItem(labelAdditionalOptionsName, widget.NewLabel("coming soon..."))

	buttonRemove := widget.NewButton(buttonRemoveRequestName, nil)
	buttonRemove.Importance = widget.DangerImportance

	timeTrackerCh := make(chan struct{}, 1)
	cancelSignal := make(chan struct{}, 1)
	mtrcs := metrics.InitMetrics()

	form := &FormRequest{
		LogPath:         le,
		MetricsPath:     me,
		RPS:             rps,
		StopAfter:       sa,
		DeadlineReq:     dr,
		ServicesMethods: sm,
		TimeTrackerCh:   timeTrackerCh,
		CancelCh:        cancelSignal,
		ButtonRemove:    buttonRemove,
		Host:            containerCards.Host,
		Metadata:        metadata,
		Metrics:         mtrcs,
		ParsedProto:     containerCards.Proto,
	}
	if preloadRequest != nil {
		r.setupFormFromPreloadRequest(form, preloadRequest)
	}

	layerHeader := container.NewGridWithColumns(2, lb, mb)
	layerTop := container.NewGridWithColumns(2, bsm, hBoxRS)
	layerMiddle := container.NewGridWithColumns(2, rb, bmd)
	layerAdditional := widget.NewAccordion(ao)
	layerController := r.makeControllerRequest(form)

	mainBox := container.NewVBox(
		layerHeader, utils.NewLine(),
		layerTop, utils.NewLine(),
		layerMiddle, utils.NewLine(),
		layerAdditional, utils.NewLine(),
		layerController)

	card := widget.NewCard("", labelRequestCardName, mainBox)
	buttonRemove.OnTapped = func() {
		close(timeTrackerCh)
		containerCards.Parent.Remove(card)
	}

	r.card = card
	r.buttonRemove = buttonRemove
	r.Form = form
	return r
}

// setupFormFromPreloadRequest preload request from specified config.
func (r *RequestCard) setupFormFromPreloadRequest(fr *FormRequest, request *config.Request) {
	if request.LogPath != "" {
		fr.LogPath.Value.SetText(request.LogPath)
	}
	if request.MetricsPath != "" {
		fr.MetricsPath.Value.SetText(request.MetricsPath)
	}
	if request.RPS != "" {
		fr.RPS.Value.SetText(request.RPS)
	}
	if request.StopAfter.Duration != "" && request.StopAfter.Type != "" {
		fr.StopAfter.FindAndSetOption(request.StopAfter.Duration, request.StopAfter.Type)
	}
	if request.RequestDeadline.Duration != "" && request.RequestDeadline.Type != "" {
		fr.DeadlineReq.FindAndSetOption(request.RequestDeadline.Duration, request.RequestDeadline.Type)
	}

	if len(request.Metadata) != 0 {
		for _, m := range request.Metadata {
			fr.Metadata.AddKeyValue(m.Key, m.Value)
		}
	}
	fr.ServicesMethods.preset(request.Service, request.Method, request.Message)
}

// makeControllerRequest make controller in GUI for control request.
func (r *RequestCard) makeControllerRequest(fr *FormRequest) *fyne.Container {
	buttonStop := widget.NewButton(buttonStopRequestName, nil)
	buttonStop.Disable()
	buttonStop.Importance = widget.MediumImportance

	buttonStart := widget.NewButton(buttonStartRequestName, nil)
	buttonStart.Importance = widget.HighImportance

	r.buttonStop = buttonStop
	r.buttonStart = buttonStart
	infoLabel := widget.NewLabel("")
	infoLabel.Wrapping = fyne.TextWrapBreak
	vsInfoLabel := container.NewVScroll(infoLabel)
	vsInfoLabel.Resize(fyne.NewSize(550, 30))
	vsInfoLabel.Hide()

	timeLabel := widget.NewLabel("")
	buttonStart.OnTapped = func() {
		vsInfoLabel.Hide()
		err := r.startLoadingRequests(fr)
		if err != nil {
			vsInfoLabel.Show()
			infoLabel.SetText(err.Error())
			return
		}
		buttonStop.Enable()
		buttonStart.Disable()
		fr.ButtonRemove.Disable()
		startTime := time.Now()
		utils.DurationLabel(fr.TimeTrackerCh, startTime, timeLabel)
	}
	buttonStop.OnTapped = func() {
		r.stopLoadingRequests(fr)
	}
	vfEntryErrFn := func() {
		buttonStart.Disable()
	}
	vfEntryPassFn := func() {
		buttonStart.Enable()
	}

	vf := utils.NewValidationForm(vfEntryErrFn, vfEntryPassFn)
	vf.AddValidationEntries(
		&utils.ValidationEntry{Entry: fr.RPS.Value, Validator: utils.NumberValidation()},
		&utils.ValidationEntry{Entry: fr.StopAfter.Entry.Value, Validator: utils.NumberValidation()},
		&utils.ValidationEntry{Entry: fr.DeadlineReq.Entry.Value, Validator: utils.NumberValidation()})
	vf.SetOrRefreshValidate()

	return container.NewHBox(
		buttonStart, buttonStop,
		fr.ButtonRemove, timeLabel,
		container.NewWithoutLayout(vsInfoLabel))
}

// stopLoadingRequests stop requests.
func (r *RequestCard) stopLoadingRequests(fr *FormRequest) {
	if fr.Logger != nil {
		fr.Logger.Close()
	}
	fr.CancelCh <- struct{}{}
	fr.TimeTrackerCh <- struct{}{}
	fyne.DoAndWait(func() {
		r.buttonStop.Disable()
		r.buttonStart.Enable()
		r.buttonRemove.Enable()
	})
}

// startLoadingRequests run requests with specified params from Form.
func (r *RequestCard) startLoadingRequests(fr *FormRequest) error {
	ctx := context.Background()
	if fr.LogPath.Value.Text != "" {
		log, err := logger.NewLogFile(fr.LogPath.Value.Text)
		if err != nil {
			return err
		}
		fr.Logger = log
		ctx = logger.ContextWithLogger(ctx, log.GetLogger())
	}
	rps, err := strconv.Atoi(fr.RPS.GetValue())
	if err != nil {
		return fmt.Errorf("could not parse rps: %w", err)
	}

	req := &entity.RequestParams{
		Service:  fr.ParsedProto.Package + "." + fr.ServicesMethods.Services.Selected,
		Method:   fr.ServicesMethods.Methods.Selected,
		Message:  fr.ServicesMethods.MessageEntry.Text,
		Metadata: fr.Metadata.MapString(),
		RPS:      rps,
		Host:     fr.Host.Text,
		Proto:    fr.ParsedProto,
	}
	method, ok := req.Proto.FindMethodByName(fr.ServicesMethods.Services.Selected, req.Method)
	if !ok {
		return fmt.Errorf("could not find method with name %s", req.Method)
	}
	req.MethodType = method.Type

	if fr.DeadlineReq.GetValue() != 0 {
		req.RequestDeadline = ptr.ToPtr(fr.DeadlineReq.GetValue())
	}

	loader, err := r.loaderFactory.NewLoader(req, fr.Metrics)
	if err != nil {
		return fmt.Errorf("could not create loader: %w", err)
	}

	ctx, cancel := context.WithCancel(ctx)
	r.stopRequestsManager(fr, cancel)

	go func() {
		if err := loader.Run(ctx); err != nil {
			r.stopLoadingRequests(fr)
			return
		}
	}()

	return nil
}

func (r *RequestCard) stopRequestsManager(fr *FormRequest, cancel context.CancelFunc) {
	if fr.StopAfter.GetValue() != 0 {
		go func() {
			for {
				select {
				case <-time.After(fr.StopAfter.GetValue()):
					r.stopLoadingRequests(fr)
					return
				}
			}
		}()
	}
	go func() {
		<-fr.CancelCh
		cancel()
	}()
}

// makeServicesMethods make services and methods for GUI form.
func (r *RequestCard) makeServicesMethods(parsedProto *entity.ParsedProto) *ServicesMethods {
	mapperUI := mapper.NewMapper(parsedProto)
	protoMapUI := mapperUI.MakeProtoMapUI(parsedProto)
	messageEntry := widget.NewMultiLineEntry()

	optionsServices := protoMapUI.GetServicesNames()
	methods := widget.NewSelect(nil, func(value string) {
		msg, ok := protoMapUI.GetMessageByMethodName(value)
		if !ok {
			messageEntry.SetText("Example message not found")
			return
		}
		exampleMessage := mapperUI.MakeExampleMessage(msg)
		j, err := json.MarshalIndent(exampleMessage, "", "  ")
		if err != nil {
			messageEntry.SetText("Error marshalling exampleMessage")
			return
		}
		messageEntry.SetText(string(j))
	})

	services := widget.NewSelect(optionsServices, func(value string) {
		methods.Options = protoMapUI.GetMethodsNamesByService(value)
		methods.SetSelectedIndex(0)
	})
	services.SetSelectedIndex(0)

	return &ServicesMethods{
		Services:     services,
		Methods:      methods,
		MessageEntry: messageEntry,
	}
}
