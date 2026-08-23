# Absence is the only authoritative empty state a control record may report

Every control record the kit reads — the learnings journal, the roadmap, the decision maps, the specs — can fail to be read in several distinguishable ways. The kit's standing rule is that only one of them may render as "nothing to report". A file that is genuinely not there is an authoritative empty state and prints the empty result on exit zero.

A file is a *failure*, not an empty result, in any of these cases:

- present but unreadable;
- present but zero-byte;
- present but invalid bytes;
- present but not the shape its parser recognizes;
- not a regular file at all.

No surface may collapse it onto the empty result. The defect this closes is that a silent empty state is indistinguishable from an authoritative one. An agent reading a clean board cannot tell "no open learnings" from "I could not open the journal". Every downstream consumer — the dashboard, the drain, the review base — inherits that lie.

One classified reader owns the distinction. It does four things:

- it stats the path before reading, so a dangling symlink classifies as unreadable rather than reporting the not-found error that a plain read would surface;
- it rejects special files before opening them. So a FIFO parked where a control record belongs can never block a command on a read that never reaches end-of-file;
- it follows a symlink to a regular file, so a linked journal is read rather than refused;
- it preserves the underlying reason on a permission failure.

The shared state vocabulary is absent, empty, malformed, unreadable, wrong-type, and unsupported-schema. The last of those is a parser's verdict rather than the reader's. The bytes arrive intact, but the document is not the one the parser was written for. So it is produced by each parser and asserted at the surfaces, not in the reader.

The posture on those states is split by what the surface is for. Query commands — the ones an agent runs to get an answer — fail closed: exit one with a structured error line naming the state. Where the failure is per-entry rather than whole-document, they render every well-formed row *plus* an explicit row naming the broken one. They still exit one, so a malformed entry is visible rather than dropped.

The ambient dashboard keeps exit zero, because its job is to render whatever it can. But it renders an explicit unknown row for any signal whose read failed, instead of a fabricated zero. That row is never suppressed by the zero-count gate every other row passes through. The roadmap context snapshot likewise survives one bad source, degrading that source's block to unknown rather than losing the whole picture.

Read safety and the fail-closed posture are separable, and only read safety extends everywhere. The housekeeping counters — spec retirement, the roadmap reconcile scan — classify before they open. A special file at a spec path would otherwise hang the command the session-start hook runs. A board that blocks forever is strictly worse than one missing an advisory row. They keep their bare-count shape and degrade quietly, because a wrong housekeeping count costs little.

The same rule governs the default-branch fact. There is one owner that resolves it and reports whether it resolved. No code path may substitute a guessed branch name for a repository whose default cannot be determined. Callers take an explicit posture on the unresolved case:

- a command computing a review base fails closed and names the per-branch git-config escape hatch, so the reader is told how to proceed;
- advisory signals tolerate it and skip;
- the pre-push hook is the one deliberate exception. It bakes a fail-safe branch token, because a guard protecting nothing is worse than one on a guessed branch. The installed hook re-resolves the real default before it ever reaches the baked value.

## Considered options

- **Fail closed everywhere, including the dashboard.** Rejected. The ambient board is run by a session-start hook on every session. Turning an unreadable optional file into a hard failure would make a routine repository state block the one surface that orients a cold session.
- **Degrade everywhere, including query commands.** Rejected for the opposite reason. A query command's whole contract is that its answer is trustworthy. A degraded row inside an otherwise successful exit-zero answer is exactly the ambiguity this decision removes.
- **Hard-fail on a malformed entry without rendering the good ones.** Rejected: the reader loses every well-formed entry to one broken heading, which trades a silent drop for a total one.
- **Schema or version markers in control records.** Rejected: the state means "this is not the document the parser expects", a property of the document's shape. Marker-based detection would require every hand-edited control record to carry ceremony it has never needed.
- **Fail closed on a zero-byte file's absent-looking emptiness.** Not an option, but a consequence worth stating. A zero-byte file is something someone created and left unwritten, so it takes the same failure posture as any other non-absent state. Reporting it as never-created is the specific collapse this decision forbids.

## Consequences

- The symbol-indexing command is deliberately outside the migration. Its no-follow rule for tracked symlinks is a recorded decision: following a link would index the target's symbols under the link's path. It would emit source anchors that do not hold, so the follow-a-symlink rule is right for control records and wrong there. Its failures are also per-file, within a listing derived from the index. There, a tracked-but-deleted file is a routine dirty-tree state that failing closed would turn into an error. Its existing skip rows already surface degradation honestly, which is the property this decision buys everywhere else.
- A signal's count and its listing derive from one scan, so the dashboard and the listing can no longer disagree about what was readable. A count surface fails closed on the same tree its listing fails on, rather than printing a number nothing supports.
- The unresolved-map tally counts only files that claim to be maps. An index file is exempt by name, before any read. A content-based rule cannot separate a directory's README from a map that is simply broken. Every other ordinary name in that directory remains a candidate and still fails closed.
- The hostile fixtures this posture requires — permission-denied, FIFOs, symlinks — depend on filesystem features a host may lack and a root user may ignore. So they declare capability requirements and skip honestly. Because a skip and a deleted assertion both look green, the capability-strict run is what proves the states are implemented rather than merely unasserted.
