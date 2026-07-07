# Dashboard (FT7)

Status: staged

## Problem
Every moving piece of a Bench project is legible to an agent through the CLI —
`bench status` ranks the ambient signals, `bench roadmap` prints the working
roadmap, `bench learnings` lists open journal entries — but there is no surface a
*human* can open and read at a glance. A reviewer who wants the whole board at
once has to run several commands in a terminal and reassemble the picture in their
head. The kit ships no human-viewable snapshot of where the project stands: the
gate verdict, the pending signals, the roadmap and its recommended sequence, the
parked ideas, the open-learnings count, and the worktree pool.

## Solution
`bench dashboard` renders one self-contained static HTML page — the same data the
existing CLI surfaces already expose, laid out for a human to scan — and writes it
to `<git-dir>/bench-dashboard.html`, printing that path on stdout. `--stdout`
emits the HTML to stdout instead of writing the file.

The page is fully self-contained: inline CSS only, zero external requests, no
JavaScript. It stores nothing of its own — every call re-reads the live sources
and overwrites the file, so it is one more view over the same facts, never a
second source of truth. It never lands in the tree: `<git-dir>` is inside `.git`,
so git ignores it and it is regenerated per call.

Styling is deliberately minimal: neutral typography, a light and a dark palette
via `prefers-color-scheme`, nothing more. The richer look — animated characters,
the `ui_examples` treatment — is explicitly out of scope for v1 and recorded as a
reviewer-taste follow-up (see Out of scope). v1 is data-faithful minimalism.

## User stories

1. As a reviewer, I want `bench dashboard` to write a self-contained HTML snapshot
   to `<git-dir>/bench-dashboard.html` and print that path on stdout (exit 0), so I
   can open the whole project board in a browser from one command. Line:
   claude-opus-4-8 / medium. The path convention and the write are a new module
   boundary whose write-then-print contract the runtime test pins against the built
   binary.

2. As a reviewer piping into another tool, I want `bench dashboard --stdout` to emit
   the HTML document on stdout and write no file, so I can redirect or inspect it
   without touching the git dir. Line: claude-sonnet-5 / low. A single output-mode
   toggle the Command and the runtime contract pin exactly.

3. As a reviewer, I want a gate-verdict section showing the cached status, the tree
   sha, its age, and the same stale/red honesty `bench status` applies, so a stale
   green never reads as a clean bill. Line: claude-opus-4-8 / medium. The staleness
   comparison is the load-bearing judgment and must come from the one gate-cache
   reader `bench status` already uses, not a second parse.

4. As a reviewer, I want an ambient-signals section listing the severity-ladder
   signals in the same order and membership `bench status` produces, so the page and
   the terminal board never disagree. Line: claude-opus-4-8 / medium. Reusing the
   status snapshot rather than recomputing the ladder is the one-source constraint
   this story exists to hold.

5. As a reviewer, I want a roadmap section rendering ROADMAP.md's rows and its
   `## Recommended sequence`, so I can read what is prioritized and what comes next.
   Line: claude-sonnet-5 / low. Reuses the roadmap reader; the render is mechanical
   once the reader returns the text and the sequence.

6. As a reviewer, I want an ideas section listing the parked lines from IDEAS.md, so
   I can see the capture inbox at a glance. Line: claude-sonnet-5 / low. Reuses the
   same IDEAS.md line reader the drain counter uses.

7. As a reviewer, I want the open-learnings count shown, so I know whether a
   `/bench-what-next` drain is pending. Line: claude-sonnet-5 / low. Reuses
   `roadmap.DrainCounts` / the learnings row reader; a rendered integer.

8. As a reviewer, I want a worktree section showing the pool state — out-of-pool,
   leased, and warm entries — so I can see stray or in-flight worktrees. Line:
   claude-sonnet-5 / low. Reuses the worktree classifier the status board reads.

9. As a reviewer, I want the page to carry the timestamp it was generated at, so I
   know how fresh the snapshot is. Line: claude-sonnet-5 / low. The time is injected
   into the snapshot as data so the renderer stays a pure, deterministic function.

10. As a reviewer, I want every git-sourced and file-sourced string — commit
    subjects, branch names, paths, roadmap and idea lines — HTML-escaped and stripped
    of C0 control bytes, so markup or a terminal-escape sequence in a commit subject
    cannot inject into the page. Line: claude-opus-4-8 / medium. This is the security
    boundary of the feature; the escaping is contextual (`html/template`) and the
    control-byte sanitize is the one manual step the template does not do.

11. As a reviewer, I want the page to be one self-contained document — inline
    `<style>` only, a `prefers-color-scheme` dark palette, readable neutral
    typography, no external request and no JavaScript — so it renders identically
    offline and leaks nothing to a network. Line: claude-opus-4-8 / medium. The
    self-contained and dark-mode guarantees are assertable properties of the emitted
    string.

12. As a reviewer, I want the file written atomically (temp file then rename), so an
    interrupt mid-write never leaves a half-written page at the published path. Line:
    claude-sonnet-5 / low. A temp-plus-rename write whose leftover-temp and
    partial-file absence the runtime test observes.

13. As a reviewer who mistypes or runs outside a repo, I want `bench dashboard`
    outside a git repository to fail with the structured not-in-repo error (exit 1)
    and an unrecognized argument to fail with a usage line (exit 2), matching sibling
    commands, so a bad invocation never writes a stray file. Line: claude-sonnet-5 /
    low. Argument and repo-resolution handling is fully gate-observable at the Command
    and through the built binary.

## Implementation decisions

**New package `internal/dashboard`; one command added to the dispatch.** The
package exposes `Command([]string) (string, int)` — added to the `commands` map in
`cmd/bench/main.go` — and a pure renderer `Render(Snapshot) string`. `bin/bench.sh`
gains one `dashboard) route_porcelain "$@" ;;` routing line and one help line; the
by-path CLI, hooks, and adapters inherit that single routing, so every shipped
surface reaches the one Go implementation. No other package changes behavior; the
dashboard is a consumer of existing readers, not a new source.

**The pure renderer is the test seam.** `Render` takes a `Snapshot` value and
returns the complete HTML document as a string — no IO, no clock, no git calls
inside it. Everything time- or environment-dependent (the generation timestamp,
the gate verdict, the signal rows, the roadmap text, the ideas, the counts, the
worktree state) is gathered into the `Snapshot` by `Command` and passed in as
data. This is what lets the whole rendering — including every hostile-input case —
be tested without a browser and without a repo.

**One source per fact: the Snapshot is assembled from the existing readers.** The
gatherer composes, never re-parses:

- **Ambient signals + gate verdict** come from `internal/status`. The status
  package's per-signal computation is promoted to an exported structured accessor
  (a `Snapshot`/`Signals` function returning the severity-sorted rows as data, plus
  a shared gate-cache reader exposing status, cached tree, work tree, staleness, and
  age); `render` (the text board) is refactored to format that same accessor, so the
  board and the dashboard cannot diverge. The dashboard reads the severity ladder
  and the gate detail from there — it does not re-open `<git-dir>/bench-last-gate`
  itself.
- **Roadmap rows + recommended sequence** come from `internal/roadmap` (the same
  reader `bench roadmap` renders, and the `## Recommended sequence` extractor);
  where that extractor is currently unexported it is promoted so both callers share
  one parser.
- **Parked ideas** come from the same IDEAS.md line reader the drain counter uses.
- **Open-learnings count** comes from `roadmap.DrainCounts` / the learnings row
  reader.
- **Worktree pool state** comes from the `internal/worktree` classifier
  (`RegisteredWorktrees` and its class enum), the same one the status board reads.

Where a currently-private reader must be shared, it is promoted to its package's
exported one-source accessor and the original surface is refactored to consume the
same function — no second parser is introduced anywhere.

**HTML generation and escaping.** `Render` builds the page with `html/template`
(stdlib; no new dependency), whose contextual auto-escaping neutralizes markup and
quote injection in every interpolated field. Because `html/template` does not strip
C0 control bytes, the gatherer/renderer runs one sanitize pass that removes control
characters other than tab and newline from every git- and file-sourced string
before it reaches the template — the HTML analog of `toon.Table` refusing an
unrepresentable cell. Any template-execution error is a template-source bug, not an
input condition, so it cannot arise from repo data; the seam stays a total function.

**File location and atomic write.** `Command` resolves the git dir with
`rev-parse --absolute-git-dir` (the same call `appendGate` uses) and writes
`bench-dashboard.html` there. The write is atomic: render to a sibling temp file in
the git dir, then rename over the target, so an interrupt never leaves a partial
page at the published path and a reader never sees a half-written file. `--stdout`
skips the write entirely and returns the HTML as the command's stdout.

**Empty and error postures follow the siblings.** Each section renders a definitive
empty state when its source is absent or empty (no ROADMAP.md, empty IDEAS.md, no
gate cache, no non-trivial worktrees) rather than a blank gap or a crash — the same
show-a-definitive-empty-state posture the AXI surfaces hold. Outside a git
repository the command returns `toon.NotInRepo()` with exit 1; an unrecognized
argument returns `toon.Usage("bench dashboard", arg)` with exit 2; `--stdout` is
the only accepted flag.

**BENCH.md CLI Inventory.** `bench dashboard` is listed under "Ambient context and
capture" in `.bench/BENCH.md`'s CLI Inventory, beside `bench status` — its natural
neighbor, since it is the human-facing view of the same ambient sources.

## Testing decisions

- **What a good test is here.** Drive the pure `Render` with a constructed
  `Snapshot` and assert properties of the emitted HTML string (escaping, dark-mode
  rule, self-containment, section presence, exact timestamp from the injected
  clock); and drive the built `bench` binary for the IO behavior (file written at
  the path, path printed, `--stdout` emits and writes nothing, atomic write, exit
  codes). Never assert internal render state or pixel layout.
- **Seams and prior art.** Two seams. The renderer unit tests live in
  `internal/dashboard` (a fresh package test file, patterned on the pure-function
  tests in `internal/status`). The IO and routing behavior attaches in
  `internal/contract/runtime` as a new `runtime_dashboard_test.go`, patterned on
  `runtime_status_test.go`, which drives the built binary through a throwaway
  fixture repo. Argument-validation edges attach at `dashboard.Command`.
- **Gate command.** `.bench/gate.sh`.

### Seam diagram

    trigger: reviewer runs `bench dashboard` [--stdout]
        │
        ▼
    live sources                    ┌───────────────────────────┐
    (status ladder + gate cache,    │  Command:                 │
     roadmap + sequence, ideas,     │   gather ──▶ Snapshot     │──▶ write <git-dir>/
     learnings count, worktree,     │   Render(Snapshot)=HTML   │    bench-dashboard.html
     generation time)          ──▶  │   write | --stdout        │    + print path  (exit 0)
                                    └───────────────────────────┘    or emit HTML (--stdout)
                                                                     outside repo → err (exit 1)
                                                                     bad arg → usage (exit 2)
        Snapshot ──▶ [ Render (pure: data ──▶ HTML string) ] ──▶ self-contained document
              ◀ tests attach here: unit test feeds a hand-built Snapshot (hostile bytes,
                absent/empty sources, fixed clock) and asserts the HTML string; runtime
                contract drives the built binary for the write, --stdout, atomic write, exits

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | `bench dashboard` writes HTML to `<git-dir>/bench-dashboard.html` and prints that exact path, exit 0 | runtime contract (built binary) | New `runtime_dashboard_test.go`: assert stdout equals the git-dir path and the file exists non-empty; red on HEAD because the subcommand does not exist | A stub that prints nothing or writes elsewhere fails both the path-equals-stdout and the file-exists assertion |
| 2 | `--stdout` emits the HTML on stdout and writes no file | runtime + `dashboard.Command` unit | Assert stdout starts an HTML document and no `bench-dashboard.html` is created; red because the flag is unhandled | An impl that always writes the file (or ignores `--stdout`) fails the no-file assertion |
| 3 | Gate section shows cached status, tree sha, age, and marks a stale verdict as stale | `dashboard.Render` unit via the shared gate reader | Feed a Snapshot whose gate verdict is stale (cached tree ≠ work tree); assert the section shows the sha and a `stale` marker, not a clean green; red because `Render` is absent | A renderer that prints the cached status verbatim omits the staleness and shows a stale green as clean — the assertion on `stale` fails |
| 4 | Ambient signals render in the same ascending-severity order and membership as `bench status` | `dashboard.Render` unit consuming `status.Snapshot` | Feed a multi-signal Snapshot; assert the signals appear in severity order and none is dropped; red because `Render` is absent | A renderer that re-sorts or drops a signal diverges from the status ladder and fails the ordered-membership assertion |
| 5 | Roadmap section renders ROADMAP.md rows and the `## Recommended sequence` | `dashboard.Render` unit | Feed a Snapshot with roadmap text and a sequence; assert both appear; red because `Render` is absent | A renderer that shows only the rows and drops the sequence fails the sequence-present assertion |
| 6 | Ideas section lists the IDEAS.md parked lines | `dashboard.Render` unit | Feed parked-idea lines; assert each appears; red because `Render` is absent | A renderer that omits the ideas section fails the per-line assertion |
| 7 | Open-learnings count is shown | `dashboard.Render` unit | Feed a count of N; assert the rendered integer appears; red because `Render` is absent | A renderer that hard-codes zero or omits the count fails the exact-integer assertion |
| 8 | Worktree section shows out-of-pool / leased / warm pool state | `dashboard.Render` unit | Feed a Snapshot with each worktree class; assert each class appears; red because `Render` is absent | A renderer that shows only the root worktree fails the per-class assertion |
| 9 | The generation timestamp comes from the Snapshot's injected clock (renderer is deterministic) | `dashboard.Render` unit | Feed a fixed time; assert the exact timestamp string; red because `Render` is absent | A renderer that reads the wall clock cannot match the fixed-time assertion, exposing an impure seam |
| 10 | Every git/file-sourced field is HTML-escaped and stripped of C0 control bytes | `dashboard.Render` unit | Feed a commit subject `<script>alert(1)</script>` plus ESC/BEL bytes; assert output contains the escaped `&lt;script&gt;` form, contains no literal `<script>` tag, and no raw ESC/BEL bytes; red because `Render` is absent | A renderer that interpolates raw text injects live markup and passes control bytes through, failing the escaped-and-sanitized assertions |
| 11 | Output is one self-contained document: inline `<style>`, a `prefers-color-scheme: dark` rule, no external URL reference, no `<script>` | `dashboard.Render` unit + runtime | Assert the `<style>` block contains `prefers-color-scheme: dark`; assert output contains no `http://`, `https://`, `src=`, or `@import`; red because `Render` is absent | A renderer that pulls a CDN font or an external stylesheet adds an external URL and fails the self-containment assertion |
| 12 | The file write is atomic — a temp file renamed over the target, no leftover temp | runtime contract | After a successful run, assert no `bench-dashboard.html.tmp*` (or sibling temp) remains and the target is complete HTML; red because the subcommand does not exist | A direct truncate-write leaves a partial file on interrupt and may leave a temp artifact — the no-leftover and complete-document assertions catch it |
| 13 | Outside a repo → structured not-in-repo error, exit 1; unknown arg (`--bogus`, `--stdout` plus extra token) → usage line, exit 2, no file written | `dashboard.Command` unit + runtime | Run outside a repo → exit 1 with `not in a git repository`; `--bogus` → exit 2 usage; red because the subcommand does not exist | An impl that renders regardless of repo state, or that swallows a trailing junk token, fails the exit-code and no-file assertions |
| edge of 3,5,6,8 | Absent-vs-empty sources (no gate cache, no/empty ROADMAP.md, absent IDEAS.md, no non-trivial worktrees) each render a definitive empty state, not a blank gap or crash | `dashboard.Render` unit | Feed a Snapshot with each source absent and each present-but-empty; assert an explicit empty-state string per section; red because `Render` is absent | A renderer that nil-derefs on an absent source, or emits a silent blank section, fails the definitive-empty-state assertion |
| edge of 10 | Paths and branch names containing spaces or glob characters render intact and escaped | `dashboard.Render` unit | Feed a path with a space and a `*`; assert it appears escaped and unmangled; red because `Render` is absent | A renderer that shell-splits or drops such a value fails the intact-and-escaped assertion |
| edge of 1 | Re-running overwrites the same path with equivalent content (idempotent) | runtime contract | Run twice; assert one file at the path and the second run's non-timestamp content equals the first; red because the subcommand does not exist | A renderer that appends, or writes a differently-named file per run, fails the single-file idempotence assertion |

### Edge inventory

Edge classes walked per the profile's shell-CLI hostile-input checklist; each lands
as a coverage row above or a **Won't handle** line here.

- Control bytes (ESC, BEL) in git-sourced text — coverage row 10.
- Markup in commit subjects / roadmap / idea lines — coverage row 10.
- Paths / directory names with spaces or glob characters — coverage `edge of 10`.
- Absent file vs present-but-empty file (gate cache, ROADMAP.md, IDEAS.md,
  worktrees) — coverage `edge of 3,5,6,8`.
- Hand-edited file whose last line lacks a trailing newline — **Won't handle** as a
  new test: the roadmap and ideas readers already normalize this and are reused
  unchanged, so no new parsing surface is introduced.
- Unquoted / multi-word arguments, unknown flags, `--stdout` plus extra token —
  coverage row 13.
- SIGINT mid-write leaving partial scratch state — coverage row 12 (atomic
  temp-plus-rename write).
- Re-run idempotency — coverage `edge of 1`.
- Required tool missing from PATH (no git) / invocation through a symlink / cwd
  deeper than the repo root — **Won't handle** as new tests: repo and git-dir
  resolution reuses the shared `git.Root` / `rev-parse` helpers the sibling commands
  use, so these cases are already covered where that resolution is tested and add no
  new surface here.
- Reach through every shipped surface (kit CLI, by-path CLI, hooks, adapters) —
  **Won't handle** as a new multi-surface test: `dashboard) route_porcelain "$@"`
  forwards to the one Go implementation exactly as every other porcelain command
  does, so there is no second routing path to regress; the runtime contract exercises
  that one implementation.

## Out of scope

- **The rich visual treatment (animated characters, the `ui_examples` look, brand
  styling).** A separate design capability gated on a reviewer-owned design source
  this repo does not have (the profile records "No UI / no design source"); v1 is
  neutral minimalism by decision. Reviewer-taste follow-up, no estimate until a
  design direction is chosen.
- **Live refresh / auto-reload.** A separate capability (a watch loop or embedded
  poller) that would make the page stateful and re-introduce a background process;
  the page is a per-call snapshot by design. Estimate to add later: ~4 edits, 2 gate
  runs.
- **Serving the dashboard over HTTP.** A distinct capability (an HTTP server
  subcommand with its own lifecycle, port, and security surface); the file-plus-open
  model needs none of it. Estimate: ~8 edits, 3 gate runs.
- **Screenshot / image export.** A separate capability requiring a headless browser
  or renderer dependency, outside the self-contained-HTML premise. Estimate: not
  sized — blocked on a dependency decision (headless-browser dependency is a reviewer
  call).
