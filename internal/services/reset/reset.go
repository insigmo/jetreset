package reset

import (
	"fmt"
	"os"
	"runtime"
)

func Reset(home string, products []string) {
	switch runtime.GOOS {
	case "darwin":
		resetMacos(home, products)
	case "linux":
		resetLinux(home, products)
	case "windows":
		resetWindows(home, products)
	default:
		fmt.Fprintln(os.Stderr, "Unsupported OS")
		os.Exit(1)
	}
}
