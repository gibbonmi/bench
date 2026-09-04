# Retro — binary-freshness

## Outcome

`binary-freshness` (FT177) landed as `89f4706f` from the reviewed source tip
`4eb50690` over the frozen base `ee339a80`. The diff held 62 files, 2577
insertions, and 186 deletions across 17 commits, and the landing released the
integration source with `census=0`.

Every consumer of `dist/bench` now grades it through the one digest primitive.
`bench commands --brief`, `bench doctor`, and `bench preflight build` each
carry a freshness verdict. The stamped build publishes the executable, its
seal, and the broker manifest as one transaction. The land route recovers an
unstamped manifest through one rebuild and keeps a digest mismatch at exit
127. The follow-on guard splits its posture on a stale answer, and the
working agreement names `bench test --check system` as the hand-run route.

All 27 acceptance rows resolve to a real test. Four were partial at the first
review and closed in repair.

## Gate-stage timings

The landing gate, in milliseconds:

| phase | verdict | elapsed |
|---|---|---|
| gofmt | green | 103 |
| vet | green | 926 |
| test | green | 58004 |
| race | green | 2747 |
| system | green | 25714 |
| shellcheck | green | 536 |

Sixteen whole-project gates ran across the build, one per merge. The test
phase held between 52 and 88 seconds throughout, and no phase went red on a
merge after its ticket's focused checks passed.

## Ticket-versus-spec-slice and delegate performance

Nine ticket charges ran as nine Opus delegates in nine worktrees, two at a
time under the test-parallel cap. Four repair charges and one guidance repair
followed. Three review axes and two repair-scoped re-reviews ran at Opus
medium.

A ticket-sized charge outperformed a spec-slice charge on premise
verification. Every one of the nine verified its ticket's stated premise with
file:line citations before it edited, and four found the premise wrong or
incomplete. The publish delegate found the verb needed one argument rather
than the two the spec named. The repair-root delegate found the kit root had
to be derived rather than threaded, because the call site sat outside every
fence. The guard delegate found the shared library was not sourced by both
hooks, which the ticket asserted.

Six of the nine landed their behavior first-pass. Three stopped and reported
a shortfall rather than editing out of fence, which is the behavior the
charge asks for.

The strongest single result came from a repair charge. It was charged to
default the manifest directory to `bin/`. The delegate found that an
unconditional default poisons the wrapper manifest on every gate run,
because the gate's private build is a throwaway. It added
`--manifest-dir`, proved the poisoning, and flagged the design fork.

## Coordinator catches

The coordinator ran an independent mutation probe for every accepted
done-claim, at a site and a kind distinct from the delegate's own. Thirteen
probes ran. Ten bit as expected. Three came back silently green, and each one
was a missing row:

- The build script's version argument was ungraded. BF14's test called
  `freshness.Publish` directly with a literal, so a build stamping `dev`
  passed. A real subject-mode build row closed it.
- The shared guard classifier table pinned no `|` or `&` segment. Narrowing
  the splitter left every row green, so `rg foo | bench help` would have run
  stale. Six rows closed it.
- `dist/bench` was spelled three times in production. Renaming the preflight
  constant left its own new test green, because the test read the same
  constant. `freshness.PublishedExecutable` now owns it, and a fourth
  derivation surfaced during that collapse.

Two coordinator traces produced blocking findings the delegates had not seen.
The shim's new envelope reader classified a plain `ls` as a Bench call on a
real envelope carrying `cwd`. The build's manifest landed in `dist/` while
every reader looked in `bin/`.

The landing itself falsified a recorded decision. It refused at exit 127 for
an absent broker manifest, the state the doctor row reports as `ok:`.

## Repair attribution

| ticket | rounds | cause per round |
|---|---|---|
| add-a-binary-seal-row-to-preflight | 1 | ticket-slicing |
| add-seal-and-broker-rows-to-doctor | 0 | none |
| name-the-system-suite-route-in-guidance | 1 | ticket-slicing |
| name-the-repair-root-in-the-landing-refusal | 0 | none |
| split-the-guard-posture-on-a-stale-answer | 4 | ticket-slicing, ticket-slicing, delegate-error, spec-row |
| publish-the-broker-from-the-stamped-build | 1 | spec-row |
| recover-a-dev-manifest-in-the-land-route | 0 | none |
| compose-the-landing-line-by-checkout-kind | 0 | none |
| verify-identity-in-commands-brief | 1 | delegate-error |
| repair-the-shim-envelope-read | 0 | none |
| repair-the-preflight-facts-and-the-comment-record | 2 | spec-row, spec-row |
| publish-the-manifest-where-the-routes-read-it | 2 | spec-row, delegate-error |
| unquote-the-head-word-in-the-shell-test | 0 | none |

Seven of thirteen tickets landed in one pass. `ticket-slicing` dominates the
first half: three rounds went to a fence that omitted a file the work had to
touch, and every one was a test or a registry the ticket's own change moved.

## Agent-experience improvements

### Bench CLI

- The landing's census entry records 64 raw calls across 11 worktrees; add a
  `bench worktree read` verb for the 53 file reads.
  Feeds: FT254
- Let `bench worktree exec` accept a redirection after the `--` separator,
  which `.bench/BENCH.md` already permits.
  Feeds: FT254
- Give `bench test --check` a `--run` filter and a `--failures full`
  projection.
  Feeds: new

### Skills

- Make `craft-delegate` require `go vet` in any charge that writes a Go file
  under `tests/canary/`.
  Feeds: new
- Make `craft-review` require the Coverage axis to probe a test's fixture
  source, not only its assertion.
  Feeds: new

### Process

- Make the coordinator probe a delegate's fixture independence as well as its
  production code.
  Feeds: new
- Require a landing rehearsal before the first landing of any spec that
  changes the promotion broker.
  Feeds: new
- Add the file a ticket's own change will move to its fence at slicing time.
  Feeds: new
