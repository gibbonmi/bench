# Require the Feeds marker on retro improvement items

Blocked by: none
Writes: internal/retros/retros.go, internal/retros/retros_test.go, internal/retros/recommendations.go, internal/retros/recommendations_test.go, internal/conformance/registry/registry.go, internal/conformance/checks.go, internal/conformance/checks_test.go, internal/conformance/registry_test.go, internal/conformance/retro_improvement_markers_test.go, tests/canary/retro-improvement-markers/, .agents/commands/bench-final-check.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go

## What to build

A new Dev conformance check `retro-improvement-markers` over a new
`capture-retros` input source grades every item under
`## Agent-experience improvements` in each `capture/retros/*.md` file. A valid
item ends with one line that matches `^Feeds: (FT[1-9][0-9]*|new|none)$`. The
check reuses `retros.Facts` and `retros.Recommendations`. The
`/bench-final-check` retro template carries the marker and the one-sentence
change test, and one anchors-registry entry pins the template clause. Spec
group C and story 35, rows RF20 to RF25 and RF28.

## Acceptance

- [ ] An item with no `Feeds:` line yields one diagnostic naming the retro path and the item's line number.
- [ ] `Feeds: FT12`, `Feeds: new`, and `Feeds: none` pass; `Feeds: maybe` yields a diagnostic.
- [ ] An absent `capture/retros/` and an empty one each yield no diagnostic; a dangling symbolic link yields a diagnostic naming the path and its state.
- [ ] An item under a different `## ` heading in the same retro yields no diagnostic.
- [ ] The canary fixture makes the owner check red, and restoring the line clears it.
- [ ] Removing the `Feeds:` marker from the `/bench-final-check` template yields the anchor diagnostic.
