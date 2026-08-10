# AXI outcome and action migration

Status: staged

Decision source: `decisions/byte-preserving-axi-foundation/ft173-axi-contract.md`

## Problem

After the shared carriers and registry exist, production commands still return local `(string, int)` pairs or domain-specific action strings. Shared types with no production callers are a degenerate foundation: later contributors can continue bypassing them while all public bytes remain green.

## Solution

After `axi-carriers-and-registry` is implemented, migrate outcome and action construction one independently green family at a time. Each domain keeps its kind, exit, payload, and renderer authority; registry-derived conformance proves every declared member reaches the shared route. A final contraction removes superseded local carrier derivations while the compatibility oracle requires exact public observations.

## User stories

1. As an agent using ordinary query and inspection commands, I want their current outcomes to pass through one typed route without changing output or exits. Line: gpt-5.6-terra / high. The adapter is simple, but complete per-member reachability is an omission-sensitive semantic contract.
2. As an agent using ambient and lifecycle commands, I want status, handoff, dashboard, and spec-build outcome/action facts typed without changing their distinct empty, refusal, or next-step renderings. Line: gpt-5.6-terra / high. Lifecycle actions mix executable commands and orchestration prose and must not be normalized.
3. As an agent using publication, worktree, or shift operations, I want each domain's specialized state and exit policy preserved through shared carriers. Line: gpt-5.6-terra / high. These operational domains have durable state and non-AXI exits that syntactically valid generic mappings can corrupt.

## Implementation decisions

- The `cmd/bench` output adapter migrates first, followed by fence-disjoint query and inspection domain batches. Every ordinary member constructs a shared outcome before its existing renderer; a registry-derived bypass mutation exists per member.
- Status, handoff, dashboard, and spec-build are separate tickets because dashboard consumes status and handoff facts while spec-build has an independent renderer and nine-operation grammar.
- Publication, worktree, and shift remain separate writers. `shift.Outcome` continues to own `complete`, `failed`, `usage`, `incomplete`, `no-op`, `interrupted` and exits 0/1/2/3/4/130.
- Typed actions record fixed tokens and open placeholders but existing `action`, `next`, `next_action`, hint, and prose bytes remain exact. This spec emits no `help[]`.
- The contract ticket removes obsolete local result/action carriers only after every family migration and route junction is green. Zero-consumer exports are residue unless a named public API owner remains.

## Testing decisions

- Each migration attaches to the domain's real command/renderer seam and the pre-existing compatibility oracle, not to mocks of `internal/axi`.
- Registry-derived route observation proves byte-equal bypasses red. Each independently bypassable member receives its own mutation identity at ticket time.
- Existing renderer tests retain exact literals; outcome/action tests mutate domain producer facts, kinds, exits, and action inputs.
- Promotion is the sole full gate after all route and residue rows are classified.

### Seam diagram

    trigger: production command computes a result or next-action fact
        │
        ▼
    domain-owned facts ──▶ [ shared outcome/action carrier ] ──▶ existing renderer bytes + exit
                                  ◀ tests attach here: real producer reachability and compatibility delta

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| OA1 | 1 | The `cmd/bench` adapter and every cmd-local returned-output member construct a shared outcome while preserving exact output and exit. | production `Command` adapter | not TDD-able until the prerequisite carrier exists; current `outputCommand` forwards local `(string, int)` | Per-member route mutations reject an adapter marker that no real result constructs. |
| OA2 | 1 | Learnings, maps, diff, and coverage construct shared outcomes before their existing renderers. | real query command seams | not TDD-able until OA1; current package tests pin only bytes | One bypass mutation per enumerated command rejects a partial query migration. |
| OA3 | 1 | Structure, models, testreport, guards, outline, roadmap, status, handoff, dashboard, and worktree query members construct shared outcomes before rendering. | real inspection command seams | not TDD-able until OA1; current package tests pin only local renderers | Complete registry-derived enumeration rejects every missing member without a sample. |
| OA4 | 2 | Status and handoff route current kinds and executable action facts through shared carriers while dashboard composes those typed facts without changing prose or rows. | status/handoff producer and dashboard consumer | not TDD-able until prerequisite carriers exist; current exact tests are positive controls | Producer-driven tests catch flattening before composition and the oracle catches any byte delta. |
| OA5 | 2 | All nine spec-build operations route outcome and action facts while preserving usage/refusal exits, one-row empty status, full detail empties, and non-invokable orchestration prose. | spec-build grammar, state, and renderer | not TDD-able until prerequisite carriers exist | Enumeration plus executable/prose classification rejects a partial or over-eager `help[]` migration. |
| OA6 | 3 | Publication prepare, submit, promote, rollback, and status route durable outcome/action facts without moving state-machine authority. | publication record and renderer | not TDD-able until prerequisite carriers exist | Bypassing any operation or re-deriving next action outside the record fails route/state tests. |
| OA7 | 3 | Worktree create, list, refresh, cleanup, release, recovery, and resume route current outcome/action facts without changing fingerprints, recovery authority, or empty forms. | worktree lifecycle and renderer | not TDD-able until the adapter and prerequisite carriers exist | Per-mode route mutations and the oracle reject partial adoption or reconstructed recovery advice. |
| OA8 | 3 | All six `shift.Outcome` kinds retain exits 0/1/2/3/4/130, cross the shared route, and render the exact current `shift_result`. | shift outcome and result renderer | observed red: the prerequisite shared-owner probe is absent; current shift tests pin the local map | Exact kind/exit enumeration rejects both bypass and universal AXI normalization. |
| OA9 | 1,2,3 | Per production member declared by the root/nested registries and `ft173-axi-surface-census.md`, every outcome/action route is reached, and superseded local carrier derivations or zero-consumer exports are absent. | registry-derived conformance and residue census | not TDD-able until all migrations land | A restored bypass or legacy derivation makes conformance or the exact literal/symbol census red. |

### Ticket derivation

Each domain row is a complete byte-preserving tracer and may repeat a covers ID per independently green member batch. Each member receives its own closure token and bypass mutation. OA9 is a final junction plus contraction, never a substitute for the migration rows.

| row | tracer acceptance and atomic facts | approved fence | subject mutation | independent owner and public operation |
|---|---|---|---|---|
| OA1 | route cmd-local returned outputs / member / owner kind / exit / exact bytes | `cmd/bench` | restore legacy `(string, int)` construction for one named cmd-local member | command adapter route test plus compatibility case; invoke the real `Command` member |
| OA2 | route learnings / maps / diff / coverage outcomes independently | `internal/learnings`, `internal/maps`, `internal/diff`, `internal/coverage` | bypass shared construction in one named package at a time | per-package route test plus compatibility case; run the public command fixture |
| OA3 | route structure / models / testreport / guards / outline / roadmap / status / handoff / dashboard / worktree query outcomes independently | `internal/structure`, `internal/models`, `internal/testreport`, `internal/guards`, `internal/outline`, `internal/roadmap`, `internal/status`, `internal/handoff`, `internal/dashboard`, `internal/worktree` | bypass one enumerated producer at a time | registry-derived per-member route tests plus compatibility cases; invoke each public query |
| OA4 | preserve status/handoff kinds and executable actions plus dashboard composition / kind / action / prose / composition | `internal/status`, `internal/handoff`, `internal/dashboard` | flatten a fixed action before composition or locally recompute a dashboard fact | real status/handoff producer-to-dashboard tests; execute clean, active, and stale fixtures |
| OA5 | route nine spec-build operations / outcome / exit / action / orchestration prose / empty | `internal/specbuild`, `cmd/bench` | bypass one named operation or mark one orchestration label executable | spec-build grammar/state/renderer tests plus compatibility cases; invoke all nine operations under bounded fixtures |
| OA6 | route five publication operations / durable kind / exit / next action / exact bytes | `internal/publication` | derive next action outside the durable record or bypass one named operation | publication state-machine/renderer tests; run prepare, submit, promote, rollback, and status fixtures |
| OA7 | route worktree create / list / refresh / cleanup / release / recovery / resume with authority unchanged | `internal/worktree`, `cmd/bench` | reconstruct a recovery action from prose or bypass one named mode | worktree lifecycle/renderer tests; invoke each bounded public mode fixture |
| OA8 | route six shift kinds / exact exits / result bytes / bounded detail / interrupted persistence | `internal/shift` | normalize one named exit, omit shared construction, or lose interrupted state on reload | shift outcome and fresh-process tests; drive all six terminal outcomes with named deadlines |
| OA9 | prove every registry/census member route and delete each superseded derivation / route / legacy symbol / zero-consumer export | `internal/conformance`, `internal/axi`, `cmd/bench`, `internal/learnings`, `internal/maps`, `internal/diff`, `internal/coverage`, `internal/structure`, `internal/models`, `internal/testreport`, `internal/guards`, `internal/outline`, `internal/roadmap`, `internal/status`, `internal/handoff`, `internal/dashboard`, `internal/specbuild`, `internal/publication`, `internal/worktree`, `internal/shift`, `projects/benchkit.md` | restore one exact legacy literal/symbol or bypass one enumerated route | conformance and residue tests; enumerate production declarations and exact moved symbols after every migration |

### Edge inventory

- Error path — OA1-OA8 retain current usage, refusal, render failure, stderr, and specialized exits.
- Empty or absent input — OA3-OA7 preserve zero-row, one-row, prose-clean, and absent/refusal forms.
- Boundary values — payload size is unchanged; bounded content stays owned by the later projection spec.
- Malformed input — argv parsing remains unchanged and the compatibility oracle compares accepted/rejected cases.
- Interrupted or partial state — OA6-OA8 retain durable partial publication/worktree/shift state and bounded interrupt tests.
- Re-run idempotency — existing no-op semantics and lifecycle replay remain exact under the oracle.
- Process-boundary lifecycle — publication, worktree, spec-build, and shift cases reload durable state in fresh processes.
- Hostile environment — the oracle retains cwd, PATH, symlink, and control-bearing fixtures.
- Command self-observation — route observation is test-only and cannot alter command output.
- Special files and dangling symlinks — existing domain refusals remain unchanged; no new reader is added.

## Out of scope

- Bounded projection migration — 8 edits, 1 promotion gate run in `axi-bounded-projection-migration`.
- Aggregate and empty migration — 14 edits, 1 promotion gate run in `axi-aggregate-empty-migration`.
- Rendering `help[]`, adding flags/schema, changing family homes, or normalizing operational exits.

## Ownership fences

- Adapter/query batches: `cmd/bench/**`, `internal/learnings/**`, `internal/maps/**`, `internal/diff/**`, `internal/coverage/**`, `internal/structure/**`, `internal/models/**`, `internal/testreport/**`, `internal/guards/**`, `internal/outline/**`, `internal/roadmap/**`.
- Ambient/lifecycle: `internal/status/**`, `internal/handoff/**`, `internal/dashboard/**`, `internal/specbuild/**`.
- Operational domains: `internal/publication/**`, `internal/worktree/**`, `internal/shift/**`.
- Route conformance and contraction: `internal/conformance/**`, `internal/axi/**`, `projects/benchkit.md`.
- The compatibility oracle, bounds owners, aggregate semantics, and public bytes remain unchanged.
