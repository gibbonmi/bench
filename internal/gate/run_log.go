package gate

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const gateLogSchema = 1

const (
	gateLogPathEnv = "BENCH_GATE_LOG_PATH"
	gateLogRootEnv = "BENCH_GATE_LOG_ROOT"
	gateLogRunEnv  = "BENCH_GATE_LOG_RUN"
)

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
	logs := filepath.Join(root, ".logs")
	if err := os.MkdirAll(logs, 0o700); err != nil {
		fmt.Fprintf(stderr, "gate: progress logging unavailable: %v\n", err)
		return ctx, func(Result) {}
	}
	started := time.Now().UTC()
	run := fmt.Sprintf("%s-%d", started.Format("20060102T150405.000000000Z"), os.Getpid())
	path := filepath.Join(logs, "gate-"+run+".jsonl")
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL|os.O_APPEND, 0o600)
	if err != nil {
		fmt.Fprintf(stderr, "gate: progress logging unavailable: %v\n", err)
		return ctx, func(Result) {}
	}
	log := &gateRunLog{file: file, stderr: stderr, run: run, root: root, started: started}
	ctx = context.WithValue(ctx, gateRunLogKey{}, log)
	log.write(gateLogRecord{Event: "gate.start", Root: root, Mode: mode, Path: path})
	fmt.Fprintf(stderr, "gate: progress log %s\n", path)
	return ctx, func(result Result) {
		exit := result.ActionExit
		log.write(gateLogRecord{Event: "gate.finish", Root: root, Exit: &exit, ElapsedMS: time.Since(started).Milliseconds()})
		if err := file.Close(); err != nil {
			log.warn(err)
		}
	}
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
		path != filepath.Join(root, ".logs", "gate-"+run+".jsonl") || !gateLogPathIgnored(root) {
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
