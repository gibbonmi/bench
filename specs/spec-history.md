# spec-history

Status: implemented

## Problem

`/bench-debug`'s "Finding a retired spec" section and `/bench-write-spec`'s
promote-then-delete step both point an agent at a hand-run incantation —
`git log --grep=spec-retire` plus `git log --diff-filter=D -- specs/` — to
recover a spec's origin after it has been promote-then-deleted. Every session
that needs this re-derives the same two git calls and reconciles them by eye:
duplicated knowledge, per the repo's one-source-per-fact standard.

## Solution

`bench spec history <slug>` folds both queries into one CLI call: it merges
commits that carry a `spec-retire: <slug>` message with commits that deleted
`specs/<slug>.md`, dedupes by commit, and renders one newest-first TOON table
— the FT9 (`bench diff --full`) pattern of compiling a hand-run git
incantation into the compiled Go core.

## User stories

1. As an agent in `/bench-debug`'s "Finding a retired spec" step, I want
   `bench spec history <slug>` to show me the commits that retired or deleted
   that spec, so I can recover its origin without hand-running two `git log`
   calls and reconciling them myself.
   Line: claude-sonnet-5 / medium. `bench` CLI shell plumbing at an established
   seam (FT9's compiled-core pattern) is the profile's cached cheap-to-mid
   routing, and this is the line the orchestrator already declared for the
   whole build.

2. As an agent, I want `bench spec history` to accept a bare slug or a
   `specs/<slug>.md` path, exactly like `bench spec implemented` and
   `bench spec retire` already do, so I don't have to remember a third
   argument convention.
   Line: claude-sonnet-5 / medium. Same seam, reusing `slugOf`/`specArg`
   already proven by the sibling subcommands.

3. As an agent, I want the retire-message matches and the file-deletion
   matches merged into one deduped, newest-first list (not two separate
   tables I reconcile by hand), so a commit that both deletes the file and
   carries the retirement message appears once, tagged `retire`.
   Line: claude-sonnet-5 / medium. Pure merge/dedupe logic, unit-testable
   without a git process — cheap-tier work at a known seam.

4. As an agent, I want the output to match the rest of the AXI query surface
   — a `history[N]{hash,date,kind,subject}:` TOON table, a definitive empty
   state when the slug has no history, structured stdout errors, and
   honest 0/1/2 exit codes — so this command composes with `bench diff` and
   `bench learnings` instead of inventing its own shape.
   Line: claude-sonnet-5 / medium. Conforms to the existing `toon.Table`/
   `toon.Errorf`/`toon.Usage` seam `bench diff` already established; no new
   output convention.

5. As an agent, I want hostile inputs — a commit subject with a control byte,
   a slug containing spaces or glob characters, a missing argument, and
   invocation from a subdirectory — handled the same way the rest of the spec
   subcommands and AXI surface handle them, so the command doesn't crash or
   silently misbehave on real repo history.
   Line: claude-sonnet-5 / medium. Same seam; the control-byte refusal is
   `toon.Table`'s existing guard, reused rather than reinvented.

## Implementation decisions

- New subcommand `history` added to `bench spec`'s dispatch in
  `internal/spec/spec.go`'s `Command()`. `bin/bench.sh` already routes the
  entire `spec` subcommand generically to the compiled binary
  (`route_porcelain "spec"`), so no shell-side routing change is needed —
  only the help text gains a line.
- Argument handling reuses the existing `specArg` helper (one positional,
  `-h`/`--help`, usage-error shape) and `slugOf` (bare-slug-or-path ->
  slug), so the argument convention has one source across `implemented`,
  `retire`, and `history`.
- Two `git log` queries, both anchored at the repo's top level (not the
  process cwd) via pathspec/grep so a subdirectory invocation still finds
  the whole history:
  - retire matches: `git log --fixed-strings --grep=spec-retire: <slug>
    --format=%H%x00%h%x00%cI%x00%s` — fixed-string grep so a slug
    containing regex metacharacters (`.`, `*`, `+`) matches literally, not as
    a pattern.
  - delete matches: `git log --diff-filter=D --format=%H%x00%h%x00%cI%x00%s
    -- ':(literal,top)specs/<slug>.md'` — the `literal` pathspec magic
    stops a slug containing `*`/`?`/`[` from being read as a glob, and `top`
    anchors the path at the repository root regardless of cwd.
- Merge/dedupe: each query's rows become entries tagged with their query's
  kind (`retire` or `delete`), keyed by full commit hash. Entries are merged
  retire-list-first so a hash appearing in both lists keeps the `retire` tag
  (the common case — a `bench spec retire` commit both deletes the file and
  carries the message). The merged set is sorted newest-first by the full
  ISO-8601 committer timestamp (lexical sort on ISO-8601 is chronological);
  the rendered `date` column is that timestamp's first 10 characters
  (`YYYY-MM-DD`), keeping the row narrow per the AXI minimal-schema
  principle while sorting on full precision internally.
- Rendering reuses `toon.Table("history", []string{"hash","date","kind","subject"},
  rows)` exactly as `bench diff`'s log table does; a control-byte subject
  makes `toon.Table` refuse, and the command surfaces that as
  `toon.RenderError(err)` at exit 1 — the same posture `bench diff --full`
  already has for the same failure, not a new one.
- No file-existence check on `specs/<slug>.md`: the whole point of the
  command is recovering a spec that no longer exists in the working tree.

## Testing decisions

- Good tests here exercise the command at the shell seam (`bench spec
  history <slug>` through the built binary against a throwaway git fixture),
  matching the existing `internal/contract/runtime/runtime_spec_test.go`
  pattern for the sibling `spec` subcommands, plus a package-level unit
  seam for the pure merge/dedupe logic (no git process), matching
  `internal/spec/spec_test.go`'s existing style of testing `Resolve`/`Flip`
  directly.
- Seams tested:
  - `internal/spec` package tests: the merge/dedupe function directly, fed
    synthetic rows — no git subprocess, no filesystem.
  - `internal/contract/runtime/runtime_spec_test.go`: `bench spec history
    <slug>` through the built `dist/bench` binary against real throwaway git
    repos — this is where the git integration, TOON shape, hostile inputs,
    and exit codes actually get proven.
- Gate: the project gate, `.bench/gate.sh` (root conformance suite +
  runtime/behavior contracts). `docs_workflow_helpers_test.go`'s anchor
  checks additionally require the substrings `spec-retire:` (in
  `bench-write-spec.md`) and `diff-filter=D` (in `bench-debug.md`) to survive
  byte-for-byte — the new prose in those two files is additive, verified by
  keeping those substrings before editing and confirming the gate stays
  green after.

### Seam diagram

    trigger: an agent debugging a retired feature, or /bench-debug's
             "Finding a retired spec" step
        │
        ▼
    <slug or specs/<slug>.md>
        │
        ▼
    [ bench spec history <slug> ]
        │           │
        ▼           ▼
    git log       git log
    --grep        --diff-filter=D
    (retire)      (delete)
        │           │
        └─────┬─────┘
              ▼
        [ merge + dedupe, newest-first ]
              │
              ▼
    history[N]{hash,date,kind,subject}: TOON table
    (or history[0]{...} — definitive empty state)

    ◀ tests attach here: `bench spec history <slug>` through the built
      binary against a throwaway git fixture (runtime contract); the
      merge+dedupe step alone via a direct package-level unit test.

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | `bench spec history <slug>` lists commits that retired or deleted a spec, newest first | runtime contract (`bench spec history`) | `go test ./internal/contract/runtime -run TestRuntimeSpecHistoryContracts` fails before `history` exists (unknown subcommand, exit 2) | a missing or degenerate (always-empty) implementation never lists the real commits, so the fixture's known retire/delete commits would be absent from stdout |
| 2 | bare slug and `specs/<slug>.md` path both resolve the same history | runtime contract | same test, path-vs-slug case | a degenerate implementation that only accepts one argument form fails the other |
| 3 | a commit that both deletes the file and carries the retire message appears exactly once, tagged `retire` | `internal/spec` unit test (merge/dedupe) + runtime contract | `go test ./internal/spec -run TestMergeHistoryDedupes` fails red before the dedupe logic exists (naive concatenation would double the row) | proves the dedupe collapses the common case instead of listing the same commit twice |
| 3 (edge) | a delete-only commit (file removed without a `spec-retire:` message) is tagged `delete`, not `retire` | `internal/spec` unit test | same test file, delete-only case | a degenerate always-`retire` implementation would mislabel it |
| 4 | empty history renders `history[0]{hash,date,kind,subject}:` (definitive empty state) | runtime contract | fixture with a slug that has no matching commits; asserts the exact `history[0]{...}` line | a silent/blank-stdout implementation would look indistinguishable from a crash |
| 4 | TOON header names all four fields; row count matches the fixture's known commit count | runtime contract | same suite, asserts `history[N]{hash,date,kind,subject}:` and exact row lines | a wrong field set or wrong count is the shape regression this line pins |
| 4 | missing argument exits 2 with a usage message, before any git call runs | runtime contract | `bench spec history` (no arg) — already covered by `specArg`'s existing usage path, exercised fresh for `history` | a build that skips argument validation would instead attempt a git call with an empty slug |
| 4 | unknown subcommand / extra positional / unknown flag exit 2 | runtime contract | same suite, args-table case | pins `specArg`'s existing usage contract holds for the new subcommand too |
| 5 | a commit subject containing a control byte (ESC) makes the command exit 1 with the unrepresentable-TOON-cell error, not a mangled row | runtime contract | fixture commits a control-byte subject via the retire-grep path; asserts `error: unrepresentable TOON cell` at exit 1 | proves `toon.Table`'s existing refusal is reached, not silently bypassed by a hand-built string instead of the shared renderer |
| 5 | a slug containing spaces or glob characters (`*`) still finds its exact commit, not a different one | runtime contract | fixture with `specs/weird name.md` and `specs/weird*name.md`, each retired; asserts each slug's history contains only its own commit | a naive (non-literal) pathspec would glob-match the wrong file or error out |
| 5 | invocation from a subdirectory still finds the full-repo history | runtime contract | `bench spec history` run with cwd set to a subdirectory; asserts the same rows as from the repo root | a cwd-relative (non-`:(top)`) pathspec would silently return nothing from a subdirectory |
| 5 | a slug with no history at all (never existed, never retired) still exits 0 with the empty state, not an error | runtime contract | fixture with an unrelated repo history; asserts `history[0]{...}` at exit 0 | distinguishes "definitively no history" from "command failed" — an error posture here would be a false negative on every clean repo |

### Edge inventory

- error path — an invocation outside a git repo: **row** — `toon.NotInRepo()` at
  exit 1, reusing the existing not-in-repo posture every other `bench spec`
  subcommand and the AXI surface already share.
- empty/absent input — slug with no matching commits: **row** (definitive
  empty state, above).
- boundary values — a slug whose only match is the very first commit in the
  repo (no parent): **Won't handle** — `git log --grep`/`--diff-filter=D`
  don't need a parent to match a commit's own message or its own diff
  against an empty tree, so this is not actually a boundary the command
  special-cases; no code path assumes a parent exists.
- malformed input — slug containing shell-meta characters (`;`, `$(`),
  spaces, glob characters: **row** (above) — arguments reach `exec.Command`
  as a single argv element, never through a shell, so shell metacharacters
  are inert; the glob/literal-pathspec case is the one that needed a
  deliberate choice and is covered.
- interrupted/partial state — n/a: the command is read-only (no writes, no
  partial-completion state to corrupt).
- re-run idempotency — running `bench spec history <slug>` twice in a row
  produces byte-identical output on an unchanged repo: **Won't handle** as a
  dedicated row — it falls directly out of the command being a pure read
  over git history with no side effects; the runtime contract's repeated
  assertions across the same fixture already exercise this implicitly.
- hostile environment — required tool missing from PATH (no `git`): **Won't
  handle** — every other `bench spec`/AXI command shares this same
  unhandled dependency; not a regression this feature introduces, and no
  existing sibling command carries a dedicated row for it either.
- control bytes in git-sourced text (profile checklist): **row** (above).
- paths/slugs with spaces or glob characters (profile checklist): **row**
  (above).
- invocation through every shipped surface (profile checklist): **Won't
  handle** as a separate row — `bin/bench.sh`'s generic `spec)
  route_porcelain "$@"` already routes every subcommand under `spec`
  uniformly; there is no `history`-specific routing path to diverge, so
  this is structurally covered rather than a new case to test.
- cwd deeper than repo root (profile checklist): **row** (above).

## Out of scope

- **A `--full` flag that appends the raw diff body of the deletion commit** —
  a separate capability (mirrors `bench diff --full`'s own append-on-flag
  design, but for a different range shape: one historical commit vs. a
  branch-relative range) — `2 edits, 2 gate runs` to add if a future session
  needs to read the deleted spec's content inline rather than via
  `git show <hash>:specs/<slug>.md` by hand.
- **Recovering a spec's content automatically (e.g. `bench spec history
  --restore <slug>`)** — a distinct write capability with its own safety
  questions (does it re-stage the file? overwrite a same-named live spec?),
  not the read-only recovery this spec scopes — `5 edits, 4 gate runs`
  estimate if ever requested.
