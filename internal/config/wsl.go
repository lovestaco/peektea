package config

import (
	"os"
	"os/exec"
	"strings"
)

// IsWSL reports whether we're running inside Windows Subsystem for Linux.
func IsWSL() bool {
	if os.Getenv("WSL_DISTRO_NAME") != "" {
		return true
	}
	data, err := os.ReadFile("/proc/version")
	return err == nil && strings.Contains(strings.ToLower(string(data)), "microsoft")
}

// WSLOpener returns the best program for opening files with Windows apps:
// wslview (handles path conversion itself) if installed, otherwise explorer.exe.
// Returns "" when neither is reachable.
func WSLOpener() string {
	if _, err := exec.LookPath("wslview"); err == nil {
		return "wslview"
	}
	if _, err := exec.LookPath("explorer.exe"); err == nil {
		return "explorer.exe"
	}
	if _, err := os.Stat("/mnt/c/Windows/explorer.exe"); err == nil {
		return "/mnt/c/Windows/explorer.exe"
	}
	return ""
}

// WindowsPath converts a Linux path to a Windows path (\\wsl.localhost\...)
// via wslpath. Windows programs like explorer.exe can't read /home/... paths.
// Returns the original path if conversion fails.
func WindowsPath(path string) string {
	out, err := exec.Command("wslpath", "-w", path).Output()
	if err != nil {
		return path
	}
	return strings.TrimSpace(string(out))
}
