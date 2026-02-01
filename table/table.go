package table

import (
	"strings"

	"github.com/eberle1080/go-rich"
)

// Table represents a table with headers, rows, and borders.
type Table struct {
	columns     []*Column
	rows        [][]string
	title       string
	box         Box
	showHeader  bool
	showEdge    bool
	padding     int
	borderStyle rich.Style
	titleStyle  rich.Style
}

// New creates a new table.
func New() *Table {
	return &Table{
		box:         BoxSimple,
		showHeader:  true,
		showEdge:    true,
		padding:     1,
		borderStyle: rich.NewStyle().Dim(),
		titleStyle:  rich.NewStyle().Bold(),
	}
}

// Title sets the table title.
func (t *Table) Title(title string) *Table {
	t.title = title
	return t
}

// Box sets the border style.
func (t *Table) Box(box Box) *Table {
	t.box = box
	return t
}

// ShowHeader controls whether to show the header row.
func (t *Table) ShowHeader(show bool) *Table {
	t.showHeader = show
	return t
}

// ShowEdge controls whether to show the outer border.
func (t *Table) ShowEdge(show bool) *Table {
	t.showEdge = show
	return t
}

// Padding sets the cell padding.
func (t *Table) Padding(padding int) *Table {
	t.padding = padding
	return t
}

// BorderStyle sets the border style.
func (t *Table) BorderStyle(style rich.Style) *Table {
	t.borderStyle = style
	return t
}

// TitleStyle sets the title style.
func (t *Table) TitleStyle(style rich.Style) *Table {
	t.titleStyle = style
	return t
}

// AddColumn adds a column to the table.
func (t *Table) AddColumn(column *Column) *Table {
	t.columns = append(t.columns, column)
	return t
}

// Headers is a convenience method to add columns from header strings.
func (t *Table) Headers(headers ...string) *Table {
	for _, header := range headers {
		t.AddColumn(NewColumn(header))
	}
	return t
}

// Row adds a row to the table.
func (t *Table) Row(cells ...string) *Table {
	t.rows = append(t.rows, cells)
	return t
}

// Render implements rich.Renderable.
func (t *Table) Render(console *rich.Console, width int) rich.Segments {
	if len(t.columns) == 0 {
		return nil
	}

	// Calculate column widths
	widths := t.calculateWidths(width)

	var segments rich.Segments

	// Render top border
	if t.showEdge {
		segments = append(segments, t.renderTopBorder(widths)...)
		segments = append(segments, rich.Segment{Text: "\n"})
	}

	// Render title if present
	if t.title != "" {
		segments = append(segments, t.renderTitle(widths)...)
		segments = append(segments, rich.Segment{Text: "\n"})
	}

	// Render header
	if t.showHeader {
		segments = append(segments, t.renderHeader(widths)...)
		segments = append(segments, rich.Segment{Text: "\n"})

		// Header separator
		segments = append(segments, t.renderHeaderSeparator(widths)...)
		segments = append(segments, rich.Segment{Text: "\n"})
	}

	// Render rows
	for i, row := range t.rows {
		segments = append(segments, t.renderRow(row, widths)...)
		if i < len(t.rows)-1 {
			segments = append(segments, rich.Segment{Text: "\n"})
		}
	}

	// Render bottom border
	if len(t.rows) > 0 {
		segments = append(segments, rich.Segment{Text: "\n"})
	}
	if t.showEdge {
		segments = append(segments, t.renderBottomBorder(widths)...)
	}

	return segments
}

// calculateWidths determines the width of each column.
func (t *Table) calculateWidths(totalWidth int) []int {
	widths := make([]int, len(t.columns))

	// Start with minimum widths (header length)
	for i, col := range t.columns {
		widths[i] = len(col.Header)
		if col.MinWidth > widths[i] {
			widths[i] = col.MinWidth
		}
	}

	// Check content widths
	for _, row := range t.rows {
		for i := 0; i < len(row) && i < len(widths); i++ {
			cellLen := len(row[i])
			if cellLen > widths[i] {
				widths[i] = cellLen
			}
		}
	}

	// Apply fixed widths and max widths
	for i, col := range t.columns {
		if col.Width > 0 {
			widths[i] = col.Width
		} else if col.MaxWidth > 0 && widths[i] > col.MaxWidth {
			widths[i] = col.MaxWidth
		}
	}

	return widths
}

// renderTopBorder renders the top border of the table.
func (t *Table) renderTopBorder(widths []int) rich.Segments {
	var segments rich.Segments

	if t.showEdge {
		segments = append(segments, rich.Segment{
			Text:  t.box.TopLeft,
			Style: t.borderStyle,
		})
	}

	for i, width := range widths {
		padding := strings.Repeat(t.box.Top, width+t.padding*2)
		segments = append(segments, rich.Segment{
			Text:  padding,
			Style: t.borderStyle,
		})

		if i < len(widths)-1 {
			segments = append(segments, rich.Segment{
				Text:  t.box.MidTop,
				Style: t.borderStyle,
			})
		}
	}

	if t.showEdge {
		segments = append(segments, rich.Segment{
			Text:  t.box.TopRight,
			Style: t.borderStyle,
		})
	}

	return segments
}

// renderBottomBorder renders the bottom border of the table.
func (t *Table) renderBottomBorder(widths []int) rich.Segments {
	var segments rich.Segments

	if t.showEdge {
		segments = append(segments, rich.Segment{
			Text:  t.box.BottomLeft,
			Style: t.borderStyle,
		})
	}

	for i, width := range widths {
		padding := strings.Repeat(t.box.Bottom, width+t.padding*2)
		segments = append(segments, rich.Segment{
			Text:  padding,
			Style: t.borderStyle,
		})

		if i < len(widths)-1 {
			segments = append(segments, rich.Segment{
				Text:  t.box.MidBottom,
				Style: t.borderStyle,
			})
		}
	}

	if t.showEdge {
		segments = append(segments, rich.Segment{
			Text:  t.box.BottomRight,
			Style: t.borderStyle,
		})
	}

	return segments
}

// renderHeaderSeparator renders the separator between header and rows.
func (t *Table) renderHeaderSeparator(widths []int) rich.Segments {
	var segments rich.Segments

	if t.showEdge {
		segments = append(segments, rich.Segment{
			Text:  t.box.HeaderLeft,
			Style: t.borderStyle,
		})
	}

	for i, width := range widths {
		padding := strings.Repeat(t.box.HeaderRow, width+t.padding*2)
		segments = append(segments, rich.Segment{
			Text:  padding,
			Style: t.borderStyle,
		})

		if i < len(widths)-1 {
			segments = append(segments, rich.Segment{
				Text:  t.box.Mid,
				Style: t.borderStyle,
			})
		}
	}

	if t.showEdge {
		segments = append(segments, rich.Segment{
			Text:  t.box.HeaderRight,
			Style: t.borderStyle,
		})
	}

	return segments
}

// renderTitle renders the table title.
func (t *Table) renderTitle(widths []int) rich.Segments {
	totalWidth := 0
	for i, w := range widths {
		totalWidth += w + t.padding*2
		if i < len(widths)-1 {
			totalWidth += 1 // separator
		}
	}

	var segments rich.Segments

	if t.showEdge {
		segments = append(segments, rich.Segment{
			Text:  t.box.Left,
			Style: t.borderStyle,
		})
	}

	titleLen := len(t.title)
	leftPad := (totalWidth - titleLen) / 2
	rightPad := totalWidth - titleLen - leftPad

	if leftPad > 0 {
		segments = append(segments, rich.Segment{Text: strings.Repeat(" ", leftPad)})
	}

	segments = append(segments, rich.Segment{
		Text:  t.title,
		Style: t.titleStyle,
	})

	if rightPad > 0 {
		segments = append(segments, rich.Segment{Text: strings.Repeat(" ", rightPad)})
	}

	if t.showEdge {
		segments = append(segments, rich.Segment{
			Text:  t.box.Right,
			Style: t.borderStyle,
		})
	}

	return segments
}

// renderHeader renders the header row.
func (t *Table) renderHeader(widths []int) rich.Segments {
	var segments rich.Segments

	if t.showEdge {
		segments = append(segments, rich.Segment{
			Text:  t.box.Left,
			Style: t.borderStyle,
		})
	}

	for i, col := range t.columns {
		width := widths[i]

		// Left padding
		segments = append(segments, rich.Segment{
			Text: strings.Repeat(" ", t.padding),
		})

		// Header text (aligned)
		text := t.alignText(col.Header, width, col.Align)
		segments = append(segments, rich.Segment{
			Text:  text,
			Style: col.HeaderStyle,
		})

		// Right padding
		segments = append(segments, rich.Segment{
			Text: strings.Repeat(" ", t.padding),
		})

		// Column separator
		if i < len(t.columns)-1 {
			segments = append(segments, rich.Segment{
				Text:  t.box.Left,
				Style: t.borderStyle,
			})
		}
	}

	if t.showEdge {
		segments = append(segments, rich.Segment{
			Text:  t.box.Right,
			Style: t.borderStyle,
		})
	}

	return segments
}

// renderRow renders a data row.
func (t *Table) renderRow(row []string, widths []int) rich.Segments {
	var segments rich.Segments

	if t.showEdge {
		segments = append(segments, rich.Segment{
			Text:  t.box.Left,
			Style: t.borderStyle,
		})
	}

	for i := 0; i < len(t.columns); i++ {
		col := t.columns[i]
		width := widths[i]

		cellText := ""
		if i < len(row) {
			cellText = row[i]
		}

		// Left padding
		segments = append(segments, rich.Segment{
			Text: strings.Repeat(" ", t.padding),
		})

		// Cell text (aligned and truncated if needed)
		if len(cellText) > width {
			cellText = cellText[:width]
		}
		text := t.alignText(cellText, width, col.Align)
		segments = append(segments, rich.Segment{
			Text:  text,
			Style: col.CellStyle,
		})

		// Right padding
		segments = append(segments, rich.Segment{
			Text: strings.Repeat(" ", t.padding),
		})

		// Column separator
		if i < len(t.columns)-1 {
			segments = append(segments, rich.Segment{
				Text:  t.box.Left,
				Style: t.borderStyle,
			})
		}
	}

	if t.showEdge {
		segments = append(segments, rich.Segment{
			Text:  t.box.Right,
			Style: t.borderStyle,
		})
	}

	return segments
}

// alignText aligns text within a given width.
func (t *Table) alignText(text string, width int, align Align) string {
	textLen := len(text)
	if textLen >= width {
		return text
	}

	padding := width - textLen

	switch align {
	case AlignLeft:
		return text + strings.Repeat(" ", padding)
	case AlignRight:
		return strings.Repeat(" ", padding) + text
	case AlignCenter:
		leftPad := padding / 2
		rightPad := padding - leftPad
		return strings.Repeat(" ", leftPad) + text + strings.Repeat(" ", rightPad)
	default:
		return text + strings.Repeat(" ", padding)
	}
}
