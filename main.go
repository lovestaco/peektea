package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lovestaco/peektea/internal/cmd"
	"github.com/lovestaco/peektea/internal/config"
)

var (
	cursorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	dirStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	fileStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	pathStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Italic(true)
	selectedBg   = lipgloss.NewStyle().Background(lipgloss.Color("236"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	taglineStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7DAD5C")).Bold(true)
)

type openResultMsg struct{ err error }

type model struct {
	dir     string
	entries []os.DirEntry
	cursor  int
	err     error
	status  string
	config  config.Config
}

func newModel(dir string) model {
	m := model{dir: dir, config: config.Load()}
	m.entries, m.err = os.ReadDir(dir)
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case openResultMsg:
		if msg.err != nil {
			m.status = fmt.Sprintf("open failed: %v", msg.err)
		}
		return m, nil

	case tea.KeyMsg:
		m.status = ""
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

		case "o":
			if len(m.entries) > 0 {
				entry := m.entries[m.cursor]
				prog := m.config.ProgramFor(entry)
				path := filepath.Join(m.dir, entry.Name())
				if m.config.IsTerminalApp(prog) {
					c := exec.Command(prog, path)
					return m, tea.ExecProcess(c, func(err error) tea.Msg {
						return openResultMsg{err: err}
					})
				}
				// GUI app — launch in background, don't block the TUI
				c := exec.Command(prog, path)
				if err := c.Start(); err != nil {
					m.status = fmt.Sprintf("open failed: %v", err)
				} else {
					go c.Wait() //nolint — reap the child so it doesn't become a zombie
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
	sb.WriteString(taglineStyle.Render("peek-a-boo, filesystem.") + "\n\n")
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

	sb.WriteString("\n" + pathStyle.Render("↑/↓ navigate  →/enter go inside  o open  ←/backspace up  q quit"))
	if m.status != "" {
		sb.WriteString("\n" + errorStyle.Render(m.status))
	}
	return sb.String()
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "init":
			cmd.RunInit()
		case "uninstall":
			cmd.RunUninstall()
		case "-h", "--help", "help":
			cmd.RunHelp()
		default:
			fmt.Fprintf(os.Stderr, "unknown command: %s\n\nrun 'peektea -h' for help\n", os.Args[1])
			os.Exit(1)
		}
		return
	}

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
