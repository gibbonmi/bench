# Migrate guard envelope consumers

Blocked by: 1-share-git-command-projections.md
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/guard_projection_test.go (new), cmd/bench/help_inventory_test.go, internal/benchguard/benchguard.go, internal/benchguard/projection_contract_test.go (new), internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/gitguard/gitguard.go, internal/gitguard/projection_contract_test.go (new), internal/shellcommand/envelope.go (new), internal/shellcommand/envelope_test.go (new), internal/writeguard/envelope_contract_test.go (new), internal/writeguard/writeguard.go, tests/canary/canonical-path-owner/second-derivation, tests/canary/injected-ports/unregistered-port
Covers: GP04, GP11, GP12, GP13, GP14, GP15, GP16, GP17, GP18, GP19, GP20, GP21, GP22, GP23, GP24, GP25, GP26, GP27, GP28, GP29, GP47, GP48, GP49, GP50, GP51, GP52, GP53, GP54, GP55, GP56, GP57, GP58, GP59, GP60, GP61

## What to build

Deliver common structural envelope facts to the three exported readers and shared projections to Bench's actual callers.
Ticket 1 supplies the reviewed ShellChildren and WorktreeExec value contract and complete Git migration.
Its independent chunk review closes before this shared Git contract-test successor starts.
Bench retains alias resolution, matching-child search, outer exec exceptions, follow-on precedence, and pool-reference priority.
Its InvokesBench callers and census exclusion continue to consume the same migrated path.
The complete Git-and-Bench W4 predicate closes here, without delaying GP1's Git-side witness.

DecodeEnvelope uses encoding/json.Unmarshal with the three original anonymous target shapes and exact JSON tags.
CommandStruct preserves folded nested keys, ordered struct merging, scalar nulls, and earlier typed-field errors.
InputMap and CwdInputMap preserve exact nested map keys, ordered raw-map replacement, map-null clearing, and final string validation.
CwdInputMap also preserves folded scalar cwd keys, scalar-null retention, and earlier outer errors even for absolute paths.
Keep anonymous-target error context and ignored unrelated fields; do not normalize keys, reject duplicates, or adopt strict JSON.

Each exported wrapper retains its signature, required fields, error text, path cleanup, authority, and failure posture.
Git stays silent-allow on decode errors; Bench and Write retain warning-allow behavior at their actual core entries.
Exercise every ordered raw decoder-preservation case before this ticket's green checkpoint.
Keep those raw bytes ordered without JSON remarshal, including every permissive case and both sides of each duplicate/null boundary.

Actual guardGit and guardFileWrite witnesses cover the seven specified entry rows now.
The file-write cases use a tracked relative path in a real primary checkout with no ignore rule, exercising RootAt, IsPrimaryCheckout, and tracked reads.
These actual entry witnesses remain mandatory in GP3's final regression, even though their Covers ownership is here.

Removing Git's encoding/json import changes a real canary mutation anchor.
Re-anchor only the affected canonical-path mutation here, at the first changing consumer, and retain both pinned fixture units.
Do not wait for GP3 to repair the fixture or its assertions.

The two canonical canary units retain their named diagnostic predicates and sufficient assertions.
Keep BASE, EXPECT, Checker targets, and their meanings intact.
Only the canonical-path mutation needle invalidated by removing Git's JSON import may be re-anchored to the migrated source.
The injected-ports mutation retains its existing Checker target.
Require each existing fixture to bite through its registered owner and restore its original subject.

Keep all existing assertions and independently authored expectations effective.
Do not rewrite old tests to match the new owner or derive expected results from the candidate producer.
For each required omission, pin the landed source and exact diagnostic before running it.
Accept only a compiling behavioral red, byte-identical restoration, and the same focused green afterward.
Invalid, compilation-only, skipped, or restore-failed proof closes no obligation.

Every first-use caller, fixture, registry, and headroom obligation closes in its introducing checkpoint.
A final family audit cannot supply a missing earlier test or repair.
Co-owned registry files may stay unchanged when their existing bindings and assertions suffice.
No new scanner, injected port, count expectation, structure grant, or policy authority enters this slice.

## Acceptance

- [ ] Git and Bench preserve their different W4 child-selection policies after both consumers migrate.
- [ ] Bench aliases, outer exec/heredoc exceptions, pool-over-follow-on diagnostic priority, and existing census tests retain their behavior.
- [ ] Invalid JSON, scalar or array roots and tool_input, initially absent or null objects, unrelated invalid cwd, and first null raw strings retain the structural result.
- [ ] Uppercase COMMAND blocks git reset through guardGit but remains missing to Bench's exact map reader.
- [ ] Command-string then null preserves Git's blocking command but gives Bench its empty-field error.
- [ ] Input-object then null preserves Git's nested command but clears Bench and Write maps with their missing-field errors.
- [ ] Numeric-command then string leaves Git at silent exit 0 with its earlier decode error while Bench accepts the final raw string.
- [ ] Cwd-string then null and uppercase CWD/TOOL_INPUT retain the actual tracked-relative-primary Write refusal.
- [ ] Uppercase TOOL_INPUT retains folded root matching for both actual command readers.
- [ ] Repeated, empty-later, and replaced-member tool_input objects preserve the original ordered merge rules under all three target shapes.
- [ ] Numeric-cwd then string remains Write's earlier outer decode error; even an absolute file_path cannot bypass cwd decoding.
- [ ] Git, Bench, and Write preserve required-field, type, empty, NUL, missing-cwd, cleaned-path, and unreadable-fact dispositions with unchanged error context.
- [ ] The seven real entry cases GP47/49/51/54/56/57/59 execute now, including real repository authority reads for GP56/59.
- [ ] Substituting Git's selected form with InputMap compiles and reds TestGuardGitEnvelopePolicy's ordered decoder cases.
- [ ] Substituting Bench's selected form with CommandStruct compiles and reds TestBenchEnvelopeRequiredField and TestBenchEnvelopeValue.
- [ ] Substituting Write's selected form with InputMap compiles and reds the tracked-relative-primary case in TestGuardWriteEnvelopePolicy.
- [ ] Each mutation has pinned source and diagnostic, byte-identical restoration, and the same focused green before its evidence is retained.
- [ ] Both canary units bite and restore at this checkpoint; the re-anchor changes no diagnostic, Checker target, or sufficient assertion.
- [ ] Existing command, AXI, routing, and help assertions remain intact; no command or injected port is added.
- [ ] Caller, decoder, actual core entry, and fixture obligations pass while GP3 ownership enforcement and final system integration remain unbuilt.
- [ ] New files stay within 400 lines; existing over-budget command and AXI expectation files receive no growth authority.
