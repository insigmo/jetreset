package reset

import (
	"path/filepath"
	"regexp"
)

var (
	reJBUserIdOnMachine = regexp.MustCompile(`<entry key="JetBrains\.UserIdOnMachine"`)
	reDeviceId          = regexp.MustCompile(`<entry key="device_id"`)
	reUserIdOnMachine   = regexp.MustCompile(`<entry key="user_id_on_machine"`)
)

func Reset(home string, products []string) {
	cleanDir(filepath.Join(home, ".config", "JetBrains"), products)
	jprefs := filepath.Join(home, ".java", ".userPrefs")

	removeLine(filepath.Join(jprefs, "prefs.xml"), reJBUserIdOnMachine)
	removeLine(filepath.Join(jprefs, "jetbrains", "prefs.xml"), reDeviceId)
	removeLine(filepath.Join(jprefs, "jetbrains", "prefs.xml"), reUserIdOnMachine)
}
