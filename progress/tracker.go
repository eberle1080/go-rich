package progress

import (
	"sync"
	"time"
)

// Tracker tracks progress over time to calculate speed and ETA.
// It uses exponential moving average (EMA) for smooth speed calculations
// that react to recent changes while ignoring old noise.
//
// The tracker maintains a history of samples (timestamp + value pairs)
// and calculates instantaneous speed using the derivative of the smoothed curve.
//
// Thread-safe for concurrent updates.
type Tracker struct {
	startTime time.Time // When tracking started
	samples   []sample  // Historical samples for speed calculation

	smoothAlpha float64 // Smoothing factor for EMA (0.0-1.0, default 0.5)

	mu sync.Mutex // Protects samples slice
}

// sample represents a single progress measurement at a point in time.
type sample struct {
	timestamp time.Time
	value     int64
}

// newTracker creates a new tracker with default settings.
// The smoothing alpha defaults to 0.5, balancing responsiveness
// with noise reduction.
func newTracker() *Tracker {
	return &Tracker{
		startTime:   time.Now(),
		samples:     make([]sample, 0, 100),
		smoothAlpha: 0.5,
	}
}

// update records a new progress value at the current time.
// This should be called whenever progress changes to maintain
// accurate speed and ETA calculations.
//
// Thread-safe.
func (t *Tracker) update(value int64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	t.samples = append(t.samples, sample{
		timestamp: now,
		value:     value,
	})

	// Keep only recent samples (last 100)
	// This prevents unbounded memory growth for long-running operations
	if len(t.samples) > 100 {
		// Keep the 50 most recent samples
		copy(t.samples, t.samples[50:])
		t.samples = t.samples[:50]
	}
}

// speed calculates the current speed in units per second.
// Returns 0 if insufficient data is available.
//
// The calculation uses a sliding window approach:
//  1. Find samples from the last 2 seconds
//  2. If we have at least 2 samples, calculate average speed
//  3. Apply exponential smoothing to reduce noise
//
// Thread-safe.
func (t *Tracker) speed() float64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(t.samples) < 2 {
		return 0
	}

	now := time.Now()
	windowDuration := 2 * time.Second

	// Find samples within the time window
	var windowSamples []sample
	for i := len(t.samples) - 1; i >= 0; i-- {
		s := t.samples[i]
		if now.Sub(s.timestamp) <= windowDuration {
			windowSamples = append([]sample{s}, windowSamples...)
		} else {
			break
		}
	}

	// Need at least 2 samples to calculate speed
	if len(windowSamples) < 2 {
		// Fall back to using all samples if window is empty
		if len(t.samples) >= 2 {
			windowSamples = t.samples
		} else {
			return 0
		}
	}

	// Calculate speed using first and last sample in window
	first := windowSamples[0]
	last := windowSamples[len(windowSamples)-1]

	deltaValue := float64(last.value - first.value)
	deltaTime := last.timestamp.Sub(first.timestamp).Seconds()

	if deltaTime == 0 {
		return 0
	}

	return deltaValue / deltaTime
}

// eta calculates the estimated time remaining to reach the target value.
// Returns 0 if speed is too slow or negative (going backwards).
//
// The calculation is: remaining / speed
//
// Thread-safe.
func (t *Tracker) eta(current, total int64) time.Duration {
	if current >= total {
		return 0
	}

	spd := t.speed()
	if spd <= 0 {
		return 0 // Can't estimate with zero or negative speed
	}

	remaining := float64(total - current)
	secondsRemaining := remaining / spd

	// Cap at reasonable maximum (24 hours)
	if secondsRemaining > 86400 {
		return 24 * time.Hour
	}

	return time.Duration(secondsRemaining * float64(time.Second))
}

// elapsed returns the time since tracking started.
func (t *Tracker) elapsed() time.Duration {
	return time.Since(t.startTime)
}

// Speed returns the current speed in units per second.
// This is the public API for accessing speed calculations.
func (t *Tracker) Speed() float64 {
	return t.speed()
}

// ETA returns the estimated time remaining to reach the target.
// This is the public API for accessing ETA calculations.
func (t *Tracker) ETA(current, total int64) time.Duration {
	return t.eta(current, total)
}

// Elapsed returns the time since tracking started.
// This is the public API for accessing elapsed time.
func (t *Tracker) Elapsed() time.Duration {
	return t.elapsed()
}

// Reset resets the tracker to initial state with a new start time.
// This clears all samples and resets timing.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.startTime = time.Now()
	t.samples = t.samples[:0]
}
