# Pin removal

Status: staged

Decision source: the named reviewed artifact `specs/agent-push-guard/spec.md`, whose Out of scope prices the pin removal the reviewer closed on 2026-09-04

Verification log: 1 iteration(s) to accept — the round marked the hook ticket's five closure paths as a blocking overlap, and the author folded them as preflight-demanded closure under the landed no-blocker-edge convention, with seven prose findings folded

## Problem

The managed pre-push hook carries a drift clause. The clause compares the
pushed `.bench` tree with a local pin file that `bench gate pin` writes. The
reviewer owns the merge and reads every gate change in review, so the pin is
a second manual copy of that review. The pin verb needs a TTY and a typed
confirmation, and no automation can run it.

Without a pin file the hook prints a warning on every push. The local pin
file left this repository on 2026-09-04, so every push now warns. The roadmap
row FT141 proposes a command name that the pin porcelain holds, and the row
reads as blocked by that collision.

## Solution

The drift clause, the pin verb, the pin file, and the warning leave the kit.
The hook keeps one clause: it blocks a direct push to the protected branch,
and `bench.allowProtectedPush` lifts that clause for one repository. The gate
grammar is `bench gate [--fresh]`, and `bench help` lists no pin row.

The reference guide, the decision record, the FT141 detail, and the
distillation doc state the resulting state. The planted-reason proofs stay
the defence for a working-tree gate change, and the reviewer's merge review
is the human check.

## User stories

### Group A — the hook

Line: opus / medium. The rendered hook runs over a temp repository under one
existing test harness, and a wrong clause reds that harness.

1. As a reviewer, I want the pre-push hook to block a direct push to the protected branch, so that the merge stays mine.
2. As a reviewer, I want `bench.allowProtectedPush=true` to lift the branch clause for one repository, so that my escape hatch stays.
3. As a pusher, I want a topic push to pass with nothing on stderr, so that no `gate unpinned` warning appears.
4. As a pusher, I want a topic push to pass when a legacy pin file names a foreign tree, so that no drift clause runs.
5. As a pusher, I want a push of a commit with no `.bench` tree to pass, so that the hook reads no tree.
6. As a session, I want `bench guards` to report the hook's deny surface as `direct push to the protected branch`, so that the disclosure matches the enforcement.

### Group B — the command surface

Line: opus / medium. Each surface has a byte-exact fixture or a census test,
and a survivor reds it.

7. As an agent, I want `bench gate pin` to exit 2 with the usage `bench gate [--fresh]`, so that no verb writes a pin.
8. As an agent, I want `bench help` to list `bench gate [--fresh]` and no pin row, so that the inventory matches the executable.
9. As a maintainer, I want the `gate-pin` plumbing verb out of the registry, so that the core routes only live verbs.
10. As a maintainer, I want the wrapper to stop its `gate pin` route, so that no route reaches a dead key.
11. As a maintainer, I want the routing and parity tables to name no `gate-pin`, so that an exemption names only a live verb.
12. As a maintainer, I want the gate package to hold no pin reader or writer, so that no code path touches `bench-gate-pin`.
13. As a maintainer, I want the gate-owned record list in the test-report tests to drop the pin file, so that it names live records only.

### Group C — the guidance

Line: opus / medium. Guidance prose routes mid and high in the profile, and
the reviewer's 2026-08-26 rule holds every subagent at medium.

14. As a session, I want the hook bullet in the reference guide to name no pin verb, so that I never seek one.
15. As a maintainer, I want the anchored push-rule sentence in that bullet list to stay intact, so that the guard rule survives the edit.
16. As a teammate, I want ADR 0001 to record the planted-reason proofs as the gate-change defence, so that I read the decided state.
17. As a reviewer, I want the FT141 detail to say the proposed name is free, so that the row's blocker reads true.
18. As a reader, I want the distillation doc to name only live verbs, so that a stale reference never misleads.

## Implementation decisions

- The hook reads no pin path and no pin tree. It resolves the protected branch live from `origin/HEAD` with the baked token as the fallback, reads `bench.allowProtectedPush`, and loops the stdin ref lines once. The branch clause is the only exit 1.
- The hook header's `denies` field reads `direct push to the protected branch`. The `why` field names no repo-only path, because the package-core guard reds a claim word beside such a path.
- The gate package deletes the pin command, its review, its writer, and its path helper. The gate grammar is `bench gate [--fresh]`, and an argument `pin` is an unknown argument.
- The command registry drops the `gate-pin` entry and the pin help row. The routing and parity tables drop their `gate-pin` rows in the same commit. The routing check reds a registry row that the dispatch no longer names.
- The wrapper's `gate_command` keeps the bare run, `--fresh`, and the help spellings. Every other shape prints the usage and exits 2.
- The test-report tests drop the pin file from the gate-owned record list, because the list mirrors the live records.
- The reference guide's hook-layer bullet states the branch clause, the config lift, and the static deny surface of guard discovery. The anchored push-rule sentence keeps its bytes.
- ADR 0001 keeps its file name, because the guidance token sweep lists that path as an ordinary span. Its title and body record the planted-reason proofs alone.
- The FT141 detail says the name the row proposes is free and keeps the reviewer's decision open. The distillation doc names `internal/git` alone as the tree-hash owner.
- The hook keeps its own live branch resolution. One shared protected-branch source for the hook and the guard is a reviewer decision, priced in Out of scope.

## Testing decisions

- The hook tests run the rendered script the way git does, with one ref line on stdin, over a temp repository. The prior art is `TestPrePushHookAllowProtectedPushConfig`.
- The guards test renders the real shipped asset, so a header that keeps the drift clause reds the expected cell.
- The gate command test drives `Command` with a `pin` argument and expects the unknown-argument path.
- The help fixture, the disposition census, and the routing check observe the registry.
- The prose lane grades every edited Markdown file, and the review reads the guidance rows.

### Seam diagram

    git push (human or agent)
        │
        ▼
    stdin ref lines ──▶ [ installed pre-push hook, rendered from internal/adopt/prepush.sh ] ──▶ exit 0 | exit 1 + blocked:
                             ◀ tests attach here: runPrePushHook over a temp repository

    bench gate <args>
        │
        ▼
    args ──▶ [ gate.Command ] ──▶ RunCommand | usage on stderr + exit 2
                             ◀ tests attach here: TestCommandHandlesPublicGateUsageWithoutStartingTheOracle

    bench help ──▶ [ commandRegistry ] ──▶ inventory text
                             ◀ tests attach here: TestHelpInventoryIsComplete and TestCommandDispositionsAreComplete

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| PR1 | 1 | the rendered hook exits 1 on `refs/heads/main` and prints `blocked: direct push to main` | `internal/adopt/link_hook_test.go` (`TestPrePushHookAllowProtectedPushConfig`) | a hook cut to an empty script passes the push |
| PR2 | 2 | with `bench.allowProtectedPush` set to `true` the hook exits 0 on `refs/heads/main`, and set to `false` it exits 1 | `internal/adopt/link_hook_test.go` (`TestPrePushHookAllowProtectedPushConfig`) | a hook that drops the config read blocks both |
| PR3 | 3 | the hook exits 0 on `refs/heads/topic` with an empty stderr | a new test `TestPrePushHookTopicPushIsSilent` in the adopt package's hook tests | the old hook prints `bench: gate unpinned` on stderr |
| PR4 | 4, 5 | with a `bench-gate-pin` file in the git dir that names a foreign tree, the hook exits 0 on `refs/heads/topic` for a local oid that is not a commit, with an empty stderr | a new test `TestPrePushHookIgnoresALegacyPin` in the adopt package's hook tests | the old hook exits 1 with `pushed commit has no .bench tree` |
| PR5 | 6 | `bench guards` renders the pre-push denies cell as `direct push to the protected branch` | `internal/guards/guards_test.go` (`TestCommandRendersRealStaleManagedPrePushHookAndRepairAction`) | the fixture is the shipped asset, so a header that keeps the drift clause reds the expected cell |
| PR6 | 7 | `bench gate pin` exits 2 and prints `usage: bench gate [--fresh]` on stderr | `internal/gate/command_test.go` (`TestCommandHandlesPublicGateUsageWithoutStartingTheOracle`) | a surviving pin case prints `usage: bench gate pin` |
| PR7 | 7 | `bench gate --brief` prints `usage: bench gate [--fresh]` on stderr | `cmd/bench/main_test.go` (`TestRunGateRejectsBriefUsage`) | an unchanged usage string keeps `pin` |
| PR8 | 8 | `bench help` prints `bench gate [--fresh]` and no line that holds `bench gate pin` | `cmd/bench/main_test.go` (`TestHelpInventoryIsComplete`) | a surviving help row reds the byte-exact fixture |
| PR9 | 9 | the system-attachment census holds no `gate-pin` | `cmd/bench/command_registry_test.go` (`TestCommandDispositionsAreComplete`) | a surviving registry entry reds the census |
| PR10 | 10 | `bench gate pin` through the wrapper exits 2 and prints the usage line | review-owned: the reviewer reads `gate_command` in the wrapper | the route check names no pin route, so a dead `gate-pin` route survives the gate |
| PR11 | 11 | the subcommand-routing registry names no `gate-pin` | `internal/conformance/subcommand_routing_test.go` (`checkSubcommandRouting`) | the check reds a registry row that the dispatch no longer names |
| PR12 | 11 | the entry-point-parity exemption table names no `gate-pin` | review-owned: the parity check reads the registry only | a stale exemption survives the gate |
| PR13 | 12 | a search for `bench-gate-pin` and `PinCommand` over the Go sources returns no line | review-owned: the reviewer runs the search | the compiler reds a caller of a deleted function, and the search reds a survivor |
| PR14 | 13 | the gate-owned record list in the two test-report tests holds no `bench-gate-pin` row | review-owned: the reviewer reads the two tests | a survivor keeps a record the gate never writes |
| PR15 | 14 | the reference guide's hook-layer bullet holds no `gate pin`, `drift`, or `pinned` | review-owned: the reviewer reads the bullet | no anchor forbids the bytes, so a survivor passes the gate |
| PR16 | 15 | the canary fixture `reference-agent-push-rule` reports its diagnostic only when its mutation runs | `internal/anchors/registry_data_test.go` (`TestAgentPushRuleAnchorRedOnRemoval`) over the fixture `tests/canary/workflow-guidance-anchors/reference-agent-push-rule` | an edit that drops the anchored sentence reds the fixture bite |
| PR17 | 16 | ADR 0001 holds no `pin`, and its title ends with `planted-reason proofs` | review-owned: the reviewer reads the ADR | no check reads the ADR body |
| PR18 | 17 | the first paragraph of the FT141 detail holds the word `free`, and the file keeps `Next: decide` | `roadmap-detail-integrity` for the shape, and review for the paragraph | a row edit that breaks the detail shape reds the check |
| PR19 | 18 | the distillation doc names `internal/git` alone as the tree-hash owner | review-owned: the reviewer reads the line | no check reads the doc |

### Edge inventory

- Error paths: the hook's `git symbolic-ref` fails with no remote, and the baked token stays the protected branch (existing behavior, `TestInspectPrePush*`). The hook's config read fails with no config, and the clause stays armed (PR1).
- Empty input: an empty stdin yields no ref line, and the hook exits 0. A last line with no trailing newline still loops, because the read keeps its `|| [ -n "$line" ]` guard.
- Boundary values: an all-zero local oid is a deletion push, and the branch clause still reads its remote ref. A two-line stdin with one protected ref exits 1.
- Legacy state: a `bench-gate-pin` file in a linked repository's git dir stays inert (PR4).
- Currency: a linked repository's installed hook reads as `stale` after the kit upgrade, and `bench status` names `bench link` as the repair. The existing test `TestCommandRendersRealStaleManagedPrePushHookAndRepairAction` proves the derivation.
- Re-run idempotency: the hook writes nothing, and the gate command writes nothing on the usage path.
- Partial implementation: a build that removes the verb and keeps the drift clause reds PR3 and PR4. A build that removes the clause and keeps the verb reds PR6 to PR9.
- Audience: the hook and the verb serve every repository that links the kit. The guidance serves every session in a linked repository.
- Package-variable swaps: no test swaps a package variable.
- Absent versus empty: an absent pin file and an empty pin file both leave the hook silent, because the hook reads neither (PR3, PR4).
- Hostile paths: a git dir path with a space reaches no hook line, because the hook resolves no path.

**Won't handle** — a deletion of a legacy pin file from a linked repository's git dir — the hook reads nothing there, and `bench link` stays the caller.

**Won't handle** — a force or delete clause in the pre-push hook — the hook guards a human pusher, and `git push origin topic` stays its surviving caller.

**Won't handle** — a rename of the ADR file — the guidance token sweep lists `docs/adr/0001-working-tree-gate-tripwire.md` as an ordinary span, and the file keeps that name.

**Won't handle** — the anchored sentence `The destructive-git guard allows an agent push to any branch other than the default branch.` — the edit keeps those bytes, and PR16 proves it.

**Won't handle** — the FT141 baseline record itself — that row keeps its own decision, and the detail edit changes one paragraph.

## Ownership fences

- `specs/pin-removal/`
- `reviews/pin-removal.md`
- `internal/adopt/prepush.sh`
- `internal/adopt/link_hook_test.go`
- `internal/guards/guards_test.go`
- `internal/gate/phases.go`
- `internal/gate/gate.go`
- `internal/gate/command_test.go`
- `cmd/bench/main.go`
- `cmd/bench/main_test.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_test.go`
- `internal/conformance/entry_point_parity_test.go`
- `bin/bench.sh`
- `internal/testreport/check_test.go`
- `internal/testreport/testreport_test.go`
- `tests/canary/package-core-guard/unrouted-subcommand`
- `tests/canary/injected-ports/unregistered-port` — closure headroom only
- `tests/canary/docs-currency-token-diet/missing-cli-inventory` — closure headroom only
- `tests/canary/docs-currency-token-diet/stale-cli-doc-reference` — closure headroom only
- `tests/canary/docs-currency-token-diet/stale-skill-cli-reference` — closure headroom only
- `tests/canary/load-validity-metadata/extensionless-gate-ref` — closure headroom only
- `tests/canary/package-core-guard/bounds-duplicate-owner` — closure headroom only
- `tests/canary/package-core-guard/reintroduced-bare-skip` — closure headroom only
- `tests/canary/workflow-guidance-anchors/distillation-reduced-schema-refactor-shape` — closure headroom only
- `.bench/BENCH-reference.md`
- `docs/adr/0001-working-tree-gate-tripwire.md`
- `roadmap/FT141.md`
- `docs/reporesident-distillation.md`
- `tests/canary/workflow-guidance-anchors/` — the fixtures that pin the reference guide, closure headroom only
- `tests/canary/docs-currency-token-diet/` — closure headroom only
- `tests/canary/skills-index-command-adapters/` — closure headroom only

The fence is the union of the three tickets' `Writes:` lines, closed by
`bench preflight build` over the fixture and registry pins. A closure
headroom file creates no blocker edge. The hook ticket carries five
registry paths as closure only, and the verb ticket owns their edits.

## Out of scope

- One shared protected-branch source for the hook and the guard. The hook must run with no `bench` on `PATH`, so the recommendation is to keep the two resolvers. 2 edits, 1 gate run.
- A `bench worktree land` step that pushes the landed branch. 3 edits, 1 gate run.
- The FT141 baseline record. That row's own spec prices it.

## Further notes

Flagged additions beyond the decision source:

- The FT141 detail edit. The 2026-09-04 grill named the freed name, and the reviewer can veto the paragraph.
- The distillation doc edit. The doc names the pin verb as a live owner, and the source names "the docs" without a list.
- The test-report record-list edit. The list mirrors live records, and the source names no test.

Build decisions recorded for reviewer veto:

- The hook keeps its own live branch resolution, and the shared source stays out of scope.
- No new anchor forbids the pin bytes in the reference guide. The removal is review-owned, and a Forbid row would add one registry row and one fixture.

Source-sentence-to-row table:

| source sentence | rows |
|---|---|
| the drift clause leaves the kit | PR3, PR4, PR5 |
| `bench gate pin` leaves the kit | PR6, PR7, PR8, PR9, PR10, PR11, PR12, PR13 |
| the pin file leaves the kit | PR4, PR13, PR14 |
| the `gate unpinned` warning leaves the kit | PR3 |
| the ADR and the docs rewritten | PR15, PR16, PR17, PR18, PR19 |
| the branch clause and `bench.allowProtectedPush` stay | PR1, PR2 |

Pre-review proof checklist:

- Cited symbols: each of these resolves in the tree at the spec commit. `PinCommand`, `pinCommand`, `showPinReview`, `dirtyBench`, `existingPinnedCommit`, `writePinFromHead`, `pinPath`, `pinFileName`, `commandUsage`, `gate.Command`, `gate_command`, `gate_usage`. `TestPrePushHookAllowProtectedPushConfig`, `runPrePushHook`, `writeHook`, `hookTestRepo`, `TestCommandRendersRealStaleManagedPrePushHookAndRepairAction`, `TestCommandHandlesPublicGateUsageWithoutStartingTheOracle`, `TestRunGateRejectsBriefUsage`, `TestHelpInventoryIsComplete`, `TestCommandDispositionsAreComplete`, `checkSubcommandRouting`, `parityExemptCommands`. The two new hook tests are new.
- Import edges: none new. The build removes the `terminal` import from the gate package's phases file, and the adopt package still imports `internal/terminal`.
- Source-row clauses and occurrences: the source is one Out of scope row in `specs/agent-push-guard/spec.md`, and the table above lists each clause once.
- Promised field labels: the hook header field `denies` reads `direct push to the protected branch`.
- Changed-function callers: `PinCommand` has two callers, the registry entry and `gate.Command`, and both leave. `gate.Command` keeps its one registry caller. `pinCommand` and its four helpers have no caller outside the pin path.
- Copy survival: PR13 reds a surviving pin reader or writer in the Go sources.

Reader sweep of the pin. The readers are:

- the hook asset, the hook test, and the guards test, which take PR1 to PR5
- the gate package, the registry, the help fixture, the disposition census, the wrapper, the routing table, and the parity table, which take PR6 to PR13
- the two test-report tests, which take PR14
- the reference guide, ADR 0001, the FT141 detail, and the distillation doc, which take PR15 to PR19
- `internal/guards`, `internal/status`, and the adopt hook inspection, which read the hook as data and take no edit
- the ft173 command-help inventory asset, a dated research snapshot pinned to its subject commit, which is excluded
- `capture/FIXES.md` and `docs/audits/`, historical records, which are excluded
- the `pinContent` fixture in the lines tests, a different pin, which is excluded

The stale-command sweep grades `/bench-*` and `$bench-*` names, and the cold-pickup CLI sweep reads only the first token after `bench`. So a doc that names `bench gate pin` survives the gate, and the sweep above is the complete list of such docs.

The reference guide and ADR 0001 are kit-guidance files, and the guidance ticket carries the shipped-surface claim words the package-core guard reads.

The trust chain is unchanged. Git runs the installed hook, and `bench link` authenticates the managed hook by its marker before it replaces the hook.

Every subagent runs `opus` at low or medium effort. The review round runs `opus` at low effort with a cap of two iterations.
