package table

import "github.com/eberle1080/go-rich"

// Align specifies how content is aligned within a column.
// Alignment affects how text is positioned when the cell content
// is shorter than the column width.
type Align int

const (
	// AlignLeft aligns content to the left side of the column.
	// This is the default and most common alignment for text content.
	AlignLeft Align = iota

	// AlignCenter centers content within the column.
	// Useful for headers, short labels, or creating visual symmetry.
	AlignCenter

	// AlignRight aligns content to the right side of the column.
	// Commonly used for numeric data to align decimal points.
	AlignRight
)

// Column represents a table column configuration.
// A column defines how a particular column in a table should be rendered,
// including its header text, width constraints, alignment, and styling.
//
// Width behavior:
//   - If Width > 0: Column has a fixed width
//   - If Width == 0: Width is calculated from content, constrained by MinWidth/MaxWidth
//   - MinWidth sets a floor (column won't be narrower than this)
//   - MaxWidth sets a ceiling (column won't be wider than this, 0 = unlimited)
//
// Styling is applied separately to headers and data cells, allowing for
// visual distinction between them.
type Column struct {
	Header string // Header text displayed in the header row

	Width    int // Fixed width in characters (0 = auto-size from content)
	MinWidth int // Minimum width in characters (0 = no minimum)
	MaxWidth int // Maximum width in characters (0 = unlimited)

	Align Align // How content is aligned within the column

	HeaderStyle rich.Style // Style applied to the header cell
	CellStyle   rich.Style // Style applied to data cells in this column

	NoWrap bool // If true, prevent text wrapping (truncate instead)
}

// NewColumn creates a new column with the given header.
// The column is initialized with sensible defaults:
//   - Auto-sized width (calculated from content)
//   - Left alignment
//   - Bold header style
//   - Unstyled cell content
//
// Use the With* methods to customize the column configuration.
//
// Example:
//
//	col := table.NewColumn("Name")
func NewColumn(header string) *Column {
	return &Column{
		Header:      header,
		Align:       AlignLeft,
		HeaderStyle: rich.NewStyle().Bold(),
		CellStyle:   rich.NewStyle(),
	}
}

// WithWidth sets a fixed width for the column.
// When set, the column will always be exactly this width,
// regardless of content. Content longer than the width will be truncated.
//
// Example:
//
//	col := table.NewColumn("Name").WithWidth(20)
func (c *Column) WithWidth(width int) *Column {
	c.Width = width
	return c
}

// WithMinWidth sets the minimum width for the column.
// The column will be at least this wide, even if content is shorter.
// Useful for ensuring columns don't become too narrow.
//
// Example:
//
//	col := table.NewColumn("ID").WithMinWidth(5)
func (c *Column) WithMinWidth(width int) *Column {
	c.MinWidth = width
	return c
}

// WithMaxWidth sets the maximum width for the column.
// The column won't exceed this width. Content longer than this will be truncated.
// Set to 0 for unlimited width.
//
// Example:
//
//	col := table.NewColumn("Description").WithMaxWidth(50)
func (c *Column) WithMaxWidth(width int) *Column {
	c.MaxWidth = width
	return c
}

// WithAlign sets the alignment for the column.
// Affects both header and cell content alignment.
//
// Example:
//
//	col := table.NewColumn("Price").WithAlign(table.AlignRight)
func (c *Column) WithAlign(align Align) *Column {
	c.Align = align
	return c
}

// WithHeaderStyle sets the header style for the column.
// This style is applied only to the header cell in this column.
//
// Example:
//
//	col := table.NewColumn("Status").
//		WithHeaderStyle(rich.NewStyle().Foreground(rich.Yellow))
func (c *Column) WithHeaderStyle(style rich.Style) *Column {
	c.HeaderStyle = style
	return c
}

// WithCellStyle sets the cell style for the column.
// This style is applied to all data cells in this column.
//
// Example:
//
//	col := table.NewColumn("Error").
//		WithCellStyle(rich.NewStyle().Foreground(rich.Red))
func (c *Column) WithCellStyle(style rich.Style) *Column {
	c.CellStyle = style
	return c
}

// WithNoWrap disables text wrapping for the column.
// When set, long text will be truncated instead of wrapping to multiple lines.
// Currently this is the only supported behavior (wrapping is not yet implemented).
//
// Example:
//
//	col := table.NewColumn("Path").WithNoWrap()
func (c *Column) WithNoWrap() *Column {
	c.NoWrap = true
	return c
}
