package rich

import (
	"testing"
)

func TestStyleImmutability(t *testing.T) {
	s1 := NewStyle()
	s2 := s1.Bold()
	s3 := s2.Foreground(Red)

	// Original should be unchanged
	if s1.bold {
		t.Error("Original style was modified")
	}

	// s2 should have bold but not color
	if !s2.bold {
		t.Error("s2 should be bold")
	}
	if s2.fg != nil {
		t.Error("s2 should not have color")
	}

	// s3 should have both
	if !s3.bold {
		t.Error("s3 should be bold")
	}
	if s3.fg == nil {
		t.Error("s3 should have color")
	}
}

func TestStyleComposition(t *testing.T) {
	style := NewStyle().
		Foreground(Red).
		Background(White).
		Bold().
		Italic().
		Underline()

	if style.fg == nil {
		t.Error("Foreground not set")
	}
	if style.bg == nil {
		t.Error("Background not set")
	}
	if !style.bold {
		t.Error("Bold not set")
	}
	if !style.italic {
		t.Error("Italic not set")
	}
	if !style.underline {
		t.Error("Underline not set")
	}
}

func TestStyleRender(t *testing.T) {
	style := NewStyle().Bold().Foreground(Red)
	st := style.Render("test")

	if st.Text != "test" {
		t.Errorf("Text = %q, want %q", st.Text, "test")
	}
	if !st.Style.bold {
		t.Error("Style lost bold attribute")
	}
}

func TestStyleToANSI(t *testing.T) {
	tests := []struct {
		name  string
		style Style
		mode  ColorMode
		want  string
	}{
		{
			name:  "none mode",
			style: NewStyle().Bold(),
			mode:  ColorModeNone,
			want:  "",
		},
		{
			name:  "bold only",
			style: NewStyle().Bold(),
			mode:  ColorModeStandard,
			want:  "\x1b[1m",
		},
		{
			name:  "italic only",
			style: NewStyle().Italic(),
			mode:  ColorModeStandard,
			want:  "\x1b[3m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.style.toANSI(tt.mode)
			if got != tt.want {
				t.Errorf("toANSI() = %q, want %q", got, tt.want)
			}
		})
	}
}
