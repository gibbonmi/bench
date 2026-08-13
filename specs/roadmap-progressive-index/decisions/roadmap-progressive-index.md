# Progressive roadmap index (FT198)

Status: ready

## Destination

`ROADMAP.md` stays the single durable owner of every row's status, evidence,
and rationale; every agent-facing read becomes a bounded projection over the
one existing parse. Bare `bench roadmap` renders a top-10 index view, default
`bench roadmap --context` emits an index with all bodies omitted, and a
fail-closed row selector serves complete bodies on demand — so no ordinary
read pays the full board, no body transports silently partial, and the
command joins the approved AXI query set as one member covering both forms. Scope: `internal/roadmap` projections, the
AXI table/registry/conformance membership, and the `/bench-what-next`
index-first evidence prose. No file split, no new verb, no parser-grammar
change.

## #1: Who is the durable detail owner for a row's evidence and rationale?

Blocked by: none
Type: Grill

### Question

Does row detail move to per-row files, a single details file, or stay in
`ROADMAP.md` with the index and detail served as projections?

### Answer

`ROADMAP.md` remains the only durable owner. The index and per-row detail are
CLI projections over the one existing `ParseDocument` pass; row presence in
the file remains the status source. There is no content migration, history
stays the file's ordinary git history, and index-to-detail completeness is
structural — both views render from the same parse of the same bytes, so no
pairing proof or second status source can exist to drift.

## #2: What surface serves per-row detail on demand?

Blocked by: #1
Type: Grill

### Question

New `show` subcommand, a positional argument on the bare verb, or a selector
on the existing snapshot?

### Answer

A row selector on `bench roadmap --context` (shape like `--row FT198,FT189`).
The snapshot stays the one schema owner; no new verb and no positional
argument on the bare verb.

## #3: What does default `bench roadmap --context` emit?

Blocked by: #1
Type: Grill

### Question

Index-only, or keep today's `limited()` truncated-body preview — and does the
rule cover the ideas, learnings, and retro bodies that also rode the 179 KB
transport failure?

### Answer

Default `--context` omits bodies uniformly — roadmap rows, ideas, learnings,
and retros — while keeping every structured field: ids, titles, dates, paths,
states, counts, byte sizes, and all cross-check blocks (occurrences,
discrepancies, specs, git, gate cache, parse failures). The `limited()`
truncated-body middle ground is removed from the snapshot: a body is either
absent-with-declared-size or complete, never silently partial. Three
consequences ride the same rule: `parse_failures.raw` follows it too (in the
unsupported-schema path today's raw is the entire file, which must not
reintroduce the full-board default read); the per-block `truncated` column is
dropped or replaced by a body-present flag, since nothing truncates anymore;
and every capture unit keeps one identifying field — the ideas block gains
its parsed line number in every mode, the key `capture_occurrences` already
refers to, since one schema owns all modes — so the index remains a complete
capture-unit inventory. A drain
fetches capture content whole via the named per-file paths or `--full`.

## #4: Where is the seam against FT172's parser work?

Blocked by: none
Type: Grill

### Question

Does this work absorb FT172's row-grammar contract and discrepancy checks?

### Answer

No. The grammar contract-test, discrepancy blocks, and workload-boundary
evidence stay FT172's. This build consumes `ParseDocument` exactly as it
stands.

## #5: What does bare `bench roadmap` print?

Blocked by: #1
Type: Grill

### Question

Keep the verbatim 141 KB file dump, or render an index — and how much of it?

### Answer

Bare `bench roadmap` renders a top-10 index view under the AXI contract: the
first 10 parsed rows in document order, index fields only, and a pre-computed
shown-of-total aggregate (`10 of 62 rows`; a board with fewer rows shows what
exists, and a drained board is a definitive zero-row result). Document order is the rule regardless
of section — only `## Features, in priority order` declares a ranking, but
introducing section awareness would be a parser change #4 excludes, so the
aggregate's denominator is every parsed row and the view leans on the board's
convention that earlier means more urgent. Because the whole stdout must
decode as one TOON document with a terminal `help[N]{cmd,why}` block, the
drain-status and next-action facts are re-expressed as TOON blocks or as
`help[]` disclosure rows (today's markdown callouts do not survive), and the
absent-file state becomes a definitive zero-row result on exit 0 whose
disclosure carries the `/bench-what-next` action. The non-document error
postures (empty, failed read, unsupported schema) keep their structured
exit-1 refusals.

## #6: What is the row selector's edge contract?

Blocked by: #2
Type: Grill

### Question

What happens on an absent ID, a malformed ID, and what do selected rows
carry?

### Answer

Fail closed: a requested ID with no matching row exits 1 naming the missing
ID — unsatisfied intent, never an empty success — and a request mixing
present and missing IDs emits no rows at all before that refusal. A malformed
ID is a usage error, exit 2. Selected rows carry complete untruncated bodies.

## #7: What does the drain's complete-evidence claim rest on?

Blocked by: #3
Type: Grill

### Question

`/bench-what-next` currently orders one `--context --full` snapshot as its
complete evidence; what replaces that once the default is index-only?

### Answer

Index-first with targeted fetches: the phase loads the index (the complete
row and capture-unit inventory with byte sizes), then fetches only the rows
and capture units the reconcile touches — the index proves what exists,
fetches prove content. `--context --full` remains the one consistent whole
snapshot for callers that redirect it to a file: unchanged in content —
every body complete — at the new schema version the field changes require.

## #8: Does `bench roadmap` join the approved AXI query set?

Blocked by: #5
Type: Grill

### Question

Formal gate-checked membership (craft-cli table, registry-derived
conformance, `help[N]{cmd,why}` envelope, contextual disclosure), or
principles only?

### Answer

`bench roadmap` joins the approved set as one member — membership is per
command name, so one craft-cli table row, one `axiApprovedRoot` registry
entry, and one `projects/benchkit.md` seam entry cover both the bare index
and `--context`, exactly as `bench coverage`'s single row covers `--check`.
The member's contextual disclosure owns the "request the complete value"
navigation (the exact `--context` / row-selector commands) by contract rather
than by prose, and conformance grades the complete stdout envelope of both
forms.

## Not yet specified

## Spec-writer discretion

- The selector flag spelling and list syntax (`--row FT1,FT2` versus repeated
  flags), within #2's one-schema-owner answer.
- The snapshot schema version bump that #3's field change requires.
- The exact index field list per block, within #3's constraint: every
  structured field but the dropped `truncated` column stays, no body rides.
- The single usage string covering the bare form, `--context`, and the row
  selector — the envelope fixture binds one usage prefix per member.
- The `--full` × row-selector combination rule — refuse as usage, or let the
  selector win; either is reversible and neither changes what the gate proves
  about the two primary modes.
- The contextual-disclosure cell wording for the one new table entry.

## Out of scope

- Any file split: per-row detail files or a second details file (#1 rejected
  both).
- FT172's row-grammar contract-test, discrepancy blocks, and
  workload-boundary evidence (#4).
- A new `show` verb or positional row argument (#2 rejected both).
- The dashboard's `RoadmapText` rendering, `bench status`'s spec-path
  cross-check, and `bench idea --owner` validation — all keep reading the
  same file and parser, untouched.
- Retaining any truncated-body preview mode in the snapshot (#3 removed the
  middle ground).
- ROADMAP.md editing and drain-write conventions; the file's row prose and
  size are unchanged by this work.

## Sources

- Path: `.claude/skills/bench-craft-cli/SKILL.md`
  Supports: #5's minimal-schema/aggregate/disclosure shape and #8's approved query-set mechanics — the ten principles, the Bench application table, and the registry-derived conformance description, read 2026-08-13.
  Drift: re-verify if the AXI skill's approved-set table or conformance description changes before the spec reads this map.
- Path: `internal/roadmap/context_parse.go`
  Supports: #1 and #3's factual premise — `ParseDocument` already projects every row (ID, title, spec, occurrence ledger, body, body bytes) and `BuildContext`/`limited()` own today's truncation; read 2026-08-13.
  Drift: re-verify if FT172's grammar work lands before the spec reads this map.
- Path: `internal/roadmap/roadmap.go`
  Supports: #5's factual premise — the bare verb's verbatim dump, drain-status callout, absent-file prompt, and error postures; read 2026-08-13.
  Drift: re-verify on any bare-verb change landing first.
- Path: `internal/roadmap/context_render.go`
  Supports: #3's per-block field lists — the blocks and columns the index mode amends, including the `truncated` column and the ideas block's missing line number; read 2026-08-13.
  Drift: re-verify if the snapshot renderer changes before the spec reads this map.
- Path: `cmd/bench/main.go`
  Supports: #8's starting state — the `roadmap` registry entry is currently `axiExempt(axiReasonOperational)` and every argument-bearing invocation routes to `ContextCommand`; read 2026-08-13 via the round-1 review.
  Drift: re-verify if the command registry entry or roadmap routing changes.
- Path: `internal/conformance/axi_query_registry_test.go`
  Supports: #8's membership model — per-command AXI disposition (`axiApprovedRoot`, child grammar), and the guidance/registry/profile equality checks; read 2026-08-13.
  Drift: re-verify if the AXI registry conformance model changes before the spec reads this map.
- Path: `cmd/bench/command_registry_test.go`
  Supports: #5's envelope premise — the complete-stdout TOON decode, terminal `help[N]{cmd,why}` block, and the definitive zero-row empty fixture; read 2026-08-13.
  Drift: re-verify if the envelope conformance test changes before the spec reads this map.
