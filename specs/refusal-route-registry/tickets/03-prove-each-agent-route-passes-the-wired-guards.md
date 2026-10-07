# Prove each agent route passes the wired guards

Blocked by: 02-move-the-landing-faces-into-the-shared-registry.md
Writes: internal/conformance/refusal_route_guard_test.go (new), internal/refusalroute/registry.go (new), internal/refusalroute/route.go (new), internal/worktree/land_refusal.go, internal/worktree/land_rerun.go, internal/worktree/land_resume.go, internal/worktree/refusal_route_test.go (new), internal/worktree/identity_component_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/worktree/parallel_census_test.go, internal/worktree/land_surface_test.go, internal/worktree/land_resume_test.go, internal/worktree/land_reauthorization_test.go, internal/worktree/refusal_route_follow_test.go (new)
Covers: RR08, RR09, RR10, RR11, RR12, RR15, RR58, RR60

## What to build

No landing agent route runs a Bench child through `bench worktree exec`.
For a source path that is not line-safe, the route prints a `bench worktree path <id>` step and then the `<checkout>` placeholder.
`atSourceWorktree` and `landingResumeNext` take this form.
So the `land-incomplete` resume and each preflight re-run obey it.

Add the guard conformance test in `internal/conformance`.
It renders each agent route of the registry with sample facts.
It fills every slot with a sample value, operator slots included.
A path fact takes a sample under the pool prefix.
The test classifies each command step through `gitguard.Classify`, `benchguard.PoolReference`, and `benchguard.Classify`.
It also refuses a step that runs `bench worktree exec <target> -- bench ...`.

The production guards do not change.
The test bites: an injected agent face whose step is `git merge <commit>` turns it red.

Extend each land producing fixture with a follow step.
The follow step carries out the printed route, runs each printed agent command step verbatim, and reruns the landing.
Then the face's sentence no longer prints.

## Acceptance

- [ ] Every agent route command step, with sample facts, gets the empty label from `gitguard.Classify`, `Blocked == false` from `benchguard.Classify`, and the empty string from `benchguard.PoolReference`.
- [ ] The guard check refuses an agent step that runs `bench worktree exec <target> -- bench <verb>`.
- [ ] An injected agent face whose route step is `git merge <commit>` makes the guard check fail.
- [ ] Each land fixture follows its printed route, reruns the landing, and the face's sentence no longer prints.
- [ ] A landing refusal and an incomplete landing on a source path that is not line-safe print a route that contains `<checkout>` and no `bench worktree exec`.
