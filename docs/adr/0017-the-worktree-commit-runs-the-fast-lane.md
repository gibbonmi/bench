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

The landing runs the one whole-project gate on the composed tree with the spec
flip applied. The gate alone authorizes a landing, and the reviewed landing
accepts only a green gate verdict. A lane pass authorizes nothing at the
landing, which derives its own evidence chain.

The lane writes one lane record in its own record class. It writes no gate
verdict and no evidence, and a lane record is never reusable green. The kit
root carries a built-in lane, and that lane is the selective one. A linked repo
declares its lane in the phase manifest, and a manifest lane runs as declared
with no selection. A repo with no declared lane keeps the full-gate commit.

## Consequence

One landed change pays the whole-project gate once, at the landing. A Markdown
commit runs no vet check and no build check, and a Go commit runs no prose
check. A lane pass never reads as green, so no consumer can mistake a worktree
commit for a graded tree.
