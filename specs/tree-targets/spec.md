# Explicit tree targets for tree-scoped Bench verbs (spec A)

Status: staged

Roadmap: FT341

Decision source: `specs/tree-targets/decisions/tree-targets.md` (ready compiled map).

Verification log: 1 iteration(s) to accept — fable/high round 1 accepted after fixes; folded F1–F10 and reviewer decisions F2–F4.

## Problem

A Bench verb finds its tree from the current directory. The wrapper selects
the kit from that directory, and each verb derives its root through
`git.Root()`. An agent that works in a pool worktree must therefore run
`bench worktree exec <label> -- bench <verb>`, or it grades the primary checkout
by mistake. The FT341 occurrences record wrong-tip reviews, red preflights
that were green, and shell variables that hold pool paths.

A linked worktree has no untracked binary. The wrapper runs the primary
checkout's executable, so registries compiled into Go come from the primary
source and not from the worktree. No response states which tree it read, so a
wrong-tree grade looks like a real verdict.

## Solution

Each public registry leaf declares `tree` or `repository`. The gate refuses a
public leaf with no declaration. A tree-scoped verb takes `--in <label|primary>`
as its first argument. It then runs as one child in that tree target, and root
help shows the flag on each tree-scoped row.

Each tree-scoped response starts with one `tree[1]{target,head,dirty}` row. In
the kit repository, a worktree target runs its own worktree build, and a
missing or stale build refuses before any child starts. The kit-source
predicate takes its canonical spelling from `internal/canonicalpath`.

This spec is spec A of the map. It keeps directory inference, and it keeps the
exec route for Bench commands. Spec B owns the exec refusal, the end of
directory inference, and the migration of each hook and skill. When spec A
retires, its retirement promotes map tickets 1, 2, 6, and 7 onto the FT341 row
as the residual work for spec B.

## User stories

Line: opus / high.
Implementation-line reason: TT-C4 is the hardest chunk, because it starts a child executable, routes its environment, and refuses before candidate code runs. The spec fixes each output byte and each refusal. The child route has one real seam uncertainty, and the gate only observes it through the system suite.
Harder chunks: TT-C4.

The scope declaration:

1. As a maintainer, I want each public registry leaf to declare `tree` or `repository`, so that one field owns the target grammar.
2. As a maintainer, I want the gate to refuse a public leaf with no scope, so that no new verb ships undeclared.
3. As a maintainer, I want the gate to refuse a scope on a plumbing leaf, so that hook and gate internals stay outside.
4. As an agent, I want a repository-scoped verb to refuse `--in` at exit 2, so that an ignored target never looks accepted.
5. As an agent, I want root help to show the `--in` form on each tree-scoped row only, so that I see which verbs take it.

The identity row:

6. As an agent, I want each tree-scoped response to start with one `tree` row, so that I know which tree the verb read.
7. As an agent, I want the row to name `primary` for the primary checkout, so that I can tell the primary checkout from a worktree.
8. As an agent in a Bench worktree, I want the row to name the assignment label, so that I can confirm the worktree I meant.
9. As an agent in an unowned linked worktree, I want the row to name `unassigned`, so that the row never guesses a label.
10. As an agent, I want the row to hold the HEAD commit or `none`, so that I can pin the graded tip.
11. As an agent, I want the row to state `dirty` as `true`, `false`, or `unknown`, so that uncommitted changes are visible.
12. As an agent, I want an artifact, terminal, or ship-tier call to print the row on stderr, so that the artifact bytes stay clean.
13. As an agent, I want a help answer to carry no row, so that the grammar reply stays unchanged.
14. As an agent, I want a repository-scoped response to carry no row, so that only a tree read names a tree.
15. As an agent, I want the row to print no pool path, so that no caller copies a pool path into shell text.

The tree target:

16. As an agent, I want `bench <verb> --in <label>` to run the verb in that worktree, so that I grade the worktree that I name.
17. As an agent in a pool worktree, I want `--in primary` to target the primary checkout, so that I can read `main`.
18. As an agent, I want a missing `--in` value to exit 2, so that a typo starts no child.
19. As an agent, I want a path operand to `--in` to exit 2 with no child, so that a pool path never becomes a target.
20. As an agent, I want an unknown label to refuse and name `bench worktree list`, so that I can find the right label.
21. As an agent, I want a label of an inactive assignment to refuse with its state, so that I never grade a retired tree.
22. As an agent, I want an ambiguous label to refuse and name each colliding assignment, so that no guess selects a tree.
23. As an agent, I want a control character in the `--in` value to refuse in escaped form, so that the reply stays one safe line.
24. As an agent, I want `--in` in any other argument position to reach the verb's own grammar, so that one position owns the flag.
25. As an agent, I want the child's exit code and output to pass through unchanged, so that the response is the verb's own.
26. As an agent in a linked project repository, I want the child to run the installed executable, so that my repository needs no worktree build.
27. As an agent, I want a relative operand after `--in <label>` to resolve in that worktree, so that a repo-relative path names the target's file.

The running executable:

28. As a maintainer, I want a kit worktree target to run its own worktree build, so that compiled registries match the tree.
29. As a maintainer, I want a missing worktree build to refuse and name `bench worktree build <label>`, so that I know the repair.
30. As a maintainer, I want a stale worktree build to refuse before candidate code runs, so that stale registries never grade the tree.
31. As a maintainer, I want a build whose bytes differ from its seal to refuse, so that a changed executable never runs.
32. As a maintainer, I want the child to inherit no run binary and no kit, so that its gate builds its own run binary.

The path derivation:

33. As a maintainer, I want the kit-source predicate to use `internal/canonicalpath`, so that one derivation owns the canonical path.
34. As a maintainer, I want the child directory to come from `internal/canonicalpath`, so that a symlinked pool path runs in one physical tree.

Reviewed exclusions:

35. As a reviewer, I want `bench worktree exec <label> -- bench <verb>` to keep working in this spec, so that no route breaks before spec B migrates the guidance.
36. As a reviewer, I want directory inference kept in this spec, so that only spec B changes the default target.
37. As a reviewer, I want the exec refusal, the inference change, and the migration left out, so that spec B owns them.
38. As a reviewer, I want the `bench worktree list` filter by spec to stay out of this spec, so that FT125 owns it.

Added in review round 1:

39. As an agent, I want a grammar refusal at exit 2 to print only its usage line, so that the usage answer stays unchanged.
40. As an agent, I want an empty `--in` value to exit 2 as a missing value, so that an empty label never reaches the lookup.
41. As an agent, I want an `--in` value that starts with `-` to exit 2, so that a misplaced flag never reads as a label.

## Implementation decisions

### Late answer from the reviewer, 2026-09-29

This run authors spec A only. Spec B takes its decision source from the FT341
row after spec A retires. Map ticket 7 fixes the split.

### The scope field

The command registry gains one scope type with two values, `tree` and
`repository`. Its zero value means undeclared. The type sits on the command
definition and on the leaf row of a command family. The dispatcher reads it once
for each call. At run time, an undeclared definition prints no row and takes no
tree target. So a planted test registry keeps its current output.

Three rules fix where a declaration sits:

- A public definition that is not a family declares exactly one scope. A public definition has the public inventory visibility.
- A family definition declares no scope of its own. Each of its leaf rows declares one. `worktree` is the only family, and every `worktree` leaf is repository-scoped.
- A plumbing definition declares no scope. A plumbing definition has the internal inventory visibility. `gate-prose`, the hook verbs, and the gate internals are plumbing.

The classification predicate is map ticket 2. A leaf is tree-scoped when it
reads or grades tracked content of one checkout. A leaf that writes tracked
content of one checkout is tree-scoped too, by the reviewer's decision of
2026-09-29. A leaf is repository-scoped when it reads or writes only Git refs,
the worktree ledger, git-ignored capture files, or state under the Bench home.
The Bench-home clause extends map ticket 2, and Flagged additions lists it.

`retro` is tree-scoped, because it writes the tracked retrospective under
`capture/retros`. `capture` stays repository-scoped, because it writes only
git-ignored inboxes and the drain state.

| scope | public leaves |
| --- | --- |
| tree | `anchors`, `canary`, `commit`, `consumers`, `coverage`, `dashboard`, `diff`, `gate`, `guards`, `handoff`, `harnesses`, `init`, `link`, `maps`, `outline`, `preflight`, `prep-release`, `probe`, `release`, `release-preflight`, `retro`, `roadmap`, `setup`, `skills-index`, `spec`, `status`, `structure`, `test`, `unlink`, `upgrade` |
| repository | `assessment`, `cache`, `capture`, `commands`, `doctor`, `help`, `idea`, `learning`, `learnings`, `models`, `repair`, `repair-pilot`, `shift`, `version`, and each `worktree` leaf |

The command registry is the one source of this classification. This table is
the review surface for it, and no test restates it.

### The scope check

The existing `subcommand-routing` conformance check grades the declarations.
It already parses the command registry, and its input source is `go-source`.
The check reports one diagnostic for each broken rule:

- `command "<name>" declares no scope` for a public definition that is not a family.
- `worktree leaf "<name>" declares no scope` for a leaf row with no scope.
- `command family "<name>" declares a scope` for a family definition with a scope.
- `plumbing command "<name>" declares a scope` for a plumbing definition with a scope.

The registry parser moves out of `subcommand_routing_test.go` into a new file,
and it gains a reader for the `worktreeLeaves` table. The move keeps the old file
under its current line count.

### The repository refusal

A repository-scoped definition that is not a family refuses `--in` as its first
argument. The refusal prints `toon.Usage("bench <name>", "--in")` on stdout at
exit 2, and the verb does not run. A family routes its first argument as a leaf
name, so `bench worktree --in <x>` keeps its unknown-leaf answer.

A wrapper-only definition, such as `repair`, never reaches the dispatcher. The
wrapper answers `--in` with its own usage line on stderr at exit 2, and the verb
does not run. This form is the reviewer's decision of 2026-09-29.

### The identity row

The row renders through `toon.Table` as the block `tree[1]{target,head,dirty}:`
with one row. It is the first block of each tree-scoped response that resolves a
repository root and exits with a code other than 2. The dispatcher computes the
row before the verb runs. The bounded response owner writes it at finish as the
first block, and only when the exit is not 2. So a spilled response keeps the
row as its first inline line.

A grammar refusal at exit 2 prints only its usage line, as it does today. This
rule is the reviewer's decision of 2026-09-29, and TT56 pins it.

The two row lines do not count toward the response bound. At finish, the owner
prints the row before the replayed or projected output. When the spill file
cannot open, the owner has already flushed the verb output, so it prints no
row. This rule is the reviewer's decision of 2026-09-29, and TT60 and TT61 pin
it.

A Bench verb can start a Bench child, as the gate starts `bench test`. That child
prints its own row, because each child call is its own response. So a
parent response can grow past the bound and spill. This rule is the reviewer's
decision of 2026-09-29.

The cells have these values:

- `target` is `primary` when the root is the primary checkout by `git.IsPrimaryCheckout`. Otherwise it is the label of the active assignment that `intent.AssignmentForWorktree` returns. Otherwise it is `unassigned`.
- `head` is the full commit that `git.ResolveCommit(root, "HEAD")` returns, or `none` when HEAD does not resolve.
- `dirty` is `true` or `false` from `git.WorktreeDirty`, or `unknown` when that query fails.

A call whose bound disposition is exempt prints the row on stderr after the verb
returns, and only when the exit is not 2. The exempt calls are
`dashboard --stdout`, `setup`, and the ship tier.
A help form prints no row. A call outside a repository prints no row, and the
verb keeps its own not-in-repo answer.

### The tree target

`--in` is a tree target only as the first argument after a tree-scoped verb. In
any other position, the verb's own parser receives it. The dispatcher consumes
`--in` and its value, and it resolves the value in this order:

1. A missing or empty value exits 2 with `toon.MissingArg("bench <name> --in", "<label|primary>")` on stdout.
2. A value that starts with `-`, such as `--help`, exits 2 with `toon.Usage("bench <name> --in", <value>)` on stdout.
3. A value with a control character refuses at exit 1 through the worktree target refusal.
4. The keyword `primary` names the primary checkout, the first entry of `git.Worktrees`.
5. The lookup compares the value with the label of every ledger row. Exactly one active match names that worktree. With no active match, exactly one match in another state gives the state refusal. Two or more active matches, or two or more inactive matches with no active match, give the ambiguity refusal. This order is the reviewer's decision of 2026-09-29.
6. A value that matches no label and that `targetPath` in `internal/worktree` reads as a path exits 2 with `toon.Usage("bench <name> --in", <value>)` on stdout. `targetPath` reads each of these values as a path shape: an absolute path, `~`, `~user`, a `~/` prefix, `.`, and each value with `/`. A `~user` value gives a path refusal.
7. Any other value refuses at exit 1 through the worktree target refusal, with the verb `bench <name> --in`.

Step 6 reuses `targetPath` and restates none of its rules. The label resolver in
`internal/worktree` calls it after the label lookup fails, and it answers a typed
path-shape error that `internal/treetarget` renders.

The label lookup matches exact labels only. It uses no id, no prefix, and no
path. It keeps the state check, the missing-tree check, and the creation-bundle
check of the current worktree resolver. A new exported function in
`internal/worktree` owns this lookup, so the target authority stays in one
package.

The dispatcher then runs the verb as one child. The child argv is the selected
executable, the verb name, and the remaining arguments. The child directory is
`canonicalpath.Resolve` of the target root. The child environment is the exec
child environment of `internal/worktree`, so the child inherits no
`BENCH_KIT`, no `BENCH_RUN_BINARY`, and no `PWD`.

The parent prints no row and applies no response bound. The child prints its own
row by directory inference, and it bounds its own response. The parent returns
the child's exit code and passes its output through. A child failure prints no
`worktree:` path line, so a new exported runner in `internal/worktree` shares
the exec child runner without that line.

A new package, `internal/treetarget`, owns the flag, the value order, the
identity row, and the executable choice. It imports `internal/worktree`,
`internal/intent`, `internal/git`, `internal/freshness`,
`internal/canonicalpath`, `internal/toon`, and `internal/sanitize`. Only
`cmd/bench` imports it.

### The running executable

A label target whose root declares Bench build inputs is a kit worktree target.
`freshness.DeclaresBuildInputs` is the predicate. For a kit worktree target, the
child executable is `freshness.PublishedExecutable` of the target root, after
`freshness.Verify` accepts it against that root.

A refusal prints two stderr lines at exit 1, and no child starts:

- `bench <name> --in: worktree build is missing` when the executable is absent, or `bench <name> --in: worktree build does not match the tree` for any other `Verify` refusal.
- `next=bench worktree build <label>`, with the label shell-quoted when it needs quoting.

Neither line prints the executable path.

Every other target runs the invoking wrapper that `BENCH_WRAPPER` names in the
parent's environment. With no wrapper, it runs the running executable. The
`primary` target never runs a worktree build.

### Bootstrap authority

The refusal-before-execution claim has these hops:

1. The operator's shell starts the invoking wrapper and the running executable. Both are trusted as they are today.
2. The running executable resolves the label through the ledger and reads the target root.
3. The running executable computes the source digest of the target root and the digest of the worktree build. It compares both with the adjacent seal through `freshness.Verify`.
4. Only on acceptance does it start the worktree build as the child.

The seal is candidate-controlled, so it detects staleness and a changed
executable. It does not detect forgery. The trust assumption for reviewer
review is this: a worktree build is candidate code that the operator chose to
run, as `bench worktree exec <label> -- ./dist/bench` runs it today. No gate and
no lane trusts it. The child's environment drops `BENCH_RUN_BINARY`, so its gate
builds its own run binary from the graded tree.

### The path derivation

The gate's `resolvedPath` helper delegates to `canonicalpath.Resolve`, so
`KitSourceCheckout` compares absolute canonical spellings. A resolve error
answers no match. A relative kit directory now matches the root it names, where
the old helper compared an unabsolute spelling.

For a relative path, `canonicalpath.Resolve` makes the path absolute before it
resolves symlinks. So a working directory that a process entered through a
symlink gives the physical spelling, not the symlink spelling. Ticket 4 owns
this change, by the reviewer's decision of 2026-09-29, and TT62 pins it.

### The help rows

Root help renders each row of a tree-scoped definition as
`bench <name> [--in <label|primary>]<suffix>`. A repository-scoped row and a
plumbing definition render as they do today. The help row renderer reads the
scope field, so the help and the dispatcher share one source.

### The structure budget

`cmd/bench/command_registry.go` holds 389 lines against a budget of 400. The
scope type, the refusal, the row hook, and the help insertion live in a new
file, `cmd/bench/tree_scope.go`. The old file only gains the two fields and
the calls into that file.

`cmd/bench/main.go`, `cmd/bench/main_test.go`,
`internal/conformance/subcommand_routing_test.go`, and
`internal/conformance/axi_query_registry_test.go` are over budget already. An
edit there keeps each file at or under its current line count. Each registry
line in `main.go` takes its scope inline.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| TT-C1 / `1-declare-command-scope.md` | Each public leaf declares its scope, the gate refuses an undeclared leaf, and a repository verb refuses `--in`. | TT1, TT2, TT3, TT4, TT5, TT6, TT7 | `bench test --package ./cmd/bench`, `bench test --check subcommand-routing` | no |
| TT-C2 / `2-derive-kit-source-path.md` | The kit-source predicate uses the canonical derivation. | TT51, TT52 | `bench test --package ./internal/gate` | no |
| TT-C3 / `3-name-the-graded-tree.md` | Each tree-scoped response names the tree that it read. | TT10, TT11, TT12, TT13, TT14, TT15, TT16, TT17, TT18, TT19, TT20, TT21, TT22, TT23, TT24, TT25, TT56, TT60, TT61 | `bench test --package ./cmd/bench`, `bench test --package ./internal/treetarget` | no |
| TT-C4 / `4-run-verbs-in-tree-target.md`, `5-run-kit-worktree-build.md` | A tree-scoped verb runs in a named tree target, and a kit worktree target runs its own current build. | TT8, TT9, TT26, TT27, TT28, TT29, TT30, TT31, TT32, TT33, TT34, TT35, TT36, TT37, TT38, TT39, TT40, TT41, TT57, TT58, TT59, TT43, TT44, TT45, TT46, TT47, TT48, TT49, TT50, TT53, TT54, TT55, TT62 | `bench test --package ./internal/treetarget`, `bench test --package ./internal/worktree`, `bench test --package ./internal/canonicalpath`, `bench test --package ./cmd/bench`, `bench test --check system` | yes |

TT-C1 creates the scope field that TT-C3 and TT-C4 consume, so its review
closes first. TT-C3 creates the identity row that each TT-C4 child prints.

## Testing decisions

- A good test drives the real dispatcher, `Command{}.Run`, in a temporary repository and compares the printed bytes and the exit code.
- `internal/treetarget` tests drive its run function with an explicit executable path. A script records its argv, its directory, and its environment in a marker file. These tests swap no package variable.
- TT9 and TT25 plant a registry, so they swap the package variable `commandRegistry` in the test process. The venue is `Command{}.Run` in that same process, so the swap reaches the dispatcher. `help_inventory_test.go` and `response_bound_test.go` are the precedents.
- The kit rows use a Go-module fixture with a `./cmd/bench` package and the auxiliary build-input manifest `scripts/go-build.inputs`. The manifest makes `freshness.DeclaresBuildInputs` true, and the digest covers its entries. `commandsBriefCheckout` in `cmd/bench/commands_brief_test.go` is the precedent, and it was read in this session.
- The kit rows publish a script as the worktree build through `freshness.Publish`, so the seal is real. A stale row edits one listed build input after the publication.
- Each refusal row asserts that the marker file is absent, so a refusal after the child starts fails the row.
- The system rows run the real built executable against a scaffolded repository with real assignments. `internal/systemtest/exec_bound_test.go` is the prior art, and it was read in this session.
- The conformance rows plant a registry source and assert the exact diagnostic, the way `TestAXIRegistryParserFailsClosed` does.
- The gate observes the package rows in its `test` phase and the system rows in its `system` phase.

### Posture change: tests that the new output reds

The identity row adds a first block to each tree-scoped response through the
dispatcher that exits with a code other than 2. The enumeration below covers
each fixture the row reds. Its needles are each literal `.Run([]string{"<verb>"`
call, each `Run(nil)` call, and each call of the four dispatch helpers in
`cmd/bench` tests. The system needle is each exact or prefix comparison of a
tree-scoped stdout in `internal/systemtest`.

Four helpers run the in-process dispatcher and return stdout. Ticket 3 adds one
test helper that removes a leading `tree[1]{target,head,dirty}:` block and its
row, and each of the four helpers calls it in place:

- `runAXICommandAsAt` in `cmd/bench/spill_support_test.go`, which `runAXICommandAt` calls.
- `dispatch` in `cmd/bench/census_output_test.go`.
- `runKeptRoute` in `cmd/bench/command_registry_test.go`, edited in place at its current line count.
- `runPreflight` in `cmd/bench/preflight_version_test.go`.

That one edit covers each helper call site. The sites include
`anchor_help_test.go` lines 76 to 254, `anchors_dir_test.go` lines 94 to 138,
`axiEnvelopeRows` in `command_registry_test.go` lines 133 to 232, and
`spec_retire_listing_test.go`. They also include `worktree_leaves_test.go`,
`help_inventory_test.go` lines 195 to 237, `response_bound_test.go` line 132,
and `selected_queries_test.go` lines 40 to 135.

These direct sites compare an exact or prefix stdout, and ticket 3 edits each in
place:

- `cmd/bench/main_test.go` line 87, the `Run(nil)` prefix of `status --route`.
- `cmd/bench/main_test.go` line 122, the prefix of the built binary's root route.
- `cmd/bench/main_test.go` line 163, the prefix of the built binary's `harnesses` probes, when the response does not spill.
- `cmd/bench/commit_chain_test.go` lines 79 and 136, the exact stdout of `commit` through `runCommitChain`.
- `cmd/bench/selected_queries_test.go` line 183, the exact line count of `spec history`.
- `cmd/bench/isolated_command_fixtures_test.go` line 39, the prefix of `status --route`.
- `internal/systemtest/owner_selection_test.go` lines 250 and 251, the exact `canary` stdout.
- `internal/systemtest/charge_evidence_test.go` lines 228 and 344, the prefixes of `preflight evidence`.
- `internal/systemtest/charge_evidence_test.go` line 318, which compares the reads of two worktrees and the `page[1]` prefix. Both comparisons drop the row.
- `cmd/bench/response_bound_exempt_test.go` lines 93 and 145, where `requireComplete` demands an empty stderr, and the planted `dashboard` and ship-tier calls now print the exempt row there. The bounded `dashboard` case also gains the two row lines.
- `internal/systemtest/adoption_test.go` line 150 and later, where the scaffolded gate's nested Bench calls each print a row, so the gate response can spill.

These sites stay green, and ticket 3 does not edit them:

- `cmd/bench/gate_route_test.go` lines 60 to 66 and line 98 take exit 2, so they print no row.
- The help forms at `commit_chain_test.go` line 149, `help_inventory_test.go` line 134, and `main_test.go` lines 200 and 264 print no row.
- `main_test.go` line 54 and `isolated_command_fixtures_test.go` line 58 compare substrings.
- `main_test.go` lines 125 and 129 and `charge_evidence_test.go` lines 341 and 386 compare two reads of one tree.
- `status_route_converge_test.go` line 50 finds its command by a substring index.
- The `land_route_test.go`, `owner_land_race_test.go`, and `version` comparisons run repository-scoped verbs.

### Seam diagram

    trigger: `bench <verb> [--in <label|primary>] <args>`
        │
        ▼
    argv, current directory  ──▶  [ dispatcher: scope → tree target → child or row ]  ──▶  stdout, stderr, exit code
                      ◀ tests attach here: `Command{}.Run` in a temporary repository,
                        the `internal/treetarget` run function with a marker script,
                        and the real executable in the system suite

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| TT1 | 1 | The live registry declares one scope on each public non-family definition and on each `worktreeLeaves` row | `internal/conformance/gate_entry_test.go` (`TestRootConformance`), through the `subcommand-routing` check over the live tree | A registry line without a scope reds the check on the live tree. |
| TT2 | 2 | A planted registry with one public definition and no scope reports `command "x" declares no scope` | planned TestCommandScopeCheckBites in internal/conformance, through a planted registry source | A check that reads only family leaves prints no diagnostic for this input. |
| TT3 | 2 | A planted `worktreeLeaves` row with no scope reports `worktree leaf "y" declares no scope` | planned TestCommandScopeCheckBites in internal/conformance, through a planted leaf table | A check that reads only `commandRegistry` misses the leaf table. |
| TT4 | 3 | A planted internal-inventory definition with a scope reports `plumbing command "z" declares a scope` | planned TestCommandScopeCheckBites in internal/conformance, through a planted registry source | A check that only demands presence accepts a plumbing scope. |
| TT5 | 2 | A planted family definition with a scope reports `command family "f" declares a scope` | planned TestCommandScopeCheckBites in internal/conformance, through a planted registry source | A check that treats a family as a plain definition accepts the extra scope. |
| TT6 | 4 | `bench version --in primary` exits 2, stdout is `usage: bench version (unknown argument: --in)`, and no version line prints | planned TestRepositoryVerbRefusesTreeTarget in cmd/bench, through `Command{}.Run` | `versionCommand` ignores its arguments and prints the version at exit 0 today. |
| TT7 | 4 | `bench idea --in primary x` exits 2 and leaves `capture/IDEAS.md` absent | planned TestRepositoryVerbRefusesTreeTarget in cmd/bench, through `Command{}.Run` in a temporary repository | A refusal after the verb runs leaves the parked line in the inbox. The `idea` grammar already refuses `--in`, so TT6 and the refusal test over every repository verb carry the bite. |
| TT8 | 5 | Root help renders each tree-scoped row with the `--in` insertion that the help-rows section states | `cmd/bench/help_inventory_test.go` (`TestHelpInventoryIsComplete`), with its expectation changed in place | The independent expectation holds the insertion, so a row without it fails the whole-text match. |
| TT9 | 5 | A registry with one tree and one repository definition renders the insertion on the tree row only | planned TestHelpRendersTreeTargetFromScope in cmd/bench, through a planted registry | A renderer that inserts on every row, or on none, fails one side. |
| TT10 | 6 | In a primary checkout, the first stdout block of `bench roadmap` is `tree[1]{target,head,dirty}:` with one row | planned TestTreeRowLeadsTreeResponse in cmd/bench, through `Command{}.Run` in a temporary repository | A row after the verb output, or no row, fails the prefix match. |
| TT11 | 7 | That row is `primary,<HEAD commit>,false` | planned TestTreeRowLeadsTreeResponse in cmd/bench, through `Command{}.Run` in a temporary repository | A target that reads the branch name prints `main`. |
| TT12 | 8 | `Identify` of a worktree that an active assignment labeled `alpha` owns answers target `alpha` | planned TestIdentifyNamesActiveLabel in internal/treetarget, through a ledger written by `intent.PutAssignment` | A target that reads the directory name prints the pool directory. |
| TT13 | 9 | `Identify` of a linked worktree that no assignment owns answers target `unassigned` | planned TestIdentifyUnownedWorktree in internal/treetarget, through `git worktree add` | A fallback to `primary` misnames the tree. |
| TT14 | 9 | `Identify` of a worktree whose assignment is released answers target `unassigned` | planned TestIdentifyUnownedWorktree in internal/treetarget, through a released ledger row | A lookup that ignores the state prints the stale label. |
| TT15 | 10 | The `head` cell equals the 40-character HEAD commit | planned TestTreeRowLeadsTreeResponse in cmd/bench, through `Command{}.Run` in a temporary repository | An abbreviated commit fails the length match. |
| TT16 | 10 | A repository with no commit renders `head` as `none` | planned TestIdentifyUnbornHead in internal/treetarget, through `git init` with no commit | A failed resolve that prints an empty cell renders an empty field. |
| TT17 | 11 | An untracked file renders `dirty` as `true` | planned TestIdentifyDirtyStates in internal/treetarget, through a temporary repository | A query that reads only tracked changes prints `false`. |
| TT18 | 11 | A modified tracked file renders `dirty` as `true` | planned TestIdentifyDirtyStates in internal/treetarget, through a temporary repository | A dirty value that reads only untracked files prints `false`. |
| TT19 | 11 | A corrupt index file renders `dirty` as `unknown` | planned TestIdentifyDirtyStates in internal/treetarget, through a corrupt `.git/index` | A query error read as clean prints `false`. |
| TT20 | 12 | `bench dashboard --stdout` prints no `tree[` line on stdout, and its stderr carries the `tree[1]{target,head,dirty}:` block | planned TestExemptTreeCallPrintsRowOnStderr in cmd/bench, through `Command{}.Run` in a temporary repository | A row on stdout corrupts the HTML artifact and fails the stdout match. |
| TT21 | 13 | `bench gate --help` prints no `tree[` line | planned TestHelpFormPrintsNoTreeRow in cmd/bench, through `Command{}.Run` | A row hook that runs before the help check prints the row. |
| TT22 | 14 | `bench version` prints no `tree[` line | planned TestHelpFormPrintsNoTreeRow in cmd/bench, through `Command{}.Run` | A row hook on every verb prints the row. |
| TT23 | 15 | The row for the `alpha` worktree holds no byte of the worktree path | planned TestIdentifyNamesActiveLabel in internal/treetarget, through a ledger written by `intent.PutAssignment` | A target that prints the path fails the substring check. |
| TT24 | 6 | Outside a repository, `bench coverage x` prints no `tree[` line and keeps its not-in-repo answer | planned TestTreeRowOutsideRepository in cmd/bench, through `Command{}.Run` in a temporary directory | A row hook that prints before the root resolves prints an empty row. |
| TT25 | 6 | A planted tree-scoped verb that prints 40 lines spills, and its first inline line is the `tree[1]{target,head,dirty}:` header | planned TestTreeRowSurvivesSpill in cmd/bench, through a planted registry | A row printed outside the bound owner lands after the spill line or nowhere. |
| TT26 | 16 | In a linked fixture, `bench status --in alpha` from the primary checkout exits 0, and its row is `alpha,<worktree HEAD>,false` | planned TestTreeTargetRunsInNamedWorktree in internal/systemtest, through the real executable | A dispatcher that ignores `--in` prints `primary` and the primary HEAD. |
| TT27 | 16 | `Run` with target `alpha` starts the given executable with argv `status --route` in the canonical `alpha` root | planned TestRunStartsChildInTarget in internal/treetarget, through a marker script | A child that keeps `--in` in its argv, or runs in the parent directory, fails the marker match. |
| TT28 | 17 | With its directory inside `alpha`, `bench status --in primary` prints the row `primary,<primary HEAD>,false` | planned TestTreeTargetRunsInNamedWorktree in internal/systemtest, through the real executable | A dispatcher that keeps the current directory prints `alpha`. |
| TT29 | 18 | `bench gate --in` exits 2, stdout is the step 1 `toon.MissingArg` line of the tree-target section, and no marker exists | planned TestRunRefusesBeforeChild in internal/treetarget, through a marker script | A parser that reads the next absent argument as a label starts a child. |
| TT30 | 19 | `--in` with an absolute pool path exits 2 with `usage: bench gate --in (unknown argument: <path>)`, and no marker exists | planned TestRunRefusesBeforeChild in internal/treetarget, through a marker script | A resolver that accepts paths starts the child. |
| TT31 | 19 | `--in ./alpha` exits 2 with the same usage form, and no marker exists | planned TestRunRefusesBeforeChild in internal/treetarget, through a marker script | A path test that checks only absolute paths accepts the relative form. |
| TT59 | 19 | `--in ~` exits 2 with the same usage form, and no marker exists | planned TestRunRefusesBeforeChild in internal/treetarget, through a marker script | A path test that skips `targetPath` misses the home form and gives the unassigned refusal. |
| TT57 | 40 | `bench gate --in ""` exits 2 with the step 1 `toon.MissingArg` line, and no marker exists | planned TestRunRefusesBeforeChild in internal/treetarget, through a marker script | A lookup of the empty label gives the unassigned refusal at exit 1. |
| TT58 | 41 | `bench gate --in --help` exits 2, stdout is `usage: bench gate --in (unknown argument: --help)`, and no marker exists | planned TestRunRefusesBeforeChild in internal/treetarget, through a marker script | A resolver without the dash step gives the unassigned refusal at exit 1. |
| TT32 | 19 | A value `team/alpha` that equals an active label starts the child in that worktree | planned TestRunLabelWinsOverPathShape in internal/treetarget, through a marker script | A path-shape test that runs before the label lookup refuses a real label. |
| TT33 | 20 | An unknown label exits 1, stderr is `bench gate --in: target is unassigned` then `next=bench worktree list`, and no marker exists | planned TestRunRefusesBeforeChild in internal/treetarget, through a marker script | A lookup that falls back to the primary checkout starts a child. |
| TT34 | 21 | A released label exits 1, stderr names the released state, and no marker exists | planned TestRunRefusesBeforeChild in internal/treetarget, through a released ledger row | A lookup without the state check starts a child in a retired tree. |
| TT35 | 22 | Two active assignments labeled `alpha` exit 1, stderr names both ids, and no marker exists | planned TestRunRefusesBeforeChild in internal/treetarget, through two ledger rows | A lookup that takes the first match starts a child. |
| TT36 | 23 | A value with U+0001 exits 1, stderr carries `target contains control characters`, and stderr holds no raw U+0001 | planned TestRunRefusesBeforeChild in internal/treetarget, through a marker script | A refusal that prints the raw value writes the control byte. |
| TT37 | 24 | `bench gate --fresh --in primary` exits 2 with the gate usage and starts no child | planned TestTreeTargetOnlyAsFirstArgument in cmd/bench, through `Command{}.Run` in a temporary repository | A dispatcher that scans every position consumes the late flag and runs the gate. |
| TT38 | 25 | A child that prints `x` and exits 3 gives the parent exit 3, stdout `x`, and stderr with no `worktree:` line | planned TestRunPassesChildResult in internal/treetarget, through a marker script | The exec child runner prints a `worktree:` path line on a failure. |
| TT39 | 26 | A non-kit target with `BENCH_WRAPPER` set starts that wrapper as the child | planned TestRunSelectsChildExecutable in internal/treetarget, through a marker script as the wrapper | A selection that always takes the running executable bypasses the wrapper's kit resolution. |
| TT40 | 26 | A non-kit target with no `BENCH_WRAPPER` starts the running executable | planned TestRunSelectsChildExecutable in internal/treetarget, through a marker script as the running executable | A selection that demands a wrapper refuses a direct binary run. |
| TT41 | 27 | `bench coverage --in alpha specs/only-in-alpha/spec.md` reads the spec that exists only in `alpha` | planned TestTreeTargetRunsInNamedWorktree in internal/systemtest, through the real executable | A child in the parent directory reports the spec path as missing. |
| TT43 | 28 | A kit worktree target with a published build starts `<root>/dist/bench` as the child | planned TestRunKitWorktreeBuild in internal/treetarget, through a script published by `freshness.Publish` | A selection that ignores the build inputs runs the wrapper. |
| TT44 | 29 | A kit worktree target with no build exits 1, stderr is `bench status --in: worktree build is missing` then `next=bench worktree build alpha`, and no marker exists | planned TestRunKitWorktreeBuild in internal/treetarget, through a marker script | A selection that falls back to the wrapper starts a child. |
| TT45 | 30 | After one listed build input changes, the target exits 1 with `worktree build does not match the tree` and `next=bench worktree build alpha`, and no marker exists | planned TestRunKitWorktreeBuild in internal/treetarget, through a script published by `freshness.Publish` | A check of the executable digest only starts the stale child. |
| TT46 | 31 | After one byte is appended to the published build, the target exits 1 with the same two lines, and no marker exists | planned TestRunKitWorktreeBuild in internal/treetarget, through a script published by `freshness.Publish` | A check of the source digest only starts the changed child. |
| TT47 | 29 | For the label `my alpha`, the refusal prints `next=bench worktree build 'my alpha'` | planned TestRunKitWorktreeBuild in internal/treetarget, through a marker script | A remedy without quoting splits the label into two operands. |
| TT48 | 32 | A kit child sees no `BENCH_RUN_BINARY`, no `BENCH_KIT`, and `BENCH_WRAPPER` set to `<root>/bin/bench.sh` | planned TestRunKitWorktreeBuild in internal/treetarget, through a script published by `freshness.Publish` | A child that inherits the parent environment carries the parent's run binary. |
| TT49 | 28 | A target without build inputs and with a `dist/bench` present runs the wrapper and not `dist/bench` | planned TestRunSelectsChildExecutable in internal/treetarget, through two marker scripts | A selection that reads only the presence of `dist/bench` runs a linked repository's stray file. |
| TT50 | 30 | In the system suite, a kit worktree with a stale build refuses `bench status --in alpha` and starts no child | planned TestTreeTargetRefusesStaleKitBuild in internal/systemtest, through the real executable | A dispatcher route that skips `freshness.Verify` in the real binary runs the stale build. |
| TT51 | 33 | With `BENCH_KIT` set to `.` and the current directory at the root, `KitSourceCheckout(root)` answers true | planned TestKitSourceCheckoutResolvesARelativeKit in internal/gate, through `t.Chdir` | The old helper compares the unabsolute spelling `.` with the root and answers false. |
| TT52 | 33 | A symlinked root spelling still matches the kit, and an unrelated directory does not | `internal/gate/kit_source_test.go` (`TestKitSourceCheckoutMatchesThroughASymlinkSpelling`) | A delegation that drops the symlink step answers false for the link. |
| TT53 | 34 | With the Bench home reached through a symlink, the child's `PWD` equals `canonicalpath.Resolve` of the recorded worktree path | planned TestRunStartsChildInTarget in internal/treetarget, through a symlinked home | A child directory that keeps the recorded spelling prints the symlink path. |
| TT54 | 35 | `bench worktree exec alpha -- <executable> status` exits 0, and its row names `alpha` | planned TestTreeTargetRunsInNamedWorktree in internal/systemtest, through the real executable | A dispatcher that refuses an exec child breaks the route this spec keeps. |
| TT55 | 36 | The real executable run with its directory inside `alpha` and no `--in` prints the row target `alpha` | planned TestTreeTargetRunsInNamedWorktree in internal/systemtest, through the real executable | A default that already targets the primary checkout prints `primary`. |
| TT56 | 39 | `bench gate --brief` exits 2 with an empty stdout and the gate usage on stderr | `cmd/bench/gate_route_test.go` (`TestRunGateRejectsBriefUsage`) | A row that the owner writes at every finish makes stdout non-empty. |
| TT60 | 6 | A planted tree-scoped verb that prints exactly 10 lines does not spill, and its stdout is the two row lines and then the 10 verb lines | planned TestTreeRowOutsideResponseBound in cmd/bench, through a planted registry | A row that counts toward the bound spills a response of 10 verb lines. |
| TT61 | 6 | When the spill file cannot open, the response carries the verb output and no `tree[` line | planned TestOwnerPrintsNoRowWhenSpillCannotOpen in internal/responsebound, through an unwritable spill directory | A row written after the passing flush lands after the verb output. |
| TT62 | 34 | With `BENCH_KIT` set to `.` and the working directory entered through a symlink to the root, `KitSourceCheckout` of the physical root answers true | planned TestResolveRelativeUnderSymlinkedWorkingDirectory in internal/canonicalpath, through `t.Chdir` to a symlink | A resolve that makes the path absolute after the symlink step keeps the symlink spelling and answers false. |

Not covered: story 37 — spec B owns these behaviors, and map ticket 6 holds the decision.
Not covered: story 38 — FT125 owns the filter, and map ticket 7 excludes it.

### Edge inventory

Audience: the scope field, the row, and `--in` serve every repository that
links the kit. The running-executable rule serves this repository only, because
only a kit source declares Bench build inputs. TT49 fixes the linked answer.

Absent versus empty: a repository with no ledger answers `unassigned` for an
unowned worktree, by TT13. An empty ledger gives an unknown label the unassigned
refusal, by TT33. A missing `dist/bench` refuses by TT44.

Handled edges, each with its row:

- The unborn HEAD: TT16.
- The failed status query: TT19.
- The exempt artifact call: TT20.
- The spilled response: TT25.
- The call outside a repository: TT24.
- The late `--in` position: TT37.
- The label with `/`: TT32.
- The label that needs shell quoting: TT47.
- The released assignment: TT14 and TT34.
- The label collision: TT35.
- The control character: TT36.
- The symlinked home: TT53.
- The home-form path `~`: TT59.
- The empty value `--in ""`: TT57, where `toon.MissingArg` wins over the lookup.
- The value that starts with `-`, such as `--in --help`: TT58, where `toon.Usage` at exit 2 wins over the lookup.
- The grammar refusal at exit 2 with no row: TT56.

The ordering promise is TT32: the label lookup wins over the path shape. A
missing or empty value wins over every lookup, by TT29 and TT57. A value that
starts with `-` wins over the lookup, by TT58.

Tests that swap a package variable: TT9 and TT25 swap `commandRegistry` in the
test process, and their venue is `Command{}.Run` in that process. Each
`internal/treetarget` test passes its executable and its home as arguments, and
it swaps no package variable.

Hostile input, shell CLI surface: the `--in` value is a caller-controlled
string. A control character refuses in escaped form through `sanitize.Controls`,
by TT36. A tab, a newline, or a return is a control character there. A label
cell with a comma prints through `toon.Table`, which quotes it. The label in the
repair command passes through `sanitize.ShellQuote`, by TT47. The value
never reaches a shell, because the child argv is a direct exec.

**Won't handle:**

- A tree target by assignment id or by prefix — map ticket 1 limits the value to a label or `primary`. `bench worktree exec <id> -- bench <verb>` stays the surviving route in spec A, and TT54 keeps it.
- A label spelled `primary` — the keyword wins. `bench worktree exec <id> -- bench <verb>` stays its surviving route, by TT54.
- Staleness of the primary checkout's own executable — the map excludes it. The current freshness route stays its caller.
- A seal that a candidate forges — the seal detects staleness only. The gate's own run binary stays the authority, by TT48.
- A `canary` root operand in the row — the row names the current tree, not the operand, and the operand keeps its current meaning. `bench canary` with no operand, or with `--in <label>`, is the supported route, by TT26.

## Ownership fences

- `cmd/bench/tree_scope.go`
- `cmd/bench/tree_scope_test.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/main.go`
- `cmd/bench/worktree_leaves.go`
- `cmd/bench/help_inventory_test.go`
- `cmd/bench/main_test.go`
- `cmd/bench/selected_queries_test.go`
- `cmd/bench/isolated_command_fixtures_test.go`
- `cmd/bench/spill_support_test.go`
- `cmd/bench/census_output_test.go`
- `cmd/bench/preflight_version_test.go`
- `cmd/bench/commit_chain_test.go`
- `cmd/bench/census_output.go`
- `internal/responsebound/owner.go`
- `internal/responsebound/owner_test.go`
- `cmd/bench/response_bound_exempt_test.go`
- `internal/conformance/subcommand_routing_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/command_registry_parse_test.go`
- `internal/conformance/command_scope_test.go`
- `internal/gate/kit_source.go`
- `internal/gate/kit_source_test.go`
- `internal/treetarget/`
- `internal/worktree/tree_target.go`
- `internal/worktree/tree_target_test.go`
- `internal/canonicalpath/canonicalpath.go`
- `internal/canonicalpath/canonicalpath_test.go`
- `internal/worktree/exec.go`
- `internal/systemtest/`
- `tests/canary/package-core-guard/unrouted-subcommand`
- `reviews/tree-targets.md`

## Ticket graph

| ticket | blocked by | chunk |
| --- | --- | --- |
| `1-declare-command-scope.md` | none | TT-C1 |
| `2-derive-kit-source-path.md` | none | TT-C2 |
| `3-name-the-graded-tree.md` | `1-declare-command-scope.md` | TT-C3 |
| `4-run-verbs-in-tree-target.md` | `3-name-the-graded-tree.md` | TT-C4 |
| `5-run-kit-worktree-build.md` | `4-run-verbs-in-tree-target.md` | TT-C4 |

## Out of scope

- Spec B: the exec refusal for a Bench child, the end of directory inference, and the migration of each hook and skill. Map tickets 1, 2, 6, and 7 hold its decisions. Estimate: 40 edits, 4 gate runs.
- A `bench worktree list` filter by spec: FT125 owns precise readers. Estimate: 5 edits, 1 gate run.
- A tree target by assignment id: the map limits the value to a label or `primary`. Estimate: 3 edits, 1 gate run.

## Further notes

### Source trace

| source sentence | rows |
| --- | --- |
| #1: `--in` accepts a worktree label or `primary`, and a path operand exits 2 | TT26 to TT36, TT57 to TT59 |
| Reviewer decision F2, 2026-09-29: no row on an exit-2 grammar refusal | TT56 |
| Reviewer decision F3, 2026-09-29: a leaf that writes tracked content is tree-scoped | the classification table, review-owned |
| Reviewer decision F4, 2026-09-29: the `canary` root operand stays, and the row names the current tree | Won't handle, TT26 |
| #1: `--in` follows the verb leaf, and a repository-scoped leaf that gets `--in` exits 2 | TT6, TT7, TT37 |
| #1: no `--in` targets the primary checkout after spec B | story 37, Not covered |
| #1: no environment variable selects a tree target | TT48 |
| #2: each public leaf declares `tree` or `repository`, and a conformance check refuses a missing declaration | TT1 to TT3, TT5 |
| #2: the tracked-content predicate | the classification table, review-owned |
| #2: the lifecycle verbs stay repository-scoped with their positional `<target>` | TT1, TT54 |
| #2: plumbing leaves stay outside the contract | TT4 |
| #2: the `--in` grammar and help derive from the field | TT8, TT9 |
| #3: a kit worktree target runs its worktree build | TT43, TT48 |
| #3: a missing or stale build refuses and names `bench worktree build <label>` | TT44 to TT47, TT50 |
| #3: a linked repository keeps the installed executable | TT39, TT40, TT49 |
| #4: one `tree{target,head,dirty}` row and no pool path | TT10 to TT25 |
| #5: target resolution uses `internal/canonicalpath`, and `resolvedPath` becomes its caller | TT51 to TT53 |
| #6: spec A keeps directory inference and the exec route | TT54, TT55 |
| #6: spec B adds the exec refusal and removes inference in one landing | story 37, Not covered |
| #7: two specs, and spec A moves every tree-scoped leaf | TT1, TT8 |
| #7: the list filter is out of scope | story 38, Not covered |

### Reader sweep and proof checklist

Readers of the command registry fields and of the help rows:

- `cmd/bench/command_registry.go`: `Command.Run`, `renderCommandHelp`, `leafRoot`, and `dispatchLeafFamily`.
- `cmd/bench/census_output.go`: `runBounded` opens the bound owner that the row writes into.
- `internal/responsebound/owner.go`: `Owner.Finish` writes the retained or projected response, so the row reaches it at finish.
- `internal/conformance/subcommand_routing_test.go`: `parseCommandRegistry` and `dispatchRegistry`.
- `internal/conformance/axi_query_registry_test.go`, `entry_point_parity_test.go`, and `help_inventory_single_source_test.go` call `parseCommandRegistry`. The move keeps its name and signature, so they do not change.
- `cmd/bench/help_inventory_test.go`: `TestHelpInventoryIsComplete` pins the whole root help. `internal/conformance/help_inventory_single_source_test.go` reads source bytes, and no row changes them.
- The guidance files that name `bench coverage <spec>` and `bench gate` spell invocations, not help rows. Two root help rows are also pinned in `cmd/bench/main_test.go` lines 137 and 193. Ticket 4 edits both in place at the current line count. The search found no other copy of a root help row outside `CHANGELOG.md` history and dated audit inputs.
- `.bench/hooks`, `.bench/adapters`, `bin`, `scripts`, and `.github/workflows` run no tree-scoped verb whose stdout they parse. The session-start hook runs the plumbing verb `session-inspect`.

Proof checklist:

- Cited symbols: each symbol resolves in the tree at `f981cd3d`.
  - In `cmd/bench`: `commandDefinition`, `commandLeaf`, `Command.Run`, `renderCommandHelp`, `runBounded`, `worktreeLeaves`, `commandRegistry`, `publicInventory`, `internalInventory`, `versionCommand`, and `boundExemptWith`.
  - In `internal/worktree`: `resolveAssignmentIn`, `selectAssignment`, `printTargetRefusal`, `runWorktreeChild`, `execEnv`, `nameWorktree`, `errTargetUnassigned`, and `errTargetControls`.
  - In other packages: `git.Root`, `git.IsPrimaryCheckout`, `git.Worktrees`, `git.ResolveCommit`, `git.WorktreeDirty`, `intent.AssignmentForWorktree`, `intent.AssignmentsOwning`, `freshness.DeclaresBuildInputs`, `freshness.Verify`, `freshness.Publish`, `freshness.PublishedExecutable`, `canonicalpath.Resolve`, `gate.KitSourceCheckout`, `toon.Table`, `toon.Usage`, `toon.MissingArg`, `sanitize.ShellQuote`, and `sanitize.Controls`.
- Import edges: `internal/treetarget` is new, and only `cmd/bench` imports it. `internal/gate` already imports `internal/canonicalpath` through `subject.go`. `internal/worktree` already imports `internal/intent`, `internal/freshness`, and `internal/sanitize`.
- Source-row clauses and occurrences: the source trace table above. The FT341 occurrences motivate the problem, and each one is a wrong-tree or exec-route incident.
- Promised field labels: `tree{target,head,dirty}` and the cell values `primary`, `unassigned`, `none`, `true`, `false`, and `unknown`.
- Changed-function callers: `resolvedPath` has one caller, `KitSourceCheckout`. `KitSourceCheckout` has callers in `cmd/bench/command_registry.go`, in `internal/adopt`, and in `internal/worktree/joins.go` line 119, and their answers change only for a relative kit spelling. `runWorktreeChild` has one caller, `execAttributed`, and it keeps the `worktree:` line.
- Copy survival: TT51 fails when `resolvedPath` keeps its own spelling.
- Rendered-shape readers: the new first block reaches the dispatcher-level tests under the posture change. Ticket 3 searches each `Command{...}.Run` call in `cmd/bench` tests that names a tree-scoped verb without a help form. It also searches each exact stdout comparison in `internal/systemtest`.

Sources re-read in the authoring session: every `## Sources` entry of the map,
and no source was left unread. `roadmap/FT341.md` still carries the destination
and the occurrences. `bin/bench.sh` still holds `kit_dir`, `main_tree_kit`,
`bench_binary_path`, and `route_binary`. `internal/worktree/exec.go` still holds
`execEnv`, which drops the wrapper routing.

`internal/freshness/freshness.go` still holds `Digest` over `BuildInputs`.
`.bench/gate.sh` still passes its root as an operand and its binary through
`BENCH_RUN_BINARY`. `internal/git/git.go` still holds `Root` and `RootAt`.

### Fence disposition

Reviewer disposition of the ownership fences: open.

The fence holds `internal/systemtest/` as a prefix, because ticket 3 sweeps its
exact comparisons and tickets 4 and 5 add the system rows. The fence holds no
file of `internal/intent`, `internal/freshness`, or `internal/canonicalpath`,
because the build only calls them.

The build preflight binds four more paths to the `cmd/bench` and
`internal/worktree` files. They are `cmd/bench/command_registry_test.go`, the two
conformance registry tests, and the unrouted-subcommand canary fixture. They
are in the fence for that closure only. `cmd/bench/command_registry_test.go` is
over budget, so an edit there keeps its current line count.

The fence holds the four dispatch-helper files, because ticket 3 edits each
helper once. It holds no helper call site, such as `anchor_help_test.go` or
`anchors_dir_test.go`, because the helper edit covers each call site. It holds
no `gate_route_test.go`, because its calls exit 2 and print no row.

The staged spec `specs/ft290-test-projection/spec.md` also fences
`cmd/bench/main.go`, `cmd/bench/command_registry.go`,
`cmd/bench/command_registry_test.go`, `cmd/bench/help_inventory_test.go`,
`internal/conformance/axi_query_registry_test.go`,
`internal/conformance/subcommand_routing_table_test.go`, and
`tests/canary/package-core-guard/unrouted-subcommand`. None of its work has
landed. The build that starts second takes the
mid-tier staleness audit that `.agents/commands/bench-implement-spec.md` line 19
requires over those files. The reviewer recommends that one of the two builds
lands before the other dispatches.

### Completion plan

```bench-completion-plan
{"version":2,"execution":{"mode":"delegate","run_id":"tree-targets-full-20260929","orchestrator_session":"claude:session-08f8b392-1af7-4e9e-9671-d6d921f3325b","author_limit":1,"assignments":{"1-declare-command-scope.md":[{"session":"claude:bench-writer/tt-t1-author","assignment":"tt-t1-author","model":"opus","effort":"high","source":"f981cd3db1a4f27eddba5feede40a80e205bf1c1","native_ref":"claude:agent/tt-t1-author-20260929@f981cd3db1a4f27eddba5feede40a80e205bf1c1"},{"session":"claude:bench-writer/tt-t1-repair-1","assignment":"tt-t1-repair-1","model":"opus","effort":"high","source":"f84ddebca7406bd6f7ea1e33e391235fc4e344e2","native_ref":"claude:agent/tt-t1-repair-1-20260929@f84ddebca7406bd6f7ea1e33e391235fc4e344e2","predecessor":"claude:bench-writer/tt-t1-author","trigger":"user-directed","stopped":"The author returned its final report after record commit a08359f9 and holds no write.","preserved":"f84ddebca7406bd6f7ea1e33e391235fc4e344e2"}],"2-derive-kit-source-path.md":[{"session":"claude:bench-writer/tt-t2-author","assignment":"tt-t2-author","model":"opus","effort":"high","source":"f84c016ec166fe2fcb4a8e0691248ece461d04ee","native_ref":"claude:agent/tt-t2-author-20260929@f84c016ec166fe2fcb4a8e0691248ece461d04ee"},{"session":"claude:bench-writer/tt-t2-author-2","assignment":"tt-t2-author-2","model":"opus","effort":"high","source":"5685e1878a659a8f6fdeb2d2ece1f073f651a989","native_ref":"claude:agent/tt-t2-author-2-20260929@5685e1878a659a8f6fdeb2d2ece1f073f651a989","predecessor":"claude:bench-writer/tt-t2-author","trigger":"terminal-failure","stopped":"The author returned a blocked report because the plan probe 2-kit-probe could not compile, and it committed nothing.","preserved":"5685e1878a659a8f6fdeb2d2ece1f073f651a989"},{"session":"claude:bench-writer/tt-t2-repair-1","assignment":"tt-t2-repair-1","model":"opus","effort":"high","source":"54e5b7f1653d6e81945636d29e93956577004604","native_ref":"claude:agent/tt-t2-repair-1-20260929@54e5b7f1653d6e81945636d29e93956577004604","predecessor":"claude:bench-writer/tt-t2-author-2","trigger":"user-directed","stopped":"The author returned its final report after record commit c8b9bb94 and holds no write.","preserved":"54e5b7f1653d6e81945636d29e93956577004604"}],"3-name-the-graded-tree.md":[{"session":"claude:bench-writer/tt-t3-author","assignment":"tt-t3-author","model":"opus","effort":"high","source":"b4f4d260410334dc937af7631e36d1b43a9c7a68","native_ref":"claude:agent/tt-t3-author-20260929@b4f4d260410334dc937af7631e36d1b43a9c7a68"},{"session":"claude:bench-writer/tt-t3-author-2","assignment":"tt-t3-author-2","model":"opus","effort":"high","source":"c050dc968e92de7d937a8263f08b8b084a5040f5","native_ref":"claude:agent/tt-t3-author-2-20260929@c050dc968e92de7d937a8263f08b8b084a5040f5","predecessor":"claude:bench-writer/tt-t3-author","trigger":"terminal-failure","stopped":"The author returned a blocked report for a fence gap and a behavioral question, and it committed nothing.","preserved":"c050dc968e92de7d937a8263f08b8b084a5040f5"}],"4-run-verbs-in-tree-target.md":[],"5-run-kit-worktree-build.md":[]}},"chunks":[{"id":"TT-C1","tickets":["1-declare-command-scope.md"],"verification":[{"id":"1-cmd","command":"bench test --package ./cmd/bench","ticket":"1-declare-command-scope.md"},{"id":"1-routing","command":"bench test --check subcommand-routing","ticket":"1-declare-command-scope.md"},{"id":"1-conformance","command":"bench test --package ./internal/conformance","ticket":"1-declare-command-scope.md"},{"id":"1-scope-probe","command":"bench probe cmd/bench/tree_scope.go --swap '\"--in\"' --with '\"--in-x\"' --package ./cmd/bench --run TestRepositoryVerbRefusesTreeTarget","probe":"swap","ticket":"1-declare-command-scope.md"}]},{"id":"TT-C2","tickets":["2-derive-kit-source-path.md"],"verification":[{"id":"2-gate","command":"bench test --package ./internal/gate","ticket":"2-derive-kit-source-path.md"},{"id":"2-kit-probe","command":"bench probe internal/gate/kit_source.go --swap 'canonicalpath.Resolve' --with 'func(p string) (string, error) { _ = canonicalpath.Resolve; return filepath.EvalSymlinks(p) }' --package ./internal/gate --run TestKitSourceCheckoutResolvesARelativeKit","probe":"swap","ticket":"2-derive-kit-source-path.md"}]},{"id":"TT-C3","tickets":["3-name-the-graded-tree.md"],"verification":[{"id":"3-cmd","command":"bench test --package ./cmd/bench","ticket":"3-name-the-graded-tree.md"},{"id":"3-treetarget","command":"bench test --package ./internal/treetarget","ticket":"3-name-the-graded-tree.md"},{"id":"3-responsebound","command":"bench test --package ./internal/responsebound","ticket":"3-name-the-graded-tree.md"},{"id":"3-unassigned-probe","command":"bench probe internal/treetarget/identify.go --swap '\"unassigned\"' --with '\"primary\"' --package ./internal/treetarget --run TestIdentifyUnownedWorktree","probe":"swap","ticket":"3-name-the-graded-tree.md"}]},{"id":"TT-C4","tickets":["4-run-verbs-in-tree-target.md","5-run-kit-worktree-build.md"],"verification":[{"id":"4-treetarget","command":"bench test --package ./internal/treetarget","ticket":"4-run-verbs-in-tree-target.md"},{"id":"4-worktree","command":"bench test --package ./internal/worktree","ticket":"4-run-verbs-in-tree-target.md"},{"id":"4-canonicalpath","command":"bench test --package ./internal/canonicalpath","ticket":"4-run-verbs-in-tree-target.md"},{"id":"4-cmd","command":"bench test --package ./cmd/bench","ticket":"4-run-verbs-in-tree-target.md"},{"id":"4-system","command":"bench test --check system","ticket":"4-run-verbs-in-tree-target.md"},{"id":"5-treetarget","command":"bench test --package ./internal/treetarget","ticket":"5-run-kit-worktree-build.md"},{"id":"5-system","command":"bench test --check system","ticket":"5-run-kit-worktree-build.md"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/tree-targets/spec.md"},{"id":"cmd","command":"bench test --package ./cmd/bench"},{"id":"treetarget","command":"bench test --package ./internal/treetarget"},{"id":"worktree","command":"bench test --package ./internal/worktree"},{"id":"gate","command":"bench test --package ./internal/gate"},{"id":"routing","command":"bench test --check subcommand-routing"},{"id":"system","command":"bench test --check system"}]}
```

### Flagged additions

Each addition below is not in a map ticket answer. The reviewer can veto each one.

- The target cell `unassigned` for a linked worktree that no active assignment owns, TT13 and TT14. Spec A keeps directory inference, so the row needs a value for that tree.
- The cell values `none` for an unborn HEAD and `unknown` for a failed status query, TT16 and TT19.
- The row on stderr for an exempt call, TT20, so that an HTML artifact stays byte-clean.
- The TOON table form `tree[1]{target,head,dirty}:` for the one row.
- The `--in` position rule: only the first argument after the verb, TT37.
- The label-before-path-shape order, TT32, so that a real label with `/` stays addressable.
- No `worktree:` path line on a child failure, TT38.
- The refusal texts `worktree build is missing` and `worktree build does not match the tree`, which print no executable path.
- `freshness.DeclaresBuildInputs` as the kit worktree predicate, because `KitSourceCheckout` compares a root with the running kit and answers false for any worktree.
- The scope check inside `subcommand-routing`, not a new registered check, because a new check pulls the anchored profile table into the fence.
- The classification of `status` and `handoff` as tree-scoped. The FT341 row names both as repository-scoped, but map ticket 2's predicate classifies both as tree readers.
- The classification of `harnesses` as tree-scoped, because it reads a registry compiled from the kit source.
- The Bench-home clause of the repository predicate. It classifies `assessment`, `cache`, `models`, and `repair-pilot`, which read only state under the Bench home.
- The empty-value and dash-value steps of the value order, TT57 and TT58.
- The one test helper that removes a leading row inside the four dispatch helpers.
- The run-time rule for a definition with no declared scope: no row and no tree target. The staleness pass added it, because planted test registries declare no scope.

### Open reviewer decision

A label collision has no `--in` address. `bench worktree create` refuses no
duplicate label, and map ticket 1 accepts only a label or `primary`. In spec A,
`bench worktree exec <id> -- bench <verb>` resolves the collision, by TT35 and
TT54. Spec B removes that route, so spec B needs a collision route. The two
options are an assignment id as a third `--in` form, or a duplicate-label refusal
at create. This spec changes neither, and it records the question for spec B.

The reviewer leans to a duplicate-label refusal at `bench worktree create`,
because map ticket 1 closed the `--in` value set.

### Reviewer decisions during the build

The build found five behavioral gaps, and the reviewer decided each one on
2026-09-29. Each rule now lives in its spec section.

1. The row and the response bound: the two row lines do not count toward the
   bound. The identity-row section holds the rule, and TT60 and TT61 pin it.
2. The label lookup precedence: exactly one active match wins. Step 5 of the
   tree-target section holds the rule, and TT34 and TT35 pin it.
3. The wrapper-only `repair` verb: the wrapper's own refusal stands. The
   repository-refusal section holds the rule.
4. The symlinked working directory: ticket 4 fixes `canonicalpath.Resolve`.
   The path-derivation section holds the rule, and TT62 pins it.
5. The nested Bench call: each child prints its own row. The identity-row
   section holds the rule.
