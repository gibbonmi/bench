package probe

import (
	"os"
	"strconv"
	"time"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/sanitize"
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

// selectionFields is the schema of the row that names what the probe ran. ran is a genuine
// count, so the row renders through toon.TableTyped and that cell stays an integer.
var selectionFields = []string{"form", "target", "run", "baseline", "ran"}

// verdictCells is one probe answer's variable half: the four cells the verdict row takes
// from a run, and the two cells the selection row takes from the baseline and the mutated
// run. The six travel together, because they come from two different runs and one outcome
// value cannot carry both.
type verdictCells struct {
	verdict  string
	cause    string
	failed   int
	restored string
	baseline string
	ran      int
}

// run applies the refusal order, grades the baseline, and then runs the probe itself. Every
// refusal below answers before preserve writes anything, which is what makes a refused probe
// leave the tree, the Bench home, and the process table exactly as it found them. The
// baseline is the last of them, because it is the only one that starts a run child.
func run(root string, parsed usage.Result) (string, int) {
	subject, line := resolveSubject(root, parsed.Positionals[0])
	if line != "" {
		return line, 1
	}
	old, replacement, kind, line := mutationInput(parsed)
	if line != "" {
		return line, 1
	}
	mutationResult, line := mutate(subject.start, old, replacement, kind)
	if line != "" {
		return line, 1
	}
	if mutationResult.matches == 0 {
		return renderInvalidMutation(subject, kind, parsed)
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
	mutation := kind
	line, code, baseline := gradeBaseline(root, subject, mutation, request)
	if line != "" {
		return line, code
	}
	return probe(root, subject, mutationResult.mutated, mutation, request, baseline)
}

func renderInvalidMutation(subject subject, mutation string, parsed usage.Result) (string, int) {
	cells := verdictCells{verdict: "invalid", cause: "substring-miss", restored: untouchedCell}
	form, target, run := "package", parsed.Flags["--package"], testreport.AllTests
	if check, ok := parsed.Flags["--check"]; ok {
		form, target = "check", check
	} else if pattern, ok := parsed.Flags["--run"]; ok {
		run = pattern
	}
	out, err := rowsWithSelection(subject, mutation, cells, form, target, run)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return out, 1
}

func mutationForm(parsed usage.Result) (string, string, string) {
	if old, omit := parsed.Flags["--omit"]; omit {
		return old, "", "omit"
	}
	if call, unwrap := parsed.Flags["--unwrap"]; unwrap {
		return call, "", "unwrap"
	}
	return parsed.Flags["--swap"], parsed.Flags["--with"], "swap"
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
func probe(root string, subject subject, mutated []byte, mutation string, request testreport.Request, baseline testreport.OutcomeKind) (string, int) {
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
	return render(subject, mutation, outcome, request, report, preserved, restored, reason, baseline)
}

// render prints the verdict row first, so a caller reads the answer before the evidence.
// A subject the encoder cannot carry answers with the shared render error in place of the
// row and the report, and the restore has already run by then, which is why an unprintable
// name still leaves a clean tree. The restore-failed verdict overrides every other answer,
// so a failed restore keeps the preserved row and exit 2 even when the row itself refuses.
func render(subject subject, mutation string, outcome testreport.Outcome, request testreport.Request, report string, preserved preservation, restored bool, reason string, baseline testreport.OutcomeKind) (string, int) {
	// The baseline cell carries the kind the baseline run observed, so the row joins the two
	// runs it reports rather than asserting the kind the refusal order implies.
	cells := verdictCells{
		cause:    string(outcome.Kind),
		failed:   outcome.FailedTests,
		restored: "yes",
		baseline: string(baseline),
		ran:      outcome.Ran,
	}
	var code int
	cells.verdict, code = verdictFor(outcome.Kind)
	if !restored {
		cells.verdict, code, cells.restored = "restore-failed", 2, "no"
	}
	out, err := rows(subject, mutation, cells, request)
	if err != nil {
		out, report = toon.RenderError(err)+"\n", ""
		if restored {
			code = 1
		}
	}
	if !restored {
		// The reason is diagnostic prose, and a failed write quotes the subject path inside
		// it, so a subject the row already refused would refuse here too and the caller would
		// lose the copy path. The path cell stays verbatim, because a stripped path names a
		// file that is not there.
		preservedRow := [][]string{{preserved.file, sanitize.Strip(reason)}}
		block, err := toon.Table("preserved", []string{"path", "reason"}, preservedRow)
		if err != nil {
			return toon.RenderError(err) + "\n", code
		}
		out += block
	}
	return out + report, code
}

// rows renders the verdict row and the selection row that follows it. The pair renders in
// one call, so a cell the encoder cannot carry refuses both rather than printing half an
// answer. The selection row names what the focused run selected, what the baseline reached,
// and how many tests the mutated run ran, which is what separates a mutation no test
// observed from a run that started none.
func rows(subject subject, mutation string, cells verdictCells, request testreport.Request) (string, error) {
	return rowsWithSelection(subject, mutation, cells, request.Form(), request.Target(), request.Run())
}

func rowsWithSelection(subject subject, mutation string, cells verdictCells, form, target, run string) (string, error) {
	verdictRow := []any{cells.verdict, subject.display, mutation, cells.cause, cells.failed, cells.restored}
	out, err := toon.TableTyped("probe", probeFields, [][]any{verdictRow})
	if err != nil {
		return "", err
	}
	selectionRow := []any{form, target, run, cells.baseline, cells.ran}
	block, err := toon.TableTyped("selection", selectionFields, [][]any{selectionRow})
	if err != nil {
		return "", err
	}
	return out + block, nil
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
