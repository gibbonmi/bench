# A commitment binding names its obligations

A deliverable binding in an outcome that owns sources names the outcome
sources that the deliverable completely satisfies. A binding that names none
can record no delivery fact, so it can never close its outcome.

One rule decides when a binding is obligation-free, and two owners apply it.
A plan refuses each new or changed obligation-free binding in every milestone.
The completion landing refuses an obligation-free binding of the active
milestone before the gate and before publication. Its refusal names the plan
step as the repair.

Policy validation does not take this rule. Every policy read validates the
policy, and the plan that repairs a legacy policy must read that policy first.
A plan therefore keeps a legacy binding that it leaves unchanged. An outcome
with no sources has no obligations to name. Each of its deliverables is an
obligation of the outcome itself.

A plan binds only the content that it keeps open. That content is the
unsettled sources and deliverables of the proposed policy, and each unsettled
outcome source of the current policy that the proposal drops. A deliverable is
settled when the policy records a delivery fact for it. A plan that settles or
drops a deliverable does not need that deliverable on the default branch.

The plan identity sorts its bound sources once, in one canonical order. The
order compares the identifier, then the path, then the identity, byte by byte.
One owner sorts the sources and encodes them, so no traversal change can move
the identity. A plan approved before a change to that order must be planned
again.
