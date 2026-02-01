package rich

import "testing"

func TestMeasurementClamp(t *testing.T) {
	m := Measurement{Minimum: 10, Maximum: 50}

	// Clamp within range
	clamped := m.Clamp(0, 100)
	if clamped.Minimum != 10 || clamped.Maximum != 50 {
		t.Errorf("Clamp(0, 100) = %+v, want {10, 50}", clamped)
	}

	// Clamp with lower minimum
	clamped = m.Clamp(20, 100)
	if clamped.Minimum != 20 {
		t.Errorf("Clamp(20, 100).Minimum = %d, want 20", clamped.Minimum)
	}

	// Clamp with upper maximum
	clamped = m.Clamp(0, 30)
	if clamped.Maximum != 30 {
		t.Errorf("Clamp(0, 30).Maximum = %d, want 30", clamped.Maximum)
	}
}

func TestMeasurementNormalize(t *testing.T) {
	// Already normalized
	m := Measurement{Minimum: 10, Maximum: 50}
	normalized := m.Normalize()
	if normalized.Minimum != 10 || normalized.Maximum != 50 {
		t.Error("Normalize should not change already normalized measurement")
	}

	// Needs normalization
	m = Measurement{Minimum: 50, Maximum: 10}
	normalized = m.Normalize()
	if normalized.Minimum != 50 || normalized.Maximum != 50 {
		t.Errorf("Normalize() = %+v, want {50, 50}", normalized)
	}
}

func TestMeasurementAdd(t *testing.T) {
	m1 := Measurement{Minimum: 10, Maximum: 20}
	m2 := Measurement{Minimum: 5, Maximum: 15}

	result := m1.Add(m2)

	if result.Minimum != 15 {
		t.Errorf("Add().Minimum = %d, want 15", result.Minimum)
	}
	if result.Maximum != 35 {
		t.Errorf("Add().Maximum = %d, want 35", result.Maximum)
	}
}

func TestMeasurementMax(t *testing.T) {
	m1 := Measurement{Minimum: 10, Maximum: 20}
	m2 := Measurement{Minimum: 15, Maximum: 25}

	result := m1.Max(m2)

	if result.Minimum != 15 {
		t.Errorf("Max().Minimum = %d, want 15", result.Minimum)
	}
	if result.Maximum != 25 {
		t.Errorf("Max().Maximum = %d, want 25", result.Maximum)
	}
}

func TestMeasurementGet(t *testing.T) {
	m := Measurement{Minimum: 10, Maximum: 50}

	// Available space fits maximum
	if got := m.Get(100); got != 50 {
		t.Errorf("Get(100) = %d, want 50", got)
	}

	// Available space between min and max
	if got := m.Get(30); got != 30 {
		t.Errorf("Get(30) = %d, want 30", got)
	}

	// Available space less than minimum
	if got := m.Get(5); got != 10 {
		t.Errorf("Get(5) = %d, want 10 (minimum)", got)
	}
}
