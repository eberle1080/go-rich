package rich

import (
	"bytes"
	"strings"
	"testing"
)

func TestConsoleBasicOutput(t *testing.T) {
	var buf bytes.Buffer
	console := NewConsole(&buf)
	console.SetColorMode(ColorModeNone) // Disable colors for predictable output

	console.Print("Hello")
	console.Println(" World")

	got := buf.String()
	want := "Hello World\n"

	if got != want {
		t.Errorf("Output = %q, want %q", got, want)
	}
}

func TestConsolePrintf(t *testing.T) {
	var buf bytes.Buffer
	console := NewConsole(&buf)
	console.SetColorMode(ColorModeNone)

	console.Printf("Hello %s, you are %d years old", "Alice", 30)

	got := buf.String()
	want := "Hello Alice, you are 30 years old"

	if got != want {
		t.Errorf("Output = %q, want %q", got, want)
	}
}

func TestConsolePrintStyled(t *testing.T) {
	var buf bytes.Buffer
	console := NewConsole(&buf)
	console.SetColorMode(ColorModeStandard)

	style := NewStyle().Foreground(Red).Bold()
	console.PrintStyledln(style.Render("Error"))

	got := buf.String()

	// Should contain ANSI codes
	if !strings.Contains(got, "\x1b[") {
		t.Error("Output should contain ANSI escape codes")
	}

	// Should contain the text
	if !strings.Contains(got, "Error") {
		t.Error("Output should contain 'Error'")
	}
}

func TestConsoleColorModeDetection(t *testing.T) {
	var buf bytes.Buffer
	console := NewConsole(&buf)

	// Should detect some mode (likely None for buffer)
	mode := console.ColorMode()
	if mode < ColorModeNone || mode > ColorModeTrueColor {
		t.Errorf("Invalid color mode: %d", mode)
	}
}

func TestConsoleSetColorMode(t *testing.T) {
	var buf bytes.Buffer
	console := NewConsole(&buf)

	console.SetColorMode(ColorMode256)
	if console.ColorMode() != ColorMode256 {
		t.Error("SetColorMode did not update mode")
	}
}

func TestConsoleRule(t *testing.T) {
	var buf bytes.Buffer
	console := NewConsole(&buf)
	console.SetColorMode(ColorModeNone)

	console.Rule("Test")

	got := buf.String()

	// Should contain title
	if !strings.Contains(got, "Test") {
		t.Error("Rule should contain title")
	}

	// Should contain horizontal lines
	if !strings.Contains(got, "─") {
		t.Error("Rule should contain horizontal line characters")
	}
}

func TestConsoleRuleEmpty(t *testing.T) {
	var buf bytes.Buffer
	console := NewConsole(&buf)
	console.SetColorMode(ColorModeNone)

	console.Rule("")

	got := buf.String()

	// Should just be a line
	if !strings.Contains(got, "─") {
		t.Error("Empty rule should contain horizontal line characters")
	}

	// Should be console width
	lines := strings.Split(strings.TrimSpace(got), "\n")
	if len(lines) != 1 {
		t.Errorf("Expected 1 line, got %d", len(lines))
	}
}

func TestConsoleWidth(t *testing.T) {
	var buf bytes.Buffer
	console := NewConsole(&buf)

	width := console.Width()
	if width <= 0 {
		t.Errorf("Width should be positive, got %d", width)
	}
}

func TestConsoleHeight(t *testing.T) {
	var buf bytes.Buffer
	console := NewConsole(&buf)

	height := console.Height()
	if height <= 0 {
		t.Errorf("Height should be positive, got %d", height)
	}
}
