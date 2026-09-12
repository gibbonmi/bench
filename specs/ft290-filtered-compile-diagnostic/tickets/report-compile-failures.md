# Report compile failures under a test filter

Blocked by: none
Writes: internal/testreport, CHANGELOG.md
Covers: none

## What to build

Classify a failed build before applying the no-match refusal.
Return the compiler diagnostic when --run selected a package that cannot compile.
Keep the existing no-match refusal for a successful run that starts no tests.
Use the existing report and outcome owners.
Add a Fixed entry under a dedicated Filtered compile failures changelog heading.

The coordinator reconciles the roadmap after this fix lands.

## Acceptance

- [ ] A filtered compile failure reports its compiler diagnostic at exit 1.
- [ ] The typed outcome is build-failed, with no failed tests and no test runs.
- [ ] A genuine no-match run retains its refusal.
- [ ] Matched failures and matched skips retain their existing results.
