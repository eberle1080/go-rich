package table

// Box defines the characters used to draw table borders.
// Each field represents a specific position or edge in the table structure.
// By customizing these characters, you can create tables with different visual styles.
//
// Box structure:
//
//	TopLeft ──Top── MidTop ──Top── TopRight
//	│                               │
//	Left     cell    Left    cell   Right
//	│                               │
//	HeaderLeft ─HeaderRow─ Mid ─HeaderRow─ HeaderRight
//	│                               │
//	Left     cell    Left    cell   Right
//	│                               │
//	BottomLeft ──Bottom── MidBottom ──Bottom── BottomRight
//
// The package provides predefined box styles (BoxSimple, BoxRounded, BoxDouble, etc.)
// for common use cases.
type Box struct {
	TopLeft  string // Top-left corner character
	Top      string // Top edge character (repeated horizontally)
	TopRight string // Top-right corner character

	Left  string // Left edge character (repeated vertically)
	Right string // Right edge character (repeated vertically)

	BottomLeft  string // Bottom-left corner character
	Bottom      string // Bottom edge character (repeated horizontally)
	BottomRight string // Bottom-right corner character

	MidLeft   string // Left T-junction (where header separator meets left edge)
	MidRight  string // Right T-junction (where header separator meets right edge)
	MidTop    string // Top T-junction (where column separator meets top edge)
	MidBottom string // Bottom T-junction (where column separator meets bottom edge)
	Mid       string // Cross junction (where header separator meets column separator)

	HeaderRow   string // Header separator row character (repeated horizontally)
	HeaderLeft  string // Header separator left junction
	HeaderRight string // Header separator right junction
}

// Predefined box styles.
// These provide ready-to-use border styles for common aesthetic preferences.
var (
	// BoxASCII uses ASCII characters (compatible everywhere).
	// This style works on all terminals and systems, including those without
	// Unicode support. It's the most compatible option but least visually appealing.
	//
	// Example:
	//   +------+------+
	//   | Name | Age  |
	//   +------+------+
	//   | Bob  | 25   |
	//   +------+------+
	BoxASCII = Box{
		TopLeft:     "+",
		Top:         "-",
		TopRight:    "+",
		Left:        "|",
		Right:       "|",
		BottomLeft:  "+",
		Bottom:      "-",
		BottomRight: "+",
		MidLeft:     "+",
		MidRight:    "+",
		MidTop:      "+",
		MidBottom:   "+",
		Mid:         "+",
		HeaderRow:   "-",
		HeaderLeft:  "+",
		HeaderRight: "+",
	}

	// BoxRounded uses rounded corners for a softer appearance.
	// This style uses Unicode box-drawing characters with rounded corners,
	// creating a modern, friendly look. This is the default style.
	//
	// Example:
	//   ╭──────┬──────╮
	//   │ Name │ Age  │
	//   ├──────┼──────┤
	//   │ Bob  │ 25   │
	//   ╰──────┴──────╯
	BoxRounded = Box{
		TopLeft:     "╭",
		Top:         "─",
		TopRight:    "╮",
		Left:        "│",
		Right:       "│",
		BottomLeft:  "╰",
		Bottom:      "─",
		BottomRight: "╯",
		MidLeft:     "├",
		MidRight:    "┤",
		MidTop:      "┬",
		MidBottom:   "┴",
		Mid:         "┼",
		HeaderRow:   "─",
		HeaderLeft:  "├",
		HeaderRight: "┤",
	}

	// BoxDouble uses double-line box drawing characters for emphasis.
	// This style uses thick double-line Unicode box-drawing characters,
	// creating a bold, formal appearance. Good for highlighting important tables.
	//
	// Example:
	//   ╔══════╦══════╗
	//   ║ Name ║ Age  ║
	//   ╠══════╬══════╣
	//   ║ Bob  ║ 25   ║
	//   ╚══════╩══════╝
	BoxDouble = Box{
		TopLeft:     "╔",
		Top:         "═",
		TopRight:    "╗",
		Left:        "║",
		Right:       "║",
		BottomLeft:  "╚",
		Bottom:      "═",
		BottomRight: "╝",
		MidLeft:     "╠",
		MidRight:    "╣",
		MidTop:      "╦",
		MidBottom:   "╩",
		Mid:         "╬",
		HeaderRow:   "═",
		HeaderLeft:  "╠",
		HeaderRight: "╣",
	}

	// BoxHeavy uses heavy box drawing characters for maximum visual weight.
	// This style uses thick single-line Unicode box-drawing characters,
	// creating a strong, bold appearance while remaining cleaner than double-line.
	//
	// Example:
	//   ┏━━━━━━┳━━━━━━┓
	//   ┃ Name ┃ Age  ┃
	//   ┣━━━━━━╋━━━━━━┫
	//   ┃ Bob  ┃ 25   ┃
	//   ┗━━━━━━┻━━━━━━┛
	BoxHeavy = Box{
		TopLeft:     "┏",
		Top:         "━",
		TopRight:    "┓",
		Left:        "┃",
		Right:       "┃",
		BottomLeft:  "┗",
		Bottom:      "━",
		BottomRight: "┛",
		MidLeft:     "┣",
		MidRight:    "┫",
		MidTop:      "┳",
		MidBottom:   "┻",
		Mid:         "╋",
		HeaderRow:   "━",
		HeaderLeft:  "┣",
		HeaderRight: "┫",
	}

	// BoxSimple uses simple single-line characters for a clean look.
	// This style uses thin single-line Unicode box-drawing characters,
	// creating a minimal, clean appearance. Good for dense information displays.
	//
	// Example:
	//   ┌──────┬──────┐
	//   │ Name │ Age  │
	//   ├──────┼──────┤
	//   │ Bob  │ 25   │
	//   └──────┴──────┘
	BoxSimple = Box{
		TopLeft:     "┌",
		Top:         "─",
		TopRight:    "┐",
		Left:        "│",
		Right:       "│",
		BottomLeft:  "└",
		Bottom:      "─",
		BottomRight: "┘",
		MidLeft:     "├",
		MidRight:    "┤",
		MidTop:      "┬",
		MidBottom:   "┴",
		Mid:         "┼",
		HeaderRow:   "─",
		HeaderLeft:  "├",
		HeaderRight: "┤",
	}

	// BoxNone has no borders at all.
	// This style removes all border characters, creating a borderless table.
	// Useful for simple data display where visual separation isn't needed,
	// or when the table is part of a larger bordered container.
	//
	// Example:
	//   Name   Age
	//   Bob    25
	BoxNone = Box{}
)
