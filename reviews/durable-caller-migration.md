# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/durable-caller-migration/spec.md",
  "plan_digest": "sha256:0a9c82011a15380bda3c64cc4ee87594038caa6080e5ec3d4f7ff5bef371daf6",
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
    },
    {
      "id": "D-B6",
      "base": "3cf04078eb7bf8169f1af53cde9074f0a1da206e",
      "tip": "91954b3279d3ab379a62f92630c0be55b62e0421",
      "plan_digest": "sha256:b46bfce6efc8865c4e196a7c89854bdcb94e436d4d1b1c1b5bd58a67d573c3b2",
      "source_digest": "3d38043c20090c1b33f3011d740349e61b132497",
      "acceptance_rows": [
        "D23",
        "D46",
        "D47",
        "D48",
        "D49",
        "D50",
        "D51"
      ],
      "verification": [
        {
          "id": "dcm-t07-repository-tests",
          "performer": "claude:dcm_t07",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6db04c672387f77749e07f16235df9cbe9e9566a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "67d8443bb5fd5574b6496be39fb9f7668632664e",
            "digest": "sha256:6440fabb3d85da00ca076778d4631c148bdbd52489ce38cf86175927901b503f",
            "excerpt": "# bench test --package ./internal/commitment/repository\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,pass,5095,143\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "repository-tests",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0
        },
        {
          "id": "dcm-t07-approval-tests",
          "performer": "claude:dcm_t07",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6db04c672387f77749e07f16235df9cbe9e9566a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "67d8443bb5fd5574b6496be39fb9f7668632664e",
            "digest": "sha256:6a9237446b2588ca9ea33567c00ec3b191fd4a021730c20b5a18e12d96b3c680",
            "excerpt": "# bench test --package ./internal/commitment\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment,pass,4322,175\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "approval-tests",
          "command": "bench test --package ./internal/commitment",
          "exit_code": 0
        },
        {
          "id": "dcm-t07-command-contract",
          "performer": "claude:dcm_t07",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6db04c672387f77749e07f16235df9cbe9e9566a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "67d8443bb5fd5574b6496be39fb9f7668632664e",
            "digest": "sha256:68dd62a5b1877b9af51df8397a857e16db528c89e53bb77977255b6a89b77a4d",
            "excerpt": "# bench test --package ./cmd/bench --run CommandRegistry\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,4,5\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "command-contract",
          "command": "bench test --package ./cmd/bench --run 'CommandRegistry'",
          "exit_code": 0
        },
        {
          "id": "dcm-t07-axi-registry",
          "performer": "claude:dcm_t07",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6db04c672387f77749e07f16235df9cbe9e9566a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "67d8443bb5fd5574b6496be39fb9f7668632664e",
            "digest": "sha256:7b949ee851a8f8107325becd4da7c44bbc81fae79444dc46bafc1e2e267ceab1",
            "excerpt": "# bench test --check axi-query-registry\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\ncheck[1]{name,kind,tests_run,subjects}:\n  axi-query-registry,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,10,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "axi-registry",
          "command": "bench test --check axi-query-registry",
          "exit_code": 0
        },
        {
          "id": "dcm-t07-routing",
          "performer": "claude:dcm_t07",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6db04c672387f77749e07f16235df9cbe9e9566a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "67d8443bb5fd5574b6496be39fb9f7668632664e",
            "digest": "sha256:15f96663f606691aa55d68456ca019df1dece1f29f54fb464ccfa275707db5de",
            "excerpt": "# bench test --check subcommand-routing\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\ncheck[1]{name,kind,tests_run,subjects}:\n  subcommand-routing,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,24,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "routing",
          "command": "bench test --check subcommand-routing",
          "exit_code": 0
        },
        {
          "id": "dcm-t07-cmd-bench",
          "performer": "claude:dcm_t07",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6db04c672387f77749e07f16235df9cbe9e9566a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "67d8443bb5fd5574b6496be39fb9f7668632664e",
            "digest": "sha256:a9c3f487c61d0df2481eb93bade6a43dca2bbafd9440fb1556935fdd4e00d6d5",
            "excerpt": "# bench test --package ./cmd/bench\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,13851,335\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "dcm-t07-conformance",
          "performer": "claude:dcm_t07",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6db04c672387f77749e07f16235df9cbe9e9566a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "67d8443bb5fd5574b6496be39fb9f7668632664e",
            "digest": "sha256:fdcbb3718840b70649483633acd3d8fe4c8977584fe90d8fe640b01e642e2f4c",
            "excerpt": "# bench test --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,53988,1690\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "dcm-t07-root-conformance",
          "performer": "claude:dcm_t07",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6db04c672387f77749e07f16235df9cbe9e9566a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "67d8443bb5fd5574b6496be39fb9f7668632664e",
            "digest": "sha256:1aeb19e3772f0e5b05f79cfa3a8b0eba9b88c26c9acec022ff1af23293eccf97",
            "excerpt": "# bench test --check root-conformance\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14751,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "dcm-t07-omission-proof",
          "performer": "claude:dcm_t07",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6db04c672387f77749e07f16235df9cbe9e9566a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "67d8443bb5fd5574b6496be39fb9f7668632664e",
            "digest": "sha256:75fdc3c53256a298fe5fdb5cb0757673bfa639ebf6c53716b340f60edd3884f8",
            "excerpt": "# omission-proof: the four named mutations plus owner bypass and lock placement, bench probe at 67d8443bb5fd5574b6496be39fb9f7668632664e, whole package ./internal/commitment/repository, no --run, each restored=yes\n## 1 restore built only after the first policy write (central self-probe). subject internal/commitment/repository/replacement.go\nold: \\trestore := func() error { return restoreAll(replacer, policyBefore, roadmapBefore) }\\n\\tif err := replace(replacer, policyPath, plan.Proposed, 0o644); err != nil {\\n\\t\\treturn nil, errors.Join(err, restore())\\n\\t}\nnew: \\tif err := replace(replacer, policyPath, plan.Proposed, 0o644); err != nil {\\n\\t\\treturn nil, err\\n\\t}\\n\\trestore := func() error { return restoreAll(replacer, policyBefore, roadmapBefore) }\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=13,bytes=1589,omitted_lines=4,cut_lines=1,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835367839543955-cdf7aade2b6e5d3f.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalFailureWithoutBoard/policy#1/directory-sync,\"replacement_test.go:390: /dev/shm/64EFQS/t/TestApprovalFailureWithoutBoardpolicy#1directory-sync1232576004/001/.bench/commitment.json = \\\"{\\\\\\\\n  \\\\\\\\\\\"version\\\\\\\\\\\": 1,\\\\\\\\n  \\\\\\\\\\\"milestones\\\\\\\\\\\": [\\\\\\\\n    {\\\\\\\\n      \\\\\\\\\\\"id\\\\\\\\\\\": \\\\\\\\\\\"M1\\\\\\\\\\\",\\\\\\\\n      \\\\\\\\\\\"outcomes\\\\\\\\\\\": [\\\\\\\\n \n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalInitialPolicyFailure,\"replacement_test.go:226: /dev/shm/64EFQS/t/TestApprovalInitialPolicyFailure2586894039/001/.bench/commitment.json mode = -rw-r--r--, want a regular file at -rw-------\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/stage,\"replacement_test.go:358: err = replace /dev/shm/64EFQS/t/TestApprovalRestoreFailurestage2092015824/001/.bench/commitment.json: directory-sync: injected policy#1 failure; the rename completed, so new bytes may already be visible, want it to … (270 bytes)\",1\nskips[0]{package,test,reason}:\n## 2 compensate a published ledger. subject internal/commitment/repository/repository.go\nold: \\treturn changed, errors.Join(err, restoreErr)\nnew: \\tif err != nil && restore != nil {\\n\\t\\trestoreErr = errors.Join(restoreErr, restore())\\n\\t}\\n\\treturn changed, errors.Join(err, restoreErr)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/repository.go,swap,failed,7,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=17,bytes=2721,omitted_lines=8,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835383403531896-4f2b9b33e3fa7e2e.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-open,\"replacement_test.go:122: policy#2 ran without the ledger lock: stat /dev/shm/FCOCUH/t/TestApprovalPublishedLedgerFailuredirectory-open3585750263/001/.git/bench-intent.json.lock: no such file or directory\",3\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-sync,\"replacement_test.go:122: policy#2 ran without the ledger lock: stat /dev/shm/FCOCUH/t/TestApprovalPublishedLedgerFailuredirectory-sync4001034437/001/.git/bench-intent.json.lock: no such file or directory\",3\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/compensation,\"replacement_test.go:122: policy#3 ran without the ledger lock: stat /dev/shm/FCOCUH/t/TestApprovalRestoreFailurecompensation1888372961/001/.git/bench-intent.json.lock: no such file or directory\",2\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalUnpublishedLedgerFailure,\"replacement_test.go:122: policy#3 ran without the ledger lock: stat /dev/shm/FCOCUH/t/TestApprovalUnpublishedLedgerFailure3350151100/001/.git/bench-intent.json.lock: no such file or directory\",3\nskips[0]{package,test,reason}:\n## 3 restore only one file. subject internal/commitment/repository/replacement.go\nold: restoreAll(replacer, policyBefore, roadmapBefore)\nnew: restoreAll(replacer, policyBefore)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,4,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=14,bytes=1769,omitted_lines=5,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835401270286860-0b9b03c650804701.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalBoardFailure,\"replacement_test.go:243: /dev/shm/VI4OE4/t/TestApprovalBoardFailure2914500064/001/ROADMAP.md = \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n2. B\\\\\\\\n\\\", want \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. old\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/compensation,\"replacement_test.go:358: err = replace /dev/shm/VI4OE4/t/TestApprovalRestoreFailurecompensation2816159136/001/.git/bench-intent.json: rename: injected ledger#1 failure, want it to hold injected board#2 failure\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back,\"replacement_test.go:358: err = replace /dev/shm/VI4OE4/t/TestApprovalRestoreFailureread-back3103315451/001/ROADMAP.md: directory-sync: injected board#1 failure; the rename completed, so new bytes may already be visible, want it to name the … (329 bytes)\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalUnpublishedLedgerFailure,\"replacement_test.go:258: replacements = [policy#1 board#1 ledger#1 policy#2], want [policy#1 board#1 ledger#1 policy#2 board#2]\",1\nskips[0]{package,test,reason}:\n## 4a discard restore failure in restoreAll. subject internal/commitment/repository/replacement.go\nold: \\treturn errors.Join(errs...)\nnew: \\treturn nil\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=13,bytes=1561,omitted_lines=4,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835416102453931-790acf5652ddb0b6.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/compensation,\"replacement_test.go:358: err = replace /dev/shm/RKYJLO/t/TestApprovalRestoreFailurecompensation2517643027/001/.git/bench-intent.json: rename: injected ledger#1 failure, want it to hold injected board#2 failure\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back,\"replacement_test.go:358: err = replace /dev/shm/RKYJLO/t/TestApprovalRestoreFailureread-back578465981/001/ROADMAP.md: directory-sync: injected board#1 failure; the rename completed, so new bytes may already be visible, want it to name the u… (327 bytes)\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/stage,\"replacement_test.go:358: err = replace /dev/shm/RKYJLO/t/TestApprovalRestoreFailurestage1335700746/001/.bench/commitment.json: directory-sync: injected policy#1 failure; the rename completed, so new bytes may already be visible, want it to … (270 bytes)\",1\nskips[0]{package,test,reason}:\n## 4b discard the compensation restore failure in approve. subject internal/commitment/repository/repository.go\nold: \\treturn changed, errors.Join(err, restoreErr)\nnew: \\t_ = restoreErr\\n\\treturn changed, err\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/repository.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=11,bytes=851,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835632645647061-becaedd1adc75988.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5098,143\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/compensation,\"replacement_test.go:358: err = replace /dev/shm/EV3YK2/t/TestApprovalRestoreFailurecompensation3796896505/001/.git/bench-intent.json: rename: injected ledger#1 failure, want it to hold injected board#2 failure\",1\nskips[0]{package,test,reason}:\n## 5 native owner in place of the injected replacer. subject internal/commitment/repository/replacement.go\nold: \\treturn replacer.Replace(path, data, mode)\nnew: \\treturn durablefile.Native().Replace(path, data, mode)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,11,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=21,bytes=2975,omitted_lines=12,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835433269905653-b9c0c3522f2f4064.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCommitmentReplacement/directory-close,\"replacement_test.go:205: err = <nil>, want a *durablefile.Error at stage directory-close with Published=true that wraps injected policy directory-close failure\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCommitmentReplacement/directory-open,\"replacement_test.go:205: err = <nil>, want a *durablefile.Error at stage directory-open with Published=true that wraps injected policy directory-open failure\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCommitmentReplacement/directory-sync,\"replacement_test.go:205: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected policy directory-sync failure\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCommitmentReplacement/rename,\"replacement_test.go:205: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected policy rename failure\",1\nskips[0]{package,test,reason}:\n## 6 restore after the transaction released the lock. subject internal/commitment/repository/repository.go\nold: \\t}, func() {\\n\\t\\tif restore != nil {\\n\\t\\t\\trestoreErr = restore()\\n\\t\\t}\\n\\t})\nnew: \\t}, func() {\\n\\t\\trestoreErr = errors.ErrUnsupported\\n\\t})\\n\\tif errors.Is(restoreErr, errors.ErrUnsupported) {\\n\\t\\trestoreErr = restore()\\n\\t}\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/repository.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=13,bytes=1440,omitted_lines=4,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835648442554770-380b1d2603646147.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalFailureWithoutBoard/ledger#1/rename,\"replacement_test.go:122: policy#2 ran without the ledger lock: stat /dev/shm/422XEB/t/TestApprovalFailureWithoutBoardledger#1rename2293033941/001/.git/bench-intent.json.lock: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/compensation,\"replacement_test.go:122: policy#2 ran without the ledger lock: stat /dev/shm/422XEB/t/TestApprovalRestoreFailurecompensation3260744936/001/.git/bench-intent.json.lock: no such file or directory\",2\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalUnpublishedLedgerFailure,\"replacement_test.go:122: policy#2 ran without the ledger lock: stat /dev/shm/422XEB/t/TestApprovalUnpublishedLedgerFailure2608250465/001/.git/bench-intent.json.lock: no such file or directory\",2\nskips[0]{package,test,reason}:\n## Package passes after each restore (bench test --package ./internal/commitment/repository: pass, 143 tests).\n"
          },
          "requirement": "omission-proof",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0,
          "probe": {
            "mutation": "Create restore only after policy write, compensate a published ledger, restore only one file, or discard restore failure. Each mutation must fail the corresponding real approval witness.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "67d8443bb5fd5574b6496be39fb9f7668632664e",
              "digest": "sha256:75fdc3c53256a298fe5fdb5cb0757673bfa639ebf6c53716b340f60edd3884f8",
              "excerpt": "# omission-proof: the four named mutations plus owner bypass and lock placement, bench probe at 67d8443bb5fd5574b6496be39fb9f7668632664e, whole package ./internal/commitment/repository, no --run, each restored=yes\n## 1 restore built only after the first policy write (central self-probe). subject internal/commitment/repository/replacement.go\nold: \\trestore := func() error { return restoreAll(replacer, policyBefore, roadmapBefore) }\\n\\tif err := replace(replacer, policyPath, plan.Proposed, 0o644); err != nil {\\n\\t\\treturn nil, errors.Join(err, restore())\\n\\t}\nnew: \\tif err := replace(replacer, policyPath, plan.Proposed, 0o644); err != nil {\\n\\t\\treturn nil, err\\n\\t}\\n\\trestore := func() error { return restoreAll(replacer, policyBefore, roadmapBefore) }\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=13,bytes=1589,omitted_lines=4,cut_lines=1,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835367839543955-cdf7aade2b6e5d3f.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalFailureWithoutBoard/policy#1/directory-sync,\"replacement_test.go:390: /dev/shm/64EFQS/t/TestApprovalFailureWithoutBoardpolicy#1directory-sync1232576004/001/.bench/commitment.json = \\\"{\\\\\\\\n  \\\\\\\\\\\"version\\\\\\\\\\\": 1,\\\\\\\\n  \\\\\\\\\\\"milestones\\\\\\\\\\\": [\\\\\\\\n    {\\\\\\\\n      \\\\\\\\\\\"id\\\\\\\\\\\": \\\\\\\\\\\"M1\\\\\\\\\\\",\\\\\\\\n      \\\\\\\\\\\"outcomes\\\\\\\\\\\": [\\\\\\\\n \n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalInitialPolicyFailure,\"replacement_test.go:226: /dev/shm/64EFQS/t/TestApprovalInitialPolicyFailure2586894039/001/.bench/commitment.json mode = -rw-r--r--, want a regular file at -rw-------\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/stage,\"replacement_test.go:358: err = replace /dev/shm/64EFQS/t/TestApprovalRestoreFailurestage2092015824/001/.bench/commitment.json: directory-sync: injected policy#1 failure; the rename completed, so new bytes may already be visible, want it to … (270 bytes)\",1\nskips[0]{package,test,reason}:\n## 2 compensate a published ledger. subject internal/commitment/repository/repository.go\nold: \\treturn changed, errors.Join(err, restoreErr)\nnew: \\tif err != nil && restore != nil {\\n\\t\\trestoreErr = errors.Join(restoreErr, restore())\\n\\t}\\n\\treturn changed, errors.Join(err, restoreErr)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/repository.go,swap,failed,7,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=17,bytes=2721,omitted_lines=8,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835383403531896-4f2b9b33e3fa7e2e.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-open,\"replacement_test.go:122: policy#2 ran without the ledger lock: stat /dev/shm/FCOCUH/t/TestApprovalPublishedLedgerFailuredirectory-open3585750263/001/.git/bench-intent.json.lock: no such file or directory\",3\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-sync,\"replacement_test.go:122: policy#2 ran without the ledger lock: stat /dev/shm/FCOCUH/t/TestApprovalPublishedLedgerFailuredirectory-sync4001034437/001/.git/bench-intent.json.lock: no such file or directory\",3\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/compensation,\"replacement_test.go:122: policy#3 ran without the ledger lock: stat /dev/shm/FCOCUH/t/TestApprovalRestoreFailurecompensation1888372961/001/.git/bench-intent.json.lock: no such file or directory\",2\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalUnpublishedLedgerFailure,\"replacement_test.go:122: policy#3 ran without the ledger lock: stat /dev/shm/FCOCUH/t/TestApprovalUnpublishedLedgerFailure3350151100/001/.git/bench-intent.json.lock: no such file or directory\",3\nskips[0]{package,test,reason}:\n## 3 restore only one file. subject internal/commitment/repository/replacement.go\nold: restoreAll(replacer, policyBefore, roadmapBefore)\nnew: restoreAll(replacer, policyBefore)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,4,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=14,bytes=1769,omitted_lines=5,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835401270286860-0b9b03c650804701.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalBoardFailure,\"replacement_test.go:243: /dev/shm/VI4OE4/t/TestApprovalBoardFailure2914500064/001/ROADMAP.md = \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n2. B\\\\\\\\n\\\", want \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. old\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/compensation,\"replacement_test.go:358: err = replace /dev/shm/VI4OE4/t/TestApprovalRestoreFailurecompensation2816159136/001/.git/bench-intent.json: rename: injected ledger#1 failure, want it to hold injected board#2 failure\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back,\"replacement_test.go:358: err = replace /dev/shm/VI4OE4/t/TestApprovalRestoreFailureread-back3103315451/001/ROADMAP.md: directory-sync: injected board#1 failure; the rename completed, so new bytes may already be visible, want it to name the … (329 bytes)\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalUnpublishedLedgerFailure,\"replacement_test.go:258: replacements = [policy#1 board#1 ledger#1 policy#2], want [policy#1 board#1 ledger#1 policy#2 board#2]\",1\nskips[0]{package,test,reason}:\n## 4a discard restore failure in restoreAll. subject internal/commitment/repository/replacement.go\nold: \\treturn errors.Join(errs...)\nnew: \\treturn nil\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=13,bytes=1561,omitted_lines=4,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835416102453931-790acf5652ddb0b6.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/compensation,\"replacement_test.go:358: err = replace /dev/shm/RKYJLO/t/TestApprovalRestoreFailurecompensation2517643027/001/.git/bench-intent.json: rename: injected ledger#1 failure, want it to hold injected board#2 failure\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back,\"replacement_test.go:358: err = replace /dev/shm/RKYJLO/t/TestApprovalRestoreFailureread-back578465981/001/ROADMAP.md: directory-sync: injected board#1 failure; the rename completed, so new bytes may already be visible, want it to name the u… (327 bytes)\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/stage,\"replacement_test.go:358: err = replace /dev/shm/RKYJLO/t/TestApprovalRestoreFailurestage1335700746/001/.bench/commitment.json: directory-sync: injected policy#1 failure; the rename completed, so new bytes may already be visible, want it to … (270 bytes)\",1\nskips[0]{package,test,reason}:\n## 4b discard the compensation restore failure in approve. subject internal/commitment/repository/repository.go\nold: \\treturn changed, errors.Join(err, restoreErr)\nnew: \\t_ = restoreErr\\n\\treturn changed, err\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/repository.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=11,bytes=851,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835632645647061-becaedd1adc75988.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5098,143\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/compensation,\"replacement_test.go:358: err = replace /dev/shm/EV3YK2/t/TestApprovalRestoreFailurecompensation3796896505/001/.git/bench-intent.json: rename: injected ledger#1 failure, want it to hold injected board#2 failure\",1\nskips[0]{package,test,reason}:\n## 5 native owner in place of the injected replacer. subject internal/commitment/repository/replacement.go\nold: \\treturn replacer.Replace(path, data, mode)\nnew: \\treturn durablefile.Native().Replace(path, data, mode)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,11,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=21,bytes=2975,omitted_lines=12,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835433269905653-b9c0c3522f2f4064.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCommitmentReplacement/directory-close,\"replacement_test.go:205: err = <nil>, want a *durablefile.Error at stage directory-close with Published=true that wraps injected policy directory-close failure\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCommitmentReplacement/directory-open,\"replacement_test.go:205: err = <nil>, want a *durablefile.Error at stage directory-open with Published=true that wraps injected policy directory-open failure\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCommitmentReplacement/directory-sync,\"replacement_test.go:205: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected policy directory-sync failure\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCommitmentReplacement/rename,\"replacement_test.go:205: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected policy rename failure\",1\nskips[0]{package,test,reason}:\n## 6 restore after the transaction released the lock. subject internal/commitment/repository/repository.go\nold: \\t}, func() {\\n\\t\\tif restore != nil {\\n\\t\\t\\trestoreErr = restore()\\n\\t\\t}\\n\\t})\nnew: \\t}, func() {\\n\\t\\trestoreErr = errors.ErrUnsupported\\n\\t})\\n\\tif errors.Is(restoreErr, errors.ErrUnsupported) {\\n\\t\\trestoreErr = restore()\\n\\t}\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/repository.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=13,bytes=1440,omitted_lines=4,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835648442554770-380b1d2603646147.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalFailureWithoutBoard/ledger#1/rename,\"replacement_test.go:122: policy#2 ran without the ledger lock: stat /dev/shm/422XEB/t/TestApprovalFailureWithoutBoardledger#1rename2293033941/001/.git/bench-intent.json.lock: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/compensation,\"replacement_test.go:122: policy#2 ran without the ledger lock: stat /dev/shm/422XEB/t/TestApprovalRestoreFailurecompensation3260744936/001/.git/bench-intent.json.lock: no such file or directory\",2\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalUnpublishedLedgerFailure,\"replacement_test.go:122: policy#2 ran without the ledger lock: stat /dev/shm/422XEB/t/TestApprovalUnpublishedLedgerFailure2608250465/001/.git/bench-intent.json.lock: no such file or directory\",2\nskips[0]{package,test,reason}:\n## Package passes after each restore (bench test --package ./internal/commitment/repository: pass, 143 tests).\n"
            }
          }
        },
        {
          "id": "dcm-t07-preservation",
          "performer": "claude:dcm_t07",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6db04c672387f77749e07f16235df9cbe9e9566a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "67d8443bb5fd5574b6496be39fb9f7668632664e",
            "digest": "sha256:4839223b04606c3e03775fb65854350f5811cb647b87429117358862a90c5b9d",
            "excerpt": "# preservation: payload, mode, refusal, and authority necessity reds, bench probe at 67d8443bb5fd5574b6496be39fb9f7668632664e, whole package ./internal/commitment/repository, no --run, each restored=yes. The existing approval tests in ./internal/commitment pass unchanged.\n## 1 D30 refusal: swallow the decision validation. subject internal/commitment/repository/repository.go\nold: if strings.TrimSpace(decision) == \"\" || !sanitize.LineSafe(decision) {\nnew: if false && (strings.TrimSpace(decision) == \"\" || !sanitize.LineSafe(decision)) {\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/repository.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=11,bytes=694,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835449800306709-b1e84cd7ad91c36c.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5129,143\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRefusal/decision,\"replacement_test.go:296: err = <nil>, want the caller's refusal\",1\nskips[0]{package,test,reason}:\n## 2 D30 refusal: swallow the non-regular policy destination. subject internal/commitment/repository/replacement.go\nold: \\tpolicyBefore, err := snapshotFile(policyPath, true)\\n\\tif err != nil {\\n\\t\\treturn nil, err\\n\\t}\nnew: \\tpolicyBefore, _ := snapshotFile(policyPath, true)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=11,bytes=895,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835468597292486-27ae2d57c06e3146.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5131,143\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRefusal/policy-directory,\"replacement_test.go:296: err = replace /dev/shm/D4YRUH/t/TestApprovalRefusalpolicy-directory194445960/001/.bench/commitment.json: rename: rename /dev/shm/D4YRUH/t/TestApprovalRefusalpolicy-directory194445960/001/.bench/.commitment.json.2016… (349 bytes)\",2\nskips[0]{package,test,reason}:\n## 3 staged policy mode drift. subject internal/commitment/repository/replacement.go\nold: policyPath, plan.Proposed, 0o644)\nnew: policyPath, plan.Proposed, 0o640)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=13,bytes=1427,omitted_lines=4,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835483566109756-8c3b67bb10373a1e.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-close,\"replacement_test.go:320: /dev/shm/6BYVUK/t/TestApprovalPublishedLedgerFailuredirectory-close3780255087/001/.bench/commitment.json mode = -rw-r-----, want a regular file at -rw-r--r--\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-open,\"replacement_test.go:320: /dev/shm/6BYVUK/t/TestApprovalPublishedLedgerFailuredirectory-open2685484973/001/.bench/commitment.json mode = -rw-r-----, want a regular file at -rw-r--r--\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-sync,\"replacement_test.go:320: /dev/shm/6BYVUK/t/TestApprovalPublishedLedgerFailuredirectory-sync1347081464/001/.bench/commitment.json mode = -rw-r-----, want a regular file at -rw-r--r--\",1\nskips[0]{package,test,reason}:\n## 4 restore drops the snapshot mode. subject internal/commitment/repository/replacement.go\nold: snapshot.path, snapshot.data, snapshot.mode)\nnew: snapshot.path, snapshot.data, 0o644)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=12,bytes=1035,omitted_lines=3,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835500671802953-e189d9826ae88456.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5037,143\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalBoardFailure,\"replacement_test.go:243: /dev/shm/5VBRDO/t/TestApprovalBoardFailure3793655477/001/ROADMAP.md mode = -rw-r--r--, want a regular file at -rw-r-----\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalInitialPolicyFailure,\"replacement_test.go:226: /dev/shm/5VBRDO/t/TestApprovalInitialPolicyFailure2128476510/001/.bench/commitment.json mode = -rw-r--r--, want a regular file at -rw-------\",1\nskips[0]{package,test,reason}:\n## 5 projected board drops its mode. subject internal/commitment/repository/replacement.go\nold: roadmapPath, projected, roadmapBefore.mode)\nnew: roadmapPath, projected, 0o644)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=13,bytes=1391,omitted_lines=4,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835515477392283-7ddc9fcddd618d71.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-close,\"replacement_test.go:321: /dev/shm/4X75OV/t/TestApprovalPublishedLedgerFailuredirectory-close3378494688/001/ROADMAP.md mode = -rw-r--r--, want a regular file at -rw-r-----\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-open,\"replacement_test.go:321: /dev/shm/4X75OV/t/TestApprovalPublishedLedgerFailuredirectory-open2339263155/001/ROADMAP.md mode = -rw-r--r--, want a regular file at -rw-r-----\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-sync,\"replacement_test.go:321: /dev/shm/4X75OV/t/TestApprovalPublishedLedgerFailuredirectory-sync2328475446/001/ROADMAP.md mode = -rw-r--r--, want a regular file at -rw-r-----\",1\nskips[0]{package,test,reason}:\n## 6 projected board bytes drift. subject internal/commitment/repository/replacement.go\nold: roadmapPath, projected, roadmapBefore.mode)\nnew: roadmapPath, projected[:len(projected)-1], roadmapBefore.mode)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=13,bytes=1690,omitted_lines=4,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835532849019815-dea421b0feae9ee2.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-close,\"replacement_test.go:321: /dev/shm/ZYT7W4/t/TestApprovalPublishedLedgerFailuredirectory-close487574298/001/ROADMAP.md = \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n2. B\\\", want \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n2. B\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-open,\"replacement_test.go:321: /dev/shm/ZYT7W4/t/TestApprovalPublishedLedgerFailuredirectory-open3596106842/001/ROADMAP.md = \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n2. B\\\", want \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n2. B\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-sync,\"replacement_test.go:321: /dev/shm/ZYT7W4/t/TestApprovalPublishedLedgerFailuredirectory-sync1583446190/001/ROADMAP.md = \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n2. B\\\", want \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n2. B\\\\\\\\n\\\"\",1\nskips[0]{package,test,reason}:\n## 7 fallback write of the new policy after the restore (old-bytes checks). subject internal/commitment/repository/replacement.go\nold: \\trestore := func() error { return restoreAll(replacer, policyBefore, roadmapBefore) }\nnew: \\trestore := func() error { err := restoreAll(replacer, policyBefore, roadmapBefore); _ = durablefile.Replace(policyPath, plan.Proposed, 0o644); return err }\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,6,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=16,bytes=2939,omitted_lines=7,cut_lines=3,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835548257103635-ec86cdeb9767e088.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalFailureWithoutBoard/policy#1/directory-sync,\"replacement_test.go:390: /dev/shm/AVKRU6/t/TestApprovalFailureWithoutBoardpolicy#1directory-sync3716141262/001/.bench/commitment.json = \\\"{\\\\\\\\n  \\\\\\\\\\\"version\\\\\\\\\\\": 1,\\\\\\\\n  \\\\\\\\\\\"milestones\\\\\\\\\\\": [\\\\\\\\n    {\\\\\\\\n      \\\\\\\\\\\"id\\\\\\\\\\\": \\\\\\\\\\\"M1\\\\\\\\\\\",\\\\\\\\n      \\\\\\\\\\\"outcomes\\\\\\\\\\\": [\\\\\\\\n \n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalInitialPolicyFailure,\"replacement_test.go:226: /dev/shm/AVKRU6/t/TestApprovalInitialPolicyFailure2226542620/001/.bench/commitment.json mode = -rw-r--r--, want a regular file at -rw-------\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalUnpublishedLedgerFailure,\"replacement_test.go:260: /dev/shm/AVKRU6/t/TestApprovalUnpublishedLedgerFailure4220542775/001/.bench/commitment.json = \\\"{\\\\\\\\n  \\\\\\\\\\\"version\\\\\\\\\\\": 1,\\\\\\\\n  \\\\\\\\\\\"milestones\\\\\\\\\\\": [\\\\\\\\n    {\\\\\\\\n      \\\\\\\\\\\"id\\\\\\\\\\\": \\\\\\\\\\\"M1\\\\\\\\\\\",\\\\\\\\n      \\\\\\\\\\\"outcomes\\\\\\\\\\\": [\\\\\\\\n        {\\\\\\\\n          \\\\\\\\\\\"i… (1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCommitmentReplacement/rename,\"replacement_test.go:207: /dev/shm/AVKRU6/t/TestCommitmentReplacementrename1767104404/001/.bench/commitment.json = \\\"{\\\\\\\\n  \\\\\\\\\\\"version\\\\\\\\\\\": 1,\\\\\\\\n  \\\\\\\\\\\"milestones\\\\\\\\\\\": [\\\\\\\\n    {\\\\\\\\n      \\\\\\\\\\\"id\\\\\\\\\\\": \\\\\\\\\\\"M1\\\\\\\\\\\",\\\\\\\\n      \\\\\\\\\\\"outcomes\\\\\\\\\\\": [\\\\\\\\n        {\\\\\\\\n          \\\\\\\\\\\"id\\\\\\\\\\\": … (1\nskips[0]{package,test,reason}:\n## 8 board-absent branch omitted (absence check). subject internal/commitment/repository/replacement.go\nold: omit: \\tif !roadmapBefore.present {\\n\\t\\treturn restore, nil\\n\\t}\\n\nnew: (omitted)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,omit,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=11,bytes=808,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835565806380386-349f8f1160d9eecd.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5011,143\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalFailureWithoutBoard/ledger#1/directory-sync,\"replacement_test.go:391: /dev/shm/Y54HZI/t/TestApprovalFailureWithoutBoardledger#1directory-sync711308251/001/ROADMAP.md exists after approval: <nil>\",1\nskips[0]{package,test,reason}:\n## 9 restore error does not name its path. subject internal/commitment/repository/replacement.go\nold: \"commitment approval recovery incomplete: restore %s: %w\", snapshot.path, err\nnew: \"commitment approval recovery incomplete: %w\", err\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=13,bytes=1453,omitted_lines=4,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835581083955557-ac37331b16d26208.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/compensation,\"replacement_test.go:358: err = replace /dev/shm/VRAQCK/t/TestApprovalRestoreFailurecompensation2580197929/001/.git/bench-intent.json: rename: injected ledger#1 failure\",2\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back,\"replacement_test.go:358: err = replace /dev/shm/VRAQCK/t/TestApprovalRestoreFailureread-back252311982/001/ROADMAP.md: directory-sync: injected board#1 failure; the rename completed, so new bytes may already be visible\",2\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/stage,\"replacement_test.go:358: err = replace /dev/shm/VRAQCK/t/TestApprovalRestoreFailurestage3550306848/001/.bench/commitment.json: directory-sync: injected policy#1 failure; the rename completed, so new bytes may already be visible\",2\nskips[0]{package,test,reason}:\n## 10 read-back ignores the mode. subject internal/commitment/repository/replacement.go\nold: if current.present != snapshot.present || !bytes.Equal(current.data, snapshot.data) || current.mode != snapshot.mode {\nnew: if current.present != snapshot.present || !bytes.Equal(current.data, snapshot.data) {\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=11,bytes=895,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835598756301401-522575946c2ae739.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5101,143\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back,\"replacement_test.go:358: err = replace /dev/shm/DANX4M/t/TestApprovalRestoreFailureread-back943075638/001/ROADMAP.md: directory-sync: injected board#1 failure; the rename completed, so new bytes may already be visible, want it to name the u… (327 bytes)\",1\nskips[0]{package,test,reason}:\n## 11 a failed stage still approves the receipt (unapproved-receipt checks). subject internal/commitment/repository/repository.go\nold: \\t\\t\\trestore, err = store.stage(replacer, plan)\\n\\t\\t\\tif err != nil {\\n\\t\\t\\t\\treturn ledger, false, err\\n\\t\\t\\t}\nnew: \\t\\t\\trestore, err = store.stage(replacer, plan)\\n\\t\\t\\tif err != nil {\\n\\t\\t\\t\\trestoreErr = err\\n\\t\\t\\t\\treceipt.Decision = decision\\n\\t\\t\\t\\treceipt.Approved = true\\n\\t\\t\\t\\treturn ledger, true, nil\\n\\t\\t\\t}\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/repository.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=13,bytes=1115,omitted_lines=4,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835614841896578-37e9de1e0232086d.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalFailureWithoutBoard/policy#1/directory-sync,\"replacement_test.go:392: receipt approved=true decision=\\\"decision-1\\\", want approved=false\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalInitialPolicyFailure,\"replacement_test.go:227: receipt approved=true decision=\\\"decision-1\\\", want approved=false\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRefusal/policy-directory,\"replacement_test.go:299: replacements = [ledger#1] after a refusal, want none\",1\nskips[0]{package,test,reason}:\n## Package passes after each restore (bench test --package ./internal/commitment/repository: pass, 143 tests).\n"
          },
          "requirement": "preservation",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0,
          "probe": {
            "mutation": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "67d8443bb5fd5574b6496be39fb9f7668632664e",
              "digest": "sha256:4839223b04606c3e03775fb65854350f5811cb647b87429117358862a90c5b9d",
              "excerpt": "# preservation: payload, mode, refusal, and authority necessity reds, bench probe at 67d8443bb5fd5574b6496be39fb9f7668632664e, whole package ./internal/commitment/repository, no --run, each restored=yes. The existing approval tests in ./internal/commitment pass unchanged.\n## 1 D30 refusal: swallow the decision validation. subject internal/commitment/repository/repository.go\nold: if strings.TrimSpace(decision) == \"\" || !sanitize.LineSafe(decision) {\nnew: if false && (strings.TrimSpace(decision) == \"\" || !sanitize.LineSafe(decision)) {\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/repository.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=11,bytes=694,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835449800306709-b1e84cd7ad91c36c.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5129,143\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRefusal/decision,\"replacement_test.go:296: err = <nil>, want the caller's refusal\",1\nskips[0]{package,test,reason}:\n## 2 D30 refusal: swallow the non-regular policy destination. subject internal/commitment/repository/replacement.go\nold: \\tpolicyBefore, err := snapshotFile(policyPath, true)\\n\\tif err != nil {\\n\\t\\treturn nil, err\\n\\t}\nnew: \\tpolicyBefore, _ := snapshotFile(policyPath, true)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=11,bytes=895,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835468597292486-27ae2d57c06e3146.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5131,143\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRefusal/policy-directory,\"replacement_test.go:296: err = replace /dev/shm/D4YRUH/t/TestApprovalRefusalpolicy-directory194445960/001/.bench/commitment.json: rename: rename /dev/shm/D4YRUH/t/TestApprovalRefusalpolicy-directory194445960/001/.bench/.commitment.json.2016… (349 bytes)\",2\nskips[0]{package,test,reason}:\n## 3 staged policy mode drift. subject internal/commitment/repository/replacement.go\nold: policyPath, plan.Proposed, 0o644)\nnew: policyPath, plan.Proposed, 0o640)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=13,bytes=1427,omitted_lines=4,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835483566109756-8c3b67bb10373a1e.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-close,\"replacement_test.go:320: /dev/shm/6BYVUK/t/TestApprovalPublishedLedgerFailuredirectory-close3780255087/001/.bench/commitment.json mode = -rw-r-----, want a regular file at -rw-r--r--\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-open,\"replacement_test.go:320: /dev/shm/6BYVUK/t/TestApprovalPublishedLedgerFailuredirectory-open2685484973/001/.bench/commitment.json mode = -rw-r-----, want a regular file at -rw-r--r--\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-sync,\"replacement_test.go:320: /dev/shm/6BYVUK/t/TestApprovalPublishedLedgerFailuredirectory-sync1347081464/001/.bench/commitment.json mode = -rw-r-----, want a regular file at -rw-r--r--\",1\nskips[0]{package,test,reason}:\n## 4 restore drops the snapshot mode. subject internal/commitment/repository/replacement.go\nold: snapshot.path, snapshot.data, snapshot.mode)\nnew: snapshot.path, snapshot.data, 0o644)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=12,bytes=1035,omitted_lines=3,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835500671802953-e189d9826ae88456.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5037,143\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalBoardFailure,\"replacement_test.go:243: /dev/shm/5VBRDO/t/TestApprovalBoardFailure3793655477/001/ROADMAP.md mode = -rw-r--r--, want a regular file at -rw-r-----\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalInitialPolicyFailure,\"replacement_test.go:226: /dev/shm/5VBRDO/t/TestApprovalInitialPolicyFailure2128476510/001/.bench/commitment.json mode = -rw-r--r--, want a regular file at -rw-------\",1\nskips[0]{package,test,reason}:\n## 5 projected board drops its mode. subject internal/commitment/repository/replacement.go\nold: roadmapPath, projected, roadmapBefore.mode)\nnew: roadmapPath, projected, 0o644)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=13,bytes=1391,omitted_lines=4,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835515477392283-7ddc9fcddd618d71.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-close,\"replacement_test.go:321: /dev/shm/4X75OV/t/TestApprovalPublishedLedgerFailuredirectory-close3378494688/001/ROADMAP.md mode = -rw-r--r--, want a regular file at -rw-r-----\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-open,\"replacement_test.go:321: /dev/shm/4X75OV/t/TestApprovalPublishedLedgerFailuredirectory-open2339263155/001/ROADMAP.md mode = -rw-r--r--, want a regular file at -rw-r-----\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-sync,\"replacement_test.go:321: /dev/shm/4X75OV/t/TestApprovalPublishedLedgerFailuredirectory-sync2328475446/001/ROADMAP.md mode = -rw-r--r--, want a regular file at -rw-r-----\",1\nskips[0]{package,test,reason}:\n## 6 projected board bytes drift. subject internal/commitment/repository/replacement.go\nold: roadmapPath, projected, roadmapBefore.mode)\nnew: roadmapPath, projected[:len(projected)-1], roadmapBefore.mode)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=13,bytes=1690,omitted_lines=4,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835532849019815-dea421b0feae9ee2.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-close,\"replacement_test.go:321: /dev/shm/ZYT7W4/t/TestApprovalPublishedLedgerFailuredirectory-close487574298/001/ROADMAP.md = \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n2. B\\\", want \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n2. B\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-open,\"replacement_test.go:321: /dev/shm/ZYT7W4/t/TestApprovalPublishedLedgerFailuredirectory-open3596106842/001/ROADMAP.md = \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n2. B\\\", want \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n2. B\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-sync,\"replacement_test.go:321: /dev/shm/ZYT7W4/t/TestApprovalPublishedLedgerFailuredirectory-sync1583446190/001/ROADMAP.md = \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n2. B\\\", want \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n2. B\\\\\\\\n\\\"\",1\nskips[0]{package,test,reason}:\n## 7 fallback write of the new policy after the restore (old-bytes checks). subject internal/commitment/repository/replacement.go\nold: \\trestore := func() error { return restoreAll(replacer, policyBefore, roadmapBefore) }\nnew: \\trestore := func() error { err := restoreAll(replacer, policyBefore, roadmapBefore); _ = durablefile.Replace(policyPath, plan.Proposed, 0o644); return err }\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,6,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=16,bytes=2939,omitted_lines=7,cut_lines=3,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835548257103635-ec86cdeb9767e088.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalFailureWithoutBoard/policy#1/directory-sync,\"replacement_test.go:390: /dev/shm/AVKRU6/t/TestApprovalFailureWithoutBoardpolicy#1directory-sync3716141262/001/.bench/commitment.json = \\\"{\\\\\\\\n  \\\\\\\\\\\"version\\\\\\\\\\\": 1,\\\\\\\\n  \\\\\\\\\\\"milestones\\\\\\\\\\\": [\\\\\\\\n    {\\\\\\\\n      \\\\\\\\\\\"id\\\\\\\\\\\": \\\\\\\\\\\"M1\\\\\\\\\\\",\\\\\\\\n      \\\\\\\\\\\"outcomes\\\\\\\\\\\": [\\\\\\\\n \n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalInitialPolicyFailure,\"replacement_test.go:226: /dev/shm/AVKRU6/t/TestApprovalInitialPolicyFailure2226542620/001/.bench/commitment.json mode = -rw-r--r--, want a regular file at -rw-------\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalUnpublishedLedgerFailure,\"replacement_test.go:260: /dev/shm/AVKRU6/t/TestApprovalUnpublishedLedgerFailure4220542775/001/.bench/commitment.json = \\\"{\\\\\\\\n  \\\\\\\\\\\"version\\\\\\\\\\\": 1,\\\\\\\\n  \\\\\\\\\\\"milestones\\\\\\\\\\\": [\\\\\\\\n    {\\\\\\\\n      \\\\\\\\\\\"id\\\\\\\\\\\": \\\\\\\\\\\"M1\\\\\\\\\\\",\\\\\\\\n      \\\\\\\\\\\"outcomes\\\\\\\\\\\": [\\\\\\\\n        {\\\\\\\\n          \\\\\\\\\\\"i… (1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCommitmentReplacement/rename,\"replacement_test.go:207: /dev/shm/AVKRU6/t/TestCommitmentReplacementrename1767104404/001/.bench/commitment.json = \\\"{\\\\\\\\n  \\\\\\\\\\\"version\\\\\\\\\\\": 1,\\\\\\\\n  \\\\\\\\\\\"milestones\\\\\\\\\\\": [\\\\\\\\n    {\\\\\\\\n      \\\\\\\\\\\"id\\\\\\\\\\\": \\\\\\\\\\\"M1\\\\\\\\\\\",\\\\\\\\n      \\\\\\\\\\\"outcomes\\\\\\\\\\\": [\\\\\\\\n        {\\\\\\\\n          \\\\\\\\\\\"id\\\\\\\\\\\": … (1\nskips[0]{package,test,reason}:\n## 8 board-absent branch omitted (absence check). subject internal/commitment/repository/replacement.go\nold: omit: \\tif !roadmapBefore.present {\\n\\t\\treturn restore, nil\\n\\t}\\n\nnew: (omitted)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,omit,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=11,bytes=808,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835565806380386-349f8f1160d9eecd.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5011,143\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalFailureWithoutBoard/ledger#1/directory-sync,\"replacement_test.go:391: /dev/shm/Y54HZI/t/TestApprovalFailureWithoutBoardledger#1directory-sync711308251/001/ROADMAP.md exists after approval: <nil>\",1\nskips[0]{package,test,reason}:\n## 9 restore error does not name its path. subject internal/commitment/repository/replacement.go\nold: \"commitment approval recovery incomplete: restore %s: %w\", snapshot.path, err\nnew: \"commitment approval recovery incomplete: %w\", err\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=13,bytes=1453,omitted_lines=4,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835581083955557-ac37331b16d26208.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/compensation,\"replacement_test.go:358: err = replace /dev/shm/VRAQCK/t/TestApprovalRestoreFailurecompensation2580197929/001/.git/bench-intent.json: rename: injected ledger#1 failure\",2\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back,\"replacement_test.go:358: err = replace /dev/shm/VRAQCK/t/TestApprovalRestoreFailureread-back252311982/001/ROADMAP.md: directory-sync: injected board#1 failure; the rename completed, so new bytes may already be visible\",2\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/stage,\"replacement_test.go:358: err = replace /dev/shm/VRAQCK/t/TestApprovalRestoreFailurestage3550306848/001/.bench/commitment.json: directory-sync: injected policy#1 failure; the rename completed, so new bytes may already be visible\",2\nskips[0]{package,test,reason}:\n## 10 read-back ignores the mode. subject internal/commitment/repository/replacement.go\nold: if current.present != snapshot.present || !bytes.Equal(current.data, snapshot.data) || current.mode != snapshot.mode {\nnew: if current.present != snapshot.present || !bytes.Equal(current.data, snapshot.data) {\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=11,bytes=895,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835598756301401-522575946c2ae739.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5101,143\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back,\"replacement_test.go:358: err = replace /dev/shm/DANX4M/t/TestApprovalRestoreFailureread-back943075638/001/ROADMAP.md: directory-sync: injected board#1 failure; the rename completed, so new bytes may already be visible, want it to name the u… (327 bytes)\",1\nskips[0]{package,test,reason}:\n## 11 a failed stage still approves the receipt (unapproved-receipt checks). subject internal/commitment/repository/repository.go\nold: \\t\\t\\trestore, err = store.stage(replacer, plan)\\n\\t\\t\\tif err != nil {\\n\\t\\t\\t\\treturn ledger, false, err\\n\\t\\t\\t}\nnew: \\t\\t\\trestore, err = store.stage(replacer, plan)\\n\\t\\t\\tif err != nil {\\n\\t\\t\\t\\trestoreErr = err\\n\\t\\t\\t\\treceipt.Decision = decision\\n\\t\\t\\t\\treceipt.Approved = true\\n\\t\\t\\t\\treturn ledger, true, nil\\n\\t\\t\\t}\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,67d8443bb5fd5574b6496be39fb9f7668632664e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/repository.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,143\nspilled{lines=13,bytes=1115,omitted_lines=4,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791835614841896578-37e9de1e0232086d.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalFailureWithoutBoard/policy#1/directory-sync,\"replacement_test.go:392: receipt approved=true decision=\\\"decision-1\\\", want approved=false\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalInitialPolicyFailure,\"replacement_test.go:227: receipt approved=true decision=\\\"decision-1\\\", want approved=false\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRefusal/policy-directory,\"replacement_test.go:299: replacements = [ledger#1] after a refusal, want none\",1\nskips[0]{package,test,reason}:\n## Package passes after each restore (bench test --package ./internal/commitment/repository: pass, 143 tests).\n"
            }
          }
        },
        {
          "id": "dcm-t07-rep1-repository-tests",
          "performer": "claude:dcm_t07_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3d38043c20090c1b33f3011d740349e61b132497",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "91954b3279d3ab379a62f92630c0be55b62e0421",
            "digest": "sha256:6c2c988bd8c00af84ff9564f4b5a9e1522ce55f65073bd02def5f7062dab4469",
            "excerpt": "# repository-tests at 91954b3279d3ab379a62f92630c0be55b62e0421 (tree clean)\n$ bench test --package ./internal/commitment/repository\npackages: github.com/gibbonmi/bench/internal/commitment/repository,pass,5363,152; failures[0]; skips[0]\n"
          },
          "requirement": "repository-tests",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0
        },
        {
          "id": "dcm-t07-rep1-approval-tests",
          "performer": "claude:dcm_t07_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3d38043c20090c1b33f3011d740349e61b132497",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "91954b3279d3ab379a62f92630c0be55b62e0421",
            "digest": "sha256:79989895a6f45c32b4575e9e6996f57cddb0ed3571a184a5bae2381a67f62a59",
            "excerpt": "# approval-tests at 91954b3279d3ab379a62f92630c0be55b62e0421 (tree clean)\n$ bench test --package ./internal/commitment\npackages: github.com/gibbonmi/bench/internal/commitment,pass,4318,175; failures[0]; skips[0]\n"
          },
          "requirement": "approval-tests",
          "command": "bench test --package ./internal/commitment",
          "exit_code": 0
        },
        {
          "id": "dcm-t07-rep1-command-contract",
          "performer": "claude:dcm_t07_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3d38043c20090c1b33f3011d740349e61b132497",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "91954b3279d3ab379a62f92630c0be55b62e0421",
            "digest": "sha256:4c210123fefeaca340e6fa12965578bbc566694646f7dade61b41f6527768d4e",
            "excerpt": "# command-contract at 91954b3279d3ab379a62f92630c0be55b62e0421 (tree clean)\n$ bench test --package ./cmd/bench --run 'CommandRegistry'\npackages: github.com/gibbonmi/bench/cmd/bench,pass,4,5; failures[0]; skips[0]\n"
          },
          "requirement": "command-contract",
          "command": "bench test --package ./cmd/bench --run 'CommandRegistry'",
          "exit_code": 0
        },
        {
          "id": "dcm-t07-rep1-axi-registry",
          "performer": "claude:dcm_t07_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3d38043c20090c1b33f3011d740349e61b132497",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "91954b3279d3ab379a62f92630c0be55b62e0421",
            "digest": "sha256:1fbcf1111fa19ff36aca0f2bda2aa021080b8422333ac50743a530ddbb46e546",
            "excerpt": "# axi-registry at 91954b3279d3ab379a62f92630c0be55b62e0421 (tree clean)\n$ bench test --check axi-query-registry\ncheck: axi-query-registry,conformance,1,0; packages: internal/conformance,pass,10,1; failures[0]; skips[0]\n"
          },
          "requirement": "axi-registry",
          "command": "bench test --check axi-query-registry",
          "exit_code": 0
        },
        {
          "id": "dcm-t07-rep1-routing",
          "performer": "claude:dcm_t07_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3d38043c20090c1b33f3011d740349e61b132497",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "91954b3279d3ab379a62f92630c0be55b62e0421",
            "digest": "sha256:c038bc70e21fd4bfa7c3da46e374e394fae9c7790314c131cc91aa3f6bf6eadb",
            "excerpt": "# routing at 91954b3279d3ab379a62f92630c0be55b62e0421 (tree clean)\n$ bench test --check subcommand-routing\ncheck: subcommand-routing,conformance,1,0; packages: internal/conformance,pass,24,1; failures[0]; skips[0]\n"
          },
          "requirement": "routing",
          "command": "bench test --check subcommand-routing",
          "exit_code": 0
        },
        {
          "id": "dcm-t07-rep1-cmd-bench",
          "performer": "claude:dcm_t07_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3d38043c20090c1b33f3011d740349e61b132497",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "91954b3279d3ab379a62f92630c0be55b62e0421",
            "digest": "sha256:22a3592b424d8e954e511fdb2c551af9dfb5214ee466c73b8d42aac1cda44f43",
            "excerpt": "# cmd-bench at 91954b3279d3ab379a62f92630c0be55b62e0421 (tree clean)\n$ bench test --package ./cmd/bench\npackages: github.com/gibbonmi/bench/cmd/bench,pass,13196,335; failures[0]; skips[0]\n"
          },
          "requirement": "cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "dcm-t07-rep1-conformance",
          "performer": "claude:dcm_t07_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3d38043c20090c1b33f3011d740349e61b132497",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "91954b3279d3ab379a62f92630c0be55b62e0421",
            "digest": "sha256:87d1b030bf409c26eb38291812989de3e7b9ba4ab95989191bea9c3e34f58bef",
            "excerpt": "# conformance at 91954b3279d3ab379a62f92630c0be55b62e0421 (tree clean)\n$ bench test --package ./internal/conformance\npackages: github.com/gibbonmi/bench/internal/conformance,pass,53896,1690; failures[0]; skips[1]: TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device (capability: privilege: cannot create a character device: operation not permitted; environment skip, not owned by this repair)\n"
          },
          "requirement": "conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "dcm-t07-rep1-durablefiletest",
          "performer": "claude:dcm_t07_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3d38043c20090c1b33f3011d740349e61b132497",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "91954b3279d3ab379a62f92630c0be55b62e0421",
            "digest": "sha256:1712b0587db2ef1bbae668fcbca479e23aafac4d120bbfd4633a3926ec583e2e",
            "excerpt": "# durablefiletest at 91954b3279d3ab379a62f92630c0be55b62e0421 (tree clean)\n$ bench test --package ./internal/durablefile/durablefiletest\npackages: github.com/gibbonmi/bench/internal/durablefile/durablefiletest,no-tests,0,0; failures[0]; skips[0]. The package holds the shared Router, which the commitment repository and gate packages exercise (152 and 450 tests pass).\n"
          },
          "requirement": "durablefiletest",
          "command": "bench test --package ./internal/durablefile/durablefiletest",
          "exit_code": 0
        },
        {
          "id": "dcm-t07-rep1-root-conformance",
          "performer": "claude:dcm_t07_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3d38043c20090c1b33f3011d740349e61b132497",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "91954b3279d3ab379a62f92630c0be55b62e0421",
            "digest": "sha256:ac02cf5ef5b9ac100c994c942be04b2b274936cf40bdb43601bc3084e66d95e5",
            "excerpt": "# root-conformance at 91954b3279d3ab379a62f92630c0be55b62e0421 (tree clean)\n$ bench test --check root-conformance\ncheck: root-conformance,conformance,1,0; packages: internal/conformance,pass,17250,1; failures[0]; skips[0]\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "dcm-t07-rep1-omission-proof",
          "performer": "claude:dcm_t07_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3d38043c20090c1b33f3011d740349e61b132497",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "91954b3279d3ab379a62f92630c0be55b62e0421",
            "digest": "sha256:20821690163ca23e31250623447a8a2247987c7db2f278731d8ac20c78752385",
            "excerpt": "# omission-proof: the named restore omissions, owner bypass, lock placement, and the D-B6 C1-C3 restore reds, bench probe at 91954b3279d3ab379a62f92630c0be55b62e0421, whole package ./internal/commitment/repository, no --run, each restored=yes. Before the repair, at 635834ac2df7300ad16700eed971bea68c161e06, the C1, C2, and C3 swaps were each silent (0 of 143 failed).\n## 1 restore built only after the first policy write. subject internal/commitment/repository/replacement.go\nold: \\trestore := func() error { return restoreAll(replacer, policyBefore, roadmapBefore) }\\n\\tif err := replace(replacer, policyPath, plan.Proposed, 0o644); err != nil {\\n\\t\\treturn nil, errors.Join(err, restore())\\n\\t}\nnew: \\tif err := replace(replacer, policyPath, plan.Proposed, 0o644); err != nil {\\n\\t\\treturn nil, err\\n\\t}\\n\\trestore := func() error { return restoreAll(replacer, policyBefore, roadmapBefore) }\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,4,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=14,bytes=1843,omitted_lines=5,cut_lines=1,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837478651761567-14d8b52bb72a9bfe.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalFailureWithoutBoard/policy#1/directory-sync,\"replacement_test.go:373: /dev/shm/RXVX7M/t/TestApprovalFailureWithoutBoardpolicy#1directory-sync2506257037/001/.bench/commitment.json = \\\"{\\\\\\\\n  \\\\\\\\\\\"version\\\\\\\\\\\": 1,\\\\\\\\n  \\\\\\\\\\\"milestones\\\\\\\\\\\": [\\\\\\\\n    {\\\\\\\\n      \\\\\\\\\\\"id\\\\\\\\\\\": \\\\\\\\\\\"M1\\\\\\\\\\\",\\\\\\\\n      \\\\\\\\\\\"outcomes\\\\\\\\\\\": [\\\\\\\\n \n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalFailureWithoutPolicy/policy#1,\"replacement_test.go:395: /dev/shm/RXVX7M/t/TestApprovalFailureWithoutPolicypolicy#14268682645/001/.bench/commitment.json exists after approval: <nil>\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalInitialPolicyFailure,\"replacement_test.go:190: /dev/shm/RXVX7M/t/TestApprovalInitialPolicyFailure288769389/001/.bench/commitment.json mode = -rw-r--r--, want a regular file at -rw-------\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/stage,\"replacement_test.go:341: err = replace /dev/shm/RXVX7M/t/TestApprovalRestoreFailurestage3386697470/001/.bench/commitment.json: directory-sync: injected policy#1 failure; the rename completed, so new bytes may already be visible, want it to … (409 bytes)\",1\nskips[0]{package,test,reason}:\n## 2 compensate a published ledger. subject internal/commitment/repository/repository.go\nold: \\treturn changed, errors.Join(err, restoreErr)\nnew: \\tif err != nil && restore != nil {\\n\\t\\trestoreErr = errors.Join(restoreErr, restore())\\n\\t}\\n\\treturn changed, errors.Join(err, restoreErr)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/repository.go,swap,failed,9,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=19,bytes=3315,omitted_lines=10,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837495553301539-b6c007d328891f95.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-sync,\"replacement_test.go:110: policy#2 ran without the ledger lock: stat /dev/shm/I2GREK/t/TestApprovalPublishedLedgerFailuredirectory-sync879508210/001/.git/bench-intent.json.lock: no such file or directory\",3\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/compensation,\"replacement_test.go:110: policy#3 ran without the ledger lock: stat /dev/shm/I2GREK/t/TestApprovalRestoreFailurecompensation3449979163/001/.git/bench-intent.json.lock: no such file or directory\",2\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/first-of-two,\"replacement_test.go:110: policy#3 ran without the ledger lock: stat /dev/shm/I2GREK/t/TestApprovalRestoreFailurefirst-of-two2580375588/001/.git/bench-intent.json.lock: no such file or directory\",2\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalUnpublishedLedgerFailure,\"replacement_test.go:110: policy#3 ran without the ledger lock: stat /dev/shm/I2GREK/t/TestApprovalUnpublishedLedgerFailure4046528577/001/.git/bench-intent.json.lock: no such file or directory\",3\nskips[0]{package,test,reason}:\n## 3 restore only one file. subject internal/commitment/repository/replacement.go\nold: restoreAll(replacer, policyBefore, roadmapBefore)\nnew: restoreAll(replacer, policyBefore)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,6,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=16,bytes=2572,omitted_lines=7,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837513830754473-c4b3ab1a94930aaa.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/first-of-two,\"replacement_test.go:343: /dev/shm/VAUNAC/t/TestApprovalRestoreFailurefirst-of-two2428827075/001/ROADMAP.md = \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n2. B\\\\\\\\n\\\", want \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. old\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back-bytes,\"replacement_test.go:341: err = replace /dev/shm/VAUNAC/t/TestApprovalRestoreFailureread-back-bytes1043763799/001/ROADMAP.md: directory-sync: injected board#1 failure; the rename completed, so new bytes may already be visible, want it to hol… (377 bytes)\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back-mode,\"replacement_test.go:341: err = replace /dev/shm/VAUNAC/t/TestApprovalRestoreFailureread-back-mode774234325/001/ROADMAP.md: directory-sync: injected board#1 failure; the rename completed, so new bytes may already be visible, want it to hold … (373 bytes)\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalUnpublishedLedgerFailure,\"replacement_test.go:225: replacements = [policy#1 board#1 ledger#1 policy#2], want [policy#1 board#1 ledger#1 policy#2 board#2]\",1\nskips[0]{package,test,reason}:\n## 4a discard restore failure in restoreAll. subject internal/commitment/repository/replacement.go\nold: \\treturn errors.Join(errs...)\nnew: \\treturn nil\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,5,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=15,bytes=2333,omitted_lines=6,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837531886590386-9a4c92f441ebb9d9.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/first-of-two,\"replacement_test.go:341: err = replace /dev/shm/SFAZIB/t/TestApprovalRestoreFailurefirst-of-two4232607666/001/.git/bench-intent.json: rename: injected ledger#1 failure, want it to hold [injected ledger#1 failure injected policy#2 failure] a… (356 bytes)\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back-bytes,\"replacement_test.go:341: err = replace /dev/shm/SFAZIB/t/TestApprovalRestoreFailureread-back-bytes450133398/001/ROADMAP.md: directory-sync: injected board#1 failure; the rename completed, so new bytes may already be visible, want it to hold… (375 bytes)\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back-mode,\"replacement_test.go:341: err = replace /dev/shm/SFAZIB/t/TestApprovalRestoreFailureread-back-mode309819791/001/ROADMAP.md: directory-sync: injected board#1 failure; the rename completed, so new bytes may already be visible, want it to hold … (373 bytes)\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/stage,\"replacement_test.go:341: err = replace /dev/shm/SFAZIB/t/TestApprovalRestoreFailurestage3680833375/001/.bench/commitment.json: directory-sync: injected policy#1 failure; the rename completed, so new bytes may already be visible, want it to … (409 bytes)\",1\nskips[0]{package,test,reason}:\n## 4b discard the compensation restore failure in approve. subject internal/commitment/repository/repository.go\nold: \\treturn changed, errors.Join(err, restoreErr)\nnew: \\t_ = restoreErr\\n\\treturn changed, err\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/repository.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=12,bytes=1256,omitted_lines=3,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837550300659048-66507ec1ebee64e8.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5280,152\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/compensation,\"replacement_test.go:341: err = replace /dev/shm/DWVDZL/t/TestApprovalRestoreFailurecompensation504571960/001/.git/bench-intent.json: rename: injected ledger#1 failure, want it to hold [injected ledger#1 failure injected board#2 failure] and… (341 bytes)\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/first-of-two,\"replacement_test.go:341: err = replace /dev/shm/DWVDZL/t/TestApprovalRestoreFailurefirst-of-two1271811107/001/.git/bench-intent.json: rename: injected ledger#1 failure, want it to hold [injected ledger#1 failure injected policy#2 failure] a… (356 bytes)\",1\nskips[0]{package,test,reason}:\n## 5 native owner in place of the injected replacer. subject internal/commitment/repository/replacement.go\nold: \\treturn replacer.Replace(path, data, mode)\nnew: \\treturn durablefile.Native().Replace(path, data, mode)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,15,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=25,bytes=4552,omitted_lines=16,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837567253567147-025c6af74ee75da9.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCommitmentReplacement/directory-close,\"replacement_test.go:174: err = <nil>, want a *durablefile.Error at stage directory-close with Published=true that wraps injected policy#1 directory-close failure\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCommitmentReplacement/directory-open,\"replacement_test.go:174: err = <nil>, want a *durablefile.Error at stage directory-open with Published=true that wraps injected policy#1 directory-open failure\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCommitmentReplacement/directory-sync,\"replacement_test.go:174: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected policy#1 directory-sync failure\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCommitmentReplacement/rename,\"replacement_test.go:174: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected policy#1 rename failure\",1\nskips[0]{package,test,reason}:\n## 6 restore after the transaction released the lock. subject internal/commitment/repository/repository.go\nold: \\t}, func() {\\n\\t\\tif restore != nil {\\n\\t\\t\\trestoreErr = restore()\\n\\t\\t}\\n\\t})\nnew: \\t}, func() {\\n\\t\\trestoreErr = errors.ErrUnsupported\\n\\t})\\n\\tif errors.Is(restoreErr, errors.ErrUnsupported) {\\n\\t\\trestoreErr = restore()\\n\\t}\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/repository.go,swap,failed,5,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=15,bytes=2034,omitted_lines=6,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837586458270675-77064c7551dbe62c.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalFailureWithoutPolicy/ledger#1,\"replacement_test.go:110: board#2 ran without the ledger lock: stat /dev/shm/TODTIC/t/TestApprovalFailureWithoutPolicyledger#12846771948/001/.git/bench-intent.json.lock: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/compensation,\"replacement_test.go:110: policy#2 ran without the ledger lock: stat /dev/shm/TODTIC/t/TestApprovalRestoreFailurecompensation590204423/001/.git/bench-intent.json.lock: no such file or directory\",2\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/first-of-two,\"replacement_test.go:110: policy#2 ran without the ledger lock: stat /dev/shm/TODTIC/t/TestApprovalRestoreFailurefirst-of-two128535701/001/.git/bench-intent.json.lock: no such file or directory\",2\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalUnpublishedLedgerFailure,\"replacement_test.go:110: policy#2 ran without the ledger lock: stat /dev/shm/TODTIC/t/TestApprovalUnpublishedLedgerFailure155305779/001/.git/bench-intent.json.lock: no such file or directory\",2\nskips[0]{package,test,reason}:\n## C1 restoreAll stops at the first failed restore (TestApprovalRestoreFailure/first-of-two). subject internal/commitment/repository/replacement.go\nold: snapshot.path, err))\nnew: snapshot.path, err)); break\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=11,bytes=900,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837332026756848-93619235a88bed9e.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5368,152\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/first-of-two,\"replacement_test.go:343: /dev/shm/T5YEOQ/t/TestApprovalRestoreFailurefirst-of-two4010155728/001/ROADMAP.md = \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n2. B\\\\\\\\n\\\", want \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. old\\\\\\\\n\\\"\",1\nskips[0]{package,test,reason}:\n## C2 skip the restore after an unpublished board error (TestApprovalBoardFailure/rename). subject internal/commitment/repository/replacement.go\nold: roadmapBefore.mode); err != nil {\\n\\t\\treturn nil, errors.Join(err, restore())\nnew: roadmapBefore.mode); err != nil {\\n\\t\\tif replaceErr := (*durablefile.Error)(nil); errors.As(err, &replaceErr) && !replaceErr.Published {\\n\\t\\t\\treturn nil, err\\n\\t\\t}\\n\\t\\treturn nil, errors.Join(err, restore())\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=11,bytes=798,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837350135379568-bd2332e9b17bb004.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5433,152\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalBoardFailure/rename,\"replacement_test.go:207: /dev/shm/4Y7VDP/t/TestApprovalBoardFailurerename1302350695/001/.bench/commitment.json mode = -rw-r--r--, want a regular file at -rw-------\",1\nskips[0]{package,test,reason}:\n## C3 an absent snapshot restore returns nil (TestApprovalFailureWithoutPolicy). subject internal/commitment/repository/replacement.go\nold: \\tif !snapshot.present {\\n\\t\\tif err := os.Remove\nnew: \\tif !snapshot.present {\\n\\t\\treturn nil\\n\\t}\\n\\tif !snapshot.present {\\n\\t\\tif err := os.Remove\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=12,bytes=1049,omitted_lines=3,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837368676051446-6ad8c8dcc2df208b.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5371,152\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalFailureWithoutPolicy/ledger#1,\"replacement_test.go:395: /dev/shm/ZNMR2F/t/TestApprovalFailureWithoutPolicyledger#11167735530/001/.bench/commitment.json exists after approval: <nil>\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalFailureWithoutPolicy/policy#1,\"replacement_test.go:395: /dev/shm/ZNMR2F/t/TestApprovalFailureWithoutPolicypolicy#13936603499/001/.bench/commitment.json exists after approval: <nil>\",1\nskips[0]{package,test,reason}:\n## Package passes after each restore (bench test --package ./internal/commitment/repository at 91954b3279d3ab379a62f92630c0be55b62e0421: pass, 152 tests).\n"
          },
          "requirement": "omission-proof",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0,
          "probe": {
            "mutation": "Create restore only after policy write, compensate a published ledger, restore only one file, or discard restore failure. Each mutation must fail the corresponding real approval witness.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "91954b3279d3ab379a62f92630c0be55b62e0421",
              "digest": "sha256:20821690163ca23e31250623447a8a2247987c7db2f278731d8ac20c78752385",
              "excerpt": "# omission-proof: the named restore omissions, owner bypass, lock placement, and the D-B6 C1-C3 restore reds, bench probe at 91954b3279d3ab379a62f92630c0be55b62e0421, whole package ./internal/commitment/repository, no --run, each restored=yes. Before the repair, at 635834ac2df7300ad16700eed971bea68c161e06, the C1, C2, and C3 swaps were each silent (0 of 143 failed).\n## 1 restore built only after the first policy write. subject internal/commitment/repository/replacement.go\nold: \\trestore := func() error { return restoreAll(replacer, policyBefore, roadmapBefore) }\\n\\tif err := replace(replacer, policyPath, plan.Proposed, 0o644); err != nil {\\n\\t\\treturn nil, errors.Join(err, restore())\\n\\t}\nnew: \\tif err := replace(replacer, policyPath, plan.Proposed, 0o644); err != nil {\\n\\t\\treturn nil, err\\n\\t}\\n\\trestore := func() error { return restoreAll(replacer, policyBefore, roadmapBefore) }\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,4,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=14,bytes=1843,omitted_lines=5,cut_lines=1,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837478651761567-14d8b52bb72a9bfe.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalFailureWithoutBoard/policy#1/directory-sync,\"replacement_test.go:373: /dev/shm/RXVX7M/t/TestApprovalFailureWithoutBoardpolicy#1directory-sync2506257037/001/.bench/commitment.json = \\\"{\\\\\\\\n  \\\\\\\\\\\"version\\\\\\\\\\\": 1,\\\\\\\\n  \\\\\\\\\\\"milestones\\\\\\\\\\\": [\\\\\\\\n    {\\\\\\\\n      \\\\\\\\\\\"id\\\\\\\\\\\": \\\\\\\\\\\"M1\\\\\\\\\\\",\\\\\\\\n      \\\\\\\\\\\"outcomes\\\\\\\\\\\": [\\\\\\\\n \n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalFailureWithoutPolicy/policy#1,\"replacement_test.go:395: /dev/shm/RXVX7M/t/TestApprovalFailureWithoutPolicypolicy#14268682645/001/.bench/commitment.json exists after approval: <nil>\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalInitialPolicyFailure,\"replacement_test.go:190: /dev/shm/RXVX7M/t/TestApprovalInitialPolicyFailure288769389/001/.bench/commitment.json mode = -rw-r--r--, want a regular file at -rw-------\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/stage,\"replacement_test.go:341: err = replace /dev/shm/RXVX7M/t/TestApprovalRestoreFailurestage3386697470/001/.bench/commitment.json: directory-sync: injected policy#1 failure; the rename completed, so new bytes may already be visible, want it to … (409 bytes)\",1\nskips[0]{package,test,reason}:\n## 2 compensate a published ledger. subject internal/commitment/repository/repository.go\nold: \\treturn changed, errors.Join(err, restoreErr)\nnew: \\tif err != nil && restore != nil {\\n\\t\\trestoreErr = errors.Join(restoreErr, restore())\\n\\t}\\n\\treturn changed, errors.Join(err, restoreErr)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/repository.go,swap,failed,9,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=19,bytes=3315,omitted_lines=10,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837495553301539-b6c007d328891f95.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-sync,\"replacement_test.go:110: policy#2 ran without the ledger lock: stat /dev/shm/I2GREK/t/TestApprovalPublishedLedgerFailuredirectory-sync879508210/001/.git/bench-intent.json.lock: no such file or directory\",3\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/compensation,\"replacement_test.go:110: policy#3 ran without the ledger lock: stat /dev/shm/I2GREK/t/TestApprovalRestoreFailurecompensation3449979163/001/.git/bench-intent.json.lock: no such file or directory\",2\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/first-of-two,\"replacement_test.go:110: policy#3 ran without the ledger lock: stat /dev/shm/I2GREK/t/TestApprovalRestoreFailurefirst-of-two2580375588/001/.git/bench-intent.json.lock: no such file or directory\",2\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalUnpublishedLedgerFailure,\"replacement_test.go:110: policy#3 ran without the ledger lock: stat /dev/shm/I2GREK/t/TestApprovalUnpublishedLedgerFailure4046528577/001/.git/bench-intent.json.lock: no such file or directory\",3\nskips[0]{package,test,reason}:\n## 3 restore only one file. subject internal/commitment/repository/replacement.go\nold: restoreAll(replacer, policyBefore, roadmapBefore)\nnew: restoreAll(replacer, policyBefore)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,6,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=16,bytes=2572,omitted_lines=7,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837513830754473-c4b3ab1a94930aaa.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/first-of-two,\"replacement_test.go:343: /dev/shm/VAUNAC/t/TestApprovalRestoreFailurefirst-of-two2428827075/001/ROADMAP.md = \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n2. B\\\\\\\\n\\\", want \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. old\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back-bytes,\"replacement_test.go:341: err = replace /dev/shm/VAUNAC/t/TestApprovalRestoreFailureread-back-bytes1043763799/001/ROADMAP.md: directory-sync: injected board#1 failure; the rename completed, so new bytes may already be visible, want it to hol… (377 bytes)\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back-mode,\"replacement_test.go:341: err = replace /dev/shm/VAUNAC/t/TestApprovalRestoreFailureread-back-mode774234325/001/ROADMAP.md: directory-sync: injected board#1 failure; the rename completed, so new bytes may already be visible, want it to hold … (373 bytes)\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalUnpublishedLedgerFailure,\"replacement_test.go:225: replacements = [policy#1 board#1 ledger#1 policy#2], want [policy#1 board#1 ledger#1 policy#2 board#2]\",1\nskips[0]{package,test,reason}:\n## 4a discard restore failure in restoreAll. subject internal/commitment/repository/replacement.go\nold: \\treturn errors.Join(errs...)\nnew: \\treturn nil\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,5,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=15,bytes=2333,omitted_lines=6,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837531886590386-9a4c92f441ebb9d9.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/first-of-two,\"replacement_test.go:341: err = replace /dev/shm/SFAZIB/t/TestApprovalRestoreFailurefirst-of-two4232607666/001/.git/bench-intent.json: rename: injected ledger#1 failure, want it to hold [injected ledger#1 failure injected policy#2 failure] a… (356 bytes)\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back-bytes,\"replacement_test.go:341: err = replace /dev/shm/SFAZIB/t/TestApprovalRestoreFailureread-back-bytes450133398/001/ROADMAP.md: directory-sync: injected board#1 failure; the rename completed, so new bytes may already be visible, want it to hold… (375 bytes)\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back-mode,\"replacement_test.go:341: err = replace /dev/shm/SFAZIB/t/TestApprovalRestoreFailureread-back-mode309819791/001/ROADMAP.md: directory-sync: injected board#1 failure; the rename completed, so new bytes may already be visible, want it to hold … (373 bytes)\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/stage,\"replacement_test.go:341: err = replace /dev/shm/SFAZIB/t/TestApprovalRestoreFailurestage3680833375/001/.bench/commitment.json: directory-sync: injected policy#1 failure; the rename completed, so new bytes may already be visible, want it to … (409 bytes)\",1\nskips[0]{package,test,reason}:\n## 4b discard the compensation restore failure in approve. subject internal/commitment/repository/repository.go\nold: \\treturn changed, errors.Join(err, restoreErr)\nnew: \\t_ = restoreErr\\n\\treturn changed, err\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/repository.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=12,bytes=1256,omitted_lines=3,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837550300659048-66507ec1ebee64e8.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5280,152\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/compensation,\"replacement_test.go:341: err = replace /dev/shm/DWVDZL/t/TestApprovalRestoreFailurecompensation504571960/001/.git/bench-intent.json: rename: injected ledger#1 failure, want it to hold [injected ledger#1 failure injected board#2 failure] and… (341 bytes)\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/first-of-two,\"replacement_test.go:341: err = replace /dev/shm/DWVDZL/t/TestApprovalRestoreFailurefirst-of-two1271811107/001/.git/bench-intent.json: rename: injected ledger#1 failure, want it to hold [injected ledger#1 failure injected policy#2 failure] a… (356 bytes)\",1\nskips[0]{package,test,reason}:\n## 5 native owner in place of the injected replacer. subject internal/commitment/repository/replacement.go\nold: \\treturn replacer.Replace(path, data, mode)\nnew: \\treturn durablefile.Native().Replace(path, data, mode)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,15,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=25,bytes=4552,omitted_lines=16,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837567253567147-025c6af74ee75da9.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCommitmentReplacement/directory-close,\"replacement_test.go:174: err = <nil>, want a *durablefile.Error at stage directory-close with Published=true that wraps injected policy#1 directory-close failure\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCommitmentReplacement/directory-open,\"replacement_test.go:174: err = <nil>, want a *durablefile.Error at stage directory-open with Published=true that wraps injected policy#1 directory-open failure\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCommitmentReplacement/directory-sync,\"replacement_test.go:174: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected policy#1 directory-sync failure\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCommitmentReplacement/rename,\"replacement_test.go:174: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected policy#1 rename failure\",1\nskips[0]{package,test,reason}:\n## 6 restore after the transaction released the lock. subject internal/commitment/repository/repository.go\nold: \\t}, func() {\\n\\t\\tif restore != nil {\\n\\t\\t\\trestoreErr = restore()\\n\\t\\t}\\n\\t})\nnew: \\t}, func() {\\n\\t\\trestoreErr = errors.ErrUnsupported\\n\\t})\\n\\tif errors.Is(restoreErr, errors.ErrUnsupported) {\\n\\t\\trestoreErr = restore()\\n\\t}\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/repository.go,swap,failed,5,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=15,bytes=2034,omitted_lines=6,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837586458270675-77064c7551dbe62c.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalFailureWithoutPolicy/ledger#1,\"replacement_test.go:110: board#2 ran without the ledger lock: stat /dev/shm/TODTIC/t/TestApprovalFailureWithoutPolicyledger#12846771948/001/.git/bench-intent.json.lock: no such file or directory\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/compensation,\"replacement_test.go:110: policy#2 ran without the ledger lock: stat /dev/shm/TODTIC/t/TestApprovalRestoreFailurecompensation590204423/001/.git/bench-intent.json.lock: no such file or directory\",2\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/first-of-two,\"replacement_test.go:110: policy#2 ran without the ledger lock: stat /dev/shm/TODTIC/t/TestApprovalRestoreFailurefirst-of-two128535701/001/.git/bench-intent.json.lock: no such file or directory\",2\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalUnpublishedLedgerFailure,\"replacement_test.go:110: policy#2 ran without the ledger lock: stat /dev/shm/TODTIC/t/TestApprovalUnpublishedLedgerFailure155305779/001/.git/bench-intent.json.lock: no such file or directory\",2\nskips[0]{package,test,reason}:\n## C1 restoreAll stops at the first failed restore (TestApprovalRestoreFailure/first-of-two). subject internal/commitment/repository/replacement.go\nold: snapshot.path, err))\nnew: snapshot.path, err)); break\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=11,bytes=900,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837332026756848-93619235a88bed9e.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5368,152\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/first-of-two,\"replacement_test.go:343: /dev/shm/T5YEOQ/t/TestApprovalRestoreFailurefirst-of-two4010155728/001/ROADMAP.md = \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n2. B\\\\\\\\n\\\", want \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. old\\\\\\\\n\\\"\",1\nskips[0]{package,test,reason}:\n## C2 skip the restore after an unpublished board error (TestApprovalBoardFailure/rename). subject internal/commitment/repository/replacement.go\nold: roadmapBefore.mode); err != nil {\\n\\t\\treturn nil, errors.Join(err, restore())\nnew: roadmapBefore.mode); err != nil {\\n\\t\\tif replaceErr := (*durablefile.Error)(nil); errors.As(err, &replaceErr) && !replaceErr.Published {\\n\\t\\t\\treturn nil, err\\n\\t\\t}\\n\\t\\treturn nil, errors.Join(err, restore())\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=11,bytes=798,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837350135379568-bd2332e9b17bb004.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5433,152\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalBoardFailure/rename,\"replacement_test.go:207: /dev/shm/4Y7VDP/t/TestApprovalBoardFailurerename1302350695/001/.bench/commitment.json mode = -rw-r--r--, want a regular file at -rw-------\",1\nskips[0]{package,test,reason}:\n## C3 an absent snapshot restore returns nil (TestApprovalFailureWithoutPolicy). subject internal/commitment/repository/replacement.go\nold: \\tif !snapshot.present {\\n\\t\\tif err := os.Remove\nnew: \\tif !snapshot.present {\\n\\t\\treturn nil\\n\\t}\\n\\tif !snapshot.present {\\n\\t\\tif err := os.Remove\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=12,bytes=1049,omitted_lines=3,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837368676051446-6ad8c8dcc2df208b.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5371,152\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalFailureWithoutPolicy/ledger#1,\"replacement_test.go:395: /dev/shm/ZNMR2F/t/TestApprovalFailureWithoutPolicyledger#11167735530/001/.bench/commitment.json exists after approval: <nil>\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalFailureWithoutPolicy/policy#1,\"replacement_test.go:395: /dev/shm/ZNMR2F/t/TestApprovalFailureWithoutPolicypolicy#13936603499/001/.bench/commitment.json exists after approval: <nil>\",1\nskips[0]{package,test,reason}:\n## Package passes after each restore (bench test --package ./internal/commitment/repository at 91954b3279d3ab379a62f92630c0be55b62e0421: pass, 152 tests).\n"
            }
          }
        },
        {
          "id": "dcm-t07-rep1-preservation",
          "performer": "claude:dcm_t07_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3d38043c20090c1b33f3011d740349e61b132497",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "91954b3279d3ab379a62f92630c0be55b62e0421",
            "digest": "sha256:4abce273a81abe841497769cc82345a50ecf53a5a2393e666b04045aa1f9b670",
            "excerpt": "# preservation: payload, mode, refusal, read-back, and no-file-written necessity reds, bench probe at 91954b3279d3ab379a62f92630c0be55b62e0421, whole package ./internal/commitment/repository, no --run, each restored=yes. Before the repair, at 635834ac2df7300ad16700eed971bea68c161e06, the C4, C5, and C6 swaps were each silent (0 of 143 failed). The existing approval tests in ./internal/commitment pass unchanged (175 tests).\n## 1 D30 refusal: swallow the decision validation. subject internal/commitment/repository/repository.go\nold: if strings.TrimSpace(decision) == \"\" || !sanitize.LineSafe(decision) {\nnew: if false && (strings.TrimSpace(decision) == \"\" || !sanitize.LineSafe(decision)) {\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/repository.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=11,bytes=694,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837605081750829-2a6a67186de60953.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5369,152\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRefusal/decision,\"replacement_test.go:262: err = <nil>, want the caller's refusal\",1\nskips[0]{package,test,reason}:\n## 2 D30 refusal: swallow the non-regular policy destination. subject internal/commitment/repository/replacement.go\nold: \\tpolicyBefore, err := snapshotFile(policyPath, true)\\n\\tif err != nil {\\n\\t\\treturn nil, err\\n\\t}\nnew: \\tpolicyBefore, _ := snapshotFile(policyPath, true)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=11,bytes=895,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837622032355509-affabaa8c77a06f2.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5291,152\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRefusal/policy-directory,\"replacement_test.go:262: err = replace /dev/shm/I4HZAD/t/TestApprovalRefusalpolicy-directory1089754129/001/.bench/commitment.json: rename: rename /dev/shm/I4HZAD/t/TestApprovalRefusalpolicy-directory1089754129/001/.bench/.commitment.json.39… (352 bytes)\",2\nskips[0]{package,test,reason}:\n## 3 staged policy mode drift. subject internal/commitment/repository/replacement.go\nold: policyPath, plan.Proposed, 0o644)\nnew: policyPath, plan.Proposed, 0o640)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=13,bytes=1426,omitted_lines=4,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837640165223234-7807c31482a58834.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-close,\"replacement_test.go:292: /dev/shm/BIGAF5/t/TestApprovalPublishedLedgerFailuredirectory-close4250791544/001/.bench/commitment.json mode = -rw-r-----, want a regular file at -rw-r--r--\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-open,\"replacement_test.go:292: /dev/shm/BIGAF5/t/TestApprovalPublishedLedgerFailuredirectory-open3052374036/001/.bench/commitment.json mode = -rw-r-----, want a regular file at -rw-r--r--\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-sync,\"replacement_test.go:292: /dev/shm/BIGAF5/t/TestApprovalPublishedLedgerFailuredirectory-sync816648865/001/.bench/commitment.json mode = -rw-r-----, want a regular file at -rw-r--r--\",1\nskips[0]{package,test,reason}:\n## 4 restore drops the snapshot mode. subject internal/commitment/repository/replacement.go\nold: snapshot.path, snapshot.data, snapshot.mode)\nnew: snapshot.path, snapshot.data, 0o644)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,8,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=18,bytes=2682,omitted_lines=9,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837656558639954-91e17a8f66abe72f.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/first-of-two,\"replacement_test.go:343: /dev/shm/2J2YJH/t/TestApprovalRestoreFailurefirst-of-two2911027674/001/ROADMAP.md mode = -rw-r--r--, want a regular file at -rw-r-----\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back-bytes,\"replacement_test.go:343: /dev/shm/2J2YJH/t/TestApprovalRestoreFailureread-back-bytes3560346476/001/.bench/commitment.json mode = -rw-r--r--, want a regular file at -rw-------\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back-mode,\"replacement_test.go:343: /dev/shm/2J2YJH/t/TestApprovalRestoreFailureread-back-mode2092301234/001/.bench/commitment.json mode = -rw-r--r--, want a regular file at -rw-------\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/stage,\"replacement_test.go:343: /dev/shm/2J2YJH/t/TestApprovalRestoreFailurestage2421170419/001/ROADMAP.md mode = -rw-r--r--, want a regular file at -rw-r-----\",1\nskips[0]{package,test,reason}:\n## 5 projected board drops its mode. subject internal/commitment/repository/replacement.go\nold: roadmapPath, projected, roadmapBefore.mode)\nnew: roadmapPath, projected, 0o644)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=13,bytes=1390,omitted_lines=4,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837674777245646-2cc9f386b824692d.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-close,\"replacement_test.go:293: /dev/shm/5CAOBI/t/TestApprovalPublishedLedgerFailuredirectory-close3050224955/001/ROADMAP.md mode = -rw-r--r--, want a regular file at -rw-r-----\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-open,\"replacement_test.go:293: /dev/shm/5CAOBI/t/TestApprovalPublishedLedgerFailuredirectory-open3179641104/001/ROADMAP.md mode = -rw-r--r--, want a regular file at -rw-r-----\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-sync,\"replacement_test.go:293: /dev/shm/5CAOBI/t/TestApprovalPublishedLedgerFailuredirectory-sync923476706/001/ROADMAP.md mode = -rw-r--r--, want a regular file at -rw-r-----\",1\nskips[0]{package,test,reason}:\n## 6 projected board bytes drift. subject internal/commitment/repository/replacement.go\nold: roadmapPath, projected, roadmapBefore.mode)\nnew: roadmapPath, projected[:len(projected)-1], roadmapBefore.mode)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=13,bytes=1691,omitted_lines=4,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837691835560232-227032b0ca4eee0e.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-close,\"replacement_test.go:293: /dev/shm/3DLKKM/t/TestApprovalPublishedLedgerFailuredirectory-close3753619624/001/ROADMAP.md = \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n2. B\\\", want \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n2. B\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-open,\"replacement_test.go:293: /dev/shm/3DLKKM/t/TestApprovalPublishedLedgerFailuredirectory-open3410808696/001/ROADMAP.md = \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n2. B\\\", want \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n2. B\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-sync,\"replacement_test.go:293: /dev/shm/3DLKKM/t/TestApprovalPublishedLedgerFailuredirectory-sync2330268297/001/ROADMAP.md = \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n2. B\\\", want \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n2. B\\\\\\\\n\\\"\",1\nskips[0]{package,test,reason}:\n## 7 fallback write of the new policy after the restore (old-bytes checks). subject internal/commitment/repository/replacement.go\nold: \\trestore := func() error { return restoreAll(replacer, policyBefore, roadmapBefore) }\nnew: \\trestore := func() error { err := restoreAll(replacer, policyBefore, roadmapBefore); _ = durablefile.Replace(policyPath, plan.Proposed, 0o644); return err }\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,12,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=22,bytes=4393,omitted_lines=13,cut_lines=2,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837710876627007-91b2a55458e1b71b.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back-bytes,\"replacement_test.go:343: /dev/shm/NHHIIL/t/TestApprovalRestoreFailureread-back-bytes3649829969/001/.bench/commitment.json mode = -rw-r--r--, want a regular file at -rw-------\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back-mode,\"replacement_test.go:343: /dev/shm/NHHIIL/t/TestApprovalRestoreFailureread-back-mode2389222185/001/.bench/commitment.json mode = -rw-r--r--, want a regular file at -rw-------\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalUnpublishedLedgerFailure,\"replacement_test.go:227: /dev/shm/NHHIIL/t/TestApprovalUnpublishedLedgerFailure908049803/001/.bench/commitment.json = \\\"{\\\\\\\\n  \\\\\\\\\\\"version\\\\\\\\\\\": 1,\\\\\\\\n  \\\\\\\\\\\"milestones\\\\\\\\\\\": [\\\\\\\\n    {\\\\\\\\n      \\\\\\\\\\\"id\\\\\\\\\\\": \\\\\\\\\\\"M1\\\\\\\\\\\",\\\\\\\\n      \\\\\\\\\\\"outcomes\\\\\\\\\\\": [\\\\\\\\n        {\\\\\\\\n          \\\\\\\\\\\"id… (1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCommitmentReplacement/rename,\"replacement_test.go:176: /dev/shm/NHHIIL/t/TestCommitmentReplacementrename1611191402/001/.bench/commitment.json = \\\"{\\\\\\\\n  \\\\\\\\\\\"version\\\\\\\\\\\": 1,\\\\\\\\n  \\\\\\\\\\\"milestones\\\\\\\\\\\": [\\\\\\\\n    {\\\\\\\\n      \\\\\\\\\\\"id\\\\\\\\\\\": \\\\\\\\\\\"M1\\\\\\\\\\\",\\\\\\\\n      \\\\\\\\\\\"outcomes\\\\\\\\\\\": [\\\\\\\\n        {\\\\\\\\n          \\\\\\\\\\\"id\\\\\\\\\\\": … (1\nskips[0]{package,test,reason}:\n## 8 board-absent branch omitted (absence check). subject internal/commitment/repository/replacement.go\nold: omit: \\tif !roadmapBefore.present {\\n\\t\\treturn restore, nil\\n\\t}\\n\nnew: (omitted)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,omit,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=11,bytes=809,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837727441628528-300c72f2bcacca7a.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5344,152\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalFailureWithoutBoard/ledger#1/directory-sync,\"replacement_test.go:374: /dev/shm/ZNV5ZO/t/TestApprovalFailureWithoutBoardledger#1directory-sync4097011476/001/ROADMAP.md exists after approval: <nil>\",1\nskips[0]{package,test,reason}:\n## 9 restore error does not name its path. subject internal/commitment/repository/replacement.go\nold: \"commitment approval recovery incomplete: restore %s: %w\", snapshot.path, err\nnew: \"commitment approval recovery incomplete: %w\", err\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,5,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=15,bytes=2066,omitted_lines=6,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837745729675273-bc03f5169082c595.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/first-of-two,\"replacement_test.go:341: err = replace /dev/shm/WEZYRB/t/TestApprovalRestoreFailurefirst-of-two1528604111/001/.git/bench-intent.json: rename: injected ledger#1 failure\",2\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back-bytes,\"replacement_test.go:341: err = replace /dev/shm/WEZYRB/t/TestApprovalRestoreFailureread-back-bytes1677050986/001/ROADMAP.md: directory-sync: injected board#1 failure; the rename completed, so new bytes may already be visible\",2\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back-mode,\"replacement_test.go:341: err = replace /dev/shm/WEZYRB/t/TestApprovalRestoreFailureread-back-mode4234979015/001/ROADMAP.md: directory-sync: injected board#1 failure; the rename completed, so new bytes may already be visible\",2\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/stage,\"replacement_test.go:341: err = replace /dev/shm/WEZYRB/t/TestApprovalRestoreFailurestage1182649700/001/.bench/commitment.json: directory-sync: injected policy#1 failure; the rename completed, so new bytes may already be visible\",2\nskips[0]{package,test,reason}:\n## 10 read-back ignores the mode. subject internal/commitment/repository/replacement.go\nold: if current.present != snapshot.present || !bytes.Equal(current.data, snapshot.data) || current.mode != snapshot.mode {\nnew: if current.present != snapshot.present || !bytes.Equal(current.data, snapshot.data) {\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=11,bytes=900,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837762205315143-6859eab721f51f81.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5365,152\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back-mode,\"replacement_test.go:341: err = replace /dev/shm/5QPTNN/t/TestApprovalRestoreFailureread-back-mode2773697350/001/ROADMAP.md: directory-sync: injected board#1 failure; the rename completed, so new bytes may already be visible, want it to hold… (375 bytes)\",1\nskips[0]{package,test,reason}:\n## 11 a failed stage still approves the receipt (unapproved-receipt and empty-trace checks). subject internal/commitment/repository/repository.go\nold: \\t\\t\\trestore, err = store.stage(replacer, plan)\\n\\t\\t\\tif err != nil {\\n\\t\\t\\t\\treturn ledger, false, err\\n\\t\\t\\t}\nnew: \\t\\t\\trestore, err = store.stage(replacer, plan)\\n\\t\\t\\tif err != nil {\\n\\t\\t\\t\\trestoreErr = err\\n\\t\\t\\t\\treceipt.Decision = decision\\n\\t\\t\\t\\treceipt.Approved = true\\n\\t\\t\\t\\treturn ledger, true, nil\\n\\t\\t\\t}\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/repository.go,swap,failed,5,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=15,bytes=1477,omitted_lines=6,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837781684105802-8a44b32f3272925b.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalInitialPolicyFailure,\"replacement_test.go:191: receipt approved=true decision=\\\"decision-1\\\", want approved=false\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRefusal/board-directory,\"replacement_test.go:265: replacements = [ledger#1] after a refusal, want none\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRefusal/board-without-sequence,\"replacement_test.go:265: replacements = [ledger#1] after a refusal, want none\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRefusal/policy-directory,\"replacement_test.go:265: replacements = [ledger#1] after a refusal, want none\",1\nskips[0]{package,test,reason}:\n## C4 the read-back ignores bytes (TestApprovalRestoreFailure/read-back-bytes). subject internal/commitment/repository/replacement.go\nold: !bytes.Equal(current.data, snapshot.data)\nnew: !bytes.Equal(current.data, current.data)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=11,bytes=901,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837384718577242-11bf81a3bd165b96.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5181,152\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back-bytes,\"replacement_test.go:341: err = replace /dev/shm/TNPD55/t/TestApprovalRestoreFailureread-back-bytes2587620370/001/ROADMAP.md: directory-sync: injected board#1 failure; the rename completed, so new bytes may already be visible, want it to hol… (377 bytes)\",1\nskips[0]{package,test,reason}:\n## C5 drop the board snapshot error (TestApprovalRefusal/board-directory). subject internal/commitment/repository/replacement.go\nold: roadmapBefore, err := snapshotFile(roadmapPath, true)\\n\\tif err != nil {\\n\\t\\treturn nil, err\\n\\t}\nnew: roadmapBefore, _ := snapshotFile(roadmapPath, true)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=11,bytes=702,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837403384874016-b1fa2c5afffddb24.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5484,152\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRefusal/board-directory,\"replacement_test.go:262: err = <nil>, want the caller's refusal\",1\nskips[0]{package,test,reason}:\n## C6 drop the ProjectSequence error (TestApprovalRefusal/board-without-sequence). subject internal/commitment/repository/replacement.go\nold: projected, err = roadmap.ProjectSequence(roadmapBefore.data, commitment.Remaining(policy))\\n\\t\\tif err != nil {\\n\\t\\t\\treturn nil, err\\n\\t\\t}\nnew: projected, _ = roadmap.ProjectSequence(roadmapBefore.data, commitment.Remaining(policy))\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=11,bytes=709,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837420180722593-10b6d722c153c84e.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5166,152\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRefusal/board-without-sequence,\"replacement_test.go:262: err = <nil>, want the caller's refusal\",1\nskips[0]{package,test,reason}:\n## C8 H1 a native board write before the policy-snapshot refusal (no-file-written check, policy-directory row only). subject internal/commitment/repository/replacement.go\nold: policyBefore, err := snapshotFile(policyPath, true)\\n\\tif err != nil {\\n\\t\\treturn nil, err\nnew: policyBefore, err := snapshotFile(policyPath, true)\\n\\tif err != nil {\\n\\t\\t_ = durablefile.Replace(roadmapPath, []byte(\"written\\n\"), 0o644)\\n\\t\\treturn nil, err\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=11,bytes=836,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837440599041077-2c749a2d61415ef3.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5263,152\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRefusal/policy-directory,\"replacement_test.go:268: /dev/shm/XEWDPR/t/TestApprovalRefusalpolicy-directory1685335839/001/ROADMAP.md = \\\"written\\\\\\\\n\\\", want \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. old\\\\\\\\n\\\"\",1\nskips[0]{package,test,reason}:\n## C8 H2 a native policy write before the decision refusal (no-file-written check, decision row only). subject internal/commitment/repository/repository.go\nold: return false, errors.New(\"commitment approval refused: decision must be one nonempty line\")\nnew: _ = durablefile.Replace(store.Root+\"/\"+commitment.PolicyPath, []byte(\"written\\n\"), 0o644)\\n\\t\\treturn false, errors.New(\"commitment approval refused: decision must be one nonempty line\")\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/repository.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=11,bytes=953,omitted_lines=2,cut_lines=1,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837458389031187-d22e692f1846c191.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5314,152\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRefusal/decision,\"replacement_test.go:268: /dev/shm/JX6DO7/t/TestApprovalRefusaldecision646836053/001/.bench/commitment.json = \\\"written\\\\\\\\n\\\", want \\\"{\\\\\\\\n  \\\\\\\\\\\"version\\\\\\\\\\\": 1,\\\\\\\\n  \\\\\\\\\\\"milestones\\\\\\\\\\\": [\\\\\\\\n    {\\\\\\\\n      \\\\\\\\\\\"id\\\\\\\\\\\": \\\\\\\\\\\"M1\\\\\\\\\\\",\\\\\\\\n      \\\\\\\\\\\"outcomes\\\\\\\\\\\": [\\\\\\\\n        {\\\\\\\\n     … (520 byte\nskips[0]{package,test,reason}:\n## Package passes after each restore (bench test --package ./internal/commitment/repository at 91954b3279d3ab379a62f92630c0be55b62e0421: pass, 152 tests).\n"
          },
          "requirement": "preservation",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0,
          "probe": {
            "mutation": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "91954b3279d3ab379a62f92630c0be55b62e0421",
              "digest": "sha256:4abce273a81abe841497769cc82345a50ecf53a5a2393e666b04045aa1f9b670",
              "excerpt": "# preservation: payload, mode, refusal, read-back, and no-file-written necessity reds, bench probe at 91954b3279d3ab379a62f92630c0be55b62e0421, whole package ./internal/commitment/repository, no --run, each restored=yes. Before the repair, at 635834ac2df7300ad16700eed971bea68c161e06, the C4, C5, and C6 swaps were each silent (0 of 143 failed). The existing approval tests in ./internal/commitment pass unchanged (175 tests).\n## 1 D30 refusal: swallow the decision validation. subject internal/commitment/repository/repository.go\nold: if strings.TrimSpace(decision) == \"\" || !sanitize.LineSafe(decision) {\nnew: if false && (strings.TrimSpace(decision) == \"\" || !sanitize.LineSafe(decision)) {\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/repository.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=11,bytes=694,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837605081750829-2a6a67186de60953.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5369,152\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRefusal/decision,\"replacement_test.go:262: err = <nil>, want the caller's refusal\",1\nskips[0]{package,test,reason}:\n## 2 D30 refusal: swallow the non-regular policy destination. subject internal/commitment/repository/replacement.go\nold: \\tpolicyBefore, err := snapshotFile(policyPath, true)\\n\\tif err != nil {\\n\\t\\treturn nil, err\\n\\t}\nnew: \\tpolicyBefore, _ := snapshotFile(policyPath, true)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=11,bytes=895,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837622032355509-affabaa8c77a06f2.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5291,152\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRefusal/policy-directory,\"replacement_test.go:262: err = replace /dev/shm/I4HZAD/t/TestApprovalRefusalpolicy-directory1089754129/001/.bench/commitment.json: rename: rename /dev/shm/I4HZAD/t/TestApprovalRefusalpolicy-directory1089754129/001/.bench/.commitment.json.39… (352 bytes)\",2\nskips[0]{package,test,reason}:\n## 3 staged policy mode drift. subject internal/commitment/repository/replacement.go\nold: policyPath, plan.Proposed, 0o644)\nnew: policyPath, plan.Proposed, 0o640)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=13,bytes=1426,omitted_lines=4,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837640165223234-7807c31482a58834.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-close,\"replacement_test.go:292: /dev/shm/BIGAF5/t/TestApprovalPublishedLedgerFailuredirectory-close4250791544/001/.bench/commitment.json mode = -rw-r-----, want a regular file at -rw-r--r--\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-open,\"replacement_test.go:292: /dev/shm/BIGAF5/t/TestApprovalPublishedLedgerFailuredirectory-open3052374036/001/.bench/commitment.json mode = -rw-r-----, want a regular file at -rw-r--r--\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-sync,\"replacement_test.go:292: /dev/shm/BIGAF5/t/TestApprovalPublishedLedgerFailuredirectory-sync816648865/001/.bench/commitment.json mode = -rw-r-----, want a regular file at -rw-r--r--\",1\nskips[0]{package,test,reason}:\n## 4 restore drops the snapshot mode. subject internal/commitment/repository/replacement.go\nold: snapshot.path, snapshot.data, snapshot.mode)\nnew: snapshot.path, snapshot.data, 0o644)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,8,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=18,bytes=2682,omitted_lines=9,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837656558639954-91e17a8f66abe72f.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/first-of-two,\"replacement_test.go:343: /dev/shm/2J2YJH/t/TestApprovalRestoreFailurefirst-of-two2911027674/001/ROADMAP.md mode = -rw-r--r--, want a regular file at -rw-r-----\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back-bytes,\"replacement_test.go:343: /dev/shm/2J2YJH/t/TestApprovalRestoreFailureread-back-bytes3560346476/001/.bench/commitment.json mode = -rw-r--r--, want a regular file at -rw-------\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back-mode,\"replacement_test.go:343: /dev/shm/2J2YJH/t/TestApprovalRestoreFailureread-back-mode2092301234/001/.bench/commitment.json mode = -rw-r--r--, want a regular file at -rw-------\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/stage,\"replacement_test.go:343: /dev/shm/2J2YJH/t/TestApprovalRestoreFailurestage2421170419/001/ROADMAP.md mode = -rw-r--r--, want a regular file at -rw-r-----\",1\nskips[0]{package,test,reason}:\n## 5 projected board drops its mode. subject internal/commitment/repository/replacement.go\nold: roadmapPath, projected, roadmapBefore.mode)\nnew: roadmapPath, projected, 0o644)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=13,bytes=1390,omitted_lines=4,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837674777245646-2cc9f386b824692d.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-close,\"replacement_test.go:293: /dev/shm/5CAOBI/t/TestApprovalPublishedLedgerFailuredirectory-close3050224955/001/ROADMAP.md mode = -rw-r--r--, want a regular file at -rw-r-----\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-open,\"replacement_test.go:293: /dev/shm/5CAOBI/t/TestApprovalPublishedLedgerFailuredirectory-open3179641104/001/ROADMAP.md mode = -rw-r--r--, want a regular file at -rw-r-----\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-sync,\"replacement_test.go:293: /dev/shm/5CAOBI/t/TestApprovalPublishedLedgerFailuredirectory-sync923476706/001/ROADMAP.md mode = -rw-r--r--, want a regular file at -rw-r-----\",1\nskips[0]{package,test,reason}:\n## 6 projected board bytes drift. subject internal/commitment/repository/replacement.go\nold: roadmapPath, projected, roadmapBefore.mode)\nnew: roadmapPath, projected[:len(projected)-1], roadmapBefore.mode)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=13,bytes=1691,omitted_lines=4,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837691835560232-227032b0ca4eee0e.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-close,\"replacement_test.go:293: /dev/shm/3DLKKM/t/TestApprovalPublishedLedgerFailuredirectory-close3753619624/001/ROADMAP.md = \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n2. B\\\", want \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n2. B\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-open,\"replacement_test.go:293: /dev/shm/3DLKKM/t/TestApprovalPublishedLedgerFailuredirectory-open3410808696/001/ROADMAP.md = \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n2. B\\\", want \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n2. B\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalPublishedLedgerFailure/directory-sync,\"replacement_test.go:293: /dev/shm/3DLKKM/t/TestApprovalPublishedLedgerFailuredirectory-sync2330268297/001/ROADMAP.md = \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n2. B\\\", want \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. A\\\\\\\\n2. B\\\\\\\\n\\\"\",1\nskips[0]{package,test,reason}:\n## 7 fallback write of the new policy after the restore (old-bytes checks). subject internal/commitment/repository/replacement.go\nold: \\trestore := func() error { return restoreAll(replacer, policyBefore, roadmapBefore) }\nnew: \\trestore := func() error { err := restoreAll(replacer, policyBefore, roadmapBefore); _ = durablefile.Replace(policyPath, plan.Proposed, 0o644); return err }\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,12,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=22,bytes=4393,omitted_lines=13,cut_lines=2,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837710876627007-91b2a55458e1b71b.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back-bytes,\"replacement_test.go:343: /dev/shm/NHHIIL/t/TestApprovalRestoreFailureread-back-bytes3649829969/001/.bench/commitment.json mode = -rw-r--r--, want a regular file at -rw-------\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back-mode,\"replacement_test.go:343: /dev/shm/NHHIIL/t/TestApprovalRestoreFailureread-back-mode2389222185/001/.bench/commitment.json mode = -rw-r--r--, want a regular file at -rw-------\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalUnpublishedLedgerFailure,\"replacement_test.go:227: /dev/shm/NHHIIL/t/TestApprovalUnpublishedLedgerFailure908049803/001/.bench/commitment.json = \\\"{\\\\\\\\n  \\\\\\\\\\\"version\\\\\\\\\\\": 1,\\\\\\\\n  \\\\\\\\\\\"milestones\\\\\\\\\\\": [\\\\\\\\n    {\\\\\\\\n      \\\\\\\\\\\"id\\\\\\\\\\\": \\\\\\\\\\\"M1\\\\\\\\\\\",\\\\\\\\n      \\\\\\\\\\\"outcomes\\\\\\\\\\\": [\\\\\\\\n        {\\\\\\\\n          \\\\\\\\\\\"id… (1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestCommitmentReplacement/rename,\"replacement_test.go:176: /dev/shm/NHHIIL/t/TestCommitmentReplacementrename1611191402/001/.bench/commitment.json = \\\"{\\\\\\\\n  \\\\\\\\\\\"version\\\\\\\\\\\": 1,\\\\\\\\n  \\\\\\\\\\\"milestones\\\\\\\\\\\": [\\\\\\\\n    {\\\\\\\\n      \\\\\\\\\\\"id\\\\\\\\\\\": \\\\\\\\\\\"M1\\\\\\\\\\\",\\\\\\\\n      \\\\\\\\\\\"outcomes\\\\\\\\\\\": [\\\\\\\\n        {\\\\\\\\n          \\\\\\\\\\\"id\\\\\\\\\\\": … (1\nskips[0]{package,test,reason}:\n## 8 board-absent branch omitted (absence check). subject internal/commitment/repository/replacement.go\nold: omit: \\tif !roadmapBefore.present {\\n\\t\\treturn restore, nil\\n\\t}\\n\nnew: (omitted)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,omit,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=11,bytes=809,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837727441628528-300c72f2bcacca7a.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5344,152\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalFailureWithoutBoard/ledger#1/directory-sync,\"replacement_test.go:374: /dev/shm/ZNV5ZO/t/TestApprovalFailureWithoutBoardledger#1directory-sync4097011476/001/ROADMAP.md exists after approval: <nil>\",1\nskips[0]{package,test,reason}:\n## 9 restore error does not name its path. subject internal/commitment/repository/replacement.go\nold: \"commitment approval recovery incomplete: restore %s: %w\", snapshot.path, err\nnew: \"commitment approval recovery incomplete: %w\", err\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,5,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=15,bytes=2066,omitted_lines=6,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837745729675273-bc03f5169082c595.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/first-of-two,\"replacement_test.go:341: err = replace /dev/shm/WEZYRB/t/TestApprovalRestoreFailurefirst-of-two1528604111/001/.git/bench-intent.json: rename: injected ledger#1 failure\",2\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back-bytes,\"replacement_test.go:341: err = replace /dev/shm/WEZYRB/t/TestApprovalRestoreFailureread-back-bytes1677050986/001/ROADMAP.md: directory-sync: injected board#1 failure; the rename completed, so new bytes may already be visible\",2\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back-mode,\"replacement_test.go:341: err = replace /dev/shm/WEZYRB/t/TestApprovalRestoreFailureread-back-mode4234979015/001/ROADMAP.md: directory-sync: injected board#1 failure; the rename completed, so new bytes may already be visible\",2\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/stage,\"replacement_test.go:341: err = replace /dev/shm/WEZYRB/t/TestApprovalRestoreFailurestage1182649700/001/.bench/commitment.json: directory-sync: injected policy#1 failure; the rename completed, so new bytes may already be visible\",2\nskips[0]{package,test,reason}:\n## 10 read-back ignores the mode. subject internal/commitment/repository/replacement.go\nold: if current.present != snapshot.present || !bytes.Equal(current.data, snapshot.data) || current.mode != snapshot.mode {\nnew: if current.present != snapshot.present || !bytes.Equal(current.data, snapshot.data) {\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=11,bytes=900,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837762205315143-6859eab721f51f81.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5365,152\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back-mode,\"replacement_test.go:341: err = replace /dev/shm/5QPTNN/t/TestApprovalRestoreFailureread-back-mode2773697350/001/ROADMAP.md: directory-sync: injected board#1 failure; the rename completed, so new bytes may already be visible, want it to hold… (375 bytes)\",1\nskips[0]{package,test,reason}:\n## 11 a failed stage still approves the receipt (unapproved-receipt and empty-trace checks). subject internal/commitment/repository/repository.go\nold: \\t\\t\\trestore, err = store.stage(replacer, plan)\\n\\t\\t\\tif err != nil {\\n\\t\\t\\t\\treturn ledger, false, err\\n\\t\\t\\t}\nnew: \\t\\t\\trestore, err = store.stage(replacer, plan)\\n\\t\\t\\tif err != nil {\\n\\t\\t\\t\\trestoreErr = err\\n\\t\\t\\t\\treceipt.Decision = decision\\n\\t\\t\\t\\treceipt.Approved = true\\n\\t\\t\\t\\treturn ledger, true, nil\\n\\t\\t\\t}\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/repository.go,swap,failed,5,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=15,bytes=1477,omitted_lines=6,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837781684105802-8a44b32f3272925b.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalInitialPolicyFailure,\"replacement_test.go:191: receipt approved=true decision=\\\"decision-1\\\", want approved=false\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRefusal/board-directory,\"replacement_test.go:265: replacements = [ledger#1] after a refusal, want none\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRefusal/board-without-sequence,\"replacement_test.go:265: replacements = [ledger#1] after a refusal, want none\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRefusal/policy-directory,\"replacement_test.go:265: replacements = [ledger#1] after a refusal, want none\",1\nskips[0]{package,test,reason}:\n## C4 the read-back ignores bytes (TestApprovalRestoreFailure/read-back-bytes). subject internal/commitment/repository/replacement.go\nold: !bytes.Equal(current.data, snapshot.data)\nnew: !bytes.Equal(current.data, current.data)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=11,bytes=901,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837384718577242-11bf81a3bd165b96.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5181,152\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRestoreFailure/read-back-bytes,\"replacement_test.go:341: err = replace /dev/shm/TNPD55/t/TestApprovalRestoreFailureread-back-bytes2587620370/001/ROADMAP.md: directory-sync: injected board#1 failure; the rename completed, so new bytes may already be visible, want it to hol… (377 bytes)\",1\nskips[0]{package,test,reason}:\n## C5 drop the board snapshot error (TestApprovalRefusal/board-directory). subject internal/commitment/repository/replacement.go\nold: roadmapBefore, err := snapshotFile(roadmapPath, true)\\n\\tif err != nil {\\n\\t\\treturn nil, err\\n\\t}\nnew: roadmapBefore, _ := snapshotFile(roadmapPath, true)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=11,bytes=702,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837403384874016-b1fa2c5afffddb24.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5484,152\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRefusal/board-directory,\"replacement_test.go:262: err = <nil>, want the caller's refusal\",1\nskips[0]{package,test,reason}:\n## C6 drop the ProjectSequence error (TestApprovalRefusal/board-without-sequence). subject internal/commitment/repository/replacement.go\nold: projected, err = roadmap.ProjectSequence(roadmapBefore.data, commitment.Remaining(policy))\\n\\t\\tif err != nil {\\n\\t\\t\\treturn nil, err\\n\\t\\t}\nnew: projected, _ = roadmap.ProjectSequence(roadmapBefore.data, commitment.Remaining(policy))\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=11,bytes=709,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837420180722593-10b6d722c153c84e.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5166,152\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRefusal/board-without-sequence,\"replacement_test.go:262: err = <nil>, want the caller's refusal\",1\nskips[0]{package,test,reason}:\n## C8 H1 a native board write before the policy-snapshot refusal (no-file-written check, policy-directory row only). subject internal/commitment/repository/replacement.go\nold: policyBefore, err := snapshotFile(policyPath, true)\\n\\tif err != nil {\\n\\t\\treturn nil, err\nnew: policyBefore, err := snapshotFile(policyPath, true)\\n\\tif err != nil {\\n\\t\\t_ = durablefile.Replace(roadmapPath, []byte(\"written\\n\"), 0o644)\\n\\t\\treturn nil, err\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/replacement.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=11,bytes=836,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837440599041077-2c749a2d61415ef3.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5263,152\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRefusal/policy-directory,\"replacement_test.go:268: /dev/shm/XEWDPR/t/TestApprovalRefusalpolicy-directory1685335839/001/ROADMAP.md = \\\"written\\\\\\\\n\\\", want \\\"# Roadmap\\\\\\\\n\\\\\\\\n## Recommended sequence\\\\\\\\n\\\\\\\\n1. old\\\\\\\\n\\\"\",1\nskips[0]{package,test,reason}:\n## C8 H2 a native policy write before the decision refusal (no-file-written check, decision row only). subject internal/commitment/repository/repository.go\nold: return false, errors.New(\"commitment approval refused: decision must be one nonempty line\")\nnew: _ = durablefile.Replace(store.Root+\"/\"+commitment.PolicyPath, []byte(\"written\\n\"), 0o644)\\n\\t\\treturn false, errors.New(\"commitment approval refused: decision must be one nonempty line\")\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/repository.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,all,passed,152\nspilled{lines=11,bytes=953,omitted_lines=2,cut_lines=1,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837458389031187-d22e692f1846c191.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,5314,152\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestApprovalRefusal/decision,\"replacement_test.go:268: /dev/shm/JX6DO7/t/TestApprovalRefusaldecision646836053/001/.bench/commitment.json = \\\"written\\\\\\\\n\\\", want \\\"{\\\\\\\\n  \\\\\\\\\\\"version\\\\\\\\\\\": 1,\\\\\\\\n  \\\\\\\\\\\"milestones\\\\\\\\\\\": [\\\\\\\\n    {\\\\\\\\n      \\\\\\\\\\\"id\\\\\\\\\\\": \\\\\\\\\\\"M1\\\\\\\\\\\",\\\\\\\\n      \\\\\\\\\\\"outcomes\\\\\\\\\\\": [\\\\\\\\n        {\\\\\\\\n     … (520 byte\nskips[0]{package,test,reason}:\n## Package passes after each restore (bench test --package ./internal/commitment/repository at 91954b3279d3ab379a62f92630c0be55b62e0421: pass, 152 tests).\n"
            }
          }
        },
        {
          "id": "dcm-t07-rep1-gate-tests",
          "performer": "claude:dcm_t07_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3d38043c20090c1b33f3011d740349e61b132497",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "91954b3279d3ab379a62f92630c0be55b62e0421",
            "digest": "sha256:fbfbf0b7f189df049077ad3daf57f543a8f9dfe13a62d0070e947d4b13af2d6c",
            "excerpt": "# gate-tests: the D-B1 omission probes and owner-bypass probes rerun through the shared durablefiletest.Router, bench probe at 91954b3279d3ab379a62f92630c0be55b62e0421, whole package ./internal/gate, no --run, each restored=yes.\n## S2 former CreateTemp/Rename sequence in durableReplace (TestVerdictReplacement and the router tests). subject internal/gate/verdict.go\nold: return durableReplaceAt(replacer, gitdir, benchgit.GateCacheFile, rec)\nnew: data, err := json.Marshal(rec)\\n\\tif err != nil {\\n\\t\\treturn err\\n\\t}\\n\\ttemp, err := replacer.CreateTemp(gitdir, \".gate-*\")\\n\\tif err != nil {\\n\\t\\treturn err\\n\\t}\\n\\t_ = temp.Chmod(0o600)\\n\\t_, _ = temp.Write(append(data, 10))\\n\\t_ = temp.Close()\\n\\treturn replacer.Rename(temp.Name(), filepath.Join(gitdir, benchgit.GateCacheFile))\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/verdict.go,swap,failed,8,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,all,passed,450\nspilled{lines=18,bytes=3015,omitted_lines=9,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837852930664313-679f16d0e5b2c42e.out}\n  github.com/gibbonmi/bench/internal/gate,TestGateTerminalDurabilityFailure/green,\"durable_replacement_test.go:117: result = gate.Result{GateExit:0, ActionExit:0, Inspection:gate.Inspection{State:\\\"ready\\\", Status:\\\"green\\\", PendingStatus:\\\"\\\", CachedTree:\\\"64820344ccccd238079e59fe74e5c73051b21d4e\\\", CurrentTree:\\\"64820344ccccd238… (474 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestGateTerminalDurabilityFailure/red,\"durable_replacement_test.go:117: result = gate.Result{GateExit:7, ActionExit:7, Inspection:gate.Inspection{State:\\\"ready\\\", Status:\\\"red\\\", PendingStatus:\\\"\\\", CachedTree:\\\"cc0ee7d63a2767a80f9f85815c27fca47dd2040b\\\", CurrentTree:\\\"cc0ee7d63a2767a80f… (485 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestGateTerminalDurabilityFailure/timeout,\"durable_replacement_test.go:117: result = gate.Result{GateExit:124, ActionExit:124, Inspection:gate.Inspection{State:\\\"ready\\\", Status:\\\"timeout\\\", PendingStatus:\\\"\\\", CachedTree:\\\"b8523df4f9f38376de0a2129b9b9d6287dc3eb20\\\", CurrentTree:\\\"b8523df4f9… (514 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestVerdictReplacement,\"durable_replacement_test.go:29: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected verdict directory-sync failure\",1\nskips[0]{package,test,reason}:\n## O1 leaf cause. subject internal/gate/run_durability.go\nold: return operational(root, gateExit, stderr, category+\": \"+err.Error())\nnew: return operational(root, gateExit, stderr, category)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_durability.go,swap,failed,6,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,all,passed,450\nspilled{lines=16,bytes=2562,omitted_lines=7,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837919336063638-4f74f0566abe6970.out}\n  github.com/gibbonmi/bench/internal/gate,TestGatePendingRecoveryFailure,\"durable_replacement_test.go:169: stderr = \\\"gate final persistence failed\\\\\\\\n\\\", want \\\"gate final persistence failed: replace /dev/shm/QVSBKF/t/TestGatePendingRecoveryFailure328802085/001/.git/bench-last-gate: directory-sync: injected final dir… (478 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestGateTerminalDurabilityFailure/green,\"durable_replacement_test.go:120: stderr = \\\"gate final persistence failed\\\\\\\\n\\\", want \\\"gate final persistence failed: replace /dev/shm/QVSBKF/t/TestGateTerminalDurabilityFailuregreen1381103681/001/.git/bench-last-gate: directory-sync: injected … (355 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestGateTerminalDurabilityFailure/red,\"durable_replacement_test.go:120: stderr = \\\"gate final persistence failed\\\\\\\\n\\\", want \\\"gate final persistence failed: replace /dev/shm/QVSBKF/t/TestGateTerminalDurabilityFailurered3152255634/001/.git/bench-last-gate: directory-sync: injected te… (353 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestGateTerminalDurabilityFailure/timeout,\"durable_replacement_test.go:120: stderr = \\\"gate: timeout\\\\\\\\ngate timeout persistence failed\\\\\\\\n\\\", want \\\"gate timeout persistence failed: replace /dev/shm/QVSBKF/t/TestGateTerminalDurabilityFailuretimeout3274670368/001/.git/bench-last-gate: dire… (376 bytes)\",1\nskips[0]{package,test,reason}:\n## O2 oracle start after pending failure. subject internal/gate/run_transaction.go\nold: return persistenceFailure(storageRoot, 0, stderr, \"gate pending persistence failed\", err)\nnew: fmt.Fprintln(stderr, \"gate pending persistence failed: \"+err.Error())\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_transaction.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,all,passed,450\nspilled{lines=12,bytes=874,omitted_lines=3,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791837988428509589-8865cee3fdeaee12.out}\n  github.com/gibbonmi/bench/internal/gate,fail,26300,450\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGatePendingDurabilityFailure,\"durable_replacement_test.go:81: run counter stat = <nil>, want no oracle start after the pending failure\",1\n  github.com/gibbonmi/bench/internal/gate,TestGateRunRefusesInitialPendingPersistenceWhenCacheIsDirectory,\"run_failure_outcomes_test.go:136: oracle ran after initial pending persistence refusal\",1\nskips[0]{package,test,reason}:\n## O3 hidden evidence failure. subject internal/gate/run_transaction.go\nold: return persistenceFailure(storageRoot, rc, stderr, \"gate evidence persistence failed\", err)\nnew: _ = err\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_transaction.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,all,passed,450\nspilled{lines=12,bytes=1008,omitted_lines=3,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791838052693485493-f1b86111ef2602fa.out}\n  github.com/gibbonmi/bench/internal/gate,fail,26246,450\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateEvidenceDurabilityFailure,\"durable_replacement_test.go:141: result = gate.Result{GateExit:0, ActionExit:0, Inspection:gate.Inspection{State:\\\"ready\\\", Status:\\\"green\\\", PendingStatus:\\\"\\\", CachedTree:\\\"0d2a3785b954205718351e3217cad42bc120fe4b\\\", CurrentTree:\\\"0d2a3785b9542057… (474 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestGateRunRejectsGreenWhenEvidenceDirectoryModeChanges,\"run_failure_outcomes_test.go:183: mode-change exit = 0, stderr=\\\"\\\"\",1\nskips[0]{package,test,reason}:\n## O4 discarded recovery error. subject internal/gate/run_durability.go\nold: return fmt.Errorf(\"%w; pending recovery failed: %w\", err, recovery)\nnew: return fmt.Errorf(\"%w\", err)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_durability.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,all,passed,450\nspilled{lines=13,bytes=1534,omitted_lines=4,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791838114268169801-7c9f4fa25cf56ccc.out}\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateDemotionDurabilityFailure/lock_contention,\"durable_replacement_test.go:221: stderr = \\\"gate execution already in progress; interrupted-green demotion failed: replace /dev/shm/UTSUIY/t/TestGateDemotionDurabilityFailurelock_contention3202101326/001/.git/bench-last-gate: rename: injecte… (719 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestGateDemotionDurabilityFailure/lock_unavailable,\"durable_replacement_test.go:221: stderr = \\\"gate lock unavailable; interrupted-green demotion failed: replace /dev/shm/UTSUIY/t/TestGateDemotionDurabilityFailurelock_unavailable449224453/001/.git/bench-last-gate: rename: injected demotion re… (660 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestGatePendingRecoveryFailure,\"durable_replacement_test.go:169: stderr = \\\"gate final persistence failed: replace /dev/shm/UTSUIY/t/TestGatePendingRecoveryFailure1613531864/001/.git/bench-last-gate: directory-sync: injected final directory-sync failure; the rename complet… (687 bytes)\",1\nskips[0]{package,test,reason}:\n## C1 timeout cause. subject internal/gate/run_transaction.go\nold: return persistenceFailure(storageRoot, 124, stderr, \"gate timeout persistence failed\", err)\nnew: return operational(storageRoot, 124, stderr, \"gate timeout persistence failed\")\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_transaction.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,all,passed,450\nspilled{lines=11,bytes=846,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791838178064631418-6c7987ca7337436d.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,25846,450\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateTerminalDurabilityFailure/timeout,\"durable_replacement_test.go:120: stderr = \\\"gate: timeout\\\\\\\\ngate timeout persistence failed\\\\\\\\n\\\", want \\\"gate timeout persistence failed: replace /dev/shm/YJCDVL/t/TestGateTerminalDurabilityFailuretimeout4268163484/001/.git/bench-last-gate: dire… (376 bytes)\",1\nskips[0]{package,test,reason}:\n## C2 contention demotion. subject internal/gate/run_transaction.go\nold: fmt.Fprintln(stderr, besideDemotion(\"gate execution already in progress\", persistInterruptedIfGreen(replacer, storageRoot, gitdir, plan)))\nnew: fmt.Fprintln(stderr, besideDemotion(\"gate execution already in progress\", func() error { _ = persistInterruptedIfGreen(replacer, storageRoot, gitdir, plan); return nil }()))\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_transaction.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,all,passed,450\nspilled{lines=11,bytes=854,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791838242560488812-a4747e0a2dcfd42d.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,26204,450\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateDemotionDurabilityFailure/lock_contention,\"durable_replacement_test.go:221: stderr = \\\"gate execution already in progress\\\\\\\\ngate owner: pid 3976742 (alive)\\\\\\\\n\\\", want \\\"gate execution already in progress; interrupted-green demotion failed: replace /dev/shm/Q6CQPW/t/TestGateDemotionDurabi… (531 bytes)\",1\nskips[0]{package,test,reason}:\n## C3 D39 cache-replacement count (Router.Count). subject internal/gate/run_transaction.go\nold: return persistenceFailure(storageRoot, rc, stderr, \"gate evidence persistence failed\", err)\nnew: _ = persistRecovering(replacer, gitdir, ready, pending); return persistenceFailure(storageRoot, rc, stderr, \"gate evidence persistence failed\", err)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_transaction.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,all,passed,450\nspilled{lines=12,bytes=1026,omitted_lines=3,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791838308029558170-319dbcc8449d4eb9.out}\n  github.com/gibbonmi/bench/internal/gate,fail,26488,450\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGateEvidenceDurabilityFailure,\"durable_replacement_test.go:147: cache replacements = 2, want only the pending record\",1\n  github.com/gibbonmi/bench/internal/gate,TestGateRunRejectsGreenWhenEvidenceDirectoryModeChanges,\"run_failure_outcomes_test.go:189: inspection after failed retain = gate.Inspection{State:\\\"ready\\\", Status:\\\"green\\\", PendingStatus:\\\"\\\", CachedTree:\\\"ff3d8140275dea459e1057201b4e75496da7bb15\\\", CurrentTree:\\\"ff3d8140275dea459e1057201b4e75496da7bb15… (375 bytes)\",1\nskips[0]{package,test,reason}:\n## C3 D40 visible-green inspection. subject internal/gate/run_durability.go\nold: return fmt.Errorf(\"%w; pending recovery failed: %w\", err, recovery)\nnew: _ = durableReplace(durablefile.Native(), gitdir, pending); return fmt.Errorf(\"%w; pending recovery failed: %w\", err, recovery)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/run_durability.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,all,passed,450\nspilled{lines=11,bytes=834,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791838373724437719-1002d1c80d0fcc00.out}\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,fail,26467,450\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gate,TestGatePendingRecoveryFailure,\"durable_replacement_test.go:172: inspection after the failed recovery = gate.Inspection{State:\\\"pending\\\", Status:\\\"\\\", PendingStatus:\\\"interrupted-pending\\\", CachedTree:\\\"4fc5245cb7869fc41b07686312163a136a780dda\\\", CurrentTree:\\\"4fc5245cb7869fc41b0… (422 bytes)\",1\nskips[0]{package,test,reason}:\n## C3 verdict-writer owner bypass. subject internal/gate/verdict.go\nold: return durableReplaceAt(replacer, gitdir, benchgit.GateCacheFile, rec)\nnew: return durableReplaceAt(durablefile.Native(), gitdir, benchgit.GateCacheFile, rec)\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,91954b3279d3ab379a62f92630c0be55b62e0421,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/verdict.go,swap,failed,9,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,all,passed,450\nspilled{lines=19,bytes=3186,omitted_lines=10,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/00b4e1c8d78f5645bd35ac462e410538/1791838439351201878-c46470850ed5306f.out}\n  github.com/gibbonmi/bench/internal/gate,TestGateTerminalDurabilityFailure/green,\"durable_replacement_test.go:117: result = gate.Result{GateExit:0, ActionExit:0, Inspection:gate.Inspection{State:\\\"ready\\\", Status:\\\"green\\\", PendingStatus:\\\"\\\", CachedTree:\\\"34b1c9453cb41d9c3b160f5dd85d58804695c65d\\\", CurrentTree:\\\"34b1c9453cb41d9c… (474 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestGateTerminalDurabilityFailure/red,\"durable_replacement_test.go:117: result = gate.Result{GateExit:7, ActionExit:7, Inspection:gate.Inspection{State:\\\"ready\\\", Status:\\\"red\\\", PendingStatus:\\\"\\\", CachedTree:\\\"428683da5faae8aac76262ed29e1b27c0cc11662\\\", CurrentTree:\\\"428683da5faae8aac7… (485 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestGateTerminalDurabilityFailure/timeout,\"durable_replacement_test.go:117: result = gate.Result{GateExit:124, ActionExit:124, Inspection:gate.Inspection{State:\\\"ready\\\", Status:\\\"timeout\\\", PendingStatus:\\\"\\\", CachedTree:\\\"5b8676172c2898677d677707d553e15cdf0e6447\\\", CurrentTree:\\\"5b8676172c… (514 bytes)\",1\n  github.com/gibbonmi/bench/internal/gate,TestVerdictReplacement,\"durable_replacement_test.go:29: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected verdict directory-sync failure\",1\nskips[0]{package,test,reason}:\n## restored pass at 91954b3279d3ab379a62f92630c0be55b62e0421: bench test --package ./internal/gate -> pass, 450 tests\n## bench test --package ./internal/gate at 91954b3279d3ab379a62f92630c0be55b62e0421: pass, 450 tests, failures[0], skips[0]\n"
          },
          "requirement": "gate-tests",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "dcm-b6-r1-standards",
          "performer": "claude:dcm_rev_b6_std",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "6db04c672387f77749e07f16235df9cbe9e9566a",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b6_std",
            "digest": "sha256:6a8b3c6592e83b10475bb59d154315425c72ae181aa32d9b039f9b0c30978d4d",
            "excerpt": "## Standards (D-B6, 3cf04078..c9670e29) — 5 findings; worst S1\nS1 | auto-fix | 8 | internal/commitment/repository/replacement_test.go:96 | the test builds the intent lock path a second time; intent.LockPath owns it | use intent.LockPath\nS2 | ask-user (orchestrator: auto-fix by plan expansion) | 5 | replacement_test.go:78-138 vs internal/gate/durable_replacement_test.go:231-269 | two routers send a replacement by destination and ordinal to a Recorder; D-C4 would add a third | one router in durablefiletest, used by both tests\nS3 | ask-user (orchestrator: auto-fix by plan expansion) | 5 | replacement_test.go:24,36-76 | the repo fixture repeats commitmenttest.Repo | expose approve through export_test.go, move the test to repository_test, start from commitmenttest.Repo\nS4 | auto-fix | 6 | replacement_test.go:277-283 | setup closures built outside t.Run call the parent t.Fatal from a subtest | pass the subtest t into setup\nS5 | auto-fix | 4 | replacement_test.go:262-264 | the lock-release check after the call repeats a fact intent owns and tests | delete it\nClean: no algorithm copy; one published classifier; Compensation signature unchanged; restore errors joined and named; repository.go 289 lines.\n"
          },
          "axis": "Standards",
          "base": "3cf04078eb7bf8169f1af53cde9074f0a1da206e",
          "tip": "67d8443bb5fd5574b6496be39fb9f7668632664e",
          "finding_ids": [
            "S1",
            "S2",
            "S3",
            "S4",
            "S5"
          ],
          "supersedes": []
        },
        {
          "id": "dcm-b6-r1-spec",
          "performer": "claude:dcm_rev_b6_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "6db04c672387f77749e07f16235df9cbe9e9566a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b6_spec",
            "digest": "sha256:5109c0d0b008cdaae763c680a5279e3e3d7f010be07afccdf80b0f197d6da0c4",
            "excerpt": "## Spec (D-B6, 3cf04078..c9670e29) — 0 findings\nD23, D46-D51 meet the commitment and restore paragraphs, the seam paragraph, and Source headroom. Restoring after every stage failure matches spec line 130. Plan reorder 661c5ec7 changes only the chunk order and keeps every ordering constraint.\nAdvice: D47 checks policy bytes only; D51 has no policy-absent case; the spec chunk table order differs from the plan order.\n"
          },
          "axis": "Spec",
          "base": "3cf04078eb7bf8169f1af53cde9074f0a1da206e",
          "tip": "67d8443bb5fd5574b6496be39fb9f7668632664e",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "dcm-b6-r1-coverage",
          "performer": "claude:dcm_rev_b6_cov",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "6db04c672387f77749e07f16235df9cbe9e9566a",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b6_cov",
            "digest": "sha256:b4ce360ddc313ab37fba80128b0f2094b9a147980461f366eafecae133db5db1",
            "excerpt": "## Coverage (D-B6, 3cf04078..c9670e29) — 8 findings; worst C1. Nine probes, each restored=yes; author probes 18/18 bite.\nC1 | auto-fix | 9 | replacement.go:68 | restoreAll stopping at the first failed restore is silent | D50 row: ledger#1 rename then policy#2 restore fault; assert the board restored and the policy named\nC2 | auto-fix | 9 | replacement.go:49 | skipping the restore after an unpublished board error is silent | board#1 rename row asserting the old policy bytes and mode\nC3 | auto-fix | 8 | replacement.go:103 | an absent snapshot restore returning nil is silent | no-policy fixture; policy#1 sync and ledger#1 rename failures leave the policy absent\nC4 | auto-fix | 8 | replacement.go:113 | a read-back that ignores bytes is silent | read-back row that rewrites the bytes\nC5 | auto-fix | 8 | replacement.go:26 | dropping the board snapshot error is silent | board-directory refusal row\nC6 | auto-fix | 8 | replacement.go:36 | dropping the ProjectSequence error is silent | refusal row for a board with no sequence; no trace, both files kept\nC7 | ask-user (orchestrator: auto-fix, same as S5) | 6 | replacement_test.go:259-264 | the lock-release check has no safe red and repeats intent's tests | delete it\nC8 | auto-fix | 9 | replacement_test.go:302 | the D30 no-file-written checks had no recorded red; H1 and H2 demonstrated | record H1 and H2\nCovered: probe G, an empty Replacer in the exported Approve, panics in ./internal/commitment.\n"
          },
          "axis": "Coverage",
          "base": "3cf04078eb7bf8169f1af53cde9074f0a1da206e",
          "tip": "67d8443bb5fd5574b6496be39fb9f7668632664e",
          "finding_ids": [
            "C1",
            "C2",
            "C3",
            "C4",
            "C5",
            "C6",
            "C7",
            "C8"
          ],
          "supersedes": []
        },
        {
          "id": "dcm-b6-r2-standards",
          "performer": "claude:dcm_rev_b6_std2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "3d38043c20090c1b33f3011d740349e61b132497",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b6_std2",
            "digest": "sha256:c8db0df55c74f46a61261beff97db2b8f058d4115ff4399157cd5a1221c11e6e",
            "excerpt": "## Standards confirming round (D-B6, repair delta 635834ac..77835db3 and plan e97d88e9) — 0 findings\nS1-S5 closed: intent.LockPath; one durablefiletest.Router used by the approval and gate tests; ApproveWith through export_test.go with commitmenttest.Repo; subtest t in setup; lock-release check removed. The router doc states its rename-keyed limit.\nAdvice: the D43 pre-rename row fits a plain Recorder; Router.Fail lacks a field comment; long lines from the 398-line compression; requireAbsent could move to durablefiletest.\n"
          },
          "axis": "Standards",
          "base": "3cf04078eb7bf8169f1af53cde9074f0a1da206e",
          "tip": "91954b3279d3ab379a62f92630c0be55b62e0421",
          "finding_ids": [],
          "supersedes": [
            "dcm-b6-r1-standards"
          ]
        },
        {
          "id": "dcm-b6-r2-spec",
          "performer": "claude:dcm_rev_b6_spec2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "3d38043c20090c1b33f3011d740349e61b132497",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b6_spec2",
            "digest": "sha256:16b10ecc580a6f7cb7c15c8bde9911c41c81a97014e8a494ab55afca92403e00",
            "excerpt": "## Spec confirming round (D-B6, repair delta 635834ac..77835db3 and plan e97d88e9) — 0 findings\nD23, D46-D51 intact and stronger; the D-B1 gate rows D16, D17, D37-D40 keep their behavior through the shared router; citations name existing tests; plan e97d88e9 adds fence paths and two floor rows only.\n"
          },
          "axis": "Spec",
          "base": "3cf04078eb7bf8169f1af53cde9074f0a1da206e",
          "tip": "91954b3279d3ab379a62f92630c0be55b62e0421",
          "finding_ids": [],
          "supersedes": [
            "dcm-b6-r1-spec"
          ]
        },
        {
          "id": "dcm-b6-r2-coverage",
          "performer": "claude:dcm_rev_b6_cov2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "3d38043c20090c1b33f3011d740349e61b132497",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_b6_cov2",
            "digest": "sha256:f169ef439e617f5bdbad1694a56f6c687b9b34688a8af06a2009fe1a8ac73707",
            "excerpt": "## Coverage confirming round (D-B6, repair delta 635834ac..77835db3) — 0 findings. New-site probes, each restored=yes.\nC1 restoreAll returns at the first error: bit (first-of-two :343). C2 restore dropped on any board failure: bit (4 tests). C3 absent snapshot keeps the policy: bit (:395). C6 ProjectSequence error discarded: bit (:262). C4, C5, C7, C8 closed by recorded evidence. The gate router behaves as the deleted faultPlan, and the recorded gate reds bite through it.\n"
          },
          "axis": "Coverage",
          "base": "3cf04078eb7bf8169f1af53cde9074f0a1da206e",
          "tip": "91954b3279d3ab379a62f92630c0be55b62e0421",
          "finding_ids": [],
          "supersedes": [
            "dcm-b6-r1-coverage"
          ]
        }
      ]
    },
    {
      "id": "D-C2",
      "base": "275c82dfc9fa03d1c87978ebd59c9be4282136bc",
      "tip": "19cb292348cf4aae0f895a139ec39cf8a6ba8edb",
      "plan_digest": "sha256:93a6235e5beda74d1bde28cbbad80df6658b3235cc2436e55544d1e14d979dce",
      "source_digest": "de2e4bde5f117c8f461ebe3ebf1bbbdb21a05cff",
      "acceptance_rows": [
        "D25"
      ],
      "verification": [
        {
          "id": "dcm-t09-caller-tests",
          "performer": "claude:dcm_t09",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "c01e29c180ed8f9fa4cdff27588e8b70c0a14243",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "cf1112f9cb755606a6bf6b58593958e33f2b6b2d",
            "digest": "sha256:12d63e212999cc47dbd2a566b26eccf60d1e7bfe8f45dad814088daed1d44114",
            "excerpt": "bench test --package ./internal/brokermanifest at cf1112f9cb755606a6bf6b58593958e33f2b6b2d (FT354-durable-callers). TestBrokerReplacement (D25 post-rename directory-sync; mode and bytes), TestBrokerPreRenameFailure (D36 rename fault; no manifest), and TestBrokerRefusalSkipsReplacement (D30 digest refusal; non-replacement error, empty trace, no manifest) pass. Callers of Write also pass: bench test --package ./internal/freshness --package ./internal/adopt --package ./internal/adopt/repairtest (79, 244, 28 tests).\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,cf1112f9cb755606a6bf6b58593958e33f2b6b2d,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/brokermanifest,pass,2,3\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "caller-tests",
          "command": "bench test --package ./internal/brokermanifest",
          "exit_code": 0
        },
        {
          "id": "dcm-t09-root-conformance",
          "performer": "claude:dcm_t09",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "c01e29c180ed8f9fa4cdff27588e8b70c0a14243",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "cf1112f9cb755606a6bf6b58593958e33f2b6b2d",
            "digest": "sha256:c68620f2254e06d4d431a20426584e2851bc519a3f89bd80708ae0bbaa0f5428",
            "excerpt": "bench test --check root-conformance at cf1112f9cb755606a6bf6b58593958e33f2b6b2d (FT354-durable-callers)\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16198,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "dcm-t09-omission-proof",
          "performer": "claude:dcm_t09",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "c01e29c180ed8f9fa4cdff27588e8b70c0a14243",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "cf1112f9cb755606a6bf6b58593958e33f2b6b2d",
            "digest": "sha256:08cf22248a6a968bdc4ec12e1df5bad6a758158539299917f62cfa9ffe95a05a",
            "excerpt": "D-C2 omission-proof probes at cf1112f9cb755606a6bf6b58593958e33f2b6b2d (FT354-durable-callers), subject internal/brokermanifest/brokermanifest.go, bench test --package ./internal/brokermanifest, no --run filter.\n\nProbe 1 (bypass the leaf with the former manifest writer).\nSwap old (exact): \"\\tif err := replacer.Replace(path, []byte(content), 0o644); err != nil {\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"publish broker manifest at %s: %w\\\", path, err)\\n\\t}\"\nSwap new (exact, the pre-migration writer from cbccb028e brokermanifest.go:47-68): \"\\ttmp, err := os.CreateTemp(dir, TemporaryPrefix)\\n\\tif err != nil {\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"publish broker manifest at %s: %w\\\", path, err)\\n\\t}\\n\\tname := tmp.Name()\\n\\tif _, err := tmp.WriteString(content); err != nil {\\n\\t\\t_ = tmp.Close()\\n\\t\\t_ = os.Remove(name)\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"publish broker manifest at %s: %w\\\", path, err)\\n\\t}\\n\\tif err := tmp.Close(); err != nil {\\n\\t\\t_ = os.Remove(name)\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"publish broker manifest at %s: %w\\\", path, err)\\n\\t}\\n\\tif err := os.Chmod(name, 0o644); err != nil {\\n\\t\\t_ = os.Remove(name)\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"publish broker manifest at %s: %w\\\", path, err)\\n\\t}\\n\\tif err := os.Rename(name, path); err != nil {\\n\\t\\t_ = os.Remove(name)\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"publish broker manifest at %s: %w\\\", path, err)\\n\\t}\"\nResult: verdict bit, 2 failed tests, restored yes.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/brokermanifest/brokermanifest.go,swap,failed,2,yes\n  github.com/gibbonmi/bench/internal/brokermanifest,fail,2,3\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/brokermanifest,TestBrokerPreRenameFailure,\"replacement_test.go:54: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected broker rename failure\",1\n  github.com/gibbonmi/bench/internal/brokermanifest,TestBrokerReplacement,\"replacement_test.go:40: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected broker directory-sync failure\",1\n\nProbe 2 (native owner in place of the injected replacer).\nSwap old (exact): \"\\tif err := replacer.Replace(path, []byte(content), 0o644); err != nil {\"\nSwap new (exact): \"\\tif err := durablefile.Replace(path, []byte(content), 0o644); err != nil {\"\nResult: verdict bit, 2 failed tests, restored yes.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/brokermanifest/brokermanifest.go,swap,failed,2,yes\n  github.com/gibbonmi/bench/internal/brokermanifest,fail,2,3\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/brokermanifest,TestBrokerPreRenameFailure,\"replacement_test.go:54: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected broker rename failure\",1\n  github.com/gibbonmi/bench/internal/brokermanifest,TestBrokerReplacement,\"replacement_test.go:40: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected broker directory-sync failure\",1\n\nThe package passes 3 tests after each restore.\n"
          },
          "requirement": "omission-proof",
          "command": "bench test --package ./internal/brokermanifest",
          "exit_code": 0,
          "probe": {
            "mutation": "Bypass the leaf with the former manifest writer. The manifest composition test must report the injected missing-failure red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "cf1112f9cb755606a6bf6b58593958e33f2b6b2d",
              "digest": "sha256:08cf22248a6a968bdc4ec12e1df5bad6a758158539299917f62cfa9ffe95a05a",
              "excerpt": "D-C2 omission-proof probes at cf1112f9cb755606a6bf6b58593958e33f2b6b2d (FT354-durable-callers), subject internal/brokermanifest/brokermanifest.go, bench test --package ./internal/brokermanifest, no --run filter.\n\nProbe 1 (bypass the leaf with the former manifest writer).\nSwap old (exact): \"\\tif err := replacer.Replace(path, []byte(content), 0o644); err != nil {\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"publish broker manifest at %s: %w\\\", path, err)\\n\\t}\"\nSwap new (exact, the pre-migration writer from cbccb028e brokermanifest.go:47-68): \"\\ttmp, err := os.CreateTemp(dir, TemporaryPrefix)\\n\\tif err != nil {\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"publish broker manifest at %s: %w\\\", path, err)\\n\\t}\\n\\tname := tmp.Name()\\n\\tif _, err := tmp.WriteString(content); err != nil {\\n\\t\\t_ = tmp.Close()\\n\\t\\t_ = os.Remove(name)\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"publish broker manifest at %s: %w\\\", path, err)\\n\\t}\\n\\tif err := tmp.Close(); err != nil {\\n\\t\\t_ = os.Remove(name)\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"publish broker manifest at %s: %w\\\", path, err)\\n\\t}\\n\\tif err := os.Chmod(name, 0o644); err != nil {\\n\\t\\t_ = os.Remove(name)\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"publish broker manifest at %s: %w\\\", path, err)\\n\\t}\\n\\tif err := os.Rename(name, path); err != nil {\\n\\t\\t_ = os.Remove(name)\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"publish broker manifest at %s: %w\\\", path, err)\\n\\t}\"\nResult: verdict bit, 2 failed tests, restored yes.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/brokermanifest/brokermanifest.go,swap,failed,2,yes\n  github.com/gibbonmi/bench/internal/brokermanifest,fail,2,3\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/brokermanifest,TestBrokerPreRenameFailure,\"replacement_test.go:54: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected broker rename failure\",1\n  github.com/gibbonmi/bench/internal/brokermanifest,TestBrokerReplacement,\"replacement_test.go:40: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected broker directory-sync failure\",1\n\nProbe 2 (native owner in place of the injected replacer).\nSwap old (exact): \"\\tif err := replacer.Replace(path, []byte(content), 0o644); err != nil {\"\nSwap new (exact): \"\\tif err := durablefile.Replace(path, []byte(content), 0o644); err != nil {\"\nResult: verdict bit, 2 failed tests, restored yes.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/brokermanifest/brokermanifest.go,swap,failed,2,yes\n  github.com/gibbonmi/bench/internal/brokermanifest,fail,2,3\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/brokermanifest,TestBrokerPreRenameFailure,\"replacement_test.go:54: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected broker rename failure\",1\n  github.com/gibbonmi/bench/internal/brokermanifest,TestBrokerReplacement,\"replacement_test.go:40: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected broker directory-sync failure\",1\n\nThe package passes 3 tests after each restore.\n"
            }
          }
        },
        {
          "id": "dcm-t09-preservation",
          "performer": "claude:dcm_t09",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "c01e29c180ed8f9fa4cdff27588e8b70c0a14243",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "cf1112f9cb755606a6bf6b58593958e33f2b6b2d",
            "digest": "sha256:59e32b8cafd4882deaf3f0684b9f18c7ca148c791a32a8495e23c9208340f736",
            "excerpt": "D-C2 preservation probes at cf1112f9cb755606a6bf6b58593958e33f2b6b2d (FT354-durable-callers), subject internal/brokermanifest/brokermanifest.go, bench test --package ./internal/brokermanifest, no --run filter. One necessity red for each independent check; each red names exactly one test and one assertion line. Every verdict bit, every restored yes; the package passes 3 tests after each restore. Mode 0644 and the tab-framed payload are frozen from the pre-migration writer; the sha256 literal is the independent digest of the fixture broker bytes.\n\nTemporaryPrefix judgment: the leaf stages under \".\" + base + \".*\", so the exported sweep prefix is \".\" + Name + \".\" to stay tied to the manifest name; it restates the leaf naming rule because internal/durablefile exports no pattern (outside this fence). The freshness residue sweep (freshness_publish_test.go:430) still matches and its package passes.\n\nProbe 1 (mode literal, TestBrokerReplacement RequireFile mode). Swap old (exact): \"replacer.Replace(path, []byte(content), 0o644)\" new (exact): \"replacer.Replace(path, []byte(content), 0o640)\"\n  bit,1 failed,restored yes: TestBrokerReplacement,\"replacement_test.go:41: .../bench-broker.manifest mode = -rw-r-----, want a regular file at -rw-r--r--\"\n\nProbe 2 (payload literal, TestBrokerReplacement RequireFile bytes). Swap old (exact): \"\\\"sha256\\\\t\\\" + digest + \\\"\\\\n\\\"\" new (exact): \"\\\"sha256\\\\t\\\" + digest[:63] + \\\"\\\\n\\\"\"\n  bit,1 failed,restored yes: TestBrokerReplacement,\"replacement_test.go:41: .../bench-broker.manifest = \\\"path\\\\t.../bench\\\\nversion\\\\t1.2.3\\\\nplatform\\\\tlinux-x64\\\\nsha256\\\\tc153e59aa5f0b544d81e95344… (468 bytes)\"\n\nProbe 3 (pre-rename absence, D36: fallback write of the new bytes after an unpublished error). Swap old (exact): \"\\tif err := replacer.Replace(path, []byte(content), 0o644); err != nil {\" new (exact): \"\\tif err := replacer.Replace(path, []byte(content), 0o644); err != nil {\\n\\t\\tif replaceErr := (*durablefile.Error)(nil); errors.As(err, &replaceErr) && !replaceErr.Published {\\n\\t\\t\\t_ = os.WriteFile(path, []byte(content), 0o644)\\n\\t\\t}\"\n  bit,1 failed,restored yes: TestBrokerPreRenameFailure,\"replacement_test.go:55: manifest stat err = <nil>, want it absent\"\n\nProbe 4 (D30 refusal swallowed). Swap old (exact): \"\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"digest promotion broker %s: %w\\\", broker, err)\" new (exact): \"\\t\\t_ = err\"\n  bit,1 failed,restored yes: TestBrokerRefusalSkipsReplacement,\"replacement_test.go:69: err = <nil>, want the caller's own refusal, not a replacement failure\"\n\nProbe 5 (D30 refusal reported as a replacement failure). Same old; new (exact): \"\\t\\treturn \\\"\\\", \\\"\\\", &durablefile.Error{Path: broker, Stage: \\\"digest\\\", Err: err}\"\n  bit,1 failed,restored yes: TestBrokerRefusalSkipsReplacement,\"replacement_test.go:69: err = replace .../001: digest: not a regular file, want the caller's own refusal, not a replacement failure\"\n\nProbe 6 (D30 refusal starts a replacement). Same old; new (exact): \"\\t\\t_ = replacer.Replace(filepath.Join(dir, \\\".probe\\\"), nil, 0o644)\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"digest promotion broker %s: %w\\\", broker, err)\"\n  bit,1 failed,restored yes: TestBrokerRefusalSkipsReplacement,\"replacement_test.go:72: trace = [create mode write file-sync file-close rename directory-open directory-sync directory-close], want no replacement stage\"\n\nProbe 7 (D30 refusal writes a manifest). Same old; new (exact): \"\\t\\t_ = os.WriteFile(filepath.Join(dir, Name), nil, 0o644)\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"digest promotion broker %s: %w\\\", broker, err)\"\n  bit,1 failed,restored yes: TestBrokerRefusalSkipsReplacement,\"replacement_test.go:74: manifest stat err = <nil>, want it absent\"\n"
          },
          "requirement": "preservation",
          "command": "bench test --package ./internal/brokermanifest",
          "exit_code": 0,
          "probe": {
            "mutation": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "cf1112f9cb755606a6bf6b58593958e33f2b6b2d",
              "digest": "sha256:59e32b8cafd4882deaf3f0684b9f18c7ca148c791a32a8495e23c9208340f736",
              "excerpt": "D-C2 preservation probes at cf1112f9cb755606a6bf6b58593958e33f2b6b2d (FT354-durable-callers), subject internal/brokermanifest/brokermanifest.go, bench test --package ./internal/brokermanifest, no --run filter. One necessity red for each independent check; each red names exactly one test and one assertion line. Every verdict bit, every restored yes; the package passes 3 tests after each restore. Mode 0644 and the tab-framed payload are frozen from the pre-migration writer; the sha256 literal is the independent digest of the fixture broker bytes.\n\nTemporaryPrefix judgment: the leaf stages under \".\" + base + \".*\", so the exported sweep prefix is \".\" + Name + \".\" to stay tied to the manifest name; it restates the leaf naming rule because internal/durablefile exports no pattern (outside this fence). The freshness residue sweep (freshness_publish_test.go:430) still matches and its package passes.\n\nProbe 1 (mode literal, TestBrokerReplacement RequireFile mode). Swap old (exact): \"replacer.Replace(path, []byte(content), 0o644)\" new (exact): \"replacer.Replace(path, []byte(content), 0o640)\"\n  bit,1 failed,restored yes: TestBrokerReplacement,\"replacement_test.go:41: .../bench-broker.manifest mode = -rw-r-----, want a regular file at -rw-r--r--\"\n\nProbe 2 (payload literal, TestBrokerReplacement RequireFile bytes). Swap old (exact): \"\\\"sha256\\\\t\\\" + digest + \\\"\\\\n\\\"\" new (exact): \"\\\"sha256\\\\t\\\" + digest[:63] + \\\"\\\\n\\\"\"\n  bit,1 failed,restored yes: TestBrokerReplacement,\"replacement_test.go:41: .../bench-broker.manifest = \\\"path\\\\t.../bench\\\\nversion\\\\t1.2.3\\\\nplatform\\\\tlinux-x64\\\\nsha256\\\\tc153e59aa5f0b544d81e95344… (468 bytes)\"\n\nProbe 3 (pre-rename absence, D36: fallback write of the new bytes after an unpublished error). Swap old (exact): \"\\tif err := replacer.Replace(path, []byte(content), 0o644); err != nil {\" new (exact): \"\\tif err := replacer.Replace(path, []byte(content), 0o644); err != nil {\\n\\t\\tif replaceErr := (*durablefile.Error)(nil); errors.As(err, &replaceErr) && !replaceErr.Published {\\n\\t\\t\\t_ = os.WriteFile(path, []byte(content), 0o644)\\n\\t\\t}\"\n  bit,1 failed,restored yes: TestBrokerPreRenameFailure,\"replacement_test.go:55: manifest stat err = <nil>, want it absent\"\n\nProbe 4 (D30 refusal swallowed). Swap old (exact): \"\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"digest promotion broker %s: %w\\\", broker, err)\" new (exact): \"\\t\\t_ = err\"\n  bit,1 failed,restored yes: TestBrokerRefusalSkipsReplacement,\"replacement_test.go:69: err = <nil>, want the caller's own refusal, not a replacement failure\"\n\nProbe 5 (D30 refusal reported as a replacement failure). Same old; new (exact): \"\\t\\treturn \\\"\\\", \\\"\\\", &durablefile.Error{Path: broker, Stage: \\\"digest\\\", Err: err}\"\n  bit,1 failed,restored yes: TestBrokerRefusalSkipsReplacement,\"replacement_test.go:69: err = replace .../001: digest: not a regular file, want the caller's own refusal, not a replacement failure\"\n\nProbe 6 (D30 refusal starts a replacement). Same old; new (exact): \"\\t\\t_ = replacer.Replace(filepath.Join(dir, \\\".probe\\\"), nil, 0o644)\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"digest promotion broker %s: %w\\\", broker, err)\"\n  bit,1 failed,restored yes: TestBrokerRefusalSkipsReplacement,\"replacement_test.go:72: trace = [create mode write file-sync file-close rename directory-open directory-sync directory-close], want no replacement stage\"\n\nProbe 7 (D30 refusal writes a manifest). Same old; new (exact): \"\\t\\t_ = os.WriteFile(filepath.Join(dir, Name), nil, 0o644)\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"digest promotion broker %s: %w\\\", broker, err)\"\n  bit,1 failed,restored yes: TestBrokerRefusalSkipsReplacement,\"replacement_test.go:74: manifest stat err = <nil>, want it absent\"\n"
            }
          }
        },
        {
          "id": "dcm-t09-r1-caller-tests",
          "performer": "claude:dcm_t09_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "de2e4bde5f117c8f461ebe3ebf1bbbdb21a05cff",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "19cb292348cf4aae0f895a139ec39cf8a6ba8edb",
            "digest": "sha256:ec2d35ef1591c3c576aa1bacecd832548f53bd2f13910e112e667ec022b1bfb4",
            "excerpt": "bench test --package ./internal/brokermanifest at 19cb292348cf4aae0f895a139ec39cf8a6ba8edb (FT354-durable-callers). TestBrokerReplacement (D25 post-rename directory-sync; mode and bytes; staging name pins TemporaryPrefix), TestBrokerPreRenameFailure (D36 rename fault; no manifest via durablefiletest.RequireAbsent), TestBrokerRefusalSkipsReplacement (D30 digest refusal), and TestBrokerWriteReturnsRefusal (public Write returns its refusal; no manifest) pass. Callers of Write also pass: bench test --package ./internal/brokermanifest --package ./internal/durablefile/durablefiletest --package ./internal/freshness --package ./internal/adopt (4, no-tests, 79, 244).\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,19cb292348cf4aae0f895a139ec39cf8a6ba8edb,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/brokermanifest,pass,2,4\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n(dirty: reviews/durable-caller-migration.md record edits only)\n"
          },
          "requirement": "caller-tests",
          "command": "bench test --package ./internal/brokermanifest",
          "exit_code": 0
        },
        {
          "id": "dcm-t09-r1-durablefiletest",
          "performer": "claude:dcm_t09_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "de2e4bde5f117c8f461ebe3ebf1bbbdb21a05cff",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "19cb292348cf4aae0f895a139ec39cf8a6ba8edb",
            "digest": "sha256:5c64e8ffa3e1ea1e6b13c79be6a26617eaf9e9e0f94d9564d58d60abdfc9a71d",
            "excerpt": "bench test --package ./internal/durablefile/durablefiletest at 19cb292348cf4aae0f895a139ec39cf8a6ba8edb (FT354-durable-callers). The package builds with the new RequireAbsent(t testing.TB, path string) helper; its callers in internal/brokermanifest pass (4 tests), and the helper's red is recorded in the preservation excerpt (probes 3 and 7).\ntree[1]{target,head,dirty}:\n  FT354-durable-callers,19cb292348cf4aae0f895a139ec39cf8a6ba8edb,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile/durablefiletest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n(dirty: reviews/durable-caller-migration.md record edits only)\n"
          },
          "requirement": "durablefiletest",
          "command": "bench test --package ./internal/durablefile/durablefiletest",
          "exit_code": 0
        },
        {
          "id": "dcm-t09-r1-root-conformance",
          "performer": "claude:dcm_t09_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "de2e4bde5f117c8f461ebe3ebf1bbbdb21a05cff",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "19cb292348cf4aae0f895a139ec39cf8a6ba8edb",
            "digest": "sha256:072f830015ed034ba8716fddb90f31b9d4806d560d81af5b2bc3f0e8e7e80c5e",
            "excerpt": "bench test --check root-conformance at 19cb292348cf4aae0f895a139ec39cf8a6ba8edb (FT354-durable-callers)\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16144,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "dcm-t09-r1-omission-proof",
          "performer": "claude:dcm_t09_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "de2e4bde5f117c8f461ebe3ebf1bbbdb21a05cff",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "19cb292348cf4aae0f895a139ec39cf8a6ba8edb",
            "digest": "sha256:d15f35de396a7df01ebcc6d28c3f8942c3325eb2a28fafb894f9b7cc6b37e0bf",
            "excerpt": "D-C2 omission-proof probes at 19cb292348cf4aae0f895a139ec39cf8a6ba8edb (FT354-durable-callers), subject internal/brokermanifest/brokermanifest.go, bench test --package ./internal/brokermanifest, no --run filter.\n\nProbe 1 (bypass the leaf with the former manifest writer).\nSwap old (exact): \"\\tif err := replacer.Replace(path, []byte(content), 0o644); err != nil {\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"publish broker manifest at %s: %w\\\", path, err)\\n\\t}\"\nSwap new (exact, the pre-migration writer): \"\\ttmp, err := os.CreateTemp(dir, TemporaryPrefix)\\n\\tif err != nil {\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"publish broker manifest at %s: %w\\\", path, err)\\n\\t}\\n\\tname := tmp.Name()\\n\\tif _, err := tmp.WriteString(content); err != nil {\\n\\t\\t_ = tmp.Close()\\n\\t\\t_ = os.Remove(name)\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"publish broker manifest at %s: %w\\\", path, err)\\n\\t}\\n\\tif err := tmp.Close(); err != nil {\\n\\t\\t_ = os.Remove(name)\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"publish broker manifest at %s: %w\\\", path, err)\\n\\t}\\n\\tif err := os.Chmod(name, 0o644); err != nil {\\n\\t\\t_ = os.Remove(name)\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"publish broker manifest at %s: %w\\\", path, err)\\n\\t}\\n\\tif err := os.Rename(name, path); err != nil {\\n\\t\\t_ = os.Remove(name)\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"publish broker manifest at %s: %w\\\", path, err)\\n\\t}\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/brokermanifest/brokermanifest.go,swap,failed,2,yes\n  github.com/gibbonmi/bench/internal/brokermanifest,fail,3,4\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/brokermanifest,TestBrokerPreRenameFailure,\"replacement_test.go:62: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected broker rename failure\",1\n  github.com/gibbonmi/bench/internal/brokermanifest,TestBrokerReplacement,\"replacement_test.go:40: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected broker directory-sync failure\",1\n\nProbe 2 (native owner in place of the injected replacer).\nSwap old (exact): \"\\tif err := replacer.Replace(path, []byte(content), 0o644); err != nil {\"\nSwap new (exact): \"\\tif err := durablefile.Replace(path, []byte(content), 0o644); err != nil {\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/brokermanifest/brokermanifest.go,swap,failed,2,yes\n  github.com/gibbonmi/bench/internal/brokermanifest,fail,2,4\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/brokermanifest,TestBrokerPreRenameFailure,\"replacement_test.go:62: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected broker rename failure\",1\n  github.com/gibbonmi/bench/internal/brokermanifest,TestBrokerReplacement,\"replacement_test.go:40: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected broker directory-sync failure\",1\n\nThe package passes 4 tests after each restore.\n"
          },
          "requirement": "omission-proof",
          "command": "bench test --package ./internal/brokermanifest",
          "exit_code": 0,
          "probe": {
            "mutation": "Bypass the leaf with the former manifest writer. The manifest composition test must report the injected missing-failure red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "19cb292348cf4aae0f895a139ec39cf8a6ba8edb",
              "digest": "sha256:d15f35de396a7df01ebcc6d28c3f8942c3325eb2a28fafb894f9b7cc6b37e0bf",
              "excerpt": "D-C2 omission-proof probes at 19cb292348cf4aae0f895a139ec39cf8a6ba8edb (FT354-durable-callers), subject internal/brokermanifest/brokermanifest.go, bench test --package ./internal/brokermanifest, no --run filter.\n\nProbe 1 (bypass the leaf with the former manifest writer).\nSwap old (exact): \"\\tif err := replacer.Replace(path, []byte(content), 0o644); err != nil {\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"publish broker manifest at %s: %w\\\", path, err)\\n\\t}\"\nSwap new (exact, the pre-migration writer): \"\\ttmp, err := os.CreateTemp(dir, TemporaryPrefix)\\n\\tif err != nil {\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"publish broker manifest at %s: %w\\\", path, err)\\n\\t}\\n\\tname := tmp.Name()\\n\\tif _, err := tmp.WriteString(content); err != nil {\\n\\t\\t_ = tmp.Close()\\n\\t\\t_ = os.Remove(name)\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"publish broker manifest at %s: %w\\\", path, err)\\n\\t}\\n\\tif err := tmp.Close(); err != nil {\\n\\t\\t_ = os.Remove(name)\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"publish broker manifest at %s: %w\\\", path, err)\\n\\t}\\n\\tif err := os.Chmod(name, 0o644); err != nil {\\n\\t\\t_ = os.Remove(name)\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"publish broker manifest at %s: %w\\\", path, err)\\n\\t}\\n\\tif err := os.Rename(name, path); err != nil {\\n\\t\\t_ = os.Remove(name)\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"publish broker manifest at %s: %w\\\", path, err)\\n\\t}\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/brokermanifest/brokermanifest.go,swap,failed,2,yes\n  github.com/gibbonmi/bench/internal/brokermanifest,fail,3,4\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/brokermanifest,TestBrokerPreRenameFailure,\"replacement_test.go:62: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected broker rename failure\",1\n  github.com/gibbonmi/bench/internal/brokermanifest,TestBrokerReplacement,\"replacement_test.go:40: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected broker directory-sync failure\",1\n\nProbe 2 (native owner in place of the injected replacer).\nSwap old (exact): \"\\tif err := replacer.Replace(path, []byte(content), 0o644); err != nil {\"\nSwap new (exact): \"\\tif err := durablefile.Replace(path, []byte(content), 0o644); err != nil {\"\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/brokermanifest/brokermanifest.go,swap,failed,2,yes\n  github.com/gibbonmi/bench/internal/brokermanifest,fail,2,4\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/brokermanifest,TestBrokerPreRenameFailure,\"replacement_test.go:62: err = <nil>, want a *durablefile.Error at stage rename with Published=false that wraps injected broker rename failure\",1\n  github.com/gibbonmi/bench/internal/brokermanifest,TestBrokerReplacement,\"replacement_test.go:40: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected broker directory-sync failure\",1\n\nThe package passes 4 tests after each restore.\n"
            }
          }
        },
        {
          "id": "dcm-t09-r1-preservation",
          "performer": "claude:dcm_t09_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "de2e4bde5f117c8f461ebe3ebf1bbbdb21a05cff",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "19cb292348cf4aae0f895a139ec39cf8a6ba8edb",
            "digest": "sha256:bee8515c2344b53b5683866a95031abea4da13045c68d862b5e48cc77f481b88",
            "excerpt": "D-C2 preservation probes at 19cb292348cf4aae0f895a139ec39cf8a6ba8edb (FT354-durable-callers), subject internal/brokermanifest/brokermanifest.go, bench test --package ./internal/brokermanifest, no --run filter. Every verdict bit, every restored yes; the package passes 4 tests after each restore. Probes 1 to 7 rerun the dcm_t09 set; probes 8 and 9 are the repair's new reds (S1/C1 and C2).\n\nProbe 1 (mode literal). Swap old (exact): \"replacer.Replace(path, []byte(content), 0o644)\" new (exact): \"replacer.Replace(path, []byte(content), 0o640)\"\n  bit,1 failed,restored yes: TestBrokerReplacement,\"replacement_test.go:41: .../bench-broker.manifest mode = -rw-r-----, want a regular file at -rw-r--r--\"\n\nProbe 2 (payload literal). Swap old (exact): \"\\\"sha256\\\\t\\\" + digest + \\\"\\\\n\\\"\" new (exact): \"\\\"sha256\\\\t\\\" + digest[:63] + \\\"\\\\n\\\"\"\n  bit,1 failed,restored yes: TestBrokerReplacement,\"replacement_test.go:41: .../bench-broker.manifest = \\\"path\\\\t.../bench\\\\nversion\\\\t1.2.3\\\\nplatform\\\\tlinux-x64\\\\nsha256\\\\tc153e59aa5f0b544d81e95344... (468 bytes)\"\n\nProbe 3 (D36 fallback write after an unpublished error). Swap old (exact): \"\\tif err := replacer.Replace(path, []byte(content), 0o644); err != nil {\" new (exact): \"\\tif err := replacer.Replace(path, []byte(content), 0o644); err != nil {\\n\\t\\tif replaceErr := (*durablefile.Error)(nil); errors.As(err, &replaceErr) && !replaceErr.Published {\\n\\t\\t\\t_ = os.WriteFile(path, []byte(content), 0o644)\\n\\t\\t}\"\n  bit,1 failed,restored yes: TestBrokerPreRenameFailure,\"replacement_test.go:63: .../bench-broker.manifest stat err = <nil>, want it absent\" (durablefiletest.RequireAbsent)\n\nProbe 4 (D30 refusal swallowed). Swap old (exact): \"\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"digest promotion broker %s: %w\\\", broker, err)\" new (exact): \"\\t\\t_ = err\"\n  bit,2 failed,restored yes: TestBrokerRefusalSkipsReplacement,\"replacement_test.go:77: err = <nil>, want the caller's own refusal, not a replacement failure\"; TestBrokerWriteReturnsRefusal,\"replacement_test.go:91: Write err = nil, want the refusal of a broker that is not a regular file\"\n\nProbe 5 (D30 refusal reported as a replacement failure). Same old; new (exact): \"\\t\\treturn \\\"\\\", \\\"\\\", &durablefile.Error{Path: broker, Stage: \\\"digest\\\", Err: err}\"\n  bit,1 failed,restored yes: TestBrokerRefusalSkipsReplacement,\"replacement_test.go:77: err = replace .../001: digest: not a regular file, want the caller's own refusal, not a replacement failure\"\n\nProbe 6 (D30 refusal starts a replacement). Same old; new (exact): \"\\t\\t_ = replacer.Replace(filepath.Join(dir, \\\".probe\\\"), nil, 0o644)\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"digest promotion broker %s: %w\\\", broker, err)\"\n  bit,1 failed,restored yes: TestBrokerRefusalSkipsReplacement,\"replacement_test.go:80: trace = [create mode write file-sync file-close rename directory-open directory-sync directory-close], want no replacement stage\"\n\nProbe 7 (D30 refusal writes a manifest). Same old; new (exact): \"\\t\\t_ = os.WriteFile(filepath.Join(dir, Name), nil, 0o644)\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"digest promotion broker %s: %w\\\", broker, err)\"\n  bit,2 failed,restored yes: TestBrokerRefusalSkipsReplacement,\"replacement_test.go:82: .../bench-broker.manifest stat err = <nil>, want it absent\"; TestBrokerWriteReturnsRefusal,\"replacement_test.go:93: .../bench-broker.manifest stat err = <nil>, want it absent\"\n\nProbe 8 (S1/C1 prefix drift). Swap old (exact): \"const TemporaryPrefix = \\\".\\\" + Name + \\\".\\\"\" new (exact): \"const TemporaryPrefix = \\\".bench-broker.\\\"\"\n  bit,1 failed,restored yes: TestBrokerReplacement,\"replacement_test.go:49: staging file \\\".bench-broker.manifest.2834163907\\\", want it named \\\".bench-broker.\\\" plus a random suffix\"\n  Debug note: at 2fd8872c the pin used strings.HasPrefix and this probe was silent (passed, 0 failed, restored yes), because \".bench-broker.\" is still a prefix of the staging name. 19cb2923 compares the staging name up to its last dot.\n\nProbe 9 (C2 public Write swallows every error). Swap old (exact): \"\\treturn write(durablefile.Native(), dir, broker, version)\" new (exact): \"\\tpath, resolved, _ := write(durablefile.Native(), dir, broker, version)\\n\\treturn path, resolved, nil\"\n  bit,1 failed,restored yes: TestBrokerWriteReturnsRefusal,\"replacement_test.go:91: Write err = nil, want the refusal of a broker that is not a regular file\"\n"
          },
          "requirement": "preservation",
          "command": "bench test --package ./internal/brokermanifest",
          "exit_code": 0,
          "probe": {
            "mutation": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "19cb292348cf4aae0f895a139ec39cf8a6ba8edb",
              "digest": "sha256:bee8515c2344b53b5683866a95031abea4da13045c68d862b5e48cc77f481b88",
              "excerpt": "D-C2 preservation probes at 19cb292348cf4aae0f895a139ec39cf8a6ba8edb (FT354-durable-callers), subject internal/brokermanifest/brokermanifest.go, bench test --package ./internal/brokermanifest, no --run filter. Every verdict bit, every restored yes; the package passes 4 tests after each restore. Probes 1 to 7 rerun the dcm_t09 set; probes 8 and 9 are the repair's new reds (S1/C1 and C2).\n\nProbe 1 (mode literal). Swap old (exact): \"replacer.Replace(path, []byte(content), 0o644)\" new (exact): \"replacer.Replace(path, []byte(content), 0o640)\"\n  bit,1 failed,restored yes: TestBrokerReplacement,\"replacement_test.go:41: .../bench-broker.manifest mode = -rw-r-----, want a regular file at -rw-r--r--\"\n\nProbe 2 (payload literal). Swap old (exact): \"\\\"sha256\\\\t\\\" + digest + \\\"\\\\n\\\"\" new (exact): \"\\\"sha256\\\\t\\\" + digest[:63] + \\\"\\\\n\\\"\"\n  bit,1 failed,restored yes: TestBrokerReplacement,\"replacement_test.go:41: .../bench-broker.manifest = \\\"path\\\\t.../bench\\\\nversion\\\\t1.2.3\\\\nplatform\\\\tlinux-x64\\\\nsha256\\\\tc153e59aa5f0b544d81e95344... (468 bytes)\"\n\nProbe 3 (D36 fallback write after an unpublished error). Swap old (exact): \"\\tif err := replacer.Replace(path, []byte(content), 0o644); err != nil {\" new (exact): \"\\tif err := replacer.Replace(path, []byte(content), 0o644); err != nil {\\n\\t\\tif replaceErr := (*durablefile.Error)(nil); errors.As(err, &replaceErr) && !replaceErr.Published {\\n\\t\\t\\t_ = os.WriteFile(path, []byte(content), 0o644)\\n\\t\\t}\"\n  bit,1 failed,restored yes: TestBrokerPreRenameFailure,\"replacement_test.go:63: .../bench-broker.manifest stat err = <nil>, want it absent\" (durablefiletest.RequireAbsent)\n\nProbe 4 (D30 refusal swallowed). Swap old (exact): \"\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"digest promotion broker %s: %w\\\", broker, err)\" new (exact): \"\\t\\t_ = err\"\n  bit,2 failed,restored yes: TestBrokerRefusalSkipsReplacement,\"replacement_test.go:77: err = <nil>, want the caller's own refusal, not a replacement failure\"; TestBrokerWriteReturnsRefusal,\"replacement_test.go:91: Write err = nil, want the refusal of a broker that is not a regular file\"\n\nProbe 5 (D30 refusal reported as a replacement failure). Same old; new (exact): \"\\t\\treturn \\\"\\\", \\\"\\\", &durablefile.Error{Path: broker, Stage: \\\"digest\\\", Err: err}\"\n  bit,1 failed,restored yes: TestBrokerRefusalSkipsReplacement,\"replacement_test.go:77: err = replace .../001: digest: not a regular file, want the caller's own refusal, not a replacement failure\"\n\nProbe 6 (D30 refusal starts a replacement). Same old; new (exact): \"\\t\\t_ = replacer.Replace(filepath.Join(dir, \\\".probe\\\"), nil, 0o644)\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"digest promotion broker %s: %w\\\", broker, err)\"\n  bit,1 failed,restored yes: TestBrokerRefusalSkipsReplacement,\"replacement_test.go:80: trace = [create mode write file-sync file-close rename directory-open directory-sync directory-close], want no replacement stage\"\n\nProbe 7 (D30 refusal writes a manifest). Same old; new (exact): \"\\t\\t_ = os.WriteFile(filepath.Join(dir, Name), nil, 0o644)\\n\\t\\treturn \\\"\\\", \\\"\\\", fmt.Errorf(\\\"digest promotion broker %s: %w\\\", broker, err)\"\n  bit,2 failed,restored yes: TestBrokerRefusalSkipsReplacement,\"replacement_test.go:82: .../bench-broker.manifest stat err = <nil>, want it absent\"; TestBrokerWriteReturnsRefusal,\"replacement_test.go:93: .../bench-broker.manifest stat err = <nil>, want it absent\"\n\nProbe 8 (S1/C1 prefix drift). Swap old (exact): \"const TemporaryPrefix = \\\".\\\" + Name + \\\".\\\"\" new (exact): \"const TemporaryPrefix = \\\".bench-broker.\\\"\"\n  bit,1 failed,restored yes: TestBrokerReplacement,\"replacement_test.go:49: staging file \\\".bench-broker.manifest.2834163907\\\", want it named \\\".bench-broker.\\\" plus a random suffix\"\n  Debug note: at 2fd8872c the pin used strings.HasPrefix and this probe was silent (passed, 0 failed, restored yes), because \".bench-broker.\" is still a prefix of the staging name. 19cb2923 compares the staging name up to its last dot.\n\nProbe 9 (C2 public Write swallows every error). Swap old (exact): \"\\treturn write(durablefile.Native(), dir, broker, version)\" new (exact): \"\\tpath, resolved, _ := write(durablefile.Native(), dir, broker, version)\\n\\treturn path, resolved, nil\"\n  bit,1 failed,restored yes: TestBrokerWriteReturnsRefusal,\"replacement_test.go:91: Write err = nil, want the refusal of a broker that is not a regular file\"\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "dcm-c2-r1-standards",
          "performer": "claude:dcm_rev_c2_std",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "c01e29c180ed8f9fa4cdff27588e8b70c0a14243",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_c2_std",
            "digest": "sha256:708dd4b541bef0ab2fc888ede93c991b4928b3067dd5084f037bca8d35c04261",
            "excerpt": "## Standards (D-C2, 275c82df..c3c67fba) — 2 findings; worst S1\nS1 | auto-fix | 7 | internal/brokermanifest/brokermanifest.go:25-29 | TemporaryPrefix restates the leaf temp naming (internal/durablefile/replace.go:97) with no production consumer; the freshness sweep would pass vacuously after a leaf drift; its comment claims the publisher owns the naming | pin the prefix against the Recorder's temp name in a test, record the red, and rewrite the comment; a leaf export is a later reviewer decision\nS2 | auto-fix (plan expansion) | 5 | internal/brokermanifest/replacement_test.go:77-82 | the Lstat-plus-ErrNotExist absence assertion repeats per caller test (prospectiveartifact:38, commitment repository:130, here) | add durablefiletest.RequireAbsent and use it here; log the closed copies for D-D\n"
          },
          "axis": "Standards",
          "base": "275c82dfc9fa03d1c87978ebd59c9be4282136bc",
          "tip": "cf1112f9cb755606a6bf6b58593958e33f2b6b2d",
          "finding_ids": [
            "S1",
            "S2"
          ],
          "supersedes": []
        },
        {
          "id": "dcm-c2-r1-spec",
          "performer": "claude:dcm_rev_c2_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "c01e29c180ed8f9fa4cdff27588e8b70c0a14243",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_c2_spec",
            "digest": "sha256:500999f484b1702e18af0a0c7aa7ebebfc01e48784758d64aee138e9e3dbcc9a",
            "excerpt": "## Spec (D-C2, 275c82df..c3c67fba) — 0 findings\nD25 holds: digest refusal and tab-framed fields at 0644; owner syncs file and directory; adopt and freshness callers and exclusions unchanged; narrow write(replacer, ...) seam; D30 refusal witness; D36 pre-rename evidence. Plan commit c4d1c95f adds only the 11 assignment.\nAdvice: the caller-package passes are claimed but not shown in the native excerpt.\n"
          },
          "axis": "Spec",
          "base": "275c82dfc9fa03d1c87978ebd59c9be4282136bc",
          "tip": "cf1112f9cb755606a6bf6b58593958e33f2b6b2d",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "dcm-c2-r1-coverage",
          "performer": "claude:dcm_rev_c2_cov",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "c01e29c180ed8f9fa4cdff27588e8b70c0a14243",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_c2_cov",
            "digest": "sha256:df09615d4498fec1997a33638446463d89ff4e65382688795ac9a524bec2c8ac",
            "excerpt": "## Coverage (D-C2, 275c82df..c3c67fba) — 3 findings; worst C1. Probes each restored=yes.\nC1 | auto-fix (same as S1) | 9 | internal/brokermanifest/brokermanifest.go:29 | const TemporaryPrefix = \".bench-broker.\" is silent in brokermanifest and freshness | assert the Recorder temp name starts with TemporaryPrefix; record the red\nC2 | auto-fix | 8 | internal/brokermanifest/brokermanifest.go:37 | a public Write that swallows every error is silent in brokermanifest, adopt, and freshness | test that Write(dir, dir, version) returns an error and leaves no manifest; record the swallow red\nC3 | ask-user (orchestrator: no-op, pre-existing and outside D25; logged as a learning) | 6 | brokermanifest.go:42-44 | dropping symlink resolution is silent | a later ticket may pin the resolved path\n"
          },
          "axis": "Coverage",
          "base": "275c82dfc9fa03d1c87978ebd59c9be4282136bc",
          "tip": "cf1112f9cb755606a6bf6b58593958e33f2b6b2d",
          "finding_ids": [
            "C1",
            "C2",
            "C3"
          ],
          "supersedes": []
        },
        {
          "id": "dcm-c2-r2-standards",
          "performer": "claude:dcm_rev_c2_std2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "de2e4bde5f117c8f461ebe3ebf1bbbdb21a05cff",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_c2_std2",
            "digest": "sha256:751fe13ca9f7e99fb2919c8894c8a534d41fd90d6a39932878cd20207b9e0440",
            "excerpt": "## Standards confirming round (D-C2, repair delta 7b7e9cfd..5d295f64) — 1 finding; S1, S2 closed\nS1 closed: replacement_test.go:45-50 pins TemporaryPrefix against the staging name cut at its last dot; the constant comment states the current fact. S2 closed: RequireAbsent (assert.go:40-46) replaces the local helper.\nS3 | auto-fix (comment-only, evidence-only correction) | 6 | internal/brokermanifest/replacement_test.go:85 | the comment names the refusal of the native replacement path, but the refusal comes from Digest before any replacement | reword to the digest refusal of a broker that is not a regular file\n"
          },
          "axis": "Standards",
          "base": "275c82dfc9fa03d1c87978ebd59c9be4282136bc",
          "tip": "19cb292348cf4aae0f895a139ec39cf8a6ba8edb",
          "finding_ids": [
            "S3"
          ],
          "supersedes": [
            "dcm-c2-r1-standards"
          ]
        },
        {
          "id": "dcm-c2-r2-spec",
          "performer": "claude:dcm_rev_c2_spec2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "de2e4bde5f117c8f461ebe3ebf1bbbdb21a05cff",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_c2_spec2",
            "digest": "sha256:b4a200142e0f7d9bb0409ecf63b0dbbbf7bbcdf53edfd2ab18ff1aaf374eb5d1",
            "excerpt": "## Spec confirming round (D-C2, repair delta 7b7e9cfd..5d295f64 and plan d723f0c3) — 0 findings\nD25, D30, and the D36 evidence intact; TestBrokerWriteReturnsRefusal meets D30 for the exported Write; plan d723f0c3 adds assert.go and one floor row only.\n"
          },
          "axis": "Spec",
          "base": "275c82dfc9fa03d1c87978ebd59c9be4282136bc",
          "tip": "19cb292348cf4aae0f895a139ec39cf8a6ba8edb",
          "finding_ids": [],
          "supersedes": [
            "dcm-c2-r1-spec"
          ]
        },
        {
          "id": "dcm-c2-r2-coverage",
          "performer": "claude:dcm_rev_c2_cov2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "de2e4bde5f117c8f461ebe3ebf1bbbdb21a05cff",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_c2_cov2",
            "digest": "sha256:b2a85b190064d99408b820fa65506984e21f6296fe53f0fc76b855a7afa68afc",
            "excerpt": "## Coverage confirming round (D-C2, repair delta 7b7e9cfd..5d295f64) — 0 findings\nC1 closed: new probe \".\" + Name (no trailing dot): bit at replacement_test.go:49, restored=yes. C2 closed: new probe, Write swallows every error except replacement errors: bit at :91, restored=yes.\nAdvice: a public Write that hides only replacement errors is silent; the native replacement cannot be faulted through Write.\n"
          },
          "axis": "Coverage",
          "base": "275c82dfc9fa03d1c87978ebd59c9be4282136bc",
          "tip": "19cb292348cf4aae0f895a139ec39cf8a6ba8edb",
          "finding_ids": [],
          "supersedes": [
            "dcm-c2-r1-coverage"
          ]
        },
        {
          "id": "dcm-c2-r3-standards",
          "performer": "claude:dcm_rev_c2_std3",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "de2e4bde5f117c8f461ebe3ebf1bbbdb21a05cff",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_c2_std3",
            "digest": "sha256:31bc929c23969ad08992b88ca6f5a6c52be97c81762833b97ac3b9ae5ffa88b6",
            "excerpt": "## Standards reaffirmation (D-C2, comment-only commit dccf92f2..b92b47f9) — 0 findings\nS3 closed: internal/brokermanifest/replacement_test.go:85-86 now names the digest refusal; Digest refuses the directory at brokermanifest.go:92-93 before Replace at :53; the commit changes only the comment.\n"
          },
          "axis": "Standards",
          "base": "275c82dfc9fa03d1c87978ebd59c9be4282136bc",
          "tip": "19cb292348cf4aae0f895a139ec39cf8a6ba8edb",
          "finding_ids": [],
          "supersedes": [
            "dcm-c2-r2-standards"
          ]
        }
      ]
    },
    {
      "id": "D-C4",
      "base": "19cb292348cf4aae0f895a139ec39cf8a6ba8edb",
      "tip": "c377e632d4056ebb6853ed1087efbb38e010c824",
      "plan_digest": "sha256:5f3b8ff2cb08f262254461d3182e5bd26885f2f80ef93f6792bdd716587ef9a2",
      "source_digest": "808f0b78dfb2cc543b4bcadac7683c7900b809d4",
      "acceptance_rows": [
        "D27",
        "D34",
        "D41",
        "D42",
        "D43"
      ],
      "verification": [
        {
          "id": "dcm-t11-caller-tests",
          "performer": "claude:dcm_t11",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e6ce538abb8886a1123ff7ec066fac3b7e74dbef",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "0e86fec98de40c9c5ff0e21cd719255731934188",
            "digest": "sha256:84207ed2f0c273cc3bf1d8ee2d2fa5881ea6cac7ce642e70ea87a03e15d035d7",
            "excerpt": "bench test --package ./internal/probe at 0e86fec98de40c9c5ff0e21cd719255731934188 (FT354-durable-callers, clean): internal/probe pass, 114 tests run, failures[0], skips[0].\nThe five ticket 11 witnesses are in internal/probe/replacement_test.go: TestProbeReplacement (D27), TestProbeFailedRestore (D34), TestProbeMutationRecoveryFailure (D41), TestProbeMutationRecoverySuccess (D42, over durablefiletest.PostRenameStages), TestProbeMutationPreparationFailure (D43).\n"
          },
          "requirement": "caller-tests",
          "command": "bench test --package ./internal/probe",
          "exit_code": 0
        },
        {
          "id": "dcm-t11-command-contract",
          "performer": "claude:dcm_t11",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e6ce538abb8886a1123ff7ec066fac3b7e74dbef",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "0e86fec98de40c9c5ff0e21cd719255731934188",
            "digest": "sha256:f0086d1339026ee09e29b43e071f029aec1f94ccf95f9e6b8e0a737ad0bfe5b3",
            "excerpt": "bench test --package ./cmd/bench --run CommandRegistry at 0e86fec98de40c9c5ff0e21cd719255731934188 (FT354-durable-callers, clean): cmd/bench pass, 5 tests run, failures[0], skips[0]. Command membership, routes, and help are unchanged; probe.Command keeps its signature.\n"
          },
          "requirement": "command-contract",
          "command": "bench test --package ./cmd/bench --run 'CommandRegistry'",
          "exit_code": 0
        },
        {
          "id": "dcm-t11-axi-registry",
          "performer": "claude:dcm_t11",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e6ce538abb8886a1123ff7ec066fac3b7e74dbef",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "0e86fec98de40c9c5ff0e21cd719255731934188",
            "digest": "sha256:43de5ffaeb3b65b18c344b766a861714566be79f85728a72abeacf1271a8b4d7",
            "excerpt": "bench test --check axi-query-registry at 0e86fec98de40c9c5ff0e21cd719255731934188 (FT354-durable-callers, clean): internal/conformance pass, 1 test run, failures[0], skips[0].\n"
          },
          "requirement": "axi-registry",
          "command": "bench test --check axi-query-registry",
          "exit_code": 0
        },
        {
          "id": "dcm-t11-routing",
          "performer": "claude:dcm_t11",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e6ce538abb8886a1123ff7ec066fac3b7e74dbef",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "0e86fec98de40c9c5ff0e21cd719255731934188",
            "digest": "sha256:bcb1fd2fa80dea84f99ea2aeabfd2ddc67bf87737674c80c81695de698e8db4a",
            "excerpt": "bench test --check subcommand-routing at 0e86fec98de40c9c5ff0e21cd719255731934188 (FT354-durable-callers, clean): internal/conformance pass, 1 test run, failures[0], skips[0].\n"
          },
          "requirement": "routing",
          "command": "bench test --check subcommand-routing",
          "exit_code": 0
        },
        {
          "id": "dcm-t11-cmd-bench",
          "performer": "claude:dcm_t11",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e6ce538abb8886a1123ff7ec066fac3b7e74dbef",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "0e86fec98de40c9c5ff0e21cd719255731934188",
            "digest": "sha256:8829a137cafc68e46ca1588f5a923ce60d72cfba21352ac6914573316fc3201b",
            "excerpt": "bench test --package ./cmd/bench at 0e86fec98de40c9c5ff0e21cd719255731934188 (FT354-durable-callers, clean): cmd/bench pass, 335 tests run, failures[0], skips[0].\n"
          },
          "requirement": "cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "dcm-t11-conformance",
          "performer": "claude:dcm_t11",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e6ce538abb8886a1123ff7ec066fac3b7e74dbef",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "0e86fec98de40c9c5ff0e21cd719255731934188",
            "digest": "sha256:047d82c2bbb7877493df1ecb7315f48f602d45e1b3549600a2fff7f7c2662dd9",
            "excerpt": "bench test --package ./internal/conformance at 0e86fec98de40c9c5ff0e21cd719255731934188 (FT354-durable-callers, clean): internal/conformance pass, 1690 tests run, failures[0], skips[1]: TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device (capability: privilege: cannot create a character device: operation not permitted), an environment capability skip not owned by this ticket.\n"
          },
          "requirement": "conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "dcm-t11-probetest",
          "performer": "claude:dcm_t11",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e6ce538abb8886a1123ff7ec066fac3b7e74dbef",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "0e86fec98de40c9c5ff0e21cd719255731934188",
            "digest": "sha256:d5f1d4d597721dddb7117b80104534d80b9cb9e873aa55e6f4beb55db39903aa",
            "excerpt": "bench test --package ./internal/probe/probetest at 0e86fec98de40c9c5ff0e21cd719255731934188 (FT354-durable-callers, clean): internal/probe/probetest no-tests, 0 tests run, failures[0], skips[0]. The fixture is unchanged: startCount and requireBaselineOnly in the probe package already separate the baseline start from the focused start.\n"
          },
          "requirement": "probetest",
          "command": "bench test --package ./internal/probe/probetest",
          "exit_code": 0
        },
        {
          "id": "dcm-t11-root-conformance",
          "performer": "claude:dcm_t11",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e6ce538abb8886a1123ff7ec066fac3b7e74dbef",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "0e86fec98de40c9c5ff0e21cd719255731934188",
            "digest": "sha256:3be31271cf580676dfbe2729d81ac0dc1ef22486cb7ff18dcfebecad25ae04e5",
            "excerpt": "bench test --check root-conformance at 0e86fec98de40c9c5ff0e21cd719255731934188 (FT354-durable-callers, clean): internal/conformance pass, 1 test run, failures[0], skips[0].\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "dcm-t11-omission-proof",
          "performer": "claude:dcm_t11",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e6ce538abb8886a1123ff7ec066fac3b7e74dbef",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "0e86fec98de40c9c5ff0e21cd719255731934188",
            "digest": "sha256:1ea91b2a2f61d60f38c301e05c1b4f15d1f8e5e779e706e3f5ff03d006d933a2",
            "excerpt": "D-C4 omission-proof probes at 0e86fec98de40c9c5ff0e21cd719255731934188 (FT354-durable-callers), bench probe ... --package ./internal/probe, no --run filter. Every verdict bit, every restored yes; the package passes 114 tests after each restore.\nNamed probe: release preservation on every mutation error, or treat failed restore sync as success; the preservation-path and focused-start witnesses detect each mutation.\n\nO1 (central self-probe: release the copy on the failed-restore exit). internal/probe/probe.go\nSwap old (exact): \"return line + block, 2\"\nSwap new (exact): \"preserved.release(); return line + block, 2\"\nVerdict bit, 1 failed, restored yes. TestProbeMutationRecoveryFailure replacement_test.go:104: preserved files = ([], <nil>), want exactly one.\n\nO2 (release without restore on every mutation error). internal/probe/probe.go\nSwap old: \"if errors.As(err, &replaceErr) && replaceErr.Published {\"\nSwap new: \"if false && errors.As(err, &replaceErr) && replaceErr.Published {\"\nVerdict bit, 4 failed, restored yes. TestProbeMutationRecoveryFailure :102 exit = 1, want 2; TestProbeMutationRecoverySuccess/directory-open, /directory-sync, /directory-close :83 subject holds the mutated bytes.\n\nO3 (treat failed restore sync as success). internal/probe/subject.go\nSwap old: \"if err := replacer.Replace(subject.path, saved, subject.mode); err != nil {\"\nSwap new: \"if err := replacer.Replace(subject.path, saved, subject.mode); err != nil && !strings.Contains(err.Error(), \\\"rename completed\\\") {\"\nVerdict bit, 1 failed, restored yes. TestProbeFailedRestore :123 row bit,...,yes in place of restore-failed at exit 2.\n\nO4 (owner bypass, mutation writer; D27). internal/probe/probe.go\nSwap old: \"if err := replacer.Replace(subject.path, mutated, subject.mode); err != nil {\"\nSwap new: \"if err := durablefile.Replace(subject.path, mutated, subject.mode); err != nil {\"\nVerdict bit, 7 failed, restored yes. TestProbeReplacement :58 (the injected failure is hidden; the first routed rename is the restore); TestProbeMutationRecoverySuccess x3 :84 Bench home holds 1 entries.\n\nO5 (owner bypass, restore writer). internal/probe/subject.go\nSwap old: \"if err := replacer.Replace(subject.path, saved, subject.mode); err != nil {\"\nSwap new: \"if err := durablefile.Replace(subject.path, saved, subject.mode); err != nil {\"\nVerdict bit, 2 failed, restored yes. TestProbeFailedRestore :123; TestProbeMutationRecoveryFailure :102 exit = 1, want 2.\n\nO6 (classification swap; Published replaced by a stage literal). internal/probe/probe.go\nSwap old: \"errors.As(err, &replaceErr) && replaceErr.Published {\"\nSwap new: \"errors.As(err, &replaceErr) && replaceErr.Stage == \\\"directory-sync\\\" {\"\nVerdict bit, 2 failed, restored yes. TestProbeMutationRecoverySuccess/directory-open and /directory-close :83 subject holds the mutated bytes.\n\nO7 (focused-start witness: a focused run starts after a failed mutation). internal/probe/probe.go\nSwap old: \"return mutationFailure(replacer, subject, preserved, err)\"\nSwap new: \"testreport.Execute(root, request); return mutationFailure(replacer, subject, preserved, err)\"\nVerdict bit, 2 failed, restored yes. TestProbeReplacement :69 Go test starts = 2, want only the baseline's one; existing PB45 TestProbeRefusesWhenTheMutationWriteFails :286 same.\n\nClassification note: mutationFailure classifies inline with errors.As(err, &replaceErr) && replaceErr.Published. The probe chooses between an immediate copy release and a restore-first path from that fact. The same derivation is the private published() in internal/intent/transaction.go; a shared durablefile predicate is outside this ticket's fence, and review decides where it belongs.\n"
          },
          "requirement": "omission-proof",
          "command": "bench test --package ./internal/probe",
          "exit_code": 1,
          "probe": {
            "mutation": "Release preservation on every mutation error or treat failed restore sync as success. The preservation-path and focused-start witnesses must detect each mutation.",
            "outcome": "bit",
            "exit_code": 0,
            "restore": "pass",
            "native_ref": {
              "ref": "0e86fec98de40c9c5ff0e21cd719255731934188",
              "digest": "sha256:1ea91b2a2f61d60f38c301e05c1b4f15d1f8e5e779e706e3f5ff03d006d933a2",
              "excerpt": "D-C4 omission-proof probes at 0e86fec98de40c9c5ff0e21cd719255731934188 (FT354-durable-callers), bench probe ... --package ./internal/probe, no --run filter. Every verdict bit, every restored yes; the package passes 114 tests after each restore.\nNamed probe: release preservation on every mutation error, or treat failed restore sync as success; the preservation-path and focused-start witnesses detect each mutation.\n\nO1 (central self-probe: release the copy on the failed-restore exit). internal/probe/probe.go\nSwap old (exact): \"return line + block, 2\"\nSwap new (exact): \"preserved.release(); return line + block, 2\"\nVerdict bit, 1 failed, restored yes. TestProbeMutationRecoveryFailure replacement_test.go:104: preserved files = ([], <nil>), want exactly one.\n\nO2 (release without restore on every mutation error). internal/probe/probe.go\nSwap old: \"if errors.As(err, &replaceErr) && replaceErr.Published {\"\nSwap new: \"if false && errors.As(err, &replaceErr) && replaceErr.Published {\"\nVerdict bit, 4 failed, restored yes. TestProbeMutationRecoveryFailure :102 exit = 1, want 2; TestProbeMutationRecoverySuccess/directory-open, /directory-sync, /directory-close :83 subject holds the mutated bytes.\n\nO3 (treat failed restore sync as success). internal/probe/subject.go\nSwap old: \"if err := replacer.Replace(subject.path, saved, subject.mode); err != nil {\"\nSwap new: \"if err := replacer.Replace(subject.path, saved, subject.mode); err != nil && !strings.Contains(err.Error(), \\\"rename completed\\\") {\"\nVerdict bit, 1 failed, restored yes. TestProbeFailedRestore :123 row bit,...,yes in place of restore-failed at exit 2.\n\nO4 (owner bypass, mutation writer; D27). internal/probe/probe.go\nSwap old: \"if err := replacer.Replace(subject.path, mutated, subject.mode); err != nil {\"\nSwap new: \"if err := durablefile.Replace(subject.path, mutated, subject.mode); err != nil {\"\nVerdict bit, 7 failed, restored yes. TestProbeReplacement :58 (the injected failure is hidden; the first routed rename is the restore); TestProbeMutationRecoverySuccess x3 :84 Bench home holds 1 entries.\n\nO5 (owner bypass, restore writer). internal/probe/subject.go\nSwap old: \"if err := replacer.Replace(subject.path, saved, subject.mode); err != nil {\"\nSwap new: \"if err := durablefile.Replace(subject.path, saved, subject.mode); err != nil {\"\nVerdict bit, 2 failed, restored yes. TestProbeFailedRestore :123; TestProbeMutationRecoveryFailure :102 exit = 1, want 2.\n\nO6 (classification swap; Published replaced by a stage literal). internal/probe/probe.go\nSwap old: \"errors.As(err, &replaceErr) && replaceErr.Published {\"\nSwap new: \"errors.As(err, &replaceErr) && replaceErr.Stage == \\\"directory-sync\\\" {\"\nVerdict bit, 2 failed, restored yes. TestProbeMutationRecoverySuccess/directory-open and /directory-close :83 subject holds the mutated bytes.\n\nO7 (focused-start witness: a focused run starts after a failed mutation). internal/probe/probe.go\nSwap old: \"return mutationFailure(replacer, subject, preserved, err)\"\nSwap new: \"testreport.Execute(root, request); return mutationFailure(replacer, subject, preserved, err)\"\nVerdict bit, 2 failed, restored yes. TestProbeReplacement :69 Go test starts = 2, want only the baseline's one; existing PB45 TestProbeRefusesWhenTheMutationWriteFails :286 same.\n\nClassification note: mutationFailure classifies inline with errors.As(err, &replaceErr) && replaceErr.Published. The probe chooses between an immediate copy release and a restore-first path from that fact. The same derivation is the private published() in internal/intent/transaction.go; a shared durablefile predicate is outside this ticket's fence, and review decides where it belongs.\n"
            }
          }
        },
        {
          "id": "dcm-t11-preservation",
          "performer": "claude:dcm_t11",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e6ce538abb8886a1123ff7ec066fac3b7e74dbef",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "0e86fec98de40c9c5ff0e21cd719255731934188",
            "digest": "sha256:41e55eecdcb1e8d6fd32845cf2ae9937bfaf8f46c40a855a3658754e2ae8b34d",
            "excerpt": "D-C4 preservation probes at 0e86fec98de40c9c5ff0e21cd719255731934188 (FT354-durable-callers), bench probe ... --package ./internal/probe, no --run filter. One necessity red for each independent check of the new tests, plus the payload and mode of each migrated writer and the native wrapper's error path. Every verdict bit, every restored yes; the package passes 114 tests after each restore. The omission-proof entry holds O1-O7; their reds also serve D41 copy existence (O1), D42 bytes (O2, O6), D42 home (O4), D27 focused start (O7), and D34 row (O3).\n\nP1 D27 output equality (the mutation cause reaches the caller). probe.go\nOld: \"line := toon.Errorf(\\\"probe mutation failed\\\", err.Error()) + \\\"\\\\n\\\"\"\nNew: \"line := toon.Errorf(\\\"probe mutation failed\\\", \\\"the subject replacement failed\\\") + \\\"\\\\n\\\"\"\nbit, 1 failed, restored yes. TestProbeReplacement :67 only.\n\nP2 D27 mutation mode (constant drift). probe.go\nOld: \"replacer.Replace(subject.path, mutated, subject.mode)\"\nNew: \"replacer.Replace(subject.path, mutated, 0o644)\"\nbit, 1 failed, restored yes. TestProbeReplacement :58 mode = -rw-r--r--, want -rw-r----- only.\n\nP3 D27 mutation payload. probe.go\nOld: \"replacer.Replace(subject.path, mutated, subject.mode)\"\nNew: \"replacer.Replace(subject.path, append(mutated, 10), subject.mode)\"\nbit, 2 failed, restored yes. TestProbeReplacement :58 bytes; TestOmitFileUsesExactBytesAndReportsOmit.\n\nP4 restore mode (owned by existing PB17). subject.go\nOld: \"replacer.Replace(subject.path, saved, subject.mode)\"\nNew: \"replacer.Replace(subject.path, saved, 0o644)\"\nbit, 1 failed, restored yes. TestProbeKeepsTheSubjectMode outcome_test.go:130 only.\n\nP5 D34 copy retained after a failed restore sync. probe.go\nOld: \"\\tif restored {\\n\\t\\tpreserved.release()\\n\\t}\\n\\treturn render(\"\nNew: \"\\tif true {\\n\\t\\tpreserved.release()\\n\\t}\\n\\treturn render(\"\nbit, 3 failed, restored yes. TestProbeFailedRestore :124 preserved files = ([]); PB14, PB47.\n\nP6 D41 preserved row names the copy path. probe.go\nOld: \"[][]string{{preserved.file, sanitize.Strip(reason)}}\"\nNew: \"[][]string{{preserved.dir, sanitize.Strip(reason)}}\"\nbit, 3 failed, restored yes. TestProbeMutationRecoveryFailure :108; PB14, PB47.\n\nP7 D41 reason carries the mutation cause. probe.go\nOld: \"err.Error()+\\\"; restore: \\\"+reason\"  New: \"\\\"restore: \\\"+reason\"\nbit, 1 failed, restored yes. TestProbeMutationRecoveryFailure :108 only.\n\nP8 D41 reason carries the restore cause. probe.go\nOld: \"err.Error()+\\\"; restore: \\\"+reason\"  New: \"err.Error()+\\\"; restore: \\\"+reason[:0]\"\nbit, 1 failed, restored yes. TestProbeMutationRecoveryFailure :108 only.\n\nP9 D41 exit 2. probe.go\nOld: \"return line + block, 2\"  New: \"return line + block, 1\"\nbit, 1 failed, restored yes. TestProbeMutationRecoveryFailure :102 only.\n\nP10 D41 exact copy bytes (shared with PB12, PB14). subject.go\nOld: \"if err := os.WriteFile(file, subject.start, 0o600); err != nil {\"\nNew: \"if err := os.WriteFile(file, subject.start[1:], 0o600); err != nil {\"\nbit, 29 failed, restored yes. Includes TestProbeMutationRecoveryFailure copy bytes and TestProbeReportsARestoreFailure.\n\nP11 D42 home: a proven restore keeps the copy. probe.go\nOld: \"if restored, reason := restore(replacer, subject, preserved); !restored {\"\nNew: \"if restored, reason := restore(replacer, subject, preserved); restored { return line, 1 } else {\"\nbit, 3 failed, restored yes. TestProbeMutationRecoverySuccess/directory-open, -sync, -close :84 only.\n\nP12 D43 old bytes: fallback write of the new bytes after an unpublished error. probe.go\nOld: \"return mutationFailure(replacer, subject, preserved, err)\"\nNew: \"_ = os.WriteFile(subject.path, mutated, subject.mode); return mutationFailure(replacer, subject, preserved, err)\"\nbit, 2 failed, restored yes. TestProbeMutationPreparationFailure :135 old-bytes check; existing PB45 :284.\n\nP13 D43 home: the rename stage is classified as published. probe.go\nOld: \"errors.As(err, &replaceErr) && replaceErr.Published {\"\nNew: \"errors.As(err, &replaceErr) && (replaceErr.Published || replaceErr.Stage == \\\"rename\\\") {\"\nbit, 1 failed, restored yes. TestProbeMutationPreparationFailure :136 Bench home holds 1 entries only.\n\nP14 native wrapper error path (run supplies durablefile.Native(); PB45 drives it through Command). probe.go\nOld: \"\\treturn runWith(durablefile.Native(), root, parsed)\"\nNew: \"\\tout, code := runWith(durablefile.Native(), root, parsed)\\n\\tif len(out) > 28 && out[:28] == \\\"error: probe mutation failed\\\" {\\n\\t\\treturn \\\"\\\", 0\\n\\t}\\n\\treturn out, code\"\nbit, 1 failed, restored yes. TestProbeRefusesWhenTheMutationWriteFails refusal_test.go:282 only.\n"
          },
          "requirement": "preservation",
          "command": "bench test --package ./internal/probe",
          "exit_code": 1,
          "probe": {
            "mutation": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass.",
            "outcome": "bit",
            "exit_code": 0,
            "restore": "pass",
            "native_ref": {
              "ref": "0e86fec98de40c9c5ff0e21cd719255731934188",
              "digest": "sha256:41e55eecdcb1e8d6fd32845cf2ae9937bfaf8f46c40a855a3658754e2ae8b34d",
              "excerpt": "D-C4 preservation probes at 0e86fec98de40c9c5ff0e21cd719255731934188 (FT354-durable-callers), bench probe ... --package ./internal/probe, no --run filter. One necessity red for each independent check of the new tests, plus the payload and mode of each migrated writer and the native wrapper's error path. Every verdict bit, every restored yes; the package passes 114 tests after each restore. The omission-proof entry holds O1-O7; their reds also serve D41 copy existence (O1), D42 bytes (O2, O6), D42 home (O4), D27 focused start (O7), and D34 row (O3).\n\nP1 D27 output equality (the mutation cause reaches the caller). probe.go\nOld: \"line := toon.Errorf(\\\"probe mutation failed\\\", err.Error()) + \\\"\\\\n\\\"\"\nNew: \"line := toon.Errorf(\\\"probe mutation failed\\\", \\\"the subject replacement failed\\\") + \\\"\\\\n\\\"\"\nbit, 1 failed, restored yes. TestProbeReplacement :67 only.\n\nP2 D27 mutation mode (constant drift). probe.go\nOld: \"replacer.Replace(subject.path, mutated, subject.mode)\"\nNew: \"replacer.Replace(subject.path, mutated, 0o644)\"\nbit, 1 failed, restored yes. TestProbeReplacement :58 mode = -rw-r--r--, want -rw-r----- only.\n\nP3 D27 mutation payload. probe.go\nOld: \"replacer.Replace(subject.path, mutated, subject.mode)\"\nNew: \"replacer.Replace(subject.path, append(mutated, 10), subject.mode)\"\nbit, 2 failed, restored yes. TestProbeReplacement :58 bytes; TestOmitFileUsesExactBytesAndReportsOmit.\n\nP4 restore mode (owned by existing PB17). subject.go\nOld: \"replacer.Replace(subject.path, saved, subject.mode)\"\nNew: \"replacer.Replace(subject.path, saved, 0o644)\"\nbit, 1 failed, restored yes. TestProbeKeepsTheSubjectMode outcome_test.go:130 only.\n\nP5 D34 copy retained after a failed restore sync. probe.go\nOld: \"\\tif restored {\\n\\t\\tpreserved.release()\\n\\t}\\n\\treturn render(\"\nNew: \"\\tif true {\\n\\t\\tpreserved.release()\\n\\t}\\n\\treturn render(\"\nbit, 3 failed, restored yes. TestProbeFailedRestore :124 preserved files = ([]); PB14, PB47.\n\nP6 D41 preserved row names the copy path. probe.go\nOld: \"[][]string{{preserved.file, sanitize.Strip(reason)}}\"\nNew: \"[][]string{{preserved.dir, sanitize.Strip(reason)}}\"\nbit, 3 failed, restored yes. TestProbeMutationRecoveryFailure :108; PB14, PB47.\n\nP7 D41 reason carries the mutation cause. probe.go\nOld: \"err.Error()+\\\"; restore: \\\"+reason\"  New: \"\\\"restore: \\\"+reason\"\nbit, 1 failed, restored yes. TestProbeMutationRecoveryFailure :108 only.\n\nP8 D41 reason carries the restore cause. probe.go\nOld: \"err.Error()+\\\"; restore: \\\"+reason\"  New: \"err.Error()+\\\"; restore: \\\"+reason[:0]\"\nbit, 1 failed, restored yes. TestProbeMutationRecoveryFailure :108 only.\n\nP9 D41 exit 2. probe.go\nOld: \"return line + block, 2\"  New: \"return line + block, 1\"\nbit, 1 failed, restored yes. TestProbeMutationRecoveryFailure :102 only.\n\nP10 D41 exact copy bytes (shared with PB12, PB14). subject.go\nOld: \"if err := os.WriteFile(file, subject.start, 0o600); err != nil {\"\nNew: \"if err := os.WriteFile(file, subject.start[1:], 0o600); err != nil {\"\nbit, 29 failed, restored yes. Includes TestProbeMutationRecoveryFailure copy bytes and TestProbeReportsARestoreFailure.\n\nP11 D42 home: a proven restore keeps the copy. probe.go\nOld: \"if restored, reason := restore(replacer, subject, preserved); !restored {\"\nNew: \"if restored, reason := restore(replacer, subject, preserved); restored { return line, 1 } else {\"\nbit, 3 failed, restored yes. TestProbeMutationRecoverySuccess/directory-open, -sync, -close :84 only.\n\nP12 D43 old bytes: fallback write of the new bytes after an unpublished error. probe.go\nOld: \"return mutationFailure(replacer, subject, preserved, err)\"\nNew: \"_ = os.WriteFile(subject.path, mutated, subject.mode); return mutationFailure(replacer, subject, preserved, err)\"\nbit, 2 failed, restored yes. TestProbeMutationPreparationFailure :135 old-bytes check; existing PB45 :284.\n\nP13 D43 home: the rename stage is classified as published. probe.go\nOld: \"errors.As(err, &replaceErr) && replaceErr.Published {\"\nNew: \"errors.As(err, &replaceErr) && (replaceErr.Published || replaceErr.Stage == \\\"rename\\\") {\"\nbit, 1 failed, restored yes. TestProbeMutationPreparationFailure :136 Bench home holds 1 entries only.\n\nP14 native wrapper error path (run supplies durablefile.Native(); PB45 drives it through Command). probe.go\nOld: \"\\treturn runWith(durablefile.Native(), root, parsed)\"\nNew: \"\\tout, code := runWith(durablefile.Native(), root, parsed)\\n\\tif len(out) > 28 && out[:28] == \\\"error: probe mutation failed\\\" {\\n\\t\\treturn \\\"\\\", 0\\n\\t}\\n\\treturn out, code\"\nbit, 1 failed, restored yes. TestProbeRefusesWhenTheMutationWriteFails refusal_test.go:282 only.\n"
            }
          }
        },
        {
          "id": "dcm-t11-omission-proof-r2",
          "performer": "claude:dcm_t11",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e6ce538abb8886a1123ff7ec066fac3b7e74dbef",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "0e86fec98de40c9c5ff0e21cd719255731934188",
            "digest": "sha256:1ea91b2a2f61d60f38c301e05c1b4f15d1f8e5e779e706e3f5ff03d006d933a2",
            "excerpt": "D-C4 omission-proof probes at 0e86fec98de40c9c5ff0e21cd719255731934188 (FT354-durable-callers), bench probe ... --package ./internal/probe, no --run filter. Every verdict bit, every restored yes; the package passes 114 tests after each restore.\nNamed probe: release preservation on every mutation error, or treat failed restore sync as success; the preservation-path and focused-start witnesses detect each mutation.\n\nO1 (central self-probe: release the copy on the failed-restore exit). internal/probe/probe.go\nSwap old (exact): \"return line + block, 2\"\nSwap new (exact): \"preserved.release(); return line + block, 2\"\nVerdict bit, 1 failed, restored yes. TestProbeMutationRecoveryFailure replacement_test.go:104: preserved files = ([], <nil>), want exactly one.\n\nO2 (release without restore on every mutation error). internal/probe/probe.go\nSwap old: \"if errors.As(err, &replaceErr) && replaceErr.Published {\"\nSwap new: \"if false && errors.As(err, &replaceErr) && replaceErr.Published {\"\nVerdict bit, 4 failed, restored yes. TestProbeMutationRecoveryFailure :102 exit = 1, want 2; TestProbeMutationRecoverySuccess/directory-open, /directory-sync, /directory-close :83 subject holds the mutated bytes.\n\nO3 (treat failed restore sync as success). internal/probe/subject.go\nSwap old: \"if err := replacer.Replace(subject.path, saved, subject.mode); err != nil {\"\nSwap new: \"if err := replacer.Replace(subject.path, saved, subject.mode); err != nil && !strings.Contains(err.Error(), \\\"rename completed\\\") {\"\nVerdict bit, 1 failed, restored yes. TestProbeFailedRestore :123 row bit,...,yes in place of restore-failed at exit 2.\n\nO4 (owner bypass, mutation writer; D27). internal/probe/probe.go\nSwap old: \"if err := replacer.Replace(subject.path, mutated, subject.mode); err != nil {\"\nSwap new: \"if err := durablefile.Replace(subject.path, mutated, subject.mode); err != nil {\"\nVerdict bit, 7 failed, restored yes. TestProbeReplacement :58 (the injected failure is hidden; the first routed rename is the restore); TestProbeMutationRecoverySuccess x3 :84 Bench home holds 1 entries.\n\nO5 (owner bypass, restore writer). internal/probe/subject.go\nSwap old: \"if err := replacer.Replace(subject.path, saved, subject.mode); err != nil {\"\nSwap new: \"if err := durablefile.Replace(subject.path, saved, subject.mode); err != nil {\"\nVerdict bit, 2 failed, restored yes. TestProbeFailedRestore :123; TestProbeMutationRecoveryFailure :102 exit = 1, want 2.\n\nO6 (classification swap; Published replaced by a stage literal). internal/probe/probe.go\nSwap old: \"errors.As(err, &replaceErr) && replaceErr.Published {\"\nSwap new: \"errors.As(err, &replaceErr) && replaceErr.Stage == \\\"directory-sync\\\" {\"\nVerdict bit, 2 failed, restored yes. TestProbeMutationRecoverySuccess/directory-open and /directory-close :83 subject holds the mutated bytes.\n\nO7 (focused-start witness: a focused run starts after a failed mutation). internal/probe/probe.go\nSwap old: \"return mutationFailure(replacer, subject, preserved, err)\"\nSwap new: \"testreport.Execute(root, request); return mutationFailure(replacer, subject, preserved, err)\"\nVerdict bit, 2 failed, restored yes. TestProbeReplacement :69 Go test starts = 2, want only the baseline's one; existing PB45 TestProbeRefusesWhenTheMutationWriteFails :286 same.\n\nClassification note: mutationFailure classifies inline with errors.As(err, &replaceErr) && replaceErr.Published. The probe chooses between an immediate copy release and a restore-first path from that fact. The same derivation is the private published() in internal/intent/transaction.go; a shared durablefile predicate is outside this ticket's fence, and review decides where it belongs.\n"
          },
          "requirement": "omission-proof",
          "command": "bench test --package ./internal/probe",
          "exit_code": 0,
          "probe": {
            "mutation": "Release preservation on every mutation error or treat failed restore sync as success. The preservation-path and focused-start witnesses must detect each mutation.",
            "outcome": "bit",
            "exit_code": 0,
            "restore": "pass",
            "native_ref": {
              "ref": "0e86fec98de40c9c5ff0e21cd719255731934188",
              "digest": "sha256:1ea91b2a2f61d60f38c301e05c1b4f15d1f8e5e779e706e3f5ff03d006d933a2",
              "excerpt": "D-C4 omission-proof probes at 0e86fec98de40c9c5ff0e21cd719255731934188 (FT354-durable-callers), bench probe ... --package ./internal/probe, no --run filter. Every verdict bit, every restored yes; the package passes 114 tests after each restore.\nNamed probe: release preservation on every mutation error, or treat failed restore sync as success; the preservation-path and focused-start witnesses detect each mutation.\n\nO1 (central self-probe: release the copy on the failed-restore exit). internal/probe/probe.go\nSwap old (exact): \"return line + block, 2\"\nSwap new (exact): \"preserved.release(); return line + block, 2\"\nVerdict bit, 1 failed, restored yes. TestProbeMutationRecoveryFailure replacement_test.go:104: preserved files = ([], <nil>), want exactly one.\n\nO2 (release without restore on every mutation error). internal/probe/probe.go\nSwap old: \"if errors.As(err, &replaceErr) && replaceErr.Published {\"\nSwap new: \"if false && errors.As(err, &replaceErr) && replaceErr.Published {\"\nVerdict bit, 4 failed, restored yes. TestProbeMutationRecoveryFailure :102 exit = 1, want 2; TestProbeMutationRecoverySuccess/directory-open, /directory-sync, /directory-close :83 subject holds the mutated bytes.\n\nO3 (treat failed restore sync as success). internal/probe/subject.go\nSwap old: \"if err := replacer.Replace(subject.path, saved, subject.mode); err != nil {\"\nSwap new: \"if err := replacer.Replace(subject.path, saved, subject.mode); err != nil && !strings.Contains(err.Error(), \\\"rename completed\\\") {\"\nVerdict bit, 1 failed, restored yes. TestProbeFailedRestore :123 row bit,...,yes in place of restore-failed at exit 2.\n\nO4 (owner bypass, mutation writer; D27). internal/probe/probe.go\nSwap old: \"if err := replacer.Replace(subject.path, mutated, subject.mode); err != nil {\"\nSwap new: \"if err := durablefile.Replace(subject.path, mutated, subject.mode); err != nil {\"\nVerdict bit, 7 failed, restored yes. TestProbeReplacement :58 (the injected failure is hidden; the first routed rename is the restore); TestProbeMutationRecoverySuccess x3 :84 Bench home holds 1 entries.\n\nO5 (owner bypass, restore writer). internal/probe/subject.go\nSwap old: \"if err := replacer.Replace(subject.path, saved, subject.mode); err != nil {\"\nSwap new: \"if err := durablefile.Replace(subject.path, saved, subject.mode); err != nil {\"\nVerdict bit, 2 failed, restored yes. TestProbeFailedRestore :123; TestProbeMutationRecoveryFailure :102 exit = 1, want 2.\n\nO6 (classification swap; Published replaced by a stage literal). internal/probe/probe.go\nSwap old: \"errors.As(err, &replaceErr) && replaceErr.Published {\"\nSwap new: \"errors.As(err, &replaceErr) && replaceErr.Stage == \\\"directory-sync\\\" {\"\nVerdict bit, 2 failed, restored yes. TestProbeMutationRecoverySuccess/directory-open and /directory-close :83 subject holds the mutated bytes.\n\nO7 (focused-start witness: a focused run starts after a failed mutation). internal/probe/probe.go\nSwap old: \"return mutationFailure(replacer, subject, preserved, err)\"\nSwap new: \"testreport.Execute(root, request); return mutationFailure(replacer, subject, preserved, err)\"\nVerdict bit, 2 failed, restored yes. TestProbeReplacement :69 Go test starts = 2, want only the baseline's one; existing PB45 TestProbeRefusesWhenTheMutationWriteFails :286 same.\n\nClassification note: mutationFailure classifies inline with errors.As(err, &replaceErr) && replaceErr.Published. The probe chooses between an immediate copy release and a restore-first path from that fact. The same derivation is the private published() in internal/intent/transaction.go; a shared durablefile predicate is outside this ticket's fence, and review decides where it belongs.\n"
            }
          }
        },
        {
          "id": "dcm-t11-preservation-r2",
          "performer": "claude:dcm_t11",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e6ce538abb8886a1123ff7ec066fac3b7e74dbef",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "0e86fec98de40c9c5ff0e21cd719255731934188",
            "digest": "sha256:41e55eecdcb1e8d6fd32845cf2ae9937bfaf8f46c40a855a3658754e2ae8b34d",
            "excerpt": "D-C4 preservation probes at 0e86fec98de40c9c5ff0e21cd719255731934188 (FT354-durable-callers), bench probe ... --package ./internal/probe, no --run filter. One necessity red for each independent check of the new tests, plus the payload and mode of each migrated writer and the native wrapper's error path. Every verdict bit, every restored yes; the package passes 114 tests after each restore. The omission-proof entry holds O1-O7; their reds also serve D41 copy existence (O1), D42 bytes (O2, O6), D42 home (O4), D27 focused start (O7), and D34 row (O3).\n\nP1 D27 output equality (the mutation cause reaches the caller). probe.go\nOld: \"line := toon.Errorf(\\\"probe mutation failed\\\", err.Error()) + \\\"\\\\n\\\"\"\nNew: \"line := toon.Errorf(\\\"probe mutation failed\\\", \\\"the subject replacement failed\\\") + \\\"\\\\n\\\"\"\nbit, 1 failed, restored yes. TestProbeReplacement :67 only.\n\nP2 D27 mutation mode (constant drift). probe.go\nOld: \"replacer.Replace(subject.path, mutated, subject.mode)\"\nNew: \"replacer.Replace(subject.path, mutated, 0o644)\"\nbit, 1 failed, restored yes. TestProbeReplacement :58 mode = -rw-r--r--, want -rw-r----- only.\n\nP3 D27 mutation payload. probe.go\nOld: \"replacer.Replace(subject.path, mutated, subject.mode)\"\nNew: \"replacer.Replace(subject.path, append(mutated, 10), subject.mode)\"\nbit, 2 failed, restored yes. TestProbeReplacement :58 bytes; TestOmitFileUsesExactBytesAndReportsOmit.\n\nP4 restore mode (owned by existing PB17). subject.go\nOld: \"replacer.Replace(subject.path, saved, subject.mode)\"\nNew: \"replacer.Replace(subject.path, saved, 0o644)\"\nbit, 1 failed, restored yes. TestProbeKeepsTheSubjectMode outcome_test.go:130 only.\n\nP5 D34 copy retained after a failed restore sync. probe.go\nOld: \"\\tif restored {\\n\\t\\tpreserved.release()\\n\\t}\\n\\treturn render(\"\nNew: \"\\tif true {\\n\\t\\tpreserved.release()\\n\\t}\\n\\treturn render(\"\nbit, 3 failed, restored yes. TestProbeFailedRestore :124 preserved files = ([]); PB14, PB47.\n\nP6 D41 preserved row names the copy path. probe.go\nOld: \"[][]string{{preserved.file, sanitize.Strip(reason)}}\"\nNew: \"[][]string{{preserved.dir, sanitize.Strip(reason)}}\"\nbit, 3 failed, restored yes. TestProbeMutationRecoveryFailure :108; PB14, PB47.\n\nP7 D41 reason carries the mutation cause. probe.go\nOld: \"err.Error()+\\\"; restore: \\\"+reason\"  New: \"\\\"restore: \\\"+reason\"\nbit, 1 failed, restored yes. TestProbeMutationRecoveryFailure :108 only.\n\nP8 D41 reason carries the restore cause. probe.go\nOld: \"err.Error()+\\\"; restore: \\\"+reason\"  New: \"err.Error()+\\\"; restore: \\\"+reason[:0]\"\nbit, 1 failed, restored yes. TestProbeMutationRecoveryFailure :108 only.\n\nP9 D41 exit 2. probe.go\nOld: \"return line + block, 2\"  New: \"return line + block, 1\"\nbit, 1 failed, restored yes. TestProbeMutationRecoveryFailure :102 only.\n\nP10 D41 exact copy bytes (shared with PB12, PB14). subject.go\nOld: \"if err := os.WriteFile(file, subject.start, 0o600); err != nil {\"\nNew: \"if err := os.WriteFile(file, subject.start[1:], 0o600); err != nil {\"\nbit, 29 failed, restored yes. Includes TestProbeMutationRecoveryFailure copy bytes and TestProbeReportsARestoreFailure.\n\nP11 D42 home: a proven restore keeps the copy. probe.go\nOld: \"if restored, reason := restore(replacer, subject, preserved); !restored {\"\nNew: \"if restored, reason := restore(replacer, subject, preserved); restored { return line, 1 } else {\"\nbit, 3 failed, restored yes. TestProbeMutationRecoverySuccess/directory-open, -sync, -close :84 only.\n\nP12 D43 old bytes: fallback write of the new bytes after an unpublished error. probe.go\nOld: \"return mutationFailure(replacer, subject, preserved, err)\"\nNew: \"_ = os.WriteFile(subject.path, mutated, subject.mode); return mutationFailure(replacer, subject, preserved, err)\"\nbit, 2 failed, restored yes. TestProbeMutationPreparationFailure :135 old-bytes check; existing PB45 :284.\n\nP13 D43 home: the rename stage is classified as published. probe.go\nOld: \"errors.As(err, &replaceErr) && replaceErr.Published {\"\nNew: \"errors.As(err, &replaceErr) && (replaceErr.Published || replaceErr.Stage == \\\"rename\\\") {\"\nbit, 1 failed, restored yes. TestProbeMutationPreparationFailure :136 Bench home holds 1 entries only.\n\nP14 native wrapper error path (run supplies durablefile.Native(); PB45 drives it through Command). probe.go\nOld: \"\\treturn runWith(durablefile.Native(), root, parsed)\"\nNew: \"\\tout, code := runWith(durablefile.Native(), root, parsed)\\n\\tif len(out) > 28 && out[:28] == \\\"error: probe mutation failed\\\" {\\n\\t\\treturn \\\"\\\", 0\\n\\t}\\n\\treturn out, code\"\nbit, 1 failed, restored yes. TestProbeRefusesWhenTheMutationWriteFails refusal_test.go:282 only.\n"
          },
          "requirement": "preservation",
          "command": "bench test --package ./internal/probe",
          "exit_code": 0,
          "probe": {
            "mutation": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass.",
            "outcome": "bit",
            "exit_code": 0,
            "restore": "pass",
            "native_ref": {
              "ref": "0e86fec98de40c9c5ff0e21cd719255731934188",
              "digest": "sha256:41e55eecdcb1e8d6fd32845cf2ae9937bfaf8f46c40a855a3658754e2ae8b34d",
              "excerpt": "D-C4 preservation probes at 0e86fec98de40c9c5ff0e21cd719255731934188 (FT354-durable-callers), bench probe ... --package ./internal/probe, no --run filter. One necessity red for each independent check of the new tests, plus the payload and mode of each migrated writer and the native wrapper's error path. Every verdict bit, every restored yes; the package passes 114 tests after each restore. The omission-proof entry holds O1-O7; their reds also serve D41 copy existence (O1), D42 bytes (O2, O6), D42 home (O4), D27 focused start (O7), and D34 row (O3).\n\nP1 D27 output equality (the mutation cause reaches the caller). probe.go\nOld: \"line := toon.Errorf(\\\"probe mutation failed\\\", err.Error()) + \\\"\\\\n\\\"\"\nNew: \"line := toon.Errorf(\\\"probe mutation failed\\\", \\\"the subject replacement failed\\\") + \\\"\\\\n\\\"\"\nbit, 1 failed, restored yes. TestProbeReplacement :67 only.\n\nP2 D27 mutation mode (constant drift). probe.go\nOld: \"replacer.Replace(subject.path, mutated, subject.mode)\"\nNew: \"replacer.Replace(subject.path, mutated, 0o644)\"\nbit, 1 failed, restored yes. TestProbeReplacement :58 mode = -rw-r--r--, want -rw-r----- only.\n\nP3 D27 mutation payload. probe.go\nOld: \"replacer.Replace(subject.path, mutated, subject.mode)\"\nNew: \"replacer.Replace(subject.path, append(mutated, 10), subject.mode)\"\nbit, 2 failed, restored yes. TestProbeReplacement :58 bytes; TestOmitFileUsesExactBytesAndReportsOmit.\n\nP4 restore mode (owned by existing PB17). subject.go\nOld: \"replacer.Replace(subject.path, saved, subject.mode)\"\nNew: \"replacer.Replace(subject.path, saved, 0o644)\"\nbit, 1 failed, restored yes. TestProbeKeepsTheSubjectMode outcome_test.go:130 only.\n\nP5 D34 copy retained after a failed restore sync. probe.go\nOld: \"\\tif restored {\\n\\t\\tpreserved.release()\\n\\t}\\n\\treturn render(\"\nNew: \"\\tif true {\\n\\t\\tpreserved.release()\\n\\t}\\n\\treturn render(\"\nbit, 3 failed, restored yes. TestProbeFailedRestore :124 preserved files = ([]); PB14, PB47.\n\nP6 D41 preserved row names the copy path. probe.go\nOld: \"[][]string{{preserved.file, sanitize.Strip(reason)}}\"\nNew: \"[][]string{{preserved.dir, sanitize.Strip(reason)}}\"\nbit, 3 failed, restored yes. TestProbeMutationRecoveryFailure :108; PB14, PB47.\n\nP7 D41 reason carries the mutation cause. probe.go\nOld: \"err.Error()+\\\"; restore: \\\"+reason\"  New: \"\\\"restore: \\\"+reason\"\nbit, 1 failed, restored yes. TestProbeMutationRecoveryFailure :108 only.\n\nP8 D41 reason carries the restore cause. probe.go\nOld: \"err.Error()+\\\"; restore: \\\"+reason\"  New: \"err.Error()+\\\"; restore: \\\"+reason[:0]\"\nbit, 1 failed, restored yes. TestProbeMutationRecoveryFailure :108 only.\n\nP9 D41 exit 2. probe.go\nOld: \"return line + block, 2\"  New: \"return line + block, 1\"\nbit, 1 failed, restored yes. TestProbeMutationRecoveryFailure :102 only.\n\nP10 D41 exact copy bytes (shared with PB12, PB14). subject.go\nOld: \"if err := os.WriteFile(file, subject.start, 0o600); err != nil {\"\nNew: \"if err := os.WriteFile(file, subject.start[1:], 0o600); err != nil {\"\nbit, 29 failed, restored yes. Includes TestProbeMutationRecoveryFailure copy bytes and TestProbeReportsARestoreFailure.\n\nP11 D42 home: a proven restore keeps the copy. probe.go\nOld: \"if restored, reason := restore(replacer, subject, preserved); !restored {\"\nNew: \"if restored, reason := restore(replacer, subject, preserved); restored { return line, 1 } else {\"\nbit, 3 failed, restored yes. TestProbeMutationRecoverySuccess/directory-open, -sync, -close :84 only.\n\nP12 D43 old bytes: fallback write of the new bytes after an unpublished error. probe.go\nOld: \"return mutationFailure(replacer, subject, preserved, err)\"\nNew: \"_ = os.WriteFile(subject.path, mutated, subject.mode); return mutationFailure(replacer, subject, preserved, err)\"\nbit, 2 failed, restored yes. TestProbeMutationPreparationFailure :135 old-bytes check; existing PB45 :284.\n\nP13 D43 home: the rename stage is classified as published. probe.go\nOld: \"errors.As(err, &replaceErr) && replaceErr.Published {\"\nNew: \"errors.As(err, &replaceErr) && (replaceErr.Published || replaceErr.Stage == \\\"rename\\\") {\"\nbit, 1 failed, restored yes. TestProbeMutationPreparationFailure :136 Bench home holds 1 entries only.\n\nP14 native wrapper error path (run supplies durablefile.Native(); PB45 drives it through Command). probe.go\nOld: \"\\treturn runWith(durablefile.Native(), root, parsed)\"\nNew: \"\\tout, code := runWith(durablefile.Native(), root, parsed)\\n\\tif len(out) > 28 && out[:28] == \\\"error: probe mutation failed\\\" {\\n\\t\\treturn \\\"\\\", 0\\n\\t}\\n\\treturn out, code\"\nbit, 1 failed, restored yes. TestProbeRefusesWhenTheMutationWriteFails refusal_test.go:282 only.\n"
            }
          }
        },
        {
          "id": "dcm-t11-rep1-caller-tests",
          "performer": "claude:dcm_t11_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "808f0b78dfb2cc543b4bcadac7683c7900b809d4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "c377e632d4056ebb6853ed1087efbb38e010c824",
            "digest": "sha256:3aa916dce6aac7c3056bda1672082390f2096fe992d0c2d9c7f0dfdaa824998c",
            "excerpt": "bench test --package ./internal/probe at c377e632d4056ebb6853ed1087efbb38e010c824 (FT354-durable-callers; only reviews/durable-caller-migration.md dirty, no Go source differs from the commit): internal/probe pass, 117 tests run, failures[0], skips[0]. Repair witnesses in internal/probe/replacement_test.go: TestProbeNamesACopyTheEncoderRefuses (D41, both table-refusal branches), TestProbeMutationRecoveryFailure (D41, block scoped after the cause line), TestProbeMutationRecoverySuccess (D42, now asserts exit 1 and the baseline-only run per stage).\n"
          },
          "requirement": "caller-tests",
          "command": "bench test --package ./internal/probe",
          "exit_code": 0
        },
        {
          "id": "dcm-t11-rep1-command-contract",
          "performer": "claude:dcm_t11_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "808f0b78dfb2cc543b4bcadac7683c7900b809d4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "c377e632d4056ebb6853ed1087efbb38e010c824",
            "digest": "sha256:ea7159bb724169f7ef0ad0d46864cfaa502dfbc0111a018acbd53f86d2689e78",
            "excerpt": "bench test --package ./cmd/bench --run CommandRegistry on the repair source (Go bytes identical to c377e632d4056ebb6853ed1087efbb38e010c824): cmd/bench pass, 5 tests run, failures[0], skips[0]. Command membership, routes, and help are unchanged.\n"
          },
          "requirement": "command-contract",
          "command": "bench test --package ./cmd/bench --run 'CommandRegistry'",
          "exit_code": 0
        },
        {
          "id": "dcm-t11-rep1-axi-registry",
          "performer": "claude:dcm_t11_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "808f0b78dfb2cc543b4bcadac7683c7900b809d4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "c377e632d4056ebb6853ed1087efbb38e010c824",
            "digest": "sha256:9466abd804480910dbf3c7a64e995a5c81923e6dfc4e238f64a6751b474625dd",
            "excerpt": "bench test --check axi-query-registry on the repair source (Go bytes identical to c377e632d4056ebb6853ed1087efbb38e010c824): internal/conformance pass, 1 test run, failures[0], skips[0].\n"
          },
          "requirement": "axi-registry",
          "command": "bench test --check axi-query-registry",
          "exit_code": 0
        },
        {
          "id": "dcm-t11-rep1-routing",
          "performer": "claude:dcm_t11_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "808f0b78dfb2cc543b4bcadac7683c7900b809d4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "c377e632d4056ebb6853ed1087efbb38e010c824",
            "digest": "sha256:d3c1fbb515c43493e00c2e9a58803052956870e184978581f5d005e2b89599ff",
            "excerpt": "bench test --check subcommand-routing on the repair source (Go bytes identical to c377e632d4056ebb6853ed1087efbb38e010c824): internal/conformance pass, 1 test run, failures[0], skips[0].\n"
          },
          "requirement": "routing",
          "command": "bench test --check subcommand-routing",
          "exit_code": 0
        },
        {
          "id": "dcm-t11-rep1-cmd-bench",
          "performer": "claude:dcm_t11_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "808f0b78dfb2cc543b4bcadac7683c7900b809d4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "c377e632d4056ebb6853ed1087efbb38e010c824",
            "digest": "sha256:1d256b501a6ccf55fe12d08ef59eb72200268156f76a6b715d4e7b53fe77e969",
            "excerpt": "bench test --package ./cmd/bench on the repair source (Go bytes identical to c377e632d4056ebb6853ed1087efbb38e010c824): cmd/bench pass, 335 tests run, failures[0], skips[0].\n"
          },
          "requirement": "cmd-bench",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "dcm-t11-rep1-conformance",
          "performer": "claude:dcm_t11_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "808f0b78dfb2cc543b4bcadac7683c7900b809d4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "c377e632d4056ebb6853ed1087efbb38e010c824",
            "digest": "sha256:e7e40f10e37b5ea0f929f50969a1cc1a3a5a5b5adcee88b0b4bfd1507fb289e5",
            "excerpt": "bench test --package ./internal/conformance on the repair source (Go bytes identical to c377e632d4056ebb6853ed1087efbb38e010c824): internal/conformance pass, 1690 tests run, failures[0], skips[1]: TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device (capability: privilege: cannot create a character device), an environment capability skip not owned by this ticket.\n"
          },
          "requirement": "conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "dcm-t11-rep1-probetest",
          "performer": "claude:dcm_t11_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "808f0b78dfb2cc543b4bcadac7683c7900b809d4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "c377e632d4056ebb6853ed1087efbb38e010c824",
            "digest": "sha256:fb90135624e8a3bc3622da431e7bcbd68b9c4c32435c0b01e5da3b4888d90013",
            "excerpt": "bench test --package ./internal/probe/probetest at c377e632d4056ebb6853ed1087efbb38e010c824 (FT354-durable-callers; only reviews/durable-caller-migration.md dirty, no Go source differs from the commit): internal/probe/probetest no-tests, 0 tests run, failures[0], skips[0]. The fixture is unchanged by the repair.\n"
          },
          "requirement": "probetest",
          "command": "bench test --package ./internal/probe/probetest",
          "exit_code": 0
        },
        {
          "id": "dcm-t11-rep1-root-conformance",
          "performer": "claude:dcm_t11_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "808f0b78dfb2cc543b4bcadac7683c7900b809d4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "c377e632d4056ebb6853ed1087efbb38e010c824",
            "digest": "sha256:a150a649252c79ba3775a7a253070215a83d460abffa1f5817cb9a3946180b8c",
            "excerpt": "bench test --check root-conformance on the repair source (Go bytes identical to c377e632d4056ebb6853ed1087efbb38e010c824): internal/conformance pass, 1 test run, failures[0], skips[0].\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "dcm-t11-rep1-omission-proof",
          "performer": "claude:dcm_t11_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "808f0b78dfb2cc543b4bcadac7683c7900b809d4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "c377e632d4056ebb6853ed1087efbb38e010c824",
            "digest": "sha256:a8f25404c54e9e008d4aa89f38e6e3725da3681492a7a97b82b478a6e71cbed0",
            "excerpt": "D-C4 repair 1 omission-proof probes at c377e632d4056ebb6853ed1087efbb38e010c824 (FT354-durable-callers), each through bench probe <file> --swap <old> --with <new> --package ./internal/probe with no --run filter. Every verdict bit, every restored yes; the package passes 117 tests after each restore.\nNamed probe: release preservation on every mutation error, or treat failed restore sync as success; the preservation-path and focused-start witnesses detect each mutation.\n\nReplayed at the repair tip (O1 swap text adapted, because the repair changed its line):\nO1 (release the copy on the failed-restore exit). internal/probe/probe.go\nOld: \"return line + keptCopy(\"\nNew: \"preserved.release(); return line + keptCopy(\"\nbit, 2 failed, restored yes. TestProbeMutationRecoveryFailure replacement_test.go:111 and TestProbeNamesACopyTheEncoderRefuses/mutation_failure :152: preserved files = ([], <nil>), want exactly one.\n\nO2 (release without restore on every mutation error). internal/probe/probe.go\nOld: \"if errors.As(err, &replaceErr) && replaceErr.Published {\"\nNew: \"if false && errors.As(err, &replaceErr) && replaceErr.Published {\"\nbit, 5 failed, restored yes. TestProbeMutationRecoverySuccess/directory-open, -sync, -close :90 subject holds the mutated bytes; TestProbeNamesACopyTheEncoderRefuses/mutation_failure :150 exit = 1, want 2; TestProbeMutationRecoveryFailure.\n\nO3 (treat failed restore sync as success). internal/probe/subject.go\nOld: \"if err := replacer.Replace(subject.path, saved, subject.mode); err != nil {\"\nNew: \"if err := replacer.Replace(subject.path, saved, subject.mode); err != nil && !strings.Contains(err.Error(), \\\"rename completed\\\") {\"\nbit, 1 failed, restored yes. TestProbeFailedRestore :177 row bit,...,yes in place of restore-failed.\n\nO4 (owner bypass, mutation writer; D27). internal/probe/probe.go\nOld: \"if err := replacer.Replace(subject.path, mutated, subject.mode); err != nil {\"\nNew: \"if err := durablefile.Replace(subject.path, mutated, subject.mode); err != nil {\"\nbit, 9 failed, restored yes. TestProbeReplacement :61; TestProbeMutationRecoverySuccess rows :87 exit = 2, want 1; both TestProbeNamesACopyTheEncoderRefuses rows.\n\nO5 (owner bypass, restore writer). internal/probe/subject.go\nOld: \"if err := replacer.Replace(subject.path, saved, subject.mode); err != nil {\"\nNew: \"if err := durablefile.Replace(subject.path, saved, subject.mode); err != nil {\"\nbit, 4 failed, restored yes. TestProbeFailedRestore :177; TestProbeMutationRecoveryFailure :109 exit = 1, want 2; both TestProbeNamesACopyTheEncoderRefuses rows :150.\n\nCited, not replayed (the repair changes neither their production line nor their failing check): O6 and O7 in /tmp/claude-1000/-home-mgibs-workspace-bench/a142325d-036b-4e84-b387-f6b875be269c/scratchpad/dcm_t11/rec/omission-proof.txt at 0e86fec98de40c9c5ff0e21cd719255731934188.\n\nNew repair reds (P1/C1: a table refusal must not hide the copy path or a cause):\nN1 (fallback omits the copy path). internal/probe/probe.go\nOld: \"toon.Errorf(\\\"probe copy kept\\\", strconv.Quote(preserved.file)+\\\": \\\"+reason)\"\nNew: \"toon.Errorf(\\\"probe copy kept\\\", reason)\"\nbit, 2 failed, restored yes. TestProbeNamesACopyTheEncoderRefuses/mutation_failure and /focused_run :161 want it to name the copy. Whole package, no --run: only this test holds the quoted-path literal.\n\nN2 (fallback answers nothing on a table refusal, the pre-repair shape). internal/probe/probe.go\nOld: \"\\tif err == nil {\\n\\t\\treturn block\\n\\t}\"\nNew: \"\\tif true {\\n\\t\\treturn block\\n\\t}\"\nbit, 2 failed, restored yes. TestProbeNamesACopyTheEncoderRefuses/mutation_failure and /focused_run :156 want the restore cause.\n\nN3 (fallback drops the reason). internal/probe/probe.go\nOld: \"strconv.Quote(preserved.file)+\\\": \\\"+reason)\"\nNew: \"strconv.Quote(preserved.file))\"\nbit, 2 failed, restored yes. Both TestProbeNamesACopyTheEncoderRefuses rows :156 want \"injected restore rename failure\".\n\nN4 (mutationFailure branch drops the kept block). internal/probe/probe.go\nOld: \"return line + keptCopy(preserved, err.Error()+\\\"; restore: \\\"+reason), 2\"\nNew: \"return line + keptCopy(preserved, err.Error()+\\\"; restore: \\\"+reason)[:0], 2\"\nbit, 2 failed, restored yes. TestProbeMutationRecoveryFailure :116; TestProbeNamesACopyTheEncoderRefuses/mutation_failure :156.\n\nN5 (render branch drops the kept block behind an envelope refusal). internal/probe/probe.go\nOld: \"return toon.RenderError(err) + \\\"\\\\n\\\" + preservedBlock, code\"\nNew: \"return toon.RenderError(err) + \\\"\\\\n\\\", code\"\nbit, 2 failed, restored yes. TestProbeNamesACopyTheEncoderRefuses/focused_run :156; existing PB47 outcome_test.go:108.\n"
          },
          "requirement": "omission-proof",
          "command": "bench test --package ./internal/probe",
          "exit_code": 0,
          "probe": {
            "mutation": "Release preservation on every mutation error or treat failed restore sync as success. The preservation-path and focused-start witnesses must detect each mutation.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "c377e632d4056ebb6853ed1087efbb38e010c824",
              "digest": "sha256:a8f25404c54e9e008d4aa89f38e6e3725da3681492a7a97b82b478a6e71cbed0",
              "excerpt": "D-C4 repair 1 omission-proof probes at c377e632d4056ebb6853ed1087efbb38e010c824 (FT354-durable-callers), each through bench probe <file> --swap <old> --with <new> --package ./internal/probe with no --run filter. Every verdict bit, every restored yes; the package passes 117 tests after each restore.\nNamed probe: release preservation on every mutation error, or treat failed restore sync as success; the preservation-path and focused-start witnesses detect each mutation.\n\nReplayed at the repair tip (O1 swap text adapted, because the repair changed its line):\nO1 (release the copy on the failed-restore exit). internal/probe/probe.go\nOld: \"return line + keptCopy(\"\nNew: \"preserved.release(); return line + keptCopy(\"\nbit, 2 failed, restored yes. TestProbeMutationRecoveryFailure replacement_test.go:111 and TestProbeNamesACopyTheEncoderRefuses/mutation_failure :152: preserved files = ([], <nil>), want exactly one.\n\nO2 (release without restore on every mutation error). internal/probe/probe.go\nOld: \"if errors.As(err, &replaceErr) && replaceErr.Published {\"\nNew: \"if false && errors.As(err, &replaceErr) && replaceErr.Published {\"\nbit, 5 failed, restored yes. TestProbeMutationRecoverySuccess/directory-open, -sync, -close :90 subject holds the mutated bytes; TestProbeNamesACopyTheEncoderRefuses/mutation_failure :150 exit = 1, want 2; TestProbeMutationRecoveryFailure.\n\nO3 (treat failed restore sync as success). internal/probe/subject.go\nOld: \"if err := replacer.Replace(subject.path, saved, subject.mode); err != nil {\"\nNew: \"if err := replacer.Replace(subject.path, saved, subject.mode); err != nil && !strings.Contains(err.Error(), \\\"rename completed\\\") {\"\nbit, 1 failed, restored yes. TestProbeFailedRestore :177 row bit,...,yes in place of restore-failed.\n\nO4 (owner bypass, mutation writer; D27). internal/probe/probe.go\nOld: \"if err := replacer.Replace(subject.path, mutated, subject.mode); err != nil {\"\nNew: \"if err := durablefile.Replace(subject.path, mutated, subject.mode); err != nil {\"\nbit, 9 failed, restored yes. TestProbeReplacement :61; TestProbeMutationRecoverySuccess rows :87 exit = 2, want 1; both TestProbeNamesACopyTheEncoderRefuses rows.\n\nO5 (owner bypass, restore writer). internal/probe/subject.go\nOld: \"if err := replacer.Replace(subject.path, saved, subject.mode); err != nil {\"\nNew: \"if err := durablefile.Replace(subject.path, saved, subject.mode); err != nil {\"\nbit, 4 failed, restored yes. TestProbeFailedRestore :177; TestProbeMutationRecoveryFailure :109 exit = 1, want 2; both TestProbeNamesACopyTheEncoderRefuses rows :150.\n\nCited, not replayed (the repair changes neither their production line nor their failing check): O6 and O7 in /tmp/claude-1000/-home-mgibs-workspace-bench/a142325d-036b-4e84-b387-f6b875be269c/scratchpad/dcm_t11/rec/omission-proof.txt at 0e86fec98de40c9c5ff0e21cd719255731934188.\n\nNew repair reds (P1/C1: a table refusal must not hide the copy path or a cause):\nN1 (fallback omits the copy path). internal/probe/probe.go\nOld: \"toon.Errorf(\\\"probe copy kept\\\", strconv.Quote(preserved.file)+\\\": \\\"+reason)\"\nNew: \"toon.Errorf(\\\"probe copy kept\\\", reason)\"\nbit, 2 failed, restored yes. TestProbeNamesACopyTheEncoderRefuses/mutation_failure and /focused_run :161 want it to name the copy. Whole package, no --run: only this test holds the quoted-path literal.\n\nN2 (fallback answers nothing on a table refusal, the pre-repair shape). internal/probe/probe.go\nOld: \"\\tif err == nil {\\n\\t\\treturn block\\n\\t}\"\nNew: \"\\tif true {\\n\\t\\treturn block\\n\\t}\"\nbit, 2 failed, restored yes. TestProbeNamesACopyTheEncoderRefuses/mutation_failure and /focused_run :156 want the restore cause.\n\nN3 (fallback drops the reason). internal/probe/probe.go\nOld: \"strconv.Quote(preserved.file)+\\\": \\\"+reason)\"\nNew: \"strconv.Quote(preserved.file))\"\nbit, 2 failed, restored yes. Both TestProbeNamesACopyTheEncoderRefuses rows :156 want \"injected restore rename failure\".\n\nN4 (mutationFailure branch drops the kept block). internal/probe/probe.go\nOld: \"return line + keptCopy(preserved, err.Error()+\\\"; restore: \\\"+reason), 2\"\nNew: \"return line + keptCopy(preserved, err.Error()+\\\"; restore: \\\"+reason)[:0], 2\"\nbit, 2 failed, restored yes. TestProbeMutationRecoveryFailure :116; TestProbeNamesACopyTheEncoderRefuses/mutation_failure :156.\n\nN5 (render branch drops the kept block behind an envelope refusal). internal/probe/probe.go\nOld: \"return toon.RenderError(err) + \\\"\\\\n\\\" + preservedBlock, code\"\nNew: \"return toon.RenderError(err) + \\\"\\\\n\\\", code\"\nbit, 2 failed, restored yes. TestProbeNamesACopyTheEncoderRefuses/focused_run :156; existing PB47 outcome_test.go:108.\n"
            }
          }
        },
        {
          "id": "dcm-t11-rep1-preservation",
          "performer": "claude:dcm_t11_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "808f0b78dfb2cc543b4bcadac7683c7900b809d4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "c377e632d4056ebb6853ed1087efbb38e010c824",
            "digest": "sha256:dd30b199e3591fa81531ca56057ad18197463514a5c6a7f93ed80ba6dbe4f5b5",
            "excerpt": "D-C4 repair 1 preservation probes at c377e632d4056ebb6853ed1087efbb38e010c824 (FT354-durable-callers), each through bench probe internal/probe/probe.go --swap <old> --with <new> --package ./internal/probe. Every verdict bit, every restored yes; the package passes 117 tests after each restore. The omission-proof entry of this repair holds O1-O5 and the new reds N1-N5; N1 is the whole-package necessity red of the one kept literal, the Go-quoted copy path.\n\nReplayed at the repair tip, because the repair moved their production line into the keptCopy call or changed the failing check (swap text adapted):\nP6 D41 preserved row names the copy path. Old: \"[][]string{{preserved.file, reason}}\" New: \"[][]string{{preserved.dir, reason}}\" (no --run). bit, 5 failed, restored yes. TestProbeMutationRecoveryFailure; both TestProbeNamesACopyTheEncoderRefuses rows :161; PB14 outcome_test.go:53; PB47 :108.\nP7 D41 reason carries the mutation cause. Old: \"keptCopy(preserved, err.Error()+\\\"; restore: \\\"+reason)\" New: \"keptCopy(preserved, \\\"restore: \\\"+reason)\" (--run '^TestProbeMutationRecoveryFailure$'). bit, 1 failed, restored yes. replacement_test.go:116 preserved block lacks the mutation cause.\nP8 D41 reason carries the restore cause. Old: \"keptCopy(preserved, err.Error()+\\\"; restore: \\\"+reason)\" New: \"keptCopy(preserved, err.Error()+\\\"; restore: \\\"+reason[:0])\" (--run '^TestProbeMutationRecoveryFailure$'). bit, 1 failed, restored yes. replacement_test.go:116.\nP9 D41 exit 2. Old: \"\\\"; restore: \\\"+reason), 2\" New: \"\\\"; restore: \\\"+reason), 1\" (--run '^TestProbeMutationRecoveryFailure$'). bit, 1 failed, restored yes. replacement_test.go:109 exit = 1, want 2.\n\nNew repair red for the D42 exit check the repair adds to TestProbeMutationRecoverySuccess:\nN6 Old: \"\\tpreserved.release()\\n\\treturn line, 1\" New: \"\\tpreserved.release()\\n\\treturn line, 3\" (no --run). bit, 5 failed, restored yes. TestProbeMutationRecoverySuccess/directory-open, -sync, -close :87 exit = 3, want 1; TestProbeReplacement :70; PB45 refusal_test.go:282. The added baseline-only check is the O7 witness, cited below.\n\nCited, not replayed (the repair changes neither their production line nor their failing check): P1-P5 and P10-P14 in /tmp/claude-1000/-home-mgibs-workspace-bench/a142325d-036b-4e84-b387-f6b875be269c/scratchpad/dcm_t11/rec/preservation.txt, and O6-O7 in the omission-proof.txt beside it, all at 0e86fec98de40c9c5ff0e21cd719255731934188.\n"
          },
          "requirement": "preservation",
          "command": "bench test --package ./internal/probe",
          "exit_code": 0,
          "probe": {
            "mutation": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "c377e632d4056ebb6853ed1087efbb38e010c824",
              "digest": "sha256:dd30b199e3591fa81531ca56057ad18197463514a5c6a7f93ed80ba6dbe4f5b5",
              "excerpt": "D-C4 repair 1 preservation probes at c377e632d4056ebb6853ed1087efbb38e010c824 (FT354-durable-callers), each through bench probe internal/probe/probe.go --swap <old> --with <new> --package ./internal/probe. Every verdict bit, every restored yes; the package passes 117 tests after each restore. The omission-proof entry of this repair holds O1-O5 and the new reds N1-N5; N1 is the whole-package necessity red of the one kept literal, the Go-quoted copy path.\n\nReplayed at the repair tip, because the repair moved their production line into the keptCopy call or changed the failing check (swap text adapted):\nP6 D41 preserved row names the copy path. Old: \"[][]string{{preserved.file, reason}}\" New: \"[][]string{{preserved.dir, reason}}\" (no --run). bit, 5 failed, restored yes. TestProbeMutationRecoveryFailure; both TestProbeNamesACopyTheEncoderRefuses rows :161; PB14 outcome_test.go:53; PB47 :108.\nP7 D41 reason carries the mutation cause. Old: \"keptCopy(preserved, err.Error()+\\\"; restore: \\\"+reason)\" New: \"keptCopy(preserved, \\\"restore: \\\"+reason)\" (--run '^TestProbeMutationRecoveryFailure$'). bit, 1 failed, restored yes. replacement_test.go:116 preserved block lacks the mutation cause.\nP8 D41 reason carries the restore cause. Old: \"keptCopy(preserved, err.Error()+\\\"; restore: \\\"+reason)\" New: \"keptCopy(preserved, err.Error()+\\\"; restore: \\\"+reason[:0])\" (--run '^TestProbeMutationRecoveryFailure$'). bit, 1 failed, restored yes. replacement_test.go:116.\nP9 D41 exit 2. Old: \"\\\"; restore: \\\"+reason), 2\" New: \"\\\"; restore: \\\"+reason), 1\" (--run '^TestProbeMutationRecoveryFailure$'). bit, 1 failed, restored yes. replacement_test.go:109 exit = 1, want 2.\n\nNew repair red for the D42 exit check the repair adds to TestProbeMutationRecoverySuccess:\nN6 Old: \"\\tpreserved.release()\\n\\treturn line, 1\" New: \"\\tpreserved.release()\\n\\treturn line, 3\" (no --run). bit, 5 failed, restored yes. TestProbeMutationRecoverySuccess/directory-open, -sync, -close :87 exit = 3, want 1; TestProbeReplacement :70; PB45 refusal_test.go:282. The added baseline-only check is the O7 witness, cited below.\n\nCited, not replayed (the repair changes neither their production line nor their failing check): P1-P5 and P10-P14 in /tmp/claude-1000/-home-mgibs-workspace-bench/a142325d-036b-4e84-b387-f6b875be269c/scratchpad/dcm_t11/rec/preservation.txt, and O6-O7 in the omission-proof.txt beside it, all at 0e86fec98de40c9c5ff0e21cd719255731934188.\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "dcm-c4-r1-standards",
          "performer": "claude:dcm_rev_c4_std",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "e6ce538abb8886a1123ff7ec066fac3b7e74dbef",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_c4_std",
            "digest": "sha256:81dded4bed91eb462d17e1803794c93a8fd43a218bfd6fdbfc84ef15f523d614",
            "excerpt": "## Standards (D-C4, 68e390e1..0e86fec9) — 3 findings; worst S1.\nS1 | auto-fix | 4 | internal/probe/replacement_test.go:105 | strings.Cut on a hand-written \"preserved[1]{path,reason}:\" header is a third copy of the preservedTable encoder fact | scope the row without restating the header (derive through preservedTable/toon)\nS2 | auto-fix | 4 | internal/probe/probe.go:217-219 | the shared preservedTable comment speaks of \"a subject the row already refused\", true only for render | state the why without a caller-specific subject\nS3 | auto-fix (orchestrator: same as C2; charge-common corrected) | 5 | reviews/durable-caller-migration.md:3697-3811 | dcm-t11 probe entries disagree with their -r2 twins | record corrected entries at the repair source\nAdvice: shared mutated-source constant; one 0o640 constant; requireSubjectBytes vs RequireBytes is outside 11's Writes (D-D).\n"
          },
          "axis": "Standards",
          "base": "68e390e155699b5aae10125808e86f8eb5eb3b93",
          "tip": "0e86fec98de40c9c5ff0e21cd719255731934188",
          "finding_ids": [
            "S1",
            "S2",
            "S3"
          ],
          "supersedes": []
        },
        {
          "id": "dcm-c4-r1-spec",
          "performer": "claude:dcm_rev_c4_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "e6ce538abb8886a1123ff7ec066fac3b7e74dbef",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_c4_spec",
            "digest": "sha256:36de25f98bd50d9e12aec74712d4bee88fe1594125d2e0ed93794df23e3ba08e",
            "excerpt": "## Spec (D-C4, 68e390e1..0e86fec9) — 1 finding; worst P1.\nD27, D34, D41, D42, D43 met; coverage cells cite existing tests; no behavior change outside the acceptance.\nP1 | auto-fix (orchestrator, --auto-approve; same branch as C1) | 5 | internal/probe/probe.go:206-209 | when preservedTable refuses (control byte in the subject base name) the exit-2 output drops the mutation cause and the copy path, against spec.md:113-118; render (probe.go:246-249) has the same gap | on a table refusal still print the cause line and name the copy path outside TOON, in both branches\nAdvice: TestProbeMutationRecoverySuccess discards out and code; assert exit 1 and baseline-only per stage.\n"
          },
          "axis": "Spec",
          "base": "68e390e155699b5aae10125808e86f8eb5eb3b93",
          "tip": "0e86fec98de40c9c5ff0e21cd719255731934188",
          "finding_ids": [
            "P1"
          ],
          "supersedes": []
        },
        {
          "id": "dcm-c4-r1-coverage",
          "performer": "claude:dcm_rev_c4_cov",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "e6ce538abb8886a1123ff7ec066fac3b7e74dbef",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_c4_cov",
            "digest": "sha256:7cb51d1094c805ad17e12fffc99b2b9018f3bf8cee7f9b860b35eed830ceeaff",
            "excerpt": "## Coverage (D-C4, 68e390e1..0e86fec9) — 3 findings; worst C2. Replays O6 and P13 bit, restored yes.\nC2 | auto-fix | 8 | reviews/durable-caller-migration.md:3774,3803; internal/reviewrecord/check.go:59 | the graded -r2 probe entries record probe exit 0, which the checkpoint refuses | record entries with --exit-code 0 and a nonzero --probe-exit-code\nC1 | auto-fix | 7 | internal/probe/probe.go:207-209 | the table-refusal branch of mutationFailure has no red (swap to return \"\", 0 stayed green) | add a BEL base-name witness: published mutation fault plus restore fault, exit 2, copy kept, path named\nC3 | no-op (D30 routed to ticket 14 by a logged learning) | 6 | internal/probe/replacement_test.go | no D30 refusal witness through the replacer-taking form; native refusal tests catch the swallow | none in this chunk\nAdvice: the leading \"probe mutation failed\" line on exit 2 is not pinned.\n"
          },
          "axis": "Coverage",
          "base": "68e390e155699b5aae10125808e86f8eb5eb3b93",
          "tip": "0e86fec98de40c9c5ff0e21cd719255731934188",
          "finding_ids": [
            "C1",
            "C2",
            "C3"
          ],
          "supersedes": []
        },
        {
          "id": "dcm-c4-r2-standards",
          "performer": "claude:dcm_rev_c4_std2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "808f0b78dfb2cc543b4bcadac7683c7900b809d4",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_c4_std2",
            "digest": "sha256:a8686cedde438890b5bded068574b765532355907d45cce1df4aec316d881e16",
            "excerpt": "## Standards (D-C4 confirming, f3e417f8..c377e632) — S1, S2, S3 closed; 1 new finding.\nS4 | auto-fix (orchestrator: evidence-only comment correction; Standards reaffirms) | 5 | internal/probe/probe.go:231-232 | the render doc comment says a failed restore keeps \"the preserved row\", but keptCopy answers with two fallback error lines when the copy path refuses | say that a failed restore keeps the answer that names the kept copy\nVeto items graded with no defect: keptCopy rename, the fallback shape, the variadic probeWith argument.\n"
          },
          "axis": "Standards",
          "base": "68e390e155699b5aae10125808e86f8eb5eb3b93",
          "tip": "c377e632d4056ebb6853ed1087efbb38e010c824",
          "finding_ids": [
            "S4"
          ],
          "supersedes": [
            "dcm-c4-r1-standards"
          ]
        },
        {
          "id": "dcm-c4-r2-spec",
          "performer": "claude:dcm_rev_c4_spec2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "808f0b78dfb2cc543b4bcadac7683c7900b809d4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_c4_spec2",
            "digest": "sha256:be690513cf0d98313120a5afa836217528f9fd157efc2a2008ff7980686eeaa6",
            "excerpt": "## Spec (D-C4 confirming, f3e417f8..c377e632) — P1 closed (probe.go:206, 218-225, 248, 257); 0 new findings.\nNo behavior change outside ticket 11's acceptance; encodable paths are byte-identical.\nAdvice: the mutation-cause line prints err.Error() unstripped; older behavior, not required by the spec.\n"
          },
          "axis": "Spec",
          "base": "68e390e155699b5aae10125808e86f8eb5eb3b93",
          "tip": "c377e632d4056ebb6853ed1087efbb38e010c824",
          "finding_ids": [],
          "supersedes": [
            "dcm-c4-r1-spec"
          ]
        },
        {
          "id": "dcm-c4-r2-coverage",
          "performer": "claude:dcm_rev_c4_cov2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "808f0b78dfb2cc543b4bcadac7683c7900b809d4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_c4_cov2",
            "digest": "sha256:cb3db9e7d0f276eb8e161d7958825e4b99de7016d08fff75c17be143dec88895",
            "excerpt": "## Coverage (D-C4 confirming, f3e417f8..c377e632) — C1 closed (replacement_test.go:124-166; reds N1-N5); C2 closed (dcm-t11-rep1 probe entries: exit 0, bit, probe exit 1, restore pass); 0 new findings.\nNo probe replayed; the orchestrator's independent sanitize.Strip probe bit, restored yes.\n"
          },
          "axis": "Coverage",
          "base": "68e390e155699b5aae10125808e86f8eb5eb3b93",
          "tip": "c377e632d4056ebb6853ed1087efbb38e010c824",
          "finding_ids": [],
          "supersedes": [
            "dcm-c4-r1-coverage"
          ]
        },
        {
          "id": "dcm-c4-r3-standards",
          "performer": "claude:dcm_rev_c4_std3",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "808f0b78dfb2cc543b4bcadac7683c7900b809d4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_c4_std3",
            "digest": "sha256:3bb3b674f9a882b0f1aa54e40a27450165b6296b1dc013780e85a1d8964dbf0f",
            "excerpt": "## Standards (D-C4 reaffirmation, 2453d334..5d56efc1) — S4 closed (internal/probe/probe.go:231-232); comment-only diff; 0 new findings.\n"
          },
          "axis": "Standards",
          "base": "68e390e155699b5aae10125808e86f8eb5eb3b93",
          "tip": "c377e632d4056ebb6853ed1087efbb38e010c824",
          "finding_ids": [],
          "supersedes": [
            "dcm-c4-r2-standards"
          ]
        },
        {
          "id": "dcm-c4-r3-spec",
          "performer": "claude:dcm_rev_c4_spec2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "808f0b78dfb2cc543b4bcadac7683c7900b809d4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_c4_spec2",
            "digest": "sha256:1c89af8c1b82416cfc9747c517826cdff9d3357fe0e84b6e4a7dcf17aa619199",
            "excerpt": "## Spec (D-C4 refrozen pair, 19cb2923..c377e632) — pass; P1 closed; the added delta 19cb2923..68e390e1 is a comment-only test edit (internal/brokermanifest/replacement_test.go:85-86) plus review records; 0 open findings.\n"
          },
          "axis": "Spec",
          "base": "19cb292348cf4aae0f895a139ec39cf8a6ba8edb",
          "tip": "c377e632d4056ebb6853ed1087efbb38e010c824",
          "finding_ids": [],
          "supersedes": [
            "dcm-c4-r2-spec"
          ]
        },
        {
          "id": "dcm-c4-r3-coverage",
          "performer": "claude:dcm_rev_c4_cov2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "808f0b78dfb2cc543b4bcadac7683c7900b809d4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_c4_cov2",
            "digest": "sha256:c449521f53ce499306c325ee44bcf78370ad0a180b011ae560a233fe79d91b52",
            "excerpt": "## Coverage (D-C4 refrozen pair, 19cb2923..c377e632) — pass; C1, C2 closed; C3 no-op; the added delta is a comment-only test edit plus review records and needs no witness; dcm-t11-rep1 entries still graded; 0 open findings.\n"
          },
          "axis": "Coverage",
          "base": "19cb292348cf4aae0f895a139ec39cf8a6ba8edb",
          "tip": "c377e632d4056ebb6853ed1087efbb38e010c824",
          "finding_ids": [],
          "supersedes": [
            "dcm-c4-r2-coverage"
          ]
        },
        {
          "id": "dcm-c4-r4-standards",
          "performer": "claude:dcm_rev_c4_std3",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "808f0b78dfb2cc543b4bcadac7683c7900b809d4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_c4_std3",
            "digest": "sha256:e242c9db73241d9cab42494bed553a19d7c502ee1e39c9ea66114ef87578a512",
            "excerpt": "## Standards (D-C4 refrozen pair, 19cb2923..c377e632) — pass; S1-S4 closed; the added delta 19cb2923..68e390e1 is an accurate comment-only test edit plus review records; the post-tip S4 comment fix 5d56efc1 is graded; 0 open findings.\n"
          },
          "axis": "Standards",
          "base": "19cb292348cf4aae0f895a139ec39cf8a6ba8edb",
          "tip": "c377e632d4056ebb6853ed1087efbb38e010c824",
          "finding_ids": [],
          "supersedes": [
            "dcm-c4-r3-standards"
          ]
        }
      ]
    },
    {
      "id": "D-C3",
      "base": "c377e632d4056ebb6853ed1087efbb38e010c824",
      "tip": "b20919b10799df0f4014df0a588bf70c3cea986c",
      "plan_digest": "sha256:caaaa230b74e8cf3a79cb22fa805c3a28bafb054e1b4053e68684ed4a0dee7cf",
      "source_digest": "2edd23a96795019c66d36851373250f4b7b5329f",
      "acceptance_rows": [
        "D26"
      ],
      "verification": [
        {
          "id": "dc3-caller-tests-1",
          "performer": "claude:dcm_t10",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "4d77c1ef8c235ca9bc90f4a80f95c22bea5f2751",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_t10",
            "digest": "sha256:00416b58a9af6b99c8e4ca2504b3ca1e4ed7ac63e75da55b9e6854a5284329b1",
            "excerpt": "D-C3 caller-tests at ee901c01b9c522f654a8524b675f1e94c4e148c2 (FT354-durable-callers): bench test --package ./internal/dashboard\npackages: github.com/gibbonmi/bench/internal/dashboard,pass,tests_run=23; failures 0; skips 0.\n"
          },
          "requirement": "caller-tests",
          "command": "bench test --package ./internal/dashboard",
          "exit_code": 0
        },
        {
          "id": "dc3-existing-canary-owner-1",
          "performer": "claude:dcm_t10",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "4d77c1ef8c235ca9bc90f4a80f95c22bea5f2751",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_t10",
            "digest": "sha256:ef5c2c42aca225bb44d097836b2a46d83073e8002c256ada6133bee7d504c6b3",
            "excerpt": "D-C3 existing-canary-owner at ee901c01b9c522f654a8524b675f1e94c4e148c2 (FT354-durable-callers): bench test --check git-plumbing-owner\ncheck git-plumbing-owner, conformance, tests_run=1; internal/conformance pass; failures 0; skips 0.\nThe git-flag-retyped fixture is unchanged: BASE names internal/dashboard/dashboard.go, and its MUTATE.json old line \"\\tgitDir, err := git.AdminDir(root)\\n\" stays intact in dashboard.go. bench test --check package-core-guard also passed at this tip (tests_run=1, failures 0).\n"
          },
          "requirement": "existing-canary-owner",
          "command": "bench test --check git-plumbing-owner",
          "exit_code": 0
        },
        {
          "id": "dc3-root-conformance-1",
          "performer": "claude:dcm_t10",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "4d77c1ef8c235ca9bc90f4a80f95c22bea5f2751",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_t10",
            "digest": "sha256:5a3db8ef3168393ea75073d56e76bcc41a470852e3e357f7542a21e81b56101c",
            "excerpt": "D-C3 root-conformance at ee901c01b9c522f654a8524b675f1e94c4e148c2 (FT354-durable-callers): bench test --check root-conformance\ncheck root-conformance, conformance, tests_run=1; internal/conformance pass; failures 0; skips 0.\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "dc3-omission-proof-1",
          "performer": "claude:dcm_t10",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "4d77c1ef8c235ca9bc90f4a80f95c22bea5f2751",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_t10",
            "digest": "sha256:c0d02c8c9fedef957ed27d6322ad524b36ac73dff0e6b6b8c12582e558f3ad69",
            "excerpt": "D-C3 omission-proof probes at ee901c01b9c522f654a8524b675f1e94c4e148c2 (FT354-durable-callers), subject internal/dashboard/dashboard.go, bench test --package ./internal/dashboard. Plan probe: \"Select 0644 by assumption or restore the former writer. The mode or injected-failure witness must turn red.\"\n\nProbe 1 (restore the former writer; owner bypass; omission probe with --run '^TestDashboardReplacement$').\nSwap old (exact): \"\\treturn r.Replace(target, []byte(content), pageMode)\"\nSwap new (exact): \"\\tnative := durablefile.Native()\\n\\ttmp, err := native.CreateTemp(filepath.Dir(target), \\\"bench-dashboard-*.html.tmp\\\")\\n\\tif err != nil {\\n\\t\\treturn err\\n\\t}\\n\\tname := tmp.Name()\\n\\tif _, err := tmp.Write([]byte(content)); err != nil {\\n\\t\\ttmp.Close()\\n\\t\\tnative.Remove(name)\\n\\t\\treturn err\\n\\t}\\n\\tif err := tmp.Close(); err != nil {\\n\\t\\tnative.Remove(name)\\n\\t\\treturn err\\n\\t}\\n\\tif err := native.Rename(name, target); err != nil {\\n\\t\\tnative.Remove(name)\\n\\t\\treturn err\\n\\t}\\n\\treturn nil\"\n(The former CreateTemp/Write/Close/Rename sequence with no sync, through the native operations, ignoring the injected replacer.)\nResult: verdict bit, cause failed, failed_tests 1, restored yes. Mutated package run exit 1.\nFailure: TestDashboardReplacement replacement_test.go:33: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected dashboard directory-sync failure.\n\nProbe 2 (select 0644 by assumption; kept mode literal 0o600, whole package, no --run filter).\nSwap old (exact): \"pageMode fs.FileMode = 0o600\"\nSwap new (exact): \"pageMode fs.FileMode = 0o644\"\nResult: verdict bit, cause failed, failed_tests 1 of 23, restored yes. Mutated package run exit 1.\nFailure: TestDashboardReplacement replacement_test.go:34: .../bench-dashboard.html mode = -rw-r--r--, want a regular file at -rw-------.\n\nRestoration: bench probe restored the source (restored=yes); bench test --package ./internal/dashboard passes at ee901c01 (23 tests).\n"
          },
          "requirement": "omission-proof",
          "command": "bench test --package ./internal/dashboard",
          "exit_code": 0,
          "probe": {
            "mutation": "Select 0644 by assumption or restore the former writer. The mode or injected-failure witness must turn red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:dcm_t10",
              "digest": "sha256:c0d02c8c9fedef957ed27d6322ad524b36ac73dff0e6b6b8c12582e558f3ad69",
              "excerpt": "D-C3 omission-proof probes at ee901c01b9c522f654a8524b675f1e94c4e148c2 (FT354-durable-callers), subject internal/dashboard/dashboard.go, bench test --package ./internal/dashboard. Plan probe: \"Select 0644 by assumption or restore the former writer. The mode or injected-failure witness must turn red.\"\n\nProbe 1 (restore the former writer; owner bypass; omission probe with --run '^TestDashboardReplacement$').\nSwap old (exact): \"\\treturn r.Replace(target, []byte(content), pageMode)\"\nSwap new (exact): \"\\tnative := durablefile.Native()\\n\\ttmp, err := native.CreateTemp(filepath.Dir(target), \\\"bench-dashboard-*.html.tmp\\\")\\n\\tif err != nil {\\n\\t\\treturn err\\n\\t}\\n\\tname := tmp.Name()\\n\\tif _, err := tmp.Write([]byte(content)); err != nil {\\n\\t\\ttmp.Close()\\n\\t\\tnative.Remove(name)\\n\\t\\treturn err\\n\\t}\\n\\tif err := tmp.Close(); err != nil {\\n\\t\\tnative.Remove(name)\\n\\t\\treturn err\\n\\t}\\n\\tif err := native.Rename(name, target); err != nil {\\n\\t\\tnative.Remove(name)\\n\\t\\treturn err\\n\\t}\\n\\treturn nil\"\n(The former CreateTemp/Write/Close/Rename sequence with no sync, through the native operations, ignoring the injected replacer.)\nResult: verdict bit, cause failed, failed_tests 1, restored yes. Mutated package run exit 1.\nFailure: TestDashboardReplacement replacement_test.go:33: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected dashboard directory-sync failure.\n\nProbe 2 (select 0644 by assumption; kept mode literal 0o600, whole package, no --run filter).\nSwap old (exact): \"pageMode fs.FileMode = 0o600\"\nSwap new (exact): \"pageMode fs.FileMode = 0o644\"\nResult: verdict bit, cause failed, failed_tests 1 of 23, restored yes. Mutated package run exit 1.\nFailure: TestDashboardReplacement replacement_test.go:34: .../bench-dashboard.html mode = -rw-r--r--, want a regular file at -rw-------.\n\nRestoration: bench probe restored the source (restored=yes); bench test --package ./internal/dashboard passes at ee901c01 (23 tests).\n"
            }
          }
        },
        {
          "id": "dc3-preservation-1",
          "performer": "claude:dcm_t10",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "4d77c1ef8c235ca9bc90f4a80f95c22bea5f2751",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_t10",
            "digest": "sha256:e0e9d4fcb6a55b89e22a959cfc423026cf8df24672fe911195040e639debe792",
            "excerpt": "D-C3 preservation probes at ee901c01b9c522f654a8524b675f1e94c4e148c2 (FT354-durable-callers), subject internal/dashboard/dashboard.go, bench test --package ./internal/dashboard.\n\nBaseline comparison against the frozen caller: the payload is the rendered page from Render (an unchanged encoder), so TestDashboardReplacement derives the expected bytes from it and needs no red. Mode 0600 is the frozen baseline (the implicit os.CreateTemp mode), kept as an independent literal; its necessity red is omission-proof probe 2 (pageMode 0o600 -> 0o644, whole package, only TestDashboardReplacement:34 red). Refusal category \"cannot write bench-dashboard.html\" with exit 1 and the stdout and git-dir authority paths are unchanged; the existing TestCommandArgs and TestDashboardWriteRefusesWhenTheReaderFails pass.\n\nProbe A (D36 refusal: fallback write of the new bytes after an unpublished error; whole package, no --run filter).\nSwap old (exact): \"\\treturn r.Replace(target, []byte(content), pageMode)\"\nSwap new (exact): \"\\terr := r.Replace(target, []byte(content), pageMode)\\n\\tif err != nil {\\n\\t\\tdurablefile.Replace(target, []byte(content), pageMode)\\n\\t}\\n\\treturn err\"\nResult: verdict bit, cause failed, failed_tests 1 of 23, restored yes. Mutated package run exit 1.\nFailure: only TestDashboardPreRenameFailure replacement_test.go:51 (the RequireBytes prior-page check): the page holds the rendered HTML, want \"<!DOCTYPE html>prior page\".\n\nProbe B (the exported Command entry that supplies durablefile.Native() must report its error; omission probe with --run '^TestDashboardCommandReportsReplacementFailure$').\nSwap old (exact): \"atomicWrite(durablefile.Native(), target, page); err != nil {\"\nSwap new (exact): \"atomicWrite(durablefile.Native(), target, page); err != nil && false {\"\nResult: verdict bit, cause failed, failed_tests 1, restored yes. Mutated package run exit 1.\nFailure: TestDashboardCommandReportsReplacementFailure replacement_test.go:70: Command(nil) = \".../.git/bench-dashboard.html\\n\", 0, want exit 1 and \"cannot write bench-dashboard.html\".\n\nRestoration: bench probe restored the source each time (restored=yes); bench test --package ./internal/dashboard passes at ee901c01 (23 tests).\n"
          },
          "requirement": "preservation",
          "command": "bench test --package ./internal/dashboard",
          "exit_code": 0,
          "probe": {
            "mutation": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:dcm_t10",
              "digest": "sha256:e0e9d4fcb6a55b89e22a959cfc423026cf8df24672fe911195040e639debe792",
              "excerpt": "D-C3 preservation probes at ee901c01b9c522f654a8524b675f1e94c4e148c2 (FT354-durable-callers), subject internal/dashboard/dashboard.go, bench test --package ./internal/dashboard.\n\nBaseline comparison against the frozen caller: the payload is the rendered page from Render (an unchanged encoder), so TestDashboardReplacement derives the expected bytes from it and needs no red. Mode 0600 is the frozen baseline (the implicit os.CreateTemp mode), kept as an independent literal; its necessity red is omission-proof probe 2 (pageMode 0o600 -> 0o644, whole package, only TestDashboardReplacement:34 red). Refusal category \"cannot write bench-dashboard.html\" with exit 1 and the stdout and git-dir authority paths are unchanged; the existing TestCommandArgs and TestDashboardWriteRefusesWhenTheReaderFails pass.\n\nProbe A (D36 refusal: fallback write of the new bytes after an unpublished error; whole package, no --run filter).\nSwap old (exact): \"\\treturn r.Replace(target, []byte(content), pageMode)\"\nSwap new (exact): \"\\terr := r.Replace(target, []byte(content), pageMode)\\n\\tif err != nil {\\n\\t\\tdurablefile.Replace(target, []byte(content), pageMode)\\n\\t}\\n\\treturn err\"\nResult: verdict bit, cause failed, failed_tests 1 of 23, restored yes. Mutated package run exit 1.\nFailure: only TestDashboardPreRenameFailure replacement_test.go:51 (the RequireBytes prior-page check): the page holds the rendered HTML, want \"<!DOCTYPE html>prior page\".\n\nProbe B (the exported Command entry that supplies durablefile.Native() must report its error; omission probe with --run '^TestDashboardCommandReportsReplacementFailure$').\nSwap old (exact): \"atomicWrite(durablefile.Native(), target, page); err != nil {\"\nSwap new (exact): \"atomicWrite(durablefile.Native(), target, page); err != nil && false {\"\nResult: verdict bit, cause failed, failed_tests 1, restored yes. Mutated package run exit 1.\nFailure: TestDashboardCommandReportsReplacementFailure replacement_test.go:70: Command(nil) = \".../.git/bench-dashboard.html\\n\", 0, want exit 1 and \"cannot write bench-dashboard.html\".\n\nRestoration: bench probe restored the source each time (restored=yes); bench test --package ./internal/dashboard passes at ee901c01 (23 tests).\n"
            }
          }
        },
        {
          "id": "dcm-t10-rep1-caller-tests",
          "performer": "claude:dcm_t10_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2edd23a96795019c66d36851373250f4b7b5329f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_t10_rep1",
            "digest": "sha256:cfc5d4dfe6b33fc1c796daed19ba1a835789a8dd075bf59e6b9a14af01c1459b",
            "excerpt": "D-C3 repair 1 caller-tests at b20919b10799df0f4014df0a588bf70c3cea986c (FT354-durable-callers, clean): bench test --package ./internal/dashboard\npackages: github.com/gibbonmi/bench/internal/dashboard,pass,tests_run=23; failures 0; skips 0.\n"
          },
          "requirement": "caller-tests",
          "command": "bench test --package ./internal/dashboard",
          "exit_code": 0
        },
        {
          "id": "dcm-t10-rep1-existing-canary-owner",
          "performer": "claude:dcm_t10_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2edd23a96795019c66d36851373250f4b7b5329f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_t10_rep1",
            "digest": "sha256:da334bfc639968b46555f5a48aa534b315c3f4c12311a77689c684bcf437a32a",
            "excerpt": "D-C3 repair 1 existing-canary-owner at b20919b10799df0f4014df0a588bf70c3cea986c (FT354-durable-callers, clean): bench test --check git-plumbing-owner\ncheck git-plumbing-owner, conformance, tests_run=1; internal/conformance pass; failures 0; skips 0.\nThe git-flag-retyped fixture's MUTATE.json old line \"\\tgitDir, err := git.AdminDir(root)\\n\" stays intact in dashboard.go; the repair does not touch it. bench test --check package-core-guard also passed at this tip (tests_run=1, failures 0).\n"
          },
          "requirement": "existing-canary-owner",
          "command": "bench test --check git-plumbing-owner",
          "exit_code": 0
        },
        {
          "id": "dcm-t10-rep1-root-conformance",
          "performer": "claude:dcm_t10_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2edd23a96795019c66d36851373250f4b7b5329f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_t10_rep1",
            "digest": "sha256:13e23f7ed338429afd33e71802dd1453faeb10834023680a17c57feed482843f",
            "excerpt": "D-C3 repair 1 root-conformance at b20919b10799df0f4014df0a588bf70c3cea986c (FT354-durable-callers, clean): bench test --check root-conformance\ncheck root-conformance, conformance, tests_run=1; internal/conformance pass; failures 0; skips 0.\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "dcm-t10-rep1-omission-proof",
          "performer": "claude:dcm_t10_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2edd23a96795019c66d36851373250f4b7b5329f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_t10_rep1",
            "digest": "sha256:578338373de1e045d5a33278e54674f3612808ff926532fb60e36e220052881f",
            "excerpt": "D-C3 repair 1 omission-proof at b20919b10799df0f4014df0a588bf70c3cea986c (FT354-durable-callers), subject internal/dashboard/dashboard.go, bench test --package ./internal/dashboard.\n\nThe repair changes no line of atomicWrite, pageMode, TestDashboardReplacement, or TestDashboardPreRenameFailure, so no omission-proof red is replayed.\nCited unchanged: /tmp/claude-1000/-home-mgibs-workspace-bench/a142325d-036b-4e84-b387-f6b875be269c/scratchpad/dcm_t10/p2/omission.txt Probe 1 (restore the former writer; owner bypass; --run '^TestDashboardReplacement$'; bit, restored yes, exit 1).\nCited unchanged: /tmp/claude-1000/-home-mgibs-workspace-bench/a142325d-036b-4e84-b387-f6b875be269c/scratchpad/dcm_t10/p2/omission.txt Probe 2 (pageMode 0o600 -> 0o644, whole package; only TestDashboardReplacement:34 red; bit, restored yes, exit 1).\nAfter the repair, bench test --package ./internal/dashboard passes at b20919b1 (23 tests).\n"
          },
          "requirement": "omission-proof",
          "command": "bench test --package ./internal/dashboard",
          "exit_code": 0,
          "probe": {
            "mutation": "Select 0644 by assumption or restore the former writer. The mode or injected-failure witness must turn red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:dcm_t10_rep1",
              "digest": "sha256:578338373de1e045d5a33278e54674f3612808ff926532fb60e36e220052881f",
              "excerpt": "D-C3 repair 1 omission-proof at b20919b10799df0f4014df0a588bf70c3cea986c (FT354-durable-callers), subject internal/dashboard/dashboard.go, bench test --package ./internal/dashboard.\n\nThe repair changes no line of atomicWrite, pageMode, TestDashboardReplacement, or TestDashboardPreRenameFailure, so no omission-proof red is replayed.\nCited unchanged: /tmp/claude-1000/-home-mgibs-workspace-bench/a142325d-036b-4e84-b387-f6b875be269c/scratchpad/dcm_t10/p2/omission.txt Probe 1 (restore the former writer; owner bypass; --run '^TestDashboardReplacement$'; bit, restored yes, exit 1).\nCited unchanged: /tmp/claude-1000/-home-mgibs-workspace-bench/a142325d-036b-4e84-b387-f6b875be269c/scratchpad/dcm_t10/p2/omission.txt Probe 2 (pageMode 0o600 -> 0o644, whole package; only TestDashboardReplacement:34 red; bit, restored yes, exit 1).\nAfter the repair, bench test --package ./internal/dashboard passes at b20919b1 (23 tests).\n"
            }
          }
        },
        {
          "id": "dcm-t10-rep1-preservation",
          "performer": "claude:dcm_t10_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2edd23a96795019c66d36851373250f4b7b5329f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_t10_rep1",
            "digest": "sha256:7245b308f1af75faa84a7a22f9805dc87a43613702b2e8872a1a4d77c4136e3b",
            "excerpt": "D-C3 repair 1 preservation probes at b20919b10799df0f4014df0a588bf70c3cea986c (FT354-durable-callers, clean), subject internal/dashboard/dashboard.go, bench test --package ./internal/dashboard.\n\nThe refusal now composes from pageName (\"cannot write \"+pageName); output bytes are unchanged. The test keeps \"cannot write bench-dashboard.html\" as one frozen baseline literal (const refusal), proven by Probe C.\n\nProbe C (new; necessity red for the kept refusal literal; whole package, no --run filter).\nSwap old (exact): \"\\\"cannot write \\\"+pageName\"\nSwap new (exact): \"\\\"cannot save \\\"+pageName\"\nResult: verdict bit, cause failed, failed_tests 1 of 23, restored yes. Mutated package run exit 1.\nFailure: only TestDashboardCommandReportsReplacementFailure replacement_test.go:71: Command(nil) = \"error: cannot save bench-dashboard.html — replace .../.git/bench-dashboard.html: rename: ...\".\n\nProbe B (replayed, because the repair moves its witness lines; --run '^TestDashboardCommandReportsReplacementFailure$'; swap text unchanged).\nSwap old (exact): \"atomicWrite(durablefile.Native(), target, page); err != nil {\"\nSwap new (exact): \"atomicWrite(durablefile.Native(), target, page); err != nil && false {\"\nResult: verdict bit, cause failed, failed_tests 1, restored yes. Mutated package run exit 1.\nFailure: TestDashboardCommandReportsReplacementFailure replacement_test.go:71: Command(nil) = \".../.git/bench-dashboard.html\\n\", 0, want exit 1 and \"cannot write bench-dashboard.html\".\n\nCited unchanged: /tmp/claude-1000/-home-mgibs-workspace-bench/a142325d-036b-4e84-b387-f6b875be269c/scratchpad/dcm_t10/p2/preservation.txt Probe A (fallback write after an unpublished error; whole package; only TestDashboardPreRenameFailure:51 red; bit, restored yes, exit 1).\n\nRestoration: bench probe restored the source each time (restored=yes, tree dirty=false); bench test --package ./internal/dashboard passes at b20919b1 (23 tests).\n"
          },
          "requirement": "preservation",
          "command": "bench test --package ./internal/dashboard",
          "exit_code": 0,
          "probe": {
            "mutation": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:dcm_t10_rep1",
              "digest": "sha256:7245b308f1af75faa84a7a22f9805dc87a43613702b2e8872a1a4d77c4136e3b",
              "excerpt": "D-C3 repair 1 preservation probes at b20919b10799df0f4014df0a588bf70c3cea986c (FT354-durable-callers, clean), subject internal/dashboard/dashboard.go, bench test --package ./internal/dashboard.\n\nThe refusal now composes from pageName (\"cannot write \"+pageName); output bytes are unchanged. The test keeps \"cannot write bench-dashboard.html\" as one frozen baseline literal (const refusal), proven by Probe C.\n\nProbe C (new; necessity red for the kept refusal literal; whole package, no --run filter).\nSwap old (exact): \"\\\"cannot write \\\"+pageName\"\nSwap new (exact): \"\\\"cannot save \\\"+pageName\"\nResult: verdict bit, cause failed, failed_tests 1 of 23, restored yes. Mutated package run exit 1.\nFailure: only TestDashboardCommandReportsReplacementFailure replacement_test.go:71: Command(nil) = \"error: cannot save bench-dashboard.html — replace .../.git/bench-dashboard.html: rename: ...\".\n\nProbe B (replayed, because the repair moves its witness lines; --run '^TestDashboardCommandReportsReplacementFailure$'; swap text unchanged).\nSwap old (exact): \"atomicWrite(durablefile.Native(), target, page); err != nil {\"\nSwap new (exact): \"atomicWrite(durablefile.Native(), target, page); err != nil && false {\"\nResult: verdict bit, cause failed, failed_tests 1, restored yes. Mutated package run exit 1.\nFailure: TestDashboardCommandReportsReplacementFailure replacement_test.go:71: Command(nil) = \".../.git/bench-dashboard.html\\n\", 0, want exit 1 and \"cannot write bench-dashboard.html\".\n\nCited unchanged: /tmp/claude-1000/-home-mgibs-workspace-bench/a142325d-036b-4e84-b387-f6b875be269c/scratchpad/dcm_t10/p2/preservation.txt Probe A (fallback write after an unpublished error; whole package; only TestDashboardPreRenameFailure:51 red; bit, restored yes, exit 1).\n\nRestoration: bench probe restored the source each time (restored=yes, tree dirty=false); bench test --package ./internal/dashboard passes at b20919b1 (23 tests).\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "dcm-c3-r1-standards",
          "performer": "claude:dcm_rev_c3_std",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "4d77c1ef8c235ca9bc90f4a80f95c22bea5f2751",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_c3_std",
            "digest": "sha256:b2c100900aca7d5d7bf043c48f9da45e831f9aa21dc9932eb0de2dfbf0700588",
            "excerpt": "## Standards (D-C3, c377e632..ee901c01) — 2 findings; worst S2.\nS1 | auto-fix (orchestrator, --auto-approve) | 6 | internal/dashboard/dashboard.go:88 vs :136 | the refusal literal \"cannot write bench-dashboard.html\" restates pageName | compose the message from pageName; output bytes unchanged\nS2 | auto-fix (same as C1) | 6 | internal/dashboard/replacement_test.go:69-70 | the kept refusal-category literal has no whole-package necessity red | record one whole-package probe that drifts only the category text\nAdvice: test doc comments restate body and mode (D-D sweep); dashboard.go:61-62 doc comment spells the file name.\n"
          },
          "axis": "Standards",
          "base": "c377e632d4056ebb6853ed1087efbb38e010c824",
          "tip": "ee901c01b9c522f654a8524b675f1e94c4e148c2",
          "finding_ids": [
            "S1",
            "S2"
          ],
          "supersedes": []
        },
        {
          "id": "dcm-c3-r1-spec",
          "performer": "claude:dcm_rev_c3_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "4d77c1ef8c235ca9bc90f4a80f95c22bea5f2751",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_c3_spec",
            "digest": "sha256:89ea53ce145b273df2141c8d4fa6f79951982d90d5924bd56ff3c98f7ead4144",
            "excerpt": "## Spec (D-C3, c377e632..ee901c01) — 0 findings.\nD26 met: shared owner, 0600 mode, page bytes, post-rename failure through Command, pre-rename old bytes; coverage cell cites TestDashboardReplacement. spec.md:199 is pinned to 85baca4d and is not stale.\nAdvice: no post-rename failure through Command; D36 tag noted for D-D.\n"
          },
          "axis": "Spec",
          "base": "c377e632d4056ebb6853ed1087efbb38e010c824",
          "tip": "ee901c01b9c522f654a8524b675f1e94c4e148c2",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "dcm-c3-r1-coverage",
          "performer": "claude:dcm_rev_c3_cov",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "4d77c1ef8c235ca9bc90f4a80f95c22bea5f2751",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_c3_cov",
            "digest": "sha256:5d3ebb23b1f8e0db21a013f30feb30c359ca80a46f6b77c2e4be14aa83376432",
            "excerpt": "## Coverage (D-C3, c377e632..ee901c01) — 1 finding; worst C1.\nC1 | auto-fix | 6 | internal/dashboard/replacement_test.go:69 | the refusal-message literal is a kept independent expectation with no whole-package necessity red | record one whole-package red that changes the message at dashboard.go:88, or derive the expectation\nD26, mode red, D36 fallback red, Command path red, and the five verification entries hold; D30 needs no witness.\nAdvice: Command success path untested (pre-existing); unsynced direct writer would stay green (D-D census).\n"
          },
          "axis": "Coverage",
          "base": "c377e632d4056ebb6853ed1087efbb38e010c824",
          "tip": "ee901c01b9c522f654a8524b675f1e94c4e148c2",
          "finding_ids": [
            "C1"
          ],
          "supersedes": []
        },
        {
          "id": "dcm-c3-r2-spec",
          "performer": "claude:dcm_rev_c3_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "2edd23a96795019c66d36851373250f4b7b5329f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_c3_spec",
            "digest": "sha256:acb4f8b4859fdbc2e068f83ea14f4a481ad3fcf1469d8c07e02a148ca8415a7c",
            "excerpt": "## Spec (D-C3 confirming, be27b3b8..b20919b1) — pass; 0 findings; refusal output bytes unchanged; D26 still met.\n"
          },
          "axis": "Spec",
          "base": "c377e632d4056ebb6853ed1087efbb38e010c824",
          "tip": "b20919b10799df0f4014df0a588bf70c3cea986c",
          "finding_ids": [],
          "supersedes": [
            "dcm-c3-r1-spec"
          ]
        },
        {
          "id": "dcm-c3-r2-standards",
          "performer": "claude:dcm_rev_c3_std",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "2edd23a96795019c66d36851373250f4b7b5329f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_c3_std",
            "digest": "sha256:16cca3a57f60462b58e401cc6c44876e6ab57cdc6c4e69bf3160cedbc04f8cf7",
            "excerpt": "## Standards (D-C3 confirming, be27b3b8..b20919b1) — S1 closed (dashboard.go:88); S2 closed (replacement_test.go:69-71, probe C whole package); 0 new findings.\nAdvice: the Command doc comment says \"that path\" with no named path; optional wording fix.\n"
          },
          "axis": "Standards",
          "base": "c377e632d4056ebb6853ed1087efbb38e010c824",
          "tip": "b20919b10799df0f4014df0a588bf70c3cea986c",
          "finding_ids": [],
          "supersedes": [
            "dcm-c3-r1-standards"
          ]
        },
        {
          "id": "dcm-c3-r2-coverage",
          "performer": "claude:dcm_rev_c3_cov",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "2edd23a96795019c66d36851373250f4b7b5329f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dcm_rev_c3_cov",
            "digest": "sha256:fc879f778756b63ede545d7b802334360da5741063ba006f490f5c7db15c33f4",
            "excerpt": "## Coverage (D-C3 confirming, be27b3b8..b20919b1) — C1 closed (replacement_test.go:69-71; probe C whole package, B replayed, orchestrator pageName drift bit); dcm-t10-rep1 entries graded; 0 new findings.\n"
          },
          "axis": "Coverage",
          "base": "c377e632d4056ebb6853ed1087efbb38e010c824",
          "tip": "b20919b10799df0f4014df0a588bf70c3cea986c",
          "finding_ids": [],
          "supersedes": [
            "dcm-c3-r1-coverage"
          ]
        }
      ]
    },
    {
      "id": "D-C6",
      "base": "3210553dcf40b324eeb975dd04758adbff2a02e9",
      "tip": "e2a36eb3bb2726e9c3edbdb3c6662898cf6b82c6",
      "plan_digest": "sha256:0a9c82011a15380bda3c64cc4ee87594038caa6080e5ec3d4f7ff5bef371daf6",
      "source_digest": "f2d751e5d1918bc60b6724a738649a3747a748d4",
      "acceptance_rows": [
        "D29"
      ],
      "verification": [
        {
          "id": "dc6-caller-tests-1",
          "performer": "claude:dcm_t13",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f2d751e5d1918bc60b6724a738649a3747a748d4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "e2a36eb3bb2726e9c3edbdb3c6662898cf6b82c6",
            "digest": "sha256:7917750030b24558683794abd46aabf0cfe58f491fdfe357e126e65ccb9d117e",
            "excerpt": "bench test --package ./internal/repairpilot at FT354-durable-callers e2a36eb3bb2726e9c3edbdb3c6662898cf6b82c6 (clean tree)\npackages: github.com/gibbonmi/bench/internal/repairpilot,pass,175,136\nfailures: 0, skips: 0\n"
          },
          "requirement": "caller-tests",
          "command": "bench test --package ./internal/repairpilot",
          "exit_code": 0
        },
        {
          "id": "dc6-root-conformance-1",
          "performer": "claude:dcm_t13",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f2d751e5d1918bc60b6724a738649a3747a748d4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "e2a36eb3bb2726e9c3edbdb3c6662898cf6b82c6",
            "digest": "sha256:2d26c84cf179dc83d6f0daeceec0381f53441dbfe40b12af89e5c45440dfd728",
            "excerpt": "bench test --check root-conformance at FT354-durable-callers e2a36eb3bb2726e9c3edbdb3c6662898cf6b82c6 (clean tree)\ncheck: root-conformance,conformance,1 test\npackages: github.com/gibbonmi/bench/internal/conformance,pass,17932,1\nfailures: 0, skips: 0\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "dc6-omission-proof-1",
          "performer": "claude:dcm_t13",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f2d751e5d1918bc60b6724a738649a3747a748d4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "e2a36eb3bb2726e9c3edbdb3c6662898cf6b82c6",
            "digest": "sha256:de2cb074847965624660558a48d873e59290353394e6f379efcffafad0923987",
            "excerpt": "omission-proof for D29 at FT354-durable-callers e2a36eb3bb2726e9c3edbdb3c6662898cf6b82c6 (clean tree).\nPlan probe: restore the local file-only sequence. Package run: bench test --package ./internal/repairpilot passes (136 tests).\n\nbench probe internal/repairpilot/command.go --package ./internal/repairpilot --run TestPilotReplacement\n  OLD: return files.Replacer.Replace(path, append(data, '\\n'), 0o600)\n  NEW: temporary, err := files.CreateTemp(filepath.Dir(path), \".pilot-*\"); if err != nil { return err }; defer files.Remove(temporary.Name()); if err := temporary.Chmod(0o600); err != nil { return err }; if _, err := temporary.Write(append(data, '\\n')); err != nil { return err }; if err := temporary.Sync(); err != nil { return err }; if err := temporary.Close(); err != nil { return err }; return files.Rename(temporary.Name(), path)\n  verdict=bit cause=failed failed_tests=1 restored=yes mutated package exit 1\n  TestPilotReplacement replacement_test.go:22: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected pilot directory-sync failure\n\nThe weakened short-write alternative ran in phase 1 at dcm-t13 3e9301f7 (internal/durablefile/replace.go \"if err == nil && n < len(data) {\" -> \"if err == nil && n < 0 {\"): bit, TestRepairPilotStorage/write-failure command_test.go:193 \"partial write = <nil>; want short write\", restored=yes. The plan names it as an alternative (\"or\"), so it was not rerun.\n"
          },
          "requirement": "omission-proof",
          "command": "bench test --package ./internal/repairpilot",
          "exit_code": 0,
          "probe": {
            "mutation": "Restore the local file-only sequence or weaken the short-write assertion. The corresponding real Store.Replace witness must turn red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "e2a36eb3bb2726e9c3edbdb3c6662898cf6b82c6",
              "digest": "sha256:de2cb074847965624660558a48d873e59290353394e6f379efcffafad0923987",
              "excerpt": "omission-proof for D29 at FT354-durable-callers e2a36eb3bb2726e9c3edbdb3c6662898cf6b82c6 (clean tree).\nPlan probe: restore the local file-only sequence. Package run: bench test --package ./internal/repairpilot passes (136 tests).\n\nbench probe internal/repairpilot/command.go --package ./internal/repairpilot --run TestPilotReplacement\n  OLD: return files.Replacer.Replace(path, append(data, '\\n'), 0o600)\n  NEW: temporary, err := files.CreateTemp(filepath.Dir(path), \".pilot-*\"); if err != nil { return err }; defer files.Remove(temporary.Name()); if err := temporary.Chmod(0o600); err != nil { return err }; if _, err := temporary.Write(append(data, '\\n')); err != nil { return err }; if err := temporary.Sync(); err != nil { return err }; if err := temporary.Close(); err != nil { return err }; return files.Rename(temporary.Name(), path)\n  verdict=bit cause=failed failed_tests=1 restored=yes mutated package exit 1\n  TestPilotReplacement replacement_test.go:22: err = <nil>, want a *durablefile.Error at stage directory-sync with Published=true that wraps injected pilot directory-sync failure\n\nThe weakened short-write alternative ran in phase 1 at dcm-t13 3e9301f7 (internal/durablefile/replace.go \"if err == nil && n < len(data) {\" -> \"if err == nil && n < 0 {\"): bit, TestRepairPilotStorage/write-failure command_test.go:193 \"partial write = <nil>; want short write\", restored=yes. The plan names it as an alternative (\"or\"), so it was not rerun.\n"
            }
          }
        },
        {
          "id": "dc6-preservation-1",
          "performer": "claude:dcm_t13",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f2d751e5d1918bc60b6724a738649a3747a748d4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "e2a36eb3bb2726e9c3edbdb3c6662898cf6b82c6",
            "digest": "sha256:6812936303c0c8e5f5f9bac475864b5cabd34995e8903fe880bbf00a4bb25ab5",
            "excerpt": "preservation for D29 at FT354-durable-callers e2a36eb3bb2726e9c3edbdb3c6662898cf6b82c6 (clean tree). Package run: bench test --package ./internal/repairpilot passes (136 tests).\nKept literals (spec inventory row D29 freezes 0600 and the indented document plus final newline): the mode 0o600 and the expected bytes from the test-side documentBytes encoder in TestPilotReplacement. Each restored=yes, mutated package exit 1.\n\n1. Bytes necessity (whole package, no --run): command.go\n  OLD: data, err := json.MarshalIndent(document, \"\", \"  \")\n  NEW: data, err := json.MarshalIndent(document, \"\", \"\\t\")\n  bit, failed_tests=1: TestPilotReplacement replacement_test.go:23 (bytes mismatch). Only the kept expectation catches it.\n\n2. Mode necessity (whole package, no --run): command.go\n  OLD: return files.Replacer.Replace(path, append(data, '\\n'), 0o600)\n  NEW: return files.Replacer.Replace(path, append(data, '\\n'), 0o640)\n  bit, failed_tests=2: TestPilotReplacement replacement_test.go:23 (mode -rw-r----- want -rw-------), TestRepairPilotStorage/custody command_test.go:146.\n\n3. Pre-rename refusal (D36 row, --run TestRepairPilotStorage/replace-failure): fallback write of the new bytes after an unpublished error, command.go\n  OLD: return files.Replacer.Replace(path, append(data, '\\n'), 0o600)\n  NEW: if err := files.Replacer.Replace(path, append(data, '\\n'), 0o600); err != nil { var failure *durablefile.Error; if errors.As(err, &failure) && !failure.Published { _ = os.WriteFile(path, append(data, '\\n'), 0o600) }; return err }; return nil\n  bit, failed_tests=1: TestRepairPilotStorage/replace-failure command_test.go:214, only at the old-bytes check (RequireBytes).\n\n4. D30 refusal witness (--run TestPilotDirectoryRefusal): swallow the pilot directory creation refusal, command.go withPilotLock\n  OLD: if err := files.MkdirAll(directory, 0o700); err != nil {\n  NEW: if err := files.MkdirAll(directory, 0o700); err != nil && false {\n  bit, failed_tests=1: TestPilotDirectoryRefusal replacement_test.go:40 (output names the chmod failure, not the injected refusal).\n"
          },
          "requirement": "preservation",
          "command": "bench test --package ./internal/repairpilot",
          "exit_code": 0,
          "probe": {
            "mutation": "Compare real caller payload, mode, refusal, authority, and output assertions with the frozen baseline. Keep independent expected bytes and demonstrate the named omission red, restoration, and pass.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "e2a36eb3bb2726e9c3edbdb3c6662898cf6b82c6",
              "digest": "sha256:6812936303c0c8e5f5f9bac475864b5cabd34995e8903fe880bbf00a4bb25ab5",
              "excerpt": "preservation for D29 at FT354-durable-callers e2a36eb3bb2726e9c3edbdb3c6662898cf6b82c6 (clean tree). Package run: bench test --package ./internal/repairpilot passes (136 tests).\nKept literals (spec inventory row D29 freezes 0600 and the indented document plus final newline): the mode 0o600 and the expected bytes from the test-side documentBytes encoder in TestPilotReplacement. Each restored=yes, mutated package exit 1.\n\n1. Bytes necessity (whole package, no --run): command.go\n  OLD: data, err := json.MarshalIndent(document, \"\", \"  \")\n  NEW: data, err := json.MarshalIndent(document, \"\", \"\\t\")\n  bit, failed_tests=1: TestPilotReplacement replacement_test.go:23 (bytes mismatch). Only the kept expectation catches it.\n\n2. Mode necessity (whole package, no --run): command.go\n  OLD: return files.Replacer.Replace(path, append(data, '\\n'), 0o600)\n  NEW: return files.Replacer.Replace(path, append(data, '\\n'), 0o640)\n  bit, failed_tests=2: TestPilotReplacement replacement_test.go:23 (mode -rw-r----- want -rw-------), TestRepairPilotStorage/custody command_test.go:146.\n\n3. Pre-rename refusal (D36 row, --run TestRepairPilotStorage/replace-failure): fallback write of the new bytes after an unpublished error, command.go\n  OLD: return files.Replacer.Replace(path, append(data, '\\n'), 0o600)\n  NEW: if err := files.Replacer.Replace(path, append(data, '\\n'), 0o600); err != nil { var failure *durablefile.Error; if errors.As(err, &failure) && !failure.Published { _ = os.WriteFile(path, append(data, '\\n'), 0o600) }; return err }; return nil\n  bit, failed_tests=1: TestRepairPilotStorage/replace-failure command_test.go:214, only at the old-bytes check (RequireBytes).\n\n4. D30 refusal witness (--run TestPilotDirectoryRefusal): swallow the pilot directory creation refusal, command.go withPilotLock\n  OLD: if err := files.MkdirAll(directory, 0o700); err != nil {\n  NEW: if err := files.MkdirAll(directory, 0o700); err != nil && false {\n  bit, failed_tests=1: TestPilotDirectoryRefusal replacement_test.go:40 (output names the chmod failure, not the injected refusal).\n"
            }
          }
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
    },
    {
      "from": "sha256:5cb50643a5295119fe76284a8cc58c0cc7bfd797ae01bfecc20e1b30a5df98af",
      "to": "sha256:b4e04ec302a9d6cac2eb040e717eda8884a27403379c7b451d89fbf778b7bda6",
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
      "from": "sha256:b4e04ec302a9d6cac2eb040e717eda8884a27403379c7b451d89fbf778b7bda6",
      "to": "sha256:b46bfce6efc8865c4e196a7c89854bdcb94e436d4d1b1c1b5bd58a67d573c3b2",
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
        "D-B6": [
          "D-B6"
        ],
        "D-C1": [
          "D-C1"
        ]
      }
    },
    {
      "from": "sha256:b46bfce6efc8865c4e196a7c89854bdcb94e436d4d1b1c1b5bd58a67d573c3b2",
      "to": "sha256:0b967f69bea388bd20a0fb6062db10215905fe7652985d627efd96fbef58f01d",
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
        "D-B6": [
          "D-B6"
        ],
        "D-C1": [
          "D-C1"
        ]
      }
    },
    {
      "from": "sha256:0b967f69bea388bd20a0fb6062db10215905fe7652985d627efd96fbef58f01d",
      "to": "sha256:93a6235e5beda74d1bde28cbbad80df6658b3235cc2436e55544d1e14d979dce",
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
        "D-B6": [
          "D-B6"
        ],
        "D-C1": [
          "D-C1"
        ],
        "D-C2": [
          "D-C2"
        ]
      }
    },
    {
      "from": "sha256:93a6235e5beda74d1bde28cbbad80df6658b3235cc2436e55544d1e14d979dce",
      "to": "sha256:72b9fa252493911657c1f94490e46fcf1041db627c4c8802b35acba37d1e8a29",
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
        "D-B6": [
          "D-B6"
        ],
        "D-C1": [
          "D-C1"
        ],
        "D-C2": [
          "D-C2"
        ]
      }
    },
    {
      "from": "sha256:72b9fa252493911657c1f94490e46fcf1041db627c4c8802b35acba37d1e8a29",
      "to": "sha256:2649d87242eda6170ea76efadda6c82fe48de463f63982d2b8d7f9e9c64f9323",
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
        "D-B6": [
          "D-B6"
        ],
        "D-C1": [
          "D-C1"
        ],
        "D-C2": [
          "D-C2"
        ]
      }
    },
    {
      "from": "sha256:2649d87242eda6170ea76efadda6c82fe48de463f63982d2b8d7f9e9c64f9323",
      "to": "sha256:5f3b8ff2cb08f262254461d3182e5bd26885f2f80ef93f6792bdd716587ef9a2",
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
        "D-B6": [
          "D-B6"
        ],
        "D-C1": [
          "D-C1"
        ],
        "D-C2": [
          "D-C2"
        ],
        "D-C4": [
          "D-C4"
        ]
      }
    },
    {
      "from": "sha256:5f3b8ff2cb08f262254461d3182e5bd26885f2f80ef93f6792bdd716587ef9a2",
      "to": "sha256:016cf63ba39a08352f69ad366b45a159ad6aca725a2e36a7620f02367440b05b",
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
        "D-B6": [
          "D-B6"
        ],
        "D-C1": [
          "D-C1"
        ],
        "D-C2": [
          "D-C2"
        ],
        "D-C4": [
          "D-C4"
        ]
      }
    },
    {
      "from": "sha256:016cf63ba39a08352f69ad366b45a159ad6aca725a2e36a7620f02367440b05b",
      "to": "sha256:caaaa230b74e8cf3a79cb22fa805c3a28bafb054e1b4053e68684ed4a0dee7cf",
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
        "D-B6": [
          "D-B6"
        ],
        "D-C1": [
          "D-C1"
        ],
        "D-C2": [
          "D-C2"
        ],
        "D-C3": [
          "D-C3"
        ],
        "D-C4": [
          "D-C4"
        ]
      }
    },
    {
      "from": "sha256:caaaa230b74e8cf3a79cb22fa805c3a28bafb054e1b4053e68684ed4a0dee7cf",
      "to": "sha256:0a9c82011a15380bda3c64cc4ee87594038caa6080e5ec3d4f7ff5bef371daf6",
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
        "D-B6": [
          "D-B6"
        ],
        "D-C1": [
          "D-C1"
        ],
        "D-C2": [
          "D-C2"
        ],
        "D-C3": [
          "D-C3"
        ],
        "D-C4": [
          "D-C4"
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

## Standards (D-B6)

Chunk D-B6 has 5 findings. The worst issue is S1.

- S1, auto-fix, confidence 8: `internal/commitment/repository/replacement_test.go:96` builds the intent lock path again. Use `intent.LockPath`.
- S2, auto-fix by the --auto-approve plan expansion, confidence 5: the test router at `replacement_test.go:78-138` repeats the router at `internal/gate/durable_replacement_test.go:231-269`. Move one router into `internal/durablefile/durablefiletest`, and use it in both tests.
- S3, auto-fix by the --auto-approve plan expansion, confidence 5: the repo fixture at `replacement_test.go:24,36-76` repeats `commitmenttest.Repo`. Expose `approve` through an `export_test.go`, move the test to the external package, and start from `commitmenttest.Repo`.
- S4, auto-fix, confidence 6: the setup closures at `replacement_test.go:277-283` call the parent test from a subtest. Pass the subtest into each setup.
- S5, auto-fix, confidence 4: the lock-release check after the call at `replacement_test.go:262-264` repeats a fact that intent owns. Delete it.

## Spec (D-B6)

Chunk D-B6 has 0 findings.

## Coverage (D-B6)

Chunk D-B6 has 8 findings. The worst issue is C1.

- C1, auto-fix, confidence 9: `restoreAll` at `internal/commitment/repository/replacement.go:68` can stop at the first failed restore unseen. Add a D50 row in which a first restore fails and a later restore succeeds.
- C2, auto-fix, confidence 9: a skipped restore after an unpublished board failure at `replacement.go:49` is silent. Add a board rename row that asserts the old policy bytes and mode.
- C3, auto-fix, confidence 8: an absent snapshot restore at `replacement.go:103` that returns early is silent. Add a fixture with no policy, and assert the policy stays absent after the policy and ledger failures.
- C4, auto-fix, confidence 8: a read-back at `replacement.go:113` that ignores bytes is silent. Add a read-back row that rewrites the bytes.
- C5, auto-fix, confidence 8: a dropped board snapshot error at `replacement.go:26` is silent. Add a board-directory refusal row.
- C6, auto-fix, confidence 8: a dropped sequence projection error at `replacement.go:36` is silent. Add a refusal row for a board with no sequence section.
- C7, auto-fix, confidence 6: this finding is the same as S5.
- C8, auto-fix, confidence 9: the two D30 no-file-written checks need their recorded reds. Record the demonstrated reds H1 and H2.

## Repair allowance (D-B6)

Chunk D-B6 used 0 of its 2 repair cycles before repair 1. Every finding falls on ticket 07, so one combined repair session takes them.

- S1 closes in commit 91954b32. The approval test gets the lock path from `intent.LockPath`. The lock-placement swap still turns the lock check red at the repair tip.
- S2 closes in commit 91954b32. `durablefiletest.Router` routes each replacement by destination and ordinal to its own Recorder. The approval test and the gate test use it. The `dcm-t07-rep1-gate-tests` excerpt records each D-B1 omission red through the router.
- S3 closes in commit 91954b32. `ApproveWith` in the export file exposes `approve`. The approval tests run in the external package from `commitmenttest.Repo`, and the local repository fixture is deleted.
- S4 closes in commit 91954b32. Each `TestApprovalRefusal` setup takes the subtest `t`.
- S5 and C7 close in commit 91954b32. The lock-release check after the call is deleted, because intent owns that fact.
- C1 closes in commit 91954b32 with `TestApprovalRestoreFailure/first-of-two`. A `break` after the first failed restore turns it red at the restored board.
- C2 closes in commit 91954b32 with `TestApprovalBoardFailure/rename`. A skipped restore after an unpublished board error turns it red at the policy mode.
- C3 closes in commit 91954b32 with `TestApprovalFailureWithoutPolicy`. An absent-snapshot restore that returns nil turns both rows red at the policy absence.
- C4 closes in commit 91954b32 with `TestApprovalRestoreFailure/read-back-bytes`. A read-back that ignores bytes turns it red.
- C5 closes in commit 91954b32 with `TestApprovalRefusal/board-directory`. A dropped board snapshot error turns it red.
- C6 closes in commit 91954b32 with `TestApprovalRefusal/board-without-sequence`. A dropped ProjectSequence error turns it red.
- C8 closes in commit 91954b32. H1 writes the board natively before the policy-snapshot refusal and turns only the policy-directory row red. H2 writes the policy natively before the decision refusal and turns only the decision row red.

Before the repair, each C1 to C6 swap was silent over the whole package. The `claude:dcm_t07_rep1` omission-proof and preservation excerpts record each red at the repair tip.

Chunk D-B6 used 1 of its 2 repair cycles.

The D-B6 confirming round at 77835db3 closed S1 through S5 and C1 through C8 with no new finding. Chunk D-B6 closed with 1 of its 2 repair cycles used.

## Standards (D-C2)

Chunk D-C2 has 2 findings. The worst issue is S1.

- S1, auto-fix, confidence 7: `TemporaryPrefix` at `internal/brokermanifest/brokermanifest.go:25-29` restates the temp naming that the leaf owns, and no test ties the two. Pin the prefix against the Recorder temp name, record the red, and correct the comment.
- S2, auto-fix by the --auto-approve plan expansion, confidence 5: the absence check at `internal/brokermanifest/replacement_test.go:77-82` repeats a check of other caller tests. Add `durablefiletest.RequireAbsent`, and use it here.

## Spec (D-C2)

Chunk D-C2 has 0 findings.

## Coverage (D-C2)

Chunk D-C2 has 3 findings. The worst issue is C1.

- C1, auto-fix, confidence 9: this finding is the same as S1.
- C2, auto-fix, confidence 8: a public `Write` at `internal/brokermanifest/brokermanifest.go:37` that swallows every error is silent. Test that the public entry returns its refusal, and record the red of the swallow.
- C3, no-op, confidence 6: the symlink resolution at `brokermanifest.go:42-44` has no witness, but the gap predates this delta and sits outside D25. A learning records it.

## Repair allowance (D-C2)

Chunk D-C2 used 0 of its 2 repair cycles before repair 1. S1, S2, and C2 fall on ticket 09, so one combined repair session takes them.

- S1 and C1 close in commit 19cb2923. `TestBrokerReplacement` compares the staging name up to its last dot with `TemporaryPrefix`, and the drifted prefix `".bench-broker."` turns only that test red. The constant's comment states that its value matches the staging name for `Name` and that a test pins it.
- S2 closes in commit 2fd8872c. `durablefiletest.RequireAbsent` holds the absence check, and the local `requireNoManifest` is gone.
- C2 closes in commit 2fd8872c. `TestBrokerWriteReturnsRefusal` calls the public `Write` with a directory as the broker. A public `Write` that returns nil for every error turns only that test red.

The `claude:dcm_t09_rep1` omission-proof and preservation excerpts record each red at 19cb2923. Chunk D-C2 used 1 of its 2 repair cycles.

The D-C2 confirming round at 5d295f64 closed S1, S2, C1, and C2. It raised one comment-only finding.

- S3, auto-fix as an evidence-only comment correction, confidence 6: the comment at `internal/brokermanifest/replacement_test.go:85` names the wrong refusal step. The orchestrator corrects it, and the Standards axis reaffirms it.

The Standards reaffirmation at b92b47f9 closed S3. Chunk D-C2 closed with 1 of its 2 repair cycles used, and the comment-only correction used none.

## Standards (D-C4)

Chunk D-C4 has 3 findings. The worst issue is S1.

- S1, auto-fix, confidence 4: the test at `internal/probe/replacement_test.go:105` cuts the output on a hand-written `preserved` table header. That header is a third copy of the `preservedTable` encoder fact. Scope the row without a restated header.
- S2, auto-fix, confidence 4: the shared `preservedTable` comment at `internal/probe/probe.go:217-219` speaks of a verdict row that only `render` has. State the why without a subject that one caller owns.
- S3, auto-fix, confidence 5: the two first `dcm-t11` probe entries disagree with their `-r2` twins. This finding is the same as C2.

## Spec (D-C4)

Chunk D-C4 has 1 finding. The worst issue is P1.

- P1, auto-fix by --auto-approve, confidence 5: `preservedTable` refuses a subject base name with a control byte. The exit-2 output at `internal/probe/probe.go:206-209` then drops the mutation cause and the copy path. The spec says that the refusal policy must not hide the preserved path. The `render` branch at `probe.go:246-249` has the same gap, and the repair closes both branches.

## Coverage (D-C4)

Chunk D-C4 has 3 findings. The worst issue is C2.

- C2, auto-fix, confidence 8: the graded `-r2` probe entries record probe exit 0, and the checkpoint refuses that value. Record entries with exit code 0 and a nonzero probe exit code.
- C1, auto-fix, confidence 7: the table-refusal branch of `mutationFailure` has no red. Add a witness with a BEL byte in the base name, a published mutation fault, and a restore fault.
- C3, no-op, confidence 6: no D30 refusal witness goes through the replacer-taking form. A logged learning routes D30 to ticket 14.

## Repair allowance (D-C4)

Chunk D-C4 used 0 of its 2 repair cycles before repair 1. S1, S2, S3, P1, C1, and C2 fall on ticket 11, so one combined repair session takes them.

- P1 closes in commit c377e632. `TestProbeNamesACopyTheEncoderRefuses` proves that both exit-2 branches name the kept copy and each cause.
- C1 closes in commit c377e632. A fallback that omits the copy path turns only `TestProbeNamesACopyTheEncoderRefuses` red over the whole package.
- S1 closes in commit c377e632. `TestProbeMutationRecoveryFailure` scopes the preserved block after the cause line and restates no table header.
- S2 closes in commit c377e632. The shared `keptCopy` comment states why the copy path stays visible, and it names no caller-specific subject.
- S3 and C2 close in the `dcm-t11-rep1` entries. Each probe entry records exit code 0 and probe exit code 1.

Chunk D-C4 used 1 of its 2 repair cycles.

The D-C4 confirming round at 08eb67b3 closed S1, S2, S3, P1, C1, and C2. It raised one comment-only finding.

- S4, auto-fix as an evidence-only comment correction, confidence 5: the `render` comment at `internal/probe/probe.go:231-232` says that a failed restore keeps the preserved row. When the copy path refuses, `keptCopy` answers with two error lines instead. The orchestrator corrects the comment, and the Standards axis reaffirms it.

The Standards reaffirmation at 5d56efc1 closed S4. Chunk D-C4 closed with 1 of its 2 repair cycles used, and the comment-only correction used none.

D-C4 is refrozen on the D-C2 tip 19cb2923, because the D-C2 comment-only correction changes the base digest that the chain check compares. The three axes confirmed their results for the wider pair with no new findings.

## Standards (D-C3)

Chunk D-C3 has 2 findings. The worst issue is S2.

- S1, auto-fix by --auto-approve, confidence 6: the refusal at `internal/dashboard/dashboard.go:88` spells the file name that `pageName` owns. Compose the message from `pageName`, so the output bytes stay the same.
- S2, auto-fix, confidence 6: this finding is the same as C1.

## Spec (D-C3)

Chunk D-C3 has 0 findings.

## Coverage (D-C3)

Chunk D-C3 has 1 finding. The worst issue is C1.

- C1, auto-fix, confidence 6: the test at `internal/dashboard/replacement_test.go:69` keeps the refusal text as an independent literal, and no whole-package red proves it. Record one whole-package red that changes only that text.

## Repair allowance (D-C3)

Chunk D-C3 used 0 of its 2 repair cycles before repair 1. S1, S2, and C1 fall on ticket 10, so one combined repair session takes them.

- S1 closes in commit b20919b1. The `Command` refusal composes its text from `pageName`, and `TestDashboardCommandReportsReplacementFailure` shows the same output bytes.
- S2 closes in commit b20919b1. `TestDashboardCommandReportsReplacementFailure` names the frozen refusal literal once.
- C1 closes in commit b20919b1. The swap of "cannot write" to "cannot save" turns only `TestDashboardCommandReportsReplacementFailure` red over the whole package.

Chunk D-C3 used 1 of its 2 repair cycles.

The D-C3 confirming round at 3316f631 closed S1, S2, and C1 with no new findings. Chunk D-C3 closed with 1 of its 2 repair cycles used.
