package logx_test

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/insigmo/jetreset/internal/logx"
)

func capture(t *testing.T, fn func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}

	old := os.Stderr
	os.Stderr = w

	defer func() { os.Stderr = old }()

	fn()

	err = w.Close()
	if err != nil {
		return ""
	}

	out, _ := io.ReadAll(r)
	return string(out)
}

func TestSilentByDefault(t *testing.T) {
	logx.SetVerbose(false)
	got := capture(t, func() {
		logx.Debugf("k=%s", "v")
	})
	if got != "" {
		t.Errorf("expected no output when disabled, got %q", got)
	}
}

func TestTraceEmits(t *testing.T) {
	logx.SetVerbose(true)
	t.Cleanup(func() { logx.SetVerbose(false) })
	got := capture(t, func() {
		logx.Debugf("k=%s", "v")
	})
	if !strings.Contains(got, "[debug]") || !strings.Contains(got, "k=v") {
		t.Errorf("unexpected output: %q", got)
	}
}
