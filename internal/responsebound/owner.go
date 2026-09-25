// Package responsebound owns the bound on one public Bench response. The owner takes the
// response's stdout and stderr writes in arrival order. A response within the line value
// prints unchanged. A longer response prints its head, one spill line, and its tail, and a
// private spill file holds the complete output. The owner states the projection once, so
// the dispatcher and `bench worktree exec` cannot drift apart.
package responsebound

import (
	"fmt"
	"io"
	"sync"

	"github.com/gibbonmi/bench/internal/bounds"
)

// tailLines is how many last lines an over-bound response prints. The head takes the rest
// of the line value after the spill line, so the projection is exactly the line value.
const tailLines = 5

const headLines = bounds.ResponseLines - tailLines - 1

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
	root           func() string
	create         func(path string) (io.WriteCloser, error)

	phase    phase
	retained []streamWrite
	lines    lineTracker
	// stdoutOpen reports that the last stdout byte is not a newline. The create-failure
	// line reads it, because stderr writes change the response's line state but not
	// stdout's.
	stdoutOpen bool
	bytes      int64
	file       io.WriteCloser
	path       string
	written    int64
	refused    []byte
	reason     string
}

// New returns the owner of one response whose original streams are stdout and stderr.
// root answers the repository root of the process, or the empty string outside a
// repository. The owner calls it only when a spill starts, so a bounded response never
// pays for the lookup.
func New(stdout, stderr io.Writer, root func() string) *Owner {
	return &Owner{stdout: stdout, stderr: stderr, root: root, create: exclusiveCreate}
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
	if stdout && len(p) > 0 {
		o.stdoutOpen = p[len(p)-1] != '\n'
	}
	o.bytes += int64(len(p))
	o.lines.add(p)
	switch o.phase {
	case retaining:
		o.retained = append(o.retained, streamWrite{stream: stream, data: append([]byte(nil), p...)})
		if o.lines.count > bounds.ResponseLines {
			o.startSpill()
		}
	case spilling, unspilled:
		o.spill(p)
	case passing:
		_, _ = stream.Write(p)
	}
}

// startSpill moves the retained writes into a new spill file. A file that cannot be
// created sends every retained byte to its original stream instead, so no byte is lost.
func (o *Owner) startSpill() {
	retained := o.retained
	o.retained = nil
	file, path, err := openSpill(o.root, o.create)
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

// Finish prints the response. A response within the bound replays each write to its own
// stream in arrival order. An over-bound response prints its head, the spill line, and
// its tail on stdout.
func (o *Owner) Finish() {
	o.mu.Lock()
	defer o.mu.Unlock()
	switch o.phase {
	case retaining:
		for _, w := range o.retained {
			_, _ = w.stream.Write(w.data)
		}
		o.retained = nil
	case spilling:
		_ = o.file.Close()
		o.printProjection(fmt.Sprintf("spilled{lines=%d,bytes=%d,omitted_lines=%d,cut_lines=%d,path=%s}\n",
			o.lines.count, o.bytes, o.lines.omitted(), 0, o.path), o.lines.tail())
	case unspilled:
		_ = o.file.Close()
		o.printProjection(fmt.Sprintf("spill-failed{path=%s,written_bytes=%d,reason=%s}\n", o.path, o.written, o.reason), o.refused)
	case passing:
		separator := ""
		if o.stdoutOpen {
			separator = "\n"
		}
		_, _ = fmt.Fprintf(o.stdout, "%sspill-failed{reason=%s}\n", separator, o.reason)
	}
	o.phase = finished
}

func (o *Owner) printProjection(line string, rest []byte) {
	for _, head := range o.lines.head {
		_, _ = o.stdout.Write(head)
	}
	_, _ = io.WriteString(o.stdout, line)
	_, _ = o.stdout.Write(rest)
}
