package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	cursorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	dirStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	fileStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	pathStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Italic(true)
	selectedBg     = lipgloss.NewStyle().Background(lipgloss.Color("236"))
)

type model struct {
	dir     string
	entries []os.DirEntry
	cursor  int
	err     error
}

func newModel(dir string) model {
	m := model{dir: dir}
	m.entries, m.err = os.ReadDir(dir)
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.entries)-1 {
				m.cursor++
			}

		case "right", "l", "enter":
			if len(m.entries) > 0 && m.entries[m.cursor].IsDir() {
				next := filepath.Join(m.dir, m.entries[m.cursor].Name())
				entries, err := os.ReadDir(next)
				if err == nil {
					m.dir = next
					m.entries = entries
					m.cursor = 0
				}
			}

		case "left", "h", "backspace":
			parent := filepath.Dir(m.dir)
			if parent != m.dir {
				entries, err := os.ReadDir(parent)
				if err == nil {
					// restore cursor to the dir we came from
					oldName := filepath.Base(m.dir)
					m.dir = parent
					m.entries = entries
					m.cursor = 0
					for i, e := range entries {
						if e.Name() == oldName {
							m.cursor = i
							break
						}
					}
				}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("error: %v\n", m.err)
	}

	var sb strings.Builder
	sb.WriteString(pathStyle.Render(m.dir) + "\n\n")

	if len(m.entries) == 0 {
		sb.WriteString(fileStyle.Render("  (empty)") + "\n")
	}

	for i, e := range m.entries {
		cursor := "  "
		if i == m.cursor {
			cursor = cursorStyle.Render("▶ ")
		}

		name := e.Name()
		var nameStyled string
		if e.IsDir() {
			nameStyled = dirStyle.Render(name + "/")
		} else {
			nameStyled = fileStyle.Render(name)
		}

		line := cursor + nameStyled
		if i == m.cursor {
			line = selectedBg.Render(line)
		}

		sb.WriteString(line + "\n")
	}

	sb.WriteString("\n" + pathStyle.Render("↑/↓ navigate  →/enter open  ←/backspace up  q quit"))
	return sb.String()
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	p := tea.NewProgram(newModel(dir), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
