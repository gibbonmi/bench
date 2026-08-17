# Review: progressive-roadmap

Base `1a8b5a6f` → tip `25fa73a6`. Three axes, run in parallel, `opus`/medium,
~1 iteration each.

## Standards

Raw findings: 7. De-duplicated repair targets: 4. Worst issue: F1.

- **F1 (auto-fix) — the split-board test fixture harness is authored five
  times**, one shadowing itself. Sites: `internal/roadmap/tree_helpers_test.go:50`
  (`writeSplitBoard`), `internal/conformance/docs_workflow_checks_test.go:579`
  (`writeSplitRoadmap`, same signature/body bar an early return),
  `cmd/bench/command_registry_test.go:383`, two inline copies in
  `internal/status/status_test.go`. `internal/roadmap/tree_test.go:218`
  declares a local `writeBoard` closure that shadows the package-level
  `writeBoard` at `tree_helpers_test.go:69`. AGENTS.md's code standard: a
  fixture harness pasted N times must collapse to one source. Fold F5 (the
  `board()` helper's `[2]string` data clump and space-index panic risk at
  `tree_helpers_test.go:78`) into this repair.
- **F2 (auto-fix) — the diagnostic `"<path>: <reason>"` format has two
  derivations.** Produced by `fmt.Sprintf` at 7 sites in
  `internal/roadmap/tree.go`, re-parsed by `strings.Cut(d, ": ")` in
  `internal/roadmap/context_parse.go:189` (whose comment concedes it
  re-derives "ParseDocument's own convention"). One `Diagnostic{Path, Reason}`
  type with a `String()` should serve both call sites.
- **F3 (auto-fix) — PR-talk comments that die at merge**, violating
  `craft-comments`' register: `internal/status/status_test.go:177`,
  `internal/conformance/docs_workflow_checks_test.go:42`,
  `internal/conformance/recurrence_maintenance_contract_test.go:157`,
  `internal/roadmap/tree.go:111`, `internal/roadmap/context_test.go:428`.
- **F4 (ask-user, low) — the row-owner triple** ("body, `Occurrence:` ledger,
  `Sources:` line") is now authored in `.bench/BENCH-reference.md:18`,
  `projects/benchkit.md:64`, `.agents/commands/bench-what-next.md:47`, and the
  layout again in `CONTEXT.md:61`. Arguably honest repetition across
  distinct-role docs — reviewer's call which doc owns it, if any should.

Refuted: the near-identical canary fixture bodies (harness-forced, not a
violation); no new Go dependency; no shell files in the diff;
`.bench/BENCH.md`/`bench-write-spec.md` line budgets held.

## Spec

Raw findings: 4. De-duplicated repair targets: 2. Worst issue: #1.

Verified clean: the per-class disposition table (spec.md:125-130) matches
`internal/roadmap/tree.go` exactly; all seven diagnostic strings match story
text; migration is structurally sound (67/67 rows, zero diagnostics);
`.bench/BENCH.md`/`bench-write-spec.md` budgets held; PR22's differential
evidence is in commit `5ffb0470`'s message.

1. **(auto-fix) PR18/story 21's pin is in the wrong package.** The coverage
   row's seam is "dashboard render test" and `internal/dashboard/` is a
   declared ownership fence, but `TestDashboardRoadmapTextAndSequenceRenderFromSplitTree`
   landed in `internal/status/status_test.go` and asserts
   `roadmap.RoadmapText`/`roadmap.RecommendedSequence` directly rather than
   `dashboard.Snapshot`. Behavior is covered transitively (the dashboard
   composes those same two calls), but a dashboard-side reader swap would not
   red it.
2. **(ask-user) Story 18 — two of eight canary fixtures carry no
   `MUTATE.json`** (`roadmap-unrecognized-file`, `roadmap-unreadable-detail`;
   both plant their fault through the `files/` overlay alone, per the shared
   `internal/canary/mutation.go` harness's inability to add/retype paths via a
   content mutation). Non-behavioral spec-vs-tree-convention contradiction —
   the spec's literal "each with `EXPECT` and `MUTATE.json`" is stale against
   how ~95 other fixtures repo-wide already ship. Flag for veto rather than
   silently reinterpret.

Refuted: spec.md:127 "inline text discarded" reads narrower than the code (a
no-op — only reachable on an already-red tree); the PR5 `{"index absent", ""}`
subtest models absence correctly since `ParseDocument` never reads
`Index.State`.

## Coverage

Raw findings: 8. De-duplicated repair targets: 6 (C1+C6 collapse; C7 is
no-op). Worst issue: C1.

- **C1 (ask-user) — a degraded `roadmap/` directory state is never exercised
  and the oracle ignores it.** `LoadTree` captures `DirState`/`DirReason`
  (`internal/roadmap/tree.go:44`) but `ValidateRoadmapTree`
  (`internal/roadmap/tree_validation.go:10`) reads only diagnostics, never
  `DirState`. With `roadmap/` a regular file/FIFO/unreadable directory and the
  index present, the check emits 67 misleading "missing detail owner" lines
  instead of naming the real cause; with the index absent it is silently
  green. The new `RoadmapDir + "/"` entry in the trust list
  (`internal/roadmap/occurrences.go:94`) is provably unexercised — deleting it
  keeps the suite green.
- **C2 (auto-fix) — an empty row file is silent-capable.** `tree.go:164`
  routes `StateEmpty` into `owners` per design (an empty file is a heading
  mismatch, not an absence — comment at `tree.go:155-156`), but only the
  absent half is asserted (`TestTreeParseReportsMissingDetailOwner`); no test
  drives the empty-file path to its heading-mismatch diagnostic.
- **C3 (ask-user) — a wrapped heading with its row file present also emits a
  spurious orphan diagnostic.** `tree.go:95-98` continues before
  `indexed[id] = true` (line 103), so `roadmap-wrapped-heading`'s own fixture
  (whose `files/roadmap/FT9001.md` exists) emits both `wrapped heading` and
  `orphan detail file with no ROADMAP.md row FT9001` today — the fixture's
  substring `EXPECT` check doesn't notice the second diagnostic. Is a
  double-diagnostic on this class intended?
- **C4 (ask-user) — the diagnostic-source round-trip is lossy for a
  colon-bearing basename.** `context_parse.go:189`'s `strings.Cut(d, ": ")`
  against a legal filename like `roadmap/x: y.md` (reachable via the
  unrecognized-file class) yields `parse_failures.source = "roadmap/x"`, a
  path that does not exist.
- **C5 (ask-user) — a live symlink row file is accepted as authoritative,
  contradicting the spec's own disposition of the edge.** `tree.go:46`'s
  `bounds.Classify` follows symlinks; only a link to a non-regular target
  reports wrong-type. `spec.md:237` disposes of this as "Won't handle … the
  classifier's wrong-type state reports it" — a premise the code refutes for
  a live in-tree-target link. No writer of `roadmap/<ID>.md` exists today, so
  the destructive write-through half doesn't apply yet.
- **C6 (auto-fix, folds with C1's repair) — absent vs. empty `roadmap/` is
  unasserted.** `TestContextCommandSourcesListsRoadmapDirectory` only pins the
  `parsed` state for the new `roadmap/` sources row.
- **C8 (auto-fix) — a test discards the return that would show a real
  fault.** `TestOccurrenceLedgerMalformedAndLineEndings` drives an LF index
  against a CRLF row file and discards the returned diagnostics
  (`valid, failures, _`), silently accepting the `heading does not match`
  fault the mixed-line-ending case actually produces.

Refuted: the migration differential (independently re-derived: 67/67 rows
match base under whitespace normalization; content loss is covered by the
standing `ValidateRoadmapTree` gate check); directory/FIFO/oversized/non-UTF8
row files (one parameterized branch, exercised); trailing-newline-less files,
`roadmap/` subdirectories, NBSP headings, hostile root paths (resolve
correctly). C7 (two fixtures lacking `MUTATE.json`) is the same instance as
Spec finding #2 — recorded there, not double-counted here.
