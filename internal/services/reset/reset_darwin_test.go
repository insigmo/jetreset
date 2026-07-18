package reset

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResetDarwin(t *testing.T) {
	home := t.TempDir()
	products := []string{"IntelliJIdea", "GoLand"}

	// ~/Library/Preferences/<Product><version>/eval/
	libPrefs := filepath.Join(home, "Library", "Preferences")
	// ~/Library/Application Support/JetBrains/<Product><version>/eval/
	appSupport := filepath.Join(home, "Library", "Application Support", "JetBrains")

	for _, name := range []string{"IntelliJIdea2024.1", "GoLand2023.3"} {
		os.MkdirAll(filepath.Join(libPrefs, name, "eval"), 0755)
		os.MkdirAll(filepath.Join(appSupport, name, "eval"), 0755)
	}

	Reset(home, products)

	t.Run("removes eval dirs from Library/Preferences", func(t *testing.T) {
		for _, name := range []string{"IntelliJIdea2024.1", "GoLand2023.3"} {
			evalDir := filepath.Join(libPrefs, name, "eval")
			if _, err := os.Stat(evalDir); !os.IsNotExist(err) {
				t.Errorf("Library/Preferences/%s/eval should be removed", name)
			}
		}
	})

	t.Run("removes eval dirs from Application Support/JetBrains", func(t *testing.T) {
		for _, name := range []string{"IntelliJIdea2024.1", "GoLand2023.3"} {
			evalDir := filepath.Join(appSupport, name, "eval")
			if _, err := os.Stat(evalDir); !os.IsNotExist(err) {
				t.Errorf("Application Support/JetBrains/%s/eval should be removed", name)
			}
		}
	})
}

func TestResetDarwinMissingDirs(t *testing.T) {
	// Reset on an empty home directory must not crash
	// (plutil calls will silently fail when the plist does not exist)
	Reset(t.TempDir(), []string{"IntelliJIdea", "GoLand"})
}
