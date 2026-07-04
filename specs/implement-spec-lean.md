# implement-spec-lean

Status: staged

## Problem

The `/bench-implement-spec` command phase loads on every implementation session, and a large share of its tokens restate facts that are canonical elsewhere: the `> Line:` declaration template (owned by `craft-line`), TDD-at-seams and the seam-set defect (owned by `craft-tdd`), the delegation rules — bound alias, effort in the charge, per-delegate worktrees (owned by `craft-delegate`), and BENCH.md invariants 1 and 4. The same phase hand-waves at the acceptance-coverage table as if it were hand-assembled, when `bench coverage <spec>` / `--check` already produces and validates it. Separately, `bench coverage` only accepts a filesystem path — `bench coverage go-hooks-port` errors — so a caller who knows the spec slug must retype the `specs/….md` path.

Both are one-source-per-fact defects on a leverage artifact: a restatement that drifts multiplies through every session that loads the phase.

## Solution

One slice, two changes.

1. **Dedupe the command prose.** Collapse each restatement in `.agents/commands/bench-implement-spec.md` into a one-line pointer that names its canonical source and the single fact it holds, never re-deriving it. Phase-local content stays as prose. The coverage-table requirement now cites `bench coverage <spec>` / `--check` as its source. A new gate anchor pins that citation so a later edit cannot silently drop it.

2. **Teach `bench coverage` the slug form.** `internal/coverage` `Command` resolves its argument path-first, then falls back to `specs/<slug>.md` for a separator-free argument, so `bench coverage go-hooks-port` works. Every other behavior of the command is untouched.

## User stories

1. As an agent loading `/bench-implement-spec`, I want the phase to point at `craft-line`, `craft-tdd`, `craft-delegate`, and `.bench/BENCH.md` for the facts those sources own — instead of restating them — so that the phase costs fewer tokens and no restatement can drift from its source, while the phase still reads standalone for a cold session (exit handoff, stop-short routing, spec status-flip discipline, and the coverage-table requirement keep their prose). `Line: claude-fable-5 / high / ~40k.` This edits a command phase that steers every session which loads it, so `craft-line`'s leverage override routes it top and high regardless of how mechanical the diff looks.

2. As an agent driving `bench coverage`, I want to pass a spec slug (`bench coverage go-hooks-port`) as well as a path, resolved path-first with a `specs/<slug>.md` fallback, so that I don't retype the full path when I already know the slug. `Line: claude-sonnet-5 / low / ~30k.` The spec is exact, the seam is a single known function with table-test prior art, and `go test` fully observes the behavior, so the decision table's first row selects cheap and low.

3. As a reviewer, I want `gate-docs-contracts.sh` to carry one new `require_anchor` line enforcing the `bench coverage` citation in the command file, so that a future prose edit that drops the citation turns the gate red. `Line: claude-sonnet-5 / low / ~15k.` The change copies an established `require_anchor` pattern and a wrong needle turns the gate red immediately, so the gate covers the change — cheap and low, with the `craft-gate` bite-proof as the one manual verification.

## Implementation decisions

**Story 1 — prose dedupe (`.agents/commands/bench-implement-spec.md`, the single file edited).** `.claude/commands` is a symlink to it in the kit repo; linked repos receive it via `bench link`, so no second file changes. The collapse targets, each becoming a one-line pointer at its named source:

| Restated today | Collapses to a pointer at | Fact the pointer keeps |
|---|---|---|
| The `> Line: <model> / <effort> / ~<cap>` declaration template | `craft-line` | declare the line before touching code |
| TDD-only-at-seams and the seam-set defect | `craft-tdd` | TDD only at the pre-agreed seams |
| Bound alias, effort in the charge, per-delegate worktrees | `craft-delegate` (model half: `craft-line`) | every delegation carries its own line |
| Gate-is-oracle / never weaken a test (invariant 1) | `.bench/BENCH.md` | done means the gate is green |
| Smallest diff, read before write, compose seams (invariant 4) | `.bench/BENCH.md` | one small change, repo stays green |

Phase-local prose that stays: the entry orientation, the exit handoff, the "When the build stops short" routing, the `Status: staged → implemented` flip discipline, and the coverage-table requirement — the last now reading that the table comes from `bench coverage <spec>` / `--check`, not hand-assembly. Each pointer names its target by the token the stale-reference sweep recognizes (`craft-line`, `craft-tdd`, `craft-delegate`, `.bench/BENCH.md`). A pointer that half-restates its target recreates the drift this slice removes — the kit's one-source-per-fact standard is the grading rule for the diff. The five existing anchor phrases survive verbatim inside the rewritten prose: `coverage table`, `already covered`, `turning red-to-green`, `When the build stops short`, `Status: implemented`.

**Story 2 — slug fallback (`internal/coverage.Command`, arg handling only).** `parse`, `State`, `Rows`, `Check`, and the validator are untouched; only the argument-to-content resolution changes. The existing flag rejection (`strings.HasPrefix(a, "-")` → usage error, exit 2) and the "`<spec.md>` is required" empty-arg error stay ahead of resolution. New resolution rule, applied to the single collected `spec` argument:

- Read `spec` as given. If readable, use it (path-first — this preserves today's path behavior and lets a same-named file in CWD shadow the fallback).
- Otherwise, if `spec` contains **no** path separator, retry `specs/<spec>.md`, appending `.md` only when `spec` does not already end in it. If that reads, use it.
- If neither reads, return a not-found error that names each form actually tried — both `<spec>` and `specs/<spec>.md` for a slug, the single path form for an argument that contained a separator (no fallback was attempted). This replaces the current single-form `spec not found:` message; it is the not-found error, not a `Check` validation phrasing, so downstream substring consumers are unaffected.

Resolution runs before the `--check` branch, so `bench coverage --check <slug>` resolves identically. The `-h`/`--help` usage line gains the slug form (e.g. `usage: bench coverage [--check] <spec.md | slug>`); the `toon.Usage` arg-shape errors are unchanged. `bench coverage`'s path behavior, `--check` semantics, and TOON output shape are otherwise identical.

**Story 3 — gate anchor (`.bench/gate-docs-contracts.sh`).** One new `require_anchor ".agents/commands/bench-implement-spec.md" "<citation needle>"` line, grouped with the existing five anchors on that file, where the needle is the distinctive citation substring introduced by story 1 (e.g. `bench coverage <spec>`). Per `craft-gate`, prove it bites: remove the citation once → run the docs fragment → it errors → restore. Sequencing: story 1 lands the citation text first, then story 3 pins it.

No new units, no new modules. Every fact the command file drops keeps exactly one surviving source; `coverage.Command` stays the one parser/validator for the acceptance-map convention.

## Testing decisions

- **What a good test is here.** For story 2, drive `Command` through its public `(args []string) (string, int)` interface and assert the returned output and exit code — never reach into `parse`/`Check`. For stories 1 and 3, the observable surface is the gate's own output: the `require_anchor` grep set (existing five plus the new citation anchor) and the stale-reference sweep, run through `.bench/gate-docs-contracts.sh`. Byte-count of the command file is reviewer-checked at review, not gated.
- **Seams tested, and prior art.** (1) `coverage.Command` — new table-test function; the existing `internal/coverage/coverage_test.go` (`TestStateAndRows`, `TestCheck`) is the table-test prior art, but those parse in-memory, so the new `TestCommand` adds fixtures on disk: create `specs/<slug>.md` under `t.TempDir()` and run with CWD set there (`t.Chdir`). (2) The docs-contracts gate over `.agents/commands/bench-implement-spec.md` — no Go test; the anchor set and sweep are the assertion, exercised by running the gate fragment.
- **Gate command.** `bench gate` (project default) must stay green: the Go seam adds cases to `go test ./...`; the prose and anchor seams run through `.bench/gate-docs-contracts.sh` (docs anchors + stale-reference sweep + structure budgets).
- **Not TDD-able, recorded here per the map's Handoff.** Cold-session legibility of the deduped command file, and whether each pointer half-restates its target (one-source-per-fact), are review-phase judgments the gate cannot grade — carried as `not TDD-able` rows below, not TDD coverage. Story 3's anchor bite-proof (remove citation → red → restore) is a runnable red command but can only run after the anchor exists, so it is `not TDD-able` as a pre-implementation red and is verified by the `craft-gate` bite ritual instead.

### Seam diagram

Seam 1 — `coverage.Command` (story 2):

    trigger: `bench coverage [--check] <arg>` from the shell or the gate's docs layer
        │
        ▼
    <arg>  ──▶  [ Command: reject flag → resolve arg (path-first,        ]  ──▶  TOON table (default)
                [   then specs/<slug>.md) → read → parse/State/Rows/Check ]  ──▶  violations (--check)
                                                                            ──▶  not-found error (names forms tried)
                                                                            + exit code (0 / 1 / 2)
                      ◀ tests attach here: call Command([]string{…}) with fixtures under a
                        temp CWD (t.TempDir + t.Chdir), assert the (output, exitcode) pair

Seam 2 — docs-contracts gate over the command file (stories 1 and 3):

    trigger: `bench gate` → .bench/gate-docs-contracts.sh
        │
        ▼
    .agents/commands/bench-implement-spec.md  ──▶  [ require_anchor grep -qF ×6 ]  ──▶  pass
                                              ──▶  [ stale-reference sweep      ]  ──▶  err (missing anchor / stale ref)
                      ◀ tests attach here: run the docs fragment; with the new anchor present,
                        drop the citation line → the citation anchor errs (bite-proof)

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 2 | a readable path argument resolves and prints its TOON table exactly as today | Command | `go test -run TestCommand ./internal/coverage` — red on the new path case before the resolve step is added | if resolution regressed the existing path, the round-trip output would differ |
| 2 | a separator-free slug resolves `specs/<slug>.md` and prints its table | Command | same run — red: today the slug returns a not-found error | the assertion fails unless the fallback reads `specs/<slug>.md` |
| edge of 2 | a slug already ending `.md` is not double-appended | Command | same run — red on the `.md`-suffixed slug case | double-append would look up `specs/foo.md.md` and miss, so the table would be empty/error |
| edge of 2 | a slug that matches no file errors, naming both `<slug>` and `specs/<slug>.md` | Command | same run — red: today's error names only one form | asserting both substrings fails unless the two-form error is produced |
| edge of 2 | a separator-bearing argument gets no fallback (single-form path error) | Command | same run — red: assert the error does not mention a `specs/…` form | if fallback ran on a path, the error would carry the extra form and the assertion would fail |
| edge of 2 | a slug shadowed by a same-named readable file in CWD resolves the CWD file (path-first) | Command | same run — red on the shadow case | if fallback preceded the path read, the CWD file would be ignored and the wrong content returned |
| edge of 2 | a flag-shaped argument stays a usage error, exit 2 | Command | already covered — the existing `strings.HasPrefix(a, "-")` branch returns exit 2 | resolution runs after flag rejection, so no new code path touches this; pinned by an assertion, not new coverage |
| 1 | the five existing anchors survive the rewrite and the stale-reference sweep stays green | docs-contracts gate | already covered — `require_anchor` ×5 and the sweep already run in `bench gate`; a dropped anchor or a pointer to a missing target turns them red | the pre-existing checks fire on any dropped anchor phrase or dangling pointer token |
| 1 | each restatement collapses to a single-source pointer, and the deduped file still reads standalone | docs-contracts gate | not TDD-able — half-restatement and cold-session legibility are review-phase judgments the gate cannot grade | recorded per Handoff item 6; caught at `/bench-review-implementation`, not by a test |
| 3 | the new anchor enforces the `bench coverage` citation — dropping it turns the gate red | docs-contracts gate | not TDD-able — the bite-proof (remove citation → run docs fragment → red → restore) can only run once the anchor exists | verified by the `craft-gate` bite ritual, which proves the anchor bites rather than a pre-implementation red |

### Edge inventory

Edge classes walked for the `Command` seam, each resolved above or excluded here:

- Error path → slug miss with the two-form not-found error (coverage row).
- Empty/absent input → no spec argument: **already covered** by the existing `<spec.md> is required` usage error.
- Boundary values → slug already ending `.md` (coverage row).
- Malformed input → flag-shaped argument (coverage row, already covered).
- Hostile environment → repo with no `specs/` dir (folds into the slug-miss two-form error row); slug shadowed by a same-named CWD file (path-first row).
- Interrupted/partial state — **Won't handle** — `Command` performs one read with no partial or mutable state; there is nothing to interrupt.
- Re-run idempotency — **Won't handle** — resolution and parsing are pure reads of the same inputs, so output is deterministic across re-runs by construction.

Every excluded edge leaves the primary calling convention (`bench coverage <path|slug>`, with and without `--check`) fully exercisable.

## Out of scope

- **`bench spec implemented`** — a wrapper that flips a spec's `Status:` line in the green-gate commit. A separate capability (it changes what a build commit does and needs its own shaping); parked on the roadmap. Estimate: ~3 edits, 2 gate runs.
- **`bench commit`** — a wrapper that stages the build's touched files and refuses on an unexplained working-tree file. A separate commit-behavior capability with its own shaping; roadmap-parked. Estimate: ~4 edits, 2 gate runs.
- **`bench refs`** — a reference-refactor helper that dry-runs and verifies old stems in every form before a rename. A distinct capability; roadmap-parked. Estimate: ~4 edits, 2 gate runs.
- **Symbol index** — a precomputed symbol map for the CLI. Rejected for this slice on portability and staleness cost in POSIX shell versus `rg` already resolving symbols cheaply; a separate capability if ever revived, roadmap-parked. Estimate: ~5 edits, 3 gate runs.
