package fileutil

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

/*

  File:    directoryView.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Create a DirectoryEntry (list of FileEntry) from the path.
*/

func NewDirectoryView(path string, sel FileSelectFilter) (*DirectoryEntry, error) {
	ext := strings.ToUpper(filepath.Ext(path))
	// "comma ok" form
	var p ProtocolType
	p, ok := extMap[ext]
	if !ok {
		p = FILE
	}
	var de *DirectoryEntry
	var err error
	switch p {
	case FILE:
		de, err = openFileImpl(path, sel)
	default:
		return nil, errors.New(fmt.Sprintf("Unable to find protocol for %s", ext))
	}
	if err != nil {
		return de, err
	}
	if sel.Date {
		if sel.Descending {
			sortDateTimeDescending(de.files)
		} else {
			sortDateTimeAscending(de.files)
		}
	} else {
		if sel.Descending {
			sortNameDescending(de.files)
		} else {
			sortNameAscending(de.files)
		}
	}
	return de, nil
}

//goland:noinspection GoUnusedFunction
func sortNameAscending(slice []FileEntry) {
	sort.SliceStable(slice, func(i, j int) bool {
		return strings.ToLower(slice[i].DisplayName()) < strings.ToLower(slice[j].DisplayName())
	})
}
func sortNameDescending(slice []FileEntry) {
	sort.SliceStable(slice, func(i, j int) bool {
		return !(strings.ToLower(slice[i].DisplayName()) < strings.ToLower(slice[j].DisplayName()))
	})
}

func sortDateTimeAscending(slice []FileEntry) {
	sort.SliceStable(slice, func(i, j int) bool {
		return slice[i].ModTime().Before(slice[j].ModTime())
	})
}
func sortDateTimeDescending(slice []FileEntry) {
	sort.SliceStable(slice, func(i, j int) bool {
		return slice[i].ModTime().After(slice[j].ModTime())
	})
}
