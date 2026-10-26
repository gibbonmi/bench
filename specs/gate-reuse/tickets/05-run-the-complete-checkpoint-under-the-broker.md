# Run the complete checkpoint under the landing broker

Blocked by: 02-name-evidence-by-the-grade-key.md
Writes: internal/gate/gate.go, internal/gate/complete_checkpoint.go, internal/gate/checkpoint.go, internal/gate/complete_checkpoint_broker_test.go (new), internal/gate/refusal_route_test.go, internal/gate/refusal_route_fixtures_test.go (new), internal/refusalroute/faces_gate.go
Covers: GR59, GR60, GR61, GR62, GR63, GR64

## What to build

In a repository that declares Bench build inputs, `bench gate --checkpoint <spec-path> --complete` compares its running executable with the primary checkout's broker manifest.
It reads the manifest and digests both executables through `internal/brokermanifest`, as `checkCandidate` in `internal/commitment/repository/window.go` does.
The comparison runs in `RunCommand` before the run log and the span begin, so a process that execs leaves no open run record.

- When the running executable's digest equals the manifest's, the checkpoint runs in its own process.
- When they differ, it execs the manifest's broker with the same arguments and the same working directory. A replaceable exec function carries the call.
- When the manifest does not read, or the broker's digest differs from the manifest, it refuses with the `checkpoint-broker-unavailable` face and runs no gate.

A repository that declares no Bench build inputs keeps its running executable and reads no manifest. A chunk checkpoint keeps its running executable.
The new face joins `internal/refusalroute/faces_gate.go`. Its route repeats the landing's broker repair advice, and a fixture in `gateFaceFixtures` produces and follows it.
`refusal_route_test.go` is at its 400-line limit. Move fixture entries to `refusal_route_fixtures_test.go` first, so the file does not grow.

## Acceptance

- [ ] A complete checkpoint whose running executable differs from the broker execs the broker with the same arguments and working directory.
- [ ] A complete checkpoint whose running executable equals the broker runs in its own process.
- [ ] An unreadable manifest and a broker whose digest differs from its manifest each refuse with `checkpoint-broker-unavailable` and run no gate.
- [ ] A repository that declares no build inputs runs the complete checkpoint under its running executable.
- [ ] A chunk checkpoint runs under its running executable when it differs from the broker.
- [ ] `TestCheckpointFacesFollowTheirRoutes` produces and follows the new face.
