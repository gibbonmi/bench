# AXI bounded projection migration

Status: staged

Decision source: `decisions/byte-preserving-axi-foundation/ft173-axi-contract.md`

## Problem

Bench has four independently implemented truncation policies with different caps, counting units, metadata, full-mode behavior, and safety authority. Leaving them local duplicates mechanics; replacing them with one universal preview would silently change meaning.

## Solution

After the outcome/action migration, route each owner through the parameterized shared projection while leaving selection and semantic facts with the domain. Sanitize retains code-point selection and original byte reporting; roadmap retains UTF-8-safe byte selection and uncapped full mode; worktree projects only after complete safety/fingerprint enumeration; outline retains owner-derived row and skip facts. A final contraction removes only superseded mechanics.

## User stories

1. As an agent reading diagnostic and roadmap text, I want current bounds and full behavior preserved through shared mechanics. Line: gpt-5.6-terra / high. Code-point versus byte semantics and double-truncation are easy to get plausibly wrong.
2. As an agent inspecting worktrees and outlines, I want display bounds shared without moving authority-bearing inventory or semantic totals. Line: gpt-5.6-terra / high. Worktree fingerprint safety and outline skip meaning require separate owner-level probes.

## Implementation decisions

- `sanitize.Preview` retains its 120 Unicode-code-point cap, original byte suffix, control/backslash escaping, and uncapped `Controls` path.
- Roadmap `limited` retains 4096 UTF-8-safe bytes, original byte count, truncation bit, raw content, and uncapped `--full`; already bounded values are not projected twice.
- Worktree inventory remains fully sorted, safety-checked, fingerprinted, and subject to authority byte/entry ceilings before the shared display projection selects 20 default or 1000 full visible entries.
- Outline retains 200 default rows, all rows under `--full`, and owner-derived tracked/scanned/skipped/total/emitted/omitted/truncated facts.
- The contraction removes local selection/count plumbing only where the domain now supplies equivalent facts to the shared projection.

## Testing decisions

- TDD attaches at each existing owner seam and drives below/at/above limits, multibyte cut points, and default/full pairs.
- Mutations change cap inputs, units, totals, full branches, or authority inputs; shared-helper golden tests do not replace owner probes.
- The compatibility oracle proves exact bytes, while route conformance proves every owner actually reaches the shared projection.

### Seam diagram

    trigger: domain has complete content and owner-derived size facts
        │
        ▼
    owner policy + content ──▶ [ shared projection ] ──▶ selected content + unchanged facts
                                   ◀ tests attach here: real owner boundaries and route bypasses

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| BP1 | 1 | `sanitize.Preview` retains 120 code points, control/backslash escaping, original-byte suffix, and uncapped `Controls` while reaching the shared projection. | sanitize preview interface | already covered by `TestPreviewBoundariesAndControls` and `TestControlsEscapesWithoutCapping`; migration adds route observation | Cap, unit, escaping, suffix, and bypass mutations all have independent reds. |
| BP2 | 1 | Roadmap context retains its 4096-byte UTF-8-safe cap, original byte facts, truncation bit, raw content, and uncapped full mode without double projection. | roadmap context owner | already covered by `TestContextBodyLimitBoundaries` and `TestBuildContextCarriesRetrosAndDegradedEvidence`; migration adds route observation | Multibyte, total, full, and double-bound mutations remain observable. |
| BP3 | 2 | Worktree retains complete sorted safety/fingerprint enumeration and applies only the final 20/1000 visible-entry view through shared projection, including lower-bound and stat-race refusal meaning. | worktree ignored inventory and display | already covered by `TestIgnoredInventoryEntryAndByteBoundaries` and `TestIgnoredInventoryStatRaceRetains`; migration adds route observation | Mutating an authority input or projecting early changes the fingerprint/refusal before display bytes can hide it. |
| BP4 | 2 | Outline retains 200 default rows, all rows under full, and owner-derived tracked/scanned/skipped/total/emitted/omitted/truncated facts through shared projection. | outline command result | already covered by `TestCommandBoundsRowsAndFullRetainsMetadata` and skip tests; migration adds route observation | Visible-row-derived totals or ignored skips fail even if the table body remains plausible. |
| BP5 | 1,2 | Per helper-census owner—sanitize preview, roadmap context, worktree ignored display, and outline rows—all four reach their declared projection route and no superseded local derivation or second truncation remains. | registry-derived route conformance and residue census | not TDD-able until BP1-BP4 land | Per-owner bypass and exact symbol/literal residue mutations reject partial contraction. |

### Ticket derivation

Each owner row becomes one independently green migration ticket with one covers ID and atomic policy facts. Existing boundary tests are positive controls; every ticket also owns a direct shared-route bypass mutation. BP5 is the explicit junction/contraction after all four.

| row | tracer acceptance and atomic facts | approved fence | subject mutation | independent owner and public operation |
|---|---|---|---|---|
| BP1 | preserve 120 code points / original bytes / controls / backslash / suffix / uncapped Controls / shared route | `internal/sanitize` | count bytes, alter the cap, drop one escape, cap Controls, or bypass projection | sanitize boundary and route tests; call `Preview`/`Controls` below, at, and above multibyte limits |
| BP2 | preserve 4096 UTF-8-safe bytes / original bytes / truncation / raw content / full uncapped / one projection | `internal/roadmap` | count runes, cut UTF-8, cap full, derive wrong total, double-project, or bypass | roadmap context boundary/route tests; build default and full context fixtures |
| BP3 | preserve complete sorted inventory / safety / fingerprint inputs / authority ceilings / 20 / 1000 / lower-bound / stat-race refusal / shared display route | `internal/worktree` | project before fingerprinting, drop one unsafe preimage entry, change a visible cap, or turn stat-race refusal into success | worktree inventory ownership/route tests; run bounded list/cleanup fixtures with hostile entries |
| BP4 | preserve 200 default / full all / tracked / scanned / skipped / total / emitted / omitted / truncated / shared route | `internal/outline` | derive totals from visible rows, omit a skip, cap full, or bypass projection | outline boundary/skip/route tests; invoke default and full command fixtures |
| BP5 | prove four named helper-census routes and delete each superseded selector/count path / owner / route / residue | `internal/conformance`, `internal/axi`, `internal/sanitize`, `internal/roadmap`, `internal/worktree`, `internal/outline`, `projects/benchkit.md` | restore one exact old selector/count symbol or bypass one named owner | conformance and exact residue census; enumerate the four helper-census owners after migrations |

### Edge inventory

- Error path — BP3 preserves unsafe/stat-race refusal; other owner validation failures retain current sinks.
- Empty or absent input — all owners preserve empty content with zero, unknown, or no-projection semantics as currently defined.
- Boundary values — BP1-BP4 cover below, exactly at, and above every cap and multibyte cut point.
- Malformed input — control bytes and invalid path entries retain current sanitize/refusal behavior.
- Interrupted or partial state — BP3 preserves partial inventory refusal; outline skipped inputs remain owner facts.
- Re-run idempotency — repeated projection returns identical content and facts without accumulating truncation.
- Process-boundary lifecycle — worktree authority tests reload the inventory/fingerprint through a fresh command.
- Hostile environment — spaces/globs, deep cwd, control text, stat races, and symlink/special entries use the project checklist.
- Command self-observation — no owner mutates the source it projects during observation.
- Special files and dangling symlinks — BP3 refuses unsafe inventory entries before projection.

## Out of scope

- Aggregate and empty migration — 14 edits, 1 promotion gate run in `axi-aggregate-empty-migration`.
- New caps, new escape hatches, universal units, or changes to authority-bearing enumeration.
- Any public byte, field, flag, schema, or `help[]` change.

## Ownership fences

- Sanitize: `internal/sanitize/**`.
- Roadmap: `internal/roadmap/**`.
- Worktree: `internal/worktree/**`.
- Outline: `internal/outline/**`.
- Route conformance and contraction: `internal/conformance/**`, `internal/axi/**`, `projects/benchkit.md`.
- Aggregate meaning and all unrelated command domains remain unchanged.
