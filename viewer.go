// Copyright (C) 2026 Murilo Gomes <profmugomes.com.br>
// SPDX-License-Identifier: MIT

// Site: https://www.profmugomes.com.br

package mgreport

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	minZoom float32 = 0.1
	maxZoom float32 = 5.0
)

// ─── ReportViewer Widget ─────────────────────────────────────────

// ReportViewer is a Fyne widget that displays a rendered report with
// zoom, page navigation, and toolbar.
type ReportViewer struct {
	widget.BaseWidget

	report      *Report
	pages       []*PageContent
	fynePages   []*fynePage
	zoomFactor  float32
	currentPage int
	totalPages  int
	onPrint     func()
	onExport    func()

	scroll *container.Scroll
}

// NewReportViewer creates a new viewer for the given report.
func NewReportViewer(report *Report) *ReportViewer {
	rv := &ReportViewer{
		report:      report,
		zoomFactor:  1.0,
		currentPage: 0,
	}
	rv.ExtendBaseWidget(rv)
	rv.renderPages()
	return rv
}

// renderPages re-renders all pages from the report model.
func (rv *ReportViewer) renderPages() {
	rv.pages = RenderReport(rv.report, rv.zoomFactor)
	rv.totalPages = len(rv.pages)
	if rv.currentPage >= rv.totalPages {
		rv.currentPage = rv.totalPages - 1
	}
	if rv.currentPage < 0 {
		rv.currentPage = 0
	}

	rv.fynePages = nil
	pageW := rv.report.PageSize.W
	pageH := rv.report.PageSize.H
	for _, pc := range rv.pages {
		rv.fynePages = append(rv.fynePages, PageToCanvas(pc, pageW, pageH, rv.zoomFactor))
	}
}

// CreateRenderer implements the Fyne WidgetRenderer interface.
func (rv *ReportViewer) CreateRenderer() fyne.WidgetRenderer {
	rv.scroll = container.NewVScroll(rv.buildPageContent())
	return &reportViewerRenderer{viewer: rv}
}

// buildPageContent creates the visual content for the current page.
func (rv *ReportViewer) buildPageContent() fyne.CanvasObject {
	if rv.currentPage >= len(rv.fynePages) || rv.currentPage < 0 {
		return canvas.NewRectangle(color.NRGBA{R: 45, G: 45, B: 48, A: 255})
	}
	p := rv.fynePages[rv.currentPage]
	return container.NewCenter(p.content)
}

// refreshScroll rebuilds the scroll content and refreshes the widget.
func (rv *ReportViewer) refreshScroll() {
	if rv.scroll == nil {
		return
	}
	rv.scroll.Content = rv.buildPageContent()
	rv.scroll.Refresh()
}

// SetZoom sets the zoom factor and re-renders all pages.
func (rv *ReportViewer) SetZoom(factor float32) {
	if factor < minZoom {
		factor = minZoom
	}
	if factor > maxZoom {
		factor = maxZoom
	}
	rv.zoomFactor = factor
	rv.renderPages()
	rv.refreshScroll()
	rv.Refresh()
}

// GetZoom returns the current zoom factor.
func (rv *ReportViewer) GetZoom() float32 {
	return rv.zoomFactor
}

// ZoomIn increases zoom by 10%.
func (rv *ReportViewer) ZoomIn() {
	rv.SetZoom(rv.zoomFactor + 0.1)
}

// ZoomOut decreases zoom by 10%.
func (rv *ReportViewer) ZoomOut() {
	rv.SetZoom(rv.zoomFactor - 0.1)
}

// ZoomToFit sets zoom so the page width fits the container width.
func (rv *ReportViewer) ZoomToFit(containerWidth float32) {
	if containerWidth <= 0 {
		return
	}
	factor := containerWidth / float32(rv.report.PageSize.W)
	rv.SetZoom(factor)
}

// NextPage advances to the next page.
func (rv *ReportViewer) NextPage() {
	if rv.currentPage < rv.totalPages-1 {
		rv.currentPage++
		rv.refreshScroll()
		rv.Refresh()
	}
}

// PrevPage goes to the previous page.
func (rv *ReportViewer) PrevPage() {
	if rv.currentPage > 0 {
		rv.currentPage--
		rv.refreshScroll()
		rv.Refresh()
	}
}

// FirstPage goes to the first page.
func (rv *ReportViewer) FirstPage() {
	rv.currentPage = 0
	rv.refreshScroll()
	rv.Refresh()
}

// LastPage goes to the last page.
func (rv *ReportViewer) LastPage() {
	rv.currentPage = rv.totalPages - 1
	if rv.currentPage < 0 {
		rv.currentPage = 0
	}
	rv.refreshScroll()
	rv.Refresh()
}

// CurrentPageNumber returns the 1-based page number.
func (rv *ReportViewer) CurrentPageNumber() int {
	return rv.currentPage + 1
}

// TotalPagesCount returns the total number of pages.
func (rv *ReportViewer) TotalPagesCount() int {
	return rv.totalPages
}

// SetOnPrint sets the callback for the print action.
func (rv *ReportViewer) SetOnPrint(fn func()) {
	rv.onPrint = fn
}

// SetOnExport sets the callback for the export action.
func (rv *ReportViewer) SetOnExport(fn func()) {
	rv.onExport = fn
}

// ─── Widget Renderer ─────────────────────────────────────────────

type reportViewerRenderer struct {
	viewer *ReportViewer
}

func (r *reportViewerRenderer) Destroy() {}

func (r *reportViewerRenderer) Layout(size fyne.Size) {
	if r.viewer.scroll != nil {
		r.viewer.scroll.Resize(size)
	}
}

func (r *reportViewerRenderer) MinSize() fyne.Size {
	return fyne.NewSize(400, 300)
}

func (r *reportViewerRenderer) Objects() []fyne.CanvasObject {
	if r.viewer.scroll != nil {
		return []fyne.CanvasObject{r.viewer.scroll}
	}
	return nil
}

func (r *reportViewerRenderer) Refresh() {
	if r.viewer.scroll != nil {
		r.viewer.scroll.Refresh()
	}
}

// ─── Print Setup Dialog ──────────────────────────────────────────

// PrintSetupDialog opens a dialog to configure paper size and orientation.
func PrintSetupDialog(report *Report, parent fyne.Window, onApply func()) {
	paperSizes := []string{"A4 (210×297mm)", "Letter (216×279mm)", "Legal (216×356mm)", "A5 (148×210mm)"}
	paperKeys := []PaperSize{A4, Letter, Legal, A5}

	orientations := []string{"Retrato (Portrait)", "Paisagem (Landscape)"}

	// Find current selection
	currentPaper := 0
	for i, ps := range paperKeys {
		if ps.W == report.PageSize.W && ps.H == report.PageSize.H {
			currentPaper = i
			break
		}
	}
	currentOrient := 0
	if report.Orientation == Landscape {
		currentOrient = 1
	}

	paperSelect := widget.NewSelect(paperTitles(paperSizes), nil)
	paperSelect.SetSelected(paperSizes[currentPaper])

	orientSelect := widget.NewSelect(orientations, nil)
	orientSelect.SetSelected(orientations[currentOrient])

	content := container.NewVBox(
		widget.NewLabel("Configurações de Impressão"),
		widget.NewSeparator(),
		widget.NewLabel("Tamanho do Papel:"),
		paperSelect,
		widget.NewLabel("Orientação:"),
		orientSelect,
	)

	dlg := dialog.NewCustom("Configurações de Impressão", "Cancelar", content, parent)
	dlg.Resize(fyne.NewSize(350, 250))

	// Add Apply button by replacing content with buttons
	btnApply := widget.NewButton("Aplicar", func() {
		// Apply paper size
		for i, title := range paperSizes {
			if paperSelect.Selected == title {
				report.PageSize = paperKeys[i]
				break
			}
		}
		// Apply orientation
		if orientSelect.Selected == orientations[1] {
			report.Orientation = Landscape
		} else {
			report.Orientation = Portrait
		}

		dlg.Hide()
		if onApply != nil {
			onApply()
		}
	})

	fullContent := container.NewVBox(
		widget.NewLabel("Configurações de Impressão"),
		widget.NewSeparator(),
		widget.NewLabel("Tamanho do Papel:"),
		paperSelect,
		widget.NewLabel("Orientação:"),
		orientSelect,
		widget.NewSeparator(),
		container.NewHBox(layout.NewSpacer(), btnApply),
	)

	dlg = dialog.NewCustom("Configurações de Impressão", "Fechar", fullContent, parent)
	dlg.Resize(fyne.NewSize(350, 250))
	dlg.Show()
}

func paperTitles(titles []string) []string {
	return titles
}

// ─── Preview Window ──────────────────────────────────────────────

// ShowPreview opens a standalone preview window with the full toolbar.
func ShowPreview(app fyne.App, report *Report) {
	win := app.NewWindow("MGReport Preview")
	// Window size to approximate A4 proportions at ~96 DPI
	// A4 at 96DPI ≈ 794×1123px, but we use a reasonable screen size
	win.Resize(fyne.NewSize(900, 1050))
	win.CenterOnScreen()

	viewer := NewReportViewer(report)

	zoomLabel := widget.NewLabel("100%")
	pageLabel := widget.NewLabel("Página 1 de 1")

	updatePageLabel := func() {
		pageLabel.SetText(fmt.Sprintf("Página %d de %d", viewer.CurrentPageNumber(), viewer.TotalPagesCount()))
	}

	// Export PDF callback
	viewer.SetOnExport(func() {
		ExportDialog(report, win)
	})

	// Print callback
	viewer.SetOnPrint(func() {
		PrintReport(report)
	})

	// Navigation buttons
	btnFirst := widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), func() {
		viewer.FirstPage()
		updatePageLabel()
	})
	btnPrev := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		viewer.PrevPage()
		updatePageLabel()
	})
	btnNext := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		viewer.NextPage()
		updatePageLabel()
	})
	btnLast := widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), func() {
		viewer.LastPage()
		updatePageLabel()
	})

	// Zoom buttons
	btnZoomOut := widget.NewButtonWithIcon("", theme.ZoomOutIcon(), func() {
		viewer.ZoomOut()
		zoomLabel.SetText(fmt.Sprintf("%.0f%%", viewer.GetZoom()*100))
	})
	btnZoomIn := widget.NewButtonWithIcon("", theme.ZoomInIcon(), func() {
		viewer.ZoomIn()
		zoomLabel.SetText(fmt.Sprintf("%.0f%%", viewer.GetZoom()*100))
	})
	btnFitWidth := widget.NewButtonWithIcon("Ajustar", theme.ViewFullScreenIcon(), func() {
		w := win.Canvas().Size().Width
		viewer.ZoomToFit(w - 20)
		zoomLabel.SetText(fmt.Sprintf("%.0f%%", viewer.GetZoom()*100))
	})

	// Export, Print, and Print Setup buttons
	btnExport := widget.NewButtonWithIcon("PDF", theme.DocumentSaveIcon(), func() {
		ExportDialog(report, win)
	})
	btnPrintSetup := widget.NewButtonWithIcon("Configurar", theme.SettingsIcon(), func() {
		PrintSetupDialog(report, win, func() {
			// Re-render and re-fit after settings change
			viewer.renderPages()
			w := win.Canvas().Size().Width
			viewer.ZoomToFit(w - 20)
			zoomLabel.SetText(fmt.Sprintf("%.0f%%", viewer.GetZoom()*100))
			updatePageLabel()
		})
	})
	btnPrint := widget.NewButtonWithIcon("Imprimir", theme.ViewRefreshIcon(), func() {
		PrintReport(report)
	})

	// Toolbar row 1: title + navigation
	toolbar1 := container.NewHBox(
		widget.NewLabelWithStyle(report.Title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewLabel("Página"),
		btnFirst,
		btnPrev,
		pageLabel,
		btnNext,
		btnLast,
	)

	// Toolbar row 2: zoom + actions
	toolbar2 := container.NewHBox(
		widget.NewLabel("Zoom:"),
		btnZoomOut,
		zoomLabel,
		btnZoomIn,
		btnFitWidth,
		layout.NewSpacer(),
		btnPrintSetup,
		btnExport,
		btnPrint,
	)

	toolbar := container.NewVBox(toolbar1, toolbar2)

	// Main layout: toolbar on top, viewer fills remaining space
	win.SetContent(container.NewBorder(
		container.NewPadded(toolbar), // top
		nil,                          // bottom
		nil,                          // left
		nil,                          // right
		viewer,                       // center
	))

	// Keyboard shortcuts
	win.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		switch ke.Name {
		case fyne.KeyRight:
			viewer.NextPage()
			updatePageLabel()
		case fyne.KeyLeft:
			viewer.PrevPage()
			updatePageLabel()
		case fyne.KeyMinus:
			viewer.ZoomOut()
			zoomLabel.SetText(fmt.Sprintf("%.0f%%", viewer.GetZoom()*100))
		case fyne.KeyEqual:
			viewer.ZoomIn()
			zoomLabel.SetText(fmt.Sprintf("%.0f%%", viewer.GetZoom()*100))
		case fyne.KeyHome:
			viewer.FirstPage()
			updatePageLabel()
		case fyne.KeyEnd:
			viewer.LastPage()
			updatePageLabel()
		}
	})

	// Auto-fit A4 on first show: pre-compute zoom from the known window width
	initialZoom := (900 - 20) / float32(report.PageSize.W)
	viewer.SetZoom(initialZoom)
	zoomLabel.SetText(fmt.Sprintf("%.0f%%", viewer.GetZoom()*100))
	updatePageLabel()

	win.Show()
}
