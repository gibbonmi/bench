# Repair hostile ownership-fence grammar coverage

Blocked by: none
Ownership fence: `.agents/skills/bench-craft-spec/SKILL.md`, `internal/anchors/registry_data.go`, `internal/conformance/fixture_bite_test.go`, `tests/canary/workflow-guidance-anchors`
Integration surfaces: ownership-fence grammar→`internal/anchors/registry_data.go`; fixture inventory→`internal/conformance/fixture_bite_test.go`; hostile grammar mutations→`tests/canary/workflow-guidance-anchors`
Contracts: the approved ownership-fence grammar crosses `.agents/skills/bench-craft-spec/SKILL.md`→the registered workflow fixtures; type is Markdown path declarations, domain is exact repo-relative files or path prefixes, order is section-local, absence or an invalid/glob-shaped declaration leaves the fence incomplete
Closure: RH1/exact-literal-fence, RH2/empty-or-invalid-fence

## What to build

Close accepted review finding `COV-001-hostile-fence-grammar` by registering section-sensitive checks for both independent ownership-fence predicates and adding auto-discovered fixtures that fail when either predicate is weakened. Keep `craft-spec` as the normative owner; add no parser or lifecycle authority.

## Acceptance

- [ ] [RH1] (covers SH3; repairs COV-001 exact/glob/out-of-repo member) an anchor and fixture require an exact repo-relative file or path prefix and preserve the rule that a glob is never a fence declaration.
- [ ] [RH2] (covers SH3; repairs COV-001 empty/invalid member) an independently failing anchor and fixture require an empty or invalid fence section to remain incomplete.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RH1/exact-literal-fence | relocate the exact repo-relative and never-a-glob predicate outside the ownership-fence rule section while leaving its words elsewhere | the registered section-sensitive workflow anchor | apply the relocation, run the focused workflow-anchor conformance test, require the exact-literal-fence diagnostic, restore the subject |
| RH2/empty-or-invalid-fence | replace the incomplete-fence predicate with language that permits an empty or invalid fence | the registered section-sensitive workflow anchor | apply the subject swap, run the focused workflow-anchor conformance test, require the empty-or-invalid-fence diagnostic, restore the subject |
