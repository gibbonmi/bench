// The record of one cleanup baseline case, and the identity mapping it renders through. Only
// fixture-dependent identities are renamed, each through a named mapping: the repository root
// and home, each member path, branch, identity, and owner, each ref tip, and each digest whose
// role the case knows. Actions, row order, details, errors, recovery refs, receipts, and
// deletion outcomes stay literal, and equal digests keep equal names.
package worktree

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

// firstRecordDifference names the first line where two records part.
func firstRecordDifference(got, frozen string) string {
	a, b := strings.Split(got, "\n"), strings.Split(frozen, "\n")
	for i := 0; i < len(a) || i < len(b); i++ {
		var x, y string
		if i < len(a) {
			x = a[i]
		}
		if i < len(b) {
			y = b[i]
		}
		if x != y {
			return fmt.Sprintf("line %d: got %q, frozen %q", i+1, x, y)
		}
	}
	return "no line"
}

// namePair is one fixture identity and the name a record gives it.
type namePair struct{ raw, name string }

// baselineNames is the bijection of one case: each identity to its name, applied longest
// first, so a path is renamed before the identity it contains.
type baselineNames struct {
	pairs   []namePair
	members []baselineMember
}

// baselineMember is one recorded assignment and its name.
type baselineMember struct {
	name       string
	assignment intent.Assignment
}

func (n *baselineNames) add(t *testing.T, raw, name string) {
	t.Helper()
	if raw == "" {
		return
	}
	for _, pair := range n.pairs {
		if pair.raw == raw && pair.name != name {
			t.Fatalf("identity %q has two names, %q and %q", raw, pair.name, name)
		}
		if pair.raw == raw {
			return
		}
		if pair.name == name {
			t.Fatalf("name %q names two identities, %q and %q", name, pair.raw, raw)
		}
	}
	n.pairs = append(n.pairs, namePair{raw, name})
}

func (n *baselineNames) replace(text string) string {
	pairs := append([]namePair(nil), n.pairs...)
	sort.SliceStable(pairs, func(a, b int) bool { return len(pairs[a].raw) > len(pairs[b].raw) })
	args := make([]string, 0, 2*len(pairs))
	for _, pair := range pairs {
		args = append(args, pair.raw, pair.name)
	}
	return strings.NewReplacer(args...).Replace(text)
}

// baselineIdentities names the case's fixture identities from the repository as the apply
// finds it. Members take their identity rank, or their label under byRole. Each ref tip takes
// the name of the first ref that points at it, and each known digest takes its role.
func baselineIdentities(t *testing.T, run baselineRun) *baselineNames {
	t.Helper()
	names := &baselineNames{}
	names.add(t, run.root, "<root>")
	names.add(t, run.home, "<home>")
	if real, err := filepath.EvalSymlinks(run.root); err == nil && real != run.root {
		names.add(t, real, "<root:real>")
	}
	assignments, err := intent.Assignments(run.root)
	mustNoError(t, err)
	sort.Slice(assignments, func(a, b int) bool { return assignments[a].ID < assignments[b].ID })
	for i, assignment := range assignments {
		name := fmt.Sprintf("m%d", i+1)
		if run.byRole {
			name = strings.ReplaceAll(assignment.Label, " ", "-")
		}
		names.members = append(names.members, baselineMember{name, assignment})
	}
	sort.Slice(names.members, func(a, b int) bool { return names.members[a].name < names.members[b].name })
	// An owner takes the name of the first member that carries it, so members that share an
	// owner keep that relation.
	owners := map[string]bool{}
	for _, member := range names.members {
		assignment := member.assignment
		names.add(t, assignment.Worktree, "<"+member.name+":path>")
		names.add(t, assignment.Branch, "<"+member.name+":branch>")
		names.add(t, strings.TrimPrefix(assignment.Branch, "refs/heads/"), "<"+member.name+":branch-name>")
		names.add(t, assignment.ID, "<"+member.name+":id>")
		if !owners[assignment.OwnerID] && assignment.OwnerID != "" {
			owners[assignment.OwnerID] = true
			names.add(t, assignment.OwnerID, "<"+member.name+":owner>")
		}
	}
	tips := map[string]string{}
	for _, line := range baselineRefs(t, run.root) {
		oid, ref, _ := strings.Cut(line, " ")
		ref = names.replace(ref)
		if current, known := tips[oid]; !known || ref < current {
			tips[oid] = ref
		}
	}
	for oid, ref := range tips {
		names.add(t, oid, "<tip:"+ref+">")
	}
	names.add(t, run.planDigest, "<digest:plan>")
	for digest, id := range run.targetDigests {
		names.add(t, digest, "<digest:target:"+names.replace(id)+">")
	}
	return names
}

// baselineRefs is every ref in the repository as "<oid> <ref>".
func baselineRefs(t *testing.T, root string) []string {
	t.Helper()
	listing := gitOutput(t, root, "for-each-ref", "--format=%(objectname) %(refname)")
	if listing == "" {
		return nil
	}
	return strings.Split(listing, "\n")
}

// rankedTip is the name of a member branch tip under identity rank.
var rankedTip = regexp.MustCompile(`<tip:<m[0-9]+:branch>>`)

// commitRelations records each named tip's tree and parents, so a renamed commit still has
// to keep its content and its ancestry. A parent no ref names reads as its tree. A fixture
// creates its members in an order that their random identities do not follow, so under
// identity rank the relations read as one set of member tips, without the rank.
func commitRelations(t *testing.T, root string, names *baselineNames, byRole bool) []string {
	t.Helper()
	lines := commitRelationLines(t, root, names)
	if !byRole {
		for i, line := range lines {
			lines[i] = rankedTip.ReplaceAllString(names.replace(line), "<member-tip>")
		}
	}
	return lines
}

func commitRelationLines(t *testing.T, root string, names *baselineNames) []string {
	t.Helper()
	tree := func(oid string) string { return "tree:" + gitOutput(t, root, "rev-parse", oid+"^{tree}")[:12] }
	var lines []string
	for _, pair := range names.pairs {
		if !strings.HasPrefix(pair.name, "<tip:") {
			continue
		}
		fields := strings.Fields(gitOutput(t, root, "rev-list", "--parents", "-n", "1", pair.raw))
		parents := make([]string, 0, len(fields)-1)
		for _, parent := range fields[1:] {
			if named := names.replace(parent); named != parent {
				parents = append(parents, named)
				continue
			}
			parents = append(parents, "("+tree(parent)+")")
		}
		lines = append(lines, fmt.Sprintf("commit %s %s parents=%s", pair.name, tree(pair.raw), strings.Join(parents, ",")))
	}
	return lines
}

// selectionOrder states how the response rows follow the planned selection: as its prefix of
// n rows, or not at all.
func selectionOrder(planned []string, plans []CleanupPlan) string {
	if len(plans) > len(planned) {
		return "diverges"
	}
	for i, plan := range plans {
		if plan.Target != planned[i] {
			return "diverges"
		}
	}
	return fmt.Sprintf("prefix %d of %d", len(plans), len(planned))
}

// durableState is the repository after the apply: each member's tree and ledger state, every
// ref, and every cleanup receipt.
func durableState(t *testing.T, root string, names *baselineNames) []string {
	t.Helper()
	ledger, err := intent.Read(root)
	mustNoError(t, err)
	states := map[string]string{}
	for _, assignment := range ledger.Assignments {
		states[assignment.ID] = string(assignment.State)
	}
	var lines []string
	for _, member := range names.members {
		tree := "present"
		if _, statErr := os.Lstat(member.assignment.Worktree); errors.Is(statErr, os.ErrNotExist) {
			tree = "absent"
		}
		state, recorded := states[member.assignment.ID]
		if !recorded {
			state = "absent"
		}
		lines = append(lines, fmt.Sprintf("member %s tree=%s ledger=%s", member.name, tree, state))
	}
	for _, ref := range baselineRefs(t, root) {
		lines = append(lines, "ref "+ref)
	}
	for _, r := range ledger.CleanupReceipts {
		lines = append(lines, fmt.Sprintf("receipt op=%s target=%s state=%s phase=%s action=%s recovery=%s fingerprint=%s checkpoint=%s branch=%s branch-oid=%s assignment=%s owner=%s owned=%t repo=%s detail=%q",
			r.Operation, r.Target, r.State, r.Phase, r.Action, r.Recovery, r.Fingerprint, r.Checkpoint, r.Branch, r.BranchOID, r.Assignment, r.Owner, r.Owned, r.Repo, r.Detail))
	}
	return lines
}

var (
	baselineDigest = regexp.MustCompile(`\b[0-9a-f]{64}\b`)
	baselineOID    = regexp.MustCompile(`\b[0-9a-f]{40}\b`)
)

// renameSorted renames each line, sorts the lines when their raw order follows random
// identities, and gives each digest or object the case does not know a name by its first
// appearance, so equal unknown values keep equal names.
func renameSorted(names *baselineNames, lines []string, sorted bool) []string {
	renamed := make([]string, len(lines))
	for i, line := range lines {
		renamed[i] = names.replace(line)
	}
	if sorted {
		mask := func(line string) string {
			return baselineOID.ReplaceAllString(baselineDigest.ReplaceAllString(line, "#"), "#")
		}
		sort.SliceStable(renamed, func(a, b int) bool { return mask(renamed[a]) < mask(renamed[b]) })
	}
	return renamed
}

// nameUnknownValues renames each digest and object no mapping names, by first appearance.
func nameUnknownValues(record string) string {
	for _, kind := range []struct {
		pattern *regexp.Regexp
		prefix  string
	}{{baselineDigest, "digest"}, {baselineOID, "oid"}} {
		seen := map[string]string{}
		record = kind.pattern.ReplaceAllStringFunc(record, func(value string) string {
			if _, known := seen[value]; !known {
				seen[value] = fmt.Sprintf("<%s:unmapped%d>", kind.prefix, len(seen)+1)
			}
			return seen[value]
		})
	}
	return record
}

// renderBaselineRecord builds one case, drives its adapter once, and renders its record.
func renderBaselineRecord(t *testing.T, c baselineCase) string {
	t.Helper()
	run := c.build(t)
	names := baselineIdentities(t, run)
	before := baselineRefs(t, run.root)
	relations := commitRelations(t, run.root, names, run.byRole)
	plans, applyErr := run.apply()

	rows := make([]string, 0, len(plans))
	for _, plan := range plans {
		rows = append(rows, fmt.Sprintf("row %s | reason-code=%s | assignment=%s", strings.Join(cleanupRow(plan), " | "), plan.ReasonCode, plan.Assignment))
	}
	errText := "none"
	if applyErr != nil {
		errText = applyErr.Error()
	}
	var b strings.Builder
	section := func(title string, lines []string) {
		fmt.Fprintf(&b, "%s\n", title)
		// A multi-line Git error keeps one record line, so each line stays one fact.
		for _, line := range lines {
			fmt.Fprintf(&b, "  %s\n", strings.ReplaceAll(line, "\n", `\n`))
		}
	}
	fmt.Fprintf(&b, "case %s/%s\n", c.mode, c.name)
	section("inputs", renameSorted(names, run.inputs, run.byRole))
	section("commits", renameSorted(names, relations, true))
	section("refs-before", renameSorted(names, prefixed("ref ", before), true))
	section("response", append([]string{
		fmt.Sprintf("error %s", names.replace(errText)),
		fmt.Sprintf("stale %t", errors.Is(applyErr, errStaleFingerprint)),
		fmt.Sprintf("rows %d nil=%t order=%s", len(plans), plans == nil, selectionOrder(run.planned, plans)),
	}, renameSorted(names, rows, run.byRole)...))
	section("after", renameSorted(names, durableState(t, run.root, names), true))
	return nameUnknownValues(b.String())
}

// prefixed puts prefix ahead of each line.
func prefixed(prefix string, lines []string) []string {
	out := make([]string, len(lines))
	for i, line := range lines {
		out[i] = prefix + line
	}
	return out
}
