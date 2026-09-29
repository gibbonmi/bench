package spec

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// historyFields is the schema of the positional history table.
var historyFields = []string{"hash", "date", "kind", "subject"}

type historyFixtureData struct {
	root     string
	hash     map[string]string
	baseline []historyBaseline
}

// historyFixture builds one repository whose histories cover the producer matrix:
// delete-only, retire-only, mixed flat and folder deletions, a commit that both retires
// and deletes, a prefix slug, and an empty history. Each commit has its own date.
func historyFixture(t *testing.T) historyFixtureData {
	t.Helper()
	baseline := readHistoryBaseline(t)
	f := emptyHistoryRepo(t)
	f.baseline = baseline
	root := f.root
	writeFolderSpec(t, root, "only-delete", "delete-only\n")
	if err := os.WriteFile(filepath.Join(root, "specs", "mixed.md"), []byte("flat\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	f.commit(t, "base", 1, "initial specs")
	f.commit(t, "retire", 2, "spec-retire: retire-only")
	runGit(t, root, "rm", "-q", "specs/mixed.md")
	f.commit(t, "both", 3, "spec-retire: mixed")
	runGit(t, root, "rm", "-q", "specs/only-delete/spec.md")
	f.commit(t, "delete", 4, "remove obsolete spec")
	writeFolderSpec(t, root, "mixed", "folder\n")
	f.commit(t, "folder", 5, "add spec folder")
	runGit(t, root, "rm", "-q", "specs/mixed/spec.md")
	f.commit(t, "mixed-delete", 6, "remove mixed spec")
	f.commit(t, "prefix", 7, "spec-retire: mixed-extra")
	f.commit(t, "newest", 8, "quoted \"café\", tab\tvalue spec-retire: mixed")
	return f
}

// emptyHistoryRepo initializes a repository and makes it the working directory.
func emptyHistoryRepo(t *testing.T) historyFixtureData {
	t.Helper()
	root := gittest.RepoOnBranch(t, "main")
	t.Chdir(root)
	return historyFixtureData{root: root, hash: map[string]string{}}
}

// commit stages every change and commits it at noon UTC on the given January 2025 day.
func (f historyFixtureData) commit(t *testing.T, name string, day int, subject string) {
	t.Helper()
	date := fmt.Sprintf("2025-01-%02dT12:00:00+00:00", day)
	t.Setenv("GIT_AUTHOR_DATE", date)
	t.Setenv("GIT_COMMITTER_DATE", date)
	runGit(t, f.root, "add", "-A")
	runGit(t, f.root, "commit", "--allow-empty", "-qm", subject)
	hash, err := git.Output("-C", f.root, "rev-parse", "--short", "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	f.hash[name] = hash
}

// expand replaces the fixture root and commit-name placeholders of a recorded value.
func (f historyFixtureData) expand(value string) string {
	value = strings.ReplaceAll(value, "{{ROOT}}", f.root)
	for name, hash := range f.hash {
		value = strings.ReplaceAll(value, "{{"+name+"}}", hash)
	}
	return value
}

// capturedRows returns the recorded positional rows of one baseline case.
func (f historyFixtureData) capturedRows(t *testing.T, name string) [][]string {
	t.Helper()
	for _, capture := range f.baseline {
		if capture.Name != name {
			continue
		}
		rows := make([][]string, len(capture.Rows))
		for i, row := range capture.Rows {
			for _, value := range row {
				rows[i] = append(rows[i], f.expand(value))
			}
		}
		return rows
	}
	t.Fatalf("missing baseline %q", name)
	return nil
}

// historyBaseline is one positional invocation that the unchanged history command
// answered. Exactly one of Rows, History, NoRepository, Help, MissingArgument, and
// UnknownArgument names the expected answer; History reuses another case's rows.
type historyBaseline struct {
	Name            string
	Args            []string
	Rows            [][]string
	History         string
	Exit            int
	NoRepository    bool
	Help            bool
	MissingArgument bool
	UnknownArgument string
}

func readHistoryBaseline(t *testing.T) []historyBaseline {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "selected-defaults.json"))
	if err != nil {
		t.Fatal(err)
	}
	var baseline []historyBaseline
	if err := json.Unmarshal(data, &baseline); err != nil {
		t.Fatal(err)
	}
	return baseline
}

// want renders the expected output of one baseline case through its owner.
func (f historyFixtureData) want(t *testing.T, c historyBaseline) string {
	t.Helper()
	switch {
	case c.NoRepository:
		return toon.NotInRepo() + "\n"
	case c.Help:
		return historyUsage + "\n"
	case c.MissingArgument:
		return toon.MissingArg(historyCmd, "argument") + "\n"
	case c.UnknownArgument != "":
		return toon.Usage(historyCmd, c.UnknownArgument) + "\n"
	}
	name := c.Name
	if c.History != "" {
		name = c.History
	}
	want, err := toon.Table("history", historyFields, f.capturedRows(t, name))
	if err != nil {
		t.Fatal(err)
	}
	return want
}

func historyRows(t *testing.T, output, table string) []map[string]any {
	t.Helper()
	document, err := axitest.DecodeDocument(output)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := document.Rows(table)
	if err != nil {
		t.Fatal(err)
	}
	rows := make([]map[string]any, len(raw))
	for i, row := range raw {
		rows[i] = row.(map[string]any)
	}
	return rows
}

// selectedHistory runs the selected view with one --spec flag for each target.
func selectedHistory(t *testing.T, limit string, targets ...string) (string, int) {
	t.Helper()
	args := []string{"history", "--limit", limit}
	for _, target := range targets {
		args = append(args, "--spec", target)
	}
	out, code, _ := Command(args)
	return out, code
}

func TestSelectedHistoriesPreserveDefault(t *testing.T) {
	f := historyFixture(t)
	for _, c := range f.baseline {
		t.Run(c.Name, func(t *testing.T) {
			if c.NoRepository {
				t.Chdir(t.TempDir())
			}
			args := []string{"history"}
			for _, arg := range c.Args {
				args = append(args, f.expand(arg))
			}
			want := f.want(t, c)
			out, code, _ := Command(args)
			if code != c.Exit || out != want {
				t.Fatalf("positional baseline: exit=%d output=%q; want exit=%d output=%q", code, out, c.Exit, want)
			}
		})
	}
}

func TestSelectedSpecPartialFailure(t *testing.T) {
	historyFixture(t)
	realGit, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	// The wrapper fails only the retire query of the slug `failed` and passes every other
	// call to the real executable.
	bin := t.TempDir()
	wrapper := "#!/bin/sh\nfor arg do\n  [ \"$arg\" = '--grep=spec-retire: failed' ] && exit 17\ndone\nexec " + sanitize.ShellQuote(realGit) + " \"$@\"\n"
	if err := os.WriteFile(filepath.Join(bin, "git"), []byte(wrapper), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	out, code := selectedHistory(t, "1", "mixed", "failed", "absent", "specs/failed/spec.md", "retire-only")
	if code != 1 {
		t.Fatalf("partial failure exit=%d output=%q", code, out)
	}
	rows := historyRows(t, out, "histories")
	if len(rows) != 4 {
		t.Fatalf("partial rows=%#v", rows)
	}
	for i, target := range []string{"mixed", "failed", "absent", "retire-only"} {
		row := rows[i]
		if row["target"] != target {
			t.Fatalf("request order=%#v", rows)
		}
		if target == "failed" {
			if row["error"] != historyDerivationFailed || row["detail"] != historyDetail(target, SlugOf(target)) || row["total_events"] != nil || row["total_bytes"] != nil || row["omitted_events"] != nil {
				t.Fatalf("unknown failed counts=%#v", row)
			}
		} else if row["error"] != "" {
			t.Fatalf("erased successful target=%#v", row)
		}
	}
	if fmt.Sprint(rows[2]["total_events"]) != "0" || len(historyRows(t, out, "history")) != 2 {
		t.Fatalf("lost empty or successful events: %s", out)
	}
}

func TestSelectedHistoryGrammar(t *testing.T) {
	t.Chdir(t.TempDir())
	help := selectedHistoryGrammar.Help
	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"--spec", "mixed"}, help},
		{[]string{"--limit", "1"}, help},
		{[]string{"--spec"}, toon.MissingArg(historyCmd, "--spec")},
		{[]string{"--spec", ""}, toon.Usage(historyCmd, usage.EmptyFlagValue("--spec"))},
		{[]string{"--spec", "mixed", "--limit"}, toon.MissingArg(historyCmd, "--limit")},
		{[]string{"--spec", "mixed", "--limit", ""}, toon.Usage(historyCmd, usage.EmptyFlagValue("--limit"))},
		{[]string{"--spec", "mixed", "--limit", "0"}, help},
		{[]string{"--spec", "mixed", "--limit", "-1"}, help},
		{[]string{"--spec", "mixed", "--limit", "NaN"}, help},
		{[]string{"--spec", "mixed", "--limit", "1.5"}, help},
		{[]string{"--spec", "mixed", "--limit", "99999999999999999999999999"}, help},
		{[]string{"--spec", "mixed", "--limit", "1", "--limit", "2"}, toon.Usage(historyCmd, "--limit")},
		{[]string{"--spec", "mixed", "--limit", "1", "positional"}, toon.Usage(historyCmd, "positional")},
		{[]string{"positional", "--spec", "mixed", "--limit", "1"}, toon.Usage(historyCmd, "positional")},
		{[]string{"--spec", "mixed", "--limit", "1", "--"}, help},
		{[]string{"--spec", "mixed", "--limit", "1", "--unknown"}, toon.Usage(historyCmd, "--unknown")},
	} {
		out, code, _ := Command(append([]string{"history"}, c.args...))
		if code != 2 || out != c.want+"\n" {
			t.Fatalf("grammar %q: exit=%d output=%q; want exit=2 output=%q", c.args, code, out, c.want+"\n")
		}
	}
	for _, spelling := range []string{"--help", "-h"} {
		out, code, _ := Command([]string{"history", "--spec", "mixed", spelling})
		if code != 0 || out != help+"\n" {
			t.Fatalf("selected help exit=%d output=%q; want %q", code, out, help+"\n")
		}
	}
	out, code := selectedHistory(t, "1", "mixed")
	if code != 1 || out != toon.NotInRepo()+"\n" {
		t.Fatalf("selected no-repository exit=%d output=%q", code, out)
	}
}

func TestSelectedHistoryHostileTarget(t *testing.T) {
	f := historyFixture(t)
	for _, bad := range []string{"bad\x1btarget", "tab\ttarget", "line\ntarget", "return\rtarget", "nul\x00target", "del\x7ftarget", "c1\u0085target"} {
		for _, c := range []struct {
			targets []string
			ordinal int
			failed  int
		}{
			{[]string{bad, "mixed", bad, "absent"}, 1, 0},
			{[]string{"mixed", "absent", bad, "retire-only", bad}, 3, 2},
			{[]string{"mixed", "specs/mixed/spec.md", bad}, 3, 1},
		} {
			out, code := selectedHistory(t, "1", c.targets...)
			if code != 1 {
				t.Fatalf("unsafe %q exit=%d output=%q", bad, code, out)
			}
			rows := historyRows(t, out, "histories")
			if len(rows) != len(c.targets)-1 {
				t.Fatalf("unsafe duplicate rows=%#v", rows)
			}
			row := rows[c.failed]
			pointer := sanitize.TargetPointer(c.ordinal)
			if row["target"] != pointer || row["slug"] != "" || row["total_events"] != nil || row["detail"] != "" || row["error"] != selectedTargetControls {
				t.Fatalf("unsafe stable identity=%#v; want %s", row, pointer)
			}
			for i, row := range rows {
				if i != c.failed && row["error"] != "" {
					t.Fatalf("erased other history=%#v", row)
				}
			}
			events := historyRows(t, out, "history")
			if len(events) < 1 || events[0]["slug"] != "mixed" || events[0]["hash"] != f.hash["newest"] {
				t.Fatalf("lost successful event: %s", out)
			}
		}
	}
}
