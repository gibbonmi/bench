# Learnings — usage journal

Append one entry when you deviate from the workflow, make a judgment call you're
unsure about, catch a should-have-asked in hindsight, or catch yourself assembling
the same ad-hoc check a second time (a codification candidate — name the `bench`
subcommand it wants to be). You capture; the reviewer
decides. `/bench-what-next` verdicts every open entry in its reviewed batch
diff: work-shaped and rule-shaped entries become roadmap items (rule-shaped ones
built later under the synthesis discipline), the rest are dismissed with one line
of why. A resolved entry leaves this file, and its verdict is recorded in the
drain commit. The journal holds open entries only; history lives in git. Never
rewrite a kit rule yourself — that is the whole point of capturing here instead.

Format per entry. Heading: `## YYYY-MM-DD — short title  [open]`

- **What happened:** …
- **Right behavior:** …
- **Proposed rule change:** … (or "none")

## 2026-07-27 — a returned delegate's background children still load the machine  [open]

- **What happened:** A write-delegate reported done, but its own `go test
  ./internal/...` sweep and two shell wait-loops were still running. I started
  `bench commit` on its worktree immediately; the gate went red on
  `internal/intent`'s concurrency test timing out — a load flake answering for
  machine contention rather than for the diff. The same commit passed on a quiet
  machine with no code change.
- **Right behavior:** Before running the whole-tree gate on a delegate's work,
  confirm the machine is quiet — the delegate's done-claim says its *report* is
  finished, not that its subprocesses have drained. A `pgrep` for `go test` and
  the delegate's wait-loops costs one call and saves a ~12-minute false red.
- **Proposed rule change:** `craft-delegate` already names the whole-tree gate a
  serialized resource for *concurrent* delegates. Extend that clause to cover the
  sequential case: a returned delegate is not necessarily a drained one, so the
  coordinator checks for live test processes before gating. Possible codification:
  `bench commit` refusing (or warning) when it observes another `go test` against
  the same module.
