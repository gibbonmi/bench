package roadmap

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/otelrecord"
	"github.com/gibbonmi/bench/internal/retros"
)

// costEpoch is the committer time of the first build commit of the fixture. Every other
// commit time is an offset from it, so the wall arithmetic is exact.
const costEpoch = 1_700_000_000

// costRepo is one canned repository with fixed commit times. Each commit sets both the
// author and the committer time, so no expectation reads the wall clock.
type costRepo struct {
	t    *testing.T
	root string
}

func (r costRepo) git(at int, args ...string) string {
	r.t.Helper()
	cmd := exec.Command("git", append([]string{"-C", r.root}, args...)...)
	stamp := "@" + strconv.Itoa(costEpoch+at) + " +0000"
	cmd.Env = append(os.Environ(), "GIT_AUTHOR_DATE="+stamp, "GIT_COMMITTER_DATE="+stamp)
	out, err := cmd.CombinedOutput()
	if err != nil {
		r.t.Fatalf("git %v: %v: %s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func (r costRepo) write(files map[string]string) {
	r.t.Helper()
	for name, body := range files {
		path := filepath.Join(r.root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			r.t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			r.t.Fatal(err)
		}
	}
}

// commit writes files and commits them at the offset at.
func (r costRepo) commit(at int, subject string, files map[string]string) {
	r.t.Helper()
	r.write(files)
	r.git(at, "add", "-A")
	r.git(at, "commit", "-q", "--allow-empty", "-m", subject)
}

// land merges branch into main with the spec flip in the merge commit, as a landing does,
// and returns the landing commit.
func (r costRepo) land(at int, branch, slug string) string {
	r.t.Helper()
	r.git(at, "checkout", "-q", "main")
	r.git(at, "merge", "-q", "--no-ff", "--no-commit", branch)
	path := filepath.Join(r.root, "specs", slug, "spec.md")
	data, err := os.ReadFile(path)
	if err != nil {
		r.t.Fatal(err)
	}
	r.write(map[string]string{"specs/" + slug + "/spec.md": strings.Replace(string(data), "Status: staged", "Status: implemented", 1)})
	r.git(at, "add", "-A")
	r.git(at, "commit", "-q", "-m", "land "+branch)
	return r.git(at, "rev-parse", "HEAD")
}

func costSpec(plan string) string {
	return "# Spec\n\nStatus: staged\n\n```bench-completion-plan\n" + plan + "\n```\n"
}

func costTicket(blocker string) string {
	return "# Ticket\n\nBlocked by: " + blocker + "\nWrites: internal/x.go\nCovers: BC1\n\n## What to build\n\nIt.\n\n## Acceptance\n\n- [ ] It.\n"
}

// newCostRepo builds the fixture history on main:
//
//   - a fold commit on main 60 minutes before the first build commit;
//   - the build branch of slug s: five non-merge commits, three of them bookkeeping, one
//     fold merge of main, and the landing L at 960 minutes;
//   - the landing of slug o with a version 1 plan;
//   - the retirement of slug s, which deletes its folder.
func newCostRepo(t *testing.T) (costRepo, string) {
	t.Helper()
	r := costRepo{t: t, root: t.TempDir()}
	r.git(-7200, "init", "-q", "-b", "main")
	r.git(-7200, "config", "user.email", "a@b.c")
	r.git(-7200, "config", "user.name", "a")
	r.commit(-7200, "base", map[string]string{"README.md": "base\n"})
	r.git(-7200, "branch", "build")
	r.commit(-3600, "fold source", map[string]string{"main.txt": "fold\n"})
	r.git(0, "checkout", "-q", "build")
	plan := `{"version":2,"execution":{"mode":"delegate","run_id":"r","orchestrator_session":"o","author_limit":1,"assignments":{` +
		`"a.md":[{"session":"a1","assignment":"asg-a"},{"session":"a2","assignment":"asg-a","predecessor":"a1"}],` +
		`"b.md":[{"session":"b1","assignment":"asg-b"}],"c.md":[{"session":"c1"},{"session":"c2","predecessor":"c1"}],"d.md":[]}},` +
		`"chunks":[{"id":"X","tickets":["a.md","b.md"]},{"id":"Y","tickets":["c.md"]},{"id":"Z","tickets":["d.md"]}],` +
		`"final_verification":[{"id":"g","command":"bench gate"}]}`
	r.commit(0, "spec", map[string]string{
		"specs/s/spec.md":      costSpec(plan),
		"specs/s/tickets/a.md": costTicket("none"),
		"specs/s/tickets/b.md": costTicket("a.md"),
		"specs/s/tickets/c.md": costTicket("b.md"),
		"specs/s/tickets/d.md": costTicket("none"),
	})
	r.commit(600, "change", map[string]string{"internal/x.go": "package x\n"})
	r.git(1200, "merge", "-q", "--no-ff", "-m", "fold main", "main")
	r.commit(1800, "record", map[string]string{"reviews/s.md": "record\n"})
	r.commit(2400, "mixed", map[string]string{"internal/x.go": "package x // mixed\n", "specs/s/spec.md": costSpec(plan) + "More.\n"})
	r.commit(3000, "ticket", map[string]string{"specs/s/tickets/a.md": costTicket("none") + "More.\n"})
	landing := r.land(960*60, "build", "s")

	r.git(960*60+60, "checkout", "-q", "-b", "other")
	r.commit(960*60+120, "other spec", map[string]string{
		"specs/o/spec.md":      costSpec(`{"version":1,"chunks":[{"id":"V","tickets":["e.md"]}],"final_verification":[{"id":"g","command":"bench gate"}]}`),
		"specs/o/tickets/e.md": costTicket("none"),
	})
	r.land(960*60+180, "other", "o")
	if err := os.RemoveAll(filepath.Join(r.root, "specs", "s")); err != nil {
		t.Fatal(err)
	}
	r.commit(960*60+240, "spec-retire: s", nil)
	return r, landing
}

func TestRetroCostSection(t *testing.T) {
	r, landing := newCostRepo(t)
	t.Chdir(r.root)
	t.Setenv("BENCH_HOME", t.TempDir())
	for _, tc := range []struct {
		slug, want string
	}{
		{slug: "s", want: strings.Join([]string{
			"- landing: " + landing + ", range " + landing + "^1.." + landing + "^2",
			"- wall: 960 min",
			"- commits: 5 total, 3 bookkeeping, 2 change",
			"- graph: 4 tickets, chain 3, 3 chunks, 2 serial checkpoints",
			"- dispatch: 3 authors, 2 repairs",
			"- trace: unknown",
			"- wall per ticket: 240.0 min against 60 min, over budget",
			"- bookkeeping per ticket: 0.8 against 7.1",
		}, "\n")},
		{slug: "o", want: "- dispatch: unknown authors, unknown repairs"},
		{slug: "missing", want: strings.Join([]string{
			"- landing: unknown, range unknown",
			"- wall: unknown min",
			"- commits: unknown total, unknown bookkeeping, unknown change",
			"- graph: unknown tickets, chain unknown, unknown chunks, unknown serial checkpoints",
			"- dispatch: unknown authors, unknown repairs",
			"- trace: unknown",
			"- wall per ticket: unknown min against 60 min",
			"- bookkeeping per ticket: unknown against 7.1",
		}, "\n")},
	} {
		t.Run(tc.slug, func(t *testing.T) {
			body, code := RetroCommand([]string{tc.slug, "--scaffold"})
			if code != 0 {
				t.Fatalf("scaffold exit = %d, want 0: %q", code, body)
			}
			if !strings.Contains(body, "## Outcome\n\n"+retros.CostHeading+"\n\n- landing: ") {
				t.Fatalf("scaffold prints no Cost section directly after the Outcome heading:\n%s", body)
			}
			got := strings.TrimSpace(sectionOf(t, body, retros.CostHeading))
			if tc.slug == "o" {
				got = stageLine(t, got, "- dispatch:")
			}
			if got != tc.want {
				t.Fatalf("cost section =\n%s\nwant\n%s", got, tc.want)
			}
		})
	}
}

func TestRetroCostBudget(t *testing.T) {
	for _, tc := range []struct {
		name                  string
		tickets, wall, booked int
		want                  []string
	}{
		{name: "wall over", tickets: 20, wall: 1300, booked: 140, want: []string{
			"- wall per ticket: 65.0 min against 60 min, over budget",
			"- bookkeeping per ticket: 7.0 against 7.1",
		}},
		{name: "bookkeeping over", tickets: 20, wall: 1200, booked: 144, want: []string{
			"- wall per ticket: 60.0 min against 60 min",
			"- bookkeeping per ticket: 7.2 against 7.1, over budget",
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := budgetLines(costFacts{tickets: tc.tickets, wall: tc.wall, bookkeeping: tc.booked, ranged: true, graphed: true})
			if strings.Join(got, "\n") != strings.Join(tc.want, "\n") {
				t.Fatalf("budget lines =\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(tc.want, "\n"))
			}
		})
	}
}

// costSpan is one canned finished span of the record. start and the elapsed ms are fixed,
// so the sums are exact and no expectation reads the wall clock.
type costSpan struct {
	trace, seam, assignment string
	ms                      int
}

// writeCostRecord writes spans as the live record of root below home, one OTLP line each.
func writeCostRecord(t *testing.T, home, root string, spans []costSpan) {
	t.Helper()
	var lines []string
	for i, s := range spans {
		attrs := []map[string]any{{"key": otelrecord.AttrSeam, "value": map[string]string{"stringValue": s.seam}}}
		if s.assignment != "" {
			attrs = append(attrs, map[string]any{"key": otelrecord.AttrAssignmentID, "value": map[string]string{"stringValue": s.assignment}})
		}
		start := int64(costEpoch) * 1e9
		line, err := json.Marshal(map[string]any{"resourceSpans": []any{map[string]any{"resource": map[string]any{}, "scopeSpans": []any{map[string]any{"scope": map[string]any{}, "spans": []any{map[string]any{
			"traceId": s.trace, "spanId": strconv.Itoa(i + 1), "name": s.seam, "kind": 1,
			"startTimeUnixNano": strconv.FormatInt(start, 10), "endTimeUnixNano": strconv.FormatInt(start+int64(s.ms)*1e6, 10),
			"attributes": attrs, "status": map[string]int{"code": 0},
		}}}}}}})
		if err != nil {
			t.Fatal(err)
		}
		lines = append(lines, string(line))
	}
	path := otelrecord.Path(home, root)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestRetroCostTraceTimes(t *testing.T) {
	r, _ := newCostRepo(t)
	t.Chdir(r.root)
	root := r.git(0, "rev-parse", "--show-toplevel")
	planned := []costSpan{
		{trace: "t1", seam: "commit", assignment: "asg-a", ms: 100},
		{trace: "t1", seam: "commit", ms: 200},
		{trace: "t1", seam: "lane", ms: 300},
		{trace: "t1", seam: otelrecord.SeamGatePhase, ms: 50},
		{trace: "t2", seam: otelrecord.SeamLanding, assignment: "asg-b", ms: 400},
	}
	foreign := []costSpan{
		{trace: "t3", seam: "commit", assignment: "asg-other", ms: 1000},
		{trace: "t4", seam: "commit", ms: 2000},
	}
	for _, tc := range []struct {
		name, slug string
		spans      []costSpan
		want       string
	}{
		{name: "plan spans", slug: "s", spans: append(append([]costSpan{}, planned...), foreign...), want: "- trace: commit 2 spans in 300 ms, lane 1 spans in 300 ms, landing 1 spans in 400 ms"},
		{name: "absent record", slug: "s", want: "- trace: unknown"},
		{name: "no plan span", slug: "s", spans: foreign, want: "- trace: unknown"},
		{name: "version 1 plan", slug: "o", spans: planned, want: "- trace: unknown"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			t.Setenv("BENCH_HOME", home)
			if tc.spans != nil {
				writeCostRecord(t, home, root, tc.spans)
			}
			body, code := RetroCommand([]string{tc.slug, "--scaffold"})
			if code != 0 {
				t.Fatalf("scaffold exit = %d, want 0: %q", code, body)
			}
			if got := stageLine(t, sectionOf(t, body, retros.CostHeading), "- trace:"); got != tc.want {
				t.Fatalf("trace line = %q, want %q", got, tc.want)
			}
		})
	}
}
