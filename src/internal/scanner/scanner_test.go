package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"gotree/src/internal/model"
)

func TestIncludeWinsOverIgnore(t *testing.T) {
	root := t.TempDir()
	mustWriteFile(t, filepath.Join(root, "keep.go"), "package main\n")
	mustWriteFile(t, filepath.Join(root, "skip.txt"), "skip\n")

	tree, err := BuildTree(root, Config{
		IncludePatterns: []string{"*.go"},
		IgnorePatterns:  []string{"*.go", "*.txt"},
		SortBy:          "name",
	})
	if err != nil {
		t.Fatalf("BuildTree failed: %v", err)
	}

	children := childNames(tree)
	if len(children) != 1 || children[0] != "keep.go" {
		t.Fatalf("expected only keep.go, got %v", children)
	}
}

func TestIncludeMatchesNestedFiles(t *testing.T) {
	root := t.TempDir()
	mustWriteFile(t, filepath.Join(root, "a", "b", "target.go"), "package main\n")
	mustWriteFile(t, filepath.Join(root, "a", "b", "skip.txt"), "skip\n")

	tree, err := BuildTree(root, Config{
		IncludePatterns: []string{"*.go"},
		SortBy:          "name",
	})
	if err != nil {
		t.Fatalf("BuildTree failed: %v", err)
	}

	if len(tree.Children) != 1 || tree.Children[0].Name != "a" {
		t.Fatalf("expected intermediate directory a to remain, got %+v", childNames(tree))
	}
	if len(tree.Children[0].Children) != 1 || tree.Children[0].Children[0].Name != "b" {
		t.Fatalf("expected nested directory b to remain")
	}
	names := childNames(tree.Children[0].Children[0])
	if len(names) != 1 || names[0] != "target.go" {
		t.Fatalf("expected only target.go under nested folder, got %v", names)
	}
}

func TestSortBySizeReverseAndDirPriority(t *testing.T) {
	root := t.TempDir()
	mustWriteFile(t, filepath.Join(root, "small.txt"), "123")
	mustWriteFile(t, filepath.Join(root, "large.txt"), "123456789")
	mustMkdir(t, filepath.Join(root, "dir"))
	mustWriteFile(t, filepath.Join(root, "dir", "nested.txt"), "x")

	tree, err := BuildTree(root, Config{
		SortBy:  "size",
		Reverse: true,
	})
	if err != nil {
		t.Fatalf("BuildTree failed: %v", err)
	}

	names := childNames(tree)
	if len(names) != 3 {
		t.Fatalf("expected three children, got %v", names)
	}
	if names[0] != "dir" {
		t.Fatalf("expected directory first even on reverse sort, got %v", names)
	}
	if names[1] != "large.txt" || names[2] != "small.txt" {
		t.Fatalf("unexpected size reverse order: %v", names)
	}
}

func TestDirectorySizeIsRecursiveTotal(t *testing.T) {
	root := t.TempDir()
	mustWriteFile(t, filepath.Join(root, "a.txt"), "1234")
	mustWriteFile(t, filepath.Join(root, "nested", "b.txt"), "12")

	tree, err := BuildTree(root, Config{SortBy: "name"})
	if err != nil {
		t.Fatalf("BuildTree failed: %v", err)
	}

	if tree.SizeBytes != 6 {
		t.Fatalf("expected root size 6, got %d", tree.SizeBytes)
	}
}

func mustMkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
}

func mustWriteFile(t *testing.T, path, content string) {
	t.Helper()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func childNames(node *model.Node) []string {
	names := make([]string, 0, len(node.Children))
	for _, child := range node.Children {
		if child.IsSummary {
			continue
		}
		names = append(names, child.Name)
	}
	return names
}
