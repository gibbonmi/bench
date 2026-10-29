## Outcome

The `git-query-cache` build landed at 2fb468a50 and closed FT402-git-cache. A command run now caches its location answers, its administration path, and each read whose revision operands are all full object ids. The ledger, the commitment store, the roadmap revision reads, and every worktree verb share the verb's one run.

The `internal/worktree` package ran in 128.6 to 129.9 s against the recorded 146.2 to 153.8 s baseline. The smallest gain is 16.3 s and the median gain is 19.6 s, so GQ42 holds in its literal reading. A same-session control at `main` gave 141.5 to 143.9 s.

GQ41 missed its literal reading after GQ-C2B: the smallest gain was 6.6 s against the 7.6 s cross-day spread. A same-session control showed 4.6 s of host drift, and a Fable tiebreak continued the build as a reviewer-vetoable deviation.

## Cost

- landing: 2fb468a50293a6e6b31415bca9db4d6550ae01be, range 2fb468a50293a6e6b31415bca9db4d6550ae01be^1..2fb468a50293a6e6b31415bca9db4d6550ae01be^2
- wall: 269 min
- commits: 43 total, 34 bookkeeping, 9 change
- graph: 5 tickets, chain 5, 4 chunks, 4 serial checkpoints
- dispatch: 5 authors, 4 repairs
- trace: commit 46 spans in 410446 ms, lane 0 spans in 0 ms, landing 1 spans in 11742 ms
- wall per ticket: 54.0 min against 60 min
- bookkeeping per ticket: 6.8 against 7.1

Provider usage and charges remain unknown. The harness gave only one total token count for each delegate.

## Gate-stage timings

- landing: commit 2fb468a50293a6e6b31415bca9db4d6550ae01be, trace 7c57d01e670def93de6e8afb80701926
- The landing gate reused the fresh green of the final whole gate, so it ran no phase.
- The first GQ-C2A checkpoint went red once on the FT372 evidence-absent flake, and the retry ran green.
- Each GQ41 and GQ42 timing took about 12 minutes for five runs, and each control took the same.

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| ticket 02 | GQ27 to GQ33 each have a demonstrated red | verified | 9 | 1 | Opus / xhigh / implementer |
| ticket 02 | the plan probe-02 reads bit | verified | 10 | 1 | Opus / xhigh / implementer |
| ticket 03 | the threaded landing pins 54 starts | verified | 8 | 1 | Opus / xhigh / implementer |
| ticket 04 | GQ35 to GQ37 each have a demonstrated red | verified | 9 | 1 | Opus / xhigh / implementer |
| ticket 05 | GQ38 to GQ40 each have a demonstrated red | verified | 9 | 1 | Opus / xhigh / implementer |
| ticket 05 | the GQ34 pin reaches 33 | verified | 8 | 1 | Opus / xhigh / implementer |

Brier mean 0.018 over 6 pairs; 0 abstained. Ticket 01 returned no integer confidences, so its claims are not labeled.

Ticket 01 met its acceptance in one pass, but Coverage found four surviving key and guard mutants. Tickets 03 and 05 each stopped before the commit on a fence gap, with the threading done and uncommitted. The same author continued after a widen. Ticket 04 needed no repair.

## Coordinator catches

- The checkpoint flake was the known FT372 evidence-absent red; the test passed alone, and one retry ran green.
- The GQ-C1 repair touched the ADR, so `bench record confirmation` refused, and the chunk took a confirming round.
- Ticket 03 needed the stub-count wrapper in the one file that the serial journey census admits. Ticket 05 exposed the whole target-resolve chain when it removed the last helper run.
- After the confirming rounds, two later findings only lost cache at a site that no row counts. The orchestrator recorded them as advice under the bounded repair policy.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 01-cache-full-object-id-reads.md | 1 | edge-found |
| 02-thread-the-run-through-the-ledger-and-store.md | 1 | one-source |
| 03-thread-the-landing.md | 1 | delegate-error |
| 04-thread-the-cleanup-and-release-paths.md | 0 | none |
| 05-thread-reset-lifecycle-and-create.md | 1 | delegate-error |

## Agent-experience improvements

### Bench CLI

- Add a build step to `bench test --in`, so that a test run after a Go edit rebuilds the worktree binary instead of refusing. Three delegates asked for it.
  Feeds: FT406
- Let `bench learning` take `--in <label>`, as the other capture verbs do.
  Feeds: new
- Print the first failing assertion line in each row of a `bench probe --mutants` answer, as the single form does.
  Feeds: new
- Census account: the 44 `rg` calls are sanctioned reads and need no Bench form. The 4 `cd` calls are a missed form, because each was a refused `cd` into the pool path where `bench worktree exec` serves. So name `bench worktree exec` and the pool-path refusal in the delegate agent definitions.
  Feeds: new
- Census account: the 2 `git diff` calls are sanctioned reads of the frozen review pair. The 1 `wc` call is a missed form, because `bench structure` reports file budgets. So name `bench structure` for a line-budget read in the author charge.
  Feeds: none

### Skills

- In `craft-tickets`, require the slicer to trace the whole caller chain of each removed construction and each package harness census before ticket approval.
  Feeds: new
- In `craft-spec`, require that a slice timing row grade against a same-session paired control.
  Feeds: new

### Process

- Run each timing before a record commit, because a commit starts a background grade that loads the host.
  Feeds: none
- Let `preflight --charge` bind a continuation of the same author on a dirty tree, or name the continuation as exempt, because each fence-gap stop left uncommitted work.
  Feeds: new
- Omit the staleness pass when no code path changed since the spec's measured base, as this build and `build-cost` both did.
  Feeds: none
- Optimization used: skip the candidate install when the spec declares no `Persisted fields:` line, because the build is then not a declared build.
  Feeds: none
- Optimization used: give one combined repair session every accepted finding of a round, and widen its ticket for a stray path, so that each chunk closed in one repair cycle.
  Feeds: none
- Optimization used: run `bench consumers` for the next ticket while the chunk checkpoint runs, and widen in the assignment commit, so that no author stopped on a known caller.
  Feeds: none
- Optimization used: record a cache-only finding at an uncounted site as advice, and narrow the row text instead of adding serial count tests.
  Feeds: none
- Optimization used: let the orchestrator write the author witness reds into the review record, which is evidence-only and costs no repair cycle.
  Feeds: none
- Optimization used: continue a stopped author in its own session after a widen, instead of a fresh author that rereads the seam.
  Feeds: none
- Optimization used: land the retro, the scorecard, and the spec retirement from one closing worktree, so that the close pays one landing gate.
  Feeds: none