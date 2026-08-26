package gate

// The engine's phase output plumbing. The engine buffers each phase's lines instead
// of relaying them, so a green phase costs stdout nothing and a red phase's stream is
// there for the report to project when the phase settles. The fast lane keeps the
// prefixed relay in lane.go; only the engine buffers.

import (
	"bytes"
	"io"
	"os"
	"sync"
)

// phaseStreams holds one run's per-phase line buffers. A phase's stdout and its stderr
// share one buffer, because a tool splits a single diagnosis across both streams and the
// reader needs the lines in the order they arrived.
type phaseStreams struct {
	mu      sync.Mutex
	buffers map[string][]string
	// stderr carries the buffer's own diagnostics, which is a stream separate from any
	// phase's. It stays quiet: the run log owner names the retained stream file when the
	// file opens, and reports the directory it cannot use when the file does not.
	stderr io.Writer
	// file is the run's retained stream, or nil when the run opened none. The buffer
	// writes each line through as the line arrives, so a killed run keeps what its
	// phases already said.
	file *os.File
}

func newPhaseStreams(stderr io.Writer) *phaseStreams {
	return &phaseStreams{buffers: make(map[string][]string), stderr: stderr}
}

// retain names the file every line is written through to. A run that opened none stays
// in memory, and its report says the stream is unavailable rather than naming a file
// that holds nothing.
func (s *phaseStreams) retain(file *os.File) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.file = file
}

// open is the writer factory the engine hands to schedule, in place of the lane's
// prefixed relay. The returned close flushes a last line that carried no newline, so a
// tool killed mid-line still contributes what it managed to say.
func (s *phaseStreams) open(phase Phase) (io.Writer, io.Writer, func()) {
	out := &phaseLineWriter{streams: s, phase: phase.Name}
	errOut := &phaseLineWriter{streams: s, phase: phase.Name}
	return out, errOut, func() { out.Close(); errOut.Close() }
}

// lines answers one phase's buffered lines in arrival order. The copy keeps a caller
// from aliasing a buffer a still-running phase may append to.
func (s *phaseStreams) lines(phase string) []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.buffers[phase]...)
}

// path names the file that holds every buffered line, or answers "" when this run
// retained none. It is the one place the report asks where the lines the row cap left
// out can be read.
func (s *phaseStreams) path() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.file == nil {
		return ""
	}
	return s.file.Name()
}

// appendLine files one line under its phase and writes it through to the retained
// stream in the same step. The write happens as the line arrives rather than when the
// phase settles, so a killed run's file holds everything that reached it. Each line
// carries its phase the way the fast lane's relay prefixes one, so one file reads as
// one run. A write that fails costs the reader the whole, never the verdict, so it is
// silent: the run log owner already reported a .logs it cannot use.
func (s *phaseStreams) appendLine(phase, line string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.buffers[phase] = append(s.buffers[phase], line)
	if s.file != nil {
		_, _ = io.WriteString(s.file, "["+phase+"] "+line+"\n")
	}
}

// phaseLineWriter splits one stream into lines the way prefixWriter does, and appends
// them to the phase's buffer rather than writing them through. The pending bytes need no
// lock: exec gives each stream its own copier goroutine, so one writer has one author.
type phaseLineWriter struct {
	streams *phaseStreams
	phase   string
	buf     []byte
}

func (w *phaseLineWriter) Write(p []byte) (int, error) {
	w.buf = append(w.buf, p...)
	for {
		idx := bytes.IndexByte(w.buf, '\n')
		if idx < 0 {
			return len(p), nil
		}
		w.streams.appendLine(w.phase, string(w.buf[:idx]))
		w.buf = w.buf[idx+1:]
	}
}

func (w *phaseLineWriter) Close() {
	if len(w.buf) == 0 {
		return
	}
	w.streams.appendLine(w.phase, string(w.buf))
	w.buf = nil
}
