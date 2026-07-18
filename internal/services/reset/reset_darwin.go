package reset

import (
	"context"
	"os/exec"
	"path/filepath"

	"github.com/insigmo/jetreset/internal/logx"
)

// Reset wipes JetBrains trial state for the given products on macOS.
func Reset(home string, products []string) {
	logx.Trace("reset: Library prefs + plist")
	CleanDir(filepath.Join(home, "Library", "Preferences"), products)
	CleanDir(filepath.Join(home, "Library", "Application Support", "JetBrains"), products)
	plist := filepath.Join(home, "Library", "Preferences", "com.apple.java.util.prefs.plist")
	for _, k := range []string{
		"/.JetBrains.UserIdOnMachine",
		"/.jetbrains/.user_id_on_machine",
		"/.jetbrains/.device_id",
	} {
		_ = exec.CommandContext(context.Background(), "plutil", "-remove", k, plist).Run() //nolint:gosec // G204: fixed command
	}
}
