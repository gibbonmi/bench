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
