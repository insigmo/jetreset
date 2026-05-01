package reset

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRemoveLine(t *testing.T) {
	t.Run("removes matching lines", func(t *testing.T) {
		f := filepath.Join(t.TempDir(), "other.xml")
		content := strings.Join([]string{
			`<application>`,
			`  <property name="evlsprtInExpiration" value="1234567890"/>`,
			`  <property name="evlsprtVersion" value="3"/>`,
			`  <property name="keep.this" value="yes"/>`,
			`</application>`,
		}, "\n")
		os.WriteFile(f, []byte(content), 0644)

		removeLine(f, evlRe)

		got, _ := os.ReadFile(f)
		if strings.Contains(string(got), "evlsprt") {
			t.Errorf("expected evlsprt lines removed, got:\n%s", got)
		}
		if !strings.Contains(string(got), "keep.this") {
			t.Error("expected non-matching lines to be kept")
		}
	})

	t.Run("does not modify file when no match", func(t *testing.T) {
		f := filepath.Join(t.TempDir(), "other.xml")
		content := `<property name="keep.this" value="yes"/>`
		os.WriteFile(f, []byte(content), 0644)

		info1, _ := os.Stat(f)
		removeLine(f, evlRe)
		info2, _ := os.Stat(f)

		if !info1.ModTime().Equal(info2.ModTime()) {
			t.Error("file should not be modified when no lines match")
		}
	})

	t.Run("silently handles missing file", func(t *testing.T) {
		removeLine(filepath.Join(t.TempDir(), "nonexistent.xml"), evlRe)
	})
}

func TestCleanDir(t *testing.T) {
	base := t.TempDir()
	products := []string{"IntelliJIdea", "GoLand"}

	for _, name := range []string{"IntelliJIdea2024.1", "GoLand2023.3"} {
		dir := filepath.Join(base, name)
		os.MkdirAll(filepath.Join(dir, "eval"), 0755)
		os.MkdirAll(filepath.Join(dir, "options"), 0755)

		otherXML := filepath.Join(dir, "options", "other.xml")
		os.WriteFile(otherXML, []byte(strings.Join([]string{
			`<application>`,
			`  <property name="evlsprtInExpiration" value="9999"/>`,
			`  <property name="keep.this" value="yes"/>`,
			`</application>`,
		}, "\n")), 0644)
	}

	cleanDir(base, products)

	for _, name := range []string{"IntelliJIdea2024.1", "GoLand2023.3"} {
		if _, err := os.Stat(filepath.Join(base, name, "eval")); !os.IsNotExist(err) {
			t.Errorf("%s/eval/ should be removed", name)
		}
		got, _ := os.ReadFile(filepath.Join(base, name, "options", "other.xml"))
		if strings.Contains(string(got), "evlsprt") {
			t.Errorf("%s/options/other.xml: evlsprt entries should be removed", name)
		}
		if !strings.Contains(string(got), "keep.this") {
			t.Errorf("%s/options/other.xml: unrelated entries should be kept", name)
		}
	}
}

func TestCleanDirIgnoresUnrelatedProducts(t *testing.T) {
	base := t.TempDir()
	evalDir := filepath.Join(base, "SomeOtherApp2024", "eval")
	os.MkdirAll(evalDir, 0755)

	cleanDir(base, []string{"IntelliJIdea", "GoLand"})

	if _, err := os.Stat(evalDir); os.IsNotExist(err) {
		t.Error("eval/ of unrelated product should not be removed")
	}
}

func TestCleanDirMissingBaseDir(t *testing.T) {
	cleanDir(filepath.Join(t.TempDir(), "nonexistent"), []string{"IntelliJIdea"})
}
