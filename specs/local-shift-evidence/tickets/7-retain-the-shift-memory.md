# 7. Retain the shift memory

Blocked by: 3-rotate-and-retain-the-record.md, 4-record-the-shift-boundaries.md
Writes: internal/otelrecord/writer.go, internal/otelrecord/retention_test.go (new), internal/otelrecord/attributes.go, internal/bounds/bounds.go, tests/canary/package-core-guard/bounds-duplicate-owner, internal/shift/session.go, internal/shift/record.go (new), internal/shift/record_test.go (new)
Covers: LE49, LE50, LE51, LE52, LE53, LE54, LE55, LE56, LE57, LE58, LE88, LE59

## What to build

Chunk: LE-B2.

Add `bounds.RecordMemoryRetained` to the bounds registry. Propose 64, as the spec's reviewer-owned size.

Add a memory store to the record package. It writes given bytes to `memory/<UTC stamp>-<trace id>.md` below the record directory, with mode 0600 in a 0700 directory. It grades each path level below the Bench home as the record writer does, so a symlinked directory is refused. It keeps the newest retained count of memory files and returns the SHA-256 digest.

On every path that reached the first iteration, the shift retains its notes before any scratch cleanup and on the retain path. It reads `.bench-notes.md` through `bounds.ClassifyNoFollow`. An absent file is `absent`. A parsed or empty file is `retained`. A symlink, a FIFO, an oversized file, or invalid UTF-8 is `refused`, and nothing is written.

The `shift` span carries `bench.memory.state`, and for a retained file `bench.memory.bytes` and `bench.memory.digest`. Declare the keys. No prompt reads the memory directory.

## Acceptance

- [ ] After a green shift whose adapter appended `MEMMARK` to the notes, one memory file holds the notes bytes.
- [ ] After a red shift that retained its worktree, one memory file holds the notes bytes.
- [ ] The span carries state `retained`, the byte count, and a digest equal to the SHA-256 of the memory file, and no record line holds `MEMMARK`.
- [ ] A notes symlink to a file that holds `SECRETMARK` records `refused`, and no memory file holds `SECRETMARK`.
- [ ] A notes FIFO records `refused`, and the shift exits within the test deadline.
- [ ] Notes past the control-record limit record `refused` and write no memory file.
- [ ] Deleted notes record `absent` and write no memory file, and empty notes record `retained` with 0 bytes.
- [ ] After the retained count plus one writes, exactly the retained count of memory files remains, and a symlinked memory directory refuses the write.
- [ ] The first adapter of a second shift reads an empty `.bench-notes.md`, and its stdin prompt holds no `MEMMARK`.
