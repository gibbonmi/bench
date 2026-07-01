# tdd acceptance coverage

## Problem

Bench's feature-build path says TDD happens at pre-agreed seams, but it does not
force the spec to name the observable behaviors that must be protected by tests.
That leaves room for an agent to write a few seam tests, get them green, and still
miss parts of the promised story breadth.

## Solution

Add acceptance coverage maps to the feature-build workflow. `/bench-write-spec`
writes the map as part of Testing decisions, `craft-tdd` defines valid rows and
red signals, `/bench-implement-spec` closes every mapped behavior during the
build, and `/bench-review-implementation` treats uncovered mapped behavior as a
Spec finding. A small gate anchor check keeps the content contract from silently
drifting, while dogfooding proves the workflow actually works.

## User stories

1. As a reviewer, I want every non-trivial feature spec to map user stories to
   concrete acceptance coverage rows, so the implementation target is fixed before
   code is written.
2. As a reviewer, I want each coverage row to name `story`, `behavior`, `seam`,
   `red signal`, and `why it catches the failure`, so I can reject vague or
   decorative tests in one pass.
3. As an implementing agent, I want `craft-tdd` to define valid acceptance rows and
   red signals, so I know what counts as TDD coverage and what must be rejected.
4. As an implementing agent, I want feature red signals to be observed red before
   implementation, so passing tests prove a missing behavior was actually closed.
5. As an implementing agent, I want rows that are already covered or not TDD-able
   to be explicitly classified, so the coverage map stays truthful instead of
   pretending every behavior started red.
6. As an implementing agent, I want `/bench-implement-spec` to make me name the
   row each vertical slice turns red-to-green, so story breadth stays visible while
   I work.
7. As a reviewer, I want `/bench-implement-spec` to end with a compact coverage
   table before the gate, so I can see every mapped behavior as `green`, `already
   covered`, or `not TDD-able`.
8. As a reviewer, I want `/bench-review-implementation` to include acceptance-map
   omissions in the Spec axis, so semantic review catches a clean implementation
   that missed promised behavior.
9. As a Bench maintainer, I want the gate to assert the acceptance-coverage anchors
   exist in the workflow content, so the new contract cannot silently disappear.
10. As a Bench maintainer, I want this change dogfooded through its own spec and a
   small real follow-up planning/build path, so adoption is based on behavior, not
   attractive prose.

## Implementation decisions

- **Spec owns the map.** Extend `/bench-write-spec` so Testing decisions include
  an acceptance coverage map for every non-trivial feature. The map is the source
  of truth for TDD rows; implementation and review consume it.
- **Approval table includes coverage.** The `/bench-write-spec` approval surface
  expands from user stories / seams / out of scope to user stories / seams /
  acceptance coverage / out of scope, so the reviewer can veto both breadth and
  test intent before implementation.
- **`craft-tdd` defines validity.** Add one compact section to `craft-tdd` defining
  the five fields, feature red-signal semantics, `already covered` and `not
  TDD-able` classifications, and the rejection bar for shallow/internal tests.
  Keep existing seam discipline; do not import the whole bug-diagnosis workflow.
- **`/bench-implement-spec` closes the map.** Add build instructions that each
  vertical slice names the row it is closing and the phase ends with the compact
  coverage table before the final gate.
- **Review audits omissions.** Update `/bench-review-implementation` so the Spec
  axis treats missing, partial, or falsely-classified acceptance rows as findings.
  The review remains advisory; it does not replace the gate.
- **Gate anchors the content contract.** Add a small kit-conformance check that
  the affected workflow files retain the acceptance-coverage anchors: the
  five-field map in `/bench-write-spec`, row/red-signal validity in `craft-tdd`,
  closeout statuses in `/bench-implement-spec`, and acceptance-map audit language
  in `/bench-review-implementation`. The check is structural, not a semantic
  parser.
- **Dogfood is required before adoption.** After the content edits, run a small
  real Bench planning/build path with the changed kit and report whether it
  completed. If it cannot be completed, the implementation remains proposed, not
  adopted.

## Testing decisions

- **Good tests here** exercise the kit content surface as future agents see it:
  commands and skills contain the instructions they must follow, the gate catches
  missing anchors, and a dogfood run proves the workflow is usable.
- **Seams:** kit content surface (`/bench-write-spec`, `craft-tdd`,
  `/bench-implement-spec`, `/bench-review-implementation`), project gate
  (`bench gate`), and the dogfood workflow (`/bench-write-spec` -> implement/shift
  -> `/bench-review-implementation` -> `bench gate`).
- **Prior art:** `.bench/gate.sh` already has kit-conformance checks for command
  and skill content drift. Add this as another narrow content-contract check.
- **Gate command:** `bench gate`.

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1, 2 | `/bench-write-spec` requires the acceptance coverage map and approval summary for non-trivial feature specs. | Kit content surface: `/bench-write-spec` | Observed red before implementation: `rg -n "acceptance coverage map\|why it catches the failure\|red signal" .agents/commands/bench-write-spec.md` exits 1. | The command currently has no acceptance-map language, so this fails until the spec template and approval table require the row shape. |
| 3, 4, 5 | `craft-tdd` defines valid rows, observed feature red signals, truthful `already covered` / `not TDD-able` classifications, and shallow-test rejection criteria. | Kit content surface: `craft-tdd` | Observed red before implementation: `rg -n "acceptance row\|red signal\|not TDD-able\|call count\|internal test double" .agents/skills/bench-craft-tdd/SKILL.md` exits 1. | The skill currently has seam-level TDD guidance but not acceptance-row validity, so this catches the missing test-quality contract. |
| 6, 7 | `/bench-implement-spec` makes each vertical slice name the row it closes and emits a final coverage table. | Kit content surface: `/bench-implement-spec` | Observed red before implementation: `rg -n "coverage table\|already covered\|not TDD-able\|turning red-to-green" .agents/commands/bench-implement-spec.md` exits 1. | The command currently says vertical slices and gate, but not row-by-row closeout, so this catches the accounting gap. |
| 8 | `/bench-review-implementation` audits missing, partial, or falsely-classified acceptance rows in the Spec axis. | Kit content surface: `/bench-review-implementation` | Observed red before implementation: `rg -n "acceptance coverage\|coverage row\|mapped behavior" .agents/commands/bench-review-implementation.md` exits 1. | The review command currently checks spec requirements generally but does not name acceptance-map omissions, so this catches review drift. |
| 9 | `bench gate` fails if the acceptance-coverage content anchors disappear. | Project gate: kit conformance | Observed red before implementation: `rg -n "acceptance coverage map\|not TDD-able\|why it catches the failure" .bench/gate.sh` exits 1. | The current gate cannot see this content contract, so this catches the missing structural guard. |
| 10 | The changed workflow is dogfooded on a small real follow-up planning/build path before adoption. | Dogfood workflow | Not TDD-able before implementation: this proof requires the changed kit to exist first. | The proof is behavioral rather than static; it prevents adopting a prose change that reads well but fails in use. |

## Out of scope

- Rewriting existing historical specs or decision maps to use acceptance coverage
  rows. They record past workflow state; migrating them is a separate cleanup, ~45-60
  minutes.
- Building a semantic parser that proves every future spec's coverage map is
  complete. This is a different conformance product and would need false-positive
  design, ~1-2 hours.
- Changing `/bench-debug`. The bug path already has a tight repro/regression-test
  discipline; altering it is a separate workflow decision, ~30-45 minutes.
- Automating a full dogfood shift inside `bench gate`. That would make the gate
  slow and environment-sensitive; a future harness benchmark could explore it,
  ~1-2 hours.
