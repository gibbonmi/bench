package preflight

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/consumers"
	"github.com/gibbonmi/bench/internal/coverage"
	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/toon"
	toonlib "github.com/toon-format/toon-go"
)

const (
	reviewSkill = ".agents/skills/bench-craft-review/SKILL.md"
	reviewPhase = ".agents/commands/bench-review-implementation.md"
)

// reviewEvidenceObserver reports each real collector as it starts. Tests use this port
// to prove that one preparation attempt does not repeat collection for each review axis.
type reviewEvidenceObserver func(string)

var reviewEvidenceObserved reviewEvidenceObserver = func(string) {}

func setReviewEvidenceObserverForTest(observer reviewEvidenceObserver) func() {
	previous := reviewEvidenceObserved
	reviewEvidenceObserved = observer
	return func() { reviewEvidenceObserved = previous }
}

type reviewEvidence struct {
	diff, consumers, coverage chargeSource
}

func (e reviewEvidence) sources() []chargeSource {
	return []chargeSource{e.diff, e.consumers, e.coverage}
}

func (e reviewEvidence) shared() chargeSource {
	return chargeSource{
		path: "shared-review-evidence",
		data: append(append(append([]byte{}, e.diff.data...), e.consumers.data...), e.coverage.data...),
	}
}

func collectReviewEvidence(root, version string, facts Facts, observe reviewEvidenceObserver) (reviewEvidence, string) {
	diffArgs := []string{"--base", facts.SourceBase, "--source-tip", facts.SourceTip, "--full"}
	observe("diff")
	diffOut, code := diff.Command(diffArgs)
	if code != 0 {
		return reviewEvidence{}, "diff evidence failed: " + strings.TrimSpace(diffOut)
	}

	consumerArgs := []string{"--changed", "--base", facts.SourceBase, "--source-tip", facts.SourceTip, "--full"}
	observe("consumers")
	consumerOut, code := consumers.CommandWithVersion(version)(consumerArgs)
	if code != 0 {
		return reviewEvidence{}, "consumer evidence failed: " + strings.TrimSpace(consumerOut)
	}
	complete, err := completeConsumerEvidence(consumerOut)
	if err != nil {
		return reviewEvidence{}, "consumer evidence metadata failed: " + err.Error()
	}
	if !complete {
		return reviewEvidence{}, "consumer evidence is incomplete: full collection omitted unrepresentable rows"
	}

	// The coverage collector resolves a slashed argument against the caller's working
	// directory, not against the repository root. The charge anchors the spec path at the
	// root, so the packet stays identical from any working directory. Coverage renders the
	// result repo-relative again, so the anchored argument changes no output byte.
	observe("coverage")
	coverageOut, code := coverage.Command([]string{filepath.Join(root, filepath.FromSlash(facts.SpecPath))})
	if code != 0 {
		return reviewEvidence{}, "coverage evidence failed: " + strings.TrimSpace(coverageOut)
	}
	return reviewEvidence{
		diff:      chargeSource{path: "diff", data: []byte(diffOut)},
		consumers: chargeSource{path: "consumers", data: []byte(consumerOut)},
		coverage:  chargeSource{path: "coverage", data: []byte(coverageOut)},
	}, ""
}

func completeConsumerEvidence(output string) (bool, error) {
	decoded, err := toonlib.DecodeString(output)
	if err != nil {
		return false, err
	}
	document, ok := decoded.(map[string]any)
	if !ok {
		return false, fmt.Errorf("decoded as %T, want an object", decoded)
	}
	rows, ok := document["meta"].([]any)
	if !ok || len(rows) != 1 {
		return false, fmt.Errorf("meta decoded as %T with %d rows, want one row", document["meta"], len(rows))
	}
	row, ok := rows[0].(map[string]any)
	if !ok {
		return false, fmt.Errorf("meta row decoded as %T, want an object", rows[0])
	}
	truncated, ok := row["truncated"].(bool)
	if !ok {
		return false, fmt.Errorf("meta.truncated decoded as %T, want a boolean", row["truncated"])
	}
	return !truncated, nil
}

func renderReviewCharge(root string, facts Facts, verdict Verdict, full bool, version string) (string, int) {
	if refusal := preparationCheckoutRefusal(root, facts, "charge"); refusal != "" {
		return refusal, 1
	}
	if verdict.Red {
		return chargeVerdictRefusal(verdict), 1
	}
	sources, failure := reviewChargeSources(root, facts.SourceTip, facts.SpecPath)
	if failure != "" {
		return chargeRefusal("source", failure, "restore the named canonical source and rerun the exact charge"), 1
	}
	evidence, failure := collectReviewEvidence(root, version, facts, reviewEvidenceObserved)
	if failure != "" {
		return chargeRefusal("evidence", failure, "repair the named collector input and rerun the exact charge"), 1
	}
	return renderReviewPacket(root, facts, sources, evidence, full)
}

// reviewChargeSourceSet names each frozen review source, for the reason
// buildChargeSourceSet does: a column reads a field name, never a list position.
type reviewChargeSourceSet struct {
	spec, reviewSkill, reviewPhase, delegateSkill, delegateProcedure chargeSource
}

func (set reviewChargeSourceSet) list() []chargeSource {
	return []chargeSource{set.spec, set.reviewSkill, set.reviewPhase, set.delegateSkill, set.delegateProcedure}
}

func reviewChargeSources(root, sourceTip, specPath string) (reviewChargeSourceSet, string) {
	var set reviewChargeSourceSet
	failure := loadChargeSources(root, sourceTip, []namedChargeSource{
		{specPath, &set.spec},
		{reviewSkill, &set.reviewSkill},
		{reviewPhase, &set.reviewPhase},
		{delegateSkill, &set.delegateSkill},
		{delegateProcedure, &set.delegateProcedure},
	})
	if failure != "" {
		return reviewChargeSourceSet{}, failure
	}
	return set, ""
}

func renderReviewPacket(
	root string,
	facts Facts,
	sources reviewChargeSourceSet,
	evidence reviewEvidence,
	full bool,
) (string, int) {
	shared := evidence.shared()
	chargeRows := make([][]string, 0, 3)
	for _, axis := range []string{"Standards", "Spec", "Coverage"} {
		chargeRows = append(chargeRows, []string{
			axis, facts.AssignmentTarget, root, facts.SourceBase, facts.SourceTip,
			chargeFenceCell(facts), sources.spec.handle(), "read-only", shared.handle(),
			sources.reviewSkill.handle(),
			sourceHandles(sources.reviewPhase, sources.delegateSkill, sources.delegateProcedure),
		})
	}
	sharedRows := make([][]string, 0, 3)
	for _, item := range evidence.sources() {
		sharedRows = append(sharedRows, []string{item.path, item.identity()})
	}
	sharedTable, err := toon.Table("shared_evidence", []string{"kind", "identity"}, sharedRows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return renderChargePacket(chargePacket{
		fields: []string{
			"axis", "assignment", "checkout", "base", "source_tip", "fence", "ticket",
			"writes", "evidence", "checks", "return",
		},
		rows:           chargeRows,
		middle:         []string{sharedTable},
		sources:        append(append([]chargeSource{}, sources.list()...), evidence.sources()...),
		identitiesOnly: []chargeSource{shared},
		next:           chargeInvocation("review", facts, ""),
	}, full)
}
