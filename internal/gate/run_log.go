package gate

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const gateLogSchema = 1

// gateLogRetainedRuns is how many gate runs .logs keeps. Twenty is a
// reviewer-owned size: enough to hold a full build session's runs, small enough
// that the directory stays readable without a separate pruning chore. A run owns
// both of its files, so the retention counts runs and never half a run.
const gateLogRetainedRuns = 20

// The file name shape, spelled once. A run leaves two files: the record of what it
// did, and the stream of what its phases said. Every producer builds a name through
// gateLogName and every consumer recognizes one by round-tripping the candidate back
// through it, so no caller restates the pattern.
const (
	gateLogRecordPrefix  = "gate-"
	gateLogRecordSuffix  = ".jsonl"
	gateLogStreamSuffix  = ".out"
	gateLogRunTimeLayout = "20060102T150405.000000000Z"
)

const (
	gateLogPathEnv   = "BENCH_GATE_LOG_PATH"
	gateLogRootEnv   = "BENCH_GATE_LOG_ROOT"
	gateLogRunEnv    = "BENCH_GATE_LOG_RUN"
	gateLogStreamEnv = "BENCH_GATE_STREAM_PATH"
)

// gateLogEnv is the whole set a child inherits, so the stripper and the composer
// below read one list rather than two that can disagree.
var gateLogEnv = []string{gateLogPathEnv, gateLogRootEnv, gateLogRunEnv, gateLogStreamEnv}

func gateLogDir(root string) string { return filepath.Join(root, ".logs") }

// gateLogRunToken names one gate run: the run's start instant and the process
// writing it, which together stay unique across concurrent gates on one tree.
func gateLogRunToken(started time.Time, pid int) string {
	return fmt.Sprintf("%s-%d", started.UTC().Format(gateLogRunTimeLayout), pid)
}

func gateLogName(run, suffix string) string {
	return gateLogRecordPrefix + run + suffix
}

func gateLogRecordName(run string) string {
	return gateLogName(run, gateLogRecordSuffix)
}

func gateLogStreamName(run string) string {
	return gateLogName(run, gateLogStreamSuffix)
}

func gateLogRecordPath(root, run string) string {
	return filepath.Join(gateLogDir(root), gateLogRecordName(run))
}

func gateLogStreamPath(root, run string) string {
	return filepath.Join(gateLogDir(root), gateLogStreamName(run))
}

// gateLogRunFromRecordName recognizes a file the gate itself wrote by rebuilding
// the name from its parsed parts and requiring the rebuild to match. Either of a
// run's two suffixes answers that run's token, so the pruner groups a run's files
// without a second parser. Anything that does not round-trip is somebody else's
// file and is not a prune candidate.
func gateLogRunFromRecordName(name string) (run string, started time.Time, ok bool) {
	for _, suffix := range []string{gateLogRecordSuffix, gateLogStreamSuffix} {
		if run, started, ok := gateLogRunFromName(name, suffix); ok {
			return run, started, true
		}
	}
	return "", time.Time{}, false
}

func gateLogRunFromName(name, suffix string) (run string, started time.Time, ok bool) {
	run = strings.TrimSuffix(strings.TrimPrefix(name, gateLogRecordPrefix), suffix)
	if gateLogName(run, suffix) != name {
		return "", time.Time{}, false
	}
	stamp, pid, found := strings.Cut(run, "-")
	if !found {
		return "", time.Time{}, false
	}
	started, err := time.Parse(gateLogRunTimeLayout, stamp)
	if err != nil {
		return "", time.Time{}, false
	}
	id, err := strconv.Atoi(pid)
	if err != nil || gateLogRunToken(started, id) != run {
		return "", time.Time{}, false
	}
	return run, started, true
}

type gateRunLogKey struct{}

type gateRunLog struct {
	mu   sync.Mutex
	file *os.File
	// streamFile holds every line this run's phases wrote, and is nil when the run
	// opened none. It is evidence, not verdict: a run without one keeps its record,
	// its table, and its exit code, and the report says the stream is unavailable.
	streamFile *os.File
	stderr     io.Writer
	run        string
	root       string
	started    time.Time
	warned     bool
}

// streamPath names this run's retained phase stream, or "" when the run opened none.
func (l *gateRunLog) streamPath() string {
	if l == nil || l.streamFile == nil {
		return ""
	}
	return l.streamFile.Name()
}

// gateRunStreamFile answers the retained stream of the run this context carries, so the
// engine's phase buffer writes through to the file the run log owns. A context outside a
// logged run answers nil, and the buffer then stays in memory.
func gateRunStreamFile(ctx context.Context) *os.File {
	log, _ := ctx.Value(gateRunLogKey{}).(*gateRunLog)
	if log == nil {
		return nil
	}
	return log.streamFile
}

type gateLogRecord struct {
	Schema    int    `json:"schema"`
	Time      string `json:"time"`
	Run       string `json:"run"`
	Event     string `json:"event"`
	Root      string `json:"root,omitempty"`
	Mode      string `json:"mode,omitempty"`
	Phase     string `json:"phase,omitempty"`
	Path      string `json:"path,omitempty"`
	Detail    string `json:"detail,omitempty"`
	Exit      *int   `json:"exit,omitempty"`
	ElapsedMS int64  `json:"elapsed_ms,omitempty"`
	// The build-cache footprint fields. They are pointers, like Exit, because a
	// measured zero is a real answer: a fresh machine's cache is empty, and an
	// omitted field would report that run as one that measured nothing.
	Bytes     *int64 `json:"bytes,omitempty"`
	Files     *int64 `json:"files,omitempty"`
	OverBound *bool  `json:"over_bound,omitempty"`
}

func beginGateRunLog(ctx context.Context, root string, stderr io.Writer, mode string) (context.Context, func(Result)) {
	if !gateLogPathIgnored(root) {
		return ctx, func(Result) {}
	}
	if err := os.MkdirAll(gateLogDir(root), 0o700); err != nil {
		fmt.Fprintf(stderr, "gate: progress logging unavailable: %v\n", err)
		return ctx, func(Result) {}
	}
	started := time.Now().UTC()
	run := gateLogRunToken(started, os.Getpid())
	path := gateLogRecordPath(root, run)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL|os.O_APPEND, 0o600)
	if err != nil {
		fmt.Fprintf(stderr, "gate: progress logging unavailable: %v\n", err)
		return ctx, func(Result) {}
	}
	log := &gateRunLog{file: file, streamFile: openGateStreamFile(root, run), stderr: stderr, run: run, root: root, started: started}
	ctx = context.WithValue(ctx, gateRunLogKey{}, log)
	log.write(gateLogRecord{Event: "gate.start", Root: root, Mode: mode, Path: path})
	fmt.Fprintf(stderr, "gate: progress log %s\n", path)
	if stream := log.streamPath(); stream != "" {
		fmt.Fprintf(stderr, "gate: stream %s\n", stream)
	}
	return ctx, log.finish
}

// openGateStreamFile opens the run's phase stream beside its record, and answers nil
// when it cannot. A refused stream is not a refused run: the record, the table, and the
// verdict are unchanged, and the only cost is that the bounded table has no whole to
// point at. So a failure here adds no diagnosis of its own.
//
// A .logs that is itself a symlink receives nothing. MkdirAll follows such a link when
// its target is a directory, and a run must not write its stream through a link out of
// the tree the run names. Lstat does not follow the link, so it fails the test below.
func openGateStreamFile(root, run string) *os.File {
	info, err := os.Lstat(gateLogDir(root))
	if err != nil || !info.Mode().IsDir() {
		return nil
	}
	file, err := os.OpenFile(gateLogStreamPath(root, run), os.O_WRONLY|os.O_CREATE|os.O_EXCL|os.O_APPEND, 0o600)
	if err != nil {
		return nil
	}
	return file
}

// finish closes out the run's record and then prunes, in that order. The record
// this run produced must exist on disk before anything counts what to retain.
func (l *gateRunLog) finish(result Result) {
	exit := result.ActionExit
	l.write(gateLogRecord{Event: "gate.finish", Root: l.root, Exit: &exit, ElapsedMS: time.Since(l.started).Milliseconds()})
	if err := l.file.Close(); err != nil {
		l.warn(err)
	}
	l.closeStream()
	pruneGateRunLogs(l.root, l.run, l.stderr)
}

func (l *gateRunLog) closeStream() {
	if l.streamFile == nil {
		return
	}
	if err := l.streamFile.Close(); err != nil {
		l.warn(err)
	}
}

func withGateRunLogEnv(ctx context.Context, base []string) []string {
	log, _ := ctx.Value(gateRunLogKey{}).(*gateRunLog)
	env := withoutGateRunLogEnv(base)
	if log == nil {
		return env
	}
	env = append(env,
		gateLogPathEnv+"="+log.file.Name(),
		gateLogRootEnv+"="+log.root,
		gateLogRunEnv+"="+log.run,
	)
	// A run whose stream never opened hands the child nothing to append to, and the
	// child then buffers its phases in memory.
	if stream := log.streamPath(); stream != "" {
		env = append(env, gateLogStreamEnv+"="+stream)
	}
	return env
}

func inheritGateRunLog(ctx context.Context, stderr io.Writer) (context.Context, func()) {
	path, root, run := os.Getenv(gateLogPathEnv), os.Getenv(gateLogRootEnv), os.Getenv(gateLogRunEnv)
	if path == "" || root == "" || run == "" || !filepath.IsAbs(path) || !filepath.IsAbs(root) ||
		filepath.Clean(path) != path || filepath.Clean(root) != root ||
		path != gateLogRecordPath(root, run) || !gateLogPathIgnored(root) {
		return ctx, func() {}
	}
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() {
		return ctx, func() {}
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		fmt.Fprintf(stderr, "gate: progress logging unavailable: %v\n", err)
		return ctx, func() {}
	}
	log := &gateRunLog{file: file, streamFile: inheritGateStreamFile(root, run), stderr: stderr, run: run, root: root, started: time.Now().UTC()}
	return context.WithValue(ctx, gateRunLogKey{}, log), func() {
		if err := file.Close(); err != nil {
			log.warn(err)
		}
		log.closeStream()
	}
}

// inheritGateStreamFile reopens the parent's phase stream, so the child that actually
// runs the phases appends to the one file the run names. The record alone decides
// whether a run is inherited: a parent that opened no stream still hands its child a
// run to write records into.
func inheritGateStreamFile(root, run string) *os.File {
	path := gateLogStreamPath(root, run)
	if os.Getenv(gateLogStreamEnv) != path {
		return nil
	}
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() {
		return nil
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		return nil
	}
	return file
}

func withoutGateRunLogEnv(base []string) []string {
	env := make([]string, 0, len(base))
	for _, item := range base {
		if slices.ContainsFunc(gateLogEnv, func(name string) bool { return strings.HasPrefix(item, name+"=") }) {
			continue
		}
		env = append(env, item)
	}
	return env
}

// gateLogPathIgnored reports whether the repository ignores the log directory, which is
// the precondition for writing a record there at all. It is a var so a test can drive
// the run below it without standing up a repository. internal/gate's ordinary tests own
// no repository adapter and start no processes.
var gateLogPathIgnored = func(root string) bool {
	cmd := exec.Command("git", "-C", root, "check-ignore", "-q", "--no-index", ".logs/gate.jsonl")
	return cmd.Run() == nil
}

func logGateEvent(ctx context.Context, record gateLogRecord) {
	log, _ := ctx.Value(gateRunLogKey{}).(*gateRunLog)
	if log != nil {
		log.write(record)
	}
}

func (l *gateRunLog) write(record gateLogRecord) {
	l.mu.Lock()
	defer l.mu.Unlock()
	record.Schema = gateLogSchema
	record.Time = time.Now().UTC().Format(time.RFC3339Nano)
	record.Run = l.run
	data, err := json.Marshal(record)
	if err == nil {
		data = append(data, '\n')
		_, err = l.file.Write(data)
	}
	if err != nil && !l.warned {
		l.warned = true
		fmt.Fprintf(l.stderr, "gate: progress logging unavailable: %v\n", err)
	}
}

func (l *gateRunLog) warn(err error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if !l.warned {
		l.warned = true
		fmt.Fprintf(l.stderr, "gate: progress logging unavailable: %v\n", err)
	}
}

// pruneGateRunLogs bounds .logs to the newest gateLogRetainedRuns gate runs. It
// counts runs, not files: a run owns a record and, when its stream opened, a stream
// beside it, and both go together. A pruner that counted files would keep half that
// many runs, and would leave a table pointing at a stream it had just removed.
//
// Only files the gate itself named are candidates, and only regular ones. A dangling
// symlink or special file wearing a record name is somebody else's problem, so the
// pruner stats before it removes. The current run is ordered with the rest but never
// removed, because pruning must not truncate the evidence being written. Every failure
// is housekeeping, not a verdict, so pruneGateRunLogs warns once and returns.
func pruneGateRunLogs(root, currentRun string, stderr io.Writer) {
	dir := gateLogDir(root)
	warned := false
	warn := func(err error) {
		if !warned {
			warned = true
			fmt.Fprintf(stderr, "gate: log pruning unavailable: %v\n", err)
		}
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		warn(err)
		return
	}
	type candidate struct {
		run     string
		started time.Time
		files   []string
	}
	byRun := make(map[string]*candidate, len(entries))
	runs := make([]*candidate, 0, len(entries))
	for _, entry := range entries {
		run, started, ok := gateLogRunFromRecordName(entry.Name())
		if !ok {
			continue
		}
		info, err := os.Lstat(filepath.Join(dir, entry.Name()))
		if err != nil || !info.Mode().IsRegular() {
			continue
		}
		found := byRun[run]
		if found == nil {
			found = &candidate{run: run, started: started}
			byRun[run] = found
			runs = append(runs, found)
		}
		found.files = append(found.files, entry.Name())
	}
	sort.Slice(runs, func(i, j int) bool {
		if runs[i].started.Equal(runs[j].started) {
			return runs[i].run > runs[j].run
		}
		return runs[i].started.After(runs[j].started)
	})
	if len(runs) <= gateLogRetainedRuns {
		return
	}
	for _, run := range runs[gateLogRetainedRuns:] {
		if run.run == currentRun {
			continue
		}
		for _, name := range run.files {
			if err := os.Remove(filepath.Join(dir, name)); err != nil {
				warn(err)
			}
		}
	}
}
