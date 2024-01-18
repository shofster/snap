package app

/*

  File:    input.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: Handle adding paths for snap shots.
*/

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"snap/element"
	"snap/fileutil"
)

func Input(c *element.Console, cmd func(t string)) {
	for {
		p := c.Ask("required!")
		cmd(p)
	}
}

func ShowCount(c *element.Console, count int) {
	c.Speak(fmt.Sprintf("Path count is %d", count))
}
func ShowText(c *element.Console, h string, txt []string) {
	c.Speak("\n" + h)
	for _, t := range txt {
		c.Speak(t)
	}
}
func ErrorText(c *element.Console, txt string) {
	c.Speak(fmt.Sprintf("!! Error: %s", txt))
}

func GetNextInputPath(window fyne.Window,
	last binding.ExternalString,
	cb func(string)) {

	fs := fileutil.FileSelectFilter{
		Title:      "Next Input Path",
		FileType:   fileutil.Dir,
		FileSelect: fileutil.Open,
		Multiple:   true,
		Hidden:     fileutil.DefaultHiddenFiles,
		Ext:        "",
		Date:       false,
		DtFormat:   "",
		Descending: false,
	}
	_ = fileutil.FileSelect(fs, last, window, func(dirs []string) {
		for _, dir := range dirs {
			cb(dir)
		}
	})
}
