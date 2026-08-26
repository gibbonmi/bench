package git

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// PorcelainEntry is one record of `git status --porcelain -z --no-renames`: the two
// status characters (XY) and the path. With --no-renames a record is always
// `XY <path>`, so the path begins at byte 3.
type PorcelainEntry struct {
	Status string // the two XY status characters
	Path   string // the record's path, verbatim
}

// RepoFacts is the typed local repository state used by read-only query owners.
// DefaultResolved makes the default-branch cells readable. When it is false,
// DefaultBranch is empty and Ahead/Behind are zero, because there is no branch to
// measure against — an unknown, not a measurement.
type RepoFacts struct {
	Branch, DefaultBranch string
	DefaultResolved       bool
	Dirty                 bool
	Ahead, Behind         int
	Changes               []PorcelainEntry
}

// DiffFacts is the additive facts path for bench diff. It leaves Facts unchanged for
// existing consumers. It expands untracked directories into the individual entries a
// coherent patch can actually show.
type DiffFacts struct {
	RepoFacts
	Head, DefaultTip, RecordedBase string
	Porcelain                      []byte
}

// LandedStateFact is the offline git verdict used by status. DirtyPaths describes the named
// checkout; commit and branch counts describe the repository.
type LandedStateFact struct {
	DirtyPaths        int
	UnpushedCommits   int
	UniqueBranches    int
	UniqueBranchNames []string
}

// LandedState derives checkout-local dirtiness and repository-wide commit and branch
// facts. The dirty count omits excludedDirtyPaths, but only for the named checkout.
func LandedState(root string, excludedDirtyPaths ...string) (LandedStateFact, error) {
	excluded := make(map[string]bool, len(excludedDirtyPaths))
	for _, path := range excludedDirtyPaths {
		excluded[path] = true
	}
	dirty := map[string]bool{}
	raw, err := Raw("-C", root, "status", "--porcelain=v1", "-z", "--no-renames")
	if err != nil {
		return LandedStateFact{}, fmt.Errorf("git status %s: %w", root, err)
	}
	for _, entry := range ParsePorcelainZ(raw) {
		if excluded[entry.Path] {
			continue
		}
		dirty[entry.Path] = true
	}
	def, ok := ResolvedDefault(root)
	if !ok {
		return LandedStateFact{}, errors.New("git repository has no resolvable default branch")
	}
	branches, err := LocalBranches(root)
	if err != nil {
		return LandedStateFact{}, fmt.Errorf("git local branches: %w", err)
	}
	commits := map[string]bool{}
	unique := map[string]bool{}
	for _, branch := range branches {
		if branch != def {
			landed, _, err := LandedInDefault(root, branch, def)
			if err != nil {
				return LandedStateFact{}, fmt.Errorf("git landedness %s: %w", branch, err)
			}
			if !landed {
				unique[branch] = true
			}
		}
		upstream, err := Output("-C", root, "for-each-ref", "--format=%(upstream:short)", "refs/heads/"+branch)
		if err != nil {
			return LandedStateFact{}, fmt.Errorf("git upstream %s: %w", branch, err)
		}
		if upstream == "" {
			continue
		}
		ahead, err := Output("-C", root, "rev-list", upstream+".."+branch)
		if err != nil {
			return LandedStateFact{}, fmt.Errorf("git ahead %s: %w", branch, err)
		}
		for _, commit := range strings.Split(ahead, "\n") {
			if commit != "" {
				commits[commit] = true
			}
		}
	}
	uniqueNames := make([]string, 0, len(unique))
	for branch := range unique {
		uniqueNames = append(uniqueNames, branch)
	}
	sort.Strings(uniqueNames)
	return LandedStateFact{DirtyPaths: len(dirty), UnpushedCommits: len(commits), UniqueBranches: len(uniqueNames), UniqueBranchNames: uniqueNames}, nil
}

// CheckedOutBranch names the branch HEAD points at, or the literal "HEAD" when detached.
// `rev-parse --abbrev-ref` fails outright on an unborn branch, so the symbolic ref settles
// that case. A repository with no commits still has a named branch, and losing the whole
// snapshot over a missing commit is the worse answer. The function is exported because
// every caller shares the probe chain, not the phrasing built from it. A caller that wants
// detachment reported as "no branch" tests the returned literal instead of running the two
// git queries again.
func CheckedOutBranch(root string) (string, error) {
	if name, err := Output("-C", root, "rev-parse", "--abbrev-ref", "HEAD"); err == nil && name != "" {
		return name, nil
	}
	return Output("-C", root, "symbolic-ref", "--quiet", "--short", "HEAD")
}

// Facts derives repository state without mutating the worktree or index.
func Facts(root string) (RepoFacts, error) {
	branch, err := CheckedOutBranch(root)
	if err != nil {
		return RepoFacts{}, err
	}
	raw, err := Raw("-C", root, "status", "--porcelain=v1", "-z", "--no-renames")
	if err != nil {
		return RepoFacts{}, err
	}
	f := RepoFacts{Branch: branch, Changes: ParsePorcelainZ(raw)}
	f.Dirty = len(f.Changes) > 0
	def, ok := ResolvedDefault(root)
	if !ok {
		// The code derives divergence from `rev-list <default>...HEAD`, which errors against
		// a branch that does not exist. The caller gets the unresolved state and the rest
		// of the snapshot instead of a failed read of the whole thing.
		return f, nil
	}
	f.DefaultBranch, f.DefaultResolved = def, true
	counts, err := Output("-C", root, "rev-list", "--left-right", "--count", f.DefaultBranch+"...HEAD")
	if err != nil {
		return RepoFacts{}, fmt.Errorf("git rev-list: %w", err)
	}
	if _, err := fmt.Sscanf(counts, "%d\t%d", &f.Behind, &f.Ahead); err != nil {
		return RepoFacts{}, fmt.Errorf("parse git divergence: %w", err)
	}
	return f, nil
}

// AllFilesFacts derives the diff-specific status facts with Git's all-files
// untracked policy. Existing Facts callers retain their collapsed-directory output.
func AllFilesFacts(root string) (DiffFacts, error) {
	branch, err := CheckedOutBranch(root)
	if err != nil {
		return DiffFacts{}, err
	}
	head, _ := Output("-C", root, "rev-parse", "HEAD")
	recordedBase := ""
	if branch != "HEAD" {
		recordedBase, _ = Output("-C", root, "config", "branch."+branch+".benchBase")
	}
	raw, changes, err := AllFilesStatus(root)
	if err != nil {
		return DiffFacts{}, err
	}
	f := DiffFacts{
		RepoFacts:    RepoFacts{Branch: branch, Changes: changes},
		Head:         head,
		RecordedBase: recordedBase,
		Porcelain:    raw,
	}
	f.Dirty = len(f.Changes) > 0
	def, ok := ResolvedDefault(root)
	if !ok {
		return f, nil
	}
	defaultTip, err := Output("-C", root, "rev-parse", def)
	if err != nil {
		return DiffFacts{}, fmt.Errorf("git default tip: %w", err)
	}
	f.DefaultBranch, f.DefaultTip, f.DefaultResolved = def, defaultTip, true
	counts, err := Output("-C", root, "rev-list", "--left-right", "--count", defaultTip+"..."+head)
	if err != nil {
		return DiffFacts{}, fmt.Errorf("git rev-list: %w", err)
	}
	if _, err := fmt.Sscanf(counts, "%d\t%d", &f.Behind, &f.Ahead); err != nil {
		return DiffFacts{}, fmt.Errorf("parse git divergence: %w", err)
	}
	return f, nil
}

// AllFilesStatus returns Git's raw all-files porcelain plus untracked special entries
// that Git omits from that stream. The supplemental walk only classifies non-regular,
// non-directory, non-symlink nodes. It asks Git whether each is tracked or ignored, and
// it never opens the node.
func AllFilesStatus(root string) ([]byte, []PorcelainEntry, error) {
	raw, err := Raw("-C", root, "status", "--porcelain=v1", "-z", "--no-renames", "--untracked-files=all")
	if err != nil {
		return nil, nil, err
	}
	changes := ParsePorcelainZ(raw)
	seen := make(map[string]bool, len(changes))
	for _, change := range changes {
		seen[change.Path] = true
	}
	err = filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.Name() == ".git" {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.IsDir() || entry.Type().IsRegular() || entry.Type()&os.ModeSymlink != 0 {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if seen[rel] || OK("-C", root, "ls-files", "--error-unmatch", "--", rel) || OK("-C", root, "check-ignore", "-q", "--", rel) {
			return nil
		}
		seen[rel] = true
		changes = append(changes, PorcelainEntry{Status: "??", Path: rel})
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	sort.Slice(changes, func(i, j int) bool { return changes[i].Path < changes[j].Path })
	return raw, changes, nil
}

// ParsePorcelainZ splits `git status --porcelain -z --no-renames` output into entries.
// The -z framing is NUL-delimited and never C-quotes. A path with spaces, glob
// characters, or a literal newline survives whole. This function is the one source of
// that framing knowledge for every caller — the shift staging diff and the commit
// block-check.
func ParsePorcelainZ(raw []byte) []PorcelainEntry {
	entries, _ := ParsePorcelainZStrict(raw)
	return entries
}

// ParsePorcelainZStrict parses and validates NUL-framed porcelain-v1 records. A rename
// or copy record carries a second, path-only NUL record. The parser returns that record
// with an empty Status, so callers can preserve the framing while filtering.
func ParsePorcelainZStrict(raw []byte) ([]PorcelainEntry, error) {
	var entries []PorcelainEntry
	for offset := 0; offset < len(raw); {
		end := bytes.IndexByte(raw[offset:], 0)
		if end < 0 {
			return nil, errors.New("malformed checkout status")
		}
		record := raw[offset : offset+end]
		offset += end + 1
		if len(record) < 4 || record[2] != ' ' {
			return nil, errors.New("malformed checkout status")
		}
		entries = append(entries, PorcelainEntry{Status: string(record[:2]), Path: string(record[3:])})
		status := record[:2]
		if status[0] == 'R' || status[0] == 'C' || status[1] == 'R' || status[1] == 'C' {
			end = bytes.IndexByte(raw[offset:], 0)
			if end < 1 {
				return nil, errors.New("malformed checkout status")
			}
			entries = append(entries, PorcelainEntry{Path: string(raw[offset : offset+end])})
			offset += end + 1
		}
	}
	return entries, nil
}
