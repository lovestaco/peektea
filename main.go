package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lovestaco/peektea/internal/cmd"
	"github.com/lovestaco/peektea/internal/config"
)

var version = "dev"

const (
	sortName = iota
	sortSize
	sortMod
)

var sortLabels = [3]string{"name", "size", "modified"}

var (
	cursorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#9AB65A")).Bold(true)
	dirStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	fileStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	pathStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Italic(true)
	selectedBg      = lipgloss.NewStyle().Background(lipgloss.Color("236"))
	errorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	taglineStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#7DAD5C")).Bold(true)
	previewHdrStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#9AB65A")).Bold(true)
	filterTagStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#9AB65A")).Bold(true)
	scrollTrackStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	scrollThumbStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#9AB65A"))
	panelBorder      = lipgloss.NewStyle().
				BorderLeft(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("241")).
				PaddingLeft(1)
)

type openResultMsg struct{ err error }
type previewMsg struct{ content string }

type model struct {
	dir            string
	allEntries     []os.DirEntry
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
	showHidden     bool
	filterInput    textinput.Model
	filtering      bool
	sortMode       int
	previewScroll  int
	listOffset     int
}

func newModel(dir string) model {
	ti := textinput.New()
	ti.Placeholder = "type to filter…"
	ti.CharLimit = 64
	ti.Prompt = "/ "
	ti.PromptStyle = filterTagStyle
	ti.TextStyle = fileStyle

	m := model{dir: dir, config: config.Load(), filterInput: ti}
	all, err := os.ReadDir(dir)
	m.allEntries = all
	m.err = err
	return m.withFilters()
}

// withFilters recomputes entries from allEntries applying the hidden and search filters.
func (m model) withFilters() model {
	q := strings.ToLower(m.filterInput.Value())
	var filtered []os.DirEntry
	for _, e := range m.allEntries {
		if !m.showHidden && strings.HasPrefix(e.Name(), ".") {
			continue
		}
		if q != "" && !strings.Contains(strings.ToLower(e.Name()), q) {
			continue
		}
		filtered = append(filtered, e)
	}
	sort.SliceStable(filtered, func(i, j int) bool {
		switch m.sortMode {
		case sortSize:
			ii, _ := filtered[i].Info()
			jj, _ := filtered[j].Info()
			if ii == nil || jj == nil {
				return false
			}
			return ii.Size() > jj.Size()
		case sortMod:
			ii, _ := filtered[i].Info()
			jj, _ := filtered[j].Info()
			if ii == nil || jj == nil {
				return false
			}
			return ii.ModTime().After(jj.ModTime())
		default:
			return strings.ToLower(filtered[i].Name()) < strings.ToLower(filtered[j].Name())
		}
	})
	m.entries = filtered
	if m.cursor >= len(m.entries) {
		if len(m.entries) == 0 {
			m.cursor = 0
		} else {
			m.cursor = len(m.entries) - 1
		}
	}
	return m.clampListOffset()
}

func (m model) visibleFileCount() int {
	if m.height == 0 {
		return 20
	}
	bottomRows := 2
	if m.filtering || m.filterInput.Value() != "" {
		bottomRows = 3
	}
	v := m.height - 4 - bottomRows // 4 header rows (tagline+blank+path+blank)
	if v < 1 {
		v = 1
	}
	return v
}

func (m model) clampListOffset() model {
	visible := m.visibleFileCount()
	if m.cursor < m.listOffset {
		m.listOffset = m.cursor
	}
	if m.cursor >= m.listOffset+visible {
		m.listOffset = m.cursor - visible + 1
	}
	if m.listOffset < 0 {
		m.listOffset = 0
	}
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
		m.previewScroll = 0
		return m, nil

	case openResultMsg:
		if msg.err != nil {
			m.status = fmt.Sprintf("open failed: %v", msg.err)
		}
		return m, nil

	case tea.KeyMsg:
		m.status = ""

		if m.filtering {
			switch msg.String() {
			case "esc":
				m.filtering = false
				m.filterInput.Blur()
				m.filterInput.SetValue("")
				m = m.withFilters()
			case "enter":
				m.filtering = false
				m.filterInput.Blur()
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
					m = m.clampListOffset()
				}
			case "down", "j":
				if m.cursor < len(m.entries)-1 {
					m.cursor++
					m = m.clampListOffset()
				}
			default:
				var tiCmd tea.Cmd
				m.filterInput, tiCmd = m.filterInput.Update(msg)
				m = m.withFilters()
				if m.showPreview && len(m.entries) > 0 {
					m.previewLoading = true
					m.previewContent = ""
					return m, tea.Batch(tiCmd, m.previewCmd())
				}
				return m, tiCmd
			}
			if m.showPreview && len(m.entries) > 0 {
				m.previewLoading = true
				m.previewContent = ""
				return m, m.previewCmd()
			}
			return m, nil
		}

		var needPreview bool
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				m = m.clampListOffset()
				needPreview = true
			}

		case "down", "j":
			if m.cursor < len(m.entries)-1 {
				m.cursor++
				m = m.clampListOffset()
				needPreview = true
			}

		case "right", "l", "enter":
			if len(m.entries) > 0 && m.entries[m.cursor].IsDir() {
				next := filepath.Join(m.dir, m.entries[m.cursor].Name())
				all, err := os.ReadDir(next)
				if err == nil {
					m.dir = next
					m.allEntries = all
					m.cursor = 0
					m.filterInput.SetValue("")
					m.filtering = false
					m = m.withFilters()
					needPreview = true
				}
			}

		case "left", "h", "backspace":
			parent := filepath.Dir(m.dir)
			if parent != m.dir {
				all, err := os.ReadDir(parent)
				if err == nil {
					oldName := filepath.Base(m.dir)
					m.dir = parent
					m.allEntries = all
					m.filterInput.SetValue("")
					m.filtering = false
					m = m.withFilters()
					m.cursor = 0
					for i, e := range m.entries {
						if e.Name() == oldName {
							m.cursor = i
							break
						}
					}
					needPreview = true
				}
			}

		case "H":
			home, err := os.UserHomeDir()
			if err == nil && m.dir != home {
				all, err := os.ReadDir(home)
				if err == nil {
					m.dir = home
					m.allEntries = all
					m.cursor = 0
					m.filterInput.SetValue("")
					m.filtering = false
					m = m.withFilters()
					needPreview = true
				}
			}

		case "o":
			if len(m.entries) > 0 {
				entry := m.entries[m.cursor]
				prog := m.config.ProgramFor(entry)
				path := filepath.Join(m.dir, entry.Name())
				if strings.HasSuffix(prog, ".exe") {
					// Windows programs (WSL interop) can't read /home/... paths.
					path = config.WindowsPath(path)
				}
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

		case "/":
			m.filtering = true
			m.filterInput.Focus()
			return m, textinput.Blink

		case "esc":
			if m.filterInput.Value() != "" {
				m.filterInput.SetValue("")
				m = m.withFilters()
				needPreview = true
			}

		case ".":
			m.showHidden = !m.showHidden
			m = m.withFilters()
			needPreview = true

		case "s":
			m.sortMode = (m.sortMode + 1) % 3
			m = m.withFilters()
			needPreview = true

		case "[":
			if m.showPreview {
				m.previewScroll -= m.height / 4
				if m.previewScroll < 0 {
					m.previewScroll = 0
				}
			}

		case "]":
			if m.showPreview {
				m.previewScroll += m.height / 4
			}
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
		nameW := 2 + len([]rune(e.Name()))
		if e.IsDir() {
			nameW++
		}
		if nameW > w {
			w = nameW
		}
	}
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
	rw := m.width - lw - 2
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

	var top strings.Builder
	top.WriteString(taglineStyle.Render("peek-a-boo, filesystem.") + "\n\n")
	top.WriteString(pathStyle.Render(m.dir) + "\n\n")

	visibleCount := m.visibleFileCount()
	showListBar := len(m.entries) > visibleCount

	// line width for padding when scrollbar is present
	lineWidth := lw
	if !m.showPreview && m.width > 0 {
		lineWidth = m.width
	}

	if len(m.entries) == 0 {
		label := "  (empty)"
		if m.filterInput.Value() != "" {
			label = "  no matches"
		}
		top.WriteString(fileStyle.Render(label) + "\n")
	}

	var barChars []string
	if showListBar {
		barChars = buildScrollbar(len(m.entries), visibleCount, m.listOffset)
	}

	end := m.listOffset + visibleCount
	if end > len(m.entries) {
		end = len(m.entries)
	}
	for i, e := range m.entries[m.listOffset:end] {
		absIdx := m.listOffset + i
		cursor := "  "
		if absIdx == m.cursor {
			cursor = cursorStyle.Render("▶ ")
		}
		name := e.Name()
		maxLen := lineWidth - 4
		if showListBar {
			maxLen -= 2
		}
		if m.showPreview || showListBar {
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
		if showListBar {
			targetW := lineWidth - 2
			if pad := targetW - lipgloss.Width(line); pad > 0 {
				line += strings.Repeat(" ", pad)
			}
		}
		if absIdx == m.cursor {
			line = selectedBg.Render(line)
		}
		if showListBar && i < len(barChars) {
			line += " " + barChars[i]
		}
		top.WriteString(line + "\n")
	}

	dotLabel := "show hidden"
	if m.showHidden {
		dotLabel = "hide hidden"
	}
	var hint string
	if m.showPreview {
		hint = fmt.Sprintf("enter:go in  o:open  /:filter  .:%-11s  [/]:scroll  p:close  s:sorted by %s  H:home  q:quit", dotLabel, sortLabels[m.sortMode])
	} else {
		hint = fmt.Sprintf("enter:go in  o:open  /:filter  .:%-11s  p:preview  s:sorted by %s  H:home  q:quit", dotLabel, sortLabels[m.sortMode])
	}

	var sb strings.Builder
	sb.WriteString(top.String())

	if m.height > 0 {
		linesUsed := strings.Count(top.String(), "\n")
		// Filter bar (when active or set) sits above the hint — takes one extra row.
		bottomRows := 2
		if m.filtering || m.filterInput.Value() != "" {
			bottomRows = 3
		}
		padding := m.height - linesUsed - bottomRows
		for i := 0; i < padding; i++ {
			sb.WriteByte('\n')
		}
	}

	if m.filtering {
		sb.WriteString("\n" + m.filterInput.View())
	} else if m.filterInput.Value() != "" {
		sb.WriteString("\n" + filterTagStyle.Render("/"+m.filterInput.Value()) +
			pathStyle.Render("  esc to clear"))
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
		ph := m.height - 6
		if ph < 5 {
			ph = 10
		}
		lines := strings.Split(m.previewContent, "\n")
		scroll := m.previewScroll
		if maxScroll := len(lines) - ph; scroll > maxScroll {
			scroll = maxScroll
		}
		if scroll < 0 {
			scroll = 0
		}
		visible := lines[scroll:]
		if len(visible) > ph {
			visible = visible[:ph]
		}
		// content width = total - border(1) - padding(1) - space(1) - scrollbar(1)
		rw := m.width - m.leftWidth() - 4
		if rw < 10 {
			rw = 30
		}
		barChars := buildScrollbar(len(lines), ph, scroll)
		var bodySb strings.Builder
		for i, line := range visible {
			pad := rw - lipgloss.Width(line)
			if pad > 0 {
				line = line + strings.Repeat(" ", pad)
			}
			barChar := " "
			if i < len(barChars) {
				barChar = barChars[i]
			}
			bodySb.WriteString(line + " " + barChar + "\n")
		}
		body = bodySb.String()
	}
	if m.height > 0 {
		return panelBorder.Height(m.height).Render(header + body)
	}
	return panelBorder.Render(header + body)
}

func buildScrollbar(total, visible, scroll int) []string {
	chars := make([]string, visible)
	track := scrollTrackStyle.Render("│")
	thumb := scrollThumbStyle.Render("┃")
	if total <= visible {
		for i := range chars {
			chars[i] = track
		}
		return chars
	}
	thumbSize := visible * visible / total
	if thumbSize < 1 {
		thumbSize = 1
	}
	maxScroll := total - visible
	if maxScroll < 1 {
		maxScroll = 1
	}
	thumbPos := scroll * (visible - thumbSize) / maxScroll
	for i := range chars {
		if i >= thumbPos && i < thumbPos+thumbSize {
			chars[i] = thumb
		} else {
			chars[i] = track
		}
	}
	return chars
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
	for i, e := range entries {
		if i >= 1000 {
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

func previewText(path string, width, _ int) string {
	if isBinary(path) {
		return pathStyle.Render("[binary file]")
	}
	const maxLines = 500
	if _, err := exec.LookPath("bat"); err == nil {
		out, err := exec.Command("bat",
			"--color=always",
			"--style=plain",
			"--line-range", fmt.Sprintf(":%d", maxLines),
			"--terminal-width", fmt.Sprintf("%d", width),
			path,
		).Output()
		if err == nil {
			return strings.TrimRight(string(out), "\n")
		}
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
	for scanner.Scan() && len(lines) < maxLines {
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
		case "version", "--version", "-v":
			cmd.RunVersion(version)
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
