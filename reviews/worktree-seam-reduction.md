# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/worktree-seam-reduction/spec.md",
  "plan_digest": "sha256:04d1f52fa7224e0708fedd02e7c1118f8404219f710c1aea080eb347e589fc9f",
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
      "tip": "1fde1ac102e6b169c46bc40ae03510b29469d1fd",
      "plan_digest": "sha256:04d1f52fa7224e0708fedd02e7c1118f8404219f710c1aea080eb347e589fc9f",
      "source_digest": "995f6f325872860920aebf1e60aaee419f87a7cb",
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
        },
        {
          "id": "sr-c2-2-worktree-repair1",
          "performer": "claude:bench-writer/sr-t2-repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "995f6f325872860920aebf1e60aaee419f87a7cb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/sr-t2-repair1-20261002@d68dd224133551cdf60152916b9f86842c5ab911",
            "digest": "sha256:b54fa6ab9c0451340a3e0d484d417c88dab23efbd67ab4b668ce382b0bbf0df2",
            "excerpt": "tree[1]{target,head,dirty}:\n  sr-integration,1fde1ac102e6b169c46bc40ae03510b29469d1fd,true\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,pass,57473\nfailures[0]{package,test,line}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/worktree,TestCleanLandedSpecialPathsRetainedWithoutOpening/socket,\"capability: fifo: unix sockets unavailable\"\n  github.com/gibbonmi/bench/internal/worktree,TestLandedConsumersRejectSpecialGitMetadataBeforePlanning/socket,\"capability: fifo: unix sockets unavailable\"\n"
          },
          "requirement": "2-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "sr-c2-r1-standards",
          "performer": "claude:bench-reviewer/sr-c2-r1-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "b46d93c63bb58354d46a8869edec68fe5fb1e180",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/sr-c2-r1-standards-20261002@f620ccd6e416501f2d4cf504d5090f3cb5104c80",
            "digest": "sha256:ff08fa9bc125b1568d49305197d8ff0e8cb0fa3a1c70e21e2c032a864c2f9bf0",
            "excerpt": "## Standards\nS1. craft-comments \"describe the current state of the code\". internal/worktree/verb_runner_test.go:163, the checkVerbCall doc comment, says the refusal holds because a test must not trust a value the verb ignores. The check refuses only a verb with no joins form. No joins form reads a.kit at this tip, and merge, reset, reauthorize, and build never read a.now, so a kit or clock value for a joins-form verb passes and is ignored. Fix: state what the check enforces and drop the guarantee. Disposition auto-fix, confidence 7.\nS2. The same rule. internal/worktree/effects.go:32-33, the newAmbient comment, says the home and the stderr writer come from the verb entry's own parameters. resume.go:76, :79, and :368 pass newAmbient(Home(), os.Stderr). Fix: describe what the constructor takes from its caller. Disposition auto-fix, confidence 5.\ncount: 2\nworst: S1.\nAdvice: the ambient type comment at effects.go:22-24 ends \"so nothing below the entry reads the process\", which tickets 3, 7, and 11 make true; examine it again at SR-C6. The verbCall doc comment repeats the dispatch rule of runVerb and callAmbient.\nAuthor claims verified: newAmbient is the one constructor; discardCall replaces discardJoins; requalifyUnrecordedRow has no joins parameter; j.now and j.home are gone.\n"
          },
          "axis": "Standards",
          "base": "37e3f72f0fb6b0280b5106bbff0b2133f2ec9996",
          "tip": "f620ccd6e416501f2d4cf504d5090f3cb5104c80",
          "finding_ids": [
            "S1",
            "S2"
          ],
          "supersedes": []
        },
        {
          "id": "sr-c2-r1-spec",
          "performer": "claude:bench-reviewer/sr-c2-r1-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "b46d93c63bb58354d46a8869edec68fe5fb1e180",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/sr-c2-r1-spec-20261002@f620ccd6e416501f2d4cf504d5090f3cb5104c80",
            "digest": "sha256:733b4f499ea002c11d283fde77aa314851a8b96e279107b7bc24e943bcccab87",
            "excerpt": "## Spec\nBind: the --check-current call ran from the primary checkout and exited 1 (base not an ancestor). The axis read each file at f620ccd6 through git show.\nP1. Spec line 455: \"A path with a space: the resume-clean fixtures use `auto clean` and `auto dirty` paths. WS12 reuses that fixture shape.\" internal/worktree/resume_clean_ambient_test.go:18 builds newPendingAssignment, whose pool path holds no space, so the stated edge is not exercised. Fix: reuse the space-path fixture, or amend the edge line. Disposition auto-fix, confidence 5.\nRows WS9 to WS16: closed.\nThe other predicates hold: 9 joins-form entries call newAmbient once; now and home left the joins value; executeCleanup drops records under a.home (lifecycle.go:447); ListCommand takes no joins home; the runner and checkVerbCall match spec lines 166 to 172; worktreeTestCount went from 699 to 701; the commit touches only internal/worktree; no over-budget file grew.\ncount: 1\nworst: P1.\nAdvice: ticket 11's ApplyAutomatic item is already met; two resume_test.go helpers still build newAmbient(Home(), os.Stderr); checkVerbCall with an unknown key and a kit value reports the kit refusal first.\n"
          },
          "axis": "Spec",
          "base": "37e3f72f0fb6b0280b5106bbff0b2133f2ec9996",
          "tip": "f620ccd6e416501f2d4cf504d5090f3cb5104c80",
          "finding_ids": [
            "P1"
          ],
          "supersedes": []
        },
        {
          "id": "sr-c2-r1-coverage",
          "performer": "claude:bench-reviewer/sr-c2-r1-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "b46d93c63bb58354d46a8869edec68fe5fb1e180",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/sr-c2-r1-coverage-20261002@f620ccd6e416501f2d4cf504d5090f3cb5104c80",
            "digest": "sha256:d90a24cdec111c4e4efdc07e8a7dd6811229d1a59ba5a5327a0772bee8d98c98",
            "excerpt": "## Coverage\nC1. State: resume-clean over a registered, active, unlanded checkout whose age is 8 days by the clock value. The per-checkout plan at internal/worktree/resume.go:414 can read currentTime() and WS13 stays green; the checkout then counts as retained active=1, not stale-active=1 (worktree.go:429). TestResumeCleanJudgesTheClockValue asserts only the clean line of the orphan sweep. The test should also assert the retained stale-active=1 cell (story 11). Disposition auto-fix, confidence 7.\nC2 (classified as optional advice by the coordinator, see the pickup): a landing resumed under an explicit home that reaches its release step. No test drops the census record under the explicit home through land --resume. The code is correct today. The axis gave ask-user, confidence 5.\nProbes, each --package ./internal/worktree, each restored yes: (1) clean_discard.go plan instant to currentTime(), bit; (2) ownership.go release-resume replan instant to currentTime(), silent; (3) land.go release ambient to newAmbient(Home(), stderr), bit; (4) land_resume.go the same swap, silent; (5) resume.go planAutomaticAt instant to currentTime(), silent.\ncount: 2\nworst: C1.\nAdvice: the release-resume replan instant has no behavioral test; no production code reads a.kit or a.warnings at this tip; checkVerbCall is tested only for the path key.\nFinal git status: clean.\n"
          },
          "axis": "Coverage",
          "base": "37e3f72f0fb6b0280b5106bbff0b2133f2ec9996",
          "tip": "f620ccd6e416501f2d4cf504d5090f3cb5104c80",
          "finding_ids": [
            "C1"
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
    },
    {
      "from": "sha256:f16cc9cb766f3c75b004517d02033ea11e7831e516954910beca41fbcdbe298d",
      "to": "sha256:04d1f52fa7224e0708fedd02e7c1118f8404219f710c1aea080eb347e589fc9f",
      "chunk_ids": {
        "SR-C1": [
          "SR-C1"
        ],
        "SR-C2": [
          "SR-C2"
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

## SR-C2 chunk review, round 1

Three fresh opus / high sessions reviewed the frozen pair `37e3f72f..f620ccd6`. The coordinator probe omitted `home: home` in `newAmbient`, and it returned `bit` with 3 failed tests and `restored=yes`. The raw finding count is 4, and the repair-target count is 4. The repair allowance of SR-C2 is 2 cycles, and 0 cycles are used.

### Standards

Count: 2. Worst: S1.

- S1 (`auto-fix`, confidence 7): the `checkVerbCall` doc comment at `internal/worktree/verb_runner_test.go` line 163 states a guarantee that the check does not give. The check refuses only a verb with no joins form. The rule is the current-state rule of `craft-comments`.
- S2 (`auto-fix`, confidence 5): the `newAmbient` comment at `internal/worktree/effects.go` lines 32 to 33 says that the home comes from the verb entry's parameters. Three callers in `resume.go` pass `Home()`.

### Spec

Count: 1. Worst: P1. The axis closed the rows WS9 to WS16.

- P1 (`auto-fix`, confidence 5): the edge inventory says that WS12 reuses the space-path fixture shape. The test in `internal/worktree/resume_clean_ambient_test.go` builds a path with no space.

### Coverage

Count: 1. Worst: C1.

- C1 (`auto-fix`, confidence 7): `TestResumeCleanJudgesTheClockValue` asserts only the clean line of the orphan sweep. A probe that swaps the instant of `planAutomaticAt` in `resume.go` for `currentTime()` stayed silent. Story 11 requires one instant for each plan. The repair makes the test assert the retained `stale-active=1` cell.

### Advice

- The Coverage axis named a second item with `ask-user`: no test drops the census record under an explicit home through `land --resume`. The code is correct at this tip, and no approved row names that path. The coordinator holds it as optional advice under the bounded repair policy, for reviewer veto. The live census of SR-C6 refuses a `Home()` read below an entry.
- The release-resume replan instant in `ownership.go` has no behavioral test. The live census of SR-C6 refuses that read.
- The ambient type comment in `effects.go` says that nothing below the entry reads the process. The tickets 3, 7, and 11 make that true, so SR-C6 examines it again.
- The item `ApplyAutomatic` of ticket 11 is already met at this tip.

## SR-C2 repair cycle 1

One fresh opus / high repair session, `claude:bench-writer/sr-t2-repair1`, repaired S1, S2, P1, and C1 in commit `1fde1ac1`. The session used 1 of 2 attempts and 1 lane pass. The repair allowance of SR-C2 is 2 cycles, and 1 cycle is used.

- S1 and S2: the two comments now state what the code does. No code changed.
- P1: the WS12 test puts the home at a sibling path with a space. A guard fails the test when the checkout path holds no space.
- C1: the WS13 test also requires the `retained stale-active=1` cell. The named swap probe in `resume.go` was `silent` before the repair and `bit` after it.

The coordinator probe swapped the home argument of `census.Drop` in `lifecycle.go`. It returned `bit` with 1 failed test and `restored=yes`. An omission of that line did not compile, so that probe was `invalid`.
