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

func main() {
    showHidden := flag.Bool("a", false, "show hidden files")
    maxDepth := flag.Int("L", 0, "max depth")
    flag.Parse()

    root := "."
    if flag.NArg() > 0 {
        root = flag.Arg(0)
    }

    tree, err := scanner.BuildTree(root, scanner.Config{
        ShowHidden: *showHidden,
        MaxDepth:   *maxDepth,
    })
    if err != nil {
        fmt.Println("error:", err)
        os.Exit(1)
    }

    printer.Print(tree)
}
