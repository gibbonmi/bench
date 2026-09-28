# Test determinism review

## Run

The user directs one inline implementation author and independent Sol/high reviews.
The version 1 plan retains the inline author identity.
The native implementation model identifier is unknown.
Review line: gpt-6-sol / high / one iteration per axis / three readers
Expected repair rounds: 1 / confidence 7

## TD-C1a

The environment owner and gate composition cover TD1 through TD6 and TD8 through TD15.
The focused environment and gate packages pass.
The root conformance test passes with no skips.
The gittest package has no standalone tests; its probe runs through the environment and gate tests.

The independent expectations detect the recorded mutations below.
Each valid probe reports a behavioral failure and confirms restoration.
The local probe transcript is .logs/test-determinism-probes.json.

| rows | mutation | observed failure |
| --- | --- | --- |
| TD1 | Remove the global git configuration override. | The global marker remains visible. |
| TD2 | Remove the system git configuration override. | The system marker remains visible. |
| TD3 | Remove the maintenance entries. | The commit starts auto-maintenance. |
| TD4 | Preserve the operator home. | HOME is outside the private run. |
| TD5, TD13 | Preserve the operator temporary directory. | TMPDIR is outside the private run. |
| TD6 | Empty the Go setting pins. | The probe reports changed pins. |
| TD8 | Omit directory permission restoration. | Close reports permission denied. |
| TD9 | Use a long run-directory prefix. | The private path exceeds the 16-byte limit. |
| TD10, TD11 | Bypass HOME validation. | Open accepts absent, empty, and relative values. |
| TD12 | Create the run outside the specified TMPDIR. | Open accepts a regular file as TMPDIR. |
| TD14 | Apply the kit policy to a linked root. | The linked phase loses the operator home. |
| TD15 | Omit the run entries from the phase environment. | The phase reports that HOME is not private. |

The first phase-composition mutation failed to compile and supplied no behavioral evidence.
The replacement mutation compiles and produces the intended failure.

The HOME-refusal test uses a temporary working directory.
A repeated validation-bypass probe fails without a checkout write.

Repair cycles consumed: 2 of 2

## Standards

TD-C1a: zero blocking findings; worst issue: none.

## Spec

TD-C1a: zero blocking findings; worst issue: none.
Later chunks retain their planned rows; this chunk claims only the rows named above.

## Coverage

TD-C1a: zero blocking findings; worst issue: none.

## Review reads

Each axis read the frozen diff, the whole approved spec, ticket 1, and the current evidence metadata and consumers.
Each axis also read its standards and the targeted owner, tests, probe, gate composer, and runner.
Standards inspected the untouched lane consumer and the recorded probe output.
Spec audited all 14 chunk rows.
Coverage enumerated the environment inputs, phase outcomes, cleanup states, and permitted writes.
All three axes report no implementation-command contribution.

## Optional advice

Coverage suggests hostile git markers in the gate fixture.
The owner tests already exercise those markers, and the current gate composer merges every owner entry.
This advice has no finding ID or repair disposition.

## TD-C1a repair 1

Checkpoint gate-20260928T115939.135031035Z-2068436 failed in TestVerifyRefusesSpecialArtifactsBeforeReading/seal_Unix_socket.
The private temporary path extended the socket path to 110 bytes.
The race and system phases passed.
A focused run under the same temporary-directory shape reproduced the failure.

The owner uses a shorter random name and creates the directory atomically.
It retries a name collision and retains private HOME and TMPDIR subdirectories.
The existing socket fixture passes under the actual repaired owner.
The environment and gate packages pass, and the root conformance test passes without skips.
A mutation that extends the random name produces the expected path-length failure and restores the source.

The source-bound author results below supersede the earlier focused results.
All three confirming reviews pass with no blocking findings.
Each reviewer affirms the complete current pair and reports confidence 9.
The worktree remains clean after their read-only returns.

## TD-C1a repair 2

Checkpoint gate-20260928T121225.513556246Z-2346173 failed only on the formatting of the temporary socket diagnostic in .logs.
The ordinary tests, race checks, and system suite passed.
The diagnostic is now formatted with gofmt.
The tracked source and its review identities are unchanged.
No review axis needs a semantic refresh for this formatting repair.

## TD-C1b plan expansion

Ticket 2 includes the child lifetime owner, its Go fixtures, and the gate test that retires with the forwarding helper.
TD14 keeps its linked-root guarantee through the actual phase test.
All chunk IDs and acceptance assignments stay unchanged.
The chunk adds gate-package verification for its gate edits.
TD-C1a passed checkpoint gate-20260928T121615.630583303Z-2591234.

```bench-review-record
{"version":1,"spec":"specs/test-determinism/spec.md","plan_digest":"sha256:a0a05390b237129415afa3bc4ccab193d0e0c7d672e8ebe3c22c00b1cf4c3b77","implementation_session":"codex/test-determinism-inline-20260928","chunks":[{"id":"TD-C1a","base":"72e751b7a020a76aa35286852b4ea981c95cb865","tip":"7feee3e4a4714d96429719ef25d52c5dc3e9582b","plan_digest":"sha256:0c1f19eeefd300b73018e190ea2fb3acf0baf416e4171f4b8717cfa5aa651a63","source_digest":"416f75ff593d9f8fc90658e8cfe6c07489cf206b","acceptance_rows":["TD1","TD2","TD3","TD4","TD5","TD6","TD8","TD9","TD10","TD11","TD12","TD13","TD14","TD15"],"verification":[{"id":"TD-C1a-env-1","performer":"codex/test-determinism-inline-20260928","role":"author-verification","model":"unknown","effort":"high","source_digest":"706142f4f32e66ff15d4c863aa66e285a0612add","state":"completed","outcome":"pass","native_ref":{"ref":"codex:exec-session-82687","digest":"sha256:c57f65257e3b204524fd8a7c37d94fb38b4db9ed278cd03c9132d4df16fa1924","excerpt":"packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/env,pass,417\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"},"requirement":"env","command":"bench test --package ./internal/env","exit_code":0},{"id":"TD-C1a-gate-1","performer":"codex/test-determinism-inline-20260928","role":"author-verification","model":"unknown","effort":"high","source_digest":"706142f4f32e66ff15d4c863aa66e285a0612add","state":"completed","outcome":"pass","native_ref":{"ref":"codex:exec-session-87505","digest":"sha256:a280ab2fcb946d6709f0f12cf81dfdaff4805543d53fe6dd129ffcc659b34f35","excerpt":"packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,pass,10349\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"},"requirement":"gate","command":"bench test --package ./internal/gate","exit_code":0},{"id":"TD-C1a-env-repair-1","performer":"codex/test-determinism-inline-20260928","role":"author-verification","model":"unknown","effort":"high","source_digest":"416f75ff593d9f8fc90658e8cfe6c07489cf206b","state":"completed","outcome":"pass","native_ref":{"ref":"codex:exec-session-19637","digest":"sha256:17e67e3581fb1a138cc3a4f1f6c8279e0492541832bfe5c7110078b5adba8c33","excerpt":"packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/env,pass,443\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"},"requirement":"env","command":"bench test --package ./internal/env","exit_code":0},{"id":"TD-C1a-gate-repair-1","performer":"codex/test-determinism-inline-20260928","role":"author-verification","model":"unknown","effort":"high","source_digest":"416f75ff593d9f8fc90658e8cfe6c07489cf206b","state":"completed","outcome":"pass","native_ref":{"ref":"codex:exec-session-38455","digest":"sha256:63647d8de429b7d0ab8ada5ead0246ad758dbed8eb46631d1bfbf3ff8240609b","excerpt":"packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,pass,10949\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"},"requirement":"gate","command":"bench test --package ./internal/gate","exit_code":0}],"reviews":[{"id":"TD-C1a-Standards-1","performer":"/root/td_c1a_standards","role":"independent-review","model":"gpt-6-sol","effort":"high","source_digest":"706142f4f32e66ff15d4c863aa66e285a0612add","state":"completed","outcome":"pass","native_ref":{"ref":"collaboration:/root/td_c1a_standards/final","digest":"sha256:f4c12987e642d07b44a442623dc7288c99f9d23396423c5ea49130e4eeb8ba7b","excerpt":"Standards: **0 blocking findings; worst: none.** The retained git policy remains single-sourced, the shared probe avoids duplicated fixture harnesses, and the review record names behavioral mutations for the independent expectations.\nNo tests or probes executed, and no files changed. This verdict covers TD-C1a; later chunks remain pending.\nclaim[1]{status,confidence}:\n  claimed,9"},"axis":"Standards","base":"72e751b7a020a76aa35286852b4ea981c95cb865","tip":"4ad8b5319043875ea8dd74a6e5195a69929930b7","finding_ids":[],"supersedes":[]},{"id":"TD-C1a-Spec-1","performer":"/root/td_c1a_spec","role":"independent-review","model":"gpt-6-sol","effort":"high","source_digest":"706142f4f32e66ff15d4c863aa66e285a0612add","state":"completed","outcome":"pass","native_ref":{"ref":"collaboration:/root/td_c1a_spec/final","digest":"sha256:d367cd431dcad69e9d6260e29df6af1e9405a979b72a4e42373a639808e88318","excerpt":"Spec review: **0 findings; worst issue: none.** Line: gpt-6-sol / high / 1 iteration.\nAudited all 14 requested rows:\nThe author record documents reds for every requested row. I ran no tests or probes; this is source review, not independent execution verification. Later chunks remain intentionally pending.\nOptional advice: none. Implementation-command contribution: none; no command change is necessary.\nClaim: {\"status\":\"claimed\",\"confidence\":9}"},"axis":"Spec","base":"72e751b7a020a76aa35286852b4ea981c95cb865","tip":"4ad8b5319043875ea8dd74a6e5195a69929930b7","finding_ids":[],"supersedes":[]},{"id":"TD-C1a-Coverage-1","performer":"/root/td_c1a_coverage","role":"independent-review","model":"gpt-6-sol","effort":"high","source_digest":"706142f4f32e66ff15d4c863aa66e285a0612add","state":"completed","outcome":"pass","native_ref":{"ref":"collaboration:/root/td_c1a_coverage/final","digest":"sha256:dcf9d02235fa97e630df2a301a129cc08067fe15fa22df44e3fe0228d884bf42","excerpt":"Coverage: 0 blocking findings; worst: none. TD-C1a covers TD1\u2013TD6 and TD8\u2013TD15; later chunks remain pending.\nAdvice: TestKitPhaseRunsInTheKitTestRun uses ambient configuration (phases_test.go:169). Explicit hostile markers there would strengthen TD15 against selective omission of Git overrides; the shared probe currently checks marker absence (gittest.go:219) without establishing that gate fixture\u2019s marker precondition.\nCommand contribution: none. No tests or probes ran; executable refutation was prohibited by the charge.\n{\"status\":\"claimed\",\"confidence\":8}"},"axis":"Coverage","base":"72e751b7a020a76aa35286852b4ea981c95cb865","tip":"4ad8b5319043875ea8dd74a6e5195a69929930b7","finding_ids":[],"supersedes":[]},{"id":"TD-C1a-Standards-repair-1","performer":"/root/td_c1a_standards_confirm","role":"independent-review","model":"gpt-6-sol","effort":"high","source_digest":"416f75ff593d9f8fc90658e8cfe6c07489cf206b","state":"completed","outcome":"pass","native_ref":{"ref":"collaboration:/root/td_c1a_standards_confirm/final","digest":"sha256:72d4c28b451e395904fae9aa57cfe2393dfbc73b47afbfcd0850053eec32a3e2","excerpt":"Standards: **0 blocking findings; worst issue: none.** Repair targets: 0. Optional advice: none.\nI affirm the Standards pass for the complete current TD-C1a pair, `72e751b7a020a76aa35286852b4ea981c95cb865..7feee3e4a4714d96429719ef25d52c5dc3e9582b`. The repair preserves prior evidence for TD1\u2013TD6 and TD8\u2013TD15.\nNo tests or probes ran, and no files changed. This is source review; author execution remains separately recorded."},"axis":"Standards","base":"72e751b7a020a76aa35286852b4ea981c95cb865","tip":"7feee3e4a4714d96429719ef25d52c5dc3e9582b","finding_ids":[],"supersedes":["TD-C1a-Standards-1"]},{"id":"TD-C1a-Spec-repair-1","performer":"/root/td_c1a_spec_confirm","role":"independent-review","model":"gpt-6-sol","effort":"high","source_digest":"416f75ff593d9f8fc90658e8cfe6c07489cf206b","state":"completed","outcome":"pass","native_ref":{"ref":"collaboration:/root/td_c1a_spec_confirm/final","digest":"sha256:47a1db8e100de3190c28573489492ad1802d28f352951f4b81c1700b0d439fc7","excerpt":"Spec: **0 blocking findings; worst issue: none.** Finding IDs/dispositions: none. Line: gpt-6-sol / high / 1 iteration.\nI reaffirm the Spec pass for the complete TD-C1a pair `72e751b7a020a76aa35286852b4ea981c95cb865` \u2192 `7feee3e4a4714d96429719ef25d52c5dc3e9582b`, covering TD1\u2013TD6 and TD8\u2013TD15. Later chunks remain pending.\nClaim: {\"status\":\"claimed\",\"confidence\":9}"},"axis":"Spec","base":"72e751b7a020a76aa35286852b4ea981c95cb865","tip":"7feee3e4a4714d96429719ef25d52c5dc3e9582b","finding_ids":[],"supersedes":["TD-C1a-Spec-1"]},{"id":"TD-C1a-Coverage-repair-1","performer":"/root/td_c1a_coverage_confirm","role":"independent-review","model":"gpt-6-sol","effort":"high","source_digest":"416f75ff593d9f8fc90658e8cfe6c07489cf206b","state":"completed","outcome":"pass","native_ref":{"ref":"collaboration:/root/td_c1a_coverage_confirm/final","digest":"sha256:e752cf7033ee24ab00f4a7c8f31277aab1ed34a84dbf890f35f6863515effd81","excerpt":"Coverage: **0 blocking findings; worst: none.** Findings, dispositions, and finding confidences: none.\nI reaffirm the complete pair `72e751b7a020a76aa35286852b4ea981c95cb865` \u2192 `7feee3e4a4714d96429719ef25d52c5dc3e9582b`. The repair preserves the initial TD-C1a evidence for TD1\u2013TD6 and TD8\u2013TD15; later chunks remain pending.\nNo blocking missing evidence identified. The coordinator still owes the whole-project checkpoint. Command contribution: none. No edits, tests, or probes ran.\n{\"status\":\"claimed\",\"confidence\":9}"},"axis":"Coverage","base":"72e751b7a020a76aa35286852b4ea981c95cb865","tip":"7feee3e4a4714d96429719ef25d52c5dc3e9582b","finding_ids":[],"supersedes":["TD-C1a-Coverage-1"]}]}],"completion":{"state":"pending","source_digest":"","performer":"codex/test-determinism-inline-20260928","reconciliation":{},"verification":[]},"amendments":[{"from":"sha256:0c1f19eeefd300b73018e190ea2fb3acf0baf416e4171f4b8717cfa5aa651a63","to":"sha256:a0a05390b237129415afa3bc4ccab193d0e0c7d672e8ebe3c22c00b1cf4c3b77","chunk_ids":{"TD-C1a":["TD-C1a"],"TD-C1b":["TD-C1b"],"TD-C2":["TD-C2"],"TD-C3":["TD-C3"]}}]}
```
