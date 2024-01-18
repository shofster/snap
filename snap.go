package main

/*

  File:    snap.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.
*/
/*
	Description: Desktop program to generate PDF images from files.
*/

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"snap/app"
	"sort"
	"strings"
	"time"
)

/*

  File:    snap.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: GOlang / fyne.io Program to generate PDF "contact sheets"

  To create resource: fyne bundle -o images.go --pkg app images

*/

import (
	"snap/element"
)

func main() {
	// system has global variables
	system := app.GetSystem()
	defer func() { // remove TempDir, if normal exit
		app.DeleteTemp()
	}()

	system.MainWindow.SetIcon(resourcePDFphotoPng)
	system.App.Settings().SetTheme(element.NewTheme(system.App.Preferences()))

	file := filepath.Join(system.Storage, system.AppName) + ".log"
	var logger *os.File
	if logf, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666); err == nil {
		logger = logf
		defer func() {
			_ = logger.Close()
		}()
		log.SetOutput(logger)
	}
	log.Printf("snap - System:  %s\n", system)

	system.App.Settings().SetTheme(element.NewTheme(system.App.Preferences()))
	console := element.NewConsole(system.MainWindow, "command", 15)
	content := container.NewBorder(nil, nil, nil, nil, console.Content)

	prefs := system.App.Preferences()
	lastPath := prefs.StringWithFallback("last", app.UserHomeDir())
	boundLast := binding.BindString(&lastPath)
	pdfPath := prefs.StringWithFallback("pdf", app.UserHomeDir())
	boundPDF := binding.BindString(&pdfPath)

	paths := make([]string, 0)
	// unique list of sorted paths
	var addPath = func(p string) {
		for i, v := range paths {
			if v == p {
				paths = append(paths[:i], paths[i+1:]...)
				break
			}
		}
		paths = append(paths, p)
		sort.SliceStable(paths, func(i, j int) bool {
			return paths[i] < paths[j]
		})
		app.ShowCount(console, len(paths))
	}

	var addAction = func() {
		app.GetNextInputPath(system.MainWindow, boundLast, func(d string) {
			console.Speak(fmt.Sprintf("** Added Path: %s", d))
			prefs.SetString("last", lastPath)
			addPath(d)
		})
		console.Focus()
	}
	var urlFromString = func(str string) (*url.URL, error) {
		u, err := url.Parse(str)
		if err != nil {
			return nil, err
		}
		return u, nil
	}
	var browse = func(file string) error {
		u, err := urlFromString(file)
		if err != nil {
			return err
		}
		return fyne.CurrentApp().OpenURL(u)
	}

	var pdfAction = func() {
		if len(paths) < 1 {
			app.ShowCount(console, len(paths))
			app.ErrorText(console, "NO PATHS. Use \"(a) Add Path...\"")
		} else {
			app.GetNextOutputPath(system.MainWindow, boundPDF, func(d string) {
				prefs.SetString("pdf", pdfPath)
				err := app.CreatePDF(paths, d)
				if err != nil {
					app.ErrorText(console, fmt.Sprintf("%v", err))
				} else {
					console.Speak(fmt.Sprintf("** PDF Written: %s", d))
					err = browse(d)
					if err != nil {
						app.ErrorText(console, fmt.Sprintf("%v", err))
					}
				}
				console.Focus()
			})
		}
	}

	photo := canvas.NewImageFromResource(resourcePDFphotoPng)
	photo.FillMode = canvas.ImageFillContain
	splash := container.NewStack(photo)
	system.MainWindow.SetContent(splash)
	system.MainWindow.Resize(fyne.NewSize(800, 600))
	system.MainWindow.SetOnClosed(func() {
	})

	// process typed commands
	var action = func(typed string) {
		switch strings.ToLower(typed) {
		case "x", "exit", "q", "quit":
			system.App.Quit()
		case "a", "add":
			addAction()
		case "p", "pdf":
			pdfAction()
		case "l", "list":
			app.ShowText(console, "Paths:", paths)
		case "c", "clear":
			paths = nil
			app.ShowCount(console, len(paths))
		default:
			app.ShowText(console, "Valid Commands:", help)
		}
		console.Focus()
	}
	go func() { // show photo for a bit
		time.Sleep(time.Millisecond * 3000)
		splash.Objects[0] = content
		splash.Refresh()
		app.ShowText(console, "", first)
		app.ShowText(console, "Valid Commands", help)
		// start the input console
		go app.Input(console, action)
		console.Focus()
	}()

	system.MainWindow.ShowAndRun()
	system.App.Quit()
}

var first = []string{
	"\nSNAP (by Bob) is a program to create a PDF file",
	"  showing thumbnails of  files on your hard drive.",
}
var help = []string{
	"(x) eXit  or (q) Quit",
	"(l) List PATHs",
	"(a) Add PATHs ...",
	"(c) Clear PATHs",
	"(p) generate PDF ...",
	"(h) Help",
}
