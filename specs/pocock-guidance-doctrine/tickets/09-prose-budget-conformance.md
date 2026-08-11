# Enforce the profile's prose-budget table with a fail-closed conformance check

Blocked by: 08-spec-seams-debug-lessons.md
Ownership fence: `projects/benchkit.md`, `internal/conformance`, `tests/canary/guidance-prose-budgets`, `internal/anchors/registry_data.go`
Integration surfaces: registry row→`internal/conformance` (registry.go `Checks` slice, `familyChecks` map, checks_test.go `conformanceChecks` binding, registry_test.go `canaryFixtureFamilyRegistry`, tier_test.go `classifiedLiveTreeTests`); canary family→`tests/canary/guidance-prose-budgets`; budget table and input-bindings advertisement row→`projects/benchkit.md`; profile anchors→`internal/anchors/registry_data.go`
Contracts: the budget table (subject-path and limit rows) crossing `projects/benchkit.md`→`internal/conformance`, asserted by PB1 against the real parser with no second hard-coded budget table
Closure: PB1/table-single-source, PB2a/malformed-policy, PB2b/missing-policy, PB2c/duplicate-row, PB2d/invalid-limit, PB2e/unclassified-new-skill, PB2f/symlink-refused, PB2g/special-file-refused, PB3a/all-paths-one-run, PB3b/at-budget-green, PB3c/one-over-red, PB3d/newline-parity, PB4/canary-bites, PB5/cadence-row-retired

## What to build

Add one mechanically parseable prose-budget table to `projects/benchkit.md`:
`.bench/BENCH.md` 150, `.agents/commands/bench-implement-spec.md` 60,
`.agents/skills/bench-craft-tickets/SKILL.md` 100, and every other real
regular file matching `.agents/skills/*/SKILL.md` 120 — the complete
enumeration universe; other command files are outside the budget and
`.claude/skills/*` adapter symlinks are distribution surfaces, never subjects.
Retire the obsolete "Spec-build guidance cadence" cached routing (its
`promote` contract no longer exists) and reconcile the "Ticket-breakdown
review pass" cached routing with ticket 04's reviewer-approved breakdown. Add
the new check's row to the profile's conformance input-bindings advertisement
table. Implement the checker in `internal/conformance` (mirror
`checkLineBinding` in `line_routing_static_test.go`; consume the pre-staged
`registry.InputBenchkitProfile` input source): parse the table rather than
repeating its numbers, apply the exact craft-tickets exception before the
all-skills default, classify every newly added skill automatically, refuse
malformed/missing/duplicate policy and invalid limits, refuse any symlink or
special file inside the enumerated canonical paths before reading, count
logical lines identically with and without trailing newlines, and report
every over-budget path in one run with path-specific diagnostics. Register
it in the `Checks` slice, `conformanceChecks`, `familyChecks`
(`guidance-prose-budgets`→the new check), `canaryFixtureFamilyRegistry`, and
`classifiedLiveTreeTests` if the harness requires it. PG15 sequencing: this
ticket's worktree is cut from the pre-prune base `2401049` recorded at
breakdown time; register the checker there first and observe its red on
today's genuinely over-budget subjects (BENCH 256, implement-spec 318,
craft-tickets 410, delegate 275, line 147, spec 174, tdd 122), record that
red in the done-report, then rebase onto post-08 main — where the same
checker is green — before landing. TDD the checker at the mechanical seam:
focused tests feed temporary profile-and-guidance trees to the real checker,
one exact-diagnostic case per closure token below, each observed red before
its implementation exists. Add the `tests/canary/guidance-prose-budgets`
fixture whose MUTATE.json pushes a classified guidance file over its limit,
with a fixture-owner mutation test mirroring
`TestWorkflowCadenceAnchorsRejectDeletionAndSwap`. Raising a budget stays a
reviewer edit to the profile.

## Acceptance

- [ ] [PB1] (covers PG13) the profile table is the one budget source; the checker parses it; every enumerated subject is within its limit on the landed tree; no second budget table exists in Go or prose.
- [ ] [PB2] (covers PG14) each fail-closed class reds with its own path-specific diagnostic, one focused temporary-tree test per class: malformed policy, missing policy, duplicate row, invalid numeric limit, newly added unclassified skill (auto-classified to the default, and refused only when the table's default row is absent), symlink inside the canonical universe, special file inside the canonical universe.
- [ ] [PB3] (covers local) a tree with three over-budget subjects reports all three in one run; exactly-at-budget is green; one line over is red; both hold with and without a trailing newline.
- [ ] [PB4] (covers PG15) the pre-prune red on today's subjects is recorded, and the canary fixture's over-budget mutation receives the budget check's own diagnostic via the fixture-owner test; `bench canary` validates the new family.
- [ ] [PB5] (covers local) the cadence cached-routing row is gone and the breakdown-review routing no longer contradicts ticket 04.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PB1/table-single-source | change the craft-tickets limit to 90 in the profile only | budget-check focused test | edit the cell in a temp tree, run the checker, expect it to enforce 90 (parses, not hard-codes) |
| PB2a/malformed-policy | corrupt the table header in a temp profile | budget-check focused test | run the checker, expect the malformed-policy diagnostic naming the profile path |
| PB2b/missing-policy | delete the table from a temp profile | budget-check focused test | run the checker, expect the missing-policy diagnostic |
| PB2c/duplicate-row | duplicate the BENCH.md row | budget-check focused test | run the checker, expect the duplicate-policy diagnostic |
| PB2d/invalid-limit | set a limit to `abc` | budget-check focused test | run the checker, expect the invalid-limit diagnostic |
| PB2e/unclassified-new-skill | add a new skill dir to a temp tree with no default row in the table | budget-check focused test | run the checker, expect the unclassified-subject diagnostic |
| PB2f/symlink-refused | plant a symlinked SKILL.md in a temp tree | budget-check focused test | run the checker, expect the symlink-refusal diagnostic before any read |
| PB2g/special-file-refused | plant a FIFO as a SKILL.md in a temp tree | budget-check focused test | run the checker, expect the special-file refusal without blocking |
| PB3a/all-paths-one-run | make the checker return after its first over-budget hit | multi-violation focused test | run the three-subject case, expect the missing-second-path failure |
| PB3b/at-budget-green | pad a temp subject to exactly its limit | boundary focused test | run the checker, expect no diagnostic |
| PB3c/one-over-red | add one more line to that subject | boundary focused test | run the checker, expect the over-budget diagnostic |
| PB3d/newline-parity | strip the trailing newline at the boundary | boundary focused test | run both PB3b/PB3c cases again, expect identical verdicts |
| PB4/canary-bites | weaken the fixture's MUTATE.json so the subject stays within budget | fixture-owner mutation test | run the owner test, expect its missing-expected-diagnostic failure |
| PB5/cadence-row-retired | restore the cadence row | semantic review reread | reviewer-graded: revived retired contract against the spec |
