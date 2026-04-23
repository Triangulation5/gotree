package printer

import (
	"fmt"

	"gotree/src/config/icon"
	"gotree/src/internal/model"
)

type Config struct {
	IconMode icon.Mode
}

func Print(root *model.Node, cfg Config) {
	glyph, color := icon.ResolveIcon(root.Name, root.IsDir, cfg.IconMode)
	if glyph != "" {
		fmt.Print(icon.Colorize(color, glyph) + " ")
	}
	fmt.Println(root.Name)
	for i, child := range root.Children {
		printNode(child, "", i == len(root.Children)-1, cfg)
	}
}

func PrintSummary(root *model.Node) {
	dirs, files := root.Count()
	fmt.Printf("\n%d directories, %d files\n", dirs, files)
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
	if glyph != "" {
		fmt.Print(icon.Colorize(color, glyph) + " ")
	}
	fmt.Println(n.Name)

	for i, child := range n.Children {
		last := i == len(n.Children)-1
		printNode(child, prefix+nextPrefix, last, cfg)
	}
}
