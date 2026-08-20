# Review pickup — worktree-exec-run-binary

Base `8c8aca66`, tip `290febc7`. Three axes, cheap tier, read-only.

Raw findings: 5. De-duplicated repair targets: 4 (one finding is `no-op`).

## Standards

1 finding. Worst: the `execEnv` doc comment's proportion.

- **`execEnv`'s doc comment breaks the file's comment density and its opening
  sentence restates the code.** ~250 words over four paragraphs for a six-line
  function, in a file where every other function carries one sentence or none
  (`internal/worktree/exec.go:79-96`). `craft-comments` requires matching the
  surrounding file and deleting a comment that restates the line below it. The
  rejected-alternative and "nothing executes the value" paragraphs are genuine
  why-content and earn their place; only the proportion and the opener are at
  issue. No documented cap applies — the profile's prose budget covers skill and
  command markdown, not Go doc comments. Disposition: **ask-user**.

## Spec

1 finding. Worst: WX20 claimed but unrecorded.

- **WX20's composed-run evidence exists nowhere in the tree.** The spec promises
  "The build demonstrates it once and records the evidence"
  (`specs/worktree-exec-run-binary/spec.md:214-215`); the ticket checkbox is
  unchecked and the only trace is a commit-message assertion. WX20 is the sole
  seam closing the gate-owner half of WX3 and WX12, so three rows rest on an
  unrecorded claim. Disposition: **ask-user**.

Row tally: 17 of 20 delivered and observed; WX20 claimed but unobserved; WX3 and
WX12 partial pending it. WX18's reliance on the pre-existing
`TestGateEntryRejectsLegacyBeforeRunningOldTableAndRunsReplacementOnce` was
verified sound — that test supplies a real absolute executable and drives the
phase table to green.

## Coverage

3 findings. Worst: the live-symlink disposition has no test.

- **A live symlink at the wrapper path is a decided behavior with zero test
  backing.** The spec's edge inventory decides it is emitted as given, but no WX
  row maps to it and `TestExecChildTakesNoMarkerFromANonWrapper`
  (`internal/worktree/exec_test.go:104-120`) covers only absent, directory,
  FIFO, and dangling symlink. Swapping `isRegularFile`'s `os.Stat` for
  `os.Lstat` — a plausible symlink-hardening change, given WX10 reads as "reject
  symlinks" — would silently invert the spec's own decision with nothing red.
  Disposition: **auto-fix**.

- **`assignment()`'s line-exact reader silently truncates a value containing an
  embedded newline** (`internal/worktree/exec_test.go:283-290`), returning a
  wrong-but-plausible value rather than missing cleanly. No fixture builds a
  worktree path from anything but a plain `t.TempDir()`, so the permitted-byte
  half of the profile's checklist is untested for the path-composition code this
  diff added. Production risk is low: both consumers
  (`internal/adopt/doctor.go:397`, `internal/gate/run_transaction.go:36`) only
  test non-emptiness, and neither line-splits. Disposition: **ask-user**.

- **The two-line gate-entry refusal is not asserted as two lines.** A mutation
  merging both `echo` calls onto one line, keeping both substrings, passed all
  three gate-entry tests. Cosmetic structure, not a behavioral contract.
  Disposition: **no-op**, recorded so it is not rediscovered.
