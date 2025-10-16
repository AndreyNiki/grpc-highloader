package highloader

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/AndreyNiki/grpc-highloader/internal/gui/components/highloader/cards"
	"github.com/AndreyNiki/grpc-highloader/internal/gui/components/highloader/config"
	guierrs "github.com/AndreyNiki/grpc-highloader/internal/gui/errors"
	"github.com/AndreyNiki/grpc-highloader/internal/gui/interfaces"
)

const (
	protoExtension        = ".proto"
	buttonUploadProtoName = "Upload Proto"
	buttonOpenConfigName  = "Open Config"
	buttonSaveConfigName  = "Save Config"
	labelHostName         = "Host"
)

// HighLoader struct for init highloader component.
type HighLoader struct {
	window        fyne.Window
	loaderFactory interfaces.LoaderFactory
	parser        interfaces.Parser
}

// New create HighLoader.
func New(w fyne.Window, loaderFactory interfaces.LoaderFactory, parser interfaces.Parser) *HighLoader {
	return &HighLoader{
		window:        w,
		loaderFactory: loaderFactory,
		parser:        parser,
	}
}

// InitComponent init component.
func (h *HighLoader) InitComponent() fyne.CanvasObject {
	box := container.NewVBox()

	lineEntryHost := widget.NewEntry()
	var currentErr *guierrs.GUIError
	protoCardHolder := cards.NewProtoCardsHolder()
	buttonUploadProto := widget.NewButton(buttonUploadProtoName, func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if reader != nil {
				if currentErr != nil {
					box.Remove(currentErr.Text)
				}

				fp := reader.URI().Path()
				extension := filepath.Ext(fp)
				if extension != protoExtension {
					dialog.ShowError(fmt.Errorf("extension file is not %s", protoExtension), h.window)
					return
				}

				parsedProto, err := h.parser.ParseProto(fp)
				if err != nil {
					guiErr := guierrs.NewGUIError(err)
					box.Add(guiErr.Text)
					currentErr = guiErr
					return
				}
				c := &cards.ContainerCards{
					Parent:        box,
					Proto:         parsedProto,
					LoaderFactory: h.loaderFactory,
					Host:          lineEntryHost,
				}
				protoCardHolder.Add(c, nil)
			}
		}, h.window)
	})
	buttonOpenConfig := widget.NewButton(buttonOpenConfigName, func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if reader != nil {
				if currentErr != nil {
					box.Remove(currentErr.Text)
				}
				b, er := io.ReadAll(reader)
				if er != nil {
					guiErr := guierrs.NewGUIError(er)
					box.Add(guiErr.Text)
					currentErr = guiErr
					return
				}
				var preloadConfig config.PreloadConfig
				er = json.Unmarshal(b, &preloadConfig)
				if er != nil {
					guiErr := guierrs.NewGUIError(er)
					box.Add(guiErr.Text)
					currentErr = guiErr
					return
				}
				lineEntryHost.SetText(preloadConfig.Host)
				for _, p := range preloadConfig.Proto {
					parsedProto, err := h.parser.ParseProto(p.FilePath)
					if err != nil {
						guiErr := guierrs.NewGUIError(err)
						box.Add(guiErr.Text)
						currentErr = guiErr
						return
					}
					c := &cards.ContainerCards{
						Parent:        box,
						Proto:         parsedProto,
						LoaderFactory: h.loaderFactory,
						Host:          lineEntryHost,
					}
					protoCardHolder.Add(c, &p)
				}
			}
		}, h.window)
	})
	buttonSaveConfig := widget.NewButton(buttonSaveConfigName, func() {
		var proto []config.Proto
		for _, v := range protoCardHolder.Cards.Holder {
			var requests []config.Request
			for _, req := range v.RequestCardHolder.Cards.Holder {
				r := config.Request{
					LogPath:     req.Form.LogPath.Value.Text,
					MetricsPath: req.Form.MetricsPath.Value.Text,
					Message:     req.Form.ServicesMethods.MessageEntry.Text,
					StopAfter: config.Time{
						Duration: req.Form.StopAfter.Entry.Value.Text,
						Type:     req.Form.StopAfter.Select.Selected,
					},
					RequestDeadline: config.Time{
						Duration: req.Form.DeadlineReq.Entry.Value.Text,
						Type:     req.Form.DeadlineReq.Select.Selected,
					},
					Service: req.Form.ServicesMethods.Services.Selected,
					Method:  req.Form.ServicesMethods.Methods.Selected,
					RPS:     req.Form.RPS.Value.Text,
				}

				var metadata []config.MetaData
				for _, m := range req.Form.Metadata.KeyValues {
					md := config.MetaData{
						Key:   m.Key.Text,
						Value: m.Value.Text,
					}
					metadata = append(metadata, md)
				}
				r.Metadata = metadata
				requests = append(requests, r)
			}

			p := config.Proto{
				FilePath: v.FilePath,
				Requests: requests,
			}
			proto = append(proto, p)
		}

		cfg := config.PreloadConfig{
			Host:  lineEntryHost.Text,
			Proto: proto,
		}

		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			jsonData, err := json.MarshalIndent(cfg, "", "  ")
			if err != nil {
				dialog.ShowError(fmt.Errorf("error save config %w", err), h.window)
				return
			}
			_, err = writer.Write(jsonData)
			if err != nil {
				dialog.ShowError(fmt.Errorf("error save config %w", err), h.window)
				return
			}

			dialog.ShowInformation("Success", "Successfully saved config", h.window)
		}, h.window)
	})
	buttonSaveConfig.Importance = widget.HighImportance

	buttonUploadProto.Importance = widget.WarningImportance
	buttonOpenConfig.Importance = widget.WarningImportance
	buttonBox := container.NewGridWithColumns(2, buttonUploadProto, buttonOpenConfig)
	scroll := container.NewVScroll(
		container.NewVBox(widget.NewLabel(labelHostName), lineEntryHost, box, buttonBox, buttonSaveConfig))
	return scroll
}
