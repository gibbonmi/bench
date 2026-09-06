# 20. A decision lives in its ticket file

Status: accepted (2026-09-06)

## Decision

A decision of a decision map lives in its ticket file. Each decision is one
file under the map's tickets folder, and the file is named by the decision's
number. That file holds the question, the answer, the type, and the blockers of
that decision. The ticket file is the one canonical home of that answer. Thus a
session that resumes one decision reads one file.

The map file is an index. It holds the title, the status, the destination, the
notes, the decisions made so far, and the four terminal sections. Each decision
made so far is one gist line. A gist line links its ticket file and condenses
that answer into one sentence. Thus the index gives the low-resolution view of
the map, and it holds no canonical answer.

The topic folder moves as one unit. The tickets and the map-owned assets sit in
one folder beside the index. Spec authoring moves that folder into the spec,
and the index moves with it. Retirement removes the folder and the index
together. Thus a reference between the index, a ticket, and an asset stays
correct after the move.

## Consequence

A shaping session loads the index once, then reads only the tickets that it
works. The index can disagree with the tickets, so the gate refuses three drift
faults: a resolved decision with no gist, a gist that links an unresolved
decision, and a gist that links no ticket file. The gate also refuses a map
with no ticket file and a tickets folder with no index.
