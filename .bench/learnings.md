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

## 2026-07-19 — transient compiled-core flake burned a full gated commit

During the FT83 slice 3 close, the conformance phase's compiled-core inner `go test` failed once during a gated `bench commit` and passed on immediate retry with an identical tree; D1 independently hit a one-off `TestDoctorShimContracts` failure that also passed standalone. What happened: a transient/isolation-sensitive unit test can red a whole ~6-minute gated commit. Right behavior: retry-once-at-same-tier was correct and worked both times. Proposed change: identify and deflake the isolation-sensitive test(s) (or run the doctor-shim contract with tighter env isolation) so the oracle stops charging a full gate rerun for noise.

## 2026-07-19 — parallel worktree delegates collide with one-assignment-per-session  [open]

- **What happened:** Launching three parallel write delegates with harness worktree isolation made the WorktreeCreate hook grant one assignment and refuse the other two ("conflicts with its existing assignment"); a retry also failed. Workaround: `bench worktree create --request <distinct-id>` per extra delegate, pointing each agent at its path, releasing with the same request id after porting diffs.
- **Right behavior:** The workaround respected the lifecycle (ownership, recovery refs) and worked; parallel fan-out per the delegate discipline should not require it.
- **Proposed rule change:** Either let the worktree-hook adapter key assignments per delegate rather than per session, or document the manual `--request` route as the canonical parallel-delegate path in the delegate skill.
