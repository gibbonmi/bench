# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/worktree-seam-reduction/spec.md",
  "plan_digest": "sha256:f16cc9cb766f3c75b004517d02033ea11e7831e516954910beca41fbcdbe298d",
  "implementation_session": "",
  "chunks": [
    {
      "id": "SR-C1",
      "base": "6ea6b7e6fe86e3fee0d0fc9a1ff01b6808486a69",
      "tip": "37e3f72f0fb6b0280b5106bbff0b2133f2ec9996",
      "plan_digest": "sha256:909699cd78b9c5c2fc464daae1bfef4d6fcb70aebf74afc38ed4d48385d8d976",
      "source_digest": "fefdc8ecd0a7d5876ca46611dcb71ca1538f7390",
      "acceptance_rows": [
        "WS1",
        "WS2",
        "WS3",
        "WS4",
        "WS5",
        "WS6",
        "WS7",
        "WS8"
      ],
      "verification": [
        {
          "id": "sr-c1-1-gate",
          "performer": "claude:bench-writer/sr-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e315aa279be4052e543a681cea96b92372bc44ef",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/sr-t1-author-20261002@6ea6b7e6fe86e3fee0d0fc9a1ff01b6808486a69",
            "digest": "sha256:abb0c3abebfe2a00e47a897b117b8226eeff62b78d6c476537509a6c63a05ed7",
            "excerpt": "tree[1]{target,head,dirty}:\n  sr-integration,0d7b40665f232bb348aa81d52e4c1360a9650c83,true\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,pass,20177\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "1-gate",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "sr-c1-1-gate-repair1",
          "performer": "claude:bench-writer/sr-t1-repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "fefdc8ecd0a7d5876ca46611dcb71ca1538f7390",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/sr-t1-repair1-20261002@d3637be279ad27405c6621bc3d22e356df777d1e",
            "digest": "sha256:e1dcc5f648c862587fa4d2c45b98dc3a5099263dfac865b0f6868d855e2ba976",
            "excerpt": "tree[1]{target,head,dirty}:\n  sr-integration,37e3f72f0fb6b0280b5106bbff0b2133f2ec9996,true\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,pass,19998\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "1-gate",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "sr-c1-r1-standards",
          "performer": "claude:bench-reviewer/sr-c1-r1-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "e315aa279be4052e543a681cea96b92372bc44ef",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/sr-c1-r1-standards-20261002@0d7b40665f232bb348aa81d52e4c1360a9650c83",
            "digest": "sha256:4f3c27a81b14e36247ffd42837603739870710ad1f6be3c60a1a7fc4b85e40a8",
            "excerpt": "## Standards\nS1. AGENTS.md code standard \"one source per fact\" (\"a fixture harness pasted N times\"). internal/gate/kit_value_test.go:101-103 (manifestLaneRoot) repeats internal/gate/lane_test.go:209-211 byte for byte; both build a root whose manifest declares a lane. lane_test.go lines 208-212 could call manifestLaneRoot(t). lane_test.go is not on ticket 1's Writes line, so the repair needs a plan commit. Disposition auto-fix, confidence 5.\ncount: 1\nworst: S1, a duplicated manifest-lane fixture literal.\nAdvice: internal/testreport/command.go:364-369 and internal/preprelease/preprelease.go:154 read BENCH_KIT outside the gate package; both predate this chunk. benchKitReaders repeats a non-test Go file scan that four other test files hold; no shared helper exists across packages.\n"
          },
          "axis": "Standards",
          "base": "6ea6b7e6fe86e3fee0d0fc9a1ff01b6808486a69",
          "tip": "0d7b40665f232bb348aa81d52e4c1360a9650c83",
          "finding_ids": [
            "S1"
          ],
          "supersedes": []
        },
        {
          "id": "sr-c1-r1-spec",
          "performer": "claude:bench-reviewer/sr-c1-r1-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "e315aa279be4052e543a681cea96b92372bc44ef",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/sr-c1-r1-spec-20261002@0d7b40665f232bb348aa81d52e4c1360a9650c83",
            "digest": "sha256:56c14989b127c4a2788d929b611e1d2aa30ef671a03e86a1610f4e00a12a5554",
            "excerpt": "## Spec\nNo findings survived refutation in SR-C1, at base 6ea6b7e6 and tip 0d7b4066.\nNeither AtKit form reaches KitValue (kit_source.go:78, kit_source.go:55). The empty-kit fallback goes through kitDirAt (kit_source.go:27) and does not call KitDir. The new forms are in kit_source.go. LaneForCommit, KitSourceCheckout, and KitDir keep their signatures; no change in internal/commit, internal/adopt, or cmd/bench. Commit 0d7b4066 touches only kit_source.go, kit_value_test.go, lane_select.go, and phases.go, all on the Writes line.\nRows WS1 to WS8: closed.\ncount: 0\nworst: none\nAdvice: the WS1 census matches only a string-literal argument; a read through a named constant or os.Environ would pass it. The row's wording accepts this.\n"
          },
          "axis": "Spec",
          "base": "6ea6b7e6fe86e3fee0d0fc9a1ff01b6808486a69",
          "tip": "0d7b40665f232bb348aa81d52e4c1360a9650c83",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "sr-c1-r1-coverage",
          "performer": "claude:bench-reviewer/sr-c1-r1-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "e315aa279be4052e543a681cea96b92372bc44ef",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/sr-c1-r1-coverage-20261002@0d7b40665f232bb348aa81d52e4c1360a9650c83",
            "digest": "sha256:d015d845e9acb964606c0a9ea0df15a7539a8e73e4fd5dba45cb94b9d63c8527",
            "excerpt": "## Coverage\nC1. State: KitSourceCheckoutAtKit(root, \"\") when the kitDirAt fallback answers the working directory in place of the executable's grandparent (internal/gate/kit_source.go:31-34). Expected break: with BENCH_KIT unset, a consumer repository run from its own root reads as the kit checkout. WS6 uses a temporary root apart from the working directory, so it passes. The test that should exist: WS6 also calls t.Chdir(root). Disposition auto-fix, confidence 6.\nProbes on internal/gate/kit_source.go, each restored=yes: (1) swap `err == nil` to `err != nil` in kitDirAt, ./internal/gate, silent; (2) the same swap, ./internal/adopt, silent; (3) swap `if root == \"\"` to `if false`, silent; (4) swap `err != nil || checks == nil` to `err != nil`, silent.\ncount: 1\nworst: C1.\nAdvice: the empty-root guard has no row; the nil-lane return has no gate-package row; a relative or missing kit path has no row and predates this delta.\nFinal git status: clean.\n"
          },
          "axis": "Coverage",
          "base": "6ea6b7e6fe86e3fee0d0fc9a1ff01b6808486a69",
          "tip": "0d7b40665f232bb348aa81d52e4c1360a9650c83",
          "finding_ids": [
            "C1"
          ],
          "supersedes": []
        },
        {
          "id": "sr-c1-r2-standards",
          "performer": "claude:bench-reviewer/sr-c1-r2-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "fefdc8ecd0a7d5876ca46611dcb71ca1538f7390",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/sr-c1-r2-standards-20261002@37e3f72f0fb6b0280b5106bbff0b2133f2ec9996",
            "digest": "sha256:1d4e9b5946b0ee3a777d4742454dc928e0bcdf21e00fc4a2a181ba91808d1057",
            "excerpt": "Standards axis, SR-C1 confirming round, repair delta 0d7b4066..37e3f72f: pass, 0 findings.\nS1: confirmed. internal/gate/lane_test.go:208 calls LaneForCommit(manifestLaneRoot(t)). The manifest literal exists in one place, manifestLaneRoot at internal/gate/kit_value_test.go:99-105.\nThe changed comment at internal/gate/kit_value_test.go:86-87 meets bench-craft-comments.\nblocking findings: none.\n"
          },
          "axis": "Standards",
          "base": "6ea6b7e6fe86e3fee0d0fc9a1ff01b6808486a69",
          "tip": "37e3f72f0fb6b0280b5106bbff0b2133f2ec9996",
          "finding_ids": [],
          "supersedes": [
            "sr-c1-r1-standards"
          ]
        },
        {
          "id": "sr-c1-r2-spec",
          "performer": "claude:bench-reviewer/sr-c1-r2-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "fefdc8ecd0a7d5876ca46611dcb71ca1538f7390",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/sr-c1-r2-spec-20261002@37e3f72f0fb6b0280b5106bbff0b2133f2ec9996",
            "digest": "sha256:e086005dd4dc1d8fc05e085f82bf5b6a3c9d3c4070efb584fd3cd7267a4cc78c",
            "excerpt": "Spec axis, SR-C1 confirming round, repair delta 0d7b4066..37e3f72f: pass, 0 findings.\nC1 confirmed: the WS6 test keeps its name, adds t.Chdir(root), and still fails if KitSourceCheckoutAtKit(root, \"\") is true (spec.md:364).\nS1 confirmed: internal/gate/lane_test.go:208 replaces the inline manifest with LaneForCommit(manifestLaneRoot(t)); the assertions are unchanged. The axis did not read the manifest body of manifestLaneRoot within its budget.\nScope confirmed: no production code changed; ticket 1's Writes line holds internal/gate/lane_test.go.\nblocking findings: none.\n"
          },
          "axis": "Spec",
          "base": "6ea6b7e6fe86e3fee0d0fc9a1ff01b6808486a69",
          "tip": "37e3f72f0fb6b0280b5106bbff0b2133f2ec9996",
          "finding_ids": [],
          "supersedes": [
            "sr-c1-r1-spec"
          ]
        },
        {
          "id": "sr-c1-r2-coverage",
          "performer": "claude:bench-reviewer/sr-c1-r2-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "fefdc8ecd0a7d5876ca46611dcb71ca1538f7390",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/sr-c1-r2-coverage-20261002@37e3f72f0fb6b0280b5106bbff0b2133f2ec9996",
            "digest": "sha256:8396dfa309b40ebe7e795aabf6e4097bc4fdcb6b5431e008eaf396452d5b0ef3",
            "excerpt": "Coverage axis, SR-C1 confirming round, repair delta 0d7b4066..37e3f72f: pass, 0 findings.\nC1: confirmed. internal/gate/kit_value_test.go adds t.Chdir(root) to the WS6 test.\nProbe: bench probe internal/gate/kit_source.go --swap 'filepath.Join(filepath.Dir(exe), \"..\")' --with 'filepath.Join(\".\", exe[:0])' --package ./internal/gate --run TestKitSourceCheckoutAtKitFallsBackToTheExecutableParent. Verdict bit, 1 failed test (kit_value_test.go:95), restored yes. The first swap in the charge did not compile and returned invalid, restored yes.\nblocking findings: none.\nFinal git status: empty.\nAn earlier session for this axis stalled with no output and the coordinator stopped it; that transport is incomplete.\n"
          },
          "axis": "Coverage",
          "base": "6ea6b7e6fe86e3fee0d0fc9a1ff01b6808486a69",
          "tip": "37e3f72f0fb6b0280b5106bbff0b2133f2ec9996",
          "finding_ids": [],
          "supersedes": [
            "sr-c1-r1-coverage"
          ]
        }
      ]
    },
    {
      "id": "SR-C2",
      "base": "37e3f72f0fb6b0280b5106bbff0b2133f2ec9996",
      "tip": "f620ccd6e416501f2d4cf504d5090f3cb5104c80",
      "plan_digest": "sha256:f16cc9cb766f3c75b004517d02033ea11e7831e516954910beca41fbcdbe298d",
      "source_digest": "b46d93c63bb58354d46a8869edec68fe5fb1e180",
      "acceptance_rows": [
        "WS9",
        "WS10",
        "WS11",
        "WS12",
        "WS13",
        "WS14",
        "WS15",
        "WS16"
      ],
      "verification": [
        {
          "id": "sr-c2-2-worktree",
          "performer": "claude:bench-writer/sr-t2-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "b46d93c63bb58354d46a8869edec68fe5fb1e180",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/sr-t2-author-20261002@3ecdaffff22967433363e9495c8807f14e2f85ea",
            "digest": "sha256:b3be7a9e9dbb0eb6edc3bf2baa842efe38557f107ba93a8cabed79d58d2cdedb",
            "excerpt": "tree[1]{target,head,dirty}:\n  sr-integration,f620ccd6e416501f2d4cf504d5090f3cb5104c80,true\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,pass,56225\nfailures[0]{package,test,line}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanLandedSpecialPathsRetainedWithoutOpening/socket,\"clean_landed_hostile_test.go:98: unix sockets unavailable: listen unix ... (280 bytes)\"\n  github.com/gibbonmi/bench/internal/worktree,TestLandedConsumersRejectSpecialGitMetadataBeforePlanning/socket,\"capability: fifo: unix sockets unavailable: listen unix ... (279 bytes)\"\n"
          },
          "requirement": "2-worktree",
          "command": "bench test --package ./internal/worktree",
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
      "from": "sha256:c9feebcde5a12aebbaeb53a8476ca2a1b81cd3483ee036c1544f4df7b7b01e8f",
      "to": "sha256:909699cd78b9c5c2fc464daae1bfef4d6fcb70aebf74afc38ed4d48385d8d976",
      "chunk_ids": {
        "SR-C1": [
          "SR-C1"
        ]
      }
    },
    {
      "from": "sha256:909699cd78b9c5c2fc464daae1bfef4d6fcb70aebf74afc38ed4d48385d8d976",
      "to": "sha256:f16cc9cb766f3c75b004517d02033ea11e7831e516954910beca41fbcdbe298d",
      "chunk_ids": {
        "SR-C1": [
          "SR-C1"
        ]
      }
    }
  ]
}
```

## SR-C1 chunk review, round 1

Three fresh opus / high sessions reviewed the frozen pair `6ea6b7e6..0d7b4066`. Each axis bound the review evidence `sha256:1f8096c6` with `--check-current`. The raw finding count is 2, and the repair-target count is 2. The repair allowance of SR-C1 is 2 cycles, and 0 cycles are used.

### Standards

Count: 1. Worst: S1.

- S1 (`auto-fix`, confidence 5): `manifestLaneRoot` in `internal/gate/kit_value_test.go` lines 101 to 103 repeats the manifest literal in `internal/gate/lane_test.go` lines 209 to 211. The rule is the `AGENTS.md` code standard, one source per fact. The repair makes `lane_test.go` call the helper, so ticket 1's `Writes:` line takes `internal/gate/lane_test.go`.

### Spec

Count: 0. The axis closed the rows WS1 to WS8.

### Coverage

Count: 1. Worst: C1.

- C1 (`auto-fix`, confidence 6): no test pins the empty-kit fallback of `kitDirAt` to the executable's parent. A probe that swaps `err == nil` for `err != nil` at `internal/gate/kit_source.go` line 31 stayed silent in `./internal/gate` and in `./internal/adopt`. The spec binds that fallback in "The gate kit reader", and the edge inventory says that WS6 covers it. The repair makes the WS6 test red for a fallback to the working directory.

### Advice

- The WS1 census matches only a string-literal argument.
- `internal/testreport/command.go` and `internal/preprelease/preprelease.go` read `BENCH_KIT` outside the gate package. Both reads predate this chunk.
- The empty-root guard and the nil-lane return have no gate-package row.

## SR-C1 repair cycle 1

One fresh opus / high repair session, `claude:bench-writer/sr-t1-repair1`, repaired S1 and C1 in commit `37e3f72f`. The session used 1 of 2 attempts and 1 lane pass. The repair allowance of SR-C1 is 2 cycles, and 1 cycle is used.

- S1: `TestLaneForCommitMarksOnlyTheKitLaneSelective` now calls `manifestLaneRoot`, and its assertions did not change.
- C1: the WS6 test makes the root the working directory. The named swap probe at `kit_source.go` line 31 was `silent` before the repair and `bit` after it.

The coordinator probe omitted `Selective: true` in `LaneForCommitAtKit`. It returned `bit` with 1 failed test and `restored=yes`.

## SR-C1 confirming round

Three fresh opus / high sessions read the repair delta `0d7b4066..37e3f72f`. Each axis confirmed its folds and reported no blocking finding. The Coverage axis ran one swap probe on the `kitDirAt` fallback, and it returned `bit` with `restored=yes`. An earlier Coverage session stalled with no output, and the coordinator stopped it. SR-C1 has no open finding, and 1 of 2 repair cycles is used.
