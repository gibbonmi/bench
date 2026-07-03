# The gate runs from the working tree, defended by a canary tripwire rather than write protection

The gate is honored from the working tree by design, so it can evolve in the same diff as the code it grades — which also means the agent being graded can weaken its own oracle. We defend against that with a canary tripwire: the gate runs itself against known-broken fixtures and goes red if any fixture stops failing for its planted reason, or if the fixture set is missing or empty. The defense makes a weakening loud — a red gate and a reviewable diff — rather than impossible; the threat model is the lazy or accidental agent that guts its gate to get unstuck, not a determined adversary.

## Consequences

A freshly scaffolded consumer gate is born defended and red: an un-canaried sentinel keeps it red until the consumer configures real checks, and a canaried example check seeds the fixture set. Between deleting the sentinel and adding real checks, the gate is green on the example alone; that window is the consumer's deliberate, visible configure step and is accepted under the threat model.

## Considered options

Pinning the gate outside the writable tree and hash-verifying it before push would cover the adversarial case — an agent that deletes the tripwire in the same edit that weakens the gate, which nothing inside a writable tree can stop. Pinning is deliberately out of scope: a separate future capability, not an extension of the tripwire.
