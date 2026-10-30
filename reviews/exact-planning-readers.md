# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/exact-planning-readers/spec.md",
  "plan_digest": "sha256:31861a49dc8cf784d33f054307f9a9c5d708ffb298f803b6b748e5518d78c1f7",
  "implementation_session": "",
  "chunks": [
    {
      "id": "ER-C1",
      "base": "c02a0324c839ab93f360781d06da7e8b7c179179",
      "tip": "97ef9ee0e680e9cc4de6f1709b09664dd6cafd0e",
      "plan_digest": "sha256:31861a49dc8cf784d33f054307f9a9c5d708ffb298f803b6b748e5518d78c1f7",
      "source_digest": "2edc7d7c59d6ac67660c1aa947cecf41b3c9c39e",
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
