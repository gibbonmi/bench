# Workflow guidance repair review record

## GR-A author evidence

Each ticket had a fresh `bench-writer` author on opus at high effort, with a cap of 3 attempts. Each author added its registry rows before the guidance edit. Then the author ran `bench anchors` on the unchanged guidance, and each new row went red. After the edit, each row went green.

| Ticket | Author session | Start tip | Ticket commit | Attempts |
|---|---|---|---|---|
| 1 | `claude:bench-writer/gr-t1-author` | `8af51e9b` | `d7e74518` | 2 of 3 |
| 2 | `claude:bench-writer/gr-t2-author` | `6e0c8203` | `eb1cb1c0` | 1 of 3 |
| 3 | `claude:bench-writer/gr-t3-author` | `1ec3325b` | `581ff5da` | 1 of 3 |

The first lane run of ticket 1 failed on `structure`, because `registry_data.go` is over its budget. The orchestrator amended the Testing decisions at `6e0c8203`. A new row now goes in `registry_retained_workflow.go`, and an in-place replacement stays in `registry_data.go`.

The first ticket 3 run stopped before any edit. The helper `checkReviewConvergenceContract` required three sentences that the ticket removes. The orchestrator added `docs_workflow_helpers_test.go` to the ticket 3 fence at `1ec3325b`. The same author then removed those three requirements. Two `bench learning` entries record these plan expansions.

The plan commit `a55c1749` records the assignments of tickets 4 to 8 before the chunk freeze. So the plan digest stays the same until a repair assignment changes it. Each of those assignments names `581ff5da` as its source.

### Probe verdicts

Each probe ran through `bench probe --check docs-currency-workflow`. Each probe bit, and each restore reads `yes`.

| Ticket | File | Mutation | Row |
|---|---|---|---|
| 1 | `bench-craft-delegate/SKILL.md` | swap: the canary plants the GR16 sentence | GR16 |
| 1 | `bench-craft-delegate/SKILL.md` | swap: the canary drops the coverage-row contents | GR6 |
| 1 | `.claude/agents/bench-writer.md` | omission: the lane-pass commit sentence | GR109 |
| 1 | `bench-craft-delegate/SKILL.md` | swap: the charge example stops at diff ready | GR113 |
| 1 | `bench-craft-delegate/SKILL.md` | omission: the guard deny-surface clause | GR14 |
| 2 | `delegation-discipline.md` | swap: the `delegate-probe-mutated-bytes` canary | probe-verdict row |
| 2 | `delegation-discipline.md` | swap: each of the four transfer-trigger canaries | trigger rows |
| 2 | `delegation-discipline.md` | omission: the GR17 bullet | GR17 |
| 2 | `delegation-discipline.md` | swap: the rebase premise returns | GR27 |
| 2 | `bounded-repair-policy.md` | swap: the retired mode list returns | GR31 |
| 3 | `bench-review-implementation.md` | swap: the `review-repair-ticket-owner` canary | GR33 |
| 3 | `bench-review-implementation.md` | swap: the `review-repair-ticket-covers` canary | GR34 |
| 3 | `.bench/BENCH.md` | swap: the `delegated-chunk-tip-review` canary | GR116 |
| 3 | `bench-debug.md` | swap: the coordinator reslice returns | GR38 |
| 3 | `bench-final-check.md` | omission: the GR37 sentence | GR37 |
| 3 | `bench-implement-spec.md` | swap: the old tier ladder returns | GR40 |

### Verification

Each author ran the six GR-A checks at the chunk tip `a55c1749`, and each check passed. The JSON payload holds each result.

## GR-A chunk review, round 1

The frozen pair is base `8a107bf5c6d005cd6f730baf24dc565d97d865ae` and tip `a55c1749c214b2c5332dc530bbeae113e5dafd34`. The shared evidence is `sha256:361e66c28bc54343200722b0bcde99febc8e50bc1104f075edd7fa7a175d4b35`. Each axis ran in a fresh `bench-reviewer` session on opus at high effort. Only the Coverage axis ran probes, and it left the tree clean.

The consumer table has one row outside the diff: `docs_workflow_checks_test.go:30` calls `checkRetainedWorkflow`. That function keeps its signature, and the `docs-currency-workflow` check runs it.

The raw finding count is 15: Standards 8, Spec 3, and Coverage 4. After the orchestrator merges the findings that name the same fix, 8 repair targets remain. The reviewer approved all work for this build, so the orchestrator decided each `ask-user` finding. These decisions are open to reviewer veto.

## Standards

Findings: 8. The worst issue is an anchor needle across two physical lines.

- `.agents/skills/bench-craft-delegate/SKILL.md:54-55` splits the GR6 needle across a line break. Target R1. `auto-fix`. Confidence 7.
- `internal/anchors/registry_retained_workflow.go:21-25` has a doc comment that does not describe the rows that moved from `registry_data.go`. Target R3 fixes the comment. The prefixes stay, because canary files pin the diagnostics. `auto-fix`. Confidence 6.
- `.claude/agents/bench-writer.md:9,20` says "repair author", but line 3 says "repair session". The GR109 needle and the spec use both terms. `no-op` by orchestrator decision. Confidence 5.
- `.agents/skills/bench-craft-delegate/SKILL.md:97-98` states the share rule twice. Ticket 1 requires both needles. `no-op` by orchestrator decision. Confidence 5.
- `.agents/skills/bench-craft-delegate/references/delegation-discipline.md:100` holds a third copy of the tier-move stop. GR17 requires the sentence. `no-op` by orchestrator decision. Confidence 5.
- `.agents/commands/bench-implement-spec.md:60` keeps the final reconciliation step. That line is the build phase's own action step. `no-op` by orchestrator decision. Confidence 4.
- `SKILL.md:69,98` and `bench-writer.md:20-22` state the same delegate rules. The agent file is the delegate's boot surface. `no-op`. Confidence 3.
- The new test rows restate registry needles without a recorded red. The pattern predates this chunk. `no-op`. Confidence 3.

## Spec

Findings: 3. The worst issue is a stale routing claim in the delegation discipline.

- `.agents/skills/bench-craft-delegate/references/delegation-discipline.md:10` says that `craft-delegate` routes each spec-backed ticket. Ticket 1 removed that rule. Target R5. `auto-fix`. Confidence 7.
- `.agents/skills/bench-craft-delegate/SKILL.md:66` restates the share rule of lines 97-98. Target R2. `auto-fix`. Confidence 7.
- `.agents/skills/bench-craft-delegate/references/delegation-discipline.md:190-191` allows an umbrella repair ticket. No audit finding names that sentence, and the scope is closed. `no-op` by orchestrator decision; a `bench idea` parks it. Confidence 4.

## Coverage

Findings: 4. The worst issue is that the new Forbid rows catch only the retired wording.

- `.agents/commands/bench-review-implementation.md:46,59,72` and `.agents/commands/bench-implement-spec.md:60` accept a verbatim copy of the owner sentence from `.bench/BENCH.md`. Four swap probes were silent. Target R6 adds Forbid rows under the plan-expansion policy. `auto-fix`. Confidence 8.
- `.agents/commands/bench-review-implementation.md:44` has an owner pointer that no row pins. The omission probe was silent. Target R7. `auto-fix`. Confidence 7.
- `.agents/commands/bench-final-check.md:211` has the scope "For other work," that no row pins. The swap probe was silent. Target R8. `auto-fix`. Confidence 7.
- `.claude/agents/bench-writer.md:22` has a landing ban for a user-directed delegate that no row pins. Ticket 1 keeps the rule that no role lands. Target R4. `auto-fix` by orchestrator decision. Confidence 5.

## GR-A repair routing

Each repair goes to a fresh `bench-writer` repair session for the ticket whose `Writes:` line holds the path.

| Target | Ticket | Repair |
|---|---|---|
| R1 | 1 | Join the GR6 needle onto one line, and update the canary `old` bytes. |
| R2 | 1 | Remove the duplicate share sentence at `SKILL.md:66`. |
| R3 | 1 | Make the `delegatedWorkflowAnchors` doc comment name the moved rows. |
| R4 | 1 | Add a Require row for the landing ban of a user-directed delegate. |
| R5 | 2 | Point the routing claim to its owner. |
| R6 | 3 | Add Forbid rows for the owner sentences in the phase commands. |
| R7 | 3 | Add a Require row for the review phase owner pointer. |
| R8 | 3 | Pin the "For other work," scope in the final check. |

```bench-review-record
{
  "version": 2,
  "spec": "specs/workflow-guidance-repair/spec.md",
  "plan_digest": "sha256:fbef67ab845114298d98e72b74f9813faebb0c3872382294ca5b401cc9eec89b",
  "implementation_session": "",
  "chunks": [
    {
      "id": "GR-A",
      "base": "8a107bf5c6d005cd6f730baf24dc565d97d865ae",
      "tip": "a55c1749c214b2c5332dc530bbeae113e5dafd34",
      "plan_digest": "sha256:fbef67ab845114298d98e72b74f9813faebb0c3872382294ca5b401cc9eec89b",
      "source_digest": "56aa36e8e96230600f18450efa1aba35f2775cad",
      "acceptance_rows": [
        "GR1",
        "GR2",
        "GR3",
        "GR4",
        "GR5",
        "GR6",
        "GR7",
        "GR8",
        "GR9",
        "GR10",
        "GR11",
        "GR12",
        "GR13",
        "GR14",
        "GR15",
        "GR16",
        "GR109",
        "GR110",
        "GR111",
        "GR112",
        "GR113",
        "GR17",
        "GR18",
        "GR19",
        "GR20",
        "GR21",
        "GR22",
        "GR23",
        "GR24",
        "GR25",
        "GR26",
        "GR27",
        "GR28",
        "GR29",
        "GR30",
        "GR31",
        "GR32",
        "GR115",
        "GR33",
        "GR34",
        "GR35",
        "GR36",
        "GR37",
        "GR38",
        "GR39",
        "GR40",
        "GR41",
        "GR42",
        "GR43",
        "GR44",
        "GR116"
      ],
      "verification": [
        {
          "id": "gr-a-1-workflow",
          "performer": "claude:bench-writer/gr-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "56aa36e8e96230600f18450efa1aba35f2775cad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/gr-t1-author-20260924/1-workflow@a55c1749",
            "digest": "sha256:08d2cd55717ae9ad20ba0feda218a2c42f932e5a2027e21e789964a764826d92",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,1196\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "1-workflow",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "gr-a-1-conformance",
          "performer": "claude:bench-writer/gr-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "56aa36e8e96230600f18450efa1aba35f2775cad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/gr-t1-author-20260924/1-conformance@a55c1749",
            "digest": "sha256:b1f6ee073d06f0fd4458eaac0fc28ffe68d44062115f21eb3c05a502edaeaa1b",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,6255\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "1-conformance",
          "command": "bench test --package ./internal/conformance --run TestRootConformance",
          "exit_code": 0
        },
        {
          "id": "gr-a-1-anchors",
          "performer": "claude:bench-writer/gr-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "56aa36e8e96230600f18450efa1aba35f2775cad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/gr-t1-author-20260924/1-anchors@a55c1749",
            "digest": "sha256:09cb5eb796932d882cc65a75e4da256211b85ad410788e6b1704a7da3d71c5dc",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/anchors,pass,745\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "1-anchors",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "gr-a-1-agents",
          "performer": "claude:bench-writer/gr-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "56aa36e8e96230600f18450efa1aba35f2775cad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/gr-t1-author-20260924/1-agents@a55c1749",
            "digest": "sha256:6f72c99a424331a3381e284e55ffa703b490f81e19b498e848e8bfc5d9ca25ed",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,5\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "1-agents",
          "command": "bench test --check claude-agent-definitions",
          "exit_code": 0
        },
        {
          "id": "gr-a-1-budgets",
          "performer": "claude:bench-writer/gr-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "56aa36e8e96230600f18450efa1aba35f2775cad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/gr-t1-author-20260924/1-budgets@a55c1749",
            "digest": "sha256:c224ae6aff6b0057319e254320ee168971646da0b2de6eae4c92751c0643cf54",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,7\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "1-budgets",
          "command": "bench test --check guidance-prose-budgets",
          "exit_code": 0
        },
        {
          "id": "gr-a-1-fixture-bite",
          "performer": "claude:bench-writer/gr-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "56aa36e8e96230600f18450efa1aba35f2775cad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/gr-t1-author-20260924/1-fixture-bite@a55c1749",
            "digest": "sha256:39517813005dfd2ec034d7ee16030aff2d5b6c130af2e0132aa194f5565451c2",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,1088\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "1-fixture-bite",
          "command": "bench test --package ./internal/conformance --run 'TestWorkflowCadenceAnchorsRejectDeletionAndSwap|TestSpecTicketHandoffWorkflowFixturesAreComplete'",
          "exit_code": 0
        },
        {
          "id": "gr-a-2-workflow",
          "performer": "claude:bench-writer/gr-t2-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "56aa36e8e96230600f18450efa1aba35f2775cad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/gr-t2-author-20260924/2-workflow@a55c1749",
            "digest": "sha256:9bac0220cfd9dc6c5ba880c37e219bbf2deb2f632dcbab662f471211e2f45c61",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,844\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "2-workflow",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "gr-a-2-conformance",
          "performer": "claude:bench-writer/gr-t2-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "56aa36e8e96230600f18450efa1aba35f2775cad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/gr-t2-author-20260924/2-conformance@a55c1749",
            "digest": "sha256:faecbd0445e487b8737da0572605e6fdb665521141d76350d0269f0186816d6f",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,6455\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "2-conformance",
          "command": "bench test --package ./internal/conformance --run TestRootConformance",
          "exit_code": 0
        },
        {
          "id": "gr-a-2-anchors",
          "performer": "claude:bench-writer/gr-t2-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "56aa36e8e96230600f18450efa1aba35f2775cad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/gr-t2-author-20260924/2-anchors@a55c1749",
            "digest": "sha256:881bfbf9b55a1e0b857da719aab3c0eaa1b3b8f654395f1224953340efe6f589",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/anchors,pass,760\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "2-anchors",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "gr-a-2-agents",
          "performer": "claude:bench-writer/gr-t2-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "56aa36e8e96230600f18450efa1aba35f2775cad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/gr-t2-author-20260924/2-agents@a55c1749",
            "digest": "sha256:6f72c99a424331a3381e284e55ffa703b490f81e19b498e848e8bfc5d9ca25ed",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,5\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "2-agents",
          "command": "bench test --check claude-agent-definitions",
          "exit_code": 0
        },
        {
          "id": "gr-a-2-budgets",
          "performer": "claude:bench-writer/gr-t2-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "56aa36e8e96230600f18450efa1aba35f2775cad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/gr-t2-author-20260924/2-budgets@a55c1749",
            "digest": "sha256:6f72c99a424331a3381e284e55ffa703b490f81e19b498e848e8bfc5d9ca25ed",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,5\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "2-budgets",
          "command": "bench test --check guidance-prose-budgets",
          "exit_code": 0
        },
        {
          "id": "gr-a-2-fixture-bite",
          "performer": "claude:bench-writer/gr-t2-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "56aa36e8e96230600f18450efa1aba35f2775cad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/gr-t2-author-20260924/2-fixture-bite@a55c1749",
            "digest": "sha256:fb12ef14fff4a050f53d4e2e3c38f1f0564e188c7abd27d6d314294fa1783c3a",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,1145\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "2-fixture-bite",
          "command": "bench test --package ./internal/conformance --run 'TestWorkflowCadenceAnchorsRejectDeletionAndSwap|TestSpecTicketHandoffWorkflowFixturesAreComplete'",
          "exit_code": 0
        },
        {
          "id": "gr-a-3-workflow",
          "performer": "claude:bench-writer/gr-t3-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "56aa36e8e96230600f18450efa1aba35f2775cad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/gr-t3-author-20260924/3-workflow@a55c1749",
            "digest": "sha256:f9eab972a608f17d9413c624e575cc404340f42ee1da998ed245ffa0dcd8db7a",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,721\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "3-workflow",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "gr-a-3-conformance",
          "performer": "claude:bench-writer/gr-t3-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "56aa36e8e96230600f18450efa1aba35f2775cad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/gr-t3-author-20260924/3-conformance@a55c1749",
            "digest": "sha256:01d7bb29ec0dd0fb588c080ba43187c120897cc99818427e94a1c1e55a4f975a",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,6932\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "3-conformance",
          "command": "bench test --package ./internal/conformance --run TestRootConformance",
          "exit_code": 0
        },
        {
          "id": "gr-a-3-anchors",
          "performer": "claude:bench-writer/gr-t3-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "56aa36e8e96230600f18450efa1aba35f2775cad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/gr-t3-author-20260924/3-anchors@a55c1749",
            "digest": "sha256:569bc56dc9ec4b686b652be22bbd49811d4ee87be772eb19e1aab700b00bc3cd",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/anchors,pass,705\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "3-anchors",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "gr-a-3-agents",
          "performer": "claude:bench-writer/gr-t3-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "56aa36e8e96230600f18450efa1aba35f2775cad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/gr-t3-author-20260924/3-agents@a55c1749",
            "digest": "sha256:03ca50a13c5e89336b4b24e5071430a4d4a62c890ebe7e010807480de7aaf0ce",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "3-agents",
          "command": "bench test --check claude-agent-definitions",
          "exit_code": 0
        },
        {
          "id": "gr-a-3-budgets",
          "performer": "claude:bench-writer/gr-t3-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "56aa36e8e96230600f18450efa1aba35f2775cad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/gr-t3-author-20260924/3-budgets@a55c1749",
            "digest": "sha256:03ca50a13c5e89336b4b24e5071430a4d4a62c890ebe7e010807480de7aaf0ce",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "3-budgets",
          "command": "bench test --check guidance-prose-budgets",
          "exit_code": 0
        },
        {
          "id": "gr-a-3-fixture-bite",
          "performer": "claude:bench-writer/gr-t3-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "56aa36e8e96230600f18450efa1aba35f2775cad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/gr-t3-author-20260924/3-fixture-bite@a55c1749",
            "digest": "sha256:9b6e17ca3068630b67cb386f1905c8a62d6ffc98ce986e789720d9426a85e206",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,1203\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "3-fixture-bite",
          "command": "bench test --package ./internal/conformance --run 'TestWorkflowCadenceAnchorsRejectDeletionAndSwap|TestSpecTicketHandoffWorkflowFixturesAreComplete'",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "gr-a-r1-standards",
          "performer": "claude:bench-reviewer/gr-a-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "56aa36e8e96230600f18450efa1aba35f2775cad",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/gr-a-standards@a55c1749",
            "digest": "sha256:af4ee3fba2b6a7ad66352d0644a2d1a10de1cb09ac310858bb9b6fae138436db",
            "excerpt": "Standards: 8 findings. Worst: an anchor needle is split across two physical lines in an edited hunk, which breaks a hard ste-prose rule."
          },
          "axis": "Standards",
          "base": "8a107bf5c6d005cd6f730baf24dc565d97d865ae",
          "tip": "a55c1749c214b2c5332dc530bbeae113e5dafd34",
          "finding_ids": [
            "R1",
            "R3"
          ],
          "supersedes": []
        },
        {
          "id": "gr-a-r1-spec",
          "performer": "claude:bench-reviewer/gr-a-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "56aa36e8e96230600f18450efa1aba35f2775cad",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/gr-a-spec@a55c1749",
            "digest": "sha256:425a2bf066a3568c04af3261ac250cffbe3979a1e1a7902c2264fa0f42dd7deb",
            "excerpt": "Spec: 3 findings. Worst: the delegation discipline still says `craft-delegate` routes each spec-backed ticket to its author, but ticket 1 removed that rule from `craft-delegate`."
          },
          "axis": "Spec",
          "base": "8a107bf5c6d005cd6f730baf24dc565d97d865ae",
          "tip": "a55c1749c214b2c5332dc530bbeae113e5dafd34",
          "finding_ids": [
            "R2",
            "R5"
          ],
          "supersedes": []
        },
        {
          "id": "gr-a-r1-coverage",
          "performer": "claude:bench-reviewer/gr-a-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "56aa36e8e96230600f18450efa1aba35f2775cad",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/gr-a-coverage@a55c1749",
            "digest": "sha256:e8e4e1a4bf00327c36495e5548fec22fa32fcb1fee6321305a57a2853f97d521",
            "excerpt": "Coverage: 4 findings. Worst: the new Forbid rows guard only the retired wording, so a word-for-word copy of the `.bench/BENCH.md` rule pasted back into a phase command still passes the gate."
          },
          "axis": "Coverage",
          "base": "8a107bf5c6d005cd6f730baf24dc565d97d865ae",
          "tip": "a55c1749c214b2c5332dc530bbeae113e5dafd34",
          "finding_ids": [
            "R4",
            "R6",
            "R7",
            "R8"
          ],
          "supersedes": []
        }
      ]
    }
  ],
  "completion": {
    "state": "pending"
  }
}
```
