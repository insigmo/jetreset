//go:build darwin

package reset_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/insigmo/jetreset/internal/services/reset"
)

func TestResetDarwin(t *testing.T) {
	home := t.TempDir()
	products := []string{"IntelliJIdea", "GoLand"}

	libPrefs := filepath.Join(home, "Library", "Preferences")
	appSupport := filepath.Join(home, "Library", "Application Support", "JetBrains")

	for _, name := range []string{"IntelliJIdea2024.1", "GoLand2023.3"} {
		os.MkdirAll(filepath.Join(libPrefs, name, "eval"), 0o755)
		os.MkdirAll(filepath.Join(appSupport, name, "eval"), 0o755)
	}

	reset.Reset(home, products)

	t.Run("removes eval dirs from Library/Preferences", func(t *testing.T) {
		for _, name := range []string{"IntelliJIdea2024.1", "GoLand2023.3"} {
			evalDir := filepath.Join(libPrefs, name, "eval")
			_, err := os.Stat(evalDir)
			if !os.IsNotExist(err) {
				t.Errorf("Library/Preferences/%s/eval should be removed", name)
			}
		}
	})

	t.Run("removes eval dirs from Application Support/JetBrains", func(t *testing.T) {
		for _, name := range []string{"IntelliJIdea2024.1", "GoLand2023.3"} {
			evalDir := filepath.Join(appSupport, name, "eval")
			_, err := os.Stat(evalDir)
			if !os.IsNotExist(err) {
				t.Errorf("Application Support/JetBrains/%s/eval should be removed", name)
			}
		}
	})
}

func TestResetDarwinMissingDirs(t *testing.T) {
	reset.Reset(t.TempDir(), []string{"IntelliJIdea", "GoLand"})
}
