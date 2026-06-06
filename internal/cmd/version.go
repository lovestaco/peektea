package cmd

import (
	"fmt"
	"runtime"
)

func RunVersion(v string) {
	fmt.Printf("peektea %s (%s/%s)\n", v, runtime.GOOS, runtime.GOARCH)
}
