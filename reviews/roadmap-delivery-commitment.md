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
  "plan_digest": "sha256:326be25226515d9d41f40d116e837f042c3d445cb05a0e99f1a5b4cc3d31eda5",
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
