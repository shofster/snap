package fileutil

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

/*

  File:    fileImpl.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*

 */

var _ fileView = (*fileImpl)(nil)

type fileImpl struct {
	path string
}

func (f fileImpl) Open(path string, sel FileSelectFilter) (*DirectoryEntry, error) {
	return openFileImpl(path, sel)
}

func (f fileImpl) Close() {
	log.Println("implement me")
}

func openFileImpl(path string, sel FileSelectFilter) (*DirectoryEntry, error) {
	de := NewDirectoryEntry(path)
	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			entries, e := os.ReadDir(path)
			for _, entry := range entries {
				// skip if a FILE matches the hidden expression
				if sel.Hidden != "" {
					//						if sel.Hidden != "" && !entry.IsDir() {
					match, e := regexp.MatchString(sel.Hidden, filepath.Base(entry.Name()))
					if match || e != nil {
						continue
					}
				}
				// skip non-directories if only want directories
				if sel.FileType == Dir && !entry.IsDir() {
					continue
				}
				if !entry.IsDir() && sel.Ext != "" {
					want := strings.ToUpper(filepath.Ext(sel.Ext))
					have := strings.ToUpper(filepath.Ext(entry.Name()))
					if want != ".*" && want != have {
						continue
					}
				}
				de.files = append(de.files, FileEntry{parent: path, entry: entry})
			}
			if e != nil {
				return de, e
			}
		} else {
			err = errors.New(fmt.Sprintf("path %s is NOT a Directory", path))
		}
	}
	return de, err
}
