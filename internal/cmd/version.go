package cmd

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

func RunVersion(v string) {
	fmt.Printf("peektea %s (%s/%s)\n", resolveVersion(v), runtime.GOOS, runtime.GOARCH)
}

// resolveVersion returns the ldflags-injected version when available,
// falling back to the module version embedded by go install.
func resolveVersion(ldflagsVersion string) string {
	if ldflagsVersion != "dev" {
		return ldflagsVersion
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return ldflagsVersion
}
