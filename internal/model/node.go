// Package model defines the core in-memory representation of a filesystem tree.
package model

type Node struct {
    Name     string
    Path     string
    IsDir    bool
    Children []*Node
}
