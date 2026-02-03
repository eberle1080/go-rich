package rich

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/eberle1080/go-rich/internal/ansi"
	"golang.org/x/term"
)

// Console is the central orchestrator for all rich terminal output.
// It manages the output writer, color mode detection, and terminal dimensions.
//
// A Console is typically created once per application using NewConsole and then
// reused for all output operations. It provides methods for:
//   - Plain text output (Print, Println, Printf)
//   - Styled text output (PrintStyled, PrintStyledln)
//   - Markup output (PrintMarkup, PrintMarkupln)
//   - Segment output (PrintSegments, PrintSegmentsln)
//   - Renderable output (Render, Renderln)
//   - Horizontal rules (Rule)
//
// The Console automatically detects terminal color capabilities and terminal
// dimensions, but both can be overridden if needed.
type Console struct {
	writer    io.Writer // Underlying writer (usually os.Stdout)
	colorMode ColorMode // Detected or explicitly set color mode
	width     int       // Terminal width in characters
	height    int       // Terminal height in characters
}

// NewConsole creates a new Console writing to the specified writer.
// If writer is nil, os.Stdout is used.
//
// The console automatically:
//   - Detects color mode based on environment variables and terminal capabilities
//   - Queries terminal dimensions if the writer is a terminal
//   - Falls back to 80×24 if dimensions cannot be determined
//
// Example:
//
//	// Use stdout with auto-detection
//	console := rich.NewConsole(nil)
//
//	// Use a custom writer
//	var buf bytes.Buffer
//	console := rich.NewConsole(&buf)
//
//	// Use stderr
//	console := rich.NewConsole(os.Stderr)
func NewConsole(writer io.Writer) *Console {
	// Default to stdout if no writer specified
	if writer == nil {
		writer = os.Stdout
	}

	console := &Console{
		writer:    writer,
		colorMode: detectColorMode(writer),
		width:     80, // Default terminal width
		height:    24, // Default terminal height
	}

	// Try to get actual terminal size if writer is a terminal file
	if f, ok := writer.(*os.File); ok {
		if w, h, err := term.GetSize(int(f.Fd())); err == nil {
			console.width = w
			console.height = h
		}
	}

	return console
}

// detectColorMode determines the color support level of the terminal.
// This function implements the color detection hierarchy:
//
//  1. NO_COLOR environment variable → ColorModeNone
//     Respects the NO_COLOR standard (https://no-color.org/)
//
//  2. Non-terminal writer → ColorModeNone
//     Colors are disabled when writing to files, pipes, etc.
//
//  3. COLORTERM=truecolor or COLORTERM=24bit → ColorModeTrueColor
//     Most modern terminals set this for 24-bit RGB support
//
//  4. TERM contains "256color" → ColorMode256
//     Common values: xterm-256color, screen-256color
//
//  5. TERM is set and not "dumb" → ColorModeStandard
//     Fallback to 16-color ANSI for basic terminals
//
//  6. Otherwise → ColorModeNone
//     Unknown or dumb terminals get no colors
func detectColorMode(w io.Writer) ColorMode {
	// Check NO_COLOR environment variable (https://no-color.org/)
	// If set to any value (even empty string), disable colors
	if os.Getenv("NO_COLOR") != "" {
		return ColorModeNone
	}

	// Check if writer is a terminal file
	f, ok := w.(*os.File)
	if !ok {
		// Not a file (e.g., bytes.Buffer), can't detect colors
		return ColorModeNone
	}

	// Check if it's actually a terminal (not a pipe or file)
	if !term.IsTerminal(int(f.Fd())) {
		return ColorModeNone
	}

	// Check COLORTERM for truecolor support
	// Many modern terminals set this: GNOME Terminal, iTerm2, VS Code, etc.
	colorTerm := os.Getenv("COLORTERM")
	if colorTerm == "truecolor" || colorTerm == "24bit" {
		return ColorModeTrueColor
	}

	// Check TERM for 256 color support
	// Common values: xterm-256color, screen-256color, tmux-256color
	termEnv := os.Getenv("TERM")
	if strings.Contains(termEnv, "256color") {
		return ColorMode256
	}

	// Default to standard ANSI colors for known terminals
	// "dumb" terminals get no colors
	if termEnv != "" && termEnv != "dumb" {
		return ColorModeStandard
	}

	// Unknown terminal, disable colors
	return ColorModeNone
}

// SetColorMode overrides the detected color mode.
// Use this to force a specific color mode instead of relying on auto-detection.
//
// Common use cases:
//   - Testing color output in different modes
//   - Forcing colors even when not detected (e.g., in CI)
//   - Disabling colors explicitly
//
// Example:
//
//	console.SetColorMode(rich.ColorModeTrueColor)  // Force 24-bit color
//	console.SetColorMode(rich.ColorModeNone)       // Disable all colors
func (c *Console) SetColorMode(mode ColorMode) {
	c.colorMode = mode
}

// ColorMode returns the current color mode.
// This will be either the auto-detected mode or one set via SetColorMode.
//
// Example:
//
//	if console.ColorMode() == rich.ColorModeTrueColor {
//		// Terminal supports full RGB colors
//	}
func (c *Console) ColorMode() ColorMode {
	return c.colorMode
}

// Width returns the console width in characters.
// This is either the detected terminal width or the default of 80.
//
// Used by renderables like tables and panels to determine layout.
//
// Example:
//
//	maxWidth := console.Width()
func (c *Console) Width() int {
	return c.width
}

// Height returns the console height in characters.
// This is either the detected terminal height or the default of 24.
//
// Example:
//
//	maxHeight := console.Height()
func (c *Console) Height() int {
	return c.height
}

// Print writes plain text to the console without styling.
// Behaves like fmt.Print, writing the string representation of the arguments.
// Returns the number of bytes written and any write error.
//
// Example:
//
//	console.Print("Hello", " ", "world")
func (c *Console) Print(a ...interface{}) (n int, err error) {
	s := fmt.Sprint(a...)
	return c.writer.Write([]byte(s))
}

// Println writes plain text to the console followed by a newline.
// Behaves like fmt.Println, adding a newline after the output.
// Returns the number of bytes written and any write error.
//
// Example:
//
//	console.Println("Hello world")
func (c *Console) Println(a ...interface{}) (n int, err error) {
	s := fmt.Sprintln(a...)
	return c.writer.Write([]byte(s))
}

// Printf writes formatted text to the console.
// Behaves like fmt.Printf, using format string and arguments.
// Returns the number of bytes written and any write error.
//
// Example:
//
//	console.Printf("Count: %d\n", 42)
func (c *Console) Printf(format string, a ...interface{}) (n int, err error) {
	s := fmt.Sprintf(format, a...)
	return c.writer.Write([]byte(s))
}

// PrintStyled writes styled text to the console.
// The text is rendered with ANSI escape sequences according to the console's color mode.
// Returns the number of bytes written and any write error.
//
// Example:
//
//	style := rich.NewStyle().Bold().Foreground(rich.Red)
//	styled := style.Render("Error!")
//	console.PrintStyled(styled)
func (c *Console) PrintStyled(st StyledText) (n int, err error) {
	seg := Segment{Text: st.Text, Style: st.Style}
	return c.PrintSegments(Segments{seg})
}

// PrintStyledln writes styled text to the console followed by a newline.
// The text is rendered with ANSI escape sequences according to the console's color mode.
// Returns the total number of bytes written (text + newline) and any write error.
//
// Example:
//
//	style := rich.NewStyle().Bold().Foreground(rich.Green)
//	styled := style.Render("Success!")
//	console.PrintStyledln(styled)
func (c *Console) PrintStyledln(st StyledText) (n int, err error) {
	// Write the styled text
	n, err = c.PrintStyled(st)
	if err != nil {
		return n, err
	}

	// Add newline
	n2, err := c.writer.Write([]byte("\n"))
	return n + n2, err
}

// PrintSegments writes segments to the console.
// Each segment is rendered with its own style, and the segments are concatenated.
// This is a lower-level API typically used by renderables or advanced use cases.
// Returns the number of bytes written and any write error.
//
// Example:
//
//	segments := rich.Segments{
//		{Text: "Hello ", Style: rich.NewStyle().Bold()},
//		{Text: "world", Style: rich.NewStyle().Foreground(rich.Red)},
//	}
//	console.PrintSegments(segments)
func (c *Console) PrintSegments(segments Segments) (n int, err error) {
	s := segments.ToANSI(c.colorMode)
	return c.writer.Write([]byte(s))
}

// PrintSegmentsln writes segments to the console followed by a newline.
// Each segment is rendered with its own style, and a newline is appended.
// Returns the total number of bytes written and any write error.
//
// Example:
//
//	segments := rich.Segments{
//		{Text: "Status: ", Style: rich.NewStyle().Bold()},
//		{Text: "OK", Style: rich.NewStyle().Foreground(rich.Green)},
//	}
//	console.PrintSegmentsln(segments)
func (c *Console) PrintSegmentsln(segments Segments) (n int, err error) {
	// Write the segments
	n, err = c.PrintSegments(segments)
	if err != nil {
		return n, err
	}

	// Add newline
	n2, err := c.writer.Write([]byte("\n"))
	return n + n2, err
}

// Rule prints a horizontal rule across the console.
// If a title is provided, it's centered in the rule with padding on both sides.
// The rule extends across the full console width.
//
// Rules are useful for visually separating sections of output.
// The line is rendered in dim style, and the title (if any) in bold.
//
// Example:
//
//	console.Rule("")           // Plain horizontal line
//	console.Rule("Section 1")  // ─────── Section 1 ───────
func (c *Console) Rule(title string) (n int, err error) {
	var segments Segments

	if title == "" {
		// No title: just print a full-width line
		line := strings.Repeat("─", c.width)
		segments = Segments{{Text: line, Style: NewStyle().Dim()}}
	} else {
		// With title: center it with lines on both sides
		titleLen := len(title)

		// Check if title fits with padding (at least 2 chars on each side)
		if titleLen+4 > c.width {
			// Title too long, just print it without the rule
			segments = Segments{{Text: title, Style: NewStyle().Bold()}}
		} else {
			// Calculate line lengths on each side
			// Format: "─────── Title ───────"
			// Title has 1 space on each side
			leftLen := (c.width - titleLen - 2) / 2
			rightLen := c.width - titleLen - 2 - leftLen

			segments = Segments{
				{Text: strings.Repeat("─", leftLen), Style: NewStyle().Dim()},
				{Text: " " + title + " ", Style: NewStyle().Bold()},
				{Text: strings.Repeat("─", rightLen), Style: NewStyle().Dim()},
			}
		}
	}

	return c.PrintSegmentsln(segments)
}

// Writer returns the underlying io.Writer.
// This provides direct access to the output destination for advanced use cases.
//
// Example:
//
//	w := console.Writer()
//	w.Write([]byte("raw output\n"))
func (c *Console) Writer() io.Writer {
	return c.writer
}

// ANSIWriter returns an ANSI-aware writer.
// The returned writer provides utilities for working with ANSI escape sequences.
//
// Example:
//
//	w := console.ANSIWriter()
//	w.WriteString("text")
func (c *Console) ANSIWriter() *ansi.Writer {
	return ansi.NewWriter(c.writer)
}

// PrintMarkup writes markup text to the console.
// Markup provides an easy way to add inline styling using tags.
//
// Markup format: [style]text[/] where style can be:
//   - Color names: red, blue, green, yellow, magenta, cyan, white, black, etc.
//   - Hex colors: #FF0000, #00FF00
//   - RGB colors: rgb(255,0,0)
//   - Attributes: bold, italic, underline, strikethrough, dim, reverse
//   - Background: "red on blue", "bold on white"
//   - Combined: "bold red", "italic blue on yellow"
//
// Special characters:
//   - [/] closes the current tag
//   - [[ escapes to a literal [
//
// Examples:
//
//	console.PrintMarkup("[bold]Bold text[/]")
//	console.PrintMarkup("[red]Red[/] and [blue]blue[/]")
//	console.PrintMarkup("[bold red on white]Styled[/]")
//	console.PrintMarkup("[[This is not a tag]]") // Prints: [This is not a tag]
func (c *Console) PrintMarkup(m string) (n int, err error) {
	return c.printMarkupInternal(m)
}

// PrintMarkupln writes markup text to the console followed by a newline.
// See PrintMarkup for markup syntax documentation.
//
// Example:
//
//	console.PrintMarkupln("[green]Success![/]")
func (c *Console) PrintMarkupln(markup string) (n int, err error) {
	// Write the markup
	n, err = c.PrintMarkup(markup)
	if err != nil {
		return n, err
	}

	// Add newline
	n2, err := c.writer.Write([]byte("\n"))
	return n + n2, err
}

// Render renders a Renderable to the console.
// The renderable is converted to segments using the console's width,
// then the segments are printed.
//
// This is the primary method for rendering complex widgets like tables and panels.
//
// Example:
//
//	table := table.New().Headers("Name", "Age").Row("Alice", "30")
//	console.Render(table)
//
//	panel := panel.New("Content").Title("Box")
//	console.Render(panel)
func (c *Console) Render(r Renderable) (n int, err error) {
	segments := r.Render(c, c.width)
	return c.PrintSegments(segments)
}

// Renderln renders a Renderable to the console followed by a newline.
// This is the same as Render but adds a trailing newline for convenience.
//
// Example:
//
//	table := table.New().Headers("Name", "Age").Row("Alice", "30")
//	console.Renderln(table)
func (c *Console) Renderln(r Renderable) (n int, err error) {
	// Render the content
	n, err = c.Render(r)
	if err != nil {
		return n, err
	}

	// Add newline
	n2, err := c.writer.Write([]byte("\n"))
	return n + n2, err
}
