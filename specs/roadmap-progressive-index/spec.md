# roadmap-progressive-index

Status: staged
Roadmap: FT198

Decision source: specs/roadmap-progressive-index/decisions/roadmap-progressive-index.md (ready compiled map, Opus-approved 2026-08-13)

## Problem

`ROADMAP.md` is 141 KB across 62 rows, and every read surface pays for all of
it: bare `bench roadmap` dumps the file verbatim, default `--context`
truncates bodies silently (the FT172 partial-evidence trap), and
`--context --full` exceeded one agent-tool response at 179 KB, which broke a
drain's required evidence capture. An agent who needs ten index facts loads
the whole board or trusts a silently cut body.

## Solution

The file stays the single durable owner; every read becomes a bounded
projection over the one existing parse. Bare `bench roadmap` shows a top-10
TOON board with a shown-of-total aggregate. Default `--context` becomes an
index: every structured field, no bodies, nothing silently partial. A
fail-closed `--row` selector returns complete bodies for exactly the named
rows, and `--full` remains the one consistent whole snapshot. `bench roadmap`
joins the approved AXI query set as one member covering both forms, and
`/bench-what-next` switches to index-first evidence with targeted fetches.

## User stories

1. **As an agent, I load the roadmap index without bodies.** Default
   `bench roadmap --context` emits every row and capture unit with all
   structured fields and byte sizes but no body content, at snapshot schema 4.
   Line: sonnet / low. Exact predicates from the map at a known, test-covered
   package seam.
2. **As an agent, I fetch complete detail for named rows.**
   `bench roadmap --context --row FT198,FT189` returns complete untruncated
   bodies for exactly those rows, failing closed on absent or malformed IDs.
   Line: sonnet / low. Exact contract, existing grammar machinery, package
   tests observe every exit.
3. **As a drain, I still get one consistent whole snapshot.**
   `bench roadmap --context --full` emits every body complete at schema 4,
   redirectable to a file; like every form of the member it now ends with the
   terminal `help` block.
   Line: sonnet / low. Rides story 1's renderer with an inverted body switch.
4. **As a session, I see the top of the board in one bounded read.** Bare
   `bench roadmap` emits one TOON document: the first 10 parsed rows in
   document order as index rows, a shown-of-total aggregate, sequence and
   drain facts as blocks, and a terminal `help` table; absent file is a
   definitive zero-row result on exit 0.
   Line: opus / medium. A surface rewrite replacing markdown porcelain with
   the envelope contract — partial precision at a known shape, gate-covered
   (the decision-table mid row; no cached cheap routing covers a Go renderer
   rewrite).
5. **As the gate, I hold `bench roadmap` to the AXI contract.** The registry
   entry flips to `axiApprovedRoot`, an envelope fixture case binds the member,
   and the guidance/registry/profile equality checks pass with the new entries.
   Line: sonnet / medium. Conformance logic takes the profile's cached mid
   effort; membership mechanics are exact and gate-graded.
6. **As a future session, the guidance matches the shipped surface.** The
   craft-cli approved-query table row and disclosure cell, the
   `projects/benchkit.md` seam entry, the `/bench-what-next` index-first
   evidence doctrine (replacing the ordered `--context --full` single
   snapshot and its schema-3 references), the two conformance contracts that
   pin that doctrine's anchors
   (`internal/conformance/recurrence_maintenance_contract_test.go`,
   `internal/conformance/docs_workflow_helpers_test.go`), the retro-evidence
   workflow anchor in `internal/anchors/registry_data.go`, and one
   `CONTEXT.md` glossary entry for the index/detail vocabulary.
   Line: fable / high. Leverage override — guidance prose compounds through
   every session; cached routing in `projects/benchkit.md` Lines.

## Implementation decisions

- **Selector grammar.** One `--row` flag on the `--context` form taking a
  comma-separated ID list; duplicates dedupe (an idempotent request), the
  empty value refuses as usage. An ID is well-formed when it matches the
  parser's row-start grammar (`[A-Za-z]+[0-9]+`, single-sourced from
  `ParseDocument`'s recognizer, never a second FT-only regex — linked-repo
  boards may carry other prefixes); anything else is usage, exit 2, and a
  well-formed ID with no matching row is unsatisfied intent, exit 1. `--row`
  with `--full` is a usage refusal, exit 2 — two modes, no precedence.
- **One usage string.** The member carries a single usage string covering
  the bare form, `--context [--full]`, and `--row` (decision discretion #4's
  grant, taken here); it begins `usage: bench roadmap` and every grammar's
  help resolves to it, so the envelope's one usage prefix holds across all
  forms.
- **Schema.** The context snapshot advertises schema 4. The `truncated`
  column leaves every block; `body`/`text`/`raw` are empty in index mode with
  their `*_bytes` fields always carrying the true size, and complete in
  `--row` (selected rows) and `--full`. `parse_failures.raw` follows the same
  rule — the unsupported-schema whole-file raw never rides the default mode.
  The `ideas` and `learnings` blocks gain their parsed `line` fields in
  every mode — the keys `capture_occurrences` already names for both unit
  kinds (reviewer-approved extension of decision #3's ideas-only wording;
  same asymmetry, same one-line fix). Every `--context` mode — default,
  `--row`, `--full` — emits the same block list (today's sixteen blocks plus
  the terminal `help`), so `--full` gains the `help` block the member
  envelope requires; "unchanged in content" means every body complete, not
  byte-identical output.
- **Doctrine anchors are oracle edits.** The schema-4 bump, the
  `/bench-what-next` index-first rewrite, and the two conformance contracts
  that pin its anchors land as one atomic sequence: the recurrence contract's
  `` `context.schema = 3` `` and evidence-doctrine anchors move to schema 4
  and index-first wording, and the helpers check's
  `bench roadmap --context` occurrence expectation updates to the new
  doctrine's exact invocation set. Both keep their fail posture and their
  bite tests are updated in the same change, never weakened — landing the
  schema bump alone bricks `/bench-what-next` (the phase stops on any schema
  but the pinned one), and editing the prose alone turns the gate red.
  Consequence for slicing: stories 1–3 and 6 sequence together; stories 4
  and 5 could ship independently green, and bundling them here is a chosen
  bundle surfaced for reviewer sign-off, not a default.
- **Capture bodies under index-first.** The drain obtains retro, learning,
  and idea bodies by reading each file at the path the index names (decision
  #3's sanctioned route — each capture file is small; `--full` stays the
  whole-snapshot alternative for callers that redirect to a file). The
  one-snapshot retro rule therefore narrows: the index is the sole retro
  *inventory* and the run must not re-enumerate `capture/retros/` into a
  second, potentially different listing, but fetching a named file's body is
  the doctrine, not a violation. The `retros`-bodies anchor in
  `internal/anchors/registry_data.go` moves to needle that narrowed rule in
  the same atomic sequence (semantics fixed here; exact prose is story 6's);
  its `Diagnostic` string stays byte-identical, since the
  `implementation-retro-drain-anchor` canary `EXPECT` matches that literal
  against a fixture carrying none of the new prose.
  The doctrine's `bench` invocation set is exactly two spellings: the index
  call `bench roadmap --context` and the fetch form
  `bench roadmap --context --row <ids>`, so the helpers' occurrence
  expectation moves from 1 to the count those two produce, and its other two
  needles ("If the query fails, stop the phase", "manual evidence
  reconstruction") survive the rewritten entry paragraph.
- **Bare verb document.** Blocks, in order: `roadmap` (index rows, first 10
  in document order), `board` (one aggregate row: `rows_shown`, `rows_total`,
  `sequence_trusted`), `sequence` (rank, text, command), `drain` (idea,
  learning, retro counts with source states), `help`. The markdown
  drain-status and next-action callouts are removed from this surface; their
  facts ride the `drain`/`sequence` blocks and `help` disclosure (a pending
  drain discloses `/bench-what-next`). Absent `ROADMAP.md` renders the same
  document with zero `roadmap`/`sequence` rows and the `/bench-what-next`
  disclosure, exit 0. Present-but-empty, failed reads, and unsupported schema
  keep their structured exit-1 refusals. Help spellings (`--help`, `-h`,
  `help`) answer usage on stdout, exit 0 — already true through each
  grammar's parser; the unified usage string is what changes.
- **AXI membership.** `{Name: "roadmap", AXI: axiApprovedRoot}` in
  `commandRegistry`; one `axiEnvelopeCase` (`successMarker: "roadmap["`,
  empty setup = no `ROADMAP.md`); one guidance table row and one profile seam
  command graded by the equality checks; and the independently authored
  `approvedAXIQueries` expectation in
  `internal/conformance/axi_query_registry_test.go` gains `"roadmap": nil` —
  that map is the check's own expectation, not derived, which is why the
  registry flip alone goes red. The `--context`, `--row`, and `--full` forms
  are graded at the `ContextCommand` package seam — acceptable because
  `roadmapCommand` is a two-line passthrough; the bare form rides the real
  dispatcher through the envelope harness. The new disclosure wording is
  pinned by `checkAXIGuidance`'s required-phrase list plus a
  `TestAXIGuidanceContractBites` mutation case (the equality check grades
  membership, not wording): the required phrase is the row-selector command
  in the roadmap disclosure cell, the same way the coverage row's disclosure
  contract is held today.
- **Untouched consumers.** The dashboard's `RoadmapText`, `bench status`
  drain counts, and `bench idea --owner` validation keep their current reads.
  `RecommendedSequence` stays (the dashboard consumes it); only the bare
  verb's markdown rendering of it goes.
- **Roadmap row pairing.** The FT198 row gains this spec's path in the
  staging commit, per the roadmap preamble's cross-check contract; this spec
  carries the matching `Roadmap: FT198` header.

## Testing decisions

- Good tests drive the public command functions (`RoadmapCommand`,
  `ContextCommand`) or the real dispatcher (`Command.Run` via the envelope
  harness) and grade complete stdout plus exit code — never renderer
  internals. Body-omission tests assert both halves: the body is empty *and*
  the byte count is the true size (the hostile-input rule: assert what the
  contract permits, not only what it refuses).
- Seams and prior art: `internal/roadmap` command-level tests
  (`context_test.go`, existing `RoadmapCommand` tests) for stories 1–4;
  `cmd/bench`'s `TestAXIRegistryBindsEachRealCommandEnvelope` fixture table
  and `internal/conformance`'s AXI equality checks for story 5. No new seam
  is introduced.
- Gate seam: `bench gate` (go test over `internal/roadmap`, `cmd/bench`,
  `internal/conformance`) observes every code behavior; the stale-command
  sweep and anchor fixtures observe the prose edits the gate can reach;
  story 6's semantic quality is review-graded.
- TDD applies at the package seams: each coverage row's test lands red before
  its behavior.

### Seam diagram

    trigger: agent shell call / envelope+conformance tests
        │
        ▼
    argv ──▶ [ cmd/bench roadmapCommand routing ] ──▶ stdout + exit code
                │                                   ◀ envelope harness drives
                ▼                                     Command.Run, decodes whole
    ROADMAP.md, capture files, git, gate cache        stdout as one TOON document
        │
        ▼
    bytes ──▶ [ internal/roadmap ParseDocument + BuildContext ] ──▶ ContextSnapshot
                                        │
                                        ▼
              [ renderContext / bare-verb renderer (mode: index|row|full|board) ]
                  ◀ package tests attach here: call RoadmapCommand/ContextCommand
                    with fixture repos, assert stdout, exit, and decoded blocks

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| PI1 | 1 | default `--context` emits every roadmap row with empty `body` and the true `body_bytes` | `internal/roadmap` ContextCommand test | observed red | fails while `limited()` still emits truncated bodies |
| PI2 | 1 | default mode also empties idea, learning, and retro bodies while keeping dates, titles, paths, states, and sizes | ContextCommand test | observed red | catches the cheapest wrong story-1 implementation: roadmap-only omission |
| PI3 | 1 | `ideas` and `learnings` rows carry their parsed `line` in default, `--row`, and `--full` modes | ContextCommand test | observed red | neither field exists today; their absence breaks the `capture_occurrences` join key for both unit kinds |
| PI4 | 1 | the snapshot advertises schema 4 and no block carries a `truncated` column | ContextCommand test | observed red | catches a stale schema advertisement or the dead column surviving |
| PI5 | 1 | over a section-less file, default `parse_failures` rows carry `raw_bytes` with empty `raw`; `--full` carries the complete raw | ContextCommand test | observed red | catches the whole-file raw reintroducing the full-board default read |
| PI6 | 2 | `--row FT<n>` returns complete untruncated bodies for exactly the named rows | ContextCommand test | observed red | catches a selector returning index rows or truncated bodies |
| PI7 | 2 | an absent ID exits 1 naming it and emits no rows, including when mixed with present IDs | ContextCommand test | observed red | catches empty-success filter semantics — the typo-reads-as-no-work false green |
| PI8 | 2 | a malformed ID, an empty `--row` value, and `--row` with `--full` each refuse as usage, exit 2 | ContextCommand test | observed red | catches a lenient parse admitting an ambiguous request |
| PI9 | 3 | `--full` emits every body complete at schema 4 | ContextCommand test | observed red | catches `--full` inheriting index omission — the cheapest wrong story-3 implementation |
| PI10 | 4 | bare-verb success decodes as one TOON document with blocks `roadmap`,`board`,`sequence`,`drain`,`help`, index rows for the first 10 document-order rows, and a true shown-of-total aggregate | envelope fixture + package test | observed red | catches surviving markdown callouts, wrong row selection, and a fabricated aggregate |
| PI11 | 4 | boards with 0 (present, sections, no rows), 9, 10, and 11 rows show min(10, N) rows and the true total, exit 0 | package test | observed red | the top-10 boundary including the drained board — a present row-less file is `StateParsed` with zero rows, a different input state than absence |
| PI12 | 4 | absent `ROADMAP.md` → zero-row document, exit 0, `/bench-what-next` disclosure; empty file, failed read, and unsupported schema → structured exit 1 | package test | observed red | keeps absent and empty distinct and pins the definitive-empty posture |
| PI13 | 4 | the member's single usage string names the bare, `--context [--full]`, and `--row` forms | package test on each grammar's help (`RoadmapCommand(["--help"])` and `ContextCommand` help — the dispatcher routes argument-bearing calls to the context grammar, so the bare grammar is reachable only at this seam) | observed red | today `roadmapGrammar.Help` is `usage: bench roadmap` alone and no help names `--row`; the help-spelling exit-0 behavior itself already passes and rides PI15's envelope subtests |
| PI14 | 4 | a control byte in a row title yields the structured render error, never a corrupt document | package test | observed red | the bare verb newly pushes titles through `toon.Table`; asserts the refusal on the new surface |
| PI15 | 5 | `roadmap` is an `axiApprovedRoot` member: guidance/registry/profile equality green and the envelope fixture passes all six member behaviors | `internal/conformance` + `cmd/bench` envelope harness | observed red | flipping the registry without fixture and table entries (or vice versa) fails the equality and fixture-membership checks |
| PI16 | 2, 3 | `--context`, `--row`, and `--full` stdout each decode as one TOON document with a terminal `help` block | ContextCommand test | observed red | membership promises the envelope on every form; only the bare form rides the fixture table |
| PI17 | 6 | `/bench-what-next` requires the schema-4 index snapshot, index-first evidence doctrine, and the narrowed retro-inventory rule; the recurrence-contract, helpers, and workflow-anchors checks pin the new anchors | `internal/conformance` recurrence-contract, helpers, and `checkWorkflowAnchors` checks against the live prose | observed red | `checkRecurrenceMaintenanceContract`, the helpers' occurrence expectation, and the retro-evidence anchor stay red until prose and all three pinning sites move together — the cheapest wrong implementation (leave the prose alone) is exactly what the current anchors mandate |
| PI18 | 1 | the schema-4 `--context` document enumerates its complete block list — the sixteen current blocks plus terminal `help` — in every mode | ContextCommand test decoding block names | observed red | pins the unenumerated "every structured field stays": an index renderer that drops `capture_occurrences`, `specs`, `git`, or any cross-check block goes red |
| PI19 | 1, 5 | index-mode success disclosure names the exact row-selector command (and `--full`), and the craft-cli disclosure cell carries the row-selector command as a `checkAXIGuidance` required phrase with a bite mutation | ContextCommand test + `checkAXIGuidance` required-phrase list and `TestAXIGuidanceContractBites` mutation | observed red | the body-omission contract rests on disclosure owning "request the complete value"; an empty `help[0]` on success or a vague disclosure cell would otherwise pass — the equality check alone grades membership, not wording |

### Edge inventory

- Absent vs present-but-empty `ROADMAP.md` — PI12 (distinct postures asserted).
- Malformed / unsupported-schema document — PI5, PI12.
- Boundary values — PI11 (0/9/10/11 rows, the drained board as its own
  asserted input state).
- Malformed input (selector) — PI8.
- Control bytes in git- or markdown-sourced text — PI14.
- Error path — PI7, PI8, PI12.
- Empty/absent input — PI12; empty `--row` value in PI8.
- Process-boundary lifecycle — the envelope harness drives the real
  dispatcher from a fresh process-equivalent `Command.Run` with fixture
  repos (PI10, PI15); deep-cwd and symlink invocation ride the fixture's
  standard behaviors.
- **Won't handle:** re-run idempotency — every surface here is read-only;
  repeated invocation cannot change what it reports.
- **Won't handle:** a concurrent writer mutating `ROADMAP.md` mid-read —
  `bounds.Classify` performs one bounded read; a torn read lands in the
  existing malformed posture (structured exit 1), inherited, not new surface.
- **Won't handle:** FIFO/device/socket at `ROADMAP.md` — `bounds.Classify`
  refuses special files before reading today; PI12's failed-read row asserts
  the posture survives through the new renderer.
- **Won't handle:** NBSP and occurrence-grammar hostilities — FT172 owns the
  parser grammar; this build consumes `ParseDocument` unchanged.
- **Won't handle:** paths with spaces or glob characters — the surface takes
  no path arguments; repo-root discovery is inherited.

## Ownership fences

- `internal/roadmap/`
- `cmd/bench/main.go`
- `cmd/bench/command_registry_test.go`
- `.agents/skills/bench-craft-cli/SKILL.md`
- `.claude/skills/bench-craft-cli/SKILL.md`
- `.agents/commands/bench-what-next.md`
- `.claude/commands/bench-what-next.md`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/recurrence_maintenance_contract_test.go`
- `internal/conformance/docs_workflow_helpers_test.go`
- `internal/anchors/registry_data.go`
- `projects/benchkit.md`
- `CONTEXT.md`
- `ROADMAP.md`
- `specs/roadmap-progressive-index/`

## Out of scope

- FT172's parser grammar contract-test, discrepancy blocks, and
  workload-boundary evidence — its own roadmap row (~10 edits, ~3 gate runs).
- A `bench spec show` / `bench outline` symbol selector (FT125's reader
  surfaces) — separate capability (~8 edits, ~3 gate runs).
- Per-row detail files or any file split — rejected by decision #1, not
  deferred.
- Dashboard rendering changes — `RoadmapText` keeps the verbatim contract
  (FT38 owns any identity pass).
- A transcript/session query surface (FT204) — separate operational-surface
  decision.
