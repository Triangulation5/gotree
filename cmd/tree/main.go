// Command gotree is a lightweight CLI that prints a directory hierarchy.
//
// It is a thin orchestration layer that:
//   - Parses CLI flags
//   - Delegates filesystem traversal to the scanner package
//   - Delegates rendering to the printer package
package main

import (
    "flag"
    "fmt"
    "os"

    "gotree/internal/printer"
    "gotree/internal/scanner"
)

const version = "1.0.0"

func main() {
	showHidden := flag.Bool("a", false, "show hidden files")
	all := flag.Bool("all", false, "show hidden files")
	maxDepth := flag.Int("L", 0, "max depth")
	directoriesOnly := flag.Bool("d", false, "list directories only")
	dirsOnlyLong := flag.Bool("directories", false, "list directories only")
	sortedS := flag.Bool("s", false, "sorted files")
	sortedLong := flag.Bool("sorted", false, "sorted files")
	summary := flag.Bool("m", false, "show summary")
	summaryLong := flag.Bool("summary", false, "show summary")
	showVersion := flag.Bool("v", false, "show version")
	showVersionLong := flag.Bool("version", false, "show version")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: gotree [options] [path]\n\nOptions:\n")
		fmt.Fprintf(os.Stderr, "  -h, --help         Show help\n")
		fmt.Fprintf(os.Stderr, "  -v, --version      Show version\n")
		fmt.Fprintf(os.Stderr, "  -a, --all          Show hidden files\n")
		fmt.Fprintf(os.Stderr, "  -s, --sorted       Alphabetically sorted files\n")
		fmt.Fprintf(os.Stderr, "  -m, --summary      Show summary\n")
		fmt.Fprintf(os.Stderr, "  -d, --directories  List directories only\n")
		fmt.Fprintf(os.Stderr, "  -L <depth>         Max depth\n")
	}

	flag.Parse()

	if *showVersion || *showVersionLong {
		fmt.Printf("gotree version %s\n", version)
		os.Exit(0)
	}

	root := "."
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}

	tree, error := scanner.BuildTree(root, scanner.Config{
		ShowHidden:      *showHidden || *all,
		MaxDepth:        *maxDepth,
		DirectoriesOnly: *directoriesOnly || *dirsOnlyLong,
		Sorted:          *sortedS || *sortedLong,
	})
	if error != nil {
		fmt.Println("error:", error)
		os.Exit(1)
	}

	printer.Print(tree)

	if *summary || *summaryLong {
		printer.PrintSummary(tree)
	}
}
