package init

import (
	remotesigner "github.com/quan-to/chevron"
	"github.com/quan-to/slog"
)

// Load here to guarantee this will load before log is used.
// This package should be imported by the main package so the code gets executed.
var _ = func() error {
	slog.SetLogFormat(remotesigner.LogFormat)
	return nil
}()
