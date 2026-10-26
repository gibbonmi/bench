# Gate reuse under one grade key (FT314)

Status: ready

## Destination

One grade key names each gate verdict. A later run of an unchanged composition looks up that verdict and runs no second whole gate.
The complete checkpoint runs under the landing broker, so its green answers the landing of the same tree.
The FT314-gate-reuse criteria in the commitment ledger fix the key members, and they require a measurement of the saving before any code change.
This map reconciles the staged draft spec on this branch with those criteria.

## Notes

Domain: the gate subject, the prospective evidence store, and the checkpoint, merge, and landing verbs.
Use `bench-craft-domain` for each ticket. Use `bench-craft-research` for tickets #2 and #7.

The reviewer signed off a draft spec at commit `9b9150893` before this map existed.
Ticket #12 sorts its rows into kept, changed, and retired rows.

These decisions stay closed. Cross-run reuse is approved (2026-10-14). One key serves all runs, the capability posture binds every reuse, and the reader-table identity is path-free (2026-10-26).
The cost disclosure form is `bench worktree merge --plan` (2026-10-26).

Do not cite the FT314 roadmap row as a source, because the landing that closes the row deletes it.

Term: **grade key** — the digest that names one gate verdict by its verdict inputs. It frames no run purpose.
Avoid: authorization key, evidence key, reuse key. The authorization of a verb is the grade lookup plus that verb's own checks.
The build adds this term to `CONTEXT.md` when the key ships.

A map-owned asset stays in the map's assets folder, specs/gate-reuse/decisions/gate-reuse/assets/.

## Decisions so far

- [Grade key content](gate-reuse/tickets/1.md): verdict inputs only; each verb runs its purpose and evidence checks on every call.
- [Merge plan disclosure](gate-reuse/tickets/8.md): parked as its own intake; FT314 ships no plan form.
- [Light-path stack owner](gate-reuse/tickets/9.md): FT450 owns several sources under one gate.
- [Evidence paths](gate-reuse/tickets/4.md): only the chunk checkpoint's review record, as FT445 strips it today.
- [Schedule member](gate-reuse/tickets/10.md): the baseline phase manifest bytes; the store holds whole grades only.
- [Lookup verbs](gate-reuse/tickets/5.md): every whole-gate verb looks up a whole grade before it picks a reduced schedule.
- [Executable roles](gate-reuse/tickets/2.md): a worktree checkpoint can run under a candidate build; the landing runs under the broker.
- [Runner members](gate-reuse/tickets/3.md): the run binary digest and the running executable digest replace the baseline runner identity.
- [Measurement design](gate-reuse/tickets/6.md): a git upper bound over spec landings and retries since 2026-10-01; ticket #11 decides.
- [Measured saving](gate-reuse/tickets/7.md): 53 minutes per week at most; 14 minutes per week while checkpoints run under worktree builds.
- [Build scope](gate-reuse/tickets/11.md): full scope, and the complete checkpoint runs under the authenticated landing broker.
- [Kept draft rows](gate-reuse/tickets/12.md): 27 rows kept, six changed, and the merge plan rows retired.
- [Broker seam](gate-reuse/tickets/13.md): in the kit, the complete checkpoint execs the manifest's broker through `internal/brokermanifest`.
- [Run binary member](gate-reuse/tickets/14.md): the seal's source digest of the run binary's source root, computed before the lookup.

## Not yet specified

## Spec-writer discretion

## Out of scope

- The background grade of each lane-passed tip. FT451 owns it, and it consumes the key that this map decides.
- The `bench worktree merge --plan` disclosure. It is its own intake (ticket #8).
- Several light-path sources landed in order under one gate. FT450 owns it (ticket #9).

## Sources

- Path: `specs/gate-reuse/decisions/gate-reuse/assets/executable-roles.md`
  Supports: tickets #2 and #3, the running executable, run binary, and identity root of each run.
  Drift: a change to the wrapper selection, the landing route, or the run binary source.
- Path: `specs/gate-reuse/decisions/gate-reuse/assets/saving-measurement.md`
  Supports: tickets #7 and #11, the measured saving of the grade key.
  Drift: a change to the destination-moved refusal, the wrapper selection, or the final-check order.
- Path: `.bench/commitment.json`
  Supports: the three FT314-gate-reuse criteria that this map reconciles.
  Drift: a commitment plan that changes those criteria.
- Path: `docs/adr/0016-the-installed-broker-owns-promotion.md`
  Supports: ticket #3, the guard that keys green evidence to the baseline runner identity.
  Drift: an amendment of ADR 16.
- Path: `internal/gate/prospective.go`
  Supports: tickets #2 and #3, the run binary selection and the baseline runner identity.
  Drift: a change to `prospectiveRunBinaryOwnerAt` or `baselineRunnerIdentity`.
- Path: `internal/gate/checkpoint.go`
  Supports: ticket #1. A chunk green answers for its subject alone because the evidence pass runs on every call, and a complete checkpoint frames its purpose.
  Drift: a change to `applyCheckpoint`.
- Path: `specs/gate-reuse/decisions/gate-reuse/assets/draft-rows.md`
  Supports: ticket #12, the kept, changed, and retired draft rows.
  Drift: a change to a decision of this map.
- Path: `internal/commitment/repository/window.go`
  Supports: ticket #13, the digest comparison of the running executable with the broker manifest.
  Drift: a change to `checkCandidate`.
