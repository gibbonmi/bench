# 5. Fault the landing follow-on steps with real fixtures

Blocked by: 4-interrupt-the-landing-marker-with-a-gate-script.md
Writes: internal/worktree/land_marker_fixture_test.go (new), internal/worktree/land_resume_test.go, internal/worktree/land_freshness_test.go, internal/worktree/land_prunes_landed_siblings_test.go, internal/worktree/land_release_refusal_test.go, internal/worktree/land_fixtures_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: WS34, WS35, WS36, WS37, WS38

## What to build

Chunk: SR-C4.

Convert the reconcile test to a nested repository in the destination, which fails the
real residue guard. Convert the post-publication freshness test to the marker fixture of
ticket 4 and the nested-repository fixture. Convert the two release-refusal tests to the
public landing fixture, so that the real source authority supplies the fences.

The prune test depends on a probe. Plant a stale `refs/heads/<name>.lock` for a landed
sibling, then run `bench probe` on `land.go` with the prune error branch omitted. If the
test turns red, keep the conversion. If the probe stays green, keep
`pruneLandedBranches` and its test unchanged, and record one `bench learning` entry.

Each converted test keeps its name. Record each probe command and its red in the
verification note.

## Acceptance

- [ ] A nested destination repository leaves the published checkout unreconciled, and the resume reconciles it.
- [ ] Each post-publication failure resumes without a second publication.
- [ ] A stale branch lock makes the landing report an incomplete prune, or the learning entry records the failed probe.
- [ ] The two release-refusal tests land through the public landing fixture.
- [ ] No test in this ticket sets `reconcileLanding`, `authorizeLandingSource`, or a converted probe field.
