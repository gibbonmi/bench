# Retro: lifecycle-refusal-names-component

## Outcome

The spec landed as `16bb950f` (tree `2b3cabe9`) from source pair `45626d91..bd40acb3`. It closes FT224. Every lifecycle identity refusal now names one component from one registry. `worktree exec` and `worktree path` print the resolver's reason through one printer. All 19 acceptance rows are green, and two package tests fence the property: a registry walk and a source scan.

Three tickets landed on one retained source: the registry and the landing sites, the resolver reason, and the review repair. The review found 6 raw findings and 5 repair targets; the repair-scoped re-review was clean.

## Gate-stage timings

The landing's gate run (`gate-20260825T092916`): gofmt 0.08 s, vet 1.4 s, test 79.8 s, race 4.9 s, system 10.6 s, shellcheck 0.4 s. The test phase is 82% of the run. The build paid ten full gate runs; every mutation probe ran a focused `go test` and cost no gate.

## Ticket-versus-spec-slice and delegate performance

Three ticket-sized charges ran at Opus/medium, one per ticket, serial on the retained source. All three returned diff-ready with green focused tests and a biting self-probe. No charge received a spec slice. Ticket 02's charge composed `identityBundleRefusal` into the creation-bundle validator and flagged the precedence change it caused, which the review then graded. Four Sonnet/high review passes (three axes and one repair-scoped re-review) held the citation standard. One axis marked a finding `ask-user`; the coordinator re-disposed it as `auto-fix` with an exact predicate.

## Coordinator catches

- The coordinator's first charge offered a `cd` into the worktree path beside `bench worktree exec`, and the delegate took it. The reviewer caught it. The rule is now in `craft-delegate`.
- The coordinator's review artifact broke the prose bounds twice (a 33-word sentence, then an 8-sentence paragraph). It paid two red gates before it ran the focused prose check first.
- Every delegate done-claim was verified: `git status`, the focused tests, and a probe of a different kind and site. All four coordinator probes bit. No delegate claim was false.
- The debug phase found the flaky `TestRootAndHelpAlignWrapperAndBinary` read the live tree three times; the fix isolates it.

## Repair attribution

| ticket | repair rounds | cause per round |
| --- | --- | --- |
| 01-name-the-landing-identity-component | 0 | none |
| 02-surface-the-resolver-reason-in-exec-and-path | 1 | shaping-ambiguity |
| 03-repair-review-findings | 0 | none |

## Agent-experience improvements

### Bench CLI

- Project gate and commit output to the red details and the phase verdicts, with the ten slowest stages on request only.
  Feeds: new
- Run a scoped check on a ticket commit and keep the full gate for the composed landing.
  Feeds: FT215
- Accept the spec path form for `bench worktree land --spec`, because today only the slug form resolves.
  Feeds: new
- Hold one landing lease in the ledger from composition through publish, so a second landing waits or names the owner's assignment id.
  Feeds: new
- Make the follow-on guard fail open with a warning when the binary answers `unknown subcommand`, because a stale `dist/bench` blocked every Bash call.
  Feeds: new

### Skills

- Keep the `craft-delegate` rule that a charge names `bench worktree exec` as the only command form; it landed as `efbef44a`.
  Feeds: none

### Process

- Run the focused prose checks before `bench commit` on any authored Markdown, so a bound violation costs seconds instead of a full gate.
  Feeds: none
- Make `bench handoff` rewrite the State section too, because the pin refresh left the prior phase's State in place.
  Feeds: new
