# Structure splits and the accept mechanism

Status: implemented

## Problem

`bench structure` flags four over-budget source files. Two are genuine
structural debt worth splitting; two (`internal/status/status.go`,
`internal/contract/runtime/runtime_status_test.go`) are barely over budget and
cohesive — splitting them would hurt. But leaving them flagged means the
structural-debt count reads a permanent "2 issues," and a signal that is always
non-zero trains every reader to ignore the row. There is no way to say "this
file is over budget on purpose, and here's why" without either fragmenting a
cohesive file or silently raising a global threshold that hides real debt.

## Solution

Two independent moves that leave the structure signal both under budget and
credible:

1. **Split the two genuine files along their responsibility seams** — pure
   relocation of test functions into command-family / concern files, each
   landing under the 400-line budget, with identical coverage. No test logic
   changes.
2. **Add a reviewer-visible accept mechanism to `bench structure`** — a
   `.bench/structure-accept` sidecar listing `<path> <one-clause reason>` rows.
   An accepted file is excluded from the violation count and printed in a
   separate `accepted:` section, so the suppression is per-file, reasoned on the
   page, and reversible — never a silent threshold bump. Seed it with the two
   cohesive noise files so the count reads clean.

After both, `bench structure` reports zero violations, `bench status` shows no
structure signal, and every remaining over-budget file is either genuinely fine
(split) or accepted with a stated reason on the page.

## User stories

1. As a reviewer, I want the oversized AXI wave-2 contract file split into one
   file per covered command family, each under the length budget, with every
   sub-test still registered and passing, so the debt drops without any coverage
   loss.
   Line: claude-sonnet-5 / low. A pure relocation of whole test functions at a
   known package seam, fully re-run by the existing suite, needs only the cheap
   tier.

2. As a reviewer, I want the oversized line-routing conformance file split along
   the static-parse vs subprocess-exec seam, each half under budget, with every
   check still wired into root conformance, so the debt drops without dropping a
   check.
   Line: claude-sonnet-5 / low. Same mechanical relocation at a known seam,
   proven by the unchanged conformance run.

3. As a reviewer, I want `bench structure` to read `.bench/structure-accept` and
   exclude each accepted file from the violation count, so an accepted
   over-budget file stops inflating the structural-debt number and the count
   `bench status` reads.
   Line: claude-opus-4-8 / medium. Oracle-adjacent logic feeding the count a
   status board trusts, where correctness of the signal outweighs speed.

4. As a reviewer, I want every suppressed violation printed as
   `accepted: <path> — <reason>` with a count, in a section separate from the
   live violations, so each acceptance is visible and its reason is on the page.
   Line: claude-opus-4-8 / medium. The visibility guarantee is the whole point
   of the suppression surface; it is graded on output shape a status hook parses.

5. As a reviewer, I want a malformed accept row (a path with no reason) reported
   and not honored, so acceptance always carries a stated reason and can never be
   granted silently.
   Line: claude-opus-4-8 / medium. A fail-closed hostile-input rule the count's
   integrity depends on.

6. As a reviewer, I want a stale accept row (a path that is no longer a scanned
   source file) reported as a warning, so the accept list cannot quietly
   accumulate dead entries.
   Line: claude-opus-4-8 / medium. Keeping the list honest over time is a
   correctness property of the suppression surface.

7. As a reviewer, I want a present-but-unreadable accept file to be loud — a
   non-zero exit, never a silently empty accept list — so a vanished or
   unreadable accept file can never change counts behind my back (the false-empty
   rule, FT29).
   Line: claude-opus-4-8 / medium. This is the exact false-empty defect class the
   platform forbids; getting the read-error posture right is the reason this
   story is not cheap.

8. As a reviewer, I want `.bench/structure-accept` seeded with the two cohesive
   noise files, each with a one-clause reason, so `bench structure` and
   `bench status` read clean the moment the mechanism lands.
   Line: claude-opus-4-8 / medium. The seed's paths and reasons are the reasoned
   content that proves the mechanism end to end and closes the count to zero.

## Implementation decisions

**The two splits are pure moves.** Function bodies relocate verbatim between
files in the same package; Go's cross-file same-package visibility means the
existing single-parent registrar can call relocated helpers and sub-tests
unchanged, and each new file follows the package's established
one-`Test…`-per-family convention. No assertion, label, `NoteContractFailure`
string, or `RunParallel`/`t.Run` registration is edited. The non-descriptive
`wave2` batch name is retired in the same move — every function lands in a
family-named file.

**AXI wave-2 split** (`internal/contract/axi/axi_wave2_test.go`, 639 lines →
deleted) by command family:

- `axi_diff_test.go` — the `bench diff` branch-relative review family:
  recorded-base, fallback/shape, error-posture, `--full`,
  git-failure-propagation, and control-byte-posture contracts, under a new
  `TestAXIDiffContracts` registrar.
- `axi_diff_commit_test.go` — the `bench diff --commit` landed-commit review
  family: commit-range, commit-merge, and commit-error-posture contracts, under
  `TestAXIDiffCommitContracts`. (The `bench diff` family is split in two because
  the diff contracts alone exceed 400 lines; `--commit` is a distinct
  command-mode with its own story cluster, so it is the natural inner seam.)
- `axi_coverage_test.go` — the `bench coverage` family: extraction, state/error,
  and `--check` validation contracts, under `TestAXICoverageContracts`.
- `testAXIMapsGuardsHelp` (the `bench maps`/`bench guards` help contract) folds
  into the existing `axi_guards_test.go`, registered under its
  `TestAXIGuardsContracts` — the maps/guards family already has a home there, and
  a standalone one-test file would be the fragmentation the structure check warns
  against.
- The shared assertion helpers (`requireOutputLine`, `requireOutputPrefix`,
  `requireNoOutput`, `requireLogRow`, `shellQuote`) move into the package's
  existing shared-helper file `axi_helpers_test.go` (already home to
  `requireContainsFold`), since `requireLogRow` is used by both diff files.

**Line-routing split** (`internal/conformance/line_routing_checks_test.go`, 423
lines → deleted) along static-parse vs subprocess-exec:

- `line_routing_static_test.go` — the `checkLineRouting` dispatcher plus the
  static-parse checks that only read and parse files: `checkLineBinding` (and
  `TestLineBindingAcceptsOpaqueSafeModelTokens`,
  `TestLineBindingRejectsUnsafeModelTokens`) and `checkClaudeHookWiring` (and
  `TestClaudeHookWiringBites`).
- `line_routing_exec_test.go` — the checks that run shipped scripts as
  subprocesses: `checkAgentHookBehavior` and `checkAdapterLineGuards`. Both rely
  only on package-shared harness helpers (`runWithInput`, `runAtEnv`,
  `tempGitRepoWithLines`, …) already defined elsewhere in the package, so no
  helper is duplicated.

**The accept mechanism lives in `internal/structure`, in `Check`.** A new
`loadAccepts` reads `.bench/structure-accept` into a `path → reason` map plus
warnings, mirroring `loadBudgets`' comment/whitespace discipline (strip from the
first `#`; blank remainder skipped; unterminated last line parsed;
first-wins on a duplicate path) with three deliberate differences:

- The **reason is the remainder** of the stripped line after the first
  whitespace-delimited path token, not a `Fields` split.
- A row with a path but **empty reason is malformed** — reported as a warning and
  not added to the accept set (a reason is the price of acceptance).
- The **read-error posture is fail-closed**: `os.IsNotExist` returns an empty
  accept list with no error (missing file = empty, exit 0), but any *other* read
  error (present but unreadable) is surfaced loudly — `Check` appends a named
  error line and forces a non-zero result through the same violation count both
  the report and `ViolationCount` return, so no surface can silently observe an
  empty accept list. This is the one intentional exception to "non-zero only on
  real violations": the normal accept states (absent, malformed, stale) never
  change the exit code; only a genuine I/O failure does.

**Matching and reporting** (full-tree `all` mode):

- An accept row matches a violation by its subject path — a file path for
  `FILE TOO LONG`, or a `dir/` key for `DIR CROWDED` — reusing `budgetFor`'s
  keying, so one match covers both violation kinds though only file rows are
  seeded.
- A matched violation is **excluded from the count** and instead recorded for the
  `accepted:` section, which prints one `accepted: <path> — <reason>` line per
  suppressed violation plus a count of suppressed files. Live (non-accepted)
  violations print unchanged.
- A valid accept row whose path is **not among the scanned source files** is
  reported stale and otherwise ignored.
- Exit-code mapping is unchanged: non-zero iff any live violation remains (plus
  the fail-closed read-error case above).

**Touched (`--since`) mode** honors accept-set exclusion (so the shift refactor
gate does not re-flag an accepted file) but emits neither the `accepted:` section
nor stale warnings, because a partial file set cannot judge staleness or list
completeness.

**Single-engine property preserved.** Because `bench status` reads
`ViolationCount(root)` off the same `Check(root, "all", nil)`, the accepted
exclusion flows to the status count with no change in `internal/status` — the
report and the count cannot disagree.

**The seed file `.bench/structure-accept`** carries a self-documenting header
comment (mirroring `structure.budgets`) and two rows: `internal/status/status.go`
and `internal/contract/runtime/runtime_status_test.go`, each with a one-clause
reason naming the file as cohesive and barely over budget. It is consumer-owned
and stays out of `package.json` `files[]`, exactly like `structure.budgets`.

**Sequencing stays green throughout.** `bench structure` violations are a
`bench status` signal, not a gate-failing layer, so each step commits green: the
splits drop the count 4 → 2, the mechanism plus tests land (count still 2, a
non-failure), and the seed closes it 2 → 0. The compiled core is rebuilt before
the runtime/status and AXI contract fragments run, since those invoke the built
`bench` binary.

## Testing decisions

- **A good test here** exercises `bench structure`'s external behavior at the
  `Check` seam — feed a real temp git tree plus a `.bench/structure-accept`, read
  the report string and the returned violation count — never an internal parse
  detail. This matches the package's prior art: `TestFileTooLongAndBudget`,
  `TestMalformedBudgetWarns`, and `TestViolationCount` all init a temp repo, write
  files, `git add`, and assert `Check`'s report substrings and count.
- **Seams tested.** One seam gets new tests: `structure.Check` (accept
  mechanism), in `internal/structure/structure_test.go` alongside the existing
  budget tests. The two splits get **no new tests** — the unchanged suite is the
  oracle, with one review caveat below.
- **The move-integrity caveat.** A green suite alone cannot prove a split
  preserved coverage: a silently *dropped* sub-test also leaves the suite green
  and the file short. So the splits' "same coverage" claim is verified by diff
  inspection — each function added exactly where it was removed, no body edited,
  and the total `RunParallel`/`t.Run` registration count conserved across the new
  files. This is the one check the cheap line still needs beyond the gate.
- **Gate command:** `.bench/gate.sh` (the project gate). It must be green before
  and after. Useful focused runs while building:
  `go test ./internal/contract/axi/... ./internal/conformance/... ./internal/structure/...`
  and `bench structure` for the debt count.

### Seam diagram

Seam 1 — the accept mechanism (`structure.Check`):

    trigger: `bench structure` (no args → all) · `bench status` via ViolationCount · `--since` (touched)
        │
        ▼
    .bench/structure-accept   ──▶  [ loadAccepts → Check: scan tree,          ]  ──▶  report string
    tracked source tree       ──▶  [ suppress accepted violations from count, ]  ──▶  violation count
                                   [ warn on malformed / stale rows,          ]
                                   [ print `accepted:` section + count,        ]
                                   [ fail loud on an unreadable accept file    ]
                                       ◀ tests attach here: internal/structure/structure_test.go drives
                                         Check(root,"all",nil) on a temp repo holding an over-budget file
                                         + a .bench/structure-accept, asserting the count and report lines

Seam 2 — the two splits (the unchanged gate is the oracle):

    trigger: `.bench/gate.sh`  (go test ./internal/... · TestRootConformance · bench structure)
        │
        ▼
    relocated test funcs (verbatim)  ──▶  [ same package compiles;     ]  ──▶  identical pass/fail set
    oversized source files           ──▶  [ each file re-measured       ]  ──▶  each split file ≤ 400 lines
                                              ◀ tests attach here: no NEW test — go test proves no retained
                                                sub-test broke, `bench structure` goes green on the split
                                                paths, and a diff review confirms a pure relocation

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | `axi_wave2_test.go` becomes command-family files each ≤400 lines, every sub-test still passing | unchanged gate: `go test ./internal/contract/axi/...` + `bench structure` | `bench structure` today prints `FILE TOO LONG` for `internal/contract/axi/axi_wave2_test.go` (639 lines); the axi suite passes | a broken move fails `go test`; a still-oversized file stays flagged by `bench structure` |
| 1 | the move preserves coverage (no sub-test silently dropped) | diff review | already covered / not TDD-able — a dropped test keeps the suite green; verified by conserved `RunParallel` count and body-identical relocation | a green suite can't see a deleted test; the registration-count/diff check is what catches it |
| 2 | `line_routing_checks_test.go` becomes static-parse + subprocess-exec files each ≤400 lines, every check still wired | unchanged gate: `go test ./internal/conformance/...` + `bench structure` | `bench structure` today prints `FILE TOO LONG` for `internal/conformance/line_routing_checks_test.go` (423 lines); `TestRootConformance` passes | a lost check makes root conformance miss its diagnostics; an oversized half stays flagged |
| 3 | a valid accept row excludes its file from the violation count | `structure.Check` | new test: over-budget file + valid accept row → want `v==0`; a no-op `loadAccepts` leaves `v==1` | the cheapest wrong impl (ignore the accept file) leaves the file counted → `v` stays 1 |
| 3 | absent accept file → empty list, exit 0 (distinct from unreadable) | `structure.Check` | new test: no `.bench/structure-accept` → want unchanged count, no warning, exit 0 | an impl that errors on absence would fail the exit-0 assertion |
| 4 | each suppressed violation prints `accepted: <path> — <reason>` plus a count | `structure.Check` report | new test asserts the `accepted:` line and count; before impl the report has no such line | an exclude-but-don't-print impl passes the count yet fails the line assertion → invisible suppression |
| 5 | a malformed accept row (no reason) is reported and not honored | `structure.Check` | new test: `<path>` with no reason on an over-budget file → want a malformed warning AND `v==1` | an impl that honors reasonless rows drops the violation (`v==0`) → acceptance without a reason |
| 6 | a stale accept row (path not a scanned source file) is reported | `structure.Check` (all mode) | new test: accept row for an absent path → want a stale warning | an impl that skips missing paths silently emits no warning → dead rows accumulate |
| 6 | a path containing whitespace is owned deterministically | `structure.Check` | new test: accept row whose path has a space → first token is not a scanned file → stale warning, no crash | a naive split would misattribute or panic; the stale path proves deterministic ownership |
| 7 | a present-but-unreadable accept file is loud (non-zero exit) | `structure.Check` / `Command` | new test: 0-perm `.bench/structure-accept` → want non-zero exit + named loud line; a `loadBudgets`-style swallow returns empty at exit 0 | the FT29 false-empty defect: swallowing the read error yields exit 0 with an empty list → counts change silently |
| 8 | the two seeded noise files drop out of the count and `bench status` | `structure.ViolationCount` + `bench status` | after seeding, `bench structure` count 2→0; before seeding it is 2 | wrong seed paths would be stale, suppress nothing → count stays 2 and status still shows the signal |

Cheapest-wrong-implementation check: an always-empty `loadAccepts` (the
no-op stub) is red on stories 3 and 8; an exclude-but-don't-print impl is red on
story 4; an honor-reasonless-rows impl is red on story 5; a swallow-all-errors
loader is red on story 7. No degenerate implementation passes the map.

### Edge inventory

Walked per behavior; each edge is a coverage row above or a **Won't handle** line
here.

- **Empty (present, zero-byte) accept file** — resolved as a coverage-adjacent
  assertion under story 3's absent-file row: parsed as an empty list, exit 0.
- **Comment-only / blank-line accept file** — empty list; same handling as
  `loadBudgets`, asserted with the empty-file case.
- **Inline `#` comment on a row** — stripped from the first `#`, mirroring
  `loadBudgets`; asserted in the valid-row test.
- **Duplicate accept rows for one path** — first-wins (set membership), mirroring
  the budgets loader; the path is excluded once and prints one `accepted:` line.
- **Reason containing a literal `#`** — **Won't handle**: `#` always starts a
  comment, so a reason is truncated at it; reasons are plain one-clause English,
  one parse discipline shared with `structure.budgets`.
- **Accept path containing whitespace** — handled: the first token is treated as
  the path, and since a spaced token is not a scanned source file it is reported
  stale (covered by story 6's edge row), satisfying the Handoff's "structure owns
  paths with spaces."
- **CRLF line endings in a hand-edited accept file** — **Won't handle**: a
  trailing `\r` rides along in the reason string, cosmetic only, mirroring
  `structure.budgets` which also does not strip `\r`.
- **`--since` (touched) mode + accept list** — exclusion is honored so the shift
  refactor gate does not re-flag an accepted file; the `accepted:` section and
  stale warnings are all-mode only. Covered by the implementation decision; not
  separately tested because touched-mode staleness is intentionally undecidable.
- **DIR CROWDED from adding files** — **Won't handle as a risk**: headroom
  verified before the fact — `internal/conformance/` goes 11 → 12 tracked source
  files (at the 12 cap, not over) and `internal/contract/axi/` goes 7 → 9; the
  gate's own structure check confirms no new `DIR CROWDED` line.
- **Interrupted / partial state, re-run idempotency** — not applicable: `Check` is
  a pure read with no scratch state, worktree, or lease; a re-run is identical.

## Out of scope

- **Warn on an inert accept row** (a valid row whose path is present but the file
  is now *under* budget, so it suppresses nothing) — a distinct
  accept-list-honesty check beyond the decided stale-on-missing-path rule; it
  needs its own signal and message. Estimated ~2 edits, 2 gate runs.
