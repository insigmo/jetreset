package reset

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/insigmo/jetreset/internal/logx"
)

// Reset wipes JetBrains trial state for the given products on Windows.
func Reset(home string, products []string) {
	logx.Debugf("reset: AppData + registry")
	_ = os.RemoveAll(filepath.Join(home, "AppData", "Roaming", "JetBrains"))
	_ = os.RemoveAll(filepath.Join(home, "AppData", "Local", "JetBrains"))
	cleanRegistry()
}

func cleanRegistry() {
	err := exec.CommandContext(context.Background(), "reg", "delete", `HKEY_CURRENT_USER\Software\JavaSoft`, "/f").Run() //nolint:gosec // G204: fixed command
	if err == nil {
		fmt.Println("🧹 Registry: cleaned")
	} else {
		logx.Debugf("registry: nothing to clean")
	}
}
