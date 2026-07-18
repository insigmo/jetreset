//go:build darwin

package jbproc

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/insigmo/jetreset/internal/config"
	"github.com/insigmo/jetreset/internal/logx"
)

const pollInterval = 150 * time.Millisecond

// Running lists currently-running JetBrains IDEs via ps.
func Running() []Proc {
	out, err := exec.CommandContext(context.Background(), "ps", "-ax", "-o", "pid=,command=").Output() //nolint:gosec // G204: fixed command
	if err != nil {
		return nil
	}
	return ProcsFromPS(string(out))
}

func kill(p Proc, sig syscall.Signal) error {
	logx.Tracef("signal %d to pid=%d (%s)", sig, p.PID, p.Product)
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
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		err := syscall.Kill(p.PID, 0)
		if errors.Is(err, syscall.ESRCH) {
			return nil
		}
		time.Sleep(pollInterval)
	}
	return ErrTimeout
}

// HasGUI reports whether a graphical session is available.
func HasGUI() bool { return true }

// Relaunch starts the IDE again via LaunchServices (open).
func Relaunch(p Proc) error {
	name := MacAppName(p.Product)
	err := exec.CommandContext(context.Background(), "open", "-a", name).Start() //nolint:gosec // G204: launching the IDE is intended
	if err == nil {
		logx.Tracef("relaunch %s via open -a %q", p.Product, name)
		return nil
	}
	if m, _ := filepath.Glob(filepath.Join("/Applications", name+"*.app")); len(m) > 0 {
		logx.Tracef("relaunch %s via open %s", p.Product, m[0])
		return exec.CommandContext(context.Background(), "open", m[0]).Start() //nolint:gosec // G204: launching the IDE is intended
	}
	return fmt.Errorf("%w: %s", ErrNoLauncher, p.Product)
}

// MacAppName maps a canonical product name to its macOS application name.
func MacAppName(product string) string {
	if product == config.IntelliJIdea {
		return "IntelliJ IDEA"
	}
	return product
}
