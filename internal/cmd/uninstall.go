package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func RunUninstall() {
	bin, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not find binary path: %v\n", err)
		os.Exit(1)
	}

	home, _ := os.UserHomeDir()
	config := filepath.Join(home, ".peektea.toml")

	fmt.Println("Uninstalling peektea...")
	fmt.Printf("  binary: %s\n", bin)
	fmt.Println()

	fmt.Print("Are you sure? [y/N]: ")
	var confirm string
	fmt.Scanln(&confirm)
	if strings.ToLower(strings.TrimSpace(confirm)) != "y" {
		fmt.Println("aborted.")
		os.Exit(0)
	}

	deleteConfig := false
	if _, err := os.Stat(config); err == nil {
		fmt.Printf("\nDelete your config (%s)? [y/N]: ", config)
		var ans string
		fmt.Scanln(&ans)
		deleteConfig = strings.ToLower(strings.TrimSpace(ans)) == "y"
	}

	fmt.Println()

	if err := os.Remove(bin); err != nil {
		fmt.Fprintf(os.Stderr, "  failed to remove binary: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  removed %s\n", bin)

	if deleteConfig {
		if err := os.Remove(config); err != nil {
			fmt.Fprintf(os.Stderr, "  failed to remove config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("  removed %s\n", config)
	}

	fmt.Println()
	fmt.Println("Done. Goodbye!")
}
