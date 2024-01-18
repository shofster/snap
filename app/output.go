package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"snap/fileutil"
)

/*

  File:    output.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: Handle file output selection of PDFs.
*/

func GetNextOutputPath(window fyne.Window,
	last binding.ExternalString,
	cb func(string)) {

	fs := fileutil.FileSelectFilter{
		Title:      "Output PDF File Path",
		FileType:   fileutil.File,
		FileSelect: fileutil.Save,
		Multiple:   false,
		Hidden:     fileutil.DefaultHiddenFiles,
		Ext:        "*.pdf",
		Date:       false,
		DtFormat:   "",
		Descending: false,
	}
	_ = fileutil.FileSelect(fs, last, window, func(files []string) {
		if len(files) == 1 {
			cb(files[0])
		}
	})

}
