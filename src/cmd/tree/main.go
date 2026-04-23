package main

import (
	"flag"
	"fmt"
	"os"

	"gotree/src/config/icon"
	"gotree/src/internal/printer"
	"gotree/src/internal/scanner"
)

const version = "1.1.0"

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
	iconModeStrAlias := flag.String("i", "nerd", "icon set: nerd, unicode, ascii, none") // Alias for --icons

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: gotree [options] [path]\n\nOptions:\n")
		fmt.Fprintf(os.Stderr, "  -h, --help         Show help\n")
		fmt.Fprintf(os.Stderr, "  -v, --version      Show version\n")
		fmt.Fprintf(os.Stderr, "  -a, --all          Show hidden files\n")
		fmt.Fprintf(os.Stderr, "  -s, --sorted       Alphabetically sorted files\n")
		fmt.Fprintf(os.Stderr, "  -o, --order-by-extension  Sort files by extension\n")
		fmt.Fprintf(os.Stderr, "  -m, --summary      Summary mode (collapses large/noise directories)\n")
		fmt.Fprintf(os.Stderr, "  -d, --directories  List directories only\n")
		fmt.Fprintf(os.Stderr, "  -L <depth>         Max depth\n")
		fmt.Fprintf(os.Stderr, "  -i, --icons <mode>     Icon mode: nerd (default), unicode, ascii, none\n")
	}

	flag.Parse()

	// Determine the final icon mode. Prioritize the alias if it's explicitly set.
	iconsMode := *iconModeStr
	// If the alias flag was set and its value is different from the default "nerd",
	// it means the user explicitly provided the alias. In this case, use the alias value.
	// Otherwise, use the value from the --icons flag.
	if *iconModeStrAlias != "nerd" {
		iconsMode = *iconModeStrAlias
	}

	if *showVersion || *showVersionLong {
		fmt.Printf("gotree version %s\n", version)
		os.Exit(0)
	}

	root := "."
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}

	tree, err := scanner.BuildTree(root, scanner.Config{
		ShowHidden:       *showHidden || *all,
		MaxDepth:         *maxDepth,
		DirectoriesOnly:  *directoriesOnly || *dirsOnlyLong,
		Sorted:           *sortedS || *sortedLong,
		SummaryMode:      *summaryMode || *summaryLong,
		OrderByExtension: *orderByExt || *orderByExtLong,
	})
	if err != nil {
		fmt.Printf("error building tree: %v\n", err)
		os.Exit(1)
	}

	printer.Print(tree, printer.Config{
		IconMode: icon.ParseMode(iconsMode),
	})

	// Final summary line
	printer.PrintSummary(tree)
	}
