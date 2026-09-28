package gate

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/canonicalpath"
	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/git"
)

type checkoutGuardKey struct{}

type checkoutFile struct {
	size     int64
	modified time.Time
}

type checkoutGuard struct {
	root    string
	tracked map[string]bool
	allowed map[string]bool
	before  map[string]checkoutFile
	err     error
}

func withCheckoutGuard(ctx context.Context, root string) context.Context {
	root, err := canonicalpath.Resolve(root)
	guard := &checkoutGuard{root: root, tracked: make(map[string]bool), allowed: make(map[string]bool)}
	var tracked []byte
	if err == nil {
		tracked, err = git.Raw("-C", root, "ls-files", "-z", "--cached")
	}
	if err == nil {
		// Keep the initial index: a phase must not hide its new files by staging them.
		for _, path := range strings.Split(string(tracked), "\x00") {
			guard.tracked[filepath.FromSlash(path)] = true
		}
		var admin string
		admin, err = git.AdminDir(root)
		if err == nil {
			paths := []string{registry.TimingPath(root), gateLockPath(admin), gateOwnerPath(admin)}
			if log, _ := ctx.Value(gateRunLogKey{}).(*gateRunLog); log != nil {
				paths = append(paths, log.file.Name(), log.streamPath())
			}
			for _, path := range paths {
				if path == "" {
					continue
				}
				var resolved string
				resolved, err = canonicalpath.Resolve(path)
				if err != nil {
					break
				}
				guard.allowed[resolved] = true
			}
			if err == nil {
				guard.before, err = guard.snapshot()
			}
		}
	}
	guard.err = err
	return context.WithValue(ctx, checkoutGuardKey{}, guard)
}

func checkoutReport(ctx context.Context) func() ([]string, string, bool) {
	return func() ([]string, string, bool) {
		guard, _ := ctx.Value(checkoutGuardKey{}).(*checkoutGuard)
		if guard == nil {
			return nil, "", false
		}
		err := guard.err
		var after map[string]checkoutFile
		if err == nil {
			after, err = guard.snapshot()
		}
		if err != nil {
			return []string{fmt.Sprintf("checkout guard: unreadable checkout state: %q", err.Error())}, "", true
		}
		changed := make(map[string]bool)
		for path, before := range guard.before {
			if current, exists := after[path]; !exists || current != before {
				changed[path] = true
			}
		}
		for path := range after {
			if _, exists := guard.before[path]; !exists {
				changed[path] = true
			}
		}
		paths := make([]string, 0, len(changed))
		for path := range changed {
			paths = append(paths, path)
		}
		sort.Strings(paths)
		var rows []string
		for _, path := range paths {
			rows = append(rows, fmt.Sprintf("checkout guard: changed path %q", path))
		}
		return rows, "", len(rows) != 0
	}
}

func (g *checkoutGuard) snapshot() (map[string]checkoutFile, error) {
	if _, err := os.Lstat(filepath.Join(g.root, ".git")); err != nil {
		return nil, err
	}
	admin, err := git.AdminDir(g.root)
	if err != nil {
		return nil, err
	}
	state := make(map[string]checkoutFile)
	for _, tree := range []struct{ root, prefix string }{{g.root, ""}, {admin, ".git"}} {
		err := filepath.WalkDir(tree.root, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			rel, err := filepath.Rel(tree.root, path)
			if err != nil {
				return err
			}
			if tree.prefix == "" && (rel == ".git" || g.tracked[rel]) {
				if entry.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if entry.IsDir() || g.allowed[path] || (tree.prefix != "" && !benchAdminFile(rel)) {
				return nil
			}
			info, err := entry.Info()
			if err != nil {
				return err
			}
			name := filepath.ToSlash(filepath.Join(tree.prefix, rel))
			state[name] = checkoutFile{size: info.Size(), modified: info.ModTime()}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return state, nil
}

// Bench owns its named administration files and namespaces, including nested refs.
func benchAdminFile(path string) bool {
	for _, part := range strings.Split(filepath.ToSlash(path), "/") {
		if part == "bench" || strings.HasPrefix(part, "bench-") || strings.HasPrefix(part, ".bench-") {
			return true
		}
	}
	return false
}
