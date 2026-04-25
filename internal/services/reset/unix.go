package reset

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
)

func resetLinux(h string, products []string) {
	fmt.Println("Linux:")
	cleanDir(filepath.Join(h, ".config", "JetBrains"), products)
	jprefs := filepath.Join(h, ".java", ".userPrefs")

	removeLine(
		filepath.Join(jprefs, "prefs.xml"),
		regexp.MustCompile(`<entry key="JetBrains\.UserIdOnMachine"`),
	)
	removeLine(
		filepath.Join(jprefs, "jetbrains", "prefs.xml"),
		regexp.MustCompile(`<entry key="device_id"`),
	)
	removeLine(
		filepath.Join(jprefs, "jetbrains", "prefs.xml"),
		regexp.MustCompile(`<entry key="user_id_on_machine"`),
	)
}

func resetMacos(h string, products []string) {
	fmt.Println("macOS:")
	cleanDir(filepath.Join(h, "Library", "Preferences"), products)
	cleanDir(filepath.Join(h, "Library", "Application Support", "JetBrains"), products)

	plist := filepath.Join(h, "Library", "Preferences", "com.apple.java.util.prefs.plist")
	filesToRemove := []string{
		"/.JetBrains.UserIdOnMachine",
		"/.jetbrains/.user_id_on_machine",
		"/.jetbrains/.device_id",
	}

	for _, k := range filesToRemove {
		exec.Command("plutil", "-remove", k, plist).Run()
	}
}
