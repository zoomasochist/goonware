package main

import (
	"time"

	g "github.com/AllenDang/giu"
)

func loop() {
	g.SingleWindow().Layout(
		g.Label("Test"),
	)
}

func MakeWindow() {
	wnd := g.NewMasterWindow("Window", 400, 400, 0)
	wnd.Run(loop)
}

func main() {
	go MakeWindow()
	time.Sleep(500 * time.Millisecond)
	MakeWindow()
}
