# Lock the probe subject and re-read it before the baseline

Blocked by: 1-add-probe-mutation-forms.md, 2-probe-prose-and-system-checks.md, 3-name-the-first-failure-and-expected-witness.md, 4-name-a-mutant-build-failure.md
Writes: cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/gocache/lock.go, internal/gocache/lock_test.go, internal/probe/command.go, internal/probe/lock.go (new), internal/probe/lock_holder_test.go (new), internal/probe/lock_test.go (new), internal/probe/probe.go, internal/probe/refusal_test.go
Covers: PS58, PS59, PS60, PS61, PS62, PS63, PS64, PS65, PS66, PS67, PS68, PS75

## What to build

A second probe on one subject refuses while the first probe runs, and probes on two subjects run at the same time.
A killed probe never blocks its subject, and a completed probe leaves no lock entry under the Bench home.

Tickets 1 to 4 supply the forms, the subjects, and the verdict faces.
This ticket adds the subject lock to the run sequence in `probe.go` after those tickets commit.

Create `internal/probe/lock.go` as the owner of the subject lock and the re-read.
The lock file is the SHA-256 of the absolute subject path, inside `$BENCH_HOME/probe-locks/`.
The owner reads `BENCH_HOME` through `benchhome.Dir`, as `preserve` does.
The probe takes the lock after the gate-lock read and before the baseline.

The acquire protocol and the release protocol follow the spec's subject-lock decisions:

- The acquire makes at most three attempts.
- Each attempt opens the lock file read-write with create, no-follow, and close-on-exec at mode 0600.
- Each attempt takes an exclusive, non-blocking write lock through `gocache.RecordLock`.
- Each attempt compares the held descriptor's identity with the path's identity through `fstat` and `lstat`.
- The release removes the lock file, closes the descriptor, and then removes the empty lock directory.
- The deferred release also runs on an interrupt, on every refusal after the lock, and on the exit-2 answer that keeps a copy.
- The lock directory is a sibling of the preserve directory, so the release never removes a kept copy.
- Keep the lock calls in `lock.go`, because `probe.go` has 332 lines on `main` before tickets 1 to 4.

After the lock succeeds, the probe re-reads the subject through `resolveSubject`.
A changed byte or mode refuses with the subject-changed line, and an unusable path gives its current subject line.
Update the refusal order in the `Command` doc comment to match the spec.

Two package hooks give the lock rows their windows.
The start-read hook runs between the start read and the lock.
The lock-owner hook runs between the open and the identity check.

A lock row holds the lock from another process, as `TestProbeRefusesUnderALiveGateRun` holds the gate lock.
That holder lives at `refusal_test.go:93-106` and `:116`, and `refusal_test.go` stays on this line because the subject-lock holder may generalize it.
Put the subject-lock holder process and its hook helpers in `internal/probe/lock_holder_test.go`, so `lock_test.go` stays at 400 lines or fewer.
`TestProbeWritesNoRecord`, `TestProbeRefusesWhenPreservationFails`, and `TestProbePreservesTheSubjectBeforeTheRun` stay green without an edit.
The command registry files join this `Writes:` line through the binding registry, and this ticket does not change them.

## Acceptance

- [ ] With another process that holds the subject lock, a probe refuses with the in-use line and exit 1.
- [ ] That refusal leaves the subject bytes, the Bench home, and the go marker unchanged.
- [ ] With the gate lock and the subject lock both held by other processes, the gate-lock line wins.
- [ ] A subject that changed inside the start-read window refuses with the changed line, and the tree ends unchanged.
- [ ] A subject that the start-read hook deletes reads the absent-subject line, and no run child starts.
- [ ] A probe on a second subject runs while another process holds the lock of the first subject.
- [ ] After an interrupted probe, a new probe on the same subject reaches its verdict.
- [ ] A seeded lock file with no holder grants the lock, and the probe removes the file at its end.
- [ ] A lock file that the lock-owner hook replaces once gives a retry that holds the new file.
- [ ] Replacement on all three attempts gives the lock-unavailable line before the baseline.
- [ ] A regular file at `$BENCH_HOME/probe-locks` refuses with the lock-unavailable line and starts no run child.
- [ ] A completed probe leaves the record surfaces unchanged.
- [ ] Each new file stays at 400 lines or fewer.
