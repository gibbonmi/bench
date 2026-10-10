# Print the published commit on a successful bench commit

Blocked by: none
Writes: internal/commit/commit.go, internal/commit/landing_test.go, CHANGELOG.md
Covers: none

## What to build

A successful `bench commit` prints only `committed N path(s)`. The line does not name the
commit that the call published. Ticket authors thus run a second `git` call after each commit
to read the hash.

Change the success line to one record in the `name{key=value,...}` grammar of the verb's
exit-3 record. The record is `committed{published_commit=<hash>,path_count=<n>}`. The
`published_commit` key is the key that the exit-3 record uses. The hash is the commit that
the landing owner published. Do not change the dry-run lines or the exit-3 record.

## Acceptance

- [ ] A green `bench commit` exits 0 and its last stdout line is
  `committed{published_commit=<hash>,path_count=<n>}`, where `<hash>` is the new `HEAD`.
- [ ] `path_count` is the number of attributed paths that the commit published.
- [ ] The dry-run lines and the exit-3 `committed{published_commit=...,path=...,next=...}`
  record are unchanged.
- [ ] `CHANGELOG.md` has one bullet for the change under `## [Unreleased]`.
