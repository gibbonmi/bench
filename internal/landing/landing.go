// Package landing owns exact prospective Git-tree composition and publication.
package landing

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/gate/authorization"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/spec"
)

// Request is the complete, immutable input to one prospective landing.
type Request struct {
	Root, Destination, Expected, Message string
	Paths                                []string
	Spec                                 string
	Stdout, Stderr                       io.Writer
}

// Result identifies a successfully published commit and its authorized tree.
type Result struct{ Base, Commit, Tree string }

// ReviewedRequest names the immutable, reviewer-approved source pair and the
// destination snapshot it is composed onto. Checkout fingerprints guard state;
// every published byte still comes from Git objects.
type ReviewedRequest struct {
	Root, Destination, DestinationBase string
	Source, SourceTip, ReviewBase      string
	SourceWorktree                     string
	SourceFingerprint                  string
	DestinationFingerprint             string
	SpecPath                           string
	SpecBytes                          []byte
	SpecMode                           os.FileMode
	Message                            string
	Stdout, Stderr                     io.Writer
}

// ReviewedResult is the immutable publication receipt needed by the lifecycle owner.
type ReviewedResult struct {
	SourceBase, SourceTip, DestinationBase, Commit, Tree string
}

// CompositionRequest identifies two immutable commits to merge without checkout state.
type CompositionRequest struct {
	Root, Destination, Source, ReviewBase string
}

// CompositionResult is either a prospective tree or one bounded conflict kind.
type CompositionResult struct {
	Base, Tree string
	Conflict   Conflict
}

// Conflict describes why Git could not produce one prospective tree.
type Conflict struct{ Kind string }

// Compose performs Git's three-way tree merge using the repository's real merge base.
// ReviewBase is metadata only and is never used as the merge base.
func (o Owner) Compose(r CompositionRequest) (CompositionResult, error) {
	if r.Root == "" || r.Destination == "" || r.Source == "" {
		return CompositionResult{}, errors.New("composition request is incomplete")
	}
	destination, err := compositionCommit(r.Root, r.Destination, "destination")
	if err != nil {
		return CompositionResult{}, err
	}
	source, err := compositionCommit(r.Root, r.Source, "source")
	if err != nil {
		return CompositionResult{}, err
	}
	base, err := output(r.Root, "merge-base", destination, source)
	if err != nil {
		return CompositionResult{}, fmt.Errorf("find merge base: %w", err)
	}
	out, err := mergeTree(r.Root, destination, source)
	if err == nil {
		tree, err := mergeTreeResult(out)
		if err != nil {
			return CompositionResult{}, err
		}
		return CompositionResult{Base: base, Tree: tree}, nil
	}
	kind, parseErr := conflictKind(out)
	if parseErr != nil {
		return CompositionResult{}, parseErr
	}
	return CompositionResult{Base: base, Conflict: Conflict{Kind: kind}}, nil
}

func compositionCommit(root, value, role string) (string, error) {
	commit, err := output(root, "rev-parse", "--verify", value+"^{commit}")
	if err != nil {
		return "", fmt.Errorf("composition %s is not a commit", role)
	}
	return commit, nil
}

func mergeTree(root, destination, source string) (string, error) {
	// merge-tree finds the merge base itself; this avoids checkout and index state.
	return outputCombined(root, "merge-tree", "--write-tree", "-z", destination, source)
}

func mergeTreeResult(output string) (string, error) {
	tree, _, found := strings.Cut(output, "\x00")
	if !found || len(tree) != 40 {
		return "", errors.New("merge-tree returned no tree")
	}
	return tree, nil
}

func conflictKind(output string) (string, error) {
	parts := bytes.Split([]byte(output), []byte{0})
	separator := -1
	for i, part := range parts {
		if len(part) == 0 {
			separator = i
			break
		}
	}
	if separator < 1 {
		return "", errors.New("merge-tree returned no conflict records")
	}
	modes := make([]string, 0, separator-1)
	for _, record := range parts[1:separator] {
		fields := strings.Fields(string(record))
		if len(fields) < 3 || fields[2][0] < '1' || fields[2][0] > '3' {
			return "", errors.New("merge-tree returned malformed conflict record")
		}
		modes = append(modes, fields[0])
	}
	for _, record := range parts[separator+1:] {
		kind := string(record)
		switch kind {
		case "CONFLICT (modify/delete)":
			return "modify/delete", nil
		case "CONFLICT (rename/rename)":
			return "rename/rename", nil
		case "CONFLICT (directory/file)", "CONFLICT (file/directory)":
			return "file/directory", nil
		case "CONFLICT (distinct modes)":
			kind := contentConflictKind(modes)
			if kind == "textual" {
				return "mode", nil
			}
			return kind, nil
		case "CONFLICT (contents)":
			return contentConflictKind(modes), nil
		}
	}
	return "", errors.New("merge-tree returned an unrecognized conflict kind")
}

func contentConflictKind(modes []string) string {
	for _, mode := range modes {
		if mode == "160000" {
			return "gitlink"
		}
		if mode == "120000" {
			return "symlink"
		}
	}
	for _, mode := range modes {
		for _, other := range modes {
			if mode != other {
				return "mode"
			}
		}
	}
	return "textual"
}

type composedSnapshot struct {
	tree            string
	specPath        string
	specPermissions os.FileMode
	// closePath is the repository-relative tickets-only folder this landing consumes,
	// empty on every other landing. It is composed out of the published tree and
	// removed from the checkout only after that tree is authorized and published.
	closePath string
}

// Owner composes only Request.Paths from Request.Expected, authorizes that tree, then
// publishes it by expected-old update. The function fields are narrow operational seams
// for deterministic fault coverage; New supplies the real owners.
type Owner struct {
	authorize func(context.Context, string, string, io.Writer, io.Writer) authorization.Result
	updateRef func(string, string, string, string) error
	reconcile func(Request, []string, composedSnapshot) error
}

// New returns the production landing owner.
func New() Owner {
	return Owner{authorize: authorization.AuthorizeWithWriters, updateRef: updateRef, reconcile: reconcile}
}

// Land publishes a commit only when the exact composed tree receives green authorization.
func (o Owner) Land(ctx context.Context, r Request) (Result, error) {
	if err := validRequest(r); err != nil {
		return Result{}, err
	}
	paths, err := attributedPaths(r.Root, r.Expected, r.Paths)
	if err != nil {
		return Result{}, err
	}
	// --spec has two effects on one green landing. A folder carrying a spec.md takes
	// the staged->implemented flip; a tickets-only folder is deleted instead, in the
	// same commit. The close path stays out of paths: it is composed and reconciled by
	// removal, not by attribution.
	closePath := ""
	if r.Spec != "" {
		if TicketsOnlyFolder(r.Root, r.Spec) {
			closePath = specsDir + "/" + r.Spec
		} else {
			resolved, err := spec.CheckStaged(r.Root, r.Spec)
			if err != nil {
				return Result{}, err
			}
			rel, err := repositoryPath(r.Root, resolved)
			if err != nil {
				return Result{}, err
			}
			paths = append(paths, rel)
		}
	}
	paths = unique(paths)
	snapshot, err := compose(r, paths, closePath)
	if err != nil {
		return Result{}, err
	}
	tree := snapshot.tree
	baseTree, err := output(r.Root, "rev-parse", r.Expected+"^{tree}")
	if err != nil {
		return Result{}, fmt.Errorf("read expected base tree: %w", err)
	}
	if tree == baseTree {
		return Result{}, errors.New("nothing to commit")
	}
	if got := o.authorize(ctx, r.Root, tree, r.Stdout, r.Stderr); got.Kind != authorization.Green {
		return Result{}, fmt.Errorf("prospective authorization refused: %s", got.Kind)
	}
	commit, err := output(r.Root, "commit-tree", tree, "-p", r.Expected, "-m", r.Message)
	if err != nil {
		return Result{}, fmt.Errorf("create landing commit: %w", err)
	}
	if err := o.updateRef(r.Root, r.Destination, commit, r.Expected); err != nil {
		return Result{}, destinationUpdateFailure(r.Root, r.Destination, r.Expected, err)
	}
	if err := o.reconcile(r, paths, snapshot); err != nil {
		return Result{Base: r.Expected, Commit: commit, Tree: tree}, fmt.Errorf("landed-but-checkout-incomplete: %w", err)
	}
	return Result{Base: r.Expected, Commit: commit, Tree: tree}, nil
}

// LandReviewed composes an exact reviewed source, applies its staged-spec
// transition only to that prospective tree, and publishes a two-parent commit.
// The worktree lifecycle owns authentication, marker advancement, reconciliation,
// and release around this irreversible operation.
func (o Owner) LandReviewed(ctx context.Context, r ReviewedRequest) (ReviewedResult, error) {
	if r.Root == "" || r.Destination == "" || r.DestinationBase == "" || r.Source == "" || r.SourceTip == "" || r.ReviewBase == "" || r.SourceWorktree == "" || r.SourceFingerprint == "" || r.DestinationFingerprint == "" || r.SpecPath == "" || strings.TrimSpace(r.Message) == "" {
		return ReviewedResult{}, errors.New("reviewed landing request is incomplete")
	}
	destination, err := compositionCommit(r.Root, r.DestinationBase, "destination")
	if err != nil || destination != r.DestinationBase {
		return ReviewedResult{}, errors.New("destination base is not an exact commit")
	}
	source, err := compositionCommit(r.Root, r.SourceTip, "source")
	if err != nil || source != r.SourceTip {
		return ReviewedResult{}, errors.New("source tip is not an exact commit")
	}
	branchTip, err := output(r.Root, "rev-parse", "--verify", r.Source+"^{commit}")
	if err != nil || branchTip != source {
		return ReviewedResult{}, errors.New("reviewed source tip moved")
	}
	if _, kind, _ := diff.ResolveSourceRange(r.Root, r.ReviewBase, source); kind != "" {
		return ReviewedResult{}, errors.New("reviewed source base is invalid")
	}
	if err := stagedSpecMatches(r.Root, destination, source, r.SpecPath, r.SpecBytes); err != nil {
		return ReviewedResult{}, err
	}
	composition, err := o.Compose(CompositionRequest{Root: r.Root, Destination: destination, Source: source, ReviewBase: r.ReviewBase})
	if err != nil {
		return ReviewedResult{}, err
	}
	if composition.Conflict.Kind != "" {
		return ReviewedResult{}, fmt.Errorf("composition conflict: %s", composition.Conflict.Kind)
	}
	implemented, err := spec.Implemented(r.SpecBytes)
	if err != nil {
		return ReviewedResult{}, err
	}
	tree, err := replaceTreeFile(r.Root, composition.Tree, r.SpecPath, implemented, r.SpecMode)
	if err != nil {
		return ReviewedResult{}, fmt.Errorf("transition staged spec: %w", err)
	}
	if got := o.authorize(ctx, r.Root, tree, r.Stdout, r.Stderr); got.Kind != authorization.Green {
		return ReviewedResult{}, fmt.Errorf("prospective authorization refused: %s", got.Kind)
	}
	// Recheck the two moving identities after the gate and before creating an
	// otherwise unreachable object. Tree equality is insufficient: review binds a commit.
	if branchTip, err = output(r.Root, "rev-parse", "--verify", r.Source+"^{commit}"); err != nil || branchTip != source {
		return ReviewedResult{}, errors.New("reviewed source tip moved")
	}
	if fingerprint, fingerprintErr := CheckoutFingerprint(r.SourceWorktree); fingerprintErr != nil || fingerprint != r.SourceFingerprint {
		return ReviewedResult{}, errors.New("reviewed source checkout changed")
	}
	if fingerprint, fingerprintErr := CheckoutFingerprint(r.Root); fingerprintErr != nil || fingerprint != r.DestinationFingerprint {
		return ReviewedResult{}, errors.New("landing destination checkout changed; rerun the landing to recompose onto the moved destination")
	}
	commit, err := output(r.Root, "commit-tree", tree, "-p", destination, "-p", source, "-m", r.Message)
	if err != nil {
		return ReviewedResult{}, fmt.Errorf("create landing commit: %w", err)
	}
	if err := o.updateRef(r.Root, r.Destination, commit, destination); err != nil {
		return ReviewedResult{}, destinationUpdateFailure(r.Root, r.Destination, destination, err)
	}
	return ReviewedResult{SourceBase: r.ReviewBase, SourceTip: source, DestinationBase: destination, Commit: commit, Tree: tree}, nil
}

// CheckoutFingerprint binds the attached branch, commit, index, worktree,
// untracked set, ignored set, and nested-repository state observed at a checkout.
func CheckoutFingerprint(root string) (string, error) {
	branch, err := output(root, "symbolic-ref", "--quiet", "HEAD")
	if err != nil {
		return "", err
	}
	head, err := output(root, "rev-parse", "HEAD^{commit}")
	if err != nil {
		return "", err
	}
	status, err := benchgit.Raw("-C", root, "status", "--porcelain=v1", "-z", "--untracked-files=all", "--ignored")
	if err != nil {
		return "", err
	}
	status, err = fingerprintStatus(status)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(bytes.Join([][]byte{[]byte(branch), []byte(head), status}, []byte{0}))
	return fmt.Sprintf("%x", sum), nil
}

// RuntimeIgnoredPath reports whether path is a safe spelling within Bench's ignored runtime log root.
func RuntimeIgnoredPath(path string) bool {
	for _, r := range path {
		if r < 0x20 || r == 0x7f {
			return false
		}
	}
	if path == ".logs/" {
		return true
	}
	native := filepath.FromSlash(path)
	if path == "" || filepath.IsAbs(native) || filepath.Clean(native) != native || filepath.ToSlash(native) != path {
		return false
	}
	return path == ".logs" || strings.HasPrefix(path, ".logs/")
}

func fingerprintStatus(raw []byte) ([]byte, error) {
	entries, err := benchgit.ParsePorcelainZStrict(raw)
	if err != nil {
		return nil, err
	}
	var filtered bytes.Buffer
	for _, entry := range entries {
		if entry.Status == "" {
			filtered.WriteString(entry.Path)
			filtered.WriteByte(0)
			continue
		}
		if entry.Status == "!!" && RuntimeIgnoredPath(entry.Path) {
			continue
		}
		filtered.WriteString(entry.Status)
		filtered.WriteByte(' ')
		filtered.WriteString(entry.Path)
		filtered.WriteByte(0)
	}
	return filtered.Bytes(), nil
}

func stagedSpecMatches(root, destination, source, path string, want []byte) error {
	for _, commit := range []string{destination, source} {
		got, err := benchgit.Raw("-C", root, "show", commit+":"+path)
		if err != nil || !bytes.Equal(got, want) {
			return errors.New("source and destination do not carry identical staged spec bytes")
		}
	}
	return nil
}

func replaceTreeFile(root, baseTree, path string, content []byte, mode os.FileMode) (string, error) {
	dir, err := os.MkdirTemp("", "bench-reviewed-landing-index-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)
	idx := filepath.Join(dir, "index")
	if err := indexRun(root, idx, "read-tree", baseTree); err != nil {
		return "", err
	}
	blob, err := outputInput(root, content, "hash-object", "-w", "--stdin")
	if err != nil {
		return "", err
	}
	if err := indexRun(root, idx, "update-index", "--add", "--cacheinfo", gitRegularFileMode(mode)+","+blob+","+path); err != nil {
		return "", err
	}
	return indexOutput(root, idx, "write-tree")
}

func validRequest(r Request) error {
	if r.Root == "" || r.Destination == "" || r.Expected == "" || strings.TrimSpace(r.Message) == "" {
		return errors.New("landing request is incomplete")
	}
	if _, err := output(r.Root, "rev-parse", "--verify", r.Expected+"^{commit}"); err != nil {
		return errors.New("expected base is not a commit")
	}
	return nil
}

func attributedPaths(root, expected string, raw []string) ([]string, error) {
	if len(raw) == 0 {
		return nil, errors.New("at least one path is required")
	}
	paths := make([]string, 0, len(raw))
	for _, p := range raw {
		rel, err := repositoryPath(root, p)
		if err != nil {
			return nil, err
		}
		if rel == "." || rel == "" {
			return nil, errors.New("repository root is not an attributed path")
		}
		if err := safePath(root, expected, rel); err != nil {
			return nil, err
		}
		paths = append(paths, rel)
	}
	return unique(paths), nil
}

func repositoryPath(root, path string) (string, error) {
	abs := path
	if !filepath.IsAbs(abs) {
		abs = filepath.Join(root, path)
	}
	abs, err := filepath.Abs(abs)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(root, abs)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path %q escapes repository", path)
	}
	return filepath.ToSlash(rel), nil
}

func safePath(root, expected, rel string) error {
	for current := filepath.Join(root, filepath.FromSlash(rel)); ; current = filepath.Dir(current) {
		info, err := os.Lstat(current)
		if err == nil && info.Mode().Perm()&0o444 == 0 {
			return fmt.Errorf("unreadable path %q is not attributable", rel)
		}
		if err == nil && special(info.Mode()) {
			return fmt.Errorf("special file %q is not attributable", rel)
		}
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("inspect %q: %w", rel, err)
		}
		if current == root {
			break
		}
	}
	path := filepath.Join(root, filepath.FromSlash(rel))
	info, err := os.Lstat(path)
	if err == nil && info.IsDir() && !gitlinkAt(root, expected, rel) {
		err = filepath.WalkDir(path, func(current string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if entry.Type()&os.ModeSymlink != 0 {
				return nil
			}
			entryInfo, statErr := entry.Info()
			if statErr != nil {
				return statErr
			}
			if entryInfo.Mode().Perm()&0o444 == 0 || special(entryInfo.Mode()) {
				return errors.New("special or unreadable descendant")
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("inspect attributed directory %q: %w", rel, err)
		}
	}
	return nil
}

func special(mode os.FileMode) bool {
	return mode&os.ModeType != 0 && mode&os.ModeSymlink == 0 && !mode.IsDir()
}

func gitlinkAt(root, expected, path string) bool {
	for _, args := range [][]string{
		{"ls-tree", expected, "--", path},
		{"ls-files", "--stage", "--", path},
	} {
		line, err := output(root, args...)
		if err == nil && strings.HasPrefix(line, "160000 ") {
			return true
		}
	}
	return false
}

func compose(r Request, paths []string, closePath string) (composedSnapshot, error) {
	dir, err := os.MkdirTemp("", "bench-landing-index-")
	if err != nil {
		return composedSnapshot{}, err
	}
	defer os.RemoveAll(dir)
	idx := filepath.Join(dir, "index")
	if err := indexRun(r.Root, idx, "read-tree", r.Expected); err != nil {
		return composedSnapshot{}, err
	}
	for _, path := range paths {
		if err := indexRun(r.Root, idx, "add", "-A", "--", ":(literal)"+path); err != nil {
			if !trackedAt(r.Root, r.Expected, path) {
				return composedSnapshot{}, fmt.Errorf("named path %q not found in worktree, index, or expected base", path)
			}
		}
	}
	snapshot := composedSnapshot{closePath: closePath}
	if closePath != "" {
		if err := removeIndexTree(r.Root, idx, closePath); err != nil {
			return composedSnapshot{}, err
		}
	} else if r.Spec != "" {
		resolved, content, fileMode, err := transitionedSpec(r.Root, r.Spec)
		if err != nil {
			return composedSnapshot{}, err
		}
		rel, _ := repositoryPath(r.Root, resolved)
		blob, err := outputInput(r.Root, content, "hash-object", "-w", "--stdin")
		if err != nil {
			return composedSnapshot{}, err
		}
		if err := indexRun(r.Root, idx, "update-index", "--add", "--cacheinfo", gitRegularFileMode(fileMode)+","+blob+","+rel); err != nil {
			return composedSnapshot{}, err
		}
		snapshot.specPath = resolved
		snapshot.specPermissions = fileMode
	}
	snapshot.tree, err = indexOutput(r.Root, idx, "write-tree")
	return snapshot, err
}

func gitRegularFileMode(mode os.FileMode) string {
	if mode.Perm()&0o111 != 0 {
		return "100755"
	}
	return "100644"
}

func transitionedSpec(root, slug string) (string, []byte, os.FileMode, error) {
	resolved, err := spec.CheckStaged(root, slug)
	if err != nil {
		return "", nil, 0, err
	}
	content, err := os.ReadFile(resolved)
	if err != nil {
		return "", nil, 0, err
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return "", nil, 0, err
	}
	content, err = spec.Implemented(content)
	if err != nil {
		return "", nil, 0, err
	}
	return resolved, content, info.Mode().Perm(), nil
}

func reconcile(r Request, paths []string, snapshot composedSnapshot) error {
	for _, path := range paths {
		literal := ":(literal)" + path
		if err := run(r.Root, "restore", "--source="+snapshot.tree, "--staged", "--worktree", "--", literal); err != nil {
			if resetErr := run(r.Root, "reset", "-q", r.Expected, "--", literal); resetErr != nil {
				return err
			}
			if retryErr := run(r.Root, "restore", "--source="+snapshot.tree, "--staged", "--worktree", "--", literal); retryErr != nil {
				return retryErr
			}
		}
		if err := run(r.Root, "clean", "-f", "-d", "--", literal); err != nil {
			return err
		}
		if err := run(r.Root, "diff", "--quiet", "--cached", snapshot.tree, "--", literal); err != nil {
			return err
		}
		if err := run(r.Root, "diff", "--quiet", "--ignore-submodules=dirty", snapshot.tree, "--", literal); err != nil {
			return err
		}
		untracked, err := output(r.Root, "ls-files", "--others", "--exclude-standard", "--", literal)
		if err != nil {
			return err
		}
		if untracked != "" {
			return fmt.Errorf("named path %q still has untracked content", path)
		}
	}
	if snapshot.closePath != "" {
		if err := os.RemoveAll(filepath.Join(r.Root, filepath.FromSlash(snapshot.closePath))); err != nil {
			return fmt.Errorf("remove closed spec folder %q: %w", snapshot.closePath, err)
		}
	}
	if snapshot.specPath != "" {
		return os.Chmod(snapshot.specPath, snapshot.specPermissions)
	}
	return nil
}

// removeIndexTree drops every entry beneath rel from the prospective index, so the
// published tree carries the deletion rather than the checkout carrying it afterwards.
// The pathspec is literal and the removals name exact index paths, so a folder name
// holding a space or a glob character resolves to itself.
func removeIndexTree(root, idx, rel string) error {
	listed, err := indexOutputRaw(root, idx, "ls-files", "-z", "--cached", "--", ":(literal)"+rel)
	if err != nil {
		return fmt.Errorf("list tracked entries under %q: %w", rel, err)
	}
	for _, path := range strings.Split(string(listed), "\x00") {
		if path == "" {
			continue
		}
		if err := indexRun(root, idx, "update-index", "--force-remove", "--", path); err != nil {
			return fmt.Errorf("remove %q from prospective index: %w", path, err)
		}
	}
	return nil
}

func trackedAt(root, base, path string) bool {
	return run(root, "cat-file", "-e", base+":"+path) == nil
}
func unique(values []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, v := range values {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	sort.Strings(out)
	return out
}
func updateRef(root, ref, new, old string) error { return run(root, "update-ref", ref, new, old) }
func destinationUpdateFailure(root, ref, expected string, updateErr error) error {
	actual, err := output(root, "rev-parse", "--verify", ref+"^{commit}")
	if err != nil {
		return fmt.Errorf("read destination after failed ref update: %w", err)
	}
	if actual != expected {
		return fmt.Errorf("destination compare-and-swap refused; rerun the landing to recompose onto the moved destination: %w", updateErr)
	}
	return updateErr
}
func output(root string, args ...string) (string, error) {
	return benchgit.Output(append([]string{"-C", root}, args...)...)
}
func outputCombined(root string, args ...string) (string, error) {
	c := exec.Command("git", append([]string{"-C", root}, args...)...)
	b, err := c.CombinedOutput()
	return strings.TrimSpace(string(b)), err
}
func outputInput(root string, input []byte, args ...string) (string, error) {
	c := exec.Command("git", append([]string{"-C", root}, args...)...)
	c.Stdin = strings.NewReader(string(input))
	b, err := c.Output()
	return strings.TrimSpace(string(b)), err
}
func run(root string, args ...string) error {
	return exec.Command("git", append([]string{"-C", root}, args...)...).Run()
}
func indexRun(root, idx string, args ...string) error {
	c := exec.Command("git", append([]string{"-C", root}, args...)...)
	c.Env = append(os.Environ(), "GIT_INDEX_FILE="+idx)
	return c.Run()
}
func indexOutputRaw(root, idx string, args ...string) ([]byte, error) {
	c := exec.Command("git", append([]string{"-C", root}, args...)...)
	c.Env = append(os.Environ(), "GIT_INDEX_FILE="+idx)
	return c.Output()
}
func indexOutput(root, idx string, args ...string) (string, error) {
	c := exec.Command("git", append([]string{"-C", root}, args...)...)
	c.Env = append(os.Environ(), "GIT_INDEX_FILE="+idx)
	b, err := c.Output()
	return strings.TrimSpace(string(b)), err
}
