package rich

// Measurement represents the size requirements of a renderable.
type Measurement struct {
	// Minimum is the minimum width needed to render without truncation.
	Minimum int
	// Maximum is the maximum width that would be used (before wrapping).
	Maximum int
}

// Clamp constrains the measurement to the given width.
func (m Measurement) Clamp(minWidth, maxWidth int) Measurement {
	min := m.Minimum
	max := m.Maximum

	if min < minWidth {
		min = minWidth
	}
	if max > maxWidth {
		max = maxWidth
	}
	if max < min {
		max = min
	}

	return Measurement{
		Minimum: min,
		Maximum: max,
	}
}

// Normalize ensures minimum <= maximum.
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
func (m Measurement) Add(other Measurement) Measurement {
	return Measurement{
		Minimum: m.Minimum + other.Minimum,
		Maximum: m.Maximum + other.Maximum,
	}
}

// Max returns the maximum of two measurements.
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
// Prefers the maximum if it fits within the available width.
func (m Measurement) Get(available int) int {
	if m.Maximum <= available {
		return m.Maximum
	}
	if m.Minimum > available {
		return m.Minimum
	}
	return available
}
