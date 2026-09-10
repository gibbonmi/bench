package roadmap

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/otelrecord"
	"github.com/gibbonmi/bench/internal/retros"
	"go.opentelemetry.io/otel/attribute"
)

// wantScaffoldHeadings is the canonical heading list spelled independently of the
// renderer's own source. A heading swapped in the exported list reds here rather than
// passing a parser that grades order alone.
var wantScaffoldHeadings = []string{
	"## Outcome",
	"## Gate-stage timings",
	"## Ticket-versus-spec-slice and delegate performance",
	"## Coordinator catches",
	"## Repair attribution",
	"## Agent-experience improvements",
	"### Bench CLI",
	"### Skills",
	"### Process",
}

// newScaffoldRepo returns a temporary repository with its own Bench home bound, so the
// scaffold reads the record this test writes and no other.
func newScaffoldRepo(t *testing.T) (root, home string) {
	t.Helper()
	root = newRepo(t)
	home = t.TempDir()
	t.Setenv("BENCH_HOME", home)
	return root, home
}

// recordLanding appends one complete landing trace: the landing span, its published
// subject, and one gate phase below it. The phase runs for hold, so its recorded
// elapsed time has a floor this test can assert.
func recordLanding(home, root, subject, stage string, hold time.Duration) {
	landing, span, finish := otelrecord.BeginIn(context.Background(), home, root, otelrecord.SeamLanding, otelrecord.SeamLanding)
	span.SetAttributes(attribute.String(otelrecord.AttrSubjectID, subject))
	_, _, finishPhase := otelrecord.BeginIn(landing, home, root, otelrecord.SeamGatePhase, stage)
	time.Sleep(hold)
	finishPhase()
	finish()
}

func writeTickets(t *testing.T, root, slug string, names ...string) {
	t.Helper()
	dir := filepath.Join(root, "specs", slug, "tickets")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range names {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("Ticket.\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestRetroScaffoldParses(t *testing.T) {
	newScaffoldRepo(t)
	body, code := RetroCommand([]string{"shape", "--scaffold"})
	if code != 0 {
		t.Fatalf("scaffold exit = %d, want 0: %q", code, body)
	}
	if err := retros.Parse([]byte(body)); err != nil {
		t.Fatalf("scaffold body does not parse: %v\n%s", err, body)
	}
	position := 0
	for _, heading := range wantScaffoldHeadings {
		offset := strings.Index(body[position:], heading+"\n")
		if offset < 0 {
			t.Fatalf("scaffold body carries no %q heading in order:\n%s", heading, body)
		}
		position += offset + len(heading)
	}
}

func TestRetroScaffoldNamesTheLandingStages(t *testing.T) {
	root, home := newScaffoldRepo(t)
	recordLanding(home, root, "older-commit", "older-build", 0)
	recordLanding(home, root, "newer-commit", "newer-build", 25*time.Millisecond)
	// The operator's own gate run shares the record and belongs to no landing.
	foreign, _, finishForeign := otelrecord.BeginIn(context.Background(), home, root, "gate", "gate.ordinary")
	_, _, finishForeignPhase := otelrecord.BeginIn(foreign, home, root, otelrecord.SeamGatePhase, "foreign-build")
	finishForeignPhase()
	finishForeign()

	body, code := RetroCommand([]string{"stages", "--scaffold"})
	if code != 0 {
		t.Fatalf("scaffold exit = %d, want 0: %q", code, body)
	}
	timings := sectionOf(t, body, retros.TimingsHeading)
	if !strings.Contains(timings, "newer-commit") {
		t.Fatalf("timings name no published commit:\n%s", timings)
	}
	if strings.Contains(timings, "older-build") || strings.Contains(timings, "foreign-build") {
		t.Fatalf("timings list a stage outside the newest landing's trace:\n%s", timings)
	}
	stage := stageLine(t, timings, "newer-build")
	elapsed := stageElapsed(t, stage)
	if elapsed < 20 {
		t.Fatalf("stage line %q reports %d ms, want at least the 25ms the phase ran", stage, elapsed)
	}
}

func TestRetroScaffoldReportsUnknownTimings(t *testing.T) {
	for _, name := range []string{"absent", "unreadable", "no-landing"} {
		t.Run(name, func(t *testing.T) {
			root, home := newScaffoldRepo(t)
			switch name {
			case "unreadable":
				if os.Geteuid() == 0 {
					capability.Capability(t, capability.Privilege, "root reads mode 0000 files; the unreadable record is unobservable")
				}
				recordLanding(home, root, "commit", "build", 0)
				record := otelrecord.Path(home, root)
				if err := os.Chmod(record, 0o000); err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() { _ = os.Chmod(record, 0o600) })
			case "no-landing":
				_, _, finish := otelrecord.BeginIn(context.Background(), home, root, "gate", "gate.ordinary")
				finish()
			}
			body, code := RetroCommand([]string{"missing", "--scaffold"})
			if code != 0 {
				t.Fatalf("scaffold exit = %d, want 0: %q", code, body)
			}
			if got := strings.TrimSpace(sectionOf(t, body, retros.TimingsHeading)); got != "unknown" {
				t.Fatalf("%s timings = %q, want unknown", name, got)
			}
		})
	}
}

func TestRetroScaffoldListsTheTickets(t *testing.T) {
	root, _ := newScaffoldRepo(t)
	writeTickets(t, root, "sliced", "2.md", "10.md", "1.md")
	body, code := RetroCommand([]string{"sliced", "--scaffold"})
	if code != 0 {
		t.Fatalf("scaffold exit = %d, want 0: %q", code, body)
	}
	want := "| ticket | rounds | causes |\n|---|---|---|\n| 1.md | unknown | unknown |\n| 10.md | unknown | unknown |\n| 2.md | unknown | unknown |"
	if got := strings.TrimSpace(sectionOf(t, body, retros.RepairHeading)); got != want {
		t.Fatalf("repair table =\n%s\nwant\n%s", got, want)
	}
}

func TestRetroScaffoldListsUnknownWithoutATicketsDirectory(t *testing.T) {
	newScaffoldRepo(t)
	body, code := RetroCommand([]string{"unsliced", "--scaffold"})
	if code != 0 {
		t.Fatalf("scaffold exit = %d, want 0: %q", code, body)
	}
	want := "| ticket | rounds | causes |\n|---|---|---|\n| unknown | unknown | unknown |"
	if got := strings.TrimSpace(sectionOf(t, body, retros.RepairHeading)); got != want {
		t.Fatalf("repair table =\n%s\nwant\n%s", got, want)
	}
}

func TestRetroScaffoldWritesNoFile(t *testing.T) {
	root, _ := newScaffoldRepo(t)
	body, code := RetroCommand([]string{"later", "--scaffold"})
	if code != 0 {
		t.Fatalf("scaffold exit = %d, want 0: %q", code, body)
	}
	path := filepath.Join(root, retros.Path("later"))
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("scaffold created %s: %v", path, err)
	}
	if out, code := RetroCommand([]string{"later", "--body", body}); code != 0 || out != "captured: later\n" {
		t.Fatalf("capture after scaffold = %q/%d, want captured on exit 0", out, code)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("capture after scaffold wrote no file: %v", err)
	}
}

func TestRetroScaffoldProposesNoImprovement(t *testing.T) {
	root, _ := newScaffoldRepo(t)
	body, code := RetroCommand([]string{"quiet", "--scaffold"})
	if code != 0 {
		t.Fatalf("scaffold exit = %d, want 0: %q", code, body)
	}
	if _, code := RetroCommand([]string{"quiet", "--body", body}); code != 0 {
		t.Fatalf("capture of the scaffold exit = %d, want 0", code)
	}
	if diags := retros.ValidateImprovementMarkers(root); len(diags) != 0 {
		t.Fatalf("captured scaffold raised %q, want no diagnostic", diags)
	}
}

func TestRetroRefusesBothForms(t *testing.T) {
	root, _ := newScaffoldRepo(t)
	out, code := RetroCommand([]string{"both", "--scaffold", "--body", eligibleRetro(t)})
	if code != 2 || out != retroGrammar.Help+"\n" {
		t.Fatalf("both forms = %q/%d, want the grammar at exit 2", out, code)
	}
	if _, err := os.Stat(filepath.Join(root, retros.Path("both"))); !os.IsNotExist(err) {
		t.Fatalf("refused retro created a file: %v", err)
	}
}

func TestRetroRefusesNeitherForm(t *testing.T) {
	root, _ := newScaffoldRepo(t)
	out, code := RetroCommand([]string{"neither"})
	if code != 2 || out != retroGrammar.Help+"\n" {
		t.Fatalf("neither form = %q/%d, want the grammar at exit 2", out, code)
	}
	if _, err := os.Stat(filepath.Join(root, retros.Path("neither"))); !os.IsNotExist(err) {
		t.Fatalf("refused retro created a file: %v", err)
	}
}

// sectionOf returns the body below one heading, up to the next heading.
func sectionOf(t *testing.T, body, heading string) string {
	t.Helper()
	_, after, found := strings.Cut(body, heading+"\n")
	if !found {
		t.Fatalf("body carries no %q heading:\n%s", heading, body)
	}
	var kept []string
	for _, line := range strings.Split(after, "\n") {
		if strings.HasPrefix(line, "#") {
			break
		}
		kept = append(kept, line)
	}
	return strings.Join(kept, "\n")
}

func stageLine(t *testing.T, timings, stage string) string {
	t.Helper()
	for _, line := range strings.Split(timings, "\n") {
		if strings.Contains(line, stage) {
			return line
		}
	}
	t.Fatalf("timings name no %q stage:\n%s", stage, timings)
	return ""
}

// stageElapsed reads the millisecond count a stage line reports.
func stageElapsed(t *testing.T, line string) int {
	t.Helper()
	fields := strings.Fields(line)
	if len(fields) < 2 || fields[len(fields)-1] != "ms" {
		t.Fatalf("stage line %q states no millisecond count", line)
	}
	count := 0
	for _, digit := range fields[len(fields)-2] {
		if digit < '0' || digit > '9' {
			t.Fatalf("stage line %q states a non-numeric elapsed time", line)
		}
		count = count*10 + int(digit-'0')
	}
	return count
}
