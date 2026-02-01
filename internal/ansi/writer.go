package ansi

import (
	"io"
	"regexp"
)

// Writer wraps an io.Writer and provides ANSI-aware writing.
type Writer struct {
	w io.Writer
}

// NewWriter creates a new ANSI writer.
func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

// Write implements io.Writer.
func (w *Writer) Write(p []byte) (n int, err error) {
	return w.w.Write(p)
}

// WriteString writes a string.
func (w *Writer) WriteString(s string) (n int, err error) {
	return io.WriteString(w.w, s)
}

// StripANSI removes all ANSI escape sequences from a string.
func StripANSI(s string) string {
	// Regex to match ANSI escape sequences
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return ansiRegex.ReplaceAllString(s, "")
}

// Length returns the visible length of a string (excluding ANSI codes).
func Length(s string) int {
	return len(StripANSI(s))
}
