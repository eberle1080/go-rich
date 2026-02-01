package table

import (
	"strings"
	"testing"

	"github.com/eberle1080/go-rich"
)

func TestTableBasic(t *testing.T) {
	table := New().
		Headers("A", "B", "C").
		Row("1", "2", "3").
		Row("4", "5", "6")

	console := rich.NewConsole(nil)
	segments := table.Render(console, 80)

	output := segments.String()

	// Should contain headers
	if !strings.Contains(output, "A") {
		t.Error("Output should contain header 'A'")
	}

	// Should contain data
	if !strings.Contains(output, "1") {
		t.Error("Output should contain data '1'")
	}
}

func TestTableTitle(t *testing.T) {
	table := New().
		Title("Test Table").
		Headers("A", "B").
		Row("1", "2")

	console := rich.NewConsole(nil)
	segments := table.Render(console, 80)

	output := segments.String()

	if !strings.Contains(output, "Test Table") {
		t.Error("Output should contain title")
	}
}

func TestTableBoxStyles(t *testing.T) {
	boxes := []Box{
		BoxASCII,
		BoxRounded,
		BoxDouble,
		BoxHeavy,
		BoxSimple,
	}

	for _, box := range boxes {
		table := New().
			Box(box).
			Headers("A", "B").
			Row("1", "2")

		console := rich.NewConsole(nil)
		segments := table.Render(console, 80)

		if len(segments) == 0 {
			t.Errorf("Table with box style should render segments")
		}
	}
}

func TestTableNoHeader(t *testing.T) {
	table := New().
		ShowHeader(false).
		Headers("A", "B").
		Row("1", "2")

	console := rich.NewConsole(nil)
	segments := table.Render(console, 80)

	output := segments.String()

	// Should not contain header text (but headers define columns)
	// The headers are still needed to know how many columns
	if !strings.Contains(output, "1") {
		t.Error("Output should contain data")
	}
}

func TestTableNoEdge(t *testing.T) {
	table := New().
		ShowEdge(false).
		Headers("A", "B").
		Row("1", "2")

	console := rich.NewConsole(nil)
	segments := table.Render(console, 80)

	output := segments.String()

	// Should still have content
	if !strings.Contains(output, "A") {
		t.Error("Output should contain header")
	}
}

func TestTableAlignment(t *testing.T) {
	table := New().
		AddColumn(NewColumn("Left").WithAlign(AlignLeft)).
		AddColumn(NewColumn("Center").WithAlign(AlignCenter)).
		AddColumn(NewColumn("Right").WithAlign(AlignRight)).
		Row("L", "C", "R")

	console := rich.NewConsole(nil)
	segments := table.Render(console, 80)

	if len(segments) == 0 {
		t.Error("Table should render segments")
	}
}

func TestTableFixedWidth(t *testing.T) {
	table := New().
		AddColumn(NewColumn("Fixed").WithWidth(10)).
		Row("text")

	console := rich.NewConsole(nil)
	segments := table.Render(console, 80)

	if len(segments) == 0 {
		t.Error("Table should render segments")
	}
}

func TestTableEmpty(t *testing.T) {
	table := New().
		Headers("A", "B")

	console := rich.NewConsole(nil)
	segments := table.Render(console, 80)

	// Should render even with no data rows
	if len(segments) == 0 {
		t.Error("Empty table should still render header")
	}
}

func TestTableNoColumns(t *testing.T) {
	table := New()

	console := rich.NewConsole(nil)
	segments := table.Render(console, 80)

	// Should return nil or empty segments
	if segments != nil && len(segments) > 0 {
		t.Error("Table with no columns should not render")
	}
}

func TestColumnChaining(t *testing.T) {
	col := NewColumn("Test").
		WithWidth(20).
		WithMinWidth(10).
		WithMaxWidth(30).
		WithAlign(AlignCenter).
		WithNoWrap()

	if col.Width != 20 {
		t.Error("Width not set")
	}
	if col.MinWidth != 10 {
		t.Error("MinWidth not set")
	}
	if col.MaxWidth != 30 {
		t.Error("MaxWidth not set")
	}
	if col.Align != AlignCenter {
		t.Error("Align not set")
	}
	if !col.NoWrap {
		t.Error("NoWrap not set")
	}
}

func TestTableChaining(t *testing.T) {
	table := New().
		Title("Title").
		Box(BoxRounded).
		ShowHeader(false).
		ShowEdge(false).
		Padding(2).
		Headers("A", "B").
		Row("1", "2")

	if table.title != "Title" {
		t.Error("Title not set")
	}
	if table.showHeader {
		t.Error("ShowHeader should be false")
	}
	if table.showEdge {
		t.Error("ShowEdge should be false")
	}
	if table.padding != 2 {
		t.Error("Padding not set")
	}
}

func TestAlignText(t *testing.T) {
	table := New()

	tests := []struct {
		text  string
		width int
		align Align
		check func(string) bool
	}{
		{"test", 10, AlignLeft, func(s string) bool {
			return strings.HasPrefix(s, "test") && len(s) == 10
		}},
		{"test", 10, AlignRight, func(s string) bool {
			return strings.HasSuffix(s, "test") && len(s) == 10
		}},
		{"test", 10, AlignCenter, func(s string) bool {
			return strings.Contains(s, "test") && len(s) == 10
		}},
	}

	for _, tt := range tests {
		result := table.alignText(tt.text, tt.width, tt.align)
		if !tt.check(result) {
			t.Errorf("alignText(%q, %d, %v) = %q failed check", tt.text, tt.width, tt.align, result)
		}
	}
}
