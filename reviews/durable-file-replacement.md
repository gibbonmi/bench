# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/durable-file-replacement/spec.md",
  "plan_digest": "sha256:7a4c1e409bfe8a274ca0cdaca77ce869f2efab34752d3692aa7cb5e7e6057754",
  "implementation_session": "",
  "chunks": [
    {
      "id": "D-A",
      "base": "9433df1bc7cc583dd77f237856b07c9cbf29250f",
      "tip": "5217d90168cdf1a39d8dba48483e8be1c9ed6860",
      "plan_digest": "sha256:d500c83005b78119c2c91ab82b12285c5ebe048c3f6f88d845ad83b896fca219",
      "source_digest": "00bb572250b671e1e40ce1ca18587a1a5bd5793e",
      "acceptance_rows": [
        "D01",
        "D02",
        "D03",
        "D04",
        "D05",
        "D06",
        "D07",
        "D08",
        "D09",
        "D10",
        "D11",
        "D12",
        "D13",
        "D14",
        "D15",
        "D19",
        "D32"
      ],
      "verification": [
        {
          "id": "dfr-t1-owner-tests",
          "performer": "claude:dfr_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "00bb572250b671e1e40ce1ca18587a1a5bd5793e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1/owner-tests",
            "digest": "sha256:46940969bafb916898a9c395971a14f377914c4cb93f88baa8ad3ffa3af1b1ec",
            "excerpt": "tree: ft354-durable-build 5217d90168cdf1a39d8dba48483e8be1c9ed6860 clean\nbench test --package ./internal/durablefile\ngithub.com/gibbonmi/bench/internal/durablefile,pass,15,16\nfailures[0] skips[0]\n"
          },
          "requirement": "owner-tests",
          "command": "bench test --package ./internal/durablefile",
          "exit_code": 0
        },
        {
          "id": "dfr-t1-caller-tests",
          "performer": "claude:dfr_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "00bb572250b671e1e40ce1ca18587a1a5bd5793e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1/caller-tests",
            "digest": "sha256:96d25e4dc0403ac10f86ff577f9d66a5c51480c1c46be9fe5d881f3b2ddda4b5",
            "excerpt": "tree: ft354-durable-build 5217d90168cdf1a39d8dba48483e8be1c9ed6860 clean\nbench test --package ./internal/reviewrecord\ngithub.com/gibbonmi/bench/internal/reviewrecord,pass,1565,88\nfailures[0] skips[0]\n"
          },
          "requirement": "caller-tests",
          "command": "bench test --package ./internal/reviewrecord",
          "exit_code": 0
        },
        {
          "id": "dfr-t1-linux-native",
          "performer": "claude:dfr_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "00bb572250b671e1e40ce1ca18587a1a5bd5793e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1/linux-native",
            "digest": "sha256:a91a1a5bb40d06d28e8c30bb2d7a215ccb056d6969403718b6164abb3b51f503",
            "excerpt": "tree: ft354-durable-build 5217d90168cdf1a39d8dba48483e8be1c9ed6860 clean\nhost: Linux 6.6.87.2-microsoft-standard-WSL2 x86_64 (WSL2), GOOS linux/amd64\nbench test --package ./internal/durablefile --run 'TestNativeReplacement'\nrun 1, default kit TMPDIR under /dev/shm, filesystem tmpfs: github.com/gibbonmi/bench/internal/durablefile,pass,2,1; failures[0] skips[0]\nrun 2, caller TMPDIR on a disk directory, filesystem ext2/ext3 magic (ext4): github.com/gibbonmi/bench/internal/durablefile,pass,15,1; failures[0] skips[0]\nTestNativeReplacement asserts the native trace create, mode, write, file-sync, file-close, rename, directory-open, directory-sync, directory-close with no error, then the native Replace bytes and mode.\nmacOS: not run; no macOS host is available.\n"
          },
          "requirement": "linux-native",
          "command": "bench test --package ./internal/durablefile --run 'TestNativeReplacement'",
          "exit_code": 0
        },
        {
          "id": "dfr-t1-root-conformance",
          "performer": "claude:dfr_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "00bb572250b671e1e40ce1ca18587a1a5bd5793e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1/root-conformance",
            "digest": "sha256:19d13478d6a8f932194a26f1bd3162df4cf516fc1450986ad5edf2e742466bf7",
            "excerpt": "tree: ft354-durable-build 5217d90168cdf1a39d8dba48483e8be1c9ed6860 clean\nbench test --check root-conformance\nroot-conformance,conformance,1,0\ngithub.com/gibbonmi/bench/internal/conformance,pass,19207,1\nfailures[0] skips[0]\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "dfr-t1-omit-file-sync",
          "performer": "claude:dfr_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "00bb572250b671e1e40ce1ca18587a1a5bd5793e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1/omit-file-sync",
            "digest": "sha256:02a4f8dc062b2b5aaee8283cd7858f309c5baf45e772e3905da0ae610e30d2d7",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,omit,failed,6,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1395\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,14,16\nfailures[6]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestNativeReplacement,\"platform_test.go:19: host linux/amd64, destination /dev/shm/WEP5ND/t/TestNativeReplacement424110334/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md\",2\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceAbsentFailures,\"replace_test.go:216: Replace with a file-sync failure = <nil>, want a *durablefile.Error\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceBeforeRenameFailures,\"replace_test.go:204: Replace with a file-sync failure = <nil>, want a *durablefile.Error\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceErrorCause,\"replace_test.go:242: Replace with a file-sync failure = <nil>, want a *durablefile.Error\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceOrder,\"replace_test.go:147: trace = [\\\"create\\\" \\\"mode\\\" \\\"write\\\" \\\"file-close\\\" \\\"rename\\\" \\\"directory-open\\\" \\\"directory-sync\\\" \\\"directory-close\\\"], want the prefix [\\\"create\\\" \\\"mode\\\" \\\"write\\\" \\\"file-sync\\\" \\\"file-close\\\" \\\"rename\\\"]\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceTemporaryCleanup,\"replace_test.go:252: Replace with a file-sync failure = <nil>, want a *durablefile.Error\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "owner-omission-proof",
          "command": "bench test --package ./internal/durablefile",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit file sync, omit directory sync, rename before file close, or accept a short write. Each named owner assertion must fail, then pass after restoring the source.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:dfr_t1/omit-file-sync",
              "digest": "sha256:02a4f8dc062b2b5aaee8283cd7858f309c5baf45e772e3905da0ae610e30d2d7",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,omit,failed,6,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1395\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,14,16\nfailures[6]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestNativeReplacement,\"platform_test.go:19: host linux/amd64, destination /dev/shm/WEP5ND/t/TestNativeReplacement424110334/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md\",2\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceAbsentFailures,\"replace_test.go:216: Replace with a file-sync failure = <nil>, want a *durablefile.Error\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceBeforeRenameFailures,\"replace_test.go:204: Replace with a file-sync failure = <nil>, want a *durablefile.Error\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceErrorCause,\"replace_test.go:242: Replace with a file-sync failure = <nil>, want a *durablefile.Error\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceOrder,\"replace_test.go:147: trace = [\\\"create\\\" \\\"mode\\\" \\\"write\\\" \\\"file-close\\\" \\\"rename\\\" \\\"directory-open\\\" \\\"directory-sync\\\" \\\"directory-close\\\"], want the prefix [\\\"create\\\" \\\"mode\\\" \\\"write\\\" \\\"file-sync\\\" \\\"file-close\\\" \\\"rename\\\"]\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceTemporaryCleanup,\"replace_test.go:252: Replace with a file-sync failure = <nil>, want a *durablefile.Error\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dfr-t1-omit-directory-sync",
          "performer": "claude:dfr_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "00bb572250b671e1e40ce1ca18587a1a5bd5793e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1/omit-directory-sync",
            "digest": "sha256:44849f50e798a54a5945da4ffd3ea4f10092fddd0a9aad4adb319e085a7c2a02",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,swap,failed,6,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1579\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,16,16\nfailures[6]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestNativeReplacement,\"platform_test.go:19: host linux/amd64, destination /dev/shm/VPWIAX/t/TestNativeReplacement890424444/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md\",2\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceAfterRenameFailures,\"replace_test.go:225: Replace with a directory-sync failure = <nil>, want a *durablefile.Error\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceAfterRenameState,\"replace_test.go:234: Replace with a directory-sync failure = <nil>, want a *durablefile.Error\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceCleanupFailure,\"replace_test.go:266: Replace = replace /dev/shm/VPWIAX/t/TestReplaceCleanupFailure2897667435/003/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md: directory-close: injected directory-close failure; the rename completed, so new bytes may a… (335 bytes)\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceDirectorySync,\"replace_test.go:160: trace = [\\\"create\\\" \\\"mode\\\" \\\"write\\\" \\\"file-sync\\\" \\\"file-close\\\" \\\"rename\\\" \\\"directory-open\\\" \\\"directory-close\\\"], want [\\\"directory-open\\\" \\\"directory-sync\\\" \\\"directory-close\\\"] after the rename\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceErrorCause,\"replace_test.go:242: Replace with a directory-sync failure = <nil>, want a *durablefile.Error\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "owner-omission-proof",
          "command": "bench test --package ./internal/durablefile",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit file sync, omit directory sync, rename before file close, or accept a short write. Each named owner assertion must fail, then pass after restoring the source.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:dfr_t1/omit-directory-sync",
              "digest": "sha256:44849f50e798a54a5945da4ffd3ea4f10092fddd0a9aad4adb319e085a7c2a02",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,swap,failed,6,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1579\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,16,16\nfailures[6]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestNativeReplacement,\"platform_test.go:19: host linux/amd64, destination /dev/shm/VPWIAX/t/TestNativeReplacement890424444/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md\",2\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceAfterRenameFailures,\"replace_test.go:225: Replace with a directory-sync failure = <nil>, want a *durablefile.Error\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceAfterRenameState,\"replace_test.go:234: Replace with a directory-sync failure = <nil>, want a *durablefile.Error\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceCleanupFailure,\"replace_test.go:266: Replace = replace /dev/shm/VPWIAX/t/TestReplaceCleanupFailure2897667435/003/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md: directory-close: injected directory-close failure; the rename completed, so new bytes may a… (335 bytes)\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceDirectorySync,\"replace_test.go:160: trace = [\\\"create\\\" \\\"mode\\\" \\\"write\\\" \\\"file-sync\\\" \\\"file-close\\\" \\\"rename\\\" \\\"directory-open\\\" \\\"directory-close\\\"], want [\\\"directory-open\\\" \\\"directory-sync\\\" \\\"directory-close\\\"] after the rename\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceErrorCause,\"replace_test.go:242: Replace with a directory-sync failure = <nil>, want a *durablefile.Error\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dfr-t1-rename-before-close",
          "performer": "claude:dfr_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "00bb572250b671e1e40ce1ca18587a1a5bd5793e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1/rename-before-close",
            "digest": "sha256:89e38413b39ff30f3a35cb281ffa4aa1c9cfb40cf2842aae7f10124c45025737",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,swap,failed,5,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1408\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,14,16\nfailures[5]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestNativeReplacement,\"platform_test.go:19: host linux/amd64, destination /dev/shm/GK7IQ3/t/TestNativeReplacement696317319/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md\",2\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceAbsentFailures,\"replace_test.go:218: after a file-close failure the absent destination is present: <nil>\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceBeforeRenameFailures,\"replace_test.go:208: after a file-close failure the destination = \\\"new\\\", want the old \\\"old bytes\\\"\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceDirectorySync,\"replace_test.go:160: trace = [\\\"create\\\" \\\"mode\\\" \\\"write\\\" \\\"file-sync\\\" \\\"rename\\\" \\\"file-close\\\" \\\"directory-open\\\" \\\"directory-sync\\\" \\\"directory-close\\\"], want [\\\"directory-open\\\" \\\"directory-sync\\\" \\\"directory-close\\\"] after the rename\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceOrder,\"replace_test.go:147: trace = [\\\"create\\\" \\\"mode\\\" \\\"write\\\" \\\"file-sync\\\" \\\"rename\\\" \\\"file-close\\\" \\\"directory-open\\\" \\\"directory-sync\\\" \\\"directory-close\\\"], want the prefix [\\\"create\\\" \\\"mode\\\" \\\"write\\\" \\\"file-sync\\\" \\\"file-close\\\" \\\"rename\\\"]\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "owner-omission-proof",
          "command": "bench test --package ./internal/durablefile",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit file sync, omit directory sync, rename before file close, or accept a short write. Each named owner assertion must fail, then pass after restoring the source.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:dfr_t1/rename-before-close",
              "digest": "sha256:89e38413b39ff30f3a35cb281ffa4aa1c9cfb40cf2842aae7f10124c45025737",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,swap,failed,5,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1408\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,14,16\nfailures[5]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestNativeReplacement,\"platform_test.go:19: host linux/amd64, destination /dev/shm/GK7IQ3/t/TestNativeReplacement696317319/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md\",2\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceAbsentFailures,\"replace_test.go:218: after a file-close failure the absent destination is present: <nil>\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceBeforeRenameFailures,\"replace_test.go:208: after a file-close failure the destination = \\\"new\\\", want the old \\\"old bytes\\\"\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceDirectorySync,\"replace_test.go:160: trace = [\\\"create\\\" \\\"mode\\\" \\\"write\\\" \\\"file-sync\\\" \\\"rename\\\" \\\"file-close\\\" \\\"directory-open\\\" \\\"directory-sync\\\" \\\"directory-close\\\"], want [\\\"directory-open\\\" \\\"directory-sync\\\" \\\"directory-close\\\"] after the rename\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceOrder,\"replace_test.go:147: trace = [\\\"create\\\" \\\"mode\\\" \\\"write\\\" \\\"file-sync\\\" \\\"rename\\\" \\\"file-close\\\" \\\"directory-open\\\" \\\"directory-sync\\\" \\\"directory-close\\\"], want the prefix [\\\"create\\\" \\\"mode\\\" \\\"write\\\" \\\"file-sync\\\" \\\"file-close\\\" \\\"rename\\\"]\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dfr-t1-accept-short-write",
          "performer": "claude:dfr_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "00bb572250b671e1e40ce1ca18587a1a5bd5793e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1/accept-short-write",
            "digest": "sha256:5d3e873d1fb0ae6771d5e145779a153e219f7b1cf1641bdf9897283df746f367",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,334\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,14,16\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceShortWrite,\"replace_test.go:173: Replace after a short write = <nil>, want io.ErrShortWrite\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "owner-omission-proof",
          "command": "bench test --package ./internal/durablefile",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit file sync, omit directory sync, rename before file close, or accept a short write. Each named owner assertion must fail, then pass after restoring the source.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:dfr_t1/accept-short-write",
              "digest": "sha256:5d3e873d1fb0ae6771d5e145779a153e219f7b1cf1641bdf9897283df746f367",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,334\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,14,16\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceShortWrite,\"replace_test.go:173: Replace after a short write = <nil>, want io.ErrShortWrite\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dfr-t1-corrupt-payload",
          "performer": "claude:dfr_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "00bb572250b671e1e40ce1ca18587a1a5bd5793e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1/corrupt-payload",
            "digest": "sha256:decf17da8ac8d00b18a390ea165be6e497a72340a7447ab63054ec2238154f3e",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,swap,failed,4,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,988\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,9,16\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestNativeReplacement,\"platform_test.go:19: host linux/amd64, destination /dev/shm/LT2SBW/t/TestNativeReplacement2727354860/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md\",2\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceAfterRenameState,\"replace_test.go:236: after a directory-sync failure the destination = \\\"new p\\\", want the new \\\"new payload\\\"\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceBytes,\"replace_test.go:53: destination = \\\"new\\\\\\\\x00payload\\\\\\\\n\\\\\\\\twith a\\\", want \\\"new\\\\\\\\x00payload\\\\\\\\n\\\\\\\\twith a NUL, a tab, and ü\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceVisibility,\"replace_test.go:117: a concurrent reader observed 524288 bytes that are neither complete payload\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "owner-preservation",
          "command": "bench test --package ./internal/durablefile",
          "exit_code": 0,
          "probe": {
            "mutation": "Corrupt supplied payload bytes, requested mode, temporary cleanup, or wrapped error cause. The corresponding independent owner assertion must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:dfr_t1/corrupt-payload",
              "digest": "sha256:decf17da8ac8d00b18a390ea165be6e497a72340a7447ab63054ec2238154f3e",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,swap,failed,4,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,988\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,9,16\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestNativeReplacement,\"platform_test.go:19: host linux/amd64, destination /dev/shm/LT2SBW/t/TestNativeReplacement2727354860/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md\",2\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceAfterRenameState,\"replace_test.go:236: after a directory-sync failure the destination = \\\"new p\\\", want the new \\\"new payload\\\"\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceBytes,\"replace_test.go:53: destination = \\\"new\\\\\\\\x00payload\\\\\\\\n\\\\\\\\twith a\\\", want \\\"new\\\\\\\\x00payload\\\\\\\\n\\\\\\\\twith a NUL, a tab, and ü\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceVisibility,\"replace_test.go:117: a concurrent reader observed 524288 bytes that are neither complete payload\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dfr-t1-corrupt-mode",
          "performer": "claude:dfr_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "00bb572250b671e1e40ce1ca18587a1a5bd5793e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1/corrupt-mode",
            "digest": "sha256:b7c216b9654c4cccea8d6e88d878f7824d4d42b6dd5fa86702d5926dd4910c53",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,326\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,15,16\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceMode,\"replace_test.go:82: destination permission bits = -rw-------, want -rwxr-xr-x\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "owner-preservation",
          "command": "bench test --package ./internal/durablefile",
          "exit_code": 0,
          "probe": {
            "mutation": "Corrupt supplied payload bytes, requested mode, temporary cleanup, or wrapped error cause. The corresponding independent owner assertion must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:dfr_t1/corrupt-mode",
              "digest": "sha256:b7c216b9654c4cccea8d6e88d878f7824d4d42b6dd5fa86702d5926dd4910c53",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,326\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,15,16\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceMode,\"replace_test.go:82: destination permission bits = -rw-------, want -rwxr-xr-x\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dfr-t1-corrupt-cleanup",
          "performer": "claude:dfr_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "00bb572250b671e1e40ce1ca18587a1a5bd5793e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1/corrupt-cleanup",
            "digest": "sha256:d8d61858f13d1b8da08bbc630ceb018f5507617306b63558c9d082e615b20bdc",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,700\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,16,16\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceCleanupFailure,\"replace_test.go:266: Replace = replace /dev/shm/OQAQYG/t/TestReplaceCleanupFailure3229413224/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md: write: injected write failure, want both the write cause and the remove cleanup cause\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceTemporaryCleanup,\"replace_test.go:254: after a mode failure the directory holds [\\\".record 'one' *.md.891042990\\\"] beside the destination\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "owner-preservation",
          "command": "bench test --package ./internal/durablefile",
          "exit_code": 0,
          "probe": {
            "mutation": "Corrupt supplied payload bytes, requested mode, temporary cleanup, or wrapped error cause. The corresponding independent owner assertion must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:dfr_t1/corrupt-cleanup",
              "digest": "sha256:d8d61858f13d1b8da08bbc630ceb018f5507617306b63558c9d082e615b20bdc",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,700\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,16,16\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceCleanupFailure,\"replace_test.go:266: Replace = replace /dev/shm/OQAQYG/t/TestReplaceCleanupFailure3229413224/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md: write: injected write failure, want both the write cause and the remove cleanup cause\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceTemporaryCleanup,\"replace_test.go:254: after a mode failure the directory holds [\\\".record 'one' *.md.891042990\\\"] beside the destination\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dfr-t1-corrupt-cause",
          "performer": "claude:dfr_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "00bb572250b671e1e40ce1ca18587a1a5bd5793e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1/corrupt-cause",
            "digest": "sha256:869ab07133ed29c214a5d40c7d072dd174923b2bfa1ec1b58ad6b929fa939f05",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1038\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,15,16\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceCleanupFailure,\"replace_test.go:266: Replace = replace /dev/shm/VNDNX2/t/TestReplaceCleanupFailure3019778760/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md: write: injected write failure\",2\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceErrorCause,\"replace_test.go:244: a create failure = replace /dev/shm/VNDNX2/t/TestReplaceErrorCause3212068268/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md: create: injected create failure, which does not hold the injected cause\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceShortWrite,\"replace_test.go:173: Replace after a short write = replace /dev/shm/VNDNX2/t/TestReplaceShortWrite1216901929/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md: write: short write, want io.ErrShortWrite\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "owner-preservation",
          "command": "bench test --package ./internal/durablefile",
          "exit_code": 0,
          "probe": {
            "mutation": "Corrupt supplied payload bytes, requested mode, temporary cleanup, or wrapped error cause. The corresponding independent owner assertion must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:dfr_t1/corrupt-cause",
              "digest": "sha256:869ab07133ed29c214a5d40c7d072dd174923b2bfa1ec1b58ad6b929fa939f05",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1038\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,15,16\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceCleanupFailure,\"replace_test.go:266: Replace = replace /dev/shm/VNDNX2/t/TestReplaceCleanupFailure3019778760/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md: write: injected write failure\",2\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceErrorCause,\"replace_test.go:244: a create failure = replace /dev/shm/VNDNX2/t/TestReplaceErrorCause3212068268/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md: create: injected create failure, which does not hold the injected cause\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceShortWrite,\"replace_test.go:173: Replace after a short write = replace /dev/shm/VNDNX2/t/TestReplaceShortWrite1216901929/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md: write: short write, want io.ErrShortWrite\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dfr-t1-caller-bypass",
          "performer": "claude:dfr_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "00bb572250b671e1e40ce1ca18587a1a5bd5793e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1/caller-bypass",
            "digest": "sha256:aed8fcab14b5066945e1fb718dccd62637375480601655f38670d02f842c0a14",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/write.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,all,passed,88\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,620\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,1521,88\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewReplacement/directory_sync_failure_after_the_rename,\"replacement_test.go:119: write with a failed directory sync = <nil>, want the published directory-sync cause\",1\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewReplacement/file_sync_failure_before_the_rename,\"replacement_test.go:147: write with a failed file sync = <nil>, want the unpublished file-sync cause\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "caller-bypass-proof",
          "command": "bench test --package ./internal/reviewrecord",
          "exit_code": 0,
          "probe": {
            "mutation": "Restore the former review-record writer or bypass the shared replacement operation in the real record transaction. TestReviewReplacement must fail its directory-sync cause assertion, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:dfr_t1/caller-bypass",
              "digest": "sha256:aed8fcab14b5066945e1fb718dccd62637375480601655f38670d02f842c0a14",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/write.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,all,passed,88\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,620\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,1521,88\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewReplacement/directory_sync_failure_after_the_rename,\"replacement_test.go:119: write with a failed directory sync = <nil>, want the published directory-sync cause\",1\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewReplacement/file_sync_failure_before_the_rename,\"replacement_test.go:147: write with a failed file sync = <nil>, want the unpublished file-sync cause\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dfr-t1-caller-payload",
          "performer": "claude:dfr_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "00bb572250b671e1e40ce1ca18587a1a5bd5793e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1/caller-payload",
            "digest": "sha256:648c654edc07077a487ff261a46d35208d67622d73c97977e1cf799d38a4d1ff",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/write.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,all,passed,88\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,644\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,1534,88\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewReplacement/directory_sync_failure_after_the_rename,\"replacement_test.go:125: record = \\\"# Review outcomes\\\\\\\\n\\\\\\\\n```bench-review-record\\\\\\\\n{\\\\\\\\n  \\\\\\\\\\\"version\\\\\\\\\\\": 1,\\\\\\\\n  \\\\\\\\\\\"spec\\\\\\\\\\\": \\\\\\\\\\\"specs/example/spec.md\\\\\\\\\\\",\\\\\\\\n  \\\\\\\\\\\"plan_digest\\\\\\\\\\\": \\\\\\\\\\\"sha256:plan\\\\\\\\\\\",\\\\\\\\n  \\\\\\\\\\\"implementation_session\\\\\\\\\\\": \\\\\\\\\\\"author\\\\\\\\\\\",\\\\\\\\n  \\\\\\\\\\\"chunks\\\\\\\\\\\": null… (1020 bytes)\",3\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "caller-preservation",
          "command": "bench test --package ./internal/reviewrecord",
          "exit_code": 0,
          "probe": {
            "mutation": "Corrupt the review-record rendered payload, 0644 mode, or pre-write refusal policy. The real caller assertion must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:dfr_t1/caller-payload",
              "digest": "sha256:648c654edc07077a487ff261a46d35208d67622d73c97977e1cf799d38a4d1ff",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/write.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,all,passed,88\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,644\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,1534,88\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewReplacement/directory_sync_failure_after_the_rename,\"replacement_test.go:125: record = \\\"# Review outcomes\\\\\\\\n\\\\\\\\n```bench-review-record\\\\\\\\n{\\\\\\\\n  \\\\\\\\\\\"version\\\\\\\\\\\": 1,\\\\\\\\n  \\\\\\\\\\\"spec\\\\\\\\\\\": \\\\\\\\\\\"specs/example/spec.md\\\\\\\\\\\",\\\\\\\\n  \\\\\\\\\\\"plan_digest\\\\\\\\\\\": \\\\\\\\\\\"sha256:plan\\\\\\\\\\\",\\\\\\\\n  \\\\\\\\\\\"implementation_session\\\\\\\\\\\": \\\\\\\\\\\"author\\\\\\\\\\\",\\\\\\\\n  \\\\\\\\\\\"chunks\\\\\\\\\\\": null… (1020 bytes)\",3\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dfr-t1-caller-mode",
          "performer": "claude:dfr_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "00bb572250b671e1e40ce1ca18587a1a5bd5793e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1/caller-mode",
            "digest": "sha256:4607f3f24c6d2378c7f3bbeae73a77e56ecc78a9f1bba4bc96c4e439dd6c6589",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/write.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,all,passed,88\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,370\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,1962,88\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewReplacement/directory_sync_failure_after_the_rename,\"replacement_test.go:132: record permission bits = -rw-------, want 0644\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "caller-preservation",
          "command": "bench test --package ./internal/reviewrecord",
          "exit_code": 0,
          "probe": {
            "mutation": "Corrupt the review-record rendered payload, 0644 mode, or pre-write refusal policy. The real caller assertion must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:dfr_t1/caller-mode",
              "digest": "sha256:4607f3f24c6d2378c7f3bbeae73a77e56ecc78a9f1bba4bc96c4e439dd6c6589",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/write.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,all,passed,88\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,370\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,1962,88\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewReplacement/directory_sync_failure_after_the_rename,\"replacement_test.go:132: record permission bits = -rw-------, want 0644\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dfr-t1-caller-refusal",
          "performer": "claude:dfr_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "00bb572250b671e1e40ce1ca18587a1a5bd5793e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1/caller-refusal",
            "digest": "sha256:49cf94705e5308208fdb64a90041dd05fadc9392f9f82250e9cd5fdc9d9c9f75",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/write.go,omit,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,all,passed,88\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,423\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,1526,88\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewReplacement/refusal_before_the_replacement,\"replacement_test.go:171: write = <nil> after 1 replacements, want the refusal \\\"invalid record spec identity\\\" before any replacement\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "caller-preservation",
          "command": "bench test --package ./internal/reviewrecord",
          "exit_code": 0,
          "probe": {
            "mutation": "Corrupt the review-record rendered payload, 0644 mode, or pre-write refusal policy. The real caller assertion must fail, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:dfr_t1/caller-refusal",
              "digest": "sha256:49cf94705e5308208fdb64a90041dd05fadc9392f9f82250e9cd5fdc9d9c9f75",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/write.go,omit,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,all,passed,88\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,423\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,1526,88\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewReplacement/refusal_before_the_replacement,\"replacement_test.go:171: write = <nil> after 1 replacements, want the refusal \\\"invalid record spec identity\\\" before any replacement\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "dfr-da-standards-1",
          "performer": "claude:dfr_rev_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "00bb572250b671e1e40ce1ca18587a1a5bd5793e",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:dfr_rev_standards",
            "digest": "sha256:48b32d739d8dd530828b1391c6fb2a3e1586e0d8b644b325847961a472f38d12",
            "excerpt": "Standards axis, chunk D-A (9433df1b..5217d901), opus/high. 2 findings; worst S1.\nS1 blocking ask-user conf 6: internal/reviewrecord/replacement_test.go:15-54 (failingFile, failingDirectory, failingSync) duplicates the fault-injection harness of internal/durablefile/recorder_test.go:42-83,107-132; AGENTS.md keeps fixture harnesses single-sourced. Fix: export the recorder as a test-support package under internal/durablefile/.\nS2 non-blocking auto-fix conf 5: internal/durablefile/platform_test.go:15-16 comment restates the review-record obligation that ticket line 60 owns. Fix: delete the sentence.\nRefuted: ordered sequence exists once (replace.go:95-134); error wrapping keeps errors.Is/As; stage literals are the spec contract; comment register matches.\n"
          },
          "axis": "Standards",
          "base": "9433df1bc7cc583dd77f237856b07c9cbf29250f",
          "tip": "5217d90168cdf1a39d8dba48483e8be1c9ed6860",
          "finding_ids": [
            "S1",
            "S2"
          ],
          "supersedes": []
        },
        {
          "id": "dfr-da-spec-1",
          "performer": "claude:dfr_rev_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "00bb572250b671e1e40ce1ca18587a1a5bd5793e",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:dfr_rev_spec",
            "digest": "sha256:80683704f7549323802f126191be677d8e8c28a3fa7042719c99007edd0457b3",
            "excerpt": "Spec axis, chunk D-A (9433df1b..5217d901), opus/high. 2 findings; worst P1.\nAll 17 rows D01-D15, D19, D32 audited: each cited test exists and asserts its row; implementation decisions met (replace.go:16-139; write.go:107,121-126; assignment.go:85).\nP1 blocking ask-user conf 10: D32 macos-native requirement (spec:133,186; ticket:25) has no evidence; review record line 84 records Linux (WSL2 tmpfs and ext4) only.\nP2 non-blocking no-op conf 4: assignment.go:85 passes a fixed replacer, so no test injects a failure into RecordAssignment; spec line 96 accepts the shared wrapper.\n"
          },
          "axis": "Spec",
          "base": "9433df1bc7cc583dd77f237856b07c9cbf29250f",
          "tip": "5217d90168cdf1a39d8dba48483e8be1c9ed6860",
          "finding_ids": [
            "P1",
            "P2"
          ],
          "supersedes": []
        },
        {
          "id": "dfr-da-coverage-1",
          "performer": "claude:dfr_rev_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "00bb572250b671e1e40ce1ca18587a1a5bd5793e",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:dfr_rev_coverage",
            "digest": "sha256:51ca19ef931857fc0d970ceb0c5d8aa66288fc0b7af3ff7811900d09bf18c917",
            "excerpt": "Coverage axis, chunk D-A (9433df1b..5217d901), opus/high. 3 findings; worst C1. Ran bench test ./internal/durablefile (16 pass) and ./internal/reviewrecord (88 pass).\nC1 blocking ask-user conf 9: D32 needs native Linux and macOS qualification; no macos-native entry (reviews:84; spec:119,186).\nC2 non-blocking auto-fix conf 7: TestReviewReplacement calls writeWith directly (replacement_test.go:116,144,170; write.go:196-198); no test pins the replacement that the production entry write binds.\nC3 non-blocking auto-fix conf 7: RecordAssignment binds its own replacement (assignment.go:85); no injected failure reaches it through the real transaction. Same fix as C2.\nClean: every stage has a failure witness; edge fixtures present; real transaction driven; recorded probes name plan mutations.\n"
          },
          "axis": "Coverage",
          "base": "9433df1bc7cc583dd77f237856b07c9cbf29250f",
          "tip": "5217d90168cdf1a39d8dba48483e8be1c9ed6860",
          "finding_ids": [
            "C1",
            "C2",
            "C3"
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
      "from": "sha256:d500c83005b78119c2c91ab82b12285c5ebe048c3f6f88d845ad83b896fca219",
      "to": "sha256:7a4c1e409bfe8a274ca0cdaca77ce869f2efab34752d3692aa7cb5e7e6057754",
      "chunk_ids": {
        "D-A": [
          "D-A"
        ]
      }
    }
  ]
}
```

## Standards

Chunk D-A has 2 findings. The worst issue is S1.

- S1, auto-fix, confidence 6: `internal/reviewrecord/replacement_test.go:15-54` repeats the fault-injection harness of `internal/durablefile/recorder_test.go`. Move the harness to one test-support package under `internal/durablefile/`, and let both test files use it.
- S2, auto-fix, confidence 5: `internal/durablefile/platform_test.go:15-16` states the evidence obligation that ticket line 60 owns. Delete that sentence.

## Spec

Chunk D-A has 2 findings. The worst issue is P1.

- P1, ask-user, confidence 10: row D32 needs native macOS evidence, and this host has no macOS. The reviewer decides the route.
- P2, no-op, confidence 4: the C3 repair closes the same gap in `internal/reviewrecord/assignment.go:85`.

## Coverage

Chunk D-A has 3 findings. The worst issue is C1.

- C1, ask-user, confidence 9: this finding is the same as P1.
- C2, auto-fix, confidence 7: `TestReviewReplacement` calls `writeWith` directly, so no test examines the replacement that the `write` entry binds (`internal/reviewrecord/write.go:196-198`). Inject the failure through the real `write` entry.
- C3, auto-fix, confidence 7: `RecordAssignment` binds its own replacement (`internal/reviewrecord/assignment.go:85`), and no injected failure reaches it. Inject the failure through the real `RecordAssignment` transaction. C2 and C3 have one repair target.
