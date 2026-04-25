package scheduler

import "runtime"

func Schedule(processName string) {
	switch runtime.GOOS {
	case "darwin", "linux":
		startCron(processName)
	case "windows":
		startWindowsScheduler(processName)
	}
}

func Unschedule(processName string) {
	switch runtime.GOOS {
	case "darwin", "linux":
		stopCron(processName)
	case "windows":
		stopWindowsScheduler()
	}
}
