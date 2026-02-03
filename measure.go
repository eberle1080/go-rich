package rich

// Measurement represents the size requirements of a renderable.
// This is used by the layout system to determine optimal sizing for
// renderables, especially in contexts like table columns.
//
// A measurement defines a range of acceptable widths:
//   - Minimum: The smallest width that displays content without truncation
//   - Maximum: The preferred width if unlimited space is available
//
// For example, a text string "Hello" might report:
//   - Minimum: 5 (to show the whole word)
//   - Maximum: 5 (it doesn't need more space)
//
// A table column with cells "A", "BB", "CCC" might report:
//   - Minimum: 3 (to fit the longest cell)
//   - Maximum: 3 (no benefit to being wider)
type Measurement struct {
	// Minimum is the minimum width needed to render without truncation.
	// Content rendered at less than this width may be cut off or wrapped.
	Minimum int

	// Maximum is the maximum width that would be used (before wrapping).
	// This is the "natural" or "preferred" width of the content.
	// Providing more width than this typically doesn't improve the rendering.
	Maximum int
}

// Clamp constrains the measurement to the given width range.
// This adjusts both Minimum and Maximum to fit within [minWidth, maxWidth],
// ensuring the result is valid (Minimum <= Maximum).
//
// Use cases:
//   - Enforcing minimum column widths in tables
//   - Limiting maximum widths to fit within terminal bounds
//
// Example:
//
//	m := Measurement{Minimum: 5, Maximum: 100}
//	clamped := m.Clamp(10, 50)
//	// Result: {Minimum: 10, Maximum: 50}
func (m Measurement) Clamp(minWidth, maxWidth int) Measurement {
	min := m.Minimum
	max := m.Maximum

	// Raise Minimum if it's below the constraint
	if min < minWidth {
		min = minWidth
	}

	// Lower Maximum if it exceeds the constraint
	if max > maxWidth {
		max = maxWidth
	}

	// Ensure Maximum is at least Minimum
	if max < min {
		max = min
	}

	return Measurement{
		Minimum: min,
		Maximum: max,
	}
}

// Normalize ensures minimum <= maximum.
// If Maximum < Minimum, sets Maximum = Minimum.
// This is useful for correcting invalid measurements.
//
// Example:
//
//	m := Measurement{Minimum: 100, Maximum: 50}
//	normalized := m.Normalize()
//	// Result: {Minimum: 100, Maximum: 100}
func (m Measurement) Normalize() Measurement {
	if m.Maximum < m.Minimum {
		return Measurement{
			Minimum: m.Minimum,
			Maximum: m.Minimum,
		}
	}
	return m
}

// Add adds two measurements together.
// Both Minimum and Maximum values are summed.
// This is useful for calculating the combined width of multiple elements.
//
// Example:
//
//	m1 := Measurement{Minimum: 5, Maximum: 10}
//	m2 := Measurement{Minimum: 3, Maximum: 8}
//	sum := m1.Add(m2)
//	// Result: {Minimum: 8, Maximum: 18}
func (m Measurement) Add(other Measurement) Measurement {
	return Measurement{
		Minimum: m.Minimum + other.Minimum,
		Maximum: m.Maximum + other.Maximum,
	}
}

// Max returns the maximum of two measurements.
// Takes the larger Minimum and the larger Maximum from both measurements.
// This is useful for finding the size requirements of the largest element.
//
// Example:
//
//	m1 := Measurement{Minimum: 5, Maximum: 10}
//	m2 := Measurement{Minimum: 3, Maximum: 15}
//	max := m1.Max(m2)
//	// Result: {Minimum: 5, Maximum: 15}
func (m Measurement) Max(other Measurement) Measurement {
	min := m.Minimum
	if other.Minimum > min {
		min = other.Minimum
	}

	max := m.Maximum
	if other.Maximum > max {
		max = other.Maximum
	}

	return Measurement{
		Minimum: min,
		Maximum: max,
	}
}

// Get returns a width within the measurement range.
// The selection logic:
//   - If Maximum fits within available width, return Maximum (preferred width)
//   - If Minimum exceeds available width, return Minimum (can't go smaller)
//   - Otherwise, return available (best fit within constraints)
//
// This is used by layout algorithms to determine actual rendering widths.
//
// Example:
//
//	m := Measurement{Minimum: 10, Maximum: 50}
//	m.Get(100) // Returns: 50 (Maximum fits)
//	m.Get(30)  // Returns: 30 (between Minimum and Maximum)
//	m.Get(5)   // Returns: 10 (below Minimum, use Minimum)
func (m Measurement) Get(available int) int {
	// Prefer Maximum if it fits
	if m.Maximum <= available {
		return m.Maximum
	}

	// Can't go below Minimum even if space is tight
	if m.Minimum > available {
		return m.Minimum
	}

	// Use all available space (between Minimum and Maximum)
	return available
}
