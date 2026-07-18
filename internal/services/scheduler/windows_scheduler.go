//go:build windows

// Package scheduler installs and removes the monthly reset schedule.
package scheduler

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

const appName = "JetReset"

// Schedule creates a monthly Task Scheduler entry that runs processName.
func Schedule(processName string) {
	args := []string{"/create", "/tn", appName, "/tr", `"` + processName + `"`, "/sc", "MONTHLY", "/f"}
	err := exec.CommandContext(context.Background(), "schtasks", args...).Run() //nolint:gosec // G204: launching the scheduler is intended
	if err != nil {
		fmt.Fprintln(os.Stderr, "schtasks error:", err)
		return
	}
	fmt.Println("Scheduled: Task Scheduler (monthly)")
}

// Unschedule removes the Task Scheduler entry.
func Unschedule(_ string) {
	_ = exec.CommandContext(context.Background(), "schtasks", "/delete", "/tn", appName, "/f").Run() //nolint:gosec // G204: fixed command
	fmt.Println("Task Scheduler entry removed.")
}
