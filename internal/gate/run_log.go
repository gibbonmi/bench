package gate

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const gateLogSchema = 1

// gateLogRetainedRecords is how many gate run records .logs keeps. Twenty is a
// reviewer-owned size: enough to hold a full build session's runs, small enough
// that the directory stays readable without a separate pruning chore.
const gateLogRetainedRecords = 20

// The record name shape, spelled once. Every producer builds a name through
// gateLogRecordName and every consumer recognizes one by round-tripping the
// candidate back through it, so no caller restates the pattern.
const (
	gateLogRecordPrefix  = "gate-"
	gateLogRecordSuffix  = ".jsonl"
	gateLogRunTimeLayout = "20060102T150405.000000000Z"
)

const (
	gateLogPathEnv = "BENCH_GATE_LOG_PATH"
	gateLogRootEnv = "BENCH_GATE_LOG_ROOT"
	gateLogRunEnv  = "BENCH_GATE_LOG_RUN"
)

func gateLogDir(root string) string { return filepath.Join(root, ".logs") }

// gateLogRunToken names one gate run: the run's start instant and the process
// writing it, which together stay unique across concurrent gates on one tree.
func gateLogRunToken(started time.Time, pid int) string {
	return fmt.Sprintf("%s-%d", started.UTC().Format(gateLogRunTimeLayout), pid)
}

func gateLogRecordName(run string) string {
	return gateLogRecordPrefix + run + gateLogRecordSuffix
}

func gateLogRecordPath(root, run string) string {
	return filepath.Join(gateLogDir(root), gateLogRecordName(run))
}

// gateLogRunFromRecordName recognizes a file the gate itself wrote by rebuilding
// the name from its parsed parts and requiring the rebuild to match. Anything
// that does not round-trip is somebody else's file and is not a prune candidate.
func gateLogRunFromRecordName(name string) (run string, started time.Time, ok bool) {
	run = strings.TrimSuffix(strings.TrimPrefix(name, gateLogRecordPrefix), gateLogRecordSuffix)
	if gateLogRecordName(run) != name {
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
	mu      sync.Mutex
	file    *os.File
	stderr  io.Writer
	run     string
	root    string
	started time.Time
	warned  bool
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
	log := &gateRunLog{file: file, stderr: stderr, run: run, root: root, started: started}
	ctx = context.WithValue(ctx, gateRunLogKey{}, log)
	log.write(gateLogRecord{Event: "gate.start", Root: root, Mode: mode, Path: path})
	fmt.Fprintf(stderr, "gate: progress log %s\n", path)
	return ctx, log.finish
}

// finish closes out the run's record and then prunes, in that order: the record
// this run produced must exist on disk before anything counts what to retain.
func (l *gateRunLog) finish(result Result) {
	exit := result.ActionExit
	l.write(gateLogRecord{Event: "gate.finish", Root: l.root, Exit: &exit, ElapsedMS: time.Since(l.started).Milliseconds()})
	if err := l.file.Close(); err != nil {
		l.warn(err)
	}
	pruneGateRunLogs(l.root, l.run, l.stderr)
}

func withGateRunLogEnv(ctx context.Context, base []string) []string {
	log, _ := ctx.Value(gateRunLogKey{}).(*gateRunLog)
	env := withoutGateRunLogEnv(base)
	if log == nil {
		return env
	}
	return append(env,
		gateLogPathEnv+"="+log.file.Name(),
		gateLogRootEnv+"="+log.root,
		gateLogRunEnv+"="+log.run,
	)
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
	log := &gateRunLog{file: file, stderr: stderr, run: run, root: root, started: time.Now().UTC()}
	return context.WithValue(ctx, gateRunLogKey{}, log), func() {
		if err := file.Close(); err != nil {
			log.warn(err)
		}
	}
}

func withoutGateRunLogEnv(base []string) []string {
	env := make([]string, 0, len(base))
	for _, item := range base {
		if strings.HasPrefix(item, gateLogPathEnv+"=") || strings.HasPrefix(item, gateLogRootEnv+"=") || strings.HasPrefix(item, gateLogRunEnv+"=") {
			continue
		}
		env = append(env, item)
	}
	return env
}

func gateLogPathIgnored(root string) bool {
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

// pruneGateRunLogs bounds .logs to the newest gateLogRetainedRecords gate records.
// Only files the gate itself named are candidates, and only regular ones — a
// dangling symlink or special file wearing a record name is somebody else's
// problem, so the pruner stats before it removes. The current run is ordered with
// the rest but never removed: pruning must not truncate the evidence being written.
// Every failure is housekeeping, not a verdict, so it warns once and returns.
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
		name    string
		run     string
		started time.Time
	}
	records := make([]candidate, 0, len(entries))
	for _, entry := range entries {
		run, started, ok := gateLogRunFromRecordName(entry.Name())
		if !ok {
			continue
		}
		info, err := os.Lstat(filepath.Join(dir, entry.Name()))
		if err != nil || !info.Mode().IsRegular() {
			continue
		}
		records = append(records, candidate{name: entry.Name(), run: run, started: started})
	}
	sort.Slice(records, func(i, j int) bool {
		if records[i].started.Equal(records[j].started) {
			return records[i].name > records[j].name
		}
		return records[i].started.After(records[j].started)
	})
	if len(records) <= gateLogRetainedRecords {
		return
	}
	for _, record := range records[gateLogRetainedRecords:] {
		if record.run == currentRun {
			continue
		}
		if err := os.Remove(filepath.Join(dir, record.name)); err != nil {
			warn(err)
		}
	}
}
