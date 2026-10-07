# Share guard command projection and envelope facts

Status: staged
Decision source: `decisions/architecture-guards.md`, answers 1 and 3, reviewed 2026-10-06
Subject: `a9c395e77fec36d1f60f6d057c654aecf5720e24` with production source `37ac80b8adeff4ecbd57f3d813ff8937f2bcbc68`
Audience: every repository that links the kit
Related roadmap: FT366
Verification log: independent SPEC acceptance precedes slicing. The high reviewer accepted the decoder/source repairs at 02c799aa with zero findings and confidence 10. The xhigh placement confirmer accepted 6c13b4b3 with zero findings and confidence 9.
Independent graph review accepted aa86c64d on all three axes with zero findings and confidence 9.

## Problem

Git and Bench guards share `shellcommand.Parse` but each recognizes shell wrappers itself.
The Git guard also projects worktree children in `internal/gitguard/scan.go` (`scanWords`, `worktreeExecChild`).
The Bench guard recognizes the same exec head in `internal/benchguard/benchguard.go` (`isExecHead`).
Their envelope readers repeat JSON structure in those packages and `internal/writeguard/writeguard.go` (`PathFromEnvelope`).
These copies can drift while each guard's different policy must remain intact.

## Solution

Make `internal/shellcommand` own command projection and common envelope facts.
Keep each guard's authority, required fields, diagnostics, and failure policy.
Preserve the supported command language and one-level traversal through actual hook and core entries.
Keep the independent degraded shell rim because it must run without the core.

## User stories

Line: `gpt-5.6-sol` / high
Implementation-line reason: caller policy is explicit, but structural ownership checks and real degraded-entry evidence require careful integration at resolved seams
Harder chunks: GP3, use xhigh effort on the declared model
Author line: `gpt-6.1-sol` / high, one draft and at most two bounded review repairs

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

## Implementation decisions

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
An invalid executable index or missing flag operand yields no child at that position.
GP02, GP03, and GP04 observe these cases.

`WorktreeExec` receives an executable index whose identity the caller has already accepted.
It recognizes adjacent `worktree` and `exec` words and finds the first later `--`.
`Child` contains the following argv unchanged, without string joining or reparsing.
`HasSeparator` records that delimiter even when the child is empty.
`DifferentCheckout` is true exactly when the projected child is nonempty.
GP05 through GP08 observe these cases.

Git retains its shell-keyword trimming before `ResolveRoutinePrefix`.
Git consumes the first matching shell child, including a child that finds no Git invocation.
Bench retains its existing search across matching children for a Bench invocation.
Each consumer limits traversal to one level across shell and worktree child forms.
GP03, GP04, and GP09 prevent the shared facts from replacing those policy choices.

Git retains `FindSubcommand`, global-option parsing, `redirectsRepository`, and verdict classification.
It propagates `DifferentCheckout` into the existing `elsewhere` decision before consulting repository facts.
Bench retains `isBench`, resolver calls, follow-on precedence, `judgeExec`, and the pool-reference guard.
It uses `ExecProjection.Head` only after its own Bench executable recognition.
GP10 through GP13 preserve these decisions.

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
InputString decodes the selected map value as a string, matching the current two-stage map readers.
GP14–GP19 and GP47–GP61 observe these contracts.

| Form | Exact target shape | Current caller | Retained decoder facts |
|---|---|---|---|
| CommandStruct | anonymous struct with tagged tool_input nested struct and tagged command string | gitguard.CommandFromEnvelope | folded root and command keys, nested-struct merging, scalar-null preservation, retained earliest type error |
| InputMap | anonymous struct with tagged tool_input map[string]json.RawMessage | benchguard.CommandFromEnvelope | folded root key, exact input-map keys, ordered map merging, map-null clearing, final raw value validation |
| CwdInputMap | anonymous struct with tagged cwd string and tagged tool_input map[string]json.RawMessage | writeguard.PathFromEnvelope | folded root keys, exact input-map keys, scalar cwd-null preservation, map-null clearing, retained outer type errors |

Keep the existing field names and JSON tags in these anonymous wire targets.
This also preserves encoding/json error context rather than introducing named-target error text.
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
Write selects CwdInputMap and retains outer cwd validation before path handling, including absolute paths.
It retains path cleanup and the existing relative-path requirement.
Neither Git nor Bench begins validating an unrelated cwd field.
GP20–GP29 and the decoder preservation matrix below observe each branch.

A folded lookup plus one final RawMessage map is insufficient for CommandStruct and scalar cwd semantics.
Do not substitute that representation for the selected standard decoder target.
Do not migrate to strict JSON, duplicate rejection, or uniform null handling.
Those changes would alter the approved guard policies.

### Migration and enforcement

| Current site | Shared fact after migration | Policy that stays local | Rows |
|---|---|---|---|
| `gitguard.scanWords` wrapper predicate and flag loop | `ShellChildren` | first matching child and one-level traversal | GP02–GP04, GP09 |
| `gitguard.worktreeExecChild` | `WorktreeExec` | Bench basename recognition and child authority | GP05–GP10 |
| `benchguard.judge`, `spanInvokesBench`, `isWrapper` | `ShellChildren` | Bench identity and matching-child selection | GP03–GP04, GP11 |
| `benchguard.isExecHead` | `WorktreeExec.Head` | outer exec exception and follow-on order | GP12–GP13 |
| Git, Bench, and write envelope readers | `DecodeEnvelope` with the current target form | required fields and failure posture | GP14–GP29, GP47–GP61 |

Retire the private wrapper recognizers and private exec-child grammar after their callers migrate.
Do not introduce a second executable-prefix helper.
Keep `ResolveRoutinePrefix` as the existing owner.
The compiler and GP29 detect missing migrations.
GP30 detects surviving private copies, including renamed unused helpers.

Add `checkGuardProjectionOwner` beneath the existing registered `package-core-guard` check.
It grades production Go files in shellcommand, gitguard, benchguard, and writeguard, using AST import bindings.
Outside shellcommand, reject JSON unmarshal or decoder calls used to decode these envelopes.
Reject the migrated shell-name predicate, wrapper-c flag recognizer, and adjacent worktree-exec child extraction outside the owner.
The checker grades executable AST nodes, including unused function bodies, rather than comments or fixture string data.
GP30 requires renamed copies of each retired recognizer to make the registered check red.

The subcheck applies only to kit source surfaces containing conformance and at least one governed guard package.
Absent consumer source surfaces produce no diagnostic.
Present kit surfaces with a missing or unreadable owner produce a named diagnostic.
Do not require optional source directories in repositories that only link the kit.
GP31 and GP32 prove both dispositions.

Retain the existing registry name, tier, subject, and input binding.
Do not copy a registry count into this spec or add a new count expectation.
Classify the planned live-tree ownership test in `tier_live_tree_test.go`.
GP33 observes registered reachability and that classification.

### Conformance orchestration placement

Before extending the registered check, move only checkPackageCoreAndGuards from package_core_checks_test.go to checks_test.go.
Both existing ownership fences include these files.
Keep its name, signature, existing subcheck order, and callers unchanged.
Keep the other producers, npm formatter call, and their tests in package_core_checks_test.go.
Add the new ownership subcheck to the relocated function.

The accepted source has 470 lines in package_core_checks_test.go against the default 400-line limit.
It has no file grant, so adding an orchestration call there would fail the required growth lane.
checks_test.go has 641 lines against its existing 709-line grant, leaving 68 lines for the move and extensions.
The moved declaration spans 16 lines; preserve the destination cap and keep each new file within 400 lines.
Do not add a budget or acceptance grant.

Apply this placement only if the declaration remains in the original file when implementation begins.
If the other accepted guard outcome already moved it, extend that sole declaration in checks_test.go.
Retain both ownership subchecks if both outcomes have landed.
Either outcome can land first; neither requires the other's unlanded implementation.
[Placement evidence](placement-growth-repair.md) records the lane mechanism, source readers, and closure disposition.

### Ordering and independent rim

Implement the accepted `degraded-guard-refusal` outcome before the grammar migration reaches its integration checkpoint.
Its accepted plan is a prerequisite, not evidence that production already contains the repair.
The subject still contains the permissive degraded fallback described in answer 2.
GP34 and GP35 preserve the repaired refusal and parse-only advice through both actual entries.

Leave shell resolver and degraded grammar implementations in place.
They cannot import the Go owner when the core is missing.
Keep existing classifier-table, resolver-order, and entry-point-parity checks as independent evidence.
Resolver and classifier checks remain beneath package-core-guard.
Entry-point parity keeps its separate entry-point-parity registry binding.
GP36 through GP42 preserve the missing-library, missing-core, core-error, stale-core, and recovery cases.

## Implementation chunks

Each independently accepted outcome now has one complete vertical ticket and review checkpoint.
Slicing preserves every acceptance predicate and the reviewed implementation line.
Each chunk carries its source claims, consumer migration, tests, and a green ordinary test phase.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
|---|---|---|---|---|
| GP1 / 1-share-git-command-projections.md | Git uses shared wrapper and worktree projections with unchanged verdicts | GP01–GP03, GP05–GP10 | projection and Git contract tests | no |
| GP2 / 2-migrate-guard-envelope-consumers.md | common envelope facts serve the three readers and Bench uses projections | GP04, GP11–GP29, GP47–GP61 | envelope, Bench, write, and core entry tests | no |
| GP3 / 3-enforce-guard-projection-ownership.md | copy refusal and actual-entry integration complete the migration | GP30–GP46 | owner mutation tests and system hook witnesses | yes, xhigh |

GP1 precedes GP2 because the shared projection API must exist before the second consumer migrates.
GP3 enables its final ownership invariant after all consumer migrations.
GP2 owns GP47–GP61 exactly once and closes every decoder predicate at its first-use green checkpoint.
GP3 repeats those cases through final differential and actual-entry witnesses under GP46, without duplicating Covers ownership.
Run GP3 after the degraded repair prerequisite is implemented.
The shared conformance fences require serial writes during this build.

### Graph ordering and execution boundary

The stable old-to-new mapping is GP1 to GP1, GP2 to GP2, GP3 to GP3.
Each ticket retains its accepted outcome and independent Standards, Spec, and Coverage checkpoint.
Shared projection tickets run GP1, then GP2, then GP3; each successor waits for its predecessor's green commit and independent review.
GS1 is independently useful and consumes no unlanded projection interface.

Cross-spec conformance and command-registry writes require a serial integrated-source handoff.
The coordinator chooses either complete outcome first and refreshes source and headroom before overlapping successor charges.
The first ticket extending package-core-guard performs the accepted orchestration move; the other extends that sole declaration.
Neither plan replaces the other outcome's already landed subcheck or its independent npm producer assertions.

This graph's Writes union has 25 paths, equal to its unchanged path fence minus the review pickup.
Each introducing ticket owns its canonical fixture, binding, and live-tree closure and closes its own headroom obligations.
GP04 closes with both consumers in GP2, while GP1 still proves the first-use Git selection contract.
Every GP47–GP61 decoder case closes in GP2 and remains a mandatory GP3 actual-entry regression.

The version-1 completion plan is a supported planning verification inventory, not a dispatch plan.
Before any implementation charge, the coordinator binds actual run/session identities and ticket verification ownership in version 2 with author-limit 1.
Each ticket receives a fresh author on gpt-5.6-sol/high, with GP3 and GS1 at the accepted xhigh effort.
Refresh current owner interfaces, source closure, and actual headroom before that charge.
The complete degraded repair and independently captured pre-migration executable precede GP1, not only final integration.
Staged status and planning validation grant no build authority.

### Completion plan

These commands and mutations are future obligations; none ran during authoring.
Before retaining any independent expectation, pin its mutation source and exact diagnostic and demonstrate a compiling behavioral red.
Require byte-identical restoration and the same focused green; invalid or restore-failed proof closes no obligation.
Final verification runs only after this plan's complete independently reviewed outcome.

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
        "2-migrate-guard-envelope-consumers.md"
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
        }
      ]
    },
    {
      "id": "GP3",
      "tickets": [
        "3-enforce-guard-projection-ownership.md"
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
          "probe": "Omit only the registered checkGuardProjectionOwner call while retaining its function and tests; renamed unused copies must stop producing the expected registered diagnostic and red this witness."
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
`TestPathFromEnvelope` supplies the path-reader seam.
`TestGuardGitBlockAllow` supplies the core entry seam.
`TestBenchFollowOnHookDegradedRim` and `TestDegradedGuardRimDecidesFromTheCommandField` supply actual degraded-entry seams.

### Seam diagram

```text
Bash or file-tool envelope -> caller reader -> shellcommand structural facts
                                              |
command words -> shared projection -> caller authority and verdict -> core exit
                                                               |
actual harness hook --------------------------------------> shim exit
```

Tests attach to returned values, caller verdicts, core exits, and actual shim exits.
The registered ordinary conformance test detects private copies in the executed root.

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| GP01 | 1 | Routine-prefix family P retains its projected executable index | `internal/shellcommand/shellcommand_test.go` (`TestResolveRoutinePrefix`) | A new private prefix rule differs on query or option forms |
| GP02 | 2 | ShellChildren returns one ordered child per supported c flag in family W | planned TestShellChildren in internal/shellcommand/projection_test.go | A wrapper basename or flag omission loses a child |
| GP03 | 2 | Invalid indexes and missing wrapper operands yield no child | planned TestShellChildrenBounds in internal/shellcommand/projection_test.go | Unchecked indexes panic or invent a child |
| GP04 | 2 | Git and Bench retain their different matching-child selection in W4 | planned TestGuardWrapperSelection in internal/benchguard/projection_contract_test.go and internal/gitguard/projection_contract_test.go | A universal first-child or recursive policy changes a verdict |
| GP05 | 3 | WorktreeExec returns child argv byte-for-byte in family X | planned TestWorktreeExec in internal/shellcommand/projection_test.go | Joining and reparsing loses a quoted argv element |
| GP06 | 3 | A recognized exec head without a delimiter has no child | planned TestWorktreeExecMissingDelimiter in internal/shellcommand/projection_test.go | Header recognition is confused with runnable child discovery |
| GP07 | 3 | An empty child after a delimiter has DifferentCheckout=false | planned TestWorktreeExecEmptyChild in internal/shellcommand/projection_test.go | An empty projection fabricates child authority |
| GP08 | 4 | A nonempty projected child has DifferentCheckout=true | planned TestWorktreeExecContext in internal/shellcommand/projection_test.go | An omitted context bit makes parent facts available |
| GP09 | 2, 3 | Nested shell and exec combinations in X3 retain one-level verdicts | `internal/gitguard/verdict_test.go` (`TestClassifyVerdicts`) | Recursive expansion changes the documented depth |
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
| GP30 | 21 | Each renamed private projection or envelope copy makes package-core-guard red | planned TestGuardProjectionOwnerBitesOnPrivateCopies in internal/conformance/guard_projection_owner_test.go | API presence alone permits the retired implementation to survive |
| GP31 | 21 | A consumer tree without kit guard source passes the ownership subcheck | planned TestGuardProjectionOwnerOptionalSurface in internal/conformance/guard_projection_owner_test.go | The new invariant blocks ordinary linked repositories |
| GP32 | 21 | A kit source tree with a missing owner receives a named diagnostic | planned TestGuardProjectionOwnerMissingOwner in internal/conformance/guard_projection_owner_test.go | Missing source makes the ownership check vacuously green |
| GP33 | 21 | The registered live-tree check grades the root with explicit live-tree classification | planned TestGuardProjectionOwnerHoldsOnTheLiveTree in internal/conformance/guard_projection_owner_test.go | An unbound helper or unclassified test escapes ordinary conformance |
| GP34 | 6, 7 | Attached-separator and unclosed-quote family A retains refusal through actual Git hook and core entries | planned TestGuardProjectionMalformedEntryParity in internal/systemtest/guard_projection_test.go | Unit-only parity misses the independent degraded parser |
| GP35 | 8 | Malformed-entry refusal retains the accepted parse-repair diagnostic | planned TestGuardProjectionParseRepairAdvice in internal/systemtest/guard_projection_test.go | A migrated label reintroduces unrelated authority advice |
| GP36 | 15 | Missing-library hooks retain their established exit and diagnostic tuples | planned TestGuardProjectionMissingLibrary in internal/systemtest/guard_projection_test.go | Removing the independent rim changes cold recovery behavior |
| GP37 | 16 | Missing-core Git commands retain the coarse refusal tuple | `internal/systemtest/guard_rim_test.go` (`TestDegradedGuardRimDecidesFromTheCommandField`) | Core-only tests cannot observe the missing-core rim |
| GP38 | 16 | Missing-core Bench and file-write hooks retain warning-allow tuples | planned TestGuardProjectionMissingCore in internal/systemtest/guard_projection_test.go | Shared authority imposes Git's refusal on another hook |
| GP39 | 17 | Core-error hooks retain their established exit and diagnostic tuples | planned TestGuardProjectionCoreError in internal/systemtest/guard_projection_test.go | Refactoring stderr forwarding changes degraded posture |
| GP40 | 14 | Valid non-Git missing-core recovery commands retain exit 0 | `internal/systemtest/guard_rim_test.go` (`TestDegradedGuardRimDecidesFromTheCommandField`) | Blanket refusal disables recovery |
| GP41 | 14 | Stale-core Bench and non-Bench calls retain their different refusal tuples | `internal/systemtest/bench_follow_on_test.go` (`TestBenchFollowOnHookDegradedRim`) | Treating exit 2 uniformly loses stale-core recovery |
| GP42 | 15, 16 | Existing resolver, classifier-table, and parity checks retain their registered execution paths | planned TestGuardProjectionPreservesIndependentChecks in internal/conformance/guard_projection_owner_test.go | A migration amputates the shell/core duplication evidence |
| GP43 | 22 | Both gitguard-pinned canary families retain their original diagnostic predicates | planned TestGuardProjectionFixtureClosure in internal/conformance/guard_projection_owner_test.go | Removing encoding/json invalidates a mutation anchor without a visible behavior change |
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


GP46 uses an independently preserved pre-migration executable and the candidate executable.
Prepare the reference after the degraded repair prerequisite, before GP1 begins.
The same canned events and fixture repositories feed both executables through actual core and hook entries.
Include the ordered raw envelopes in the decoder preservation matrix, without JSON remarshal or key normalization.
GP47, GP49, GP51, GP54, GP56, GP57, and GP59 require actual guard entry witnesses.

Compare exit, stdout, and stderr exactly through executable entries.
Compare reader results and decode errors in their existing owner-package contract tests.
Review owns the reference capture and each intentional difference, and no intentional difference is authorized here.

### Edge inventory

| Family | Exact inputs or constructors | Winning policy and observation |
|---|---|---|
| P | assignment, env options, command -v/-V, nohup, timeout, xargs from TestResolveRoutinePrefix | existing prefix projection, including nonexecuting command queries |
| W | sh, bash, zsh, absolute shell paths, -c/-lc, missing operand, invalid index | supported c-flag projection |
| W4 | `bash -c 'echo ok' -c 'git reset --hard'`, corresponding Bench follow-on child, and a nested wrapper | Git takes first flag, Bench searches matching children, neither recurses twice |
| X | `bench worktree exec X --env A=1 -- git commit -m 'two words'`, no --, empty child, ./dist/bench | direct argv and caller executable recognition |
| X3 | exec inside wrapper, wrapper inside exec, nested wrappers from TestClassifyVerdicts | existing one-level allowance |
| A | attached ;/&&/newline before Git reset, unclosed single quote, unclosed double quote, trailing escape | accepted degraded repair refuses malformed input |
| G | -C, --git-dir, --work-tree, -c and equals forms from TestClassifyVerdicts | Git owns option operands and repository uncertainty |
| E0 | empty stdin, invalid JSON, initially absent/null input, initially absent/null command, numeric/array command, empty command, NUL | Git silent allow, Bench warning allow |
| EP | exact raw envelope bytes in the decoder preservation matrix below | selected standard target semantics |
| F | absolute path, relative path plus cwd, missing cwd, wrong cwd type, null cwd, unreadable tracked/root fact | writeguard owns validation and warning allow |
| H | actual hooks with missing library, missing core, erroring core, stale unknown-subcommand core, valid recovery | each shim's own established failure tuple |
| D | P, W, W4, X, X3, A, G, E0, EP, F, H plus alias, pool conflict, heredoc, redirect, quoted operator | exact before/after observable tuples |

### Decoder preservation matrix

These inputs use ordered raw JSON bytes rather than maps or remarshaled fixture objects.
The /repo spelling is a unit fixture root.
Actual file-write entry witnesses substitute a real primary checkout with tracked README.md and no ignore rule.
They call guardFileWrite so RootAt, IsPrimaryCheckout, and tracked-path reads precede the observed exit 2.
Actual Git witnesses call guardGit, where git reset has an unconditional refusal label after successful decoding.

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

Case-folding applies to tagged struct fields, not to the nested maps' command and file_path keys.
Repeated nested objects merge members in all forms, but null distinguishes structs from maps.
A retained early decoder error is not erased by later valid bytes.
These are existing observable policies, including permissive cases, rather than new validation rules.

Tests that swap package variables stay in their own package process.
System tests use owner.runWithInput, owner.observeSelected, and privateToolPath instead of those substitutions.
Use existing Checker, Resolver, and writeguard Checker injection without adding an injected port.
Existing `injected_ports_registry_test.go` remains the authority for port contracts.

**Won't handle** deeper wrapper recursion — gitguard.scanWords and benchguard.judge keep the one-level contract.
**Won't handle** shell expansion or execution — shellcommand.Parse remains a structural parser.
**Won't handle** rg operand parsing — the separate FT366 rg outcome retains permissive operand uncertainty.
**Won't handle** shared authority policy — gitguard.Classify, benchguard.Classify, and writeguard.Classify retain their separate decisions.
**Won't handle** check-agent-line model fields — internal/lines owns that separate envelope contract.
**Won't handle** pool path expansion — benchguard.PoolReference retains its existing literal-path policy.

## Ownership fences

These 26 exact path entries preserve the independently accepted implementation fence.
Current author writes planning artifacts only; the accepted graph grants no implementation authority.
The future graph must co-name closure entries in each affected ticket and preserve assertions unless an approved fixture re-anchor requires a text change.

- `internal/shellcommand/projection.go` (new)
- `internal/shellcommand/projection_test.go` (new)
- `internal/shellcommand/envelope.go` (new)
- `internal/shellcommand/envelope_test.go` (new)
- `internal/gitguard/scan.go`
- `internal/gitguard/gitguard.go`
- `internal/gitguard/projection_contract_test.go` (new)
- `internal/benchguard/benchguard.go`
- `internal/benchguard/projection_contract_test.go` (new)
- `internal/writeguard/writeguard.go`
- `internal/writeguard/envelope_contract_test.go` (new)
- `cmd/bench/guard_projection_test.go` (new)
- `internal/systemtest/guard_projection_test.go` (new)
- `internal/conformance/guard_projection_owner_test.go` (new)
- `internal/conformance/package_core_checks_test.go`
- `internal/conformance/checks_test.go`
- `internal/conformance/registry/checks.go`
- `internal/conformance/tier_live_tree_test.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `tests/canary/canonical-path-owner/second-derivation/`
- `tests/canary/injected-ports/unregistered-port/`
- `reviews/shared-guard-projection.md` (new)

## Preserved closure and hook anchors

`internal/tickets/registry_data.go` (`BoundFiles`) derives the command registry closure for the new cmd test.
`internal/canary/inventory.go` (`FixturePins`) derives both canary directory closures for gitguard.go.
The canonical-path mutation currently anchors its import edit on encoding/json.
Re-anchor that mutation after envelope migration without changing its canonical-path diagnostic.
Retain the injected-ports mutation's Checker target and its diagnostic.
No other prospective production path above appears in the current BASE or MUTATE path sweep.

No anchored hook prose moves.
The Git hook keeps these exact anchor bytes:

- `this is an honest-mistake layer, not an evasion-resistant`
- `exactly one level deep by design (see internal/gitguard)`
- `misaligned agent are the git pre-push hook and bench's pooled-worktree`

The future owner tests classify TestGuardProjectionOwnerHoldsOnTheLiveTree explicitly.
Registry co-names permit consistency repair, not a new command, assertion weakening, or new semantic registry row.

## Out of scope

The degraded malformed-input repair has its own accepted spec, fences, and independent checkpoint.
Its remaining price follows `specs/degraded-guard-refusal/spec.md`, not a duplicate estimate here.
Guard scan simplification is a separate capability in `specs/guard-scan-loop/spec.md`, with one vertical chunk and one final integration run.
The rg operand guard remains separate roadmap work, estimated at 3 production edits and 2 gate runs after its own specification.
A full shell interpreter is excluded, with no implementation price because it is outside the approved threat model.

## Further notes

### Source-clause accounting

The shared ready map remains in place because it also owns the independent degraded and scan outcomes.
Each answer 1 clause occurs once there.
Answer 3 repeats the integration preservation obligation for the overall sequence.

| Source clause | Disposition |
|---|---|
| “Deepen shellcommand with shared executable, wrapper, and worktree exec-child projection.” | GP01–GP09 |
| “Use one common envelope reader for structural input facts.” | GP14–GP19, GP30, GP47–GP61 |
| “Each guard retains required fields, authority, diagnostics, and fail posture.” | GP10–GP13, GP20–GP29, GP36–GP41, GP46–GP61 |
| “Preserve the current one-level wrapper contract” | GP04, GP09 |
| “Keep Git global-option meaning in the Git guard.” | GP10, family G in GP46 |
| “carry the child argv and its different-checkout context.” | GP05–GP10 |
| “It must not evaluate expansions or execute shell input.” | GP45 |
| “The new rg guard retains its explicitly permissive operand-uncertainty rule.” | Won't handle rg operand parsing, GP46 preserves existing guard posture |
| “That rule must not weaken dangerous-command parse refusal.” | GP34–GP35 |
| “First, repair malformed degraded input and its diagnostic.” | predecessor spec prerequisite |
| “Second, migrate common grammar and envelope facts.” | GP1–GP3 |
| “Third, simplify the guard inventory scan if its behavior stays fixed.” | separate guard-scan-loop spec |
| “retain tests through actual hook and guard entries.” | GP21, GP25, GP29, GP34–GP41, GP46 |
| “Cover attached separators, unclosed quotes, wrapper depth, Git options, and worktree exec children.” | GP02–GP10, GP34–GP35, GP46 |
| “Cover missing library, missing core, core error, and valid recovery commands.” | GP36–GP41 |
| “Prove that each guard retains its own uncertainty and diagnostic policy.” | GP10, GP21, GP25, GP28–GP29, GP35–GP41, GP46 |

Flagged additions: none in product behavior.
The structural-copy ownership check and differential witness enforce the approved single owner and preservation promises.
Their extra fence files are explicit implementation verification, not a new user capability.

### Pre-review proof checklist

- Cited symbols: current owner and test definitions were read in the subject tree. Future tests use the planned marker.
- Import edges: go list confirms shellcommand is a stdlib leaf. Adding its import to writeguard creates no cycle.
- Source-row clauses and occurrences: the table quotes every answer 1 clause and relevant answer 3 clause. Answer 2 stays with its accepted predecessor.
- Promised field labels: structural forms retain command, tool_input, and cwd tags. Existing file_path, segment, and operator contracts retain their producers.
- Changed-function callers: bench consumers finds scanWords, worktreeExecChild, judge, spanInvokesBench, isWrapper, isExecHead, and each exported reader's callers.
- Copy survival: GP30 restores renamed private recognizers and envelope decoders beneath the executed package-core-guard check.
- Rendered-shape readers: no rendered label changes. cmd/bench/guards.go, internal/census, entry-point parity, classifier table, and system hook readers retain existing contracts.
- Pin operators: FixturePins keys exact repo-relative source paths. Preflight uses tickets.Covers for fixture and registry co-name containment.
- Entry reads: cmd guards own stdin, cwd, resolver and repository reads. Hook resolver and core-error branches retain their existing entry graders.
- Derived expectations: independent reference capture and named current tests own GP46 expectations. Candidate projection output does not author its own oracle.
- Consolidated rules: the migration and decoder-form tables enumerate each replaced site and its retained target semantics.
- Quantified obligations: family D includes ordered decoder cases EP. GP30 enumerates retired copy forms, including unused bodies.
- Workflow-step writes: none. This spec adds no publisher, command, cache, or execution workflow.

The caller sweep is tied to subject a9c395e77fec36d1f60f6d057c654aecf5720e24 and citation c1ce54b3baafa91dfcefa87a26b1c453c53e17d67cdc95d394c0b3d7e165803b.
`gitguard.scanWords` callers are scan and its direct-child recursion.
`worktreeExecChild` is called only by scanWords.
`benchguard.judge` is called by Classify and its wrapper recursion.
`spanInvokesBench` is called by judgeExec and invokes, which serves InvokesBench and census.

`isWrapper` is called by judge and spanInvokesBench.
`isExecHead` is called by judge.
The three envelope wrappers are called by cmd/bench/guards.go and their owner-package reader tests.

Preserve core entry source bytes in cmd/bench/guards.go because the exported wrappers keep their signatures.
Preserve internal/census.Record and InvokesBench callers through GP44 and existing census tests.
Preserve .bench/lib/resolve-bench.sh and the three hook entry files through independent conformance and system witnesses.
Preserve registry metadata and existing injected-port declarations without adding a new injected parameter.

### Review repair record

[GP-C1 repair evidence](review-r1-repairs.md) records the reviewed source, decoder definitions, standard-library contract, and bounded change.
Independent review accepted the repaired spec at 02c799aa8e7ba9f1dfcf38a14906f51ab4a6d8db with no material findings.
The xhigh placement confirmer accepted the bounded amendment at 6c13b4b before ticket slicing.
Independent graph review accepted aa86c64d; no runtime preservation proof is claimed.

### Review disposition

| Item | Proposed disposition |
|---|---|
| Implementation line | gpt-5.6-sol/high, GP3 xhigh, accepted at 02c799aa |
| Seams and coverage | value, caller, core, actual hook, and registered ownership seams, accepted at 02c799aa |
| Ownership fences | exact path fence above, placement accepted at 6c13b4b; graph accepted at aa86c64d |
| Exclusions | independent rim, rg outcome, scan outcome, and full interpreter boundaries, accepted at 02c799aa |
| Ticket graph | complete vertical graph authored after independent SPEC acceptance; graph accepted at aa86c64d |

### Numbered breakdown

| ticket | title | Blocked by | delivered outcome |
|---|---|---|---|
| 1-share-git-command-projections.md | Share Git command projections | none | Shared projection facts through the real Git classifier |
| 2-migrate-guard-envelope-consumers.md | Migrate guard envelope consumers | 1-share-git-command-projections.md | All three standard decoder forms through migrated readers and actual core entries |
| 3-enforce-guard-projection-ownership.md | Enforce guard projection ownership | 1-share-git-command-projections.md, 2-migrate-guard-envelope-consumers.md | Registered copy refusal and complete actual hook/core regression |
