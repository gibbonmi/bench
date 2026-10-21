# Migrate guard envelope readers

Blocked by: 1-share-git-command-projections.md
Writes: internal/shellcommand/envelope.go (new), internal/shellcommand/envelope_test.go (new), internal/gitguard/gitguard.go, internal/gitguard/projection_contract_test.go (new), internal/benchguard/benchguard.go, internal/benchguard/projection_contract_test.go (new), internal/writeguard/writeguard.go, internal/writeguard/envelope_contract_test.go (new), cmd/bench/guard_projection_test.go (new), tests/canary/canonical-path-owner/second-derivation/
Covers: GP14, GP15, GP16, GP17, GP18, GP19, GP20, GP21, GP22, GP23, GP24, GP25, GP26, GP27, GP28, GP29, GP43, GP47, GP48, GP49, GP50, GP51, GP52, GP53, GP54, GP55, GP56, GP57, GP58, GP59, GP60, GP61

## What to build

Move the three hook-envelope decoders onto one structural owner, and keep each guard's required fields, errors, and failure posture.
Add `DecodeEnvelope`, the three `EnvelopeForm` values, and the `Command`, `Cwd`, and `InputString` accessors in a new `internal/shellcommand/envelope.go`, as the spec's Envelope API states.
The owner calls `encoding/json.Unmarshal` with the original anonymous target shapes and tags, so folded keys, merges, null rules, and early type errors stay as they are.
It does not require a field, resolve a path, inspect a repository, or choose an exit.

`gitguard.CommandFromEnvelope` selects `CommandStruct` and keeps its string-only signature and its empty result on any decode error.
`benchguard.CommandFromEnvelope` selects `InputMap` and keeps its error signature and its missing, type, empty, and NUL error text.
`writeguard.PathFromEnvelope` selects `CwdInputMap` and keeps its outer cwd validation before path handling, its path cleanup, and its relative-path rule.
Change only `CommandFromEnvelope` and its `encoding/json` import in `benchguard.go`; ticket 3 migrates its projection consumers.
`benchguard.go` has 13 lines of headroom and must end at 400 lines or fewer.

`gitguard.go` loses its `encoding/json` import, which the canonical-path-owner canary anchors on.
Re-anchor `tests/canary/canonical-path-owner/second-derivation/MUTATE.json` to a surviving byte sequence, and keep its expected diagnostic.
Keep `type Checker struct` in `gitguard.go` byte for byte, because the injected-ports canary anchors on it.
Each new file stays within 400 lines. The core entry tests go in the new `cmd/bench/guard_projection_test.go`; do not edit `cmd/bench/main_test.go`, which is over its limit with no grant.

## Acceptance

- [ ] Invalid JSON, a non-object root, a non-object `tool_input`, absent and null objects, an unrelated invalid cwd, and a first null string give the decode results that GP14–GP19 state.
- [ ] Git unreadable, missing, null, and empty commands return an empty string, and the Git core stays at exit 0 with empty stderr for family E0 (GP20, GP21).
- [ ] Bench missing, non-string, empty, and NUL commands keep their error text, and the Bench core keeps its warning-allow output (GP22–GP25).
- [ ] `TestPathFromEnvelope` and `TestPathFromEnvelopeRefusesAnUnreadablePath` stay unchanged and pass, and an invalid cwd type stays unreadable for an absolute path (GP26–GP28).
- [ ] The file-write core keeps its warning-allow output for unreadable inputs (GP29).
- [ ] Each row of the decoder preservation matrix gives its stated result at its stated seam (GP47–GP61).
- [ ] Both gitguard-pinned canary families keep biting through their registered owners (GP43).
- [ ] Swapping one caller's envelope form for another makes that caller's named witness red, for Git, Bench, and Write separately.
- [ ] No reader outside `internal/shellcommand` calls `json.Unmarshal` on a hook envelope.
