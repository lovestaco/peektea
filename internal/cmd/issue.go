package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/lovestaco/peektea/internal/config"
)

const issuesURL = "https://github.com/lovestaco/peektea/issues"

// RunIssue opens the project's GitHub issues page in the default browser.
func RunIssue() {
	opener, args := browserOpener()

	if _, err := exec.LookPath(opener); err != nil {
		fmt.Println("Open this URL to view or file an issue:")
		fmt.Println(issuesURL)
		return
	}

	c := exec.Command(opener, append(args, issuesURL)...)
	if err := c.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to open browser: %v\n", err)
		fmt.Println(issuesURL)
		os.Exit(1)
	}
}

// browserOpener returns the command and leading args used to open a URL
// in the default browser for the current platform.
func browserOpener() (string, []string) {
	if runtime.GOOS == "darwin" {
		return "open", nil
	}
	if wslOpener := config.WSLOpener(); wslOpener != "" {
		return wslOpener, nil
	}
	return "xdg-open", nil
}
