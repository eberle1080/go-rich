package rich

// Renderable is the interface for objects that can be rendered to the console.
// Renderables convert themselves into a series of styled segments that can be
// displayed, taking into account the available width.
//
// This is the primary abstraction for creating custom rich content. Tables,
// panels, and other complex widgets implement this interface.
//
// Example implementation:
//
//	type MyWidget struct {}
//
//	func (w *MyWidget) Render(console *Console, width int) Segments {
//		return Segments{{Text: "Hello", Style: NewStyle().Bold()}}
//	}
type Renderable interface {
	// Render converts the renderable into segments for the given width.
	// The console parameter provides access to color mode and other console properties.
	// The width parameter indicates the maximum available width in characters.
	//
	// Implementations should respect the width constraint where possible,
	// wrapping or truncating content as appropriate.
	Render(console *Console, width int) Segments
}

// Measurable is the interface for renderables that can report their size requirements.
// Implementing this interface allows renderables to participate in automatic layout
// calculations, such as column sizing in tables.
//
// This is an optional interface - renderables that don't implement it will simply
// be rendered at the requested width without size negotiation.
type Measurable interface {
	// Measure returns the minimum and maximum width requirements.
	// The console parameter provides access to color mode and other properties.
	// The maxWidth parameter indicates the maximum width available for measurement.
	//
	// The returned Measurement describes the range of acceptable widths:
	//   - Minimum: The smallest width that can display the content without truncation
	//   - Maximum: The ideal width if space is unlimited
	Measure(console *Console, maxWidth int) Measurement
}

// RenderableString is a simple string that implements Renderable and Measurable.
// This is a basic implementation used internally, primarily for wrapping plain
// strings in a Renderable interface.
//
// For most use cases, you should use styled text or markup instead of creating
// RenderableString instances directly.
type RenderableString struct {
	Text  string // The text content
	Style Style  // The style to apply to the text
}

// NewRenderableString creates a new renderable string.
// This is typically used internally when converting strings to renderables.
//
// Example:
//
//	rs := NewRenderableString("Hello", NewStyle().Bold())
//	console.Render(rs)
func NewRenderableString(text string, style Style) *RenderableString {
	return &RenderableString{
		Text:  text,
		Style: style,
	}
}

// Render implements Renderable.
// Returns a single segment containing the text with the associated style.
// The width parameter is ignored since the text is returned as-is.
func (r *RenderableString) Render(console *Console, width int) Segments {
	return Segments{{Text: r.Text, Style: r.Style}}
}

// Measure implements Measurable.
// Returns the byte length of the text as both minimum and maximum width.
// Note: This uses byte length, not rune count, which may not be accurate
// for multi-byte Unicode characters. This is a known limitation of the
// current implementation.
func (r *RenderableString) Measure(console *Console, maxWidth int) Measurement {
	length := len(r.Text)
	return Measurement{
		Minimum: length,
		Maximum: length,
	}
}

// Lines is a renderable that represents multiple lines of content.
// Each line is itself a Renderable, allowing for complex multi-line layouts.
// Lines are automatically separated by newline characters during rendering.
//
// Example:
//
//	lines := Lines{
//		NewRenderableString("Line 1", NewStyle().Bold()),
//		NewRenderableString("Line 2", NewStyle().Italic()),
//		NewRenderableString("Line 3", NewStyle()),
//	}
//	console.Renderln(lines)
type Lines []Renderable

// Render implements Renderable.
// Renders each line in sequence, inserting newline segments between them.
// Each line is rendered with the full width available.
//
// The final line does not have a trailing newline - that should be added
// by the caller if needed (e.g., using Console.Renderln instead of Render).
func (l Lines) Render(console *Console, width int) Segments {
	var result Segments

	// Render each line
	for i, line := range l {
		// Add the line's segments
		result = append(result, line.Render(console, width)...)

		// Add newline between lines (but not after the last line)
		if i < len(l)-1 {
			result = append(result, Segment{Text: "\n", Style: NewStyle()})
		}
	}

	return result
}
