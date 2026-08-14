# worktree-enumeration-hang

Status: staged

Decision source: specs/worktree-enumeration-hang/decisions/worktree-enumeration-hang.md (compiled ready map, reviewer-resolved 2026-08-14)

Verification log: spec 33 + tickets 2 iteration(s) to accept — the spec loop grew the map from 11 to 24 rows; the largest catches were the linked-worktree degenerate, the typed-error rendering contract across three surfaces, and the fail-closed resolution path.

## Problem

A FIFO (or any non-regular file) planted as a worktree admin entry —
`gitdir`, `HEAD`, `commondir`, `locked` under `<git-common-dir>/worktrees/<id>/`
— makes `git worktree list --porcelain` block forever on git 2.43.0
(blocking open-for-read with no writer). Every Bench command that discovers
worktrees inherits the hang through the sole enumeration owner
`git.Worktrees`, including `bench status` and the session-start `bench resume`,
before any Bench guard runs — from the primary checkout and equally from a
linked pool worktree, whose admin dir is shared. The failure is a wedged
command with no attribution, and no Bench-owned deadline bounds it today.

## Solution

Discovery fails fast and says why. Before enumerating, Bench refuses any
malformed admin entry by shape, naming the exact entry and its observed kind
with the repair action; every git child the enumeration launches runs under a
named bounds-registry deadline as the backstop for blocking shapes the scan
cannot classify. `bench doctor` reports the same finding. Bench never deletes
anything under `<git-common-dir>/worktrees/` — removal stays a human act. The
mitigation
names the upstream git 2.43.0 behavior it works around so it retires if git
bounds its own admin reads.

## User stories

1. **Discovery refuses a malformed admin entry with attribution.** Running any
   discovering command (`bench status`, `bench worktree list`, the
   session-start resume) against a repository whose shared admin dir holds a
   non-regular entry — from the primary checkout or from a linked worktree —
   returns promptly with a structured refusal naming `worktrees/<id>/<entry>`
   and its observed shape, with "inspect and remove it" as the next action,
   and never launches `git worktree list`.
   Line: sonnet / medium. A filesystem shape scan at an existing seam, fully
   gate-observable through fixture repositories.
2. **The enumeration call is bounded.** With clean-shaped admin entries and a
   git that blocks anyway — even one that leaves a descendant holding stdout,
   and whichever of its two children (the common-dir rev-parse or the
   porcelain list) blocks — `git.Worktrees` returns a timeout refusal within
   the named `internal/bounds` deadline per launched child (worst case two
   windows), and the `bounds-policy` conformance check binds that consumption.
   Line: opus / medium. Child-process deadline handling and stdout-framing
   preservation are easy to get subtly wrong in ways only a well-built test
   observes.
3. **`bench doctor` reports the malformed entry.** A new doctor row goes red
   naming the offending admin entry on an adopted repo, leaving the entry
   untouched; a healthy repo prints no new row.
   Line: sonnet / low. One registry row plus one evaluator at the existing
   `doctorRows` seam.

## Implementation decisions

- **One shape owner.** A single scanner in `internal/git` owns the predicate
  the reviewer fixed: every entry under `<git-common-dir>/worktrees/` — each
  `<id>` directory and each entry directly inside it — must be a regular file
  or a directory. Exported contract (lands in the scanner ticket; the doctor ticket only
  calls it): `ScanWorktreeAdmin(commonDir string) error`, nil when `worktrees/` is
  absent, empty, or a non-symlink non-directory (git skips those shapes, probe
  asset) — but a **symlinked** `worktrees/` is refused, because git follows
  the link and hangs on a FIFO behind it (probe asset): symlinks fail open,
  every other non-directory fails closed to absence. Nothing in the admin
  tree is ever opened for read — the same "must never be opened" rule
  `internal/adopt`'s `isSpecialFile` states: the scanner first `os.Lstat`s
  `worktrees/` itself (never relying on `os.ReadDir`'s incidental
  ENOTDIR-on-FIFO behavior), then classifies with `os.ReadDir` /
  `entry.Type()` / `os.Lstat` only, exactly two levels deep — deeper content
  is git-tolerated (probe asset) and stays unscanned. `git.Worktrees` calls the scanner before launching
  `git worktree list`; the doctor row calls the same exported scanner rather
  than re-deriving the predicate.
- **Why not `internal/bounds/classify.go`.** The tree already owns a
  non-opening shape classifier there, but its `resolve` follows symlinks, so
  a symlinked `gitdir` with a regular target would classify healthy and the
  reviewer's regular-file-or-directory predicate (WE2) could never refuse it.
  The scanner needs lstat semantics, so the predicate lives in `internal/git`
  beside its only consumer — a considered non-reuse, not an oversight.
- **One common-dir owner.** The scan needs `<git-common-dir>`, resolved by
  `git rev-parse --path-format=absolute --git-common-dir` (clean under the
  hang from both primary and linked roots, probe asset). `internal/git`
  exports `CommonDir(root string) (string, error)`, which keeps today's
  plain `git.Output` runner — trimmed stdout only, git stderr never reaches
  the returned path, and the eight migrated callers keep their current
  unbounded failure modes exactly. What the two consumers share is the
  derivation fact: one unexported argv helper owns the rev-parse invocation
  form, consumed by `CommonDir` through `Output` and by `git.Worktrees`
  through the bounded stream-preserving variant. The resolution ticket lands `CommonDir`
  and the argv helper; the bound ticket's `git.Worktrees` drives the same
  helper through the variant — the helper is the stated contract between the two
  tickets. Bounded-versus-plain is deliberate behavior at two call sites,
  not a second copy of the derivation. `git.Worktrees` fails closed on the
  resolution (lands in the resolution ticket, WE19/WE20): a rev-parse error propagates, and a
  resolved path that does not `Lstat` as a directory is refused naming it —
  a failed or corrupted resolution can never silently skip the scan and
  fall through to the hanging enumeration. The resolution ticket lands first, so WE19 and WE20 go red against
  today's tree; the scanner ticket's baseline rows (WE1, WE2, WE3, WE6, WE12, WE13, WE18, WE24 — WE3's red
  needs the scanner still absent; its green adds the typed-refusal routing
  in `appendWorktree`) build on it, and WE4 alone reds only after the
  scanner lands. Eight production sites run that rev-parse today
  (`internal/intent/intent.go`, `internal/status`'s `isPrimaryCheckout`,
  two in `internal/worktree/subshell.go`, `internal/worktree/worktree.go`,
  `internal/worktree/classifier.go`, `internal/worktree/lifecycle.go`, and
  `internal/gate/engine.go`'s `commonGitDir` wrapper — the wrapper is deleted
  and its three callers call `git.CommonDir` directly, keeping one source);
  all eight migrate to `git.CommonDir` in this build — mechanical
  replacements priced at 8 sites / ~10 edits, 1 gate run, landing inside
  the resolution ticket's serial commit (`isPrimaryCheckout` keeps its separate
  `--git-dir` call and consumes `CommonDir` for the common half). The
  resolution then has one production Go source — three worktree test
  scaffolds also spell the rev-parse and stay outside this build's fences,
  a consolidation review may take up separately;
  the two bootstrap-layer shell derivations (`bin/bench.sh`'s
  `main_tree_kit` and the doctor-emitted shim) run before any Go code exists
  to call and stay as they are.
- **Refusal shape.** The scanner's error names the offending entry's
  repo-relative admin path — `worktrees/<id>/<entry>`, `worktrees/<id>`, or
  `worktrees` itself, whichever level offends — the observed shape word
  (`fifo`, `symlink`, `socket`, `device`), and the action
  `inspect and remove it`. Every failure `git.Worktrees` produces **except a
  genuine porcelain-child nonzero exit** is one exported typed error
  carrying a message and its own action: the scanner's shape refusal as
  above; the fail-closed resolution class (WE19's non-directory path,
  WE20's propagated rev-parse error) with the rev-parse failure text and
  the action `investigate the git failure`; and both bound expiries (WE17's
  rev-parse, WE9's porcelain) with the timed-out invocation, the enforced
  bound, and the same `investigate the git failure` action — a timeout must
  never advise re-running the invocation that just wedged. A start failure
  at either child (git absent, exec refused) joins the same typed class
  with the `investigate the git failure` action, so no status falls outside
  the partition into a silent nil result. `git.Worktrees` always supplies a
  `context.Background()` parent, so the `canceled` status is unreachable by
  construction; the consumption switch still maps it into the same typed
  class rather than a default nil. A rendering
  surface renders the typed class as its own framing — detail = the
  message, action = the type's action field verbatim; a per-surface action
  constant is the duplication defect review grades — and keeps its existing
  generic wrapper for untyped errors, which after this build are only
  porcelain-child exits, where retry advice is honest and the
  `git worktree list failed` label stays truthful.
  The typed error's `Error()` text itself carries all three elements —
  path, shape word, and the action clause — so every `%w`/`%v` surface
  renders the complete refusal with no routing; the structured action field
  exists for surfaces with a separate action cell. That yields the label
  rule: a surface that routes by class (status, `worktree list`) may keep a
  class-specific label on its untyped branch, while a surface that does not
  route must use a class-agnostic label. The session-start resume's one
  wrap site (`git worktree list failed: %w`), the dashboard's template
  label, and `PruneLandedBranches`' `git worktree list: %w` wrap (the one
  other class-specific label among the eight call sites) therefore all
  change, unconditionally, to the neutral `worktree discovery failed` —
  truthful for every class, no discriminator to game; resume's second, bare
  enumeration route already renders the full typed text unlabeled. One consequence the scanner ticket owns openly: the
  FT29 prior-art fixture (`chmod 000` on `.git`) becomes a typed resolution
  failure once the rev-parse lands first, so
  `TestAppendWorktreeSurfacesClassifyFailure`'s detail assertion moves to
  the typed framing while keeping its no-silent-empty contract — and its
  leading comment's "rather than a PATH-shimmed git" convention is
  deliberately reversed by this spec's stub fixtures, so the scanner ticket rewrites
  that comment too. The
  dashboard needs no routing — no action cell to misadvise — and its
  template takes the same neutral label swap above, so its framing stays
  truthful for every class it now receives. Existing
  callers propagate the error with `%w`/`%v` and keep the detail (verified
  for resume, intent, classifier, subshell, the dashboard, and the harness
  hook); `bench worktree list` today swallows the cause in a fixed
  "cannot read registered worktrees" / "run git worktree list and retry"
  message — the one advice string that names the exact invocation that
  wedges — and is edited to carry the typed error's detail and its action
  field verbatim, per class, never a surface constant.
- **The bound.** A new `WorktreeListTimeout = 15 * time.Second` constant joins
  the `internal/bounds` policy registry const block (15 s because the profile
  documents multi-second WSL2 VHDX fsync stalls a healthy host can hit).
  `internal/git` consumes it through the `internal/worktree/refresh/refresh.go`
  precedent — a package var initialized from `bounds.WorktreeListTimeout`,
  overridden to milliseconds in tests — so no test waits out the production
  value. Because WE21 and WE22 override the bound from other packages,
  `internal/git` also exports a test hook in the tree's existing shape
  (`diff.SetSnapshotAfterReadForTest` /
  `releaseevidence.SetArchiveMemberLimitForTesting` are the prior art) that
  sets and restores that var; it is the one cross-package knob, and
  production code never calls it. The bound is scoped to `git.Worktrees`' own children and applied
  per child — its rev-parse (the argv helper driven through the variant) and
  the porcelain list each get a fresh window, and the expiry error names
  which invocation
  timed out — so no child the enumeration launches is Bench-unbounded, while
  the eight plain `CommonDir` callers are untouched. `bounds.Run` merges
  stdout and stderr into one
  buffer, and the parser's NUL-framed stdout must stay unmixed, so
  `internal/bounds` gains a stream-preserving variant mirroring `Run`'s
  shape — same required `(parent, limit, cmd)` and the same five-status
  classification (complete / timeout / canceled / exit / start) — plus a
  caller-supplied stdout writer, keeping the `Setpgid` + process-group SIGKILL logic single-sourced
  in the registry package: `internal/git` re-deriving group kill would be a
  second copy of that knowledge. On expiry
  the refusal names the timed-out git invocation and renders the package
  var's live value — never the registry constant formatted independently, so
  an overridden test run cannot print a duration it did not enforce (the
  `refresh.go` precedent formats the constant there and must not be copied on
  that point).
  The cheapest wrong implementation — calling `bounds.Run` and parsing the
  merged buffer — is WE15's target, and substituting `bounds.Run` for the
  variant is the bound ticket's named mutation probe.
- **Conformance binding.** The `bounds-policy` check's `required` list gains
  `WorktreeListTimeout` and its `owners` map gains the `internal/git/git.go` →
  `bounds.WorktreeListTimeout` row. This is a substring tripwire, not a
  semantic proof (the `injected-port-registry` precedent frames it the same
  way); the behavioral half is WE9, and the var's default is pinned to the
  registry constant by WE14. The build proves the tripwire bites (red with
  the owner row present and the consumption absent) before relying on it.
- **Retirement trigger.** The scanner and the constant carry one comment
  naming the observed git 2.43.0 blocking-open behavior and the retire
  condition — remove the mitigation when git bounds its own admin reads. The
  probe matrix stays in the compiled map's assets as provenance.
- **The local probe bound stays local.** `internal/git`'s existing
  `refCheckTimeout` remains a package const: it is a hook-scoped fail-safe
  whose callers resolve a timeout to per-caller safe defaults, not a policy
  the registry owns; the enumeration bound is policy. One sentence of comment
  distinguishes them at the declaration.
- **No repair route.** Nothing deletes or rewrites entries under
  `<git-common-dir>/worktrees/` (reviewer decision #4; a plan/apply repair
  command was rejected). WE16 asserts the corrupted entry survives a doctor
  run.

## Testing decisions

- The behavior under test is external: a caller of `git.Worktrees` (or of a
  discovering command) either gets the worktree list promptly or gets a
  prompt, attributable refusal. Tests build real repositories and linked
  worktrees through `internal/git/git_test.go`'s existing helpers — the test
  census allows exactly one repository/process constructor site in
  `internal/git`, so story-1/2 tests reuse it and add no second
  `exec.Command` there. Admin entries are corrupted with `syscall.Mkfifo`
  under `capability.Capability(t, capability.Fifo, …)` and `os.Symlink` under
  the `Symlink` class (prior art: `internal/worktree/land_test.go`,
  `internal/worktree/classifier_shape_test.go`).
- Seams receiving tests: the `git.Worktrees` interface (primary; TDD), the
  `doctorRows` registry via `reportDoctorRows` (TDD; the fixture must chdir
  into an adopted repo, since it resolves `git.Root()` and requires
  `.bench/`), the `bench worktree list` command surface, the resume command
  surface (WE24), and `appendWorktree`
  in `internal/status` (prior art:
  `TestAppendWorktreeSurfacesClassifyFailure`; the seam sits below `render`'s
  five-row budget, which these rows leave unproven).
- A hang cannot be a test failure by itself, so every row whose failure mode
  is a block observes its red through a bounded wait: run the call in a
  goroutine and fail when the deadline expires before it returns. Scanner
  rows involve no inner wait, so their deadline is `bounds.TestDeadline(0)`
  (the floor); bound rows derive theirs from the test-overridden package var.
  These reds are first-write reds: no Go test fails on today's tree — the
  underlying hang is the reproduced shell probe
  (`timeout 5 git worktree list --porcelain` → exit 124, probe asset) — and
  each row's signal runs red immediately before its slice lands.
- One PATH-stub `git` design serves stories 1 and 2, and never execs real
  git, so no assertion races a real process against a milliseconds override.
  Its builder is pure file-writing (script + `t.Setenv` PATH, no process
  launch) and lives once in `internal/gittest` — whose package charter
  widens from repositories-only to the shared git test scaffolds — so the
  three consuming test packages (`internal/git`, `internal/status`,
  `internal/worktree`) share one harness instead of pasting three.
  It always appends its argv to a log file, and its `rev-parse` handling is
  argv-aware: `--show-toplevel` is always answered with the fixture root
  (command surfaces resolve the root before enumerating), while each mode's
  rev-parse behavior binds to the `--git-common-dir` invocation. It runs one
  fixture-selected mode: **block-worktree** answers `rev-parse` by printing
  the fixture's common dir (recorded into the stub at creation) and blocks on
  `worktree`, spawning a non-exec'd child that holds stdout — the child
  shares the pgid, so only a process-group kill releases the pipe (WE9,
  WE13); **block-rev-parse** blocks on `rev-parse` itself with the same
  non-exec'd stdout-holding child (WE17); **noisy-list** writes noise to
  stderr on **both** children — emitted before and interleaved with the
  stdout bytes, so a merged buffer cannot stay parseable by luck — while
  answering `rev-parse` and emitting valid NUL-framed porcelain on stdout,
  then exits 0 (WE15); **bad-rev-parse** answers `rev-parse` with a path
  that exists nowhere and serves `worktree` empty porcelain at exit 0
  (WE19); **fail-rev-parse** exits non-zero on `rev-parse` and serves
  `worktree` empty porcelain at exit 0 (WE20); **fail-worktree** answers
  `rev-parse` and exits non-zero on `worktree` — the untyped
  porcelain-failure contrast fixture (WE3, WE4); **vanish-after-rev-parse**
  answers `rev-parse`, then removes its own script so the porcelain child
  fails to start (WE23 — the two children run sequentially, so the second
  exec fails deterministically).
  Story 1 uses the stub for the non-invocation assertion (WE13,
  block-worktree), the corrupted-resolution refusal (WE19, bad-rev-parse),
  the failed-resolution propagation (WE20, fail-rev-parse), the typed
  contrast fixtures (WE3, fail-rev-parse; WE4, bad-rev-parse), and the
  untyped contrast fixtures (WE3/WE4, fail-worktree) — so the resolution
  ticket lands bad-rev-parse and fail-rev-parse (with the `internal/gittest`
  builder), the scanner ticket adds block-worktree and fail-worktree, and
  the bound ticket adds block-rev-parse, noisy-list, and
  vanish-after-rev-parse (WE23 is a story-2 row; both of its fixtures are
  asserted at the bound ticket's variant slice, where the PATH-empty failure lands
  on the variant's start arm). Story 1's other rows (WE1/WE2/WE6/WE12
  among the refusal rows, WE3/WE4/WE5/WE7/WE24 at the surface and
  over-refusal rows) run against real git, where the refusal or the normal
  enumeration is the observable.
- The gate observes everything through the ordinary `test` phase; WE8
  additionally runs the `bounds-policy` check red before the consumption
  lands.

### Seam diagram

    trigger: any discovering command (status, resume, worktree *, dashboard)
        │
        ▼
    root ──▶ [ git.Worktrees                      ] ──▶ []Worktree, or a refusal
             │  1. rev-parse via the argv helper     │   naming the offending admin
             │     (bounded variant; fail-closed)    │   path,
             │  2. ScanWorktreeAdmin: lstat-only     │   its shape, and the action —
             │     regular-file-or-dir predicate     │   or a timeout refusal naming
             │  3. porcelain -z (bounded, same       │   the timed-out invocation
             │     WorktreeListTimeout-backed var)   │   and the bound
             └───────────────────────────────────────┘
                  ◀ tests attach here: fixture repo (primary and linked roots)
                    with one corrupted admin entry; assert prompt attributable
                    error and no `git worktree` launch, never a hang

    trigger: bench doctor (adopted repo)
        │
        ▼
    root ──▶ [ doctorRows: admin-entries evaluator ] ──▶ red row naming the entry,
             (delegates to git.ScanWorktreeAdmin)        entry left untouched, or
                                                         no row when healthy
                  ◀ tests attach here: reportDoctorRows over the fixture repo

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| WE1 | 1 | a FIFO `gitdir` admin entry makes `git.Worktrees` (real git) return, within the bounded wait, a refusal whose message contains `worktrees/<id>/gitdir`, `fifo`, and `inspect and remove it` | `git.Worktrees` interface | first-write red: the call blocks (shell probe exit 124) and the bounded wait expires | a hang and an unattributed error both fail the name-shape-action assertions |
| WE2 | 1 | a first-level symlink (`worktrees/<id>` itself a link), a first-level FIFO at `worktrees/<id>`, and a second-level symlinked `gitdir` with a regular target are each refused with a message containing the offending path and its shape word (all three fixtures asserted) | `git.Worktrees` interface | first-write red: git follows the links and skips the FIFO (exit 0 each, probe asset), `Worktrees` returns no refusal, the contains-assertions fail | catches the FIFO-only degenerate, the level-1 branch that refuses symlinks but skips other non-regular entries, the scanner that classifies only inside `<id>` dirs, a first-level refusal that drops its attribution, and a level-2 classifier that stats through links — which admits the symlink-to-FIFO `gitdir` that hangs git (probe asset) |
| WE3 | 1 | with the FIFO planted, `appendWorktree` returns a row whose detail contains `worktrees/<id>/gitdir` and `fifo` and whose action contains `inspect and remove it`; with the stub in fail-rev-parse mode, the detail does **not** contain `git worktree list failed`, and the action contains `investigate the git failure` and not `inspect and remove it`; with the stub in fail-worktree mode (untyped porcelain exit), the row keeps the `git worktree list failed` detail and the re-run action (all three fixtures asserted) | `appendWorktree` (internal/status) | first-write red twice over: the FIFO fixture's row-building call blocks and the bounded wait expires, and the fail-rev-parse fixture's no-`git worktree list failed` assertion fails against the generic wrapper; the fail-worktree fixture is the untyped branch's green guard | proves attribution survives the ambient surface, and the contrast fixtures force a real per-class discriminator — an unconditional action swap, a substring sniff, or a single hardcoded typed action each fails one fixture |
| WE4 | 1 | with the FIFO planted, `bench worktree list` returns a structured error containing `worktrees/<id>/gitdir` and `inspect and remove it`; with the stub in bad-rev-parse mode (garbage common dir passes root resolution and the empty intent ledger, then `git.Worktrees` refuses), its error carries the resolution class's detail naming that path and an action containing `investigate the git failure`, not `inspect and remove it`; with the stub in fail-worktree mode, it keeps the fixed "cannot read registered worktrees" / retry message (all three fixtures asserted) | `worktree list` command surface | first-write red: once `git.Worktrees` refuses, the command still prints the fixed message and the typed-side contains-assertions fail on both typed fixtures; the fail-worktree fixture is the untyped branch's green guard | catches the one caller that swallows the cause, pins the per-class action strings at this surface (field provenance stays the review-checked exception), and kills the unconditional swap that would advise removing admin entries on every generic failure |
| WE24 | 1 | with the FIFO planted (real git), the session-start resume command's output contains `worktrees/<id>/gitdir`, `fifo`, and `inspect and remove it`, and does not contain `git worktree list failed` | resume command surface (internal/worktree) | first-write red: the resume's enumeration blocks and the bounded wait expires | proves the story's third named surface carries the complete refusal text and the neutral label — not a claim that the never-launched invocation failed |
| WE5 | 1 | a repository whose `worktrees/` admin dir is absent — or is a FIFO, which the scanner treats as absence — enumerates normally (both fixtures asserted) | `git.Worktrees` interface | not TDD-able: guards the new scanner against over-refusal, so it cannot start red for the right reason; lands green with the same ticket | catches a scanner that errors on the absent directory every fresh repo has, or blocks on the FIFO shape git itself skips |
| WE6 | 1 | a FIFO at a name outside git's read set (`stray`), under an `<id>` containing a space and a glob character (`x y*`), is refused naming the id verbatim | `git.Worktrees` interface | first-write red: git tolerates the stray FIFO (exit 0, probe asset), `Worktrees` returns no refusal, the contains-id assertion fails | catches the scanner that lstats only git's known entry names, plus quoting, path-joining, or `filepath.Glob`-based scanning that mangles or skips the attribution |
| WE7 | 1 | well-shaped or deep states — a prunable entry (worktree dir deleted), an `<id>` dir whose `gitdir` was removed, an empty `<id>` dir, and a depth-3 FIFO at `logs/HEAD` — all still enumerate (all four fixtures asserted) | `git.Worktrees` interface | not TDD-able: over-refusal guard on the new predicate (git tolerates all four at exit 0, probe asset); lands green with the same ticket | catches a scanner that refuses the states the refusal's own "remove it" instruction produces, and a recursive walk that scans past the two-level predicate |
| WE8 | 2 | removing the `bounds.WorktreeListTimeout` consumption turns the `bounds-policy` check red | conformance check (`internal/conformance`) | observed red: the check runs with the owner row present before the consumption lands and reports `internal/git/git.go does not consume bounds.WorktreeListTimeout` | substring tripwire binding the deadline to the policy registry; WE9 and WE14 carry the behavioral halves |
| WE9 | 2 | with clean-shaped admin entries and the PATH-stub git blocking on `worktree` while a non-exec'd child holds stdout, `git.Worktrees` returns within the bounded wait an error containing the timed-out invocation and the enforced bound value — the test-overridden duration, not the registry constant's 15s | `git.Worktrees` interface | first-write red: the bounded wait expires while the stub still blocks | catches a missing deadline and a deadline that kills only the direct child while a descendant holds the pipe |
| WE10 | 3 | on an adopted repo with a FIFO admin entry, `bench doctor` prints a red row containing `worktrees/<id>/gitdir` | `reportDoctorRows` (internal/adopt) | first-write red: no such row exists, so the contains-assertion fails | proves the reviewer-chosen doctor surface reports the finding |
| WE11 | 3 | on an adopted repo with healthy worktrees, `bench doctor` prints no admin-entry row | `reportDoctorRows` (internal/adopt) | not TDD-able: false-positive guard on the new row; lands green with the same ticket | catches an evaluator that is always red or reports on healthy repos |
| WE12 | 1 | from a **linked worktree's** root with the FIFO planted in the shared admin dir, `git.Worktrees` (real git) returns the same attributable refusal within the bounded wait | `git.Worktrees` interface | first-write red: from a linked root the enumeration hangs identically (exit 124, probe asset) and the bounded wait expires | catches the cheapest wrong scanner — `<root>/.git/worktrees` — which scans nothing from a linked checkout, where `.git` is a file, and leaves pool sessions wedged |
| WE13 | 1 | with the FIFO planted and the PATH-stub git logging argv, the refusal contains `worktrees/<id>/gitdir`, `fifo`, and `inspect and remove it`, and the log contains no `worktree` invocation | `git.Worktrees` interface | first-write red: without the scanner the stub's `worktree` call blocks, the bounded wait expires, and the log records the `worktree` launch | proves refusal-before-execution — the scan verdict, not some earlier failure, lands before the hanging invocation can start |
| WE14 | 2 | the enumeration package var's default equals `bounds.WorktreeListTimeout` | `internal/git` package assertion | not TDD-able: registry-pinning assertion on the new var; lands green with the same ticket | catches the var drifting from the registry value while the substring tripwire (WE8) stays green |
| WE15 | 2 | with the stub in noisy-list mode (stderr noise on both children), `git.Worktrees` parses the correct worktree list | `git.Worktrees` interface | not TDD-able as a pre-implementation red (today's `Raw` already preserves stdout); the bound ticket's named mutation runs as two single-site `bounds.Run` substitutions, each redding this row — at the porcelain child the merge corrupts the NUL-framed parse; at the rev-parse child the merged noise makes the resolved path a non-directory, the fail-closed refusal fires, and the parse assertion fails | catches the cheapest wrong bound implementation, one child per probe site |
| WE16 | 3 | after the `bench doctor` run over the FIFO fixture, the FIFO admin entry still exists (lstat) | `reportDoctorRows` (internal/adopt) | not TDD-able: guards the rejected repair route, so it cannot start red honestly; lands green with the same ticket | catches any path where reporting mutates git-owned state, honoring reviewer decision #4's no-delete promise |
| WE19 | 1 | with the stub in bad-rev-parse mode (resolved path exists nowhere), `git.Worktrees` returns a refusal naming that path, and the argv log contains no `worktree` invocation | `git.Worktrees` interface | first-write red against today's tree (the resolution slice's shared baseline): no rev-parse runs, the porcelain child answers empty, and `Worktrees` returns an empty list with no refusal — both assertions fail | catches an implementation that swallows a corrupted resolution or launches the enumeration before the resolution verdict lands |
| WE21 | 2 | with the stub in block-worktree mode and the bound overridden to milliseconds, `appendWorktree` returns a row whose detail names the timed-out invocation, and whose action contains `investigate the git failure` with no re-run advice | `appendWorktree` (internal/status) | not TDD-able: the scanner ticket's typed routing is generic over the class, so the row arrives green when the bound ticket's typed timeout lands; guards routing generality | catches an implementation that routes only the shape refusal and lets a bound expiry fall through to the generic wrapper's re-run advice — the exact advice a timeout must never give |
| WE22 | 2 | with the stub in block-worktree mode and the bound overridden to milliseconds, `bench worktree list` returns a structured error naming the timed-out invocation, whose action contains `investigate the git failure` and not `inspect and remove it` or the retry advice | `worktree list` command surface | not TDD-able: the scanner ticket's typed routing is generic over the class, so the row arrives green when the bound ticket's typed timeout lands; guards routing generality | catches the route-only-shape degenerate at the second rendering surface, where a bound expiry would otherwise print "run git worktree list and retry" — retry advice on the wedged invocation |
| WE23 | 2 | a git that cannot start yields the typed error carrying the exec failure text and naming the failing invocation — `worktree list` under vanish-after-rev-parse mode, `rev-parse` under a PATH with no git — with the `investigate the git failure` action, never a silent empty list or a borrowed timeout message (per-fixture assertions on the named invocation) | `git.Worktrees` interface | first-write red on the vanish fixture at the bound ticket's variant slice: the porcelain child still runs through unbounded `Raw`, so its failed exec surfaces untyped and the typed-action assertion fails; the PATH-empty fixture arrives green off the resolution ticket's typed class and stays as the routing-generality guard for the variant's rev-parse start arm — where a default-nil start collapses into WE19's fail-closed refusal, losing the exec failure text and the `rev-parse` invocation name the per-fixture assertions require | catches a variant-consumption switch whose `start` arm defaults to nil — zero worktrees reported silently where today's tree at least propagates the exec error |
| WE20 | 1 | with the stub in fail-rev-parse mode (rev-parse exits non-zero), `git.Worktrees` returns the typed resolution failure naming the rev-parse invocation, and the argv log contains no `worktree` invocation | `git.Worktrees` interface | first-write red against today's tree (the resolution slice's shared baseline): no rev-parse runs, the porcelain child answers empty, and `Worktrees` returns an empty list with no error — both assertions fail | catches an implementation that maps a failed resolution to "no admin dir" or launches the enumeration before the resolution verdict lands |
| WE18 | 1 | `worktrees/` itself a symlink is refused with a message containing `worktrees` and `symlink` | `git.Worktrees` interface | first-write red: git follows the link (exit 0 healthy / exit 124 with a FIFO behind it, probe asset), `Worktrees` returns no refusal, the contains-assertions fail | catches the fail-open shape a treat-all-non-directories-as-absence scanner would silently skip while git walks through it |
| WE17 | 2 | with the stub in block-rev-parse mode while a non-exec'd child holds stdout, `git.Worktrees` returns within the bounded wait a timeout refusal naming the rev-parse invocation and the enforced bound value — the test-overridden duration, not the registry constant's 15s | `git.Worktrees` interface | first-write red: the rev-parse child blocks and the bounded wait expires while the stub still blocks | catches a bound applied only to the porcelain child, and a rev-parse runner without process-group kill that a surviving descendant would outlive |

### Edge inventory

- Error path — WE1–WE4, WE9, WE10, WE12, WE13, WE17, WE18, WE19, WE20,
  WE21, WE22, WE23, WE24. Empty/absent
  input — WE5 (absent or FIFO `worktrees/`), WE7 (present-but-empty `<id>`).
  Boundary values — WE6. Malformed input — WE1, WE2, WE6, WE12, WE18.
  Interrupted/partial state — WE7. Hostile environment — WE9 (PATH-stubbed
  git with a surviving descendant), WE15 (stderr noise), WE17 (blocking
  rev-parse), WE23 (git that cannot start). Process-boundary lifecycle — the
  refusal derives from a stateless read each invocation; nothing Bench
  serializes crosses a process boundary, and the child-process edges are
  WE9, WE15, WE17, WE19, WE20, and WE23.
- **Exception (decision-source promise without a row):** the retirement
  trigger is a source comment the gate cannot grade; review checks it against
  the map's Destination.
- **Deviation flagged for reviewer veto:** decision #2's answer names
  `bounds.Run` as the vehicle; `bounds.Run` merges stdout into stderr and
  would corrupt the NUL-framed parse, so the build uses a stream-preserving
  variant added to the same `internal/bounds` owner instead. Behavior and
  the bound's registry ownership are unchanged; only the named API differs.
- **Exception (review-checked property):** per-child windows — each of the
  two children getting a fresh deadline rather than sharing one span — has no
  row: distinguishing it needs elapsed-time assertions the suite would flake
  on. Review owns checking it in the diff; WE9 and WE17 prove each child is
  bounded at all.
- **Exception (review-checked property):** action-field provenance — WE3,
  WE4, WE21, and WE22 pin the per-class action strings at both surfaces,
  but only review can check the strings are read from the type's action
  field rather than re-derived as per-surface constants (the duplication
  defect the code standard names).
- **Exception (review-checked property):** the `canceled` arm of the
  variant-consumption switch — unreachable by construction under the
  `context.Background()` parent, so no fixture can produce it. Review owns
  checking the switch maps it into the typed class rather than a default
  nil; the `start` arm is fixture-reachable and carries WE23.
- **Exception (no per-caller rows):** the eight production `git.Worktrees`
  call sites are `PruneLandedBranches` (internal/git), `subshell.go`,
  `ClassifyRegisteredWorktrees` (worktree.go), `classifier.go`, `resume.go`,
  `list.go`, `intent.go`, and `harness/worktree.go`; all but `list.go` (the
  one swallower, WE4), `resume.go` (a relabeling propagator, WE24), and
  `PruneLandedBranches` (whose `git worktree list:` wrap takes the neutral
  relabel in the resolution ticket, inside its `internal/git/` fence, with no row — the
  one-line change rides the resolution slice and review checks it) propagate the
  refusal with `%w`/`%v`/bare returns under neutral labels (verified at
  review, the dashboard and status rendering included). WE1 at the sole
  producer plus WE3/WE4/WE24 at the three rendering surfaces cover the
  composition.
- **Won't handle:** socket or device admin entries — the same non-regular
  branch WE1/WE2 already prove; a device node needs privilege every host test
  would skip, and a socket buys a test dependency for an identical code path.
- **Won't handle:** control bytes in an admin entry id — the toon layer's
  control-byte policy (owned by `single-control-escaper`) refuses rendering,
  so the refusal degrades to a generic structured error and stays fail-closed.
- **Won't handle:** admin content deeper than two levels (e.g. `logs/HEAD` as
  a FIFO) — clean at exit 0 on git 2.43.0 (probe asset); the story-2 bound is
  the declared backstop for read-set drift.
- **Won't handle:** re-run idempotency — the scan and the refusal are
  stateless reads that write nothing, so repetition cannot diverge.
- **Won't handle:** worktree-mutating git calls (`add`/`lock`/`unlock`/
  `prune`) — reviewer decision #3 scoped this build to enumeration; the
  exposure is parked in `capture/IDEAS.md`.

## Ownership fences

One writer per ticket. The graph is resolve → refuse → {bound, doctor};
the two frontier tickets share no fence path (`internal/adopt/` is disjoint
from the bound ticket's paths), and every shared path sits on the serial
spine, so no path ever has two writers at once.

- `resolve-git-common-dir`: `internal/git/`, `internal/gittest/`,
  `internal/intent/intent.go`, `internal/worktree/subshell.go`,
  `internal/worktree/worktree.go`, `internal/worktree/classifier.go`,
  `internal/worktree/lifecycle.go`, `internal/worktree/resume.go`,
  `internal/gate/engine.go`, `internal/dashboard/render.go`,
  `internal/status/status.go`
- `refuse-malformed-admin-entries`: `internal/git/`, `internal/gittest/`,
  `internal/worktree/list.go`, `internal/worktree/list_actions_test.go`,
  `internal/worktree/resume_test.go`, `internal/status/status.go`,
  `internal/status/status_test.go`
- `bound-worktree-enumeration`: `internal/git/`, `internal/gittest/`,
  `internal/bounds/`, `internal/conformance/bounds_policy_test.go`,
  `internal/status/status_test.go`, `internal/worktree/list_actions_test.go`
- `report-admin-entry-in-doctor`: `internal/adopt/`

## Out of scope

- Fronting the worktree-mutating git call sites with the same pre-scan
  (`internal/worktree/{ownership,lifecycle,reauthorize,snapshot}.go`,
  `internal/gate/engine.go`) — a separate capability parked in
  `capture/IDEAS.md`; roughly 10 edits, 4 gate runs.
- A plan/apply repair command for malformed admin entries — rejected by
  reviewer decision #4, not deferred.
- Fixing or patching upstream git.
