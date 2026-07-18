//go:build darwin || linux

package scheduler

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Schedule(processName string) {
	out, _ := exec.Command("crontab", "-l").Output()
	if strings.Contains(string(out), processName) {
		fmt.Println("Already in crontab, skipping.")
		return
	}
	entry := strings.TrimRight(string(out), "\n") + "\n@monthly " + processName + "\n"
	cmd := exec.Command("crontab", "-")
	cmd.Stdin = strings.NewReader(entry)
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "crontab error:", err)
		return
	}
	fmt.Println("Scheduled: crontab @monthly")
}

func Unschedule(processName string) {
	out, err := exec.Command("crontab", "-l").Output()
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
	cmd := exec.Command("crontab", "-")
	cmd.Stdin = strings.NewReader(strings.Join(lines, "\n") + "\n")
	cmd.Run()
	fmt.Println("Removed from crontab.")
}
