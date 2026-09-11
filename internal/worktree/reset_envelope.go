package worktree

import (
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

func writeResetEnvelope(root string, plan resetPlan) (intent.Recovery, error) {
	rootOID, payloads, _, err := captureLayers(root, plan.assignment.Worktree, true, plan.tip)
	if err != nil {
		return intent.Recovery{}, err
	}
	ref, err := nextEnvelopeRef(root, intent.ResetRefPrefix(plan.assignment.OwnerID, plan.assignment.ID))
	if err != nil {
		return intent.Recovery{}, err
	}
	recovery := intent.Recovery{Ref: ref, Root: rootOID, Payloads: payloads}
	_, err = git.Output("-C", root, "update-ref", ref, rootOID, strings.Repeat("0", len(rootOID)))
	return recovery, err
}

func resetEnvelopeValid(root string, envelope intent.Recovery) bool {
	resolved, err := git.Output("-C", root, "rev-parse", "--verify", envelope.Ref+"^{commit}")
	if err != nil || resolved != envelope.Root {
		return false
	}
	manifest, ok := readRecoveryManifest(root, envelope.Root)
	if !ok || manifest.Tip == "" {
		return false
	}
	parents, err := git.Output("-C", root, "show", "-s", "--format=%P", envelope.Root)
	if err != nil {
		return false
	}
	parentSet := map[string]bool{}
	for _, parent := range strings.Fields(parents) {
		parentSet[parent] = true
	}
	for _, payload := range manifest.Layers {
		if !parentSet[payload] {
			return false
		}
	}
	return true
}
