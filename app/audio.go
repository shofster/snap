package app

import (
	"fmt"
	"fyne.io/fyne/v2"
	"github.com/bogem/id3v2"
	"path/filepath"
	"strings"
	"unicode"
)

/*

  File:    audio.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: retrieve audio files info.
*/

type AudioInfo struct {
	title      string
	artist     string
	copyright  string
	collection string
	genre      string
	year       string
	length     uint16 // seconds
	audioType  string // MP3, ...
	mime       string
	cover      *fyne.StaticResource
}

func ID3Details(path string) (audioInfo AudioInfo, err error) {
	tag, err := id3v2.Open(path, id3v2.Options{Parse: true})
	if err == nil {
		defer func() {
			_ = tag.Close()
		}()
		picFrames := tag.GetFrames(tag.CommonID("Attached picture"))
		for _, f := range picFrames {
			pf, ok := f.(id3v2.PictureFrame)
			if ok { // some used jpg
				switch strings.ToLower(pf.MimeType) {
				case "image/jpeg", "image/jpg":
					audioInfo.mime = "jpeg"
				case "image/png":
					audioInfo.mime = "png"
				}
			}
			if audioInfo.mime != "" {
				audioInfo.cover = &fyne.StaticResource{
					StaticName:    filepath.Base(path),
					StaticContent: pf.Picture,
				}
			}
			break
		}
	}
	noHidden := func(s string) string {
		return strings.Map(func(r rune) rune {
			if r > unicode.MaxASCII {
				return -1
			}
			return r
		}, s)
	}
	audioInfo.artist = noHidden(tag.Artist())
	audioInfo.title = noHidden(tag.Title())
	audioInfo.collection = noHidden(tag.Album())
	audioInfo.year = noHidden(tag.Year())
	audioInfo.genre = noHidden(tag.Genre())

	return
}

const audioFmt = "Artist: %s\nTitle: %s\nAlbum: %s\n" + "Year: %s, Genre: %s"

func (a AudioInfo) String() string {
	return fmt.Sprintf(audioFmt,
		a.artist, a.title, a.collection, a.year, a.genre)
}
