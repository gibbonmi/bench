# Install the two verification loops and their record

Blocked by: move-slicing-into-write-spec.md
Writes: .agents/commands/bench-write-spec.md, projects/benchkit.md, internal/anchors/registry_data.go, internal/conformance/fixture_bite_test.go, tests/canary/workflow-guidance-anchors/

## What to build

Step 9 of bench-write-spec.md becomes two verification loops run by a
read-only same-family mid-tier delegate at high effort through the harness's
native agent surface: loop 1 takes the spec alone with the existing
falsification questions, the slicing step runs only after loop 1 accepts, and
loop 2 takes the breakdown alone against `craft-tickets` rules. Each loop
repeats author-fix/re-review uncapped until no blocking findings — a recorded
reviewer exception to the iteration-cap declaration, each round still reported
in one line — and reviewer sign-off stays the hard stop. This ticket also
re-points the slicing step's trigger (landed by move-slicing wording it
against the old falsification review) to loop 1, and re-points benchkit.md's
ticket-breakdown Lines row from `/bench-implement-spec` to the write-spec
slicing step. The retired
one-iteration wording lands as a Forbid+Require pair, rewriting the
`write-spec-review-made-conditional` and `write-spec-review-tier-escalated`
fixtures. On close the author writes `Verification log: spec <n> + tickets <m>
iteration(s) to accept — <note>` into the spec (the Template gains the line
under `Decision source:`, pinned RequireInSection by a new fixture), and a
count above 1 in either loop appends one `capture/learnings.md` entry (which
stage missed, what was caught, why, the proposed rule change — new fixture).
The rewrite keeps `Bootstrap authority before execution` at exactly its
current two occurrences in the file (a bespoke check counts them), and the
step-9 reflow updates the byte-exact step-9 substring hardcoded in
`TestWorkflowCadenceAnchorsRejectDeletionAndSwap`
(`internal/conformance/fixture_bite_test.go`) in the same change.
benchkit.md's falsification-pass and ticket-breakdown-review Lines rows are
rewritten to the two-loop state — both section headings (`Spec falsification
pass`, `Ticket-breakdown review pass`) stay, each carrying a live Require
needle no fixture backs; the `benchkit-spec-ownership` fixture is
owned by main-session-authorship.md (its mutation string sits in the
spec-authoring bullet), so this ticket leaves that bullet and its pinned
string untouched. Blocked by move-slicing-into-write-spec.md because the
loop clause states "the slicing step runs only after loop 1 accepts" and
"loop 2 takes the breakdown" — both presuppose the write-spec slicing step
that ticket installs. Shares bench-write-spec.md and benchkit.md with
main-session-authorship.md (serial by Writes overlap), and registry and
fixture paths land serially across the whole spec.

## Acceptance

- [ ] the two-loop clause names read-only, same-family mid, high effort,
      native agent surface, spec-loop-before-slicing, breakdown-loop-after,
      uncapped-until-clean, and sign-off as the hard stop, with the retired
      one-iteration text tripping a Forbid (covers WF6)
- [ ] the Template requires the `Verification log: spec <n> + tickets <m>`
      line and the process writes it at close (covers WF7)
- [ ] a count above 1 in either loop requires the `capture/learnings.md` entry
      (covers WF15)
- [ ] the slicing step's trigger reads loop 1, no falsification-review
      reference remaining (covers WF6)
- [ ] benchkit.md's falsification and breakdown rows match the two-loop
      process, the breakdown row pointing at the write-spec slicing step,
      not `/bench-implement-spec` (covers WF10)
