//go:build darwin || linux

// Package scheduler installs and removes the monthly reset schedule.
package scheduler

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Schedule adds a monthly crontab entry that runs processName.
func Schedule(processName string) {
	out, _ := exec.CommandContext(context.Background(), "crontab", "-l").Output()
	if strings.Contains(string(out), processName) {
		fmt.Println("Already in crontab, skipping.")
		return
	}
	entry := strings.TrimRight(string(out), "\n") + "\n@monthly " + processName + "\n"
	cmd := exec.CommandContext(context.Background(), "crontab", "-")
	cmd.Stdin = strings.NewReader(entry)
	err := cmd.Run()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "crontab error:", err)
		return
	}
	fmt.Println("Scheduled: crontab @monthly")
}

// Unschedule removes the crontab entry for processName.
func Unschedule(processName string) {
	out, err := exec.CommandContext(context.Background(), "crontab", "-l").Output()
	if err != nil {
		fmt.Println("No crontab found.")
		return
	}
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(string(out)))
	for sc.Scan() {
		if !strings.Contains(sc.Text(), processName) {
			lines = append(lines, sc.Text())
		}
	}
	cmd := exec.CommandContext(context.Background(), "crontab", "-")
	cmd.Stdin = strings.NewReader(strings.Join(lines, "\n") + "\n")
	_ = cmd.Run()
	fmt.Println("Removed from crontab.")
}
