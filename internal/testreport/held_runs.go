package testreport

import (
	"context"

	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/toon"
)

// HeldRuns serves several focused runs of one prepared request over one held run-binary
// selection. The caller's context governs each selection, build, and Go child.
type HeldRuns struct {
	ctx     context.Context
	root    string
	request focusedRequest
	held    *runbinary.Selection
	builds  int
}

// Hold prepares the held runs of request from root under ctx. The first run selects the run
// binary as Execute does, and Close releases it.
func Hold(ctx context.Context, root string, request Request) *HeldRuns {
	return &HeldRuns{ctx: ctx, root: root, request: request.focused}
}

// Run executes the held request once and returns its outcome, report, and exit. A failed
// selection under a cancelled context reads interrupted, and a fresh build closes after its run.
func (h *HeldRuns) Run() (Outcome, string, int) {
	if outcome, out, code, answered := h.request.answerWithoutSelection(h.root); answered {
		return outcome, out, code
	}
	selection, err := h.selection()
	if err != nil {
		if h.ctx.Err() != nil {
			return interruptedOutcome(toon.Errorf("Bench executable selection interrupted", err.Error())+"\n", 1)
		}
		return selectionFailure(err)
	}
	if selection != h.held {
		defer selection.Close()
	}
	return runSelected(h.ctx, h.root, h.request, selection)
}

// Builds counts the run binaries that the held runs built.
func (h *HeldRuns) Builds() int {
	return h.builds
}

// Close releases the held selection after the last run.
func (h *HeldRuns) Close() error {
	return h.held.Close()
}

// selection answers the run binary of the next run. A held binary that no longer validates
// refuses the run when it is inherited, and an owned one yields to a fresh build.
func (h *HeldRuns) selection() (*runbinary.Selection, error) {
	if h.held == nil {
		selection, err := selectRunBinary(h.ctx, h.root)
		if err != nil {
			return nil, err
		}
		h.held = selection
		if selection.Owned() {
			h.builds++
		}
		return selection, nil
	}
	err := runBinaryFactory.Validate(h.held.SourceRoot, h.held.Path)
	if err == nil {
		return h.held, nil
	}
	if !h.held.Owned() {
		return nil, err
	}
	fresh, err := runBinaryFactory.Own(h.ctx, h.held.SourceRoot)
	if err != nil {
		return nil, err
	}
	h.builds++
	return fresh, nil
}
