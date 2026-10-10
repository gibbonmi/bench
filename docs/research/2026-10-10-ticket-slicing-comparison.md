# Ticket slicing comparison: whole graph against chunk by chunk

Recommendation: keep the per-chunk research, but forbid the slicer to drop a path that research names. The slicer can drop such a path only when an acceptance line binds the ticket author to the design that avoids it. Give the reviewer the research table and the drop table. Run the research for all chunks in one round, and review the tickets in one round on the conditional line. A replay on the saved research decides the slicing seat.

On FT349, the NEW method cost 2.6 times the wall time of the OLD method and 4.8 times its uncached input. Its recall was lower (0.738 against 0.787), but its research alone reached 0.869. The main session lost the difference: it dropped 10 real paths, each by a design choice that no ticket line bound.

Scope: the ticket-slicing step of `/bench-write-spec`, on one already-built spec
Evidence status: measured on one subject with one run per arm; the causes are inferences
Source repository: /home/mgibs/workspace/bench
Source branch: main
Source commit: e406423f9d58ba8dcaf6a6541303f8ea9afb0e17
Retrieved: 2026-10-10
Consumed by: the reviewer's decision about the slicing rules in `/bench-write-spec` and `craft-tickets`
Drift: refresh when the slicing step, the slicing checks, or the review lines of `craft-line` change
Retire when: an approved slicing rule replaces this comparison and a second subject validates it

This report changes no command, skill, or check. It proposes rule changes for the reviewer to decide.

## Question graph

1. Which subject and which commits make a fair test?
2. How did each arm slice and review?
3. What did each arm cost, and how many review rounds did it take?
4. How well did each arm predict the real build fence?
5. Which real paths did each arm miss, and why?
6. Which workflow should Bench use?

Questions 1 and 2 fix the method. Questions 3 to 5 hold the facts. Question 6 holds the proposal.

## 1. Subject and commits

### Facts

The subject is `ticket-checkpoint-floor` (FT349). It meets each selection rule:

- It is retired, and its build landed at `d0b62266`.
- It has 5 tickets in 5 chunks (TF-C1 to TF-C5).
- Its build widened the fence many times. Plan commit `26639d1d` widened ticket 04 by four owner files, after its author reported four copies of a rule outside the fence. Plan commits `938f04f9`, `5b58d27c`, `7c308b1b`, `64c8dd77`, and `7545640e` widened tickets 01 to 03.

BASE is `bc389247`. That commit staged the spec and its research asset, and no ticket file existed yet. The tickets came next, at `561da614`.

The ground truth is the set of files that the build changed: `git diff --name-only 9b99dc0c d0b62266`. The diff starts at the landing's first parent, not at BASE, because the build branch merged main twice. A diff from BASE carried unrelated main files. The truth set leaves out `specs/`, `reviews/`, `roadmap/`, `ROADMAP.md`, and `.bench/commitment.json`. It holds 61 paths.

## 2. Method of each arm

### Shared setup

- Each arm worked in its own Bench worktree (`slice-old`, `slice-new`) on the new slug `ticket-checkpoint-floor-exp`. The worktree held the spec at BASE and the arm's own copy of the slicing guidance.
- The worktrees stayed at main. The `block-dangerous-git` guard refused a restore of the tree to BASE, and the reviewer chose a source snapshot instead. Each arm read source only from a `git archive bc389247` export.
- The main tree holds the built code, so both arms were forbidden to run `--propose-writes`, `bench consumers`, `bench anchors`, and `bench preflight`. Each arm did the closure checks by hand.
- The spec body was frozen in both arms. A slicer could change only the tickets and the fence list. The NEW spec-alone review did not run, because the spec was approved at BASE.
- The arms ran one after the other, OLD first. Neither arm saw the other's tickets.

### OLD (command text at `4cb8d010`)

1. One slicer in the spec-fork seat (opus, high) sliced all five tickets and wrote the fence.
2. One review round (opus, high: the old conditional review line, which is the Claude mid binding) graded the spec and the tickets together.
3. The slicer folded the findings, and the same reviewer graded the repair.

### NEW (command text at `e406423f`)

For each chunk in order:

1. Two research readers (opus, medium) returned the likely fence expansions in the six classes of the slicing checks.
2. One long-lived main session (opus, high) verified the findings, wrote the chunk's ticket, and folded the last chunk review.
3. One chunk reviewer (sonnet, high: the new cheap chunk review line) returned accept or reject. The research of the next chunk ran during this review.

After the last chunk, the main session wrote the fence.

The coordinator ran `bench preflight build ticket-checkpoint-floor-exp --plan-only` after each arm and relayed nothing.

## 3. Cost and review rounds

### Tested results

| Measure | OLD | NEW |
| --- | --- | --- |
| Wall time, start to accepted ticket set | 1577 s | 4029 s |
| Delegates | 2 | 16 (1 main session, 10 research readers, 5 chunk reviewers) |
| Output tokens | 24,594 | 40,963 |
| Uncached input and cache writes | 735,175 | 3,492,396 |
| Cache reads | 30.6 M | 108.7 M |
| Review rounds and iterations | 1 round, 2 iterations | 5 rounds, 5 iterations |
| Blocking findings | 3 in iteration 1 (1 against the frozen spec) | 0 |
| All findings | 8, then 3 | 31 |
| `--plan-only` reds | 2 | 2 |

The token counts sum the final usage of each assistant message in each delegate transcript. The harness also reports a context size at each stop. Its sum is 488,460 for OLD and 2,431,232 for NEW.

Both `--plan-only` runs show the same two reds:

- `check-floor`: the frozen TF-C3 completion plan omits `./internal/conformance`. The OLD reviewer found this defect as a blocking finding. The real build fixed it with plan commit `7545640e`.
- `anchor-closure`: it names `internal/anchors/registry_check_floor_test.go`, a file that the build added after BASE. This red is drift from the main tree.

The coordinator relayed neither red, because neither red can be repaired under the experiment's rules.

## 4. Fence prediction

### Tested results

| Ticket set | Recall | Precision | Precision without forced closure |
| --- | --- | --- | --- |
| OLD final | 0.787 (48 of 61) | 0.581 (36 of 62) | 0.837 |
| NEW final | 0.738 (45 of 61) | 0.579 (33 of 57) | 0.868 |
| Real approved tickets at `f3301951` | 0.803 (49 of 61) | 0.607 (37 of 61) | 0.881 |
| NEW research union (counterfactual) | 0.869 (53 of 61) | 0.500 (42 of 84) | - |
| NEW research, entries both readers named | 0.787 | 0.561 | - |
| NEW final plus the research union | 0.902 | 0.500 | - |

Recall is the share of truth paths that a `Writes:` entry covers. A directory entry covers each file below it. Precision is the share of entries that hit a truth path.

All three ticket sets name the same 19 forced-closure entries: the five command-registry files and fourteen anchor-registry files. The enforced closure rules require these entries, and the real build changed none of them. The last column leaves them out.

The real approved tickets came from the OLD-style method on 2026-10-08. They are a reference point, not an arm.

### Research yield

- The ten readers named 84 distinct entries, and 42 of them hit a truth path.
- The research union covers 53 truth paths. The NEW main session dropped 10 of them from its final tickets.
- The main session added 2 truth paths that no reader named: `internal/anchors/registry_check_floor.go` and `internal/reviewrecord/recordcmd/chunk.go`. Both are files that it split out of over-budget files.
- Six truth paths sit outside the provisional fence. The research named 4 of them. The NEW tickets kept 1 (`cmd/bench/preflight_version_test.go`), and the OLD tickets covered the same 1.
- The four owner files of plan commit `26639d1d` are `internal/coverage/citations.go`, `internal/reviewrecord/coverage.go`, `internal/reviewrecord/completion.go`, and `internal/reviewrecord/recordtest/fixture.go`. The research named 3 of them. Neither arm's final tickets named any of them.

## 5. Missed paths

### Tested results

| Truth path | Class | OLD | NEW | Research |
| --- | --- | --- | --- | --- |
| `internal/anchors/registry_check_floor.go` (new) | anchor registry | miss | hit | miss |
| `internal/anchors/registry_check_floor_test.go` (new) | anchor registry | miss | miss | miss |
| `internal/gate/gate.go` | caller file | miss | miss | miss |
| `internal/preflight/source_tip_test.go` | count pin | miss | miss | miss |
| `internal/preflight/ticket_rows.go` (new) | over-budget file | miss | miss | miss |
| `internal/reviewrecord/recordcmd/chunk.go` (new) | over-budget file | miss | hit | miss |
| `internal/reviewrecord/recordcmd/planned_cell.go` (new) | over-budget file | miss | miss | hit |
| `internal/coverage/citations.go` | caller file | miss | miss | hit |
| `internal/reviewrecord/coverage.go` | caller file | miss | miss | hit |
| `internal/reviewrecord/completion.go` | caller file | miss | miss | miss |
| `internal/reviewrecord/recordtest/fixture.go` | fixture or canary | miss | miss | hit |
| `internal/testreport/outcome_test.go` | other | miss | miss | miss |
| `internal/testreport/selection_facts_test.go` | other | miss | hit | hit |
| `internal/preflight/preflighttest/testdata/review-*.toon` (4 files) | fixture or canary | hit | miss | hit |
| `internal/testreport/check_test.go` | over-budget file | hit | miss | hit |
| `internal/testreport/environment.go` | other | hit | miss | hit |

| Class | OLD misses | NEW misses |
| --- | --- | --- |
| Caller file | 4 | 4 |
| Fixture or canary | 1 | 5 |
| Anchor registry | 2 | 1 |
| Count pin | 1 | 1 |
| Over-budget file | 3 | 3 |
| Shared conflict file | 0 | 0 |
| Other | 2 | 2 |
| Total | 13 | 16 |

The caller-file class here holds the owners of a rule. The real build folded duplicated copies of a rule into those owners, or it exported a rule for a second caller. The count pin is a row count in `source_tip_test.go` (14 became 15).

### Inference: why the NEW main session dropped real paths

The main session gave a reason for each drop. The reasons share one pattern. A path was needed only under one design, so the slicer chose the other design and removed the path. No acceptance line bound the ticket author to that design.

| Dropped path | Stated reason | What the real build did |
| --- | --- | --- |
| Four `review-*.toon` goldens | `PlanFence` keeps its bytes | It changed the `PlanFence` command, so every pinned hash changed |
| `environment.go` | The caller sets the empty scope | It added a `conformanceScope` helper in this file |
| `check_test.go` | The file is over budget and must not grow | It edited the file, net -15 lines |
| `citations.go`, `reviewrecord/coverage.go`, `recordtest/fixture.go` | The file is only read, or an export elsewhere avoids it | It folded duplicated copies into these owners |

Two facts made the drops stick:

- The step "verify the findings" lets the slicer prune by design. The ticket author is a fresh session, and it can choose a different design.
- The chunk reviewer saw the tickets, but it did not see the research returns or the slicer's drop list. The TF-C3 reviewer approved the "keep `PlanFence`" design.

The OLD slicer kept more of its design-conditional paths. Its opus reviewer also found the goldens gap and the `command.go` gap as blocking findings. The sonnet chunk reviewers found no blocking defect in five rounds.

### Inference: context growth

| Stop | NEW main-session context |
| --- | --- |
| After TF-C1 | 140k |
| After TF-C2 | 186k |
| After TF-C3 | 269k |
| After TF-C4 | 316k |
| After TF-C5 | 371k |
| After the fence | 388k |

The OLD slicer ended at 269k. Each chunk added 45k to 85k tokens: two research returns, one review report, and source reads. Only a small state must cross a chunk boundary. The ledger in section 6 lists it.

## 6. Recommended workflow

### Proposal

Two read-only consultants answered this question independently on the same evidence pack. One ran on fable at high effort, and one ran on sonnet at high effort. The coordinator reconciled the two returns, checked the disputed facts, and sent the draft back to each consultant. Both consultants signed the result below. Each gave a confidence of 6 out of 10.

Both consultants found that the drop rule is the lever, not the seat. Research found the real paths, and the slicer removed them. A fresh slicer would make the same design choices without a rule against them.

1. **Spec stage.** The spec fork writes the spec and its chunk table. The spec-alone review runs on the cheap line at high effort. This stage does not change, because the experiment did not test it.
2. **Ticket research.** After the spec review accepts, one round charges two readers for each chunk on the mid line at medium effort. All chunks run at the same time. The coordinator compares each pair of returns and writes one research table: the agreed set and the frontier. If the spec review changes the chunk table, the research reruns for each changed chunk.
3. **Slicing.** The current main-session seat stays until a replay decides the seat. The slicer obeys the keep rule below. For each ticket, it writes a drop table with three columns: the research path that it removed, the reason, and the binding acceptance line.
4. **Ticket review.** One round on the conditional line (mid at high effort), with a cap of 2 iterations. The reviewer receives the spec, the tickets, the research table, and the drop table. It grades the `craft-tickets` quiz and each drop. If the ticket set does not fit one review context, the round splits at chunk boundaries. This round replaces the per-chunk review on the cheap line.
5. **Fence and plan-only.** The slicer writes the fence, and `bench preflight build <slug> --plan-only` grades the whole graph before sign-off.

### Proposed rule changes

- **Keep rule** (`slicing-checks.md`): a path that a research return names stays in the ticket's `Writes:` line. The slicer can drop it only when an acceptance line binds the author to the design that avoids it. The drop table names that line.
- **Review inputs** (the command's review step): the ticket review charge carries the research table and the drop table, and the reviewer grades each dropped path.
- **Over-budget files** (`slicing-checks.md`): an over-budget file that the ticket's behavior touches stays in `Writes:`, even when the ticket plans no growth. A ticket that splits such a file names each candidate home that research proposed.
- **Review line** (`craft-line`): the ticket review row moves from the cheap line to the conditional line. The spec review row stays on the cheap line.
- **Research timing** (the command's research step): the ticket research for all chunks runs as one round after the spec review.

These changes leave four misses open: `gate.go`, `completion.go`, `source_tip_test.go`, and `registry_check_floor_test.go`. Research missed them too. Two candidate rules belong to a separate kit idea:

- a seventh evidence class for an existing copy of a rule, with its owner file;
- a rule that a new registry file names its test file and the count pin that it moves.

The second-caller rule in `slicing-checks.md` already exists and was still missed, so that gap is one of enforcement, not wording.

### Inference: why each choice

- **Two readers, not one.** Reader A alone gave recall 0.803. Reader B alone gave 0.852. The pair gave 0.869. That spread between two readers on the same chunks is the variance that the pair absorbs. `craft-research` already requires the pair.
- **One review round on the conditional line.** The one opus round cost 218,117 cache-write tokens and 604 s, and it found 3 blocking defects. The five sonnet rounds cost 758,802 cache-write tokens and 1,738 s, and they found none. They also approved the drop that lost the four goldens. The tier and the method changed together, so the finding power is an inference. The cost difference is a measurement.
- **Seat left open.** The drops came from design choices that a fresh slicer would also make without the keep rule. A fresh slicer avoids the context growth, but no run has tested it.

### Context management

If the main-session seat stays, compact or start a new session at a chunk-accept boundary when the context passes about 300 k tokens. In this run that point is after TF-C3. The threshold is a judgment, not a measurement. Never compact inside a chunk. A fresh session needs this carry-over ledger:

- the spec path, the chunk table, and the accepted chunks with their ticket paths;
- the research table (the agreed set and the frontier);
- each ticket's drop table;
- the cross-chunk deferrals, such as "`CHANGELOG.md` goes to ticket 05";
- the ordering edges and the shared-file owners;
- the flags against the frozen spec, and the paths that wait for the fence.

Today no file holds the research table. `craft-research` requires one durable output for each research run, and neither command version names its place. That place is a separate decision for the reviewer.

If a fresh slicer takes the seat, the only restart point is after the spec review. The slicer and the reviewer then start with fresh contexts.

### Expected cost

The consultants estimate about 2,000 to 2,700 s of wall time and about 2.2 to 2.6 times the OLD token cost. That is about half of NEW. They expect recall between 0.84 and 0.90, with raw precision near 0.52 to 0.55. These figures come from this run's numbers and are not measured.

## Threats to validity

- **One subject.** One spec, with one run per arm, cannot separate a method effect from model variance.
- **Kit knowledge at main.** Both arms ran the kit binary and the guidance at main. Only the source tree was at BASE.
- **The real fence is not minimal.** The truth set is one build's design choices. A design-conditional miss can be a valid alternative design.
- **No `--propose-writes`.** Both arms did the closure checks by hand. The command would have closed some closure classes mechanically.
- **Two variables changed together.** The method changed, and the review tier changed with it (opus in OLD, sonnet in NEW).
- **Coordinator latency.** Wall time includes the coordinator's relays, and NEW had more relays.

## What was read and what was not

The coordinator read these sources:

- the OLD and NEW command text, and `craft-tickets` with its slicing checks at both commits;
- `craft-line`, `craft-research`, and `craft-delegate`;
- the spec and its research asset at BASE;
- the real approved tickets at `f3301951` (the coordinator only, never an arm);
- the build diff and its commit subjects;
- every delegate return, the token use in each transcript, and the `internal/assessment` schema.

The two consultants read the evidence pack, both command versions, and the current skills. They also read some of the NEW tickets, research returns, and review reports. The fable consultant did not read reviews C1, C2, and C5 or tickets 01, 02, and 05.

The coordinator did not read these sources:

- the real build's review record and its retro;
- any arm's ticket text beyond its `Writes:` line and its returned reasons;
- the full diffs of the build commits, except the files that section 5 names.

## Verification record

- [x] The output opens with the recommendation, the scope, and the evidence status.
- [x] It keeps facts, inferences, and proposals in separate sections.
- [x] It holds one section for each question.
- [x] It uses tables where the comparison matters.
- [x] It keeps each contradiction and each residual unknown.
- [x] It states what could not be verified.
- [x] It ends with a validation plan.
- [x] Each load-bearing claim names a commit, a path, or a recorded measurement.

The two arms are recorded as the assessment runs `ticket-slicing-old-ft349-20261010` and `ticket-slicing-new-ft349-20261010`. `bench assessment compare` refused the pair. The verb pins one line for each attempt role and allows one declared variable. These arms differ in the method and also in the review model that each method names.

## Validation plan

1. **Replay the slicing seat.** Freeze this run's spec and its ten saved research returns. Run a fresh slicer against the main session, with the keep rule, the review-inputs rule, and a reviewer on the mid line. Grade both arms against the same 61-path truth set. Each arm costs about 25 minutes.

   Record whether any defect that the single review round finds comes from a design choice in TF-C1 or TF-C2. If one does, restore a per-chunk review for the chunks that create a seam.
2. **Test the reviewer.** In the replay, give the mid reviewer the TF-C3 drop of the four goldens. If it still accepts the drop, a review cannot police design-conditional drops. In that case, the agreed research set goes into `Writes:` by rule.
3. **Repeat on a second subject.** Replay the chosen workflow on a second retired spec, such as `refusal-route-registry` (FT393), with `--propose-writes` available. Score recall, precision, and the count of plan-expansion commits during the build.
4. **Measure the precision cost.** Count the build-time effect of extra `Writes:` entries: a wider check floor and a looser fence. If that cost is material, narrow the keep rule to the classes that research predicts well.
