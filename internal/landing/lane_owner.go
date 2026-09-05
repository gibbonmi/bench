package landing

import (
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/gate/authorization"
)

// NewLane returns the worktree-commit owner whose authority is the root's declared fast
// lane rather than the whole-project gate.
func NewLane(lane authorization.LaneAuthority) Owner {
	owner := New()
	owner.authorize = lane.Authorize
	return owner
}

// laneAuthority maps a resolved lane and the base it is measured against onto the
// authority that grades a composed tree. It is the one place the lane's fields become an
// authority, so no caller states the mapping a second time.
func laneAuthority(lane *gate.Lane, base string) authorization.LaneAuthority {
	return authorization.LaneAuthority{
		Checks:    lane.Checks,
		Kit:       lane.Kit,
		Selective: lane.Selective,
		Base:      base,
	}
}

// NewForLane returns the owner a caller lands under given the lane its root declares. A
// nil lane is a root that declares none, so it answers the whole-project gate owner.
func NewForLane(lane *gate.Lane, base string) Owner {
	if lane == nil {
		return New()
	}
	return NewLane(laneAuthority(lane, base))
}
