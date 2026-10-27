# Record the background grade rules

Blocked by: 12-let-a-landing-wait-on-a-background-grade.md, 13-stop-background-grades-when-work-goes-away.md, 14-show-background-outcomes.md, 15-retire-the-prepared-run.md
Writes: docs/adr, .bench/BENCH-reference.md, projects/benchkit.md, internal/anchors, CHANGELOG.md, .bench/commitment.json
Covers: BG96, BG97, BG98, BG104, BG105

## What to build

Write down the current decided state, so a cold session can act without this spec.

- Record the background grade in an ADR. Either write a new ADR or amend the existing ones. An ADR holds no file path and no code snippet.
- ADR 0002 posture 5 says that gate execution serializes on one lock in each checkout, and it describes the prepared chunk run. It then says that a foreground run takes that lock and that a background grade takes none, and it names no prepared run.
- ADR 0017 says that one landed change pays one gate run, at the landing. It then says that a background grade can pay that run before the landing.
- ADR 0024 says that no Bench-owned concurrent phase workload remains. It then names the background grade as that workload, its bound, and the census that admitted it.
- ADR 0031 names the consumers that look up a green. It then names the background grade as a producer whose green each consumer reuses.
- `.bench/BENCH-reference.md` states the claim, the wait, the background record, and the queue.
- The profile's note "Never mutate the repository while a gate is running" names a foreground gate (`projects/benchkit.md:700-701`).
- The profile's gate section retires its prepared-run sentence and describes the background grade (`projects/benchkit.md:501-504`).
- One Forbid anchor keeps each retired prepared-run sentence out of the profile and out of ADR 0002.

## Acceptance

- [ ] ADR 0002 posture 5 says that a background grade takes no checkout lock, and it holds no sentence about a prepared chunk run.
- [ ] ADR 0017 says that a background grade can pay the gate run of a landed change before the landing.
- [ ] ADR 0024 names the background grade, its machine bound, and its census.
- [ ] ADR 0031 names the background grade as a producer of a reusable green.
- [ ] `.bench/BENCH-reference.md` states the claim, the wait, the background record, and the queue.
- [ ] The profile's note on a running gate names a foreground gate.
- [ ] The profile holds no sentence about a prepared chunk run, and that sentence put back alone makes `bench test --check docs-currency-workflow` red.
