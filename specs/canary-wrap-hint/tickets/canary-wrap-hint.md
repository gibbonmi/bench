# Name the line wrap when a canary mutation anchor is found only under collapsed whitespace

Blocked by: none
Writes: internal/canary/mutation.go, the internal/canary test file that covers the exactly-once anchor refusal

## What to build

The materializer in `internal/canary/mutation.go` refuses a fixture whose
`mutation.Old` text does not occur exactly once in the target body. The refusal
message is the exactly-once message quoted in the acceptance rows. The anchor
evaluator in `internal/anchors/match.go` matches under collapsed whitespace. So
a needle that wraps across a line in the target file passes the evaluator, and
the materializer refuses it with no hint.

Keep the byte-exact match. Do not collapse whitespace in the materializer,
because the materializer cannot reconstruct the bytes it would rewrite. When the
byte count is zero, compare again under collapsed whitespace on both sides. If
that collapsed comparison finds the needle, extend the refusal.

The extended refusal says the anchor spans a line wrap in the target file. It also says the
fixture must quote the physical line bytes. When the collapsed comparison also
misses, or when the count is greater than one, keep the current message.

Read the collapse helper in `internal/anchors/match.go` before you call it.
Make sure the new import from the canary package to the anchors package creates
no cycle. If it does, stop and report.

## Acceptance

- [ ] A fixture whose `Old` text matches the target only after whitespace collapse is refused with a message that names the line wrap.
- [ ] A fixture whose `Old` text matches nowhere keeps the existing `did not occur exactly once` message with no wrap hint.
- [ ] A fixture whose `Old` text occurs twice keeps the existing message with no wrap hint.
- [ ] Removal of the collapsed comparison turns the new test red.
- [ ] `go vet` and `gofmt` are clean on the edited files.
