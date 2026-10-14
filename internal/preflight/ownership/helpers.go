package ownership

import (
	"encoding/json"
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
		if match.PkgPath == owner.Package && len(pkgs) > 0 {
			position := pkgs[0].Fset.Position(match.Obj.Pos())
			rel, err := filepath.Rel(c.root, position.Filename)
			if err != nil {
				return nil, "", 0, err
			}
			return match.Obj, filepath.ToSlash(rel), position.Line, nil
		}
	}
	return nil, "", 0, nil
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
		case obj == nil && strings.Contains(localName(change.Owner.Symbol), "."):
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

// localName is a symbol's name inside its package, after the package qualifier.
func localName(symbol string) string {
	_, name, _ := strings.Cut(symbol, ".")
	return name
}

// ownerSource is one verified reused owner: the premise that declares it and the reason
// that locates its rule source.
type ownerSource struct {
	premise string
	reason  Provenance
}

// premises grades each current premise of a kind this census owns: a reused owner, and the
// retirement, posture, and printed-field member censuses. It answers the graded row of each
// such premise by ID, and each verified reused owner that the closure reads. A fault
// leaves its premise unknown with that fault as the reason.
func (c *census) premises(items []PremiseItem) (map[string]Premise, []ownerSource) {
	graded := map[string]Premise{}
	var owners []ownerSource
	for _, item := range items {
		var reason string
		var err error
		switch {
		case item.State != "current":
			continue
		case item.Kind == "reused-owner":
			var owner ownerSource
			if owner, reason, err = c.reusedOwner(item); err == nil {
				owners = append(owners, owner)
			}
		case item.Kind == "retirement" || item.Kind == "posture" || item.Kind == "printed-field":
			reason, err = c.memberCensus(item)
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
// rule source that the closure reads.
func (c *census) reusedOwner(item PremiseItem) (ownerSource, string, error) {
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
		return ownerSource{}, "", err
	}
	loc := subject.Owner
	at := fmt.Sprintf("%s line %d at the source %s", loc.Path, loc.Line, loc.Source)
	reason := Provenance{Owner: loc.Path, Line: loc.Line, Rule: "reused owner", Reason: "premise " + item.ID + " reuses the rule that " + at + " holds",
		Base: c.source.Base, Tip: c.source.Tip, Source: rev, Scope: ScopeTree}
	return ownerSource{premise: item.ID, reason: reason}, "reused owner " + at + " holds its needle", nil
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

// memberCensus grades one member census against the static references at the source tip.
// A retirement premise names every consumer of its owner, each with one disposition; a
// posture premise names every static caller of its writer; a printed-field premise names
// every declaration that references its field. A missing or extra member is no complete
// census. A complete posture census still waits for its writer probe, which only a
// reused-test collector certifies.
func (c *census) memberCensus(item PremiseItem) (string, error) {
	// Each kind's subject keys and expected keys, with the subject or expected key that
	// lists the members.
	shape := map[string]struct {
		subject, expected []string
		owner, members    string
	}{
		"retirement":    {[]string{"owner", "consumers"}, []string{"dispositions"}, "owner", "consumers"},
		"posture":       {[]string{"owner", "writer_probe"}, []string{"callers"}, "owner", "callers"},
		"printed-field": {[]string{"field"}, []string{"producers"}, "field", "producers"},
	}[item.Kind]
	subject, err := fields(item.Subject, "subject: ", shape.subject...)
	if err != nil {
		return "", err
	}
	expected, err := fields(item.Expected, "expected: ", shape.expected...)
	if err != nil {
		return "", err
	}
	owner, err := identity(subject[shape.owner], "subject: "+shape.owner+": ")
	if err != nil {
		return "", err
	}
	members := subject[shape.members]
	if members == nil {
		members = expected[shape.members]
	}
	listed, err := identities(members, shape.members+": ")
	if err != nil {
		return "", err
	}
	switch item.Kind {
	case "retirement":
		err = dispositions(listed, expected["dispositions"])
	case "posture":
		var probe map[string]json.RawMessage
		if json.Unmarshal(subject["writer_probe"], &probe) != nil || probe == nil {
			err = fmt.Errorf("subject: writer_probe is not a tagged object")
		}
	}
	if err != nil {
		return "", err
	}
	actual, err := c.members(owner, item.Kind == "posture")
	if err != nil {
		return "", err
	}
	if err := censusGaps(listed, actual); err != nil {
		return "", err
	}
	if item.Kind == "posture" {
		return "", fmt.Errorf("the caller census of %s is complete, but no collector certifies its writer_probe evidence", owner.Symbol)
	}
	return fmt.Sprintf("the %s census of %s is complete at the source tip: %d member(s)", item.Kind, owner.Symbol, len(actual)), nil
}

// members is every declaration that statically references owner at the source tip, or
// only every static caller. Each member is a file and its enclosing declaration name.
func (c *census) members(owner Owner, calls bool) ([]Owner, error) {
	obj, file, _, err := c.declaration(owner)
	if err != nil {
		return nil, err
	}
	if obj == nil || file != owner.Path {
		return nil, fmt.Errorf("%s does not resolve in %s at the source tip", owner.Symbol, owner.Path)
	}
	var found []Owner
	for _, row := range consumers.Rows(c.pkgs, obj, c.root) {
		member := Owner{Path: row.File, Symbol: row.Enclosing}
		if (!calls || row.Via == consumers.ViaCall) && !slices.Contains(found, member) {
			found = append(found, member)
		}
	}
	return found, nil
}

// censusGaps names each actual member that listed omits and each listed identity that is
// no member. A listed identity matches a member by path and by its name inside its package.
func censusGaps(listed, actual []Owner) error {
	matches := func(owner, member Owner) bool {
		return owner.Path == member.Path && localName(owner.Symbol) == member.Symbol
	}
	var gaps []string
	for _, member := range actual {
		if !slices.ContainsFunc(listed, func(owner Owner) bool { return matches(owner, member) }) {
			gaps = append(gaps, "omits "+member.Path+" "+member.Symbol)
		}
	}
	for _, owner := range listed {
		if !slices.ContainsFunc(actual, func(member Owner) bool { return matches(owner, member) }) {
			gaps = append(gaps, "names "+owner.Symbol+", which is no member")
		}
	}
	if len(gaps) > 0 {
		return fmt.Errorf("the census %s", strings.Join(gaps, "; "))
	}
	return nil
}

// dispositions grades a retirement's dispositions: one valid action for each consumer.
func dispositions(consumers []Owner, raw json.RawMessage) error {
	var items []json.RawMessage
	if err := json.Unmarshal(raw, &items); err != nil || items == nil {
		return fmt.Errorf("expected: dispositions is not a JSON array")
	}
	actions := map[string][]string{}
	for i, item := range items {
		decoded, err := disposition(item, fmt.Sprintf("expected: dispositions[%d]: ", i))
		if err != nil {
			return err
		}
		actions[decoded.Symbol] = append(actions[decoded.Symbol], decoded.Action)
	}
	for _, consumer := range consumers {
		if len(actions[consumer.Symbol]) != 1 {
			return fmt.Errorf("retirement consumer %s has %d dispositions, want 1", consumer.Symbol, len(actions[consumer.Symbol]))
		}
	}
	return nil
}
