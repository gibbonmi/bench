# Grade a new system test file as system-tagged

Blocked by: none
Writes: internal/preflight/closure.go, internal/preflight/system_tag.go (new), internal/preflight/kit_pin_new_test.go (new)
Covers: none

## What to build

The `kit-pin` row reads the system tag from the bytes of each `Writes:` file. A `(new)` test
file does not exist before the build, so build preflight cannot read its tag. The row then
stays green, and review preflight goes red after the ticket creates the file.

The gatherer grades a `(new)` Go test path that the tree does not hold from its directory. When
that directory holds a Go test file with the system tag, the new path is system-tagged. The
directory's own files are the one source, so no second list of system-suite paths exists. When
the tree holds the path, the bytes of the file stay the only source.

## Acceptance

- [ ] A ticket that writes a `(new)` test path beside system-tagged test files and does not state `BENCH_KIT` makes `kit-pin` red in `bench preflight build`. The unchanged code gave green.
- [ ] A `(new)` test path beside test files with no system tag keeps `kit-pin` green.
- [ ] A `(new)` test path that the tree holds with no system tag keeps `kit-pin` green.
- [ ] `decision.go` and `gather.go` do not grow.
