//go:build windows

package scheduler

import (
	"fmt"
	"os"
	"os/exec"
)

const appName = "JetReset"

func Schedule(processName string) {
	err := exec.Command("schtasks",
		"/create", "/tn", appName,
		"/tr", `"`+processName+`"`,
		"/sc", "MONTHLY", "/f",
	).Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "schtasks error:", err)
		return
	}
	fmt.Println("Scheduled: Task Scheduler (monthly)")
}

func Unschedule(_ string) {
	exec.Command("schtasks", "/delete", "/tn", appName, "/f").Run()
	fmt.Println("Task Scheduler entry removed.")
}
