# Record minimal probing as a known issue

Blocked by: none
Writes: docs/known-issues.md
Covers: none

## What to build

`docs/known-issues.md` gains one row for the cost of `bench probe`. Each mutant and each baseline of a probe without `--run` runs the whole package. A package such as `./internal/gate` takes about 70 seconds for each run, so probe sets take most of an author's time.

The workaround is minimal probing. Each ticket runs its planned probe and adds a witness probe only for an expectation that has no other demonstrated red. Every probe and every mutant names `--run` with the killing test. The author runs the whole package once, at the commit.

The retire condition names the fix: `bench probe` narrows each mutant run to the test that `--expect` names.

## Acceptance

- [ ] `docs/known-issues.md` holds the row, and each charge that names the file inherits the workaround.
- [ ] `bench test --check docs-currency-workflow` passes.
