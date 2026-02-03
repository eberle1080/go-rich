package rich

import (
	"fmt"
	"strconv"
	"strings"
)

// ColorMode represents the color capability of the terminal.
// Different terminals support different levels of color output, from no colors
// to full 24-bit RGB (16.7 million colors). The ColorMode determines how color
// values are rendered as ANSI escape sequences.
type ColorMode int

const (
	// ColorModeNone indicates no color support. All color sequences are stripped,
	// resulting in plain text output. This is used for non-terminal outputs or
	// when the NO_COLOR environment variable is set.
	ColorModeNone ColorMode = iota

	// ColorModeStandard indicates support for the 16 standard ANSI colors (8 colors
	// plus 8 bright variants). This is the most widely supported color mode and works
	// on virtually all terminals.
	ColorModeStandard

	// ColorMode256 indicates support for the 256-color palette (216 colors in a 6×6×6
	// cube plus 24 grayscale colors). This is detected via TERM environment variables
	// containing "256color".
	ColorMode256

	// ColorModeTrueColor indicates support for 24-bit RGB colors (16.7 million colors).
	// This is detected via COLORTERM=truecolor or COLORTERM=24bit environment variables.
	// Most modern terminals support this mode.
	ColorModeTrueColor
)

// Color represents a color that can be rendered in various terminal modes.
// This interface abstracts over different color representations (ANSI, 256-color, RGB)
// and allows them to be downgraded gracefully to match the terminal's capabilities.
//
// The toANSI method converts the color to an ANSI escape sequence appropriate for
// the given ColorMode, with automatic fallback to simpler color modes when needed.
type Color interface {
	// toANSI converts the color to an ANSI escape sequence.
	// The mode parameter specifies the target color mode.
	// The foreground parameter indicates whether this is a foreground (true) or background (false) color.
	// Returns an empty string for ColorModeNone.
	toANSI(mode ColorMode, foreground bool) string
}

// ANSIColor represents one of the 16 standard ANSI colors (0-15).
// These colors are the most widely supported and work on virtually all terminals.
// The first 8 colors (Black through White) are the standard colors, while
// the next 8 (BrightBlack through BrightWhite) are their bright variants.
//
// Standard colors use SGR codes 30-37 (foreground) or 40-47 (background).
// Bright colors use SGR codes 90-97 (foreground) or 100-107 (background).
type ANSIColor int

const (
	Black   ANSIColor = iota // Standard black (SGR 30/40)
	Red                      // Standard red (SGR 31/41)
	Green                    // Standard green (SGR 32/42)
	Yellow                   // Standard yellow (SGR 33/43)
	Blue                     // Standard blue (SGR 34/44)
	Magenta                  // Standard magenta (SGR 35/45)
	Cyan                     // Standard cyan (SGR 36/46)
	White                    // Standard white (SGR 37/47)

	BrightBlack   // Bright black/gray (SGR 90/100)
	BrightRed     // Bright red (SGR 91/101)
	BrightGreen   // Bright green (SGR 92/102)
	BrightYellow  // Bright yellow (SGR 93/103)
	BrightBlue    // Bright blue (SGR 94/104)
	BrightMagenta // Bright magenta (SGR 95/105)
	BrightCyan    // Bright cyan (SGR 96/106)
	BrightWhite   // Bright white (SGR 97/107)
)

// toANSI converts the ANSI color to an escape sequence.
// For standard colors (0-7), uses SGR codes 30-37 (foreground) or 40-47 (background).
// For bright colors (8-15), uses SGR codes 90-97 (foreground) or 100-107 (background).
// Returns an empty string if mode is ColorModeNone.
func (c ANSIColor) toANSI(mode ColorMode, foreground bool) string {
	// No color output in ColorModeNone
	if mode == ColorModeNone {
		return ""
	}

	// Foreground colors start at 30, background at 40
	base := 30
	if !foreground {
		base = 40
	}

	if c >= BrightBlack {
		// Bright colors (8-15) use codes 90-97 (fg) or 100-107 (bg)
		// Add 60 to the base (30+60=90 or 40+60=100)
		base += 60
		return fmt.Sprintf("\x1b[%dm", base+int(c-BrightBlack))
	}

	// Standard colors (0-7) use codes 30-37 (fg) or 40-47 (bg)
	return fmt.Sprintf("\x1b[%dm", base+int(c))
}

// ANSI256Color represents a color from the 256-color palette.
// The 256-color palette consists of:
//   - Colors 0-15: Standard ANSI colors (same as ANSIColor)
//   - Colors 16-231: A 6×6×6 RGB cube (216 colors)
//   - Colors 232-255: Grayscale ramp (24 shades from black to white)
//
// This color mode is widely supported on modern terminals and provides
// a good balance between color fidelity and compatibility.
type ANSI256Color int

// toANSI converts the 256-color to an escape sequence.
// Uses SGR codes 38;5;n (foreground) or 48;5;n (background) where n is 0-255.
// Automatically downgrades to standard ANSI colors when mode is ColorModeStandard.
// Returns an empty string if mode is ColorModeNone.
func (c ANSI256Color) toANSI(mode ColorMode, foreground bool) string {
	// No color output in ColorModeNone
	if mode == ColorModeNone {
		return ""
	}

	// Downgrade to 16-color mode if terminal doesn't support 256 colors
	if mode == ColorModeStandard {
		return c.toStandardANSI().toANSI(mode, foreground)
	}

	// Use 256-color escape sequence: ESC[38;5;n m (fg) or ESC[48;5;n m (bg)
	code := 38 // Foreground
	if !foreground {
		code = 48 // Background
	}

	return fmt.Sprintf("\x1b[%d;5;%dm", code, int(c))
}

// toStandardANSI converts a 256-color to the closest standard ANSI color.
// This is used for downgrading colors when the terminal only supports 16 colors.
//
// The conversion logic:
//   - Colors 0-15: Map directly to the same ANSIColor values
//   - Colors 232-255 (grayscale): Map to Black or White based on brightness
//   - Colors 16-231 (RGB cube): Use a threshold-based approximation on each RGB component
//
// The RGB cube uses a 6×6×6 color space. Each component (R, G, B) is compared against
// a threshold of 3 (middle of 0-5 range) to determine if it should be 0 or 1 in the
// standard color space, building a 3-bit color index (RGB → 0b_R_G_B).
func (c ANSI256Color) toStandardANSI() ANSIColor {
	n := int(c)

	// Colors 0-15 are the standard ANSI colors, map directly
	if n < 8 {
		return ANSIColor(n)
	}
	if n < 16 {
		return ANSIColor(n)
	}

	// Colors 232-255 are grayscale
	// These range from near-black to near-white in 24 steps
	if n >= 232 {
		// Use brightness threshold: darker shades → Black, lighter → White
		if n < 244 { // Midpoint of grayscale ramp
			return Black
		}
		return White
	}

	// Colors 16-231 form a 6×6×6 RGB cube
	// Position in cube: color = 16 + 36×r + 6×g + b (where r,g,b ∈ [0,5])
	n -= 16
	r := n / 36       // Red component (0-5)
	g := (n % 36) / 6 // Green component (0-5)
	b := n % 6        // Blue component (0-5)

	// Map to the closest standard color using threshold-based approximation
	// Build a 3-bit index where each bit represents whether that component is "bright"
	colors := []ANSIColor{Black, Red, Green, Yellow, Blue, Magenta, Cyan, White}
	idx := 0
	if r >= 3 { // Red threshold
		idx |= 1 // Set bit 0
	}
	if g >= 3 { // Green threshold
		idx |= 2 // Set bit 1
	}
	if b >= 3 { // Blue threshold
		idx |= 4 // Set bit 2
	}
	return colors[idx]
}

// RGBColor represents a 24-bit true color (RGB).
// Each component (Red, Green, Blue) is an 8-bit value (0-255), providing
// access to approximately 16.7 million colors (256³).
//
// This is the highest quality color mode but requires terminal support.
// Most modern terminals support true color via the COLORTERM environment variable.
//
// RGBColor automatically downgrades to 256-color or 16-color modes when needed.
type RGBColor struct {
	R uint8 // Red component (0-255)
	G uint8 // Green component (0-255)
	B uint8 // Blue component (0-255)
}

// toANSI converts the RGB color to an escape sequence.
// Uses SGR codes 38;2;r;g;b (foreground) or 48;2;r;g;b (background) for true color.
// Automatically downgrades to 256-color when mode is ColorMode256.
// Automatically downgrades to standard ANSI when mode is ColorModeStandard.
// Returns an empty string if mode is ColorModeNone.
func (c RGBColor) toANSI(mode ColorMode, foreground bool) string {
	// No color output in ColorModeNone
	if mode == ColorModeNone {
		return ""
	}

	// Downgrade to 16-color mode if needed
	if mode == ColorModeStandard {
		return c.toStandardANSI().toANSI(mode, foreground)
	}

	// Downgrade to 256-color mode if needed
	if mode == ColorMode256 {
		return c.toANSI256().toANSI(mode, foreground)
	}

	// Use true color escape sequence: ESC[38;2;r;g;b m (fg) or ESC[48;2;r;g;b m (bg)
	code := 38 // Foreground
	if !foreground {
		code = 48 // Background
	}

	return fmt.Sprintf("\x1b[%d;2;%d;%d;%dm", code, c.R, c.G, c.B)
}

// toANSI256 converts an RGB color to the closest 256-color palette entry.
// This is used when the terminal supports 256 colors but not true color.
//
// The conversion logic:
//   - If all RGB components are equal (grayscale), use the grayscale ramp (colors 232-255)
//   - Otherwise, map to the 6×6×6 RGB cube (colors 16-231)
//
// The grayscale ramp covers values from RGB(8,8,8) to RGB(238,238,238) in 24 steps.
// Very dark grays map to color 16, very bright grays to color 231.
//
// The RGB cube maps each 8-bit component (0-255) to a 6-level component (0-5),
// then calculates the color index as: 16 + 36×r + 6×g + b
func (c RGBColor) toANSI256() ANSI256Color {
	// Check for grayscale (all components equal)
	if c.R == c.G && c.G == c.B {
		// Map to grayscale ramp (colors 232-255)
		// Very dark values → color 16 (from standard colors)
		if c.R < 8 {
			return 16
		}
		// Very bright values → color 231 (from RGB cube)
		if c.R > 247 {
			return 231
		}
		// Map the range [8, 247] to 24 grayscale steps
		// Each step covers about 10 units: (247-8)/24 ≈ 10
		return ANSI256Color(232 + (int(c.R)-8)/10)
	}

	// Map to the 6×6×6 RGB cube (colors 16-231)
	// Scale each 8-bit component (0-255) to a 6-level value (0-5)
	r := int(c.R) * 5 / 255
	g := int(c.G) * 5 / 255
	b := int(c.B) * 5 / 255

	// Calculate the cube position: index = 16 + 36*r + 6*g + b
	return ANSI256Color(16 + 36*r + 6*g + b)
}

// toStandardANSI converts an RGB color to the closest standard ANSI color.
// This is used when the terminal only supports 16 colors.
//
// Uses a simple threshold-based approach: each component is compared against 128
// (the midpoint of 0-255). Components >= 128 are considered "bright", creating
// a 3-bit index that maps to one of the 8 standard colors.
//
// The mapping creates this color table:
//   - 0b000 (0) = Black   (RGB all < 128)
//   - 0b001 (1) = Red     (R >= 128, G < 128, B < 128)
//   - 0b010 (2) = Green   (R < 128, G >= 128, B < 128)
//   - 0b011 (3) = Yellow  (R >= 128, G >= 128, B < 128)
//   - 0b100 (4) = Blue    (R < 128, G < 128, B >= 128)
//   - 0b101 (5) = Magenta (R >= 128, G < 128, B >= 128)
//   - 0b110 (6) = Cyan    (R < 128, G >= 128, B >= 128)
//   - 0b111 (7) = White   (RGB all >= 128)
func (c RGBColor) toStandardANSI() ANSIColor {
	// Build a 3-bit index from RGB components using a threshold of 128
	idx := 0
	if c.R >= 128 {
		idx |= 1 // Set bit 0 for red
	}
	if c.G >= 128 {
		idx |= 2 // Set bit 1 for green
	}
	if c.B >= 128 {
		idx |= 4 // Set bit 2 for blue
	}

	// Map the 3-bit index to the corresponding standard color
	colors := []ANSIColor{Black, Red, Green, Yellow, Blue, Magenta, Cyan, White}
	return colors[idx]
}

// RGB creates an RGBColor from individual red, green, and blue components.
// Each component should be in the range 0-255.
//
// Example:
//
//	red := rich.RGB(255, 0, 0)
//	purple := rich.RGB(128, 0, 128)
func RGB(r, g, b uint8) RGBColor {
	return RGBColor{R: r, G: g, B: b}
}

// Hex creates an RGBColor from a hexadecimal color string.
// Accepts formats with or without the leading "#": "#FF0000" or "FF0000".
// The string must be exactly 6 hexadecimal digits (RRGGBB format).
//
// Returns an error if the string format is invalid or contains non-hex characters.
//
// Example:
//
//	red, _ := rich.Hex("#FF0000")
//	blue, _ := rich.Hex("0000FF")
//	invalid, err := rich.Hex("#FFF") // Error: invalid length
func Hex(hex string) (RGBColor, error) {
	// Remove optional leading "#"
	hex = strings.TrimPrefix(hex, "#")

	// Validate length (must be exactly 6 characters for RRGGBB)
	if len(hex) != 6 {
		return RGBColor{}, fmt.Errorf("invalid hex color: %s", hex)
	}

	// Parse red component (first 2 hex digits)
	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return RGBColor{}, fmt.Errorf("invalid hex color: %s", hex)
	}

	// Parse green component (middle 2 hex digits)
	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return RGBColor{}, fmt.Errorf("invalid hex color: %s", hex)
	}

	// Parse blue component (last 2 hex digits)
	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return RGBColor{}, fmt.Errorf("invalid hex color: %s", hex)
	}

	return RGBColor{R: uint8(r), G: uint8(g), B: uint8(b)}, nil
}

// namedColors maps color names to RGB values.
// This provides a convenient way to reference common colors by name.
// The names are case-insensitive when used with the Named function.
//
// Includes basic CSS color names for familiarity. Both "gray" and "grey"
// spellings are supported.
var namedColors = map[string]RGBColor{
	"black":   {0, 0, 0},       // Pure black
	"red":     {255, 0, 0},     // Pure red
	"green":   {0, 128, 0},     // CSS green (darker than pure green)
	"yellow":  {255, 255, 0},   // Pure yellow
	"blue":    {0, 0, 255},     // Pure blue
	"magenta": {255, 0, 255},   // Pure magenta
	"cyan":    {0, 255, 255},   // Pure cyan
	"white":   {255, 255, 255}, // Pure white
	"gray":    {128, 128, 128}, // Medium gray (American spelling)
	"grey":    {128, 128, 128}, // Medium gray (British spelling)
	"orange":  {255, 165, 0},   // Orange
	"purple":  {128, 0, 128},   // Purple (darker than magenta)
	"pink":    {255, 192, 203}, // Pink
}

// Named creates an RGBColor from a named color string.
// Color names are case-insensitive. Supported names include standard CSS colors
// like "red", "blue", "green", as well as "orange", "purple", "pink", and both
// "gray" and "grey".
//
// Returns an error if the color name is not recognized.
//
// Example:
//
//	red, _ := rich.Named("red")
//	gray, _ := rich.Named("GRAY")     // Case-insensitive
//	grey, _ := rich.Named("grey")     // British spelling
//	invalid, err := rich.Named("foo") // Error: unknown color
func Named(name string) (RGBColor, error) {
	// Lookup is case-insensitive
	color, ok := namedColors[strings.ToLower(name)]
	if !ok {
		return RGBColor{}, fmt.Errorf("unknown color name: %s", name)
	}
	return color, nil
}
