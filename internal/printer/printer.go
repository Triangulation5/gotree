// Package printer is responsible for rendering a model.Node tree into a
// human-readable ASCII representation.
//
// It does not perform any filesystem operations or tree construction.
// Its only responsibility is presentation of an already-built tree.
package printer

import (
    "fmt"

    "gotree/internal/model"
)

func Print(root *model.Node) {
	fmt.Println(root.Name)
	for i, child := range root.Children {
		printNode(child, "", i == len(root.Children)-1)
}
}

func printNode(n *model.Node, prefix string, isLast bool) {
    connector := "├── "
    nextPrefix := "│   "

    if isLast {
        connector = "└── "
        nextPrefix = "    "
    }

    fmt.Println(prefix + connector + n.Name)

    for i, child := range n.Children {
        last := i == len(n.Children)-1
        printNode(child, prefix+nextPrefix, last)
    }
}
