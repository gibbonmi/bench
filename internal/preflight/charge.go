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
func chargeCommand(
	root, mode, slug, base, sourceTip, ticket string,
	full bool,
	version string,
	args []string,
) (string, int) {
	return preparedCommand(root, mode, slug, base, sourceTip, ticket, full, version, args, chargePreparation)
}

// chargePacket is the one source for every charge tail. Build and review differ only in
// their charge columns and in the table between the charge and the sources listing. The
// complete/next flip, the sources listing, and the full-or-omitted projection are the
// same fact in both modes, so they are authored once here.
type chargePacket struct {
	// fields and rows exclude the trailing complete and next columns; the renderer owns
	// that pair because it owns the flip that fills them.
	fields []string
	rows   [][]string

	// middle holds already-rendered tables that sit between the charge and the sources
	// listing: the build coverage rows, or the review shared-evidence identities.
	middle []string

	// sources are the frozen inputs the packet lists, retrieves under --full, and names
	// as omitted otherwise. identitiesOnly are listed but never retrieved: the review
	// packet's derived shared-evidence handle is one.
	sources        []chargeSource
	identitiesOnly []chargeSource

	// next is the exact retrieval invocation a compact packet advertises.
	next string
}

func renderChargePacket(packet chargePacket, full bool) (string, int) {
	complete, next := "false", packet.next
	if full {
		complete, next = "true", ""
	}
	chargeRows := make([][]string, len(packet.rows))
	for i, row := range packet.rows {
		chargeRows[i] = append(append([]string{}, row...), complete, next)
	}
	charge, err := toon.Table("charge", append(append([]string{}, packet.fields...), "complete", "next"), chargeRows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	listed := append(append([]chargeSource{}, packet.sources...), packet.identitiesOnly...)
	sourceRows := make([][]string, len(listed))
	for i, source := range listed {
		sourceRows[i] = []string{source.path, source.identity()}
	}
	identities, err := toon.Table("sources", []string{"path", "identity"}, sourceRows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	var b strings.Builder
	b.WriteString(charge)
	for _, table := range packet.middle {
		b.WriteString(table)
	}
	b.WriteString(identities)
	if full {
		evidenceRows := make([][]string, len(packet.sources))
		for i, source := range packet.sources {
			evidenceRows[i] = []string{source.path, string(source.data)}
		}
		evidence, err := toon.Table("evidence", []string{"source", "content"}, evidenceRows)
		if err != nil {
			return toon.RenderError(err) + "\n", 1
		}
		b.WriteString(evidence)
		return b.String(), 0
	}
	omittedRows := make([][]string, len(packet.sources))
	for i, source := range packet.sources {
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

// chargeFenceCell is the spec's declared ownership fence, the one fact the fence column
// carries. It is deliberately not a source handle: the ticket and evidence columns
// already carry handles, and a column that repeats one of them grades nothing.
func chargeFenceCell(facts Facts) string {
	return strings.Join(facts.FenceEntries, ", ")
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
	coverage, err := toon.Table("coverage", []string{"row"}, rows(parsed.Covers))
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return renderChargePacket(chargePacket{
		fields: []string{"assignment", "checkout", "base", "source_tip", "fence", "ticket", "writes", "evidence", "checks", "return"},
		rows: [][]string{{
			facts.AssignmentTarget, root, facts.SourceBase, facts.SourceTip,
			chargeFenceCell(facts), sources.ticket.handle(), strings.Join(parsed.Writes, ", "),
			sourceHandles(sources.list()...),
			sourceHandles(sources.ticket, sources.buildPhase),
			sourceHandles(sources.delegateSkill, sources.delegateProcedure),
		}},
		middle:  []string{coverage},
		sources: sources.list(),
		next:    chargeInvocation(modeBuild, facts, name),
	}, full)
}

// buildChargeSourceSet names each frozen build source. Every charge column reads a field
// name, so a reorder of the load list below cannot silently reassign a column.
type buildChargeSourceSet struct {
	ticket, spec, delegateSkill, buildPhase, delegateProcedure chargeSource
}

func (set buildChargeSourceSet) list() []chargeSource {
	return []chargeSource{set.ticket, set.spec, set.delegateSkill, set.buildPhase, set.delegateProcedure}
}

func chargeSources(root, sourceTip, specPath string, selected *tickets.Entry) (buildChargeSourceSet, string) {
	var set buildChargeSourceSet
	failure := loadChargeSources(root, sourceTip, []namedChargeSource{
		{filepath.ToSlash(filepath.Join(filepath.Dir(specPath), "tickets", selected.Rel)), &set.ticket},
		{specPath, &set.spec},
		{delegateSkill, &set.delegateSkill},
		{buildPhase, &set.buildPhase},
		{delegateProcedure, &set.delegateProcedure},
	})
	if failure != "" {
		return buildChargeSourceSet{}, failure
	}
	return set, ""
}

// namedChargeSource binds one canonical path to the field that holds it. The binding is
// one literal, so a path and its meaning never drift apart.
type namedChargeSource struct {
	path string
	into *chargeSource
}

func loadChargeSources(root, sourceTip string, named []namedChargeSource) string {
	for _, item := range named {
		data, failure := readChargeSource(root, sourceTip, item.path)
		if failure != "" {
			return failure
		}
		if !toon.Representable(string(data)) {
			return item.path + " contains a byte spec-TOON cannot represent"
		}
		*item.into = chargeSource{path: item.path, data: data}
	}
	return ""
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

// chargeInvocation is the exact full-retrieval command a compact packet advertises. An
// empty name is the review form, which takes no ticket.
func chargeInvocation(mode string, facts Facts, name string) string {
	args := []string{"bench", "preflight", mode, facts.SpecPath, "--charge"}
	if name != "" {
		args = append(args, "--ticket", name)
	}
	args = append(args, "--base", facts.SourceBase, "--source-tip", facts.SourceTip, "--full")
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
