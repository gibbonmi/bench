# CLI and desktop consistency review

C1 has five findings and five repair targets.
Repair cycle 1 of 2 is authorized in the user-selected session.
No repair cycle has completed yet.
The source stays unqualified until current reviews and its checkpoint pass.

## Standards

Finding count: 1.
The worst issue is duplicated interface vocabulary.
STD-C1-01: auto-fix, confidence 9.
Derive accepted operands and both help presentations from one vocabulary.
Sources: AGENTS.md, internal/compatibility/inspect.go:17, internal/adopt/compatibility.go:21, cmd/bench/main.go:142.

The first Standards return did not read the complete frozen diff and spec.
The confirming reviewer must cover that unread scope and the repair delta.
The current finding is supported; the partial read is not a complete Standards pass.

## Spec

Finding count: 2.
The worst issue is unchanged fingerprints after a policy change.
C1-SPEC-01: auto-fix, confidence 10.
Include observed policy identity in the context fingerprint.
Sources: CD13, internal/adopt/compatibility.go:116, internal/compatibility/inspect.go:120.

C1-SPEC-02: auto-fix, confidence 9.
Keep the configuration home unknown when its lookup fails.
Sources: CD03, internal/adopt/compatibility.go:135.

## Coverage

Finding count: 2.
The worst issue is an untested production configuration write.
COV-C1-01: auto-fix, confidence 10.
Exercise both selected interfaces through the real collector.
Sources: CD01, internal/adopt/compatibility_test.go:15, internal/systemtest/compatibility_test.go:38.

COV-C1-02: auto-fix, confidence 10.
Compare both configuration homes before and after production inspection.
Sources: CD05, internal/adopt/compatibility_test.go:32, internal/adopt/compatibility.go:87.

The corrected Coverage binding succeeded before its native reaffirmation.
No optional advice was returned.
No review requested a Bench command change.

```bench-review-record
{
  "version": 1,
  "spec": "specs/cli-desktop-consistency/spec.md",
  "plan_digest": "sha256:8a7ffa8720bf9e2ab3de0c7c3392ee62bdc9e868652deb7adcf02b80bb07fe85",
  "implementation_session": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
  "chunks": [
    {
      "id": "C1",
      "base": "6ea6b7e6fe86e3fee0d0fc9a1ff01b6808486a69",
      "tip": "1217c75334b173c8a2d039b876039555738beda1",
      "plan_digest": "sha256:8a7ffa8720bf9e2ab3de0c7c3392ee62bdc9e868652deb7adcf02b80bb07fe85",
      "source_digest": "422d7577bfc454c113673a6232d4151a881fbf71",
      "acceptance_rows": [
        "CD01",
        "CD02",
        "CD03",
        "CD04",
        "CD05",
        "CD06",
        "CD07",
        "CD08",
        "CD09",
        "CD10",
        "CD11",
        "CD12",
        "CD13",
        "CD14",
        "CD15",
        "CD16",
        "CD17"
      ],
      "verification": [
        {
          "id": "c1-compatibility-author",
          "performer": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "422d7577bfc454c113673a6232d4151a881fbf71",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
            "digest": "sha256:3856d224fa2fb2b321317277f865600d81a175f23b9f25c751cf7ca19abc4a9f",
            "excerpt": "Author result in Codex chat 01a0fc58-0f59-7f52-9c4c-a25247d5fa21: bench test --package ./internal/compatibility passed with no skips in 5 ms after the effective-config correction. The current chunk retains the verified package bytes. The ledger retains the row mutation results.\n"
          },
          "requirement": "compatibility",
          "command": "bench test --package ./internal/compatibility",
          "exit_code": 0
        },
        {
          "id": "c1-adopt-author",
          "performer": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "422d7577bfc454c113673a6232d4151a881fbf71",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
            "digest": "sha256:f86de95151685cd755c8e6b92b5df722aefd90d7f3295b25a94b0e625c2b67d8",
            "excerpt": "Author result in Codex chat 01a0fc58-0f59-7f52-9c4c-a25247d5fa21: bench test --package ./internal/adopt passed with no skips in 29298 ms after the canonical-payload collector correction. The current chunk retains the verified package bytes.\n"
          },
          "requirement": "adopt",
          "command": "bench test --package ./internal/adopt",
          "exit_code": 0
        },
        {
          "id": "c1-system-author",
          "performer": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "422d7577bfc454c113673a6232d4151a881fbf71",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
            "digest": "sha256:b0d398543c61313251addc5fffd37bc6a9a8958097e65e558452151aea7734e8",
            "excerpt": "Author result in Codex chat 01a0fc58-0f59-7f52-9c4c-a25247d5fa21: bench test --check system passed with no skips in 90110 ms after exact doctor-route source restoration. The current chunk retains the verified system and doctor source bytes.\n"
          },
          "requirement": "system",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "c1-doctor-route-probe-author",
          "performer": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "422d7577bfc454c113673a6232d4151a881fbf71",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
            "digest": "sha256:0c9c050abe0a0e3599f58f06d53809c519fbbbe93a52231c62c8d6e71da2b3db",
            "excerpt": "Author result in Codex chat 01a0fc58-0f59-7f52-9c4c-a25247d5fa21: the spec's manual swap changed Run: adoptCommand(\"doctor\") to Run: adoptCommand(\"setup\"). bench test --check system exited 1 in 72514 ms; TestCompatibilityMissingPath and TestCompatibilityPaths failed, with one existing wrapper-reload failure. The backup was restored and byte/mode identity verified. The restored system suite exited 0 with no skips in 90110 ms. This is the corrected probe; the earlier unquoted-path assertion failure was not accepted.\n"
          },
          "requirement": "doctor-route-probe",
          "command": "doctor-route-system-swap: follow the Doctor-route mutation procedure",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
              "digest": "sha256:0c9c050abe0a0e3599f58f06d53809c519fbbbe93a52231c62c8d6e71da2b3db",
              "excerpt": "Author result in Codex chat 01a0fc58-0f59-7f52-9c4c-a25247d5fa21: the spec's manual swap changed Run: adoptCommand(\"doctor\") to Run: adoptCommand(\"setup\"). bench test --check system exited 1 in 72514 ms; TestCompatibilityMissingPath and TestCompatibilityPaths failed, with one existing wrapper-reload failure. The backup was restored and byte/mode identity verified. The restored system suite exited 0 with no skips in 90110 ms. This is the corrected probe; the earlier unquoted-path assertion failure was not accepted.\n"
            }
          }
        },
        {
          "id": "c1-evidence-command-author",
          "performer": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "422d7577bfc454c113673a6232d4151a881fbf71",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
            "digest": "sha256:c1bc8cf142afb5d93ec141804f72ded1c4205bd408b76ea04615362161a750d2",
            "excerpt": "Native result in this chat: bench test --package ./internal/preflight/evidencecmd passed with no failures or skips in 39242 ms. Before that, the compiled wrong-base mutation in internal/diff/diff.go failed TestReviewFileReconstruction/predecessor_base at the exact-base assertion, and bench probe reported bit and restored yes. The current chunk retains those source bytes.\n"
          },
          "requirement": "evidence-command",
          "command": "bench test --package ./internal/preflight/evidencecmd",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "c1-standards-r1",
          "performer": "codex-collaboration:/root/c1_standards",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "422d7577bfc454c113673a6232d4151a881fbf71",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "codex-collaboration:/root/c1_standards",
            "digest": "sha256:f2504b22179cf082035f2faf4859d8eb68b83f6fbeab89ef70e533a497b75902",
            "excerpt": "Standards review: 1 finding. Worst issue: duplicated interface-vocabulary knowledge.\nSTD-C1-01 — Interface vocabulary has three production sources.\nRule: AGENTS.md, “Code standard — one source per fact”: “Two derivations of the same fact must collapse into one source,” explicitly including “an enforcement and its advertisement.”\nEvidence: internal/compatibility/inspect.go:17 defines and parses codex-cli / codex-desktop; internal/adopt/compatibility.go:21 independently embeds both values in diagnostic usage; cmd/bench/main.go:142 independently embeds them in the CLI advertisement. Adding or renaming an interface can therefore leave accepted values, diagnostics, and help inconsistent. Derive both presentations from the compatibility owner’s vocabulary.\nConfidence: 9/10. Disposition: auto-fix.\nLimitation: coordinator requested an immediate bounded return before I paged the entire 2,924-line saved diff, the whole 537-line spec export, and every project-profile section. Those unread portions were not used to assert findings; unsupported concerns remain uncertain.\nThe corrected evidence binding reported current=true, source71afaadb, and a clean tree.\nDone claim: status=claimed, confidence=9—the finding is source-cited, but this read-only axis ran no red-to-green check.\n"
          },
          "axis": "Standards",
          "base": "6ea6b7e6fe86e3fee0d0fc9a1ff01b6808486a69",
          "tip": "1217c75334b173c8a2d039b876039555738beda1",
          "finding_ids": [
            "STD-C1-01"
          ],
          "supersedes": []
        },
        {
          "id": "c1-spec-r1",
          "performer": "codex-collaboration:/root/c1_spec",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "422d7577bfc454c113673a6232d4151a881fbf71",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "codex-collaboration:/root/c1_spec",
            "digest": "sha256:5bbf3d653c8b3ebf44c7ccbb89a928cc7e564f57caa3c7f143c46a0a35d32768",
            "excerpt": "Spec axis, gpt-5.6-sol/high, frozen pair 6ea6b7e6fe86e3fee0d0fc9a1ff01b6808486a69..71afaadb2dd28a50e37d47e4b83d72866d92fd16. Native current binding succeeded. Two held auto-fix findings: C1-SPEC-01 (confidence10): collectCompatibility reduces hook policy to state and path, so changed valid policy bytes leave Fingerprint unchanged, violating CD13 and policy invalidation. Sources: internal/adopt/compatibility.go:116-118; internal/compatibility/inspect.go:120-126. C1-SPEC-02 (confidence9): configurationHome ignores os.UserHomeDir errors, producing relative .codex and reading unintended cwd configuration, violating unknown metadata and CD03. Source: internal/adopt/compatibility.go:135-136. CD01,CD02,CD04-CD12,CD14-CD17 source-aligned; CD18-CD64 intentionally pending. No optional advice. One frozen diff and targeted rules/spec/compiled decisions/consumers read. No tests or edits. Claim: claimed, confidence9.\n"
          },
          "axis": "Spec",
          "base": "6ea6b7e6fe86e3fee0d0fc9a1ff01b6808486a69",
          "tip": "1217c75334b173c8a2d039b876039555738beda1",
          "finding_ids": [
            "C1-SPEC-01",
            "C1-SPEC-02"
          ],
          "supersedes": []
        },
        {
          "id": "c1-coverage-r1",
          "performer": "codex-collaboration:/root/c1_coverage",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "422d7577bfc454c113673a6232d4151a881fbf71",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "codex-collaboration:/root/c1_coverage",
            "digest": "sha256:5b360a8a978237719fae8784fd1c6c8625d0af6c701ec29d2c8ed0bb1b3a8309",
            "excerpt": "Coverage axis, gpt-5.6-sol/high, frozen pair 6ea6b7e6fe86e3fee0d0fc9a1ff01b6808486a69..71afaadb2dd28a50e37d47e4b83d72866d92fd16. First current-binding call omitted sha256: and was rejected. Native same-round reaffirmation corrected that transcription, obtained current=true and clean source71afaadb, then rechecked cited implementation/tests/spec. Two findings reaffirmed unchanged: COV-C1-01 (auto-fix, claimed, confidence10): substituting CLI in the real collector would preserve CD01 unit positives because they inject a collector; real process tests invoke only CLI. Sources: internal/adopt/compatibility.go:89-93; compatibility_test.go:15-28; internal/systemtest/compatibility_test.go:38-69; spec.md:278. COV-C1-02 (auto-fix, claimed, confidence10): a config write in collectCompatibility would preserve CD05 positives because its test substitutes the collector and process tests do not compare pre/post bytes. Sources: internal/adopt/compatibility.go:87-104; compatibility_test.go:32-53; spec.md:90,282; finding-discipline.md:32-35. Worst COV-C1-02. Producer input family and empty operational write set enumerated; untouched consumers walked; all64 rows and263 consumer rows read. No optional advice or command change. No tests, probes, builds, or edits.\n"
          },
          "axis": "Coverage",
          "base": "6ea6b7e6fe86e3fee0d0fc9a1ff01b6808486a69",
          "tip": "1217c75334b173c8a2d039b876039555738beda1",
          "finding_ids": [
            "COV-C1-01",
            "COV-C1-02"
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
  }
}
```
