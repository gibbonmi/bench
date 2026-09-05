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
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/gate/authorization"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/spec"
)

// Request is the complete, immutable input to one prospective landing.
type Request struct {
	Root, Destination, Expected, Message string
	Paths                                []string
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
	// empty on every other landing. It never coexists with SpecPath: the landing's
	// --spec names either a staged spec.md to transition or a tickets-only folder to
	// close.
	ClosePath      string
	Message        string
	Stdout, Stderr io.Writer
}

// ReviewedResult is the immutable publication receipt needed by the lifecycle owner.
type ReviewedResult struct {
	SourceBase, SourceTip, DestinationBase, Commit, Tree string
}

type composedSnapshot struct{ tree string }

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
// the checkout is not reconciled. It carries the published commit, the path that did not
// reconcile, and every named path in the owner's sorted, deduplicated order, so a caller
// reports the remainder and the repair without a second read.
type PublishedUnreconciledError struct {
	Commit string
	Path   string
	Paths  []string
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
	for _, accepted := range k {
		if kind == accepted {
			return true
		}
	}
	return false
}

// New returns the production landing owner. The prospective landing of a worktree commit
// publishes on a graded green or on a fast-lane pass. The reviewed landing onto main
// publishes on a graded green alone, so a lane pass never reaches main.
func New() Owner {
	return Owner{
		authorize:         authorization.AuthorizeWithWriters,
		publishes:         acceptedKinds{authorization.Green, authorization.LanePass},
		reviewedPublishes: acceptedKinds{authorization.Green},
		updateRef:         updateRef,
		reconcile:         reconcile,
	}
}

// NewLane returns the worktree-commit owner whose authority is the root's declared fast
// lane rather than the whole-project gate.
func NewLane(lane authorization.LaneAuthority) Owner {
	owner := New()
	owner.authorize = lane.Authorize
	return owner
}

// laneAuthority maps a resolved lane and the base it is measured against onto the
// authority that grades a composed tree. It is the one place the lane's fields become an
// authority, so no caller states the mapping a second time.
func laneAuthority(lane *gate.Lane, base string) authorization.LaneAuthority {
	return authorization.LaneAuthority{
		Checks:    lane.Checks,
		Kit:       lane.Kit,
		Selective: lane.Selective,
		Base:      base,
	}
}

// NewForLane returns the owner a caller lands under given the lane its root declares. A
// nil lane is a root that declares none, so it answers the whole-project gate owner.
func NewForLane(lane *gate.Lane, base string) Owner {
	if lane == nil {
		return New()
	}
	return NewLane(laneAuthority(lane, base))
}

// composeAuthorized runs the prospective half every landing verdict shares: attribute,
// compose, refuse an empty diff, then authorize the exact composed tree.
func (o Owner) composeAuthorized(ctx context.Context, r Request) ([]string, composedSnapshot, error) {
	if err := validRequest(r); err != nil {
		return nil, composedSnapshot{}, err
	}
	paths, err := attributedPaths(r.Root, r.Expected, r.Paths)
	if err != nil {
		return nil, composedSnapshot{}, err
	}
	snapshot, err := compose(r, paths)
	if err != nil {
		return nil, composedSnapshot{}, err
	}
	baseTree, err := output(r.Root, "rev-parse", r.Expected+"^{tree}")
	if err != nil {
		return nil, composedSnapshot{}, fmt.Errorf("read expected base tree: %w", err)
	}
	if snapshot.tree == baseTree {
		return nil, composedSnapshot{}, errors.New("nothing to commit")
	}
	if got := o.authorize(ctx, r.Root, snapshot.tree, r.Stdout, r.Stderr); !o.publishes.permits(got.Kind) {
		return nil, composedSnapshot{}, errors.New(refusalMessage(got))
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
func (o Owner) Land(ctx context.Context, r Request) (Result, error) {
	paths, snapshot, err := o.composeAuthorized(ctx, r)
	if err != nil {
		return Result{}, err
	}
	tree := snapshot.tree
	commit, err := output(r.Root, "commit-tree", tree, "-p", r.Expected, "-m", r.Message)
	if err != nil {
		return Result{}, fmt.Errorf("create landing commit: %w", err)
	}
	if err := o.updateRef(r.Root, r.Destination, commit, r.Expected); err != nil {
		return Result{}, destinationUpdateFailure(r.Root, r.Destination, r.Expected, err)
	}
	if err := o.reconcile(r, paths, snapshot); err != nil {
		var failed *ReconcileError
		remainder := &PublishedUnreconciledError{Commit: commit, Paths: paths, Err: err}
		if errors.As(err, &failed) {
			remainder.Path = failed.Path
		}
		return Result{Base: r.Expected, Commit: commit, Tree: tree}, remainder
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
	// The close consumes the tickets-only folder from the published tree by index
	// removal. A folder the destination already removed lists no entries, so the removal
	// writes the composed tree back unchanged.
	if r.ClosePath != "" {
		if tree, err = removeTreeFolder(r.Root, tree, r.ClosePath); err != nil {
			return ReviewedResult{}, fmt.Errorf("close tickets-only folder: %w", err)
		}
	}
	if got := o.authorize(ctx, r.Root, tree, r.Stdout, r.Stderr); !o.reviewedPublishes.permits(got.Kind) {
		return ReviewedResult{}, errors.New(refusalMessage(got))
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

// refusalMessage renders the one refusal line every authorization caller prints. The
// literal prefix and the kind stay, because two tests and the operator's own memory read
// them; the sentence after names what the attribution means and what to run next.
func refusalMessage(got authorization.Result) string {
	line := "prospective authorization refused: " + string(got.Kind)
	explanation, action := refusalGuidance(got.Kind)
	// An Infrastructure attribution carries the gate's own reason, which is more exact than
	// anything this renderer could state about the kind.
	if got.Reason != "" {
		explanation = got.Reason
	}
	if explanation != "" {
		line += " (" + explanation + ")"
	}
	if action != "" {
		line += "; " + action
	}
	return line
}

// refusalGuidance answers what a refused kind means to the operator and what to do next.
// The two lane outcomes answer nothing: a lane states its own outcome line already.
func refusalGuidance(kind authorization.Kind) (explanation, action string) {
	switch kind {
	case authorization.Inherited:
		return "the gate ran red on the composed tree and no green baseline attributes the red to this diff", "run bench gate --fresh"
	case authorization.Candidate:
		return "the gate ran red on the composed tree and the green baseline attributes the red to this diff", "fix the failures above"
	case authorization.Infrastructure:
		return "", "run bench doctor"
	}
	return "", ""
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

// LocalCapturePath reports whether path is one exact project-local capture file.
func LocalCapturePath(path string) bool {
	switch path {
	case "capture/IDEAS.md", "capture/learnings.md", "capture/session-handoff.md":
		return true
	default:
		return false
	}
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
