# go-status-port — slice 3 of the Go rewrite

Status: implemented

Map: `decisions/go-status-port.md` (closed; child of `decisions/go-rewrite.md`).
This spec builds slice 3 of the 8-slice strangler: the whole of
`bin/bench-status.sh` ported to the Go core and the file deleted. That is the
status renderer plus every helper entangled with it — the structure checker, the
merged-spec retirement counter, the gate tree-hash, the worktree pool/lease path
facts, and the trivial `models`/`idea`/`roadmap` riders. Two live duplication
defects the parent map named die in this slice: the `gate_tree_hash` mirror in
`.bench/hooks/stop.sh`, and the worktree pool/lease facts otherwise owned by
`bin/bench-worktree.sh` (slice 7).

## Problem

`bin/bench-status.sh` (329 lines) is the last big block of executable parsing and
rendering still in shell. It carries the ambient dashboard the session-start hook
consumes verbatim, the structure-budget checker the shift loop gates on, and four
support parsers. Slice 2 left it as shell that sources its two counts from the Go
binary through thin adapters; the renderer, the structure engine, and the
tree-hash still live in bash.

The file is entangled across three future slices, which is why the map settled its
ownership now rather than porting a clean sub-unit:

- `gate_tree_hash` has three consumers — status (gate-cache read), the shift
  loop's gate-cache write (`gate_record` in `bin/bench.sh`), and a **hand-mirrored
  copy** in `.bench/hooks/stop.sh` that cannot source `bench.sh`. The mirror is
  one of the two live two-source defects the parent map named.
- the structure checker (`structure_check`) feeds three surfaces: `bench
  structure`, `status`'s violation count, and the shift loop's touched-scope
  refactor gate (`structure_touched_since`).
- status's worktree row reads the pool-path and lease-file conventions owned by
  `bin/bench-worktree.sh` (slice 7), so status can't render its worktree signal
  without either duplicating those two facts or inverting their ownership.

Left in shell, these stay untestable except through the black-box gate, the
tree-hash rule stays implemented in two languages, and the worktree facts can't
move to Go without a duplication window held open four slices.

## Solution

One Go package per unit under `internal/`, dispatched by new subcommand cases in
`cmd/bench`, composed by a `internal/status` renderer. The shell router flips
`status|structure|models|idea|roadmap` to `route_binary` (the seam slice 1 built
and slice 2 reused), adds three plumbing subcommands (`tree-hash`,
`worktree-pool`, `worktree-lease-file`), and `bin/bench-status.sh` is **deleted**.

Ownership inverts as the map decided (#2, #3): the tree-hash and the worktree
pool/lease path facts become Go's, and the remaining shell callers — `gate_record`
and `structure_touched_since` in `bin/bench.sh`, the pool/lease helpers in
`bin/bench-worktree.sh`, and `record_gate` in `.bench/hooks/stop.sh` — become
one-glance adapters that call the binary. No fact is implemented twice; the repo
standard grades that against the diff.

The existing black-box contracts — ~13 status/structure/gate-cache contracts in
`gate-runtime-contracts.sh` and the shift touched-scope/lease contracts in
`gate-runtime-shift-contracts.sh` — drive `bench.sh <cmd>` and assert stdout and
exit codes. They run unchanged against the ported binary and are the regression
net. The parsers and renderer gain `go test` table tests, the new unit layer.

## User stories

1. As a kit developer, I want a `internal/structure` engine — whole-tree and
   touched-scope modes, the `.bench/structure.budgets` per-path override parser
   (trailing-`/` dir budgets, `#` comments, warn-and-skip on malformed lines,
   last line without a trailing newline), `FILE TOO LONG`/`DIR CROWDED` line
   emission, the stderr debt summary, exit 1 on violations, and a violation count
   — so that `bench structure`, `status`, and the shift loop read structural debt
   from one detector. Line: claude-opus-4-8 / medium. The dir-crowding boundary
   logic and the budget-override rule are the intricate, load-bearing core three
   surfaces depend on, so this sets the idiom at the mid tier as slice 2's `maps`
   engine did.

2. As the gate and the Stop hook, I want the gate tree-hash — the throwaway-index
   content hash of tracked-plus-untracked-unignored files, printing the hash or
   `none` — ported to `internal/git` and exposed as `bench tree-hash [root]`, with
   `gate_record` (`bin/bench.sh`) and `record_gate` (`.bench/hooks/stop.sh`) both
   calling it and the stop.sh mirror deleted, so that the verdict-cache key has one
   source. Line: claude-opus-4-8 / medium. The hash sits on the verdict-record
   path and a wrong or forged key silently corrupts every dashboard's freshness
   read, which is why the map flagged its fail-safe as safety-critical.

3. As the Stop hook running with no reachable binary, I want `record_gate` to fail
   safe — skip the cache write loudly and never record a verdict keyed to a guessed
   tree — so that a missing platform binary degrades to "no verdict" rather than a
   forged green the ambient surface reports as real. Line: claude-opus-4-8 /
   medium. This is the map's named hostile-input owner for the write path and the
   gate-blind spot; getting the safe posture wrong is worse than a crash.

4. As a bench user, I want `internal/worktree` to own the pool directory path
   (`$BENCH_HOME/worktrees/<basename>-<cksum>`, the POSIX-`cksum` figure
   reproduced byte-for-byte) and the lease-file path (`git rev-parse --git-path
   bench-lease`), exposed as `bench worktree-pool <root>` and `bench
   worktree-lease-file <path>`, with `bin/bench-worktree.sh`'s two helpers becoming
   adapters over them, so that the pool/lease convention has one source across the
   slice-3→7 window. Line: claude-opus-4-8 / medium. The cksum reproduction is a
   silent-divergence risk — a wrong figure breaks warm-pool exclusion and the
   worktree lease with no error — so it takes the mid tier and a pinned table test.

5. As every agent session, I want `bench status` ported to `internal/status` —
   composing the structure violation count, the unresolved-maps count
   (`maps.UnresolvedCount`), the open-learnings count (`learnings.Rows`), the
   worktree active count (leased-pool and out-of-pool entries only), the
   merged-spec retirement count, and the parked-idea footer; reading the gate cache
   and comparing its tree against `git.TreeHash`; and emitting the lead `▶` line,
   the severity-ordered rows under a five-row budget with `+N more`, and the footer
   — byte-identical to the shell renderer, so that the ambient dashboard reads from
   the Go core. Line: claude-opus-4-8 / medium. This is the composition crux of the
   slice: a wrong compose, sort, or budget silently mis-renders the surface every
   session and the session-start hook consume verbatim.

6. As a kit developer, I want the merged-spec retirement counter — the count of
   `specs/*.md` carrying a line-start `Status: implemented` outside a ``` fence,
   CRLF-safe, positive-marker-only, gated to the default branch by status — ported
   into `internal/status`, so that the retirement signal keeps its exact
   silent-when-absent-or-fenced behavior. Line: claude-sonnet-4-6 / medium. The
   fence/CRLF/no-`Status:` edges are mechanical and fully pinned by the retirement
   contract plus a table test, which is the cheap-tier case.

7. As a bench user, I want `internal/roadmap` to own `ROADMAP.md` — `bench idea
   "<text>"` (dated append, trailing-newline normalization so a hand-edited last
   line never swallows the entry, exit 2 on empty text), `bench roadmap`
   (cat-or-empty), and the `^- ` parked count status's footer reads — so that the
   append normalization and the count it protects share one source. Line:
   claude-sonnet-4-6 / medium. Mechanical file I/O fully pinned by the idea/roadmap
   contract and the footer contract.

8. As a reviewer binding tiers, I want `bench models` ported to `internal/models` —
   the Anthropic Models API query (Go `net/http`, dropping the python3+curl
   dependency) with its indented id list, and the deterministic no-key guidance
   text — so that the model-discovery command is one language. Line:
   claude-sonnet-4-6 / medium. The JSON id-extraction and the no-key text are pure
   and table-tested; the live HTTP call is the untested subprocess-equivalent
   boundary, manually smoke-verified.

9. As a kit maintainer, I want the strangler router to send
   `status|structure|models|idea|roadmap` and the three plumbing subcommands to
   the Go binary through `route_binary`, `bin/bench-status.sh` deleted with its
   `source` line removed, and every dangling reference fixed in the same commit —
   `package.json` ships `bin/` wholesale so nothing changes there, but the
   package-contract's expected-file list, the link contract's installed-helper
   assertion, `bin/bench-link.sh`'s planted-CLI copy list, the two runtime-contract
   fixtures that `cp` the file, and the `bench-status.sh` mentions in `bin/bench.sh`
   and `stop.sh` comments all update — so that no command has two live
   implementations and the stale-reference sweep stays green. Line:
   claude-sonnet-4-6 / medium. Mechanical dispatch and reference edits at the one
   seam slice 1 built.

10. As a kit maintainer, I want the two runtime-contract fixtures that copy
    `bench.sh` and run a **copy** to provision the Go binary the ported code now
    needs (the stop.sh gate-cache-write fixture and its sibling), and a **new**
    fail-safe contract asserting a missing-binary stop.sh writes no forged cache,
    so that the port's gate-cache write stays gate-locked and the safe posture from
    story 3 is enforced rather than trusted. Line: claude-opus-4-8 / medium.
    Touching a gate fixture is the worst defect class in this kit (`craft-gate`);
    provisioning without weakening the assertion, and adding the missing safety
    coverage, takes the mid tier.

11. As the teammate who just walked in, I want `README.md`'s CLI tree and any prose
    naming `bench-status.sh` updated to the Go package layout (new `internal/`
    packages: `status`, `structure`, `worktree`, `models`, `roadmap`), so that the
    repo map matches the repo. Line: claude-sonnet-4-6 / medium. Mechanical doc
    edits verified by the docs gate's stale-reference sweep.

## Implementation decisions

- **Package layout.** Five new packages under `internal/`: `structure` (checker +
  budgets parser + touched-scope), `worktree` (pool/lease path facts),
  `models` (API query + no-key text), `roadmap` (`ROADMAP.md` owner: idea, roadmap,
  parked count), and `status` (the renderer that composes the others plus the
  retirement counter). `internal/git` gains `TreeHash(root)`. `cmd/bench` adds one
  dispatch entry each for `status`, `structure`, `models`, `idea`, `roadmap`,
  `tree-hash`, `worktree-pool`, `worktree-lease-file` — names into the existing
  `commands` map, no new mechanism. Each package stays under the file/dir budgets it
  itself enforces.

- **`internal/status` composes in-process; no self-exec.** Slice 2's shell adapters
  shelled out to `bench maps --count` / `bench learnings` because the renderer was
  bash. In Go, `status` calls `maps.UnresolvedCount`, `learnings.Rows`,
  `structure` (violation count), `worktree` (pool/lease paths for the warm-pool
  exclusion), `roadmap.ParkedCount`, and its own retirement counter directly. The
  close-readiness rule, the TOON count, and the structure detector each still live
  in exactly one package; status counts, it does not re-derive.

- **Tree-hash is one function, two exposures.** `git.TreeHash(root)` runs the
  throwaway-index computation (`GIT_INDEX_FILE` in the OS temp dir — never inside
  the repo, or it joins the tree it hashes; `read-tree HEAD` or `--empty`, `add
  -A`, `write-tree`), returning the hash or `none`. `bench tree-hash` exposes it;
  `status` calls it in-process for the gate-cache freshness compare. Both callers
  read one implementation.

- **The stop.sh mirror dies by reusing stop.sh's own resolver — no duplicated
  binary-path knowledge.** `.bench/hooks/stop.sh` already resolves the `bench.sh`
  wrapper via its `bench_cmd()` helper (searches `.bench/bin/bench.sh`,
  `bin/bench.sh`, then PATH). `record_gate` calls `"$cmd" tree-hash` through that
  same wrapper, and the wrapper owns Go-binary resolution (`bench_binary_path`).
  So the hook needs **no** second copy of the resolution logic — it reuses the
  resolver it already has. This settles the map's one uncertainty flag (item 7)
  without duplication and without escalation. If `tree-hash` returns non-hash
  output (binary missing → `route_binary` exits 127, stdout empty), `record_gate`
  warns and skips the write (story 3). `gate_record` in `bin/bench.sh` calls
  `"$(bench_binary_path)" tree-hash "$root"` directly (it is already inside the
  sourced wrapper) with the same skip-on-non-hash guard.

- **The router flip is additive at one seam.** The five user-facing cases change
  from their shell-function calls to `route_binary "$@"` (matching `version` and
  the slice-2 five); `idea` drops its `shift; idea "$@"` for `route_binary "$@"`
  (the binary reads the text from argv). Three plumbing cases (`tree-hash`,
  `worktree-pool`, `worktree-lease-file`) are added the same way. The `. "$…"/
  bench-status.sh` source line is removed and the file deleted. `repo_root` and
  `default_branch` stay in `bench.sh` (they are the wrapper's, not the ported
  file's).

- **Surviving shell callers become thin adapters.** `structure_touched_since`
  (used only by `shift_loop`) is redefined in `bin/bench.sh` as a one-line call to
  `"$(bench_binary_path)" structure --since "$base"` (captured, not `exec`, so the
  loop continues) — the Go `--since` mode does the `git diff --diff-filter=ACMR
  base..HEAD` internally and checks only touched files. `gate_record`'s tree-hash
  line becomes the binary call above. `bin/bench-worktree.sh`'s `worktree_pool` and
  `worktree_lease_file` become adapters over `bench worktree-pool` / `bench
  worktree-lease-file`. Every adapter is a one-glance pass-through with no logic.

- **`models` uses `net/http`, not curl+python3.** `internal/models` GETs
  `https://api.anthropic.com/v1/models` with the `x-api-key` and
  `anthropic-version` headers, parses `data[].id`, and prints the indented list;
  the no-key branch prints the same guidance text the shell heredoc did. The HTTP
  call is the untested boundary (as git subprocess is elsewhere); the JSON
  id-extraction and the no-key text are pure and table-tested. Dropping the two
  binaries is the map's named incidental win, not a contract change.

- **Deletion touches distribution — fixed in the same commit.** `bin/bench-status.sh`
  is a shipped file: `package.json` `files[]` includes `bin/` wholesale (unchanged),
  but `gate-package-contracts.sh`'s expected-file list names it explicitly, the
  link contract asserts `.bench/bin/bench-status.sh` is installed, and
  `bin/bench-link.sh` copies it into `.bench/bin/`. All three drop the file. This
  is the delete discipline: promote nothing (no durable decision lives only here —
  the tree-hash and worktree comments move with their code), delete, and fix every
  dangling reference in the same change so the sweep stays green.

- **Deviations from the map: none.** Scope is the map's #1 answer (whole file,
  deleted). Tree-hash ownership is #2, worktree ownership is #3, and the package
  homes are the spec-writer's call the Handoff (item 1) delegated. The one
  uncertainty flag (item 7, stop.sh resolution) is settled above without
  duplication.

## Testing decisions

- **What a good test is here.** Acceptance drives the public `bench` entry and
  asserts stdout/stderr text and exit codes — never Go internals. The existing
  runtime and shift contracts already do exactly this and carry over. Go table
  tests are additional, at the pure-function seam (the structure engine, the
  tree-hash, the cksum pool key, the retirement scan, the models parser), and are
  where the shell-untestability tax retires.
- **Seams.** Four: the **CLI acceptance seam** (the runtime/shift contracts driving
  `bench.sh <cmd>` → Go binary — prior art: they exist and pass today against
  shell); the **`go test` unit seam** (table tests beside each `internal/` package
  — prior art: slice 2's `internal/maps` and `internal/toon` tests); the
  **gate-cache coexistence seam** (the stop.sh and verdict-record contracts driving
  the cache write through `bench tree-hash` — prior art: the gate-cache-write
  contract at `gate-runtime-contracts.sh:169`); and the **worktree adapter seam**
  (the shift lease contracts and the warm-pool status contract exercising the
  pool/lease facts — prior art: the red-rollback lease-release assertion at
  `gate-runtime-shift-contracts.sh:113` and the warm-pool contract at
  `gate-runtime-contracts.sh:150`).
- **Gate command:** `.bench/gate.sh` (the project gate), whose Go layer already
  runs `gofmt -l`, `go vet ./...`, `go test ./...`, `go build`, and the four
  cross-compile targets — so new `internal/` packages and their tests are graded
  with no new gate wiring.

### Seam diagram — CLI acceptance seam (status / structure / models / idea / roadmap)

    trigger: runtime gate contract (or a user / session-start hook) runs `bench <cmd> [args]`
        │
        ▼
    argv ─────────▶ [ bin/bench.sh: route_binary ] ──▶ exec dist/bench
    repo state ───▶ [ cmd/bench dispatch → internal/<pkg> ]
    (git, files,    [ internal/status composes structure/maps/ ] ──▶ stdout: `▶ …` + rows + footer, exit 0/1
     ROADMAP.md,    [   learnings/worktree/roadmap/retirement   ] ──▶ structure: `FILE TOO LONG`/`DIR CROWDED`, exit 1
     gate cache)    [                                          ] ──▶ idea: exit 2 on empty; usage on stdout
        ◀ tests attach here: the existing runtime contracts run `bash bench.sh <cmd>`
          in a fixture repo and assert stdout shape + exit code, unchanged from the
          shell surface. Red before the Go subcommand exists: route_binary reaches a
          binary that answers `unknown subcommand` (exit 2), so the expected line is absent.

### Seam diagram — `go test` unit seam (engines + pure helpers)

    trigger: gate Go layer runs `go test ./...`
        │
        ▼
    files / budgets / ──▶ [ internal/structure.Check / Budgets   ] ──▶ violation lines + count
    root / cells /        [ internal/git.TreeHash                 ] ──▶ hash or `none`
    JSON body             [ internal/worktree.Pool / LeaseFile    ] ──▶ pool dir (cksum) / lease path
    (temp-dir or          [ internal/status.retirementCount       ] ──▶ count (fence/CRLF/default-branch)
     in-memory fixtures)  [ internal/models.parse                 ] ──▶ id lines / no-key text
        ◀ tests attach here: table tests assert the dir-crowding boundary and budget
          override; TreeHash's temp index lives outside the repo; Pool reproduces the
          POSIX cksum figure exactly; the retirement scan is fence-aware and
          positive-marker-only. Red before the function exists: the package/function
          is absent and the test does not compile.

### Seam diagram — gate-cache coexistence seam (stop.sh + gate_record)

    trigger: Stop hook (armed shift) or `bench gate` finishes a gate run
        │
        ▼
    verdict + root ──▶ [ .bench/hooks/stop.sh: record_gate ] ──("$cmd" tree-hash)──▶ [ dist/bench: git.TreeHash ]
                       [ bin/bench.sh: gate_record         ] ──($(bench_binary_path) tree-hash)──▶ [           ]
                            │  hash valid → write `<status> <hash> <iso8601>` to <git-dir>/bench-last-gate
                            │  hash absent (binary missing) → warn, write NOTHING (never forge)
        ◀ tests attach here: the gate-cache-write contract provisions dist/bench in the
          fixture and asserts the cache format; a NEW contract omits the binary and asserts
          no forged cache is written. Red before the port: stop.sh's inline mirror still
          computes the hash (no binary call), so the new fail-safe assertion cannot pass and
          the mirror-deletion sweep still finds the duplicated computation.

### Seam diagram — worktree adapter seam (bench-worktree.sh over the binary)

    trigger: shift loop lease claim/release, or `bench status` worktree row
        │
        ▼
    root / wt path ──▶ [ bin/bench-worktree.sh: worktree_pool / worktree_lease_file ]
                       [   → "$(bench_binary_path)" worktree-pool / worktree-lease-file ]
                       [ internal/status (in-process): worktree.Pool / LeaseFile      ] ──▶ warm-pool exclusion
        ◀ tests attach here: the shift red-rollback contract asserts the lease is released
          (bench-worktree.sh adapter path); the warm-pool status contract asserts a leased
          pool entry surfaces and a warm one does not (status in-process path). Red before
          the port: status's shell pool/lease helpers are deleted and unresolved, so the
          worktree row errors or the warm-pool exclusion mis-fires.

### Acceptance coverage map

Per-item granularity is stated where a behavior quantifies over a set (each budget
kind, each gate-cache branch, each retirement edge).

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | `bench structure` flags an over-budget shell file (`FILE TOO LONG`), exit 1 | CLI acceptance | "bench structure shell-file" contract (`:261`) before the Go subcommand → `unknown subcommand`, `FILE TOO LONG` absent | the length-budget branch is the base structure assertion |
| 1 | per-path budgets: a granted file/dir budget applied, a tightening override applied, a malformed line warned-and-skipped, last line without trailing newline parsed — each kind | CLI acceptance | "bench structure budgets" contract (`:274`) before the subcommand → warn line and override figures absent | the contract asserts all four budget-kinds in one fixture; a wrong parse breaks a named substring |
| 1 | `DIR CROWDED` on a directory whose path contains spaces, path preserved un-split | CLI acceptance | "bench structure path-with-spaces" contract (`:293`) before the subcommand → `space dir/` row absent | catches a Go port that splits on whitespace instead of argv/exact-path handling |
| 1 | dir-crowding boundary + budget-override edges keep coverage at the pure seam | go test unit | `internal/structure` table test before `Check`/`Budgets` exist → does not compile | pins the boundary count and override rule below the CLI, where the shell contracts can't reach every case |
| 1 | `structure --since <base>` checks only files touched since base (touched-scope), driving the shift refactor gate | CLI acceptance (shift) | "bench shift touched-scope structure" contract (`shift:118`,`:132`) before the subcommand → pre-existing debt wrongly triggers refactor, or touched debt missed | the shift contracts assert exactly the touched-vs-repo-wide distinction |
| 2 | `bench tree-hash` prints the throwaway-index content hash (or `none`), temp index outside the repo | go test unit + CLI acceptance | `internal/git.TreeHash` table test before it exists → fails; verdict-record contract (`:221`) asserting `<tree>` equals `HEAD^{tree}` after a same-tree commit | the same-tree-survival and dirty-tree cases pin the hash; the temp-index-location test guards the tree from including its own index |
| 2 | `gate_record` and `record_gate` write `<status> <tree> <iso8601>` keyed to the binary's hash; the stop.sh inline mirror is gone | gate-cache coexistence | gate-cache-write contract (`:169`) with dist/bench provisioned (story 10) → format assert; grep for the deleted mirror in stop.sh | one hash source feeds both writers; the mirror-deletion is verified by absence, not trust |
| 3 | a missing binary makes `record_gate` skip the write (no forged verdict) and warn | gate-cache coexistence | **new** fail-safe contract: stop.sh fixture with no dist/bench, green stub gate → assert `bench-last-gate` not written | closes the map's gate-blind spot; catches a port that writes `<status> none <date>` and mis-reports it as real |
| 4 | `bench worktree-pool <root>` = `$BENCH_HOME/worktrees/<basename>-<cksum>` with the cksum reproduced exactly; `bench worktree-lease-file <path>` = the git-path lease | go test unit | `internal/worktree` table test before `Pool`/`LeaseFile` exist → does not compile; a known-`root`→known-cksum vector | a diverging cksum silently breaks warm-pool exclusion and lease naming; the vector pins it byte-for-byte |
| 4 | `bench status` excludes warm pooled worktrees (released, no lease) and surfaces leased-pool + out-of-pool ones, via the in-process facts | CLI acceptance (worktree) | "bench status warm-pool" contract (`:150`) after the shell helpers are deleted, before the Go compose → warm entry wrongly surfaces or leased one omitted | the contract hardcodes the pool formula (`:158`); a wrong Go pool path fails the exclusion |
| 5 | `bench status` clean/footer/stale-gate/fresh-green renders byte-identical (lead `▶`, footer separate, never-leads-on-footer) | CLI acceptance | clean/footer/stale/fresh contracts (`:61`,`:68`,`:78`,`:88`) before the Go subcommand → expected lines absent | the four framing contracts pin the lead line, footer placement, and stale-vs-clean read |
| 5 | severity sort, five-row budget + `+N more`, lowest-severity dropped first | CLI acceptance | "bench status budget" contract (`:132`) before the subcommand → wrong lead / >5 rows / lowest signal retained | the budget contract asserts the ordering and the row cap together |
| 5 | unresolved-maps count and open-learnings floor sourced in-process, byte-identical | CLI acceptance | "unresolved-maps count" (`:108`) and "learnings-floor" (`:246`) contracts before the compose → counts absent/wrong | one engine feeds status and `bench maps`; the contracts pin the shared figure |
| 6 | retirement row counts unfenced line-start `Status: implemented` specs, default-branch only, silent when fenced/absent/off-branch, self-clears on deletion — each edge | CLI acceptance | "retirement-signal" contract (`:184`) before the Go counter → row absent/miscounted across its five sub-cases | the contract walks default-branch, fenced, plain, off-branch, and deletion in one fixture |
| 6 | fence/CRLF/no-`Status:` scan edges keep coverage at the pure seam | go test unit | `internal/status` retirement table test before it exists → does not compile | pins the CRLF and fence edges below the CLI |
| 7 | `bench idea` appends one dated `- YYYY-MM-DD  <text>` line, creates the file, normalizes a newline-less last line, exit 2 on empty, joins unquoted multi-word args | CLI acceptance | "bench idea/roadmap" contract (`:11`) before the Go subcommand → dated line / exit-2 / count-after-handedit asserts fail | the contract pins the append format, the empty-guard, and the normalization the footer count depends on |
| 7 | `bench roadmap` prints the parked ideas or the empty state | CLI acceptance | same contract's `roadmap` assert (`:19`) before the subcommand → parked line absent | pins cat-or-empty |
| 8 | `bench models` no-key branch prints the tier-binding guidance text | go test unit | `internal/models` table test before it exists → does not compile | deterministic text, table-pinned; the live API branch is not gate-observable (network) → manual smoke, see Won't handle |
| 8 | `bench models` extracts `data[].id` from an API JSON body into the indented list | go test unit | `internal/models` parse table test before the parser → fails | pins the JSON extraction without a network dependency |
| 9 | `status\|structure\|models\|idea\|roadmap` + the three plumbing names route to the binary; `bench-status.sh` gone and unsourced; no dangling reference | CLI acceptance (all above) + parse layer + docs gate | any runtime contract run while a case still calls the deleted shell fn → function-not-found; `bash -n` / gate load error on a `source` of a deleted file; package/link/docs sweeps naming the file → red | the whole runtime suite plus the parse and stale-reference layers fail if a command still routes to shell or a reference dangles |
| 10 | the stop.sh gate-cache-write fixture (and its `:419` sibling) provision dist/bench and stay green; the new fail-safe contract goes red-then-green | gate-cache coexistence | the provisioned write contract fails until the fixture copies dist/bench; the fail-safe contract fails until story 3 lands | proves the port's write path is gate-locked and the safe posture is enforced, not assumed |
| 11 | `README.md`'s CLI tree and prose match the Go layout (no `bench-status.sh`, new `internal/` packages named) | docs gate | the docs gate's stale-reference sweep against a tree still naming `bench-status.sh` → red | the conformance layer fails when docs reference a deleted file |

### Edge inventory

Walked per behavior against the profile's shell-CLI hostile-input checklist and the
map's item-6 owners; each resolved as a coverage row above or a **Won't handle**
line here.

- **paths/dirs with spaces** — covered: the structure path-with-spaces contract
  (`:293`) and the worktree pool key over a spaced basename; Go uses `os/exec` argv
  and exact-path matching, never a shell split.
- **hand-edited file without a trailing newline** — covered: the budgets contract's
  unterminated last line (`:281`, tightening override on it) and the idea
  normalization assert (`:26`, `^- ` count stays correct); table-tested in
  `internal/structure` and `internal/roadmap`.
- **CRLF / carriage returns** — covered: the retirement scan strips `\r` (table
  test) and the structure/budgets parsers handle it; the retirement contract's
  fenced/plain fixtures exercise line handling.
- **absent vs present-but-empty file** — covered: `roadmap` empty-state, `status`
  with no gate cache / no ROADMAP.md / no `specs/` and `decisions/` (clean
  contract), and the structure "no tracked source files" path.
- **fenced/anchored markers** — covered: the retirement scan counts only unfenced
  line-start `Status: implemented` (retirement contract's `fenced.md`,
  `plain.md`); table-tested.
- **required tool / binary missing** — covered: story 3's fail-safe contract
  (stop.sh, no dist/bench) and slice 1's `route_binary` missing-binary path
  (exit 127) reused unchanged by the five user-facing commands.
- **cwd deeper than root / invocation through a symlink** — covered: every command
  resolves the root via `git rev-parse` in the binary (as the shell did); symlink
  resolution is slice 1's router property, unchanged.
- **cksum divergence across platforms** — covered: the worktree pool table test
  pins a known `root`→cksum vector, and the warm-pool contract hardcodes the
  formula; a diverging Go implementation fails both.
- **SIGINT mid-shift leaving a stale lease** — covered by the existing shift
  red-rollback lease-release contract (`shift:113`) exercising the ported
  `bench-worktree.sh` adapter; the interrupt/cleanup logic itself stays shell
  (slice 7), so this slice only proves the lease path still resolves.
- **Won't handle: a repository path containing a literal newline** — parity with
  the shell (which already misreads it); no in-scope caller can produce one, and
  the worktree `git worktree list --porcelain` parse in status preserves the shell's
  behavior rather than fixing it. Fixing it is not this slice.
- **Won't handle: the live `bench models` API branch under the gate** — it hits the
  network, so it is not gate-observable; the JSON parser and the no-key text are
  table-tested and the live call is manually smoke-verified once, exactly as the
  shell version was never gate-tested against the live API.
- **Won't handle: worktree adapter subprocess-spawn cost in `bench-worktree.sh`'s
  lease loops** — each `worktree_pool`/`worktree_lease_file` call is now a binary
  spawn; the binary runs in milliseconds and the lease loop already spawns git per
  iteration, so this is not optimized here. Slice 7 (worktree/shift port) collapses
  these into in-process Go calls.

## Out of scope

- **Slices 4–8** — `git-guard.py` absorption, hook logic behind shims,
  `doctor`/`link`, the worktree/shift loop, and the gate-fragment port. Each is a
  distinct capability with its own spec by the map's dependency order; the
  shift/worktree contract backfill (parent Handoff item 7) precedes slice 7.
  Estimate to build later: per the parent map, ~10–15 spec-sized sessions total.

- **In-process collapse of the worktree adapters** — this slice inverts ownership
  (Go owns the facts, `bench-worktree.sh` adapts) but leaves `bench-worktree.sh`'s
  lease claim/release logic in shell. Folding those call sites into in-process Go is
  slice 7's worktree port, a separate capability. Estimate: subsumed by slice 7,
  ~15 edits, ~6 gate runs on top of its own scope.

- **In-process collapse of `gate_record` / `structure_touched_since`** — both stay
  shell adapters over the binary this slice; they fold into Go when the shift loop
  and gate port (slices 7–8) move. Separate capability. Estimate: subsumed by
  slices 7–8.
