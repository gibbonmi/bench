# Recovery-aware ref inventory

Status: ready

## Destination

Grow the unclaimed-ref sweep into one repository-wide ref inventory that classifies each Bench-namespace branch by its content against `main`.
The sweep discards only a ref whose commits stay reachable after the discard.
A ref with content that `main` lacks leaves the repository only through an explicit reviewer discard, and that discard writes a dated discarded ref first.
`bench status` never names a destructive command.
The roadmap owner is FT199.

## Notes

The domain terms are in `CONTEXT.md`: unclaimed ref, landed ref, subsumed ref, holder, unique ref, and discarded ref.
The phase used craft-domain, craft-grill, and the read-only probes of 2026-09-27 that ticket 1 records.
The reviewer answered every round on 2026-09-27 and confirmed each recommendation.
Standing preferences: worktree retirement goes through `bench worktree clean`; no raw Git route exists; the reviewer owns every destructive choice.
A map-owned asset stays in the map's assets folder, specs/ref-inventory/decisions/ref-inventory/assets/.

## Decisions so far

- [Which classes does the inventory give a ref?](ref-inventory/tickets/1.md): four classes; the bulk sweep discards landed and subsumed rows only.
- [How does a ref become one kept on purpose?](ref-inventory/tickets/2.md): no marker; the unique class is the keep, and the citing roadmap row owns the reason.
- [What is a discarded ref?](ref-inventory/tickets/3.md): a dated ref at the exact tip, written before the delete, swept 30 days after its date.
- [How does the inventory handle a superseded ref?](ref-inventory/tickets/4.md): it classifies as unique, and only an explicit reviewer discard removes it.
- [What does the status git row route to while unclaimed refs exist?](ref-inventory/tickets/5.md): the plan-only unclaimed form, always.
- [How does the reviewer discard one unique unrecorded ref?](ref-inventory/tickets/6.md): `--target` resolves the branch, and the apply needs the set fingerprint.
- [How does the classifier pick the holder of a subsumed ref?](ref-inventory/tickets/7.md): the lexically first equal tip is the root, and the row names the holder.
- [Does `--apply-current` stay in the unclaimed grammar?](ref-inventory/tickets/8.md): it stays, discards landed and subsumed rows only, and leaves the status action table.
- [What does `bench spec retire` list?](ref-inventory/tickets/9.md): recorded rows whose label or request token contains the slug, plus the unique count; it discards nothing.
- [One spec or two?](ref-inventory/tickets/10.md): one spec with two chunks, the classifier chunk before the discard chunk.
- [Does this map include the disposition of today's 42 refs?](ref-inventory/tickets/11.md): no; fixtures prove the classes, and the live rows stay with their roadmap owners.

## Not yet specified

## Spec-writer discretion

- The row schema of the plan output: the column names for class and holder, and the row order.
- The seam that keeps the discarded namespace out of the lifecycle-namespace emptying rule while the 30-day rule applies.
- The fixture repositories that prove each class, including the equal-tip case and the squash-fold case.
- The exact spelling of the retire listing lines.

## Out of scope

- The disposition of the seven cited tips: FT346, FT347, and FT348 own it.
- The discard of the three FT336 ticket drafts: the reviewer runs the explicit route after the landing.
- A hold marker, a keep list, or a hold verb: the unique class is the keep.
- An automatic supersession proof from a spec retirement: a human claim is not a discard proof.
- A raw Git route for any retirement: closed by FT345.
- A change to `clean --landed` or to the landing prune: both already use the four proofs.

## Sources

- Path: `roadmap/FT199.md`
  Supports: the destination, the occurrences, and the closed decisions on the retirement routes.
  Drift: the row retires or its Next line changes.
- Path: `docs/research/parallel-implementation-wave.md`
  Supports: the recovery table that cites four of the live tips, and the FT348 disposition that holds them.
  Drift: a stream lands, is re-authored, or is dropped.
- Path: `roadmap/FT345.md`
  Supports: the closed decision that `clean --discard-branch <path>` is the unlanded release route.
  Drift: the row reopens a retirement surface.
- Path: `specs/ref-inventory/decisions/ref-inventory/tickets/1.md`
  Supports: the probe record of 2026-09-27 over the 43 live refs.
  Drift: a ref lands, moves, or is discarded.
