package reset

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var evlRe = regexp.MustCompile(`<property name="evlsprt[^"]*" value="[^"]*"`)

func removeLine(path string, re *regexp.Regexp) {
	b, err := os.ReadFile(path)
	if err != nil {
		return
	}
	var out []string
	changed := false
	for _, line := range strings.Split(string(b), "\n") {
		if re.MatchString(line) {
			changed = true
			continue
		}
		out = append(out, line)
	}
	if changed {
		os.WriteFile(path, []byte(strings.Join(out, "\n")), 0644)
	}
}

func cleanDir(base string, products []string) {
	for _, p := range products {
		dirs, _ := filepath.Glob(filepath.Join(base, p+"*"))
		for _, d := range dirs {
			os.RemoveAll(filepath.Join(d, "eval"))
			removeLine(filepath.Join(d, "options", "other.xml"), evlRe)
		}
	}
}
