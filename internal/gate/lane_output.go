package gate

// The lane's output plumbing: the per-check diagnostic keeper and the tap writer that
// feeds it while the run streams. The lane's table and its run live beside this file.

import (
	"bytes"
	"io"
	"strings"
	"sync"
)

// laneDiagnostics keeps each check's first output line while the run streams. A caller
// naming a failing check needs one line of why, and re-reading the stream afterwards
// would need the whole stream retained.
type laneDiagnostics struct {
	mu    sync.Mutex
	first map[string]string
}

func (d *laneDiagnostics) record(name, line string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, held := d.first[name]; !held {
		d.first[name] = line
	}
}

// firstLine answers the check's first output line, and falls back to the start error for
// a check that never ran and so wrote nothing.
func (d *laneDiagnostics) firstLine(name string, startErr error) string {
	d.mu.Lock()
	defer d.mu.Unlock()
	if line, held := d.first[name]; held {
		return line
	}
	if startErr != nil {
		return startErr.Error()
	}
	return ""
}

// laneWriters is the lane's output plumbing. It keeps the gate's prefixed streams and
// taps each of them on the way through. The tap holds only the current partial line, so
// a chatty check costs the lane no memory.
func laneWriters(stdout, stderr io.Writer, diagnostics *laneDiagnostics) func(Phase) (io.Writer, io.Writer, func()) {
	inner := prefixedPhaseWriters(discardIfNil(stdout), discardIfNil(stderr))
	return func(check Phase) (io.Writer, io.Writer, func()) {
		out, errOut, closeWriters := inner(check)
		tapOut := &laneTapWriter{name: check.Name, dst: out, diagnostics: diagnostics}
		tapErr := &laneTapWriter{name: check.Name, dst: errOut, diagnostics: diagnostics}
		return tapOut, tapErr, func() {
			// A final line with no newline after it still names the failure, so it is
			// offered before the underlying writers close.
			tapOut.flush()
			tapErr.flush()
			closeWriters()
		}
	}
}

func discardIfNil(w io.Writer) io.Writer {
	if w == nil {
		return io.Discard
	}
	return w
}

type laneTapWriter struct {
	name        string
	dst         io.Writer
	diagnostics *laneDiagnostics
	pending     []byte
}

func (w *laneTapWriter) Write(p []byte) (int, error) {
	w.pending = append(w.pending, p...)
	for {
		idx := bytes.IndexByte(w.pending, '\n')
		if idx < 0 {
			break
		}
		w.diagnostics.record(w.name, string(w.pending[:idx]))
		w.pending = w.pending[idx+1:]
	}
	return w.dst.Write(p)
}

func (w *laneTapWriter) flush() {
	if len(w.pending) == 0 {
		return
	}
	w.diagnostics.record(w.name, string(w.pending))
	w.pending = nil
}
