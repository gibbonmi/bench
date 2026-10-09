package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/probe/probetest"
	"github.com/gibbonmi/bench/internal/responsebound/responseboundtest"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/treetarget/kittest"
	"github.com/gibbonmi/bench/internal/usage"
)

// TestNativeRecordCLIJourney drives the native record forms through the production
// dispatcher in a linked worktree, so the primary-checkout guard stays active.
func TestNativeRecordCLIJourney(t *testing.T) {
	f, primary := recordtest.NewLinked(t, 2, recordtest.Delegate, recordtest.Undispatched)
	slug, err := rr.Slug(recordtest.Spec)
	if err != nil {
		t.Fatal(err)
	}
	spec := filepath.Join(f.Root, filepath.FromSlash(recordtest.Spec))
	read := func() []byte {
		data, err := os.ReadFile(spec)
		if err != nil {
			t.Fatal(err)
		}
		return data
	}

	// TestHelpInventoryIsComplete pins the complete grammar in root help. The form help
	// and the family help must advertise that same grammar.
	form := runAXICommandAt(t, f.Root, []string{"record", "assignment", "--help"})
	grammar, ok := strings.CutPrefix(strings.TrimSuffix(form.stdout, "\n"), "usage: bench record assignment ")
	root := runAXICommandAt(t, f.Root, []string{"help"})
	if form.code != 0 || !ok || !strings.Contains(root.stdout, "\n  bench record [--in <label|primary>] assignment "+grammar+"  ") {
		t.Fatalf("record assignment --help = exit %d, stdout %q; want exit 0 and the grammar that root help advertises", form.code, form.stdout)
	}
	if family := runAXICommandAt(t, f.Root, []string{"record", "--help"}); family.code != 0 || !slices.Contains(strings.Split(family.stdout, "\n"), "usage: bench record assignment "+grammar) {
		t.Fatalf("record --help = exit %d, stdout %q; want exit 0 with the assignment form", family.code, family.stdout)
	}

	execution := f.Plan.Execution
	dispatch := func(ticket string) []string {
		return []string{"record", "assignment", slug, "--source", "HEAD", "--ticket", ticket, "--performer", "native-author-1", "--assignment", "native-dispatch-1",
			"--model", "opus", "--effort", "high", "--ref", "native:dispatch-1", "--run", execution.RunID, "--orchestrator", execution.OrchestratorSession,
			"--author-limit", strconv.Itoa(execution.AuthorLimit)}
	}
	if result := runAXICommandAt(t, primary, dispatch("1.md")); result.code != 1 || !strings.Contains(result.stdout+result.stderr, usage.PrimaryCheckoutRefusal()) {
		t.Fatalf("primary checkout = exit %d, stdout %q, stderr %q; want exit 1 and the primary-checkout refusal", result.code, result.stdout, result.stderr)
	}
	source := f.Tip()
	if result := runAXICommandAt(t, f.Root, dispatch("1.md")); result.code != 0 {
		t.Fatalf("first dispatch = exit %d, stdout %q, stderr %q; want exit 0", result.code, result.stdout, result.stderr)
	}
	f.Commit("dispatch 1.md")
	f.Reload()
	want := []rr.Assignment{{Session: "native-author-1", Assignment: "native-dispatch-1", Model: "opus", Effort: "high", Source: source, NativeRef: "native:dispatch-1"}}
	if got := f.Plan.Execution.Assignments; !reflect.DeepEqual(got["1.md"], want) || len(got["2.md"]) != 0 {
		t.Fatalf("assignments = %+v, want 1.md %+v and 2.md empty", got, want)
	}

	// The same session cannot also author 2.md, and the full-plan validator refuses it
	// before the spec changes.
	before := read()
	if result := runAXICommandAt(t, f.Root, dispatch("2.md")); result.code != 1 || !strings.Contains(result.stdout, "already authors ticket 1.md") {
		t.Fatalf("cross-ticket dispatch = exit %d, stdout %q; want exit 1 naming the session's ticket", result.code, result.stdout)
	}
	if !bytes.Equal(read(), before) {
		t.Fatal("cross-ticket refusal changed the spec")
	}

	// Chunk 1 records against the plan that holds its dispatch. The 2.md dispatch then
	// commits a changed plan, so the chunk 2 entry refuses before any record byte changes,
	// and the amendment admits it.
	path, err := rr.RecordPath(recordtest.Spec)
	if err != nil {
		t.Fatal(err)
	}
	record := filepath.Join(f.Root, filepath.FromSlash(path))
	chunk := func(id string) (axiCommandResult, []byte) {
		base := f.Tip()
		f.Write("source.txt", "chunk "+id+"\n")
		f.Commit("chunk " + id)
		result := runAXICommandAt(t, f.Root, []string{"record", "chunk", slug, "--chunk", id, "--base", base, "--tip", f.Tip()})
		data, _ := os.ReadFile(record)
		return result, data
	}
	if result, _ := chunk("1"); result.code != 0 {
		t.Fatalf("chunk 1 = exit %d, stdout %q; want exit 0", result.code, result.stdout)
	}
	f.Commit("record chunk 1")
	second := dispatch("2.md")
	second[slices.Index(second, "--performer")+1] = "native-author-2"
	if result := runAXICommandAt(t, f.Root, second); result.code != 0 {
		t.Fatalf("second dispatch = exit %d, stdout %q; want exit 0", result.code, result.stdout)
	}
	f.Commit("dispatch 2.md")
	recorded, err := os.ReadFile(record)
	if err != nil {
		t.Fatal(err)
	}
	result, after := chunk("2")
	if result.code != 1 || !strings.Contains(result.stdout, "stale plan digest") || !strings.Contains(result.stdout, "bench record amendment") || !bytes.Equal(after, recorded) {
		t.Fatalf("stale chunk 2 = exit %d, stdout %q; want exit 1 naming the amendment form and unchanged record bytes", result.code, result.stdout)
	}
	if result := runAXICommandAt(t, f.Root, []string{"record", "amendment", slug, "--source", "HEAD"}); result.code != 0 {
		t.Fatalf("amendment = exit %d, stdout %q; want exit 0", result.code, result.stdout)
	}
	if result := runAXICommandAt(t, f.Root, []string{"record", "chunk", slug, "--chunk", "2", "--base", f.Git("rev-parse", "HEAD~1"), "--tip", f.Tip()}); result.code != 0 {
		t.Fatalf("amended chunk 2 = exit %d, stdout %q; want exit 0", result.code, result.stdout)
	}
	probeImportJourney(t, f, slug, record)
}

// The journey's Go streams and exits are authored here, apart from the probe tests. The
// baseline passes. A biting mutated run fails one test and exits 1, so the probe command
// exits 0 while its focused execution exits 1; a silent mutated run passes, so the command
// exits 1 while its focused execution exits 0.
const (
	journeyPass = `{"Action":"run","Package":"benchkit/clamp","Test":"TestClampNegative"}
{"Action":"pass","Package":"benchkit/clamp","Test":"TestClampNegative","Elapsed":0}
{"Action":"pass","Package":"benchkit/clamp","Elapsed":0.01}`
	journeyFailFormat = `{"Action":"run","Package":"benchkit/clamp","Test":"TestClampNegative"}
{"Action":"output","Package":"benchkit/clamp","Test":"TestClampNegative","Output":"    %s\n"}
{"Action":"fail","Package":"benchkit/clamp","Test":"TestClampNegative","Elapsed":0}
{"Action":"fail","Package":"benchkit/clamp","Elapsed":0.01}`
	journeySubject = "package clamp\n\n// Clamp keeps n at or above zero.\nfunc Clamp(n int) int {\n\tif n < 0 {\n\t\treturn 0\n\t}\n\treturn n\n}\n"
)

// probeImportJourney runs the real probe producer through the production dispatcher in the
// fixture's linked Go root, then imports its selected bytes through the real record entry.
// Each incomplete selection refuses with the record bytes unchanged.
func probeImportJourney(t *testing.T, f *recordtest.Fixture, slug, record string) {
	root, err := filepath.EvalSymlinks(f.Root)
	if err != nil {
		t.Fatal(err)
	}
	kittest.WriteTree(t, root)
	f.Write("clamp/clamp.go", journeySubject)
	f.Write("clamp/clamp_test.go", "package clamp\n\nimport \"testing\"\n\nfunc TestClampNegative(t *testing.T) {\n\tif Clamp(-1) != 0 {\n\t\tt.Fatal(\"negative\")\n\t}\n}\n")
	staged, executable := filepath.Join(t.TempDir(), "staged"), filepath.Join(t.TempDir(), "bench")
	if err := os.WriteFile(staged, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := freshness.Publish(root, staged, executable, filepath.Dir(executable), "native-journey"); err != nil {
		t.Fatal(err)
	}
	t.Setenv(runbinary.Env, executable)
	t.Setenv("BENCH_KIT", root)
	home := t.TempDir()
	t.Setenv(benchhome.Env, home)
	t.Chdir(root)
	marker := filepath.Join(t.TempDir(), "go-started")
	produce := func(mutated probetest.Start, full bool) (string, int) {
		probetest.InstallStarts(t, marker, probetest.Start{Events: journeyPass}, mutated)
		argv := []string{"probe", "clamp/clamp.go", "--swap", "n < 0", "--with", "n > 0", "--package", "./clamp", "--run", "^TestClampNegative$"}
		if full {
			argv = append(argv, "--full")
		}
		var stdout, stderr bytes.Buffer
		code := Command{Stdout: &stdout, Stderr: &stderr, Executable: "bench"}.Run(argv)
		if got, _ := os.ReadFile(filepath.Join(root, "clamp", "clamp.go")); string(got) != journeySubject {
			t.Fatalf("probe left the subject %q", got)
		}
		return stdout.String(), code
	}
	selected := func(name, data string) string {
		path := filepath.Join(t.TempDir(), name)
		if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
			t.Fatal(err)
		}
		return path
	}
	id := 0
	importing := func(path string, exit int) axiCommandResult {
		id++
		return runAXICommandAt(t, root, []string{"record", "verification", slug, "--chunk", "2", "--requirement", "tests", "--id", "probe-" + strconv.Itoa(id),
			"--performer", "native-author-2", "--model", "opus", "--effort", "high", "--exit-code", strconv.Itoa(exit), "--ref", "native:probe", "--probe-result", path})
	}
	imported := func(path string, command int, want rr.Probe) {
		t.Helper()
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if result := importing(path, command); result.code != 0 {
			t.Fatalf("import %s = exit %d, stdout %q; want exit 0", path, result.code, result.stdout)
		}
		retained, err := rr.Read(root, recordtest.Spec)
		if err != nil {
			t.Fatal(err)
		}
		results := retained.Chunks[1].Verification
		got := results[len(results)-1]
		want.Mutation, want.Restore, want.NativeRef = "omit source check", "pass", got.NativeRef
		if got.Probe == nil || *got.Probe != want || got.NativeRef.Excerpt != string(data) || got.NativeRef.Digest != rr.Digest(data) {
			t.Fatalf("imported probe = %+v with excerpt %q; want %+v and the exact selected bytes", got.Probe, got.NativeRef.Excerpt, want)
		}
	}

	// A silent answer is within the bound, so stdout is complete after its tree lead.
	ordinary, code := produce(probetest.Start{Events: journeyPass}, false)
	if _, spilled := responseboundtest.Find(ordinary); code != 1 || spilled || !strings.HasPrefix(ordinary, "tree[1]{target,head,dirty}:\n") {
		t.Fatalf("silent probe = exit %d, stdout %q; want exit 1 and complete tree-prefixed stdout", code, ordinary)
	}
	imported(selected("ordinary.out", ordinary), 1, rr.Probe{Outcome: "silent", ExitCode: 0})

	// The --full biting answer carries a long diagnostic, so it spills; the body is read only
	// after the production marker names it.
	biting := probetest.Start{Events: fmt.Sprintf(journeyFailFormat, strings.Repeat("d", 2*bounds.ResponseBytes)), Exit: 1}
	projected, code := produce(biting, true)
	spill, spilled := responseboundtest.Find(projected)
	if code != 0 || !spilled || !strings.HasPrefix(spill.Path, home) {
		t.Fatalf("--full probe = exit %d, stdout %q; want exit 0 and a spill under the private home", code, projected)
	}
	imported(spill.Path, 0, rr.Probe{Outcome: "bit", ExitCode: 1})
	body, err := os.ReadFile(spill.Path)
	if err != nil {
		t.Fatal(err)
	}

	// A spill store that cannot open sends the complete output and a spill-failed line.
	blocked := t.TempDir()
	if err := os.WriteFile(filepath.Join(blocked, "responses"), nil, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv(benchhome.Env, blocked)
	failed, code := produce(biting, true)
	if code != 0 || !strings.Contains(failed, "\nspill-failed{") {
		t.Fatalf("blocked --full probe = exit %d, stdout %q; want exit 0 and a spill-failed line", code, failed)
	}
	envelopeEnd := 0
	for range 6 {
		envelopeEnd += bytes.IndexByte(body[envelopeEnd:], '\n') + 1
	}
	lines := strings.SplitAfter(string(body), "\n")
	count := strings.Split(strings.TrimSuffix(lines[5], "\n"), ",")
	recount := func(delta int) string {
		n, _ := strconv.Atoi(count[len(count)-1])
		cells := append(append([]string{}, count[:len(count)-1]...), strconv.Itoa(n+delta))
		return strings.Join(lines[:5], "") + strings.Join(cells, ",") + "\n" + strings.Join(lines[6:], "")
	}
	before, err := os.ReadFile(record)
	if err != nil {
		t.Fatal(err)
	}
	// The projected marker names the complete spill body outside the tree, so a reader that
	// followed it would import a valid answer; each refusal below proves it read only the
	// selected bytes.
	for _, tc := range []struct{ name, data, want string }{
		{"projected stdout", projected, "incomplete response"},
		{"spill-failed stdout", failed, "incomplete response"},
		{"truncated envelope", string(body[:envelopeEnd-3]), "envelope is incomplete"},
		{"truncated suffix", string(body[:len(body)-7]), "trailing_bytes declares"},
		{"count one over", recount(1), "trailing_bytes declares"},
		{"count one under", recount(-1), "trailing_bytes declares"},
	} {
		if result := importing(selected("refused.out", tc.data), 0); result.code != 1 || !strings.Contains(result.stdout, tc.want) {
			t.Fatalf("%s import = exit %d, stdout %q; want exit 1 naming %q", tc.name, result.code, result.stdout, tc.want)
		}
		if after, _ := os.ReadFile(record); !bytes.Equal(after, before) {
			t.Fatalf("%s import changed the record", tc.name)
		}
	}
}
