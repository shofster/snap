package fileutil

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

/*

  File:    places.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:
*/

// LoadPlaces enumerates the logical drives / standard places on the system.
//
//goland:noinspection GoUnusedExportedFunction
func LoadPlaces() (drives []string) {
	switch runtime.GOOS {
	case "windows", "linux":
		drives, _ = getDrives()
	default:
		drives = append(drives, "/")
	}
	return
}

type DiskUsage struct {
	All   uint64
	Used  uint64
	Free  uint64
	Avail uint64
}

//goland:noinspection GoUnusedExportedFunction
func PrettyDiskSize(s uint64) string {
	var B uint64 = 1073741824
	if s > B {
		return fmt.Sprintf("%4.1fGB", float32(s)/float32(B))
	}
	B /= 1024
	if s > B {
		return fmt.Sprintf("%4.1fMB", float32(s/1024./B))
	}
	B /= 1024
	if s > B {
		return fmt.Sprintf("%4.1fKB", float32(s/1024./1024./B))
	}
	return fmt.Sprintf("%db", s)
}

//goland:noinspection GoUnusedExportedFunction
func GetDiskUsage(filename string) DiskUsage {
	switch runtime.GOOS {
	case "windows", "linux":
		return getDiskUsage(filename)
	default:
		return DiskUsage{}
	}
}

// TooManyFilesError is a custom error when trying to display 3000 files bogs down
// error is given, but first 256 file names are returned
//
//goland:noinspection GoUnusedExportedFunction
type TooManyFilesError struct {
	Path  string
	Count int
}

func (e *TooManyFilesError) Error() string {
	return fmt.Sprintf("%d", e.Count)
}

// DirContents gets a []os.FileInfo's from a path.
//
//goland:noinspection GoUnusedExportedFunction
func DirContents(dir string, maxFiles int) ([]os.FileInfo, error) {
	p := make([]os.FileInfo, 0)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return p, errors.New("sys.DirContents: " + dir + "," + err.Error())
	}
	d, err := os.ReadDir(dir)
	if len(d) > maxFiles {
		err = &TooManyFilesError{
			Count: len(d)}
		d = d[0 : maxFiles-1]
	}
	for _, de := range d {
		fi, _ := de.Info()
		p = append(p, fi)
	}
	return p, err
}

// PathContents gets the names of all files in a directory/
//
//goland:noinspection GoUnusedExportedFunction
func PathContents(path string) ([]string, error) {
	_, x := os.Stat(path)
	if os.IsNotExist(x) {
		return nil, errors.New("sys.PathContents: " + path + "," + x.Error())
	}
	fi, err := os.ReadDir(path)
	p := make([]string, 0, len(fi))
	for _, file := range fi {
		p = append(p, filepath.Join(path, file.Name()))
	}
	return p, err
}

// DirTreeSize gets (recursive) count and size of all files.
//
//goland:noinspection GoUnusedExportedFunction
func DirTreeSize(path string, sz *uint64) int {
	// I have seen the 1.15 version go belly-up for no reason :-((
	// catch and release
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error in DirTreeSize", r)
		}
	}()
	n := 0
	_ = filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
		info, err := d.Info()
		if info != nil && err == nil {
			n++
			*(sz) += uint64(info.Size())
		}
		return nil
	})
	return n
}

// DirTreeList gets (recursive) names of all files and directories.
//
//goland:noinspection GoUnusedExportedFunction
func DirTreeList(path string, f func(string) error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error in DirTreeList", r)
		}
	}()
	_ = filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
		return f(p)
	})
}

// Rename changes just the name a file
//
//goland:noinspection GoUnusedExportedFunction
func Rename(path string, base string) error {
	err := os.Rename(path, filepath.Join(filepath.Dir(path), base))
	return err
}

// Replace replaces a file in a new directory.
//
//goland:noinspection GoUnusedExportedFunction
func Replace(path string, dir string) error {
	base := filepath.Base(path)
	err := os.Rename(path, filepath.Join(dir, base))
	return err
}

type PlaceType int

const (
	EmptyDirPlace = iota
	FilePlace
	DirPlace
	OtherPlace
)

// a "toString" of the PlaceType
func (t PlaceType) String() string {
	return [...]string{"EmptyDir", "File", "Dir", "Other"}[t]
}

// GetPlaceType return a limited PlaceType.
//
//goland:noinspection GoUnusedExportedFunction
func GetPlaceType(place string) (PlaceType, error) {
	fi, err := os.Stat(place)
	if err != nil {
		return OtherPlace, err
	}
	if fi.IsDir() {
		sub, ep := PathContents(place)
		if ep != nil {
			return DirPlace, ep
		}
		if len(sub) == 0 {
			return EmptyDirPlace, nil
		} else {
			return DirPlace, nil
		}
	}
	if fi.Mode().IsRegular() {
		return FilePlace, nil
	}
	return OtherPlace, nil
}

// CopyPlace copies a source file from the destination directory.
// and reset the new file's time to the original.
//
//goland:noinspection GoUnusedExportedFunction
func CopyPlace(source, destination string, time time.Time) (uint64, error) {
	in, err := os.Open(source)
	if err != nil {
		return 0, err
	}
	defer func(in *os.File) {
		if e := in.Close(); e != nil {
		}
	}(in)
	out, err1 := os.Create(destination)
	if err1 != nil {
		return 0, err1
	}
	defer func(out *os.File) {
		if e := out.Close(); e != nil {
			log.Println("close copy close", e, time)
		}
		_ = os.Chtimes(destination, time, time)
	}(out)
	nBytes, err2 := io.Copy(out, in)
	if err2 != nil {
		return uint64(nBytes), err2
	}
	return uint64(nBytes), err2
}
