package reset

import (
	"path/filepath"
	"regexp"

	"github.com/insigmo/jetreset/internal/logx"
)

var (
	reJBUserIDOnMachine = regexp.MustCompile(`<entry key="JetBrains\.UserIdOnMachine"`)
	reDeviceID          = regexp.MustCompile(`<entry key="device_id"`)
	reUserIDOnMachine   = regexp.MustCompile(`<entry key="user_id_on_machine"`)
)

// Reset wipes JetBrains trial state for the given products on Linux.
func Reset(home string, products []string) {
	logx.Debugf("reset: ~/.config/JetBrains + Java prefs")
	CleanDir(filepath.Join(home, ".config", "JetBrains"), products)
	jprefs := filepath.Join(home, ".java", ".userPrefs")
	RemoveLine(filepath.Join(jprefs, "prefs.xml"), reJBUserIDOnMachine)
	RemoveLine(filepath.Join(jprefs, "jetbrains", "prefs.xml"), reDeviceID)
	RemoveLine(filepath.Join(jprefs, "jetbrains", "prefs.xml"), reUserIDOnMachine)
}
