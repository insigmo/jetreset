package reset

import (
	"fmt"
	"path/filepath"
)

func resetWindows(h string, products []string) {
	fmt.Println("Windows:")
	cleanDir(filepath.Join(h, "AppData", "Roaming", "JetBrains"), products)
	cleanDir(filepath.Join(h, "AppData", "Local", "JetBrains"), products)
	fmt.Println("Note: registry cleanup requires manual action (regedit/PowerShell).")
}
