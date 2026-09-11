# FT311 recoverable reset review

Frozen base: `cf2caa6d5c195c6b6ed5852fb4478902d1d99fbd`

Reviewed tip: `8016bde9ac69a19bc0b8163e0bcb39068f47aba7`

Current findings: Standards 0, Spec 0, Coverage 0.

De-duplicated decision targets: 0.

## Decision

On 2026-09-11 the reviewer decided the ignored-path collision policy.
The reset and the restore refuse before movement when a materialized tracked path would overwrite ignored content.
The refusal names the affected paths and preserves the existing bytes.
Story 53 and rows RR70 and RR71 carry that decision.

## Repairs after the frozen tip

- The plan refuses `ignored content would be overwritten` with the colliding paths in the refusal table, in both modes, before any envelope or move. RR70 and RR71 are the oracles.
- The fingerprint binds the raw index entries, so a staged-only change stales the plan in both modes. RR72 is the oracle. This closes the earlier benchmark's staged-only stale-plan finding for this candidate.
- The reference paragraph names the detached checkout and the drifted registration lock as repaired states, and it names the collision refusal.

## Earlier benchmark findings carried by this candidate

TestResumeReconcileKeepsResetRefsForEveryRecordedState and TestResumeReconcileRefusesAResetRefMovedAfterListing cover the cross-state survival and the moved-ref conditional deletion at the frozen tip.

## Review round of 2026-09-11

Three fresh Sol/high axes reviewed `aed6f18c` over the frozen base: Standards 14, Spec 3, Coverage 9.
The reviewer decided the three behavior targets.
The move cleans before the reset, so an ignore-rule change never deletes bytes.
An index with hidden flags refuses, and a symbolic checkpoint resolves in the target.
The reviewer also extended the fence to the cleanup planner file, so the status argv and the ignored listing have one owner each.

The fold closed every auto-fix target.
Those targets are the command constructor, the single index read, the below-path collision, and the manifest read.
They also include the dropped capture return, the reconcile read order, the stale seam names, the approval line, three comments, and two prose items.
One judgment call stays open as no-op: the long plan function.
The six in-range fence widenings from the benchmark stay as veto surface.

## Scoped re-review of the fold

Three fresh opus/medium axes re-reviewed the fold commit `8016bde9` alone: Standards 3, Spec 2, Coverage 2, none blocking.
Every predicate the fold claimed was verified, and five mutations each turned a new test red.
The fold of this round repaired the status reader's comment, the reference's exit-3 bullet, and ticket 3's `Writes` line.
The Coverage axis found that the restore RR73 names preserves the drifted file in a second envelope and reports it as `preserved`.
That is the RR35 shape, so it stays as decided, and RR73 now follows the named restore and asserts the second envelope holds the bytes.
The three-tag set of the hidden-flag reader is complete for its argv, and the reviewer may pin the `s` tag later.
