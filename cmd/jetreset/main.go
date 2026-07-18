// Command jetreset resets the trial period of JetBrains IDEs.
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/insigmo/jetreset/internal/config"
	"github.com/insigmo/jetreset/internal/logx"
	"github.com/insigmo/jetreset/internal/services/jbproc"
	"github.com/insigmo/jetreset/internal/services/reset"
	"github.com/insigmo/jetreset/internal/services/scheduler"
)

const (
	gracefulTimeout = 8 * time.Second
	forceTimeout    = 3 * time.Second
)

func main() {
	runSchedule := flag.Bool("run-schedule", false, "Enable automatic monthly reset")
	stopSchedule := flag.Bool("stop", false, "Disable automatic monthly reset")
	verbose := flag.Bool("v", false, "Show debug logs")

	flag.Parse()
	logx.SetVerbose(*verbose)

	fmt.Printf("🚀 Jetreset launched on %s\n", runtime.GOOS)

	home, _ := os.UserHomeDir()
	exe, _ := os.Executable()

	logx.Debugf("home=%s executable=%s", home, exe)

	switch {
	case *stopSchedule:
		scheduler.Unschedule(exe)
		fmt.Println("🛑 Scheduler stopped. Auto-reset disabled. Check via `crontab -l`")
	case *runSchedule:
		fmt.Println("📅 Scheduler started. Trial will be reset automatically every month. Check via `crontab -l`")
		scheduler.Schedule(exe)
	default:
		runReset(home)
	}
}

func runReset(home string) {
	running := jbproc.Running()
	logx.Debugf("detected %d running IDE(s)", len(running))

	if len(running) > 0 && isInteractive() {
		fmt.Println("⚠️  Closing running IDEs — any unsaved changes will be lost.")
	}
	closeRunning(running)

	reset.Reset(home, config.Products)
	fmt.Println("🧹 Trial data wiped.")

	relaunchRunning(running)

	fmt.Println("✅ Done! Trial period reset for all found products.")
}

func closeRunning(running []jbproc.Proc) {
	for _, p := range running {
		fmt.Printf("🛑 Closing %s...\n", p.Product)
		err := jbproc.Kill(p)
		if err != nil {
			logx.Debugf("kill %s: %v", p.Product, err)
			continue
		}
		err = jbproc.Wait(p, gracefulTimeout)
		if err != nil {
			logx.Debugf("%s force-killing: %v", p.Product, err)
			_ = jbproc.ForceKill(p)
			_ = jbproc.Wait(p, forceTimeout)
		}
	}
}

func relaunchRunning(running []jbproc.Proc) {
	if len(running) == 0 {
		return
	}
	if !jbproc.HasGUI() {
		fmt.Println("ℹ️  No graphical session detected; skipping relaunch. Run jetreset from your desktop session to reopen IDEs.")
		return
	}
	for _, p := range running {
		fmt.Printf("🚀 Relaunching %s...\n", p.Product)
		err := jbproc.Relaunch(p)
		if err != nil {
			logx.Debugf("relaunch %s: %v", p.Product, err)
			fmt.Printf("⚠️  Could not relaunch %s (open it manually).\n", p.Product)
		}
	}
}

func isInteractive() bool {
	fi, err := os.Stdin.Stat()
	return err == nil && fi.Mode()&os.ModeCharDevice != 0
}
