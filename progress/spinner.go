package progress

import (
	"time"

	"github.com/eberle1080/go-rich"
)

// Spinner represents an animated spinner for indeterminate progress.
// Unlike a progress bar, spinners are used when you don't know how long
// an operation will take or can't measure progress.
//
// Spinners cycle through a series of frames at a fixed interval,
// creating an animation effect. Several built-in spinner styles are provided,
// or you can create custom spinners with your own frames.
//
// Example:
//
//	spinner := progress.NewSpinner(progress.SpinnerDots).
//		Description("Loading...").
//		Style(rich.NewStyle().Foreground(rich.Cyan))
//	console.Render(spinner)
type Spinner struct {
	frames      []string      // Animation frames to cycle through
	frameIndex  int           // Current frame index
	interval    time.Duration // Time between frame updates
	style       rich.Style    // Style applied to the spinner
	description string        // Text description displayed with spinner
}

// Predefined spinner styles
var (
	// SpinnerDots is a Braille-pattern dots spinner: ⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏
	SpinnerDots = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

	// SpinnerLine is a simple line spinner: -\|/
	SpinnerLine = []string{"-", "\\", "|", "/"}

	// SpinnerArc is an arc spinner: ◜◠◝◞◡◟
	SpinnerArc = []string{"◜", "◠", "◝", "◞", "◡", "◟"}

	// SpinnerArrow is an arrow spinner rotating clockwise: ←↖↑↗→↘↓↙
	SpinnerArrow = []string{"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"}

	// SpinnerCircle is a circle quadrant spinner: ◐◓◑◒
	SpinnerCircle = []string{"◐", "◓", "◑", "◒"}

	// SpinnerBounce is a bouncing ball: ⠁⠂⠄⡀⢀⠠⠐⠈
	SpinnerBounce = []string{"⠁", "⠂", "⠄", "⡀", "⢀", "⠠", "⠐", "⠈"}

	// SpinnerBoxBounce is a bouncing box animation
	SpinnerBoxBounce = []string{"▖", "▘", "▝", "▗"}

	// SpinnerSimple is a simple ASCII spinner: |/-\
	SpinnerSimple = []string{"|", "/", "-", "\\"}

	// SpinnerGrowVertical grows vertically: ▁▃▄▅▆▇▆▅▄▃
	SpinnerGrowVertical = []string{"▁", "▃", "▄", "▅", "▆", "▇", "▆", "▅", "▄", "▃"}

	// SpinnerGrowHorizontal grows horizontally: ▏▎▍▌▋▊▉▊▋▌▍▎
	SpinnerGrowHorizontal = []string{"▏", "▎", "▍", "▌", "▋", "▊", "▉", "▊", "▋", "▌", "▍", "▎"}
)

// NewSpinner creates a new spinner with the specified frame sequence.
// Use one of the predefined spinner styles (SpinnerDots, SpinnerLine, etc.)
// or provide your own custom frames.
//
// The default interval between frames is 80ms (12.5 FPS).
//
// Example:
//
//	spinner := progress.NewSpinner(progress.SpinnerDots)
//
// Custom frames:
//
//	spinner := progress.NewSpinner([]string{".", "..", "..."})
func NewSpinner(frames []string) *Spinner {
	if len(frames) == 0 {
		frames = SpinnerDots // Default to dots if empty
	}

	return &Spinner{
		frames:      frames,
		frameIndex:  0,
		interval:    80 * time.Millisecond,
		style:       rich.NewStyle(),
		description: "",
	}
}

// Description sets the text description displayed with the spinner.
// This typically appears after the spinner animation.
//
// Example:
//
//	spinner := progress.NewSpinner(progress.SpinnerDots).
//		Description("Loading data...")
func (s *Spinner) Description(desc string) *Spinner {
	s.description = desc
	return s
}

// Style sets the style for the spinner frames.
// This affects the color and formatting of the animated character.
//
// Example:
//
//	spinner := progress.NewSpinner(progress.SpinnerDots).
//		Style(rich.NewStyle().Foreground(rich.Cyan).Bold())
func (s *Spinner) Style(style rich.Style) *Spinner {
	s.style = style
	return s
}

// Interval sets the time between frame updates.
// Lower values create faster animation, higher values create slower animation.
// Default is 80ms (12.5 FPS).
//
// Example:
//
//	spinner := progress.NewSpinner(progress.SpinnerDots).
//		Interval(100 * time.Millisecond) // Slower animation
func (s *Spinner) Interval(interval time.Duration) *Spinner {
	s.interval = interval
	return s
}

// Next advances to the next frame in the animation.
// This should be called periodically (typically by a Progress manager)
// to animate the spinner.
//
// The method wraps around to the first frame after reaching the end.
func (s *Spinner) Next() {
	s.frameIndex = (s.frameIndex + 1) % len(s.frames)
}

// CurrentFrame returns the current animation frame character.
func (s *Spinner) CurrentFrame() string {
	return s.frames[s.frameIndex]
}

// Render implements rich.Renderable.
// Converts the spinner into styled segments for display.
//
// The spinner is rendered as:
//
//	[frame] [description]
func (s *Spinner) Render(console *rich.Console, width int) rich.Segments {
	segments := rich.Segments{}

	// Render current frame
	segments = append(segments, rich.Segment{
		Text:  s.CurrentFrame(),
		Style: s.style,
	})

	// Render description if present
	if s.description != "" {
		segments = append(segments, rich.Segment{
			Text:  " " + s.description,
			Style: rich.NewStyle(),
		})
	}

	return segments
}

// Measure implements rich.Measurable.
// Returns the size requirements for the spinner.
func (s *Spinner) Measure(console *rich.Console, maxWidth int) rich.Measurement {
	// Frame is typically 1-2 characters, description is variable
	frameLen := 2 // Assume max 2 chars for Unicode spinners
	descLen := len(s.description)
	if descLen > 0 {
		descLen++ // Space before description
	}

	size := frameLen + descLen
	return rich.Measurement{
		Minimum: size,
		Maximum: size,
	}
}
