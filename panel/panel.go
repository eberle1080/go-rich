// Package panel provides bordered containers for rich terminal content.
//
// A panel wraps content (text or other renderables) in a bordered box with
// optional title and subtitle. Panels are useful for:
//   - Highlighting important information
//   - Grouping related content
//   - Creating visual hierarchy in terminal output
//   - Drawing attention to alerts, warnings, or status messages
//
// # Basic Usage
//
// Create a simple panel with text content:
//
//	p := panel.New("Hello, World!")
//	console.Renderln(p)
//
// # Customization
//
// Panels support extensive customization:
//
//	p := panel.New("Important message").
//		Title("Alert").
//		Subtitle("Press any key").
//		Box(table.BoxDouble).
//		BorderStyle(rich.NewStyle().Foreground(rich.Red)).
//		Padding(2)
//
// # Border Styles
//
// Panels use the same Box styles as tables:
//
//	panel.New("Content").Box(table.BoxRounded)   // Rounded corners
//	panel.New("Content").Box(table.BoxDouble)    // Double-line borders
//	panel.New("Content").Box(table.BoxHeavy)     // Heavy borders
//	panel.New("Content").Box(table.BoxSimple)    // Simple borders
//
// # Width Control
//
// By default, panels expand to fill available width. Control this with Width() and Expand():
//
//	panel.New("Text").Width(50)            // Fixed 50-character width
//	panel.New("Text").Expand(false)        // Auto-size to content
//	panel.New("Text").Expand(true)         // Expand to fill width (default)
//
// # Alignment
//
// Control how content is aligned within the panel:
//
//	panel.New("Text").Align(panel.AlignLeft)    // Left-aligned (default)
//	panel.New("Text").Align(panel.AlignCenter)  // Centered
//	panel.New("Text").Align(panel.AlignRight)   // Right-aligned
//
// # Renderables as Content
//
// Panels can contain any Renderable, not just strings:
//
//	tbl := table.New().Headers("Name", "Age").Row("Alice", "30")
//	p := panel.New(tbl).Title("User Data")
//	console.Renderln(p)
package panel

import (
	"strings"

	"github.com/eberle1080/go-rich"
	"github.com/eberle1080/go-rich/table"
)

// Align specifies how content is aligned within a panel.
// Alignment affects how content is positioned horizontally when it's
// narrower than the panel's inner width.
type Align int

const (
	// AlignLeft aligns content to the left side of the panel.
	// This is the default and most common alignment.
	AlignLeft Align = iota

	// AlignCenter centers content within the panel.
	// Useful for titles, short messages, or creating symmetry.
	AlignCenter

	// AlignRight aligns content to the right side of the panel.
	// Less common but useful for specific layouts.
	AlignRight
)

// Panel represents a bordered container for content.
// Panels wrap content in a box with borders, optional title and subtitle,
// and customizable styling. They implement rich.Renderable and can be
// rendered to a Console.
//
// A panel consists of:
//   - Content: Any Renderable (text, tables, or custom renderables)
//   - Optional title: Displayed at the top, centered
//   - Optional subtitle: Displayed at the bottom, centered
//   - Border: Using table.Box styles (rounded, double-line, etc.)
//   - Styling: Separate styles for borders, title, and content
//   - Layout: Width, padding, alignment, and expand behavior
//
// Panels are built using a fluent API:
//
//	panel.New("Message").
//		Title("Alert").
//		BorderStyle(rich.NewStyle().Foreground(rich.Red))
type Panel struct {
	content rich.Renderable // The content to display inside the panel

	title    string // Optional title displayed at top (centered)
	subtitle string // Optional subtitle displayed at bottom (centered)

	box table.Box // Border characters (from table package)

	width   int   // Fixed width (0 = auto-size based on content/expand)
	padding int   // Internal padding (spaces between border and content)
	align   Align // Content alignment (left, center, right)

	borderStyle  rich.Style // Style for border characters
	titleStyle   rich.Style // Style for title and subtitle text
	contentStyle rich.Style // Style for content (if content is a string)

	expand bool // If true, expand to fill available width; if false, fit to content
}

// New creates a new panel with the given content.
// The content can be either:
//   - A string: Will be wrapped in a RenderableString
//   - A rich.Renderable: Used directly (e.g., tables, custom renderables)
//   - Any other type: Treated as empty content
//
// The panel is initialized with sensible defaults:
//   - Rounded box borders (BoxRounded)
//   - 1 character of padding
//   - Left-aligned content
//   - Dim border style
//   - Bold title style
//   - Expand enabled (fills available width)
//
// Example:
//
//	// String content
//	p := panel.New("Hello, World!")
//
//	// Renderable content
//	tbl := table.New().Headers("A", "B")
//	p := panel.New(tbl)
func New(content interface{}) *Panel {
	var renderable rich.Renderable

	// Convert content to Renderable
	switch c := content.(type) {
	case string:
		// Wrap string in a renderable
		renderable = rich.NewRenderableString(c, rich.NewStyle())
	case rich.Renderable:
		// Already a renderable, use directly
		renderable = c
	default:
		// Unknown type, use empty content
		renderable = rich.NewRenderableString("", rich.NewStyle())
	}

	return &Panel{
		content:      renderable,
		box:          table.BoxRounded,
		padding:      1,
		align:        AlignLeft,
		borderStyle:  rich.NewStyle().Dim(),
		titleStyle:   rich.NewStyle().Bold(),
		contentStyle: rich.NewStyle(),
		expand:       true, // Fill available width by default
	}
}

// Title sets the panel title displayed at the top.
// The title is centered and shown in its own row above the content.
// If empty (default), no title row is displayed.
//
// Example:
//
//	panel.New("Message").Title("Alert")
func (p *Panel) Title(title string) *Panel {
	p.title = title
	return p
}

// Subtitle sets the panel subtitle displayed at the bottom.
// The subtitle is centered and shown in its own row below the content.
// If empty (default), no subtitle row is displayed.
//
// Example:
//
//	panel.New("Message").Subtitle("Press any key")
func (p *Panel) Subtitle(subtitle string) *Panel {
	p.subtitle = subtitle
	return p
}

// Box sets the border style using predefined or custom box characters.
// See table.Box and predefined styles (BoxRounded, BoxDouble, etc.) for options.
// Default is BoxRounded.
//
// Example:
//
//	panel.New("Message").Box(table.BoxDouble)
func (p *Panel) Box(box table.Box) *Panel {
	p.box = box
	return p
}

// Width sets a fixed width for the panel.
// If width is 0 (default), the panel width is determined by:
//   - The expand setting (fill available width or fit to content)
//   - The content size (when expand is false)
//
// If width > 0, the panel will be exactly that width.
//
// Example:
//
//	panel.New("Message").Width(50) // Fixed 50-character width
func (p *Panel) Width(width int) *Panel {
	p.width = width
	return p
}

// Padding sets the internal padding in characters.
// Padding is added to all four sides (top, right, bottom, left) of the content.
// Default is 1.
//
// Example:
//
//	panel.New("Message").Padding(2) // 2 spaces on all sides
func (p *Panel) Padding(padding int) *Panel {
	p.padding = padding
	return p
}

// Align sets the content alignment.
// Determines how content is positioned horizontally when it's narrower
// than the panel's inner width.
// Default is AlignLeft.
//
// Example:
//
//	panel.New("Message").Align(panel.AlignCenter)
func (p *Panel) Align(align Align) *Panel {
	p.align = align
	return p
}

// BorderStyle sets the style for all border characters.
// This affects the visual appearance of the borders but not their shape.
// Default is dim (faint) style.
//
// Example:
//
//	panel.New("Error").BorderStyle(rich.NewStyle().Foreground(rich.Red))
func (p *Panel) BorderStyle(style rich.Style) *Panel {
	p.borderStyle = style
	return p
}

// TitleStyle sets the style for the title and subtitle text.
// Affects both title (if set) and subtitle (if set).
// Default is bold style.
//
// Example:
//
//	panel.New("Message").
//		Title("Alert").
//		TitleStyle(rich.NewStyle().Bold().Foreground(rich.Yellow))
func (p *Panel) TitleStyle(style rich.Style) *Panel {
	p.titleStyle = style
	return p
}

// ContentStyle sets the style for the content.
// Note: This only affects string content. If content is a Renderable,
// the Renderable is responsible for its own styling.
// Default is unstyled.
//
// Example:
//
//	panel.New("Message").ContentStyle(rich.NewStyle().Foreground(rich.Green))
func (p *Panel) ContentStyle(style rich.Style) *Panel {
	p.contentStyle = style
	return p
}

// Expand controls whether the panel expands to fill available width.
// When true (default): Panel stretches to use all available width.
// When false: Panel auto-sizes to fit content.
//
// This setting is ignored if Width() is used to set a fixed width.
//
// Example:
//
//	panel.New("Small").Expand(false) // Fit to content width
func (p *Panel) Expand(expand bool) *Panel {
	p.expand = expand
	return p
}

// Render implements rich.Renderable.
// Converts the panel into styled segments that can be displayed on the console.
//
// The rendering process:
//  1. Determine panel width (from Width, expand setting, or content measurement)
//  2. Calculate content width (panel width minus borders and padding)
//  3. Render top border
//  4. Render title row (if title is set)
//  5. Render content lines with borders and padding
//  6. Render subtitle row (if subtitle is set)
//  7. Render bottom border
//
// The maxWidth parameter is the maximum available width for the panel.
// The console parameter provides access to color mode and other settings.
func (p *Panel) Render(console *rich.Console, maxWidth int) rich.Segments {
	// Determine the panel's total width
	width := p.width

	// If no fixed width is set, calculate based on expand and content
	if width == 0 || width > maxWidth {
		if p.expand {
			// Use full available width
			width = maxWidth
		} else {
			// Auto-size to fit content
			width = p.measureContent(console, maxWidth)
		}
	}

	// Ensure panel is wide enough for borders
	if width < 3 {
		width = 3 // Minimum: left border + 1 char + right border
	}

	// Calculate content width (total width minus borders and padding)
	// Formula: contentWidth = width - 2 (borders) - 2*padding
	contentWidth := width - 2 - (p.padding * 2)
	if contentWidth < 1 {
		contentWidth = 1 // Ensure at least 1 char of content width
	}

	var segments rich.Segments

	// Render top border
	segments = append(segments, p.renderTopBorder(width)...)
	segments = append(segments, rich.Segment{Text: "\n"})

	// Render title if present
	if p.title != "" {
		segments = append(segments, p.renderTitle(width)...)
		segments = append(segments, rich.Segment{Text: "\n"})
	}

	// Render content
	// First, get the content as segments
	contentSegments := p.content.Render(console, contentWidth)

	// Split content into lines (handle newlines)
	contentLines := p.splitIntoLines(contentSegments)

	// Render each line with borders and padding
	for _, line := range contentLines {
		segments = append(segments, p.renderContentLine(line, width, contentWidth)...)
		segments = append(segments, rich.Segment{Text: "\n"})
	}

	// Render subtitle if present
	if p.subtitle != "" {
		segments = append(segments, p.renderSubtitle(width)...)
		segments = append(segments, rich.Segment{Text: "\n"})
	}

	// Render bottom border
	segments = append(segments, p.renderBottomBorder(width)...)

	return segments
}

// measureContent measures the optimal width for the panel based on content.
// Used when expand is false and no fixed width is set.
//
// The measurement strategy:
//  1. If content implements Measurable, use its Measure method (efficient)
//  2. Otherwise, render content and measure the longest line (fallback)
//
// Returns the total panel width including borders and padding:
//
//	contentWidth + 2 (borders) + 2*padding
//
// The maxWidth parameter constrains the measurement to available space.
func (p *Panel) measureContent(console *rich.Console, maxWidth int) int {
	// Try to measure efficiently if content supports it
	if measurable, ok := p.content.(rich.Measurable); ok {
		// Measure content with available inner width
		measurement := measurable.Measure(console, maxWidth-2-(p.padding*2))
		// Add borders and padding to get total panel width
		return measurement.Maximum + 2 + (p.padding * 2)
	}

	// Fallback: render and measure manually
	// Render content at max available width
	segments := p.content.Render(console, maxWidth-2-(p.padding*2))

	// Split into lines
	lines := p.splitIntoLines(segments)

	// Find the longest line
	maxLen := 0
	for _, line := range lines {
		lineLen := line.Length()
		if lineLen > maxLen {
			maxLen = lineLen
		}
	}

	// Add borders and padding to content width
	return maxLen + 2 + (p.padding * 2)
}

// splitIntoLines splits segments into lines based on newline characters.
// This is necessary because content may contain newlines, and each line
// needs to be rendered separately with its own borders.
//
// The algorithm:
//  1. Process each segment
//  2. If segment contains "\n", split it into parts
//  3. Each newline creates a new line in the output
//  4. Empty parts (consecutive newlines) are skipped
//  5. Preserve the segment style for each part
//
// Returns a slice of segment slices, one per line.
//
// Example:
//
//	Input: [{"Hello\nWorld", style}]
//	Output: [[{"Hello", style}], [{"World", style}]]
func (p *Panel) splitIntoLines(segments rich.Segments) []rich.Segments {
	var lines []rich.Segments
	var currentLine rich.Segments

	for _, seg := range segments {
		if strings.Contains(seg.Text, "\n") {
			// Segment contains newlines, split it
			parts := strings.Split(seg.Text, "\n")

			for i, part := range parts {
				if i > 0 {
					// New line: save current line and start fresh
					lines = append(lines, currentLine)
					currentLine = nil
				}

				// Add non-empty parts to current line
				if part != "" {
					currentLine = append(currentLine, rich.Segment{
						Text:  part,
						Style: seg.Style,
					})
				}
			}
		} else {
			// No newlines, add to current line
			currentLine = append(currentLine, seg)
		}
	}

	// Don't forget the last line
	if len(currentLine) > 0 {
		lines = append(lines, currentLine)
	}

	return lines
}

// renderTopBorder renders the top border.
func (p *Panel) renderTopBorder(width int) rich.Segments {
	innerWidth := width - 2
	line := p.box.TopLeft + strings.Repeat(p.box.Top, innerWidth) + p.box.TopRight

	return rich.Segments{{Text: line, Style: p.borderStyle}}
}

// renderBottomBorder renders the bottom border.
func (p *Panel) renderBottomBorder(width int) rich.Segments {
	innerWidth := width - 2
	line := p.box.BottomLeft + strings.Repeat(p.box.Bottom, innerWidth) + p.box.BottomRight

	return rich.Segments{{Text: line, Style: p.borderStyle}}
}

// renderTitle renders the title line.
func (p *Panel) renderTitle(width int) rich.Segments {
	innerWidth := width - 2
	titleLen := len(p.title)

	var segments rich.Segments

	segments = append(segments, rich.Segment{
		Text:  p.box.Left,
		Style: p.borderStyle,
	})

	if titleLen >= innerWidth {
		// Title too long, truncate
		segments = append(segments, rich.Segment{
			Text:  p.title[:innerWidth],
			Style: p.titleStyle,
		})
	} else {
		// Center the title
		leftPad := (innerWidth - titleLen) / 2
		rightPad := innerWidth - titleLen - leftPad

		if leftPad > 0 {
			segments = append(segments, rich.Segment{Text: strings.Repeat(" ", leftPad)})
		}

		segments = append(segments, rich.Segment{
			Text:  p.title,
			Style: p.titleStyle,
		})

		if rightPad > 0 {
			segments = append(segments, rich.Segment{Text: strings.Repeat(" ", rightPad)})
		}
	}

	segments = append(segments, rich.Segment{
		Text:  p.box.Right,
		Style: p.borderStyle,
	})

	return segments
}

// renderSubtitle renders the subtitle line.
func (p *Panel) renderSubtitle(width int) rich.Segments {
	innerWidth := width - 2
	subtitleLen := len(p.subtitle)

	var segments rich.Segments

	segments = append(segments, rich.Segment{
		Text:  p.box.Left,
		Style: p.borderStyle,
	})

	if subtitleLen >= innerWidth {
		// Subtitle too long, truncate
		segments = append(segments, rich.Segment{
			Text:  p.subtitle[:innerWidth],
			Style: p.titleStyle,
		})
	} else {
		// Center the subtitle
		leftPad := (innerWidth - subtitleLen) / 2
		rightPad := innerWidth - subtitleLen - leftPad

		if leftPad > 0 {
			segments = append(segments, rich.Segment{Text: strings.Repeat(" ", leftPad)})
		}

		segments = append(segments, rich.Segment{
			Text:  p.subtitle,
			Style: p.titleStyle,
		})

		if rightPad > 0 {
			segments = append(segments, rich.Segment{Text: strings.Repeat(" ", rightPad)})
		}
	}

	segments = append(segments, rich.Segment{
		Text:  p.box.Right,
		Style: p.borderStyle,
	})

	return segments
}

// renderContentLine renders a single line of content.
func (p *Panel) renderContentLine(line rich.Segments, width int, contentWidth int) rich.Segments {
	var segments rich.Segments

	// Left border
	segments = append(segments, rich.Segment{
		Text:  p.box.Left,
		Style: p.borderStyle,
	})

	// Left padding
	if p.padding > 0 {
		segments = append(segments, rich.Segment{Text: strings.Repeat(" ", p.padding)})
	}

	// Content (aligned)
	lineLen := line.Length()
	if lineLen > contentWidth {
		// Truncate
		line = p.truncateLine(line, contentWidth)
		segments = append(segments, line...)
	} else {
		// Align
		padding := contentWidth - lineLen

		switch p.align {
		case AlignLeft:
			segments = append(segments, line...)
			if padding > 0 {
				segments = append(segments, rich.Segment{Text: strings.Repeat(" ", padding)})
			}

		case AlignRight:
			if padding > 0 {
				segments = append(segments, rich.Segment{Text: strings.Repeat(" ", padding)})
			}
			segments = append(segments, line...)

		case AlignCenter:
			leftPad := padding / 2
			rightPad := padding - leftPad
			if leftPad > 0 {
				segments = append(segments, rich.Segment{Text: strings.Repeat(" ", leftPad)})
			}
			segments = append(segments, line...)
			if rightPad > 0 {
				segments = append(segments, rich.Segment{Text: strings.Repeat(" ", rightPad)})
			}
		}
	}

	// Right padding
	if p.padding > 0 {
		segments = append(segments, rich.Segment{Text: strings.Repeat(" ", p.padding)})
	}

	// Right border
	segments = append(segments, rich.Segment{
		Text:  p.box.Right,
		Style: p.borderStyle,
	})

	return segments
}

// truncateLine truncates a line of segments to fit within a given width.
// This is used when a line is longer than the available content width.
//
// The algorithm:
//  1. Process segments in order
//  2. Include full segments while there's space
//  3. Truncate the first segment that doesn't fit
//  4. Discard all following segments
//
// The width parameter is in characters (not bytes).
//
// Note: This uses byte-based truncation (seg.Text[:remaining]), which may
// split multi-byte Unicode characters. This is a known limitation.
//
// Example:
//
//	Input: [{"Hello", style1}, {" ", style2}, {"World!", style3}]
//	Width: 8
//	Output: [{"Hello", style1}, {" ", style2}, {"Wo", style3}]
func (p *Panel) truncateLine(line rich.Segments, width int) rich.Segments {
	var result rich.Segments
	remaining := width

	for _, seg := range line {
		segLen := len(seg.Text)

		if segLen <= remaining {
			// Segment fits completely
			result = append(result, seg)
			remaining -= segLen
		} else if remaining > 0 {
			// Segment needs truncation
			result = append(result, rich.Segment{
				Text:  seg.Text[:remaining],
				Style: seg.Style,
			})
			// Stop processing after truncation
			break
		} else {
			// No space left, stop
			break
		}
	}

	return result
}
