package table

// Box defines the characters used to draw table borders.
type Box struct {
	TopLeft      string
	Top          string
	TopRight     string
	Left         string
	Right        string
	BottomLeft   string
	Bottom       string
	BottomRight  string
	MidLeft      string
	MidRight     string
	MidTop       string
	MidBottom    string
	Mid          string
	HeaderRow    string
	HeaderLeft   string
	HeaderRight  string
}

// Predefined box styles.
var (
	// BoxASCII uses ASCII characters (compatible everywhere).
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

	// BoxRounded uses rounded corners.
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

	// BoxDouble uses double-line box drawing characters.
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

	// BoxHeavy uses heavy box drawing characters.
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

	// BoxSimple uses simple single-line characters.
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

	// BoxNone has no borders.
	BoxNone = Box{}
)
