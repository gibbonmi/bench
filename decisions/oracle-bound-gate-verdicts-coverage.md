# Oracle-verdict coverage map

The highest existing seam is the real `bench` CLI in throwaway Git
repositories. Runtime contracts already invoke the subject wrapper selected by
`BENCH_CONTRACT_ROOT`, so nested behavior canaries exercise the same contract
against their deliberately broken fixture surface
([subject selection](../internal/contract/helper.go#L40),
[contract phase](../internal/gate/phases.go#L62),
[behavior-canary routing](../internal/canary/canary.go#L24)).

## Black-box contracts

| Behavior | Real surface and assertion | Named red signal |
|---|---|---|
| Exact reuse | `bench gate` then `bench commit` on the same tree/subject runs the gate once and commits; command, declared environment, ignored input, tool identity, or recognized script/executable-target mutation runs it again and refuses when the new oracle is red. Auto-detect mutation uses ignored selection files so tree equality cannot mask a missing oracle comparison. | `oracle-bound gate verdict contract failed` |
| Closed-subject ceiling | Opaque or remote-marked gates without a complete `.bench/gate-inputs.json` run on every request. A complete manifest enables reuse; absent, empty, malformed, wrong-schema, missing-input, symlink-escape, and control-byte cases remain non-reusable with a reason. | `gate input manifest closure contract failed` |
| Freshness and legacy | Matching green inside ten minutes reuses; expired, future, old three-field, unknown-schema, malformed, truncated, oversized, and wrong-type cache records rerun. Readers do not rewrite them. | `versioned gate verdict cache contract failed` |
| Durable invalidation | A planted old green is durably replaced by `pending` before a marker proves the gate starts. Pre-run replacement failure leaves the marker absent. Green/red finalization failure, subject drift, and cancellation return non-success and never expose the old green. | `fail-closed gate verdict persistence contract failed` |
| Execution ownership | A blocking first gate produces locked-pending status; a second gate, commit, shift, and armed Stop fail immediately without running or staging. Killing the owner leaves interrupted-pending; the next run acquires the released lock and completes. | `single-owner gate execution contract failed` |
| Consumer projection | Status, dashboard, and roadmap context render the same typed absent, reusable-green, red, stale, locked-pending, interrupted-pending, invalid, and unavailable states without re-parsing the file or running the gate. | `gate verdict projection contract failed` |

The commit suite already counts real gate executions through a Git-directory
marker and distinguishes identical-tree reuse from a stale-tree rerun
([fixture](../internal/contract/runtime/runtime_commit_test.go#L11),
[reuse assertion](../internal/contract/runtime/runtime_commit_test.go#L60),
[rerun assertion](../internal/contract/runtime/runtime_commit_test.go#L78)).
Extend that seam for authorization. The runtime gate suite owns resolver, Stop,
recording, lock, crash, drift, and persistence cases; status owns projections;
shift gets one preservation assertion for an operational verdict failure. Each
test invokes the real wrapper and asserts exit code, output, cache bytes/state,
run count, HEAD/index, and retained work as applicable.

## Lowest extra seam

`internal/gate` gets one injected clock/filesystem/lock seam only for failures
that cannot be made deterministic and portable through the CLI: file sync,
directory sync, failure after rename, future-clock evaluation, lock acquisition
errors, and exact strict-decoder byte limits. Table tests assert ordered calls
and the resulting typed state. All other cases stay at the CLI seam; source
inspection or a second cache parser is not evidence.

Runtime fixtures can still force reliable filesystem failures without the
injection seam where possible: make the cache path a directory to fail pre-run
replacement, or let the gate replace the pending file with a directory before
finalization. Mid-run tracked mutation, blocking markers, signals, and external
symlink targets are likewise observable through the real process.

## Canary attachment

Add two `behavior-owned` fixtures, because that family runs only the contract
phase while the empty baseline still prevents a vacuous expectation
([family ownership](../internal/canary/canary.go#L28),
[vacuity check](../internal/canary/canary.go#L150)).

1. `gate-verdict-oracle-binding-bypassed` supplies a subject wrapper that trusts
   matching tree/status while ignoring the oracle digest. It must trigger
   `oracle-bound gate verdict contract failed`.
2. `gate-verdict-invalidation-bypassed` supplies a subject wrapper that runs
   before durable pending and permits an older green after a final-write
   failure. It must trigger
   `fail-closed gate verdict persistence contract failed`.

Register both as behavior-owned; the registry requires every fixture to have
one owner matching its family
([registry contract](../internal/conformance/registry_test.go#L130)). These two
canaries defend the independent authorization and persistence invariants; the
remaining rows are permutations of those real contracts rather than new
tripwires.

## Hostile-input ownership

- The manifest/cache codec owns absent versus empty, no trailing newline,
  truncation/oversize, wrong type/schema, control bytes, and paths containing
  spaces or glob characters.
- Subject construction owns symlink chains, ignored paths, missing tools,
  executable modes, and deep-CWD root resolution.
- The execution protocol owns signals, crashes, repeat runs, lock contention,
  and plan/run drift.
- The runtime matrix reaches standalone gate, commit, shift, Stop, status,
  dashboard, and roadmap context; adapters remain thin because they reach those
  same commands.
