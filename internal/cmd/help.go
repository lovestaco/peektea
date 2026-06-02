package cmd

import "fmt"

func RunHelp() {
	fmt.Println("peektea — a terminal file browser")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  peektea              start the file browser")
	fmt.Println("  peektea init         interactive setup of ~/.peektea.toml")
	fmt.Println("  peektea uninstall    remove peektea and optionally its config")
	fmt.Println("  peektea -h           show this help")
	fmt.Println()
	fmt.Println("Keys:")
	fmt.Println("  ↑/↓  k/j        navigate")
	fmt.Println("  →    l/enter    go inside directory")
	fmt.Println("  ←    h/backspace  go to parent")
	fmt.Println("  o               open with configured program")
	fmt.Println("  p               toggle preview panel (text + images via chafa)")
	fmt.Println("  q    ctrl+c     quit")
}
