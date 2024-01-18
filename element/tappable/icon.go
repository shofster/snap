package tappable

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"time"
)

/*

  File:    icon.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: A tappable widget.Icon.
*/

type Icon struct {
	widget.Icon
	ID       int
	OnTapped OnTap
	lastTime time.Time
}

//goland:noinspection GoUnusedExportedFunction
func NewIcon(res fyne.Resource, onTap OnTap) *Icon {
	icon := &Icon{OnTapped: onTap}
	icon.ExtendBaseWidget(icon)
	icon.SetResource(res)
	return icon
}
func (t *Icon) Tapped(pe *fyne.PointEvent) {
	duration := time.Now().Sub(t.lastTime)
	if duration <= TapperDoubleClickTime {
		t.OnTapped(Double, t.ID, pe)
		t.lastTime = time.Time{}
	} else {
		t.OnTapped(Primary, t.ID, pe)
		t.lastTime = time.Now()
	}
}
func (t *Icon) TappedSecondary(pe *fyne.PointEvent) {
	t.OnTapped(Secondary, t.ID, pe)
}
