// Copyright (C) 2026 Murilo Gomes <profmugomes.com.br>
// SPDX-License-Identifier: MIT

// Site: https://www.profmugomes.com.br

package mgreport

// PaperSize defines the dimensions of a paper in points (1/72 inch).
type PaperSize struct {
	W float64 // width in points
	H float64 // height in points
}

// Orientation defines the page orientation.
type Orientation int

const (
	Portrait  Orientation = iota // vertical
	Landscape                    // horizontal
)

// Margins defines the page margins in points.
type Margins struct {
	Left   float64
	Top    float64
	Right  float64
	Bottom float64
}

// Position defines an (x, y) coordinate in points relative to the section.
type Position struct {
	X float64
	Y float64
}

// Size defines width and height in points.
type Size struct {
	W float64
	H float64
}

// SectionType identifies which band of the report a section represents.
type SectionType int

const (
	ReportHeader SectionType = iota
	PageHeader
	Detail
	ReportFooter
	PageFooter
)

// String returns the name of the section type.
func (st SectionType) String() string {
	switch st {
	case ReportHeader:
		return "ReportHeader"
	case PageHeader:
		return "PageHeader"
	case Detail:
		return "Detail"
	case ReportFooter:
		return "ReportFooter"
	case PageFooter:
		return "PageFooter"
	default:
		return "Unknown"
	}
}

// FormatFunc is a function that formats a raw value for display.
type FormatFunc func(interface{}) string

// Common paper sizes in points (1 inch = 72 points).
var (
	A4      = PaperSize{W: 595.0, H: 842.0}
	Letter  = PaperSize{W: 612.0, H: 792.0}
	Legal   = PaperSize{W: 612.0, H: 1008.0}
	A5      = PaperSize{W: 420.0, H: 595.0}
)

// Default margins in points.
var DefaultMargins = Margins{
	Left:   36.0, // ~12.7mm
	Top:    36.0,
	Right:  36.0,
	Bottom: 36.0,
}

// Text alignment constants.
type Alignment int

const (
	AlignLeft   Alignment = iota
	AlignCenter
	AlignRight
)

// FontWeight constants.
type FontWeight int

const (
	FontNormal FontWeight = iota
	FontBold
)

// Color represents an RGBA color.
type Color struct {
	R, G, B, A uint8
}

// Predefined colors.
var (
	ColorBlack = Color{0, 0, 0, 255}
	ColorWhite = Color{255, 255, 255, 255}
	ColorGray  = Color{128, 128, 128, 255}
	ColorRed   = Color{255, 0, 0, 255}
	ColorBlue  = Color{0, 0, 255, 255}
)

// EffectiveDimensions returns the usable area inside the margins.
func (m Margins) EffectiveDimensions(page PaperSize) Size {
	return Size{
		W: page.W - m.Left - m.Right,
		H: page.H - m.Top - m.Bottom,
	}
}
