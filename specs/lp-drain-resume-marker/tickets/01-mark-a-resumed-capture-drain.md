# Mark a resumed capture drain

Blocked by: none
Writes: internal/capturetx/public.go, internal/capturetx/recovery_test.go, internal/roadmap/capture.go, internal/roadmap/capture_test.go, CHANGELOG.md
Covers: none

## What to build

Source: drain d-359accf567c1, learning L2.

On 2026-10-16, `bench capture drain` found a sealed drain that an earlier session did not commit.
The command printed the same drain ID with the state `sealed`, and it gave no sign that the drain was an old one.
Thus the new session read a stale snapshot.

When `bench capture drain` finds an open drain, the command resumes that drain.
Make the output of a resumed drain state the resume and the time that the drain was sealed.
Add a `resumed[1]{sealed_at}:` block after the `drain` block, only for a resumed drain.
The `sealed_at` cell is an RFC 3339 UTC time.

The drain manifest records no seal time.
The last manifest write of an open drain is the write that sets the state `sealed`, so the modification time of the manifest file is the seal time.
Use that recorded file time, and add no persisted field.

The output of a fresh drain does not change.
The output of `bench capture drain show <drain-id>` does not change.
This ticket does not change `bench status`.

## Acceptance

- [ ] A repeated `bench capture drain` prints the `resumed[1]{sealed_at}:` block with the modification time of the drain manifest.
- [ ] A fresh `bench capture drain` prints the `drain` block and the `sources` block, and no `resumed` block.
- [ ] `capturetx.Begin` reports a resumed generation and its seal time, and it reports a fresh generation as not resumed.
