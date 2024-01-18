package fileutil

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"time"
)

/*

  File:    entry.go
Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Directory and File Entries.
*/

const DirSeparator string = "\n"

type DirectoryEntry struct {
	base  string
	url   string // concatenated, logical path (may be archive "file")
	path  string // physical path
	files []FileEntry
}

func NewDirectoryEntry(url string) (de *DirectoryEntry) {
	de = &DirectoryEntry{url: url, path: url}
	dirs := strings.Split(url, DirSeparator)
	de.base = filepath.Base(filepath.Ext(dirs[len(dirs)-1]))
	de.files = make([]FileEntry, 0)
	return
}
func (d *DirectoryEntry) Name() string {
	return d.path
}
func (d *DirectoryEntry) Count() int {
	return len(d.files)
}
func (d *DirectoryEntry) File(ix int) *FileEntry {
	return &d.files[ix]
}
func (d *DirectoryEntry) MaxDisplayName() (max int) {
	for _, file := range d.files {
		size := len(file.DisplayName())
		if size > max {
			max = size
		}
	}
	return
}
func (d *DirectoryEntry) GetFilteredGroup(_ FileSelectFilter) DirectoryEntry {
	group := DirectoryEntry{base: d.base, url: d.url, path: d.path}
	group.files = make([]FileEntry, 0)
	return group
}
func (d *DirectoryEntry) SelectAll(s bool) {
	for _, file := range d.files {
		file.selected = s
	}
}
func (d *DirectoryEntry) GetSelected() (files []FileEntry) {
	files = make([]FileEntry, 0)
	for _, file := range d.files {
		if file.selected {
			files = append(files, file)
		}
	}
	return
}

type FileEntry struct {
	selected bool
	parent   string
	index    int
	entry    fs.DirEntry
}

func (f *FileEntry) String() string {
	return fmt.Sprintf("Name: %s, selected %t", filepath.Base(f.entry.Name()), f.selected)
}
func (f *FileEntry) Name() string {
	return filepath.Join(f.parent, f.entry.Name())
}
func (f *FileEntry) Index() int {
	return f.index
}
func (f *FileEntry) DisplayName() string {
	return f.entry.Name()
}
func (f *FileEntry) Info() (fs.FileInfo, error) {
	return f.entry.Info()
}
func (f *FileEntry) ModTime() time.Time {
	i, err := f.Info()
	if err != nil {
		return time.Now()
	}
	return i.ModTime()
}
func (f *FileEntry) IsDir() bool {
	return f.entry.IsDir()
}
func (f *FileEntry) IsSelected() bool {
	return f.selected
}
func (f *FileEntry) SetSelected(sel bool) {
	f.selected = sel
}
