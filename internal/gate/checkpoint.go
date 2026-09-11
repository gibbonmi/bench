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

func (e *gateEvaluation) applyCheckpoint(generation *treeGeneration, plan subject) (subject, error) {
	if err := e.checkpoint.validate(); err != nil {
		return subject{}, err
	}
	if e.checkpoint.Spec == "" {
		return plan, nil
	}
	tip, err := benchgit.Output("-C", e.identityRoot, "rev-parse", "--verify", "HEAD^{commit}")
	if err != nil {
		return subject{}, err
	}
	if e.checkpointTip == "" {
		e.checkpointTip = tip
	} else if e.checkpointTip != tip {
		return subject{}, errors.New("checkpoint source tip changed")
	}
	if err := reviewrecord.Check(e.identityRoot, generation.tree, tip, e.checkpoint.Spec, e.checkpoint.Chunk, e.checkpoint.Complete); err != nil {
		return subject{}, fmt.Errorf("completion evidence: %w", err)
	}
	purpose, _ := json.Marshal(e.checkpoint)
	hash := sha256.New()
	frame(hash, plan.Oracle)
	frame(hash, string(purpose))
	frame(hash, tip)
	plan.Oracle = hex.EncodeToString(hash.Sum(nil))
	return plan, nil
}
