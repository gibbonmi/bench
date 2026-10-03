# C1 artifact repair evidence

The reviewer approved one additional C1 repair cycle on 2026-10-02.
The original author session retained all implementation and repair work.
This cycle addresses writes to ignored build artifacts during the checkpoint.
The initial two cycles remain consumed, and their findings remain closed.

## Observed writers

The quiet checkpoint at cf90f1df10e213c11fa708e0d6e91bf1590e8f3c changed three ignored artifacts.
No concurrent worktree call ran during that checkpoint.
All test phases passed, but the checkout guard refused the changed files.
The tracked worktree remained clean.

The desktop shell selected Node 18.19.1 and npm 9.2.0.
The repository declares Node 24 or later.
An isolated package-core-guard run reproduced all three artifact writes.
The process observer traced npm pack to prepare and then to the live build script.
The installed npm 9 directory packer does not check ignoreScripts before prepare.

Node 25.8.1 and npm 11.11.0 were already installed.
Their directory was prepended to PATH for these commands, preserving the remaining harness paths.
With that selection, package-core-guard passed in 2246 ms without skips.
Complete binary, seal, and manifest snapshots stayed unchanged.
No user configuration or installed runtime file was edited.

The next quiet checkpoint preserved the binary and seal but changed the manifest timestamp.
The run-binary manifest test unconditionally rewrote that live file during cleanup.
The single test reproduced the write in 7316 ms without skips.
The repair uses the existing private kit-copy fixture and removes the live-file cleanup.
It explicitly tests both an absent manifest and an existing manifest.

## Regression proof

The repaired manifest test passed in 9182 ms without skips.
The live artifact bytes, modes, sizes, and modification times stayed unchanged.
The planned manifest-directory swap failed both cases in 9899 ms.
The absent case reported an unwanted publication; the present case reported changed bytes.
The probe reported bit and verified source restoration.

The first fixture attempt used equals signs in its subtest names.
That temporary executable path was parsed as an assignment by env, so the build failed before the assertions.
The corrected names are absent and present.
This setup failure is not a behavioral mutation result.

The required doctor-to-setup route swap failed the sealed system suite in 59905 ms.
TestCompatibilityMissingPath reported setup usage and exit 2 instead of the expected doctor result.
The original bytes and mode were restored from a verified private backup.
The restored SHA-256 is df42cb91f526209d9f1f6512e0b362b63a72e63742af06f9014703d633e563cf.
The restored sealed system suite passed in 61053 ms without skips.

## Current verification

| Route | Result | Elapsed milliseconds |
| --- | --- | --- |
| internal/runbinary | Pass, no skips | 22766 |
| internal/adopt | Pass, no skips | 21003 |
| internal/compatibility | Pass, no skips | 4 |
| internal/preflight/evidencecmd | Pass, no skips | 19397 |
| cmd/bench | Pass, no skips | 17227 |
| ordinary-build-census | Pass, no skips | 285 |
| Restored system suite | Pass, no skips | 61053 |

The exploratory full conformance runs had three declared environment skips.
Those runs served diagnosis and do not replace the required checks above.
The implementation still requires current independent review and a green C1 checkpoint.
Normal desktop process startup still fails under its usual sandbox permissions.
The recovered toolchain and escalated diagnostics do not qualify that interface.
