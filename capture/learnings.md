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

## 2026-08-03 — Skipped two charged craft-tickets steps at FT181 ticket derivation  [open]

- **What happened:** The coordinator derived four tickets without executing "Discover the contracts before writing files" (no junction row for the specbuild→worktree abandon contract) and without the prefactor rule (the shape classifier was a shared primitive of two tickets, stated in prose to both delegates, each of which derived its own copy). Ticket 3's Assumptions line actively encoded the gap: "specbuild abandon tests drive a counting fake owner, so no internal/worktree change is needed here" — a spec ambiguity converted into a build instruction. The semantic review then paid two repair rounds for what the skipped steps would have caught at derivation (findings P1/C1, S1).
- **Right behavior:** Execute both steps at derivation and land their output in the ticket artifact, so a skip is visible.
- **Proposed rule change:** Partially landed 2026-08-03 in 2892501 (reviewer-directed): the ticket template now requires a `Contracts:` field; craft-spec gained the composition degenerate and the existing-control edge rule. Residual to verdict: the skip was possible because prose-only steps leave no artifact — any remaining prose-only steps in the phase commands deserve the same landing-site treatment (prefactor has no artifact slot yet).

## 2026-08-03 — Auto-repaired every review finding without triaging blocking vs contestable  [open]

- **What happened:** The round-1 review returned 8 blocking findings; the coordinator routed all into repair delegates without distinguishing must-fix production defects (P1, P2) from judgment-priced hardening (S1, C2, C4) and without surfacing the triage to the reviewer before spending the repair rounds.
- **Right behavior:** When a review returns material blocking findings, present the triage — "N must-fix, M contestable, promote-now or harden-first" — before charging repairs; the finding's disposition may be the review delegate's, but build-priority is the reviewer's.
- **Proposed rule change:** /bench-implement-spec's finding-disposition paragraph could require the triage be shown when a review returns more than a threshold of blocking findings; reviewer decides the threshold or rejects the rule.

## 2026-08-03 — Did not audit my own process compliance when the extra work appeared  [open]

- **What happened:** Two direct signals (the ticket-3 delegate's fake-owner caveat at integration; round-1's "pinned against fakes" findings) pointed at coordinator process skips, and the coordinator absorbed the repair work without asking where it came from until the reviewer forced the question.
- **Right behavior:** When a review produces findings the build's own discipline claims to prevent, the first question is "which charged step did not run," answered against the skill text before the repair pass starts — invariant 1 applied to process, not just code.
- **Proposed rule change:** none (behavioral; the Contracts-field landing makes the specific skip visible, and the retro carries the pattern).

## 2026-08-03 — Review-receipt disposition "accepted" means an open defect and blocks promote  [open]

- **What happened:** The coordinator wrote a review receipt marking closed findings "accepted" (intending endorsed/closed); `hasAcceptedFinding` reads that as an unrepaired defect, promote refused, and clearing the recorded review required a recomposition round.
- **Right behavior:** Closed findings carry a non-"accepted" disposition (resolved/endorsed); "accepted" is reserved for a defect the review accepts as real and unrepaired.
- **Proposed rule change:** The receipt-skeleton helper (parked in capture/IDEAS.md 2026-08-03) should enumerate the disposition vocabulary, or the receipt schema should validate the field so the wrong word is refused at submission.

- 2026-08-03 — Coordinator probe read a capability-skip as green: a mutation probe of the injected-port conformance check ran `go test | tail -1` without `BENCH_CONFORMANCE_ROOT`, the test skipped, and two runs read "ok" as "probe survived" (one also hit the test cache). Caught only because the probe expected red and green was the anomaly. Right behavior: probe evidence must prove the test *executed* — run probes through `bench test` (renders skip evidence) or `go test -v -count=1`, and treat `(cached)` or near-zero runtimes as invalid probe evidence. Proposed rule: add a probe-evidence clause to craft-delegate/craft-tdd naming execution proof (not exit status) as the probe's acceptance.

- 2026-08-03 — T6's fence missed three registration surfaces (executable binding in checks_test.go, canary family in registry_test.go, profile check-input row), forcing a follow-up ticket. The spec review had flagged the registry.Check row (finding 12) but the ticket derivation never asked what else a check registration touches; the previous check's landing commit (`git show --stat`) was a complete wiring manifest, and the conformance suite's fail-closed diagnostics enumerate the same surfaces. Right behavior: a ticket that adds an entry to an extensible registry derives its fence from the previous entry's landing commit, and contract discovery traces every crossing of the new entry's *name*, not just its data rows. Proposed rule: add this fence-derivation clause to craft-tickets' contract-discovery step.

- 2026-08-03 — FT181's "delegate caveat relayed but not acted on" repeated: the abandon re-drive delegate reported "unreadable metadata refuses before the planner is reached" and the coordinator recorded it as no-divergence instead of checking it against the coverage row's promise (real-planner composition), leaving a vacuous junction test for sol to catch (review round 1 FAIL, findings SP1/C1/S1). Right behavior: a delegate observation that contradicts a coverage row's mechanism is a row-disposition event — reclassify or surface before checkpoint, never checkpoint over it. Proposed rule: craft-delegate's verification step adds "match every delegate caveat against the charged rows' mechanisms, not just their outcomes."

- 2026-08-03 — Sol round 2's four findings were all blast radius of round 1's repairs: repair tickets were scoped to the findings' cited lines, so amending a coverage row left the story prose contradicting it and adding a fifth check message left the spec contracting four. Right behavior: a repair ticket's scope derives from the invariant the finding names, and before resubmitting for review the coordinator re-reads each repaired artifact whole for internal agreement (story prose vs rows vs implementation decisions), plus greps for the stale contract the repair obsoletes. Proposed rule: craft-tickets' repair-ticket guidance adds "scope by invariant, not citation; sweep the edited artifact for agreement before the next review round."

- 2026-08-03 — The round-2 findings were not a scoping *judgment* failure but a scoping *capability* failure: `repair-orphan-registry-rows.md`'s ownership fence was `internal/conformance/injected_ports_test.go` alone, while the fact it changed — the number of distinct failure messages the check emits — is advertised in `specs/.../spec.md` and in `add-injected-port-conformance-check.md`'s own Contracts line ("the check's four failure messages"). That ticket could not have closed SP4 no matter how well its delegate scoped by invariant, because the advertising artifact was outside its fence. The already-captured "scope by invariant, not citation" rule is necessary but not sufficient: a fence derived from the finding's citation silently caps what any amount of good scoping can reach. Right behavior: derive the fence from the Contracts field — every artifact that advertises a fact the ticket changes joins the fence, or the ticket declares an explicit successor ticket that owns the advertisement. Proposed rule: `craft-tickets` states that a Contracts entry of the form `X crosses A→B` requires both A and B in the fence or a named successor; and, because `assign` already parses both fields, a mechanical refusal at `bench spec build assign` is the codification candidate — a ticket whose Contracts names a crossing its fence cannot reach is refused at lease time rather than discovered a review round later.

- 2026-08-03 — The failure-message count was duplicated knowledge across three artifacts (spec implementation-decisions, a ticket's Contracts line, the implementation's message constants) with only the implementation as truth, and nothing mechanical graded the agreement: `bench coverage --check` reported "coverage map valid — 13 row(s)" on the exact candidate where SP4 was live, because it grades the map's shape and not its agreement with the implementation. The repo's code standard already forbids this shape ("an enforcement and its advertisement must collapse to one source") but applies it to code, not to specs and tickets. A 7-line shell loop comparing the spec's advertised count against `grep -c` of the message constants reproduces SP4 in 7ms and goes green at the promoted tree. Right behavior: a spec should not restate a count the implementation derives — it names the enumeration source instead. Proposed rule: extend the one-source-per-fact standard in `AGENTS.md` to spec and ticket prose, and let `/bench-review-implementation`'s Standards axis grade a restated count as duplicated knowledge; the reviewer decides whether any of it becomes a gate check rather than a review judgment.
