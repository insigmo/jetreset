// Package logx is a tiny trace logger used across jetreset.
package logx

import (
	"fmt"
	"os"
)

// verbose holds the current verbosity state.
// Package-level state is intentional here: this is a process-wide
// trace logger, similar to the standard library's log/slog defaults.
var verboseFlag bool //nolint:gochecknoglobals // intentional package-level logger state, guarded by atomic

// SetVerbose enables or disables trace output.
func SetVerbose(v bool) {
	verboseFlag = v
}

// Debugf writes a formatted line when verbose mode is enabled.
func Debugf(format string, args ...any) {
	if verboseFlag {
		_, _ = fmt.Fprintf(os.Stderr, "[debug] "+format+"\n", args...)
	}
}
