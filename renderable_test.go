package rich

import (
	"testing"
)

func TestRenderableString(t *testing.T) {
	rs := NewRenderableString("Hello", NewStyle().Bold())

	console := NewConsole(nil)
	segments := rs.Render(console, 80)

	if len(segments) != 1 {
		t.Errorf("Expected 1 segment, got %d", len(segments))
	}

	if segments[0].Text != "Hello" {
		t.Errorf("Expected text 'Hello', got %q", segments[0].Text)
	}

	if !segments[0].Style.bold {
		t.Error("Expected bold style")
	}
}

func TestRenderableStringMeasure(t *testing.T) {
	rs := NewRenderableString("Hello", NewStyle())

	console := NewConsole(nil)
	measurement := rs.Measure(console, 100)

	if measurement.Minimum != 5 {
		t.Errorf("Expected minimum 5, got %d", measurement.Minimum)
	}

	if measurement.Maximum != 5 {
		t.Errorf("Expected maximum 5, got %d", measurement.Maximum)
	}
}

func TestLines(t *testing.T) {
	lines := Lines{
		NewRenderableString("Line 1", NewStyle()),
		NewRenderableString("Line 2", NewStyle()),
		NewRenderableString("Line 3", NewStyle()),
	}

	console := NewConsole(nil)
	segments := lines.Render(console, 80)

	output := segments.String()

	if output != "Line 1\nLine 2\nLine 3" {
		t.Errorf("Expected multiline output, got %q", output)
	}
}

func TestLinesEmpty(t *testing.T) {
	lines := Lines{}

	console := NewConsole(nil)
	segments := lines.Render(console, 80)

	if len(segments) != 0 {
		t.Errorf("Expected 0 segments for empty lines, got %d", len(segments))
	}
}

func TestLinesSingle(t *testing.T) {
	lines := Lines{
		NewRenderableString("Single", NewStyle()),
	}

	console := NewConsole(nil)
	segments := lines.Render(console, 80)

	// Should not have trailing newline for single line
	output := segments.String()
	if output != "Single" {
		t.Errorf("Expected 'Single', got %q", output)
	}
}
