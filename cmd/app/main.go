package main

import (
	"github.com/AndreyNiki/grpc-highloader/internal/app"
)

func main() {
	a := app.NewApp()
	a.Run()
}
