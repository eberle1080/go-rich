package rich

// Style represents an immutable text style with colors and formatting attributes.
type Style struct {
	fg         Color
	bg         Color
	bold       bool
	italic     bool
	underline  bool
	strikethrough bool
	dim        bool
	reverse    bool
}

// NewStyle creates a new empty style.
func NewStyle() Style {
	return Style{}
}

// Foreground returns a new style with the specified foreground color.
func (s Style) Foreground(color Color) Style {
	s.fg = color
	return s
}

// Background returns a new style with the specified background color.
func (s Style) Background(color Color) Style {
	s.bg = color
	return s
}

// Bold returns a new style with bold enabled.
func (s Style) Bold() Style {
	s.bold = true
	return s
}

// Italic returns a new style with italic enabled.
func (s Style) Italic() Style {
	s.italic = true
	return s
}

// Underline returns a new style with underline enabled.
func (s Style) Underline() Style {
	s.underline = true
	return s
}

// Strikethrough returns a new style with strikethrough enabled.
func (s Style) Strikethrough() Style {
	s.strikethrough = true
	return s
}

// Dim returns a new style with dim/faint enabled.
func (s Style) Dim() Style {
	s.dim = true
	return s
}

// Reverse returns a new style with reverse video enabled.
func (s Style) Reverse() Style {
	s.reverse = true
	return s
}

// Render applies this style to the given text.
func (s Style) Render(text string) StyledText {
	return StyledText{
		Text:  text,
		Style: s,
	}
}

// toANSI generates the ANSI escape sequence for this style.
func (s Style) toANSI(mode ColorMode) string {
	if mode == ColorModeNone {
		return ""
	}

	var codes []string

	if s.bold {
		codes = append(codes, "1")
	}
	if s.dim {
		codes = append(codes, "2")
	}
	if s.italic {
		codes = append(codes, "3")
	}
	if s.underline {
		codes = append(codes, "4")
	}
	if s.reverse {
		codes = append(codes, "7")
	}
	if s.strikethrough {
		codes = append(codes, "9")
	}

	seq := ""
	if len(codes) > 0 {
		seq = "\x1b[" + codes[0]
		for _, code := range codes[1:] {
			seq += ";" + code
		}
		seq += "m"
	}

	if s.fg != nil {
		seq += s.fg.toANSI(mode, true)
	}

	if s.bg != nil {
		seq += s.bg.toANSI(mode, false)
	}

	return seq
}

// StyledText represents text with an associated style.
type StyledText struct {
	Text  string
	Style Style
}
