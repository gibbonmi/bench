# Share Git command projections

Blocked by: none
Writes: internal/gitguard/projection_contract_test.go (new), internal/gitguard/scan.go, internal/shellcommand/projection.go (new), internal/shellcommand/projection_test.go (new)
Covers: GP01, GP02, GP03, GP05, GP06, GP07, GP08, GP09, GP10

## What to build

Deliver shared shell-wrapper and worktree-child facts through the existing Git classifier with its current verdicts.
Add ShellChildren and WorktreeExec in shellcommand, preserving Parse, ProjectCommandWords, and ResolveRoutinePrefix.
ShellChildren returns matching c-flag children in argv order without recursion.
WorktreeExec preserves the direct child argv and reports head, delimiter, and nonempty-child context separately.
Git retains keyword trimming, its first matching shell child, one-level traversal, executable recognition, global-option interpretation, and repository authority.
Remove its private wrapper and exec-child recognizers once these actual callers migrate.

Before this ticket starts, the complete accepted degraded-guard-refusal outcome must have landed.
Independently review and preserve that pre-migration executable after the repair and before GP1 changes source.
The coordinator retains its identity and exact observable reference tuples for GP3's final differential witness.
This prerequisite and reference capture are future implementation work, not observed evidence in this plan.

The Git-only W4 selection witness closes now: a first benign shell child followed by a dangerous Git child retains the current first-child result.
Nested wrapper and exec combinations remain one level deep, and child push facts cannot use the parent's checkout.
GP04's complete Git-and-Bench predicate closes in GP2 after its Bench consumer migrates.
GP2 receives a reviewed shared projection API with all Git-side tests green while its own consumer remains unbuilt.

Keep all existing assertions and independently authored expectations effective.
Do not rewrite old tests to match the new owner or derive expected results from the candidate producer.
For each required omission, pin the landed source and exact diagnostic before running it.
Accept only a compiling behavioral red, byte-identical restoration, and the same focused green afterward.
Invalid, compilation-only, skipped, or restore-failed proof closes no obligation.

Every first-use caller, fixture, registry, and headroom obligation closes in its introducing checkpoint.
A final family audit cannot supply a missing earlier test or repair.
Co-owned registry files may stay unchanged when their existing bindings and assertions suffice.
No new scanner, injected port, count expectation, structure grant, or policy authority enters this slice.

## Acceptance

- [ ] Existing routine-prefix tests retain executable indexes and nonexecuting query forms.
- [ ] ShellChildren preserves sh/bash/zsh basenames, supported c flags, argv order, invalid indexes, and missing operands without recursion or panic.
- [ ] WorktreeExec preserves quoted child elements byte-for-byte and separates missing delimiter, empty child, and nonempty DifferentCheckout.
- [ ] The actual Git classifier preserves its first matching wrapper child, including the benign-first W4 case, before Bench migration exists.
- [ ] Existing nested wrapper/exec verdicts and Git option forms retain the one-level contract and unresolved child or redirected push destination.
- [ ] Omitting DifferentCheckout propagation at the Git consumer compiles and reds TestGitProjectionAuthority without changing the shared fact owner.
- [ ] Omitting a supported wrapper basename or c-flag projection separately reds TestShellChildren; exact restoration returns the same focused green.
- [ ] No private Git wrapper or exec-child derivation remains after its callers migrate; final all-package enforcement waits for GP3.
- [ ] Existing Git policy, lexer failure, parse advice, and guard core entry assertions remain effective.
- [ ] New files stay within 400 lines and the modified scan file retains headroom in this checkpoint.
- [ ] The source-backed reference is captured after degraded repair and before migration, without candidate-derived expectations.
- [ ] Git and projection obligations pass while tickets 2 and 3 remain unbuilt.
