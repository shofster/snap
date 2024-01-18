package fileutil

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"io/fs"
	"snap/element/tappable"
	"strings"
	"time"
)

/*

  File:    fileList.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*

 */

var defaultDateTimeFormat = "01/02/06 15:04"
var dtFormat = ""

func NewFileList(path string, sel FileSelectFilter, action FileSelectAction,
	pretty func(name string, info fs.FileInfo, err error) string) (*widget.List, *DirectoryEntry, error) {
	// get specific DirectoryEntry
	dir, err := NewDirectoryView(path, sel)
	if err != nil {
		return nil, nil, err
	}
	dtFormat = sel.DtFormat
	if pretty == nil { // use standard display
		pretty = ls_al
	}
	var fileList = &widget.List{}
	var clearSelections = func() {
		//fileList.UnselectAll()  // forces updates
		for ix := range dir.files {
			fp := dir.File(ix)
			if fp.selected {
				fp.SetSelected(false)
			}
		}
		fileList.Refresh()
	}
	fileList = widget.NewList(
		// length
		func() int {
			return dir.Count()
		},
		// create
		func() fyne.CanvasObject {
			// build the unique tappable Label
			tl := tappable.NewLabel("", nil)
			co := container.NewHBox(widget.NewIcon(theme.FileIcon()), tl)
			tl.OnTapped = func(t tappable.Tapper, id int, pe *fyne.PointEvent) {
				// fmt.Println("create a new label", t, " tap on", id)
				file := dir.File(id) // current file POINTER
				switch t {
				case tappable.Primary:
					if sel.FileType == File && file.IsDir() { // must have File
						break
					}
					if sel.Multiple { // allow many selections
						file.SetSelected(!file.selected)
						fileList.UpdateItem(id, co)
					} else { // only a single File or Dir, toggle select
						if !file.selected {
							clearSelections()
							file.SetSelected(true)
						} else {
							file.SetSelected(false)
						}
						fileList.UpdateItem(id, co)
					}
					file.index = id
					action.OnClick(*file)
				case tappable.Double:
					if file.IsDir() {
						clearSelections()
					}
					if action.OnDoubleClick != nil {
						action.OnDoubleClick(*file)
					}
				case tappable.Secondary:
					if action.OnSecondaryClick != nil {
						action.OnSecondaryClick(*file, pe)
					}
				}
			}
			return co
		},
		// update
		func(id widget.ListItemID, item fyne.CanvasObject) {
			file := dir.File(id)
			info, err := file.Info()
			item.(*fyne.Container).Objects[1].(*tappable.Label).Text = pretty(file.DisplayName(), info, err)
			item.(*fyne.Container).Objects[1].(*tappable.Label).ID = id
			if file.IsSelected() {
				item.(*fyne.Container).Objects[0].(*widget.Icon).SetResource(theme.ConfirmIcon())
			} else {
				if file.IsDir() {
					item.(*fyne.Container).Objects[0].(*widget.Icon).SetResource(theme.FolderIcon())
				} else {
					item.(*fyne.Container).Objects[0].(*widget.Icon).SetResource(theme.FileIcon())
				}
			}
			item.(*fyne.Container).Objects[0].Refresh()
			item.(*fyne.Container).Objects[1].Refresh()
		})
	fileList.OnSelected = func(id widget.ListItemID) {}

	return fileList, dir, nil
}

// ls_al. LINUX ls -Al
//
//goland:noinspection GoSnakeCaseUsage,SpellCheckingInspection
func ls_al(name string, info fs.FileInfo, err error) string {
	if err != nil {
		return fmt.Sprintf("Unable to get FileInfo for %s", name)
	}
	dfmt := defaultDateTimeFormat
	switch strings.ToUpper(dtFormat) {
	case "UNIX":
		dfmt = time.UnixDate
	case "RFC822":
		dfmt = time.RFC822
	case "RFC822Z":
		dfmt = time.RFC822Z
	case "RFC1123":
		dfmt = time.RFC1123
	case "RFC1123Z":
		dfmt = time.RFC1123Z
	case "RFC3339":
		dfmt = time.RFC3339
	}
	dt := info.ModTime().Format(dfmt)
	m := fmt.Sprintf("%s", info.Mode())
	f := "%4s %10d %s %s"
	return fmt.Sprintf(f, m[0:4], info.Size(), dt, name)
}
