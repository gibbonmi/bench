# Six enforcement postures are accepted risk, not gaps to fix

Bench's enforcement is deliberately partial. The six postures below are accepted
decisions: each names the risk it accepts, why the acceptance is right under the
kit's threat model — the lazy or accidental weakening, not a determined
adversary — and what evidence would reopen it. None of these is mitigated;
recording them as accepted is the point, because a record that implies coverage
it does not have is worse than silence.

1. **Interactive commit-on-red remains possible.** The gate wrapper and the
   shift loop are the enforced commit paths; raw git in an interactive session
   is not fenced, so a session can still commit on a red gate. Accepted because
   fencing every raw git invocation is out of proportion to the threat, and the
   repo-stays-green invariant governs the paths that matter. Reopens on
   evidence that interactive commit-on-red actually lands red work with any
   frequency.

2. **Done-claims made outside a shift are honor-system.** The stop hook
   verifies a shift's end against the gate; an interactive session's "done" is
   discipline, not enforcement — hooks can only grade what they observe.
   Accepted because the gate-is-the-oracle invariant governs behavior and the
   shift path exists for work that needs the enforced loop. Reopens on evidence
   of unverified done-claims shipping defects off the shift path.

3. **The declared line's effort level has no enforcement surface.** Model
   membership is enforced — the agent-line guard denies delegations off the
   bound tiers — but effort exists only in the declaration, because effort is
   not observable to a hook the way a model id is. Accepted as declaration
   discipline under the no-silent-escalation invariant. Reopens if effort
   becomes mechanically observable, or if effort-declaration drift is shown to
   cause real cost.

4. **The agent-line guard's residual rims fail open.** In a routed repo with a
   complete tier binding, a delegation without a model is denied — that branch
   was the silent-escalation path and is now closed. The remaining degraded
   branches — an unrouted repo, an incomplete binding, a malformed or absent
   envelope — still allow, because there is no binding to enforce and a broken
   guard must never brick delegation. Accepted as the fail-open rim of a guard
   whose closed branch is the one that mattered. Reopens on evidence that a
   residual rim is exploited as a silent-escalation path.

5. **The gated commit reuses a fresh green verdict; every other consumer of
   the verdict cache treats it as advisory.** The cache is keyed to the content
   hash of the exact tree the gate judged and has a single writer — only a
   finished gate run records, and an unverifiable tree hash records nothing —
   so a fresh green verdict proves precisely the tree the commit is about to
   land, which the commit's own block-check has already pinned to the named
   set. Anything less — stale, red, untrusted, or absent — pays a real gate
   run. Accepted because re-judging a byte-identical tree buys no correctness
   for the gate's full cost, and the exact-tree key, fresh-only rule, and
   single writer together close the lie-a-tree-green risk that made the
   earlier always-recompute posture the safe default. The commit contract
   suite regression-tests both directions: a fresh green verdict commits
   without re-running the gate, and a verdict recorded for any other tree
   forces a re-run. Reopens on evidence that a reused verdict authorized a
   commit the gate would have refused.

6. **Canary coverage is family-level.** The tripwire proves one planted needle
   per check family still bites; it does not plant a needle per individual
   check, so unwiring a single check inside a proven family can escape the
   canary layer. Accepted because per-check needles up front over-fit the
   canary and cost maintenance out of proportion, while family granularity
   catches the always-pass rot the tripwire exists for. Reopens on observed
   anchor rot inside a family that the family-level needle failed to catch —
   that check graduates to its own needle.

## Consequences

A session that finds one of these gaps should cite this record instead of
filing it as a defect or quietly building enforcement for it; changing any of
these postures is a reviewer decision that supersedes this record. Two
placements were considered and rejected: per-package code comments (a comment
is not a citable decision) and the project profile (it describes the repo, not
platform decisions).
