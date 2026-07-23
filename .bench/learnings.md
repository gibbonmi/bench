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

## 2026-07-23 — a spec's load-bearing tooling claim went unverified into three stories

**What happened.** `specs/cli-grammar-and-capability-evidence.md` asserts, in its
Solution section and again in its Implementation decisions, that a structured line
"is written to stdout before the skip, so it survives non-verbose `go test`". The
gate collector (story 10) and strict mode (story 11) are both built on that claim.
It is false: `go test` without `-v` discards a passing or skipping package's stdout
*and* stderr. The line is only visible under `-v`, under `-json`, or when the test
binary is executed directly. The gate's own phases invoke plain
`go test -count=1 <pkg>`, so nothing teeing that stream could ever observe a skip.
The build discovered it only because the coordinator ran an independent probe
during done-claim verification — the delegate's package was correct and its two
coverage rows bit cleanly under mutation, so gate-green and row-green would both
have passed while the feature did nothing.

**What the right behavior was.** A spec claim of the form "tool X behaves like Y"
that three stories depend on should be executed once before the spec closes, not
reasoned about. The `/bench-write-spec` step-9 falsification pass is the designated
place for that, and it was deliberately skipped here on the grounds that its only
firing trigger was the Handoff's uncertainty flags, all of which belonged to earlier
slices. That reasoning was sound on its own terms and still missed this, because the
claim was not flagged as uncertain by anyone — it read as an obvious fact.

**Proposed rule change.** Add a trigger to the falsification pass, or to the spec
author's checklist: any spec sentence asserting observable third-party tool behavior
that a story's seam depends on must carry either a cited command whose output was
actually run, or an explicit uncertainty flag. "Survives non-verbose `go test`" is
a one-command check; the cost of verifying it is far below the cost of building
three stories on it.

**Also worth noting.** The delegate-verification discipline is what caught this. The
`craft-delegate` requirement to "probe at least one accepted behavior independently
of the delegate's own tests" was the only step in the chain that could have — the
unit tests, the mutation testing, and a full gate run would all have been green.
