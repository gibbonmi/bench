# Harden the census heads evidence

Blocked by: none
Writes: internal/census/census.go, internal/census/census_test.go

## What to build

Three residuals from the FT271 landing close in one slice, all inside
`internal/census`.

First, the record-line layout gets one codec. The writer composes
`<timestamp>\t<head>\n` and the reader splits on `\t` and takes the second
field. A shared record/parse pair becomes the one source of that layout, and
both ends call it.

Second, the `heads{}` render becomes delimiter-safe. The line's grammar uses
`=` and `,` as its own delimiters, and `sanitize.Controls` passes both
verbatim. The render applies `sanitize.Controls` first, then escapes `,` and
`=` with a backslash.

This order is pinned. After `Controls`, every backslash
is an escape introducer. Thus a left-to-right reader parses `\\` as a literal
backslash, and `\,` or `\=` as a literal delimiter. The reverse order doubles
the inserted backslash and forges a delimiter again.

Third, `resolvedHead` stops recording a bare `KEY=VALUE` assignment text as a
head. When the pool-path-bearing simple command is only assignments, the head
comes from the next simple command in the same text that has a command word.
When the whole text has no command word, the head is the key of the
assignment whose value carries the pool path. The `=` stays and the value
drops, for example `W=`. An unrelated co-assignment's key is never the
head. The same
degenerate-head rule applies to `verbHead`'s `first` fallback (the heredoc
case): an assignment-only first command yields its degenerate key there too.

## Acceptance

- [ ] A written record for head `sed` reads back as head `sed` through the shared codec. No other census call site splits a record on `\t`.
- [ ] A record file with a foreign head `sed=9,rm` renders `census heads{}` as the one entry `sed\=9\,rm=1`, and a test pins that exact text.
- [ ] A head with a literal backslash and a delimiter, `a\,b`, renders as `a\\\,b=1`, pinning the Controls-then-delimiter order.
- [ ] A raw text `W=<pool path>; git -C "$W" rev-parse HEAD` records head `git rev-parse`.
- [ ] A raw text that is only `W=<pool path>` records head `W=`, and it renders as `W\==1`.
- [ ] A raw text `X=1 W=<pool path> Y=2` with no command word records head `W=`.
- [ ] The existing landing evidence shape stays: `census heads{sed=2,awk=1}` for plain heads is unchanged.
