// Package responsebound owns the bound on one public Bench response. The owner takes the
// response's stdout and stderr writes in arrival order. A response within the line value
// and the byte value prints unchanged. A larger response prints its head, one spill line,
// and its tail, each line cut to lineCut bytes, and a private spill file holds the
// complete output. The owner states the projection once, so the dispatcher and
// `bench worktree exec` cannot drift apart.
package responsebound

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/gibbonmi/bench/internal/bounds"
)

// tailLines is how many last lines an over-bound response prints. The head takes the rest
// of the line value after the spill line, so the projection is exactly the line value.
const tailLines = 5

const headLines = bounds.ResponseLines - tailLines - 1

// lineCut is the most content bytes one printed line keeps, so the printed lines of a
// projection share the byte value evenly.
const lineCut = bounds.ResponseBytes / bounds.ResponseLines

// phase is where one response stands between its first write and its finish.
type phase uint8

const (
	// retaining keeps every write in memory while the response is within the bound.
	retaining phase = iota
	// spilling streams every byte to the spill file after the bound is passed.
	spilling
	// passing writes every byte to its original stream after the spill could not start.
	passing
	// unspilled keeps every byte the spill file refused after a midway write failure.
	unspilled
	// finished passes a late write through to its own stream.
	finished
)

// streamWrite is one retained write and the stream it arrived on.
type streamWrite struct {
	stream io.Writer
	data   []byte
}

// Owner holds one process's public response to the bound. Its two writers serialize
// their writes, because os/exec copies a child's two pipes on two goroutines.
type Owner struct {
	mu             sync.Mutex
	stdout, stderr io.Writer
	home           string
	root           func() string
	retiring       bool
	create         func(path string) (io.WriteCloser, error)

	phase    phase
	retained []streamWrite
	lines    lineTracker
	// stdoutOpen and stderrOpen report that the stream's last byte is not a newline. The
	// create-failure route closes each open stream's line before its stdout line, so the
	// line is whole in each stream and in a sink that takes both.
	stdoutOpen bool
	stderrOpen bool
	bytes      int64
	file       io.WriteCloser
	path       string
	written    int64
	refused    []byte
	reason     string
	// lead is the block that Finish prints first on stdout, outside the bound.
	lead string
	// whole reports a finished response that one copy holds complete: replayed within the
	// bound, or the spill file at path.
	whole    bool
	replayed []byte
}

// New returns the owner of one response whose original streams are stdout and stderr.
// home is the Bench home that holds the spill store. root answers the repository root of
// the process, or the empty string outside a repository. The owner calls it only when a
// spill starts, so a bounded response never pays for the lookup. retiring reports a verb
// that retires an assignment, and that verb spills to the primary scope.
func New(home string, stdout, stderr io.Writer, root func() string, retiring bool) *Owner {
	return &Owner{stdout: stdout, stderr: stderr, home: home, root: root, retiring: retiring, create: exclusiveCreate}
}

// Stdout answers the writer that takes the response's stdout bytes.
func (o *Owner) Stdout() io.Writer { return streamWriter{owner: o, stdout: true} }

// Stderr answers the writer that takes the response's stderr bytes.
func (o *Owner) Stderr() io.Writer { return streamWriter{owner: o} }

// streamWriter tags each write with its original stream.
type streamWriter struct {
	owner  *Owner
	stdout bool
}

// Write accepts every byte. A spill failure changes where the bytes go, never whether
// they arrive, so the command that writes never sees an error from the bound.
func (w streamWriter) Write(p []byte) (int, error) {
	w.owner.write(w.stdout, p)
	return len(p), nil
}

func (o *Owner) write(stdout bool, p []byte) {
	o.mu.Lock()
	defer o.mu.Unlock()
	stream := o.stderr
	if stdout {
		stream = o.stdout
	}
	if o.phase == finished {
		_, _ = stream.Write(p)
		return
	}
	if len(p) > 0 {
		open := p[len(p)-1] != '\n'
		if stdout {
			o.stdoutOpen = open
		} else {
			o.stderrOpen = open
		}
	}
	o.bytes += int64(len(p))
	o.lines.add(p)
	if o.phase == retaining {
		// The write that passes the bound goes straight to the spill, so the owner never
		// copies a long write into memory.
		if !o.overBound() {
			o.retained = append(o.retained, streamWrite{stream: stream, data: append([]byte(nil), p...)})
			return
		}
		o.startSpill()
	}
	switch o.phase {
	case spilling, unspilled:
		o.spill(p)
	case passing:
		_, _ = stream.Write(p)
	}
}

// overBound reports a response past the line value or the byte value.
func (o *Owner) overBound() bool {
	return o.lines.count > bounds.ResponseLines || o.bytes > bounds.ResponseBytes
}

// startSpill moves the retained writes into a new spill file. A file that cannot be
// created sends every retained byte to its original stream instead, so no byte is lost.
func (o *Owner) startSpill() {
	retained := o.retained
	o.retained = nil
	file, path, err := openSpill(o.home, o.root, o.retiring, o.create)
	if err != nil {
		o.phase, o.reason = passing, failureReason(err)
		for _, w := range retained {
			_, _ = w.stream.Write(w.data)
		}
		return
	}
	o.phase, o.file, o.path = spilling, file, path
	for _, w := range retained {
		o.spill(w.data)
	}
}

// spill streams p to the file. After a failed write, it keeps every byte the file did
// not receive, so the file and the response together hold the complete output.
func (o *Owner) spill(p []byte) {
	if o.phase == unspilled {
		o.refused = append(o.refused, p...)
		return
	}
	n, err := o.file.Write(p)
	o.written += int64(n)
	if err == nil && n < len(p) {
		err = io.ErrShortWrite
	}
	if err != nil {
		o.phase, o.reason = unspilled, failureReason(err)
		o.refused = append(o.refused, p[n:]...)
	}
}

// Lead sets the block that Finish prints first on stdout. The block is outside the bound:
// it counts toward neither value, and the spill file does not hold it. A response whose
// spill file could not open has already passed its bytes through, and a lead after them
// would not come first, so Finish prints no lead for that response.
func (o *Owner) Lead(block string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.lead = block
}

// Finish prints the response after the lead. A response within the bound replays each
// write to its own stream in arrival order. An over-bound response prints its head, the
// spill line, and its tail on stdout.
func (o *Owner) Finish() {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.phase != passing && o.phase != finished {
		_, _ = io.WriteString(o.stdout, o.lead)
	}
	switch o.phase {
	case retaining:
		for _, w := range o.retained {
			_, _ = w.stream.Write(w.data)
			o.replayed = append(o.replayed, w.data...)
		}
		o.retained, o.whole = nil, true
	case spilling:
		o.whole = true
		_ = o.file.Close()
		o.printProjection(fmt.Sprintf("%s{lines=%d,bytes=%d,omitted_lines=%d,cut_lines=%d,"+pathField+"%s}\n",
			Spilled, o.lines.count, o.bytes, o.lines.omitted(), o.lines.cutLines(), o.path), o.lines.tailBytes())
	case unspilled:
		_ = o.file.Close()
		o.printProjection(fmt.Sprintf("%s{"+pathField+"%s"+writtenField+"%d,"+reasonField+"%s}\n", SpillFailed, o.path, o.written, o.reason), o.refused)
	case passing:
		if o.stderrOpen {
			_, _ = io.WriteString(o.stderr, "\n")
		}
		separator := ""
		if o.stdoutOpen {
			separator = "\n"
		}
		_, _ = fmt.Fprintf(o.stdout, "%s%s{"+reasonField+"%s}\n", separator, SpillFailed, o.reason)
	}
	o.phase = finished
}

// completeOutput answers the complete bytes of a finished response in arrival order, without
// the lead. A response whose spill failed, at the create or midway, answers false, because no
// one copy holds every byte.
func (o *Owner) completeOutput() ([]byte, bool) {
	o.mu.Lock()
	defer o.mu.Unlock()
	switch {
	case !o.whole:
		return nil, false
	case o.path == "":
		return o.replayed, true
	}
	data, err := os.ReadFile(o.path)
	return data, err == nil
}

func (o *Owner) printProjection(line string, rest []byte) {
	_, _ = o.stdout.Write(o.lines.headBytes())
	_, _ = io.WriteString(o.stdout, line)
	_, _ = o.stdout.Write(rest)
}

// Kind names one reserved frame line that the owner prints in place of complete output.
type Kind string

const (
	// Spilled frames a projection whose complete output is in the named spill file.
	Spilled Kind = "spilled"
	// SpillFailed frames a response whose spill file could not take the complete output.
	SpillFailed Kind = "spill-failed"
)

// pathField, writtenField, and reasonField are the field spellings that a frame reader
// keys on.
const (
	pathField    = "path="
	writtenField = ",written_bytes="
	reasonField  = "reason="
)

// ErrIncomplete marks data that holds a reserved frame line, or a malformed one.
var ErrIncomplete = errors.New("incomplete response")

// Projection is one reserved frame line. Path is the file that the frame names, and it is
// empty for a spill that could not open.
type Projection struct {
	Kind Kind
	Path string
}

// ProjectionLine classifies one physical line, with or without its newline. A line outside
// the reserved frames answers false. A line that opens a reserved frame but does not close
// it, or whose fields do not match the shape that Finish prints for its kind, answers true
// with an error. The reader checks only the frame, the path, and the field order; Finish
// owns the count fields.
func ProjectionLine(line string) (Projection, bool, error) {
	line = strings.TrimSuffix(line, "\n")
	for _, kind := range []Kind{Spilled, SpillFailed} {
		body, ok := strings.CutPrefix(line, string(kind)+"{")
		if !ok {
			continue
		}
		projection := Projection{Kind: kind}
		body, ok = strings.CutSuffix(body, "}")
		if !ok {
			return projection, true, fmt.Errorf("%s frame does not close", kind)
		}
		switch {
		case kind == Spilled:
			_, projection.Path, ok = strings.Cut(body, ","+pathField)
			ok = ok && projection.Path != ""
		case strings.HasPrefix(body, pathField):
			var rest string
			projection.Path, rest, ok = strings.Cut(strings.TrimPrefix(body, pathField), writtenField)
			ok = ok && projection.Path != "" && strings.Contains(rest, ","+reasonField)
		default:
			ok = strings.HasPrefix(body, reasonField)
		}
		if !ok {
			return projection, true, fmt.Errorf("%s frame is malformed", kind)
		}
		return projection, true, nil
	}
	return Projection{}, false, nil
}

// First finds the first physical line of data that opens a reserved frame. It answers the
// frame, its one-based line number, and the error of a malformed frame. Line 0 means that
// no line opens a frame.
func First(data []byte) (Projection, int, error) {
	for number, line := range strings.SplitAfter(string(data), "\n") {
		if projection, ok, err := ProjectionLine(line); ok {
			return projection, number + 1, err
		}
	}
	return Projection{}, 0, nil
}

// RequireComplete refuses data that holds a reserved frame on any physical line, so a
// caller never reads a projection as the complete output. It never opens a path that a
// frame names, because the selected bytes are the only evidence the caller supplied.
func RequireComplete(data []byte) error {
	projection, line, err := First(data)
	switch {
	case line == 0:
		return nil
	case err != nil:
		return fmt.Errorf("%w: line %d: %v", ErrIncomplete, line, err)
	}
	return fmt.Errorf("%w: line %d is a %s frame", ErrIncomplete, line, projection.Kind)
}
