package worktree

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

var setFingerprint = regexp.MustCompile(`[0-9a-f]{64}`)

// cleanupRows returns the rendered worktree_cleanup rows alone. A cleanup response can
// carry an ignored preview and a help block after the table, and both indent their own
// rows the same way, so the reader tracks which table it is inside.
func cleanupRows(output string) []string {
	var rows []string
	inTable := false
	for _, line := range strings.Split(output, "\n") {
		switch {
		case strings.HasPrefix(line, "worktree_cleanup["):
			inTable = true
		case line == "":
		case strings.HasPrefix(line, "  "):
			if inTable {
				rows = append(rows, strings.TrimPrefix(line, "  "))
			}
		default:
			inTable = false
		}
	}
	return rows
}

// cleanupRowFields splits one rendered row. Every fixture path and detail below is
// comma-free, so the split addresses the same fields the table header names.
func cleanupRowFields(row string) []string { return strings.Split(row, ",") }

func cleanupRowsField(t *testing.T, output string, index int) []string {
	t.Helper()
	rows := cleanupRows(output)
	values := make([]string, 0, len(rows))
	for _, row := range rows {
		fields := cleanupRowFields(row)
		if index >= len(fields) {
			t.Fatalf("row %q has no field %d", row, index)
		}
		values = append(values, fields[index])
	}
	return values
}

func TestCleanExplicitSetPlan(t *testing.T) {
	t.Parallel()
	root, home, first, second, _ := landedSetFixture(t)
	stdout, stderr, code := runCleanup(t, root, home, "--target", first.Assignment.ID, "--target", second.Assignment.ID)
	if code != 0 || stderr != "" {
		t.Fatalf("set plan exit=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	rows := cleanupRows(stdout)
	if len(rows) != 2 {
		t.Fatalf("set plan rows = %#v, want one row per selected target", rows)
	}
	shared := setFingerprint.FindAllString(stdout, -1)
	if len(shared) < 2 || shared[0] != shared[1] {
		t.Fatalf("set plan = %q, want one shared set fingerprint", stdout)
	}
	if !strings.Contains(stdout, "bench worktree clean --target ") || !strings.Contains(stdout, "--apply "+shared[0]) {
		t.Fatalf("set plan = %q, want the whole-set apply action", stdout)
	}
	for _, creation := range []Creation{first, second} {
		if _, err := os.Stat(creation.Path); err != nil {
			t.Fatalf("bare set plan removed %s: %v", creation.Path, err)
		}
	}

	narrow, narrowErr, narrowCode := runCleanup(t, root, home, "--target", first.Assignment.ID)
	if narrowCode != 0 || narrowErr != "" {
		t.Fatalf("narrow plan exit=%d stdout=%q stderr=%q", narrowCode, narrow, narrowErr)
	}
	if setFingerprint.FindString(narrow) == shared[0] {
		t.Fatalf("set fingerprint did not bind complete membership: pair=%q one=%q", stdout, narrow)
	}

	applied, applyErr, applyCode := runCleanup(t, root, home, "--target", first.Assignment.ID, "--target", second.Assignment.ID, "--apply", shared[0])
	if applyCode != 0 || applyErr != "" || strings.Count(applied, ",removed,") != 2 {
		t.Fatalf("set apply = (%d, %q, %q), want two removals", applyCode, applied, applyErr)
	}
	for _, creation := range []Creation{first, second} {
		if _, err := os.Lstat(creation.Path); !os.IsNotExist(err) {
			t.Fatalf("set apply left %s: %v", creation.Path, err)
		}
	}
}

func TestCleanExplicitSetAliases(t *testing.T) {
	t.Parallel()
	root, home, first, second, _ := landedSetFixture(t)
	aliases := []string{
		"--target", first.Assignment.ID, "--target", first.Path, "--target", first.Assignment.Label,
		"--target", second.Assignment.Label, "--target", second.Assignment.ID,
	}
	stdout, stderr, code := runCleanup(t, root, home, aliases...)
	if code != 0 || stderr != "" {
		t.Fatalf("alias plan exit=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	rows := cleanupRows(stdout)
	if len(rows) != 2 {
		t.Fatalf("alias plan rows = %#v, want one row per cleanup identity", rows)
	}
	lower, higher := first, second
	if higher.Assignment.ID < lower.Assignment.ID {
		lower, higher = second, first
	}
	if targets := cleanupRowsField(t, stdout, 0); targets[0] != lower.Path || targets[1] != higher.Path {
		t.Fatalf("alias plan targets = %#v, want canonical identity order", targets)
	}
	canonical, _, canonicalCode := runCleanup(t, root, home, "--target", first.Assignment.ID, "--target", second.Assignment.ID)
	if canonicalCode != 0 || setFingerprint.FindString(stdout) != setFingerprint.FindString(canonical) {
		t.Fatalf("alias plan = %q and canonical plan = %q, want one identity-ordered fingerprint", stdout, canonical)
	}

	applied, applyErr, applyCode := runCleanup(t, root, home, append(aliases, "--apply", setFingerprint.FindString(stdout))...)
	if applyCode != 0 || applyErr != "" || strings.Count(applied, ",removed,") != 2 || len(cleanupRows(applied)) != 2 {
		t.Fatalf("alias apply = (%d, %q, %q), want one removal per identity", applyCode, applied, applyErr)
	}
	for _, creation := range []Creation{first, second} {
		if _, err := os.Lstat(creation.Path); !os.IsNotExist(err) {
			t.Fatalf("alias apply left %s: %v", creation.Path, err)
		}
	}
}

func TestCleanExplicitSetSelectionFailure(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		target func(*testing.T, string, string) string
		detail string
	}{
		{name: "unassigned", target: func(*testing.T, string, string) string { return "absent-selection-target" }, detail: "target is unassigned"},
		{name: "ambiguous", target: func(t *testing.T, root, home string) string {
			mustCreate(t, root, home, "ambiguous-alias-one", "shared alias")
			mustCreate(t, root, home, "ambiguous-alias-two", "shared alias")
			return "shared alias"
		}, detail: "target is ambiguous"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root, home, first, second, _ := landedSetFixture(t)
			failing := tc.target(t, root, home)
			stdout, stderr, code := runCleanup(t, root, home, "--target", first.Assignment.ID, "--target", failing)
			if code != 1 || stderr != "" {
				t.Fatalf("failed selection exit=%d stdout=%q stderr=%q", code, stdout, stderr)
			}
			if rows := cleanupRows(stdout); len(rows) != 2 {
				t.Fatalf("failed selection rows = %#v, want every selection outcome", rows)
			}
			if !strings.Contains(stdout, tc.detail) || !strings.Contains(stdout, ",error,") {
				t.Fatalf("failed selection = %q, want the reported failure", stdout)
			}
			if setFingerprint.MatchString(stdout) || strings.Contains(stdout, "bench worktree clean --target") {
				t.Fatalf("failed selection = %q, want no applicable fingerprint and no apply action", stdout)
			}
			_, applyStderr, applyCode := runCleanup(t, root, home, "--target", first.Assignment.ID, "--target", failing, "--apply", strings.Repeat("a", 64))
			if applyCode != 1 || applyStderr != "" {
				t.Fatalf("failed selection apply exit=%d stderr=%q, want a refusal", applyCode, applyStderr)
			}
			for _, creation := range []Creation{first, second} {
				if _, err := os.Stat(creation.Path); err != nil {
					t.Fatalf("failed selection removed %s: %v", creation.Path, err)
				}
			}
		})
	}
}

func TestCleanSetRetainsAuthority(t *testing.T) {
	t.Parallel()
	root, home, first, _, dirty := landedSetFixture(t)
	active := mustCreate(t, root, home, "set-authority-active", "active label")
	commitInWorktree(t, active.Path, "active.txt", "active\n", "active work")
	lease, err := LeaseFile(active.Path)
	mustNoError(t, err)
	mustWrite(t, lease, []byte(strconv.Itoa(os.Getpid())+" 2026-07-15T00:00:00Z\n"), 0o600)
	members := []Creation{first, dirty, active}

	// Every explicit target keeps the verdict the single-target form reaches for it, so a
	// set selection cannot widen what the command is allowed to remove.
	verdicts := make([]string, 0, len(members))
	for _, member := range members {
		single, singleErr, singleCode := runCleanup(t, root, home, member.Path)
		if singleCode != 0 || singleErr != "" {
			t.Fatalf("single-target plan for %s exit=%d stdout=%q stderr=%q", member.Path, singleCode, single, singleErr)
		}
		set, setErr, setCode := runCleanup(t, root, home, "--target", member.Assignment.ID)
		if setCode != 0 || setErr != "" {
			t.Fatalf("set plan for %s exit=%d stdout=%q stderr=%q", member.Path, setCode, set, setErr)
		}
		singleRow, setRow := cleanupRowFields(cleanupRows(single)[0]), cleanupRowFields(cleanupRows(set)[0])
		singleRow[5], setRow[5] = "fingerprint", "fingerprint"
		if strings.Join(singleRow, ",") != strings.Join(setRow, ",") {
			t.Fatalf("set row = %q, want the single-target verdict %q", setRow, singleRow)
		}
		verdicts = append(verdicts, setRow[1])
	}
	if verdicts[2] != string(ActionRetain) {
		t.Fatalf("a live-leased target planned %q, want the existing refusal", verdicts[2])
	}

	args := []string{"--target", first.Assignment.ID, "--target", dirty.Assignment.ID, "--target", active.Assignment.ID}
	stdout, stderr, code := runCleanup(t, root, home, args...)
	if code != 0 || stderr != "" {
		t.Fatalf("mixed set plan = (%d, %q, %q), want one applicable plan", code, stdout, stderr)
	}
	sorted := append([]string(nil), verdicts...)
	sort.Strings(sorted)
	planned := cleanupRowsField(t, stdout, 1)
	sort.Strings(planned)
	if strings.Join(planned, ",") != strings.Join(sorted, ",") {
		t.Fatalf("mixed set verdicts = %#v, want each target's own verdict %#v", planned, sorted)
	}
	_, applyErr, applyCode := runCleanup(t, root, home, append(args, "--apply", setFingerprint.FindString(stdout))...)
	if applyCode != 0 || applyErr != "" {
		t.Fatalf("mixed set apply = (%d, %q), want success", applyCode, applyErr)
	}
	if _, statErr := os.Stat(active.Path); statErr != nil {
		t.Fatalf("set apply removed the live-leased target %s: %v", active.Path, statErr)
	}
}

// cleanupEffects is one clean form's normalized lifecycle result: the exit code, every
// rendered verdict, and the durable state each fixture checkout, branch, and assignment
// record is left in. The verdicts sort, because a set orders its rows by an assignment
// identity the fixture draws at random.
func cleanupEffects(t *testing.T, root string, creations []Creation, stdout, stderr string, code int) string {
	t.Helper()
	verdicts := make([]string, 0, len(creations))
	for _, row := range cleanupRows(stdout) {
		fields := cleanupRowFields(row)
		verdicts = append(verdicts, fields[1]+"/"+fields[2])
	}
	sort.Strings(verdicts)
	trees := make([]string, 0, len(creations))
	for _, creation := range creations {
		tree := "absent"
		if _, err := os.Stat(creation.Path); err == nil {
			tree = "present"
		}
		trees = append(trees, fmt.Sprintf("%s/%t", tree, git.OK("-C", root, "show-ref", "--verify", "--quiet", creation.Assignment.Branch)))
	}
	assignments, err := intent.Assignments(root)
	mustNoError(t, err)
	return fmt.Sprintf("exit=%d stderr=%t; rows=%s; trees=%s; assignments=%d",
		code, stderr == "", strings.Join(verdicts, ","), strings.Join(trees, ","), len(assignments))
}

// TestCleanSetCompatibility is the differential that holds every existing cleanup form to
// the lifecycle effects it produced before explicit sets joined the command. The baseline
// below was captured from the unedited tree, so a changed branch, receipt, or retention
// verdict on any existing form turns this row red.
func TestCleanSetCompatibility(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		run  func(*testing.T, string, string, []Creation) (string, string, int)
		want string
	}{
		{
			name: "single target plan",
			run: func(t *testing.T, root, home string, c []Creation) (string, string, int) {
				return runCleanup(t, root, home, c[0].Path)
			},
			want: "exit=0 stderr=true; rows=remove/clean; trees=present/true,present/true,present/true; assignments=3",
		},
		{
			name: "single target apply",
			run: func(t *testing.T, root, home string, c []Creation) (string, string, int) {
				plan, _, _ := runCleanup(t, root, home, c[0].Path)
				return runCleanup(t, root, home, c[0].Path, "--apply", cleanupRowFingerprint(t, plan))
			},
			want: "exit=0 stderr=true; rows=removed/clean; trees=absent/false,present/true,present/true; assignments=2",
		},
		{
			name: "landed plan",
			run: func(t *testing.T, root, home string, _ []Creation) (string, string, int) {
				return runCleanup(t, root, home, "--landed")
			},
			want: "exit=0 stderr=true; rows=remove/clean,remove/clean,retain/dirty; trees=present/true,present/true,present/true; assignments=3",
		},
		{
			name: "landed apply",
			run: func(t *testing.T, root, home string, _ []Creation) (string, string, int) {
				plan, _, _ := runCleanup(t, root, home, "--landed")
				return runCleanup(t, root, home, "--landed", "--apply", cleanupRowFingerprint(t, plan))
			},
			want: "exit=0 stderr=true; rows=removed/clean,removed/clean,retain/dirty; trees=absent/false,absent/false,present/true; assignments=1",
		},
		{
			name: "unclaimed plan",
			run: func(t *testing.T, root, home string, _ []Creation) (string, string, int) {
				return runCleanup(t, root, home, "--discard-branch", "--unclaimed")
			},
			want: "exit=0 stderr=true; rows=discard-remove/unclaimed; trees=present/true,present/true,present/true; assignments=3",
		},
		{
			name: "unclaimed apply",
			run: func(t *testing.T, root, home string, _ []Creation) (string, string, int) {
				plan, _, _ := runCleanup(t, root, home, "--discard-branch", "--unclaimed")
				return runCleanup(t, root, home, "--discard-branch", "--unclaimed", "--apply", cleanupRowFingerprint(t, plan))
			},
			want: "exit=0 stderr=true; rows=removed/unclaimed; trees=present/true,present/true,present/true; assignments=3",
		},
		{
			name: "dirty single target apply",
			run: func(t *testing.T, root, home string, c []Creation) (string, string, int) {
				plan, _, _ := runCleanup(t, root, home, c[2].Path)
				return runCleanup(t, root, home, c[2].Path, "--apply", cleanupRowFingerprint(t, plan))
			},
			want: "exit=0 stderr=true; rows=removed/dirty; trees=present/true,present/true,absent/false; assignments=3",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root, home, first, second, dirty := landedSetFixture(t)
			orphan := intent.AssignmentBranchRef(strings.Repeat("c", 32), strings.Repeat("d", 32))
			gitRun(t, root, "branch", strings.TrimPrefix(orphan, "refs/heads/"))
			creations := []Creation{first, second, dirty}
			stdout, stderr, code := tc.run(t, root, home, creations)
			if got := cleanupEffects(t, root, creations, stdout, stderr, code); got != tc.want {
				t.Fatalf("%s effects = %q, want the captured baseline %q", tc.name, got, tc.want)
			}
		})
	}
}

func TestCleanSetPresentEmptyInventory(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	creation := mustCreate(t, root, home, "present-empty-inventory", "present empty")
	mustNoError(t, intent.DeleteAssignment(root, creation.Assignment.ID))
	ledger, err := intent.Address(root)
	mustNoError(t, err)
	if _, statErr := os.Stat(ledger); statErr != nil {
		t.Fatalf("present-empty fixture has no ledger at %q: %v", ledger, statErr)
	}
	requireEmptyInventoryOutcome(t, root, home, creation)
}

func TestCleanSetAbsentInventory(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	ledger, err := intent.Address(root)
	mustNoError(t, err)
	if _, statErr := os.Lstat(ledger); !os.IsNotExist(statErr) {
		t.Fatalf("absent fixture already holds a ledger at %q: %v", ledger, statErr)
	}
	requireEmptyInventoryOutcome(t, root, home, Creation{Assignment: intent.Assignment{ID: strings.Repeat("e", 32)}})
}

// requireEmptyInventoryOutcome holds both empty landed inventories to the same answer: the
// selector reports an empty plan, an explicit target reports its own failed selection, and
// neither mutates the repository.
func requireEmptyInventoryOutcome(t *testing.T, root, home string, absent Creation) {
	t.Helper()
	before, err := git.Output("-C", root, "rev-parse", "HEAD")
	mustNoError(t, err)
	stdout, stderr, code := runCleanup(t, root, home, "--landed")
	if code != 0 || stderr != "" || !strings.HasPrefix(stdout, "worktree_cleanup[0]") || len(cleanupRows(stdout)) != 0 {
		t.Fatalf("landed plan = (%d, %q, %q), want an empty plan", code, stdout, stderr)
	}
	targeted, targetedErr, targetedCode := runCleanup(t, root, home, "--target", absent.Assignment.ID)
	if targetedCode != 1 || targetedErr != "" || len(cleanupRows(targeted)) != 1 || setFingerprint.MatchString(targeted) {
		t.Fatalf("explicit plan = (%d, %q, %q), want one unapplicable selection outcome", targetedCode, targeted, targetedErr)
	}
	assignments, err := intent.Assignments(root)
	mustNoError(t, err)
	after, err := git.Output("-C", root, "rev-parse", "HEAD")
	mustNoError(t, err)
	if len(assignments) != 0 || after != before {
		t.Fatalf("empty inventory mutated the repository: assignments=%#v head %q -> %q", assignments, before, after)
	}
}
