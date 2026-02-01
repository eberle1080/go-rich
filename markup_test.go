package rich

import (
	"strings"
	"testing"
)

func TestStripMarkup(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"[bold]text[/]", "text"},
		{"[red]Hello[/] World", "Hello World"},
		{"[[escaped]", "[escaped]"},
		{"plain text", "plain text"},
		{"[bold red]Error:[/] message", "Error: message"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := StripMarkup(tt.input)
			if got != tt.want {
				t.Errorf("StripMarkup(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestEscapeMarkup(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"text", "text"},
		{"[bold]", "[[bold]"},
		{"a[b]c", "a[[b]c"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := EscapeMarkup(tt.input)
			if got != tt.want {
				t.Errorf("EscapeMarkup(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateMarkup(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"[bold]text[/]", false},
		{"[red]Hello[/]", false},
		{"plain text", false},
		{"[bold][italic]text[/][/]", false},
		{"[bold]unclosed", true},
		{"[/]unmatched", true},
		{"[bold][/][/]", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			err := ValidateMarkup(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMarkup(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestParseMarkup(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(Segments) bool
	}{
		{
			name:  "plain text",
			input: "Hello World",
			check: func(s Segments) bool {
				return len(s) == 1 && s[0].Text == "Hello World"
			},
		},
		{
			name:  "bold text",
			input: "[bold]Hello[/]",
			check: func(s Segments) bool {
				return len(s) == 1 && s[0].Text == "Hello" && s[0].Style.bold
			},
		},
		{
			name:  "colored text",
			input: "[red]Error[/]",
			check: func(s Segments) bool {
				return len(s) == 1 && s[0].Text == "Error" && s[0].Style.fg != nil
			},
		},
		{
			name:  "multiple segments",
			input: "[bold]Hello[/] World",
			check: func(s Segments) bool {
				return len(s) == 2 && s[0].Text == "Hello" && s[1].Text == " World"
			},
		},
		{
			name:  "nested styles",
			input: "[bold][red]Error[/][/]",
			check: func(s Segments) bool {
				return len(s) == 1 && s[0].Text == "Error" && s[0].Style.bold && s[0].Style.fg != nil
			},
		},
		{
			name:  "escaped brackets",
			input: "[[not a tag]",
			check: func(s Segments) bool {
				return s.String() == "[not a tag]"
			},
		},
		{
			name:  "combined styles",
			input: "[bold red]Error:[/] message",
			check: func(s Segments) bool {
				return len(s) == 2 && s[0].Text == "Error:" && s[0].Style.bold && s[0].Style.fg != nil
			},
		},
		{
			name:  "background color",
			input: "[red on blue]text[/]",
			check: func(s Segments) bool {
				return len(s) == 1 && s[0].Style.fg != nil && s[0].Style.bg != nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segments, err := parseMarkup(tt.input)
			if err != nil {
				t.Errorf("parseMarkup(%q) error = %v", tt.input, err)
				return
			}
			if !tt.check(segments) {
				t.Errorf("parseMarkup(%q) check failed, got segments: %+v", tt.input, segments)
			}
		})
	}
}

func TestConsolePrintMarkup(t *testing.T) {
	var buf strings.Builder
	console := NewConsole(&buf)
	console.SetColorMode(ColorModeNone)

	console.PrintMarkup("[bold]Hello[/] World")

	got := buf.String()
	want := "Hello World"

	if got != want {
		t.Errorf("PrintMarkup output = %q, want %q", got, want)
	}
}

func TestConsolePrintMarkupln(t *testing.T) {
	var buf strings.Builder
	console := NewConsole(&buf)
	console.SetColorMode(ColorModeNone)

	console.PrintMarkupln("[bold]Hello[/]")

	got := buf.String()
	want := "Hello\n"

	if got != want {
		t.Errorf("PrintMarkupln output = %q, want %q", got, want)
	}
}

func TestConsolePrintMarkupWithColors(t *testing.T) {
	var buf strings.Builder
	console := NewConsole(&buf)
	console.SetColorMode(ColorModeStandard)

	console.PrintMarkup("[red]Error[/]")

	got := buf.String()

	// Should contain ANSI codes and the text
	if !strings.Contains(got, "Error") {
		t.Error("Output should contain 'Error'")
	}

	if !strings.Contains(got, "\x1b[") {
		t.Error("Output should contain ANSI escape codes")
	}
}

func TestMarkupColorParsing(t *testing.T) {
	tests := []struct {
		markup string
		valid  bool
	}{
		{"[red]text[/]", true},
		{"[#FF0000]text[/]", true},
		{"[rgb(255,0,0)]text[/]", true},
		{"[bright_red]text[/]", true},
		{"[orange]text[/]", true},
	}

	for _, tt := range tests {
		t.Run(tt.markup, func(t *testing.T) {
			segments, err := parseMarkup(tt.markup)
			if err != nil {
				t.Errorf("parseMarkup(%q) error = %v", tt.markup, err)
				return
			}
			if len(segments) == 0 {
				t.Error("Expected at least one segment")
				return
			}
			// Just check that it parsed without error
			if tt.valid && segments[0].Style.fg == nil {
				t.Errorf("Expected color to be set for %q", tt.markup)
			}
		})
	}
}
