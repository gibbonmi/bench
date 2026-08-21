# Diagnose abbreviated commit identities at landing comparisons

Blocked by: enrich-refusals-through-one-emitter.md
Writes: internal/worktree/land.go, internal/worktree/land_test.go
Line: opus / medium — literal comparisons at one known seam.

## What to build

At every literal identity comparison in the landing family, a case-insensitive
hex prefix of length 4–39 of the compared identity classifies as abbreviated:
the refusal names the flag, states the abbreviation, and prints the full
identity as `wanted` through the emitter ticket's typed refusal. The value is
never accepted. A full-length value that differs takes the drift shape with
both identities (`observed` and `wanted`) and no abbreviation claim, as does a
prefix shorter than 4 or non-hex input. For `--base` this adds a refusal the
tree does not have today (a short `--base` currently resolves and lands); the
spec flags that posture decision and sign-off settles it before this ticket
starts.

## Acceptance

- [ ] A land with --source-tip shortened to 12 hex and every other value correct exits 1 with a refusal naming the abbreviation and printing the full worktree HEAD (covers LR1).
- [ ] A land with --base abbreviated and a full correct --source-tip exits 1 with a refusal that names --base as abbreviated and contains no tip-mismatch detail (covers LR2).
- [ ] A resume with a short --resume value exits 1 with the resolved full identity in the refusal (covers LR3).
- [ ] A land with a full-length --source-tip that differs from the worktree HEAD exits 1 with both identities in the refusal and no abbreviation claim (covers LR4).
- [ ] After a land with an abbreviated --source-tip, destination HEAD, source branch tip, and green marker are unchanged (covers LR5).
- [ ] An uppercase abbreviated identity classifies the same as lowercase (edge under LR1).
