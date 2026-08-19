# Retro — ft226-test-home-isolation

Landed `073e1f5a` on `main`, 2026-08-19. Line: `opus` / medium, orchestrator and
three write/review delegates. Three tickets, one review round, one repair pass.

## What the build was

The kit's own `internal/worktree` tests wrote pool worktrees into the operator's
real `BENCH_HOME` and never removed them — 1,719 keys and 91 MB at build time.
One fixture was the leak; the build fixed it and then made the whole package
structurally unable to reach the operator's pool under any driver.

## What worked

**The two-level fix was worth the extra ticket.** Ticket 01 alone would have
closed the leak. Ticket 02's `TestMain` is what makes a *future* fixture fail
loudly instead of accumulating silently, and the mutation probe proves it bites.
The spec was right to price the detector as oracle logic deserving the cheap
row's bump.

**Independent verification caught nothing wrong, and that is still the result
worth having.** The coordinator re-ran probe (a) itself rather than accepting the
delegate's transcript. It reproduced exactly. The discipline cost one test run.

**The write delegate reported a spec defect instead of satisfying it.** Told to
run `TMPDIR=/nonexistent go test …`, it found the command never reaches
`TestMain`, said so plainly, and supplied the `GOTMPDIR` form that does. A
delegate that had quietly "passed" the probe would have left DT4 unproven behind
a green checkbox.

**The Coverage axis earned its slot.** Its worst finding (C1: a walk error
suppressing the entire residue listing) was invisible to the gate, to the Spec
axis, and to the coordinator, and it defeated the report exactly when a leak was
hardest to diagnose.

## What went wrong

**A blanket checkbox tick asserted something false.** The coordinator ticked all
three tickets' acceptance boxes with one search-and-replace. That marked ticket
02's probe (b) box complete while the same commit's verification log recorded
that command as not working. The Spec axis caught it. Ticking a box is a claim
per box; a bulk edit cannot make a per-box claim.

**Two axes disagreed and one was wrong.** Standards cited a real duplicate of the
new `withinDir` in production `insidePool`; the Spec axis asserted no duplicate
existed. The coordinator read the source and resolved it against the delegate.
Axis output is evidence to check, not a verdict to merge.

**The landing needed two unplanned steps.** A concurrent session moved `main`
past the frozen review base, so the source needed a reauthorize; and
`bench worktree land` refuses unless source and destination carry identical
staged spec bytes, so the verification log needed its own destination commit
first. Neither is in the spec's process notes. The second is captured as a
learning.

**The destructive step could not be run by the agent at all.** The auto-mode
classifier denied every `rm -rf` of the operator's pool, so the reviewer ran the
sweep. Plan-before-apply made that handoff cheap — the plan was already printed,
verified against planted hostile shapes, and reviewable.

## Evidence quality

Eight of twelve acceptance rows are backed by code the gate runs. Four (EV1, SW1,
SW2, SW3) are backed only by the spec's prose verification log, because their
subject is the operator's machine. SW3 was ultimately taken twice: once in the
integration worktree, once as a fresh full gate on landed `main`, both leaving
the pool listing byte-identical. That second run is the stronger evidence and
was only available after landing.
