package rich

// Renderable is the interface for objects that can be rendered to the console.
// Renderables convert themselves into a series of styled segments that can be
// displayed, taking into account the available width.
type Renderable interface {
	// Render converts the renderable into segments for the given width.
	// The width parameter indicates the maximum available width in characters.
	Render(console *Console, width int) Segments
}

// Measurable is the interface for renderables that can report their size requirements.
type Measurable interface {
	// Measure returns the minimum and maximum width requirements.
	Measure(console *Console, maxWidth int) Measurement
}

// RenderableString is a simple string that implements Renderable.
type RenderableString struct {
	Text  string
	Style Style
}

// NewRenderableString creates a new renderable string.
func NewRenderableString(text string, style Style) *RenderableString {
	return &RenderableString{
		Text:  text,
		Style: style,
	}
}

// Render implements Renderable.
func (r *RenderableString) Render(console *Console, width int) Segments {
	return Segments{{Text: r.Text, Style: r.Style}}
}

// Measure implements Measurable.
func (r *RenderableString) Measure(console *Console, maxWidth int) Measurement {
	length := len(r.Text)
	return Measurement{
		Minimum: length,
		Maximum: length,
	}
}

// Lines is a renderable that represents multiple lines.
type Lines []Renderable

// Render implements Renderable.
func (l Lines) Render(console *Console, width int) Segments {
	var result Segments
	for i, line := range l {
		result = append(result, line.Render(console, width)...)
		if i < len(l)-1 {
			result = append(result, Segment{Text: "\n", Style: NewStyle()})
		}
	}
	return result
}
