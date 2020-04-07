package main

import (
	"flag"
	"github.com/asticode/go-astilectron"
	"github.com/quan-to/chevron/internal/tools"
	"github.com/quan-to/slog"
)

// Vars
var (
	AppName string
	BuiltAt string
	debug   = flag.Bool("d", false, "enables the debug mode")
	w       *astilectron.Window
	log     = slog.Scope("AgentUI").Tag(tools.DefaultTag)
)

func main() {
	// Init
	flag.Parse()
	slog.SetDebug(*debug)

	Migrate()

	// Run bootstrap
	log.Debug("Running app built at %s", BuiltAt)
	Begin()
	Run()
}
