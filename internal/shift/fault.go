package shift

// shiftStep names an injectable point in the loop's post-mutation path. It mirrors
// internal/worktree's Fault/LifecycleStep pattern. An in-process test then forces a
// failure the shell cannot reach, such as staging or teardown, without a shell-visible knob.
type shiftStep string

const (
	stepStage        shiftStep = "stage"
	stepTeardown     shiftStep = "teardown"
	stepIntentUpsert shiftStep = "intent-upsert"
)

// fault is a step-keyed injection hook. It is nil in every production path; a test
// in this package sets it to force an error at the named step.
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
