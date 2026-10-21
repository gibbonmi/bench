# Migrate the Bench and Write envelope readers

Blocked by: 2-add-envelope-owner-for-git.md
Writes: internal/benchguard/benchguard.go, internal/benchguard/projection_contract_test.go (new), internal/writeguard/writeguard.go, internal/writeguard/envelope_contract_test.go (new), cmd/bench/guard_projection_test.go (new)
Covers: GP22, GP23, GP24, GP25, GP26, GP27, GP28, GP29, GP48, GP50, GP52, GP53, GP55, GP56, GP58, GP59, GP61

## What to build

Move the Bench and Write readers onto the envelope owner, and keep each reader's required fields, error text, and failure posture.
Ticket 2 supplies the reviewed `DecodeEnvelope` owner with all three forms, and `cmd/bench/guard_projection_test.go`, which this ticket extends.

`benchguard.CommandFromEnvelope` selects `InputMap` and keeps its error signature and its missing, type, empty, and NUL error text.
Change only `CommandFromEnvelope` and its `encoding/json` import in `benchguard.go`; ticket 4 migrates its projection consumers.
`benchguard.go` has 13 lines of headroom and must end at 400 lines or fewer.

`writeguard.PathFromEnvelope` selects `CwdInputMap` and keeps its outer cwd validation before path handling, including for an absolute path, its path cleanup, and its relative-path rule.
Neither reader begins to validate a field it does not validate today.
Each new file stays within 400 lines.

## Acceptance

- [ ] Bench missing, non-string, empty, and NUL commands keep their error text, and the Bench core keeps its warning-allow output for family E0 (GP22–GP25).
- [ ] `TestPathFromEnvelope` and `TestPathFromEnvelopeRefusesAnUnreadablePath` stay unchanged and pass, and an invalid cwd type stays unreadable for an absolute path (GP26–GP28).
- [ ] The file-write core keeps its warning-allow output for unreadable inputs (GP29).
- [ ] The Bench and Write rows of the decoder preservation matrix give their stated results at their stated seams (GP48, GP50, GP52, GP53, GP55, GP56, GP58, GP59, GP61).
- [ ] Replacing the Bench caller's `InputMap` selection with `CommandStruct` makes `TestBenchEnvelopeRequiredField` or `TestBenchEnvelopeValue` red, and restoration returns the same green.
- [ ] Replacing the Write caller's `CwdInputMap` selection with `InputMap` makes `TestGuardWriteEnvelopePolicy` red, and restoration returns the same green.
- [ ] No reader outside `internal/shellcommand` calls `json.Unmarshal` on a hook envelope.
