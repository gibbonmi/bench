# Share the operand resolution between bench anchors and bench probe

Blocked by: none
Writes: internal/canonicalpath/canonicalpath.go, internal/canonicalpath/canonicalpath_test.go, cmd/bench/anchors_command.go, internal/probe/subject.go
Covers: none

## What to build

`anchorQueryPath` in `cmd/bench/anchors_command.go` and `resolveSubject` in
`internal/probe/subject.go` each derive the same fact. An absolute operand
stays. A relative operand joins onto the working directory. The result is
cleaned, and the display spelling is the root-relative path with forward
slashes. That is two sources for one rule, which the code standard names a
defect.

`internal/canonicalpath` owns the canonical-path derivation, and its purity
census forbids `os.Getwd`. So the shared function is pure. It takes the root,
the working directory, and the operand. It returns the cleaned absolute path
and the display spelling. When the root cannot relativize the path, the display
keeps the cleaned operand with forward slashes.

Both callers read the working directory themselves and pass it in. The anchors verb keeps its existence
branch, which returns the cleaned operand for an absent path. The probe verb
keeps every refusal and its `withinRoot` guard.

## Acceptance

- [ ] `canonicalpath` exports one operand resolver, and a table test covers the absolute operand, the relative operand, the `..` operand, and the operand outside the root.
- [ ] `anchorQueryPath` and `resolveSubject` call that resolver, and neither file joins the working directory or relativizes the root itself.
- [ ] `bench anchors .agents/skills/bench-craft-spec/SKILL.md` from the worktree root and `bench anchors ../bench-craft-spec/SKILL.md` from inside `.agents/skills/bench-craft-tickets` print the same rows.
- [ ] Every test under `internal/probe`, `internal/canonicalpath`, and `cmd/bench` stays green, and `bench test --check canonical-path-owner` exits 0.
