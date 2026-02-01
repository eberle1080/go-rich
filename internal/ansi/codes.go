package ansi

// ANSI escape sequences and control codes.
const (
	Reset = "\x1b[0m"

	// Text attributes
	Bold          = "\x1b[1m"
	Dim           = "\x1b[2m"
	Italic        = "\x1b[3m"
	Underline     = "\x1b[4m"
	Blink         = "\x1b[5m"
	Reverse       = "\x1b[7m"
	Hidden        = "\x1b[8m"
	Strikethrough = "\x1b[9m"

	// Reset specific attributes
	ResetBold          = "\x1b[22m"
	ResetDim           = "\x1b[22m"
	ResetItalic        = "\x1b[23m"
	ResetUnderline     = "\x1b[24m"
	ResetBlink         = "\x1b[25m"
	ResetReverse       = "\x1b[27m"
	ResetHidden        = "\x1b[28m"
	ResetStrikethrough = "\x1b[29m"

	// Cursor control
	CursorUp       = "\x1b[A"
	CursorDown     = "\x1b[B"
	CursorForward  = "\x1b[C"
	CursorBack     = "\x1b[D"
	CursorNextLine = "\x1b[E"
	CursorPrevLine = "\x1b[F"

	// Screen control
	ClearScreen     = "\x1b[2J"
	ClearLine       = "\x1b[2K"
	ClearLineToEnd  = "\x1b[K"
	ClearLineToStart = "\x1b[1K"

	// Hide/show cursor
	HideCursor = "\x1b[?25l"
	ShowCursor = "\x1b[?25h"
)
