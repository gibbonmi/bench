package worktree

// This file holds the one way a worktree test runs a verb: a typed key selects the verb,
// a call value carries the inputs, and a result carries the exit code and both streams.
// The core readers return an error; a test reads a result only through the must forms,
// which fail the test on a reader result they do not accept.

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/axi/axitest"
)

// verbKey names one worktree verb that a test runs.
type verbKey string

const (
	verbCreate      verbKey = "create"
	verbRelease     verbKey = "release"
	verbClean       verbKey = "clean"
	verbReclaim     verbKey = "reclaim"
	verbReauthorize verbKey = "reauthorize"
	verbMerge       verbKey = "merge"
	verbReset       verbKey = "reset"
	verbLand        verbKey = "land"
	verbLandResume  verbKey = "land-resume"
	verbResumeClean verbKey = "resume-clean"
	verbList        verbKey = "list"
	verbPath        verbKey = "path"
	verbShow        verbKey = "show"
	verbBuild       verbKey = "build"
	verbExec        verbKey = "exec"
	verbPool        verbKey = "pool"
	verbLeaseFile   verbKey = "lease-file"
)

// verbCall is the input of one verb run. A call with no joins value, no kit value, and no
// clock value runs the verb entry. Any other call runs the verb's joins form, with the
// joins value unchanged or defaultJoins() when the call holds none. The kit value and the
// clock value replace the reads of the ambient value that the joins form receives.
type verbCall struct {
	root  string
	home  string
	kit   string
	clock func() time.Time
	stdin io.Reader
	joins *joins
	args  []string
}

// verbResult is the output of one verb run. A verb that returns its output as a string
// fills stdout and leaves stderr empty. Only the exec verb fills assignment.
type verbResult struct {
	exit       int
	stdout     string
	stderr     string
	assignment string
}

// verbEntry runs a verb's public entry, and joinsForm runs the internal form that takes
// the seam set and the ambient value.
type (
	verbEntry func(call verbCall, stdout, stderr io.Writer) (int, string)
	joinsForm func(j joins, a ambient, root string, args []string, stdout, stderr io.Writer) int
)

// verbForm is one key's verb entry and, for a verb that has one, its joins form.
type verbForm struct {
	entry  verbEntry
	joined joinsForm
}

// streamed adapts an entry that writes both streams.
func streamed(entry func(root, home string, args []string, stdout, stderr io.Writer) int) verbEntry {
	return func(call verbCall, stdout, stderr io.Writer) (int, string) {
		return entry(call.root, call.home, call.args, stdout, stderr), ""
	}
}

// rendered adapts an entry that returns its output as a string.
func rendered(entry func(call verbCall) (string, int)) verbEntry {
	return func(call verbCall, stdout, _ io.Writer) (int, string) {
		out, exit := entry(call)
		_, _ = io.WriteString(stdout, out)
		return exit, ""
	}
}

var verbForms = map[verbKey]verbForm{
	verbCreate:      {entry: streamed(CreateCommand)},
	verbRelease:     {entry: streamed(ReleaseCommand), joined: releaseCommandWith},
	verbClean:       {entry: streamed(CleanCommand), joined: cleanCommandWith},
	verbReclaim:     {entry: streamed(ReclaimCommand)},
	verbReauthorize: {entry: streamed(ReauthorizeCommand), joined: reauthorizeWith},
	verbMerge:       {entry: streamed(MergeCommand), joined: mergeWith},
	verbReset:       {entry: streamed(ResetCommand), joined: resetWith},
	verbLand:        {entry: streamed(LandCommand), joined: landWith},
	verbLandResume:  {entry: streamed(ResumeLandCommand), joined: resumeLandWith},
	verbResumeClean: {entry: streamed(ResumeCleanCommand), joined: resumeCleanCommandWith},
	verbList:        {entry: rendered(func(call verbCall) (string, int) { return ListCommand(call.root, call.home, call.args) })},
	verbPath:        {entry: streamed(PathCommand)},
	verbShow:        {entry: streamed(ShowCommand)},
	verbBuild:       {entry: streamed(BuildCommand), joined: buildWith},
	verbExec: {entry: func(call verbCall, stdout, stderr io.Writer) (int, string) {
		return ExecCommandResolving(call.root, call.home, call.args, call.stdin, stdout, stderr)
	}},
	verbPool:      {entry: rendered(func(call verbCall) (string, int) { return PoolCommand(call.home, call.args) })},
	verbLeaseFile: {entry: rendered(func(call verbCall) (string, int) { return LeaseFileCommand(call.args) })},
}

// runVerb runs the verb that key names and returns its result. It fails t on a call that
// checkVerbCall refuses, on an unknown key, and on a joins value for a verb without a
// joins form.
func runVerb(t testing.TB, key verbKey, call verbCall) verbResult {
	t.Helper()
	if err := checkVerbCall(key, call); err != nil {
		t.Fatalf("%v", err)
		return verbResult{}
	}
	form, ok := verbForms[key]
	if !ok {
		t.Fatalf("verb runner: no verb key %q", key)
		return verbResult{}
	}
	var stdout, stderr bytes.Buffer
	var result verbResult
	if call.joins == nil && call.kit == "" && call.clock == nil {
		result.exit, result.assignment = form.entry(call, &stdout, &stderr)
	} else {
		if form.joined == nil {
			t.Fatalf("verb runner: the %s verb has no joins form", key)
			return verbResult{}
		}
		j := defaultJoins()
		if call.joins != nil {
			j = *call.joins
		}
		result.exit = form.joined(j, callAmbient(call, &stderr), call.root, call.args, &stdout, &stderr)
	}
	result.stdout, result.stderr = stdout.String(), stderr.String()
	return result
}

// callAmbient builds the ambient value with the constructor that each verb entry calls.
// The call's kit value and clock value replace the constructor's reads.
func callAmbient(call verbCall, stderr io.Writer) ambient {
	a := newAmbient(call.home, stderr)
	if call.kit != "" {
		a.kit = call.kit
	}
	if call.clock != nil {
		a.now = call.clock()
	}
	return a
}

// checkVerbCall refuses a kit value or a clock value for a verb key without a joins form,
// because only a joins form receives the ambient value. It does not check that a joins
// form reads either value.
func checkVerbCall(key verbKey, call verbCall) error {
	if verbForms[key].joined != nil {
		return nil
	}
	if call.kit != "" {
		return fmt.Errorf("verb runner: the %s verb takes no kit value", key)
	}
	if call.clock != nil {
		return fmt.Errorf("verb runner: the %s verb takes no clock value", key)
	}
	return nil
}

// errNoVerbFingerprint is the one reader error that mustNoFingerprint accepts.
var errNoVerbFingerprint = errors.New("verb result carries no fingerprint")

// readVerbRows decodes the whole of stdout and returns the rows of the named table block.
func readVerbRows(stdout, block string) ([]any, error) {
	document, err := axitest.DecodeDocument(stdout)
	if err != nil {
		return nil, err
	}
	return document.Rows(block)
}

// readVerbFingerprint returns the one fingerprint that stdout carries. It reads the
// `fingerprint` cell of every table row when stdout decodes, and otherwise the
// `fingerprint=` cell of each record line. A plan that no apply can name still writes the
// cell, as an empty value or as the unapplicable placeholder, so both read as absent.
func readVerbFingerprint(stdout string) (string, error) {
	var values []string
	if document, err := axitest.DecodeDocument(stdout); err == nil {
		for _, block := range document.Blocks {
			rows, _ := document.Values[block].([]any)
			for _, row := range rows {
				fields, _ := row.(map[string]any)
				cell, present := fields["fingerprint"]
				if !present {
					continue
				}
				value, ok := cell.(string)
				if !ok {
					return "", fmt.Errorf("fingerprint cell %#v in block %q is not text", cell, block)
				}
				values = append(values, value)
			}
		}
	}
	if len(values) == 0 {
		values = recordFingerprints(stdout)
	}
	for _, value := range values {
		if value != values[0] {
			return "", fmt.Errorf("verb result carries two fingerprints, %q and %q", values[0], value)
		}
	}
	if len(values) == 0 || values[0] == "" || values[0] == unapplicableFingerprint {
		return "", errNoVerbFingerprint
	}
	return values[0], nil
}

// recordFingerprints returns the `fingerprint=` cell of each `name{key=value,...}` record
// line. The cell ends at the next comma or at the closing brace.
func recordFingerprints(stdout string) []string {
	var values []string
	for _, line := range strings.Split(stdout, "\n") {
		if !strings.HasSuffix(line, "}") {
			continue
		}
		for _, lead := range []string{"{fingerprint=", ",fingerprint="} {
			if _, rest, found := strings.Cut(line, lead); found {
				value, _, _ := strings.Cut(strings.TrimSuffix(rest, "}"), ",")
				values = append(values, value)
				break
			}
		}
	}
	return values
}

// mustRows returns the rows of the named table block, or fails t.
func (r verbResult) mustRows(t testing.TB, block string) []any {
	t.Helper()
	rows, err := readVerbRows(r.stdout, block)
	if err != nil {
		t.Fatalf("verb result rows of %q: %v\nstdout:\n%s", block, err, r.stdout)
		return nil
	}
	return rows
}

// mustFingerprint returns the one fingerprint of the result, or fails t.
func (r verbResult) mustFingerprint(t testing.TB) string {
	t.Helper()
	value, err := readVerbFingerprint(r.stdout)
	if err != nil {
		t.Fatalf("verb result fingerprint: %v\nstdout:\n%s", err, r.stdout)
		return ""
	}
	return value
}

// mustNoFingerprint fails t unless the result carries no fingerprint. Any other reader
// error also fails t, so a conflict does not pass as an absent value.
func (r verbResult) mustNoFingerprint(t testing.TB) {
	t.Helper()
	value, err := readVerbFingerprint(r.stdout)
	if errors.Is(err, errNoVerbFingerprint) {
		return
	}
	t.Fatalf("verb result fingerprint = %q, %v; want no fingerprint\nstdout:\n%s", value, err, r.stdout)
}
