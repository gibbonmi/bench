# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/durable-caller-migration/spec.md",
  "plan_digest": "sha256:e952476f7701550400bde51f65f6db4fa0961f9d5b42c481a78e78b23e168ed5",
  "implementation_session": "",
  "chunks": [
    {
      "id": "D-B1",
      "base": "69cffaf6a7416abab014e5e1b352a683f3cbd566",
      "tip": "c87faf66abc2a8eb228b7179aeb57d0a075f9dd6",
      "plan_digest": "sha256:e952476f7701550400bde51f65f6db4fa0961f9d5b42c481a78e78b23e168ed5",
      "source_digest": "2c57c2d6e8525a84409078f842647f628de3faa9",
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
  }
}
```
