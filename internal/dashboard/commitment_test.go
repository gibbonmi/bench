package dashboard

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/intent"
)

// TestRenderCommitmentOutlook holds the commitment card to the outlook in the snapshot. The
// card names the eligible outcome and each blocker with its reason, and it escapes hostile
// cells. Only an absent commitment marks the recommended sequence as unapproved input.
func TestRenderCommitmentOutlook(t *testing.T) {
	const unapproved = "Recommended sequence (unapproved input: no commitment is adopted)"
	s := baseSnapshot()
	s.RoadmapPresent = true
	s.RoadmapText = "# Roadmap\n"
	s.Sequence = "## Recommended sequence\n\n1. A\n"
	s.Commitment = commitment.Outlook{State: "eligible", Milestone: "M1", Next: "B", Blocked: []intent.OutcomeBlocker{{Outcome: "A", Reason: "vendor <b>fix</b>\x1b"}}, Command: "bench commitment start --outcome B"}
	out := Render(s)
	card := commitmentCard(out)
	for _, want := range []string{"<dt>state</dt><dd>eligible</dd>", "<dt>next_outcome</dt><dd>B</dd>", "<dt>blocked</dt><dd>A</dd>", "<dt>A</dt><dd>vendor &lt;b&gt;fix&lt;/b&gt;\\u001b</dd>", "<dt>command</dt><dd>bench commitment start --outcome B</dd>"} {
		if !strings.Contains(card, want) {
			t.Errorf("commitment card lacks %q:\n%s", want, card)
		}
	}
	if strings.Contains(card, "<dt>waiting</dt>") || strings.ContainsRune(out, 0x1b) || strings.Contains(out, unapproved) {
		t.Errorf("eligible card shows an empty field, a raw control byte, or an unapproved sequence:\n%s", card)
	}

	s.Commitment = commitment.Project(nil, intent.CommitmentState{})
	if out := Render(s); !strings.Contains(out, unapproved) || !strings.Contains(commitmentCard(out), "<dd>adoption-required</dd>") {
		t.Errorf("absent commitment did not mark the sequence as unapproved input:\n%s", commitmentCard(out))
	}
}

// commitmentCard returns the commitment card of a rendered page.
func commitmentCard(out string) string {
	_, card, _ := strings.Cut(out, "<h2>Delivery commitment</h2>")
	card, _, _ = strings.Cut(card, "</section>")
	return card
}
