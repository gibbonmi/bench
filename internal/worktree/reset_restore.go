package worktree

import (
	"errors"
	"strings"

	"github.com/gibbonmi/bench/internal/gate/authorization"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

func readResetRestore(root string, assignment intent.Assignment, ref string) (recoveryManifest, error) {
	if !strings.HasPrefix(ref, intent.ResetRefPrefix(assignment.OwnerID, assignment.ID)) {
		return recoveryManifest{}, errors.New("restore ref is not this assignment's")
	}
	rootOID, err := git.Output("-C", root, "rev-parse", "--verify", ref+"^{commit}")
	if err != nil || !resetEnvelopeValid(root, intent.Recovery{Ref: ref, Root: rootOID}) {
		return recoveryManifest{}, errors.New("reset envelope does not verify")
	}
	manifest, ok := readRecoveryManifest(root, ref)
	if !ok {
		return recoveryManifest{}, errors.New("reset envelope does not verify")
	}
	reached, err := authorization.IsAncestor(root, assignment.Start, manifest.Tip)
	if err != nil {
		return recoveryManifest{}, err
	}
	if !reached {
		return recoveryManifest{}, refusalError{refusal{detail: "envelope tip is not reached by the assignment start", wanted: assignment.Start}}
	}
	return manifest, nil
}

func restoreResetLayers(path string, manifest recoveryManifest) error {
	if _, err := git.Output("-C", path, "read-tree", "-u", "--reset", manifest.Layers["working"]+"^{tree}"); err != nil {
		return err
	}
	staged := manifest.Layers["staged"]
	if staged == "" {
		staged = manifest.Base
	}
	_, err := git.Output("-C", path, "read-tree", staged+"^{tree}")
	return err
}

func resetRestoredLayersMatch(root string, plan resetPlan) bool {
	rootOID, _, _, err := captureLayers(root, plan.assignment.Worktree, true, plan.manifest.Tip)
	if err != nil {
		return false
	}
	actual, ok := readRecoveryManifest(root, rootOID)
	if !ok || len(actual.Layers) != len(plan.manifest.Layers) {
		return false
	}
	for layer, payload := range plan.manifest.Layers {
		want, err := git.Output("-C", root, "rev-parse", payload+"^{tree}")
		if err != nil {
			return false
		}
		got, err := git.Output("-C", root, "rev-parse", actual.Layers[layer]+"^{tree}")
		if err != nil || got != want {
			return false
		}
	}
	return true
}
