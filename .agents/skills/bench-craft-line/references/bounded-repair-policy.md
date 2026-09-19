# Bounded repair policy

## Scope and allowance

Each implementation chunk permits at most two repair cycles after its initial review.
This allowance applies to retained, full, delegated, unattended, and light-path implementation runs.
Light-path work counts as one chunk when it receives review findings.
This policy adds no mandatory initial review for light-path work.
Initial implementation, pre-review checks, and the first review consume no repair cycles.

A repair cycle is one repair attempt and verification of its affected findings. It can address several findings.
Avoid: tool call, individual finding, fresh review alone.
Individual tool calls and unchanged verification reruns consume no additional repair cycles.

The allowance takes precedence over continuation while progress holds after initial review.
Progress does not extend the allowance. A fresh review does not reset the same chunk's count.
An uncapped implementation line does not remove this allowance.
Only an explicit reviewer decision extends the allowance; record the additional work it permits.
The existing no-progress, cancellation, user-budget, required-decision, and external-blocker stops can stop work earlier.

## Classification and completion

Before accepting a blocker, cite its binding requirement or concrete defect evidence.
The following unresolved conditions remain blocking, including at exhaustion:

- A required check fails.
- An approved acceptance requirement fails.
- A concrete correctness defect remains, even with a green gate.
- A concrete safety defect remains, even with a green gate.
- A documented mandatory standard fails, even without an automated check.

A preference without a binding requirement or concrete defect remains optional advice.
Optional advice is neither a finding nor a repair target.
Give optional advice no disposition or finding ID, and exclude it from finding totals.
Retain optional advice in a separate advice section of the native excerpt and the review pickup.

After the acceptance rows of a chunk prove, the chunk permits at most one hardening cycle.
A hardening cycle adds checks beyond the approved acceptance rows.
A later finding that grades only a check from that hardening cycle is optional advice, unless it meets a blocking condition above.

Use the existing review phase's dispositions for findings; preserve earlier findings and their supersession history.
Do not label an unrefuted suggestion as `no-op` to obtain a pass.
A native review can pass with no finding IDs while its prose retains optional advice.

A review-record prose correction is evidence-only only when it changes no finding, source identity, observation, disposition, or verification claim.
Evidence-only corrections consume no repair cycle. Batch all cited corrections before verification.
Only the issuing axis reaffirms an evidence-only correction unless it invalidates another axis's evidence.

After substantive repairs, obtain current results or permitted native reaffirmations from every review axis.
Required verification, acceptance reconciliation, and completion checkpoints remain mandatory.
Unresolved findings, stale source identity, missing axes, and incomplete verification still block completion.
Repeat independent review only for a later semantic change or a cross-chunk concern that invalidates earlier evidence.
The coordinator commits the final review record before starting the whole-project gate. It does not edit the gate subject while the gate runs.

## State and handoff

Retain the consumed count and any explicit reviewer extension in the existing review pickup.
For light-path work without a pickup, retain that state in the existing session handoff.
Use ordinary prose; this policy adds no counter schema or CLI command.
Before a resumed session starts another repair cycle, read the retained state.
If the count is missing, reconstruct it from available evidence; never assume zero.
If evidence cannot establish the remaining allowance, return that uncertainty to the reviewer before another repair cycle.

At exhaustion with blockers, report the consumed allowance, completed repairs, remaining blockers, and the needed reviewer decision.
While that decision is pending, stop further repair cycles and dependent chunks.
When blockers close within the allowance, continue through the existing verification and landing requirements.
The allowance is a maximum, not a target.
