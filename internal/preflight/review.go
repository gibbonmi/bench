package preflight

import (
	"fmt"
	"strings"

	"github.com/gibbonmi/bench/internal/axi"
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

func collectReviewEvidence(version string, facts Facts, observe reviewEvidenceObserver) (reviewEvidence, string) {
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

	observe("coverage")
	coverageOut, code := coverage.Command([]string{facts.SpecPath})
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
	evidence, failure := collectReviewEvidence(version, facts, reviewEvidenceObserved)
	if failure != "" {
		return chargeRefusal("evidence", failure, "repair the named collector input and rerun the exact charge"), 1
	}
	return renderReviewPacket(root, facts, sources, evidence, full)
}

func reviewChargeSources(root, sourceTip, specPath string) ([]chargeSource, string) {
	paths := []string{specPath, reviewSkill, reviewPhase, delegateSkill, delegateProcedure}
	return loadChargeSources(root, sourceTip, paths)
}

func renderReviewPacket(
	root string,
	facts Facts,
	sources []chargeSource,
	evidence reviewEvidence,
	full bool,
) (string, int) {
	shared := evidence.shared()
	next := reviewChargeInvocation(facts)
	complete := "false"
	if full {
		complete, next = "true", ""
	}
	chargeRows := make([][]string, 0, 3)
	for _, axis := range []string{"Standards", "Spec", "Coverage"} {
		chargeRows = append(chargeRows, []string{
			axis, facts.AssignmentTarget, root, facts.SourceBase, facts.SourceTip,
			sources[0].handle(), sources[0].handle(), "read-only", shared.handle(),
			sources[1].handle(), sourceHandles(sources[2], sources[3], sources[4]), complete, next,
		})
	}
	chargeFields := []string{
		"axis", "assignment", "checkout", "base", "source_tip", "fence", "ticket",
		"writes", "evidence", "checks", "return", "complete", "next",
	}
	charge, err := toon.Table("charge", chargeFields, chargeRows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	sharedRows := make([][]string, 0, 3)
	for _, item := range evidence.sources() {
		sharedRows = append(sharedRows, []string{item.path, item.identity()})
	}
	sharedTable, err := toon.Table("shared_evidence", []string{"kind", "identity"}, sharedRows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	all := append(append([]chargeSource{}, sources...), evidence.sources()...)
	sourceRows := make([][]string, 0, len(all)+1)
	for _, source := range all {
		sourceRows = append(sourceRows, []string{source.path, source.identity()})
	}
	sourceRows = append(sourceRows, []string{shared.path, shared.identity()})
	identities, err := toon.Table("sources", []string{"path", "identity"}, sourceRows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}

	var b strings.Builder
	b.WriteString(charge)
	b.WriteString(sharedTable)
	b.WriteString(identities)
	if full {
		rows := make([][]string, len(all))
		for i, source := range all {
			rows[i] = []string{source.path, string(source.data)}
		}
		table, err := toon.Table("evidence", []string{"source", "content"}, rows)
		if err != nil {
			return toon.RenderError(err) + "\n", 1
		}
		b.WriteString(table)
		return b.String(), 0
	}
	omitted := make([][]string, len(all))
	for i, source := range all {
		omitted[i] = []string{source.path}
	}
	omittedTable, err := toon.Table("omitted", []string{"source"}, omitted)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	help, err := toon.Table("help", []string{"cmd", "why"}, [][]string{{next, "retrieve every omitted source before dispatch"}})
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	b.WriteString(omittedTable)
	b.WriteString(help)
	return b.String(), 0
}

func reviewChargeInvocation(facts Facts) string {
	args := []string{
		"bench", "preflight", "review", facts.SpecPath, "--charge", "--base",
		facts.SourceBase, "--source-tip", facts.SourceTip, "--full",
	}
	for i := range args {
		args[i] = axi.ShellQuote(args[i])
	}
	return strings.Join(args, " ")
}
