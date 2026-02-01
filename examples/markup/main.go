package main

import (
	"os"

	"github.com/eberle1080/go-rich"
)

func main() {
	console := rich.NewConsole(os.Stdout)

	console.Println()
	console.Rule("Markup Examples")
	console.Println()

	// Basic colors
	console.PrintMarkupln("[bold]1. Basic Colors[/]")
	console.PrintMarkupln("   [red]Red text[/]")
	console.PrintMarkupln("   [green]Green text[/]")
	console.PrintMarkupln("   [blue]Blue text[/]")
	console.PrintMarkupln("   [yellow]Yellow text[/]")
	console.Println()

	// Text attributes
	console.PrintMarkupln("[bold]2. Text Attributes[/]")
	console.PrintMarkupln("   [bold]Bold text[/]")
	console.PrintMarkupln("   [italic]Italic text[/]")
	console.PrintMarkupln("   [underline]Underlined text[/]")
	console.PrintMarkupln("   [dim]Dim text[/]")
	console.Println()

	// Combined styles
	console.PrintMarkupln("[bold]3. Combined Styles[/]")
	console.PrintMarkupln("   [bold red]Bold red text[/]")
	console.PrintMarkupln("   [italic green]Italic green text[/]")
	console.PrintMarkupln("   [underline blue]Underlined blue text[/]")
	console.PrintMarkupln("   [bold italic yellow]Bold italic yellow[/]")
	console.Println()

	// Background colors
	console.PrintMarkupln("[bold]4. Background Colors[/]")
	console.PrintMarkupln("   [red on white]Red on white[/]")
	console.PrintMarkupln("   [white on red]White on red[/]")
	console.PrintMarkupln("   [yellow on blue]Yellow on blue[/]")
	console.PrintMarkupln("   [bold white on magenta]Bold white on magenta[/]")
	console.Println()

	// Hex colors
	console.PrintMarkupln("[bold]5. Hex Colors[/]")
	console.PrintMarkupln("   [#FF1493]Deep pink (#FF1493)[/]")
	console.PrintMarkupln("   [#00CED1]Dark turquoise (#00CED1)[/]")
	console.PrintMarkupln("   [#FFD700]Gold (#FFD700)[/]")
	console.Println()

	// RGB colors
	console.PrintMarkupln("[bold]6. RGB Colors[/]")
	console.PrintMarkupln("   [rgb(255,100,50)]Custom RGB (255,100,50)[/]")
	console.PrintMarkupln("   [rgb(50,200,100)]Custom RGB (50,200,100)[/]")
	console.Println()

	// Nested styles
	console.PrintMarkupln("[bold]7. Nested Styles[/]")
	console.PrintMarkupln("   [bold]Bold [red]and red[/] and bold again[/]")
	console.PrintMarkupln("   [green]Green [bold]and bold[/] green again[/]")
	console.Println()

	// Practical examples
	console.PrintMarkupln("[bold]8. Practical Examples[/]")
	console.PrintMarkupln("   [bold red]Error:[/] File not found")
	console.PrintMarkupln("   [bold yellow]Warning:[/] Low disk space")
	console.PrintMarkupln("   [bold green]Success:[/] Operation completed")
	console.PrintMarkupln("   [bold blue]Info:[/] Processing data...")
	console.Println()

	// Log levels
	console.PrintMarkupln("[bold]9. Log Levels[/]")
	console.PrintMarkupln("   [dim]DEBUG:[/] Application started")
	console.PrintMarkupln("   [cyan]INFO:[/] Connected to database")
	console.PrintMarkupln("   [yellow]WARN:[/] Deprecated API used")
	console.PrintMarkupln("   [red]ERROR:[/] Failed to save file")
	console.PrintMarkupln("   [bold red on white]FATAL:[/] System crash")
	console.Println()

	// Escaped brackets
	console.PrintMarkupln("[bold]10. Escaped Brackets[/]")
	console.PrintMarkupln("   Use [[[ for literal brackets")
	console.PrintMarkupln("   [[bold] is not a tag")
	console.PrintMarkupln("   [[red]text[[/] shows as [[red]text[[/]")
	console.Println()

	// Status indicators
	console.PrintMarkupln("[bold]11. Status Indicators[/]")
	console.PrintMarkupln("   [green]✓[/] Test passed")
	console.PrintMarkupln("   [red]✗[/] Test failed")
	console.PrintMarkupln("   [yellow]⚠[/] Test skipped")
	console.PrintMarkupln("   [blue]ℹ[/] Information")
	console.Println()

	// Progress messages
	console.PrintMarkupln("[bold]12. Progress Messages[/]")
	console.PrintMarkup("[bold]Downloading:[/] ")
	console.PrintMarkupln("[green]████████████[/][dim]░░░░░░░░[/] 60%")
	console.PrintMarkup("[bold]Processing:[/] ")
	console.PrintMarkupln("[cyan]■■■■■■■■■■[/][dim]□□□□□□□□□□[/] 50%")
	console.Println()

	// Mixed content
	console.PrintMarkupln("[bold]13. Mixed Inline Styles[/]")
	console.PrintMarkup("The ")
	console.PrintMarkup("[bold]quick[/] ")
	console.PrintMarkup("[italic red]brown[/] ")
	console.PrintMarkup("fox ")
	console.PrintMarkup("[underline blue]jumps[/] ")
	console.PrintMarkupln("over the lazy dog")
	console.Println()

	console.Rule("End of Examples")
	console.Println()

	// Show helper functions
	console.PrintMarkupln("[dim]Helper functions:[/]")
	original := "[bold red]formatted[/] text"
	console.Printf("   Original: %s\n", original)
	console.Printf("   Stripped: %s\n", rich.StripMarkup(original))
	console.Printf("   Escaped:  %s\n", rich.EscapeMarkup("[bold]"))
	console.Println()
}
