package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"gotree/src/config/icon"
	"gotree/src/internal/printer"
	"gotree/src/internal/scanner"
)

const version = "2.0.0"

func main() {
	showHidden := flag.Bool("a", false, "show hidden files")
	all := flag.Bool("all", false, "show hidden files")
	maxDepth := flag.Int("L", 0, "max depth")
	directoriesOnly := flag.Bool("d", false, "list directories only")
	dirsOnlyLong := flag.Bool("directories", false, "list directories only")
	sortedS := flag.Bool("s", false, "sorted files")
	sortedLong := flag.Bool("sorted", false, "sorted files")
	summaryMode := flag.Bool("m", false, "summary mode (collapses deep structures)")
	summaryLong := flag.Bool("summary", false, "summary mode (collapses deep structures)")
	showVersion := flag.Bool("v", false, "show version")
	showVersionLong := flag.Bool("version", false, "show version")
	orderByExt := flag.Bool("o", false, "sort files by extension")
	orderByExtLong := flag.Bool("order-by-extension", false, "sort files by extension")
	iconModeStr := flag.String("icons", "nerd", "icon set: nerd, unicode, ascii, none")
	iconModeStrAlias := flag.String("i", "nerd", "icon set: nerd, unicode, ascii, none")
	theme := flag.String("theme", "auto", "theme mode: auto, color, mono")
	du := flag.Bool("du", false, "show disk usage sizes")
	jsonOut := flag.Bool("json", false, "output as JSON")
	yamlOut := flag.Bool("yaml", false, "output as YAML")
	sortBy := flag.String("sort", "", "sort mode: name, ext, size, mtime")
	reverse := flag.Bool("reverse", false, "reverse sort order")
	completion := flag.String("completion", "", "generate shell completion: bash, zsh, fish, powershell")

	var includePatterns stringListFlag
	var ignorePatterns stringListFlag
	flag.Var(&includePatterns, "include", "include glob pattern (repeatable)")
	flag.Var(&ignorePatterns, "ignore", "ignore glob pattern (repeatable)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: gotree [options] [path]\n\nOptions:\n")
		fmt.Fprintf(os.Stderr, "  -h, --help                     Show help\n")
		fmt.Fprintf(os.Stderr, "  -v, --version                  Show version\n")
		fmt.Fprintf(os.Stderr, "  -a, --all                      Show hidden files\n")
		fmt.Fprintf(os.Stderr, "  -s, --sorted                   Sort by name (alias)\n")
		fmt.Fprintf(os.Stderr, "  -o, --order-by-extension       Sort by extension (alias)\n")
		fmt.Fprintf(os.Stderr, "  --sort <name|ext|size|mtime>   Sort mode\n")
		fmt.Fprintf(os.Stderr, "  --reverse                      Reverse sort order\n")
		fmt.Fprintf(os.Stderr, "  --include <glob>               Include glob pattern (repeatable)\n")
		fmt.Fprintf(os.Stderr, "  --ignore <glob>                Ignore glob pattern (repeatable)\n")
		fmt.Fprintf(os.Stderr, "  -m, --summary                  Summary mode\n")
		fmt.Fprintf(os.Stderr, "  -d, --directories              List directories only\n")
		fmt.Fprintf(os.Stderr, "  -L <depth>                     Max depth\n")
		fmt.Fprintf(os.Stderr, "  -i, --icons <mode>             Icon mode: nerd, unicode, ascii, none\n")
		fmt.Fprintf(os.Stderr, "  --theme <auto|color|mono>      Color theme mode\n")
		fmt.Fprintf(os.Stderr, "  --du                           Show recursive size totals\n")
		fmt.Fprintf(os.Stderr, "  --json                         Output structured JSON\n")
		fmt.Fprintf(os.Stderr, "  --yaml                         Output structured YAML\n")
		fmt.Fprintf(os.Stderr, "  --completion <shell>           Print shell completion script\n")
	}

	flag.Parse()

	if *showVersion || *showVersionLong {
		fmt.Printf("gotree version %s\n", version)
		os.Exit(0)
	}

	if *completion != "" {
		script, err := completionScript(*completion)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(script)
		os.Exit(0)
	}

	if *jsonOut && *yamlOut {
		fmt.Println("error: --json and --yaml are mutually exclusive")
		os.Exit(1)
	}

	iconsMode := *iconModeStr
	if *iconModeStrAlias != "nerd" {
		iconsMode = *iconModeStrAlias
	}

	sortMode, err := resolveSortMode(*sortBy, *sortedS || *sortedLong, *orderByExt || *orderByExtLong)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	colorEnabled, err := resolveColorEnabled(*theme)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	root := "."
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}

	tree, err := scanner.BuildTree(root, scanner.Config{
		ShowHidden:      *showHidden || *all,
		MaxDepth:        *maxDepth,
		DirectoriesOnly: *directoriesOnly || *dirsOnlyLong,
		SummaryMode:     *summaryMode || *summaryLong,
		SortBy:          sortMode,
		Reverse:         *reverse,
		IncludePatterns: includePatterns.Values(),
		IgnorePatterns:  ignorePatterns.Values(),
	})
	if err != nil {
		fmt.Printf("error building tree: %v\n", err)
		os.Exit(1)
	}

	if *jsonOut || *yamlOut {
		if err := printStructured(tree, *jsonOut, *yamlOut); err != nil {
			fmt.Printf("error writing structured output: %v\n", err)
			os.Exit(1)
		}
		return
	}

	printCfg := printer.Config{
		IconMode:      icon.ParseMode(iconsMode),
		EnableColor:   colorEnabled,
		ShowSize:      *du,
		ShowSizeTotal: *du,
	}
	printer.Print(tree, printCfg)
	printer.PrintSummary(tree, printCfg)
}

func resolveSortMode(sortBy string, oldSorted, oldExt bool) (string, error) {
	mode := strings.ToLower(strings.TrimSpace(sortBy))
	if mode == "" {
		if oldExt {
			return "ext", nil
		}
		if oldSorted {
			return "name", nil
		}
		return "name", nil
	}
	switch mode {
	case "name", "ext", "size", "mtime":
		return mode, nil
	default:
		return "", fmt.Errorf("invalid --sort value %q (expected: name, ext, size, mtime)", sortBy)
	}
}

func resolveColorEnabled(theme string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(theme)) {
	case "auto", "":
		return os.Getenv("NO_COLOR") == "", nil
	case "color":
		return true, nil
	case "mono":
		return false, nil
	default:
		return false, errors.New("invalid --theme value (expected: auto, color, mono)")
	}
}
