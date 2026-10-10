# Give the adopt reasons and the hook refusal one source

Blocked by: none
Writes: internal/adopt/broker.go, internal/adopt/link_hook.go, internal/adopt/link_transaction.go, internal/adopt/manifest.go, internal/adopt/lifecycle_special_test.go
Covers: none

## What to build

The lifecycle reason words and the foreign-hook refusal text each get one production source. Today the repair checks type the reason words as literals, although the lifecycle decision owns them as constants. The shim check types both words, and the doctor-asset, pre-push, and manifest checks type the modified-managed word. The repair checks read the decision's reason constants instead.

Today two callers also type the foreign-hook refusal text: the link transaction and the hook installer. The hook domain helper gets one producer of that text, and both callers print its result. The lifecycle hook test compares the link output with that producer.

Behavior does not change. Each command prints the same reason words and the same refusal text on the same stream, and it exits with the same code.

## Acceptance

- [ ] No production file in the adopt package types the modified-managed or project-owned reason word outside the decision constants. A sweep of the adopt package shows the decision constants as the only production source.
- [ ] The foreign-hook refusal text has one producer in the hook domain helper. The link transaction, the hook installer, and the lifecycle hook test read that producer.
- [ ] The lifecycle hook test fails when a caller prints a refusal text that differs from the producer. A mutation probe records that red.
- [ ] The test expectations that type the reason words stay independent. A mutation probe of a reason constant records their red.
- [ ] The adopt package tests and the repair tests stay green, and the author reports their elapsed time.
