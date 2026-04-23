package scanner

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gotree/src/internal/model"
)

type Config struct {
	ShowHidden       bool
	MaxDepth         int
	DirectoriesOnly  bool
	Sorted           bool
	SummaryMode      bool
	OrderByExtension bool
}

func BuildTree(root string, cfg Config) (*model.Node, error) {
	ignoreList := loadGitIgnore(root)
	return walk(root, 0, cfg, ignoreList)
}

func walk(path string, depth int, cfg Config, ignores []string) (*model.Node, error) {
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

	var visibleEntries []os.DirEntry
	var hiddenFiles, hiddenDirs int
	var ignoredFiles, ignoredDirs int

	for _, e := range entries {
		name := e.Name()

		// .gitignore support - only filter in summary mode
		if cfg.SummaryMode && isIgnored(name, ignores) {
			if e.IsDir() {
				ignoredDirs++
			} else {
				ignoredFiles++
			}
			continue
		}

		// Hidden files support
		if !cfg.ShowHidden && name[0] == '.' {
			if e.IsDir() {
				hiddenDirs++
			} else {
				hiddenFiles++
			}
			continue
		}

		if cfg.DirectoriesOnly && !e.IsDir() {
			continue
		}

		visibleEntries = append(visibleEntries, e)
	}

	// Sort visible entries
	if cfg.OrderByExtension {
		sort.Slice(visibleEntries, func(i, j int) bool {
			ei, ej := visibleEntries[i], visibleEntries[j]
			if ei.IsDir() != ej.IsDir() {
				return ei.IsDir() // directories first
			}
			if !ei.IsDir() {
				extI := getExtension(ei.Name())
				extJ := getExtension(ej.Name())
				if extI != extJ {
					if extI == "" {
						return false
					}
					if extJ == "" {
						return true
					}
					return extI < extJ
				}
			}
			return ei.Name() < ej.Name()
		})
	} else if cfg.Sorted {
		sort.Slice(visibleEntries, func(i, j int) bool {
			return visibleEntries[i].Name() < visibleEntries[j].Name()
		})
	}

	// Summary mode: scale-aware pruning & noise suppression
	if cfg.SummaryMode {
		if depth > 0 && (isNoise(node.Name) || len(visibleEntries) > 50) {
			node.Children = append(node.Children, &model.Node{
				IsSummary:  true,
				SummaryMsg: fmtSummary(len(visibleEntries)),
			})
			return node, nil
		}
	}

	for _, e := range visibleEntries {
		childPath := filepath.Join(path, e.Name())
		child, err := walk(childPath, depth+1, cfg, ignores)
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, child)
	}

	// Add aggregate summaries for unlisted items
	if ignoredFiles > 0 || ignoredDirs > 0 {
		node.Children = append(node.Children, &model.Node{
			IsSummary:  true,
			SummaryMsg: fmtUnlisted("ignored by .gitignore", ignoredFiles, ignoredDirs),
		})
	}
	if hiddenFiles > 0 || hiddenDirs > 0 {
		node.Children = append(node.Children, &model.Node{
			IsSummary:  true,
			SummaryMsg: fmtUnlisted("hidden", hiddenFiles, hiddenDirs),
		})
	}

	return node, nil
}

func isNoise(name string) bool {
	noise := []string{"node_modules", ".git", "vendor", "dist", "build", ".vscode", ".idea"}
	for _, n := range noise {
		if name == n {
			return true
		}
	}
	return false
}

func fmtSummary(count int) string {
	return fmt.Sprintf("... (%d items collapsed)", count)
}

func fmtUnlisted(reason string, files, dirs int) string {
	var parts []string
	if dirs > 0 {
		dStr := "directories"
		if dirs == 1 {
			dStr = "directory"
		}
		parts = append(parts, fmt.Sprintf("%d %s", dirs, dStr))
	}
	if files > 0 {
		fStr := "files"
		if files == 1 {
			fStr = "file"
		}
		parts = append(parts, fmt.Sprintf("%d %s", files, fStr))
	}
	return fmt.Sprintf("unlisted (%s: %s)", reason, strings.Join(parts, ", "))
}

// Simple gitignore support
func loadGitIgnore(root string) []string {
	var ignores []string
	f, err := os.Open(filepath.Join(root, ".gitignore"))
	if err != nil {
		return nil
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		ignores = append(ignores, line)
	}
	return ignores
}

func isIgnored(name string, ignores []string) bool {
	for _, pattern := range ignores {
		// Very basic pattern matching
		if name == pattern || strings.HasSuffix(name, pattern) || strings.HasPrefix(name, pattern) {
			return true
		}
		if strings.HasPrefix(pattern, "*.") && strings.HasSuffix(name, pattern[1:]) {
			return true
		}
	}
	return false
}

func getExtension(name string) string {
	i := strings.LastIndex(name, ".")
	if i == -1 {
		return ""
	}
	return name[i+1:]
}
