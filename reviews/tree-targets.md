# Tree targets review record

## TT-C1 author evidence

Ticket 1 had a fresh `bench-writer` author, `tt-t1-author`, on opus at high effort, with a cap of 3 attempts. The author started at `1634be1e` and committed `166a3a8f` on a lane pass in the first attempt. The author then committed this record in a second commit.

The chunk pair is `f981cd3d..166a3a8f`. The base is the `main` tip, because a plan commit is never a chunk base. The payload names the ticket commit as the tip, because a record cannot name its own commit. The coordinator moves the tip to the record commit. The source digest does not change, because the record file is outside the graded source.

The author wrote each test before the code that it grades. The red and green log for each row follows:

- TT6: before the first production edit, `TestRepositoryVerbRefusesTreeTarget/version` failed. The verb printed `bench dev (linux/amd64)` at exit 0. After the refusal in `cmd/bench/tree_scope.go`, the subtest passed.
- TT7: before the change, the `idea` subtest passed, because the idea grammar refuses `--in` itself. A probe made the idea parser skip its first two arguments, and the subtest failed with `parked: x` at exit 0. After the change, the same probe was silent, because the dispatcher refuses before the verb runs.
- TT2 to TT5: with an empty check, `TestCommandScopeCheckBites` failed with no scope diagnostic. After the check, the test passed.
- TT1: after the check and before the registry declarations, `TestRootConformance` failed with `gate: command "anchors" declares no scope`. The named check with the worktree build gave the same red. After the declarations, both passed.

The author reported these deviations from the ticket. The Spec axis grades each one.

- The scope check reads a classified entry that is not internal as public. The help-inventory check refuses any third classification. So the scope check reuses `parityInternalCommand` and adds no second test for `publicInventory`. This change follows the duplicated-facts sweep.
- The author's check read the presence of a `Scope` field and did not grade its value. The author's reason was false: Go accepts an untyped `0` for the scope type, so a line can set the undeclared zero value. Repair R2 grades the value.
- A tree without the leaf file grades the registry only. The routing check applies the same rule to an absent package, and the canary fixture copies only `cmd/bench/main.go`.
- Each registry line takes its scope after its bound disposition. So the AXI mutation anchors still match, and `internal/conformance/axi_query_registry_test.go` stays unchanged.
- No check required an edit to these `Writes:` paths, so the diff leaves them unchanged:
  - `cmd/bench/command_registry_test.go` and `cmd/bench/help_inventory_test.go`;
  - `internal/conformance/axi_query_registry_test.go` and `internal/conformance/subcommand_routing_table_test.go`.
- The canary fixture change is in `MUTATE.json`, because that file quotes the `learnings` registry line byte for byte.
- The TT7 red probe changed `internal/roadmap/roadmap.go`, which is outside the fence. The probe restored the file, and the diff does not touch it.

The author computed the plan digest from the plan rule. The input is the JSON array of the spec bytes and the five ticket bytes, in plan order. The same rule gives the digest `48ff3619` for the session context queries record at `a5a527c8`, which confirms the method. The source digest is the tree of `166a3a8f` without this record file.

### Probe verdicts

Each probe ran through `bench probe`, and each restore reads `yes`. The first row is the plan probe `1-scope-probe`. The JSON payload holds the exact command and output of the first three rows. The last two rows ran on the source before the ticket commit or record the silent result, so only this table holds them.

| File | Mutation | Test | Row | Verdict |
|---|---|---|---|---|
| `cmd/bench/tree_scope.go` | swap: `"--in"` to `"--in-x"` | TestRepositoryVerbRefusesTreeTarget | TT6 | bit |
| `internal/conformance/command_scope_test.go` | swap: `if len(leaf.fields["Scope"]) == 0 {` to `if len(leaf.fields["Scope"]) < 0 {` | TestCommandScopeCheckBites | TT3 | bit |
| `internal/conformance/command_scope_test.go` | swap: `case len(entry.fields["Inventory"]) == 1 && !declared:` to `case len(entry.fields["Inventory"]) == 2 && !declared:` | TestCommandScopeCheckBites | TT2 | bit |
| `internal/roadmap/roadmap.go` | swap: `usage.Parse(ideaGrammar, args)` to `usage.Parse(ideaGrammar, args[min(2, len(args)):])`, before the change | TestRepositoryVerbRefusesTreeTarget/idea | TT7 | bit |
| `internal/roadmap/roadmap.go` | the same swap, after the change | TestRepositoryVerbRefusesTreeTarget/idea | TT7 | silent |

### Verification

The author ran each TT-C1 verification on the source of `166a3a8f`, and each passed. The JSON payload holds each result. The conformance excerpt omits its three skip rows. Each skip is an environment capability skip for unix sockets or device nodes. The author also ran these checks, and each passed: `bench structure --growth f981cd3d`, `bench canary`, and the named checks `canary-fixture-compliance`, `package-core-guard`, `conformance-canary-families`, and `axi-query-registry`. After the commit, `bench preflight build tree-targets` reported 14 green checks, 1 check that does not apply, and 0 red checks.

## TT-C1 chunk review, round 1

The frozen pair is base `f981cd3db1a4f27eddba5feede40a80e205bf1c1` and tip `a08359f9550f64d92d524dfebc3f7e6143a31118`. The coordinator moved the chunk tip from the ticket commit to the record commit, and the source digest stays the same. The shared evidence is `sha256:bd1140010578a62d29360a6a33d4647e322f10a032106ee66ede731b1edd3df3`. Each axis ran in a fresh `bench-reviewer` session on opus at high effort, on the conditional review line. Only the Coverage axis ran probes, and it left the tree clean. The whole-project gate was green at the frozen tip.

The raw finding count is 2: Standards 1, Spec 0, and Coverage 1. The Spec axis gave the Coverage finding as advice. The coordinator adds one target that the chunk shares with every other author of this ticket: no test drives the refusal over each repository verb. So 3 repair targets remain.

## Standards

Findings: 1. The worst issue is a new private copy of a fixture-writing helper.

- `internal/conformance/command_scope_test.go:66-75` adds a `write` closure that does the work of `writeFixtureFile` at `internal/conformance/package_core_checks_test.go:295`, in the same package. `AGENTS.md` names a fixture harness pasted N times as duplicated knowledge. Target R1. `auto-fix`. Confidence 5.

## Spec

Findings: 0. All seven rows TT1 to TT7 are met, and the classification matches the spec table. The axis gave two items as advice. The first is the zero-value scope gap, which the Coverage axis files as a finding. The second is the TT7 premise below.

## Coverage

Findings: 1. The worst issue is that the scope check grades only the presence of the `Scope` field.

- `internal/conformance/command_scope_test.go:26` and `:52` test only that the key exists. A probe that set `Scope: 0` on the `models` line in `cmd/bench/main.go` stayed silent under `--check subcommand-routing` and under `./cmd/bench`. Spec lines 120-121 say "Its zero value means undeclared", and line 128 says a public definition "declares exactly one scope". The author's reason, that an undeclared value needs an explicit conversion, is false, because Go accepts an untyped `0`. Target R2, `auto-fix`, confidence 6.

## TT-C1 coordinator additions

- R3: `TestRepositoryVerbRefusesTreeTarget` drives only `version` and `idea`. In another author's tree for this ticket, a probe that narrowed the refusal to `version` stayed silent across `./cmd/bench`. With that mutation, `bench shift --in primary x` would start a loop with the objective `--in primary x`. The same two-verb test is in this tree. `auto-fix`.
- The TT7 premise is a non-behavioral spec contradiction. The `idea` grammar already refuses `--in` at exit 2, so the row's "why it catches" clause cannot occur. The coordinator follows the tree convention: TT7 keeps its behavior, TT6 and R3 carry the bite, and the spec row states this. The reviewer can veto this reading.

## TT-C1 repair routing

Each repair goes to one fresh `bench-writer` repair session for ticket 1 on opus at high effort. Ticket 1's `Writes:` line holds every path. This is cycle 1 of the two repair cycles for chunk TT-C1.

| Target | Ticket | Repair |
|---|---|---|
| R1 | 1 | Use the package's `writeFixtureFile` helper instead of the new closure. |
| R2 | 1 | Accept only `scopeTree` and `scopeRepository` as a declaration, and plant a `Scope: 0` row in `TestCommandScopeCheckBites`. |
| R3 | 1 | Drive the refusal over every repository-scoped definition that is not a family, read from `commandRegistry`, and show that a refusal narrowed to `version` bites. |

## TT-C1 ticket 1 repair evidence, cycle 1

The session `claude:bench-writer/tt-t1-repair-1` ran on opus at high effort, with a cap of 3 attempts. It started at `56caf792` and committed `258edfd9` on a lane pass in the first attempt. This repair is cycle 1 of the two repair cycles for chunk TT-C1.

- R1: `TestCommandScopeCheckBites` writes both planted files through the package's `writeFixtureFile` helper. The private closure is gone, and the test file makes no other file write. This change has no behavior, so no probe applies. The R2 probes below show that the test still bites through the helper.
- R2: a public definition and a `worktreeLeaves` row declare a scope only when the `Scope` field names `scopeTree` or `scopeRepository`. A family or a plumbing line still carries no `Scope` field at all. The bite test plants `Scope: 0` on one public line and on one leaf row. Before the fix, the test failed without the two undeclared diagnostics. After the fix, it passed. The reviewer's live probe, `Scope: 0` on the `models` line, now reports `command "models" declares no scope`.
- R3: `TestRepositoryVerbRefusesTreeTarget` reads each repository-scoped definition from `commandRegistry` and skips a wrapper-only definition. For each verb, it runs `--in primary` and a bare `--in`, and it asserts exit 2, the one usage line on stdout, and an empty stderr. The test sets `BENCH_AGENT` empty, so `bench shift` fails fast when a mutation lets it run. The loop covers `version`, so its separate subtest is gone. The `idea` subtest stays for its file side effect.

The sweep of duplicated facts found one second statement. The scope check names the two scope identifiers, because a conformance check cannot import `package main`. The live `models` probe records the red of that independent expectation.

### Probe verdicts

Each probe ran through `bench probe` at the repair source, and each restore reads `yes`. The JSON payload holds the exact command and output of each probe. The first row is the plan probe `1-scope-probe`. The live R2 probe ran through the worktree build, because a named check compiles from the run binary's source.

| Target | File | Mutation | Test | Verdict |
|---|---|---|---|---|
| TT6 | `cmd/bench/tree_scope.go` | swap: `"--in"` to `"--in-x"` | TestRepositoryVerbRefusesTreeTarget, 10 failed | bit |
| R3 | `cmd/bench/tree_scope.go` | swap: `definition.Scope != scopeRepository ||` to `definition.Name != "version" ||` | TestRepositoryVerbRefusesTreeTarget, 8 failed, `shift` included | bit |
| R2 | `internal/conformance/command_scope_test.go` | swap: `return ok && declaredScopes[value.Name]` to `return value != nil || !ok` | TestCommandScopeCheckBites | bit |
| R2 | `cmd/bench/main.go` | swap: `Scope: scopeRepository` to `Scope: 0` on the `models` line | named check `subcommand-routing` | bit |

A first R2 predicate probe swapped the return to `return true`. It did not compile, so its verdict was `invalid`. The table holds the rerun.

### Verification

The session ran each TT-C1 plan verification on the source of `258edfd9`, and each passed. The JSON payload holds each result. The conformance excerpt omits its three skip rows. Each skip is an environment capability skip for unix sockets or device nodes. `bench structure --growth f981cd3d` reported that no source file grew past its budget. After the commit, `bench preflight build tree-targets` reported 14 green checks, 1 check that does not apply, and 0 red checks.

The chunk tip is now the repair commit `258edfd9`. The source digest is the tree of `258edfd9` without this record file. The same rule at `a08359f9` gives the round 1 digest `48cf94fb`, which confirms the method. The spec changed at `56caf792`, so the plan digest changed. The `ReadPlan` rule at `a08359f9` gives the round 1 digest `7355203a`, which confirms the method. The round 1 entries keep their earlier source digest as history.

## TT-C1 chunk review, round 2, and close

The frozen pair is base `f981cd3db1a4f27eddba5feede40a80e205bf1c1` and tip `196e1cda291a9d974f881435f6e1284b1838b241`. The shared evidence is `sha256:c994d11ef9b198c5abebdfd70863b3bb017e5b2e2bdb24b4785d9c73873a4199`. This round is the confirming round of all three axes. Each axis ran in a new fresh `bench-reviewer` session on opus at high effort, and each read only the repair delta `a08359f9..196e1cda`. Only the Coverage axis ran probes, and it left the tree clean.

## Standards

Findings: 0. R1, R2, and R3 are confirmed. The repair delta adds no duplicated knowledge.

## Spec

Findings: 0. All seven ticket 1 rows stay met, and the repair writes only fenced paths and this record.

## Coverage

Findings: 0. The R2 and R3 probes bite. The coordinator asked about a later `--in` on a repository verb. No test pins it, and a probe that also refuses a later `--in` stayed silent. The approved row for story 24 is TT37 in ticket 4, so this gap goes to the ticket 4 charge.

## Advice

- Ticket 4 adds a repository-verb case to `TestTreeTargetOnlyAsFirstArgument`, such as `bench version x --in primary`, which reaches the verb grammar.
- `TestCommandScopeCheckBites` pins only a literal `0`. A planted row with an undeclared identifier would also pin the membership test. No such constant exists today.

Chunk TT-C1 closes after one repair cycle.

## TT-C2 author evidence

Ticket 2 had a fresh `bench-writer` successor author, `tt-t2-author-2`, on opus at high effort, with a cap of 3 attempts. The first author, `tt-t2-author`, stopped because the plan probe could not compile, and it committed nothing. The coordinator replaced the plan probe in the plan commit `9100a74a`.

The successor started at `9100a74a` with the uncommitted diff of the first author. The successor verified that diff again with its own runs and kept it without a change. It committed `e43254d5` on a lane pass in the first attempt. The successor then committed this record in a second commit.

The chunk pair is `196e1cda..e43254d5`. The base is the accepted TT-C1 tip, because a plan commit is never a chunk base. The payload names the ticket commit as the tip, because a record cannot name its own commit. The source digest does not change at the record commit, because the record file is outside the graded source.

The production edit was already in the tree at the start. So each red came from a `bench probe` and not from a run before the edit. The red and green log for each row follows:

- TT51: a probe restored the exact old helper, `filepath.EvalSymlinks` with a `filepath.Clean` fallback. `TestKitSourceCheckoutResolvesARelativeKit` failed with `with BENCH_KIT=. = false, want true`, and the symlink test stayed green. The plan probe `2-kit-probe` gave the same red. After the restore, the test passed.
- TT52: the symlink test is unchanged. A probe that dropped the symlink step for the root failed it with `kit-link = false, want true`. A probe that made any two resolved paths match failed it with `consumer = true, want false`. After each restore, the test passed.

The author reported these deviations from the ticket. The Spec axis grades each one.

- `resolvedPath` returns a spelling and a success flag. So two failed resolves never compare equal as two empty spellings.
- No test drives the resolve-error branch. The only refusal of `canonicalpath.Resolve` is a failed working-directory read, and no acceptance row names that branch.
- The first post-commit preflight reported `binary-seal` red, because the worktree binary was older than the source. The author ran `bench worktree build`, and the second preflight was green.

The sweep of duplicated facts found no second statement. `kit_source.go` spells `canonicalpath.Resolve` once, and the new test uses its own fixture of three lines with no copied helper. Each independent expectation has a red above.

The author applied the `SourceDigest` and `ReadPlan` rules to the tree of `e43254d5`. The same steps at `196e1cda` give the TT-C1 digests `f0adccfe` and `82131a58`, which confirms the method. The plan changed at `5685e187` and at `9100a74a`, so the payload has one amendment from `82131a58` to `ab6c7f85`. The amendment maps each chunk ID to itself, because the plan delta changes only the ticket 2 author assignments and the ticket 2 plan probe.

### Probe verdicts

Each probe ran through `bench probe`, and each restore reads `yes`. The first row is the plan probe `2-kit-probe`. The JSON payload holds the exact command and output of the first three rows at `e43254d5`. The last two rows ran on the source before the ticket commit, so only this table holds them.

| File | Mutation | Test | Row | Verdict |
|---|---|---|---|---|
| `internal/gate/kit_source.go` | swap: `canonicalpath.Resolve` to a function that returns `filepath.EvalSymlinks(p)` | TestKitSourceCheckoutResolvesARelativeKit | TT51 | bit |
| `internal/gate/kit_source.go` | swap: `resolvedPath(root)` to `filepath.Clean(root), true` | TestKitSourceCheckoutMatchesThroughASymlinkSpelling | TT52 | bit |
| `internal/gate/kit_source.go` | swap: `return ok && resolvedRoot == resolvedKit` to `return ok && resolvedRoot+resolvedKit != ""` | TestKitSourceCheckoutMatchesThroughASymlinkSpelling | TT52 | bit |
| `internal/gate/kit_source.go` | swap: `canonicalpath.Resolve(path)` to the old helper body, before the commit | TestKitSourceCheckoutResolvesARelativeKit | TT51 | bit |
| `internal/gate/kit_source.go` | swap: `return ok && resolvedRoot == resolvedKit` to `return ok && resolvedRoot != ""`, before the commit | none | TT52 | invalid |

The last probe did not compile, because it left `resolvedKit` unused. The third row is its replacement.

### Verification

The author ran each TT-C2 plan verification on the source of `e43254d5`, and each passed. The JSON payload holds each result. The author also ran these checks, and each passed:

- `bench test --package ./internal/conformance --run TestBranchNativeArchitectureCensus`;
- `bench test --package ./internal/worktree --run TestParallelCensusOnTheLiveTree`;
- `bench test --check canonical-path-owner`;
- `bench structure --growth f981cd3d`, which reported that no source file grew past its budget.

A verbose run of each census test showed that the test ran and passed. After the commit and the worktree build, `bench preflight build tree-targets` reported 14 green checks, 1 check that does not apply, and 0 red checks.

## TT-C2 chunk review, round 1

The frozen pair is base `196e1cda291a9d974f881435f6e1284b1838b241` and tip `c8b9bb94eb15029a84abd137e1fbac93e8990480`. The shared evidence is `sha256:c69050cd5095eff839318373b2c019d536391f96ff3e71d7cbfef921a4ff6c30`. Each axis ran in a fresh `bench-reviewer` session on opus at high effort. Only the Coverage axis ran probes, and it left the tree clean.

The raw finding count is 2: Standards 1, Spec 0, and Coverage 1. The coordinator classifies the Standards item as advice, because row tags in test doc comments are the tree convention, as the TT-C1 review also found. The coordinator adds one target that every author of this ticket shares. So 2 repair targets remain.

## Standards

Findings: 0 after classification. The axis reported the `(TT51)` row tag in the new test doc comment at `internal/gate/kit_source_test.go:38`. More than 20 test files use row tags in this way, so the tag stays.

## Spec

Findings: 0. TT51 and TT52 are met, and the three recorded deviations are acceptable.

## Coverage

Findings: 1. The worst issue is that no test pins the resolve-error rule.

- Spec lines 276-277 say "A resolve error answers no match." A probe that swapped `return resolved, err == nil` for `_ = err; return resolved, true` in `internal/gate/kit_source.go` stayed silent. With a deleted working directory, `BENCH_KIT` set to `.`, and the root `.`, both resolves fail, and the mutant answers true. Target R1, `auto-fix`, confidence 6.

## TT-C2 coordinator additions

- R2: the new test pins only the positive relative case. In another author's tree for this ticket, a probe that resolved a relative kit against the root and not the working directory stayed silent. With `BENCH_KIT` set to `.`, every root would then count as the kit source. The same positive-only test is in this tree. `auto-fix`.

## TT-C2 repair routing

Each repair goes to one fresh `bench-writer` repair session for ticket 2 on opus at high effort. Ticket 2's `Writes:` line holds every path. This is cycle 1 of the two repair cycles for chunk TT-C2, and it is the chunk's one hardening cycle.

| Target | Ticket | Repair |
|---|---|---|
| R1 | 2 | Pin the resolve-error rule with a test that makes a resolve fail, and show that the probe above now bites. |
| R2 | 2 | Pin that a relative `BENCH_KIT` resolves against the working directory: with the directory elsewhere, `KitSourceCheckout(root)` answers false. |

## TT-C2 ticket 2 repair evidence, cycle 1

The session `claude:bench-writer/tt-t2-repair-1` ran on opus at high effort, with a cap of 3 attempts. It started at `fa840bcf` and committed `9958ef57` on a lane pass in the first attempt. This repair is cycle 1 of the two repair cycles for chunk TT-C2. The repair adds two tests to `internal/gate/kit_source_test.go` and changes no production code.

- R1: `TestKitSourceCheckoutAnswersFalseOnAResolveError` makes a directory, moves into it with `t.Chdir`, and removes it. Then `BENCH_KIT` is `.`, and `KitSourceCheckout(".")` must answer false. A host that cannot remove its working directory, or that still gives `.` an absolute spelling, skips through `capability.Environment`. On this Linux host the test runs, and a verbose run showed a pass with no skip. The review probe for R1 was silent in round 1. It now bites.
- R2: `TestKitSourceCheckoutResolvesARelativeKitAgainstTheWorkingDirectory` sets `BENCH_KIT` to `.` and moves to a second temporary directory. Then `KitSourceCheckout(root)` must answer false. A probe that joins only a relative kit onto the root fails this new test alone. The two earlier tests stay green under that probe, so the new test is necessary. `TestKitSourceCheckoutMatchesThroughASymlinkSpelling` is unchanged.

The sweep of duplicated facts found no second statement. Each new test builds its own fixture from `t.TempDir` and copies no helper. Each new expectation is independent, and the R1 and R2 probes below record its red.

### Probe verdicts

Each probe ran through `bench probe` at the source of `9958ef57`, and each restore reads `yes`. The JSON payload holds the exact command and output of each probe. The first row is the plan probe `2-kit-probe`, which the session ran again with the exact plan command.

| Target | File | Mutation | Test | Verdict |
|---|---|---|---|---|
| TT51 | `internal/gate/kit_source.go` | swap: `canonicalpath.Resolve` to a function that returns `filepath.EvalSymlinks(p)` | TestKitSourceCheckoutResolvesARelativeKit | bit |
| R1 | `internal/gate/kit_source.go` | swap: `return resolved, err == nil` to `_ = err; return resolved, true` | TestKitSourceCheckoutAnswersFalseOnAResolveError | bit |
| R2 | `internal/gate/kit_source.go` | swap: `kit := KitDir()` to a line that joins a relative kit onto the root | TestKitSourceCheckoutResolvesARelativeKitAgainstTheWorkingDirectory | bit |

A first R2 probe swapped `resolvedPath(kit)` to `resolvedPath(filepath.Join(root, kit))`. It bit, but it also failed the symlink test, because it joined an absolute kit too. The table holds the narrower probe, which fails only the new test.

### Verification

The session ran each TT-C2 plan verification on the source of `9958ef57`, and each passed. The JSON payload holds each result. The session also ran these checks, and each passed:

- `bench test --package ./internal/conformance --run TestBranchNativeArchitectureCensus`;
- `bench test --check skip-ownership`, because the R1 test can skip;
- `bench structure --growth f981cd3d`, which reported that no source file grew past its budget.

After the commit, `bench preflight build tree-targets` reported 14 green checks, 1 check that does not apply, and 0 red checks.

The chunk tip is now the repair commit `9958ef57`. The source digest is the tree of `9958ef57` without this record file, which is `80fbc1c6`. The same rule at `c8b9bb94` gives the round 1 digest `673171b5`, which confirms the method.

The spec changed at `fa840bcf`, so the plan digest changed from `ab6c7f85` to `834182d7`. The `ReadPlan` rule at `c8b9bb94` gives the earlier digest `ab6c7f85`, which confirms the method. The payload keeps the earlier amendment and adds one amendment from `ab6c7f85` to `834182d7`. That amendment maps each chunk ID to itself, because the plan change at `fa840bcf` changes only the ticket 2 assignments. The TT-C1 chunk keeps the digest of its own tip. The round 1 entries keep their earlier source digest as history.

## TT-C2 chunk review, round 2, and close

The frozen pair is base `196e1cda291a9d974f881435f6e1284b1838b241` and tip `74daf800f8a11d85083711b31aaabc47797ebc73`. The shared evidence is `sha256:51ffd8bccc3a331ce2f82d1887d076e60c37d5fd9c833c9323ddd5001b418730`. This round is the confirming round of all three axes. Each axis ran in a new fresh `bench-reviewer` session on opus at high effort, and each read only the repair delta `fa840bcf..74daf800`. Only the Coverage axis ran probes, and it left the tree clean.

## Standards

Findings: 0. R1 and R2 are confirmed, and both skips go through `capability.Environment`.

## Spec

Findings: 0. TT51 and TT52 stay met, and the repair delta stays inside the ticket 2 fence.

## Coverage

Findings: 0. Two independent probes, one for each fold, each bit on the new test alone.

## Advice

- The comment on the resolve-error test says both resolves fail. Only the root resolve runs, because the check returns before it resolves the kit.
- The skip-ownership diagnostic and the capability package describe the seam in two different ways. A later drain can decide which wording owns the rule.
- `canonicalpath.Resolve` keeps the symlink spelling of a working directory for a relative path. The spec records this as pending reviewer decision 4.

Chunk TT-C2 closes after one repair cycle. That cycle is the chunk's one hardening cycle.

The first checkpoint run refused the TT-C2 verification entries. Each entry held the digest of its full command output, not the digest of its embedded excerpt. The coordinator set each digest to the SHA-256 of its excerpt. This correction is evidence-only: it changes no finding, observation, source identity, or verification outcome. The second checkpoint run was green.

## TT-C3 author evidence

Ticket 3 had a fresh `bench-writer` successor author, `tt-t3-author-2`, on opus at high effort, with a cap of 3 attempts. The first author, `tt-t3-author`, stopped blocked on a fence gap and a behavioral question, and it committed nothing. The reviewer decided both questions, and the coordinator committed the plan change `af654c55`.

The successor started at `af654c55` with the uncommitted diff of the first author. The successor verified that diff again with its own runs. It kept the production behavior, and it changed four test files and one comment. It committed `93b7c071` on a lane pass in the first attempt. The successor then committed this record in a second commit.

The chunk pair is `74daf800..93b7c071`. The base is the accepted TT-C2 tip, because a plan commit is never a chunk base. The payload names the ticket commit as the tip, because a record cannot name its own commit. The source digest does not change at the record commit, because the record file is outside the graded source.

The successor made these changes to the preserved diff:

- `cmd/bench/response_bound_exempt_test.go`: the planted `dashboard` and ship-tier calls print the exempt row on stderr, and the bounded `dashboard` call prints it on stdout. The new `withoutRows` helper removes the leading row of each stream before `requireComplete` and `requireBounded` grade the verb output.
- `internal/systemtest/adoption_test.go`: each nested Bench call of the scaffolded gate prints its own row. So the gate with an empty `tests/canary` prints 11 lines and spills, and the canary message moves into the spill file. The new `boundedStderr` helper reads the spill file of a spilled response. `assertPrivateHomeEmpty` accepts the `responses` entry, which holds that spill file.
- `internal/treetarget/identify_test.go`: the path check of TT23 runs before the label check of TT12, so each row has its own red.
- `cmd/bench/tree_scope_test.go`: the row fixture fixes its commit dates, so HEAD is always the commit `2151c564`. The deviations below give the reason.
- `cmd/bench/census_output.go`: one comment gets its missing article.

### Red and green log

The production code was already in the tree at the start. So each red came from a `bench probe`, from a run before a test edit, or from a copy-aside edit of a system test. After each probe, the restore read `yes`, and the test passed again.

- TT10, TT25, and TT60: a probe wrote the row into the bounded stream after the verb. `TestTreeRowLeadsTreeResponse` lost the leading row. `TestTreeRowSurvivesSpill` printed `line 01` first, and `TestTreeRowOutsideResponseBound` spilled at 12 lines.
- TT11: a probe set the primary target to `main`. `TestTreeRowLeadsTreeResponse` failed on the row `main,2151c564…,false`.
- TT15: a probe cut the head cell to 12 characters. `TestTreeRowLeadsTreeResponse` failed on the row `primary,2151c5641124,false`.
- TT12: a probe set the target to the assignment ID. `TestIdentifyNamesActiveLabel` failed on the label check.
- TT23: a probe set the target to the worktree path. `TestIdentifyNamesActiveLabel` failed on the path check.
- TT13 and TT14: the plan probe `3-unassigned-probe` failed both subtests of `TestIdentifyUnownedWorktree` with the target `primary`.
- TT14 alone: a probe in `internal/intent/assignment.go` accepted an assignment in any state. The released subtest failed with the target `alpha`.
- TT16: a probe set the head fallback to an empty string. `TestIdentifyUnbornHead` failed on a quoted empty cell.
- TT17: a probe added `--untracked-files=no` to the status query in `internal/git/status.go`. The untracked subtest of `TestIdentifyDirtyStates` failed with `false`.
- TT18: a probe replaced the status query with a query of untracked files only. The modified-file subtest failed with `false`.
- TT19: a probe set the dirty fallback to `false`. The corrupt-index subtest failed with `false`.
- TT20: a probe printed the exempt row on stdout. `TestExemptTreeCallPrintsRowOnStderr` failed.
- TT21: a probe removed the help-form check from the row hook. `bench gate --help` printed the row on stderr.
- TT22: a probe removed the scope check from the row hook. `bench version` printed the row.
- TT24: a probe skipped the empty-root check. `bench coverage x` printed the row `unassigned,none,unknown` before its answer.
- TT56: a probe changed the exit-2 rule to exit 99. `TestRunGateRejectsBriefUsage` failed, because `bench gate --brief` printed the row on stdout.
- TT61: a probe removed the passing-phase guard in `internal/responsebound/owner.go`. `TestOwnerPrintsNoRowWhenSpillCannotOpen` failed, because the row came after the verb output.
- The posture sites: before the test edits, `TestDashboardStdoutStaysComplete` and `TestShipTierStaysComplete` failed at lines 93 and 145, and `TestAdoptionSmokeJourney` failed at line 151. With the `responses` allowance removed, the journey failed with `private BENCH_HOME … is not empty: [d otel/ d responses/]`. After each edit, the tests passed.
- The header expectation in `cmd/bench/tree_scope_test.go`: a probe renamed the field `head` to `commit`. `TestTreeRowLeadsTreeResponse` and `TestTreeRowSurvivesSpill` failed.

### Deviations

- Non-behavioral spec contradiction, for reviewer veto: TOON quotes a string that starts with a zero and a digit. So for a HEAD such as `025072b2…`, the row is `primary,"025072b2…",false`, and not the bare form of the acceptance row. About 4 in 100 commits print quoted. The successor follows the tree convention of TOON quoting, as `internal/diff/review_base_test.go` does. A final run on a random fixture commit gave this red, so the row tests now pin one commit that prints bare.
- `Row` renders through `toon.TableTyped`, not `toon.Table`. The dirty cell must print a bare `true` or `false`, and `toon.Table` quotes the string `true`. Both functions share one block contract.
- The row-removal logic lives in the test package `internal/treetarget/treetargettest`, and the `cmd/bench` helper `withoutTreeRow` calls it. The system tests cannot import package `main`, so they call the same function.
- The adoption journey accepts a `responses` entry in its private home after every leg. By the nested-call decision, any gate response can spill, and a spill writes only that entry.
- The ticket marks `cmd/bench/tree_scope.go` as new, but ticket 1 created it. No check required an edit to `cmd/bench/help_inventory_test.go` or to the two conformance files in `Writes:`, so the diff leaves them unchanged.

The sweep of duplicated facts found no second source. `withoutRows` and `boundedStderr` compose the existing `withoutTreeRow`, `responseboundtest.Find`, and `readSpillFile`. The row header in `treetargettest` comes from the renderer. The independent expectations are the header, the cell values, and the `responses` entry name. The log above records a red for each.

The source digest is the tree of `93b7c071` without this record file, which is `8bf44343`. The same rule at `74daf800` gives `80fbc1c6`, which confirms the method.

The plan commits after `74daf800` fold the reviewer decisions and add the rows TT60, TT61, and TT62. They also change the fences of tickets 3 and 4 and the ticket 3 assignments. So the plan digest changed from `834182d7` to `ea3ee9d0`. The `ReadPlan` rule at `74daf800` gives `834182d7`, which confirms the method. The payload adds one amendment that maps each chunk ID to itself, because no chunk ID changed.

### Probe verdicts

Each probe ran through `bench probe`, and each restore reads `yes`. The first row is the plan probe `3-unassigned-probe`. The JSON payload holds the exact command and output of the first five rows at `93b7c071`. The other rows ran on the source before the fixture commit dates were fixed, so only this table holds them.

| File | Mutation | Test | Row | Verdict |
|---|---|---|---|---|
| `internal/treetarget/identify.go` | swap: `"unassigned"` to `"primary"` | TestIdentifyUnownedWorktree | TT13, TT14 | bit |
| `cmd/bench/tree_scope.go` | swap: `if exit == 2 {` to `if exit == 99 {` | TestRunGateRejectsBriefUsage | TT56 | bit |
| `cmd/bench/census_output.go` | swap: `owner.Lead(shownRow(row, exit))` to a write into the bounded stream | the three row placement tests | TT10, TT25, TT60 | bit |
| `internal/treetarget/identify.go` | swap: `identity.Target = "primary"` to `identity.Target = "main"` | TestTreeRowLeadsTreeResponse | TT11 | bit |
| `internal/treetarget/identify.go` | swap: `identity.Head = head` to `identity.Head = head[:12]` | TestTreeRowLeadsTreeResponse | TT15 | bit |
| `internal/treetarget/identify.go` | swap: `assignment.Label` to `assignment.ID` | TestIdentifyNamesActiveLabel | TT12 | bit |
| `internal/treetarget/identify.go` | swap: `assignment.Label` to `assignment.Worktree` | TestIdentifyNamesActiveLabel | TT23 | bit |
| `internal/intent/assignment.go` | swap: `a.State == StateActive` to `a.State != ""` | TestIdentifyUnownedWorktree | TT14 | bit |
| `internal/treetarget/identify.go` | swap: `Head: "none"` to `Head: ""` | TestIdentifyUnbornHead | TT16 | bit |
| `internal/git/status.go` | swap: add `--untracked-files=no` to the status query | TestIdentifyDirtyStates | TT17 | bit |
| `internal/git/status.go` | swap: the status query to `ls-files --others --exclude-standard` | TestIdentifyDirtyStates | TT18 | bit |
| `internal/treetarget/identify.go` | swap: `Dirty: "unknown"` to `Dirty: false` | TestIdentifyDirtyStates | TT19 | bit |
| `cmd/bench/tree_scope.go` | swap: `c.Stderr` to `c.Stdout` in `finishExempt` | TestExemptTreeCallPrintsRowOnStderr | TT20 | bit |
| `cmd/bench/tree_scope.go` | swap: remove the help-form clause | TestHelpFormPrintsNoTreeRow | TT21 | bit |
| `cmd/bench/tree_scope.go` | swap: remove the scope clause | TestHelpFormPrintsNoTreeRow | TT22 | bit |
| `cmd/bench/tree_scope.go` | swap: `if root == "" {` to `if root == "-" {` | TestTreeRowOutsideRepository | TT24 | bit |
| `internal/responsebound/owner.go` | swap: remove the `passing` clause of the lead guard | TestOwnerPrintsNoRowWhenSpillCannotOpen | TT61 | bit |
| `internal/treetarget/identify.go` | swap: the field `"head"` to `"commit"` | TestTreeRowLeadsTreeResponse, TestTreeRowSurvivesSpill | header | bit |

A first placement probe used `fmt.Fprint`, which `census_output.go` does not import. Its verdict was `invalid`, and the table holds the rerun. One more probe added a no-op statement to the exit-2 rule to confirm that the new fixture compiles. Its verdict was `silent`, and it grades no row.

### Verification

The successor ran each TT-C3 plan verification on the source of `93b7c071`, and each passed. The JSON payload holds each result. The successor also ran these checks on the same source, and each passed:

- `bench test --package ./internal/conformance`, with three environment capability skips for unix sockets or device nodes;
- `bench test --check system`;
- `bench structure --growth f981cd3d`, which reported that no source file grew past its budget;
- `gofmt -l` and `go vet -tags system` on the changed packages.

The first post-commit preflight reported `binary-seal` red, because the worktree build was older than the source. After `bench worktree build`, `bench preflight build tree-targets` reported 14 green checks, 1 check that does not apply, and 0 red checks.

## TT-C3 chunk review, round 1

The frozen pair is base `74daf800f8a11d85083711b31aaabc47797ebc73` and tip `5444f36ac57e3310026b62439bcf7648d49871df`. The shared evidence is `sha256:6aae40c43e205e132c18b04f6744eb48ca19059cae6432db32034816d803e5e7`. Each axis ran in a fresh `bench-reviewer` session on opus at high effort. Only the Coverage axis ran probes, and it left the tree clean.

The raw finding count is 3: Standards 0, Spec 1, and Coverage 2. Each finding names its own fix, so 3 repair targets remain. The coordinator also corrects the spec text on the table renderer, which is a plan change and not a repair target.

## Standards

Findings: 0. Each independent header expectation has a recorded red.

## Spec

Findings: 1. The worst issue is that a label with a control byte drops the row.

- `cmd/bench/tree_scope.go:73-76` returns no row when `treetarget.Row` fails. The ledger accepts a label with a control byte, so such an assignment gives a tree-scoped response with no row. The spec says the row is the first block of each tree-scoped response that resolves a root and exits other than 2. The spec's import list names `internal/sanitize`, so the tree convention renders such a label through it. Target R1, `auto-fix`, confidence 6. The coordinator routes this as a spec predicate and flags the rendering choice for reviewer veto.

## Coverage

Findings: 2. The worst issue is two untested timing rules of the row.

- The exempt path's exit-2 rule has no test. A probe that printed the row on stderr at every exit stayed silent across `./cmd/bench`. `bench release-preflight` with no arguments is exempt and tree-scoped, and it exits 2. Target R2, `auto-fix`, confidence 7.
- No test pins that the dispatcher computes the row before the verb runs. A probe that computed the row after the verb stayed silent. A verb that moves HEAD, such as `bench commit`, would then name the new commit. Target R3, `auto-fix`, confidence 7.

## TT-C3 repair routing

Each repair goes to one fresh `bench-writer` repair session for ticket 3 on opus at high effort. Ticket 3's `Writes:` line holds every path. This is cycle 1 of the two repair cycles for chunk TT-C3.

| Target | Ticket | Repair |
|---|---|---|
| R1 | 3 | Render a label that TOON cannot carry through `sanitize`, so the row always prints, and pin it. |
| R2 | 3 | Test an exempt tree-scoped call that exits 2, with no row on either stream. |
| R3 | 3 | Test that the row names the tree state from before the verb, with a planted verb that moves HEAD. |

## TT-C3 ticket 3 repair evidence, cycle 1

The session `claude:bench-writer/tt-t3-repair-1` ran on opus at high effort, with a cap of 3 attempts. It started at `fe44447f` and committed `744ef656` on a lane pass in the first attempt. This repair is cycle 1 of the two repair cycles for chunk TT-C3.

- R1: `treetarget.Row` escapes a target label that fails `sanitize.LineSafe` through `sanitize.Controls`, so the row always prints. The predicate is the one that `internal/worktree/list_selected.go` uses for a hostile operand. A pointer such as `target-1` names a request ordinal, and a row has no request, so the label prints escaped. `sanitize.Strip` could make a hostile label read as `primary`, so the session did not use it. `TestRowEscapesHostileLabel` records a ledger label with a BEL byte and asserts the escaped cell and no raw control byte. The coordinator flagged the rendering choice for reviewer veto.
- R2: `TestExemptTreeRowFollowsExitRule` runs one planted exempt tree-scoped verb at exit 1 and at exit 2. At exit 1, stderr holds the row, so the checkout prints a row. At exit 2, neither stream holds a `tree[` line. The review probe for R2 was silent in round 1. It now bites.
- R3: `TestTreeRowPrecedesVerb` plants a tree-scoped verb that writes an untracked file. On the bounded path and on the exempt path, the row still reads `dirty` `false`. The review probe for R3 was silent in round 1. It now bites, and a second probe that computes the exempt row after the verb also bites.

`plantedTreeVerb` now takes a bound disposition and a handler, so four row tests share one planted registry. The new `exitingHandler` repeats the three-line closure in `runBoundFixture` of `cmd/bench/response_bound_test.go`. That file is outside the ticket 3 fence, so the session did not fold the two. The sweep of duplicated facts found no other second source. The independent expectations are the escaped cell, the exit 1 row, no row at exit 2, and the clean row after a dirtying verb. The probes below record a red for each.

### Probe verdicts

Each probe ran through `bench probe` on the source of `744ef656`, and each restore reads `yes`. The JSON payload holds the exact command and output of each probe. The first row is the plan probe `3-unassigned-probe`, which the session ran again with the exact plan command.

| Target | File | Mutation | Test | Verdict |
|---|---|---|---|---|
| TT13, TT14 | `internal/treetarget/identify.go` | swap: `"unassigned"` to `"primary"` | TestIdentifyUnownedWorktree | bit |
| R1 | `internal/treetarget/identify.go` | swap: the `LineSafe` condition to `false && …` | TestRowEscapesHostileLabel | bit |
| R2 | `cmd/bench/tree_scope.go` | swap: `shownRow(row, exit)` to `row` in `finishExempt` | TestExemptTreeRowFollowsExitRule | bit |
| R2 | `cmd/bench/tree_scope.go` | swap: `if exit == 2 {` to `if exit != 0 {` | TestExemptTreeRowFollowsExitRule | bit |
| R3 | `cmd/bench/census_output.go` | swap: `row` to `definition.treeRow(args)` in the owner lead | TestTreeRowPrecedesVerb | bit |
| R3 | `cmd/bench/command_registry.go` | swap: compute the exempt row after the verb | TestTreeRowPrecedesVerb | bit |

The two review probes for R2 and R3 ran on all of `./cmd/bench`, and each failed only the new test.

### Verification

The session ran each TT-C3 plan verification on the source of `744ef656`, and each passed. The JSON payload holds each result. The session also ran these checks on the same source, and each passed:

- `bench test --check system`;
- `bench structure --growth f981cd3d`, which reported that no source file grew past its budget;
- `gofmt -l` on the changed packages.

The first post-commit preflight reported `binary-seal` red, because the worktree build was older than the source. After `bench worktree build`, `bench preflight build tree-targets` reported 14 green checks, 1 check that does not apply, and 0 red checks.

The chunk tip is now the repair commit `744ef656`. The source digest is the tree of `744ef656` without this record file, which is `564566b2`. The same rule at `5444f36a` gives the round 1 digest `8bf44343`, which confirms the method.

The spec changed at `fe44447f`, so the plan digest changed from `ea3ee9d0` to `51bdf6fe`. The `ReadPlan` rule at `156beff0` gives `ea3ee9d0`, which confirms the method. The payload adds one amendment that maps each chunk ID to itself, because no chunk ID changed. The TT-C1 and TT-C2 chunks keep the digests of their own tips. The round 1 entries keep their earlier source digest as history.

## TT-C3 chunk review, round 2

The frozen pair is base `74daf800f8a11d85083711b31aaabc47797ebc73` and tip `7251523e83a6993521cbd51db6531daf2297d9b9`. The shared evidence is `sha256:37b6823245848942030c2bd1d8a39475060da0aa2ba9d57f0cf168f53b2f8ac1`. This round is the confirming round of all three axes, and each axis read only the repair delta `fe44447f..7251523e`.

## Standards

Findings: 2. The worst issue is that the label escape is conditional.

- `internal/treetarget/identify.go:49-50` calls `sanitize.Controls` only when `sanitize.LineSafe` fails. So a hostile label and a literal label with a backslash escape can print the same cell. The `LineSafe` contract says a caller that fails the predicate emits a pointer, not an escaped value. The `sanitize` package names `Strip` as its duty for a table cell that must read as its source text. Target R4, `auto-fix`, confidence 7. The coordinator routes this to `sanitize.Strip` and flags the choice for reviewer veto.
- The repair widened `plantedTreeVerb` in `cmd/bench/tree_scope_test.go:28-40` into a second copy of `runBoundFixture` in `cmd/bench/response_bound_test.go:67-79`. `AGENTS.md` names a fixture harness pasted N times as duplicated knowledge. Target R5, `auto-fix`, confidence 6. The fix needs `cmd/bench/response_bound_test.go` in the ticket 3 fence, which the coordinator adds in a plan commit first.

## Spec

Findings: 0. R1, R2, and R3 are confirmed, and the rows stay met.

## Coverage

Findings: 0. R1, R2, and R3 are pinned, and an independent probe on the shared exit-2 guard bit.

## TT-C3 repair routing, cycle 2

Each repair goes to one fresh `bench-writer` repair session for ticket 3 on opus at high effort. This is cycle 2 of the two repair cycles for chunk TT-C3.

| Target | Ticket | Repair |
|---|---|---|
| R4 | 3 | Render every label through `sanitize.Strip`, and pin a hostile label with a test. |
| R5 | 3 | Give `plantedTreeVerb` and `runBoundFixture` one harness with a scope parameter. |

## TT-C3 ticket 3 repair evidence, cycle 2

The session `claude:bench-writer/tt-t3-repair-2` ran on opus at high effort, with a cap of 3 attempts. It started at `92e86104` and committed `9b7df67e` on a lane pass in the first attempt. This repair is cycle 2 of the two repair cycles for chunk TT-C3.

- R4: `treetarget.Row` renders every label through `sanitize.Strip`, as the amended identity-row section states. The `LineSafe` and `Controls` branch is gone, and a clean label passes unchanged. `TestRowStripsHostileLabel` replaces `TestRowEscapesHostileLabel`. It pins two cells: a label with a BEL byte prints `alpha`, and a literal backslash label prints its own text, escaped once by TOON. Before the fix, the hostile row was red and printed the same cell as the literal label. After the fix, both rows are green.
- R5: `runPlanted` in `cmd/bench/response_bound_test.go` is now the one planted-command harness, with a bound, a scope, and a handler. The four bound tests call it with the zero scope, and `plantedTreeVerb` calls it with `scopeTree`. `exitingHandler` moved to the same file and takes a line count, so `runBoundFixture` and its copied closure are gone. Every current test passes.

The sweep of duplicated facts found no other second source in the delta. The `treeRow` comment in `cmd/bench/tree_scope.go` now names the strip instead of the escape. The independent expectations are the stripped cell and the literal cell, and the R4 probe below records a red for each. The R5 probe changes one line of the shared harness and reds tests in both files, which shows that one harness serves both callers.

### Probe verdicts

Each probe ran through `bench probe` on the source of `9b7df67e`, and each restore reads `yes`. The JSON payload holds the exact command and output of each probe. The first row is the plan probe `3-unassigned-probe`, which the session ran again with the exact plan command.

| Target | File | Mutation | Test | Verdict |
|---|---|---|---|---|
| TT13, TT14 | `internal/treetarget/identify.go` | swap: `"unassigned"` to `"primary"` | TestIdentifyUnownedWorktree | bit |
| R4 | `internal/treetarget/identify.go` | swap: `sanitize.Strip` to `sanitize.Controls` on the target | TestRowStripsHostileLabel | bit, both subtests |
| R5 | `cmd/bench/response_bound_test.go` | swap: `bound` to `boundResponse` in `runPlanted` | three tests in two files | bit |

### Verification

The session ran each TT-C3 plan verification on the source of `9b7df67e`, and each passed. The JSON payload holds each result. The session also ran these checks on the same source, and each passed:

- `bench test --check system`;
- `bench test --check single-control-escaper`;
- `bench structure --growth f981cd3d`, which reported that no source file grew past its budget;
- `gofmt -l` on the changed packages.

The commit ran with `--preflight-build tree-targets`. The worktree build was green, and `bench preflight build tree-targets` reported 15 green checks, 1 check that does not apply, and 0 red checks.

The chunk tip is now the repair commit `9b7df67e`. The source digest is the tree of `9b7df67e` without this record file, which is `ef1d5964`. The same rule at `744ef656` gives the cycle 1 digest `564566b2`, which confirms the method.

The spec changed at `92e86104`, so the plan digest changed from `51bdf6fe` to `49d7c345`. The `ReadPlan` rule at `744ef656` gives `51bdf6fe`, which confirms the method. The payload adds one amendment that maps each chunk ID to itself, because no chunk ID changed. The earlier entries keep their earlier source digests as history.

## TT-C3 chunk review, round 3, and close

The frozen pair is base `74daf800f8a11d85083711b31aaabc47797ebc73` and tip `9e742c5434c327e7174f85dd427331004b2bcef6`. The shared evidence is `sha256:a185d294f8016f0b67c880a7f5b8c78859c637767c3d68bf96a6cc1173cb0664`. This round is the confirming round of all three axes after repair cycle 2, and each axis read only the repair delta `92e86104..9e742c54`.

## Standards

Findings: 0. R4 and R5 are confirmed, and one planted-command harness remains.

## Spec

Findings: 0. R4 and R5 are confirmed, and the rows of ticket 3 stay met.

## Coverage

Findings: 0. R4 and R5 are pinned, and an independent scope probe bit.

## Advice

- `TestRowStripsHostileLabel` tests only the BEL byte. A table over more than one refused byte would pin the whole strip rule.
- `sanitize.Strip` passes DEL and the C1 controls through, because TOON accepts them. A stripped label can also equal another label. The rendering choice stays flagged for reviewer veto.
- On the bounded path, no test covers a tree call that exits with a code other than 0 or 2. The exempt path has an exit-1 control.

Chunk TT-C3 closes after two repair cycles, the full allowance.

## TT-C4 ticket 4 author evidence

Ticket 4 had a fresh `bench-writer` successor author, `tt-t4-author-3`, on opus at xhigh effort, with a cap of 4 attempts. The first author, `tt-t4-author`, stopped blocked on a fence gap and committed nothing. The second author, `tt-t4-author-2`, lost its session and left an uncommitted partial diff. The coordinator kept that diff outside the tree.

The successor started at `ada9184e` with a clean tree. It committed `d5bad995` on a lane pass in the first attempt. The successor then committed this record in a second commit.

The successor read the preserved diff as prior art and verified each kept part again with its own runs. It kept these parts:

- `internal/treetarget/flag.go`, `internal/worktree/tree_target.go`, and the constant `primaryTarget`, with no change.
- The split of `runChild` out of `runWorktreeChild` in `internal/worktree/exec.go`, and the constant `WrapperEnv`.
- The change to `canonicalpath.Resolve` and its test `TestResolveRelativeUnderSymlinkedWorkingDirectory`.
- The help insertion, the test `TestHelpRendersTreeTargetFromScope`, and the help-row edits in four `cmd/bench` test files.
- `internal/treetarget/run.go`, its tests, the lookup tests in `internal/worktree`, `TestTreeTargetOnlyAsFirstArgument`, and the system test, with the changes below.

The successor changed these parts of the preserved diff:

- `treetarget.Call` takes the invoking wrapper as a field. So the `internal/treetarget` tests bind no process environment, and the dispatcher reads `BENCH_WRAPPER` through `worktree.WrapperEnv`.
- `TestRunStartsChildInTarget` does not bind `BENCH_KIT` or `BENCH_RUN_BINARY`, because TT48 in ticket 5 owns the child environment.
- The system test reads the row cells through `treetargettest.WithoutRow` and `systemTOONCell`. The literal row of the preserved test failed two system runs, because TOON quotes a HEAD that starts with a zero and a digit.
- `TestTreeTargetPrefersOneActiveLabel` creates the retired assignment first, so a lookup that takes the first match fails.
- `TestTreeTargetOnlyAsFirstArgument` binds no Bench home, because `TestMain` already gives each test a private home.

### Red and green log

The production code was in the tree before the first test run. So each red came from a `bench probe`, from a revert of one file to HEAD, or from a copy-aside edit for the system suite. After each probe and each restore, the test passed again.

- TT8: a probe omitted the help insertion. `TestHelpInventoryIsComplete` failed, and so did the pinned help rows in `TestHelpKeepsStatusPublicRoute`, `TestRootAndHelpAlignWrapperAndBinary`, and four other help tests.
- TT9: the same omission failed the tree row of `TestHelpRendersTreeTargetFromScope`. A probe that inserted the flag on every row failed its repository row.
- TT26, TT28, and TT41: a copy-aside edit made the dispatcher ignore `--in`. The system suite failed the three rows with `usage: bench status (unknown argument: --in)` and the coverage equivalent.
- TT54: in the same system run, a copy-aside edit made `bench worktree exec` refuse a child named `bench`. The row failed at exit 2, and TT55 stayed green.
- TT55: a copy-aside edit made `boundaryRoot` answer the primary checkout. The row failed with the target `primary`.
- TT27: a probe kept `--in alpha` in the child argv, and a probe ran the child in the parent root. Each failed `TestRunStartsChildInTarget`. A probe that printed a `worktree:` line in the parent failed the check for no parent output.
- TT53: a probe gave the child the directory spelling `<root>/.`, and the PWD check failed. A probe that dropped `canonicalpath.Resolve` stayed silent, as the spec states, because the create verb records the physical path.
- TT29 and TT57: the plan probe `4-flag-probe` failed both rows. A probe that removed the empty-value test failed TT57 alone with the unassigned refusal at exit 1.
- TT58: a probe that removed the dash step failed TT58 and the dash row with a control character. A probe that printed the raw value failed that dash row alone.
- TT30, TT31, and TT59: three probes of `pathShaped` each dropped one path shape. They failed TT30 with TT59, TT31 with the `~nobody` row, and TT59 alone.
- TT32: a probe ran the path-shape test before the label lookup, and `team/alpha` gave the usage line at exit 2. A probe that answered the primary root failed the child directory check.
- TT33: a probe answered the primary root for an unknown label, and a child started at exit 0.
- TT34: a probe accepted each ledger state, and the released label started a child.
- TT35: a probe took the first of two colliding matches, and a child started.
- TT36: a probe removed the control-character test, and the value gave the unassigned refusal.
- The marker checks of `TestRunRefusesBeforeChild`: a probe started a child before the lookup, and the 8 lookup rows failed on the marker. A probe started a child before the grammar steps, and all 12 rows failed on the marker.
- TT37: a probe moved a late `--in` to the front, and `bench gate --fresh --in primary` started the child at exit 0. A probe refused a late `--in` on a repository verb, and `bench version x --in primary` exited 2.
- TT37 controls: a probe that ignored `--in` failed the first-position call. A probe that consumed `--in` and started no child failed the marker control. A probe that started a silent child for a late flag failed the late-flag marker check.
- TT38: three probes in `internal/worktree/exec.go` failed `TestRunPassesChildResult`. The probes added the `worktree:` line, capped the exit at 1, and dropped the child stdout.
- TT39: a probe bypassed the wrapper, and the wrapper marker was absent. A probe that also started the running executable failed the check for its absent marker.
- TT40: a probe demanded a wrapper, and the direct run failed with `exec: no command`.
- TT49: a probe ran `<root>/dist/bench`, and the stray marker existed in both subtests. A probe that kept `--in alpha` in the argv failed both argv checks.
- TT62: a revert of `internal/canonicalpath/canonicalpath.go` to HEAD failed the `.` row with the link spelling. A probe that joined the working directory through `filepath.Join` failed the `jump/..` row.
- The lookup tests in `internal/worktree`: a probe that ignored the active match gave the ambiguity refusal. A probe that took the first match gave the state refusal instead of the ambiguity.
- A probe that ran the general resolver first failed the id, prefix, label-prefix, and absolute-path rows of `TestTreeTargetTakesOnlyAnExactLabel`. A probe that dropped the path-shape outcome failed its two path rows.

### Deviations

- Non-behavioral reading, for reviewer veto: TT34 asks that stderr names the released state. The state refusal of step 5 prints `assignment <id> is not active`, and the shared printer drops the observed state. `internal/worktree/path.go` owns that printer and is outside the fence.
- `internal/treetarget/identify.go` gains the constant `primaryTarget`, which the row cell and the keyword share. The ticket marks `internal/treetarget/` as new, but ticket 3 created it, and the fence holds it as a prefix.
- Outside a repository, a call with `--in` prints the not-in-repo line on stderr at exit 1 after the grammar steps. The spec states no rule for this case.
- The usage line of a dash value or a path value escapes the value through `sanitize.Controls`, so a backslash prints doubled. Only a dash value can carry a control character, because step 3 runs before step 6.
- `TestTreeTargetOnlyAsFirstArgument` binds `BENCH_WRAPPER` to a marker script in `cmd/bench`. Without it, a late-flag mutation would start the test binary itself as the child. No `internal/worktree` test binds the environment, so the serial census does not change.
- `internal/systemtest/tree_target_test.go` imports `internal/worktree` for `WrapperEnv`.

The sweep of duplicated facts found one repeated fragment. `TreeTarget` builds the ambiguity id list with the same four lines as `selectAssignment` in `internal/worktree/path.go`, which is outside the fence. `Operand`, `primaryTarget`, `treeTargetFlag`, `leadsWithTreeTarget`, `WrapperEnv`, and `runChild` each hold one fact that two callers share. The independent expectations are the usage lines, the refusal lines, the child argv, directory, and PWD, the help rows, and the row cells. The log above records a red for each.

### Probe verdicts

Each probe ran through `bench probe`, and each restore reads `yes`. The probes ran on the source that `d5bad995` commits, with one exception. Before the commit, the successor removed a redundant raw-byte clause from `TestRunRefusesBeforeChild`. So the plan probe and the before-lookup probe ran again at `d5bad995`, and both bit. The first row is the plan probe `4-flag-probe`, with the exact plan command. The second row is the successor's own probe of the central property.

| File | Mutation | Test | Row | Verdict |
|---|---|---|---|---|
| `internal/treetarget/flag.go` | swap: `<label\|primary>` to `<label>` | TestRunRefusesBeforeChild | TT29, TT57 | bit |
| `cmd/bench/tree_scope.go` | self-probe, swap: the target guard to `if true \|\| …` | TestTreeTargetOnlyAsFirstArgument | TT37 | bit |
| `internal/treetarget/run.go` | swap: start a child before the lookup | TestRunRefusesBeforeChild | 8 lookup rows | bit |
| `internal/treetarget/run.go` | swap: start a child before the grammar steps | TestRunRefusesBeforeChild | 12 rows | bit |
| `internal/treetarget/run.go` | swap: keep `--in alpha` in the child argv | TestRunStartsChildInTarget | TT27 | bit |
| `internal/treetarget/run.go` | swap: resolve the parent root | TestRunStartsChildInTarget | TT27 | bit |
| `internal/treetarget/run.go` | swap: the child directory to `dir+"/."` | TestRunStartsChildInTarget | TT53 | bit |
| `internal/treetarget/run.go` | swap: drop `canonicalpath.Resolve` | TestRunStartsChildInTarget | TT53 | silent |
| `internal/treetarget/flag.go` | swap: remove the empty-value test | TestRunRefusesBeforeChild | TT57 | bit |
| `internal/treetarget/flag.go` | swap: remove the dash step | TestRunRefusesBeforeChild | TT58 | bit |
| `internal/treetarget/flag.go` | swap: print the raw value | TestRunRefusesBeforeChild | dash row | bit |
| `internal/worktree/tree_target.go` | swap: three `pathShaped` omissions | TestRunRefusesBeforeChild | TT30, TT31, TT59 | bit |
| `internal/worktree/tree_target.go` | swap: path shape before the label lookup | TestRunLabelWinsOverPathShape | TT32 | bit |
| `internal/worktree/tree_target.go` | swap: answer the primary root for no match | TestRunRefusesBeforeChild | TT33 | bit |
| `internal/worktree/tree_target.go` | swap: accept each ledger state | TestRunRefusesBeforeChild | TT34 | bit |
| `internal/worktree/tree_target.go` | swap: `> 1` to `> 2` for the ambiguity | TestRunRefusesBeforeChild | TT35 | bit |
| `internal/worktree/tree_target.go` | swap: remove the control-character test | TestRunRefusesBeforeChild | TT36 | bit |
| `internal/worktree/exec.go` | swap: three child-result faults | TestRunPassesChildResult | TT38 | bit |
| `internal/treetarget/run.go` | swap: bypass the wrapper | TestRunSelectsChildExecutable | TT39 | bit |
| `internal/treetarget/run.go` | swap: demand a wrapper | TestRunSelectsChildExecutable | TT40 | bit |
| `internal/treetarget/run.go` | swap: run `<root>/dist/bench` | TestRunSelectsChildExecutable | TT49 | bit |
| `cmd/bench/command_registry.go` | omit: the help insertion | eight help tests | TT8, TT9 | bit |
| `cmd/bench/tree_scope.go` | swap: insert the flag on every row | TestHelpRendersTreeTargetFromScope | TT9 | bit |
| `cmd/bench/tree_scope.go` | swap: move a late `--in` to the front | TestTreeTargetOnlyAsFirstArgument | TT37 | bit |
| `internal/canonicalpath/canonicalpath.go` | swap: join through `filepath.Join` | TestResolveRelativeUnderSymlinkedWorkingDirectory | TT62 | bit |
| `cmd/bench/main.go` | swap: remove the `anchors` scope | the `subcommand-routing` check | liveness | bit |

A first try of three probes did not compile, and each verdict was `invalid`. The table holds the rerun of each. The system rows TT26, TT28, TT41, TT54, and TT55 took copy-aside edits, because `bench probe` does not take the system suite. The log above names each edit. Those runs came before `runTreeTarget` took `worktree.WrapperEnv` in place of the same literal name.

### Verification

The successor ran each check below on the source of `d5bad995`, and each passed:

- `bench test --package ./internal/treetarget`, in 456 ms;
- `bench test --package ./internal/worktree`, in 52022 ms, with two environment capability skips for unix sockets;
- `bench test --package ./internal/canonicalpath`, in 4 ms, where the revert run failed in 4 ms, so the tests run;
- `bench test --package ./cmd/bench`, in 12893 ms;
- `bench test --check system`, in 52478 ms;
- `bench test --check subcommand-routing`, in 20 ms, where the liveness probe bit in 20 ms;
- `bench test --package ./internal/conformance --run TestRootConformance`, in 6658 ms;
- `bench test --check skip-ownership`, in 41 ms;
- `bench test --package ./internal/conformance`, in 38202 ms, with three environment capability skips;
- `bench structure --growth f981cd3d`, which reported that no source file grew past its budget;
- `gofmt -l` and `go vet -tags system` on the changed packages.

The commit ran with `--preflight-build tree-targets`. The lane passed, the worktree build was green, and `bench preflight build tree-targets` reported 15 green checks, 1 check that does not apply, and 0 red checks.

## TT-C4 ticket 5 author evidence

Ticket 5 had a fresh `bench-writer` author, `tt-t5-author`, on opus at xhigh effort, with a cap of 4 attempts. The author started at `79e79c51` with a clean tree. It committed `c587623b` on a lane pass in the third attempt. The first two attempts made the package rows green, but the system row failed on the observation line that the fixture environment turns on. The author then committed this record in a second commit.

The ticket changed these parts:

- `internal/treetarget/build.go` is new. It holds the executable choice of ticket 4 and extends that choice to a kit worktree target. It also holds the two refusal texts and the refusal printer.
- `internal/treetarget/run.go` calls the new choice after it resolves the child directory. A refusal prints before any child starts.
- `internal/treetarget/kittest` is a new test support package. It writes the smallest tree that declares Bench build inputs, and it edits the one listed build input. The internal tests of `internal/treetarget` and the system suite both import it. The package `treetargettest` cannot hold it, because `treetargettest` imports `internal/treetarget`, so an internal test of `internal/treetarget` cannot import `treetargettest`.
- `internal/treetarget/build_test.go` adds `TestRunKitWorktreeBuild`. `internal/systemtest/tree_target_test.go` adds `TestTreeTargetRefusesStaleKitBuild`.
- `runTreeTarget` in the system suite removes `BENCH_COMMAND_OBSERVE` from the fixture environment, so stderr holds only the lines of the verb. The five rows of ticket 4 read only stdout, so their result does not change.

### Red and green log

The tests were in the tree before the first production edit. The author ran them against the source of `79e79c51`, and each row failed:

- TT43: the wrapper started, and the check for an absent wrapper marker failed.
- TT44, TT45, TT46, and TT47: each call exited 0 through the wrapper, where each row wants exit 1 and the two refusal lines.
- TT50: the system suite ran `bench status --in alpha` at exit 0, and stdout started with the row of `alpha`.

After the production edit, each row passed. Each red below came from a `bench probe` or from a copy-aside edit, and each row passed again after the restore:

- TT43: a probe that ignored the build inputs started the wrapper. The current-build row failed on the wrapper marker, and the four refusal rows exited 0.
- TT44: the plan probe failed TT44 and TT47. A probe that started the wrapper after a refusal failed the four refusal rows at exit 0. A probe that always quoted the label failed TT44, TT45, and TT46 with `'alpha'`.
- TT45: a probe that checked only the executable digest through `freshness.VerifyExecutable` started the stale child, and TT45 alone failed.
- TT46: a probe that compared only the source digest of the seal started the changed child, and TT46 alone failed.
- The mismatch line: a probe that read each refusal as a missing build failed TT45 and TT46.
- TT47: a probe that printed the label with no quotes failed TT47 alone with `my alpha`.
- TT48: three probes of `internal/worktree/exec.go` failed the current-build row. In the first probe, the child took the parent environment and carried `BENCH_RUN_BINARY`. In the second probe, the child took only `BENCH_KIT` from the parent and carried the kit root. In the third probe, the child did not get the wrapper of its tree, and `BENCH_WRAPPER` was empty.
- The primary control: a probe that removed the primary guard refused `--in primary` in a kit primary checkout as a missing build.
- TT50: a copy-aside edit made `worktreeBuild` ignore the `freshness.Verify` refusal. The system suite then started the stale build at exit 0. The author restored the file from the copy and confirmed the bytes with `cmp`.

The two absence checks of TT48 bite only when the parent environment carries `BENCH_RUN_BINARY` and `BENCH_KIT`. `bench test` gives both to the Go child, and the probe output shows both values. A plain `go test` run gives neither, so there the two checks pass with no bite.

### Deviations

- Non-behavioral reading, for reviewer veto: the spec and the ticket name `sanitize.ShellQuote` for the label in the repair command. That function always quotes, but TT44 fixes `alpha` with no quotes. No exported function of `internal/sanitize` tells when a value needs quoting. `axi.ShellQuote` quotes a value only when the value needs it. Its comment names it the one derivation of the kit's shell quoting for a command line that a reader can run again. So the build uses `axi.ShellQuote` and adds no second copy of the rule for safe characters.
- The output of `axi.ShellQuote` is the same for each label that the spec names. For a label with a single quote, the two functions write different escapes that the shell reads as the same word. `internal/treetarget` now imports `internal/axi`, which the import list of the spec does not name.
- The kit check reads the root that `canonicalpath.Resolve` gives, because ticket 4 already runs the child in that root. `freshness` refuses a path with a symbolic link in it, so the physical root keeps a linked Bench home usable.
- A build is missing when `os.Lstat` of the published path reports no file after a `Verify` refusal. A missing seal or a missing build input gives the line `worktree build does not match the tree`.
- `TestRunKitWorktreeBuild` adds one control that no coverage row names: in a kit primary checkout with no build, `--in primary` starts the wrapper. The spec states that the primary target never runs a worktree build.
- The system row writes its build as a script of one line that makes a marker. The marker-script harness of ticket 4 is in a test file of `internal/treetarget`, so the system suite cannot import it. The system row needs only the presence of the marker.

The sweep of duplicated facts found no second copy in the delta. `kittest.WriteTree` is the one kit tree for three callers in two packages. `kittest.EditBuildInput` is the one stale edit for two callers. `kittest.Wrapper` is the one path that the fixture writes and that the TT48 expectation reads. `freshness.PublishedExecutable`, `runbinary.Env`, and `worktree.WrapperEnv` name the build path and two of the three variables.

The independent expectations are the two refusal lines and the repair command with and without quotes. The child argv, the child directory, and the three facts of the child environment are independent expectations too. The log above records a red for each. The system row repeats the stale-build expectation of `TestRunKitWorktreeBuild`, because it grades the real executable. Its red is the copy-aside run. The two `freshness.Publish` calls stay in test files, because the publication topology check refuses a call from any other file that is not a test.

### Probe verdicts

Each probe ran through `bench probe`, and each restore reads `yes`. The probes ran on the source that `c587623b` commits, with one exception: the package comment of `internal/treetarget/identify.go` changed after the probes. The first row is the plan probe `5-build-probe`, with the exact plan command. The second row is the author's own probe of the central property.

| File | Mutation | Test | Row | Verdict |
|---|---|---|---|---|
| `internal/treetarget/build.go` | swap: `worktree build is missing` to `worktree build is absent` | TestRunKitWorktreeBuild | TT44, TT47 | bit |
| `internal/treetarget/build.go` | self-probe, swap: ignore the `freshness.Verify` refusal | TestRunKitWorktreeBuild | TT44, TT45, TT46, TT47 | bit |
| `internal/treetarget/build.go` | swap: ignore the build inputs | TestRunKitWorktreeBuild | TT43 and the four refusal rows | bit |
| `internal/treetarget/build.go` | swap: start the wrapper after a refusal | TestRunKitWorktreeBuild | TT44, TT45, TT46, TT47 | bit |
| `internal/treetarget/build.go` | swap: check the executable digest only | TestRunKitWorktreeBuild | TT45 | bit |
| `internal/treetarget/build.go` | swap: check the source digest only | TestRunKitWorktreeBuild | TT46 | bit |
| `internal/treetarget/build.go` | swap: read each refusal as a missing build | TestRunKitWorktreeBuild | TT45, TT46 | bit |
| `internal/treetarget/build.go` | swap: print the label with no quotes | TestRunKitWorktreeBuild | TT47 | bit |
| `internal/treetarget/build.go` | swap: always quote the label | TestRunKitWorktreeBuild | TT44, TT45, TT46 | bit |
| `internal/treetarget/build.go` | swap: remove the primary guard | TestRunKitWorktreeBuild | primary control | bit |
| `internal/worktree/exec.go` | swap: take the parent environment | TestRunKitWorktreeBuild | TT48 | bit |
| `internal/worktree/exec.go` | swap: take `BENCH_KIT` from the parent | TestRunKitWorktreeBuild | TT48 | bit |
| `internal/worktree/exec.go` | swap: remove the wrapper of the tree | TestRunKitWorktreeBuild | TT48 | bit |

The system row TT50 took a copy-aside edit, because `bench probe` does not take the system suite. The log above names the edit.

### Verification

The author ran each check below on the source of `c587623b`, and each passed:

- `bench test --package ./internal/treetarget`, in 1394 ms, where the run before the production edit failed 5 subtests;
- `bench test --check system`, in 54232 ms, where the run before the production edit and the copy-aside run each failed TT50;
- `bench test --package ./cmd/bench`, in 14032 ms;
- `bench test --package ./internal/conformance --run TestRootConformance`, in 7039 ms;
- `bench structure --growth 79e79c51`, which reported that no source file grew past its budget.

No new test can skip, so the author did not run `bench test --check skip-ownership`. The commit ran with `--preflight-build tree-targets`. The lane ran gofmt, vet, build, and structure, and it passed. The worktree build was green, and `bench preflight build tree-targets` reported 15 green checks, 1 check that does not apply, and 0 red checks.

## TT-C4 chunk review, round 1

The frozen pair is base `9e742c5434c327e7174f85dd427331004b2bcef6` and tip `3831024c361eab0bad94374b8fb864bd85a76c4e`. The shared evidence is `sha256:4121fa0407fc63e9321af8d4f62d5ca3dc09bdf7f151335bb751bd23008c2ce0`. By the reviewer's direction, each axis ran as a fresh `bench-reviewer` session on fable at high effort. Only the Coverage axis ran probes, and the tree stayed clean.

The raw finding counts are 4 for Standards, 2 for Spec, and 3 for Coverage. No two findings name one fix. The de-duplicated repair-target count is 8, because P2 needs no repair.

## Standards

Findings: 4. Worst: the ticket 5 build refusal prints its own refusal shape and a third spelling of the `bench worktree build` verb.

- C4-S1, confidence 7, auto-fix: `TreeTarget` builds the ambiguity id list with the same four lines as `selectAssignment` (`internal/worktree/tree_target.go:42-46`, `internal/worktree/path.go:152-156`). The rule is one source per fact in `AGENTS.md`. The fix needs `internal/worktree/path.go` in the fence.
- C4-S2, confidence 6, auto-fix: `printBuildRefusal` prints the refusal and its `next=` line by hand (`internal/treetarget/build.go:50-53`). `printTargetRefusal` owns that shape (`internal/worktree/path.go:171-189`), and `usage.WorktreeBuild` owns the verb spelling (`internal/usage/worktree.go:15`). The fix needs an owner seam in `internal/worktree`.
- C4-S3, confidence 6, auto-fix: the new constant `WrapperEnv` (`internal/worktree/exec.go:212`) names the variable that `env.WrapperRouting` already owns (`internal/env/wrapper.go:10`). The fix needs `internal/env/wrapper.go` in the fence.
- C4-S4, confidence 5, ask-user: `kittest.WriteTree` restates the manifest path and the line grammar that `internal/freshness` owns (`internal/treetarget/kittest/kittest.go:26-46`, `internal/freshness/freshness_buildinputs.go:16,98-119`). It is also a third kit-tree fixture. The home of the helper is a reviewer decision.

## Spec

Findings: 2. Worst: the spec names `sanitize.ShellQuote` for the repair label, but that function always quotes, and TT44 fixes a bare `alpha`.

- C4-P1, confidence 6, ask-user: `specs/tree-targets/spec.md` names `sanitize.ShellQuote` for the label and omits `internal/axi` from the import list of `internal/treetarget`. The build uses `axi.ShellQuote`, which meets TT44 and TT47. The behavior is met, and the spec text is wrong at two places.
- C4-P2, confidence 6, no-op: TT34 says that stderr names the released state. The printed line is the state refusal of the current worktree resolver, which step 5 and the lookup rule of the spec require. The row text is a paraphrase, and the tree convention stays, flagged for reviewer veto.

Each other TT-C4 row is met at its named seam.

## Coverage

Findings: 3. Worst: no test pins the missing-tree or the creation-bundle refusal on the `--in` path, so a lookup that drops both checks stays green.

- C4-C1, confidence 8, auto-fix: no `--in` test covers an active label with a removed worktree or an invalid creation bundle. The probe that kept only the state check was silent in `./internal/worktree` and `./internal/treetarget`. A row must pin the `worktree tree is missing` refusal, its `next=` line, and no marker.
- C4-C2, confidence 7, auto-fix: no test asserts the child's `BENCH_HOME`. The probe that passed an empty home to `RunTreeChild` was silent in `./internal/treetarget`. The child environment row must name the resolved home.
- C4-C3, confidence 5, ask-user: `bench setup --in <label>` at an interactive terminal can hang. The shared runner puts the child in its own process group, so a terminal read stops the child, and the parent waits with no deadline. The same defect exists for `bench worktree exec <label> -- bench setup`. The route is a reviewer decision.

## TT-C4 round 1 advice

- `TreeTarget` selects the record again by id through `resolveAssignmentIn`. So a sibling whose label equals that id can turn a unique label match into an ambiguity refusal.
- TT36 tests only U+0001. A table over tab, newline, and return would pin `lineSafe` for the new caller.
- Step 4 is pinned only by the system row TT28. A unit row that calls `--in primary` from a linked worktree root would pin it at the fast seam.
- The TT48 absence checks bite only when the parent sets `BENCH_RUN_BINARY` and `BENCH_KIT`. The subtest can set both itself.
- No test pins `--in <label>` outside a repository, which prints the not-in-repo line at exit 1.
- The `test` help row shows the flag only on its first alternative.
- `valueUsage` escapes a backslash, so the usage value differs from the literal `toon.Usage` form.
- `kittest.WriteTree` writes each file with mode `0o755`, and only the two scripts need the execute bit.
- The kit holds two `ShellQuote` derivations, in `internal/axi` and `internal/sanitize`.
- `internal/treetarget/treetargettest` imports `internal/treetarget`, against the rule that only `cmd/bench` imports it. TT-C3 landed this import.

## TT-C4 repair routing

This is cycle 1 of the two repair cycles for chunk TT-C4. The reviewer decided C4-S4 and C4-C3 on 2026-09-30. Each affected ticket gets one fresh `bench-writer` repair session on opus at xhigh effort, with a cap of 3 attempts. That line matches the ticket 4 and ticket 5 authors. The ticket 4 session runs first, and the ticket 5 session starts after the ticket 4 repair commits green.

| Target | Finding | Ticket | Repair |
|---|---|---|---|
| R1 | C4-S1 | 4 | `TreeTarget` and `selectAssignment` build the ambiguity id list through one helper. |
| R2 | C4-S3 | 4 | `internal/env` names the wrapper variable once. `execEnv` and `cmd/bench/tree_scope.go` use that name, and `worktree.WrapperEnv` goes. |
| R3 | C4-C1 | 4 | Tests pin `--in` on an active label with a removed worktree and on an invalid creation bundle. Each test asserts the refusal, its `next=` line, and no marker. |
| R4 | C4-S2 | 5 | `printBuildRefusal` prints through the `internal/worktree` refusal printer, and `usage.WorktreeBuild` gives the verb spelling. TT44 and TT47 keep their exact lines. |
| R5 | C4-S4 | 5 | `internal/freshness` exports the manifest path and the line form, and `kittest.WriteTree` uses them. |
| R6 | C4-C2 | 5 | The TT48 test asserts that the child's `BENCH_HOME` is the resolved home. |

The reviewer decisions and the other dispositions are these:

- C4-S4: the reviewer chose the narrow repair R5. The other kit-tree fixtures stay, and a parked idea records their consolidation.
- C4-C3: the reviewer routed the terminal hang to spec B, which owns the exec route. A learning records it, and TT-C4 has no repair target for it.
- C4-P1: the plan commit corrects the spec text to name `axi.ShellQuote` and the `internal/axi` import edge. This is a non-behavioral contradiction, so the tree convention stays, flagged for reviewer veto.
- C4-P2: no repair. The tree convention stays, flagged for reviewer veto.

The plan commit expands the fences before dispatch. Ticket 4 gains `internal/worktree/path.go` and `internal/env/wrapper.go`. Ticket 5 gains `internal/worktree/path.go`, `internal/worktree/tree_target.go`, and `internal/freshness/freshness_buildinputs.go`. A learning records the expansion.

## TT-C4 ticket 4 repair evidence, cycle 1

The session `claude:bench-writer/tt-t4-repair-1` ran on opus at xhigh effort, with a cap of 3 attempts. It started at `63666844` and committed `764a1966` on a lane pass in the first attempt. This repair is cycle 1 of the two repair cycles for chunk TT-C4.

- R1: the new helper `ambiguousAssignments` in `internal/worktree/path.go` makes the ambiguity refusal from the matched assignments. `TreeTarget` and `selectAssignment` both call it, so the id list has one source.
- R2: `internal/env/wrapper.go` declares the constant `WrapperEnv`, and `WrapperRouting` uses it. `execEnv`, `runInTreeTarget`, the system helper `runTreeTarget`, and two tests read that constant. The constant `worktree.WrapperEnv` is removed. The other `"BENCH_WRAPPER"` literals are outside the fence, so they stay.
- R3: `TestRunRefusesBeforeChild` has two new rows. Each row asserts exit 1, one refusal line, its `next=` line, and no marker. The row `missing tree` removes the tree of the active label `gone`, and it wants `worktree tree is missing` and `next=bench worktree clean --landed`. The branch of `gone` is the tip of `main`, so the recovery route is the batch clean. The row `creation bundle` detaches HEAD in the tree of the active label `detached`, and it wants `assignment branch is not checked out` and `next=bench worktree list`. The review probe for C4-C1 was silent in round 1, and it now bites.

The sweep of duplicated facts found no second source in the delta. Two exec tests in `internal/worktree` keep the literal `BENCH_WRAPPER`, and their file is outside the fence. They are independent expectations, and the R2 probe records a red for each. The independent expectations of the new rows are the two refusal lines and the two `next=` lines. The probes below record a red for each.

### Probe verdicts

Each probe ran through `bench probe` on the source that `764a1966` commits, and each restore reads `yes`. The first row is the plan probe `4-flag-probe`, which the session ran again with the exact plan command.

| Target | File | Mutation | Test | Verdict |
|---|---|---|---|---|
| TT29, TT57 | `internal/treetarget/flag.go` | swap: `<label\|primary>` to `<label>` | TestRunRefusesBeforeChild | bit |
| R3 | `internal/worktree/tree_target.go` | swap: the shared resolver call to the state check alone | TestRunRefusesBeforeChild | bit |
| R3 | `internal/worktree/path.go` | swap: the missing-tree condition to `false && …` | TestRunRefusesBeforeChild | bit |
| R3 | `internal/worktree/path.go` | swap: the bundle condition to `false && …` | TestRunRefusesBeforeChild | bit |
| R3 | `internal/worktree/path.go` | swap: the recovery line to its `why` text | TestRunRefusesBeforeChild | bit |
| R1 | `internal/worktree/path.go` | swap: `a.ID` to `a.Label` in `ambiguousAssignments` | TestRunRefusesBeforeChild | bit |
| R1 | `internal/worktree/path.go` | swap: `a.ID` to `a.Label` in `ambiguousAssignments` | TestTargetVerbsNameTheResolverReason | bit |
| R2 | `internal/env/wrapper.go` | swap: the `WrapperEnv` value to `BENCH_WRAPPER_X` | TestExecChildDropsWrapperRouting, TestExecPWDMatchesChildDirectory | bit |

The state-check probe failed both new rows. The `detached` row started a child at exit 0, and the `gone` row reached the child start. The missing-tree probe failed only the `missing tree` row, with the owner-marker refusal. The bundle probe failed only the `creation bundle` row, and a child started. The two R1 probes failed TT35 and the `ambiguous` resolver case, so both callers use the helper. The R2 probe failed four tests, because the strip and the child value both read the one constant.

A first try of the recovery-line probe did not compile, and its verdict was `invalid`. The table holds the rerun. The session did not probe the name in `runInTreeTarget`, because a wrong name there starts the `cmd/bench` test binary as the child.

### Verification

The session ran each ticket 4 plan verification on the source that `764a1966` commits, and each passed:

- `bench test --package ./internal/treetarget`, in 1051 ms;
- `bench test --package ./internal/worktree`, in 48931 ms, with two environment capability skips for unix sockets;
- `bench test --package ./internal/canonicalpath`, in 4 ms, where the ticket 4 author record shows a red revert in 4 ms, so the tests run;
- `bench test --package ./cmd/bench`, in 15028 ms;
- `bench test --check system`, in 49940 ms, and in 49738 ms again after the worktree build;
- the plan probe `4-flag-probe`, which bit.

The session also ran these checks, and each passed:

- `bench test --package ./internal/env`, in 854 ms;
- `bench test --package ./internal/conformance --run TestRootConformance`, in 6298 ms;
- `bench structure --growth f981cd3d`, which reported that no source file grew past its budget;
- `gofmt -l` and `go vet -tags system` on the changed packages.

No new row can skip, so the session did not run `bench test --check skip-ownership`. The commit ran with `--preflight-build tree-targets`. The lane ran gofmt, vet, build, and structure, and it passed. The worktree build was green, and the preflight reported 15 green checks, 1 check that does not apply, and 0 red checks. A later `bench preflight build tree-targets` reported 14 green checks, 1 check that does not apply, and 0 red checks.

`internal/treetarget/build_test.go` is a ticket 5 file, but the fence holds `internal/treetarget/` as a prefix. The removal of `worktree.WrapperEnv` needs its change, or the package does not compile. The `creation bundle` row detaches HEAD, which fails the branch predicate of the bundle. The owner-marker component needs a helper of `internal/worktree` or a second spelling of the marker path, so the session did not use it.

```bench-review-record
{
  "version": 2,
  "spec": "specs/tree-targets/spec.md",
  "plan_digest": "sha256:8e3a09402b758dab396e7e37828fde90bc8eee75ecd288aab9535bffd0c92886",
  "implementation_session": "",
  "chunks": [
    {
      "id": "TT-C1",
      "base": "f981cd3db1a4f27eddba5feede40a80e205bf1c1",
      "tip": "196e1cda291a9d974f881435f6e1284b1838b241",
      "plan_digest": "sha256:82131a586bf62c0c3278a5dd81e06c3709beb7d6f263935fa90651ae0c7eccc8",
      "source_digest": "f0adccfef21bee43abd80bd3bdc3e92edfc7ecac",
      "acceptance_rows": [
        "TT1",
        "TT2",
        "TT3",
        "TT4",
        "TT5",
        "TT6",
        "TT7"
      ],
      "verification": [
        {
          "id": "tt-c1-1-cmd-r1",
          "performer": "claude:bench-writer/tt-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "48cf94fbd9d931f23b9390333303f8cabd76fad6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t1-author-20260929/1-cmd@166a3a8f",
            "digest": "sha256:b8706b37336cad03b084fddd17aff5f746edb6b943cfc5c83ebfe709a7766696",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,pass,16090\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "1-cmd",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "tt-c1-1-routing-r1",
          "performer": "claude:bench-writer/tt-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "48cf94fbd9d931f23b9390333303f8cabd76fad6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t1-author-20260929/1-routing@166a3a8f",
            "digest": "sha256:915a63e3adf95500b076f300919bf9a71615b398618b71408f9c5bc26dab14d8",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,20\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "1-routing",
          "command": "bench test --check subcommand-routing",
          "exit_code": 0
        },
        {
          "id": "tt-c1-1-conformance-r1",
          "performer": "claude:bench-writer/tt-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "48cf94fbd9d931f23b9390333303f8cabd76fad6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t1-author-20260929/1-conformance@166a3a8f",
            "digest": "sha256:2ccce730465fdae97c867b5b5ecec44425db5d7d0e68002847654e251e563f68",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,51673\nfailures[0]{package,test,line}:\nskips[3]{package,test,reason}:"
          },
          "requirement": "1-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "tt-c1-1-scope-probe-r1",
          "performer": "claude:bench-writer/tt-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "48cf94fbd9d931f23b9390333303f8cabd76fad6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t1-author-20260929/1-scope-probe@166a3a8f",
            "digest": "sha256:ca4cb4929604998da4cf09f2ae158ca735484565de97eef0e59bd357d08297ab",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,cmd/bench/tree_scope.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestRepositoryVerbRefusesTreeTarget,passed,3\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,fail,55"
          },
          "requirement": "1-scope-probe",
          "command": "bench probe cmd/bench/tree_scope.go --swap '\"--in\"' --with '\"--in-x\"' --package ./cmd/bench --run TestRepositoryVerbRefusesTreeTarget",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/tt-t1-author-20260929/1-scope-probe@166a3a8f",
              "digest": "sha256:ca4cb4929604998da4cf09f2ae158ca735484565de97eef0e59bd357d08297ab",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,cmd/bench/tree_scope.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestRepositoryVerbRefusesTreeTarget,passed,3\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,fail,55"
            }
          }
        },
        {
          "id": "tt-c1-author-probe-tt3-leaf-r1",
          "performer": "claude:bench-writer/tt-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "48cf94fbd9d931f23b9390333303f8cabd76fad6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t1-author-20260929/tt3-leaf-probe@166a3a8f",
            "digest": "sha256:d52c3fad433dc890602431309f3217fb4f9944bad20b8cb108351e881f4f6f51",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/conformance/command_scope_test.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/conformance,TestCommandScopeCheckBites,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,fail,5"
          },
          "requirement": "author-probe-TT3-leaf",
          "command": "bench probe internal/conformance/command_scope_test.go --swap 'if len(leaf.fields[\"Scope\"]) == 0 {' --with 'if len(leaf.fields[\"Scope\"]) < 0 {' --package ./internal/conformance --run TestCommandScopeCheckBites",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/tt-t1-author-20260929/tt3-leaf-probe@166a3a8f",
              "digest": "sha256:d52c3fad433dc890602431309f3217fb4f9944bad20b8cb108351e881f4f6f51",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/conformance/command_scope_test.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/conformance,TestCommandScopeCheckBites,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,fail,5"
            }
          }
        },
        {
          "id": "tt-c1-author-probe-tt2-public-r1",
          "performer": "claude:bench-writer/tt-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "48cf94fbd9d931f23b9390333303f8cabd76fad6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t1-author-20260929/tt2-public-probe@166a3a8f",
            "digest": "sha256:3064fcd3c3bc49a39214018ad6c41a263ea660bb39faaaab910d795894018d0a",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/conformance/command_scope_test.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/conformance,TestCommandScopeCheckBites,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,fail,10"
          },
          "requirement": "author-probe-TT2-public",
          "command": "bench probe internal/conformance/command_scope_test.go --swap 'case len(entry.fields[\"Inventory\"]) == 1 && !declared:' --with 'case len(entry.fields[\"Inventory\"]) == 2 && !declared:' --package ./internal/conformance --run TestCommandScopeCheckBites",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/tt-t1-author-20260929/tt2-public-probe@166a3a8f",
              "digest": "sha256:3064fcd3c3bc49a39214018ad6c41a263ea660bb39faaaab910d795894018d0a",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/conformance/command_scope_test.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/conformance,TestCommandScopeCheckBites,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,fail,10"
            }
          }
        },
        {
          "id": "tt-c1-1-cmd-r2",
          "performer": "claude:bench-writer/tt-t1-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f0adccfef21bee43abd80bd3bdc3e92edfc7ecac",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t1-repair-1-20260929/1-cmd@258edfd9",
            "digest": "sha256:47a0b8ad22f673b21149bd428e41c53819faf4f56bc3ba64aa2a6b856eeccd60",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,pass,12624\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "1-cmd",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "tt-c1-1-routing-r2",
          "performer": "claude:bench-writer/tt-t1-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f0adccfef21bee43abd80bd3bdc3e92edfc7ecac",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t1-repair-1-20260929/1-routing@258edfd9",
            "digest": "sha256:08f59d6fd9e57035cbbe09a2be129894319125639069e3eee8be894e9faac2ef",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,19\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "1-routing",
          "command": "bench test --check subcommand-routing",
          "exit_code": 0
        },
        {
          "id": "tt-c1-1-conformance-r2",
          "performer": "claude:bench-writer/tt-t1-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f0adccfef21bee43abd80bd3bdc3e92edfc7ecac",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t1-repair-1-20260929/1-conformance@258edfd9",
            "digest": "sha256:56bb76195d74696b17d505649cd1ca014cd6b755a8a134c857cd4aa119968d9e",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,39349\nfailures[0]{package,test,line}:\nskips[3]{package,test,reason}:"
          },
          "requirement": "1-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "tt-c1-1-scope-probe-r2",
          "performer": "claude:bench-writer/tt-t1-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f0adccfef21bee43abd80bd3bdc3e92edfc7ecac",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t1-repair-1-20260929/1-scope-probe@258edfd9",
            "digest": "sha256:edb47168d909814f26bffeeb231f307a16745017c23a062dd3caf538ed65d630",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,cmd/bench/tree_scope.go,swap,failed,10,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestRepositoryVerbRefusesTreeTarget,passed,28"
          },
          "requirement": "1-scope-probe",
          "command": "bench probe cmd/bench/tree_scope.go --swap '\"--in\"' --with '\"--in-x\"' --package ./cmd/bench --run TestRepositoryVerbRefusesTreeTarget",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/tt-t1-repair-1-20260929/1-scope-probe@258edfd9",
              "digest": "sha256:edb47168d909814f26bffeeb231f307a16745017c23a062dd3caf538ed65d630",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,cmd/bench/tree_scope.go,swap,failed,10,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestRepositoryVerbRefusesTreeTarget,passed,28"
            }
          }
        },
        {
          "id": "tt-c1-repair-probe-r3-narrow-r2",
          "performer": "claude:bench-writer/tt-t1-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f0adccfef21bee43abd80bd3bdc3e92edfc7ecac",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t1-repair-1-20260929/r3-narrow-probe@258edfd9",
            "digest": "sha256:d7c512e292145ea15ee442f873b72db9dd94de05ce0e0eabf19e1dd3c6382cc1",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,cmd/bench/tree_scope.go,swap,failed,8,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestRepositoryVerbRefusesTreeTarget,passed,28"
          },
          "requirement": "repair-probe-R3-narrow",
          "command": "bench probe cmd/bench/tree_scope.go --swap 'definition.Scope != scopeRepository ||' --with 'definition.Name != \"version\" ||' --package ./cmd/bench --run TestRepositoryVerbRefusesTreeTarget",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/tt-t1-repair-1-20260929/r3-narrow-probe@258edfd9",
              "digest": "sha256:d7c512e292145ea15ee442f873b72db9dd94de05ce0e0eabf19e1dd3c6382cc1",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,cmd/bench/tree_scope.go,swap,failed,8,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestRepositoryVerbRefusesTreeTarget,passed,28"
            }
          }
        },
        {
          "id": "tt-c1-repair-probe-r2-predicate-r2",
          "performer": "claude:bench-writer/tt-t1-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f0adccfef21bee43abd80bd3bdc3e92edfc7ecac",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t1-repair-1-20260929/r2-predicate-probe@258edfd9",
            "digest": "sha256:d52c3fad433dc890602431309f3217fb4f9944bad20b8cb108351e881f4f6f51",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/conformance/command_scope_test.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/conformance,TestCommandScopeCheckBites,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,fail,5"
          },
          "requirement": "repair-probe-R2-predicate",
          "command": "bench probe internal/conformance/command_scope_test.go --swap 'return ok && declaredScopes[value.Name]' --with 'return value != nil || !ok' --package ./internal/conformance --run TestCommandScopeCheckBites",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/tt-t1-repair-1-20260929/r2-predicate-probe@258edfd9",
              "digest": "sha256:d52c3fad433dc890602431309f3217fb4f9944bad20b8cb108351e881f4f6f51",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/conformance/command_scope_test.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/conformance,TestCommandScopeCheckBites,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,fail,5"
            }
          }
        },
        {
          "id": "tt-c1-repair-probe-r2-live-r2",
          "performer": "claude:bench-writer/tt-t1-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "f0adccfef21bee43abd80bd3bdc3e92edfc7ecac",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t1-repair-1-20260929/r2-live-probe@258edfd9",
            "digest": "sha256:64da265920299c029c9bb310031089464a36dd4b8ebf73dd8cdf73ac3c53db5e",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,cmd/bench/main.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,subcommand-routing,^TestRootConformance$,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,fail,20\nfailures[1]{package,test,line}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:29: gate: command \\\"models\\\" declares no scope\""
          },
          "requirement": "repair-probe-R2-live",
          "command": "./dist/bench probe cmd/bench/main.go --swap 'Bound: boundResponse, Scope: scopeRepository, Run: outputCommand(models.Command)' --with 'Bound: boundResponse, Scope: 0, Run: outputCommand(models.Command)' --check subcommand-routing",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/tt-t1-repair-1-20260929/r2-live-probe@258edfd9",
              "digest": "sha256:64da265920299c029c9bb310031089464a36dd4b8ebf73dd8cdf73ac3c53db5e",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,cmd/bench/main.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,subcommand-routing,^TestRootConformance$,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,fail,20\nfailures[1]{package,test,line}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:29: gate: command \\\"models\\\" declares no scope\""
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "tt-c1-standards-r1",
          "performer": "claude:bench-reviewer/tt-c1-standards-r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "48cf94fbd9d931f23b9390333303f8cabd76fad6",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/tt-c1-standards-r1@a08359f9",
            "digest": "sha256:1b7c505b57a1d115541cfa1a921376528a7b0546793765610614dfbe8da848b1",
            "excerpt": "Standards: 1 finding. Worst: the bite test adds a new private copy of a file-writing fixture helper that the same package already has."
          },
          "axis": "Standards",
          "base": "f981cd3db1a4f27eddba5feede40a80e205bf1c1",
          "tip": "a08359f9550f64d92d524dfebc3f7e6143a31118",
          "finding_ids": [
            "R1"
          ],
          "supersedes": []
        },
        {
          "id": "tt-c1-spec-r1",
          "performer": "claude:bench-reviewer/tt-c1-spec-r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "48cf94fbd9d931f23b9390333303f8cabd76fad6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-c1-spec-r1@a08359f9",
            "digest": "sha256:ab5b650410c319daa5c3c75b3d5c285fe8a123842460ae08651799828e0bd5e2",
            "excerpt": "Spec: 0 findings. Every TT1-TT7 row is met, and the classification matches the spec table exactly."
          },
          "axis": "Spec",
          "base": "f981cd3db1a4f27eddba5feede40a80e205bf1c1",
          "tip": "a08359f9550f64d92d524dfebc3f7e6143a31118",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "tt-c1-coverage-r1",
          "performer": "claude:bench-reviewer/tt-c1-coverage-r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "48cf94fbd9d931f23b9390333303f8cabd76fad6",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/tt-c1-coverage-r1@a08359f9",
            "digest": "sha256:c02fe3282a13d4e8a6123d01bc66a60a820c74944fd3afb5d131c0ca1c1e05d9",
            "excerpt": "Coverage: 1 finding. Worst: the scope check checks only that a Scope key is present, so a registry line with Scope: 0 passes both the gate and the full cmd/bench suite."
          },
          "axis": "Coverage",
          "base": "f981cd3db1a4f27eddba5feede40a80e205bf1c1",
          "tip": "a08359f9550f64d92d524dfebc3f7e6143a31118",
          "finding_ids": [
            "R2"
          ],
          "supersedes": []
        },
        {
          "id": "tt-c1-standards-r2",
          "performer": "claude:bench-reviewer/tt-c1-standards-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "f0adccfef21bee43abd80bd3bdc3e92edfc7ecac",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-c1-standards-r2@196e1cda",
            "digest": "sha256:1edd83e66a84c21d094be64f07753ea71a6b2bb24a2f7734f1995e7d8e3101d9",
            "excerpt": "Standards: 0 findings. All three folds hold, and the repair delta adds no duplicated knowledge that blocks."
          },
          "axis": "Standards",
          "base": "f981cd3db1a4f27eddba5feede40a80e205bf1c1",
          "tip": "196e1cda291a9d974f881435f6e1284b1838b241",
          "finding_ids": [],
          "supersedes": [
            "tt-c1-standards-r1"
          ]
        },
        {
          "id": "tt-c1-spec-r2",
          "performer": "claude:bench-reviewer/tt-c1-spec-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "f0adccfef21bee43abd80bd3bdc3e92edfc7ecac",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-c1-spec-r2@196e1cda",
            "digest": "sha256:7e4bda2be89752b8ef5cb073ac14fd3ec66eb7094f8cedecb7549e7ecfa6b1ee",
            "excerpt": "Spec: 0 findings. Nothing in the repair delta a08359f9..196e1cda breaks TT1-TT7."
          },
          "axis": "Spec",
          "base": "f981cd3db1a4f27eddba5feede40a80e205bf1c1",
          "tip": "196e1cda291a9d974f881435f6e1284b1838b241",
          "finding_ids": [],
          "supersedes": [
            "tt-c1-spec-r1"
          ]
        },
        {
          "id": "tt-c1-coverage-r2",
          "performer": "claude:bench-reviewer/tt-c1-coverage-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "f0adccfef21bee43abd80bd3bdc3e92edfc7ecac",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-c1-coverage-r2@196e1cda",
            "digest": "sha256:5684929f3d87c592dea74ce2ff19b7ee826e0dbb818dfa448499e24a6826d261",
            "excerpt": "Coverage: 0 findings. R1, R2, and R3 are confirmed, and no candidate in the repair delta survived refutation."
          },
          "axis": "Coverage",
          "base": "f981cd3db1a4f27eddba5feede40a80e205bf1c1",
          "tip": "196e1cda291a9d974f881435f6e1284b1838b241",
          "finding_ids": [],
          "supersedes": [
            "tt-c1-coverage-r1"
          ]
        }
      ]
    },
    {
      "id": "TT-C2",
      "base": "196e1cda291a9d974f881435f6e1284b1838b241",
      "tip": "74daf800f8a11d85083711b31aaabc47797ebc73",
      "plan_digest": "sha256:834182d7bb87b440ee0493ba5fd18fb70360c237bae9e7a24acfe73bc8fc91a0",
      "source_digest": "80fbc1c6ab2706d7744da73a60e0cb634a2b61aa",
      "acceptance_rows": [
        "TT51",
        "TT52"
      ],
      "verification": [
        {
          "id": "tt-c2-2-gate-r1",
          "performer": "claude:bench-writer/tt-t2-author-2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "673171b5be48cb53a7038841dc6fe6cd56c7a727",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t2-author-2-20260929/2-gate@e43254d5",
            "digest": "sha256:17c46a859802dedd68119f61c2de7da193d5c035dfabf55cae2ec5a68144950d",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,pass,14924\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "2-gate",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "tt-c2-2-kit-probe-r1",
          "performer": "claude:bench-writer/tt-t2-author-2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "673171b5be48cb53a7038841dc6fe6cd56c7a727",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t2-author-2-20260929/2-kit-probe@e43254d5",
            "digest": "sha256:7f3561247e94d373bec956747a595f7c2b89bc5a2bfb30f195fab60d2e04d586",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/kit_source.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestKitSourceCheckoutResolvesARelativeKit,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,fail,3"
          },
          "requirement": "2-kit-probe",
          "command": "bench probe internal/gate/kit_source.go --swap 'canonicalpath.Resolve' --with 'func(p string) (string, error) { _ = canonicalpath.Resolve; return filepath.EvalSymlinks(p) }' --package ./internal/gate --run TestKitSourceCheckoutResolvesARelativeKit",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/tt-t2-author-2-20260929/2-kit-probe@e43254d5",
              "digest": "sha256:7f3561247e94d373bec956747a595f7c2b89bc5a2bfb30f195fab60d2e04d586",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/kit_source.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestKitSourceCheckoutResolvesARelativeKit,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,fail,3"
            }
          }
        },
        {
          "id": "tt-c2-author-probe-tt52-link-r1",
          "performer": "claude:bench-writer/tt-t2-author-2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "673171b5be48cb53a7038841dc6fe6cd56c7a727",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t2-author-2-20260929/tt52-link@e43254d5",
            "digest": "sha256:fda4978fa5de5e49400f8434ca9a382d4c46e90359b85cd10f027ee244631b8f",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/kit_source.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestKitSourceCheckout,passed,2\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,fail,4"
          },
          "requirement": "author-probe-TT52-link",
          "command": "bench probe internal/gate/kit_source.go --swap 'resolvedPath(root)' --with 'filepath.Clean(root), true' --package ./internal/gate --run 'TestKitSourceCheckout'",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/tt-t2-author-2-20260929/tt52-link@e43254d5",
              "digest": "sha256:fda4978fa5de5e49400f8434ca9a382d4c46e90359b85cd10f027ee244631b8f",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/kit_source.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestKitSourceCheckout,passed,2\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,fail,4"
            }
          }
        },
        {
          "id": "tt-c2-author-probe-tt52-unrelated-r1",
          "performer": "claude:bench-writer/tt-t2-author-2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "673171b5be48cb53a7038841dc6fe6cd56c7a727",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t2-author-2-20260929/tt52-unrelated@e43254d5",
            "digest": "sha256:fda4978fa5de5e49400f8434ca9a382d4c46e90359b85cd10f027ee244631b8f",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/kit_source.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestKitSourceCheckout,passed,2\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,fail,4"
          },
          "requirement": "author-probe-TT52-unrelated",
          "command": "bench probe internal/gate/kit_source.go --swap 'return ok && resolvedRoot == resolvedKit' --with 'return ok && resolvedRoot+resolvedKit != \"\"' --package ./internal/gate --run 'TestKitSourceCheckout'",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/tt-t2-author-2-20260929/tt52-unrelated@e43254d5",
              "digest": "sha256:fda4978fa5de5e49400f8434ca9a382d4c46e90359b85cd10f027ee244631b8f",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/kit_source.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestKitSourceCheckout,passed,2\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,fail,4"
            }
          }
        },
        {
          "id": "tt-c2-2-gate-r2",
          "performer": "claude:bench-writer/tt-t2-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "80fbc1c6ab2706d7744da73a60e0cb634a2b61aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t2-repair-1-20260929/2-gate@9958ef57",
            "digest": "sha256:4c41fe5ef4b8170f76c7de34b809b9abdc3232c0c369b0ce6008255a9e770db0",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,pass,15012\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "2-gate",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "tt-c2-2-kit-probe-r2",
          "performer": "claude:bench-writer/tt-t2-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "80fbc1c6ab2706d7744da73a60e0cb634a2b61aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t2-repair-1-20260929/2-kit-probe@9958ef57",
            "digest": "sha256:dbd03c9fd1d394189f42d4b936d0487b0567f29294ff483509ae30bd5628b202",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/kit_source.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestKitSourceCheckoutResolvesARelativeKit,passed,2\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,fail,4"
          },
          "requirement": "2-kit-probe",
          "command": "bench probe internal/gate/kit_source.go --swap 'canonicalpath.Resolve' --with 'func(p string) (string, error) { _ = canonicalpath.Resolve; return filepath.EvalSymlinks(p) }' --package ./internal/gate --run TestKitSourceCheckoutResolvesARelativeKit",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/tt-t2-repair-1-20260929/2-kit-probe@9958ef57",
              "digest": "sha256:dbd03c9fd1d394189f42d4b936d0487b0567f29294ff483509ae30bd5628b202",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/kit_source.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestKitSourceCheckoutResolvesARelativeKit,passed,2\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,fail,4"
            }
          }
        },
        {
          "id": "tt-c2-repair-probe-r1-resolve-error-r2",
          "performer": "claude:bench-writer/tt-t2-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "80fbc1c6ab2706d7744da73a60e0cb634a2b61aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t2-repair-1-20260929/r1-resolve-error-probe@9958ef57",
            "digest": "sha256:4cce62584c926c2604b00e1671176b0bdc5f482fb158517f751c30384e59c54b",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/kit_source.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestKitSourceCheckout,passed,4\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,fail,5\nfailures[1]{package,test,line}:\n  github.com/gibbonmi/bench/internal/gate,TestKitSourceCheckoutAnswersFalseOnAResolveError,\"kit_source_test.go:84: KitSourceCheckout(\\\".\\\") with a deleted working directory = true, want false\""
          },
          "requirement": "repair-probe-R1-resolve-error",
          "command": "bench probe internal/gate/kit_source.go --swap 'return resolved, err == nil' --with '_ = err; return resolved, true' --package ./internal/gate --run TestKitSourceCheckout",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/tt-t2-repair-1-20260929/r1-resolve-error-probe@9958ef57",
              "digest": "sha256:4cce62584c926c2604b00e1671176b0bdc5f482fb158517f751c30384e59c54b",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/kit_source.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestKitSourceCheckout,passed,4\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,fail,5\nfailures[1]{package,test,line}:\n  github.com/gibbonmi/bench/internal/gate,TestKitSourceCheckoutAnswersFalseOnAResolveError,\"kit_source_test.go:84: KitSourceCheckout(\\\".\\\") with a deleted working directory = true, want false\""
            }
          }
        },
        {
          "id": "tt-c2-repair-probe-r2-root-join-r2",
          "performer": "claude:bench-writer/tt-t2-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "80fbc1c6ab2706d7744da73a60e0cb634a2b61aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t2-repair-1-20260929/r2-root-join-probe@9958ef57",
            "digest": "sha256:cb3abc032eaba68dc252a9f2b0032298ada73de1277d4750d32f0ac45e5179ad",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/kit_source.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestKitSourceCheckout,passed,4\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,fail,5\nfailures[1]{package,test,line}:\n  github.com/gibbonmi/bench/internal/gate,TestKitSourceCheckoutResolvesARelativeKitAgainstTheWorkingDirectory,..."
          },
          "requirement": "repair-probe-R2-root-join",
          "command": "bench probe internal/gate/kit_source.go --swap 'kit := KitDir()' --with 'kit := KitDir(); if !filepath.IsAbs(kit) { kit = filepath.Join(root, kit) }' --package ./internal/gate --run TestKitSourceCheckout",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/tt-t2-repair-1-20260929/r2-root-join-probe@9958ef57",
              "digest": "sha256:cb3abc032eaba68dc252a9f2b0032298ada73de1277d4750d32f0ac45e5179ad",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/kit_source.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestKitSourceCheckout,passed,4\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,fail,5\nfailures[1]{package,test,line}:\n  github.com/gibbonmi/bench/internal/gate,TestKitSourceCheckoutResolvesARelativeKitAgainstTheWorkingDirectory,..."
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "tt-c2-standards-r1",
          "performer": "claude:bench-reviewer/tt-c2-standards-r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "673171b5be48cb53a7038841dc6fe6cd56c7a727",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-c2-standards-r1@c8b9bb94",
            "digest": "sha256:cd216dc517b69fdccdfbc6d81170c4adf06ee6dfb62fe5259c168a90cf6e46db",
            "excerpt": "Standards: 1 item, classified as advice by the coordinator. The new test doc comment carries a row tag, which the tree uses in more than 20 test files."
          },
          "axis": "Standards",
          "base": "196e1cda291a9d974f881435f6e1284b1838b241",
          "tip": "c8b9bb94eb15029a84abd137e1fbac93e8990480",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "tt-c2-spec-r1",
          "performer": "claude:bench-reviewer/tt-c2-spec-r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "673171b5be48cb53a7038841dc6fe6cd56c7a727",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-c2-spec-r1@c8b9bb94",
            "digest": "sha256:964b583f87b5e59b3d03729acbb796195fe6cc03d7e826a7000b70cc96302976",
            "excerpt": "Spec: 0 findings. TT51 and TT52 are met, and the three recorded deviations are acceptable."
          },
          "axis": "Spec",
          "base": "196e1cda291a9d974f881435f6e1284b1838b241",
          "tip": "c8b9bb94eb15029a84abd137e1fbac93e8990480",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "tt-c2-coverage-r1",
          "performer": "claude:bench-reviewer/tt-c2-coverage-r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "673171b5be48cb53a7038841dc6fe6cd56c7a727",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/tt-c2-coverage-r1@c8b9bb94",
            "digest": "sha256:883be48740407f42cb1a7459945ac9b4fa8304bca2145f037bdb94150069786b",
            "excerpt": "Coverage: 1 finding. Worst: nothing guards the spec's resolve-error rule, and a mutation that removes it stays green."
          },
          "axis": "Coverage",
          "base": "196e1cda291a9d974f881435f6e1284b1838b241",
          "tip": "c8b9bb94eb15029a84abd137e1fbac93e8990480",
          "finding_ids": [
            "R1"
          ],
          "supersedes": []
        },
        {
          "id": "tt-c2-standards-r2",
          "performer": "claude:bench-reviewer/tt-c2-standards-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "80fbc1c6ab2706d7744da73a60e0cb634a2b61aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-c2-standards-r2@74daf800",
            "digest": "sha256:5d837ef413d374d65c9cbc67ba30448ceec30b9ccbbb83e605662944cabcdf60",
            "excerpt": "Standards: 0 findings. The repair delta adds two independent tests, their reds are recorded, and it duplicates no knowledge."
          },
          "axis": "Standards",
          "base": "196e1cda291a9d974f881435f6e1284b1838b241",
          "tip": "74daf800f8a11d85083711b31aaabc47797ebc73",
          "finding_ids": [],
          "supersedes": [
            "tt-c2-standards-r1"
          ]
        },
        {
          "id": "tt-c2-spec-r2",
          "performer": "claude:bench-reviewer/tt-c2-spec-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "80fbc1c6ab2706d7744da73a60e0cb634a2b61aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-c2-spec-r2@74daf800",
            "digest": "sha256:7689a934a4416af7c4fccc3144660a5029377a3e83cd2dbe97b40eb9270f1050",
            "excerpt": "Spec: 0 findings. Both repair folds meet their routing rows, and the repair delta stays inside ticket 2's Writes: fence."
          },
          "axis": "Spec",
          "base": "196e1cda291a9d974f881435f6e1284b1838b241",
          "tip": "74daf800f8a11d85083711b31aaabc47797ebc73",
          "finding_ids": [],
          "supersedes": [
            "tt-c2-spec-r1"
          ]
        },
        {
          "id": "tt-c2-coverage-r2",
          "performer": "claude:bench-reviewer/tt-c2-coverage-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "80fbc1c6ab2706d7744da73a60e0cb634a2b61aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-c2-coverage-r2@74daf800",
            "digest": "sha256:8cfa87072452331a0032e77d4728c7b97624e2a0a2e5b4b8093b0283e6a204c9",
            "excerpt": "Coverage: 0 findings. R1 and R2 are both confirmed by independent bypasses that the new tests catch."
          },
          "axis": "Coverage",
          "base": "196e1cda291a9d974f881435f6e1284b1838b241",
          "tip": "74daf800f8a11d85083711b31aaabc47797ebc73",
          "finding_ids": [],
          "supersedes": [
            "tt-c2-coverage-r1"
          ]
        }
      ]
    },
    {
      "id": "TT-C3",
      "base": "74daf800f8a11d85083711b31aaabc47797ebc73",
      "tip": "9e742c5434c327e7174f85dd427331004b2bcef6",
      "plan_digest": "sha256:49d7c3455ec98c0ab034c4bb473f14b0515ff390730eefb988e644e95f92b708",
      "source_digest": "ef1d5964df6ea99a501c802627517db74e027f95",
      "acceptance_rows": [
        "TT10",
        "TT11",
        "TT12",
        "TT13",
        "TT14",
        "TT15",
        "TT16",
        "TT17",
        "TT18",
        "TT19",
        "TT20",
        "TT21",
        "TT22",
        "TT23",
        "TT24",
        "TT25",
        "TT56",
        "TT60",
        "TT61"
      ],
      "verification": [
        {
          "id": "tt-c3-3-cmd-r1",
          "performer": "claude:bench-writer/tt-t3-author-2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8bf4434336ae7fb7c412c3cde6930803bbe93b01",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t3-author-2-20260929/3-cmd@93b7c071",
            "digest": "sha256:6638158506e5b2b04674be3479a32837a35ae1dcd20e6c1a635248f9fd8c0652",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,pass,12705\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "3-cmd",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "tt-c3-3-treetarget-r1",
          "performer": "claude:bench-writer/tt-t3-author-2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8bf4434336ae7fb7c412c3cde6930803bbe93b01",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t3-author-2-20260929/3-treetarget@93b7c071",
            "digest": "sha256:7d1c96ef15ef0ae02ae229b00910d58708306ff8aae8107af5e4539143cb8aed",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/treetarget,pass,114\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "3-treetarget",
          "command": "bench test --package ./internal/treetarget",
          "exit_code": 0
        },
        {
          "id": "tt-c3-3-responsebound-r1",
          "performer": "claude:bench-writer/tt-t3-author-2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8bf4434336ae7fb7c412c3cde6930803bbe93b01",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t3-author-2-20260929/3-responsebound@93b7c071",
            "digest": "sha256:c8a94ec4b6eb46bd99be2214c1d7119e94454dfc722ed4ff6a38b776dbe536a8",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/responsebound,pass,239\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "3-responsebound",
          "command": "bench test --package ./internal/responsebound",
          "exit_code": 0
        },
        {
          "id": "tt-c3-3-unassigned-probe-r1",
          "performer": "claude:bench-writer/tt-t3-author-2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8bf4434336ae7fb7c412c3cde6930803bbe93b01",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t3-author-2-20260929/3-unassigned-probe@93b7c071",
            "digest": "sha256:bd6a7cad29435f241c06ba9eb4ebf9c48af94713cff4652fef179e60999866c1",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/treetarget/identify.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/treetarget,TestIdentifyUnownedWorktree,passed,3\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/treetarget,fail,48"
          },
          "requirement": "3-unassigned-probe",
          "command": "bench probe internal/treetarget/identify.go --swap '\"unassigned\"' --with '\"primary\"' --package ./internal/treetarget --run TestIdentifyUnownedWorktree",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/tt-t3-author-2-20260929/3-unassigned-probe@93b7c071",
              "digest": "sha256:bd6a7cad29435f241c06ba9eb4ebf9c48af94713cff4652fef179e60999866c1",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/treetarget/identify.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/treetarget,TestIdentifyUnownedWorktree,passed,3\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/treetarget,fail,48"
            }
          }
        },
        {
          "id": "tt-c3-author-probe-tt56-exit2-r1",
          "performer": "claude:bench-writer/tt-t3-author-2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8bf4434336ae7fb7c412c3cde6930803bbe93b01",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t3-author-2-20260929/tt56-exit2@93b7c071",
            "digest": "sha256:d823d7b0c34bf474d2f56182712a411227b11e2ed5a4a821d0c3123d564c7857",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,cmd/bench/tree_scope.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestRunGateRejectsBriefUsage,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,fail,33"
          },
          "requirement": "author-probe-TT56-exit2",
          "command": "bench probe cmd/bench/tree_scope.go --swap 'if exit == 2 {' --with 'if exit == 99 {' --package ./cmd/bench --run TestRunGateRejectsBriefUsage",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/tt-t3-author-2-20260929/tt56-exit2@93b7c071",
              "digest": "sha256:d823d7b0c34bf474d2f56182712a411227b11e2ed5a4a821d0c3123d564c7857",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,cmd/bench/tree_scope.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestRunGateRejectsBriefUsage,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,fail,33"
            }
          }
        },
        {
          "id": "tt-c3-author-probe-placement-r1",
          "performer": "claude:bench-writer/tt-t3-author-2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8bf4434336ae7fb7c412c3cde6930803bbe93b01",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t3-author-2-20260929/placement@93b7c071",
            "digest": "sha256:a79c87b937584f4f8674823262a1c9647a49c3a1b3ff5d0e9f125ce80459032c",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,cmd/bench/census_output.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestTreeRowLeadsTreeResponse|TestTreeRowSurvivesSpill|TestTreeRowOutsideResponseBound,passed,3"
          },
          "requirement": "author-probe-TT25-TT60-placement",
          "command": "bench probe cmd/bench/census_output.go --swap 'owner.Lead(shownRow(row, exit))' --with '_, _ = c.Stdout.Write([]byte(shownRow(row, exit)))' --package ./cmd/bench --run 'TestTreeRowLeadsTreeResponse|TestTreeRowSurvivesSpill|TestTreeRowOutsideResponseBound'",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/tt-t3-author-2-20260929/placement@93b7c071",
              "digest": "sha256:a79c87b937584f4f8674823262a1c9647a49c3a1b3ff5d0e9f125ce80459032c",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,cmd/bench/census_output.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestTreeRowLeadsTreeResponse|TestTreeRowSurvivesSpill|TestTreeRowOutsideResponseBound,passed,3"
            }
          }
        },
        {
          "id": "tt-c3-author-probe-tt11-target-r1",
          "performer": "claude:bench-writer/tt-t3-author-2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8bf4434336ae7fb7c412c3cde6930803bbe93b01",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t3-author-2-20260929/tt11-target@93b7c071",
            "digest": "sha256:df2e26f046773628a2ebb545e4b9013bae4a6295d0dc02ffd4a218fd6429660e",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/treetarget/identify.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestTreeRowLeadsTreeResponse,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,fail,20"
          },
          "requirement": "author-probe-TT11-target",
          "command": "bench probe internal/treetarget/identify.go --swap 'identity.Target = \"primary\"' --with 'identity.Target = \"main\"' --package ./cmd/bench --run TestTreeRowLeadsTreeResponse",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/tt-t3-author-2-20260929/tt11-target@93b7c071",
              "digest": "sha256:df2e26f046773628a2ebb545e4b9013bae4a6295d0dc02ffd4a218fd6429660e",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/treetarget/identify.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestTreeRowLeadsTreeResponse,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,fail,20"
            }
          }
        },
        {
          "id": "tt-c3-author-probe-tt15-head-r1",
          "performer": "claude:bench-writer/tt-t3-author-2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8bf4434336ae7fb7c412c3cde6930803bbe93b01",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t3-author-2-20260929/tt15-head@93b7c071",
            "digest": "sha256:df2e26f046773628a2ebb545e4b9013bae4a6295d0dc02ffd4a218fd6429660e",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/treetarget/identify.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestTreeRowLeadsTreeResponse,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,fail,20"
          },
          "requirement": "author-probe-TT15-head",
          "command": "bench probe internal/treetarget/identify.go --swap 'identity.Head = head' --with 'identity.Head = head[:12]' --package ./cmd/bench --run TestTreeRowLeadsTreeResponse",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/tt-t3-author-2-20260929/tt15-head@93b7c071",
              "digest": "sha256:df2e26f046773628a2ebb545e4b9013bae4a6295d0dc02ffd4a218fd6429660e",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/treetarget/identify.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestTreeRowLeadsTreeResponse,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,fail,20"
            }
          }
        },
        {
          "id": "tt-c3-3-cmd-r2",
          "performer": "claude:bench-writer/tt-t3-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "564566b25b32a0cd36399364a16fd537671887c9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t3-repair-1-20260929/3-cmd@744ef656",
            "digest": "sha256:5eb85dd2dc451b0c96b38206d8df9517571cb0928b4f0636d81e20d6db3d9114",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,pass,13480\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "3-cmd",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "tt-c3-3-treetarget-r2",
          "performer": "claude:bench-writer/tt-t3-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "564566b25b32a0cd36399364a16fd537671887c9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t3-repair-1-20260929/3-treetarget@744ef656",
            "digest": "sha256:639b706779298abc145a0e675fd3def5df15c324c4f707dd4e280d349cf7873f",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/treetarget,pass,172\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "3-treetarget",
          "command": "bench test --package ./internal/treetarget",
          "exit_code": 0
        },
        {
          "id": "tt-c3-3-responsebound-r2",
          "performer": "claude:bench-writer/tt-t3-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "564566b25b32a0cd36399364a16fd537671887c9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t3-repair-1-20260929/3-responsebound@744ef656",
            "digest": "sha256:8c1f8bc69daea97ad2fc2fa98df0e2306d2fbd59f3128b89c4df97188e4a437d",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/responsebound,pass,268\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "3-responsebound",
          "command": "bench test --package ./internal/responsebound",
          "exit_code": 0
        },
        {
          "id": "tt-c3-3-unassigned-probe-r2",
          "performer": "claude:bench-writer/tt-t3-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "564566b25b32a0cd36399364a16fd537671887c9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t3-repair-1-20260929/3-unassigned-probe@744ef656",
            "digest": "sha256:bd6a7cad29435f241c06ba9eb4ebf9c48af94713cff4652fef179e60999866c1",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/treetarget/identify.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/treetarget,TestIdentifyUnownedWorktree,passed,3\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/treetarget,fail,48"
          },
          "requirement": "3-unassigned-probe",
          "command": "bench probe internal/treetarget/identify.go --swap '\"unassigned\"' --with '\"primary\"' --package ./internal/treetarget --run TestIdentifyUnownedWorktree",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/tt-t3-repair-1-20260929/3-unassigned-probe@744ef656",
              "digest": "sha256:bd6a7cad29435f241c06ba9eb4ebf9c48af94713cff4652fef179e60999866c1",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/treetarget/identify.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/treetarget,TestIdentifyUnownedWorktree,passed,3\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/treetarget,fail,48"
            }
          }
        },
        {
          "id": "tt-c3-repair-probe-r1-hostile-label-r2",
          "performer": "claude:bench-writer/tt-t3-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "564566b25b32a0cd36399364a16fd537671887c9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t3-repair-1-20260929/r1-hostile-label-probe@744ef656",
            "digest": "sha256:b7d039e0d721c873fcc235600faf27dc0e8221b79652844047e43c0df6d4839a",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/treetarget/identify.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/treetarget,TestRowEscapesHostileLabel,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/treetarget,fail,40\nfailures[1]{package,test,line}:\n  github.com/gibbonmi/bench/internal/treetarget,TestRowEscapesHostileLabel,\"identify_test.go:90: hostile label row = (\\\"\\\", toon: unsupported control character U+0007 in string), ..."
          },
          "requirement": "repair-probe-R1-hostile-label",
          "command": "bench probe internal/treetarget/identify.go --swap '\tif !sanitize.LineSafe(target) {' --with '\tif false && !sanitize.LineSafe(target) {' --package ./internal/treetarget --run TestRowEscapesHostileLabel",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/tt-t3-repair-1-20260929/r1-hostile-label-probe@744ef656",
              "digest": "sha256:b7d039e0d721c873fcc235600faf27dc0e8221b79652844047e43c0df6d4839a",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/treetarget/identify.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/treetarget,TestRowEscapesHostileLabel,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/treetarget,fail,40\nfailures[1]{package,test,line}:\n  github.com/gibbonmi/bench/internal/treetarget,TestRowEscapesHostileLabel,\"identify_test.go:90: hostile label row = (\\\"\\\", toon: unsupported control character U+0007 in string), ..."
            }
          }
        },
        {
          "id": "tt-c3-repair-probe-r2-exempt-exit2-r2",
          "performer": "claude:bench-writer/tt-t3-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "564566b25b32a0cd36399364a16fd537671887c9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t3-repair-1-20260929/r2-exempt-exit2-probe@744ef656",
            "digest": "sha256:9160321bccd66ade4d48f658c6fc870a7ae5f9b2ad4041159bd034200c6ea498",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,cmd/bench/tree_scope.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,all,passed,322\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,fail,16945\nfailures[1]{package,test,line}:\n  github.com/gibbonmi/bench/cmd/bench,TestExemptTreeRowFollowsExitRule,\"tree_scope_test.go:184: exempt exit 2 = (2, ..."
          },
          "requirement": "repair-probe-R2-exempt-exit2",
          "command": "bench probe cmd/bench/tree_scope.go --swap 'fmt.Fprint(c.Stderr, shownRow(row, exit))' --with 'fmt.Fprint(c.Stderr, row)' --package ./cmd/bench",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/tt-t3-repair-1-20260929/r2-exempt-exit2-probe@744ef656",
              "digest": "sha256:9160321bccd66ade4d48f658c6fc870a7ae5f9b2ad4041159bd034200c6ea498",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,cmd/bench/tree_scope.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,all,passed,322\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,fail,16945\nfailures[1]{package,test,line}:\n  github.com/gibbonmi/bench/cmd/bench,TestExemptTreeRowFollowsExitRule,\"tree_scope_test.go:184: exempt exit 2 = (2, ..."
            }
          }
        },
        {
          "id": "tt-c3-repair-probe-r2-exit1-control-r2",
          "performer": "claude:bench-writer/tt-t3-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "564566b25b32a0cd36399364a16fd537671887c9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t3-repair-1-20260929/r2-exit1-control-probe@744ef656",
            "digest": "sha256:9dc5d0cc7ac9f78be47eada48f1adf12f66db9e06434e84e381ece63a33513f8",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,cmd/bench/tree_scope.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestExemptTreeRowFollowsExitRule,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,fail,17\nfailures[1]{package,test,line}:\n  github.com/gibbonmi/bench/cmd/bench,TestExemptTreeRowFollowsExitRule,\"tree_scope_test.go:180: exempt exit 1 = (1, ..."
          },
          "requirement": "repair-probe-R2-exit1-control",
          "command": "bench probe cmd/bench/tree_scope.go --swap 'if exit == 2 {' --with 'if exit != 0 {' --package ./cmd/bench --run TestExemptTreeRowFollowsExitRule",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/tt-t3-repair-1-20260929/r2-exit1-control-probe@744ef656",
              "digest": "sha256:9dc5d0cc7ac9f78be47eada48f1adf12f66db9e06434e84e381ece63a33513f8",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,cmd/bench/tree_scope.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestExemptTreeRowFollowsExitRule,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,fail,17\nfailures[1]{package,test,line}:\n  github.com/gibbonmi/bench/cmd/bench,TestExemptTreeRowFollowsExitRule,\"tree_scope_test.go:180: exempt exit 1 = (1, ..."
            }
          }
        },
        {
          "id": "tt-c3-repair-probe-r3-bounded-after-r2",
          "performer": "claude:bench-writer/tt-t3-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "564566b25b32a0cd36399364a16fd537671887c9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t3-repair-1-20260929/r3-bounded-after-probe@744ef656",
            "digest": "sha256:3692c62f5111c79f344103b9089a9caf16e89fed387118113d3de3fab9147691",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,cmd/bench/census_output.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,all,passed,322\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,fail,15134\nfailures[1]{package,test,line}:\n  github.com/gibbonmi/bench/cmd/bench,TestTreeRowPrecedesVerb,\"tree_scope_test.go:197: bounded dirtying verb = (0, ..."
          },
          "requirement": "repair-probe-R3-bounded-after",
          "command": "bench probe cmd/bench/census_output.go --swap 'owner.Lead(shownRow(row, exit))' --with 'owner.Lead(shownRow(definition.treeRow(args), exit))' --package ./cmd/bench",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/tt-t3-repair-1-20260929/r3-bounded-after-probe@744ef656",
              "digest": "sha256:3692c62f5111c79f344103b9089a9caf16e89fed387118113d3de3fab9147691",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,cmd/bench/census_output.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,all,passed,322\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,fail,15134\nfailures[1]{package,test,line}:\n  github.com/gibbonmi/bench/cmd/bench,TestTreeRowPrecedesVerb,\"tree_scope_test.go:197: bounded dirtying verb = (0, ..."
            }
          }
        },
        {
          "id": "tt-c3-repair-probe-r3-exempt-after-r2",
          "performer": "claude:bench-writer/tt-t3-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "564566b25b32a0cd36399364a16fd537671887c9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t3-repair-1-20260929/r3-exempt-after-probe@744ef656",
            "digest": "sha256:dff5feaef5f3854751bcf06a32a8f3b5a7f4388ad0ed008d7754c8243561ee9d",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,cmd/bench/command_registry.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestTreeRowPrecedesVerb,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,fail,37\nfailures[1]{package,test,line}:\n  github.com/gibbonmi/bench/cmd/bench,TestTreeRowPrecedesVerb,\"tree_scope_test.go:201: exempt dirtying verb = (0, ..."
          },
          "requirement": "repair-probe-R3-exempt-after",
          "command": "bench probe cmd/bench/command_registry.go --swap 'return c.finishExempt(row, definition.run(c, args[1:]))' --with 'exit := definition.run(c, args[1:]); return c.finishExempt(definition.treeRow(args[1:]), exit)' --package ./cmd/bench --run TestTreeRowPrecedesVerb",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/tt-t3-repair-1-20260929/r3-exempt-after-probe@744ef656",
              "digest": "sha256:dff5feaef5f3854751bcf06a32a8f3b5a7f4388ad0ed008d7754c8243561ee9d",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,cmd/bench/command_registry.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestTreeRowPrecedesVerb,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,fail,37\nfailures[1]{package,test,line}:\n  github.com/gibbonmi/bench/cmd/bench,TestTreeRowPrecedesVerb,\"tree_scope_test.go:201: exempt dirtying verb = (0, ..."
            }
          }
        },
        {
          "id": "tt-c3-3-cmd-r3",
          "performer": "claude:bench-writer/tt-t3-repair-2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ef1d5964df6ea99a501c802627517db74e027f95",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t3-repair-2-20260929/3-cmd@9b7df67e",
            "digest": "sha256:c1a294d0257822aeee7382426f5b2267e090252da52ea9e2d84e30be5926892f",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,pass,12450\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "3-cmd",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "tt-c3-3-treetarget-r3",
          "performer": "claude:bench-writer/tt-t3-repair-2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ef1d5964df6ea99a501c802627517db74e027f95",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t3-repair-2-20260929/3-treetarget@9b7df67e",
            "digest": "sha256:ad8dface2c2f05e8e9847a516cf5ff97b9492d83655aa432f2df7660499e5cac",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/treetarget,pass,117\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "3-treetarget",
          "command": "bench test --package ./internal/treetarget",
          "exit_code": 0
        },
        {
          "id": "tt-c3-3-responsebound-r3",
          "performer": "claude:bench-writer/tt-t3-repair-2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ef1d5964df6ea99a501c802627517db74e027f95",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t3-repair-2-20260929/3-responsebound@9b7df67e",
            "digest": "sha256:465867c960aad630941cd003c829bae589650660753c6fd040cf8701ffba37ee",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/responsebound,pass,250\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "3-responsebound",
          "command": "bench test --package ./internal/responsebound",
          "exit_code": 0
        },
        {
          "id": "tt-c3-3-unassigned-probe-r3",
          "performer": "claude:bench-writer/tt-t3-repair-2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ef1d5964df6ea99a501c802627517db74e027f95",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t3-repair-2-20260929/3-unassigned-probe@9b7df67e",
            "digest": "sha256:68800e4362461085f809a9446c65dd8fb5eab15a2e057df6eb7269a08d0946fd",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/treetarget/identify.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/treetarget,TestIdentifyUnownedWorktree,passed,3\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/treetarget,fail,46"
          },
          "requirement": "3-unassigned-probe",
          "command": "bench probe internal/treetarget/identify.go --swap '\"unassigned\"' --with '\"primary\"' --package ./internal/treetarget --run TestIdentifyUnownedWorktree",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/tt-t3-repair-2-20260929/3-unassigned-probe@9b7df67e",
              "digest": "sha256:68800e4362461085f809a9446c65dd8fb5eab15a2e057df6eb7269a08d0946fd",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/treetarget/identify.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/treetarget,TestIdentifyUnownedWorktree,passed,3\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/treetarget,fail,46"
            }
          }
        },
        {
          "id": "tt-c3-repair-probe-r4-strip-label-r3",
          "performer": "claude:bench-writer/tt-t3-repair-2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ef1d5964df6ea99a501c802627517db74e027f95",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t3-repair-2-20260929/r4-strip-label-probe@9b7df67e",
            "digest": "sha256:9b3e28bdd384ebeed4861f360cc2ca83c701ad470ad63a07a01bb1e7c9e9795e",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/treetarget/identify.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/treetarget,TestRowStripsHostileLabel,passed,3\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/treetarget,fail,3\nfailures[2]{package,test,line}:\n  github.com/gibbonmi/bench/internal/treetarget,TestRowStripsHostileLabel/hostile,\"identify_test.go:93: label ...\"\n  github.com/gibbonmi/bench/internal/treetarget,TestRowStripsHostileLabel/literal,\"identify_test.go:93: label ...\""
          },
          "requirement": "repair-probe-R4-strip-label",
          "command": "bench probe internal/treetarget/identify.go --swap 'sanitize.Strip(identity.Target)' --with 'sanitize.Controls(identity.Target)' --package ./internal/treetarget --run TestRowStripsHostileLabel",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/tt-t3-repair-2-20260929/r4-strip-label-probe@9b7df67e",
              "digest": "sha256:9b3e28bdd384ebeed4861f360cc2ca83c701ad470ad63a07a01bb1e7c9e9795e",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/treetarget/identify.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/treetarget,TestRowStripsHostileLabel,passed,3\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/treetarget,fail,3\nfailures[2]{package,test,line}:\n  github.com/gibbonmi/bench/internal/treetarget,TestRowStripsHostileLabel/hostile,\"identify_test.go:93: label ...\"\n  github.com/gibbonmi/bench/internal/treetarget,TestRowStripsHostileLabel/literal,\"identify_test.go:93: label ...\""
            }
          }
        },
        {
          "id": "tt-c3-repair-probe-r5-shared-harness-r3",
          "performer": "claude:bench-writer/tt-t3-repair-2",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ef1d5964df6ea99a501c802627517db74e027f95",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-t3-repair-2-20260929/r5-shared-harness-probe@9b7df67e",
            "digest": "sha256:8ebfa700b15ebb20408f887b3bcd55392c5826f5e77cbbd9558ee70072d88a06",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,cmd/bench/response_bound_test.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestDispatcherPassesExemptResponse|TestExemptTreeRowFollowsExitRule|TestTreeRowPrecedesVerb,passed,3\nfailures[3]{package,test,line}:\n  github.com/gibbonmi/bench/cmd/bench,TestDispatcherPassesExemptResponse,\"response_bound_test.go:126: exempt fixture stdout = ...\"\n  github.com/gibbonmi/bench/cmd/bench,TestExemptTreeRowFollowsExitRule,\"tree_scope_test.go:166: exempt exit 1 = (1, ...\"\n  github.com/gibbonmi/bench/cmd/bench,TestTreeRowPrecedesVerb,\"tree_scope_test.go:187: exempt dirtying verb = (0, ...\""
          },
          "requirement": "repair-probe-R5-shared-harness",
          "command": "bench probe cmd/bench/response_bound_test.go --swap 'Bound:     bound,' --with 'Bound:     boundResponse,' --package ./cmd/bench --run 'TestDispatcherPassesExemptResponse|TestExemptTreeRowFollowsExitRule|TestTreeRowPrecedesVerb'",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/tt-t3-repair-2-20260929/r5-shared-harness-probe@9b7df67e",
              "digest": "sha256:8ebfa700b15ebb20408f887b3bcd55392c5826f5e77cbbd9558ee70072d88a06",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,cmd/bench/response_bound_test.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestDispatcherPassesExemptResponse|TestExemptTreeRowFollowsExitRule|TestTreeRowPrecedesVerb,passed,3\nfailures[3]{package,test,line}:\n  github.com/gibbonmi/bench/cmd/bench,TestDispatcherPassesExemptResponse,\"response_bound_test.go:126: exempt fixture stdout = ...\"\n  github.com/gibbonmi/bench/cmd/bench,TestExemptTreeRowFollowsExitRule,\"tree_scope_test.go:166: exempt exit 1 = (1, ...\"\n  github.com/gibbonmi/bench/cmd/bench,TestTreeRowPrecedesVerb,\"tree_scope_test.go:187: exempt dirtying verb = (0, ...\""
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "tt-c3-standards-r1",
          "performer": "claude:bench-reviewer/tt-c3-standards-r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "8bf4434336ae7fb7c412c3cde6930803bbe93b01",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-c3-standards-r1@5444f36a",
            "digest": "sha256:32cdaa396edc864c7068c324872251161c1d2dbc21308b3305b6250d96b7b904",
            "excerpt": "Standards: 0 findings. No candidate survived refutation."
          },
          "axis": "Standards",
          "base": "74daf800f8a11d85083711b31aaabc47797ebc73",
          "tip": "5444f36ac57e3310026b62439bcf7648d49871df",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "tt-c3-spec-r1",
          "performer": "claude:bench-reviewer/tt-c3-spec-r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "8bf4434336ae7fb7c412c3cde6930803bbe93b01",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/tt-c3-spec-r1@5444f36a",
            "digest": "sha256:13bac55949b81ea1e0f4437864611c99d03c8b5abc1cff0fb3eb02279b8d8ca0",
            "excerpt": "Spec: 1 finding. Worst: a ledger label that TOON cannot carry silently drops the identity row."
          },
          "axis": "Spec",
          "base": "74daf800f8a11d85083711b31aaabc47797ebc73",
          "tip": "5444f36ac57e3310026b62439bcf7648d49871df",
          "finding_ids": [
            "R1"
          ],
          "supersedes": []
        },
        {
          "id": "tt-c3-coverage-r1",
          "performer": "claude:bench-reviewer/tt-c3-coverage-r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "8bf4434336ae7fb7c412c3cde6930803bbe93b01",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/tt-c3-coverage-r1@5444f36a",
            "digest": "sha256:c0be25e82cf536b1609a854f06d8770ed52fb4eab45340c583cd3f9feb3bd9c1",
            "excerpt": "Coverage: 2 findings. Worst: two binding spec rules about when the row is printed have no test, and two mutations that break them pass the whole ./cmd/bench suite."
          },
          "axis": "Coverage",
          "base": "74daf800f8a11d85083711b31aaabc47797ebc73",
          "tip": "5444f36ac57e3310026b62439bcf7648d49871df",
          "finding_ids": [
            "R2",
            "R3"
          ],
          "supersedes": []
        },
        {
          "id": "tt-c3-standards-r2",
          "performer": "claude:bench-reviewer/tt-c3-standards-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "564566b25b32a0cd36399364a16fd537671887c9",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/tt-c3-standards-r2@7251523e",
            "digest": "sha256:0a8f16ad1bde8930661c043809d66d48c45c70a37864e7a0fdc6f6bf85c02e49",
            "excerpt": "Standards: 2 findings. Worst: the escape is conditional, so a hostile label and a plain label can print the same row."
          },
          "axis": "Standards",
          "base": "74daf800f8a11d85083711b31aaabc47797ebc73",
          "tip": "7251523e83a6993521cbd51db6531daf2297d9b9",
          "finding_ids": [
            "R4",
            "R5"
          ],
          "supersedes": [
            "tt-c3-standards-r1"
          ]
        },
        {
          "id": "tt-c3-spec-r2",
          "performer": "claude:bench-reviewer/tt-c3-spec-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "564566b25b32a0cd36399364a16fd537671887c9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-c3-spec-r2@7251523e",
            "digest": "sha256:9ea04d04426becf63cc9f6918508ec0849c26b48cef82cfabbc964e7c4fd3311",
            "excerpt": "Spec: 0 findings. All three folds are confirmed against the repair delta fe44447f..7251523e."
          },
          "axis": "Spec",
          "base": "74daf800f8a11d85083711b31aaabc47797ebc73",
          "tip": "7251523e83a6993521cbd51db6531daf2297d9b9",
          "finding_ids": [],
          "supersedes": [
            "tt-c3-spec-r1"
          ]
        },
        {
          "id": "tt-c3-coverage-r2",
          "performer": "claude:bench-reviewer/tt-c3-coverage-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "564566b25b32a0cd36399364a16fd537671887c9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-c3-coverage-r2@7251523e",
            "digest": "sha256:53f5c0d9920812caf87f92aefeafd2e4679f173f73e0a653a8c20de84882b60d",
            "excerpt": "Coverage: 0 findings. All three folds are pinned, and an independent R2 bypass is caught."
          },
          "axis": "Coverage",
          "base": "74daf800f8a11d85083711b31aaabc47797ebc73",
          "tip": "7251523e83a6993521cbd51db6531daf2297d9b9",
          "finding_ids": [],
          "supersedes": [
            "tt-c3-coverage-r1"
          ]
        },
        {
          "id": "tt-c3-standards-r3",
          "performer": "claude:bench-reviewer/tt-c3-standards-r3",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "ef1d5964df6ea99a501c802627517db74e027f95",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-c3-standards-r3@9e742c54",
            "digest": "sha256:7863fde2cc7954de942cd6c3b4d04a59d8b3cde5e6b6bd8a1a401ce6bf450689",
            "excerpt": "Standards: 0 findings. No candidate in repair delta 92e86104..9e742c54 survived refutation."
          },
          "axis": "Standards",
          "base": "74daf800f8a11d85083711b31aaabc47797ebc73",
          "tip": "9e742c5434c327e7174f85dd427331004b2bcef6",
          "finding_ids": [],
          "supersedes": [
            "tt-c3-standards-r2"
          ]
        },
        {
          "id": "tt-c3-spec-r3",
          "performer": "claude:bench-reviewer/tt-c3-spec-r3",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "ef1d5964df6ea99a501c802627517db74e027f95",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-c3-spec-r3@9e742c54",
            "digest": "sha256:0d72e169f72ca467a8e2d1eb74d0b3e35fba86359bfa9893555dd01c8aad4c8a",
            "excerpt": "Spec: 0 findings. R4 and R5 are confirmed, and the rows of this ticket stay met."
          },
          "axis": "Spec",
          "base": "74daf800f8a11d85083711b31aaabc47797ebc73",
          "tip": "9e742c5434c327e7174f85dd427331004b2bcef6",
          "finding_ids": [],
          "supersedes": [
            "tt-c3-spec-r2"
          ]
        },
        {
          "id": "tt-c3-coverage-r3",
          "performer": "claude:bench-reviewer/tt-c3-coverage-r3",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "ef1d5964df6ea99a501c802627517db74e027f95",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/tt-c3-coverage-r3@9e742c54",
            "digest": "sha256:64e9cd5d5624348bc4eb3ff69f5b32caa9458f92b3e117250071658bd0aa1454",
            "excerpt": "Coverage: 0 findings. Both folds hold, and the only surviving bypass narrows a test sample of a requirement the code already meets."
          },
          "axis": "Coverage",
          "base": "74daf800f8a11d85083711b31aaabc47797ebc73",
          "tip": "9e742c5434c327e7174f85dd427331004b2bcef6",
          "finding_ids": [],
          "supersedes": [
            "tt-c3-coverage-r2"
          ]
        }
      ]
    },
    {
      "id": "TT-C4",
      "base": "9e742c5434c327e7174f85dd427331004b2bcef6",
      "tip": "3831024c361eab0bad94374b8fb864bd85a76c4e",
      "plan_digest": "sha256:8e3a09402b758dab396e7e37828fde90bc8eee75ecd288aab9535bffd0c92886",
      "source_digest": "4cd183192b9630a45f19387992ba0e1c68cab351",
      "acceptance_rows": [
        "TT8",
        "TT9",
        "TT26",
        "TT27",
        "TT28",
        "TT29",
        "TT30",
        "TT31",
        "TT32",
        "TT33",
        "TT34",
        "TT35",
        "TT36",
        "TT37",
        "TT38",
        "TT39",
        "TT40",
        "TT41",
        "TT49",
        "TT53",
        "TT54",
        "TT55",
        "TT57",
        "TT58",
        "TT59",
        "TT62",
        "TT43",
        "TT44",
        "TT45",
        "TT46",
        "TT47",
        "TT48",
        "TT50"
      ],
      "verification": [],
      "reviews": [
        {
          "id": "tt-c4-standards-r1",
          "performer": "claude:bench-reviewer/tt-c4-standards-r1",
          "role": "independent-review",
          "model": "fable",
          "effort": "high",
          "source_digest": "4cd183192b9630a45f19387992ba0e1c68cab351",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/tt-c4-standards-r1@3831024c",
            "digest": "sha256:f6f5608cf246b6418c9010f248025b56dd199cf31f72b280745c787bc1e3252a",
            "excerpt": "Standards: 4 findings. Worst: ticket 5's build refusal re-derives the worktree target refusal printer and the `bench worktree build` verb spelling across the fence."
          },
          "axis": "Standards",
          "base": "9e742c5434c327e7174f85dd427331004b2bcef6",
          "tip": "3831024c361eab0bad94374b8fb864bd85a76c4e",
          "finding_ids": [
            "C4-S1",
            "C4-S2",
            "C4-S3",
            "C4-S4"
          ],
          "supersedes": []
        },
        {
          "id": "tt-c4-spec-r1",
          "performer": "claude:bench-reviewer/tt-c4-spec-r1",
          "role": "independent-review",
          "model": "fable",
          "effort": "high",
          "source_digest": "4cd183192b9630a45f19387992ba0e1c68cab351",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/tt-c4-spec-r1@3831024c",
            "digest": "sha256:4af61a40043681b2e37ca3754df54824ba3325c0b46203b5c1a24d4f492ebcb2",
            "excerpt": "Spec: 2 findings. Worst: the spec contradicts itself on the label quoting function, and the build followed the acceptance rows (TT44/TT47) over the named symbol and the import list."
          },
          "axis": "Spec",
          "base": "9e742c5434c327e7174f85dd427331004b2bcef6",
          "tip": "3831024c361eab0bad94374b8fb864bd85a76c4e",
          "finding_ids": [
            "C4-P1",
            "C4-P2"
          ],
          "supersedes": []
        },
        {
          "id": "tt-c4-coverage-r1",
          "performer": "claude:bench-reviewer/tt-c4-coverage-r1",
          "role": "independent-review",
          "model": "fable",
          "effort": "high",
          "source_digest": "4cd183192b9630a45f19387992ba0e1c68cab351",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/tt-c4-coverage-r1@3831024c",
            "digest": "sha256:a8ad58eb9660360ab65060fc3c71efaea9b9e71f559a960f2fc9a73ad3dcfd48",
            "excerpt": "Coverage: 3 findings. Worst: no test pins the missing-tree or creation-bundle refusal on the `--in` path, so a lookup that drops both checks stays green in both packages."
          },
          "axis": "Coverage",
          "base": "9e742c5434c327e7174f85dd427331004b2bcef6",
          "tip": "3831024c361eab0bad94374b8fb864bd85a76c4e",
          "finding_ids": [
            "C4-C1",
            "C4-C2",
            "C4-C3"
          ],
          "supersedes": []
        }
      ]
    }
  ],
  "completion": {
    "state": "pending",
    "source_digest": "",
    "performer": "",
    "reconciliation": {},
    "verification": []
  },
  "amendments": [
    {
      "from": "sha256:82131a586bf62c0c3278a5dd81e06c3709beb7d6f263935fa90651ae0c7eccc8",
      "to": "sha256:ab6c7f85535262e57e3235e5da67d2460d1b56586d98b7ce78c6a487073539fb",
      "chunk_ids": {
        "TT-C1": [
          "TT-C1"
        ],
        "TT-C2": [
          "TT-C2"
        ],
        "TT-C3": [
          "TT-C3"
        ],
        "TT-C4": [
          "TT-C4"
        ]
      }
    },
    {
      "from": "sha256:ab6c7f85535262e57e3235e5da67d2460d1b56586d98b7ce78c6a487073539fb",
      "to": "sha256:834182d7bb87b440ee0493ba5fd18fb70360c237bae9e7a24acfe73bc8fc91a0",
      "chunk_ids": {
        "TT-C1": [
          "TT-C1"
        ],
        "TT-C2": [
          "TT-C2"
        ],
        "TT-C3": [
          "TT-C3"
        ],
        "TT-C4": [
          "TT-C4"
        ]
      }
    },
    {
      "from": "sha256:834182d7bb87b440ee0493ba5fd18fb70360c237bae9e7a24acfe73bc8fc91a0",
      "to": "sha256:ea3ee9d0bf9cd3b6386e8aff692c80c514b0ee1075b723a9256ba5a30f8d0430",
      "chunk_ids": {
        "TT-C1": [
          "TT-C1"
        ],
        "TT-C2": [
          "TT-C2"
        ],
        "TT-C3": [
          "TT-C3"
        ],
        "TT-C4": [
          "TT-C4"
        ]
      }
    },
    {
      "from": "sha256:ea3ee9d0bf9cd3b6386e8aff692c80c514b0ee1075b723a9256ba5a30f8d0430",
      "to": "sha256:51bdf6fe67e2a33a0106d4f13992621fbee9b36369e2c1c19dd499a8582856e3",
      "chunk_ids": {
        "TT-C1": [
          "TT-C1"
        ],
        "TT-C2": [
          "TT-C2"
        ],
        "TT-C3": [
          "TT-C3"
        ],
        "TT-C4": [
          "TT-C4"
        ]
      }
    },
    {
      "from": "sha256:51bdf6fe67e2a33a0106d4f13992621fbee9b36369e2c1c19dd499a8582856e3",
      "to": "sha256:49d7c3455ec98c0ab034c4bb473f14b0515ff390730eefb988e644e95f92b708",
      "chunk_ids": {
        "TT-C1": [
          "TT-C1"
        ],
        "TT-C2": [
          "TT-C2"
        ],
        "TT-C3": [
          "TT-C3"
        ],
        "TT-C4": [
          "TT-C4"
        ]
      }
    },
    {
      "from": "sha256:49d7c3455ec98c0ab034c4bb473f14b0515ff390730eefb988e644e95f92b708",
      "to": "sha256:8e3a09402b758dab396e7e37828fde90bc8eee75ecd288aab9535bffd0c92886",
      "chunk_ids": {
        "TT-C1": [
          "TT-C1"
        ],
        "TT-C2": [
          "TT-C2"
        ],
        "TT-C3": [
          "TT-C3"
        ],
        "TT-C4": [
          "TT-C4"
        ]
      }
    }
  ]
}
```
