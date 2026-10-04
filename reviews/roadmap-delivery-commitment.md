# Commitment implementation review

The reviewer directs implementation in the current session.
The previous ticket 01 author stopped before its first commit.
The current session preserves that source and performs fresh verification.
Independent review uses three separate GPT-6.1 Sol sessions at high effort.

Repair cycles consumed: 1 of 2 for DC-C1.

## Standards

The initial review found two issues. Repair cycle 1 closes both; zero current findings remain.

- `DC-C1-Standards-S1`: Auto-fix, confidence 10. Derive grammar and help from one flag declaration, per `AGENTS.md:35` and `commitcmd/command.go:29`.
- `DC-C1-Standards-S2`: Auto-fix, confidence 10. Commit the existing probe evidence, as `.agents/commands/bench-implement-spec.md:50` requires.

## Spec

The initial review found two issues. Repair cycle 1 closes both; zero current findings remain.

- `DC-C1-Spec-S1`: Auto-fix, confidence 10. Bind affected predecessor sources, per spec lines 158 and 164 and `authority.go:48`.
- `DC-C1-Spec-S2`: Auto-fix, confidence 10. Refuse case aliases in JSON field names, per spec line 464 and `parse.go:20`.

## Coverage

The initial review found two issues, including the shared source-binding defect. Repair cycle 1 closes both; zero current findings remain.

- `COV-1`: Auto-fix, confidence 9. Add the removed-source command regression, per spec lines 158 and 164 and `repository.go:152`.
- `COV-2`: Auto-fix, confidence 10. Add a two-outcome dependency cycle test for DC56. The cycle-guard omission passed all 27 tests.

Six raw findings produce five repair targets after the shared source-binding defect is counted once.
The record uses axis-qualified names for the reviewers' overlapping S1 and S2 identifiers.
The native occurrences supersede coordinator summaries with literal native excerpts.
The source identity, findings, observations, and dispositions are unchanged.

The Coverage reviewer suggests that evidence pages expose navigation outside oversized content cells.
This suggestion is optional advice and has no repair disposition.

```bench-review-record
{
  "version": 1,
  "spec": "specs/roadmap-delivery-commitment/spec.md",
  "plan_digest": "sha256:44dc3d96f2c82e35a8c0f95cf64d9a988c14b3a3fd3df478e6b343d2e3bf5b57",
  "implementation_session": "/root",
  "chunks": [
    {
      "id": "DC-C1",
      "base": "ca339ea83ef401f5edea17a012ba6dc891ec2995",
      "tip": "3df02a0850105a852ca308134f41d822f34991a0",
      "plan_digest": "sha256:326be25226515d9d41f40d116e837f042c3d445cb05a0e99f1a5b4cc3d31eda5",
      "source_digest": "6c06ee1b5cec9a450b56aed10a2f3c80d414e53e",
      "acceptance_rows": [
        "DC1",
        "DC2",
        "DC3",
        "DC4",
        "DC5",
        "DC7",
        "DC8",
        "DC9",
        "DC55",
        "DC56",
        "DC57",
        "DC58",
        "DC59",
        "DC60",
        "DC62",
        "DC71"
      ],
      "verification": [
        {
          "id": "dc-c1-root-commitment",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "9fdad244f80cf037eda7934b2887907cafff4c27",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:root:dc-c1:commitment",
            "digest": "sha256:56a9db8201153ffe0a11b4e4ad364a61d2dbfab22db19f57a730e11fc2cb51c4",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,pass,744\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "commitment",
          "command": "bench test --package ./internal/commitment",
          "exit_code": 0
        },
        {
          "id": "dc-c1-root-intent",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "9fdad244f80cf037eda7934b2887907cafff4c27",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:root:dc-c1:intent",
            "digest": "sha256:6dc6a3fdb3c5ae187cfdb8c5e31c71b13f01c37707774ec6c2f026004cbeabd3",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/intent,pass,4335\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "intent",
          "command": "bench test --package ./internal/intent",
          "exit_code": 0
        },
        {
          "id": "dc-c1-root-roadmap",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "9fdad244f80cf037eda7934b2887907cafff4c27",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:root:dc-c1:roadmap",
            "digest": "sha256:e47fd560266a997f684e900b955674a1da981ea94ee0ca180dcda5c871827a38",
            "excerpt": "tree[1]{target,head,dirty}:\n  dc-integration,75bdb76c551f026c58369e4d92156318933d8e24,true\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/roadmap,pass,3212\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "roadmap",
          "command": "bench test --package ./internal/roadmap",
          "exit_code": 0
        },
        {
          "id": "dc-c1-root-bench",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "9fdad244f80cf037eda7934b2887907cafff4c27",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:root:dc-c1:bench",
            "digest": "sha256:4341e06a75c5240ce31d068dcf71c1be47a43fc6f0b8f1656eb66fb166e25d65",
            "excerpt": "tree[1]{target,head,dirty}:\n  dc-integration,75bdb76c551f026c58369e4d92156318933d8e24,true\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17141\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "dc-c1-r1-commitment",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "6c06ee1b5cec9a450b56aed10a2f3c80d414e53e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:root:dc-c1-r1:commitment",
            "digest": "sha256:dbb1481dec8c691146065a2cd73df023af7958d27cdddae2f011262653850713",
            "excerpt": "tree[1]{target,head,dirty}:\n  dc-integration,5d207a8865063b055206af588b3b31e900993cce,true\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,pass,661\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "commitment",
          "command": "bench test --package ./internal/commitment",
          "exit_code": 0
        },
        {
          "id": "dc-c1-r1-jsonfile",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "6c06ee1b5cec9a450b56aed10a2f3c80d414e53e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:root:dc-c1-r1:jsonfile",
            "digest": "sha256:50970728fec739ec8b87baf6328300f0225ef772877f8655957799b789cd1cf5",
            "excerpt": "tree[1]{target,head,dirty}:\n  dc-integration,5d207a8865063b055206af588b3b31e900993cce,true\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/jsonfile,pass,2\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "jsonfile",
          "command": "bench test --package ./internal/jsonfile",
          "exit_code": 0
        },
        {
          "id": "dc-c1-r1-intent",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "6c06ee1b5cec9a450b56aed10a2f3c80d414e53e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:root:dc-c1-r1:intent",
            "digest": "sha256:e8b447fe28ff5eac39a2d9837981e3d63b5ef5c74dcb5a44d947a2e9af3b9b54",
            "excerpt": "tree[1]{target,head,dirty}:\n  dc-integration,5d207a8865063b055206af588b3b31e900993cce,true\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/intent,pass,4480\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "intent",
          "command": "bench test --package ./internal/intent",
          "exit_code": 0
        },
        {
          "id": "dc-c1-r1-bench",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "6c06ee1b5cec9a450b56aed10a2f3c80d414e53e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:root:dc-c1-r1:bench",
            "digest": "sha256:c612599a3903a57d207902bb4a7a2124c21a69e52a3b741858a3bc8d48270295",
            "excerpt": "tree[1]{target,head,dirty}:\n  dc-integration,5d207a8865063b055206af588b3b31e900993cce,true\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,pass,18197\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "dc-c1-r1-roadmap",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "6c06ee1b5cec9a450b56aed10a2f3c80d414e53e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:root:dc-c1-r1:roadmap",
            "digest": "sha256:da84b81c2e2bf0ff4b8c19a1e349f217b95d0cb4bce853b1683afcc400cc458e",
            "excerpt": "tree[1]{target,head,dirty}:\n  dc-integration,5d207a8865063b055206af588b3b31e900993cce,true\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/roadmap,pass,2790\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "roadmap",
          "command": "bench test --package ./internal/roadmap",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "dc-c1-standards-initial",
          "performer": "/root/dc_c1_standards",
          "role": "independent-review",
          "model": "gpt-6.1-sol",
          "effort": "high",
          "source_digest": "9fdad244f80cf037eda7934b2887907cafff4c27",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "native:collaboration:/root/dc_c1_standards",
            "digest": "sha256:8feb019790eca18a9ec33b180c3368ceec8c8360fa5cd730d0aecca30fafb85d",
            "excerpt": "Standards: 2 findings; 2 repair targets. Worst issue: duplicated command grammar.\nDC-C1-S1 — Grammar has two sources. Confidence 10/10, disposition auto-fix. AGENTS.md:35 requires one source per fact, explicitly including enforcement and advertisement. command.go:29 independently declares flag names, value operands, and requiredness in suffix strings and usage.Flag entries; line 59 consumes those separately for help and parsing. Derive both from one declaration.\nDC-C1-S2 — Probe evidence was not committed before dispatch. Confidence 10/10, disposition auto-fix. bench-implement-spec.md:50 requires a committed verification and probe record before axis dispatch. The frozen review pickup:53 enumerates only four passing package verifications. The coordinator confirms native red/probe evidence exists outside that record. This finding does not claim that reds were never observed.\nNo command change is necessary.\ncheckout: dc-c1-standards,8a33f77423dee97e1bb18f60390448b756b92696,clean\n"
          },
          "axis": "Standards",
          "base": "ca339ea83ef401f5edea17a012ba6dc891ec2995",
          "tip": "75bdb76c551f026c58369e4d92156318933d8e24",
          "finding_ids": [
            "DC-C1-Standards-S1",
            "DC-C1-Standards-S2"
          ],
          "supersedes": []
        },
        {
          "id": "dc-c1-spec-initial",
          "performer": "/root/dc_c1_spec",
          "role": "independent-review",
          "model": "gpt-6.1-sol",
          "effort": "high",
          "source_digest": "9fdad244f80cf037eda7934b2887907cafff4c27",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "native:collaboration:/root/dc_c1_spec",
            "digest": "sha256:d1d9e1117ddc0a532561c1cd033f70ff4c6b6442f2681168f4d400f704505c38",
            "excerpt": "Spec: 2 findings; 2 repair targets. Both auto-fix, confidence 10/10.\nDC-C1-S1 — Predecessor-only sources escape approval invalidation. Spec:158 binds affected source identities; spec:164 states a changed protected source invalidates the receipt. authority.go:48 uses sources := policySources(proposed), and repository.go:152 rechecks only that list. A restored probe of TestCommitmentPublishedSourceIdentity observed approval with stale working copy = <nil>, want changed published source refusal.\nDC-C1-S2 — Policy field aliases violate exact JSON grammar. Spec:464 requires exact field names. A restored probe containing \"version\":1,\"Version\":1 observed Parse() accepted duplicate field input.\nBoth probes restored=yes; final commitment package passed. No phase-command change is necessary.\ndc-c1-spec HEAD 8a33f77423dee97e1bb18f60390448b756b92696; status empty.\n"
          },
          "axis": "Spec",
          "base": "ca339ea83ef401f5edea17a012ba6dc891ec2995",
          "tip": "75bdb76c551f026c58369e4d92156318933d8e24",
          "finding_ids": [
            "DC-C1-Spec-S1",
            "DC-C1-Spec-S2"
          ],
          "supersedes": []
        },
        {
          "id": "dc-c1-coverage-initial",
          "performer": "/root/dc_c1_coverage",
          "role": "independent-review",
          "model": "gpt-6.1-sol",
          "effort": "high",
          "source_digest": "9fdad244f80cf037eda7934b2887907cafff4c27",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "native:collaboration:/root/dc_c1_coverage",
            "digest": "sha256:3ff1fbc34ecf170211b63a7a778b82d5c3a910127aec65bfcfea4f331e1fd4e6",
            "excerpt": "Coverage: 2 findings. Worst: removal approval remains usable after its removed source changes.\nCOV-1 — Removed sources escape receipt invalidation. Confidence 9, auto-fix. Plan removal of B with source obligation.md, publish changed source without policy change, then approve --removed B succeeds exit0 changed=true. Spec.md:158,164 and ticket01:13 require affected source binding. Temporary public-command TestCoverageRemovalBindsRemovedSource failed:98ms, zero skips, approval exit0, policy unchanged=false. Source spill:1791109844811149249-7da256bc6b6d40b6.out.\nCOV-2 — DC56 never exercises multi-node cycle rejection. Confidence10, auto-fix. Existing malformed cycle is A->A and fails before dependencyCycle. Omission probe removed cycle-rejection block parse.go:106; silent, baseline pass, 27 tests ran, mutated package pass580ms, zero failures/skips, restored=yes. Native spill:1791109857649689359-488dc9b4cd420972.out.\nOptional CLI advice: expose evidence successor cursor outside oversized content cell.\ndc-c1-coverage HEAD 8a33f77423dee97e1bb18f60390448b756b92696; final status empty; temporary test removed.\n"
          },
          "axis": "Coverage",
          "base": "ca339ea83ef401f5edea17a012ba6dc891ec2995",
          "tip": "75bdb76c551f026c58369e4d92156318933d8e24",
          "finding_ids": [
            "COV-1",
            "COV-2"
          ],
          "supersedes": []
        },
        {
          "id": "dc-c1-standards-native",
          "performer": "/root/dc_c1_standards",
          "role": "independent-review",
          "model": "gpt-6.1-sol",
          "effort": "high",
          "source_digest": "9fdad244f80cf037eda7934b2887907cafff4c27",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "native:collaboration:/root/dc_c1_standards",
            "digest": "sha256:6ac5d1d2120f4db426af8c123e983e14048fed5e6052cfc75bf07c0e61c34232",
            "excerpt": "Standards: **2 findings; 2 repair targets.** Worst issue: duplicated command grammar.\n\n1. **DC-C1-S1 — Grammar has two sources.** Confidence **10/10**, disposition **auto-fix**. [AGENTS.md:35](AGENTS.md:35) requires one source per fact, explicitly including enforcement and advertisement. [command.go:29](internal/commitment/commitcmd/command.go:29) independently declares flag names, value operands, and requiredness in `suffix` strings and `usage.Flag` entries; line 59 consumes those separately for help and parsing. Derive both from one declaration. Contrary-source check: existing [recordcmd/command.go:93](internal/reviewrecord/recordcmd/command.go:93) derives its suffix and grammar from declared flags. This static ownership defect has no executable refutation; the current help conformance check verifies projection wiring rather than this duplication.\n\n2. **DC-C1-S2 — Probe evidence was not committed before dispatch.** Confidence **10/10**, disposition **auto-fix**. [bench-implement-spec.md:50](.agents/commands/bench-implement-spec.md:50) requires a committed verification **and probe** record before axis dispatch. The frozen [review pickup:53](reviews/roadmap-delivery-commitment.md:53) enumerates only four passing package verifications. The coordinator confirms native red/probe evidence exists outside that record. Retain those existing source-bound outcomes and restoration results in the pickup and commit the record; this finding does **not** claim that reds were never observed.\n"
          },
          "axis": "Standards",
          "base": "ca339ea83ef401f5edea17a012ba6dc891ec2995",
          "tip": "75bdb76c551f026c58369e4d92156318933d8e24",
          "finding_ids": [
            "DC-C1-Standards-S1",
            "DC-C1-Standards-S2"
          ],
          "supersedes": [
            "dc-c1-standards-initial"
          ]
        },
        {
          "id": "dc-c1-spec-native",
          "performer": "/root/dc_c1_spec",
          "role": "independent-review",
          "model": "gpt-6.1-sol",
          "effort": "high",
          "source_digest": "9fdad244f80cf037eda7934b2887907cafff4c27",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "native:collaboration:/root/dc_c1_spec",
            "digest": "sha256:084144ec94c032bc82b473a64fd9d93f27f2f8d6961ef989ffeda4cf4a906472",
            "excerpt": "**2 findings; 2 repair targets. Worst: DC-C1-S1 permits approval after a protected source changes.** Both are `auto-fix`, confidence **10/10**.\n\n- **DC-C1-S1 — Predecessor-only sources escape approval invalidation.** Spec:158 binds “affected source identities”; spec:164 states, “A changed proposal, protected source, or predecessor policy invalidates the receipt.” `internal/commitment/authority.go:48` uses `sources := policySources(proposed)`, and `repository/repository.go:152` rechecks only that list. Trace: current A binds `obligation.md`; the proposal retains A but removes its source binding and adds B; the published source changes without changing policy bytes; approval succeeds. A restored probe of `TestCommitmentPublishedSourceIdentity` observed `approval with stale working copy = <nil>, want changed published source refusal`. Bind and recheck the affected predecessor sources as well as proposed sources.\n\n- **DC-C1-S2 — Policy field aliases violate exact JSON grammar.** Spec:464 states, “JSON field names and identity tokens use exact grammar”; DC56 requires ambiguous and duplicate fields to refuse. `internal/commitment/parse.go:20` calls `jsonfile.DecodeDocument`, whose scanner detects literal duplicate names (`internal/jsonfile/decode.go:69`) while its decoder accepts case-insensitive aliases (`:40`). Trace: an otherwise valid policy containing `\"version\":1,\"Version\":1` is accepted. A restored probe of `TestCommitmentMalformedPolicy` observed `Parse() accepted duplicate field input`. Enforce exact field names through the existing JSON owner; an in-scope fence expansion may be needed.\n"
          },
          "axis": "Spec",
          "base": "ca339ea83ef401f5edea17a012ba6dc891ec2995",
          "tip": "75bdb76c551f026c58369e4d92156318933d8e24",
          "finding_ids": [
            "DC-C1-Spec-S1",
            "DC-C1-Spec-S2"
          ],
          "supersedes": [
            "dc-c1-spec-initial"
          ]
        },
        {
          "id": "dc-c1-coverage-native",
          "performer": "/root/dc_c1_coverage",
          "role": "independent-review",
          "model": "gpt-6.1-sol",
          "effort": "high",
          "source_digest": "9fdad244f80cf037eda7934b2887907cafff4c27",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "native:collaboration:/root/dc_c1_coverage",
            "digest": "sha256:e48c95c54b9d75e0c1eefea9e85724bb69b833284af2839c42d3ef4e79ad19e2",
            "excerpt": "Coverage: **2 findings**. Worst: removal approval remains usable after its removed source changes.\n\n1. **COV-1 — Removed sources escape receipt invalidation.** Confidence **9**, disposition **auto-fix**.\n\n   Input: M1 contains A and B; B binds `obligation.md`. Plan removal of B, then publish changed bytes to `obligation.md` without changing the policy. `approve --removed B` succeeds with exit 0 and `changed=true`, staging B’s removal.\n\n   Expected: source refusal with unchanged policy and receipt. The spec says the plan binds “affected source identities” and a changed “protected source” invalidates its receipt (`spec.md:158,164`; ticket 01:13). `BuildPlan` takes only `policySources(proposed)` (`authority.go:48`), and approval validates only that list (`repository.go:152`). A removed binding therefore disappears from the check.\n\n2. **COV-2 — DC56 never exercises multi-node cycle rejection.** Confidence **10**, disposition **auto-fix**.\n\n   Input: valid distinct outcomes with dependencies A→B→A. Expected: malformed-policy refusal before a receipt write. DC56 explicitly requires dependency cycles to refuse.\n\n   `TestCommitmentMalformedPolicy` supplies only A→A (`parse_test.go:40–42`). The earlier self-dependency guard rejects it, so it does not exercise `dependencyCycle`.\n\n   Independent probe omitted the entire cycle-rejection block at `parse.go:106`. Result: **silent**, baseline passed, **27 tests ran**, mutated package passed in **580 ms**, zero failures/skips, restored=yes. [Exact native probe evidence](/home/mgibs/.bench/responses/bench-2826441890/primary/1791109857649689359-488dc9b4cd420972.out).\n"
          },
          "axis": "Coverage",
          "base": "ca339ea83ef401f5edea17a012ba6dc891ec2995",
          "tip": "75bdb76c551f026c58369e4d92156318933d8e24",
          "finding_ids": [
            "COV-1",
            "COV-2"
          ],
          "supersedes": [
            "dc-c1-coverage-initial"
          ]
        },
        {
          "id": "dc-c1-standards-r1",
          "performer": "/root/dc_c1_r1_standards",
          "role": "independent-review",
          "model": "gpt-6.1-sol",
          "effort": "high",
          "source_digest": "6c06ee1b5cec9a450b56aed10a2f3c80d414e53e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:collaboration:/root/dc_c1_r1_standards",
            "digest": "sha256:b5b92a6b0afd0e678414cabe60051a5243b6b65d800f78701572a3097d2613c5",
            "excerpt": "Standards: **0 findings; worst: none; 0 repair targets.**\n\nThe five requested folds confirm:\n\n- **Standards-S1 resolved**, confidence 10, `no-op`: `commitcmd/command.go:28` owns flag declarations; `suffix` and `grammar` derive help and parsing from them, satisfying `AGENTS.md:35`.\n- **Standards-S2 resolved**, confidence 10, `no-op`: `reviews/roadmap-delivery-commitment.md:403` retains seven author probe outcomes, their source context, and restoration results, satisfying `.agents/commands/bench-implement-spec.md:50`.\n- **Removed-source binding confirmed:** `authority.go:48` includes predecessor sources; `repository.go:152` validates that plan set before staging. The regression preserves policy and receipt bytes.\n- **Exact fields confirmed:** `parse.go:20` calls `DecodeExactDocument`; `jsonfile/fields.go:35` derives field names from struct tags. The existing scanner remains the document and duplicate-field owner. Existing decoder entry points pass a nil schema.\n- **Two-outcome cycle confirmed:** `parse_test.go:94` supplies A→B→A. The retained repair probe records that omitting graph-cycle rejection turns this regression red.\n\nNo duplicated production grammar, JSON field inventory, scanner, or fixture harness was found in the repair.\n"
          },
          "axis": "Standards",
          "base": "ca339ea83ef401f5edea17a012ba6dc891ec2995",
          "tip": "3df02a0850105a852ca308134f41d822f34991a0",
          "finding_ids": [],
          "supersedes": [
            "dc-c1-standards-native"
          ]
        },
        {
          "id": "dc-c1-spec-r1",
          "performer": "/root/dc_c1_r1_spec",
          "role": "independent-review",
          "model": "gpt-6.1-sol",
          "effort": "high",
          "source_digest": "6c06ee1b5cec9a450b56aed10a2f3c80d414e53e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:collaboration:/root/dc_c1_r1_spec",
            "digest": "sha256:6174d4acdc6d063c5e2ee9ec36177689df74a1f93983495ca8df0dfa94044f40",
            "excerpt": "Spec: **0 findings; 0 repair targets. Worst issue: none.** Repair cycle 1 accepted for this axis.\n\n- **DC-C1-Spec-S1 closed.** Spec:158 binds affected sources; spec:164 requires changed protected sources to invalidate approval. `internal/commitment/authority.go:50` now includes predecessor sources. Approval rechecks them at `repository/repository.go:152`. The removed-source command regression passed and checks unchanged policy and receipt bytes.\n- **DC-C1-Spec-S2 closed.** Spec:464 requires exact field names. `parse.go:20` now uses the existing scanner with the typed schema; `internal/jsonfile/fields.go:47` rejects unmatched names. Nested case-alias regressions passed. `Store.Plan` parses before writing receipts.\n- The two-outcome cycle regression and existing published-source, grammar, framing, and control-input checks passed. Shared JSON package checks passed.\n\nIndependent verification: focused commitment selection passed in 304 ms; JSON package passed in 2 ms. Both reported zero failures and skips. No new scope change or repair regression found. Successor-ticket requirements remain future work.\n"
          },
          "axis": "Spec",
          "base": "ca339ea83ef401f5edea17a012ba6dc891ec2995",
          "tip": "3df02a0850105a852ca308134f41d822f34991a0",
          "finding_ids": [],
          "supersedes": [
            "dc-c1-spec-native"
          ]
        },
        {
          "id": "dc-c1-coverage-r1",
          "performer": "/root/dc_c1_r1_coverage",
          "role": "independent-review",
          "model": "gpt-6.1-sol",
          "effort": "high",
          "source_digest": "6c06ee1b5cec9a450b56aed10a2f3c80d414e53e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:collaboration:/root/dc_c1_r1_coverage",
            "digest": "sha256:dbd8af4e1e769118345eae5fcc3485fdde27817ed7275539d4018f93d6fb587e",
            "excerpt": "Coverage: **0 current findings; 0 repair targets. COV-1 and COV-2 are closed.** No blocking issue remains in the repair delta.\n\n- **COV-1 — Refuted; no-op; confidence 10.** `TestCommitmentRemovalBindsRemovedSource` passes. It changes B’s published source after planning removal, then requires approval refusal with unchanged policy and receipt bytes (`removed_source_test.go:16–59`). `BuildPlan` now includes predecessor-only bindings (`authority.go:48–55`), which approval rechecks (`repository.go:152`). An independent command test also deleted B’s source after planning; approval refused and preserved policy, roadmap, and receipt bytes.\n- **COV-2 — Refuted; no-op; confidence 10.** The distinct A→B→A fixture reaches graph-cycle rejection (`parse_test.go:94–103`). My omission of the recursive traversal at `dependencyCycle`, a different site from the author’s guard omission, returned **bit**, one failed test, zero skips, **restored=yes**. [Native evidence](/home/mgibs/.bench/responses/bench-2826441890/primary/1791111184397485760-2055ba7821a62582.out).\n- **Exact-field repair confirmed.** Independently substituted case aliases at every field occurrence in a valid fixture containing criteria, sources, dependencies, parallel grants, and delivery facts. Every command refused before receipt writes. Swapping the matcher to `strings.EqualFold` returned **bit**, four failing alias cases, zero skips, **restored=yes**. [Native evidence](/home/mgibs/.bench/responses/bench-2826441890/3369b698aea02234ecca57bb5a85d77e/1791111207992668368-f68395298a64dea1.out).\n"
          },
          "axis": "Coverage",
          "base": "ca339ea83ef401f5edea17a012ba6dc891ec2995",
          "tip": "3df02a0850105a852ca308134f41d822f34991a0",
          "finding_ids": [],
          "supersedes": [
            "dc-c1-coverage-native"
          ]
        }
      ]
    },
    {
      "id": "DC-C2",
      "base": "2225afe0312a73980267cd5278fd500dfb883e33",
      "tip": "ae1ebbce906d728e1faa41d3fd790962324265ee",
      "plan_digest": "sha256:843b9b52c5196117901caa37f3b2bfe0978ef99e77453b5a6c5fffa7a62c0274",
      "source_digest": "76648b9d99695133d13af2cebd9220e73c0f17b7",
      "acceptance_rows": [
        "DC16",
        "DC17",
        "DC18",
        "DC19",
        "DC20",
        "DC21",
        "DC22",
        "DC33",
        "DC48",
        "DC61",
        "DC66",
        "DC73",
        "DC74"
      ],
      "verification": [
        {
          "id": "dc-c2-intent-author",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "68b6b25f193cb486e8535cef8286a4992289257f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:exec:ceb03d",
            "digest": "sha256:3ede76fc0f6d85e9ad558fb10d93aa96933030d87b47e332bad0e09df02430f7",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/intent,pass,3419\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "intent",
          "command": "bench test --package ./internal/intent",
          "exit_code": 0
        },
        {
          "id": "dc-c2-commitment-author",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "68b6b25f193cb486e8535cef8286a4992289257f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:exec:ce7cb9",
            "digest": "sha256:08539c63ab69fed1415de61780da016bc799dc924031fc5029845335f0d5085a",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,pass,1997\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "commitment",
          "command": "bench test --package ./internal/commitment",
          "exit_code": 0
        },
        {
          "id": "dc-c2-r1-commitment",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "60d5bb1cbe637bbe50ff1c5e6c062b04570cd35b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:exec:12a758",
            "digest": "sha256:02f01d7cd82e08a2a2f72aa072d7f344010951a806a437181357c1086f044b77",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,pass,2486\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n\ntree[1]{target,head,dirty}:\n  dc-integration,e1afc2b9dc4a9103ee342f23748917a2a653381c,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/admission.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment,^TestCommitmentStartPublishedIdentity$,passed,3\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,fail,163\nfailures[2]{package,test,line}:\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentStartPublishedIdentity/changed,\"admission_test.go:296: stale deliverable=(commitment_admission[1]{operation,outcome}:\\\\nstart,A\\\\n,0), unchanged=false\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentStartPublishedIdentity/deleted,\"admission_test.go:296: stale deliverable=(commitment_admission[1]{operation,outcome}:\\\\nstart,A\\\\n,0), unchanged=false\"\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "commitment",
          "command": "bench test --package ./internal/commitment",
          "exit_code": 0,
          "probe": {
            "mutation": "Ignore deliverable validation errors. TestCommitmentStartPublishedIdentity must fail on the reported admission result and preserve an exact source restore.",
            "outcome": "bit",
            "exit_code": 0,
            "restore": "pass",
            "native_ref": {
              "ref": "native:exec:12a758",
              "digest": "sha256:02f01d7cd82e08a2a2f72aa072d7f344010951a806a437181357c1086f044b77",
              "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,pass,2486\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n\ntree[1]{target,head,dirty}:\n  dc-integration,e1afc2b9dc4a9103ee342f23748917a2a653381c,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/admission.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment,^TestCommitmentStartPublishedIdentity$,passed,3\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,fail,163\nfailures[2]{package,test,line}:\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentStartPublishedIdentity/changed,\"admission_test.go:296: stale deliverable=(commitment_admission[1]{operation,outcome}:\\\\nstart,A\\\\n,0), unchanged=false\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentStartPublishedIdentity/deleted,\"admission_test.go:296: stale deliverable=(commitment_admission[1]{operation,outcome}:\\\\nstart,A\\\\n,0), unchanged=false\"\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dc-c2-r1-intent",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "60d5bb1cbe637bbe50ff1c5e6c062b04570cd35b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:exec:2169dd",
            "digest": "sha256:5738960b30c78b2ac6d4a45b714e861af31289133966accedbf2513806ad235d",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/intent,pass,3777\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "intent",
          "command": "bench test --package ./internal/intent",
          "exit_code": 0
        },
        {
          "id": "dc-c2-r1-spec",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "60d5bb1cbe637bbe50ff1c5e6c062b04570cd35b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:exec:62e9fc",
            "digest": "sha256:49c0e6577b7b515bcb885e506aba6ad9a75cbab6c634f6cd0380c46bb426369c",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/spec,pass,1575\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "spec",
          "command": "bench test --package ./internal/spec",
          "exit_code": 0
        },
        {
          "id": "dc-c2-r1-landing",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "60d5bb1cbe637bbe50ff1c5e6c062b04570cd35b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:exec:966200",
            "digest": "sha256:ae2e5881fc3cafd4dd6e16bab6ab32f50ce063f41f0ba59a26a45f227ef1acf9",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/landing,pass,11883\nfailures[0]{package,test,line}:\n"
          },
          "requirement": "landing",
          "command": "bench test --package ./internal/landing",
          "exit_code": 0
        },
        {
          "id": "dc-c2-r1-bench",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "60d5bb1cbe637bbe50ff1c5e6c062b04570cd35b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:exec:56c3e8",
            "digest": "sha256:8929cc139dc8a1c604bec5894b33ff247ee013c2d48652c6d24357743ee4851e",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,pass,18328\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "dc-c2-r1-conformance",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "60d5bb1cbe637bbe50ff1c5e6c062b04570cd35b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:exec:99e657",
            "digest": "sha256:064ab2f2cfb8d81818e884262375a1494212fbbdd2e4e84176b9f0648a413a1d",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,40159\nfailures[0]{package,test,line}:\n"
          },
          "requirement": "conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "dc-c2-r2-commitment",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "76648b9d99695133d13af2cebd9220e73c0f17b7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:exec:c2-r2-commitment",
            "digest": "sha256:0415f2b4fc962208bc437d6ce4ab9cb4f1732b16a938ea0ce7760a2752672b76",
            "excerpt": "tree[1]{target,head,dirty}:\n  dc-integration,4a461e6e8c5cd7586b01da22ac025b6b68c96601,true\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,pass,2979\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/admission.go,swap,failed,4,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment,^TestCommitmentStartPublishedIdentity$,passed,5\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,fail,267\nfailures[4]{package,test,line}:\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentStartPublishedIdentity/spec/changed,\"admission_test.go:303: stale deliverable=(commitment_admission[1]{operation,outcome}:\\\\nstart,A\\\\n,0), unchanged=false\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentStartPublishedIdentity/spec/deleted,\"admission_test.go:303: stale deliverable=(commitment_admission[1]{operation,outcome}:\\\\nstart,A\\\\n,0), unchanged=false\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentStartPublishedIdentity/tickets-only/changed,\"admission_test.go:303: stale deliverable=(commitment_admission[1]{operation,outcome}:\\\\nstart,A\\\\n,0), unchanged=false\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentStartPublishedIdentity/tickets-only/deleted,\"admission_test.go:303: stale deliverable=(commitment_admission[1]{operation,outcome}:\\\\nstart,A\\\\n,0), unchanged=false\"\nskips[0]{package,test,reason}:\n\n"
          },
          "requirement": "commitment",
          "command": "bench test --package ./internal/commitment",
          "exit_code": 0,
          "probe": {
            "mutation": "Ignore deliverable validation errors. TestCommitmentStartPublishedIdentity must fail on the reported admission result and preserve an exact source restore.",
            "outcome": "bit",
            "exit_code": 0,
            "restore": "pass",
            "native_ref": {
              "ref": "native:exec:c2-r2-commitment",
              "digest": "sha256:0415f2b4fc962208bc437d6ce4ab9cb4f1732b16a938ea0ce7760a2752672b76",
              "excerpt": "tree[1]{target,head,dirty}:\n  dc-integration,4a461e6e8c5cd7586b01da22ac025b6b68c96601,true\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,pass,2979\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/admission.go,swap,failed,4,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment,^TestCommitmentStartPublishedIdentity$,passed,5\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,fail,267\nfailures[4]{package,test,line}:\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentStartPublishedIdentity/spec/changed,\"admission_test.go:303: stale deliverable=(commitment_admission[1]{operation,outcome}:\\\\nstart,A\\\\n,0), unchanged=false\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentStartPublishedIdentity/spec/deleted,\"admission_test.go:303: stale deliverable=(commitment_admission[1]{operation,outcome}:\\\\nstart,A\\\\n,0), unchanged=false\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentStartPublishedIdentity/tickets-only/changed,\"admission_test.go:303: stale deliverable=(commitment_admission[1]{operation,outcome}:\\\\nstart,A\\\\n,0), unchanged=false\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentStartPublishedIdentity/tickets-only/deleted,\"admission_test.go:303: stale deliverable=(commitment_admission[1]{operation,outcome}:\\\\nstart,A\\\\n,0), unchanged=false\"\nskips[0]{package,test,reason}:\n\n"
            }
          }
        },
        {
          "id": "dc-c2-r2-intent",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "76648b9d99695133d13af2cebd9220e73c0f17b7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:exec:c2-r2-intent",
            "digest": "sha256:8855b6156744b3f558d7997276066001562df61106c9ec634148136a44e5f911",
            "excerpt": "tree[1]{target,head,dirty}:\n  dc-integration,4a461e6e8c5cd7586b01da22ac025b6b68c96601,true\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/intent,pass,4223\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "intent",
          "command": "bench test --package ./internal/intent",
          "exit_code": 0
        },
        {
          "id": "dc-c2-r2-spec",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "76648b9d99695133d13af2cebd9220e73c0f17b7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:exec:c2-r2-spec",
            "digest": "sha256:00714d43ffe1dca417319bb323d880d461f284f087411d93764a8e3dda7fa2ef",
            "excerpt": "tree[1]{target,head,dirty}:\n  dc-integration,4a461e6e8c5cd7586b01da22ac025b6b68c96601,true\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/spec,pass,1764\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "spec",
          "command": "bench test --package ./internal/spec",
          "exit_code": 0
        },
        {
          "id": "dc-c2-r2-landing",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "76648b9d99695133d13af2cebd9220e73c0f17b7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:exec:c2-r2-landing",
            "digest": "sha256:3f8ae203ca5d2e1c598da44cd655c6bc42d6ca86113d530ef8049cdd8f8f74ce",
            "excerpt": "tree[1]{target,head,dirty}:\n  dc-integration,4a461e6e8c5cd7586b01da22ac025b6b68c96601,true\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/landing,pass,12733\nfailures[0]{package,test,line}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/descendant-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/direct-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "landing",
          "command": "bench test --package ./internal/landing",
          "exit_code": 0
        },
        {
          "id": "dc-c2-r2-bench",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "76648b9d99695133d13af2cebd9220e73c0f17b7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:exec:c2-r2-bench",
            "digest": "sha256:59f082f29de05d70e2e92581432b3dbbd3f3f6bfe184b0626e08987d5301d451",
            "excerpt": "tree[1]{target,head,dirty}:\n  dc-integration,4a461e6e8c5cd7586b01da22ac025b6b68c96601,true\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,pass,21496\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "dc-c2-r2-conformance",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "76648b9d99695133d13af2cebd9220e73c0f17b7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:exec:c2-r2-conformance",
            "digest": "sha256:d2c99296e4fb7ac801ef1fe04301c2acefe2b8e748165f08957786aeee0163aa",
            "excerpt": "tree[1]{target,head,dirty}:\n  dc-integration,4a461e6e8c5cd7586b01da22ac025b6b68c96601,true\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,41948\nfailures[0]{package,test,line}:\nskips[3]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceProseBudgetRefusesNonRegularSubjects/socket,\"capability: fifo: unix sockets unavailable on this filesystem: listen unix /tmp/FVIEI3/t/TestGuidanceProseBudgetRefusesNonRegularSubjectssocket4020550135/001/.agents/skills/bench-craft-linked/SKILL.md: bind: invalid argument\"\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/conformance,TestSkillDescriptionBudgetRefusesNonRegularSubjects/socket,\"capability: fifo: unix sockets unavailable on this filesystem: listen unix /tmp/FVIEI3/t/TestSkillDescriptionBudgetRefusesNonRegularSubjectssocket3225391931/001/.agents/skills/bench-craft-planted/SKILL.md: bind: invalid argument\"\n"
          },
          "requirement": "conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "dc-c2-r2-commitment-probe-exit",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "76648b9d99695133d13af2cebd9220e73c0f17b7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:exec:1c3d51+faab0e",
            "digest": "sha256:0415f2b4fc962208bc437d6ce4ab9cb4f1732b16a938ea0ce7760a2752672b76",
            "excerpt": "tree[1]{target,head,dirty}:\n  dc-integration,4a461e6e8c5cd7586b01da22ac025b6b68c96601,true\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,pass,2979\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/admission.go,swap,failed,4,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment,^TestCommitmentStartPublishedIdentity$,passed,5\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,fail,267\nfailures[4]{package,test,line}:\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentStartPublishedIdentity/spec/changed,\"admission_test.go:303: stale deliverable=(commitment_admission[1]{operation,outcome}:\\\\nstart,A\\\\n,0), unchanged=false\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentStartPublishedIdentity/spec/deleted,\"admission_test.go:303: stale deliverable=(commitment_admission[1]{operation,outcome}:\\\\nstart,A\\\\n,0), unchanged=false\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentStartPublishedIdentity/tickets-only/changed,\"admission_test.go:303: stale deliverable=(commitment_admission[1]{operation,outcome}:\\\\nstart,A\\\\n,0), unchanged=false\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentStartPublishedIdentity/tickets-only/deleted,\"admission_test.go:303: stale deliverable=(commitment_admission[1]{operation,outcome}:\\\\nstart,A\\\\n,0), unchanged=false\"\nskips[0]{package,test,reason}:\n\n"
          },
          "requirement": "commitment",
          "command": "bench test --package ./internal/commitment",
          "exit_code": 0,
          "probe": {
            "mutation": "Ignore deliverable validation errors. TestCommitmentStartPublishedIdentity must fail on the reported admission result and preserve an exact source restore.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "native:exec:1c3d51+faab0e",
              "digest": "sha256:0415f2b4fc962208bc437d6ce4ab9cb4f1732b16a938ea0ce7760a2752672b76",
              "excerpt": "tree[1]{target,head,dirty}:\n  dc-integration,4a461e6e8c5cd7586b01da22ac025b6b68c96601,true\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,pass,2979\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/admission.go,swap,failed,4,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment,^TestCommitmentStartPublishedIdentity$,passed,5\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commitment,fail,267\nfailures[4]{package,test,line}:\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentStartPublishedIdentity/spec/changed,\"admission_test.go:303: stale deliverable=(commitment_admission[1]{operation,outcome}:\\\\nstart,A\\\\n,0), unchanged=false\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentStartPublishedIdentity/spec/deleted,\"admission_test.go:303: stale deliverable=(commitment_admission[1]{operation,outcome}:\\\\nstart,A\\\\n,0), unchanged=false\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentStartPublishedIdentity/tickets-only/changed,\"admission_test.go:303: stale deliverable=(commitment_admission[1]{operation,outcome}:\\\\nstart,A\\\\n,0), unchanged=false\"\n  github.com/gibbonmi/bench/internal/commitment,TestCommitmentStartPublishedIdentity/tickets-only/deleted,\"admission_test.go:303: stale deliverable=(commitment_admission[1]{operation,outcome}:\\\\nstart,A\\\\n,0), unchanged=false\"\nskips[0]{package,test,reason}:\n\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "dc-c2-standards",
          "performer": "/root/dc_c2_standards",
          "role": "independent-review",
          "model": "gpt-6.1-sol",
          "effort": "high",
          "source_digest": "68b6b25f193cb486e8535cef8286a4992289257f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:collaboration:/root/dc_c2_standards",
            "digest": "sha256:48b95a985dcbbeda682e12026ec0f021cad39dabb82583cabd329cdcbdf541c3",
            "excerpt": "Standards: **0 findings**, **0 deduplicated repair targets**. Worst issue: none. No implementation-command change is necessary.\n\nReviewed one frozen diff, `2225afe0312a73980267cd5278fd500dfb883e33..66796aa93693bbc5ada27130b3b62a5a7ff541ca`, against AGENTS.md, BENCH.md, the profile, review/delegation rules, comment and CLI standards, the whole approved spec, and ticket 02. Examined admission, deliverable validation, command forms, durable records, cleanup preservation, shared fixture changes, and recorded mutation evidence. Walked relevant unchanged consumers, including the root command adapter, repository readers, assignment lifecycle, and ledger policy owners.\n\nEvidence bound once on `dc-integration`. Retrieved complete source streams: **s1 page 0; s26 pages 0–2; s27 pages 0–1**.\n\nVerification in `dc-c2-standards`, with Go parallelism limited to 2:\n\n- `bench test --package ./internal/commitment`: pass, 2675 ms; no failures or skips.\n- `bench test --package ./internal/intent`: pass, 3580 ms; no failures or skips.\n- Normal shell, repository wrapper, scratch write/read/remove, and repository-rule/skill access probes passed. Doctor inspection retained unavailable runtime and effective-configuration metadata as unknown.\n\nFinal `git status --short`: empty. HEAD remains `66796aa93693bbc5ada27130b3b62a5a7ff541ca`. No production edits, commits, mutation probes, spawned agents, or live commands/tests remain. This return supplies Standards review evidence, not a gate verdict.\n"
          },
          "axis": "Standards",
          "base": "2225afe0312a73980267cd5278fd500dfb883e33",
          "tip": "f3d907da1cdf1a423bafd3aa17937371b8b31c71",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "dc-c2-spec",
          "performer": "/root/dc_c2_spec",
          "role": "independent-review",
          "model": "gpt-6.1-sol",
          "effort": "high",
          "source_digest": "68b6b25f193cb486e8535cef8286a4992289257f",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "native:collaboration:/root/dc_c2_spec",
            "digest": "sha256:fd9c1ed1e4aaf1b4cfdf25e01001d2ba163a51e382657200a7a8bb341ea9758c",
            "excerpt": "## Spec\n\nTwo findings; worst issue: the admission producer handles the wrong deliverable types.\n\n1. **C2-S1 — Deliverable admission violates the approved type contract.** Confidence **10/10**; disposition **auto-fix**.\n   - Spec `specs/roadmap-delivery-commitment/spec.md:131`: “Each delivery binding names one approved spec or one tickets-only folder…”\n   - `internal/commitment/repository/admission.go:39` passes deliverables to `validateSources`; `repository/repository.go:217` reads them through `git.ReadTreeFile`. `internal/git/tree.go:171` permits only regular files. Every tickets-only folder therefore refuses before its identity is compared.\n   - Conversely, `internal/commitment/parse.go:173` validates deliverables as generic source paths. A policy naming `src/main.go` permits a successful start.\n   - Independent temporary command probes confirmed both results: folder start exited **1** with `missing or nonregular tree file specs/light/tickets`; ordinary source-file start exited **0**.\n   - Repair the shared deliverable validation and identity contract for approved specs and tickets-only folders. Ticket07 defers closure, while ticket02 owns this admission producer.\n\n2. **C2-S2 — Three current acceptance seams point to different files than their implemented tests.** Confidence **10/10**; disposition **auto-fix**.\n   - DC33, spec **:418**, names `TestCommitmentRetainedClaim` in `store_test.go`; actual definition: `internal/commitment/admission_test.go:191`.\n   - DC48, spec **:433**, names `TestCommitmentAdoptionRequired` in `command_test.go`; actual definition: `internal/commitment/admission_test.go:178`.\n   - DC61, spec **:446**, names `TestCommitmentConcurrentStarts` in `store_test.go`; actual definition: `internal/commitment/concurrent_test.go:15`.\n   - Update the three seam cells to their actual owners. These are evidence-map mismatches, not absent tests.\n\nRaw finding count: **2**. Deduplicated repair targets: **2**. CLI improvement contribution: **0**.\n\nAudited DC16–22, DC33, DC48, DC61, and DC66. The remaining current behavior matches the specified ordering, blockers, exact grants, shared claims, restart retention, and atomic competition. Later consumer, publication, closure, and projection requirements were excluded as charged.\n\nRead the whole spec, ticket02, the single frozen diff, relevant producer definitions, and untouched command/ledger consumers. Evidence fetched completely: **s1 initial; s26 initial plus cursors `.s.26.1` and `.s.26.2`; s27 initial plus `.s.27.1`**, all under the charged evidence identity.\n\nVerification in `dc-c2-spec`:\n\n- Temporary Spec probes: **2 executed failures**, confirming C2-S1; probe file removed.\n- `bench test --package ./internal/commitment`: **pass**, 2,085 ms; no skips.\n- `bench test --package ./internal/intent`: **pass**, 3,531 ms; no skips.\n\nFrozen tip remains `66796aa93693bbc5ada27130b3b62a5a7ff541ca`. Final `git status --porcelain=v1` is empty. No production edits, commits, delegates, or live commands remain.\n"
          },
          "axis": "Spec",
          "base": "2225afe0312a73980267cd5278fd500dfb883e33",
          "tip": "f3d907da1cdf1a423bafd3aa17937371b8b31c71",
          "finding_ids": [
            "C2-S1",
            "C2-S2"
          ],
          "supersedes": []
        },
        {
          "id": "dc-c2-coverage",
          "performer": "/root/dc_c2_coverage",
          "role": "independent-review",
          "model": "gpt-6.1-sol",
          "effort": "high",
          "source_digest": "68b6b25f193cb486e8535cef8286a4992289257f",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "native:collaboration:/root/dc_c2_coverage",
            "digest": "sha256:9b37e80e74323061b253a590daa35a81523551289ae4789e50b1c0c7e69aae84",
            "excerpt": "Coverage: **1 finding, 1 unique repair target.** Disposition: `auto-fix`, held. Confidence: **9/10**. Severity: P2.\n\n**DC-C2-COV-1 — Start-time deliverable identity lacks a behavioral regression.**\n\nThe spec requires: “The start command checks … deliverable identity” ([spec.md:190](/home/mgibs/workspace/bench/specs/roadmap-delivery-commitment/spec.md:190)). The implementation enforces this at [repository/admission.go:39](/home/mgibs/workspace/bench/internal/commitment/repository/admission.go:39).\n\nThe concrete missing state is an approved deliverable whose default-branch bytes change or whose file disappears before `commitment start`. Start must refuse and preserve the ledger bytes. [TestCommitmentStartIdentityAndAtomicRefusal:218](/home/mgibs/workspace/bench/internal/commitment/admission_test.go:218) tests a wrong request and a different deliverable path, but never changes the approved path’s content. The existing [TestCommitmentPublishedSourceIdentity:188](/home/mgibs/workspace/bench/internal/commitment/store_test.go:188) checks approval, not start.\n\nIndependent refutation:\n\n- Omitting the validation call reported `bit`, restored `yes`, but failed only because `TestCommitmentConcurrentStarts` could not reach its Git barrier. This is incidental fixture detection.\n- A swap kept the source read and changed its guard to `err != nil && source.Identity == \"\"`. Approved identities are nonempty, so source errors were ignored. **All 48 commitment tests passed:** `silent`, zero failed tests, restored `yes`.\n\nRepair target: add command-level changed/deleted approved-deliverable cases in `internal/commitment/admission_test.go`, asserting refusal and unchanged bindings/claims. Add the explicit acceptance-map coverage and ticket ownership, then demonstrate that the error-ignoring mutation turns this regression red. Ticket 02 already owns this package ([ticket:4](/home/mgibs/workspace/bench/specs/roadmap-delivery-commitment/tickets/02-admit-committed-outcomes.md:4)).\n\nThe other current rows DC16–22, DC33, DC48, DC61, and DC66 have no retained Coverage findings.\n\nRead and verification scope:\n\n- Read the whole approved spec, ticket 02, frozen diff, governing rules/profile, relevant untouched consumers, admission tests, policy producers, and ledger lifecycle/transaction definitions.\n- Enumerated inputs: default-branch policy and source identities; assignment identity/state/request; durable bindings, claims, and blockers; dependency/order/grant facts; literal command operands.\n- Fetched evidence cursors: `s1/0`, `s26/0–2`, `s27/0–1`, all through their terminal pages. Checked integration binding once: `current=true`, tip `66796aa93693bbc5ada27130b3b62a5a7ff541ca`.\n- Both probe baselines passed. The second mutant passed all 48 tests. I did not independently rerun the intent package or whole gate.\n- Command contribution: **0** new CLI improvements.\n\nFinal venue status: `git status --short` empty. Both mutations restored exactly. No production changes, commits, retained scratch files, or live tools.\n"
          },
          "axis": "Coverage",
          "base": "2225afe0312a73980267cd5278fd500dfb883e33",
          "tip": "f3d907da1cdf1a423bafd3aa17937371b8b31c71",
          "finding_ids": [
            "DC-C2-COV-1"
          ],
          "supersedes": []
        },
        {
          "id": "dc-c2-r1-standards",
          "performer": "/root/dc_c2_r1_standards",
          "role": "independent-review",
          "model": "gpt-6.1-sol",
          "effort": "high",
          "source_digest": "60d5bb1cbe637bbe50ff1c5e6c062b04570cd35b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:collaboration:/root/dc_c2_r1_standards",
            "digest": "sha256:a66eaddd2656afd586282f42a396c28623a255777966b87d49fff5a2249d93b1",
            "excerpt": "Standards: **0 findings, 0 unique repair targets**. Worst issue: none. No implementation-command change is necessary.\n\nReviewed the single repair diff `66796aa93693bbc5ada27130b3b62a5a7ff541ca..d9796c107278e8a35231d17590d34618d902a9cd`, the whole approved spec, ticket02, governing standards, comment/CLI rules, and targeted source consumers.\n\nThe repair satisfies the one-source rule: `spec.TicketsOnly` owns classification, landing forwards to that owner, and admission uses it. Source identity and identity comparison each have one production owner. The changed/deleted regression’s independent expectations have a committed behavioral red under the error-ignoring mutation. Corrected seams and DC73/DC74 amendments introduce no retained Standards issue.\n\nEvidence bound once through `dc-integration`: `current=true`. Retrieved **s1/0 and s35/0,2,3**; intentionally omitted unrelated consumer page 1. Examined untouched command, status, landing identity/resume, immutable Git reader, and source-validation consumers.\n\nIndependent checks in `dc-c2-r1-standards`, with Go parallelism 2:\n\n- Commitment package: **pass**, 2396 ms.\n- Spec package: **pass**, 1626 ms.\n- Landing `^TestTicketsOnly`: **pass**, 26 ms.\n- All three reported zero failures and zero skips.\n- Normal shell, repository wrapper, scratch write/read/remove, and rules/skill access passed. Unavailable runtime/configuration observations remain unknown.\n\nFinal `git status --short` is empty. HEAD remains `d9796c107278e8a35231d17590d34618d902a9cd`. No production edits, commits, mutation probes, delegates, retained scratch files, or live tools remain. This is confirming Standards evidence, not a gate verdict.\n"
          },
          "axis": "Standards",
          "base": "2225afe0312a73980267cd5278fd500dfb883e33",
          "tip": "67413d275ba84dfbb930187c456b04895ea6f538",
          "finding_ids": [],
          "supersedes": [
            "dc-c2-standards"
          ]
        },
        {
          "id": "dc-c2-r1-spec-review",
          "performer": "/root/dc_c2_r1_spec",
          "role": "independent-review",
          "model": "gpt-6.1-sol",
          "effort": "high",
          "source_digest": "60d5bb1cbe637bbe50ff1c5e6c062b04570cd35b",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "native:collaboration:/root/dc_c2_r1_spec",
            "digest": "sha256:f12e14c51cc07410fa08d29024b81e463306d99101fdbdb5ceedbbbedf7e778c",
            "excerpt": "## Spec\n\n**1 blocking finding; 1 unique repair target.** Worst issue: approved specs with literal `.md` directory names cannot start.\n\n**C2-R1-S1 — Deliverable classification normalizes a literal folder name.** Confidence **10/10**; disposition **auto-fix**, held.\n\n- Spec `specs/roadmap-delivery-commitment/spec.md:131` requires: “Each delivery binding names one approved spec or one tickets-only folder…”\n- New `internal/commitment/repository/sources.go:44` compares the path with `LiveSpecPath(LiveSpecSlug(path))`. For `specs/release.md/spec.md`, this produces `specs/release/spec.md`.\n- `internal/spec/spec.go:45–46` explicitly reserves the `.md` trim for CLI arguments. Existing `internal/spec/spec_test.go:308–312` preserves a literal `.md` directory in spec facts.\n- Independent temporary command probes confirmed that an approved staged `specs/release.md/spec.md` refuses both planning and start. The equivalent `specs/release/spec.md` passes. Start reports exit **1** and “not an approved spec or tickets-only folder.”\n- Repair the path classification to preserve the literal parent name, and extend the deliverable regression through planning and start.\n\nC2-S2’s DC33/DC48/DC61 seam paths are corrected. DC73 and its changed/deleted cases close the previous identity regression; the committed ignored-error probe records behavioral failures and restoration. The shared tickets-only classifier, folder tree identity, landing delegates, ownership expansion, and later ticket07 contract match the approved repair.\n\nRaw findings: **1**. Unique repair targets: **1**. Optional advice: **0**. Implementation-command contribution: **0**; no command change is necessary.\n\nRead the whole approved spec and ticket02, one repair diff `66796aa93693bbc5ada27130b3b62a5a7ff541ca..d9796c107278e8a35231d17590d34618d902a9cd`, governing rules/profile, producer definitions, recorded repair evidence, and relevant untouched command, landing, resume, status, and spec consumers. Evidence bound once through `dc-integration`. Retrieved **s1/0, s35/0–1, s36/0–1**; s35 retrieval was intentionally limited to relevant consumer rows.\n\nVerification in `dc-c2-r1-spec`, Go parallelism **2**:\n\n- Temporary literal-folder probe reproduced planning and start refusals; scratch file removed.\n- Focused commitment checks for published identity, deliverable types, retained claims, missing adoption, and concurrent starts: **pass**, 731 ms; zero failures/skips.\n\nFinal `git status --porcelain=v1`: **empty**. HEAD remains `d9796c107278e8a35231d17590d34618d902a9cd`. No production edits, commits, spawned agents, retained scratch files, or live tools/tests remain.\n"
          },
          "axis": "Spec",
          "base": "2225afe0312a73980267cd5278fd500dfb883e33",
          "tip": "67413d275ba84dfbb930187c456b04895ea6f538",
          "finding_ids": [
            "C2-R1-S1"
          ],
          "supersedes": [
            "dc-c2-spec"
          ]
        },
        {
          "id": "dc-c2-r1-coverage-review",
          "performer": "/root/dc_c2_r1_coverage",
          "role": "independent-review",
          "model": "gpt-6.1-sol",
          "effort": "high",
          "source_digest": "60d5bb1cbe637bbe50ff1c5e6c062b04570cd35b",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "native:collaboration:/root/dc_c2_r1_coverage",
            "digest": "sha256:bd8a87ca1bf1e1f63176623d3b831aacf2395ebdbf368294abf35f88a03254d8",
            "excerpt": "Coverage: **1 finding, 1 unique repair target.** Worst issue: changed tickets-only tree identity has no behavioral regression. CLI improvement contribution: **0**.\n\n**DC-C2-COV-2 — DC73 omits changed approved tickets-only folders.** Confidence **10/10**; disposition **auto-fix**, held.\n\nThe repair introduces committed folder identities at `internal/commitment/repository/sources.go:25`. Ticket 02 explicitly requires “exact identities for files and committed folder trees” (`tickets/02-admit-committed-outcomes.md:13`). DC73 requires changed approved deliverables to refuse without runtime writes (`spec.md:462`).\n\nThe concrete missing input is an approved `specs/light` folder whose ticket content changes on the default branch before start. Start must refuse and preserve the ledger bytes.\n\n- `TestCommitmentStartPublishedIdentity`, `admission_test.go:266`, changes/deletes only `specs/A/spec.md`, through `admissionRepo` and `deliverable`.\n- `TestCommitmentDeliverableTypes`, `deliverable_test.go:20`, starts an unchanged tickets-only folder. Its approved tree identity is captured at **:36**, and no descendant changes before the command at **:50**.\n- The complete tree-identity test inventory search found this folder fixture as the only `git-tree:` expectation in the commitment tests. Existing approval identity tests use regular files.\n\nIndependent bypass at `repository/sources.go:91`:\n\n```go\nif got != source.Identity && !strings.HasPrefix(got, \"git-tree:\") {\n```\n\nThis preserves regular-file enforcement while ignoring changed folder identities. The full commitment package remained green:\n\n```text\nprobe: silent; failed_tests=0; restored=yes\nselection: all; baseline=passed; ran=57\ncommitment package: pass, 2415 ms\nfailures=0; skips=0\n```\n\nRepair target: extend DC73’s command regression with a changed approved tickets-only descendant; assert refusal and unchanged ledger bytes. Demonstrate that the tree-specific bypass becomes a behavioral red. This falls within ticket 02’s existing `internal/commitment` fence (`ticket:4`).\n\nThe original regular-file COV-1 repair is verified: ignoring admission validation errors produced **two behavioral failures**, `changed` and `deleted`, with `restored=yes`; baseline ran three tests and passed.\n\n```text\nclaims[2]{id,status,confidence}:\n  DC-C2-COV-1-file-regression,verified,10\n  DC-C2-COV-2,verified,10\n```\n\nReviewed one repair diff, `66796aa93693bbc5ada27130b3b62a5a7ff541ca..d9796c107278e8a35231d17590d34618d902a9cd`, the whole spec, ticket 02, producer/type/identity owners, command adapters, fixture owners, and existing classifier tests. Enumerated regular spec blobs, tickets-only committed trees, invalid file/folder types, changed/deleted sources, and the approved write fence. Later ticket 03–10 behavior was excluded.\n\nEvidence bound current **once** on integration. Complete evidence pages fetched: **s1/0, s35/0–3, s36/0–1**, through terminal pages.\n\nFinal venue `git status --porcelain=v1`: empty. HEAD remains `d9796c107278e8a35231d17590d34618d902a9cd`. Probes restored exactly. No retained scratch files, production edits, commits, spawned agents, or live commands/tests remain.\n"
          },
          "axis": "Coverage",
          "base": "2225afe0312a73980267cd5278fd500dfb883e33",
          "tip": "67413d275ba84dfbb930187c456b04895ea6f538",
          "finding_ids": [
            "DC-C2-COV-2"
          ],
          "supersedes": [
            "dc-c2-coverage"
          ]
        },
        {
          "id": "dc-c2-r2-standards-review",
          "performer": "/root/dc_c2_r2_standards",
          "role": "independent-review",
          "model": "gpt-6.1-sol",
          "effort": "high",
          "source_digest": "76648b9d99695133d13af2cebd9220e73c0f17b7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:collaboration:/root/dc_c2_r2_standards",
            "digest": "sha256:e02d7dfb9eab8b2e713501e64fe4edb438cf624a96047afc21babc00bc3eb9ab",
            "excerpt": "Standards: **pass — 0 raw findings, 0 unique repair targets**. Worst issue: none. Optional advice: 0. Repair dispositions: none. Implementation-command contribution: 0.\n\n```text\nclaims[1]{id,status,confidence}:\n  dc-c2-r2-standards,claimed,9\n```\n\nReviewed the single repair delta `d9796c107278e8a35231d17590d34618d902a9cd..607f9db7f1684028e6b84a4d0e5b9c2647fbb640`, the whole approved spec, ticket 02, governing rules/profile, review/delegate discipline, and comment/CLI/gate standards.\n\nThe repair satisfies the one-source rule:\n\n- `internal/spec/resolve.go:119–125` shares the literal folder layout between construction and recognition; admission consumes it at `repository/sources.go:44`.\n- `deliverable_test.go:62` supplies the common fixture builder to both test families.\n- `admission_test.go:266–303` covers changed/deleted file and folder members. The retained second ticket preserves the folder for the deleted-member case.\n- The committed repair record retains behavioral reds and exact restoration for the independent expectations. Moved comments describe the current contract.\n\nEvidence bound **once** through `dc-integration`: `current=true`. Complete evidence pages fetched and read: **s1/0, s36/0–3, s37/0–1**. Relevant untouched path consumers in status, roadmap, preflight, retirement, landing refusal, and resume were inspected; no cross-chunk concern invalidates earlier evidence.\n\nIndependent checks in `dc-c2-r2-standards`, Go parallelism 2:\n\n- Deliverable type and published-identity tests: **pass**, 846 ms.\n- Spec package: **pass**, 1671 ms.\n- Both report zero failures and zero skips.\n\nNormal shell, repository wrapper, rules/skill access, and scratch write/read/remove passed. Unavailable runtime observations remain unknown.\n\nFinal `git status --porcelain=v1` is empty. HEAD remains `607f9db7f1684028e6b84a4d0e5b9c2647fbb640`. No production edits, commits, mutation probes, delegates, retained scratch files, or live tools/tests remain. This confirms Standards; it is not a gate verdict.\n"
          },
          "axis": "Standards",
          "base": "2225afe0312a73980267cd5278fd500dfb883e33",
          "tip": "ae1ebbce906d728e1faa41d3fd790962324265ee",
          "finding_ids": [],
          "supersedes": [
            "dc-c2-r1-standards"
          ]
        },
        {
          "id": "dc-c2-r2-spec-review",
          "performer": "/root/dc_c2_r2_spec",
          "role": "independent-review",
          "model": "gpt-6.1-sol",
          "effort": "high",
          "source_digest": "76648b9d99695133d13af2cebd9220e73c0f17b7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:collaboration:/root/dc_c2_r2_spec",
            "digest": "sha256:06e08623c5f0de19053712402c369b4fc6d7bb88de3159fd9e6b8a4ffaaef5d0",
            "excerpt": "Spec: **0 findings; 0 unique repair targets.** Worst issue: none. Optional advice: 0. Implementation-command contribution: 0.\n\nBoth repair targets conform to the approved requirements:\n\n- **C2-R1-S1:** Spec `:131` permits approved specs. `repository/sources.go:44` now uses `spec.IsLiveSpecPath`; `resolve.go:124` preserves literal folder names through the existing layout owner. `deliverable_test.go:21,33,46` checks planning and start for `specs/release.md/spec.md`.\n- **DC-C2-COV-2:** DC73 (`spec.md:462`) requires changed or deleted deliverables to refuse without runtime writes. `admission_test.go:277` retains two folder members; `:289–303` changes or deletes one, invokes start, and checks refusal plus identical ledger bytes.\n\nReviewed the whole approved spec, ticket02, governing rules/profile, the single repair delta `d9796c1..607f9db`, changed owners, relevant untouched consumers, and committed verification/probe records. No cross-chunk concern invalidated earlier evidence. Tickets03–10 remain outside this checkpoint.\n\nEvidence bound current once on `dc-integration`: **true**. Retrieved **s1/0, s36/0–3, s37/0–1**, including terminal pages; inspected relevant consumer rows.\n\nIndependent checks in `dc-c2-r2-spec`, Go parallelism 2:\n\n- Focused commitment identity/type/persistence/adoption/concurrency tests: **pass**, 1062 ms, zero failures/skips.\n- Spec path and facts tests: **pass**, 11 ms, zero failures/skips.\n- Normal shell, wrapper, scratch write/read/remove, rules and skill access passed. Runtime identity remains unknown.\n\nNo mutation probe was independently executed this round; committed author probe evidence was inspected.\n\n```text\nclaims[1]{id,status,confidence}:\n  DC-C2-final-Spec-conformance,claimed,9\n```\n\nFinal status is empty; HEAD remains `607f9db7f1684028e6b84a4d0e5b9c2647fbb640`. No production edits, commits, delegates, retained scratch files, or live tools/tests remain.\n"
          },
          "axis": "Spec",
          "base": "2225afe0312a73980267cd5278fd500dfb883e33",
          "tip": "ae1ebbce906d728e1faa41d3fd790962324265ee",
          "finding_ids": [],
          "supersedes": [
            "dc-c2-r1-spec-review"
          ]
        },
        {
          "id": "dc-c2-r2-coverage-review",
          "performer": "/root/dc_c2_r2_coverage",
          "role": "independent-review",
          "model": "gpt-6.1-sol",
          "effort": "high",
          "source_digest": "76648b9d99695133d13af2cebd9220e73c0f17b7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:collaboration:/root/dc_c2_r2_coverage",
            "digest": "sha256:87cb3669a1880c45a6fe9ea7d59cd52838f54d8426a687437e511b0edea40ce5",
            "excerpt": "Coverage: **0 findings, 0 unique repair targets**. Worst issue: none. Optional advice: 0. Implementation-command contribution: 0.\n\nReviewed only `d9796c107278e8a35231d17590d34618d902a9cd..607f9db7f1684028e6b84a4d0e5b9c2647fbb640`, against the whole approved spec and ticket 02. Later-ticket behavior remains outside this round.\n\nThe repaired input family covers staged spec blobs, literal `.md` spec directories, committed tickets-only trees, invalid source files, unstaged specs, and folders containing specs. DC73 now exercises both changed and deleted file/folder members. Its deleted-folder-member fixture retains a second ticket, ensuring identity comparison remains reachable (`internal/commitment/admission_test.go:266`). Both refusal cases assert identical ledger bytes. The shared fixture builder supplies committed identities (`internal/commitment/deliverable_test.go:62`). Changes remain within ticket 02’s commitment/spec fences and the authorized review pickup.\n\nIndependent bypass attempt: in `repository.sourceAt`, changed the tree read from `revision+\":\"+path` to `revision+\"^:\"+path`. This preserved unchanged deliverables while admitting stale folder content through a different production site from the author’s comparison bypass.\n\n```text\nprobe: bit\nsubject: internal/commitment/repository/sources.go\nmutation: swap\nbaseline: passed, ran=12\nfailures:\n  TestCommitmentStartPublishedIdentity/tickets-only/changed\n  TestCommitmentStartPublishedIdentity/tickets-only/deleted\nrestored: yes\nskips: 0\n```\n\nThe literal-directory case checks planning and start through the production command (`TestCommitmentDeliverableTypes`, `deliverable_test.go:15`). The moved helpers retain existing CLI normalization coverage and literal-directory fact coverage (`internal/spec/spec_test.go:308`, `:332`, `:342`).\n\nIndependent checks, serialized with Go parallelism 2:\n\n- Commitment package: pass, 2779 ms; zero failures/skips.\n- Spec package: pass, 1610 ms; zero failures/skips.\n\n```text\nclaims[2]{id,status,confidence}:\n  DC-C2-COV-2-repair,verified,10\n  literal-folder-command-regression,claimed,10\n```\n\nEvidence bound current once through `dc-integration`. Complete pages fetched: **s1/0, s36/0–3, s37/0–1**, through terminal pages. Read governing rules, Session compatibility, project profile, review/delegate/gate instructions, complete spec/ticket, the single repair diff, fixture and identity owners, and relevant untouched path-helper consumers.\n\nFinal venue status is empty; HEAD remains `607f9db7f1684028e6b84a4d0e5b9c2647fbb640`. Probe restored exactly. No retained scratch files, production edits, commits, spawned agents, or live commands/tests remain.\n"
          },
          "axis": "Coverage",
          "base": "2225afe0312a73980267cd5278fd500dfb883e33",
          "tip": "ae1ebbce906d728e1faa41d3fd790962324265ee",
          "finding_ids": [],
          "supersedes": [
            "dc-c2-r1-coverage-review"
          ]
        }
      ]
    },
    {
      "id": "DC-C3",
      "base": "975be583753195efcc790927fc8019b0e18367ff",
      "tip": "ca141db2e1e6813ebc6b0ac7c522fb855e52951b",
      "plan_digest": "sha256:44dc3d96f2c82e35a8c0f95cf64d9a988c14b3a3fd3df478e6b343d2e3bf5b57",
      "source_digest": "fe4d7c5ea584850cb933745cef9b4a4c5406309c",
      "acceptance_rows": [
        "DC6",
        "DC10",
        "DC11",
        "DC13",
        "DC15",
        "DC26",
        "DC27",
        "DC31"
      ],
      "verification": [
        {
          "id": "c3-author-commit",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "fe4d7c5ea584850cb933745cef9b4a4c5406309c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:c3-commit",
            "digest": "sha256:db6c6ae5feed9bb8dcfd13419ffa2f8cc903eb90e0941d2dc12bf9630ca0bfeb",
            "excerpt": "native 46a429: internal/commit pass 11240ms; failures 0; skips 0.\n"
          },
          "requirement": "commit",
          "command": "bench test --package ./internal/commit",
          "exit_code": 0
        },
        {
          "id": "c3-author-preflight",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "fe4d7c5ea584850cb933745cef9b4a4c5406309c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:c3-preflight",
            "digest": "sha256:43cd2dbad756a35dabe8df5515a8cd45b3a1e862272d03a2a4cce5a4718f85aa",
            "excerpt": "native 745aaf: internal/preflight pass 33666ms; evidencecmd pass 24695ms; chargesource and preflighttest no-tests; failures 0; skips 0.\n"
          },
          "requirement": "preflight",
          "command": "bench test --package ./internal/preflight",
          "exit_code": 0
        },
        {
          "id": "c3-author-roadmap",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "fe4d7c5ea584850cb933745cef9b4a4c5406309c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:c3-roadmap",
            "digest": "sha256:ee0725717d52f99816fe8638641f578d2edebd14673971a4be28ea8e9447617c",
            "excerpt": "native 583002: internal/roadmap pass 2732ms; failures 0; skips 0.\n"
          },
          "requirement": "roadmap",
          "command": "bench test --package ./internal/roadmap",
          "exit_code": 0
        },
        {
          "id": "c3-author-commitment",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "fe4d7c5ea584850cb933745cef9b4a4c5406309c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:c3-commitment",
            "digest": "sha256:7a405230284a367bb4e70a1ba114c728dde45bc45c4216ac043058af3ded7999",
            "excerpt": "native 632600: internal/commitment pass 3866ms; child packages compiled; failures 0; skips 0. Later candidate change only adds an error remedy and removes an unused parameter; final commit and system runs exercise that source.\n"
          },
          "requirement": "commitment",
          "command": "bench test --package ./internal/commitment",
          "exit_code": 0
        },
        {
          "id": "c3-author-charge-evidence-system",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "fe4d7c5ea584850cb933745cef9b4a4c5406309c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:c3-charge-evidence-system",
            "digest": "sha256:d45384d33eb20ff724d6ae1b4afedcb601537852774fe8d8708d779266d7defd",
            "excerpt": "native 29e968: internal/systemtest pass 78515ms; failures 0; skips 0.\n"
          },
          "requirement": "charge-evidence-system",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "c3-author-landing",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "fe4d7c5ea584850cb933745cef9b4a4c5406309c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:c3-landing",
            "digest": "sha256:6b25a3d16f4c0b14a758a5f374b636a922094204316040363a77616f70b710d9",
            "excerpt": "native 2ac0f6: internal/landing pass 12382ms; failures 0; capability skips 2 (character device); environment skips 0.\n"
          },
          "requirement": "landing",
          "command": "bench test --package ./internal/landing",
          "exit_code": 0
        },
        {
          "id": "c3-author-bench",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "fe4d7c5ea584850cb933745cef9b4a4c5406309c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:c3-bench",
            "digest": "sha256:941b96113f7ff79df26bb4f3b2968f7dfa9c3fedcdd1ad4fba875b7f074fe4a2",
            "excerpt": "native 6ec016: cmd/bench pass 19047ms; failures 0; skips 0. Later plan-only bounded-response selection is covered by final evidencecmd run 745aaf.\n"
          },
          "requirement": "bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "c3-author-conformance",
          "performer": "/root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "fe4d7c5ea584850cb933745cef9b4a4c5406309c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:c3-conformance",
            "digest": "sha256:294dec6c251cc5a9041eb6fded9d6240a9ca810f363777cf8f15c8544d1e4336",
            "excerpt": "native c1559b: internal/conformance pass 38971ms; failures 0; capability skips 3 (two long socket paths, one character device); environment skips 0.\n"
          },
          "requirement": "conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        }
      ],
      "reviews": []
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
      "from": "sha256:90d07b9062b87212aeafb5aa3d16850c1107cffd9762b127f7e1c25b38a3f331",
      "to": "sha256:326be25226515d9d41f40d116e837f042c3d445cb05a0e99f1a5b4cc3d31eda5",
      "chunk_ids": {
        "DC-C1": [
          "DC-C1"
        ]
      }
    },
    {
      "from": "sha256:326be25226515d9d41f40d116e837f042c3d445cb05a0e99f1a5b4cc3d31eda5",
      "to": "sha256:949bb0ba6bde7b19f76aaf51a4c0aeda8f20264281f6281f21a07d924868a5f2",
      "chunk_ids": {
        "DC-C1": [
          "DC-C1"
        ],
        "DC-C2": [
          "DC-C2"
        ]
      }
    },
    {
      "from": "sha256:949bb0ba6bde7b19f76aaf51a4c0aeda8f20264281f6281f21a07d924868a5f2",
      "to": "sha256:843b9b52c5196117901caa37f3b2bfe0978ef99e77453b5a6c5fffa7a62c0274",
      "chunk_ids": {
        "DC-C1": [
          "DC-C1"
        ],
        "DC-C2": [
          "DC-C2"
        ]
      }
    },
    {
      "from": "sha256:843b9b52c5196117901caa37f3b2bfe0978ef99e77453b5a6c5fffa7a62c0274",
      "to": "sha256:44dc3d96f2c82e35a8c0f95cf64d9a988c14b3a3fd3df478e6b343d2e3bf5b57",
      "chunk_ids": {
        "DC-C1": [
          "DC-C1"
        ],
        "DC-C2": [
          "DC-C2"
        ]
      }
    }
  ]
}
```

## Initial author probes

The current session ran these probes before ticket 01 committed as `75bdb76`.
Each probe returned `bit` and `restored=yes`.
The seven native excerpts below retain those observed outcomes.
Their run tip was `184c7ee`, with the ticket implementation in the working tree.
The committed-source package checks above passed after formatting.

```text
probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:
  bit,internal/commitment/commitcmd/command.go,swap,failed,1,yes
failures[1]{package,test,line}:
  github.com/gibbonmi/bench/internal/commitment,TestCommitmentUnsafeInput,"command_test.go:142: Command(plan \"/tmp/UJWKMH/t/TestCommitmentUnsafeInput3156862455/001/live.json\") = (\"commitment_plan[1]{id,predecessor,proposal}:\\\\n  \\\\\"sha256:2dbc22e0b525db8bf98117337a12791afdec85f787eb0458f373d507687d4e03\\\\\",\\\\\"sha256:81ab4ea3944f8cb056adacff73ac14c235e9ae5182162436a3fc5c501c65b9cd\\\\\",\\\\\"s
```

```text
probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:
  bit,internal/commitment/repository/repository.go,omit,failed,1,yes
failures[1]{package,test,line}:
  github.com/gibbonmi/bench/internal/commitment,TestCommitmentPublishedSourceIdentity,"store_test.go:204: approval with stale working copy = <nil>, want changed published source refusal"
```

```text
probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:
  bit,internal/commitment/repository/repository.go,swap,failed,1,yes
failures[1]{package,test,line}:
  github.com/gibbonmi/bench/internal/commitment,TestCommitmentUnrelatedAdvance,"store_test.go:135: commitment approval refused: predecessor changed from sha256:bd1b1d1dff9f4e4b164062709ab307bd31ee6528698e8bf1659b6f4d4085dd2c to sha256:bd1b1d1dff9f4e4b164062709ab307bd31ee6528698e8bf1659b6f4d4085dd2c"
```

```text
probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:
  bit,internal/commitment/commitcmd/command.go,swap,failed,1,yes
failures[1]{package,test,line}:
  github.com/gibbonmi/bench/internal/commitment,TestCommitmentLiteralInput,"command_test.go:51: Command(plan literal) = (\"error: bench commitment plan refused — input \\\\\"/tmp/EWJGJ6/t/TestCommitmentLiteralInput2249645367/001/proposal [*].json\\\\\" is absent: \\\\n\", 1), want plan"
```

```text
probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:
  bit,internal/commitment/commitcmd/command.go,swap,failed,1,yes
failures[1]{package,test,line}:
  github.com/gibbonmi/bench/internal/commitment,TestCommitmentGrammar,"command_test.go:81: Command(plan --input approve) = (\"error: bench commitment approve refused — use an owned planning worktree: bench worktree create --request <request> --label <label>\\\\n\", 1), want plan refusal"
```

```text
probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:
  bit,internal/commitment/repository/repository.go,swap,failed,1,yes
failures[1]{package,test,line}:
  github.com/gibbonmi/bench/internal/commitment,TestCommitmentControlInput,"command_test.go:107: Command(approve control) = (\"commitment_approval[1]{plan,decision,changed}:\\\\n  \\\\\"sha256:9e3d5e35da6d1386796177646ff95cbc30b1d769535285a02a17830f556b6757\\\\\",\\\\\"bad\\\\\\\\nreference\\\\\",\\\\\"true\\\\\"\\\\n\\\\n\", 0), want refusal"
```

```text
probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:
  bit,internal/commitment/commitcmd/command.go,swap,failed,1,yes
failures[1]{package,test,line}:
  github.com/gibbonmi/bench/internal/commitment,TestCommitmentInputFraming,"command_test.go:73: framed plan = (\"commitment_plan[1]{id,predecessor,proposal}:\\\\n  \\\\\"sha256:9e3d5e35da6d1386796177646ff95cbc30b1d769535285a02a17830f556b6757\\\\\",\\\\\"sha256:bd1b1d1dff9f4e4b164062709ab307bd31ee6528698e8bf1659b6f4d4085dd2c\\\\\",\\\\\"sha256:2609e8f85fbbcd7bcf43c0de2cf98c39004a9ec7ef1600067cafcb12ff8fd7c6
```

## Repair cycle 1

The repair source is `3df02a0850105a852ca308134f41d822f34991a0`.
Current verification passes for commitment, intent, roadmap, the CLI, and JSON.
The conformance package also passes, including its root conformance test.
Three existing capability cases skip: two Unix socket paths and one character device.
A root-package selection matched no tests. It supplies no verification credit.

The source defect came from omission of predecessor-only bindings in the plan.
The alternative causes were a wrong published revision and a stale receipt lookup.
The command regression changes only the protected source, which isolates the omitted binding.
The plan now binds the union of current and proposed sources.
The refusal preserves both policy and receipt bytes.

The field-name defect came from the decoder's case-insensitive field matching.
The alternatives were literal duplicate scanning and an invalid fixture.
The literal duplicate case already refused, while a valid uppercase alias was accepted.
The existing JSON scanner now takes an optional typed schema for exact field matching.
The standard decoder entry points retain their existing behavior.

The command form now derives its help and grammar from one flag declaration.
The new dependency test uses two distinct outcomes, which reaches the graph-cycle guard.
Independent confirmation closes all five repair targets with zero current findings.

### Repair probes

Each baseline passed. Each mutation returned `bit`, and each source restored exactly.
These runs used the working tree that committed as the repair source above.

| target | mutation | test | observed red |
| --- | --- | --- | --- |
| Removed source | Omit predecessor-source binding from `BuildPlan` | `TestCommitmentRemovalBindsRemovedSource` | Approval succeeded and changed policy bytes. |
| Exact fields | Replace the exact decoder with the existing document decoder | `TestCommitmentExactFieldNames` | Four case-alias inputs were accepted. |
| DC56 cycle | Omit graph-cycle rejection | `TestCommitmentMultiOutcomeCycle` | The two-outcome cycle was accepted. |

## Confirmed DC-C1 source

Three fresh GPT-6.1 Sol reviewers at high effort accepted repair cycle 1.
Standards, Spec, and Coverage each report zero findings and zero repair targets.
The exact source remains `3df02a0850105a852ca308134f41d822f34991a0`.
The confirming review pair ends at its record commit, `d619163d2107790c4aea024fd8d73033484557e6`.
The full chunk checkpoint remains required before ticket 02.

Coverage also checked deleted predecessor sources and aliases at every field occurrence in a complete policy fixture.
The temporary checks passed and left a clean tree.
Its independent recursive-cycle omission and field-matcher swap both bit and restored exactly.
The native Coverage excerpt retains those results.

## DC-C2 author probes

The initial admission tests failed because the three command forms were absent.
The current admission and intent packages pass. The CLI package also passes.
The root conformance check passes after the test waits use the shared bound.

The race check pauses B after its runtime snapshot. It then races an unblock and A start against B.
A split check and write lets both starts succeed. The test refuses that result.

The separate help expectation detects an omitted command form.
These native probe excerpts record both failures and their exact restores.

```text
probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:
  bit,internal/commitment/repository/admission.go,swap,failed,1,yes
selection[1]{form,target,run,baseline,ran}:
  package,./internal/commitment,^TestCommitmentConcurrentStarts$,passed,1
```

```text
probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:
  bit,internal/commitment/commitcmd/command.go,omit,failed,1,yes
selection[1]{form,target,run,baseline,ran}:
  package,./cmd/bench,^TestHelpInventoryIsComplete$,passed,1
spilled{lines=9,bytes=23905,omitted_lines=0,cut_lines=1,path=/home/mgibs/.bench/responses/bench-2826441890/af49ac1c59888c026ae64548dc90b756/1791112556490443470-9b8be7a4d000ec81.out}
```


## DC-C2 repair pickup

The first round has three raw findings and three repair targets. All are accepted for auto-fix.
The current session performs the repair under the user override. Repair cycles consumed: 0 of 2.

### Standards

Zero findings. No repair target remains on this axis.

### Spec

Two findings. The worst issue is the deliverable type contract.

- C2-S1: Auto-fix, confidence 10. Admit approved specs and tickets-only folders, per spec line 131. The current source reader accepts only files.
- C2-S2: Auto-fix, confidence 10. Correct the DC33, DC48, and DC61 seam paths to their test owners.

### Coverage

One finding. The worst issue is the absent source-identity regression.

- DC-C2-COV-1: Auto-fix, confidence 9. Add changed and deleted deliverable cases for the identity check at admission.go:39, per spec line 190.

The native return records contain the complete citations and probe results.
The repair adds explicit acceptance coverage and reuses the existing tickets-only classifier.

## DC-C2 repair verification

Repair cycles consumed: 1 of 2 for DC-C2.

The spec owner now supplies the tickets-only classifier. The landing API delegates to that owner.
Admission accepts staged specs and tickets-only folders. It binds file bytes or the complete committed folder tree.
The three acceptance seam paths now name their test owners. DC73 and DC74 add the missing coverage.

The new type tests first reproduced the folder refusal and ordinary-file admission.
A further case reproduced admission of an unstaged spec. All three now pass with the repaired producer.
The changed and deleted source cases pass, and ignoring validation errors makes both fail on admission behavior.

Commitment, intent, spec, landing, CLI, and conformance package checks pass.
Landing reports two existing privilege skips. Conformance reports two socket-path skips and one privilege skip.
No environment skip supplies verification credit.

These native excerpts record the current probes and their exact restores.

```text
probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:
  bit,internal/commitment/repository/admission.go,swap,failed,2,yes
selection[1]{form,target,run,baseline,ran}:
  package,./internal/commitment,^TestCommitmentStartPublishedIdentity$,passed,3
```

```text
probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:
  bit,internal/commitment/repository/admission.go,swap,failed,1,yes
selection[1]{form,target,run,baseline,ran}:
  package,./internal/commitment,^TestCommitmentConcurrentStarts$,passed,1
```

```text
probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:
  bit,internal/commitment/repository/sources.go,swap,failed,1,yes
selection[1]{form,target,run,baseline,ran}:
  package,./internal/commitment,^TestCommitmentDeliverableTypes/source-file$,passed,2
```


## DC-C2 final repair pickup

The confirming round has two raw findings and two repair targets. Both are accepted for auto-fix.
Standards has zero findings. The prior three targets are closed. Repair cycles consumed: 1 of 2.

- C2-R1-S1: Spec, confidence 10. Preserve literal spec folder names ending in `.md`, per the approved deliverable contract.
- DC-C2-COV-2: Coverage, confidence 10. Extend DC73 to changed tickets-only descendants and prove that a tree-identity bypass fails.

The current session performs repair cycle 2 under the recorded user override.
The native records retain the source citations and independent probes.

## DC-C2 final repair verification

Repair cycles consumed: 2 of 2 for DC-C2.
The literal spec path uses the spec owner's folder layout without CLI normalization.
The command regression first reproduced both planning and start refusals for a directory ending in `.md`.
Both commands now accept that exact path.

DC73 now changes or deletes a member of an approved tickets-only tree.
A second ticket keeps the deleted-member folder present, so this case tests its changed tree identity.
The common fixture builder supplies the approved file or folder binding for both type and identity tests.

Commitment, intent, spec, landing, CLI, and conformance package checks pass.
Landing reports two existing privilege skips. Conformance reports two socket-path skips and one privilege skip.
No environment skip supplies verification credit.

Each probe baseline passed. Each mutation failed on behavior and restored exactly.
The table records current results from the native probe outputs.

| target | mutation | failed tests | restore |
| --- | --- | --- | --- |
| Published identity | Ignore deliverable validation errors | Four changed or deleted file and folder cases | Exact |
| Folder identity | Bypass only the Git-tree identity comparison | Two changed or deleted folder-member cases | Exact |
| Literal spec path | Restore the CLI-normalization comparison | Planning and start in the literal-folder case | Exact |
| Atomic admission | Split the runtime check from its write | Concurrent starts both succeeded | Exact |

The native verification entries retain the current package results and required source-error probe.
The independent test expectations detect the named omissions above.

The commit lane refused growth in the existing oversized spec file.
Its path helpers now share the existing resolver file, within the directory file budget.
All six package checks and all four probes passed again on the final layout.

## Confirmed DC-C2 source

Three fresh GPT-6.1 Sol reviewers at high effort accepted repair cycle 2.
Standards, Spec, and Coverage each report zero findings and zero repair targets.
The exact source is `ae1ebbce906d728e1faa41d3fd790962324265ee`.
The confirming pair ends at its record commit, `607f9db7f1684028e6b84a4d0e5b9c2647fbb640`.

Coverage independently read the predecessor tree identity instead of the current tree.
Both changed-folder cases failed on admission behavior, and the probe restored exactly.
All earlier C2 findings are closed. The full chunk checkpoint remains required before ticket 03.

The first checkpoint stopped at evidence validation before running the gate.
The prior probe record used the wrapper exit of zero. The current entry records the mutated test exit of one.
The retained native failures, successful baseline, and exact restore are unchanged.
This metadata correction changes no source, finding, or diagnostic verdict.

## DC-C3 author mutation evidence

The three author probes caught their intended failures. Each probe restored its source exactly.
These supplemental probes do not replace the eight planned checks.

```text
native bea16c: probe bit; baseline passed; inner test exit 1; TestCommitmentCommitBeforeEffects fails in both cases when candidate authorization errors are ignored; unbound commit returns 0 and formats Go; restored yes; wrapper exit 0.
```

```text
native 228685: probe bit; baseline passed; inner test exit 1; TestCommitmentOccurrenceUpdate fails when the stripped label changes from Occurrences to Unknown; restored yes; wrapper exit 0.
```

```text
native 1566ca: probe bit; baseline passed; inner test exit 1; TestHelpInventoryIsComplete fails when --plan-only becomes --plan-draft; restored yes; wrapper exit 0. This proves the independent help expectation detects a public grammar mutation.
```

The initial review starts with zero repair cycles consumed. The chunk permits two repair cycles.
