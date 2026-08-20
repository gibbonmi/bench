# Retro — ft230-release-through-bench

## Outcome

Landed and retired. The publish job of `.github/workflows/release.yml` now
runs one `dist/bench release submit --adapter npm --provenance` invocation
over downloaded preflight evidence; `bench release submit/promote/rollback`
select their adapter from `--adapter npm|fixture` (default `fixture`);
`checkReleaseWorkflow` flipped with five red-capable bite subtests; two
step-name byte contracts retired to a record-level ordering test; the runbook
records the tag-push presence rule. Landing commit `b48c4609` over reviewed
pair base `1766e4d1` / tip `736fb3a6`; retirement commit follows it. All 13
coverage rows landed; stories 7, 12, and 19 are review-graded as the spec
declares.

## Gate-stage timings

A full gate run held near 3 minutes wall throughout; the test phase dominates
(the `internal/worktree` package alone runs 48–56s), race adds ~3s of shown
tests, and gofmt/vet/system/shellcheck stay in single-digit seconds. The
build paid 10 full gate runs (5 on the integration source, 5 on the
destination — one red) plus 2 reused fresh verdicts on the ticket commits.
Three of the five destination runs graded markdown-only diffs (handoff pin,
spec mirror, capture), which motivates the capture-lane idea in
`capture/learnings.md`.

## Ticket-versus-spec-slice and delegate performance

Three write charges, all ticket-sized, all first-pass at the diff level.
T1 (opus/high) delivered rows R1–R7 and R13 with five attributed mutation
probes and two good unforced judgment calls (reviving the dead
`canary_shared_test.go` harness as the single approved-set fixture source;
collapsing the thrice-copied registry-base fallback). T2 (opus/medium)
delivered the workflow swap, the conformance flip with the job-scoping cheat
tests, both retirements, and the runbook amendment, and self-reported its two
out-of-fence mechanical edits rather than hiding them. The repair charge
(opus/medium) read the state machine before asserting and corrected the
ticket's "completes" sketch to the real `in_progress`/`promote` terminal
state instead of bending production. Review: three sonnet axes returned
0/0/3 findings; the coverage axis alone produced both accepted findings, and
the scoped sonnet re-review verified the repair empirically.

## Coordinator catches

No delegate done-claim was false this build; every verification (focused
tests, gate, diff inspection) confirmed the claims. The coordinator-side
catches were process reds, not delegate errors: the review preflight caught
T2's two out-of-fence paths (spec fence omission, amended and flagged for
veto); the land refusal caught the destination's divergent spec bytes; the
retirement gate caught the orphaned `roadmap/FT230.md` that
`bench spec retire`'s own next-step text does not mention.

## Repair attribution

| ticket | repair rounds | cause per round |
|---|---|---|
| wire-adapter-selection | 1 | spec-row |
| swap-workflow-and-flip-conformance | 1 | spec-row |
| repair-npm-failure-paths | none | none |

T1's round is the post-review failure-path repair: the spec's edge inventory
decided the absent-binary and mid-sequence behaviors but mapped no coverage
row to them. T2's round is the ownership-fence amendment: two gate-forced
paths (`tier_test.go` registration, canary anchor uniqueness) the fence did
not anticipate.

## Agent-experience improvements

### Bench CLI

- `bench worktree land` refuses on divergent staged spec bytes while the
  review skill and the scorecard's decision line both say the landing
  publishes the source's spec bytes; the workaround cost a hand mirror
  commit and two full destination gate runs. Proposed shapes are in
  `capture/learnings.md` (2026-08-20 entry).
- `bench spec retire`'s next-step names the ROADMAP row but not the
  `roadmap/FT<n>.md` detail file; the orphan check then reds a full gate run
  later. The next-step should name both deletions.
- `bench handoff` replaces the file's Next-command section with the board's
  leading signal (`git push`), which erases the phase-resume invocation the
  handoff contract asks for mid-build.

### Skills

- `bench-review-implementation`'s amendment sentence ("an amendment never
  routes through a hand commit on the destination") contradicts the land
  contract as shipped; one of the two must change with the CLI decision
  above.

### Process

- Batch capture files (handoff, retro, learnings, scorecards) into one
  gate-priced commit at phase close; this build initially paid a full gate
  run for a handoff pin alone.
