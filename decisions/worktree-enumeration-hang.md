# Worktree enumeration hang

Status: ready

## Destination

FT189: no malformed worktree admin entry can wedge a Bench command. Every
command that discovers worktrees fails fast with an attributable structured
refusal instead of inheriting the upstream `git worktree list --porcelain`
hang, and the mitigation names the upstream behavior it works around so it
retires if git bounds its own admin reads.

## #1: What exactly hangs, and where does Bench inherit it?

Blocked by: none
Type: Research

### Question

Reproduce the FT189 hang on this tree, bound the set of admin entries and git
commands affected, and map the Bench call graph that inherits it.

### Answer

Reproduced 2026-08-14, git 2.43.0: any of `gitdir`, `HEAD`, `commondir`,
`locked` under `.git/worktrees/<id>/` as a FIFO hangs `git worktree list
--porcelain` (blocking open-for-read with no writer). With a FIFO `gitdir`,
`git worktree add`/`lock`/`unlock`/`prune` and `git branch --list` also hang;
`git status`, `git rev-parse` (including `--git-common-dir`), and
`git for-each-ref` stay clean. Bench enumeration has one owner —
`git.Worktrees` in `internal/git/git.go`, seven production callers, reached by
`bench status`, the session-start `bench resume`, `bench worktree *`, and the
dashboard — and runs with no deadline today. The established bound seam is the
`internal/bounds` policy registry plus `bounds.Run`. Full matrix and repro
script: `decisions/assets/worktree-enumeration-hang-probe.md`.

## #2: Which mitigation posture does Bench own — pre-scan refusal, execution bound, or both?

Blocked by: #1
Type: Grill

### Question

The roadmap row offers a pre-scan that refuses a malformed admin entry by
shape, a bound on the enumeration call, or both. The pre-scan gives an
attributable refusal naming the exact entry and costs no deadline, but git's
admin read set is version-dependent and a scan can itself block on a
pathological filesystem; the bound catches every blocking shape including
unknown ones, but attributes only "timed out", and needs a deadline value.
Recommendation: both — the pre-scan (every entry under
`<git-common-dir>/worktrees/` must be a regular file or directory) for
attribution, the bound as backstop, each composing an existing seam: the scan
ahead of `git.Worktrees`, the bound as a named `internal/bounds` registry
constant driven through `bounds.Run`.

### Answer

Both (reviewer, 2026-08-14). Before enumeration, every entry under
`<git-common-dir>/worktrees/` must be a regular file or directory or Bench
refuses, naming the entry; the enumeration call itself runs under a named
`internal/bounds` registry deadline through `bounds.Run` as the backstop for
blocking shapes the scan cannot classify.

## #3: Does the mitigation cover enumeration only, or every worktree-admin git call?

Blocked by: #2
Type: Grill

### Question

`git worktree add`/`lock`/`unlock`/`prune` hang on the same malformed entries
and are invoked from `internal/worktree/{ownership,lifecycle,reauthorize,snapshot}.go`
and `internal/gate/engine.go`. Guarding only `git.Worktrees` closes FT189 as
written and covers the discovery path that runs first in practically every
session, but a mutation call that runs without prior discovery can still wedge.
Recommendation: enumeration owner only in this build — smallest diff that
removes the ambient hang — with the mutation-site exposure parked via
`bench idea` for a reviewed drain rather than widening this scope silently.

### Answer

Enumeration only (reviewer, 2026-08-14). The pre-scan and bound attach to the
sole enumeration owner `git.Worktrees` and nothing else; worktree-mutating git
calls stay unguarded in this build, and that exposure is parked in
`capture/IDEAS.md` for a reviewed drain.

## #4: What does the refusal disclose, and is there a repair route?

Blocked by: #2
Type: Grill

### Question

On a malformed admin entry the enumeration cannot partially succeed — git owns
the listing — so the refusal is fail-closed. What does it name, and does any
Bench surface offer repair? Recommendation: the structured refusal names the
offending entry path and its observed shape with "inspect and remove it" as
the next action, and `bench doctor` reports the same finding; no Bench command
deletes anything under `.git/worktrees/` — the entry is git-owned state and
removal stays a human act.

### Answer

Name + doctor, no delete (reviewer, 2026-08-14). The refusal names the
offending entry path and its observed shape with "inspect and remove it" as
the next action; `bench doctor` reports the same finding; no Bench command
deletes anything under `.git/worktrees/`. A plan/apply repair command was
rejected.

## Not yet specified

## Spec-writer discretion

- The bound's constant name and value, placed in the `internal/bounds` registry
  per its convention — never a call-site literal.
- Exact refusal wording and TOON shape, consistent with existing structured
  errors.
- Test fixture design; the FIFO fixture from the probe asset is the natural
  red-capable seed.
- How the retirement trigger is recorded: the mitigation's owner names the
  observed git 2.43.0 blocking-open behavior and retires when git bounds its
  own admin reads.

## Out of scope

- Fixing or patching upstream git; no upstream report is part of this work.
- Guarding worktree-mutating git calls (`add`/`lock`/`unlock`/`prune`) — #3
  scoped this build to enumeration; the exposure is parked in
  `capture/IDEAS.md`.
- Any Bench command that deletes or repairs entries under `.git/worktrees/` —
  a plan/apply repair route was rejected in #4.
- `git branch --list` exposure — Bench never invokes it (`LocalBranches` uses
  `for-each-ref`, which stays clean).
- Repairing `decisions/spec-build-review-gate-cadence.md`; that invalid map is
  its own shaping resume.

## Sources

- Path: `decisions/assets/worktree-enumeration-hang-probe.md`
  Supports: #1's full matrix and repro script, and the exposure facts behind #2 and #3's recommendations.
  Drift: git-version- and host-dependent; re-run its script after a git upgrade before trusting the matrix.
- Path: `ROADMAP.md`
  Supports: the FT189 row with the original 2026-08-03 repro and the mitigation framing this map decides.
  Drift: a working prioritization document; the row retires when this ships.
- Path: `internal/git/git.go`
  Supports: `Worktrees` as the sole enumeration owner every discovery caller shares, the seam #3's recommendation guards.
  Drift: re-read if enumeration ownership moves.
- Path: `internal/bounds/bounds.go`
  Supports: the named-constant policy registry and bounded process runner the #2 bound arm composes.
  Drift: none expected while the registry stands.
