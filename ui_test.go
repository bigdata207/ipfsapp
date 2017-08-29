package main

import (
	"testing"

	"github.com/andlabs/ui"
)

func initUI() {
	name := ui.NewEntry()
	button := ui.NewButton("Greet")
	greeting := ui.NewLabel("")
	box := ui.NewVerticalBox()
	box.Append(ui.NewLabel("Enter your name:"), false)
	box.Append(name, false)
	box.Append(button, false)
	box.Append(greeting, false)
	window := ui.NewWindow("Hello", 200, 100, false)
	window.SetChild(box)
	button.OnClicked(func(*ui.Button) {
		greeting.SetText("Hello, " + name.Text() + "!")
	})
	window.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	window.Show()
}
func TestUI(t *testing.T) {
	err := ui.Main(initUI)
	if err != nil {
		panic(err)
	}
}
