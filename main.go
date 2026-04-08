package main

import (
	"embed"
	"io/fs"

	"github.com/williamsantosa/cli-repl-template/cmd"
	"github.com/williamsantosa/cli-repl-template/internal/app"
)

//go:embed all:assets
var embeddedAssets embed.FS

func main() {
	sub, _ := fs.Sub(embeddedAssets, "assets")
	app.EmbeddedAssets = sub
	cmd.Execute()
}
