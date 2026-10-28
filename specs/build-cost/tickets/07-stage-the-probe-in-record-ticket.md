# Stage the retained probe in the ticket record

Blocked by: 01-accept-probe-only-plans-and-retire-check-floor.md, 02-refuse-plain-chunk-requirements.md, 03-state-the-probe-only-rule.md, 06-pass-the-chunk-base-to-the-commit-chain.md
Writes: internal/reviewrecord/recordcmd/ticket.go, internal/reviewrecord/recordcmd/stage.go, internal/reviewrecord/recordcmd/ticket_probe_test.go (new), internal/reviewrecord/recordcmd/command_test.go, cmd/bench/help_inventory_test.go, cmd/bench/record_operations_test.go, .agents/commands/bench-implement-spec.md, .agents/commands/bench-review-implementation.md, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go
Covers: BC41, BC42, BC43, BC44, BC45, BC46, BC47, BC48

## What to build

`bench record ticket` stages the retained response of the worktree when exactly one owed requirement has no stage and the retained call is `bench probe`. The guard line in `ticket.go` reads `if len(unstaged) == 1 && retained.Verb == probe.Verb {`, because the probe-07 command swaps that exact line.
The stage applies the checks of the excerpt form: the bound, the probe import, and `--mutant`. Move those checks into one function in `stage.go` that both forms call, so that the rule has one source.

`bench record ticket` takes `--mutant <n>` with the meaning of the excerpt form. A retained call that is not `bench probe` refuses and names `bench record excerpt`. Two or more unstaged owed requirements keep the current refusal. Each refusal leaves the record bytes unchanged.
After its table, the verb prints one line: `next=bench commit -m <msg> --preflight-build <slug> --base <chunk-base> reviews/<slug>.md`, with the slug and the chunk base of the written entry. The whole response fits in 10 lines.

The Build section of `.agents/commands/bench-implement-spec.md` replaces the exec preflight sentence with: "After the ticket commit, run the probe, then `bench record ticket`, then the record commit that its `next=` line prints."
The Land section states that the ticket form stages the last probe response, so that the author runs `bench record excerpt` only for an earlier probe of the same ticket. `.agents/commands/bench-review-implementation.md` states the same staging rule in its record list. Each edit to `bench-implement-spec.md` stays inside its paragraph, so the file keeps its 97-line budget.
`internal/anchors/registry_chunk_chain.go` replaces its exec preflight Require entry with the new sentence, and it adds the Forbid entry of BC48.

## Acceptance

- [ ] One unstaged probe requirement and a retained `bench probe` response give the entry with the response bytes (BC41).
- [ ] `--mutant 2` over a mutant-set response records row 2 (BC42).
- [ ] A retained `bench test` response refuses, names the excerpt form, and leaves the record unchanged (BC43).
- [ ] Two unstaged probe requirements refuse and leave the record unchanged (BC44).
- [ ] The `next=` line names the record commit with the slug and the chunk base (BC45).
- [ ] A call that writes three entries prints at most 10 lines (BC46).
- [ ] The Build section holds the new sequence sentence and no exec preflight step (BC47, BC48).
