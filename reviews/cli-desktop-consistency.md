# CLI and desktop consistency review

C1 confirmation has two findings and two repair targets.
The first review had five findings and five repair targets.
Repair cycle 1 of 2 completed its author verification and confirmation.
Repair cycle 2 is authorized in the user-selected session.
One repair cycle remains before that work starts.
The source stays unqualified until current reviews and its checkpoint pass.

## Standards

Current finding count: 1.
The full confirmation closed STD-C1-01.
STD-C1-02: auto-fix, confidence 9.
Make the system fixture consume one interface vocabulary.
Sources: AGENTS.md:35-48 and internal/systemtest/compatibility_test.go:49,86,135,156,184,192.
The confirmation read the complete frozen diff, spec, and profile.
The earlier partial-read limitation is closed.

First review:

Finding count: 1.
The worst issue is duplicated interface vocabulary.
STD-C1-01: auto-fix, confidence 9.
Derive accepted operands and both help presentations from one vocabulary.
Sources: AGENTS.md, internal/compatibility/inspect.go:17, internal/adopt/compatibility.go:21, cmd/bench/main.go:142.

The first Standards return did not read the complete frozen diff and spec.
The confirming reviewer must cover that unread scope and the repair delta.
The current finding is supported; the partial read is not a complete Standards pass.

## Spec

Current finding count: 0.
Confirmation closed both Spec findings.
The Spec axis found no remaining requirement mismatch in the repair delta.

First review:

Finding count: 2.
The worst issue is unchanged fingerprints after a policy change.
C1-SPEC-01: auto-fix, confidence 10.
Include observed policy identity in the context fingerprint.
Sources: CD13, internal/adopt/compatibility.go:116, internal/compatibility/inspect.go:120.

C1-SPEC-02: auto-fix, confidence 9.
Keep the configuration home unknown when its lookup fails.
Sources: CD03, internal/adopt/compatibility.go:135.

## Coverage

Current finding count: 1.
Confirmation closed COV-C1-01.
COV-C1-02 remains: auto-fix, confidence 10.
Exercise the real CLI collector with CODEX_HOME absent and HOME set.
Compare complete home inventories to detect added or changed cache files.
Sources: spec.md:90 and internal/systemtest/compatibility_test.go:83-93,145-169.
The Spec pass does not override this independent Coverage finding.

First review:

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
  "plan_digest": "sha256:088394145d38657fce2c8fe4b1c1f995daf2613993625c7922a93b1636ea30f4",
  "implementation_session": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
  "chunks": [
    {
      "id": "C1",
      "base": "6ea6b7e6fe86e3fee0d0fc9a1ff01b6808486a69",
      "tip": "51e6851df49f4bf6175c0925b09c3d68748e79ef",
      "plan_digest": "sha256:088394145d38657fce2c8fe4b1c1f995daf2613993625c7922a93b1636ea30f4",
      "source_digest": "f4cc0130bd36c3a62dfc52b20094585520555204",
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
        },
        {
          "id": "c1-compatibility-repair1",
          "performer": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "f4cc0130bd36c3a62dfc52b20094585520555204",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
            "digest": "sha256:7e8204315b18544c2491c4c98256a88ac1f81203e8c4096c061e5b3fd53121ee",
            "excerpt": "Native author result after repair cycle1: bench test --package ./internal/compatibility passed in4ms with no failures or skips. The repaired package bytes are unchanged at51e6851d. Policy-file identity and compatibility context share one framing encoder.\n"
          },
          "requirement": "compatibility",
          "command": "bench test --package ./internal/compatibility",
          "exit_code": 0
        },
        {
          "id": "c1-adopt-repair1",
          "performer": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "f4cc0130bd36c3a62dfc52b20094585520555204",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
            "digest": "sha256:85d94f1d9cf9dae65e03b01dbc7ecf4bbf6be44ff9e11958f5f6c6a19d38c1bd",
            "excerpt": "Native author result after repair cycle1: bench test --package ./internal/adopt passed in20764ms with no failures or skips. Before repair TestCompatibilityMissingConfigurationHome failed in3ms because missing HOME became .codex. The focused compatibility family passed in21ms after correction.\n"
          },
          "requirement": "adopt",
          "command": "bench test --package ./internal/adopt",
          "exit_code": 0
        },
        {
          "id": "c1-system-repair1",
          "performer": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "f4cc0130bd36c3a62dfc52b20094585520555204",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
            "digest": "sha256:274acadfef29f0dfebe075d24527a5d27753093126d4d6e8b98ffbefb957cfc7",
            "excerpt": "Native author result after repair cycle1: the final bench test --check system passed in67102ms with no failures or skips after both manual source restorations. It exercises real CLI/desktop identity, both configuration homes unchanged, and policy changes through the selected executable.\n"
          },
          "requirement": "system",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "c1-doctor-route-probe-repair1",
          "performer": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "f4cc0130bd36c3a62dfc52b20094585520555204",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
            "digest": "sha256:03cb0b993da6e2f97d7a33f721d20ada2e7e44b51fd8c2d2f474bb9b5c4bf8d1",
            "excerpt": "Native author result after repair cycle1: exactly Run: adoptCommand(\"doctor\") was replaced by Run: adoptCommand(\"setup\"). The sealed system suite exited1 in71105ms. TestCompatibilityMissingPath failed with code2 instead of1 and setup usage, plus other compatibility and wrapper-reload failures. Exact registry bytes and mode were restored; SHA256 df42cb91f526209d9f1f6512e0b362b63a72e63742af06f9014703d633e563cf. The restored system suite exited0 in67102ms with no skips.\n"
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
              "digest": "sha256:03cb0b993da6e2f97d7a33f721d20ada2e7e44b51fd8c2d2f474bb9b5c4bf8d1",
              "excerpt": "Native author result after repair cycle1: exactly Run: adoptCommand(\"doctor\") was replaced by Run: adoptCommand(\"setup\"). The sealed system suite exited1 in71105ms. TestCompatibilityMissingPath failed with code2 instead of1 and setup usage, plus other compatibility and wrapper-reload failures. Exact registry bytes and mode were restored; SHA256 df42cb91f526209d9f1f6512e0b362b63a72e63742af06f9014703d633e563cf. The restored system suite exited0 in67102ms with no skips.\n"
            }
          }
        },
        {
          "id": "c1-evidence-command-repair1",
          "performer": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "f4cc0130bd36c3a62dfc52b20094585520555204",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21",
            "digest": "sha256:478e5684babd09fedbb482a68ab4c008430bbb6206ff284012014fa84757f6d0",
            "excerpt": "Native author result after repair cycle1: bench test --package ./internal/preflight/evidencecmd passed in20201ms with no failures or skips. The approved predecessor-base assertion uses the shared TOON encoder and the earlier observed wrong-base probe remains recorded.\n"
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
        },
        {
          "id": "c1-standards-confirm1",
          "performer": "codex-collaboration:/root/c1_standards_confirm",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "f4cc0130bd36c3a62dfc52b20094585520555204",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "native:c1_standards_confirm",
            "digest": "sha256:c67257b110fdf053fc43f03a54fb30e0ad89242fd24963a1df2cdc573eff9e2a",
            "excerpt": "Standards confirmation: one current finding.\nSTD-C1-01 repaired: internal/compatibility/inspect.go:16-41 owns interface values/enumeration/operand/parser; internal/adopt/compatibility.go:21,29-41,78-80 derives usage and parsing; cmd/bench/main.go:142 derives root help. Exact root-help expectation cmd/bench/help_inventory_test.go:63,111 has demonstrated omission red in assets/repair-cycle-1.md:56-57.\nSTD-C1-02: held, auto-fix, confidence 9. AGENTS.md:35-48 requires fixture harnesses single-sourced. internal/systemtest/compatibility_test.go repeats CLI default at49, desktop routing at86, two-member loops at135/156, CLI policy calls184/192. An interface rename/addition can omit or misroute a process path. Consume canonical enumeration or one fixture table; retain root-help as independent omission oracle. Unit tests use constants, and help exception does not cover repeated system fixture routing. No automated check run; semantic mandatory standard.\nNo other Standards finding survived. Full frozen diff4435lines, whole spec537lines, whole profile674lines, AGENTS/BENCH/craft-review/finding discipline/bounded repair/delegate/Claim schema/synthesis/CLI/comments read. Trusted s1,s44 all279consumerrows,s45 all64coverage rows; untouched consumers first. Targeted compatibility/adopt/legacy/cmd/system/help/records/evidence read.\nBinding current=true, clean: evidence sha256:ba0f51e4c0dec746697b3076c372a975b37263c863cd86e97395831cf79dcb88; assignment f7123d5b2ad592acc6e5a239c1f1f389; base6ea6b7e6fe86e3fee0d0fc9a1ff01b6808486a69; tip8eca522fa925d2a6307cfb22144529abcac597c6. Delivery unverified; canonical export used, no cursors.\nFuture C2/C3 ungraded. No edits/tests/probes/commits/stash/poolpath calls. Optional advice none; command contribution none requested. Claim status=claimed confidence=9.\n"
          },
          "axis": "Standards",
          "base": "6ea6b7e6fe86e3fee0d0fc9a1ff01b6808486a69",
          "tip": "51e6851df49f4bf6175c0925b09c3d68748e79ef",
          "finding_ids": [
            "STD-C1-02"
          ],
          "supersedes": [
            "c1-standards-r1"
          ]
        },
        {
          "id": "c1-spec-confirm1",
          "performer": "codex-collaboration:/root/c1_spec_confirm",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "f4cc0130bd36c3a62dfc52b20094585520555204",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:c1_spec_confirm",
            "digest": "sha256:1dc3eff14586384efd27f6aca176c9040c3b678790f600faedbdd5eb78a1ab8e",
            "excerpt": "Spec confirmation: positive terminal result, zero findings, confidence10.\nCurrent binding true, clean, head8eca522fa925d2a6307cfb22144529abcac597c6; evidence sha256:ba0f51e4c0dec746697b3076c372a975b37263c863cd86e97395831cf79dcb88; assignmentf7123d5b2ad592acc6e5a239c1f1f389. Frozen chunk6ea6b7e6fe86e3fee0d0fc9a1ff01b6808486a69..8eca522fa925d2a6307cfb22144529abcac597c6; confirmation delta c219a98fcd227a7eab852e4e8bc42597039707d9..8eca522fa925d2a6307cfb22144529abcac597c6.\nRead own metadata/source-1, consumers/source-44, coverage/source-45; complete spec,ticket1,CD01..17,repair/prior verification/review/plan; single saved delta; untouched legacy dispatch/test/help consumers first; compatibility/adopt/system tests and bounds/grammar/TOON/spill/process helpers.\nConfirmed STD-C1-01 production interface owner (inspect.go:21,adopt/compatibility.go:21,main.go:142); C1-SPEC-01 content-backed policy fingerprint (spec:290,inspect.go:71,147,collector:103,systemtest:172,repair evidence:33); C1-SPEC-02 missing home unknown (spec:280,collector:134,adopt test:200,evidence:28); COV-C1-01 actual wrapper identities (spec:278,systemtest:52,133,evidence:45); COV-C1-02 actual collector config comparisons (spec:282,systemtest:145,evidence:47). This axis judged the existing comparison sufficient; independent Coverage retains the broader home-inventory gap.\nPlan amendment preserves C1 and updates CD01,03,05,13 seams. Future C2/C3 not graded.\nNo tests/probes/edits/commits/recollection. No optional advice. No implementation-command change warranted. Claim status=claimed confidence10.\n"
          },
          "axis": "Spec",
          "base": "6ea6b7e6fe86e3fee0d0fc9a1ff01b6808486a69",
          "tip": "51e6851df49f4bf6175c0925b09c3d68748e79ef",
          "finding_ids": [],
          "supersedes": [
            "c1-spec-r1"
          ]
        },
        {
          "id": "c1-coverage-confirm1",
          "performer": "codex-collaboration:/root/c1_coverage_confirm",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "f4cc0130bd36c3a62dfc52b20094585520555204",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "native:c1_coverage_confirm",
            "digest": "sha256:54d89967dc371101622664c9c1fb1c3348d830d4b63097b648b19ecb95e7b599",
            "excerpt": "Coverage confirmation: one fold confirmed; one finding remains.\nBinding: full chunk 6ea6b7e6fe86e3fee0d0fc9a1ff01b6808486a69..8eca522fa925d2a6307cfb22144529abcac597c6; repair delta c219a98fcd227a7eab852e4e8bc42597039707d9..8eca522fa925d2a6307cfb22144529abcac597c6; evidence sha256:ba0f51e4c0dec746697b3076c372a975b37263c863cd86e97395831cf79dcb88. Current binding true; clean.\nRead scope: whole spec, profile hostile-input checklist, full confirmation delta, producer/collector/tests, repair evidence, trusted s1/s7/s22/s44/s45, 279-row consumer inventory.\nCOV-C1-01 confirmed: TestCompatibilityCollectorInterfaces invokes the sealed binary via the repository wrapper for both identities (internal/systemtest/compatibility_test.go:83-143). Actual collector carries selected identity (internal/adopt/compatibility.go:93-100).\nCOV-C1-02 remains held: auto-fix, confidence 10. With CODEX_HOME absent and HOME set, the CLI selects HOME/.codex (internal/adopt/compatibility.go:134-145); every sealed fixture instead sets CODEX_HOME (internal/systemtest/compatibility_test.go:83-93). Missing-home unit test calls only configurationHome with both variables empty (internal/adopt/compatibility_test.go:200-208). Independent bypass: create or alter a sibling health-cache file while config.toml stays unchanged. TestCompatibilityCollectorReadOnly snapshots only config.toml (internal/systemtest/compatibility_test.go:145-169), so it misses this violation of spec.md:90. Add real CLI fallback invocation and recursive before/after inventories of both fixture homes.\nPolicy fingerprint test exercises both hook and selected-config changes. Vocabulary probe does not close surviving write bypass.\nOne repair cycle consumed; one remains. Future C2/C3 ungraded. No implementation command, test, probe, edit, or commit contributed.\nClaim: status=claimed, confidence=10.\n"
          },
          "axis": "Coverage",
          "base": "6ea6b7e6fe86e3fee0d0fc9a1ff01b6808486a69",
          "tip": "51e6851df49f4bf6175c0925b09c3d68748e79ef",
          "finding_ids": [
            "COV-C1-02"
          ],
          "supersedes": [
            "c1-coverage-r1"
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
      "from": "sha256:8a7ffa8720bf9e2ab3de0c7c3392ee62bdc9e868652deb7adcf02b80bb07fe85",
      "to": "sha256:088394145d38657fce2c8fe4b1c1f995daf2613993625c7922a93b1636ea30f4",
      "chunk_ids": {
        "C1": [
          "C1"
        ]
      }
    }
  ]
}
```
