# Saving of the grade key on checkpoint-then-landing pairs and retries

Consumed by: ticket #11 of this map, then the FT314 spec stage.
Drift: a change to the landing's destination-moved refusal, the wrapper's executable selection, or the final-check order.
Retire when: the FT314 spec retires.

## Result

- Recommendation: the grade key removes about one whole gate for each spec landing only when the complete checkpoint runs under the landing broker's executable.
- With the current wrapper selection, most complete checkpoints run under the worktree's own build, so the key of ticket #3 reuses few of them.
- Scope: `main` from 2026-10-01 to 2026-10-26, the retained gate run logs of the primary checkout, and the Claude session records of this repository.
- Evidence status: the counts are tested results from git, the run logs, and the session records. The saving figures are inferences.

## Q1. How many spec landings reproduce the tree that their complete checkpoint graded?

Facts:

- 40 first-parent landing merges on `main` since 2026-10-01 flip a spec from `staged` to `implemented`.
- In 39 of them, the pre-landing `main` commit is an ancestor of the source tip. The exception is `860706428` on 2026-10-01.
- 30 of them landed from 2026-10-08 to 2026-10-26, after the complete checkpoint began to grade the published tree on 2026-10-07.
- `/bench-final-check` obtains `bench gate --checkpoint <spec-path> --complete` on the source before a spec landing (`.agents/commands/bench-final-check.md:19`).
- When the destination moved past the reviewed source, the landing refuses with `land-destination-moved` (`.agents/commands/bench-final-check.md:24-26`).
- The complete checkpoint grades `published.Tree` of the committed HEAD (`internal/gate/complete_checkpoint.go:15-24`).
- In 24 session records, each final complete checkpoint was followed by the landing with no commit or record write between them.

Inference: the composed tree of each spec landing since 2026-10-08 equals the tree that its final complete checkpoint graded.

## Q2. Under which running executable did each complete checkpoint run?

Facts:

- The session records invoked the complete checkpoint 52 times as `bench worktree exec <target> -- bench gate ...` and 2 times as `./dist/bench gate ...`.
- In a worktree, the wrapper tries `<worktree>/dist/bench` before `<primary>/dist/bench` (`bin/bench.sh:82-95`, `bin/bench.sh:241-244`).
- `bench worktree build` writes `<worktree>/dist/bench` (`internal/worktree/build.go:1-4`).
- 20 of the 27 session records that ran a complete checkpoint also ran `bench worktree build`. None ran `bench worktree build --install`.
- The landing runs only under the broker whose digest equals `bin/bench-broker.manifest` (`bin/bench.sh:406-458`).

Inference: under the key of ticket #3, a complete checkpoint in a worktree with its own build keys to another running executable. Its green does not answer the landing.

## Q3. How long does one whole gate take?

Facts:

- The primary checkout keeps the last 20 gate runs (`internal/gate/run_log.go:25`).
- 19 of the 20 retained runs ran all six phases: gofmt, vet, test, race, system, and shellcheck.
- Their median wall time is 288 s. The range is 242 s to 341 s.

## Q4. What do retries add?

Facts:

- A landing grades from the primary checkout as its identity root (`internal/worktree/land_identity.go:45-58`).
- Today the evidence store reuses a green of the same tree and subject in the same checkout (`internal/gate/engine.go:195`).
- The session records show complete checkpoints that ran again with no change between them, in the same worktree.

Inference: a retry on an unchanged composition stays in one checkout, so today's key already reuses it. The grade key adds no retry saving.

## Saving

| case | landings in 19 days | saved whole gates | saved time |
|---|---|---|---|
| upper bound: every final checkpoint answers its landing | 30 | 30 | about 144 min, or 53 min per week |
| current selection: only sessions with no worktree build | about 8 (7 of 27) | about 8 | about 38 min, or 14 min per week |
| retries | not counted | 0 | 0 |

## Contradictions and unknowns

- The handoff of 2026-10-26 counted 274 landings and 36 spec landings since 2026-10-01. This run counts 330 first-parent merges and 40 spec flips. The handoff does not state its method, so this asset keeps its own count.
- The session records do not show whether a `dist/bench` was present at the moment of each complete checkpoint. A worktree build earlier in the session is the proxy.
- Codex session records were not read for the order of calls.
- No byte comparison of two run binaries from one tree was made.

## Validation plan

- Rerun the git count with the method below, and compare the counts.
- Before the spec relies on equal run binaries, build one tree in two checkouts and compare the digests.
- After the build lands, count the reuse lines of the next ten spec landings.

## Method

- Spec landings: for each first-parent merge `M` since 2026-10-01, look for `-Status: staged` and `+Status: implemented` in `git diff M^1 M -- 'specs/*/spec.md'`. Then run `git merge-base --is-ancestor M^1 M^2`.
- Order of calls: in each session record, list the Bash calls in order.
  Keep each call that runs a complete checkpoint, a commit, a record write, a merge, or a landing.
- Wall time: read the `gate.finish` event of each `.logs/gate-*.jsonl` record in the primary checkout.

## Verification record

- [x] Opens with the recommendation, the scope, and the evidence status.
- [x] Keeps facts and inferences apart.
- [x] One section for each question.
- [x] A table where comparison matters.
- [x] Keeps every contradiction and unknown.
- [ ] Diagram: not drawn, because the relations are linear.
- [x] States what could not be verified.
- [x] Ends with a validation plan.
- [x] Each load-bearing claim cites a path and line, or names its method.
