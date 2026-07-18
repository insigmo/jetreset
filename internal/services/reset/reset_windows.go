package reset

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func Reset(home string, products []string) {
	os.RemoveAll(filepath.Join(home, "AppData", "Roaming", "JetBrains"))
	os.RemoveAll(filepath.Join(home, "AppData", "Local", "JetBrains"))
	cleanRegistry()
}

func cleanRegistry() {
	err := exec.Command("reg", "delete", `HKEY_CURRENT_USER\Software\JavaSoft`, "/f").Run()
	if err == nil {
		fmt.Println("🧹 Registry: cleaned")
	}
}
