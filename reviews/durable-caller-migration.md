# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/durable-caller-migration/spec.md",
  "plan_digest": "sha256:5cb50643a5295119fe76284a8cc58c0cc7bfd797ae01bfecc20e1b30a5df98af",
  "implementation_session": "",
  "chunks": [
    {
      "id": "D-B1",
      "base": "69cffaf6a7416abab014e5e1b352a683f3cbd566",
      "tip": "e886531e4923b056e469e1fd9b06daa698b39b8e",
      "plan_digest": "sha256:4b8198b6b5417276e4d9ccfeb3ddbd2f4b4439b06c5f3a27abc84d74e6645358",
      "source_digest": "2b201e4e9ad334e3d653d68e027c3c415b1aecf6",
      "acceptance_rows": [
        "D16",
        "D17",
        "D37",
        "D38",
        "D39",
        "D40"
      ],
      "verification": [
        {
          "id": "dcm-t02-gate-tests",
          "performer": "claude:dcm_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2c57c2d6e8525a84409078f842647f628de3faa9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "c87faf66abc2a8eb228b7179aeb57d0a075f9dd6",
            "digest": "sha256:04964e228789f40d10c8ee97aef690d9fd11e3c787a962b9ad34608eb3822de2",
            "excerpt": "bench test --package ./internal/gate at c87faf66abc2a8eb228b7179aeb57d0a075f9dd6\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,24691,447\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "gate-tests",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "dcm-t02-root-conformance",
          "performer": "claude:dcm_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2c57c2d6e8525a84409078f842647f628de3faa9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "c87faf66abc2a8eb228b7179aeb57d0a075f9dd6",
            "digest": "sha256:0c921c5cc503c80b7ea0ab5060ef4fab76b90ac2d9f11c48cd79fb57af1c5335",
            "excerpt": "bench test --check root-conformance (ticket tree, content identical to c87faf66abc2a8eb228b7179aeb57d0a075f9dd6)\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17629,1\nfailures[0]{package,test,line,lines}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "dcm-t02-omission-proof",
          "performer": "claude:dcm_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2c57c2d6e8525a84409078f842647f628de3faa9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "c87faf66abc2a8eb228b7179aeb57d0a075f9dd6",
            "digest": "sha256:791e565bdbd2ab1f088af63800e1378d5008ace87a2719d1af9a416f17d671b1",
            "excerpt": "# omission-proof: four mutations through bench probe, each restored\n## 1791820490243913036-fd2513653c6ca820\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,69cffaf6a7416abab014e5e1b352a683f3cbd566,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_durability.go,swap,failed,5,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestGatePendingDurabilityFailure|TestGateTerminalDurabilityFailure|TestGateEvidenceDurabilityFailure|TestGatePendingRecoveryFailure,passed,6\nspilled{lines=15,bytes=2331,omitted_lines=6,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791820490242029864-40eb5c973985f9b2.out}\n  github.com/gibbonmi/bench/internal/gate,TestGatePendingDurabilityFailure,\"durable_replacement_test.go:83: stderr = \\\"gate pending persistence failed\\\\\\\\n\\\", want \\\"gate pending persistence failed: replace /dev/shm/GL2WST/t/TestGatePendingDurabilityFailure3638066449/001/.git/bench-last-gate: directory-sync: injected pen… (327 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestGatePendingRecoveryFailure,\"durable_replacement_test.go:160: stderr = \\\"gate final persistence failed\\\\\\\\n\\\", want \\\"gate final persistence failed: replace /dev/shm/GL2WST/t/TestGatePendingRecoveryFailure513445909/001/.git/bench-last-gate: directory-sync: injected final dir… (478 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestGateTerminalDurabilityFailure/green,\"durable_replacement_test.go:111: stderr = \\\"gate final persistence failed\\\\\\\\n\\\", want \\\"gate final persistence failed: replace /dev/shm/GL2WST/t/TestGateTerminalDurabilityFailuregreen552843700/001/.git/bench-last-gate: directory-sync: injected f… (351 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestGateTerminalDurabilityFailure/red,\"durable_replacement_test.go:111: stderr = \\\"gate final persistence failed\\\\\\\\n\\\", want \\\"gate final persistence failed: replace /dev/shm/GL2WST/t/TestGateTerminalDurabilityFailurered1363136781/001/.git/bench-last-gate: directory-sync: injected fi… (350 bytes)\",1\nskips[0]{package,test,reason}:\n## 1791820515396595618-f1ac02f09ffdaecf\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,69cffaf6a7416abab014e5e1b352a683f3cbd566,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_transaction.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestGatePendingDurabilityFailure,passed,1\nspilled{lines=11,bytes=699,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791820515394703718-955f5d5234f32396.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,34,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGatePendingDurabilityFailure,\"durable_replacement_test.go:81: run counter stat = <nil>, want no oracle start after the pending failure\",1\nskips[0]{package,test,reason}:\n## 1791820525390246433-91b93c48a7237799\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,69cffaf6a7416abab014e5e1b352a683f3cbd566,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_transaction.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestGateEvidenceDurabilityFailure,passed,1\nspilled{lines=11,bytes=861,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791820525388434471-6639e4deaced3b4f.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,34,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateEvidenceDurabilityFailure,\"durable_replacement_test.go:132: result = gate.Result{GateExit:0, ActionExit:0, Inspection:gate.Inspection{State:\\\"ready\\\", Status:\\\"green\\\", PendingStatus:\\\"\\\", CachedTree:\\\"8dcab660b579c8a554fcbd2b5838e00021651b7c\\\", CurrentTree:\\\"8dcab660b579c8a5… (474 bytes)\",1\nskips[0]{package,test,reason}:\n## 1791820534108047915-d81542c36f359997\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,69cffaf6a7416abab014e5e1b352a683f3cbd566,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_durability.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestGatePendingRecoveryFailure|TestGateDemotionDurabilityFailure,passed,2\nspilled{lines=12,bytes=1221,omitted_lines=3,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791820534106206975-2032b6fe8e497192.out}\n  github.com/gibbonmi/bench/internal/gate,fail,91,2\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateDemotionDurabilityFailure,\"durable_replacement_test.go:194: stderr = \\\"gate lock unavailable; interrupted-green demotion failed: replace /dev/shm/6KOGNW/t/TestGateDemotionDurabilityFailure4157545712/001/.git/bench-last-gate: rename: injected demotion rename failure\\\\\\\\n\\\"… (615 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestGatePendingRecoveryFailure,\"durable_replacement_test.go:160: stderr = \\\"gate final persistence failed: replace /dev/shm/6KOGNW/t/TestGatePendingRecoveryFailure840713086/001/.git/bench-last-gate: directory-sync: injected final directory-sync failure; the rename complete… (684 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "omission-proof",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0,
          "probe": {
            "mutation": "Discard the leaf cause, allow oracle start after pending failure, hide evidence failure, or discard the recovery error. Each mutation must fail its real execution witness.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "c87faf66abc2a8eb228b7179aeb57d0a075f9dd6",
              "digest": "sha256:791e565bdbd2ab1f088af63800e1378d5008ace87a2719d1af9a416f17d671b1",
              "excerpt": "# omission-proof: four mutations through bench probe, each restored\n## 1791820490243913036-fd2513653c6ca820\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,69cffaf6a7416abab014e5e1b352a683f3cbd566,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_durability.go,swap,failed,5,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestGatePendingDurabilityFailure|TestGateTerminalDurabilityFailure|TestGateEvidenceDurabilityFailure|TestGatePendingRecoveryFailure,passed,6\nspilled{lines=15,bytes=2331,omitted_lines=6,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791820490242029864-40eb5c973985f9b2.out}\n  github.com/gibbonmi/bench/internal/gate,TestGatePendingDurabilityFailure,\"durable_replacement_test.go:83: stderr = \\\"gate pending persistence failed\\\\\\\\n\\\", want \\\"gate pending persistence failed: replace /dev/shm/GL2WST/t/TestGatePendingDurabilityFailure3638066449/001/.git/bench-last-gate: directory-sync: injected pen… (327 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestGatePendingRecoveryFailure,\"durable_replacement_test.go:160: stderr = \\\"gate final persistence failed\\\\\\\\n\\\", want \\\"gate final persistence failed: replace /dev/shm/GL2WST/t/TestGatePendingRecoveryFailure513445909/001/.git/bench-last-gate: directory-sync: injected final dir… (478 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestGateTerminalDurabilityFailure/green,\"durable_replacement_test.go:111: stderr = \\\"gate final persistence failed\\\\\\\\n\\\", want \\\"gate final persistence failed: replace /dev/shm/GL2WST/t/TestGateTerminalDurabilityFailuregreen552843700/001/.git/bench-last-gate: directory-sync: injected f… (351 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestGateTerminalDurabilityFailure/red,\"durable_replacement_test.go:111: stderr = \\\"gate final persistence failed\\\\\\\\n\\\", want \\\"gate final persistence failed: replace /dev/shm/GL2WST/t/TestGateTerminalDurabilityFailurered1363136781/001/.git/bench-last-gate: directory-sync: injected fi… (350 bytes)\",1\nskips[0]{package,test,reason}:\n## 1791820515396595618-f1ac02f09ffdaecf\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,69cffaf6a7416abab014e5e1b352a683f3cbd566,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_transaction.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestGatePendingDurabilityFailure,passed,1\nspilled{lines=11,bytes=699,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791820515394703718-955f5d5234f32396.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,34,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGatePendingDurabilityFailure,\"durable_replacement_test.go:81: run counter stat = <nil>, want no oracle start after the pending failure\",1\nskips[0]{package,test,reason}:\n## 1791820525390246433-91b93c48a7237799\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,69cffaf6a7416abab014e5e1b352a683f3cbd566,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_transaction.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestGateEvidenceDurabilityFailure,passed,1\nspilled{lines=11,bytes=861,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791820525388434471-6639e4deaced3b4f.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,34,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateEvidenceDurabilityFailure,\"durable_replacement_test.go:132: result = gate.Result{GateExit:0, ActionExit:0, Inspection:gate.Inspection{State:\\\"ready\\\", Status:\\\"green\\\", PendingStatus:\\\"\\\", CachedTree:\\\"8dcab660b579c8a554fcbd2b5838e00021651b7c\\\", CurrentTree:\\\"8dcab660b579c8a5… (474 bytes)\",1\nskips[0]{package,test,reason}:\n## 1791820534108047915-d81542c36f359997\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,69cffaf6a7416abab014e5e1b352a683f3cbd566,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_durability.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestGatePendingRecoveryFailure|TestGateDemotionDurabilityFailure,passed,2\nspilled{lines=12,bytes=1221,omitted_lines=3,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791820534106206975-2032b6fe8e497192.out}\n  github.com/gibbonmi/bench/internal/gate,fail,91,2\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateDemotionDurabilityFailure,\"durable_replacement_test.go:194: stderr = \\\"gate lock unavailable; interrupted-green demotion failed: replace /dev/shm/6KOGNW/t/TestGateDemotionDurabilityFailure4157545712/001/.git/bench-last-gate: rename: injected demotion rename failure\\\\\\\\n\\\"… (615 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestGatePendingRecoveryFailure,\"durable_replacement_test.go:160: stderr = \\\"gate final persistence failed: replace /dev/shm/6KOGNW/t/TestGatePendingRecoveryFailure840713086/001/.git/bench-last-gate: directory-sync: injected final directory-sync failure; the rename complete… (684 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dcm-t02-preservation",
          "performer": "claude:dcm_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2c57c2d6e8525a84409078f842647f628de3faa9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "c87faf66abc2a8eb228b7179aeb57d0a075f9dd6",
            "digest": "sha256:a9bbe79d939f18df8f8d50516928a3268493a826d5217462ea90e594e6850093",
            "excerpt": "# preservation: drop the final newline from the verdict and store writers; independent expected bytes turn red; restored\n## 1791820544916413834-78a0c22da8180ec0\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,69cffaf6a7416abab014e5e1b352a683f3cbd566,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/verdict.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestVerdictReplacement,passed,1\nspilled{lines=11,bytes=951,omitted_lines=2,cut_lines=1,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791820544914541511-07ce19bc4ad140f2.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,4,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestVerdictReplacement,\"durable_replacement_test.go:30: record = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":1,\\\\\\\\\\\"state\\\\\\\\\\\":\\\\\\\\\\\"ready\\\\\\\\\\\",\\\\\\\\\\\"status\\\\\\\\\\\":\\\\\\\\\\\"green\\\\\\\\\\\",\\\\\\\\\\\"tree\\\\\\\\\\\":\\\\\\\\\\\"abc123\\\\\\\\\\\",\\\\\\\\\\\"oracle\\\\\\\\\\\":\\\\\\\\\\\"feed\\\\\\\\\\\",\\\\\\\\\\\"recorded_at\\\\\\\\\\\":\\\\\\\\\\\"2026-10-12T00:00:00Z\\\\\\\\\\\"}\\\", want \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":1,\\\\\\\\\\\"state\\\\\\\\\\\":\\\\\\\\\\\"ready\\\\\\\\\\\",\\\\\\\\\\\"s\nskips[0]{package,test,reason}:\n## 1791820618400582478-15c6df5b9c04a125\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,69cffaf6a7416abab014e5e1b352a683f3cbd566,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/engine.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestEvidenceReplacement,passed,3\nspilled{lines=12,bytes=1419,omitted_lines=3,cut_lines=2,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791820618398735692-39ebbefeaa7f6c09.out}\n  github.com/gibbonmi/bench/internal/gate,fail,4,3\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestEvidenceReplacement/evidence,\"durable_replacement_test.go:67: record = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":1,\\\\\\\\\\\"state\\\\\\\\\\\":\\\\\\\\\\\"ready\\\\\\\\\\\",\\\\\\\\\\\"status\\\\\\\\\\\":\\\\\\\\\\\"green\\\\\\\\\\\",\\\\\\\\\\\"tree\\\\\\\\\\\":\\\\\\\\\\\"abc123\\\\\\\\\\\",\\\\\\\\\\\"oracle\\\\\\\\\\\":\\\\\\\\\\\"feed\\\\\\\\\\\",\\\\\\\\\\\"recorded_at\\\\\\\\\\\":\\\\\\\\\\\"2026-10-12T00:00:00Z\\\\\\\\\\\"}\\\", want \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":1,\\\\\\\\\\\"state\\\\\\\\\\\":\\\\\\\\\\\"ready\\\\\\\\\n  github.com/gibbonmi/bench/internal/gate,TestEvidenceReplacement/lane,\"durable_replacement_test.go:67: record = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":1,\\\\\\\\\\\"tree\\\\\\\\\\\":\\\\\\\\\\\"abc123\\\\\\\\\\\",\\\\\\\\\\\"recorded_at\\\\\\\\\\\":\\\\\\\\\\\"2026-10-12T00:00:00Z\\\\\\\\\\\",\\\\\\\\\\\"lane\\\\\\\\\\\":\\\\\\\\\\\"fast\\\\\\\\\\\",\\\\\\\\\\\"outcome\\\\\\\\\\\":\\\\\\\\\\\"pass\\\\\\\\\\\",\\\\\\\\\\\"run_binary\\\\\\\\\\\":\\\\\\\\\\\"beef\\\\\\\\\\\"}\\\", want \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":1,\\\\\\\\\\\"tree\\\\\\\\\\\":\\\\\\\\\\\"abc123\\\\\\\\\\\"\nskips[0]{package,test,reason}:\n## restored pass at c87faf66: bench test --package ./internal/gate -> pass, 447 tests\n"
          },
          "requirement": "preservation",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0,
          "probe": {
            "mutation": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "c87faf66abc2a8eb228b7179aeb57d0a075f9dd6",
              "digest": "sha256:a9bbe79d939f18df8f8d50516928a3268493a826d5217462ea90e594e6850093",
              "excerpt": "# preservation: drop the final newline from the verdict and store writers; independent expected bytes turn red; restored\n## 1791820544916413834-78a0c22da8180ec0\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,69cffaf6a7416abab014e5e1b352a683f3cbd566,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/verdict.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestVerdictReplacement,passed,1\nspilled{lines=11,bytes=951,omitted_lines=2,cut_lines=1,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791820544914541511-07ce19bc4ad140f2.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,4,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestVerdictReplacement,\"durable_replacement_test.go:30: record = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":1,\\\\\\\\\\\"state\\\\\\\\\\\":\\\\\\\\\\\"ready\\\\\\\\\\\",\\\\\\\\\\\"status\\\\\\\\\\\":\\\\\\\\\\\"green\\\\\\\\\\\",\\\\\\\\\\\"tree\\\\\\\\\\\":\\\\\\\\\\\"abc123\\\\\\\\\\\",\\\\\\\\\\\"oracle\\\\\\\\\\\":\\\\\\\\\\\"feed\\\\\\\\\\\",\\\\\\\\\\\"recorded_at\\\\\\\\\\\":\\\\\\\\\\\"2026-10-12T00:00:00Z\\\\\\\\\\\"}\\\", want \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":1,\\\\\\\\\\\"state\\\\\\\\\\\":\\\\\\\\\\\"ready\\\\\\\\\\\",\\\\\\\\\\\"s\nskips[0]{package,test,reason}:\n## 1791820618400582478-15c6df5b9c04a125\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,69cffaf6a7416abab014e5e1b352a683f3cbd566,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/engine.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestEvidenceReplacement,passed,3\nspilled{lines=12,bytes=1419,omitted_lines=3,cut_lines=2,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791820618398735692-39ebbefeaa7f6c09.out}\n  github.com/gibbonmi/bench/internal/gate,fail,4,3\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestEvidenceReplacement/evidence,\"durable_replacement_test.go:67: record = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":1,\\\\\\\\\\\"state\\\\\\\\\\\":\\\\\\\\\\\"ready\\\\\\\\\\\",\\\\\\\\\\\"status\\\\\\\\\\\":\\\\\\\\\\\"green\\\\\\\\\\\",\\\\\\\\\\\"tree\\\\\\\\\\\":\\\\\\\\\\\"abc123\\\\\\\\\\\",\\\\\\\\\\\"oracle\\\\\\\\\\\":\\\\\\\\\\\"feed\\\\\\\\\\\",\\\\\\\\\\\"recorded_at\\\\\\\\\\\":\\\\\\\\\\\"2026-10-12T00:00:00Z\\\\\\\\\\\"}\\\", want \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":1,\\\\\\\\\\\"state\\\\\\\\\\\":\\\\\\\\\\\"ready\\\\\\\\\n  github.com/gibbonmi/bench/internal/gate,TestEvidenceReplacement/lane,\"durable_replacement_test.go:67: record = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":1,\\\\\\\\\\\"tree\\\\\\\\\\\":\\\\\\\\\\\"abc123\\\\\\\\\\\",\\\\\\\\\\\"recorded_at\\\\\\\\\\\":\\\\\\\\\\\"2026-10-12T00:00:00Z\\\\\\\\\\\",\\\\\\\\\\\"lane\\\\\\\\\\\":\\\\\\\\\\\"fast\\\\\\\\\\\",\\\\\\\\\\\"outcome\\\\\\\\\\\":\\\\\\\\\\\"pass\\\\\\\\\\\",\\\\\\\\\\\"run_binary\\\\\\\\\\\":\\\\\\\\\\\"beef\\\\\\\\\\\"}\\\", want \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":1,\\\\\\\\\\\"tree\\\\\\\\\\\":\\\\\\\\\\\"abc123\\\\\\\\\\\"\nskips[0]{package,test,reason}:\n## restored pass at c87faf66: bench test --package ./internal/gate -> pass, 447 tests\n"
            }
          }
        },
        {
          "id": "dcm-t02-rep1-gate-tests",
          "performer": "claude:dcm_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2b201e4e9ad334e3d653d68e027c3c415b1aecf6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "e886531e4923b056e469e1fd9b06daa698b39b8e",
            "digest": "sha256:c30dfac506001eea59930ca80f59b1b7e188541a15103256a423c2054ac7c3af",
            "excerpt": "bench test --package ./internal/gate at e886531e4923b056e469e1fd9b06daa698b39b8e (repair tip; only reviews/durable-caller-migration.md dirty)\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,25482,450\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "gate-tests",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "dcm-t02-rep1-root-conformance",
          "performer": "claude:dcm_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2b201e4e9ad334e3d653d68e027c3c415b1aecf6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "e886531e4923b056e469e1fd9b06daa698b39b8e",
            "digest": "sha256:f9006bb4cd5bfa78a3569815a53cd03fd9f84dc60c00d649be7184159479b341",
            "excerpt": "bench test --check root-conformance at e886531e4923b056e469e1fd9b06daa698b39b8e (repair tip; only reviews/durable-caller-migration.md dirty)\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,18007,1\nfailures[0]{package,test,line,lines}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "dcm-t02-rep1-omission-proof",
          "performer": "claude:dcm_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2b201e4e9ad334e3d653d68e027c3c415b1aecf6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "e886531e4923b056e469e1fd9b06daa698b39b8e",
            "digest": "sha256:7745843b6533e95db1a2f97d494fe825fa98631f58c1c203bc20e183f38ef5b8",
            "excerpt": "# omission-proof at e886531e4923b056e469e1fd9b06daa698b39b8e: eight bench probe swaps, each bit and restored=yes\n## O1 leaf cause: run_durability.go swap 'return operational(root, gateExit, stderr, category+\": \"+err.Error())' with 'return operational(root, gateExit, stderr, category)'\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,e886531e4923b056e469e1fd9b06daa698b39b8e,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_durability.go,swap,failed,6,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestGatePendingDurabilityFailure|TestGateTerminalDurabilityFailure|TestGateEvidenceDurabilityFailure|TestGatePendingRecoveryFailure,passed,7\nspilled{lines=16,bytes=2684,omitted_lines=7,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791822030503880926-5022f7e0830f383d.out}\n  github.com/gibbonmi/bench/internal/gate,TestGatePendingRecoveryFailure,\"durable_replacement_test.go:169: stderr = \\\"gate final persistence failed\\\\\\\\n\\\", want \\\"gate final persistence failed: replace /dev/shm/ITHBFS/t/TestGatePendingRecoveryFailure1537287849/001/.git/bench-last-gate: directory-sync: injected final di… (480 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestGateTerminalDurabilityFailure/green,\"durable_replacement_test.go:120: stderr = \\\"gate final persistence failed\\\\\\\\n\\\", want \\\"gate final persistence failed: replace /dev/shm/ITHBFS/t/TestGateTerminalDurabilityFailuregreen3526316563/001/.git/bench-last-gate: directory-sync: injected … (355 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestGateTerminalDurabilityFailure/red,\"durable_replacement_test.go:120: stderr = \\\"gate final persistence failed\\\\\\\\n\\\", want \\\"gate final persistence failed: replace /dev/shm/ITHBFS/t/TestGateTerminalDurabilityFailurered3976651526/001/.git/bench-last-gate: directory-sync: injected te… (353 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestGateTerminalDurabilityFailure/timeout,\"durable_replacement_test.go:120: stderr = \\\"gate: timeout\\\\\\\\ngate timeout persistence failed\\\\\\\\n\\\", want \\\"gate timeout persistence failed: replace /dev/shm/ITHBFS/t/TestGateTerminalDurabilityFailuretimeout299311702/001/.git/bench-last-gate: direc… (375 bytes)\",1\nskips[0]{package,test,reason}:\n## O2 oracle start after pending failure: run_transaction.go swap 'return persistenceFailure(storageRoot, 0, stderr, \"gate pending persistence failed\", err)' with 'fmt.Fprintln(stderr, \"gate pending persistence failed: \"+err.Error())'\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,e886531e4923b056e469e1fd9b06daa698b39b8e,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_transaction.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestGatePendingDurabilityFailure,passed,1\nspilled{lines=11,bytes=699,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791822038339941802-9d920358da8c5aff.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,34,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGatePendingDurabilityFailure,\"durable_replacement_test.go:81: run counter stat = <nil>, want no oracle start after the pending failure\",1\nskips[0]{package,test,reason}:\n## O3 hidden evidence failure: run_transaction.go swap 'return persistenceFailure(storageRoot, rc, stderr, \"gate evidence persistence failed\", err)' with '_ = err'\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,e886531e4923b056e469e1fd9b06daa698b39b8e,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_transaction.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestGateEvidenceDurabilityFailure,passed,1\nspilled{lines=11,bytes=861,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791822045899192195-a519204e542bc7a9.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,39,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateEvidenceDurabilityFailure,\"durable_replacement_test.go:141: result = gate.Result{GateExit:0, ActionExit:0, Inspection:gate.Inspection{State:\\\"ready\\\", Status:\\\"green\\\", PendingStatus:\\\"\\\", CachedTree:\\\"94ee8e54d5b6e8382f5c399f026581643384395e\\\", CurrentTree:\\\"94ee8e54d5b6e838… (474 bytes)\",1\nskips[0]{package,test,reason}:\n## O4 discarded recovery error: run_durability.go swap 'return fmt.Errorf(\"%w; pending recovery failed: %w\", err, recovery)' with 'return fmt.Errorf(\"%w\", err)'\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,e886531e4923b056e469e1fd9b06daa698b39b8e,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_durability.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestGatePendingRecoveryFailure|TestGateDemotionDurabilityFailure,passed,4\nspilled{lines=13,bytes=1589,omitted_lines=4,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791822060620158943-7f97bf7553d3d777.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateDemotionDurabilityFailure/lock_contention,\"durable_replacement_test.go:221: stderr = \\\"gate execution already in progress; interrupted-green demotion failed: replace /dev/shm/ZEQ2QQ/t/TestGateDemotionDurabilityFailurelock_contention2739114308/001/.git/bench-last-gate: rename: injecte… (719 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestGateDemotionDurabilityFailure/lock_unavailable,\"durable_replacement_test.go:221: stderr = \\\"gate lock unavailable; interrupted-green demotion failed: replace /dev/shm/ZEQ2QQ/t/TestGateDemotionDurabilityFailurelock_unavailable1640129025/001/.git/bench-last-gate: rename: injected demotion r… (663 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestGatePendingRecoveryFailure,\"durable_replacement_test.go:169: stderr = \\\"gate final persistence failed: replace /dev/shm/ZEQ2QQ/t/TestGatePendingRecoveryFailure727242966/001/.git/bench-last-gate: directory-sync: injected final directory-sync failure; the rename complete… (684 bytes)\",1\nskips[0]{package,test,reason}:\n## C1 timeout cause (pre-repair at ce013f61 the same swap over the whole package was silent, 0 of 447 failed): run_transaction.go swap 'return persistenceFailure(storageRoot, 124, stderr, \"gate timeout persistence failed\", err)' with 'return operational(storageRoot, 124, stderr, \"gate timeout persistence failed\")'\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,e886531e4923b056e469e1fd9b06daa698b39b8e,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_transaction.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestGateTerminalDurabilityFailure,passed,4\nspilled{lines=11,bytes=870,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791822067069652853-096cce86b2304a05.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,212,4\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateTerminalDurabilityFailure/timeout,\"durable_replacement_test.go:120: stderr = \\\"gate: timeout\\\\\\\\ngate timeout persistence failed\\\\\\\\n\\\", want \\\"gate timeout persistence failed: replace /dev/shm/5B7WIQ/t/TestGateTerminalDurabilityFailuretimeout793381759/001/.git/bench-last-gate: direc… (375 bytes)\",1\nskips[0]{package,test,reason}:\n## C2 contention demotion (pre-repair at ce013f61 the same swap over the whole package was silent, 0 of 447 failed): run_transaction.go swap 'fmt.Fprintln(stderr, besideDemotion(\"gate execution already in progress\", persistInterruptedIfGreen(replacer, storageRoot, gitdir, plan)))' with 'fmt.Fprintln(stderr, besideDemotion(\"gate execution already in progress\", func() error { _ = persistInterruptedIfGreen(replacer, storageRoot, gitdir, plan); return nil }()))'\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,e886531e4923b056e469e1fd9b06daa698b39b8e,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_transaction.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestGateDemotionDurabilityFailure,passed,3\nspilled{lines=11,bytes=878,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791822070776722175-959b92163759d8f4.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,125,3\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateDemotionDurabilityFailure/lock_contention,\"durable_replacement_test.go:221: stderr = \\\"gate execution already in progress\\\\\\\\ngate owner: pid 2319563 (alive)\\\\\\\\n\\\", want \\\"gate execution already in progress; interrupted-green demotion failed: replace /dev/shm/3FUDLS/t/TestGateDemotionDurabi… (531 bytes)\",1\nskips[0]{package,test,reason}:\n## C3 D39 cache-replacement count: run_transaction.go swap 'return persistenceFailure(storageRoot, rc, stderr, \"gate evidence persistence failed\", err)' with '_ = persistRecovering(replacer, gitdir, ready, pending); return persistenceFailure(storageRoot, rc, stderr, \"gate evidence persistence failed\", err)'\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,e886531e4923b056e469e1fd9b06daa698b39b8e,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_transaction.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestGateEvidenceDurabilityFailure,passed,1\nspilled{lines=11,bytes=682,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791822078836771803-5f857aec59edffd5.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,38,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateEvidenceDurabilityFailure,\"durable_replacement_test.go:147: cache replacements = 2, want only the pending record\",1\nskips[0]{package,test,reason}:\n## C3 D40 visible-green inspection: run_durability.go swap 'return fmt.Errorf(\"%w; pending recovery failed: %w\", err, recovery)' with '_ = durableReplace(durablefile.Native(), gitdir, pending); return fmt.Errorf(\"%w; pending recovery failed: %w\", err, recovery)'\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,e886531e4923b056e469e1fd9b06daa698b39b8e,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_durability.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestGatePendingRecoveryFailure,passed,1\nspilled{lines=11,bytes=854,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791822086501993626-95148c6b77aef22f.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,47,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGatePendingRecoveryFailure,\"durable_replacement_test.go:172: inspection after the failed recovery = gate.Inspection{State:\\\"pending\\\", Status:\\\"\\\", PendingStatus:\\\"interrupted-pending\\\", CachedTree:\\\"089de20a152f10a33d06130cce91b0f368f4598c\\\", CurrentTree:\\\"089de20a152f10a33d0… (422 bytes)\",1\nskips[0]{package,test,reason}:\n## restored pass at e886531e4923b056e469e1fd9b06daa698b39b8e: bench test --package ./internal/gate -> pass, 450 tests\n"
          },
          "requirement": "omission-proof",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0,
          "probe": {
            "mutation": "Discard the leaf cause, allow oracle start after pending failure, hide evidence failure, or discard the recovery error. Each mutation must fail its real execution witness.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "e886531e4923b056e469e1fd9b06daa698b39b8e",
              "digest": "sha256:7745843b6533e95db1a2f97d494fe825fa98631f58c1c203bc20e183f38ef5b8",
              "excerpt": "# omission-proof at e886531e4923b056e469e1fd9b06daa698b39b8e: eight bench probe swaps, each bit and restored=yes\n## O1 leaf cause: run_durability.go swap 'return operational(root, gateExit, stderr, category+\": \"+err.Error())' with 'return operational(root, gateExit, stderr, category)'\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,e886531e4923b056e469e1fd9b06daa698b39b8e,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_durability.go,swap,failed,6,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestGatePendingDurabilityFailure|TestGateTerminalDurabilityFailure|TestGateEvidenceDurabilityFailure|TestGatePendingRecoveryFailure,passed,7\nspilled{lines=16,bytes=2684,omitted_lines=7,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791822030503880926-5022f7e0830f383d.out}\n  github.com/gibbonmi/bench/internal/gate,TestGatePendingRecoveryFailure,\"durable_replacement_test.go:169: stderr = \\\"gate final persistence failed\\\\\\\\n\\\", want \\\"gate final persistence failed: replace /dev/shm/ITHBFS/t/TestGatePendingRecoveryFailure1537287849/001/.git/bench-last-gate: directory-sync: injected final di… (480 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestGateTerminalDurabilityFailure/green,\"durable_replacement_test.go:120: stderr = \\\"gate final persistence failed\\\\\\\\n\\\", want \\\"gate final persistence failed: replace /dev/shm/ITHBFS/t/TestGateTerminalDurabilityFailuregreen3526316563/001/.git/bench-last-gate: directory-sync: injected … (355 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestGateTerminalDurabilityFailure/red,\"durable_replacement_test.go:120: stderr = \\\"gate final persistence failed\\\\\\\\n\\\", want \\\"gate final persistence failed: replace /dev/shm/ITHBFS/t/TestGateTerminalDurabilityFailurered3976651526/001/.git/bench-last-gate: directory-sync: injected te… (353 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestGateTerminalDurabilityFailure/timeout,\"durable_replacement_test.go:120: stderr = \\\"gate: timeout\\\\\\\\ngate timeout persistence failed\\\\\\\\n\\\", want \\\"gate timeout persistence failed: replace /dev/shm/ITHBFS/t/TestGateTerminalDurabilityFailuretimeout299311702/001/.git/bench-last-gate: direc… (375 bytes)\",1\nskips[0]{package,test,reason}:\n## O2 oracle start after pending failure: run_transaction.go swap 'return persistenceFailure(storageRoot, 0, stderr, \"gate pending persistence failed\", err)' with 'fmt.Fprintln(stderr, \"gate pending persistence failed: \"+err.Error())'\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,e886531e4923b056e469e1fd9b06daa698b39b8e,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_transaction.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestGatePendingDurabilityFailure,passed,1\nspilled{lines=11,bytes=699,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791822038339941802-9d920358da8c5aff.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,34,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGatePendingDurabilityFailure,\"durable_replacement_test.go:81: run counter stat = <nil>, want no oracle start after the pending failure\",1\nskips[0]{package,test,reason}:\n## O3 hidden evidence failure: run_transaction.go swap 'return persistenceFailure(storageRoot, rc, stderr, \"gate evidence persistence failed\", err)' with '_ = err'\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,e886531e4923b056e469e1fd9b06daa698b39b8e,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_transaction.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestGateEvidenceDurabilityFailure,passed,1\nspilled{lines=11,bytes=861,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791822045899192195-a519204e542bc7a9.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,39,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateEvidenceDurabilityFailure,\"durable_replacement_test.go:141: result = gate.Result{GateExit:0, ActionExit:0, Inspection:gate.Inspection{State:\\\"ready\\\", Status:\\\"green\\\", PendingStatus:\\\"\\\", CachedTree:\\\"94ee8e54d5b6e8382f5c399f026581643384395e\\\", CurrentTree:\\\"94ee8e54d5b6e838… (474 bytes)\",1\nskips[0]{package,test,reason}:\n## O4 discarded recovery error: run_durability.go swap 'return fmt.Errorf(\"%w; pending recovery failed: %w\", err, recovery)' with 'return fmt.Errorf(\"%w\", err)'\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,e886531e4923b056e469e1fd9b06daa698b39b8e,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_durability.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestGatePendingRecoveryFailure|TestGateDemotionDurabilityFailure,passed,4\nspilled{lines=13,bytes=1589,omitted_lines=4,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791822060620158943-7f97bf7553d3d777.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateDemotionDurabilityFailure/lock_contention,\"durable_replacement_test.go:221: stderr = \\\"gate execution already in progress; interrupted-green demotion failed: replace /dev/shm/ZEQ2QQ/t/TestGateDemotionDurabilityFailurelock_contention2739114308/001/.git/bench-last-gate: rename: injecte… (719 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestGateDemotionDurabilityFailure/lock_unavailable,\"durable_replacement_test.go:221: stderr = \\\"gate lock unavailable; interrupted-green demotion failed: replace /dev/shm/ZEQ2QQ/t/TestGateDemotionDurabilityFailurelock_unavailable1640129025/001/.git/bench-last-gate: rename: injected demotion r… (663 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestGatePendingRecoveryFailure,\"durable_replacement_test.go:169: stderr = \\\"gate final persistence failed: replace /dev/shm/ZEQ2QQ/t/TestGatePendingRecoveryFailure727242966/001/.git/bench-last-gate: directory-sync: injected final directory-sync failure; the rename complete… (684 bytes)\",1\nskips[0]{package,test,reason}:\n## C1 timeout cause (pre-repair at ce013f61 the same swap over the whole package was silent, 0 of 447 failed): run_transaction.go swap 'return persistenceFailure(storageRoot, 124, stderr, \"gate timeout persistence failed\", err)' with 'return operational(storageRoot, 124, stderr, \"gate timeout persistence failed\")'\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,e886531e4923b056e469e1fd9b06daa698b39b8e,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_transaction.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestGateTerminalDurabilityFailure,passed,4\nspilled{lines=11,bytes=870,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791822067069652853-096cce86b2304a05.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,212,4\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateTerminalDurabilityFailure/timeout,\"durable_replacement_test.go:120: stderr = \\\"gate: timeout\\\\\\\\ngate timeout persistence failed\\\\\\\\n\\\", want \\\"gate timeout persistence failed: replace /dev/shm/5B7WIQ/t/TestGateTerminalDurabilityFailuretimeout793381759/001/.git/bench-last-gate: direc… (375 bytes)\",1\nskips[0]{package,test,reason}:\n## C2 contention demotion (pre-repair at ce013f61 the same swap over the whole package was silent, 0 of 447 failed): run_transaction.go swap 'fmt.Fprintln(stderr, besideDemotion(\"gate execution already in progress\", persistInterruptedIfGreen(replacer, storageRoot, gitdir, plan)))' with 'fmt.Fprintln(stderr, besideDemotion(\"gate execution already in progress\", func() error { _ = persistInterruptedIfGreen(replacer, storageRoot, gitdir, plan); return nil }()))'\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,e886531e4923b056e469e1fd9b06daa698b39b8e,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_transaction.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestGateDemotionDurabilityFailure,passed,3\nspilled{lines=11,bytes=878,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791822070776722175-959b92163759d8f4.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,125,3\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateDemotionDurabilityFailure/lock_contention,\"durable_replacement_test.go:221: stderr = \\\"gate execution already in progress\\\\\\\\ngate owner: pid 2319563 (alive)\\\\\\\\n\\\", want \\\"gate execution already in progress; interrupted-green demotion failed: replace /dev/shm/3FUDLS/t/TestGateDemotionDurabi… (531 bytes)\",1\nskips[0]{package,test,reason}:\n## C3 D39 cache-replacement count: run_transaction.go swap 'return persistenceFailure(storageRoot, rc, stderr, \"gate evidence persistence failed\", err)' with '_ = persistRecovering(replacer, gitdir, ready, pending); return persistenceFailure(storageRoot, rc, stderr, \"gate evidence persistence failed\", err)'\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,e886531e4923b056e469e1fd9b06daa698b39b8e,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_transaction.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestGateEvidenceDurabilityFailure,passed,1\nspilled{lines=11,bytes=682,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791822078836771803-5f857aec59edffd5.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,38,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateEvidenceDurabilityFailure,\"durable_replacement_test.go:147: cache replacements = 2, want only the pending record\",1\nskips[0]{package,test,reason}:\n## C3 D40 visible-green inspection: run_durability.go swap 'return fmt.Errorf(\"%w; pending recovery failed: %w\", err, recovery)' with '_ = durableReplace(durablefile.Native(), gitdir, pending); return fmt.Errorf(\"%w; pending recovery failed: %w\", err, recovery)'\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,e886531e4923b056e469e1fd9b06daa698b39b8e,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_durability.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestGatePendingRecoveryFailure,passed,1\nspilled{lines=11,bytes=854,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791822086501993626-95148c6b77aef22f.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,47,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGatePendingRecoveryFailure,\"durable_replacement_test.go:172: inspection after the failed recovery = gate.Inspection{State:\\\"pending\\\", Status:\\\"\\\", PendingStatus:\\\"interrupted-pending\\\", CachedTree:\\\"089de20a152f10a33d06130cce91b0f368f4598c\\\", CurrentTree:\\\"089de20a152f10a33d0… (422 bytes)\",1\nskips[0]{package,test,reason}:\n## restored pass at e886531e4923b056e469e1fd9b06daa698b39b8e: bench test --package ./internal/gate -> pass, 450 tests\n"
            }
          }
        },
        {
          "id": "dcm-t02-rep1-preservation",
          "performer": "claude:dcm_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2b201e4e9ad334e3d653d68e027c3c415b1aecf6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "e886531e4923b056e469e1fd9b06daa698b39b8e",
            "digest": "sha256:f735c9993df15146247359ce44502bef5d8594356e194ddb7ca68ade5bcd24a5",
            "excerpt": "# preservation at e886531e4923b056e469e1fd9b06daa698b39b8e: four bench probe swaps on the one store writer, each bit and restored=yes\n## S1 loop (pre-repair at ce013f61): engine.go swap \"append(data, '\\n'), 0o600)\" with \"append(data, '\\n'), 0o644)\" left TestVerdictReplacement passing (silent, 0 failed), because the verdict writer held its own 0o600; at the repair tip the same swap turns it red (C3 mode row below)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,ce013f61c673796c3dbca37926960364c9475214,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  silent,internal/gate/engine.go,swap,passed,0,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestVerdictReplacement,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,1,0,170\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,4,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nworktree: /home/mgibs/.bench/worktrees/bench-2826441890/5edaf1c9a73f8054055267cdcca55bdc-00b4e1c8d78f5645bd35ac462e410538\n## P1 final newline: engine.go swap \"append(data, '\\n'), 0o600)\" with \"data, 0o600)\"\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,e886531e4923b056e469e1fd9b06daa698b39b8e,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/engine.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestVerdictReplacement|TestEvidenceReplacement,passed,4\nspilled{lines=13,bytes=1897,omitted_lines=4,cut_lines=3,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791822094191991909-2d98ec173cd77027.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestEvidenceReplacement/evidence,\"durable_replacement_test.go:67: record = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":1,\\\\\\\\\\\"state\\\\\\\\\\\":\\\\\\\\\\\"ready\\\\\\\\\\\",\\\\\\\\\\\"status\\\\\\\\\\\":\\\\\\\\\\\"green\\\\\\\\\\\",\\\\\\\\\\\"tree\\\\\\\\\\\":\\\\\\\\\\\"abc123\\\\\\\\\\\",\\\\\\\\\\\"oracle\\\\\\\\\\\":\\\\\\\\\\\"feed\\\\\\\\\\\",\\\\\\\\\\\"recorded_at\\\\\\\\\\\":\\\\\\\\\\\"2026-10-12T00:00:00Z\\\\\\\\\\\"}\\\", want \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":1,\\\\\\\\\\\"state\\\\\\\\\\\":\\\\\\\\\\\"ready\\\\\\\\\n  github.com/gibbonmi/bench/internal/gate,TestEvidenceReplacement/lane,\"durable_replacement_test.go:67: record = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":1,\\\\\\\\\\\"tree\\\\\\\\\\\":\\\\\\\\\\\"abc123\\\\\\\\\\\",\\\\\\\\\\\"recorded_at\\\\\\\\\\\":\\\\\\\\\\\"2026-10-12T00:00:00Z\\\\\\\\\\\",\\\\\\\\\\\"lane\\\\\\\\\\\":\\\\\\\\\\\"fast\\\\\\\\\\\",\\\\\\\\\\\"outcome\\\\\\\\\\\":\\\\\\\\\\\"pass\\\\\\\\\\\",\\\\\\\\\\\"run_binary\\\\\\\\\\\":\\\\\\\\\\\"beef\\\\\\\\\\\"}\\\", want \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":1,\\\\\\\\\\\"tree\\\\\\\\\\\":\\\\\\\\\\\"abc123\\\\\\\\\\\"\n  github.com/gibbonmi/bench/internal/gate,TestVerdictReplacement,\"durable_replacement_test.go:30: record = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":1,\\\\\\\\\\\"state\\\\\\\\\\\":\\\\\\\\\\\"ready\\\\\\\\\\\",\\\\\\\\\\\"status\\\\\\\\\\\":\\\\\\\\\\\"green\\\\\\\\\\\",\\\\\\\\\\\"tree\\\\\\\\\\\":\\\\\\\\\\\"abc123\\\\\\\\\\\",\\\\\\\\\\\"oracle\\\\\\\\\\\":\\\\\\\\\\\"feed\\\\\\\\\\\",\\\\\\\\\\\"recorded_at\\\\\\\\\\\":\\\\\\\\\\\"2026-10-12T00:00:00Z\\\\\\\\\\\"}\\\", want \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":1,\\\\\\\\\\\"state\\\\\\\\\\\":\\\\\\\\\\\"ready\\\\\\\\\\\",\\\\\\\\\\\"s\nskips[0]{package,test,reason}:\n## C3 0600 mode literal (verdict and store writers share engine.go): engine.go swap \"append(data, '\\n'), 0o600)\" with \"append(data, '\\n'), 0o644)\"\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,e886531e4923b056e469e1fd9b06daa698b39b8e,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/engine.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestVerdictReplacement|TestEvidenceReplacement,passed,4\nspilled{lines=13,bytes=925,omitted_lines=4,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791822098087056210-0e5ca39493621b04.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestEvidenceReplacement/evidence,\"durable_replacement_test.go:67: mode = -rw-r--r--, want 0600\",1\n  github.com/gibbonmi/bench/internal/gate,TestEvidenceReplacement/lane,\"durable_replacement_test.go:67: mode = -rw-r--r--, want 0600\",1\n  github.com/gibbonmi/bench/internal/gate,TestVerdictReplacement,\"durable_replacement_test.go:30: mode = -rw-r--r--, want 0600\",1\nskips[0]{package,test,reason}:\n## C3 store-writer owner bypass: engine.go swap \"return replacer.Replace(filepath.Join(dir, name), append(data, '\\n'), 0o600)\" with \"return durablefile.Replace(filepath.Join(dir, name), append(data, '\\n'), 0o600)\" (after S1 this is also the verdict writer's only Replace call)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,e886531e4923b056e469e1fd9b06daa698b39b8e,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/engine.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestVerdictReplacement|TestEvidenceReplacement,passed,4\nspilled{lines=13,bytes=1181,omitted_lines=4,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791822103668687342-ba717c1517f09d07.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestEvidenceReplacement/evidence,\"durable_replacement_test.go:66: err = <nil>, want a published directory-sync *durablefile.Error that wraps injected evidence directory-sync failure\",1\n  github.com/gibbonmi/bench/internal/gate,TestEvidenceReplacement/lane,\"durable_replacement_test.go:66: err = <nil>, want a published directory-sync *durablefile.Error that wraps injected lane directory-sync failure\",1\n  github.com/gibbonmi/bench/internal/gate,TestVerdictReplacement,\"durable_replacement_test.go:29: err = <nil>, want a published directory-sync *durablefile.Error that wraps injected verdict directory-sync failure\",1\nskips[0]{package,test,reason}:\n## C3 verdict-writer owner bypass (the verdict.go replacer.Replace site no longer exists after S1; the bypass is the replacer it forwards): verdict.go swap 'return durableReplaceAt(replacer, gitdir, benchgit.GateCacheFile, rec)' with 'return durableReplaceAt(durablefile.Native(), gitdir, benchgit.GateCacheFile, rec)'\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,e886531e4923b056e469e1fd9b06daa698b39b8e,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/verdict.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestVerdictReplacement,passed,1\nspilled{lines=11,bytes=712,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791822109312385924-6f59a1e258854444.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,4,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestVerdictReplacement,\"durable_replacement_test.go:29: err = <nil>, want a published directory-sync *durablefile.Error that wraps injected verdict directory-sync failure\",1\nskips[0]{package,test,reason}:\n## restored pass at e886531e4923b056e469e1fd9b06daa698b39b8e: bench test --package ./internal/gate -> pass, 450 tests\n"
          },
          "requirement": "preservation",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0,
          "probe": {
            "mutation": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "e886531e4923b056e469e1fd9b06daa698b39b8e",
              "digest": "sha256:f735c9993df15146247359ce44502bef5d8594356e194ddb7ca68ade5bcd24a5",
              "excerpt": "# preservation at e886531e4923b056e469e1fd9b06daa698b39b8e: four bench probe swaps on the one store writer, each bit and restored=yes\n## S1 loop (pre-repair at ce013f61): engine.go swap \"append(data, '\\n'), 0o600)\" with \"append(data, '\\n'), 0o644)\" left TestVerdictReplacement passing (silent, 0 failed), because the verdict writer held its own 0o600; at the repair tip the same swap turns it red (C3 mode row below)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,ce013f61c673796c3dbca37926960364c9475214,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  silent,internal/gate/engine.go,swap,passed,0,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestVerdictReplacement,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,1,0,170\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,4,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nworktree: /home/mgibs/.bench/worktrees/bench-2826441890/5edaf1c9a73f8054055267cdcca55bdc-00b4e1c8d78f5645bd35ac462e410538\n## P1 final newline: engine.go swap \"append(data, '\\n'), 0o600)\" with \"data, 0o600)\"\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,e886531e4923b056e469e1fd9b06daa698b39b8e,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/engine.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestVerdictReplacement|TestEvidenceReplacement,passed,4\nspilled{lines=13,bytes=1897,omitted_lines=4,cut_lines=3,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791822094191991909-2d98ec173cd77027.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestEvidenceReplacement/evidence,\"durable_replacement_test.go:67: record = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":1,\\\\\\\\\\\"state\\\\\\\\\\\":\\\\\\\\\\\"ready\\\\\\\\\\\",\\\\\\\\\\\"status\\\\\\\\\\\":\\\\\\\\\\\"green\\\\\\\\\\\",\\\\\\\\\\\"tree\\\\\\\\\\\":\\\\\\\\\\\"abc123\\\\\\\\\\\",\\\\\\\\\\\"oracle\\\\\\\\\\\":\\\\\\\\\\\"feed\\\\\\\\\\\",\\\\\\\\\\\"recorded_at\\\\\\\\\\\":\\\\\\\\\\\"2026-10-12T00:00:00Z\\\\\\\\\\\"}\\\", want \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":1,\\\\\\\\\\\"state\\\\\\\\\\\":\\\\\\\\\\\"ready\\\\\\\\\n  github.com/gibbonmi/bench/internal/gate,TestEvidenceReplacement/lane,\"durable_replacement_test.go:67: record = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":1,\\\\\\\\\\\"tree\\\\\\\\\\\":\\\\\\\\\\\"abc123\\\\\\\\\\\",\\\\\\\\\\\"recorded_at\\\\\\\\\\\":\\\\\\\\\\\"2026-10-12T00:00:00Z\\\\\\\\\\\",\\\\\\\\\\\"lane\\\\\\\\\\\":\\\\\\\\\\\"fast\\\\\\\\\\\",\\\\\\\\\\\"outcome\\\\\\\\\\\":\\\\\\\\\\\"pass\\\\\\\\\\\",\\\\\\\\\\\"run_binary\\\\\\\\\\\":\\\\\\\\\\\"beef\\\\\\\\\\\"}\\\", want \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":1,\\\\\\\\\\\"tree\\\\\\\\\\\":\\\\\\\\\\\"abc123\\\\\\\\\\\"\n  github.com/gibbonmi/bench/internal/gate,TestVerdictReplacement,\"durable_replacement_test.go:30: record = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":1,\\\\\\\\\\\"state\\\\\\\\\\\":\\\\\\\\\\\"ready\\\\\\\\\\\",\\\\\\\\\\\"status\\\\\\\\\\\":\\\\\\\\\\\"green\\\\\\\\\\\",\\\\\\\\\\\"tree\\\\\\\\\\\":\\\\\\\\\\\"abc123\\\\\\\\\\\",\\\\\\\\\\\"oracle\\\\\\\\\\\":\\\\\\\\\\\"feed\\\\\\\\\\\",\\\\\\\\\\\"recorded_at\\\\\\\\\\\":\\\\\\\\\\\"2026-10-12T00:00:00Z\\\\\\\\\\\"}\\\", want \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":1,\\\\\\\\\\\"state\\\\\\\\\\\":\\\\\\\\\\\"ready\\\\\\\\\\\",\\\\\\\\\\\"s\nskips[0]{package,test,reason}:\n## C3 0600 mode literal (verdict and store writers share engine.go): engine.go swap \"append(data, '\\n'), 0o600)\" with \"append(data, '\\n'), 0o644)\"\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,e886531e4923b056e469e1fd9b06daa698b39b8e,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/engine.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestVerdictReplacement|TestEvidenceReplacement,passed,4\nspilled{lines=13,bytes=925,omitted_lines=4,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791822098087056210-0e5ca39493621b04.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestEvidenceReplacement/evidence,\"durable_replacement_test.go:67: mode = -rw-r--r--, want 0600\",1\n  github.com/gibbonmi/bench/internal/gate,TestEvidenceReplacement/lane,\"durable_replacement_test.go:67: mode = -rw-r--r--, want 0600\",1\n  github.com/gibbonmi/bench/internal/gate,TestVerdictReplacement,\"durable_replacement_test.go:30: mode = -rw-r--r--, want 0600\",1\nskips[0]{package,test,reason}:\n## C3 store-writer owner bypass: engine.go swap \"return replacer.Replace(filepath.Join(dir, name), append(data, '\\n'), 0o600)\" with \"return durablefile.Replace(filepath.Join(dir, name), append(data, '\\n'), 0o600)\" (after S1 this is also the verdict writer's only Replace call)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,e886531e4923b056e469e1fd9b06daa698b39b8e,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/engine.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestVerdictReplacement|TestEvidenceReplacement,passed,4\nspilled{lines=13,bytes=1181,omitted_lines=4,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791822103668687342-ba717c1517f09d07.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestEvidenceReplacement/evidence,\"durable_replacement_test.go:66: err = <nil>, want a published directory-sync *durablefile.Error that wraps injected evidence directory-sync failure\",1\n  github.com/gibbonmi/bench/internal/gate,TestEvidenceReplacement/lane,\"durable_replacement_test.go:66: err = <nil>, want a published directory-sync *durablefile.Error that wraps injected lane directory-sync failure\",1\n  github.com/gibbonmi/bench/internal/gate,TestVerdictReplacement,\"durable_replacement_test.go:29: err = <nil>, want a published directory-sync *durablefile.Error that wraps injected verdict directory-sync failure\",1\nskips[0]{package,test,reason}:\n## C3 verdict-writer owner bypass (the verdict.go replacer.Replace site no longer exists after S1; the bypass is the replacer it forwards): verdict.go swap 'return durableReplaceAt(replacer, gitdir, benchgit.GateCacheFile, rec)' with 'return durableReplaceAt(durablefile.Native(), gitdir, benchgit.GateCacheFile, rec)'\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,e886531e4923b056e469e1fd9b06daa698b39b8e,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/verdict.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestVerdictReplacement,passed,1\nspilled{lines=11,bytes=712,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791822109312385924-6f59a1e258854444.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,4,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestVerdictReplacement,\"durable_replacement_test.go:29: err = <nil>, want a published directory-sync *durablefile.Error that wraps injected verdict directory-sync failure\",1\nskips[0]{package,test,reason}:\n## restored pass at e886531e4923b056e469e1fd9b06daa698b39b8e: bench test --package ./internal/gate -> pass, 450 tests\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "dcm-b1-r1-standards",
          "performer": "claude:dcm_rev_b1_std",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "2c57c2d6e8525a84409078f842647f628de3faa9",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b1_std",
            "digest": "sha256:d8cbab0194337ce7fdcc9cdc95becb3467a5f20937e3e45d42abe8bc25f62486",
            "excerpt": "## Standards (D-B1, 69cffaf6..90d67d63) — 1 finding; worst S1\nS1 | auto-fix | 8 | internal/gate/verdict.go:95 | durableReplace repeats json.Marshal + '\\n' + 0o600 + filepath.Join(dir,name), the same encoding and mode as durableReplaceAt/durableReplaceRecordAt (internal/gate/engine.go:173-187); the comments at engine.go:182 and lane.go:317 say every record class publishes through the store writer, which is false for the verdict cache | make durableReplace return durableReplaceAt(replacer, gitdir, benchgit.GateCacheFile, rec); keep the TestVerdictReplacement preservation red\nDeclared deviations 1-4 and 7: accepted. 5-6: not this axis.\nAdvice: the verdictRecord literal and expected JSON appear twice in the test; a production constant for bench-gate-evidence; executeSubjectReplacing has 11 parameters; one named 0600 mode.\n"
          },
          "axis": "Standards",
          "base": "69cffaf6a7416abab014e5e1b352a683f3cbd566",
          "tip": "c87faf66abc2a8eb228b7179aeb57d0a075f9dd6",
          "finding_ids": [
            "S1"
          ],
          "supersedes": []
        },
        {
          "id": "dcm-b1-r1-spec",
          "performer": "claude:dcm_rev_b1_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "2c57c2d6e8525a84409078f842647f628de3faa9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b1_spec",
            "digest": "sha256:6bee411119beae598cf34dcd7aed8a00dc70ced51c8310dbb2f78aee91fd29d9",
            "excerpt": "## Spec (D-B1, 69cffaf6..90d67d63) — 0 findings\nD16, D17, D37, D38, D39, D40: each behavior and seam holds against the spec text and the tip source. ActionExit 1 and non-reusable Inspection through operational (gate.go:335-341). Demotion returns its error to both callers (run_transaction.go:79, :83). Seams: the native replacer reaches executeSubjectReplacing on every production route; faults key by destination and rename ordinal. run_transaction.go is 346 lines. Coverage-map seam cells cite the real tests.\nAdvice: the native replacer comes from the unexported wrapper; the lock-contention caller, the timeout cause, and runLane's %w propagation have no execution witness; the omission-proof excerpt names subjects but not each swap.\n"
          },
          "axis": "Spec",
          "base": "69cffaf6a7416abab014e5e1b352a683f3cbd566",
          "tip": "c87faf66abc2a8eb228b7179aeb57d0a075f9dd6",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "dcm-b1-r1-coverage",
          "performer": "claude:dcm_rev_b1_cov",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "2c57c2d6e8525a84409078f842647f628de3faa9",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b1_cov",
            "digest": "sha256:e3a8e8bbeef7041e2586c1b5b76074e9bfb52dfe07ab87c0286cec5377c00442",
            "excerpt": "## Coverage (D-B1, 69cffaf6..90d67d63) — 3 findings; worst C1. 10 bench probe runs, each restored=yes.\nC1 | auto-fix | 9 | internal/gate/run_transaction.go:186-187 | no test fails the timeout persistence branch; swapping persistenceFailure(...,\"gate timeout persistence failed\", err) for operational(...,\"gate timeout persistence failed\") is silent (0 failed) | add an execution test with a short gateTimeout and a cache fault at ordinal 2 asserting the classified cause, GateExit 124 / ActionExit 1, and pending after recovery\nC2 | auto-fix | 9 | internal/gate/run_transaction.go:83 | only the lock-failure demotion caller is tested; discarding the demotion error at the contention caller is silent (0 failed) | add a contention case from a reusable green with a cache fault at ordinal 1 asserting \"gate execution already in progress; interrupted-green demotion failed: ...\"\nC3 | auto-fix | 7 | reviews/durable-caller-migration.md:62,91 | independent expectations without a recorded red: the 0600 mode literal (durable_replacement_test.go:277), the D39 cache count (:137), the D40 visible-green inspection (:162), and the per-family owner-bypass mutation (spec testing decisions) | record the reds: verdict.go 0600->0644, engine.go 0600->0644, verdict.go replacer.Replace->os.WriteFile, engine.go replacer.Replace->durablefile.Replace, final write after evidence failure, silent native re-rollback after failed recovery (all bit, restored)\nAdvice: record exact swap text; runLane persistence-error return removal is silent (pre-existing); bare durableReplace on the pending branch is silent (no row).\n"
          },
          "axis": "Coverage",
          "base": "69cffaf6a7416abab014e5e1b352a683f3cbd566",
          "tip": "c87faf66abc2a8eb228b7179aeb57d0a075f9dd6",
          "finding_ids": [
            "C1",
            "C2",
            "C3"
          ],
          "supersedes": []
        },
        {
          "id": "dcm-b1-r2-standards",
          "performer": "claude:dcm_rev_b1_std2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "2b201e4e9ad334e3d653d68e027c3c415b1aecf6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b1_std2",
            "digest": "sha256:dfa11d09f99f619139d8721cd830e11ea8954525651ec9b38b106fb2c0a25dca",
            "excerpt": "## Standards confirming round (D-B1, repair delta ce013f61..9c9ae1c0) — 0 findings\nS1 closed, confidence 9: durableReplace calls durableReplaceAt(replacer, gitdir, benchgit.GateCacheFile, rec) at internal/gate/verdict.go:91; the encoding lives only at engine.go:174-178 and the newline and 0o600 only at engine.go:186; the comments at engine.go:182-184 and lane.go:316-317 are now true. The new test rows assert production literals with recorded reds.\nAdvice: a withGateTimeout test helper for the four gateTimeout overrides; lane.go could call durableReplaceAt directly.\n"
          },
          "axis": "Standards",
          "base": "69cffaf6a7416abab014e5e1b352a683f3cbd566",
          "tip": "e886531e4923b056e469e1fd9b06daa698b39b8e",
          "finding_ids": [],
          "supersedes": [
            "dcm-b1-r1-standards"
          ]
        },
        {
          "id": "dcm-b1-r2-spec",
          "performer": "claude:dcm_rev_b1_spec2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "2b201e4e9ad334e3d653d68e027c3c415b1aecf6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b1_spec2",
            "digest": "sha256:a53e485fee16c9ad4ac6d977c448d1bb87064cbbbe108f37684a693c34d0f708",
            "excerpt": "## Spec confirming round (D-B1, repair delta ce013f61..9c9ae1c0) — 0 findings\nThe repair delta keeps D16, D17, D37-D40 and the gate seam rules intact. The timeout row meets spec.md:95 and :98; the lock-contention row meets the two-caller demotion rule. The repair touches only verdict.go and durable_replacement_test.go inside ticket 02 Writes.\nAdvice: the gateTimeout swap is a timing knob, not the replacement seam; the D38 why-clause names only green and red branches.\n"
          },
          "axis": "Spec",
          "base": "69cffaf6a7416abab014e5e1b352a683f3cbd566",
          "tip": "e886531e4923b056e469e1fd9b06daa698b39b8e",
          "finding_ids": [],
          "supersedes": [
            "dcm-b1-r1-spec"
          ]
        },
        {
          "id": "dcm-b1-r2-coverage",
          "performer": "claude:dcm_rev_b1_cov2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "2b201e4e9ad334e3d653d68e027c3c415b1aecf6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b1_cov2",
            "digest": "sha256:31a0021f89e5574b1588c2fc94133e96350bba5696b219cea5ea8b5c611a135e",
            "excerpt": "## Coverage confirming round (D-B1, repair delta ce013f61..9c9ae1c0) — 0 findings\nC1 closed: timeout row at durable_replacement_test.go:100; probe run_transaction.go:187 persistenceFailure(storageRoot, 124 -> 0 ...) bit, restored=yes.\nC2 closed: lock contention row at durable_replacement_test.go:196-200; probe run_transaction.go:83 replacer -> durablefile.Native() bit, restored=yes.\nC3 closed: dcm-t02-rep1-preservation and -omission-proof record exact swap text; probe lane.go:332 durableReplaceRecordAt(replacer -> durablefile.Native() ...) bit in TestEvidenceReplacement/lane, restored=yes.\nAdvice: the preservation excerpt heading count is off by one; the gateTimeout override needs non-parallel tests.\n"
          },
          "axis": "Coverage",
          "base": "69cffaf6a7416abab014e5e1b352a683f3cbd566",
          "tip": "e886531e4923b056e469e1fd9b06daa698b39b8e",
          "finding_ids": [],
          "supersedes": [
            "dcm-b1-r1-coverage"
          ]
        }
      ]
    },
    {
      "id": "D-B2",
      "base": "ea32f1f8c853243a523b0113b0b9d3453c8dd135",
      "tip": "6eed4ae5973ef9b9a9e9c77d8b6fcba5d3c698d2",
      "plan_digest": "sha256:46a171dbfa4f9e5b92760bbc064c364febc8bd9da34e0bd2e74698781183dd38",
      "source_digest": "b087942b45c60e7624a0737fd8018970569f1ce1",
      "acceptance_rows": [
        "D18"
      ],
      "verification": [
        {
          "id": "dcm-b2-t03-caller-tests",
          "performer": "claude:dcm_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6bc8eab39ace6f03dde72a6e14e14d96f675f046",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "649a980b5139c79657092f7987aec1aa6ea49edc",
            "digest": "sha256:d1dd65f2db545473c9aeb3bbd231cb858b10d5b625fad400077322cbb525de20",
            "excerpt": "# caller-tests at 649a980b5139c79657092f7987aec1aa6ea49edc: bench test --package ./internal/gate/prospectiveartifact\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,pass,274,37\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n# consumer package: bench test --package ./internal/gate -> pass, 450 tests\n"
          },
          "requirement": "caller-tests",
          "command": "bench test --package ./internal/gate/prospectiveartifact",
          "exit_code": 0
        },
        {
          "id": "dcm-b2-t03-root-conformance",
          "performer": "claude:dcm_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6bc8eab39ace6f03dde72a6e14e14d96f675f046",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "649a980b5139c79657092f7987aec1aa6ea49edc",
            "digest": "sha256:f3282407dde0fd892b57d3b38c50cd0d08cb31e33a98e7ba0ef7b163b54e4b3c",
            "excerpt": "# root-conformance at 649a980b5139c79657092f7987aec1aa6ea49edc: bench test --check root-conformance\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16673,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "dcm-b2-t03-omission-proof",
          "performer": "claude:dcm_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6bc8eab39ace6f03dde72a6e14e14d96f675f046",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "649a980b5139c79657092f7987aec1aa6ea49edc",
            "digest": "sha256:3e8be1747e3af634fa80264c9183c180687779b61b76c8dc6b6f7f422c2ab1f4",
            "excerpt": "# omission-proof at 649a980b5139c79657092f7987aec1aa6ea49edc: bench probe swap, bit and restored=yes\n## O1 owner bypass: prospectiveartifact.go swap 'return replacer.Replace(filepath.Join(root, RecordName), append(data, '\\n'), RecordMode)' with 'temporary, err := os.CreateTemp(root, \".owner-\"); if err != nil { return err }; name := temporary.Name(); defer func() { _ = temporary.Close(); _ = os.Remove(name) }(); if err := temporary.Chmod(RecordMode); err != nil { return err }; if _, err := temporary.Write(append(data, '\\n')); err != nil { return err }; if err := temporary.Sync(); err != nil { return err }; if err := temporary.Close(); err != nil { return err }; if err := os.Rename(name, filepath.Join(root, RecordName)); err != nil { return err }; directory, err := os.Open(root); if err != nil { return err }; defer directory.Close(); return directory.Sync()'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate/prospectiveartifact,all,passed,37\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,fail,280,37\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestOwnerReplacement,\"replacement_test.go:25: err = <nil>, want a published directory-sync *durablefile.Error that wraps injected owner directory-sync failure\",1\n## restored pass at 649a980b5139c79657092f7987aec1aa6ea49edc: bench test --package ./internal/gate/prospectiveartifact -> pass, 37 tests\n"
          },
          "requirement": "omission-proof",
          "command": "bench test --package ./internal/gate/prospectiveartifact",
          "exit_code": 0,
          "probe": {
            "mutation": "Restore the former local replacement sequence. The real owner-record junction test must detect the missing shared failure.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "649a980b5139c79657092f7987aec1aa6ea49edc",
              "digest": "sha256:3e8be1747e3af634fa80264c9183c180687779b61b76c8dc6b6f7f422c2ab1f4",
              "excerpt": "# omission-proof at 649a980b5139c79657092f7987aec1aa6ea49edc: bench probe swap, bit and restored=yes\n## O1 owner bypass: prospectiveartifact.go swap 'return replacer.Replace(filepath.Join(root, RecordName), append(data, '\\n'), RecordMode)' with 'temporary, err := os.CreateTemp(root, \".owner-\"); if err != nil { return err }; name := temporary.Name(); defer func() { _ = temporary.Close(); _ = os.Remove(name) }(); if err := temporary.Chmod(RecordMode); err != nil { return err }; if _, err := temporary.Write(append(data, '\\n')); err != nil { return err }; if err := temporary.Sync(); err != nil { return err }; if err := temporary.Close(); err != nil { return err }; if err := os.Rename(name, filepath.Join(root, RecordName)); err != nil { return err }; directory, err := os.Open(root); if err != nil { return err }; defer directory.Close(); return directory.Sync()'\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate/prospectiveartifact,all,passed,37\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,fail,280,37\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestOwnerReplacement,\"replacement_test.go:25: err = <nil>, want a published directory-sync *durablefile.Error that wraps injected owner directory-sync failure\",1\n## restored pass at 649a980b5139c79657092f7987aec1aa6ea49edc: bench test --package ./internal/gate/prospectiveartifact -> pass, 37 tests\n"
            }
          }
        },
        {
          "id": "dcm-b2-t03-preservation",
          "performer": "claude:dcm_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6bc8eab39ace6f03dde72a6e14e14d96f675f046",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "649a980b5139c79657092f7987aec1aa6ea49edc",
            "digest": "sha256:6224e068bc2b7a84abd96c62e3e64fd6532f3bbacc87acd4fee4dab9349900be",
            "excerpt": "# preservation at 649a980b5139c79657092f7987aec1aa6ea49edc: two bench probe swaps on the independent expectations of TestOwnerReplacement, each bit and restored=yes\n# Frozen baseline: TestOpenPublishesPrivateOwnerRecordBeforeCheckout (regular RecordMode file, decoded record), TestOpenRefusesInvalidTemporaryRootsWithoutPublishing (refused root publishes nothing), and the sweep tests that read records through Publish stay unchanged and pass (37 tests).\n## P1 encoded newline: prospectiveartifact.go swap \"append(data, '\\n'), RecordMode)\" with \"data, RecordMode)\" over the whole package; only the independent bytes literal catches it\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestOwnerReplacement,\"replacement_test.go:40: owner record = {\"schema\":1,\"owner_pid\":4242,\"common_dir\":\"/repo/.git\"}, want {\"schema\":1,\"owner_pid\":4242,\"common_dir\":\"/repo/.git\"}\\n\",1\n## P2 record mode: prospectiveartifact.go swap \"append(data, '\\n'), RecordMode)\" with \"append(data, '\\n'), 0o644)\" --run TestOwnerReplacement\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestOwnerReplacement,\"replacement_test.go:33: owner record mode = -rw-r--r--, want a regular 0600 file\",1\n## restored pass at 649a980b5139c79657092f7987aec1aa6ea49edc: bench test --package ./internal/gate/prospectiveartifact -> pass, 37 tests\n"
          },
          "requirement": "preservation",
          "command": "bench test --package ./internal/gate/prospectiveartifact",
          "exit_code": 0,
          "probe": {
            "mutation": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "649a980b5139c79657092f7987aec1aa6ea49edc",
              "digest": "sha256:6224e068bc2b7a84abd96c62e3e64fd6532f3bbacc87acd4fee4dab9349900be",
              "excerpt": "# preservation at 649a980b5139c79657092f7987aec1aa6ea49edc: two bench probe swaps on the independent expectations of TestOwnerReplacement, each bit and restored=yes\n# Frozen baseline: TestOpenPublishesPrivateOwnerRecordBeforeCheckout (regular RecordMode file, decoded record), TestOpenRefusesInvalidTemporaryRootsWithoutPublishing (refused root publishes nothing), and the sweep tests that read records through Publish stay unchanged and pass (37 tests).\n## P1 encoded newline: prospectiveartifact.go swap \"append(data, '\\n'), RecordMode)\" with \"data, RecordMode)\" over the whole package; only the independent bytes literal catches it\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestOwnerReplacement,\"replacement_test.go:40: owner record = {\"schema\":1,\"owner_pid\":4242,\"common_dir\":\"/repo/.git\"}, want {\"schema\":1,\"owner_pid\":4242,\"common_dir\":\"/repo/.git\"}\\n\",1\n## P2 record mode: prospectiveartifact.go swap \"append(data, '\\n'), RecordMode)\" with \"append(data, '\\n'), 0o644)\" --run TestOwnerReplacement\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestOwnerReplacement,\"replacement_test.go:33: owner record mode = -rw-r--r--, want a regular 0600 file\",1\n## restored pass at 649a980b5139c79657092f7987aec1aa6ea49edc: bench test --package ./internal/gate/prospectiveartifact -> pass, 37 tests\n"
            }
          }
        },
        {
          "id": "dcm-b2-t03-rep1-caller-tests",
          "performer": "claude:dcm_t03_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "123e6c7765ecd010a41bd5ebfadaf1c2ee6e5a2a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "3a3dbbe1a471a4fdd48e24abe38359ff138a1e73",
            "digest": "sha256:5e1b6dadbf14d96d3c017a341124061f8f50e167f548541f4910bf98483a27b2",
            "excerpt": "# caller-tests at 3a3dbbe1a471a4fdd48e24abe38359ff138a1e73: bench test --package ./internal/gate/prospectiveartifact --package ./internal/durablefile/durablefiletest --package ./internal/gate\npackages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile/durablefiletest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/gate,pass,25643,450\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,pass,329,38\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n# census: bench test --package ./internal/worktree --run 'TestParallelCensusOnTheLiveTree|TestSerialSetStaysBelowTheCeiling' -> pass, 2 tests\n"
          },
          "requirement": "caller-tests",
          "command": "bench test --package ./internal/gate/prospectiveartifact",
          "exit_code": 0
        },
        {
          "id": "dcm-b2-t03-rep1-durablefiletest",
          "performer": "claude:dcm_t03_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "123e6c7765ecd010a41bd5ebfadaf1c2ee6e5a2a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "3a3dbbe1a471a4fdd48e24abe38359ff138a1e73",
            "digest": "sha256:1aa789754fb58d492dfb98a4b7caac73eee38f7861a4061d7afefbd9c3806c00",
            "excerpt": "# durablefiletest at 3a3dbbe1a471a4fdd48e24abe38359ff138a1e73: bench test --package ./internal/durablefile/durablefiletest (run in the one focused call)\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile/durablefiletest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\n# The package compiles and has no own tests; its RequireFailure and RequireFile assertions run in the gate and prospectiveartifact caller tests, and the probes in preservation and omission-proof show each bite.\n"
          },
          "requirement": "durablefiletest",
          "command": "bench test --package ./internal/durablefile/durablefiletest",
          "exit_code": 0
        },
        {
          "id": "dcm-b2-t03-rep1-gate-tests",
          "performer": "claude:dcm_t03_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "123e6c7765ecd010a41bd5ebfadaf1c2ee6e5a2a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "3a3dbbe1a471a4fdd48e24abe38359ff138a1e73",
            "digest": "sha256:6021053aa9ad208fadc7c8c29808dba860be9fd66667d165f5d2a958a9d3cd2b",
            "excerpt": "# gate-tests at 3a3dbbe1a471a4fdd48e24abe38359ff138a1e73: bench test --package ./internal/gate (run in the one focused call)\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,25643,450\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "gate-tests",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "dcm-b2-t03-rep1-root-conformance",
          "performer": "claude:dcm_t03_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "123e6c7765ecd010a41bd5ebfadaf1c2ee6e5a2a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "3a3dbbe1a471a4fdd48e24abe38359ff138a1e73",
            "digest": "sha256:cde52bcc3e10c6fd3c63e867cc139bd62ed55973613a322f242589727ce09a84",
            "excerpt": "# root-conformance at 3a3dbbe1a471a4fdd48e24abe38359ff138a1e73: bench test --check root-conformance\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15484,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "dcm-b2-t03-rep1-omission-proof",
          "performer": "claude:dcm_t03_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "123e6c7765ecd010a41bd5ebfadaf1c2ee6e5a2a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "3a3dbbe1a471a4fdd48e24abe38359ff138a1e73",
            "digest": "sha256:50cb030b957cf41bf740ed5cc9627300da0653a3db593bc567c710f5bb61d056",
            "excerpt": "# omission-proof at 3a3dbbe1a471a4fdd48e24abe38359ff138a1e73: two bench probe swaps over the whole ./internal/gate/prospectiveartifact package (no --run), each bit and restored=yes\n## O1 owner bypass: prospectiveartifact.go swap \"return replacer.Replace(filepath.Join(root, RecordName), append(data, '\\n'), RecordMode)\" with \"temporary, err := os.CreateTemp(root, \\\".owner-\\\"); if err != nil { return err }; name := temporary.Name(); defer func() { _ = temporary.Close(); _ = os.Remove(name) }(); if err := temporary.Chmod(RecordMode); err != nil { return err }; if _, err := temporary.Write(append(data, '\\n')); err != nil { return err }; if err := temporary.Sync(); err != nil { return err }; if err := temporary.Close(); err != nil { return err }; if err := os.Rename(name, filepath.Join(root, RecordName)); err != nil { return err }; directory, err := os.Open(root); if err != nil { return err }; defer directory.Close(); return directory.Sync()\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact.go,swap,failed,2,yes\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,fail,245,38\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestOwnerPreRenameFailure,\"replacement_test.go:37: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected owner rename failure\",1\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestOwnerReplacement,\"replacement_test.go:23: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected owner directory-sync failure\",1\n## C1 pre-rename bypass (Phase 1 loop at db1f156a: silent, 0 failed of 37; at 3a3dbbe1: bit): prospectiveartifact.go swap \"return replacer.Replace(filepath.Join(root, RecordName), append(data, '\\n'), RecordMode)\" with \"if err := replacer.Replace(filepath.Join(root, RecordName), append(data, '\\n'), RecordMode); err != nil { var replaceErr *durablefile.Error; if errors.As(err, &replaceErr) && !replaceErr.Published { return nil }; return err }; return nil\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact.go,swap,failed,1,yes\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,fail,246,38\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestOwnerPreRenameFailure,\"replacement_test.go:37: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected owner rename failure\",1\n## restored pass at 3a3dbbe1a471a4fdd48e24abe38359ff138a1e73: bench test --package ./internal/gate/prospectiveartifact -> pass, 38 tests\n"
          },
          "requirement": "omission-proof",
          "command": "bench test --package ./internal/gate/prospectiveartifact",
          "exit_code": 0,
          "probe": {
            "mutation": "Restore the former local replacement sequence. The real owner-record junction test must detect the missing shared failure.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "3a3dbbe1a471a4fdd48e24abe38359ff138a1e73",
              "digest": "sha256:50cb030b957cf41bf740ed5cc9627300da0653a3db593bc567c710f5bb61d056",
              "excerpt": "# omission-proof at 3a3dbbe1a471a4fdd48e24abe38359ff138a1e73: two bench probe swaps over the whole ./internal/gate/prospectiveartifact package (no --run), each bit and restored=yes\n## O1 owner bypass: prospectiveartifact.go swap \"return replacer.Replace(filepath.Join(root, RecordName), append(data, '\\n'), RecordMode)\" with \"temporary, err := os.CreateTemp(root, \\\".owner-\\\"); if err != nil { return err }; name := temporary.Name(); defer func() { _ = temporary.Close(); _ = os.Remove(name) }(); if err := temporary.Chmod(RecordMode); err != nil { return err }; if _, err := temporary.Write(append(data, '\\n')); err != nil { return err }; if err := temporary.Sync(); err != nil { return err }; if err := temporary.Close(); err != nil { return err }; if err := os.Rename(name, filepath.Join(root, RecordName)); err != nil { return err }; directory, err := os.Open(root); if err != nil { return err }; defer directory.Close(); return directory.Sync()\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact.go,swap,failed,2,yes\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,fail,245,38\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestOwnerPreRenameFailure,\"replacement_test.go:37: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected owner rename failure\",1\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestOwnerReplacement,\"replacement_test.go:23: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected owner directory-sync failure\",1\n## C1 pre-rename bypass (Phase 1 loop at db1f156a: silent, 0 failed of 37; at 3a3dbbe1: bit): prospectiveartifact.go swap \"return replacer.Replace(filepath.Join(root, RecordName), append(data, '\\n'), RecordMode)\" with \"if err := replacer.Replace(filepath.Join(root, RecordName), append(data, '\\n'), RecordMode); err != nil { var replaceErr *durablefile.Error; if errors.As(err, &replaceErr) && !replaceErr.Published { return nil }; return err }; return nil\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact.go,swap,failed,1,yes\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,fail,246,38\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestOwnerPreRenameFailure,\"replacement_test.go:37: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected owner rename failure\",1\n## restored pass at 3a3dbbe1a471a4fdd48e24abe38359ff138a1e73: bench test --package ./internal/gate/prospectiveartifact -> pass, 38 tests\n"
            }
          }
        },
        {
          "id": "dcm-b2-t03-rep1-preservation",
          "performer": "claude:dcm_t03_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "123e6c7765ecd010a41bd5ebfadaf1c2ee6e5a2a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "3a3dbbe1a471a4fdd48e24abe38359ff138a1e73",
            "digest": "sha256:62776d53bd89c7285fe2d2b35a23411457a90a2d153f5477a05490a19c19f248",
            "excerpt": "# preservation at 3a3dbbe1a471a4fdd48e24abe38359ff138a1e73: six bench probe swaps, each bit and restored=yes; the owner rows run over the whole package with no --run\n# Frozen baseline: TestOpenPublishesPrivateOwnerRecordBeforeCheckout, TestOpenRefusesInvalidTemporaryRootsWithoutPublishing, and the sweep tests stay unchanged and pass (38 tests with the new TestOwnerPreRenameFailure).\n## P1 owner encoded newline: prospectiveartifact.go swap \"append(data, '\\n'), RecordMode)\" with \"data, RecordMode)\" --package ./internal/gate/prospectiveartifact\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestOwnerReplacement,\"replacement_test.go:24: .../owner.json = {\"schema\":1,\"owner_pid\":4242,\"common_dir\":\"/repo/.git\"}, want {\"schema\":1,\"owner_pid\":4242,\"common_dir\":\"/repo/.git\"}\\n\",1\n## P2 owner 0600 literal (C3 necessity red): prospectiveartifact.go swap \"const RecordMode os.FileMode = 0o600\" with \"const RecordMode os.FileMode = 0o640\" --package ./internal/gate/prospectiveartifact\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact.go,swap,failed,1,yes\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,fail,246,38\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestOwnerReplacement,\"replacement_test.go:24: .../owner.json mode = -rw-r-----, want a regular file at -rw-------\",1\n## D16/D17 reds after the move to durablefiletest.RequireFailure and RequireFile (--package ./internal/gate --run 'TestVerdictReplacement|TestEvidenceReplacement')\n## G1 mode: engine.go swap \"append(data, '\\n'), 0o600)\" with \"append(data, '\\n'), 0o644)\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/engine.go,swap,failed,3,yes\nfailures[3]: TestEvidenceReplacement/evidence, TestEvidenceReplacement/lane, TestVerdictReplacement: \"durable_replacement_test.go:67|30: .../<entry> mode = -rw-r--r--, want a regular file at -rw-------\"\n# The same swap over the whole ./internal/gate package (no --run): bit, failed 22, restored=yes, TestVerdictReplacement among them.\n## G2 newline: engine.go swap \"append(data, '\\n'), 0o600)\" with \"data, 0o600)\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/engine.go,swap,failed,3,yes\nfailures[3]: TestEvidenceReplacement/evidence, TestEvidenceReplacement/lane, TestVerdictReplacement: \"durable_replacement_test.go:67|30: .../<entry> = <record without newline>, want <record>\\n\"\n## G3 store-writer bypass: engine.go swap \"return replacer.Replace(filepath.Join(dir, name), append(data, '\\n'), 0o600)\" with \"return durablefile.Replace(filepath.Join(dir, name), append(data, '\\n'), 0o600)\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/engine.go,swap,failed,3,yes\nfailures[3]: TestEvidenceReplacement/evidence, TestEvidenceReplacement/lane (line 66), TestVerdictReplacement (line 29): \"err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected <name> directory-sync failure\"\n## G4 verdict-writer bypass: verdict.go swap \"return durableReplaceAt(replacer, gitdir, benchgit.GateCacheFile, rec)\" with \"return durableReplaceAt(durablefile.Native(), gitdir, benchgit.GateCacheFile, rec)\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/verdict.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestVerdictReplacement,\"durable_replacement_test.go:29: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected verdict directory-sync failure\",1\n## restored pass at 3a3dbbe1a471a4fdd48e24abe38359ff138a1e73: bench test --package ./internal/gate/prospectiveartifact --package ./internal/durablefile/durablefiletest --package ./internal/gate -> pass, 38 + 450 tests\n"
          },
          "requirement": "preservation",
          "command": "bench test --package ./internal/gate/prospectiveartifact",
          "exit_code": 0,
          "probe": {
            "mutation": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "3a3dbbe1a471a4fdd48e24abe38359ff138a1e73",
              "digest": "sha256:62776d53bd89c7285fe2d2b35a23411457a90a2d153f5477a05490a19c19f248",
              "excerpt": "# preservation at 3a3dbbe1a471a4fdd48e24abe38359ff138a1e73: six bench probe swaps, each bit and restored=yes; the owner rows run over the whole package with no --run\n# Frozen baseline: TestOpenPublishesPrivateOwnerRecordBeforeCheckout, TestOpenRefusesInvalidTemporaryRootsWithoutPublishing, and the sweep tests stay unchanged and pass (38 tests with the new TestOwnerPreRenameFailure).\n## P1 owner encoded newline: prospectiveartifact.go swap \"append(data, '\\n'), RecordMode)\" with \"data, RecordMode)\" --package ./internal/gate/prospectiveartifact\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestOwnerReplacement,\"replacement_test.go:24: .../owner.json = {\"schema\":1,\"owner_pid\":4242,\"common_dir\":\"/repo/.git\"}, want {\"schema\":1,\"owner_pid\":4242,\"common_dir\":\"/repo/.git\"}\\n\",1\n## P2 owner 0600 literal (C3 necessity red): prospectiveartifact.go swap \"const RecordMode os.FileMode = 0o600\" with \"const RecordMode os.FileMode = 0o640\" --package ./internal/gate/prospectiveartifact\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact.go,swap,failed,1,yes\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,fail,246,38\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestOwnerReplacement,\"replacement_test.go:24: .../owner.json mode = -rw-r-----, want a regular file at -rw-------\",1\n## D16/D17 reds after the move to durablefiletest.RequireFailure and RequireFile (--package ./internal/gate --run 'TestVerdictReplacement|TestEvidenceReplacement')\n## G1 mode: engine.go swap \"append(data, '\\n'), 0o600)\" with \"append(data, '\\n'), 0o644)\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/engine.go,swap,failed,3,yes\nfailures[3]: TestEvidenceReplacement/evidence, TestEvidenceReplacement/lane, TestVerdictReplacement: \"durable_replacement_test.go:67|30: .../<entry> mode = -rw-r--r--, want a regular file at -rw-------\"\n# The same swap over the whole ./internal/gate package (no --run): bit, failed 22, restored=yes, TestVerdictReplacement among them.\n## G2 newline: engine.go swap \"append(data, '\\n'), 0o600)\" with \"data, 0o600)\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/engine.go,swap,failed,3,yes\nfailures[3]: TestEvidenceReplacement/evidence, TestEvidenceReplacement/lane, TestVerdictReplacement: \"durable_replacement_test.go:67|30: .../<entry> = <record without newline>, want <record>\\n\"\n## G3 store-writer bypass: engine.go swap \"return replacer.Replace(filepath.Join(dir, name), append(data, '\\n'), 0o600)\" with \"return durablefile.Replace(filepath.Join(dir, name), append(data, '\\n'), 0o600)\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/engine.go,swap,failed,3,yes\nfailures[3]: TestEvidenceReplacement/evidence, TestEvidenceReplacement/lane (line 66), TestVerdictReplacement (line 29): \"err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected <name> directory-sync failure\"\n## G4 verdict-writer bypass: verdict.go swap \"return durableReplaceAt(replacer, gitdir, benchgit.GateCacheFile, rec)\" with \"return durableReplaceAt(durablefile.Native(), gitdir, benchgit.GateCacheFile, rec)\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/verdict.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestVerdictReplacement,\"durable_replacement_test.go:29: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected verdict directory-sync failure\",1\n## restored pass at 3a3dbbe1a471a4fdd48e24abe38359ff138a1e73: bench test --package ./internal/gate/prospectiveartifact --package ./internal/durablefile/durablefiletest --package ./internal/gate -> pass, 38 + 450 tests\n"
            }
          }
        },
        {
          "id": "dcm-b2-t03-rep2-caller-tests",
          "performer": "claude:dcm_t03_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "b087942b45c60e7624a0737fd8018970569f1ce1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "6eed4ae5973ef9b9a9e9c77d8b6fcba5d3c698d2",
            "digest": "sha256:5d3b18c666359eac9cdea4570ce8a9e94979bf917556096657b65244cfd56e3d",
            "excerpt": "# caller-tests at 6eed4ae5973ef9b9a9e9c77d8b6fcba5d3c698d2: bench test --package ./internal/gate/prospectiveartifact (run in the one focused call with ./internal/durablefile/durablefiletest and ./internal/gate)\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,pass,276,38\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "caller-tests",
          "command": "bench test --package ./internal/gate/prospectiveartifact",
          "exit_code": 0
        },
        {
          "id": "dcm-b2-t03-rep2-durablefiletest",
          "performer": "claude:dcm_t03_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "b087942b45c60e7624a0737fd8018970569f1ce1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "6eed4ae5973ef9b9a9e9c77d8b6fcba5d3c698d2",
            "digest": "sha256:5d4f6eb0d36f8b6465ed1d6cacebefc06407f0fadd117200cdeaf49f031876c8",
            "excerpt": "# durablefiletest at 6eed4ae5973ef9b9a9e9c77d8b6fcba5d3c698d2: bench test --package ./internal/durablefile/durablefiletest (run in the one focused call)\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile/durablefiletest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\n"
          },
          "requirement": "durablefiletest",
          "command": "bench test --package ./internal/durablefile/durablefiletest",
          "exit_code": 0
        },
        {
          "id": "dcm-b2-t03-rep2-gate-tests",
          "performer": "claude:dcm_t03_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "b087942b45c60e7624a0737fd8018970569f1ce1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "6eed4ae5973ef9b9a9e9c77d8b6fcba5d3c698d2",
            "digest": "sha256:0dbd3588dc6a253e524d2e8a45131bf0fbe9a5e82852c57dddd7ac6af4d52b8c",
            "excerpt": "# gate-tests at 6eed4ae5973ef9b9a9e9c77d8b6fcba5d3c698d2: bench test --package ./internal/gate (run in the one focused call)\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,25540,450\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "gate-tests",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "dcm-b2-t03-rep2-root-conformance",
          "performer": "claude:dcm_t03_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "b087942b45c60e7624a0737fd8018970569f1ce1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "6eed4ae5973ef9b9a9e9c77d8b6fcba5d3c698d2",
            "digest": "sha256:d81a95887e0d98d6014e279ea88804ce0b82a2b1c510990f5eb1e96af163ea8e",
            "excerpt": "# root-conformance at 6eed4ae5973ef9b9a9e9c77d8b6fcba5d3c698d2: bench test --check root-conformance\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16267,1\nfailures[0]{package,test,line,lines}:\n# also green: bench test --package ./internal/worktree --run 'TestParallelCensusOnTheLiveTree|TestSerialSetStaysBelowTheCeiling' -> pass, 2 tests\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "dcm-b2-t03-rep2-omission-proof",
          "performer": "claude:dcm_t03_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "b087942b45c60e7624a0737fd8018970569f1ce1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "6eed4ae5973ef9b9a9e9c77d8b6fcba5d3c698d2",
            "digest": "sha256:7f29e2c4ed976af641bb1332c56fb2b5bf0f07e55d9f5fc19f957c72e0763a20",
            "excerpt": "# omission-proof at 6eed4ae5973ef9b9a9e9c77d8b6fcba5d3c698d2: three bench probe swaps over the whole ./internal/gate/prospectiveartifact package (no --run), each bit and restored=yes\n## O1 owner bypass (former local sequence): prospectiveartifact.go swap \"return replacer.Replace(filepath.Join(root, RecordName), append(data, '\\n'), RecordMode)\" with \"temporary, err := os.CreateTemp(root, \\\".owner-\\\"); if err != nil { return err }; name := temporary.Name(); defer func() { _ = temporary.Close(); _ = os.Remove(name) }(); if err := temporary.Chmod(RecordMode); err != nil { return err }; if _, err := temporary.Write(append(data, '\\n')); err != nil { return err }; if err := temporary.Sync(); err != nil { return err }; if err := temporary.Close(); err != nil { return err }; if err := os.Rename(name, filepath.Join(root, RecordName)); err != nil { return err }; directory, err := os.Open(root); if err != nil { return err }; defer directory.Close(); return directory.Sync()\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact.go,swap,failed,2,yes\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestOwnerPreRenameFailure,\"replacement_test.go:37: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected owner rename failure\",1\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestOwnerReplacement,\"replacement_test.go:23: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected owner directory-sync failure\",1\n## C1 pre-rename bypass: prospectiveartifact.go swap \"return replacer.Replace(filepath.Join(root, RecordName), append(data, '\\n'), RecordMode)\" with \"if err := replacer.Replace(filepath.Join(root, RecordName), append(data, '\\n'), RecordMode); err != nil { var replaceErr *durablefile.Error; if errors.As(err, &replaceErr) && !replaceErr.Published { return nil }; return err }; return nil\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestOwnerPreRenameFailure,\"replacement_test.go:37: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected owner rename failure\",1\n## R1 record written after a pre-rename failure (necessity red of the kept record-absence check): prospectiveartifact.go swap \"return replacer.Replace(filepath.Join(root, RecordName), append(data, '\\n'), RecordMode)\" with \"if err := replacer.Replace(filepath.Join(root, RecordName), append(data, '\\n'), RecordMode); err != nil { var replaceErr *durablefile.Error; if errors.As(err, &replaceErr) && !replaceErr.Published { _ = os.WriteFile(filepath.Join(root, RecordName), append(data, '\\n'), RecordMode) }; return err }; return nil\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestOwnerPreRenameFailure,\"replacement_test.go:39: owner record stat = <nil>, want no published record\",1\n## restored pass at 6eed4ae5973ef9b9a9e9c77d8b6fcba5d3c698d2: bench test --package ./internal/gate/prospectiveartifact -> pass, 38 tests\n"
          },
          "requirement": "omission-proof",
          "command": "bench test --package ./internal/gate/prospectiveartifact",
          "exit_code": 0,
          "probe": {
            "mutation": "Restore the former local replacement sequence. The real owner-record junction test must detect the missing shared failure.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "6eed4ae5973ef9b9a9e9c77d8b6fcba5d3c698d2",
              "digest": "sha256:7f29e2c4ed976af641bb1332c56fb2b5bf0f07e55d9f5fc19f957c72e0763a20",
              "excerpt": "# omission-proof at 6eed4ae5973ef9b9a9e9c77d8b6fcba5d3c698d2: three bench probe swaps over the whole ./internal/gate/prospectiveartifact package (no --run), each bit and restored=yes\n## O1 owner bypass (former local sequence): prospectiveartifact.go swap \"return replacer.Replace(filepath.Join(root, RecordName), append(data, '\\n'), RecordMode)\" with \"temporary, err := os.CreateTemp(root, \\\".owner-\\\"); if err != nil { return err }; name := temporary.Name(); defer func() { _ = temporary.Close(); _ = os.Remove(name) }(); if err := temporary.Chmod(RecordMode); err != nil { return err }; if _, err := temporary.Write(append(data, '\\n')); err != nil { return err }; if err := temporary.Sync(); err != nil { return err }; if err := temporary.Close(); err != nil { return err }; if err := os.Rename(name, filepath.Join(root, RecordName)); err != nil { return err }; directory, err := os.Open(root); if err != nil { return err }; defer directory.Close(); return directory.Sync()\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact.go,swap,failed,2,yes\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestOwnerPreRenameFailure,\"replacement_test.go:37: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected owner rename failure\",1\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestOwnerReplacement,\"replacement_test.go:23: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected owner directory-sync failure\",1\n## C1 pre-rename bypass: prospectiveartifact.go swap \"return replacer.Replace(filepath.Join(root, RecordName), append(data, '\\n'), RecordMode)\" with \"if err := replacer.Replace(filepath.Join(root, RecordName), append(data, '\\n'), RecordMode); err != nil { var replaceErr *durablefile.Error; if errors.As(err, &replaceErr) && !replaceErr.Published { return nil }; return err }; return nil\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestOwnerPreRenameFailure,\"replacement_test.go:37: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected owner rename failure\",1\n## R1 record written after a pre-rename failure (necessity red of the kept record-absence check): prospectiveartifact.go swap \"return replacer.Replace(filepath.Join(root, RecordName), append(data, '\\n'), RecordMode)\" with \"if err := replacer.Replace(filepath.Join(root, RecordName), append(data, '\\n'), RecordMode); err != nil { var replaceErr *durablefile.Error; if errors.As(err, &replaceErr) && !replaceErr.Published { _ = os.WriteFile(filepath.Join(root, RecordName), append(data, '\\n'), RecordMode) }; return err }; return nil\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestOwnerPreRenameFailure,\"replacement_test.go:39: owner record stat = <nil>, want no published record\",1\n## restored pass at 6eed4ae5973ef9b9a9e9c77d8b6fcba5d3c698d2: bench test --package ./internal/gate/prospectiveartifact -> pass, 38 tests\n"
            }
          }
        },
        {
          "id": "dcm-b2-t03-rep2-preservation",
          "performer": "claude:dcm_t03_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "b087942b45c60e7624a0737fd8018970569f1ce1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "6eed4ae5973ef9b9a9e9c77d8b6fcba5d3c698d2",
            "digest": "sha256:090d53fadfe2a882af7d18f6aa2404cbf1f5a211106f5ef13ca7f20a42dc0ab0",
            "excerpt": "# preservation at 6eed4ae5973ef9b9a9e9c77d8b6fcba5d3c698d2: three bench probe swaps over the whole ./internal/gate/prospectiveartifact package (no --run); P1 and P2 bit, L1 silent as intended; each restored=yes\n# Frozen baseline: TestOpenPublishesPrivateOwnerRecordBeforeCheckout, TestOpenRefusesInvalidTemporaryRootsWithoutPublishing, and the sweep tests stay unchanged and pass (38 tests).\n## P1 owner encoded newline: prospectiveartifact.go swap \"append(data, '\\n'), RecordMode)\" with \"data, RecordMode)\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact.go,swap,failed,1,yes\nfailures[1]: TestOwnerReplacement \"replacement_test.go:24: .../owner.json = <record without newline>, want <record>\\n\"\n## P2 owner 0600 literal: prospectiveartifact.go swap \"const RecordMode os.FileMode = 0o600\" with \"const RecordMode os.FileMode = 0o640\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestOwnerReplacement,\"replacement_test.go:24: .../owner.json mode = -rw-r-----, want a regular file at -rw-------\",1\n## L1 leaf cleanup ownership (DCM-B2-S2): internal/durablefile/replace.go swap \"Err: errors.Join(err, r.Remove(temp.Name()))}\" with \"Err: err}\"\n# at 40aef905 over ./internal/gate/prospectiveartifact: bit, failed 1, restored=yes: TestOwnerPreRenameFailure \"replacement_test.go:43: root entries = [.owner.json.405312725], want no temporary file\"\n# at 40aef905 over ./internal/durablefile: bit, failed 2, restored=yes: TestReplaceTemporaryCleanup (replace_test.go:263), TestReplaceCleanupFailure (replace_test.go:275)\n# at 6eed4ae5 over ./internal/gate/prospectiveartifact: silent, failed 0, restored=yes; the leaf tests alone own temporary-file cleanup\n## restored pass at 6eed4ae5973ef9b9a9e9c77d8b6fcba5d3c698d2: bench test --package ./internal/gate/prospectiveartifact --package ./internal/durablefile/durablefiletest --package ./internal/gate -> pass, 38 + 450 tests\n"
          },
          "requirement": "preservation",
          "command": "bench test --package ./internal/gate/prospectiveartifact",
          "exit_code": 0,
          "probe": {
            "mutation": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "6eed4ae5973ef9b9a9e9c77d8b6fcba5d3c698d2",
              "digest": "sha256:090d53fadfe2a882af7d18f6aa2404cbf1f5a211106f5ef13ca7f20a42dc0ab0",
              "excerpt": "# preservation at 6eed4ae5973ef9b9a9e9c77d8b6fcba5d3c698d2: three bench probe swaps over the whole ./internal/gate/prospectiveartifact package (no --run); P1 and P2 bit, L1 silent as intended; each restored=yes\n# Frozen baseline: TestOpenPublishesPrivateOwnerRecordBeforeCheckout, TestOpenRefusesInvalidTemporaryRootsWithoutPublishing, and the sweep tests stay unchanged and pass (38 tests).\n## P1 owner encoded newline: prospectiveartifact.go swap \"append(data, '\\n'), RecordMode)\" with \"data, RecordMode)\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact.go,swap,failed,1,yes\nfailures[1]: TestOwnerReplacement \"replacement_test.go:24: .../owner.json = <record without newline>, want <record>\\n\"\n## P2 owner 0600 literal: prospectiveartifact.go swap \"const RecordMode os.FileMode = 0o600\" with \"const RecordMode os.FileMode = 0o640\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/prospectiveartifact/prospectiveartifact.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,TestOwnerReplacement,\"replacement_test.go:24: .../owner.json mode = -rw-r-----, want a regular file at -rw-------\",1\n## L1 leaf cleanup ownership (DCM-B2-S2): internal/durablefile/replace.go swap \"Err: errors.Join(err, r.Remove(temp.Name()))}\" with \"Err: err}\"\n# at 40aef905 over ./internal/gate/prospectiveartifact: bit, failed 1, restored=yes: TestOwnerPreRenameFailure \"replacement_test.go:43: root entries = [.owner.json.405312725], want no temporary file\"\n# at 40aef905 over ./internal/durablefile: bit, failed 2, restored=yes: TestReplaceTemporaryCleanup (replace_test.go:263), TestReplaceCleanupFailure (replace_test.go:275)\n# at 6eed4ae5 over ./internal/gate/prospectiveartifact: silent, failed 0, restored=yes; the leaf tests alone own temporary-file cleanup\n## restored pass at 6eed4ae5973ef9b9a9e9c77d8b6fcba5d3c698d2: bench test --package ./internal/gate/prospectiveartifact --package ./internal/durablefile/durablefiletest --package ./internal/gate -> pass, 38 + 450 tests\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "dcm-b2-r1-standards",
          "performer": "claude:dcm_rev_b2_std",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "6bc8eab39ace6f03dde72a6e14e14d96f675f046",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b2_std",
            "digest": "sha256:f883e1d98be0f7a2519e306c8e55fa3fbd922f962dc3d216588504e6737e4b77",
            "excerpt": "## Standards (D-B2, ea32f1f8..6415ae0f) — 1 finding; worst DCM-B2-S1\nDCM-B2-S1 | ask-user (orchestrator accepted as auto-fix under --auto-approve) | 6 | internal/gate/prospectiveartifact/replacement_test.go:22-41 | copies the published-failure predicate and the record-file check of internal/gate/durable_replacement_test.go:292-312 (requirePublishedFailure, requireRecordFile); ten more caller packages will need them | move one exported assertion into internal/durablefile/durablefiletest and call it from both packages\nChecks: no copy of the durablefile algorithm; the seam follows the gate idiom; error handling unchanged; independent literals have recorded reds.\nAdvice: the temp name change to .owner.json.* is acceptable.\n"
          },
          "axis": "Standards",
          "base": "ea32f1f8c853243a523b0113b0b9d3453c8dd135",
          "tip": "649a980b5139c79657092f7987aec1aa6ea49edc",
          "finding_ids": [
            "DCM-B2-S1"
          ],
          "supersedes": []
        },
        {
          "id": "dcm-b2-r1-spec",
          "performer": "claude:dcm_rev_b2_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "6bc8eab39ace6f03dde72a6e14e14d96f675f046",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b2_spec",
            "digest": "sha256:1bb970bef690f0b465800b9300224d9c440405e78fee60c30dcd36df30588476",
            "excerpt": "## Spec (D-B2, ea32f1f8..6415ae0f) — 0 findings\nD18 holds: RecordMode and Marshal plus final newline at prospectiveartifact.go:236-240; exported Publish keeps its signature and calls publish(durablefile.Native(), root, record); caller authority stays in Factory.Open; read-back observes the exact new record. D30 claim holds: Publish had no pre-write refusal beyond an unreachable marshal error.\nAdvice: seed an older record before the fault so read-back proves replacement.\n"
          },
          "axis": "Spec",
          "base": "ea32f1f8c853243a523b0113b0b9d3453c8dd135",
          "tip": "649a980b5139c79657092f7987aec1aa6ea49edc",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "dcm-b2-r1-coverage",
          "performer": "claude:dcm_rev_b2_cov",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "6bc8eab39ace6f03dde72a6e14e14d96f675f046",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b2_cov",
            "digest": "sha256:5e6651cd65f038c37b49800c950ff7ad9b8d81111bd6b8a3661c39256800127b",
            "excerpt": "## Coverage (D-B2, ea32f1f8..6415ae0f) — 3 findings; worst C1. Five probes, each restored=yes.\nC1 | auto-fix | 8 | internal/gate/prospectiveartifact/replacement_test.go:19-25; prospectiveartifact.go:240 | only the post-rename branch is tested; a publish that returns nil on every unpublished *durablefile.Error is silent (0 of 505) | add a pre-rename row (for example a rename fault) asserting Published=false, errors.Is the cause, owner.json absent, and no siblings; record its bypass red\nC2 | no-op | 6 | prospectiveartifact.go:231 | no test detects a plain os.WriteFile in exported Publish; the seam accepts this limit | the D-D source census checks the native binding\nC3 | auto-fix | 7 | reviews/durable-caller-migration.md (dcm-b2-t03-preservation P2) | the P2 mode red ran with --run TestOwnerReplacement, so it does not show the 0600 literal is necessary; const RecordMode 0o600 -> 0o640 over the whole package bites only TestOwnerReplacement:33 | record the whole-package constant-drift red\nAdvice: the marshal refusal is unreachable for Record's fields.\n"
          },
          "axis": "Coverage",
          "base": "ea32f1f8c853243a523b0113b0b9d3453c8dd135",
          "tip": "649a980b5139c79657092f7987aec1aa6ea49edc",
          "finding_ids": [
            "C1",
            "C2",
            "C3"
          ],
          "supersedes": []
        },
        {
          "id": "dcm-b2-r2-standards",
          "performer": "claude:dcm_rev_b2_std2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "123e6c7765ecd010a41bd5ebfadaf1c2ee6e5a2a",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b2_std2",
            "digest": "sha256:3afbef56bfa0d76d684cad4773b4e63ea57550fbe0c9fb4bfc1dd9b87a732d1b",
            "excerpt": "## Standards confirming round (D-B2, repair delta db1f156a..f0c35260) — 1 finding; DCM-B2-S1 closed\nDCM-B2-S1 closed: RequireFailure and RequireFile in internal/durablefile/durablefiletest/assert.go:18-46; both caller tests use them; no local copy remains.\nDCM-B2-S2 | auto-fix | 5 | internal/gate/prospectiveartifact/replacement_test.go:42-44 | TestOwnerPreRenameFailure asserts durablefiletest.Siblings is empty, a leaf postcondition that TestReplaceTemporaryCleanup (internal/durablefile/replace_test.go:262) owns; no caller mutation shows the check is necessary | remove the sibling check; keep the record-absence check\nAdvice: the published bool reads unclearly; a RequireAbsent helper when a second caller needs it; the durablefiletest package comment names only the harness; ticket 03 says one assertion where two exist.\n"
          },
          "axis": "Standards",
          "base": "ea32f1f8c853243a523b0113b0b9d3453c8dd135",
          "tip": "3a3dbbe1a471a4fdd48e24abe38359ff138a1e73",
          "finding_ids": [
            "DCM-B2-S2"
          ],
          "supersedes": [
            "dcm-b2-r1-standards"
          ]
        },
        {
          "id": "dcm-b2-r2-spec",
          "performer": "claude:dcm_rev_b2_spec2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "123e6c7765ecd010a41bd5ebfadaf1c2ee6e5a2a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b2_spec2",
            "digest": "sha256:840f617b85efbbf813954d8fddcbd39aa221547aeaa6c6e26ec4762e51ab4034",
            "excerpt": "## Spec confirming round (D-B2, repair delta db1f156a..f0c35260 and plan b511def3..db1f156a) — 0 findings\nD16, D17, D18 hold after the move to the shared helpers. TestOwnerPreRenameFailure meets the expanded ticket 03 text. The plan expansion adds no acceptance row and weakens no guarantee.\nAdvice: ticket 03 wording says one assertion; ticket 14 rationale sentence is universal; the pre-rename test uses an empty root.\n"
          },
          "axis": "Spec",
          "base": "ea32f1f8c853243a523b0113b0b9d3453c8dd135",
          "tip": "3a3dbbe1a471a4fdd48e24abe38359ff138a1e73",
          "finding_ids": [],
          "supersedes": [
            "dcm-b2-r1-spec"
          ]
        },
        {
          "id": "dcm-b2-r2-coverage",
          "performer": "claude:dcm_rev_b2_cov2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "123e6c7765ecd010a41bd5ebfadaf1c2ee6e5a2a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b2_cov2",
            "digest": "sha256:6c624de091ab9e410a0df1b68b756f31ac7f52463753330f3b91df0a90fef740",
            "excerpt": "## Coverage confirming round (D-B2, repair delta db1f156a..f0c35260) — 0 findings. Six probes, each restored=yes.\nC1 closed: publish returns nil on an unpublished *durablefile.Error, whole package: bit at replacement_test.go:37.\nC3 closed: RecordMode 0o600 -> 0o640, whole package: bit, only TestOwnerReplacement.\nD16/D17 after the move: engine.go mode 0o644: bit, 3 failed via RequireFile.\nHelpers are not vacuous: wrong stage (replace.go:114 rename -> file-close) bit; Published true on pre-rename bit; fresh Err at replace.go:130 bit in 9 tests.\nAdvice: a self-test for durablefiletest assertions; hand-summarized G2-G4 failure rows in the preservation excerpt.\n"
          },
          "axis": "Coverage",
          "base": "ea32f1f8c853243a523b0113b0b9d3453c8dd135",
          "tip": "3a3dbbe1a471a4fdd48e24abe38359ff138a1e73",
          "finding_ids": [],
          "supersedes": [
            "dcm-b2-r1-coverage"
          ]
        },
        {
          "id": "dcm-b2-r3-standards",
          "performer": "claude:dcm_rev_b2_std3",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "b087942b45c60e7624a0737fd8018970569f1ce1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b2_std3",
            "digest": "sha256:b5073de6ceb431020e1433a6a4d8bdc0a74aa4781fc7f9a520a59d01f1c21c42",
            "excerpt": "## Standards second confirming round (D-B2, repair delta 40aef905..79926b97) — 0 findings\nDCM-B2-S2 closed, confidence 9: TestOwnerPreRenameFailure no longer calls durablefiletest.Siblings and keeps RequireFailure and the record-absence check (replacement_test.go:37-40); R1 is the necessity red of the absence check; L1 shows only the leaf tests own cleanup.\n"
          },
          "axis": "Standards",
          "base": "ea32f1f8c853243a523b0113b0b9d3453c8dd135",
          "tip": "6eed4ae5973ef9b9a9e9c77d8b6fcba5d3c698d2",
          "finding_ids": [],
          "supersedes": [
            "dcm-b2-r2-standards"
          ]
        },
        {
          "id": "dcm-b2-r3-spec",
          "performer": "claude:dcm_rev_b2_spec3",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "b087942b45c60e7624a0737fd8018970569f1ce1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b2_spec3",
            "digest": "sha256:9402952984dc6c6ac3029588447317410bb58f2f968a448f757037fe80fcc93d",
            "excerpt": "## Spec second confirming round (D-B2, repair delta 40aef905..79926b97) — 0 findings\nD18 intact (replacement_test.go:15-25). The D36 pre-rename refusal evidence is intact: the unpublished rename-stage failure (:37) and no published record (:38-40), with recorded reds C1 and R1. Cleanup belongs to the leaf rows D01-D05.\n"
          },
          "axis": "Spec",
          "base": "ea32f1f8c853243a523b0113b0b9d3453c8dd135",
          "tip": "6eed4ae5973ef9b9a9e9c77d8b6fcba5d3c698d2",
          "finding_ids": [],
          "supersedes": [
            "dcm-b2-r2-spec"
          ]
        },
        {
          "id": "dcm-b2-r3-coverage",
          "performer": "claude:dcm_rev_b2_cov3",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "b087942b45c60e7624a0737fd8018970569f1ce1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b2_cov3",
            "digest": "sha256:55d75daf980d59c8939bae927ba747a7f135c74242d08970af8afc27bd677b8e",
            "excerpt": "## Coverage second confirming round (D-B2, repair delta 40aef905..79926b97) — 0 findings\nD36 binds the pre-rename refusal result (error and absent record), both asserted. Probe: the caller writes root/.owner.json.stray after a failed Replace and returns the error: silent, restored=yes; not a gap, because it adds caller behavior that no acceptance row forbids.\n"
          },
          "axis": "Coverage",
          "base": "ea32f1f8c853243a523b0113b0b9d3453c8dd135",
          "tip": "6eed4ae5973ef9b9a9e9c77d8b6fcba5d3c698d2",
          "finding_ids": [],
          "supersedes": [
            "dcm-b2-r2-coverage"
          ]
        }
      ]
    },
    {
      "id": "D-B3",
      "base": "0cc52e2103645332e77432be7fc6404f977e4bb5",
      "tip": "3f43adb796df97c91b491c258db5365355d1aff6",
      "plan_digest": "sha256:0c585206337b1291a396cf3b404bad1b24a48b69e43d8dc10a70727f573ea631",
      "source_digest": "675566c7ba8a376743b520da7edd5780db0a13f3",
      "acceptance_rows": [
        "D20"
      ],
      "verification": [
        {
          "id": "dcm-t04-caller-tests",
          "performer": "claude:dcm_t04",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "b034f5f9f948e705edf9e57aa049f3cbfb8e95ab",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "c7a940a3e0bdf6af8653c2d8f726ec2fadd7f301",
            "digest": "sha256:7d0de11d341e727b81c32c403f37634774db503b107790fcd8835a937a24713a",
            "excerpt": "bench test --package ./internal/capturetx at c7a940a3e0bdf6af8653c2d8f726ec2fadd7f301 (FT354-durable-callers). TestCaptureReplacement (D20 post-rename directory-sync) and TestCapturePreRenameFailure (D36 pre-rename rename fault) pass.\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/capturetx,pass,212,26\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "caller-tests",
          "command": "bench test --package ./internal/capturetx",
          "exit_code": 0
        },
        {
          "id": "dcm-t04-root-conformance",
          "performer": "claude:dcm_t04",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "b034f5f9f948e705edf9e57aa049f3cbfb8e95ab",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "c7a940a3e0bdf6af8653c2d8f726ec2fadd7f301",
            "digest": "sha256:fba72a65af14b533a62edeb80fd7960fbd6ab09918316876e9140594ac7c5cd8",
            "excerpt": "bench test --check root-conformance at c7a940a3e0bdf6af8653c2d8f726ec2fadd7f301 (FT354-durable-callers)\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,18498,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "dcm-t04-omission-proof",
          "performer": "claude:dcm_t04",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "b034f5f9f948e705edf9e57aa049f3cbfb8e95ab",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "c7a940a3e0bdf6af8653c2d8f726ec2fadd7f301",
            "digest": "sha256:b4a2b9d7bf89dcce54c9a3f1e6b1a01c8da6ed42f86ef440656496e1a4c9105a",
            "excerpt": "D-B3 omission-proof probe at c7a940a3e0bdf6af8653c2d8f726ec2fadd7f301 (FT354-durable-callers), subject internal/capturetx/transaction.go, bench test --package ./internal/capturetx, no --run filter.\nSwap old (exact): \"\\treturn r.Replace(path, data, 0o644)\\n}\"\nSwap new (exact): the former sequence: \"\\ttmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+\\\".tmp-*\\\")\\n\\tif err != nil {\\n\\t\\treturn err\\n\\t}\\n\\tname := tmp.Name()\\n\\tdefer os.Remove(name)\\n\\tif _, err := tmp.Write(data); err != nil {\\n\\t\\ttmp.Close()\\n\\t\\treturn err\\n\\t}\\n\\tif err := tmp.Sync(); err != nil {\\n\\t\\ttmp.Close()\\n\\t\\treturn err\\n\\t}\\n\\tif err := tmp.Close(); err != nil {\\n\\t\\treturn err\\n\\t}\\n\\tif err := os.Chmod(name, 0o644); err != nil {\\n\\t\\treturn err\\n\\t}\\n\\treturn os.Rename(name, path)\\n}\"\nResult: verdict bit, 2 failed tests (TestCaptureReplacement D20 witness; TestCapturePreRenameFailure D36 pre-rename evidence), restored yes; package passes 26 tests after restore.\n\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,c7a940a3e0bdf6af8653c2d8f726ec2fadd7f301,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/capturetx/transaction.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/capturetx,all,passed,26\nspilled{lines=12,bytes=954,omitted_lines=3,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791826326202107373-a73004d9cef18703.out}\n  github.com/gibbonmi/bench/internal/capturetx,fail,131,26\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/capturetx,TestCapturePreRenameFailure,\"replacement_test.go:67: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected capture rename failure\",1\n  github.com/gibbonmi/bench/internal/capturetx,TestCaptureReplacement,\"replacement_test.go:46: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected capture directory-sync failure\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "omission-proof",
          "command": "bench test --package ./internal/capturetx",
          "exit_code": 0,
          "probe": {
            "mutation": "Bypass the shared owner in capture replace. The capture operation failure witness must turn red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "c7a940a3e0bdf6af8653c2d8f726ec2fadd7f301",
              "digest": "sha256:b4a2b9d7bf89dcce54c9a3f1e6b1a01c8da6ed42f86ef440656496e1a4c9105a",
              "excerpt": "D-B3 omission-proof probe at c7a940a3e0bdf6af8653c2d8f726ec2fadd7f301 (FT354-durable-callers), subject internal/capturetx/transaction.go, bench test --package ./internal/capturetx, no --run filter.\nSwap old (exact): \"\\treturn r.Replace(path, data, 0o644)\\n}\"\nSwap new (exact): the former sequence: \"\\ttmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+\\\".tmp-*\\\")\\n\\tif err != nil {\\n\\t\\treturn err\\n\\t}\\n\\tname := tmp.Name()\\n\\tdefer os.Remove(name)\\n\\tif _, err := tmp.Write(data); err != nil {\\n\\t\\ttmp.Close()\\n\\t\\treturn err\\n\\t}\\n\\tif err := tmp.Sync(); err != nil {\\n\\t\\ttmp.Close()\\n\\t\\treturn err\\n\\t}\\n\\tif err := tmp.Close(); err != nil {\\n\\t\\treturn err\\n\\t}\\n\\tif err := os.Chmod(name, 0o644); err != nil {\\n\\t\\treturn err\\n\\t}\\n\\treturn os.Rename(name, path)\\n}\"\nResult: verdict bit, 2 failed tests (TestCaptureReplacement D20 witness; TestCapturePreRenameFailure D36 pre-rename evidence), restored yes; package passes 26 tests after restore.\n\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,c7a940a3e0bdf6af8653c2d8f726ec2fadd7f301,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/capturetx/transaction.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/capturetx,all,passed,26\nspilled{lines=12,bytes=954,omitted_lines=3,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791826326202107373-a73004d9cef18703.out}\n  github.com/gibbonmi/bench/internal/capturetx,fail,131,26\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/capturetx,TestCapturePreRenameFailure,\"replacement_test.go:67: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected capture rename failure\",1\n  github.com/gibbonmi/bench/internal/capturetx,TestCaptureReplacement,\"replacement_test.go:46: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected capture directory-sync failure\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dcm-t04-preservation",
          "performer": "claude:dcm_t04",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "b034f5f9f948e705edf9e57aa049f3cbfb8e95ab",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "c7a940a3e0bdf6af8653c2d8f726ec2fadd7f301",
            "digest": "sha256:3e2027a2fcf9db842d3fd8534381d7b19d9dca64dd88109bb4e5edfe40a19a95",
            "excerpt": "D-B3 preservation probes at c7a940a3e0bdf6af8653c2d8f726ec2fadd7f301 (FT354-durable-callers), subject internal/capturetx/transaction.go, bench test --package ./internal/capturetx, no --run filter. Each red names only the new independent caller assertions; restored yes, and the package passes 26 tests after restore.\nProbe 1 (payload literal necessity). Swap old (exact): \"return replaceWith(r, filepath.Join(dir, manifestName), append(data, '\\n'))\" new (exact): \"return replaceWith(r, filepath.Join(dir, manifestName), data)\". Verdict bit, 1 failed test: TestCaptureReplacement (sealedManifest bytes), restored yes.\n\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,c7a940a3e0bdf6af8653c2d8f726ec2fadd7f301,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/capturetx/transaction.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/capturetx,all,passed,26\nspilled{lines=11,bytes=909,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791826329705781023-7a98ad56d30c2d2d.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/capturetx,fail,170,26\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/capturetx,TestCaptureReplacement,\"replacement_test.go:47: /dev/shm/BYTEJG/t/TestCaptureReplacement3696050737/001/bench-capture-drain/manifest.json = \\\"{\\\\\\\\n  \\\\\\\\\\\"schema\\\\\\\\\\\": 1,\\\\\\\\n  \\\\\\\\\\\"id\\\\\\\\\\\": \\\\\\\\\\\"d-0123456789ab\\\\\\\\\\\",\\\\\\\\n  \\\\\\\\\\\"state\\\\\\\\\\\": \\\\\\\\\\\"sealed\\\\\\\\\\\",\\\\\\\\n  \\\\\\\\\\\"sources\\\\\\\\\\\": [\\\\\\\\n    {\\\\\\\\n      \\\\\\\\\\\"name\\\\\\\\\\\": \\\\\\\\\\\"ide… (672 bytes)\",1\nskips[0]{package,test,reason}:\n\nProbe 2 (mode literal necessity). Swap old (exact): \"return r.Replace(path, data, 0o644)\" new (exact): \"return r.Replace(path, data, 0o640)\". Verdict bit, 2 failed tests: TestCaptureReplacement and TestCapturePreRenameFailure (D36 pre-rename: former bytes at 0644), restored yes.\n\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,c7a940a3e0bdf6af8653c2d8f726ec2fadd7f301,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/capturetx/transaction.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/capturetx,all,passed,26\nspilled{lines=12,bytes=970,omitted_lines=3,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791826332990128801-881eb3fe5a3028a0.out}\n  github.com/gibbonmi/bench/internal/capturetx,fail,165,26\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/capturetx,TestCapturePreRenameFailure,\"replacement_test.go:68: /dev/shm/JX4MH5/t/TestCapturePreRenameFailure1487139520/001/manifest.json mode = -rw-r-----, want a regular file at -rw-r--r--\",1\n  github.com/gibbonmi/bench/internal/capturetx,TestCaptureReplacement,\"replacement_test.go:47: /dev/shm/JX4MH5/t/TestCaptureReplacement2756281889/001/bench-capture-drain/manifest.json mode = -rw-r-----, want a regular file at -rw-r--r--\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "preservation",
          "command": "bench test --package ./internal/capturetx",
          "exit_code": 0,
          "probe": {
            "mutation": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "c7a940a3e0bdf6af8653c2d8f726ec2fadd7f301",
              "digest": "sha256:3e2027a2fcf9db842d3fd8534381d7b19d9dca64dd88109bb4e5edfe40a19a95",
              "excerpt": "D-B3 preservation probes at c7a940a3e0bdf6af8653c2d8f726ec2fadd7f301 (FT354-durable-callers), subject internal/capturetx/transaction.go, bench test --package ./internal/capturetx, no --run filter. Each red names only the new independent caller assertions; restored yes, and the package passes 26 tests after restore.\nProbe 1 (payload literal necessity). Swap old (exact): \"return replaceWith(r, filepath.Join(dir, manifestName), append(data, '\\n'))\" new (exact): \"return replaceWith(r, filepath.Join(dir, manifestName), data)\". Verdict bit, 1 failed test: TestCaptureReplacement (sealedManifest bytes), restored yes.\n\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,c7a940a3e0bdf6af8653c2d8f726ec2fadd7f301,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/capturetx/transaction.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/capturetx,all,passed,26\nspilled{lines=11,bytes=909,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791826329705781023-7a98ad56d30c2d2d.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/capturetx,fail,170,26\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/capturetx,TestCaptureReplacement,\"replacement_test.go:47: /dev/shm/BYTEJG/t/TestCaptureReplacement3696050737/001/bench-capture-drain/manifest.json = \\\"{\\\\\\\\n  \\\\\\\\\\\"schema\\\\\\\\\\\": 1,\\\\\\\\n  \\\\\\\\\\\"id\\\\\\\\\\\": \\\\\\\\\\\"d-0123456789ab\\\\\\\\\\\",\\\\\\\\n  \\\\\\\\\\\"state\\\\\\\\\\\": \\\\\\\\\\\"sealed\\\\\\\\\\\",\\\\\\\\n  \\\\\\\\\\\"sources\\\\\\\\\\\": [\\\\\\\\n    {\\\\\\\\n      \\\\\\\\\\\"name\\\\\\\\\\\": \\\\\\\\\\\"ide… (672 bytes)\",1\nskips[0]{package,test,reason}:\n\nProbe 2 (mode literal necessity). Swap old (exact): \"return r.Replace(path, data, 0o644)\" new (exact): \"return r.Replace(path, data, 0o640)\". Verdict bit, 2 failed tests: TestCaptureReplacement and TestCapturePreRenameFailure (D36 pre-rename: former bytes at 0644), restored yes.\n\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,c7a940a3e0bdf6af8653c2d8f726ec2fadd7f301,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/capturetx/transaction.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/capturetx,all,passed,26\nspilled{lines=12,bytes=970,omitted_lines=3,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791826332990128801-881eb3fe5a3028a0.out}\n  github.com/gibbonmi/bench/internal/capturetx,fail,165,26\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/capturetx,TestCapturePreRenameFailure,\"replacement_test.go:68: /dev/shm/JX4MH5/t/TestCapturePreRenameFailure1487139520/001/manifest.json mode = -rw-r-----, want a regular file at -rw-r--r--\",1\n  github.com/gibbonmi/bench/internal/capturetx,TestCaptureReplacement,\"replacement_test.go:47: /dev/shm/JX4MH5/t/TestCaptureReplacement2756281889/001/bench-capture-drain/manifest.json mode = -rw-r-----, want a regular file at -rw-r--r--\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dcm-t04-rep1-caller-tests",
          "performer": "claude:dcm_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "675566c7ba8a376743b520da7edd5780db0a13f3",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "3f43adb796df97c91b491c258db5365355d1aff6",
            "digest": "sha256:0679688923dc85448757f5fd1e362d19488da04592c9907c199d11e26aede367",
            "excerpt": "D-B3 caller-tests at repair commit 3f43adb796df97c91b491c258db5365355d1aff6 (FT354-durable-callers), command bench test --package ./internal/capturetx. TestCaptureReplacement (D20, published directory-sync fault, mode 0644 and encoded bytes) and TestCapturePreRenameFailure (D36, unpublished rename fault, former manifest bytes) pass.\n\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,3f43adb796df97c91b491c258db5365355d1aff6,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/capturetx,pass,154,26\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "caller-tests",
          "command": "bench test --package ./internal/capturetx",
          "exit_code": 0
        },
        {
          "id": "dcm-t04-rep1-root-conformance",
          "performer": "claude:dcm_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "675566c7ba8a376743b520da7edd5780db0a13f3",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "3f43adb796df97c91b491c258db5365355d1aff6",
            "digest": "sha256:f133bb11f1ae378b33c4568a61d44524303894b7bc6e0c1e2c77cf3df1bc4d37",
            "excerpt": "D-B3 root-conformance at repair commit 3f43adb796df97c91b491c258db5365355d1aff6 (FT354-durable-callers), command bench test --check root-conformance. The dirty cell is reviews/durable-caller-migration.md only.\n\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,3f43adb796df97c91b491c258db5365355d1aff6,true\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17171,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "dcm-t04-rep1-omission-proof",
          "performer": "claude:dcm_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "675566c7ba8a376743b520da7edd5780db0a13f3",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "3f43adb796df97c91b491c258db5365355d1aff6",
            "digest": "sha256:e7e6882b43728c95990201efae6e880f504d7197e10510adce0b05842a2b0008",
            "excerpt": "D-B3 omission-proof probe at repair commit 3f43adb796df97c91b491c258db5365355d1aff6 (FT354-durable-callers), subject internal/capturetx/transaction.go, bench test --package ./internal/capturetx, no --run filter.\nSwap old (exact): \"\\treturn r.Replace(path, data, 0o644)\\n}\"\nSwap new (exact): the former sequence: \"\\ttmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+\\\".tmp-*\\\")\\n\\tif err != nil {\\n\\t\\treturn err\\n\\t}\\n\\tname := tmp.Name()\\n\\tdefer os.Remove(name)\\n\\tif _, err := tmp.Write(data); err != nil {\\n\\t\\ttmp.Close()\\n\\t\\treturn err\\n\\t}\\n\\tif err := tmp.Sync(); err != nil {\\n\\t\\ttmp.Close()\\n\\t\\treturn err\\n\\t}\\n\\tif err := tmp.Close(); err != nil {\\n\\t\\treturn err\\n\\t}\\n\\tif err := os.Chmod(name, 0o644); err != nil {\\n\\t\\treturn err\\n\\t}\\n\\treturn os.Rename(name, path)\\n}\"\nResult: verdict bit, 2 failed tests (TestCaptureReplacement D20 witness; TestCapturePreRenameFailure D36 pre-rename evidence), restored yes.\n\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,3f43adb796df97c91b491c258db5365355d1aff6,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/capturetx/transaction.go,swap,failed,2,yes\n  github.com/gibbonmi/bench/internal/capturetx,fail,205,26\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/capturetx,TestCapturePreRenameFailure,\"replacement_test.go:66: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected capture rename failure\",1\n  github.com/gibbonmi/bench/internal/capturetx,TestCaptureReplacement,\"replacement_test.go:45: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected capture directory-sync failure\",1\n"
          },
          "requirement": "omission-proof",
          "command": "bench test --package ./internal/capturetx",
          "exit_code": 0,
          "probe": {
            "mutation": "Bypass the shared owner in capture replace. The capture operation failure witness must turn red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "3f43adb796df97c91b491c258db5365355d1aff6",
              "digest": "sha256:e7e6882b43728c95990201efae6e880f504d7197e10510adce0b05842a2b0008",
              "excerpt": "D-B3 omission-proof probe at repair commit 3f43adb796df97c91b491c258db5365355d1aff6 (FT354-durable-callers), subject internal/capturetx/transaction.go, bench test --package ./internal/capturetx, no --run filter.\nSwap old (exact): \"\\treturn r.Replace(path, data, 0o644)\\n}\"\nSwap new (exact): the former sequence: \"\\ttmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+\\\".tmp-*\\\")\\n\\tif err != nil {\\n\\t\\treturn err\\n\\t}\\n\\tname := tmp.Name()\\n\\tdefer os.Remove(name)\\n\\tif _, err := tmp.Write(data); err != nil {\\n\\t\\ttmp.Close()\\n\\t\\treturn err\\n\\t}\\n\\tif err := tmp.Sync(); err != nil {\\n\\t\\ttmp.Close()\\n\\t\\treturn err\\n\\t}\\n\\tif err := tmp.Close(); err != nil {\\n\\t\\treturn err\\n\\t}\\n\\tif err := os.Chmod(name, 0o644); err != nil {\\n\\t\\treturn err\\n\\t}\\n\\treturn os.Rename(name, path)\\n}\"\nResult: verdict bit, 2 failed tests (TestCaptureReplacement D20 witness; TestCapturePreRenameFailure D36 pre-rename evidence), restored yes.\n\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,3f43adb796df97c91b491c258db5365355d1aff6,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/capturetx/transaction.go,swap,failed,2,yes\n  github.com/gibbonmi/bench/internal/capturetx,fail,205,26\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/capturetx,TestCapturePreRenameFailure,\"replacement_test.go:66: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected capture rename failure\",1\n  github.com/gibbonmi/bench/internal/capturetx,TestCaptureReplacement,\"replacement_test.go:45: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected capture directory-sync failure\",1\n"
            }
          }
        },
        {
          "id": "dcm-t04-rep1-preservation",
          "performer": "claude:dcm_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "675566c7ba8a376743b520da7edd5780db0a13f3",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "3f43adb796df97c91b491c258db5365355d1aff6",
            "digest": "sha256:e64e8aacf1e1fa2c20a5a6f7b7cff0d269c047fdeaabdabbf4a50194c6cb5767",
            "excerpt": "D-B3 preservation probes at repair commit 3f43adb796df97c91b491c258db5365355d1aff6 (FT354-durable-callers), subject internal/capturetx/transaction.go, bench test --package ./internal/capturetx, no --run filter. Each red names only the independent caller assertion it guards; every probe restored yes on a clean tree.\n\nProbe 1 (payload newline necessity). Swap old (exact): \"return replaceWith(r, filepath.Join(dir, manifestName), append(data, '\\n'))\" new (exact): \"return replaceWith(r, filepath.Join(dir, manifestName), data)\". Verdict bit, 1 failed test, restored yes:\n  github.com/gibbonmi/bench/internal/capturetx,TestCaptureReplacement,\"replacement_test.go:46: .../bench-capture-drain/manifest.json = \\\"{\\\\n  \\\\\\\"schema\\\\\\\": 1, ... (672 bytes)\",1\n\nProbe 2 (mode literal necessity). Swap old (exact): \"return r.Replace(path, data, 0o644)\" new (exact): \"return r.Replace(path, data, 0o640)\". Verdict bit, 1 failed test, restored yes. The pre-rename row no longer repeats the mode check (S1):\n  github.com/gibbonmi/bench/internal/capturetx,TestCaptureReplacement,\"replacement_test.go:46: .../bench-capture-drain/manifest.json mode = -rw-r-----, want a regular file at -rw-r--r--\",1\n\nProbe 3 (MkdirAll omission, C3; guards \"creates its missing transaction directory\"). Omit (exact): \"\\tif err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {\\n\\t\\treturn err\\n\\t}\\n\". Verdict bit, 2 failed tests, restored yes:\n  github.com/gibbonmi/bench/internal/capturetx,TestAppendSerializesConcurrentWriters,\"store_test.go:35: replace .../capture/IDEAS.md: create: open .../capture/.IDEAS.md.4257275427: no such file ...\",1\n  github.com/gibbonmi/bench/internal/capturetx,TestCaptureReplacement,\"replacement_test.go:45: err = replace .../bench-capture-drain/manifest.json: create: open .../bench-capture-drain/.manifest.json.198650151... (390 bytes)\",1\n\nProbe 4 (former-bytes necessity, S1; fallback write after an unpublished error). Swap old (exact): \"\\treturn r.Replace(path, data, 0o644)\\n}\" new (exact): \"\\terr := r.Replace(path, data, 0o644)\\n\\tvar replaceErr *durablefile.Error\\n\\tif errors.As(err, &replaceErr) && !replaceErr.Published {\\n\\t\\t_ = os.WriteFile(path, data, 0o644)\\n\\t}\\n\\treturn err\\n}\". Verdict bit, 1 failed test, restored yes:\n  github.com/gibbonmi/bench/internal/capturetx,TestCapturePreRenameFailure,\"replacement_test.go:72: manifest = \\\"{\\\\n  \\\\\\\"schema\\\\\\\": 1, ... \\\\\\\"state\\\\\\\": \\\\\\\"sealed\\\\\\\", ...\",1\n\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,3f43adb796df97c91b491c258db5365355d1aff6,false\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/capturetx,all,passed,26\n"
          },
          "requirement": "preservation",
          "command": "bench test --package ./internal/capturetx",
          "exit_code": 0,
          "probe": {
            "mutation": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "3f43adb796df97c91b491c258db5365355d1aff6",
              "digest": "sha256:e64e8aacf1e1fa2c20a5a6f7b7cff0d269c047fdeaabdabbf4a50194c6cb5767",
              "excerpt": "D-B3 preservation probes at repair commit 3f43adb796df97c91b491c258db5365355d1aff6 (FT354-durable-callers), subject internal/capturetx/transaction.go, bench test --package ./internal/capturetx, no --run filter. Each red names only the independent caller assertion it guards; every probe restored yes on a clean tree.\n\nProbe 1 (payload newline necessity). Swap old (exact): \"return replaceWith(r, filepath.Join(dir, manifestName), append(data, '\\n'))\" new (exact): \"return replaceWith(r, filepath.Join(dir, manifestName), data)\". Verdict bit, 1 failed test, restored yes:\n  github.com/gibbonmi/bench/internal/capturetx,TestCaptureReplacement,\"replacement_test.go:46: .../bench-capture-drain/manifest.json = \\\"{\\\\n  \\\\\\\"schema\\\\\\\": 1, ... (672 bytes)\",1\n\nProbe 2 (mode literal necessity). Swap old (exact): \"return r.Replace(path, data, 0o644)\" new (exact): \"return r.Replace(path, data, 0o640)\". Verdict bit, 1 failed test, restored yes. The pre-rename row no longer repeats the mode check (S1):\n  github.com/gibbonmi/bench/internal/capturetx,TestCaptureReplacement,\"replacement_test.go:46: .../bench-capture-drain/manifest.json mode = -rw-r-----, want a regular file at -rw-r--r--\",1\n\nProbe 3 (MkdirAll omission, C3; guards \"creates its missing transaction directory\"). Omit (exact): \"\\tif err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {\\n\\t\\treturn err\\n\\t}\\n\". Verdict bit, 2 failed tests, restored yes:\n  github.com/gibbonmi/bench/internal/capturetx,TestAppendSerializesConcurrentWriters,\"store_test.go:35: replace .../capture/IDEAS.md: create: open .../capture/.IDEAS.md.4257275427: no such file ...\",1\n  github.com/gibbonmi/bench/internal/capturetx,TestCaptureReplacement,\"replacement_test.go:45: err = replace .../bench-capture-drain/manifest.json: create: open .../bench-capture-drain/.manifest.json.198650151... (390 bytes)\",1\n\nProbe 4 (former-bytes necessity, S1; fallback write after an unpublished error). Swap old (exact): \"\\treturn r.Replace(path, data, 0o644)\\n}\" new (exact): \"\\terr := r.Replace(path, data, 0o644)\\n\\tvar replaceErr *durablefile.Error\\n\\tif errors.As(err, &replaceErr) && !replaceErr.Published {\\n\\t\\t_ = os.WriteFile(path, data, 0o644)\\n\\t}\\n\\treturn err\\n}\". Verdict bit, 1 failed test, restored yes:\n  github.com/gibbonmi/bench/internal/capturetx,TestCapturePreRenameFailure,\"replacement_test.go:72: manifest = \\\"{\\\\n  \\\\\\\"schema\\\\\\\": 1, ... \\\\\\\"state\\\\\\\": \\\\\\\"sealed\\\\\\\", ...\",1\n\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,3f43adb796df97c91b491c258db5365355d1aff6,false\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/capturetx,all,passed,26\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "dcm-b3-r1-standards",
          "performer": "claude:dcm_rev_b3_std",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "b034f5f9f948e705edf9e57aa049f3cbfb8e95ab",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b3_std",
            "digest": "sha256:60585ff770be7e5499d50c236d5f2a8091243df55fe009181ab548326934a911",
            "excerpt": "## Standards (D-B3, 0cc52e21..cbccb028) — 3 findings; worst S1\nS1 | auto-fix (orchestrator: in part) | 5 | internal/capturetx/replacement_test.go:68 | the pre-rename row repeats the 0644 mode check that :47 owns; the former-bytes check needs a recorded necessity red | drop the repeated mode check; keep the caller-observable former-bytes check with its necessity red\nS2 | auto-fix | 6 | internal/capturetx/transaction.go:245-246, :208; replacement_test.go:12 | comments restate the code and the mode literal | keep only the reason for parent creation; delete the restating comments\nS3 | ask-user (orchestrator: no-op) | 4 | replacement_test.go:35, :50 | spec-row tags in test doc comments | no-op: row tags are an established tree convention (D-B1 Standards counted 92 test files)\nAdvice: two Native supply sites; a Middle Man chain in the manifest writers.\n"
          },
          "axis": "Standards",
          "base": "0cc52e2103645332e77432be7fc6404f977e4bb5",
          "tip": "c7a940a3e0bdf6af8653c2d8f726ec2fadd7f301",
          "finding_ids": [
            "S1",
            "S2",
            "S3"
          ],
          "supersedes": []
        },
        {
          "id": "dcm-b3-r1-spec",
          "performer": "claude:dcm_rev_b3_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "b034f5f9f948e705edf9e57aa049f3cbfb8e95ab",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b3_spec",
            "digest": "sha256:605482547dc4c3c392d479b60f93da63611110af63d00d4db0e532ce405be328",
            "excerpt": "## Spec (D-B3, 0cc52e21..cbccb028) — 0 findings\nD20 holds: 0644 and caller bytes through replaceWith; Append, Abort, recovery, and the manifest writer reach the shared owner; writeNew and directory promotion unchanged; parent creation stays with the caller; the manifest writer is a sound narrow caller writer; TestCapturePreRenameFailure is the D36 evidence. The four plan commits change only the execution block.\nAdvice: A1 the witness reaches replace only through replaceWith; A2 a post-rename Append failure leaves visible bytes, which the spec gives to the caller.\n"
          },
          "axis": "Spec",
          "base": "0cc52e2103645332e77432be7fc6404f977e4bb5",
          "tip": "c7a940a3e0bdf6af8653c2d8f726ec2fadd7f301",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "dcm-b3-r1-coverage",
          "performer": "claude:dcm_rev_b3_cov",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "b034f5f9f948e705edf9e57aa049f3cbfb8e95ab",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b3_cov",
            "digest": "sha256:df821e91679f3dc0ed01ce7f6de8748192eeef47d10b7eeef4984cc831152e93",
            "excerpt": "## Coverage (D-B3, 0cc52e21..cbccb028) — 3 findings; worst C1. Six probes, each restored=yes.\nC1 | ask-user (orchestrator: no-op) | 6 | internal/capturetx/transaction.go:241-242, :204-205 | a bypass of the native bindings in replace or writeManifestDir is silent | no-op: the spec forbids ambient swaps, so the D-D source census owns each native binding (ticket 14)\nC2 | ask-user (orchestrator: no-op) | 6 | internal/capturetx/public.go:49 | an Append that discards its replacement error is silent | no-op: spec line 83 accepts a narrow caller writer, and public.go is outside the fence\nC3 | auto-fix | 7 | reviews/durable-caller-migration.md (dcm-t04-preservation) | the MkdirAll omission red that guards \"creates its missing transaction directory\" is not recorded | record the MkdirAll omission red with exact swap text\nBit: replaceWith returns nil when Published (TestCaptureReplacement:46); a fallback os.WriteFile after an unpublished error (TestCapturePreRenameFailure:68). Silent, no requirement: parent mode 0o700.\n"
          },
          "axis": "Coverage",
          "base": "0cc52e2103645332e77432be7fc6404f977e4bb5",
          "tip": "c7a940a3e0bdf6af8653c2d8f726ec2fadd7f301",
          "finding_ids": [
            "C1",
            "C2",
            "C3"
          ],
          "supersedes": []
        },
        {
          "id": "dcm-b3-r2-standards",
          "performer": "claude:dcm_rev_b3_std2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "675566c7ba8a376743b520da7edd5780db0a13f3",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b3_std2",
            "digest": "sha256:9161869db1e9db25ec61f0d278b76fd69805714bbfb5a35be22ad9aa4bdcd155",
            "excerpt": "## Standards confirming round (D-B3, repair delta 69e7ecc7..13925a7d) — 0 findings\nS1 closed: the pre-rename row compares the former bytes only (replacement_test.go:67-73); mode is checked only at :46; necessity red probe 4 recorded. S2 closed: the restating comments are gone; transaction.go:244-245 gives only the parent-creation reason.\nAdvice: test doc comments at :34-36 still repeat \"at mode 0644\".\n"
          },
          "axis": "Standards",
          "base": "0cc52e2103645332e77432be7fc6404f977e4bb5",
          "tip": "3f43adb796df97c91b491c258db5365355d1aff6",
          "finding_ids": [],
          "supersedes": [
            "dcm-b3-r1-standards"
          ]
        },
        {
          "id": "dcm-b3-r2-spec",
          "performer": "claude:dcm_rev_b3_spec2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "675566c7ba8a376743b520da7edd5780db0a13f3",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b3_spec2",
            "digest": "sha256:6bc6f84ea456e719b989a068ab43b8bde60b187ebf8ed938e7744aee4270581b",
            "excerpt": "## Spec confirming round (D-B3, repair delta 69e7ecc7..13925a7d) — 0 findings\nD20 holds (TestCaptureReplacement unchanged, :37-47). D36 holds (TestCapturePreRenameFailure :51-74: unpublished rename failure and former bytes). The delta stays inside ticket 04 Writes.\n"
          },
          "axis": "Spec",
          "base": "0cc52e2103645332e77432be7fc6404f977e4bb5",
          "tip": "3f43adb796df97c91b491c258db5365355d1aff6",
          "finding_ids": [],
          "supersedes": [
            "dcm-b3-r1-spec"
          ]
        },
        {
          "id": "dcm-b3-r2-coverage",
          "performer": "claude:dcm_rev_b3_cov2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "675566c7ba8a376743b520da7edd5780db0a13f3",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b3_cov2",
            "digest": "sha256:d987f7a0d56bd15b41c6504cf1a57f877f063a29172b7f8953664674cf56ea3b",
            "excerpt": "## Coverage confirming round (D-B3, repair delta 69e7ecc7..13925a7d) — 0 findings\nC3 closed: MkdirAll omission recorded (probe 3, 2 tests). New probe: MkdirAll of the grandparent only: bit, 2 tests, restored=yes.\nS1 necessity: new probe writes the new bytes before the replacement in writeManifestWith: bit, TestCapturePreRenameFailure:72, restored=yes.\n"
          },
          "axis": "Coverage",
          "base": "0cc52e2103645332e77432be7fc6404f977e4bb5",
          "tip": "3f43adb796df97c91b491c258db5365355d1aff6",
          "finding_ids": [],
          "supersedes": [
            "dcm-b3-r1-coverage"
          ]
        }
      ]
    },
    {
      "id": "D-B4",
      "base": "c05092c603ff080970397570ac6a5c497d3663f4",
      "tip": "5e49631e1c32bff364b2b48a3c1b0afc17cbf8d7",
      "plan_digest": "sha256:908fc1e6123bcee6feb067edaa58110a42f1ba077bfcf2abe60805d1e61bc234",
      "source_digest": "253d32d0972a9499f4d7c52f935285154c9173c6",
      "acceptance_rows": [
        "D21"
      ],
      "verification": [
        {
          "id": "dcm-t05-caller-tests",
          "performer": "claude:dcm_t05",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d3d801af894c7a5d78dd5468fb21bc6811a0bb63",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "4fc0e3162cdc3e33779be1c49dcfd7ba51818eff",
            "digest": "sha256:e4b42df12401dea2e307497feff464f8d674b4c7d9b7a125682ea9c8f2154505",
            "excerpt": "bench test --package ./internal/handoffdoc at 4fc0e3162cdc3e33779be1c49dcfd7ba51818eff (FT354-durable-callers). TestHandoffReplacement (D21 post-rename directory-sync) and TestHandoffPreRenameFailure (D36 pre-rename rename fault; prior document kept) pass.\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/handoffdoc,pass,32,26\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "caller-tests",
          "command": "bench test --package ./internal/handoffdoc",
          "exit_code": 0
        },
        {
          "id": "dcm-t05-root-conformance",
          "performer": "claude:dcm_t05",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d3d801af894c7a5d78dd5468fb21bc6811a0bb63",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "4fc0e3162cdc3e33779be1c49dcfd7ba51818eff",
            "digest": "sha256:104cc037d93da491dafc821004e9b2227e1f67ae83ab6fee535c30c7c3ade446",
            "excerpt": "bench test --check root-conformance at 4fc0e3162cdc3e33779be1c49dcfd7ba51818eff (FT354-durable-callers)\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,19669,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "dcm-t05-omission-proof",
          "performer": "claude:dcm_t05",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d3d801af894c7a5d78dd5468fb21bc6811a0bb63",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "4fc0e3162cdc3e33779be1c49dcfd7ba51818eff",
            "digest": "sha256:7df322b02a7c16f53fa08b94b6be8433e475192b559c03f52bc0e2fc03e27302",
            "excerpt": "D-B4 omission-proof probes at 4fc0e3162cdc3e33779be1c49dcfd7ba51818eff (FT354-durable-callers), subject internal/handoffdoc/store.go, bench test --package ./internal/handoffdoc, no --run filter.\nProbe 1 (restore ignored directory errors). Swap old (exact): \"\tif err := replacer.Replace(path, []byte(doc.Render()), 0o644); err != nil {\"\nSwap new (exact): \"\tif err := replacer.Replace(path, []byte(doc.Render()), 0o644); err != nil && !func() bool { var published *durablefile.Error; return errors.As(err, &published) && published.Published }() {\"\nResult: verdict bit, 1 failed test (TestHandoffReplacement, D21 witness), restored yes.\n\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,4fc0e3162cdc3e33779be1c49dcfd7ba51818eff,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/handoffdoc,all,passed,26\nspilled{lines=11,bytes=729,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791828022592018805-3900910443a601c2.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/handoffdoc,fail,31,26\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/handoffdoc,TestHandoffReplacement,\"replacement_test.go:29: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected handoff directory-sync failure\",1\nskips[0]{package,test,reason}:\n\nProbe 2 (bypass the injected leaf with the native owner). Swap old (exact): \"\tif err := replacer.Replace(path, []byte(doc.Render()), 0o644); err != nil {\"\nSwap new (exact): \"\tif err := durablefile.Replace(path, []byte(doc.Render()), 0o644); err != nil {\"\nResult: verdict bit, 2 failed tests (TestHandoffReplacement D21 witness; TestHandoffPreRenameFailure D36 pre-rename evidence), restored yes; package passes 26 tests after restore.\n\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,4fc0e3162cdc3e33779be1c49dcfd7ba51818eff,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/handoffdoc,all,passed,26\nspilled{lines=12,bytes=952,omitted_lines=3,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791828033001097924-ddd758aac05191da.out}\n  github.com/gibbonmi/bench/internal/handoffdoc,fail,33,26\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/handoffdoc,TestHandoffPreRenameFailure,\"replacement_test.go:51: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected handoff rename failure\",1\n  github.com/gibbonmi/bench/internal/handoffdoc,TestHandoffReplacement,\"replacement_test.go:29: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected handoff directory-sync failure\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "omission-proof",
          "command": "bench test --package ./internal/handoffdoc",
          "exit_code": 0,
          "probe": {
            "mutation": "Restore ignored directory errors or bypass the leaf. The section-writing composition test must detect the missing failure.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "4fc0e3162cdc3e33779be1c49dcfd7ba51818eff",
              "digest": "sha256:7df322b02a7c16f53fa08b94b6be8433e475192b559c03f52bc0e2fc03e27302",
              "excerpt": "D-B4 omission-proof probes at 4fc0e3162cdc3e33779be1c49dcfd7ba51818eff (FT354-durable-callers), subject internal/handoffdoc/store.go, bench test --package ./internal/handoffdoc, no --run filter.\nProbe 1 (restore ignored directory errors). Swap old (exact): \"\tif err := replacer.Replace(path, []byte(doc.Render()), 0o644); err != nil {\"\nSwap new (exact): \"\tif err := replacer.Replace(path, []byte(doc.Render()), 0o644); err != nil && !func() bool { var published *durablefile.Error; return errors.As(err, &published) && published.Published }() {\"\nResult: verdict bit, 1 failed test (TestHandoffReplacement, D21 witness), restored yes.\n\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,4fc0e3162cdc3e33779be1c49dcfd7ba51818eff,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/handoffdoc,all,passed,26\nspilled{lines=11,bytes=729,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791828022592018805-3900910443a601c2.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/handoffdoc,fail,31,26\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/handoffdoc,TestHandoffReplacement,\"replacement_test.go:29: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected handoff directory-sync failure\",1\nskips[0]{package,test,reason}:\n\nProbe 2 (bypass the injected leaf with the native owner). Swap old (exact): \"\tif err := replacer.Replace(path, []byte(doc.Render()), 0o644); err != nil {\"\nSwap new (exact): \"\tif err := durablefile.Replace(path, []byte(doc.Render()), 0o644); err != nil {\"\nResult: verdict bit, 2 failed tests (TestHandoffReplacement D21 witness; TestHandoffPreRenameFailure D36 pre-rename evidence), restored yes; package passes 26 tests after restore.\n\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,4fc0e3162cdc3e33779be1c49dcfd7ba51818eff,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/handoffdoc,all,passed,26\nspilled{lines=12,bytes=952,omitted_lines=3,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791828033001097924-ddd758aac05191da.out}\n  github.com/gibbonmi/bench/internal/handoffdoc,fail,33,26\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/handoffdoc,TestHandoffPreRenameFailure,\"replacement_test.go:51: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected handoff rename failure\",1\n  github.com/gibbonmi/bench/internal/handoffdoc,TestHandoffReplacement,\"replacement_test.go:29: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected handoff directory-sync failure\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dcm-t05-preservation",
          "performer": "claude:dcm_t05",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d3d801af894c7a5d78dd5468fb21bc6811a0bb63",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "4fc0e3162cdc3e33779be1c49dcfd7ba51818eff",
            "digest": "sha256:73c58abd129962b5946db2e1674209eab380cfca816e095971e768b439a77739",
            "excerpt": "D-B4 preservation probes at 4fc0e3162cdc3e33779be1c49dcfd7ba51818eff (FT354-durable-callers), subject internal/handoffdoc/store.go, bench test --package ./internal/handoffdoc, no --run filter. Each red names only TestHandoffReplacement, the owner of the independent mode and bytes expectations (frozen from the pre-migration writer); restored yes, and the package passes 26 tests after restore. TestHandoffPreRenameFailure checks only the caller refusal outcome (prior document bytes).\nProbe 1 (mode literal necessity). Swap old (exact): \"[]byte(doc.Render()), 0o644)\" new (exact): \"[]byte(doc.Render()), 0o640)\"\nResult: verdict bit, 1 failed test (TestHandoffReplacement mode), restored yes.\n\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,4fc0e3162cdc3e33779be1c49dcfd7ba51818eff,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/handoffdoc,all,passed,26\nspilled{lines=11,bytes=722,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791828041671539024-337aa15fdd8b958b.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/handoffdoc,fail,36,26\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/handoffdoc,TestHandoffReplacement,\"replacement_test.go:30: /dev/shm/4573PW/t/TestHandoffReplacement1551975842/001/session-handoff.md mode = -rw-r-----, want a regular file at -rw-r--r--\",1\nskips[0]{package,test,reason}:\n\nProbe 2 (payload literal necessity). Swap old (exact): \"[]byte(doc.Render()), 0o644)\" new (exact): \"[]byte(doc.Render()+\\\"<!-- -->\\\\n\\\"), 0o644)\"\nResult: verdict bit, 1 failed test (TestHandoffReplacement bytes), restored yes.\n\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,4fc0e3162cdc3e33779be1c49dcfd7ba51818eff,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/handoffdoc,all,passed,26\nspilled{lines=11,bytes=867,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791828050717363019-7be50526b97d176b.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/handoffdoc,fail,34,26\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/handoffdoc,TestHandoffReplacement,\"replacement_test.go:30: /dev/shm/HV7SYV/t/TestHandoffReplacement1198845858/001/session-handoff.md = \\\"## main\\\\\\\\n\\\\\\\\nNext command:\\\\\\\\n\\\\\\\\n### State\\\\\\\\n\\\\\\\\n## request 9f2ab77\\\\\\\\n\\\\\\\\nRequest token: `9f2ab77`\\\\\\\\nNext command: `bench gate`\\\\\\\\n\\\\\\\\n### State\\\\\\\\n\\\\\\\\nThe pha… (429 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "preservation",
          "command": "bench test --package ./internal/handoffdoc",
          "exit_code": 0,
          "probe": {
            "mutation": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "4fc0e3162cdc3e33779be1c49dcfd7ba51818eff",
              "digest": "sha256:73c58abd129962b5946db2e1674209eab380cfca816e095971e768b439a77739",
              "excerpt": "D-B4 preservation probes at 4fc0e3162cdc3e33779be1c49dcfd7ba51818eff (FT354-durable-callers), subject internal/handoffdoc/store.go, bench test --package ./internal/handoffdoc, no --run filter. Each red names only TestHandoffReplacement, the owner of the independent mode and bytes expectations (frozen from the pre-migration writer); restored yes, and the package passes 26 tests after restore. TestHandoffPreRenameFailure checks only the caller refusal outcome (prior document bytes).\nProbe 1 (mode literal necessity). Swap old (exact): \"[]byte(doc.Render()), 0o644)\" new (exact): \"[]byte(doc.Render()), 0o640)\"\nResult: verdict bit, 1 failed test (TestHandoffReplacement mode), restored yes.\n\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,4fc0e3162cdc3e33779be1c49dcfd7ba51818eff,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/handoffdoc,all,passed,26\nspilled{lines=11,bytes=722,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791828041671539024-337aa15fdd8b958b.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/handoffdoc,fail,36,26\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/handoffdoc,TestHandoffReplacement,\"replacement_test.go:30: /dev/shm/4573PW/t/TestHandoffReplacement1551975842/001/session-handoff.md mode = -rw-r-----, want a regular file at -rw-r--r--\",1\nskips[0]{package,test,reason}:\n\nProbe 2 (payload literal necessity). Swap old (exact): \"[]byte(doc.Render()), 0o644)\" new (exact): \"[]byte(doc.Render()+\\\"<!-- -->\\\\n\\\"), 0o644)\"\nResult: verdict bit, 1 failed test (TestHandoffReplacement bytes), restored yes.\n\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,4fc0e3162cdc3e33779be1c49dcfd7ba51818eff,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/handoffdoc,all,passed,26\nspilled{lines=11,bytes=867,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791828050717363019-7be50526b97d176b.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/handoffdoc,fail,34,26\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/handoffdoc,TestHandoffReplacement,\"replacement_test.go:30: /dev/shm/HV7SYV/t/TestHandoffReplacement1198845858/001/session-handoff.md = \\\"## main\\\\\\\\n\\\\\\\\nNext command:\\\\\\\\n\\\\\\\\n### State\\\\\\\\n\\\\\\\\n## request 9f2ab77\\\\\\\\n\\\\\\\\nRequest token: `9f2ab77`\\\\\\\\nNext command: `bench gate`\\\\\\\\n\\\\\\\\n### State\\\\\\\\n\\\\\\\\nThe pha… (429 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dcm-t05-rep1-caller-tests",
          "performer": "claude:dcm_t05_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "54e5145003126d1570f2fe5cd55b825e71b6f3ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "5e6c8994901d776291446a2836e669b1997d7005",
            "digest": "sha256:7da247e7af7febc99432010b44438f44de218bd9c23a55a11ac05b33999c74b0",
            "excerpt": "bench test --package ./internal/handoffdoc at 5e6c8994901d776291446a2836e669b1997d7005 (FT354-durable-callers). TestHandoffReplacement (D21 post-rename rows directory-open, directory-sync, directory-close; Published=true; one bytes and 0644 mode assertion) and TestHandoffPreRenameFailure (D36 pre-rename rename fault; prior document kept) pass.\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,5e6c8994901d776291446a2836e669b1997d7005,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/handoffdoc,pass,30,29\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "caller-tests",
          "command": "bench test --package ./internal/handoffdoc",
          "exit_code": 0
        },
        {
          "id": "dcm-t05-rep1-root-conformance",
          "performer": "claude:dcm_t05_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "54e5145003126d1570f2fe5cd55b825e71b6f3ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "5e6c8994901d776291446a2836e669b1997d7005",
            "digest": "sha256:e02eb628c77033c34d441647b5b86a79836f71c6d6e57818bf7b9f1320b5a537",
            "excerpt": "bench test --check root-conformance at 5e6c8994901d776291446a2836e669b1997d7005 (FT354-durable-callers)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,5e6c8994901d776291446a2836e669b1997d7005,false\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15581,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "dcm-t05-rep1-omission-proof",
          "performer": "claude:dcm_t05_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "54e5145003126d1570f2fe5cd55b825e71b6f3ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "5e6c8994901d776291446a2836e669b1997d7005",
            "digest": "sha256:f70c696160d4194dfc24660275a851f0067ebb310af68aebd4f2e49e5017560a",
            "excerpt": "D-B4 omission-proof probes at 5e6c8994901d776291446a2836e669b1997d7005 (FT354-durable-callers), subject internal/handoffdoc/store.go, bench probe ... --package ./internal/handoffdoc, no --run filter. Each probe restored yes; the package passes 29 tests after restore.\nBefore the repair, at 12864fc1b52a7b3c331bd8f22ab27776760b1371, Probe 1 (directory-open) was silent: verdict silent, 0 failed tests, 26 tests passed, restored yes. That was the F2 gap.\n\nProbe 1 (ignore a directory-open error). Swap old (exact): \"\\tif err := replacer.Replace(path, []byte(doc.Render()), 0o644); err != nil {\"\nSwap new (exact): \"\\tif err := replacer.Replace(path, []byte(doc.Render()), 0o644); err != nil && !func() bool { var e *durablefile.Error; return errors.As(err, &e) && e.Stage == \\\"directory-open\\\" }() {\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/handoffdoc,TestHandoffReplacement/directory-open,\"replacement_test.go:31: err = <nil>, want a *durablefile.Error at stage directory-open with Published=true that wraps injected handoff directory-open failure\",1\n\nProbe 2 (ignore a directory-sync error). Swap old: same as Probe 1. Swap new (exact): the Probe 1 new text with \"directory-sync\" in place of \"directory-open\".\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/handoffdoc,TestHandoffReplacement/directory-sync,\"replacement_test.go:31: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected handoff directory-sync failure\",1\n\nProbe 3 (ignore a directory-close error). Swap old: same as Probe 1. Swap new (exact): the Probe 1 new text with \"directory-close\" in place of \"directory-open\".\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/handoffdoc,TestHandoffReplacement/directory-close,\"replacement_test.go:31: err = <nil>, want a *durablefile.Error at stage directory-close with Published=true that wraps injected handoff directory-close failure\",1\n\nProbe 4 (bypass the injected leaf with the native owner). Swap old (exact): \"\\tif err := replacer.Replace(path, []byte(doc.Render()), 0o644); err != nil {\"\nSwap new (exact): \"\\tif err := durablefile.Replace(path, []byte(doc.Render()), 0o644); err != nil {\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,4,yes\nfailures[4]: TestHandoffPreRenameFailure (replacement_test.go:55, err = <nil>, want stage rename Published=false), TestHandoffReplacement/directory-open, TestHandoffReplacement/directory-sync, TestHandoffReplacement/directory-close (replacement_test.go:31, err = <nil>).\n"
          },
          "requirement": "omission-proof",
          "command": "bench test --package ./internal/handoffdoc",
          "exit_code": 0,
          "probe": {
            "mutation": "Restore ignored directory errors or bypass the leaf. The section-writing composition test must detect the missing failure.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "5e6c8994901d776291446a2836e669b1997d7005",
              "digest": "sha256:f70c696160d4194dfc24660275a851f0067ebb310af68aebd4f2e49e5017560a",
              "excerpt": "D-B4 omission-proof probes at 5e6c8994901d776291446a2836e669b1997d7005 (FT354-durable-callers), subject internal/handoffdoc/store.go, bench probe ... --package ./internal/handoffdoc, no --run filter. Each probe restored yes; the package passes 29 tests after restore.\nBefore the repair, at 12864fc1b52a7b3c331bd8f22ab27776760b1371, Probe 1 (directory-open) was silent: verdict silent, 0 failed tests, 26 tests passed, restored yes. That was the F2 gap.\n\nProbe 1 (ignore a directory-open error). Swap old (exact): \"\\tif err := replacer.Replace(path, []byte(doc.Render()), 0o644); err != nil {\"\nSwap new (exact): \"\\tif err := replacer.Replace(path, []byte(doc.Render()), 0o644); err != nil && !func() bool { var e *durablefile.Error; return errors.As(err, &e) && e.Stage == \\\"directory-open\\\" }() {\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/handoffdoc,TestHandoffReplacement/directory-open,\"replacement_test.go:31: err = <nil>, want a *durablefile.Error at stage directory-open with Published=true that wraps injected handoff directory-open failure\",1\n\nProbe 2 (ignore a directory-sync error). Swap old: same as Probe 1. Swap new (exact): the Probe 1 new text with \"directory-sync\" in place of \"directory-open\".\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/handoffdoc,TestHandoffReplacement/directory-sync,\"replacement_test.go:31: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected handoff directory-sync failure\",1\n\nProbe 3 (ignore a directory-close error). Swap old: same as Probe 1. Swap new (exact): the Probe 1 new text with \"directory-close\" in place of \"directory-open\".\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/handoffdoc,TestHandoffReplacement/directory-close,\"replacement_test.go:31: err = <nil>, want a *durablefile.Error at stage directory-close with Published=true that wraps injected handoff directory-close failure\",1\n\nProbe 4 (bypass the injected leaf with the native owner). Swap old (exact): \"\\tif err := replacer.Replace(path, []byte(doc.Render()), 0o644); err != nil {\"\nSwap new (exact): \"\\tif err := durablefile.Replace(path, []byte(doc.Render()), 0o644); err != nil {\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,4,yes\nfailures[4]: TestHandoffPreRenameFailure (replacement_test.go:55, err = <nil>, want stage rename Published=false), TestHandoffReplacement/directory-open, TestHandoffReplacement/directory-sync, TestHandoffReplacement/directory-close (replacement_test.go:31, err = <nil>).\n"
            }
          }
        },
        {
          "id": "dcm-t05-rep1-preservation",
          "performer": "claude:dcm_t05_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "54e5145003126d1570f2fe5cd55b825e71b6f3ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "5e6c8994901d776291446a2836e669b1997d7005",
            "digest": "sha256:a2e9e3a0e19b0953b1600842d1712f05810b9b1138fafa2197ec2eeba1cfff44",
            "excerpt": "D-B4 preservation probes at 5e6c8994901d776291446a2836e669b1997d7005 (FT354-durable-callers), subject internal/handoffdoc/store.go, bench probe ... --package ./internal/handoffdoc, no --run filter. Each probe restored yes; the package passes 29 tests after restore. TestHandoffReplacement holds the one independent mode and bytes assertion (frozen from the pre-migration writer) for its three post-rename rows. TestHandoffPreRenameFailure checks only the caller refusal outcome (prior document bytes).\n\nProbe 1 (mode literal necessity). Swap old (exact): \"[]byte(doc.Render()), 0o644)\" new (exact): \"[]byte(doc.Render()), 0o640)\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,3,yes\nfailures[3]: TestHandoffReplacement/directory-open, /directory-sync, /directory-close, each \"replacement_test.go:32: .../session-handoff.md mode = -rw-r-----, want a regular file at -rw-r--r--\".\n\nProbe 2 (payload literal necessity). Swap old (exact): \"[]byte(doc.Render()), 0o644)\" new (exact): \"[]byte(doc.Render()+\\\"<!-- -->\\\\n\\\"), 0o644)\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,3,yes\nfailures[3]: TestHandoffReplacement/directory-open, /directory-sync, /directory-close, each \"replacement_test.go:32: .../session-handoff.md = \\\"## main\\\\n\\\\nNext command: ... (442-444 bytes)\".\n\nProbe 3 (prior-bytes necessity: a fallback write after an unpublished error). Swap old (exact): \"\\t\\treturn fmt.Errorf(\\\"write handoff document: %w\\\", err)\"\nSwap new (exact): \"\\t\\t_ = os.WriteFile(path, []byte(doc.Render()), 0o644)\\n\\t\\treturn fmt.Errorf(\\\"write handoff document: %w\\\", err)\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/handoffdoc,TestHandoffPreRenameFailure,\"replacement_test.go:61: document = \\\"## main\\\\n\\\\nNext command:\\\\n\\\\n### State\\\\n\\\\n## request 9f2ab77 ... The phase is live.\\\\n\\\\n\\\", want the prior document \\\"## main\\\\n\\\\nNext comman… (260 bytes)\",1\nThe red falls only at the prior-bytes check (line 61); the RequireFailure check at line 55 passes, because the error is still wrapped.\n"
          },
          "requirement": "preservation",
          "command": "bench test --package ./internal/handoffdoc",
          "exit_code": 0,
          "probe": {
            "mutation": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "5e6c8994901d776291446a2836e669b1997d7005",
              "digest": "sha256:a2e9e3a0e19b0953b1600842d1712f05810b9b1138fafa2197ec2eeba1cfff44",
              "excerpt": "D-B4 preservation probes at 5e6c8994901d776291446a2836e669b1997d7005 (FT354-durable-callers), subject internal/handoffdoc/store.go, bench probe ... --package ./internal/handoffdoc, no --run filter. Each probe restored yes; the package passes 29 tests after restore. TestHandoffReplacement holds the one independent mode and bytes assertion (frozen from the pre-migration writer) for its three post-rename rows. TestHandoffPreRenameFailure checks only the caller refusal outcome (prior document bytes).\n\nProbe 1 (mode literal necessity). Swap old (exact): \"[]byte(doc.Render()), 0o644)\" new (exact): \"[]byte(doc.Render()), 0o640)\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,3,yes\nfailures[3]: TestHandoffReplacement/directory-open, /directory-sync, /directory-close, each \"replacement_test.go:32: .../session-handoff.md mode = -rw-r-----, want a regular file at -rw-r--r--\".\n\nProbe 2 (payload literal necessity). Swap old (exact): \"[]byte(doc.Render()), 0o644)\" new (exact): \"[]byte(doc.Render()+\\\"<!-- -->\\\\n\\\"), 0o644)\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,3,yes\nfailures[3]: TestHandoffReplacement/directory-open, /directory-sync, /directory-close, each \"replacement_test.go:32: .../session-handoff.md = \\\"## main\\\\n\\\\nNext command: ... (442-444 bytes)\".\n\nProbe 3 (prior-bytes necessity: a fallback write after an unpublished error). Swap old (exact): \"\\t\\treturn fmt.Errorf(\\\"write handoff document: %w\\\", err)\"\nSwap new (exact): \"\\t\\t_ = os.WriteFile(path, []byte(doc.Render()), 0o644)\\n\\t\\treturn fmt.Errorf(\\\"write handoff document: %w\\\", err)\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/handoffdoc,TestHandoffPreRenameFailure,\"replacement_test.go:61: document = \\\"## main\\\\n\\\\nNext command:\\\\n\\\\n### State\\\\n\\\\n## request 9f2ab77 ... The phase is live.\\\\n\\\\n\\\", want the prior document \\\"## main\\\\n\\\\nNext comman… (260 bytes)\",1\nThe red falls only at the prior-bytes check (line 61); the RequireFailure check at line 55 passes, because the error is still wrapped.\n"
            }
          }
        },
        {
          "id": "dcm-t05-rep2-caller-tests",
          "performer": "claude:dcm_t05_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "253d32d0972a9499f4d7c52f935285154c9173c6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "5e49631e1c32bff364b2b48a3c1b0afc17cbf8d7",
            "digest": "sha256:8be8889ffd3884b26217152364cd4f6d21d37ebf7d513d90eb99265e0e07eb34",
            "excerpt": "bench test --package ./internal/handoffdoc at 5e49631e1c32bff364b2b48a3c1b0afc17cbf8d7 (FT354-durable-callers). TestHandoffReplacement ranges its D21 post-rename rows over durablefiletest.PostRenameStages (directory-open, directory-sync, directory-close; Published=true; one bytes and 0644 mode assertion), and TestHandoffPreRenameFailure (D36 rename fault; prior document kept) passes.\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,5e49631e1c32bff364b2b48a3c1b0afc17cbf8d7,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/handoffdoc,pass,29,29\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "caller-tests",
          "command": "bench test --package ./internal/handoffdoc",
          "exit_code": 0
        },
        {
          "id": "dcm-t05-rep2-durablefiletest",
          "performer": "claude:dcm_t05_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "253d32d0972a9499f4d7c52f935285154c9173c6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "5e49631e1c32bff364b2b48a3c1b0afc17cbf8d7",
            "digest": "sha256:627e57eb7c8e80758fdb42aaf252362b65e05d8c2ffb7508cae9cec21982b345",
            "excerpt": "bench test --package ./internal/durablefile/durablefiletest at 5e49631e1c32bff364b2b48a3c1b0afc17cbf8d7 (FT354-durable-callers). The package compiles and has no own tests; it exports PostRenameStages, RequireFailure, and RequireFile to the caller tests, which pass in the focused call over ./internal/handoffdoc, ./internal/durablefile, ./internal/capturetx, ./internal/gate/prospectiveartifact, and ./internal/gate.\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,5e49631e1c32bff364b2b48a3c1b0afc17cbf8d7,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile/durablefiletest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "durablefiletest",
          "command": "bench test --package ./internal/durablefile/durablefiletest",
          "exit_code": 0
        },
        {
          "id": "dcm-t05-rep2-root-conformance",
          "performer": "claude:dcm_t05_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "253d32d0972a9499f4d7c52f935285154c9173c6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "5e49631e1c32bff364b2b48a3c1b0afc17cbf8d7",
            "digest": "sha256:015ee3a4240bc24c1de0ee744cb5096c295b4bfa9250d10cb7b7bc9c5ee2a2fa",
            "excerpt": "bench test --check root-conformance at 5e49631e1c32bff364b2b48a3c1b0afc17cbf8d7 (FT354-durable-callers)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,5e49631e1c32bff364b2b48a3c1b0afc17cbf8d7,false\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15772,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "dcm-t05-rep2-omission-proof",
          "performer": "claude:dcm_t05_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "253d32d0972a9499f4d7c52f935285154c9173c6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "5e49631e1c32bff364b2b48a3c1b0afc17cbf8d7",
            "digest": "sha256:01eb4420884dbd2a8339f0b00cf0d2b923f00fb9bbe25c752f99b9d7d7285be5",
            "excerpt": "D-B4 omission-proof probes at 5e49631e1c32bff364b2b48a3c1b0afc17cbf8d7 (FT354-durable-callers), bench probe ... --package ./internal/handoffdoc, no --run filter. Each probe restored yes; the package passes 29 tests after restore. TestHandoffReplacement ranges over durablefiletest.PostRenameStages.\n\nProbe 0 (S2 set-ownership pair; subject internal/durablefile/durablefiletest/assert.go). Swap old (exact): \"var PostRenameStages = []string{\\\"directory-open\\\", \\\"directory-sync\\\", \\\"directory-close\\\"}\"\nSwap new (exact): \"var PostRenameStages = []string{\\\"directory-sync\\\", \\\"directory-close\\\"}\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  silent,internal/durablefile/durablefiletest/assert.go,swap,passed,0,yes\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/handoffdoc,pass,30,28\nThe handoff test loses its directory-open row (28 tests, not 29), so the row set comes from the exported set. Probe 1 is the pair: with the set intact, an ignored directory-open error turns that row red.\n\nProbe 1 (ignore a directory-open error; subject internal/handoffdoc/store.go). Swap old (exact): \"\\tif err := replacer.Replace(path, []byte(doc.Render()), 0o644); err != nil {\"\nSwap new (exact): \"\\tif err := replacer.Replace(path, []byte(doc.Render()), 0o644); err != nil && !func() bool { var e *durablefile.Error; return errors.As(err, &e) && e.Stage == \\\"directory-open\\\" }() {\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/handoffdoc,TestHandoffReplacement/directory-open,\"replacement_test.go:31: err = <nil>, want a *durablefile.Error at stage directory-open with Published=true that wraps injected handoff directory-open failure\",1\n\nProbe 2 (ignore a directory-sync error). Swap old: same as Probe 1. Swap new (exact): the Probe 1 new text with \"directory-sync\" in place of \"directory-open\".\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/handoffdoc,TestHandoffReplacement/directory-sync,\"replacement_test.go:31: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected handoff directory-sync failure\",1\n\nProbe 3 (ignore a directory-close error). Swap old: same as Probe 1. Swap new (exact): the Probe 1 new text with \"directory-close\" in place of \"directory-open\".\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/handoffdoc,TestHandoffReplacement/directory-close,\"replacement_test.go:31: err = <nil>, want a *durablefile.Error at stage directory-close with Published=true that wraps injected handoff directory-close failure\",1\n\nProbe 4 (bypass the injected leaf with the native owner). Swap old (exact): \"\\tif err := replacer.Replace(path, []byte(doc.Render()), 0o644); err != nil {\"\nSwap new (exact): \"\\tif err := durablefile.Replace(path, []byte(doc.Render()), 0o644); err != nil {\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,4,yes\nfailures[4]: TestHandoffPreRenameFailure (replacement_test.go:55, err = <nil>, want stage rename Published=false), TestHandoffReplacement/directory-open, TestHandoffReplacement/directory-sync, TestHandoffReplacement/directory-close (replacement_test.go:31, err = <nil>).\n"
          },
          "requirement": "omission-proof",
          "command": "bench test --package ./internal/handoffdoc",
          "exit_code": 0,
          "probe": {
            "mutation": "Restore ignored directory errors or bypass the leaf. The section-writing composition test must detect the missing failure.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "5e49631e1c32bff364b2b48a3c1b0afc17cbf8d7",
              "digest": "sha256:01eb4420884dbd2a8339f0b00cf0d2b923f00fb9bbe25c752f99b9d7d7285be5",
              "excerpt": "D-B4 omission-proof probes at 5e49631e1c32bff364b2b48a3c1b0afc17cbf8d7 (FT354-durable-callers), bench probe ... --package ./internal/handoffdoc, no --run filter. Each probe restored yes; the package passes 29 tests after restore. TestHandoffReplacement ranges over durablefiletest.PostRenameStages.\n\nProbe 0 (S2 set-ownership pair; subject internal/durablefile/durablefiletest/assert.go). Swap old (exact): \"var PostRenameStages = []string{\\\"directory-open\\\", \\\"directory-sync\\\", \\\"directory-close\\\"}\"\nSwap new (exact): \"var PostRenameStages = []string{\\\"directory-sync\\\", \\\"directory-close\\\"}\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  silent,internal/durablefile/durablefiletest/assert.go,swap,passed,0,yes\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/handoffdoc,pass,30,28\nThe handoff test loses its directory-open row (28 tests, not 29), so the row set comes from the exported set. Probe 1 is the pair: with the set intact, an ignored directory-open error turns that row red.\n\nProbe 1 (ignore a directory-open error; subject internal/handoffdoc/store.go). Swap old (exact): \"\\tif err := replacer.Replace(path, []byte(doc.Render()), 0o644); err != nil {\"\nSwap new (exact): \"\\tif err := replacer.Replace(path, []byte(doc.Render()), 0o644); err != nil && !func() bool { var e *durablefile.Error; return errors.As(err, &e) && e.Stage == \\\"directory-open\\\" }() {\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/handoffdoc,TestHandoffReplacement/directory-open,\"replacement_test.go:31: err = <nil>, want a *durablefile.Error at stage directory-open with Published=true that wraps injected handoff directory-open failure\",1\n\nProbe 2 (ignore a directory-sync error). Swap old: same as Probe 1. Swap new (exact): the Probe 1 new text with \"directory-sync\" in place of \"directory-open\".\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/handoffdoc,TestHandoffReplacement/directory-sync,\"replacement_test.go:31: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected handoff directory-sync failure\",1\n\nProbe 3 (ignore a directory-close error). Swap old: same as Probe 1. Swap new (exact): the Probe 1 new text with \"directory-close\" in place of \"directory-open\".\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/handoffdoc,TestHandoffReplacement/directory-close,\"replacement_test.go:31: err = <nil>, want a *durablefile.Error at stage directory-close with Published=true that wraps injected handoff directory-close failure\",1\n\nProbe 4 (bypass the injected leaf with the native owner). Swap old (exact): \"\\tif err := replacer.Replace(path, []byte(doc.Render()), 0o644); err != nil {\"\nSwap new (exact): \"\\tif err := durablefile.Replace(path, []byte(doc.Render()), 0o644); err != nil {\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,4,yes\nfailures[4]: TestHandoffPreRenameFailure (replacement_test.go:55, err = <nil>, want stage rename Published=false), TestHandoffReplacement/directory-open, TestHandoffReplacement/directory-sync, TestHandoffReplacement/directory-close (replacement_test.go:31, err = <nil>).\n"
            }
          }
        },
        {
          "id": "dcm-t05-rep2-preservation",
          "performer": "claude:dcm_t05_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "253d32d0972a9499f4d7c52f935285154c9173c6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "5e49631e1c32bff364b2b48a3c1b0afc17cbf8d7",
            "digest": "sha256:69ebb3be7bad3f85787f032b73e6eaaa10a1a4fa4b02d3ea634eeae9c047858c",
            "excerpt": "D-B4 preservation probes at 5e49631e1c32bff364b2b48a3c1b0afc17cbf8d7 (FT354-durable-callers), subject internal/handoffdoc/store.go, bench probe ... --package ./internal/handoffdoc, no --run filter. Each probe restored yes; the package passes 29 tests after restore. TestHandoffReplacement holds the one independent mode and bytes assertion (frozen from the pre-migration writer) for each row of durablefiletest.PostRenameStages. TestHandoffPreRenameFailure checks only the caller refusal outcome (prior document bytes).\n\nProbe 1 (mode literal necessity). Swap old (exact): \"[]byte(doc.Render()), 0o644)\" new (exact): \"[]byte(doc.Render()), 0o640)\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,3,yes\nfailures[3]: TestHandoffReplacement/directory-open, /directory-sync, /directory-close, each \"replacement_test.go:32: .../session-handoff.md mode = -rw-r-----, want a regular file at -rw-r--r--\".\n\nProbe 2 (payload literal necessity). Swap old (exact): \"[]byte(doc.Render()), 0o644)\" new (exact): \"[]byte(doc.Render()+\\\"<!-- -->\\\\n\\\"), 0o644)\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,3,yes\nfailures[3]: TestHandoffReplacement/directory-open, /directory-sync, /directory-close, each \"replacement_test.go:32: .../session-handoff.md = \\\"## main\\\\n\\\\nNext command: ... (443 bytes)\".\n\nProbe 3 (prior-bytes necessity: a fallback write after an unpublished error). Swap old (exact): \"\\t\\treturn fmt.Errorf(\\\"write handoff document: %w\\\", err)\"\nSwap new (exact): \"\\t\\t_ = os.WriteFile(path, []byte(doc.Render()), 0o644)\\n\\t\\treturn fmt.Errorf(\\\"write handoff document: %w\\\", err)\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/handoffdoc,TestHandoffPreRenameFailure,\"replacement_test.go:61: document = \\\"## main\\\\n\\\\nNext command:\\\\n\\\\n### State\\\\n\\\\n## request 9f2ab77 ... The phase is live.\\\\n\\\\n\\\", want the prior document \\\"## main\\\\n\\\\nNext comman… (260 bytes)\",1\nThe red falls only at the prior-bytes check (line 61); the RequireFailure check at line 55 passes, because the error is still wrapped.\n"
          },
          "requirement": "preservation",
          "command": "bench test --package ./internal/handoffdoc",
          "exit_code": 0,
          "probe": {
            "mutation": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "5e49631e1c32bff364b2b48a3c1b0afc17cbf8d7",
              "digest": "sha256:69ebb3be7bad3f85787f032b73e6eaaa10a1a4fa4b02d3ea634eeae9c047858c",
              "excerpt": "D-B4 preservation probes at 5e49631e1c32bff364b2b48a3c1b0afc17cbf8d7 (FT354-durable-callers), subject internal/handoffdoc/store.go, bench probe ... --package ./internal/handoffdoc, no --run filter. Each probe restored yes; the package passes 29 tests after restore. TestHandoffReplacement holds the one independent mode and bytes assertion (frozen from the pre-migration writer) for each row of durablefiletest.PostRenameStages. TestHandoffPreRenameFailure checks only the caller refusal outcome (prior document bytes).\n\nProbe 1 (mode literal necessity). Swap old (exact): \"[]byte(doc.Render()), 0o644)\" new (exact): \"[]byte(doc.Render()), 0o640)\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,3,yes\nfailures[3]: TestHandoffReplacement/directory-open, /directory-sync, /directory-close, each \"replacement_test.go:32: .../session-handoff.md mode = -rw-r-----, want a regular file at -rw-r--r--\".\n\nProbe 2 (payload literal necessity). Swap old (exact): \"[]byte(doc.Render()), 0o644)\" new (exact): \"[]byte(doc.Render()+\\\"<!-- -->\\\\n\\\"), 0o644)\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,3,yes\nfailures[3]: TestHandoffReplacement/directory-open, /directory-sync, /directory-close, each \"replacement_test.go:32: .../session-handoff.md = \\\"## main\\\\n\\\\nNext command: ... (443 bytes)\".\n\nProbe 3 (prior-bytes necessity: a fallback write after an unpublished error). Swap old (exact): \"\\t\\treturn fmt.Errorf(\\\"write handoff document: %w\\\", err)\"\nSwap new (exact): \"\\t\\t_ = os.WriteFile(path, []byte(doc.Render()), 0o644)\\n\\t\\treturn fmt.Errorf(\\\"write handoff document: %w\\\", err)\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/handoffdoc/store.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/handoffdoc,TestHandoffPreRenameFailure,\"replacement_test.go:61: document = \\\"## main\\\\n\\\\nNext command:\\\\n\\\\n### State\\\\n\\\\n## request 9f2ab77 ... The phase is live.\\\\n\\\\n\\\", want the prior document \\\"## main\\\\n\\\\nNext comman… (260 bytes)\",1\nThe red falls only at the prior-bytes check (line 61); the RequireFailure check at line 55 passes, because the error is still wrapped.\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "dcm-b4-r1-standards",
          "performer": "claude:dcm_rev_b4_std",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "d3d801af894c7a5d78dd5468fb21bc6811a0bb63",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b4_std",
            "digest": "sha256:5c9b3e7574f1d97b12a649fbb789d276a5ef840ce37c57de04d40bbbc161d1d5",
            "excerpt": "## Standards (D-B4, c05092c6..bce8ea45) — 1 finding; worst S1\nS1 | auto-fix | 7 | internal/handoffdoc/replacement_test.go:19-20 | the D21 comment ends \"which the former writer ignored.\", which narrates history | delete the clause\nChecks: no copy of the algorithm; one %w wrap; replacer seam with Native from the exported functions; shared assertions; independent expectations with recorded reds; each fact checked once.\nAdvice: replacementSection is a package-level fixture; a constructor would be safer.\n"
          },
          "axis": "Standards",
          "base": "c05092c603ff080970397570ac6a5c497d3663f4",
          "tip": "4fc0e3162cdc3e33779be1c49dcfd7ba51818eff",
          "finding_ids": [
            "S1"
          ],
          "supersedes": []
        },
        {
          "id": "dcm-b4-r1-spec",
          "performer": "claude:dcm_rev_b4_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "d3d801af894c7a5d78dd5468fb21bc6811a0bb63",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b4_spec",
            "digest": "sha256:31dcc66b8b7d3f2c76e4691b8e878bc267bb609c2d4d21efa623538bbf8e8879",
            "excerpt": "## Spec (D-B4, c05092c6..bce8ea45) — 0 findings\nD21 holds: 0644 and the rendered document at store.go:76; a directory failure returns an error through the %w wrap; replacer seam with Native from Update and WriteSection; the lock and section edits stay with the caller; TestHandoffPreRenameFailure is the D36 evidence. The two plan commits add only the 08 and 09 assignments.\n"
          },
          "axis": "Spec",
          "base": "c05092c603ff080970397570ac6a5c497d3663f4",
          "tip": "4fc0e3162cdc3e33779be1c49dcfd7ba51818eff",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "dcm-b4-r1-coverage",
          "performer": "claude:dcm_rev_b4_cov",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "d3d801af894c7a5d78dd5468fb21bc6811a0bb63",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b4_cov",
            "digest": "sha256:cfda4c90de727aa15bd1b1e3b75a46c378332050d56ed1c5b4406a8571a6c907",
            "excerpt": "## Coverage (D-B4, c05092c6..bce8ea45) — 2 findings; worst F1. Three probes, each restored=yes.\nF1 | auto-fix | 7 | internal/handoffdoc/replacement_test.go:56 | the D36 prior-bytes check has no recorded necessity red; probe A (os.WriteFile of the rendered document before the wrapped return) bit only at :57 | record probe A with exact swap text\nF2 | ask-user (orchestrator: auto-fix) | 4 | tickets/05-migrate-handoff-documents.md:10, store.go:76 | only directory-sync has a caller witness; ignoring directory-open errors is silent (probe C) | make TestHandoffReplacement table-driven over directory-open, directory-sync, and directory-close\nSilent, advice: dropping the caller prefix (probe B); the native bindings belong to the D-D census.\n"
          },
          "axis": "Coverage",
          "base": "c05092c603ff080970397570ac6a5c497d3663f4",
          "tip": "4fc0e3162cdc3e33779be1c49dcfd7ba51818eff",
          "finding_ids": [
            "F1",
            "F2"
          ],
          "supersedes": []
        },
        {
          "id": "dcm-b4-r2-standards",
          "performer": "claude:dcm_rev_b4_std2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "54e5145003126d1570f2fe5cd55b825e71b6f3ad",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b4_std2",
            "digest": "sha256:a42a856ef2bbc14bd8d967e17831668b38a3fcf0dd32681ec59dd557acddee86",
            "excerpt": "## Standards confirming round (D-B4, repair delta 12864fc1..92da1312) — 1 finding; S1 closed\nS1 closed: replacement_test.go:19 states the current behavior in present tense.\nS2 | ask-user (orchestrator: auto-fix by plan expansion) | 6 | internal/handoffdoc/replacement_test.go:22 | the table writes out the post-rename stage set that internal/durablefile/replace.go:121-127, durablefiletest/recorder.go, and afterRename in internal/durablefile/replace_test.go:25 already define | export the set once from durablefiletest and range over it\n"
          },
          "axis": "Standards",
          "base": "c05092c603ff080970397570ac6a5c497d3663f4",
          "tip": "5e6c8994901d776291446a2836e669b1997d7005",
          "finding_ids": [
            "S2"
          ],
          "supersedes": [
            "dcm-b4-r1-standards"
          ]
        },
        {
          "id": "dcm-b4-r2-spec",
          "performer": "claude:dcm_rev_b4_spec2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "54e5145003126d1570f2fe5cd55b825e71b6f3ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b4_spec2",
            "digest": "sha256:a0a91f8cff43295139b638fa95f63b7538a5681089ba964f65f604300b623aed",
            "excerpt": "## Spec confirming round (D-B4, repair delta 12864fc1..92da1312) — 0 findings\nD21 and its citation intact; the three directory-stage rows meet the behavior-change rule (spec.md:233); D36 pre-rename evidence unchanged; the delta stays inside ticket 05 Writes.\n"
          },
          "axis": "Spec",
          "base": "c05092c603ff080970397570ac6a5c497d3663f4",
          "tip": "5e6c8994901d776291446a2836e669b1997d7005",
          "finding_ids": [],
          "supersedes": [
            "dcm-b4-r1-spec"
          ]
        },
        {
          "id": "dcm-b4-r2-coverage",
          "performer": "claude:dcm_rev_b4_cov2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "54e5145003126d1570f2fe5cd55b825e71b6f3ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b4_cov2",
            "digest": "sha256:cc8c5791967c118548fe9c94a4e92385c0292549cf82aa604ef5206ef921ce34",
            "excerpt": "## Coverage confirming round (D-B4, repair delta 12864fc1..92da1312) — 0 findings\nF1 closed: new probe, a fallback durablefile.Native().Replace before the wrapped return: bit only at replacement_test.go:61, restored=yes.\nF2 closed: new probe, swallow every Published error: bit in the three stage rows at :31, restored=yes.\n"
          },
          "axis": "Coverage",
          "base": "c05092c603ff080970397570ac6a5c497d3663f4",
          "tip": "5e6c8994901d776291446a2836e669b1997d7005",
          "finding_ids": [],
          "supersedes": [
            "dcm-b4-r1-coverage"
          ]
        },
        {
          "id": "dcm-b4-r3-standards",
          "performer": "claude:dcm_rev_b4_std3",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "253d32d0972a9499f4d7c52f935285154c9173c6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b4_std3",
            "digest": "sha256:76ccf4f49fff9e5ec64897daa4c7b443970d65aca0fe056b1f01ab928d881fba",
            "excerpt": "## Standards second confirming round (D-B4, repair delta 8b01ad82..2477c580) — 0 findings\nS2 closed: TestHandoffReplacement ranges over durablefiletest.PostRenameStages (assert.go:14), the one in-fence copy of the stage set; the leaf test's afterRename copy is out of fence and recorded as a learning.\nAdvice: the doc comment could say \"can fail\" and Published=true; an exported func would avoid a mutable slice.\n"
          },
          "axis": "Standards",
          "base": "c05092c603ff080970397570ac6a5c497d3663f4",
          "tip": "5e49631e1c32bff364b2b48a3c1b0afc17cbf8d7",
          "finding_ids": [],
          "supersedes": [
            "dcm-b4-r2-standards"
          ]
        },
        {
          "id": "dcm-b4-r3-spec",
          "performer": "claude:dcm_rev_b4_spec3",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "253d32d0972a9499f4d7c52f935285154c9173c6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b4_spec3",
            "digest": "sha256:62bce732ad60f538c687ae8e7b17ddcdbedfc92778064bf9538878a13044c6e2",
            "excerpt": "## Spec second confirming round (D-B4, repair delta 8b01ad82..2477c580 and plan ae2055d6..8b01ad82) — 0 findings\nD21, its citation, and the behavior-change rule hold; the plan expansion adds assert.go and a floor row only; the delta stays inside the expanded ticket 05 Writes.\n"
          },
          "axis": "Spec",
          "base": "c05092c603ff080970397570ac6a5c497d3663f4",
          "tip": "5e49631e1c32bff364b2b48a3c1b0afc17cbf8d7",
          "finding_ids": [],
          "supersedes": [
            "dcm-b4-r2-spec"
          ]
        },
        {
          "id": "dcm-b4-r3-coverage",
          "performer": "claude:dcm_rev_b4_cov3",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "253d32d0972a9499f4d7c52f935285154c9173c6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b4_cov3",
            "digest": "sha256:f11025067d758ab4810603c560ef3544ef7dd8bb6274541a9671237b0dab645c",
            "excerpt": "## Coverage second confirming round (D-B4, repair delta 8b01ad82..2477c580) — 0 findings\nProbes at store.go:76-77, each restored=yes: ignore directory-sync only (bit, its row only); ignore every Published error (bit, three rows); %w -> %v (bit, four tests). No silent caller-side mutation.\n"
          },
          "axis": "Coverage",
          "base": "c05092c603ff080970397570ac6a5c497d3663f4",
          "tip": "5e49631e1c32bff364b2b48a3c1b0afc17cbf8d7",
          "finding_ids": [],
          "supersedes": [
            "dcm-b4-r2-coverage"
          ]
        }
      ]
    },
    {
      "id": "D-B5",
      "base": "2d06e7594c9d31f1ce9f22102d03da660601cea2",
      "tip": "19422ae3a3d4e5de070edbe562bbf45f201f5421",
      "plan_digest": "sha256:1328d3cecb091573d96d49938b9049c6cc47592195da6220c313c0c4dbef2abd",
      "source_digest": "c4af640671be8b9686510058a2b060a939e82f09",
      "acceptance_rows": [
        "D22",
        "D44",
        "D45",
        "D52"
      ],
      "verification": [
        {
          "id": "dcm-t06-caller-tests",
          "performer": "claude:dcm_t06",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f31faf8613dc4f4f281fc93c6c44040dd859c815",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "a6c0ffe56d68d2499f336761de2589f5988801da",
            "digest": "sha256:ed3a05260fc593191ceeb0764fc9d893081715cf541f96fc3a247b9c2529dbf6",
            "excerpt": "bench test --package ./internal/intent at a6c0ffe56d68d2499f336761de2589f5988801da (FT354-durable-callers). TestIntentReplacement (D22, one row per durablefiletest.PostRenameStages stage), TestIntentPreRenameFailure (D36 rename fault; prior ledger kept), TestTransactionPublishedFailure (D44), TestTransactionUnclassifiedFailureCompensates (unclassified branch), TestReauthorizePublishedFailure (D52), and the unchanged D45 TestLedgerTransactionTerminalWriteFailureKeepsThePreviousBytes pass.\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,pass,999,72\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "caller-tests",
          "command": "bench test --package ./internal/intent",
          "exit_code": 0
        },
        {
          "id": "dcm-t06-existing-canary-owner",
          "performer": "claude:dcm_t06",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f31faf8613dc4f4f281fc93c6c44040dd859c815",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "a6c0ffe56d68d2499f336761de2589f5988801da",
            "digest": "sha256:dc64111165aae6395ed9de511e6ae2b09bb76d5aef29616be85a3885286c8bca",
            "excerpt": "bench test --check bounds-policy at a6c0ffe56d68d2499f336761de2589f5988801da (FT354-durable-callers). The canary tests/canary/package-core-guard/bounds-intent-window-fixed is unchanged: its BASE pins paths only, and intent.go keeps bounds.VerdictWindow(bounds.IntentLockTimeout).\ncheck[1]{name,kind,tests_run,subjects}:\n  bounds-policy,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,250,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "existing-canary-owner",
          "command": "bench test --check bounds-policy",
          "exit_code": 0
        },
        {
          "id": "dcm-t06-root-conformance",
          "performer": "claude:dcm_t06",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f31faf8613dc4f4f281fc93c6c44040dd859c815",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "a6c0ffe56d68d2499f336761de2589f5988801da",
            "digest": "sha256:3b22fe21353e5e8137e9824e4201d8cf2e06b3587cdd407f5f67611c2e31e949",
            "excerpt": "bench test --check root-conformance at a6c0ffe56d68d2499f336761de2589f5988801da (FT354-durable-callers)\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16553,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "dcm-t06-omission-proof",
          "performer": "claude:dcm_t06",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f31faf8613dc4f4f281fc93c6c44040dd859c815",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "82697755375c32932a82993f9f42432c744a7b66",
            "digest": "sha256:767ae134625b6455041f2a65d9fce31b798a7ad15fa6f5782de8823855831186",
            "excerpt": "D-B5 omission-proof probes at 82697755375c32932a82993f9f42432c744a7b66 (FT354-durable-callers; a6c0ffe5 adds only spec text), bench test --package ./internal/intent, no --run filter.\nProbe q09 (compensate every write error). internal/intent/transaction.go swap old (exact): \"if compensate != nil && !published(err) {\" new (exact): \"if compensate != nil {\"\nResult: verdict bit, 2 failed tests (TestTransactionPublishedFailure compensation count; TestReauthorizePublishedFailure rollback count), restored yes.\n\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/transaction.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,72\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,501\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,1362,72\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestReauthorizePublishedFailure,\"replacement_test.go:138: the rollback ran 1 times after a completed rename, want 0\",1\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure,\"replacement_test.go:78: compensation ran 1 times after a completed rename, want 0\",1\nskips[0]{package,test,reason}:\n\nProbe q13 (omit the reauthorization effect). internal/intent/assignment.go swap old (exact): \"step, err := transition(stored, next)\" new (exact): \"step, err := transition(stored, stored)\"\nResult: verdict bit, 1 failed test (TestReauthorizePublishedFailure at the external-witness check: effect ccc... against ledger request 80f3b5dd...), restored yes.\n\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/assignment.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,72\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,461\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,983,72\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestReauthorizePublishedFailure,\"replacement_test.go:150: external effect = \\\"cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc\\\", want the ledger request \\\"80f3b5ddef868a2b4422c3203ef39bec9a3a4db048441b8c968ad92acfd966ab\\\"\",1\nskips[0]{package,test,reason}:\n\nProbe q14 (roll back a published ledger without compensation). internal/intent/transaction.go swap old (exact): \"if err := writePathWith(replacer, path, next); err != nil {\" new (exact): \"old, _ := os.ReadFile(path); if err := writePathWith(replacer, path, next); err != nil { if published(err) { _ = durablefile.Replace(path, old, 0o600) };\"\nResult: verdict bit, 2 failed tests (TestTransactionPublishedFailure and TestReauthorizePublishedFailure at the new-ledger bytes check), restored yes. The package passes 72 tests after restore.\n\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/transaction.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,72\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1067\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,949,72\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestReauthorizePublishedFailure,\"replacement_test.go:140: ledger = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":2,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"seed\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"shift\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:01Z\\\\\\\\\\\"}],\\\\\\\\\\\"assignments\\\\\\\\\\\":[{\\\\\\\\\\\"schema\\\\\\\\\\\":\\\\\\\\\\\"bench-assignment/v1\\\\\\\\\\\",\\\\\\\\\\\"id\\\\\\\\\\\":\\\\\\\\\\\"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\\\\\\\\\\\",\\\\\\\\\\\"… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure,\"replacement_test.go:80: ledger = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":2,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"seed\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"shift\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:01Z\\\\\\\\\\\"}],\\\\\\\\\\\"assignments\\\\\\\\\\\":[{\\\\\\\\\\\"schema\\\\\\\\\\\":\\\\\\\\\\\"bench-assignment/v1\\\\\\\\\\\",\\\\\\\\\\\"id\\\\\\\\\\\":\\\\\\\\\\\"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\\\\\\\\\\\",\\\\\\\\\\\"o… (1020 bytes)\",2\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "omission-proof",
          "command": "bench test --package ./internal/intent",
          "exit_code": 0,
          "probe": {
            "mutation": "Compensate every write error or omit the reauthorization effect check. The new-ledger and external-witness test must detect their mismatch.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "82697755375c32932a82993f9f42432c744a7b66",
              "digest": "sha256:767ae134625b6455041f2a65d9fce31b798a7ad15fa6f5782de8823855831186",
              "excerpt": "D-B5 omission-proof probes at 82697755375c32932a82993f9f42432c744a7b66 (FT354-durable-callers; a6c0ffe5 adds only spec text), bench test --package ./internal/intent, no --run filter.\nProbe q09 (compensate every write error). internal/intent/transaction.go swap old (exact): \"if compensate != nil && !published(err) {\" new (exact): \"if compensate != nil {\"\nResult: verdict bit, 2 failed tests (TestTransactionPublishedFailure compensation count; TestReauthorizePublishedFailure rollback count), restored yes.\n\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/transaction.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,72\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,501\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,1362,72\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestReauthorizePublishedFailure,\"replacement_test.go:138: the rollback ran 1 times after a completed rename, want 0\",1\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure,\"replacement_test.go:78: compensation ran 1 times after a completed rename, want 0\",1\nskips[0]{package,test,reason}:\n\nProbe q13 (omit the reauthorization effect). internal/intent/assignment.go swap old (exact): \"step, err := transition(stored, next)\" new (exact): \"step, err := transition(stored, stored)\"\nResult: verdict bit, 1 failed test (TestReauthorizePublishedFailure at the external-witness check: effect ccc... against ledger request 80f3b5dd...), restored yes.\n\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/assignment.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,72\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,461\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,983,72\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestReauthorizePublishedFailure,\"replacement_test.go:150: external effect = \\\"cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc\\\", want the ledger request \\\"80f3b5ddef868a2b4422c3203ef39bec9a3a4db048441b8c968ad92acfd966ab\\\"\",1\nskips[0]{package,test,reason}:\n\nProbe q14 (roll back a published ledger without compensation). internal/intent/transaction.go swap old (exact): \"if err := writePathWith(replacer, path, next); err != nil {\" new (exact): \"old, _ := os.ReadFile(path); if err := writePathWith(replacer, path, next); err != nil { if published(err) { _ = durablefile.Replace(path, old, 0o600) };\"\nResult: verdict bit, 2 failed tests (TestTransactionPublishedFailure and TestReauthorizePublishedFailure at the new-ledger bytes check), restored yes. The package passes 72 tests after restore.\n\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/transaction.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,72\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1067\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,949,72\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestReauthorizePublishedFailure,\"replacement_test.go:140: ledger = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":2,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"seed\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"shift\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:01Z\\\\\\\\\\\"}],\\\\\\\\\\\"assignments\\\\\\\\\\\":[{\\\\\\\\\\\"schema\\\\\\\\\\\":\\\\\\\\\\\"bench-assignment/v1\\\\\\\\\\\",\\\\\\\\\\\"id\\\\\\\\\\\":\\\\\\\\\\\"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\\\\\\\\\\\",\\\\\\\\\\\"… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure,\"replacement_test.go:80: ledger = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":2,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"seed\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"shift\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:01Z\\\\\\\\\\\"}],\\\\\\\\\\\"assignments\\\\\\\\\\\":[{\\\\\\\\\\\"schema\\\\\\\\\\\":\\\\\\\\\\\"bench-assignment/v1\\\\\\\\\\\",\\\\\\\\\\\"id\\\\\\\\\\\":\\\\\\\\\\\"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\\\\\\\\\\\",\\\\\\\\\\\"o… (1020 bytes)\",2\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dcm-t06-preservation",
          "performer": "claude:dcm_t06",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f31faf8613dc4f4f281fc93c6c44040dd859c815",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "82697755375c32932a82993f9f42432c744a7b66",
            "digest": "sha256:1ddb828f035d67c3675971734bf8d986caa908dcc9c4155f1a85a499f965238c",
            "excerpt": "D-B5 preservation probes at 82697755375c32932a82993f9f42432c744a7b66 (FT354-durable-callers; a6c0ffe5 adds only spec text), bench test --package ./internal/intent, no --run filter. TestIntentReplacement owns the mode and bytes expectations; the other new tests check bytes only. Restored yes for each; the package passes 72 tests after restore.\nProbe q04 (mode literal). internal/intent/intent.go swap old (exact): \"return replacer.Replace(path, append(data, '\\n'), 0o600)\" new (exact): \"return replacer.Replace(path, append(data, '\\n'), 0o640)\"\nResult: verdict bit, 3 failed tests (the three TestIntentReplacement stage rows at the mode check only).\n\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/intent.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,72\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,937\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,1322,72\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestIntentReplacement/directory-close,\"replacement_test.go:42: /dev/shm/OYCMYV/t/TestIntentReplacementdirectory-close1076546497/001/.git/bench-intent.json mode = -rw-r-----, want a regular file at -rw-------\",1\n  github.com/gibbonmi/bench/internal/intent,TestIntentReplacement/directory-open,\"replacement_test.go:42: /dev/shm/OYCMYV/t/TestIntentReplacementdirectory-open3840080956/001/.git/bench-intent.json mode = -rw-r-----, want a regular file at -rw-------\",1\n  github.com/gibbonmi/bench/internal/intent,TestIntentReplacement/directory-sync,\"replacement_test.go:42: /dev/shm/OYCMYV/t/TestIntentReplacementdirectory-sync2017259731/001/.git/bench-intent.json mode = -rw-r-----, want a regular file at -rw-------\",1\nskips[0]{package,test,reason}:\n\nProbe q05 (payload literal). internal/intent/intent.go --omit (exact): \"ledger.Schema = Schema\"\nResult: verdict bit, the three TestIntentReplacement rows at the bytes check plus the existing TestPurgeAssignmentsDropsEveryLegacySchemaRecord.\n\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/intent.go,omit,failed,4,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,72\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1784\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,1347,72\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestIntentReplacement/directory-close,\"replacement_test.go:42: /dev/shm/AOLA2L/t/TestIntentReplacementdirectory-close2908410681/001/.git/bench-intent.json = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":0,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"added\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"worktree\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:02Z\\\\\\\\\\\"},{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"s… (509 bytes)\",1\n  github.com/gibbonmi/bench/internal/intent,TestIntentReplacement/directory-open,\"replacement_test.go:42: /dev/shm/AOLA2L/t/TestIntentReplacementdirectory-open3620186243/001/.git/bench-intent.json = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":0,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"added\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"worktree\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:02Z\\\\\\\\\\\"},{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"se… (508 bytes)\",1\n  github.com/gibbonmi/bench/internal/intent,TestIntentReplacement/directory-sync,\"replacement_test.go:42: /dev/shm/AOLA2L/t/TestIntentReplacementdirectory-sync281249581/001/.git/bench-intent.json = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":0,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"added\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"worktree\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:02Z\\\\\\\\\\\"},{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"see… (507 bytes)\",1\n  github.com/gibbonmi/bench/internal/intent,TestPurgeAssignmentsDropsEveryLegacySchemaRecord,\"purge_test.go:86: ledger after legacy purge = ledger.Ledger{Schema:1, Entries:[]ledger.Entry{}, Assignments:[]ledger.Assignment{}, CleanupReceipts:[]ledger.CleanupReceipt(nil), CommitmentReceipts:[]ledger.CommitmentReceipt(nil), MilestoneRe… (328 bytes)\",1\nskips[0]{package,test,reason}:\n\nProbe q07 (pre-rename refusal: fallback write of the new bytes after an unpublished error). internal/intent/intent.go swap old (exact): \"return replacer.Replace(path, append(data, '\\n'), 0o600)\" new (exact): \"if err := replacer.Replace(path, append(data, '\\n'), 0o600); err != nil { var e *durablefile.Error; if errors.As(err, &e) && !e.Published { _ = durablefile.Replace(path, append(data, '\\n'), 0o600) }; return err }; return nil\"\nResult: verdict bit, 1 failed test (TestIntentPreRenameFailure at the prior-ledger check only).\n\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/intent.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,72\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,637\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,1489,72\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestIntentPreRenameFailure,\"replacement_test.go:60: ledger = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":2,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"added\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"worktree\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:02Z\\\\\\\\\\\"}]}\\\\\\\\n\\\", want \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":2,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"seed\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"shift\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970… (748 bytes)\",1\nskips[0]{package,test,reason}:\n\nProbe q17 (unclassified failure overwrites the ledger). internal/intent/intent.go swap old (exact): \"return fmt.Errorf(\\\"encode intent ledger: %w\\\", err)\" new (exact): \"_ = durablefile.Replace(path, []byte(\\\"{}\\n\\\"), 0o600); return fmt.Errorf(\\\"encode intent ledger: %w\\\", err)\"\nResult: verdict bit, 1 failed test (TestTransactionUnclassifiedFailureCompensates at the prior-ledger check only).\n\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/intent.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,72\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,632\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,1001,72\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestTransactionUnclassifiedFailureCompensates,\"replacement_test.go:108: ledger = \\\"{}\\\\\\\\n\\\", want \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":2,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"seed\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"shift\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:01Z\\\\\\\\\\\"}],\\\\\\\\\\\"assignments\\\\\\\\\\\":[{\\\\\\\\\\\"schema\\\\\\\\\\\":\\\\\\\\\\\"bench-assignment/v1\\\\\\\\\\\",\\\\\\\\\\\"id\\\\\\\\\\\":\\\\\\\\\\\"bbbbbbbbbbbbbbbbbbbbbbbb… (641 bytes)\",1\nskips[0]{package,test,reason}:\n\nEvery other check and its red (q01-q18, with exact swaps) is in the author's probe log; q11 reds the D44 lock check only, q15 the unclassified-error check only, q01-q03 one ignored-error red per PostRenameStages row.\n"
          },
          "requirement": "preservation",
          "command": "bench test --package ./internal/intent",
          "exit_code": 0,
          "probe": {
            "mutation": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "82697755375c32932a82993f9f42432c744a7b66",
              "digest": "sha256:1ddb828f035d67c3675971734bf8d986caa908dcc9c4155f1a85a499f965238c",
              "excerpt": "D-B5 preservation probes at 82697755375c32932a82993f9f42432c744a7b66 (FT354-durable-callers; a6c0ffe5 adds only spec text), bench test --package ./internal/intent, no --run filter. TestIntentReplacement owns the mode and bytes expectations; the other new tests check bytes only. Restored yes for each; the package passes 72 tests after restore.\nProbe q04 (mode literal). internal/intent/intent.go swap old (exact): \"return replacer.Replace(path, append(data, '\\n'), 0o600)\" new (exact): \"return replacer.Replace(path, append(data, '\\n'), 0o640)\"\nResult: verdict bit, 3 failed tests (the three TestIntentReplacement stage rows at the mode check only).\n\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/intent.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,72\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,937\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,1322,72\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestIntentReplacement/directory-close,\"replacement_test.go:42: /dev/shm/OYCMYV/t/TestIntentReplacementdirectory-close1076546497/001/.git/bench-intent.json mode = -rw-r-----, want a regular file at -rw-------\",1\n  github.com/gibbonmi/bench/internal/intent,TestIntentReplacement/directory-open,\"replacement_test.go:42: /dev/shm/OYCMYV/t/TestIntentReplacementdirectory-open3840080956/001/.git/bench-intent.json mode = -rw-r-----, want a regular file at -rw-------\",1\n  github.com/gibbonmi/bench/internal/intent,TestIntentReplacement/directory-sync,\"replacement_test.go:42: /dev/shm/OYCMYV/t/TestIntentReplacementdirectory-sync2017259731/001/.git/bench-intent.json mode = -rw-r-----, want a regular file at -rw-------\",1\nskips[0]{package,test,reason}:\n\nProbe q05 (payload literal). internal/intent/intent.go --omit (exact): \"ledger.Schema = Schema\"\nResult: verdict bit, the three TestIntentReplacement rows at the bytes check plus the existing TestPurgeAssignmentsDropsEveryLegacySchemaRecord.\n\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/intent.go,omit,failed,4,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,72\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1784\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,1347,72\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestIntentReplacement/directory-close,\"replacement_test.go:42: /dev/shm/AOLA2L/t/TestIntentReplacementdirectory-close2908410681/001/.git/bench-intent.json = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":0,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"added\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"worktree\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:02Z\\\\\\\\\\\"},{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"s… (509 bytes)\",1\n  github.com/gibbonmi/bench/internal/intent,TestIntentReplacement/directory-open,\"replacement_test.go:42: /dev/shm/AOLA2L/t/TestIntentReplacementdirectory-open3620186243/001/.git/bench-intent.json = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":0,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"added\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"worktree\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:02Z\\\\\\\\\\\"},{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"se… (508 bytes)\",1\n  github.com/gibbonmi/bench/internal/intent,TestIntentReplacement/directory-sync,\"replacement_test.go:42: /dev/shm/AOLA2L/t/TestIntentReplacementdirectory-sync281249581/001/.git/bench-intent.json = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":0,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"added\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"worktree\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:02Z\\\\\\\\\\\"},{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"see… (507 bytes)\",1\n  github.com/gibbonmi/bench/internal/intent,TestPurgeAssignmentsDropsEveryLegacySchemaRecord,\"purge_test.go:86: ledger after legacy purge = ledger.Ledger{Schema:1, Entries:[]ledger.Entry{}, Assignments:[]ledger.Assignment{}, CleanupReceipts:[]ledger.CleanupReceipt(nil), CommitmentReceipts:[]ledger.CommitmentReceipt(nil), MilestoneRe… (328 bytes)\",1\nskips[0]{package,test,reason}:\n\nProbe q07 (pre-rename refusal: fallback write of the new bytes after an unpublished error). internal/intent/intent.go swap old (exact): \"return replacer.Replace(path, append(data, '\\n'), 0o600)\" new (exact): \"if err := replacer.Replace(path, append(data, '\\n'), 0o600); err != nil { var e *durablefile.Error; if errors.As(err, &e) && !e.Published { _ = durablefile.Replace(path, append(data, '\\n'), 0o600) }; return err }; return nil\"\nResult: verdict bit, 1 failed test (TestIntentPreRenameFailure at the prior-ledger check only).\n\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/intent.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,72\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,637\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,1489,72\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestIntentPreRenameFailure,\"replacement_test.go:60: ledger = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":2,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"added\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"worktree\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:02Z\\\\\\\\\\\"}]}\\\\\\\\n\\\", want \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":2,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"seed\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"shift\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970… (748 bytes)\",1\nskips[0]{package,test,reason}:\n\nProbe q17 (unclassified failure overwrites the ledger). internal/intent/intent.go swap old (exact): \"return fmt.Errorf(\\\"encode intent ledger: %w\\\", err)\" new (exact): \"_ = durablefile.Replace(path, []byte(\\\"{}\\n\\\"), 0o600); return fmt.Errorf(\\\"encode intent ledger: %w\\\", err)\"\nResult: verdict bit, 1 failed test (TestTransactionUnclassifiedFailureCompensates at the prior-ledger check only).\n\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/intent.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,72\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,632\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,1001,72\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestTransactionUnclassifiedFailureCompensates,\"replacement_test.go:108: ledger = \\\"{}\\\\\\\\n\\\", want \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":2,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"seed\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"shift\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:01Z\\\\\\\\\\\"}],\\\\\\\\\\\"assignments\\\\\\\\\\\":[{\\\\\\\\\\\"schema\\\\\\\\\\\":\\\\\\\\\\\"bench-assignment/v1\\\\\\\\\\\",\\\\\\\\\\\"id\\\\\\\\\\\":\\\\\\\\\\\"bbbbbbbbbbbbbbbbbbbbbbbb… (641 bytes)\",1\nskips[0]{package,test,reason}:\n\nEvery other check and its red (q01-q18, with exact swaps) is in the author's probe log; q11 reds the D44 lock check only, q15 the unclassified-error check only, q01-q03 one ignored-error red per PostRenameStages row.\n"
            }
          }
        },
        {
          "id": "dcm-t06-rep1-caller-tests",
          "performer": "claude:dcm_t06_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "c4af640671be8b9686510058a2b060a939e82f09",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "19422ae3a3d4e5de070edbe562bbf45f201f5421",
            "digest": "sha256:a8e4f6b3ad2f5a4772c32330acd5de8e3fb0cc6a6d86bf64db5d04f0d27305dc",
            "excerpt": "bench test --package ./internal/intent at 19422ae3a3d4e5de070edbe562bbf45f201f5421 (FT354-durable-callers). TestIntentReplacement (D22, one row per durablefiletest.PostRenameStages stage), TestIntentPreRenameFailure (D36 rename fault; prior ledger kept), TestTransactionPublishedFailure (D44 and D45 through TransactWithReplacer: a rename row with one compensation under the held lock and the prior ledger, and one row per PostRenameStages stage with no compensation and the new ledger; the lock is gone after every row), TestTransactionUnclassifiedFailureCompensates (unclassified branch), TestReauthorizePublishedFailure (D52), and the unchanged D45 TestLedgerTransactionTerminalWriteFailureKeepsThePreviousBytes pass.\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,pass,1005,76\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "caller-tests",
          "command": "bench test --package ./internal/intent",
          "exit_code": 0
        },
        {
          "id": "dcm-t06-rep1-existing-canary-owner",
          "performer": "claude:dcm_t06_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "c4af640671be8b9686510058a2b060a939e82f09",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "19422ae3a3d4e5de070edbe562bbf45f201f5421",
            "digest": "sha256:f69306f5e97a04ec00de4b4a6c26e8caa3bf7e5589902b54e9a1c29da6008a4f",
            "excerpt": "bench test --check bounds-policy at 19422ae3a3d4e5de070edbe562bbf45f201f5421 (FT354-durable-callers). The canary tests/canary/package-core-guard/bounds-intent-window-fixed is unchanged: its BASE pins paths only, and intent.go keeps bounds.VerdictWindow(bounds.IntentLockTimeout).\ncheck[1]{name,kind,tests_run,subjects}:\n  bounds-policy,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,223,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "existing-canary-owner",
          "command": "bench test --check bounds-policy",
          "exit_code": 0
        },
        {
          "id": "dcm-t06-rep1-root-conformance",
          "performer": "claude:dcm_t06_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "c4af640671be8b9686510058a2b060a939e82f09",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "19422ae3a3d4e5de070edbe562bbf45f201f5421",
            "digest": "sha256:f10dbdf97b28d25f313528923b1b3b8e144a710a8189c45766708de850c69862",
            "excerpt": "bench test --check root-conformance at 19422ae3a3d4e5de070edbe562bbf45f201f5421 (FT354-durable-callers)\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15435,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "dcm-t06-rep1-omission-proof",
          "performer": "claude:dcm_t06_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "c4af640671be8b9686510058a2b060a939e82f09",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "19422ae3a3d4e5de070edbe562bbf45f201f5421",
            "digest": "sha256:0899da93c64b86e7f476da2cbd7b12c1711c34767326da295ac99d37c81ee3fb",
            "excerpt": "D-B5 omission-proof. D-B5 repair 1 probes at 19422ae3a3d4e5de070edbe562bbf45f201f5421 (FT354-durable-callers), bench probe <file> ... --package ./internal/intent, no --run filter. The package passes 76 tests at the tip after every restore. At base c73da9c875d1ffe9e8cf7f7b8c8730cf14af5a0e, before the repair, the C1c and C2 swaps were silent (72 tests passed, restored yes); the repair turned each red.\nProbe q09 (compensate every write error). internal/intent/transaction.go swap old (exact): \"if compensate != nil && !published(err) {\" new (exact): \"if compensate != nil {\"\nResult: verdict bit, 4 failed: the three TestTransactionPublishedFailure post-rename rows and TestReauthorizePublishedFailure at the compensation or rollback count, restored yes.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/transaction.go,swap,failed,4,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,76\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,873\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,981,76\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestReauthorizePublishedFailure,\"replacement_test.go:168: the rollback ran 1 times after a completed rename, want 0\",1\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/directory-close,\"replacement_test.go:100: compensation ran 1 times after a completed rename, want 0\",1\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/directory-open,\"replacement_test.go:100: compensation ran 1 times after a completed rename, want 0\",1\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/directory-sync,\"replacement_test.go:100: compensation ran 1 times after a completed rename, want 0\",1\nskips[0]{package,test,reason}:\n\nProbe q13 (omit the reauthorization effect). internal/intent/assignment.go swap old (exact): \"step, err := transition(stored, next)\" new (exact): \"step, err := transition(stored, stored)\"\nResult: verdict bit, 1 failed: TestReauthorizePublishedFailure at the external-witness check, restored yes.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/assignment.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,76\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,461\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,981,76\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestReauthorizePublishedFailure,\"replacement_test.go:180: external effect = \\\"cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc\\\", want the ledger request \\\"80f3b5ddef868a2b4422c3203ef39bec9a3a4db048441b8c968ad92acfd966ab\\\"\",1\nskips[0]{package,test,reason}:\n\nProbe q14 (roll back a published ledger without compensation). internal/intent/transaction.go swap old (exact): \"if err := writePathWith(replacer, path, next); err != nil {\" new (exact): \"old, _ := os.ReadFile(path); if err := writePathWith(replacer, path, next); err != nil { if published(err) { _ = durablefile.Replace(path, old, 0o600) };\"\nResult: verdict bit, 4 failed: the three TestTransactionPublishedFailure post-rename rows and TestReauthorizePublishedFailure at the new-ledger bytes check, restored yes.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/transaction.go,swap,failed,4,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,76\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,2005\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,978,76\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestReauthorizePublishedFailure,\"replacement_test.go:170: ledger = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":2,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"seed\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"shift\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:01Z\\\\\\\\\\\"}],\\\\\\\\\\\"assignments\\\\\\\\\\\":[{\\\\\\\\\\\"schema\\\\\\\\\\\":\\\\\\\\\\\"bench-assignment/v1\\\\\\\\\\\",\\\\\\\\\\\"id\\\\\\\\\\\":\\\\\\\\\\\"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\\\\\\\\\\\",\\\\\\\\\\\"… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/directory-close,\"replacement_test.go:102: ledger = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":2,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"seed\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"shift\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:01Z\\\\\\\\\\\"}],\\\\\\\\\\\"assignments\\\\\\\\\\\":[{\\\\\\\\\\\"schema\\\\\\\\\\\":\\\\\\\\\\\"bench-assignment/v1\\\\\\\\\\\",\\\\\\\\\\\"id\\\\\\\\\\\":\\\\\\\\\\\"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\\\\\\\\\\\",\\\\\\\\\\\"… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/directory-open,\"replacement_test.go:102: ledger = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":2,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"seed\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"shift\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:01Z\\\\\\\\\\\"}],\\\\\\\\\\\"assignments\\\\\\\\\\\":[{\\\\\\\\\\\"schema\\\\\\\\\\\":\\\\\\\\\\\"bench-assignment/v1\\\\\\\\\\\",\\\\\\\\\\\"id\\\\\\\\\\\":\\\\\\\\\\\"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\\\\\\\\\\\",\\\\\\\\\\\"… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/directory-sync,\"replacement_test.go:102: ledger = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":2,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"seed\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"shift\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:01Z\\\\\\\\\\\"}],\\\\\\\\\\\"assignments\\\\\\\\\\\":[{\\\\\\\\\\\"schema\\\\\\\\\\\":\\\\\\\\\\\"bench-assignment/v1\\\\\\\\\\\",\\\\\\\\\\\"id\\\\\\\\\\\":\\\\\\\\\\\"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\\\\\\\\\\\",\\\\\\\\\\\"… (1020 bytes)\",2\nskips[0]{package,test,reason}:\n\nProbe C1a (every stage but create counts as published). internal/intent/transaction.go swap old (exact): \"return errors.As(err, &replaceErr) && replaceErr.Published\" new (exact): \"return errors.As(err, &replaceErr) && replaceErr.Stage != \"create\"\"\nResult: verdict bit, 1 failed: TestTransactionPublishedFailure/rename at the compensation count (0, want 1), restored yes.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/transaction.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,76\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,346\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,977,76\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/rename,\"replacement_test.go:105: compensation ran 0 times before a completed rename, want 1\",1\nskips[0]{package,test,reason}:\n\nProbe C1b (a rename failure counts as published). internal/intent/transaction.go swap old (exact): \"return errors.As(err, &replaceErr) && replaceErr.Published\" new (exact): \"return errors.As(err, &replaceErr) && (replaceErr.Published || replaceErr.Stage == \"rename\")\"\nResult: verdict bit, 1 failed: TestTransactionPublishedFailure/rename at the compensation count (0, want 1), restored yes.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/transaction.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,76\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,346\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,967,76\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/rename,\"replacement_test.go:105: compensation ran 0 times before a completed rename, want 1\",1\nskips[0]{package,test,reason}:\n\nProbe C1c (only a directory-sync failure counts as published). internal/intent/transaction.go swap old (exact): \"return errors.As(err, &replaceErr) && replaceErr.Published\" new (exact): \"return errors.As(err, &replaceErr) && replaceErr.Stage == \"directory-sync\"\"\nResult: verdict bit, 2 failed: TestTransactionPublishedFailure/directory-open and /directory-close at the compensation count (1, want 0), restored yes.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/transaction.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,76\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,532\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,997,76\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/directory-close,\"replacement_test.go:100: compensation ran 1 times after a completed rename, want 0\",1\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/directory-open,\"replacement_test.go:100: compensation ran 1 times after a completed rename, want 0\",1\nskips[0]{package,test,reason}:\n\nProbe C2 (release the lock before the compensation; the deferred release stays). internal/intent/transaction.go swap old (exact): \"compensate()\" new (exact): \"release(); compensate()\"\nResult: verdict bit, 1 failed: TestTransactionPublishedFailure/rename at the lock check inside the compensation, restored yes.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/transaction.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,76\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,451\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,988,76\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/rename,\"replacement_test.go:93: lock file during the compensation: stat /dev/shm/CORSYI/t/TestTransactionPublishedFailurerename2407171961/001/.git/bench-intent.json.lock: no such file or directory\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "omission-proof",
          "command": "bench test --package ./internal/intent",
          "exit_code": 0,
          "probe": {
            "mutation": "Compensate every write error or omit the reauthorization effect check. The new-ledger and external-witness test must detect their mismatch.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "19422ae3a3d4e5de070edbe562bbf45f201f5421",
              "digest": "sha256:0899da93c64b86e7f476da2cbd7b12c1711c34767326da295ac99d37c81ee3fb",
              "excerpt": "D-B5 omission-proof. D-B5 repair 1 probes at 19422ae3a3d4e5de070edbe562bbf45f201f5421 (FT354-durable-callers), bench probe <file> ... --package ./internal/intent, no --run filter. The package passes 76 tests at the tip after every restore. At base c73da9c875d1ffe9e8cf7f7b8c8730cf14af5a0e, before the repair, the C1c and C2 swaps were silent (72 tests passed, restored yes); the repair turned each red.\nProbe q09 (compensate every write error). internal/intent/transaction.go swap old (exact): \"if compensate != nil && !published(err) {\" new (exact): \"if compensate != nil {\"\nResult: verdict bit, 4 failed: the three TestTransactionPublishedFailure post-rename rows and TestReauthorizePublishedFailure at the compensation or rollback count, restored yes.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/transaction.go,swap,failed,4,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,76\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,873\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,981,76\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestReauthorizePublishedFailure,\"replacement_test.go:168: the rollback ran 1 times after a completed rename, want 0\",1\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/directory-close,\"replacement_test.go:100: compensation ran 1 times after a completed rename, want 0\",1\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/directory-open,\"replacement_test.go:100: compensation ran 1 times after a completed rename, want 0\",1\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/directory-sync,\"replacement_test.go:100: compensation ran 1 times after a completed rename, want 0\",1\nskips[0]{package,test,reason}:\n\nProbe q13 (omit the reauthorization effect). internal/intent/assignment.go swap old (exact): \"step, err := transition(stored, next)\" new (exact): \"step, err := transition(stored, stored)\"\nResult: verdict bit, 1 failed: TestReauthorizePublishedFailure at the external-witness check, restored yes.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/assignment.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,76\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,461\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,981,76\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestReauthorizePublishedFailure,\"replacement_test.go:180: external effect = \\\"cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc\\\", want the ledger request \\\"80f3b5ddef868a2b4422c3203ef39bec9a3a4db048441b8c968ad92acfd966ab\\\"\",1\nskips[0]{package,test,reason}:\n\nProbe q14 (roll back a published ledger without compensation). internal/intent/transaction.go swap old (exact): \"if err := writePathWith(replacer, path, next); err != nil {\" new (exact): \"old, _ := os.ReadFile(path); if err := writePathWith(replacer, path, next); err != nil { if published(err) { _ = durablefile.Replace(path, old, 0o600) };\"\nResult: verdict bit, 4 failed: the three TestTransactionPublishedFailure post-rename rows and TestReauthorizePublishedFailure at the new-ledger bytes check, restored yes.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/transaction.go,swap,failed,4,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,76\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,2005\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,978,76\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestReauthorizePublishedFailure,\"replacement_test.go:170: ledger = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":2,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"seed\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"shift\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:01Z\\\\\\\\\\\"}],\\\\\\\\\\\"assignments\\\\\\\\\\\":[{\\\\\\\\\\\"schema\\\\\\\\\\\":\\\\\\\\\\\"bench-assignment/v1\\\\\\\\\\\",\\\\\\\\\\\"id\\\\\\\\\\\":\\\\\\\\\\\"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\\\\\\\\\\\",\\\\\\\\\\\"… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/directory-close,\"replacement_test.go:102: ledger = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":2,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"seed\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"shift\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:01Z\\\\\\\\\\\"}],\\\\\\\\\\\"assignments\\\\\\\\\\\":[{\\\\\\\\\\\"schema\\\\\\\\\\\":\\\\\\\\\\\"bench-assignment/v1\\\\\\\\\\\",\\\\\\\\\\\"id\\\\\\\\\\\":\\\\\\\\\\\"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\\\\\\\\\\\",\\\\\\\\\\\"… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/directory-open,\"replacement_test.go:102: ledger = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":2,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"seed\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"shift\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:01Z\\\\\\\\\\\"}],\\\\\\\\\\\"assignments\\\\\\\\\\\":[{\\\\\\\\\\\"schema\\\\\\\\\\\":\\\\\\\\\\\"bench-assignment/v1\\\\\\\\\\\",\\\\\\\\\\\"id\\\\\\\\\\\":\\\\\\\\\\\"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\\\\\\\\\\\",\\\\\\\\\\\"… (1020 bytes)\",2\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/directory-sync,\"replacement_test.go:102: ledger = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":2,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"seed\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"shift\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:01Z\\\\\\\\\\\"}],\\\\\\\\\\\"assignments\\\\\\\\\\\":[{\\\\\\\\\\\"schema\\\\\\\\\\\":\\\\\\\\\\\"bench-assignment/v1\\\\\\\\\\\",\\\\\\\\\\\"id\\\\\\\\\\\":\\\\\\\\\\\"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\\\\\\\\\\\",\\\\\\\\\\\"… (1020 bytes)\",2\nskips[0]{package,test,reason}:\n\nProbe C1a (every stage but create counts as published). internal/intent/transaction.go swap old (exact): \"return errors.As(err, &replaceErr) && replaceErr.Published\" new (exact): \"return errors.As(err, &replaceErr) && replaceErr.Stage != \"create\"\"\nResult: verdict bit, 1 failed: TestTransactionPublishedFailure/rename at the compensation count (0, want 1), restored yes.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/transaction.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,76\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,346\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,977,76\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/rename,\"replacement_test.go:105: compensation ran 0 times before a completed rename, want 1\",1\nskips[0]{package,test,reason}:\n\nProbe C1b (a rename failure counts as published). internal/intent/transaction.go swap old (exact): \"return errors.As(err, &replaceErr) && replaceErr.Published\" new (exact): \"return errors.As(err, &replaceErr) && (replaceErr.Published || replaceErr.Stage == \"rename\")\"\nResult: verdict bit, 1 failed: TestTransactionPublishedFailure/rename at the compensation count (0, want 1), restored yes.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/transaction.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,76\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,346\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,967,76\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/rename,\"replacement_test.go:105: compensation ran 0 times before a completed rename, want 1\",1\nskips[0]{package,test,reason}:\n\nProbe C1c (only a directory-sync failure counts as published). internal/intent/transaction.go swap old (exact): \"return errors.As(err, &replaceErr) && replaceErr.Published\" new (exact): \"return errors.As(err, &replaceErr) && replaceErr.Stage == \"directory-sync\"\"\nResult: verdict bit, 2 failed: TestTransactionPublishedFailure/directory-open and /directory-close at the compensation count (1, want 0), restored yes.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/transaction.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,76\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,532\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,997,76\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/directory-close,\"replacement_test.go:100: compensation ran 1 times after a completed rename, want 0\",1\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/directory-open,\"replacement_test.go:100: compensation ran 1 times after a completed rename, want 0\",1\nskips[0]{package,test,reason}:\n\nProbe C2 (release the lock before the compensation; the deferred release stays). internal/intent/transaction.go swap old (exact): \"compensate()\" new (exact): \"release(); compensate()\"\nResult: verdict bit, 1 failed: TestTransactionPublishedFailure/rename at the lock check inside the compensation, restored yes.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/transaction.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,76\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,451\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,988,76\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/rename,\"replacement_test.go:93: lock file during the compensation: stat /dev/shm/CORSYI/t/TestTransactionPublishedFailurerename2407171961/001/.git/bench-intent.json.lock: no such file or directory\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dcm-t06-rep1-preservation",
          "performer": "claude:dcm_t06_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "c4af640671be8b9686510058a2b060a939e82f09",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "19422ae3a3d4e5de070edbe562bbf45f201f5421",
            "digest": "sha256:87b246bc2c1a5150c67d01572e31066e6b6308706cbcdb4bd6b45403386a06fb",
            "excerpt": "D-B5 preservation. TestIntentReplacement owns the mode and bytes expectations; the other tests check bytes only. D-B5 repair 1 probes at 19422ae3a3d4e5de070edbe562bbf45f201f5421 (FT354-durable-callers), bench probe <file> ... --package ./internal/intent, no --run filter. The package passes 76 tests at the tip after every restore.\nProbe q04 (mode literal). internal/intent/intent.go swap old (exact): \"return replacer.Replace(path, append(data, '\\n'), 0o600)\" new (exact): \"return replacer.Replace(path, append(data, '\\n'), 0o640)\"\nResult: verdict bit, 3 failed: the three TestIntentReplacement rows at the mode check only, restored yes.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/intent.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,76\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,935\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,965,76\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestIntentReplacement/directory-close,\"replacement_test.go:42: /dev/shm/LY762D/t/TestIntentReplacementdirectory-close2426194141/001/.git/bench-intent.json mode = -rw-r-----, want a regular file at -rw-------\",1\n  github.com/gibbonmi/bench/internal/intent,TestIntentReplacement/directory-open,\"replacement_test.go:42: /dev/shm/LY762D/t/TestIntentReplacementdirectory-open1859290290/001/.git/bench-intent.json mode = -rw-r-----, want a regular file at -rw-------\",1\n  github.com/gibbonmi/bench/internal/intent,TestIntentReplacement/directory-sync,\"replacement_test.go:42: /dev/shm/LY762D/t/TestIntentReplacementdirectory-sync656104851/001/.git/bench-intent.json mode = -rw-r-----, want a regular file at -rw-------\",1\nskips[0]{package,test,reason}:\n\nProbe q05 (payload literal; --omit). internal/intent/intent.go --omit (exact): \"ledger.Schema = Schema\"\nResult: verdict bit, 4 failed: the three TestIntentReplacement rows at the bytes check and the existing TestPurgeAssignmentsDropsEveryLegacySchemaRecord, restored yes.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/intent.go,omit,failed,4,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,76\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1783\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,993,76\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestIntentReplacement/directory-close,\"replacement_test.go:42: /dev/shm/K3G64A/t/TestIntentReplacementdirectory-close4189935453/001/.git/bench-intent.json = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":0,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"added\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"worktree\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:02Z\\\\\\\\\\\"},{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"s… (509 bytes)\",1\n  github.com/gibbonmi/bench/internal/intent,TestIntentReplacement/directory-open,\"replacement_test.go:42: /dev/shm/K3G64A/t/TestIntentReplacementdirectory-open1180453917/001/.git/bench-intent.json = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":0,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"added\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"worktree\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:02Z\\\\\\\\\\\"},{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"se… (508 bytes)\",1\n  github.com/gibbonmi/bench/internal/intent,TestIntentReplacement/directory-sync,\"replacement_test.go:42: /dev/shm/K3G64A/t/TestIntentReplacementdirectory-sync190257201/001/.git/bench-intent.json = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":0,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"added\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"worktree\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:02Z\\\\\\\\\\\"},{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"see… (507 bytes)\",1\n  github.com/gibbonmi/bench/internal/intent,TestPurgeAssignmentsDropsEveryLegacySchemaRecord,\"purge_test.go:86: ledger after legacy purge = ledger.Ledger{Schema:1, Entries:[]ledger.Entry{}, Assignments:[]ledger.Assignment{}, CleanupReceipts:[]ledger.CleanupReceipt(nil), CommitmentReceipts:[]ledger.CommitmentReceipt(nil), MilestoneRe… (328 bytes)\",1\nskips[0]{package,test,reason}:\n\nProbe q07 (pre-rename refusal: fallback write of the new bytes after an unpublished error). internal/intent/intent.go swap old (exact): \"return replacer.Replace(path, append(data, '\\n'), 0o600)\" new (exact): \"if err := replacer.Replace(path, append(data, '\\n'), 0o600); err != nil { var e *durablefile.Error; if errors.As(err, &e) && !e.Published { _ = durablefile.Replace(path, append(data, '\\n'), 0o600) }; return err }; return nil\"\nResult: verdict bit, 2 failed: TestIntentPreRenameFailure and TestTransactionPublishedFailure/rename, each at its prior-ledger check only, restored yes.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/intent.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,76\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1102\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,1010,76\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestIntentPreRenameFailure,\"replacement_test.go:60: ledger = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":2,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"added\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"worktree\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:02Z\\\\\\\\\\\"}]}\\\\\\\\n\\\", want \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":2,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"seed\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"shift\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970… (748 bytes)\",1\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/rename,\"replacement_test.go:107: ledger = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":2,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"added\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"worktree\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:02Z\\\\\\\\\\\"},{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"seed\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"shift\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:01Z\\\\\\\\\\\"}],\\\\\\\\\\\"assignments\\\\\\\\\\\":[… (1020 bytes)\",2\nskips[0]{package,test,reason}:\n\nProbe q17 (unclassified failure overwrites the ledger). internal/intent/intent.go swap old (exact): \"return fmt.Errorf(\"encode intent ledger: %w\", err)\" new (exact): \"_ = durablefile.Replace(path, []byte(\"{}\\n\"), 0o600); return fmt.Errorf(\"encode intent ledger: %w\", err)\"\nResult: verdict bit, 1 failed: TestTransactionUnclassifiedFailureCompensates at the prior-ledger check only, restored yes.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/intent.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,76\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,632\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,1020,76\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestTransactionUnclassifiedFailureCompensates,\"replacement_test.go:138: ledger = \\\"{}\\\\\\\\n\\\", want \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":2,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"seed\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"shift\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:01Z\\\\\\\\\\\"}],\\\\\\\\\\\"assignments\\\\\\\\\\\":[{\\\\\\\\\\\"schema\\\\\\\\\\\":\\\\\\\\\\\"bench-assignment/v1\\\\\\\\\\\",\\\\\\\\\\\"id\\\\\\\\\\\":\\\\\\\\\\\"bbbbbbbbbbbbbbbbbbbbbbbb… (641 bytes)\",1\nskips[0]{package,test,reason}:\n\nProbe P1 (D22 owner bypass: the native owner in place of the injected replacer). internal/intent/intent.go swap old (exact): \"return replacer.Replace(path, append(data, '\\n'), 0o600)\" new (exact): \"return durablefile.Native().Replace(path, append(data, '\\n'), 0o600)\"\nResult: verdict bit, 9 failed at RequireFailure: the three TestIntentReplacement rows, TestIntentPreRenameFailure, TestReauthorizePublishedFailure, and the four TestTransactionPublishedFailure rows, restored yes.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/intent.go,swap,failed,9,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,76\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,2347\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,991,76\nfailures[9]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestIntentPreRenameFailure,\"replacement_test.go:59: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected ledger rename failure\",1\n  github.com/gibbonmi/bench/internal/intent,TestIntentReplacement/directory-close,\"replacement_test.go:41: err = <nil>, want a *durablefile.Error at stage directory-close with Published=true that wraps injected ledger directory-close failure\",1\n  github.com/gibbonmi/bench/internal/intent,TestIntentReplacement/directory-open,\"replacement_test.go:41: err = <nil>, want a *durablefile.Error at stage directory-open with Published=true that wraps injected ledger directory-open failure\",1\n  github.com/gibbonmi/bench/internal/intent,TestIntentReplacement/directory-sync,\"replacement_test.go:41: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected ledger directory-sync failure\",1\n  github.com/gibbonmi/bench/internal/intent,TestReauthorizePublishedFailure,\"replacement_test.go:166: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected ledger directory-sync failure\",1\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/directory-close,\"replacement_test.go:97: err = <nil>, want a *durablefile.Error at stage directory-close with Published=true that wraps injected ledger directory-close failure\",1\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/directory-open,\"replacement_test.go:97: err = <nil>, want a *durablefile.Error at stage directory-open with Published=true that wraps injected ledger directory-open failure\",1\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/directory-sync,\"replacement_test.go:97: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected ledger directory-sync failure\",1\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/rename,\"replacement_test.go:97: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected ledger rename failure\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "preservation",
          "command": "bench test --package ./internal/intent",
          "exit_code": 0,
          "probe": {
            "mutation": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "19422ae3a3d4e5de070edbe562bbf45f201f5421",
              "digest": "sha256:87b246bc2c1a5150c67d01572e31066e6b6308706cbcdb4bd6b45403386a06fb",
              "excerpt": "D-B5 preservation. TestIntentReplacement owns the mode and bytes expectations; the other tests check bytes only. D-B5 repair 1 probes at 19422ae3a3d4e5de070edbe562bbf45f201f5421 (FT354-durable-callers), bench probe <file> ... --package ./internal/intent, no --run filter. The package passes 76 tests at the tip after every restore.\nProbe q04 (mode literal). internal/intent/intent.go swap old (exact): \"return replacer.Replace(path, append(data, '\\n'), 0o600)\" new (exact): \"return replacer.Replace(path, append(data, '\\n'), 0o640)\"\nResult: verdict bit, 3 failed: the three TestIntentReplacement rows at the mode check only, restored yes.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/intent.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,76\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,935\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,965,76\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestIntentReplacement/directory-close,\"replacement_test.go:42: /dev/shm/LY762D/t/TestIntentReplacementdirectory-close2426194141/001/.git/bench-intent.json mode = -rw-r-----, want a regular file at -rw-------\",1\n  github.com/gibbonmi/bench/internal/intent,TestIntentReplacement/directory-open,\"replacement_test.go:42: /dev/shm/LY762D/t/TestIntentReplacementdirectory-open1859290290/001/.git/bench-intent.json mode = -rw-r-----, want a regular file at -rw-------\",1\n  github.com/gibbonmi/bench/internal/intent,TestIntentReplacement/directory-sync,\"replacement_test.go:42: /dev/shm/LY762D/t/TestIntentReplacementdirectory-sync656104851/001/.git/bench-intent.json mode = -rw-r-----, want a regular file at -rw-------\",1\nskips[0]{package,test,reason}:\n\nProbe q05 (payload literal; --omit). internal/intent/intent.go --omit (exact): \"ledger.Schema = Schema\"\nResult: verdict bit, 4 failed: the three TestIntentReplacement rows at the bytes check and the existing TestPurgeAssignmentsDropsEveryLegacySchemaRecord, restored yes.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/intent.go,omit,failed,4,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,76\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1783\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,993,76\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestIntentReplacement/directory-close,\"replacement_test.go:42: /dev/shm/K3G64A/t/TestIntentReplacementdirectory-close4189935453/001/.git/bench-intent.json = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":0,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"added\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"worktree\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:02Z\\\\\\\\\\\"},{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"s… (509 bytes)\",1\n  github.com/gibbonmi/bench/internal/intent,TestIntentReplacement/directory-open,\"replacement_test.go:42: /dev/shm/K3G64A/t/TestIntentReplacementdirectory-open1180453917/001/.git/bench-intent.json = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":0,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"added\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"worktree\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:02Z\\\\\\\\\\\"},{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"se… (508 bytes)\",1\n  github.com/gibbonmi/bench/internal/intent,TestIntentReplacement/directory-sync,\"replacement_test.go:42: /dev/shm/K3G64A/t/TestIntentReplacementdirectory-sync190257201/001/.git/bench-intent.json = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":0,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"added\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"worktree\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:02Z\\\\\\\\\\\"},{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"see… (507 bytes)\",1\n  github.com/gibbonmi/bench/internal/intent,TestPurgeAssignmentsDropsEveryLegacySchemaRecord,\"purge_test.go:86: ledger after legacy purge = ledger.Ledger{Schema:1, Entries:[]ledger.Entry{}, Assignments:[]ledger.Assignment{}, CleanupReceipts:[]ledger.CleanupReceipt(nil), CommitmentReceipts:[]ledger.CommitmentReceipt(nil), MilestoneRe… (328 bytes)\",1\nskips[0]{package,test,reason}:\n\nProbe q07 (pre-rename refusal: fallback write of the new bytes after an unpublished error). internal/intent/intent.go swap old (exact): \"return replacer.Replace(path, append(data, '\\n'), 0o600)\" new (exact): \"if err := replacer.Replace(path, append(data, '\\n'), 0o600); err != nil { var e *durablefile.Error; if errors.As(err, &e) && !e.Published { _ = durablefile.Replace(path, append(data, '\\n'), 0o600) }; return err }; return nil\"\nResult: verdict bit, 2 failed: TestIntentPreRenameFailure and TestTransactionPublishedFailure/rename, each at its prior-ledger check only, restored yes.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/intent.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,76\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1102\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,1010,76\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestIntentPreRenameFailure,\"replacement_test.go:60: ledger = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":2,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"added\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"worktree\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:02Z\\\\\\\\\\\"}]}\\\\\\\\n\\\", want \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":2,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"seed\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"shift\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970… (748 bytes)\",1\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/rename,\"replacement_test.go:107: ledger = \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":2,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"added\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"worktree\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:02Z\\\\\\\\\\\"},{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"seed\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"shift\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:01Z\\\\\\\\\\\"}],\\\\\\\\\\\"assignments\\\\\\\\\\\":[… (1020 bytes)\",2\nskips[0]{package,test,reason}:\n\nProbe q17 (unclassified failure overwrites the ledger). internal/intent/intent.go swap old (exact): \"return fmt.Errorf(\"encode intent ledger: %w\", err)\" new (exact): \"_ = durablefile.Replace(path, []byte(\"{}\\n\"), 0o600); return fmt.Errorf(\"encode intent ledger: %w\", err)\"\nResult: verdict bit, 1 failed: TestTransactionUnclassifiedFailureCompensates at the prior-ledger check only, restored yes.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/intent.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,76\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,632\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,1020,76\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestTransactionUnclassifiedFailureCompensates,\"replacement_test.go:138: ledger = \\\"{}\\\\\\\\n\\\", want \\\"{\\\\\\\\\\\"schema\\\\\\\\\\\":2,\\\\\\\\\\\"entries\\\\\\\\\\\":[{\\\\\\\\\\\"key\\\\\\\\\\\":\\\\\\\\\\\"seed\\\\\\\\\\\",\\\\\\\\\\\"kind\\\\\\\\\\\":\\\\\\\\\\\"shift\\\\\\\\\\\",\\\\\\\\\\\"created_at\\\\\\\\\\\":\\\\\\\\\\\"1970-01-01T00:00:01Z\\\\\\\\\\\"}],\\\\\\\\\\\"assignments\\\\\\\\\\\":[{\\\\\\\\\\\"schema\\\\\\\\\\\":\\\\\\\\\\\"bench-assignment/v1\\\\\\\\\\\",\\\\\\\\\\\"id\\\\\\\\\\\":\\\\\\\\\\\"bbbbbbbbbbbbbbbbbbbbbbbb… (641 bytes)\",1\nskips[0]{package,test,reason}:\n\nProbe P1 (D22 owner bypass: the native owner in place of the injected replacer). internal/intent/intent.go swap old (exact): \"return replacer.Replace(path, append(data, '\\n'), 0o600)\" new (exact): \"return durablefile.Native().Replace(path, append(data, '\\n'), 0o600)\"\nResult: verdict bit, 9 failed at RequireFailure: the three TestIntentReplacement rows, TestIntentPreRenameFailure, TestReauthorizePublishedFailure, and the four TestTransactionPublishedFailure rows, restored yes.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/intent/intent.go,swap,failed,9,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,all,passed,76\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,2347\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,991,76\nfailures[9]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestIntentPreRenameFailure,\"replacement_test.go:59: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected ledger rename failure\",1\n  github.com/gibbonmi/bench/internal/intent,TestIntentReplacement/directory-close,\"replacement_test.go:41: err = <nil>, want a *durablefile.Error at stage directory-close with Published=true that wraps injected ledger directory-close failure\",1\n  github.com/gibbonmi/bench/internal/intent,TestIntentReplacement/directory-open,\"replacement_test.go:41: err = <nil>, want a *durablefile.Error at stage directory-open with Published=true that wraps injected ledger directory-open failure\",1\n  github.com/gibbonmi/bench/internal/intent,TestIntentReplacement/directory-sync,\"replacement_test.go:41: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected ledger directory-sync failure\",1\n  github.com/gibbonmi/bench/internal/intent,TestReauthorizePublishedFailure,\"replacement_test.go:166: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected ledger directory-sync failure\",1\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/directory-close,\"replacement_test.go:97: err = <nil>, want a *durablefile.Error at stage directory-close with Published=true that wraps injected ledger directory-close failure\",1\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/directory-open,\"replacement_test.go:97: err = <nil>, want a *durablefile.Error at stage directory-open with Published=true that wraps injected ledger directory-open failure\",1\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/directory-sync,\"replacement_test.go:97: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected ledger directory-sync failure\",1\n  github.com/gibbonmi/bench/internal/intent,TestTransactionPublishedFailure/rename,\"replacement_test.go:97: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected ledger rename failure\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "dcm-b5-r1-standards",
          "performer": "claude:dcm_rev_b5_std",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "f31faf8613dc4f4f281fc93c6c44040dd859c815",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b5_std",
            "digest": "sha256:49cba09eafd1d55ca218b17adc51ab66277d2cf325665851857dc5960815d7a9",
            "excerpt": "## Standards (D-B5, 2d06e759..ba35cc83) — 0 findings\nOne transact body serves Transact, TransactWithReplacer, and PurgeAssignments; no copy of the durablefile algorithm; shared durablefiletest assertions and PostRenameStages; comments state only the why; each fact checked once. Veto items accepted: operand order, leaf error text, repeated upsert and reauthorization wants (both with recorded reds).\nAdvice: a bytes-only assertion now has three spellings; name the upsert expectations as fixture constants; the published classifier may be derived again in D-B6 or D-C4.\n"
          },
          "axis": "Standards",
          "base": "2d06e7594c9d31f1ce9f22102d03da660601cea2",
          "tip": "a6c0ffe56d68d2499f336761de2589f5988801da",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "dcm-b5-r1-spec",
          "performer": "claude:dcm_rev_b5_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "f31faf8613dc4f4f281fc93c6c44040dd859c815",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b5_spec",
            "digest": "sha256:6172111808a85e9032ffe6d4d42e751af97c22a4311266086df39bec84138ccb",
            "excerpt": "## Spec (D-B5, 2d06e759..ba35cc83) — 1 finding; worst P1\nP1 | auto-fix | 6 | reviews/durable-caller-migration.md (D-B5 omission-proof); spec.md:288-289 | the record holds no D22 owner-bypass red with exact swap text | record a probe that swaps intent.go:216 for the former rename sequence or durablefile.Native(), with the red TestIntentReplacement rows, restore, and pass\nD22, D44, D45, D52, the unclassified branch, the seam shape, and the D36 evidence pass. Ticket 07 can consume TransactWithReplacer. The two plan commits add only the 10 and 13 assignments.\nAdvice: a %w wrap would keep the intent ledger label in diagnostics.\n"
          },
          "axis": "Spec",
          "base": "2d06e7594c9d31f1ce9f22102d03da660601cea2",
          "tip": "a6c0ffe56d68d2499f336761de2589f5988801da",
          "finding_ids": [
            "P1"
          ],
          "supersedes": []
        },
        {
          "id": "dcm-b5-r1-coverage",
          "performer": "claude:dcm_rev_b5_cov",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "f31faf8613dc4f4f281fc93c6c44040dd859c815",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b5_cov",
            "digest": "sha256:75b6fdea2526cfb93dc75cb8c415c3ee3757fde166898738fdcc70e391245099",
            "excerpt": "## Coverage (D-B5, 2d06e759..ba35cc83) — 2 findings; worst C1. Probes q01-q18 re-run: each bit, restored=yes.\nC1 | auto-fix | 9 | internal/intent/transaction.go:88-91 | no rename-stage or directory-open/close fault reaches TransactWithReplacer; three swaps of the published predicate are silent (Stage != \"create\"; Published || Stage == \"rename\"; Stage == \"directory-sync\") | through TransactWithReplacer with the Recorder, range PostRenameStages (compensation 0, new bytes) and inject rename (compensation 1, old bytes); record the three swaps as reds\nC2 | auto-fix | 8 | internal/intent/transaction.go:65,78-79 | releasing the lock before compensation is silent; the D45 test cannot catch it | in the rename-stage test, assert inside the compensation that the lock exists, and after the call that it is gone; record the release-before-compensate swap as its red\nq11 judged acceptable as a bounded necessity red.\n"
          },
          "axis": "Coverage",
          "base": "2d06e7594c9d31f1ce9f22102d03da660601cea2",
          "tip": "a6c0ffe56d68d2499f336761de2589f5988801da",
          "finding_ids": [
            "C1",
            "C2"
          ],
          "supersedes": []
        },
        {
          "id": "dcm-b5-r2-standards",
          "performer": "claude:dcm_rev_b5_std2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "c4af640671be8b9686510058a2b060a939e82f09",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b5_std2",
            "digest": "sha256:3781ab193f12a9a3d56ecd20a488eb08800bd6273fa5d04003c1fec32fe27634",
            "excerpt": "## Standards confirming round (D-B5, repair delta c73da9c8..d12060c6) — 0 findings\nThe table form of TestTransactionPublishedFailure takes its stages from durablefiletest.PostRenameStages, classifies through RequireFailure, and checks the held lock once inside the compensation; the doc comment states only the constraint.\nAdvice: the rename row's prior-bytes check overlaps TestIntentPreRenameFailure; the test name now also holds an unpublished row.\n"
          },
          "axis": "Standards",
          "base": "2d06e7594c9d31f1ce9f22102d03da660601cea2",
          "tip": "19422ae3a3d4e5de070edbe562bbf45f201f5421",
          "finding_ids": [],
          "supersedes": [
            "dcm-b5-r1-standards"
          ]
        },
        {
          "id": "dcm-b5-r2-spec",
          "performer": "claude:dcm_rev_b5_spec2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "c4af640671be8b9686510058a2b060a939e82f09",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b5_spec2",
            "digest": "sha256:447213e8be3f5f44c24a83851748e503959519fc1a693f092453595606d3c246",
            "excerpt": "## Spec confirming round (D-B5, repair delta c73da9c8..d12060c6) — 0 findings\nP1 closed: dcm-t06-rep1-preservation records replacer.Replace -> durablefile.Native().Replace in intent.go (bit, 9 failed). D44, D45, D52 and the intent paragraph (compensation once under the held lock) stay intact.\n"
          },
          "axis": "Spec",
          "base": "2d06e7594c9d31f1ce9f22102d03da660601cea2",
          "tip": "19422ae3a3d4e5de070edbe562bbf45f201f5421",
          "finding_ids": [],
          "supersedes": [
            "dcm-b5-r1-spec"
          ]
        },
        {
          "id": "dcm-b5-r2-coverage",
          "performer": "claude:dcm_rev_b5_cov2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "c4af640671be8b9686510058a2b060a939e82f09",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b5_cov2",
            "digest": "sha256:fe666f1b5ccc51f1e8f02a88e94128d1c196fa1dcbafc18a9298bfec77e84de0",
            "excerpt": "## Coverage confirming round (D-B5, repair delta c73da9c8..d12060c6) — 0 findings\nC1 closed: new probe, replaceErr.Published -> replaceErr.Stage != \"rename\" in published: bit, 2 tests, restored=yes. C2 closed: the recorded release-before-compensate probe fails at the lock check inside the compensation.\n"
          },
          "axis": "Coverage",
          "base": "2d06e7594c9d31f1ce9f22102d03da660601cea2",
          "tip": "19422ae3a3d4e5de070edbe562bbf45f201f5421",
          "finding_ids": [],
          "supersedes": [
            "dcm-b5-r1-coverage"
          ]
        }
      ]
    },
    {
      "id": "D-C1",
      "base": "641ddbaa6c7f22023c6453abee122954564b8eca",
      "tip": "d816ae02466c9c241ae8330dbb02d894fb6f0344",
      "plan_digest": "sha256:95a4fa23af538a1a3b21d7b5f8d34d8497b07fa508e4ad0dad253db7925d183c",
      "source_digest": "2c2a7a21b89b514ed42ea5f51281c1ef521b6dd0",
      "acceptance_rows": [
        "D24"
      ],
      "verification": [
        {
          "id": "dcm-t08-caller-tests",
          "performer": "claude:dcm_t08",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d5010a7c7ea7de4e46f377d692d2ed05a9a2adcc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "17ffac99dfbfaf0ac0d69e340c172839cffe9b29",
            "digest": "sha256:fe85a28cfcf8c98a512faeda2fee75f16c3fd2e9512fbce89861009b502fb425",
            "excerpt": "bench test --package ./internal/publication at 17ffac99dfbfaf0ac0d69e340c172839cffe9b29\nD24 witness: TestPublicationReplacement. D36 pre-rename evidence: TestPublicationPreRenameFailure.\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,17ffac99dfbfaf0ac0d69e340c172839cffe9b29,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/publication,pass,1157,28\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "caller-tests",
          "command": "bench test --package ./internal/publication",
          "exit_code": 0
        },
        {
          "id": "dcm-t08-port-registry",
          "performer": "claude:dcm_t08",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d5010a7c7ea7de4e46f377d692d2ed05a9a2adcc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "17ffac99dfbfaf0ac0d69e340c172839cffe9b29",
            "digest": "sha256:f7d8ee532d68dc59498803ff9549ab3a6127e0b3494d4673611b24f2c7e9d440",
            "excerpt": "bench test --check injected-port-registry at 17ffac99dfbfaf0ac0d69e340c172839cffe9b29 (SaveRecord adds no audited port; registry unchanged)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,17ffac99dfbfaf0ac0d69e340c172839cffe9b29,false\ncheck[1]{name,kind,tests_run,subjects}:\n  injected-port-registry,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,26,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "port-registry",
          "command": "bench test --check injected-port-registry",
          "exit_code": 0
        },
        {
          "id": "dcm-t08-conformance",
          "performer": "claude:dcm_t08",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d5010a7c7ea7de4e46f377d692d2ed05a9a2adcc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "17ffac99dfbfaf0ac0d69e340c172839cffe9b29",
            "digest": "sha256:74ebd3d854171ce7ce08b31475170da97a13e30f1dca2dc948c113f7ba8188dd",
            "excerpt": "bench test --package ./internal/conformance at 17ffac99dfbfaf0ac0d69e340c172839cffe9b29\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,17ffac99dfbfaf0ac0d69e340c172839cffe9b29,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,74430,1690\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "dcm-t08-root-conformance",
          "performer": "claude:dcm_t08",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d5010a7c7ea7de4e46f377d692d2ed05a9a2adcc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "17ffac99dfbfaf0ac0d69e340c172839cffe9b29",
            "digest": "sha256:1b15f1751fd97fcd1f149ffbe5eafe89b61f24badc448d75d6f7a85d27703a5a",
            "excerpt": "bench test --check root-conformance at 17ffac99dfbfaf0ac0d69e340c172839cffe9b29\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,17ffac99dfbfaf0ac0d69e340c172839cffe9b29,false\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,23383,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "dcm-t08-omission-proof",
          "performer": "claude:dcm_t08",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d5010a7c7ea7de4e46f377d692d2ed05a9a2adcc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "17ffac99dfbfaf0ac0d69e340c172839cffe9b29",
            "digest": "sha256:885a24e1b474b80b34e269a70d91a4b0665795093bc5fedef3541c7421c7710e",
            "excerpt": "# omission-proof: owner-bypass mutations of SaveRecord through bench probe at 17ffac99dfbfaf0ac0d69e340c172839cffe9b29, whole package ./internal/publication, no --run, each restored\n## 1 former writer. subject internal/publication/record.go\nold: \tif err := r.Replace(path, append(data, '\\n'), 0o644); err != nil {\nnew: \tif err := func() error { temp := path + \".tmp\"; if err := os.WriteFile(temp, append(data, '\\n'), 0o644); err != nil { return err }; return os.Rename(temp, path) }(); err != nil {\nred: TestPublicationReplacement (D24, directory-sync failure hidden) and TestPublicationPreRenameFailure (D36, rename failure hidden); restored yes; package pass after restore\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,17ffac99dfbfaf0ac0d69e340c172839cffe9b29,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/publication/record.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/publication,all,passed,28\nspilled{lines=12,bytes=976,omitted_lines=3,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791833680617547039-750e43e1e2fe7999.out}\n  github.com/gibbonmi/bench/internal/publication,fail,1323,28\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/publication,TestPublicationPreRenameFailure,\"replacement_test.go:70: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected publication rename failure\",1\n  github.com/gibbonmi/bench/internal/publication,TestPublicationReplacement,\"replacement_test.go:34: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected publication directory-sync failure\",1\nskips[0]{package,test,reason}:\n## 2 native owner in place of the injected replacer\nold: if err := r.Replace(path,\nnew: if err := durablefile.Native().Replace(path,\nred: the same two tests; restored yes\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,17ffac99dfbfaf0ac0d69e340c172839cffe9b29,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/publication/record.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/publication,all,passed,28\nspilled{lines=12,bytes=976,omitted_lines=3,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791833692344765788-e7a01be492253800.out}\n  github.com/gibbonmi/bench/internal/publication,fail,1136,28\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/publication,TestPublicationPreRenameFailure,\"replacement_test.go:70: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected publication rename failure\",1\n  github.com/gibbonmi/bench/internal/publication,TestPublicationReplacement,\"replacement_test.go:34: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected publication directory-sync failure\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "omission-proof",
          "command": "bench test --package ./internal/publication",
          "exit_code": 0,
          "probe": {
            "mutation": "Restore the former unsynchronized writer. The publication record junction test must detect the missing failure.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "17ffac99dfbfaf0ac0d69e340c172839cffe9b29",
              "digest": "sha256:885a24e1b474b80b34e269a70d91a4b0665795093bc5fedef3541c7421c7710e",
              "excerpt": "# omission-proof: owner-bypass mutations of SaveRecord through bench probe at 17ffac99dfbfaf0ac0d69e340c172839cffe9b29, whole package ./internal/publication, no --run, each restored\n## 1 former writer. subject internal/publication/record.go\nold: \tif err := r.Replace(path, append(data, '\\n'), 0o644); err != nil {\nnew: \tif err := func() error { temp := path + \".tmp\"; if err := os.WriteFile(temp, append(data, '\\n'), 0o644); err != nil { return err }; return os.Rename(temp, path) }(); err != nil {\nred: TestPublicationReplacement (D24, directory-sync failure hidden) and TestPublicationPreRenameFailure (D36, rename failure hidden); restored yes; package pass after restore\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,17ffac99dfbfaf0ac0d69e340c172839cffe9b29,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/publication/record.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/publication,all,passed,28\nspilled{lines=12,bytes=976,omitted_lines=3,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791833680617547039-750e43e1e2fe7999.out}\n  github.com/gibbonmi/bench/internal/publication,fail,1323,28\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/publication,TestPublicationPreRenameFailure,\"replacement_test.go:70: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected publication rename failure\",1\n  github.com/gibbonmi/bench/internal/publication,TestPublicationReplacement,\"replacement_test.go:34: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected publication directory-sync failure\",1\nskips[0]{package,test,reason}:\n## 2 native owner in place of the injected replacer\nold: if err := r.Replace(path,\nnew: if err := durablefile.Native().Replace(path,\nred: the same two tests; restored yes\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,17ffac99dfbfaf0ac0d69e340c172839cffe9b29,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/publication/record.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/publication,all,passed,28\nspilled{lines=12,bytes=976,omitted_lines=3,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791833692344765788-e7a01be492253800.out}\n  github.com/gibbonmi/bench/internal/publication,fail,1136,28\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/publication,TestPublicationPreRenameFailure,\"replacement_test.go:70: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected publication rename failure\",1\n  github.com/gibbonmi/bench/internal/publication,TestPublicationReplacement,\"replacement_test.go:34: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected publication directory-sync failure\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dcm-t08-preservation",
          "performer": "claude:dcm_t08",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "d5010a7c7ea7de4e46f377d692d2ed05a9a2adcc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "17ffac99dfbfaf0ac0d69e340c172839cffe9b29",
            "digest": "sha256:81c0b3c34e17f7bd22a31c10a7d478ff1003785bc4d8ff554ca24a588043445c",
            "excerpt": "# preservation: necessity reds of each independent check at 17ffac99dfbfaf0ac0d69e340c172839cffe9b29, whole package ./internal/publication, no --run, each restored\n## 3 mode drift: only TestPublicationReplacement (D24 mode check) fails\nold: append(data, '\\n'), 0o644)\nnew: append(data, '\\n'), 0o640)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,17ffac99dfbfaf0ac0d69e340c172839cffe9b29,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/publication/record.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/publication,all,passed,28\nspilled{lines=11,bytes=759,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791833699894384376-023334cc254ed1bc.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/publication,fail,1085,28\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/publication,TestPublicationReplacement,\"replacement_test.go:35: /dev/shm/DZURNZ/t/TestPublicationReplacement3995141172/001/dist/publication/publication-record.json mode = -rw-r-----, want a regular file at -rw-r--r--\",1\nskips[0]{package,test,reason}:\n## 4 final-newline omission: only TestPublicationReplacement (D24 payload bytes check) fails\nold: r.Replace(path, append(data, '\\n'), 0o644)\nnew: r.Replace(path, data, 0o644)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,17ffac99dfbfaf0ac0d69e340c172839cffe9b29,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/publication/record.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/publication,all,passed,28\nspilled{lines=11,bytes=903,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791833707300272744-10870a42ea522570.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/publication,fail,1039,28\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/publication,TestPublicationReplacement,\"replacement_test.go:35: /dev/shm/YV6I3X/t/TestPublicationReplacement2687722160/001/dist/publication/publication-record.json = \\\"{\\\\\\\\n  \\\\\\\\\\\"schema_version\\\\\\\\\\\": 1,\\\\\\\\n  \\\\\\\\\\\"release_index_sha256\\\\\\\\\\\": \\\\\\\\\\\"abc\\\\\\\\\\\",\\\\\\\\n  \\\\\\\\\\\"path\\\\\\\\\\\": \\\\\\\\\\\"first\\\\\\\\\\\",\\\\\\\\n  \\\\\\\\\\\"profile\\\\\\\\\\\": \\\\\\\\\\\"publi… (747 bytes)\",1\nskips[0]{package,test,reason}:\n## 5 fallback write of the new bytes after an unpublished error: only TestPublicationPreRenameFailure (D36 prior-bytes check) fails at replacement_test.go:72\nold: \tif err := r.Replace(path, append(data, '\\n'), 0o644); err != nil {\nnew: \tif err := r.Replace(path, append(data, '\\n'), 0o644); err != nil {\n\t\t_ = os.WriteFile(path, append(data, '\\n'), 0o644)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,17ffac99dfbfaf0ac0d69e340c172839cffe9b29,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/publication/record.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/publication,all,passed,28\nspilled{lines=11,bytes=959,omitted_lines=2,cut_lines=1,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791833714879015462-073e62a81c8b30d0.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/publication,fail,1045,28\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/publication,TestPublicationPreRenameFailure,\"replacement_test.go:72: record = \\\"{\\\\\\\\n  \\\\\\\\\\\"schema_version\\\\\\\\\\\": 1,\\\\\\\\n  \\\\\\\\\\\"release_index_sha256\\\\\\\\\\\": \\\\\\\\\\\"abc\\\\\\\\\\\",\\\\\\\\n  \\\\\\\\\\\"path\\\\\\\\\\\": \\\\\\\\\\\"first\\\\\\\\\\\",\\\\\\\\n  \\\\\\\\\\\"profile\\\\\\\\\\\": \\\\\\\\\\\"public\\\\\\\\\\\",\\\\\\\\n  \\\\\\\\\\\"transitions\\\\\\\\\\\": [],\\\\\\\\n  \\\\\\\\\\\"provenance\\\\\\\\\\\": [\\\\\\\\n    {\\\\\\\\n      \\\\\\\\\\\"package\\\\\\\\\\\"\nskips[0]{package,test,reason}:\n## The RequireFailure checks of both rows go red under probes 1 and 2 (omission-proof excerpt).\n"
          },
          "requirement": "preservation",
          "command": "bench test --package ./internal/publication",
          "exit_code": 0,
          "probe": {
            "mutation": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "17ffac99dfbfaf0ac0d69e340c172839cffe9b29",
              "digest": "sha256:81c0b3c34e17f7bd22a31c10a7d478ff1003785bc4d8ff554ca24a588043445c",
              "excerpt": "# preservation: necessity reds of each independent check at 17ffac99dfbfaf0ac0d69e340c172839cffe9b29, whole package ./internal/publication, no --run, each restored\n## 3 mode drift: only TestPublicationReplacement (D24 mode check) fails\nold: append(data, '\\n'), 0o644)\nnew: append(data, '\\n'), 0o640)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,17ffac99dfbfaf0ac0d69e340c172839cffe9b29,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/publication/record.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/publication,all,passed,28\nspilled{lines=11,bytes=759,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791833699894384376-023334cc254ed1bc.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/publication,fail,1085,28\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/publication,TestPublicationReplacement,\"replacement_test.go:35: /dev/shm/DZURNZ/t/TestPublicationReplacement3995141172/001/dist/publication/publication-record.json mode = -rw-r-----, want a regular file at -rw-r--r--\",1\nskips[0]{package,test,reason}:\n## 4 final-newline omission: only TestPublicationReplacement (D24 payload bytes check) fails\nold: r.Replace(path, append(data, '\\n'), 0o644)\nnew: r.Replace(path, data, 0o644)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,17ffac99dfbfaf0ac0d69e340c172839cffe9b29,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/publication/record.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/publication,all,passed,28\nspilled{lines=11,bytes=903,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791833707300272744-10870a42ea522570.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/publication,fail,1039,28\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/publication,TestPublicationReplacement,\"replacement_test.go:35: /dev/shm/YV6I3X/t/TestPublicationReplacement2687722160/001/dist/publication/publication-record.json = \\\"{\\\\\\\\n  \\\\\\\\\\\"schema_version\\\\\\\\\\\": 1,\\\\\\\\n  \\\\\\\\\\\"release_index_sha256\\\\\\\\\\\": \\\\\\\\\\\"abc\\\\\\\\\\\",\\\\\\\\n  \\\\\\\\\\\"path\\\\\\\\\\\": \\\\\\\\\\\"first\\\\\\\\\\\",\\\\\\\\n  \\\\\\\\\\\"profile\\\\\\\\\\\": \\\\\\\\\\\"publi… (747 bytes)\",1\nskips[0]{package,test,reason}:\n## 5 fallback write of the new bytes after an unpublished error: only TestPublicationPreRenameFailure (D36 prior-bytes check) fails at replacement_test.go:72\nold: \tif err := r.Replace(path, append(data, '\\n'), 0o644); err != nil {\nnew: \tif err := r.Replace(path, append(data, '\\n'), 0o644); err != nil {\n\t\t_ = os.WriteFile(path, append(data, '\\n'), 0o644)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,17ffac99dfbfaf0ac0d69e340c172839cffe9b29,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/publication/record.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/publication,all,passed,28\nspilled{lines=11,bytes=959,omitted_lines=2,cut_lines=1,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791833714879015462-073e62a81c8b30d0.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/publication,fail,1045,28\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/publication,TestPublicationPreRenameFailure,\"replacement_test.go:72: record = \\\"{\\\\\\\\n  \\\\\\\\\\\"schema_version\\\\\\\\\\\": 1,\\\\\\\\n  \\\\\\\\\\\"release_index_sha256\\\\\\\\\\\": \\\\\\\\\\\"abc\\\\\\\\\\\",\\\\\\\\n  \\\\\\\\\\\"path\\\\\\\\\\\": \\\\\\\\\\\"first\\\\\\\\\\\",\\\\\\\\n  \\\\\\\\\\\"profile\\\\\\\\\\\": \\\\\\\\\\\"public\\\\\\\\\\\",\\\\\\\\n  \\\\\\\\\\\"transitions\\\\\\\\\\\": [],\\\\\\\\n  \\\\\\\\\\\"provenance\\\\\\\\\\\": [\\\\\\\\n    {\\\\\\\\n      \\\\\\\\\\\"package\\\\\\\\\\\"\nskips[0]{package,test,reason}:\n## The RequireFailure checks of both rows go red under probes 1 and 2 (omission-proof excerpt).\n"
            }
          }
        },
        {
          "id": "dcm-t08-rep1-caller-tests",
          "performer": "claude:dcm_t08_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2c2a7a21b89b514ed42ea5f51281c1ef521b6dd0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "d816ae02466c9c241ae8330dbb02d894fb6f0344",
            "digest": "sha256:4f717c5329614451135fb88a2383415ff3c84b2b010673162ae6bf0362eddcf8",
            "excerpt": "bench test --package ./internal/publication at d816ae02466c9c241ae8330dbb02d894fb6f0344 (repair 1; dirty path is the review record only)\nD24 witness: TestPublicationReplacement. D36 pre-rename evidence: TestPublicationPreRenameFailure (prior bytes through durablefiletest.RequireBytes). D30 refusal witness: TestPublicationDirectoryRefusal.\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,d816ae02466c9c241ae8330dbb02d894fb6f0344,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/publication,pass,1042,29\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "caller-tests",
          "command": "bench test --package ./internal/publication",
          "exit_code": 0
        },
        {
          "id": "dcm-t08-rep1-port-registry",
          "performer": "claude:dcm_t08_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2c2a7a21b89b514ed42ea5f51281c1ef521b6dd0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "d816ae02466c9c241ae8330dbb02d894fb6f0344",
            "digest": "sha256:a46e1c7cf5999a69bd8eea8daa90082b181814404bb7e7b0813daf5119ef42e7",
            "excerpt": "bench test --check injected-port-registry at d816ae02466c9c241ae8330dbb02d894fb6f0344 (repair 1 adds no audited port; registry unchanged)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,d816ae02466c9c241ae8330dbb02d894fb6f0344,true\ncheck[1]{name,kind,tests_run,subjects}:\n  injected-port-registry,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,23,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "port-registry",
          "command": "bench test --check injected-port-registry",
          "exit_code": 0
        },
        {
          "id": "dcm-t08-rep1-conformance",
          "performer": "claude:dcm_t08_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2c2a7a21b89b514ed42ea5f51281c1ef521b6dd0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "d816ae02466c9c241ae8330dbb02d894fb6f0344",
            "digest": "sha256:c8638482e3b3b8e3e906ee6ca3dae874cf2f847b11b686ee603423d31f1a39cb",
            "excerpt": "bench test --package ./internal/conformance at d816ae02466c9c241ae8330dbb02d894fb6f0344\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,d816ae02466c9c241ae8330dbb02d894fb6f0344,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,55518,1690\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "dcm-t08-rep1-durablefiletest",
          "performer": "claude:dcm_t08_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2c2a7a21b89b514ed42ea5f51281c1ef521b6dd0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "d816ae02466c9c241ae8330dbb02d894fb6f0344",
            "digest": "sha256:dac683fa16593e7426944b527071f21564801730279f21f9dda942782b79fbd4",
            "excerpt": "bench test --package ./internal/durablefile/durablefiletest at d816ae02466c9c241ae8330dbb02d894fb6f0344 (the package compiles and has no tests of its own; RequireBytes and RequireFile run in the caller packages, whose focused run passed)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,d816ae02466c9c241ae8330dbb02d894fb6f0344,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile/durablefiletest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "durablefiletest",
          "command": "bench test --package ./internal/durablefile/durablefiletest",
          "exit_code": 0
        },
        {
          "id": "dcm-t08-rep1-root-conformance",
          "performer": "claude:dcm_t08_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2c2a7a21b89b514ed42ea5f51281c1ef521b6dd0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "d816ae02466c9c241ae8330dbb02d894fb6f0344",
            "digest": "sha256:80a4a9d5e914d2ae408e3d5f72e63987a4ea249843c549951a35a4e22a63ea94",
            "excerpt": "bench test --check root-conformance at d816ae02466c9c241ae8330dbb02d894fb6f0344\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,d816ae02466c9c241ae8330dbb02d894fb6f0344,true\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17321,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "dcm-t08-rep1-omission-proof",
          "performer": "claude:dcm_t08_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2c2a7a21b89b514ed42ea5f51281c1ef521b6dd0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "d816ae02466c9c241ae8330dbb02d894fb6f0344",
            "digest": "sha256:fac62bdda7787332a12c4d19de80b4944a51ef3a2f329eac1f3b173064cd0d75",
            "excerpt": "# omission-proof: owner-bypass mutations of SaveRecord through bench probe at d816ae02466c9c241ae8330dbb02d894fb6f0344, whole package ./internal/publication, no --run, each restored\n## 1 former writer. subject internal/publication/record.go\nold: if err := r.Replace(path, append(data, '\\n'), 0o644); err != nil {\nnew: if err := func() error { temp := path + \".tmp\"; if err := os.WriteFile(temp, append(data, '\\n'), 0o644); err != nil { return err }; return os.Rename(temp, path) }(); err != nil {\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,d816ae02466c9c241ae8330dbb02d894fb6f0344,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/publication/record.go,swap,failed,2,yes\n  github.com/gibbonmi/bench/internal/publication,fail,1146,29\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/publication,TestPublicationPreRenameFailure,\"replacement_test.go:72: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected publication rename failure\",1\n  github.com/gibbonmi/bench/internal/publication,TestPublicationReplacement,\"replacement_test.go:36: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected publication directory-sync failure\",1\n## 2 native owner in place of the injected replacer\nold: if err := r.Replace(path,\nnew: if err := durablefile.Native().Replace(path,\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,d816ae02466c9c241ae8330dbb02d894fb6f0344,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/publication/record.go,swap,failed,2,yes\n  github.com/gibbonmi/bench/internal/publication,fail,1002,29\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/publication,TestPublicationPreRenameFailure,\"replacement_test.go:72: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected publication rename failure\",1\n  github.com/gibbonmi/bench/internal/publication,TestPublicationReplacement,\"replacement_test.go:36: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected publication directory-sync failure\",1\n## Package passes after each restore (bench test --package ./internal/publication: pass, 29 tests).\n"
          },
          "requirement": "omission-proof",
          "command": "bench test --package ./internal/publication",
          "exit_code": 0,
          "probe": {
            "mutation": "Restore the former unsynchronized writer. The publication record junction test must detect the missing failure.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "d816ae02466c9c241ae8330dbb02d894fb6f0344",
              "digest": "sha256:fac62bdda7787332a12c4d19de80b4944a51ef3a2f329eac1f3b173064cd0d75",
              "excerpt": "# omission-proof: owner-bypass mutations of SaveRecord through bench probe at d816ae02466c9c241ae8330dbb02d894fb6f0344, whole package ./internal/publication, no --run, each restored\n## 1 former writer. subject internal/publication/record.go\nold: if err := r.Replace(path, append(data, '\\n'), 0o644); err != nil {\nnew: if err := func() error { temp := path + \".tmp\"; if err := os.WriteFile(temp, append(data, '\\n'), 0o644); err != nil { return err }; return os.Rename(temp, path) }(); err != nil {\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,d816ae02466c9c241ae8330dbb02d894fb6f0344,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/publication/record.go,swap,failed,2,yes\n  github.com/gibbonmi/bench/internal/publication,fail,1146,29\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/publication,TestPublicationPreRenameFailure,\"replacement_test.go:72: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected publication rename failure\",1\n  github.com/gibbonmi/bench/internal/publication,TestPublicationReplacement,\"replacement_test.go:36: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected publication directory-sync failure\",1\n## 2 native owner in place of the injected replacer\nold: if err := r.Replace(path,\nnew: if err := durablefile.Native().Replace(path,\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,d816ae02466c9c241ae8330dbb02d894fb6f0344,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/publication/record.go,swap,failed,2,yes\n  github.com/gibbonmi/bench/internal/publication,fail,1002,29\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/publication,TestPublicationPreRenameFailure,\"replacement_test.go:72: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected publication rename failure\",1\n  github.com/gibbonmi/bench/internal/publication,TestPublicationReplacement,\"replacement_test.go:36: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected publication directory-sync failure\",1\n## Package passes after each restore (bench test --package ./internal/publication: pass, 29 tests).\n"
            }
          }
        },
        {
          "id": "dcm-t08-rep1-preservation",
          "performer": "claude:dcm_t08_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2c2a7a21b89b514ed42ea5f51281c1ef521b6dd0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "d816ae02466c9c241ae8330dbb02d894fb6f0344",
            "digest": "sha256:7ed614153340dd16fc2e2dd61bf55b6a35f5a35c16653fc2d16a7382108306d4",
            "excerpt": "# preservation: necessity reds of each independent check at d816ae02466c9c241ae8330dbb02d894fb6f0344, whole package ./internal/publication, no --run, each restored\n## 3 mode drift: only TestPublicationReplacement (D24 mode check, RequireFile) fails\nold: append(data, '\\n'), 0o644)\nnew: append(data, '\\n'), 0o640)\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/publication/record.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/publication,TestPublicationReplacement,\"replacement_test.go:37: .../dist/publication/publication-record.json mode = -rw-r-----, want a regular file at -rw-r--r--\",1\n## 4 final-newline omission: only TestPublicationReplacement fails at its RequireFile call (byte step through RequireBytes)\nold: r.Replace(path, append(data, '\\n'), 0o644)\nnew: r.Replace(path, data, 0o644)\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/publication/record.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/publication,TestPublicationReplacement,\"replacement_test.go:37: .../dist/publication/publication-record.json = \\\"{\\\\n  \\\\\\\"schema_version\\\\\\\": 1, ... (747 bytes)\",1\n## 5 fallback write of the new bytes after an unpublished error: only TestPublicationPreRenameFailure fails, at its RequireBytes call (replacement_test.go:73)\nold: r.Replace(path, append(data, '\\n'), 0o644); err != nil {\nnew: r.Replace(path, append(data, '\\n'), 0o644); err != nil { _ = os.WriteFile(path, append(data, '\\n'), 0o644)\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/publication/record.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/publication,TestPublicationPreRenameFailure,\"replacement_test.go:73: .../dist/publication/publication-record.json = \\\"{\\\\n  \\\\\\\"schema_version\\\\\\\": 1, ... (757 bytes)\",1\n## 6 refusal swallow (D30): only TestPublicationDirectoryRefusal fails; before the repair this swap was silent (verdict silent, 28 tests passed at f46af6a5)\nold: return fmt.Errorf(\"could not create publication record directory: %w\", err)\nnew: _ = err\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/publication/record.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/publication,TestPublicationDirectoryRefusal,\"replacement_test.go:90: err = could not write publication record: replace .../dist/publication/publication-record.json: create: open ... (369 bytes)\",1\n## Each probe tree: FT354-durable-callers,d816ae02466c9c241ae8330dbb02d894fb6f0344,false; each package run 29 tests; package passes after each restore.\n"
          },
          "requirement": "preservation",
          "command": "bench test --package ./internal/publication",
          "exit_code": 0,
          "probe": {
            "mutation": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "d816ae02466c9c241ae8330dbb02d894fb6f0344",
              "digest": "sha256:7ed614153340dd16fc2e2dd61bf55b6a35f5a35c16653fc2d16a7382108306d4",
              "excerpt": "# preservation: necessity reds of each independent check at d816ae02466c9c241ae8330dbb02d894fb6f0344, whole package ./internal/publication, no --run, each restored\n## 3 mode drift: only TestPublicationReplacement (D24 mode check, RequireFile) fails\nold: append(data, '\\n'), 0o644)\nnew: append(data, '\\n'), 0o640)\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/publication/record.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/publication,TestPublicationReplacement,\"replacement_test.go:37: .../dist/publication/publication-record.json mode = -rw-r-----, want a regular file at -rw-r--r--\",1\n## 4 final-newline omission: only TestPublicationReplacement fails at its RequireFile call (byte step through RequireBytes)\nold: r.Replace(path, append(data, '\\n'), 0o644)\nnew: r.Replace(path, data, 0o644)\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/publication/record.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/publication,TestPublicationReplacement,\"replacement_test.go:37: .../dist/publication/publication-record.json = \\\"{\\\\n  \\\\\\\"schema_version\\\\\\\": 1, ... (747 bytes)\",1\n## 5 fallback write of the new bytes after an unpublished error: only TestPublicationPreRenameFailure fails, at its RequireBytes call (replacement_test.go:73)\nold: r.Replace(path, append(data, '\\n'), 0o644); err != nil {\nnew: r.Replace(path, append(data, '\\n'), 0o644); err != nil { _ = os.WriteFile(path, append(data, '\\n'), 0o644)\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/publication/record.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/publication,TestPublicationPreRenameFailure,\"replacement_test.go:73: .../dist/publication/publication-record.json = \\\"{\\\\n  \\\\\\\"schema_version\\\\\\\": 1, ... (757 bytes)\",1\n## 6 refusal swallow (D30): only TestPublicationDirectoryRefusal fails; before the repair this swap was silent (verdict silent, 28 tests passed at f46af6a5)\nold: return fmt.Errorf(\"could not create publication record directory: %w\", err)\nnew: _ = err\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/publication/record.go,swap,failed,1,yes\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/publication,TestPublicationDirectoryRefusal,\"replacement_test.go:90: err = could not write publication record: replace .../dist/publication/publication-record.json: create: open ... (369 bytes)\",1\n## Each probe tree: FT354-durable-callers,d816ae02466c9c241ae8330dbb02d894fb6f0344,false; each package run 29 tests; package passes after each restore.\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "dcm-c1-r1-standards",
          "performer": "claude:dcm_rev_c1_std",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "d5010a7c7ea7de4e46f377d692d2ed05a9a2adcc",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_c1_std",
            "digest": "sha256:3e9e2032fbb196ffbf50325370894f2f1eb3ed12a98959705c80339cd1f0a78b",
            "excerpt": "## Standards (D-C1, 641ddbaa..dafb4e2d) — 1 finding; worst S1\nS1 | auto-fix (orchestrator: RequireBytes route by plan expansion) | 5 | internal/publication/replacement_test.go:61-64,71-73 | the prior-bytes check copies the read-back that durablefiletest.RequireFile owns | add durablefiletest.RequireBytes, let RequireFile call it for its byte step, and use it here; log the closed copies for D-D\nClean: no copy of the algorithm; one %w wrap; saveRecordWith seam idiom; independent expectations with recorded reds; comment register.\nAdvice: two comments restate the owner contract or the code.\n"
          },
          "axis": "Standards",
          "base": "641ddbaa6c7f22023c6453abee122954564b8eca",
          "tip": "17ffac99dfbfaf0ac0d69e340c172839cffe9b29",
          "finding_ids": [
            "S1"
          ],
          "supersedes": []
        },
        {
          "id": "dcm-c1-r1-spec",
          "performer": "claude:dcm_rev_c1_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "d5010a7c7ea7de4e46f377d692d2ed05a9a2adcc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_c1_spec",
            "digest": "sha256:e7ff6c862a3cc59473be27d6f05d82d66332754c47ab2db7780d9f09cabf0cab",
            "excerpt": "## Spec (D-C1, 641ddbaa..dafb4e2d) — 0 findings\nD24 holds: MarshalIndent plus newline at 0644; parent creation stays with the caller; the state machine and release lock are unchanged; no new audited port; narrow saveRecordWith seam; TestPublicationPreRenameFailure is the D36 evidence. Plan commit 607638cd adds only the 07 assignment.\nAdvice: no run compared the former writer's bytes to the literal; best-effort SaveRecord calls discard the published signal.\n"
          },
          "axis": "Spec",
          "base": "641ddbaa6c7f22023c6453abee122954564b8eca",
          "tip": "17ffac99dfbfaf0ac0d69e340c172839cffe9b29",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "dcm-c1-r1-coverage",
          "performer": "claude:dcm_rev_c1_cov",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "d5010a7c7ea7de4e46f377d692d2ed05a9a2adcc",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_c1_cov",
            "digest": "sha256:8251519c492fa74025443a2116db03d441c9415883e8b807b357d178488bf638",
            "excerpt": "## Coverage (D-C1, 641ddbaa..dafb4e2d) — 1 finding; worst C1. Three probes, each restored=yes.\nC1 | auto-fix | 7 | internal/publication/record.go:90-92 | swallowing the MkdirAll destination refusal is silent; the replacement observer then fires (D30, D36) | add a test that makes root/dist a regular file, calls saveRecordWith, and asserts a non-durablefile error, an empty Recorder trace, and no record; record the swallow red\nBit: omit SchemaVersion (3 tests). Silent, no requirement: record directory mode 0o700.\n"
          },
          "axis": "Coverage",
          "base": "641ddbaa6c7f22023c6453abee122954564b8eca",
          "tip": "17ffac99dfbfaf0ac0d69e340c172839cffe9b29",
          "finding_ids": [
            "C1"
          ],
          "supersedes": []
        },
        {
          "id": "dcm-c1-r2-standards",
          "performer": "claude:dcm_rev_c1_std2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "2c2a7a21b89b514ed42ea5f51281c1ef521b6dd0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_c1_std2",
            "digest": "sha256:a48d1e55d48d578f59e23239922b4881a17b2d4b0720c93baa2c2d0c7e84bdae",
            "excerpt": "## Standards confirming round (D-C1, repair delta f46af6a5..41bc07d3) — 0 findings\nS1 closed: RequireBytes (assert.go:41-50) holds the only byte read-back; RequireFile calls it (assert.go:37); the pre-rename row uses it (replacement_test.go:73). TestPublicationDirectoryRefusal follows the file layout with a recorded red.\nAdvice: derive the blocking path from RecordPath rather than the \"dist\" literal.\n"
          },
          "axis": "Standards",
          "base": "641ddbaa6c7f22023c6453abee122954564b8eca",
          "tip": "d816ae02466c9c241ae8330dbb02d894fb6f0344",
          "finding_ids": [],
          "supersedes": [
            "dcm-c1-r1-standards"
          ]
        },
        {
          "id": "dcm-c1-r2-spec",
          "performer": "claude:dcm_rev_c1_spec2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "2c2a7a21b89b514ed42ea5f51281c1ef521b6dd0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_c1_spec2",
            "digest": "sha256:d859d219b621bb4b29e3d426c3e047831d61f8c918c431a783608a595417296c",
            "excerpt": "## Spec confirming round (D-C1, repair delta f46af6a5..41bc07d3 and plan 7720aa24) — 0 findings\nD24 and the D36 evidence intact. D30 met for this caller: the refusal precedes Replace and the empty Recorder trace shows the observer never ran. Plan 7720aa24 adds only assert.go to Writes and a durablefiletest floor row.\n"
          },
          "axis": "Spec",
          "base": "641ddbaa6c7f22023c6453abee122954564b8eca",
          "tip": "d816ae02466c9c241ae8330dbb02d894fb6f0344",
          "finding_ids": [],
          "supersedes": [
            "dcm-c1-r1-spec"
          ]
        },
        {
          "id": "dcm-c1-r2-coverage",
          "performer": "claude:dcm_rev_c1_cov2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "2c2a7a21b89b514ed42ea5f51281c1ef521b6dd0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_c1_cov2",
            "digest": "sha256:d2f48d40a4cc95490129de3eb6eeae348c905e379dfae348bbb590fe6123a5bd",
            "excerpt": "## Coverage confirming round (D-C1, repair delta f46af6a5..41bc07d3) — 0 findings. Four probes, each restored=yes.\nC1 closed: swallow at record.go:91 fails only TestPublicationDirectoryRefusal:90; new probe calls r.Replace then returns the refusal: bit at :93 (trace = [create]), so the trace check bites on its own.\nS1 closed: fallback write fails only at the RequireBytes call (:73); newline drop fails only at the RequireFile call (:37).\nAdvice: the no-record check at :95 is vacuous while dist is a file; delete it or reword the D30 comment.\n"
          },
          "axis": "Coverage",
          "base": "641ddbaa6c7f22023c6453abee122954564b8eca",
          "tip": "d816ae02466c9c241ae8330dbb02d894fb6f0344",
          "finding_ids": [],
          "supersedes": [
            "dcm-c1-r1-coverage"
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
      "from": "sha256:e952476f7701550400bde51f65f6db4fa0961f9d5b42c481a78e78b23e168ed5",
      "to": "sha256:4b8198b6b5417276e4d9ccfeb3ddbd2f4b4439b06c5f3a27abc84d74e6645358",
      "chunk_ids": {
        "D-B1": [
          "D-B1"
        ]
      }
    },
    {
      "from": "sha256:4b8198b6b5417276e4d9ccfeb3ddbd2f4b4439b06c5f3a27abc84d74e6645358",
      "to": "sha256:dc91be82e1b6ac107878451efa8d5156e3743e0ddef4341f9fb07a7495dfd3ae",
      "chunk_ids": {
        "D-B1": [
          "D-B1"
        ]
      }
    },
    {
      "from": "sha256:dc91be82e1b6ac107878451efa8d5156e3743e0ddef4341f9fb07a7495dfd3ae",
      "to": "sha256:304dad8ec31f6ef14776693f8e342efca2547b34f3d96073b07864b5d715b986",
      "chunk_ids": {
        "D-B1": [
          "D-B1"
        ],
        "D-B2": [
          "D-B2"
        ]
      }
    },
    {
      "from": "sha256:304dad8ec31f6ef14776693f8e342efca2547b34f3d96073b07864b5d715b986",
      "to": "sha256:46a171dbfa4f9e5b92760bbc064c364febc8bd9da34e0bd2e74698781183dd38",
      "chunk_ids": {
        "D-B1": [
          "D-B1"
        ],
        "D-B2": [
          "D-B2"
        ]
      }
    },
    {
      "from": "sha256:46a171dbfa4f9e5b92760bbc064c364febc8bd9da34e0bd2e74698781183dd38",
      "to": "sha256:92aeb66c2921bc9839cb0631f660389236b0a1cc8f5551f4743c244d38d0e6e7",
      "chunk_ids": {
        "D-B1": [
          "D-B1"
        ],
        "D-B2": [
          "D-B2"
        ]
      }
    },
    {
      "from": "sha256:92aeb66c2921bc9839cb0631f660389236b0a1cc8f5551f4743c244d38d0e6e7",
      "to": "sha256:0c585206337b1291a396cf3b404bad1b24a48b69e43d8dc10a70727f573ea631",
      "chunk_ids": {
        "D-B1": [
          "D-B1"
        ],
        "D-B2": [
          "D-B2"
        ],
        "D-B3": [
          "D-B3"
        ]
      }
    },
    {
      "from": "sha256:0c585206337b1291a396cf3b404bad1b24a48b69e43d8dc10a70727f573ea631",
      "to": "sha256:d9f4749fa22291bd7c2bec9da3823c664d5f1398f310b8f905792fb89a8fe34f",
      "chunk_ids": {
        "D-B1": [
          "D-B1"
        ],
        "D-B2": [
          "D-B2"
        ],
        "D-B3": [
          "D-B3"
        ]
      }
    },
    {
      "from": "sha256:d9f4749fa22291bd7c2bec9da3823c664d5f1398f310b8f905792fb89a8fe34f",
      "to": "sha256:fd92d1ea9c14b9c23a33447baa5f7cc1afe5b2c99dc7c09f54cc9c666808e6fc",
      "chunk_ids": {
        "D-B1": [
          "D-B1"
        ],
        "D-B2": [
          "D-B2"
        ],
        "D-B3": [
          "D-B3"
        ],
        "D-B4": [
          "D-B4"
        ]
      }
    },
    {
      "from": "sha256:fd92d1ea9c14b9c23a33447baa5f7cc1afe5b2c99dc7c09f54cc9c666808e6fc",
      "to": "sha256:908fc1e6123bcee6feb067edaa58110a42f1ba077bfcf2abe60805d1e61bc234",
      "chunk_ids": {
        "D-B1": [
          "D-B1"
        ],
        "D-B2": [
          "D-B2"
        ],
        "D-B3": [
          "D-B3"
        ],
        "D-B4": [
          "D-B4"
        ]
      }
    },
    {
      "from": "sha256:908fc1e6123bcee6feb067edaa58110a42f1ba077bfcf2abe60805d1e61bc234",
      "to": "sha256:1f71fb49502390cf5c1ccabedb96986f67bf4e6e59abffc47a1a0f1e0f25de2b",
      "chunk_ids": {
        "D-B1": [
          "D-B1"
        ],
        "D-B2": [
          "D-B2"
        ],
        "D-B3": [
          "D-B3"
        ],
        "D-B4": [
          "D-B4"
        ]
      }
    },
    {
      "from": "sha256:1f71fb49502390cf5c1ccabedb96986f67bf4e6e59abffc47a1a0f1e0f25de2b",
      "to": "sha256:1328d3cecb091573d96d49938b9049c6cc47592195da6220c313c0c4dbef2abd",
      "chunk_ids": {
        "D-B1": [
          "D-B1"
        ],
        "D-B2": [
          "D-B2"
        ],
        "D-B3": [
          "D-B3"
        ],
        "D-B4": [
          "D-B4"
        ],
        "D-B5": [
          "D-B5"
        ]
      }
    },
    {
      "from": "sha256:1328d3cecb091573d96d49938b9049c6cc47592195da6220c313c0c4dbef2abd",
      "to": "sha256:e6d4e71219c6d5808b3d70faeec627fdabed0455e2f06870b25523a882c731d7",
      "chunk_ids": {
        "D-B1": [
          "D-B1"
        ],
        "D-B2": [
          "D-B2"
        ],
        "D-B3": [
          "D-B3"
        ],
        "D-B4": [
          "D-B4"
        ],
        "D-B5": [
          "D-B5"
        ]
      }
    },
    {
      "from": "sha256:e6d4e71219c6d5808b3d70faeec627fdabed0455e2f06870b25523a882c731d7",
      "to": "sha256:95a4fa23af538a1a3b21d7b5f8d34d8497b07fa508e4ad0dad253db7925d183c",
      "chunk_ids": {
        "D-B1": [
          "D-B1"
        ],
        "D-B2": [
          "D-B2"
        ],
        "D-B3": [
          "D-B3"
        ],
        "D-B4": [
          "D-B4"
        ],
        "D-B5": [
          "D-B5"
        ],
        "D-C1": [
          "D-C1"
        ]
      }
    },
    {
      "from": "sha256:95a4fa23af538a1a3b21d7b5f8d34d8497b07fa508e4ad0dad253db7925d183c",
      "to": "sha256:5cb50643a5295119fe76284a8cc58c0cc7bfd797ae01bfecc20e1b30a5df98af",
      "chunk_ids": {
        "D-B1": [
          "D-B1"
        ],
        "D-B2": [
          "D-B2"
        ],
        "D-B3": [
          "D-B3"
        ],
        "D-B4": [
          "D-B4"
        ],
        "D-B5": [
          "D-B5"
        ],
        "D-C1": [
          "D-C1"
        ]
      }
    }
  ],
  "folds": [
    {
      "id": "F1",
      "after": "D-B4",
      "from": "5aabee17c744adc872eca23dc08fe1a67c46c0b0",
      "commit": "51709c6942f766ae6e2c2af36e7ae96c92a2592a",
      "incoming": "91fde7529c228c0b2276530d5864e5157df169c6",
      "source_digest": "f4f96ea11906fa29f6ba8c088b237d1f4442b35d",
      "resolution": [],
      "resolution_digest": "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
      "test_packages": [],
      "reviews": []
    }
  ]
}
```

## Standards

Chunk D-B1 has 1 finding. The worst issue is S1.

- S1, auto-fix, confidence 8: `internal/gate/verdict.go:95` repeats the encoding and mode of `durableReplaceAt` in `internal/gate/engine.go:173-187`. The comments at `internal/gate/engine.go:182` and `internal/gate/lane.go:317` are false for the verdict cache. Make `durableReplace` call `durableReplaceAt`, and keep the preservation red of `TestVerdictReplacement`.

## Spec

Chunk D-B1 has 0 findings.

## Coverage

Chunk D-B1 has 3 findings. The worst issue is C1.

- C1, auto-fix, confidence 9: no test fails the timeout persistence branch at `internal/gate/run_transaction.go:186-187`. A swap that discards its cause stays silent. Add an execution witness for the timeout cause, the exits 124 and 1, and the pending record after recovery.
- C2, auto-fix, confidence 9: no test fails the lock-contention demotion caller at `internal/gate/run_transaction.go:83`. A swap that discards its error stays silent. Add a contention witness from a reusable green with a demotion fault.
- C3, auto-fix, confidence 7: four independent expectations and the owner-bypass mutation have no recorded red. They are the 0600 mode literal, the D39 cache count, the D40 visible-green inspection, and the bypass of each writer. Record each red with its exact swap text.

## Repair allowance

Chunk D-B1 used 0 of its 2 repair cycles before repair 1. S1, C1, C2, and C3 fall on ticket 02, so one combined repair session takes them.

- S1 closes in commit e886531e. The verdict writer now calls the store writer, and the 0600 swap in the store writer turns `TestVerdictReplacement` red.
- C1 closes in commit e886531e. `TestGateTerminalDurabilityFailure/timeout` turns red when the timeout branch drops its cause.
- C2 closes in commit e886531e. `TestGateDemotionDurabilityFailure/lock_contention` turns red when the contention caller drops its demotion error.
- C3 closes in commit e886531e. The preservation and omission-proof excerpts of `claude:dcm_t02_rep1` record the exact swap and the red for each expectation.

Chunk D-B1 used 1 of its 2 repair cycles.

The D-B1 confirming round at 9c9ae1c0 closed S1, C1, C2, and C3 with no new finding. Chunk D-B1 closed with 1 of its 2 repair cycles used.

## Standards (D-B2)

Chunk D-B2 has 1 finding. The worst issue is DCM-B2-S1.

- DCM-B2-S1, auto-fix by the --auto-approve plan expansion, confidence 6: `internal/gate/prospectiveartifact/replacement_test.go:22-41` copies the published-failure and record-file checks of `internal/gate/durable_replacement_test.go:292-312`. Move one exported assertion into `internal/durablefile/durablefiletest`, and let both caller tests use it.

## Spec (D-B2)

Chunk D-B2 has 0 findings.

## Coverage (D-B2)

Chunk D-B2 has 3 findings. The worst issue is C1.

- C1, auto-fix, confidence 8: no test fails the pre-rename branch of `publish` at `internal/gate/prospectiveartifact/prospectiveartifact.go:240`. Add a pre-rename row that asserts an unpublished classified error, the cause, and an absent record, and record its bypass red.
- C2, no-op, confidence 6: no caller test detects a plain write in exported `Publish`, and the seam accepts that limit. The D-D source census confirms the native binding.
- C3, auto-fix, confidence 7: the mode red ran only one test. Record the whole-package red of the `RecordMode` constant drift to 0o640.

## Repair allowance (D-B2)

Chunk D-B2 used 0 of its 2 repair cycles before repair 1. DCM-B2-S1, C1, and C3 fall on ticket 03, so one combined repair session takes them.

- DCM-B2-S1 closes in commit 3a3dbbe1. `durablefiletest.RequireFailure` and `durablefiletest.RequireFile` own the checks, and `TestVerdictReplacement`, `TestEvidenceReplacement`, and `TestOwnerReplacement` use them with no local copy.
- C1 closes in commit 3a3dbbe1. `TestOwnerPreRenameFailure` turns red when `publish` returns nil for an unpublished replacement failure.
- C3 closes in commit 3a3dbbe1. The `RecordMode` drift to 0o640 turns `TestOwnerReplacement` red over the whole package, as the `claude:dcm_t03_rep1` preservation excerpt records.

Chunk D-B2 used 1 of its 2 repair cycles.

The D-B2 confirming round at f0c35260 closed DCM-B2-S1, C1, and C3. It raised one new finding.

- DCM-B2-S2, auto-fix, confidence 5: `internal/gate/prospectiveartifact/replacement_test.go:42-44` asserts the temporary-file cleanup that the leaf test `TestReplaceTemporaryCleanup` owns, and no caller mutation shows that the check is necessary. Remove the sibling check, and keep the record-absence check.

Chunk D-B2 takes repair cycle 2 of 2 for DCM-B2-S2.

DCM-B2-S2 closes in commit 6eed4ae5, and chunk D-B2 used 2 of its 2 repair cycles.

The second D-B2 confirming round at 79926b97 closed DCM-B2-S2 with no new finding. Chunk D-B2 closed with 2 of its 2 repair cycles used.

## Standards (D-B3)

Chunk D-B3 has 3 findings. The worst issue is S1.

- S1, auto-fix in part, confidence 5: the pre-rename row at `internal/capturetx/replacement_test.go:68` repeats the 0644 check that line 47 owns. Drop the repeated mode check, and keep the former-bytes check with a recorded necessity red.
- S2, auto-fix, confidence 6: the comments at `internal/capturetx/transaction.go:208`, `:245-246`, and the test line 12 restate the code. Keep only the reason for the parent creation.
- S3, no-op, confidence 4: spec-row tags in test doc comments are an established tree convention.

## Spec (D-B3)

Chunk D-B3 has 0 findings.

## Coverage (D-B3)

Chunk D-B3 has 3 findings. The worst issue is C1.

- C1, no-op, confidence 6: a bypass of the native bindings in `replace` or `writeManifestDir` is silent. The spec forbids ambient swaps, so the D-D source census owns each native binding.
- C2, no-op, confidence 6: an Append that discards its replacement error is silent. The spec accepts a narrow caller writer, and Append sits outside the fence.
- C3, auto-fix, confidence 7: the MkdirAll omission red is not recorded. Record it with its exact swap text.

## Repair allowance (D-B3)

Chunk D-B3 used 0 of its 2 repair cycles before repair 1. S1, S2, and C3 fall on ticket 04, so one combined repair session takes them.

- S1 closes in commit 3f43adb7. `TestCapturePreRenameFailure` keeps only the former-bytes check. A fallback write of the new bytes after an unpublished error turns it red over the whole package. The `claude:dcm_t04_rep1` preservation excerpt records this red.
- S2 closes in commit 3f43adb7. The restating comments are gone, and the `replaceWith` comment states only the reason for the parent creation.
- C3 closes in commit 3f43adb7. The `os.MkdirAll` omission in `replaceWith` turns `TestCaptureReplacement` red over the whole package, as the `claude:dcm_t04_rep1` preservation excerpt records with its exact omitted text.

Chunk D-B3 used 1 of its 2 repair cycles.

The D-B3 confirming round at 13925a7d closed S1, S2, and C3 with no new finding. Chunk D-B3 closed with 1 of its 2 repair cycles used.

## Standards (D-B4)

Chunk D-B4 has 1 finding. The worst issue is S1.

- S1, auto-fix, confidence 7: the comment at `internal/handoffdoc/replacement_test.go:19-20` narrates history. Delete the clause about the former writer.

## Spec (D-B4)

Chunk D-B4 has 0 findings.

## Coverage (D-B4)

Chunk D-B4 has 2 findings. The worst issue is F1.

- F1, auto-fix, confidence 7: the prior-bytes check of `TestHandoffPreRenameFailure` has no recorded necessity red. Record the fallback-write red with its exact swap text.
- F2, auto-fix, confidence 4: only directory-sync has a caller witness, and an ignored directory-open error is silent. Make `TestHandoffReplacement` table-driven over directory-open, directory-sync, and directory-close.

## Repair allowance (D-B4)

Chunk D-B4 used 0 of its 2 repair cycles before repair 1. S1, F1, and F2 fall on ticket 05, so one combined repair session takes them.

- S1 closes in commit 5e6c8994. The `TestHandoffReplacement` comment states only the current behavior and names no former writer.
- F1 closes in commit 5e6c8994. A fallback write of the new bytes after an unpublished error turns `TestHandoffPreRenameFailure` red over the whole package. The red falls only at its prior-bytes check. The `claude:dcm_t05_rep1` preservation excerpt records this red with its exact swap text.
- F2 closes in commit 5e6c8994. `TestHandoffReplacement` is table-driven over directory-open, directory-sync, and directory-close. An ignored error at each stage turns its own row red over the whole package. The `claude:dcm_t05_rep1` omission-proof excerpt records each red.

Chunk D-B4 used 1 of its 2 repair cycles.

The D-B4 confirming round at 92da1312 closed S1, F1, and F2. It raised one new finding.

- S2, auto-fix by the --auto-approve plan expansion, confidence 6: `internal/handoffdoc/replacement_test.go:22` writes out the post-rename stage set that the leaf owns. Export the set once from `internal/durablefile/durablefiletest`, and range over it.

Chunk D-B4 takes repair cycle 2 of 2 for S2.

S2 closes in commit 5e49631e, where `TestHandoffReplacement` ranges over `durablefiletest.PostRenameStages`, and chunk D-B4 used 2 of its 2 repair cycles.

The second D-B4 confirming round at 2477c580 closed S2 with no new finding. Chunk D-B4 closed with 2 of its 2 repair cycles used.

## Standards (D-B5)

Chunk D-B5 has 0 findings.

## Spec (D-B5)

Chunk D-B5 has 1 finding. The worst issue is P1.

- P1, auto-fix, confidence 6: the D-B5 record holds no D22 owner-bypass red with its exact swap text. Record a probe that swaps the owner call at `internal/intent/intent.go:216` for the former rename sequence or the native owner.

## Coverage (D-B5)

Chunk D-B5 has 2 findings. The worst issue is C1.

- C1, auto-fix, confidence 9: no rename-stage or directory-open or directory-close fault reaches `TransactWithReplacer`, so three swaps of the `published` predicate at `internal/intent/transaction.go:88-91` stay silent. Range `PostRenameStages` through the transaction with no compensation, inject a rename fault with one compensation, and record the three swaps as reds.
- C2, auto-fix, confidence 8: releasing the lock before compensation at `internal/intent/transaction.go:78-79` stays silent. Assert inside the compensation that the lock exists, assert after the call that it is gone, and record that swap as its red.

## Repair allowance (D-B5)

Chunk D-B5 used 0 of its 2 repair cycles before repair 1. P1, C1, and C2 fall on ticket 06, so one combined repair session takes them.

- P1 closes in commit 19422ae3. The swap of the owner call in `writePathWith` for `durablefile.Native().Replace` turns each `TestIntentReplacement` row red at its failure check over the whole package. The `claude:dcm_t06_rep1` preservation excerpt records this red with its exact swap text.
- C1 closes in commit 19422ae3. `TestTransactionPublishedFailure` drives a rename fault and each `durablefiletest.PostRenameStages` fault through `TransactWithReplacer`. Each of the three predicate swaps turns a row red over the whole package. The `claude:dcm_t06_rep1` omission-proof excerpt records each red.
- C2 closes in commit 19422ae3. The rename row of `TestTransactionPublishedFailure` checks the lock inside the compensation and its absence after the call. A release before the compensation turns that row red. The `claude:dcm_t06_rep1` omission-proof excerpt records this red.

Chunk D-B5 used 1 of its 2 repair cycles.

The D-B5 confirming round at d12060c6 closed P1, C1, and C2 with no new finding. Chunk D-B5 closed with 1 of its 2 repair cycles used.

## Standards (D-C1)

Chunk D-C1 has 1 finding. The worst issue is S1.

- S1, auto-fix by the --auto-approve plan expansion, confidence 5: the prior-bytes check at `internal/publication/replacement_test.go:61-73` copies the read-back that `durablefiletest.RequireFile` owns. Add `durablefiletest.RequireBytes`, let `RequireFile` call it for its byte step, and use it here.

## Spec (D-C1)

Chunk D-C1 has 0 findings.

## Coverage (D-C1)

Chunk D-C1 has 1 finding. The worst issue is C1.

- C1, auto-fix, confidence 7: swallowing the parent-creation refusal at `internal/publication/record.go:90-92` is silent, and the replacement observer then fires. Add the D30 refusal witness, and record the red of the swallow.

## Repair allowance (D-C1)

Chunk D-C1 used 0 of its 2 repair cycles before repair 1. S1 and C1 fall on ticket 08, so one combined repair session takes them.

- S1 closes in commit d816ae02. `durablefiletest.RequireBytes` owns the byte read-back, `RequireFile` calls it, and `TestPublicationPreRenameFailure` checks the prior record through it. The fallback-write swap turns that test red only at its `RequireBytes` call, and the final-newline swap turns `TestPublicationReplacement` red at its `RequireFile` call. The `claude:dcm_t08_rep1` preservation excerpt records both reds.
- C1 closes in commit d816ae02. `TestPublicationDirectoryRefusal` makes the record directory uncreatable and checks a refusal that is not a `*durablefile.Error`, an empty Recorder trace, and no record. The swap of the parent-creation return for `_ = err` was silent before the repair and turns that test red after it. The `claude:dcm_t08_rep1` preservation excerpt records this red.

Chunk D-C1 used 1 of its 2 repair cycles.

The D-C1 confirming round at 41bc07d3 closed S1 and C1 with no new finding. Chunk D-C1 closed with 1 of its 2 repair cycles used.
