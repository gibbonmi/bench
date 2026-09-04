# Publish the manifest where the routes read it

Blocked by: none
Writes: scripts/release-preflight.sh, internal/brokermanifest/brokermanifest.go, internal/runbinary/runbinary.go, internal/runbinary/runbinary_test.go, internal/worktree/build_test.go, internal/adopt/adopt_test.go, cmd/bench/main_test.go, internal/preflight/decision_test.go, internal/freshness/freshness_publish.go, internal/freshness/freshness_publish_test.go, internal/freshness/freshness_verify_test.go, internal/runbinary/runbinary_test.go, internal/conformance/gate_entry_test.go, internal/systemtest/owner_artifact_recovery_test.go, cmd/bench/freshness_publish.go, cmd/bench/freshness_publish_test.go, cmd/bench/build_artifact_mode_test.go, cmd/bench/build_subject_mode_test.go, internal/adopt/broker_test.go, scripts/go-build.sh, specs/binary-freshness/spec.md
Covers: BF10, BF14

## What to build

This ticket repairs review finding F1.

The spec says the manifest lands beside the resolved wrapper. The publish
transaction instead derives the directory from the published executable, so
an ordinary build writes `dist/bench-broker.manifest`. The complete reader
enumeration is `bin/bench.sh` and `internal/adopt/broker.go`, and both name
the wrapper's `bin/`. `.gitignore` and `.bench/build-outputs.json` name
`bin/bench-broker.manifest` as well. So stories 10 and 13 deliver nothing on
the ordinary build path.

The reviewer decided the route on 2026-09-04. Give `freshness.Publish` and
the `freshness-publish` verb a manifest-directory operand, and make
`scripts/go-build.sh` pass the kit's `bin/`. The manifest still lands inside
the same transaction, so a rollback still restores the prior one. Artifact
mode stays excluded.

Refuse an empty manifest directory, the way the version operand is refused.
A caller that names no directory has nothing to publish against.

## Acceptance

- [ ] After a subject-mode build the manifest is at `<root>/bin/bench-broker.manifest`, and its digest equals the published executable's digest.
- [ ] A publication rolled back leaves the prior manifest unchanged.
- [ ] Artifact mode still writes no manifest and executes nothing.
- [ ] `bench doctor` reads the same manifest the land route reads.
- [ ] Self-probe: pass the published executable's own directory again, and report which row reds.
