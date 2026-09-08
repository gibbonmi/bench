package preflight

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/tickets"
	"github.com/gibbonmi/bench/internal/toon"
)

const (
	delegateSkill     = ".agents/skills/bench-craft-delegate/SKILL.md"
	delegateProcedure = ".agents/skills/bench-craft-delegate/references/delegation-discipline.md"
	buildPhase        = ".agents/commands/bench-implement-spec.md"
)

type chargeSource struct {
	path string
	data []byte
}

func (source chargeSource) identity() string {
	sum := sha256.Sum256(source.data)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func (source chargeSource) handle() string {
	return source.path + " " + source.identity()
}

// chargeCommand keeps every required read in the source snapshot's one retry. A charge
// must not turn GatherPinned's retry into nested retries with mixed source evidence.
func chargeCommand(root, mode, slug, base, sourceTip, ticket string, full bool, args []string) (string, int) {
	return preparedCommand(root, mode, slug, base, sourceTip, ticket, full, args, chargePreparation)
}

func renderCharge(root string, facts Facts, verdict Verdict, name string, full bool) (string, int) {
	if refusal := preparationCheckoutRefusal(root, facts, "charge"); refusal != "" {
		return refusal, 1
	}
	if verdict.Red {
		return chargeVerdictRefusal(verdict), 1
	}
	selected, parsed, detail, selectionNext := preparationTicket(root, facts, name)
	if detail != "" {
		return chargeRefusal("ticket", detail, selectionNext), 1
	}
	sources, failure := chargeSources(root, facts.SourceTip, facts.SpecPath, selected)
	if failure != "" {
		return chargeRefusal("source", failure, "restore the named canonical source and rerun the exact charge"), 1
	}
	next := chargeInvocation(facts, name)
	complete := "false"
	if full {
		complete, next = "true", ""
	}
	charge, err := toon.Table("charge", []string{"assignment", "checkout", "base", "source_tip", "fence", "ticket", "writes", "evidence", "checks", "return", "complete", "next"}, [][]string{{
		facts.AssignmentTarget, root, facts.SourceBase, facts.SourceTip,
		sources[1].handle(), sources[0].handle(), strings.Join(parsed.Writes, ", "),
		sources[1].handle(), sourceHandles(sources[0], sources[3]), sourceHandles(sources[2], sources[4]), complete, next,
	}})
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	coverage, err := toon.Table("coverage", []string{"row"}, rows(parsed.Covers))
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	sourceRows := make([][]string, len(sources))
	for i, source := range sources {
		sourceRows[i] = []string{source.path, source.identity()}
	}
	identities, err := toon.Table("sources", []string{"path", "identity"}, sourceRows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	var b strings.Builder
	b.WriteString(charge)
	b.WriteString(coverage)
	b.WriteString(identities)
	if full {
		evidenceRows := make([][]string, len(sources))
		for i, source := range sources {
			evidenceRows[i] = []string{source.path, string(source.data)}
		}
		evidence, err := toon.Table("evidence", []string{"path", "content"}, evidenceRows)
		if err != nil {
			return toon.RenderError(err) + "\n", 1
		}
		b.WriteString(evidence)
		return b.String(), 0
	}
	omittedRows := make([][]string, len(sources))
	for i, source := range sources {
		omittedRows[i] = []string{source.path}
	}
	omitted, err := toon.Table("omitted", []string{"source"}, omittedRows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	help, err := toon.Table("help", []string{"cmd", "why"}, [][]string{{next, "retrieve every omitted source before dispatch"}})
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	b.WriteString(omitted)
	b.WriteString(help)
	return b.String(), 0
}

func chargeSources(root, sourceTip, specPath string, selected *tickets.Entry) ([]chargeSource, string) {
	paths := []string{
		filepath.ToSlash(filepath.Join(filepath.Dir(specPath), "tickets", selected.Rel)),
		specPath,
		delegateSkill,
		buildPhase,
		delegateProcedure,
	}
	sources := make([]chargeSource, 0, len(paths))
	for _, path := range paths {
		data, failure := readChargeSource(root, sourceTip, path)
		if failure != "" {
			return nil, failure
		}
		if !toon.Representable(string(data)) {
			return nil, path + " contains a byte spec-TOON cannot represent"
		}
		sources = append(sources, chargeSource{path: path, data: data})
	}
	return sources, ""
}

func selectedTicket(entries []tickets.Entry, name string) *tickets.Entry {
	for i := range entries {
		if entries[i].Name == name {
			return &entries[i]
		}
	}
	return nil
}

func selectedParsedTicket(facts Facts, name string) *tickets.Ticket {
	for i := range facts.Tickets {
		if facts.Tickets[i].Name == name {
			return &facts.Tickets[i]
		}
	}
	return nil
}

func readChargeSource(root, sourceTip, rel string) ([]byte, string) {
	read := bounds.ClassifyNoFollow(filepath.Join(root, filepath.FromSlash(rel)))
	if read.State != bounds.StateParsed {
		return nil, rel + " is " + string(read.State) + ": " + read.Reason
	}
	pinned, err := git.Raw("-C", root, "show", sourceTip+":"+rel)
	if err != nil {
		return nil, rel + " is absent or unreadable at source tip " + sourceTip
	}
	if !bytes.Equal(read.Data, pinned) {
		return nil, rel + " does not match source tip " + sourceTip
	}
	return pinned, ""
}

func sourceHandles(sources ...chargeSource) string {
	handles := make([]string, len(sources))
	for i, source := range sources {
		handles[i] = source.handle()
	}
	return strings.Join(handles, "; ")
}

func chargeInvocation(facts Facts, name string) string {
	args := []string{"bench", "preflight", "build", facts.SpecPath, "--charge", "--ticket", name, "--base", facts.SourceBase, "--source-tip", facts.SourceTip, "--full"}
	for i := range args {
		args[i] = axi.ShellQuote(args[i])
	}
	return strings.Join(args, " ")
}

func chargeVerdictRefusal(verdict Verdict) string {
	for _, check := range verdict.Checks {
		if check.Verdict == verdictRed {
			next := check.Next
			if next == "" {
				next = "repair " + check.Check + " and rerun the exact charge"
			}
			return chargeRefusal("preflight", check.Check+": "+check.Detail, next)
		}
	}
	return chargeRefusal("preflight", "required preflight checks are red", "repair the reported check and rerun the exact charge")
}

func rows(values []string) [][]string {
	result := make([][]string, len(values))
	for i, value := range values {
		result[i] = []string{value}
	}
	return result
}

func chargeRefusal(input, detail, next string) string {
	return toon.Errorf(input+" required: "+detail, next) + "\n"
}
