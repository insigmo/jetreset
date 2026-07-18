//go:build windows

package jbproc

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/insigmo/jetreset/internal/config"
	"github.com/insigmo/jetreset/internal/logx"
)

const pollInterval = 200 * time.Millisecond

// Running lists currently-running JetBrains IDEs via tasklist.
func Running() []Proc {
	out, err := exec.CommandContext(context.Background(), "tasklist", "/FO", "CSV", "/NH").Output() //nolint:gosec // G204: fixed command
	if err != nil {
		return nil
	}
	return ProcsFromTasklist(string(out))
}

func taskkill(p Proc, force bool) error {
	args := []string{"/PID", strconv.Itoa(p.PID), "/T"}
	if force {
		args = append(args, "/F")
	}
	logx.Tracef("taskkill %v (%s)", args, p.Product)
	return exec.CommandContext(context.Background(), "taskkill", args...).Run() //nolint:gosec // G204: terminating the IDE is intended
}

// Kill gracefully closes the IDE (and its child JVM tree) via taskkill /T.
func Kill(p Proc) error { return taskkill(p, false) }

// ForceKill force-terminates the process tree after a graceful close timed out.
func ForceKill(p Proc) error { return taskkill(p, true) }

// Wait blocks until the process is gone or the timeout elapses.
func Wait(p Proc, timeout time.Duration) error {
	pid := strconv.Itoa(p.PID)
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		out, _ := exec.CommandContext(context.Background(), "tasklist", "/FI", "PID eq "+pid, "/FO", "CSV", "/NH").Output() //nolint:gosec // G204: fixed command
		if !strings.Contains(string(out), pid) {
			return nil
		}
		time.Sleep(pollInterval)
	}
	return ErrTimeout
}

// HasGUI reports whether a graphical session is available.
func HasGUI() bool { return true }

// Relaunch starts the IDE again via its Start Menu shortcut or launcher exe.
func Relaunch(p Proc) error {
	name := winStartMenuName(p.Product)
	root := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs", "JetBrains")
	for _, g := range []string{
		filepath.Join(root, name+"*.lnk"),
		filepath.Join(root, name+"*", "*.lnk"),
	} {
		m, _ := filepath.Glob(g)
		if len(m) > 0 {
			logx.Tracef("relaunch %s via start %s", p.Product, m[0])
			return exec.CommandContext(context.Background(), "cmd", "/c", "start", "", m[0]).Start() //nolint:gosec // G204: launching the IDE is intended
		}
	}
	if exe := config.ProductExe[p.Product]; exe != "" {
		path, err := exec.LookPath(exe)
		if err == nil {
			logx.Tracef("relaunch %s via %s", p.Product, path)
			return exec.CommandContext(context.Background(), path).Start() //nolint:gosec // G204: launching the IDE is intended
		}
		for _, c := range winInstallCandidates(p.Product, exe) {
			_, statErr := os.Stat(c)
			if statErr != nil {
				continue
			}
			logx.Tracef("relaunch %s via %s", p.Product, c)
			return exec.CommandContext(context.Background(), c).Start() //nolint:gosec // G204: launching the IDE is intended
		}
	}
	return fmt.Errorf("%w: %s", ErrNoLauncher, p.Product)
}

func winStartMenuName(product string) string {
	if product == config.IntelliJIdea {
		return "IntelliJ IDEA"
	}
	return product
}

func winInstallCandidates(product, exe string) []string {
	return []string{
		filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", product, exe),
		filepath.Join(os.Getenv("ProgramFiles"), "JetBrains", product, exe),
		filepath.Join(os.Getenv("ProgramFiles(x86)"), "JetBrains", product, exe),
	}
}
