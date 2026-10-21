# Share guard command projection and envelope facts

Status: staged
Decision source: `decisions/architecture-guards.md`, answers 1 and 3, reviewed 2026-10-06
Subject: `f07508b2add49469d686614c70f575531e8ea445` (refresh of 2026-10-20)
Audience: every repository that links the kit
Related roadmap: FT366
Verification log: pending — refresh review at f07508b2a

## Problem

Git and Bench guards share `shellcommand.Parse`, but each guard recognizes shell wrappers itself.
The Git guard also projects worktree children in `internal/gitguard/scan.go` (`scanWords`, `worktreeExecChild`).
The Bench guard recognizes the same exec head in `internal/benchguard/benchguard.go` (`isExecHead`).
The reviewer-verb rule in `internal/benchguard/reviewer.go` repeats the wrapper, exec-head, and `--` child grammar a third time.
A conformance check, `runsBenchUnderExec`, repeats the exec-child grammar a fourth time.
The envelope readers repeat JSON structure in `gitguard`, `benchguard`, and `internal/writeguard/writeguard.go` (`PathFromEnvelope`).
These copies can drift, while each guard's different policy must stay intact.

## Solution

Make `internal/shellcommand` own command projection and common envelope facts.
Keep each guard's authority, required fields, diagnostics, traversal depth, and failure policy.
Preserve the supported command language through actual hook and core entries.
Keep the independent degraded shell rim, because it must run without the core.

## User stories

Line: opus / high
Implementation-line reason: the hardest material chunk is GP2, which moves three decoder forms with fifteen ordered decoder rows and the reviewer-verb walk onto new APIs at once. The spec fixes each decoder target and caller policy exactly, and the seams are resolved owner-package and entry tests. Test strength is high, but the ownership check and the actual-hook witnesses need careful integration.
Harder chunks: GP2 and GP3, at xhigh effort on the same model.

1. As an agent, I want routine executable prefixes preserved, so that ordinary commands receive the same guard verdict.
2. As an agent, I want one wrapper level recognized, so that supported shell children remain visible.
3. As an agent, I want worktree child argv preserved, so that direct child commands retain their meaning.
4. As an agent, I want child checkout context preserved, so that parent repository facts cannot authorize child writes.
5. As an agent, I want Git option policy retained, so that repository redirection remains a Git decision.
6. As an agent, I want attached separators recognized, so that later dangerous commands cannot disappear.
7. As an agent, I want malformed dangerous commands refused, so that lexer failure cannot permit a partial command.
8. As an agent, I want parse-repair advice retained, so that I can correct malformed input.
9. As an agent, I want each guard's required fields retained, so that structural sharing cannot broaden authority.
10. As an agent, I want Git's unreadable-envelope policy retained, so that the core still allows unreadable input silently.
11. As an agent, I want Bench's unreadable-envelope policy retained, so that the warning remains useful.
12. As an agent, I want file-write path validation retained, so that relative paths use the correct cwd.
13. As an agent, I want file-write uncertainty retained, so that unreadable repository facts do not create a new refusal.
14. As an agent, I want valid recovery commands retained, so that I can repair a missing core.
15. As an agent, I want missing-library handling retained, so that each hook keeps its own recovery posture.
16. As an agent, I want missing-core handling retained, so that the independent rim still protects Git operations.
17. As an agent, I want core-error handling retained, so that hooks keep their established diagnostics.
18. As an agent, I want alias recognition retained, so that Bench resolver policy remains consistent.
19. As an agent, I want the exec exception retained, so that worktree child operations remain usable.
20. As an agent, I want pool-reference precedence retained, so that a competing refusal names the pool route.
21. As a maintainer, I want private grammar copies rejected, so that a migration cannot leave a second owner.
22. As a maintainer, I want fixture bindings retained, so that source movement cannot amputate existing checks.
23. As a maintainer, I want deterministic behavior comparison, so that unseen changes block the migration.
24. As an agent, I want census exclusion retained, so that a Bench invocation does not become a raw shell record.
25. As an agent, I want commands parsed without execution, so that guard inspection cannot run shell input.
26. As an agent, I want rg operand uncertainty kept separate, so that permissive operand analysis cannot weaken dangerous-command refusal.
27. As a reviewer, I want the reviewer-verb refusal to keep its recursive reach, so that a nested wrapper or a nested exec child cannot run my verb.
28. As an agent, I want the reviewer-verb refusal to stay between the pool refusal and the follow-on verdict, so that a competing refusal names the same cause as before.
29. As a maintainer, I want the agent-route conformance check to read exec children through the shared projection, so that no test keeps a second exec grammar.

## Implementation decisions

Reviewer decision (conversation, 2026-10-20): `benchguard.ReviewerVerb` keeps its recursive walk as guard policy. It consumes the shared one-level projection facts and recurses through them itself.

### Projection API and caller policy

Keep `Parse`, `ProjectCommandWords`, and `ResolveRoutinePrefix` unchanged.
They remain the executable-word owner and the token stream owner.
Add `projection.go` with these concrete value APIs:

```go
type WrapperChild struct {
    ArgumentIndex int
    Stream Stream
}
func ShellChildren(words []string, executable int) []WrapperChild

type ExecProjection struct {
    Head bool
    HasSeparator bool
    Child []string
    DifferentCheckout bool
}
func WorktreeExec(words []string, executable int) ExecProjection
```

`ShellChildren` recognizes basename `sh`, `bash`, or `zsh` and the existing `^-[A-Za-z]*c[A-Za-z]*$` flag shape.
It returns matching flag children in argv order, parsed by `Parse`.
It performs no recursion.
An invalid executable index or a missing flag operand yields no child at that position.
GP02, GP03, and GP04 observe these cases.

`WorktreeExec` receives an executable index whose identity the caller has already accepted.
It recognizes adjacent `worktree` and `exec` words and finds the first later `--`.
`Child` contains the following argv unchanged, without string joining or reparsing.
`HasSeparator` records that delimiter even when the child is empty.
`DifferentCheckout` is true exactly when the projected child is nonempty.
GP05 through GP08 observe these cases.

Git retains its shell-keyword trimming before `ResolveRoutinePrefix`.
Git consumes the first matching shell child, including a child that finds no Git invocation.
Bench `Classify` retains its existing search across matching children for a Bench invocation.
Git and Bench `Classify` limit traversal to one level across shell and worktree child forms.
GP03, GP04, GP09, and GP74 prevent the shared facts from replacing those policy choices.
Each caller applies its own Bench identity test before it uses `WorktreeExec`; GP72 and GP73 give a non-Bench head on the Git and reviewer paths.

`ReviewerVerb` is the one consumer whose traversal is not one level.
It calls `ShellChildren` and `WorktreeExec` at each level and recurses into each child stream or child argv.
An exec head without a separator yields no child and no verb match, as `reviewer.go` decides today.
It keeps `isBench`, `namesVerb`, and the resolver as its own identity policy.
GP62 through GP64 and GP70 observe this walk.

Git retains `FindSubcommand`, global-option parsing, `redirectsRepository`, and verdict classification.
It propagates `DifferentCheckout` into the existing `elsewhere` decision before it consults repository facts.
Bench retains `isBench`, resolver calls, follow-on precedence, `judgeExec`, and the pool-reference guard.
It uses `ExecProjection.Head` only after its own Bench executable recognition.
The follow-on core entry keeps its order: pool reference, then reviewer verb, then `Classify` (`cmd/bench/guards.go:170-178`).
GP10 through GP13, GP65, and GP66 preserve these decisions.

The agent-route conformance check (`runsBenchUnderExec`) keeps its own Bench identity test and its agent-route scope.
It reads the exec child through `WorktreeExec` instead of its own `worktree`, `exec`, and `--` scan.
GP68 preserves its finding, and GP69 grades the migration.

### Envelope API

Add `envelope.go` in `shellcommand`, with private decoded storage and these APIs:

```go
type EnvelopeForm uint8
const (
    CommandStruct EnvelopeForm = iota
    InputMap
    CwdInputMap
)
type Envelope struct { /* private decoded facts */ }
func DecodeEnvelope(data []byte, form EnvelopeForm) (Envelope, error)
func (e Envelope) Command() string
func (e Envelope) Cwd() string
func (e Envelope) InputString(name string) (value string, present bool, err error)
```

A form selects structural decoder semantics, not required fields or guard authority.
The common owner uses encoding/json.Unmarshal with the original anonymous target shapes.
It does not implement a second JSON decoder.
Each call decodes the outer envelope once into the selected shape.
InputString decodes the selected map value as a string, as the current two-stage map readers do.
GP14–GP19 and GP47–GP61 observe these contracts.

| Form | Exact target shape | Current caller | Retained decoder facts |
|---|---|---|---|
| CommandStruct | anonymous struct with tagged tool_input nested struct and tagged command string | gitguard.CommandFromEnvelope | folded root and command keys, nested-struct merging, scalar-null preservation, retained earliest type error |
| InputMap | anonymous struct with tagged tool_input map[string]json.RawMessage | benchguard.CommandFromEnvelope | folded root key, exact input-map keys, ordered map merging, map-null clearing, final raw value validation |
| CwdInputMap | anonymous struct with tagged cwd string and tagged tool_input map[string]json.RawMessage | writeguard.PathFromEnvelope | folded root keys, exact input-map keys, scalar cwd-null preservation, map-null clearing, retained outer type errors |

Keep the existing field names and JSON tags in these anonymous wire targets.
This also preserves encoding/json error context, rather than introducing named-target error text.
Command returns the decoded nested string for CommandStruct.
Cwd returns the decoded scalar for CwdInputMap.
InputString reads exact map keys for InputMap and CwdInputMap.
Callers use only the accessors for their selected form.

Root keys use encoding/json's folded struct-field matching in all three forms.
CommandStruct also folds nested command keys.
The two map forms keep command and file_path lookup exact.
Do not normalize keys or reject duplicates.
GP47, GP48, and GP57–GP59 observe casing at actual caller seams.

Process duplicate members in their original order through the standard decoder.
Repeated tool_input objects merge into the existing nested struct or map.
Later raw map entries replace earlier entries with that exact key.
Null preserves a scalar string or nested struct but clears a map.
Earlier typed-field errors remain errors even when a later value would be valid.
GP49–GP56 and GP60–GP61 distinguish these cases.

Initially absent or null objects produce the selected form's zero values.
An unsupported root or tool_input type returns the standard decoder error.
Unknown sibling fields remain ignored unless the selected target contains that field.
A first null raw string value returns an empty string with present=true and no string-type error.
GP14–GP19 preserve these initial-state facts.

InputString distinguishes missing map keys from invalid final string types.
It does not require a field or reject empty strings and NUL bytes.
The common owner does not resolve paths, inspect repositories, choose an exit, or produce a guard diagnostic.
These decisions remain in the three existing exported reader wrappers.
GP20–GP29 and GP47–GP61 preserve caller contracts.

Git selects CommandStruct and keeps its string-only signature and empty result on any retained decode error.
Bench selects InputMap and keeps its error signature and existing missing, type, empty, and NUL error text.
Write selects CwdInputMap and retains outer cwd validation before path handling, including for absolute paths.
It retains path cleanup and the existing relative-path requirement.
Neither Git nor Bench begins to validate an unrelated cwd field.
GP20–GP29 and the decoder preservation matrix below observe each branch.

A folded lookup plus one final RawMessage map is insufficient for CommandStruct and scalar cwd semantics.
Do not substitute that representation for the selected standard decoder target.
Do not migrate to strict JSON, duplicate rejection, or uniform null handling.
Those changes would alter the approved guard policies.

### Migration and enforcement

| Current site | Shared fact after migration | Policy that stays local | Rows |
|---|---|---|---|
| `gitguard.scanWords` wrapper predicate and flag loop (`scan.go:67-75`) | `ShellChildren` | first matching child and one-level traversal | GP02, GP03, GP09, GP74 |
| `gitguard.worktreeExecChild` (`scan.go:88-98`) | `WorktreeExec` | Bench basename recognition and child authority | GP05–GP10, GP72 |
| `benchguard.judge`, `spanInvokesBench`, `isWrapper`, `wrapperFlag` (`benchguard.go:17, 141-155, 315-332`) | `ShellChildren` | Bench identity and matching-child selection | GP03–GP04, GP11 |
| `benchguard.isExecHead` (`benchguard.go:134, 251`) | `WorktreeExec.Head` | outer exec exception and follow-on order | GP12–GP13 |
| `benchguard.runsReviewerVerb` (`reviewer.go:30-38`) | `ShellChildren` and `WorktreeExec` | recursive walk, verb naming, and Bench identity | GP62–GP64, GP67, GP70, GP73 |
| conformance `runsBenchUnderExec` (`refusal_route_guard_test.go:84-91`) | `WorktreeExec` | agent-route scope and Bench identity | GP68–GP69 |
| Git, Bench, and write envelope readers | `DecodeEnvelope` with the current target form | required fields and failure posture | GP14–GP29, GP47–GP61 |

Retire the private wrapper recognizers and the private exec-child grammar after their callers migrate.
Do not introduce a second executable-prefix helper.
Keep `ResolveRoutinePrefix` as the existing owner.
The compiler and GP29 detect missing migrations.
GP30 and GP67 detect surviving private copies, including renamed unused helpers.

Add `checkGuardProjectionOwner` beneath the existing registered `package-core-guard` check.
It grades production Go files in shellcommand, gitguard, benchguard, and writeguard, through AST import bindings.
Outside shellcommand, it rejects JSON unmarshal or decoder calls that decode these envelopes.
It rejects the migrated shell-name predicate, the wrapper-c flag recognizer, and adjacent worktree-exec child extraction outside the owner.
The checker grades executable AST nodes, including unused function bodies, rather than comments or fixture string data.
GP30 and GP67 require renamed copies of each retired recognizer to make the registered check red.

The subcheck applies only to kit source surfaces that contain conformance and at least one governed guard package.
Absent consumer source surfaces produce no diagnostic.
A present kit surface with a missing or unreadable owner produces a named diagnostic.
Do not require optional source directories in repositories that only link the kit.
GP31 and GP32 prove both dispositions.

Retain the existing registry name, tier, subject, and input binding (`internal/conformance/registry/checks.go:19`).
Do not copy a registry count into this spec or add a new count expectation.
Classify every planned test in `guard_projection_owner_test.go` that reads live tree sources in `tier_live_tree_test.go`, including the GP30 and GP67 mutation tests if they build fixtures from live sources.
GP33 observes registered reachability and that classification.

The `stdlib-reimplementation` check (`internal/conformance/stdlib_reimplementation_test.go:39`) also grades the new production Go.
New shellcommand code must not hand-write a loop that `slices.Contains`, `maps.Copy`, `min`, `max`, `strings.Cut`, or `strings.CutPrefix` replaces.
The `bounds-policy` check pins no file in this scope.

### Conformance orchestration placement

Before you extend the registered check, move only checkPackageCoreAndGuards from package_core_checks_test.go to checks_test.go.
Keep its name, signature, existing subcheck order, and callers unchanged.
Keep the other producers, the npm formatter call, and their tests in package_core_checks_test.go.
Add the new ownership subcheck to the relocated function.

package_core_checks_test.go has 470 lines against the default 400-line limit, and no grant covers it.
The required commit lane runs `structure --growth` (`internal/gate/lane.go:83`).
That ratchet reports FILE GREW for a source file over its limit that also grew (`internal/structure/structure.go:309-331`).
An added orchestration call there would therefore turn the lane red.
checks_test.go has 642 lines against its 739-line grant (`.bench/structure.budgets:51`), which leaves 97 lines.
The moved declaration spans 16 lines (`package_core_checks_test.go:22-37`).
Keep checks_test.go within 739 lines after both guard outcomes, and keep each new file within 400 lines.
Do not add a budget or acceptance grant.

The move changes no reader contract.
`checks_test.go:49` binds package-core-guard by function symbol, so the declaration file does not matter.
`node_runtime_policy_test.go:143` and `:147` call the same package symbol and need no edit.
`package_core_diagnostics_test.go:202-228` counts formatter calls in package_core_checks_test.go; the moved function holds none.
`TestSkillsIndexConformanceCarriesNoSecondReader` reads checks_test.go for copied skills-index literals; the moved function holds none.

Apply this placement only if the declaration remains in the original file when implementation begins.
If the guard-scan-loop outcome already moved it, extend that sole declaration in checks_test.go.
Retain both ownership subchecks if both outcomes have landed.
Either outcome can land first; neither requires the other's unlanded implementation.

### Ordering and independent rim

The degraded-guard-refusal outcome landed before this refresh (commit a9324deb7).
The degraded Git hook now refuses a malformed command with the parse-repair sentence (`.bench/hooks/block-dangerous-git.sh:105-107`).
The core emits the same sentence (`internal/gitguard/gitguard.go:143-146`, `internal/shellcommand/unlexed.go:7`).
`internal/systemtest/guard_refusal_test.go` holds its witnesses, and this spec cites them as preserved regressions.
No prerequisite outcome remains.

Leave the shell resolver and the degraded grammar implementations in place.
They cannot import the Go owner when the core is missing.
Keep the existing classifier-table, resolver-order, and entry-point-parity checks as independent evidence.
Resolver and classifier checks remain beneath package-core-guard.
Entry-point parity keeps its separate entry-point-parity registry binding.
GP36 through GP42 and GP71 preserve the missing-library, missing-core, core-error, stale-core, and recovery cases.

## Implementation chunks

Each chunk carries its source claims, consumer migration, tests, and a green ordinary test phase.
Slicing preserves every acceptance predicate and the implementation line.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
|---|---|---|---|---|
| GP1 / 1-share-git-command-projections.md | Git uses shared wrapper and worktree projections with unchanged verdicts | GP01–GP03, GP05–GP10, GP72, GP74 | projection and Git contract tests | no |
| GP2 / 2-migrate-guard-envelope-consumers.md, 3-migrate-bench-projection-consumers.md | common envelope facts serve the three readers, and Bench, the reviewer verb, and the agent-route check use the projections | GP04, GP11–GP29, GP43, GP47–GP66, GP68–GP70, GP73 | envelope, Bench, reviewer-verb, write, route-check, canary, and core entry tests | yes, xhigh |
| GP3 / 4-enforce-guard-projection-ownership.md, 5-prove-guard-entry-regression.md | copy refusal and actual-entry integration complete the migration | GP30–GP42, GP44–GP46, GP67, GP71 | owner mutation tests and system hook witnesses | yes, xhigh |

GP1 precedes GP2, because the shared projection API must exist before the second consumer migrates.
GP1 retires the Git recognizers in scan.go.
GP2 retires the Bench recognizers in one step, because `reviewer.go` and `benchguard.go` share `isWrapper`, `isExecHead`, and `wrapperFlag`.
GP2 owns GP47–GP61 exactly once and closes every decoder predicate at its first-use green checkpoint.
GP3 repeats those cases through the final differential and actual-entry witnesses under GP46, without duplicate Covers ownership.
GP3 enables its final ownership invariant after all consumer migrations.
The slicer split GP2 into two serial tickets: ticket 2 migrates the envelope readers, and ticket 3 migrates the Bench projection consumers.
Ticket 2 owns GP14–GP29, GP43, and GP47–GP61; ticket 3 owns GP04, GP11–GP13, GP62–GP66, GP68–GP70, and GP73.
Both tickets write `benchguard.go`, `benchguard/projection_contract_test.go`, and `cmd/bench/guard_projection_test.go`, so ticket 3 starts after ticket 2 commits green.

### Graph ordering and execution boundary

The stable old-to-new mapping is GP1 to GP1, GP2 to GP2, GP3 to GP3.
Each chunk keeps its independent Standards, Spec, and Coverage checkpoint.
GP1, GP2, and GP3 run in that order; each successor waits for its predecessor's green commit and review.
The guard-scan-loop outcome (GS1) is independently useful and consumes no unlanded projection interface.

Both outcomes write checks_test.go, package_core_checks_test.go, and tier_live_tree_test.go.
These shared writes require a serial integrated-source handoff between the two outcomes.
The coordinator chooses either complete outcome first and refreshes source and headroom before an overlapping successor charge.
The first ticket that extends package-core-guard performs the orchestration move; the other extends that sole declaration.

Derived Writes scope for the slicer, at f07508b2a:

- GP1: `internal/shellcommand/projection.go` (new), `internal/shellcommand/projection_test.go` (new), `internal/gitguard/scan.go`, `internal/gitguard/projection_contract_test.go` (new).
- GP2: `internal/shellcommand/envelope.go` (new), `internal/shellcommand/envelope_test.go` (new), `internal/gitguard/gitguard.go`, `internal/gitguard/projection_contract_test.go`, `internal/benchguard/benchguard.go`, `internal/benchguard/reviewer.go`, `internal/benchguard/projection_contract_test.go` (new), `internal/writeguard/writeguard.go`, `internal/writeguard/envelope_contract_test.go` (new), `cmd/bench/guard_projection_test.go` (new), `internal/conformance/refusal_route_guard_test.go`, `tests/canary/canonical-path-owner/second-derivation/`.
- GP3: `internal/conformance/guard_projection_owner_test.go` (new), `internal/conformance/package_core_checks_test.go`, `internal/conformance/checks_test.go`, `internal/conformance/tier_live_tree_test.go`, `internal/systemtest/guard_projection_test.go` (new), `internal/systemtest/testdata/guard-projection-family-d.golden` (new), `cmd/bench/guard_projection_test.go`.

GP3 holds two tickets with disjoint writes: ticket 4 owns the registered copy refusal (GP30–GP33, GP42, GP67), and ticket 5 owns the actual-entry regression (GP34–GP41, GP44–GP46, GP71).

The refresh drops seven paths from the 2026-10-06 fence.
Commit 1e8846bb0 deleted the ownership closure that required the five command-registry co-names: `cmd/bench/command_registry.go`, `cmd/bench/command_registry_test.go`, `cmd/bench/help_inventory_test.go`, `internal/conformance/axi_query_registry_test.go`, and `internal/conformance/subcommand_routing_table_test.go`.
No row changes `internal/conformance/registry/checks.go`, because the registry row keeps its metadata.
The injected-ports canary keeps its `type Checker struct` anchor, so `tests/canary/injected-ports/unregistered-port/` needs no edit.
The slicer's research confirms or restores each dropped path.

The version-1 completion plan is a planning verification inventory, not a dispatch plan.
Before any implementation charge, the coordinator binds actual run and session identities and ticket verification ownership in version 2 with author-limit 1.
Each ticket receives a fresh author on opus/high, with GP2 and GP3 at xhigh.
Refresh current owner interfaces, source closure, and actual headroom before that charge.
Build and retain the pre-migration executable for GP46 before GP1 begins.
Staged status and planning validation grant no build authority.

### Completion plan

These commands and mutations are future obligations; none ran during authoring.
Before you retain an independent expectation, pin its mutation source and exact diagnostic and demonstrate a compiling behavioral red.
Require byte-identical restoration and the same focused green; an invalid or restore-failed proof closes no obligation.
Final verification runs only after the complete, independently reviewed outcome.

```bench-completion-plan
{
  "version": 1,
  "chunks": [
    {
      "id": "GP1",
      "tickets": [
        "1-share-git-command-projections.md"
      ],
      "verification": [
        {
          "id": "projection",
          "command": "bench test --package ./internal/shellcommand"
        },
        {
          "id": "git-policy",
          "command": "bench test --package ./internal/gitguard"
        },
        {
          "id": "root-conformance",
          "command": "bench test --check root-conformance"
        },
        {
          "id": "git-context-omission",
          "command": "bench test --package ./internal/gitguard --run '^TestGitProjectionAuthority$'",
          "probe": "Omit only DifferentCheckout propagation at the Git consumer, keeping WorktreeExec and its other facts; child push must red the existing unresolved-destination assertion."
        },
        {
          "id": "wrapper-projection-omission",
          "command": "bench test --package ./internal/shellcommand --run '^TestShellChildren$'",
          "probe": "Omit a supported shell basename or c-flag projection while preserving the API; the independent ordered child cases must red."
        }
      ]
    },
    {
      "id": "GP2",
      "tickets": [
        "2-migrate-guard-envelope-consumers.md",
        "3-migrate-bench-projection-consumers.md"
      ],
      "verification": [
        {
          "id": "owner-values",
          "command": "bench test --package ./internal/shellcommand"
        },
        {
          "id": "git-policy",
          "command": "bench test --package ./internal/gitguard"
        },
        {
          "id": "bench-policy",
          "command": "bench test --package ./internal/benchguard"
        },
        {
          "id": "write-policy",
          "command": "bench test --package ./internal/writeguard"
        },
        {
          "id": "actual-core-entries",
          "command": "bench test --package ./cmd/bench"
        },
        {
          "id": "existing-census",
          "command": "bench test --package ./internal/census"
        },
        {
          "id": "conformance-closure",
          "command": "bench test --package ./internal/conformance"
        },
        {
          "id": "root-conformance",
          "command": "bench test --check root-conformance"
        },
        {
          "id": "canary-bites",
          "command": "bench test --package ./internal/conformance --run '^TestEveryRetainedFixtureBitesThroughRegisteredOwner$'"
        },
        {
          "id": "git-decoder-form-omission",
          "command": "bench test --package ./cmd/bench --run '^TestGuardGitEnvelopePolicy$'",
          "probe": "Replace only the Git caller CommandStruct selection with InputMap; ordered folded/null/early-error cases must red through the actual guardGit entry."
        },
        {
          "id": "bench-decoder-form-omission",
          "command": "bench test --package ./internal/benchguard --run '^TestBenchEnvelope(RequiredField|Value)$'",
          "probe": "Replace only the Bench caller InputMap selection with CommandStruct; exact-map and map-null cases must red."
        },
        {
          "id": "write-decoder-form-omission",
          "command": "bench test --package ./cmd/bench --run '^TestGuardWriteEnvelopePolicy$'",
          "probe": "Replace only the Write caller CwdInputMap selection with InputMap; actual tracked-relative-primary cwd cases must red."
        },
        {
          "id": "reviewer-recursion-omission",
          "command": "bench test --package ./internal/benchguard --run '^TestReviewerVerbKeepsRecursiveReach$'",
          "probe": "Replace only the reviewer walk's recursion into a wrapper child stream with a one-level verb check; the nested-wrapper case must red."
        }
      ]
    },
    {
      "id": "GP3",
      "tickets": [
        "4-enforce-guard-projection-ownership.md",
        "5-prove-guard-entry-regression.md"
      ],
      "verification": [
        {
          "id": "registered-owner",
          "command": "bench test --check package-core-guard"
        },
        {
          "id": "conformance-closure",
          "command": "bench test --package ./internal/conformance"
        },
        {
          "id": "registry-closure",
          "command": "bench test --package ./internal/conformance/registry"
        },
        {
          "id": "root-conformance",
          "command": "bench test --check root-conformance"
        },
        {
          "id": "core-regression",
          "command": "bench test --package ./cmd/bench"
        },
        {
          "id": "actual-hook-family",
          "command": "bench test --check system"
        },
        {
          "id": "copy-check-omission",
          "command": "bench test --package ./internal/conformance --run '^TestGuardProjectionOwnerBitesOnPrivateCopies$'",
          "probe": "Omit only the registered checkGuardProjectionOwner call while retaining its function and tests; renamed unused copies, including the reviewer.go copy, must stop producing the expected registered diagnostic and red this witness."
        },
        {
          "id": "differential-golden-bite",
          "command": "bench test --check system",
          "probe": "Swap the Git caller's first-child selection for an every-child search; the independently captured guard-projection-family-d.golden must red TestGuardProjectionDifferentialFamilies on the W4 member."
        }
      ]
    }
  ],
  "final_verification": [
    {
      "id": "coverage",
      "command": "bench coverage --check shared-guard-projection"
    },
    {
      "id": "shared-owner",
      "command": "bench test --package ./internal/shellcommand"
    },
    {
      "id": "git-policy",
      "command": "bench test --package ./internal/gitguard"
    },
    {
      "id": "bench-policy",
      "command": "bench test --package ./internal/benchguard"
    },
    {
      "id": "write-policy",
      "command": "bench test --package ./internal/writeguard"
    },
    {
      "id": "census",
      "command": "bench test --package ./internal/census"
    },
    {
      "id": "core-entries",
      "command": "bench test --package ./cmd/bench"
    },
    {
      "id": "registered-owner",
      "command": "bench test --check package-core-guard"
    },
    {
      "id": "conformance-closure",
      "command": "bench test --package ./internal/conformance"
    },
    {
      "id": "root-conformance",
      "command": "bench test --check root-conformance"
    },
    {
      "id": "actual-hooks",
      "command": "bench test --check system"
    }
  ]
}
```

## Testing decisions

Use value tests for projections and structural decoding.
Use each guard's existing injected facts for policy tests.
Use direct cmd entry tests for exact exit and warning contracts.
Use the system process owner for actual hook and selected-core witnesses.
Do not introduce a second subprocess harness or change existing test assertions.

Existing precedents are `TestResolveRoutinePrefix`, `TestScanWrapperDepthAndGlobalOpts`, and `TestInvokesBenchWalksOneWrapperLevel` in their owner packages.
`TestGuardDeniesTheOverrideVerb` supplies the reviewer-verb value seam.
`TestPathFromEnvelope` supplies the path-reader seam.
`TestGuardGitBlockAllow` and `TestFollowOnGuardDeniesTheOverrideEnvelope` supply the core entry seams.
`TestDegradedGuardRefusesMalformedCommands`, `TestBenchFollowOnHookDegradedRim`, and `TestDegradedGuardRimDecidesFromTheCommandField` supply actual degraded-entry seams.

### Seam diagram

```text
Bash or file-tool envelope -> caller reader -> shellcommand structural facts
                                              |
command words -> shared projection -> caller authority and verdict -> core exit
                     |                      (reviewer verb recurses here)
actual harness hook --------------------------------------------------> shim exit
```

Tests attach to returned values, caller verdicts, core exits, and actual shim exits.
The registered ordinary conformance test detects private copies in the executed root.

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| GP01 | 1 | Routine-prefix family P retains its projected executable index | `internal/shellcommand/shellcommand_test.go` (`TestResolveRoutinePrefix`) | A new private prefix rule differs on query or option forms |
| GP02 | 2 | ShellChildren returns one ordered child per supported c flag in family W | planned TestShellChildren in internal/shellcommand/projection_test.go | A wrapper basename or flag omission loses a child |
| GP03 | 2 | Invalid indexes and missing wrapper operands yield no child | planned TestShellChildrenBounds in internal/shellcommand/projection_test.go | Unchecked indexes panic or invent a child |
| GP04 | 2 | Bench Classify searches every matching child in W4 and refuses the follow-on in the second `-c` child | planned TestGuardWrapperSelection in internal/benchguard/projection_contract_test.go | A universal first-child or recursive policy changes the Bench verdict |
| GP05 | 3 | WorktreeExec returns child argv byte-for-byte in family X | planned TestWorktreeExec in internal/shellcommand/projection_test.go | Joining and reparsing loses a quoted argv element |
| GP06 | 3 | A recognized exec head without a delimiter has no child | planned TestWorktreeExecMissingDelimiter in internal/shellcommand/projection_test.go | Header recognition is confused with runnable child discovery |
| GP07 | 3 | An empty child after a delimiter has DifferentCheckout=false | planned TestWorktreeExecEmptyChild in internal/shellcommand/projection_test.go | An empty projection fabricates child authority |
| GP08 | 4 | A nonempty projected child has DifferentCheckout=true | planned TestWorktreeExecContext in internal/shellcommand/projection_test.go | An omitted context bit makes parent facts available |
| GP09 | 2, 3 | Nested shell and exec combinations in X3 retain one-level Git verdicts | `internal/gitguard/verdict_test.go` (`TestClassifyVerdicts`) | Recursive expansion changes the documented depth |
| GP10 | 4, 5 | Child and redirected push destinations retain the unresolved-destination label | planned TestGitProjectionAuthority in internal/gitguard/projection_contract_test.go | Parent facts incorrectly permit the child push |
| GP11 | 18 | Bare and symlink Bench aliases retain classification | `internal/benchguard/benchguard_test.go` (`TestClassifyResolvesBareAliasFromProcessPath`) | Shared name matching replaces caller resolver policy |
| GP12 | 19 | Existing outer exec follow-on and heredoc exceptions retain verdicts | `internal/benchguard/benchguard_test.go` (`TestClassifySpanScopedFollowOns`) | Child projection is incorrectly used as a new refusal policy |
| GP13 | 20 | A pool reference competing with a follow-on emits the pool refusal | planned TestGuardProjectionPoolPrecedence in cmd/bench/guard_projection_test.go | Reordered consumers emit the wrong winning diagnostic |
| GP14 | 9 | Invalid JSON returns a structural decode error | planned TestEnvelopeInvalidJSON in internal/shellcommand/envelope_test.go | An unreadable envelope becomes a fabricated empty object |
| GP15 | 9 | Array or scalar root values return a structural decode error | planned TestEnvelopeRootShape in internal/shellcommand/envelope_test.go | A decoder accepts an unsupported root |
| GP16 | 9 | Array or scalar tool_input values return a structural decode error | planned TestEnvelopeInputShape in internal/shellcommand/envelope_test.go | A decoder accepts a non-object input |
| GP17 | 9 | Initially absent and null objects produce each form’s zero values | planned TestEnvelopeAbsentAndNullObjects in internal/shellcommand/envelope_test.go | Structural sharing invents required fields |
| GP18 | 9 | An unrelated invalid cwd type leaves both command forms readable | planned TestEnvelopeIndependentFields in internal/shellcommand/envelope_test.go | Eager sibling validation changes Git and Bench posture |
| GP19 | 9 | A first null input-map string field returns the empty present value | planned TestEnvelopeNullString in internal/shellcommand/envelope_test.go | Null becomes a new type refusal |
| GP20 | 10 | Git initially unreadable, missing, null, and empty command inputs return an empty string | planned TestGitEnvelopeContract in internal/gitguard/projection_contract_test.go | Structural errors become a new required-field policy |
| GP21 | 10 | The Git core emits exit 0 with empty stderr for unreadable-envelope family E0 | planned TestGuardGitEnvelopePolicy in cmd/bench/guard_projection_test.go | A shared reader leaks Bench's warning or refusal |
| GP22 | 11 | Bench missing command retains its current error text | planned TestBenchEnvelopeRequiredField in internal/benchguard/projection_contract_test.go | The accessor accepts an absent required field |
| GP23 | 11 | Bench non-string command retains its current error text | planned TestBenchEnvelopeType in internal/benchguard/projection_contract_test.go | Invalid types silently become command text |
| GP24 | 11 | Bench empty and NUL command values retain their current error text | planned TestBenchEnvelopeValue in internal/benchguard/projection_contract_test.go | Shared structure removes caller validation |
| GP25 | 11 | The Bench core retains warning-allow output for family E0 | planned TestGuardBenchEnvelopePolicy in cmd/bench/guard_projection_test.go | Sharing imports Git's silent allow or a new refusal |
| GP26 | 12 | Absolute and cwd-relative paths retain the current cleaned path result | `internal/writeguard/writeguard_test.go` (`TestPathFromEnvelope`) | Moving decoding changes path interpretation |
| GP27 | 12 | Missing, non-string, empty, and NUL path values remain unreadable | `internal/writeguard/writeguard_test.go` (`TestPathFromEnvelopeRefusesAnUnreadablePath`) | Field access removes caller path validation |
| GP28 | 12, 13 | Invalid cwd types remain unreadable even for an absolute file_path | planned TestWriteEnvelopeCwdContract in internal/writeguard/envelope_contract_test.go | Lazy field validation broadens the write reader |
| GP29 | 13 | File-write core warning-allow output survives unreadable path and fact inputs | planned TestGuardWriteEnvelopePolicy in cmd/bench/guard_projection_test.go | Shared structure changes the write guard's failure posture |
| GP30 | 21 | Each renamed private projection or envelope copy in gitguard, benchguard.go, or writeguard makes package-core-guard red | planned TestGuardProjectionOwnerBitesOnPrivateCopies in internal/conformance/guard_projection_owner_test.go | API presence alone permits the retired implementation to survive |
| GP31 | 21 | A consumer tree without kit guard source passes the ownership subcheck | planned TestGuardProjectionOwnerOptionalSurface in internal/conformance/guard_projection_owner_test.go | The new invariant blocks ordinary linked repositories |
| GP32 | 21 | A kit source tree with a missing owner receives a named diagnostic | planned TestGuardProjectionOwnerMissingOwner in internal/conformance/guard_projection_owner_test.go | Missing source makes the ownership check vacuously green |
| GP33 | 21 | The registered live-tree check grades the root with explicit live-tree classification | planned TestGuardProjectionOwnerHoldsOnTheLiveTree in internal/conformance/guard_projection_owner_test.go | An unbound helper or unclassified test escapes ordinary conformance |
| GP34 | 6, 7 | Family A refuses with the parse diagnostic at the missing-core Git hook | `internal/systemtest/guard_refusal_test.go` (`TestDegradedGuardRefusesMalformedCommands`, GR01–GR05) | A migration that touches the shared lexer lets a malformed command reach the coarse path |
| GP35 | 8 | The live-core Git hook and the guard-git entry emit the same parse-repair refusal | `internal/systemtest/guard_refusal_test.go` (`TestGitGuardHookAndCoreParseRefusal`, GR26–GR27) | A migrated label reintroduces unrelated authority advice or splits hook and core text |
| GP36 | 15 | The Git hook without its shared library keeps exit 2 and the missing-library refusal | `internal/systemtest/guard_refusal_test.go` (`TestGitGuardMissingLibraryRefuses`, GR20) | Removing the independent rim changes cold recovery behavior |
| GP37 | 16 | Missing-core Git commands retain the coarse refusal tuple | `internal/systemtest/guard_rim_test.go` (`TestDegradedGuardRimDecidesFromTheCommandField`) | Core-only tests cannot observe the missing-core rim |
| GP38 | 15, 16 | The file-write hook without its library or core keeps its warning-allow tuple | planned TestGuardProjectionWriteHookDegraded in internal/systemtest/guard_projection_test.go | Shared authority imposes Git's refusal on the write hook |
| GP39 | 17 | The Git hook keeps exit 2 and the analyzer-error refusal when its core exits 3 | `internal/systemtest/guard_refusal_test.go` (`TestGitGuardCoreErrorRefuses`, GR21–GR22) | Refactoring stderr forwarding changes degraded posture |
| GP40 | 14 | Valid non-Git missing-core recovery commands retain exit 0 | `internal/systemtest/guard_rim_test.go` (`TestDegradedGuardRimDecidesFromTheCommandField`) | Blanket refusal disables recovery |
| GP41 | 14 | Stale-core Bench and non-Bench calls retain their different refusal tuples | `internal/systemtest/bench_follow_on_test.go` (`TestBenchFollowOnHookDegradedRim`) | Treating exit 2 uniformly loses stale-core recovery |
| GP42 | 15, 16 | Existing resolver, classifier-table, and parity checks retain their registered execution paths | planned TestGuardProjectionPreservesIndependentChecks in internal/conformance/guard_projection_owner_test.go | A migration amputates the shell and core duplication evidence |
| GP43 | 22 | Both gitguard-pinned canary families keep biting through their registered owners | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | Removing the encoding/json import breaks a mutation anchor without a visible behavior change |
| GP44 | 24 | A supported Bench invocation produces no raw census record | planned TestGuardProjectionCensusExclusion in cmd/bench/guard_projection_test.go | A stale invocation consumer records a Bench command as raw shell |
| GP45 | 25 | Projection of an expansion-shaped command leaves a sentinel absent | planned TestGuardProjectionNeverExecutes in internal/systemtest/guard_projection_test.go | An evaluator executes inspected shell text |
| GP46 | 23, 26 | Differential family D produces identical caller result tuples before and after migration | planned TestGuardProjectionDifferentialFamilies in internal/systemtest/guard_projection_test.go | Existing narrow suites miss a changed supported command or failure policy |
| GP47 | 9, 10 | An uppercase COMMAND carrying git reset still makes guard-git exit 2 | planned TestGuardGitEnvelopePolicy in cmd/bench/guard_projection_test.go | Exact-only lookup silently permits a dangerous command |
| GP48 | 9, 11 | An uppercase COMMAND remains missing to the Bench reader | planned TestBenchEnvelopeRequiredField in internal/benchguard/projection_contract_test.go | Uniform folded matching broadens exact-map policy |
| GP49 | 9, 10 | A command string followed by duplicate null still makes guard-git exit 2 | planned TestGuardGitEnvelopePolicy in cmd/bench/guard_projection_test.go | Last-raw-value lookup discards the retained Git command |
| GP50 | 9, 11 | A command string followed by duplicate null gives Bench its empty-field error | planned TestBenchEnvelopeValue in internal/benchguard/projection_contract_test.go | Scalar-null preservation is incorrectly applied to a raw map entry |
| GP51 | 9, 10 | Nonempty tool_input followed by duplicate null still makes guard-git exit 2 | planned TestGuardGitEnvelopePolicy in cmd/bench/guard_projection_test.go | Uniform map-null clearing loses nested-struct state |
| GP52 | 9, 11 | Nonempty tool_input followed by duplicate null gives Bench its missing-field error | planned TestBenchEnvelopeRequiredField in internal/benchguard/projection_contract_test.go | Nested-struct null preservation is incorrectly applied to the map |
| GP53 | 9, 12 | Nonempty tool_input followed by duplicate null gives Write its missing-field error | planned TestWriteEnvelopeCwdContract in internal/writeguard/envelope_contract_test.go | A retained path survives the current map-clearing boundary |
| GP54 | 9, 10 | Numeric command followed by a valid string leaves guard-git at silent exit 0 | planned TestGuardGitEnvelopePolicy in cmd/bench/guard_projection_test.go | Last-value validation discards the earlier typed-field error |
| GP55 | 9, 11 | Numeric command followed by a valid string gives Bench the final string without a decode error | planned TestBenchEnvelopeType in internal/benchguard/projection_contract_test.go | Eager typed decoding rejects a replaced RawMessage |
| GP56 | 9, 12 | A valid cwd followed by null still makes guard-file-write exit 2 for a tracked relative primary path | planned TestGuardWriteEnvelopePolicy in cmd/bench/guard_projection_test.go | Raw root-map replacement loses cwd and permits the write |
| GP57 | 9, 10 | Uppercase TOOL_INPUT still makes guard-git exit 2 for git reset | planned TestGuardGitEnvelopePolicy in cmd/bench/guard_projection_test.go | Exact root lookup loses existing folded struct matching |
| GP58 | 9, 11 | Uppercase TOOL_INPUT gives Bench the exact lowercase command value | planned TestBenchEnvelopeType in internal/benchguard/projection_contract_test.go | A uniform exact-map root representation loses the input |
| GP59 | 9, 12 | Uppercase CWD and TOOL_INPUT still make guard-file-write exit 2 for a tracked relative primary path | planned TestGuardWriteEnvelopePolicy in cmd/bench/guard_projection_test.go | Exact root matching loses the existing write authority facts |
| GP60 | 9 | Repeated tool_input objects preserve earlier members under each form's standard merge rules | planned TestEnvelopeInputShape in internal/shellcommand/envelope_test.go | A last-object representation discards members absent from the later object |
| GP61 | 9, 12 | A numeric cwd followed by a valid cwd remains an outer Write decode error | planned TestWriteEnvelopeCwdContract in internal/writeguard/envelope_contract_test.go | Last-value decoding silently clears an earlier typed-field error |
| GP62 | 27 | ReviewerVerb refuses `bash -c "sh -c 'bench commitment override --assignment x --reason y'"` | planned TestReviewerVerbKeepsRecursiveReach in internal/benchguard/projection_contract_test.go | A one-level walk on the reviewer path passes the nested wrapper |
| GP63 | 27 | ReviewerVerb refuses `bench worktree exec t -- bench worktree exec u -- bench commitment override` | planned TestReviewerVerbKeepsRecursiveReach in internal/benchguard/projection_contract_test.go | A one-level exec walk passes the nested exec child |
| GP64 | 27 | Every TestGuardDeniesTheOverrideVerb command keeps its denied or allowed result | `internal/benchguard/benchguard_test.go` (`TestGuardDeniesTheOverrideVerb`) | The migration changes the recognized reviewer-verb set |
| GP65 | 28 | A command with a pool reference and the reviewer verb emits the pool refusal | planned TestGuardProjectionReviewerPrecedence in cmd/bench/guard_projection_test.go | A reordered entry names the reviewer verb instead of the pool cause |
| GP66 | 28 | A command with the reviewer verb and a shell follow-on emits the reviewer-verb refusal | planned TestGuardProjectionReviewerPrecedence in cmd/bench/guard_projection_test.go | Classify before the reviewer check emits the follow-on refusal |
| GP67 | 21 | A renamed private wrapper-flag or exec-child copy in reviewer.go makes package-core-guard red | planned TestGuardProjectionOwnerBitesOnPrivateCopies in internal/conformance/guard_projection_owner_test.go | A file list that omits reviewer.go lets its private grammar survive |
| GP68 | 29 | The agent-route guard check still flags a Bench child under bench worktree exec | `internal/conformance/refusal_route_guard_test.go` (`TestAgentRouteGuardCheckBites`, RR12) | The migration loses exec-child detection in the route check |
| GP69 | 29 | runsBenchUnderExec extracts the exec child only through WorktreeExec | review-owned (Standards axis) | GP30 grades production files only, so only review sees a test-side copy |
| GP70 | 27 | The follow-on core entry keeps exit 2 and the reviewer-verb message for each override envelope | `cmd/bench/main_test.go` (`TestFollowOnGuardDeniesTheOverrideEnvelope`) | The migrated walk breaks the actual entry path |
| GP71 | 15, 16 | The Bench follow-on hook keeps warning-allow when its resolver, wrapper, or core is missing | `internal/systemtest/bench_follow_on_test.go` (`TestBenchFollowOnHookDegradedRim`, FOG19 and FOG33) | Shared authority imposes Git's refusal on the follow-on hook |
| GP72 | 18 | `mybench worktree exec X -- git reset --hard` keeps the empty Git label | planned TestGitProjectionAuthority in internal/gitguard/projection_contract_test.go | A Git consumer that calls WorktreeExec without its Bench basename test refuses a non-Bench head |
| GP73 | 18, 27 | ReviewerVerb allows `x worktree exec t -- bench commitment override` | planned TestReviewerVerbKeepsRecursiveReach in internal/benchguard/projection_contract_test.go | A reviewer walk that calls WorktreeExec without isBench denies a non-Bench head |
| GP74 | 2 | The Git guard reads only the first matching child in W4 and keeps the allow verdict for `bash -c 'echo ok' -c 'git reset --hard'` | planned TestGuardWrapperSelection in internal/gitguard/projection_contract_test.go | A Git consumer that searches every child refuses the second-flag reset |

GP46 compares the candidate executable with tuples recorded from a pre-migration executable.
Build that executable at the integrated source at build entry, before GP1 begins, and retain it outside the tree.
Ticket 5 runs its driver once against that executable and stores the tuples as `internal/systemtest/testdata/guard-projection-family-d.golden`, reviewed apart from the candidate code.
The landed test compares the candidate with that golden, so it needs no retained executable.
The same canned events and fixture repositories feed both runs through actual core and hook entries.
Include the ordered raw envelopes from the decoder preservation matrix, without JSON remarshal or key normalization.
GP47, GP49, GP51, GP54, GP56, GP57, and GP59 require actual guard entry witnesses.

Compare exit, stdout, and stderr exactly through executable entries.
Compare reader results and decode errors in their existing owner-package contract tests.
Review owns the reference capture and each intentional difference, and this spec authorizes no intentional difference.

### Edge inventory

| Family | Exact inputs or constructors | Winning policy and observation |
|---|---|---|
| P | assignment, env options, command -v/-V, nohup, timeout, xargs from TestResolveRoutinePrefix | existing prefix projection, including nonexecuting command queries |
| W | sh, bash, zsh, absolute shell paths, -c/-lc, missing operand, invalid index, `bash -c -c 'git status'` | supported c-flag projection; each argv position is read independently, so a flag-shaped operand yields a child and also starts its own child |
| W4 | `bash -c 'echo ok' -c 'git reset --hard'`, corresponding Bench follow-on child, and a nested wrapper | Git takes the first flag, Bench Classify searches matching children, neither recurses twice |
| X | `bench worktree exec X --env A=1 -- git commit -m 'two words'`, no --, empty child, ./dist/bench, `mybench worktree exec X -- git reset --hard` | direct argv and caller executable recognition; a non-Bench head has no exec child (GP72) |
| X3 | exec inside wrapper, wrapper inside exec, nested wrappers from TestClassifyVerdicts | existing one-level allowance for Git and Bench Classify |
| R | the TestGuardDeniesTheOverrideVerb commands, a wrapper inside a wrapper, an exec inside an exec, an exec head without `--`, `x worktree exec t -- bench commitment override` | the reviewer walk recurses through each child and denies only a named verb |
| A | attached ;/&&/newline before Git reset, unclosed single quote, unclosed double quote, trailing escape | the landed degraded repair refuses malformed input |
| G | -C, --git-dir, --work-tree, -c and equals forms from TestClassifyVerdicts | Git owns option operands and repository uncertainty |
| E0 | empty stdin, invalid JSON, initially absent/null input, initially absent/null command, numeric/array command, empty command, NUL | Git silent allow, Bench warning allow |
| EP | exact raw envelope bytes in the decoder preservation matrix below | selected standard target semantics |
| F | absolute path, relative path plus cwd, missing cwd, wrong cwd type, null cwd, unreadable tracked/root fact | writeguard owns validation and warning allow |
| H | actual hooks with missing library, missing core, erroring core, stale unknown-subcommand core, valid recovery | each shim's own established failure tuple |
| D | P, W, W4, X, X3, R, A, G, E0, EP, F, H plus alias, pool conflict, heredoc, redirect, quoted operator | exact before/after observable tuples |

### Decoder preservation matrix

These inputs use ordered raw JSON bytes, not maps or remarshaled fixture objects.
The /repo spelling is a unit fixture root.
Actual file-write entry witnesses substitute a real primary checkout with a tracked README.md and no ignore rule.
They call guardFileWrite, so RootAt, the IsPrimary checker field, and tracked-path reads precede the observed exit 2.
Actual Git witnesses call guardGit, where git reset has an unconditional refusal label after a successful decode.

Existing planned reader, core-entry, and differential tests own these cases.
No parallel suite or new test harness is required.

| Case | Ordered raw input | Preserved result | Rows |
|---|---|---|---|
| nested key casing | `{"tool_input":{"COMMAND":"git reset --hard"}}` | Git decodes the command and blocks, Bench reports command field missing | GP47–GP48 |
| duplicate command null | `{"tool_input":{"command":"git reset --hard","command":null}}` | Git retains the command and blocks, Bench reports command field empty | GP49–GP50 |
| duplicate input null | `{"tool_input":{"command":"git reset --hard","file_path":"README.md"},"tool_input":null,"cwd":"/repo"}` | Git retains its nested command, Bench and Write clear their maps | GP51–GP53 |
| earlier command type error | `{"tool_input":{"command":7,"command":"git reset --hard"}}` | Git retains the earlier decode error and silently allows, Bench validates the final raw string | GP54–GP55 |
| duplicate cwd null | `{"cwd":"/repo","cwd":null,"tool_input":{"file_path":"README.md"}}` | Write retains cwd and blocks the tracked relative primary path | GP56 |
| root key casing | `{"TOOL_INPUT":{"command":"git reset --hard"}}` | Git and Bench retain folded root-key matching | GP57–GP58 |
| write root key casing | `{"CWD":"/repo","TOOL_INPUT":{"file_path":"README.md"}}` | Write retains folded root-key matching and blocks the tracked path | GP59 |
| repeated input objects | `{"tool_input":{"command":"git reset --hard","file_path":"README.md"},"TOOL_INPUT":{"other":1},"cwd":"/repo"}` | All forms retain their earlier relevant member under ordered object merging | GP60, GP46 |
| empty later object | `{"tool_input":{"command":"git reset --hard","file_path":"README.md"},"tool_input":{},"cwd":"/repo"}` | An empty later object does not clear earlier struct or map members | GP60, GP46 |
| later matching member | `{"tool_input":{"command":"git status","file_path":"old.md"},"tool_input":{"command":"git reset --hard","file_path":"README.md"},"cwd":"/repo"}` | Later matching members replace prior string or raw-map values | GP60, GP46 |
| earlier cwd type error | `{"cwd":7,"cwd":"/repo","tool_input":{"file_path":"README.md"}}` | Write retains its outer decode error and core warning-allow policy | GP61, GP29, GP46 |

Case folding applies to tagged struct fields, not to the nested maps' command and file_path keys.
Repeated nested objects merge members in all forms, but null distinguishes structs from maps.
A retained early decoder error is not erased by later valid bytes.
These are existing observable policies, including permissive cases, not new validation rules.

Tests that swap package variables stay in their own package process.
System tests use owner.runWithInput, owner.observeSelected, and privateToolPath instead of those substitutions.
Use the existing Checker, Resolver, and writeguard Checker injection without a new injected port.
The existing `injected_ports_registry_test.go` remains the authority for port contracts.
A new system test file must not redeclare `parseRepairWant`, `occurrenceCommand`, `coarseGitWant`, or `coreMissingMarker`, which `guard_refusal_test.go` declares in the same package.

**Won't handle** deeper wrapper recursion for Git and Bench Classify — gitguard.scanWords and benchguard.judge keep the one-level contract.
**Won't handle** a depth bound on the reviewer walk — the reviewer decision of 2026-10-20 keeps its recursion as guard policy.
**Won't handle** shell expansion or execution — shellcommand.Parse remains a structural parser.
**Won't handle** rg operand parsing — the separate FT366 rg outcome retains permissive operand uncertainty.
**Won't handle** shared authority policy — gitguard.Classify, benchguard.Classify, and writeguard.Classify retain their separate decisions.
**Won't handle** check-agent-line model fields — internal/lines (`lines.go:109-125`) owns that separate envelope contract.
**Won't handle** the agent-intent capture envelope — `captureClaudeAgentIntent` (`cmd/bench/main.go:324-333`) reads non-guard fields and stays outside the governed packages.
**Won't handle** the shell rim's exact-key envelope match — `.bench/lib/resolve-bench.sh:252-271` stays independent under decision answer 2, so its keys may differ from Go's folded keys.
**Won't handle** pool path expansion — benchguard.PoolReference retains its existing literal-path policy.

## Preserved closure and hook anchors

`internal/canary/pins.go:81` (`FixturePins`) maps gitguard.go to both canary directories.
The canonical-path mutation anchors its import edit on `"encoding/json"` (`MUTATE.json:4`).
Re-anchor that mutation in GP2, when gitguard.go loses the import, without a change to its canonical-path diagnostic.
The injected-ports mutation keeps its `type Checker struct` target (`gitguard.go:25-31`) and its diagnostic.
GP43 observes both families through the registered canary owner.

The build does not edit `.bench/hooks/block-dangerous-git.sh`.
No anchor registry entry pins its prose at f07508b2a.
Its threat-model and one-level comments (lines 10, 14, and 16) stay as written.

## Out of scope

The degraded malformed-input repair landed as its own outcome; this spec cites its tests as preserved regressions.
Guard scan simplification is a separate capability in `specs/guard-scan-loop/spec.md`, with one vertical chunk.
The rg operand guard remains separate roadmap work, estimated at 3 production edits and 2 gate runs after its own specification.
A full shell interpreter is excluded, with no price, because it is outside the approved threat model.

## Further notes

### Source-clause accounting

The shared ready map remains in place, because it also owns the scan outcome and the rg residual.
Each answer 1 clause occurs once there.
Answer 3 repeats the integration preservation obligation for the overall sequence.

| Source clause | Disposition |
|---|---|
| “Deepen shellcommand with shared executable, wrapper, and worktree exec-child projection.” | GP01–GP09, GP62–GP69, GP72–GP74 |
| “Use one common envelope reader for structural input facts.” | GP14–GP19, GP30, GP47–GP61 |
| “Each guard retains required fields, authority, diagnostics, and fail posture.” | GP10–GP13, GP20–GP29, GP36–GP41, GP46–GP61, GP65–GP66, GP71 |
| “Preserve the current one-level wrapper contract” | GP04, GP09, GP74; the reviewer decision of 2026-10-20 scopes it to Git and Bench Classify |
| “Keep Git global-option meaning in the Git guard.” | GP10, family G in GP46 |
| “carry the child argv and its different-checkout context.” | GP05–GP10 |
| “It must not evaluate expansions or execute shell input.” | GP45 |
| “The new rg guard retains its explicitly permissive operand-uncertainty rule.” | Won't handle rg operand parsing; GP46 preserves existing guard posture |
| “That rule must not weaken dangerous-command parse refusal.” | GP34–GP35 |
| “Preserve this distinction from `roadmap/FT366.md`.” | Won't handle rg operand parsing; GP34–GP35 keep the dangerous-command refusal separate |
| “Use three separately reviewable outcomes under FT366.” | the landed degraded outcome, this spec, and the guard-scan-loop spec |
| “First, repair malformed degraded input and its diagnostic.” | landed predecessor outcome |
| “Second, migrate common grammar and envelope facts.” | GP1–GP3 |
| “Third, simplify the guard inventory scan if its behavior stays fixed.” | separate guard-scan-loop spec |
| “retain tests through actual hook and guard entries.” | GP21, GP25, GP29, GP34–GP41, GP46, GP70–GP71 |
| “Cover attached separators, unclosed quotes, wrapper depth, Git options, and worktree exec children.” | GP02–GP10, GP34–GP35, GP46 |
| “Cover missing library, missing core, core error, and valid recovery commands.” | GP36–GP41, GP71 |
| “Prove that each guard retains its own uncertainty and diagnostic policy.” | GP10, GP21, GP25, GP28–GP29, GP35–GP41, GP46 |

Flagged additions for reviewer veto:

- The agent-route check migration (GP68–GP69) is not in the decision source. It applies the one-source code standard in `AGENTS.md` to a test-side copy of the exec grammar.
- The refresh drops seven paths from the 2026-10-06 fence, as Graph ordering states.
- The spec keeps `Related roadmap: FT366` instead of a `Roadmap:` line. FT366 closes only when both bound specs land, so a retiring spec must not remove the shared row.

The structural-copy ownership check and the differential witness enforce the approved single owner and the preservation promises.
They are implementation verification, not a new user capability.

### Seam research citations

Each seam beyond the decision source has one cited site.
The map's spec-writer discretion (`decisions/architecture-guards.md:27-28`) covers projection and envelope placement inside shellcommand.

- `ShellChildren` replaces `internal/gitguard/scan.go:67-75`, `internal/benchguard/benchguard.go:141-155` and `:315-322`, and `internal/benchguard/reviewer.go:36-38`.
- `WorktreeExec` replaces `internal/gitguard/scan.go:88-98`, `internal/benchguard/benchguard.go:251`, `internal/benchguard/reviewer.go:30-32`, and `internal/conformance/refusal_route_guard_test.go:84-91`.
- The envelope API replaces `internal/gitguard/gitguard.go:114-124`, `internal/benchguard/benchguard.go:29-51`, and `internal/writeguard/writeguard.go:23-52`.
- `checkGuardProjectionOwner` extends the package-core-guard orchestration at `internal/conformance/package_core_checks_test.go:22-37`, registered at `internal/conformance/registry/checks.go:19`.

### Build-process constraints

These rows of `docs/known-issues.md` bind the build of this spec:

- Row 10 (`docs/known-issues.md:10`): rebuild the worktree binary after each Go ticket commit, before a system run or preflight. GP3 runs `bench test --check system`.
- Row 13 (`docs/known-issues.md:13`): after each assignment record, run the plan amendment before the chunk record.
- Row 14 (`docs/known-issues.md:14`): a test-only repair needs a confirming round on all three axes.
- Row 15 (`docs/known-issues.md:15`): run the evidence currency check through `bench worktree exec`.
- Row 16 (`docs/known-issues.md:16`): the seam research citations above satisfy it.

Structure obligations: package_core_checks_test.go must not grow, and it shrinks by the moved declaration.
benchguard.go has 387 lines and must stay within 400 after GP2.
checks_test.go must stay within its 739-line grant.
Each new file must stay within 400 lines.
No file in the derived scope needs a new grant.

### Pre-review proof checklist

- Cited symbols: the refresh research read each current owner and test definition at f07508b2a. Future tests use the planned marker.
- Import edges: shellcommand imports the standard library only (`shellcommand.go:4-8`). The writeguard import of shellcommand creates no cycle.
- Source-row clauses and occurrences: the table quotes every answer 1 clause and each relevant answer 3 clause. Answer 2 stays with its landed outcome.
- Promised field labels: structural forms retain the command, tool_input, and cwd tags. Existing file_path, segment, and operator contracts retain their producers.
- Changed-function callers: the sweep below lists each caller of each retired recognizer and reader.
- Copy survival: GP30 and GP67 restore renamed private recognizers and envelope decoders beneath the executed package-core-guard check.
- Rendered-shape readers: no rendered label changes. cmd/bench/guards.go, internal/census, entry-point parity, the classifier table, and system hook readers retain their contracts.
- Pin operators: FixturePins keys exact repo-relative source paths. The check floor adds `./cmd/bench` for a Writes path under a commandPackages prefix (`internal/tickets/registry_data.go:8-28`).
- Entry reads: the cmd guards own stdin, cwd, resolver, and repository reads. Hook resolver and core-error branches retain their existing entry graders.
- Derived expectations: the independent reference capture and the named current tests own GP46 expectations. Candidate projection output does not author its own oracle.
- Consolidated rules: the migration and decoder-form tables list each replaced site and its retained target semantics.
- Quantified obligations: family D includes the ordered decoder cases EP and the reviewer family R. GP30 and GP67 list the retired copy forms, including unused bodies.
- Workflow-step writes: none. This spec adds no publisher, command, cache, or execution workflow.

The caller sweep is tied to f07508b2add49469d686614c70f575531e8ea445.
`gitguard.scanWords` callers are scan and its direct-child recursion.
`worktreeExecChild` is called only by scanWords (`scan.go:77`).
`benchguard.judge` is called by Classify and its wrapper recursion.
`spanInvokesBench` is called by judgeExec and invokes, which serves InvokesBench and census.

`isWrapper` is called by judge (`benchguard.go:141`), spanInvokesBench (`:315`), and runsReviewerVerb (`reviewer.go:36`).
`isExecHead` is called by judge (`benchguard.go:134`) and runsReviewerVerb (`reviewer.go:30`).
`wrapperFlag` is read by judge (`benchguard.go:146`), spanInvokesBench (`:319`), and runsReviewerVerb (`reviewer.go:38`).
`ReviewerVerb` is called by the follow-on core entry (`cmd/bench/guards.go:174`) and its own wrapper recursion.
The three envelope wrappers are called by cmd/bench/guards.go (`:115`, `:162`, `:201`) and their owner-package reader tests.

Preserve the core entry source bytes in cmd/bench/guards.go, because the exported wrappers keep their signatures.
Preserve internal/census.Record and the InvokesBench callers through GP44 and the existing census tests.
Preserve .bench/lib/resolve-bench.sh and the three hook entry files through independent conformance and system witnesses.
Preserve registry metadata and the existing injected-port declarations without a new injected parameter.

### Review disposition

| Item | Proposed disposition |
|---|---|
| Implementation line | opus/high, GP2 (tickets 2 and 3) and GP3 at xhigh |
| Seams and coverage | value, caller, core, actual hook, and registered ownership seams |
| Derived Writes scope | three chunk scopes above, seven paths dropped from the prior fence |
| Exclusions | independent rim, reviewer-walk depth, rg outcome, scan outcome, and full interpreter boundaries |
| Ticket graph | five tickets; GP2 holds tickets 2 and 3 in series, and GP3 holds tickets 4 and 5 with disjoint writes |

### Numbered breakdown

| ticket | title | Blocked by | delivered outcome |
|---|---|---|---|
| 1-share-git-command-projections.md | Share Git command projections | none | Shared projection facts through the real Git classifier |
| 2-migrate-guard-envelope-consumers.md | Migrate guard envelope readers | 1-share-git-command-projections.md | All three standard decoder forms through migrated readers and actual core entries |
| 3-migrate-bench-projection-consumers.md | Migrate Bench projection consumers | 1-share-git-command-projections.md, 2-migrate-guard-envelope-consumers.md | Bench Classify, the recursive reviewer walk, and the route check through the shared projection |
| 4-enforce-guard-projection-ownership.md | Enforce guard projection ownership | 1-share-git-command-projections.md, 2-migrate-guard-envelope-consumers.md, 3-migrate-bench-projection-consumers.md | Registered copy refusal for every retired private grammar copy |
| 5-prove-guard-entry-regression.md | Prove guard entry regression | 1-share-git-command-projections.md, 2-migrate-guard-envelope-consumers.md, 3-migrate-bench-projection-consumers.md | Actual hook and core regression against a checked-in pre-migration golden |
