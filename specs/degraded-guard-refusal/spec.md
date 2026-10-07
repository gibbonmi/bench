# Refuse malformed commands in the degraded Git guard

Status: staged
Decision source: `decisions/architecture-guards/tickets/2.md` in `ft362-process-lifetime-map`
Subject: `37ac80b8adeff4ecbd57f3d813ff8937f2bcbc68`
Audience: every repository that links the kit
Related roadmap: FT366
Verification log: one spec-review iteration accepted before ticket slicing; ticket review pending

## Problem

An unclosed quote causes the degraded Git hook to use a permissive whitespace split.
An attached separator can then hide a Git invocation.
The reachable core refuses that command but adds unrelated authority advice.

## Solution

Refuse commands when the degraded parser cannot establish the supported command structure.
Give a parse-repair diagnostic for that refusal and the core's existing unlexed Git class.
Retain the coarse Git refusal and the supported non-Git recovery path.

## User stories

Line: `gpt-6.1-sol` / high
Implementation-line reason: the shell must propagate parse failure through recursion, while the real-hook fixtures provide a resolved test seam
Harder chunks: GR1


1. As an agent, I want the degraded Git guard to refuse malformed commands, so that partial commands cannot bypass the guard.
2. As an agent, I want a parse-repair diagnostic, so that I can repair the command without a false authority warning.
3. As an agent, I want malformed wrapper children to refuse, so that the supported wrapper depth cannot hide a parse failure.
4. As an agent, I want valid Git commands to refuse without the core, so that the coarse guard retains its protection.
5. As an agent, I want valid non-Git recovery commands to pass without the core, so that I can restore the core.
6. As an agent, I want the guard to distinguish arguments from executable words, so that recovery reads remain available.
7. As an agent, I want the existing wrapper depth, so that this repair does not widen the shell interpreter.
8. As an agent, I want unreadable degraded envelopes to refuse, so that absent commands cannot grant access.
9. As an agent, I want a missing shared library to refuse, so that a hook error cannot permit the command.
10. As an agent, I want a core error to refuse, so that a failed analyzer cannot grant access.
11. As an agent, I want the core to retain its malformed-command verdict, so that the diagnostic repair does not relax refusal.
12. As an agent, I want destructive-command authority advice to remain, so that syntax advice does not replace policy advice.
13. As an agent, I want the core's non-Git uncertainty rule, so that this degraded repair does not change another input policy.
14. As a maintainer, I want tests through the real hook and core entries, so that an internal helper cannot mask a broken route.
15. As a maintainer, I want the existing Git options and worktree-child policy, so that a diagnostic change cannot alter command classification.
16. As a maintainer, I want separate shell and core implementations, so that the recovery guard works without the core.

## Implementation decisions

The intended shell result distinguishes recognized Git, supported non-Git, and an unparseable command.
Only the supported non-Git state permits execution without the core.
The parser reports failure for unclosed single quotes, unclosed double quotes, and trailing escapes.
The existing one-level wrapper scan propagates a child parse failure.

The shell retains its separate implementation because it runs without the core.
The reachable core retains its current classification and envelope policy.
Only its existing unlexed refusal receives the parse-only message.

The core's unlexed message reads the repair sentence from its existing shared owner.
The shell owns its independent repair text.
Tests through both entries justify that duplication.
This repair adds no shared parser, command, or runtime dependency.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
|---|---|---|---|---|
| GR1 / 01-refuse-malformed-guard-commands.md | The Git guard refuses malformed degraded input with parse-repair diagnostics | GR01-GR41 | The real hook, core entry, and Git message tests | yes |

This is one behavior outcome with one independent implementation review checkpoint.
Independent spec acceptance precedes this one-ticket slice.
A future ticket author gets three iterations on the declared implementation line.
The coordinator reconciles the complete outcome after its checkpoint.

## Testing decisions

The external behavior is the hook exit code and its stderr diagnostic.
The hook must refuse malformed input before the caller executes the proposed command.
The hook itself does not execute that command.
This spec adds no executable-authentication guarantee.

Use the existing system owner and `privateToolPath` from `internal/systemtest/guard_rim_test.go`.
Add system tests in `internal/systemtest/guard_refusal_test.go` instead of copying the fixture harness.
A private PATH contains only the hook dependencies, with no reachable Bench wrapper or core.
Run the hook outside a repository for the missing-core fixtures.

Use a temporary hook layout for the missing-library fixture.
Copy the current hook into that layout and omit its adjacent shared library.
Use a temporary repository with a local executable wrapper that returns exit 3 for the core-error fixture.
A wrapper that returns exit 127 reaches the missing-binary rim through the same fixture owner.

For the live-core fixture, the local executable resolves to the system owner's sealed selected binary.
Send a real PreToolUse envelope through the unchanged hook entry.
Also call the selected binary's `guard-git` entry directly.

The concrete occurrence input is `printf hi;git push --force\necho 'x`.
The attached semicolon precedes an executable Git word.
The `\n` spelling denotes one newline byte.
The final quote remains open.
The current degraded fallback joins `hi;git` and loses that executable word.

The parse diagnostic starts with `BLOCKED:` and identifies a parse failure.
It contains `Close every quote and escape in the command, then run it again.`
It excludes `you don't have authority`, `Stop and hand back`, and the destructive-command history explanation.
The core's `unlexed` class retains its current label.
Other core labels retain their current messages and advice.
The shell owns its repair text because it must work without the core.


### Seam diagram

    trigger: PreToolUse Bash envelope
        │
        ▼
    block-dangerous-git.sh
        ├── absent core ──▶ shell parser ──▶ allow / Git refusal / parse refusal
        ├── reachable core ──▶ guard-git ──▶ gitguard ──▶ allow / refusal
        └── core error ──▶ existing fail-closed refusal

    tests attach at the real hook entry and the selected binary's guard-git entry

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| GR01 | 1, 14 | The missing-core hook exits 2 for `printf hi;git push --force\necho 'x` | planned TestDegradedGuardRefusesMalformedCommands in internal/systemtest/guard_refusal_test.go | The current fallback permits the occurrence input |
| GR02 | 1 | The missing-core hook exits 2 for `printf hi;git push --force\necho "x` | planned TestDegradedGuardRefusesMalformedCommands in internal/systemtest/guard_refusal_test.go | A single-quote-only repair permits this input |
| GR03 | 1 | The missing-core hook exits 2 for `printf hi;git push --force\nprintf x\` | planned TestDegradedGuardRefusesMalformedCommands in internal/systemtest/guard_refusal_test.go | The current trailing-escape fallback succeeds |
| GR04 | 1, 5 | The missing-core hook exits 2 for `echo 'x` | planned TestDegradedGuardRefusesMalformedCommands in internal/systemtest/guard_refusal_test.go | A Git-word-only check permits this malformed recovery command |
| GR05 | 3, 7 | The missing-core hook exits 2 for `bash -c "echo 'x"` | planned TestDegradedGuardRefusesMalformedCommands in internal/systemtest/guard_refusal_test.go | A discarded child parse status permits this supported wrapper |
| GR06 | 2 | The malformed degraded refusal contains the quote-and-escape repair sentence | planned TestDegradedGuardRefusesMalformedCommands in internal/systemtest/guard_refusal_test.go | A generic degraded refusal omits the actionable repair |
| GR07 | 2 | The malformed degraded refusal excludes the authority preamble | planned TestDegradedGuardRefusesMalformedCommands in internal/systemtest/guard_refusal_test.go | Appending syntax advice to the old authority message fails |
| GR08 | 2 | The malformed degraded refusal excludes `Stop and hand back` | planned TestDegradedGuardRefusesMalformedCommands in internal/systemtest/guard_refusal_test.go | The authority termination cannot survive the parse diagnostic |
| GR09 | 2 | The malformed degraded refusal excludes the history-rewrite explanation | planned TestDegradedGuardRefusesMalformedCommands in internal/systemtest/guard_refusal_test.go | A partial preamble edit leaves unrelated policy advice |
| GR10 | 2 | The malformed degraded refusal identifies the parse failure after `BLOCKED:` | planned TestDegradedGuardRefusesMalformedCommands in internal/systemtest/guard_refusal_test.go | Advice alone cannot identify the refusal class |
| GR11 | 2, 11 | The core unlexed refusal contains the existing quote-and-escape repair sentence | `internal/gitguard/unlexed_test.go` (`TestBlockMessageCarriesUnlexedAdvice`) | The test detects advice removal |
| GR12 | 2, 11 | The core unlexed refusal excludes the authority preamble | planned TestBlockMessageSeparatesParseRefusal in internal/gitguard/unlexed_test.go | The current uniform BlockMessage prefix fails |
| GR13 | 4 | The missing-core hook refuses the existing recognized Git fixtures | `internal/systemtest/guard_rim_test.go` (`TestDegradedGuardRimDecidesFromTheCommandField`) | The existing destructive and wrapper cases detect a removed Git branch |
| GR14 | 4 | The missing-core hook exits 2 for well-formed `git status` | planned TestDegradedGuardRefusesMalformedCommands in internal/systemtest/guard_refusal_test.go | A destructive-only degraded policy permits read-only Git |
| GR15 | 5, 6 | The missing-core hook permits the existing non-Git recovery fixtures | `internal/systemtest/guard_rim_test.go` (`TestDegradedGuardRimDecidesFromTheCommandField`) | Blanket degraded refusal breaks recovery reads |
| GR16 | 7 | The missing-core hook preserves the existing nested-wrapper exit 0 | `internal/systemtest/guard_rim_test.go` (`TestDegradedGuardRimDecidesFromTheCommandField`) | Deeper automatic recursion changes the pinned depth |
| GR17 | 4, 6 | The missing-core hook preserves the existing operator-run verdicts | `internal/systemtest/guard_rim_test.go` (`TestDegradedGuardRimDecidesFromTheCommandField`) | Collapsing operators into words loses command boundaries |
| GR18 | 5, 6 | The missing-core hook preserves the existing escaped-envelope verdicts | `internal/systemtest/guard_rim_test.go` (`TestDegradedGuardRimDecidesFromTheCommandField`) | Bypassing the decoder changes literal and escaped operators |
| GR19 | 8 | The missing-core hook refuses each existing unreadable-envelope fixture | `internal/systemtest/guard_rim_test.go` (`TestDegradedGuardRimDecidesFromTheCommandField`) | Treating unreadable as empty permits these fixtures |
| GR20 | 9 | The missing-library hook exits 2 with its existing missing-library diagnostic | planned TestGitGuardMissingLibraryRefuses in internal/systemtest/guard_refusal_test.go | Sourcing a missing library without refusal grants access |
| GR21 | 10 | The hook exits 2 when its reachable core returns exit 3 | planned TestGitGuardCoreErrorRefuses in internal/systemtest/guard_refusal_test.go | Reinterpreting a core error as a parse result or permission fails |
| GR22 | 10 | The core-error refusal retains the analyzer-error diagnostic | planned TestGitGuardCoreErrorRefuses in internal/systemtest/guard_refusal_test.go | The parse diagnostic cannot replace an analyzer failure |
| GR23 | 11 | The core preserves the existing unlexable-Git labels | `internal/gitguard/unlexed_test.go` (`TestClassifyRefusesAnUnlexableGitCommand`) | Removing the existing fail-closed classification fails |
| GR24 | 13 | The core permits the existing malformed non-Git fixtures | `internal/gitguard/unlexed_test.go` (`TestClassifyRefusesAnUnlexableGitCommand`) | Applying the degraded policy to the core fails |
| GR25 | 12 | The core force-push message retains its authority advice | planned TestBlockMessageSeparatesParseRefusal in internal/gitguard/unlexed_test.go | A global prefix replacement weakens a destructive refusal |
| GR26 | 14 | The live-core hook returns a parse-only refusal for the occurrence envelope | planned TestGitGuardHookAndCoreParseRefusal in internal/systemtest/guard_refusal_test.go | A helper-only fix cannot mask stale entry wiring |
| GR27 | 14 | The selected binary returns a parse-only refusal through `guard-git` | planned TestGitGuardHookAndCoreParseRefusal in internal/systemtest/guard_refusal_test.go | A shell-only diagnostic fix leaves the direct core entry wrong |
| GR28 | 15 | The core retains the existing normal push allow and block verdicts | `cmd/bench/main_test.go` (`TestGuardGitBlockAllow`) | Changing classification while changing advice fails the old entry cases |
| GR29 | 16 | The real missing-core fixture observes its verdict without a reachable Bench executable | planned TestDegradedGuardRefusesMalformedCommands in internal/systemtest/guard_refusal_test.go | Accidentally reaching the core cannot prove the degraded repair |
| GR30 | 1 | The hook exits 2 for the occurrence envelope when its wrapper returns exit 127 | planned TestGitGuardMissingBinaryPreservesRecovery in internal/systemtest/guard_refusal_test.go | Fixing only unresolved-wrapper dispatch misses the missing-binary rim |
| GR31 | 4 | The missing-core hook exits 2 for well-formed `printf hi;git status` | planned TestDegradedGuardRefusesMalformedCommands in internal/systemtest/guard_refusal_test.go | Rejecting only malformed text loses the attached-separator Git case |
| GR32 | 2 | The core unlexed refusal excludes `Stop and hand back` | planned TestBlockMessageSeparatesParseRefusal in internal/gitguard/unlexed_test.go | Removing only the first authority clause leaves the old termination |
| GR33 | 2 | The core unlexed refusal excludes the history-rewrite explanation | planned TestBlockMessageSeparatesParseRefusal in internal/gitguard/unlexed_test.go | A partial prefix edit leaves unrelated policy advice |
| GR34 | 11 | The core unlexed refusal retains its existing parse label after `BLOCKED:` | planned TestBlockMessageSeparatesParseRefusal in internal/gitguard/unlexed_test.go | A generic syntax message loses the existing refusal identity |
| GR35 | 5 | The hook permits `cat .github/workflows/gate.yml` when its wrapper returns exit 127 | planned TestGitGuardMissingBinaryPreservesRecovery in internal/systemtest/guard_refusal_test.go | Blanket refusal at the missing-binary rim breaks the recovery path |
| GR36 | 15 | The core classifies `git -C /tmp reset --hard` as `git reset` | `internal/gitguard/scan_test.go` (`TestScanWrapperDepthAndGlobalOpts`) | A global option value cannot hide the destructive verb |
| GR37 | 15 | The core refuses `bench worktree exec X -- git commit -m x` | `internal/gitguard/verdict_test.go` (`TestClassifyVerdicts`) | A diagnostic change cannot discard the exec-child classification |
| GR38 | 15 | The core permits `bench worktree exec X -- git status` | `internal/gitguard/verdict_test.go` (`TestClassifyVerdicts`) | Blanket exec-child refusal loses ordinary recovery reads |
| GR39 | 15 | The core classifies the exec-child push as unresolved | `internal/gitguard/verdict_test.go` (`TestClassifyVerdicts`) | The target checkout cannot inherit the current checkout's push facts |
| GR40 | 7 | The core refuses `bash -c 'git push'` | `internal/gitguard/scan_test.go` (`TestScanWrapperDepthAndGlobalOpts`) | A removed wrapper branch loses the supported first level |
| GR41 | 7 | The core permits `bash -c 'sh -c "git push"'` | `internal/gitguard/scan_test.go` (`TestScanWrapperDepthAndGlobalOpts`) | Deeper recursion changes the core's pinned uncertainty rule |

### Edge inventory

The profile's hostile-input checklist applies to envelope bytes and shell command text.
The system fixtures use directory names with spaces and brackets through the existing owner.
The old tests distinguish empty commands from unreadable envelopes.
They cover literal and escaped operators, Git arguments, redirect targets, and shell wrappers.
GR01-GR05 cover both quote types, attached separators, trailing escapes, and malformed supported children.

No test substitutes a package variable across a subprocess boundary.
The missing-core fixture controls PATH and the working directory.
The missing-library and core-error fixtures control exact files inside temporary layouts.
This outcome changes neither cleanup behavior nor durable records.

- **Won't handle**: malformed non-Git input under the reachable core remains permitted by `gitguard.Classify`.
- **Won't handle**: a second wrapper level retains the depth fixture in `TestDegradedGuardRimDecidesFromTheCommandField`.
- **Won't handle**: shell expansions and arbitrary executable aliases remain outside `gitguard.scanWords` and `invokes_git`.
- **Won't handle**: full shell syntax validation remains outside the supported subset in `lex_command`.
- **Won't handle**: common envelope validation remains with `CommandFromEnvelope` and `bench_envelope_command` until the separate grammar outcome.
- **Won't handle**: the core retains its current unreadable-envelope allow case in `TestGuardGitBlockAllow`.
- **Won't handle**: new Git-option or worktree-child recognition remains outside this repair of `gitguard.scanWords`.
- **Won't handle**: TOON cells, Git patch paths, symbolic refs, and Markdown fields do not reach this seam.
- **Won't handle**: filesystem discovery, FIFOs, symlink policy, and transactional writes remain with the existing wrapper resolver and system owner.
- **Won't handle**: cancellation deadlines remain with the separate guard-inventory scan outcome.

## Ownership fences

The prospective build owns these exact paths:

- `.bench/hooks/block-dangerous-git.sh`
- `internal/gitguard/gitguard.go`
- `internal/gitguard/unlexed_test.go`
- `internal/systemtest/guard_refusal_test.go`
- `tests/canary/package-core-guard/guard-describe-boundary-dropped`
- `tests/canary/package-core-guard/guard-resolver-order-drift`
- `tests/canary/canonical-path-owner/second-derivation`
- `tests/canary/injected-ports/unregistered-port`
- `internal/anchors/registry_data.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `reviews/degraded-guard-refusal.md`

The ticket owns four implementation paths and ten closure paths.
The coordinator owns the review pickup.
The planning author owns this spec folder during the spec phase.
The spec fence equals the ticket write expectation plus the review pickup.
The closure paths join because their existing declarations pin the production paths.
The author keeps their fixture and anchor assertions intact.

## Out of scope

FT366 remains partial after this outcome.
This spec has no Roadmap retirement field.
These separate capabilities retain their source requirements:

- Shared command projection and envelope facts: 12 edits, 3 gate runs.
- A context-aware guard inventory scan: 3 edits, 2 gate runs.
- The `rg` operand guard and its permissive operand-uncertainty rule: 4 edits, 2 gate runs.

These are planning estimates, not approved priorities or implementation commitments.
The shared-grammar estimate includes the existing shellcommand, gitguard, benchguard, census, and worktree callers.
The scan estimate includes its owner, tests, and decision documentation.
The `rg` estimate includes its classifier, hook entry, entry tests, and invocation registry.
No `rg` guard exists in the subject tree's command registry.
Its roadmap requirement survives this spec.

## Further notes

Evidence status: source-backed proposal, with no runtime compatibility claim
Remaining uncertainty: future red-to-green evidence and the independent reviewer verdict

### Reader inventory

The build does not add a command, hook path, injected port, or fixture family.
The existing invocation registrations consume the same paths and exit codes.
The ten closure paths retain their current rules and assertion meanings.

The registry sweep includes these unchanged consumers:

- `cmd/bench/main.go:134`: the `guard-git` command registration.
- `.claude/settings.json:68`: the Claude hook invocation.
- `.codex/hooks.json`: the Codex hook invocation.
- `internal/conformance/entry_point_parity_test.go:56-58`: the hook/core parity registration.
- `internal/conformance/package_core_checks_test.go:22-35`: the resolver, header, and classifier checks.
- `internal/conformance/injected_ports_registry_test.go:54-55`: the unchanged Checker injection.
- `internal/conformance/tier_test.go:41`: the inventory scans only `internal/conformance` tests.

The build retains the existing system fixture helper and the core repair-sentence owner.
The future build must preserve each registry path listed above.

The fixture closure sources are the two hook overlay files and the two Git guard BASE records.
The anchor closure source is `internal/anchors/registry_data.go:403-405`.
The anchor package requires five command registries through `internal/tickets/registry_data.go:20-26` and line 45.
Those entries pin the threat model, wrapper depth, and named backstops.
The author read the fixture overlays, BASE records, mutation records, and EXPECT records before the fence amendment.

### Source coverage

The single reviewed input is decision ticket 2 in the retained `ft362-process-lifetime-map` assignment.
Its lines 12-14 require the coarse Git rim and non-Git recovery path.
GR13-GR19, GR14, and GR31 preserve those behaviors.
Its lines 16-19 require malformed-input refusal and parse-repair advice.
GR01-GR12, GR26-GR27, and GR32-GR34 cover those requirements.

Its lines 21-23 require independent degraded-entry tests and the existing dangerous-command cases.
GR13, GR29, and the separate shell/core implementations preserve those requirements.

Decision ticket 3:12-24 places this outcome before grammar and scan work.
GR36-GR41 preserve representative option, child, and wrapper verdicts through the existing tests.
Their shared ownership changes belong to the grammar outcome.
Decision ticket 1:21-23 retains the `rg` uncertainty distinction.
That future permissive rule cannot relax this malformed-command refusal.

The primary sources below establish the current behavior at the subject commit:

- `roadmap/FT366.md:12-17`: the degraded occurrence and incorrect core authority preamble.
- `.bench/hooks/block-dangerous-git.sh:91-102`: the coarse recovery verdict.
- `.bench/hooks/block-dangerous-git.sh:109-153`: the supported one-level child scan.
- `.bench/hooks/block-dangerous-git.sh:250-310`: the malformed-input fallback.
- `.bench/hooks/block-dangerous-git.sh:313-323`: the core dispatch and error rim.
- `internal/gitguard/gitguard.go:125-149`: the unchanged classifier and uniform message prefix.
- `internal/shellcommand/unlexed.go:5-7`: the core repair sentence.
- `internal/systemtest/guard_rim_test.go:21-169`: the actual degraded seam and its single fixture helper.
- `cmd/bench/guards.go:99-131`: the direct core entry and exit contract.
- `internal/gitguard/scan.go:11-101`: the existing depth, options, and worktree-child ownership.

A source edit to those functions invalidates these current-code claims.
The author inspected sources without executing production tests or a malformed command.
No runtime red is claimed during specification.

### Validation plan

The build first observes GR01 red through the real missing-core hook.
It then observes each independent new expectation red through a named omission or behavioral mutation.

The central author probe restores the old permissive malformed fallback at the trailing-escape branch.
GR03 must fail through `bench probe`.
The coordinator probe instead discards the malformed child's status at the wrapper recursion site.
GR05 must fail through the real degraded hook.
Each probe must report a restored subject.

Run `bench test --package ./internal/gitguard` for the core message and classifier cases.
Run `bench test --check system` for the sealed binary and real hook entries.
Run `bench test --check package-core-guard` for resolver and manifest compatibility.
Run the required lane checks before a future ticket commit.
The coordinator runs the complete gate only for the authorized future landing.

### Completion plan

The authored version 1 plan records verification obligations.
An approved implementation records its fresh author in a version 2 amendment before dispatch.
System verification uses BENCH_KIT and the sealed selected binary supplied by `bench test`.

```bench-completion-plan
{"version":1,"chunks":[{"id":"GR1","tickets":["01-refuse-malformed-guard-commands.md"],"verification":[{"id":"gitguard","command":"bench test --package ./internal/gitguard"},{"id":"core-entry","command":"bench test --package ./cmd/bench"},{"id":"hook-entries","command":"bench test --check system"},{"id":"hook-contracts","command":"bench test --check package-core-guard"},{"id":"author-parse-success-omission","command":"bench test --check system","probe":"restore permissive malformed success at the trailing-escape branch; GR03 must fail"},{"id":"canonical-path-policy","command":"bench test --check canonical-path-owner"},{"id":"injected-port-policy","command":"bench test --check injected-port-registry"},{"id":"guard-anchors","command":"bench test --check docs-currency-workflow"},{"id":"axi-registry","command":"bench test --check axi-query-registry"},{"id":"command-routing","command":"bench test --check subcommand-routing"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/degraded-guard-refusal/spec.md"},{"id":"gitguard","command":"bench test --package ./internal/gitguard"},{"id":"core-entry","command":"bench test --package ./cmd/bench"},{"id":"hook-entries","command":"bench test --check system"},{"id":"hook-contracts","command":"bench test --check package-core-guard"},{"id":"coordinator-child-status-omission","command":"bench test --check system","probe":"discard malformed-child status at wrapper recursion; GR05 must fail"},{"id":"canonical-path-policy","command":"bench test --check canonical-path-owner"},{"id":"injected-port-policy","command":"bench test --check injected-port-registry"},{"id":"guard-anchors","command":"bench test --check docs-currency-workflow"},{"id":"axi-registry","command":"bench test --check axi-query-registry"},{"id":"command-routing","command":"bench test --check subcommand-routing"}]}
```

### Review chronology

An independent xhigh review accepted the frozen spec before ticket slicing on 2026-10-06.
It found no material Standards, Spec, or Coverage findings at confidence 9.
Its accepted spec digest is `ded9e33e6e9841c2ddbea4b4098cef35c1255dbcda6c2147f568980c2bfcc15b`.
The review inspected static evidence only.
The ticket graph contains one complete outcome and passed independent ticket review.

A subsequent plan-only preflight found four omitted fixture directories and one omitted anchor registry.
The coordinator authorized their fence addition as an in-scope plan correction.
The correction adds no behavior promise and weakens no existing assertion.
A second closure pass adds the five command registries bound to the anchor package.

Independent ticket review accepted the corrected fence and complete graph at static confidence 9.
Its source checkpoint is bcf7e5d705ecd8d56bcd66d55645850a7245c949.
Clean-source preflight passed with fourteen green checks and two not-applicable checks.
The ticket write proposal found no missing paths or ordering requirements.
No implementation check or benchmark ran during this planning phase.

### Spec approval table

| subject | proposal | disposition |
|---|---|---|
| implementation line | Sol 6.1 at high effort | spec and ticket accepted |
| seam | the real degraded hook and selected core entry | spec and ticket accepted |
| coverage and edges | GR01-GR41 plus the named exclusions | spec and ticket accepted |
| ownership fences | fourteen ticket paths plus the review pickup | closure proposal and independent review passed |
| scope | malformed degraded refusal and parse-only Git diagnostics | spec and ticket accepted |
| residuals | grammar, scan, and `rg` outcomes remain separate | retained |
| ticket breakdown | 1. Refuse malformed guard commands with repair advice, blocked by none, complete GR1 outcome | independently accepted |
