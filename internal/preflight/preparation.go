package preflight

import (
	"fmt"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/tickets"
	"github.com/gibbonmi/bench/internal/toon"
)

type preparationMode int

const (
	chargePreparation preparationMode = iota
	proposalPreparation
)

func preparedCommand(
	root, mode, slug, base, sourceTip, name string,
	full bool,
	version string,
	args []string,
	form preparationMode,
) (string, int) {
	action := "charge"
	if form == proposalPreparation {
		action = "proposal"
	}
	return preparedAttempts(root, mode, slug, base, sourceTip, action, args, func(facts Facts) (string, int) {
		if form == proposalPreparation {
			return renderWritesProposal(root, facts, name)
		}
		return renderReviewCharge(root, facts, Decide(facts), full, version)
	})
}

// preparedAttempts gathers the pinned facts inside the movement-checked retry and renders
// each attempt through render. A drift or read failure replaces the attempt's output, so
// only an unmoved attempt's answer reaches the caller.
func preparedAttempts(root, mode, slug, base, sourceTip, action string, args []string, render func(Facts) (string, int)) (string, int) {
	var out string
	code := 1
	result := diff.MovementCheckedRetry(root, func(snapshot diff.MovementSnapshot) (string, string) {
		out, code = "", 1
		source, kind, hint := snapshot.ResolveSourceRange(base)
		if kind != "" {
			out = toon.Errorf(kind, hint) + "\n"
			return "", ""
		}
		paths, err := snapshot.SourceSnapshotPaths(source)
		if err != nil {
			out = toon.Errorf("source paths failed", err.Error()) + "\n"
			return "", ""
		}
		facts, failure := gatherCharge(root, mode, slug, &source, paths, sourceTip)
		if failure != nil {
			out = chargeRefusal("source", failure.Kind+": "+failure.Hint, "restore the named canonical source and rerun the exact "+action)
			return "", ""
		}
		if err := unrepresentableChangedPath(facts.ChangedPaths); err != nil {
			out = toon.RenderError(err) + "\n"
			return "", ""
		}
		out, code = render(facts)
		return "", ""
	})
	if result.DriftKind != "" {
		return snapshotDriftRefusal(args, result.DriftHint), 1
	}
	if result.Kind != "" {
		return toon.Errorf(result.Kind, result.Hint) + "\n", 1
	}
	return out, code
}

func preparationCheckoutRefusal(root string, facts Facts, action string) string {
	if facts.AssignmentTarget == "" {
		return chargeRefusal("assignment", "active assignment is required", "run from the assigned worktree")
	}
	if _, dirty, err := git.AllFilesStatus(root); err != nil {
		return chargeRefusal("checkout", err.Error(), "repair the checkout and rerun the exact "+action)
	} else if len(dirty) != 0 {
		return chargeRefusal("checkout", "source checkout is dirty", "commit or remove local changes and rerun the exact "+action)
	}
	return ""
}

func preparationTicket(root string, facts Facts, name string) (*tickets.Entry, *tickets.Ticket, string, string) {
	dir := filepath.Join(root, filepath.Dir(facts.SpecPath), "tickets")
	classified := bounds.ClassifyDirNoFollow(dir)
	if classified.State != bounds.StateParsed && classified.State != bounds.StateEmpty {
		return nil, nil, "tickets directory is " + string(classified.State) + ": " + classified.Reason, "repair the selected ticket directory"
	}
	entries, _, refusal := tickets.EnumerateNoFollow(dir, classified.Entries)
	if refusal != nil {
		return nil, nil, refusal.Message(dir), "repair the selected ticket"
	}
	selected := selectedTicket(entries, name)
	if selected == nil {
		return nil, nil, fmt.Sprintf("selected ticket %q was not found", name), "pass a ticket basename from the spec tickets directory"
	}
	parsed := selectedParsedTicket(facts, name)
	if parsed == nil {
		return nil, nil, "selected ticket has no parsed evidence", "repair ticket grammar and rerun the exact command"
	}
	return selected, parsed, "", ""
}
