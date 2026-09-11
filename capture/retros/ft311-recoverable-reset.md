## Outcome

The FT311 recoverable-reset spec landed on `main` at `5201c9fc21f0dd85563edb704f401294816edeb7` on 2026-09-11, from source tip `579c2e7380582b6af0c44ef4301eda2b33422ff7` over base `cf2caa6d5c195c6b6ed5852fb4478902d1d99fbd`. The landing gate ran green on tree `5ba37eaa`, and the spec flipped to `Status: implemented` in that commit.
The build delivered `bench worktree reset` with its plan, apply, and restore modes, the reset envelope and its namespace, the record-bound sweep rule, and the guidance fold. The source was the frozen benchmark candidate `f5667b6b` from the Sol-orchestrated Astra arm, plus one repair round in this session. The repair added the ignored-path collision refusal in both directions, the index-bound fingerprint, the clean-before-reset move order, the hidden-flag refusal, and the target-side checkpoint resolution.

The landing retired the frozen candidate worktree and its branch through the eligible-sibling cleanup. The candidate tip stays reachable in main's history, and the benchmark evidence directory stays intact outside the repository until the phase-close landing declares it.

## Gate-stage timings

The landing's own gate, from its stdout on 2026-09-11:

| stage | elapsed ms |
|---|---|
| gofmt | 112 |
| vet | 1145 |
| test | 85057 |
| race | 3066 |
| system | 27268 |
| shellcheck | 544 |

The scaffold named the landing commit and trace `6b77784da6f7972894248fbe980a8e8e` and carried the same six stages.

## Ticket-versus-spec-slice and delegate performance

The five tickets ran in the benchmark arm as ticket-sized charges to one retained Astra/medium author under a Sol/high coordinator, with Terra/medium axes. Ticket 3 needed a same-model replacement context after two auto-review refusals of an authorized test replacement. Every ticket landed with a biting author probe and a biting coordinator probe, except ticket 5, whose prose probe was silent.

This session ran no write delegate. The repairs ran in the coordinator's own context at the reviewer's direction. The three Sol/high axes and the three Opus/medium re-review axes ran as fresh isolated delegates. The Sol/high round returned 26 findings and 14 accepted targets. The Opus/medium re-review returned 7 minor findings, verified every fold predicate, and observed five mutations red.

## Coordinator catches

- The Spec axis found that the first collision check walked ancestors only, so an ignored file under a tracked directory was still replaced. The fold added the descendant direction and a fixture for both.
- The Coverage axis found that `clean -fd` after the hard reset deleted files that a newer ignore rule covered. The reviewer chose the clean-first order, and RR73 pins it.
- The Coverage axis found that assume-unchanged and skip-worktree entries escaped every fingerprint input, and that a symbolic `--to` resolved in the primary checkout. The reviewer chose refusal and target-side resolution.
- The Standards axis found the reset command layout spelled twice, the staged listing read twice, and the status argv spelled three times. The fold gave each one owner, which needed the cleanup planner file in the fence.
- The landing refused twice on the review base, because an assignment created from a sibling has that sibling's tip as its start. An integration assignment from main folded the repair and landed.
- The landing refused on undeclared ignored residue, because the benchmark evidence directory ignores itself. The evidence moved aside and back, and the phase-close landing declares the path.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 1.md | none | none |
| 2.md | 1 | spec-row |
| 3.md | 3 | spec-row, delegate-error, shaping-ambiguity |
| 4.md | 2 | spec-row, spec-row |
| 5.md | none | none |

## Agent-experience improvements

### Bench CLI

- The landing's census entry reads `census=0` for the integration assignment, and the repair assignment recorded one raw call before its release.
  Feeds: none
- Let `bench worktree land` accept main's tip as the review base for a `--from` sibling, so a repair of a frozen candidate lands without an integration assignment.
  Feeds: new
- Let the landing's residue refusal name a `bench`-native declaration step instead of `rm -rf` over the offending paths.
  Feeds: new
- Add one create form that opens the three review-axis worktrees from one source assignment and prints one table.
  Feeds: new

### Skills

- Add to the review discipline that a Coverage axis chains an exit-3 record's named `next` command through its apply before it accepts the row.
  Feeds: none

### Process

- State to the reviewer before a landing that the landing retires every sibling whose commits it carries, when a sibling must survive as frozen evidence.
  Feeds: none
- Declare a self-ignoring capture directory in the build-outputs file before the first landing that meets it.
  Feeds: none
