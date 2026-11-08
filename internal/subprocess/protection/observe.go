// Package protection owns the observations that decide whether a process group can still
// use a resource. It is a dependency leaf: subprocess and resource owners import it.
package protection

import (
	"errors"
	"fmt"
	"syscall"
)

// GroupState is one current observation of a process group.
type GroupState string

const (
	// GroupAbsent means a current probe answered ESRCH.
	GroupAbsent GroupState = "absent"
	// GroupPresent means a current probe reached at least one member.
	GroupPresent GroupState = "present"
	// GroupUncertain means the probe failed for a reason other than ESRCH.
	GroupUncertain GroupState = "uncertain"
)

// Observation is the state of one group probe and the error that made it uncertain.
type Observation struct {
	State GroupState
	Err   error
}

// Probe sends a signal to a process or, with a negative target, to a process group.
// ObserveGroup sends plain signal zero through syscall.Kill, and a protocol's default probe
// is systemProbe.
type Probe func(target int, signal syscall.Signal) error

// ObserveGroup probes process group pgid with signal zero.
func ObserveGroup(pgid int) Observation { return observeGroup(syscall.Kill, pgid) }

// observeGroup sends signal zero to group pgid through probe and grades the answer.
func observeGroup(probe Probe, pgid int) Observation {
	if pgid <= 1 {
		// Zero and one select the caller's own group and every process, not one group.
		return Observation{State: GroupUncertain, Err: fmt.Errorf("process group %d is not one group", pgid)}
	}
	return Classify(probe(-pgid, 0))
}

// Classify grades the answer to one signal sent to a group. Only ESRCH proves absence:
// EPERM and every other error leave a member possibly alive, so they stay uncertain.
func Classify(err error) Observation {
	switch {
	case err == nil:
		return Observation{State: GroupPresent}
	case errors.Is(err, syscall.ESRCH):
		return Observation{State: GroupAbsent}
	default:
		return Observation{State: GroupUncertain, Err: err}
	}
}
