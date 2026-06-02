package main

import (
	"bufio"
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
	cursorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	dirStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	fileStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	pathStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Italic(true)
	selectedBg      = lipgloss.NewStyle().Background(lipgloss.Color("236"))
	errorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	taglineStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#7DAD5C")).Bold(true)
	previewHdrStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	panelBorder     = lipgloss.NewStyle().
				BorderLeft(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("241")).
				PaddingLeft(1)
)

type openResultMsg struct{ err error }
type previewMsg struct{ content string }

type model struct {
	dir            string
	entries        []os.DirEntry
	cursor         int
	err            error
	status         string
	config         config.Config
	width          int
	height         int
	showPreview    bool
	previewContent string
	previewLoading bool
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
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case previewMsg:
		m.previewContent = msg.content
		m.previewLoading = false
		return m, nil

	case openResultMsg:
		if msg.err != nil {
			m.status = fmt.Sprintf("open failed: %v", msg.err)
		}
		return m, nil

	case tea.KeyMsg:
		m.status = ""
		var needPreview bool
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				needPreview = true
			}

		case "down", "j":
			if m.cursor < len(m.entries)-1 {
				m.cursor++
				needPreview = true
			}

		case "right", "l", "enter":
			if len(m.entries) > 0 && m.entries[m.cursor].IsDir() {
				next := filepath.Join(m.dir, m.entries[m.cursor].Name())
				entries, err := os.ReadDir(next)
				if err == nil {
					m.dir = next
					m.entries = entries
					m.cursor = 0
					needPreview = true
				}
			}

		case "left", "h", "backspace":
			parent := filepath.Dir(m.dir)
			if parent != m.dir {
				entries, err := os.ReadDir(parent)
				if err == nil {
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
					needPreview = true
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
				c := exec.Command(prog, path)
				if err := c.Start(); err != nil {
					m.status = fmt.Sprintf("open failed: %v", err)
				} else {
					go c.Wait() //nolint — reap child to avoid zombie
				}
			}

		case "p":
			m.showPreview = !m.showPreview
			if m.showPreview && len(m.entries) > 0 {
				m.previewLoading = true
				m.previewContent = ""
				return m, m.previewCmd()
			}
			return m, nil
		}

		if m.showPreview && needPreview && len(m.entries) > 0 {
			m.previewLoading = true
			m.previewContent = ""
			return m, m.previewCmd()
		}
	}
	return m, nil
}

func (m model) leftWidth() int {
	const minWidth = 50
	if m.width == 0 {
		return minWidth
	}
	if m.width < 60 {
		return m.width / 2
	}
	w := minWidth
	for _, e := range m.entries {
		// 2 for cursor prefix, +1 for trailing slash on dirs
		nameW := 2 + len([]rune(e.Name()))
		if e.IsDir() {
			nameW++
		}
		if nameW > w {
			w = nameW
		}
	}
	// Always leave at least 30 chars for the preview panel.
	if cap := m.width - 32; w > cap {
		w = cap
	}
	return w
}

func (m model) previewCmd() tea.Cmd {
	if len(m.entries) == 0 {
		return nil
	}
	entry := m.entries[m.cursor]
	path := filepath.Join(m.dir, entry.Name())
	lw := m.leftWidth()
	rw := m.width - lw - 2 // 2 = border + padding
	if rw < 10 {
		rw = 40
	}
	ph := m.height - 6
	if ph < 5 {
		ph = 20
	}
	return loadPreview(path, entry, rw, ph)
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("error: %v\n", m.err)
	}
	fileList := m.renderFileList()
	if !m.showPreview || m.width == 0 {
		return fileList
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, fileList, m.renderPreview())
}

func (m model) renderFileList() string {
	lw := m.leftWidth()

	// Build the scrollable content area first so we can measure its line count.
	var top strings.Builder
	top.WriteString(taglineStyle.Render("peek-a-boo, filesystem.") + "\n\n")
	top.WriteString(pathStyle.Render(m.dir) + "\n\n")

	if len(m.entries) == 0 {
		top.WriteString(fileStyle.Render("  (empty)") + "\n")
	}
	for i, e := range m.entries {
		cursor := "  "
		if i == m.cursor {
			cursor = cursorStyle.Render("▶ ")
		}
		name := e.Name()
		if m.showPreview {
			maxLen := lw - 4 // cursor(2) + slash(1) + margin(1)
			runes := []rune(name)
			if len(runes) > maxLen {
				name = string(runes[:maxLen-1]) + "…"
			}
		}
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
		top.WriteString(line + "\n")
	}

	var hint string
	if m.showPreview {
		hint = "↑/↓  enter  o open  p close  q quit"
	} else {
		hint = "↑/↓ navigate  →/enter go in  o open  p preview  ←/h up  q quit"
	}

	var sb strings.Builder
	sb.WriteString(top.String())

	// Pad so the hint line sits at the very bottom of the panel.
	// top already ends each line with \n, so strings.Count gives lines consumed.
	// We need 1 more \n (blank gap) before the hint, so reserve 2 rows for that.
	if m.height > 0 {
		linesUsed := strings.Count(top.String(), "\n")
		padding := m.height - linesUsed - 2
		for i := 0; i < padding; i++ {
			sb.WriteByte('\n')
		}
	}

	sb.WriteString("\n" + pathStyle.Render(hint))
	if m.status != "" {
		sb.WriteString("\n" + errorStyle.Render(m.status))
	}

	if m.showPreview {
		return lipgloss.NewStyle().Width(lw).Render(sb.String())
	}
	return sb.String()
}

func (m model) renderPreview() string {
	var header string
	if len(m.entries) > 0 {
		header = previewHdrStyle.Render(m.entries[m.cursor].Name()) + "\n\n"
	}

	var body string
	switch {
	case m.previewLoading:
		body = pathStyle.Render("loading…")
	case m.previewContent == "":
		body = pathStyle.Render("(no preview)")
	default:
		body = m.previewContent
	}

	return panelBorder.Render(header + body)
}

func loadPreview(path string, entry os.DirEntry, width, height int) tea.Cmd {
	return func() tea.Msg {
		if entry.IsDir() {
			return previewMsg{content: previewDir(path, width, height)}
		}
		if isImageExt(entry.Name()) {
			return previewMsg{content: previewImage(path, width, height)}
		}
		return previewMsg{content: previewText(path, width, height)}
	}
}

func previewDir(path string, width, height int) string {
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	if len(entries) == 0 {
		return pathStyle.Render("(empty directory)")
	}
	var sb strings.Builder
	limit := height
	for i, e := range entries {
		if i >= limit {
			sb.WriteString(pathStyle.Render(fmt.Sprintf("  … %d more", len(entries)-i)) + "\n")
			break
		}
		if e.IsDir() {
			sb.WriteString(dirStyle.Render("  "+e.Name()+"/") + "\n")
		} else {
			sb.WriteString(fileStyle.Render("  "+e.Name()) + "\n")
		}
	}
	return sb.String()
}

func previewImage(path string, width, height int) string {
	if _, err := exec.LookPath("chafa"); err != nil {
		return pathStyle.Render("[image — install chafa for inline preview]")
	}
	out, err := exec.Command("chafa",
		"--size", fmt.Sprintf("%dx%d", width, height),
		path,
	).Output()
	if err != nil {
		return pathStyle.Render(fmt.Sprintf("[image preview failed: %v]", err))
	}
	return strings.TrimRight(string(out), "\n")
}

func previewText(path string, width, height int) string {
	if isBinary(path) {
		return pathStyle.Render("[binary file]")
	}
	f, err := os.Open(path)
	if err != nil {
		return pathStyle.Render(fmt.Sprintf("error: %v", err))
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var lines []string
	maxW := width - 1
	if maxW < 1 {
		maxW = 40
	}
	for scanner.Scan() && len(lines) < height {
		line := scanner.Text()
		line = strings.ReplaceAll(line, "\t", "    ")
		runes := []rune(line)
		if len(runes) > maxW {
			line = string(runes[:maxW-1]) + "…"
		}
		lines = append(lines, line)
	}
	if len(lines) == 0 {
		return pathStyle.Render("(empty file)")
	}
	return strings.Join(lines, "\n")
}

func isBinary(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	buf := make([]byte, 512)
	n, _ := f.Read(buf)
	for _, b := range buf[:n] {
		if b == 0 {
			return true
		}
	}
	return false
}

var imageExts = map[string]bool{
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true,
	".webp": true, ".bmp": true, ".svg": true, ".tiff": true, ".tif": true,
}

func isImageExt(name string) bool {
	return imageExts[strings.ToLower(filepath.Ext(name))]
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
