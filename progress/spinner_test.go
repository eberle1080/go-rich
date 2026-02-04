package progress

import (
	"testing"
	"time"

	"github.com/eberle1080/go-rich"
)

func TestNewSpinner(t *testing.T) {
	spinner := NewSpinner(SpinnerDots)

	if len(spinner.frames) != len(SpinnerDots) {
		t.Errorf("Expected %d frames, got %d", len(SpinnerDots), len(spinner.frames))
	}

	if spinner.frameIndex != 0 {
		t.Errorf("Expected frameIndex=0, got %d", spinner.frameIndex)
	}

	if spinner.interval != 80*time.Millisecond {
		t.Errorf("Expected interval=80ms, got %v", spinner.interval)
	}
}

func TestNewSpinnerEmpty(t *testing.T) {
	spinner := NewSpinner([]string{})

	// Should default to SpinnerDots
	if len(spinner.frames) != len(SpinnerDots) {
		t.Errorf("Expected default SpinnerDots frames, got %d frames", len(spinner.frames))
	}
}

func TestSpinnerNext(t *testing.T) {
	spinner := NewSpinner([]string{"A", "B", "C"})

	if spinner.CurrentFrame() != "A" {
		t.Errorf("Expected current frame='A', got '%s'", spinner.CurrentFrame())
	}

	spinner.Next()
	if spinner.CurrentFrame() != "B" {
		t.Errorf("Expected current frame='B', got '%s'", spinner.CurrentFrame())
	}

	spinner.Next()
	if spinner.CurrentFrame() != "C" {
		t.Errorf("Expected current frame='C', got '%s'", spinner.CurrentFrame())
	}

	// Should wrap around
	spinner.Next()
	if spinner.CurrentFrame() != "A" {
		t.Errorf("Expected current frame='A' (wrapped), got '%s'", spinner.CurrentFrame())
	}
}

func TestSpinnerFluentAPI(t *testing.T) {
	spinner := NewSpinner(SpinnerLine).
		Description("Loading...").
		Style(rich.NewStyle().Bold()).
		Interval(100 * time.Millisecond)

	if spinner.description != "Loading..." {
		t.Errorf("Expected description='Loading...', got '%s'", spinner.description)
	}

	if spinner.interval != 100*time.Millisecond {
		t.Errorf("Expected interval=100ms, got %v", spinner.interval)
	}
}

func TestSpinnerRender(t *testing.T) {
	console := rich.NewConsole(nil)
	spinner := NewSpinner([]string{"|", "/"}).Description("Test")

	segments := spinner.Render(console, 80)

	// Should have frame + description
	if len(segments) < 2 {
		t.Errorf("Expected at least 2 segments, got %d", len(segments))
	}

	// First segment is the frame
	if segments[0].Text != "|" {
		t.Errorf("Expected first segment='|', got '%s'", segments[0].Text)
	}

	// Second segment is space + description
	if segments[1].Text != " Test" {
		t.Errorf("Expected second segment=' Test', got '%s'", segments[1].Text)
	}

	// After Next(), frame should change
	spinner.Next()
	segments = spinner.Render(console, 80)
	if segments[0].Text != "/" {
		t.Errorf("Expected first segment='/' after Next(), got '%s'", segments[0].Text)
	}
}

func TestSpinnerPredefinedStyles(t *testing.T) {
	styles := []struct {
		name   string
		frames []string
	}{
		{"SpinnerDots", SpinnerDots},
		{"SpinnerLine", SpinnerLine},
		{"SpinnerArc", SpinnerArc},
		{"SpinnerArrow", SpinnerArrow},
		{"SpinnerCircle", SpinnerCircle},
		{"SpinnerBounce", SpinnerBounce},
		{"SpinnerBoxBounce", SpinnerBoxBounce},
		{"SpinnerSimple", SpinnerSimple},
		{"SpinnerGrowVertical", SpinnerGrowVertical},
		{"SpinnerGrowHorizontal", SpinnerGrowHorizontal},
	}

	for _, style := range styles {
		if len(style.frames) == 0 {
			t.Errorf("%s has no frames", style.name)
		}

		spinner := NewSpinner(style.frames)
		if len(spinner.frames) != len(style.frames) {
			t.Errorf("%s: expected %d frames, got %d", style.name, len(style.frames), len(spinner.frames))
		}
	}
}
