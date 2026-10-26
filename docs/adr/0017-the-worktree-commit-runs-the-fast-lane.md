# 17. The worktree commit runs the fast lane

Status: accepted (2026-08-25)

## Decision

A worktree commit runs the fast lane and not the whole-project gate. The lane
is a declared check list. The composed changes select from that list by path
class, so a commit pays only the checks its own changes need. An unknown path,
a symbolic link, and a gitlink each select every declared check, so the
selection never narrows below the declared lane.

The lane runs in a private checkout of the composed snapshot and takes seconds.
A lane pass authorizes the commit onto the worktree branch. A lane fail refuses
the commit and names the check.

The landing runs the gate on the composed tree with the spec flip applied. A
kit landing whose every changed path sits in a covered family of the reader
table runs the reader set. Every other landing runs the whole-project gate. The
gate alone authorizes a landing, and the reviewed landing accepts a green gate
verdict or a green reader set. A lane pass authorizes nothing at the landing,
which derives its own evidence chain.

A green reader set writes a reader record. That record authorizes its own
landing only, and the green marker stays at the last whole green. After five
reader landings or a changed oracle, the next landing runs the whole-project
gate. After a red whole landing run, each landing runs the whole-project gate
until a whole green moves the marker.

The lane writes one lane record in its own record class. It writes no gate
verdict and no evidence, and a lane record is never reusable green. The kit
root carries a built-in lane, and that lane is the selective one. A linked repo
declares its lane in the phase manifest, and a manifest lane runs as declared
with no selection. A repo with no declared lane keeps the full-gate commit.

## Consequence

One landed change pays one gate run, at the landing. That run is the reader set
for a covered kit change and the whole-project gate otherwise. A Markdown
commit runs no vet check and no build check, and a Go commit runs no prose
check. A lane pass never reads as green, so no consumer can mistake a worktree
commit for a graded tree.
