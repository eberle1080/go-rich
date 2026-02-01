package rich

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Markup tokenizer and parser integrated into rich package to avoid import cycles.

// markupTokenType represents the type of a markup token.
type markupTokenType int

const (
	markupTokenText markupTokenType = iota
	markupTokenOpenTag
	markupTokenCloseTag
	markupTokenEOF
)

// markupToken represents a lexical token in markup.
type markupToken struct {
	typ   markupTokenType
	value string
	pos   int
}

// markupLexer tokenizes markup strings.
type markupLexer struct {
	input string
	pos   int
	start int
	width int
}

// newMarkupLexer creates a new lexer for the given input.
func newMarkupLexer(input string) *markupLexer {
	return &markupLexer{
		input: input,
	}
}

// next returns the next rune in the input.
func (l *markupLexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return 0
	}
	r, width := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = width
	l.pos += l.width
	return r
}

// backup steps back one rune.
func (l *markupLexer) backup() {
	l.pos -= l.width
}

// peek returns but does not consume the next rune.
func (l *markupLexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// emit creates a token from the current position.
func (l *markupLexer) emit(t markupTokenType) markupToken {
	token := markupToken{
		typ:   t,
		value: l.input[l.start:l.pos],
		pos:   l.start,
	}
	l.start = l.pos
	return token
}

// nextToken returns the next token from the input.
func (l *markupLexer) nextToken() markupToken {
	for {
		r := l.peek()

		if r == 0 {
			if l.pos > l.start {
				return l.emit(markupTokenText)
			}
			return markupToken{typ: markupTokenEOF, pos: l.pos}
		}

		if r == '[' {
			// Emit any pending text
			if l.pos > l.start {
				return l.emit(markupTokenText)
			}

			// Check for escaped bracket [[
			l.next() // consume '['
			if l.peek() == '[' {
				l.next() // consume second '['
				return l.emit(markupTokenText)
			}

			// Scan until ']'
			for {
				r := l.next()
				if r == 0 {
					// Unclosed tag, treat as text
					return l.emit(markupTokenText)
				}
				if r == ']' {
					break
				}
			}

			// Determine if it's a close tag
			value := l.input[l.start:l.pos]
			if strings.HasPrefix(value, "[/") {
				return l.emit(markupTokenCloseTag)
			}
			return l.emit(markupTokenOpenTag)
		}

		// Regular text
		l.next()

		// Continue until we hit a '[' or EOF
		for {
			r := l.peek()
			if r == 0 || r == '[' {
				break
			}
			l.next()
		}

		return l.emit(markupTokenText)
	}
}

// markupParser parses markup tokens into styled segments.
type markupParser struct {
	tokens     []markupToken
	pos        int
	styleStack []Style
}

// newMarkupParser creates a new parser for the given tokens.
func newMarkupParser(tokens []markupToken) *markupParser {
	return &markupParser{
		tokens:     tokens,
		styleStack: []Style{NewStyle()},
	}
}

// currentStyle returns the current style (top of stack).
func (p *markupParser) currentStyle() Style {
	if len(p.styleStack) == 0 {
		return NewStyle()
	}
	return p.styleStack[len(p.styleStack)-1]
}

// pushStyle adds a new style to the stack.
func (p *markupParser) pushStyle(style Style) {
	p.styleStack = append(p.styleStack, style)
}

// popStyle removes the top style from the stack.
func (p *markupParser) popStyle() {
	if len(p.styleStack) > 1 {
		p.styleStack = p.styleStack[:len(p.styleStack)-1]
	}
}

// parse converts tokens into styled segments.
func (p *markupParser) parse() (Segments, error) {
	var segments Segments

	for p.pos < len(p.tokens) {
		token := p.tokens[p.pos]
		p.pos++

		switch token.typ {
		case markupTokenText:
			// Handle escaped brackets [[
			text := strings.ReplaceAll(token.value, "[[", "[")
			if text != "" {
				segments = append(segments, Segment{
					Text:  text,
					Style: p.currentStyle(),
				})
			}

		case markupTokenOpenTag:
			// Parse the tag and apply style
			style, err := p.parseTag(token.value)
			if err != nil {
				// Invalid tag, treat as text
				segments = append(segments, Segment{
					Text:  token.value,
					Style: p.currentStyle(),
				})
			} else {
				p.pushStyle(style)
			}

		case markupTokenCloseTag:
			// Pop style
			p.popStyle()

		case markupTokenEOF:
			// Done
			return segments, nil
		}
	}

	return segments, nil
}

// parseTag parses a tag string like "[bold red]" into a style.
func (p *markupParser) parseTag(tag string) (Style, error) {
	// Remove [ and ]
	tag = strings.TrimPrefix(tag, "[")
	tag = strings.TrimSuffix(tag, "]")
	tag = strings.TrimSpace(tag)

	if tag == "" {
		return p.currentStyle(), nil
	}

	// Start with current style
	style := p.currentStyle()

	// Split by spaces to get individual style components
	parts := strings.Fields(tag)

	var i int
	for i < len(parts) {
		part := parts[i]

		switch strings.ToLower(part) {
		case "bold", "b":
			style = style.Bold()
		case "italic", "i":
			style = style.Italic()
		case "underline", "u":
			style = style.Underline()
		case "strikethrough", "strike", "s":
			style = style.Strikethrough()
		case "dim":
			style = style.Dim()
		case "reverse":
			style = style.Reverse()

		case "on":
			// Background color follows
			if i+1 < len(parts) {
				i++
				color, err := parseMarkupColor(parts[i])
				if err == nil {
					style = style.Background(color)
				}
			}

		default:
			// Try to parse as a color
			color, err := parseMarkupColor(part)
			if err == nil {
				style = style.Foreground(color)
			}
			// If not a valid color, ignore it
		}

		i++
	}

	return style, nil
}

// parseMarkupColor parses a color string.
func parseMarkupColor(s string) (Color, error) {
	s = strings.ToLower(s)

	// Check for hex color
	if strings.HasPrefix(s, "#") {
		return Hex(s)
	}

	// Check for rgb(r,g,b)
	if strings.HasPrefix(s, "rgb(") && strings.HasSuffix(s, ")") {
		rgb := strings.TrimPrefix(s, "rgb(")
		rgb = strings.TrimSuffix(rgb, ")")
		parts := strings.Split(rgb, ",")
		if len(parts) == 3 {
			r, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
			g, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
			b, err3 := strconv.Atoi(strings.TrimSpace(parts[2]))
			if err1 == nil && err2 == nil && err3 == nil {
				return RGB(uint8(r), uint8(g), uint8(b)), nil
			}
		}
	}

	// Check for ANSI color names
	switch s {
	case "black":
		return Black, nil
	case "red":
		return Red, nil
	case "green":
		return Green, nil
	case "yellow":
		return Yellow, nil
	case "blue":
		return Blue, nil
	case "magenta":
		return Magenta, nil
	case "cyan":
		return Cyan, nil
	case "white":
		return White, nil
	case "bright_black", "gray", "grey":
		return BrightBlack, nil
	case "bright_red":
		return BrightRed, nil
	case "bright_green":
		return BrightGreen, nil
	case "bright_yellow":
		return BrightYellow, nil
	case "bright_blue":
		return BrightBlue, nil
	case "bright_magenta":
		return BrightMagenta, nil
	case "bright_cyan":
		return BrightCyan, nil
	case "bright_white":
		return BrightWhite, nil
	}

	// Try named color
	return Named(s)
}

// parseMarkup parses markup into styled segments.
func parseMarkup(markup string) (Segments, error) {
	lexer := newMarkupLexer(markup)
	var tokens []markupToken

	for {
		token := lexer.nextToken()
		tokens = append(tokens, token)
		if token.typ == markupTokenEOF {
			break
		}
	}

	parser := newMarkupParser(tokens)
	return parser.parse()
}

// printMarkupInternal is the internal implementation of PrintMarkup.
func (c *Console) printMarkupInternal(m string) (n int, err error) {
	segments, err := parseMarkup(m)
	if err != nil {
		// On error, just print the raw markup
		return c.writer.Write([]byte(m))
	}
	return c.PrintSegments(segments)
}

// StripMarkup removes all markup tags from a string.
func StripMarkup(markup string) string {
	lexer := newMarkupLexer(markup)
	var result strings.Builder

	for {
		token := lexer.nextToken()
		if token.typ == markupTokenEOF {
			break
		}
		if token.typ == markupTokenText {
			result.WriteString(strings.ReplaceAll(token.value, "[[", "["))
		}
	}

	return result.String()
}

// EscapeMarkup escapes markup by doubling brackets.
func EscapeMarkup(s string) string {
	return strings.ReplaceAll(s, "[", "[[")
}

// ValidateMarkup checks if markup is well-formed.
func ValidateMarkup(markup string) error {
	lexer := newMarkupLexer(markup)
	depth := 0

	for {
		token := lexer.nextToken()
		if token.typ == markupTokenEOF {
			break
		}

		switch token.typ {
		case markupTokenOpenTag:
			depth++
		case markupTokenCloseTag:
			depth--
			if depth < 0 {
				return fmt.Errorf("unmatched close tag at position %d", token.pos)
			}
		}
	}

	if depth > 0 {
		return fmt.Errorf("unclosed tags: %d tags remain open", depth)
	}

	return nil
}
