package git

import "sync"

// Run caches the stable location answers of one command run: the common directory, the
// checkout's administration directory, and the top level. The key is the root spelling
// that the caller passes. A command creates one Run at its entry and passes it down. The
// zero value is ready to use.
//
// The cache saves the git start, not the filesystem check. Each directory read validates
// the stored answer again, so a directory that a write replaced after the first read still
// refuses. A Run holds no ref, HEAD, registration, index, default-branch, or ancestry
// answer, because a write inside the run changes those answers. Never hold a Run in a
// package variable: a test that swaps PATH or recreates a repository needs a fresh answer.
type Run struct {
	mu      sync.Mutex
	answers map[runQuery]string
}

// runQuery names one cached raw answer: the location query and the root spelling.
type runQuery struct {
	location location
	root     string
}

type location int

const (
	commonDirLocation location = iota
	adminDirLocation
	topLevelLocation
)

// CommonDir answers as the package CommonDir does, and it starts git once for each root
// spelling in the run.
func (r *Run) CommonDir(root string) (string, error) {
	return r.directory(runQuery{commonDirLocation, root}, subjectCommonDir, commonDirArgs(root))
}

// AdminDir answers as the package AdminDir does, and it starts git once for each root
// spelling in the run.
func (r *Run) AdminDir(root string) (string, error) {
	return r.directory(runQuery{adminDirLocation, root}, subjectAdminDir, adminDirArgs(root))
}

// RootAt answers as the package RootAt does, and it starts git once for each directory
// spelling in the run.
func (r *Run) RootAt(dir string) (string, error) {
	return r.answer(runQuery{topLevelLocation, dir}, func() (string, error) { return RootAt(dir) })
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

func (r *Run) directory(query runQuery, subject string, args []string) (string, error) {
	raw, err := r.answer(query, func() (string, error) { return locationAnswer(subject, args) })
	if err != nil {
		return "", err
	}
	return validateCommonDir(raw, subject)
}

// answer returns the stored raw answer of query, or it starts the query and stores a
// successful answer. A failed start stores nothing, so a later read starts git again.
func (r *Run) answer(query runQuery, start func() (string, error)) (string, error) {
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
