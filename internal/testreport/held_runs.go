package testreport

import (
	"context"

	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/toon"
)

// HeldRuns serves several focused runs over one held run-binary selection. Each run takes the
// held request or a request of its own. The caller's context governs each selection, build,
// and Go child.
type HeldRuns struct {
	ctx     context.Context
	root    string
	request Request
	held    *runbinary.Selection
	builds  int
}

// Hold prepares the held runs of request from root under ctx. The first run selects the run
// binary as Execute does, and Close releases it.
func Hold(ctx context.Context, root string, request Request) *HeldRuns {
	return &HeldRuns{ctx: ctx, root: root, request: request}
}

// Run executes the held request once through RunRequest.
func (h *HeldRuns) Run() (Outcome, string, int) {
	return h.RunRequest(h.request)
}

// RunRequest executes request once over the held selection and returns its outcome, report,
// and exit. A failed selection under a cancelled context reads interrupted, and a fresh build
// closes after its run. A failed close of that build stays beside the run's report and fails
// the run.
func (h *HeldRuns) RunRequest(request Request) (outcome Outcome, out string, code int) {
	if outcome, out, code, answered := request.focused.answerWithoutSelection(h.root); answered {
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
		defer func() {
			if err := selection.Close(); err != nil {
				outcome, out, code = withCleanupFailure(outcome, out, code, SelectionCleanupFailure(err))
			}
		}()
	}
	return runSelected(h.ctx, h.root, request.focused, selection)
}

// Builds counts the run binaries that the held runs built.
func (h *HeldRuns) Builds() int {
	return h.builds
}

// Close releases the held selection after the last run. The caller consumes its error as a
// cleanup failure beside the runs' results.
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
