package table

import "github.com/eberle1080/go-rich"

// Align specifies how content is aligned within a column.
type Align int

const (
	AlignLeft Align = iota
	AlignCenter
	AlignRight
)

// Column represents a table column configuration.
type Column struct {
	Header    string
	Width     int  // Fixed width (0 = auto)
	MinWidth  int  // Minimum width
	MaxWidth  int  // Maximum width (0 = no limit)
	Align     Align
	HeaderStyle rich.Style
	CellStyle   rich.Style
	NoWrap    bool
}

// NewColumn creates a new column with the given header.
func NewColumn(header string) *Column {
	return &Column{
		Header:      header,
		Align:       AlignLeft,
		HeaderStyle: rich.NewStyle().Bold(),
		CellStyle:   rich.NewStyle(),
	}
}

// WithWidth sets a fixed width for the column.
func (c *Column) WithWidth(width int) *Column {
	c.Width = width
	return c
}

// WithMinWidth sets the minimum width for the column.
func (c *Column) WithMinWidth(width int) *Column {
	c.MinWidth = width
	return c
}

// WithMaxWidth sets the maximum width for the column.
func (c *Column) WithMaxWidth(width int) *Column {
	c.MaxWidth = width
	return c
}

// WithAlign sets the alignment for the column.
func (c *Column) WithAlign(align Align) *Column {
	c.Align = align
	return c
}

// WithHeaderStyle sets the header style for the column.
func (c *Column) WithHeaderStyle(style rich.Style) *Column {
	c.HeaderStyle = style
	return c
}

// WithCellStyle sets the cell style for the column.
func (c *Column) WithCellStyle(style rich.Style) *Column {
	c.CellStyle = style
	return c
}

// WithNoWrap disables text wrapping for the column.
func (c *Column) WithNoWrap() *Column {
	c.NoWrap = true
	return c
}
