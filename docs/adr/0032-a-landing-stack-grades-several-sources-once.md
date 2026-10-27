# A landing stack grades several sources once

Status: accepted

## Decision

A landing stack lands two to sixteen reviewed sources with one landing command.
Each source is a stack member, and the stack file lists the members in landing order.
The stack holds one landing lease for all its members, and that lease names the assignment of the first member.
The stack grades its top tree with one whole gate, and one update of the destination publishes every green member.

A stack admits only a spec-less member or a tickets-only member.
A member that names a staged spec lands alone, because its landing transitions that spec.
A member whose assignment holds a delivery binding also lands alone, because its landing closes that deliverable.
Each member tip contains the destination tip.
Two members change one path only when that path is in the capture directory, whose settle policy resolves it.
The stack checks each of these rules before composition, so a known cause never costs a gate.

The stack composes the members in file order.
Each member gets its own landing merge commit with its own message.
Each tickets-only member closes its own folder in its own merge commit.
The commitment admission of each member grades the destination with that member alone, so a member answers only for its own change.

The gate proves that a stack tree holds only reviewed bytes.
Each path keeps the destination's bytes or the bytes of the one member that changes it.
A path in the closed tickets folder of a member is absent.
A capture path that several members change passes, because the settle policy owns its bytes.
One owner rule maps each changed path to its members, and the preflight and the gate proof both read it.

Every grade of a stack runs the whole gate and never the reader set.
A green stack moves the green marker to the stack tip.

A red top tree bisects the ordered prefixes of the stack, and each grade is a whole gate.
The stack attributes the red to the first member whose prefix is red on top of a green prefix.
The stack publishes that green prefix, refuses the attributed member, and leaves each later member unlanded for a rerun.
The bisect costs at most one plus the ceiling of log2 of the member count whole gates, so sixteen members cost at most five.
An infrastructure outcome stops the stack with no further grade and no publication.
A red over a destination that held no green also stops the stack with no publication, because no member can carry that repair.

The build window keeps every rule.
Each member passes the window on its own, and a landing override admits only its own member.
Each preflight refusal and each window refusal stops the whole stack before the lease and names its member.

The stack keeps no persistent queue and adds no field to the intent ledger or the review record.
A second session that wants to land meets the held lease, and it runs its own landing after the lease ends.

## Consequence

A set of ready light-path sources pays one whole gate instead of one whole gate for each source.
A partial publication moves the destination, so each unlanded member folds the destination before its rerun.
A flaky red can name a member, and the single landing of that member grades it again.

## Considered options

A landing queue that lasts across sessions stays out.
It needs its own ledger field and its own recovery, and one command with one lease gives most of the saving.
