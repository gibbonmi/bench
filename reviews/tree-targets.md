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

```bench-review-record
{
  "version": 2,
  "spec": "specs/tree-targets/spec.md",
  "plan_digest": "sha256:49d7c3455ec98c0ab034c4bb473f14b0515ff390730eefb988e644e95f92b708",
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
    }
  ]
}
```
