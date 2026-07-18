//go:build windows

package reset

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResetWindows(t *testing.T) {
	home := t.TempDir()
	products := []string{"IntelliJIdea", "GoLand"}

	roaming := filepath.Join(home, "AppData", "Roaming", "JetBrains")
	local := filepath.Join(home, "AppData", "Local", "JetBrains")
	os.MkdirAll(roaming, 0755)
	os.MkdirAll(local, 0755)

	Reset(home, products)

	t.Run("removes AppData/Roaming/JetBrains", func(t *testing.T) {
		if _, err := os.Stat(roaming); !os.IsNotExist(err) {
			t.Error("AppData/Roaming/JetBrains should be removed")
		}
	})

	t.Run("removes AppData/Local/JetBrains", func(t *testing.T) {
		if _, err := os.Stat(local); !os.IsNotExist(err) {
			t.Error("AppData/Local/JetBrains should be removed")
		}
	})
}

func TestResetWindowsMissingDirs(t *testing.T) {
	// Reset on an empty home directory must not crash
	// (os.RemoveAll on nonexistent path is a no-op; cleanRegistry silently fails)
	Reset(t.TempDir(), []string{"IntelliJIdea", "GoLand"})
}
