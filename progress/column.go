package progress

import (
	"strings"
	"time"

	"github.com/eberle1080/go-rich"
)

// Column represents a column in a progress display.
// Columns define what information to show alongside the progress bar
// (e.g., description, percentage, speed, ETA).
//
// Multiple columns can be combined to create custom progress displays.
type Column interface {
	// Render generates the segments for this column for the given task.
	Render(bar *ProgressBar, console *rich.Console) rich.Segments

	// Width returns the width this column will occupy.
	// This is used for layout calculations.
	Width(bar *ProgressBar, console *rich.Console) int
}

// DescriptionColumn displays the task description.
type DescriptionColumn struct {
	style rich.Style
}

// NewDescriptionColumn creates a new description column.
func NewDescriptionColumn() *DescriptionColumn {
	return &DescriptionColumn{
		style: rich.NewStyle(),
	}
}

// Style sets the style for the description text.
func (c *DescriptionColumn) Style(style rich.Style) *DescriptionColumn {
	c.style = style
	return c
}

// Render implements Column.
func (c *DescriptionColumn) Render(bar *ProgressBar, console *rich.Console) rich.Segments {
	if bar.description == "" {
		return rich.Segments{}
	}

	return rich.Segments{
		{Text: bar.description, Style: c.style},
	}
}

// Width implements Column.
func (c *DescriptionColumn) Width(bar *ProgressBar, console *rich.Console) int {
	return len(bar.description)
}

// BarColumn displays the visual progress bar.
type BarColumn struct {
	width          int
	completeChar   string
	remainingChar  string
	completeStyle  rich.Style
	remainingStyle rich.Style
}

// NewBarColumn creates a new bar column with default settings.
func NewBarColumn() *BarColumn {
	return &BarColumn{
		width:          40,
		completeChar:   "█",
		remainingChar:  "░",
		completeStyle:  rich.NewStyle(),
		remainingStyle: rich.NewStyle(),
	}
}

// SetWidth sets the fixed width of the bar.
func (c *BarColumn) SetWidth(width int) *BarColumn {
	c.width = width
	return c
}

// CompleteChar sets the character for the completed portion.
func (c *BarColumn) CompleteChar(char string) *BarColumn {
	c.completeChar = char
	return c
}

// RemainingChar sets the character for the remaining portion.
func (c *BarColumn) RemainingChar(char string) *BarColumn {
	c.remainingChar = char
	return c
}

// CompleteStyle sets the style for the completed portion.
func (c *BarColumn) CompleteStyle(style rich.Style) *BarColumn {
	c.completeStyle = style
	return c
}

// RemainingStyle sets the style for the remaining portion.
func (c *BarColumn) RemainingStyle(style rich.Style) *BarColumn {
	c.remainingStyle = style
	return c
}

// Render implements Column.
func (c *BarColumn) Render(bar *ProgressBar, console *rich.Console) rich.Segments {
	percentage := bar.Percentage()
	fillWidth := int(float64(c.width) * percentage)
	emptyWidth := c.width - fillWidth

	segments := rich.Segments{}

	// Render completed portion
	if fillWidth > 0 {
		segments = append(segments, rich.Segment{
			Text:  strings.Repeat(c.completeChar, fillWidth),
			Style: c.completeStyle,
		})
	}

	// Render remaining portion
	if emptyWidth > 0 {
		segments = append(segments, rich.Segment{
			Text:  strings.Repeat(c.remainingChar, emptyWidth),
			Style: c.remainingStyle,
		})
	}

	return segments
}

// Width implements Column (returns the configured width).
func (c *BarColumn) Width(bar *ProgressBar, console *rich.Console) int {
	return c.width
}

// PercentageColumn displays the completion percentage.
type PercentageColumn struct {
	style rich.Style
}

// NewPercentageColumn creates a new percentage column.
func NewPercentageColumn() *PercentageColumn {
	return &PercentageColumn{
		style: rich.NewStyle(),
	}
}

// Style sets the style for the percentage text.
func (c *PercentageColumn) Style(style rich.Style) *PercentageColumn {
	c.style = style
	return c
}

// Render implements Column.
func (c *PercentageColumn) Render(bar *ProgressBar, console *rich.Console) rich.Segments {
	percentage := bar.Percentage()
	text := formatPercentage(percentage)

	return rich.Segments{
		{Text: text, Style: c.style},
	}
}

// Width implements Column.
func (c *PercentageColumn) Width(bar *ProgressBar, console *rich.Console) int {
	return 5 // "100%"
}

// SpeedColumn displays the current speed in units per second.
type SpeedColumn struct {
	style rich.Style
	unit  string // Optional unit suffix (e.g., "it", "files")
}

// NewSpeedColumn creates a new speed column.
func NewSpeedColumn() *SpeedColumn {
	return &SpeedColumn{
		style: rich.NewStyle(),
		unit:  "it", // Default: iterations
	}
}

// Style sets the style for the speed text.
func (c *SpeedColumn) Style(style rich.Style) *SpeedColumn {
	c.style = style
	return c
}

// Unit sets the unit suffix for the speed display.
// Example: "it" shows "125 it/s", "files" shows "125 files/s"
func (c *SpeedColumn) Unit(unit string) *SpeedColumn {
	c.unit = unit
	return c
}

// Render implements Column.
func (c *SpeedColumn) Render(bar *ProgressBar, console *rich.Console) rich.Segments {
	speed := bar.tracker.Speed()
	text := formatSpeed(speed, c.unit)

	return rich.Segments{
		{Text: text, Style: c.style},
	}
}

// Width implements Column.
func (c *SpeedColumn) Width(bar *ProgressBar, console *rich.Console) int {
	return 12 // "999.9 it/s" with some padding
}

// ETAColumn displays the estimated time remaining.
type ETAColumn struct {
	style rich.Style
}

// NewETAColumn creates a new ETA column.
func NewETAColumn() *ETAColumn {
	return &ETAColumn{
		style: rich.NewStyle(),
	}
}

// Style sets the style for the ETA text.
func (c *ETAColumn) Style(style rich.Style) *ETAColumn {
	c.style = style
	return c
}

// Render implements Column.
func (c *ETAColumn) Render(bar *ProgressBar, console *rich.Console) rich.Segments {
	eta := bar.tracker.ETA(bar.current, bar.total)
	text := formatDuration(eta)

	return rich.Segments{
		{Text: text, Style: c.style},
	}
}

// Width implements Column.
func (c *ETAColumn) Width(bar *ProgressBar, console *rich.Console) int {
	return 10 // "23h 59m" or similar
}

// ElapsedColumn displays the time elapsed since start.
type ElapsedColumn struct {
	style rich.Style
}

// NewElapsedColumn creates a new elapsed time column.
func NewElapsedColumn() *ElapsedColumn {
	return &ElapsedColumn{
		style: rich.NewStyle(),
	}
}

// Style sets the style for the elapsed time text.
func (c *ElapsedColumn) Style(style rich.Style) *ElapsedColumn {
	c.style = style
	return c
}

// Render implements Column.
func (c *ElapsedColumn) Render(bar *ProgressBar, console *rich.Console) rich.Segments {
	elapsed := bar.tracker.Elapsed()
	text := formatDuration(elapsed)

	return rich.Segments{
		{Text: text, Style: c.style},
	}
}

// Width implements Column.
func (c *ElapsedColumn) Width(bar *ProgressBar, console *rich.Console) int {
	return 10 // "23h 59m" or similar
}

// TransferSpeedColumn displays speed with byte unit conversion (B/s, KB/s, MB/s, etc.).
type TransferSpeedColumn struct {
	style rich.Style
}

// NewTransferSpeedColumn creates a new transfer speed column.
func NewTransferSpeedColumn() *TransferSpeedColumn {
	return &TransferSpeedColumn{
		style: rich.NewStyle(),
	}
}

// Style sets the style for the transfer speed text.
func (c *TransferSpeedColumn) Style(style rich.Style) *TransferSpeedColumn {
	c.style = style
	return c
}

// Render implements Column.
func (c *TransferSpeedColumn) Render(bar *ProgressBar, console *rich.Console) rich.Segments {
	speed := bar.tracker.Speed()
	text := formatTransferSpeed(speed)

	return rich.Segments{
		{Text: text, Style: c.style},
	}
}

// Width implements Column.
func (c *TransferSpeedColumn) Width(bar *ProgressBar, console *rich.Console) int {
	return 12 // "999.9 MB/s"
}

// Utility functions for formatting

// formatSpeed formats a speed value with unit.
func formatSpeed(speed float64, unit string) string {
	if speed == 0 {
		return "0 " + unit + "/s"
	}

	// Format with one decimal place
	speedInt := int(speed*10+0.5) / 10 // Round to nearest 0.1
	whole := speedInt / 10
	decimal := speedInt % 10

	result := formatInt(whole)
	if decimal > 0 {
		result += "." + formatInt(decimal)
	}

	return result + " " + unit + "/s"
}

// formatTransferSpeed formats a byte speed with appropriate unit (B/s, KB/s, MB/s, GB/s).
func formatTransferSpeed(bytesPerSecond float64) string {
	if bytesPerSecond == 0 {
		return "0 B/s"
	}

	units := []string{"B/s", "KB/s", "MB/s", "GB/s", "TB/s"}
	value := bytesPerSecond
	unitIndex := 0

	// Find appropriate unit
	for unitIndex < len(units)-1 && value >= 1024 {
		value /= 1024
		unitIndex++
	}

	// Format with one decimal place
	valueInt := int(value*10+0.5) / 10
	whole := valueInt / 10
	decimal := valueInt % 10

	result := formatInt(whole)
	if decimal > 0 && whole < 100 {
		result += "." + formatInt(decimal)
	}

	return result + " " + units[unitIndex]
}

// formatDuration formats a duration for display.
// Shows hours and minutes for long durations, minutes and seconds for short ones.
func formatDuration(d time.Duration) string {
	if d == 0 {
		return "0s"
	}

	seconds := int(d.Seconds())
	if seconds < 0 {
		seconds = 0
	}

	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	if hours > 0 {
		return formatInt(hours) + "h " + formatInt(minutes) + "m"
	} else if minutes > 0 {
		return formatInt(minutes) + "m " + formatInt(secs) + "s"
	} else {
		return formatInt(secs) + "s"
	}
}
