// Package ansi provides ANSI escape sequence constants and utilities.
//
// This internal package contains the low-level ANSI escape codes used for
// terminal control and formatting. It is used internally by the rich package
// to generate styled output.
//
// # ANSI Escape Sequences
//
// ANSI escape sequences are special character sequences that control terminal
// behavior, such as text formatting, cursor movement, and screen clearing.
// They all start with the ESC character (0x1B or \x1b) followed by specific codes.
//
// # Text Formatting
//
// The package provides constants for common text attributes:
//
//	ansi.Bold          // Make text bold
//	ansi.Dim           // Make text dim/faint
//	ansi.Italic        // Make text italic
//	ansi.Underline     // Underline text
//	ansi.Strikethrough // Strike through text
//	ansi.Reverse       // Swap foreground and background colors
//
// # Cursor Control
//
// Constants for cursor movement:
//
//	ansi.CursorUp       // Move cursor up one line
//	ansi.CursorDown     // Move cursor down one line
//	ansi.CursorForward  // Move cursor forward (right)
//	ansi.CursorBack     // Move cursor backward (left)
//
// # Screen Control
//
// Constants for screen manipulation:
//
//	ansi.ClearScreen     // Clear entire screen
//	ansi.ClearLine       // Clear entire line
//	ansi.HideCursor      // Hide the cursor
//	ansi.ShowCursor      // Show the cursor
//
// # Writer
//
// The Writer type wraps an io.Writer and provides ANSI-aware writing utilities:
//
//	w := ansi.NewWriter(os.Stdout)
//	w.WriteString(ansi.Bold + "Bold text" + ansi.Reset)
//
// # Utilities
//
// Helper functions for working with ANSI sequences:
//
//	// Remove all ANSI escape sequences from a string
//	plain := ansi.StripANSI(styledText)
//
//	// Get the visible length of a string (excluding ANSI codes)
//	length := ansi.Length(styledText)
//
// # Note
//
// This is an internal package and should not be imported directly by users
// of the go-rich library. Use the higher-level APIs in the rich package instead.
package ansi
