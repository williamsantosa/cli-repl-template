# CLI REPL Template

A configurable Go command-line interface template featuring a loading screen with ANSI block art.

## Quick Start

```bash
# Build (use -o to name the output binary whatever you like)
go build -o myapp.exe .

# Show help
./myapp --help

# Display the art
./myapp show

# Run the loading animation demo (3 seconds by default)
./myapp loading

# Run for a custom duration
./myapp loading -d 5

# Print version
./myapp version
```

## Configuration

The CLI reads its config from `config.yaml`, searched in this order:

1. Path given via `--config` flag
2. Current working directory
3. `$HOME/.cli-repl/config.yaml`

All settings have defaults, so no config file is required.

### Config Reference

```yaml
art:
  source: "built-in"       # "built-in", path to image (.png/.jpg/.gif/.bmp), or path to .txt
  width: 40                 # character width for image rendering (ignored for .txt/built-in)
  border: true              # wrap art in a rounded border
  border_color: "63"        # ANSI 256 color for the border

loader:
  spinner: "dots"           # dots, line, minidot, jump, pulse, points, globe, moon, monkey, meter, hamburger
  spinner_color: "205"      # ANSI 256 color for the spinner
  speed_ms: 100             # milliseconds between spinner frames
  message_color: "252"      # ANSI 256 color for the status text
```

### Custom Art from an Image

Point `art.source` at any image file and it will be rendered in true color using Unicode half-block characters (`▄`), similar to chafa. Each character cell displays two pixels vertically for high-quality output:

```yaml
art:
  source: "./my-image.png"
  width: 50
```

Supported formats: PNG, JPEG, GIF, BMP.

### Custom Art from a Text File

For pre-made ANSI art, point to a `.txt` file instead:

```yaml
art:
  source: "./my-art.txt"
```

## Adding Commands

Create a new file in `cmd/`:

```go
package cmd

import "github.com/spf13/cobra"

var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Description here",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Use app.RunLoader("Working...", func() error { ... })
        // to wrap long-running tasks with the loader.
        return nil
    },
}

func init() {
    rootCmd.AddCommand(myCmd)
}
```

## Using the Loader in Code

```go
import "github.com/williamsantosa/cli-repl-template/internal/app"

err := app.RunLoader("Downloading...", func() error {
    // your long-running work here
    return nil
})
```

## Cross-Platform

Builds and runs on Windows, macOS, and Linux with a single `go build`. ANSI colors work in Windows Terminal, iTerm2, and most modern terminal emulators.

## Build with Version Info

```bash
go build -ldflags "-X github.com/williamsantosa/cli-repl-template/internal/app.Version=1.0.0 -X github.com/williamsantosa/cli-repl-template/cmd.BuildDate=$(date -u +%Y-%m-%d)" -o myapp.exe .
```
