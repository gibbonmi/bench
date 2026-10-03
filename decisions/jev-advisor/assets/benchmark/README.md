# Jev skill-selection benchmark

Status: sixth offline repair cycle; paid restart and policy adoption remain unapproved

This benchmark measures whether Jev-assisted initial skill selection improves complete task cost without reducing task quality.
It does not adopt a production skill policy.
The earlier trials remain unchanged.

## Evidence contract

The freeze command reads one committed repository revision.
It copies every regular tracked file, including governing documents, source, tests, and skill resources.
It also preserves committed symbolic links and seals their target text.
Links must resolve within the frozen source and cannot point into Git metadata.
The validator refuses absolute, dangling, cyclic, and escaping links.
Harness links and canary fixture overlays are executable inputs; canonical bodies alone are insufficient.

Each task packet contains full governing documents, task source evidence, skill descriptions, protected skill bodies, and the action scope.
Selected task sources use line-numbered excerpts with checked anchors and whole-file digests.
Other task source files remain complete in their packets.
The reference guide and glossary remain complete on disk as lookup documents.

The validator checks each source digest, executable mode, excerpt, and the complete frozen file inventory.
It rejects missing rules, missing source, altered descriptions, unavailable protected skills, and ambiguous scope.
It reads catalog membership from the canonical generated index.
It reads names and descriptions from the indexed skill sources.

The provider documents separate request and state-plus-question token limits.
The retrieved provider documentation exposes no exact offline tokenizer.
Preparation therefore applies conservative byte ceilings and reports that token fit remains uncertified.

An oversized packet fails before a provider call; the runner never truncates evidence silently.
An authentication or input rejection stops the comparison without native fallback.
It cannot count as a model-selection failure.
The plan includes the observed sizes and the context-limit source.

Both arms receive the same packet and a separate workspace from the same snapshot.
Both arms include source inspection, verification, self-review, and repairs within the task session.
The runner performs the mandatory final gate in a separate copy of the exact authored files.
The task agent can request earlier gate feedback through the supplied verification command.
Focused checks and self-review remain available during the task.

The replay excludes production ticket or spec authoring, commits, landings, and publication.
This explicit replay exception applies equally to both arms.

The baseline chooses initial skills through the native task agent.
The treatment receives the Jev selection or the accepted native fallback selection.
Both arms retain access to the complete catalog and skill bodies.
The task agent must load additional required guidance when it discovers a requirement.
Its final record names those rescues.

The packet lists unindexed phase-skill names in `additional_guidance`.
The frozen source owns their bodies and file hashes.
Packet validation derives these names from the frozen skill files.
They can qualify as rescues but do not enter the initial-selection questions.

Each workspace has its own Git repository and baseline commit.
Copies preserve the frozen links.
That local baseline supports diff inspection and package tools; it is not a production commit.

Full project rules expose some skill-discovery guidance.
Therefore, this benchmark does not claim to remove all catalog tokens from the treatment.
It measures the remaining work under the complete rules.

## Questions and routing

Each candidate receives two independent Noul questions over the same state.
The first asks whether the current action requires the skill.
The second asks whether the skill provides optional benefit, assuming it is not mandatory.
Each question carries its candidate identity and complete description.

The policy constants in `evidence.py` own the thresholds.
Required uncertainty triggers native fallback.
Optional uncertainty alone does not trigger fallback.
Explicitly protected skills survive every provider outcome.
The terminology task retains its exact reviewer-approved optional exception.

A failed, malformed, incomplete, or mismatched Jev response triggers fallback.
An invalid fallback stops the dependent task.
The runner records failures and makes no automatic retry.
Missing evidence stops the run before a provider call.

## Task population

`corpus.json` owns each task's origin, family, action, write fence, protection evidence, and quality criteria.
The corpus contains source-derived tasks across implementation, documentation, diagnosis, and policy analysis.
The first three reproduce the exposed consumers, onboarding, and terminology tasks.
The other tasks cover assessment diagnosis, assessment documentation, comment maintenance, skill-trigger analysis, and gate authority.

All supplied cases are development material.
They are not a random sample of production work or an independent holdout.
A later holdout needs new task families, independent labels, and a freeze before prediction inspection.
The validator refuses a family that crosses the development and holdout partitions.

The schedule randomizes task pairs with a recorded seed.
It alternates the leading arm across pairs and requires repeated pairs.
Both arms use the same native model, effort, service tier, timeout, and harness version.
Provider cache effects and model randomness remain uncontrolled.
Report per-task paired outcomes and variation; do not infer equivalence from a small average difference.

## Offline preparation

Run these commands through the owning Bench worktree.
The following paths are relative to that worktree.

```sh
python3 decisions/jev-advisor/assets/benchmark/test_benchmark.py -v
python3 decisions/jev-advisor/assets/benchmark/runner.py freeze --repo . --out /tmp/jev-benchmark-inputs
python3 decisions/jev-advisor/assets/benchmark/runner.py validate /tmp/jev-benchmark-inputs
python3 decisions/jev-advisor/assets/benchmark/runner.py plan /tmp/jev-benchmark-inputs --out /tmp/jev-benchmark-config.json --model gpt-6-astra --effort high --harness 'codex-cli 0.155.1' --budget-reference 'Pending reviewer approval'
```

Preparation makes no model call.
The default plan has three repetitions per task.
The plan records the maximum adapter calls, source fingerprint, and configuration fingerprint.
Review the complete config and request packets before authorizing a paid run.
Populate rate provenance when estimated monetary comparisons are required.
An empty rate remains unknown; it never means free usage.

## Execution

The runner accepts explicit adapter argument arrays.
The supplied adapters call hosted Jev and the installed Codex executable.
The Codex adapter checks its version and executable digest, then invokes the pinned absolute path with the specified model and effort.
It retains native events, final output, and every completed parent turn's usage.

The pinned ephemeral exec interface does not certify nested-review usage coverage.
Every native export therefore carries an unknown additional-cost component.
Known parent-token estimates remain visible.
The benchmark claims no complete native bill or complete total-cost comparison.
User configuration is excluded; project rules and the replay packet remain present.

The live command requires the SHA-256 of the exact reviewed configuration.
The configuration also pins the frozen evidence fingerprint.
It pins executable, adapter, and local implementation bytes and checks them before each trial.
The hash is a deliberate launch acknowledgement, not an independent consent authority.
Use the hash from the plan command only after approval.

```sh
python3 decisions/jev-advisor/assets/benchmark/runner.py run /tmp/jev-benchmark-inputs --config /tmp/jev-benchmark-config.json --out /tmp/jev-benchmark-run --authorize-config-sha256 REVIEWED_CONFIG_SHA256
```

`TYPESAFE_API_KEY` supplies the Jev credential.
The frozen source must also contain a committed `.bench/jev-benchmark-opt-in.json` file.
Its complete value must be `{"hosted_jev_benchmark": true}`.
The user authorized the first live run and committed this opt-in.
A new run still requires its exact configuration acknowledgement and a green baseline preflight.

`BENCH_OFFLINE` prevents a live run and prevents the hosted Jev adapter from posting.
The adapter posts the same bytes that the runner saves.
The runner does not save authentication headers or pass the Jev key to native tasks.

Before the first adapter call, the runner builds the frozen source and runs its full gate in a fresh verification workspace.
A red, interrupted, or unavailable baseline gate stops all paid calls.
The baseline receipt, outputs, source inventory, and generated files join the run ledger.
Each successful task submission receives an independent final gate from the runner.
Task-created verification receipts cannot replace this final check.

Verification copies are distinct from authored workspaces.
The verifier accepts Codex's protected metadata directories at the output root.
It rejects other entries that lack a valid receipt.

Each verification run owns its subprocess home, build caches, and temporary files.
Go records the resolved module locations in that private home before the gate filters the environment.
The verifier disables module fetching and retains checksum verification.
It disables Go telemetry in the private home so background writes cannot race evidence capture.

The builder and gate commands, input digest, outputs, generated-file inventory, and elapsed time remain sealed evidence.
The authored workspace retains its exact write fence, including unknown ignored files and mode changes.
There is no blanket exemption for ignored files, gate logs, or build directories.

The runner bounds adapter calls and elapsed execution time.
These limits are not a hard monetary billing cap.
The native provider can perform several inference calls within one task adapter call.
Do not describe the adapter count as the provider request count.
The output directory must be new, so a rerun cannot silently replace a failed attempt.

## Measurement and quality

Each adapter call has one monotonic duration from launch through response capture and decode.
Each trial has one monotonic duration from workspace setup through artifact capture.
Ledger sealing and report generation are separate benchmark overhead, outside that task duration.
The baseline gate is preparation overhead, outside paired task duration.
Each trial includes its independent final gate and any earlier task-requested gates.

UTC timestamps identify events; they do not supply a second elapsed-time estimate.
Do not sum concurrent spans or compare historical conflicting clocks with these measurements.

The runner retains every attempted role, including failures, interruptions, and fallbacks.
An interrupted adapter retains its result and any captured native events.
Known completed turns survive; uncaptured usage remains unknown.
A coordinator interrupt stops later calls after the captured trial is sealed.

A stopped run has missing rows in the schedule summary.
If interruption leaves a trial without a completed integrity record, reports refuse to score the run and retain its raw evidence for diagnosis.
Frozen-source preparation and independent quality assessment remain separate research overhead.
Charge that overhead separately when assessing deployment economics.

```sh
python3 decisions/jev-advisor/assets/benchmark/report.py summary --runs /tmp/jev-benchmark-run
python3 decisions/jev-advisor/assets/benchmark/report.py reviews --frozen /tmp/jev-benchmark-inputs --runs /tmp/jev-benchmark-run --out /tmp/jev-benchmark-review
python3 decisions/jev-advisor/assets/benchmark/report.py export --frozen /tmp/jev-benchmark-inputs --runs /tmp/jev-benchmark-run --out /tmp/jev-benchmark-assessment --repo-key REPOSITORY_POOL_KEY
```

Give the independent reviewer only the review packets and frozen source.
Keep the private index away from that reviewer until criterion judgments are complete.
The packets omit the condition, selection, provider probabilities, and measured cost.
They retain condition-neutral completion validity, the changed-file list, write-fence failures, and counts of nonblank verification and self-review claims.
Raw claims stay in the private run evidence; the task adapter receives no literal arm label.
Trials that never reached the task adapter stay in the overall failure accounting but produce no quality-review packet.

The run ledger binds the plan, baseline verification, and every captured trial file.
It covers metadata, requests, responses, native events, workspace content, and file modes.
Review, export, and summary verify that ledger and the configuration fingerprint before producing output.
Changed, deleted, added, or unsealed trial evidence causes refusal.
Changed baseline evidence also causes refusal.

Freeze version four and ledger version two require a new run; use the archived verifier for older evidence.
The ledger detects evidence drift; it is not a signature against an operator who rewrites both evidence and ledger.
The runner returns the ledger digest for retention outside the run directory.

Task-authored artifacts may still suggest the arm, so reviewers must record any suspected disclosure before unblinding.
Identical evidence packets share one review and retain every associated trial in the private index.

The reviewer records pass, fail, or unknown for every frozen criterion, with supporting source or test evidence.
The reviewer checks project standards as well as artifact function.
Run the full named package for code tasks, and verify write fences for all tasks.
Do not accept a native agent's verification claim as independent quality evidence.

Import exported run records with `bench assessment record --input <file>`.
Bench's existing assessment owner derives token totals and monetary estimates.
Exported measures retain each role's elapsed time and the native task's validated rescue count.
Claims must be nonblank strings.
Rescues must name unique catalog skills or frozen additional guidance, absent from the supplied initial selection.

The exporter leaves quality and completion assurance unresolved.
Fill those fields from the independent review before a quality-adjusted comparison.
Actual charges require separate billing evidence.
The benchmark never substitutes token estimates for actual charges.

The summary deliberately reports no adoption evidence.
A complete evaluation must account for failed tasks, required-skill rescues, fallback costs, repairs, independent quality, and missing measurements.
The first live run stopped after three sealed trials and one interrupted trial.
It exposed missing snapshot links and incorrect classification of gate-generated files.
Those observations diagnose the benchmark; they do not establish Jev performance.

## Verification and source references

`verification.json` records the offline checks and observed mutation failures.
`cycle-six.md` records native sandbox verification and the remaining adoption limits.
The recorded corpus pins source excerpts from the archived benchmark revision.
Freeze that revision in a separate local checkout; changed source must pass fresh excerpt checks.
Do not relax the context ceiling or remove governing documents to make a packet fit.

The standalone Python suite exercises the benchmark interface without model calls.
The repository gate does not discover this map-owned Python suite; run it explicitly before the gate.

The HTTP request shape follows the [TypeSafe API](https://docs.typesafe.ai/api.md).
The independent question meanings follow the [Noul documentation](https://docs.typesafe.ai/primitives/noul.md).
The byte preflight uses the [model limits](https://docs.typesafe.ai/models.md) and their approximate English-text guidance.
The native adapter follows the installed `codex exec --help` contract and [non-interactive documentation](https://learn.chatgpt.com/docs/non-interactive-mode).
These sources were read on 2026-09-21.
