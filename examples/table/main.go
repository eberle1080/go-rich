package main

import (
	"os"

	"github.com/eberle1080/go-rich"
	"github.com/eberle1080/go-rich/table"
)

func main() {
	console := rich.NewConsole(os.Stdout)

	console.Println()
	console.Rule("Table Examples")
	console.Println()

	// Example 1: Simple table
	console.PrintMarkupln("[bold]1. Simple Table[/]")
	console.Println()

	t1 := table.New().
		Headers("Name", "Age", "City").
		Row("Alice", "30", "New York").
		Row("Bob", "25", "Los Angeles").
		Row("Charlie", "35", "Chicago")

	console.Renderln(t1)

	// Example 2: Table with title
	console.Println()
	console.PrintMarkupln("[bold]2. Table with Title[/]")
	console.Println()

	t2 := table.New().
		Title("Employee Directory").
		Headers("Name", "Department", "Email").
		Row("Alice Smith", "Engineering", "alice@example.com").
		Row("Bob Johnson", "Sales", "bob@example.com").
		Row("Carol White", "Marketing", "carol@example.com")

	console.Renderln(t2)

	// Example 3: Different box styles
	console.Println()
	console.PrintMarkupln("[bold]3. Box Styles[/]")
	console.Println()

	console.PrintMarkupln("[dim]ASCII:[/]")
	t3a := table.New().
		Box(table.BoxASCII).
		Headers("ID", "Status").
		Row("1", "Active").
		Row("2", "Pending")
	console.Renderln(t3a)

	console.Println()
	console.PrintMarkupln("[dim]Rounded:[/]")
	t3b := table.New().
		Box(table.BoxRounded).
		Headers("ID", "Status").
		Row("1", "Active").
		Row("2", "Pending")
	console.Renderln(t3b)

	console.Println()
	console.PrintMarkupln("[dim]Double:[/]")
	t3c := table.New().
		Box(table.BoxDouble).
		Headers("ID", "Status").
		Row("1", "Active").
		Row("2", "Pending")
	console.Renderln(t3c)

	console.Println()
	console.PrintMarkupln("[dim]Heavy:[/]")
	t3d := table.New().
		Box(table.BoxHeavy).
		Headers("ID", "Status").
		Row("1", "Active").
		Row("2", "Pending")
	console.Renderln(t3d)

	// Example 4: Column alignment
	console.Println()
	console.PrintMarkupln("[bold]4. Column Alignment[/]")
	console.Println()

	t4 := table.New().
		AddColumn(table.NewColumn("Product").WithAlign(table.AlignLeft)).
		AddColumn(table.NewColumn("Price").WithAlign(table.AlignRight)).
		AddColumn(table.NewColumn("Qty").WithAlign(table.AlignCenter)).
		Row("Apple", "$1.99", "10").
		Row("Banana", "$0.99", "25").
		Row("Orange", "$2.49", "15")

	console.Renderln(t4)

	// Example 5: Custom styles
	console.Println()
	console.PrintMarkupln("[bold]5. Custom Styles[/]")
	console.Println()

	t5 := table.New().
		BorderStyle(rich.NewStyle().Foreground(rich.Blue)).
		TitleStyle(rich.NewStyle().Bold().Foreground(rich.Green)).
		Title("Styled Table").
		AddColumn(table.NewColumn("Status").
			WithHeaderStyle(rich.NewStyle().Bold().Foreground(rich.Yellow)).
			WithCellStyle(rich.NewStyle().Foreground(rich.Green))).
		AddColumn(table.NewColumn("Count").
			WithHeaderStyle(rich.NewStyle().Bold().Foreground(rich.Cyan))).
		Row("Active", "42").
		Row("Pending", "7")

	console.Renderln(t5)

	// Example 6: Fixed column widths
	console.Println()
	console.PrintMarkupln("[bold]6. Fixed Column Widths[/]")
	console.Println()

	t6 := table.New().
		AddColumn(table.NewColumn("ID").WithWidth(5)).
		AddColumn(table.NewColumn("Description").WithWidth(30)).
		AddColumn(table.NewColumn("Status").WithWidth(10)).
		Row("1", "This is a very long description that will be truncated", "Done").
		Row("2", "Short desc", "Pending")

	console.Renderln(t6)

	// Example 7: No header, no edge
	console.Println()
	console.PrintMarkupln("[bold]7. Minimal Table (No Header/Edge)[/]")
	console.Println()

	t7 := table.New().
		ShowHeader(false).
		ShowEdge(false).
		Headers("A", "B", "C").
		Row("1", "2", "3").
		Row("4", "5", "6").
		Row("7", "8", "9")

	console.Renderln(t7)

	// Example 8: Server status table
	console.Println()
	console.PrintMarkupln("[bold]8. Server Status Dashboard[/]")
	console.Println()

	t8 := table.New().
		Box(table.BoxRounded).
		Title("Server Status").
		TitleStyle(rich.NewStyle().Bold().Foreground(rich.BrightCyan)).
		BorderStyle(rich.NewStyle().Foreground(rich.BrightBlack)).
		AddColumn(table.NewColumn("Service").WithAlign(table.AlignLeft)).
		AddColumn(table.NewColumn("Status").WithAlign(table.AlignCenter)).
		AddColumn(table.NewColumn("Uptime").WithAlign(table.AlignRight)).
		Row("API Server", "🟢 Running", "99.9%").
		Row("Database", "🟢 Running", "99.8%").
		Row("Cache", "🟡 Degraded", "95.2%").
		Row("Queue", "🟢 Running", "99.7%")

	console.Renderln(t8)

	console.Println()
	console.Rule("End of Examples")
	console.Println()
}
