# Retro: bench-probe

## Outcome

The landing `d4085c6c` published the spec on 2026-09-06 over the source pair `963aa880`
to `d35641eb`. `bench probe` is a root verb. It preserves one subject under the Bench
home and applies one exact swap or omission. It runs one focused test or named check
through `testreport.Prepare` and `testreport.Execute`. Then it restores the subject and
proves the restore byte-exact.

It prints `bit` at 0, `silent` at 1, `invalid` at 1, or `restore-failed` at 2 with the
preserved copy named. `gate.ExecutionInProgress` refuses a probe under a live gate run.
`bench worktree path` notes the exec route on stderr. The reference guide and
`CONTEXT.md` state the verb and the **probe verdict** term. The spec stays `implemented`
as the veto surface. `reviews/bench-probe.md` holds four `ask-user` findings for the
reviewer.

## Gate-stage timings

| stage | landing gate |
| --- | --- |
| gofmt | 89 ms |
| vet | 936 ms |
| test | 67288 ms |
| race | 2392 ms |
| system | 23570 ms |
| shellcheck | 521 ms |

Four fold gates and one whole-tree gate ran green before the landing. The test stage ran
between 64860 ms and 68198 ms on each.

## Ticket-versus-spec-slice and delegate performance

Six ticket charges and one repair charge ran on `opus`, one per ticket. Tickets 01, 02,
and 03 ran in parallel in three sibling worktrees. Tickets 04, 05, and 06 ran after the
folds, with 06 in a fourth sibling. Every charge landed first-pass on behavior, and every
self-probe bit.

Two charges returned to their delegate once for the growth lane. The lane grades growth
against the current tip, so the spec base gave no headroom at the commit. The verb charge
reported the PB4 shortfall with its cause and amended the row under the batch approval.
Three review axes ran on `opus`. The Coverage axis at medium found the one behavior
defect: a render refusal over a failed restore lost the `restore-failed` verdict. One
repair-scoped re-review at low found one undecided edge, which the spec now lists as a
Won't handle.

## Coordinator catches

- Ticket 02 edited `internal/testreport/check_test.go` outside its fence, because that
  test calls the changed `runGoTest`. The fence took the file as an amendment.
- The named check `bench test --check subcommand-routing` compiles its table from the
  run binary's source root. A probe of a moved test table needs the worktree's own
  binary or the root override.
- A probe that omits a line and leaves an unused variable does not compile. The verb
  reported it as `invalid`, and the coordinator replaced it with a swap.
- The lane's growth base is the current tip, not the spec base. A registry line in an
  over-budget file needs its headroom in the same commit.
- The preflight's registry closure requires the five bound registry files on a repair
  ticket that writes `internal/probe`.

## Repair attribution

| ticket | repair rounds | cause per round |
| --- | --- | --- |
| 01-make-headroom-in-the-registry-files | 0 | none |
| 02-expose-the-focused-run-outcome | 1 | tree-drift |
| 03-export-the-gate-execution-probe | 0 | none |
| 04-add-the-bench-probe-verb | 2 | tree-drift, spec-row |
| 05-note-the-file-tool-route-on-the-path-verb | 0 | none |
| 06-state-the-probe-verb-in-the-guidance | 0 | none |
| 07-repair-the-verdict-override-rows | 0 | none |

The two `tree-drift` rounds are the growth lane's tip base. The `spec-row` round is the
render refusal over a failed restore, which row PB47 now owns.

## Agent-experience improvements

### Bench CLI

- Let `bench test` print a `tests_run` count in its packages row, so a green package is
  distinguishable from a no-test-run without a second `-v` call. The census entry
  `bench-probe census 5` records the five raw shell reads.
  Feeds: new
- Let the lane's growth refusal name its base commit, because the spec's growth check
  and the lane's growth check take different bases.
  Feeds: new
- Let `bench test --checks` list the check inventory at exit 0, because the only route
  today is an unknown check at exit 2.
  Feeds: new

### Skills

- Let `bench-craft-tickets` state that a ticket which adds a line to an over-budget file
  moves its headroom in the same ticket. The lane grades growth per commit.
  Feeds: none
- Let `bench-craft-delegate` name `bench probe` as the coordinator's probe form, so the
  copy-aside sequence leaves the charges.
  Feeds: none

### Process

- A coordinator probe runs through `bench probe` from the first fold of the verb. A
  probe the verb reports as `invalid` is replaced before a verdict is read.
  Feeds: none
- A material acceptance shortfall under a batch approval amends the row and records the
  reason in the spec before the review round. The Spec axis then reads one source.
  Feeds: none
