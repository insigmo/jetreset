package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/insigmo/jetreset/internal/services/reset"
	"github.com/insigmo/jetreset/internal/services/scheduler"
)

var (
	products = []string{
		"IntelliJIdea", "CLion", "PhpStorm", "GoLand", "PyCharm",
		"WebStorm", "Rider", "DataGrip", "RubyMine", "AppCode",
	}
)

func main() {
	fmt.Printf("🚀 Jetreset launched on %s\n", runtime.GOOS)
	home, _ := os.UserHomeDir()
	processName, _ := os.Executable()

	if len(os.Args) > 1 && os.Args[1] == "--stop" {
		scheduler.Unschedule(processName)
		fmt.Println("🛑 Scheduler stopped. Auto-reset disabled. Check via `crontab -l`")
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "--run-schedule" {
		fmt.Println("📅 Scheduler started. Trial will be reset automatically every month. Check via `crontab -l`")
		scheduler.Schedule(processName)
		return
	}

	reset.Reset(home, products)

	fmt.Println("✅ Done! Trial period reset for all found products.")
}
