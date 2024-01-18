//go:build windows

package fileutil

import (
	"golang.org/x/sys/windows"
	"syscall"
)

/*

  File:    windrives.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:
*/

var winDrives = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

// Unused unexported function
//
//goland:noinspection GoUnusedFunction
func getDrives() (drives []string, err error) {
	var sysHandle syscall.Handle
	sysHandle, _ = syscall.LoadLibrary("kernel32.dll")
	var logicalDrivesHandle uintptr
	logicalDrivesHandle, _ = syscall.GetProcAddress(sysHandle, "GetLogicalDrives")
	ret, _, _ := syscall.SyscallN(logicalDrivesHandle, 0, 0, 0, 0)
	for i := range winDrives {
		if ret&1 == 1 {
			drives = append(drives, winDrives[i]+":\\")
		}
		ret >>= 1
	}
	return
}

func getDiskUsage(vol string) (disk DiskUsage) {
	u16fname, err := syscall.UTF16FromString(vol)
	if err == nil {
		_ = windows.GetDiskFreeSpaceEx(&u16fname[0], &disk.Avail, &disk.All, &disk.Free)
		disk.Used = disk.All - disk.Free
	}
	return
}
