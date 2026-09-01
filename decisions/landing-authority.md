# Landing authority (FT169)

Status: ready

## Destination

The reviewer decides the landing command's authority: its relation to
`bench commit`, `bench worktree release`, and `bench spec retire`. The
reviewer also decides the hook pool's one-active limit, the board-conflict
and spec-staging behavior, the primary-checkout landing route, the source-tip
interface, and the `--spec` normalization. The map also settles the refusal
standard: a refusal names the failed path, the recovery state, and the exact
next command. The map ends in one reviewed decision source for
`/bench-write-spec`. FT249 owns the capture-storage decision.

## #1: What is the current authority split?

Blocked by: none
Type: Research

### Question

Which state does each of `bench worktree land`, `bench commit`,
`bench worktree release`, and `bench spec retire` mutate today, and which
refusals does each enforce? Where does the one-active limit live?

### Answer

Resolved 2026-09-01. See
`decisions/assets/landing-authority-authority.md`. The landing is the only
verb that writes `main`, and the only Bench verb that runs on the primary
checkout. `bench commit` writes one branch tip and refuses the primary
checkout. `bench spec retire` deletes files and commits nothing. No
one-active pool cap exists in code; a request digest, a create flock, the
session-derived request id, and a destination compare-and-swap give the
effect.

## #2: Which recorded refusal faces are still open?

Blocked by: none
Type: Research

### Question

For each FT169 occurrence class, does the current refusal name the failed
path, the recovery state, and the next command? Classify each as fixed,
partial, or open.

### Answer

Resolved 2026-09-01. See
`decisions/assets/landing-authority-refusals.md`. Three occurrences are
fixed by cause removal, and eight are partial. Two are open: the
conflict-versus-pending-merge distinction, and the `--spec` path form on
the staged branch. The gap is structural: the one refusal formatter makes
`next` optional, and 4 of 29 sites set it. One adjacent correctness
fault: `spec.Resolve` reads an explicit path against the process working
directory, not the passed base.

## #3: Which input forms does the landing parse today?

Blocked by: none
Type: Research

### Question

Which source-tip forms parse, how does `--spec` resolve a slug and a path,
and where does `--base` validation run relative to the review?

### Answer

Resolved 2026-09-01. See
`decisions/assets/landing-authority-parse-seams.md`. A 4-to-39-character
hex prefix already expands to the worktree head; a symbolic tip refuses
with a mismatch detail, not with the accepted form. The `--spec` slug
resolves on both landing branches; the path form works on the tickets-only
branch and breaks on the staged branch. The landing validates `--base`
against the recorded start and the destination; `bench preflight review`
validates neither, so a review can pass on a base the landing refuses.

## #4: How do the request token and the recovery surfaces behave?

Blocked by: none
Type: Research

### Question

Where does the request token live, what does a lost token refuse, when is
the census record deleted, and what continues an incomplete publication?

### Answer

Resolved 2026-09-01. See
`decisions/assets/landing-authority-recovery.md`. The intent ledger stores
the token digest and the plain token; a lost token refuses with a named
`bench worktree reauthorize` or `bench worktree list` next step. The
release step deletes the census record, so the record survives an
interrupted landing. A post-publication failure prints an exact `--resume`
continuation. A textual conflict refuses before publication, and nothing
distinguishes a committed conflict resolution from a pending
`git merge --continue` (FT258 holds that contract).

## #5: When does `--base` validation run?

Blocked by: none
Type: Grill

### Question

What validates the landing's `--base`, and when?

### Answer

The landing's `--base` is the assignment's recorded start.
`bench preflight review` refuses a `--base` that is not an ancestor of the
destination, before the review runs. The review then runs on the base the
landing accepts. A refusal after a complete review is a defect. Source: the
FT169 row, recorded by the reviewer.

## #6: What is the landing command's authority?

Blocked by: #1, #4
Type: Grill

### Question

What is the decided relation of `bench worktree land` to `bench commit` and
`bench worktree release`? What is the hook pool's one-active limit?

### Answer

Decided 2026-09-01. `bench worktree land` is the only verb that writes
`main`, and the only Bench verb that runs on the primary checkout.
`bench commit` writes one branch tip and never merges. `bench worktree
release` stays the landing's cleanup step. The one-active limit is the
four existing mechanisms — the request digest replay, the create flock,
the session-derived request id, and the destination compare-and-swap. No
new counter or pool cap is added.

## #7: Who resolves board conflicts and stages specs?

Blocked by: #1
Type: Grill

### Question

How does the landing handle a board conflict and spec staging, and what is
the primary-checkout landing route?

### Answer

Decided 2026-09-01. A staged spec lands as ordinary content through a
spec-less landing. `--spec` is reserved for the built spec's flip. When
`--spec` names a spec whose status is still staged, the landing refuses,
and the refusal's next names the spec-less form. The board stays outside
the landing's writes; the composition rules (handoff source-wins, capture
union) stand. The primary checkout is the decided landing venue: the
exec-route refusal states that the landing runs outside `bench worktree
exec` and names the exact command.

## #8: What is the source-tip and `--spec` interface?

Blocked by: #3
Type: Grill

### Question

Does the interface accept a symbolic tip and resolve a named worktree tip?
Does it accept or refuse an unambiguous short source tip? What is the one
normalization for the `--spec` slug and path forms?

### Answer

Decided 2026-09-01. `--source-tip` accepts three forms: a full sha, an
unambiguous hex prefix of 4 to 39 characters, and a symbolic tip such as
`HEAD`. The landing resolves the symbolic tip in the source worktree. A
mismatch refusal names the accepted forms. Both `--spec` branches route through one normalization
seam (`spec.LiveSpecSlug`), so the slug and path forms resolve
identically on the tickets-only and staged branches. This one seam also
removes the CWD-relative read in `spec.Resolve` for the landing path.

## #9: What is the refusal and recovery standard?

Blocked by: #2, #4
Type: Grill

### Question

What must every landing refusal name? How does the record distinguish a
committed conflict resolution from a pending `git merge --continue`? What
is the recovery route for a lost request token, a stale executable, and an
undeclared ignored residue?

### Answer

Decided 2026-09-01. The recovery route becomes mandatory at the refusal
seam: every refusal carries the failed paths and an exact next command,
and the three formatter families align on that shape. The landing's
conflict recovery reads `MERGE_HEAD` itself: a pending merge names
`git merge --continue`, and a committed resolution names the
commit-and-review route. This does not wait for FT258, which keeps the
`bench commit` contract. The lost-token, stale-seal, and residue routes
keep their current next commands; the residue refusal gains a
declare-or-clean next.

## #10: One spec or a split?

Blocked by: #6, #7, #8, #9
Type: Grill

### Question

Do the settled decisions ship as one spec, or do the refusal-standard fixes
and the interface changes split into two independently useful specs?

### Answer

Decided 2026-09-01. Two specs, in order. The refusal-standard spec lands
first: the mandatory-next formatter seam, the reworded partial refusals,
and the `MERGE_HEAD` distinction. The interface spec lands second and
adds its new refusals through that seam. It carries the symbolic tip,
the one `--spec` normalization seam, the staged-spec refusal, and the
preflight base checks.

## #11: Which base validations does preflight mirror?

Blocked by: #5
Type: Grill

### Question

Does `bench preflight review` mirror one landing base check or all three?

### Answer

Decided 2026-09-01. Preflight review mirrors all three landing
validations from one shared validation source. The base resolves as an
ancestor of the source tip. The base is the assignment's recorded start
or a descendant of it. The base is an ancestor of the destination. A
green review can then never precede a base refusal at land time.

## #12: How far does the refusal standard reach?

Blocked by: #9
Type: Grill

### Question

Does the mandatory-next standard bind the landing path only, the
occurrence-named verbs, or every lifecycle verb?

### Answer

Decided 2026-09-01. The standard binds every refusal face an FT169
occurrence names: the landing path and its resume, the wrapper's exec
route, and `bench spec retire`'s primary-checkout refusal. Other verbs
adopt the standard only when their code is touched.

## Not yet specified

## Spec-writer discretion

- Choose the exact refusal message text, if the message names the failed
  path, the recovery state, and the exact next command.

## Out of scope

- The capture-storage decision. FT249 owns it. FT169 only requires the
  landing to state its capture-related refusal and recovery route after
  that decision.

## Sources

- Path: `decisions/assets/landing-authority-authority.md`
  Supports: #6 and #7 — the authority split and the four one-active mechanisms. Produced 2026-09-01 by a read-only research delegate, with coordinator spot checks.
  Drift: re-read if `internal/worktree/land.go`, `internal/commit/commit.go`, or `internal/landing/landing.go` changes before the spec reads this map.
- Path: `decisions/assets/landing-authority-refusals.md`
  Supports: #9 and #12 — the per-occurrence refusal grades and the optional-next formatter gap. Produced 2026-09-01 by a read-only research delegate, with coordinator spot checks.
  Drift: re-grade after any landing-path refusal wording changes.
- Path: `decisions/assets/landing-authority-parse-seams.md`
  Supports: #5, #8, and #11 — the source-tip, `--spec`, and `--base` seams. Produced 2026-09-01 by a read-only research delegate, with coordinator spot checks.
  Drift: re-read if `internal/worktree/land_identity.go`, `internal/spec/spec.go`, or `internal/preflight` changes before the spec reads this map.
- Path: `decisions/assets/landing-authority-recovery.md`
  Supports: #6 and #9 — the token, census, seal, residue, and resume lifecycle. Produced 2026-09-01 by a read-only research delegate, with coordinator spot checks.
  Drift: re-read if the intent ledger, the resume path, or `.bench/build-outputs.json` changes before the spec reads this map.
