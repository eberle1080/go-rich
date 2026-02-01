package main

import (
	"os"

	"github.com/eberle1080/go-rich"
)

func main() {
	console := rich.NewConsole(os.Stdout)

	// Plain text
	console.Println("Hello, world!")

	// Colored text
	style := rich.NewStyle().Foreground(rich.Green).Bold()
	console.PrintStyledln(style.Render("Hello, colorful world!"))

	// Red error message
	errorStyle := rich.NewStyle().Foreground(rich.Red).Bold()
	console.PrintStyledln(errorStyle.Render("Error: something went wrong!"))

	// Blue info with underline
	infoStyle := rich.NewStyle().Foreground(rich.Blue).Underline()
	console.PrintStyledln(infoStyle.Render("Info: this is underlined"))

	// Horizontal rule
	console.Println()
	console.Rule("Welcome")
	console.Println()

	// Multiple styles
	console.Print("This is ")
	console.PrintStyled(rich.NewStyle().Bold().Render("bold"))
	console.Print(" and this is ")
	console.PrintStyled(rich.NewStyle().Italic().Foreground(rich.Magenta).Render("italic magenta"))
	console.Println(".")
}
