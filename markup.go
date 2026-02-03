package rich

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Markup tokenizer and parser integrated into rich package to avoid import cycles.
//
// This file implements a lexer and parser for the markup language, which allows
// inline styling using tags like [bold]text[/]. The implementation follows a
// two-phase approach:
//  1. Lexer (markupLexer): Tokenizes the input string into tokens
//  2. Parser (markupParser): Converts tokens into styled segments
//
// The markup language supports:
//   - Style tags: [bold], [italic], [underline], etc.
//   - Color tags: [red], [#FF0000], [rgb(255,0,0)]
//   - Background colors: [red on blue]
//   - Combined styles: [bold red on white]
//   - Close tags: [/]
//   - Escaped brackets: [[

// markupTokenType represents the type of a markup token.
// The lexer categorizes input into these token types during tokenization.
type markupTokenType int

const (
	// markupTokenText represents plain text content (not a tag).
	// This includes regular text and escaped brackets ([[ → [).
	markupTokenText markupTokenType = iota

	// markupTokenOpenTag represents an opening style tag like [bold] or [red].
	// The tag content (text between [ and ]) is stored in the token value.
	markupTokenOpenTag

	// markupTokenCloseTag represents a closing tag [/].
	// This pops a style from the parser's style stack.
	markupTokenCloseTag

	// markupTokenEOF represents the end of the input.
	// The lexer emits this token when all input has been consumed.
	markupTokenEOF
)

// markupToken represents a lexical token in markup.
// Each token has a type, the matched text, and its position in the input.
type markupToken struct {
	typ   markupTokenType // Type of token (text, open tag, close tag, EOF)
	value string          // The actual text matched by this token
	pos   int             // Position in the input where this token starts
}

// markupLexer tokenizes markup strings into tokens.
// The lexer scans through the input character by character, identifying
// brackets and building tokens. It handles:
//   - Regular text runs (everything outside brackets)
//   - Opening tags ([style])
//   - Closing tags ([/])
//   - Escaped brackets ([[)
//   - Unclosed tags (treated as text)
type markupLexer struct {
	input string // The complete input string being lexed
	pos   int    // Current position in input (next rune to read)
	start int    // Start position of current token being built
	width int    // Width of last rune read (for backup)
}

// newMarkupLexer creates a new lexer for the given input.
// The lexer starts at position 0 and is ready to emit tokens.
func newMarkupLexer(input string) *markupLexer {
	return &markupLexer{
		input: input,
	}
}

// next returns the next rune in the input and advances the position.
// Returns 0 (null rune) when at EOF.
// The width field is updated to store the byte width of the rune for backup.
func (l *markupLexer) next() rune {
	// Check for EOF
	if l.pos >= len(l.input) {
		l.width = 0
		return 0
	}

	// Decode the next UTF-8 rune
	r, width := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = width
	l.pos += l.width
	return r
}

// backup steps back one rune.
// Can only back up one rune (the last one read with next).
func (l *markupLexer) backup() {
	l.pos -= l.width
}

// peek returns but does not consume the next rune.
// This is used for lookahead when deciding how to tokenize.
// Returns 0 (null rune) at EOF.
func (l *markupLexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// emit creates a token from the accumulated characters.
// The token spans from l.start to l.pos in the input.
// After emitting, l.start is moved to l.pos (ready for next token).
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
// This is the main lexer loop that categorizes characters into tokens.
//
// Logic flow:
//  1. If at EOF and no pending text, emit EOF token
//  2. If at '[':
//     a. Emit any pending text first
//     b. Check for escaped bracket [[ → emit as text token
//     c. Scan until ']' to get full tag
//     d. Classify as open tag or close tag based on content
//  3. Otherwise, accumulate text until next '[' or EOF
//
// The lexer handles these cases:
//   - Plain text: "hello" → TEXT("hello")
//   - Tags: "[bold]" → OPEN("[bold]")
//   - Close tags: "[/]" → CLOSE("[/]")
//   - Escaped brackets: "[[text]]" → TEXT("[[") + TEXT("text") + TEXT("]]")
//   - Unclosed tags: "[bold" → TEXT("[bold") (treated as text)
func (l *markupLexer) nextToken() markupToken {
	for {
		r := l.peek()

		// Handle EOF
		if r == 0 {
			// Emit pending text if any
			if l.pos > l.start {
				return l.emit(markupTokenText)
			}
			// Otherwise emit EOF token
			return markupToken{typ: markupTokenEOF, pos: l.pos}
		}

		// Handle opening bracket
		if r == '[' {
			// First, emit any pending text before this bracket
			if l.pos > l.start {
				return l.emit(markupTokenText)
			}

			// Consume the opening bracket
			l.next()

			// Check for escaped bracket [[
			if l.peek() == '[' {
				l.next() // consume second '['
				// Emit as text (will be converted to single [ by parser)
				return l.emit(markupTokenText)
			}

			// Scan until ']' to get the complete tag
			for {
				r := l.next()
				if r == 0 {
					// Unclosed tag (reached EOF without ']')
					// Treat the entire thing as text
					return l.emit(markupTokenText)
				}
				if r == ']' {
					// Found closing bracket, tag is complete
					break
				}
			}

			// Determine tag type based on content
			value := l.input[l.start:l.pos]
			if strings.HasPrefix(value, "[/") {
				// Close tag: [/] or [/style]
				return l.emit(markupTokenCloseTag)
			}
			// Open tag: [style]
			return l.emit(markupTokenOpenTag)
		}

		// Regular text: consume characters until '[' or EOF
		l.next()

		// Continue consuming text
		for {
			r := l.peek()
			if r == 0 || r == '[' {
				// Stop at EOF or next bracket
				break
			}
			l.next()
		}

		// Emit the accumulated text
		return l.emit(markupTokenText)
	}
}

// markupParser parses markup tokens into styled segments.
// The parser maintains a style stack to handle nested tags.
// When an opening tag is encountered, a new style is pushed onto the stack.
// When a closing tag is encountered, a style is popped from the stack.
//
// The style stack starts with one empty style (no formatting), so closing
// all tags returns to the default style rather than causing an error.
//
// Example transformation:
//
//	Input: "[bold]Hello [red]world[/][/]"
//	Tokens: OPEN("[bold]"), TEXT("Hello "), OPEN("[red]"), TEXT("world"), CLOSE("[/]"), CLOSE("[/]")
//	Segments: {"Hello ", bold}, {"world", bold+red}
type markupParser struct {
	tokens     []markupToken // All tokens to parse
	pos        int           // Current position in tokens array
	styleStack []Style       // Stack of active styles (innermost at end)
}

// newMarkupParser creates a new parser for the given tokens.
// The style stack is initialized with one empty style as the base.
func newMarkupParser(tokens []markupToken) *markupParser {
	return &markupParser{
		tokens:     tokens,
		styleStack: []Style{NewStyle()}, // Start with base (unstyled) style
	}
}

// currentStyle returns the current style (top of stack).
// This is the style that will be applied to the next text segment.
// The current style is the result of all active (unclosed) tags.
func (p *markupParser) currentStyle() Style {
	if len(p.styleStack) == 0 {
		// Shouldn't happen, but provide a safe default
		return NewStyle()
	}
	return p.styleStack[len(p.styleStack)-1]
}

// pushStyle adds a new style to the stack.
// This happens when an opening tag is encountered.
// The new style is based on the current style with additional attributes.
func (p *markupParser) pushStyle(style Style) {
	p.styleStack = append(p.styleStack, style)
}

// popStyle removes the top style from the stack.
// This happens when a closing tag [/] is encountered.
// The base style is never popped (stack size stays >= 1).
func (p *markupParser) popStyle() {
	if len(p.styleStack) > 1 {
		p.styleStack = p.styleStack[:len(p.styleStack)-1]
	}
	// If stack is already at base (len==1), do nothing
}

// parse converts tokens into styled segments.
// Processes each token in sequence, building up segments with appropriate styles.
//
// Token handling:
//   - TEXT: Create a segment with current style, handle escaped brackets
//   - OPEN TAG: Parse the style and push it onto the stack
//   - CLOSE TAG: Pop a style from the stack
//   - EOF: Return completed segments
//
// Invalid tags (parse errors) are treated as literal text rather than
// causing the entire parse to fail. This provides graceful degradation.
//
// Example:
//
//	Input tokens: OPEN("[bold]"), TEXT("Hi"), CLOSE("[/]")
//	Output: Segments{{Text:"Hi", Style:bold}}
func (p *markupParser) parse() (Segments, error) {
	var segments Segments

	// Process each token
	for p.pos < len(p.tokens) {
		token := p.tokens[p.pos]
		p.pos++

		switch token.typ {
		case markupTokenText:
			// Convert escaped brackets [[ → [
			text := strings.ReplaceAll(token.value, "[[", "[")

			// Create a segment with the current style
			if text != "" {
				segments = append(segments, Segment{
					Text:  text,
					Style: p.currentStyle(),
				})
			}

		case markupTokenOpenTag:
			// Parse the tag to extract style attributes
			style, err := p.parseTag(token.value)
			if err != nil {
				// Invalid tag syntax, treat it as literal text
				// This allows graceful handling of malformed markup
				segments = append(segments, Segment{
					Text:  token.value,
					Style: p.currentStyle(),
				})
			} else {
				// Valid tag, push the style onto the stack
				p.pushStyle(style)
			}

		case markupTokenCloseTag:
			// Close the most recent tag
			p.popStyle()

		case markupTokenEOF:
			// End of input, return what we've parsed
			return segments, nil
		}
	}

	return segments, nil
}

// parseTag parses a tag string like "[bold red]" into a style.
// The tag content (without brackets) is split into space-separated parts,
// and each part is interpreted as a style attribute or color.
//
// Supported tag components:
//   - Attributes: bold, italic, underline, strikethrough, dim, reverse
//   - Foreground colors: red, #FF0000, rgb(255,0,0)
//   - Background colors: on blue, on #0000FF
//   - Combinations: "bold red on blue"
//
// The resulting style is based on the current style with new attributes added.
// This allows tags to accumulate styles: [bold][red]text[/][/] applies both.
//
// Invalid components are silently ignored rather than causing errors.
//
// Example:
//
//	"[bold red on white]" → bold + foreground:red + background:white
func (p *markupParser) parseTag(tag string) (Style, error) {
	// Remove brackets and whitespace
	tag = strings.TrimPrefix(tag, "[")
	tag = strings.TrimSuffix(tag, "]")
	tag = strings.TrimSpace(tag)

	// Empty tag maintains current style
	if tag == "" {
		return p.currentStyle(), nil
	}

	// Start with current style (inherit from outer tags)
	style := p.currentStyle()

	// Split tag into space-separated parts
	parts := strings.Fields(tag)

	// Process each part
	var i int
	for i < len(parts) {
		part := parts[i]

		// Check for style attributes (case-insensitive)
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
			// "on" keyword indicates the next part is a background color
			// Example: "red on blue" → foreground:red, background:blue
			if i+1 < len(parts) {
				i++ // Move to the color part
				color, err := parseMarkupColor(parts[i])
				if err == nil {
					style = style.Background(color)
				}
				// If color parsing fails, silently ignore
			}

		default:
			// Try to parse as a foreground color
			color, err := parseMarkupColor(part)
			if err == nil {
				style = style.Foreground(color)
			}
			// If not a valid color or attribute, ignore it
		}

		i++
	}

	return style, nil
}

// parseMarkupColor parses a color string from markup.
// Supports multiple color formats:
//   - Hex colors: #FF0000, #00ff00 (case-insensitive)
//   - RGB function: rgb(255,0,0), rgb(0, 255, 0)
//   - ANSI color names: red, blue, green, etc.
//   - Bright colors: bright_red, bright_blue
//   - Gray/grey: both spellings accepted
//   - Named colors: orange, purple, pink (from Named function)
//
// Color names are case-insensitive.
// Returns an error if the color string doesn't match any known format.
//
// Examples:
//
//	"#FF0000" → RGBColor{255, 0, 0}
//	"rgb(255,0,0)" → RGBColor{255, 0, 0}
//	"red" → ANSIColor Red
//	"orange" → RGBColor{255, 165, 0}
func parseMarkupColor(s string) (Color, error) {
	// Normalize to lowercase for case-insensitive matching
	s = strings.ToLower(s)

	// Check for hex color (#RRGGBB format)
	if strings.HasPrefix(s, "#") {
		return Hex(s)
	}

	// Check for rgb(r,g,b) function format
	if strings.HasPrefix(s, "rgb(") && strings.HasSuffix(s, ")") {
		// Extract the content between parentheses
		rgb := strings.TrimPrefix(s, "rgb(")
		rgb = strings.TrimSuffix(rgb, ")")

		// Split into R, G, B components
		parts := strings.Split(rgb, ",")
		if len(parts) == 3 {
			// Parse each component as an integer
			r, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
			g, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
			b, err3 := strconv.Atoi(strings.TrimSpace(parts[2]))

			// All components must parse successfully
			if err1 == nil && err2 == nil && err3 == nil {
				return RGB(uint8(r), uint8(g), uint8(b)), nil
			}
		}
	}

	// Check for ANSI color names (standard 16 colors)
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

	// Bright variants (with multiple accepted names)
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

	// Fall back to named color lookup (orange, purple, pink, etc.)
	return Named(s)
}

// parseMarkup parses markup into styled segments.
// This is the main entry point for markup processing, combining lexing and parsing.
//
// Process:
//  1. Create a lexer for the input string
//  2. Tokenize the entire input into a token array
//  3. Create a parser with the tokens
//  4. Parse tokens into styled segments
//
// Returns an error only if parsing fails (currently doesn't happen as
// invalid tags are treated as text).
//
// Example:
//
//	parseMarkup("[bold]Hi[/]") → Segments{{Text:"Hi", Style:bold}}
func parseMarkup(markup string) (Segments, error) {
	// Phase 1: Lexing - convert string to tokens
	lexer := newMarkupLexer(markup)
	var tokens []markupToken

	for {
		token := lexer.nextToken()
		tokens = append(tokens, token)
		if token.typ == markupTokenEOF {
			break
		}
	}

	// Phase 2: Parsing - convert tokens to styled segments
	parser := newMarkupParser(tokens)
	return parser.parse()
}

// printMarkupInternal is the internal implementation of PrintMarkup.
// Parses the markup into segments and writes them to the console.
// If parsing fails, falls back to printing the raw markup as plain text.
func (c *Console) printMarkupInternal(m string) (n int, err error) {
	segments, err := parseMarkup(m)
	if err != nil {
		// On error, print raw markup without styling
		return c.writer.Write([]byte(m))
	}
	return c.PrintSegments(segments)
}

// StripMarkup removes all markup tags from a string, leaving only the text content.
// This is useful for:
//   - Measuring the display length of markup text
//   - Extracting plain text for logging or storage
//   - Generating non-styled versions of marked-up content
//
// Escaped brackets [[ are correctly converted to single brackets [.
//
// Example:
//
//	StripMarkup("[bold]Hello[/] [red]world[/]") → "Hello world"
//	StripMarkup("[[not a tag]]") → "[not a tag]"
func StripMarkup(markup string) string {
	lexer := newMarkupLexer(markup)
	var result strings.Builder

	for {
		token := lexer.nextToken()
		if token.typ == markupTokenEOF {
			break
		}

		// Only keep text tokens, discard tag tokens
		if token.typ == markupTokenText {
			// Convert escaped brackets [[ → [
			result.WriteString(strings.ReplaceAll(token.value, "[[", "["))
		}
	}

	return result.String()
}

// EscapeMarkup escapes markup by doubling all brackets.
// This makes literal bracket characters safe to use in markup.
//
// Use this when you want to display text that contains [ or ] characters
// without them being interpreted as markup tags.
//
// Example:
//
//	EscapeMarkup("Use [bold] for bold") → "Use [[bold]] for bold"
//	// When printed with markup: "Use [bold] for bold" (literal, not bold)
func EscapeMarkup(s string) string {
	// Replace every [ with [[
	return strings.ReplaceAll(s, "[", "[[")
}

// ValidateMarkup checks if markup is well-formed.
// Well-formed markup has balanced opening and closing tags:
//   - Every [tag] should have a corresponding [/]
//   - Close tags should not appear before their matching open tags
//
// Returns nil if the markup is valid, or an error describing the problem.
//
// Note: This checks tag balance but doesn't validate tag content.
// Invalid color names or attributes will be ignored during rendering,
// not caught by this validation.
//
// Example:
//
//	ValidateMarkup("[bold]text[/]")        → nil (valid)
//	ValidateMarkup("[bold]text")           → error (unclosed tag)
//	ValidateMarkup("text[/]")              → error (unmatched close)
//	ValidateMarkup("[bold][red]text[/]")   → error (one tag unclosed)
func ValidateMarkup(markup string) error {
	lexer := newMarkupLexer(markup)
	depth := 0 // Track nesting depth

	for {
		token := lexer.nextToken()
		if token.typ == markupTokenEOF {
			break
		}

		switch token.typ {
		case markupTokenOpenTag:
			// Opening tag increases depth
			depth++

		case markupTokenCloseTag:
			// Closing tag decreases depth
			depth--

			// Check for unmatched close tag
			if depth < 0 {
				return fmt.Errorf("unmatched close tag at position %d", token.pos)
			}
		}
	}

	// Check for unclosed tags
	if depth > 0 {
		return fmt.Errorf("unclosed tags: %d tags remain open", depth)
	}

	return nil
}
