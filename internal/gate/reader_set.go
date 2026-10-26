package gate

// The reader set: a kit landing whose changed paths sit in covered families runs the test
// phase over the readers of those families, the race table when a reader is a race-table
// package, and the system suite. A green reader set writes only the reader record of its
// tree, which never makes the tree reusable green.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/durablefile"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// readerScheduleEnv carries the reader set, as a JSON object, to the one gate-phases process
// that runs it.
const readerScheduleEnv = "BENCH_GATE_READER"

// readerSchedule is the reader set that readerScheduleEnv carries: the reader packages, and
// whether the race table runs.
type readerSchedule struct {
	Packages []string `json:"packages"`
	Race     bool     `json:"race"`
}

// readerScheduleFields is the exact field set of an encoded reader schedule.
var readerScheduleFields = []string{"packages", "race"}

// readerRun is an eligible reader set: the subject whose environment names the reader
// schedule, the schedule as that environment carries it, the destination tree, and the green
// marker commit.
type readerRun struct {
	plan                          subject
	schedule, destination, marker string
}

// newReaderRun answers the reader run of plan over schedule, with the schedule in the
// environment of the phase child.
func newReaderRun(plan subject, schedule readerSchedule, destination, marker string) readerRun {
	// A list of strings and a flag always encode.
	encoded, _ := json.Marshal(schedule)
	plan.Env = mergeEnv(plan.Env, []string{readerScheduleEnv + "=" + string(encoded)})
	return readerRun{plan: plan, schedule: string(encoded), destination: destination, marker: marker}
}

// parseReaderSchedule reads an encoded reader schedule. It refuses any other field set, a
// duplicate or trailing value, and a package operand that does not start with "./", which a
// test command would read as a flag.
func parseReaderSchedule(value string) (readerSchedule, error) {
	var schedule readerSchedule
	data := []byte(value)
	if requireObjectFields(data, readerScheduleFields) != nil || strictJSON(data, &schedule) != nil {
		return readerSchedule{}, errMalformedReaderSchedule
	}
	for _, operand := range schedule.Packages {
		if !strings.HasPrefix(operand, "./") {
			return readerSchedule{}, errMalformedReaderSchedule
		}
	}
	return schedule, nil
}

var errMalformedReaderSchedule = fmt.Errorf("gate: malformed reader schedule in %s", readerScheduleEnv)

// readerScheduled reports whether this process runs the reader set.
func readerScheduled() bool { return os.Getenv(readerScheduleEnv) != "" }

// readerScheduledTable answers the reader set that this process runs, or the refusal of a
// malformed schedule.
func readerScheduledTable(root, kit string) ([]Phase, error) {
	schedule, err := parseReaderSchedule(os.Getenv(readerScheduleEnv))
	if err != nil {
		return nil, err
	}
	return readerPhases(root, kit, schedule), nil
}

// readerPhases is the reader set: the test phase over the reader packages and the
// root-conformance package, each once, and then the race and system phases of the kit table.
// The race phase runs only when the schedule flags it.
func readerPhases(root, kit string, schedule readerSchedule) []Phase {
	var packages []string
	for _, operand := range append(slices.Clone(schedule.Packages), "./"+conformancePackagePath) {
		if !slices.Contains(packages, operand) {
			packages = append(packages, operand)
		}
	}
	phases := []Phase{testPhase(root, kit, packages...)}
	for _, phase := range BenchkitPhases(root, kit) {
		if phase.Name == SystemPhaseName || (schedule.Race && phase.Name == canary.PhaseRace) {
			phases = append(phases, phase)
		}
	}
	return phases
}

// runReader runs the oracle over an eligible reader set. A green run writes the reader record
// of its tree through replacer, and no pending record, no evidence, and no verdict cache. A
// red reader set writes no record.
func runReader(ctx, runCtx context.Context, runtimeRoot, storageRoot string, declaredPaths []string, reader readerRun, evaluation executionEvaluation, local []protection.Descriptor, selection *runbinary.Selection, release func() error, replacer durablefile.Replacer, stdout, stderr io.Writer) Result {
	return runReduced(ctx, runCtx, runtimeRoot, storageRoot, declaredPaths, reader.plan, ScheduleReader, evaluation, local, selection, release, stdout, stderr, func(runBinary string) error {
		gitdir, err := benchgit.CommonDir(storageRoot)
		if err != nil {
			return err
		}
		return writeReaderGreen(replacer, gitdir, reader, runBinary, time.Now())
	})
}

// writeReaderGreen installs the green of run, which the executable at the content address
// runBinary graded at recordedAt, as the reader record of its tree in the store of gitdir.
func writeReaderGreen(replacer durablefile.Replacer, gitdir string, run readerRun, runBinary string, recordedAt time.Time) error {
	dir := filepath.Join(gitdir, benchgit.ReaderRecordDir)
	if err := ensureEvidenceDir(gitdir, dir); err != nil {
		return err
	}
	return durableReplaceAt(replacer, dir, run.plan.Tree, verdictRecord{
		Schema: verdictSchema, Tree: run.plan.Tree, RecordedAt: recordedAt.UTC().Truncate(time.Second).Format(time.RFC3339),
		RunBinary: runBinary, ReaderSet: run.schedule, DestinationTree: run.destination, GreenMarker: run.marker,
	})
}

// ReaderRecorded answers whether the store of root holds a valid reader record for tree. A
// tree that is not a tree hash has no record.
func ReaderRecorded(root, tree string) (bool, error) {
	if !treeHashRE.MatchString(tree) {
		return false, nil
	}
	gitdir, err := benchgit.CommonDir(root)
	if err != nil {
		return false, err
	}
	loaded := loadVerdict(filepath.Join(gitdir, benchgit.ReaderRecordDir, tree), time.Now().UTC())
	if loaded.state == Unavailable {
		return false, errors.New(loaded.reason)
	}
	// Only the reader class carries a green marker, so no other class answers true.
	return loaded.state == Ready && loaded.record.GreenMarker != "" && loaded.record.Tree == tree, nil
}
