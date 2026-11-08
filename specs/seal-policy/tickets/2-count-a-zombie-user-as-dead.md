# Count a zombie actor and a zombie-only group as dead on Linux

Blocked by: 1-grade-before-the-cache-release-seals.md
Writes: internal/subprocess/protection/store.go, internal/subprocess/protection/root_linux.go, internal/subprocess/protection/root_darwin.go, internal/subprocess/protection/root_other.go, internal/subprocess/protection/observe.go, internal/subprocess/protection/observe_test.go, internal/subprocess/protection/observe_linux_test.go (new), internal/subprocess/protection/registration.go, internal/subprocess/protection/recovery_test.go, internal/conformance/injected_ports_registry_test.go, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, docs/known-issues.md, docs/adr/0030-one-lifetime-owner-and-process-user-protection.md, CHANGELOG.md
Covers: SP11, SP12, SP13, SP14, SP15, SP16, SP17, SP18, SP19, SP20, SP33, SP35

## What to build

Review chunk: SP-C2.

Change the default of `Protocol.observer` (`internal/subprocess/protection/store.go:49-54`) from `syscall.Kill` to a protection-owned probe. Put the Linux probe and its proc reader in `root_linux.go`. Give `root_darwin.go` and `root_other.go` the plain `syscall.Kill` probe. Add no other production file, because the package directory already holds 18 source files against a limit of 12.

On Linux, the probe sends the signal first. A non-nil answer returns unchanged. After a nil answer to signal zero for a pid, the reader reads `<root>/<pid>/task/*/stat`, and the pid counts as absent (ESRCH) only when every task that it finds is in state `Z`. After a nil answer to signal zero for a group, the reader scans `<root>/*/stat` for each process whose process group is that group. The group counts as absent only when no member counts as present under the pid rule. The reader parses the state as the first field after the last `)`, and the process group as the third field. A record that it cannot read or parse keeps the nil answer. An entry that is gone (ENOENT) decides nothing. The probe never sends a second signal. The reader takes its root as a parameter; production passes `/proc`.

Keep `ObserveGroup`, `Classify`, and the lifetime controller on the signal-zero answer. `Registration.Finish` (`registration.go:146`) reads the new default too; change its comment if it states the old rule.

Remove the FT482 row from `docs/known-issues.md`. Add a Forbid entry with the bare needle `FT482` on the known-issues file to `internal/anchors/registry_ticket_passes.go`, and update its comment list. Add the independent forbidden rule to `TestTicketSlicingAnchors` with this exact row text as its needle:

```text
| A live descendant of a `bench test` run seals the cache generation of its HOME when its CLI releases the cache. Each later `bench test` on that HOME then refuses with `the resource generation is sealed`. Under the gate, a system test that waits for a child mark can show this refusal as a 20 s timeout. | Stop the leaked descendant or wait until it exits, then run the command again. If a system test times out on a child mark, run that test alone before you run the gate again. | FT482: a cache release applies the decided seal policy for a live, leaked, or dead protected user. |
```

Record the red of that rule before you delete the row. Add one sentence to the platform scope of ADR 0030: on Linux, a process whose tasks are all zombies counts as absent. Add one `CHANGELOG.md` entry.

## Acceptance

- [ ] A release observation of a known user whose actor and group leader are one unreaped zombie retires the generation (SP11).
- [ ] A pending user whose actor is an unreaped zombie grades safe (SP12).
- [ ] A user whose zombie group leader has one live group member grades protected (SP13).
- [ ] The reader counts a process present when its leader task is `Z` and another task is `S` (SP14).
- [ ] The reader reads state `S` from a stat record whose command name is `a) Z 1 (b` (SP15).
- [ ] The reader keeps a pid present when its stat record is unreadable or malformed (SP16).
- [ ] The reader counts a group absent when its only entries are one zombie member and one gone member (SP17).
- [ ] `TestRecoveryPendingActorAbsent`, `TestRecoveryWorkerActorDeath`, and `TestRecoveryNoStaleSignal` stay green (SP18, SP19, SP20).
- [ ] `docs/known-issues.md` holds no FT482 row, and `TestTicketSlicingAnchors` reds when the row returns (SP33).
- [ ] ADR 0030 states the Linux zombie rule (SP35).

## Verification

Run `bench test --package ./internal/subprocess/protection` and `bench test --package ./internal/anchors` whole. The real-zombie tests start a child with `Setpgid`, do not reap it until cleanup, and wait until its stat record shows `Z`. Observe each planned row red before its production edit. Self-probe: make the reader read only the leader task state, and report that SP14 goes red.
