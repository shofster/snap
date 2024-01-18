package app

/*

  File:    sys.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

var system *System
var doOnce sync.Once

type System struct {
	ARtype     string
	OStype     string
	AppName    string
	App        fyne.App
	Storage    string
	MainWindow fyne.Window
	TempDir    string
}

func GetSystem() *System {
	// Singleton System.
	doOnce.Do(func() {
		system = &System{
			AppName: "snap",
			App:     app.NewWithID("com.scsi.snap"),
			OStype:  runtime.GOOS,
			ARtype:  runtime.GOARCH}
		system.Storage = system.App.Storage().RootURI().Path()
		system.MainWindow = system.App.NewWindow("Photo SnapShot (by Bob)")
		dir, err := os.MkdirTemp("", system.AppName)
		if err == nil {
			system.TempDir = dir
		} else {
			name := filepath.Join(UserHomeDir(), "snapTemp")
			err := os.Mkdir(name, 0700)
			if err != nil {
				panic(errors.New(fmt.Sprintf("Unable to create TEMP %s", err)))
			}
			system.TempDir = name
		}
	})
	return system
}

// "toString"
func (s System) String() string {
	return fmt.Sprintf("Application %s, ARCH %s, OS %s\n STORAGE %s\n",
		s.AppName, s.ARtype, s.OStype, s.Storage)
}

// UserHomeDir provides the path to the user's Home directory (windows or other).
func UserHomeDir() string {
	home, err := os.UserHomeDir()
	if err == nil {
		return home
	}
	home = os.Getenv("HOMEDRIVE")
	if home == "" {
		panic("Unable to get User Home")
	}
	return home
}

func DeleteTemp() {
	if system.TempDir != "" {
		_ = os.RemoveAll(system.TempDir)
	}
}
