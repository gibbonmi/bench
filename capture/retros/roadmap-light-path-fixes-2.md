# Retro: roadmap-light-path-fixes-2

## Outcome

The spec landed at `dbbec276` on `main`, base `70ccbe8a`. Ten tickets landed:
the eight original stories, one repair ticket (`trim-over-length-prose-in-the-live-tree`)
added mid-build for a material LP5 shortfall, and one post-review repair
ticket (`repair-review-round-1`). The spec then retired at `c04d3410`,
closing seven roadmap rows (FT266, FT272, FT282, FT286, FT288, FT289,
FT291). FT115 keeps its remaining clause; the drain reduces that row later.

## Gate-stage timings

Spec landing (`dbbec276`): gofmt 98ms, vet 919ms, test 57827ms, race
5511ms, system 22752ms, shellcheck 506ms. All six phases green.

Spec-retire landing (`c04d3410`): gofmt 92ms, vet 894ms, test 56534ms,
race 2375ms, system 21439ms, shellcheck 483ms. All six phases green.

## Ticket-versus-spec-slice and delegate performance

Six of nine story tickets landed first-pass with no repair round. Three
tickets carried a mid-flight shortfall the delegate itself surfaced and
reported rather than silently patching. The label-line ticket found LP5
reds ten out-of-fence sites. The sweep ticket found one out-of-fence wait
site, and separately, two spellings that could not widen without breaking
a closed decision. The post-review repair ticket found its own named
migration site held no real literal. Every stop-and-report held; no
delegate edited outside its fence to force a green.

The two concurrent-worktree rounds (label-rule plus tickets-route; sweep
plus craft-spec) both landed clean with no merge conflict. Their
`Writes:` fences were genuinely file-disjoint, even where a fixture-closure
pin later made them look like they overlapped.

## Coordinator catches

The coordinator ran one independent mutation probe per landed ticket, each
a distinct kind and site from the delegate's own self-probe. All killed
correctly; no probe came back silently green. Three coordinator catches
went beyond probing:

- The label-line ticket's corrected rule surfaced ten genuine over-length
  sites outside its fence (LP5). Routed to the reviewer rather than folded
  in silently; the reviewer chose a ninth ticket over a fence widening.
- A repair-scoped re-review's read-only Coverage delegate left the shared
  `lpf2-integration` worktree dirty — a self-probe on
  `internal/gate/prospective_owner_test.go` it did not restore, breaking
  `go vet`. Caught before the next commit; restored with
  `git show HEAD:path > path` and reverified.
- Four separate rounds of `bench preflight review` fixture-closure and
  registry-closure reds, each naming an exact missing `Writes:` pin.
  Fixed one round at a time rather than guessing the full closure upfront.

## Repair attribution

| ticket | repair rounds | cause |
|---|---|---|
| fix-the-label-line-rule-in-the-prose-grader | 1 | spec-row |
| route-the-tickets-only-close-to-the-landing-verb | none | — |
| trim-over-length-prose-in-the-live-tree | none | — |
| derive-test-wait-deadlines-from-bounds | none | — |
| sweep-literal-wait-deadlines-in-tests | 1 | spec-row |
| state-the-fence-order-and-the-claim-words-in-craft-spec | none | — |
| bind-the-exec-only-form-to-every-caller-in-craft-delegate | none | — |
| state-the-helper-return-rule-in-craft-tdd | none | — |
| state-the-census-read-the-changelog-rule-and-the-review-base | none | — |
| repair-review-round-1 | 1 | spec-row |

## Agent-experience improvements

### Bench CLI

- Add an `--env KEY=VALUE` flag to `bench worktree exec`, so a charge needing `BENCH_CONFORMANCE_ROOT` or `BENCH_CONFORMANCE_TIER` skips the `sh -c` wrapper five delegates hit independently (census: `lpf2-integration raw calls (n=47)`).
  Feeds: new
- Have `bench test --check` resolve a check from the worktree's own registry source, not only the installed `dist/bench` snapshot (census: `lpf2-integration raw calls (n=47)`).
  Feeds: new
- Stop `bench test --check` truncating a long diagnostic line (census: `lpf2-integration raw calls (n=47)`).
  Feeds: new
- Add a per-fixture proof route, such as `bench probe --fixture <name>`, so a delegate can prove one new canary bites without the full fixture sweep (census: `lpf2-integration raw calls (n=47)`).
  Feeds: new
- Have `bench spec retire` also remove the closed `Roadmap:` rows from `ROADMAP.md` and their `roadmap/FT<n>.md` files, since the spec already names them (census: `lpf2-spec-retire raw calls (n=10)`).
  Feeds: new

### Skills

- State in `craft-delegate`'s Isolation section that a read-only delegate on a shared retained worktree restores any probe byte-exact before it reports.
  Feeds: new
- State that the coordinator confirms `git status` clean on a shared worktree before trusting a read-only delegate's finding.
  Feeds: new

### Process

- Name the five registry-closure files and relevant canary fixture pins in a mid-build ticket's `Writes:` line from the start, when it touches an already-bound registry.
  Feeds: none
