package surface

import (
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

// TestDoctorBindingContracts pins the two binding rows of bench doctor: the hard-cut
// migration report over the retired BENCH_TIER_*/BENCH_ALIAS_* schema, and the naming of a
// known harness whose column binds nothing. Both are graded as black-box CLI results in a
// fixture repo — the command runs for real, and the assertions read its output and the
// binding file's bytes rather than any internal state.
func TestDoctorBindingContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "bench doctor retired-key rewrite contract", testDoctorRetiredKeyRewrites)
	contract.RunParallel(t, "bench doctor migration re-run contract", testDoctorMigrationRerun)
	contract.RunParallel(t, "bench doctor unbound column contract", testDoctorUnboundColumn)
}

// retiredLinesEnv is a binding written entirely in the retired schema — the exact state a
// maintainer migrating from it starts in.
const retiredLinesEnv = "# binding\n" +
	"BENCH_TIER_TOP=gpt-5.6-sol\n" +
	"BENCH_TIER_MID=gpt-5.6-terra\n" +
	"BENCH_TIER_CHEAP=gpt-5.6-luna\n" +
	"BENCH_ALIAS_TOP=fable\n" +
	"BENCH_ALIAS_MID=opus\n" +
	"BENCH_ALIAS_CHEAP=sonnet\n"

// wantRetiredRewrites is authored independently of the production mapping, one line per
// retired key. That independence is what makes the omission the story names — reporting
// five of six keys and leaving a silent survivor — turn this test red, so it is not the
// duplicated implementation knowledge the code standard rejects.
var wantRetiredRewrites = []string{
	"BENCH_TIER_TOP=gpt-5.6-sol  ->  BENCH_CODEX_TOP=gpt-5.6-sol",
	"BENCH_TIER_MID=gpt-5.6-terra  ->  BENCH_CODEX_MID=gpt-5.6-terra",
	"BENCH_TIER_CHEAP=gpt-5.6-luna  ->  BENCH_CODEX_CHEAP=gpt-5.6-luna",
	"BENCH_ALIAS_TOP=fable  ->  BENCH_CLAUDE_TOP=fable",
	"BENCH_ALIAS_MID=opus  ->  BENCH_CLAUDE_MID=opus",
	"BENCH_ALIAS_CHEAP=sonnet  ->  BENCH_CLAUDE_CHEAP=sonnet",
}

// testDoctorRetiredKeyRewrites drives the migration row from an otherwise green tree, so
// the red row and the exit code answer for the binding alone. The byte comparison is the
// half that separates report-and-offer from silent mutation: a doctor that helpfully
// rewrote the file would satisfy every text assertion above it.
func testDoctorRetiredKeyRewrites(t *testing.T) {
	f := contract.NewFixture(t)
	sb := newDoctorSandbox(t, f)
	writeDoctorGreenTree(t, f)
	f.BenchWrapperEnv(sb.env, "doctor", "--fix").RequireExit(0)
	f.WriteFile(".bench/lines.env", retiredLinesEnv)

	probe := f.BenchWrapperEnv(sb.env, "doctor")
	probe.RequireExit(1)
	probe.RequireContains(probe.Stdout, "red:")
	for _, want := range wantRetiredRewrites {
		probe.RequireContains(probe.Stdout, want)
	}
	if got := f.ReadFile(".bench/lines.env"); got != retiredLinesEnv {
		t.Fatalf("bench doctor rewrote the binding it only reports on:\nwant %q\ngot  %q", retiredLinesEnv, got)
	}
}

// testDoctorMigrationRerun proves the report is safe to leave in a routed repo: running it
// twice reports the same rewrites and still changes nothing.
func testDoctorMigrationRerun(t *testing.T) {
	f := contract.NewFixture(t)
	sb := newDoctorSandbox(t, f)
	writeDoctorGreenTree(t, f)
	f.BenchWrapperEnv(sb.env, "doctor", "--fix").RequireExit(0)
	f.WriteFile(".bench/lines.env", retiredLinesEnv)

	first := f.BenchWrapperEnv(sb.env, "doctor")
	second := f.BenchWrapperEnv(sb.env, "doctor")
	second.RequireExit(1)
	for _, want := range wantRetiredRewrites {
		first.RequireContains(first.Stdout, want)
		second.RequireContains(second.Stdout, want)
	}
	if got := f.ReadFile(".bench/lines.env"); got != retiredLinesEnv {
		t.Fatalf("a second bench doctor run rewrote the binding:\nwant %q\ngot  %q", retiredLinesEnv, got)
	}
}

// testDoctorUnboundColumn pins story 3's doctor row against a binding whose codex and
// claude columns are complete: the unadopted harness is named along with the keys that
// would bind it, and naming it is not itself a defect — an unbound column fails closed by
// design, so the row must not turn doctor red.
func testDoctorUnboundColumn(t *testing.T) {
	f := contract.NewFixture(t)
	sb := newDoctorSandbox(t, f)
	writeDoctorGreenTree(t, f)
	f.BenchWrapperEnv(sb.env, "doctor", "--fix").RequireExit(0)
	f.WriteFile(".bench/lines.env", "BENCH_CODEX_TOP=gpt-5.6-sol\nBENCH_CODEX_MID=gpt-5.6-terra\n"+
		"BENCH_CODEX_CHEAP=gpt-5.6-luna\nBENCH_CLAUDE_TOP=fable\nBENCH_CLAUDE_MID=opus\nBENCH_CLAUDE_CHEAP=sonnet\n")

	probe := f.BenchWrapperEnv(sb.env, "doctor")
	probe.RequireExit(0)
	probe.RequireContains(probe.Stdout, "opencode")
	for _, key := range []string{"BENCH_OPENCODE_TOP", "BENCH_OPENCODE_MID", "BENCH_OPENCODE_CHEAP"} {
		probe.RequireContains(probe.Stdout, key)
	}
}
