# Learnings — usage journal

- 2026-08-09  Found and killed orphaned child processes (`mv`, `go-build.sh`,
  `git worktree list`) left running for 3+ days. The first read of this blamed
  two named Go tests and missing `t.Cleanup` teardown in test helpers; both
  parts were wrong. Neither test name exists in the tree, and the owner is
  production code, not a test helper: `internal/testreport/testreport.go` armed
  `signal.NotifyContext` for `os.Interrupt` alone while
  `runbinary.canonicalBuild` detaches the builder into its own process group via
  `Setpgid`, which also removes the kernel's fallback — a detached group never
  receives the terminal's SIGINT or SIGHUP. So any owner killed by the SIGTERM
  that session and harness teardown sends left `go-build.sh` and its `go build`
  children alive with no cleanup path (repro: SIGINT green, SIGTERM and SIGKILL
  red). Fixed by spec `light-path-cancel-signal-parity`, ticket
  `arm-every-owner-for-term-and-hup`: `subprocess.CancelSignals` is now the one
  source of the trapped set and all thirteen production registrations take it.
  SIGKILL stays open — it needs `Pdeathsig`, which is a reviewer decision.

  The judgment worth keeping: a symptom read off `ps` output named a plausible
  owner, and the plausible owner was not the real one. The repro loop is what
  separated them, and it should come before the attribution, not after it.
  No rule change proposed — `/bench-debug` already says exactly this.
