# per-component-gate-scoping

Status: implemented

Decision source: the reviewer's parallel-session decision block (five decisions:
invocation model unchanged, per-component skip predicates, build artifact reuse,
prose-triggers-conformance-only, per-component ancestors), pasted and read back
in this session 2026-08-01, plus one in-session reviewer ruling the same day:
the canary's declared input set excludes the bench binary.

## Problem

The gate pays for work the changeset cannot have affected. A `ROADMAP.md` edit
compiles the binary; a prose edit runs the canary. Neither could change either
outcome. Today's reduced run asks one question of the whole changeset — did it
stay inside the capture allowlist — and that single answer switches all seven
excludable phases together, so components whose inputs disagree cannot be
expressed: editing only a canary fixture should run the canary and skip the
build. This is a cost problem, not a correctness problem: the gate grades
correctly today and must continue to. A skipped component must always be one
that something already graded — never one nobody ran.

## Solution

One gate run, same trigger, same phase ordering, one verdict — nothing about
invocation moves. Inside the run, each skippable component carries its own
input declaration and independently answers whether its inputs moved.
Components skipped on evidence inherit from their own content-addressed
ancestor slot with no freshness bound. The build phase skips through artifact
reuse: the freshness seal beside `dist/bench` records a content digest of the
binary's exact input set, and a gate-authored attestation binds the seal to a
build the gate itself ran, so a planted artifact never counts as evidence. A
capture-only commit runs conformance and the conformance suite against a
sealed binary and nothing else.

## User stories

1. As a session committing a capture-only changeset (`capture/`, `specs/`,
   `ROADMAP.md`, `.bench-notes.md`), I pay only the conformance components:
   build skips on its attested seal, every other component skips on its own
   ancestor evidence, and the verdict names each skip with the evidence that
   covers it. Line: opus / high. The inheritance decision is the one seam where
   a wrong answer credits ungraded work.
2. As an ordinary gate run (never `--fresh`), the build phase executes only
   when `dist/bench` is absent, its seal mismatches the current build-input
   set, or the seal is not attested by a prior gate-run build; a valid attested
   seal skips the phase and the binary readers exec the sealed binary. mtime
   plays no part. `bench gate --fresh` executes everything, build included.
   Line: opus / medium. The seal machinery exists; the new work is the in-gate
   skip decision, the attestation, and their fail-closed edges.
3. As the kit maintainer, every component's input declaration is single-sourced
   in the gate package, and the components whose inputs have a derivable
   source — build (`go list -deps ./cmd/bench`), gofmt/vet/test/race/contract
   (the module-wide `go list -deps -test ./...` closure plus listed packages'
   `testdata/` directories), shellcheck (its own argv enumeration) — must take
   their sets from that derivation; only canary's three surfaces and the
   wrapper scripts are hand-declared. The profile's rendered table is
   conformance-bound to the source so drift turns the gate red, and the
   declaration-honesty width is stated in the profile with the same candor as
   the current construction prose. Line: opus / high. Declaration honesty is
   the load-bearing property of the whole feature.
4. As the oracle, a component skipped on evidence resolves its own ancestor
   slot: content-addressed over its declared inputs plus its execution closure,
   no freshness bound, authored only by a run that executed that component
   green, never re-stamped by a skip. A slot record is a distinct
   component-domain class naming its component; any record in a slot that
   carries inherited fields, names the wrong component, or fails validation
   refuses — the component runs. Line: opus / high. Slot soundness is the
   anti-forgery core.
5. As the reviewer, the canary runs exactly when canary-owned surfaces change —
   `internal/canary/` sources, `tests/canary/` fixtures, or the wrapper wiring
   its phase execs (`bin/bench.sh` and the canary dispatch behind it) — and
   skips on ordinary Go edits because its declaration excludes the binary
   (reviewer ruling 2026-08-01, a recorded tripwire narrowing). Composition
   inherits the same narrowing: two changes graded separately may land together
   with the canary never run against the combined tree, and the verdict still
   names every such skip. `bench gate --fresh` and the ship tier remain the
   paths that re-prove the tripwire against a new binary. Line: opus / medium.
   The declaration is data; the recorded-weakening prose must be exact.
6. As an operator reading output, every skipped component is announced on its
   own line with its evidence (ancestor identity and recorded time, or the
   attested seal), the verdict record carries the per-component partition in a
   validated record class, a partial verdict is never reusable as a whole-tree
   green by `bench commit` or anything else, `bench status` renders it as its
   own class, and `bench prep-release` refuses it naming the skipped
   components and pointing at `bench gate --fresh`. Line: opus / medium.
   Record-class validation has strong prior art.
7. As a linked repository or a root with no Go module, I never scope: the
   declarations are the kit's own, matched by directory identity rather than
   spelling. Every skip decision routes through one decision function whose
   every error return — slot unreadable, seal unreadable, derivation failure,
   identity computation failure, domain mismatch — answers run-the-component,
   so fail-closed is a property of one site rather than N. Line: opus /
   medium. The single site is what makes the every-path quantifier checkable.

## Ticket routing note

Story lines above are ceilings. At ticket slicing, these land on sonnet —
each mechanical, fully gate-observable, with direct prior art: the profile
table rendering (story 3's prose half), the shellcheck argv derivation, the
`bench status` partial-class rendering, the prep-release refusal wording, and
the RB1/R20 root-guard test extensions. The decision function and identities,
the Go derivation closures, the slot record class and authorship, the build
attestation, and the whole-tree reuse guard stay at their story lines.

**Ticket conventions (from FT164 — the skill does not teach these yet, so this
build's tickets carry them explicitly).** Blast-radius classification is the
first branch of the breakdown: the decision-function rework touches consumers
in status, commit, and prep-release, so its contract ticket is blocked by
every consumer-migration ticket, each sized by an ownership fence. Every
ticket carries single-line acceptance criteria with requirement identifiers,
an `Ownership fence:` field, and an `Assumptions:` field; fence-disjointness
is the mechanical check that a group may run concurrently. No ticket carries
gate-state checkboxes — cadence is focused seam checks while iterating, one
full gate at the atomic commit, one composed gate at final check. Each
acceptance row binds to one concrete red mutation, its independent owner, and
the exact public operation sequence that proves it — no umbrella claims.
Ticket claims are re-derived from the tree after earlier tickets land, never
from this spec's account of the base. A ticket adding a member to an
enumerated family (the partial record class, the per-component slots, new
profile rows) names every registry found by tracing an existing sibling —
verdict field sets, record validation, the conformance registry, the profile
table, the status renderer. And because this build changes evidence
authorship, every cadence-touching ticket names which command authors
evidence and which consumes it: `bench gate` is the canonical producer;
`gate-run --fresh` prints a valid result without publishing what promotion
consumes.

## Implementation decisions

- **Selection is the gate's own — never the caller's.** The partition is
  computed inside the gate run by the decision function, from the tree's
  content identities against the retained evidence. No caller, flag, agent, or
  session names components to run or skip: `bench commit` keeps asking the
  gate unconditionally and reads the returned verdict alone, and the only
  operator lever remains `bench gate --fresh`, which forces more grading,
  never less.
- **Components and their skip predicates.** `conformance` and
  `conformance-suite` are unconditional: they grade the prose surfaces and
  enforce every binding this feature adds, so the enforcement is never
  skippable by the surfaces it enforces (this is also decision 4's mapping:
  prose and roadmap changes trigger conformance only). `build` skips on the
  attested seal (artifact reuse). `gofmt`, `vet`, `test`, `race`, `contract`,
  `shellcheck`, and `canary` skip on per-component ancestor evidence.
- **Input declarations, derivation mandatory where a source exists.** The
  toolchain and contract components derive from the module-wide
  `go list -deps -test ./...` closure — not the binary's `./cmd/bench`
  closure, which excludes the conformance and contract packages those
  components grade — plus the `testdata/` directories of listed packages.
  Build derives from the binary closure the freshness digest already resolves.
  `contract` additionally includes the seal's source digest — it execs the
  CLI, so the binary is one of its inputs. `shellcheck` derives from the file
  list its argv already enumerates. `canary` is hand-declared:
  `internal/canary/`, `tests/canary/`, and the wrapper scripts its phase
  execs, with the binary digest deliberately excluded. Hand-declared entries
  are part of the conformance-bound table; derivation-sourced components are
  conformance-checked to actually call their derivation, not a copied list.
- **Per-component identity.** Each component's identity hashes its declared
  input contents (positive selection over the same `git add -A` snapshot the
  stripped identity reads today) plus its execution closure (argv shape, env
  contract, and — where the binary is an input — the seal's source digest),
  under a per-component policy domain so no slot can answer for another.
- **`ReducedScope()` is not retired.** It remains the capture-surface
  declaration: the stripped-worktree construction keeps consuming it to prove
  every scoped component blind to `capture/`, `specs/`, `ROADMAP.md`, and
  `.bench-notes.md` on every full run, and `internal/status`'s capture-only
  staleness signal keeps reading it. The new per-component declarations layer
  on top of that floor. Honesty width, stated plainly: the construction proves
  capture-surface blindness only; for the per-component declarations, honesty
  rests on mandatory derivation plus conformance binding, and a component that
  reads an undeclared non-capture path skips wrongly — recorded in the profile
  as the accepted residual, mirroring today's "proves nothing beyond those two
  signatures" prose.
- **Ancestor slots.** N slots in the existing retained-evidence store, one per
  evidence-skipped component. Content addressing with no freshness bound;
  every refusal runs the component. Authorship generalizes to component
  granularity: a run that executed component C green authors C's slot; a run
  that skipped C leaves C's slot byte-identical; a red C invalidates C's slot
  and touches no other component's. The on-disk witness replaces the old
  whole-run marker: a slot record is its own validated class carrying the
  component name and no inherited fields, so "authored at execution" has a
  checkable shape and the forgery refusal (wrong class, wrong component,
  inherited fields present) stays meaningful. (The decision block's "only a
  full run may author" is preserved as intent at component granularity —
  literal whole-run authorship starves every slot once any component scopes,
  because runs are then almost never full.)
- **Build skip mechanics and attestation.** The skip decision requires the
  existing `freshness.Check` contract to pass (missing binary, missing seal,
  digest mismatch, symlinked or irregular sidecar each run the build) AND the
  seal's executable digest to match the gate's own build attestation — a
  gate-evidence record authored only when the build phase runs green inside a
  gate, in the same store as the ancestor slots. A self-consistent seal the
  gate never produced is therefore not evidence: planted binary plus planted
  seal fails attestation and the build runs, overwriting both. A green build
  republishes the seal (existing `Publish`) and re-authors the attestation.
  The `Needs` edges stay in the table and are satisfied trivially when build
  skips — no writer, no race. The gate-entry untrusted-runner check in
  `gate.sh` is a different surface and does not change.
- **Verdict record.** The reduced record class generalizes: the record carries
  the executed set and, per skipped component, the evidence that covered it
  (ancestor identity and recorded time, or the seal's source digest for
  build). Validation follows the existing two-class discipline — exact field
  sets, no spectrum. A partial record is never a reusable whole-tree green:
  the reuse guard that refuses `Reduced` today refuses the partial class
  identically, covered by its own row. `bench status`, `bench commit`, and
  `bench prep-release` consume the class through the existing `Inspection`
  surface.
- **Whole-tree reuse and `--fresh`.** The whole-tree fresh-green reuse
  (60-minute window) still answers first. `--fresh` forces a full run —
  every component, build included — and re-authors every slot and the build
  attestation.
- **Test harness.** The synthetic `reducedRunFixture` root (two phases, no Go
  module, no build phase) cannot exercise build, seal, or canary scoping. The
  build includes a richer fixture: a kit-shaped temp root with a real Go
  module, a sealable binary, and a canary-bearing phase table, alongside the
  existing fixture rather than replacing it. This is test-harness work inside
  this spec's scope; the gate's own canary fixture set is untouched.

## Testing decisions

- The behavior under test is always the external partition: which phases a run
  executed, which it skipped, what evidence each skip named, and what the
  record round-trips.
- Prior art per seam: `reduced_run_test.go` for inheritance decisions and
  authorship (extended with the kit-shaped fixture above),
  `reduced_verdict_test.go` for record-class round-trips and the
  never-reusable guard, `freshness_test.go` for seal edges,
  `scope_binding_test.go` for declaration↔profile binding,
  `runtime_gate_reduced_test.go` for the CLI-observable announcement lines.
- The gate observes the feature through the conformance binding checks (drift
  red, derivation-source red) and through the retained stripped-worktree
  construction, which keeps proving capture-surface blindness on every full
  run.

### Seam diagram

    trigger: bench gate / bench commit (unchanged invocation)
        │
        ▼
    changeset tree ──▶ [ per-component decision (one function, one       ] ──▶ phases to run
    + attested seal      fail-closed site): identity vs slot;            │     + skip lines
    + N ancestor slots   build: freshness.Check + attestation            │     + verdict record
                              │                                          ▼
                        [ runner (unchanged ordering) ] ──▶ green: author slots + attestation
                                                            for executed components only
                        ◀ tests attach here: kit-shaped fixture root + fault-engine clock;
                          assert executed set, skip announcements, slot bytes, record class

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | a capture-only changeset executes only conformance components (build skipped on attested seal; gofmt, vet, test, race, contract, shellcheck, canary skipped on evidence — set enumerated) | kit-shaped fixture: executed-set assertion | observed red planned: red before the partition exists | a global-allowlist implementation runs build; a stub partition runs nothing — both fail the exact executed-set assertion |
| 1 | the verdict names every skipped component with its evidence, one announcement line and one record entry per skip | record round-trip + CLI output contract | observed red planned | a silent skip reads as a gate that never ran |
| 2 | valid attested seal skips build; readers exec the sealed binary unchanged | kit-shaped fixture with gate-built sealed binary; byte-compare | observed red planned | an mtime or always-rebuild implementation rewrites the binary; byte-compare catches it |
| 2 | each of: absent binary, absent seal, source-digest mismatch, executable-digest mismatch, symlinked/irregular sidecar, missing attestation, attestation/seal mismatch — runs the build (seven cases, each its own subtest) | seal + attestation fixtures through the gate decision | observed red planned per case | fail-open on any one case execs a stale, unknown, or planted binary |
| 2 | a self-consistent seal the gate never authored (planted binary + recomputed seal) fails attestation and build runs | planted-artifact fixture | observed red planned | this is the case the unauthenticated seal cannot catch; attestation exists exactly for it |
| 3 | each of the seven scoped components' declarations is single-sourced and the profile table matches it | scope-binding conformance check extended per component | already covered pattern; new table rows get observed red by desync mutation | a hand-edited profile row that drifts from the source reds the gate, not review |
| 3 | build, gofmt, vet, test, race, contract, and shellcheck take their input sets from their named derivations (enumerated); a literal path list in any of them reds conformance | derivation-source conformance check | observed red planned | a hand-copied list passes a table check and drifts from what the component grades |
| 3 | the toolchain/contract derivation covers the module-wide `./...` closure including test files: breaking a test file in a package outside the binary closure moves their identities and they run | kit-shaped fixture with a package outside `cmd`'s closure | observed red planned | the binary-closure derivation skips test/vet/gofmt on a tree whose test suite does not build — the delegate-demonstrated hole |
| 4 | a skipped component's slot is byte-identical after the run | slot-bytes assertion (prior art: reduced-run re-author test) | already covered for the stripped slot; per-component red planned | re-stamping dresses old evidence as new |
| 4 | a slot record is the component-domain class: inherited fields present, wrong component name, or failed validation each refuse and the component runs | slot-forgery fixtures, one per refusal | observed red planned | without the class witness, "authored at execution" has no on-disk meaning and forged evidence inherits |
| 4 | a red component invalidates its own slot and no other | fixture forcing one component red | observed red planned | invalidating everything re-pays the full gate; invalidating nothing reuses evidence the red contradicted |
| 5 | editing only a canary fixture runs canary and skips build | kit-shaped fixture, edit under `tests/canary/` | observed red planned | the decisions' worked example; a global predicate cannot express it |
| 5 | editing the wrapper script the canary phase execs runs canary | edit `bin/bench.sh` in fixture | observed red planned | pins the hand-declared wrapper surface so the cheapest declaration (two directories only) fails |
| 5 | an ordinary Go edit skips canary (binary excluded) and runs the toolchain components | fixture edit under a Go source | observed red planned | pins the ruled narrowing so a later refactor cannot silently widen canary's inputs back |
| 6 | a partial verdict is never reusable as a whole-tree green (`bench commit` reuse guard refuses it) | reuse-guard test (prior art: reduced not-reusable row) | already covered pattern; partial-class red planned | a missed reuse-guard edit turns a partial green into a release-path green for 60 minutes |
| 6 | prep-release refuses a partial verdict naming the skipped components and pointing at `--fresh` | prep-release refusal test | already covered pattern; wording red planned | a release resting on skipped components ships ungraded work |
| 6 | `bench status` renders a partial verdict as its own class | status fixture | observed red planned | a partial verdict rendered as full green hides narrowness from the board |
| 7 | a non-kit root and a no-Go root never scope | RB1-pattern tests extended | already covered for the global guard; per-component red planned | a linked repo inheriting against kit declarations grades against a scope it never chose |
| 7 | every error class in the decision function (slot unreadable, seal unreadable, derivation failure, identity failure, domain mismatch — enumerated) answers run-the-component | error-injection subtests at the single decision site | observed red planned per class | one fail-open error path anywhere silently credits ungraded work |
| edge of 1 | first run, pruned evidence, and `--fresh` each execute everything (build included) and author all slots plus the attestation | R20/--fresh prior-art tests extended | already covered pattern; build-included red planned | absence treated as reusable emits evidence-free skips; `--fresh` skipping build makes the tripwire escape a no-op |
| edge of 4 | slots and the attestation authored by one process are honored, and forged ones refused, by a fresh CLI process reloading them | runtime contract test execing the CLI twice against one evidence store | observed red planned | unit-level green hides serialization defects that appear only on reload — the FT164-flagged process-boundary class, and this build is all serialized evidence |

### Edge inventory

- Error path — every decision-site failure: story 7 row (enumerated classes).
- Empty/absent input — no binary, no seal, no attestation, no slots: story 2
  and edge-of-1 rows; an unchanged tree is answered first by whole-tree reuse.
- Boundary — a path in two declarations (`internal/canary` source is both a
  build input and a canary input): both components run when it moves; identity
  is per-component, no arbitration.
- Malformed input — forged slot records, planted seal, symlinked sidecar:
  story 2 and 4 rows.
- Interrupted state — a run killed between executing a component and authoring
  its slot leaves the old slot: safe, the next run re-executes; the inverse of
  the byte-identical assertion.
- Re-run idempotency — two consecutive identical runs: the second skips
  everything skippable and authors nothing; slot-bytes rows.
- Hostile environment — non-kit root, no-Go root, `BENCH_KIT` elsewhere:
  story 7 row.
- Composition — changes graded in separate partial runs landing together:
  sound for every component whose declaration is honest (each component's
  inputs were graded at exactly the composed content); for canary it inherits
  the ruled narrowing — the combined tree may land with canary never run
  against it, visible in the verdict's skip lines (story 5 prose, story 6
  rows). Accepted by the same ruling, backstopped by `--fresh` and ship tier.
- **Won't handle:** a component that reads an undeclared non-capture path
  skips wrongly — honesty for the per-component declarations rests on
  mandatory derivation plus conformance binding, not on construction proof;
  recorded in the profile as the accepted residual.
- **Won't handle:** adversarial substitution beyond the attestation (an
  attacker who can also write gate evidence owns the store outright) — the
  kit's guard contract is loudness against honest mistakes, not defense
  against a writer inside `.git`.
- **Won't handle:** mtime anywhere in any predicate — banned by decision.

## Out of scope

- **Check-level conformance scoping** (running only the conformance checks
  whose observed paths moved): separate capability over the conformance
  registry, ~14 edits, ~6 gate runs.
- **Linked-repo scope declarations** (a repo declaring its own component
  inputs): joins FT144's two-audience work, ~10 edits, ~4 gate runs.
- **Focused canary invocation** (`bench canary` running one fixture/family):
  already owned by FT168, unchanged here.
- **Per-component stripped-worktree constructions** (proving each declaration
  by running its component against a complement-stripped tree): a genuine
  enforcement capability with its own design (the complement of a derived
  input set is most of the tree), ~12 edits, ~5 gate runs; the capture-surface
  construction stays as the floor meanwhile.
