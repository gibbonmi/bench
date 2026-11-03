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
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/gate/authorization"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/landing/published"
)

// Request is the complete, immutable input to one prospective landing.
type Request struct {
	Root, Destination, Expected, Message string
	Paths                                []string
	Stdout, Stderr                       io.Writer
	// Check is the optional pre-publication check. Land calls it with the candidate commit
	// before the ref update, and DryRun never calls it.
	Check func(candidate string) error
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
	// empty on every other landing. It never coexists with SpecPath: the landing's
	// --spec names either a staged spec.md to transition or a tickets-only folder to
	// close.
	ClosePath string
	// GreenMarker is the destination's green marker commit, empty when the destination has
	// none. The gate run reads it with the destination base to plan the reader set.
	GreenMarker    string
	Message        string
	Stdout, Stderr io.Writer
}

// ReviewedResult is the immutable publication receipt needed by the lifecycle owner. Kind is
// the authorization kind that published the tree.
type ReviewedResult struct {
	SourceBase, SourceTip, DestinationBase, Commit, Tree string
	Kind                                                 authorization.Kind
}

// composedSnapshot is the composed tree and the merge heads that its name-set check read.
type composedSnapshot struct {
	tree       string
	mergeHeads []string
}

// ReconcileError names the attributed path whose checkout reconciliation failed. Every
// return of the reconcile step carries one, so a caller renders the remaining path from
// a value rather than by parsing the message.
type ReconcileError struct {
	Path string
	Err  error
}

func (e *ReconcileError) Error() string {
	return fmt.Sprintf("reconcile named path %q: %v", e.Path, e.Err)
}
func (e *ReconcileError) Unwrap() error { return e.Err }

// PublishedUnreconciledError is the publication boundary: the commit is published and
// the checkout is not reconciled. It carries the published commit and the path that did
// not reconcile, so a caller reports the remainder without a second read.
type PublishedUnreconciledError struct {
	Commit string
	Path   string
	Err    error
}

func (e *PublishedUnreconciledError) Error() string {
	return fmt.Sprintf("landed-but-checkout-incomplete: %v", e.Err)
}
func (e *PublishedUnreconciledError) Unwrap() error { return e.Err }

// Owner composes only Request.Paths from Request.Expected, authorizes that tree, then
// publishes it by expected-old update. The function fields are narrow operational seams
// for deterministic fault coverage; New supplies the real owners.
type Owner struct {
	authorize func(context.Context, string, string, io.Writer, io.Writer) authorization.Result
	// publishes and reviewedPublishes name the authorization kinds each of this owner's
	// two landings publishes on. They are separate because the two landings answer to
	// different oracles: the worktree commit's fast lane, and main's whole-project gate.
	publishes         acceptedKinds
	reviewedPublishes acceptedKinds
	updateRef         func(string, string, string, string) error
	reconcile         func(Request, []string, composedSnapshot) error
}

// acceptedKinds is one landing's publication policy stated as data.
type acceptedKinds []authorization.Kind

func (k acceptedKinds) permits(kind authorization.Kind) bool {
	return slices.Contains(k, kind)
}

// New returns the production landing owner. The prospective landing of a worktree commit
// publishes on a graded green or on a fast-lane pass. The reviewed landing onto main
// publishes on a graded green or on a green reader set, so a lane pass never reaches main.
func New() Owner {
	return Owner{
		authorize:         authorization.AuthorizeWithWriters,
		publishes:         acceptedKinds{authorization.Green, authorization.LanePass},
		reviewedPublishes: acceptedKinds{authorization.Green, authorization.ReaderGreen},
		updateRef:         updateRef,
		reconcile:         reconcile,
	}
}

// composeAuthorized runs the prospective half every landing verdict shares: attribute,
// compose, refuse an empty diff, then authorize the exact composed tree. A pending merge
// waives the empty-diff refusal, because the merge parent is the change.
func (o Owner) composeAuthorized(ctx context.Context, r Request) ([]string, composedSnapshot, error) {
	if err := validRequest(r); err != nil {
		return nil, composedSnapshot{}, err
	}
	paths, heads, err := attributedPaths(r.Root, r.Expected, r.Paths)
	if err != nil {
		return nil, composedSnapshot{}, err
	}
	snapshot, err := compose(r, paths)
	if err != nil {
		return nil, composedSnapshot{}, err
	}
	snapshot.mergeHeads = heads
	baseTree, err := output(r.Root, "rev-parse", r.Expected+"^{tree}")
	if err != nil {
		return nil, composedSnapshot{}, fmt.Errorf("read expected base tree: %w", err)
	}
	if snapshot.tree == baseTree && len(heads) == 0 {
		return nil, composedSnapshot{}, NamedPathError{errors.New("nothing to commit")}
	}
	if got := o.authorize(ctx, r.Root, snapshot.tree, r.Stdout, r.Stderr); !o.publishes.permits(got.Kind) {
		return nil, composedSnapshot{}, AuthorizationRefusal{got}
	}
	return paths, snapshot, nil
}

// DryRun composes and authorizes exactly what Land would, then stops: no commit is
// created and no ref moves. The verdict is the value, so a caller diagnoses a composed
// path set without publishing a junk commit.
func (o Owner) DryRun(ctx context.Context, r Request) error {
	_, _, err := o.composeAuthorized(ctx, r)
	return err
}

// Land publishes a commit only when the exact composed tree receives green authorization.
// Parent 1 is the expected base, and each merge head that the composition read follows in
// file order. Land reads the pending merge again after authorization and refuses a changed
// answer before publication. A request's check then grades the candidate commit, and its
// refusal publishes nothing. After the ref update, Git ends a recorded merge before the
// named paths reconcile, so a reconcile failure leaves no pending merge.
func (o Owner) Land(ctx context.Context, r Request) (Result, error) {
	paths, snapshot, err := o.composeAuthorized(ctx, r)
	if err != nil {
		return Result{}, err
	}
	heads, err := benchgit.PendingMerge(r.Root)
	if err != nil {
		return Result{}, err
	}
	if !slices.Equal(heads, snapshot.mergeHeads) {
		return Result{}, errors.New("pending merge changed during the commit")
	}
	tree := snapshot.tree
	args := []string{"commit-tree", tree}
	for _, parent := range append([]string{r.Expected}, snapshot.mergeHeads...) {
		args = append(args, "-p", parent)
	}
	commit, err := output(r.Root, append(args, "-m", r.Message)...)
	if err != nil {
		return Result{}, fmt.Errorf("create landing commit: %w", err)
	}
	if err := checkCandidate(r, commit); err != nil {
		return Result{}, err
	}
	if err := o.updateRef(r.Root, r.Destination, commit, r.Expected); err != nil {
		return Result{}, destinationUpdateFailure(r.Root, r.Destination, r.Expected, err)
	}
	result := Result{Base: r.Expected, Commit: commit, Tree: tree}
	// Git owns its merge-state files, so Git removes them.
	if len(snapshot.mergeHeads) > 0 {
		if err := run(r.Root, "merge", "--quit"); err != nil {
			return result, &PublishedUnreconciledError{Commit: commit, Err: fmt.Errorf("end the pending merge: %w", err)}
		}
	}
	if err := o.reconcile(r, paths, snapshot); err != nil {
		var failed *ReconcileError
		remainder := &PublishedUnreconciledError{Commit: commit, Err: err}
		if errors.As(err, &failed) {
			remainder.Path = failed.Path
		}
		return result, remainder
	}
	return result, nil
}

// LandReviewed composes an exact reviewed source, applies its staged-spec
// transition and verified closure only to that prospective tree, and publishes a
// two-parent commit.
// The worktree lifecycle owns authentication, marker advancement, reconciliation,
// and release around this irreversible operation. It applies no admission.
func (o Owner) LandReviewed(ctx context.Context, r ReviewedRequest) (ReviewedResult, error) {
	return o.landReviewed(ctx, r, unadmitted{})
}

func (o Owner) landReviewed(ctx context.Context, r ReviewedRequest, admission Admission) (ReviewedResult, error) {
	destination, source, err := reviewedSource(r)
	if err != nil {
		return ReviewedResult{}, err
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
	tree, err := o.landingTree(r, composed, source)
	if err != nil {
		return ReviewedResult{}, err
	}
	if err := admission.Check(tree); err != nil {
		return ReviewedResult{}, err
	}
	if deliverable := r.Deliverable(); deliverable != "" {
		ctx = gate.WithCompletion(ctx, deliverable, source)
	}
	ctx = gate.WithLanding(ctx, destination, r.GreenMarker)
	got := o.authorize(ctx, r.Root, tree, r.Stdout, r.Stderr)
	if !o.reviewedPublishes.permits(got.Kind) {
		return ReviewedResult{}, AuthorizationRefusal{got}
	}
	if err := sourceUnchanged(r, source); err != nil {
		return ReviewedResult{}, err
	}
	if err := destinationUnchanged(r.Root, r.Destination, destination, r.DestinationFingerprint); err != nil {
		return ReviewedResult{}, err
	}
	var commit string
	err = admission.Publish(tree, func() error {
		created, err := landingMerge(r.Root, tree, destination, source, r.Message)
		if err != nil {
			return err
		}
		if err := o.publishMerge(r.Root, r.Destination, created, destination); err != nil {
			return err
		}
		commit = created
		return nil
	})
	if err != nil {
		return ReviewedResult{}, err
	}
	return ReviewedResult{SourceBase: r.ReviewBase, SourceTip: source, DestinationBase: destination, Commit: commit, Tree: tree, Kind: got.Kind}, nil
}

// reviewedSource proves a reviewed request complete, its destination base and source tip
// exact commits, its branch at the reviewed tip, and its review base valid. It answers the
// destination and the source commits.
func reviewedSource(r ReviewedRequest) (string, string, error) {
	if r.Root == "" || r.Destination == "" || r.DestinationBase == "" || r.Source == "" || r.SourceTip == "" || r.ReviewBase == "" || r.SourceWorktree == "" || r.SourceFingerprint == "" || r.DestinationFingerprint == "" || strings.TrimSpace(r.Message) == "" {
		return "", "", errors.New("reviewed landing request is incomplete")
	}
	destination, err := compositionCommit(r.Root, r.DestinationBase, "destination")
	if err != nil || destination != r.DestinationBase {
		return "", "", errors.New("destination base is not an exact commit")
	}
	source, err := compositionCommit(r.Root, r.SourceTip, "source")
	if err != nil || source != r.SourceTip {
		return "", "", errors.New("source tip is not an exact commit")
	}
	branchTip, err := benchgit.ResolveCommit(r.Root, r.Source)
	if err != nil || branchTip != source {
		return "", "", errors.New("reviewed source tip moved")
	}
	if _, kind, _ := diff.ResolveSourceRange(r.Root, r.ReviewBase, source); kind != "" {
		return "", "", errors.New("reviewed source base is invalid")
	}
	return destination, source, nil
}

// landingTree composes source onto the commit onto and applies the verified delivery of
// r, so it answers the tree that the landing of r publishes there.
func (o Owner) landingTree(r ReviewedRequest, onto, source string) (string, error) {
	tree, err := published.Subject(r.Root, onto, source, r.ReviewBase, r.Deliverable(), r.Stderr)
	var conflict published.ConflictError
	if errors.As(err, &conflict) {
		return "", ConflictError(conflict)
	}
	return tree, err
}

// AuthorizationRefusal is the refusal of a tree the gate did not authorize. It carries the
// attributed result, so a caller picks the refusal's route by the kind and the reason and
// never by the sentence.
type AuthorizationRefusal struct{ Result authorization.Result }

func (r AuthorizationRefusal) Error() string { return refusalMessage(r.Result) }

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
	return withinIgnoredDir(path, ".logs/")
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
		if entry.Status == "!!" && (RuntimeIgnoredPath(entry.Path) || LocalCapturePath(entry.Path)) {
			continue
		}
		filtered.WriteString(entry.Status)
		filtered.WriteByte(' ')
		filtered.WriteString(entry.Path)
		filtered.WriteByte(0)
	}
	return filtered.Bytes(), nil
}
