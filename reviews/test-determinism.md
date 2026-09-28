# Test determinism review

## Run

The user directs one inline implementation author and independent Sol/high reviews.
The version 1 plan retains the inline author identity.
The native implementation model identifier is unknown.
Review line: gpt-6-sol / high / one iteration per axis / three readers
Expected repair rounds: 1 / confidence 7

## TD-C1a

The environment owner and gate composition cover TD1 through TD6 and TD8 through TD15.
The focused environment and gate packages pass.
The root conformance test passes with no skips.
The gittest package has no standalone tests; its probe runs through the environment and gate tests.

The independent expectations detect the recorded mutations below.
Each valid probe reports a behavioral failure and confirms restoration.
The local probe transcript is .logs/test-determinism-probes.json.

| rows | mutation | observed failure |
| --- | --- | --- |
| TD1 | Remove the global git configuration override. | The global marker remains visible. |
| TD2 | Remove the system git configuration override. | The system marker remains visible. |
| TD3 | Remove the maintenance entries. | The commit starts auto-maintenance. |
| TD4 | Preserve the operator home. | HOME is outside the private run. |
| TD5, TD13 | Preserve the operator temporary directory. | TMPDIR is outside the private run. |
| TD6 | Empty the Go setting pins. | The probe reports changed pins. |
| TD8 | Omit directory permission restoration. | Close reports permission denied. |
| TD9 | Use a long run-directory prefix. | The private path exceeds the 16-byte limit. |
| TD10, TD11 | Bypass HOME validation. | Open accepts absent, empty, and relative values. |
| TD12 | Create the run outside the specified TMPDIR. | Open accepts a regular file as TMPDIR. |
| TD14 | Apply the kit policy to a linked root. | The linked phase loses the operator home. |
| TD15 | Omit the run entries from the phase environment. | The phase reports that HOME is not private. |

The first phase-composition mutation failed to compile and supplied no behavioral evidence.
The replacement mutation compiles and produces the intended failure.

Repair cycles consumed: 0 of 2

## Standards

Pending independent review.

## Spec

Pending independent review.
Later chunks retain their planned rows; this chunk claims only the rows named above.

## Coverage

Pending independent review.

```bench-review-record
{"version":1,"spec":"specs/test-determinism/spec.md","plan_digest":"sha256:0c1f19eeefd300b73018e190ea2fb3acf0baf416e4171f4b8717cfa5aa651a63","implementation_session":"codex/test-determinism-inline-20260928","chunks":[{"id":"TD-C1a","base":"72e751b7a020a76aa35286852b4ea981c95cb865","tip":"72230a46458b3a4be1161e20a1efe7abbf0943e2","plan_digest":"sha256:0c1f19eeefd300b73018e190ea2fb3acf0baf416e4171f4b8717cfa5aa651a63","source_digest":"98e65974fca9d5eed474afbb219587d1150d3d82","acceptance_rows":["TD1","TD2","TD3","TD4","TD5","TD6","TD8","TD9","TD10","TD11","TD12","TD13","TD14","TD15"],"verification":[{"id":"TD-C1a-env-1","performer":"codex/test-determinism-inline-20260928","role":"author-verification","model":"unknown","effort":"high","source_digest":"98e65974fca9d5eed474afbb219587d1150d3d82","state":"completed","outcome":"pass","native_ref":{"ref":"codex:exec-session-36166","digest":"sha256:74e3f8f12fd1460396b8a205d6fdbae43a45e08c0f2b3bf64b663d18158ed388","excerpt":"packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/env,pass,510\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"},"requirement":"env","command":"bench test --package ./internal/env","exit_code":0},{"id":"TD-C1a-gate-1","performer":"codex/test-determinism-inline-20260928","role":"author-verification","model":"unknown","effort":"high","source_digest":"98e65974fca9d5eed474afbb219587d1150d3d82","state":"completed","outcome":"pass","native_ref":{"ref":"codex:exec-session-50094","digest":"sha256:009d83ab1d9703f56d5a60d3f12aebb478f053b62e871be5bcb48211dca46b18","excerpt":"packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,pass,10115\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"},"requirement":"gate","command":"bench test --package ./internal/gate","exit_code":0}],"reviews":[]}],"completion":{"state":"pending","source_digest":"","performer":"codex/test-determinism-inline-20260928","reconciliation":{},"verification":[]}}
```
