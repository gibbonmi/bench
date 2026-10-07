# Refuse malformed guard commands with repair advice

Blocked by: none
Writes: .bench/hooks/block-dangerous-git.sh, internal/gitguard/gitguard.go, internal/gitguard/unlexed_test.go, internal/systemtest/guard_refusal_test.go (new), tests/canary/package-core-guard/guard-describe-boundary-dropped, tests/canary/package-core-guard/guard-resolver-order-drift, tests/canary/canonical-path-owner/second-derivation, tests/canary/injected-ports/unregistered-port, internal/anchors/registry_data.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: GR01, GR02, GR03, GR04, GR05, GR06, GR07, GR08, GR09, GR10, GR11, GR12, GR13, GR14, GR15, GR16, GR17, GR18, GR19, GR20, GR21, GR22, GR23, GR24, GR25, GR26, GR27, GR28, GR29, GR30, GR31, GR32, GR33, GR34, GR35, GR36, GR37, GR38, GR39, GR40, GR41

## What to build

Deliver the complete GR1 outcome through the real Git hook and the reachable core entry.
The degraded hook refuses a command when its parser cannot establish the supported command structure.
Its supported wrapper scan propagates the same refusal for a malformed child.
The degraded refusal identifies the parse failure and gives quote-and-escape repair advice.

The core keeps its existing classification and envelope policy.
Its unlexed Git class receives the same repair advice without the destructive-command authority explanation.
Other core labels retain their authority messages and advice.
The shell remains a separate implementation because it must work without the core.

The concrete scenario sends `printf hi;git push --force\necho 'x` through a real missing-core hook.
The `\n` spelling denotes one newline byte.
The hook must exit 2 with a parse-repair diagnostic.
The same envelope reaches the live-core hook and the direct selected binary's `guard-git` entry.
Both entries must retain refusal and emit parse-only advice.

Retain the coarse refusal for recognized Git and the supported non-Git recovery path.
Retain the old wrapper depth, operator, decoder, envelope, Git-option, and worktree-child cases.
Do not widen the reachable core's malformed non-Git or unreadable-envelope policy.
The grammar migration, scanner simplification, and future `rg` operand guard remain separate outcomes.

Read the source definitions and the cited existing tests before implementation.
Use the existing system owner and `privateToolPath` helper without copying their harness.
The tagged system suite receives BENCH_KIT and the sealed selected binary from `bench test --check system`.
Do not start an unsealed system binary.

## Acceptance

- [ ] The missing-core hook exits 2 for `printf hi;git push --force\necho 'x` (GR01).
- [ ] The missing-core hook exits 2 for `printf hi;git push --force\necho "x` (GR02).
- [ ] The missing-core hook exits 2 for `printf hi;git push --force\nprintf x\` (GR03).
- [ ] The missing-core hook exits 2 for `echo 'x` (GR04).
- [ ] The missing-core hook exits 2 for `bash -c "echo 'x"` (GR05).
- [ ] The malformed degraded refusal contains the quote-and-escape repair sentence (GR06).
- [ ] The malformed degraded refusal excludes the authority preamble (GR07).
- [ ] The malformed degraded refusal excludes `Stop and hand back` (GR08).
- [ ] The malformed degraded refusal excludes the history-rewrite explanation (GR09).
- [ ] The malformed degraded refusal identifies the parse failure after `BLOCKED:` (GR10).
- [ ] The core unlexed refusal contains the existing quote-and-escape repair sentence (GR11).
- [ ] The core unlexed refusal excludes the authority preamble (GR12).
- [ ] The missing-core hook refuses the existing recognized Git fixtures (GR13).
- [ ] The missing-core hook exits 2 for well-formed `git status` (GR14).
- [ ] The missing-core hook permits the existing non-Git recovery fixtures (GR15).
- [ ] The missing-core hook preserves the existing nested-wrapper exit 0 (GR16).
- [ ] The missing-core hook preserves the existing operator-run verdicts (GR17).
- [ ] The missing-core hook preserves the existing escaped-envelope verdicts (GR18).
- [ ] The missing-core hook refuses each existing unreadable-envelope fixture (GR19).
- [ ] The missing-library hook exits 2 with its existing missing-library diagnostic (GR20).
- [ ] The hook exits 2 when its reachable core returns exit 3 (GR21).
- [ ] The core-error refusal retains the analyzer-error diagnostic (GR22).
- [ ] The core preserves the existing unlexable-Git labels (GR23).
- [ ] The core permits the existing malformed non-Git fixtures (GR24).
- [ ] The core force-push message retains its authority advice (GR25).
- [ ] The live-core hook returns a parse-only refusal for the occurrence envelope (GR26).
- [ ] The selected binary returns a parse-only refusal through `guard-git` (GR27).
- [ ] The core retains the existing normal push allow and block verdicts (GR28).
- [ ] The real missing-core fixture observes its verdict without a reachable Bench executable (GR29).
- [ ] The hook exits 2 for the occurrence envelope when its wrapper returns exit 127 (GR30).
- [ ] The missing-core hook exits 2 for well-formed `printf hi;git status` (GR31).
- [ ] The core unlexed refusal excludes `Stop and hand back` (GR32).
- [ ] The core unlexed refusal excludes the history-rewrite explanation (GR33).
- [ ] The core unlexed refusal retains its existing parse label after `BLOCKED:` (GR34).
- [ ] The hook permits `cat .github/workflows/gate.yml` when its wrapper returns exit 127 (GR35).
- [ ] The core classifies `git -C /tmp reset --hard` as `git reset` (GR36).
- [ ] The core refuses `bench worktree exec X -- git commit -m x` (GR37).
- [ ] The core permits `bench worktree exec X -- git status` (GR38).
- [ ] The core classifies the exec-child push as unresolved (GR39).
- [ ] The core refuses `bash -c 'git push'` (GR40).
- [ ] The core permits `bash -c 'sh -c "git push"'` (GR41).

## Checkpoint verification

The ticket completes production behavior, tests, and both executable entries in one green checkpoint.
No successor supplies a missing producer or test seam.
Record each coverage row's red-to-green evidence in the review pickup.
Name the exact fixture input and the behavioral mutation for each independent expectation.
An invalid probe or a compile failure proves no acceptance row.

Run these future checks from the candidate integration source:

- `bench test --package ./internal/gitguard`
- `bench test --package ./cmd/bench`
- `bench test --check system`
- `bench test --check package-core-guard`
- `bench test --check canonical-path-owner`
- `bench test --check injected-port-registry`
- `bench test --check docs-currency-workflow`
- `bench test --check axi-query-registry`
- `bench test --check subcommand-routing`
- `bench coverage --check specs/degraded-guard-refusal/spec.md`

The system fixtures prove missing core, missing binary, missing library, core error, and reachable-core entries.
The new system file contains the planned tests that the spec names.
Every launch uses the existing system owner's process ledger.
A missing-core fixture must prove that no Bench executable is reachable on disk or PATH.

Run the author mutation through `bench probe` at the malformed trailing-escape branch.
Restore permissive parse success for the exact GR03 command.
GR03 must turn red through the real degraded hook.
Report the mutation kind, site, observed row, and restored subject.

The coordinator probes a different site after the author returns.
It discards the malformed child's status at the wrapper recursion site.
GR05 must turn red through the real degraded hook.
The coordinator records its separate restored subject.

Run the duplicate-facts sweep before the ticket commit.
Keep the core repair sentence in its existing shared owner.
Keep the shell repair sentence only in its independent degraded implementation.
Preserve the fixture helpers and the invocation registries outside this ticket's write expectation.

The four canary directories and anchor registry join because their existing declarations pin this ticket's production paths.
Keep their fixture and anchor assertions intact.
The anchor package binds five command registries through the existing ticket binding table.
Their entries join the write expectation without a command or help change.

- Preserve `this is an honest-mistake layer, not an evasion-resistant`.
- Preserve `exactly one level deep by design (see internal/gitguard)`.
- Preserve `misaligned agent are the git pre-push hook and bench's pooled-worktree`.

Do not remove a pin or weaken an assertion to make the lane green.
If a closure proposal requires another file, report it before any out-of-fence edit.
The coordinator owns any approved plan expansion.

## Review checkpoint

This ticket is the sole member of GR1.
Freeze its predecessor and committed tip for Standards, Spec, and Coverage review.
Reconcile GR01-GR41 against the complete approved spec after that checkpoint.
The review pickup is `reviews/degraded-guard-refusal.md`.
The coordinator owns that pickup and the final acceptance record.
