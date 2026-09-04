# Agent push guard

Status: staged

Decision source: the reviewer-confirmed conversation of 2026-09-04, a one-round grill that closed five questions

Verification log: 2 iteration(s) to accept — the first round blocked on the detached-`HEAD` refspec, which the literal `HEAD` from the git owner let through

## Problem

The Bash guard `block-dangerous-git` denies every `git push` an agent types.
The reviewer owns the merge, and the guard protects that ownership by a blanket
rule. The blanket rule also stops the agent from publishing a topic branch for
review. The reviewer then pushes every branch by hand.

The git pre-push hook already knows the protected branch. The guard does not
read it. So the guard cannot tell a topic branch from the default branch.

## Solution

The guard allows a plain push whose every destination is not the default
branch. It keeps the denial for a forced push, a deletion, and the three
broadcast forms. A bare `git push` resolves its destination the way git does,
from `push.default` and the checked-out branch. When the guard cannot resolve
the destination, or the default branch, it denies the push and names the fix.

The reviewer's config `bench.allowProtectedPush` never lifts the guard. The
protected branch comes from the one Go owner of the default-branch fact. The
reference guide states the rule for every session that reads it.

## User stories

### Group A — the destination rule

Line: opus / medium. The seam is one classifier function with a table test,
and a wrong verdict reds that table.

1. As an agent, I want `git push origin topic` allowed when `topic` is not the default branch, so that I can publish a phase branch.
2. As a reviewer, I want `git push origin main` denied when `main` is the default branch, so that the merge stays mine.
3. As a reviewer, I want `git push origin HEAD:main` denied, so that a source-colon-destination spelling cannot reach the default branch.
4. As an agent, I want `git push origin main:topic` allowed, so that the source name never decides the verdict.
5. As a reviewer, I want `git push origin refs/heads/main` denied, so that the full ref spelling cannot reach the default branch.
6. As an agent, I want `git push origin refs/tags/main` allowed, so that a tag that shares the branch name is not read as the branch.
7. As an agent, I want `git push origin HEAD` allowed when the checked-out branch is a topic branch, so that the common spelling works.
8. As a reviewer, I want `git push origin HEAD` denied when the checked-out branch is the default branch, so that `HEAD` cannot alias it.
9. As a reviewer, I want a two-refspec push denied when one refspec targets the default branch, so that a mixed list cannot carry it.
10. As a reviewer, I want the default branch read from the repository's live default, so that a repository whose default is `master` is protected too.
11. As a reviewer, I want every push denied when the default branch cannot be resolved, so that the guard fails closed.
12. As a reviewer, I want `git push --dry-run origin main` denied, so that no option exempts the destination rule.

### Group B — the denied forms

Line: opus / medium. Each form is one lexical predicate under the same table
test.

13. As a reviewer, I want `-f` and `--force` denied on every destination, so that the agent never rewrites a remote branch.
14. As a reviewer, I want `--force-with-lease` denied, with and without a value, so that the lease spelling is not a loophole.
15. As a reviewer, I want a refspec that starts with `+` denied, so that the refspec spelling of force is not a loophole.
16. As a reviewer, I want `--delete`, `-d`, and an empty-source refspec denied, so that the agent never deletes a remote ref.
17. As a reviewer, I want `--all`, `--mirror`, and `--tags` denied, so that the agent never pushes refs it did not name.
18. As a reviewer, I want `bench.allowProtectedPush=true` to leave the guard in force, so that my escape hatch stays mine.

### Group C — the bare push

Line: opus / medium. The destination logic reads git config and the checked-out
branch, and a temp repository proves each branch of it.

19. As an agent, I want a bare `git push` on a topic branch allowed under `simple` or `current`, so that a fresh branch publishes.
20. As a reviewer, I want a bare `git push` denied when the checked-out branch is the default branch, so that the bare form cannot reach it.
21. As a reviewer, I want a bare `git push` under `upstream` denied when the upstream is the default branch, so that tracking cannot reach it.
22. As a reviewer, I want a bare `git push` denied on a detached `HEAD`, so that an unresolvable destination fails closed.
23. As a reviewer, I want a bare `git push` denied under `push.default` `matching` or `nothing`, so that a broadcast or an empty destination fails closed.
24. As an agent, I want the unresolved refusal to name the fix, so that I type the remote and the branch next.
25. As an agent, I want `git push origin` with no refspec to obey the bare rule, so that a named remote changes nothing.
26. As a reviewer, I want a push outside a git repository denied, so that the guard never guesses a destination.
27. As a reviewer, I want a bare push denied when the git probe hangs, so that a slow probe fails closed.
28. As a reviewer, I want `git push origin HEAD` on a detached `HEAD` denied, so that the literal `HEAD` never passes as a branch name.

### Group D — the refusal and the scan

Line: opus / medium. The message and the scan are existing seams.

29. As an agent, I want the refusal to name the deny class, so that I know which rule I hit.
30. As a reviewer, I want `env git push origin main`, `bash -c 'git push origin main'`, and `git status && git push origin main` denied, so that a prefix or a wrapper is not a loophole.
31. As an agent, I want `git push origin topic && ls` allowed, so that an allowed push composes with an ordinary follow-on.
32. As a reviewer, I want the guard subcommand to exit 2 on `git push --force` and print `BLOCKED:`, so that the shim still reads a block.

### Group E — the guidance

Line: opus / medium. A guidance sentence steers every session that reads it, and the scorecard routes guidance prose to medium.

33. As a session, I want the reference guide's hook-layer list to state the push rule, so that I push a topic branch without a hand-back.
34. As a maintainer, I want an anchor to pin that sentence, so that a later edit cannot drop the rule in silence.

## Implementation decisions

- The `push` case in the classifier becomes a verdict function. The function skips each option that takes a separate value. It reads the first free arg as the remote and the rest as refspecs.
- A refspec that starts with `+` is force. A refspec with an empty source is a deletion. A refspec `src:dst` targets `dst`. A refspec `src` targets `src`, except `HEAD`, which targets the checked-out branch.
- A destination is protected when it equals the default branch after removal of a `refs/heads/` prefix. A destination under another `refs/` namespace is never protected.
- With no refspec, the destination comes from one new function in the git package. Under `simple` or `current`, or an unset `push.default`, it is the checked-out branch. Under `upstream` or `tracking`, it is the upstream branch name. Under `matching` or `nothing`, on a detached `HEAD`, outside a repository, or on a probe error, the function reports no destination.
- The default branch comes from the existing Go owner `git.ResolvedDefault`. No destination and no default branch both deny with the unresolved class.
- The checked-out fact reports no branch on a detached `HEAD`. The git owner returns the literal `HEAD` there, and the adapter maps that literal and any error to no branch. A `HEAD` refspec with no branch denies with the unresolved class.
- The `Checker` gains three injected facts: the default branch, the checked-out branch, and the bare destination. The guard subcommand wires them from the git package. The fakes in the classifier tests set them by hand.
- The deny table replaces its one push row with seven classes: default branch, force, delete, all, mirror, tags, and unresolved. The unresolved class carries one advice sentence, and the block message appends it. The table stays the one source of labels and advice.
- Every value-taking option is skipped with its value, so `-o x` and `--repo x` never read `x` as the remote. No option exempts a push from the destination rule.
- The reference guide's hook-layer list gains one sentence that states the rule. One anchor row pins the sentence, with its canary fixture.
- The pre-push hook and its pin logic stay as they are. The second spec removes the pin.

## Testing decisions

- The classifier table test drives every lexical case through the full tokenize, scan, and classify path with a fake `Checker`.
- A composition test over a temp repository drives the real facts: the checked-out branch, `push.default`, the upstream, a detached `HEAD`, and the `bench.allowProtectedPush` config.
- The git package test drives the destination function over the same temp repository states.
- The stub-git timeout test extends to the bare push.
- The guard subcommand test keeps its exit and message assertions on a lexical case.
- The anchor registry test and one canary fixture pin the guidance sentence.

### Seam diagram

    PreToolUse envelope (Claude, Codex)
        │
        ▼
    block-dangerous-git.sh ──▶ bench guard-git ──▶ [ gitguard.Classify + pushVerdict ] ──▶ exit 0 | exit 2 + BLOCKED
                                                        ◀ tests attach here: TestClassifyVerdicts with a fake Checker
        │
        ▼
    Checker facts ──▶ [ git.ResolvedDefault, git.CheckedOutBranch, git.BarePushDestination ] ──▶ branch names
                                                        ◀ tests attach here: temp-repository composition tests

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| PG1 | 1 | `git push origin topic` with default `main` returns the allow verdict | `internal/gitguard/verdict_test.go` (`TestClassifyVerdicts`) with a fake default of `main` | the old blanket denial reds the row |
| PG2 | 2 | `git push origin main` with default `main` returns the label `git push to the default branch` | `internal/gitguard/verdict_test.go` (`TestClassifyVerdicts`) | an allow-all rewrite reds the row |
| PG3 | 3 | `git push origin HEAD:main` returns the default-branch label | `internal/gitguard/verdict_test.go` (`TestClassifyVerdicts`) | a rule that reads the source name allows it |
| PG4 | 4 | `git push origin main:topic` returns the allow verdict | `internal/gitguard/verdict_test.go` (`TestClassifyVerdicts`) | a rule that reads the source name denies it |
| PG5 | 5 | `git push origin refs/heads/main` returns the default-branch label | `internal/gitguard/verdict_test.go` (`TestClassifyVerdicts`) | a bare string compare allows the full ref |
| PG6 | 6 | `git push origin refs/tags/main` returns the allow verdict | `internal/gitguard/verdict_test.go` (`TestClassifyVerdicts`) | a suffix compare denies the tag |
| PG7 | 7 | `git push origin HEAD` with a checked-out `topic` returns the allow verdict | `internal/gitguard/verdict_test.go` (`TestClassifyVerdicts`) with a fake checked-out branch | a rule that denies `HEAD` reds the row |
| PG8 | 8 | `git push origin HEAD` with a checked-out `main` returns the default-branch label | `internal/gitguard/verdict_test.go` (`TestClassifyVerdicts`) | a rule that reads `HEAD` as a literal allows it |
| PG9 | 9 | `git push origin topic main` returns the default-branch label | `internal/gitguard/verdict_test.go` (`TestClassifyVerdicts`) | a rule that reads only the first refspec allows it |
| PG10 | 10 | `git push origin master` with a fake default of `master` returns the default-branch label | `internal/gitguard/verdict_test.go` (`TestClassifyVerdicts`) | a literal `main` allows it |
| PG11 | 10 | `git push origin main` with a fake default of `master` returns the allow verdict | `internal/gitguard/verdict_test.go` (`TestClassifyVerdicts`) | a literal `main` denies it |
| PG12 | 11 | `git push origin topic` with no resolvable default returns the label `git push with an unresolved destination` | `internal/gitguard/verdict_test.go` (`TestClassifyVerdicts`) | a fail-open default allows it |
| PG13 | 12 | `git push --dry-run origin main` returns the default-branch label | `internal/gitguard/verdict_test.go` (`TestClassifyVerdicts`) | a dry-run exemption allows it |
| PG14 | 13 | `git push -f origin topic` and `git push --force origin topic` each return the label `git push --force` | `internal/gitguard/verdict_test.go` (`TestClassifyVerdicts`) | a destination-only rule allows both |
| PG15 | 13 | `git push -fu origin topic` returns the force label | `internal/gitguard/verdict_test.go` (`TestClassifyVerdicts`) | a whole-token compare misses the cluster |
| PG16 | 14 | `git push --force-with-lease origin topic` and `git push --force-with-lease=topic origin topic` each return the force label | `internal/gitguard/verdict_test.go` (`TestClassifyVerdicts`) | an exact-token compare misses the value form |
| PG17 | 15 | `git push origin +topic` returns the force label | `internal/gitguard/verdict_test.go` (`TestClassifyVerdicts`) | a rule that reads `+topic` as a name allows it |
| PG18 | 16 | `git push --delete origin topic`, `git push -d origin topic`, and `git push origin :topic` each return the label `git push --delete` | `internal/gitguard/verdict_test.go` (`TestClassifyVerdicts`) | a rule without the empty-source case allows the third |
| PG19 | 17 | `git push --all origin`, `git push --mirror origin`, and `git push --tags origin` each return their own label | `internal/gitguard/verdict_test.go` (`TestClassifyVerdicts`) | a rule with no refspec treats each as a bare push |
| PG20 | 18 | with `bench.allowProtectedPush` set to `true` in a temp repository, `git push origin main` returns the default-branch label | `internal/gitguard/checker_junction_test.go` (`TestClassifyRealCheckerResolvedComposition`) | a guard that reads the config allows it |
| PG21 | 19 | a bare `git push` on a checked-out `topic` with `push.default` unset returns the allow verdict | `internal/gitguard/checker_junction_test.go` (`TestClassifyRealCheckerResolvedComposition`) | a rule that needs `@{push}` denies the fresh branch |
| PG22 | 19 | the same bare push with `push.default` set to `current` returns the allow verdict | `internal/gitguard/checker_junction_test.go` (`TestClassifyRealCheckerResolvedComposition`) | a rule that handles only the unset case denies it |
| PG23 | 20 | a bare `git push` on a checked-out `main` returns the default-branch label | `internal/gitguard/checker_junction_test.go` (`TestClassifyRealCheckerResolvedComposition`) | a rule that allows every bare push reds the row |
| PG24 | 21 | a bare `git push` on `topic` with upstream `origin/main` and `push.default` set to `upstream` returns the default-branch label | `internal/gitguard/checker_junction_test.go` (`TestClassifyRealCheckerResolvedComposition`) | a rule that reads only the checked-out name allows it |
| PG25 | 22 | a bare `git push` on a detached `HEAD` returns the unresolved label | `internal/gitguard/checker_junction_test.go` (`TestClassifyRealCheckerResolvedComposition`) | a rule that reads `HEAD` as a name allows it |
| PG26 | 23 | a bare `git push` with `push.default` set to `matching` returns the unresolved label | `internal/gitguard/checker_junction_test.go` (`TestClassifyRealCheckerResolvedComposition`) | a rule that reads the checked-out name allows the broadcast |
| PG27 | 23 | a bare `git push` with `push.default` set to `nothing` returns the unresolved label | `internal/gitguard/checker_junction_test.go` (`TestClassifyRealCheckerResolvedComposition`) | a rule that reads the checked-out name allows the empty push |
| PG28 | 24 | the block message for the unresolved label ends with the sentence `Name the remote and the branch: git push <remote> <branch>.` | `internal/gitguard/gitguard_test.go` (`TestBlockMessageNamesLabel`) extended | a generic message leaves the agent without the fix |
| PG29 | 25 | `git push origin` with no refspec on a checked-out `main` returns the default-branch label | `internal/gitguard/checker_junction_test.go` (`TestClassifyRealCheckerResolvedComposition`) | a rule that reads `origin` as a refspec allows it |
| PG30 | 26 | a bare `git push` outside a repository returns the unresolved label | `internal/gitguard/checker_junction_test.go` (`TestClassifyRealCheckerResolvedComposition`) with a non-repository directory | a fail-open probe allows it |
| PG31 | 27 | a bare `git push` with a stub `git` that sleeps past the probe bound returns the unresolved label | `internal/gitguard/checker_junction_test.go` (`TestClassifyRealCheckerTimeoutComposition`) | a probe without a bound hangs the guard |
| PG32 | 29 | the block message for each of the seven push labels names that label inside backticks | `internal/gitguard/gitguard_test.go` (`TestBlockMessageNamesLabel`) extended | a shared label hides which rule fired |
| PG33 | 30 | `env git push origin main`, `bash -c 'git push origin main'`, and `git status && git push origin main` each return the default-branch label | `internal/gitguard/verdict_test.go` (`TestClassifyVerdicts`) | a verdict reached only in bare command position misses the three forms |
| PG34 | 31 | `git push origin topic && ls` returns the allow verdict | `internal/gitguard/verdict_test.go` (`TestClassifyVerdicts`) | a rule that denies any push in a chain reds the row |
| PG35 | 32 | the guard subcommand exits 2 on `git push --force` and prints `BLOCKED:` | `cmd/bench/main_test.go` (`TestGuardGitBlockAllow`) | a subcommand that maps the new labels to exit 0 reds the row |
| PG36 | 33, 34 | the anchor registry holds one row whose needle is the push-rule sentence in the reference guide's hook-layer list | a new red-on-removal test in the anchor registry package | a sentence without an anchor can vanish in silence |
| PG37 | 19 | the destination function returns `topic` for a checked-out `topic` with `push.default` unset | a new test in the git package | a function that requires an upstream returns no destination |
| PG38 | 21, 22, 23 | the destination function returns the upstream branch name under `upstream`, and no destination on a detached `HEAD`, under `matching`, and under `nothing` | the new git package test | a function that ignores `push.default` returns the checked-out name for all four |
| PG39 | 28 | `git push origin HEAD` with a fake checked-out fact that reports no branch returns the unresolved label | `internal/gitguard/verdict_test.go` (`TestClassifyVerdicts`) | a rule that reads `HEAD` as a name allows it |
| PG40 | 28 | `git push origin HEAD` on a detached `HEAD` in a temp repository returns the unresolved label | `internal/gitguard/checker_junction_test.go` (`TestClassifyRealCheckerResolvedComposition`) | an adapter that passes the literal `HEAD` through allows it |
| PG41 | 34 | the canary fixture reports the anchor diagnostic when its mutation replaces the sentence | the new canary fixture under the workflow-guidance-anchors family | a registry row without a fixture never proves it bites |

### Edge inventory

- Error paths: a git probe error, a timeout, a detached `HEAD`, and a missing repository each yield no destination. The guard then denies with the unresolved label (PG25, PG30, PG31, PG40).
- Empty input: a bare `git push` reads no refspec and takes the bare rule (PG21 to PG27).
- Boundary values: a refspec of `+` alone is force; a refspec of `:` alone is a deletion; both deny.
- Value-taking options: `-o`, `--push-option`, `--repo`, `--receive-pack`, and `--exec` skip their value, so `git push -o ci.skip origin topic` allows.
- Option terminator: tokens after `--` are free args.
- Re-run idempotency: the verdict is pure over its inputs, and no probe writes.
- Hostile paths: a remote spelled as a URL is the first free arg and never a refspec.
- Partial implementation: a build that ships the destination rule without the bare rule reds PG21 to PG27.
- Audience: the classifier serves every repository that links the kit, and the guidance sentence serves every session in a linked repository.
- Package-variable swaps: no test swaps a package variable; the fakes go through the injected `Checker`.

**Won't handle** — `--follow-tags` — it adds only annotated tags reachable from the pushed commits, and `git push --follow-tags origin topic` stays an allowed caller.

**Won't handle** — a remote named like a branch, as in `git push main` — the first free arg is always the remote, and `git push origin main` stays the denied caller.

**Won't handle** — a force or delete clause in the pre-push hook — the hook guards a human pusher, and `git push origin topic` stays its surviving caller.

**Won't handle** — the degraded shim rim — a missing core still refuses every git call, and a session with the core stays the surviving caller.

**Won't handle** — a `remote.<name>.push` refspec config — the destination function reports no destination, and the agent names the branch instead.

**Won't handle** — the route row `git push` and the `/bench` push offer — the row fires for the default branch, and that push stays the reviewer's.

## Ownership fences

- `specs/agent-push-guard/`
- `reviews/agent-push-guard.md`
- `internal/gitguard/`
- `internal/git/push_destination.go`
- `internal/git/push_destination_test.go`
- `cmd/bench/main.go`
- `cmd/bench/main_test.go`
- `.bench/BENCH-reference.md`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `tests/canary/workflow-guidance-anchors/reference-agent-push-rule/`
- `tests/canary/workflow-guidance-anchors/` — the fixtures that pin the reference guide, closure headroom only
- `tests/canary/docs-currency-token-diet/` — closure headroom only
- `tests/canary/skills-index-command-adapters/` — closure headroom only
- `tests/canary/canonical-path-owner/second-derivation`
- `tests/canary/injected-ports/unregistered-port`
- `tests/canary/package-core-guard/unrouted-subcommand`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_test.go`

The fence is the union of the four tickets' `Writes:` lines, closed by
`bench preflight build` over the fixture and registry pins. A closure
headroom file creates no blocker edge; only a file a ticket's `What to build`
names does.

## Out of scope

- The pin removal: the drift clause, `bench gate pin`, the pin file, and the `gate unpinned` warning leave the kit, with the ADR and the docs rewritten. 11 edits, 2 gate runs.
- A `bench worktree land` step that pushes the landed branch. 3 edits, 1 gate run.
- A `--follow-tags` deny class. 2 edits, 1 gate run.
- One shared protected-branch source for the hook and the guard. The pin removal rewrites the hook, and the reviewer decides the source there. 2 edits, 1 gate run.

## Further notes

Flagged additions beyond the decision source:

- The bare rule reads `push.default`. The grill said "the way git does". Git's own `@{push}` fails for a topic branch with no remote-tracking ref, so the guard reads the same config git reads.
- `matching`, `nothing`, and a `remote.<name>.push` refspec deny as unresolved. The grill named only the unresolvable case, and these three have no single destination.
- A destination is protected by branch name on any remote. The grill said "the protected branch" and named no remote.
- `--dry-run` receives no exemption. The grill named no option exemption.
- A destination under `refs/tags/` or another `refs/` namespace is never protected. The grill spoke of branches.
- The guard fails closed with no repository and with no resolvable default branch. The grill named the unresolvable destination only.
- The guidance sentence and its anchor. The grill said every session must know the rule, and the reference guide is where the hook layers are stated.

Build decisions recorded for reviewer veto:

- Row PG31 names a probe bound that the destination probes do not carry. Only the ref probe behind `RefResolves` and `BranchExists` has a timeout. `Output` and `OK` run git with no bound. The row passes because the sleeping stub exits with empty output, and the destination fact then reports no branch. A git that never exits blocks the guard on the push path. A bound on `Output` and `OK` changes every caller in the git package, so it is a reviewer decision and is not in this build.
- One exported helper, `git.CheckedOutName`, owns the mapping from the checked-out probe to a branch name. The wire ticket wrote that mapping three times, and the build collapsed the copies before review.

Source-sentence-to-row table:

| source sentence | rows |
|---|---|
| allow a plain push when every destination ref is not the protected branch | PG1, PG3, PG4, PG5, PG6, PG7, PG8, PG9, PG34 |
| deny a push to the protected branch | PG2, PG10, PG11, PG13, PG33 |
| deny force in every spelling | PG14, PG15, PG16, PG17 |
| deny a deletion | PG18 |
| deny `--all`, `--mirror`, `--tags` | PG19 |
| resolve a bare push the way git does, allow a non-protected target | PG21, PG22, PG37 |
| deny a bare push that targets the protected branch | PG23, PG24, PG29 |
| deny an unresolvable destination and name the fix | PG12, PG25, PG26, PG27, PG28, PG30, PG31, PG38, PG39, PG40 |
| `bench.allowProtectedPush` never lifts the guard | PG20 |
| one Go owner of the protected branch | PG10, PG11, PG12 |
| the refusal names the class and the shim reads a block | PG32, PG35 |
| the reference guide states the rule | PG36, PG41 |

Pre-review proof checklist:

- Cited symbols: each of these resolves in the tree at the spec commit. `Classify`, `classify`, `Checker`, `denyTable`, `BlockMessage`, `guardGit`, `git.ResolvedDefault`, `git.CheckedOutBranch`. `TestClassifyVerdicts`, `TestClassifyRealCheckerResolvedComposition`, `TestClassifyRealCheckerTimeoutComposition`, `TestBlockMessageNamesLabel`, `TestGuardGitBlockAllow`, `TestPrePushHookAllowProtectedPushConfig`. `git.BarePushDestination` is new.
- Import edges: `cmd/bench` imports `internal/gitguard` and `internal/git` today. The classifier's junction test imports `internal/git` today. No new edge.
- Source-row clauses and occurrences: the decision source is the 2026-09-04 conversation, and the source-sentence table above lists each clause once.
- Promised field labels: none.
- Changed-function callers: `classify` has one caller, the scan. `Classify` has one production caller, `guardGit`, and the test files in the package. `BlockMessage` has one caller, `guardGit`. The `Checker` literal has one production site, `guardGit`, and three test sites in the package.
- Copy survival: none.

Reader sweep of the push deny label. The readers are:

- the classifier table test and the guard subcommand test, which take rows above
- the system rim test, which reads the degraded rim this spec leaves as it is
- the reference guide, which takes PG36 and PG41
- the README hook table and the field guide, which name destructive authority in general and stay true
- the route command, whose push offer stays true for the default branch
- the audit inputs under `docs/audits/`, which are historical records and take no edit

The session memory that says the agent never pushes is outside the tree. The reviewer's session updates it after the landing.

The reference guide and the anchor registry are kit-guidance files, and the guidance ticket carries the shipped-surface claim words the package-core guard reads.

The shim's trust chain is unchanged. The harness runs the shim. The shim sources the shared library and resolves the wrapper. The wrapper runs the core. This spec changes only the core's verdict.

Every subagent runs `opus` at low or medium effort. The review round runs `opus` at low effort with a cap of two iterations, as the reviewer named on 2026-09-04.
