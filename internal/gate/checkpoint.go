package gate

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"strings"

	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/spec"
)

// Checkpoint names an explicit evidence obligation; its zero value is ordinary work.
type Checkpoint struct {
	Spec, Chunk string
	Complete    bool
}
type checkpointKey struct{}
type verbKey struct{}

// chunkSourceKey carries the committed source that a chunk checkpoint accepted.
type chunkSourceKey struct{}

// chunkSource is the HEAD commit and the HEAD tree of a chunk checkpoint.
type chunkSource struct{ tip, tree string }

// WithCheckpoint carries the obligation through the existing gate execution owner.
func WithCheckpoint(ctx context.Context, checkpoint Checkpoint) context.Context {
	return context.WithValue(ctx, checkpointKey{}, checkpoint)
}

// asVerb marks a run that the gate verb itself runs. Only that run prints the gate's route.
// Another caller, such as the authorization run of a landing or a commit, prints its own
// route for the refusal that the run returns, so one route prints per refusal.
func asVerb(ctx context.Context) context.Context { return context.WithValue(ctx, verbKey{}, true) }

func checkpointEvaluation(ctx context.Context, evaluation *gateEvaluation) *gateEvaluation {
	evaluation.checkpoint, _ = ctx.Value(checkpointKey{}).(Checkpoint)
	evaluation.completionSource, _ = ctx.Value(completionSourceKey{}).(string)
	evaluation.completionFolder, _ = ctx.Value(completionFolderKey{}).(string)
	// The accepted tip starts the tip-moved guard, so a tip that moves before the first
	// capture refuses too.
	source, _ := ctx.Value(chunkSourceKey{}).(chunkSource)
	evaluation.checkpointTip, evaluation.checkpointTree = source.tip, source.tree
	return evaluation
}

func (checkpoint Checkpoint) validate() error {
	if checkpoint.Spec == "" {
		if checkpoint.Chunk != "" || checkpoint.Complete {
			return errors.New("checkpoint spec is required")
		}
		return nil
	}
	if _, err := reviewrecord.RecordPath(checkpoint.Spec); err != nil {
		return err
	}
	if (checkpoint.Chunk != "") == checkpoint.Complete || hasUnsafeText(checkpoint.Chunk) {
		return errors.New("select one checkpoint chunk or complete")
	}
	return nil
}

// The gate's flags. The parser, the usage line, and a route's rerun read these spellings.
const (
	flagFresh      = "--fresh"
	flagCheckpoint = "--checkpoint"
	flagChunk      = "--chunk"
	flagComplete   = "--complete"
)

func parseGateArgs(args []string, plumbing bool) (string, runMode, Checkpoint, error) {
	root, mode, checkpoint := "", reuseFreshGreen, Checkpoint{}
	seen := map[string]bool{}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "--") {
			if seen[arg] {
				return "", mode, checkpoint, errors.New("duplicate gate argument")
			}
			seen[arg] = true
		}
		switch arg {
		case flagFresh:
			mode = forceRun
		case flagCheckpoint, flagChunk:
			i++
			if i == len(args) || args[i] == "" || strings.HasPrefix(args[i], "--") {
				return "", mode, checkpoint, errors.New("missing checkpoint value")
			}
			if arg == flagCheckpoint {
				checkpoint.Spec = args[i]
			} else {
				checkpoint.Chunk = args[i]
			}
		case flagComplete:
			checkpoint.Complete = true
		default:
			if !plumbing || root != "" || strings.HasPrefix(arg, "-") {
				return "", mode, checkpoint, errors.New("unknown gate argument")
			}
			root = arg
		}
	}
	return root, mode, checkpoint, checkpoint.validate()
}

// The checkpoint causes that the funnel before the oracle routes by type. Every other
// fault there captured no subject.
type (
	// evidenceError is a completion-evidence refusal.
	evidenceError struct{ err error }
	// completionProofError is a completion whose graded tree is not the exact transform of
	// its reviewed source. A fault that captures or reads either tree proves nothing about
	// the transform, so it is not a proof fault.
	completionProofError struct{ err error }
)

func (e evidenceError) Error() string        { return e.err.Error() }
func (e evidenceError) Unwrap() error        { return e.err }
func (e completionProofError) Error() string { return e.err.Error() }
func (e completionProofError) Unwrap() error { return e.err }

// errCheckpointTipMoved is a source tip that moved after the run accepted its subject.
var errCheckpointTipMoved = errors.New("checkpoint source tip changed")

// checkpointFace is one of the checkpoint's refusal faces, which the shared registry
// declares, and the record form that the refused cause names, if any.
type checkpointFace struct{ name, form string }

// The checkpoint's refusal faces.
var (
	faceCompletionEvidence = checkpointFace{name: "checkpoint-completion-evidence"}
	faceDirtyCheckout      = checkpointFace{name: "checkpoint-dirty-checkout"}
	faceComposition        = checkpointFace{name: "checkpoint-composition"}
	faceTipMoved           = checkpointFace{name: "checkpoint-tip-moved"}
	faceSubjectUnavailable = checkpointFace{name: "checkpoint-subject-unavailable"}
	faceHandback           = checkpointFace{name: CompletionProofFace}
)

// CompletionProofFace is the face of a completion proof fault, the one funnel cause that
// hands back to the reviewer. A caller outside the gate verb reads it on Result.Face.
const CompletionProofFace = "gate-handback"

// faceRecordRepair is the face of an evidence refusal that one record form repairs.
const faceRecordRepair = "checkpoint-record-repair"

// funnelFace picks the face of a refusal before the oracle by the type of its cause. An
// evidence refusal that one record form repairs names that form; every other evidence
// refusal keeps the completion-evidence face.
func funnelFace(err error) checkpointFace {
	var repair reviewrecord.Repair
	switch {
	case errors.As(err, new(evidenceError)) && errors.As(err, &repair):
		return checkpointFace{name: faceRecordRepair, form: repair.Form}
	case errors.As(err, new(evidenceError)):
		return faceCompletionEvidence
	case errors.Is(err, errCheckpointTipMoved):
		return faceTipMoved
	case errors.As(err, new(completionProofError)):
		return faceHandback
	}
	return faceSubjectUnavailable
}

// refuse prints a refusal's reason. In a run of the gate verb, the route of face follows on
// its own next= line. The route reruns the caller's gate in the worktree of the active
// assignment that owns root; with no owner, the label prints its slot.
func refuse(ctx context.Context, root string, stderr io.Writer, mode runMode, face checkpointFace, reason string) Result {
	return refuseWith(ctx, root, stderr, mode, face, reason, nil)
}

// refuseWith is refuse with the values that the raising policy observed, which join the
// checkpoint's own. The caller's gate re-run is one of the checkpoint's values, for a face
// whose route ends with the printing verb's re-run.
func refuseWith(ctx context.Context, root string, stderr io.Writer, mode runMode, face checkpointFace, reason string, raised map[string]string) Result {
	if verb, _ := ctx.Value(verbKey{}).(bool); !verb {
		result := operational(root, 0, stderr, reason)
		result.Face = face.name
		return result
	}
	checkpoint, _ := ctx.Value(checkpointKey{}).(Checkpoint)
	owner, _ := intent.AssignmentForWorktree(root)
	slug, _ := reviewrecord.Slug(checkpoint.Spec)
	values := map[string]string{
		refusalroute.FactLabel:     owner.Label,
		refusalroute.FactSlug:      slug,
		refusalroute.FactForm:      face.form,
		refusalroute.FactArguments: rerunArguments(checkpoint, mode == forceRun || face == faceSubjectUnavailable),
	}
	values[refusalroute.FactRerun] = refusalroute.GateRerun(values)
	maps.Copy(values, raised)
	refusal := refusalroute.New(face.name, refusalroute.Facts{Sentence: reason, Values: values})
	result := operational(root, 0, stderr, refusal.Sentence)
	fmt.Fprintln(stderr, refusalroute.NextField+"="+refusal.Route)
	return result
}

// refuseUninstalled runs the candidate check of a checkpoint of a declared build, before
// any ledger read or write. The declaration comes from the spec in the checkpoint source
// tree, the commit that HEAD names. It answers whether the check refused, and the refusal.
// The reviewer signs off a declaration, so a declaration that does not parse hands back,
// and a spec that does not read leaves the checkpoint no subject.
func refuseUninstalled(ctx context.Context, root string, checkpoint Checkpoint, stderr io.Writer, mode runMode) (Result, bool) {
	if checkpoint.Spec == "" {
		return Result{}, false
	}
	err := commitrepo.CheckEntry(root, checkpoint.Spec)
	var raised refusalroute.Raised
	face := faceSubjectUnavailable
	switch {
	case err == nil:
		return Result{}, false
	case errors.As(err, &raised):
		face = checkpointFace{name: raised.Name}
	case errors.Is(err, spec.ErrPersistedFieldsInvalid):
		face = faceHandback
	}
	return refuseWith(ctx, root, stderr, mode, face, err.Error(), raised.Values), true
}

// cleanCheckoutRefusal names the recovery for a checkpoint on a dirty checkout.
const cleanCheckoutRefusal = "checkpoint requires a clean checkout: commit or remove each uncommitted change"

// committedSource answers the HEAD commit of root and its tree. A checkpoint grades only
// committed bytes, so a working tree that differs from that tree refuses, and the review
// record has no exemption. It also answers whether it refused, and the refusal. A caller runs
// it before the evidence pass and the oracle.
func committedSource(ctx context.Context, root string, stderr io.Writer, mode runMode) (string, string, Result, bool) {
	unavailable := func(err error) (string, string, Result, bool) {
		return "", "", sourceUnavailable(ctx, root, stderr, mode, err), true
	}
	tip, err := benchgit.ResolveCommit(root, "HEAD")
	if err != nil {
		return unavailable(err)
	}
	tree, err := benchgit.Output("-C", root, "rev-parse", "--verify", tip+"^{tree}")
	if err != nil {
		return unavailable(err)
	}
	working, err := workingTreeSource{root: root}.tree()
	if err != nil {
		return unavailable(err)
	}
	if working != tree {
		return "", "", refuse(ctx, root, stderr, mode, faceDirtyCheckout, cleanCheckoutRefusal), true
	}
	return tip, tree, Result{}, false
}

// sourceUnavailable refuses a checkpoint whose committed source does not read.
func sourceUnavailable(ctx context.Context, root string, stderr io.Writer, mode runMode, err error) Result {
	return refuse(ctx, root, stderr, mode, faceSubjectUnavailable, "checkpoint source unavailable: "+err.Error())
}

// executeChunkCheckpoint grades the HEAD tree of root without the review record of spec,
// once committedSource accepts the checkout. That tree runs in a private checkout with root
// as the storage and identity root, so a record-only commit keeps the subject of an earlier
// chunk green. The evidence pass and the record grade still read the record from the HEAD
// tree on every call.
func executeChunkCheckpoint(ctx context.Context, root, spec string, stdout, stderr io.Writer, arm postAcquireContextArm, mode runMode) Result {
	tip, tree, refusal, refused := committedSource(ctx, root, stderr, mode)
	if refused {
		return refusal
	}
	stripped, err := reviewrecord.SourceDigest(root, tree, spec)
	if err != nil {
		return sourceUnavailable(ctx, root, stderr, mode, err)
	}
	ctx = context.WithValue(ctx, chunkSourceKey{}, chunkSource{tip: tip, tree: tree})
	return executeCheckout(ctx, root, stripped, stdout, stderr, arm, mode, nil, newCheckoutTreeEvaluation)
}

// rerunArguments composes the caller's gate arguments, each value rendered by Arg. A fresh
// rerun names --fresh.
func rerunArguments(checkpoint Checkpoint, fresh bool) string {
	var args []string
	if fresh {
		args = append(args, flagFresh)
	}
	if checkpoint.Spec != "" {
		args = append(args, flagCheckpoint, refusalroute.Arg("spec-path", checkpoint.Spec))
		if checkpoint.Complete {
			args = append(args, flagComplete)
		} else {
			args = append(args, flagChunk, refusalroute.Arg("id", checkpoint.Chunk))
		}
	}
	return strings.Join(args, " ")
}

func (e *gateEvaluation) applyCheckpoint(generation *treeGeneration, plan subject) (subject, error) {
	if err := validateCompletionContext(e); err != nil {
		return subject{}, err
	}
	if err := e.checkpoint.validate(); err != nil {
		return subject{}, err
	}
	if e.checkpoint.Spec == "" && e.completionFolder == "" {
		return plan, nil
	}
	tip, err := benchgit.ResolveCommit(e.identityRoot, "HEAD")
	if e.completionSource != "" {
		tip, err = benchgit.ResolveCommit(e.identityRoot, e.completionSource)
		if err == nil && (!e.prospective || !e.completing() || tip != e.completionSource) {
			return subject{}, completionProofError{errors.New("invalid prospective completion source")}
		}
	}
	if err != nil {
		return subject{}, err
	}
	if e.checkpointTip == "" {
		e.checkpointTip = tip
	} else if e.checkpointTip != tip {
		return subject{}, errCheckpointTipMoved
	}
	sourceTree, evidenceTree := generation.tree, generation.tree
	if e.checkpointTree != "" {
		// The graded tree of a chunk checkpoint holds no record, so its committed tree is
		// both its source and its evidence.
		sourceTree, evidenceTree = e.checkpointTree, e.checkpointTree
	}
	if e.completionSource != "" {
		sourceTree, err = e.completionTree(generation)
		if err != nil {
			return subject{}, err
		}
	}
	// A tickets-only close has no completion record; its completion tree is its evidence.
	if e.checkpoint.Spec != "" {
		if err := reviewrecord.CheckTrees(e.identityRoot, sourceTree, evidenceTree, tip, e.checkpoint.Spec, e.checkpoint.Chunk, e.checkpoint.Complete); err != nil {
			return subject{}, evidenceError{fmt.Errorf("completion evidence: %w", err)}
		}
	}
	// A chunk green answers for its subject alone, because the evidence pass above runs on
	// every call. So the chunk checkpoint keeps the base oracle of its subject.
	if e.checkpoint.Chunk != "" {
		return plan, nil
	}
	purpose, _ := json.Marshal(e.checkpoint)
	hash := sha256.New()
	frame(hash, plan.Oracle)
	frame(hash, string(purpose))
	if e.completionFolder != "" {
		frame(hash, e.completionFolder)
	}
	frame(hash, tip)
	plan.Oracle = hex.EncodeToString(hash.Sum(nil))
	return plan, nil
}
