# Session context cleanup review

## Chunk CL-C1

The frozen pair is base `4e98e581083562e284ba20802d168dbc575fd321` and tip `4a1ca2928776035d04be0ca123010065d924b80c`.
Three Opus/medium axes ran on 2026-09-12. Each ran in a separate native context and a separate read-only venue.
Raw findings: Standards 5, Spec 3, Coverage 2.

De-duplicated repair targets: 4 auto-fix (ST1, ST2, ST3, COV-1).
Two further findings (SPEC-1, SPEC-2) are ask-user. They await the reviewer's decision before any disposition.

The blast walk found no unlisted-consumer defect. `resolveAssignment` is shared by `exec.go:87`, `build.go:86`, and `path.go:50`. Each of these reads the changed identity resolver unchanged.

## Standards

Findings: 5. Worst issue: medium.

- ST1 (medium, auto-fix): `staleSetPlan`'s doc comment claims "every set mode shares it." `staleUnclaimedPlans` still builds its own refusal row, so the claim is false. Citations at tip 4a1ca292: clean_set.go:340-342; clean_unclaimed.go:96-98.
- ST2 (auto-fix): the new constant `unapplicableFingerprint = "none"` claims to match the spelling the invalid-invocation row already uses. That row still writes the literal `"none"`, not the constant. Citations at tip 4a1ca292: clean_set.go:193; worktree.go:247.
- ST3 (auto-fix): a test comment narrates provenance in past tense: "The baseline below was captured from the unedited tree." The comment-register rule requires timeless present. Citation at tip 4a1ca292: clean_set_test.go:263-266.
- ST4 (no-op, judgment): `cleanupRows` and `cleanupRowFields` hand-parse the producer's own rendered table. They do not derive the expectation through the producer's own call. The fixtures are controlled, and the constraint is stated inline. Citation at tip 4a1ca292: clean_set_test.go:22-43.
- ST5 (no-op, judgment): `cleanSelection` is destructured back into six locals right after parsing, at its one consumer. This is a Lazy Element smell baseline. Citation at tip 4a1ca292: worktree.go:270-271.

## Spec

Findings: 3. Worst issue: ask-user — a spec-prose versus code contradiction on a veto surface.

- SPEC-1 (ask-user): the rendered apply command names members by canonical assignment id, not the typed operand (clean_set.go:313-322 appends `--target <row.assignment.ID>`). Spec line 76 promises the operand is safely quoted in that command, which reads as echoing the operand text. The code substitutes a safer value instead. `axi.KnownArgument` already shell-quotes independently, so no injection risk exists either way. This is a spec-prose mismatch, not a safety gap, and spec lines 75-76 actually describe the deferred CL9 re-plan action, not this apply command. Citations: spec.md:76; clean_set.go:313-322; internal/axi/action_test.go:101.
- SPEC-2 (ask-user): the diff silently resolves an open reviewer question. The spec's "Flagged additions" section says repeated explicit targets are a proposed scope clarification, with no recorded reviewer answer. `clean_set.go:218-221` already collapses repeats silently, without that answer. Citations: spec.md, section "Flagged additions"; clean_set.go:218-221.
- SPEC-3 (no-op): CL16's re-plan-output half is vacuously covered. A hostile operand can never resolve, so it never reaches a rendered command. The test asserts the action's absence, not neutralization of hostile text under render. Citation: clean_set_command_test.go:520.

All 10 claimed rows are delivered:

- CL1 `TestCleanExplicitSetPlan` (clean_set_test.go:601).
- CL2 `TestCleanExplicitSetAliases` (:644).
- CL3 `TestCleanExplicitSetSelectionFailure` (:682).
- CL10 `TestCleanSetRetainsAuthority` (:726).
- CL11 `TestCleanSetCompatibility` (:810).
- CL13 `TestCleanSetPresentEmptyInventory` (:885).
- CL14 `TestCleanSetGrammar` (clean_set_command_test.go:446).
- CL15 is confirmed by inspection: clean_set.go:82-91 imports no budget or measurement package.
- CL16 `TestCleanSetHostileOperand` (:493).
- CL17 `TestCleanSetAbsentInventory` (:899).

The fence is clean. All 7 changed files sit in the ticket's `Writes:` list and the spec's Ownership fences. No `Won't handle` line is violated.

## Coverage

Findings: 2. Worst issue: auto-fix — a modifier path has no biting assertion.

- COV-1 (auto-fix): discard-modifier effects through the explicit-set path are untested. The spec's edge inventory puts branch- and ignored-preservation cases under CL10. The existing tests exercise only default options, not the discard modifiers. Concrete break: dropping `--discard-branch` from `cleanupModifierFlags` changes the rendered command and the set digest, yet every current test stays green. Missing rows: a set apply with `--discard-branch` asserting the branch is gone, and one with `--discard-ignored` asserting the ignored file is discarded. Citations: spec.md:203; clean_set.go:326-338; clean_set_command_test.go:55.
- COV-2 (no-op): CL16's renderer clause is unreachable in its own test. Every hostile operand fails resolution before `renderExplicitSet` reaches the help block. The renderer is safe by construction: `applyArguments` emits `row.assignment.ID`, never operand text. The space case is independently covered through the alias test's own operands. Citations: clean_set.go:240-249; clean_set_test.go:105-107; clean_set_command_test.go:60-107.

The suite is verified green. `bench test --package ./internal/worktree --run TestClean` passes in 5.1 seconds, with one unrelated host-capability skip. All nine promised tests exist: clean_set_test.go:59,102,140,184,268,343,357 and clean_set_command_test.go:13,60. The unscoped-call fixture in `TestCleanSetGrammar` was confirmed non-trivial. Relaxing its `modes == 0` guard to `modes > 1` lets each spelling fall through and lose its required exit-2 refusal.

## Author verification

The retained continuation session re-ran verification at tip 4a1ca292. The closed session's own pre-commit runs were not preserved.

- `bench test --package ./internal/worktree --run TestClean`: pass, 5086 ms, one unrelated unix-socket host-capability skip.
- `bench test --package ./cmd/bench`: pass, 8106 ms.
- Mutation probe (completion plan CL-C1, clean-tests): bypassed the alias-collapse guard at `clean_set.go:139`. `TestCleanExplicitSetAliases` failed as expected; the probe bit. The revert used `git checkout -- internal/worktree/clean_set.go`. The test passed again in 671 ms, and the tree returned to clean.

## Reaffirmation at 1b40ca08

The retained author repaired all six findings in one commit on the integration source.
The reviewer decided SPEC-1 and SPEC-2 before this repair. Keep the id-substitution behavior and fix the spec prose. Confirm silent target collapse and close the flagged addition.

- ST1, ST2, ST3: closed. Each doc comment now states its scope or spelling accurately, with no behavior change.
- COV-1: closed. `TestCleanSetDiscardModifiers` exists at clean_set_command_test.go and exercises both `--discard-branch` and `--discard-ignored` through an explicit-set apply, asserting the branch and the ignored residue. Its own mutation probe (guard `--discard-branch` out of `cleanupModifierFlags`) bit and was restored.
- SPEC-1, SPEC-2: closed. spec.md:76 now names canonical assignment identity as the rendering rule, citing the apply command as the existing precedent, and does not redefine CL9. The "Flagged additions" section records the reviewer's 2026-09-12 confirmation.

The fresh Standards pass found one new nit: a doc comment the repair added narrated in past tense, the same register ST3 corrected. The retained author fixed the one word (`carried` to `carries`) in a follow-up commit, `1b40ca08`, on top of the repair commit `518910e5`.
All three axes then re-read the chunk pair fresh, in the same independent venues, and each returned zero findings. The pair `4a1ca292`..`1b40ca08` is now the chunk's frozen tip.

Author verification at the tip: `bench test --package ./internal/worktree --run TestClean` passes (5930 ms, the same pre-existing host-capability skip). `bench test --package ./cmd/bench` passes (8655 ms, 0 skips). `bench test --package ./internal/worktree --run TestCleanSetDiscardModifiers` passes (348 ms, per the Coverage axis). `bench gate-prose` on the spec passes. `bench preflight review` reports all 13 checks green for the repair pair.

## Record
```bench-review-record
{
  "version": 1,
  "spec": "specs/session-context-cleanup/spec.md",
  "plan_digest": "sha256:acdac92b1c935f8a91069e59869d663e6094070c148cddb0916a0c74cb86cf44",
  "implementation_session": "claude:sonnet:retained-continuation",
  "chunks": [
    {
      "id": "CL-C1",
      "base": "4a1ca2928776035d04be0ca123010065d924b80c",
      "tip": "1b40ca08c6c7b723b9b60882bd8cc4271f6f836d",
      "plan_digest": "sha256:acdac92b1c935f8a91069e59869d663e6094070c148cddb0916a0c74cb86cf44",
      "source_digest": "7ce06d236a1c3e19c6f6b82eb195afb7c18c30aa",
      "acceptance_rows": [
        "CL1", "CL2", "CL3", "CL10", "CL11", "CL13", "CL14", "CL15", "CL16", "CL17"
      ],
      "verification": [
        {
          "id": "cl-c1-clean-tests",
          "performer": "claude:sonnet:retained-continuation",
          "role": "author-verification",
          "model": "sonnet",
          "effort": "unknown",
          "source_digest": "753a91da5c3b9606a733fda7c67222611ea7cb3f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:coordinator/session-context-cleanup/clean-tests@4a1ca292",
            "digest": "sha256:acfd07e4b516de9d713682922eec6b725d5d03c3e9fb76869114764d9890b489",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,pass,5086\nfailures[0]{package,test,line}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanLandedSpecialPathsRetainedWithoutOpening/socket,\"unix sockets unavailable (host-capability skip)\"\n"
          },
          "requirement": "clean-tests",
          "command": "bench test --package ./internal/worktree --run TestClean",
          "exit_code": 0,
          "probe": {
            "mutation": "bypass the alias-collapse guard (selected[assignment.ID]) in planExplicitSet, clean_set.go:139",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:coordinator/session-context-cleanup/probe-cl1@4a1ca292",
              "digest": "sha256:91e2959c23ab7260cee8ede2da95b18b932039477264ee458110dd5f162c6c57",
              "excerpt": "mutation: bypass the alias-collapse guard (selected[assignment.ID]) in planExplicitSet, clean_set.go:139\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,fail,304\nfailures[1]{package,test,line}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanExplicitSetAliases,\"clean_set_test.go:115: alias plan rows duplicated\"\nrestore: git checkout -- internal/worktree/clean_set.go; TestCleanExplicitSetAliases pass, 671ms\n"
            }
          }
        },
        {
          "id": "cl-c1-command-tests",
          "performer": "claude:sonnet:retained-continuation",
          "role": "author-verification",
          "model": "sonnet",
          "effort": "unknown",
          "source_digest": "753a91da5c3b9606a733fda7c67222611ea7cb3f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:coordinator/session-context-cleanup/command-tests@4a1ca292",
            "digest": "sha256:a15efd6d06a8d37b75e0276fd4dcae1e9913f1f07a667e4d589fadde1f3d51f6",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,pass,8106\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "command-tests",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "cl-c1-repair-clean-tests",
          "performer": "claude:sonnet:retained-continuation",
          "role": "author-verification",
          "model": "sonnet",
          "effort": "unknown",
          "source_digest": "0559c89ccdb90a6ac5ffa7044285dd19b3103d9a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:coordinator/session-context-cleanup/clean-tests@518910e5",
            "digest": "sha256:dfb62eadf49b1d2f88789dd569639dd06511806fc972ea1193ae3d8b482969d7",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,pass,5930\nfailures[0]{package,test,line}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanLandedSpecialPathsRetainedWithoutOpening/socket,\"unix sockets unavailable (host-capability skip)\"\n"
          },
          "requirement": "clean-tests",
          "command": "bench test --package ./internal/worktree --run TestClean",
          "exit_code": 0
        },
        {
          "id": "cl-c1-repair-command-tests",
          "performer": "claude:sonnet:retained-continuation",
          "role": "author-verification",
          "model": "sonnet",
          "effort": "unknown",
          "source_digest": "0559c89ccdb90a6ac5ffa7044285dd19b3103d9a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:coordinator/session-context-cleanup/command-tests@518910e5",
            "digest": "sha256:bf2b0271312752dbd762cd4a8dcc39e8c79a157009e25a4483c643711e237e19",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,pass,8655\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "command-tests",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "cl-c1-final-clean-tests",
          "performer": "claude:sonnet:retained-continuation",
          "role": "author-verification",
          "model": "sonnet",
          "effort": "unknown",
          "source_digest": "7ce06d236a1c3e19c6f6b82eb195afb7c18c30aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:coordinator/session-context-cleanup/clean-tests@1b40ca08",
            "digest": "sha256:6b6896b616de396d2b2c3eef9e94a9b3850af4c676b63e73fd142731484c2cd9",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,pass,5039\nfailures[0]{package,test,line}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanLandedSpecialPathsRetainedWithoutOpening/socket,\"unix sockets unavailable (host-capability skip)\"\n"
          },
          "requirement": "clean-tests",
          "command": "bench test --package ./internal/worktree --run TestClean",
          "exit_code": 0,
          "probe": {
            "mutation": "omit the alias collapse so a repeated identity plans twice",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:coordinator/session-context-cleanup/probe-cl1@1b40ca08",
              "digest": "sha256:929016aaa1debccc5dccb254769d0e2dd2bbefbc18c57dc45d56cf328ae065d5",
              "excerpt": "mutation: bypass the alias-collapse guard (selected[assignment.ID]) in planExplicitSet, clean_set.go:139\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,fail,4723\nfailures[1]{package,test,line}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanExplicitSetAliases,\"clean_set_test.go:115: alias plan rows duplicated\"\nrestore: git checkout -- internal/worktree/clean_set.go; TestClean pass, 5039ms\n"
            }
          }
        },
        {
          "id": "cl-c1-final-command-tests",
          "performer": "claude:sonnet:retained-continuation",
          "role": "author-verification",
          "model": "sonnet",
          "effort": "unknown",
          "source_digest": "7ce06d236a1c3e19c6f6b82eb195afb7c18c30aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:coordinator/session-context-cleanup/command-tests@1b40ca08",
            "digest": "sha256:64d2e145a263f94c40c406a53c652c5201bfa12fe36eb02de8bff25766b85430",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,pass,9393\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "command-tests",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "cl-c1-standards-1",
          "performer": "claude:opus-medium:standards-axis",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "753a91da5c3b9606a733fda7c67222611ea7cb3f",
          "state": "completed",
          "outcome": "findings",
          "native_ref": {
            "ref": "claude:agent/cl-c1-standards-1",
            "digest": "sha256:e59fabbec0dfa6f8c796974b75d172baaafe2bbc24979ea1e4feb3fa5978e72d",
            "excerpt": "result: completed; axis: Standards; findings: 5; worst: medium (staleSetPlan doc comment false universal); tip: 4a1ca2928776035d04be0ca123010065d924b80c"
          },
          "axis": "Standards",
          "base": "4e98e581083562e284ba20802d168dbc575fd321",
          "tip": "4a1ca2928776035d04be0ca123010065d924b80c",
          "finding_ids": ["ST1", "ST2", "ST3", "ST4", "ST5"],
          "supersedes": []
        },
        {
          "id": "cl-c1-standards-2",
          "performer": "claude:opus-medium:standards-axis",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7ce06d236a1c3e19c6f6b82eb195afb7c18c30aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cl-c1-standards-2",
            "digest": "sha256:2c707b09dd4e61eb30a1faa55fc8fa04141d645f4fa164129d31d85f4563a18d",
            "excerpt": "result: completed; axis: Standards; findings: 0; worst: none; tip: 1b40ca08c6c7b723b9b60882bd8cc4271f6f836d"
          },
          "axis": "Standards",
          "base": "4a1ca2928776035d04be0ca123010065d924b80c",
          "tip": "1b40ca08c6c7b723b9b60882bd8cc4271f6f836d",
          "finding_ids": [],
          "supersedes": ["cl-c1-standards-1"]
        },
        {
          "id": "cl-c1-spec-1",
          "performer": "claude:opus-medium:spec-axis",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "753a91da5c3b9606a733fda7c67222611ea7cb3f",
          "state": "completed",
          "outcome": "findings",
          "native_ref": {
            "ref": "claude:agent/cl-c1-spec-1",
            "digest": "sha256:cec5a14fd6836fd946bb2639b823d91b19fca11456bdc18b18b5cb4df2bc826f",
            "excerpt": "result: completed; axis: Spec; findings: 3; worst: ask-user (rendered apply command substitutes canonical id for typed operand); tip: 4a1ca2928776035d04be0ca123010065d924b80c"
          },
          "axis": "Spec",
          "base": "4e98e581083562e284ba20802d168dbc575fd321",
          "tip": "4a1ca2928776035d04be0ca123010065d924b80c",
          "finding_ids": ["SPEC-1", "SPEC-2", "SPEC-3"],
          "supersedes": []
        },
        {
          "id": "cl-c1-spec-2",
          "performer": "claude:opus-medium:spec-axis",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7ce06d236a1c3e19c6f6b82eb195afb7c18c30aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cl-c1-spec-2",
            "digest": "sha256:c56778b5a3709529ee32d3c3e59a32ffbaa389c4215bad46eab5c11e0c48650f",
            "excerpt": "result: completed; axis: Spec; findings: 0; worst: none; tip: 1b40ca08c6c7b723b9b60882bd8cc4271f6f836d"
          },
          "axis": "Spec",
          "base": "4a1ca2928776035d04be0ca123010065d924b80c",
          "tip": "1b40ca08c6c7b723b9b60882bd8cc4271f6f836d",
          "finding_ids": [],
          "supersedes": ["cl-c1-spec-1"]
        },
        {
          "id": "cl-c1-coverage-1",
          "performer": "claude:opus-medium:coverage-axis",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "753a91da5c3b9606a733fda7c67222611ea7cb3f",
          "state": "completed",
          "outcome": "findings",
          "native_ref": {
            "ref": "claude:agent/cl-c1-coverage-1",
            "digest": "sha256:32ab4af11f30e61a2ac6f9556ad5e452a5f09cd917d805d761a4b632c775a2a7",
            "excerpt": "result: completed; axis: Coverage; findings: 2; worst: auto-fix (discard-modifier effects through the set path are untested); tip: 4a1ca2928776035d04be0ca123010065d924b80c"
          },
          "axis": "Coverage",
          "base": "4e98e581083562e284ba20802d168dbc575fd321",
          "tip": "4a1ca2928776035d04be0ca123010065d924b80c",
          "finding_ids": ["COV-1", "COV-2"],
          "supersedes": []
        },
        {
          "id": "cl-c1-coverage-2",
          "performer": "claude:opus-medium:coverage-axis",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7ce06d236a1c3e19c6f6b82eb195afb7c18c30aa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cl-c1-coverage-2",
            "digest": "sha256:b556c96e50766b34ef6cec90b9feb59782dc73413efeac029942e02cd1943a45",
            "excerpt": "result: completed; axis: Coverage; findings: 0; worst: none; tip: 1b40ca08c6c7b723b9b60882bd8cc4271f6f836d"
          },
          "axis": "Coverage",
          "base": "4a1ca2928776035d04be0ca123010065d924b80c",
          "tip": "1b40ca08c6c7b723b9b60882bd8cc4271f6f836d",
          "finding_ids": [],
          "supersedes": ["cl-c1-coverage-1"]
        }
      ]
    }
  ]
}
```

## Chunk CL-C2

The frozen pair is base `12bbd56b349e55af711deb55cb075dadaf9c082e` and tip `0fa5fb2988b5c08b86051480cef4d234e1fee5cf`.
Three Opus/medium axes ran on 2026-09-12. Each axis ran in a separate native context and a separate venue.
Raw findings: Standards 6, Spec 4, Coverage 4.

De-duplicated repair targets: 9 auto-fix (ST6, ST7, ST8, ST9, ST10, ST11, SPEC-4, COV-5, COV-6).
One auto-fix target (COV-3) repairs the test, and its spec-row half stays ask-user.
Two further findings (SPEC-5, COV-4) are ask-user. They await the reviewer's decision before any disposition.

No axis returned a blocking finding. The Spec axis traced rows CL4 to CL9, CL12, CL18, and CL19 clean.
The Spec axis also closed the CL-C1 flag SPEC-1. The reviewer amended spec line 76, and `targetSelectors` is now the one source both the apply action and the re-plan action read.
The Spec axis confirms that both coordinator fence expansions stay inside the approved behavior.

## Standards at 0fa5fb29

Findings: 6. Worst issue: medium.

- ST6 (medium, auto-fix): three call sites each compose the stale refusal by hand. `renderOutcomes` already owns that composition. The deleted `renderLandedStale` was the landed mode's named owner, and the delta inlined its body. Citations at tip 0fa5fb29: clean_set_apply.go:104; worktree.go:323; clean_set.go:322.
- ST7 (auto-fix): explicit mode spells the row requalification twice. Landed mode collapses the same fact into `requalifyLandedRow`, and `preflightLandedSet` reuses it. Explicit mode has no symmetric owner. Citations at tip 0fa5fb29: clean_set_apply.go:117-122; clean_set.go:291-293; clean_landed.go:314-326.
- ST8 (auto-fix): the `unstarted` closure is the same code in both modes. Only the row element type differs. Citations at tip 0fa5fb29: clean_set.go:277; clean_landed.go:333.
- ST9 (auto-fix): the new file opens a second package doc comment for package `worktree`. No blank line separates it from the package clause. Citations at tip 0fa5fb29: clean_set_apply.go:1-5; worktree.go:1.
- ST10 (auto-fix): `ActionNotAttempted` is declared outside the package that owns the action vocabulary. The `lifecyclepolicy` predicates cannot see this member. Citations at tip 0fa5fb29: clean_set_apply.go:19; lifecyclepolicy.go:24-51; classifier.go:146-147.
- ST11 (auto-fix): the unclaimed mode holds a third copy of its fixed options, and the copy already differs from the other two. The same line also computes the re-plan before the guard that discards it. Citations at tip 0fa5fb29: clean_unclaimed.go:105; clean_unclaimed.go:113; clean_unclaimed.go:136.

Judged correct on this axis: the shared-apply seam, the side-by-side preflight pair, the test split boundary, the comment register, and the `cleanArguments` collapse.

## Spec at 0fa5fb29

Findings: 4. Worst issue: low.

- SPEC-4 (low, auto-fix): a preflight refusal marks every member not attempted. A member the plan retained or refused under CL10 loses its authority verdict. The apply loop already passes a non-removable row through untouched. Citations at tip 0fa5fb29: clean_set.go:283-289; clean_landed.go:336-341; clean_set_apply.go:22-30.
- SPEC-5 (ask-user): a replay after a spent apply renders a refusal with no recovery command. The selection no longer resolves, so the branch has no digest to re-plan. CL9 names a stale result, and this result is an unresolved selection. An agent reaches this case often after a partial apply. Citations at tip 0fa5fb29: clean_set.go:322-327.
- SPEC-6 (no-op, folds into ST11): the unclaimed re-plan reports a constant rather than the options the plan answered under. The rendered command is exact today, because the grammar forbids any other modifier on `--unclaimed`. Citation at tip 0fa5fb29: clean_unclaimed.go:101-108.
- SPEC-7 (no-op): the action vocabulary now has two homes. The value stays outside `Removes()`, which is what the row requires. No external surface enumerates action tokens. Placement is the Standards call ST10. Citations at tip 0fa5fb29: clean_set_apply.go:19; lifecyclepolicy.go:56-58.

## Coverage at 0fa5fb29

Findings: 4. Worst issue: medium.

- COV-3 (medium, auto-fix for the test, ask-user for the row): `TestCleanSetSpentPlan` survives three probed mutations. After the first apply the member checkouts are gone, so every refusal path exits 1 with no removed row. The test asserts a true observable that no mutation can turn red. Citations at tip 0fa5fb29: clean_set_apply_test.go:262; clean_set.go:281; clean_set.go:322; clean_set.go:330.
- COV-4 (ask-user): CL4 bites the entry fingerprint check, not the new preflight. Under the entry mutation the apply removes the clean member, because the drifted member re-plans as non-removable and the preflight skips it. The row's rationale credits the preflight for a refusal the entry check delivers. Citations at tip 0fa5fb29: clean_set_apply_test.go:106; clean_set.go:330; clean_set_apply.go:114.
- COV-5 (auto-fix): no fixture holds a retained member, so two mutations are silent across 984 tests. The retained pass-through and the preflight skip both have no test. Citations at tip 0fa5fb29: clean_set.go:287; clean_set_apply.go:114.
- COV-6 (low, auto-fix): the `not-attempted` wire token is asserted only through its own constant. A rename of the agent-facing token passes the gate, and the spec's field label breaks. Citation at tip 0fa5fb29: clean_set_outcomes_test.go:93.

Probe evidence returned with the axis. CL5 bit in both modes, CL7 bit, CL8 bit, and CL4 bit against the entry check. The axis did not probe CL6 and CL9, because each drives an injected fault or asserts an exact rendered command. CL18 and CL19 sit outside the delta as untouched regression anchors.

## Coordinator verification at 0fa5fb29

The coordinator ran the focused checks, the structure ratchet, and one independent probe before the commit.
`bench probe internal/worktree/clean_set_apply.go --omit` on the `notAttemptedPlan` action assignment returned `verdict=bit` and `restored=yes` against `TestCleanSetUnstartedOutcomes`.
`bench structure --growth HEAD` exits zero, and `.bench/structure-accept` holds no new grant.
The coordinator granted two fence expansions inside the approved behavior, and `bench learning` records each one.
The author reported `worktree.go` as unchanged. The coordinator read the diff and found five rewired call sites at an unchanged line count.

## CL-C2 repair at 52734949

The retained author repaired nine auto-fix targets. These are ST6, ST7, ST8, ST9, ST11 with SPEC-6, SPEC-4, COV-5, COV-6, and the test half of COV-3.

`staleRows` and `renderStaleSet` now own the stale refusal. `requalifyExplicitRow` gives explicit mode the owner landed mode already had.
`notAttemptedPlans` owns the unreached tail for both modes, and it passes a non-removable row through untouched.
`unclaimedOptions` is the one source of the unclaimed fixed options. The re-plan is computed inside the branch that uses it.

`retainedMemberFixture` holds one removable member and one retained member.
`TestCleanSetRetainedMember` covers the pass-through and the preflight skip.
`TestCleanSetSpentPlan` now uses that fixture. One member still resolves after the apply, so the refusal has a live source.

ST10 is reverted and stays open. The move of `ActionNotAttempted` into `lifecyclepolicy` reds the ratchet on two files.
`classifier.go` is 510 lines and `lifecyclepolicy.go` is 562 lines, against a 400 budget.
The only path through is a cap raise in `.bench/structure.budgets`, which is a reviewer-owned file.
The coordinator refused the raise. A cap loosens the ratchet for every later change to those two files.

The value keeps its declaration in `clean_set_apply.go`. Its doc comment now states why it sits outside the plan vocabulary.
The reviewer owns the disposition. A split of both files is the alternative to a cap.

The author reported five biting probes. The coordinator re-ran the COV-3 probe independently.
A `--swap` of the digest comparison to `if false` against `TestCleanSetSpentPlan` returned `verdict=bit` and `restored=yes`.
The focused checks pass across all seven packages. `bench structure --growth HEAD` exits zero with no new grant and no budget change.

Open for the reviewer: ST10, SPEC-5, and COV-4.

## Completion reconciliation at 835bf1c3

The chunk CL-C2 checkpoint passes. `bench test --package ./internal/worktree --run TestClean` and `--run TestLand` both report zero failures.
Two host skips remain. Each one reports that unix sockets are unavailable under `/tmp`, and neither touches an acceptance row.

The plan's named probe bit in both selection modes. An omission of the explicit preflight call failed `TestCleanSetPreflightAllRows/explicit`, and an omission of the landed preflight call failed the landed subtest. Both files restored.

The final verification passes. `bench test --package ./...` reports 99 packages green, and `bench test --check system` reports `internal/systemtest` green with zero skips.

The author walked all 19 acceptance rows. Eighteen rows resolve to a named test that exists and passes.
CL15 has no test by design, because the coverage map marks it review-owned for the ticket graph and the entry checks.
The author verified the CL15 claim directly. No cleanup source file names a numeric output budget, and ticket 2 depends only on ticket 1.

All 16 user stories are met by covered, passing rows. No material acceptance shortfall exists.

Three items stay open for the reviewer, and none of them is an unmet acceptance row. ST10 is a placement call with no behavior in it.
SPEC-5 asks whether CL9 extends to an unresolved selection. COV-4 asks whether CL4's stated reason should credit the entry check rather than the preflight.
