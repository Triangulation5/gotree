// Package model defines the core in-memory representation of a filesystem tree.
package model

type Node struct {
    Name     string
    Path     string
    IsDir    bool
    Children []*Node
}

func (n *Node) Count() (dirs, files int) {
	for _, child := range n.Children {
		if child.IsDir {
			dirs++
		} else {
			files++
		}
		d, f := child.Count()
		dirs += d
		files += f
	}
	return
}
