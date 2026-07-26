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

## 2026-07-26 — gate probe env leak spawned a recursive test cascade  [open]

- **What happened:** A throwaway timing probe in `internal/conformance` was
  gated on a custom env var (`BENCH_TIMING_PROBE`). `conformanceSubprocessEnv`
  scrubs only `BENCH_CONFORMANCE_ROOT`, so the flag leaked into the inner
  `go test` that `checkGoCore` spawns, re-ran the probe one level down, and
  seeded a self-perpetuating cascade of full-suite runs whose orphans outlived
  each generation's 600 s timeout. Diagnosis also surfaced that the recursion
  exists without the probe: the inner suite includes the conformance package
  itself, whose unfiltered tests re-invoke the heavyweight checks on a cache
  miss. Cascade was found by asking "what is still running", killed by hand.
- **Right behavior:** Gate any conformance-package probe on an env var the
  subprocess env scrubs (or a `-run` filter that cannot re-enter), and after
  any run that times out mid-subprocess, sweep for surviving descendants
  before trusting later measurements — contention had inflated two numbers.
- **Proposed rule change:** none for the kit prose — the structural fix
  (exclude the conformance package from the inner suite) is already decided on
  the cost-follows-project-size map and rides the FT91 tier-split spec.
