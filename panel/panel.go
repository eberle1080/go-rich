package panel

import (
	"strings"

	"github.com/eberle1080/go-rich"
	"github.com/eberle1080/go-rich/table"
)

// Align specifies how content is aligned within a panel.
type Align int

const (
	AlignLeft Align = iota
	AlignCenter
	AlignRight
)

// Panel represents a bordered container for content.
type Panel struct {
	content      rich.Renderable
	title        string
	subtitle     string
	box          table.Box
	width        int // 0 = auto
	padding      int
	align        Align
	borderStyle  rich.Style
	titleStyle   rich.Style
	contentStyle rich.Style
	expand       bool
}

// New creates a new panel with the given content.
func New(content interface{}) *Panel {
	var renderable rich.Renderable

	switch c := content.(type) {
	case string:
		renderable = rich.NewRenderableString(c, rich.NewStyle())
	case rich.Renderable:
		renderable = c
	default:
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
		expand:       true,
	}
}

// Title sets the panel title (displayed at top).
func (p *Panel) Title(title string) *Panel {
	p.title = title
	return p
}

// Subtitle sets the panel subtitle (displayed at bottom).
func (p *Panel) Subtitle(subtitle string) *Panel {
	p.subtitle = subtitle
	return p
}

// Box sets the border style.
func (p *Panel) Box(box table.Box) *Panel {
	p.box = box
	return p
}

// Width sets a fixed width for the panel (0 = auto).
func (p *Panel) Width(width int) *Panel {
	p.width = width
	return p
}

// Padding sets the internal padding.
func (p *Panel) Padding(padding int) *Panel {
	p.padding = padding
	return p
}

// Align sets the content alignment.
func (p *Panel) Align(align Align) *Panel {
	p.align = align
	return p
}

// BorderStyle sets the border style.
func (p *Panel) BorderStyle(style rich.Style) *Panel {
	p.borderStyle = style
	return p
}

// TitleStyle sets the title style.
func (p *Panel) TitleStyle(style rich.Style) *Panel {
	p.titleStyle = style
	return p
}

// ContentStyle sets the content style.
func (p *Panel) ContentStyle(style rich.Style) *Panel {
	p.contentStyle = style
	return p
}

// Expand controls whether the panel expands to fill available width.
func (p *Panel) Expand(expand bool) *Panel {
	p.expand = expand
	return p
}

// Render implements rich.Renderable.
func (p *Panel) Render(console *rich.Console, maxWidth int) rich.Segments {
	// Determine panel width
	width := p.width
	if width == 0 || width > maxWidth {
		if p.expand {
			width = maxWidth
		} else {
			// Auto-size to content
			width = p.measureContent(console, maxWidth)
		}
	}

	// Ensure minimum width
	if width < 3 {
		width = 3
	}

	// Content width (excluding borders and padding)
	contentWidth := width - 2 - (p.padding * 2)
	if contentWidth < 1 {
		contentWidth = 1
	}

	var segments rich.Segments

	// Top border
	segments = append(segments, p.renderTopBorder(width)...)
	segments = append(segments, rich.Segment{Text: "\n"})

	// Title
	if p.title != "" {
		segments = append(segments, p.renderTitle(width)...)
		segments = append(segments, rich.Segment{Text: "\n"})
	}

	// Content
	contentSegments := p.content.Render(console, contentWidth)
	contentLines := p.splitIntoLines(contentSegments)

	for _, line := range contentLines {
		segments = append(segments, p.renderContentLine(line, width, contentWidth)...)
		segments = append(segments, rich.Segment{Text: "\n"})
	}

	// Subtitle
	if p.subtitle != "" {
		segments = append(segments, p.renderSubtitle(width)...)
		segments = append(segments, rich.Segment{Text: "\n"})
	}

	// Bottom border
	segments = append(segments, p.renderBottomBorder(width)...)

	return segments
}

// measureContent measures the content width.
func (p *Panel) measureContent(console *rich.Console, maxWidth int) int {
	// Try to measure if content is Measurable
	if measurable, ok := p.content.(rich.Measurable); ok {
		measurement := measurable.Measure(console, maxWidth-2-(p.padding*2))
		return measurement.Maximum + 2 + (p.padding * 2)
	}

	// Otherwise render and measure
	segments := p.content.Render(console, maxWidth-2-(p.padding*2))
	lines := p.splitIntoLines(segments)

	maxLen := 0
	for _, line := range lines {
		lineLen := line.Length()
		if lineLen > maxLen {
			maxLen = lineLen
		}
	}

	return maxLen + 2 + (p.padding * 2)
}

// splitIntoLines splits segments into lines.
func (p *Panel) splitIntoLines(segments rich.Segments) []rich.Segments {
	var lines []rich.Segments
	var currentLine rich.Segments

	for _, seg := range segments {
		if strings.Contains(seg.Text, "\n") {
			// Split by newlines
			parts := strings.Split(seg.Text, "\n")
			for i, part := range parts {
				if i > 0 {
					lines = append(lines, currentLine)
					currentLine = nil
				}
				if part != "" {
					currentLine = append(currentLine, rich.Segment{
						Text:  part,
						Style: seg.Style,
					})
				}
			}
		} else {
			currentLine = append(currentLine, seg)
		}
	}

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

// truncateLine truncates a line to the given width.
func (p *Panel) truncateLine(line rich.Segments, width int) rich.Segments {
	var result rich.Segments
	remaining := width

	for _, seg := range line {
		segLen := len(seg.Text)
		if segLen <= remaining {
			result = append(result, seg)
			remaining -= segLen
		} else if remaining > 0 {
			result = append(result, rich.Segment{
				Text:  seg.Text[:remaining],
				Style: seg.Style,
			})
			break
		} else {
			break
		}
	}

	return result
}
