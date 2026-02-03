package table

import (
	"strings"

	"github.com/eberle1080/go-rich"
)

// Table represents a table with headers, rows, and borders.
// Tables provide a structured way to display tabular data with customizable
// appearance. They implement the rich.Renderable interface and can be rendered
// to a Console.
//
// A table consists of:
//   - Columns: Define headers, widths, alignment, and styles
//   - Rows: Data organized as string arrays (one per row)
//   - Optional title: Displayed above the table
//   - Border style (Box): Visual appearance of table borders
//   - Various display options: header visibility, edge visibility, padding
//
// Tables are built using a fluent API:
//
//	table.New().
//		Title("Users").
//		Headers("Name", "Age").
//		Row("Alice", "30").
//		Row("Bob", "25")
type Table struct {
	columns []*Column  // Column configurations (headers, styles, widths)
	rows    [][]string // Data rows (each row is array of cell values)

	title string // Optional title displayed at top
	box   Box    // Border characters to use

	showHeader bool // Whether to display the header row
	showEdge   bool // Whether to display outer borders

	padding int // Cell padding (spaces on left/right of content)

	borderStyle rich.Style // Style applied to border characters
	titleStyle  rich.Style // Style applied to the title
}

// New creates a new table with sensible defaults.
// The table is initialized with:
//   - BoxSimple border style (clean single-line borders)
//   - Header row visible
//   - Outer borders visible
//   - 1 character of padding in each cell
//   - Dim style for borders
//   - Bold style for title
//
// Use the fluent API methods to customize the table before rendering.
//
// Example:
//
//	tbl := table.New().
//		Headers("Name", "Age").
//		Row("Alice", "30").
//		Row("Bob", "25")
//	console.Renderln(tbl)
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

// Title sets the table title displayed at the top.
// The title is centered and displayed in its own row above the headers.
// If empty (default), no title row is shown.
//
// Example:
//
//	tbl := table.New().Title("User Report")
func (t *Table) Title(title string) *Table {
	t.title = title
	return t
}

// Box sets the border style using predefined or custom box characters.
// See the Box type and predefined styles (BoxSimple, BoxRounded, etc.) for options.
//
// Example:
//
//	tbl := table.New().Box(table.BoxDouble)
func (t *Table) Box(box Box) *Table {
	t.box = box
	return t
}

// ShowHeader controls whether to show the header row.
// When false, only data rows are displayed (no header row or header separator).
// Default is true.
//
// Example:
//
//	tbl := table.New().ShowHeader(false)
func (t *Table) ShowHeader(show bool) *Table {
	t.showHeader = show
	return t
}

// ShowEdge controls whether to show the outer border.
// When false, only internal column separators are shown (no top/bottom/left/right edges).
// Default is true.
//
// Example:
//
//	tbl := table.New().ShowEdge(false)
func (t *Table) ShowEdge(show bool) *Table {
	t.showEdge = show
	return t
}

// Padding sets the cell padding in characters.
// Padding is added to both left and right sides of cell content.
// Default is 1.
//
// Example:
//
//	tbl := table.New().Padding(2) // 2 spaces on each side
func (t *Table) Padding(padding int) *Table {
	t.padding = padding
	return t
}

// BorderStyle sets the style for all border characters.
// This affects the visual appearance of the borders but not their shape.
// Default is dim (faint) style.
//
// Example:
//
//	tbl := table.New().BorderStyle(rich.NewStyle().Foreground(rich.Cyan))
func (t *Table) BorderStyle(style rich.Style) *Table {
	t.borderStyle = style
	return t
}

// TitleStyle sets the style for the title text.
// Only affects the title if one is set via Title().
// Default is bold style.
//
// Example:
//
//	tbl := table.New().
//		Title("Report").
//		TitleStyle(rich.NewStyle().Bold().Foreground(rich.Green))
func (t *Table) TitleStyle(style rich.Style) *Table {
	t.titleStyle = style
	return t
}

// AddColumn adds a column to the table.
// Use this when you need fine-grained control over column configuration.
// For simple cases, use Headers() instead.
//
// Example:
//
//	col := table.NewColumn("Price").
//		WithAlign(table.AlignRight).
//		WithWidth(10)
//	tbl := table.New().AddColumn(col)
func (t *Table) AddColumn(column *Column) *Table {
	t.columns = append(t.columns, column)
	return t
}

// Headers is a convenience method to add columns from header strings.
// Creates one column per header with default settings (left-aligned, auto-width, bold header).
// This is the quickest way to set up a simple table.
//
// Example:
//
//	tbl := table.New().Headers("Name", "Age", "City")
func (t *Table) Headers(headers ...string) *Table {
	for _, header := range headers {
		t.AddColumn(NewColumn(header))
	}
	return t
}

// Row adds a data row to the table.
// Each argument becomes a cell in the row, matched to columns by position.
// If fewer cells than columns are provided, remaining cells are empty.
// If more cells than columns are provided, extra cells are ignored.
//
// Example:
//
//	tbl := table.New().
//		Headers("Name", "Age").
//		Row("Alice", "30").
//		Row("Bob", "25")
func (t *Table) Row(cells ...string) *Table {
	t.rows = append(t.rows, cells)
	return t
}

// Render implements rich.Renderable.
// Converts the table into styled segments that can be displayed on the console.
//
// The rendering process:
//  1. Calculate optimal column widths based on content and constraints
//  2. Render top border (if showEdge is true)
//  3. Render title row (if title is set)
//  4. Render header row (if showHeader is true)
//  5. Render header separator
//  6. Render data rows
//  7. Render bottom border (if showEdge is true)
//
// The width parameter is the maximum available width for the table.
// The console parameter provides access to the color mode and other settings.
func (t *Table) Render(console *rich.Console, width int) rich.Segments {
	// Empty table with no columns
	if len(t.columns) == 0 {
		return nil
	}

	// Calculate optimal widths for each column
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

// calculateWidths determines the optimal width for each column.
// The algorithm:
//  1. Start with the maximum of header length and MinWidth for each column
//  2. Expand widths to fit the longest content in each column
//  3. Apply Width (fixed) or MaxWidth (ceiling) constraints
//
// This ensures:
//   - Headers are fully visible (unless overridden by Width/MaxWidth)
//   - Content fits when possible (within MaxWidth constraints)
//   - MinWidth constraints are respected
//   - Fixed-width columns (Width > 0) override calculated widths
//
// Note: The totalWidth parameter is currently unused. Future implementations
// may use it to proportionally shrink columns when total width exceeds available space.
func (t *Table) calculateWidths(totalWidth int) []int {
	widths := make([]int, len(t.columns))

	// Phase 1: Initialize with maximum of header length and MinWidth
	for i, col := range t.columns {
		widths[i] = len(col.Header)
		if col.MinWidth > widths[i] {
			widths[i] = col.MinWidth
		}
	}

	// Phase 2: Expand to fit content (longest cell in each column)
	for _, row := range t.rows {
		for i := 0; i < len(row) && i < len(widths); i++ {
			cellLen := len(row[i])
			if cellLen > widths[i] {
				widths[i] = cellLen
			}
		}
	}

	// Phase 3: Apply width constraints
	for i, col := range t.columns {
		if col.Width > 0 {
			// Fixed width overrides calculated width
			widths[i] = col.Width
		} else if col.MaxWidth > 0 && widths[i] > col.MaxWidth {
			// MaxWidth provides a ceiling
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
// Adds padding spaces to position the text according to the alignment setting.
//
// Behavior:
//   - If text is already >= width, returns text unchanged (no truncation here)
//   - AlignLeft: text + spaces
//   - AlignRight: spaces + text
//   - AlignCenter: (leftPad) + text + (rightPad), where leftPad = padding/2
//
// The text parameter is the content to align.
// The width parameter is the target width in characters.
// The align parameter specifies the alignment strategy.
//
// Returns the text padded to exactly the specified width.
func (t *Table) alignText(text string, width int, align Align) string {
	textLen := len(text)

	// Text already fills or exceeds the width
	if textLen >= width {
		return text
	}

	// Calculate how many spaces to add
	padding := width - textLen

	switch align {
	case AlignLeft:
		// "text    "
		return text + strings.Repeat(" ", padding)

	case AlignRight:
		// "    text"
		return strings.Repeat(" ", padding) + text

	case AlignCenter:
		// "  text  " (balanced padding)
		leftPad := padding / 2
		rightPad := padding - leftPad
		return strings.Repeat(" ", leftPad) + text + strings.Repeat(" ", rightPad)

	default:
		// Fallback to left alignment
		return text + strings.Repeat(" ", padding)
	}
}
