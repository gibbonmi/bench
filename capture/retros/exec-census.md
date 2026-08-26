# Retro: exec-census

## Outcome

The spec landed as `415c99e9` on 2026-08-26 from the reviewed pair
`3d76ff1b..ace105ca`, with the whole-project gate green in the landing run.
Seven tickets, one fence amendment, two review repairs, and one anchor reflow
committed serially on one integration source. `bench guard-bench-follow-on`
records one raw call per Bash call that names a pool path without a `bench`
verb. `bench status` shows the `census` row, the landing prints `census=<n>`,
retirement drops the records, and the final check owns the learning duty.

This landing's record carries no `census=` key, because the installed binary
predates the build, and no census directory exists on this machine. The close
states `census: 0 raw calls`. The census starts to count after `bench repair`
or the release install publishes the new binary.

## Gate-stage timings

| phase | elapsed |
| --- | --- |
| gofmt | 76 ms |
| vet | 1.5 s |
| test | 43.0 s |
| race | 4.7 s |
| system | 9.7 s |
| shellcheck | 0.4 s |

Two landing attempts went red before the green one. The first red was a
canary fixture anchor that ticket 07's reflow had joined. The second red was
one `infrastructure` refusal inside a nested landing test under full parallel
load; it did not reproduce in four later runs.

## Ticket-versus-spec-slice and delegate performance

Every charge was ticket-sized; no charge received a spec slice. Opus/high
carried the signal, cleanup-authority, and guidance tickets and each landed
first-pass on behavior. Opus/medium carried the two prefactor tickets and
each needed one one-source repair round. Sonnet/low carried the two exact-spec
tickets and the coverage pin; each landed first-pass on behavior, and two of
them needed a one-source repair round. Sonnet/high ran the three review axes
and the repair-scoped re-review; each returned one cited finding.

## Coordinator catches

- Four knowledge duplications at fence boundaries: the pool join, the assignment segment shape, the `git` global-option table, and the pool parent join.
- One canary fixture red that a delegate called pre-existing; it belonged to ticket 01's file move.
- One invalid coordinator probe: a print to `os.Stderr` that the test's writer never captured; the probe was redone through the verb's own writer.
- One stray `bench` binary built into the worktree root.
- Every coordinator probe bit at a site and kind distinct from the delegate's.

## Repair attribution

| ticket | repair rounds | cause per round |
| --- | --- | --- |
| 01-derive-the-census-home-beside-the-pool | 1 | ticket-slicing |
| 02-record-one-raw-call-per-bash-call | 1 | delegate-error |
| 03-name-the-real-verb-head | 1 | delegate-error |
| 04-record-through-the-follow-on-guard-verb | 1 | delegate-error |
| 05-show-the-census-row-on-the-board | 0 | none |
| 06-carry-and-drop-the-census-across-a-landing | 0 | none |
| 07-add-the-census-duty-and-the-charge-line | 1 | ticket-slicing |
| 08-anchor-the-learning-fields-and-the-retro-citation | 0 | none |
| 09-pin-the-first-assignment-id-rule | 0 | none |

Ticket 01's round is the canary fixture that anchored on the moved `cksum`
skip. Ticket 07's round is the canary fixture that anchored on the reflowed
self-probe line. Neither ticket named the fixture registry.

## Agent-experience improvements

### Bench CLI

- Name the token that trips the follow-on guard in its refusal text; three delegates guessed at the cause of a blocked call.
  Feeds: new
- Treat a heredoc body with no `bench` word as not a follow-on, so a scripted edit can run through `bench worktree exec`.
  Feeds: new
- Add a `bench test --package <path> --run <pattern>` projection that reports only the failure set.
  Feeds: new
- Add a `bench gate --check <name>` form, so one conformance check runs without a hand-built environment.
  Feeds: new
- Publish the census binary with `bench repair` after a landing that changes the promotion broker, so the next landed record carries `census=<n>`.
  Feeds: none

### Skills

- Make `craft-tickets` require a ticket's `Writes:` to name every canary fixture that anchors on a moved or reflowed line.
  Feeds: new

### Process

- Run the fixture-bite conformance test in a delegate's focused checks when a ticket moves a test or reflows anchored prose.
  Feeds: new
- Probe through the writer the test captures. A print to `os.Stderr` proves nothing against a captured `stderr` parameter.
  Feeds: none
- The coordinator and two delegates ran raw `cd` calls into the pool path during this build. The census measures that habit from the next landing on.
  Feeds: none
