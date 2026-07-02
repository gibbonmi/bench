# Structured state surface

Decisions inherited from `decisions/state-surface.md` (closed) — hybrid output
contract, dedicated subcommands, wave-1 scope, per-script `--describe`
manifests, `--brief` SessionStart injection, one spec.

## Problem

Agents re-derive state `bench` already computes: the phases hand-assemble greps
for open learnings, unresolved decision maps, and structure violations, and
every guard is invisible until an agent collides with it. Two derivations of
the same fact is the bug class the learnings-counter defect came from; the
collisions and re-derivations burn tokens every session on every linked repo.

## Solution

Three agent-facing query subcommands — `bench learnings`, `bench maps`,
`bench guards` — emitting TOON per the AXI standard, sharing one emitter and
`status`'s parsers. Every guard answers `--describe` from the same rules it
enforces; SessionStart injects a one-line-per-guard brief. The gate grows an
AXI-conformance layer (with canaries) for the new surfaces. Existing commands
keep their plain-text contracts untouched.

## User stories

1. As an agent running `/bench-integrate-learnings`, I want `bench learnings`
   to list open journal entries (date, title) as a TOON table with a count, so
   one call replaces hand-reading the journal.
   Line: claude-sonnet-4-6 / medium. This is shell plumbing at a known seam
   whose output shape the gate fully observes.
2. As an agent in `/bench-shape-idea`, I want `bench maps` to list each
   decision map with its unresolved tickets (map, ticket, type), so the
   placeholder grep is never re-assembled.
   Line: claude-sonnet-4-6 / medium. Same known-seam plumbing as story 1 with
   full gate coverage.
3. As an agent, I want `bench guards` to aggregate every guard's manifest —
   and report a deny-capable script lacking `--describe` as `no manifest` —
   so the block surface is learnable without collision.
   Line: claude-sonnet-4-6 / medium. Aggregation is plumbing; the discovery
   convention is decided and the gate observes the output.
4. As the other wave-1 commands, I want one shared TOON emitter for flat
   tables (header + escaped CSV rows) with no new runtime dependency, so
   every surface emits the same shape.
   Line: claude-sonnet-4-6 / medium. A small pure-bash helper whose escaping
   edges are pinned by contract tests.
5. As a reviewer, I want each hook guard (`block-dangerous-git.sh`,
   `check-agent-line.sh`, `stop.sh`, and sourced `_line-guard.sh`) to answer
   `--describe` from its own enforcement tables with deny behavior unchanged,
   so the advertisement cannot drift from the enforcement.
   Line: claude-opus-4-8 / medium. This edits live enforcement scripts where
   a regression weakens a guard, which the profile routes above cheap.
6. As a linked repo, I want the generated pre-push hook to answer
   `--describe` too, and `bench guards` to report it as `not installed` when
   absent, so the git-layer guard is part of the same surface.
   Line: claude-opus-4-8 / medium. It touches the safe-link contract, a
   protected seam.
7. As an agent opening a session, I want `session-start.sh` to also run
   `bench guards --brief` — one line per guard plus one footer line — so the
   block surface arrives with the dashboard.
   Line: claude-sonnet-4-6 / low. A two-line hook change behind an existing
   never-blocking contract.
8. As the kit, I want `status` and `structure` to share one violations
   function instead of `status` grepping `structure`'s human text, so one
   parser serves both.
   Line: claude-sonnet-4-6 / low. An internal refactor already covered by the
   existing status contracts.
9. As every future session, I want the gate to assert AXI conformance on the
   three new commands (TOON stdout, definitive empty states, structured
   stdout errors, exit 0/1/2) with canary fixtures proving each check bites.
   Line: claude-opus-4-8 / medium. Oracle correctness matters more than
   speed; a wrong gate is the worst bug class in this kit.
10. As a future agent, I want the learnings charter (`.bench/learnings.md`
    header and the `bench init` scaffold) to name recurring-ad-hoc-assembly
    as an entry class, so codification candidates get captured.
    Line: claude-fable-5 / high. Guidance prose compounds through every
    session that loads it — the leverage override applies.
11. As a linked-repo agent, I want the kit prose updated to the decided
    state: `craft-cli`'s scope clause and the benchkit profile name the
    hybrid split instead of a blanket exemption, `.bench/BENCH.md` lists the
    new commands, and the phase files (`bench-integrate-learnings`,
    `bench-shape-idea`) point at `bench learnings` / `bench maps` instead of
    instructing ad-hoc greps.
    Line: claude-fable-5 / high. Same leverage override — this prose steers
    every future session.

## Implementation decisions

- New `bin/bench-query.sh`, sourced by the dispatcher like the worktree and
  status siblings, holds the TOON emitter and the three subcommands. Keeps
  `bench-status.sh` inside the structure budget; the shared violations
  function lands where both can source it.
- TOON emitter handles flat tables only: `name[N]{fields}:` header plus
  comma-separated rows; fields containing comma, quote, or newline are
  double-quoted with quotes doubled. The general TOON format is explicitly
  not implemented.
- Errors are one structured line on stdout with an actionable suggestion,
  exit 1; unknown arguments print usage, exit 2; empty results are definitive
  (`learnings[0]{...}:` with a trailing count line), exit 0. Mutations: none —
  all three commands are read-only.
- `--describe` protocol: each guard script checks its first argument before
  reading stdin, prints its manifest (name, boundary, denies, why), exits 0.
  Hook mode is untouched. The git guard prints its verb classes from the same
  python structures that deny; the line guards print the live `lines.env`
  binding. `_line-guard.sh` gets a describe entry point that leaves its
  sourced-usage contract intact.
- `bench guards` discovery: every `.bench/hooks/*.sh` plus `_line-guard.sh`
  plus the installed pre-push. A script answering `--describe` with an empty
  deny (session-start) is listed as informational, not a guard row; a script
  not answering at all is `no manifest`.
- `bench guards --brief` renders one line per deny-capable guard plus one
  footer; this is the surface `session-start.sh` calls, so the rendering
  contract lives in the CLI, not the hook.
- Gate conformance layer is a new contract fragment following the
  `gate-runtime-contracts.sh` pattern, with canary fixtures under
  `tests/canary/` per the existing canary discipline.

## Testing decisions

- A good test here runs the real command in a throwaway fixture repo and
  asserts stdout shape and exit code — external behavior at the CLI seam,
  never parser internals. Prior art: the contract blocks in
  `.bench/gate-runtime-contracts.sh` and the canary fixtures in
  `tests/canary/`.
- Seams tested: the CLI subcommand seam (all query commands), the guard
  `--describe` seam, and the gate/canary seam. The SessionStart hook is
  tested through the CLI seam it delegates to, plus its existing
  never-blocks contract.
- Gate: `.bench/gate.sh` — the project gate, extended by story 9, must be
  green.

### Seam diagram

Seam 1 — CLI query commands (stories 1–4, 7, 8):

    trigger: agent mid-phase / contract test in a fixture repo
        │
        ▼
    .bench/learnings.md      ──▶  [ bench learnings|maps|guards        ]  ──▶  TOON table + count
    decisions/*.md           ──▶  [ (bench-query.sh: shared parsers    ]       or definitive empty
    guard scripts            ──▶  [  + flat-table TOON emitter)        ]       or structured error
                                       ◀ tests attach here: run the command in a
                                         throwaway repo; assert stdout shape + exit code

Seam 2 — guard self-description (stories 5, 6):

    trigger: bench guards aggregation, or direct invocation
        │
        ▼
    the script's own          ──▶  [ <guard>.sh --describe ]  ──▶  manifest: name, boundary,
    enforcement tables                                             denies, why (exit 0)
                                       ◀ tests attach here: invoke each guard with
                                         --describe (fields present) AND replay one
                                         deny case (hook mode unchanged)

Seam 3 — gate conformance + canary (story 9):

    trigger: bench gate
        │
        ▼
    repo tree                ──▶  [ gate AXI-conformance fragment ]  ──▶  green / red + targeted error
                                       ◀ tests attach here: canary fixtures broken on
                                         purpose must go red with the expected substring

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | journal with 2 entry headings → TOON table of 2 + count | 1 | `bash bin/bench.sh learnings` in fixture: no TOON header today (unknown command) | command absent or wrong shape → header/count assertion fails |
| 1 | absent journal and template-only journal → definitive `learnings[0]`, exit 0 | 1 | same fixture, journal deleted / template-only | silence-instead-of-empty or template counted → row fails |
| 2 | maps with `— (open`, `— (deferred`, `GRILL DEFERRED` → one row per unresolved ticket | 1 | `bash bin/bench.sh maps` in fixture with seeded maps | missed placeholder form → row count wrong |
| 2 | no `decisions/` dir → definitive empty, exit 0 | 1 | same, dir absent | absent-vs-empty distinction lost → fails |
| 3 | all five guards present → one manifest row each | 1+2 | `bash bin/bench.sh guards` in linked fixture | missing guard or silent skip → row count wrong |
| 3 | executable hook script without `--describe` → `no manifest` row | 1 | fixture with a stub hook script | silent skip would pass a wrong implementation |
| 3 | pre-push absent → `not installed` row | 1 | fixture without the git hook | definitive state vs omission |
| 4 | field containing comma and quote → escaped TOON row | 1 | fixture entry titled `a, "b"` | naive join corrupts the table |
| 5 | each guard `--describe` prints name/boundary/denies/why, exit 0 | 2 | `.bench/hooks/block-dangerous-git.sh --describe` today: hangs on stdin / nonzero | mode absent → assertion fails |
| 5 | deny behavior unchanged: replay one blocked `git push origin main` case | 2 | already covered — existing guard behavior; re-run after edit | regression in the enforcement path surfaces |
| 6 | linked fixture's generated pre-push answers `--describe` | 2 | link fixture + `--describe` invocation | heredoc missed → fails |
| 7 | session-start output contains the brief block (one line per guard + footer) | 1 | run `session-start.sh` in linked fixture; no brief block today | hook not wired → block absent |
| 7 | outside a repo: prints nothing, exit 0 | 1 | already covered — existing never-blocks contract | guard addition must not break it |
| 8 | `status` structure row agrees with `structure` via shared function | 1 | already covered — existing status budget contract exercises the row | refactor regression surfaces in existing tests |
| 9 | each conformance check goes red on its broken canary fixture | 3 | canary fixture + expected substring, red by construction | a rotted always-pass check fails the canary |
| 10 | charter + scaffold name the new entry class | — | not TDD-able — semantic prose; reviewed, and drift is docs-gate territory | — |
| 11 | kit prose names only real commands post-edit | — | already covered — `gate-docs-contracts.sh` command-currency checks | dead references caught by existing gate |

### Edge inventory

Walked per the profile's hostile-input checklist:

- absent file vs present-but-empty vs template-only journal — coverage rows (story 1).
- malformed input: comma/quote fields — coverage row (story 4); heading with no
  trailing newline — folded into story 1's fixture (entry still counted).
- unknown argument to any new command — usage on stdout, exit 2; one coverage
  row folded into story 9's conformance checks.
- cwd deeper than repo root — commands resolve root via `git rev-parse` like
  `status`; folded into story 9's fixtures (invoke from a subdirectory).
- hostile environment: `python3` missing at `--describe` time — the git guard
  prints a `manifest unavailable (python3 missing)` line and exits 0; coverage
  row folded into story 5's contract block.
- paths with spaces — fixture repos created under a space-containing tempdir
  in the new contract fragment (story 9).
- **Won't handle:** SIGINT mid-command — all three commands are read-only;
  no state to corrupt.
- **Won't handle:** re-run idempotency — read-only queries are trivially
  idempotent; asserting it adds no information.
- **Won't handle:** invocation through a symlink — resolved once in the
  dispatcher (`bench.sh`), which all subcommands inherit; already exercised.

## Out of scope

- Second-wave parsers (`bench diff`, `refs`, `coverage`, `doctor`, `detect`) —
  each a separate capability with its own parser and consumers; parked on the
  roadmap. Estimate if pulled forward: ~40 edits, ~12 gate runs across the set.
- `bench status --json` — decided against in map ticket #3 (dedicated
  subcommands instead); would be ~6 edits, 3 gate runs if ever reopened.
- Public structured mode for `bench structure` — map ticket #4 fixed the
  self-parse internally; a public TOON mode is ~4 edits, 2 gate runs later.
