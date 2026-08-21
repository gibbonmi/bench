# Retro — land-executable-freshness

## Outcome

`bench worktree land` now proves its own executable before it enforces any
landing contract. Landed at `84b7c4b0` from reviewed source `60d1154b` →
`62f92b2e`, spec `Status: implemented`, gate green.

The feature is small: a presence-only `freshness.DeclaresBuildInputs` beside the
manifest path it reads, a `verifyLandingExecutable` package-var seam,
`LandCommand` growing the invoked-executable parameter, and the registry closure
forwarding `Command.Executable`. All eight acceptance rows are graded, plus a
ninth behavior the re-review found ungraded. The board rewrite records that
FT242's original ask closed as FT225.

Folded in at the reviewer's explicit direction, not because it belonged to the
feature: the `regroup` example profile retired across all four registries that
named it, `COMPLIANCE_ASSESSMENT.md` removed, `ui_example/` untracked, and the
README's Regroup walkthrough replaced by a project-neutral design-system section
so removing a profile did not delete the only documentation of a shipped skill.

## Gate-stage timings

Measured on the landing's own gate run (`.logs/gate-20260821T013533…`):

| phase | elapsed |
| --- | --- |
| gofmt | 0.06 s |
| vet | 1.14 s |
| test | 54.63 s |
| race | 4.10 s |
| system | 3.09 s |
| shellcheck | 0.38 s |

Test remains the whole cost at 84% of the run; `internal/worktree` alone is
~49 s of it. Seven full gate runs were paid this landing — five ticket and
amendment commits, one capture commit on the destination, one landing.

## Ticket-versus-spec-slice and delegate performance

One build ticket, charged to Opus/medium as a full ticket-sized slice with the
coordinator's design decisions pre-settled in the charge. It returned all seven
files first-pass: production, tests, and the board rewrite, with five attributed
mutation probes and production restored exactly. It also self-reported one edit
beyond its charge's literal wording (the stale `ROADMAP.md` sequence entry) and
flagged it for revert rather than burying it.

The charge carried more design than usual — the seam placement, the predicate's
home, and a correction to the spec's own LF4 test seam. That front-loading is
what bought the first-pass result on a ticket whose spec had one unworkable row.

Review ran three Sonnet axes plus one scoped Sonnet re-review of the repair.

## Coordinator catches

- **The spec's LF4 seam could not work.** The map said to extend the existing
  resume journey "with an erroring check seam", but that test drives a real
  subprocess, so a parent-process package var reaches nothing. Caught at charge
  time by reading the test rather than the row. The charge replaced it with a
  real failing condition — the manifest committed into the destination between
  publication and resume.
- **The probe the delegate did not run.** Its five probes covered LF1, LF3,
  LF8, LF9, and LF10, but not the mutation LF4 exists to catch. Moving the check
  before the resume dispatch turned LF4 red, confirming the row is load-bearing.
- **Two ownership-fence omissions**, each caught by a red `bench preflight
  review` rather than by reading: the freshness owner, then its test file.
- **The re-review finding was real.** Six cases all produced a nil `Lstat`
  error, so collapsing the predicate to `err == nil` survived every one. An
  unusable parent is the state that separates "absent" from "errored at all",
  and it is the direction that must never skip the proof.
- **Delegate citation drift.** The Spec axis cited `cmd/bench/main.go:94` for a
  closure at line 555 and gave placement line numbers that matched no rendering
  of the file. Every substantive claim held; the line numbers did not.

## Repair attribution

| ticket | repair rounds | cause per round |
| --- | --- | --- |
| 01-refuse-stale-landing-executable | 0 | none |
| 02-pin-the-applicability-predicate | 1 | spec-row |

Ticket 02 was itself the repair for review's one accepted finding. Its own
acceptance rows named the symlink and directory states but not the error class
the predicate's discrimination actually turns on, so the first pass left the
`errors.Is` half ungraded — a missing row, not a delegate miss.

## Agent-experience improvements

### Bench CLI

- `bench handoff` overwrote a mid-phase `Next command`. The phase says to write
  the phase reached and then refresh the pin; in that order the verb replaced
  the phase-correct next action with its own board routing (`git push`).
  Expected effect: either the phase documents the reverse order, or the verb
  leaves an agent-authored next command alone while a phase is mid-run.
- `bench worktree release` refused with "ignored residuals require
  `--discard-ignored`", but `bench worktree release --help` does not list that
  flag. The remedy named is not in the advertised grammar. Expected effect: an
  operator can act on the refusal without guessing where the flag lives.
- Sibling verbs disagree on flag surface: `bench preflight review` accepts
  `--source-tip`, `bench diff` rejects it. Reviewing an explicit pair means
  passing the tip to one verb and trusting the checkout for the other.

### Skills

- `craft-spec`'s coverage-map discipline should require a row naming its test
  **venue** when the seam is a package-var substitution. LF4 named a real test
  and a real technique that cannot combine, and nothing in the map's shape made
  that visible — the row read as complete until someone opened the test.
- `craft-review`'s citation standard holds on substance but not on line numbers.
  A delegate citing a file and line it did not re-read produces a claim that
  looks checkable and is not. Worth one line: cite the line you read this pass,
  or cite the symbol instead.

### Process

- Ownership fences went red twice for the same reason — a file the design
  needed but the spec had not listed, then its test file. Deriving the fence set
  from the tickets' `Writes:` lines at spec time would have caught both before
  any gate run.
- A reviewer-directed housekeeping batch folded into a feature landing needs its
  own fence block in the spec, and it makes the landing's diff two unrelated
  stories. It worked, but the spec now carries a section that is not about the
  feature. Worth a convention: a housekeeping change lands on its own, unless
  the reviewer chooses otherwise and the spec says why.
- Untracking a file that stays on disk interacts with the landing contract: the
  now-ignored path became ignored residue in the source worktree and stopped the
  release step after publication. The landing was complete and correct; the
  recovery was to clear the residue and resume. Worth knowing before, not after.
