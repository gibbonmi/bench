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

## 2026-07-23 — a delegate charge's verification list became the delegate's ceiling

**What happened.** A write-delegate building the gate's capability-skip collector
added `BENCH_REQUIRE_CAPABILITIES: '1'` to two GitHub workflow steps by putting an
`env:` block *above* the `run:` line. That re-indented the run line by two spaces.
Two existing canary fixtures anchor on that line's exact bytes, so both broke with
`mutation anchor ... did not occur exactly once`. Nothing in the delegate's charged
verification list could see it: unit tests, the conformance suite, and the contract
suite were all green. The coordinator's charge had said "do NOT run `bench gate`"
to keep delegate rounds cheap, and had listed the narrower commands to run instead.
The delegate treated that list as the boundary of what needed checking.

**What the right behavior was.** The canary layer is the only thing that grades a
byte change under `.github/workflows/`, and the coordinator knew that — the project
profile says so directly, and warns that the inner byte-shape is load-bearing. The
charge should have included `bench canary <root>` the moment it authorized workflow
edits. The delegate is not blameless either: a charge's verification list is a
floor, not a ceiling, and "which layer grades the artifact class I just changed"
is a question a delegate should ask itself. But the cheaper fix lives in the charge.

**Proposed rule change.** In `craft-delegate`, when a charge authorizes edits to an
artifact class that a specific gate layer owns — workflows and `.bench/` content
(canary), gate output shape (canary), skills and commands (conformance) — the
charge must name that layer's command in the verification list. Pair it with the
converse instruction to the delegate: the list is a floor; if you change an
artifact class the list does not grade, say so rather than assuming it is covered.

**Cost note.** This was caught by the coordinator running the canary during
done-claim verification, one round after the fact. Cheap here; it would have been a
red gate at merge with an unrelated-looking failure message.

## 2026-07-23 — `git commit` after `git merge --squash` swept an entire slice into a capture commit

**What happened.** During the FT87 slice 3 build the coordinator ran
`git merge --squash <branch>` to bring a delegate's work into `main`'s working
tree, intending to land it with a path-scoped `bench commit` so the gate would
grade exactly that diff. Before doing so, `bench commit` refused over an
unexplained working-tree file — a parked idea in `IDEAS.md`, written by the
coordinator moments earlier. The coordinator committed the capture file with
`git add IDEAS.md && git commit -m "capture: ..."`. But `git merge --squash` had
already staged the delegate's ENTIRE slice into the index, and a bare
`git commit` commits the index, not the file just added. The result was one
commit labelled "capture: park the gate phase-timeout headroom idea" containing
that one line plus 649 insertions across eleven files of an unrelated feature
slice — and landed at a moment when no gate had graded it.

**What the right behavior was.** Two independent fixes, either sufficient. Use
the pathspec form, `git commit -m "..." -- IDEAS.md`, which commits only the
named path and ignores the rest of the index. Or — better — capture files get
committed BEFORE a squash-merge stages anything, so the index is never a mix of
two intents. The deeper rule: after any `git merge --squash`, treat the index as
already full, and never issue a bare `git commit` until the intended landing
commit.

**Why the guard did not catch it.** `bench commit`'s attribution block-check is
exactly the mechanism that makes a green gate describe the diff that lands, and
it fired correctly. The coordinator then stepped around it with plain `git`,
which carries no such check. That is the failure mode: the enforcement lives in
`bench commit`, so every plain-`git` commit during a build is unguarded.

**Recovering.** No content was lost or wrong, and a later full gate run on the
identical tree came back green, so the tree is verified — the defect is purely
history: one mislabelled commit conflating a capture with a feature slice.
Repair is a history rewrite, which is the reviewer's call, not the worker's.

**Proposed rule change.** Note in the phase guidance that the doc-only
plain-`git` convention is safe only when the index is otherwise empty, and that
during a squash-merge landing sequence every plain-`git` commit must use the
explicit pathspec form.
