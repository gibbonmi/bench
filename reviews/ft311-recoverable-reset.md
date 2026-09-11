# FT311 recoverable reset review

Frozen base: `cf2caa6d5c195c6b6ed5852fb4478902d1d99fbd`

Reviewed tip: `f5667b6ba7b4a9506a0606c45d09a9d15ad04078`

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
