# Learnings — usage journal

## 2026-08-15 — a decision map compiled under a spec is deleted with that spec [open]

**What happened.** The deepening map (`/bench-shape-idea` over the 2026-08-15 survey)
was compiled into `specs/gate-run-transaction/decisions/`, and its ticket #5 told
later specs and light-path tickets to cite that path. `bench spec retire
gate-run-transaction` deleted the spec directory whole, so the only in-tree copy of a
map still owning six unbuilt lanes was gone; the next session recovered it from git
history and re-homed it at `decisions/deepening-2026-08.md`.

**Right behavior.** A map that outlives one spec lives in `decisions/`, never under a
spec's own directory; only a map consumed entirely by that one spec may travel with it.

**Proposed rule change.** `/bench-shape-idea` exit: a map whose tickets route to more
than one spec or light-path lane compiles to `decisions/<map>.md`; `/bench-write-spec`
cites, and does not copy, that path.

## 2026-08-15 — write-spec loop 1 needed three rounds on skills-index-reader [open]

**What happened.** Loop 1 on `specs/skills-index-reader` returned 7 blocking findings in
round 1, 2 in round 2, 1 in round 3; the reviewer stopped the loop before a clean
acceptance round and before loop 2. The misses were authoring-stage: the central
single-source promise (SI6) had no red row; `bin/bench.sh` labels and the
`classifiedLiveTreeTests` registry were enforcement surfaces I had not read before locking
rows; the ticket boundary left the SI6 guard red at ticket 1's end; `.bench/BENCH.md`
sat exactly at its prose budget.

**Right behavior.** Before locking rows for a change that adds a CLI verb or a
conformance test in this kit, read the enforcement surfaces that grade those two families
— `checkColdPickupCLILists`' `bin/bench.sh` label derivation, `checkSubcommandRouting`'s
`usage.Parse` rule, `hiddenLiveTreeDiags`' registry, and the prose-budget table — and
walk each ticket's end state against every guard the spec itself adds.

**Proposed rule change.** `projects/benchkit.md` cold-session notes: name those four
enforcement surfaces as the read-before-rows list for "new verb" and "new conformance
test"; `craft-spec` already says to read the enforcement surface — this is the project
half of that rule.
