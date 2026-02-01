package main

import (
	"os"

	"github.com/eberle1080/go-rich"
)

func main() {
	console := rich.NewConsole(os.Stdout)

	// Title
	console.Println()
	console.Rule("go-rich Phase 1 Showcase")
	console.Println()

	// Feature 1: Basic colors
	console.PrintStyledln(rich.NewStyle().Bold().Render("1. ANSI Colors"))
	console.PrintStyledln(rich.NewStyle().Foreground(rich.Red).Render("   • Red text"))
	console.PrintStyledln(rich.NewStyle().Foreground(rich.Green).Render("   • Green text"))
	console.PrintStyledln(rich.NewStyle().Foreground(rich.BrightBlue).Render("   • Bright blue text"))
	console.Println()

	// Feature 2: Text styles
	console.PrintStyledln(rich.NewStyle().Bold().Render("2. Text Styles"))
	console.PrintStyledln(rich.NewStyle().Bold().Render("   • Bold"))
	console.PrintStyledln(rich.NewStyle().Italic().Render("   • Italic"))
	console.PrintStyledln(rich.NewStyle().Underline().Render("   • Underline"))
	console.PrintStyledln(rich.NewStyle().Dim().Render("   • Dim/Faint"))
	console.Println()

	// Feature 3: Combined styles
	console.PrintStyledln(rich.NewStyle().Bold().Render("3. Combined Styles"))
	console.PrintStyledln(
		rich.NewStyle().
			Foreground(rich.BrightYellow).
			Background(rich.Red).
			Bold().
			Render("   ⚠ WARNING: Danger ahead!"),
	)
	console.PrintStyledln(
		rich.NewStyle().
			Foreground(rich.BrightGreen).
			Bold().
			Render("   ✓ SUCCESS: All systems go"),
	)
	console.Println()

	// Feature 4: RGB colors
	console.PrintStyledln(rich.NewStyle().Bold().Render("4. RGB Colors"))
	orange := rich.RGB(255, 165, 0)
	console.PrintStyledln(rich.NewStyle().Foreground(orange).Render("   • Custom RGB (255, 165, 0)"))

	hex, _ := rich.Hex("#FF1493")
	console.PrintStyledln(rich.NewStyle().Foreground(hex).Render("   • Hex color #FF1493 (Deep Pink)"))

	named, _ := rich.Named("purple")
	console.PrintStyledln(rich.NewStyle().Foreground(named).Render("   • Named color 'purple'"))
	console.Println()

	// Feature 5: Segments
	console.PrintStyledln(rich.NewStyle().Bold().Render("5. Styled Segments"))
	console.Print("   ")
	console.PrintStyled(rich.NewStyle().Foreground(rich.Red).Bold().Render("Error:"))
	console.Print(" File ")
	console.PrintStyled(rich.NewStyle().Foreground(rich.Cyan).Italic().Render("config.yaml"))
	console.Println(" not found")
	console.Println()

	// Feature 6: Horizontal rules
	console.PrintStyledln(rich.NewStyle().Bold().Render("6. Horizontal Rules"))
	console.Rule("Section Title")
	console.Rule("")
	console.Println()

	// Feature 7: Console info
	console.PrintStyledln(rich.NewStyle().Bold().Render("7. Console Information"))
	console.Printf("   • Color mode: %s\n", colorModeName(console.ColorMode()))
	console.Printf("   • Terminal size: %dx%d\n", console.Width(), console.Height())
	console.Println()

	// Footer
	console.Rule("Phase 1 Complete!")
	console.Println()
	console.PrintStyledln(
		rich.NewStyle().Dim().Render(
			"Next phases: Markup, Tables, Panels, Progress bars...",
		),
	)
	console.Println()
}

func colorModeName(mode rich.ColorMode) string {
	switch mode {
	case rich.ColorModeNone:
		return "None (plain text)"
	case rich.ColorModeStandard:
		return "Standard (16 colors)"
	case rich.ColorMode256:
		return "256 colors"
	case rich.ColorModeTrueColor:
		return "True color (16M colors)"
	default:
		return "Unknown"
	}
}
