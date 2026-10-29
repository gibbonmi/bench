package git

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// Run caches the pure git answers of one command run. A pure answer is a location answer,
// an administration path, or a read whose revision operands are all full object ids. The
// location key is the root spelling that the caller passes, and the object key is the exact
// argv. A command creates one Run at its entry and passes it down. The zero value is ready
// to use, and a nil Run starts each query on every call and stores nothing.
//
// The cache saves the git start, not the filesystem check. Each directory read validates
// the stored answer again, so a directory that a write replaced after the first read still
// refuses. A Run holds no ref, HEAD, short-id, registration, index, status, or
// default-branch answer, because a write inside the run changes those answers. Never hold a
// Run in a package variable: a test that swaps PATH or recreates a repository needs a fresh
// answer.
type Run struct {
	mu      sync.Mutex
	answers map[runQuery]string
	objects map[string][]byte
}

// runQuery names one cached raw location answer: the location query, the root spelling,
// and the file name of an administration path.
type runQuery struct {
	location location
	root     string
	name     string
}

type location int

const (
	commonDirLocation location = iota
	adminDirLocation
	topLevelLocation
	adminPathLocation
)

// FullObjectID reports whether value is a full object id: 40 or 64 hexadecimal digits of
// either case. It is the one purity predicate of the run's object reads.
func FullObjectID(value string) bool {
	if len(value) != 40 && len(value) != 64 {
		return false
	}
	for _, b := range []byte(value) {
		if !((b >= '0' && b <= '9') || (b >= 'a' && b <= 'f') || (b >= 'A' && b <= 'F')) {
			return false
		}
	}
	return true
}

// CommonDir answers as the package CommonDir does, and it starts git once for each root
// spelling in the run.
func (r *Run) CommonDir(root string) (string, error) {
	return r.directory(runQuery{location: commonDirLocation, root: root}, subjectCommonDir, commonDirArgs(root))
}

// AdminDir answers as the package AdminDir does, and it starts git once for each root
// spelling in the run.
func (r *Run) AdminDir(root string) (string, error) {
	return r.directory(runQuery{location: adminDirLocation, root: root}, subjectAdminDir, adminDirArgs(root))
}

// RootAt answers as the package RootAt does, and it starts git once for each directory
// spelling in the run.
func (r *Run) RootAt(dir string) (string, error) {
	return r.answer(runQuery{location: topLevelLocation, root: dir}, func() (string, error) { return RootAt(dir) })
}

// AdminPath starts git once for each root spelling and name in the run. On a non-nil run,
// each read checks the parent directory of the answer, the first read included, so a
// missing parent refuses where the package AdminPath answers. A nil run answers as the
// package AdminPath does.
func (r *Run) AdminPath(root, name string) (string, error) {
	if r == nil {
		return AdminPath(root, name)
	}
	path, err := r.answer(runQuery{location: adminPathLocation, root: root, name: name}, func() (string, error) { return AdminPath(root, name) })
	if err != nil {
		return "", err
	}
	if _, err := validateCommonDir(filepath.Dir(path), subjectAdminPath); err != nil {
		return "", err
	}
	return path, nil
}

// Worktrees answers as the package Worktrees does. It reads the common directory through
// the run, but it runs the admin scan and starts the porcelain on every call, because a
// registration can change inside the run.
func (r *Run) Worktrees(root string) ([]Worktree, error) {
	common, err := r.CommonDir(root)
	if err != nil {
		return nil, err
	}
	return worktreesAt(root, common)
}

// ReadControlBlob answers as the package ReadControlBlob does. It starts git once for each
// full object id in the run, and it reads any other operand live.
func (r *Run) ReadControlBlob(root, object string) ([]byte, error) {
	if !FullObjectID(object) {
		return ReadControlBlob(root, object)
	}
	return r.objectAnswer(controlBlobArgs(root, object), func() ([]byte, error) { return ReadControlBlob(root, object) })
}

// ListTree answers as the package ListTree does. It starts git once for each argv in the
// run when the listing names exactly one revision before "--", and that revision is a full
// object id with an optional "^{tree}". It lists any other revision live.
func (r *Run) ListTree(root string, args ...string) ([]TreeEntry, error) {
	argv := listTreeArgs(root, args)
	start := func() ([]byte, error) { return Raw(argv...) }
	var raw []byte
	var err error
	if pureListing(args) {
		raw, err = r.objectAnswer(argv, start)
	} else {
		raw, err = start()
	}
	if err != nil {
		return nil, err
	}
	return parseTreeEntries(raw)
}

// pureListing reports whether the ls-tree arguments name exactly one revision before "--",
// and whether that revision is a full object id with an optional "^{tree}".
func pureListing(args []string) bool {
	var revisions []string
	for _, arg := range args {
		if arg == "--" {
			break
		}
		if !strings.HasPrefix(arg, "-") {
			revisions = append(revisions, arg)
		}
	}
	return len(revisions) == 1 && FullObjectID(strings.TrimSuffix(revisions[0], "^{tree}"))
}

// ReadTreeFile reads a regular immutable blob under the control-record bound. It lists
// the path through the run's ListTree, and then it reads the blob through the run's
// ReadControlBlob.
func (r *Run) ReadTreeFile(root, tree, path string) ([]byte, error) {
	listing, err := r.ListTree(root, tree, "--", path)
	if err != nil {
		return nil, err
	}
	if len(listing) != 1 || listing[0].Path != path || !listing[0].IsRegularFile() {
		return nil, fmt.Errorf("missing or nonregular tree file %s", path)
	}
	return r.ReadControlBlob(root, listing[0].Object)
}

// ResolveCommit answers as the package ResolveCommit does. It starts git once for each
// full object id and flag set in the run, and it resolves any other revision live.
func (r *Run) ResolveCommit(root, revision string, flags ...string) (string, error) {
	argv := resolveCommitArgs(root, revision, flags)
	start := func() ([]byte, error) {
		out, err := Output(argv...)
		return []byte(out), err
	}
	if !FullObjectID(revision) {
		out, err := start()
		return string(out), err
	}
	out, err := r.objectAnswer(argv, start)
	return string(out), err
}

// IsAncestor answers as the package IsAncestor does. It stores a yes and a no when both
// operands are full object ids, and it checks any other pair live. An error stores nothing.
func (r *Run) IsAncestor(root, ancestor, descendant string) (bool, error) {
	if !FullObjectID(ancestor) || !FullObjectID(descendant) {
		return IsAncestor(root, ancestor, descendant)
	}
	raw, err := r.objectAnswer(isAncestorArgs(root, ancestor, descendant), func() ([]byte, error) {
		yes, err := IsAncestor(root, ancestor, descendant)
		if err != nil {
			return nil, err
		}
		return []byte(strconv.FormatBool(yes)), nil
	})
	if err != nil {
		return false, err
	}
	return string(raw) == strconv.FormatBool(true), nil
}

func (r *Run) directory(query runQuery, subject string, args []string) (string, error) {
	raw, err := r.answer(query, func() (string, error) { return locationAnswer(subject, args) })
	if err != nil {
		return "", err
	}
	return validateCommonDir(raw, subject)
}

// answer returns the stored raw answer of query, or it starts the query and stores a
// successful answer. A failed start stores nothing, so a later read starts git again. A
// nil run starts the query and stores nothing.
func (r *Run) answer(query runQuery, start func() (string, error)) (string, error) {
	if r == nil {
		return start()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if raw, ok := r.answers[query]; ok {
		return raw, nil
	}
	raw, err := start()
	if err != nil {
		return raw, err
	}
	if r.answers == nil {
		r.answers = map[runQuery]string{}
	}
	r.answers[query] = raw
	return raw, nil
}

// objectAnswer returns a copy of the stored answer of argv, or it starts the read and
// stores a copy of a successful answer. A failed read stores nothing. A nil run starts the
// read and stores nothing. The run holds its mutex while git runs.
func (r *Run) objectAnswer(argv []string, start func() ([]byte, error)) ([]byte, error) {
	if r == nil {
		return start()
	}
	key := strings.Join(argv, "\x00")
	r.mu.Lock()
	defer r.mu.Unlock()
	if raw, ok := r.objects[key]; ok {
		return bytes.Clone(raw), nil
	}
	raw, err := start()
	if err != nil {
		return raw, err
	}
	if r.objects == nil {
		r.objects = map[string][]byte{}
	}
	r.objects[key] = bytes.Clone(raw)
	return raw, nil
}
