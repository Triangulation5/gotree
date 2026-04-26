// Package scanner
package scanner

import (
	"bufio"
	"cmp"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gotree/src/internal/model"
)

type Config struct {
	ShowHidden      bool
	MaxDepth        int
	DirectoriesOnly bool
	SummaryMode     bool
	SortBy          string
	Reverse         bool
	IncludePatterns []string
	IgnorePatterns  []string
}

type walkContext struct {
	rootAbs string
	cfg     Config
	ignores []string
}

func BuildTree(root string, cfg Config) (*model.Node, error) {
	ignoreList := loadGitIgnore(root)
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	ctx := walkContext{
		rootAbs: absRoot,
		cfg:     cfg,
		ignores: ignoreList,
	}
	return walk(root, 0, ctx)
}

func walk(path string, depth int, ctx walkContext) (*model.Node, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}

	node := &model.Node{
		Name:         info.Name(),
		Path:         path,
		IsDir:        info.IsDir(),
		SizeBytes:    info.Size(),
		ModifiedUnix: info.ModTime().Unix(),
	}
	if info.IsDir() {
		node.SizeBytes = 0
	}

	if !info.IsDir() {
		return node, nil
	}

	if ctx.cfg.MaxDepth > 0 && depth >= ctx.cfg.MaxDepth {
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
		childPath := filepath.Join(path, name)

		// .gitignore support - only filter in summary mode
		if ctx.cfg.SummaryMode && isIgnored(name, ctx.ignores) {
			if e.IsDir() {
				ignoredDirs++
			} else {
				ignoredFiles++
			}
			continue
		}

		// Hidden files support
		if !ctx.cfg.ShowHidden && name[0] == '.' {
			if e.IsDir() {
				hiddenDirs++
			} else {
				hiddenFiles++
			}
			continue
		}

		if ctx.cfg.DirectoriesOnly && !e.IsDir() {
			continue
		}

		includeMatch, ignoreMatch := pathMatches(ctx.rootAbs, childPath, ctx.cfg.IncludePatterns, ctx.cfg.IgnorePatterns)
		if includeMatch {
			visibleEntries = append(visibleEntries, e)
			continue
		}
		if ignoreMatch {
			continue
		}
		if len(ctx.cfg.IncludePatterns) > 0 && !e.IsDir() {
			continue
		}

		visibleEntries = append(visibleEntries, e)
	}

	// Summary mode: scale-aware pruning & noise suppression
	if ctx.cfg.SummaryMode {
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
		child, err := walk(childPath, depth+1, ctx)
		if err != nil {
			return nil, err
		}
		includeMatch, _ := pathMatches(ctx.rootAbs, childPath, ctx.cfg.IncludePatterns, ctx.cfg.IgnorePatterns)
		if len(ctx.cfg.IncludePatterns) > 0 && child.IsDir && !includeMatch && len(child.Children) == 0 {
			continue
		}
		node.Children = append(node.Children, child)
		node.SizeBytes += child.SizeBytes
	}

	sortChildren(node.Children, ctx.cfg.SortBy, ctx.cfg.Reverse)

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
	return strings.ToLower(name[i+1:])
}

func pathMatches(rootAbs, fullPath string, includes, ignores []string) (bool, bool) {
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return false, false
	}
	rel, err := filepath.Rel(rootAbs, absPath)
	if err != nil {
		return false, false
	}
	rel = filepath.ToSlash(rel)
	base := filepath.Base(rel)

	matchedInclude := matchAny(includes, rel, base)
	matchedIgnore := matchAny(ignores, rel, base)
	return matchedInclude, matchedIgnore && !matchedInclude
}

func matchAny(patterns []string, rel, base string) bool {
	for _, pattern := range patterns {
		p := strings.TrimSpace(filepath.ToSlash(pattern))
		if p == "" {
			continue
		}
		if ok, _ := filepath.Match(p, rel); ok {
			return true
		}
		if ok, _ := filepath.Match(p, base); ok {
			return true
		}
	}
	return false
}

func sortChildren(children []*model.Node, sortBy string, reverse bool) {
	if len(children) < 2 {
		return
	}
	normalized := strings.ToLower(strings.TrimSpace(sortBy))
	sort.SliceStable(children, func(i, j int) bool {
		a, b := children[i], children[j]
		if a.IsSummary != b.IsSummary {
			return !a.IsSummary
		}
		if a.IsSummary && b.IsSummary {
			return a.SummaryMsg < b.SummaryMsg
		}
		if a.IsDir != b.IsDir {
			return a.IsDir
		}

		var cmpVal int
		switch normalized {
		case "ext":
			cmpVal = compareExtensions(a.Name, b.Name)
		case "size":
			cmpVal = cmp.Compare(a.SizeBytes, b.SizeBytes)
		case "mtime":
			cmpVal = cmp.Compare(a.ModifiedUnix, b.ModifiedUnix)
		default:
			cmpVal = cmp.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
		}
		if cmpVal == 0 {
			cmpVal = cmp.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
		}
		if reverse {
			return cmpVal > 0
		}
		return cmpVal < 0
	})
}

func compareExtensions(a, b string) int {
	extA := getExtension(a)
	extB := getExtension(b)
	if extA == extB {
		return 0
	}
	if extA == "" {
		return 1
	}
	if extB == "" {
		return -1
	}
	return cmp.Compare(extA, extB)
}
