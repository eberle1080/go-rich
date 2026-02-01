package rich

import (
	"testing"
)

func TestHex(t *testing.T) {
	tests := []struct {
		input    string
		expected RGBColor
		wantErr  bool
	}{
		{"#FF0000", RGBColor{255, 0, 0}, false},
		{"FF0000", RGBColor{255, 0, 0}, false},
		{"#00FF00", RGBColor{0, 255, 0}, false},
		{"#0000FF", RGBColor{0, 0, 255}, false},
		{"#FF1493", RGBColor{255, 20, 147}, false},
		{"#ZZZZZZ", RGBColor{}, true},
		{"#FF", RGBColor{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := Hex(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hex(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.expected {
				t.Errorf("Hex(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestNamed(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"red", false},
		{"blue", false},
		{"green", false},
		{"orange", false},
		{"invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Named(tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Named(%q) error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}

func TestANSIColor_toANSI(t *testing.T) {
	tests := []struct {
		color      ANSIColor
		mode       ColorMode
		foreground bool
		want       string
	}{
		{Red, ColorModeNone, true, ""},
		{Red, ColorModeStandard, true, "\x1b[31m"},
		{Green, ColorModeStandard, true, "\x1b[32m"},
		{Blue, ColorModeStandard, false, "\x1b[44m"},
		{BrightRed, ColorModeStandard, true, "\x1b[91m"},
		{BrightRed, ColorModeStandard, false, "\x1b[101m"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := tt.color.toANSI(tt.mode, tt.foreground)
			if got != tt.want {
				t.Errorf("toANSI() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRGBColor_toANSI(t *testing.T) {
	rgb := RGBColor{255, 100, 50}

	tests := []struct {
		mode       ColorMode
		foreground bool
		check      func(string) bool
	}{
		{ColorModeNone, true, func(s string) bool { return s == "" }},
		{ColorModeTrueColor, true, func(s string) bool { return s == "\x1b[38;2;255;100;50m" }},
		{ColorMode256, true, func(s string) bool { return len(s) > 0 && s != "" }},
		{ColorModeStandard, true, func(s string) bool { return len(s) > 0 && s != "" }},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := rgb.toANSI(tt.mode, tt.foreground)
			if !tt.check(got) {
				t.Errorf("toANSI() = %q, check failed", got)
			}
		})
	}
}
