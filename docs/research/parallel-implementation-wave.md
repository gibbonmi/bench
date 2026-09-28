# Parallel implementation wave, 2026-09-22

## Recommendation

Re-author or re-review the overflow stream against the current `main`, because `main` changed its spec after the wave stopped.
The reviewer dropped the shared-delegate-startup stream on 2026-09-28, because it contradicts the current fresh-author and fresh-repair rules.
The reviewer chose on 2026-09-28 to re-author the queries stream on `main`, with its tip as the reference source for fresh ticket authors.

Scope: the Codex parallel implementation wave on three staged specs, its measurements, and its model findings.
Evidence status: each stream passed its own chunk checkpoints on its branch. No wave source landed on `main`.
Read date: 2026-09-26
Wave base: `d7d4e8226d63c1996975385b77ec784c30c4794d`
Roadmap owner: FT348
Consumed by: the FT348 disposition decision for each stream
Drift: refresh after a stream composes a later `main`, after a new review, or after a spec change on `main`.
Retire when: each stream lands, is re-authored, or is dropped by a reviewer decision.

The wave source lives only on its branches. The recovery table at the end pins each branch and tip.
This report cites a branch file as `<tip>:<path>:<line>`.

## Question graph

1. What did each stream deliver, and where did it stop?
2. What did the selected-query measurements show?
3. What did the overflow runtime evidence prove?
4. What changed on `main` after the wave stopped?
5. What did the wave show about the implementation and review line?

Question 4 constrains the disposition of each stream in question 1.
Questions 2 and 3 feed the two blocked decisions.

## Q1. What did each stream deliver?

### Facts

| Stream | Retained tip | Acceptance state at the stop | Blocker at the stop |
| --- | --- | --- | --- |
| shared-delegate-startup | `d02df85bf14c360c086e2475d321a3127f83b4a6` | All 52 rows reconciled. Adoption, three review axes, and the complete checkpoint passed. | The landing refused four dirty primary-checkout paths. |
| session-context-queries | `194b7dba932a2881a085650789fa4c98d6478feb` | 23 of 25 rows accepted through chunk QU-C3. | Chunk QU-C4 needed a numeric budget decision. |
| session-context-overflow | `208371fecd1cfee5a09a1af9bba1df6a1effd109` | OV1, OV20, OV21, and OV22 accepted for chunk OV-C1. 21 product rows remain. | Chunk OV-C2 needs adapter, budget, and retrieval decisions. |

Source: `52485ec9:capture/retros/parallel-implementation-wave.md:11`.

The startup stream edits guidance only. It adds no command, flag, or default route.
Source: `1cb8dae7:specs/shared-delegate-startup/spec.md:66`.
The queries stream adds selected views to `bench worktree list` and `bench spec history`.
Source: `194b7dba:specs/session-context-queries/tickets/1-select-worktrees.md`.
The overflow stream adds runtime evidence only. It enables no adapter.
Source: `208371fe:specs/session-context-overflow/assets/runtime-evidence.md:3`.

Each stream kept one retained Astra/ultra author and serial tickets.
Independent Sol/high sessions reviewed each chunk on three axes.
Source: `52485ec9:capture/retros/parallel-implementation-wave.md:94`.

### Inference

The primary checkout was clean at `c8c444ffae2fb1578cfa54a22fa632590ffbc322`. So the original startup landing blocker is gone.
The startup stream still needs a composition with the current `main` before it can land.

## Q2. What did the selected-query measurements show?

### Tested results

The measurement ran at tip `194b7dba` on 2026-09-23.
Source: `52485ec9:capture/retros/parallel-implementation-wave.md:687`.

| Task | Bytes |
| --- | ---: |
| One-target selected worktree view | 300 |
| Complete worktree inventory, 55 rows | 16938 |
| Newest event for one selected history | 322 |
| Complete histories of three, two, and zero events | 245, 218, 36 |

Source: `52485ec9:capture/retros/parallel-implementation-wave.md:707`.

The selected worktree view answered the path task without a follow-on call.
The complete inventory needed three more path calls for the same facts.
For the small histories, the selected form used more bytes than the complete output: 640 against 499.
Source: `52485ec9:capture/retros/parallel-implementation-wave.md:729` and `:833`.

### Proposal that was open at the stop

The coordinator proposed a 10000-byte default for the worktree inventory only.
The history default was to stay unchanged. No user approved this proposal.
Source: `52485ec9:capture/retros/parallel-implementation-wave.md:788`.

The raw captures are in the ignored file `.logs/query-measurement-20260923.json` of the queries assignment.
Their bundle digest is `42ac077de98ec3738504e9f7827ffcec4cd9f240c63f85f50ab453efc84f56c8`.
Source: `52485ec9:capture/retros/parallel-implementation-wave.md:699`.

## Q3. What did the overflow runtime evidence prove?

### Tested results

The tested Claude persisted-file success path keeps and retrieves complete text and binary bytes.
The binary case has 262144 bytes, with 16384 NUL bytes and invalid UTF-8.
Source: `52485ec9:capture/retros/parallel-implementation-wave.md:157`.

The evidence covers Codex CLI 0.155.1 and Claude Code 2.1.280. OpenCode was absent.
Source: `208371fe:specs/session-context-overflow/assets/runtime-evidence.md:25`.

### Limits

The evidence proves no failed-output preservation and no adapter eligibility.
It sets no numeric budget.
Source: `52485ec9:capture/retros/parallel-implementation-wave.md:159`.

## Q4. What changed on `main` after the wave stopped?

### Facts

`main` is 359 commits past the wave base at `c8c444ffae2fb1578cfa54a22fa632590ffbc322`.
Commit `b7667d8f` made the FT336 spec own the default output budget and retired queries ticket 4.
So the QU-C4 numeric decision is now FT336's, and the queries stream carries a ticket that `main` retired.
Source: `4bdffd04:specs/session-context-queries/spec.md:20`.

No roadmap row owns the three staged specs. FT172 holds that gap and forbids a parent inferred from a shared file.
Source: `roadmap/FT172.md:27`.

### Inference

The queries stream needs a spec reconciliation before a merge from `main`.
The startup and overflow streams touch files that later work also changed, so each composition needs a reconciliation against both specs.

## Q5. What did the wave show about the line?

### Facts

The user moved the remaining reviews to GPT-6 Sol/high during the wave. The three authors stayed on Astra/ultra.
Source: `52485ec9:capture/retros/parallel-implementation-wave.md:253`.
The retained calibration was a mean Brier error of 0.007931 over 29 labeled pairs, with one abstention.
Source: `52485ec9:capture/retros/parallel-implementation-wave.md:133`.
Token counts, provider costs, and comparative latency are unknown.
Source: `52485ec9:capture/retros/parallel-implementation-wave.md:135`.

The wave used one raw `python3` call across 37 assignments.
Source: `52485ec9:capture/retros/parallel-implementation-wave.md:34`.

The retrospective and the scorecard update are on the capture branch only.
They are `capture/retros/parallel-implementation-wave.md` and `capture/agent-performance/open-ai-models.md` at `52485ec9`.

## Residual unknowns

- The cost of the wave, in tokens or money, is unknown.
- No stream has run a gate after a composition with the current `main`.
- No representative oversized history or empty worktree inventory was measured.
- The overflow eligibility of any path other than the tested Claude success path is unknown.

## Verification record

- [x] The report opens with its recommendation, scope, and evidence status.
- [x] Facts, tested results, inferences, and proposals are in separate sections.
- [x] Each question has one section.
- [x] Each load-bearing claim cites a branch blob or a `main` path.
- [x] Each stream branch and its exact tip are recorded.
- [x] Unknowns stay explicit.

## Validation plan

1. The reviewer picks a disposition for each stream: compose and land, re-author, or drop.
2. For a stream to land, compose the current `main` into its branch with `bench worktree merge` and reconcile the shared files.
3. Run the complete checkpoint on the composed tip before the landing.
4. For the queries stream, reconcile the retired ticket 4 with the FT336 budget owner first.
5. For the overflow stream, record the adapter, budget, and retrieval decisions before chunk OV-C2.

## Recovery table

Each stream branch below keeps its commits. The review-axis branches of the wave hold commits that are ancestors of these four tips.

| Label | Branch | Tip |
| --- | --- | --- |
| wave-shared-startup | `bench/assign/e2fbcdd698f090a3967ac1fecd8b2e10/71e28d5e8d145f6d8c121629d2f73093` | `d02df85bf14c360c086e2475d321a3127f83b4a6` |
| wave-context-queries | `bench/assign/3458eb8d3cac8810295a65f9fdf866f2/910100df1a45f972d086be3f35818c63` | `194b7dba932a2881a085650789fa4c98d6478feb` |
| wave-context-overflow | `bench/assign/c7f0abc2d0d559713fb6a5a32bb0af75/639b85dab8d7793586a45fdba9cb93e6` | `208371fecd1cfee5a09a1af9bba1df6a1effd109` |
| wave-capture | `bench/assign/b1674c8aa41f42ede525752183a219e3/0f91864f110f4f4130fa50e83cda39fc` | `52485ec90e3c20de0be82ccad5da74e8fe8d099f` |
