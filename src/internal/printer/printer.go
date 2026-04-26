// Package printer prints out the tree with the icons.
package printer

import (
	"fmt"
	"strings"

	"gotree/src/config/icon"
	"gotree/src/internal/model"
)

type Config struct {
	IconMode      icon.Mode
	EnableColor   bool
	ShowSize      bool
	ShowSizeTotal bool
}

type Summary struct {
	Directories int   `json:"directories" yaml:"directories"`
	Files       int   `json:"files" yaml:"files"`
	TotalSize   int64 `json:"total_size_bytes" yaml:"total_size_bytes"`
}

func Print(root *model.Node, cfg Config) {
	glyph, color := icon.ResolveIcon(root.Name, root.IsDir, cfg.IconMode)
	sizePrefix := ""
	if cfg.ShowSize {
		sizePrefix = fmt.Sprintf("%8s ", HumanizeBytes(root.SizeBytes))
	}
	fmt.Print(sizePrefix)
	if glyph != "" {
		fmt.Print(icon.Colorize(color, glyph, cfg.EnableColor) + " ")
	}
	fmt.Println(root.Name)
	for i, child := range root.Children {
		printNode(child, "", i == len(root.Children)-1, cfg)
	}
}

func BuildSummary(root *model.Node) Summary {
	dirs, files := root.Count()
	return Summary{
		Directories: dirs,
		Files:       files,
		TotalSize:   root.SizeBytes,
	}
}

func PrintSummary(root *model.Node, cfg Config) {
	s := BuildSummary(root)
	if cfg.ShowSizeTotal {
		fmt.Printf("\n%d directories, %d files, %s total\n", s.Directories, s.Files, HumanizeBytes(s.TotalSize))
		return
	}
	fmt.Printf("\n%d directories, %d files\n", s.Directories, s.Files)
}

func printNode(n *model.Node, prefix string, isLast bool, cfg Config) {
	connector := "├── "
	nextPrefix := "│   "

	if isLast {
		connector = "└── "
		nextPrefix = "    "
	}

	if n.IsSummary {
		fmt.Println(prefix + connector + "(" + n.SummaryMsg + ")")
		return
	}

	glyph, color := icon.ResolveIcon(n.Name, n.IsDir, cfg.IconMode)
	fmt.Print(prefix + connector)
	if cfg.ShowSize {
		fmt.Print(fmt.Sprintf("%8s ", HumanizeBytes(n.SizeBytes)))
	}
	if glyph != "" {
		fmt.Print(icon.Colorize(color, glyph, cfg.EnableColor) + " ")
	}
	fmt.Println(n.Name)

	for i, child := range n.Children {
		last := i == len(n.Children)-1
		printNode(child, prefix+nextPrefix, last, cfg)
	}
}

func HumanizeBytes(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%dB", size)
	}
	units := []string{"KB", "MB", "GB", "TB"}
	val := float64(size)
	for _, unit := range units {
		val /= 1024
		if val < 1024 {
			text := fmt.Sprintf("%.1f%s", val, unit)
			return strings.ReplaceAll(text, ".0", "")
		}
	}
	text := fmt.Sprintf("%.1fPB", val/1024)
	return strings.ReplaceAll(text, ".0", "")
}
