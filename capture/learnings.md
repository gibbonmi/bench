# Learnings — usage journal

## 2026-08-22 — FT242 spec review needed three iterations [open]

What happened: `/bench-write-spec` first treated a fixture that mimicked Codex
as portability evidence and described clean-login discovery as nonblocking
without a deadline row. After those were corrected, the timeout row still
could not distinguish process-group teardown from killing only the login-shell
parent; Terra/medium required a descendant sentinel before accepting.

Right behavior: a decision source requiring a real harness comparison stays a
manual evidence gate and is never weakened into a fixture claim. Any coverage
row promising process-group timeout plants a descendant whose survival is
observable; returning on time proves only the parent stopped.

Proposed rule change: add both checks to `craft-spec` review — source-required
real-environment evidence cannot be replaced by simulation, and a process-group
timeout row names a descendant-survival oracle rather than only elapsed time.
