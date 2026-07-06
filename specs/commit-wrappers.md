# commit-wrappers — `bench spec implemented` + `bench commit`

Status: implemented

## Problem

The commit discipline lives as prose duplicated across two command files
(`/bench-implement-spec`'s close-on-green and `/bench-final-check`). An agent must
*remember* to gate before committing, to stage only named files, to block on an
unexplained working-tree file, and to flip the spec's `Status:` line in the same
green commit. Prose the agent scrolls past is a footgun: a second gate run gets
skipped, a stray `git add -A` rides an untracked file in, or the status flip is
hand-typed in the wrong format and `bench status` mislabels an unfinished spec as
done. The invariant "commit on green, never on red" is only as strong as the
agent's memory of prose.

## Solution

Two CLI wrappers mechanize the discipline so the oracle lives in code, not memory:

- `bench spec implemented <slug>` — a standalone primitive that flips exactly the
  one `Status: staged` line to the retirement-detector form, and nothing else. The
  single source of the flip logic.
- `bench commit [-m <msg>] [--spec <slug>] <path>...` — block-check the tree,
  run the gate, and commit only on green, staging exactly the named paths. With
  `--spec`, it composes the flip so the finishing build-commit is one command.

Then the duplicated commit-discipline prose in the two command files collapses to
pointers at the two wrappers — one source for the fact.

## User stories

1. As an agent finishing a build, I want `bench spec implemented <slug>` to flip
   the single `Status: staged` line to `Status: implemented`, so the flip logic
   has one tested source instead of hand-typed prose.
   Line: claude-sonnet-5 / medium. The flip is an exact-spec, gate-covered Go
   primitive at a known seam, so it routes cheap, with medium effort because the
   exact retirement-regex output format is load-bearing.

2. As an agent, I want `bench spec implemented` to resolve a bare slug to
   `specs/<slug>.md` and a path argument as-given, reusing `coverage.Command`'s
   exact convention, so both commands take the spec argument the same way.
   Line: claude-sonnet-5 / low. This composes the existing `resolveSpec` convention
   rather than inventing resolution, so it is mechanical at a known seam.

3. As an agent, I want the flip to rewrite only that one line to the exact form
   `Status:<tab/space>implemented` and leave every other byte identical, so the
   retirement detector fires correctly and no other spec is disturbed.
   Line: claude-sonnet-5 / medium. Cheap because gate-observable, medium because a
   byte-preserving rewrite against a load-bearing regex is where a wrong output silently corrupts status.

4. As an agent, I want the flip to exit nonzero and name the file plus the reason
   on not-found, no-`Status: staged` line, or more than one such line, so a typo or
   a re-run on an already-implemented spec is non-destructive.
   Line: claude-sonnet-5 / low. The error branches are enumerated in the map, so
   this is exact-spec plumbing at the same seam.

5. As an agent, I want `bench commit -m <msg> <path>...` to run the project gate
   and commit only on green — reporting the first failing phase and exiting nonzero
   on red — so commit-on-red is impossible by construction.
   Line: claude-sonnet-5 / medium. Cheap because the whole behavior is black-box
   gate-observable (git state + exit code), medium because it composes the oracle and
   ordering matters.

6. As an agent, I want `bench commit` to refuse — before gating — if any
   tracked-modified or untracked file exists outside the named set (plus the
   `--spec` file), listing the offenders and exiting nonzero, so a green verdict
   describes exactly the diff that lands.
   Line: claude-sonnet-5 / medium. Cheap and gate-observable; medium because the
   block-check must precede the gate and its scope (named set plus `--spec`) is the crux.

7. As an agent, I want `bench commit` to stage only the named paths via explicit
   `:(literal)` pathspecs and never `git add -A`, so a second untracked file can
   never ride into the commit.
   Line: claude-sonnet-5 / low. The `:(literal)` pathspec pattern is prior art in
   `internal/shift`; this is mechanical reuse.

8. As an agent, I want `bench commit --spec <slug>` to flip the spec and include
   that spec file in the same single commit, so the finishing build-commit is one
   command instead of a flip-then-commit sequence.
   Line: claude-sonnet-5 / medium. Cheap and gate-observable; medium because it
   composes `internal/spec` and must add the spec path to the named set atomically.

9. As an agent, I want `bench commit` to be branch-agnostic and commit on the
   current branch including `main`, because the harness-independent branch
   protection is the pre-push hook, not a second commit-time guard.
   Line: claude-sonnet-5 / low. This is the *absence* of a guard — no new machinery,
   the simplest case.

10. As an agent, I want `bench commit` to refuse an empty commit (named paths that
    produce no staged change) with a nonzero "nothing to commit", so a re-run after
    a clean commit fails loudly instead of writing an empty commit.
    Line: claude-sonnet-5 / low. A single pre-commit staged-diff check at the known seam.

11. As an agent on any harness, I want both commands routed through `bin/bench.sh`
    into the Go core, so the kit CLI, the linked by-path CLI, and hooks all reach
    the same implementation.
    Line: claude-sonnet-5 / low. The top-level route and the `spec` namespace follow
    the existing `gate` / `gate pin` routing pattern exactly.

12. As the gate, I want the conformance phase to assert `bin/bench.sh` routes both
    `commit` and `spec implemented`, so a later edit cannot silently drop a route.
    Line: claude-sonnet-5 / medium. Gate/conformance logic routes cheap here because
    it is a route-anchor with a mechanical bite-proof, at medium effort per the profile's gate-logic rule.

13. As a maintainer, I want the duplicated commit-discipline prose in
    `/bench-implement-spec` and `/bench-final-check` collapsed to pointers at the two
    wrappers, so the discipline has exactly one source.
    Line: claude-fable-5 / high. This is the `craft-line` leverage override,
    pre-agreed at the map's close: command prose compounds through every session that
    loads it, so authoring it is worth the top tier even though the edit is small.

14. As a reviewer, I want usage errors (no paths, no `-m`, unknown flag, `spec`
    with no subcommand) to exit 2, distinct from operational failures which exit 1,
    so the exit code tells a misuse apart from a real block.
    Line: claude-sonnet-5 / low. The exit-code contract is enumerated in the map's
    item 2; this is exact-spec argument handling.

## Implementation decisions

- **`internal/spec` (new)** owns `bench spec implemented <slug>`: resolve
  (path-first, then `specs/<slug>.md` for a separator-free arg — reuse the exact
  `resolveSpec` convention from `internal/coverage`), validate exactly one
  line-start `Status: staged`, and rewrite that one line to
  `^Status:[ \t]+implemented[ \t]*$` form. It edits the file only; it never stages.
  This is the single source of the flip logic, deep over resolution + validation +
  byte-preserving rewrite.
- **`internal/commit` (new)** owns `bench commit`: a thin orchestrator sequencing
  block-check → gate → stage → commit, plus the optional `--spec` flip. It reuses
  `internal/gate`'s `RunAndRecord` (resolve + run + record the verdict cache) and
  does **not** touch `internal/shift`'s loop. When `--spec` is set it composes
  `internal/spec` for the flip and adds the spec path to the named set. Staging uses
  explicit `:(literal)` pathspecs (prior art: `internal/shift.stageTouched`); never
  `git add -A` over the tree. Root and spec resolution run after `git.Root`, so a
  cwd deeper than the repo root resolves correctly.
- **`cmd/bench/main.go`** gains `commit` as a `run()`-switch case (like `shift` /
  `gate-run`, because it streams the gate's live output) and `spec` as a namespace
  case dispatching `implemented`. `spec implemented`'s single confirmation line and
  errors return through the same case.
- **`bin/bench.sh`** gains a top-level `commit) route_porcelain "$@"` route and a
  `spec` namespace routing `spec implemented` into the core — mirroring the existing
  `gate` / `gate pin` shape — so every shipped surface reaches one implementation.
- **`.agents/commands/bench-implement-spec.md`** and **`bench-final-check.md`**: the
  duplicated commit-discipline prose collapses to pointers naming `bench commit` and
  `bench spec implemented` by the token the stale-reference sweep recognizes. The
  final-check prose branch-guard (switch-off-default) is *removed*, not ported —
  the pre-push hook owns that fact and it breaks a main-is-working-branch repo.

## Testing decisions

- **A good test here** drives each command as a black box in a throwaway git repo
  and asserts git state, file content, exit code, and stderr — never an internal
  helper. Prior art: the Go table tests in `internal/coverage/coverage_test.go` and
  `internal/shift/shift_test.go`, and the fixture-repo style already used across
  `internal/`.
- **Seams tested:** the two CLI `Command` boundaries — `spec.Command` and
  `commit.Command`. One seam each; both fully gate-observable.
- **Gate command:** the project gate, `.bench/gate.sh`. Both packages' table tests
  run under its `go test` conformance phase; the route-anchor runs under the
  conformance suite (`internal/conformance`).
- **Not TDD-able (stated):** cold-session legibility of the collapsed command prose
  is `/bench-review-implementation` judgment. The stale-reference sweep is the only
  gate-observable signal on the prose half — it bites if a pointer names a
  command token that does not route.

### Seam diagram

Seam 1 — `bench spec implemented <slug>`:

    trigger: agent runs `bench spec implemented <slug>` (or `bench commit --spec`)
        │
        ▼
    slug/path  ──▶  [ resolve (path-first, specs/<slug>.md) ]  ──▶  file rewritten:
                    [ validate exactly one `Status: staged` ]       one line staged→implemented
                    [ rewrite that one line, byte-preserve   ]  ──▶  confirmation line + exit 0
                        │                                            OR exit 1 (file + reason)
                        ◀ tests attach here: run Command in a fixture repo,
                          assert file bytes / exit code / stderr

Seam 2 — `bench commit [-m <msg>] [--spec <slug>] <path>...`:

    trigger: agent runs `bench commit -m "<msg>" [--spec <slug>] <path>...`
        │
        ▼
    named paths ──▶ [ 1. block if any dirty/untracked file outside named set ]
    -m msg      ──▶ [ 2. gate (RunAndRecord) — red ⇒ report + exit 1        ] ──▶ commit lands
    --spec slug ──▶ [ 3. flip spec (if --spec), add to named set            ]      on current branch
                    [ 4. stage named paths via :(literal); refuse empty     ] ──▶ exit 0 (summary)
                    [ 5. commit -m msg                                       ]      OR exit 1 / 2
                        ◀ tests attach here: run Command in a fixture repo,
                          assert committed tree / branch / exit code / stderr

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | staged→implemented flip writes the exact retirement-regex line | spec.Command | `go test ./internal/spec` — package absent, test fails to build | no flip primitive exists; the wrong output format leaves `Status: staged` or a non-canonical line |
| 2 | bare slug resolves to `specs/<slug>.md`; a path arg resolves as-given; a slug containing `/` gets no fallback | spec.Command | `go test ./internal/spec` red on the resolution cases | a flip that only accepts full paths (or fabricates a fallback for a slashed arg) fails these rows |
| 3 (edge: no trailing newline) | only the `Status:` line changes; every other byte, including a missing final newline, is preserved | spec.Command | `go test ./internal/spec` red on the byte-identity assertion | a rewrite that appends a line, rewrites all `Status:` lines, or normalizes the trailing newline diverges from the golden bytes |
| 4 (edge of absent/empty/re-run) | exit 1 naming file+reason on not-found, no-`staged`-line, already-implemented, and >1 staged line | spec.Command | `go test ./internal/spec` red on the error-branch rows | a flip that succeeds (or panics) on a missing/empty/already-flipped spec passes silently and corrupts state |
| 5 | gate runs; green ⇒ commit lands, red ⇒ no commit + first failure reported + nonzero | commit.Command | `go test ./internal/commit` — package absent | a commit that skips the gate lands on red; a commit that gates but ignores the verdict lands anyway |
| 6 (edge: hostile env) | an unexplained tracked-modified *or* untracked file outside the named set ⇒ block before gate, nonzero, no commit, file listed | commit.Command | `go test ./internal/commit` red on both the modified-file and untracked-file block rows | a commit missing the block-check (or gating before checking) commits a green verdict that describes files it won't land |
| 7 (edge: glob/space paths) | only named paths staged via `:(literal)`; a second untracked file never rides in; a path with a space/glob survives whole | commit.Command | `go test ./internal/commit` red on the `git add -A` and literal-pathspec rows | a `git add -A` implementation stages the stray file; a plain pathspec mis-globs a `*` path |
| 8 | `--spec` flips the spec and lands it in one commit with the named paths | commit.Command | `go test ./internal/commit` red on the `--spec` row | a commit that flips but commits separately (or forgets to add the spec path) fails the single-commit assertion |
| 9 | commit lands on the current branch including `main` | commit.Command | `go test ./internal/commit` red on the on-`main` row | a ported default-branch guard refuses on `main` and fails this row |
| 10 (edge: re-run idempotency) | named paths with no staged change ⇒ refuse empty commit, nonzero | commit.Command | `go test ./internal/commit` red on the empty-commit row | a commit that writes an empty commit on a clean re-run passes silently |
| 11 (edge: every surface) | `bin/bench.sh` routes `commit` and `spec implemented` into the core | internal/conformance route-anchor | remove one route from `bin/bench.sh`, run the conformance phase → red; restore → green | a dropped route sends a shipped surface to a dead key; the anchor is the only gate signal that a route exists |
| 12 | the route-anchor itself bites | internal/conformance | same bite-proof as row 11, recorded per `craft-gate` (remove route → red, restore → green) | an always-pass anchor would let a dropped route through; the recorded proof shows it fails red |
| 13 | the collapsed phase prose points at command tokens that route (no dangling reference) | internal/conformance stale-reference sweep | after collapse, break a pointer token → sweep red | a pointer naming a non-routing command is a dead key; cold-session legibility beyond that is review judgment (not TDD-able) |
| 14 | usage errors exit 2; operational failures exit 1 | spec.Command / commit.Command | `go test ./internal/spec` and `./internal/commit` red on the exit-code rows | a command that returns 1 for a missing `-m` (or 2 for a real block) conflates misuse with a genuine failure |

### Edge inventory

Edge classes walked per behavior, each landed above as a coverage row or below as a
**Won't handle** line:

- **error path** → rows 4, 5, 6, 10, 14.
- **empty / absent input** → row 4 (absent spec vs empty/no-staged-line spec, distinct).
- **boundary values** → row 4 (exactly-one vs zero vs >1 `Status: staged` line).
- **malformed input** → row 3 (hand-edited spec, no trailing newline); row 7 (path with space/glob).
- **interrupted / partial state** → **Won't handle:** SIGINT mid-`bench commit` —
  staging happens only after a green gate and `git commit` is effectively atomic, so
  an interrupt before that point leaves the tree exactly as found; no partial-commit
  state to clean up. A crash between stage and commit leaves staged changes the
  reviewer sees in `git status` and re-runs — the same recovery as any aborted commit.
- **re-run idempotency** → row 4 (second `spec implemented` errors); row 10 (commit
  re-run after a clean commit ⇒ empty-commit refusal).
- **hostile environment** → row 6 (unexplained working-tree files); **git missing
  from PATH** → row 5/6 fail closed: the underlying `git`/gate invocation errors
  nonzero (git is the core dependency, same posture as every `internal/git` caller),
  asserted implicitly by the operational-failure exit-1 contract. cwd deeper than
  the repo root → resolved by `git.Root` before path/spec resolution, covered by
  rows 2 and 6 running from a subdirectory fixture.
- **invocation through every surface** → row 11 (route-anchor). Symlink invocation
  is existing `route_binary` infrastructure, unchanged here — **Won't handle:** no
  new symlink logic, the shared routing already covers it.
- **unquoted multi-word arguments** → `bin/bench.sh` forwards `"$@"`, so multi-word
  paths reach the core intact; **Won't handle** as a new concern — the existing
  routing quoting is unchanged and row 7 covers a space in a path at the Command seam.

## Out of scope

None. The map is a single feature with a natural two-part build order (the Go
wrappers first, then the prose collapse that points at commands which must already
exist). Whether to slice that into two specs is the reviewer's call at sign-off —
if split, the prose-collapse story (13) and the route-anchor (11, 12) move to the
second spec; the coverage map rows partition with them.
