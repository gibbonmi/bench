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

```bench-review-record
{
  "version": 2,
  "spec": "specs/tree-targets/spec.md",
  "plan_digest": "sha256:ab6c7f85535262e57e3235e5da67d2460d1b56586d98b7ce78c6a487073539fb",
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
      "tip": "c8b9bb94eb15029a84abd137e1fbac93e8990480",
      "plan_digest": "sha256:ab6c7f85535262e57e3235e5da67d2460d1b56586d98b7ce78c6a487073539fb",
      "source_digest": "673171b5be48cb53a7038841dc6fe6cd56c7a727",
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
            "digest": "sha256:49a9464dd2afdab3deac20382c8d92a0f1a47e5781b9adfee0bd7f4f12ce2874",
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
            "digest": "sha256:ed2820500b31123ca67f43cfb395446dde80693a6108a241afa962db3c618007",
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
              "digest": "sha256:ed2820500b31123ca67f43cfb395446dde80693a6108a241afa962db3c618007",
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
            "digest": "sha256:f3a3e2bedf569ef2dfc157bd1bc3a3c2ef701515c0e6dcfb45f27cd0b6b11b78",
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
              "digest": "sha256:f3a3e2bedf569ef2dfc157bd1bc3a3c2ef701515c0e6dcfb45f27cd0b6b11b78",
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
            "digest": "sha256:f6c1775d1d958a5cc36da568a1225d9386b1a739e497dca559da93e0af2dabdc",
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
              "digest": "sha256:f6c1775d1d958a5cc36da568a1225d9386b1a739e497dca559da93e0af2dabdc",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/kit_source.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestKitSourceCheckout,passed,2\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,fail,4"
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
    }
  ]
}
```
