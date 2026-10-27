# Record the background grade rules

Blocked by: 11-let-a-landing-wait-on-a-background-grade.md, 12-stop-background-grades-when-work-goes-away.md, 13-show-background-outcomes.md, 14-retire-the-prepared-run.md
Writes: docs/adr, .bench/BENCH-reference.md, projects/benchkit.md, internal/anchors, CHANGELOG.md, .bench/commitment.json
Covers: BG96, BG97, BG98

## What to build

Write down the current decided state, so a cold session can act without this spec.

- Record the background grade in an ADR. Either write a new ADR or amend the existing ones.
- Each of four ADRs states a rule that FT451 changes:
  - ADR 0002 posture 5 says that gate execution serializes on one lock in each checkout.
  - ADR 0017 says that one landed change pays one gate run, at the landing.
  - ADR 0024 says that no Bench-owned concurrent phase workload remains.
  - ADR 0031 names the consumers that look up a green.
- `.bench/BENCH-reference.md` states the claim, the wait, the background record, and the queue.
- The profile's note "Never mutate the repository while a gate is running" names a foreground gate (`projects/benchkit.md:700-701`). The gate section describes the background grade beside the prepared-run sentence that ticket 14 retires (`projects/benchkit.md:501-504`).

An ADR holds no file path and no code snippet.

## Acceptance

- [ ] An ADR states the background grade, and ADR 0002 posture 5, ADR 0017, ADR 0024, and ADR 0031 agree with it.
- [ ] `.bench/BENCH-reference.md` states the claim, the wait, the background record, and the queue.
- [ ] The profile's note on a running gate names a foreground gate.
- [ ] `bench test --check prose` and `bench test --check docs-currency-workflow` pass.
