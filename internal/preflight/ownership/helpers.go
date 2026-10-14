package ownership

import (
	"fmt"
	"go/types"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/consumers"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/tickets"
)

// callersID names each diagnostic of a changed declaration's caller census.
const callersID = "caller-census"

// census reads the static references of declared owners. Its tip side is the typed
// packages of the checkout, which must be clean and at the source tip, loaded once on first
// use; its base side reads the source base without types. A snapshot keeps only the
// resolved identities that the census reads, never the packages themselves.
type census struct {
	root   string
	source Source
	loaded bool
	pkgs   []*consumers.Package
	err    error
}

// use is one file that statically uses one changed declaration: the file and the line of
// its first use, the revision that holds that use, and the change with the file and line
// that declare it there.
type use struct {
	change         Change
	file           string
	line           int
	rev, side      string
	declared       string
	declaredAtLine int
}

// declares reports whether subject covers the declaration that the use names: the owner
// path before the change, or its declared destination.
func (u use) declares(subject string) bool {
	return tickets.Covers(subject, u.change.Owner.Path) || u.change.Destination != nil && tickets.Covers(subject, *u.change.Destination)
}

// reason is the source reason of one use. The use file and line are the evidence, and the
// reason names the declaration that the change changes.
func (u use) reason(source Source) Provenance {
	return Provenance{Owner: u.file, Line: u.line, Rule: "caller " + u.change.Owner.Symbol,
		Reason: fmt.Sprintf("%s uses %s, which change %s changes, declared in %s line %d at the source %s", u.file, u.change.Owner.Symbol, u.change.ID, u.declared, u.declaredAtLine, u.side),
		Base:   source.Base, Tip: source.Tip, Source: u.rev, Scope: ScopeTree}
}

// packages is the typed tip side. A checkout that is not the source tip would enumerate
// bytes that the source pair does not hold, so it answers an error.
func (c *census) packages() ([]*consumers.Package, error) {
	if !c.loaded {
		c.loaded = true
		head, err := git.Output("-C", c.root, "rev-parse", "HEAD")
		switch {
		case err != nil:
			c.err = fmt.Errorf("the checkout HEAD does not resolve: %w", err)
		case head != c.source.Tip:
			c.err = fmt.Errorf("the checkout HEAD %s is not the source tip %s", head, c.source.Tip)
		default:
			c.pkgs, c.err = consumers.Load(c.root)
		}
	}
	return c.pkgs, c.err
}

// declaration resolves one owner identity at the source tip: the declaration that its
// symbol names in its exact package, and the repository path and line that declare it. An
// absent declaration answers a nil object and no error.
func (c *census) declaration(owner Owner) (types.Object, string, int, error) {
	pkgs, err := c.packages()
	if err != nil {
		return nil, "", 0, err
	}
	matches, err := consumers.Resolve(pkgs, owner.Symbol)
	if err != nil {
		return nil, "", 0, nil
	}
	for _, match := range matches {
		if match.PkgPath == owner.Package {
			file, line, err := c.position(pkgs, match.Obj)
			return match.Obj, file, line, err
		}
	}
	return nil, "", 0, nil
}

// position is the repository path and line that declare one loaded object.
func (c *census) position(pkgs []*consumers.Package, obj types.Object) (string, int, error) {
	position := pkgs[0].Fset.Position(obj.Pos())
	rel, err := filepath.Rel(c.root, position.Filename)
	if err != nil {
		return "", 0, err
	}
	return filepath.ToSlash(rel), position.Line, nil
}

// callers is every static use of each existing signature or fixture-helper change, in
// production and test code: the typed references at the source tip and the lexical uses at
// the source base. The base side keeps the former users of a relocated or deleted
// declaration. A planned change declares a helper that has no callers yet, so it reads
// none. Each census that cannot read stays diagnostic.
func (c *census) callers(changes []Change) ([]use, []Diagnostic) {
	var uses []use
	var diagnostics []Diagnostic
	diagnose := func(change Change, format string, args ...any) {
		diagnostics = append(diagnostics, Diagnostic{ID: callersID, Reason: "change " + change.ID + ": " + fmt.Sprintf(format, args...)})
	}
	for _, change := range changes {
		if change.State != "existing" || change.Kind != "signature" && change.Kind != "fixture-helper" {
			continue
		}
		obj, file, line, err := c.declaration(change.Owner)
		if err != nil {
			diagnose(change, "the tip census does not read: %v", err)
		}
		if obj != nil {
			if file != change.Owner.Path && (change.Destination == nil || file != *change.Destination) {
				diagnose(change, "%s is declared in %s, which is neither its owner path nor its declared destination", change.Owner.Symbol, file)
			}
			for _, row := range consumers.Rows(c.pkgs, obj, c.root) {
				uses = appendUse(uses, use{change: change, file: row.File, line: row.Line, rev: c.source.Tip, side: "tip", declared: file, declaredAtLine: line})
			}
		}
		before, err := c.before(change)
		switch {
		case c.source.Base == "":
			diagnose(change, "no explicit source base, so the former callers of %s are not read", change.Owner.Symbol)
		case err != nil:
			diagnose(change, "the base census does not read: %v", err)
		case before.line == 0:
			diagnose(change, "%s declares no %s at the source base", change.Owner.Path, change.Owner.Symbol)
		case strings.Contains(localName(change.Owner), "."):
			diagnose(change, "%s is a method, whose former callers need the types that the source base does not load", change.Owner.Symbol)
		}
		for _, found := range before.uses {
			uses = appendUse(uses, use{change: change, file: found.File, line: found.Line, rev: c.source.Base, side: "base", declared: change.Owner.Path, declaredAtLine: before.line})
		}
	}
	return uses, diagnostics
}

// appendUse keeps the first use of each file for each change and revision.
func appendUse(uses []use, next use) []use {
	if slices.ContainsFunc(uses, func(u use) bool { return u.change.ID == next.change.ID && u.file == next.file && u.rev == next.rev }) {
		return uses
	}
	return append(uses, next)
}

// localName is an owner's declaration spelling inside its package, by the consumers query
// grammar. Owner decoding refuses a symbol with no such spelling, so that read cannot fail
// here.
func localName(owner Owner) string {
	name, _ := consumers.LocalName(owner.Package, owner.Symbol)
	return name
}

// premises grades each current premise of a kind this census owns: a reused owner, and the
// retirement, posture, and printed-field member censuses. It answers the graded row of each
// such premise by ID, and the reason that locates each verified reused owner's rule source,
// which the closure reads. A fault
// leaves its premise unknown with that fault as the reason.
func (c *census) premises(items []PremiseItem) (map[string]Premise, []Provenance) {
	graded := map[string]Premise{}
	var owners []Provenance
	for _, item := range items {
		var reason string
		var err error
		shape, member := memberCensuses[item.Kind]
		switch {
		case item.State != "current":
			continue
		case item.Kind == "reused-owner":
			var owner Provenance
			if owner, reason, err = c.reusedOwner(item); err == nil {
				owners = append(owners, owner)
			}
		case member:
			reason, err = c.memberCensus(item, shape)
		default:
			continue
		}
		row := Premise{ID: item.ID, State: StateGreen, Source: c.source.Tip, Reason: reason}
		if err != nil {
			row.State, row.Reason = StateUnknown, item.Kind+" premise: "+err.Error()
		}
		graded[item.ID] = row
	}
	return graded, owners
}

// reusedOwner verifies one reused owner location at its pinned source and answers the
// reason that locates the rule source the closure reads.
func (c *census) reusedOwner(item PremiseItem) (Provenance, string, error) {
	var subject struct {
		Owner Location `json:"owner"`
	}
	var expected struct {
		Required bool `json:"required"`
	}
	object, err := typed(item.Subject, "subject: ", &subject, "owner")
	if err == nil {
		err = location(object, "owner", "subject: ", true)
	}
	if err == nil {
		_, err = typed(item.Expected, "expected: ", &expected, "required")
	}
	if err == nil && !expected.Required {
		err = fmt.Errorf("expected: required is not true")
	}
	var rev string
	if err == nil {
		rev, err = c.locate(subject.Owner)
	}
	if err != nil {
		return Provenance{}, "", err
	}
	loc := subject.Owner
	at := fmt.Sprintf("%s line %d at the source %s", loc.Path, loc.Line, loc.Source)
	reason := Provenance{Owner: loc.Path, Line: loc.Line, Rule: "reused owner", Reason: "premise " + item.ID + " reuses the rule that " + at + " holds",
		Base: c.source.Base, Tip: c.source.Tip, Source: rev, Scope: ScopeTree}
	return reason, "reused owner " + at + " holds its needle", nil
}

// locate verifies that one owner location's line holds its needle exactly once at its
// pinned source, and answers that source's commit. A moved or repeated needle is a fault,
// never a nearby match.
func (c *census) locate(loc Location) (string, error) {
	rev := c.source.Tip
	if loc.Source == "base" {
		rev = c.source.Base
	}
	if rev == "" {
		return "", fmt.Errorf("no source %s pins %s", loc.Source, loc.Path)
	}
	data, err := git.ReadTreeFile(c.root, rev, loc.Path)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(data), "\n")
	if loc.Line > len(lines) || strings.Count(strings.TrimSuffix(lines[loc.Line-1], "\r"), loc.Needle) != 1 {
		return "", fmt.Errorf("%s line %d does not hold its needle exactly once at the source %s", loc.Path, loc.Line, loc.Source)
	}
	return rev, nil
}
