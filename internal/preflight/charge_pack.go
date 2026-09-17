package preflight

import (
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/chargeevidence"
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
		{role: "delegate-skill", path: fixedSource(delegateSkill), returns: true},
		{role: "build-phase", path: fixedSource(buildPhase), checks: true},
		{role: "delegate-procedure", path: fixedSource(delegateProcedure), returns: true},
	}
}

const buildAccess = "write-within-fence"

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
			metadata.Charge = append(metadata.Charge, chargeevidence.ChargeRow{Ticket: id, Access: buildAccess})
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

// renderLegacyBuildPacket projects the validated pack onto the legacy build charge. It
// reads every source and column from the decoded pack, never from the policy.
func renderLegacyBuildPacket(root string, facts Facts, pack *chargeevidence.Pack, name string, full bool) (string, int) {
	manifest, metadata := pack.Manifest(), pack.Metadata()
	byID := map[string]chargeSource{}
	var listed []chargeSource
	for _, source := range manifest.Sources {
		if source.Kind != chargeevidence.KindRepository {
			continue
		}
		data, _ := pack.Source(source.ID)
		byID[source.ID] = chargeSource{path: source.Path, data: data}
		listed = append(listed, byID[source.ID])
	}
	handlesOf := func(ids []string) string {
		sources := make([]chargeSource, len(ids))
		for i, id := range ids {
			sources[i] = byID[id]
		}
		return sourceHandles(sources...)
	}
	coverage, err := toon.Table("coverage", []string{"row"}, rows(metadata.Coverage))
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	selection := manifest.Selection
	return renderChargePacket(chargePacket{
		fields: []string{"assignment", "checkout", "base", "source_tip", "fence", "ticket", "writes", "evidence", "checks", "return"},
		rows: [][]string{{
			facts.AssignmentTarget, root, selection.Base, selection.SourceTip,
			chargeFenceCell(metadata.Fence), handlesOf([]string{metadata.Charge[0].Ticket}),
			strings.Join(metadata.Writes, ", "), sourceHandles(listed...),
			handlesOf(metadata.Checks), handlesOf(metadata.Returns),
		}},
		middle:  []string{coverage},
		sources: listed,
		next:    chargeInvocation(modeBuild, facts, name),
	}, full)
}
