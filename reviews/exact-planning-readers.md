# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/exact-planning-readers/spec.md",
  "plan_digest": "sha256:3ac7ce83999acc3cf7e20f551f4a8a52889192313d46706167b7c957a273b330",
  "implementation_session": "",
  "chunks": [
    {
      "id": "ER-C1",
      "base": "c02a0324c839ab93f360781d06da7e8b7c179179",
      "tip": "8b89063e28fc3b5ede6b0d6250d23089eee0e0c9",
      "plan_digest": "sha256:3ac7ce83999acc3cf7e20f551f4a8a52889192313d46706167b7c957a273b330",
      "source_digest": "c64671b495882fb808603035026f8735babcf984",
      "acceptance_rows": [
        "ER1",
        "ER16",
        "ER17",
        "ER18",
        "ER19",
        "ER20",
        "ER21",
        "ER22",
        "ER23",
        "ER24",
        "ER25",
        "ER26",
        "ER27",
        "ER28",
        "ER35",
        "ER36",
        "ER37",
        "ER38",
        "ER39",
        "ER40",
        "ER44",
        "ER45",
        "ER46"
      ],
      "verification": [
        {
          "id": "ER-C1-successor-red",
          "performer": "ft125-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2edc7d7c59d6ac67660c1aa947cecf41b3c9c39e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft125-author-01",
            "digest": "sha256:6930f08283d31f41eda8deb2a55c8b5d3eb9957de14af7813cd22f3d797f58dc",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/responsebound/read.go,swap,failed,6,TestExactReaderCLIReconstruction/empty.txt,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestExactReaderCLIReconstruction,passed,8\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,2308\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,700,8\nfailures[6]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/empty.txt,\"planning_reader_commands_test.go:96: [\\\"read\\\" \\\"empty.txt\\\" \\\"--offset\\\" \\\"1\\\" \\\"--source\\\" \\\"61b45b3dedeb2168900a7220792c78d867d65ea2f4294ab19ca255928d15a069\\\"] = stdout=\\\"error: cursor out of range — offset 1 is outside the 0 selected bytes; read aga… (329 bytes)\",1\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/exported_evidence,\"planning_reader_commands_test.go:96: [\\\"read\\\" \\\"/dev/shm/IWND7Z/t/TestExactReaderCLIReconstruction155984913/002/evidence/source-1.txt\\\" \\\"--artifact\\\" \\\"--offset\\\" \\\"9801\\\" \\\"--source\\\" \\\"e27548a8e9e63ad4f3becf814d376bb682756968497e5454df04d75dd59cb13c… (514 bytes)\",1\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/invalid.bin,\"planning_reader_commands_test.go:96: [\\\"read\\\" \\\"invalid.bin\\\" \\\"--offset\\\" \\\"6003\\\" \\\"--source\\\" \\\"b608671a06e4befce4f041cbffc6441e0335b952a3e60caa420f89950682f624\\\"] = stdout=\\\"error: cursor out of range — offset 6003 is outside the 6002 selected byte… (342 bytes)\",1\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/mixed.txt,\"planning_reader_commands_test.go:96: [\\\"read\\\" \\\"mixed.txt\\\" \\\"--offset\\\" \\\"11239\\\" \\\"--source\\\" \\\"46b85bac9811733299542453c808777d3c375b11c6f4817cbc8032747bce16ab\\\"] = stdout=\\\"error: cursor out of range — offset 11239 is outside the 11238 selected byt… (341 bytes)\",1\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/spill,\"planning_reader_commands_test.go:96: [\\\"read\\\" \\\"/dev/shm/IWND7Z/t/bench-cmd-home-4247697596/responses/none/primary/1793400900665359220-c4b656308f3c9620.out\\\" \\\"--artifact\\\" \\\"--offset\\\" \\\"5201\\\" \\\"--source\\\" \\\"b48f8496e09065051f4598387fa3b82567a4c327d1… (558 bytes)\",1\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/unicode.txt,\"planning_reader_commands_test.go:96: [\\\"read\\\" \\\"unicode.txt\\\" \\\"--offset\\\" \\\"9601\\\" \\\"--source\\\" \\\"f5e77d45e198b9c37c8593b20cff4af9677890e880235a91536703849c4f81db\\\"] = stdout=\\\"error: cursor out of range — offset 9601 is outside the 9600 selected byte… (342 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "97ef9ee0e680e9cc4de6f1709b09664dd6cafd0e",
          "requirement": "ER-C1-successor-red",
          "command": "bench test --package ./cmd/bench --run TestExactReaderCLIReconstruction",
          "exit_code": 0,
          "probe": {
            "mutation": "skip one source byte in continuation successor",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft125-author-01",
              "digest": "sha256:6930f08283d31f41eda8deb2a55c8b5d3eb9957de14af7813cd22f3d797f58dc",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/responsebound/read.go,swap,failed,6,TestExactReaderCLIReconstruction/empty.txt,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestExactReaderCLIReconstruction,passed,8\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,2308\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,700,8\nfailures[6]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/empty.txt,\"planning_reader_commands_test.go:96: [\\\"read\\\" \\\"empty.txt\\\" \\\"--offset\\\" \\\"1\\\" \\\"--source\\\" \\\"61b45b3dedeb2168900a7220792c78d867d65ea2f4294ab19ca255928d15a069\\\"] = stdout=\\\"error: cursor out of range — offset 1 is outside the 0 selected bytes; read aga… (329 bytes)\",1\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/exported_evidence,\"planning_reader_commands_test.go:96: [\\\"read\\\" \\\"/dev/shm/IWND7Z/t/TestExactReaderCLIReconstruction155984913/002/evidence/source-1.txt\\\" \\\"--artifact\\\" \\\"--offset\\\" \\\"9801\\\" \\\"--source\\\" \\\"e27548a8e9e63ad4f3becf814d376bb682756968497e5454df04d75dd59cb13c… (514 bytes)\",1\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/invalid.bin,\"planning_reader_commands_test.go:96: [\\\"read\\\" \\\"invalid.bin\\\" \\\"--offset\\\" \\\"6003\\\" \\\"--source\\\" \\\"b608671a06e4befce4f041cbffc6441e0335b952a3e60caa420f89950682f624\\\"] = stdout=\\\"error: cursor out of range — offset 6003 is outside the 6002 selected byte… (342 bytes)\",1\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/mixed.txt,\"planning_reader_commands_test.go:96: [\\\"read\\\" \\\"mixed.txt\\\" \\\"--offset\\\" \\\"11239\\\" \\\"--source\\\" \\\"46b85bac9811733299542453c808777d3c375b11c6f4817cbc8032747bce16ab\\\"] = stdout=\\\"error: cursor out of range — offset 11239 is outside the 11238 selected byt… (341 bytes)\",1\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/spill,\"planning_reader_commands_test.go:96: [\\\"read\\\" \\\"/dev/shm/IWND7Z/t/bench-cmd-home-4247697596/responses/none/primary/1793400900665359220-c4b656308f3c9620.out\\\" \\\"--artifact\\\" \\\"--offset\\\" \\\"5201\\\" \\\"--source\\\" \\\"b48f8496e09065051f4598387fa3b82567a4c327d1… (558 bytes)\",1\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/unicode.txt,\"planning_reader_commands_test.go:96: [\\\"read\\\" \\\"unicode.txt\\\" \\\"--offset\\\" \\\"9601\\\" \\\"--source\\\" \\\"f5e77d45e198b9c37c8593b20cff4af9677890e880235a91536703849c4f81db\\\"] = stdout=\\\"error: cursor out of range — offset 9601 is outside the 9600 selected byte… (342 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ER-C1-r1-ER-C1-successor-red",
          "performer": "ft125-repair-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "c64671b495882fb808603035026f8735babcf984",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft125-repair-01",
            "digest": "sha256:33137a0238df5ce7c1c6237d98b1e397b3b2b74a4bc43026109cc1e6edc8f6ba",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/responsebound/read.go,swap,failed,10,TestExactReaderCLIReconstruction/empty.txt,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestExactReaderCLIReconstruction,passed,12\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,3773\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,404,12\nfailures[10]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/empty.txt,\"planning_reader_commands_test.go:112: [\\\"read\\\" \\\"empty.txt\\\" \\\"--offset\\\" \\\"1\\\" \\\"--source\\\" \\\"f00566d7c32be13f59bf8cc354e7e800c1f10b752a2f02be4819a86045f546e3\\\"] = stdout=\\\"error: cursor out of range — offset 1 is outside the 0 selected bytes; read ag… (330 bytes)\",1\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/exported_evidence,\"planning_reader_commands_test.go:112: [\\\"read\\\" \\\"/dev/shm/ZBPITE/t/TestExactReaderCLIReconstruction2238966278/002/evidence/source-1.txt\\\" \\\"--artifact\\\" \\\"--offset\\\" \\\"9801\\\" \\\"--source\\\" \\\"98988842053f9b987bdb62909ba693f24e5ba55917cfe29d7210f15ce28923… (517 bytes)\",1\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/from_a_subdirectory,\"planning_reader_commands_test.go:112: [\\\"read\\\" \\\"deep.txt\\\" \\\"--offset\\\" \\\"8401\\\" \\\"--source\\\" \\\"878ebc0f1b56930cc5bdbca33b3d6466eb62414acfcf5f291d8149bd1b3f299b\\\"] = stdout=\\\"error: cursor out of range — offset 8401 is outside the 8400 selected bytes;… (337 bytes)\",1\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/invalid.bin,\"planning_reader_commands_test.go:112: [\\\"read\\\" \\\"invalid.bin\\\" \\\"--offset\\\" \\\"6003\\\" \\\"--source\\\" \\\"a1ab9c4cc20654262bcc019f9e95853849ec923fa3d734b867bd196cf0449250\\\"] = stdout=\\\"error: cursor out of range — offset 6003 is outside the 6002 selected byt… (343 bytes)\",1\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/mixed.txt,\"planning_reader_commands_test.go:112: [\\\"read\\\" \\\"mixed.txt\\\" \\\"--offset\\\" \\\"11239\\\" \\\"--source\\\" \\\"d82efeadf5a5434784a9e9004c7adc9d86a54696a7b573aae694a96a314a813a\\\"] = stdout=\\\"error: cursor out of range — offset 11239 is outside the 11238 selected by… (342 bytes)\",1\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/nested/deep.txt,\"planning_reader_commands_test.go:112: [\\\"read\\\" \\\"nested/deep.txt\\\" \\\"--offset\\\" \\\"8401\\\" \\\"--source\\\" \\\"878ebc0f1b56930cc5bdbca33b3d6466eb62414acfcf5f291d8149bd1b3f299b\\\"] = stdout=\\\"error: cursor out of range — offset 8401 is outside the 8400 selected… (351 bytes)\",1\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/parent_from_a_subdirectory,\"planning_reader_commands_test.go:112: [\\\"read\\\" \\\"../mixed.txt\\\" \\\"--offset\\\" \\\"11239\\\" \\\"--source\\\" \\\"d82efeadf5a5434784a9e9004c7adc9d86a54696a7b573aae694a96a314a813a\\\"] = stdout=\\\"error: cursor out of range — offset 11239 is outside the 11238 selected… (348 bytes)\",1\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/spaced_name.txt,\"planning_reader_commands_test.go:112: [\\\"read\\\" \\\"spaced name.txt\\\" \\\"--offset\\\" \\\"8001\\\" \\\"--source\\\" \\\"3071850cc1efb28437772059bcf3947c7e6bed4424d5ff63f3b854d311bfdd71\\\"] = stdout=\\\"error: cursor out of range — offset 8001 is outside the 8000 selected… (353 bytes)\",1\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/spill,\"planning_reader_commands_test.go:112: [\\\"read\\\" \\\"/dev/shm/ZBPITE/t/bench-cmd-home-4247622525/responses/none/primary/1793403881356566172-700fdd1beb5c3b4b.out\\\" \\\"--artifact\\\" \\\"--offset\\\" \\\"5201\\\" \\\"--source\\\" \\\"3a45dc0504dd049ca5bd7056744e3f2e606300624… (559 bytes)\",1\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/unicode.txt,\"planning_reader_commands_test.go:112: [\\\"read\\\" \\\"unicode.txt\\\" \\\"--offset\\\" \\\"9601\\\" \\\"--source\\\" \\\"ea90ae20f56cc6bbadd37f8006b456e8502a51bb1a779ecb382c47b1fca20ba8\\\"] = stdout=\\\"error: cursor out of range — offset 9601 is outside the 9600 selected byt… (343 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "8b89063e28fc3b5ede6b0d6250d23089eee0e0c9",
          "requirement": "ER-C1-successor-red",
          "command": "bench test --package ./cmd/bench --run TestExactReaderCLIReconstruction",
          "exit_code": 0,
          "probe": {
            "mutation": "skip one source byte in continuation successor",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft125-repair-01",
              "digest": "sha256:33137a0238df5ce7c1c6237d98b1e397b3b2b74a4bc43026109cc1e6edc8f6ba",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/responsebound/read.go,swap,failed,10,TestExactReaderCLIReconstruction/empty.txt,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,TestExactReaderCLIReconstruction,passed,12\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,3773\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,404,12\nfailures[10]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/empty.txt,\"planning_reader_commands_test.go:112: [\\\"read\\\" \\\"empty.txt\\\" \\\"--offset\\\" \\\"1\\\" \\\"--source\\\" \\\"f00566d7c32be13f59bf8cc354e7e800c1f10b752a2f02be4819a86045f546e3\\\"] = stdout=\\\"error: cursor out of range — offset 1 is outside the 0 selected bytes; read ag… (330 bytes)\",1\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/exported_evidence,\"planning_reader_commands_test.go:112: [\\\"read\\\" \\\"/dev/shm/ZBPITE/t/TestExactReaderCLIReconstruction2238966278/002/evidence/source-1.txt\\\" \\\"--artifact\\\" \\\"--offset\\\" \\\"9801\\\" \\\"--source\\\" \\\"98988842053f9b987bdb62909ba693f24e5ba55917cfe29d7210f15ce28923… (517 bytes)\",1\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/from_a_subdirectory,\"planning_reader_commands_test.go:112: [\\\"read\\\" \\\"deep.txt\\\" \\\"--offset\\\" \\\"8401\\\" \\\"--source\\\" \\\"878ebc0f1b56930cc5bdbca33b3d6466eb62414acfcf5f291d8149bd1b3f299b\\\"] = stdout=\\\"error: cursor out of range — offset 8401 is outside the 8400 selected bytes;… (337 bytes)\",1\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/invalid.bin,\"planning_reader_commands_test.go:112: [\\\"read\\\" \\\"invalid.bin\\\" \\\"--offset\\\" \\\"6003\\\" \\\"--source\\\" \\\"a1ab9c4cc20654262bcc019f9e95853849ec923fa3d734b867bd196cf0449250\\\"] = stdout=\\\"error: cursor out of range — offset 6003 is outside the 6002 selected byt… (343 bytes)\",1\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/mixed.txt,\"planning_reader_commands_test.go:112: [\\\"read\\\" \\\"mixed.txt\\\" \\\"--offset\\\" \\\"11239\\\" \\\"--source\\\" \\\"d82efeadf5a5434784a9e9004c7adc9d86a54696a7b573aae694a96a314a813a\\\"] = stdout=\\\"error: cursor out of range — offset 11239 is outside the 11238 selected by… (342 bytes)\",1\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/nested/deep.txt,\"planning_reader_commands_test.go:112: [\\\"read\\\" \\\"nested/deep.txt\\\" \\\"--offset\\\" \\\"8401\\\" \\\"--source\\\" \\\"878ebc0f1b56930cc5bdbca33b3d6466eb62414acfcf5f291d8149bd1b3f299b\\\"] = stdout=\\\"error: cursor out of range — offset 8401 is outside the 8400 selected… (351 bytes)\",1\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/parent_from_a_subdirectory,\"planning_reader_commands_test.go:112: [\\\"read\\\" \\\"../mixed.txt\\\" \\\"--offset\\\" \\\"11239\\\" \\\"--source\\\" \\\"d82efeadf5a5434784a9e9004c7adc9d86a54696a7b573aae694a96a314a813a\\\"] = stdout=\\\"error: cursor out of range — offset 11239 is outside the 11238 selected… (348 bytes)\",1\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/spaced_name.txt,\"planning_reader_commands_test.go:112: [\\\"read\\\" \\\"spaced name.txt\\\" \\\"--offset\\\" \\\"8001\\\" \\\"--source\\\" \\\"3071850cc1efb28437772059bcf3947c7e6bed4424d5ff63f3b854d311bfdd71\\\"] = stdout=\\\"error: cursor out of range — offset 8001 is outside the 8000 selected… (353 bytes)\",1\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/spill,\"planning_reader_commands_test.go:112: [\\\"read\\\" \\\"/dev/shm/ZBPITE/t/bench-cmd-home-4247622525/responses/none/primary/1793403881356566172-700fdd1beb5c3b4b.out\\\" \\\"--artifact\\\" \\\"--offset\\\" \\\"5201\\\" \\\"--source\\\" \\\"3a45dc0504dd049ca5bd7056744e3f2e606300624… (559 bytes)\",1\n  github.com/gibbonmi/bench/cmd/bench,TestExactReaderCLIReconstruction/unicode.txt,\"planning_reader_commands_test.go:112: [\\\"read\\\" \\\"unicode.txt\\\" \\\"--offset\\\" \\\"9601\\\" \\\"--source\\\" \\\"ea90ae20f56cc6bbadd37f8006b456e8502a51bb1a779ecb382c47b1fca20ba8\\\"] = stdout=\\\"error: cursor out of range — offset 9601 is outside the 9600 selected byt… (343 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "ER-C1-standards-1",
          "performer": "ft125-review-c1-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "2edc7d7c59d6ac67660c1aa947cecf41b3c9c39e",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft125-review-c1-standards",
            "digest": "sha256:5553046294979116042ddedbbd0b0b67605d6ff0a0bb72fa40c8e12ddcd96e34",
            "excerpt": "Standards axis ER-C1 c02a0324c..97ef9ee0e: 11 findings; worst F1.\nF1 internal/tickets/registry_data.go:21-23 placeholder internal/read entry satisfies a check that derives internal/\u003cquery>; ask-user; conf 9\nF2 internal/responsebound/read.go:168 page-budget predicate re-derives Owner.overBound; auto-fix; conf 7\nF3 internal/bounds/open.go:28-29 second regular-file classifier copies gradeBytes; auto-fix; conf 8\nF4 cmd/bench/planning_reads.go:73 stand-in commandDefinition repeats the registered scope; ask-user; conf 6\nF5 internal/responsebound/read.go:92, file.go:122-123 two decimal validators; auto-fix; conf 6\nF6 internal/responsebound/read.go:163 re-derives the TOON help header to detect axi fallback; ask-user; conf 5\nF7 internal/conformance/axi_query_registry_test.go:22 bite tests not moved as the spec says; auto-fix; conf 7\nF8 cmd/bench/planning_reads.go:80 another containment derivation; no-op; conf 4\nF9 cmd/bench/planning_reads.go:90 literal \"read\" repeats the verb name; auto-fix; conf 4\nF10 cmd/bench/planning_reader_commands_test.go:33,35,134 discarded Rows error; auto-fix; conf 4\nF11 internal/responsebound/read_test.go:18 fixture embeds the build worktree label; auto-fix; conf 3\n"
          },
          "axis": "Standards",
          "base": "c02a0324c839ab93f360781d06da7e8b7c179179",
          "tip": "97ef9ee0e680e9cc4de6f1709b09664dd6cafd0e",
          "finding_ids": [
            "F1",
            "F2",
            "F3",
            "F4",
            "F5",
            "F6",
            "F7",
            "F8",
            "F9",
            "F10",
            "F11"
          ],
          "supersedes": []
        },
        {
          "id": "ER-C1-spec-1",
          "performer": "ft125-review-c1-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "2edc7d7c59d6ac67660c1aa947cecf41b3c9c39e",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft125-review-c1-spec",
            "digest": "sha256:afed380b7b1a92b7f5ac842667ead9eeda750e212f63ce0e3ffee4c5f76b24f4",
            "excerpt": "Spec axis ER-C1 c02a0324c..97ef9ee0e: 4 findings; worst S1.\nS1 cmd/bench/planning_reads.go:90 with read.go:106,211 and main.go:95 a bench read --in continuation omits --in and refuses as changed source; spec grammar has no --in on read; ask-user; conf 8\nS2 cmd/bench/planning_reads.go:79,90 continuation carries the root-relative path, so a subdirectory caller reads the wrong file; auto-fix; conf 8\nS3 internal/tickets/registry_data.go:23 internal/read entry names no package; ask-user; conf 8\nS4 spec Placement lines 338-352 and tickets/1.md:22-23 not amended for the bite tests that stayed in place; auto-fix; conf 7\nDivergences: bite tests in place accepted (IEA seam citations); SKILL.md join accepted; open choices accepted except tree-scoped read (S1).\n"
          },
          "axis": "Spec",
          "base": "c02a0324c839ab93f360781d06da7e8b7c179179",
          "tip": "97ef9ee0e680e9cc4de6f1709b09664dd6cafd0e",
          "finding_ids": [
            "S1",
            "S2",
            "S3",
            "S4"
          ],
          "supersedes": []
        },
        {
          "id": "ER-C1-coverage-1",
          "performer": "ft125-review-c1-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "2edc7d7c59d6ac67660c1aa947cecf41b3c9c39e",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft125-review-c1-coverage",
            "digest": "sha256:bdc15bbee649b0383badabbdbd48fb309055ff8de11afb5ceeb2e321e5ed1ac0",
            "excerpt": "Coverage axis ER-C1 c02a0324c..97ef9ee0e: 7 findings; worst CV1.\nCV1 cmd/bench/planning_reads.go:79,90 continuation from a subdirectory reads nested/\u003cname>; probe bit, 3 failures; auto-fix; conf 9\nCV2 internal/responsebound/file.go:67 oversized repository file refusal untested; probe silent; auto-fix; conf 8\nCV3 internal/responsebound/read.go:213 selection binding in the token untested; probe silent; auto-fix; conf 8\nCV4 internal/responsebound/file.go:101 Source.Unchanged untested; probe silent; auto-fix; conf 7\nCV5 internal/responsebound/read.go:163 help fallback guard untested (control byte in path); probe silent; auto-fix; conf 8\nCV6 internal/responsebound/read.go:182 empty nonterminal guard untested; tiny window splits a rune; probe silent; auto-fix; conf 7\nCV7 internal/responsebound/read.go:168, file.go:127 line clause and --lines overflow untested; probes silent; auto-fix; conf 7\n"
          },
          "axis": "Coverage",
          "base": "c02a0324c839ab93f360781d06da7e8b7c179179",
          "tip": "97ef9ee0e680e9cc4de6f1709b09664dd6cafd0e",
          "finding_ids": [
            "CV1",
            "CV2",
            "CV3",
            "CV4",
            "CV5",
            "CV6",
            "CV7"
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
      "from": "sha256:31861a49dc8cf784d33f054307f9a9c5d708ffb298f803b6b748e5518d78c1f7",
      "to": "sha256:3ac7ce83999acc3cf7e20f551f4a8a52889192313d46706167b7c957a273b330",
      "chunk_ids": {
        "ER-C1": [
          "ER-C1"
        ]
      }
    }
  ]
}
```

## ER-C1 author witness reds

The ticket 1 author `ft125-author-01` (opus / medium) reported these witness reds at ticket commit `97ef9ee0e`.
Each probe answered `bit`, and the author restored each subject.

| subject | mutation | failing witness | row |
| --- | --- | --- | --- |
| `internal/responsebound/read.go` | drop the tree frame from the page budget | `TestExactPageProtocol/ascii` | ER20 |
| `internal/responsebound/read.go` | disable the rune-boundary trim | `TestExactPageProtocol/unicode` | ER21 |
| `internal/responsebound/read.go` | drop the source-token check | `TestExactReaderCLIReconstruction/changed_source` | ER17 |
| `internal/responsebound/read.go` | render invalid UTF-8 as a JSON string | `TestExactPageProtocol/invalid_utf-8` | ER24 |
| `internal/responsebound/read.go` | drop `--source` from the continuation | `TestExactPageProtocol/ascii` | ER18 |
| `internal/responsebound/read.go` | accept an offset equal to the size | `TestExactPageProtocol/offset_at_the_end` | ER28 |
| `internal/bounds/open.go` | remove the `RefuseLinks` call | `TestExactFileStates/linked_parent` | ER37 |
| `cmd/bench/main.go` | make the read row AXI-exempt | `TestAXIRegistryBindsEachRealCommandEnvelope` | ER46 |
| `cmd/bench/main.go` | change the read help suffix | `TestHelpInventoryIsComplete` | ER46 |

## ER-C1 review round 1

The frozen pair is `c02a0324c..97ef9ee0e`. Each axis ran on opus / high.
The round has 22 raw findings and 16 repair targets after duplicates merge.
The reviewer decided F1 and S3 on 2026-10-30: the check takes each query's grammar package from the routing table.

### Standards

Count: 11. Worst issue: F1.

- F1 (auto-fix, reviewer decision, confidence 9): `internal/tickets/registry_data.go:21-23` adds a placeholder `internal/read` entry. Repair target R2.
- F2 (auto-fix, confidence 7): `internal/responsebound/read.go:168` derives the page budget a second time beside `Owner.overBound`. Repair target R3.
- F3 (auto-fix, confidence 8): `internal/bounds/open.go:28-29` copies the regular-file classification of `gradeBytes`. Repair target R4.
- F4 (auto-fix, confidence 6): `cmd/bench/planning_reads.go:73` builds a stand-in definition that repeats the registered scope. Repair target R5.
- F5 (auto-fix, confidence 6): `internal/responsebound/read.go:92` and `file.go:122-123` hold two decimal validators. Repair target R6.
- F6 (no-op, confidence 5): `internal/responsebound/read.go:163` detects the axi help fallback by its header. A learning records the axi predicate.
- F7 (no-op, confidence 7): the bite tests stay in place, because a staged spec cites them there. Spec target S4 amends the text.
- F8 (no-op, confidence 4): `cmd/bench/planning_reads.go:80` adds a containment test. A learning records the missing owner.
- F9 (auto-fix, confidence 4): `cmd/bench/planning_reads.go:90` repeats the verb name. Repair target R8.
- F10 (auto-fix, confidence 4): `cmd/bench/planning_reader_commands_test.go:33,35,134` discards a `Rows` error. Repair target R9.
- F11 (auto-fix, confidence 3): `internal/responsebound/read_test.go:18` embeds the build label. Repair target R10.

### Spec

Count: 4. Worst issue: S1.

- S1 (auto-fix, confidence 8): a root `bench read --in` continuation refuses, and the spec grammar gives root `read` no `--in`. Repair target R1.
- S2 (auto-fix, confidence 8): `cmd/bench/planning_reads.go:79,90` makes a continuation from a subdirectory read the wrong path. Repair target R1.
- S3 (auto-fix, reviewer decision, confidence 8): the same defect as F1. Repair target R2.
- S4 (auto-fix, confidence 7): the spec Placement text and ticket 1 still say to move the bite tests. The plan commit before the repair amends them.

### Coverage

Count: 7. Worst issue: CV1.

- CV1 (auto-fix, confidence 9): a probe from a subdirectory bit with 3 failures. Repair target R1.
- CV2 (auto-fix, confidence 8): no test reads an oversized repository file. Repair target R11.
- CV3 (auto-fix, confidence 8): no test reuses a token across selections. Repair target R12.
- CV4 (auto-fix, confidence 7): no test replaces a file during one read. Repair target R13.
- CV5 (auto-fix, confidence 8): no test puts a control byte or a space in a continued path. Repair target R14.
- CV6 (auto-fix, confidence 7): no test makes a window smaller than one rune; such a window must refuse, not split the rune. Repair target R15.
- CV7 (auto-fix, confidence 7): no test grades the 10-line budget or a `--lines` overflow. Repair targets R3 and R6.
