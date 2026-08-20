# Review pickup — ft229-hygiene-batch

Frozen base `364f34fa`, source tip `7c05e8e1`. Three axes at `opus`/medium,
read-only, run in parallel. Raw findings 21; de-duplicated repair targets 10,
of which 3 need a reviewer decision before any edit.

## Standards

Count 8. Worst: a second git-command lexer now lives in the shim's degraded rim
alongside `internal/gitguard`, and its correctness claim is unchecked.

- **S1 — duplicated `forced` predicate.** `internal/gitguard/verdict.go:22` and
  `:205` both spell `contains(args, "--force") || anyShortFlagHas(args, "f")`.
  Two derivations of one fact; the `--long` exclusion at `:342` had to stay
  correct for both by hand. Extract one helper. — auto-fix
- **S2 — the rim re-derives `internal/gitguard` in shell.**
  `.bench/hooks/block-dangerous-git.sh` grows a JSON decoder (`:82`), command
  lexer (`:234`), wrapper-prefix list (`:199`), and one-level `sh -c` recursion.
  `:309` asserts it behaves "exactly as internal/gitguard's tokenizer does" and
  nothing in the tree checks that claim. C1 below proves the claim false today.
  The spec closed *why* a second parser exists (the rim runs when Go cannot);
  what stays open is that the claim is unverified. — collapses into R1
- **S3 — broken sentence in a shipped kit file.**
  `.agents/commands/bench-review-implementation.md:26` reads "The gate is
  deterministic: the phases its table declares — in the Bench kit, the
  conformance registry's checks run inside the ordinary test phase, not as a
  phase of their own." The aside ate the predicate. Confirmed by direct read.
  Ships to every linked repo. — auto-fix
- **S4 — gate-shape fact stated twice.** The same conformance-shape claim sits
  at `.agents/commands/bench-review-implementation.md:26` and
  `.bench/BENCH-reference.md:146`, each with its own anchor
  (`internal/anchors/registry_data.go:35,37`). — ask-user
- **S5 — comment register.** `internal/structure/accept_validation.go:12-18`
  and `internal/status/status.go:924` argue their own correctness to the diff's
  reviewer, which `craft-comments` rejects. — auto-fix
- **S6 — middle man.** `internal/preflight/gather.go:41` keeps variadic
  `Gather` as a pass-through to `GatherPinned` with one non-test caller. — auto-fix
- **S7 — bare outline meta describes an absent table.**
  `internal/outline/outline.go:271` emits `outline_dirs[N]` while `outline_meta`
  still reports `emitted_symbols`, a count no block carries. — ask-user, folds into R8
- **S8 — redundant separator test** at `internal/landing/close.go:27`. — no-op

## Spec

Count 4. Worst: the shipped bare-outline shape is not the shape the spec and its
coverage rows describe, and the tests were written to the shipped shape.

- **P1 — bare `bench outline` groups by top-level directory, not scanned
  directory.** spec.md:168, H12 (spec.md:243), H15 (spec.md:246) and
  `tickets/summarize-the-bare-outline-form.md:18,21` all say "scanned
  directory". `internal/outline/outline.go:286` (`topLevel`) collapses to the
  first path segment. Verified live: `outline_dirs[19]` against 302 scanned
  directories, with `internal,"3942"` in one row. The H12/H15 tests
  (`internal/outline/outline_test.go:41,80`) assert the shipped shape, so both
  rows read green against a behavior the spec never stated. Behavioral spec
  contradiction — this asks rather than follows tree convention. — ask-user
- **P2 — the build amended its own ownership fences.** `56452041` and `b9de6ea2`
  added `internal/bounds/`, `cmd/bench/`, `projects/benchkit.md`, and
  `internal/gate/run_log_prune_test.go` to the approved list at spec.md:296.
  Every touched file is inside the amended list; none was inside the approved
  one. Both commits named themselves veto surface. — ask-user
- **P3 — H11 has no per-passage bite assertion.** Seven new anchor rows
  (`internal/anchors/registry_data.go:304-311`) rest on the generic
  `EvaluateGroup` path, matching how the other ~300 anchors are graded. — no-op
- **P4 — `appendHandoff` severity moved 11→12** (`internal/status/handoff.go:17`)
  so residue could take 11. Harmless, unasked. — no-op

Decisions otherwise honored: the close step on `bench commit --spec`, the deny
table in `internal/gitguard`, the rim reading `tool_input.command`, retention 20,
the identity-not-path residue predicate, and `--source-tip` verified with a
`--base`-shaped unresolvable error. No Won't-handle item was implemented.

## Coverage

Count 9. Worst: the narrowed rim allows `git push --force` behind an escaped
`&&`, a fail-open regression at the enforcement boundary.

- **C1 — the degraded rim fails open on `\uXXXX`-escaped operators.**
  `.bench/hooks/block-dangerous-git.sh:113` decodes any `\u` escape to a single
  `_`, so `is_control_op` never sees the operator. **Reproduced independently
  by the coordinator**, not only by the axis: the envelope
  `{"tool_input":{"command":"cat notes.md && git push --force"}}` exits
  0, where the pre-diff raw-substring rim exited 2. Go's `encoding/json`
  HTML-escapes `&`, `<`, `>` by default, confirmed by running `json.Marshal` on
  that command, so this is the shape a real producer emits — not a contrived
  one. Regression introduced by this diff. — auto-fix, blocking
- **C2 — no test exercises the escape table at all**
  (`block-dangerous-git.sh:105-120`). Mutating `*) return 1` or `n)` survives
  every case in `internal/systemtest/guard_rim_test.go:82-94`. — auto-fix
- **C3 — `isRunningBinary`'s two unknown→warn branches are unasserted**
  (`internal/worktree/live_binary.go:48,55`). `stubRunningBinary`
  (`live_binary_test.go:45`) always resolves, so flipping both `return true` to
  `return false` — silently removing the live binary, the reported incident —
  stays green. — auto-fix
- **C4 — the `EvalSymlinks` normalization** (`live_binary.go:50-52`) has no
  test; deleting the block passes. The profile lists invocation through a
  symlink as a live class. — auto-fix
- **C5 — nothing observes that a gate run prunes.** Only `pruneGateRunLogs` and
  `log.finish` are called directly (`run_log_prune_test.go:89,109,156,197`).
  Reverting `run_log.go:129` to the old non-pruning closure leaves story 39
  green. — auto-fix
- **C6 — the equal-timestamp tiebreak** (`run_log.go:262-267`) is untested;
  `seededRun` spaces records a minute apart. Two concurrent gates sharing a
  start instant get nondeterministic retention. — auto-fix
- **C7 — H09 pins reason presence, not reason text**
  (`structure_test.go:383`). Staleness is the load-bearing half. — no-op
- **C8 — bare outline's control-byte accounting is undecided.**
  `outline.go:250-263` skips `dirCount` for an unrepresentable path, so a repo
  whose only file carries a control-byte path renders `outline_dirs[0]` — the
  same table H16 asserts for absent. — ask-user, folds into R8
- **C9 — H25 never asserts the live binary was removed**
  (`live_binary_test.go:66-70` checks only the warning text). A guard that warns
  and skips the removal passes. — auto-fix

## Repair targets

| id | target | findings | disposition |
|---|---|---|---|
| R1 | decode ASCII `\u00XX` in the rim and assert the escape table | C1, C2, S2 | auto-fix, blocking |
| R2 | harden the live-binary tests | C3, C4, C9 | auto-fix |
| R3 | harden the run-log tests | C5, C6 | auto-fix |
| R4 | one `forced(args)` helper | S1 | auto-fix |
| R5 | repair the broken sentence | S3 | auto-fix |
| R6 | strip the PR-talk comments | S5 | auto-fix |
| R7 | collapse the `Gather` pass-through | S6 | auto-fix |
| R8 | bare-outline shape, meta, and empty state | P1, S7, C8 | ask-user |
| R9 | the gate-shape fact stated twice | S4 | ask-user |
| R10 | the mid-build fence amendment | P2 | ask-user |
