package main

import (
	"fmt"
	"os"

	"github.com/eberle1080/go-rich"
)

func main() {
	console := rich.NewConsole(os.Stdout)

	console.Rule("ANSI Colors")
	console.Println()

	// Standard ANSI colors
	colors := []struct {
		name  string
		color rich.ANSIColor
	}{
		{"Black", rich.Black},
		{"Red", rich.Red},
		{"Green", rich.Green},
		{"Yellow", rich.Yellow},
		{"Blue", rich.Blue},
		{"Magenta", rich.Magenta},
		{"Cyan", rich.Cyan},
		{"White", rich.White},
		{"Bright Black", rich.BrightBlack},
		{"Bright Red", rich.BrightRed},
		{"Bright Green", rich.BrightGreen},
		{"Bright Yellow", rich.BrightYellow},
		{"Bright Blue", rich.BrightBlue},
		{"Bright Magenta", rich.BrightMagenta},
		{"Bright Cyan", rich.BrightCyan},
		{"Bright White", rich.BrightWhite},
	}

	for _, c := range colors {
		style := rich.NewStyle().Foreground(c.color)
		console.PrintStyled(style.Render(fmt.Sprintf("%-20s", c.name)))

		bgStyle := rich.NewStyle().Background(c.color).Foreground(rich.White)
		console.PrintStyledln(bgStyle.Render(" ████ "))
	}

	console.Println()
	console.Rule("Text Attributes")
	console.Println()

	// Text attributes
	console.PrintStyledln(rich.NewStyle().Bold().Render("Bold text"))
	console.PrintStyledln(rich.NewStyle().Italic().Render("Italic text"))
	console.PrintStyledln(rich.NewStyle().Underline().Render("Underlined text"))
	console.PrintStyledln(rich.NewStyle().Strikethrough().Render("Strikethrough text"))
	console.PrintStyledln(rich.NewStyle().Dim().Render("Dim/faint text"))
	console.PrintStyledln(rich.NewStyle().Reverse().Render("Reverse video"))

	console.Println()
	console.Rule("Combined Styles")
	console.Println()

	// Combined styles
	console.PrintStyledln(
		rich.NewStyle().
			Foreground(rich.BrightYellow).
			Background(rich.Blue).
			Bold().
			Render(" WARNING: Combined styles! "),
	)

	console.PrintStyledln(
		rich.NewStyle().
			Foreground(rich.BrightRed).
			Bold().
			Underline().
			Render("ERROR: Critical issue detected"),
	)

	console.PrintStyledln(
		rich.NewStyle().
			Foreground(rich.BrightGreen).
			Italic().
			Render("✓ Success: Operation completed"),
	)

	console.Println()
	console.Rule("RGB Colors")
	console.Println()

	// RGB colors
	rgb1 := rich.RGB(255, 100, 50)
	console.PrintStyledln(rich.NewStyle().Foreground(rgb1).Render("Custom RGB color (255, 100, 50)"))

	rgb2, _ := rich.Hex("#FF1493")
	console.PrintStyledln(rich.NewStyle().Foreground(rgb2).Render("Deep Pink from hex (#FF1493)"))

	rgb3, _ := rich.Named("orange")
	console.PrintStyledln(rich.NewStyle().Foreground(rgb3).Render("Named color: orange"))

	console.Println()
	console.Rule("Color Gradients")
	console.Println()

	// Simple gradient effect
	for i := 0; i < 16; i++ {
		r := uint8(255 - (i * 16))
		g := uint8(i * 16)
		b := uint8(128)
		color := rich.RGB(r, g, b)
		console.PrintStyled(rich.NewStyle().Foreground(color).Render("█"))
	}
	console.Println()

	console.Println()
	console.Rule("Color Mode Detection")
	console.Println()

	modeNames := map[rich.ColorMode]string{
		rich.ColorModeNone:      "None (no colors)",
		rich.ColorModeStandard:  "Standard (16 colors)",
		rich.ColorMode256:       "256 colors",
		rich.ColorModeTrueColor: "True color (16M colors)",
	}

	console.Printf("Detected mode: %s\n", modeNames[console.ColorMode()])
	console.Printf("Terminal size: %dx%d\n", console.Width(), console.Height())
}
