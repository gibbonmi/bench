# FT362 disposable prototype verdict

Question: Can one lifetime owner bound shutdown without losing stream or artifact ownership?

Verdict: Yes, with explicit stream ownership and an incomplete-cleanup outcome.

Evidence date: 2026-10-06 session

Runtime: Go 1.25.14, linux/amd64

Command: `go run /tmp/bench-ft362-lifetime-probe/main.go`

The coordinator ran the disposable program twice.
The second run added forced stream closure and cancellation after successful start.
Every assertion in the second run passed.
The program is discarded after this record.

## Observed results

- Buffer, caller writer, direct file, and owned pipe preserved `alpha\nbeta\n` on normal completion.
- Each successful start had one `Cmd.Wait` call.
- A ready child ignored TERM and required KILL after the 100 ms probe grace.
- The resistant child's waiter completed after escalation.
- A leader exited while its descendant kept the output pipe open.
- Killing that group allowed the separate JSON decoder to reach EOF.
- Closing an owned reader interrupted the decoder with a non-EOF error.
- A blocked caller writer kept `Cmd.Wait` pending after the child ended.
- The caller returned an incomplete result after the 50 ms grace and 150 ms final window.
- Releasing that writer allowed the same waiter to complete.
- A canceled pre-start context started no child.
- Cancellation after successful start completed teardown with one waiter.
- The concurrent-start sample recorded 100 successful starts and 100 waiters.

These durations demonstrate finite probe control flow.
They do not measure Bench performance or select production policy values.

## Modeled results and limits

The prototype injected ESRCH, EPERM, and a present-group result into absence classification.
Only ESRCH selected completed cleanup.
Artifact retention and publication refusal were modeled outcome decisions.
The prototype did not exercise production artifact owners.

The concurrent-start sample used ordinary successful children.
It does not qualify every scheduler interleaving or the current shift caller.
Production tests must cover a cancellation request racing a resistant child's startup.

macOS was not available.
The prototype does not qualify native macOS process behavior or Linux parent-death behavior.
It does not guarantee termination of a process in uninterruptible kernel sleep.
It does not guarantee completion of an arbitrary caller writer.

## Source interpretation

The coordinator read the local Go 1.25.14 os/exec source.
Cmd.Wait waits for non-file stream-copy goroutines at lines 897–949.
WaitDelay closes owned pipes, then still awaits those goroutines at lines 959–998.
StdoutPipe forbids Wait before reads complete at lines 1066–1087.
These rules explain the blocked-writer and independently owned pipe outcomes.

The coordinator read all eleven production group sites named by the map census.
Linux runbinary sets Pdeathsig to SIGKILL.
Darwin runbinary sets only Setpgid.
The spec must preserve that platform distinction and plan platform-specific tests.

## Evidence disposition

The independent Sol 6.1/xhigh review accepted the stream strategy as specification input.
It did not accept production artifact lifetime as proved.
Decision ticket 6 holds the resulting design constraints and future evidence obligations.
The disposable source was removed after the probe.

Prototype source SHA256: 278bcf2c97c1ba65c677201be7dfa23c6f6001b83bd61146628dcd78306a2f4a
