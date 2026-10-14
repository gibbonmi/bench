package ownership

import (
	"encoding/json"
	"fmt"
	"go/types"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/consumers"
)

// livenessID names each diagnostic of an existing liveness change's retirement census.
const livenessID = "liveness-census"

// liveness grades each existing liveness change's dispositions against the retirement
// census of its owner. A planned change has no consumers yet, so it reads none.
func (c *census) liveness(changes []Change) []Diagnostic {
	var diagnostics []Diagnostic
	for _, change := range changes {
		if change.State != "existing" || change.Kind != "liveness" {
			continue
		}
		members, err := c.members(change.Owner, memberCensuses["retirement"].rows)
		if err == nil {
			err = c.dispositions(members, change.Dispositions)
		}
		if err != nil {
			diagnostics = append(diagnostics, Diagnostic{ID: livenessID, Reason: "change " + change.ID + ": " + err.Error()})
		}
	}
	return diagnostics
}

// memberShape is one member census kind: its exact subject and expected keys, the subject
// key that names its owner, the key that lists its members, and the static references that
// make a member. A kind whose expected keys hold the member key lists its members on both
// sides, and the two lists must be one set.
type memberShape struct {
	subject, expected []string
	owner, members    string
	rows              func(pkgs []*consumers.Package, target types.Object, root string) []consumers.Row
}

// memberCensuses is every member census kind. A retirement's members are every reference
// to its owner, a posture's members are the static callers of its writer, and a
// printed-field's members are the static writes of its field.
var memberCensuses = map[string]memberShape{
	"retirement":    {[]string{"owner", "consumers"}, []string{"dispositions"}, "owner", "consumers", consumers.Rows},
	"posture":       {[]string{"owner", "callers", "writer_probe"}, []string{"callers"}, "owner", "callers", callRows},
	"printed-field": {[]string{"field", "producers"}, []string{"producers"}, "field", "producers", consumers.WriteRows},
}

// callRows is every static call of target.
func callRows(pkgs []*consumers.Package, target types.Object, root string) []consumers.Row {
	var calls []consumers.Row
	for _, row := range consumers.Rows(pkgs, target, root) {
		if row.Via == consumers.ViaCall {
			calls = append(calls, row)
		}
	}
	return calls
}

// memberCensus grades one member census against the static references at the source tip.
// A retirement premise names every consumer of its owner, each with one disposition; a
// posture premise names every static caller of its writer, with a reused-test writer probe;
// a printed-field premise names every declaration that writes its field. A missing or extra
// member is no complete census. A complete posture census still waits for its writer probe,
// which only a reused-test collector certifies.
func (c *census) memberCensus(item PremiseItem, shape memberShape) (string, error) {
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
	listed, err := identities(subject[shape.members], "subject: "+shape.members+": ")
	if err != nil {
		return "", err
	}
	var decided []Disposition
	switch item.Kind {
	case "retirement":
		decided, err = decodeDispositions(expected["dispositions"])
	case "posture":
		_, err = fields(subject["writer_probe"], "subject: writer_probe: ", reusedTestKeys...)
	}
	if err != nil {
		return "", err
	}
	if slices.Contains(shape.expected, shape.members) {
		again, err := identities(expected[shape.members], "expected: "+shape.members+": ")
		if err != nil {
			return "", err
		}
		if err := censusGaps(again, asMembers(listed)); err != nil {
			return "", fmt.Errorf("the expected %s differ from the subject %s: %w", shape.members, shape.members, err)
		}
	}
	actual, err := c.members(owner, shape.rows)
	if err != nil {
		return "", err
	}
	if err := censusGaps(listed, actual); err != nil {
		return "", err
	}
	switch item.Kind {
	case "retirement":
		if err := c.dispositions(actual, decided); err != nil {
			return "", err
		}
	case "posture":
		return "", fmt.Errorf("the caller census of %s is complete, but no collector certifies its writer_probe evidence", owner.Symbol)
	}
	return fmt.Sprintf("the %s census of %s is complete at the source tip: %d member(s)", item.Kind, owner.Symbol, len(actual)), nil
}

// members is every declaration that one row reader finds for owner at the source tip. Each
// member is a file and its enclosing declaration name.
func (c *census) members(owner Owner, rows func([]*consumers.Package, types.Object, string) []consumers.Row) ([]Owner, error) {
	obj, file, _, err := c.declaration(owner)
	if err != nil {
		return nil, err
	}
	if obj == nil || file != owner.Path {
		return nil, fmt.Errorf("%s does not resolve in %s at the source tip", owner.Symbol, owner.Path)
	}
	var found []Owner
	for _, row := range rows(c.pkgs, obj, c.root) {
		member := Owner{Path: row.File, Symbol: row.Enclosing}
		if !slices.Contains(found, member) {
			found = append(found, member)
		}
	}
	return found, nil
}

// asMembers is each listed identity in member form: its path and its name inside its
// package.
func asMembers(listed []Owner) []Owner {
	members := make([]Owner, len(listed))
	for i, owner := range listed {
		members[i] = Owner{Path: owner.Path, Symbol: localName(owner)}
	}
	return members
}

// censusGaps names each actual member that listed omits and each listed identity that is
// no member. A listed identity matches a member by path and by its name inside its package.
func censusGaps(listed, actual []Owner) error {
	matches := func(owner, member Owner) bool {
		return owner.Path == member.Path && localName(owner) == member.Symbol
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

// decodeDispositions decodes a retirement's dispositions array, each through the one
// disposition grammar.
func decodeDispositions(raw json.RawMessage) ([]Disposition, error) {
	var items []json.RawMessage
	if err := json.Unmarshal(raw, &items); err != nil || items == nil {
		return nil, fmt.Errorf("expected: dispositions is not a JSON array")
	}
	decided := make([]Disposition, 0, len(items))
	for i, item := range items {
		decoded, err := disposition(item, fmt.Sprintf("expected: dispositions[%d]: ", i))
		if err != nil {
			return nil, err
		}
		decided = append(decided, decoded)
	}
	return decided, nil
}

// dispositions grades dispositions against the retirement census members: each member takes
// exactly one disposition, and each disposition names a member. A disposition's symbol names
// each declaration that it resolves to at the source tip.
func (c *census) dispositions(members []Owner, decided []Disposition) error {
	pkgs, err := c.packages()
	if err != nil {
		return err
	}
	counts := map[Owner]int{}
	for _, item := range decided {
		matches, _ := consumers.Resolve(pkgs, item.Symbol)
		named := false
		for _, match := range matches {
			file, _, err := c.position(pkgs, match.Obj)
			if err != nil {
				return err
			}
			if member := (Owner{Path: file, Symbol: match.Name}); slices.Contains(members, member) {
				counts[member]++
				named = true
			}
		}
		if !named {
			return fmt.Errorf("disposition %s names no retirement consumer", item.Symbol)
		}
	}
	for _, member := range members {
		if counts[member] != 1 {
			return fmt.Errorf("%s %s has %d dispositions, want 1", member.Path, member.Symbol, counts[member])
		}
	}
	return nil
}
