package recordcmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordcmd"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/toon"
)

// drop removes a flag from the argv that formArgs builds.
const drop = "\x00drop"

// excerptFile writes data to a new file and returns its path.
func excerptFile(t *testing.T, name, data string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// verifyArgs is the verification argv for requirement tests of chunk 1 with its planned
// probe.
func verifyArgs(t *testing.T, edits ...map[string]string) []string {
	t.Helper()
	return formArgs(t, "verification", map[string]string{"--chunk": "1", "--requirement": "tests", "--id": "tests-1", "--performer": recordtest.Author("1.md"),
		"--model": "unknown", "--effort": "unknown", "--exit-code": "0", "--ref": "fixture:terminal-result", "--excerpt": excerptFile(t, "excerpt.txt", "ok\n"),
		"--probe-outcome": "bit", "--probe-exit-code": "1", "--probe-restore": "pass"}, edits...)
}

// formArgs is the argv of form for the fixture slug with the flags of flags. Each edit
// replaces a flag value, an empty value names a flag without a value, and drop removes
// the flag.
func formArgs(t *testing.T, form string, flags map[string]string, edits ...map[string]string) []string {
	t.Helper()
	for _, edit := range edits {
		for name, value := range edit {
			flags[name] = value
			if value == drop {
				delete(flags, name)
			}
		}
	}
	names := make([]string, 0, len(flags))
	for name := range flags {
		names = append(names, name)
	}
	sort.Strings(names)
	args := []string{form, slug(t)}
	for _, name := range names {
		args = append(args, name)
		if flags[name] != "" {
			args = append(args, flags[name])
		}
	}
	return args
}

// unprobed edits a call to requirement additional, which plans no probe.
var unprobed = map[string]string{"--requirement": "additional", "--id": "additional-1", "--probe-outcome": drop, "--probe-exit-code": drop, "--probe-restore": drop}

// final edits a call to the completion list at source.
func final(source string) map[string]string {
	return map[string]string{"--chunk": drop, "--final": "", "--source": source, "--requirement": "acceptance", "--id": "acceptance-1", "--performer": recordtest.Orchestrator,
		"--probe-outcome": drop, "--probe-exit-code": drop, "--probe-restore": drop}
}

// recorded is a fixture with count planned chunks whose chunk 1 entry is written.
func recorded(t *testing.T, count int) *recordtest.Fixture {
	t.Helper()
	f := linked(t, count)
	base, tip := advance(f, "chunk 1")
	record(t, f, "1", base, tip)
	return f
}

// verify runs the verification form and requires exit 0.
func verify(t *testing.T, f *recordtest.Fixture, edits ...map[string]string) string {
	t.Helper()
	return succeed(t, f, verifyArgs(t, edits...))
}

// succeed runs args over the record of f and requires exit 0.
func succeed(t *testing.T, f *recordtest.Fixture, args []string) string {
	t.Helper()
	out, code := recordcmd.Command(f.Root, args)
	if code != 0 {
		t.Fatalf("%s = exit %d, output %q; want exit 0", args[0], code, out)
	}
	return out
}

// refuseVerify runs the verification form over the record of f and requires exit 1, an
// output that holds want, and unchanged record bytes.
func refuseVerify(t *testing.T, f *recordtest.Fixture, want string, edits ...map[string]string) string {
	t.Helper()
	return refuseArgs(t, f, want, verifyArgs(t, edits...))
}

// refuseArgs runs args over the record of f and requires exit 1, an output that holds
// want, and unchanged record bytes.
func refuseArgs(t *testing.T, f *recordtest.Fixture, want string, args []string) string {
	t.Helper()
	path := recordFile(t, f)
	before := bytesOf(t, path)
	out, code := recordcmd.Command(f.Root, args)
	if code != 1 || !strings.Contains(out, want) {
		t.Fatalf("%s = exit %d, output %q; want exit 1 naming %q", args[0], code, out, want)
	}
	if after := bytesOf(t, path); !bytes.Equal(before, after) {
		t.Fatalf("refusal changed the record")
	}
	return out
}

func chunkResults(t *testing.T, f *recordtest.Fixture) []rr.Verification {
	t.Helper()
	return onlyChunk(t, f).Verification
}

// onlyResult is the one verification result in the chunk 1 list.
func onlyResult(t *testing.T, f *recordtest.Fixture) rr.Verification {
	t.Helper()
	results := chunkResults(t, f)
	if len(results) != 1 {
		t.Fatalf("chunk verification results = %d, want 1", len(results))
	}
	return results[0]
}

// onlyFinal is the one result in the completion list.
func onlyFinal(t *testing.T, f *recordtest.Fixture) rr.Verification {
	t.Helper()
	results := read(t, f).Completion.Verification
	if len(results) != 1 {
		t.Fatalf("completion verification results = %d, want 1", len(results))
	}
	return results[0]
}

func sourceDigest(t *testing.T, f *recordtest.Fixture) string {
	t.Helper()
	digest, err := rr.SourceDigest(f.Root, f.Tree(), recordtest.Spec)
	if err != nil {
		t.Fatal(err)
	}
	return digest
}

func TestRecordVerificationLandsInTheChunkList(t *testing.T) {
	f := recorded(t, 1)
	verify(t, f)
	chunk := onlyChunk(t, f)
	if len(chunk.Verification) != 1 || len(chunk.Reviews) != 0 || chunk.Verification[0].ID != "tests-1" {
		t.Fatalf("chunk = %d verification and %d review results, want the one verification result tests-1", len(chunk.Verification), len(chunk.Reviews))
	}
}

func TestRecordFinalVerificationLandsInTheCompletionList(t *testing.T) {
	f := recorded(t, 1)
	verify(t, f, final(f.Tip()))
	if got := onlyFinal(t, f); got.ID != "acceptance-1" || got.Requirement != "acceptance" || len(chunkResults(t, f)) != 0 {
		t.Fatalf("completion result = %+v, want acceptance-1 and no chunk result", got)
	}
}

func TestRecordVerificationDefaultsToTheChunkTip(t *testing.T) {
	f := recorded(t, 1)
	tip := onlyChunk(t, f).SourceDigest
	advance(f, "later")
	verify(t, f)
	if got := onlyResult(t, f).SourceDigest; got != tip || got == sourceDigest(t, f) {
		t.Fatalf("source digest = %s, want the chunk tip digest %s, not the HEAD digest %s", got, tip, sourceDigest(t, f))
	}
}

func TestRecordVerificationRefusesAStaleSource(t *testing.T) {
	f := recorded(t, 1)
	advance(f, "later")
	refuseVerify(t, f, "bench record chunk", map[string]string{"--source": "HEAD"})
}

func TestRecordFinalVerificationDigestsTheNamedSource(t *testing.T) {
	f := recorded(t, 1)
	tip := onlyChunk(t, f).SourceDigest
	advance(f, "past the last chunk")
	verify(t, f, final("HEAD"))
	if got := onlyFinal(t, f).SourceDigest; got != sourceDigest(t, f) || got == tip {
		t.Fatalf("source digest = %s, want the HEAD digest %s, not the chunk digest %s", got, sourceDigest(t, f), tip)
	}
}

func TestRecordVerificationCopiesThePlannedCommand(t *testing.T) {
	f := recorded(t, 1)
	verify(t, f)
	planned := f.Plan.Chunks[0].Verification[0]
	if got := onlyResult(t, f); planned.ID != "tests" || got.Command != planned.Command || got.Command == "" {
		t.Fatalf("command = %q, want the planned command %q of requirement %s", got.Command, planned.Command, planned.ID)
	}
}

func TestRecordVerificationRefusesAnUnplannedRequirement(t *testing.T) {
	f := recorded(t, 1)
	out := refuseVerify(t, f, "tests", map[string]string{"--requirement": "nosuch"})
	if !strings.Contains(out, "additional") {
		t.Fatalf("output = %q, want the planned IDs tests and additional", out)
	}
}

func TestRecordVerificationWritesTheAuthorRole(t *testing.T) {
	f := recorded(t, 1)
	verify(t, f)
	if role := onlyResult(t, f).Role; role != "author-verification" {
		t.Fatalf("role = %q, want author-verification", role)
	}
}

func TestRecordFinalVerificationWritesTheIntegrationRole(t *testing.T) {
	f := recorded(t, 1)
	verify(t, f, final(f.Tip()))
	if role := onlyFinal(t, f).Role; role != "integration-verification" {
		t.Fatalf("role = %q, want integration-verification", role)
	}
}

func TestRecordVerificationRefusesAnUndispatchedTicket(t *testing.T) {
	f := linked(t, 1)
	f.Plan.Execution.Assignments["1.md"] = []rr.Assignment{}
	f.RewritePlan()
	base, tip := advance(f, "chunk 1")
	record(t, f, "1", base, tip)
	refuseVerify(t, f, "undispatched ticket")
}

func TestRecordVerificationPassesOnExitZero(t *testing.T) {
	f := recorded(t, 1)
	verify(t, f)
	if got := onlyResult(t, f); got.Outcome != "pass" || got.ExitCode == nil || *got.ExitCode != 0 {
		t.Fatalf("result = outcome %q, exit code %v; want pass and 0", got.Outcome, got.ExitCode)
	}
}

func TestRecordVerificationFailsOnNonzeroExit(t *testing.T) {
	f := recorded(t, 1)
	verify(t, f, map[string]string{"--exit-code": "3"})
	if got := onlyResult(t, f); got.Outcome != "fail" || got.ExitCode == nil || *got.ExitCode != 3 {
		t.Fatalf("result = outcome %q, exit code %v; want fail and 3", got.Outcome, got.ExitCode)
	}
}

func TestRecordVerificationWritesThePlannedProbe(t *testing.T) {
	f := recorded(t, 1)
	verify(t, f, map[string]string{"--probe-outcome": "bit", "--probe-exit-code": "4", "--probe-restore": "fail"})
	got := onlyResult(t, f)
	want := rr.Probe{Mutation: f.Plan.Chunks[0].Verification[0].Probe, Outcome: "bit", ExitCode: 4, Restore: "fail", NativeRef: got.NativeRef}
	if got.Probe == nil || want.Mutation == "" || !reflect.DeepEqual(*got.Probe, want) {
		t.Fatalf("probe = %+v, want %+v", got.Probe, want)
	}
}

func TestRecordVerificationRequiresThePlannedProbe(t *testing.T) {
	f := recorded(t, 1)
	refuseVerify(t, f, "--probe-outcome", map[string]string{"--probe-outcome": drop, "--probe-exit-code": drop, "--probe-restore": drop})
}

func TestRecordVerificationRefusesAnUnplannedProbe(t *testing.T) {
	f := recorded(t, 1)
	out := refuseVerify(t, f, "--probe-outcome", unprobed, map[string]string{"--probe-outcome": "bit", "--probe-exit-code": "1", "--probe-restore": "pass"})
	if !strings.Contains(out, "additional") {
		t.Fatalf("output = %q, want requirement additional", out)
	}
}

func TestRecordVerificationNamesTheChunkForm(t *testing.T) {
	f := recorded(t, 2)
	refuseVerify(t, f, "bench record chunk", map[string]string{"--chunk": "2"})
}

func TestRecordFinalVerificationNeedsARecord(t *testing.T) {
	f := linked(t, 1)
	out, code := recordcmd.Command(f.Root, verifyArgs(t, final(f.Tip())))
	if code != 1 || !strings.Contains(out, "bench record chunk") {
		t.Fatalf("final without a record = exit %d, output %q; want exit 1 naming bench record chunk", code, out)
	}
	if _, err := os.Lstat(recordFile(t, f)); !os.IsNotExist(err) {
		t.Fatalf("record file after refusal: %v, want absent", err)
	}
}

func TestRecordVerificationRefusesADuplicateID(t *testing.T) {
	f := recorded(t, 1)
	verify(t, f, unprobed)
	refuseVerify(t, f, "additional-1", map[string]string{"--id": "additional-1"})
}

func TestRecordFinalVerificationRefusesADuplicateID(t *testing.T) {
	f := recorded(t, 1)
	verify(t, f, final(f.Tip()))
	refuseVerify(t, f, "acceptance-1", final(f.Tip()))
}

func TestRecordRefusalLeavesNoTemporaryFile(t *testing.T) {
	f := recorded(t, 1)
	refuseVerify(t, f, "native result", map[string]string{"--ref": "fixture:../outside"})
	entries, err := os.ReadDir(filepath.Dir(recordFile(t, f)))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != filepath.Base(recordFile(t, f)) {
		t.Fatalf("reviews entries = %v, want only the record", entries)
	}
}

// wantVerification is the one output table of a verification form, derived from the
// written entry.
func wantVerification(t *testing.T, list, chunk string, got rr.Verification) string {
	t.Helper()
	table, err := toon.Table("verification", []string{"list", "chunk", "id", "requirement", "role", "outcome", "source_digest", "excerpt_digest"},
		[][]string{{list, chunk, got.ID, got.Requirement, got.Role, got.Outcome, got.SourceDigest, got.NativeRef.Digest}})
	if err != nil {
		t.Fatal(err)
	}
	return table
}

func TestRecordVerificationReportsTheChunkList(t *testing.T) {
	f := recorded(t, 1)
	out := verify(t, f)
	if want := wantVerification(t, "chunk", "1", onlyResult(t, f)); out != want {
		t.Fatalf("output = %q, want %q", out, want)
	}
}

func TestRecordFinalVerificationReportsTheCompletionList(t *testing.T) {
	f := recorded(t, 1)
	out := verify(t, f, final(f.Tip()))
	if want := wantVerification(t, "completion", "", onlyFinal(t, f)); out != want {
		t.Fatalf("output = %q, want %q", out, want)
	}
}

func TestRecordedVerificationPassesTheCheckpoint(t *testing.T) {
	f := recorded(t, 1)
	verify(t, f)
	verify(t, f, unprobed)
	record := read(t, f)
	chunk := &record.Chunks[0]
	for _, axis := range rr.Axes() {
		chunk.Reviews = append(chunk.Reviews, rr.Review{Evidence: f.Evidence("review-"+axis, chunk.SourceDigest, "independent-review"), Axis: axis, Base: chunk.Base, Tip: chunk.Tip})
	}
	path := recordFile(t, f)
	document, err := rr.Render(bytesOf(t, path), record)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, document, 0o644); err != nil {
		t.Fatal(err)
	}
	f.Commit("record chunk 1")
	if err := rr.Check(f.Root, f.Tree(), f.Tip(), recordtest.Spec, "1", false); err != nil {
		t.Fatalf("checkpoint for chunk 1: %v", err)
	}
}

func TestRecordVerificationGrammarRefusals(t *testing.T) {
	f := recorded(t, 1)
	for name, edit := range map[string]map[string]string{
		"both lists":           {"--final": "", "--source": "HEAD"},
		"no list":              {"--chunk": drop},
		"final without source": {"--chunk": drop, "--final": ""},
		"one probe flag":       {"--probe-exit-code": drop, "--probe-restore": drop},
		"two probe flags":      {"--probe-outcome": drop},
		"non-integer exit":     {"--exit-code": "x"},
		"restore outside":      {"--probe-restore": "maybe"},
	} {
		out, code := recordcmd.Command(f.Root, verifyArgs(t, edit))
		if code != 2 || !strings.HasPrefix(out, "usage: bench record verification") {
			t.Errorf("%s = exit %d, output %q; want exit 2 with a usage: bench record verification line", name, code, out)
		}
	}
}

func TestRecordRefusesAControlCharacterInAFlag(t *testing.T) {
	f := recorded(t, 1)
	for _, name := range []string{"--performer", "--model", "--effort", "--id", "--requirement", "--probe-outcome"} {
		out := refuseVerify(t, f, name, map[string]string{name: "value\x1b"})
		if strings.Contains(out, "\x1b") {
			t.Errorf("%s refusal = %q, want no control character", name, out)
		}
	}
}
