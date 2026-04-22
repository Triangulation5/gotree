// Package scanner is responsible for traversing the filesystem and building
// an in-memory tree structure.
//
// It converts a real directory hierarchy into a structured model.Node tree.
// This package performs all filesystem I/O and recursion logic, but does not
// handle any output formatting or presentation concerns.
package scanner

import (
	"os"
	"path/filepath"
	"sort"

	"gotree/internal/model"
)

type Config struct {
	ShowHidden bool
	MaxDepth   int
}

func BuildTree(root string, cfg Config) (*model.Node, error) {
	return walk(root, 0, cfg)
}

func walk(path string, depth int, cfg Config) (*model.Node, error) {
    info, err := os.Lstat(path)
    if err != nil {
        return nil, err
    }

    node := &model.Node{
        Name:  info.Name(),
        Path:  path,
        IsDir: info.IsDir(),
    }

    if !info.IsDir() {
        return node, nil
    }

    if cfg.MaxDepth > 0 && depth >= cfg.MaxDepth {
        return node, nil
    }

    entries, err := os.ReadDir(path)
    if err != nil {
        return nil, err
    }

    sort.Slice(entries, func(i, j int) bool {
        return entries[i].Name() < entries[j].Name()
    })

    for _, e := range entries {
        name := e.Name()

        if !cfg.ShowHidden && name[0] == '.' {
            continue
        }

        childPath := filepath.Join(path, name)

        child, err := walk(childPath, depth+1, cfg)
        if err != nil {
            return nil, err
        }

        node.Children = append(node.Children, child)
    }

    return node, nil
}
