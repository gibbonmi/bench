# Move freshness matrix to the lower seam

Blocked by: scope-direct-conformance-fixture-bites.md
Ownership fence: `internal/freshness/freshness_test.go`, `internal/conformance/gate_entry_test.go`
Integration surfaces: artifact and seal trust→existing `internal/freshness/freshness.go` + FS1; verified freshness-subcommand result→existing `internal/freshness/freshness.go` + FS2; rebuild refusal text→existing `freshness.RebuildAction` + FS1/FS2/FS3; shell verification ordering and phase handoff→existing `.bench/gate.sh` + FS3/FS4/FS5; fixture-scoping producer→scope-direct-conformance-fixture-bites.md + OP1/OP2
Contracts: root and executable artifact pair cross lower freshness test fixtures→`internal/freshness/freshness_test.go`, asserted by FS1/FS2 against real `Verify`/`Check`; one rebuild refusal crosses freshness→`internal/conformance/gate_entry_test.go`, asserted by FS3 against the real shell entry; phase marker crosses verified replacement→`internal/conformance/gate_entry_test.go`, asserted by FS4 against the real shell entry; structural demand-reduction evidence crosses the blocker ticket→`internal/conformance/gate_entry_test.go`, asserted by OP1/OP2 against the real composed tree
Closure: FS1/missing-executable, FS1/missing-seal, FS1/unreadable-seal, FS1/malformed-complete-seal, FS1/malformed-partial-seal, FS1/executable-digest-mismatch, FS1/one-rebuild, FS2/nonzero-subcommand, FS2/suppressed-output, FS2/one-rebuild, FS3/hostile-path, FS3/refusal-before-phases, FS3/one-rebuild, FS4/legacy-refusal, FS4/replacement-first, FS4/replacement-repeat, FS5/static-order, OP1/fixture-one-check, OP1/freshness-no-shell-matrix, OP2/production-fence, OP2/behavior-controls

## What to build

The exhaustive freshness artifact/refusal classes run directly at `Verify` or
`Check`, and the over-wide eight-variant shell matrix is removed only after its
lower owner is complete. The shell layer retains one untrusted-artifact journey
and the verified-legacy-to-replacement journey, including hostile paths, one
rebuild action, refusal before phases, and exactly one phase handoff per valid
invocation. After the blocker lands, this ticket also verifies the two reductions
together and proves no production path changed.

## Acceptance

- [ ] [FS1] (covers FR1) lower `Verify` coverage enumerates missing executable; missing, unreadable, malformed complete, and malformed partial seals; executable-digest mismatch; and exactly one rebuild action for every class.
- [ ] [FS2] (covers FR2) lower `Check` coverage rejects a verified executable whose `freshness-check` exits nonzero, suppresses the child's output, and emits exactly one rebuild action.
- [ ] [FS3] (covers FR3) the representative real shell entry refuses an untrusted artifact from a nested hostile path with one rebuild action and never schedules `gate-phases`.
- [ ] [FS4] (covers FR4) the real shell entry rejects a verified legacy binary and a published replacement runs phases exactly once on the first and repeated valid invocations.
- [ ] [FS5] (covers FR5) the static gate-entry contract retains current-source verification before its one `gate-phases` handoff.
- [ ] [OP1] (covers OP1) the composed tree records one registered check per direct fixture tree and keeps the exhaustive freshness table entirely at lower seams with no eight-case shell/module matrix.
- [ ] [OP2] (covers OP2) production registry, selection, diagnostics, freshness policy, gate routing, timing format, and verdict semantics are unchanged under the existing behavior controls.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| FS1/missing-executable | remove the published executable | the lower Verify table | remove it, call `Verify`, require the untrusted-binary refusal and one rebuild action |
| FS1/missing-seal | remove the matching seal | the lower Verify table | remove it, call `Verify`, require the missing-seal refusal and one rebuild action |
| FS1/unreadable-seal | make the seal unreadable on a privilege-capable host | the lower Verify table | strip read permission, call `Verify`, require refusal or the existing capability skip rather than a false green |
| FS1/malformed-complete-seal | replace the seal with a complete but schema-invalid JSON object | the lower Verify table | write the malformed object, call `Verify`, require the malformed-contents refusal |
| FS1/malformed-partial-seal | truncate the seal JSON | the lower Verify table | write the partial object, call `Verify`, require the malformed refusal |
| FS1/executable-digest-mismatch | alter the executable after publication while preserving executable mode | the lower Verify table | replace its bytes, call `Verify`, require the executable-digest mismatch refusal |
| FS1/one-rebuild | duplicate rebuild advice in the refusal formatter through an overlay | the lower refusal assertion | apply the duplication, run the Verify table, require the exact-one count to fail |
| FS2/nonzero-subcommand | publish a sealed executable whose `freshness-check` exits nonzero | the lower Check table | call `Check`, require the freshness-subcommand refusal rather than trust |
| FS2/suppressed-output | make the failing verified child print a sentinel to stdout and stderr | the lower Check table | call `Check`, capture output, require the sentinel absent from the refusal surface |
| FS2/one-rebuild | duplicate rebuild advice for a verified-child failure | the lower Check assertion | apply the duplication, run the Check case, require the exact-one count to fail |
| FS3/hostile-path | remove quoting from one gate-entry fixture path use in an overlay | the representative shell test | run from the nested spaces/glob path, require the shell entry to fail before the intended freshness refusal |
| FS3/refusal-before-phases | move the `gate-phases` handoff ahead of freshness verification in an overlay | the representative shell test | run the untrusted-artifact journey, require the phase marker's presence to fail the no-phases assertion |
| FS3/one-rebuild | print a second rebuild line from the shell entry | the representative shell test | apply the duplicate, run the journey, require the exact-one assertion to fail |
| FS4/legacy-refusal | treat nonzero `freshness-check` as trusted | the legacy shell journey | apply the bypass, run the legacy binary, require the stale phase marker/output to fail the refusal assertion |
| FS4/replacement-first | suppress the replacement `gate-phases` handoff | the replacement shell journey | publish the replacement, invoke once, require the missing phase record to fail |
| FS4/replacement-repeat | cache or skip the second valid replacement invocation | the replacement shell journey | invoke twice, require the phase record to contain two entries |
| FS5/static-order | reorder or delete current-source verification in `.bench/gate.sh` through an overlay | `checkGateEntryContract` | apply the mutation, run its focused conformance control, require the ordering diagnostic |
| OP1/fixture-one-check | revert one migrated fixture caller to empty scope in the composed tree | the first ticket's timing identity plus this ticket's composed verification | apply the omission, run both focused package commands, require the fixture timing assertion to fail |
| OP1/freshness-no-shell-matrix | restore `TestGateEntryNormalizesIndeterminateFreshnessFailures` after lower coverage lands | the checkpoint residue command | restore the function, run `! rg -n '^func TestGateEntryNormalizesIndeterminateFreshnessFailures' internal/conformance/gate_entry_test.go`, require the command to fail, then restore the candidate |
| OP2/production-fence | modify any production registry, selection, freshness, gate script, timing, or verdict path | lifecycle fence enforcement | add the out-of-fence edit, submit the checkpoint, require the lifecycle to refuse it |
| OP2/behavior-controls | make singular conformance scope ignore its check or make freshness trust a nonzero verified child in an overlay | the existing conformance and freshness behavior controls | apply one mutation at a time, run the focused package controls, require each changed behavior to red before restoring the candidate |
