# Add the envelope owner and migrate the Git reader

Blocked by: 1-share-git-command-projections.md
Writes: internal/shellcommand/envelope.go (new), internal/shellcommand/envelope_test.go (new), internal/gitguard/gitguard.go, internal/gitguard/projection_contract_test.go (new), cmd/bench/guard_projection_test.go (new), tests/canary/canonical-path-owner/second-derivation/
Covers: GP14, GP15, GP16, GP17, GP18, GP19, GP20, GP21, GP43, GP47, GP49, GP51, GP54, GP57, GP60

## What to build

Add the one structural owner for hook envelopes, and move the Git reader onto it with its current posture.
Add `DecodeEnvelope`, all three `EnvelopeForm` values, and the `Command`, `Cwd`, and `InputString` accessors in a new `internal/shellcommand/envelope.go`, as the spec's Envelope API states.
The owner calls `encoding/json.Unmarshal` with the original anonymous target shapes and tags, so folded keys, merges, null rules, and early type errors stay as they are.
It does not require a field, resolve a path, inspect a repository, or choose an exit.
Value tests in `envelope_test.go` cover all three forms, so ticket 3 migrates its readers onto a reviewed owner.
Do not edit `shellcommand.go`, which has 3 lines of headroom. New shellcommand code must pass the `stdlib-reimplementation` check.

`gitguard.CommandFromEnvelope` selects `CommandStruct` and keeps its string-only signature and its empty result on any decode error.
Ticket 1 creates `gitguard/projection_contract_test.go`; this ticket extends it with the Git reader rows.
Create `cmd/bench/guard_projection_test.go` for the Git core entry rows; tickets 3, 4, and 6 extend it. Do not edit `cmd/bench/main_test.go`, which is over its limit with no grant.

`gitguard.go` loses its `encoding/json` import, which the canonical-path-owner canary anchors on.
Re-anchor `tests/canary/canonical-path-owner/second-derivation/MUTATE.json` to a surviving byte sequence, and keep its expected diagnostic.
Keep `type Checker struct` in `gitguard.go` byte for byte, because the injected-ports canary anchors on it.
Each new file stays within 400 lines.

## Acceptance

- [ ] Invalid JSON, a non-object root, a non-object `tool_input`, absent and null objects, an unrelated invalid cwd, a first null string, and repeated `tool_input` objects give the results that GP14–GP19 and GP60 state.
- [ ] Git unreadable, missing, null, and empty commands return an empty string, and the Git core stays at exit 0 with empty stderr for family E0 (GP20, GP21).
- [ ] The Git rows of the decoder preservation matrix give their stated results through the actual `guard-git` entry (GP47, GP49, GP51, GP54, GP57).
- [ ] Both gitguard-pinned canary families keep biting through their registered owners (GP43).
- [ ] Replacing the Git caller's `CommandStruct` selection with `InputMap` makes `TestGuardGitEnvelopePolicy` red, and restoration returns the same green.
- [ ] `gitguard` calls no `json.Unmarshal` on a hook envelope.
