package element

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"strings"
)

/*

  File:    console.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
	Description: Handle user commands and program output.
*/

type Console struct {
	Content *fyne.Container
	Buttons []*widget.Button
	Focus   func()
	wait    chan string
	view    *fyne.Container
	row     int // current row
	rowMax  int
	entry   *widget.Entry
	prompt  string
	font    fyne.TextStyle
}

var Prompt = ">> "

func NewConsole(window fyne.Window, command string, rowMax int) *Console {
	console := Console{
		Content: nil,
		wait:    make(chan string),
		view:    nil,
		row:     0,
		rowMax:  rowMax,
		entry:   widget.NewEntry(),
		prompt:  Prompt,
		font:    fyne.TextStyle{Monospace: true},
	}
	console.entry.PlaceHolder = command
	p := widget.NewLabel(Prompt)
	p.TextStyle = fyne.TextStyle{Bold: true}
	bottom := container.NewBorder(nil, nil, p, nil, console.entry)
	console.view = container.NewVBox()
	console.Content = container.NewBorder(nil, bottom, nil, nil, console.view)
	console.Focus = func() { window.Canvas().Focus(console.entry) }
	return &console
}
func (c *Console) speakResponse(required string) {
	c.Speak("\n" + c.prompt + required)
}
func (c *Console) Ask(required string) (b string) {
	c.entry.SetText("")
	c.entry.OnSubmitted = func(response string) {
		if response == "" && required != "" {
			c.speakResponse(required)
		} else {
			c.speakResponse(response)
			c.wait <- response
		}
	}
	b = <-c.wait
	return
}
func (c *Console) AskYesNo(required string) (b bool) {
	c.entry.SetText("")
	c.entry.OnSubmitted = func(response string) {
		if response != "" {
			l := strings.ToLower(response)[:1]
			if l == "y" || l == "n" {
				c.speakResponse(response)
				c.wait <- strings.ToLower(response)[:1]
				return
			}
		}
		c.entry.SetText("")
		c.Speak(Prompt + required)
	}
	c.Focus()
	ok := <-c.wait
	if ok == "y" {
		b = true
	}
	return
}
func (c *Console) Speak(txt string) {
	if txt == "" {
		return
	}
	lines := strings.Split(txt, "\n")
	if len(lines) > c.rowMax {
		lines = lines[len(lines)-c.rowMax:]
	}
	avail := c.rowMax - c.row
	if avail < len(lines) {
		// remove from top
		c.view.Objects = c.view.Objects[len(lines)-avail:]
		c.row = c.rowMax - len(lines)
	}
	for _, line := range lines {
		t := canvas.NewText(line, theme.ForegroundColor())
		t.TextStyle = c.font
		c.view.Objects = append(c.view.Objects, t)
		c.row++
	}
	c.view.Refresh()
}
