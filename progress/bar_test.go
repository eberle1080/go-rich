package progress

import (
	"testing"

	"github.com/eberle1080/go-rich"
)

func TestNewBar(t *testing.T) {
	bar := NewBar(100)

	if bar.total != 100 {
		t.Errorf("Expected total=100, got %d", bar.total)
	}

	if bar.current != 0 {
		t.Errorf("Expected current=0, got %d", bar.current)
	}

	if bar.completeChar != "█" {
		t.Errorf("Expected completeChar='█', got '%s'", bar.completeChar)
	}

	if bar.remainingChar != "░" {
		t.Errorf("Expected remainingChar='░', got '%s'", bar.remainingChar)
	}
}

func TestProgressBarSetProgress(t *testing.T) {
	bar := NewBar(100)

	// Normal case
	bar.SetProgress(50)
	if bar.current != 50 {
		t.Errorf("Expected current=50, got %d", bar.current)
	}

	// Clamping to max
	bar.SetProgress(150)
	if bar.current != 100 {
		t.Errorf("Expected current=100 (clamped), got %d", bar.current)
	}

	// Clamping to min
	bar.SetProgress(-10)
	if bar.current != 0 {
		t.Errorf("Expected current=0 (clamped), got %d", bar.current)
	}
}

func TestProgressBarAdvance(t *testing.T) {
	bar := NewBar(100)

	bar.SetProgress(30)
	bar.Advance(20)

	if bar.current != 50 {
		t.Errorf("Expected current=50, got %d", bar.current)
	}

	// Should clamp at total
	bar.Advance(100)
	if bar.current != 100 {
		t.Errorf("Expected current=100 (clamped), got %d", bar.current)
	}
}

func TestProgressBarPercentage(t *testing.T) {
	bar := NewBar(100)

	tests := []struct {
		progress int64
		expected float64
	}{
		{0, 0.0},
		{25, 0.25},
		{50, 0.5},
		{75, 0.75},
		{100, 1.0},
	}

	for _, tt := range tests {
		bar.SetProgress(tt.progress)
		pct := bar.Percentage()
		if pct != tt.expected {
			t.Errorf("SetProgress(%d): expected percentage=%f, got %f", tt.progress, tt.expected, pct)
		}
	}
}

func TestProgressBarIsComplete(t *testing.T) {
	bar := NewBar(100)

	if bar.IsComplete() {
		t.Error("Expected IsComplete()=false at start")
	}

	bar.SetProgress(50)
	if bar.IsComplete() {
		t.Error("Expected IsComplete()=false at 50%")
	}

	bar.SetProgress(100)
	if !bar.IsComplete() {
		t.Error("Expected IsComplete()=true at 100%")
	}
}

func TestProgressBarRender(t *testing.T) {
	console := rich.NewConsole(nil)
	bar := NewBar(100).Description("Test").Width(20)

	// Test at 0%
	bar.SetProgress(0)
	segments := bar.Render(console, 80)

	// Should have description + bar + percentage
	if len(segments) < 2 {
		t.Errorf("Expected at least 2 segments, got %d", len(segments))
	}

	// Check description
	if segments[0].Text != "Test " {
		t.Errorf("Expected first segment to be 'Test ', got '%s'", segments[0].Text)
	}

	// Test at 50%
	bar.SetProgress(50)
	segments = bar.Render(console, 80)

	plainText := segments.String()
	if plainText[:5] != "Test " {
		t.Errorf("Expected output to start with 'Test ', got '%s'", plainText)
	}
}

func TestProgressBarFluentAPI(t *testing.T) {
	bar := NewBar(100).
		Description("Download").
		Width(40).
		CompleteChar("=").
		RemainingChar("-").
		CompleteStyle(rich.NewStyle().Bold()).
		RemainingStyle(rich.NewStyle().Dim())

	if bar.description != "Download" {
		t.Errorf("Expected description='Download', got '%s'", bar.description)
	}

	if bar.width != 40 {
		t.Errorf("Expected width=40, got %d", bar.width)
	}

	if bar.completeChar != "=" {
		t.Errorf("Expected completeChar='=', got '%s'", bar.completeChar)
	}

	if bar.remainingChar != "-" {
		t.Errorf("Expected remainingChar='-', got '%s'", bar.remainingChar)
	}
}

func TestFormatPercentage(t *testing.T) {
	tests := []struct {
		value    float64
		expected string
	}{
		{0.0, "0%"},
		{0.1, "10%"},
		{0.25, "25%"},
		{0.5, "50%"},
		{0.755, "75.5%"},
		{0.999, "99.9%"},
		{1.0, "100%"},
	}

	for _, tt := range tests {
		result := formatPercentage(tt.value)
		if result != tt.expected {
			t.Errorf("formatPercentage(%f): expected '%s', got '%s'", tt.value, tt.expected, result)
		}
	}
}

func TestFormatInt(t *testing.T) {
	tests := []struct {
		value    int
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{10, "10"},
		{123, "123"},
		{-5, "-5"},
		{-123, "-123"},
	}

	for _, tt := range tests {
		result := formatInt(tt.value)
		if result != tt.expected {
			t.Errorf("formatInt(%d): expected '%s', got '%s'", tt.value, tt.expected, result)
		}
	}
}
