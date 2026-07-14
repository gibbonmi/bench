# Review findings — lossless-shift-recovery (FT79 implementation)

Semantic review of the implementation branch `bench/assign/facb6d11…`
(11 commits, 18 files, base = merge-base with main). Advisory; the branch gate is
green. Fix pass works on the branch; this file is deleted by the green fix commit
that closes the findings.

## Standards

**4 findings (3 hard violations — one minor, 1 judgment call). Worst: S1.**

- **S1 (hard) — outcome-recording sequence duplicated between exit paths.**
  `internal/shift/loop.go:22-29` (`finish`) and `internal/shift/session.go:245-252`
  (`exitPreserving`) both derive Emit → `entry.Outcome`/`entry.Recovery` →
  `intent.Upsert` → exit code. Same package; `exitPreserving` can call
  `os.Exit(finish(...))`. As written, the signal/deadline `os.Exit` path can
  silently drift from every normal exit if the ledger-recording protocol changes.
- **S2 (hard) — fixture harness pasted.** `internal/shift/fault_test.go:68-107` vs
  `:173-206`: `faultFixture` and `faultFixtureNoAgentOverride` duplicate the
  ~30-line repo-setup core (git init, `.bench`, `gate.sh`, `gate-inputs.json`,
  add/commit, chdir, env), differing only in the agent file/env. Extract one core
  the second composes.
- **S3 (hard, minor) — recovery-ref namespace as bare literal in two production
  sites.** `internal/shift/session.go:288` and `internal/shift/result.go:48` both
  spell `refs/bench/recovery/`; `internal/intent/assignment.go` already
  establishes the named-constant idiom for a ref namespace. Collapse to one
  constant.
- **S4 (judgment) — `runPreservingGate` is a pure forwarder** over `runGate`
  (`internal/shift/session.go:209-211`), justified only by a documentary name and
  an anticipated split. Middle Man / Speculative Generality; fold or keep is the
  reviewer's taste.

## Spec

**1 finding (low). Worst and only: the row 11 range assertion.**

- **Spec 1 — row 11's test asserts weaker than its mapped red signal.** The row
  and story 11 pin "exit 2 naming the variable **and accepted range**";
  `internal/contract/runtime/runtime_shift_test.go:674` asserts only the bare
  variable name (`tc.want` values at :653-660), so a mutation dropping the range
  clause from `validate.go`'s messages stays green. Behavior is correct today;
  the contract is under-locked. Fix: extend `tc.want` (or add a second contains)
  to include the range fragment.

All 18 stories and 25 rows otherwise verified implemented with red signals met or
exceeded; no scope creep found.

## Coverage

**2 findings (both medium). Worst: C1.**

- **C1 — simultaneous wall deadline + SIGINT precedence is undecided and
  untested.** `checkpoint()` (`internal/shift/session.go:227-234`) tests
  `deadline` before `interrupted`, so a coincident wall trip + Ctrl-C always
  exits 3, never 130 — deterministic today, but pinned by nothing: no spec line
  decides the precedence and no row sets both flags, so a future case-order flip
  is invisible to the gate. Reviewer decision: pin deadline-wins with a
  test/decision line, or declare it Won't-handle.
- **C2 — finish-time `intent.Upsert` error swallowed.** `loop.go:27` and
  `session.go:250` discard the outcome-recording Upsert error (`_ =`), while
  loop-entry Upserts are checked. A ledger write failure at finish loses exactly
  the after-the-terminal discovery pointer story 17 promises, with no signal.
  Reviewer decision: accept as best-effort, or surface a stderr warning and
  assert it via an injected failure.
