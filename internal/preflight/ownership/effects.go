package ownership

import (
	"cmp"
	"maps"
	"path"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/tickets"
)

// grammarEffectsID names each diagnostic of grammar effect classification.
const grammarEffectsID = "grammar-effects"

// effects classifies the grammar effect of the committed source range. It reads each
// classified package, and each grammar package that can reach it, at the source base and at
// the source tip, so a classification binds exactly those two commits.
type effects struct {
	root        string
	source      Source
	reads       map[string]producerRead
	modules     map[string]revRead[string]
	grammar     map[string]revRead[map[string]bool]
	packages    map[string]packageEffect
	planned     []Change
	diagnostics []Diagnostic
}

// producerRead is one package read at one revision, or the fault that stopped it.
type producerRead struct {
	producers producers
	err       error
}

// grammarChange is one changed producer declaration and one grammar root that reaches it:
// the root's selector and the file and line that declare that root.
type grammarChange struct {
	unit, root, file string
	line             int
}

// packageEffect is one package's classification: the grammar changes in each changed file,
// or the reason that the package cannot be classified.
type packageEffect struct {
	changes map[string][]grammarChange
	failure string
}

// effect is the grammar effect of one subject. applies is false only when the source range
// holds a committed change that the subject covers, and every such change is classified
// and changes no grammar producer. A write that no committed change shows keeps the
// obligation, because nothing proves it body-only.
type effect struct {
	applies bool
	changes []grammarChange
}

func newEffects(root string, source Source) *effects {
	return &effects{root: root, source: source, reads: map[string]producerRead{}, modules: map[string]revRead[string]{},
		grammar: map[string]revRead[map[string]bool]{}, packages: map[string]packageEffect{}}
}

func (e *effects) diagnose(id, reason string) {
	e.diagnostics = append(e.diagnostics, Diagnostic{ID: id, Reason: reason})
}

// classify is the grammar effect of one subject over the source range. Every changed
// production Go file is classified, so a package with no Grammar row still changes grammar
// when one of its own parser calls, or a grammar root of another package, reaches it.
func (e *effects) classify(subject string) effect {
	if e.source.Base == "" || e.source.Tip == "" {
		return effect{applies: true}
	}
	var result effect
	observed := false
	for _, changed := range e.source.Paths {
		if !tickets.Covers(subject, changed.Path) {
			continue
		}
		observed = true
		if !producerFile(changed.Path) {
			continue
		}
		pkg := e.pkg(path.Dir(changed.Path))
		if pkg.failure != "" || len(pkg.changes[changed.Path]) > 0 {
			result.applies = true
			result.changes = append(result.changes, pkg.changes[changed.Path]...)
		}
	}
	result.applies = result.applies || !observed
	return result
}

// producerFile reports whether one path is production Go source.
func producerFile(name string) bool {
	return strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")
}

// pkg classifies one package once. A producer declaration changes when its printed syntax
// differs between the base and the tip, or when only one side declares it. It changes the
// grammar of each root of its own package that reaches it, and of each root of another
// package that reaches it through in-module imports. An unchanged dependency stays
// unchanged even when its caller changes.
func (e *effects) pkg(dir string) packageEffect {
	if cached, ok := e.packages[dir]; ok {
		return cached
	}
	result := packageEffect{changes: map[string][]grammarChange{}}
	base, tip := e.read(e.source.Base, dir), e.read(e.source.Tip, dir)
	var baseRefs, tipRefs map[string][]grammarChange
	var baseFailure, tipFailure string
	switch {
	case base.err != nil:
		result.failure = "base " + base.err.Error()
	case tip.err != nil:
		result.failure = "tip " + tip.err.Error()
	case len(base.producers.unresolved) > 0:
		result.failure = "base " + base.producers.unresolved[0]
	case len(tip.producers.unresolved) > 0:
		result.failure = "tip " + tip.producers.unresolved[0]
	}
	if result.failure == "" {
		baseRefs, baseFailure = e.referred(e.source.Base, dir)
		tipRefs, tipFailure = e.referred(e.source.Tip, dir)
	}
	switch {
	case baseFailure != "":
		result.failure = "base " + baseFailure
	case tipFailure != "":
		result.failure = "tip " + tipFailure
	}
	if result.failure != "" {
		e.diagnose(grammarEffectsID, dir+": "+result.failure+", so its grammar effect is unclassified")
		e.packages[dir] = result
		return result
	}
	before, after := base.producers, tip.producers
	for _, key := range unionKeys(before.units, after.units) {
		old, inBase := before.units[key]
		current, inTip := after.units[key]
		if inBase && inTip && old.text == current.text {
			continue
		}
		var changes []grammarChange
		for _, root := range append(append([]string(nil), before.reached[key]...), after.reached[key]...) {
			at, ok := after.units[root]
			if !ok {
				at = before.units[root]
			}
			changes = append(changes, grammarChange{root: tickets.GrammarRoot{Dir: dir, Name: root}.String(), file: at.file, line: at.line})
		}
		changes = append(append(changes, tipRefs[key]...), baseRefs[key]...)
		slices.SortStableFunc(changes, func(a, b grammarChange) int { return cmp.Compare(a.root, b.root) })
		for _, change := range slices.CompactFunc(changes, func(a, b grammarChange) bool { return a.root == b.root }) {
			change.unit = tickets.GrammarRoot{Dir: dir, Name: key}.String()
			for _, file := range []string{old.file, current.file} {
				if file != "" && !slices.Contains(result.changes[file], change) {
					result.changes[file] = append(result.changes[file], change)
				}
			}
		}
	}
	e.packages[dir] = result
	return result
}

// read loads one package's producers at one revision once.
func (e *effects) read(rev, dir string) producerRead {
	key := rev + "\x00" + dir
	if cached, ok := e.reads[key]; ok {
		return cached
	}
	module, err := e.module(rev)
	var producers producers
	if err == nil {
		producers, err = readProducers(e.root, rev, dir, module)
	}
	e.reads[key] = producerRead{producers, err}
	return e.reads[key]
}

// declare resolves each declared grammar change to a registered producer at the source tip.
// A change that names no producer stays diagnostic, so a declaration whose name only looks
// like a command entry, such as a function named Command, never completes. A planned change
// that resolves joins the closure of each write beside its owner.
func (e *effects) declare(changes []Change) {
	for _, change := range changes {
		if change.Kind != "grammar" {
			continue
		}
		root := declaredRoot(change)
		read := e.read(e.source.Tip, root.Dir)
		unit, declared := read.producers.units[root.Name]
		if read.err != nil || !declared || unit.file != change.Owner.Path || !slices.Contains(read.producers.roots, root.Name) {
			e.diagnose(change.ID, "grammar change names "+change.Owner.Symbol+" in "+change.Owner.Path+", which no registered grammar producer declares at the source tip")
			continue
		}
		if change.State == "planned" {
			e.planned = append(e.planned, change)
		}
	}
}

// declaredRoot is the selector a declared change names: its owner's package directory and
// the exact declaration name after the symbol's package qualifier.
func declaredRoot(change Change) tickets.GrammarRoot {
	return tickets.GrammarRoot{Dir: path.Dir(change.Owner.Path), Name: localName(change.Owner.Symbol)}
}

// unionKeys is the sorted key union of two unit maps.
func unionKeys(a, b map[string]unit) []string {
	keys := slices.AppendSeq(slices.Collect(maps.Keys(a)), maps.Keys(b))
	slices.Sort(keys)
	return slices.Compact(keys)
}
