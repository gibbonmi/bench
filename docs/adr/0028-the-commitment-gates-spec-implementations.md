# The commitment gates spec implementations

Status: accepted

The commitment holds the approved order of outcomes, and each outcome pays a plan, an approval, a start, and a delivery landing.
That cost suits a spec implementation, but it is larger than the value of a one-ticket change.
So the commitment gates spec implementations only, and a one-ticket light-path change needs no outcome, plan, or start.

A light-path change carries one tickets-only folder that holds exactly one ticket.
The change lands with that folder as its spec, and that landing closes the folder.
Its production paths stay inside the expected writes of that ticket.
The exemption applies while a spec outcome is active, so a light-path change never waits for the commitment.

## Consequences

- The policy of the commitment still changes only through a plan that the reviewer approves.
- A pinned roadmap row still changes only through the outcome that pins it. A light-path change can remove a row that no outcome pins.
- A tickets-only folder that a milestone approves as a deliverable is committed work. It keeps the binding rule and records its delivery fact.
- A change that spans two tickets, or that writes a path outside its ticket, keeps the binding refusal.
- A bound assignment keeps its current admission, whatever light-path folders its tree holds.
- The drain dispatches each light-path fix that its verdicts keep to a fresh write delegate. The delegate runs on the mid tier at high effort, in its own bench worktree. The dispatch does not wait for the active commitment.

## Considered options

- Admit each light-path item into the commitment before it ships: rejected, because each one-ticket change then pays the full cost of an outcome.
- A light-path landing with no spec: rejected, because that landing leaves the ticket folder open.
