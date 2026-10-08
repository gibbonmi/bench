package gate

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/reviewrecord"
)

// Checkpoint names an explicit evidence obligation; its zero value is ordinary work.
type Checkpoint struct {
	Spec, Chunk string
	Complete    bool
}
type checkpointKey struct{}

// WithCheckpoint carries the obligation through the existing gate execution owner.
func WithCheckpoint(ctx context.Context, checkpoint Checkpoint) context.Context {
	return context.WithValue(ctx, checkpointKey{}, checkpoint)
}

func checkpointEvaluation(ctx context.Context, evaluation *gateEvaluation) *gateEvaluation {
	evaluation.checkpoint, _ = ctx.Value(checkpointKey{}).(Checkpoint)
	evaluation.completionSource, _ = ctx.Value(completionSourceKey{}).(string)
	evaluation.completionFolder, _ = ctx.Value(completionFolderKey{}).(string)
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
	// its reviewed source.
	completionProofError struct{ err error }
)

func (e evidenceError) Error() string        { return e.err.Error() }
func (e evidenceError) Unwrap() error        { return e.err }
func (e completionProofError) Error() string { return e.err.Error() }
func (e completionProofError) Unwrap() error { return e.err }

// errCheckpointTipMoved is a source tip that moved after the run accepted its subject.
var errCheckpointTipMoved = errors.New("checkpoint source tip changed")

// The checkpoint's refusal faces. The shared registry declares each face.
const (
	faceCompletionEvidence = "checkpoint-completion-evidence"
	faceDirtyCheckout      = "checkpoint-dirty-checkout"
	faceComposition        = "checkpoint-composition"
	faceTipMoved           = "checkpoint-tip-moved"
	faceSubjectUnavailable = "checkpoint-subject-unavailable"
	faceHandback           = "gate-handback"
)

// funnelFace picks the face of a refusal before the oracle by the type of its cause.
func funnelFace(err error) string {
	switch {
	case errors.As(err, new(evidenceError)):
		return faceCompletionEvidence
	case errors.Is(err, errCheckpointTipMoved):
		return faceTipMoved
	case errors.As(err, new(completionProofError)):
		return faceHandback
	}
	return faceSubjectUnavailable
}

// refuse prints a refusal's reason and then the route of face on its own next= line. The
// route reruns the caller's gate in the worktree of the active assignment that owns root;
// with no owner, the label prints its slot.
func refuse(ctx context.Context, root string, stderr io.Writer, mode runMode, face, reason string) Result {
	checkpoint, _ := ctx.Value(checkpointKey{}).(Checkpoint)
	owner, _ := intent.AssignmentForWorktree(root)
	slug, _ := reviewrecord.Slug(checkpoint.Spec)
	refusal := refusalroute.New(face, refusalroute.Facts{Sentence: reason, Values: map[string]string{
		refusalroute.FactLabel:     owner.Label,
		refusalroute.FactSlug:      slug,
		refusalroute.FactArguments: rerunArguments(checkpoint, mode == forceRun || face == faceSubjectUnavailable),
	}})
	result := operational(root, 0, stderr, refusal.Sentence)
	fmt.Fprintln(stderr, refusalroute.NextField+"="+refusal.Route)
	return result
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
	sourceTree := generation.tree
	if e.completionSource != "" {
		sourceTree, err = e.completionTree(generation)
		if err != nil {
			return subject{}, completionProofError{err}
		}
	}
	// A tickets-only close has no completion record; its completion tree is its evidence.
	if e.checkpoint.Spec != "" {
		if err := reviewrecord.CheckTrees(e.identityRoot, sourceTree, generation.tree, tip, e.checkpoint.Spec, e.checkpoint.Chunk, e.checkpoint.Complete); err != nil {
			return subject{}, evidenceError{fmt.Errorf("completion evidence: %w", err)}
		}
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
