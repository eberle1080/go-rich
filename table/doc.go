// Package table provides formatted table rendering for terminal output.
//
// This package enables the creation of beautiful, styled tables with customizable
// borders, alignment, and formatting. Tables are rendered as Renderables and can
// be displayed using a rich.Console.
//
// # Basic Usage
//
// Create a simple table with headers and rows:
//
//	tbl := table.New().
//		Headers("Name", "Age", "City").
//		Row("Alice", "30", "New York").
//		Row("Bob", "25", "Los Angeles").
//		Row("Carol", "35", "Chicago")
//
//	console.Renderln(tbl)
//
// # Column Configuration
//
// Customize individual columns with detailed configuration:
//
//	col1 := table.NewColumn("Name").
//		WithAlign(table.AlignLeft).
//		WithMinWidth(10)
//
//	col2 := table.NewColumn("Age").
//		WithAlign(table.AlignRight).
//		WithWidth(5)
//
//	col3 := table.NewColumn("City").
//		WithAlign(table.AlignCenter).
//		WithMaxWidth(20)
//
//	tbl := table.New().
//		AddColumn(col1).
//		AddColumn(col2).
//		AddColumn(col3).
//		Row("Alice", "30", "New York")
//
// # Border Styles
//
// Choose from predefined border styles or create custom ones:
//
//	tbl.Box(table.BoxRounded)   // Rounded corners (default)
//	tbl.Box(table.BoxDouble)    // Double-line borders
//	tbl.Box(table.BoxHeavy)     // Heavy borders
//	tbl.Box(table.BoxSimple)    // Simple single-line borders
//	tbl.Box(table.BoxASCII)     // ASCII-only characters
//	tbl.Box(table.BoxNone)      // No borders
//
// # Styling
//
// Apply custom styles to borders, headers, and titles:
//
//	tbl := table.New().
//		Title("User Report").
//		TitleStyle(rich.NewStyle().Bold().Foreground(rich.Cyan)).
//		BorderStyle(rich.NewStyle().Dim()).
//		Headers("Name", "Status")
//
// Column-specific styling:
//
//	col := table.NewColumn("Status").
//		WithHeaderStyle(rich.NewStyle().Bold().Foreground(rich.Yellow)).
//		WithCellStyle(rich.NewStyle().Italic())
//
// # Alignment
//
// Control how content is aligned within columns:
//
//	table.AlignLeft    // Left-aligned (default)
//	table.AlignCenter  // Centered
//	table.AlignRight   // Right-aligned
//
// # Width Control
//
// Tables support flexible width management:
//
//	col.WithWidth(20)      // Fixed width of 20 characters
//	col.WithMinWidth(10)   // Minimum width of 10 characters
//	col.WithMaxWidth(50)   // Maximum width of 50 characters
//
// Width calculation:
//   - If Width is set, the column uses that exact width
//   - Otherwise, width is calculated from content and constrained by MinWidth/MaxWidth
//   - By default, columns auto-size to fit their content
//
// # Options
//
// Additional table configuration:
//
//	tbl.ShowHeader(false)   // Hide the header row
//	tbl.ShowEdge(false)     // Hide outer borders
//	tbl.Padding(2)          // Set cell padding (default: 1)
//
// # Custom Box Styles
//
// Create custom border styles:
//
//	customBox := table.Box{
//		TopLeft:     "╔",
//		Top:         "═",
//		TopRight:    "╗",
//		Left:        "║",
//		Right:       "║",
//		BottomLeft:  "╚",
//		Bottom:      "═",
//		BottomRight: "╝",
//		// ... other fields
//	}
//
//	tbl.Box(customBox)
//
// # Complete Example
//
//	import (
//		"github.com/eberle1080/go-rich"
//		"github.com/eberle1080/go-rich/table"
//	)
//
//	console := rich.NewConsole(nil)
//
//	tbl := table.New().
//		Title("Sales Report").
//		Box(table.BoxDouble).
//		BorderStyle(rich.NewStyle().Foreground(rich.Cyan)).
//		Headers("Product", "Quantity", "Price").
//		Row("Widget A", "150", "$29.99").
//		Row("Widget B", "200", "$39.99").
//		Row("Widget C", "75", "$49.99")
//
//	console.Renderln(tbl)
package table
