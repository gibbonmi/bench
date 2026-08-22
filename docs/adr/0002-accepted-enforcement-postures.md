# Six enforcement postures are accepted risk, not gaps to fix

Bench's enforcement is deliberately partial. The six postures below are accepted
decisions. Each names the risk it accepts, why the acceptance is right under the
kit's threat model, and what evidence would reopen it. The threat model covers the
lazy or accidental weakening, not a determined
adversary. None of these is mitigated; recording them as accepted is the point.
A record that implies coverage it does not have is worse than silence.

1. **Interactive commit-on-red remains possible.** The gate wrapper and the
   shift loop are the enforced commit paths. Raw git in an interactive session
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
   membership is enforced: the agent-line guard denies delegations off the
   bound tiers. But effort exists only in the declaration, because effort is
   not observable to a hook the way a model id is. Accepted as declaration
   discipline under the no-silent-escalation invariant. Reopens if effort
   becomes mechanically observable, or if effort-declaration drift is shown to
   cause real cost.

4. **The agent-line guard's residual rims fail open.** A delegation without a
   model in a routed repo with complete tier binding is denied; that branch,
   the silent-escalation path, is now closed. The remaining degraded
   branches are an unrouted repo, an incomplete binding, and a malformed or absent
   envelope. They still allow, because there is no binding to enforce, and a broken
   guard must never brick delegation. Accepted as the fail-open rim of a guard
   whose closed branch is the one that mattered. Reopens on evidence that a
   residual rim is exploited as a silent-escalation path.

5. **Gate execution — and through it the gated commit — reuses only a fresh
   green verdict for the identical closed oracle subject.** **Every other
   consumer observes the verdict without granting authority.** The subject binds the
   working tree, resolved oracle, and execution policy, plus every project-declared
   environment, path, and tool input. Gate execution serializes on one lock. A
   subject that still holds a fresh green under that lock is answered from the
   record without a write. So reuse never slides its own freshness window.

   Any
   other state durably replaces the older verdict with a pending record. This
   happens before running that exact subject, and the gate atomically records
   the final ready verdict only if the subject remains unchanged. An open subject can run and report
   diagnostics but cannot authorize reuse. Stale, red, pending, invalid,
   unavailable, or mismatched state pays a real gate run or fails closed.

   One
   residual rides the closure: shellcheck, an optional phase, is deliberately
   undeclared in the subject. Declaring it would open the subject on
   every host that legitimately lacks it, and silently disable reuse there. So
   an in-place shellcheck upgrade inside the freshness window can hide behind
   a reused green. Accepted because re-judging an identical closed subject
   buys no correctness for the gate's full cost. Subject binding and
   durable invalidation prevent an older green from authorizing changed or
   interrupted work. Reopens on evidence that a reused verdict authorized a
   commit the gate would have refused.

6. **Canary coverage is family-level.** The tripwire proves one planted needle
   per check family still bites. It does not plant a needle per individual
   check, so unwiring a single check inside a proven family can escape the
   canary layer. Accepted because per-check needles up front over-fit the
   canary and cost maintenance out of proportion. Family granularity
   catches the always-pass rot the tripwire exists for. Reopens on observed
   anchor rot inside a family that the family-level needle failed to catch —
   that check graduates to its own needle.

## Consequences

A session that finds one of these gaps should cite this record instead of
filing it as a defect or quietly building enforcement for it. Changing any of
these postures is a reviewer decision that supersedes this record. Two
placements were considered and rejected. Per-package code comments were rejected, because a comment
is not a citable decision. The project profile was rejected too, because it describes the repo, not
platform decisions.
