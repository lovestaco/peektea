package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type Config struct {
	data         map[string]string
	terminalApps map[string]bool
}

var defaultTerminalApps = []string{
	"vim", "nvim", "vi", "nano", "emacs", "micro", "hx", "helix", "joe", "mcedit",
}

func Load() Config {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".peektea.toml")

	raw := map[string]interface{}{}
	if _, err := os.Stat(path); err == nil {
		toml.DecodeFile(path, &raw) //nolint
	}

	data := map[string]string{}
	terminalApps := map[string]bool{}

	for _, p := range defaultTerminalApps {
		terminalApps[p] = true
	}

	for k, v := range raw {
		if k == "terminal_programs" {
			if arr, ok := v.([]interface{}); ok {
				for _, item := range arr {
					if s, ok := item.(string); ok {
						terminalApps[s] = true
					}
				}
			}
		} else if s, ok := v.(string); ok {
			data[k] = s
		}
	}

	return Config{data: data, terminalApps: terminalApps}
}

// ProgramFor returns the configured opener for an entry.
// Key derivation: take everything after the first non-leading dot,
// replace remaining dots with underscores, wrap with _ and _config.
//
//	file.md        → _md_config
//	archive.tar.gz → _tar_gz_config
//	hello.xd.dd    → _xd_dd_config
//	directory      → _dir_config
//	no-extension   → _default_config
func (c Config) ProgramFor(entry os.DirEntry) string {
	key := keyFor(entry)
	if prog := c.data[key]; prog != "" {
		return prog
	}
	if prog := c.data["_default_config"]; prog != "" {
		return prog
	}
	return "vim"
}

func (c Config) IsTerminalApp(prog string) bool {
	return c.terminalApps[prog]
}

func keyFor(entry os.DirEntry) string {
	if entry.IsDir() {
		return "_dir_config"
	}
	name := entry.Name()
	start := 0
	if strings.HasPrefix(name, ".") {
		start = 1
	}
	idx := strings.Index(name[start:], ".")
	if idx < 0 {
		return "_default_config"
	}
	ext := name[start+idx+1:]
	return "_" + strings.ReplaceAll(ext, ".", "_") + "_config"
}
