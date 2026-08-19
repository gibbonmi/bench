# ft234-pool-key-reclaim

Status: staged

Decision source: roadmap/FT234.md, opened by the reviewer in conversation on 2026-08-19 immediately after the FT226 landing, with the `/bench-debug` repro recorded below as its root-cause evidence and the FT226 build verification log (`specs/ft226-test-home-isolation/spec.md`, SW2 survivor list) naming the eight surviving keys this spec's population is drawn from. Those keys are deliberately held back unswept as this build's acceptance fixture.

Verification log: <n> iteration(s) to accept — pending the spec-and-tickets review round.

## Problem

`$BENCH_HOME/worktrees` accumulates keys that nothing can ever reclaim.

A pool key is named by `Pool(root)` — `basename(root)` plus a checksum of the
root path — so the mapping runs one way only, from a live repository to its key.
Every reclamation path Bench has is anchored at such a root: `ConservativeCleanup`
takes `root`, and its three inputs are that repository's git registrations
(`git worktree list`), that repository's intent ledger, and that repository's
`refs/bench/…` namespaces. The ledger itself lives at
`<git-common-dir>/bench-intent.json` — inside the source repository. So when the
source repository is deleted, the only record of its assignments is deleted with
it, and no `root` exists that resolves back to the key. The key is unreachable,
not merely unswept.

The operator's pool holds nine keys today. Five are dogfood keys whose single
child carries a `.git` pointer file naming a `/tmp/bench-dogfood-*` repository
that no longer exists; three are empty directories; one is this build's own
integration key. `bench session-inspect` runs at every session start and reports
`removed 0, swept refs 0; pruned branches 0; reconciled 0` against that pool.

**This is a missing capability, not a broken predicate.** The existing orphan
machinery answers the opposite question. `OrphanCandidate` and
`sweepOrphanAssignments` walk *ledger records* and report the ones whose worktree
is still on disk; `reconcileLifecycleDebris` drops *ledger records* whose checkout
is gone. Both need a record. The population here is the mirror image: bytes on
disk with no record and no repository left to hold one. Nothing in production
code ever enumerates `$BENCH_HOME/worktrees` at all — the only `ReadDir` calls in
the package read a single already-resolved key or a repository's private git
administration directory, and the only reader of the pool parent in the tree is
`homeResidue`, a test-only detector FT226 added.

### Repro (`/bench-debug` Phase 1, run before any theory)

A synthetic pool under a temporary `BENCH_HOME` (never the operator's), holding
one key whose child `.git` pointer names a deleted repository and one empty key,
with `bench resume-clean` driven from a separate live repository:

    --- pool before
    dead-repo-1596407047
    empty-key-123
    --- bench resume-clean, run from the live repo
    bench resume: removed 0, swept refs 0; pruned branches 0; reconciled 0; failed 0; open assignments 0
    exit=0
    --- pool after
    dead-repo-1596407047
    empty-key-123
    --- VERDICT
    RED: both unreclaimable keys survive resume-clean

`bench session-inspect`, the SessionStart surface, prints the same line.

The discriminating experiment separates "unreachable key" from "wrong
predicate": a dangling-pointer child planted inside the key of a repository that
is **alive and current** also survives untouched. Reachability is therefore not
the whole gap — nothing walks pool bytes under any circumstances, so no
predicate change to `poolAssignment`, `orphaned`, or `sweepOrphanAssignments`
can reach this population. The fix is a new reader over the pool parent.

## Solution

One opt-in, plan-before-apply verb: `bench worktree reclaim`.

Bare, it plans: it enumerates `$BENCH_HOME/worktrees`, classifies every key
against one reclaimability predicate, and prints a TOON table of keys with a
verdict and a reason, an aggregate count, and the exact apply invocation carrying
the plan's fingerprint. It removes nothing.

`bench worktree reclaim --apply <fingerprint>` removes the keys that plan named,
refusing when the pool has drifted since the plan was printed, and re-checking
the predicate against each key immediately before removing it.

A key is reclaimable when it holds nothing at top level, or when every top-level
entry is a directory holding a regular `.git` file whose `gitdir:` target is
provably absent, and nothing else. Anything else — a live pointer, a `.git` that
is a real repository directory, a `.git` with no `gitdir:` line, a stray
top-level file, a key mixing a live and a dead pointer, a symlink anywhere on
the path, or any stat error that leaves existence unknown — retains the key and
says why.

Reclamation never happens automatically. `bench resume-clean` gains one reported
count and the command that acts on it, and keeps its hands off the bytes.

## User stories

### Planning the pool

Line: `opus` / medium. Exact predicate over a small filesystem shape, one new
file in a package whose conventions are established; the predicate is itself
the safety of a destructive command, so the cheap row's bump applies.

1. As an operator, I want to see every key in my pool that can never be
   reclaimed, so that I can recover the disk a deleted repository left behind.
2. As an operator, I want each key's verdict to say *why* it survived, so that a
   key I expected to be reclaimed tells me what protected it instead of leaving
   me to guess.
3. As an operator, I want the plan to name the exact apply command with its
   fingerprint, so that acting on the plan needs no invention.
4. As an operator, I want an empty pool, or a pool with no reclaimable key, to
   print a definitive zero-row result and exit zero, so that "nothing to do" is
   an answer rather than silence.
5. As an agent, I want the plan on stdout as TOON with a small row schema, so
   that reading it costs few tokens.
6. As an operator, I want the plan to remove nothing at all, so that inspecting
   my pool is never destructive.
7. As an operator, I want a key holding a live worktree — a `.git` pointer whose
   target exists — to be retained, so that the command cannot destroy work.
8. As an operator, I want a key whose child `.git` is a real repository
   directory rather than a pointer file to be retained, so that a repository
   somebody parked in the pool survives.
9. As an operator, I want a key whose child `.git` file carries no `gitdir:`
   line to be retained, so that an unparseable pointer is never read as proof of
   absence.
10. As an operator, I want a key holding any stray top-level entry that is not a
    child directory to be retained, so that anything I put there by hand is safe.
11. As an operator, I want a key mixing one live and one dead pointer to be
    retained whole, so that partial reclamation cannot amputate a live worktree.
12. As an operator, I want any stat or read error to retain the key and be
    reported, so that unknown is never treated as absent.
13. As an operator, I want a symlinked key, a symlinked child, or a symlinked
    `.git` to be retained unfollowed, so that the command cannot be aimed at
    bytes outside the pool.
14. As an operator, I want the key of the repository I am currently in to be
    retained unconditionally, so that a session between acquiring its pool
    directory and its first checkout is never swept out from under itself.

### Applying the plan

Line: `opus` / high. The subject is the operator's real pool, outside any tree
the gate observes and with no undo; FT226's uncovered-gate bump applies, and the
effort goes into the refusals rather than the removal.

15. As an operator, I want removal to require the fingerprint the plan printed,
    so that a destructive step is always a second, deliberate invocation.
16. As an operator, I want the apply to refuse when the pool changed since the
    plan, naming the re-plan command, so that I never remove a key on the
    strength of a stale reading.
17. As an operator, I want each key re-checked against the same predicate
    immediately before it is removed, so that a key that became live between
    plan and apply survives.
18. As an operator, I want the apply to report per key what it removed and what
    it retained, plus the counts, so that the destructive step leaves evidence.
19. As an operator, I want an apply against a plan with no targets to exit zero
    having done nothing, so that a repeat invocation is a safe no-op.
20. As an operator, I want a malformed or absent fingerprint to be a usage
    error, so that the handshake cannot be fumbled into a removal.
21. As an operator, I want removal to touch only direct children of
    `$BENCH_HOME/worktrees`, so that no path outside the pool is ever a target.

### Knowing it is there

Line: `opus` / medium. One count and one action line in an existing renderer.

22. As an operator, I want `bench resume-clean` to report how many pool keys are
    reclaimable and name the command that reclaims them, so that I learn about
    the debris without hunting for a verb.
23. As an operator, I want session start never to remove a pool key on its own,
    so that an unattended hook is never destructive outside the tree it opened in.
24. As a maintainer, I want the resume count and the command's plan to come from
    one predicate, so that the ambient number can never disagree with what the
    verb would do.

## Implementation decisions

- **A new verb, not an extension of `bench worktree clean`.** `clean` is
  path-addressed and repository-scoped: it resolves a git registration, consults
  the ledger, and preserves dirty work before removing anything. Pool
  reclamation is home-scoped and repository-independent, and its whole premise is
  that no registration and no ledger record exist. Overloading `clean` would put
  two unrelated subjects behind one operand grammar. `reclaim` takes the
  `--apply <fingerprint>` grammar `clean` already established, so the
  plan-before-apply handshake is one convention rather than two.
- **One predicate, one file.** `internal/worktree/pool_reclaim.go` owns the
  reclaimability predicate and nothing else derives it: the plan renderer, the
  apply path, and the resume count all call it. The residue detector FT226 added
  in `main_test.go` answers a different question (*is anything at all under a
  private home*, where a legitimate entry does not exist) and stays separate; it
  is not a second derivation of this fact and must not be collapsed into it.
- **Proven absence only.** The dangling test is an `os.Lstat` of the `gitdir:`
  target returning `IsNotExist`. Every other error, and every unreadable entry
  anywhere on the key's walk, retains the key with a reason. Symlinks are never
  followed: `Lstat` throughout, and a symlink found where a key, a child, or a
  `.git` file belongs retains the key.
- **The current repository's key is excluded before the predicate runs.** A live
  session sits between `MkdirAll(pool)` and its first `worktree add` with an
  empty key on disk, which the empty-key clause would otherwise take. The command
  runs inside a repository, like every other `bench worktree` subcommand, and
  `Pool(root)` names the one key it may never target.
- **The fingerprint is over the plan, not over a clock.** It is derived with the
  package's existing `fingerprintParts` helper over the sorted reclaimable key
  names and each key's child `gitdir:` targets, so any change to what would be
  removed changes it and a change elsewhere in the pool does not.
- **Apply re-verifies per key.** The fingerprint proves the plan is current as a
  whole; the immediate re-check proves each individual key still qualifies at the
  instant of removal. Removal is `os.RemoveAll` on a path whose parent is
  asserted to be `benchHome()/worktrees` exactly.
- **No home-wide lock.** Bench has no lock over `$BENCH_HOME`, and inventing one
  is a separate capability (Out of scope). The fingerprint plus the per-key
  re-check is the concurrency posture, and it is honest about its window: a key
  created after the re-check and before the `RemoveAll` is a race this does not
  close, which is why the current repository's key — the only key a concurrent
  Bench process in this session could be minting — is excluded outright.
- **Resume reports, never removes.** `ConservativeCleanup` gains one count from
  the same predicate and `renderResumeSummary` gains one line naming
  `bench worktree reclaim`. It is the `OrphanCandidate` posture: the sweep
  reports, the explicit path-addressed command acts.
- **AXI disposition.** `bench worktree` currently declares
  `axiApprovedChildren("list")`. `reclaim` is a state-changing surface, so it
  stays outside the approved read-only query set like `clean` and `create`, and
  keeps `clean`'s documented output contract: TOON on stdout, a definitive
  zero-row empty state, structured refusals on stdout, exit 0 success and
  idempotent no-op / 1 unsatisfied intent / 2 usage. Promoting the bare plan into
  the approved query set is a separate registry change (Out of scope).

## Testing decisions

- The external behavior a good test exercises: with a pool holding the shapes
  side by side, `bench worktree reclaim` names exactly the reclaimable ones and
  removes nothing; the apply removes exactly those and leaves every other key
  byte-identical; and a pool changed between the two refuses.
- Seams: the predicate (unit seam, a synthetic pool under a bound temporary
  `BENCH_HOME`), the command entry point driven with captured stdout/stderr and
  its exit code (prior art `internal/worktree/clean_operand_test.go` and
  `list_actions_test.go`), and `renderResumeSummary` (prior art
  `orphan_render_test.go`).
- Every test binds its own `BENCH_HOME` under its `t.TempDir`. FT226's
  `TestMain` residue red is the standing oracle for a test that forgets, so this
  spec adds no isolation machinery of its own.
- Gate seam: the ordinary `test` phase. The verb is Go, so no shell or gate
  edit; `cmd/bench/main.go`'s registry and `internal/usage/worktree.go` carry the
  help rows the conformance checks read.
- The `/bench-debug` Phase 1 repro is retained as a test rather than a
  throwaway: a pool key whose source repository is deleted, driven through the
  real command, must be planned and then removed.
- Mutation probes recorded in the verification log, not retained: (a) delete the
  `IsNotExist` narrowing so any stat error counts as absent → the
  permission-denied case reds; (b) delete the current-key exclusion → the
  empty-current-key case reds; (c) accept any fingerprint → the drift case reds.

### Seam diagram

    operator: bench worktree reclaim [--apply <fingerprint>]
        │
        ▼
    $BENCH_HOME/worktrees  ──▶  [ poolKeyReclaimable — one predicate ]  ──▶  verdict per key
        │                            ◀ unit tests: synthetic pool, every hostile shape
        ├──▶ plan:  TOON rows + aggregate + apply invocation carrying the fingerprint
        ├──▶ apply: fingerprint match ──▶ per-key re-check ──▶ RemoveAll ──▶ TOON rows + counts
        └──▶ resume: count only, plus the reclaim command   ◀ render test
                                                            ◀ probe: predicate widened → hostile shape reds

### Acceptance coverage map
| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| PL1 | 1, 5 | a pool holding a dangling-pointer key, an empty key, and a live key plans exactly the first two as reclaimable, as TOON rows carrying key, verdict, and reason | command test over a synthetic pool | a plan that reports every key, or none, fails here; a plan that reports the live key is the destructive bug |
| PL2 | 2 | every retained key carries a reason naming what protected it, distinct per protecting shape | command test asserts a distinct reason for the live-pointer, repository-directory, no-`gitdir:`, stray-entry, and mixed keys | one shared reason string for every retention makes the output useless for deciding whether the retention was correct |
| PL3 | 3 | the plan prints the apply invocation including the fingerprint the apply will accept | command test feeds the printed invocation's fingerprint straight into an apply that succeeds | a plan printing a fingerprint the apply rejects makes the handshake unusable, which a test asserting only the string's presence would miss |
| PL4 | 4 | a pool with no reclaimable key, and an absent pool directory, both print a zero-row table and exit 0 | command test on an empty home and on a home whose `worktrees` directory does not exist | silence, an error, or a non-zero exit on a clean pool turns a successful absence into a failure |
| PL5 | 6 | after a bare plan over a pool holding every shape, the sorted recursive listing of the pool is byte-identical | command test compares the listing before and after | a plan that removes anything is the worst failure this command has |
| SH1 | 7, 11 | a key whose child `.git` pointer target exists, and a key holding one live and one dead pointer, are both retained | predicate unit test | a predicate that reclaims on any dangling child rather than every child destroys a live worktree |
| SH2 | 8, 9 | a key whose child `.git` is a directory, and one whose `.git` file has no `gitdir:` line, are retained | predicate unit test | reading a repository directory or an unparseable file as a dangling pointer reclaims a real repository |
| SH3 | 10 | a key holding a stray top-level file, or a top-level entry that is not a directory, is retained | predicate unit test | a predicate that only inspects child directories and ignores other entries deletes whatever the operator put there |
| SH4 | 12 | a key whose child cannot be read, and one whose `gitdir:` target cannot be stat'd for a reason other than absence, are retained and reported | predicate unit test with a mode-0 child directory | treating any error as absence reclaims on unknown, which is exactly the unsafe direction |
| SH5 | 13 | a symlinked key, a symlinked child, and a `.git` that is a symlink are each retained and never followed | predicate unit test | following a link makes bytes outside the pool the subject, so the pool stops bounding what can be removed |
| SH6 | 14 | the current repository's own key is retained even when it is empty | command test run inside a repository whose key exists and is empty | the empty-key clause otherwise takes the key of a session that has just acquired its pool directory |
| SH7 | 21 | every removed path's parent is exactly `$BENCH_HOME/worktrees` | apply-path unit test asserting the parent of each target | a target assembled from a key name alone can escape the pool if the name is ever attacker- or accident-shaped |
| AP1 | 15, 18 | `--apply` with the plan's fingerprint removes exactly the planned keys, leaves every other key present, and prints per-key verdicts with removed and retained counts | command test: plan, apply, compare listings | an apply that removes the retained set, or reports counts it did not act on, is undetectable from the plan alone |
| AP2 | 16 | an apply whose fingerprint does not match the current pool refuses on stdout naming the re-plan command and exits 1, removing nothing | command test mutates the pool between plan and apply | accepting a stale fingerprint removes keys on the strength of a reading that is no longer true |
| AP3 | 17 | a key that becomes non-reclaimable between plan and apply — a live `.git` pointer appears in it — survives the apply | `TestReclaimApplyRetainsAKeyThatWentLiveAfterThePlan` drives `applyPoolReclaim` over a captured plan, not the command: the fingerprint is derived over the reclaimable subset only, so a planned key going live necessarily leaves that subset and changes the fingerprint, and the stale refusal fires before the re-check is ever reached (build finding, reviewer-approved correction) | the window the re-check guards opens inside the command, after its own re-plan matched and while it removes keys one at a time; a probe deleting the re-check reds this test |
| AP4 | 19 | an apply against a pool with no reclaimable keys exits 0 having removed nothing | command test | a no-op that errors makes a repeat invocation unsafe to run |
| AP5 | 20 | an absent, empty, or malformed `--apply` value exits 2 with usage and removes nothing | command test over the three forms | a parser that treats a missing value as "apply everything" turns a fumbled flag into a destructive run |
| RS1 | 22 | `bench resume-clean` reports the reclaimable-key count and names `bench worktree reclaim` as the action | render test over a resume result carrying a non-zero count | a resume that stays silent leaves the operator with no path from the debris to the verb |
| RS3 | 24 | the count `bench resume-clean` reports equals the target count `bench worktree reclaim` plans over the same pool, across the hostile-shape pool and a clean one | command test comparing the two numbers | two derivations of the predicate drift, and the ambient number is the one an operator trusts without re-checking |
| RS2 | 23 | `bench resume-clean` over a pool holding reclaimable keys leaves the pool listing byte-identical | command test compares listings across a resume | an automatic removal at session start is destructive outside any tree the gate observes |
| RP1 | 1, 15 | the debug repro's exact shape — a key whose source repository was created, keyed, and deleted — is planned and then removed by the real command | retained regression test | the bug this spec exists for must have a test that reds on its return |

### Edge inventory

- **Error path:** unreadable key, unreadable child, unreadable `.git`, and a
  `gitdir:` target whose stat fails for any reason other than absence all retain
  and report (SH4).
- **Empty/absent input:** an absent `worktrees` directory and an empty one are
  both the zero-row answer (PL4); an empty *key* is reclaimable (PL1); an empty
  *child* directory, holding no `.git` at all, is not.
- **Boundary values:** a key with exactly one child, a key with two children of
  mixed liveness (SH1), and a key with zero children.
- **Malformed input:** a `.git` file with no `gitdir:` line, with a `gitdir:`
  line and no value, or with trailing content after it, all retain (SH2).
- **Special files:** a FIFO, device, or socket found where a key, a child, or a
  `.git` belongs is retained unread (SH3, SH5) — the profile's checklist class
  for special files in a discovered path.
- **Paths with spaces or glob characters:** the pool already holds
  `project with spaces-2317837872`; keys are handled as path values, never
  interpolated into a shell, and the TOON emitter refuses control bytes in a key
  name rather than rendering them.
- **Interrupted / partial state:** an apply killed partway leaves the keys it
  already removed removed and the rest untouched; the next plan simply reports a
  smaller target set. No receipt is written, because there is no multi-step
  transaction per key to resume — one `RemoveAll` is the whole operation.
- **Re-run idempotency:** a second apply with the same fingerprint refuses as
  stale (the pool changed) or, over an unchanged clean pool, exits 0 with nothing
  to do (AP2, AP4).
- **Won't handle:** a home-wide lock serialising concurrent `reclaim` runs — the
  fingerprint and the per-key re-check bound the window, and the only key a
  concurrent Bench process in this session mints is the current repository's,
  which is excluded outright (SH6). Two operators reclaiming the same pool
  simultaneously is not a shape this tool has.
- **Won't handle:** reclaiming a key whose source repository still exists but
  whose worktrees were removed by hand — that key is reachable from its root and
  is `bench worktree clean`'s subject, not this command's.
- **Won't handle:** recovering work from a key before removing it — every child
  reaching the predicate has a dangling pointer, so git cannot read it and there
  is no branch, ref, or index left to preserve. A key with anything recoverable
  fails the predicate and is retained.
- **Won't handle:** running outside a repository. Every `bench worktree`
  subcommand requires one, and the current repository's key is what the exclusion
  needs; a homeless invocation refuses with the not-in-repo error.
- **Won't handle:** a size or age threshold on what may be reclaimed — the
  predicate is about provable deadness, and an age gate would only delay a
  verdict that is already certain.

## Ownership fences

One writer at a time; the three tickets are serial. Reviewer disposition:
approve, merge, or split at sign-off.

- `internal/worktree/pool_reclaim.go`
- `internal/worktree/pool_reclaim_test.go`
- `internal/worktree/lifecycle.go`
- `internal/worktree/resume.go`
- `internal/worktree/worktree.go`
- `internal/worktree/orphan_render_test.go`
- `internal/usage/worktree.go`
- `cmd/bench/main.go`
- `cmd/bench/main_test.go` (the help-inventory golden paired with the row above)
- `projects/benchkit.md`
- `specs/ft234-pool-key-reclaim/`
- `capture/session-handoff.md`

## Out of scope

- **A home-wide lock over `$BENCH_HOME`** so concurrent pool operations
  serialise: ~4 edits (a lock helper, its acquisition in `reclaim` and
  `Acquire`, tests), 1 gate run. The fingerprint plus per-key re-check is the
  posture this spec ships.
- **Promoting the bare plan into the approved AXI query set** as
  `axiApprovedChildren("list", "reclaim")` with the contract fragments that
  implies: ~4 edits (registry, the AXI contract fragment, its test, the profile's
  AXI seam paragraph), 1 gate run. The verb is state-changing as a whole, so it
  takes `clean`'s documented contract for now.
- **Reclaiming a live repository's own dangling pool children** — the
  unregistered, unrecorded child inside a reachable key that the discriminating
  experiment also found surviving: ~5 edits (a plan source in `ConservativeCleanup`,
  its predicate reuse, rendering, tests), 2 gate runs. Reachable, so it belongs to
  the repository-scoped cleanup path, not the home-scoped one.
- **A one-off operator sweep of the eight reclaimable keys currently in the
  pool** — superseded: a validated throwaway script for exactly this population
  was written and deliberately left unapplied, because this spec ships the verb
  that does it and the reviewer runs that instead.

## Further notes

The operator's current pool is the acceptance demo, not a test fixture: after the
build lands, the reviewer runs `bench worktree reclaim` against
`$HOME/.bench/worktrees`, reads the plan, and decides whether to apply. No build
step writes there. The population is known in advance and held back for this
purpose — five keys pointing into deleted `/tmp/bench-dogfood-*` repositories,
three empty keys, and one key (`bench-2826441890`, this repo's own) that must
survive on three independent grounds: a stray top-level file and two children
carrying no `.git` pointer. A plan that names anything else, or misses one of
the eight, is the demo's red.
