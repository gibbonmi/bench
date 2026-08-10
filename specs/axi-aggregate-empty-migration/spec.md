# AXI aggregate and empty migration

Status: staged

Decision source: `decisions/byte-preserving-axi-foundation/ft173-axi-contract.md`

## Problem

Bench command domains independently order and type aggregate facts and independently encode successful absence. A generic renderer could derive totals from visible rows, coerce unknown to zero, omit required zero classes, or normalize zero-row tables, one-row lifecycle empties, prose-clean states, and refusals.

## Solution

After bounded projections are implemented, migrate each aggregate owner separately to shared ordered facts and make every registry member's empty disposition explicit. Domain producers retain total, completeness, skip, fingerprint, record, and absence authority. Registry-derived conformance proves route reachability; contraction removes superseded aggregate/empty carriers only after exact bytes remain green.

## User stories

1. As an agent reading scan, outline, roadmap, or worktree metadata, I want totals, unknowns, zeros, skips, and ordering preserved through shared aggregate mechanics. Line: gpt-5.6-terra / high. Each family has counterexamples a generic `len(rows)` implementation would flatten.
2. As an agent reading lifecycle and ambient output, I want spec-build, publication, status, and dashboard facts typed without losing required zeros, unknowns, or durable record authority. Line: gpt-5.6-terra / high. These outputs are syntactically similar but derive meaning from different state owners.
3. As an agent receiving an empty result, I want each command to retain its exact successful-empty or refusal class, so absence never becomes zero, unknown, or failure by inference. Line: gpt-5.6-terra / high. Empty normalization can remain byte-plausible while changing semantics.

## Implementation decisions

- Guard scan, outline, roadmap, and worktree migrate as separate ownership-fenced tickets. Each supplies ordered typed facts to the shared aggregate without moving its total/completeness/fingerprint derivation.
- Spec-build, publication, and status/dashboard migrate separately. Reclamation keeps every zero-valued disposition class; publication record and status signals remain the semantic owners.
- Empty is a registry declaration plus a typed outcome fact, not one renderer. Zero-row TOON, spec-build one-row `state=empty` with full detail empties, prose-clean status, and absent/unreadable refusal remain distinct.
- Every empty-capable member declares and reaches one exact class; members with no semantic empty explicitly declare absence. Metadata cannot invent empty behavior.
- Final contraction removes only superseded local aggregate/empty carrier mechanics and retains domain derivations and existing renderers.

## Testing decisions

- TDD attaches at real producer-to-renderer seams with counterexamples: unknown guard total, outline skipped input, roadmap source total, worktree fingerprint inventory, reclamation zero classes, and publication/status durable facts.
- Empty tests drive real command fixtures through existing renderers and the compatibility oracle; no universal empty fixture stands in for multiple families.
- Per-family route mutations make byte-equal bypasses red. Residue tests enumerate removed literals/symbols and reject zero-consumer exports.

### Seam diagram

    trigger: domain computes ordered facts or an empty/absent outcome
        │
        ▼
    semantic owner ──▶ [ aggregate + empty carriers ] ──▶ existing exact renderer
                            ◀ tests attach here: real producer counterexamples and class reachability

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| AE1 | 1 | Guard scan retains complete/incomplete, inspected, total/unknown, omitted, and reason/none in exact order through shared aggregates. | guard scanner and renderer | already covered by `TestScanEnumerationTimeoutUsesUnknownCounts`; migration adds route observation | Unknown-to-zero, incomplete-to-complete, and field-order mutations remain independent. |
| AE2 | 1 | Outline retains tracked, scanned, skipped, total, emitted, omitted, and truncated owner facts through shared aggregates. | outline result and renderer | already covered by `TestCommandBoundsRowsAndFullRetainsMetadata`; migration adds route observation | A visible-row total or dropped skip fails without relying on projection internals. |
| AE3 | 1 | Roadmap context/drain retains ordered owner totals, unknowns, zeros, degraded evidence, and emitted-less-than-total meaning. | roadmap context/drain renderer | existing context tests are positive controls; shared reachability is absent | Counterexamples reject row-count derivation and evidence omission. |
| AE4 | 1 | Worktree ignored inventory retains owner count/bytes/shown/truncated and fingerprint authority through ordered aggregates. | worktree inventory and renderer | existing inventory ownership tests are positive controls; shared reachability is absent | Mutating the carrier to infer totals cannot reproduce the authority inventory. |
| AE5 | 2 | Spec-build abandonment/reclamation retains ordered counts and every required zero class through shared aggregates. | spec-build lifecycle record and renderer | already covered by `TestReclaimReceiptReportsEveryClassIncludingTheEmptyOnes`; migration adds route observation | Omitting zero classes or deriving from visible rows turns the real renderer red. |
| AE6 | 2 | Publication records retain ordered typed facts and owner-derived next-state aggregates without moving durable-record authority. | publication record and renderer | current renderer tests pin bytes but no shared route exists | A local bypass or record-independent derivation fails route/state tests. |
| AE7 | 2 | Status and dashboard retain ordered signal summaries, zero/unknown values, and composed section meaning through shared aggregates. | status producer and dashboard consumer | current exact tests pin bytes but no shared route exists | Producer/consumer tests reject fixture-only agreement and local recomputation. |
| AE8 | 3 | Zero-row TOON, spec-build one-row empty with full detail empties, prose-clean status, and absent/unreadable refusal remain distinct and exact. | typed empty outcome plus real renderers | already covered by `TestTable`, `TestSpecBuildStatusRendersDefinitiveEmptyProjection`, `TestStatusHasDefinitiveEmptyAndActiveProjections`, and `TestRenderClean`; migration adds route observation | Universal normalization breaks at least one independent public contract. |
| AE9 | 3 | Per root/nested production-registry member, every empty-capable member declares and reaches its exact class and every other member carries an exact no-semantic-empty disposition. | registry empty disposition and per-member reachability | current empty meanings are command-local | Missing, defaulted, or bypassed classification fails complete registry enumeration. |
| AE10 | 1,2,3 | Per named aggregate owner in `ft173-axi-surface-census.md` and per empty disposition in the production registry, every route is reached and no superseded local carrier or zero-consumer export remains. | registry-derived conformance and residue census | not TDD-able until AE1-AE9 land | Per-family bypass and exact literal/symbol residue mutations reject partial contraction. |

### Ticket derivation

Each semantic owner row becomes an independently green migration with direct producer facts and a shared-route mutation. AE8 may repeat its covers ID across independently green renderer families. AE9 is per production member; AE10 is the final junction/contraction. No staged behavior may be relabeled local.

| row | tracer acceptance and atomic facts | approved fence | subject mutation | independent owner and public operation |
|---|---|---|---|---|
| AE1 | preserve guard complete/incomplete / inspected / total-unknown / omitted / timeout-none / order / route | `internal/guards` | coerce unknown to zero, incomplete to complete, omit reason, reorder, or bypass | guard timeout/enumeration/route tests; run complete and bounded incomplete scans |
| AE2 | preserve outline tracked / scanned / skipped / total / emitted / omitted / truncated / order / route | `internal/outline` | derive total from visible rows, drop skip, reorder, or bypass | outline metadata/route tests; run bounded and skipped-input fixtures |
| AE3 | preserve roadmap context/drain order / totals / unknown / zero / degraded evidence / emitted-less-than-total / route | `internal/roadmap` | derive from rows, coerce unknown, omit zero/degraded evidence, or bypass | roadmap context/drain/route tests; run complete and degraded fixtures |
| AE4 | preserve worktree count / bytes / shown / truncated / fingerprint authority / order / route | `internal/worktree` | infer total from shown, move fingerprint derivation, reorder, or bypass | worktree inventory/route tests; run below/above visible limits and hostile inventory |
| AE5 | preserve spec-build abandonment/reclamation order / counts / every zero class / route | `internal/specbuild`, `cmd/bench` | omit one named zero class, derive from visible rows, reorder, or bypass | lifecycle record/renderer/route tests; invoke abandon and reclaim fixtures |
| AE6 | preserve publication record order / typed facts / durable next state / zero-unknown / route | `internal/publication` | derive outside record, stringify/omit a fact, reorder, or bypass | publication record/renderer/route tests; reload and render durable records |
| AE7 | preserve status signal order / zero-unknown / dashboard composition / route | `internal/status`, `internal/dashboard` | locally recompute a dashboard fact, coerce unknown, reorder, or bypass | real status producer-to-dashboard tests; execute clean, active, and degraded fixtures |
| AE8 | preserve zero-row TOON / spec-build one-row empty plus detail empties / prose clean / absent-unreadable refusal / route | `internal/toon`, `internal/specbuild`, `cmd/bench`, `internal/status`, `internal/dashboard` | normalize one named empty producer to another class | real renderer plus compatibility tests; execute each distinct empty/absent fixture |
| AE9 | classify and reach empty per root/nested member / exact class / explicit no-empty / route | `cmd/bench`, `internal/conformance`, `internal/axi`, `projects/benchkit.md` | remove one member class, default an empty class, or bypass one capable member | registry-derived empty conformance; enumerate and drive each production member or exact not-applicable disposition |
| AE10 | prove every surface-census aggregate and registry empty route and delete legacy carriers / owner / route / residue | `internal/conformance`, `internal/axi`, `internal/guards`, `internal/outline`, `internal/roadmap`, `internal/worktree`, `internal/specbuild`, `cmd/bench`, `internal/publication`, `internal/status`, `internal/dashboard`, `internal/toon`, `projects/benchkit.md` | restore one exact legacy symbol/literal or bypass one enumerated owner | conformance and exact residue census; enumerate source owners and production dispositions after all migrations |

### Edge inventory

- Error path — AE1, AE3, AE4, and AE8 keep incomplete, degraded, unsafe, absent, and unreadable distinctions.
- Empty or absent input — AE8 and AE9 exhaust successful empty, unknown, and refusal classes.
- Boundary values — AE1-AE7 cover zero, one, many, unknown, emitted-less-than-total, and required zero classes.
- Malformed input — existing domain parsers/refusals remain unchanged; registry validation rejects missing empty classification.
- Interrupted or partial state — guard timeout, degraded roadmap, partial worktree inventory, and durable lifecycle records retain owner facts.
- Re-run idempotency — repeated rendering preserves identical order/types and no route accumulates facts.
- Process-boundary lifecycle — spec-build and publication facts reload from durable records in fresh processes.
- Hostile environment — unreadable sources, stat races, control text, special files, and deep cwd retain current domain behavior.
- Command self-observation — aggregate capture cannot change the state being reported.
- Special files and dangling symlinks — domain readers retain current refusal before aggregate or empty rendering.

## Out of scope

- Public schema, field order, flag, `help[]`, family home, stream, or exit changes.
- Replacing domain totals, fingerprint inputs, durable records, or absence authority with generic inference.
- Full-AXI spec-build, coherent-diff, and contextual-disclosure migrations from the FT173 successor sequence.

## Ownership fences

- Query aggregates: `internal/guards/**`, `internal/outline/**`, `internal/roadmap/**`, `internal/worktree/**`.
- Lifecycle/ambient aggregates: `internal/specbuild/**`, `cmd/bench/specbuild.go`, `internal/publication/**`, `internal/status/**`, `internal/dashboard/**`.
- Empty declarations and route conformance: `internal/toon/**`, `cmd/bench/**`, `internal/conformance/**`, `internal/axi/**`, `projects/benchkit.md`. (`internal/toon` matches the AE8/AE10 grant in the ticket-derivation table; the fences list previously omitted it.)
- Existing renderers and the compatibility oracle remain exact-output owners.
