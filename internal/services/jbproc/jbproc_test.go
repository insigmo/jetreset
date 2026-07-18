package jbproc_test

import (
	"strings"
	"testing"

	"github.com/insigmo/jetreset/internal/services/jbproc"
)

func TestCanonicalize(t *testing.T) {
	cases := map[string]string{
		"GoLand": "GoLand", "Idea": "IntelliJIdea", "IdeaC": "IntelliJIdea",
		"PyCharm": "PyCharm", "PyCharmC": "PyCharm", "CLion": "CLion",
		"WebStorm": "WebStorm", "Unknown": "", "": "",
	}
	for in, want := range cases {
		if got := jbproc.Canonicalize(in); got != want {
			t.Errorf("Canonicalize(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestProductFromTokens(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"selector GoLand", "-Didea.paths.selector=GoLand2026.2 com.intellij.idea.Main", "GoLand"},
		{"prefix GoLand", "-Didea.platform.prefix=GoLand com.intellij.idea.Main", "GoLand"},
		{"install path", "/opt/GoLand-2026.1.2/jbr/bin/java com.intellij.idea.Main", "GoLand"},
		{"IDEA ultimate", "-Didea.paths.selector=Idea2024.1 com.intellij.idea.Main", "IntelliJIdea"},
		{"IDEA community", "-Didea.paths.selector=IdeaC2024.1 com.intellij.idea.Main", "IntelliJIdea"},
		{"PyCharm", "-Didea.paths.selector=PyCharm2024.1 com.intellij.idea.Main", "PyCharm"},
		{"no product", "com.intellij.idea.Main", ""},
		{"not an IDE", "java -cp app.jar foo.Main", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := jbproc.ProductFromTokens(c.in); got != c.want {
				t.Errorf("ProductFromTokens() = %q, want %q", got, c.want)
			}
		})
	}
}

func TestDedupByProduct(t *testing.T) {
	in := []jbproc.Proc{
		{Product: "GoLand", PID: 500},
		{Product: "GoLand", PID: 501},
		{Product: "PyCharm", PID: 100},
		{Product: "", PID: 999},
	}
	got := jbproc.DedupByProduct(in)
	want := []jbproc.Proc{
		{Product: "PyCharm", PID: 100},
		{Product: "GoLand", PID: 500},
	}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("idx %d: got %v, want %v", i, got[i], want[i])
		}
	}
}

func TestProcsFromPS(t *testing.T) {
	in := strings.Join([]string{
		" 19395 /opt/GoLand-2026.1.2/jbr/bin/java -Didea.paths.selector=GoLand2026.2 com.intellij.idea.Main",
		" 1234 /usr/bin/java -cp app.jar foo.Main",
		" 555 /jdk/bin/java -Didea.paths.selector=PyCharm2024.1 com.intellij.idea.Main",
		" 556 /jdk/bin/java -Didea.paths.selector=PyCharm2024.1 com.intellij.idea.Main",
		" garbage line",
		"",
	}, "\n")
	got := jbproc.ProcsFromPS(in)
	want := []jbproc.Proc{
		{Product: "PyCharm", PID: 555},
		{Product: "GoLand", PID: 19395},
	}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("idx %d: got %v, want %v", i, got[i], want[i])
		}
	}
}

func TestProcsFromTasklist(t *testing.T) {
	in := strings.Join([]string{
		`"goland64.exe","1234","Console","1","1,234 K"`,
		`"pycharm64.exe","5678","Console","1","2,000 K"`,
		`"chrome.exe","9999","Console","1","500 K"`,
		`garbage`,
		"",
	}, "\n")
	got := jbproc.ProcsFromTasklist(in)
	want := []jbproc.Proc{
		{Product: "GoLand", PID: 1234},
		{Product: "PyCharm", PID: 5678},
	}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("idx %d: got %v, want %v", i, got[i], want[i])
		}
	}
}
