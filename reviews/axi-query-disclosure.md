# AXI query disclosure implementation review

Candidate: `c2542a97..68e9416e` (reviewed again through phase-only HEAD `68e9416e`)

Raw findings: Standards 2, Spec 5, Coverage 2. De-duplicated repair targets: 7.

## Standards

Count: 2. Worst: the reusable action owner rejects one magic test sentinel rather than a defined invalid state.

1. **Medium — ask-user.** `internal/axi/action.go:141` treats the literal known argument `"unknown"` as guessed, while every other guessed literal renders and a legitimate literal `unknown` is refused. The matching assertion at `internal/axi/action_test.go:104` is fitted to that sentinel. Decide whether provenance belongs at call-site derivation tests (recommended) or whether the public action type needs a different enforceable contract.
2. **Medium — auto-fix.** The approved set and fixed 49-entry count are repeated in `internal/conformance/axi_query_registry_test.go:16` and `cmd/bench/command_registry_test.go:53`. Keep one independent conformance expectation and let package-main tests bind registry-derived members to the real envelope matrix; remove the fixed count, because exhaustive per-entry dispositions already fail closed.

## Spec

Count: 5. Worst: coverage extraction changes an existing successful exit from 0 to 1 outside the approved additive delta.

1. **High — auto-fix.** `internal/coverage/coverage.go:437` returns exit 1 for extraction-mode map violations, while the base returned the rendered table at exit 0. Spec line 25 permits only appended help and the named worktree help spellings.
2. **High — auto-fix.** A canonical mapped coverage table with zero data rows reaches the same block, emits repair busywork, and exits 1. Spec lines 17 and 28 require the honest `help[0]{cmd,why}:` empty result. Same repair target as finding 1.
3. **Medium — auto-fix.** QD6 at spec lines 38 and 60 requires checked-in old/new response pairs per changed surface/state; candidate tests hold only newly authored inline expectations. Add the required base fixtures and exact appended-block-only comparisons, including streams, exits, and the named worktree argv delta.
4. **Medium — auto-fix.** `.agents/skills/bench-craft-cli/SKILL.md:81` advertises checks for “unchecked coverage rows,” but `internal/coverage/coverage.go:442` has no row-level checked-state classifier and emits one action per mapped row. State the shipped map-level behavior precisely.
5. **Low — auto-fix.** QD3's red signal at `specs/axi-query-disclosure/spec.md:57` searches `axi-disclosure`, while the shipped check is `axi-query-registry`; reconcile the acceptance token with the live check.

## Coverage

Count: 2. Worst: anchors' deep-CWD conformance case never reaches its cwd-sensitive branch.

1. **Medium — auto-fix.** `cmd/bench/command_registry_test.go:179` invokes `.bench/BENCH.md` without materializing it in the hermetic repository. Both root and nested runs fall through `cmd/bench/main.go:185` after `Lstat` fails, so a broken repository-relative normalization remains green. Create the file and assert the branch's real result.
2. **Medium — auto-fix.** `cmd/bench/command_registry_test.go:296` uses a historical no-map fixture for coverage empty state; no test exercises a canonical mapped table with zero rows. Add that partition with the Spec findings 1–2 repair.

Refuted/no-op: multiple malformed coverage rows correctly dedupe to one byte-identical retry template under spec line 27; valid per-row carried values are independently covered.
