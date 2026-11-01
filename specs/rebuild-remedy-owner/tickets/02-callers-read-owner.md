# Read the rebuild remedy from the owner in the live-binary warning and the binary-seal row

Blocked by: 01-pool-rebuild-owner.md
Writes: internal/worktree/live_binary.go, internal/worktree/live_binary_test.go, internal/worktree/clean.go, internal/worktree/land_broker_notice_test.go, internal/preflight/binary_seal.go, internal/preflight/binary_seal_test.go, internal/preflight/decision.go, internal/preflight/gather.go, CHANGELOG.md, .bench/lib/resolve-bench.sh
Covers: RB13, RB14, RB15, RB16, RB23, RB24, RB25

## What to build

The live-binary warning and the `binary-seal` row read the remedy and the root kind from the owner that ticket 01 adds, so neither caller decides the root kind a second way.
Read these spec sections before you start: "Implementation decisions" > "The callers", "The `binary-seal` row", and "Records", the acceptance coverage map, and the edge inventory.
This ticket consumes `freshness.Rebuild`, `freshness.RebuildIn`, and `freshness.CheckoutRebuildAction` from ticket 01. Read their definitions in the tree before you call them.

Make these changes:

1. In `internal/worktree/live_binary.go`, make `warnBeforeRemovingLiveBinary` print `freshness.RebuildIn(a.home, root).Action`. `discardIgnored` already passes the ambient `a` (`internal/worktree/clean.go:294-304`), so no edit is expected in `clean.go`. The path stays on `Writes:` under the keep rule, because the slicing research named it as the caller.
2. In `internal/worktree/live_binary_test.go`, make `TestResidueGuardWarnsBeforeRemovingTheLiveBinary` expect the literal `bench worktree build '<creation path>'` and no `cd ` and no `scripts/go-build.sh` (RB13). The fixture home differs from the test-process home, so the same test grades RB14. Update the H25 comment, which today says that the warning names the `scripts/go-build.sh` invocation.
3. In `internal/worktree/land_broker_notice_test.go`, make the kit-checkout expectation and the installed-kit negative read `freshness.CheckoutRebuildAction(f.root)` (RB25). The landing code does not change.
4. In `internal/preflight/binary_seal.go`, make `binarySealFacts(home, root string)` also return `freshness.RebuildIn(home, root)`. Make `binarySealCheck` branch on that rebuild: a `Worktree` rebuild gives `red(Reason())` with `Next` set to the rebuild `Action`, and any other rebuild gives `red(Error())` with an empty `Next`. The check no longer reads `AssignmentTarget`, and the file no longer imports `usage`. Update the comment above `binarySealCheck`.
5. In `internal/preflight/decision.go`, add the `BinarySealRebuild freshness.Rebuild` field beside `BinarySealRefusal`, with a one-line doc comment. In `internal/preflight/gather.go`, pass `benchhome.Dir()` to `binarySealFacts` and set the new field. No production file in `internal/preflight` imports `benchhome` today, so this edit adds the import and the one ambient read of the home in the gather. `bench structure` grades `gather.go` within 5 lines of the 400-line maximum, so its net growth stays at 3 lines or fewer.
6. In `internal/preflight/binary_seal_test.go`, update the three `binarySealFacts` calls (`:187`, `:215`, `:242`) to the new signature. Rewrite `TestBinarySealRemedyNamesTheWorktreeBuildVerb` for RB15 (a `Worktree` rebuild fact with an empty `AssignmentTarget`), RB16 (a non-worktree rebuild fact), and RB23 (`AssignmentTarget` set with a non-worktree rebuild fact). Add the RB24 case to `TestBinarySealFactsGradeARootOnDisk`: a stale seal at `<home>/worktrees/k/x` gives a `Worktree` rebuild with `Action` equal to the literal `bench worktree build '<root>'`.
7. In `CHANGELOG.md`, add one entry to `## [Unreleased]` that names the pool-root remedy and the `binary-seal` next step.
8. In `.bench/lib/resolve-bench.sh`, reword only the clause of the comment above `bench_shell_quote` that says the guard-off hook blocks the worktree build verb. With an unreachable core, the hook allows every call, so the build verb cannot run there. With a stale core, the hook blocks every Bench call. This step is the RB-C1 review finding `rb-c1-spec-1`. Do not change a function body.

The slicing scratch run of steps 1, 4, and 5 observed exactly these reds before the test edits: `internal/preflight` fails to build at `binary_seal_test.go:187,215,242`, and `TestResidueGuardWarnsBeforeRemovingTheLiveBinary` reds on the home mismatch. `cmd/bench`, `internal/preflight/staleness`, and the rest of `internal/worktree` stayed green.

After each Go edit, run `bench worktree build <target>` before a `bench <verb> --in <target>` call (`docs/known-issues.md:10`, `docs/known-issues.md:12`).

## Acceptance

- [ ] A release of a pool worktree whose `dist/bench` is the running binary warns with `bench worktree build '<creation path>'`, with no `cd ` and no `scripts/go-build.sh`, when the test-process home differs from the fixture home.
- [ ] The `binary-seal` row of a pool root prints the refusal reason with `Next` equal to the owner's build remedy, and no cell holds `cd ` or `scripts/go-build.sh`.
- [ ] The `binary-seal` row of any other root prints the whole refusal with an empty `Next`, also when `AssignmentTarget` is set.
- [ ] `binarySealFacts` with an explicit home returns the owner's build remedy for a stale seal under that home's pool.
- [ ] The broker notice for the kit source checkout still holds the `cd` form and `bench doctor --fix`.
- [ ] `CHANGELOG.md` holds one `## [Unreleased]` entry that names the pool-root remedy and the `binary-seal` next step.
