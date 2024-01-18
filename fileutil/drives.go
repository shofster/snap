//go:build !windows
// +build !windows

package fileutil

import (
	"os"
	"syscall"
)

/*

  File:    drives.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:
*/
//goland:noinspection GoUnusedFunction
func getDrives() (drives []string, err error) {
	dirs := []string{"/", "/etc", "/home", "/media", "/mnt", "/tmp", "/usr"}
	for _, dir := range dirs {
		f, err := os.Open(dir)
		if err == nil {
			if err == nil {
				drives = append(drives, dir)
				_ = f.Close()
			}
		}
	}
	return
}

// disk usage of path/disk
func getDiskUsage(vol string) (disk DiskUsage) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(vol, &fs)
	if err == nil {
		disk.All = fs.Blocks * uint64(fs.Bsize)
		disk.Avail = fs.Bavail * uint64(fs.Bsize)
		disk.Free = fs.Bfree * uint64(fs.Bsize)
		disk.Used = disk.All - disk.Free
	}
	return
}
