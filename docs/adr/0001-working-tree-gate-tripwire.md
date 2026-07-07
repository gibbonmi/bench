# The gate runs from the working tree, defended by a canary tripwire and a push-time pin

The gate is honored from the working tree by design, so it can evolve in the same diff as the code it grades — which also means the agent being graded can weaken its own oracle. Two layers defend that, both loud rather than impossible:

- **The canary tripwire**, inside the tree: the gate runs itself against known-broken fixtures and goes red if any fixture stops failing for its planted reason, or if the fixture set is missing or empty. This catches the lazy or accidental agent that guts a check to get unstuck.
- **The push-time pin**, outside the writable tree: the managed pre-push hook verifies the committed gate tree against a human-recorded local pin and blocks drift with re-pin instructions. This covers what the tripwire cannot — a weakening that deletes the tripwire in the same edit that guts the gate. When no pin is recorded, the drift check is disarmed and the hook's own manifest says so.

The threat model for both layers is the lazy or accidental weakening, not a determined adversary: an actor with shell access can still edit local git hooks, and defending against that remains out of scope.

## Consequences

A freshly scaffolded consumer gate is born defended and red: an un-canaried sentinel keeps it red until the consumer configures real checks, and a canaried example check seeds the fixture set. Between deleting the sentinel and adding real checks, the gate is green on the example alone; that window is the consumer's deliberate, visible configure step and is accepted under the threat model.
