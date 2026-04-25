package main

import (
	"fmt"
	"os"

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
	home, _ := os.UserHomeDir()
	processName, _ := os.Executable()
	if len(os.Args) > 1 && os.Args[1] == "--stop" {
		scheduler.Unschedule(processName)
		return
	}
	reset.Reset(home, products)
	scheduler.Schedule(processName)
	fmt.Println("Done.")
}
