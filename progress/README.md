# Progress Bars and Spinners

The `progress` package provides comprehensive progress bar and spinner functionality for go-rich, supporting both simple single-bar usage and complex multi-bar scenarios with full customization.

## Features

- **Single and multi-bar support** - Display one or many progress indicators
- **Static and live rendering** - Choose between one-time snapshots or continuous updates
- **io.Reader/Writer wrappers** - Automatic progress tracking for file operations
- **Customizable spinners** - 10+ built-in styles, or create your own
- **Speed and ETA tracking** - Automatic calculation with exponential smoothing
- **Column system** - Mix and match display elements (description, bar, percentage, speed, ETA)
- **Thread-safe** - Safe for concurrent updates from multiple goroutines
- **ANSI control** - Smooth in-place updates without scrolling

## Quick Start

### Simple Progress Bar

```go
console := rich.NewConsole(nil)

bar := progress.NewBar(100).
    Description("Download").
    Width(40).
    CompleteStyle(rich.NewStyle().Foreground(rich.Green))

for i := 0; i <= 100; i++ {
    bar.SetProgress(int64(i))
    console.Render(bar)
    console.Println()
    time.Sleep(50 * time.Millisecond)
}
```

### Live Progress Updates

```go
prog := progress.New(console)
task := prog.AddBar("Processing", 1000)

prog.Start()

for i := 0; i <= 1000; i++ {
    prog.Update(task, int64(i))
    time.Sleep(10 * time.Millisecond)
}

prog.Stop()
```

### Multiple Concurrent Bars

```go
prog := progress.New(console)

download := prog.AddBar("Download", 1000)
process := prog.AddBar("Process", 500)
upload := prog.AddBar("Upload", 750)

prog.Start()

// Update from different goroutines
go func() {
    for i := 0; i <= 1000; i += 10 {
        prog.Update(download, int64(i))
        time.Sleep(20 * time.Millisecond)
    }
    prog.Complete(download)
}()

go func() {
    for i := 0; i <= 500; i += 5 {
        prog.Update(process, int64(i))
        time.Sleep(40 * time.Millisecond)
    }
    prog.Complete(process)
}()

// ... handle upload ...

prog.Stop()
```

## Spinners

For indeterminate progress (when you don't know the total):

```go
spinner := progress.NewSpinner(progress.SpinnerDots).
    Description("Loading...").
    Style(rich.NewStyle().Foreground(rich.Cyan))

prog := progress.New(console)
task := prog.AddSpinnerWithStyle(spinner)

prog.Start()
time.Sleep(5 * time.Second)
prog.Stop()
```

### Built-in Spinner Styles

- `SpinnerDots` - Braille dots: ⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏
- `SpinnerLine` - Simple line: -\|/
- `SpinnerArc` - Arc animation: ◜◠◝◞◡◟
- `SpinnerArrow` - Rotating arrow: ←↖↑↗→↘↓↙
- `SpinnerCircle` - Circle quadrants: ◐◓◑◒
- `SpinnerBounce` - Bouncing dots: ⠁⠂⠄⡀⢀⠠⠐⠈
- `SpinnerBoxBounce` - Box animation: ▖▘▝▗
- `SpinnerSimple` - ASCII-safe: |/-\
- `SpinnerGrowVertical` - Growing vertical: ▁▃▄▅▆▇▆▅▄▃
- `SpinnerGrowHorizontal` - Growing horizontal: ▏▎▍▌▋▊▉▊▋▌▍▎

### Custom Spinners

```go
spinner := progress.NewSpinner([]string{".", "..", "..."}).
    Interval(200 * time.Millisecond)
```

## Automatic Progress with io.Reader/Writer

Track file operations automatically:

```go
file, _ := os.Open("large-file.dat")
fileInfo, _ := file.Stat()

prog := progress.New(console)
task := prog.AddBar("Reading", fileInfo.Size())
prog.Start()

// Wrap reader for automatic progress
reader := progress.NewReader(file, func(n int) {
    prog.Advance(task, int64(n))
})

io.Copy(dest, reader)
prog.Stop()
```

Same for writing:

```go
writer := progress.NewWriter(file, func(n int) {
    prog.Advance(task, int64(n))
})

io.Copy(writer, source)
```

## Customization

### Bar Appearance

```go
bar := progress.NewBar(100).
    Description("Custom").
    Width(30).
    CompleteChar("━").
    RemainingChar("─").
    CompleteStyle(rich.NewStyle().Foreground(rich.Green).Bold()).
    RemainingStyle(rich.NewStyle().Dim())
```

### Progress Manager Options

```go
prog := progress.New(console).
    RefreshRate(50 * time.Millisecond).  // Faster updates (20 FPS)
    Transient(true)                       // Clear on completion
```

### Spinner Customization

```go
spinner := progress.NewSpinner(progress.SpinnerArc).
    Description("Processing data...").
    Style(rich.NewStyle().Foreground(rich.Yellow).Bold()).
    Interval(100 * time.Millisecond)
```

## Column System

The column system allows fine-grained control over what's displayed. While not yet exposed in the public API, the foundation is in place for future enhancements like:

- `DescriptionColumn` - Task description
- `BarColumn` - Visual progress bar
- `PercentageColumn` - Completion percentage
- `SpeedColumn` - Items per second
- `ETAColumn` - Estimated time remaining
- `ElapsedColumn` - Time spent
- `TransferSpeedColumn` - Bytes/sec with unit conversion

## API Reference

### ProgressBar

```go
type ProgressBar struct { ... }

func NewBar(total int64) *ProgressBar

// Fluent configuration
func (pb *ProgressBar) Description(desc string) *ProgressBar
func (pb *ProgressBar) Width(width int) *ProgressBar
func (pb *ProgressBar) CompleteChar(char string) *ProgressBar
func (pb *ProgressBar) RemainingChar(char string) *ProgressBar
func (pb *ProgressBar) CompleteStyle(style rich.Style) *ProgressBar
func (pb *ProgressBar) RemainingStyle(style rich.Style) *ProgressBar

// Progress control
func (pb *ProgressBar) SetProgress(current int64)
func (pb *ProgressBar) Advance(delta int64)

// Status
func (pb *ProgressBar) Current() int64
func (pb *ProgressBar) Total() int64
func (pb *ProgressBar) Percentage() float64
func (pb *ProgressBar) IsComplete() bool

// Implements rich.Renderable
func (pb *ProgressBar) Render(console *rich.Console, width int) rich.Segments
```

### Spinner

```go
type Spinner struct { ... }

func NewSpinner(frames []string) *Spinner

// Fluent configuration
func (s *Spinner) Description(desc string) *Spinner
func (s *Spinner) Style(style rich.Style) *Spinner
func (s *Spinner) Interval(interval time.Duration) *Spinner

// Animation
func (s *Spinner) Next()
func (s *Spinner) CurrentFrame() string

// Implements rich.Renderable
func (s *Spinner) Render(console *rich.Console, width int) rich.Segments
```

### Progress Manager

```go
type Progress struct { ... }

func New(console *rich.Console) *Progress

// Configuration
func (p *Progress) RefreshRate(rate time.Duration) *Progress
func (p *Progress) Transient(transient bool) *Progress

// Task management
func (p *Progress) AddBar(description string, total int64) TaskID
func (p *Progress) AddSpinner(description string) TaskID
func (p *Progress) Add(bar *ProgressBar) TaskID
func (p *Progress) AddSpinnerWithStyle(spinner *Spinner) TaskID

// Progress updates
func (p *Progress) Update(id TaskID, value int64)
func (p *Progress) Advance(id TaskID, delta int64)
func (p *Progress) Complete(id TaskID)
func (p *Progress) Remove(id TaskID)

// Lifecycle
func (p *Progress) Start()
func (p *Progress) Stop()
```

### io Wrappers

```go
func NewReader(r io.Reader, callback func(int)) io.Reader
func NewWriter(w io.Writer, callback func(int)) io.Writer
```

## Implementation Details

### Speed and ETA Calculation

The tracker uses exponential moving average (EMA) with a smoothing factor of 0.5 to calculate speed:

1. Maintains a rolling window of the last 2 seconds of samples
2. Calculates instantaneous speed from first and last sample in window
3. Applies smoothing to reduce noise

ETA is calculated as: `(total - current) / speed`

### ANSI Control Sequences

Live updates use ANSI cursor control for smooth in-place rendering:

- `ESC[<N>A` - Move cursor up N lines
- `ESC[1000D` - Move to line start
- `ESC[0K` - Clear line from cursor
- `ESC[?25l` / `ESC[?25h` - Hide/show cursor

### Thread Safety

The Progress manager uses `sync.RWMutex` to protect task access, making it safe for concurrent updates from multiple goroutines.

### Memory Management

The tracker automatically prunes old samples to prevent unbounded growth:
- Keeps last 100 samples
- When exceeded, prunes to 50 most recent
- Prevents memory issues in long-running operations

## Examples

See `examples/progress/main.go` for comprehensive examples:

```bash
go run examples/progress/main.go -example all       # Run all examples
go run examples/progress/main.go -example static    # Static rendering
go run examples/progress/main.go -example live      # Live updates
go run examples/progress/main.go -example multi     # Multiple bars
go run examples/progress/main.go -example spinner   # Spinners
go run examples/progress/main.go -example file      # File transfer
go run examples/progress/main.go -example custom    # Custom styling
```

## Testing

Run the test suite:

```bash
go test ./progress/... -v
```

Tests cover:
- Progress bar rendering and calculations
- Speed and ETA tracking accuracy
- Spinner frame cycling
- io.Reader/Writer wrappers
- Thread safety (via race detector: `go test -race`)

## Design Philosophy

The progress package follows go-rich conventions:

1. **Fluent API** - Chainable methods for easy configuration
2. **Immutable styles** - Styles are value types, modifications return new instances
3. **Renderable interface** - Progress bars integrate seamlessly with Console
4. **Hybrid approach** - Support both static snapshots and live updates
5. **No dependencies** - Only uses Go standard library and existing go-rich code
