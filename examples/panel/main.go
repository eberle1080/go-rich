package main

import (
	"os"

	"github.com/eberle1080/go-rich"
	"github.com/eberle1080/go-rich/panel"
	"github.com/eberle1080/go-rich/table"
)

func main() {
	console := rich.NewConsole(os.Stdout)

	console.Println()
	console.Rule("Panel Examples")
	console.Println()

	// Example 1: Simple panel
	console.PrintMarkupln("[bold]1. Simple Panel[/]")
	console.Println()

	p1 := panel.New("Hello, World!")
	console.Renderln(p1)

	// Example 2: Panel with title
	console.Println()
	console.PrintMarkupln("[bold]2. Panel with Title[/]")
	console.Println()

	p2 := panel.New("This is a panel with a title.").
		Title("Welcome")
	console.Renderln(p2)

	// Example 3: Panel with title and subtitle
	console.Println()
	console.PrintMarkupln("[bold]3. Panel with Title and Subtitle[/]")
	console.Println()

	p3 := panel.New("This panel has both a title and a subtitle.").
		Title("Important Message").
		Subtitle("Please read carefully")
	console.Renderln(p3)

	// Example 4: Different box styles
	console.Println()
	console.PrintMarkupln("[bold]4. Box Styles[/]")
	console.Println()

	console.PrintMarkupln("[dim]Simple:[/]")
	p4a := panel.New("Simple box style").Box(table.BoxSimple).Title("Simple")
	console.Renderln(p4a)

	console.Println()
	console.PrintMarkupln("[dim]Rounded:[/]")
	p4b := panel.New("Rounded corners").Box(table.BoxRounded).Title("Rounded")
	console.Renderln(p4b)

	console.Println()
	console.PrintMarkupln("[dim]Double:[/]")
	p4c := panel.New("Double-line borders").Box(table.BoxDouble).Title("Double")
	console.Renderln(p4c)

	console.Println()
	console.PrintMarkupln("[dim]Heavy:[/]")
	p4d := panel.New("Heavy borders").Box(table.BoxHeavy).Title("Heavy")
	console.Renderln(p4d)

	console.Println()
	console.PrintMarkupln("[dim]ASCII:[/]")
	p4e := panel.New("ASCII-only characters").Box(table.BoxASCII).Title("ASCII")
	console.Renderln(p4e)

	// Example 5: Content alignment
	console.Println()
	console.PrintMarkupln("[bold]5. Content Alignment[/]")
	console.Println()

	console.PrintMarkupln("[dim]Left aligned:[/]")
	p5a := panel.New("Left").Align(panel.AlignLeft).Width(40).Title("Left Align")
	console.Renderln(p5a)

	console.Println()
	console.PrintMarkupln("[dim]Center aligned:[/]")
	p5b := panel.New("Center").Align(panel.AlignCenter).Width(40).Title("Center Align")
	console.Renderln(p5b)

	console.Println()
	console.PrintMarkupln("[dim]Right aligned:[/]")
	p5c := panel.New("Right").Align(panel.AlignRight).Width(40).Title("Right Align")
	console.Renderln(p5c)

	// Example 6: Custom styles
	console.Println()
	console.PrintMarkupln("[bold]6. Custom Styles[/]")
	console.Println()

	p6 := panel.New("This panel has custom colors!").
		Title("Styled Panel").
		BorderStyle(rich.NewStyle().Foreground(rich.BrightBlue)).
		TitleStyle(rich.NewStyle().Bold().Foreground(rich.BrightYellow))
	console.Renderln(p6)

	// Example 7: Multi-line content
	console.Println()
	console.PrintMarkupln("[bold]7. Multi-line Content[/]")
	console.Println()

	multiLine := rich.Lines{
		rich.NewRenderableString("Line 1", rich.NewStyle()),
		rich.NewRenderableString("Line 2", rich.NewStyle()),
		rich.NewRenderableString("Line 3", rich.NewStyle()),
	}

	p7 := panel.New(multiLine).Title("Multiple Lines")
	console.Renderln(p7)

	// Example 8: Different padding
	console.Println()
	console.PrintMarkupln("[bold]8. Padding Options[/]")
	console.Println()

	console.PrintMarkupln("[dim]No padding:[/]")
	p8a := panel.New("No padding").Padding(0).Title("Padding: 0")
	console.Renderln(p8a)

	console.Println()
	console.PrintMarkupln("[dim]Default padding (1):[/]")
	p8b := panel.New("Default padding").Padding(1).Title("Padding: 1")
	console.Renderln(p8b)

	console.Println()
	console.PrintMarkupln("[dim]Large padding:[/]")
	p8c := panel.New("Large padding").Padding(3).Title("Padding: 3")
	console.Renderln(p8c)

	// Example 9: Fixed width
	console.Println()
	console.PrintMarkupln("[bold]9. Fixed Width[/]")
	console.Println()

	p9 := panel.New("This panel has a fixed width of 50 characters.").
		Width(50).
		Title("Fixed Width: 50")
	console.Renderln(p9)

	// Example 10: Info panel
	console.Println()
	console.PrintMarkupln("[bold]10. Info Panel[/]")
	console.Println()

	p10 := panel.New("ℹ️  This is an informational message.").
		Title("Information").
		BorderStyle(rich.NewStyle().Foreground(rich.BrightCyan)).
		TitleStyle(rich.NewStyle().Bold().Foreground(rich.BrightCyan))
	console.Renderln(p10)

	// Example 11: Warning panel
	console.Println()
	console.PrintMarkupln("[bold]11. Warning Panel[/]")
	console.Println()

	p11 := panel.New("⚠️  This is a warning message.").
		Title("Warning").
		BorderStyle(rich.NewStyle().Foreground(rich.BrightYellow)).
		TitleStyle(rich.NewStyle().Bold().Foreground(rich.BrightYellow))
	console.Renderln(p11)

	// Example 12: Error panel
	console.Println()
	console.PrintMarkupln("[bold]12. Error Panel[/]")
	console.Println()

	p12 := panel.New("❌ This is an error message.").
		Title("Error").
		BorderStyle(rich.NewStyle().Foreground(rich.BrightRed)).
		TitleStyle(rich.NewStyle().Bold().Foreground(rich.BrightRed))
	console.Renderln(p12)

	// Example 13: Success panel
	console.Println()
	console.PrintMarkupln("[bold]13. Success Panel[/]")
	console.Println()

	p13 := panel.New("✅ Operation completed successfully!").
		Title("Success").
		BorderStyle(rich.NewStyle().Foreground(rich.BrightGreen)).
		TitleStyle(rich.NewStyle().Bold().Foreground(rich.BrightGreen))
	console.Renderln(p13)

	// Example 14: Long content
	console.Println()
	console.PrintMarkupln("[bold]14. Long Content[/]")
	console.Println()

	longText := rich.Lines{
		rich.NewRenderableString("This is a panel with multiple lines of content.", rich.NewStyle()),
		rich.NewRenderableString("It demonstrates how panels can contain longer text.", rich.NewStyle()),
		rich.NewRenderableString("Each line is properly padded and bordered.", rich.NewStyle()),
		rich.NewRenderableString("", rich.NewStyle()),
		rich.NewRenderableString("You can use this for help text, documentation,", rich.NewStyle()),
		rich.NewRenderableString("or any other multi-line content you need to display.", rich.NewStyle()),
	}

	p14 := panel.New(longText).
		Title("Documentation").
		Subtitle("Version 1.0")
	console.Renderln(p14)

	console.Println()
	console.Rule("End of Examples")
	console.Println()
}
