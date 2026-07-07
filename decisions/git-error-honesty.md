# Git-Error Honesty (FT29)

## #1: What is the porcelain posture when git itself fails?

Blocked by: —
Type: Grill

### Question
`bench structure` converts git failure into false-clean: `git.Output` errors
are discarded in both the all-files and `--since` paths, so a bad ref or a
failed `ls-files` prints "no tracked source files to check" at exit 0. The
same class was fixed in `bench diff` (FT19) and re-found here — the class
outlived the instance fix.

### Answer
A porcelain command that cannot ask git the question reports the git error to
stderr and exits 1 — never an empty-but-clean answer at exit 0. The build is a
class kill, not an instance fix: audit every discarded `git.Output` (and
sibling helper) error across the porcelain packages and either handle it
meaningfully or propagate it; each audited site gets a verdict in the spec
(propagate / legitimately-optional with a comment). Advisory surfaces that
deliberately tolerate absence (e.g. status rows skipping on error) stay
tolerant, but the tolerance must be visible in the code as a decision, not a
dropped error value.

## #2: What does `bench models` do with argv?

Blocked by: —
Type: Grill

### Question
`bench models` ignores its arguments — no `-h/--help`, and unknown args emit
the full inventory at exit 0. Every sibling porcelain rejects unknown args at
exit 2.

### Answer
Match the sibling norm: unknown args → usage line at exit 2; no new flags. The
advisory exit-0 posture for *discovery results* (unreachable providers) is
unchanged — argv rejection and discovery tolerance are different contracts.

## Handoff

1. **Module boundaries.** `internal/structure` owns the loud-error fix;
   `internal/models` owns argv rejection; the audit sweeps every porcelain
   package but changes only sites where a dropped error can produce
   false-clean/false-empty output.
2. **Contracts.** `bench structure` on git failure: error naming the failed
   git operation on stderr, exit 1. `bench models <junk>`: usage on stderr,
   exit 2. Existing green-path outputs unchanged.
3. **Deep vs thin.** The audit judgment (propagate vs tolerate) is the deep
   part; each individual fix is thin error plumbing.
4. **Black-box assertables.** structure against a non-repo dir / corrupted
   ref: exit 1 + stderr; models with an unknown arg: exit 2 + usage; every
   other porcelain's existing contract stays green.
5. **Gate attachment.** Contract/runtime tests per fixed command; the audit's
   tolerate-verdicts are recorded in the spec, review-checked, not
   gate-checked.
6. **Hostile-input owners.** structure owns not-a-repo, bad `--since` ref,
   and git binary absent; models owns arbitrary argv bytes.
7. **Uncertainty flags.** How many audit sites exist is unknown until the
   sweep runs — the spec sets the rule and the delegate enumerates; if a site
   is genuinely ambiguous it comes back as a finding, not a guess.
8. **Rejected alternatives.** Fixing only structure (instance-fix repeats the
   FT19 mistake); adding flags to models.
9. **Domain watch-outs.** `bench status` is an ambient advisory surface — a
   git failure there should degrade a row, not crash the SessionStart hook;
   tolerance there is the recorded-decision case, not a defect.

Dependency order: n/a — single spec.
