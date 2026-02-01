package rich

import (
	"strings"
	"unicode/utf8"
)

// Segment represents an atomic unit of styled text.
type Segment struct {
	Text  string
	Style Style
}

// Segments is a slice of segments with helper methods.
type Segments []Segment

// String converts segments to a plain text string (no styling).
func (s Segments) String() string {
	var b strings.Builder
	for _, seg := range s {
		b.WriteString(seg.Text)
	}
	return b.String()
}

// Length returns the total visible character length of all segments.
func (s Segments) Length() int {
	length := 0
	for _, seg := range s {
		length += utf8.RuneCountInString(seg.Text)
	}
	return length
}

// ToANSI converts segments to an ANSI-escaped string.
func (s Segments) ToANSI(mode ColorMode) string {
	if mode == ColorModeNone {
		return s.String()
	}

	var b strings.Builder
	for _, seg := range s {
		ansi := seg.Style.toANSI(mode)
		if ansi != "" {
			b.WriteString(ansi)
		}
		b.WriteString(seg.Text)
		if ansi != "" {
			b.WriteString("\x1b[0m") // Reset
		}
	}
	return b.String()
}

// Append adds segments to the end.
func (s Segments) Append(segs ...Segment) Segments {
	return append(s, segs...)
}

// Join concatenates multiple segment slices.
func Join(segments ...Segments) Segments {
	total := 0
	for _, s := range segments {
		total += len(s)
	}

	result := make(Segments, 0, total)
	for _, s := range segments {
		result = append(result, s...)
	}
	return result
}
