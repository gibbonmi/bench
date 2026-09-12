package gate

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	benchgit "github.com/gibbonmi/bench/internal/git"
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
		case "--fresh":
			mode = forceRun
		case "--checkpoint", "--chunk":
			i++
			if i == len(args) || args[i] == "" || strings.HasPrefix(args[i], "--") {
				return "", mode, checkpoint, errors.New("missing checkpoint value")
			}
			if arg == "--checkpoint" {
				checkpoint.Spec = args[i]
			} else {
				checkpoint.Chunk = args[i]
			}
		case "--complete":
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

// routeError is a refusal that carries the harness-native read which answers it.
// The refusal printer adds that route under the reason. Only the completion-
// evidence refusal builds one: the checkpoint cannot show what the evidence
// lacks, and one read reports it.
type routeError struct {
	next string
	err  error
}

func (e routeError) Error() string { return e.err.Error() }
func (e routeError) Unwrap() error { return e.err }

// routedRefusal routes the completion-evidence refusal to the preflight review
// read. The slug comes from the spec-path grammar's one owner, so the route can
// never name a spec the checkpoint did not open. A spec whose slug does not
// resolve carries no route; validate already refused that path upstream.
func routedRefusal(spec string, err error) error {
	slug, slugErr := reviewrecord.Slug(spec)
	if slugErr != nil {
		return err
	}
	return routeError{next: "bench preflight review " + slug, err: err}
}

func (e *gateEvaluation) applyCheckpoint(generation *treeGeneration, plan subject) (subject, error) {
	if err := validateCompletionContext(e); err != nil {
		return subject{}, err
	}
	if err := e.checkpoint.validate(); err != nil {
		return subject{}, err
	}
	if e.checkpoint.Spec == "" {
		return plan, nil
	}
	tip, err := benchgit.Output("-C", e.identityRoot, "rev-parse", "--verify", "HEAD^{commit}")
	if e.completionSource != "" {
		tip, err = benchgit.Output("-C", e.identityRoot, "rev-parse", "--verify", e.completionSource+"^{commit}")
		if err == nil && (!e.prospective || !e.checkpoint.Complete || tip != e.completionSource) {
			return subject{}, errors.New("invalid prospective completion source")
		}
	}
	if err != nil {
		return subject{}, err
	}
	if e.checkpointTip == "" {
		e.checkpointTip = tip
	} else if e.checkpointTip != tip {
		return subject{}, errors.New("checkpoint source tip changed")
	}
	sourceTree := generation.tree
	if e.completionSource != "" {
		sourceTree, err = e.completionTree(generation)
		if err != nil {
			return subject{}, err
		}
	}
	if err := reviewrecord.CheckTrees(e.identityRoot, sourceTree, generation.tree, tip, e.checkpoint.Spec, e.checkpoint.Chunk, e.checkpoint.Complete); err != nil {
		return subject{}, routedRefusal(e.checkpoint.Spec, fmt.Errorf("completion evidence: %w", err))
	}
	purpose, _ := json.Marshal(e.checkpoint)
	hash := sha256.New()
	frame(hash, plan.Oracle)
	frame(hash, string(purpose))
	frame(hash, tip)
	plan.Oracle = hex.EncodeToString(hash.Sum(nil))
	return plan, nil
}
