# C1 comment repair evidence

The reviewer approved one additional scoped repair cycle on 2026-10-02.
Cycle 4 removes the two comment lines cited by ST-R3-1.
The build script remains the owner of the manifest-lifetime explanation.
No executable statement or test assertion changed.
The original author session retains this repair.

## Current verification

| Route | Result | Elapsed milliseconds |
| --- | --- | --- |
| internal/runbinary | Pass, no skips | 23643 |
| internal/adopt | Pass, no skips | 21846 |
| internal/compatibility | Pass, no skips | 4 |
| internal/preflight/evidencecmd | Pass, no skips | 19299 |
| cmd/bench | Pass, no skips | 17359 |
| ordinary-build-census | Pass, no skips | 259 |
| package-core-guard | Pass, no skips | 2377 |
| Restored system suite | Pass, no skips | 67363 |

The manifest-directory swap failed both intended fixture cases in 11403 ms.
The absent case rejected publication, and the present case rejected changed bytes.
The probe reported bit and verified source restoration.

The doctor-route swap failed TestCompatibilityMissingPath in the sealed system suite.
That run took 69586 ms and had no skips.
The original bytes and mode were restored from a verified private backup.
The restored SHA-256 is df42cb91f526209d9f1f6512e0b362b63a72e63742af06f9014703d633e563cf.
The restored system suite passed as recorded above.

All commands used the installed Node 25.8.1 and npm 11.11.0 through a process-local PATH prefix.
No user configuration or runtime file changed.
Normal desktop process startup still fails under its usual sandbox permissions.
Current native review and the quiet C1 checkpoint remain required.
