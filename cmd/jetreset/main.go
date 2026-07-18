package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/insigmo/jetreset/internal/services/reset"
	"github.com/insigmo/jetreset/internal/services/scheduler"
)

var products = []string{
	"IntelliJIdea", "CLion", "PhpStorm", "GoLand", "PyCharm",
	"WebStorm", "Rider", "DataGrip", "RubyMine", "AppCode",
}

func main() {
	runSchedule := flag.Bool("run-schedule", false, "Enable automatic monthly reset")
	stopSchedule := flag.Bool("stop", false, "Disable automatic monthly reset")
	flag.Parse()

	fmt.Printf("🚀 Jetreset launched on %s\n", runtime.GOOS)
	home, _ := os.UserHomeDir()
	processName, _ := os.Executable()

	switch {
	case *stopSchedule:
		scheduler.Unschedule(processName)
		fmt.Println("🛑 Scheduler stopped. Auto-reset disabled. Check via `crontab -l`")
	case *runSchedule:
		fmt.Println("📅 Scheduler started. Trial will be reset automatically every month. Check via `crontab -l`")
		scheduler.Schedule(processName)
	default:
		reset.Reset(home, products)
		fmt.Println("✅ Done! Trial period reset for all found products.")
	}
}
