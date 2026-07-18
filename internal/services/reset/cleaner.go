// Package reset wipes JetBrains trial/eval/license state.
package reset

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/insigmo/jetreset/internal/logx"
)

const configFileMode = 0o600

// EvlRe matches the evlsprt* trial properties written to options/other.xml.
var EvlRe = regexp.MustCompile(`<property name="evlsprt[^"]*" value="[^"]*"`)

// RemoveLine drops every line matching re from the file at path.
func RemoveLine(path string, re *regexp.Regexp) {
	b, err := os.ReadFile(path) //nolint:gosec // G304: path is a known config file
	if err != nil {
		return
	}
	var out []string
	removed := 0
	for line := range strings.SplitSeq(string(b), "\n") {
		if re.MatchString(line) {
			removed++
			continue
		}
		out = append(out, line)
	}
	if removed == 0 {
		return
	}
	err = os.WriteFile(path, []byte(strings.Join(out, "\n")), configFileMode)
	if err != nil {
		logx.Debugf("RemoveLine %s: %v", path, err)
		return
	}
	logx.Debugf("RemoveLine %s: stripped %d line(s)", path, removed)
}

// CleanDir removes the eval/ dir and evlsprt lines for each product under base.
func CleanDir(base string, products []string) {
	for _, p := range products {
		dirs, _ := filepath.Glob(filepath.Join(base, p+"*"))
		for _, d := range dirs {
			evalDir := filepath.Join(d, "eval")
			_, statErr := os.Stat(evalDir)
			if statErr == nil {
				rmErr := os.RemoveAll(evalDir)
				if rmErr != nil {
					logx.Debugf("CleanDir: %s: %v", evalDir, rmErr)
				} else {
					logx.Debugf("CleanDir: removed %s", evalDir)
				}
			}
			RemoveLine(filepath.Join(d, "options", "other.xml"), EvlRe)
		}
	}
}
