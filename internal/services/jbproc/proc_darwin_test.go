package jbproc_test

import (
	"testing"

	"github.com/insigmo/jetreset/internal/services/jbproc"
)

func TestMacAppName(t *testing.T) {
	cases := map[string]string{
		"IntelliJIdea": "IntelliJ IDEA",
		"GoLand":       "GoLand",
		"PyCharm":      "PyCharm",
	}
	for in, want := range cases {
		if got := jbproc.MacAppName(in); got != want {
			t.Errorf("MacAppName(%q) = %q, want %q", in, got, want)
		}
	}
}
