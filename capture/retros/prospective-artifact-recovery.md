# Retro — prospective artifact recovery

## Outcome

The landing published `0ec709aa` from source tip `c9401550` onto base
`c2b767d8`. The gate ran green across all six phases. The census recorded 136
raw calls in the integration assignment.

One prospective artifact owner now holds the private checkout and every
owner-authored run binary under one temporary-root prefix. The owner writes a
strict private record before it creates the checkout. Every prospective
producer sweeps its repository's dead bundles before it creates its own. Only
the absent-process probe result proves death; every other probe result retains
the bundle.

The review found two real defects in the first build. An authored run binary
escaped the bundle when the candidate tree had no kit of its own. A temporary
root with a symbolic-link component left the Git registration forever, because
Git registers the resolved path. Both are repaired and covered.

## Gate-stage timings

| stage | verdict | elapsed |
| --- | --- | --- |
| gofmt | green | 92 ms |
| vet | green | 735 ms |
| test | green | 43881 ms |
| race | green | 2224 ms |
| system | green | 13707 ms |
| shellcheck | green | 429 ms |

The first landing attempt was interrupted by an external signal after five
seconds. The second attempt went red on the branch-native census. The third
attempt landed.

## Ticket-versus-spec-slice and delegate performance

Every charge was ticket-sized. No charge received a spec slice, so this
landing offers no contrast between the two shapes.

Ticket 01 landed in a prior Codex session on Terra. Its code carried the two
real defects and the census violation that this session repaired. Its own
tests, probes, and lane were green, so the gap was in what the charge and the
coordinator did not ask.

Four Opus charges ran in this session: ticket 02, the two repair tickets as
one serial charge, and one fix pass. Every one landed first-pass on behavior
with a biting self-probe. Ticket 02 duplicated the owner-record shape across
three test files, which the review caught as a Standards finding. Ticket 02
also reported a spec contradiction instead of working around it: the gate lock
refuses a second authorization, so the concurrent seam is unreachable.

Six Opus review axes ran, plus three scoped re-review axes. The initial axes
returned 15 raw findings in 4 repair targets. The Coverage axis confirmed the
symlinked-root defect with a real repository run before it reported. The
scoped re-review verified every predicate with an enumeration and returned no
blocking finding.

## Coordinator catches

I verified every done-claim against the tree and re-ran the focused checks
myself. Three catches were material.

The ticket 02 delegate claimed no PAR28 producer row existed in the tree. The
system journey already graded the full-gate record; I refuted the claim before
it became a repair ticket.

The system suite needs `BENCH_KIT` and `BENCH_RUN_BINARY`, which the gate
wrapper exports and a hand run does not. My first two hand runs reported
environmental reds that were not diff-owned. I attributed them before I
believed them.

My probes bit at four distinct sites: an omitted root guard, an inverted
death guard, an inserted direct process call, and a dropped bundle root.

The honest entry on the other side: I charged two gate-package test files
without naming the branch-native census. Focused tests were green, and the
whole-project gate went red at the landing. The rule was in the tree, and I
did not read it.

## Repair attribution

| ticket | rounds | cause per round |
| --- | --- | --- |
| 01-recover-dead-prospective-bundles | 2 | delegate-error, other |
| 02-protect-shared-prospective-owners | 1 | delegate-error |
| 03-confine-run-binaries-and-resolved-bases | 0 | none |
| 04-single-source-test-knowledge | 0 | none |

Ticket 01's first round repaired the escaped run binary and the symlinked
root. Its second round routed its tests through the sanctioned Git seam after
the census red; that cause is a charge omission, not the delegate's. Ticket
02's round collapsed the duplicated record shape and narrated comments.

## Agent-experience improvements

### Bench CLI

- Make the landing print the census verb-head breakdown before it drops the
  assignment's records. The exit duty asks for per-head counts that the landing
  has already deleted.
  Feeds: new
- Make `bench handoff` rewrite the body of `capture/session-handoff.md` from
  the tree state, or refuse a stale body. This run refreshed the pin and kept a
  body that described a paused build.
  Feeds: new
- Make `bench preflight build` read the retained source's frozen base by
  default. The bare form reads main's tip and reports a false red on every
  retained source.
  Feeds: new

The census entry for this landing is
`prospective-artifact-recovery: 136 raw calls in the integration assignment`
in `capture/learnings.md`. It records that the per-head counts were dropped.
Feeds: new

### Skills

- Teach `craft-delegate` that a charge for a test under an architecture-owned
  package names the branch-native census. The focused checks run that census.
  Feeds: new
- Teach `craft-spec` that a seam named "concurrent authorization" must first
  check the gate execution lock, because two gate runs in one repository cannot
  coexist.
  Feeds: new

### Process

- Run the system suite by hand only with `BENCH_KIT` and `BENCH_RUN_BINARY`
  set from a fresh worktree build, and attribute a red before believing it.
  Feeds: none
- Ask each review axis to refute a strong finding with a real run before it
  reports. One repository probe made the symlinked-root defect undeniable.
  Feeds: none
