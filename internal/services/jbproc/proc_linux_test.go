package jbproc_test

import (
	"reflect"
	"testing"

	"github.com/insigmo/jetreset/internal/services/jbproc"
)

func TestParseDesktopExec(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []string
	}{
		{"quoted path + %f", `Exec="/opt/GoLand-2026.1.2/bin/goland.sh" %f`, []string{"/opt/GoLand-2026.1.2/bin/goland.sh"}},
		{"bare path + %U", `Exec=/usr/bin/idea %U`, []string{"/usr/bin/idea"}},
		{"quoted path + flag + %f", `Exec="path with space.sh" --flag %f`, []string{"path with space.sh", "--flag"}},
		{"multiple field codes", `Exec=idea %F %U %i %c %k`, []string{"idea"}},
		{"no args", `Exec=/opt/app`, []string{"/opt/app"}},
		{"empty exec", `Exec=`, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := jbproc.ParseDesktopExec(c.in)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("ParseDesktopExec(%q) = %v, want %v", c.in, got, c.want)
			}
		})
	}
}

func TestShellSplit(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []string
	}{
		{"simple", `a b c`, []string{"a", "b", "c"}},
		{"quoted space", `"a b" c`, []string{"a b", "c"}},
		{"escaped space", `a\ b c`, []string{"a b", "c"}},
		{"empty quoted", `"" x`, []string{"", "x"}},
		{"leading/trailing space", `  a   b  `, []string{"a", "b"}},
		{"empty", ``, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := jbproc.ShellSplit(c.in)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("ShellSplit(%q) = %v, want %v", c.in, got, c.want)
			}
		})
	}
}

func TestDesktopToken(t *testing.T) {
	cases := map[string]string{
		"IntelliJIdea": "idea",
		"GoLand":       "goland",
		"PyCharm":      "pycharm",
		"DataGrip":     "datagrip",
	}
	for in, want := range cases {
		if got := jbproc.DesktopToken(in); got != want {
			t.Errorf("DesktopToken(%q) = %q, want %q", in, got, want)
		}
	}
}
