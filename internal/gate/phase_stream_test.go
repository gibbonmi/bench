package gate

import (
	"io"
	"slices"
	"sync"
	"testing"
)

// One phase's two streams share one buffer, in arrival order. A tool splits a single
// diagnosis across both, so the reader needs the lines in the order they were written.
func TestPhaseStreamsBufferBothStreamsInArrivalOrder(t *testing.T) {
	streams := newPhaseStreams(io.Discard)
	out, errOut, closeWriters := streams.open(Phase{Name: "vet"})
	io.WriteString(out, "first\n")
	io.WriteString(errOut, "second\n")
	io.WriteString(out, "third\n")
	closeWriters()

	want := []string{"first", "second", "third"}
	if got := streams.lines("vet"); !slices.Equal(got, want) {
		t.Errorf("buffered lines = %q, want %q", got, want)
	}
}

// BG15 at the buffer: a last line that carried no newline still flushes as one line when
// the phase's writers close.
func TestPhaseStreamsFlushALastLineWithoutANewlineAtClose(t *testing.T) {
	streams := newPhaseStreams(io.Discard)
	out, _, closeWriters := streams.open(Phase{Name: "gofmt"})
	io.WriteString(out, "whole line\npartial")

	if got := streams.lines("gofmt"); !slices.Equal(got, []string{"whole line"}) {
		t.Errorf("lines before close = %q, want only the newline-terminated one", got)
	}
	closeWriters()
	want := []string{"whole line", "partial"}
	if got := streams.lines("gofmt"); !slices.Equal(got, want) {
		t.Errorf("lines after close = %q, want %q", got, want)
	}
}

// A write arriving in fragments is still one line. The scheduler hands the buffer whatever
// chunks the pipe delivers, and a line boundary is a newline, never a write boundary.
func TestPhaseStreamsSplitOnNewlinesNotOnWrites(t *testing.T) {
	streams := newPhaseStreams(io.Discard)
	out, _, closeWriters := streams.open(Phase{Name: "prose"})
	io.WriteString(out, "one li")
	io.WriteString(out, "ne\ntwo lines\n")
	closeWriters()

	want := []string{"one line", "two lines"}
	if got := streams.lines("prose"); !slices.Equal(got, want) {
		t.Errorf("fragmented lines = %q, want %q", got, want)
	}
}

// Each phase keeps its own buffer, so the report attributes a row to the phase that wrote
// it. A phase that wrote nothing answers with nothing rather than another phase's lines.
func TestPhaseStreamsKeepPhasesApart(t *testing.T) {
	streams := newPhaseStreams(io.Discard)
	out, _, closeWriters := streams.open(Phase{Name: "vet"})
	io.WriteString(out, "vet finding\n")
	closeWriters()

	if got := streams.lines("vet"); !slices.Equal(got, []string{"vet finding"}) {
		t.Errorf("vet lines = %q", got)
	}
	if got := streams.lines("build"); len(got) != 0 {
		t.Errorf("build lines = %q, want none", got)
	}
}

// Phases run concurrently, so the buffer takes lines from several at once. The race
// detector observes this test under the gate's race phase.
func TestPhaseStreamsAcceptConcurrentPhases(t *testing.T) {
	streams := newPhaseStreams(io.Discard)
	names := []string{"gofmt", "vet", "build", "prose"}
	var wg sync.WaitGroup
	for _, name := range names {
		wg.Add(1)
		go func() {
			defer wg.Done()
			out, errOut, closeWriters := streams.open(Phase{Name: name})
			io.WriteString(out, name+" out\n")
			io.WriteString(errOut, name+" err\n")
			closeWriters()
		}()
	}
	wg.Wait()

	for _, name := range names {
		if got := streams.lines(name); len(got) != 2 {
			t.Errorf("phase %s buffered %q, want two lines", name, got)
		}
	}
}
