//go:build windows

package reset_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/insigmo/jetreset/internal/services/reset"
)

func TestResetWindows(t *testing.T) {
	home := t.TempDir()
	products := []string{"IntelliJIdea", "GoLand"}

	roaming := filepath.Join(home, "AppData", "Roaming", "JetBrains")
	local := filepath.Join(home, "AppData", "Local", "JetBrains")
	os.MkdirAll(roaming, 0o755)
	os.MkdirAll(local, 0o755)

	reset.Reset(home, products)

	t.Run("removes AppData/Roaming/JetBrains", func(t *testing.T) {
		_, err := os.Stat(roaming)
		if !os.IsNotExist(err) {
			t.Error("AppData/Roaming/JetBrains should be removed")
		}
	})

	t.Run("removes AppData/Local/JetBrains", func(t *testing.T) {
		_, err := os.Stat(local)
		if !os.IsNotExist(err) {
			t.Error("AppData/Local/JetBrains should be removed")
		}
	})
}

func TestResetWindowsMissingDirs(t *testing.T) {
	reset.Reset(t.TempDir(), []string{"IntelliJIdea", "GoLand"})
}
