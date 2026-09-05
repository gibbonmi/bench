# Add the git-plumbing-owner conformance check

Blocked by: migrate-the-worktree-directory-sites.md, migrate-the-worktree-file-sites.md, migrate-the-gate-and-dashboard-sites.md, migrate-the-diff-index-identity.md, refuse-an-unresolved-hooks-directory.md
Writes: internal/conformance/git_plumbing_owner_test.go (new), internal/conformance/checks_test.go, internal/conformance/registry/registry.go, internal/conformance/registry_test.go, tests/canary/package-core-guard/git-flag-retyped/ (new), projects/benchkit.md, tests/canary/guidance-prose-budgets/over-budget-skill, tests/canary/line-routing/line-binding-prose-drift, tests/canary/workflow-guidance-anchors/benchkit-hostile-input-heading, tests/canary/workflow-guidance-anchors/benchkit-review-round-owner, tests/canary/workflow-guidance-anchors/benchkit-review-round-routing, tests/canary/workflow-guidance-anchors/benchkit-spec-ownership, tests/canary/workflow-guidance-anchors/benchkit-system-suite-route
Covers: GR10, GR27, GR28, GR29, GR30, GR31, GR32, GR33, GR37

## What to build

Verify the premise first. Read `checkCanonicalPathOwner` in
internal/conformance/canonical_path_owner_test.go as the closest precedent. Read
`conformanceChecks` at internal/conformance/checks_test.go line 45. Read `Checks`
and `familyChecks` in internal/conformance/registry/registry.go. Read
`InputGoSource` in the same file. Read `canaryFixtureRegistry` and
`fixtureRegistrationFor` in internal/conformance/registry_test.go.

Read `TestEveryRetainedFixtureBitesThroughRegisteredOwner` in
internal/conformance/fixture_bite_test.go. Read
`TestCanaryFixtureRegistryClassifiesEveryFixture` in
internal/conformance/registry_test.go. Read the `BASE`, `CHECK`, `EXPECT`, and
`MUTATE.json` files of tests/canary/package-core-guard/bounds-duplicate-owner as
the fixture shape. Read the profile lane table near projects/benchkit.md line 428.

Read the three survivors the rule must pass: internal/gittest/gittest.go line
91, internal/gitguard/scan.go lines 26 and 104, and internal/benchguard/benchguard.go
line 265. The harness helper sits under the skipped tree. The two guard tables
hold a flag literal in no `rev-parse` call. Run the survivor sweep again after
every migration ticket lands. Stop and report to the coordinator when the check
reds a file this ticket does not name.

Add the check `checkGitPlumbingOwner` in
internal/conformance/git_plumbing_owner_test.go. Parse each non-test Go file
outside internal/git and internal/gittest. The four flags are `--git-dir`,
`--absolute-git-dir`, `--git-path`, and `--git-common-dir`. Red a string literal
that equals a flag inside a call whose arguments also hold the literal `rev-parse`.
Pass a flag that sits inside a longer literal. Pass a flag literal in a map literal
or in a call with no `rev-parse` argument.

Emit the diagnostic `<file> spells the Git administration flag <flag> outside
internal/git`. The seven fixture directories in the `Writes:` line are closure
headroom for the profile pins. This ticket edits none of them.

Join the check to its four advertisements. Add the registry row with the input
class `go-source`. Add the binding in `conformanceChecks`. Add the profile table
row in projects/benchkit.md. Add the fixture classification in
`canaryFixtureRegistry`.

Add the canary fixture `git-flag-retyped` under
tests/canary/package-core-guard/. Name `git-plumbing-owner` in its `CHECK` file.
In its mutation, replace the dashboard's reader call with the old
`git.Output("-C", root, "rev-parse", "--absolute-git-dir")` call, so the planted
literal sits inside a `rev-parse` call. Expect the planted diagnostic in its
`EXPECT` file.

Write the four unit tests in internal/conformance/git_plumbing_owner_test.go.
Name them `TestGitPlumbingOwnerRedsARetypedFlag`,
`TestGitPlumbingOwnerSkipsTestsAndTheAdapter`,
`TestGitPlumbingOwnerToleratesEmbeddedFlag`, and
`TestGitPlumbingOwnerIgnoresCallsWithoutRevParse`. Drive each over a temporary
tree. The red test runs one case per flag, each inside a `rev-parse` call. The
skip test also plants a `rev-parse` call under internal/gittest.

## Acceptance

- [ ] The check returns no diagnostic over the landing tree.
- [ ] For each of the four flags, the check returns one diagnostic that names the file and the flag over a retyped `rev-parse` call.
- [ ] `bench test --help` lists `git-plumbing-owner`, and `TestNamedCheckHelpListsEverySupportedCheck` passes.
- [ ] The check returns no diagnostic for a test file, a file under internal/git, or a file under internal/gittest.
- [ ] The check returns no diagnostic for a flag inside a longer literal.
- [ ] The check returns no diagnostic for a flag literal in a map literal or in a call with no `rev-parse` argument.
- [ ] The fixture `git-flag-retyped` reds the check with its `EXPECT` diagnostic.
- [ ] The fixture goes green again on restore.
- [ ] The registry, the check map, and the profile table agree on `git-plumbing-owner`.
- [ ] The fixture map classifies `git-flag-retyped`.
- [ ] `bench test --check git-plumbing-owner` runs the check alone and prints one verdict.
- [ ] Self-probe: point the fixture mutation at a test file, and report the fixture-bite test red.
