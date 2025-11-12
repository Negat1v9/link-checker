package service

import (
	"fmt"
	"io"

	linkmodel "github.com/Negat1v9/link-checker/internal/linkChecker/model"
	"github.com/phpdave11/gofpdf"
)

// interface for give out from service to write file
type PdfFileWriter interface {
	Output(w io.Writer) error
}

// create pdf report add header and columns name
func (s *LinkCheckerService) createPDFReport() *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 14)
	pdf.MoveTo(10, 10)
	pdf.SetFontSize(16)
	pdf.SetFontStyle("B")
	pdf.Cell(10, 10, "Report of Link group number 7:")
	pdf.MoveTo(10, pdf.GetY()+8)
	pdf.SetFontSize(15)
	pdf.SetFontStyle("")
	pdf.Cell(40, 10, "Number")
	pdf.Cell(100, 10, "URL")
	pdf.Cell(50, 10, "Status")

	return pdf
}

// add in pdf file row with link state
func (s *LinkCheckerService) addReportRowLinkState(pdf *gofpdf.Fpdf, rowNubmer int, link string, status linkmodel.LinkStatus) {
	pdf.SetFontSize(13)
	pdf.MoveTo(10, pdf.GetY()+8)
	pdf.SetLineWidth(0.4)
	pdf.Line(5, pdf.GetY()+1, 205, pdf.GetY()+1)
	pdf.MoveTo(10, pdf.GetY())
	// number
	pdf.Cell(40, 10, fmt.Sprintf("%d", rowNubmer))
	pdf.Cell(100, 10, link)
	if status == linkmodel.LinkStatusNotAvailable {
		pdf.SetTextColor(232, 134, 117)
	} else {
		pdf.SetTextColor(108, 169, 99)
	}
	pdf.Cell(50, 10, string(status))
	pdf.SetTextColor(0, 0, 0)
}
