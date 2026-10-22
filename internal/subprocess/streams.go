package subprocess

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"syscall"
)

// pipe is one owned stream. The child receives a file endpoint, so the waiter never waits
// for this pipe's worker, and the worker reports its own terminal result.
type pipe struct {
	obligation   Obligation
	child, owner *os.File
	source       io.Reader
	sink         io.Writer
	capture      bool
	// buf holds captured bytes. Only the worker touches it until done is set.
	buf  bytes.Buffer
	done bool
	err  error
}

// openStreams creates the owned pipes and gives the command its child endpoints.
func openStreams(cmd *exec.Cmd, streams Streams) ([]*pipe, error) {
	var pipes []*pipe
	open := func(name Obligation, input bool) (*pipe, error) {
		read, write, err := os.Pipe()
		if err != nil {
			return nil, err
		}
		p := &pipe{obligation: name, child: write, owner: read}
		if input {
			p.child, p.owner = read, write
		}
		pipes = append(pipes, p)
		return p, nil
	}
	output := func(name Obligation, out Output) (io.Writer, error) {
		if !out.Capture {
			switch sink := out.Writer.(type) {
			case nil:
				return nil, nil
			case *os.File:
				return sink, nil
			}
		}
		p, err := open(name, false)
		if err != nil {
			return nil, err
		}
		p.sink, p.capture = out.Writer, out.Capture
		return p.child, nil
	}
	cmd.Stdin, cmd.Stdout, cmd.Stderr = nil, nil, nil
	switch source := streams.Stdin.(type) {
	case nil:
	case *os.File:
		cmd.Stdin = source
	default:
		p, err := open(ObligationStdin, true)
		if err != nil {
			return pipes, err
		}
		p.source = source
		cmd.Stdin = p.child
	}
	stdout, err := output(ObligationStdout, streams.Stdout)
	if err != nil {
		return pipes, err
	}
	cmd.Stdout = stdout
	if streams.MergeStderr {
		cmd.Stderr = stdout
		return pipes, nil
	}
	cmd.Stderr, err = output(ObligationStderr, streams.Stderr)
	return pipes, err
}

// work moves one stream until EOF or failure and then closes the owner endpoint.
func (p *pipe) work() error {
	defer p.owner.Close()
	if p.source != nil {
		_, err := io.Copy(p.owner, p.source)
		if errors.Is(err, syscall.EPIPE) {
			// The child closed its stdin; exec.Cmd treats this as a complete copy too.
			return nil
		}
		return err
	}
	if p.capture {
		_, err := p.buf.ReadFrom(p.owner)
		return err
	}
	_, err := io.Copy(p.sink, p.owner)
	return err
}

// finish records the worker result and drops the caller's reader and writer.
func (p *pipe) finish(err error) {
	p.done, p.err = true, err
	p.source, p.sink = nil, nil
}

// stream projects the pipe into an immutable snapshot. closed reports that the owner
// closed its endpoint, which turns a closed-file error into forced closure.
func (p *pipe) stream(closed bool) Stream {
	if !p.done {
		return Stream{Status: StreamPending}
	}
	s := Stream{Status: StreamComplete, Err: p.err}
	if p.capture {
		s.Data = bytes.Clone(p.buf.Bytes())
	}
	switch {
	case p.err == nil:
	case closed && errors.Is(p.err, os.ErrClosed):
		s.Status = StreamClosed
	default:
		s.Status = StreamFailed
	}
	return s
}
