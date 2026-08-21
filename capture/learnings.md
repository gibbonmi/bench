# Learnings — usage journal

## 2026-08-21 — spec review round needed two iterations (landing-refusal-diagnostics)

What happened: the /bench-write-spec round for FT233 returned nine blocking
findings. Two were red-capability defects the author should have caught: LR2
named a mechanism the tree cannot produce (a short `--base` resolves and lands;
the tip-mismatch string has one producer), and LR19 asserted hostile bytes
(ESC, BEL) that cannot split the record it protects. The author had noticed
the LR2 uncertainty and hedged it into Further notes instead of resolving it.

Right behavior: before locking a coverage row that names a mechanism, trace
the message to its one producer and confirm the claimed input can fire it. A
"the build will trace it" hedge is the signal the row is not lockable yet.
Occurrence ledgers are evidence about an operator's session, not about the
tree.

Proposed rule change: a candidate line for `craft-spec` (reviewer decides at
drain): an occurrence-sourced row that names a mechanism carries a traced
producer citation, or the row locks the observable without the mechanism.
