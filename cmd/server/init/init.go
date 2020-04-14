package init

import (
	"github.com/quan-to/chevron/internal/config"
	"github.com/quan-to/slog"
)

// Load here to guarantee this will load before log is used.
// This package should be imported by the main package so the code gets executed.
var _ = func() error {
	slog.SetLogFormat(config.LogFormat)
	return nil
}()
