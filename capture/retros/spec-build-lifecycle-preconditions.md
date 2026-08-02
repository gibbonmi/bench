# Retro: spec-build-lifecycle-preconditions

## Outcome

All six stories landed as ten path-scoped gated commits on `main` across two
sessions: `6df669d` (story 5 resolver messages), `e6652e4` (story 1 reorder),
`5e17d6a` + `5fbb241` (story 4, landed by the parallel pcgs session to
undeadlock `start`), `5fc1b55` (story 2), `9166960` (story 6), `dfcc71d` +
`eec8b22` (story 3, both packages), `e7b8569` (story 5/4 CLI enumeration),
`5cb645b` (review repairs). Light-path per-ticket commits, not a lifecycle run:
a second concurrent lifecycle on `main` would have recompose-churned against
the active per-component-gate-scoping run. Fresh-context three-axis review
found no promotion blocker; four concrete findings were repaired and landed,
three reviewer-call findings are parked in `capture/IDEAS.md` (restart marker
path, husk/symlink liveness, plus the earlier resume crash window).

## Gate-stage timings

Stage-level timings were not captured. Seven whole-tree gate runs, each green
first try; wall clock dominated by the Go test phase. Focused package suites:
`internal/specbuild` 26–31s, `internal/worktree` 21–29s, `internal/conformance`
~37s, `internal/contract/runtime` full ~58s. `dist/bench` freshness forced two
manual rebuilds before `bench commit` would run (tree had moved under the seal).

## Ticket-versus-spec-slice and delegate performance

Ticket-sized charges performed uniformly better than a spec slice would have:
every delegate stayed inside its ownership fence and returned a verifiable
claim. Token spend scaled with seam complexity, not ticket text: exempt-abandon
(sonnet) ~84k; split-ownership (opus) ~144k; enumerate-CLI (sonnet) ~210k/106
tool uses — the CLI contract suite is the expensive seam because every probe
needs a binary rebuild. Salvaged draft patches from the interrupted prior
session were useful as orientation only; both drafts that were adopted were
re-derived and one (exempt-abandon) had its fixture simplified against the
draft.

## Coordinator catches

- The exempt-abandon delegate pasted its why-comment twice; caught at diff
  inspection, deduplicated before landing.
- My own CT4 probe first hit the nearly-unreachable completion-helper bootstrap
  site (`lifecycle.go`) and passed suspiciously fast; re-derived the
  fresh-start site (`assign.go`) and got the genuine red. The spec's
  "two sites, one nearly unreachable" note was load-bearing for verification.
- Every landing used a coordinator mutation probe of a different kind than the
  delegate's own probe; all bit.
- The review delegate's C1 (restart-after-terminal still refuses an advanced
  marker) is a real residue of story 4 that no ticket fence covered — the
  strongest argument for the fresh-context review round.

## Agent-experience improvements

### Bench CLI

- `bench worktree clean` plan output has no stable way to extract the
  fingerprint (awk over a TOON row; multi-row plans shift the field). A
  `--fingerprint-only` plan mode or a keyed line would remove a scripting trap.
- The ignored-inventory destructive limit refuses `clean --discard-ignored
  --full` outright when `dist/.freshness-go-cache` exceeds it; the only route
  is a manual `rm -rf` of the cache before re-planning. The derived-cache case
  deserves a carve-out or a named override.
- The Claude WorktreeCreate hook admits one active hook-created assignment per
  session, so parallel `isolation: worktree` agent launches fail after the
  first; the workaround is manual `bench worktree create --request <uuid>` per
  delegate. Either lift the limit or document the manual route as canonical.

### Skills

- No skill defects observed. The implement-spec `--full` orchestration text was
  followable as written; the resume-from-handoff entry worked cold.

### Process

- Two concurrent sessions sharing one `main` worked under serialized gate runs
  plus rebase-forward landings; the tip moved twice mid-build with zero
  conflicts because ownership fences were disjoint. The handoff's two-thread
  pin was what made the resumption unambiguous.
- `bench idea` dirties `capture/IDEAS.md` and blocks the next path-scoped
  commit until named; folding the capture file into the landing commit worked
  but the interaction is easy to trip on.
