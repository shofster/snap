package tappable

import (
	"fyne.io/fyne/v2"
	"time"
)

/*

  File:    tapper.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: Types for Tappable widgets.
*/

type Tapper int

const (
	Primary Tapper = iota
	Secondary
	Double
)

type OnTap func(Tapper, int, *fyne.PointEvent)

var TapperDoubleClickTime = time.Millisecond * 500

// Declare conformity with WidgetRenderer interface
var _ fyne.WidgetRenderer = (*emptyRenderer)(nil)

type emptyRenderer struct {
	objects []fyne.CanvasObject
	layout  fyne.Layout
	size    fyne.Size
}

func (r *emptyRenderer) Destroy() {
}
func (r *emptyRenderer) Layout(_ fyne.Size) {
	r.layout.Layout(r.objects, r.size)
}
func (r *emptyRenderer) MinSize() fyne.Size {
	return r.size
}
func (r *emptyRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}
func (r *emptyRenderer) Refresh() {
	r.objects[0].Refresh()
}
func (r *emptyRenderer) applyTheme() {
	r.objects[0].Refresh()
}
