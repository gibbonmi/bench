# Review pickup: spec-authoring-discipline

Frozen base `34889eef2c0b00629a790d49307d08a58855dd48`, reviewed tip
`6722b151fccda92804902ad9cb056a5577edfe5a`. Raw findings: 12. Repair targets: 1.

## Standards

Findings: 7. Worst: two bare "map" words on lines the build moved.

- **auto-fix** — Row SAD42 and `CONTEXT.md` "Not "map" — decision map."
  `.agents/commands/bench-shape-idea.md:48` writes "moves a ready map" on a line the
  build reflowed. Repair ticket:
  `specs/spec-authoring-discipline/tickets/repair-bare-map-on-moved-lines.md`.
- **ask-user** — `.agents/commands/bench-write-spec.md:26` writes "the source map"
  and "spec-local map" in a pre-existing sentence that the build did not change. The
  build appended a new sentence to the same line. Two anchor needles in
  `internal/anchors/registry_data.go` (lines 299 and 301) pin those bytes. The build
  leaves the sentence as it is. Reviewer veto stands open.
- **ask-user** — `bench-write-spec.md` dropped "and asks one question of its own"
  from the review-round paragraph. The clause introduced the moved question, so the
  build reads it as a no-op. Reviewer veto stands open.
- **ask-user** — `internal/anchors/registry_data_test.go` adds two more hand-rolled
  `evaluate` closures on an 18-deep precedent. Parked as an idea.

## Spec

Findings: 3. Worst: the same dropped clause.

- **ask-user** — Rows SAD27 and SAD36 say "the kit" and "no kit file", but each
  forbid tuple scopes to one file. `roadmap/FT220.md:8` and `roadmap/FT257.md:10`
  still write the old terms. The anchor mechanism is per file. Reviewer decides
  whether the roadmap bodies change or the rows narrow.

## Coverage

Findings: 2. Worst: a wrapped Sources continuation that contains a colon.

- **ask-user** — A continuation line such as `and it continues: here.` holds a
  colon, so it reaches the unknown-field message and not the one-physical-line
  message. The spec defines a wrapped continuation as a line with no separator, so
  this is a spec-level gap. Parked as an idea.
