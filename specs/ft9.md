# bench diff --full — compiled base-relative review context

Status: implemented

## Problem

A review agent needs the whole base-relative picture — the changed-file list, the
diff body, and the commit log — but `bench diff` emits only the file list. So
`/bench-review-implementation` step 1 makes two extra hand-run git calls
(`git diff <base>...HEAD` and `git log <base>..HEAD --oneline`) using the base
`bench diff` just printed. That is the one genuinely repeated base-relative call
site: base resolution already lives in `bench diff`, yet the content fetch is
duplicated in prose the agent re-types every review.

## Solution

Add a `--full` flag to `bench diff` that appends — after the existing
`branch/base/method` preamble and `files[N]{status,path}` table — a
`log[N]{sha,subject}` TOON table (from `git log <base>..HEAD`) and a single
delimited raw diff-body block (from `git diff <base>...HEAD`). Bare `bench diff`
stays byte-for-byte unchanged. `/bench-review-implementation` step 1 then calls
`bench diff --full` in place of the two git commands, collapsing the repeated
call site.

## User stories

1. As a review agent, I want `bench diff --full` to append the base-relative diff
   body and commit log beneath the changed-file list, so that I get the whole
   review context from one command instead of two hand-run git calls.
   Line: claude-sonnet-5 / medium. The flag is a thin addition at the existing
   `diff.Command` seam and is fully observable in stdout by the AXI contract
   tests, so the cheap tier at medium effort covers it.

2. As a review agent, I want bare `bench diff` output to stay byte-for-byte
   unchanged, so that existing callers and the non-full path are unaffected.
   Line: claude-sonnet-5 / low. This is a regression guard on an output the
   contract suite already pins, so the cheap tier at low effort is right.

3. As a review agent, I want the log rendered as a TOON `log[N]{sha,subject}`
   table and the diff body as one documented raw block appended last, so that the
   structured data stays parseable while the unified diff stays readable and
   unmangled. Line: claude-sonnet-5 / medium. The render split is decided in the
   map and asserted line-by-line in the contract test, so the cheap tier at
   medium effort suffices.

4. As a review agent, I want empty-since-base to yield empty log and diff sections
   at exit 0, and base-unresolvable or unknown-argument to keep their existing
   exit 1 and exit 2 postures, so that edge states stay predictable.
   Line: claude-sonnet-5 / medium. This is edge behavior at the same
   gate-observable seam, so the cheap tier at medium effort covers it.

5. As a reviewer running `/bench-review-implementation`, I want step 1 to call
   `bench diff --full` instead of the separate `git diff` and `git log` prose, so
   that the one repeated base-relative call site collapses into the command.
   Line: claude-opus-4-8 / low. This edits command-guidance prose the gate grades
   only structurally through a conformance anchor, so it routes off the cheap tier
   to mid; it deviates downward from the profile's "doc authoring → top/high"
   cached row because the edit is a mechanical call-site swap rather than net-new
   guidance prose, so the leverage override does not apply.

## Implementation decisions

- **Where.** `internal/diff/diff.go` `Command` grows a `--full` branch; no new
  module. `diff.Command` stays the deep unit that hides base resolution and
  NUL-safe name-status parsing — the `--full` branch is a thin addition (two
  `git.Raw` calls plus rendering) attaching at the same stdout seam.

- **Arg parsing.** Today `Command` accepts only zero args, `-h`/`--help`, else a
  usage-error exit 2. Recognize `--full` as the one valid flag (order-independent
  with no other positional args); any other/unknown arg still exits 2. Bare
  `bench diff` behavior is untouched.

- **Section order and git ranges.** Shared prefix (preamble + `files` table) is
  emitted for both paths. When `--full`: append `log[N]{sha,subject}` built from
  `git log <base>..HEAD` (**two-dot**) via `toon.Table`, then the raw diff block
  from `git diff <base>...HEAD` (**three-dot**). The two-dot/three-dot asymmetry
  is intentional (map watch-out #9: two-dot log = commits on HEAD since base;
  three-dot diff = changes on the HEAD side since merge-base) — do not normalize
  the two.

- **Log rendering.** A `git log <base>..HEAD` format yielding short-sha + subject
  per commit, one row per commit, into `toon.Table("log", {"sha","subject"}, …)`.
  Subjects are raw passthrough — `toon.Table` escapes the field a single layer,
  exactly as the `files` table does, so a comma/quote in a subject is safe.

- **Raw diff block.** `git diff <base>...HEAD` appended last as a single delimited
  block, introduced by a fixed marker line so an agent and a test can find its
  start. This block is a documented output contract (the `craft-cli` exemption for
  a surface with its own contract), documented in `bench diff` help text. No
  TOON-encoding (mangles `@@`/`+`/`-`) and no truncation (map #9: the reader wants
  it whole). The diff body follows the map's plain `git diff <base>...HEAD`; note
  for the build: the file list uses `--no-renames` while the body does not, so a
  rename shows as add+delete in the table but as a rename in the body — following
  the map; do not add `--no-renames` to the body without a reviewer call.

- **TOON conformance — uncertainty flag resolved by verification, not judgment.**
  The map (Handoff item 5/7) worried a generic "stdout is TOON-shaped" gate
  assertion would need to exempt the raw tail. **There is no such assertion.**
  Conformance and the AXI contract suite assert TOON shape with *positive
  per-line* checks (`requireOutputLine`, `requireAXILine`: a specific line is
  present) — no check parses the whole stdout as TOON. So the raw diff tail needs
  no gate exemption: the `--full` contract test asserts the structured prefix
  lines plus the raw `@@` markers, and nothing rejects the non-TOON tail. The
  documented-contract framing is satisfied by the help-text note, not a gate
  carve-out. **This discharges the map's uncertainty flag and deviates from its
  framing — flagged for veto.**

- **Phase-doc repoint.** Edit `.agents/commands/bench-review-implementation.md`
  step 1 (source of truth) to replace the `git diff <base>...HEAD` +
  `git log <base>..HEAD --oneline` instruction with `bench diff --full`; mirror
  the change byte-identically into `.claude/commands/bench-review-implementation.md`
  (the claude-mirror conformance check enforces parity). `bench-debug`'s
  `--grep=spec-retire` / `--diff-filter=D` lookups are specialized, not the
  base-relative bundle, and stay untouched (map #5).

## Testing decisions

- **What a good test is here.** Drive `bench diff --full` in a real git fixture
  and assert the stdout sections at the `diff.Command` seam (black-box) — never
  internals. Prior art: `internal/contract/axi/axi_wave2_test.go`
  (`testAXIDiffRecordedBase`, `testAXIDiffFallbackShape`, `testAXIDiffErrorPosture`).
- **Seam.** One: `diff.Command` stdout, exercised by AXI contract tests in
  `internal/contract/axi`.
- **Gate.** The project gate (`bench gate`). New behavior assertions land in the
  AXI contract phase; the story-5 anchor lands in the conformance phase
  (`internal/conformance/docs_workflow_helpers_test.go`).

### Seam diagram

    trigger: review agent (or /bench-review-implementation step 1) runs `bench diff --full`
        │
        ▼
    args: [--full]  ──▶  [ diff.Command                         ]  ──▶  stdout:
    repo state       ──▶  [  resolveBase → files table (shared) ]        branch/base/method
                          [  +full: log[N]{sha,subject} table   ]        files[N]{status,path}:
                          [  +full: raw `git diff base...HEAD`  ]        log[N]{sha,subject}:
                          [                                      ]        <marker> + raw diff body
                              ◀ tests attach here: AXI contract drives `bench diff --full`
                                in a git fixture and asserts the section lines + `@@` markers
                                on the Probe's stdout; a bare-`bench diff` probe asserts the
                                new sections are absent.

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | `bench diff --full` with a commit since base emits a `log[` table and a raw diff body | diff.Command stdout (AXI contract) | `bench diff --full` in a fixture with a commit since base — today `--full` is an unknown arg so it exits 2; the `RequireExit(0)` + `log[` assertion fails | A `--full` that is unrecognized, or that ignores the flag and prints only the file list, produces neither the `log[` header nor `@@` markers, so the row fails |
| 3 | log renders as `log[N]{sha,subject}:`; diff body is a raw block with `@@` hunk markers appended last | diff.Command stdout (AXI contract) | same fixture — assert the `log[` header, a `<sha>,<subject>` row, the block marker, and a line-start `@@ ` in the tail; exits 2 today | A TOON-encoded diff body escapes `@@`, so a line-start `@@ ` assertion fails; a missing/misordered log table fails the header/row assertion |
| 3 | a commit subject containing a comma and a quote is escaped exactly one layer in the `log` table | diff.Command stdout (AXI contract) | fixture commit subject `a, "b"` — assert `log` row renders `"a, \"b\""` (mirrors the existing `testAXITOONFieldEscaping`); exits 2 today | Double-escaping or no-escaping the subject field breaks the exact-once assertion, the same guarantee the `learnings` table is held to |
| 4 | empty-since-base yields `log[0]{sha,subject}:` and an empty diff block at exit 0 | diff.Command stdout (AXI contract) | fixture with base == HEAD (no commits since base) — `bench diff --full` exits 0 with an empty log table and empty body; exits 2 today | An implementation that errors on an empty range, or omits the sections, fails `RequireExit(0)` or the `log[0]` assertion |
| 4 | base-unresolvable exits 1 and unknown arg exits 2, unchanged under `--full` | diff.Command stdout (AXI contract) | `bench diff --full` in a repo with no merge-base exits 1 with `cannot resolve a review base`; `bench diff --full bogus` exits 2 `usage` (extends `testAXIDiffErrorPosture`) | A `--full` branch that resolves the base or parses args differently from bare `diff` would change the exit posture, failing the exit-code assertions |
| 2 | bare `bench diff` carries no `log[` section and no diff body | diff.Command stdout (AXI contract) | already covered / not TDD-able: bare `diff` never emitted these, so the row cannot start red; it is a regression guard — assert `RequireNotContains(stdout, "log[")` and no `@@` on a bare-`diff` probe | If a future change leaks the `--full` sections into bare `diff`, the not-contains assertion fails — the guard is the point |
| 5 | `/bench-review-implementation` step 1 names `bench diff --full` | `.agents/commands/bench-review-implementation.md` text (conformance) | tighten `docs_workflow_helpers_test.go:38` from `require("…bench-review-implementation.md", "bench diff")` to `"bench diff --full"` — red until the doc (both mirrors) is repointed | The anchor goes red if the phase doc is not repointed, pinning the payoff instead of leaving it to prose review |

Cheapest wrong implementations, confirmed caught: a `--full` that is accepted but
emits bare-`diff` output (rows 1/3 go red — no `log[`, no `@@`); a `--full` that
TOON-encodes the diff body (row 3 line-start `@@ ` fails); a `--full` that errors
on an empty range (row 4 exit-0 fails); a doc left un-repointed (row 5 anchor red).

### Edge inventory

Walked per behavior; each resolved as a coverage row above or a **Won't handle**
line here.

- error path → rows: base-unresolvable (exit 1), unknown arg (exit 2).
- empty/absent input → row: empty-since-base yields empty sections at exit 0.
- boundary values → covered: one commit vs many rows through `toon.Table` (row 1/3).
- malformed input → row: comma/quote in a commit subject escaped exactly once.
- interrupted/partial state — **Won't handle** — `bench diff --full` is a single
  read-only command with no partial state to resume.
- re-run idempotency — **Won't handle** — read-only and deterministic for a fixed
  tree; nothing to make idempotent.
- hostile environment (paths/subjects with spaces, non-ASCII, quotes) → the files
  table keeps its existing `-z` NUL framing (existing test); log subjects escape
  once (row above); the raw diff body is verbatim passthrough with no TOON layer,
  so **Won't handle** any escaping of the body — that is the documented contract.
- working-tree / uncommitted state — **Won't handle** — out of scope (see below);
  the flag is purely base-relative (map #3).

## Out of scope

- **Working-tree `git status` in the flag.** A separate capability — a mid-shift
  uncommitted-state view is a different context, and folding it in reintroduces
  the generic bundle the map killed (map #3). Park until it earns a trigger.
  Estimate to build later: ~2 edits, 2 gate runs.
- **Diff-body truncation / paging for very large diffs.** A separate capability;
  the map decided the opt-in flag returns the diff whole (map #9). If huge-diff
  ergonomics ever bite, spec a `--max-lines`-style cap then. Estimate: ~2 edits,
  2 gate runs.
