// Package rich provides a Go library for rich text and beautiful formatting in the terminal.
//
// This package is inspired by the Python Rich library and brings similar capabilities to Go,
// enabling developers to create visually appealing terminal output with minimal effort.
//
// # Features
//
// - Multiple color modes: TrueColor (24-bit RGB), 256-color, standard ANSI, and plain text
// - Automatic color mode detection based on terminal capabilities
// - Styled text with colors, bold, italic, underline, strikethrough, dim, and reverse
// - Markup language for inline styling (e.g., "[bold red]Error[/]")
// - Tables with customizable borders, alignment, and styling
// - Panels (bordered containers) for grouping content
// - Horizontal rules and dividers
// - Renderable interface for custom content types
//
// # Basic Usage
//
// Create a console and print styled output:
//
//	console := rich.NewConsole(nil) // Uses os.Stdout by default
//	console.PrintMarkupln("[bold green]Success![/] Operation completed.")
//
// # Styling
//
// Styles can be created using the fluent API:
//
//	style := rich.NewStyle().
//		Foreground(rich.Red).
//		Bold().
//		Underline()
//
//	console.PrintStyledln(style.Render("Important message"))
//
// # Colors
//
// The package supports multiple color types:
//
//	// Standard ANSI colors (0-15)
//	rich.Red
//	rich.BrightBlue
//
//	// 256-color palette
//	rich.ANSI256Color(196)
//
//	// True color RGB
//	rich.RGB(255, 100, 50)
//	rich.Hex("#FF6432")
//	rich.Named("orange")
//
// # Markup
//
// Use markup tags for inline styling:
//
//	console.PrintMarkupln("[bold]Bold text[/]")
//	console.PrintMarkupln("[red]Red text[/]")
//	console.PrintMarkupln("[bold red on white]Styled[/]")
//	console.PrintMarkupln("[#FF0000]Hex color[/]")
//	console.PrintMarkupln("[rgb(255,0,0)]RGB color[/]")
//
// Escape brackets with double brackets:
//
//	console.PrintMarkupln("[[This is not a tag]]") // Prints: [This is not a tag]
//
// # Tables
//
// Create formatted tables with borders and styling:
//
//	import "github.com/eberle1080/go-rich/table"
//
//	tbl := table.New().
//		Headers("Name", "Age", "City").
//		Row("Alice", "30", "NYC").
//		Row("Bob", "25", "LA")
//
//	console.Renderln(tbl)
//
// # Panels
//
// Wrap content in bordered containers:
//
//	import "github.com/eberle1080/go-rich/panel"
//
//	p := panel.New("Important message").
//		Title("Alert").
//		BorderStyle(rich.NewStyle().Foreground(rich.Red))
//
//	console.Renderln(p)
//
// # Color Mode Detection
//
// The package automatically detects the terminal's color capabilities:
//
//   - Checks NO_COLOR environment variable (https://no-color.org/)
//   - Detects true color support via COLORTERM=truecolor or COLORTERM=24bit
//   - Detects 256-color support via TERM=*256color*
//   - Falls back to standard ANSI colors or no colors
//
// You can override the detected mode:
//
//	console.SetColorMode(rich.ColorModeTrueColor)
//
// # Renderables
//
// Implement the Renderable interface to create custom content:
//
//	type Renderable interface {
//		Render(console *Console, width int) Segments
//	}
//
// For advanced sizing control, also implement Measurable:
//
//	type Measurable interface {
//		Measure(console *Console, maxWidth int) Measurement
//	}
package rich
