package preflight

import (
	"path/filepath"

	"github.com/gibbonmi/bench/internal/chargeevidence"
	"github.com/gibbonmi/bench/internal/preflight/chargesource"
	"github.com/gibbonmi/bench/internal/preflight/evidencecmd"
	specref "github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/tickets"
	"github.com/gibbonmi/bench/internal/toon"
)

// buildSourceDescriptor is one required build charge source. The descriptor list is the
// one build source policy: the prepared inventory, the charge columns, and the source
// check all derive from it.
type buildSourceDescriptor struct {
	role string
	path func(specPath string, ticket *tickets.Entry) string
	// ticket marks the selected ticket; checks and returns place the source in those
	// charge columns, in source order.
	ticket, checks, returns bool
}

func fixedSource(path string) func(string, *tickets.Entry) string {
	return func(string, *tickets.Entry) string { return path }
}

func selectedTicketPath(specPath string, ticket *tickets.Entry) string {
	return filepath.ToSlash(filepath.Join(filepath.Dir(specPath), "tickets", ticket.Rel))
}

func buildSourcePolicy() []buildSourceDescriptor {
	return []buildSourceDescriptor{
		{role: "ticket", path: selectedTicketPath, ticket: true, checks: true},
		{role: "spec", path: func(specPath string, _ *tickets.Entry) string { return specPath }},
		{role: "delegate-skill", path: fixedSource(chargesource.DelegateSkill), returns: true},
		{role: "build-phase", path: fixedSource(chargesource.BuildPhase), checks: true},
		{role: "delegate-procedure", path: fixedSource(chargesource.DelegateProcedure), returns: true},
	}
}

// loadBuildSources reads every policy source at the pinned tip, in policy order.
func loadBuildSources(root string, facts Facts, selected *tickets.Entry, policy []buildSourceDescriptor) ([]chargeevidence.SourceInput, string) {
	inputs := make([]chargeevidence.SourceInput, 0, len(policy))
	for _, descriptor := range policy {
		path := descriptor.path(facts.SpecPath, selected)
		data, failure := loadChargeSource(root, facts.SourceTip, path)
		if failure != "" {
			return nil, failure
		}
		inputs = append(inputs, chargeevidence.SourceInput{
			Role: descriptor.role, Kind: chargeevidence.KindRepository, Path: path, Required: true, Data: data,
		})
	}
	return inputs, ""
}

// prepareBuildPack gathers the complete build evidence set and returns it only as a
// strictly validated in-memory pack.
func prepareBuildPack(root string, facts Facts, selected *tickets.Entry, parsed *tickets.Ticket, policy []buildSourceDescriptor) (*chargeevidence.Pack, string, error) {
	inputs, failure := loadBuildSources(root, facts, selected, policy)
	if failure != "" {
		return nil, failure, nil
	}
	metadata := chargeevidence.Metadata{
		Fence:    facts.FenceEntries,
		Writes:   parsed.Writes,
		Coverage: parsed.Covers,
	}
	for i, descriptor := range policy {
		id := chargeevidence.InputSourceID(i)
		if descriptor.ticket {
			metadata.Charge = append(metadata.Charge, chargeevidence.ChargeRow{Ticket: id, Access: chargeevidence.AccessBuild})
		}
		if descriptor.checks {
			metadata.Checks = append(metadata.Checks, id)
		}
		if descriptor.returns {
			metadata.Returns = append(metadata.Returns, id)
		}
	}
	pack, err := chargeevidence.Build(chargeevidence.Candidate{
		Selection: chargeevidence.Selection{
			Mode: modeBuild, Spec: facts.SpecPath, Ticket: selectedTicketPath(facts.SpecPath, selected),
			Base: facts.SourceBase, SourceTip: facts.SourceTip,
		},
		Metadata: metadata,
		Sources:  inputs,
	})
	return pack, "", err
}

// currentEvidenceCommand binds one prepared artifact to the current action. The manifest
// supplies the frozen selectors; the assignment, the checkout state, the source pair, and
// the required bytes come from the current checkout, so a released assignment, a dirty
// checkout, a moved source, or a changed required source refuses. The binding names the
// current assignment, never the assignment that prepared the artifact.
func currentEvidenceCommand(root, identity string, args []string) (string, int) {
	artifact, refusal, code := evidencecmd.OpenEvidence(root, identity)
	if refusal != "" {
		return refusal, code
	}
	defer artifact.Close()
	manifest, _ := artifact.Manifest()
	selection := manifest.Selection
	return preparedAttempts(root, selection.Mode, specref.LiveSpecSlug(selection.Spec), selection.Base, "", "check", args, func(facts Facts) (string, int) {
		if refusal := preparationCheckoutRefusal(root, facts, "check"); refusal != "" {
			return refusal, 1
		}
		if verdict := Decide(facts); verdict.Red {
			return chargeVerdictRefusal(boundedVerdict(verdict)), 1
		}
		if facts.SourceTip != selection.SourceTip {
			return chargeRefusal("source", "the current source tip "+facts.SourceTip+" is not the prepared source tip "+selection.SourceTip,
				"prepare evidence for the current source and rerun the exact check"), 1
		}
		if refusal := currentSourcesRefusal(root, facts, manifest); refusal != "" {
			return refusal, 1
		}
		text, err := chargeevidence.Current{Evidence: identity, Assignment: facts.AssignmentTarget,
			Base: selection.Base, SourceTip: facts.SourceTip}.Encode()
		if err != nil {
			return toon.RenderError(err) + "\n", 1
		}
		return text, 0
	})
}

// currentSourcesRefusal compares every prepared repository source with the bytes the
// current source tip holds, so unchanged Git pins alone cannot authorize action.
func currentSourcesRefusal(root string, facts Facts, manifest chargeevidence.Manifest) string {
	for _, source := range manifest.Sources {
		if source.Kind != chargeevidence.KindRepository {
			continue
		}
		data, failure := loadChargeSource(root, facts.SourceTip, source.Path)
		if failure != "" {
			return chargeRefusal("source", failure, "restore the named canonical source and rerun the exact check")
		}
		if chargeevidence.Digest(data) != source.SHA256 {
			return chargeRefusal("source", source.Path+" differs from the prepared evidence",
				"prepare evidence for the current sources and rerun the exact check")
		}
	}
	return ""
}
