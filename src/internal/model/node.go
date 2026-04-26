// Package model defines the core in-memory representation of a filesystem tree.
package model

type Node struct {
	Name         string
	Path         string
	IsDir        bool
	SizeBytes    int64
	ModifiedUnix int64
	Children     []*Node
	IsSummary    bool
	SummaryMsg   string
}

func (n *Node) Count() (dirs, files int) {
	for _, child := range n.Children {
		if child.IsSummary {
			continue
		}
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
