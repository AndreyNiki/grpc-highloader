package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/AndreyNiki/grpc-highloader/internal/gui/components/highloader"
	"github.com/AndreyNiki/grpc-highloader/internal/gui/interfaces"
)

// GUI struct for UI app.
type GUI struct {
	width         float32
	height        float32
	loaderFactory interfaces.LoaderFactory
	parser        interfaces.Parser
}

// NewGUI create a new GUI.
func NewGUI(width, height float32, loaderFactory interfaces.LoaderFactory, parser interfaces.Parser) *GUI {
	return &GUI{
		width:         width,
		height:        height,
		loaderFactory: loaderFactory,
		parser:        parser,
	}
}

// Run ui for application.
func (g *GUI) Run() {
	a := app.New()
	w := a.NewWindow("GRPC HighLoader v1.0")
	w.Resize(fyne.NewSize(g.width, g.height))

	highLoader := highloader.New(w, g.loaderFactory, g.parser)
	hlComponent := highLoader.InitComponent()

	t := container.NewAppTabs(
		container.NewTabItem("HighLoader", hlComponent),
		container.NewTabItem("About", widget.NewLabel("About")))

	t.SetTabLocation(container.TabLocationLeading)
	w.SetContent(t)
	w.ShowAndRun()
}
