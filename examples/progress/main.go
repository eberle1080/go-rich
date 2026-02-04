package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"time"

	"github.com/eberle1080/go-rich"
	"github.com/eberle1080/go-rich/progress"
)

func main() {
	// Parse command line flags
	example := flag.String("example", "all", "Which example to run: all, static, live, multi, spinner, file, custom")

	flag.Parse()

	console := rich.NewConsole(nil) // Use stdout

	switch *example {
	case "static":
		staticExample(console)
	case "live":
		liveExample(console)
	case "multi":
		multiExample(console)
	case "spinner":
		spinnerExample(console)
	case "file":
		fileExample(console)
	case "custom":
		customExample(console)
	case "all":
		runAll(console)
	default:
		fmt.Println("Unknown example:", *example)
		fmt.Println("Available examples: all, static, live, multi, spinner, file, custom")
	}
}

func runAll(console *rich.Console) {
	examples := []struct {
		name string
		fn   func(*rich.Console)
	}{
		{"Static Rendering", staticExample},
		{"Live Updates", liveExample},
		{"Multiple Bars", multiExample},
		{"Spinners", spinnerExample},
		{"File Transfer", fileExample},
		{"Custom Styling", customExample},
	}

	for i, ex := range examples {
		console.Println(fmt.Sprintf("\n=== Example %d: %s ===", i+1, ex.name))
		ex.fn(console)
		time.Sleep(500 * time.Millisecond)
	}
}

// staticExample demonstrates static progress bar rendering
func staticExample(console *rich.Console) {
	console.Println("Static rendering - single snapshot:")

	bar := progress.NewBar(100).
		Description("Download").
		Width(40).
		CompleteStyle(rich.NewStyle().Foreground(rich.Green)).
		RemainingStyle(rich.NewStyle().Dim())

	// Render at different completion levels
	for i := 0; i <= 100; i += 25 {
		bar.SetProgress(int64(i))
		console.Render(bar)
		console.Println()
	}
}

// liveExample demonstrates live progress updates
func liveExample(console *rich.Console) {
	console.Println("Live updates - single bar:")

	prog := progress.New(console).
		RefreshRate(50 * time.Millisecond)

	task := prog.AddBar("Processing", 100)
	prog.Start()

	// Simulate work
	for i := 0; i <= 100; i++ {
		prog.Update(task, int64(i))
		time.Sleep(30 * time.Millisecond)
	}

	prog.Complete(task)
	prog.Stop()
}

// multiExample demonstrates multiple concurrent progress bars
func multiExample(console *rich.Console) {
	console.Println("Multiple concurrent bars:")

	prog := progress.New(console).
		RefreshRate(50 * time.Millisecond)

	// Add multiple tasks
	download := prog.AddBar("Download", 1000)
	process := prog.AddBar("Process", 500)
	upload := prog.AddBar("Upload", 750)

	prog.Start()

	// Simulate concurrent work
	done := make(chan bool, 3)

	// Download task
	go func() {
		for i := 0; i <= 1000; i += 10 {
			prog.Update(download, int64(i))
			time.Sleep(20 * time.Millisecond)
		}

		prog.Complete(download)

		done <- true
	}()

	// Process task (slower)
	go func() {
		for i := 0; i <= 500; i += 5 {
			prog.Update(process, int64(i))
			time.Sleep(40 * time.Millisecond)
		}

		prog.Complete(process)

		done <- true
	}()

	// Upload task (medium speed)
	go func() {
		for i := 0; i <= 750; i += 10 {
			prog.Update(upload, int64(i))
			time.Sleep(25 * time.Millisecond)
		}

		prog.Complete(upload)

		done <- true
	}()

	// Wait for all tasks
	for i := 0; i < 3; i++ {
		<-done
	}

	prog.Stop()
}

// spinnerExample demonstrates spinner animations
func spinnerExample(console *rich.Console) {
	console.Println("Spinners for indeterminate progress:")

	prog := progress.New(console).
		RefreshRate(80 * time.Millisecond).
		Transient(true)

	// Add different spinner styles
	spinner1 := progress.NewSpinner(progress.SpinnerDots).
		Description("Loading data...").
		Style(rich.NewStyle().Foreground(rich.Cyan))
	task1 := prog.AddSpinnerWithStyle(spinner1)

	spinner2 := progress.NewSpinner(progress.SpinnerArc).
		Description("Processing...").
		Style(rich.NewStyle().Foreground(rich.Yellow))
	task2 := prog.AddSpinnerWithStyle(spinner2)

	spinner3 := progress.NewSpinner(progress.SpinnerCircle).
		Description("Waiting...").
		Style(rich.NewStyle().Foreground(rich.Magenta))
	task3 := prog.AddSpinnerWithStyle(spinner3)

	prog.Start()

	// Complete tasks one by one
	time.Sleep(2 * time.Second)
	prog.Complete(task1)
	prog.Remove(task1)

	time.Sleep(2 * time.Second)
	prog.Complete(task2)
	prog.Remove(task2)

	time.Sleep(2 * time.Second)
	prog.Complete(task3)
	prog.Remove(task3)

	prog.Stop()
}

// fileExample demonstrates io.Reader wrapper for automatic progress
func fileExample(console *rich.Console) {
	console.Println("File transfer with automatic progress:")

	// Create a fake file with 10MB of data
	fileSize := int64(10 * 1024 * 1024)
	fakeFile := bytes.NewReader(make([]byte, fileSize))

	prog := progress.New(console).
		RefreshRate(50 * time.Millisecond)

	task := prog.AddBar("Upload", fileSize)
	prog.Start()

	// Wrap reader for automatic progress
	reader := progress.NewReader(fakeFile, func(n int) {
		prog.Advance(task, int64(n))
	})

	// Simulate upload by reading in chunks
	buf := make([]byte, 64*1024) // 64KB chunks
	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		// Simulate network delay
		time.Sleep(10 * time.Millisecond)

		_ = n
	}

	prog.Complete(task)
	prog.Stop()
}

// customExample demonstrates custom styling and configuration
func customExample(console *rich.Console) {
	console.Println("Custom styling:")

	// Create custom styled bars
	bar1 := progress.NewBar(100).
		Description("Success").
		Width(30).
		CompleteChar("━").
		RemainingChar("─").
		CompleteStyle(rich.NewStyle().Foreground(rich.Green).Bold()).
		RemainingStyle(rich.NewStyle().Dim())

	bar2 := progress.NewBar(100).
		Description("Warning").
		Width(30).
		CompleteChar("■").
		RemainingChar("□").
		CompleteStyle(rich.NewStyle().Foreground(rich.Yellow)).
		RemainingStyle(rich.NewStyle().Dim())

	bar3 := progress.NewBar(100).
		Description("Error").
		Width(30).
		CompleteChar("▰").
		RemainingChar("▱").
		CompleteStyle(rich.NewStyle().Foreground(rich.Red).Bold()).
		RemainingStyle(rich.NewStyle().Dim())

	prog := progress.New(console).
		RefreshRate(50 * time.Millisecond)

	task1 := prog.Add(bar1)
	task2 := prog.Add(bar2)
	task3 := prog.Add(bar3)

	prog.Start()

	// Update all bars together
	for i := 0; i <= 100; i++ {
		prog.Update(task1, int64(i))
		prog.Update(task2, int64(i))
		prog.Update(task3, int64(i))
		time.Sleep(30 * time.Millisecond)
	}

	prog.Stop()
}
