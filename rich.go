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
type Console struct {
	writer    io.Writer
	colorMode ColorMode
	width     int
	height    int
}

// NewConsole creates a new Console writing to the specified writer.
// If writer is nil, os.Stdout is used.
func NewConsole(writer io.Writer) *Console {
	if writer == nil {
		writer = os.Stdout
	}

	console := &Console{
		writer:    writer,
		colorMode: detectColorMode(writer),
		width:     80,
		height:    24,
	}

	// Try to get actual terminal size
	if f, ok := writer.(*os.File); ok {
		if w, h, err := term.GetSize(int(f.Fd())); err == nil {
			console.width = w
			console.height = h
		}
	}

	return console
}

// detectColorMode determines the color support level of the terminal.
func detectColorMode(w io.Writer) ColorMode {
	// Check NO_COLOR environment variable (https://no-color.org/)
	if os.Getenv("NO_COLOR") != "" {
		return ColorModeNone
	}

	// Check if writer is a terminal
	f, ok := w.(*os.File)
	if !ok {
		return ColorModeNone
	}

	if !term.IsTerminal(int(f.Fd())) {
		return ColorModeNone
	}

	// Check COLORTERM for truecolor support
	colorTerm := os.Getenv("COLORTERM")
	if colorTerm == "truecolor" || colorTerm == "24bit" {
		return ColorModeTrueColor
	}

	// Check TERM for 256 color support
	termEnv := os.Getenv("TERM")
	if strings.Contains(termEnv, "256color") {
		return ColorMode256
	}

	// Default to standard ANSI colors
	if termEnv != "" && termEnv != "dumb" {
		return ColorModeStandard
	}

	return ColorModeNone
}

// SetColorMode overrides the detected color mode.
func (c *Console) SetColorMode(mode ColorMode) {
	c.colorMode = mode
}

// ColorMode returns the current color mode.
func (c *Console) ColorMode() ColorMode {
	return c.colorMode
}

// Width returns the console width in characters.
func (c *Console) Width() int {
	return c.width
}

// Height returns the console height in characters.
func (c *Console) Height() int {
	return c.height
}

// Print writes text to the console.
func (c *Console) Print(a ...interface{}) (n int, err error) {
	s := fmt.Sprint(a...)
	return c.writer.Write([]byte(s))
}

// Println writes text to the console followed by a newline.
func (c *Console) Println(a ...interface{}) (n int, err error) {
	s := fmt.Sprintln(a...)
	return c.writer.Write([]byte(s))
}

// Printf writes formatted text to the console.
func (c *Console) Printf(format string, a ...interface{}) (n int, err error) {
	s := fmt.Sprintf(format, a...)
	return c.writer.Write([]byte(s))
}

// PrintStyled writes styled text to the console.
func (c *Console) PrintStyled(st StyledText) (n int, err error) {
	seg := Segment{Text: st.Text, Style: st.Style}
	return c.PrintSegments(Segments{seg})
}

// PrintStyledln writes styled text to the console followed by a newline.
func (c *Console) PrintStyledln(st StyledText) (n int, err error) {
	n, err = c.PrintStyled(st)
	if err != nil {
		return n, err
	}
	n2, err := c.writer.Write([]byte("\n"))
	return n + n2, err
}

// PrintSegments writes segments to the console.
func (c *Console) PrintSegments(segments Segments) (n int, err error) {
	s := segments.ToANSI(c.colorMode)
	return c.writer.Write([]byte(s))
}

// PrintSegmentsln writes segments to the console followed by a newline.
func (c *Console) PrintSegmentsln(segments Segments) (n int, err error) {
	n, err = c.PrintSegments(segments)
	if err != nil {
		return n, err
	}
	n2, err := c.writer.Write([]byte("\n"))
	return n + n2, err
}

// Rule prints a horizontal rule across the console.
func (c *Console) Rule(title string) (n int, err error) {
	var segments Segments

	if title == "" {
		// Just print a line
		line := strings.Repeat("─", c.width)
		segments = Segments{{Text: line, Style: NewStyle().Dim()}}
	} else {
		// Calculate spacing
		titleLen := len(title)
		if titleLen+4 > c.width {
			// Title too long, just print it
			segments = Segments{{Text: title, Style: NewStyle().Bold()}}
		} else {
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
func (c *Console) Writer() io.Writer {
	return c.writer
}

// ANSIWriter returns an ANSI-aware writer.
func (c *Console) ANSIWriter() *ansi.Writer {
	return ansi.NewWriter(c.writer)
}

// PrintMarkup writes markup text to the console.
// Markup format: [style]text[/] where style can be:
//   - Color names: red, blue, green, etc.
//   - Hex colors: #FF0000
//   - RGB colors: rgb(255,0,0)
//   - Attributes: bold, italic, underline, etc.
//   - Background: "red on blue"
//   - Combined: "bold red on white"
func (c *Console) PrintMarkup(m string) (n int, err error) {
	return c.printMarkupInternal(m)
}

// PrintMarkupln writes markup text to the console followed by a newline.
func (c *Console) PrintMarkupln(markup string) (n int, err error) {
	n, err = c.PrintMarkup(markup)
	if err != nil {
		return n, err
	}
	n2, err := c.writer.Write([]byte("\n"))
	return n + n2, err
}
