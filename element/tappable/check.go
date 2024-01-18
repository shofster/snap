package tappable

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"time"
)

/*

  File:    check.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: A tappable widget.Check.
     TappedSecondary is a right click. doesn't change widget state.
     Double is a second left click within a specified time.
       Doesn't change widget state.
*/

type Check struct {
	widget.Check
	ID       int
	OnTapped OnTap
	lastTime time.Time
}

//goland:noinspection GoUnusedExportedFunction
func NewCheck(text string, onTap OnTap) *Check {
	check := &Check{OnTapped: onTap}
	check.ExtendBaseWidget(check)
	check.Text = text
	return check
}
func (t *Check) Tapped(pe *fyne.PointEvent) {
	if t.Checked {
		duration := time.Now().Sub(t.lastTime)
		if duration <= TapperDoubleClickTime {
			t.OnTapped(Double, t.ID, pe)
			t.lastTime = time.Time{}
			return
		} else {
			t.SetChecked(false)
		}
	} else {
		t.OnTapped(Primary, t.ID, pe)
		t.lastTime = time.Now()
		t.SetChecked(true)
	}
}
func (t *Check) TappedSecondary(pe *fyne.PointEvent) {
	t.OnTapped(Secondary, t.ID, pe)
}
