# Jev skill-selection benchmark trial

This report preserves the 2026-09-26 trial state and its original authority limits.
The [topic handoff](../session-handoff.md) owns the current continuation; later cycle reports supersede this report for repair status.

## Historical recommendation

Keep the benchmark trial stopped until the reviewer decides on a sixth repair cycle.
Two paid runs stopped on benchmark defects, not on Jev predictions.
Neither run supports a performance, latency, quality, or adoption conclusion.

Scope: the benchmark harness, its two live runs, and the repair boundary that remains.
Evidence status: the harness passed its offline suite and three review axes. Both live runs are diagnostic only.
Read date: 2026-09-26
Source branch: `bench/assign/fd0ac724000cdf7d2317bb953c7f64b3/9ad9d78ef1df680bcac5486c89051599`
Source tip: `c929c500501bc86606a0dae55923319a34e685d0`
Recorded base: `e667581749f23b9d12243ef4c844b6694abf25c6`
Roadmap owner: FT347
Consumed by: the FT347 decision, and the replay probe of this map's pilot
Drift: refresh after a harness change, a new live run, or a change to the Codex sandbox contract.
Retire when: the reviewer authorizes a repair and its successor report replaces this asset, or the reviewer drops the trial.

The harness code lives only on the source branch.
Read a harness file with `bench worktree show` while its assignment exists.
After the assignment is released, read the same blob with `git show <tip>:<path>`.
This asset cites a branch file as `c929c500:<path>:<line>`.

## Question graph

1. What does the benchmark measure, and under which evidence contract?
2. What did the two live runs observe?
3. Which defects stopped the runs?
4. What must a repair prove before another paid run?

Questions 2 and 3 depend on question 1. Question 4 depends on question 3.

## Q1. What does the benchmark measure?

### Facts

The benchmark asks one question.
Does Jev-assisted initial skill selection lower the complete task cost without a loss of task quality?
It adopts no production skill policy.
Source: `c929c500:decisions/jev-advisor/assets/benchmark/README.md:5`.

Both arms receive the same frozen packet and separate workspaces from one snapshot.
Both arms include source inspection, verification, self-review, and repairs.
The runner grades each submission with an independent final gate in a separate copy.
Source: `c929c500:decisions/jev-advisor/assets/benchmark/README.md:37`.

Each candidate skill receives two independent Noul questions.
The first question asks whether the action requires the skill.
The second question asks whether the skill is useful when it is not required.
Required uncertainty triggers native fallback. Optional uncertainty alone does not.
Source: `c929c500:decisions/jev-advisor/assets/benchmark/README.md:62`.

The corpus holds eight source-derived development tasks.
These tasks are exposed development material, not an independent holdout.
Source: `c929c500:decisions/jev-advisor/assets/benchmark/README.md:85`.

The live command needs the SHA-256 digest of the reviewed configuration.
It also needs a committed opt-in file in the frozen source.
Source: `c929c500:decisions/jev-advisor/assets/benchmark/README.md:128`.

### Branch content

| Commit | Change |
| --- | --- |
| `d3310333f25d029e16fe1863b52249c11afbf76f` | Rebuild the benchmark around complete and sealed evidence. |
| `7cf4105dc220638e32a7f07fd921a2bc669515fb` | Commit the hosted-use opt-in that the user authorized. |
| `c929c500501bc86606a0dae55923319a34e685d0` | Preserve replay fidelity and cancellation precedence. |

The branch adds the harness under `decisions/jev-advisor/assets/benchmark/`.
It also adds the retrospective `capture/retros/jev-benchmark-repair.md` and a scorecard update.
The opt-in file `.bench/jev-benchmark-opt-in.json` is an authorization record.
Do not copy the opt-in file to `main` without a new user authorization.

## Q2. What did the live runs observe?

### Tested results: first run

The first run used implementation `d3310333` and frozen revision `7cf4105d`.
It stopped after three sealed trials and one interrupted trial of 48.
Source: `~/.bench/experiments/jev-benchmark-live-20260921/report.md:11`.

| Trial | Task | Arm | Status |
| --- | --- | --- | --- |
| 000 | gate-authority-audit | baseline | completed; quality not reviewed |
| 001 | gate-authority-audit | treatment | completed; quality not reviewed |
| 002 | readme-guide | treatment | invalid write fence |
| 003 | readme-guide | baseline | interrupted; not sealed |

Source: `~/.bench/experiments/jev-benchmark-live-20260921/report.md:22`.
The known token estimate is $7.468124938, including $0.002926938 for Jev.
Source: `~/.bench/experiments/jev-benchmark-live-20260921/report.md:72`.

### Tested results: second run

The second run used repair commit `c929c500`.
It stopped after three of 48 trials.
Source: `~/.bench/experiments/jev-benchmark-live2-20260921/report.md:3`.

| Trial | Task | Condition | Status |
| --- | --- | --- | --- |
| 000 | gate-authority-audit | baseline | completed; final gate green |
| 001 | gate-authority-audit | treatment | invalid_rescue_evidence |
| 002 | readme-guide | treatment | interrupted |

Source: `~/.bench/experiments/jev-benchmark-live2-20260921/report.md:27`.
The known token estimate is $8.635839322, including $0.002933322 for two Jev calls.
Source: `~/.bench/experiments/jev-benchmark-live2-20260921/report.md:73`.

In both runs, both Jev calls routed to fallback on required-skill uncertainty.
The runs are too small and too damaged to grade that routing policy.
Source: `~/.bench/experiments/jev-benchmark-live2-20260921/report.md:88`.

### Retained evidence

| Item | Location |
| --- | --- |
| First run report and archive | `~/.bench/experiments/jev-benchmark-live-20260921/` |
| Second run report and archive | `~/.bench/experiments/jev-benchmark-live2-20260921/` |
| Second run configuration digest | `e944144f27c960ada8e24d207802ba14b44647b78deac954b35db2164a5bb3e6` |
| Second run ledger digest | `c2f0d0f64936926aa4270f6dc04d641421680503930fe1a79495bd9c2decbb5a` |
| Second run archive digest | `6e6073558cd06e771cb6f434d19ad58259ca649981d7a4b2ff0dc957cc7002eb` |

These archives are local to one machine. No credential or authentication header is in them.
Source: `~/.bench/experiments/jev-benchmark-live2-20260921/report.md:23`.

## Q3. Which defects stopped the runs?

### Facts: first run

The snapshot builder dropped every tracked symbolic link.
The write fence counted generated verification files as task edits.
Source: `~/.bench/experiments/jev-benchmark-live-20260921/report.md:45`.
Repair cycle four fixed both defects.
Source: `c929c500:capture/retros/jev-benchmark-repair.md:73`.

### Facts: second run

The supplied verifier read the `.git`, `.agents`, and `.codex` directories as prior gate receipts.
Alternate verification roots then failed on a read-only Go build cache and a read-only Bench cache lock.
Source: `~/.bench/experiments/jev-benchmark-live2-20260921/report.md:46`.

The outer preflight and the final gate run outside the task sandbox.
So their green results do not prove that the in-session verifier works under native task permissions.
Source: `~/.bench/experiments/jev-benchmark-live2-20260921/report.md:58`.

Trial 001 named `bench-debug` as a rescued skill. That skill is outside the frozen catalog.
The validator correctly rejected the value.
But the prompt tells the agent to load newly required guidance, and the debug skill applies to failures.
Source: `~/.bench/experiments/jev-benchmark-live2-20260921/report.md:63`.

### Inference

An incomplete benchmark contract explains the rescue failure better than model noncompliance.
Neither explanation shows a Jev prediction failure.
A separate read-only Sol session agreed, at confidence 8.
Source: `~/.bench/experiments/jev-benchmark-live2-20260921/report.md:107`.

## Q4. What must a repair prove?

### Proposals

1. Give gate receipts their own namespace, apart from sandbox-reserved directories.
2. Put bootstrap and gate caches inside writable verification storage.
3. State the rule for required guidance that is outside the selector's catalog.
4. Prove the exact verification path offline in the pinned native sandbox. Make no inference request in that proof.
5. Keep the write fence, the exact source-copy checks, the independent final gate, and the sealed evidence.

Source: `~/.bench/experiments/jev-benchmark-live2-20260921/report.md:96`.

### Authority limits

Five repair cycles are consumed: two default cycles and three user extensions.
The bounded-repair policy needs a reviewer extension before a sixth cycle.
No further paid run has authorization.
Source: `~/.bench/experiments/jev-benchmark-live2-20260921/report.md:103`.

## Residual unknowns

- Actual charges, interrupted-task usage, and nested-review usage are unknown.
- The independent quality review of the captured trials did not occur.
- The creator of the reserved sandbox directories is inferred, not observed.
- The eight development tasks give no holdout evidence.

## Verification record

- [x] The asset opens with its recommendation, scope, and evidence status.
- [x] Facts, tested results, inferences, and proposals are in separate sections.
- [x] Each question has one section.
- [x] Each load-bearing claim cites a branch blob or a local report line.
- [x] The source branch and its exact tip are recorded.
- [x] Unknowns stay explicit.

## Validation plan

1. The reviewer decides whether to authorize a sixth repair cycle, a narrower trial, or no trial.
2. If a repair is authorized, the author works on the source branch from its recorded tip.
3. The author proves the native-sandbox verification path offline before any paid call.
4. A new paid run needs its own configuration acknowledgement and a new run directory.
5. The FT331 pilot keeps its replay-probe requirement whether or not this trial resumes.
