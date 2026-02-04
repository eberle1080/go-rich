package progress

import (
	"strings"

	"github.com/eberle1080/go-rich"
)

// ProgressBar represents a visual progress bar that can be rendered to the console.
// It implements the rich.Renderable interface and displays progress as a filled bar
// with customizable appearance, characters, and styles.
//
// ProgressBars can be used in two ways:
//  1. Static rendering: Create, update, and render directly to console
//  2. Live updates: Add to a Progress manager for automatic refresh
//
// Example (static):
//
//	bar := progress.NewBar(100).
//		Description("Download").
//		Width(40)
//	bar.SetProgress(50)
//	console.Render(bar)
//
// Example (live):
//
//	prog := progress.New(console)
//	task := prog.AddBar("Download", 1000)
//	prog.Start()
//	// Update progress from goroutine
//	prog.Update(task, 500)
type ProgressBar struct {
	current int64 // Current progress value
	total   int64 // Total value for completion

	description string // Text description displayed with the bar
	width       int    // Bar width in characters (0 = auto)

	// Styles for different parts of the bar
	barStyle       rich.Style // Style for the bar container
	completeStyle  rich.Style // Style for completed portion
	remainingStyle rich.Style // Style for remaining portion

	// Characters used to draw the bar
	completeChar  string // Character for completed portion (default: "█")
	remainingChar string // Character for remaining portion (default: "░")

	// Tracker for speed and ETA calculations
	tracker *Tracker
}

// NewBar creates a new progress bar with the specified total value.
// The total represents 100% completion. The bar starts at 0 progress.
//
// Default settings:
//   - Complete character: "█" (full block)
//   - Remaining character: "░" (light shade)
//   - Width: 0 (auto-sized)
//   - No styles (uses terminal defaults)
//
// Example:
//
//	bar := progress.NewBar(1000) // Total of 1000 units
func NewBar(total int64) *ProgressBar {
	return &ProgressBar{
		current:        0,
		total:          total,
		width:          0, // Auto-size
		completeChar:   "█",
		remainingChar:  "░",
		barStyle:       rich.NewStyle(),
		completeStyle:  rich.NewStyle(),
		remainingStyle: rich.NewStyle(),
		tracker:        newTracker(),
	}
}

// Description sets the text description displayed with the progress bar.
// This typically appears before the bar itself.
//
// Example:
//
//	bar := progress.NewBar(100).Description("Downloading file.txt")
func (pb *ProgressBar) Description(desc string) *ProgressBar {
	pb.description = desc
	return pb
}

// Width sets the width of the progress bar in characters.
// If set to 0 (default), the bar will auto-size based on available terminal width.
//
// Example:
//
//	bar := progress.NewBar(100).Width(40) // Fixed 40-character bar
func (pb *ProgressBar) Width(width int) *ProgressBar {
	pb.width = width
	return pb
}

// BarStyle sets the style for the bar container.
// This style is applied to border characters if using a bordered bar.
//
// Example:
//
//	bar := progress.NewBar(100).BarStyle(rich.NewStyle().Dim())
func (pb *ProgressBar) BarStyle(style rich.Style) *ProgressBar {
	pb.barStyle = style
	return pb
}

// CompleteStyle sets the style for the completed portion of the bar.
// This affects the color and formatting of the filled part.
//
// Example:
//
//	bar := progress.NewBar(100).CompleteStyle(rich.NewStyle().Foreground(rich.Green))
func (pb *ProgressBar) CompleteStyle(style rich.Style) *ProgressBar {
	pb.completeStyle = style
	return pb
}

// RemainingStyle sets the style for the remaining (incomplete) portion of the bar.
// This affects the color and formatting of the empty part.
//
// Example:
//
//	bar := progress.NewBar(100).RemainingStyle(rich.NewStyle().Dim())
func (pb *ProgressBar) RemainingStyle(style rich.Style) *ProgressBar {
	pb.remainingStyle = style
	return pb
}

// CompleteChar sets the character used for the completed portion of the bar.
// Default is "█" (full block).
//
// Example:
//
//	bar := progress.NewBar(100).CompleteChar("━")
func (pb *ProgressBar) CompleteChar(char string) *ProgressBar {
	pb.completeChar = char
	return pb
}

// RemainingChar sets the character used for the remaining portion of the bar.
// Default is "░" (light shade).
//
// Example:
//
//	bar := progress.NewBar(100).RemainingChar("─")
func (pb *ProgressBar) RemainingChar(char string) *ProgressBar {
	pb.remainingChar = char
	return pb
}

// SetProgress sets the current progress value and updates the tracker.
// The value should be between 0 and total (inclusive).
// Values outside this range are clamped.
//
// This method is thread-safe when used with a Progress manager.
//
// Example:
//
//	bar.SetProgress(50) // Set to 50 out of total
func (pb *ProgressBar) SetProgress(current int64) {
	// Clamp to valid range
	if current < 0 {
		current = 0
	}
	if current > pb.total {
		current = pb.total
	}

	pb.current = current
	pb.tracker.update(current)
}

// Advance increments the current progress by the specified delta.
// This is a convenience method equivalent to SetProgress(current + delta).
//
// Example:
//
//	bar.Advance(10) // Add 10 to current progress
func (pb *ProgressBar) Advance(delta int64) {
	pb.SetProgress(pb.current + delta)
}

// Current returns the current progress value.
func (pb *ProgressBar) Current() int64 {
	return pb.current
}

// Total returns the total progress value (100% completion).
func (pb *ProgressBar) Total() int64 {
	return pb.total
}

// Percentage returns the completion percentage (0.0 to 1.0).
func (pb *ProgressBar) Percentage() float64 {
	if pb.total == 0 {
		return 0
	}
	return float64(pb.current) / float64(pb.total)
}

// IsComplete returns true if progress has reached 100%.
func (pb *ProgressBar) IsComplete() bool {
	return pb.current >= pb.total
}

// Render implements rich.Renderable.
// Converts the progress bar into styled segments for display.
// The width parameter indicates available terminal width.
//
// The bar is rendered as:
//
//	[description] [filled][empty] percentage%
//
// If width is 0 (auto), the bar uses all available space minus description and percentage.
func (pb *ProgressBar) Render(console *rich.Console, width int) rich.Segments {
	segments := rich.Segments{}

	// Render description if present
	if pb.description != "" {
		segments = append(segments, rich.Segment{
			Text:  pb.description + " ",
			Style: rich.NewStyle(),
		})
	}

	// Calculate bar width
	barWidth := pb.width
	if barWidth == 0 {
		// Auto-size: use available width minus description and percentage display
		descLen := len(pb.description)
		if descLen > 0 {
			descLen++ // Account for space
		}
		percentLen := 6 // " 100%"
		barWidth = width - descLen - percentLen
		if barWidth < 10 {
			barWidth = 10 // Minimum bar width
		}
	}

	// Calculate fill width based on percentage
	percentage := pb.Percentage()
	fillWidth := int(float64(barWidth) * percentage)
	emptyWidth := barWidth - fillWidth

	// Render completed portion
	if fillWidth > 0 {
		segments = append(segments, rich.Segment{
			Text:  strings.Repeat(pb.completeChar, fillWidth),
			Style: pb.completeStyle,
		})
	}

	// Render remaining portion
	if emptyWidth > 0 {
		segments = append(segments, rich.Segment{
			Text:  strings.Repeat(pb.remainingChar, emptyWidth),
			Style: pb.remainingStyle,
		})
	}

	// Render percentage
	percentText := " "
	if percentage >= 0.99995 { // Round to 100% at 99.995%
		percentText += "100%"
	} else {
		percentText += formatPercentage(percentage)
	}

	segments = append(segments, rich.Segment{
		Text:  percentText,
		Style: rich.NewStyle(),
	})

	return segments
}

// Measure implements rich.Measurable.
// Returns the size requirements for the progress bar.
func (pb *ProgressBar) Measure(console *rich.Console, maxWidth int) rich.Measurement {
	// Minimum: description + 10 char bar + percentage
	descLen := len(pb.description)
	if descLen > 0 {
		descLen++ // Space after description
	}
	minWidth := descLen + 10 + 6 // 10 char min bar, " 100%"

	// Maximum: use fixed width if set, otherwise prefer 40 chars
	maxBarWidth := pb.width
	if maxBarWidth == 0 {
		maxBarWidth = 40
	}
	maxBarWidth = descLen + maxBarWidth + 6

	return rich.Measurement{
		Minimum: minWidth,
		Maximum: maxBarWidth,
	}
}

// formatPercentage formats a percentage value for display.
// Returns a string like "42.5%" with one decimal place.
// Input p is expected to be 0.0-1.0 (0%-100%).
func formatPercentage(p float64) string {
	// Convert to percentage (0.0-1.0 -> 0-100)
	// Format with one decimal place
	pct := int(p*1000 + 0.5) // Round to nearest 0.1%
	whole := pct / 10
	decimal := pct % 10

	if decimal == 0 {
		return formatInt(whole) + "%"
	}
	return formatInt(whole) + "." + formatInt(decimal) + "%"
}

// formatInt converts an integer to a string without importing fmt or strconv.
func formatInt(n int) string {
	if n == 0 {
		return "0"
	}

	// Handle negative numbers
	negative := false
	if n < 0 {
		negative = true
		n = -n
	}

	// Build digits in reverse
	var digits []byte
	for n > 0 {
		digits = append(digits, byte('0'+n%10))
		n /= 10
	}

	// Reverse digits
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}

	if negative {
		return "-" + string(digits)
	}
	return string(digits)
}
