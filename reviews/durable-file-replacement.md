# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/durable-file-replacement/spec.md",
  "plan_digest": "sha256:286fab6eee12fd2fb29be89f7acf16d6242dd8e57be66c984d59206f683de6df",
  "implementation_session": "",
  "chunks": [
    {
      "id": "D-A",
      "base": "9433df1bc7cc583dd77f237856b07c9cbf29250f",
      "tip": "61e6ca30bcf699be0d8ef9d3d10ff9014888af82",
      "plan_digest": "sha256:286fab6eee12fd2fb29be89f7acf16d6242dd8e57be66c984d59206f683de6df",
      "source_digest": "b9ea1b2395d2ae721f63c905b289bbde1a4100b3",
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
        },
        {
          "id": "dfr-t1-rep1-owner-tests",
          "performer": "claude:dfr_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "58796930b5bb3c6483e32048fac3bdb733b175cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1_rep1/owner-tests",
            "digest": "sha256:ddbe49da3e9dd64d161b7f00780685b150f37b3e2d4caa3b995ee18d354682bb",
            "excerpt": "tree: ft354-durable-build 5f1ea9dc7b6612d8095d883eca4730bcc8853cd9 clean\nbench test --package ./internal/durablefile\ngithub.com/gibbonmi/bench/internal/durablefile,pass,16,16\nfailures[0] skips[0]\n"
          },
          "requirement": "owner-tests",
          "command": "bench test --package ./internal/durablefile",
          "exit_code": 0
        },
        {
          "id": "dfr-t1-rep1-caller-tests",
          "performer": "claude:dfr_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "58796930b5bb3c6483e32048fac3bdb733b175cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1_rep1/caller-tests",
            "digest": "sha256:d35b941b515ca1cfdf19a34e72107148e3be44f2247cb5a2e58c6a66fba88179",
            "excerpt": "tree: ft354-durable-build 5f1ea9dc7b6612d8095d883eca4730bcc8853cd9 clean\nbench test --package ./internal/reviewrecord\ngithub.com/gibbonmi/bench/internal/reviewrecord,pass,1641,89\nfailures[0] skips[0]\n"
          },
          "requirement": "caller-tests",
          "command": "bench test --package ./internal/reviewrecord",
          "exit_code": 0
        },
        {
          "id": "dfr-t1-rep1-harness-package",
          "performer": "claude:dfr_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "58796930b5bb3c6483e32048fac3bdb733b175cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1_rep1/harness-package",
            "digest": "sha256:daaa05d927bcb031b54d134032444eba03679add537cd58d40f7c366b7385665",
            "excerpt": "tree: ft354-durable-build 5f1ea9dc7b6612d8095d883eca4730bcc8853cd9 clean\nbench test --package ./internal/durablefile/durablefiletest\ngithub.com/gibbonmi/bench/internal/durablefile/durablefiletest,no-tests,0,0\nfailures[0] skips[0]\nThe package is the shared fault-injection harness; the owner and caller packages exercise it.\n"
          },
          "requirement": "harness-package",
          "command": "bench test --package ./internal/durablefile/durablefiletest",
          "exit_code": 0
        },
        {
          "id": "dfr-t1-rep1-linux-native",
          "performer": "claude:dfr_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "58796930b5bb3c6483e32048fac3bdb733b175cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1_rep1/linux-native",
            "digest": "sha256:b87aedd1381bf5866e700162b1e78c16025b23c07d834870c76f22801ae34166",
            "excerpt": "tree: ft354-durable-build 5f1ea9dc7b6612d8095d883eca4730bcc8853cd9 clean\nhost: Linux 6.6.87.2-microsoft-standard-WSL2 x86_64 (WSL2), GOOS linux/amd64\nbench test --package ./internal/durablefile --run 'TestNativeReplacement'\nrun 1, default kit TMPDIR under /dev/shm, filesystem tmpfs: github.com/gibbonmi/bench/internal/durablefile,pass,2,1; failures[0] skips[0]\nrun 2, caller TMPDIR on a scratch directory with filesystem ext2/ext3 magic (ext4): github.com/gibbonmi/bench/internal/durablefile,pass,17,1; failures[0] skips[0]\nTestNativeReplacement asserts the native trace create, mode, write, file-sync, file-close, rename, directory-open, directory-sync, directory-close with no error, then the native Replace bytes and mode.\nmacOS: not run; no macOS host is available.\n"
          },
          "requirement": "linux-native",
          "command": "bench test --package ./internal/durablefile --run 'TestNativeReplacement'",
          "exit_code": 0
        },
        {
          "id": "dfr-t1-rep1-root-conformance",
          "performer": "claude:dfr_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "58796930b5bb3c6483e32048fac3bdb733b175cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1_rep1/root-conformance",
            "digest": "sha256:e83d297a5b785644bc9b2b094186070eabee22fe483eb0b2a0e496197a2dfb27",
            "excerpt": "tree: ft354-durable-build 5f1ea9dc7b6612d8095d883eca4730bcc8853cd9 clean\nbench test --check root-conformance\nroot-conformance,conformance,1,0\ngithub.com/gibbonmi/bench/internal/conformance,pass,22465,1\nfailures[0] skips[0]\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "dfr-t1-rep1-omit-file-sync",
          "performer": "claude:dfr_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "58796930b5bb3c6483e32048fac3bdb733b175cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1_rep1/omit-file-sync",
            "digest": "sha256:b148798180d93280f6364bf3fc631b0deba48b43ec6c56937d303a6729f904cc",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,swap,failed,6,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1395\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,17,16\nfailures[6]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestNativeReplacement,\"platform_test.go:19: host linux/amd64, destination /dev/shm/JFN4B3/t/TestNativeReplacement303582856/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md\",2\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceAbsentFailures,\"replace_test.go:225: Replace with a file-sync failure = <nil>, want a *durablefile.Error\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceBeforeRenameFailures,\"replace_test.go:213: Replace with a file-sync failure = <nil>, want a *durablefile.Error\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceErrorCause,\"replace_test.go:251: Replace with a file-sync failure = <nil>, want a *durablefile.Error\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceOrder,\"replace_test.go:156: trace = [\\\"create\\\" \\\"mode\\\" \\\"write\\\" \\\"file-close\\\" \\\"rename\\\" \\\"directory-open\\\" \\\"directory-sync\\\" \\\"directory-close\\\"], want the prefix [\\\"create\\\" \\\"mode\\\" \\\"write\\\" \\\"file-sync\\\" \\\"file-close\\\" \\\"rename\\\"]\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceTemporaryCleanup,\"replace_test.go:261: Replace with a file-sync failure = <nil>, want a *durablefile.Error\",1\nskips[0]{package,test,reason}:\n"
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
              "ref": "claude-agent:dfr_t1_rep1/omit-file-sync",
              "digest": "sha256:b148798180d93280f6364bf3fc631b0deba48b43ec6c56937d303a6729f904cc",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,swap,failed,6,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1395\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,17,16\nfailures[6]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestNativeReplacement,\"platform_test.go:19: host linux/amd64, destination /dev/shm/JFN4B3/t/TestNativeReplacement303582856/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md\",2\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceAbsentFailures,\"replace_test.go:225: Replace with a file-sync failure = <nil>, want a *durablefile.Error\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceBeforeRenameFailures,\"replace_test.go:213: Replace with a file-sync failure = <nil>, want a *durablefile.Error\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceErrorCause,\"replace_test.go:251: Replace with a file-sync failure = <nil>, want a *durablefile.Error\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceOrder,\"replace_test.go:156: trace = [\\\"create\\\" \\\"mode\\\" \\\"write\\\" \\\"file-close\\\" \\\"rename\\\" \\\"directory-open\\\" \\\"directory-sync\\\" \\\"directory-close\\\"], want the prefix [\\\"create\\\" \\\"mode\\\" \\\"write\\\" \\\"file-sync\\\" \\\"file-close\\\" \\\"rename\\\"]\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceTemporaryCleanup,\"replace_test.go:261: Replace with a file-sync failure = <nil>, want a *durablefile.Error\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dfr-t1-rep1-omit-directory-sync",
          "performer": "claude:dfr_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "58796930b5bb3c6483e32048fac3bdb733b175cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1_rep1/omit-directory-sync",
            "digest": "sha256:d2583b2a2dbeb40ea2b9b67d5263deb66a1697cb1e6cd5e8bad429045d199dfa",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,swap,failed,6,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1580\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,20,16\nfailures[6]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestNativeReplacement,\"platform_test.go:19: host linux/amd64, destination /dev/shm/TBXFLL/t/TestNativeReplacement3226404670/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md\",2\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceAfterRenameFailures,\"replace_test.go:234: Replace with a directory-sync failure = <nil>, want a *durablefile.Error\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceAfterRenameState,\"replace_test.go:243: Replace with a directory-sync failure = <nil>, want a *durablefile.Error\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceCleanupFailure,\"replace_test.go:275: Replace = replace /dev/shm/TBXFLL/t/TestReplaceCleanupFailure1635106085/003/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md: directory-close: injected directory-close failure; the rename completed, so new bytes may a… (335 bytes)\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceDirectorySync,\"replace_test.go:169: trace = [\\\"create\\\" \\\"mode\\\" \\\"write\\\" \\\"file-sync\\\" \\\"file-close\\\" \\\"rename\\\" \\\"directory-open\\\" \\\"directory-close\\\"], want [\\\"directory-open\\\" \\\"directory-sync\\\" \\\"directory-close\\\"] after the rename\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceErrorCause,\"replace_test.go:251: Replace with a directory-sync failure = <nil>, want a *durablefile.Error\",1\nskips[0]{package,test,reason}:\n"
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
              "ref": "claude-agent:dfr_t1_rep1/omit-directory-sync",
              "digest": "sha256:d2583b2a2dbeb40ea2b9b67d5263deb66a1697cb1e6cd5e8bad429045d199dfa",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,swap,failed,6,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1580\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,20,16\nfailures[6]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestNativeReplacement,\"platform_test.go:19: host linux/amd64, destination /dev/shm/TBXFLL/t/TestNativeReplacement3226404670/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md\",2\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceAfterRenameFailures,\"replace_test.go:234: Replace with a directory-sync failure = <nil>, want a *durablefile.Error\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceAfterRenameState,\"replace_test.go:243: Replace with a directory-sync failure = <nil>, want a *durablefile.Error\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceCleanupFailure,\"replace_test.go:275: Replace = replace /dev/shm/TBXFLL/t/TestReplaceCleanupFailure1635106085/003/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md: directory-close: injected directory-close failure; the rename completed, so new bytes may a… (335 bytes)\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceDirectorySync,\"replace_test.go:169: trace = [\\\"create\\\" \\\"mode\\\" \\\"write\\\" \\\"file-sync\\\" \\\"file-close\\\" \\\"rename\\\" \\\"directory-open\\\" \\\"directory-close\\\"], want [\\\"directory-open\\\" \\\"directory-sync\\\" \\\"directory-close\\\"] after the rename\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceErrorCause,\"replace_test.go:251: Replace with a directory-sync failure = <nil>, want a *durablefile.Error\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dfr-t1-rep1-rename-before-close",
          "performer": "claude:dfr_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "58796930b5bb3c6483e32048fac3bdb733b175cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1_rep1/rename-before-close",
            "digest": "sha256:566b77f8065fa1521c11b6981c7568a1a60c59cbea075056c31432b7091ca146",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,swap,failed,5,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1409\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,17,16\nfailures[5]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestNativeReplacement,\"platform_test.go:19: host linux/amd64, destination /dev/shm/5ESKK4/t/TestNativeReplacement1287020960/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md\",2\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceAbsentFailures,\"replace_test.go:227: after a file-close failure the absent destination is present: <nil>\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceBeforeRenameFailures,\"replace_test.go:217: after a file-close failure the destination = \\\"new\\\", want the old \\\"old bytes\\\"\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceDirectorySync,\"replace_test.go:169: trace = [\\\"create\\\" \\\"mode\\\" \\\"write\\\" \\\"file-sync\\\" \\\"rename\\\" \\\"file-close\\\" \\\"directory-open\\\" \\\"directory-sync\\\" \\\"directory-close\\\"], want [\\\"directory-open\\\" \\\"directory-sync\\\" \\\"directory-close\\\"] after the rename\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceOrder,\"replace_test.go:156: trace = [\\\"create\\\" \\\"mode\\\" \\\"write\\\" \\\"file-sync\\\" \\\"rename\\\" \\\"file-close\\\" \\\"directory-open\\\" \\\"directory-sync\\\" \\\"directory-close\\\"], want the prefix [\\\"create\\\" \\\"mode\\\" \\\"write\\\" \\\"file-sync\\\" \\\"file-close\\\" \\\"rename\\\"]\",1\nskips[0]{package,test,reason}:\n"
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
              "ref": "claude-agent:dfr_t1_rep1/rename-before-close",
              "digest": "sha256:566b77f8065fa1521c11b6981c7568a1a60c59cbea075056c31432b7091ca146",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,swap,failed,5,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1409\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,17,16\nfailures[5]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestNativeReplacement,\"platform_test.go:19: host linux/amd64, destination /dev/shm/5ESKK4/t/TestNativeReplacement1287020960/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md\",2\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceAbsentFailures,\"replace_test.go:227: after a file-close failure the absent destination is present: <nil>\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceBeforeRenameFailures,\"replace_test.go:217: after a file-close failure the destination = \\\"new\\\", want the old \\\"old bytes\\\"\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceDirectorySync,\"replace_test.go:169: trace = [\\\"create\\\" \\\"mode\\\" \\\"write\\\" \\\"file-sync\\\" \\\"rename\\\" \\\"file-close\\\" \\\"directory-open\\\" \\\"directory-sync\\\" \\\"directory-close\\\"], want [\\\"directory-open\\\" \\\"directory-sync\\\" \\\"directory-close\\\"] after the rename\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceOrder,\"replace_test.go:156: trace = [\\\"create\\\" \\\"mode\\\" \\\"write\\\" \\\"file-sync\\\" \\\"rename\\\" \\\"file-close\\\" \\\"directory-open\\\" \\\"directory-sync\\\" \\\"directory-close\\\"], want the prefix [\\\"create\\\" \\\"mode\\\" \\\"write\\\" \\\"file-sync\\\" \\\"file-close\\\" \\\"rename\\\"]\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dfr-t1-rep1-accept-short-write",
          "performer": "claude:dfr_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "58796930b5bb3c6483e32048fac3bdb733b175cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1_rep1/accept-short-write",
            "digest": "sha256:d43d35944a909ea0a51e006092dabfa963571980addaa164bdc26fc72a7d70aa",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,334\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,15,16\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceShortWrite,\"replace_test.go:182: Replace after a short write = <nil>, want io.ErrShortWrite\",1\nskips[0]{package,test,reason}:\n"
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
              "ref": "claude-agent:dfr_t1_rep1/accept-short-write",
              "digest": "sha256:d43d35944a909ea0a51e006092dabfa963571980addaa164bdc26fc72a7d70aa",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,334\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,15,16\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceShortWrite,\"replace_test.go:182: Replace after a short write = <nil>, want io.ErrShortWrite\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dfr-t1-rep1-corrupt-payload",
          "performer": "claude:dfr_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "58796930b5bb3c6483e32048fac3bdb733b175cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1_rep1/corrupt-payload",
            "digest": "sha256:68d41c1ed93510aeeac94169e132ddea70f752e24ae975dce867e40bfcd1e94a",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,swap,failed,4,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,986\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,9,16\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestNativeReplacement,\"platform_test.go:19: host linux/amd64, destination /dev/shm/GIJIVN/t/TestNativeReplacement29938317/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md\",2\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceAfterRenameState,\"replace_test.go:245: after a directory-sync failure the destination = \\\"new p\\\", want the new \\\"new payload\\\"\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceBytes,\"replace_test.go:62: destination = \\\"new\\\\\\\\x00payload\\\\\\\\n\\\\\\\\twith a\\\", want \\\"new\\\\\\\\x00payload\\\\\\\\n\\\\\\\\twith a NUL, a tab, and ü\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceVisibility,\"replace_test.go:126: a concurrent reader observed 524288 bytes that are neither complete payload\",1\nskips[0]{package,test,reason}:\n"
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
              "ref": "claude-agent:dfr_t1_rep1/corrupt-payload",
              "digest": "sha256:68d41c1ed93510aeeac94169e132ddea70f752e24ae975dce867e40bfcd1e94a",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,swap,failed,4,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,986\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,9,16\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestNativeReplacement,\"platform_test.go:19: host linux/amd64, destination /dev/shm/GIJIVN/t/TestNativeReplacement29938317/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md\",2\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceAfterRenameState,\"replace_test.go:245: after a directory-sync failure the destination = \\\"new p\\\", want the new \\\"new payload\\\"\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceBytes,\"replace_test.go:62: destination = \\\"new\\\\\\\\x00payload\\\\\\\\n\\\\\\\\twith a\\\", want \\\"new\\\\\\\\x00payload\\\\\\\\n\\\\\\\\twith a NUL, a tab, and ü\\\\\\\\n\\\"\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceVisibility,\"replace_test.go:126: a concurrent reader observed 524288 bytes that are neither complete payload\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dfr-t1-rep1-corrupt-mode",
          "performer": "claude:dfr_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "58796930b5bb3c6483e32048fac3bdb733b175cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1_rep1/corrupt-mode",
            "digest": "sha256:45c46ea3153f5c212884c0bfee2726c3957180de5a0d165170585327eb32aed7",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,326\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,15,16\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceMode,\"replace_test.go:91: destination permission bits = -rw-------, want -rwxr-xr-x\",1\nskips[0]{package,test,reason}:\n"
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
              "ref": "claude-agent:dfr_t1_rep1/corrupt-mode",
              "digest": "sha256:45c46ea3153f5c212884c0bfee2726c3957180de5a0d165170585327eb32aed7",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,326\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,15,16\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceMode,\"replace_test.go:91: destination permission bits = -rw-------, want -rwxr-xr-x\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dfr-t1-rep1-corrupt-cleanup",
          "performer": "claude:dfr_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "58796930b5bb3c6483e32048fac3bdb733b175cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1_rep1/corrupt-cleanup",
            "digest": "sha256:005eea9cb62096077b795ea0c70e8c71abf5ae710e575d967cba8eb2870bbfce",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,700\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,15,16\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceCleanupFailure,\"replace_test.go:275: Replace = replace /dev/shm/S6RKP2/t/TestReplaceCleanupFailure107273439/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md: write: injected write failure, want both the write cause and the remove cleanup cause\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceTemporaryCleanup,\"replace_test.go:263: after a mode failure the directory holds [\\\".record 'one' *.md.1841320126\\\"] beside the destination\",1\nskips[0]{package,test,reason}:\n"
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
              "ref": "claude-agent:dfr_t1_rep1/corrupt-cleanup",
              "digest": "sha256:005eea9cb62096077b795ea0c70e8c71abf5ae710e575d967cba8eb2870bbfce",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,700\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,15,16\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceCleanupFailure,\"replace_test.go:275: Replace = replace /dev/shm/S6RKP2/t/TestReplaceCleanupFailure107273439/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md: write: injected write failure, want both the write cause and the remove cleanup cause\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceTemporaryCleanup,\"replace_test.go:263: after a mode failure the directory holds [\\\".record 'one' *.md.1841320126\\\"] beside the destination\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dfr-t1-rep1-corrupt-cause",
          "performer": "claude:dfr_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "58796930b5bb3c6483e32048fac3bdb733b175cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1_rep1/corrupt-cause",
            "digest": "sha256:fc0833973a89e133f99f8da3d1e1094b771fb3135dc55333826e36e18e7cb889",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1038\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,16,16\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceCleanupFailure,\"replace_test.go:275: Replace = replace /dev/shm/SJODBG/t/TestReplaceCleanupFailure4185829830/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md: write: injected write failure\",2\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceErrorCause,\"replace_test.go:253: a create failure = replace /dev/shm/SJODBG/t/TestReplaceErrorCause3536961005/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md: create: injected create failure, which does not hold the injected cause\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceShortWrite,\"replace_test.go:182: Replace after a short write = replace /dev/shm/SJODBG/t/TestReplaceShortWrite4032543831/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md: write: short write, want io.ErrShortWrite\",1\nskips[0]{package,test,reason}:\n"
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
              "ref": "claude-agent:dfr_t1_rep1/corrupt-cause",
              "digest": "sha256:fc0833973a89e133f99f8da3d1e1094b771fb3135dc55333826e36e18e7cb889",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/durablefile/replace.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/durablefile,all,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1038\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/durablefile,fail,16,16\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceCleanupFailure,\"replace_test.go:275: Replace = replace /dev/shm/SJODBG/t/TestReplaceCleanupFailure4185829830/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md: write: injected write failure\",2\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceErrorCause,\"replace_test.go:253: a create failure = replace /dev/shm/SJODBG/t/TestReplaceErrorCause3536961005/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md: create: injected create failure, which does not hold the injected cause\",1\n  github.com/gibbonmi/bench/internal/durablefile,TestReplaceShortWrite,\"replace_test.go:182: Replace after a short write = replace /dev/shm/SJODBG/t/TestReplaceShortWrite4032543831/001/parent \\\"quoted\\\" [glob]*? ünïcödé/record 'one' *.md: write: short write, want io.ErrShortWrite\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dfr-t1-rep1-caller-bypass",
          "performer": "claude:dfr_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "58796930b5bb3c6483e32048fac3bdb733b175cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1_rep1/caller-bypass",
            "digest": "sha256:245b39aab90bf767190c428f79f203a71cfa88abbdb183dd1f37eebfed0b0a71",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/write.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,all,passed,89\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,875\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,1656,89\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewReplacement/assignment_write_directory_sync_failure_after_the_rename,\"replacement_test.go:130: write with a failed directory sync = <nil>, want the published directory-sync cause\",1\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewReplacement/file_sync_failure_before_the_rename,\"replacement_test.go:154: write with a failed file sync = <nil>, want the unpublished file-sync cause\",1\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewReplacement/record_write_directory_sync_failure_after_the_rename,\"replacement_test.go:105: write with a failed directory sync = <nil>, want the published directory-sync cause\",1\nskips[0]{package,test,reason}:\n"
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
              "ref": "claude-agent:dfr_t1_rep1/caller-bypass",
              "digest": "sha256:245b39aab90bf767190c428f79f203a71cfa88abbdb183dd1f37eebfed0b0a71",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/write.go,swap,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,all,passed,89\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,875\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,1656,89\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewReplacement/assignment_write_directory_sync_failure_after_the_rename,\"replacement_test.go:130: write with a failed directory sync = <nil>, want the published directory-sync cause\",1\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewReplacement/file_sync_failure_before_the_rename,\"replacement_test.go:154: write with a failed file sync = <nil>, want the unpublished file-sync cause\",1\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewReplacement/record_write_directory_sync_failure_after_the_rename,\"replacement_test.go:105: write with a failed directory sync = <nil>, want the published directory-sync cause\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dfr-t1-rep1-caller-payload",
          "performer": "claude:dfr_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "58796930b5bb3c6483e32048fac3bdb733b175cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1_rep1/caller-payload",
            "digest": "sha256:31e97d1fb9cddab031b94fbc5bc93059935d12b057e43ed51936f519fdea59fe",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/write.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,all,passed,89\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,657\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,1673,89\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewReplacement/record_write_directory_sync_failure_after_the_rename,\"replacement_test.go:107: record = \\\"# Review outcomes\\\\\\\\n\\\\\\\\n```bench-review-record\\\\\\\\n{\\\\\\\\n  \\\\\\\\\\\"version\\\\\\\\\\\": 1,\\\\\\\\n  \\\\\\\\\\\"spec\\\\\\\\\\\": \\\\\\\\\\\"specs/example/spec.md\\\\\\\\\\\",\\\\\\\\n  \\\\\\\\\\\"plan_digest\\\\\\\\\\\": \\\\\\\\\\\"sha256:plan\\\\\\\\\\\",\\\\\\\\n  \\\\\\\\\\\"implementation_session\\\\\\\\\\\": \\\\\\\\\\\"author\\\\\\\\\\\",\\\\\\\\n  \\\\\\\\\\\"chunks\\\\\\\\\\\": null… (1020 bytes)\",3\nskips[0]{package,test,reason}:\n"
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
              "ref": "claude-agent:dfr_t1_rep1/caller-payload",
              "digest": "sha256:31e97d1fb9cddab031b94fbc5bc93059935d12b057e43ed51936f519fdea59fe",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/write.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,all,passed,89\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,657\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,1673,89\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewReplacement/record_write_directory_sync_failure_after_the_rename,\"replacement_test.go:107: record = \\\"# Review outcomes\\\\\\\\n\\\\\\\\n```bench-review-record\\\\\\\\n{\\\\\\\\n  \\\\\\\\\\\"version\\\\\\\\\\\": 1,\\\\\\\\n  \\\\\\\\\\\"spec\\\\\\\\\\\": \\\\\\\\\\\"specs/example/spec.md\\\\\\\\\\\",\\\\\\\\n  \\\\\\\\\\\"plan_digest\\\\\\\\\\\": \\\\\\\\\\\"sha256:plan\\\\\\\\\\\",\\\\\\\\n  \\\\\\\\\\\"implementation_session\\\\\\\\\\\": \\\\\\\\\\\"author\\\\\\\\\\\",\\\\\\\\n  \\\\\\\\\\\"chunks\\\\\\\\\\\": null… (1020 bytes)\",3\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dfr-t1-rep1-caller-mode",
          "performer": "claude:dfr_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "58796930b5bb3c6483e32048fac3bdb733b175cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1_rep1/caller-mode",
            "digest": "sha256:9de8b001a1691885b62af67e20806e9ad77cdf62c8079d10e4e12a2323539e2c",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/write.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,all,passed,89\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,851\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,1626,89\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewReplacement/assignment_write_directory_sync_failure_after_the_rename,\"replacement_test.go:135: /dev/shm/6DX2TM/t/TestReviewReplacementassignment_write_directory_sync_failure_af3749477331/002/linked/specs/example/spec.md permission bits = -rw-------, want 0644\",1\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewReplacement/record_write_directory_sync_failure_after_the_rename,\"replacement_test.go:109: /dev/shm/6DX2TM/t/TestReviewReplacementrecord_write_directory_sync_failure_after_1541071386/001/root \\\"quoted\\\" [glob]*? ünïcödé/reviews/example.md permission bits = -rw-------, want 0644\",1\nskips[0]{package,test,reason}:\n"
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
              "ref": "claude-agent:dfr_t1_rep1/caller-mode",
              "digest": "sha256:9de8b001a1691885b62af67e20806e9ad77cdf62c8079d10e4e12a2323539e2c",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/write.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,all,passed,89\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,851\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,1626,89\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewReplacement/assignment_write_directory_sync_failure_after_the_rename,\"replacement_test.go:135: /dev/shm/6DX2TM/t/TestReviewReplacementassignment_write_directory_sync_failure_af3749477331/002/linked/specs/example/spec.md permission bits = -rw-------, want 0644\",1\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewReplacement/record_write_directory_sync_failure_after_the_rename,\"replacement_test.go:109: /dev/shm/6DX2TM/t/TestReviewReplacementrecord_write_directory_sync_failure_after_1541071386/001/root \\\"quoted\\\" [glob]*? ünïcödé/reviews/example.md permission bits = -rw-------, want 0644\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dfr-t1-rep1-caller-refusal",
          "performer": "claude:dfr_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "58796930b5bb3c6483e32048fac3bdb733b175cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1_rep1/caller-refusal",
            "digest": "sha256:57fe7f4d02bcb4ecf4b9b57c5515b8664b48bd7e23213d42010cca40ac8367ae",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/write.go,omit,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,all,passed,89\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,423\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,1651,89\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewReplacement/refusal_before_the_replacement,\"replacement_test.go:179: write = <nil> after 1 replacements, want the refusal \\\"invalid record spec identity\\\" before any replacement\",1\nskips[0]{package,test,reason}:\n"
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
              "ref": "claude-agent:dfr_t1_rep1/caller-refusal",
              "digest": "sha256:57fe7f4d02bcb4ecf4b9b57c5515b8664b48bd7e23213d42010cca40ac8367ae",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/write.go,omit,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,all,passed,89\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,423\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,1651,89\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewReplacement/refusal_before_the_replacement,\"replacement_test.go:179: write = <nil> after 1 replacements, want the refusal \\\"invalid record spec identity\\\" before any replacement\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dfr-t1-rep1-write-binding",
          "performer": "claude:dfr_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "58796930b5bb3c6483e32048fac3bdb733b175cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1_rep1/write-binding",
            "digest": "sha256:ac8cc948ab4f56da915f86cbe9354ebb38d70360d179f19b5d46ee5f6c898476",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/write.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,all,passed,89\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,633\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,1780,89\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewReplacement/file_sync_failure_before_the_rename,\"replacement_test.go:154: write with a failed file sync = <nil>, want the unpublished file-sync cause\",1\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewReplacement/record_write_directory_sync_failure_after_the_rename,\"replacement_test.go:105: write with a failed directory sync = <nil>, want the published directory-sync cause\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "write-binding-proof",
          "command": "bench test --package ./internal/reviewrecord",
          "exit_code": 0,
          "probe": {
            "mutation": "In the record write entry of internal/reviewrecord/write.go, replace the replace call with os.WriteFile of the rendered bytes and 0644. The record-write case of TestReviewReplacement must fail its directory-sync cause assertion, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:dfr_t1_rep1/write-binding",
              "digest": "sha256:ac8cc948ab4f56da915f86cbe9354ebb38d70360d179f19b5d46ee5f6c898476",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/write.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,all,passed,89\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,633\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,1780,89\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewReplacement/file_sync_failure_before_the_rename,\"replacement_test.go:154: write with a failed file sync = <nil>, want the unpublished file-sync cause\",1\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewReplacement/record_write_directory_sync_failure_after_the_rename,\"replacement_test.go:105: write with a failed directory sync = <nil>, want the published directory-sync cause\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "dfr-t1-rep1-assignment-binding",
          "performer": "claude:dfr_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "58796930b5bb3c6483e32048fac3bdb733b175cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_t1_rep1/assignment-binding",
            "digest": "sha256:648b0fa7e3c071973f3002a6c596054f82decd10788118d32af70fd0e286f00d",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/assignment.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,all,passed,89\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,424\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,1616,89\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewReplacement/assignment_write_directory_sync_failure_after_the_rename,\"replacement_test.go:130: write with a failed directory sync = <nil>, want the published directory-sync cause\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "assignment-binding-proof",
          "command": "bench test --package ./internal/reviewrecord",
          "exit_code": 0,
          "probe": {
            "mutation": "In RecordAssignment of internal/reviewrecord/assignment.go, replace the replace call with os.WriteFile of the candidate bytes and 0644. The assignment case of TestReviewReplacement must fail its directory-sync cause assertion, then pass after source restoration.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:dfr_t1_rep1/assignment-binding",
              "digest": "sha256:648b0fa7e3c071973f3002a6c596054f82decd10788118d32af70fd0e286f00d",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/assignment.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,all,passed,89\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,424\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,1616,89\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewReplacement/assignment_write_directory_sync_failure_after_the_rename,\"replacement_test.go:130: write with a failed directory sync = <nil>, want the published directory-sync cause\",1\nskips[0]{package,test,reason}:\n"
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
        },
        {
          "id": "dfr-da-standards-2",
          "performer": "claude:dfr_rev_standards_2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "58796930b5bb3c6483e32048fac3bdb733b175cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_rev_standards_2",
            "digest": "sha256:0133e505ea7ad8b49502353b74c5b63fc8957abd6ea635c06812299d7fea5dd9",
            "excerpt": "Standards confirming round, chunk D-A, repair delta fc9a1d33..5f1ea9dc, opus/high. Evidence current=true at bd53453f.\nS1 closed: the harness lives once in internal/durablefile/durablefiletest/recorder.go; owner tests (replace_test.go, platform_test.go) and caller test (replacement_test.go:21,61) use it; no other copy in internal/.\nS2 closed: the sentence is gone from platform_test.go:14-16.\nreplaceFile binding (write.go:61-63) and export_test.go conform to AGENTS.md and the comment register. No new blocking defect.\nAdvice: rename the local replaceFile at replacement_test.go:122.\n"
          },
          "axis": "Standards",
          "base": "9433df1bc7cc583dd77f237856b07c9cbf29250f",
          "tip": "5f1ea9dc7b6612d8095d883eca4730bcc8853cd9",
          "finding_ids": [],
          "supersedes": [
            "dfr-da-standards-1"
          ]
        },
        {
          "id": "dfr-da-spec-2",
          "performer": "claude:dfr_rev_spec_2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "58796930b5bb3c6483e32048fac3bdb733b175cc",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:dfr_rev_spec_2",
            "digest": "sha256:589771ca19aecfd358de8a42509ffc512bae956084e5009967926fd4d160eee4",
            "excerpt": "Spec confirming round, chunk D-A, repair delta fc9a1d33..5f1ea9dc plus spec amendment 87e2713e..d1053edb, opus/high. Evidence current=true at bd53453f.\nP1 open: D32 macOS evidence absent (reviews lines 84, 522, 1071); reviewer decision.\nP2 closed: one binding write.go:63 serves write.go:100 and assignment.go:84; assignment case replacement_test.go:120-145 drives the real RecordAssignment.\nAmended implementation decisions met. Rows D01-D15 and D19 hold their cited tests. No row broken by the repair.\n"
          },
          "axis": "Spec",
          "base": "9433df1bc7cc583dd77f237856b07c9cbf29250f",
          "tip": "5f1ea9dc7b6612d8095d883eca4730bcc8853cd9",
          "finding_ids": [
            "P1"
          ],
          "supersedes": [
            "dfr-da-spec-1"
          ]
        },
        {
          "id": "dfr-da-coverage-2",
          "performer": "claude:dfr_rev_coverage_2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "58796930b5bb3c6483e32048fac3bdb733b175cc",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:dfr_rev_coverage_2",
            "digest": "sha256:36cb8cd7dc2f070deb1c8022bd6f87f4e5291bc4e01cc4ca10066fdc00fdd6f3",
            "excerpt": "Coverage confirming round, chunk D-A, repair delta fc9a1d33..5f1ea9dc, opus/high. The axis ran --check-current from the primary checkout (refused: base not an ancestor), so its currency check is unknown; it read and tested the worktree tip bd53453f.\nC1 open: D32 macOS evidence absent (reviews:1079-1080).\nC2 closed: the record-write case binds a directory-sync failure through the real write entry (replacement_test.go:104-105); write-binding-proof bit and restored.\nC3 closed: the assignment case drives the real RecordAssignment (replacement_test.go:115-142); assignment-binding-proof bit and restored.\nAll nine stage witnesses kept after the harness move. Ran bench test on ./internal/reviewrecord (89 pass) and ./internal/durablefile/... (16 pass). No new gap.\n"
          },
          "axis": "Coverage",
          "base": "9433df1bc7cc583dd77f237856b07c9cbf29250f",
          "tip": "5f1ea9dc7b6612d8095d883eca4730bcc8853cd9",
          "finding_ids": [
            "C1"
          ],
          "supersedes": [
            "dfr-da-coverage-1"
          ]
        },
        {
          "id": "dfr-da-spec-3",
          "performer": "claude:dfr_rev_spec_3",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "58796930b5bb3c6483e32048fac3bdb733b175cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_rev_spec_3",
            "digest": "sha256:e8bd183df77658ef3ddf2b8040b062f659acfc39142ff80810c8bd45a0b33ed9",
            "excerpt": "Spec reaffirmation, chunk D-A, rescope diff 5280f304..5303ca3c, opus/high. Evidence current=true at a4fc1773.\nP1 closed: the reviewer decided on 2026-10-12 that D32 is Linux only. Row D32 (spec.md:189), Testing decisions (:144), Won't handle (:201), the chunk table (:121), and the ticket say Linux only; the plan (:135) has no macos-native entry; Out of scope (:222) keeps the macOS qualification open on FT354. The linux-native entry dfr-t1-rep1-linux-native covers the rescoped row; internal/durablefile has no diff from 5f1ea9dc to a4fc1773. No new inconsistency.\n"
          },
          "axis": "Spec",
          "base": "9433df1bc7cc583dd77f237856b07c9cbf29250f",
          "tip": "5f1ea9dc7b6612d8095d883eca4730bcc8853cd9",
          "finding_ids": [],
          "supersedes": [
            "dfr-da-spec-2"
          ]
        },
        {
          "id": "dfr-da-coverage-3",
          "performer": "claude:dfr_rev_coverage_3",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "58796930b5bb3c6483e32048fac3bdb733b175cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:dfr_rev_coverage_3",
            "digest": "sha256:1d0cf9ba63940b41fca4ae48c8f65fbe8377f7ea31839bd6aa12d36fc85652b6",
            "excerpt": "Coverage reaffirmation, chunk D-A, rescope diff 5280f304..5303ca3c, opus/high. Evidence current=true at a4fc1773.\nC1 closed: the plan no longer requires macos-native (spec.md:121,135). TestNativeReplacement (platform_test.go:17-37) asserts the native trace with the directory sync, and the linux-native entry records tmpfs and ext4 passes on WSL2 x86_64. Reran TestNativeReplacement at a4fc1773: pass. No other row lost coverage. No new gap.\n"
          },
          "axis": "Coverage",
          "base": "9433df1bc7cc583dd77f237856b07c9cbf29250f",
          "tip": "5f1ea9dc7b6612d8095d883eca4730bcc8853cd9",
          "finding_ids": [],
          "supersedes": [
            "dfr-da-coverage-2"
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
    "verification": [
      {
        "id": "dfr-final-ordinary-integration",
        "performer": "claude:dfr-orchestrator-20261011",
        "role": "integration-verification",
        "model": "opus",
        "effort": "high",
        "source_digest": "9107a8aeff98e6bb835081228a1dd13b49ea1939",
        "state": "completed",
        "outcome": "pass",
        "native_ref": {
          "ref": "claude-session:dfr-orchestrator-20261011",
          "digest": "sha256:ee6b5ca0b560508b548e21a3261c1412ebff96c72e7698a2cde111cb088f7065",
          "excerpt": "bench test --package ./... at 95969942 (WSL2 linux/amd64): 127 packages pass, failures[0], skips[7] for host capability (character device creation not permitted; no bin/bench.sh in a capability subject root).\n"
        },
        "requirement": "ordinary-integration",
        "command": "bench test --package ./...",
        "exit_code": 0
      },
      {
        "id": "dfr-final-coverage",
        "performer": "claude:dfr-orchestrator-20261011",
        "role": "integration-verification",
        "model": "opus",
        "effort": "high",
        "source_digest": "9107a8aeff98e6bb835081228a1dd13b49ea1939",
        "state": "completed",
        "outcome": "pass",
        "native_ref": {
          "ref": "claude-session:dfr-orchestrator-20261011",
          "digest": "sha256:71226fbc8bd22540a2e1539d36a98e8a6aec0b8c0cf9dfa8b5543b0e87a243b0",
          "excerpt": "bench coverage --check durable-file-replacement at 95969942: ok: coverage map valid — 17 row(s)\n"
        },
        "requirement": "coverage",
        "command": "bench coverage --check durable-file-replacement",
        "exit_code": 0
      },
      {
        "id": "dfr-final-root-conformance",
        "performer": "claude:dfr-orchestrator-20261011",
        "role": "integration-verification",
        "model": "opus",
        "effort": "high",
        "source_digest": "9107a8aeff98e6bb835081228a1dd13b49ea1939",
        "state": "completed",
        "outcome": "pass",
        "native_ref": {
          "ref": "claude-session:dfr-orchestrator-20261011",
          "digest": "sha256:9c1715032bd579af748bb2edccd95d7fbaa950c2e25745cc8abd0d390905e8ea",
          "excerpt": "bench test --check root-conformance at 95969942: internal/conformance pass, tests_run 1, failures[0], skips[0]\n"
        },
        "requirement": "root-conformance",
        "command": "bench test --check root-conformance",
        "exit_code": 0
      }
    ]
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
    },
    {
      "from": "sha256:7a4c1e409bfe8a274ca0cdaca77ce869f2efab34752d3692aa7cb5e7e6057754",
      "to": "sha256:7a8be2c6a86c8f70f2ee4dd22898e1f83dee922fd0f2c5366c3d885eda7ebe5b",
      "chunk_ids": {
        "D-A": [
          "D-A"
        ]
      }
    },
    {
      "from": "sha256:7a8be2c6a86c8f70f2ee4dd22898e1f83dee922fd0f2c5366c3d885eda7ebe5b",
      "to": "sha256:13c36574a31cf4b389ddc0c9f802e8f0574155402083cdaf8f59fe459d757a1a",
      "chunk_ids": {
        "D-A": [
          "D-A"
        ]
      }
    },
    {
      "from": "sha256:13c36574a31cf4b389ddc0c9f802e8f0574155402083cdaf8f59fe459d757a1a",
      "to": "sha256:286fab6eee12fd2fb29be89f7acf16d6242dd8e57be66c984d59206f683de6df",
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
  - Repair 5f1ea9dc: fixed. The `durablefiletest` package holds the harness, and the owner tests and `TestReviewReplacement` use it.
- S2, auto-fix, confidence 5: `internal/durablefile/platform_test.go:15-16` states the evidence obligation that ticket line 60 owns. Delete that sentence.
  - Repair 5f1ea9dc: fixed. `TestNativeReplacement` does not state the evidence obligation.

## Spec

Chunk D-A has 2 findings. The worst issue is P1.

- P1, ask-user, confidence 10: row D32 needs native macOS evidence, and this host has no macOS. The reviewer decides the route.
  - Repair 5f1ea9dc: not fixed. This host has no macOS, so the reviewer decision stays open.
  - Closed at 5303ca3c: the reviewer decided on 2026-10-12 that D32 is Linux only, and FT354 keeps the macOS qualification open.
- P2, no-op, confidence 4: the C3 repair closes the same gap in `internal/reviewrecord/assignment.go:85`.
  - Repair 5f1ea9dc: fixed by the C3 repair.

## Coverage

Chunk D-A has 3 findings. The worst issue is C1.

- C1, ask-user, confidence 9: this finding is the same as P1.
  - Repair 5f1ea9dc: not fixed. This host has no macOS, so the reviewer decision stays open.
  - Closed at 5303ca3c: the reviewer decided on 2026-10-12 that D32 is Linux only, and FT354 keeps the macOS qualification open.
- C2, auto-fix, confidence 7: `TestReviewReplacement` calls `writeWith` directly, so no test examines the replacement that the `write` entry binds (`internal/reviewrecord/write.go:196-198`). Inject the failure through the real `write` entry.
  - Repair 5f1ea9dc: fixed. The record-write case of `TestReviewReplacement` injects the failure through the one `replaceFile` binding and the real `write` entry, and the write-binding-proof probe bit.
- C3, auto-fix, confidence 7: `RecordAssignment` binds its own replacement (`internal/reviewrecord/assignment.go:85`), and no injected failure reaches it. Inject the failure through the real `RecordAssignment` transaction. C2 and C3 have one repair target.
  - Repair 5f1ea9dc: fixed. The assignment case of `TestReviewReplacement` injects the failure through the real `RecordAssignment` transaction, and the assignment-binding-proof probe bit.

Repair cycles consumed: 1 of 2.
