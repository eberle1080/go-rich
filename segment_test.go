package rich

import (
	"testing"
)

func TestSegments_String(t *testing.T) {
	segments := Segments{
		{Text: "Hello", Style: NewStyle()},
		{Text: " ", Style: NewStyle()},
		{Text: "World", Style: NewStyle()},
	}

	got := segments.String()
	want := "Hello World"

	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestSegments_Length(t *testing.T) {
	segments := Segments{
		{Text: "Hello", Style: NewStyle()},
		{Text: " ", Style: NewStyle()},
		{Text: "世界", Style: NewStyle()}, // 2 wide characters
	}

	got := segments.Length()
	want := 8 // "Hello " (6) + "世界" (2)

	if got != want {
		t.Errorf("Length() = %d, want %d", got, want)
	}
}

func TestSegments_ToANSI(t *testing.T) {
	segments := Segments{
		{Text: "Error", Style: NewStyle().Foreground(Red).Bold()},
		{Text: ": ", Style: NewStyle()},
		{Text: "failed", Style: NewStyle().Foreground(Red)},
	}

	// Test with no color
	plain := segments.ToANSI(ColorModeNone)
	if plain != "Error: failed" {
		t.Errorf("Plain text incorrect: %q", plain)
	}

	// Test with colors
	colored := segments.ToANSI(ColorModeStandard)
	if colored == plain {
		t.Error("Colored output should differ from plain")
	}
	if len(colored) <= len(plain) {
		t.Error("Colored output should be longer (ANSI codes)")
	}
}

func TestJoin(t *testing.T) {
	s1 := Segments{{Text: "a", Style: NewStyle()}}
	s2 := Segments{{Text: "b", Style: NewStyle()}}
	s3 := Segments{{Text: "c", Style: NewStyle()}}

	result := Join(s1, s2, s3)

	if len(result) != 3 {
		t.Errorf("Join length = %d, want 3", len(result))
	}

	if result.String() != "abc" {
		t.Errorf("Join string = %q, want %q", result.String(), "abc")
	}
}
