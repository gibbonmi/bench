package ownership

import (
	"errors"
	"path"
	"path/filepath"
	"slices"
	"sort"

	"github.com/gibbonmi/bench/internal/anchors"
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/tickets"
)

// kitRegistryID names the diagnostic of a compiled registry location that no verified kit
// source binds.
const kitRegistryID = "kit-registry"

// inventory is every canonical owner fact one collection reads, loaded once.
type inventory struct {
	source Source
	pins   map[string][]canary.Pin
	refs   []anchors.Reference
	// lines locates each binding row in the verified kit source. It is nil, and unbound
	// says why, when no verified kit source binds the compiled table.
	lines   map[string]int
	unbound string
	// effects classifies each Grammar binding by the committed change it covers.
	effects *effects
	// uses is every static use of each declared changed declaration.
	uses []use
}

// found is one path that one probe requires, with its reason.
type found struct {
	path   string
	kind   Kind
	reason Provenance
}

// closure walks every ticket entry to its fixed point. A ticket with no write reads no
// owner, so a repository whose tickets write nothing needs no readable inventory. changes
// are the inventory's declared changes, whose grammar declarations the walk resolves and
// whose callers and liveness censuses read. Each write entry of a ticket also reads the rule
// source of every verified reused owner in owners.
func closure(root string, source Source, parsed []tickets.Ticket, changes []Change, census *census, owners []Provenance) ([]Requirement, []Diagnostic, error) {
	if !slices.ContainsFunc(parsed, func(ticket tickets.Ticket) bool { return len(ticket.Writes) > 0 }) {
		return nil, nil, nil
	}
	for _, ticket := range parsed {
		for _, entry := range ticket.Writes {
			if path, _ := tickets.WritesPath(entry); !repositoryPath(path) {
				return nil, nil, &Error{Stage: StagePath, Subject: ticket.Name + ": " + entry, Cause: errors.New("the path escapes the repository or is not clean")}
			}
		}
	}
	inv, err := load(root, source)
	if err != nil {
		return nil, nil, err
	}
	inv.effects.declare(changes)
	var callerDiagnostics []Diagnostic
	inv.uses, callerDiagnostics = census.callers(changes)
	callerDiagnostics = append(callerDiagnostics, census.liveness(changes)...)
	var requirements []Requirement
	for _, ticket := range parsed {
		if len(ticket.Writes) == 0 {
			continue
		}
		owned := OwnedPaths(ticket)
		for _, entry := range ticket.Writes {
			requirements = append(requirements, inv.walk(ticket.Name, entry, owned)...)
			for _, owner := range owners {
				requirements = append(requirements, Requirement{Ticket: ticket.Name, Entry: entry, Path: owner.Owner, Kind: KindOwner, Provenance: []Provenance{owner}})
			}
		}
	}
	return requirements, append(append(inv.diagnostics(requirements), inv.effects.diagnostics...), callerDiagnostics...), nil
}

// repositoryPath reports whether one entry path names a place inside the repository by
// its clean spelling. A lexical fence prefix covers an entry by its spelling, so a ".."
// segment would let that prefix authorize a file outside the fence or the repository.
func repositoryPath(entry string) bool {
	return entry == "" || filepath.IsLocal(filepath.FromSlash(entry)) && path.Clean(entry) == entry
}

// load reads the fixture pins, the anchor literals, and the registry's verified source.
func load(root string, source Source) (inventory, error) {
	pins, err := canary.Pins(root)
	if err != nil {
		return inventory{}, &Error{Stage: StageFixtures, Subject: inventoryPath(root), Cause: err}
	}
	refs, err := anchors.References(root)
	if err != nil {
		return inventory{}, &Error{Stage: StageAnchors, Subject: anchors.RegistryDir, Cause: err}
	}
	inv := inventory{source: source, pins: map[string][]canary.Pin{}, refs: refs, effects: newEffects(root, source)}
	for _, pin := range pins {
		inv.pins[pin.Path] = append(inv.pins[pin.Path], pin)
	}
	inv.lines, inv.unbound = registryLines(root, source.Kit)
	return inv, nil
}

// inventoryPath is the fixture inventory's repository-relative directory.
func inventoryPath(root string) string {
	rel, err := filepath.Rel(root, canary.Dir(root))
	if err != nil {
		return canary.Dir(root)
	}
	return filepath.ToSlash(rel)
}

// registryLines locates each binding row in the kit source at root, which the caller's
// verified kit context authenticates. Without that context, a built-in path spelling is no
// source proof, so the location stays unbound.
func registryLines(root string, kit Kit) (map[string]int, string) {
	if kit.Identity == "" {
		return nil, "no verified kit source binds " + tickets.BindingRegistry
	}
	read := bounds.ClassifyNoFollow(filepath.Join(root, filepath.FromSlash(tickets.BindingRegistry)))
	if read.State != bounds.StateParsed {
		return nil, tickets.BindingRegistry + " is " + string(read.State) + " in the verified kit source"
	}
	lines, err := tickets.BindingLines(read.Data)
	if err != nil {
		return nil, err.Error()
	}
	return lines, ""
}

// walk follows one entry to its fixed point. Each newly required path that the ticket does
// not own is probed in turn, until no new path appears. A path the ticket owns is its own
// entry, so its closure stays under that entry. Visited paths end each cycle, and every
// reason found on the way stays with its requirement.
func (inv inventory) walk(ticket, entry string, owned []string) []Requirement {
	start, _ := tickets.WritesPath(entry)
	var group []Requirement
	at := map[string]int{}
	visited := map[string]bool{start: true}
	for queue := []string{start}; len(queue) > 0; queue = queue[1:] {
		for _, next := range inv.probe(queue[0]) {
			key := string(next.kind) + "\x00" + next.path
			i, seen := at[key]
			if !seen {
				i, at[key] = len(group), len(group)
				group = append(group, Requirement{Ticket: ticket, Entry: entry, Path: next.path, Kind: next.kind})
			}
			if !slices.Contains(group[i].Provenance, next.reason) {
				group[i].Provenance = append(group[i].Provenance, next.reason)
			}
			if !visited[next.path] && !Owned(owned, next.path) {
				visited[next.path] = true
				queue = append(queue, next.path)
			}
		}
	}
	sort.SliceStable(group, func(a, b int) bool { return group[a].Path < group[b].Path })
	return group
}

// probe is every path that subject requires under the canonical owners: each fixture that
// pins it, each file its binding rows require, each anchor registry file whose literal
// names it or a path under it, and each file that uses a changed declaration that subject
// declares before or after its change. A Core row always applies; a Grammar row applies
// only when the subject's grammar effect does, and a grammar change also requires its
// sources.
func (inv inventory) probe(subject string) []found {
	var next []found
	for _, pin := range inv.pins[subject] {
		next = append(next, found{pin.Fixture, KindFixture, inv.tree(pin.Owner, pin.Line, "fixture "+string(pin.Rule), pinReason(pin.Rule)+" "+subject)})
	}
	for _, use := range inv.uses {
		if use.file != subject && use.declares(subject) {
			next = append(next, found{use.file, KindCaller, use.reason(inv.source)})
		}
	}
	grammar := inv.effects.classify(subject)
	bound := false
	for _, rule := range tickets.BindingRules(subject) {
		if rule.Kind == tickets.BindingGrammar && !grammar.applies {
			continue
		}
		bound = bound || rule.Kind == tickets.BindingGrammar
		next = append(next, found{rule.File, KindRegistry, inv.registry(rule, subject)})
	}
	next = append(next, inv.grammarSources(subject, grammar, bound)...)
	if subject == "" {
		return next
	}
	for _, ref := range inv.refs {
		if tickets.Covers(subject, ref.Literal) {
			next = append(next, found{ref.Holder, KindAnchor, inv.tree(ref.Holder, ref.Line, "anchor literal", "literal "+ref.Literal+" anchors "+subject)})
		}
	}
	return next
}

// pinReason is the clause that names one pin rule's producer.
func pinReason(rule canary.PinRule) string {
	switch rule {
	case canary.PinInclude:
		return "BASE include names"
	case canary.PinOverlay:
		return "overlay carries"
	case canary.PinMutation:
		return "MUTATE.json anchors"
	}
	return "BASE names"
}

// tree is one reason whose owner is a committed file of the tree at the source tip.
func (inv inventory) tree(owner string, line int, rule, reason string) Provenance {
	return Provenance{Owner: owner, Line: line, Rule: rule, Reason: reason, Base: inv.source.Base, Tip: inv.source.Tip, Source: inv.source.Tip, Scope: ScopeTree}
}

// grammarSources is every path that subject's grammar changes and planned grammar forms
// require: the file that declares each grammar root a change reaches, and the public
// inventory files when no Grammar row already bound them. A planned form binds each write
// in its owner's package, before its implementation exists.
func (inv inventory) grammarSources(subject string, grammar effect, bound bool) []found {
	var next []found
	require := func(file string, reason Provenance) {
		if file != subject {
			next = append(next, found{file, KindRegistry, reason})
		}
	}
	for _, change := range grammar.changes {
		require(change.file, inv.tree(change.file, change.line, "grammar root "+change.root, change.unit+" changes the grammar that "+change.root+" produces"))
		if bound {
			continue
		}
		for _, file := range tickets.GrammarFiles() {
			require(file, inv.kit(change.root, "grammar root "+change.root, "grammar root "+change.root+" requires "+file+" beside "+subject))
		}
	}
	for _, change := range inv.effects.planned {
		root := declaredRoot(change)
		if !tickets.Covers(subject, change.Owner.Path) && path.Dir(subject) != root.Dir {
			continue
		}
		require(change.Owner.Path, inv.tree(inv.source.SpecPath, change.Location.Line, "planned grammar change", "planned grammar change "+change.ID+" registers through "+root.String()))
		for _, file := range tickets.GrammarFiles() {
			require(file, inv.kit(root.String(), "grammar root "+root.String(), "planned grammar change "+change.ID+" requires "+file+" beside "+subject))
		}
	}
	return next
}

// registry is one reason whose owner is a compiled binding row.
func (inv inventory) registry(rule tickets.BindingRule, subject string) Provenance {
	return inv.kit(rule.Prefix, "binding "+rule.Prefix, string(rule.Kind)+" binding "+rule.Prefix+" requires "+rule.File+" beside "+subject)
}

// kit is one reason whose owner is the compiled binding registry, at the row that key
// locates. Its location binds only through the verified kit source; otherwise it names the
// owner with no line.
func (inv inventory) kit(key, rule, text string) Provenance {
	reason := Provenance{
		Owner:  tickets.BindingRegistry,
		Rule:   rule,
		Reason: text,
		Base:   inv.source.Base,
		Tip:    inv.source.Tip,
		Scope:  ScopeKit,
	}
	if line, ok := inv.lines[key]; ok {
		reason.Line, reason.Source = line, inv.source.Kit.Identity
		reason.Reason += " at kit source " + inv.source.Kit.Identity
	}
	return reason
}

// diagnostics names each registry reason that no verified kit source located.
func (inv inventory) diagnostics(requirements []Requirement) []Diagnostic {
	for _, requirement := range requirements {
		for _, reason := range requirement.Provenance {
			if reason.Scope == ScopeKit && reason.Source == "" {
				unbound := inv.unbound
				if unbound == "" {
					unbound = "the verified kit source holds no row for " + reason.Rule
				}
				return []Diagnostic{{ID: kitRegistryID, Reason: unbound}}
			}
		}
	}
	return nil
}

// OwnedPaths is every tree path one ticket declares, with each marker stripped. Closure
// and its grading read this one set, so no two of them disagree about what a ticket owns.
func OwnedPaths(ticket tickets.Ticket) []string {
	paths := make([]string, 0, len(ticket.Writes))
	for _, entry := range ticket.Writes {
		path, _ := tickets.WritesPath(entry)
		paths = append(paths, path)
	}
	return paths
}

// Owned reports whether one owned path covers path at a segment boundary.
func Owned(owned []string, path string) bool {
	return slices.ContainsFunc(owned, func(entry string) bool { return tickets.Covers(entry, path) })
}
