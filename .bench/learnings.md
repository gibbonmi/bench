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

## 2026-07-22 — gate flake class recurred and burned five blind runs (FT76 close)

During the FT76 landing, `bench gate` went red 5 of 7 runs on an identical tree, alternating between two unrelated pre-existing checks: the conformance phase's inner core `go test` (4×, output discarded — the diag says only "go test failed") and `TestBinaryRepairContracts/repair_losing-racer` (1×, "did not reach synchronization marker", a hard 2s `time.Now()` deadline). Every suite passes deterministically solo and pairwise; dmesg shows WSL2 "time jumped backwards" clock jumps under the gate's peak parallel load. Same class as the 2026-07-19 entry, now with named culprits. What made diagnosis expensive: the conformance diag swallows the failing test's output, and a PATH shim cannot instrument the inner run because `go test` prepends `$GOROOT/bin` to child PATH. Right behavior: retried at the same tier and diagnosed instead of weakening checks; the deadline widths and the diag are gate-authoring decisions left to the reviewer. Proposed changes: (1) conformance's "go test failed" diag should carry the tail of the inner output; (2) the repair sync-marker deadline (and siblings using wall-clock `time.Now()` deadlines) should scale or use monotonic-friendly generous bounds; (3) consider capping gate phase parallelism on hosts where the load provokes the class.
