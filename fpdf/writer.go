package fpdf

import (
	"bytes"
	"io"
	"os"

	"github.com/phpdave11/gofpdf"
	"github.com/phpdave11/gofpdf/contrib/gofpdi"
	"github.com/x64c/gw/pdfs"
	"github.com/x64c/gw/rw"
)

// Writer - simple PDF writer using gofpdf (Go port of FPDF)
// Supported unit: "pt" only
// Currently Custom Size Not Supported: "Letter" and "A4" Only
// ToDo: Support custom fonts
type Writer struct {
	pdfs.Writer[int] // [Embedded Interface] [To Implement]

	paperSize   pdfs.PaperSize
	orientation string

	impl      *gofpdf.Fpdf
	templates *pdfs.TemplateStore[int]
}

func NewWriter(paperSize pdfs.PaperSize, orientation string) *Writer {
	return &Writer{
		orientation: orientation,
		paperSize:   paperSize,
		impl:        gofpdf.New(orientation, "pt", paperSize.Name, ""),
		templates:   pdfs.NewTemplateStore[int](),
	}
}

func (w *Writer) PaperSize() pdfs.PaperSize {
	return w.paperSize
}

func (w *Writer) Orientation() string {
	return w.orientation
}

func (w *Writer) TemplateStore() *pdfs.TemplateStore[int] {
	return w.templates
}

func (w *Writer) ImportPageAsTemplate(filepath string, pageNum int, storeKey string) error {
	// Check filepath exist
	tplID := gofpdi.ImportPage(w.impl, filepath, pageNum, "/MediaBox")
	w.templates.Store(storeKey, tplID)
	return nil
}

func (w *Writer) AddBlankPage() {
	w.impl.AddPage()
}

func (w *Writer) AddTemplatePage(storeKey string) bool {
	template, ok := w.templates.Get(storeKey)
	if !ok {
		return false
	}
	w.impl.AddPage()
	gofpdi.UseImportedTemplate(w.impl, template, 0, 0, w.paperSize.Width, w.paperSize.Height)
	return true
}

func (w *Writer) SetFont(family string, style string, size float64) {
	w.impl.SetFont(family, style, size)
}

func (w *Writer) Text(x float64, y float64, text string) {
	w.impl.Text(x, y, text)
}

func (w *Writer) WriteTo(writer io.Writer) (int64, error) {
	cw := rw.NewCountWriter(writer)
	err := w.impl.Output(cw)
	return cw.BytesWritten(), err
}

func (w *Writer) WriteToFile(filepath string) error {
	pdfBytes, err := w.ProduceBytes()
	if err != nil {
		return err
	}
	return os.WriteFile(filepath, pdfBytes, 0644)
}

func (w *Writer) ProduceBytes() ([]byte, error) {
	var buf bytes.Buffer
	err := w.impl.Output(&buf)
	return buf.Bytes(), err
}
