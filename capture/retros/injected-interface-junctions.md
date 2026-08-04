# Retro — injected-interface-junctions

## Outcome

Promoted terminal. Subject `adb1a5deb87a8e8cb45e967d407ed2ff98b6fef9`, published
working-branch commit `6e085dd`, promotion tree `bd73880`, evidence
`v1:dd5239c3…`. The spec reads `Status: implemented` and its acceptance coverage
map enumerates 14 rows, all dispositioned implemented — the story-2
unreadable-metadata row lands as a documented pre-composition refusal rather than
a junction test, which is the disposition the round-1 repair established and
rounds 2 and 3 confirmed.

Ten assignments checkpointed, integrated, and released: six feature tickets
(ReleaseOwner, abandon decayed-family, canary runner, gitguard checker,
`ShapeUnknown` fixtures, injected-port conformance check), one binding
follow-up, and three repair tickets across two review rounds. Three semantic
review rounds ran under codex `gpt-5.6-sol` at high effort, read-only: round 1
FAIL with five findings, round 2 FAIL with four, round 3 PASS with zero. Round 3
was re-run once against a recomposed candidate after the branch tip moved; the
re-run independently confirmed the recomposition delta touched no reviewed
surface.

## Gate-stage timings

Per-stage timings are not recoverable for the promotion gate: `bench spec build
promote` emits only its TOON status line, so the phase evidence its gate run
produced never reached a retained surface. What was measured is a wall-clock
bound — the promotion exceeded the 120-second foreground limit and completed in
the background. By contrast the `bench commit` that landed the carry commit
printed per-phase evidence inline (`[conformance] 17 example-agreement 1ms`,
`phase conformance: green`, `gate: green`), which is the shape the promotion gate
should also emit. The scoped probe runs were sub-second: the conformance package
under `-count=1 -v` with `BENCH_CONFORMANCE_ROOT` set finished at 0.008–0.009s
across two independent executions with twelve subtests visibly running.

## Ticket-versus-spec-slice and delegate performance

This session ran exactly one delegate charge, and it was ticket-sized rather than
spec-sliced: `repair-spec-and-comment-consistency` at `sonnet` / low, in an
isolated assignment worktree, fenced to `spec.md` and
`internal/conformance/injected_ports_test.go`. It returned in 109 seconds having
made all four required edits plus one it correctly inferred — extending the edge
inventory's registry-drift bullet from four cases to five so it stayed consistent
with the new orphan coverage row. Nothing landed outside the fence.

The ticket's red-mutation table is what made that charge cheap. Because the
mutations were written as executable greps the delegate could run against its own
output, the charge carried its own self-check and needed no coordinator
correction. That is the contrast worth keeping: a spec-sliced charge on this same
material in earlier rounds produced findings, and the ticket-sized charge with
encoded mutations produced none.

## Coordinator catches

The delegate reported that `rg "through the real planner"` still returned two
hits, against a pickup instruction that said the grep must return nothing.
Verifying rather than acting on the instruction showed both surviving hits are
the accurate decayed/husk claim in story 2 and its line justification, not the
unreadable-path claim the finding named. Rewording true prose to satisfy the
grep would have been the exact failure mode — fixing the citation instead of the
invariant — that produced round 2 in the first place, so the hits stayed and the
reviewer was told the grep would not go to zero.

The done-claim was verified with different-kind evidence rather than a rerun of
what the delegate already ran: the coverage-map row count moved 13 → 14 under
`bench coverage --check` comparing the pre-edit candidate spec against the
post-edit worktree, which is parser-observed proof the new row exists and parses,
independent of any test outcome. The scoped conformance run was then re-executed
fresh with `-count=1 -v` to prove execution rather than a cached green.

## Agent-experience improvements

### Bench CLI

`bench spec build promote` should emit its gate run's phase evidence and timings
the way `bench commit` does. Right now a terminal retro cannot record stage
timings at all, which is a required heading with no available source.

`bench spec build assign --ticket` resolves its argument against the spec's
`tickets/` directory and refuses both repo-relative and absolute paths with the
same message, `spec build ticket must name one regular ticket file`. That message
reads as "the file is missing" when the file is present and the *form* is wrong;
it cost two failed invocations before the parse in `assign.go` disambiguated it.
A refusal naming the expected form would have cost none.

`promote`'s clean-checkout precondition composes badly with its
recompose-discards-review rule. Any unrelated dirty file blocks promotion, and
clearing it by committing moves the branch tip, which discards the bound review
and forces an entire extra review round. The two rules are individually correct
and jointly expensive. Either promote tolerates out-of-set dirt the way path-
scoped `bench commit` reasons about ownership, or its refusal names stashing as
the cheap route so the reader does not reach for the expensive one.

### Skills

A ticket's red-mutation greps must be expressed against the defect, not against a
phrase. This ticket's `rg "through the real planner"` mutation was written as an
absolute — the grep must return nothing — for a phrase that legitimately survives
in accurate prose. Written that way, a grep that stays non-empty looks like an
incomplete repair, and the tempting resolution is to edit true sentences until
the tool goes quiet. `craft-tickets` should say that a mutation grep names the
false claim, not the words it happened to be made of.

### Process

Stash, do not commit, to clear `promote`'s clean-checkout precondition. Committing
an unrelated reviewer-owned draft to unblock promotion moved the tip, discarded a
passing review, and cost a full extra sol round plus a second reviewer-run
promote. The draft had to be preserved either way; stashing preserves it without
touching the tip.

Encoding the regression loop in the repair ticket worked. Round 3 is the first
clean round of the three, and the difference between it and round 2 is that its
ticket carried explicit red mutations and an instruction to scope by invariant
rather than by citation. That pairing — a mutation the delegate can execute, plus
a whole-artifact reread for internal agreement — is what a repair ticket should
carry by default.
