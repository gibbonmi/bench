# Record carried entries in the ticket form

Blocked by: 03-decide-the-carry.md
Writes: internal/reviewrecord/recordcmd/ticket.go, internal/reviewrecord/recordcmd/ticket_carry_test.go (new), internal/reviewrecord/recordtest/carry.go (new), cmd/bench/help_inventory_test.go
Covers: CF33, CF34, CF35, CF36, CF37, CF38, CF39, CF40, CF41, CF42, CF43, CF44, CF48

## What to build

A repair session calls `bench record ticket` and stages only the checks that its repair delta reaches. The ticket form no longer refuses a missing stage before the write. In the verification step, it asks the carry decision of `03-decide-the-carry.md` about each unstaged owed requirement, against the chunk entry after the freeze. A stage always wins over a carry.

The form writes each carried entry through the one verification writer of the review record, in owed order with the staged entries. The source is the source commit of the candidate. The exit code and the probe result are those of the candidate. The evidence ID, the performer, the model, and the effort come from the call. The native reference is `carried:` followed by the evidence ID of the candidate, with the excerpt and the digest of the candidate.

The output row gains a `carried` cell after the `verifications` cell. The cell lists the carried requirement IDs in owed order, or `none`.

When an unstaged requirement cannot carry, the form refuses through the excerpt form and writes nothing. The refusal keeps the text `these owed requirements have no staged excerpt: <ids>` and lists only the requirements that cannot carry. It then adds `; carry: <id> (<reason>)` for each of them, in owed order. The help line of the form states that it writes the carried entries too.

## Acceptance

- [ ] Under a repair dispatch through `recordtest.Reassign`, after a test-only repair, a call that stages only the reached requirement exits 0 and writes a carried entry for the unstaged system requirement.
- [ ] The carried entry names the repair session as performer, keeps the source commit, the source digest, and the excerpt of the candidate, and has the reference `carried:<candidate ID>`.
- [ ] The `carried` cell lists the carried requirement IDs, and it shows `none` when each owed requirement is staged.
- [ ] A staged requirement that the delta does not reach records from its stage at the repair commit.
- [ ] The staged entry and the carried entry follow the owed order.
- [ ] The refusal names each unstaged requirement that cannot carry with its reason, omits each one that can carry, and leaves the record bytes and every stage unchanged.
- [ ] On a one-chunk fixture, after the carry call and a record commit, `rr.Check` of the chunk refuses only through a review-form repair of that chunk.
- [ ] `TestHelpInventoryIsComplete` and each existing ticket form test pass.
