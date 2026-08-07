# State the repair-ticket rationale

Blocked by: none
Ownership fence: `.agents/commands/bench-implement-spec.md`, `specs/light-path-repair-ticket-rationale/tickets/state-repair-ticket-rationale.md`
Integration surfaces: repair-route rationale→`.agents/commands/bench-implement-spec.md`; Claude command consumer→`.claude/commands/bench-implement-spec.md` + RR2
Contracts: the canonical repair-route rationale crosses `.agents/commands/bench-implement-spec.md`→`.claude/commands/bench-implement-spec.md` through the tracked command-directory symlink, asserted by RR2 against the real resolved file; absence means the two harnesses load different guidance
Closure: RR1/receipt-fixes-cause, RR1/receipt-fixes-owned-paths, RR1/receipt-fixes-proceed-condition, RR1/ticket-adds-ownership-fence, RR1/ticket-adds-acceptance-row, RR1/ticket-adds-red-mutation, RR2/coordinator-derives-ticket, RR2/delegate-does-not-author-criteria, RR2/approval-round-trip-adds-no-independence, RR2/claude-command-resolves-source

## What to build

The reviewer has fixed this light-path change as one ticket; that approval is the
grouping constraint, not a claim that a thinner cut would strand a gate red. The
existing out-of-fence repair route briefly explains why its reviewer-produced
debug receipt becomes an ownership-fenced repair ticket instead of a small spec,
and why the coordinator derives that ticket without an acceptance-criteria
round trip through the implementing delegate. Keep the rationale inside the
existing section and leave its mechanics unchanged.

## Acceptance

- [ ] [RR1] the repair route says the debug receipt already fixes cause, owned paths, and proceed-condition, while the repair ticket adds the ownership fence, an acceptance row, and a red mutation the probes run against, so a small spec would only restate the receipt.
- [ ] [RR2] the repair route says the coordinator derives the ticket and the implementing delegate does not author its own acceptance criteria, because delegate authorship followed by coordinator approval adds no independence to a reviewer-produced receipt; the Claude command symlink resolves to the edited source.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RR1/receipt-fixes-cause | omit the confirmed cause from the receipt's fixed facts | the semantic reviewer | apply the omission, read the complete out-of-fence section, expect the small-spec rejection to omit cause as an already-fixed fact |
| RR1/receipt-fixes-owned-paths | omit the owned paths from the receipt's fixed facts | the semantic reviewer | apply the omission, read the complete out-of-fence section, expect the small-spec rejection to omit owned paths as an already-fixed fact |
| RR1/receipt-fixes-proceed-condition | omit the proceed-condition from the receipt's fixed facts | the semantic reviewer | apply the omission, read the complete out-of-fence section, expect the small-spec rejection to omit the proceed-condition as an already-fixed fact |
| RR1/ticket-adds-ownership-fence | omit the ownership fence from the ticket's additions | the semantic reviewer | apply the omission, read the complete out-of-fence section, expect the route not to name the ticket's ownership boundary as new information |
| RR1/ticket-adds-acceptance-row | omit the acceptance row from the ticket's additions | the semantic reviewer | apply the omission, read the complete out-of-fence section, expect the route not to name the ticket's observable completion criterion as new information |
| RR1/ticket-adds-red-mutation | omit the red mutation from the ticket's additions | the semantic reviewer | apply the omission, read the complete out-of-fence section, expect the route not to name the ticket's falsification input as new information |
| RR2/coordinator-derives-ticket | omit that the coordinator derives the ticket from the receipt | the semantic reviewer | apply the omission, read the complete out-of-fence section, expect ticket derivation to have no named owner |
| RR2/delegate-does-not-author-criteria | assign acceptance-criteria authorship to the implementing delegate | the semantic reviewer | apply the replacement, read the complete out-of-fence section, expect the delegate to author the criteria it must satisfy |
| RR2/approval-round-trip-adds-no-independence | omit why delegate authorship followed by coordinator approval adds no independence | the semantic reviewer | apply the omission, read the complete out-of-fence section, expect the route to reject the round trip without explaining its independence failure |
| RR2/claude-command-resolves-source | alter one rationale term in the canonical `.agents` command without touching `.claude/` | the tracked `.claude/commands` symlink | apply the canonical-file mutation, read `.claude/commands/bench-implement-spec.md`, expect the same altered term, then restore the canonical file |
