package fileutil

/*

  File:    etc.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Miscellaneous methods.
*/

//goland:noinspection GoUnusedGlobalVariable
var DefaultHiddenFiles = "(^[^\\w].+)|(.+\\.bak)$"

// StringList is the type of array
type StringList []string

// Remove deletes a string from a []string.
func Remove(sl StringList, r string) []string {
	for i, v := range sl {
		if v == r {
			return append(sl[:i], sl[i+1:]...)
		}
	}
	return sl
}

// Add appends a string to the end of a []string.
func Add(sl []string, a string, max int) []string {
	sl = append([]string{a}, sl...)
	if len(sl) > max {
		sl = sl[0:max]
	}
	return sl
}
