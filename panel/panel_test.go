package panel

import (
	"strings"
	"testing"

	"github.com/eberle1080/go-rich"
	"github.com/eberle1080/go-rich/table"
)

func TestPanelBasic(t *testing.T) {
	p := New("Hello")

	console := rich.NewConsole(nil)
	segments := p.Render(console, 80)

	output := segments.String()

	if !strings.Contains(output, "Hello") {
		t.Error("Panel should contain content text")
	}
}

func TestPanelTitle(t *testing.T) {
	p := New("Content").Title("My Title")

	console := rich.NewConsole(nil)
	segments := p.Render(console, 80)

	output := segments.String()

	if !strings.Contains(output, "My Title") {
		t.Error("Panel should contain title")
	}

	if !strings.Contains(output, "Content") {
		t.Error("Panel should contain content")
	}
}

func TestPanelSubtitle(t *testing.T) {
	p := New("Content").
		Title("Title").
		Subtitle("Subtitle")

	console := rich.NewConsole(nil)
	segments := p.Render(console, 80)

	output := segments.String()

	if !strings.Contains(output, "Title") {
		t.Error("Panel should contain title")
	}

	if !strings.Contains(output, "Subtitle") {
		t.Error("Panel should contain subtitle")
	}

	if !strings.Contains(output, "Content") {
		t.Error("Panel should contain content")
	}
}

func TestPanelBoxStyles(t *testing.T) {
	boxes := []table.Box{
		table.BoxSimple,
		table.BoxRounded,
		table.BoxDouble,
		table.BoxHeavy,
		table.BoxASCII,
	}

	for _, box := range boxes {
		p := New("Test").Box(box)

		console := rich.NewConsole(nil)
		segments := p.Render(console, 80)

		if len(segments) == 0 {
			t.Error("Panel should render segments")
		}
	}
}

func TestPanelWidth(t *testing.T) {
	p := New("Test").Width(40).Expand(false)

	console := rich.NewConsole(nil)
	segments := p.Render(console, 80)

	// Just verify fixed width panels render
	if len(segments) == 0 {
		t.Error("Panel should render segments")
	}

	// The width setting should be respected (within reasonable bounds)
	output := segments.String()
	if !strings.Contains(output, "Test") {
		t.Error("Panel should contain content")
	}
}

func TestPanelPadding(t *testing.T) {
	p1 := New("Test").Padding(0).Width(40)
	p2 := New("Test").Padding(3).Width(40)

	console := rich.NewConsole(nil)
	segments1 := p1.Render(console, 80)
	segments2 := p2.Render(console, 80)

	// Both should render
	if len(segments1) == 0 || len(segments2) == 0 {
		t.Error("Both panels should render")
	}

	// Just verify padding values are set correctly
	if p1.padding != 0 {
		t.Error("Padding should be 0")
	}
	if p2.padding != 3 {
		t.Error("Padding should be 3")
	}
}

func TestPanelAlignment(t *testing.T) {
	tests := []struct {
		align Align
		name  string
	}{
		{AlignLeft, "left"},
		{AlignCenter, "center"},
		{AlignRight, "right"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New("X").Align(tt.align).Width(40)

			console := rich.NewConsole(nil)
			segments := p.Render(console, 80)

			if len(segments) == 0 {
				t.Error("Panel should render segments")
			}
		})
	}
}

func TestPanelWithRenderable(t *testing.T) {
	lines := rich.Lines{
		rich.NewRenderableString("Line 1", rich.NewStyle()),
		rich.NewRenderableString("Line 2", rich.NewStyle()),
	}

	p := New(lines)

	console := rich.NewConsole(nil)
	segments := p.Render(console, 80)

	output := segments.String()

	if !strings.Contains(output, "Line 1") {
		t.Error("Panel should contain first line")
	}

	if !strings.Contains(output, "Line 2") {
		t.Error("Panel should contain second line")
	}
}

func TestPanelChaining(t *testing.T) {
	p := New("Test").
		Title("Title").
		Subtitle("Subtitle").
		Box(table.BoxRounded).
		Width(50).
		Padding(2).
		Align(AlignCenter)

	if p.title != "Title" {
		t.Error("Title not set")
	}

	if p.subtitle != "Subtitle" {
		t.Error("Subtitle not set")
	}

	if p.width != 50 {
		t.Error("Width not set")
	}

	if p.padding != 2 {
		t.Error("Padding not set")
	}

	if p.align != AlignCenter {
		t.Error("Align not set")
	}
}

func TestPanelExpand(t *testing.T) {
	p1 := New("Test").Expand(true).Width(0)
	p2 := New("Test").Expand(false).Width(0)

	console := rich.NewConsole(nil)
	segments1 := p1.Render(console, 80)
	segments2 := p2.Render(console, 80)

	// Expanded panel should be wider
	if len(segments1.String()) <= len(segments2.String()) {
		t.Error("Expanded panel should be wider")
	}
}

func TestPanelCustomStyles(t *testing.T) {
	borderStyle := rich.NewStyle().Foreground(rich.Red)
	titleStyle := rich.NewStyle().Bold()

	p := New("Test").
		BorderStyle(borderStyle).
		TitleStyle(titleStyle).
		Title("Styled")

	console := rich.NewConsole(nil)
	segments := p.Render(console, 80)

	// Just verify it renders with custom styles
	if len(segments) == 0 {
		t.Error("Styled panel should render")
	}
}

func TestPanelEmpty(t *testing.T) {
	p := New("")

	console := rich.NewConsole(nil)
	segments := p.Render(console, 80)

	// Should still render borders
	if len(segments) == 0 {
		t.Error("Empty panel should still render")
	}
}

func TestPanelNarrowWidth(t *testing.T) {
	p := New("Test").Width(5)

	console := rich.NewConsole(nil)
	segments := p.Render(console, 80)

	// Should handle narrow widths gracefully
	if len(segments) == 0 {
		t.Error("Narrow panel should still render")
	}
}

func TestSplitIntoLines(t *testing.T) {
	p := New("Test")

	// Test with newlines
	segments := rich.Segments{
		{Text: "Line 1\nLine 2\nLine 3", Style: rich.NewStyle()},
	}

	lines := p.splitIntoLines(segments)

	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}
}

func TestSplitIntoLinesMultipleSegments(t *testing.T) {
	p := New("Test")

	segments := rich.Segments{
		{Text: "Hello ", Style: rich.NewStyle()},
		{Text: "World", Style: rich.NewStyle()},
	}

	lines := p.splitIntoLines(segments)

	if len(lines) != 1 {
		t.Errorf("Expected 1 line, got %d", len(lines))
	}

	if lines[0].String() != "Hello World" {
		t.Errorf("Expected 'Hello World', got %q", lines[0].String())
	}
}

func TestTruncateLine(t *testing.T) {
	p := New("Test")

	segments := rich.Segments{
		{Text: "This is a long line", Style: rich.NewStyle()},
	}

	truncated := p.truncateLine(segments, 10)

	if truncated.String() != "This is a " {
		t.Errorf("Expected 'This is a ', got %q", truncated.String())
	}
}
