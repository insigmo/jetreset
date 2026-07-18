//go:build linux

package jbproc

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/insigmo/jetreset/internal/config"
	"github.com/insigmo/jetreset/internal/logx"
)

const (
	procDir       = "/proc"
	cmdlineFile   = "cmdline"
	pollInterval  = 150 * time.Millisecond
	ideaDesktopID = "idea"
)

// desktopFieldCodes are the Desktop Entry spec field codes stripped from Exec=.
var desktopFieldCodes = []string{"%f", "%F", "%u", "%U", "%i", "%c", "%k"} // const

// Running lists currently-running JetBrains IDEs by walking /proc.
func Running() []Proc {
	entries, err := os.ReadDir(procDir)
	if err != nil {
		return nil
	}
	var ps []Proc
	for _, e := range entries {
		pid, err := strconv.Atoi(e.Name())
		if err != nil {
			continue
		}
		raw, err := os.ReadFile(procDir + "/" + e.Name() + "/" + cmdlineFile)
		if err != nil {
			continue
		}
		cmd := strings.ReplaceAll(string(raw), "\x00", " ")
		if !strings.Contains(cmd, ideaMain) {
			continue
		}
		product := ProductFromTokens(cmd)
		if product == "" {
			continue
		}
		logx.Debugf("detected %s pid=%d", product, pid)
		ps = append(ps, Proc{product, pid})
	}
	return DedupByProduct(ps)
}

func kill(p Proc, sig syscall.Signal) error {
	logx.Debugf("signal %d to pid=%d (%s)", sig, p.PID, p.Product)
	err := syscall.Kill(p.PID, sig)
	if errors.Is(err, syscall.ESRCH) {
		return nil
	}
	return err
}

// Kill sends a graceful terminate signal (SIGTERM).
func Kill(p Proc) error { return kill(p, syscall.SIGTERM) }

// ForceKill sends SIGKILL after a graceful close timed out.
func ForceKill(p Proc) error { return kill(p, syscall.SIGKILL) }

// Wait blocks until the process is gone or the timeout elapses.
func Wait(p Proc, timeout time.Duration) error {
	dir := procDir + "/" + strconv.Itoa(p.PID)
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			return nil
		}
		time.Sleep(pollInterval)
	}
	return ErrTimeout
}

// HasGUI reports whether a graphical session is available.
func HasGUI() bool {
	return os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != ""
}

// Relaunch starts the IDE again from its .desktop launcher; fire-and-forget.
func Relaunch(p Proc) error {
	exe, args, err := findDesktopLaunch(p.Product)
	if err != nil {
		token := DesktopToken(p.Product)
		logx.Debugf("gtk-launch jetbrains-%s", token)
		return exec.CommandContext(context.Background(), "gtk-launch", "jetbrains-"+token).Start() //nolint:gosec // G204: launching the IDE is intended
	}
	logx.Debugf("relaunch %s via %s %v", p.Product, exe, args)
	cmd := exec.CommandContext(context.Background(), exe, args...) //nolint:gosec // G204: launching the IDE is intended
	cmd.Env = os.Environ()
	return cmd.Start()
}

// DesktopToken returns the lowercase token used in JetBrains .desktop names.
func DesktopToken(product string) string {
	if product == config.IntelliJIdea {
		return ideaDesktopID
	}
	return strings.ToLower(product)
}

func findDesktopLaunch(product string) (string, []string, error) {
	home, _ := os.UserHomeDir()
	dirs := []string{
		filepath.Join(home, ".local", "share", "applications"),
		"/usr/share/applications",
		"/usr/local/share/applications",
	}
	token := DesktopToken(product)
	for _, dir := range dirs {
		if exe, args, ok := readDesktopExec(filepath.Join(dir, "jetbrains-"+token+".desktop")); ok {
			return exe, args, nil
		}
		matches, _ := filepath.Glob(filepath.Join(dir, "jetbrains-*"+token+"*.desktop"))
		for _, m := range matches {
			if exe, args, ok := readDesktopExec(m); ok {
				return exe, args, nil
			}
		}
	}
	return "", nil, fmt.Errorf("%w: %s", ErrNoLauncher, product)
}

func readDesktopExec(path string) (string, []string, bool) {
	b, err := os.ReadFile(path) //nolint:gosec // G304: path from a glob over application dirs
	if err != nil {
		return "", nil, false
	}
	for line := range strings.SplitSeq(string(b), "\n") {
		if line = strings.TrimSpace(line); strings.HasPrefix(line, "Exec=") {
			if parts := ParseDesktopExec(line); len(parts) > 0 {
				logx.Debugf("desktop %s Exec=%v", path, parts)
				return parts[0], parts[1:], true
			}
		}
	}
	return "", nil, false
}

// ParseDesktopExec parses a desktop entry "Exec=..." line into an argv.
func ParseDesktopExec(line string) []string {
	line = strings.TrimSpace(strings.TrimPrefix(line, "Exec="))
	for _, fc := range desktopFieldCodes {
		line = strings.ReplaceAll(line, fc, "")
	}
	return ShellSplit(strings.TrimSpace(line))
}

// ShellSplit splits a command line into argv, respecting quotes and escapes.
func ShellSplit(s string) []string {
	var parts []string
	var cur strings.Builder
	inQ, hasC := false, false
	flush := func() {
		if hasC {
			parts = append(parts, cur.String())
			cur.Reset()
			hasC = false
		}
	}
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c == '\\' && i+1 < len(s):
			cur.WriteByte(s[i+1])
			hasC = true
			i++
		case c == '"':
			inQ = !inQ
			hasC = true
		case (c == ' ' || c == '\t') && !inQ:
			flush()
		default:
			cur.WriteByte(c)
			hasC = true
		}
	}
	flush()
	return parts
}
