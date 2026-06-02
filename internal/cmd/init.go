package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	pickerCursor = lipgloss.NewStyle().Foreground(lipgloss.Color("#7DAD5C")).Bold(true)
	pickerMuted  = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type category struct {
	label    string
	comment  string
	programs []string
	fallback string
	keys     []string
}

var setupCategories = []category{
	{
		label:   "Text editor",
		comment: "text / code",
		programs: []string{
			"nvim", "vim", "vi", "nano", "micro", "hx", "emacs", "code", "gedit", "kate", "mousepad",
		},
		fallback: "vim",
		keys: []string{
			"_default_config",
			"_txt_config", "_md_config",
			"_go_config", "_py_config", "_sh_config", "_js_config", "_ts_config",
			"_rs_config", "_c_config", "_h_config", "_cpp_config",
			"_json_config", "_yaml_config", "_yml_config", "_toml_config",
			"_html_config", "_css_config",
		},
	},
	{
		label:    "File manager",
		comment:  "directories & archives",
		programs: []string{"nautilus", "thunar", "nemo", "dolphin", "pcmanfm"},
		fallback: "xdg-open",
		keys:     []string{"_dir_config", "_zip_config", "_tar_gz_config", "_tar_bz2_config", "_gz_config", "_xz_config"},
	},
	{
		label:    "Image viewer",
		comment:  "images",
		programs: []string{"feh", "eog", "sxiv", "imv", "viewnior", "ristretto", "gwenview", "display"},
		fallback: "xdg-open",
		keys:     []string{"_png_config", "_jpg_config", "_jpeg_config", "_gif_config", "_webp_config", "_bmp_config", "_svg_config"},
	},
	{
		label:    "PDF viewer",
		comment:  "documents",
		programs: []string{"evince", "okular", "zathura", "mupdf", "atril", "xpdf"},
		fallback: "xdg-open",
		keys:     []string{"_pdf_config"},
	},
}

func RunInit() {
	home, _ := os.UserHomeDir()
	dest := filepath.Join(home, ".peektea.toml")

	if _, err := os.Stat(dest); err == nil {
		fmt.Printf("%s already exists. Overwrite? [y/N]: ", dest)
		var ans string
		fmt.Scanln(&ans)
		if strings.ToLower(strings.TrimSpace(ans)) != "y" {
			fmt.Println("aborted.")
			os.Exit(0)
		}
	}

	selections := map[string]string{}

	fmt.Println("peektea init")
	fmt.Println()
	fmt.Println("I peeked into your installed software — here's what I found. Pick one.")
	fmt.Println()

	for _, cat := range setupCategories {
		found := installedFrom(cat.programs)
		fmt.Printf("── %s\n", cat.label)

		var chosen string
		switch len(found) {
		case 0:
			chosen = cat.fallback
			if chosen != "" {
				fmt.Printf("   none detected — using %s as fallback\n\n", chosen)
			} else {
				fmt.Println("   none detected — skipping")
				fmt.Println()
				continue
			}
		case 1:
			chosen = found[0]
			fmt.Printf("   only %s found — selected automatically\n\n", chosen)
		default:
			idx := pickOne(found)
			chosen = found[idx]
			fmt.Printf("   selected: %s\n\n", chosen)
		}

		for _, key := range cat.keys {
			selections[key] = chosen
		}
	}

	if err := writeToml(dest, selections); err != nil {
		fmt.Fprintf(os.Stderr, "error writing config: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("created %s\n", dest)

	fmt.Println()
	fmt.Println("── Image previews")
	if _, err := exec.LookPath("chafa"); err != nil {
		fmt.Println("   chafa not found — install it for inline image previews.")
		fmt.Println("   Ubuntu/Debian : sudo apt install chafa")
		fmt.Println("   Arch          : sudo pacman -S chafa")
		fmt.Println("   macOS         : brew install chafa")
	} else {
		fmt.Println("   chafa found — image previews are ready.")
	}
}

type pickerModel struct {
	options []string
	cursor  int
}

func (m pickerModel) Init() tea.Cmd { return nil }

func (m pickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "enter", " ":
			return m, tea.Quit
		case "ctrl+c":
			os.Exit(0)
		default:
			if n, err := strconv.Atoi(msg.String()); err == nil && n >= 1 && n <= len(m.options) {
				m.cursor = n - 1
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m pickerModel) View() string {
	var sb strings.Builder
	for i, opt := range m.options {
		if i == m.cursor {
			sb.WriteString(pickerCursor.Render(fmt.Sprintf("  ▶ %d) %s", i+1, opt)) + "\n")
		} else {
			sb.WriteString(pickerMuted.Render(fmt.Sprintf("    %d) %s", i+1, opt)) + "\n")
		}
	}
	sb.WriteString("\n" + pickerMuted.Render("  ↑/↓ or number  enter to confirm"))
	return sb.String()
}

func pickOne(options []string) int {
	p := tea.NewProgram(pickerModel{options: options})
	result, _ := p.Run()
	if m, ok := result.(pickerModel); ok {
		return m.cursor
	}
	return 0
}

func installedFrom(programs []string) []string {
	var found []string
	for _, p := range programs {
		if _, err := exec.LookPath(p); err == nil {
			found = append(found, p)
		}
	}
	return found
}

func writeToml(dest string, selections map[string]string) error {
	var sb strings.Builder

	sb.WriteString("# ~/.peektea.toml — generated by peektea init\n")
	sb.WriteString("# key format: file.tar.gz → _tar_gz_config  |  directory → _dir_config\n\n")
	sb.WriteString("terminal_programs = [\"vim\", \"nvim\", \"vi\", \"nano\", \"emacs\", \"micro\", \"hx\", \"helix\"]\n")

	for _, cat := range setupCategories {
		hasAny := false
		for _, key := range cat.keys {
			if _, ok := selections[key]; ok {
				hasAny = true
				break
			}
		}
		if !hasAny {
			continue
		}
		fmt.Fprintf(&sb, "\n# %s\n", cat.comment)
		for _, key := range cat.keys {
			if prog, ok := selections[key]; ok {
				fmt.Fprintf(&sb, "%-20s = %q\n", key, prog)
			}
		}
	}

	return os.WriteFile(dest, []byte(sb.String()), 0644)
}
