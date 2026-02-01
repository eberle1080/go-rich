# go-rich

A Go port of Python's [rich library](https://github.com/Textualize/rich) for beautiful terminal output.

## Features

- 🎨 **Rich Colors**: ANSI, 256-color, and true color (16M) RGB support
- 💅 **Text Styling**: Bold, italic, underline, strikethrough, dim, and reverse
- 🎯 **Automatic Detection**: Detects terminal capabilities automatically
- 🔧 **Composable**: Fluent API for building complex styles
- 📦 **Zero Dependencies**: Only uses Go standard library (except `golang.org/x/term` for terminal detection)

## Installation

```bash
go get github.com/eberle1080/go-rich
```

## Quick Start

```go
package main

import (
    "os"
    "github.com/eberle1080/go-rich"
)

func main() {
    console := rich.NewConsole(os.Stdout)

    // Simple colored text
    style := rich.NewStyle().Foreground(rich.Green).Bold()
    console.PrintStyledln(style.Render("Success!"))

    // RGB colors
    rgb := rich.RGB(255, 100, 50)
    console.PrintStyledln(rich.NewStyle().Foreground(rgb).Render("Custom color"))

    // Hex colors
    hex, _ := rich.Hex("#FF1493")
    console.PrintStyledln(rich.NewStyle().Foreground(hex).Render("Deep pink"))
}
```

## Examples

Run the example programs:

```bash
# Basic hello world
cd examples/hello
go run main.go

# Style showcase
cd examples/styles
go run main.go
```

## Usage

### Creating a Console

```go
// Write to stdout
console := rich.NewConsole(os.Stdout)

// Write to any io.Writer
var buf bytes.Buffer
console := rich.NewConsole(&buf)
```

### Colors

**ANSI Colors (16 colors):**
```go
rich.Black, rich.Red, rich.Green, rich.Yellow
rich.Blue, rich.Magenta, rich.Cyan, rich.White
rich.BrightRed, rich.BrightGreen, ... // Bright variants
```

**RGB Colors:**
```go
rgb := rich.RGB(255, 100, 50)
hex, _ := rich.Hex("#FF1493")
named, _ := rich.Named("orange")
```

### Styles

Styles are immutable and composable:

```go
style := rich.NewStyle().
    Foreground(rich.Red).
    Background(rich.White).
    Bold().
    Underline()

console.PrintStyledln(style.Render("Error!"))
```

Available style methods:
- `Foreground(color)` - Set text color
- `Background(color)` - Set background color
- `Bold()` - Bold text
- `Italic()` - Italic text
- `Underline()` - Underline text
- `Strikethrough()` - Strikethrough text
- `Dim()` - Dim/faint text
- `Reverse()` - Reverse video

### Color Modes

go-rich automatically detects terminal capabilities:

- `ColorModeNone` - No colors (plain text)
- `ColorModeStandard` - 16 ANSI colors
- `ColorMode256` - 256 colors
- `ColorModeTrueColor` - True color (16M colors)

Override detection:
```go
console.SetColorMode(rich.ColorModeTrueColor)
```

Respects the `NO_COLOR` environment variable (https://no-color.org/).

## Roadmap

This is Phase 1 (Foundation) of the implementation. Future phases will add:

- **Phase 2**: Markup support (`[bold red]text[/]`)
- **Phase 3**: Tables with borders and styling
- **Phase 4**: Panels (bordered containers)
- **Phase 5**: Progress bars and live updates

## Design Philosophy

- **Go-idiomatic**: Uses interfaces, `io.Writer`, and value types
- **Immutable styles**: Thread-safe and composable
- **Progressive enhancement**: Works everywhere, looks best where supported
- **No global state**: All instance-based for testability

## License

MIT License

## Credits

Inspired by [rich](https://github.com/Textualize/rich) by Will McGugan.
