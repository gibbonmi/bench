package probe

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/testreport"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// probeFields is the verdict row's schema. failed_tests is a genuine count, so the row
// renders through toon.TableTyped and that cell stays an integer on a round-trip.
var probeFields = []string{"verdict", "subject", "mutation", "cause", "failed_tests", "restored"}

// selectionFields is the schema of the row that names what the probe ran. ran is a genuine
// count, so the row renders through toon.TableTyped and that cell stays an integer.
var selectionFields = []string{"form", "target", "run", "baseline", "ran"}

// executionFields is the schema of the row that follows the probe and selection pair. The
// command exit is the probe's own exit. The focused exit is the normalized exit of the
// mutated run, or null when no mutated run started. trailing_bytes counts every byte after
// the row, so a reader can tell a complete answer from a cut one.
var executionFields = []string{"schema", "command_exit_code", "focused_exit_code", "trailing_bytes"}

// resultTables names the envelope's three one-row tables in their printed order.
var resultTables = []struct {
	name   string
	fields []string
}{{"probe", probeFields}, {"selection", selectionFields}, {"probe_execution", executionFields}}

// restoredCell is the restored cell of a probe whose restore was proven byte-exact.
const restoredCell = "yes"

// resultSchema is the schema cell of the execution row. A reader refuses any other value.
const resultSchema = "bench-probe-result-v1"

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

// run is the probe over the native file operations.
func run(root string, parsed usage.Result) (string, int) {
	return runWith(durablefile.Native(), root, parsed)
}

// runWith applies the refusal order, grades the baseline, and then runs the probe itself,
// with each subject replacement through replacer. Every refusal below answers before
// preserve writes anything, which is what makes a refused probe leave the tree, the Bench
// home, and the process table exactly as it found them. The baseline is the last of them,
// because it is the only one that starts a run child.
func runWith(replacer durablefile.Replacer, root string, parsed usage.Result) (string, int) {
	subject, line := resolveSubject(root, parsed.Positionals[0])
	if line != "" {
		return line, 1
	}
	input, line := mutationInput(parsed)
	if line != "" {
		return line, 1
	}
	mutant, line := mutate(subject.start, input, subject.display)
	if line != "" {
		return line, 1
	}
	request, line, code := testreport.Prepare(root, selectionArgs(parsed))
	if line != "" {
		return line, code
	}
	if mutant.miss {
		return renderInvalidMutation(subject, input.cell, request)
	}
	if line := gradeGateLock(root); line != "" {
		return line, 1
	}
	mutation := input.cell
	line, code, baseline := gradeBaseline(root, subject, mutation, request)
	if line != "" {
		return line, code
	}
	return probe(replacer, root, subject, mutant, mutation, request, baseline)
}

// renderInvalidMutation answers a substring miss. Its selection cells come from the prepared
// request, so the row names the selection that the focused run would have run.
func renderInvalidMutation(subject subject, mutation string, request testreport.Request) (string, int) {
	cells := verdictCells{verdict: "invalid", cause: "substring-miss", restored: untouchedCell}
	e := selected(subject, mutation, cells, request)
	e.command = 1
	out, err := e.answer("")
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return out, 1
}

// selectionArgs spells the focused run in `bench test`'s own grammar, so the probe and
// the verb it composes can never disagree about what a selection means.
func selectionArgs(parsed usage.Result) []string {
	var args []string
	if check, ok := parsed.Flags["--check"]; ok {
		args = append(args, "--check", check)
	} else {
		args = append(args, "--package", parsed.Flags["--package"])
	}
	if run, ok := parsed.Flags["--run"]; ok {
		args = append(args, "--run", run)
	}
	if _, ok := parsed.Flags["--full"]; ok {
		args = append(args, "--full")
	}
	return args
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
func probe(replacer durablefile.Replacer, root string, subject subject, mutant mutationResult, mutation string, request testreport.Request, baseline testreport.OutcomeKind) (string, int) {
	preserved, line := preserve(root, subject)
	if line != "" {
		return line, 1
	}
	if published, err := mutant.write(replacer, subject); err != nil {
		return mutationFailure(replacer, subject, preserved, err, published)
	}
	var outcome testreport.Outcome
	var report string
	var focused int
	restored, reason := false, ""
	func() {
		defer func() { restored, reason = restore(replacer, subject, preserved) }()
		outcome, report, focused = testreport.Execute(root, request)
	}()
	if restored {
		preserved.release()
	}
	return render(subject, mutation, outcome, focused, request, report, preserved, restored, reason, baseline)
}

// write puts the mutant in place of the subject: a removal through replacer for --delete,
// and a replacement otherwise. With an error, published reports that the subject may no
// longer hold its start bytes. A failed removal can leave the subject absent, so it always
// counts as published.
func (r mutationResult) write(replacer durablefile.Replacer, subject subject) (bool, error) {
	if r.remove {
		return true, replacer.Remove(subject.path)
	}
	err := replacer.Replace(subject.path, r.mutated, subject.mode)
	var replaceErr *durablefile.Error
	return errors.As(err, &replaceErr) && replaceErr.Published, err
}

// mutationFailure answers a mutation that returned an error, and no focused run starts.
// An unpublished failure left the start bytes in place, so the copy goes at once. After a
// published failure the mutant may be visible, so the copy stays until a restore proves
// the start bytes back, and a failed restore keeps it at exit 2.
func mutationFailure(replacer durablefile.Replacer, subject subject, preserved preservation, err error, published bool) (string, int) {
	line := toon.Errorf("probe mutation failed", err.Error()) + "\n"
	if published {
		if restored, reason := restore(replacer, subject, preserved); !restored {
			return line + keptCopy(preserved, err.Error()+"; restore: "+reason), 2
		}
	}
	preserved.release()
	return line, 1
}

// keptCopy names the kept copy and the reason it is kept, because the caller restores the
// subject from that copy by hand. The reason is diagnostic prose that can quote the subject
// path, so its control bytes are stripped. A stripped path names a file that is not there,
// so the path stays verbatim, and a path the encoder refuses is Go-quoted outside the table
// beside that refusal, so the refusal never hides the copy.
func keptCopy(preserved preservation, reason string) string {
	reason = sanitize.Strip(reason)
	block, err := toon.Table("preserved", []string{"path", "reason"}, [][]string{{preserved.file, reason}})
	if err == nil {
		return block
	}
	return toon.RenderError(err) + "\n" + toon.Errorf("probe copy kept", strconv.Quote(preserved.file)+": "+reason) + "\n"
}

// render prints the verdict row first, so a caller reads the answer before the evidence.
// focused is the normalized exit of the mutated run. A subject the encoder cannot carry
// answers with the shared render error in place of the rows and the report, and the
// restore has already run by then, which is why an unprintable name still leaves a clean
// tree. The restore-failed verdict overrides every other answer, so a failed restore keeps
// the answer that names the kept copy and exit 2 even when the rows themselves refuse.
func render(subject subject, mutation string, outcome testreport.Outcome, focused int, request testreport.Request, report string, preserved preservation, restored bool, reason string, baseline testreport.OutcomeKind) (string, int) {
	// The baseline cell carries the kind the baseline run observed, so the row joins the two
	// runs it reports rather than asserting the kind the refusal order implies.
	cells := verdictCells{
		cause:    string(outcome.Kind),
		failed:   outcome.FailedTests,
		restored: restoredCell,
		baseline: string(baseline),
		ran:      outcome.Ran,
	}
	var code int
	cells.verdict, code = verdictFor(outcome.Kind)
	preservedBlock := ""
	if !restored {
		cells.verdict, code, cells.restored = "restore-failed", 2, "no"
		preservedBlock = keptCopy(preserved, reason)
	}
	e := selected(subject, mutation, cells, request)
	e.command, e.focused = code, &focused
	out, err := e.answer(preservedBlock + report)
	if err != nil {
		if restored {
			code = 1
		}
		return toon.RenderError(err) + "\n" + preservedBlock, code
	}
	return out, code
}

// envelope is one probe answer's three leading rows: the verdict row, the selection row
// that names what the focused run selected and how many tests the mutated run ran, and the
// execution row. Every render path builds one, so the rows keep one order and one schema.
type envelope struct {
	cells             verdictCells
	subject, mutation string
	form, target, run string
	command           int
	// focused is nil when no mutated run started, so the cell prints null rather than an
	// integer's zero value.
	focused  *int
	trailing int
}

// selected fills an envelope's selection cells from the focused run's request.
func selected(subject subject, mutation string, cells verdictCells, request testreport.Request) envelope {
	return envelope{cells: cells, subject: subject.display, mutation: mutation, form: request.Form(), target: request.Target(), run: request.Run()}
}

// answer renders the envelope ahead of suffix, with the execution row counting every byte
// of suffix. The rows render in one call, so a cell the encoder cannot carry refuses all
// of them rather than printing half an answer.
func (e envelope) answer(suffix string) (string, error) {
	e.trailing = len(suffix)
	head, err := e.rows()
	if err != nil {
		return "", err
	}
	return head + suffix, nil
}

func (e envelope) rows() (string, error) {
	var focused any
	if e.focused != nil {
		focused = *e.focused
	}
	cells := [][]any{
		{e.cells.verdict, e.subject, e.mutation, e.cells.cause, e.cells.failed, e.cells.restored},
		{e.form, e.target, e.run, e.cells.baseline, e.cells.ran},
		{resultSchema, e.command, focused, e.trailing},
	}
	var out string
	for i, table := range resultTables {
		block, err := toon.TableTyped(table.name, table.fields, [][]any{cells[i]})
		if err != nil {
			return "", err
		}
		out += block
	}
	return out, nil
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
