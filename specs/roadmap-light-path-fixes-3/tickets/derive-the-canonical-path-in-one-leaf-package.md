# Derive the canonical path in one leaf package

Blocked by: none
Writes: internal/canonicalpath (new), internal/worktree/subshell.go, internal/preflight/gather.go, internal/canary/mutation.go, internal/maps/validation.go, internal/runbinary/runbinary.go, internal/gate/subject.go, internal/conformance/canonical_path_owner_test.go (new), internal/conformance/checks_test.go, internal/conformance/registry/registry.go, tests/canary/canonical-path-owner (new), internal/conformance/registry_test.go, projects/benchkit.md, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: LQ1, LQ2, LQ3, LQ27

## What to build

Verify the premise first: `canonicalPath` in internal/worktree/subshell.go,
`canonicalRoot` in internal/preflight/gather.go, and `resolvePath` in
internal/canary/mutation.go each run `Abs`, `EvalSymlinks`, and `Clean`. Three
more production functions run the same pair: `validateSourcePath` in
internal/maps/validation.go, `canonicalSourceRoot` in
internal/runbinary/runbinary.go, and `canonicalSubjectRoot` in
internal/gate/subject.go.

Then add `internal/canonicalpath` with one exported
function that returns the absolute, symlink-resolved, cleaned form and keeps
the absolute spelling when the path does not exist. Mirror the
`internal/bounds` package doc. Make all six functions one-line wrappers over
it, and keep each caller's error posture. Add a purity census in the new
package that forbids every `internal/` import, mirroring
`internal/worktree/lifecyclepolicy/purity_census_test.go`. The share ticket
folds that census into the shared helper later.

Add one dev-tier conformance check that reds a production file outside
`internal/canonicalpath` that calls both `filepath.Abs` and
`filepath.EvalSymlinks`. Mirror `checkBoundCallers` in
internal/conformance/bounds_policy_test.go, test files excluded. Register it
after the marker-wait row, mirror the row in `checks_test.go`, and add a
canary fixture. The adopt `resolvedPath` has no `Abs` call, so it stays
green, and the bootstrap test fixture is a test file.

Three sibling tickets edit the maps, runbinary, and gate subject files after
this ticket, so this ticket lands first.

Do not touch `canonicalRoot` in internal/worktree/worktree.go, which wraps
`poolkey.Canonical` and answers a different question.

## Acceptance

- [ ] `TestGatherAssignmentTarget` and `TestRestoreMutationFixtureRefusesSymlinkedRootSpelling` pass through the shared owner.
- [ ] `TestRestoreMutationFixtureReinstatesBaseAndRemovesOverlay` passes with a not-yet-existing destination.
- [ ] The new package's purity census passes and reds a synthetic `internal/` import.
- [ ] The owner check passes on the live worktree, and its fixture reds a second `Abs` plus `EvalSymlinks` pair outside the package.
- [ ] The maps, runbinary, and gate subject packages' own tests pass through the shared owner.
- [ ] Self-probe: reduce the owner to `filepath.Clean`, and report the two symlink tests red.
