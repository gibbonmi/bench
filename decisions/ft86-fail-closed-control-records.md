# FT86 — fail-closed control records and single-sourced repository facts

## Destination

Control-record reads (learnings, maps, roadmap, outline, status, coverage) and
the default-branch fact each get one owner and a fail-closed posture: only
absence is an authoritative empty state; every other failure is a distinguished,
visible state. Sources: `RR:C-01..03`, `RC:H-08`.

## #1: What exit posture does each surface take on a non-absent failure?

Type: Grill

### Question
Today every surface except `roadmap --context` collapses unreadable, wrong-type,
and malformed control records to "empty, exit 0". Gate paths are decided (fail
red); the fork is query commands vs the ambient dashboard.

### Answer
Split posture. Query commands (`learnings`, `maps`, `roadmap`, `coverage`,
`outline`) fail closed: exit 1 with a structured `error:` naming the state.
`bench status` stays exit 0 but renders an explicit per-section `unknown` row
instead of a fabricated zero — the same surfaced-degradation pattern its
worktree-list failure row already uses. Only absence renders as empty.

## #2: Where does the single-sourced file-state classifier live?

Type: Grill

### Question
Two partial helpers exist: `internal/bounds.Read` (size-bounded, typed result)
and roadmap's private `readSource` (the only one classifying file state).

### Answer
Extend `internal/bounds` into the one classified reader: typed state
(absent | empty | malformed | unreadable | wrong-type | unsupported-schema)
plus the existing size bound. Roadmap's `readSource`/`readDirSource`,
learnings, maps, and status all migrate onto it. Classification is Lstat-first:
a dangling symlink is unreadable, never absent (`os.ReadFile` alone reports it
as ENOENT); FIFOs/devices/sockets are wrong-type, rejected before reading;
symlinks to regular files are followed.

## #3: What do diff and roadmap --context do on an unresolvable default branch?

Type: Grill

### Question
`DefaultBranch` fabricates `main`; a sole-`master` repo already resolves via
`ResolvedDefault`, so the residual case is no origin/HEAD with multiple local
branches.

### Answer
`DefaultBranch` is deleted; `internal/git.ResolvedDefault` is the sole owner
and every caller handles `ok=false`. `bench diff` fails closed (exit 1) naming
reality — no resolvable default — plus the `branch.<name>.benchBase` git-config
escape hatch. `roadmap --context` keeps the snapshot and renders the git facts
block as explicit `unknown` state rather than losing the whole context.

## #4: How does a query command report malformed entries inside a readable file?

Type: Grill

### Question
`bench learnings` parses malformed `## ` headings, then silently drops them.

### Answer
Show all, exit 1: render well-formed entries plus explicit malformed-entry rows
(line, reason), and exit non-zero. Same pattern for `bench maps` when a
per-file read fails: the row names the unreadable file instead of silently
dropping it from the listing and the unresolved count.

## #5: Is "unsupported-schema" a version marker or a shape classification?

Type: Grill

### Question
No control record carries a schema marker today; the coverage historical
marker is the closest concept.

### Answer
Shape-based, no new markers. Unsupported-schema = the file reads fine but its
structure isn't one the parser recognizes (a decisions file with no ticket
headings, a roadmap with no recognizable rows) — a distinct state with a
reason, separate from byte-level malformed (invalid UTF-8).

## #6: Coverage-map validation semantics

Type: Grill

### Question
Coverage `--check` exits green with no map, validates only against the max
story number (story 0, reversed ranges, and gap references all pass).

### Answer
Decided by the roadmap row plus one fact: `specs/` is empty, so no existing
spec goes red. A spec without a map and without the existing
`<!-- coverage-map: historical -->` marker fails `--check` (exit 1). Story
references validate exact membership in the declared story set (catches a
reference to story 2 when stories are 1 and 3), require positive IDs (story 0
fails), and ranges must run forward (`3-1` fails).

## #7: One spec or sliced?

Type: Grill

### Question
Three clean seams: (A) classifier + five-surface migration, (B) coverage
semantics, (C) default-branch ownership.

### Answer
One spec. They share the hostile-filesystem fixture harness and the
fail-closed contract vocabulary; B and C are small once decided, and one
review pass grades the posture change coherently.

## Not yet specified

- Exact TOON field names for the new state facts (spec-writer derives from
  `roadmap --context`'s existing absent/empty/malformed vocabulary and the
  established `"unknown"` literal).

## Out of scope

- Outline output bounding/truncation metadata (`RR:C-06`) — separate roadmap row.
- Tamper-proofing or signing control records — FT71 owns evidence posture.
- Schema/version markers in control records (rejected in #5).
- Linked-repo migration tooling — no format changes ship, so none is needed.

## Handoff

1. **Module boundaries.** `internal/bounds` owns file-state classification (the
   one classified reader). `internal/git.ResolvedDefault` owns the
   default-branch fact; `DefaultBranch` is deleted. `internal/coverage` owns
   story-membership validation. `learnings`, `maps`, `roadmap`, `status`,
   `outline`, `diff` are thin consumers of those facts.
2. **Contracts.** Classifier: typed state absent | empty | malformed |
   unreadable | wrong-type | unsupported-schema, Lstat-first, size-bounded.
   Query commands: exit 1 + `toon.Errorf`-shaped `error:` on any non-absent
   failure; absent renders empty, exit 0. `bench status`: exit 0, explicit
   `unknown` section rows on failure. `bench coverage --check`: exit 1 on
   no-map-without-historical-marker, non-member story, story 0, reversed
   range. `bench diff`: exit 1 on unresolvable default, message names the
   `benchBase` escape. `roadmap --context`: git block degrades to `unknown`,
   rest of snapshot intact, exit 0.
3. **Deep vs thin.** Deep: the bounds classifier (hides Lstat/symlink/
   special-file/bound logic) and `ResolvedDefault` (hides origin-HEAD +
   sole-branch fallback). Thin: every surface command — they map classifier
   states to their posture with no reading logic of their own.
4. **Black-box assertables.** Exit codes; `error: <kind> — <hint>` stdout
   strings; TOON state fields (`absent`/`unknown`/`malformed` rows);
   malformed-entry rows with line and reason; coverage `--check` red on each
   membership violation; diff's escape-hatch message text; `roadmap --context`
   snapshot completeness with a degraded git block.
5. **Gate attachment.** All seams gate-visible: `go test ./...` carries the
   contract tests (`internal/contract/axi`) and conformance's
   `checkCoverageMaps` sweep of `specs/*.md`; canaries under `tests/canary/`
   add hostile fixtures (pattern: `coverage-map-validation/broken-coverage-map`
   with `EXPECT`). No TDD-only or manual-verify seam.
6. **Hostile-input owners.** Absent vs present-but-empty → classifier (two
   distinct states, both asserted). Special files (FIFO/device/socket) →
   classifier wrong-type, rejected before read. Dangling symlink → classifier
   unreadable. Permission-denied → classifier unreadable (fixture pattern:
   chmod 0o000 + cleanup, as in status's worktree test). Missing trailing
   newline → each parser. Control bytes in rendered text → existing
   `toon.Representable` posture. Story-gap/range/master-only/unknown-default →
   coverage and git seams per contracts above.
7. **Uncertainty flags.** None — every ticket closed with the reviewer.
8. **Rejected alternatives.** All-surfaces-fail-closed and all-surfaces-degrade
   postures (#1); a new classifier package and promoting `readSource` (#2);
   both-fail and both-degrade for unknown-default callers (#3); hard-fail
   without rows and warn-exit-0 for malformed entries (#4); schema version
   markers (#5); slicing into two or three specs (#7).
9. **Domain watch-outs.** `os.ReadFile` reports a dangling symlink as ENOENT —
   classification must Lstat first or a broken link reads as authoritative
   absence. Conformance's coverage sweep only globs repo-root `specs/*.md`; a
   spec elsewhere is never gate-validated. `roadmap --context`'s
   absent/empty/malformed vocabulary and the worktree `"unknown"` literal are
   the existing state names to reuse, not reinvent.

Dependency order: n/a — single spec (recommended internal build order:
classifier → surface migration → coverage semantics → default-branch
ownership).
