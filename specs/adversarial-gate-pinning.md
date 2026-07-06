# adversarial gate pinning (FT2)

Status: implemented

## Problem

The gate runs from the working tree by design (ADR 0001), which means the agent
being graded can weaken its own oracle. The canary tripwire makes a *lazy*
weakening loud, but ADR 0001 scopes out the adversarial case: an agent that
weakens `.bench/gate.sh` and deletes the tripwire fixtures in the same diff.
Nothing inside a writable tree can catch that — the check that would catch it is
itself writable.

## Solution

Pin the git tree hash of `.bench/` outside the writable tree and hash-verify it
at push time. A human runs `bench gate pin` once (TTY-interactive, typed
confirmation) to record the current `.bench` tree hash beside the existing gate
cache. The managed pre-push hook then recomputes each pushed commit's `.bench`
tree hash inline — pure git, no exec into tree-resident code — and blocks any
push whose enforcement surface drifted from the pin, with repin instructions.
Posture stays ADR 0001's: **loud, not impossible.** An agent with shell access
can edit `.git/hooks` directly; defending that is out of scope.

## User stories

1. As a reviewer, I want the pre-push hook to block a push whose pushed commit
   changed `.bench/` while a pin exists, so a weakened gate can't leave the
   machine unseen.
   Line: claude-opus-4-8 / medium. This is oracle-trust-domain logic where a
   wrong verifier silently defeats the whole feature, so correctness outranks
   speed.
2. As a reviewer, I want a push whose pushed `.bench/` tree matches the pin to
   pass cleanly, so legitimate work is never impeded.
   Line: claude-opus-4-8 / medium. Same inline-verify seam as story 1; the
   allow path and block path share one derivation and must be tested together.
3. As a consumer upgrading, I want an unpinned repo to pass with a one-line
   stderr notice naming `bench gate pin`, so pinning is opt-in and no relink
   breaks an existing consumer.
   Line: claude-opus-4-8 / medium. The unpinned branch of the same hook
   conditional; misjudging it hard-blocks every consumer on upgrade.
4. As a reviewer deleting a branch, I want a ref-delete push (zero local oid)
   skipped and allowed, so deletions aren't misread as an empty `.bench` tree.
   Line: claude-opus-4-8 / medium. A hostile-input branch of the hook loop that
   a naive `rev-parse` would crash on.
5. As a reviewer, I want a pushed commit that has no `.bench/` tree at all
   (while pinned) blocked, so deleting the enforcement surface wholesale is
   caught the same as mutating it.
   Line: claude-opus-4-8 / medium. The absent-tree branch of the verify block;
   distinct git error path from a hash mismatch.
6. As a reviewer, I want the existing default-branch block retained unchanged,
   so adding pin verification doesn't regress the guard already shipped.
   Line: claude-opus-4-8 / medium. Regression guard on the co-resident block in
   the same hook body.
7. As an agent reading `bench guards`, I want the hook's single `denies:` line to
   also name the pin denial, so the guard manifest lists every denial the hook
   actually enforces (one source per fact).
   Line: claude-opus-4-8 / medium. A guard row carries one `denies:` line
   (block-dangerous-git comma-separates its denials); the pin denial extends that
   line, and the `--describe` string and the enforcement branch — two derivations
   of one denial — must move together.
8. As a reviewer pushing several refs at once, I want any single `.bench`
   mismatch to block the whole push, so a mixed push can't smuggle a weakened
   gate past on a clean sibling ref.
   Line: claude-opus-4-8 / medium. The multi-ref loop's fail-closed aggregation,
   a hostile-input class the single-ref happy path hides.
9. As a reviewer, I want `bench gate pin` on a TTY to show the `.bench` diff from
   the pinned commit (or announce an initial pin), require typed confirmation,
   and write the pin from HEAD's committed tree, so re-pinning is a deliberate
   human-attended act.
   Line: claude-opus-4-8 / medium. The confirm path is human-only and gate-blind
   (see Testing decisions); it needs careful manual construction, not
   guess-and-check iteration.
10. As the system defending against a grinding loop, I want `bench gate pin` to
    refuse non-interactive stdin with a structured error and write no pin file,
    so a headless agent cannot repin after weakening the gate.
    Line: claude-opus-4-8 / medium. This refusal is the load-bearing security
    property of the whole feature and the one confirm-adjacent path the gate can
    drive.
11. As a reviewer, I want declining at the confirmation prompt to write nothing
    and exit non-zero, so an accidental `pin` invocation never silently records a
    bad hash.
    Line: claude-opus-4-8 / medium. The decline branch of the confirm path; a
    write-on-decline bug would defeat the human gate.
12. As a reviewer re-pinning after a legitimate gate change, I want a second
    `bench gate pin` to show the diff since the last pin and overwrite cleanly,
    so re-pin is idempotent and always reviewable.
    Line: claude-opus-4-8 / medium. Re-run idempotency on the pin writer; the
    diff base shifts from HEAD-only to last-pinned-commit→HEAD.
13. As a reviewer with uncommitted `.bench` edits, I want `bench gate pin` to
    warn about the dirty working tree but still pin from the committed tree, so
    the pin can never record an unpushable working-tree hash.
    Line: claude-opus-4-8 / medium. The committed-tree-not-working-tree contract
    (Handoff #9); the dirty-tree warning is the honesty signal.
14. As the kit dogfooding its own feature, I want this repo's managed pre-push
    hook reinstalled to carry the verify block in the same build, so the stale
    "no manifest" guards row disappears and `bench guards` shows a manifest row
    again.
    Line: claude-sonnet-5 / low. Mechanical reinstall at a known seam once the
    installer emits the new template.
15. As the teammate who walks in later, I want ADR 0001's "pinning is out of
    scope" paragraph rewritten to record the now-shipped state, so the decision
    record matches reality.
    Line: claude-fable-5 / high. Decision-record prose is the leverage-override
    surface in craft-line; it compounds through every session that reads it.

## Implementation decisions

- **Hook installer (`internal/adopt/link.go`, `installGitHook`).** The managed
  pre-push template gains an inline verify block *before* the existing
  default-branch loop. The block is pure bash + git — it must not exec `bench`
  plumbing, because the verifier has to sit in the same trust domain as the pin
  or an agent weakens the verifier instead of the gate (Handoff #5). It reads the
  pin file at `$(git rev-parse --git-path bench-gate-pin)`; absent → print the
  one-line stderr notice and fall through to allow; present → for each stdin ref
  line, skip zero-oid deletes, compute `git rev-parse "$local_oid:.bench"`
  (missing tree or mismatch against the pin's line 0 → stderr message naming
  `bench gate pin`, `exit 1`). The `--describe` branch's single `denies:` line is
  extended to name the pin denial alongside the default-branch one — a guard row
  aggregates exactly one `denies:` line (block-dangerous-git's comma-separated
  denials are the precedent), so a second line would be silently dropped. All
  ref/oid expansions stay quoted for paths and branches with spaces. The template
  stays an unexported string literal in `installGitHook`, so it remains the single
  source and tests re-spell substrings of it (today's convention).
- **Pin verb (`bench gate pin`).** New human-run porcelain. Shell wrapper: the
  `gate)` case inspects `$2`; `pin` routes to a new `gate-pin` porcelain
  subcommand (stdin passed through so TTY detection works), anything else keeps
  today's `run_gate`. Go: `cmd/bench/main.go` dispatches `case "gate-pin"` to a
  new `gate.PinCommand(args, os.Stdin, stdout, stderr)`. The verb detects a
  non-TTY stdin and refuses with a structured error (no pin written); on a TTY it
  shows `git diff <pinned-commit>..HEAD -- .bench` (or an initial-pin banner when
  no pin exists), requires a typed confirmation token, and on confirm writes the
  pin from `git rev-parse HEAD:.bench` — the *committed* tree, never the working
  tree. A dirty `.bench` working tree prints a warning but does not change what is
  pinned. Declining writes nothing and exits non-zero.
- **TTY detection.** No isatty/`term.IsTerminal` helper exists in the repo today;
  add one narrowly for the pin verb (stat `os.Stdin` for `ModeCharDevice`, the
  standard library-free check). It lives with the pin verb, not as a general util.
- **Pin file format (`<git-dir>/bench-gate-pin`).** Three newline-terminated
  lines, mirroring the line-keyed `bench-last-gate` cache precedent: line 0 = the
  `.bench` tree hash (the only line the hook reads), line 1 = the pinned-at commit
  oid, line 2 = an RFC3339 timestamp. Untracked, beside `bench-last-gate`,
  survives relink. The hook tolerates a malformed/short file by treating a
  missing line 0 as "no usable pin" and falling through to the unpinned notice
  rather than crashing.
- **Dogfood + ADR.** Reinstall this repo's managed hook via the same installer in
  this build so the kit carries the verify block. Rewrite ADR 0001's considered-
  options paragraph to record shipped pinning. FT10 (doctor *detects* a
  missing/stale kit hook) stays a separate capability — not folded in.

## Testing decisions

- **A good test here** drives the *installed* hook and the *built* verb as black
  boxes in a throwaway git repo — pipe ref lines to the hook's stdin and assert
  exit code + stderr; run `bench gate pin` with a controlled stdin and assert the
  pin file and exit code. Never assert on internal Go functions in place of the
  shell-observable behavior; the hook's byte shape is the contract.
- **Seams — two, both under `go test ./internal/contract/...` (one gate phase).**
  - *Hook behavior* → `internal/contract/surface`, the `contract.Fixture` harness
    `TestLinkContracts` uses to `bench link` into a temp repo and assert on the
    installed `pre-push`. Prior art: `link_test.go` installs the hook and asserts
    its text (`requireFixtureFileContains(... "bench:managed-pre-push")`);
    `RunAtWithInput` feeds ref lines to the installed hook's stdin. Today's
    pre-push tests are structural (assert text) plus one no-execute guard
    (`axi_guards_test.go`); FT2's hook contracts add the behavioral drive
    (pipe ref lines, assert exit code + stderr) at this seam.
  - *Pin-verb behavior* → `internal/contract/runtime`, beside the gate-cache
    tests that already read/write `<git-dir>` files via the `gitDir(t, f)` helper
    (`runtime_gate_test.go`). The verb writes `<git-dir>/bench-gate-pin`, so its
    contract belongs with the other git-dir-file behavior, not the link surface.
- **The one gate-blind seam: the TTY confirm path.** The fixture harness wires
  stdin as `strings.NewReader` — never a character device — so it can drive the
  **non-TTY refusal** path (story 10) directly but *cannot* drive the interactive
  confirm/decline path (stories 9, 11, 12, 13's confirm step). Decision: **test
  the refusal path in the gate; the confirm path is manual-verify, no PTY helper
  built.** Rationale — a PTY helper is a separate capability with no second caller
  in this repo, and the confirm path is human-only by design (an agent can't
  reach a TTY through the harness), so gate coverage of it buys nothing the
  refusal path and the write-from-committed-tree assertion don't already give.
  The write side of stories 9/12/13 (what lands in the pin file *given* a
  confirm) is tested by invoking the underlying pin-write with confirmation
  forced, so only the literal interactive prompt is manual.
- **Gate command:** `.bench/gate.sh` (unchanged; the new contracts ride the
  existing `go test ./internal/contract/...` phase).

### Seam diagram

Hook verify (stories 1–8):

    trigger: git push  (git feeds "local_ref local_oid remote_ref remote_oid" lines on stdin)
        │
        ▼
    stdin ref lines  ──▶  [ managed pre-push hook: inline bash+git      ]  ──▶  exit 0 (allow)
    <git-dir>/bench-  ──▶  [   per ref: skip zero-oid; rev-parse         ]  ──▶  exit 1 + stderr
      gate-pin        ──▶  [   local_oid:.bench; compare to pin line 0;  ]        (mismatch / absent)
    (absent → notice) ──▶  [   default-branch block retained             ]  ──▶  stderr notice (unpinned)
                              ◀ tests attach: bench link into a Fixture repo, then
                                RunAtWithInput pipes ref lines to the installed hook,
                                assert exit code + stderr

Pin verb (stories 9–13):

    trigger: human runs `bench gate pin`
        │
        ▼
    stdin (TTY?)     ──▶  [ gate.PinCommand: isTTY? no → structured   ]  ──▶  exit≠0, no file (non-TTY)
    HEAD:.bench      ──▶  [   error. yes → show .bench diff, read      ]  ──▶  <git-dir>/bench-gate-pin
    last pin (line1) ──▶  [   confirm token; write pin from committed  ]        (3 lines, on confirm)
    working-tree     ──▶  [   tree; dirty .bench → warn                ]  ──▶  exit≠0, no file (declined)
      dirtiness      ──▶  [                                            ]
                              ◀ tests attach: gate drives the non-TTY refusal via
                                Fixture stdin=pipe; confirm path is manual-verify

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | pinned + pushed `.bench` changed → block, exit 1, msg names `bench gate pin` | contract/surface fixture | new contract: link, pin (seed file), commit a `.bench` edit, pipe its ref line → currently allows (no verify block) | absent verify block → hook exits 0 on a drifted push; test goes red |
| 2 | pinned + pushed `.bench` unchanged → allow | contract/surface fixture | green today (no verify block, so a seeded pin is ignored and the push allowed); **over-block guard**, not red-first — goes red only if the added block over-blocks a matching push or reads the wrong hash source | a block that always blocks (or hashes the wrong tree) fails this row |
| 3 | unpinned → allow + one stderr notice naming the verb | contract/surface fixture | link (no pin file), pipe any ref → red: no notice emitted today | a hard-block-when-unpinned bug prints no notice / exits 1; row catches it |
| 4 | ref-delete (zero local oid) → skip, allow | contract/surface fixture | green today (zero-oid ref allowed, no verify block); **mis-skip guard**, not red-first — goes red only if the added block runs `rev-parse 0000:.bench` and blocks instead of skipping | `rev-parse 0000:.bench` errors → a naive block; the skip guard keeps this green |
| 5 | pinned + commit has no `.bench` tree → block | contract/surface fixture | pin, pipe a commit with `.bench` removed → red until absent-tree branch blocks | wholesale deletion of `.bench` must fail-closed, not be read as "no drift" |
| 6 | default-branch push still blocked | contract/surface fixture | existing assertion; re-run after adding the verify block → stays green (regression) | proves the new block didn't displace the retained default-branch guard |
| 7 | `--describe` lists the pin denial | contract/surface fixture + `bench guards` | assert `--describe` output has the pin `denies:` line → red until added | advertisement and enforcement are one fact; a missing line is drift |
| 8 | multi-ref push, one mismatch → whole push blocked | contract/surface fixture | pipe two ref lines in **both orders** — clean-then-drifted *and* drifted-then-clean — each must block → red until the loop fails closed on every ref regardless of position | a first-ref-only loop passes drifted-first but leaks clean-first (or vice-versa); testing both orders pins it |
| 9 (write) | confirm → pin file holds `HEAD:.bench` hash, not working-tree hash | contract/runtime fixture (confirm forced) | new contract: dirty `.bench`, force-confirm pin, read file line 0 == committed hash → red until write-from-committed wired | pinning the working tree records an unpushable hash; row pins the source |
| 10 | non-TTY stdin → structured error, exit≠0, no pin file | contract/runtime fixture | run `bench gate pin` with piped stdin → red: verb doesn't exist yet | the security property; a verb that pins headlessly defeats the feature |
| 11 (confirm prompt) | decline → no file, exit≠0 | manual-verify | not gate-drivable (TTY); stated exception | interactive decline can't be piped; write side covered by story 9's forced path |
| 12 | second pin → diff base = last-pinned commit; overwrites | contract/runtime fixture (write side) + manual (prompt) | pin twice via forced-confirm, assert line 1 updates and file overwrites → red until re-pin reads last pin | a re-pin that diffs against HEAD-only or appends instead of overwriting |
| 13 (warn) | dirty `.bench` → warning emitted, pin still from committed tree | manual-verify (warning text) + story 9 (write side) | warning is on the interactive path; committed-tree write asserted by story 9 | the honesty signal; its absence is cosmetic, the write correctness is story 9 |
| 14 | this repo's hook carries the verify block; `bench guards` shows a manifest row | `bench guards` on this repo | `bench guards` today reports "no manifest"; red until the hook is reinstalled | a stale hook leaves the guard unadvertised; the manifest row is the proof |
| 15 | ADR 0001 records shipped pinning (no stale "out of scope") | conformance stale-reference sweep | already covered: the docs sweep goes red on a dangling reference; ADR prose itself is not TDD-able | prose accuracy is reviewer-judged; the sweep only guards references |

Cheapest-wrong-implementation check: an always-allow hook (the sequential no-op)
fails rows 1, 5, 8; an always-block hook fails the over-block guard rows 2, 3, 4,
6; a first-ref-only loop fails row 8 (which tests both ref orders); a verb that
pins from the working tree fails row 9; a verb that pins headlessly fails row 10.
No degenerate implementation passes the map.

### Edge inventory

Walked per the profile's shell-CLI hostile-input checklist and the Handoff #6
owners:

- **paths/branches with spaces** — coverage: hook quoting exercised by the
  space-in-ref fixture folded into stories 1/8 (ref lines with a spaced branch).
- **zero oid on delete** — coverage row, story 4.
- **commit with no `.bench` tree** — coverage row, story 5.
- **multiple refs in one push** — coverage row, story 8.
- **absent vs present-but-empty vs malformed pin file** — absent → story 3;
  malformed/short (missing line 0) → **Won't handle as a distinct block**: treated
  as unpinned (falls through to the notice), asserted as an extension of story 3's
  fixture rather than its own row.
- **no `bench` on PATH** — the hook is inline bash+git and calls no `bench`;
  covered by construction (the verify block runs with an empty PATH in a fixture
  assertion under story 1).
- **re-run idempotency (pin twice)** — coverage row, story 12.
- **dirty `.bench` working tree** — coverage row, story 13 (+ 9 write side).
- **detached HEAD** — **Won't handle specially** — `git rev-parse HEAD:.bench`
  resolves in detached HEAD; the diff base still works. No special path; one
  fixture assertion under story 9 confirms it doesn't error.
- **repo with no commits** — **Won't handle** — `bench gate pin` in a repo with no
  commits has no `HEAD:.bench` to pin; the verb errors with the structured
  no-commits message. One assertion, folded under story 10's error-path fixture.
- **interrupt (SIGINT) mid-confirm** — **Won't handle** — declining-by-Ctrl-C is
  indistinguishable from decline; writes nothing (same as story 11). No leftover
  state because the write is a single final step. Safe to skip.

### Won't handle

- **`.git/hooks` tampering** — an agent with shell access can overwrite the hook
  itself; ADR 0001's posture is "loud, not impossible" and defending the hook
  file is explicitly out of scope.
- **malformed/short pin file as its own error** — treated as unpinned (notice +
  allow), not a hard block; a hard block here would strand a consumer whose pin
  file got truncated, contradicting the opt-in posture.

## Out of scope

- **FT10 — doctor detects a missing/stale kit hook.** A separate capability:
  `bench doctor` reporting (and offering to reinstall) a drifted managed hook on
  the kit repo itself. This spec only reinstalls the hook once, in-build; ongoing
  detection is FT10's own future spec. Estimate to build later: ~4 edits, 3 gate
  runs.
- **PTY helper for the interactive confirm path.** A test harness that allocates a
  pseudo-terminal so the gate could drive the confirm/decline prompt. A distinct
  testing-infrastructure capability with no other caller today; the confirm path
  is human-only and the write side is already gate-covered. Estimate if ever
  wanted: ~5 edits, 4 gate runs. (Low value — leave parked, do not roadmap.)
- **CI-side / server-side pin enforcement.** Verifying the pin on a remote or in
  CI rather than the local pre-push hook. Different trust domain and a separate
  capability; the local hook is the shipped scope. Estimate: unknown, needs its
  own grill.
