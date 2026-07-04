# implement-spec-lean — token-efficiency pass on the implement phase

Source: reviewer-requested assessment of `/bench-implement-spec` and the CLI it
touches (2026-07-04 session). The assessment found the phase's per-run cost is
dominated by prose that restates facts canonical elsewhere, plus hand-assembled
shell where `bench coverage` already exists. The larger CLI wrappers
(`bench spec implemented`, `bench commit`, `bench refs`) and a symbol index were
parked on the roadmap (2026-07-04 entries); this map covers only the approved
single slice.

## #1: One spec or several?

Type: Grill

### Question
The assessment produced four recommendations. Does this slice carry all of
them, or only the two the reviewer approved for a single spec?

### Answer
**One spec, two changes:** (a) dedupe the `/bench-implement-spec` command prose
into pointers at its canonical sources, citing `bench coverage` as the
coverage-table source, and (b) make `bench coverage` accept a spec slug as
well as a path. Rejected: folding in the commit-behavior wrappers — they change
what a build commit does and need their own shaping; parked with estimates on
the roadmap.

## #2: What does the dedupe collapse, and what must it keep?

Type: Grill

### Question
The command file restates craft-line (the `> Line:` template), craft-tdd
(TDD-at-seams, seam-set defect), craft-delegate (bound alias, effort in charge,
per-delegate worktrees), and BENCH.md invariants 1 and 4. What collapses to a
pointer, and what is phase-local and stays?

### Answer
**Restatements collapse; phase sequencing and gate anchors stay.** A pointer
names the skill and the one fact it holds ("declare the line per `craft-line`"),
never re-derives it. Phase-local content keeps its prose: the exit handoff, the
stop-short routing, the spec status-flip discipline, the coverage-table
requirement (now citing `bench coverage <spec>` / `--check` as its source
instead of implying hand-assembly). The gate's anchor contracts on this file
are load-bearing and must survive verbatim: "coverage table", "already
covered", "turning red-to-green", "When the build stops short",
"Status: implemented". Rejected: deleting whole sections — the phase must still
read standalone for a cold session.

## #3: How does a slug resolve in `bench coverage`?

Type: Proposal (not grilled — veto at spec approval)

### Question
`bench coverage go-hooks-port` currently errors ("pass a path to a spec
markdown file"). What is the resolution rule?

### Answer
**Path first, slug fallback.** An argument that names a readable file behaves
exactly as today. Otherwise, an argument with no path separator falls back to
`specs/<arg>.md` (appending `.md` only when absent). The not-found error names
both forms tried. Flag-shaped arguments stay rejected as usage errors. Rejected:
slug-only or magic search across directories — one deterministic fallback,
no scanning.

## Handoff

1. **Module boundaries.** Prose: `.agents/commands/bench-implement-spec.md` is
   the single file edited — `.claude/commands` is a symlink to it in the kit
   repo, and linked repos receive it via `bench link`. CLI: slug fallback lives
   in `internal/coverage` `Command` (arg handling only; parser and validator
   untouched); its `--help` usage line gains the slug form. Gate:
   `gate-docs-contracts.sh` gains one anchor requiring the `bench coverage`
   citation in the command file, so a later edit cannot silently drop it.
2. **Contracts.** The five existing anchor phrases on the command file survive
   verbatim (#2). `bench coverage` keeps its exact path behavior, `--check`
   semantics, TOON output shape, and validation phrasings (downstream consumers
   match them by substring); the only new observable is the slug fallback and
   the two-form not-found error. Pointers in the deduped prose name their
   target by the token the stale-reference sweep recognizes (`craft-line`,
   `craft-tdd`, `craft-delegate`, `.bench/BENCH.md`).
3. **Deep vs thin.** No new units. The command file gets thinner; each fact it
   drops must have exactly one surviving source (craft skill, BENCH.md, or the
   CLI's own output). `coverage.Command` remains the one parser/validator for
   the acceptance-map convention.
4. **Black-box assertables.** Go table tests on `Command`: path arg unchanged,
   bare slug hit, slug with `.md` hit, slug miss (error names both forms),
   slug containing a separator gets no fallback, flag-shaped arg still a usage
   error. Prose: the anchor set (existing five plus the new citation anchor)
   and the stale-reference sweep are the observable surface; byte-count of the
   command file is reviewer-checked at review, not gated.
5. **Gate attachment.** `bench gate` as-is: docs anchors, stale-reference
   sweep, structure budgets, `go test`. Prose semantics beyond anchors are not
   gate-observable — per `craft-line`'s leverage override the prose story
   routes top + high; the Go slug story is exact-spec, known-shape, gate-covered
   and routes cheap + low.
6. **Hostile-input owners.** Slug resolving under a repo with no `specs/` dir →
   the two-form not-found error, table-tested. Slug shadowed by a same-named
   file in CWD → path-first rule wins, asserted. Anchor phrases sitting inside
   rewritten sentences → the anchor contracts are substring matches; the build
   re-runs the docs fragment after every prose pass. Cold-session legibility of
   the deduped file → review-phase judgment, named in the spec's testing
   decisions as not TDD-able.
7. **Uncertainty flags.** None genuinely open. One verification the spec must
   record: confirm the new citation anchor lands in `gate-docs-contracts.sh`
   per `craft-gate` (prove it bites by removing the citation once, red, restore).
8. **Rejected alternatives.** Symbol index (portability and staleness cost in
   POSIX shell vs. `rg` already resolving symbols cheaply); commit-behavior
   wrappers in this slice (#1); slug search across arbitrary directories (#3);
   deleting whole prose sections instead of pointer replacement (#2).
9. **Domain watch-outs.** The command file is a leverage artifact — a dedupe
   defect multiplies through every session that loads the phase, which is why
   the prose story rides top tier despite being "just prose". The kit's
   one-source-per-fact standard is the grading rule for the diff: a pointer
   that half-restates its target recreates the drift this slice exists to
   remove. Linked repos see the change only on their next `bench link`.

Dependency order: n/a — single spec.
