package ansi

// ANSI escape sequences and control codes.
//
// ANSI escape sequences are special character combinations that control
// terminal behavior. They all start with ESC (0x1B, written as \x1b in Go),
// followed by specific control codes.
//
// These constants provide ready-to-use escape sequences for common operations.
// For more complex operations (like colored text), use the higher-level
// APIs in the rich package instead.
const (
	// Reset clears all text attributes and colors, returning to default terminal style.
	// SGR (Select Graphic Rendition) code: 0
	// Use this after any styled text to ensure following text is unstyled.
	Reset = "\x1b[0m"

	// Text attributes - modify the appearance of text

	// Bold makes text bold or bright. May also brighten foreground colors on some terminals.
	// SGR code: 1
	Bold = "\x1b[1m"

	// Dim makes text faint or decreased intensity.
	// SGR code: 2
	Dim = "\x1b[2m"

	// Italic makes text italic. Not all terminals support this; some show inverse or underline instead.
	// SGR code: 3
	Italic = "\x1b[3m"

	// Underline draws a line under the text.
	// SGR code: 4
	Underline = "\x1b[4m"

	// Blink makes text blink. Rarely used and may not be supported on modern terminals.
	// SGR code: 5
	Blink = "\x1b[5m"

	// Reverse swaps foreground and background colors (inverse video).
	// SGR code: 7
	Reverse = "\x1b[7m"

	// Hidden makes text invisible (same color as background). Useful for passwords.
	// SGR code: 8
	Hidden = "\x1b[8m"

	// Strikethrough draws a line through the middle of text.
	// SGR code: 9
	Strikethrough = "\x1b[9m"

	// Reset specific attributes - turn off individual attributes without affecting others

	// ResetBold turns off bold/bright without affecting other attributes.
	// SGR code: 22 (normal intensity)
	ResetBold = "\x1b[22m"

	// ResetDim turns off dim/faint without affecting other attributes.
	// SGR code: 22 (normal intensity) - same as ResetBold
	ResetDim = "\x1b[22m"

	// ResetItalic turns off italic without affecting other attributes.
	// SGR code: 23
	ResetItalic = "\x1b[23m"

	// ResetUnderline turns off underline without affecting other attributes.
	// SGR code: 24
	ResetUnderline = "\x1b[24m"

	// ResetBlink turns off blink without affecting other attributes.
	// SGR code: 25
	ResetBlink = "\x1b[25m"

	// ResetReverse turns off reverse video without affecting other attributes.
	// SGR code: 27
	ResetReverse = "\x1b[27m"

	// ResetHidden reveals hidden text without affecting other attributes.
	// SGR code: 28
	ResetHidden = "\x1b[28m"

	// ResetStrikethrough turns off strikethrough without affecting other attributes.
	// SGR code: 29
	ResetStrikethrough = "\x1b[29m"

	// Cursor control - move the cursor position

	// CursorUp moves the cursor up one line.
	// CSI code: A
	CursorUp = "\x1b[A"

	// CursorDown moves the cursor down one line.
	// CSI code: B
	CursorDown = "\x1b[B"

	// CursorForward moves the cursor forward (right) one column.
	// CSI code: C
	CursorForward = "\x1b[C"

	// CursorBack moves the cursor backward (left) one column.
	// CSI code: D
	CursorBack = "\x1b[D"

	// CursorNextLine moves the cursor to the beginning of the next line.
	// CSI code: E
	CursorNextLine = "\x1b[E"

	// CursorPrevLine moves the cursor to the beginning of the previous line.
	// CSI code: F
	CursorPrevLine = "\x1b[F"

	// Screen control - manipulate screen content

	// ClearScreen clears the entire screen.
	// CSI code: 2J
	ClearScreen = "\x1b[2J"

	// ClearLine clears the entire current line.
	// CSI code: 2K
	ClearLine = "\x1b[2K"

	// ClearLineToEnd clears from cursor to end of line.
	// CSI code: K (equivalent to 0K)
	ClearLineToEnd = "\x1b[K"

	// ClearLineToStart clears from beginning of line to cursor.
	// CSI code: 1K
	ClearLineToStart = "\x1b[1K"

	// Cursor visibility control

	// HideCursor makes the cursor invisible.
	// DEC Private Mode: ?25l
	HideCursor = "\x1b[?25l"

	// ShowCursor makes the cursor visible.
	// DEC Private Mode: ?25h
	ShowCursor = "\x1b[?25h"
)
