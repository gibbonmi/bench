# Retro: worktree-native-forms

## Outcome

FT254 slice 2 landed at `1925f8de` on 2026-08-30 through `bench worktree land`
over the reviewed pair `a8724eb..2945e199`. The spec carries 46 coverage rows,
six build tickets, and three repair tickets. It shipped `bench worktree build
<target>`, `bench test --check system`, `bench worktree create --from <target>`,
the preflight `next` column, and the `bench worktree exec <target> -- bench gate`
form in both inventories. The landing gate was green, and the census counted 226
raw calls on the integration worktree.

## Gate-stage timings

| stage | landing | fold 1 | fold 2 |
| --- | --- | --- | --- |
| gofmt | 0.1 s | 0.1 s | 0.1 s |
| vet | 0.8 s | 1.1 s | 0.9 s |
| test | 50.9 s | 71.0 s | 52.0 s |
| race | 5.1 s | 5.5 s | 5.3 s |
| system | 15.1 s | 17.0 s | 15.1 s |
| shellcheck | 0.5 s | 0.5 s | 0.5 s |

Each `bench commit` lane ran gofmt, prose, vet, and build in under a minute.
Each delegate's `bench test --changed` ran 37 packages in about 50 s.

## Ticket-versus-spec-slice and delegate performance

Every charge was ticket-sized; no charge took a spec slice. Six build charges
ran on Opus, four at low and two at medium. Five landed first-pass on behavior.
The build charge stopped at an out-of-fence pin: the serial-test census ceiling
in `parallel_census_test.go`, which the fence did not name. Three repair
charges and five review charges (three axes, two scoped re-reviews) ran on
Opus at low or medium, all first-pass. Two parallel siblings folded through
`bench worktree merge` with no collision.

## Coordinator catches

- The census ceiling file sat outside the fence; the fence and the ticket
  gained it before the commit.
- WF43 contradicted WF30 on the assignment-state component; the collapse to
  the active filter's sentence was accepted and recorded in the row.
- The Coverage axis's worst finding (C1, `BENCH_KIT` unset) was wrong:
  `selectedRunEnvironment` sets the variable, and the check mirrors the gate's
  `kitRoot`. One code read refuted it before the pickup file was written.
- Seven coordinator probes at distinct sites and kinds all bit; WF46 admits
  only one mutation, so the delegate's red proof stood for it.
- The prose lane refused four over-bound sentences in the pickup file and a
  repair ticket before their commit.

## Repair attribution

| ticket | rounds | causes |
| --- | --- | --- |
| add-worktree-build-verb | 2 | delegate-error; delegate-error |
| run-system-suite-as-named-check | 1 | spec-row |
| fill-preflight-next-column | 0 | none |
| add-create-from-sibling-start | 1 | spec-row |
| name-the-gate-form-in-both-inventories | 0 | none |
| record-new-forms-in-prose | 0 | none |
| repair-worktree-single-resolution-and-replay | 1 | delegate-error |
| repair-system-check-name-reservation | 0 | none |
| repair-re-review-follow-ups | 0 | none |

## Agent-experience improvements

### Bench CLI

The census entry is `worktree-native-forms-integration census 226` in
`capture/learnings.md`.

- Add `bench worktree revert <target> <path>...` that copies each path aside and restores HEAD, so a probe needs no raw path.
  Feeds: new
- Let the follow-on hook accept a non-Bench step before a `bench worktree exec` segment and a quoted `|` inside the exec child's argv.
  Feeds: new
- Narrow `bench test --changed` to the packages that hold edited files plus their reverse dependencies, or name why the set widened.
  Feeds: new
- Print the `worktree:` trailer of `bench worktree exec` on success as well, and report the child's wall time.
  Feeds: new
- Make `bench gate-prose --help` exit 0 on stdout, and make a pass print one `prose[N]{path,verdict}` table.
  Feeds: new
- Make `bench test --check <unknown>` name the valid check names.
  Feeds: FT270

### Skills

- `craft-delegate`: a charge that adds a test which binds `PATH` or the process environment names the serial-census ceiling file in its fence.
  Feeds: new
- `craft-review`: a Coverage finding on an environment variable cites the producer that sets it before it claims the variable is absent.
  Feeds: new

### Process

- Give every promise in a spec's edge inventory a coverage row; the replay promise shipped without one and became the C3 defect.
  Feeds: new
- A test expectation over a random id derives through the encoder that renders it, never by hand concatenation.
  Feeds: none
- Run `bench gate-prose` on a pickup file and a repair ticket before the commit that carries them.
  Feeds: none
