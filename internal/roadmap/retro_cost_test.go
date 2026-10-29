package roadmap

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/otelrecord"
	"github.com/gibbonmi/bench/internal/retros"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
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
	return r.gitAuthored(at, at, args...)
}

// gitAuthored runs git with the author time authored and the committer time at.
func (r costRepo) gitAuthored(authored, at int, args ...string) string {
	r.t.Helper()
	cmd := exec.Command("git", append([]string{"-C", r.root}, args...)...)
	stamp := func(offset int) string { return "@" + strconv.Itoa(costEpoch+offset) + " +0000" }
	cmd.Env = append(os.Environ(), "GIT_AUTHOR_DATE="+stamp(authored), "GIT_COMMITTER_DATE="+stamp(at))
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
	gittest.MustRun(r.t, r.root, "add", "-A")
	r.git(at, "commit", "-q", "--allow-empty", "-m", subject)
}

// land merges branch into main with the spec flip in the merge commit, as a landing does,
// and returns the landing commit.
func (r costRepo) land(at int, branch, slug string) string {
	r.t.Helper()
	gittest.MustRun(r.t, r.root, "checkout", "-q", "main")
	r.git(at, "merge", "-q", "--no-ff", "--no-commit", branch)
	path := filepath.Join(r.root, "specs", slug, "spec.md")
	data, err := os.ReadFile(path)
	if err != nil {
		r.t.Fatal(err)
	}
	r.write(map[string]string{"specs/" + slug + "/spec.md": strings.Replace(string(data), "Status: staged", "Status: implemented", 1)})
	gittest.MustRun(r.t, r.root, "add", "-A")
	r.git(at, "commit", "-q", "-m", "land "+branch)
	return gittest.Output(r.t, r.root, "rev-parse", "HEAD")
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
//     fold merge of main, and the landing L at 960 minutes. The first build commit has an
//     author time 30 minutes before its committer time, and one bookkeeping commit
//     changes a spec-folder path that git quotes;
//   - the landing of slug o with a version 1 plan;
//   - the retirement of slug s, which deletes its folder.
func newCostRepo(t *testing.T) (costRepo, string) {
	t.Helper()
	r := costRepo{t: t, root: gittest.RepoOnBranch(t, "main")}
	r.commit(-7200, "base", map[string]string{"README.md": "base\n"})
	r.git(-7200, "branch", "build")
	r.commit(-3600, "fold source", map[string]string{"main.txt": "fold\n"})
	gittest.MustRun(r.t, r.root, "checkout", "-q", "build")
	plan := `{"version":2,"execution":{"mode":"delegate","run_id":"r","orchestrator_session":"o","author_limit":1,"assignments":{` +
		`"a.md":[{"session":"a1","assignment":"asg-a"},{"session":"a2","assignment":"asg-a","predecessor":"a1"}],` +
		`"b.md":[{"session":"b1","assignment":"asg-b"}],"c.md":[{"session":"c1"},{"session":"c2","predecessor":"c1"}],"d.md":[]}},` +
		`"chunks":[{"id":"X","tickets":["a.md","b.md"]},{"id":"Y","tickets":["c.md"]},{"id":"Z","tickets":["d.md"]}],` +
		`"final_verification":[{"id":"g","command":"bench gate"}]}`
	r.write(map[string]string{
		"specs/s/spec.md":      costSpec(plan),
		"specs/s/tickets/a.md": costTicket("none"),
		"specs/s/tickets/b.md": costTicket("a.md"),
		"specs/s/tickets/c.md": costTicket("b.md"),
		"specs/s/tickets/d.md": costTicket("none"),
	})
	gittest.MustRun(r.t, r.root, "add", "-A")
	r.gitAuthored(-1800, 0, "commit", "-q", "-m", "spec")
	r.commit(600, "change", map[string]string{"internal/x.go": "package x\n"})
	r.git(1200, "merge", "-q", "--no-ff", "-m", "fold main", "main")
	r.commit(1800, "record", map[string]string{"reviews/s.md": "record\n", `specs/s/n"q".md`: "quoted\n"})
	r.commit(2400, "mixed", map[string]string{"internal/x.go": "package x // mixed\n", "specs/s/spec.md": costSpec(plan) + "More.\n"})
	r.commit(3000, "ticket", map[string]string{"specs/s/tickets/a.md": costTicket("none") + "More.\n"})
	landing := r.land(960*60, "build", "s")

	gittest.MustRun(r.t, r.root, "checkout", "-q", "-b", "other")
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
	home := t.TempDir()
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
			body := scaffoldBody(home, r.root, tc.slug)
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
		{name: "wall over", tickets: 20, wall: 1300 * 60, booked: 140, want: []string{
			"- wall per ticket: 65.0 min against 60 min, over budget",
			"- bookkeeping per ticket: 7.0 against 7.1",
		}},
		{name: "bookkeeping over", tickets: 20, wall: 1200 * 60, booked: 144, want: []string{
			"- wall per ticket: 60.0 min against 60 min",
			"- bookkeeping per ticket: 7.2 against 7.1, over budget",
		}},
		{name: "fractional minutes", tickets: 16, wall: 57888, booked: 0, want: []string{
			"- wall per ticket: 60.3 min against 60 min, over budget",
			"- bookkeeping per ticket: 0.0 against 7.1",
		}},
		{name: "rounded to the budget", tickets: 50, wall: 3002 * 60, booked: 357, want: []string{
			"- wall per ticket: 60.0 min against 60 min",
			"- bookkeeping per ticket: 7.1 against 7.1",
		}},
		{name: "bookkeeping rounded up over the budget", tickets: 50, wall: 3000 * 60, booked: 358, want: []string{
			"- wall per ticket: 60.0 min against 60 min",
			"- bookkeeping per ticket: 7.2 against 7.1, over budget",
		}},
		{name: "wall rounded up over the budget", tickets: 50, wall: 3003 * 60, booked: 355, want: []string{
			"- wall per ticket: 60.1 min against 60 min, over budget",
			"- bookkeeping per ticket: 7.1 against 7.1",
		}},
		{name: "zero tickets", tickets: 0, wall: 600, booked: 3, want: []string{
			"- wall per ticket: unknown min against 60 min",
			"- bookkeeping per ticket: unknown against 7.1",
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

// writeCostRecord appends spans to the record of root below home through the record's own
// encoder and writer. trace is the last hex digit of a fixed trace ID.
func writeCostRecord(t *testing.T, home, root string, spans []costSpan) {
	t.Helper()
	writer := otelrecord.NewWriter(home, root)
	start := time.Unix(costEpoch, 0).UTC()
	for i, s := range spans {
		traceID, err := trace.TraceIDFromHex("0102030405060708090a0b0c0d0e0f0" + s.trace)
		if err != nil {
			t.Fatal(err)
		}
		spanID, err := trace.SpanIDFromHex(fmt.Sprintf("%016x", i+1))
		if err != nil {
			t.Fatal(err)
		}
		attrs := []attribute.KeyValue{attribute.String(otelrecord.AttrSeam, s.seam)}
		if s.assignment != "" {
			attrs = append(attrs, attribute.String(otelrecord.AttrAssignmentID, s.assignment))
		}
		line, err := otelrecord.Encode(tracetest.SpanStub{
			Name:        s.seam,
			SpanContext: trace.NewSpanContext(trace.SpanContextConfig{TraceID: traceID, SpanID: spanID, TraceFlags: trace.FlagsSampled}),
			SpanKind:    trace.SpanKindInternal,
			StartTime:   start,
			EndTime:     start.Add(time.Duration(s.ms) * time.Millisecond),
			Attributes:  attrs,
		}.Snapshot())
		if err != nil {
			t.Fatal(err)
		}
		if err := writer.Append(line); err != nil {
			t.Fatal(err)
		}
	}
}

func TestRetroCostTraceTimes(t *testing.T) {
	r, _ := newCostRepo(t)
	planned := []costSpan{
		{trace: "1", seam: otelrecord.SeamCommit, assignment: "asg-a", ms: 100},
		{trace: "1", seam: otelrecord.SeamCommit, ms: 200},
		{trace: "1", seam: otelrecord.SeamLane, ms: 300},
		{trace: "1", seam: otelrecord.SeamGatePhase, ms: 50},
		{trace: "2", seam: otelrecord.SeamLanding, assignment: "asg-b", ms: 400},
	}
	foreign := []costSpan{
		{trace: "3", seam: otelrecord.SeamCommit, assignment: "asg-other", ms: 1000},
		{trace: "4", seam: otelrecord.SeamCommit, ms: 2000},
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
			if tc.spans != nil {
				writeCostRecord(t, home, r.root, tc.spans)
			}
			body := scaffoldBody(home, r.root, tc.slug)
			if got := stageLine(t, sectionOf(t, body, retros.CostHeading), "- trace:"); got != tc.want {
				t.Fatalf("trace line = %q, want %q", got, tc.want)
			}
		})
	}
}
