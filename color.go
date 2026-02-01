package rich

import (
	"fmt"
	"strconv"
	"strings"
)

// ColorMode represents the color capability of the terminal.
type ColorMode int

const (
	ColorModeNone ColorMode = iota
	ColorModeStandard
	ColorMode256
	ColorModeTrueColor
)

// Color represents a color that can be rendered in various terminal modes.
type Color interface {
	toANSI(mode ColorMode, foreground bool) string
}

// ANSIColor represents one of the 16 standard ANSI colors (0-15).
type ANSIColor int

const (
	Black ANSIColor = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	BrightBlack
	BrightRed
	BrightGreen
	BrightYellow
	BrightBlue
	BrightMagenta
	BrightCyan
	BrightWhite
)

func (c ANSIColor) toANSI(mode ColorMode, foreground bool) string {
	if mode == ColorModeNone {
		return ""
	}

	base := 30
	if !foreground {
		base = 40
	}

	if c >= BrightBlack {
		// Bright colors (8-15) use codes 90-97 (fg) or 100-107 (bg)
		base += 60
		return fmt.Sprintf("\x1b[%dm", base+int(c-BrightBlack))
	}

	// Standard colors (0-7)
	return fmt.Sprintf("\x1b[%dm", base+int(c))
}

// ANSI256Color represents a color from the 256-color palette.
type ANSI256Color int

func (c ANSI256Color) toANSI(mode ColorMode, foreground bool) string {
	if mode == ColorModeNone {
		return ""
	}

	if mode == ColorModeStandard {
		// Downgrade to standard ANSI color (approximate)
		return c.toStandardANSI().toANSI(mode, foreground)
	}

	code := 38
	if !foreground {
		code = 48
	}

	return fmt.Sprintf("\x1b[%d;5;%dm", code, int(c))
}

func (c ANSI256Color) toStandardANSI() ANSIColor {
	// Simple approximation: map 256 colors to 16 standard colors
	n := int(c)
	if n < 8 {
		return ANSIColor(n)
	}
	if n < 16 {
		return ANSIColor(n)
	}
	// For other colors, use a simple heuristic
	if n >= 232 {
		// Grayscale ramp
		if n < 244 {
			return Black
		}
		return White
	}
	// 216-color cube: approximate based on position
	n -= 16
	r := n / 36
	g := (n % 36) / 6
	b := n % 6

	// Map to closest standard color
	colors := []ANSIColor{Black, Red, Green, Yellow, Blue, Magenta, Cyan, White}
	idx := 0
	if r >= 3 {
		idx |= 1
	}
	if g >= 3 {
		idx |= 2
	}
	if b >= 3 {
		idx |= 4
	}
	return colors[idx]
}

// RGBColor represents a 24-bit true color (RGB).
type RGBColor struct {
	R, G, B uint8
}

func (c RGBColor) toANSI(mode ColorMode, foreground bool) string {
	if mode == ColorModeNone {
		return ""
	}

	if mode == ColorModeStandard {
		return c.toStandardANSI().toANSI(mode, foreground)
	}

	if mode == ColorMode256 {
		return c.toANSI256().toANSI(mode, foreground)
	}

	code := 38
	if !foreground {
		code = 48
	}

	return fmt.Sprintf("\x1b[%d;2;%d;%d;%dm", code, c.R, c.G, c.B)
}

func (c RGBColor) toANSI256() ANSI256Color {
	// Convert RGB to 256-color palette
	// Grayscale colors (232-255)
	if c.R == c.G && c.G == c.B {
		if c.R < 8 {
			return 16
		}
		if c.R > 247 {
			return 231
		}
		return ANSI256Color(232 + (int(c.R)-8)/10)
	}

	// Color cube (16-231): 6x6x6
	r := int(c.R) * 5 / 255
	g := int(c.G) * 5 / 255
	b := int(c.B) * 5 / 255

	return ANSI256Color(16 + 36*r + 6*g + b)
}

func (c RGBColor) toStandardANSI() ANSIColor {
	// Simple conversion: threshold at 128
	idx := 0
	if c.R >= 128 {
		idx |= 1
	}
	if c.G >= 128 {
		idx |= 2
	}
	if c.B >= 128 {
		idx |= 4
	}

	colors := []ANSIColor{Black, Red, Green, Yellow, Blue, Magenta, Cyan, White}
	return colors[idx]
}

// RGB creates an RGBColor.
func RGB(r, g, b uint8) RGBColor {
	return RGBColor{R: r, G: g, B: b}
}

// Hex creates an RGBColor from a hex string like "#FF0000" or "FF0000".
func Hex(hex string) (RGBColor, error) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return RGBColor{}, fmt.Errorf("invalid hex color: %s", hex)
	}

	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return RGBColor{}, fmt.Errorf("invalid hex color: %s", hex)
	}

	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return RGBColor{}, fmt.Errorf("invalid hex color: %s", hex)
	}

	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return RGBColor{}, fmt.Errorf("invalid hex color: %s", hex)
	}

	return RGBColor{R: uint8(r), G: uint8(g), B: uint8(b)}, nil
}

// Named color map (basic CSS colors).
var namedColors = map[string]RGBColor{
	"black":   {0, 0, 0},
	"red":     {255, 0, 0},
	"green":   {0, 128, 0},
	"yellow":  {255, 255, 0},
	"blue":    {0, 0, 255},
	"magenta": {255, 0, 255},
	"cyan":    {0, 255, 255},
	"white":   {255, 255, 255},
	"gray":    {128, 128, 128},
	"grey":    {128, 128, 128},
	"orange":  {255, 165, 0},
	"purple":  {128, 0, 128},
	"pink":    {255, 192, 203},
}

// Named creates a color from a named color string.
func Named(name string) (RGBColor, error) {
	color, ok := namedColors[strings.ToLower(name)]
	if !ok {
		return RGBColor{}, fmt.Errorf("unknown color name: %s", name)
	}
	return color, nil
}
