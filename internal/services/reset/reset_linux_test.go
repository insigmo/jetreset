package reset_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/insigmo/jetreset/internal/services/reset"
)

func TestResetLinux(t *testing.T) {
	home := t.TempDir()
	products := []string{"IntelliJIdea", "GoLand"}

	jbConfig := filepath.Join(home, ".config", "JetBrains")
	for _, name := range []string{"IntelliJIdea2024.1", "GoLand2023.3"} {
		_ = os.MkdirAll(filepath.Join(jbConfig, name, "eval"), 0o755)
	}

	jprefs := filepath.Join(home, ".java", ".userPrefs")
	_ = os.MkdirAll(filepath.Join(jprefs, "jetbrains"), 0o755)

	_ = os.WriteFile(filepath.Join(jprefs, "prefs.xml"), []byte(strings.Join([]string{
		`<?xml version="1.0" encoding="UTF-8" standalone="no"?>`,
		`<!DOCTYPE map SYSTEM "http://java.sun.com/dtd/preferences.dtd">`,
		`<map MAP_XML_VERSION="1.0">`,
		`  <entry key="JetBrains.UserIdOnMachine" value="abc-uuid-123"/>`,
		`  <entry key="other.key" value="keep"/>`,
		`</map>`,
	}, "\n")), 0o644)

	_ = os.WriteFile(filepath.Join(jprefs, "jetbrains", "prefs.xml"), []byte(strings.Join([]string{
		`<?xml version="1.0" encoding="UTF-8" standalone="no"?>`,
		`<!DOCTYPE map SYSTEM "http://java.sun.com/dtd/preferences.dtd">`,
		`<map MAP_XML_VERSION="1.0">`,
		`  <entry key="device_id" value="dev-456"/>`,
		`  <entry key="user_id_on_machine" value="uid-789"/>`,
		`  <entry key="other.key" value="keep"/>`,
		`</map>`,
	}, "\n")), 0o644)

	reset.Reset(home, products)

	t.Run("removes eval directories", func(t *testing.T) {
		for _, name := range []string{"IntelliJIdea2024.1", "GoLand2023.3"} {
			evalDir := filepath.Join(jbConfig, name, "eval")
			_, err := os.Stat(evalDir)
			if !os.IsNotExist(err) {
				t.Errorf("%s/eval should be removed", name)
			}
		}
	})

	t.Run("removes JetBrains.UserIdOnMachine from prefs.xml", func(t *testing.T) {
		got, _ := os.ReadFile(filepath.Join(jprefs, "prefs.xml"))
		if strings.Contains(string(got), "JetBrains.UserIdOnMachine") {
			t.Errorf("JetBrains.UserIdOnMachine should be removed, got:\n%s", got)
		}
		if !strings.Contains(string(got), "other.key") {
			t.Error("unrelated entries should be kept")
		}
	})

	t.Run("removes device_id and user_id_on_machine from jetbrains/prefs.xml", func(t *testing.T) {
		got, _ := os.ReadFile(filepath.Join(jprefs, "jetbrains", "prefs.xml"))
		if strings.Contains(string(got), `key="device_id"`) {
			t.Errorf("device_id should be removed, got:\n%s", got)
		}
		if strings.Contains(string(got), `key="user_id_on_machine"`) {
			t.Errorf("user_id_on_machine should be removed, got:\n%s", got)
		}
		if !strings.Contains(string(got), "other.key") {
			t.Error("unrelated entries should be kept")
		}
	})
}

func TestResetLinuxMissingDirs(t *testing.T) {
	reset.Reset(t.TempDir(), []string{"IntelliJIdea", "GoLand"})
}
