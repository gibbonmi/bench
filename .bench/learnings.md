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

## 2026-07-22 — accepting a delegate's green without checking its assertions were non-vacuous

**What happened.** During the FT85 build I accepted the payload slice (later commit
`a6dcec3`) after verifying its gate was green, its diff was scoped to its charge, and
its reported red signals were real. It landed with a defect: the kit-only allowlist row
named `.agents/skills/craft-synthesis`, but the directory is `bench-craft-synthesis`.
That skill therefore kept shipping to consumer repos and into the npm tarball — the
exact least-privilege hole the story existed to close — while the story's contract
assertions passed by asserting the absence of a path that never existed. The next
slice's delegate found it.

**What the right behavior was.** A green gate cannot detect a vacuous assertion: a test
that asserts "path X is absent" passes trivially when X is misspelled. When a delegate's
red signal is an *absence* assertion, verifying the claim requires confirming the named
paths exist in the source tree — one `ls` would have caught this. More generally, the
verification step in `craft-delegate` should treat absence-shaped assertions as a named
case needing an existence check on the identifiers, not just a gate reading.

**Proposed rule change.** Add to `craft-delegate`'s verification discipline: when a
delegate's evidence rests on an assertion that something is absent, excluded, or
withheld, verify the named identifiers resolve to real things before accepting. Consider
whether the gate itself can hold this — a conformance check that every `source` path in
`.bench/consumer-payload.json` exists in the tree would have made the defect red rather
than invisible, and that is a cheap single-source check.

## 2026-07-22 — a load-dependent flake cost three full gate runs before a green commit

- **What happened:** The FT85 review-fix commit was refused twice by a red gate, both
  times on `TestFT78Story5ProofLedger` R14 subtests (`commit gate did not start`,
  `gate owner did not reach pending state`) — different subtests each run, and the
  test passed in isolation in 10s against 80s+ under gate load. The third identical
  attempt went green with no code change. Roughly 35 minutes of wall clock bought
  nothing.
- **Right behavior:** What I did was correct under invariant 1 — I never weakened the
  test, and I re-ran it in isolation to prove the flake rather than asserting it. But
  after the second red on the same known-flaky test (already parked in `IDEAS.md`),
  the better move was to stop and surface the blocked commit with the evidence and a
  recommendation, rather than spend a third full gate run on a coin flip. A known
  flake blocking a landing is a reviewer decision, not something to grind through.
- **Proposed rule change:** Add to the gate/commit discipline: when a commit is
  refused twice by the same test that is already recorded as a known flake and passes
  in isolation, stop and hand back with the evidence instead of re-running. Retrying a
  flake is not iteration toward green.
