package shift

// shiftStep names an injectable point in the loop's post-mutation path — mirroring
// internal/worktree's Fault/LifecycleStep pattern so an in-process test can force a
// failure the shell cannot reach (staging, teardown) without adding a shell-visible knob.
type shiftStep string

const (
	stepStage    shiftStep = "stage"
	stepTeardown shiftStep = "teardown"
)

// fault is a step-keyed injection hook: nil in every production path, set only by tests
// in this package to force an error at the named step.
type fault func(shiftStep) error

// shiftFault is the one seam var a test overrides. Production code never assigns it.
var shiftFault fault

// hitShift reports the fault's error for step, or nil when no fault is armed.
func hitShift(f fault, step shiftStep) error {
	if f == nil {
		return nil
	}
	return f(step)
}
