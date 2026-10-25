// Copyright (C) 2026 Murilo Gomes <profmugomes.com.br>
// SPDX-License-Identifier: MIT

// Site: https://www.profmugomes.com.br

package mgreport

// Report is the top-level definition of a report.
type Report struct {
	Title       string
	PageSize    PaperSize
	Orientation Orientation
	Margins     Margins
	DataSource  []map[string]interface{}
	sections    map[SectionType]*Section
}

// Section represents a band within the report (e.g., ReportHeader, Detail, PageFooter).
type Section struct {
	Height  float64      // height in points
	Items   []ReportItem // items to render inside this section
	visible bool
}

// NewReport creates a new Report with sensible defaults (A4 Portrait).
func NewReport() *Report {
	return &Report{
		PageSize:    A4,
		Orientation: Portrait,
		Margins:     DefaultMargins,
		sections:    make(map[SectionType]*Section),
	}
}

// SetTitle sets the report title.
func (r *Report) SetTitle(title string) *Report {
	r.Title = title
	return r
}

// SetPageSize sets the paper size and orientation.
func (r *Report) SetPageSize(size PaperSize, orient Orientation) *Report {
	r.PageSize = size
	r.Orientation = orient
	return r
}

// SetMargins sets the page margins in points.
func (r *Report) SetMargins(left, top, right, bottom float64) *Report {
	r.Margins = Margins{Left: left, Top: top, Right: right, Bottom: bottom}
	return r
}

// SetDataSource sets the data to be rendered in the Detail section.
func (r *Report) SetDataSource(data []map[string]interface{}) *Report {
	r.DataSource = data
	return r
}

// Section returns a SectionBuilder for the given section type.
// If the section does not exist yet, it is created with the specified height.
func (r *Report) Section(typ SectionType, height float64) *SectionBuilder {
	sec, ok := r.sections[typ]
	if !ok {
		sec = &Section{Height: height, visible: true}
		r.sections[typ] = sec
	} else {
		sec.Height = height
	}
	return &SectionBuilder{report: r, section: sec, sectionType: typ}
}

// GetSection returns the section for the given type, or nil if not defined.
func (r *Report) GetSection(typ SectionType) *Section {
	return r.sections[typ]
}

// HasSection returns true if the section type is defined.
func (r *Report) HasSection(typ SectionType) bool {
	sec, ok := r.sections[typ]
	return ok && sec != nil && sec.visible
}

// VisibleSections returns all section types in display order that are defined and visible.
func (r *Report) VisibleSections() []SectionType {
	order := []SectionType{ReportHeader, PageHeader, Detail, ReportFooter, PageFooter}
	var result []SectionType
	for _, st := range order {
		if r.HasSection(st) {
			result = append(result, st)
		}
	}
	return result
}

// EffectiveWidth returns the usable width inside margins.
func (r *Report) EffectiveWidth() float64 {
	return r.PageSize.W - r.Margins.Left - r.Margins.Right
}

// EffectiveHeight returns the usable height inside margins.
func (r *Report) EffectiveHeight() float64 {
	return r.PageSize.H - r.Margins.Top - r.Margins.Bottom
}

// ─── SectionBuilder ──────────────────────────────────────────────

// SectionBuilder provides a fluent API for adding items to a section.
type SectionBuilder struct {
	report      *Report
	section     *Section
	sectionType SectionType
}

// Section returns the builder (for chaining).
func (sb *SectionBuilder) Section() *Section {
	return sb.section
}

// AddLabel adds a static text label to the section.
func (sb *SectionBuilder) AddLabel(text string, pos Position, opts ...LabelOption) *SectionBuilder {
	label := &Label{
		Text:     text,
		Pos:      pos,
		FontSize: 14.0,
		FontFace: "Arial",
		Color:    ColorBlack,
		Align:    AlignLeft,
	}
	for _, opt := range opts {
		opt(label)
	}
	sb.section.Items = append(sb.section.Items, label)
	return sb
}

// AddField adds a data-bound field to the section.
func (sb *SectionBuilder) AddField(fieldName string, pos Position, opts ...FieldOption) *SectionBuilder {
	field := &Field{
		FieldName: fieldName,
		Pos:       pos,
		FontSize:  14.0,
		FontFace:  "Arial",
		Color:     ColorBlack,
		Align:     AlignLeft,
	}
	for _, opt := range opts {
		opt(field)
	}
	sb.section.Items = append(sb.section.Items, field)
	return sb
}

// AddLine adds a line (horizontal or vertical) to the section.
func (sb *SectionBuilder) AddLine(x1, y1, x2, y2 float64, opts ...LineOption) *SectionBuilder {
	line := &Line{
		X1:      x1,
		Y1:      y1,
		X2:      x2,
		Y2:      y2,
		Color:   ColorBlack,
		Width:   0.5,
	}
	for _, opt := range opts {
		opt(line)
	}
	sb.section.Items = append(sb.section.Items, line)
	return sb
}

// AddPageNumber adds an automatic page number placeholder.
func (sb *SectionBuilder) AddPageNumber(pos Position, opts ...LabelOption) *SectionBuilder {
	pn := &PageNumber{
		Label: Label{
			Pos:      pos,
			FontSize: 10.0,
			FontFace: "Arial",
			Color:    ColorGray,
			Align:    AlignCenter,
		},
	}
	for _, opt := range opts {
		opt(&pn.Label)
	}
	sb.section.Items = append(sb.section.Items, pn)
	return sb
}
