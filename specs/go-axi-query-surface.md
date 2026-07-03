# go-axi-query-surface — slice 2 of the Go rewrite

Status: implemented

Map: `decisions/go-rewrite.md` (closed). This spec builds slice 2 of the
8-slice strangler: the shared flat-table TOON emitter and all five read-only
AXI query subcommands — `learnings`, `maps`, `guards`, `diff`, `coverage` —
ported from shell to the Go core, behind the existing strangler router. It is
the first slice that moves real parsing logic, so it sets the `internal/`
package idioms every later slice copies.

## Problem

Slice 1 proved the pipeline: a Go module, the launcher, the strangler router,
npm packaging, and the gate's Go layer, with exactly one trivial subcommand
(`version`). No real logic has moved. The AXI query surface — the agent-facing
read-only half of the hybrid output contract — is still shell: a shared TOON
emitter plus five parsers (`bin/bench-query.sh`, `bin/bench-diff.sh`,
`bin/bench-coverage.sh`). These are exactly the parsers/emitters the map
committed to Go because shell can't unit-test them: the canary layer stands in
for the missing test seam, the `LEARNINGS_OPEN_RE` and `maps_awk_prelude`
knowledge is single-sourced only by hand-discipline, and every hostile-input
edge is re-litigated in a gate fixture instead of a table test.

The emitter is the binding constraint. `toon_escape`/`toon_table` are shared by
all five commands (`bench-diff.sh:58`, `bench-coverage.sh:121` both pipe into
`toon_table`). Porting the emitter to Go without its five consumers would leave
the TOON escaping rule implemented in two languages at once — the two-source
defect the migration exists to end. So the emitter and its five consumers move
as one slice.

## Solution

One Go package per parser under `internal/`, composing a shared `internal/toon`
emitter, dispatched by new subcommand cases in `cmd/bench`. The shell router
flips `learnings|maps|guards|diff|coverage` to `route_binary` (the seam slice 1
built for `version`); the three shell libraries are deleted. The existing
black-box AXI gate contracts — which drive `bench.sh <cmd>` and assert
stdout/exit — run unchanged against the Go binary and are the regression net.
The parsers/emitter gain `go test` table tests (the new unit layer), and the
three gate contracts that reached into shell internals move down to that layer.

`bench status` (slice 3, still shell) keeps working: its two counts are sourced
from the Go binary through thin adapters, so the parsing knowledge lives once —
in Go — and the dashboard's numbers cannot drift from what `bench maps` and
`bench learnings` report.

## User stories

1. As a kit developer, I want a shared `internal/toon` emitter (flat-table
   `name[N]{fields}:` header + escaped, comma-joined rows) with the exact
   escaping rule — quote a field containing a comma, double-quote, newline, or
   leading/trailing whitespace, doubling inner quotes — so that all five
   commands compose one TOON implementation and the escaping rule has a single
   source. Line: claude-opus-4-8 / medium. This is the first `internal/`
   package and every command in this slice and later slices composes it, so it
   sets the package and table-test idioms and warrants the mid tier despite a
   small, pure surface.

2. As a bench user, I want `bench learnings` to emit open journal headings as a
   `learnings[N]{date,title}:` TOON table (date + separator-run stripped, CRLF
   stripped, definitive empty state, exit 0; usage on stdout + exit 2 for an
   unknown argument; structured error + exit 1 outside a git repo), identical to
   the shell command it replaces, so that the ported command is a drop-in.
   Line: claude-sonnet-4-6 / medium. The parser is mechanical and the existing
   AXI contract fully pins the observable output, which is the cheap-tier case.

3. As a bench user, I want `bench maps` to emit unresolved decision-map tickets
   and the close-readiness handoff row from one `internal/maps` engine — the
   marker/fence/CRLF/`## Handoff` close-readiness logic that feeds both the row
   listing and the distinct-file count as a single source — so that the ticket
   listing and the count status consumes cannot drift. Line: claude-opus-4-8 /
   medium. The awk close-readiness prelude is the most intricate parser in the
   slice and its one-source-for-two-outputs shape is load-bearing, so it takes
   the mid tier even though the gate fully observes it.

4. As a bench user, I want `bench guards` (and `--brief`) to discover every
   guard by convention, exec each one's `--describe` under a 5-second bound
   (never hanging, never running an unmanaged pre-push), and aggregate the
   manifests into a `guards[N]{guard,boundary,denies}:` table, identical to the
   shell command, so that the deny surface stays learnable. Line:
   claude-opus-4-8 / medium. It sets the exec-a-shell-guard idiom later slices
   reuse and its timeout and never-run-a-foreign-hook branches are safety-
   critical, so the mid tier is justified despite full gate coverage.

5. As a reviewer, I want `bench diff` to resolve the review base (recorded
   `branch.<name>.benchBase` when it names a reachable ancestor, else merge-base
   with the default branch, naming which method resolved and falling back loudly
   on an unreachable or divergent recorded sha) and emit changed files as a
   `files[N]{status,path}:` table with paths TOON-escaped exactly once, identical
   to the shell command, so that the review phase has one source of base truth.
   Line: claude-sonnet-4-6 / medium. Every base-resolution branch is pinned by
   the wave-2 contract, so the fully-observed port routes cheap.

6. As a spec author and the gate's docs layer, I want `bench coverage <spec>`
   (extraction: `spec`/`state`/`rows` TOON) and `bench coverage --check`
   (validation: the canonical header, five-cell rows, non-empty cells, in-range
   story references, historical opt-out — each with its exact, load-bearing
   error phrasing) ported to one `internal/coverage` parser identical to the
   shell command, so that the acceptance-coverage convention keeps one
   validator. Line: claude-opus-4-8 / medium. The validation phrasings are
   consumed by substring elsewhere and a drift silently breaks a downstream
   consumer, so the fiddly rule set takes the mid tier.

7. As a kit maintainer, I want the strangler router to send
   `learnings|maps|guards|diff|coverage` to the Go binary through the existing
   `route_binary` seam and the three shell libraries
   (`bench-query.sh`, `bench-diff.sh`, `bench-coverage.sh`) deleted with their
   `source` lines removed, so that no ported command has two live
   implementations. Line: claude-sonnet-4-6 / medium. Mechanical dispatch edits
   at the one seam slice 1 already built.

8. As a bench user, I want `bench status` to keep reporting its open-learnings
   and unresolved-maps signals with byte-identical output, sourcing both counts
   from the Go binary (learnings count from the `learnings[N]` header;
   unresolved-maps count as the distinct maps the ported `bench maps` lists) so
   that the dashboard's numbers have the same single source as the query
   commands. Line: claude-opus-4-8 / medium. This is the coexistence crux of the
   slice — a wrong adapter either re-introduces the two-derivations bug or
   silently mis-renders the dashboard — so it takes the mid tier despite a small
   diff.

9. As a kit maintainer, I want the three AXI gate contracts that reached into
   shell internals (`toon_escape` leading/trailing-space quoting, the
   `maps_unresolved_count` distinct-file figure, and the close-readiness count)
   re-homed to `internal/` table tests covering the identical edges, so that no
   contract sources a deleted shell file and the pure-helper coverage survives
   the port at its proper seam. Line: claude-opus-4-8 / medium. Editing the
   oracle is the worst defect class in this kit (`craft-gate`), so preserving
   exact coverage while moving the seam takes the mid tier.

10. As the teammate who just walked in, I want the docs that describe the CLI
    file layout (`README.md`'s tree, and any prose naming the shell parser
    internals) updated to the Go package structure, so that the repo map matches
    the repo. Line: claude-sonnet-4-6 / medium. Mechanical doc edits verified by
    the docs gate's stale-reference sweep.

## Implementation decisions

- **Package layout.** One package per parser under `internal/`: `toon`
  (emitter + the shared structured-error/usage helpers `axi_error`/`axi_usage`
  become `toon.Errorf`/`toon.Usage`, one source for the `error: … — …` and
  `usage: … (unknown argument: …)` formats), `learnings`, `maps`, `guards`,
  `diff`, `coverage`. `cmd/bench/main.go` gains one dispatch case per command
  that resolves arguments, calls the package, and maps results to
  stdout + exit code. Idioms set here (package boundary at the parser, table
  tests beside each parser, git and guard invocation via `os/exec`) are what
  slices 3–8 copy.

- **The emitter takes structured rows, not TAB-encoded lines.** The shell
  emitter round-tripped records as TAB-joined strings because it piped text;
  Go parsers return `[][]string` and hand them to `toon.Table` directly. The
  TAB intermediate was a shell artifact and does not survive the port. The
  escaping rule and the `name[N]{fields}:` framing are byte-identical to the
  shell output the contracts assert.

- **git and guard invocation stay subprocess.** The Go commands shell out to
  `git` (`rev-parse --show-toplevel`, `config branch.<n>.benchBase`,
  `merge-base`, `symbolic-ref`, `diff --name-status --no-renames -z`) and to
  each guard script's `--describe`, exactly as the shell did — git and the
  guard scripts stay the source of that truth. `bench guards` uses
  `exec.CommandContext` with a 5-second timeout for the bound the shell got
  from coreutils `timeout`/a watchdog; this removes the coreutils dependency,
  an incidental improvement, not a contract change. The hooks-dir resolution
  the shell `guards_rows` re-derived (a parked one-source finding) is written
  once in `internal/guards`; the port dissolves that duplication rather than
  carrying it.

- **Errors and usage print to stdout.** Per the AXI hybrid contract, structured
  errors and usage go to stdout (not stderr) with exit 1 and 2 respectively;
  the Go helpers preserve this. `route_binary` `exec`s the binary, so the Go
  exit code becomes `bench`'s, as with `version`.

- **The router flip is additive at one seam.** Slice 1's `route_binary` takes
  the subcommand as `$@`; the five cases change from `shift; <fn> "$@"` to
  `route_binary "$@"` (no `shift`, matching `version`). The three `.
  "$BENCH_BIN_DIR/bench-*.sh"` source lines are removed and the files deleted.
  No new resolver, no second mechanism — the map's constraint that the router
  grows names, not mechanisms.

- **Status sources its counts from the Go binary (the coexistence contract).**
  `status()` stays shell in this slice. Its two call sites keep their
  signatures; `learnings_open_count` and `maps_unresolved_count` move from the
  deleted `bench-query.sh` into thin adapters (co-located with `status()` in
  `bench-status.sh`) that invoke the Go binary against the same root and derive
  the figure from its TOON output: the learnings count is the `learnings[N]`
  header value; the unresolved-maps count is the number of distinct `map`
  values `bench maps` lists (a file appears iff it is not close-ready, so
  distinct-file count equals the ported engine's figure). The close-readiness
  **rule** lives once, in `internal/maps`; the adapter only counts emitters. No
  shell copy of the parser survives — that would be the two-source defect this
  slice removes.

- **Three gate contracts re-home to `go test`; no canary changes.** The
  five shell-fixture AXI canaries (`toon-escaping-dropped`,
  `learnings-parse-broken`, `guards-aggregation-dropped`,
  `coverage-extraction-dropped`, `diff-recorded-base-dropped`) carry
  self-contained broken shell CLIs and trip black-box CLI contracts that run
  the fixture's `bench.sh` and assert output; those contracts are unchanged by
  the port, so the canaries keep biting untouched (a property the build must
  verify, not assume). The only gate edits are the three contracts in
  `gate-axi-contracts.sh` that `. bin/bench-query.sh` to test pure helpers —
  they have no canaries and would fail on the deleted file — which move to
  `internal/toon` and `internal/maps` table tests covering the identical edges.

- **Deviations from the map: none.** Slice 2 is "AXI query surface + TOON
  emitter"; the profile defines that surface as these five commands plus the
  emitter, and the map's Handoff item 4 names "`go test` table tests on
  parsers/emitters/pure logic" as the new unit layer and "existing gate
  fragments run unchanged against the ported binary" as the net — this spec is
  that plan. No slice-2 uncertainty flags remain in the Handoff (items 7's
  flags were slice-1 or slice-7 scoped), so no top-tier spec review was spawned.

## Testing decisions

- **What a good test is here.** Acceptance drives the public `bench` entry and
  asserts stdout/stderr text, exit codes, and TOON shape — never Go internals.
  The existing AXI contracts already do exactly this and carry over. Go table
  tests are additional, at the pure-function seam (parsers/emitter), and are
  where the canary layer's compensation for shell untestability retires.
- **Seams.** Three: the **AXI CLI acceptance seam** (the existing contracts in
  `gate-axi-contracts.sh`, `gate-axi-guards-contracts.sh`,
  `gate-axi-wave2-contracts.sh`, driving `bench.sh <cmd>` → Go binary — prior
  art: they exist and pass today against shell); the **`go test` unit seam**
  (table tests beside each `internal/` package — prior art:
  `cmd/bench/main_test.go`'s `versionLine` table); and the **status
  coexistence seam** (the status contract in `gate-runtime-contracts.sh`,
  driving `bench.sh status` and asserting the signal rows — prior art: its
  learnings-signal assertion at `gate-runtime-contracts.sh:233`).
- **Gate command:** `.bench/gate.sh` (the project gate), whose Go layer already
  runs `gofmt -l`, `go vet ./...`, `go test ./...`, `go build`, and the four
  cross-compile targets — so new `internal/` packages and their tests are
  graded with no new gate wiring.

### Seam diagram — AXI CLI acceptance seam (all five commands)

    trigger: gate contract (or a user/hook) runs `bench <cmd> [args]`
        │
        ▼
    argv ─────────▶ [ bin/bench.sh: route_binary ] ──▶ exec dist/bench
    repo state ───▶ [ cmd/bench dispatch → internal/<pkg> ]
    (git, files)    [ compose internal/toon emitter ] ──▶ stdout: `<name>[N]{fields}:` + rows, exit 0
                    [                              ] ──▶ stdout: `usage: …`, exit 2
                    [                              ] ──▶ stdout: `error: … — …`, exit 1
        ◀ tests attach here: the existing AXI contracts run `bash bench.sh <cmd>`
          in a fixture repo and assert stdout shape + exit code, unchanged from
          the shell surface. Red before the Go subcommand exists: route_binary
          reaches a binary that answers `unknown subcommand` (exit 2), so the
          contract's expected TOON is absent.

### Seam diagram — `go test` unit seam (parsers + emitter)

    trigger: gate Go layer runs `go test ./...`
        │
        ▼
    field / heading / ──▶ [ internal/toon.Escape/Table      ] ──▶ TOON string
    map file / cells      [ internal/{learnings,maps,       ] ──▶ rows [][]string
    (in-memory or         [   coverage}.Parse               ] ──▶ count / violations
     temp-dir fixtures)   [ internal/{diff,guards} helpers  ]
        ◀ tests attach here: table tests assert Escape quotes leading/trailing
          whitespace and comma/quote/newline; maps engine returns the distinct-
          file count and close-readiness rows; coverage returns each violation
          phrasing. Red before the parser exists: the package/function is absent
          and the test does not compile or fails on empty output.

### Seam diagram — status coexistence seam

    trigger: gate status contract (or session-start) runs `bench status`
        │
        ▼
    repo state ──▶ [ bench-status.sh status() ] ──▶ ranked signal rows
                   [   learnings adapter  ] ──(cd root; bench learnings)──▶ [N] header
                   [   maps adapter       ] ──(cd root; bench maps)──▶ distinct map count
        ◀ tests attach here: the status contract sets up an open learning and an
          unresolved map, runs `bench status`, and asserts the `learnings` and
          `decisions` rows carry the right counts. Red before the adapters are
          wired: the deleted shell functions are unresolved, so status errors or
          omits the rows.

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | `internal/toon.Escape` quotes a field with leading/trailing whitespace and with comma/quote/newline, doubling inner quotes; leaves plain fields bare | go test unit | `go test ./internal/toon` before `Escape` exists → does not compile | any missing branch (esp. the whitespace branch the CLI contracts never exercise) breaks a table row |
| 1 | `internal/toon.Table` emits `name[N]{fields}:` + one escaped, comma-joined, two-space-indented row per record; empty input → `name[0]{fields}:` | go test unit + AXI CLI | table test before `Table` exists → fails; every CLI contract's header/row assert | a wrong count, indent, or join breaks the exact-string asserts across all five commands |
| 2 | `bench learnings` → identical `learnings[N]{date,title}:` TOON (date/separator/CRLF stripping, template/empty → 0 rows, escaping) | AXI CLI | "AXI learnings two-entry", "empty/template", "field-escaping", "ascii-separator" contracts before the Go subcommand → `unknown subcommand`, expected TOON absent | the four existing learnings contracts pin every parser edge against the binary |
| 2,4,5,6 | unknown argument → `usage` on stdout, exit 2; outside a repo / missing input → `error:` on stdout, exit 1 | AXI CLI | the usage/error-posture contracts before the Go subcommand → wrong exit / no usage line | honest exit codes are the AXI contract; each command's contract asserts them |
| 3 | `bench maps` → one row per unresolved ticket, `unknown` for a Type-less ticket, anchored/fence-aware/CRLF-safe placeholder detection, resolved tickets excluded | AXI CLI | "AXI maps unresolved-ticket", "over-match anchoring/fence", "CRLF-stripping", "no-Type-ticket" contracts before the Go subcommand → expected rows absent | the four maps contracts pin the parser against the two-derivations edge classes |
| 3 | close-readiness handoff row: zero-open map emits `<map>,handoff,handoff,{missing\|state}`; filled Handoff silent; open-ticket map emits no handoff row; non-map file never nagged | AXI CLI | "AXI maps handoff close-readiness contract" (its `bench.sh maps` half) before the engine → handoff rows absent/wrong | the close-readiness engine's row output is asserted per case at the CLI |
| 3,8 | the distinct-file count `status` consumes equals the ported engine's figure (two tickets in one file → 1) | go test unit | `internal/maps` count test before the engine → fails; the close-readiness contract's count assertion re-homed here | one engine feeds rows and count, so the two cannot drift; the table test pins the count the CLI rows imply |
| 4 | `bench guards` → one TOON row per deny-capable guard, session-start excluded, stub → `no manifest`, absent pre-push → `not installed`, foreign pre-push never executed, slow `--describe` → `no manifest (timed out)` under the bound | AXI CLI | the five `gate-axi-guards-contracts.sh` contracts before the Go subcommand → expected rows/timing absent | the guards contracts pin discovery, exclusion, the exec bound, and the never-run-a-foreign-hook safety |
| 4 | `bench guards --brief` → one line per deny-capable guard + exactly one footer; the surface session-start injects | AXI CLI + status/session | "AXI guards --brief" and "session-start guard-brief injection" contracts before the Go subcommand → brief absent | the brief is consumed by the session-start hook; its contract asserts the injected shape |
| 5 | `bench diff` base resolution: recorded reachable-ancestor base wins (`method: recorded`); unreachable, divergent, and no-key cases fall back to merge-base naming the reason; unresolvable → exit 1 | AXI CLI | "AXI diff recorded-base" and "diff fallback/shape" and "diff error-posture" contracts before the Go subcommand → base/method lines absent | the wave-2 diff contracts pin every base-resolution branch and the loud-fallback preamble |
| 5 | changed files → `files[N]{status,path}:` with space/non-ASCII/quote paths TOON-escaped exactly once (not git-quoted then re-quoted); empty diff → `files[0]…`, exit 0 | AXI CLI | the "diff fallback/shape" contract's path rows before the Go subcommand → rows absent/double-escaped | the contract asserts raw `-z` paths escaped a single layer |
| 6 | `bench coverage <spec>` extraction: `spec`/`state`/`rows` TOON for mapped, historical, no-map; CRLF-safe; hostile cells escaped | AXI CLI | "AXI coverage extraction" and "coverage state/error" contracts before the Go subcommand → preamble/rows absent | the wave-2 coverage contracts pin extraction and state classification |
| 6 | `bench coverage --check` emits each validation violation with its exact phrasing (missing header, no data rows, cell count, empty cell, out-of-range / unrecognized story ref); valid/historical/no-map silent, exit 0 | AXI CLI | "AXI coverage --check validation" contract before the Go subcommand → phrasings absent | the contract matches each phrasing by substring, as does the docs-gate consumer |
| 7 | `learnings\|maps\|guards\|diff\|coverage` route to the Go binary; the three shell libs are gone and unsourced | AXI CLI (all above) + parse layer | any AXI contract run while the case still calls the deleted shell fn → function-not-found; `bash -n`/gate load error if a `source` line points at a deleted file | the whole AXI suite plus the gate's parse layer fail if a command still routes to shell or a dangling source remains |
| 8 | `bench status` emits the open-learnings signal row with the right count, sourced from the Go binary | status coexistence | `gate-runtime-contracts.sh:233` learnings-signal contract after the shell fn is deleted and before the adapter is wired → row missing | the existing status contract asserts the `/bench-integrate-learnings` row appears for an open learning |
| 8 | `bench status` emits the unresolved-maps signal row with the right count (**new** contract — this signal is unasserted today) | status coexistence | a new status contract setting up an unresolved map, run before the maps adapter is wired → `N unresolved map(s)` row absent/wrong | closes the pre-existing gap where status's maps count had no end-to-end assertion; catches an adapter that miscounts distinct files |
| 9 | the leading/trailing-space escaping, distinct-file count, and close-readiness count edges keep coverage after leaving the shell-source contracts | go test unit | the re-homed `internal/toon`/`internal/maps` table tests before they exist → fail/don't compile | proves the pure-helper coverage moved rather than vanished when the internal-sourcing contracts were removed |
| 10 | `README.md`'s CLI tree and any prose naming the deleted shell parsers match the Go layout | docs gate | the docs gate's stale-reference / CLI-command sweep against a tree that still names `bench-query.sh` etc. → red | the conformance layer fails when docs reference files that no longer exist |

### Edge inventory

Walked per behavior against the profile's shell-CLI hostile-input checklist;
each resolved as a coverage row above or a **Won't handle** line here.

- **paths/dirs with spaces or globs** — covered: the AXI path-with-spaces
  contracts (`--space-path`) for learnings/maps/guards, and diff's `a b.txt`
  row, run against the Go binary; Go's `os/exec` passes argv, never a shell
  string.
- **absent vs present-but-empty file** — covered: learnings empty/template
  contract and maps absent-`decisions/` contract; the maps engine's fence/empty
  handling is table-tested.
- **hand-edited files without a trailing newline** — covered by the parsers'
  line handling; the maps/learnings table tests include a no-trailing-newline
  fixture (profile checklist item, asserted in `internal/` tests).
- **CRLF / carriage returns** — covered: the maps CRLF contract and coverage
  CRLF contract at the CLI, plus a `\r`-bearing row in the toon/maps table
  tests.
- **unquoted multi-word args (`$*` vs `$1`)** — dissolves: Go receives argv as
  `[]string`; the usage/exit-2 contracts assert the single-unknown-arg path.
- **required tool missing / off-matrix platform** — covered by slice 1's
  version-routing seam (unchanged); these five commands reuse `route_binary`,
  so the missing-binary (exit 127) and unsupported-platform (exit 2) paths are
  already proven and not re-tested here.
- **invocation through a symlink; cwd deeper than root** — covered: the AXI
  subdirectory root-resolution contract (runs from `sub/deeper`); symlink
  resolution is slice 1's router property, unchanged.
- **guard `--describe` that hangs or is foreign** — covered: the guards
  timeout-bound and unmanaged-pre-push safety contracts, via
  `exec.CommandContext`.
- **status called from a non-root cwd (worktree, subdir)** — covered by the
  adapter `cd`-ing to the resolved root before invoking the Go binary, so the
  count resolves the same repo status intends; the status contract exercises a
  subdir where prior art does.
- **Won't handle: a repository path containing a literal newline** — the shell
  `diff` already documents this as the one shape its `-z` reader misreads; the
  Go port preserves parity (behavior unchanged), and no in-scope caller can
  produce it. Fixing it is not this slice.
- **Won't handle: `bench status` performance from two added subprocess spawns**
  — status gains two Go-binary invocations per call; the binary runs in
  milliseconds and session-start already spawns many git subprocesses, so this
  is not optimized here. Slice 3 (status port) collapses both into in-process
  Go calls.

## Out of scope

- **`bench status` renderer port** — slice 3 of the map: the deterministic
  dashboard renderer, severity ladder, gate-cache read, and the roadmap footer
  move to Go as their own capability. This slice only adds the two count
  adapters so the shell renderer keeps working. Estimate to build later: ~20
  edits, ~8 gate runs.
- **Slices 4–8** — `git-guard.py` absorption, hook logic behind shims,
  `doctor`/`link`, worktree/shift loop, and the gate-fragment port. Each is a
  distinct capability with its own spec by the map's dependency order; the
  shift/worktree contract backfill (Handoff item 7) precedes slice 7.
- **`gate_tree_hash` de-duplication** (mirrored verbatim in
  `bin/bench.sh` and `.bench/hooks/stop.sh`) — a parked one-source finding that
  belongs with the stop-hook port (slice 5/7), not here; `internal/maps`
  touches nothing in that mirror. Estimate: ~4 edits, ~2 gate runs.
