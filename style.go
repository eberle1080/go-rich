package rich

// Style represents an immutable text style with colors and formatting attributes.
// Styles are created using a fluent builder pattern, where each method returns
// a new Style with the specified attribute enabled. This allows for easy chaining:
//
//	style := NewStyle().Bold().Foreground(Red).Background(White)
//
// Because styles are immutable, they can be safely reused and shared without
// concerns about accidental modification.
//
// All style attributes are optional. An empty style (created with NewStyle())
// renders text without any formatting.
type Style struct {
	fg            Color // Foreground (text) color
	bg            Color // Background color
	bold          bool  // Bold/bright text (SGR 1)
	italic        bool  // Italic text (SGR 3)
	underline     bool  // Underlined text (SGR 4)
	strikethrough bool  // Strikethrough text (SGR 9)
	dim           bool  // Dim/faint text (SGR 2)
	reverse       bool  // Reverse video - swap fg/bg colors (SGR 7)
}

// NewStyle creates a new empty style with no formatting.
// This is the starting point for building styled text using the fluent API.
//
// Example:
//
//	style := NewStyle().Bold().Foreground(Red)
func NewStyle() Style {
	return Style{}
}

// Foreground returns a new style with the specified foreground (text) color.
// The color will be rendered according to the terminal's color mode capabilities.
//
// Example:
//
//	style := NewStyle().Foreground(rich.Red)
//	style := NewStyle().Foreground(rich.RGB(255, 100, 50))
func (s Style) Foreground(color Color) Style {
	s.fg = color
	return s
}

// Background returns a new style with the specified background color.
// The color will be rendered according to the terminal's color mode capabilities.
//
// Example:
//
//	style := NewStyle().Background(rich.Blue)
//	style := NewStyle().Foreground(rich.White).Background(rich.Red)
func (s Style) Background(color Color) Style {
	s.bg = color
	return s
}

// Bold returns a new style with bold (increased intensity) enabled.
// Uses ANSI SGR code 1. On some terminals, this also brightens colors.
//
// Example:
//
//	style := NewStyle().Bold()
func (s Style) Bold() Style {
	s.bold = true
	return s
}

// Italic returns a new style with italic text enabled.
// Uses ANSI SGR code 3. Not all terminals support italic rendering;
// some may display italic text as inverse or underlined instead.
//
// Example:
//
//	style := NewStyle().Italic()
func (s Style) Italic() Style {
	s.italic = true
	return s
}

// Underline returns a new style with underline enabled.
// Uses ANSI SGR code 4. Draws a line under the text.
//
// Example:
//
//	style := NewStyle().Underline()
func (s Style) Underline() Style {
	s.underline = true
	return s
}

// Strikethrough returns a new style with strikethrough enabled.
// Uses ANSI SGR code 9. Draws a line through the middle of the text.
// Not all terminals support this attribute.
//
// Example:
//
//	style := NewStyle().Strikethrough()
func (s Style) Strikethrough() Style {
	s.strikethrough = true
	return s
}

// Dim returns a new style with dim/faint text enabled.
// Uses ANSI SGR code 2. Renders text with decreased intensity.
// Useful for de-emphasizing secondary information.
//
// Example:
//
//	style := NewStyle().Dim()
func (s Style) Dim() Style {
	s.dim = true
	return s
}

// Reverse returns a new style with reverse video enabled.
// Uses ANSI SGR code 7. Swaps the foreground and background colors.
// Useful for highlighting or creating inverse text effects.
//
// Example:
//
//	style := NewStyle().Reverse()
func (s Style) Reverse() Style {
	s.reverse = true
	return s
}

// Render applies this style to the given text, creating StyledText.
// This is a convenience method for creating styled text that can be
// printed using Console.PrintStyled or Console.PrintStyledln.
//
// Example:
//
//	style := NewStyle().Bold().Foreground(Red)
//	styled := style.Render("Error!")
//	console.PrintStyledln(styled)
func (s Style) Render(text string) StyledText {
	return StyledText{
		Text:  text,
		Style: s,
	}
}

// toANSI generates the ANSI escape sequence for this style.
// Returns an empty string if the color mode is ColorModeNone.
//
// The method builds a combined escape sequence by:
//  1. Collecting all text attribute codes (bold, dim, italic, etc.) into a single ESC[...m sequence
//  2. Appending the foreground color sequence (if set)
//  3. Appending the background color sequence (if set)
//
// This approach minimizes the number of escape sequences while maintaining compatibility.
//
// ANSI SGR (Select Graphic Rendition) codes used:
//   - 1: Bold/bright
//   - 2: Dim/faint
//   - 3: Italic
//   - 4: Underline
//   - 7: Reverse video
//   - 9: Strikethrough
func (s Style) toANSI(mode ColorMode) string {
	// No styling in ColorModeNone
	if mode == ColorModeNone {
		return ""
	}

	var codes []string

	// Collect all enabled text attributes
	if s.bold {
		codes = append(codes, "1") // SGR 1: Bold
	}
	if s.dim {
		codes = append(codes, "2") // SGR 2: Dim
	}
	if s.italic {
		codes = append(codes, "3") // SGR 3: Italic
	}
	if s.underline {
		codes = append(codes, "4") // SGR 4: Underline
	}
	if s.reverse {
		codes = append(codes, "7") // SGR 7: Reverse video
	}
	if s.strikethrough {
		codes = append(codes, "9") // SGR 9: Strikethrough
	}

	// Build the attribute sequence: ESC[code1;code2;...m
	seq := ""
	if len(codes) > 0 {
		seq = "\x1b[" + codes[0]
		for _, code := range codes[1:] {
			seq += ";" + code
		}
		seq += "m"
	}

	// Append foreground color sequence (if set)
	if s.fg != nil {
		seq += s.fg.toANSI(mode, true)
	}

	// Append background color sequence (if set)
	if s.bg != nil {
		seq += s.bg.toANSI(mode, false)
	}

	return seq
}

// StyledText represents text with an associated style.
// This is typically created using Style.Render() and can be printed
// using Console.PrintStyled() or Console.PrintStyledln().
//
// Example:
//
//	styled := NewStyle().Bold().Render("Hello")
//	console.PrintStyledln(styled)
type StyledText struct {
	Text  string // The text content to display
	Style Style  // The style to apply to the text
}
