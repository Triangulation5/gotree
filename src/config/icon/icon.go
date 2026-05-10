// Package icon renders all the icons for this project
package icon

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Mode int

const (
	Nerd Mode = iota
	Unicode
	ASCII
	None
)

func ParseMode(s string) Mode {
	switch strings.ToLower(s) {
	case "unicode":
		return Unicode
	case "ascii":
		return ASCII
	case "none":
		return None
	default:
		return Nerd
	}
}

type Style struct {
	Nerd    string
	Unicode string
	ASCII   string
	Color   string
}

var (
	DefaultFolder = Style{Nerd: "\uf07b", Unicode: "📁", ASCII: "[D]", Color: "#80d4ff"}
	DefaultFile   = Style{Nerd: "\uf15b", Unicode: "📄", ASCII: "[F]", Color: "#ffffff"}
	HiddenFolder  = Style{Nerd: "\uf023", Unicode: "🔒", ASCII: "[H]", Color: "#75715e"}
)

var iconMap = map[string]Style{
	// Programming Languages
	"go":   {Nerd: "\ue627", Unicode: "GO", ASCII: "go", Color: "#6ed8e5"},
	"js":   {Nerd: "\ue781", Unicode: "JS", ASCII: "js", Color: "#f39c12"},
	"ts":   {Nerd: "\ue628", Unicode: "TS", ASCII: "ts", Color: "#2980b9"},
	"py":   {Nerd: "\ue606", Unicode: "PY", ASCII: "py", Color: "#3498db"},
	"rb":   {Nerd: "\ue21e", Unicode: "RB", ASCII: "rb", Color: "#9b59b6"},
	"rs":   {Nerd: "\ue7a8", Unicode: "RS", ASCII: "rs", Color: "#f39c12"},
	"java": {Nerd: "\ue738", Unicode: "☕", ASCII: "jv", Color: "#e67e22"},
	"c":    {Nerd: "\ue649", Unicode: "C", ASCII: "c", Color: "#0188d2"},
	"cpp":  {Nerd: "\ue646", Unicode: "C+", ASCII: "cp", Color: "#0188d2"},
	"html": {Nerd: "\uf13b", Unicode: "HT", ASCII: "ht", Color: "#e67e22"},
	"css":  {Nerd: "\uf13c", Unicode: "CS", ASCII: "cs", Color: "#2d53e5"},

	// Data & Config
	"json": {Nerd: "\ue60b", Unicode: "{ }", ASCII: "js", Color: "#f1c40f"},
	"yaml": {Nerd: "\ue601", Unicode: "YM", ASCII: "ym", Color: "#f39c12"},
	"yml":  {Nerd: "\ue601", Unicode: "YM", ASCII: "ym", Color: "#f39c12"},
	"toml": {Nerd: "\ue6a2", Unicode: "TM", ASCII: "tm", Color: "#f39c12"},
	"xml":  {Nerd: "\ue796", Unicode: "XM", ASCII: "xm", Color: "#3498db"},
	"md":   {Nerd: "\uf48a", Unicode: "📝", ASCII: "md", Color: "#7f8c8d"},
	"sql":  {Nerd: "\ue706", Unicode: "DB", ASCII: "sq", Color: "#FF8400"},

	// Special Files
	"gitignore":  {Nerd: "\ue702", Unicode: "GIT", ASCII: "gi", Color: "#e67e22"},
	"dockerfile": {Nerd: "\uf308", Unicode: "🐳", ASCII: "dk", Color: "#099cec"},
	"license":    {Nerd: "\uf02d", Unicode: "📜", ASCII: "lc", Color: "#d35400"},
}

var dirMap = map[string]Style{
	".git":         {Nerd: "\ue5fb", Unicode: "GIT", ASCII: "[G]", Color: "#f14e32"},
	"node_modules": {Nerd: "\ue5fa", Unicode: "NPM", ASCII: "[N]", Color: "#cb3837"},
	".vscode":      {Nerd: "\ue70c", Unicode: "VSC", ASCII: "[V]", Color: "#007acc"},
	"src":          {Nerd: "\uf121", Unicode: "SRC", ASCII: "[S]", Color: "#ffb86c"},
	"cmd":          {Nerd: "\uf489", Unicode: "CMD", ASCII: "[C]", Color: "#2ecc71"},
	"internal":     {Nerd: "\uf023", Unicode: "INT", ASCII: "[I]", Color: "#75715e"},
}

func ResolveIcon(name string, isDir bool, mode Mode) (string, string) {
	var style Style
	found := false

	if isDir {
		style, found = dirMap[name]
		if !found {
			if len(name) > 0 && name[0] == '.' {
				style = HiddenFolder
			} else {
				style = DefaultFolder
			}
		}
	} else {
		lowerName := strings.ToLower(name)
		// Check full name first (for things like LICENSE or .gitignore)
		style, found = iconMap[strings.TrimPrefix(lowerName, ".")]
		if !found {
			ext := strings.TrimPrefix(filepath.Ext(lowerName), ".")
			style, found = iconMap[ext]
		}
		if !found {
			style = DefaultFile
		}
	}

	var glyph string
	switch mode {
	case Unicode:
		glyph = style.Unicode
	case ASCII:
		glyph = style.ASCII
	case None:
		return "", ""
	default:
		glyph = style.Nerd
	}

	return glyph, style.Color
}

func Colorize(color string, text string, enabled bool) string {
	if !enabled {
		return text
	}
	if color == "NONE" || color == "" {
		return text
	}

	hexStr := strings.TrimPrefix(color, "#")
	if len(hexStr) != 6 {
		return text
	}

	var r, g, b uint8
	_, err := fmt.Sscanf(hexStr, "%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return text
	}

	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", r, g, b, text)
}
