# Actual-interface qualification

Status: incomplete. No complete compatibility claim or landing is authorized by this record.

## Desktop shell repro

The author ran this debug pass in the actual Desktop chat.
The command was `pwd`, with normal sandbox permissions and the repository as its working directory.
Two consecutive calls failed before the shell started.
Changing only the directory to `/tmp` produced the same failure.
Selecting `/bin/bash` explicitly also produced the same failure.

```text
Failed to create unified exec process: No such file or directory (os error 2)
```

The error text has SHA-256 `fdbad5ea73e6aea9bb40573ae3923d834eb8d9d29e51d39ba3c890cf7141055b`.
The original tool results remain in chat `01a0fc58-0f59-7f52-9c4c-a25247d5fa21`.
Elevated diagnostic execution succeeds but does not qualify normal execution.
The repository, `/tmp`, `/bin/bash`, and `/bin/sh` exist.

## Runtime observations

The Desktop server PID was 3083902 when inspected.
Its executable was `/mnt/c/Users/gibbo/.codex/bin/wsl/3ac368078cf7546b/codex`.
Its declared configuration home was `/mnt/c/Users/gibbo/.codex`.
These are direct process observations, not a claim about every effective policy source.
The repository identity is `/home/mgibs/workspace/bench`.

The server held descriptor 3 for this lock:

```text
/mnt/c/Users/gibbo/.codex/tmp/arg0/codex-arg0q0yovE/.lock
```

That directory was absent during inspection.
Its `codex-linux-sandbox` and `codex-execve-wrapper` entries were also absent.
The CLI server's corresponding runtime directory existed with valid launcher links.
This comparison supports a missing Desktop launcher hypothesis.
It does not establish who removed the directory or prove the exact failed executable.

The local log database supplied no matching startup-error record.
No syscall tracer was installed.

## Ranked hypotheses

1. A missing Desktop sandbox launcher prevents process creation. The missing active runtime directory supports this hypothesis.
2. The working directory is absent or inaccessible. The same failure from `/tmp` weakens this hypothesis.
3. The configured shell is absent. The same failure with an existing explicit `/bin/bash` weakens this hypothesis.

The cause remains unconfirmed until a supported recovery changes the actual-tool result or direct tracing identifies the failed executable.

## Recovery boundary

No private runtime file, policy, trust setting, or process was changed.
A controlled Desktop restart remains a possible recovery test after active work is preserved and the reviewer approves interruption.
After restart, run `pwd` and the repository wrapper's `version` command separately through normal Desktop tools.
A terminal, hook, new subprocess, or elevated command cannot supply these results.

[Official troubleshooting](https://learn.chatgpt.com/docs/reference/troubleshooting) recommends preserving active work before an app restart for persistent terminal problems.
That general recovery advice does not establish the cause of this agent-tool failure.

## Retained implementation before resumption

C1 and C2 passed their checkpoints.
C3 has six uncommitted source and test files in the owned assignment.
The session-inspection package passed in 159 milliseconds after the initial C3 changes.
Its unknown-context and startup-obligation tests each had a behavioral red before implementation.
No C3 gate or independent review has run.
A verified private preservation copy is `/tmp/cli-desktop-c3-before-runtime-recovery.tar`.

## Qualification before resumption

The reviewer previously supplied successful normal-permission CLI `pwd` and wrapper `version` results.
That transcript does not qualify the remaining CLI workflow or any Desktop operation.
All required live rows still need a complete reconciliation against the approved coverage map.
Desktop recovery, fresh and resumed chats, file access, skill invocation, hook refusal, assignment isolation, review, and landing remain incomplete.

## Resumed Desktop observations

On 2026-10-03, separate actual Desktop calls ran `pwd` and the repository wrapper's `version` command.
Both returned exit 0 without an elevation request.
The results were `/home/mgibs/workspace/bench` and `bench 0.2.0 (linux/amd64)`.
Tool chunks `17be74` and `2b0058` retain the repeated observations.
The user then identified Custom permissions as necessary for successful execution.

The current session declares unrestricted filesystem access and approval policy `never`.
The Desktop configuration declares `sandbox_mode = "danger-full-access"` and `approval_policy = "never"`.
These results qualify current unrestricted execution only.
They do not demonstrate recovery under the standard workspace sandbox.
No configuration, trust setting, private runtime file, or running process was changed.

Separate actual Desktop shell calls wrote and read the expected scratch bytes.
Their tool chunks are `f0088f` and `c5e320`.
The scratch file remains at `.logs/cli-desktop-handoff/resume-desktop-file-probe.txt` in the assignment.
Its SHA-256 is `0fdb882cb094b721e672725d8e5b258e5e0babbb949aa329a901ae4cfe1fc847`.
The user invoked `bench-debug` and `bench-implement-spec` from their installed repository skill paths in this chat.

An actual Desktop call attempted a harmless forbidden-command fixture.
The fixture used an invalid Git directory and an invalid option, so Git could not reset a repository.
The PreToolUse hook rejected the call with `BLOCKED: git reset` before execution.
The subsequent call confirmed that `/tmp/bench-desktop-guard-d3rkk32r/sentinel` was absent; tool chunk `fe8449` retains that check.
This supplies current Desktop hook evidence for CD50.
It supplies no CLI hook evidence.

## Sandbox diagnosis

Recommendation: retain the sandbox boundary and verify the smallest required path grants before changing an active permission profile.
The current Custom setting removes the sandbox restrictions.
The evidence confirms two path conflicts in a fresh diagnostic sandbox.
It does not qualify the actual Desktop sandboxed tool route.
Refresh these observations after runtime, profile, workspace, or Bench storage changes.

### Question 1: Which operation fails?

The installed Desktop executable is `/mnt/c/Users/gibbo/.codex/bin/wsl/095c52da468c9593/codex`.
Its `sandbox --help` exposes an explicit permission-profile selector.
This command reproduced the Bench failure:

```text
/mnt/c/Users/gibbo/.codex/bin/wsl/095c52da468c9593/codex sandbox -P :workspace -C /home/mgibs/workspace/bench -- /home/mgibs/workspace/bench/bin/bench.sh worktree exec cli-desktop-consistency -- bench test --package ./internal/compatibility
```

Tool chunk `34e711` reports exit 1:

```text
error: cache lock unavailable — open /home/mgibs/.cache/bench/go-build/bench.lock: read-only file system at /home/mgibs/.cache/bench/go-build
```

The same profile successfully ran `pwd`, wrapper `version`, and assignment lookup.
Tool chunks `6a1b4b`, `0f773d`, and `67af25` retain those controls.
A scratch creation in the cache also failed with `EROFS`; tool chunk `21f11f` retains that result.
An assignment scratch creation failed with `EROFS`; tool chunk `13d4ef` retains that result.
Its error report also emitted `spill-failed{reason=open: read-only file system}`.

### Question 2: Why does the boundary conflict with Bench?

[OpenAI permission documentation](https://learn.chatgpt.com/docs/permissions) states that the workspace profile permits workspace and temporary-directory writes.
It also documents additional roots and narrower filesystem grants.
Retrieved: 2026-10-03.

| Bench resource | Source | Observation |
| --- | --- | --- |
| Home and assignment pool | `internal/benchhome/benchhome.go`, `fallbackDir`; current assignment path | The pool is outside the active repository root. Its scratch write fails under `:workspace`. |
| Go cache | `internal/gocache/gocache.go`, `Dir` and `Apply` | The cache comes from HOME. `Apply` replaces an inherited GOCACHE override. The lock write fails under `:workspace`. |
| Complete response storage | Current sandbox assignment probe | The sandbox refuses the spill write after the original assignment write fails. |

These results establish a Bench storage-boundary mismatch for the tested workspace profile.
They do not establish that unrestricted execution is the only possible configuration.
A narrower supported profile remains untested.
No path grant was applied.

### Question 3: Does this explain the original ENOENT?

No Bench write occurs before the shell starts.
A cache refusal therefore does not explain the earlier process-creation ENOENT.

The current Desktop server is PID 2180.
Its open lock points at `/mnt/c/Users/gibbo/.codex/tmp/arg0/codex-arg0U6lZhG/.lock`.
That directory and its sandbox launchers were absent during inspection; tool chunk `82eef3` retains the observation.
The current executable differs from the earlier recorded executable.
The missing-directory symptom therefore remains relevant after the runtime change.
The deletion actor and exact failing executable remain unknown.

A fresh process of the installed runtime successfully launched its workspace sandbox.
The host also exposes `/usr/bin/bwrap`.
[OpenAI sandbox documentation](https://learn.chatgpt.com/docs/sandboxing) describes the Linux and WSL2 sandbox prerequisites.
Retrieved: 2026-10-03.
The successful fresh launch weakens a general missing-sandbox-support hypothesis.
It does not validate the long-running server's missing launcher path.

### Validation plan

1. Capture `pwd` and wrapper `version` through the actual Desktop tool with its standard sandbox selected.
2. If startup fails, retain that exact error before any Bench diagnosis.
3. Obtain explicit authorization before any process recovery, runtime edit, or permission change.
4. Test the smallest required Bench path grants in a disposable diagnostic profile.
5. Repeat the failing Bench test and assignment scratch probe through that profile.
6. Verify actual Desktop behavior separately after any approved recovery or configuration change.

Independent CLI qualification, concurrent writers, fresh and resumed chats, review equivalence, and disposable landing remain incomplete.
C3 cannot pass its live qualification or land from this evidence.
## Sandbox debug disposition

The sandbox storage failure is reproduced; the approved diagnostic grants resolve both tested operations.
The original Desktop process-creation failure has no confirmed cause.
The reviewer approved the one-off diagnostic grants, which were applied only to those child processes.
No source repair for sandbox storage has been made.

The exact proposed invocations are in `.logs/cli-desktop-handoff/sandbox-probe-plan.json`.
The proposal adds writes to the canonical Go cache, this assignment, and this assignment's response directory.
It keeps network access disabled and changes no persistent Codex configuration.
These grants do not claim to support the complete Bench workflow.
The actual Desktop route requires separate qualification.

The cache refusal has SHA-256 `a72e455fcb66f59878a7a187123172aad5c684f493b7caedc7ff3822c173333e` for its UTF-8 text without a trailing newline.
The failing surface is a Bench package test inside the installed runtime's workspace sandbox.
The ranked hypotheses are storage outside writable roots, a missing long-running server launcher, and missing host sandbox support.
The first explains the reproduced storage refusal.
The second remains relevant to ENOENT; the successful fresh sandbox weakens the third.

The plan file retains all 14 dirty paths.
C3 remains in its owned assignment, with C1 and C2 accepted.
Its package, system, and conformance tests pass; full live qualification remains incomplete.
Independent C3 review has not run.
Implementation commits do not authorize complete qualification or landing.

## Approved sandbox probes

The reviewer approved the one-off sandbox proposal in this chat.
The original dotted override keys failed parsing before sandbox startup.
Passing the same three grants as one filesystem table corrected that invocation.
The grants and network restriction did not change.

Both corrected probes returned exit 0; tool chunk `ec1be8` retains the results.
The compatibility package passed in 13 milliseconds with no skips.
The assignment scratch probe wrote, read, and removed its temporary file.
The exact invocations and outputs are in `.logs/cli-desktop-handoff/sandbox-probe-corrected-results.json`.
No persistent configuration, private runtime file, or process was changed.

The successful probes confirm that these operations can run with narrower grants than unrestricted access.
They do not qualify the actual Desktop tool route under standard permissions.
They do not qualify Git writes, the complete gate, or a landing under the diagnostic profile.
The missing launcher hypothesis remains unresolved.
C3 implementation can continue through the currently working Desktop route while these live qualification limits remain explicit.

## C3 implementation verification

The session reporter rejects stale, unknown, elevated, hook-only, and subprocess observations.
The startup path emits the shared obligation before its existing phases.
The reference guide owns the procedure, and the core guide points to it.
The same payload installer ships both documents to consumers.

The required comparability omission initially survived because fingerprint drift still rejected the changed context.
The corrected test supplies matching unknown contexts, and that omission now fails `TestCompatibilityUnknownContext`.
The lifecycle omission fails the startup, retest, and unresolved-repair unit tests.
Through the real installed hook, it also fails `TestCompatibilityResume`.
Every probe restored its exact subject before subsequent work.

| Mutation or omission | Observed failing test | Evidence |
| --- | --- | --- |
| Require optional presentation for diagnosis | TestCompatibilityOptionalTool | Tool chunk ab1e43 |
| Accept hook, subprocess, elevated, or unnamed routes | TestCompatibilityEquivalentEvidence | Tool chunk f6bf4b |
| Accept a failed equivalent outcome | TestCompatibilityEquivalentEvidence/failed | Tool chunk 6a2aa0 |
| Reuse another operation's observation | TestCompatibilityEquivalentEvidence/another_operation | Tool chunk eec432 |
| Ignore context fingerprints | TestCompatibilityContextInvalidation | Tool chunk 940e1c |
| Require equal launcher labels | TestCompatibilityVersionDifference | Tool chunk 445360 |
| Omit native Windows refusal | TestCompatibilitySupportBoundary | Tool chunk 6446b4 |
| Ignore chat and lifecycle identities | TestCompatibilityContextInvalidation | Tool chunk aede67 |
| Make an expired startup return a failure exit | TestCompatibilityStartupTimeout | Tool chunk 4a2648 |
| Omit the core session pointer | TestCompatibilityPayload | Tool chunk 488eda |
| Omit the lifecycle phase | TestCompatibilityResume | Tool chunk d61813 |
| Omit context comparability | TestCompatibilityUnknownContext | Tool chunk 05b896 |

A later failed observation initially retained the earlier success's action of `none`.
`TestCompatibilityFailureAfterSuccess` reproduced that defect before the repair; tool chunk `50e0f1` retains the red.
The shared reporter now updates the failed state and its action together.
The compatibility package then passed with no skips; tool chunk `c2bf50` retains the green.

## Live acceptance disposition

These rows remain explicit acceptance limits at this implementation stage.
Automated fixtures do not close an actual-interface row.
No complete qualification or landing is authorized.

| Rows | Current evidence | Remaining requirement |
| --- | --- | --- |
| CD39, CD40 | Earlier failed Desktop calls and successful diagnostic controls remain recorded. | Independent review must confirm the provenance distinction. |
| CD42 | The workflow stopped at unverified dependent operations. | An actual missing-capability exercise must retain its untouched mutation sentinel. |
| CD43, CD60 | Static recovery guidance names the external route and unknown repair boundary. | Independent review must compare it with the recorded startup failure. |
| CD45 | Current Custom Desktop calls pass. | The previously failed standard Desktop route has no verified recovery. |
| CD46 | This Desktop chat read the repository agreement and installed skills. | A fresh independent CLI chat must read the same committed rules. |
| CD47 | This writer retains its distinct assignment. | Concurrent actual CLI and Desktop writers must exercise separate assignments. |
| CD48 | No disposable live lifecycle is claimed. | Both actual interfaces must complete review and landing in a disposable repository. |
| CD49, CD50 | The actual Desktop guard refused the harmless fixture before its sentinel. | The actual CLI must supply the corresponding refusal. |
| CD51 | No equivalent-route comparison is claimed. | Observe and compare both review outcomes. |
| CD54 | The qualification status remains incomplete. | Retain this block until every required live row closes. |
| CD57, CD58 | Actual Desktop file access and skill invocation are recorded. | Obtain corresponding actual CLI evidence. |
| CD62 | Unit probes detect context changes. | Record affected checks repeating after an actual configuration change. |

## Pre-review automated results

The retained source passed these checks before its implementation commit.
All deliberate mutation subjects were restored and verified.
The ordinary lane and independent reviews remain separate requirements.

| Check | Result | Evidence |
| --- | --- | --- |
| Compatibility package | Pass; 9 ms; no skips | Tool chunk c2bf50 |
| Session-inspection package | Pass; 187 ms; no skips | Tool chunk 09dc5a |
| Adoption package | Pass; 81238 ms; no skips | Tool chunk 4f60b6 |
| Conformance package | Pass; 62181 ms; three capability skips | Tool chunk 90ec2c |
| Anchors package | Pass; 1500 ms; no skips | Tool chunk 70cc35 |
| Command package | Pass; 28541 ms; no skips | Tool chunk 3306e1 |
| Restored sealed system suite | Pass; 179622 ms; no skips | Tool chunk 2447c0 |
| Live-root documentation checks | Four checks pass; no skips | Tool chunks af9da6, 05af0e, 3dfbee, 57ae0d |
| Coverage map | Valid; all 64 rows cited | Tool chunk 5998c7 |
| Source growth | Pass against accepted C2 base | Tool chunk 3fcd71 |

The conformance skips concern two long Unix socket paths and an unavailable character-device creation capability.
No test or skip policy was weakened.
These skips do not supply positive evidence for those fixture cases.

The author checked the delta for duplicated facts and comment contracts.
Capability requirements and actions have one registry, and context reuse has one comparison owner.
The system tests reuse the existing fixture and complete-output helper.
The red probes above demonstrate the independent behavioral expectations.
