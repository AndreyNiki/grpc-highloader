package app

import (
	"github.com/AndreyNiki/grpc-highloader/internal/gui"
	"github.com/AndreyNiki/grpc-highloader/internal/loader"
	"github.com/AndreyNiki/grpc-highloader/internal/proto"
)

const (
	width  = 1540.0
	height = 750.0
)

// App main struct for all app.
type App struct{}

// NewApp create a new App.
func NewApp() *App {
	return &App{}
}

// Run start application.
func (app *App) Run() {
	requesterFactory := proto.NewRequesterFactory()
	loaderFactory := loader.NewLoaderFactory(requesterFactory)
	parser := proto.NewProtoParser()
	ui := gui.NewGUI(width, height, loaderFactory, parser)
	ui.Run()
}
