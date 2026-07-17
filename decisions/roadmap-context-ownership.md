# Roadmap context ownership

## Result

`internal/roadmap` should own the context snapshot and command extension because it
already owns roadmap display, idea capture, drain counts, and the recommended-sequence
parser. The snapshot is a deep aggregation module: it gathers typed facts from their
owners, applies size policy, and renders AXI/TOON; it does not reparse another command's
stdout. The dashboard demonstrates this composition pattern, but its snapshot is a
lossy human view and is not the context schema ([dashboard gatherer](../internal/dashboard/dashboard.go#L79)).

## Fact inventory

| Evidence | Current source | Smallest one-source change |
|---|---|---|
| Roadmap presence and raw bytes | `RoadmapText` owns absent-versus-empty behavior ([source](../internal/roadmap/roadmap.go#L104)). | Keep this reader; add one parsed `Document` model for feature rows and sequence items. `RoadmapCommand`, the dashboard, status reconciliation, and context must consume that parser rather than add regexes. |
| Recommended sequence | `RecommendedSequence` is already the shared fence-aware parser ([source](../internal/roadmap/roadmap.go#L146)). | Preserve it as the section parser and derive typed numbered items from its result. |
| Ideas | `ideaLines` is the one reader behind both counts and `ParkedIdeas`, but it drops every malformed/non-entry line ([source](../internal/roadmap/roadmap.go#L193)). | Promote it to a typed inbox parse containing entries plus raw malformed content; keep counts and existing views derived from that parse. |
| Learnings | `learnings.Rows` owns open-heading recognition but returns only date/title and silently skips malformed content ([source](../internal/learnings/learnings.go#L17)). | Add a typed journal parser with complete entry bodies and malformed fragments; derive `Rows` and drain counts from it. |
| Structure findings | `structure.Check` is the one rule engine, but its public result is a formatted report plus count ([source](../internal/structure/structure.go#L33)). | Make the engine return typed violations, warnings, grants, and limits; render the existing report and count from that result. Parsing `bench structure` output would create a second contract. |
| Spec state | `spec.AwaitsRetirement` is the canonical implemented predicate ([source](../internal/spec/spec.go#L39)); directory enumeration is private inside status ([source](../internal/status/status.go#L461)). | Add a typed spec inventory in `internal/spec`; status retirement counts and context consume it. |
| Spec retirement history | The history parser, merge, and git queries are canonical but private to the renderer ([source](../internal/spec/history.go#L31), [source](../internal/spec/history.go#L109)). | Export a typed history query and keep `bench spec history` as its TOON renderer. Context queries history only for roadmap spec references. |
| Git branch, dirty, and ahead state | Git primitives are exported, while status privately derives a lossy signal and tolerates query errors ([source](../internal/git/git.go#L73), [source](../internal/status/status.go#L240)). | Add a typed repository-state query in `internal/git`; status reduces it to a signal and context reports the full state or fails closed. |
| Cached gate verdict | `status.GateVerdict` already returns the complete typed cache state and computes staleness once ([source](../internal/status/status.go#L171)). | Reuse unchanged; do not run the gate. |
| External-evidence markers | Roadmap rows are currently opaque Markdown; no owner classifies graduation triggers ([source](../internal/roadmap/roadmap.go#L77)). | Carry the row body losslessly and mark explicit external-trigger phrases through the single `Document` parse. Do not probe providers or infer whether the trigger fired. |

## Module shape

- `internal/roadmap` gains the typed `Document`, `ContextSnapshot`, options, gatherer,
  truncation metadata, and AXI renderer. This extends an existing seam instead of
  adding a pass-through package.
- `internal/learnings`, `internal/structure`, `internal/spec`, and `internal/git`
  expose typed facts from their existing engines. Their current commands and status
  rows become projections of those facts, preserving one source per meaning.
- `internal/dashboard.Snapshot` and `status.Signals` remain human/ambient projections.
  They omit entry bodies and deliberately tolerate several read failures, so composing
  them would violate the context command's complete, fail-closed contract
  ([dashboard snapshot](../internal/dashboard/dashboard.go#L27), [signal gatherer](../internal/status/status.go#L73)).
- The public command remains one seam: `bench roadmap --context [--full]`. Exact TOON
  fields and the truncation ceiling wait for the gate and hostile-input research.

This shape introduces no harness adapter and performs no network or gate execution.
