package progress

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/eberle1080/go-rich"
)

// ANSI control sequences for cursor manipulation
const (
	cursorUp   = "\x1b[%dA"   // Move cursor up N lines
	cursorLeft = "\x1b[1000D" // Move cursor to line start (move left 1000 chars)
	clearLine  = "\x1b[0K"    // Clear from cursor to end of line
	hideCursor = "\x1b[?25l"  // Hide cursor
	showCursor = "\x1b[?25h"  // Show cursor
)

// TaskID identifies a task in the progress manager.
type TaskID int

// Task represents a single progress task being tracked.
type Task struct {
	id        TaskID
	bar       *ProgressBar // Progress bar (nil for spinners)
	spinner   *Spinner     // Spinner (nil for bars)
	startTime time.Time    // When the task started
	completed bool         // Whether the task is complete
}

// Progress manages live progress updates for multiple tasks.
// It handles the rendering loop and ANSI cursor control to update
// progress bars and spinners in place without scrolling.
//
// Example:
//
//	prog := progress.New(console)
//	task1 := prog.AddBar("Download", 1000)
//	task2 := prog.AddSpinner("Processing...")
//	prog.Start()
//
//	// Update from goroutines
//	for i := 0; i < 1000; i++ {
//		prog.Update(task1, i)
//		time.Sleep(10 * time.Millisecond)
//	}
//
//	prog.Complete(task1)
//	prog.Stop()
type Progress struct {
	console *rich.Console // Console for rendering
	writer  io.Writer     // Direct writer (usually console.Writer())

	tasks   map[TaskID]*Task // Active tasks
	taskSeq TaskID           // Task ID sequence
	mu      sync.RWMutex     // Protects tasks map

	running     bool          // Whether the render loop is running
	ticker      *time.Ticker  // Ticker for periodic refresh
	refreshRate time.Duration // Time between refreshes
	stopChan    chan struct{} // Channel to signal stop

	transient     bool // Whether to clear progress on completion
	lastLineCount int  // Number of lines rendered in last update
}

// New creates a new progress manager.
// The console is used for rendering and color mode detection.
//
// Default settings:
//   - Refresh rate: 100ms (10 FPS)
//   - Transient: false (keep progress visible after completion)
//
// Example:
//
//	prog := progress.New(console).RefreshRate(50 * time.Millisecond)
func New(console *rich.Console) *Progress {
	return &Progress{
		console:     console,
		writer:      console.Writer(),
		tasks:       make(map[TaskID]*Task),
		taskSeq:     0,
		refreshRate: 100 * time.Millisecond,
		stopChan:    make(chan struct{}),
		transient:   false,
	}
}

// RefreshRate sets the time between progress updates.
// Lower values create smoother animation but use more CPU.
// Higher values reduce CPU usage but may appear choppy.
//
// Default is 100ms (10 FPS).
//
// Example:
//
//	prog := progress.New(console).RefreshRate(50 * time.Millisecond)
func (p *Progress) RefreshRate(rate time.Duration) *Progress {
	p.refreshRate = rate
	return p
}

// Transient sets whether to clear the progress display when stopped.
// If true, all progress bars/spinners are erased when Stop() is called.
// If false (default), they remain visible as the final state.
//
// Example:
//
//	prog := progress.New(console).Transient(true)
func (p *Progress) Transient(transient bool) *Progress {
	p.transient = transient
	return p
}

// AddBar adds a progress bar task with the given description and total.
// Returns a TaskID that can be used to update the task's progress.
//
// Example:
//
//	task := prog.AddBar("Download", 1000)
//	prog.Update(task, 500) // Set to 50%
func (p *Progress) AddBar(description string, total int64) TaskID {
	bar := NewBar(total).Description(description)
	return p.Add(bar)
}

// AddSpinner adds a spinner task with the given description.
// Returns a TaskID that can be used to control the task.
//
// Example:
//
//	task := prog.AddSpinner("Processing...")
//	prog.Complete(task) // Mark as done
func (p *Progress) AddSpinner(description string) TaskID {
	spinner := NewSpinner(SpinnerDots).Description(description)
	return p.AddSpinnerWithStyle(spinner)
}

// Add adds a progress bar to the manager.
// Returns a TaskID that can be used to update the task.
//
// Example:
//
//	bar := progress.NewBar(1000).CompleteStyle(rich.NewStyle().Foreground(rich.Green))
//	task := prog.Add(bar)
func (p *Progress) Add(bar *ProgressBar) TaskID {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.taskSeq++
	id := p.taskSeq

	p.tasks[id] = &Task{
		id:        id,
		bar:       bar,
		spinner:   nil,
		startTime: time.Now(),
		completed: false,
	}

	return id
}

// AddSpinnerWithStyle adds a spinner to the manager.
// Returns a TaskID that can be used to control the task.
//
// Example:
//
//	spinner := progress.NewSpinner(progress.SpinnerArc).Style(rich.NewStyle().Foreground(rich.Cyan))
//	task := prog.AddSpinnerWithStyle(spinner)
func (p *Progress) AddSpinnerWithStyle(spinner *Spinner) TaskID {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.taskSeq++
	id := p.taskSeq

	p.tasks[id] = &Task{
		id:        id,
		bar:       nil,
		spinner:   spinner,
		startTime: time.Now(),
		completed: false,
	}

	return id
}

// Update sets the progress for a bar task to the given value.
// This is a no-op for spinner tasks.
//
// Thread-safe.
//
// Example:
//
//	prog.Update(task, 750) // Set to 75% of 1000
func (p *Progress) Update(id TaskID, value int64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	task, ok := p.tasks[id]
	if !ok || task.bar == nil {
		return
	}

	task.bar.SetProgress(value)
}

// Advance increments the progress for a bar task by the given delta.
// This is a no-op for spinner tasks.
//
// Thread-safe.
//
// Example:
//
//	prog.Advance(task, 10) // Add 10 to current progress
func (p *Progress) Advance(id TaskID, delta int64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	task, ok := p.tasks[id]
	if !ok || task.bar == nil {
		return
	}

	task.bar.Advance(delta)
}

// Complete marks a task as completed.
// Completed tasks remain visible until Stop() is called.
//
// Thread-safe.
//
// Example:
//
//	prog.Complete(task)
func (p *Progress) Complete(id TaskID) {
	p.mu.Lock()
	defer p.mu.Unlock()

	task, ok := p.tasks[id]
	if !ok {
		return
	}

	task.completed = true
}

// Remove removes a task from the display.
// The task disappears immediately on the next refresh.
//
// Thread-safe.
//
// Example:
//
//	prog.Remove(task)
func (p *Progress) Remove(id TaskID) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.tasks, id)
}

// Start begins the live update loop.
// This spawns a goroutine that periodically refreshes the display.
// Call Stop() to stop the loop and clean up.
//
// It's safe to call Start() multiple times - subsequent calls are ignored.
//
// Example:
//
//	prog.Start()
//	// ... update progress ...
//	prog.Stop()
func (p *Progress) Start() {
	p.mu.Lock()
	if p.running {
		p.mu.Unlock()
		return
	}
	p.running = true
	p.mu.Unlock()

	// Hide cursor for cleaner display
	fmt.Fprint(p.writer, hideCursor)

	// Create ticker for periodic refresh
	p.ticker = time.NewTicker(p.refreshRate)

	// Start render loop in goroutine
	go p.renderLoop()
}

// Stop stops the live update loop and performs final cleanup.
// If transient mode is enabled, clears the progress display.
// Otherwise, leaves the final state visible.
//
// This method blocks until the render loop exits.
//
// Example:
//
//	prog.Stop()
func (p *Progress) Stop() {
	p.mu.Lock()
	if !p.running {
		p.mu.Unlock()
		return
	}
	p.running = false
	p.mu.Unlock()

	// Stop the ticker
	if p.ticker != nil {
		p.ticker.Stop()
	}

	// Signal the render loop to stop
	close(p.stopChan)

	// Final render (or clear if transient)
	if p.transient {
		p.clear()
	} else {
		p.render()
		fmt.Fprintln(p.writer) // Move to next line
	}

	// Show cursor again
	fmt.Fprint(p.writer, showCursor)
}

// renderLoop is the main render loop that runs in a goroutine.
func (p *Progress) renderLoop() {
	for {
		select {
		case <-p.ticker.C:
			p.advanceSpinners()
			p.render()
		case <-p.stopChan:
			return
		}
	}
}

// advanceSpinners advances all spinner animations to the next frame.
func (p *Progress) advanceSpinners() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, task := range p.tasks {
		if task.spinner != nil {
			task.spinner.Next()
		}
	}
}

// render renders all tasks to the console.
func (p *Progress) render() {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if len(p.tasks) == 0 {
		return
	}

	// Move cursor up to start of progress area (if we rendered before)
	if p.lastLineCount > 0 {
		fmt.Fprintf(p.writer, cursorUp, p.lastLineCount)
	}

	// Render each task
	lineCount := 0
	consoleWidth := p.console.Width()

	for _, task := range p.tasks {
		// Move to line start and clear
		fmt.Fprintf(p.writer, "%s%s", cursorLeft, clearLine)

		// Render the bar or spinner
		var segments rich.Segments
		if task.bar != nil {
			segments = task.bar.Render(p.console, consoleWidth)
		} else if task.spinner != nil {
			segments = task.spinner.Render(p.console, consoleWidth)
		}

		// Convert to ANSI and write
		ansi := segments.ToANSI(p.console.ColorMode())
		fmt.Fprint(p.writer, ansi)
		fmt.Fprintln(p.writer)

		lineCount++
	}

	p.lastLineCount = lineCount
}

// clear clears the progress display (for transient mode).
func (p *Progress) clear() {
	if p.lastLineCount == 0 {
		return
	}

	// Move cursor up to start of progress area
	fmt.Fprintf(p.writer, cursorUp, p.lastLineCount)

	// Clear each line
	for i := 0; i < p.lastLineCount; i++ {
		fmt.Fprintf(p.writer, "%s%s\n", cursorLeft, clearLine)
	}

	// Move cursor back up
	fmt.Fprintf(p.writer, cursorUp, p.lastLineCount)

	p.lastLineCount = 0
}
