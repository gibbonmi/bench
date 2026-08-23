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
	// ClosePath is the repository-relative tickets-only folder this landing closes,
	// empty on every other landing. It never coexists with SpecPath: one --spec names
	// either a staged spec.md to transition or a tickets-only folder to close.
	ClosePath      string
	Message        string
	Stdout, Stderr io.Writer
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
// Resolved lists every capture path the composition policy settled, as
// "<path>:<side>", so the landing can disclose what the merge did not decide.
type CompositionResult struct {
	Base, Tree string
	Conflict   Conflict
	Resolved   []string
}

// Conflict describes why Git could not produce one prospective tree, and names
// every path it could not merge.
type Conflict struct {
	Kind  string
	Paths []string
}

// ConflictError is the refusal a conflicted reviewed landing returns. Its message is
// the bounded kind; the paths ride typed so the caller can render them.
type ConflictError struct{ Conflict }

func (e ConflictError) Error() string { return "composition conflict: " + e.Kind }

// stageRecord is one conflicted-file line of `merge-tree -z`: the mode and object of
// one stage (1 base, 2 destination, 3 source) at one path.
type stageRecord struct {
	mode, oid, path string
	stage           int
}

// unionStage is the CaptureSide stage of a path the rule table settles by union
// rather than by taking one side; no merge-tree stage carries that number.
const unionStage = 0

// CaptureSide is the rule table for a conflicted phase-owned path: it names the verb
// that settles the path and, for a take-a-side verb, the merge-tree stage that wins.
// The phase handoff is the source session's closing state, so the source wins it. The
// two append-only journals compose as the union of both sides, so no appended entry is
// lost. Every other capture file is the destination's running ledger, so the
// destination wins. A path outside the table has no rule, and the conflict stays a
// refusal.
func CaptureSide(path string) (stage int, side string, ok bool) {
	switch path {
	case "capture/session-handoff.md":
		return 3, "source", true
	case "capture/learnings.md", "capture/IDEAS.md":
		return unionStage, "union", true
	}
	if strings.HasPrefix(path, "capture/") {
		return 2, "destination", true
	}
	return 0, "", false
}

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
	conflict, records, parseErr := parseConflict(out)
	if parseErr != nil {
		return CompositionResult{}, parseErr
	}
	tree, resolved, ok, err := resolveCaptureConflict(r.Root, out, records)
	if err != nil {
		return CompositionResult{}, err
	}
	if ok {
		return CompositionResult{Base: base, Tree: tree, Resolved: resolved}, nil
	}
	return CompositionResult{Base: base, Conflict: conflict}, nil
}

// resolveCaptureConflict settles a conflict whose every path is a regular file the
// rule table names, by applying that path's verb. The merge's own written tree carries
// Git's conflict markers at those paths; each is replaced by the settled object, or
// removed when the settled verb takes a side that deleted the file. A conflict touching
// any path the table does not name, any non-regular stage, any mode disagreement, or
// any union content Git cannot merge as text is left for the caller to refuse.
func resolveCaptureConflict(root, mergeOutput string, records []stageRecord) (string, []string, bool, error) {
	stages := map[string]map[int]stageRecord{}
	var order []string
	for _, record := range records {
		if _, _, ok := CaptureSide(record.path); !ok {
			return "", nil, false, nil
		}
		if record.mode != "100644" && record.mode != "100755" {
			return "", nil, false, nil
		}
		if _, seen := stages[record.path]; !seen {
			stages[record.path] = map[int]stageRecord{}
			order = append(order, record.path)
		}
		stages[record.path][record.stage] = record
	}
	if len(order) == 0 {
		return "", nil, false, nil
	}
	// settled holds the object each path publishes; a nil entry publishes a removal.
	settled := map[string]*stageRecord{}
	for _, path := range order {
		// Two sides that disagree on the file mode leave no settled mode to publish,
		// so the conflict stays a refusal under every verb.
		destination, hasDestination := stages[path][2]
		source, hasSource := stages[path][3]
		if hasDestination && hasSource && destination.mode != source.mode {
			return "", nil, false, nil
		}
		if stage, _, _ := CaptureSide(path); stage != unionStage {
			if record, ok := stages[path][stage]; ok {
				settled[path] = &record
			} else {
				settled[path] = nil
			}
			continue
		}
		record, ok, err := unionStages(root, path, stages[path])
		if err != nil {
			return "", nil, false, err
		}
		if !ok {
			return "", nil, false, nil
		}
		settled[path] = record
	}
	baseTree, err := mergeTreeResult(mergeOutput)
	if err != nil {
		return "", nil, false, err
	}
	tree, err := editTree(root, baseTree, func(idx string) error {
		for _, path := range order {
			record := settled[path]
			if record == nil {
				if err := indexRun(root, idx, "update-index", "--force-remove", "--", path); err != nil {
					return fmt.Errorf("resolve %q: %w", path, err)
				}
				continue
			}
			if err := indexRun(root, idx, "update-index", "--add", "--cacheinfo", record.mode+","+record.oid+","+path); err != nil {
				return fmt.Errorf("resolve %q: %w", path, err)
			}
		}
		return nil
	})
	if err != nil {
		return "", nil, false, err
	}
	resolved := make([]string, 0, len(order))
	for _, path := range order {
		_, side, _ := CaptureSide(path)
		resolved = append(resolved, path+":"+side)
	}
	return tree, resolved, true, nil
}

// unionStages composes one union path from its merge-tree stages. With both sides
// present, the result is Git's own three-way union merge over the three stage blobs,
// with an absent merge base standing in as empty content. With one side absent, the
// present side's blob is the result, so a deletion cannot erase the other side's
// entries. A false ok is a refusal: Git could not merge the content as text.
func unionStages(root, path string, stages map[int]stageRecord) (*stageRecord, bool, error) {
	destination, hasDestination := stages[2]
	source, hasSource := stages[3]
	switch {
	case !hasDestination && !hasSource:
		return nil, false, nil
	case !hasSource:
		return &destination, true, nil
	case !hasDestination:
		return &source, true, nil
	}
	dir, err := os.MkdirTemp("", "bench-landing-union-")
	if err != nil {
		return nil, false, err
	}
	defer os.RemoveAll(dir)
	names := map[int]string{}
	for _, stage := range []int{1, 2, 3} {
		file := filepath.Join(dir, fmt.Sprintf("stage%d", stage))
		content := []byte(nil)
		if record, ok := stages[stage]; ok {
			if content, err = blobContent(root, record.oid); err != nil {
				return nil, false, err
			}
		}
		if err := os.WriteFile(file, content, 0o600); err != nil {
			return nil, false, err
		}
		names[stage] = file
	}
	merged, err := outputRaw(root, "merge-file", "-p", "--union", names[2], names[1], names[3])
	if err != nil {
		// merge-file refuses content it cannot read as text, so the union has no
		// result and the whole conflict returns to the caller as a refusal.
		return nil, false, nil
	}
	oid, err := hashBlob(root, merged)
	if err != nil {
		return nil, false, fmt.Errorf("write union of %q: %w", path, err)
	}
	return &stageRecord{mode: source.mode, oid: oid, path: path, stage: 3}, true, nil
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

// parseConflict reads `merge-tree --write-tree -z` conflict output: the written tree,
// one stage record per conflicted file (`<mode> <object> <stage>\t<path>`), a NUL
// separator, then the informational messages that carry the conflict kind. It
// returns the bounded kind with every conflicted path, and the stage records.
func parseConflict(output string) (Conflict, []stageRecord, error) {
	parts := bytes.Split([]byte(output), []byte{0})
	separator := -1
	for i, part := range parts {
		if len(part) == 0 {
			separator = i
			break
		}
	}
	if separator < 1 {
		return Conflict{}, nil, errors.New("merge-tree returned no conflict records")
	}
	records := make([]stageRecord, 0, separator-1)
	modes := make([]string, 0, separator-1)
	var paths []string
	seen := map[string]bool{}
	for _, raw := range parts[1:separator] {
		header, path, found := strings.Cut(string(raw), "\t")
		fields := strings.Fields(header)
		if !found || path == "" || len(fields) != 3 || len(fields[2]) != 1 || fields[2][0] < '1' || fields[2][0] > '3' {
			return Conflict{}, nil, errors.New("merge-tree returned malformed conflict record")
		}
		records = append(records, stageRecord{mode: fields[0], oid: fields[1], stage: int(fields[2][0] - '0'), path: path})
		modes = append(modes, fields[0])
		if !seen[path] {
			seen[path] = true
			paths = append(paths, path)
		}
	}
	for _, record := range parts[separator+1:] {
		kind := ""
		switch string(record) {
		case "CONFLICT (modify/delete)":
			kind = "modify/delete"
		case "CONFLICT (rename/rename)":
			kind = "rename/rename"
		case "CONFLICT (directory/file)", "CONFLICT (file/directory)":
			kind = "file/directory"
		case "CONFLICT (distinct modes)":
			kind = contentConflictKind(modes)
			if kind == "textual" {
				kind = "mode"
			}
		case "CONFLICT (contents)":
			kind = contentConflictKind(modes)
		}
		if kind != "" {
			return Conflict{Kind: kind, Paths: paths}, records, nil
		}
	}
	return Conflict{}, nil, errors.New("merge-tree returned an unrecognized conflict kind")
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
			closePath = ClosedFolderPath(r.Spec)
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
	if r.Root == "" || r.Destination == "" || r.DestinationBase == "" || r.Source == "" || r.SourceTip == "" || r.ReviewBase == "" || r.SourceWorktree == "" || r.SourceFingerprint == "" || r.DestinationFingerprint == "" || strings.TrimSpace(r.Message) == "" {
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
	// An empty SpecPath is the spec-less landing: it has no staged spec to prove, to
	// neutralize, or to transition, so it composes against the destination itself and
	// publishes the composition unchanged.
	composed := destination
	if r.SpecPath != "" {
		if err := stagedSpecMatches(r.Root, source, r.SpecPath, r.SpecBytes); err != nil {
			return ReviewedResult{}, err
		}
		if composed, err = specNeutralizedDestination(r.Root, destination, r.SpecPath, r.SpecBytes, r.SpecMode); err != nil {
			return ReviewedResult{}, err
		}
	}
	composition, err := o.Compose(CompositionRequest{Root: r.Root, Destination: composed, Source: source, ReviewBase: r.ReviewBase})
	if err != nil {
		return ReviewedResult{}, err
	}
	if composition.Conflict.Kind != "" {
		return ReviewedResult{}, ConflictError{composition.Conflict}
	}
	if len(composition.Resolved) > 0 && r.Stderr != nil {
		fmt.Fprintf(r.Stderr, "landing composition{resolved=%s}\n", strings.Join(composition.Resolved, ","))
	}
	tree := composition.Tree
	if r.SpecPath != "" {
		implemented, err := spec.Implemented(r.SpecBytes)
		if err != nil {
			return ReviewedResult{}, err
		}
		if tree, err = replaceTreeFile(r.Root, composition.Tree, r.SpecPath, implemented, r.SpecMode); err != nil {
			return ReviewedResult{}, fmt.Errorf("transition staged spec: %w", err)
		}
	}
	// The close consumes the tickets-only folder from the published tree with the same
	// index removal the commit path composes. A folder the destination already removed
	// lists no entries, so the removal writes the composed tree back unchanged.
	if r.ClosePath != "" {
		if tree, err = removeTreeFolder(r.Root, tree, r.ClosePath); err != nil {
			return ReviewedResult{}, fmt.Errorf("close tickets-only folder: %w", err)
		}
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

// RuntimeIgnoredPath reports whether path is a safe spelling within Bench's ignored
// runtime log root.
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

// stagedSpecMatches proves provenance, not agreement: the bytes the landing will
// transition must be the reviewed source tip's committed spec. The function never
// reads the destination's copy for comparison; the source's bytes win. A stale,
// amended, or absent destination spec is not a landing question.
func stagedSpecMatches(root, source, path string, want []byte) error {
	got, err := benchgit.Raw("-C", root, "show", source+":"+path)
	if err != nil || !bytes.Equal(got, want) {
		return errors.New("staged spec bytes are not the reviewed source tip's committed spec")
	}
	return nil
}

// specNeutralizedDestination returns a commit to compose against whose tree already
// carries the source's spec bytes. The spec path is then a one-sided change no merge
// can conflict on. Its parent is the real destination, so the merge base stays
// unchanged. The returned commit is a composition input only; it never becomes a
// published parent.
func specNeutralizedDestination(root, destination, path string, want []byte, mode os.FileMode) (string, error) {
	baseTree, err := output(root, "rev-parse", destination+"^{tree}")
	if err != nil {
		return "", fmt.Errorf("read destination tree: %w", err)
	}
	tree, err := replaceTreeFile(root, baseTree, path, want, mode)
	if err != nil {
		return "", fmt.Errorf("neutralize spec path: %w", err)
	}
	if tree == baseTree {
		return destination, nil
	}
	commit, err := output(root, "commit-tree", tree, "-p", destination, "-m", "compose against the reviewed source spec")
	if err != nil {
		return "", fmt.Errorf("neutralize spec path: %w", err)
	}
	return commit, nil
}

func replaceTreeFile(root, baseTree, path string, content []byte, mode os.FileMode) (string, error) {
	return editTree(root, baseTree, func(idx string) error {
		blob, err := outputInput(root, content, "hash-object", "-w", "--stdin")
		if err != nil {
			return err
		}
		return indexRun(root, idx, "update-index", "--add", "--cacheinfo", gitRegularFileMode(mode)+","+blob+","+path)
	})
}

// removeTreeFolder writes baseTree without every entry beneath rel. The removal itself
// is removeIndexTree's, so the composed close and the commit-path close stay one fact.
func removeTreeFolder(root, baseTree, rel string) (string, error) {
	return editTree(root, baseTree, func(idx string) error { return removeIndexTree(root, idx, rel) })
}

// editTree reads baseTree into a private index, applies edit to that index, and
// writes the resulting tree. No checkout or repository index is touched.
func editTree(root, baseTree string, edit func(idx string) error) (string, error) {
	dir, err := os.MkdirTemp("", "bench-reviewed-landing-index-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)
	idx := filepath.Join(dir, "index")
	if err := indexRun(root, idx, "read-tree", baseTree); err != nil {
		return "", err
	}
	if err := edit(idx); err != nil {
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
func outputRaw(root string, args ...string) ([]byte, error) {
	return exec.Command("git", append([]string{"-C", root}, args...)...).Output()
}
func blobContent(root, oid string) ([]byte, error) {
	content, err := outputRaw(root, "cat-file", "blob", oid)
	if err != nil {
		return nil, fmt.Errorf("read blob %s: %w", oid, err)
	}
	return content, nil
}
func hashBlob(root string, content []byte) (string, error) {
	c := exec.Command("git", "-C", root, "hash-object", "-w", "--no-filters", "--stdin")
	c.Stdin = bytes.NewReader(content)
	b, err := c.Output()
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
