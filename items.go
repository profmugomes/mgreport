// Required Notice: Copyright (c) 2026 Murilo Gomes Julio. All Rights Reserved. (https://profmugomes.com.br)

// Licensed under the PolyForm Perimeter License 1.0.1.
// See LICENSE.md for details.

package mgreport

import "fmt"

// ─── ReportItem interface ────────────────────────────────────────

// ReportItem is anything that can be rendered inside a section.
type ReportItem interface {
	render(ctx *RenderContext) // internal rendering
}

// ─── Label ───────────────────────────────────────────────────────

// Label is a static text item.
type Label struct {
	Text     string
	Pos      Position
	FontSize float64
	FontFace string
	Color    Color
	Align    Alignment
	Bold     bool
}

// render draws the label on the render context.
func (l *Label) render(ctx *RenderContext) {
	absX := ctx.SectionX + l.Pos.X
	absY := ctx.SectionY + l.Pos.Y
	fontSize := l.FontSize * float64(ctx.Zoom)

	ctx.AddText(&TextElement{
		Text:     l.Text,
		X:        absX,
		Y:        absY,
		FontSize: fontSize,
		FontFace: l.FontFace,
		Color:    l.Color,
		Align:    l.Align,
		Bold:     l.Bold,
	})
}

// LabelOption is a functional option for Label configuration.
type LabelOption func(*Label)

// WithFontSize sets the font size in points.
func WithFontSize(size float64) LabelOption {
	return func(l *Label) { l.FontSize = size }
}

// WithFontFace sets the font face name.
func WithFontFace(face string) LabelOption {
	return func(l *Label) { l.FontFace = face }
}

// WithColor sets the text color.
func WithColor(c Color) LabelOption {
	return func(l *Label) { l.Color = c }
}

// WithAlign sets the text alignment.
func WithAlign(a Alignment) LabelOption {
	return func(l *Label) { l.Align = a }
}

// WithBold enables bold text.
func WithBold(b bool) LabelOption {
	return func(l *Label) { l.Bold = b }
}

// ─── Field ───────────────────────────────────────────────────────

// Field is a data-bound text item that displays a value from the data source.
type Field struct {
	FieldName string
	Pos       Position
	FontSize  float64
	FontFace  string
	Color     Color
	Align     Alignment
	Bold      bool
	Format    FormatFunc // optional formatter
}

// render draws the field value on the render context.
func (f *Field) render(ctx *RenderContext) {
	if ctx.DataRow == nil {
		return
	}
	raw, ok := ctx.DataRow[f.FieldName]
	if !ok {
		return
	}

	var display string
	if f.Format != nil {
		display = f.Format(raw)
	} else {
		display = fmt.Sprintf("%v", raw)
	}

	absX := ctx.SectionX + f.Pos.X
	absY := ctx.SectionY + f.Pos.Y
	fontSize := f.FontSize * float64(ctx.Zoom)

	ctx.AddText(&TextElement{
		Text:     display,
		X:        absX,
		Y:        absY,
		FontSize: fontSize,
		FontFace: f.FontFace,
		Color:    f.Color,
		Align:    f.Align,
		Bold:     f.Bold,
	})
}

// FieldOption is a functional option for Field configuration.
type FieldOption func(*Field)

// WithFieldFontSize sets the font size for a field.
func WithFieldFontSize(size float64) FieldOption {
	return func(f *Field) { f.FontSize = size }
}

// WithFieldFormat sets the format function for a field.
func WithFieldFormat(fn FormatFunc) FieldOption {
	return func(f *Field) { f.Format = fn }
}

// WithFieldColor sets the text color for a field.
func WithFieldColor(c Color) FieldOption {
	return func(f *Field) { f.Color = c }
}

// WithFieldAlign sets the alignment for a field.
func WithFieldAlign(a Alignment) FieldOption {
	return func(f *Field) { f.Align = a }
}

// WithFieldBold enables bold text for a field.
func WithFieldBold(b bool) FieldOption {
	return func(f *Field) { f.Bold = b }
}

// WithFieldFontFace sets the font face for a field.
func WithFieldFontFace(face string) FieldOption {
	return func(f *Field) { f.FontFace = face }
}

// ─── Line ────────────────────────────────────────────────────────

// Line is a straight line between two points.
type Line struct {
	X1, Y1 float64
	X2, Y2 float64
	Color  Color
	Width  float64
}

// render draws the line on the render context.
func (l *Line) render(ctx *RenderContext) {
	absX1 := ctx.SectionX + l.X1
	absY1 := ctx.SectionY + l.Y1
	absX2 := ctx.SectionX + l.X2
	absY2 := ctx.SectionY + l.Y2
	width := l.Width * float64(ctx.Zoom)

	ctx.AddLine(&LineElement{
		X1: absX1, Y1: absY1,
		X2: absX2, Y2: absY2,
		Color: l.Color,
		Width: width,
	})
}

// LineOption is a functional option for Line configuration.
type LineOption func(*Line)

// WithLineColor sets the line color.
func WithLineColor(c Color) LineOption {
	return func(l *Line) { l.Color = c }
}

// WithLineWidth sets the line width in points.
func WithLineWidth(w float64) LineOption {
	return func(l *Line) { l.Width = w }
}

// ─── PageNumber ──────────────────────────────────────────────────

// PageNumber is an automatic page number placeholder.
type PageNumber struct {
	Label
}

// render draws the page number on the render context.
func (pn *PageNumber) render(ctx *RenderContext) {
	text := fmt.Sprintf("%d", ctx.PageIndex+1)
	absX := ctx.SectionX + pn.Pos.X
	absY := ctx.SectionY + pn.Pos.Y
	fontSize := pn.FontSize * float64(ctx.Zoom)

	ctx.AddText(&TextElement{
		Text:     text,
		X:        absX,
		Y:        absY,
		FontSize: fontSize,
		FontFace: pn.FontFace,
		Color:    pn.Color,
		Align:    pn.Align,
		Bold:     pn.Bold,
	})
}

// ─── Predefined formatters ───────────────────────────────────────

// FormatCurrency formats a number as Brazilian Real currency.
func FormatCurrency(v interface{}) string {
	switch val := v.(type) {
	case float64:
		return fmt.Sprintf("R$ %.2f", val)
	case float32:
		return fmt.Sprintf("R$ %.2f", val)
	case int:
		return fmt.Sprintf("R$ %d.00", val)
	case int64:
		return fmt.Sprintf("R$ %d.00", val)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// FormatInteger formats a number as integer with thousand separators.
func FormatInteger(v interface{}) string {
	switch val := v.(type) {
	case float64:
		return fmt.Sprintf("%.0f", val)
	case int:
		return fmt.Sprintf("%d", val)
	case int64:
		return fmt.Sprintf("%d", val)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// FormatPercent formats a number as percentage.
func FormatPercent(v interface{}) string {
	switch val := v.(type) {
	case float64:
		return fmt.Sprintf("%.2f%%", val*100)
	case float32:
		return fmt.Sprintf("%.2f%%", float64(val)*100)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// FormatDateBR formats a string as Brazilian date (DD/MM/YYYY).
func FormatDateBR(v interface{}) string {
	s := fmt.Sprintf("%v", v)
	if len(s) >= 10 {
		// Assume YYYY-MM-DD or YYYY/MM/DD
		return s[8:10] + "/" + s[5:7] + "/" + s[0:4]
	}
	return s
}
