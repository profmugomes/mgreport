// Required Notice: Copyright (c) 2026 Murilo Gomes Julio. All Rights Reserved. (https://profmugomes.com.br)

// Licensed under the PolyForm Perimeter License 1.0.1.
// See LICENSE.md for details.

package mgreport

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// ─── Render Context ──────────────────────────────────────────────

// RenderContext holds the state during rendering of a single section.
type RenderContext struct {
	Zoom      float32
	PageIndex int
	DataRow   map[string]interface{}

	// Current section origin (absolute coordinates on the page in points)
	SectionX float64
	SectionY float64

	// Collected elements for the current page
	texts  []*TextElement
	lines  []*LineElement
}

// AddText appends a text element to the current context.
func (rc *RenderContext) AddText(t *TextElement) {
	rc.texts = append(rc.texts, t)
}

// AddLine appends a line element to the current context.
func (rc *RenderContext) AddLine(l *LineElement) {
	rc.lines = append(rc.lines, l)
}

// ─── Element Types ───────────────────────────────────────────────

// TextElement is a resolved text ready for rendering to canvas.
type TextElement struct {
	Text     string
	X        float64 // absolute X in points
	Y        float64 // absolute Y in points
	FontSize float64
	FontFace string
	Color    Color
	Align    Alignment
	Bold     bool
}

// LineElement is a resolved line ready for rendering to canvas.
type LineElement struct {
	X1, Y1 float64
	X2, Y2 float64
	Color  Color
	Width  float64
}

// ─── Page Content ────────────────────────────────────────────────

// PageContent holds all resolved elements for a single page.
type PageContent struct {
	Texts []*TextElement
	Lines []*LineElement
}

// ─── Report Rendering ────────────────────────────────────────────

// RenderReport converts a Report definition and its DataSource into a list of PageContent.
// Each PageContent contains the resolved elements for one page.
func RenderReport(report *Report, zoom float32) []*PageContent {
	r := &renderer{
		report: report,
		zoom:   zoom,
	}
	return r.render()
}

type renderer struct {
	report *Report
	zoom   float32
	pages  []*PageContent
	current *PageContent
}

func (r *renderer) newPage() *PageContent {
	pc := &PageContent{}
	r.pages = append(r.pages, pc)
	r.current = pc
	return pc
}

func (r *renderer) render() []*PageContent {
	r.pages = nil
	r.newPage()

	pageIndex := 0
	pageY := r.report.Margins.Top // current Y cursor in points on the page

	// ── ReportHeader (rendered once on first page) ──
	if r.report.HasSection(ReportHeader) {
		sec := r.report.GetSection(ReportHeader)
		r.renderSection(sec, 0, pageY)
		pageY += sec.Height
	}

	// ── PageHeader (rendered at top of every page) ──
	pageHeaderHeight := 0.0
	if r.report.HasSection(PageHeader) {
		pageHeaderHeight = r.report.GetSection(PageHeader).Height
	}

	// Render page header on first page
	if r.report.HasSection(PageHeader) {
		r.renderSection(r.report.GetSection(PageHeader), 0, pageY)
		pageY += pageHeaderHeight
	}

	// ── Detail rows ──
	detailHeight := 0.0
	if r.report.HasSection(Detail) {
		detailHeight = r.report.GetSection(Detail).Height
	}

	dataSource := r.report.DataSource
	if dataSource == nil {
		dataSource = []map[string]interface{}{}
	}

	usableHeight := r.report.EffectiveHeight()

	for i, row := range dataSource {
		// Check if the detail row fits on the current page
		if pageY+detailHeight > r.report.Margins.Top+usableHeight {
			// ── PageFooter (before page break) ──
			if r.report.HasSection(PageFooter) {
				footerSec := r.report.GetSection(PageFooter)
				footerY := r.report.PageSize.H - r.report.Margins.Bottom - footerSec.Height
				r.renderSectionAt(footerSec, 0, footerY, i)
			}

			// Start new page
			pageIndex++
			r.newPage()
			pageY = r.report.Margins.Top

			// Re-render PageHeader on new page
			if r.report.HasSection(PageHeader) {
				r.renderSection(r.report.GetSection(PageHeader), 0, pageY)
				pageY += pageHeaderHeight
			}
		}

		// Render detail row
		if r.report.HasSection(Detail) {
			r.renderSectionWithData(r.report.GetSection(Detail), 0, pageY, row, pageIndex)
			pageY += detailHeight
		}
	}

	// ── ReportFooter (rendered once on last page) ──
	if r.report.HasSection(ReportFooter) {
		sec := r.report.GetSection(ReportFooter)
		footerY := r.report.PageSize.H - r.report.Margins.Bottom - sec.Height
		// If it does not fit, start a new page
		if footerY < pageY {
			pageIndex++
			r.newPage()
			footerY = r.report.Margins.Top
		}
		lastRow := map[string]interface{}{}
		if len(dataSource) > 0 {
			lastRow = dataSource[len(dataSource)-1]
		}
		r.renderSectionWithData(sec, 0, footerY, lastRow, pageIndex)
	}

	// ── PageFooter on last page if not already rendered ──
	if !r.report.HasSection(ReportFooter) && r.report.HasSection(PageFooter) {
		footerSec := r.report.GetSection(PageFooter)
		footerY := r.report.PageSize.H - r.report.Margins.Bottom - footerSec.Height
		lastRow := map[string]interface{}{}
		if len(dataSource) > 0 {
			lastRow = dataSource[len(dataSource)-1]
		}
		r.renderSectionWithData(footerSec, 0, footerY, lastRow, pageIndex)
	}

	return r.pages
}

// renderSection renders a section without data context.
func (r *renderer) renderSection(sec *Section, x, y float64) {
	ctx := &RenderContext{
		Zoom:      r.zoom,
		PageIndex: len(r.pages) - 1,
		SectionX:  x,
		SectionY:  y,
		DataRow:   nil,
	}
	for _, item := range sec.Items {
		item.render(ctx)
	}
	r.current.Texts = append(r.current.Texts, ctx.texts...)
	r.current.Lines = append(r.current.Lines, ctx.lines...)
}

// renderSectionWithData renders a section with a data row context.
func (r *renderer) renderSectionWithData(sec *Section, x, y float64, row map[string]interface{}, pageIndex int) {
	ctx := &RenderContext{
		Zoom:      r.zoom,
		PageIndex: pageIndex,
		SectionX:  x,
		SectionY:  y,
		DataRow:   row,
	}
	for _, item := range sec.Items {
		item.render(ctx)
	}
	r.current.Texts = append(r.current.Texts, ctx.texts...)
	r.current.Lines = append(r.current.Lines, ctx.lines...)
}

// renderSectionAt renders a section at an absolute position with a specific page index.
func (r *renderer) renderSectionAt(sec *Section, x, y float64, pageIndex int) {
	ctx := &RenderContext{
		Zoom:      r.zoom,
		PageIndex: pageIndex,
		SectionX:  x,
		SectionY:  y,
		DataRow:   nil,
	}
	for _, item := range sec.Items {
		item.render(ctx)
	}
	r.current.Texts = append(r.current.Texts, ctx.texts...)
	r.current.Lines = append(r.current.Lines, ctx.lines...)
}

// ─── Canvas Conversion ───────────────────────────────────────────

// PageToCanvas converts a PageContent into a Fyne CanvasObject (a page visual).
func PageToCanvas(pc *PageContent, pageW, pageH float64, zoom float32) *fynePage {
	scaledW := float32(pageW) * zoom
	scaledH := float32(pageH) * zoom

	objects := []fyne.CanvasObject{}

	// White paper background — fills entire page
	paper := canvas.NewRectangle(color.White)
	paper.SetMinSize(fyneSize(scaledW, scaledH))
	paper.Resize(fyneSize(scaledW, scaledH))
	objects = append(objects, paper)

	// Lines
	for _, l := range pc.Lines {
		line := canvas.NewLine(color.NRGBA{
			R: l.Color.R, G: l.Color.G, B: l.Color.B, A: l.Color.A,
		})
		line.StrokeWidth = float32(l.Width)
		line.Position1 = fynePos(float32(l.X1)*zoom, float32(l.Y1)*zoom)
		line.Position2 = fynePos(float32(l.X2)*zoom, float32(l.Y2)*zoom)
		objects = append(objects, line)
	}

	// Texts
	for _, t := range pc.Texts {
		txt := canvas.NewText(t.Text, color.NRGBA{
			R: t.Color.R, G: t.Color.G, B: t.Color.B, A: t.Color.A,
		})
		txt.TextSize = float32(t.FontSize) * zoom
		txt.Text = t.Text

		// Estimate size based on character count
		estimatedW := float32(float64(len(t.Text)) * t.FontSize * 0.6)
		estimatedH := float32(t.FontSize * 1.4)
		txt.Resize(fyneSize(estimatedW*zoom, estimatedH*zoom))
		txt.Move(fynePos(float32(t.X)*zoom, float32(t.Y)*zoom))

		switch t.Align {
		case AlignCenter:
			txt.Alignment = 1
		case AlignRight:
			txt.Alignment = 2
		default:
			txt.Alignment = 0
		}

		objects = append(objects, txt)
	}

	// Use pageWidget so children keep their explicit Move() positions
	// (container.NewStack would Resize all children to the page size, breaking text positions)
	return &fynePage{
		content: newPageWidget(scaledW, scaledH, objects...),
		width:   scaledW,
		height:  scaledH,
	}
}

// ─── Page Widget (absolute positioning) ──────────────────────────

// pageWidget is a Fyne widget that renders objects at their explicit positions
// without overriding them via layout. The background rectangle fills the widget.
type pageWidget struct {
	widget.BaseWidget
	minW, minH float32
	objects    []fyne.CanvasObject
}

func newPageWidget(w, h float32, objects ...fyne.CanvasObject) *pageWidget {
	pw := &pageWidget{minW: w, minH: h, objects: objects}
	pw.ExtendBaseWidget(pw)
	return pw
}

func (pw *pageWidget) CreateRenderer() fyne.WidgetRenderer {
	return &pageWidgetRenderer{widget: pw}
}

type pageWidgetRenderer struct {
	widget *pageWidget
}

func (r *pageWidgetRenderer) Destroy() {}

func (r *pageWidgetRenderer) Layout(size fyne.Size) {
	// Only resize the background rectangle to fill the widget.
	// All other elements (text, lines) keep their Move()-set positions.
	for _, obj := range r.widget.objects {
		if _, ok := obj.(*canvas.Rectangle); ok {
			obj.Resize(size)
		}
	}
}

func (r *pageWidgetRenderer) MinSize() fyne.Size {
	return fyne.NewSize(r.widget.minW, r.widget.minH)
}

func (r *pageWidgetRenderer) Objects() []fyne.CanvasObject {
	return r.widget.objects
}

func (r *pageWidgetRenderer) Refresh() {
	for _, obj := range r.widget.objects {
		obj.Refresh()
	}
}

// ─── Helpers ─────────────────────────────────────────────────────

func fynePos(x, y float32) fyne.Position {
	return fyne.NewPos(x, y)
}

func fyneSize(w, h float32) fyne.Size {
	return fyne.NewSize(w, h)
}

// fynePage holds a rendered page as a Fyne container.
type fynePage struct {
	content fyne.CanvasObject
	width   float32
	height  float32
}
