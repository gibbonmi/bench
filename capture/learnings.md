# Learnings — usage journal

## 2026-08-21 — spec review round needed two iterations (landing-refusal-diagnostics) [open]

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

## 2026-08-21 — client inherited Envman's marker without its PATH [open]

What happened: a deterministic Codex-client repro found a native Linux Go
1.25.0 binary at `/home/mgibs/.local/opt/go/bin/go`, while the client process
carried `ENVMAN_LOAD=loaded` and omitted that directory from `PATH`. Envman's
generated `PATH.env` declares the directory, but shell initialization trusted
the inherited marker and skipped the file. Clearing only `ENVMAN_LOAD` made the
same login shell discover Go; making the shell interactive without clearing it
did not. A fresh WSL/CLI login therefore works while the affected client shell
does not. A Git-source Bench install can fail at its Go build in this state;
an installed platform binary can still run, but Go-backed project gate phases
can be omitted. The agent initially concluded that WSL lacked native Go, then
replaced the full `PATH` while repairing the run and removed the bundled `rg`
entry, causing a second, agent-owned failure.

Right behavior: diagnose tool discovery from both the executable's known
location and the process environment before concluding a dependency is
missing. Treat a loaded-marker-without-effects shape as partial environment
propagation. When repairing one tool lookup, prepend its directory to the
existing `PATH`; never replace the rest of the client-provided toolchain. Bench
must not search arbitrary home directories or silently trust a user-specific
toolchain path, and a missing required toolchain must not produce a weaker
apparently-green gate.

Proposed rule change: extend the existing SessionStart path wired by Codex to
detect a repository requiring Go when `go` is absent from the harness `PATH`
but discoverable from a clean login shell, then print a loud diagnosis and a
copy-paste recovery command. Make the gate fail visibly rather than silently
omit required Go phases, independent of harness. Add a system regression for
the inherited-marker/reconstructed-PATH shape, and add agent guidance to
preserve the complete ambient `PATH` when prepending a recovered tool
directory. Test client and CLI sessions against the same WSL repository and
toolchain before portability is claimed.
