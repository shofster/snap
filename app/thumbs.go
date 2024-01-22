package app

/*

  File:    input.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: Manage images for the PDF.
*/

import (
	"fmt"
	"fyne.io/fyne/v2"
	"os"
	"path/filepath"
	"strings"
)

// ImageResourcePath insures a file with image data exists. internal images
//
//	are in fyne storage folder.
func ImageResourcePath(dir, name string) (path string, err error) {
	path = filepath.Join(dir, name)
	switch strings.ToLower(filepath.Ext(path)) {
	case ".jpg", ".jpeg", ".png", ".gif":
		return
		//case ".heic":
		//	if e := CheckHeic(path); e == nil {
		//		return getHEICImagePath(path)
		//	}
	}
	switch ExtensionType(path) {
	case AudioExt:
		details, err := ID3Details(path)
		if err == nil && details.cover != nil {
			p := details.cover.StaticName
			details.cover.StaticName = fmt.Sprintf("%s.%s",
				strings.TrimSuffix(p, filepath.Ext(p)),
				details.mime)
			return getTempImagePath(details.cover)
		}
	}
	return getResourceImagePath(imageResourceMap[ExtensionType(path)])
}

//  to create resource: fyne bundle -o images.go --pkg app images

const AudioExt = "audio"
const BitmapExt = "bitMap"
const CameraExt = "camera"
const DocExt = "doc"
const ExeExt = "exe"
const PdfExt = "pdf"
const ZipExt = "zip"
const VideoExt = "video"
const UnknownExt = "unknown"
const FolderExt = "folder"
const HtmlExt = "html"

func ExtensionType(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return UnknownExt
	}
	if info.IsDir() {
		return FolderExt
	}
	switch strings.ToLower(filepath.Ext(path)) {
	case ".jpg", ".jpeg", ".png", ".gif", ".kdc", ".sfw", ".raw", ".heic":
		return CameraExt
	case ".mp3", ".m4a", ".flac", ".wav", ".wma", ".aac", ".ogg":
		return AudioExt
	case ".pdf":
		return PdfExt
	case ".mp4", ".m4v", ".mov", ".wmv", ".avi", ".avchd", ".hevc",
		".flv", ".f4v", ".swf", ".3gp", ".mpeg", ".mpg":
		return VideoExt
	case ".zip", ".gz", "tgz", ".gzip", ".7z", ".jar":
		return ZipExt
	case ".exe", ".com", ".bat", ".cmd", ".sh", ".bin":
		return ExeExt
	case ".bmp", ".tiff", ".tif":
		return BitmapExt
	case ".html", ".htm":
		return HtmlExt
	case ".dat", ".txt", ".csv",
		".xls", ".xlsx", ".doc", ".docx",
		".odt", ".ods", "odp", ".odg":
		return DocExt
	}
	return UnknownExt
}

var imageResourceMap = map[string]*fyne.StaticResource{
	AudioExt:   resourceAudioPng,
	BitmapExt:  resourceBitmapJpg,
	CameraExt:  resourceCameraPng,
	DocExt:     resourceDocPng,
	ExeExt:     resourceExePng,
	PdfExt:     resourcePdfPng,
	ZipExt:     resourceZipPng,
	VideoExt:   resourceVideoJpg,
	UnknownExt: resourceUnknownPng,
	FolderExt:  resourceDirPng,
	HtmlExt:    resourceHtmlPng,
}

func getResourceImagePath(resource *fyne.StaticResource) (string, error) {
	name := resource.StaticName
	content := resource.StaticContent
	path := filepath.Join(GetSystem().Storage, name)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.WriteFile(path, content, 0644)
	}
	return path, err
}
func getTempImagePath(resource *fyne.StaticResource) (string, error) {
	name := resource.StaticName
	content := resource.StaticContent
	path := filepath.Join(GetSystem().TempDir, name)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.WriteFile(path, content, 0644)
	}
	return path, err
}

/*
func getHEICImagePath(path string) (string, error) {
	name := filepath.Base(path)
	name = strings.TrimSuffix(name, filepath.Ext(name))
	o := fmt.Sprintf("%s.jpg",
		filepath.Join(GetSystem().TempDir, name))
	//	log.Printf("HEIC: %s\n   %s\n", path, o)
	return o, ConvertHeicToJpg(path, o)
}
*/
