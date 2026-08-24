# Enforce learning journal compliance

Blocked by: none
Writes: CHANGELOG.md, capture/learnings.md, internal/learnings, internal/conformance, projects/benchkit.md, tests/canary/prose-mechanics

## What to build

The registered prose-mechanics check grades the learning journal with the existing prose parser, learning schema validator, and learning-entry parser.
The check accepts an absent or empty journal. It fails closed when a present journal has an unsupported schema, malformed content, or a refused file state.
Remove the tracked learning journal so the ignore-listed journal remains local to each checkout.

## Acceptance

- [ ] An unsupported learning journal makes the registered check report the journal path, state, and schema reason.
- [ ] Malformed headings and unaccounted content make the registered check report the journal path, source line, and parser reason.
- [ ] Invalid UTF-8, oversized content, and a wrong file type make the registered check report the journal path, state, and classifier reason.
- [ ] The retained mutation fixtures prove each planted diagnostic appears after mutation and disappears after restoration.
- [ ] Valid, empty, and absent learning journals pass the registered check.
- [ ] The existing prose-mechanics checks continue to grade the journal prose.
- [ ] The repository no longer tracks `capture/learnings.md`; the ignored local journal remains outside commits.
- [ ] The project profile describes the prose and learning-journal checks that the registered owner runs.
