# Pin every excluded category with a silent site

Blocked by: none
Writes: specs/stdlib-excluded-category-fixtures/tickets/01-pin-every-excluded-category-with-a-silent-site.md (new), internal/conformance/stdlib_reimplementation_test.go
Covers: none

## What to build

The `stdlib-reimplementation` check runs five modernize analyzers. Its doc comment
excludes "the other modernize categories". No list names that excluded set. The set is
the modernize suite less the run list, plus the exported modernize analyzers outside
the suite. Today the bites fixture holds a site for one excluded category only, the
`rangeint` site. Thus a change that adds a different excluded analyzer to the run list
keeps the tests and the kit gate green.

The bites test plants one production site for each excluded category, in one table. The
test makes sure that the check reports no diagnostic on each site. The test also runs
each site's own analyzer under the check's load mode, and makes sure that the analyzer
reports the site. Thus each silent assertion holds a real finding, and it cannot pass
because the site does not fire.

Some sites cannot fire at the kit's Go version or under the check's load mode. Each of
these sites carries a dormant reason. The test fails when a dormant site starts to fire,
so the next author removes the reason and the site pins its category. The test also
fails when a member of the modernize suite has no run-list row and no excluded site.

This ticket adds tests only. The check's run list and its function do not change.

## Acceptance

- [ ] The bites fixture holds one production site for each excluded analyzer: each suite member outside the run list, and each exported analyzer outside the suite.
- [ ] The check reports no diagnostic on each excluded site, and a failure names the analyzer of that site.
- [ ] The own analyzer of each live excluded site reports that site under the check's load mode.
- [ ] Each dormant site has a reason, and the test fails when the own analyzer of a dormant site reports it.
- [ ] The test fails when a suite member has no run-list row and no excluded site.
- [ ] For each live excluded analyzer, a probe that adds it to the run list turns the bites test red and restores the source.
- [ ] The check's run list and its function body stay unchanged.
