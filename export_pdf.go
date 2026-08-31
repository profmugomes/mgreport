// Required Notice: Copyright (c) 2026 Murilo Gomes Julio. All Rights Reserved. (https://profmugomes.com.br)

// Licensed under the PolyForm Perimeter License 1.0.1.
// See LICENSE.md for details.

package mgreport

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"

	"codeberg.org/go-pdf/fpdf"
)

// ExportToPDF renders the report to a PDF file at the given path.
func ExportToPDF(report *Report, savePath string) error {
	// Determine PDF page orientation
	orient := "P"
	if report.Orientation == Landscape {
		orient = "L"
	}

	// Use fpdf predefined page sizes
	paperName := pdfPaperName(report.PageSize)

	pdf := fpdf.New(orient, "mm", paperName, "")
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	// Conversion: report points → mm
	pointToMm := 25.4 / 72.0

	// Remove all margins — coordinates from the renderer already include margins
	pdf.SetMargins(0, 0, 0)
	pdf.SetAutoPageBreak(false, 0)

	for _, page := range renderedPages(report) {
		pdf.AddPage()

		// ── Draw lines ──
		for _, l := range page.Lines {
			r, g, b := l.Color.R, l.Color.G, l.Color.B
			pdf.SetDrawColor(int(r), int(g), int(b))
			pdf.SetLineWidth(l.Width * pointToMm)
			pdf.Line(
				l.X1*pointToMm,
				l.Y1*pointToMm,
				l.X2*pointToMm,
				l.Y2*pointToMm,
			)
		}
		pdf.SetDrawColor(0, 0, 0)

		// ── Draw texts ──
		for _, t := range page.Texts {
			r, g, b := t.Color.R, t.Color.G, t.Color.B
			pdf.SetTextColor(int(r), int(g), int(b))

			style := ""
			if t.Bold {
				style = "B"
			}

			// Convert font size from points to mm
			fSize := t.FontSize * pointToMm
			if fSize < 1 {
				fSize = 1
			}

			pdf.SetFont("Arial", style, fSize)

			pdfX := t.X * pointToMm
			pdfY := t.Y * pointToMm

			// Line height: approx 1.2× font size
			lineH := fSize * 1.2

			// Alignment
			align := "L"
			switch t.Align {
			case AlignCenter:
				align = "C"
			case AlignRight:
				align = "R"
			}

			// Use CellFormat for single-line text (no wrapping)
			// Calculate a reasonable width for the text cell
			cellW := float64(len(t.Text)) * fSize * 0.55
			if cellW < 10 {
				cellW = 10
			}

			pdf.SetXY(pdfX, pdfY)
			pdf.CellFormat(cellW, lineH, tr(t.Text), "", 0, align, false, 0, "")
		}
		pdf.SetTextColor(0, 0, 0)
	}

	return pdf.OutputFileAndClose(savePath)
}

// pdfPaperName returns the fpdf predefined paper name for the given PaperSize.
func pdfPaperName(ps PaperSize) string {
	// Match by width and height (in points)
	switch {
	case ps.W == A4.W && ps.H == A4.H:
		return "A4"
	case ps.W == Letter.W && ps.H == Letter.H:
		return "Letter"
	case ps.W == Legal.W && ps.H == Legal.H:
		return "Legal"
	case ps.W == A5.W && ps.H == A5.H:
		return "A5"
	default:
		return "A4" // fallback
	}
}

// renderedPages renders the report at zoom 1.0 for PDF export.
func renderedPages(report *Report) []*PageContent {
	return RenderReport(report, 1.0)
}

// ExportDialog opens a file save dialog for PDF export.
func ExportDialog(report *Report, parent fyne.Window) {
	saveDlg := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, parent)
			return
		}
		if writer == nil {
			return
		}

		path := writer.URI().Path()
		writer.Close()

		err = ExportToPDF(report, path)
		if err != nil {
			dialog.ShowError(err, parent)
		} else {
			dialog.ShowInformation("Sucesso", "PDF exportado com sucesso!", parent)
		}
	}, parent)
	saveDlg.SetFileName("relatorio.pdf")
	saveDlg.Show()
}
