package probe

import (
	"os"
	"strconv"
	"time"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/testreport"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// proseCheck names the one `bench test` check that grades sentences instead of running a
// Go test. It and the system suite are refused as probe targets, because neither reaches
// the report a verdict derives from.
const proseCheck = "prose"

// probeFields is the verdict row's schema. failed_tests is a genuine count, so the row
// renders through toon.TableTyped and that cell stays an integer on a round-trip.
var probeFields = []string{"verdict", "subject", "mutation", "cause", "failed_tests", "restored"}

// run applies the refusal order and then the probe itself. Every refusal below answers
// before preserve writes anything, which is what makes a refused probe leave the tree,
// the Bench home, and the process table exactly as it found them.
func run(root string, parsed usage.Result) (string, int) {
	subject, line := resolveSubject(root, parsed.Positionals[0])
	if line != "" {
		return line, 1
	}
	old, replacement, omit := mutationForm(parsed)
	mutated, line := mutate(subject.start, old, replacement, omit)
	if line != "" {
		return line, 1
	}
	request, line, code := testreport.Prepare(root, selectionArgs(parsed))
	if line != "" {
		return line, code
	}
	if line := gradeCheckTarget(parsed); line != "" {
		return line, 1
	}
	if line := gradeGateLock(root); line != "" {
		return line, 1
	}
	return probe(root, subject, mutated, mutationName(omit), request)
}

func mutationForm(parsed usage.Result) (string, string, bool) {
	if old, omit := parsed.Flags["--omit"]; omit {
		return old, "", true
	}
	return parsed.Flags["--swap"], parsed.Flags["--with"], false
}

func mutationName(omit bool) string {
	if omit {
		return "omit"
	}
	return "swap"
}

// selectionArgs spells the focused run in `bench test`'s own grammar, so the probe and
// the verb it composes can never disagree about what a selection means.
func selectionArgs(parsed usage.Result) []string {
	var args []string
	if check, ok := parsed.Flags["--check"]; ok {
		args = append(args, "--check", check)
	} else {
		args = append(args, "--package", parsed.Flags["--package"])
		if run, ok := parsed.Flags["--run"]; ok {
			args = append(args, "--run", run)
		}
	}
	if _, ok := parsed.Flags["--full"]; ok {
		args = append(args, "--full")
	}
	return args
}

func gradeCheckTarget(parsed usage.Result) string {
	switch parsed.Flags["--check"] {
	case proseCheck, gate.SystemPhaseName:
		hint := "--check prose and --check system are not probe targets"
		return toon.Errorf("probe focused run unsupported", hint) + "\n"
	}
	return ""
}

// gradeGateLock refuses while a gate run holds the tree. The gate grades the bytes it
// read at its start, so a probe that mutated a file under it would make the gate's
// verdict describe a tree that no longer exists.
func gradeGateLock(root string) string {
	held, err := gate.ExecutionInProgress(root)
	if err != nil {
		return toon.Errorf("gate execution state unavailable", err.Error()) + "\n"
	}
	if held {
		hint := "wait for the gate run to finish before you mutate the tree"
		return toon.Errorf("gate execution in progress", hint) + "\n"
	}
	return ""
}

// probe preserves, mutates, runs, and restores. The restore is deferred around the run,
// so an interrupt or a panic inside the focused run still puts the subject back before
// the verb answers, and the render runs only after the restore has been proven.
func probe(root string, subject subject, mutated []byte, mutation string, request testreport.Request) (string, int) {
	preserved, line := preserve(root, subject)
	if line != "" {
		return line, 1
	}
	if err := replaceAtomic(subject.path, mutated, subject.mode); err != nil {
		preserved.release()
		return toon.Errorf("probe mutation failed", err.Error()) + "\n", 1
	}
	var outcome testreport.Outcome
	var report string
	restored, reason := false, ""
	func() {
		defer func() { restored, reason = restore(subject, preserved) }()
		outcome, report, _ = testreport.Execute(root, request)
	}()
	if restored {
		preserved.release()
	}
	return render(subject, mutation, outcome, report, preserved, restored, reason)
}

// render prints the verdict row first, so a caller reads the answer before the evidence.
// A subject the encoder cannot carry answers with the shared render error, and the
// restore has already run by then, which is why an unprintable name still leaves a clean
// tree.
func render(subject subject, mutation string, outcome testreport.Outcome, report string, preserved preservation, restored bool, reason string) (string, int) {
	verdict, code := verdictFor(outcome.Kind)
	restoredCell := "yes"
	if !restored {
		verdict, code, restoredCell = "restore-failed", 2, "no"
	}
	row := []any{verdict, subject.display, mutation, string(outcome.Kind), outcome.FailedTests, restoredCell}
	out, err := toon.TableTyped("probe", probeFields, [][]any{row})
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	if !restored {
		preservedRow := [][]string{{preserved.file, reason}}
		block, err := toon.Table("preserved", []string{"path", "reason"}, preservedRow)
		if err != nil {
			return toon.RenderError(err) + "\n", 1
		}
		out += block
	}
	return out + report, code
}

// verdictFor maps the focused run's outcome onto the word a caller cites. A failing test is
// the success case, because a mutation that no test noticed proves nothing; every kind
// that reached no verdict at all is invalid rather than silent.
func verdictFor(kind testreport.OutcomeKind) (string, int) {
	switch kind {
	case testreport.OutcomeFailed:
		return "bit", 0
	case testreport.OutcomePassed:
		return "silent", 1
	default:
		return "invalid", 1
	}
}

// stamp names one run's preserved directory. The nanosecond clock and the process id
// together keep two concurrent probes on one repository in separate directories.
func stamp() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10) + "-" + strconv.Itoa(os.Getpid())
}
