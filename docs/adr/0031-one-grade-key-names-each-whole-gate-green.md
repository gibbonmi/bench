# One grade key names each whole-gate green

Status: accepted

## Decision

One grade key names each whole-gate green by its verdict inputs only.
The members of the key are these inputs:

- the graded tree
- the toolchain: the resolved oracle, the executable search path, and each declared environment, path, and tool input
- the baseline phase schedule, where an absent schedule and an empty schedule are different
- the capability posture
- the run binary member, which follows the source root of the run binary
- the digest of the running executable

The checkout path, the run purpose, the execution policy, and the baseline runner identity are not members.
The oracle digest keeps them, and that digest names the gate cache of each checkout.
A running executable that does not read opens the subject, so that run reuses no green.

A green under a grade key answers a run of the same key from any checkout of the repository and from any whole-gate verb.
The whole-gate verbs are the gate, the commit, the merge, the landing, and the checkpoints.
Each verb looks up a green under the key before it selects a reduced schedule.
On a hit, each verb still runs its own purpose checks and evidence checks, so a reused green skips no refusal of that verb.
A reused green records the gate cache of its checkout, with the recorded time of the source run.
So a reuse never extends the freshness window of the green.

A background grade is a producer of a green under the grade key.
It grades a chunk subject or a landing subject, so its green answers the chunk checkpoint or the landing of the same key.
A chunk grade runs as the executable that the wrapper of its worktree resolves, and the chunk checkpoint runs as that executable too.
A landing grade runs as the authenticated broker, so its green answers the landing.

The capability posture binds every reuse, so a relaxed green never answers a strict run.
The running executable binds every reuse, so a green from a candidate build never answers the landing of the broker.

In the kit, the complete checkpoint runs under the installed landing broker, so its green answers the landing of the same tree.
When the running executable differs from the broker, the checkpoint replaces its process with the broker.
The broker then runs with the same arguments in the same working directory.
When the checkpoint cannot authenticate the broker, it refuses and runs no gate.

A repository that declares no build inputs runs the complete checkpoint in its own process.
The chunk checkpoint always runs in its own process, so it keeps the grammar of the candidate.

## Consequence

A spec landing that follows a green complete checkpoint of the same tree pays no second whole gate.
A change of one member makes a new key, so the gate runs again.
A red run or a timed-out run removes the green under its key, so no checkout reuses a contradicted green.
