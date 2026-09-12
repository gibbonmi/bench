# Budget evidence

Date: 2026-09-11.
Bench version: `bench 0.2.0 (linux/amd64)`.
Harness version: Codex `0.153.4`, from the pinned record's `session_meta.cli_version`.
Source tip: `9ecc1d561d918a3270a0894e5313b9ad3aabb0f5`.

This report is evidence, not policy. It proposes a byte bound for each surface.
The reviewer approves or refuses each proposal separately. No default changes
because this report exists, and each command owner keeps its policy under FT173.

## How each number was made

Each case ran as a real command in the assignment worktree at the source tip.
The measure is the byte count of the captured output, taken with `wc -c`.
The capture files were removed after the measurement.
A baseline is the output an agent receives today.

A candidate is a bounded form that exists today. A surface with no bounded form
records `unavailable`, and it receives no proposed number. A surface that is
already bounded records its baseline bytes as its candidate.

Input digest names what the command read. A command that reads the tree names the
source tip. A command that reads one file names that file's blob at the source
tip. A census-backed command names the worktree census and the run date instead,
because the census and the clock decide its output.

## Case inventory

Cases 1 to 20 cover raw file, Git, test, and shell reads, Bench queries, and the
harness wrapper. They include a large history, a long diagnostic, multibyte text,
empty output, a failed command, and a file with no sections.

| # | Surface | Task the case names | Input digest | Baseline bytes | Candidate bytes |
| --- | --- | --- | --- | ---: | ---: |
| 1 | Git read | Name what landed recently | tip `9ecc1d56` | 298,347 | 1,438 |
| 2 | Git read | Name what the repair commit changed | commit `983df3ad` | 59,869 | 1,507 |
| 3 | Raw file read | Find one section of the spec | `spec.md` blob at tip `9ecc1d56` | 19,827 | 480 |
| 4 | Shell read | Find every `Command` reference in one package | tip `9ecc1d56` | 2,278 | unavailable |
| 5 | Test read | Report one package's test result | tip `9ecc1d56` | 158 | 158 |
| 6 | Test read | Report the command package's test result | tip `9ecc1d56` | 151 | 151 |
| 7 | Long diagnostic | Report structural debt | tip `9ecc1d56` | 8,714 | unavailable |
| 8 | Empty output | Vet the whole module | tip `9ecc1d56` | 0 | 0 |
| 9 | Failed command | Refuse an unknown harness name | tip `9ecc1d56` | 50 | 50 |
| 10 | Bench query | Name the changed paths | tip `9ecc1d56` | 1,847 | 1,847 |
| 11 | Bench query | Read the roadmap index | tip `9ecc1d56` | 1,970 | 1,970 |
| 12 | Bench query | Find one symbol's reference edges | tip `9ecc1d56` | 1,166 | 1,166 |
| 13 | Bench query | Name what needs attention | worktree census, run 2026-09-11 | 395 | 395 |
| 14 | Bench query | List the owned worktrees | worktree census, run 2026-09-11 | 2,351 | unavailable |
| 15 | Bench query | Report one spec's coverage state | tip `9ecc1d56` | 5,245 | unavailable |
| 16 | Bench packet | Produce the full review charge | tip `9ecc1d56` | unavailable | unavailable |
| 17 | Bench packet | Report the review phase-entry checks | tip `9ecc1d56` | 547 | 547 |
| 18 | Multibyte read | Observe one 12.8 MB session record | `3b40464b` | 3,432 | 3,432 |
| 19 | Wrapper result | Return a bounded oversized tool result | none | unavailable | unavailable |
| 20 | Raw file read | Read one Go file with no sections | `internal/anchors/registry_data.go` blob at tip `9ecc1d56` | 146,885 | unavailable |

Case 16 is unavailable for one recorded reason. `bench preflight review --charge
--full` refuses a dirty checkout, and the checkout carried the coordinator's
in-flight spec edit. The refusal line is 113 bytes, and it is not the packet.
Case 19 is unavailable because the wrapper belongs to the overflow child, which
is staged and not built. Case 20 records no candidate because a file without
sections has no named-section route.

## Task verdicts and recovery calls

A verdict says whether the case's task still succeeds from the candidate alone.
A recovery call is one follow-on invocation that reaches the omitted detail.
A zero is a measured zero only where the task ran and needed nothing more.

| # | Verdict from the candidate alone | Recovery calls | Recovery route and its bytes |
| --- | --- | ---: | --- |
| 1 | Succeeds for recent history; fails for an older commit | 1 | `git log --oneline --grep ME-C1`, 368 |
| 2 | Succeeds for the changed paths; fails for the changed lines | 1 | `git show 983df3ad`, 59,869 |
| 3 | Succeeds for the section names; fails for the section text | 1 | one ranged read of the named lines |
| 4 | Not measured; the surface has no bounded form | unknown | none today |
| 5 | Succeeds | 0 | `bench test --full` returns the same 158 bytes |
| 6 | Succeeds | 0 | none needed |
| 7 | Not measured; the surface has no bounded form | unknown | none today |
| 8 | Succeeds | 0 | none needed |
| 9 | Succeeds; the refusal is the whole answer | 0 | none needed |
| 10 | Succeeds for the changed paths; fails for the diff body | 1 | `bench diff --full`, 483,385 |
| 11 | Succeeds for the index; fails for one row body | 1 | `bench roadmap --context --row FT204`, 18,161 |
| 12 | Succeeds; the default cap never binds on this symbol | 0 | `bench consumers harnesses.Command --full`, 1,173 |
| 13 | Succeeds; the `+2 more` line discloses the omission | 1 | `bench status --all`, 584 |
| 14 | Not measured; the surface has no bounded form | unknown | none today |
| 15 | Not measured; the surface has no bounded form | unknown | none today |
| 16 | Not measured | unknown | none today |
| 17 | Succeeds | 0 | none needed |
| 18 | Succeeds | 0 | the reader states each unknown dimension itself |
| 19 | Not measured; the surface is not built | unknown | none today |
| 20 | Not measured; no named-section route exists for a file without sections | unknown | none today |

Case 12 carries a warning for the reviewer. The default and the full route differ
by 7 bytes, so this symbol never reaches the cap. The case shows that the route
works. It is not evidence about a capped result.

Cases 2 and 3 name one command each, so another reader reaches the same bytes.
Case 2 captured `git show 983df3ad`, and its candidate captured
`git show --stat 983df3ad`. Case 3 captured
`git show 9ecc1d56…:specs/session-context-measurement/spec.md`, which is the
blob at the source tip and not the later working file.

## Harness observations

The record view read the pinned record once. These are session-level facts, and
they come from the harness transcript alone. The spec's Dogfood runs note is the
record of record for these fourteen facts, and this table repeats them.

| Fact | Value |
| --- | --- |
| Record digest | `sha256:3b40464b598e2a977fddb4f249ff01ff6ec9881e6e7435d6f22ffe1613a60468` |
| Observed interval | `2026-09-11T09:00:11.841Z` to `2026-09-11T12:08:59.066Z` |
| Result-text bytes | 1,374,543 |
| Result-text characters | 1,371,945 |
| Result-text lines | 16,071 |
| Outer calls | 403 |
| Unmatched calls | 0 |
| Turns | 5 |
| Compactions | 2 |
| Input tokens | 52,339,083 |
| Cached-input tokens | 51,375,872 |
| Output tokens | 135,788 |
| Reasoning tokens | 58,203 |

The reader reports three dimensions as unknown for this source. They are observed
nested calls, explicit read paths, and per-result token attribution. The spec
states that this reader gives no per-result attribution. So no proposal below
converts bytes into tokens or into money.

The upper tail carries most of the cost. A separate pinned slice reports that
results of at least 10,000 characters are 31 calls and 73.2 percent of all
result text. That figure comes from
`specs/session-context-efficiency/decisions/session-context-efficiency/assets/preliminary-assessment.md`,
and this report cites it rather than recomputing it.

## Bench-owned facts

These facts have Bench owners. They are recorded here separately from the harness
observations above, and no harness text supplied any of them.

| Fact | Value | Bench owner that reported it |
| --- | --- | --- |
| Gate stage elapsed, gofmt | 139 ms | `bench retro session-context-measurement --scaffold` |
| Gate stage elapsed, vet | 1,376 ms | the same scaffold |
| Gate stage elapsed, test | 111,947 ms | the same scaffold |
| Gate stage elapsed, race | 3,505 ms | the same scaffold |
| Gate stage elapsed, system | 36,558 ms | the same scaffold |
| Gate stage elapsed, shellcheck | 708 ms | the same scaffold |
| Landing commit and trace | `244a8acb`, trace `b4fee53d` | the same scaffold |
| Raw calls in this assignment | 5 | `bench status --all` |
| Raw calls in `delegate-boot-cost` | 3 | `bench status --all` |
| Registered worktrees | 7 | `bench worktree list` |

The scaffold reads the OTEL span record through its existing command. This report
adds no second reader and changes no span format.

## Proposed budgets

Each row proposes one bound for the reviewer. Each row names the owner child that
would carry the change. No row is approved, and no default changes today.

| Surface | Proposed bound | Owner | Evidence cases | State |
| --- | --- | --- | --- | --- |
| Raw file read | 10,000 bytes, with a named-section route, on sectioned Markdown evidence alone | queries child | 3, 20 | awaiting reviewer approval; a file without sections is evidence incomplete |
| Git read | 10,000 bytes, with a full-diff route | queries child | 1, 2 | awaiting reviewer approval |
| Bench query default | 10,000 bytes, with the existing full route | queries child, with each command owner | 10, 11, 12, 13, 14, 15, 17, 18 | awaiting reviewer approval |
| Explicit full route | no bound | each command owner | 10, 11 | awaiting reviewer approval |
| Test read | no change; already bounded | `bench test` owner | 5, 6, 8 | awaiting reviewer approval |
| Shell read | none proposed | queries child | 4, 7 | evidence incomplete |
| Harness wrapper result | none proposed | overflow child | 19 | evidence unavailable |

The 10,000-byte figure has one reason. The cited slice shows that the upper tail
above 10,000 characters holds 73.2 percent of all result text. Every measured
Bench query default is 5,245 bytes or fewer. So this bound would refuse the upper
tail and would change no measured default today.

Each proposed bound is in bytes, never in characters. A multibyte result reports
fewer characters than bytes, as case 18 shows with 1,374,543 bytes over 1,371,945
characters. A character bound would therefore admit more bytes than it names.

The raw-file bound rests on sectioned Markdown evidence alone. Case 3 shows the
section route on a Markdown file. Case 20 reads a 146,885-byte Go file with no
sections, and no named-section route reaches it. So the raw-file row records a
file without sections as evidence incomplete, and the reviewer decides that part
separately.

Two rows propose no number. The shell read has baselines but no candidate and no
task verdict, so nothing matches. The wrapper result has no built surface.
Unavailable evidence cannot justify a default.

## Reviewer checkpoint

This report proposes per-surface bytes and approves no default. Current policy
stays until the reviewer decides each row above.
