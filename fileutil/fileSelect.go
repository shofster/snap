package fileutil

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"io/fs"
	"os"
	"path/filepath"
)

/*

  File:    fileSelect.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*

 */

func FileSelect(sel FileSelectFilter, lastDir binding.String, window fyne.Window, cb func([]string)) *widget.PopUp {
	p := newPanel(sel, lastDir, cb)
	title := widget.NewLabelWithStyle(makeTitle(sel), fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true, Italic: true})
	content := container.NewBorder(title, nil, nil, nil, p.content)
	p.win = window
	p.popup = widget.NewModalPopUp(content, window.Canvas())
	p.popup.ShowAtPosition(fyne.NewPos(0, 0))
	p.popup.Resize(fyne.NewSize(500, 400))
	p.popup.Refresh()
	return p.popup
}

type panel struct {
	sel           FileSelectFilter
	cb            func([]string)
	lastDir       binding.String
	selected      []string
	parent        binding.String
	addDir        *widget.Button
	currentDir    binding.String
	previousDir   binding.String
	filename      binding.String
	filenameEntry *widget.Entry
	exts          *widget.Select
	ext           string
	defaultExt    string
	//	list            *widget.List
	dir             *DirectoryEntry
	win             fyne.Window
	popup           *widget.PopUp
	content         *container.Split
	placeContainer  *fyne.Container
	selectContainer *fyne.Container
	listContainer   *fyne.Container
}

func newPanel(sel FileSelectFilter, lastDir binding.String, cb func([]string)) *panel {
	p := &panel{sel: sel,
		cb:          cb,
		lastDir:     lastDir,
		parent:      binding.NewString(),
		currentDir:  binding.NewString(),
		previousDir: binding.NewString(),
		filename:    binding.NewString(),
	}
	p.selected = make([]string, 0)
	p.initPanel(lastDir)
	return p
}
func (p *panel) initPanel(lastDir binding.String) {
	p.placeContainer = p.buildPlaces(lastDir)
	currentDir := fmt.Sprintf("%30sChoose Folder%30s", " ", " ")
	filename := ""
	_ = p.currentDir.Set(currentDir)
	currentLabel := widget.NewLabelWithData(p.currentDir)
	_ = p.previousDir.Set("")
	previousButton := widget.NewButtonWithIcon("", theme.MoveUpIcon(), func() {
		parent, _ := p.parent.Get()
		p.showDir(filepath.Dir(parent))
	})
	previousButton.IconPlacement = widget.ButtonIconLeadingText
	p.parent.AddListener(binding.NewDataListener(func() {
		p, _ := p.parent.Get()
		previousButton.Text = filepath.Dir(p)
		previousButton.Refresh()
	}))
	p.addDir = widget.NewButtonWithIcon("", theme.FolderNewIcon(), func() {
		parent, _ := p.parent.Get()
		AskFolder(p.win, parent, func(name string) {
			if name != "" {
				path := filepath.Join(parent, name)
				err := os.MkdirAll(path, os.ModePerm)
				if err != nil {
					p.showError(err)
				} else {
					p.showDir(path)
				}
			}
		})
	})
	p.addDir.Disable()
	dirControl := container.NewBorder(nil, nil, currentLabel, p.addDir)
	top := container.NewVBox(dirControl, previousButton)

	list := widget.NewList(
		func() int {
			return 1
		},
		func() fyne.CanvasObject {
			return widget.NewIcon(theme.FileIcon())
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
		})
	p.listContainer = container.NewStack(list)
	_ = p.filename.Set(filename)
	p.filenameEntry = widget.NewEntryWithData(p.filename)
	p.filenameEntry.OnSubmitted = func(name string) {
		if name == "" {
			p.popup.Hide()
			p.cb(make([]string, 0))
		}
		p.checkSaveFileName(name)
	}
	p.filenameEntry.Disable()
	p.exts = widget.NewSelect(make([]string, 0), func(_ string) {})
	if p.sel.Ext != "" {
		p.exts.Options = append(p.exts.Options, p.sel.Ext)
		p.defaultExt = p.sel.Ext
	}
	all := "*.*"
	p.exts.Options = append(p.exts.Options, all)
	p.exts.OnChanged = func(ext string) {
		p.sel.Ext = ext
		parent, _ := p.parent.Get()
		p.showList(parent)
	}
	saveas := container.NewVBox(p.filenameEntry, p.exts)
	done := widget.NewButtonWithIcon("Done", theme.ConfirmIcon(), func() {
		if p.sel.FileSelect == Save {
			name, _ := p.filename.Get()
			if name != "" {
				p.checkSaveFileName(name)
				return
			}
			p.popup.Hide()
			p.cb(make([]string, 0))
			return
		}
		p.popup.Hide()
		p.cb(p.selected)
	})
	done.Importance = widget.HighImportance
	cancel := widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
		p.popup.Hide()
		p.cb(make([]string, 0))
	})
	canceldone := container.NewHBox(cancel, done)
	bottom := container.NewVBox()
	if p.sel.FileSelect == Save {
		bottom.Objects = append(bottom.Objects, saveas)
	}
	bottom.Objects = append(bottom.Objects, container.NewBorder(nil, nil, nil, canceldone))

	p.selectContainer = container.NewBorder(top, bottom, nil, nil, p.listContainer)
	p.content = container.NewHSplit(p.placeContainer, p.selectContainer)
	p.content.Offset = .3
}
func (p *panel) checkSaveFileName(name string) {
	x := filepath.Ext(name)
	if x == "" && p.defaultExt != "" { // add extension
		name += filepath.Ext(p.defaultExt)
		_ = p.filename.Set(name)
	}
	// see if overwriting
	parent, _ := p.parent.Get()
	path := filepath.Join(parent, name)
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			// doesn't exist, done - send result
			p.popup.Hide()
			s := make([]string, 0)
			s = append(s, path)
			p.cb(s)
			return
		}
	}
	VerifyOverwrite(p.win, name, func(b bool) {
		if b {
			p.popup.Hide()
			if len(p.selected) < 1 {
				p.selected = append(p.selected, name)
			}
			p.cb(p.selected)
		} else {
			_ = p.filename.Set("")
		}
	})
	return
}
func (p *panel) showError(err error) {
	list := widget.NewList(
		func() int {
			return 1
		},
		func() fyne.CanvasObject {
			return widget.NewIcon(theme.ErrorIcon())
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
		})
	p.selected = p.selected[:0]
	p.addDir.Disable()
	list.OnSelected = nil
	list.OnUnselected = nil
	p.listContainer.Objects = p.listContainer.Objects[:0]
	p.listContainer.Objects = append(p.listContainer.Objects, list)
	list.Refresh()
	_ = p.parent.Set("")
	_ = p.currentDir.Set(err.Error())
	p.filenameEntry.Disable()
}
func (p *panel) showDir(newPlace string) {
	p.showList(newPlace)
}
func (p *panel) showList(newPlace string) {
	p.selected = p.selected[:0]
	p.filenameEntry.Enable()
	p.addDir.Enable()
	if p.lastDir != nil {
		_ = p.lastDir.Set(newPlace)
	}
	fa := FileSelectAction{
		OnClick: func(file FileEntry) { // single click
			path := file.Name()
			//			.Printf("FileSelect  selected %v ", file)
			fdir, _ := file.Info()
			if (p.sel.FileType == Any) ||
				(p.sel.FileType == File && !fdir.IsDir()) ||
				(p.sel.FileType == Dir && fdir.IsDir()) {
				p.selected = Remove(p.selected, path)
				if file.IsSelected() {
					_ = p.filename.Set(filepath.Base(path))
					selectMax := 1
					if p.sel.Multiple {
						selectMax = 1000
					}
					p.selected = Add(p.selected, path, selectMax)
				} else {
					_ = p.filename.Set("")
				}
			}
		},
		OnDoubleClick: func(entry FileEntry) {
			if entry.IsDir() {
				path := entry.Name()
				p.selected = Remove(p.selected, path)
				if entry.IsSelected() {
					selectMax := 1
					if p.sel.Multiple {
						selectMax = 1000
					}
					p.selected = Add(p.selected, path, selectMax)
				}
				p.showDir(path)
			}
		},
	}

	list, dir, err := NewFileList(newPlace, p.sel, fa,
		func(name string, info fs.FileInfo, err error) string { // pretty
			if err != nil {
				return fmt.Sprintf("Unable to get FileInfo for %s", name)
			}
			dt := info.ModTime().Format("01/02/06")
			m := fmt.Sprintf("%s", info.Mode())
			f := "%4s %8s %s"
			return fmt.Sprintf(f, m[0:4], dt, name)
		})
	if err != nil {
		p.showError(err)
		return
	}
	p.listContainer.Objects = p.listContainer.Objects[:0]
	p.listContainer.Objects = append(p.listContainer.Objects, list)
	p.dir = dir
	list.Refresh()
	p.selectContainer.Refresh()
	_ = p.parent.Set(newPlace)
	place := fmt.Sprintf("%30s%s%30s", " ", filepath.Base(newPlace), " ")
	_ = p.currentDir.Set(place)
}
func (p *panel) buildPlaces(lastDir binding.String) *fyne.Container {
	// make a Button list of the available Places
	places := LoadPlaces()
	placeContainer := container.NewVBox()
	n := 0
	if lastDir != nil {
		last, _ := lastDir.Get()
		if last != "" {
			n++
			placeContainer.Objects = append(placeContainer.Objects,
				widget.NewButtonWithIcon(" "+last, theme.ContentRedoIcon(), func() {
					p.showList(last)
				}))
		}
	}
	placeContainer.Objects = append(placeContainer.Objects, widget.NewLabel(""))
	n++
	home, err := os.UserHomeDir()
	if err == nil {
		n++
		placeContainer.Objects = append(placeContainer.Objects,
			widget.NewButtonWithIcon(home, theme.FolderOpenIcon(), func() {
				p.showList(home)
			}))
	}
	for _, place := range places {
		n++
		pl := place
		placeContainer.Objects = append(placeContainer.Objects,
			widget.NewButtonWithIcon(place, theme.FolderOpenIcon(), func() {
				p.showList(pl)
			}))
	}
	for n < 13 {
		placeContainer.Objects = append(placeContainer.Objects, widget.NewLabel(""))
		n++
	}
	return placeContainer
}
func makeTitle(sel FileSelectFilter) string {
	if sel.Title != "" {
		return sel.Title
	}
	title := "Select "
	if sel.Multiple {
		title += "multiple "
	} else {
		title += "one "
	}
	switch sel.FileType {
	case Any:
		title += "File or Folder"
	case File:
		title += "File"
	case Dir:
		title += "Folder"
	}
	if sel.Multiple {
		title += "(s) "
	} else {
		title += " "
	}
	switch sel.FileSelect {
	case Open:
		title += "<Open>"
	case Save:
		title += "<Save>"
	}
	return title
}
