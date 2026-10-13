package probe

import (
	"testing"

	"github.com/gibbonmi/bench/internal/probe/probetest"
)

// unorderedFailures fails three tests in an event order that is neither the package order
// nor the test-name order: the first event, the last event, and the lowest test name each
// name a test that is not the first failure.
const unorderedFailures = `{"Action":"run","Package":"probefixture","Test":"TestClampPositive"}
{"Action":"output","Package":"probefixture","Test":"TestClampPositive","Output":"    caught\n"}
{"Action":"fail","Package":"probefixture","Test":"TestClampPositive","Elapsed":0}
{"Action":"run","Package":"probefixture","Test":"TestClampNegative"}
{"Action":"output","Package":"probefixture","Test":"TestClampNegative","Output":"    caught\n"}
{"Action":"fail","Package":"probefixture","Test":"TestClampNegative","Elapsed":0}
{"Action":"fail","Package":"probefixture","Elapsed":0.01}
{"Action":"run","Package":"probefixture/edge","Test":"TestAlpha"}
{"Action":"output","Package":"probefixture/edge","Test":"TestAlpha","Output":"    caught\n"}
{"Action":"fail","Package":"probefixture/edge","Test":"TestAlpha","Elapsed":0}
{"Action":"fail","Package":"probefixture/edge","Elapsed":0.01}`

// subtestFailure fails only the subtest TestClampNegative/zero. The parent prints no
// diagnostic of its own, so the failures table names only the subtest.
const subtestFailure = `{"Action":"run","Package":"probefixture","Test":"TestClampNegative"}
{"Action":"run","Package":"probefixture","Test":"TestClampNegative/zero"}
{"Action":"output","Package":"probefixture","Test":"TestClampNegative/zero","Output":"    caught\n"}
{"Action":"fail","Package":"probefixture","Test":"TestClampNegative/zero","Elapsed":0}
{"Action":"fail","Package":"probefixture","Test":"TestClampNegative","Elapsed":0}
{"Action":"fail","Package":"probefixture","Elapsed":0.01}`

// PS40, PS41, PS42, PS43: --expect gives bit only over a failure of the named test or of
// one of its subtests, another failure reads invalid with expect-missed, and the prose
// check, which names no test, refuses the flag as usage.
func TestExpectRequiresTheNamedWitness(t *testing.T) {
	for _, tc := range []struct {
		name, events, expect string
		code                 int
		verdict, cause       string
		first                string
	}{
		{"the named test", cannedFailure("caught"), "TestClampNegative", 0, "bit", "failed", "TestClampNegative"},
		{"a subtest of the named test", subtestFailure, "TestClampNegative", 0, "bit", "failed", "TestClampNegative/zero"},
		{"another test", cannedFailure("caught"), "TestClampPositive", 1, "invalid", "expect-missed", "TestClampNegative"},
		{"a name prefix", cannedFailure("caught"), "TestClampNeg", 1, "invalid", "expect-missed", "TestClampNegative"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newFixture(t)
			installStubGo(t, f, tc.events, 1, "")
			out, code := runProbe(t, append(probeArgs(), "--expect", tc.expect)...)
			requireRow(t, out, code, tc.code, probeRow(t, tc.verdict, "clamp.go", "swap", tc.cause, 1, tc.first, "yes"))
			requireSubjectBytes(t, f, clampSource)
		})
	}
	t.Run("prose check", func(t *testing.T) {
		f := refusalFixture(t)
		out, code := runProbe(t, "clamp.go", "--omit", "n < 0", "--check", "prose", "--expect", "TestClampNegative")
		requireRefused(t, f, out, code, grammar.Help+"\n", 2)
	})
}

// PS38, PS39: first_failure names the first failed test in package-then-test order, and
// it reads null for every run that named no failed test.
func TestFirstFailureNamesTheFirstFailedTest(t *testing.T) {
	t.Run("two failed tests", func(t *testing.T) {
		f := newFixture(t)
		installStubGo(t, f, unorderedFailures, 1, "")
		out, code := runProbe(t, probeArgs()...)
		requireRow(t, out, code, 0, probeRow(t, "bit", "clamp.go", "swap", "failed", 3, "TestClampNegative", "yes"))
	})
	unpatterned := probeArgs("clamp.go", "--swap", "n < 0", "--with", "n > 0", "--package", "./")
	for _, tc := range []struct {
		name                     string
		baseline, mutated        probetest.Start
		args                     []string
		verdict, cause, restored string
	}{
		{"silent run", probetest.Start{Events: cannedPass}, probetest.Start{Events: cannedPass}, probeArgs(), "silent", "passed", "yes"},
		{"build failure", probetest.Start{Events: cannedPass}, probetest.Start{Events: buildFailedRun, Exit: 1}, unpatterned, "invalid", "build-failed", "yes"},
		{"baseline refusal", probetest.Start{Exit: 1}, probetest.Start{Events: cannedPass}, probeArgs(), "invalid", "baseline-refused", "untouched"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newFixture(t)
			installStubStarts(t, f, tc.baseline, tc.mutated)
			out, code := runProbe(t, tc.args...)
			requireRow(t, out, code, 1, probeRow(t, tc.verdict, "clamp.go", "swap", tc.cause, 0, nil, tc.restored))
		})
	}
	t.Run("prose bite", func(t *testing.T) {
		f := newFixture(t)
		findings := trackProseNote(t, f)
		out, code := runProbe(t, proseArgs...)
		requireRow(t, out, code, 0, probeRow(t, "bit", "note.md", "swap", "failed", findings, nil, "yes"))
	})
}
