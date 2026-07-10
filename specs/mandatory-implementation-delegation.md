# Mandatory implementation delegation

Status: staged

## Problem

Invoking `/bench-implement-spec` does not itself guarantee that a spec-backed build
uses subagents. The phase currently permits inline work whenever it judges delegation
overhead too high, so the reviewer must repeat "use subagents" to get the intended
orchestration. That makes a preferred implementation venue depend on conversational
memory instead of the phase contract.

The kit also needs an honest failure posture. Some harnesses cannot spawn write
subagents, and Codex does not expose spawning on a deny-capable hook. A portable
phase can require delegation and the gate can prevent that requirement from
disappearing, but the kit cannot truthfully claim a cross-harness runtime guard.

## Solution

Make meaningful write delegation mandatory for every spec-backed implementation.
Independent vertical slices fan out in parallel within the harness limit; dependent
slices stay sequential; an atomic build goes wholly to one isolated write subagent.
Read-only helpers do not count, while a no-spec lighter-path change may remain inline.

When the current harness cannot spawn a write subagent, stop before the first edit and
give the reviewer explicit resume instructions: the repository location, working
branch or worktree, spec name, destination harness, and the exact invocation that
harness accepts. Pin each clause at the existing root-conformance guidance seam with
distinct diagnostics and a canary fixture, then dogfood the changed phase in a fresh
session. The result is portable generation guidance with a permanent bite proof, not
a fictitious spawn ledger.

## User stories

1. As a reviewer starting a spec-backed implementation, I want the phase to delegate
   meaningful write work before any implementation edit, so that I no longer have to
   append an instruction to use subagents. Independent slices use separate parallel
   subagents, dependent slices remain sequential, and an atomic spec goes wholly to
   one isolated write subagent. Line: `gpt-5.6-sol` / high. Command-phase guidance
   steers every future build, so the project leverage override routes it to the top
   tier even though the prose edit is small.

2. As a reviewer using the implementation phase, I want the delegation requirement
   to distinguish real implementation from ceremony, so that a read-only helper does
   not satisfy it and a tiny no-spec lighter-path edit can still remain inline. Line:
   `gpt-5.6-sol` / high. This is part of the same high-leverage phase contract and a
   subtle scope error would multiply through every later implementation.

3. As a reviewer on a harness without write-subagent support, I want the phase to
   stop before editing and tell me exactly where and how to resume in a capable
   harness, so that it neither silently falls back inline nor leaves me to translate
   a dead command key. Line: `gpt-5.6-sol` / high. The unavailable-venue handoff is
   user-facing command guidance whose correctness is semantic rather than fully
   gate-observable.

4. As a kit maintainer, I want the mandatory delegation clauses guarded by root
   conformance with distinct failure messages and a permanent canary, so that future
   edits cannot quietly restore optional or ceremonial delegation. Line:
   `gpt-5.6-terra` / medium. The project caches gate and conformance work on the mid
   tier because oracle correctness matters and the existing workflow-anchor seam is
   known.

## Implementation decisions

The canonical implementation phase owns **whether** a build delegates. The existing
delegation and line skills continue to own **how** a delegate is charged, routed,
isolated, and verified; this change points to those owners and does not duplicate
their rules.

The phase classifies its venue before the first implementation edit:

- A run with an approved spec must use at least one write subagent. One independent
  vertical slice is one coherent charge; independent charges run in parallel within
  the harness's concurrency limit, while dependency edges make them sequential.
- When the spec is one atomic diff, the entire implementation is one charge to a
  worktree-isolated write subagent. The invoking session still verifies the returned
  done-claim through the existing delegation contract.
- Research, review, planning, or other read-only helpers do not satisfy the write
  requirement.
- A no-spec change admitted by the existing lighter-path threshold may remain inline.
  This is the only inline exception introduced or retained by this slice.

An incapable harness has a fail-closed phase posture: it stops before editing and
emits a resume handoff containing the repository path, branch or worktree, spec name,
and one applicable destination instruction. Codex is named as
`$bench-implement-spec`, Claude Code as `/bench-implement-spec`, and another
subagent-capable AGENTS harness as the canonical implementation command file. The
agent selects the real destination and presents one executable route, not a blind
menu.

Root conformance extends the existing workflow-guidance anchor checker rather than
adding a second parser. It guards four enumerated clause classes with distinct
diagnostics: mandatory spec-backed write delegation before editing; independent,
dependent, and atomic routing; the read-only exclusion plus lighter-path exception;
and the pre-edit stop plus explicit harness-native resume route. Missing guidance
fails closed. One targeted workflow-guidance canary removes the mandatory contract
and expects its diagnostic, proving the family bites; the existing fixture registry
and bite inventory remain the single inventories for that canary.

The phase command, conformance anchors, canary, and spec status change land as one
atomic build. Stories 1–3 therefore collapse to the highest approved line during
implementation. After the gate is green, `craft-synthesis` requires a fresh-session
dogfood run of a small spec-backed build: observe a genuine write subagent start
before edits, verify its charged rows, and confirm the invoking session independently
checks the done-claim. This is the observable ceiling; no hook or gate result is
described as proof that a historical spawn occurred.

## Testing decisions

- A good automated test drives root conformance against the canonical command text
  or a planted broken fixture and observes a distinct diagnostic. It does not infer
  delegation from unrelated prose or duplicate the command contract in an adapter.
- The existing workflow-guidance anchor checker is the automated seam. Its canary
  family provides the break-it proof through the real conformance phase and gate.
- The actual spawn behavior is tested at the fresh-session phase-invocation seam.
  This is a required dogfood observation, not TDD-able gate coverage, because no
  portable event ledger exists.
- Gate command: `.bench/gate.sh`.

### Seam diagram

Automated guidance-contract seam:

    trigger: root conformance grades the kit or a planted canary tree
        │
        ▼
    command guidance  ──▶  [ workflow-guidance anchor checker ]  ──▶  diagnostics
    broken fixture    ──▶  [                                  ]  ──▶  gate exit 0/1
                              ◀ tests attach here: require each collapsed contract
                                clause; canary removes one and expects its message

Fresh-session behavior seam:

    trigger: reviewer invokes implement-spec with an approved spec
        │
        ▼
    spec + harness capability  ──▶  [ canonical implementation phase ]  ──▶  write subagent(s)
    atomic/dependency shape    ──▶  [                                ]  ──▶  verified done-claim
                                      ◀ tests attach here: fresh-session dogfood observes
                                        spawn-before-edit and independent verification

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | Every spec-backed run assigns genuine write implementation to at least one subagent before the first edit. | workflow-guidance anchor checker | Observed red: `rg -qi 'spec-backed.*write subagent'` against the current canonical implementation phase exits 1; during TDD the new root-conformance anchor must red with its own missing-mandatory-delegation diagnostic before the guidance edit. | Optional delegation and read-only-only orchestration both lack the required write-before-edit clause, so the anchor rejects the cheapest wrong wording. |
| 1 | Routing enumerates all three shapes: independent slices use separate parallel subagents within the harness limit; dependent slices run sequentially; one atomic spec is delegated whole to one isolated write subagent. | workflow-guidance anchor checker | Observed red: `rg -qi 'independent.*parallel'` against the current phase exits 1; the new collapsed routing anchor must then fail before its guidance lands. | A phase that merely spawns one ceremonial helper, serializes independent work, or inlines an atomic build omits at least one enumerated shape and fails the anchor. |
| 2 | Read-only helpers do not satisfy mandatory implementation delegation, while a no-spec lighter-path change may remain inline. | workflow-guidance anchor checker | Observed red: `rg -qi 'read-only.*does not'` against the current phase exits 1; the new scope anchor must fail before the command carries both clauses. | Requiring a write subagent and naming the sole inline exception prevents both degenerate implementations: ceremonial research and mandatory delegation for tiny no-spec work. |
| 3 | A harness without write-subagent support stops before editing and gives one explicit resume route containing repository, branch or worktree, spec, destination harness, and that harness's exact invocation. | workflow-guidance anchor checker | Observed red: `rg -qi 'stop before editing.*subagent-capable'` against the current phase exits 1; the unavailable-venue anchor must fail before the explicit handoff text lands. | An inline fallback, a post-edit stop, or a bare canonical slash command cannot satisfy the complete pre-edit resume contract. |
| 4 | Four failure classes have distinct conformance diagnostics: mandatory write delegation; routing shapes; scope and non-counting; unavailable-harness resume. Missing or empty canonical guidance fails closed. | root conformance on the real tree and planted fixture | Red-first implementation: add the four anchor checks before editing the command and run `go test -count=1 ./internal/conformance -run '^TestRootConformance$'`; it must fail with the targeted diagnostics. | Distinct checks make each missing contract attributable, while the absent-file path already fails rather than skipping the required canonical surface. |
| 4 | The workflow-guidance canary permanently proves the mandatory-delegation check bites through the conformance phase. | canary fixture registry and fixture-bite inventory | Red-first implementation: register a planted command missing the mandatory clause and run `go test -count=1 ./internal/conformance -run '^TestRootConformance$'`; without the matching EXPECT/bite wiring, the fixture contract fails. | A check deleted or weakened until the planted command passes makes the canary go red, preventing a vacuous guidance guard. |
| edge of 1 | A fresh-session spec-backed dogfood run starts genuine write delegation before edits and the invoking session verifies the charged rows and one independent behavior. | real implementation-phase invocation in a linked or isolated consumer | Not TDD-able: the portable gate has no historical spawn ledger; record the fresh-session invocation, delegate charge, returned done-claim, and independent verification in the synthesis handoff. | The manual observation distinguishes guidance that merely reads correctly from guidance a real harness follows, without pretending a grep proves runtime behavior. |

**Degenerate-implementation check.** The cheapest wrong phase says "use a subagent"
but permits a read-only helper, keeps all edits inline, and offers an inline fallback
when spawning is unavailable. Rows 1–4 independently reject those shortcuts. The
cheapest wrong gate adds prose without conformance or adds a conformance anchor with
no planted break; rows 5–6 reject those shortcuts. The dogfood row catches a
structurally conformant command that a real fresh session still ignores.

### Edge inventory

- **Error path — subagent facility unavailable:** covered by story 3; stop before
  editing and emit one complete, harness-native resume route.
- **Empty/absent input — no approved spec:** covered by story 2; only a change that
  independently qualifies for the existing lighter path may remain inline. Missing or
  malformed required spec state continues to use the phase's existing entry refusal.
- **Boundary values — one slice, several independent slices, dependent slices:**
  covered by story 1's enumerated routing row. One slice is the atomic whole-build
  charge; parallelism never exceeds the harness limit.
- **Malformed intent — read-only helper offered as compliance:** covered by story 2;
  it does not count.
- **Interrupted or partial delegate state:** **Won't handle** — existing delegation
  isolation, worktree, and stopped-build contracts own interrupted writers and
  partial green work; this slice changes venue selection only.
- **Re-run idempotency:** covered by story 3 for the unavailable route because it
  stops before editing; re-entry in a capable harness starts from the named unchanged
  branch or worktree. Re-running after an interrupted writer remains with the
  existing delegation contract.
- **Hostile environment — missing spawn capability:** covered by story 3. Missing
  Git, model bindings, or build tools remain existing phase and line failures.
- **Every shipped harness surface:** covered through the canonical command consumed
  by thin adapters. The resume contract explicitly names Codex, Claude Code, and the
  canonical-file route for other AGENTS harnesses; no adapter receives copied policy.
- **Missing or empty canonical command file:** covered by root conformance's existing
  missing-file posture plus the four new anchors.
- **Paths with spaces/globs and cwd below root:** **Won't handle as new cases** — no
  path parser is added; delegate worktree paths remain governed by the existing
  path-pinning rule.
- **Control bytes, missing trailing newline, unquoted arguments, symlink invocation,
  and SIGINT in a CLI loop:** **Won't handle** — this slice adds no CLI, wire format,
  shell argument, symlink, or loop surface.

## Out of scope

- **Shared structured conversational output** — the second reviewer-approved slice
  is a separate cross-phase communication capability with its own shared-rule seam and
  spec. Estimated later cost: `~5 edits, 2 gate runs`.
- **Hard runtime spawn policing or a spawn ledger** — a separate harness-integration
  capability, blocked until every supported harness exposes a deny-capable or
  auditable spawn event. Estimated later cost after that prerequisite exists:
  `~7 edits, 3 gate runs`.
- **Mandatory subagents for no-spec lighter-path changes** — reviewer-directed scope
  cut; those edits remain inline-capable. Estimated reversal cost: `1 command edit,
  1 gate run`.
- **Changing model-tier bindings or delegation mechanics** — separate configuration
  and orchestration capabilities already owned by the line and delegation contracts;
  this slice only changes when they are invoked. Estimated later cost if requirements
  change: `~3 edits, 2 gate runs`.
