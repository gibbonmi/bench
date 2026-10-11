# Give the published_commit record key one source

Blocked by: none
Writes: internal/commit/commit.go, internal/commit/record_keys.go
Covers: none

## What to build

`bench commit` prints a `committed{...}` record when it publishes a commit. The success record and the exit-3 record both name the published commit with the key `published_commit`. One constant in `internal/commit/record_keys.go` holds that key, and both records read it. The two records then cannot name the published commit with different keys.

## Acceptance

- [ ] A green `bench commit` prints `committed{published_commit=<sha>,path_count=<n>}` as its last line.
- [ ] A `bench commit` that publishes and cannot reconcile its checkout exits 3 and prints `committed{published_commit=<sha>,path=<path>,next=<route>}`.
- [ ] A change to the value of the one key constant turns a test red for each of the two records.
