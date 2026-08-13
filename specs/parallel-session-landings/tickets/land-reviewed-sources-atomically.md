# Land a reviewed source as one authorized merge commit

Blocked by: expose-explicit-source-review.md, compose-reviewed-source-trees.md, reauthorize-retained-assignments.md
Writes: internal/diff, internal/landing, internal/worktree, internal/intent, internal/gate/authorization, internal/usage, cmd/bench/main.go, cmd/bench/command_registry_test.go, bin/bench.sh, internal/systemtest, internal/conformance/axi_query_registry_test.go, internal/conformance/package_shipped_surface_test.go

## What to build

Add the first-run `bench worktree land` operation from a clean checkout attached
to the resolved default branch. Authenticate the request and assignment, prove
the frozen reviewed pair, require a clean committed source, reauthorize its path
set against the staged spec's fences, validate identical staged spec bytes in
source and destination, and compose through the blocker ticket's Git-native
primitive. Apply the staged-to-implemented transition only in the prospective
tree and run the existing whole-project oracle on that exact immutable tree.

Only `authorization.Green` may publish. Create a two-parent commit over the
captured destination and reviewed source tip, win the destination with one
expected-old CAS, advance project-green in order, reconcile the destination
checkout, and release or retain the assignment according to existing residue
policy. Once the destination CAS wins, every later fault returns a structured
`landed` result with the published commit and incomplete step; it never rolls
back or republishes. The terminal envelope and exit partition are part of the
public command.

This is intentionally one atomic tracer. Landing tree transition, exact-tree
authorization, publication CAS, marker order, checkout reconciliation, and
assignment disposition cannot ship as independently safe partial commands.

## Acceptance

- [ ] Conflict repair requires a new reviewed tip; movement away from the exact
      reviewed tip, authentication, final source-fence authorization, clean-
      source/destination checks, default-branch attachment, and land grammar all
      refuse before gate when invalid; every composition conflict has a zero gate
      tally and preserves refs, indexes, worktrees, and merge-state files (covers
      PL4, PL8, PL9, PL18, PL19, PL26, PL27).
- [ ] The exact prospective tree is gated, identical spec bytes transition only
      there, every non-green authorization kind preserves all state, exact
      evidence reuse stays tree-and-oracle bound, and every invalid spec-state
      partition refuses before gate (covers PL10, PL11, PL12, PL13, PL33).
- [ ] The published object has exactly the captured destination and reviewed
      source parents, including on an FT198-shaped partly ancestral graph;
      destination races preserve the winner; a rerun recomposes; and absent,
      present, or concurrently moved project-green markers follow the specified
      order (covers PL7, PL14, PL15, PL16, PL17, PL28).
- [ ] Complete success exits 0 with the six-field structured result and a clean,
      released worktree; post-CAS marker, reconcile, or release faults exit 1 as
      `landed` with recoverable state; declared output and unknown ignored residue
      take their distinct terminal dispositions (covers PL20, PL21, PL31).
- [ ] Control-bearing source paths, hostile worktree paths and Git diagnostics,
      and special files fail boundedly without forged terminal lines or blocking
      reads (covers PL22).
