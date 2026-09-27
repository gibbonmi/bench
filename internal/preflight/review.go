package preflight

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/reviewrecord"

	"github.com/gibbonmi/bench/internal/chargeevidence"
	"github.com/gibbonmi/bench/internal/consumers"
	"github.com/gibbonmi/bench/internal/coverage"
	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/preflight/chargesource"
	toonlib "github.com/toon-format/toon-go"
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

// reviewCollector is one collector's generated review evidence. It carries the fragments
// the collector returned and the exact invocation that produced them, so the manifest commits
// the provenance beside each digest. A changed producer or a changed argument therefore
// changes the evidence identity even when the captured bytes are identical.
type reviewCollector struct {
	kind      string
	producer  string
	arguments []string
	fragments []reviewFragment
}

// reviewFragment is one generated source of a collector: its declared role and path, and its
// exact bytes. The fragments of one collector rebuild its complete output in order.
type reviewFragment struct {
	role, path string
	data       []byte
}

// diffFragments is the diff collector's source table: the existing prefix, each file patch
// under its identifying path, then the existing suffix, in reconstruction order.
func diffFragments(snapshot diff.PatchSnapshot) []reviewFragment {
	fragments := []reviewFragment{{role: "diff-prefix", path: "diff-prefix", data: snapshot.Prefix}}
	for _, patch := range snapshot.Patches {
		fragments = append(fragments, reviewFragment{role: "diff-file", path: patch.Path, data: patch.Body})
	}
	return append(fragments, reviewFragment{role: "diff-suffix", path: "diff-suffix", data: snapshot.Suffix})
}

// reviewCollectorPolicy is the one review collector inventory. Each entry names the kind
// the metadata binds, the producer the manifest records, and the command that runs.
type reviewCollectorPolicy struct {
	kind, producer string
	arguments      func(facts Facts) []string
	// run returns the collector's fragments, or nil fragments and its whole output.
	run func(root, version string, args []string) ([]reviewFragment, string, int)
}

func reviewCollectorPolicies() []reviewCollectorPolicy {
	return []reviewCollectorPolicy{
		{
			kind: "diff", producer: "bench diff",
			arguments: func(facts Facts) []string {
				return []string{"--base", facts.SourceBase, "--source-tip", facts.SourceTip, "--full"}
			},
			run: func(_, _ string, args []string) ([]reviewFragment, string, int) {
				snapshot, out, code := diff.PairPatches(args)
				if code != 0 {
					return nil, out, code
				}
				return diffFragments(snapshot), "", 0
			},
		},
		{
			kind: "consumers", producer: "bench consumers",
			arguments: func(facts Facts) []string {
				return []string{"--changed", "--base", facts.SourceBase, "--source-tip", facts.SourceTip, "--full"}
			},
			run: func(_, version string, args []string) ([]reviewFragment, string, int) {
				out, code := consumers.CommandWithVersion(version)(args)
				return nil, out, code
			},
		},
		{
			// The coverage collector resolves a slashed argument against the caller's working
			// directory, not against the repository root. The charge anchors the spec path at
			// the root, so the capture stays identical from any working directory. Coverage
			// renders the result repo-relative again, so the anchored argument changes no
			// output byte.
			kind: "coverage", producer: "bench coverage",
			arguments: func(facts Facts) []string { return []string{facts.SpecPath} },
			run: func(root, _ string, args []string) ([]reviewFragment, string, int) {
				anchored := []string{filepath.Join(root, filepath.FromSlash(args[0]))}
				out, code := coverage.Command(anchored)
				return nil, out, code
			},
		},
	}
}

// collectReviewEvidence runs every registered collector once, in policy order. It runs only
// inside a preparation attempt: a read serves published bytes from the store and reaches no
// collector. A failed or incomplete collector stops the attempt, so no later collector's
// capture reaches a published artifact.
func collectReviewEvidence(root, version string, facts Facts, observe reviewEvidenceObserver) ([]reviewCollector, string) {
	policies := reviewCollectorPolicies()
	collected := make([]reviewCollector, 0, len(policies))
	for _, policy := range policies {
		arguments := policy.arguments(facts)
		observe(policy.kind)
		fragments, out, code := policy.run(root, version, arguments)
		if code != 0 {
			return nil, policy.kind + " evidence failed: " + strings.TrimSpace(out)
		}
		if failure := reviewCaptureRefusal(policy.kind, out); failure != "" {
			return nil, failure
		}
		if fragments == nil {
			// A collector that returns no fragments keeps its whole output as one source.
			fragments = []reviewFragment{{role: policy.kind, path: policy.kind, data: []byte(out)}}
		}
		collected = append(collected, reviewCollector{
			kind: policy.kind, producer: policy.producer, arguments: arguments, fragments: fragments,
		})
	}
	return collected, ""
}

// reviewCaptureRefusal reports one collector's own completeness claim. A collector that
// declares truncated output has not captured the frozen pair, so the attempt refuses rather
// than freezing a partial capture.
func reviewCaptureRefusal(kind, output string) string {
	if kind != "consumers" {
		return ""
	}
	complete, err := completeConsumerEvidence(output)
	if err != nil {
		return "consumer evidence metadata failed: " + err.Error()
	}
	if !complete {
		return "consumer evidence is incomplete: full collection omitted unrepresentable rows"
	}
	return ""
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

// reviewSourceDescriptor is one required review repository source. The descriptor list is
// the one review source policy: the prepared inventory and the metadata check and return
// columns all derive from it.
type reviewSourceDescriptor struct {
	role string
	path func(specPath string) string
	// checks and returns place the source in those metadata columns, in source order.
	checks, returns bool
}

func fixedReviewSource(path string) func(string) string {
	return func(string) string { return path }
}

func reviewSourcePolicy() []reviewSourceDescriptor {
	return []reviewSourceDescriptor{
		{role: "spec", path: func(specPath string) string { return specPath }},
		{role: "review-skill", path: fixedReviewSource(chargesource.ReviewSkill), checks: true},
		{role: "review-phase", path: fixedReviewSource(chargesource.ReviewPhase), returns: true},
		{role: "delegate-skill", path: fixedReviewSource(chargesource.DelegateSkill), returns: true},
		{role: "delegate-procedure", path: fixedReviewSource(chargesource.DelegateProcedure), returns: true},
	}
}

// reviewChargePack applies every review preparation refusal in its fixed order and returns
// the validated in-memory pack, or the refusal that stopped it. The collectors run last, so
// a dirty checkout, a red verdict, or a missing source refuses before any collector starts.
func reviewChargePack(root, version string, facts Facts, verdict Verdict) (*chargeevidence.Pack, string) {
	if refusal := preparationCheckoutRefusal(root, facts, "charge"); refusal != "" {
		return nil, refusal
	}
	if verdict.Red {
		return nil, chargeVerdictRefusal(verdict)
	}
	policy := reviewSourcePolicy()
	inputs := make([]chargeevidence.SourceInput, 0, len(policy))
	for _, descriptor := range policy {
		path := descriptor.path(facts.SpecPath)
		data, failure := loadChargeSource(root, facts.SourceTip, path)
		if failure != "" {
			return nil, chargeRefusal("source", failure, "restore the named canonical source and rerun the exact charge")
		}
		inputs = append(inputs, chargeevidence.SourceInput{
			Role: descriptor.role, Kind: chargeevidence.KindRepository, Path: path, Required: true, Data: data,
		})
	}
	collected, failure := collectReviewEvidence(root, version, facts, reviewEvidenceObserved)
	if failure != "" {
		return nil, chargeRefusal("evidence", failure, "repair the named collector input and rerun the exact charge")
	}
	for _, item := range collected {
		for _, fragment := range item.fragments {
			inputs = append(inputs, chargeevidence.SourceInput{
				Role: fragment.role, Kind: chargeevidence.KindGenerated, Path: fragment.path, Required: true, Data: fragment.data,
				Producer: &chargeevidence.Producer{
					Name: item.producer, Version: version, Cwd: root, Arguments: item.arguments,
				},
			})
		}
	}
	metadata, err := reviewMetadata(root, facts, policy, collected)
	if err != nil {
		return nil, chargeRefusal("completion-evidence", err.Error(), "repair the record path and retry")
	}
	pack, err := chargeevidence.Build(chargeevidence.Candidate{
		Selection: chargeevidence.Selection{
			Mode: modeReview, Spec: facts.SpecPath, Base: facts.SourceBase, SourceTip: facts.SourceTip,
		},
		Metadata: metadata,
		Sources:  inputs,
	})
	if err != nil {
		return nil, chargeRefusal("evidence", err.Error(), "repair the reported evidence condition and rerun the exact charge")
	}
	return pack, ""
}

// reviewMetadata binds the complete review facts to the prepared sources. Every axis reads
// one artifact, so the charge block names each axis against the same frozen captures.
func reviewMetadata(root string, facts Facts, policy []reviewSourceDescriptor, collected []reviewCollector) (chargeevidence.Metadata, error) {
	metadata := chargeevidence.Metadata{Fence: facts.FenceEntries}
	for _, axis := range reviewrecord.Axes() {
		metadata.Charge = append(metadata.Charge, chargeevidence.ChargeRow{Axis: axis, Access: chargeevidence.AccessReview})
	}
	for i, descriptor := range policy {
		id := chargeevidence.InputSourceID(i)
		if descriptor.checks {
			metadata.Checks = append(metadata.Checks, id)
		}
		if descriptor.returns {
			metadata.Returns = append(metadata.Returns, id)
		}
	}
	// Each fragment's shared row names its collector kind, in the same order the inputs
	// declare the fragments, so a split collector binds every fragment in reconstruction order.
	next := len(policy)
	for _, item := range collected {
		for range item.fragments {
			metadata.Shared = append(metadata.Shared, chargeevidence.SharedRow{
				Kind: item.kind, Source: chargeevidence.InputSourceID(next),
			})
			next++
		}
	}
	completion, err := completionEvidenceRow(root, facts)
	if err != nil {
		return chargeevidence.Metadata{}, err
	}
	metadata.Completion = append(metadata.Completion, completion)
	return metadata, nil
}

func completionEvidenceRow(root string, facts Facts) (chargeevidence.CompletionRow, error) {
	path, err := reviewrecord.RecordPath(facts.SpecPath)
	if err != nil {
		return chargeevidence.CompletionRow{}, err
	}
	tree, err := benchgit.Output("-C", root, "rev-parse", "--verify", facts.SourceTip+"^{tree}")
	if err != nil {
		return chargeevidence.CompletionRow{}, err
	}
	source, err := reviewrecord.SourceDigest(root, tree, facts.SpecPath)
	if err != nil {
		return chargeevidence.CompletionRow{}, err
	}
	state, detail := "parsed", ""
	_, readErr := reviewrecord.Read(root, facts.SpecPath)
	if errors.Is(readErr, reviewrecord.ErrMissing) {
		state = "missing"
	} else if readErr != nil {
		state, detail = "invalid", readErr.Error()
	}
	// The gatherer already read the plan for the completion-plan row. The metadata
	// records that one answer rather than parsing the fence a second time.
	if facts.CompletionPlanError != "" {
		detail = "completion plan unavailable: " + facts.CompletionPlanError
	}
	return chargeevidence.CompletionRow{
		Record: path, SourceDigest: source, PlanDigest: facts.CompletionPlanDigest,
		RecordState: state, Detail: detail,
	}, nil
}
