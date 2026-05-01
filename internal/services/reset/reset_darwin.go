package reset

import (
	"os/exec"
	"path/filepath"
)

func Reset(home string, products []string) {
	cleanDir(filepath.Join(home, "Library", "Preferences"), products)
	cleanDir(filepath.Join(home, "Library", "Application Support", "JetBrains"), products)

	plist := filepath.Join(home, "Library", "Preferences", "com.apple.java.util.prefs.plist")
	filesToRemove := []string{
		"/.JetBrains.UserIdOnMachine",
		"/.jetbrains/.user_id_on_machine",
		"/.jetbrains/.device_id",
	}

	for _, k := range filesToRemove {
		exec.Command("plutil", "-remove", k, plist).Run()
	}
}
