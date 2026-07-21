# FT87 slice 1 implementation review

## Standards

3 findings. Worst: **High — bounds ownership enforcement is only a token tripwire.**

- **High:** `internal/conformance/checks_test.go:47-82` uses `strings.Contains` for expected tokens, so a caller can keep a dead `bounds.*` reference while redeclaring the actual policy value. This violates `AGENTS.md:25-35`'s one-source-per-fact rule. Add enforcement plus a duplicate-owner mutation/canary that proves the check bites.
- **Medium:** `internal/models/models.go:123,153` and `internal/worktree/refresh.go:53` duplicate the owned `10s`, `5 MiB`, and `30s` values in observable descriptions. Derive the descriptions from `internal/bounds` so enforcement and advertisement share one fact, per `AGENTS.md:25-35`.
- **Medium:** The user-visible offline, refresh, bounded-output, and timeout changes have no Unreleased entry in `CHANGELOG.md`. `craft-synthesis` requires one concise typed changelog entry for user-visible behavior.

The ambient `GIT_TERMINAL_PROMPT=1` concern was refuted by a focused run: the child observed `0`.

## Spec

4 findings. Worst: **High — shift refresh does not affect the worktree start ref, and the workflow offline proof is absent.**

- **High:** `internal/shift/loop.go:155-168` resolves `base=HEAD` before refresh and acquires the worktree at that pre-fetch SHA. This violates `specs/ft87-bounds-core.md:133-139`, which requires refresh before the start ref is resolved. With local A and refreshed origin B, the worktree still starts at A.
- **High:** `scripts/smoke-offline.sh:149-220` exercises only version, commands, and help under `BENCH_OFFLINE=1`; it never runs models or requested refresh and therefore cannot prove the five suppressed operation classes required by acceptance row `specs/ft87-bounds-core.md:247`.
- **Medium:** `internal/conformance/checks_test.go:47-82` checks names and tokens, not the locked numeric values or duplicate caller owners required by acceptance row `specs/ft87-bounds-core.md:245`.
- **Medium:** `internal/contract/runtime/runtime_worktree_test.go:820` does not provide rows `250-253`'s three-surface fetch-marker matrix, and `internal/gate/runner_test.go:26` does not assert timeout descendant death or row `267`'s healthy near-boundary behavior.

All 26 rows, eight stories, locked numerals, and slice-2/3 exclusions were audited. The duplicate environment-key concern was refuted.

## Coverage

5 findings. Worst: **High — duplicate policy owners and shipped offline attempts can escape the current tests.**

- **High:** `internal/conformance/checks_test.go:47-82` has no mutation where a caller redeclares a bound while retaining the expected token. Add the mutation required by `specs/ft87-bounds-core.md:245`.
- **High:** `scripts/smoke-offline.sh:149-220` and `internal/conformance/native_workflow_test.go:185` do not execute or sentinel models and refresh under offline mode. Add the operation-enumerating workflow proof required by `specs/ft87-bounds-core.md:247`.
- **Medium:** `internal/worktree/refresh_test.go:11-58` tests `Refresh` directly but not worktree subshell, create, and shift routing. Add the default/explicit/failure/offline acquisition-form marker matrix required by rows `250-253`, including hostile ambient `GIT_TERMINAL_PROMPT=1`.
- **Medium:** `internal/models/models_test.go:106,194` does not prove the Codex live and bundled attempts share one partially consumed deadline. Add a real shared-budget cancellation test for row `255`.
- **Medium:** `internal/gate/runner_test.go:26-49` asserts timeout state but not descendant death; the existing PID test covers ordinary cancellation. Add timeout-specific child-PID and healthy near-boundary coverage for rows `265` and `267`.

The full acceptance map, changed and existing tests, canonical edge classes, and the project hostile-input checklist were examined.
