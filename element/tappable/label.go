package tappable

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"time"
)

/*

  File:    label.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:A tappable widget.Label.
  Text, Bindidg, and Style may be set / reset in the base Label.
*/

type Label struct {
	widget.Label
	ID       int
	OnTapped OnTap
	lastTime time.Time
}

//goland:noinspection GoUnusedExportedFunction
func NewLabel(text string, onTap OnTap) *Label {
	label := &Label{OnTapped: onTap}
	label.ExtendBaseWidget(label)
	label.Text = text
	return label
}
func (t *Label) Tapped(pe *fyne.PointEvent) {
	duration := time.Now().Sub(t.lastTime)
	if duration <= TapperDoubleClickTime {
		t.OnTapped(Double, t.ID, pe)
		t.lastTime = time.Time{}
	} else {
		t.OnTapped(Primary, t.ID, pe)
		t.lastTime = time.Now()
	}
}
func (t *Label) TappedSecondary(pe *fyne.PointEvent) {
	t.OnTapped(Secondary, t.ID, pe)
}
