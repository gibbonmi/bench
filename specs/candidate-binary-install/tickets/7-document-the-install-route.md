# Document the install route in the ADR, the spec discipline, and the build phase

Blocked by: 2-install-the-candidate.md, 3-hold-the-refresh.md, 5-refuse-an-uninstalled-writer.md, 6-restore-after-a-stopped-build.md
Writes: docs/adr/0016-the-installed-broker-owns-promotion.md, .agents/skills/bench-craft-spec/SKILL.md, .agents/commands/bench-implement-spec.md, .bench/BENCH-reference.md, projects/benchkit.md, bin/bench.sh, internal/worktree/land.go, internal/adopt/broker.go
Covers: CB44, CB45, CB46, CB47, CB58

## What to build

Deliver CB-C5: the guidance names the install route for a teammate who just walked in.

- ADR 0016 gains the decided exception text from the spec's "The ADR 0016 exception" section, under its Decision. Its Consequence sentence "Candidate code cannot authorize its own publication" gains the qualifier that the install exception is the one route. `projects/benchkit.md`'s broker sentence in the Working branch section and the `land_route` comment in `bin/bench.sh` gain the same qualifier; the `bin/bench.sh` edit changes the comment only. The `LandCommand` doc in `internal/worktree/land.go` and the `WriteBrokerManifest` doc in `internal/adopt/broker.go` gain the same qualifier as comment-only edits; `land.go` has 3 lines of headroom. The text names no file path and no code, and it adds no backticked `docs/adr/` reference elsewhere.
- The edge-inventory section of `craft-spec` requires a spec that declares a persisted field to name each pre-landing writer and the landing binary, and to name the guard of each writer. The skill has 4 lines of budget headroom.
- The Build section of `/bench-implement-spec` names the install from the primary checkout before the first ticket dispatch of a declared build, and the reinstall after a Go change or a fold reaches it. The file sits at its budget of 95 lines, so the author keeps it at 95 lines or fewer by tightening nearby text. Each Require needle in the anchor registry and each canary `old` string that names that file stays byte-identical. If that is not possible, the author stops and asks the reviewer for a budget raise.
- `.bench/BENCH-reference.md` describes the install and the restore in its command notes, and it tells the operator to run the restore at once after `bench commitment block`, because the block ends the hold. Ticket 3 owns the `held` effect entry.

Review owns rule wording, so this ticket adds no anchor row. The new sentences avoid every needle that an existing anchor forbids in these files. The craft-spec sentence stays inside the edge-inventory section, and the build-phase sentences stay inside the Build section.

## Acceptance

- [ ] ADR 0016 states the exception and ends with the sentence that the install is the only route by which candidate code gains landing authority (CB44, review-owned).
- [ ] The `craft-spec` edge-inventory section states the declared-spec rule (CB45, review-owned).
- [ ] The `/bench-implement-spec` Build section names the install before the first dispatch and the reinstall after a Go change or a fold (CB46, CB58, review-owned).
- [ ] `.bench/BENCH-reference.md` describes the install and the restore (CB47, review-owned).

## Checkpoint verification

Run CB-C5's verification commands from the spec's completion plan. The orchestrator then freezes the CB-C5 delta for independent review, and the final reconciliation follows.
