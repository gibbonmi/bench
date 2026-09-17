# Review outcomes: bounded-charge-evidence

This file holds the review pickup for the `specs/bounded-charge-evidence/spec.md` implementation run.
The fenced record at the end is the machine-readable evidence.

## Run decisions

- The reviewer approved digest-verified reads for oversized charges. The author and each review axis read each omitted repository source at the pinned tip and match its sha256 to the charge identity. Each review axis also matches the diff, consumers, and coverage captures of the frozen pair to their charge identities. A mismatch stops the reader. This approval applies until tickets 4 and 5 replace the full-charge routes.
- The retained author for each ticket is a fork of the coordinator session on opus at high effort. The ticket 1 fork ended with an API error during the CE-C1A review. By reviewer direction, a sonnet write delegate performs the CE-C1A repairs. This is a user-directed author transfer after a terminal author failure.
- The review axes run opus at medium effort, by reviewer direction.

## CE-C1A: review round 1

Frozen pair: base `8e4cf48021ec5a3a576d2f60e40480059e1e8c08`, tip `3f010d4e0ed13d3300e2ed1a47641ff317090f64`.
The raw finding count is 14. The de-duplicated repair-target count is 11.
Repair cycles used: 0 of 2.

### Standards

Finding count: 5. Worst issue: ST1.

- CE-C1A-ST1 (auto-fix): The chunk claims recorded omission probes, but no red record exists in the tree. A test comment in `internal/chargeevidence/format_test.go` states that the review pickup carries the records. The repair records every probe in this pickup and removes the red-record claim from the comment. SP1 names the same repair.
- CE-C1A-ST2 (auto-fix): The header layout and the profile values are stated twice. The reference text in `internal/chargeevidence/schema.go` hard-codes the version, reserved value, and page size. `internal/chargeevidence/pack.go` uses literal offsets. Derive both from the registry constants.
- CE-C1A-ST3 (auto-fix): The source-state and preparation-refusal fixtures in `internal/preflight/charge_evidence_test.go` repeat the setups in `internal/preflight/charge_legacy_pack_test.go`. Use one case table for both tests.
- CE-C1A-ST4 (auto-fix): `internal/preflight/charge_pack.go` computes the source ordinal with `SourceID(i + 2)`. The pack package owns that rule. Export one input-ordinal helper from `internal/chargeevidence`.
- CE-C1A-ST5 (auto-fix): A comment in `internal/preflight/charge_legacy_pack_test.go` narrates history. The build path joins the fence with its own separator, apart from `chargeFenceCell`. State the current fact, and use one fence-cell owner.

### Spec

Finding count: 4. Worst issue: SP1.

- CE-C1A-SP1 (auto-fix): This finding is the same as ST1. Spec line 490 requires a recorded red for each independent expectation before acceptance.
- CE-C1A-SP2 (auto-fix): CE52 tests no ledger record and a foreign active record. No case has a record that owns this worktree in a non-active state. Add that case.
- CE-C1A-SP3 (auto-fix): The CE148 test matches each field label anywhere in the document. An omitted field can pass through a label in another table. Check each field inside its own table row.
- CE-C1A-SP4 (auto-fix): The shipped reference omits facts that an independent reader needs. These facts are the relative page offset, zero-based indexes, the canonical string quoting rules, and the build and review access values. Add them through the registry and the projection.

Flagged for reviewer veto, with no finding ID: spec line 474 names `TestChargeProjectionAndFullRetrieval` as the CE-C1A differential, but row CE147 names `TestLegacyPreparedPackDifferential`. The implementation obeys the row. The CE18 case also changes the ticket, because review mode has an empty ticket. The axis marked that point no-op.

### Coverage

Finding count: 5. Worst issue: CV1 to CV3.

- CE-C1A-CV1 (auto-fix): No test reaches the reader's page-split refusal. A probe that disabled the check was silent.
- CE-C1A-CV2 (auto-fix): No test reaches the reader's source-digest refusal, including a zero-byte optional source. A probe that disabled the check was silent.
- CE-C1A-CV3 (auto-fix): No test reaches the reader's source-bytes refusal. A probe that disabled the check was silent.
- CE-C1A-CV4 (auto-fix): No case pins the integer upper bound or the other noncanonical integer forms. A probe that removed the bound was silent.
- CE-C1A-CV5 (auto-fix): No test covers generated-source provenance rows. A probe that disabled the provenance validation was silent. The coordinator routes this to the repair because the reader already enforces these rows in this chunk.

### Advice

- Rename the `flag` schema helper, because the name reads like the standard library package.
- The legacy differential compares bytes only, and it cannot show that the output came from the pack.
- `Metadata.references()` adds each charge ticket cell, so a later review row with an empty ticket can fail the reader.
- Consider cases for a byte order mark and for a path with a comma or a newline in a manifest cell.

### Author verification

The ticket 1 author ran the plan probe and 34 other probes, and it reported every verdict as `bit` with the file restored. Its detailed log was lost when the author session ended. The repair author reruns and records every probe at the repaired tip.
The coordinator's independent probe swapped the checks and return columns in `internal/preflight/charge_pack.go`. Five tests failed, and the file was restored.

## CE-C1A: review round 2

Frozen pair: base `bcd9eb7ff876abd681d2c71547e9cf6c5b86174b`, tip `d434bb7623126a0b632f018098c4d66b51922608`.
The base is the main tip that the coordinator merged into the source after repair cycle 1.
The raw finding count is 4. The de-duplicated repair-target count is 4.
Repair cycles used: 1 of 2. Repair cycle 1 closed ST1, ST3, ST4, ST5, SP1, SP2, SP3, and CV1 to CV5.

### Standards

Finding count: 3. Worst issue: ST6.

- CE-C1A-ST2 (auto-fix): The header marker text `BENCHEV` is still restated in the reference text in `internal/chargeevidence/schema.go` and in a reader message in `internal/chargeevidence/pack.go`. Derive both from `HeaderMarker`.
- CE-C1A-ST6 (auto-fix): `internal/chargeevidence/format_refusals_test.go` repeats the header re-frame of `repack` in a second helper, with literal offsets. Use one re-frame helper.
- CE-C1A-ST7 (auto-fix): `internal/chargeevidence/reference.go` hard-codes the access values that `internal/preflight/charge_pack.go` and `internal/preflight/review.go` set. Give the access values one owner in `internal/chargeevidence`, and format the reference from it.

### Spec

Finding count: 1. Worst issue: SP4.

- CE-C1A-SP4 (auto-fix): The quoting paragraph in the shipped reference reads as a complete list, but it omits encoder triggers. The omitted triggers are leading or trailing whitespace, a leading-zero decimal, a colon, a bracket or brace, a leading hyphen, and the comma. Give the complete trigger set, or state only that the pinned encoder owns quoting.

### Coverage

Finding count: 0.

### Advice

- The producer-source identity check in `internal/chargeevidence/manifest.go` is silent under a probe. That check predates repair cycle 1.
- Derive `HeaderBytes` from the header table, and give the metadata ordinal one named owner.
- Add comma-bearing and colon-bearing strings to the canonical profile cases.

### Author verification

The coordinator ran the three chunk verifications at `d434bb76`. Both focused suites passed, and the plan probe failed seven tests and restored the file.
The repair writer logged 31 probe runs at `76afb4fb`, and each run failed the expected tests and restored the file.
The coordinator's independent probe shifted the input ordinal in `internal/chargeevidence/manifest.go`, and six tests failed.

## CE-C1A: review round 3

Frozen pair: base `bcd9eb7ff876abd681d2c71547e9cf6c5b86174b`, tip `85de8fa38a9d90808a6034e3dda0029b860cdebc`.
The raw finding count is 2. The de-duplicated repair-target count is 1, because both findings concern one quoting paragraph.
Repair cycles used: 2 of 2. Repair cycle 2 closed ST2, ST6, ST7, and SP4.

### Standards

Finding count: 1. Worst issue: ST8.

- CE-C1A-ST8 (auto-fix): `internal/chargeevidence/reference.go` adds a sentence that answers the removed trigger list. A later reader learns nothing from it. Delete the sentence.

### Spec

Finding count: 1. Worst issue: SP5.

- CE-C1A-SP5 (auto-fix): The quoting paragraph names a Bench test file as the trigger inventory. Spec line 257 requires a reference that an independent reader can use without Bench-private helpers. Name the upstream encoder module and the version that `go.mod` pins instead.

### Coverage

Finding count: 0.

### Advice

- Derive the ASCII marker text from `HeaderMarker` without a literal length.
- Keep `HeaderLengthRange` out of the package API through a test-only export.
- A header length probe fails by a panic instead of a named assertion.

### Reviewer decision

The allowance was exhausted with SP5 open. The reviewer extended the allowance by one repair cycle. The extension covers only the quoting paragraph in `internal/chargeevidence/reference.go` and the regenerated shipped reference. A fresh opus write delegate at medium effort performs it. Round 4 reviews only that change.

### Author verification

The coordinator ran the three chunk verifications at `85de8fa3`. Both focused suites passed, and the plan probe failed seven tests and restored the file.
The repair writer logged nine probe runs at `85de8fa3`, and each run failed the expected tests and restored the file.
The coordinator's independent probe changed the ticket-cell sentence in the reference generator, and the projection test failed.

## CE-C1A: review round 4

Frozen pair: base `bcd9eb7ff876abd681d2c71547e9cf6c5b86174b`, tip `3d9932917645b955b012070c5a2862bfb560cbe8`.
The round reviewed only the scoped extension cycle. The raw finding count is 0.
Repair cycles used: 2 of 2, plus the one reviewer-approved scoped extension. The extension closed ST8 and SP5.

### Standards

Finding count: 0. The axis passed.

### Spec

Finding count: 0. The axis passed.

### Coverage

Finding count: 0. The axis passed.

### Advice

- The Standards axis reported that `EncoderModule` restates the module path that the encoder imports, and that its test reads `go.mod` instead of the import. By reviewer decision, this item is advice without an ID. The ticket 2 author ties the test to the encoder import, and the CE-C1B review verifies that change.

### Author verification

The coordinator ran the three chunk verifications at `3d993291`. Both focused suites passed, and the plan probe failed seven tests and restored the file.
The extension writer probed the module constant and the removed sentence, and both probes failed the projection test.
The coordinator's independent probe removed the version sentence from the generator, and the projection test failed.

## CE-C1B: review round 1

Frozen pair: base `88a981b9cf9952ccb59be918fe72cae3c73f16d3`, tip `1829a5c8641de6029801e87a3ca1fc9584600270`.
The raw finding count is 16. The de-duplicated repair-target count is 15, including the reviewer-directed directory move.
Repair cycles used: 0 of 2.

### Reviewer decisions for this chunk

- The root help check accepts one named preflight help projection, and the check moved into its own conformance test file. The plan expansions at `9e78a26b` and `1e7500b3` added both files to the fence and ticket 2.
- The production hook `BENCH_EVIDENCE_PAUSE` is accepted, including its resume marker file.
- CE-C1A advice ST9 was carried into ticket 2, and the Spec axis confirmed that the encoder test now reads the adapter import.
- CE75 uses a pure admission function test with a candidate above the default quota.
- The evidence command surface moves into its own subpackage under `internal/preflight` in this repair.

### Standards

Finding count: 7. Worst issue: ST1.

- CE-C1B-ST1 (auto-fix): The independent test constants for the response budget, the store name, and the default quota have no recorded red. Probe each against its production owner and record the result.
- CE-C1B-ST2 (auto-fix): A preflight test repeats the literal response budget beside the package test constant.
- CE-C1B-ST3 (auto-fix): Each operation row states its flags in its lists and again in its usage string. Derive the usage text from the lists.
- CE-C1B-ST4 (auto-fix): The selector flags, the grammar flags, and the operation flag sets restate one fact. Derive them from one source.
- CE-C1B-ST5 (auto-fix): Three test packages restate the store directory name and the pack and temporary name parts. Use the exported owners.
- CE-C1B-ST6 (auto-fix): Several comments narrate history, cite the spec as provenance, mislabel a constant block, or restate the cursor grammar.
- CE-C1B-ST7 (auto-fix): `internal/preflight` holds 41 source files against a budget of 12. By reviewer decision, the evidence command surface moves into a subpackage.

### Spec

Finding count: 4. Worst issue: SP1.

- CE-C1B-SP1 (auto-fix): The CE128 test refuses at repository resolution and never reaches source inspection. Make the missing tool fail inside source inspection, and prove that no artifact is published. Coverage CV3 names the same repair.
- CE-C1B-SP2 (auto-fix): CE75 never admits an artifact above the default quota. Test the admission function directly with a large candidate. Coverage CV4 names the same repair.
- CE-C1B-SP3 (auto-fix): The response budget rows do not fail when a guard call is removed from the command dispatch. Push an over-bound response through the command for each path.
- CE-C1B-SP4 (auto-fix): CE94 had no durable pending record. This pickup now records it below.

### Coverage

Finding count: 4 after the merge. Worst issue: CV1.

- CE-C1B-CV1 (auto-fix): A legitimate large preparation refusal becomes a generic response bound defect line and loses its check and recovery action. Bound the diagnostic operands so the refusal keeps its meaning.
- CE-C1B-CV2 (auto-fix): An empty existing store refuses as `unsafe-store`, but the spec requires a missing-artifact refusal.
- CE-C1B-CV5 (auto-fix): No test replaces the store directory between inspection and open. Add a registered in-process test port at that point.
- CE-C1B-CV6 (auto-fix): No test covers the refusal for an unrepresentable required quota.

### Advice

- `BENCH_EVIDENCE_PAUSE` waits without a limit, and a stray value can overwrite an existing file through its resume marker.
- `refusalClass` exists in two packages.
- Add cursor cases for uppercase hex and the integer boundary.
- The help check constrains only the projection name.
- A symbolic link above the Git common directory is followed.
- Spec lines 472 to 474 expect the compact differential to remain until its consumer migrates. The author replaced the compact route, and the Spec axis flagged this contradiction as non-behavioral.

### Pending evidence

- CE94: the device store-kind case skips on this host because the user cannot create a device node. The row stays pending until a privileged native run records it.

### Author verification

The coordinator ran the five chunk verifications at `1829a5c8`. All passed, and the plan probe failed 23 tests and restored the file.
The ticket 2 author logged 29 probe runs and two manual system probes, each with a failing result and a verified restore.
The coordinator's independent probe changed the required bytes in the capacity refusal, and the capacity recovery test failed.

## Record

```bench-review-record
{
  "version": 1,
  "spec": "specs/bounded-charge-evidence/spec.md",
  "plan_digest": "sha256:ad140f7fc480da95134f328448ea09d840c4b3e49c8e1c5a76ef9518d56544f5",
  "implementation_session": "bounded-charge-evidence-retained-author",
  "chunks": [
    {
      "id": "CE-C1A",
      "base": "bcd9eb7ff876abd681d2c71547e9cf6c5b86174b",
      "tip": "3d9932917645b955b012070c5a2862bfb560cbe8",
      "plan_digest": "sha256:971cc327eafed497f8ad3ae83762bcaf4c1a30538971f90736c74ba6be2031d0",
      "source_digest": "ff09af8a6fdd9a16b23d8b34074a3c3ea5492c78",
      "acceptance_rows": [
        "CE16",
        "CE18",
        "CE19",
        "CE20",
        "CE21",
        "CE22",
        "CE23",
        "CE24",
        "CE25",
        "CE26",
        "CE27",
        "CE28",
        "CE29",
        "CE30",
        "CE31",
        "CE32",
        "CE33",
        "CE34",
        "CE35",
        "CE36",
        "CE37",
        "CE38",
        "CE41",
        "CE42",
        "CE43",
        "CE44",
        "CE45",
        "CE46",
        "CE47",
        "CE48",
        "CE49",
        "CE50",
        "CE51",
        "CE52",
        "CE111",
        "CE114",
        "CE125",
        "CE126",
        "CE147",
        "CE148",
        "CE149"
      ],
      "verification": [
        {
          "id": "ce-c1a-v1-preflight",
          "performer": "bounded-charge-evidence-retained-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "894a28f8060219c74a2a94b309c391fe64bca045",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-session:coordinator-run",
            "digest": "sha256:c78422f71f008d447d6f392ea1aea0cc77aa99f6638fdefa7d4b8b9458bb399e",
            "excerpt": "bench test --package ./internal/preflight at 3f010d4e: pass, 24410 ms"
          },
          "requirement": "preflight",
          "command": "bench test --package ./internal/preflight",
          "exit_code": 0
        },
        {
          "id": "ce-c1a-v1-format",
          "performer": "bounded-charge-evidence-retained-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "894a28f8060219c74a2a94b309c391fe64bca045",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-session:coordinator-run",
            "digest": "sha256:07c4fe7861fce740985adf16004677ef54309fe2fb6033939ced9f9ed632f807",
            "excerpt": "bench test --package ./internal/chargeevidence at 3f010d4e: pass, 7 ms"
          },
          "requirement": "format",
          "command": "bench test --package ./internal/chargeevidence",
          "exit_code": 0
        },
        {
          "id": "ce-c1a-v1-mutation",
          "performer": "bounded-charge-evidence-retained-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "894a28f8060219c74a2a94b309c391fe64bca045",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-session:ticket-1-author-fork",
            "digest": "sha256:3afb2043634ebcd33a24c3cc0342b1f8996de9735be4d6f43bfe6571be1853db",
            "excerpt": "bench test --package ./internal/chargeevidence at 3f010d4e: pass"
          },
          "requirement": "mutation",
          "command": "bench test --package ./internal/chargeevidence",
          "exit_code": 0,
          "probe": {
            "mutation": "omit a declared field from the shipped format projection",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-session:ticket-1-author-fork",
              "digest": "sha256:69c25644490d231693b5eecfc125651f28dccc080fd7b7f0d3d1455f4042610e",
              "excerpt": "bench probe internal/chargeevidence/schema.go --omit ', str(\"cwd\")' --package ./internal/chargeevidence: verdict bit, 4 failed tests, restored=yes"
            }
          }
        },
        {
          "id": "ce-c1a-v2-preflight",
          "performer": "bounded-charge-evidence-retained-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7301ed96ce7df50b7f267a86cd9f50ae3e3a4f26",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-session:coordinator-run",
            "digest": "sha256:092f59109a94a2fd65291d8215d6cd98eda5de89714b4cda4593cb22afee720c",
            "excerpt": "bench test --package ./internal/preflight at d434bb76: pass, 17823 ms"
          },
          "requirement": "preflight",
          "command": "bench test --package ./internal/preflight",
          "exit_code": 0
        },
        {
          "id": "ce-c1a-v2-format",
          "performer": "bounded-charge-evidence-retained-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7301ed96ce7df50b7f267a86cd9f50ae3e3a4f26",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-session:coordinator-run",
            "digest": "sha256:3f010ff9f6e12a4dcd879fbb2cae3c77eaa17988483a612999676e810782a87e",
            "excerpt": "bench test --package ./internal/chargeevidence at d434bb76: pass, 9 ms"
          },
          "requirement": "format",
          "command": "bench test --package ./internal/chargeevidence",
          "exit_code": 0
        },
        {
          "id": "ce-c1a-v2-mutation",
          "performer": "bounded-charge-evidence-retained-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7301ed96ce7df50b7f267a86cd9f50ae3e3a4f26",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-session:coordinator-run",
            "digest": "sha256:3f010ff9f6e12a4dcd879fbb2cae3c77eaa17988483a612999676e810782a87e",
            "excerpt": "bench test --package ./internal/chargeevidence at d434bb76: pass, 9 ms"
          },
          "requirement": "mutation",
          "command": "bench test --package ./internal/chargeevidence",
          "exit_code": 0,
          "probe": {
            "mutation": "omit a declared field from the shipped format projection",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-session:coordinator-run",
              "digest": "sha256:7be9da9656ee31c775d169268d1affdab1e4c91f451548ed16de54fad4f7aec2",
              "excerpt": "bench probe internal/chargeevidence/schema.go --omit ', str(\"cwd\")' --package ./internal/chargeevidence at d434bb76: verdict bit, 7 failed tests, restored=yes"
            }
          }
        },
        {
          "id": "ce-c1a-v3-preflight",
          "performer": "bounded-charge-evidence-retained-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a9540853957b34f88b3f5ecb14fe4930f439406d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-session:coordinator-run",
            "digest": "sha256:29aa55e29ab195a78f6b19559682939596b5282e3c3ce8054058a8ad2aefd48e",
            "excerpt": "bench test --package ./internal/preflight at 85de8fa3: pass, 17411 ms"
          },
          "requirement": "preflight",
          "command": "bench test --package ./internal/preflight",
          "exit_code": 0
        },
        {
          "id": "ce-c1a-v3-format",
          "performer": "bounded-charge-evidence-retained-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a9540853957b34f88b3f5ecb14fe4930f439406d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-session:coordinator-run",
            "digest": "sha256:441191ee37556bde55a06cdf39ad37f42d052797ab6a3b3390315afe695e99d7",
            "excerpt": "bench test --package ./internal/chargeevidence at 85de8fa3: pass, 9 ms"
          },
          "requirement": "format",
          "command": "bench test --package ./internal/chargeevidence",
          "exit_code": 0
        },
        {
          "id": "ce-c1a-v3-mutation",
          "performer": "bounded-charge-evidence-retained-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a9540853957b34f88b3f5ecb14fe4930f439406d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-session:coordinator-run",
            "digest": "sha256:441191ee37556bde55a06cdf39ad37f42d052797ab6a3b3390315afe695e99d7",
            "excerpt": "bench test --package ./internal/chargeevidence at 85de8fa3: pass, 9 ms"
          },
          "requirement": "mutation",
          "command": "bench test --package ./internal/chargeevidence",
          "exit_code": 0,
          "probe": {
            "mutation": "omit a declared field from the shipped format projection",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-session:coordinator-run",
              "digest": "sha256:71513b041e3213348a018fbe124c206bb4ab0b93538fec7dc7e953ba6c6b5e8f",
              "excerpt": "bench probe internal/chargeevidence/schema.go --omit ', str(\"cwd\")' --package ./internal/chargeevidence at 85de8fa3: verdict bit, 7 failed tests, restored=yes"
            }
          }
        },
        {
          "id": "ce-c1a-v4-preflight",
          "performer": "bounded-charge-evidence-retained-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ff09af8a6fdd9a16b23d8b34074a3c3ea5492c78",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-session:coordinator-run",
            "digest": "sha256:5c12cb1a9fd00a6aae7c00b24149ffda51c462f5bb54c5df8b8042d4f2290d9e",
            "excerpt": "bench test --package ./internal/preflight at 3d993291: pass, 18108 ms"
          },
          "requirement": "preflight",
          "command": "bench test --package ./internal/preflight",
          "exit_code": 0
        },
        {
          "id": "ce-c1a-v4-format",
          "performer": "bounded-charge-evidence-retained-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ff09af8a6fdd9a16b23d8b34074a3c3ea5492c78",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-session:coordinator-run",
            "digest": "sha256:3a5f9633edeebc7becac0a9296c042aa3be12252351762bbd40c4bc3c73dbf45",
            "excerpt": "bench test --package ./internal/chargeevidence at 3d993291: pass, 8 ms"
          },
          "requirement": "format",
          "command": "bench test --package ./internal/chargeevidence",
          "exit_code": 0
        },
        {
          "id": "ce-c1a-v4-mutation",
          "performer": "bounded-charge-evidence-retained-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "ff09af8a6fdd9a16b23d8b34074a3c3ea5492c78",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-session:coordinator-run",
            "digest": "sha256:3a5f9633edeebc7becac0a9296c042aa3be12252351762bbd40c4bc3c73dbf45",
            "excerpt": "bench test --package ./internal/chargeevidence at 3d993291: pass, 8 ms"
          },
          "requirement": "mutation",
          "command": "bench test --package ./internal/chargeevidence",
          "exit_code": 0,
          "probe": {
            "mutation": "omit a declared field from the shipped format projection",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-session:coordinator-run",
              "digest": "sha256:25394e590ed5e19d64fdae77d91137ff16d20ec6a45ca809dc019b01ad42a255",
              "excerpt": "bench probe internal/chargeevidence/schema.go --omit ', str(\"cwd\")' --package ./internal/chargeevidence at 3d993291: verdict bit, 7 failed tests, restored=yes"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "ce-c1a-r1-standards",
          "performer": "claude-review-ce-c1a-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "894a28f8060219c74a2a94b309c391fe64bca045",
          "state": "completed",
          "outcome": "findings",
          "native_ref": {
            "ref": "claude-subagent:ce-c1a-standards",
            "digest": "sha256:570663f168288e857b67d9cbdf537e8898b535ae4f61a27d303f43f6e391a6e9",
            "excerpt": "Standards CE-C1A: 5 findings. Worst: claimed omission-probe red records do not exist in the tree (ST1). Also: header layout restated outside the registry (ST2), duplicated fixture setups (ST3), source ordinal re-derived in preflight (ST4), history-narrating comments and a second fence join (ST5)."
          },
          "axis": "Standards",
          "base": "8e4cf48021ec5a3a576d2f60e40480059e1e8c08",
          "tip": "3f010d4e0ed13d3300e2ed1a47641ff317090f64",
          "finding_ids": [
            "CE-C1A-ST1",
            "CE-C1A-ST2",
            "CE-C1A-ST3",
            "CE-C1A-ST4",
            "CE-C1A-ST5"
          ],
          "supersedes": []
        },
        {
          "id": "ce-c1a-r1-spec",
          "performer": "claude-review-ce-c1a-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "894a28f8060219c74a2a94b309c391fe64bca045",
          "state": "completed",
          "outcome": "findings",
          "native_ref": {
            "ref": "claude-subagent:ce-c1a-spec",
            "digest": "sha256:f9a11a5339750286fed180efcbb7073a59c74dbd57f92fac1deec623807ce171",
            "excerpt": "Spec CE-C1A: 4 findings. Worst: mutation records missing while a test comment claims them (SP1). Also: CE52 lacks a non-active owned record case (SP2), CE148 checks labels document-wide (SP3), shipped reference omits facts an independent reader needs (SP4)."
          },
          "axis": "Spec",
          "base": "8e4cf48021ec5a3a576d2f60e40480059e1e8c08",
          "tip": "3f010d4e0ed13d3300e2ed1a47641ff317090f64",
          "finding_ids": [
            "CE-C1A-SP1",
            "CE-C1A-SP2",
            "CE-C1A-SP3",
            "CE-C1A-SP4"
          ],
          "supersedes": []
        },
        {
          "id": "ce-c1a-r1-coverage",
          "performer": "claude-review-ce-c1a-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "894a28f8060219c74a2a94b309c391fe64bca045",
          "state": "completed",
          "outcome": "findings",
          "native_ref": {
            "ref": "claude-subagent:ce-c1a-coverage",
            "digest": "sha256:1bb32b562202393aa445fd1ccb91349181871e49132ecb4f2091317aa60d8f5b",
            "excerpt": "Coverage CE-C1A: 5 findings. Worst: reader page-split, source-digest, and source-bytes refusals are unexercised; probes were silent (CV1-CV3). Also: integer upper bound unpinned (CV4), generated-source provenance untested (CV5)."
          },
          "axis": "Coverage",
          "base": "8e4cf48021ec5a3a576d2f60e40480059e1e8c08",
          "tip": "3f010d4e0ed13d3300e2ed1a47641ff317090f64",
          "finding_ids": [
            "CE-C1A-CV1",
            "CE-C1A-CV2",
            "CE-C1A-CV3",
            "CE-C1A-CV4",
            "CE-C1A-CV5"
          ],
          "supersedes": []
        },
        {
          "id": "ce-c1a-r2-standards",
          "performer": "claude-review-ce-c1a-standards-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7301ed96ce7df50b7f267a86cd9f50ae3e3a4f26",
          "state": "completed",
          "outcome": "findings",
          "native_ref": {
            "ref": "claude-subagent:ce-c1a-standards-r2",
            "digest": "sha256:a1f3a9eeb036a197cb005e6e48acc5b0a6092006260eca31e517f3fe8d4b32ff",
            "excerpt": "Standards CE-C1A round 2: ST1, ST3, ST4, ST5 closed. ST2 partly open: the BENCHEV marker text is restated in schema.go and pack.go. New: header re-frame pasted twice in format_refusals_test.go (ST6); access values hard-coded in reference.go apart from their owners (ST7)."
          },
          "axis": "Standards",
          "base": "bcd9eb7ff876abd681d2c71547e9cf6c5b86174b",
          "tip": "d434bb7623126a0b632f018098c4d66b51922608",
          "finding_ids": [
            "CE-C1A-ST2",
            "CE-C1A-ST6",
            "CE-C1A-ST7"
          ],
          "supersedes": [
            "ce-c1a-r1-standards"
          ]
        },
        {
          "id": "ce-c1a-r2-spec",
          "performer": "claude-review-ce-c1a-spec-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7301ed96ce7df50b7f267a86cd9f50ae3e3a4f26",
          "state": "completed",
          "outcome": "findings",
          "native_ref": {
            "ref": "claude-subagent:ce-c1a-spec-r2",
            "digest": "sha256:cfcd44bdf3b1905b58bb8819dfb72397f62ff94c26e3273bb2afbe326fbcd76e",
            "excerpt": "Spec CE-C1A round 2: SP1, SP2, SP3 closed; plan commit changes no behavior or row. SP4 partly open: the quoting paragraph in the shipped reference omits encoder quoting triggers (whitespace, leading zero, colon, brackets, leading hyphen, comma)."
          },
          "axis": "Spec",
          "base": "bcd9eb7ff876abd681d2c71547e9cf6c5b86174b",
          "tip": "d434bb7623126a0b632f018098c4d66b51922608",
          "finding_ids": [
            "CE-C1A-SP4"
          ],
          "supersedes": [
            "ce-c1a-r1-spec"
          ]
        },
        {
          "id": "ce-c1a-r2-coverage",
          "performer": "claude-review-ce-c1a-coverage-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7301ed96ce7df50b7f267a86cd9f50ae3e3a4f26",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-subagent:ce-c1a-coverage-r2",
            "digest": "sha256:4ed39a7066a2c520bf8917d20ba87fc01e4bfe96fa8d9bfe1b5dc83a9979d198",
            "excerpt": "Coverage CE-C1A round 2: 0 findings. CV1 to CV5 closed with biting probes; the repair's header range, truncation, CE52, and input-ordinal guards all bit. Advice: the producer-source identity check is silent under a probe; it predates the repair."
          },
          "axis": "Coverage",
          "base": "bcd9eb7ff876abd681d2c71547e9cf6c5b86174b",
          "tip": "d434bb7623126a0b632f018098c4d66b51922608",
          "finding_ids": [],
          "supersedes": [
            "ce-c1a-r1-coverage"
          ]
        },
        {
          "id": "ce-c1a-r3-standards",
          "performer": "claude-review-ce-c1a-standards-r3",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "a9540853957b34f88b3f5ecb14fe4930f439406d",
          "state": "completed",
          "outcome": "findings",
          "native_ref": {
            "ref": "claude-subagent:ce-c1a-standards-r3",
            "digest": "sha256:80e72420bd321a8509befd4c80f9d17e103a158da4fc1dca46baf784e7b7879e",
            "excerpt": "Standards CE-C1A round 3: ST2, ST6, ST7 closed. New ST8: reference.go adds a sentence that argues with removed text (This reference does not restate a partial trigger list); delete it."
          },
          "axis": "Standards",
          "base": "bcd9eb7ff876abd681d2c71547e9cf6c5b86174b",
          "tip": "85de8fa38a9d90808a6034e3dda0029b860cdebc",
          "finding_ids": [
            "CE-C1A-ST8"
          ],
          "supersedes": [
            "ce-c1a-r2-standards"
          ]
        },
        {
          "id": "ce-c1a-r3-spec",
          "performer": "claude-review-ce-c1a-spec-r3",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "a9540853957b34f88b3f5ecb14fe4930f439406d",
          "state": "completed",
          "outcome": "findings",
          "native_ref": {
            "ref": "claude-subagent:ce-c1a-spec-r3",
            "digest": "sha256:157d7f997c173451e2c77dad5eb59c77a695f633e21e4891aef2a81d0c049c4d",
            "excerpt": "Spec CE-C1A round 3: SP4 closed. New SP5: the quoting paragraph names a Bench test file as the trigger inventory, which contradicts spec line 257; name the upstream encoder module and its pinned version instead."
          },
          "axis": "Spec",
          "base": "bcd9eb7ff876abd681d2c71547e9cf6c5b86174b",
          "tip": "85de8fa38a9d90808a6034e3dda0029b860cdebc",
          "finding_ids": [
            "CE-C1A-SP5"
          ],
          "supersedes": [
            "ce-c1a-r2-spec"
          ]
        },
        {
          "id": "ce-c1a-r3-coverage",
          "performer": "claude-review-ce-c1a-coverage-r3",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "a9540853957b34f88b3f5ecb14fe4930f439406d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-subagent:ce-c1a-coverage-r3",
            "digest": "sha256:582b2d5de38c364ba10b65ec25195e1eb214b3993e296856f399343ddd95d804",
            "excerpt": "Coverage CE-C1A round 3: 0 findings. Every guard changed in 85de8fa3 bit under a probe: marker text, header length range, reframe helper, access constants and uses, comma and colon cells."
          },
          "axis": "Coverage",
          "base": "bcd9eb7ff876abd681d2c71547e9cf6c5b86174b",
          "tip": "85de8fa38a9d90808a6034e3dda0029b860cdebc",
          "finding_ids": [],
          "supersedes": [
            "ce-c1a-r2-coverage"
          ]
        },
        {
          "id": "ce-c1a-r4-standards",
          "performer": "claude-review-ce-c1a-standards-r4",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ff09af8a6fdd9a16b23d8b34074a3c3ea5492c78",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-subagent:ce-c1a-standards-r4",
            "digest": "sha256:835657faa6b8dae28dc18e61dc94076d4d0464f789e6b2ac26f8ec7572764b91",
            "excerpt": "Standards CE-C1A round 4: ST8 closed. The axis raised one minor non-blocking item: EncoderModule restates the encoder import path and its test anchors on go.mod. By reviewer decision this item is advice without an ID, carried into ticket 2."
          },
          "axis": "Standards",
          "base": "bcd9eb7ff876abd681d2c71547e9cf6c5b86174b",
          "tip": "3d9932917645b955b012070c5a2862bfb560cbe8",
          "finding_ids": [],
          "supersedes": [
            "ce-c1a-r3-standards"
          ]
        },
        {
          "id": "ce-c1a-r4-spec",
          "performer": "claude-review-ce-c1a-spec-r4",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ff09af8a6fdd9a16b23d8b34074a3c3ea5492c78",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-subagent:ce-c1a-spec-r4",
            "digest": "sha256:90aca4587b65f31526fcb24cd45b8ed5ed85b916a400adc1ed57cc769069d41d",
            "excerpt": "Spec CE-C1A round 4: SP5 closed; the reference names github.com/toon-format/toon-go at the go.mod pin, which meets spec lines 153 and 257. No new findings."
          },
          "axis": "Spec",
          "base": "bcd9eb7ff876abd681d2c71547e9cf6c5b86174b",
          "tip": "3d9932917645b955b012070c5a2862bfb560cbe8",
          "finding_ids": [],
          "supersedes": [
            "ce-c1a-r3-spec"
          ]
        },
        {
          "id": "ce-c1a-r4-coverage",
          "performer": "claude-review-ce-c1a-coverage-r4",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ff09af8a6fdd9a16b23d8b34074a3c3ea5492c78",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-subagent:ce-c1a-coverage-r4",
            "digest": "sha256:42e97d94105892ed24271f64216ebdec54d3b976b01cc70d0e8471594d23d167",
            "excerpt": "Coverage CE-C1A round 4: 0 findings. The EncoderModule drift probe and the quoting paragraph omission probe bit; the weakened go.mod check stays guarded by the projection test."
          },
          "axis": "Coverage",
          "base": "bcd9eb7ff876abd681d2c71547e9cf6c5b86174b",
          "tip": "3d9932917645b955b012070c5a2862bfb560cbe8",
          "finding_ids": [],
          "supersedes": [
            "ce-c1a-r3-coverage"
          ]
        }
      ]
    },
    {
      "id": "CE-C1B",
      "base": "88a981b9cf9952ccb59be918fe72cae3c73f16d3",
      "tip": "1829a5c8641de6029801e87a3ca1fc9584600270",
      "plan_digest": "sha256:ad140f7fc480da95134f328448ea09d840c4b3e49c8e1c5a76ef9518d56544f5",
      "source_digest": "3284d0947b73d4492340feda5908eb7497a3a24e",
      "acceptance_rows": [
        "CE1",
        "CE2",
        "CE3",
        "CE4",
        "CE5",
        "CE6",
        "CE9",
        "CE10",
        "CE11",
        "CE12",
        "CE13",
        "CE14",
        "CE15",
        "CE17",
        "CE39",
        "CE40",
        "CE53",
        "CE54",
        "CE63",
        "CE64",
        "CE74",
        "CE75",
        "CE76",
        "CE77",
        "CE86",
        "CE87",
        "CE88",
        "CE89",
        "CE90",
        "CE91",
        "CE92",
        "CE93",
        "CE94",
        "CE95",
        "CE96",
        "CE97",
        "CE98",
        "CE99",
        "CE118",
        "CE128",
        "CE129",
        "CE130",
        "CE131",
        "CE132",
        "CE138",
        "CE139",
        "CE150",
        "CE151",
        "CE152",
        "CE153",
        "CE154",
        "CE155",
        "CE156",
        "CE157"
      ],
      "verification": [
        {
          "id": "ce-c1b-v1-preflight",
          "performer": "bounded-charge-evidence-retained-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "3284d0947b73d4492340feda5908eb7497a3a24e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-session:coordinator-run",
            "digest": "sha256:892af95fb638aefbc880ce96a8b9dcb022600b0d3b75c72944c46abc8aec12c8",
            "excerpt": "bench test --package ./internal/preflight at 1829a5c8: pass, 22754 ms"
          },
          "requirement": "preflight",
          "command": "bench test --package ./internal/preflight",
          "exit_code": 0
        },
        {
          "id": "ce-c1b-v1-store",
          "performer": "bounded-charge-evidence-retained-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "3284d0947b73d4492340feda5908eb7497a3a24e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-session:coordinator-run",
            "digest": "sha256:d0026688d4f8d20b25aa1aaea5f7aaee2b66acca9b5ab047f6cb55d8fdf0a705",
            "excerpt": "bench test --package ./internal/chargeevidence at 1829a5c8: pass, 88 ms; CE94 device case skipped for the privilege capability"
          },
          "requirement": "store",
          "command": "bench test --package ./internal/chargeevidence",
          "exit_code": 0
        },
        {
          "id": "ce-c1b-v1-inventory",
          "performer": "bounded-charge-evidence-retained-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "3284d0947b73d4492340feda5908eb7497a3a24e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-session:coordinator-run",
            "digest": "sha256:5de083c4cce3f78d7886d5b948c7413741756c1e1b8e62edda5a91b20709cd93",
            "excerpt": "bench test --package ./cmd/bench at 1829a5c8: pass, 7468 ms"
          },
          "requirement": "inventory",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ce-c1b-v1-system",
          "performer": "bounded-charge-evidence-retained-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "3284d0947b73d4492340feda5908eb7497a3a24e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-session:coordinator-run",
            "digest": "sha256:928db49bfbace995498270b3306017a05def8ea61cb20a207f431a16943c5ad4",
            "excerpt": "bench test --check system at 1829a5c8: pass, 34695 ms"
          },
          "requirement": "system",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "ce-c1b-v1-mutation",
          "performer": "bounded-charge-evidence-retained-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "3284d0947b73d4492340feda5908eb7497a3a24e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-session:coordinator-run",
            "digest": "sha256:892af95fb638aefbc880ce96a8b9dcb022600b0d3b75c72944c46abc8aec12c8",
            "excerpt": "bench test --package ./internal/preflight at 1829a5c8: pass, 22754 ms"
          },
          "requirement": "mutation",
          "command": "bench test --package ./internal/preflight",
          "exit_code": 0,
          "probe": {
            "mutation": "omit manifest_bytes from the prepared response",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-session:coordinator-run",
              "digest": "sha256:699570ea129787a043e2f0f3dd70916e445978a5aa180d50e6eded9ecc7326ac",
              "excerpt": "bench probe internal/chargeevidence/schema.go --omit 'num(\"manifest_bytes\"),' --package ./internal/preflight at 1829a5c8: verdict bit, 23 failed tests, restored=yes"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "ce-c1b-r1-standards",
          "performer": "claude-review-ce-c1b-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3284d0947b73d4492340feda5908eb7497a3a24e",
          "state": "completed",
          "outcome": "findings",
          "native_ref": {
            "ref": "claude-subagent:ce-c1b-standards",
            "digest": "sha256:cd36ed763ad0c19add1ce656a67df35d1446cf7a392b3adcc90c732d6af958e5",
            "excerpt": "Standards CE-C1B: 6 findings plus a crowded preflight directory. Worst: independent test constants for the response budget, store name, and default quota have no recorded red (ST1)."
          },
          "axis": "Standards",
          "base": "88a981b9cf9952ccb59be918fe72cae3c73f16d3",
          "tip": "1829a5c8641de6029801e87a3ca1fc9584600270",
          "finding_ids": [
            "CE-C1B-ST1",
            "CE-C1B-ST2",
            "CE-C1B-ST3",
            "CE-C1B-ST4",
            "CE-C1B-ST5",
            "CE-C1B-ST6",
            "CE-C1B-ST7"
          ],
          "supersedes": []
        },
        {
          "id": "ce-c1b-r1-spec",
          "performer": "claude-review-ce-c1b-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3284d0947b73d4492340feda5908eb7497a3a24e",
          "state": "completed",
          "outcome": "findings",
          "native_ref": {
            "ref": "claude-subagent:ce-c1b-spec",
            "digest": "sha256:860243751d44c9a1c798dc1e3164bf9abaf0b49aa325b7001532284aef795aee",
            "excerpt": "Spec CE-C1B: 4 findings. Worst: the CE128 test refuses at repository resolution and never reaches preparation (SP1). CE75 lacks an above-quota artifact (SP2); budget rows miss a removed guard call (SP3); CE94 lacks a pending record (SP4)."
          },
          "axis": "Spec",
          "base": "88a981b9cf9952ccb59be918fe72cae3c73f16d3",
          "tip": "1829a5c8641de6029801e87a3ca1fc9584600270",
          "finding_ids": [
            "CE-C1B-SP1",
            "CE-C1B-SP2",
            "CE-C1B-SP3",
            "CE-C1B-SP4"
          ],
          "supersedes": []
        },
        {
          "id": "ce-c1b-r1-coverage",
          "performer": "claude-review-ce-c1b-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3284d0947b73d4492340feda5908eb7497a3a24e",
          "state": "completed",
          "outcome": "findings",
          "native_ref": {
            "ref": "claude-subagent:ce-c1b-coverage",
            "digest": "sha256:02f648906be7541972d6978944e884f439093fe3a100d8e399cad4a9639d2867",
            "excerpt": "Coverage CE-C1B: 6 findings (2 merged into SP1 and SP2). Worst: a large preparation refusal becomes a generic response-bound defect line (CV1). An empty store refuses as unsafe (CV2); the store replacement check is untested (CV5); the unrepresentable-quota refusal is untested (CV6)."
          },
          "axis": "Coverage",
          "base": "88a981b9cf9952ccb59be918fe72cae3c73f16d3",
          "tip": "1829a5c8641de6029801e87a3ca1fc9584600270",
          "finding_ids": [
            "CE-C1B-CV1",
            "CE-C1B-CV2",
            "CE-C1B-CV5",
            "CE-C1B-CV6"
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
      "from": "sha256:971cc327eafed497f8ad3ae83762bcaf4c1a30538971f90736c74ba6be2031d0",
      "to": "sha256:ad140f7fc480da95134f328448ea09d840c4b3e49c8e1c5a76ef9518d56544f5",
      "chunk_ids": {
        "CE-C1A": [
          "CE-C1A"
        ],
        "CE-C1B": [
          "CE-C1B"
        ],
        "CE-C1C": [
          "CE-C1C"
        ],
        "CE-C1D": [
          "CE-C1D"
        ],
        "CE-C2": [
          "CE-C2"
        ],
        "CE-C3": [
          "CE-C3"
        ],
        "CE-C4": [
          "CE-C4"
        ]
      }
    }
  ]
}
```
