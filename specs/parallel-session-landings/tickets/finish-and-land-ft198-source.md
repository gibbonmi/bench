# Finish and land the complete FT198 source

Blocked by: restore-ft198-doctrine-commit.md
Writes: cmd/bench/main.go, cmd/bench/command_registry_test.go, internal/conformance/axi_query_registry_test.go, .agents/skills/bench-craft-cli/SKILL.md, .claude/skills/bench-craft-cli/SKILL.md, projects/benchkit.md, internal/roadmap

## What to build

Starting from the restored doctrine commit in the retained FT198 integration
source, implement its remaining `join-axi-set.md` ticket at that source's current
tip. Preserve the existing FT198 ticket contract: join `bench roadmap` to the
approved AXI set, pin its complete-value disclosure, and keep every PI row green.
Receive the exact ephemeral request token returned by the blocker and reuse it
unchanged as this ticket's `bench worktree land --request` input; do not persist
or replace it while it remains available. If it was lost at the ticket/session
boundary, rerun `bench worktree reauthorize` under the same explicit reviewer
authorization with the current FT198 source tip and a new token; never recover by
persisting a token.

After the generic parallel-session-landings implementation and tests are green,
reuse the blocker-built parallel-session-landings integration source
`./dist/bench` by absolute path, never a PATH `bench`. Invoke it with the clean
retained FT198 source worktree as cwd to review the complete
`0924e02e..<FT198-source-tip>` range. Then invoke the same binary from the clean
`main` checkout for `bench worktree land --spec roadmap-progressive-index` with
the handed-on or explicitly reminted request token. Verify the merge parents,
prospective tree, project-green marker, destination reconciliation, and
assignment release.

Do not rewrite `main`, rerun the first three FT198 tickets, alter the restored
doctrine bytes, absorb ambient handoff state, or touch either foreign assignment.
This is the actual remaining FT198 implementation and landing, not a synthetic
verification-only change.

## Acceptance

- [ ] The completed FT198 source satisfies its remaining PI15 and PI19 AXI
      contract while preserving the restored doctrine commit and every other
      roadmap-progressive-index row.
- [ ] The landing authenticates with the exact ephemeral token returned by the
      blocker or, after documented loss, one token reminted by the same reviewer-
      authorized reauthorization path; no token is committed or written into a
      handoff file.
- [ ] Explicit-base review includes the complete FT198 patch from `0924e02e` and
      excludes destination-only phase history when the absolute-path feature
      binary runs with the retained source as cwd; the retained doctrine path
      hashes remain unchanged (covers PL2, PL24).
- [ ] One public `--spec roadmap-progressive-index` landing publishes the reviewed
      FT198 source when the same binary runs from clean `main`, without rewriting
      `main` or replaying its first three commits, then closes only its own
      assignment with the expected marker and checkout state (covers PL24).
