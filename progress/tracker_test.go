package progress

import (
	"testing"
	"time"
)

func TestTrackerSpeed(t *testing.T) {
	tracker := newTracker()

	// Simulate progress over time
	start := time.Now()
	tracker.startTime = start

	// Add samples at 100ms intervals
	for i := 0; i < 10; i++ {
		tracker.samples = append(tracker.samples, sample{
			timestamp: start.Add(time.Duration(i*100) * time.Millisecond),
			value:     int64(i * 10),
		})
	}

	// Speed should be approximately 100 units/second
	// (10 units per 100ms = 100 units/sec)
	speed := tracker.speed()
	if speed < 90 || speed > 110 {
		t.Errorf("Expected speed ~100, got %f", speed)
	}
}

func TestTrackerSpeedInsufficientData(t *testing.T) {
	tracker := newTracker()

	// No samples
	speed := tracker.speed()
	if speed != 0 {
		t.Errorf("Expected speed=0 with no samples, got %f", speed)
	}

	// One sample
	tracker.update(50)
	speed = tracker.speed()
	if speed != 0 {
		t.Errorf("Expected speed=0 with one sample, got %f", speed)
	}
}

func TestTrackerETA(t *testing.T) {
	tracker := newTracker()

	// Simulate steady progress
	start := time.Now()
	tracker.startTime = start

	// 50 units in 1 second = 50 units/sec
	tracker.samples = append(tracker.samples, sample{
		timestamp: start,
		value:     0,
	})
	tracker.samples = append(tracker.samples, sample{
		timestamp: start.Add(1 * time.Second),
		value:     50,
	})

	// Current: 50, Total: 100
	// Remaining: 50 units at 50 units/sec = 1 second
	eta := tracker.eta(50, 100)
	if eta < 900*time.Millisecond || eta > 1100*time.Millisecond {
		t.Errorf("Expected ETA ~1s, got %v", eta)
	}
}

func TestTrackerETAComplete(t *testing.T) {
	tracker := newTracker()

	// Already complete
	eta := tracker.eta(100, 100)
	if eta != 0 {
		t.Errorf("Expected ETA=0 when complete, got %v", eta)
	}

	// Over complete
	eta = tracker.eta(150, 100)
	if eta != 0 {
		t.Errorf("Expected ETA=0 when over complete, got %v", eta)
	}
}

func TestTrackerElapsed(t *testing.T) {
	tracker := newTracker()
	tracker.startTime = time.Now().Add(-5 * time.Second)

	elapsed := tracker.elapsed()
	if elapsed < 4*time.Second || elapsed > 6*time.Second {
		t.Errorf("Expected elapsed ~5s, got %v", elapsed)
	}
}

func TestTrackerReset(t *testing.T) {
	tracker := newTracker()

	// Add some samples
	tracker.update(10)
	tracker.update(20)
	tracker.update(30)

	if len(tracker.samples) != 3 {
		t.Errorf("Expected 3 samples before reset, got %d", len(tracker.samples))
	}

	// Reset
	tracker.Reset()

	if len(tracker.samples) != 0 {
		t.Errorf("Expected 0 samples after reset, got %d", len(tracker.samples))
	}

	// Start time should be recent
	if time.Since(tracker.startTime) > 100*time.Millisecond {
		t.Error("Expected startTime to be reset to now")
	}
}

func TestTrackerSampleLimit(t *testing.T) {
	tracker := newTracker()

	// Add more than 100 samples
	for i := 0; i < 150; i++ {
		tracker.update(int64(i))
	}

	// After 150 adds:
	// - First 100 samples: 0-99 (no pruning)
	// - 101st sample (value 100): triggers prune, keeps 50 most recent (50-100)
	// - Samples 102-150: adds 49 more (101-149)
	// - Final count: 50 + 49 = 99 samples
	if len(tracker.samples) != 99 {
		t.Errorf("Expected 99 samples after pruning, got %d", len(tracker.samples))
	}

	// Check that the kept samples are the most recent
	if tracker.samples[len(tracker.samples)-1].value != 149 {
		t.Errorf("Expected last sample value=149, got %d", tracker.samples[len(tracker.samples)-1].value)
	}
}
