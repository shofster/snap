package tappable

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"time"
)

/*

  File:    text.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: A tappable canvas.Text.
*/

type Text struct {
	widget.BaseWidget
	OnTapped OnTap
	ID       int
	text     *canvas.Text
	lastTime time.Time
}

//goland:noinspection GoUnusedExportedFunction
func NewText(text *canvas.Text, onTap OnTap) *Text {
	t := &Text{
		text:     text,
		OnTapped: onTap,
	}
	t.ExtendBaseWidget(t)
	return t
}
func (t *Text) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)
	var objects []fyne.CanvasObject
	objects = append(objects, t.text)
	size := fyne.Size{
		Width:  t.text.TextSize,
		Height: t.text.TextSize,
	}
	r := &emptyRenderer{objects, layout.NewStackLayout(), size}
	r.applyTheme()
	return r
}

func (t *Text) Tapped(pe *fyne.PointEvent) {
	duration := time.Now().Sub(t.lastTime)
	if duration <= TapperDoubleClickTime {
		t.OnTapped(Double, t.ID, pe)
		t.lastTime = time.Time{}
	} else {
		t.OnTapped(Primary, t.ID, pe)
		t.lastTime = time.Now()
	}
}
func (t *Text) TappedSecondary(pe *fyne.PointEvent) {
	t.OnTapped(Secondary, t.ID, pe)
}
