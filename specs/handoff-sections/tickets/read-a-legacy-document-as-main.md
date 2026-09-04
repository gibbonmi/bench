# Read a legacy document as main

Blocked by: own-the-document-path-and-reclaim-the-lock.md
Writes: internal/handoffdoc
Covers: HS26

## What to build

Verify the premise first: `Parse` refuses a document whose first level-two
heading is `## State`, so the first run of the new verb over an existing
handoff exits 1 and names the file. Then read that legacy shape as one
`main` section. The header block stays the header. The `## State` body
becomes the section's State. The `## Next command` body, one backticked
line, becomes the section's Next value.

A `## Closed decisions` block, or any other legacy level-two block above
`## Shape`, joins the State body under a level-three heading of the same
name. The `## Shape` body stays the Shape. A second read of the rendered
result parses as the new grammar.

## Acceptance

- [ ] A legacy document with State, Closed decisions, Next command, and Shape reads as a document with one `main` section, and its render parses back.
- [ ] The legacy Next command's backticked line is the section's Next value byte for byte.
- [ ] A document that mixes a legacy `## State` with a `## main` section refuses with the file and line.
- [ ] Self-probe: drop the Next command carry-over, and report the byte-for-byte test red.
