package app

import (
	"fmt"
	"fyne.io/fyne/v2/dialog"
	"github.com/jung-kurt/gofpdf"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"snap/fileutil"
	"sort"
)

/*

  File:    pdf.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: Create a PDF file from list of paths.
*/

const maxPhotoCols = 5
const maxPhotoRows = 7

func getAllFiles(path string) []string {
	files := make([]string, 0)
	entries, _ := os.ReadDir(path)
	for _, entry := range entries {
		// skip hidden
		match, e := regexp.MatchString(fileutil.DefaultHiddenFiles, filepath.Base(entry.Name()))
		if match || e != nil {
			continue
		}
		if entry.IsDir() || entry.Type() == 0 {
			files = append(files, entry.Name())
		}
	}
	return files
}

func CreatePDF(dirs []string, file string) (err error) {
	pdf := gofpdf.New("P", "pt", "Letter", "")
	pdf.SetMargins(15, 15, 15)
	for _, dir := range dirs {
		e := showPdf(dir, pdf)
		if e != nil || pdf.Err() {
			log.Printf("CreatePDF error: %s\n %s\n  %s\n", e, pdf.Error(), dir)
			continue
		}
	}
	err = pdf.OutputFileAndClose(file)
	if pdf.Err() {
		log.Printf("pdf.OutputFileAndClose error: %s\n", pdf.Error())
	}
	return
}

func showPdf(dir string, pdf *gofpdf.Fpdf) (err error) {
	files := getAllFiles(dir)
	err = buildPDF(dir, pdf, files)
	if err != nil {
		dialog.ShowError(err, GetSystem().MainWindow)
	}
	return err
}

func buildPDF(dir string, pdf *gofpdf.Fpdf, sorted []string) error {
	var n int
	header := func() {
		n = 0
		pdf.AddPage()
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(570, 12, fmt.Sprintf("%s", dir), "", 0, "CM", false, 0, "")
		pdf.SetFont("Arial", "", 8)
	}
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})
	xScale := 100
	yScale := 100
	for _, s := range sorted {
		// get generic path for image
		path, err := ImageResourcePath(dir, s)
		if err != nil {
			log.Println("Got ImageResourcePath error ", s, err)
			continue
		}
		if n%(maxPhotoRows*maxPhotoCols) == 0 {
			header()
		}
		row := n / maxPhotoCols
		col := n - (row * maxPhotoCols)
		n++
		// ImageOptions(src, x, y, width, height, flow, options, link, linkStr)
		pdf.ImageOptions(
			path,
			float64(col*xScale+50), float64(row*yScale+50),
			80, 80,
			false,
			gofpdf.ImageOptions{ReadDpi: true},
			0,
			"",
		)
		if pdf.Err() {
			log.Printf("buildPDF error: %s\n  %s\n", pdf.Error(), path)
			continue
		}
		// limit length to avoid collision
		name := s //strings.TrimSuffix(s, filepath.Ext(s))
		if len(name) > 24 {
			name = name[len(s)-24:]
		}
		pdf.Text(float64(col*xScale+45), float64(row*yScale+140), name)
	}
	return nil
}
