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

## Record

```bench-review-record
{
  "version": 1,
  "spec": "specs/bounded-charge-evidence/spec.md",
  "plan_digest": "sha256:971cc327eafed497f8ad3ae83762bcaf4c1a30538971f90736c74ba6be2031d0",
  "implementation_session": "bounded-charge-evidence-retained-author",
  "chunks": [
    {
      "id": "CE-C1A",
      "base": "bcd9eb7ff876abd681d2c71547e9cf6c5b86174b",
      "tip": "d434bb7623126a0b632f018098c4d66b51922608",
      "plan_digest": "sha256:971cc327eafed497f8ad3ae83762bcaf4c1a30538971f90736c74ba6be2031d0",
      "source_digest": "7301ed96ce7df50b7f267a86cd9f50ae3e3a4f26",
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
  }
}
```
