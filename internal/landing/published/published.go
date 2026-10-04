// Package published owns the exact tree that a verified delivery publishes, and the
// private-index edits that write it. It imports no gate, so a fixture can publish the
// same tree that a landing publishes.
package published

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/spec"
)

// Tree returns the tree that the verified delivery of the reviewed deliverable at path
// publishes from tree. A staged spec takes its implemented status, with the bytes and the
// mode that the reviewed source commit holds. A tickets-only folder leaves the tree; a
// folder that tree already lacks lists no entries, so its removal changes nothing. The
// commitment owner then derives the exact closure of the delivery from the result, and
// the published tree carries the delivery and its closure together or carries neither.
func Tree(root, tree, path, source string) (string, error) {
	var err error
	if spec.IsLiveSpecPath(path) {
		if tree, err = implementedSpec(root, tree, path, source); err != nil {
			return "", fmt.Errorf("transition staged spec: %w", err)
		}
	} else if tree, err = Edit(root, tree, func(idx string) error { return RemoveFolder(root, idx, path) }); err != nil {
		return "", fmt.Errorf("close tickets-only folder: %w", err)
	}
	if tree, err = closeDelivery(root, tree, path, source); err != nil {
		return "", fmt.Errorf("close verified delivery: %w", err)
	}
	return tree, nil
}

// implementedSpec writes the implemented status of the spec that source holds at path.
func implementedSpec(root, tree, path, source string) (string, error) {
	listed, err := benchgit.Output("-C", root, "ls-tree", source, "--", path)
	fields := strings.Fields(listed)
	if err != nil || len(fields) < 3 || !(benchgit.IndexEntry{Mode: fields[0]}).IsRegularFile() {
		return "", fmt.Errorf("reviewed source spec %s is not a regular file", path)
	}
	staged, err := benchgit.ReadTreeFile(root, source, path)
	if err != nil {
		return "", err
	}
	implemented, err := spec.Implemented(staged)
	if err != nil {
		return "", err
	}
	return Edit(root, tree, func(idx string) error { return WriteFile(root, idx, path, implemented, fields[0]) })
}

// closeDelivery applies the verified closure of the reviewed deliverable to tree. The
// commitment owner derives the exact edits, and publication admission and the gate's
// completion oracle consume the same derivation.
func closeDelivery(root, tree, path, source string) (string, error) {
	edits, err := commitrepo.Store{Root: root}.Closure(tree, commitrepo.Delivery{Spec: path, Source: source})
	if err != nil || len(edits) == 0 {
		return tree, err
	}
	return Edit(root, tree, func(idx string) error {
		for _, edit := range edits {
			if edit.Delete {
				if err := benchgit.IndexCommand(root, idx, "update-index", "--force-remove", "--", edit.Path).Run(); err != nil {
					return err
				}
				continue
			}
			if err := WriteFile(root, idx, edit.Path, edit.Data, edit.Mode); err != nil {
				return err
			}
		}
		return nil
	})
}

// Edit reads baseTree into a private index, applies edit to that index, and writes the
// resulting tree. No checkout or repository index is touched.
func Edit(root, baseTree string, edit func(idx string) error) (string, error) {
	dir, err := os.MkdirTemp("", "bench-reviewed-landing-index-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)
	idx := filepath.Join(dir, "index")
	if err := benchgit.IndexCommand(root, idx, "read-tree", baseTree).Run(); err != nil {
		return "", err
	}
	if err := edit(idx); err != nil {
		return "", err
	}
	tree, err := benchgit.IndexCommand(root, idx, "write-tree").Output()
	return strings.TrimSpace(string(tree)), err
}

// WriteFile stores content as one regular index entry at path with the Git mode.
func WriteFile(root, idx, path string, content []byte, mode string) error {
	hash := exec.Command("git", "-C", root, "hash-object", "-w", "--stdin")
	hash.Stdin = bytes.NewReader(content)
	blob, err := hash.Output()
	if err != nil {
		return err
	}
	return benchgit.IndexCommand(root, idx, "update-index", "--add", "--cacheinfo", mode+","+strings.TrimSpace(string(blob))+","+path).Run()
}

// RemoveFolder drops every entry beneath rel from the private index, so the published
// tree carries the deletion rather than the checkout carrying it afterwards. The pathspec
// is literal and the removals name exact index paths, so a folder name holding a space or
// a glob character resolves to itself.
func RemoveFolder(root, idx, rel string) error {
	listed, err := benchgit.IndexCommand(root, idx, "ls-files", "-z", "--cached", "--", ":(literal)"+rel).Output()
	if err != nil {
		return fmt.Errorf("list tracked entries under %q: %w", rel, err)
	}
	for _, path := range strings.Split(string(listed), "\x00") {
		if path == "" {
			continue
		}
		if err := benchgit.IndexCommand(root, idx, "update-index", "--force-remove", "--", path).Run(); err != nil {
			return fmt.Errorf("remove %q from prospective index: %w", path, err)
		}
	}
	return nil
}
