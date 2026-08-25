// The identity component registry: the named checks a lifecycle verb runs against an
// assignment's own records, and the one constructor that turns a failed check into a
// refusal.
package worktree

import "github.com/gibbonmi/bench/internal/intent"

// The component names. An operator reads these in guidance and in a refusal's detail,
// so the string and the registry entry are the same fact.
const (
	componentRequest         = "request"
	componentAssignmentState = "assignment-state"
	componentAssignmentPath  = "assignment-path"
	componentOwnerMarker     = "owner-marker"
	componentRegistration    = "registration"
	componentLock            = "lock"
)

// identityComponent is one check in the identity bundle. detail renders the check's
// sentence for an assignment id, and recovers records whether the component can name a
// command that repairs it.
type identityComponent struct {
	name     string
	detail   func(assignment string) string
	recovers bool
}

// identityComponents is the declared registry. Its order is the precedence a bundle
// check follows, so a bundle that fails several checks names the earliest one. The
// registry test walks this slice and requires a producing fixture per entry, so a new
// component without a fixture turns the gate red.
var identityComponents = []identityComponent{
	{name: componentRequest, detail: func(string) string { return "request token matches no assignment" }, recovers: true},
	{name: componentAssignmentState, detail: func(id string) string { return "assignment " + id + " is not active" }},
	{name: componentAssignmentPath, detail: func(id string) string { return "assignment " + id + " owns another worktree" }},
	{name: componentOwnerMarker, detail: func(id string) string { return "owner marker does not match assignment " + id }},
	{name: componentRegistration, detail: func(id string) string { return "worktree registration does not match assignment " + id }},
	{name: componentLock, detail: func(id string) string { return "Bench lock does not match assignment " + id }},
}

func identityComponentByName(name string) identityComponent {
	for _, component := range identityComponents {
		if component.name == name {
			return component
		}
	}
	// A name that is not in the registry is a programming fault in this package, not an
	// operator condition. It surfaces as a refusal rather than a panic, because a
	// lifecycle verb must not abort a caller's session over its own bookkeeping.
	return identityComponent{name: name, detail: func(string) string { return "identity component " + name + " is unregistered" }}
}

// componentRefusal is the one constructor every producing site calls. No site composes
// a detail sentence of its own, so the registry stays the single source of the text an
// operator reads.
func componentRefusal(name, assignment, observed, wanted string) refusalError {
	return refusalError{refusal{
		detail:   identityComponentByName(name).detail(assignment),
		observed: observed,
		wanted:   wanted,
	}}
}

// identityBundleRefusal names the first component the assignment's own records fail, in
// registry order. active decides which ledger states the calling verb accepts: a first
// landing takes `active` alone, and a resume also takes `cleanup-pending`.
func identityBundleRefusal(root, target string, a intent.Assignment, active func(intent.AssignmentState) bool) error {
	if !active(a.State) {
		return componentRefusal(componentAssignmentState, a.ID, string(a.State), string(intent.StateActive))
	}
	if a.Worktree != target {
		return componentRefusal(componentAssignmentPath, a.ID, a.Worktree, target)
	}
	// A marker this verb cannot read at all is an owner-marker fault: the later
	// components are read from the same evidence, so nothing downstream is decided.
	evidence, err := validateOwnerMarker(root, target)
	if err != nil || evidence.marker.OwnerID != a.OwnerID || evidence.marker.Path != a.Worktree {
		return componentRefusal(componentOwnerMarker, a.ID, "", "")
	}
	if evidence.registration.BranchRef != a.Branch {
		return componentRefusal(componentRegistration, a.ID, "", "")
	}
	if !evidence.registration.Locked || evidence.registration.LockReason != lockReason(a) {
		return componentRefusal(componentLock, a.ID, "", "")
	}
	return nil
}

func landingActiveState(state intent.AssignmentState) bool { return state == intent.StateActive }

// resumeActiveState accepts the state a first landing leaves behind when it published and
// then failed before release. A resume finishes that landing rather than refusing it.
func resumeActiveState(state intent.AssignmentState) bool {
	return state == intent.StateActive || state == intent.StateCleanupPending
}
