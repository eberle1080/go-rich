package ansi

import (
	"io"
	"regexp"
)

// Writer wraps an io.Writer and provides ANSI-aware writing.
// Currently, this is a thin wrapper that passes through to the underlying writer.
// It provides a foundation for future ANSI-aware operations like:
//   - Automatic escape sequence stripping
//   - Width calculation
//   - Buffer management
type Writer struct {
	w io.Writer // Underlying writer
}

// NewWriter creates a new ANSI writer wrapping the given io.Writer.
//
// Example:
//
//	w := ansi.NewWriter(os.Stdout)
//	w.WriteString(ansi.Bold + "Bold text" + ansi.Reset)
func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

// Write implements io.Writer.
// Writes the byte slice to the underlying writer.
// Returns the number of bytes written and any error encountered.
func (w *Writer) Write(p []byte) (n int, err error) {
	return w.w.Write(p)
}

// WriteString writes a string to the underlying writer.
// This is more efficient than Write for string data as it avoids
// converting the string to a byte slice.
// Returns the number of bytes written and any error encountered.
func (w *Writer) WriteString(s string) (n int, err error) {
	return io.WriteString(w.w, s)
}

// StripANSI removes all ANSI escape sequences from a string.
// This is useful for:
//   - Calculating the visible length of styled text
//   - Saving styled output to plain text files
//   - Comparing text content without considering styling
//
// The function removes standard SGR (Select Graphic Rendition) sequences,
// cursor movement codes, and other CSI (Control Sequence Introducer) sequences.
//
// Example:
//
//	styled := "\x1b[1mBold\x1b[0m text"
//	plain := ansi.StripANSI(styled)
//	// plain == "Bold text"
func StripANSI(s string) string {
	// Regular expression to match ANSI escape sequences:
	// \x1b\[       - ESC[ (start of CSI sequence)
	// [0-9;]*      - Zero or more digits or semicolons (parameters)
	// [a-zA-Z]     - Letter (command character)
	//
	// This matches most common ANSI sequences including:
	// - SGR codes: \x1b[1m, \x1b[31m, \x1b[0m
	// - Cursor movement: \x1b[A, \x1b[2J
	// - Colors: \x1b[38;2;255;0;0m
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return ansiRegex.ReplaceAllString(s, "")
}

// Length returns the visible length of a string (excluding ANSI codes).
// This counts only the visible characters, not the escape sequence bytes.
//
// Note: This returns byte length after stripping ANSI codes, not rune count.
// For accurate character counting with multi-byte Unicode, use utf8.RuneCountInString
// on the stripped result.
//
// Example:
//
//	styled := "\x1b[1mHello\x1b[0m"
//	length := ansi.Length(styled)
//	// length == 5 (counts only "Hello")
func Length(s string) int {
	// Strip ANSI codes then measure length
	return len(StripANSI(s))
}
