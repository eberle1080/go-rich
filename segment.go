package rich

import (
	"strings"
	"unicode/utf8"
)

// Segment represents an atomic unit of styled text.
// A segment is the fundamental building block of rich text rendering.
// Each segment contains a string of text and an associated style.
//
// Segments are typically grouped into Segments slices, which can be
// rendered together to create complex styled output.
type Segment struct {
	Text  string // The text content of this segment
	Style Style  // The style to apply to this segment
}

// Segments is a slice of segments with helper methods.
// This type provides convenient operations for working with collections
// of segments, such as converting to strings, calculating lengths,
// and rendering to ANSI escape sequences.
type Segments []Segment

// String converts segments to a plain text string (no styling).
// All ANSI styling is stripped, returning only the text content.
// This is useful for measuring text length or saving to plain text files.
//
// Example:
//
//	segments := Segments{
//		{Text: "Hello ", Style: NewStyle().Bold()},
//		{Text: "world", Style: NewStyle().Foreground(Red)},
//	}
//	plain := segments.String() // Returns: "Hello world"
func (s Segments) String() string {
	var b strings.Builder
	for _, seg := range s {
		b.WriteString(seg.Text)
	}
	return b.String()
}

// Length returns the total visible character length of all segments.
// This counts Unicode runes (characters), not bytes, so multi-byte
// characters like emojis are counted as single characters.
//
// The length excludes any ANSI escape sequences that would be added
// during rendering - it only counts the actual visible text.
//
// Example:
//
//	segments := Segments{
//		{Text: "Hello", Style: NewStyle()},
//		{Text: " 世界", Style: NewStyle()}, // Chinese for "world"
//	}
//	length := segments.Length() // Returns: 8 (5 + 1 + 2)
func (s Segments) Length() int {
	length := 0
	for _, seg := range s {
		// Count Unicode runes, not bytes
		length += utf8.RuneCountInString(seg.Text)
	}
	return length
}

// ToANSI converts segments to an ANSI-escaped string.
// Each segment's style is converted to ANSI escape sequences appropriate
// for the given color mode, with a reset sequence (ESC[0m) after each segment.
//
// If mode is ColorModeNone, returns plain text with no escape sequences
// (equivalent to calling String()).
//
// The reset sequence ensures that each segment's style doesn't bleed into
// subsequent segments, maintaining style isolation.
//
// Example:
//
//	segments := Segments{
//		{Text: "Error: ", Style: NewStyle().Bold().Foreground(Red)},
//		{Text: "File not found", Style: NewStyle()},
//	}
//	ansi := segments.ToANSI(ColorModeTrueColor)
//	// Returns: "\x1b[1m\x1b[38;2;255;0;0mError: \x1b[0mFile not found"
func (s Segments) ToANSI(mode ColorMode) string {
	// No styling in ColorModeNone, just return plain text
	if mode == ColorModeNone {
		return s.String()
	}

	var b strings.Builder
	for _, seg := range s {
		// Get the ANSI sequence for this segment's style
		ansi := seg.Style.toANSI(mode)

		// Apply the style if it has any formatting
		if ansi != "" {
			b.WriteString(ansi)
		}

		// Write the text content
		b.WriteString(seg.Text)

		// Reset formatting after this segment if we applied any
		if ansi != "" {
			b.WriteString("\x1b[0m") // SGR 0: Reset all attributes
		}
	}
	return b.String()
}

// Append adds segments to the end of this segment slice.
// Returns a new Segments slice with the additional segments appended.
//
// Example:
//
//	segs1 := Segments{{Text: "Hello", Style: NewStyle()}}
//	segs2 := segs1.Append(Segment{Text: " world", Style: NewStyle()})
func (s Segments) Append(segs ...Segment) Segments {
	return append(s, segs...)
}

// Join concatenates multiple segment slices into a single Segments slice.
// This is useful for combining segments from different sources.
//
// The function pre-allocates the result slice with the exact required capacity
// for efficiency, avoiding multiple reallocations during appending.
//
// Example:
//
//	header := Segments{{Text: "Title: ", Style: NewStyle().Bold()}}
//	content := Segments{{Text: "Hello", Style: NewStyle()}}
//	footer := Segments{{Text: "\n---", Style: NewStyle().Dim()}}
//	combined := Join(header, content, footer)
func Join(segments ...Segments) Segments {
	// Calculate total capacity needed
	total := 0
	for _, s := range segments {
		total += len(s)
	}

	// Pre-allocate result slice with exact capacity
	result := make(Segments, 0, total)

	// Append all segments
	for _, s := range segments {
		result = append(result, s...)
	}

	return result
}
