package reset

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func resetWindows(h string) {
	os.RemoveAll(filepath.Join(h, "AppData", "Roaming", "JetBrains"))
	os.RemoveAll(filepath.Join(h, "AppData", "Local", "JetBrains"))
	cleanRegistry()
}

func cleanRegistry() {
	err := exec.Command("reg", "delete", `HKEY_CURRENT_USER\Software\JavaSoft`, "/f").Run()
	if err == nil {
		fmt.Println("🧹 Registry: cleaned")
	}
}
