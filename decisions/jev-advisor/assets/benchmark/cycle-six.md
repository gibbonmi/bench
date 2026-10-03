# Sixth offline Jev benchmark repair

Status: repaired offline candidate; native full-gate qualification remains blocked.

## Authority and source

The reviewer approved one sixth repair cycle on 2026-10-03.
Five earlier cycles remain consumed. This cycle adds no paid-run authority.
All author work stays in the current Astra session, with no delegates.
The CLI/Desktop assignment remains separate.

The recovered benchmark source comes from `c929c500501bc86606a0dae55923319a34e685d0`.
The research checkpoint before this repair is `c00b78ea82df94487642712ef3359cb821e62fa7`.
The local assignment is `jev-shaping`, request `jev-research-resume-20261003`.
`verification.json` identifies the repaired source and preserves earlier verification records.
`review.md` remains the historical review through cycle five.

## Repairs and evidence

The real verifier CLI first failed while reading `.agents/result.json`.
The same command reproduced that failure twice before implementation.
The verifier now distinguishes Codex's protected metadata directories from its gate receipts.
Unknown entries and malformed receipts still fail validation.

The next native run reproduced a write failure at the ambient Bench cache lock.
Each verification run now owns its subprocess home, build caches, temporary files, and Go configuration.
The verifier resolves installed module paths before changing the subprocess home.
Go records those paths in the private home so the gate's environment filter cannot lose them.
Module fetching and Go telemetry stay disabled. Checksum verification remains enabled.

Telemetry initially raced temporary-directory cleanup in the regression suite.
The cache test proved that the private Go configuration still allowed telemetry.
It failed before the correction and passed after Go disabled telemetry in that home.
A separate suite attempt correctly refused source changes made while it was running.
These failed attempts remain evidence; they are not passing qualification runs.

The frozen packet now names unindexed phase guidance separately from selection candidates.
This permits a `bench-debug` rescue when governing rules require diagnosis.
Unknown names, duplicate rescues, and rescues already supplied remain invalid.
The packet validator derives the extra names from frozen skill files; omitted guidance fails validation.
The frozen source inventory continues to seal full skill bodies and paths.

The final standalone suite passes all 43 tests inside the pinned Codex sandbox.
Five omission mutations fail their named assertions.
The telemetry test adds an observed red-before-green result.
The existing checks still require independent final verification and reject writes outside the task fence.

All eight recorded task packets freeze and validate at their original revision.
The largest state-plus-question size is 119,228 bytes against the unchanged 120,000-byte ceiling.
The provider's exact token fit remains uncertified.
The current repository source rejects a stale consumers excerpt rather than silently changing the recorded task.
The original full governing documents remain in every packet.

## Native full-gate result

The installed Codex executable matches the archived SHA-256:
`0753dfe1d8b87a52436deb13eb1c549661ef4c84fee2c5aa688385eebeccb761`.
Its version is `0.155.1-x86_64-unknown-linux-musl`.
The probe uses `codex sandbox`, with workspace writes and a separate verification root.
Networking remains disabled. No native inference session or Jev request runs.
The [permission documentation](https://learn.chatgpt.com/docs/permissions) describes the profile contract; the installed command help supplies the exact syntax.

The actual build succeeds. The gate passes formatting and vet, then reaches its native tests.
Unix-socket fixtures fail with `setsockopt: operation not permitted`.
The models package also cannot create its local TCP test server.
The same run retains a separate handoff-document concurrency failure; its cause remains unclassified.
The gate returns failure, and the verifier preserves that verdict and the unchanged source digest.

The full native run takes about 288 seconds; this is diagnostic runtime, not a benchmark comparison.

The offline repair therefore does not qualify the paid benchmark for restart.
A green Python fixture suite cannot substitute for this red full-gate result.
No gate, fixture, write fence, or provider threshold was weakened.

## Review pickup

The author checked standards, scope, and coverage against the current source.
This is an author review, not a fresh independent review of the six-cycle candidate.
The earlier independent reviews do not certify the current bytes.
The remaining blocker is full-gate execution under the benchmark's native permission contract.
Independent review and exact configuration approval remain necessary before paid execution.

One possible next design lets the controller process fixed-command verification requests while the task session waits.
Each request would seal the exact authored inputs and return a separately captured gate result.
The final controller gate would still run independently, even after earlier successful feedback.
The agent would receive no authority to change the gate command, permissions, or write fence.
That proposal needs a separate execution decision and an additional repair allowance; this cycle does not implement it.

The repair allowance is exhausted at six cycles.
The [bounded repair policy](../../../../.agents/skills/bench-craft-line/references/bounded-repair-policy.md) requires a reviewer extension for another repair cycle.
Accepted decisions 1A, 2A, and 3A remain closed.
Recorded briefs stay development-only; fresh real tasks remain necessary before adoption.

## Evidence locations

The assignment retains `.logs/jev-cycle-six/` with commands, failures, test outputs, mutations, frozen packets, and native gate receipts.
The durable preservation bundle lives under `/home/mgibs/.bench/experiments/jev-research-recovery-20261003/cycle-six/`.
The bundle manifest binds each preserved regular file or symbolic link to its content.
`verification.json` records the current source digests and the compact results.
No original probe or earlier experiment archive is overwritten.
